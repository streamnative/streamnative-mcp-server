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

package kafka

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/hamba/avro/v2"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/adapter"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
)

// KafkaProduceToolBuilder implements the ToolBuilder interface for Kafka client produce operations
// It provides functionality to build Kafka producer tools
// /nolint:revive
type KafkaProduceToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewKafkaProduceToolBuilder creates a new Kafka produce tool builder instance
func NewKafkaProduceToolBuilder() *KafkaProduceToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "kafka_produce",
		Version:     "1.0.0",
		Description: "Kafka client produce tools",
		Category:    "kafka_client",
		Tags:        []string{"kafka", "client", "produce"},
	}

	features := []string{
		"kafka-client",
		"all",
		"all-kafka",
	}

	return &KafkaProduceToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Kafka produce tool list
// This is the core method implementing the ToolBuilder interface
func (b *KafkaProduceToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ServerTool, error) {
	// Skip registration if in read-only mode
	if config.ReadOnly {
		return nil, nil
	}

	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool := b.buildKafkaProduceTool()
	handler := b.buildKafkaProduceHandler()

	return []builders.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildKafkaProduceTool builds the Kafka produce MCP tool definition
// Migrated from the original tool definition logic
func (b *KafkaProduceToolBuilder) buildKafkaProduceTool() *mcpsdk.Tool {
	toolDesc := "Produce messages to a Kafka topic.\n" +
		"This tool allows you to send messages to Kafka topics with various options for message creation.\n\n" +
		"Kafka Producer Concepts:\n" +
		"- Producers write data to Kafka topics, which can be spread across multiple partitions\n" +
		"- Messages can include a key, which determines the partition assignment (consistent hashing)\n" +
		"- Headers can be added to messages to include metadata without affecting the message payload\n" +
		"- Partitions enable parallel processing and ordered delivery within a single partition\n\n" +
		"This tool provides a simple producer instance for diagnostic and testing purposes. Do not use this tool for Pulsar protocol operations. Use 'pulsar_client_produce' instead.\n\n" +
		"Usage Examples:\n\n" +
		"1. Basic message production - Send a simple message to a topic:\n" +
		"   topic: \"my-topic\"\n" +
		"   value: \"Hello, Kafka!\"\n\n" +
		"2. Keyed message - Send a message with a key for consistent partition routing:\n" +
		"   topic: \"my-topic\"\n" +
		"   key: \"user-123\"\n" +
		"   value: \"User activity data\"\n\n" +
		"3. Multiple messages - Send several messages in one request:\n" +
		"   topic: \"my-topic\"\n" +
		"   messages: [{\"key\": \"key1\", \"value\": \"value1\"}, {\"key\": \"key2\", \"value\": \"value2\"}]\n\n" +
		"4. Message with headers - Include metadata with your message:\n" +
		"   topic: \"my-topic\"\n" +
		"   value: \"Message with headers\"\n" +
		"   headers: [\"source=mcp-tool\", \"timestamp=2023-06-01\"]\n\n" +
		"5. Specific partition - Send to a particular partition:\n" +
		"   topic: \"my-topic\"\n" +
		"   value: \"Targeted message\"\n" +
		"   partition: 2\n\n" +
		"This tool requires Kafka producer permissions on the specified topic."

	return builders.NewTool("kafka_client_produce",
		builders.WithDescription(toolDesc),
		builders.WithString("topic", builders.Required(),
			builders.Description("The name of the Kafka topic to produce messages to. "+
				"Must be an existing topic that the user has write permissions for."),
		),
		builders.WithString("key",
			builders.Description("The key for the message. "+
				"Optional. Keys are used for partition assignment and maintaining order for related messages. "+
				"Messages with the same key will be sent to the same partition."),
		),
		builders.WithString("value",
			builders.Required(),
			builders.Description("The value/content of the message to send. "+
				"This is the actual payload that will be delivered to consumers. It can be a JSON string, and the system will automatically serialize it to the appropriate format based on the schema registry if it is available."),
		),
		builders.WithArray("headers",
			builders.Description("Message headers in the format of [\"key=value\"]. "+
				"Optional. Headers allow you to attach metadata to messages without modifying the payload. "+
				"They are passed along with the message to consumers."),
			builders.Items(map[string]interface{}{
				"type":        "string",
				"description": "key value pair in the format of \"key=value\"",
			}),
		),
		builders.WithNumber("partition",
			builders.Description("The specific partition to send the message to. "+
				"Optional. If not specified, Kafka will automatically assign a partition based on the message key (if provided) or round-robin assignment. "+
				"Specifying a partition can be useful for testing or when you need guaranteed partition assignment."),
		),
		builders.WithArray("messages",
			builders.Description("An array of messages to send in batch. "+
				"Optional. Alternative to the single message parameters (key, value, headers, partition). "+
				"Each message object can contain 'key', 'value', 'headers', and 'partition' properties. "+
				"Batch sending is more efficient for multiple messages."),
			builders.Items(map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"key": map[string]interface{}{
						"type":        "string",
						"description": "Message key",
					},
					"value": map[string]interface{}{
						"type":        "string",
						"description": "Message value (required)",
					},
					"headers": map[string]interface{}{
						"type":        "array",
						"description": "Message headers as array of \"key=value\" strings",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
					"partition": map[string]interface{}{
						"type":        "number",
						"description": "Target partition number",
					},
				},
				"required": []string{"value"},
			}),
		),
		builders.WithBoolean("sync",
			builders.Description("Whether to wait for server acknowledgment before returning. "+
				"Optional. Default is true. When true, ensures the message was successfully written "+
				"to the topic before the tool returns a success response."),
		),
	)
}

// buildKafkaProduceHandler builds the Kafka produce handler function
// Migrated from the original handler logic
func (b *KafkaProduceToolBuilder) buildKafkaProduceHandler() mcpsdk.ToolHandler {
	return func(ctx context.Context, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
		// Get required topic parameter
		topicName, err := adapter.RequireString(request, "topic")
		if err != nil {
			return b.handleError("get topic name", err), nil
		}

		// Get Kafka session from context
		session := mcpCtx.GetKafkaSession(ctx)
		if session == nil {
			return b.handleError("get Kafka session not found in context", nil), nil
		}

		// Create Kafka client using the session
		kafkaClient, err := session.GetClient()
		if err != nil {
			return b.handleError("create Kafka client", err), nil
		}
		defer kafkaClient.Close()

		srClient, err := session.GetSchemaRegistryClient()
		schemaReady := false
		var serde sr.Serde
		if err == nil && srClient != nil {
			schemaReady = true
		}

		// Set timeout
		timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		if err = kafkaClient.Ping(timeoutCtx); err != nil { // check connectivity to cluster
			return b.handleError("ping Kafka cluster", err), nil
		}

		if schemaReady {
			subjSchema, err := srClient.SchemaByVersion(timeoutCtx, topicName+"-value", -1)
			if err != nil {
				return b.handleError("get schema", err), nil
			}
			switch subjSchema.Type {
			case sr.TypeAvro:
				avroSchema, err := avro.Parse(subjSchema.Schema.Schema)
				if err != nil {
					return b.handleError("parse avro schema", err), nil
				}
				serde.Register(
					subjSchema.ID,
					map[string]any{},
					sr.EncodeFn(func(v any) ([]byte, error) {
						return avro.Marshal(avroSchema, v)
					}),
					sr.DecodeFn(func(data []byte, v any) error {
						return avro.Unmarshal(avroSchema, data, v)
					}),
				)
			case sr.TypeJSON:
				serde.Register(
					subjSchema.ID,
					map[string]any{},
					sr.EncodeFn(json.Marshal),
					sr.DecodeFn(json.Unmarshal),
				)
			case sr.TypeProtobuf:
			default:
				// TODO: support other schema types
				schemaReady = false
			}
		}

		// Single message mode (simplified version)
		value, err := adapter.RequireString(request, "value")
		if err != nil {
			return b.handleError("get value", err), nil
		}

		key := adapter.GetString(request, "key", "")
		headers := adapter.GetStringSlice(request, "headers", []string{})
		sync := adapter.GetBool(request, "sync", true)

		// Prepare record
		record := &kgo.Record{
			Topic: topicName,
			Value: []byte(value),
		}

		// Add key if provided
		if key != "" {
			record.Key = []byte(key)
		}

		// Add headers if provided
		if len(headers) > 0 {
			for _, headerStr := range headers {
				parts := strings.SplitN(headerStr, "=", 2)
				if len(parts) == 2 {
					record.Headers = append(record.Headers, kgo.RecordHeader{
						Key:   parts[0],
						Value: []byte(parts[1]),
					})
				}
			}
		}

		// Handle schema encoding if available
		if schemaReady {
			var jsonValue interface{}
			if err := json.Unmarshal([]byte(value), &jsonValue); err == nil {
				encodedValue, err := serde.Encode(jsonValue)
				if err != nil {
					return b.handleError("encode value with schema", err), nil
				}
				record.Value = encodedValue
			}
		}

		// Produce the message based on sync parameter
		if sync {
			results := kafkaClient.ProduceSync(timeoutCtx, record)
			if len(results) > 0 && results[0].Err != nil {
				return b.handleError("produce message", results[0].Err), nil
			}
		} else {
			kafkaClient.Produce(timeoutCtx, record, func(_ *kgo.Record, _ error) {
				// Log async errors but don't fail since we're async
				// In the future, this could be enhanced with proper async result handling
			})
		}

		// Create result
		response := map[string]interface{}{
			"status":    "success",
			"topic":     record.Topic,
			"timestamp": time.Now().Format(time.RFC3339),
		}

		if len(record.Key) > 0 {
			response["key"] = string(record.Key)
		}

		if record.Partition != -1 {
			response["partition"] = record.Partition
		}

		return b.marshalResponse(response)
	}
}

// Helper functions

// Unified error handling and utility functions

// handleError provides unified error handling
func (b *KafkaProduceToolBuilder) handleError(operation string, err error) *mcpsdk.CallToolResult {
	if err != nil {
		return adapter.NewErrorResult("Failed to %s: %v", operation, err)
	}
	return adapter.NewErrorResult("Failed to %s", operation)
}

// marshalResponse provides unified JSON serialization for responses
func (b *KafkaProduceToolBuilder) marshalResponse(data interface{}) (*mcpsdk.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return b.handleError("marshal response", err), nil
	}
	return adapter.NewTextResult(string(jsonBytes)), nil
}
