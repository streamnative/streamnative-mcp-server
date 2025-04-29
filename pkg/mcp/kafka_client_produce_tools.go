package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	"github.com/twmb/franz-go/pkg/kgo"
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
		"This tool provides a simple producer instance for diagnostic and testing purposes.\n\n" +
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
				"Messages with the same key will be sent to the same partition. "+
				"If not provided, the message will be assigned to a random partition."),
		),
		mcp.WithString("value",
			mcp.Description("The value/content of the message to send. "+
				"Required if 'messages' is not provided. This is the actual payload that will be delivered to consumers."),
		),
		mcp.WithNumber("partition",
			mcp.Description("The specific partition number to produce to. "+
				"Optional. Partitions are zero-indexed (0, 1, 2, etc). "+
				"If specified, the message will be sent to this partition regardless of the key. "+
				"If not specified, the partition will be determined by the key hash or randomly assigned."),
		),
		mcp.WithArray("headers",
			mcp.Description("Message headers in the format of [{\"key\": \"header-key\", \"value\": \"header-value\"}]. "+
				"Optional. Headers allow you to attach metadata to messages without modifying the payload. "+
				"They are passed along with the message to consumers but are not used for routing."),
		),
		mcp.WithArray("messages",
			mcp.Description("Batch of messages to send in the format of [{\"key\": \"key1\", \"value\": \"value1\", \"headers\": [...], \"partition\": 0}, ...]. "+
				"Optional. Allows sending multiple messages with different keys, values, headers, and partitions in a single request. "+
				"If provided, the individual key/value/headers/partition parameters are ignored."),
		),
		mcp.WithBoolean("sync",
			mcp.Description("Whether to wait for server acknowledgment before returning. "+
				"Optional. Default is true. When true, ensures the message was successfully written "+
				"to the topic before the tool returns a success response."),
		),
		mcp.WithString("file",
			mcp.Description("Path to a file containing the message value. "+
				"Optional. If provided, the file content will be used as the message value instead of the 'value' parameter. "+
				"Useful for sending larger messages or binary data."),
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

	// Process batch messages first if provided
	messages, hasMessages := optionalParam[[]interface{}](request.Params.Arguments, "messages")
	if hasMessages && len(messages) > 0 {
		return handleBatchMessages(ctx, topicName, messages)
	}

	// Handle single message case
	// Get value from parameter or file
	value, hasValue := optionalParam[string](request.Params.Arguments, "value")
	filePath, hasFile := optionalParam[string](request.Params.Arguments, "file")

	if !hasValue && !hasFile {
		return mcp.NewToolResultError("Either 'value' or 'file' parameter is required when not using 'messages'"), nil
	}

	// If file is provided, read its content
	if hasFile {
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to read file: %v", err)), nil
		}
		value = string(fileContent)
	}

	// Get optional parameters
	key, hasKey := optionalParam[string](request.Params.Arguments, "key")
	partition, hasPartition := optionalParam[float64](request.Params.Arguments, "partition")
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

	// Prepare record
	record := &kgo.Record{
		Topic: topicName,
		Value: []byte(value),
	}

	// Add key if provided
	if hasKey {
		record.Key = []byte(key)
	}

	// Set partition if provided
	if hasPartition {
		record.Partition = int32(partition)
	}

	// Add headers if provided
	if hasHeaders {
		record.Headers = parseHeaders(headers)
	}

	// Produce the message
	result, err := produceRecord(ctx, kafkaClient, record, sync)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to produce message: %v", err)), nil
	}

	return result, nil
}

// handleBatchMessages processes multiple messages from the messages array
func handleBatchMessages(ctx context.Context, topicName string, messages []interface{}) (*mcp.CallToolResult, error) {
	// Create Kafka client
	kafkaClient, err := kafka.GetKafkaClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create Kafka client: %v", err)), nil
	}
	defer kafkaClient.Close()

	// Process each message
	var records []*kgo.Record
	for i, msg := range messages {
		msgMap, ok := msg.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid message format at index %d", i)), nil
		}

		// Value is required
		valueInterface, exists := msgMap["value"]
		if !exists {
			return mcp.NewToolResultError(fmt.Sprintf("Message at index %d is missing required 'value' field", i)), nil
		}
		value, ok := valueInterface.(string)
		if !ok {
			return mcp.NewToolResultError(fmt.Sprintf("Value at index %d must be a string", i)), nil
		}

		// Create record
		record := &kgo.Record{
			Topic: topicName,
			Value: []byte(value),
		}

		// Add key if provided
		if keyInterface, exists := msgMap["key"]; exists {
			if key, ok := keyInterface.(string); ok {
				record.Key = []byte(key)
			}
		}

		// Set partition if provided
		if partitionInterface, exists := msgMap["partition"]; exists {
			if partition, ok := partitionInterface.(float64); ok {
				record.Partition = int32(partition)
			}
		}

		// Add headers if provided
		if headersInterface, exists := msgMap["headers"]; exists {
			if headers, ok := headersInterface.([]interface{}); ok {
				record.Headers = parseHeaders(headers)
			}
		}

		records = append(records, record)
	}

	// No valid records to produce
	if len(records) == 0 {
		return mcp.NewToolResultError("No valid messages to produce"), nil
	}

	// Produce all records
	results := kafkaClient.ProduceSync(ctx, records...)
	// Check for errors in produce results
	for _, result := range results {
		if result.Err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to produce messages: %v", result.Err)), nil
		}
	}

	// Create result
	result := map[string]interface{}{
		"status":        "success",
		"topic":         topicName,
		"message_count": len(records),
		"timestamp":     time.Now().Format(time.RFC3339),
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// parseHeaders converts the headers array to kgo.RecordHeader format
func parseHeaders(headers []interface{}) []kgo.RecordHeader {
	var recordHeaders []kgo.RecordHeader

	for _, header := range headers {
		headerMap, ok := header.(map[string]interface{})
		if !ok {
			continue
		}

		key, hasKey := headerMap["key"].(string)
		value, hasValue := headerMap["value"].(string)

		if hasKey && hasValue {
			recordHeaders = append(recordHeaders, kgo.RecordHeader{
				Key:   key,
				Value: []byte(value),
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
		client.Produce(ctx, record, func(r *kgo.Record, err error) {
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
