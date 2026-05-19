// Copyright 2026 StreamNative
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
	"fmt"
	"time"

	"github.com/hamba/avro/v2"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/toolannotations"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
)

// KafkaConsumeToolBuilder implements the ToolBuilder interface for Kafka client consume operations
// It provides functionality to build Kafka consumer tools
// /nolint:revive
type KafkaConsumeToolBuilder struct {
	*builders.BaseToolBuilder
	logger *logrus.Logger
}

// NewKafkaConsumeToolBuilder creates a new Kafka consume tool builder instance
func NewKafkaConsumeToolBuilder() *KafkaConsumeToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "kafka_consume",
		Version:     "1.0.0",
		Description: "Kafka client consume tools",
		Category:    "kafka_client",
		Tags:        []string{"kafka", "client", "consume"},
	}

	features := []string{
		"kafka-client",
		"all",
		"all-kafka",
	}

	return &KafkaConsumeToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Kafka consume tool list
// This is the core method implementing the ToolBuilder interface
func (b *KafkaConsumeToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
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

	// Extract logger from options if provided
	if loggerOpt, ok := config.Options["logger"]; ok {
		if logger, ok := loggerOpt.(*logrus.Logger); ok {
			b.logger = logger
		}
	}

	// Build tools
	tool := b.buildKafkaConsumeTool()
	handler := b.buildKafkaConsumeHandler()

	return []server.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildKafkaConsumeTool builds the Kafka consume MCP tool definition
// Migrated from the original tool definition logic
func (b *KafkaConsumeToolBuilder) buildKafkaConsumeTool() mcp.Tool {
	toolDesc := "Consume messages from a Kafka topic.\n" +
		"This tool allows you to read messages from Kafka topics, specifying various consumption parameters.\n\n" +
		"Kafka Consumer Concepts:\n" +
		"- Consumers read data from Kafka topics, which can be spread across multiple partitions\n" +
		"- Consumer Groups enable multiple consumers to cooperatively process messages from the same topic\n" +
		"- Offsets track the position of consumers in each partition, allowing resumption after failures\n" +
		"- Partitions are independent ordered sequences of messages that enable parallel processing\n\n" +
		"This tool provides a temporary consumer instance for diagnostic and testing purposes. " +
		"It does not commit offsets back to Kafka unless the 'group' parameter is explicitly specified. Do not use this tool for Pulsar protocol operations. Use 'pulsar_client_consume' instead.\n\n" +
		"Usage Examples:\n\n" +
		"1. Basic consumption - Get 10 earliest messages from a topic:\n" +
		"   topic: \"my-topic\"\n" +
		"   max-messages: 10\n\n" +
		"2. Consumer group - Join an existing consumer group and resume from committed offset:\n" +
		"   topic: \"my-topic\"\n" +
		"   offset: \"atstart\"\n" +
		"   max-messages: 20\n\n" +
		"3. Consumer group - Join an existing consumer group and resume from committed offset:\n" +
		"   topic: \"my-topic\"\n" +
		"   group: \"my-consumer-group\"\n" +
		"   offset: \"atcommitted\"\n\n" +
		"4. Time-limited consumption - Set a longer timeout for slow topics:\n" +
		"   topic: \"my-topic\"\n" +
		"   max-messages: 100\n" +
		"   timeout: 30\n\n" +
		"This tool requires Kafka consumer permissions on the specified topic."

	return mcp.NewTool("kafka_client_consume",
		mcp.WithDescription(toolDesc),
		toolannotations.Destructive("Consume Kafka Messages"),
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
}

// buildKafkaConsumeHandler builds the Kafka consume handler function
// Migrated from the original handler logic
func (b *KafkaConsumeToolBuilder) buildKafkaConsumeHandler() func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		opts := []kgo.Opt{}
		// Get required parameters
		topicName, err := request.RequireString("topic")
		if err != nil {
			return b.handleError("get topic name", err), nil
		}
		opts = append(opts, kgo.ConsumeTopics(topicName))

		opts = append(opts, kgo.FetchIsolationLevel(kgo.ReadUncommitted()))
		opts = append(opts, kgo.KeepRetryableFetchErrors())
		if b.logger != nil {
			w := b.logger.Writer()
			opts = append(opts, kgo.WithLogger(kgo.BasicLogger(w, kgo.LogLevelInfo, nil)))
		}
		maxMessages := request.GetFloat("max-messages", 10)

		timeoutSec := request.GetFloat("timeout", 10)

		group := request.GetString("group", "")
		if group != "" {
			opts = append(opts, kgo.ConsumerGroup(group))
		}

		offsetStr := request.GetString("offset", "atstart")

		var offset kgo.Offset
		switch offsetStr {
		case "atstart":
			offset = kgo.NewOffset().AtStart()
		case "atend":
			offset = kgo.NewOffset().AtEnd()
		case "atcommitted":
			offset = kgo.NewOffset().AtCommitted()
		default:
			offset = kgo.NewOffset().AtStart()
		}
		opts = append(opts, kgo.ConsumeResetOffset(offset))
		if b.logger != nil {
			b.logger.Infof("Consuming from topic: %s, group: %s, max-messages: %d, timeout: %d", topicName, group, int(maxMessages), int(timeoutSec))
		}

		// Get Kafka session from context
		session := mcpCtx.GetKafkaSession(ctx)
		if session == nil {
			return b.handleError("get Kafka session not found in context", nil), nil
		}

		// Create Kafka client using the session
		kafkaClient, err := session.GetClient(opts...)
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
		timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSec)*time.Second)
		defer cancel()

		if err = kafkaClient.Ping(timeoutCtx); err != nil { // check connectivity to cluster
			return b.handleError("ping Kafka cluster", err), nil
		}

		if schemaReady {
			subjSchema, err := srClient.SchemaByVersion(timeoutCtx, topicName+"-value", -1)
			if err != nil {
				return b.handleError("get schema", err), nil
			}
			if b.logger != nil {
				b.logger.Infof("Schema ID: %d", subjSchema.ID)
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
				if b.logger != nil {
					b.logger.Infof("Unsupported schema type: %s", subjSchema.Type)
				}
				schemaReady = false
			}
		}

		results := make([]any, 0)
	consumerLoop:
		for {
			fetches := kafkaClient.PollRecords(timeoutCtx, int(maxMessages))
			iter := fetches.RecordIter()
			if b.logger != nil {
				b.logger.Infof("NumRecords: %d\n", fetches.NumRecords())
			}

			for _, fetchErr := range fetches.Errors() {
				if b.logger != nil {
					b.logger.Infof("error consuming from topic: topic=%s, partition=%d, err=%v\n",
						fetchErr.Topic, fetchErr.Partition, fetchErr.Err)
				}
				break consumerLoop
			}

			for !iter.Done() {
				record := iter.Next()
				if schemaReady {
					var value map[string]any
					err := serde.Decode(record.Value, &value)
					if err != nil {
						if b.logger != nil {
							b.logger.Infof("Failed to decode value: %v", err)
						}
						results = append(results, record.Value)
					} else {
						results = append(results, value)
					}
				} else {
					results = append(results, record.Value)
				}
				if len(results) >= int(maxMessages) {
					break consumerLoop
				}
			}
		}

		err = kafkaClient.CommitUncommittedOffsets(timeoutCtx)
		if err != nil {
			if err != context.Canceled && b.logger != nil {
				b.logger.Infof("Failed to commit offsets: %v", err)
			}
		}

		return b.marshalResponse(results)
	}
}

// Unified error handling and utility functions

// handleError provides unified error handling
func (b *KafkaConsumeToolBuilder) handleError(operation string, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("Failed to %s: %v", operation, err))
}

// marshalResponse provides unified JSON serialization for responses
func (b *KafkaConsumeToolBuilder) marshalResponse(data interface{}) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return b.handleError("marshal response", err), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}
