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

type pulsarClientProduceInput struct {
	Topic           string   `json:"topic"`
	Messages        []string `json:"messages,omitempty"`
	NumProduce      *float64 `json:"num-produce,omitempty"`
	Rate            *float64 `json:"rate,omitempty"`
	DisableBatching bool     `json:"disable-batching,omitempty"`
	Chunking        bool     `json:"chunking,omitempty"`
	Separator       *string  `json:"separator,omitempty"`
	Properties      []string `json:"properties,omitempty"`
	Key             *string  `json:"key,omitempty"`
}

const (
	pulsarClientProduceTopicDesc = "The fully qualified topic name to produce to (format: [persistent|non-persistent]://tenant/namespace/topic). " +
		"For partitioned topics, messages will be distributed across partitions based on the partitioning scheme. " +
		"If a message key is provided, messages with the same key will go to the same partition."
	pulsarClientProduceMessagesDesc = "List of message content to send. Each array element represents one message payload. " +
		"IMPORTANT: Use this parameter to provide message content. Multiple messages can be sent in a single operation. " +
		"Each message will be sent according to the specified num-produce parameter."
	pulsarClientProduceMessageItemDesc = "Individual message content to be sent to the topic"
	pulsarClientProduceNumProduceDesc  = "Number of times to send the entire message set. " +
		"If you have 3 messages and set num-produce to 2, a total of 6 messages will be sent. (default: 1)"
	pulsarClientProduceRateDesc = "Rate limiting in messages per second. Controls the maximum speed of message production. " +
		"Set to 0 to produce messages as fast as possible without rate limiting. " +
		"Higher rates may be limited by broker capacity and network bandwidth. (default: 0)"
	pulsarClientProduceDisableBatchingDesc = "Disable message batching. When false (default), Pulsar batches multiple messages " +
		"to improve throughput and reduce network overhead. Set to true to send each message individually. " +
		"Disabling batching may reduce throughput but provides lower latency. (default: false)"
	pulsarClientProduceChunkingDesc = "Enable message chunking for large messages. When true, messages larger than " +
		"the maximum allowed size will be automatically split into smaller chunks and reassembled on consumption. " +
		"This allows sending messages that exceed broker size limits. (default: false)"
	pulsarClientProduceSeparatorDesc = "Character or string to split message content on. When specified, each message " +
		"in the messages array will be split by this separator to create additional individual messages. " +
		"Useful for sending multiple messages from a single delimited string. (default: none)"
	pulsarClientProducePropertiesDesc = "Message properties in key=value format. Properties are metadata key-value pairs " +
		"attached to messages for filtering, routing, or application-specific processing. " +
		"Example: ['priority=high', 'source=api', 'version=1.0']. Multiple properties can be specified."
	pulsarClientProducePropertyItemDesc = "Property in key=value format"
	pulsarClientProduceKeyDesc          = "Partitioning key for message routing. Messages with the same key will be sent " +
		"to the same partition in partitioned topics, ensuring ordering for related messages. " +
		"The key is also available to consumers for processing logic. Leave empty for round-robin partitioning."
)

// PulsarClientProduceToolBuilder implements the ToolBuilder interface for Pulsar Client Producer tools
// It provides functionality to build Pulsar message production tools
// /nolint:revive
type PulsarClientProduceToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarClientProduceToolBuilder creates a new Pulsar Client Producer tool builder instance
func NewPulsarClientProduceToolBuilder() *PulsarClientProduceToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_client_produce",
		Version:     "1.0.0",
		Description: "Pulsar Client message production tools",
		Category:    "pulsar_client",
		Tags:        []string{"pulsar", "produce", "client", "messaging"},
	}

	features := []string{
		"pulsar-client",
		"all",
		"all-pulsar",
	}

	return &PulsarClientProduceToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar Client Producer tool list
// This is the core method implementing the ToolBuilder interface
func (b *PulsarClientProduceToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool, err := b.buildProduceTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildProduceHandler(config.ReadOnly)

	return []builders.ToolDefinition{
		builders.ServerTool[pulsarClientProduceInput, any]{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildProduceTool builds the Pulsar Client Producer MCP tool definition
// Migrated from the original tool definition logic
func (b *PulsarClientProduceToolBuilder) buildProduceTool() (*sdk.Tool, error) {
	inputSchema, err := buildPulsarClientProduceInputSchema()
	if err != nil {
		return nil, err
	}

	toolDesc := "Produce messages to a Pulsar topic. " +
		"This tool allows you to send messages to a specified Pulsar topic with various options " +
		"to control message format, batching, rate limiting, and properties. " +
		"Pulsar supports message batching for improved throughput, chunking for large messages, " +
		"and message properties for metadata attachment. " +
		"You can send single or multiple messages, control the production rate, and add custom properties. " +
		"The tool supports message partitioning through keys and provides detailed feedback on sent messages. " +
		"Do not use this tool for Kafka protocol operations. Use 'kafka_client_produce' instead."

	return &sdk.Tool{
		Name:        "pulsar_client_produce",
		Description: toolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildProduceHandler builds the Pulsar Client Producer handler function
// Migrated from the original handler logic
func (b *PulsarClientProduceToolBuilder) buildProduceHandler(readOnly bool) builders.ToolHandlerFunc[pulsarClientProduceInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input pulsarClientProduceInput) (*sdk.CallToolResult, any, error) {
		// Check read-only mode - producing is a write operation
		if readOnly {
			return nil, nil, fmt.Errorf("message production is not allowed in read-only mode")
		}

		// Extract required parameters with validation
		topic := strings.TrimSpace(input.Topic)
		if topic == "" {
			return nil, nil, fmt.Errorf("failed to get topic: topic is required")
		}

		// Set default values and extract optional parameters
		messages := input.Messages
		if len(messages) == 0 {
			return nil, nil, fmt.Errorf("please supply message content with 'messages' parameter")
		}

		numProduce := 1
		if input.NumProduce != nil {
			numProduce = int(*input.NumProduce)
		}
		if numProduce < 1 {
			return nil, nil, fmt.Errorf("num-produce must be at least 1")
		}

		rate := 0.0
		if input.Rate != nil {
			rate = *input.Rate
		}
		if rate < 0 {
			return nil, nil, fmt.Errorf("rate must be non-negative")
		}

		disableBatching := input.DisableBatching
		chunkingAllowed := input.Chunking
		separator := ""
		if input.Separator != nil {
			separator = *input.Separator
		}
		properties := input.Properties
		key := ""
		if input.Key != nil {
			key = *input.Key
		}

		// Split messages by separator if needed
		if separator != "" && len(messages) > 0 {
			var splitMessages []string
			for _, msg := range messages {
				parts := strings.Split(msg, separator)
				for _, part := range parts {
					if part != "" {
						splitMessages = append(splitMessages, part)
					}
				}
			}
			messages = splitMessages
		}

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

		// Prepare producer options
		producerOpts := pulsar.ProducerOptions{
			Topic: topic,
		}

		// Set batching and chunking options
		if chunkingAllowed {
			producerOpts.EnableChunking = true
			producerOpts.BatchingMaxPublishDelay = 0 * time.Millisecond
		} else if disableBatching {
			producerOpts.BatchingMaxPublishDelay = 0 * time.Millisecond
		}

		producer, err := client.CreateProducer(producerOpts)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create producer: %v", err)
		}
		defer producer.Close()

		// Generate message bodies from messages
		messagePayloads, err := b.generateMessagePayloads(messages)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate message payloads: %v", err)
		}

		// Parse properties
		propMap, err := b.parseProperties(properties)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse properties: %v", err)
		}

		// Setup rate limiter
		var limiter *time.Ticker
		if rate > 0 {
			interval := time.Duration(1000/rate) * time.Millisecond
			limiter = time.NewTicker(interval)
			defer limiter.Stop()
		}

		// Send messages
		numMessagesSent := 0
		var lastMessageID pulsar.MessageID

		for range numProduce {
			for _, payload := range messagePayloads {
				// Apply rate limiting if enabled
				if limiter != nil {
					<-limiter.C
				}

				// Create message to send
				msg := &pulsar.ProducerMessage{
					Payload: payload,
				}

				// Add properties if specified
				if len(propMap) > 0 {
					msg.Properties = propMap
				}

				// Add key if specified
				if key != "" {
					msg.Key = key
				}

				// Send the message
				msgID, err := producer.Send(ctx, msg)
				if err != nil {
					return nil, nil, fmt.Errorf("failed to send message: %v", err)
				}

				lastMessageID = msgID
				numMessagesSent++
			}
		}

		// Prepare response
		response := map[string]interface{}{
			"topic":           topic,
			"messages_sent":   numMessagesSent,
			"last_message_id": fmt.Sprintf("%v", lastMessageID),
			"success":         true,
		}

		result, err := b.marshalResponse(response)
		return result, nil, err
	}
}

func buildPulsarClientProduceInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[pulsarClientProduceInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}

	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}
	setSchemaDescription(schema, "topic", pulsarClientProduceTopicDesc)
	setSchemaDescription(schema, "messages", pulsarClientProduceMessagesDesc)
	setSchemaDescription(schema, "num-produce", pulsarClientProduceNumProduceDesc)
	setSchemaDescription(schema, "rate", pulsarClientProduceRateDesc)
	setSchemaDescription(schema, "disable-batching", pulsarClientProduceDisableBatchingDesc)
	setSchemaDescription(schema, "chunking", pulsarClientProduceChunkingDesc)
	setSchemaDescription(schema, "separator", pulsarClientProduceSeparatorDesc)
	setSchemaDescription(schema, "properties", pulsarClientProducePropertiesDesc)
	setSchemaDescription(schema, "key", pulsarClientProduceKeyDesc)

	messagesSchema := schema.Properties["messages"]
	if messagesSchema != nil && messagesSchema.Items != nil {
		messagesSchema.Items.Description = pulsarClientProduceMessageItemDesc
	}

	propertiesSchema := schema.Properties["properties"]
	if propertiesSchema != nil && propertiesSchema.Items != nil {
		propertiesSchema.Items.Description = pulsarClientProducePropertyItemDesc
	}

	normalizeAdditionalProperties(schema)
	return schema, nil
}

// Unified error handling and utility functions

// handleError provides unified error handling
func (b *PulsarClientProduceToolBuilder) handleError(operation string, err error) error {
	return fmt.Errorf("failed to %s: %v", operation, err)
}

// marshalResponse provides unified JSON serialization for responses
func (b *PulsarClientProduceToolBuilder) marshalResponse(data interface{}) (*sdk.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, b.handleError("marshal response", err)
	}
	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: string(jsonBytes)}},
	}, nil
}

// generateMessagePayloads generates message payloads from message strings
func (b *PulsarClientProduceToolBuilder) generateMessagePayloads(messages []string) ([][]byte, error) {
	var payloads [][]byte

	// Add message strings
	for _, msg := range messages {
		payloads = append(payloads, []byte(msg))
	}

	return payloads, nil
}

// parseProperties parses property strings in key=value format
func (b *PulsarClientProduceToolBuilder) parseProperties(properties []string) (map[string]string, error) {
	propMap := make(map[string]string)
	for _, prop := range properties {
		parts := strings.SplitN(prop, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid property format '%s', expected key=value", prop)
		}
		propMap[parts[0]] = parts[1]
	}
	return propMap, nil
}
