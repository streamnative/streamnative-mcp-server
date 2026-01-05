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
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/adapter"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
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
func (b *PulsarClientProduceToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ServerTool, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool := b.buildProduceTool()
	handler := b.buildProduceHandler(config.ReadOnly)

	return []builders.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildProduceTool builds the Pulsar Client Producer MCP tool definition
// Migrated from the original tool definition logic
func (b *PulsarClientProduceToolBuilder) buildProduceTool() *mcpsdk.Tool {
	toolDesc := "Produce messages to a Pulsar topic. " +
		"This tool allows you to send messages to a specified Pulsar topic with various options " +
		"to control message format, batching, rate limiting, and properties. " +
		"Pulsar supports message batching for improved throughput, chunking for large messages, " +
		"and message properties for metadata attachment. " +
		"You can send single or multiple messages, control the production rate, and add custom properties. " +
		"The tool supports message partitioning through keys and provides detailed feedback on sent messages. " +
		"Do not use this tool for Kafka protocol operations. Use 'kafka_client_produce' instead."

	return builders.NewTool("pulsar_client_produce",
		builders.WithDescription(toolDesc),
		builders.WithString("topic", builders.Required(),
			builders.Description("The fully qualified topic name to produce to (format: [persistent|non-persistent]://tenant/namespace/topic). "+
				"For partitioned topics, messages will be distributed across partitions based on the partitioning scheme. "+
				"If a message key is provided, messages with the same key will go to the same partition."),
		),
		builders.WithArray("messages",
			builders.Description("List of message content to send. Each array element represents one message payload. "+
				"IMPORTANT: Use this parameter to provide message content. Multiple messages can be sent in a single operation. "+
				"Each message will be sent according to the specified num-produce parameter."),
			builders.Items(
				map[string]interface{}{
					"type":        "string",
					"description": "Individual message content to be sent to the topic",
				},
			),
		),
		builders.WithNumber("num-produce",
			builders.Description("Number of times to send the entire message set. "+
				"If you have 3 messages and set num-produce to 2, a total of 6 messages will be sent. (default: 1)"),
		),
		builders.WithNumber("rate",
			builders.Description("Rate limiting in messages per second. Controls the maximum speed of message production. "+
				"Set to 0 to produce messages as fast as possible without rate limiting. "+
				"Higher rates may be limited by broker capacity and network bandwidth. (default: 0)"),
		),
		builders.WithBoolean("disable-batching",
			builders.Description("Disable message batching. When false (default), Pulsar batches multiple messages "+
				"to improve throughput and reduce network overhead. Set to true to send each message individually. "+
				"Disabling batching may reduce throughput but provides lower latency. (default: false)"),
		),
		builders.WithBoolean("chunking",
			builders.Description("Enable message chunking for large messages. When true, messages larger than "+
				"the maximum allowed size will be automatically split into smaller chunks and reassembled on consumption. "+
				"This allows sending messages that exceed broker size limits. (default: false)"),
		),
		builders.WithString("separator",
			builders.Description("Character or string to split message content on. When specified, each message "+
				"in the messages array will be split by this separator to create additional individual messages. "+
				"Useful for sending multiple messages from a single delimited string. (default: none)"),
		),
		builders.WithArray("properties",
			builders.Description("Message properties in key=value format. Properties are metadata key-value pairs "+
				"attached to messages for filtering, routing, or application-specific processing. "+
				"Example: ['priority=high', 'source=api', 'version=1.0']. Multiple properties can be specified."),
			builders.Items(
				map[string]interface{}{
					"type":        "string",
					"description": "Property in key=value format",
				},
			),
		),
		builders.WithString("key",
			builders.Description("Partitioning key for message routing. Messages with the same key will be sent "+
				"to the same partition in partitioned topics, ensuring ordering for related messages. "+
				"The key is also available to consumers for processing logic. Leave empty for round-robin partitioning."),
		),
	)
}

// buildProduceHandler builds the Pulsar Client Producer handler function
// Migrated from the original handler logic
func (b *PulsarClientProduceToolBuilder) buildProduceHandler(readOnly bool) func(context.Context, *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	return func(ctx context.Context, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
		// Check read-only mode - producing is a write operation
		if readOnly {
			return adapter.NewErrorResult("Message production is not allowed in read-only mode"), nil
		}

		// Extract required parameters with validation
		topic, err := adapter.RequireString(request, "topic")
		if err != nil {
			return adapter.NewErrorResult("Failed to get topic: %v", err), nil
		}

		// Set default values and extract optional parameters
		messages := adapter.GetStringSlice(request, "messages", []string{})
		if len(messages) == 0 {
			return adapter.NewErrorResult("Please supply message content with 'messages' parameter"), nil
		}

		numProduce := int(adapter.GetFloat(request, "num-produce", 1))
		if numProduce < 1 {
			return adapter.NewErrorResult("num-produce must be at least 1"), nil
		}

		rate := adapter.GetFloat(request, "rate", 0)
		if rate < 0 {
			return adapter.NewErrorResult("rate must be non-negative"), nil
		}

		disableBatching := adapter.GetBool(request, "disable-batching", false)
		chunkingAllowed := adapter.GetBool(request, "chunking", false)
		separator := adapter.GetString(request, "separator", "")
		properties := adapter.GetStringSlice(request, "properties", []string{})
		key := adapter.GetString(request, "key", "")

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
			return adapter.NewErrorResult("Pulsar session not found in context"), nil
		}

		// Setup client
		client, err := session.GetPulsarClient()
		if err != nil {
			return adapter.NewErrorResult("Failed to create Pulsar client: %v", err), nil
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
			return adapter.NewErrorResult("Failed to create producer: %v", err), nil
		}
		defer producer.Close()

		// Generate message bodies from messages
		messagePayloads, err := b.generateMessagePayloads(messages)
		if err != nil {
			return adapter.NewErrorResult("Failed to generate message payloads: %v", err), nil
		}

		// Parse properties
		propMap, err := b.parseProperties(properties)
		if err != nil {
			return adapter.NewErrorResult("Failed to parse properties: %v", err), nil
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
					return adapter.NewErrorResult("Failed to send message: %v", err), nil
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

		return b.marshalResponse(response)
	}
}

// Unified error handling and utility functions

// handleError provides unified error handling
func (b *PulsarClientProduceToolBuilder) handleError(operation string, err error) *mcpsdk.CallToolResult {
	return adapter.NewErrorResult("Failed to %s: %v", operation, err)
}

// marshalResponse provides unified JSON serialization for responses
func (b *PulsarClientProduceToolBuilder) marshalResponse(data interface{}) (*mcpsdk.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return b.handleError("marshal response", err), nil
	}
	return adapter.NewTextResult(string(jsonBytes)), nil
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
