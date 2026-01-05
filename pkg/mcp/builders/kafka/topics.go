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
	"fmt"
	"strings"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/adapter"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	"github.com/twmb/franz-go/pkg/kadm"
)

// KafkaTopicsToolBuilder implements the ToolBuilder interface for Kafka Topics
// /nolint:revive
type KafkaTopicsToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewKafkaTopicsToolBuilder creates a new Kafka Topics tool builder instance
func NewKafkaTopicsToolBuilder() *KafkaTopicsToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "kafka_topics",
		Version:     "1.0.0",
		Description: "Kafka Topics administration tools",
		Category:    "kafka_admin",
		Tags:        []string{"kafka", "topics", "admin"},
	}

	features := []string{
		"kafka-admin",
		"all",
		"all-kafka",
	}

	return &KafkaTopicsToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Kafka Topics tool list
func (b *KafkaTopicsToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ServerTool, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool := b.buildKafkaTopicsTool()
	handler := b.buildKafkaTopicsHandler(config.ReadOnly)

	return []builders.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildKafkaTopicsTool builds the Kafka Topics MCP tool definition
func (b *KafkaTopicsToolBuilder) buildKafkaTopicsTool() *mcpsdk.Tool {
	resourceDesc := "Resource to operate on. Available resources:\n" +
		"- topic: A single Kafka topic for operations on individual topics (create, get, delete)\n" +
		"- topics: Collection of Kafka topics for bulk operations (list)"

	operationDesc := "Operation to perform. Available operations:\n" +
		"- list: List all topics in the Kafka cluster, optionally including internal topics\n" +
		"- get: Get detailed configuration for a specific topic\n" +
		"- create: Create a new topic with specified partitions, replication factor, and optional configs\n" +
		"- delete: Delete an existing topic\n" +
		"- metadata: Get metadata for a specific topic\n"

	toolDesc := "Unified tool for managing Apache Kafka topics.\n" +
		"This tool provides access to various Kafka topic operations, including creation, deletion, listing, and configuration retrieval.\n" +
		"Kafka topics are the core abstraction for organizing and partitioning data streams. Topics:\n" +
		"- Organize messages into categories for producers and consumers\n" +
		"- Are divided into partitions for scalability and parallelism\n" +
		"- Can be configured with replication factors for fault tolerance\n" +
		"- Support various configuration options for retention, compression, and more\n\n" +
		"Usage Examples:\n\n" +
		"1. List all topics (excluding internal Kafka topics):\n" +
		"   resource: \"topics\"\n" +
		"   operation: \"list\"\n\n" +
		"2. List all topics including internal ones:\n" +
		"   resource: \"topics\"\n" +
		"   operation: \"list\"\n" +
		"   includeInternal: true\n\n" +
		"3. Create a new topic with default settings:\n" +
		"   resource: \"topic\"\n" +
		"   operation: \"create\"\n" +
		"   name: \"user-events\"\n" +
		"   partitions: 3\n" +
		"   replicationFactor: 2\n\n" +
		"4. Create a topic with custom configuration:\n" +
		"   resource: \"topic\"\n" +
		"   operation: \"create\"\n" +
		"   name: \"log-aggregation\"\n" +
		"   partitions: 6\n" +
		"   replicationFactor: 3\n" +
		"   configs: {\n" +
		"     \"retention.ms\": \"604800000\",\n" +
		"     \"compression.type\": \"gzip\",\n" +
		"     \"cleanup.policy\": \"compact\"\n" +
		"   }\n\n" +
		"5. Get detailed information about a topic:\n" +
		"   resource: \"topic\"\n" +
		"   operation: \"get\"\n" +
		"   name: \"user-events\"\n\n" +
		"6. Get metadata for a topic:\n" +
		"   resource: \"topic\"\n" +
		"   operation: \"metadata\"\n" +
		"   name: \"user-events\"\n\n" +
		"7. Delete a topic:\n" +
		"   resource: \"topic\"\n" +
		"   operation: \"delete\"\n" +
		"   name: \"old-topic\"\n\n" +
		"This tool requires appropriate Kafka permissions for topic management."

	return builders.NewTool("kafka_admin_topics",
		builders.WithDescription(toolDesc),
		builders.WithString("resource", builders.Required(),
			builders.Description(resourceDesc),
		),
		builders.WithString("operation", builders.Required(),
			builders.Description(operationDesc),
		),
		builders.WithString("name",
			builders.Description("The name of the Kafka topic to operate on. "+
				"Required for 'get', 'create', 'delete', and 'metadata' operations on the 'topic' resource. "+
				"Topic names should follow Kafka naming conventions (alphanumeric, dots, underscores, and hyphens).")),
		builders.WithNumber("partitions",
			builders.Description("The number of partitions for the topic. Required for 'create' operation. "+
				"Partitions determine the parallelism and scalability of the topic. "+
				"More partitions allow more concurrent consumers and higher throughput.")),
		builders.WithNumber("replicationFactor",
			builders.Description("The replication factor for the topic. Required for 'create' operation. "+
				"Replication factor determines fault tolerance - it should be at least 2 for production use. "+
				"Cannot exceed the number of available brokers in the cluster.")),
		builders.WithObject("configs",
			builders.Description("Optional configuration parameters for the topic during 'create' operation. "+
				"Common configurations include:\n"+
				"- retention.ms: How long to retain messages (milliseconds)\n"+
				"- compression.type: Compression algorithm (none, gzip, snappy, lz4, zstd)\n"+
				"- cleanup.policy: Log cleanup policy (delete, compact, compact,delete)\n"+
				"- segment.ms: Time before a new log segment is rolled out\n"+
				"- max.message.bytes: Maximum size of a message batch")),
		builders.WithBoolean("includeInternal",
			builders.Description("Whether to include internal Kafka topics in the 'list' operation. "+
				"Internal topics are used by Kafka itself (e.g., __consumer_offsets, __transaction_state). "+
				"Default: false")),
	)
}

// buildKafkaTopicsHandler builds the Kafka Topics handler function
func (b *KafkaTopicsToolBuilder) buildKafkaTopicsHandler(readOnly bool) mcpsdk.ToolHandler {
	return func(ctx context.Context, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
		// Get required parameters
		resource, err := adapter.RequireString(request, "resource")
		if err != nil {
			return b.handleError("get resource", err), nil
		}

		operation, err := adapter.RequireString(request, "operation")
		if err != nil {
			return b.handleError("get operation", err), nil
		}

		// Normalize parameters
		resource = strings.ToLower(resource)
		operation = strings.ToLower(operation)

		// Validate write operations in read-only mode
		if readOnly && (operation == "create" || operation == "delete") {
			return adapter.NewErrorResult("Write operations are not allowed in read-only mode"), nil
		}

		// Get Kafka admin client
		session := mcpCtx.GetKafkaSession(ctx)
		if session == nil {
			return b.handleError("get Kafka session not found in context", nil), nil
		}
		admin, err := session.GetAdminClient()
		if err != nil {
			return b.handleError("get admin client", err), nil
		}

		// Dispatch based on resource and operation
		switch resource {
		case "topics":
			switch operation {
			case "list":
				return b.handleKafkaTopicsList(ctx, admin, request)
			default:
				return adapter.NewErrorResult("Invalid operation for resource 'topics': %s", operation), nil
			}
		case "topic":
			switch operation {
			case "get":
				return b.handleKafkaTopicGet(ctx, admin, request)
			case "create":
				return b.handleKafkaTopicCreate(ctx, admin, request)
			case "delete":
				return b.handleKafkaTopicDelete(ctx, admin, request)
			case "metadata":
				return b.handleKafkaTopicMetadata(ctx, admin, request)
			default:
				return adapter.NewErrorResult("Invalid operation for resource 'topic': %s", operation), nil
			}
		default:
			return adapter.NewErrorResult("Invalid resource: %s. Available resources: topics, topic", resource), nil
		}
	}
}

// Utility functions

// handleError provides unified error handling
func (b *KafkaTopicsToolBuilder) handleError(operation string, err error) *mcpsdk.CallToolResult {
	if err != nil {
		return adapter.NewErrorResult("Failed to %s: %v", operation, err)
	}
	return adapter.NewErrorResult("Failed to %s", operation)
}

// marshalResponse provides unified JSON serialization for responses
func (b *KafkaTopicsToolBuilder) marshalResponse(data interface{}) (*mcpsdk.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return b.handleError("marshal response", err), nil
	}
	return adapter.NewTextResult(string(jsonBytes)), nil
}

// handleKafkaTopicsList handles listing all topics
func (b *KafkaTopicsToolBuilder) handleKafkaTopicsList(ctx context.Context, admin *kadm.Client, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	includeInternal := adapter.GetBool(request, "includeInternal", false)

	topics, err := admin.ListTopics(ctx)
	if err != nil {
		return b.handleError("list Kafka topics", err), nil
	}

	// Filter out internal topics if not requested
	if !includeInternal {
		filteredTopics := make(kadm.TopicDetails)
		for name, details := range topics {
			if !strings.HasPrefix(name, "__") {
				filteredTopics[name] = details
			}
		}
		topics = filteredTopics
	}

	return b.marshalResponse(topics)
}

// handleKafkaTopicGet handles getting detailed information about a specific topic
func (b *KafkaTopicsToolBuilder) handleKafkaTopicGet(ctx context.Context, admin *kadm.Client, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	topicName, err := adapter.RequireString(request, "name")
	if err != nil {
		return b.handleError("get topic name", err), nil
	}

	topics, err := admin.ListTopics(ctx, topicName)
	if err != nil {
		return b.handleError("get Kafka topic", err), nil
	}

	return b.marshalResponse(topics)
}

// handleKafkaTopicCreate handles creating a new topic
func (b *KafkaTopicsToolBuilder) handleKafkaTopicCreate(ctx context.Context, admin *kadm.Client, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	topicName, err := adapter.RequireString(request, "name")
	if err != nil {
		return b.handleError("get topic name", err), nil
	}

	partitionsNum, err := adapter.RequireInt(request, "partitions")
	if err != nil {
		return b.handleError("get partitions", err), nil
	}

	replicationFactorNum, err := adapter.RequireInt(request, "replicationFactor")
	if err != nil {
		return b.handleError("get replication factor", err), nil
	}

	///nolint:gosec
	partitions := int32(partitionsNum)
	///nolint:gosec
	replicationFactor := int16(replicationFactorNum)

	// Parse optional configs
	var configs map[string]*string
	arguments, err := adapter.GetArgumentsMap(request)
	if err == nil {
		if configsParam, exists := arguments["configs"]; exists {
			if configsMap, ok := configsParam.(map[string]interface{}); ok {
				configs = make(map[string]*string)
				for key, value := range configsMap {
					if strValue, ok := value.(string); ok {
						configs[key] = &strValue
					} else {
						// Convert non-string values to strings
						strValue := fmt.Sprintf("%v", value)
						configs[key] = &strValue
					}
				}
			}
		}
	}

	// Create topic using the correct CreateTopics API
	results, err := admin.CreateTopics(ctx, partitions, replicationFactor, configs, topicName)
	if err != nil {
		return b.handleError("create Kafka topic", err), nil
	}

	return b.marshalResponse(results)
}

// handleKafkaTopicDelete handles deleting a topic
func (b *KafkaTopicsToolBuilder) handleKafkaTopicDelete(ctx context.Context, admin *kadm.Client, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	topicName, err := adapter.RequireString(request, "name")
	if err != nil {
		return b.handleError("get topic name", err), nil
	}

	results, err := admin.DeleteTopics(ctx, topicName)
	if err != nil {
		return b.handleError("delete Kafka topic", err), nil
	}

	return b.marshalResponse(results)
}

// handleKafkaTopicMetadata handles getting metadata for a topic
func (b *KafkaTopicsToolBuilder) handleKafkaTopicMetadata(ctx context.Context, admin *kadm.Client, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	topicName, err := adapter.RequireString(request, "name")
	if err != nil {
		return b.handleError("get topic name", err), nil
	}

	metadata, err := admin.Metadata(ctx, topicName)
	if err != nil {
		return b.handleError("get Kafka topic metadata", err), nil
	}

	return b.marshalResponse(metadata)
}
