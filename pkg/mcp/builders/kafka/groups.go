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

	"github.com/google/jsonschema-go/jsonschema"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	"github.com/twmb/franz-go/pkg/kadm"
)

type kafkaGroupsInput struct {
	Resource  string  `json:"resource"`
	Operation string  `json:"operation"`
	Group     *string `json:"group,omitempty"`
	Members   *string `json:"members,omitempty"`
	Topic     *string `json:"topic,omitempty"`
	Partition *int    `json:"partition,omitempty"`
	Offset    *int64  `json:"offset,omitempty"`
}

const (
	kafkaGroupsResourceDesc = "Resource to operate on. Available resources:\n" +
		"- group: A single Kafka Consumer Group for operations on individual groups (describe, remove-members, set-offset, delete-offset)\n" +
		"- groups: Collection of Kafka Consumer Groups for bulk operations (list)"
	kafkaGroupsOperationDesc = "Operation to perform. Available operations:\n" +
		"- list: List all Kafka Consumer Groups in the cluster\n" +
		"- describe: Get detailed information about a specific Consumer Group, including members, offsets, and lag\n" +
		"- remove-members: Remove specific members from a Consumer Group to force rebalancing or troubleshoot issues\n" +
		"- offsets: Get offsets for a specific consumer group\n" +
		"- delete-offset: Delete a specific offset for a consumer group of a topic\n" +
		"- set-offset: Set a specific offset for a consumer group's topic-partition"
	kafkaGroupsGroupDesc = "The name of the Kafka Consumer Group to operate on. " +
		"Required for the 'describe' and 'remove-members' operations. " +
		"Must be an existing consumer group name in the Kafka cluster. " +
		"Consumer Group names are case-sensitive and typically follow a naming convention like 'application-name'."
	kafkaGroupsMembersDesc = "Comma-separated list of consumer instance IDs to remove from the group. " +
		"Required for the 'remove-members' operation. " +
		"Consumer instance IDs can be found using the 'describe' operation."
	kafkaGroupsTopicDesc     = "The topic name. Required for 'delete-offset' and 'set-offset' operations."
	kafkaGroupsPartitionDesc = "The partition number. Required for 'set-offset' operation."
	kafkaGroupsOffsetDesc    = "The offset value to set. Required for 'set-offset' operation."
)

// KafkaGroupsToolBuilder implements the ToolBuilder interface for Kafka Consumer Groups
// /nolint:revive
type KafkaGroupsToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewKafkaGroupsToolBuilder creates a new Kafka Groups tool builder instance
func NewKafkaGroupsToolBuilder() *KafkaGroupsToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "kafka_groups",
		Version:     "1.0.0",
		Description: "Kafka Consumer Groups administration tools",
		Category:    "kafka_admin",
		Tags:        []string{"kafka", "groups", "admin", "consumer"},
	}

	features := []string{
		"kafka-admin",
		"all",
		"all-kafka",
	}

	return &KafkaGroupsToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Kafka Groups tool list
func (b *KafkaGroupsToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool, err := b.buildKafkaGroupsTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildKafkaGroupsHandler(config.ReadOnly)

	return []builders.ToolDefinition{
		builders.ServerTool[kafkaGroupsInput, any]{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildKafkaGroupsTool builds the Kafka Groups MCP tool definition
func (b *KafkaGroupsToolBuilder) buildKafkaGroupsTool() (*sdk.Tool, error) {
	inputSchema, err := buildKafkaGroupsInputSchema()
	if err != nil {
		return nil, err
	}

	toolDesc := "Unified tool for managing Apache Kafka Consumer Groups.\n" +
		"This tool provides access to Kafka consumer group operations including listing, describing, and managing group membership.\n" +
		"Kafka Consumer Groups are a key concept for scalable consumption of Kafka topics. A consumer group consists of multiple consumer instances\n" +
		"that collaborate to consume data from topic partitions. Kafka ensures that:\n" +
		"- Each partition is consumed by exactly one consumer in the group\n" +
		"- When consumers join or leave, Kafka triggers a 'rebalance' to redistribute partitions\n" +
		"- Consumer groups track consumption progress through committed offsets\n\n" +
		"Usage Examples:\n\n" +
		"1. List all Kafka Consumer Groups in the cluster:\n" +
		"   resource: \"groups\"\n" +
		"   operation: \"list\"\n\n" +
		"2. Describe a specific Kafka Consumer Group to see its members and consumption details:\n" +
		"   resource: \"group\"\n" +
		"   operation: \"describe\"\n" +
		"   group: \"my-consumer-group\"\n\n" +
		"3. Remove specific members from a Kafka Consumer Group to trigger rebalancing:\n" +
		"   resource: \"group\"\n" +
		"   operation: \"remove-members\"\n" +
		"   group: \"my-consumer-group\"\n" +
		"   members: \"consumer-instance-1,consumer-instance-2\"\n\n" +
		"4. Get offsets for a specific consumer group:\n" +
		"   resource: \"group\"\n" +
		"   operation: \"offsets\"\n" +
		"   group: \"my-consumer-group\"\n\n" +
		"5. Delete a specific offset for a consumer group of a topic:\n" +
		"   resource: \"group\"\n" +
		"   operation: \"delete-offset\"\n" +
		"   group: \"my-consumer-group\"\n" +
		"   topic: \"my-topic\"\n\n" +
		"6. Set a specific offset for a consumer group's topic-partition:\n" +
		"   resource: \"group\"\n" +
		"   operation: \"set-offset\"\n" +
		"   group: \"my-consumer-group\"\n" +
		"   topic: \"my-topic\"\n" +
		"   partition: 0\n" +
		"   offset: 1000\n\n" +
		"This tool requires Kafka super-user permissions."

	return &sdk.Tool{
		Name:        "kafka_admin_groups",
		Description: toolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildKafkaGroupsHandler builds the Kafka Groups handler function
func (b *KafkaGroupsToolBuilder) buildKafkaGroupsHandler(readOnly bool) builders.ToolHandlerFunc[kafkaGroupsInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input kafkaGroupsInput) (*sdk.CallToolResult, any, error) {
		// Normalize parameters
		resource := strings.ToLower(input.Resource)
		operation := strings.ToLower(input.Operation)

		// Validate write operations in read-only mode
		if readOnly && (operation == "remove-members" || operation == "delete-offset" || operation == "set-offset") {
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
		case "groups":
			switch operation {
			case "list":
				result, err := b.handleKafkaGroupsList(ctx, admin)
				return result, nil, err
			default:
				return nil, nil, fmt.Errorf("invalid operation for resource 'groups': %s", operation)
			}
		case "group":
			switch operation {
			case "describe":
				result, err := b.handleKafkaGroupDescribe(ctx, admin, input)
				return result, nil, err
			case "remove-members":
				result, err := b.handleKafkaGroupRemoveMembers(ctx, admin, input)
				return result, nil, err
			case "offsets":
				result, err := b.handleKafkaGroupOffsets(ctx, admin, input)
				return result, nil, err
			case "delete-offset":
				result, err := b.handleKafkaGroupDeleteOffset(ctx, admin, input)
				return result, nil, err
			case "set-offset":
				result, err := b.handleKafkaGroupSetOffset(ctx, admin, input)
				return result, nil, err
			default:
				return nil, nil, fmt.Errorf("invalid operation for resource 'group': %s", operation)
			}
		default:
			return nil, nil, fmt.Errorf("invalid resource: %s. available resources: groups, group", resource)
		}
	}
}

// Utility functions

// handleError provides unified error handling
func (b *KafkaGroupsToolBuilder) handleError(operation string, err error) error {
	return fmt.Errorf("failed to %s: %v", operation, err)
}

// marshalResponse provides unified JSON serialization for responses
func (b *KafkaGroupsToolBuilder) marshalResponse(data interface{}) (*sdk.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, b.handleError("marshal response", err)
	}
	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: string(jsonBytes)}},
	}, nil
}

func requireInt64(value *int64, key string) (int64, error) {
	if value == nil {
		return 0, fmt.Errorf("required argument %q not found", key)
	}
	return *value, nil
}

// Specific operation handler functions

// handleKafkaGroupsList handles listing all consumer groups
func (b *KafkaGroupsToolBuilder) handleKafkaGroupsList(ctx context.Context, admin *kadm.Client) (*sdk.CallToolResult, error) {
	groups, err := admin.ListGroups(ctx)
	if err != nil {
		return nil, b.handleError("list Kafka consumer groups", err)
	}
	return b.marshalResponse(groups)
}

// handleKafkaGroupDescribe handles describing a specific consumer group
func (b *KafkaGroupsToolBuilder) handleKafkaGroupDescribe(ctx context.Context, admin *kadm.Client, input kafkaGroupsInput) (*sdk.CallToolResult, error) {
	groupName, err := requireString(input.Group, "group")
	if err != nil {
		return nil, b.handleError("get group name", err)
	}

	groups, err := admin.DescribeGroups(ctx, groupName)
	if err != nil {
		return nil, b.handleError("describe Kafka consumer group", err)
	}
	return b.marshalResponse(groups)
}

// handleKafkaGroupRemoveMembers handles removing members from a consumer group
func (b *KafkaGroupsToolBuilder) handleKafkaGroupRemoveMembers(ctx context.Context, admin *kadm.Client, input kafkaGroupsInput) (*sdk.CallToolResult, error) {
	groupName, err := requireString(input.Group, "group")
	if err != nil {
		return nil, b.handleError("get group name", err)
	}

	membersStr, err := requireString(input.Members, "members")
	if err != nil {
		return nil, b.handleError("get members", err)
	}

	memberIDs := strings.Split(membersStr, ",")
	for i, member := range memberIDs {
		memberIDs[i] = strings.TrimSpace(member)
	}

	builder := kadm.LeaveGroup(groupName).InstanceIDs(memberIDs...)
	response, err := admin.LeaveGroup(ctx, builder)
	if err != nil {
		return nil, b.handleError("remove members from Kafka consumer group", err)
	}

	return b.marshalResponse(response)
}

// handleKafkaGroupOffsets handles getting offsets for a consumer group
func (b *KafkaGroupsToolBuilder) handleKafkaGroupOffsets(ctx context.Context, admin *kadm.Client, input kafkaGroupsInput) (*sdk.CallToolResult, error) {
	groupName, err := requireString(input.Group, "group")
	if err != nil {
		return nil, b.handleError("get group name", err)
	}

	response, err := admin.FetchOffsets(ctx, groupName)
	if err != nil {
		return nil, b.handleError("get offsets for Kafka consumer group", err)
	}
	return b.marshalResponse(response)
}

// handleKafkaGroupDeleteOffset handles deleting a specific offset for a consumer group
func (b *KafkaGroupsToolBuilder) handleKafkaGroupDeleteOffset(ctx context.Context, admin *kadm.Client, input kafkaGroupsInput) (*sdk.CallToolResult, error) {
	groupName, err := requireString(input.Group, "group")
	if err != nil {
		return nil, b.handleError("get group name", err)
	}

	topicName, err := requireString(input.Topic, "topic")
	if err != nil {
		return nil, b.handleError("get topic name", err)
	}

	// Create a TopicsSet with the specified topic
	// This will target all partitions for the given topic
	topicsSet := make(kadm.TopicsSet)
	topicsSet.Add(topicName)

	// Call DeleteOffsets to delete the offsets for the specified topic
	responses, err := admin.DeleteOffsets(ctx, groupName, topicsSet)
	if err != nil {
		return nil, b.handleError("delete offset for Kafka consumer group", err)
	}

	// Check for errors in the response
	if err := responses.Error(); err != nil {
		return nil, b.handleError("delete offset for Kafka consumer group", err)
	}

	return b.marshalResponse(responses)
}

// handleKafkaGroupSetOffset handles setting a specific offset for a consumer group
func (b *KafkaGroupsToolBuilder) handleKafkaGroupSetOffset(ctx context.Context, admin *kadm.Client, input kafkaGroupsInput) (*sdk.CallToolResult, error) {
	groupName, err := requireString(input.Group, "group")
	if err != nil {
		return nil, b.handleError("get group name", err)
	}

	topicName, err := requireString(input.Topic, "topic")
	if err != nil {
		return nil, b.handleError("get topic name", err)
	}

	partitionNum, err := requireInt(input.Partition, "partition")
	if err != nil {
		return nil, b.handleError("get partition number", err)
	}
	//nolint:gosec
	partitionInt := int32(partitionNum)

	offsetNum, err := requireInt64(input.Offset, "offset")
	if err != nil {
		return nil, b.handleError("get offset", err)
	}

	// Create Offsets object with the specified topic, partition, and offset
	offsets := make(kadm.Offsets)
	offsets.AddOffset(topicName, partitionInt, offsetNum, -1) // Using -1 for leaderEpoch as it's optional

	// Commit the offsets
	responses, err := admin.CommitOffsets(ctx, groupName, offsets)
	if err != nil {
		return nil, b.handleError("set offset for Kafka consumer group", err)
	}

	// Check for errors in the response
	if err := responses.Error(); err != nil {
		return nil, b.handleError("set offset for Kafka consumer group", err)
	}

	return b.marshalResponse(responses)
}

func buildKafkaGroupsInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[kafkaGroupsInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}

	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}
	setSchemaDescription(schema, "resource", kafkaGroupsResourceDesc)
	setSchemaDescription(schema, "operation", kafkaGroupsOperationDesc)
	setSchemaDescription(schema, "group", kafkaGroupsGroupDesc)
	setSchemaDescription(schema, "members", kafkaGroupsMembersDesc)
	setSchemaDescription(schema, "topic", kafkaGroupsTopicDesc)
	setSchemaDescription(schema, "partition", kafkaGroupsPartitionDesc)
	setSchemaDescription(schema, "offset", kafkaGroupsOffsetDesc)

	normalizeAdditionalProperties(schema)
	return schema, nil
}
