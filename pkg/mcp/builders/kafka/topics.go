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
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/toolannotations"
	"github.com/twmb/franz-go/pkg/kadm"
)

var kafkaTopicWriteOperations = map[string]struct{}{
	"create": {},
	"delete": {},
}

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
func (b *KafkaTopicsToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	tools := []server.ServerTool{
		{
			Tool:    b.buildKafkaTopicsTool(toolModeRead),
			Handler: b.buildKafkaTopicsHandler(toolModeRead),
		},
	}
	if !config.ReadOnly {
		tools = append(tools, server.ServerTool{
			Tool:    b.buildKafkaTopicsTool(toolModeWrite),
			Handler: b.buildKafkaTopicsHandler(toolModeWrite),
		})
	}

	return tools, nil
}

// buildKafkaTopicsTool builds the Kafka Topics MCP tool definition
func (b *KafkaTopicsToolBuilder) buildKafkaTopicsTool(mode toolMode) mcp.Tool {
	resourceDesc := "Resource to operate on. Available resources:\n" +
		"- topic: A single Kafka topic for read operations (get, metadata)\n" +
		"- topics: Collection of Kafka topics for list operations"
	resourceEnum := []string{"topic", "topics"}

	operationDesc := "Operation to perform. Available operations:\n" +
		"- list: List all topics in the Kafka cluster, optionally including internal topics\n" +
		"- get: Get detailed configuration for a specific topic\n" +
		"- metadata: Get metadata for a specific topic\n"
	operationEnum := []string{"list", "get", "metadata"}
	toolName := "kafka_admin_topics_read"
	annotation := toolannotations.ReadOnly("Read Kafka Topics")
	if isToolModeWrite(mode) {
		operationDesc = "Operation to perform. Available operations:\n" +
			"- create: Create a new topic with specified partitions, replication factor, and optional configs\n" +
			"- delete: Delete an existing topic\n"
		operationEnum = []string{"create", "delete"}
		resourceDesc = "Resource to operate on. Available resources:\n" +
			"- topic: A single Kafka topic for create or delete operations"
		resourceEnum = []string{"topic"}
		toolName = "kafka_admin_topics_write"
		annotation = toolannotations.Destructive("Manage Kafka Topics")
	}

	toolDesc := "Read Apache Kafka topic metadata and lists.\n" +
		"This read-only tool lists topics and retrieves topic details or metadata without changing the cluster.\n\n" +
		"Usage Examples:\n\n" +
		"1. List all topics (excluding internal Kafka topics):\n" +
		"   resource: \"topics\"\n" +
		"   operation: \"list\"\n\n" +
		"2. Get detailed information about a topic:\n" +
		"   resource: \"topic\"\n" +
		"   operation: \"get\"\n" +
		"   name: \"user-events\"\n\n" +
		"3. Get metadata for a topic:\n" +
		"   resource: \"topic\"\n" +
		"   operation: \"metadata\"\n" +
		"   name: \"user-events\"\n\n" +
		"This tool requires appropriate Kafka permissions for topic reads."
	if isToolModeWrite(mode) {
		toolDesc = "Manage Apache Kafka topic lifecycle.\n" +
			"This write tool creates and deletes Kafka topics and may change cluster state.\n\n" +
			"Usage Examples:\n\n" +
			"1. Create a new topic with default settings:\n" +
			"   resource: \"topic\"\n" +
			"   operation: \"create\"\n" +
			"   name: \"user-events\"\n" +
			"   partitions: 3\n" +
			"   replicationFactor: 2\n\n" +
			"2. Delete a topic:\n" +
			"   resource: \"topic\"\n" +
			"   operation: \"delete\"\n" +
			"   name: \"old-topic\"\n\n" +
			"This tool requires appropriate Kafka permissions for topic management."
	}

	tool := mcp.NewTool(toolName,
		mcp.WithDescription(toolDesc),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description(resourceDesc),
			mcp.Enum(resourceEnum...),
		),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
			mcp.Enum(operationEnum...),
		),
		mcp.WithString("name",
			mcp.Description("The name of the Kafka topic to operate on. Required for operations that target one topic. "+
				"Topic names should follow Kafka naming conventions (alphanumeric, dots, underscores, and hyphens).")),
		mcp.WithNumber("partitions",
			mcp.Description("The number of partitions for the topic. Required for 'create' operation. "+
				"Partitions determine the parallelism and scalability of the topic. "+
				"More partitions allow more concurrent consumers and higher throughput.")),
		mcp.WithNumber("replicationFactor",
			mcp.Description("The replication factor for the topic. Required for 'create' operation. "+
				"Replication factor determines fault tolerance - it should be at least 2 for production use. "+
				"Cannot exceed the number of available brokers in the cluster.")),
		mcp.WithObject("configs",
			mcp.Description("Optional configuration parameters for the topic during 'create' operation. "+
				"Common configurations include:\n"+
				"- retention.ms: How long to retain messages (milliseconds)\n"+
				"- compression.type: Compression algorithm (none, gzip, snappy, lz4, zstd)\n"+
				"- cleanup.policy: Log cleanup policy (delete, compact, compact,delete)\n"+
				"- segment.ms: Time before a new log segment is rolled out\n"+
				"- max.message.bytes: Maximum size of a message batch")),
		mcp.WithBoolean("includeInternal",
			mcp.Description("Whether to include internal Kafka topics in the 'list' operation. "+
				"Internal topics are used by Kafka itself (e.g., __consumer_offsets, __transaction_state). "+
				"Default: false")),
		annotation,
	)
	if isToolModeWrite(mode) {
		pruneToolInputSchema(&tool, []string{"resource", "operation", "name", "partitions", "replicationFactor", "configs"})
	} else {
		pruneToolInputSchema(&tool, []string{"resource", "operation", "name", "includeInternal"})
	}
	return tool
}

// buildKafkaTopicsHandler builds the Kafka Topics handler function
func (b *KafkaTopicsToolBuilder) buildKafkaTopicsHandler(mode toolMode) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get required parameters
		resource, err := request.RequireString("resource")
		if err != nil {
			return b.handleError("get resource", err), nil
		}

		operation, err := request.RequireString("operation")
		if err != nil {
			return b.handleError("get operation", err), nil
		}

		// Normalize parameters
		resource = strings.ToLower(resource)
		operation = strings.ToLower(operation)

		if !validateModeOperation(mode, operation, kafkaTopicWriteOperations) {
			return mcp.NewToolResultError(fmt.Sprintf("Operation %q is not available in %s mode", operation, mode)), nil
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
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'topics': %s", operation)), nil
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
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'topic': %s", operation)), nil
			}
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resource: %s. Available resources: topics, topic", resource)), nil
		}
	}
}

// Utility functions

// handleError provides unified error handling
func (b *KafkaTopicsToolBuilder) handleError(operation string, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("Failed to %s: %v", operation, err))
}

// marshalResponse provides unified JSON serialization for responses
func (b *KafkaTopicsToolBuilder) marshalResponse(data interface{}) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return b.handleError("marshal response", err), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleKafkaTopicsList handles listing all topics
func (b *KafkaTopicsToolBuilder) handleKafkaTopicsList(ctx context.Context, admin *kadm.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	includeInternal := request.GetBool("includeInternal", false)

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
func (b *KafkaTopicsToolBuilder) handleKafkaTopicGet(ctx context.Context, admin *kadm.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := request.RequireString("name")
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
func (b *KafkaTopicsToolBuilder) handleKafkaTopicCreate(ctx context.Context, admin *kadm.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := request.RequireString("name")
	if err != nil {
		return b.handleError("get topic name", err), nil
	}

	partitionsNum, err := request.RequireInt("partitions")
	if err != nil {
		return b.handleError("get partitions", err), nil
	}

	replicationFactorNum, err := request.RequireInt("replicationFactor")
	if err != nil {
		return b.handleError("get replication factor", err), nil
	}

	///nolint:gosec
	partitions := int32(partitionsNum)
	///nolint:gosec
	replicationFactor := int16(replicationFactorNum)

	// Parse optional configs
	var configs map[string]*string
	arguments := request.GetArguments()
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

	// Create topic using the correct CreateTopics API
	results, err := admin.CreateTopics(ctx, partitions, replicationFactor, configs, topicName)
	if err != nil {
		return b.handleError("create Kafka topic", err), nil
	}

	return b.marshalResponse(results)
}

// handleKafkaTopicDelete handles deleting a topic
func (b *KafkaTopicsToolBuilder) handleKafkaTopicDelete(ctx context.Context, admin *kadm.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := request.RequireString("name")
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
func (b *KafkaTopicsToolBuilder) handleKafkaTopicMetadata(ctx context.Context, admin *kadm.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := request.RequireString("name")
	if err != nil {
		return b.handleError("get topic name", err), nil
	}

	metadata, err := admin.Metadata(ctx, topicName)
	if err != nil {
		return b.handleError("get Kafka topic metadata", err), nil
	}

	return b.marshalResponse(metadata)
}
