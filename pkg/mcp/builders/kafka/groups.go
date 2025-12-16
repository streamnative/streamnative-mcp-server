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

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	"github.com/twmb/franz-go/pkg/kadm"
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
func (b *KafkaGroupsToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool := b.buildKafkaGroupsTool()
	handler := b.buildKafkaGroupsHandler(config.ReadOnly)

	return []server.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildKafkaGroupsTool builds the Kafka Groups MCP tool definition
func (b *KafkaGroupsToolBuilder) buildKafkaGroupsTool() mcp.Tool {
	resourceDesc := "Resource to operate on. Available resources:\n" +
		"- group: A single Kafka Consumer Group for operations on individual groups (describe, remove-members, set-offset, delete-offset)\n" +
		"- groups: Collection of Kafka Consumer Groups for bulk operations (list)"

	operationDesc := "Operation to perform. Available operations:\n" +
		"- list: List all Kafka Consumer Groups in the cluster\n" +
		"- describe: Get detailed information about a specific Consumer Group, including members, offsets, and lag\n" +
		"- remove-members: Remove specific members from a Consumer Group to force rebalancing or troubleshoot issues\n" +
		"- offsets: Get offsets for a specific consumer group\n" +
		"- delete-offset: Delete a specific offset for a consumer group of a topic\n" +
		"- set-offset: Set a specific offset for a consumer group's topic-partition"

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

	return mcp.NewTool("kafka_admin_groups",
		mcp.WithDescription(toolDesc),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description(resourceDesc),
		),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
		),
		mcp.WithString("group",
			mcp.Description("The name of the Kafka Consumer Group to operate on. "+
				"Required for the 'describe' and 'remove-members' operations. "+
				"Must be an existing consumer group name in the Kafka cluster. "+
				"Consumer Group names are case-sensitive and typically follow a naming convention like 'application-name'.")),
		mcp.WithString("members",
			mcp.Description("Comma-separated list of consumer instance IDs to remove from the group. "+
				"Required for the 'remove-members' operation. "+
				"Consumer instance IDs can be found using the 'describe' operation.")),
		mcp.WithString("topic",
			mcp.Description("The topic name. Required for 'delete-offset' and 'set-offset' operations.")),
		mcp.WithNumber("partition",
			mcp.Description("The partition number. Required for 'set-offset' operation.")),
		mcp.WithNumber("offset",
			mcp.Description("The offset value to set. Required for 'set-offset' operation.")),
	)
}

// buildKafkaGroupsHandler builds the Kafka Groups handler function
func (b *KafkaGroupsToolBuilder) buildKafkaGroupsHandler(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		// Validate write operations in read-only mode
		if readOnly && (operation == "remove-members" || operation == "delete-offset" || operation == "set-offset") {
			return mcp.NewToolResultError("Write operations are not allowed in read-only mode"), nil
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
		case "groups":
			switch operation {
			case "list":
				return b.handleKafkaGroupsList(ctx, admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'groups': %s", operation)), nil
			}
		case "group":
			switch operation {
			case "describe":
				return b.handleKafkaGroupDescribe(ctx, admin, request)
			case "remove-members":
				return b.handleKafkaGroupRemoveMembers(ctx, admin, request)
			case "offsets":
				return b.handleKafkaGroupOffsets(ctx, admin, request)
			case "delete-offset":
				return b.handleKafkaGroupDeleteOffset(ctx, admin, request)
			case "set-offset":
				return b.handleKafkaGroupSetOffset(ctx, admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'group': %s", operation)), nil
			}
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resource: %s. Available resources: groups, group", resource)), nil
		}
	}
}

// Utility functions

// handleError provides unified error handling
func (b *KafkaGroupsToolBuilder) handleError(operation string, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("Failed to %s: %v", operation, err))
}

// marshalResponse provides unified JSON serialization for responses
func (b *KafkaGroupsToolBuilder) marshalResponse(data interface{}) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return b.handleError("marshal response", err), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// Specific operation handler functions

// handleKafkaGroupsList handles listing all consumer groups
func (b *KafkaGroupsToolBuilder) handleKafkaGroupsList(ctx context.Context, admin *kadm.Client, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groups, err := admin.ListGroups(ctx)
	if err != nil {
		return b.handleError("list Kafka consumer groups", err), nil
	}
	return b.marshalResponse(groups)
}

// handleKafkaGroupDescribe handles describing a specific consumer group
func (b *KafkaGroupsToolBuilder) handleKafkaGroupDescribe(ctx context.Context, admin *kadm.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groupName, err := request.RequireString("group")
	if err != nil {
		return b.handleError("get group name", err), nil
	}

	groups, err := admin.DescribeGroups(ctx, groupName)
	if err != nil {
		return b.handleError("describe Kafka consumer group", err), nil
	}
	return b.marshalResponse(groups)
}

// handleKafkaGroupRemoveMembers handles removing members from a consumer group
func (b *KafkaGroupsToolBuilder) handleKafkaGroupRemoveMembers(ctx context.Context, admin *kadm.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groupName, err := request.RequireString("group")
	if err != nil {
		return b.handleError("get group name", err), nil
	}

	membersStr, err := request.RequireString("members")
	if err != nil {
		return b.handleError("get members", err), nil
	}

	memberIDs := strings.Split(membersStr, ",")
	for i, member := range memberIDs {
		memberIDs[i] = strings.TrimSpace(member)
	}

	builder := kadm.LeaveGroup(groupName).InstanceIDs(memberIDs...)
	response, err := admin.LeaveGroup(ctx, builder)
	if err != nil {
		return b.handleError("remove members from Kafka consumer group", err), nil
	}

	return b.marshalResponse(response)
}

// handleKafkaGroupOffsets handles getting offsets for a consumer group
func (b *KafkaGroupsToolBuilder) handleKafkaGroupOffsets(ctx context.Context, admin *kadm.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groupName, err := request.RequireString("group")
	if err != nil {
		return b.handleError("get group name", err), nil
	}

	response, err := admin.FetchOffsets(ctx, groupName)
	if err != nil {
		return b.handleError("get offsets for Kafka consumer group", err), nil
	}
	return b.marshalResponse(response)
}

// handleKafkaGroupDeleteOffset handles deleting a specific offset for a consumer group
func (b *KafkaGroupsToolBuilder) handleKafkaGroupDeleteOffset(ctx context.Context, admin *kadm.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groupName, err := request.RequireString("group")
	if err != nil {
		return b.handleError("get group name", err), nil
	}

	topicName, err := request.RequireString("topic")
	if err != nil {
		return b.handleError("get topic name", err), nil
	}

	// Create a TopicsSet with the specified topic
	// This will target all partitions for the given topic
	topicsSet := make(kadm.TopicsSet)
	topicsSet.Add(topicName)

	// Call DeleteOffsets to delete the offsets for the specified topic
	responses, err := admin.DeleteOffsets(ctx, groupName, topicsSet)
	if err != nil {
		return b.handleError("delete offset for Kafka consumer group", err), nil
	}

	// Check for errors in the response
	if err := responses.Error(); err != nil {
		return b.handleError("delete offset for Kafka consumer group", err), nil
	}

	return b.marshalResponse(responses)
}

// handleKafkaGroupSetOffset handles setting a specific offset for a consumer group
func (b *KafkaGroupsToolBuilder) handleKafkaGroupSetOffset(ctx context.Context, admin *kadm.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groupName, err := request.RequireString("group")
	if err != nil {
		return b.handleError("get group name", err), nil
	}

	topicName, err := request.RequireString("topic")
	if err != nil {
		return b.handleError("get topic name", err), nil
	}

	partitionNum, err := request.RequireFloat("partition")
	if err != nil {
		return b.handleError("get partition number", err), nil
	}
	partitionInt := int32(partitionNum)

	offsetNum, err := request.RequireFloat("offset")
	if err != nil {
		return b.handleError("get offset", err), nil
	}
	offsetInt := int64(offsetNum)

	// Create Offsets object with the specified topic, partition, and offset
	offsets := make(kadm.Offsets)
	offsets.AddOffset(topicName, partitionInt, offsetInt, -1) // Using -1 for leaderEpoch as it's optional

	// Commit the offsets
	responses, err := admin.CommitOffsets(ctx, groupName, offsets)
	if err != nil {
		return b.handleError("set offset for Kafka consumer group", err), nil
	}

	// Check for errors in the response
	if err := responses.Error(); err != nil {
		return b.handleError("set offset for Kafka consumer group", err), nil
	}

	return b.marshalResponse(responses)
}
