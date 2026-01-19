// Copyright 2025 StreamNative
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pulsar

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/google/jsonschema-go/jsonschema"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

type pulsarClientConsumeInput struct {
	Topic            string   `json:"topic"`
	SubscriptionName string   `json:"subscription-name"`
	SubscriptionType *string  `json:"subscription-type,omitempty"`
	SubscriptionMode *string  `json:"subscription-mode,omitempty"`
	InitialPosition  *string  `json:"initial-position,omitempty"`
	NumMessages      *float64 `json:"num-messages,omitempty"`
	Timeout          *float64 `json:"timeout,omitempty"`
	ShowProperties   bool     `json:"show-properties,omitempty"`
	HidePayload      bool     `json:"hide-payload,omitempty"`
}

const (
	pulsarClientConsumeTopicDesc = "The fully qualified topic name to consume from (format: [persistent|non-persistent]://tenant/namespace/topic). " +
		"For partitioned topics, you can consume from all partitions by specifying the base topic name " +
		"or from a specific partition by appending -partition-N to the topic name."
	pulsarClientConsumeSubscriptionNameDesc = "The subscription name for this consumer. " +
		"A subscription represents a named cursor for tracking message consumption progress. " +
		"Multiple consumers can share the same subscription name to form a consumer group."
	pulsarClientConsumeSubscriptionTypeDesc = "Subscription type controlling message distribution among consumers:\n" +
		"- exclusive: Only one consumer can consume from the subscription at a time\n" +
		"- shared: Messages are distributed across all consumers in a round-robin fashion\n" +
		"- failover: Only one active consumer, others act as backups\n" +
		"- key_shared: Messages with the same key are delivered to the same consumer (default: exclusive)"
	pulsarClientConsumeSubscriptionModeDesc = "Subscription durability mode:\n" +
		"- durable: Subscription persists even when all consumers disconnect\n" +
		"- non-durable: Subscription is deleted when all consumers disconnect (default: durable)"
	pulsarClientConsumeInitialPositionDesc = "Initial cursor position for new subscriptions:\n" +
		"- latest: Start consuming from the latest (most recent) message\n" +
		"- earliest: Start consuming from the earliest (oldest available) message (default: latest)"
	pulsarClientConsumeNumMessagesDesc = "Maximum number of messages to consume in this session. " +
		"Set to 0 for unlimited consumption until timeout. (default: 10)"
	pulsarClientConsumeTimeoutDesc = "Maximum time to wait for messages in seconds. " +
		"The consumer will stop after this timeout even if fewer messages were received. (default: 30)"
	pulsarClientConsumeShowPropertiesDesc = "Include message properties in the output. " +
		"Message properties are key-value pairs attached to messages for metadata purposes. (default: false)"
	pulsarClientConsumeHidePayloadDesc = "Exclude message payload from the output. " +
		"Useful when you only need message metadata or are dealing with large payloads. (default: false)"
)

// PulsarClientConsumeToolBuilder implements the ToolBuilder interface for Pulsar Client Consumer tools
// It provides functionality to build Pulsar message consumption tools
// /nolint:revive
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
func (b *PulsarClientConsumeToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool, err := b.buildConsumeTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildConsumeHandler(config.ReadOnly)

	return []builders.ToolDefinition{
		builders.ServerTool[pulsarClientConsumeInput, any]{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildConsumeTool builds the Pulsar Client Consumer MCP tool definition
// Migrated from the original tool definition logic
func (b *PulsarClientConsumeToolBuilder) buildConsumeTool() (*sdk.Tool, error) {
	inputSchema, err := buildPulsarClientConsumeInputSchema()
	if err != nil {
		return nil, err
	}

	toolDesc := "Consume messages from a Pulsar topic. " +
		"This tool allows you to consume messages from a specified Pulsar topic with various options " +
		"to control the subscription behavior, message processing, and display format. " +
		"Pulsar supports multiple subscription types (Exclusive, Shared, Failover, Key_Shared) and modes " +
		"(Durable, Non-Durable) to accommodate different messaging patterns. " +
		"The tool provides comprehensive control over consumption parameters including subscription position, " +
		"timeout settings, and message display options. " +
		"Do not use this tool for Kafka protocol operations. Use 'kafka_client_consume' instead."

	return &sdk.Tool{
		Name:        "pulsar_client_consume",
		Description: toolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildConsumeHandler builds the Pulsar Client Consumer handler function
// Migrated from the original handler logic
func (b *PulsarClientConsumeToolBuilder) buildConsumeHandler(_ bool) builders.ToolHandlerFunc[pulsarClientConsumeInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input pulsarClientConsumeInput) (*sdk.CallToolResult, any, error) {
		// Extract required parameters with validation
		topic := strings.TrimSpace(input.Topic)
		if topic == "" {
			return nil, nil, fmt.Errorf("failed to get topic: topic is required")
		}

		subscriptionName := strings.TrimSpace(input.SubscriptionName)
		if subscriptionName == "" {
			return nil, nil, fmt.Errorf("failed to get subscription name: subscription-name is required")
		}

		// Set default values and extract optional parameters
		subscriptionType := "exclusive"
		if input.SubscriptionType != nil {
			value := strings.TrimSpace(*input.SubscriptionType)
			if value != "" {
				subscriptionType = value
			}
		}

		subscriptionMode := "durable"
		if input.SubscriptionMode != nil {
			value := strings.TrimSpace(*input.SubscriptionMode)
			if value != "" {
				subscriptionMode = value
			}
		}

		initialPosition := "latest"
		if input.InitialPosition != nil {
			value := strings.TrimSpace(*input.InitialPosition)
			if value != "" {
				initialPosition = value
			}
		}

		numMessages := 10
		if input.NumMessages != nil {
			numMessages = int(*input.NumMessages)
		}

		timeout := 30
		if input.Timeout != nil {
			timeout = int(*input.Timeout)
		}

		showProperties := input.ShowProperties
		hidePayload := input.HidePayload

		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return nil, nil, fmt.Errorf("pulsar session not found in context")
		}

		// Setup client
		client, err := session.GetPulsarClient()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create Pulsar client: %v", err)
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
			return nil, nil, fmt.Errorf("invalid subscription type: %s. Valid types: exclusive, shared, failover, key_shared", subscriptionType)
		}

		// Set subscription mode
		switch strings.ToLower(subscriptionMode) {
		case "durable":
			consumerOpts.SubscriptionMode = pulsar.Durable
		case "non-durable":
			consumerOpts.SubscriptionMode = pulsar.NonDurable
		default:
			return nil, nil, fmt.Errorf("invalid subscription mode: %s. Valid modes: durable, non-durable", subscriptionMode)
		}

		// Set initial position
		switch strings.ToLower(initialPosition) {
		case "latest":
			consumerOpts.SubscriptionInitialPosition = pulsar.SubscriptionPositionLatest
		case "earliest":
			consumerOpts.SubscriptionInitialPosition = pulsar.SubscriptionPositionEarliest
		default:
			return nil, nil, fmt.Errorf("invalid initial position: %s. Valid positions: latest, earliest", initialPosition)
		}

		// Create consumer
		consumer, err := client.Subscribe(consumerOpts)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create consumer: %v", err)
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
		for numMessages <= 0 || messageCount < numMessages {
			// Receive message with timeout
			msg, err := consumer.Receive(consumeCtx)
			if err != nil {
				if err == context.DeadlineExceeded || err == context.Canceled {
					break
				}
				return nil, nil, fmt.Errorf("error receiving message: %v", err)
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

		result, err := b.marshalResponse(response)
		return result, nil, err
	}
}

func buildPulsarClientConsumeInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[pulsarClientConsumeInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}

	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}
	setSchemaDescription(schema, "topic", pulsarClientConsumeTopicDesc)
	setSchemaDescription(schema, "subscription-name", pulsarClientConsumeSubscriptionNameDesc)
	setSchemaDescription(schema, "subscription-type", pulsarClientConsumeSubscriptionTypeDesc)
	setSchemaDescription(schema, "subscription-mode", pulsarClientConsumeSubscriptionModeDesc)
	setSchemaDescription(schema, "initial-position", pulsarClientConsumeInitialPositionDesc)
	setSchemaDescription(schema, "num-messages", pulsarClientConsumeNumMessagesDesc)
	setSchemaDescription(schema, "timeout", pulsarClientConsumeTimeoutDesc)
	setSchemaDescription(schema, "show-properties", pulsarClientConsumeShowPropertiesDesc)
	setSchemaDescription(schema, "hide-payload", pulsarClientConsumeHidePayloadDesc)

	normalizeAdditionalProperties(schema)
	return schema, nil
}

// Unified error handling and utility functions

// handleError provides unified error handling
func (b *PulsarClientConsumeToolBuilder) handleError(operation string, err error) error {
	return fmt.Errorf("failed to %s: %v", operation, err)
}

// marshalResponse provides unified JSON serialization for responses
func (b *PulsarClientConsumeToolBuilder) marshalResponse(data interface{}) (*sdk.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, b.handleError("marshal response", err)
	}
	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: string(jsonBytes)}},
	}, nil
}
