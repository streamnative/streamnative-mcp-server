package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	"github.com/twmb/franz-go/pkg/kgo"
)

// KafkaClientAddConsumeTools adds Kafka client consume tools to the MCP server
func KafkaClientAddConsumeTools(s *server.MCPServer, readOnly bool) {
	toolDesc := "Consume messages from a Kafka topic.\n" +
		"This tool allows you to read messages from Kafka topics, specifying various consumption parameters.\n\n" +
		"Kafka Consumer Concepts:\n" +
		"- Consumers read data from Kafka topics, which can be spread across multiple partitions\n" +
		"- Consumer Groups enable multiple consumers to cooperatively process messages from the same topic\n" +
		"- Offsets track the position of consumers in each partition, allowing resumption after failures\n" +
		"- Partitions are independent ordered sequences of messages that enable parallel processing\n\n" +
		"This tool provides a temporary consumer instance for diagnostic and testing purposes. " +
		"It does not commit offsets back to Kafka unless the 'group' parameter is explicitly specified.\n\n" +
		"Usage Examples:\n\n" +
		"1. Basic consumption - Get 10 latest messages from a topic:\n" +
		"   topic: \"my-topic\"\n" +
		"   max-messages: 10\n\n" +
		"2. Consume from beginning - Get messages from the start of a topic:\n" +
		"   topic: \"my-topic\"\n" +
		"   offset: \"atstart\"\n" +
		"   max-messages: 20\n\n" +
		"3. Consumer group - Join an existing consumer group and resume from committed offset:\n" +
		"   topic: \"my-topic\"\n" +
		"   group: \"my-consumer-group\"\n" +
		"   offset: \"atcommitted\"\n\n" +
		"4. Specific partition - Read from a particular partition starting at the beginning:\n" +
		"   topic: \"my-topic\"\n" +
		"   partition: 0\n" +
		"   offset: \"atstart\"\n" +
		"   max-messages: 50\n\n" +
		"5. Time-limited consumption - Set a longer timeout for slow topics:\n" +
		"   topic: \"my-topic\"\n" +
		"   max-messages: 100\n" +
		"   timeout: 30\n\n" +
		"This tool requires Kafka consumer permissions on the specified topic."

	// Add consume tool
	consumeTool := mcp.NewTool("kafka_client_consume",
		mcp.WithDescription(toolDesc),
		mcp.WithString("topic", mcp.Required(),
			mcp.Description("The name of the Kafka topic to consume messages from. "+
				"Must be an existing topic that the user has read permissions for. "+
				"For partitioned topics, this will consume from all partitions unless a specific partition is specified."),
		),
		mcp.WithString("group",
			mcp.Description("The consumer group ID to use for consumption. "+
				"Optional. If provided, the consumer will join this consumer group and track offsets with Kafka. "+
				"If not provided, a random group ID will be generated, and offsets won't be committed back to Kafka. "+
				"Using a meaningful group ID is important when you want to resume consumption later or coordinate multiple consumers."),
		),
		mcp.WithNumber("partition",
			mcp.Description("The specific partition number to consume from. "+
				"Optional. Partitions are zero-indexed (0, 1, 2, etc). "+
				"If not specified, messages will be consumed from all partitions of the topic. "+
				"This is useful when you need to debug or read from a specific partition. "+
				"When specified, the 'offset' parameter determines where in the partition to start reading."),
		),
		mcp.WithString("offset",
			mcp.Description("The offset position to start consuming from. "+
				"Optional. Must be one of these values:\n"+
				"- 'atstart': Begin from the earliest available message in the topic/partition\n"+
				"- 'atend': Begin from the next message that arrives after the consumer starts\n"+
				"- 'atcommitted': Begin from the last committed offset (only works with specified 'group')\n"+
				"Default: 'atstart'"),
		),
		mcp.WithNumber("max-messages",
			mcp.Description("Maximum number of messages to consume in this request. "+
				"Optional. Limits the total number of messages returned, across all partitions if no specific partition is specified. "+
				"Higher values retrieve more data but may increase response time and size. "+
				"Default: 10"),
		),
		mcp.WithNumber("timeout",
			mcp.Description("Maximum time in seconds to wait for messages. "+
				"Optional. The consumer will wait up to this long to collect the requested number of messages. "+
				"If fewer than 'max-messages' are available within this time, the available messages are returned. "+
				"Longer timeouts are useful for low-volume topics or when consuming with 'atend'. "+
				"Default: 10 seconds"),
		),
	)
	s.AddTool(consumeTool, handleKafkaConsume)
}

// handleKafkaConsume handles consuming messages from a Kafka topic
func handleKafkaConsume(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	opts := []kgo.Opt{}
	// Get required parameters
	topicName, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic name: %v", err)), nil
	}
	opts = append(opts, kgo.ConsumeTopics(topicName))

	// Get optional parameters
	groupID, hasGroup := optionalParam[string](request.Params.Arguments, "group")
	if !hasGroup {
		// Generate a random group ID if not provided
		groupID = fmt.Sprintf("snmcp-kafka-%d", time.Now().UnixNano())
	}
	opts = append(opts, kgo.ConsumerGroup(groupID))
	partitionID, hasPartition := optionalParam[float64](request.Params.Arguments, "partition")

	offsetStr, hasOffset := optionalParam[string](request.Params.Arguments, "offset")
	if !hasOffset {
		offsetStr = "atstart" // Default to starting at the beginning
	}

	var offset kgo.Offset
	switch offsetStr {
	case "atstart":
		offset = kgo.NewOffset().AtStart()
	case "atend":
		offset = kgo.NewOffset().AtEnd()
	case "atcommitted":
		offset = kgo.NewOffset().AtCommitted()
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Invalid offset value: %s. Must be one of 'atstart', 'atend', or 'atcommitted'", offsetStr)), nil
	}

	if hasPartition {
		// Consume from specific partition with specified offset
		opts = append(opts, kgo.ConsumePartitions(map[string]map[int32]kgo.Offset{
			topicName: {
				int32(partitionID): offset,
			},
		}))
	} else {
		opts = append(opts, kgo.ConsumeResetOffset(offset))
	}

	// maxMessages, hasMaxMessages := optionalParam[float64](request.Params.Arguments, "max-messages")
	// if !hasMaxMessages {
	// 	maxMessages = 10 // Default to 10 messages
	// }

	timeoutSec, hasTimeout := optionalParam[float64](request.Params.Arguments, "timeout")
	if !hasTimeout {
		timeoutSec = 10 // Default to 10 seconds
	}

	opts = append(opts, kgo.BlockRebalanceOnPoll())

	// Create Kafka client using the new Kafka package
	kafkaClient, err := kafka.GetKafkaClient(opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create Kafka client: %v", err)), nil
	}
	defer kafkaClient.Close()

	// Set timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSec)*time.Second)
	defer cancel()

	if err = kafkaClient.Ping(context.Background()); err != nil { // check connectivity to cluster
		return mcp.NewToolResultError(fmt.Sprintf("Failed to ping Kafka cluster: %v", err)), nil
	}

	// Consume messages
	fetches := kafkaClient.PollFetches(timeoutCtx)
	if fetches.IsClientClosed() {
		return mcp.NewToolResultError(fmt.Sprintf("Kafka client closed: %v", err)), nil
	}

	records := fetches.Records()

	// Format the output as JSON
	jsonBytes, err := json.Marshal(records)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal messages: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}
