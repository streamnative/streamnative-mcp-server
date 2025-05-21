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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	cliutils "github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/schema"
)

// FunctionInvoker handles function invocation and result tracking
type FunctionInvoker struct {
	client         pulsar.Client
	resultChannels map[string]chan FunctionResult
	mutex          sync.RWMutex
	manager        *PulsarFunctionManager
}

// FunctionResult represents the result of a function invocation
type FunctionResult struct {
	Data  string
	Error error
}

// NewFunctionInvoker creates a new FunctionInvoker
func NewFunctionInvoker(manager *PulsarFunctionManager) *FunctionInvoker {
	return &FunctionInvoker{
		client:         manager.pulsarClient,
		resultChannels: make(map[string]chan FunctionResult),
		mutex:          sync.RWMutex{},
		manager:        manager,
	}
}

// InvokeFunctionAndWait sends a message to the function and waits for the result
func (fi *FunctionInvoker) InvokeFunctionAndWait(ctx context.Context, fnTool *FunctionTool, params map[string]interface{}) (*mcp.CallToolResult, error) {
	schemaConverter, err := schema.ConverterFactory(fnTool.OutputSchema.Type)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get schema converter: %v", err)), nil
	}

	payload, err := schemaConverter.SerializeMCPRequestToPulsarPayload(params, fnTool.OutputSchema.PulsarSchemaInfo)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize payload: %v", err)), nil
	}

	// Create a result channel for this request
	resultChan := make(chan FunctionResult, 1)

	// Send message to input topic
	msgID, err := fi.sendMessage(ctx, fnTool.InputTopic, payload)
	if err != nil || msgID == "" {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to send message: %v", err)), nil
	}

	fi.registerResultChannel(msgID, resultChan)
	defer fi.unregisterResultChannel(msgID)

	// Set up consumer for output topic
	err = fi.setupConsumer(ctx, fnTool.InputTopic, fnTool.OutputTopic, msgID, fnTool.OutputSchema)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set up consumer: %v", err)), nil
	}

	// Wait for result or timeout
	select {
	case result := <-resultChan:
		if result.Error != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Function execution failed: %v", result.Error)), nil
		}

		return mcp.NewToolResultText(result.Data), nil
	case <-ctx.Done():
		return mcp.NewToolResultError(fmt.Sprintf("Function invocation timed out after %v", ctx.Value("timeout"))), nil
	}
}

// setupConsumer creates a consumer for the output topic
func (fi *FunctionInvoker) setupConsumer(ctx context.Context, inputTopic, outputTopic, messageID string, schema *SchemaInfo) error {
	consumerOptions := pulsar.ConsumerOptions{
		Topic:                       outputTopic,
		SubscriptionName:            fmt.Sprintf("mcp-tool-consumer-%s", messageID),
		Type:                        pulsar.Exclusive,
		SubscriptionInitialPosition: pulsar.SubscriptionPositionEarliest,
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
				defer cancel()

				if err != nil {
					if err == context.DeadlineExceeded || err == context.Canceled {
						// Timeout or cancellation, but keep trying unless the parent context is done
						continue
					}

					continue
				}

				// Process the message
				err = fi.processMessage(inputTopic, msg, messageID, schema)
				if err != nil {
					if err == ErrNotOurMessage {
						_ = consumer.Ack(msg)
						continue
					}
					continue
				}

				// Acknowledge the message
				_ = consumer.Ack(msg)

				// Stop after processing one message
				return
			}
		}
	}()

	return nil
}

// sendMessage sends a message to the input topic
func (fi *FunctionInvoker) sendMessage(ctx context.Context, inputTopic string, payload []byte) (string, error) {
	producer, err := fi.manager.GetProducer(inputTopic)
	if err != nil {
		return "", fmt.Errorf("failed to get producer for topic %s: %w", inputTopic, err)
	}

	// Send the message with properties
	msgID, err := producer.Send(ctx, &pulsar.ProducerMessage{
		Payload: payload,
	})

	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	return msgID.String(), nil
}

// processMessage processes a message received from the output topic
func (fi *FunctionInvoker) processMessage(inputTopic string, msg pulsar.Message, messageID string, schema *SchemaInfo) error {
	// Check if the message has our correlation ID
	correlationIDbytes, err := base64.StdEncoding.DecodeString(msg.Properties()["__pfn_input_msg_id__"])
	if err != nil {
		return fmt.Errorf("failed to decode correlation ID: %w", err)
	}
	correlationID, err := pulsar.DeserializeMessageID(correlationIDbytes)
	if err != nil {
		return fmt.Errorf("failed to deserialize correlation ID: %w", err)
	}
	correlationInputTopic := msg.Properties()["__pfn_input_topic__"]

	if !isCorrelationInputTopic(correlationInputTopic, inputTopic) {
		// Not our message, ignore
		return ErrNotOurMessage
	}
	if correlationID.String() != messageID {
		// Not our message, ignore
		return ErrNotOurMessage
	}

	// Get the result channel
	fi.mutex.RLock()
	resultChan, exists := fi.resultChannels[messageID]
	fi.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("result channel not found for message ID: %s", messageID)
	}

	switch schema.Type {
	case "STRING":
		result := string(msg.Payload())
		// Send the result to the channel
		resultChan <- FunctionResult{
			Data:  result,
			Error: nil,
		}
	case "JSON":
		var result map[string]interface{}
		err := json.Unmarshal(msg.Payload(), &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal message payload: %w", err)
		}
		resultString, err := json.Marshal(result)
		if err != nil {
			return fmt.Errorf("failed to marshal result to JSON: %w", err)
		}
		resultChan <- FunctionResult{
			Data:  string(resultString),
			Error: nil,
		}
	default:
		return fmt.Errorf("unsupported schema type: %s", schema.Type)
	}

	return nil
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

func isCorrelationInputTopic(correlationInputTopic string, inputTopic string) bool {
	// remove the partition index from the input topic
	if strings.Contains(correlationInputTopic, cliutils.PARTITIONEDTOPICSUFFIX) {
		correlationInputTopic = strings.Split(correlationInputTopic, cliutils.PARTITIONEDTOPICSUFFIX)[0]
	}

	// remove the partition index from the input topic
	if strings.Contains(inputTopic, cliutils.PARTITIONEDTOPICSUFFIX) {
		inputTopic = strings.Split(inputTopic, cliutils.PARTITIONEDTOPICSUFFIX)[0]
	}

	correlationInputTopicName, err := cliutils.GetTopicName(correlationInputTopic)
	if err != nil {
		return false
	}
	inputTopicName, err := cliutils.GetTopicName(inputTopic)
	if err != nil {
		return false
	}

	return correlationInputTopicName.String() == inputTopicName.String()
}
