// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package pftools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/mark3labs/mcp-go/mcp"
)

// FunctionInvoker handles function invocation and result tracking
type FunctionInvoker struct {
	client         pulsar.Client
	resultChannels map[string]chan FunctionResult
	mutex          sync.RWMutex
}

// FunctionResult represents the result of a function invocation
type FunctionResult struct {
	Data  interface{}
	Error error
}

// NewFunctionInvoker creates a new FunctionInvoker
func NewFunctionInvoker(client pulsar.Client) *FunctionInvoker {
	return &FunctionInvoker{
		client:         client,
		resultChannels: make(map[string]chan FunctionResult),
		mutex:          sync.RWMutex{},
	}
}

// InvokeFunctionAndWait sends a message to the function and waits for the result
func (fi *FunctionInvoker) InvokeFunctionAndWait(ctx context.Context, fnTool *FunctionTool, params map[string]interface{}) (*mcp.CallToolResult, error) {
	// Convert params to JSON
	jsonData, err := json.Marshal(params)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal params to JSON: %v", err)), nil
	}

	// Create a unique message ID for tracking the request
	messageID := generateMessageID(fnTool.Name)

	// Create a result channel for this request
	resultChan := make(chan FunctionResult, 1)
	fi.registerResultChannel(messageID, resultChan)
	defer fi.unregisterResultChannel(messageID)

	// Set up consumer for output topic
	err = fi.setupConsumer(ctx, fnTool.OutputTopic, messageID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set up consumer: %v", err)), nil
	}

	// Send message to input topic
	err = fi.sendMessage(ctx, fnTool.InputTopic, messageID, string(jsonData))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to send message: %v", err)), nil
	}

	// Wait for result or timeout
	select {
	case result := <-resultChan:
		if result.Error != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Function execution failed: %v", result.Error)), nil
		}

		// Convert result to JSON
		jsonResult, err := json.Marshal(result.Data)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal result to JSON: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonResult)), nil
	case <-ctx.Done():
		return mcp.NewToolResultError(fmt.Sprintf("Function invocation timed out after %v", ctx.Value("timeout"))), nil
	}
}

// setupConsumer creates a consumer for the output topic
func (fi *FunctionInvoker) setupConsumer(ctx context.Context, outputTopic, messageID string) error {
	consumerOptions := pulsar.ConsumerOptions{
		Topic:            outputTopic,
		SubscriptionName: fmt.Sprintf("mcp-tool-consumer-%s", messageID),
		Type:             pulsar.Exclusive,
	}

	// Create the consumer
	consumer, err := fi.client.Subscribe(consumerOptions)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	// Start goroutine to receive messages
	go func() {
		// Ensure we close the consumer when done
		defer func() {
			consumer.Close()
		}()

		// Get messages with a timeout
		for {
			select {
			case <-ctx.Done():
				// Context canceled, exit
				return
			default:
				// Set a short timeout for Receive to make it responsive to context cancellation
				receiveCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
				msg, err := consumer.Receive(receiveCtx)
				cancel()

				if err != nil {
					if err == context.DeadlineExceeded || err == context.Canceled {
						// Timeout or cancellation, but keep trying unless the parent context is done
						continue
					}

					log.Printf("Error receiving message: %v", err)
					continue
				}

				// Process the message
				fi.processMessage(msg, messageID)

				// Acknowledge the message
				consumer.Ack(msg)

				// Stop after processing one message
				return
			}
		}
	}()

	return nil
}

// sendMessage sends a message to the input topic
func (fi *FunctionInvoker) sendMessage(ctx context.Context, inputTopic, messageID, payload string) error {
	// Create a producer
	producer, err := fi.client.CreateProducer(pulsar.ProducerOptions{
		Topic: inputTopic,
	})
	if err != nil {
		return fmt.Errorf("failed to create producer: %w", err)
	}
	defer producer.Close()

	// Create message properties with correlation ID
	properties := map[string]string{
		"correlationId": messageID,
	}

	// Send the message with properties
	_, err = producer.Send(ctx, &pulsar.ProducerMessage{
		Payload:    []byte(payload),
		Properties: properties,
	})

	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// processMessage processes a message received from the output topic
func (fi *FunctionInvoker) processMessage(msg pulsar.Message, messageID string) {
	// Check if the message has our correlation ID
	if msg.Properties()["correlationId"] != messageID {
		// Not our message, ignore
		return
	}

	// Parse the message payload
	var result interface{}
	err := json.Unmarshal(msg.Payload(), &result)
	if err != nil {
		// If JSON parsing fails, use the raw payload as a string
		result = string(msg.Payload())
	}

	// Get the result channel
	fi.mutex.RLock()
	resultChan, exists := fi.resultChannels[messageID]
	fi.mutex.RUnlock()

	if exists {
		// Send the result to the channel
		resultChan <- FunctionResult{
			Data:  result,
			Error: nil,
		}
	}
}

// registerResultChannel registers a result channel for a message ID
func (fi *FunctionInvoker) registerResultChannel(messageID string, resultChan chan FunctionResult) {
	fi.mutex.Lock()
	defer fi.mutex.Unlock()
	fi.resultChannels[messageID] = resultChan
}

// unregisterResultChannel unregisters a result channel for a message ID
func (fi *FunctionInvoker) unregisterResultChannel(messageID string) {
	fi.mutex.Lock()
	defer fi.mutex.Unlock()
	delete(fi.resultChannels, messageID)
}

// generateMessageID generates a unique message ID
func generateMessageID(functionName string) string {
	return fmt.Sprintf("%s-%d", functionName, time.Now().UnixNano())
}
