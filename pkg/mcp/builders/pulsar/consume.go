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

package pulsar

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

// PulsarClientConsumeToolBuilder implements the ToolBuilder interface for Pulsar Client Consumer tools
// It provides functionality to build Pulsar message consumption tools
type PulsarClientConsumeToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarClientConsumeToolBuilder creates a new Pulsar Client Consumer tool builder instance
func NewPulsarClientConsumeToolBuilder() *PulsarClientConsumeToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_client_consume",
		Version:     "1.0.0",
		Description: "Pulsar Client message consumption tools",
		Category:    "pulsar_client",
		Tags:        []string{"pulsar", "consume", "client", "messaging"},
	}

	features := []string{
		"pulsar-client",
		"all",
		"all-pulsar",
	}

	return &PulsarClientConsumeToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar Client Consumer tool list
// This is the core method implementing the ToolBuilder interface
func (b *PulsarClientConsumeToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool := b.buildConsumeTool()
	handler := b.buildConsumeHandler(config.ReadOnly)

	return []server.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildConsumeTool builds the Pulsar Client Consumer MCP tool definition
// Migrated from the original tool definition logic
func (b *PulsarClientConsumeToolBuilder) buildConsumeTool() mcp.Tool {
	toolDesc := "Consume messages from a Pulsar topic. " +
		"This tool allows you to consume messages from a specified Pulsar topic with various options " +
		"to control the subscription behavior, message processing, and display format. " +
		"Pulsar supports multiple subscription types (Exclusive, Shared, Failover, Key_Shared) and modes " +
		"(Durable, Non-Durable) to accommodate different messaging patterns. " +
		"The tool provides comprehensive control over consumption parameters including subscription position, " +
		"timeout settings, and message display options. " +
		"Do not use this tool for Kafka protocol operations. Use 'kafka_client_consume' instead."

	return mcp.NewTool("pulsar_client_consume",
		mcp.WithDescription(toolDesc),
		mcp.WithString("topic", mcp.Required(),
			mcp.Description("The fully qualified topic name to consume from (format: [persistent|non-persistent]://tenant/namespace/topic). "+
				"For partitioned topics, you can consume from all partitions by specifying the base topic name "+
				"or from a specific partition by appending -partition-N to the topic name."),
		),
		mcp.WithString("subscription-name", mcp.Required(),
			mcp.Description("The subscription name for this consumer. "+
				"A subscription represents a named cursor for tracking message consumption progress. "+
				"Multiple consumers can share the same subscription name to form a consumer group."),
		),
		mcp.WithString("subscription-type",
			mcp.Description("Subscription type controlling message distribution among consumers:\\n"+
				"- exclusive: Only one consumer can consume from the subscription at a time\\n"+
				"- shared: Messages are distributed across all consumers in a round-robin fashion\\n"+
				"- failover: Only one active consumer, others act as backups\\n"+
				"- key_shared: Messages with the same key are delivered to the same consumer (default: exclusive)"),
		),
		mcp.WithString("subscription-mode",
			mcp.Description("Subscription durability mode:\\n"+
				"- durable: Subscription persists even when all consumers disconnect\\n"+
				"- non-durable: Subscription is deleted when all consumers disconnect (default: durable)"),
		),
		mcp.WithString("initial-position",
			mcp.Description("Initial cursor position for new subscriptions:\\n"+
				"- latest: Start consuming from the latest (most recent) message\\n"+
				"- earliest: Start consuming from the earliest (oldest available) message (default: latest)"),
		),
		mcp.WithNumber("num-messages",
			mcp.Description("Maximum number of messages to consume in this session. "+
				"Set to 0 for unlimited consumption until timeout. (default: 10)"),
		),
		mcp.WithNumber("timeout",
			mcp.Description("Maximum time to wait for messages in seconds. "+
				"The consumer will stop after this timeout even if fewer messages were received. (default: 30)"),
		),
		mcp.WithBoolean("show-properties",
			mcp.Description("Include message properties in the output. "+
				"Message properties are key-value pairs attached to messages for metadata purposes. (default: false)"),
		),
		mcp.WithBoolean("hide-payload",
			mcp.Description("Exclude message payload from the output. "+
				"Useful when you only need message metadata or are dealing with large payloads. (default: false)"),
		),
	)
}

// buildConsumeHandler builds the Pulsar Client Consumer handler function
// Migrated from the original handler logic
func (b *PulsarClientConsumeToolBuilder) buildConsumeHandler(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract required parameters with validation
		topic, err := request.RequireString("topic")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
		}

		subscriptionName, err := request.RequireString("subscription-name")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get subscription name: %v", err)), nil
		}

		// Set default values and extract optional parameters
		subscriptionType := request.GetString("subscription-type", "exclusive")
		subscriptionMode := request.GetString("subscription-mode", "durable")
		initialPosition := request.GetString("initial-position", "latest")
		numMessages := int(request.GetFloat("num-messages", 10))
		timeout := int(request.GetFloat("timeout", 30))
		showProperties := request.GetBool("show-properties", false)
		hidePayload := request.GetBool("hide-payload", false)

		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return mcp.NewToolResultError("Pulsar session not found in context"), nil
		}

		// Setup client
		client, err := session.GetPulsarClient()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create Pulsar client: %v", err)), nil
		}
		defer client.Close()

		// Prepare consumer options
		consumerOpts := pulsar.ConsumerOptions{
			Name:             "snmcp-consumer",
			Topic:            topic,
			SubscriptionName: subscriptionName,
		}

		// Set subscription type
		switch strings.ToLower(subscriptionType) {
		case "exclusive":
			consumerOpts.Type = pulsar.Exclusive
		case "shared":
			consumerOpts.Type = pulsar.Shared
		case "failover":
			consumerOpts.Type = pulsar.Failover
		case "key_shared":
			consumerOpts.Type = pulsar.KeyShared
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid subscription type: %s. Valid types: exclusive, shared, failover, key_shared", subscriptionType)), nil
		}

		// Set subscription mode
		switch strings.ToLower(subscriptionMode) {
		case "durable":
			consumerOpts.SubscriptionMode = pulsar.Durable
		case "non-durable":
			consumerOpts.SubscriptionMode = pulsar.NonDurable
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid subscription mode: %s. Valid modes: durable, non-durable", subscriptionMode)), nil
		}

		// Set initial position
		switch strings.ToLower(initialPosition) {
		case "latest":
			consumerOpts.SubscriptionInitialPosition = pulsar.SubscriptionPositionLatest
		case "earliest":
			consumerOpts.SubscriptionInitialPosition = pulsar.SubscriptionPositionEarliest
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid initial position: %s. Valid positions: latest, earliest", initialPosition)), nil
		}

		// Create consumer
		consumer, err := client.Subscribe(consumerOpts)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create consumer: %v", err)), nil
		}
		defer consumer.Close()

		// Set up timeout context
		timeoutDuration := time.Duration(timeout) * time.Second
		consumeCtx, cancelConsume := context.WithTimeout(ctx, timeoutDuration)
		defer cancelConsume()

		// Container for messages
		type MessageData struct {
			ID           string            `json:"id"`
			PublishTime  string            `json:"publish_time"`
			Properties   map[string]string `json:"properties,omitempty"`
			Key          string            `json:"key,omitempty"`
			Data         string            `json:"data,omitempty"`
			MessageCount int               `json:"message_count"`
		}

		messages := []MessageData{}
		messageCount := 0

		// Consume messages
		for {
			// Check if we've consumed the requested number of messages
			if numMessages > 0 && messageCount >= numMessages {
				break
			}

			// Receive message with timeout
			msg, err := consumer.Receive(consumeCtx)
			if err != nil {
				if err == context.DeadlineExceeded || err == context.Canceled {
					break
				}
				return mcp.NewToolResultError(fmt.Sprintf("Error receiving message: %v", err)), nil
			}

			// Process the message
			messageCount++

			// Create message data
			messageData := MessageData{
				ID:           msg.ID().String(),
				PublishTime:  msg.PublishTime().Format(time.RFC3339),
				MessageCount: messageCount,
			}

			// Add properties if requested
			if showProperties {
				messageData.Properties = msg.Properties()
			}

			// Add key if present
			if msg.Key() != "" {
				messageData.Key = msg.Key()
			}

			// Add payload unless hidden
			if !hidePayload {
				messageData.Data = string(msg.Payload())
			}

			messages = append(messages, messageData)

			// Acknowledge the message
			_ = consumer.Ack(msg)
		}

		// Prepare response
		response := map[string]interface{}{
			"topic":             topic,
			"subscription_name": subscriptionName,
			"messages_consumed": messageCount,
			"messages":          messages,
		}

		return b.marshalResponse(response)
	}
}

// Unified error handling and utility functions

// handleError provides unified error handling
func (b *PulsarClientConsumeToolBuilder) handleError(operation string, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("Failed to %s: %v", operation, err))
}

// marshalResponse provides unified JSON serialization for responses
func (b *PulsarClientConsumeToolBuilder) marshalResponse(data interface{}) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return b.handleError("marshal response", err), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}
