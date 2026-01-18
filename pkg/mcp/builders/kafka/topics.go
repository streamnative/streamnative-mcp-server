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
	"reflect"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	"github.com/twmb/franz-go/pkg/kadm"
)

type kafkaTopicsInput struct {
	Resource          string         `json:"resource"`
	Operation         string         `json:"operation"`
	Name              *string        `json:"name,omitempty"`
	Partitions        *int           `json:"partitions,omitempty"`
	ReplicationFactor *int           `json:"replicationFactor,omitempty"`
	Configs           map[string]any `json:"configs,omitempty"`
	IncludeInternal   bool           `json:"includeInternal,omitempty"`
}

const (
	kafkaTopicsResourceDesc = "Resource to operate on. Available resources:\n" +
		"- topic: A single Kafka topic for operations on individual topics (create, get, delete)\n" +
		"- topics: Collection of Kafka topics for bulk operations (list)"
	kafkaTopicsOperationDesc = "Operation to perform. Available operations:\n" +
		"- list: List all topics in the Kafka cluster, optionally including internal topics\n" +
		"- get: Get detailed configuration for a specific topic\n" +
		"- create: Create a new topic with specified partitions, replication factor, and optional configs\n" +
		"- delete: Delete an existing topic\n" +
		"- metadata: Get metadata for a specific topic"
	kafkaTopicsNameDesc = "The name of the Kafka topic to operate on. " +
		"Required for 'get', 'create', 'delete', and 'metadata' operations on the 'topic' resource. " +
		"Topic names should follow Kafka naming conventions (alphanumeric, dots, underscores, and hyphens)."
	kafkaTopicsPartitionsDesc = "The number of partitions for the topic. Required for 'create' operation. " +
		"Partitions determine the parallelism and scalability of the topic. " +
		"More partitions allow more concurrent consumers and higher throughput."
	kafkaTopicsReplicationFactorDesc = "The replication factor for the topic. Required for 'create' operation. " +
		"Replication factor determines fault tolerance - it should be at least 2 for production use. " +
		"Cannot exceed the number of available brokers in the cluster."
	kafkaTopicsConfigsDesc = "Optional configuration parameters for the topic during 'create' operation. " +
		"Common configurations include:\n" +
		"- retention.ms: How long to retain messages (milliseconds)\n" +
		"- compression.type: Compression algorithm (none, gzip, snappy, lz4, zstd)\n" +
		"- cleanup.policy: Log cleanup policy (delete, compact, compact,delete)\n" +
		"- segment.ms: Time before a new log segment is rolled out\n" +
		"- max.message.bytes: Maximum size of a message batch"
	kafkaTopicsIncludeInternalDesc = "Whether to include internal Kafka topics in the 'list' operation. " +
		"Internal topics are used by Kafka itself (e.g., __consumer_offsets, __transaction_state). " +
		"Default: false"
)

// KafkaTopicsToolBuilder implements the ToolBuilder interface for Kafka topics.
type KafkaTopicsToolBuilder struct { //nolint:revive
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
func (b *KafkaTopicsToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool, err := b.buildKafkaTopicsTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildKafkaTopicsHandler(config.ReadOnly)

	return []builders.ToolDefinition{
		builders.ServerTool[kafkaTopicsInput, any]{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildKafkaTopicsTool builds the Kafka Topics MCP tool definition
func (b *KafkaTopicsToolBuilder) buildKafkaTopicsTool() (*sdk.Tool, error) {
	inputSchema, err := buildKafkaTopicsInputSchema()
	if err != nil {
		return nil, err
	}

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

	return &sdk.Tool{
		Name:        "kafka_admin_topics",
		Description: toolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildKafkaTopicsHandler builds the Kafka Topics handler function
func (b *KafkaTopicsToolBuilder) buildKafkaTopicsHandler(readOnly bool) builders.ToolHandlerFunc[kafkaTopicsInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input kafkaTopicsInput) (*sdk.CallToolResult, any, error) {
		// Normalize parameters
		resource := strings.ToLower(input.Resource)
		operation := strings.ToLower(input.Operation)

		// Validate write operations in read-only mode
		if readOnly && (operation == "create" || operation == "delete") {
			return nil, nil, fmt.Errorf("write operations are not allowed in read-only mode")
		}

		// Get Kafka admin client
		session := mcpCtx.GetKafkaSession(ctx)
		if session == nil {
			return nil, nil, b.handleError("get Kafka session not found in context", nil)
		}
		admin, err := session.GetAdminClient()
		if err != nil {
			return nil, nil, b.handleError("get admin client", err)
		}

		// Dispatch based on resource and operation
		switch resource {
		case "topics":
			switch operation {
			case "list":
				result, err := b.handleKafkaTopicsList(ctx, admin, input)
				return result, nil, err
			default:
				return nil, nil, fmt.Errorf("invalid operation for resource 'topics': %s", operation)
			}
		case "topic":
			switch operation {
			case "get":
				result, err := b.handleKafkaTopicGet(ctx, admin, input)
				return result, nil, err
			case "create":
				result, err := b.handleKafkaTopicCreate(ctx, admin, input)
				return result, nil, err
			case "delete":
				result, err := b.handleKafkaTopicDelete(ctx, admin, input)
				return result, nil, err
			case "metadata":
				result, err := b.handleKafkaTopicMetadata(ctx, admin, input)
				return result, nil, err
			default:
				return nil, nil, fmt.Errorf("invalid operation for resource 'topic': %s", operation)
			}
		default:
			return nil, nil, fmt.Errorf("invalid resource: %s. available resources: topics, topic", resource)
		}
	}
}

// Utility functions

// handleError provides unified error handling
func (b *KafkaTopicsToolBuilder) handleError(operation string, err error) error {
	return fmt.Errorf("failed to %s: %v", operation, err)
}

// marshalResponse provides unified JSON serialization for responses
func (b *KafkaTopicsToolBuilder) marshalResponse(data interface{}) (*sdk.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, b.handleError("marshal response", err)
	}
	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: string(jsonBytes)}},
	}, nil
}

func requireString(value *string, key string) (string, error) {
	if value == nil {
		return "", fmt.Errorf("required argument %q not found", key)
	}
	return *value, nil
}

func requireInt(value *int, key string) (int, error) {
	if value == nil {
		return 0, fmt.Errorf("required argument %q not found", key)
	}
	return *value, nil
}

// handleKafkaTopicsList handles listing all topics
func (b *KafkaTopicsToolBuilder) handleKafkaTopicsList(ctx context.Context, admin *kadm.Client, input kafkaTopicsInput) (*sdk.CallToolResult, error) {
	includeInternal := input.IncludeInternal

	topics, err := admin.ListTopics(ctx)
	if err != nil {
		return nil, b.handleError("list Kafka topics", err)
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
func (b *KafkaTopicsToolBuilder) handleKafkaTopicGet(ctx context.Context, admin *kadm.Client, input kafkaTopicsInput) (*sdk.CallToolResult, error) {
	topicName, err := requireString(input.Name, "name")
	if err != nil {
		return nil, b.handleError("get topic name", err)
	}

	topics, err := admin.ListTopics(ctx, topicName)
	if err != nil {
		return nil, b.handleError("get Kafka topic", err)
	}

	return b.marshalResponse(topics)
}

// handleKafkaTopicCreate handles creating a new topic
func (b *KafkaTopicsToolBuilder) handleKafkaTopicCreate(ctx context.Context, admin *kadm.Client, input kafkaTopicsInput) (*sdk.CallToolResult, error) {
	topicName, err := requireString(input.Name, "name")
	if err != nil {
		return nil, b.handleError("get topic name", err)
	}

	partitionsNum, err := requireInt(input.Partitions, "partitions")
	if err != nil {
		return nil, b.handleError("get partitions", err)
	}

	replicationFactorNum, err := requireInt(input.ReplicationFactor, "replicationFactor")
	if err != nil {
		return nil, b.handleError("get replication factor", err)
	}

	//nolint:gosec
	partitions := int32(partitionsNum)
	//nolint:gosec
	replicationFactor := int16(replicationFactorNum)

	configs := b.buildConfigs(input)

	// Create topic using the correct CreateTopics API
	results, err := admin.CreateTopics(ctx, partitions, replicationFactor, configs, topicName)
	if err != nil {
		return nil, b.handleError("create Kafka topic", err)
	}

	return b.marshalResponse(results)
}

// handleKafkaTopicDelete handles deleting a topic
func (b *KafkaTopicsToolBuilder) handleKafkaTopicDelete(ctx context.Context, admin *kadm.Client, input kafkaTopicsInput) (*sdk.CallToolResult, error) {
	topicName, err := requireString(input.Name, "name")
	if err != nil {
		return nil, b.handleError("get topic name", err)
	}

	results, err := admin.DeleteTopics(ctx, topicName)
	if err != nil {
		return nil, b.handleError("delete Kafka topic", err)
	}

	return b.marshalResponse(results)
}

// handleKafkaTopicMetadata handles getting metadata for a topic
func (b *KafkaTopicsToolBuilder) handleKafkaTopicMetadata(ctx context.Context, admin *kadm.Client, input kafkaTopicsInput) (*sdk.CallToolResult, error) {
	topicName, err := requireString(input.Name, "name")
	if err != nil {
		return nil, b.handleError("get topic name", err)
	}

	metadata, err := admin.Metadata(ctx, topicName)
	if err != nil {
		return nil, b.handleError("get Kafka topic metadata", err)
	}

	return b.marshalResponse(metadata)
}

func (b *KafkaTopicsToolBuilder) buildConfigs(input kafkaTopicsInput) map[string]*string {
	if len(input.Configs) == 0 {
		return nil
	}

	configs := make(map[string]*string, len(input.Configs))
	for key, value := range input.Configs {
		strValue, ok := value.(string)
		if !ok {
			strValue = fmt.Sprintf("%v", value)
		}
		configs[key] = &strValue
	}

	return configs
}

func buildKafkaTopicsInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[kafkaTopicsInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}

	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}
	setSchemaDescription(schema, "resource", kafkaTopicsResourceDesc)
	setSchemaDescription(schema, "operation", kafkaTopicsOperationDesc)
	setSchemaDescription(schema, "name", kafkaTopicsNameDesc)
	setSchemaDescription(schema, "partitions", kafkaTopicsPartitionsDesc)
	setSchemaDescription(schema, "replicationFactor", kafkaTopicsReplicationFactorDesc)
	setSchemaDescription(schema, "configs", kafkaTopicsConfigsDesc)
	setSchemaDescription(schema, "includeInternal", kafkaTopicsIncludeInternalDesc)

	normalizeAdditionalProperties(schema)
	return schema, nil
}

func setSchemaDescription(schema *jsonschema.Schema, name, desc string) {
	if schema == nil {
		return
	}
	prop, ok := schema.Properties[name]
	if !ok || prop == nil {
		return
	}
	prop.Description = desc
}

func normalizeAdditionalProperties(schema *jsonschema.Schema) {
	visited := map[*jsonschema.Schema]bool{}
	var walk func(*jsonschema.Schema)
	walk = func(s *jsonschema.Schema) {
		if s == nil || visited[s] {
			return
		}
		visited[s] = true

		if s.Type == "object" && s.Properties != nil && isFalseSchema(s.AdditionalProperties) {
			s.AdditionalProperties = nil
		}

		for _, prop := range s.Properties {
			walk(prop)
		}
		for _, prop := range s.PatternProperties {
			walk(prop)
		}
		for _, def := range s.Defs {
			walk(def)
		}
		for _, def := range s.Definitions {
			walk(def)
		}
		if s.AdditionalProperties != nil && !isFalseSchema(s.AdditionalProperties) {
			walk(s.AdditionalProperties)
		}
		if s.Items != nil {
			walk(s.Items)
		}
		for _, item := range s.PrefixItems {
			walk(item)
		}
		if s.AdditionalItems != nil {
			walk(s.AdditionalItems)
		}
		if s.UnevaluatedItems != nil {
			walk(s.UnevaluatedItems)
		}
		if s.UnevaluatedProperties != nil {
			walk(s.UnevaluatedProperties)
		}
		if s.PropertyNames != nil {
			walk(s.PropertyNames)
		}
		if s.Contains != nil {
			walk(s.Contains)
		}
		for _, subschema := range s.AllOf {
			walk(subschema)
		}
		for _, subschema := range s.AnyOf {
			walk(subschema)
		}
		for _, subschema := range s.OneOf {
			walk(subschema)
		}
		if s.Not != nil {
			walk(s.Not)
		}
		if s.If != nil {
			walk(s.If)
		}
		if s.Then != nil {
			walk(s.Then)
		}
		if s.Else != nil {
			walk(s.Else)
		}
		for _, subschema := range s.DependentSchemas {
			walk(subschema)
		}
	}
	walk(schema)
}

func isFalseSchema(schema *jsonschema.Schema) bool {
	if schema == nil || schema.Not == nil {
		return false
	}
	if !reflect.ValueOf(*schema.Not).IsZero() {
		return false
	}
	clone := *schema
	clone.Not = nil
	return reflect.ValueOf(clone).IsZero()
}
