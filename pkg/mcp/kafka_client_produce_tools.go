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

package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/hamba/avro/v2"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
)

// KafkaClientAddProduceTools adds Kafka client produce tools to the MCP server
func KafkaClientAddProduceTools(s *server.MCPServer, readOnly bool, features []string) {
	// Skip registration if in read-only mode
	if readOnly {
		return
	}

	if !slices.Contains(features, string(FeatureKafkaClient)) && !slices.Contains(features, string(FeatureAll)) && !slices.Contains(features, string(FeatureAllKafka)) {
		return
	}

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
		"   headers: [{\"key\": \"source\", \"value\": \"mcp-tool\"}, {\"key\": \"timestamp\", \"value\": \"2023-06-01\"}]\n\n" +
		"5. Specific partition - Send to a particular partition:\n" +
		"   topic: \"my-topic\"\n" +
		"   value: \"Targeted message\"\n" +
		"   partition: 2\n\n" +
		"This tool requires Kafka producer permissions on the specified topic."

	// Add produce tool
	produceTool := mcp.NewTool("kafka_client_produce",
		mcp.WithDescription(toolDesc),
		mcp.WithString("topic", mcp.Required(),
			mcp.Description("The name of the Kafka topic to produce messages to. "+
				"Must be an existing topic that the user has write permissions for."),
		),
		mcp.WithString("key",
			mcp.Description("The key for the message. "+
				"Optional. Keys are used for partition assignment and maintaining order for related messages. "+
				"Messages with the same key will be sent to the same partition."),
		),
		mcp.WithString("value",
			mcp.Required(),
			mcp.Description("The value/content of the message to send. "+
				"This is the actual payload that will be delivered to consumers. It can be a JSON string, and the system will automatically serialize it to the appropriate format based on the schema registry if it is available."),
		),
		mcp.WithArray("headers",
			mcp.Description("Message headers in the format of [{\"key\": \"value\"}]. "+
				"Optional. Headers allow you to attach metadata to messages without modifying the payload. "+
				"They are passed along with the message to consumers."),
		),
		mcp.WithBoolean("sync",
			mcp.Description("Whether to wait for server acknowledgment before returning. "+
				"Optional. Default is true. When true, ensures the message was successfully written "+
				"to the topic before the tool returns a success response."),
		),
	)
	s.AddTool(produceTool, handleKafkaProduce)
}

// handleKafkaProduce handles producing messages to a Kafka topic
func handleKafkaProduce(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topicName, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic name: %v", err)), nil
	}

	// Handle single message case
	// Get value from parameter or file
	value, err := requiredParam[string](request.Params.Arguments, "value")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get value: %v", err)), nil
	}

	// Get optional parameters
	key, hasKey := optionalParam[string](request.Params.Arguments, "key")
	headers, hasHeaders := optionalParam[[]interface{}](request.Params.Arguments, "headers")
	sync := true
	if syncVal, hasSync := optionalParam[bool](request.Params.Arguments, "sync"); hasSync {
		sync = syncVal
	}

	// Create Kafka client
	kafkaClient, err := kafka.GetKafkaClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create Kafka client: %v", err)), nil
	}
	defer kafkaClient.Close()

	srClient, err := kafka.GetKafkaSchemaRegistryClient()
	schemaReady := false
	var serde sr.Serde
	if err == nil && srClient != nil {
		schemaReady = true
	}

	if schemaReady {
		// Set timeout
		timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		subjSchema, err := srClient.SchemaByVersion(timeoutCtx, topicName+"-value", -1)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get schema: %v", err)), nil
		}
		logger.Infof("Schema ID: %d", subjSchema.ID)
		switch subjSchema.Type {
		case sr.TypeAvro:
			avroSchema, err := avro.Parse(subjSchema.Schema.Schema)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to parse avro schema: %v", err)), nil
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
			logger.Infof("Unsupported schema type: %s", subjSchema.Type)
			schemaReady = false
		}
	}

	// Prepare record
	record := &kgo.Record{
		Topic: topicName,
		Value: []byte(value),
	}

	// Add key if provided
	if hasKey {
		record.Key = []byte(key)
	}

	// Add headers if provided
	if hasHeaders {
		record.Headers = parseHeaders(headers)
	}

	if schemaReady {
		var structValue map[string]any
		if err := json.Unmarshal([]byte(value), &structValue); err == nil {
			record.Value, err = serde.Encode(structValue)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to encode value: %v", err)), nil
			}
		} else {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to unmarshal value: %v", err)), nil
		}
	}

	// Produce the message
	result, err := produceRecord(ctx, kafkaClient, record, sync)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to produce message: %v", err)), nil
	}

	return result, nil
}

// parseHeaders converts the headers array to kgo.RecordHeader format
func parseHeaders(headers []interface{}) []kgo.RecordHeader {
	var recordHeaders []kgo.RecordHeader

	for _, header := range headers {
		headerMap, ok := header.(map[string]interface{})
		if !ok {
			continue
		}

		for key, value := range headerMap {
			recordHeaders = append(recordHeaders, kgo.RecordHeader{
				Key:   key,
				Value: []byte(value.(string)),
			})
		}
	}

	return recordHeaders
}

// produceRecord sends a single record to Kafka
func produceRecord(ctx context.Context, client *kgo.Client, record *kgo.Record, sync bool) (*mcp.CallToolResult, error) {
	if sync {
		results := client.ProduceSync(ctx, record)
		if len(results) > 0 && results[0].Err != nil {
			return nil, results[0].Err
		}
	} else {
		client.Produce(ctx, record, func(_ *kgo.Record, err error) {
			if err != nil {
				// Log the error but continue since we're async
				fmt.Printf("Async produce error: %v\n", err)
			}
		})
	}

	// Create result
	result := map[string]interface{}{
		"status":    "success",
		"topic":     record.Topic,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	if len(record.Key) > 0 {
		result["key"] = string(record.Key)
	}

	if record.Partition != -1 {
		result["partition"] = record.Partition
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to create result: %v", err)
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}
