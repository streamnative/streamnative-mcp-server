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

package pulsar

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/toolannotations"
)

var readOnlyRestrictedTopicOperations = map[string]struct{}{
	"create":             {},
	"delete":             {},
	"unload":             {},
	"terminate":          {},
	"compact":            {},
	"update":             {},
	"offload":            {},
	"grant-permissions":  {},
	"revoke-permissions": {},
}

var topicOperationAliases = map[string]string{
	"status": "compact-status",
}

const topicStatusPollInterval = time.Second

// PulsarAdminTopicToolBuilder implements the ToolBuilder interface for Pulsar Admin Topic tools
// It provides functionality to build Pulsar topic management tools
// /nolint:revive
type PulsarAdminTopicToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminTopicToolBuilder creates a new Pulsar Admin Topic tool builder instance
func NewPulsarAdminTopicToolBuilder() *PulsarAdminTopicToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_topic",
		Version:     "1.0.0",
		Description: "Pulsar Admin topic management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "topic", "admin"},
	}

	features := []string{
		"pulsar-admin-topics",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminTopicToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar Admin Topic tool list
// This is the core method implementing the ToolBuilder interface
func (b *PulsarAdminTopicToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	tools := []server.ServerTool{
		{
			Tool:    b.buildTopicTool(toolModeRead),
			Handler: b.buildTopicHandler(toolModeRead),
		},
	}
	if !config.ReadOnly {
		tools = append(tools, server.ServerTool{
			Tool:    b.buildTopicTool(toolModeWrite),
			Handler: b.buildTopicHandler(toolModeWrite),
		})
	}

	return tools, nil
}

// buildTopicTool builds the Pulsar Admin Topic MCP tool definition
// Migrated from the original tool definition logic
func (b *PulsarAdminTopicToolBuilder) buildTopicTool(mode toolMode) mcp.Tool {
	toolDesc := "Manage Apache Pulsar topics. " +
		"Topics are the core messaging entities in Pulsar that store and transmit messages. " +
		"Pulsar supports two types of topics: persistent (durable storage with guaranteed delivery) " +
		"and non-persistent (in-memory with at-most-once delivery). " +
		"Topics can be partitioned for parallel processing and higher throughput, where each partition " +
		"functions as an independent topic with its own message log. " +
		"Topics follow a hierarchical naming structure: persistent://tenant/namespace/topic. " +
		"This tool supports various operations on topics including creation, deletion, lookup, compaction, " +
		"offloading, and retrieving statistics. " +
		"Do not use this tool for Kafka protocol operations. Use 'kafka_admin_topics' instead." +
		"Most operations require namespace admin permissions."

	resourceDesc := "Resource to operate on. Available resources:\n" +
		"- topic: A Pulsar topic\n" +
		"- topics: Multiple topics within a namespace"

	operationDesc := "Operation to perform. Available operations:\n" +
		"- list: List all topics in a namespace\n" +
		"- get: Get metadata for a topic\n" +
		"- get-permissions: Get topic permissions for all roles\n" +
		"- grant-permissions: Grant topic permissions to a role\n" +
		"- revoke-permissions: Revoke topic permissions from a role\n" +
		"- create: Create a new topic with optional partitions\n" +
		"- delete: Delete a topic\n" +
		"- stats: Get stats for a topic\n" +
		"- lookup: Look up the broker serving a topic\n" +
		"- internal-stats: Get internal stats for a topic\n" +
		"- internal-info: Get internal info for a topic\n" +
		"- bundle-range: Get the bundle range of a topic\n" +
		"- last-message-id: Get the last message ID of a topic\n" +
		"- compact-status: Get compaction status for a topic (legacy alias: status)\n" +
		"- unload: Unload a topic\n" +
		"- terminate: Terminate a topic\n" +
		"- compact: Trigger compaction on a topic\n" +
		"- update: Update a topic partitions\n" +
		"- offload: Offload data from a topic to long-term storage\n" +
		"- offload-status: Check the status of data offloading for a topic"

	operationEnum := []string{"list", "get", "get-permissions", "stats", "lookup", "internal-stats", "internal-info", "bundle-range", "last-message-id", "compact-status", "offload-status"}
	toolName := "pulsar_admin_topic_read"
	annotation := toolannotations.ReadOnly("Read Pulsar Topics")
	if isToolModeWrite(mode) {
		operationEnum = []string{"grant-permissions", "revoke-permissions", "create", "delete", "unload", "terminate", "compact", "update", "offload"}
		toolName = "pulsar_admin_topic_write"
		annotation = toolannotations.Destructive("Manage Pulsar Topics")
	}

	return mcp.NewTool(toolName,
		mcp.WithDescription(toolDesc),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description(resourceDesc),
		),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
			mcp.Enum(operationEnum...),
		),
		mcp.WithString("topic",
			mcp.Description("The fully qualified topic name (format: [persistent|non-persistent]://tenant/namespace/topic). "+
				"Required for all operations except 'list'. "+
				"For partitioned topics, reference the base topic name without the partition suffix. "+
				"To operate on a specific partition, append -partition-N to the topic name."),
		),
		mcp.WithString("namespace",
			mcp.Description("The namespace name in the format 'tenant/namespace'. "+
				"Required for the 'list' operation. "+
				"A namespace is a logical grouping of topics within a tenant."),
		),
		mcp.WithNumber("partitions",
			mcp.Description("The number of partitions for the topic. Required for 'create' and 'update' operations. "+
				"Set to 0 for a non-partitioned topic. "+
				"Partitioned topics provide higher throughput by dividing message traffic across multiple brokers. "+
				"Each partition is an independent unit with its own retention and cursor positions."),
		),
		mcp.WithBoolean("force",
			mcp.Description("Force operation even if it disrupts producers or consumers. Optional for 'delete' operation. "+
				"When true, all producers and consumers will be forcefully disconnected. "+
				"Use with caution as it can interrupt active message processing."),
		),
		mcp.WithBoolean("non-partitioned",
			mcp.Description("Operate on a non-partitioned topic. Optional for 'delete' operation. "+
				"When true and operating on a partitioned topic name, only deletes the non-partitioned topic "+
				"with the same name, if it exists."),
		),
		mcp.WithBoolean("partitioned",
			mcp.Description("Get stats for a partitioned topic. Optional for 'stats' operation. "+
				"It has to be true if the topic is partitioned. Leave it empty or false for non-partitioned topic."),
		),
		mcp.WithBoolean("per-partition",
			mcp.Description("Include per-partition stats. Optional for 'stats' operation. "+
				"When true, returns statistics for each partition separately. "+
				"Requires 'partitioned' parameter to be true."),
		),
		mcp.WithString("config",
			mcp.Description("JSON configuration for the topic. Required for 'update' operation. "+
				"Set various policies like retention, compaction, deduplication, etc. "+
				"Use a JSON object format, e.g., '{\"deduplicationEnabled\": true, \"replication_clusters\": [\"us-west\", \"us-east\"]}'"),
		),
		mcp.WithString("messageId",
			mcp.Description("Message ID for operations that require a position. Required for 'offload' operation. "+
				"Format is 'ledgerId:entryId' representing a position in the topic's message log. "+
				"For offload operations, specifies the message up to which data should be moved to long-term storage."),
		),
		mcp.WithBoolean("wait",
			mcp.Description("Wait for the long-running status operation to finish. Optional for 'compact-status' (and legacy 'status') and 'offload-status'."),
		),
		mcp.WithString("role",
			mcp.Description("Role name. Required for 'grant-permissions' and 'revoke-permissions' operations."),
		),
		mcp.WithArray("actions",
			mcp.Description("List of topic permissions to grant. Required for 'grant-permissions'. "+
				"Allowed values: produce, consume, sources, sinks, functions, packages."),
			mcp.Items(
				map[string]interface{}{
					"type":        "string",
					"description": "auth action",
				},
			),
		),
		annotation,
	)
}

// buildTopicHandler builds the Pulsar Admin Topic handler function
// Migrated from the original handler logic
func (b *PulsarAdminTopicToolBuilder) buildTopicHandler(mode toolMode) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get required parameters
		resource, err := request.RequireString("resource")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get resource: %v", err)), nil
		}

		operation, err := request.RequireString("operation")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get operation: %v", err)), nil
		}

		// Normalize parameters
		resource = strings.ToLower(resource)
		operation = normalizeTopicOperation(operation)

		if !validateModeOperation(mode, operation, readOnlyRestrictedTopicOperations) {
			return mcp.NewToolResultError(fmt.Sprintf("Operation %q is not available in %s mode", operation, mode)), nil
		}

		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return mcp.NewToolResultError("Pulsar session not found in context"), nil
		}

		// Create the admin client
		admin, err := session.GetAdminClient()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
		}

		// Dispatch based on resource and operation
		switch resource {
		case "topic":
			switch operation {
			case "get":
				return b.handleTopicGet(admin, request)
			case "get-permissions":
				return b.handleTopicGetPermissions(admin, request)
			case "grant-permissions":
				return b.handleTopicGrantPermissions(admin, request)
			case "revoke-permissions":
				return b.handleTopicRevokePermissions(admin, request)
			case "create":
				return b.handleTopicCreate(admin, request)
			case "delete":
				return b.handleTopicDelete(admin, request)
			case "stats":
				return b.handleTopicStats(admin, request)
			case "lookup":
				return b.handleTopicLookup(admin, request)
			case "internal-stats":
				return b.handleTopicInternalStats(admin, request)
			case "internal-info":
				return b.handleTopicInternalInfo(admin, request)
			case "bundle-range":
				return b.handleTopicBundleRange(admin, request)
			case "last-message-id":
				return b.handleTopicLastMessageID(admin, request)
			case "compact-status":
				return b.handleTopicCompactStatus(ctx, admin, request)
			case "unload":
				return b.handleTopicUnload(admin, request)
			case "terminate":
				return b.handleTopicTerminate(admin, request)
			case "compact":
				return b.handleTopicCompact(admin, request)
			case "update":
				return b.handleTopicUpdate(admin, request)
			case "offload":
				return b.handleTopicOffload(admin, request)
			case "offload-status":
				return b.handleTopicOffloadStatus(ctx, admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Unknown topic operation: %s", operation)), nil
			}
		case "topics":
			switch operation {
			case "list":
				return b.handleTopicsList(admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Unknown topics operation: %s", operation)), nil
			}
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Unknown resource: %s", resource)), nil
		}
	}
}

func isReadOnlyRestrictedTopicOperation(operation string) bool {
	_, ok := readOnlyRestrictedTopicOperations[strings.ToLower(operation)]
	return ok
}

func normalizeTopicOperation(operation string) string {
	normalized := strings.ToLower(strings.TrimSpace(operation))
	if alias, ok := topicOperationAliases[normalized]; ok {
		return alias
	}

	return normalized
}

// Unified error handling and utility functions

// handleError provides unified error handling
func (b *PulsarAdminTopicToolBuilder) handleError(operation string, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("Failed to %s: %v", operation, err))
}

// marshalResponse provides unified JSON serialization for responses
func (b *PulsarAdminTopicToolBuilder) marshalResponse(data interface{}) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return b.handleError("marshal response", err), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleTopicsList lists all existing topics under the specified namespace
func (b *PulsarAdminTopicToolBuilder) handleTopicsList(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	namespace, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'namespace' for topics.list: %v", err)), nil
	}

	// Get namespace name
	namespaceName, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name '%s': %v", namespace, err)), nil
	}

	// List topics
	partitionedTopics, nonPartitionedTopics, err := admin.Topics().List(*namespaceName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list topics in namespace '%s': %v",
			namespace, err)), nil
	}

	// Format the output
	result := struct {
		PartitionedTopics    []string `json:"partitionedTopics"`
		NonPartitionedTopics []string `json:"nonPartitionedTopics"`
	}{
		PartitionedTopics:    partitionedTopics,
		NonPartitionedTopics: nonPartitionedTopics,
	}

	return b.marshalResponse(result)
}

// handleTopicGet gets the metadata of an existing topic
func (b *PulsarAdminTopicToolBuilder) handleTopicGet(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := request.RequireString("topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.get: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Get topic metadata
	metadata, err := admin.Topics().GetMetadata(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get metadata for topic '%s': %v",
			topic, err)), nil
	}

	return b.marshalResponse(metadata)
}

// handleTopicGetPermissions gets topic permissions for all roles.
func (b *PulsarAdminTopicToolBuilder) handleTopicGetPermissions(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topic, err := request.RequireString("topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.get-permissions: %v", err)), nil
	}

	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	permissions, err := admin.Topics().GetPermissions(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get permissions for topic '%s': %v", topic, err)), nil
	}

	return b.marshalResponse(permissions)
}

// handleTopicGrantPermissions grants topic permissions to a role.
func (b *PulsarAdminTopicToolBuilder) handleTopicGrantPermissions(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topic, err := request.RequireString("topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.grant-permissions: %v", err)), nil
	}

	role, err := request.RequireString("role")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'role' for topic.grant-permissions: %v", err)), nil
	}
	role = strings.TrimSpace(role)
	if role == "" {
		return mcp.NewToolResultError("Missing required parameter 'role' for topic.grant-permissions"), nil
	}

	actions, err := request.RequireStringSlice("actions")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'actions' for topic.grant-permissions: %v", err)), nil
	}
	if len(actions) == 0 {
		return mcp.NewToolResultError("Missing required parameter 'actions' for topic.grant-permissions"), nil
	}

	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	authActions, err := parseTopicActions(actions)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to parse actions: %v", err)), nil
	}

	if err := admin.Topics().GrantPermission(*topicName, role, authActions); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to grant permissions for topic '%s': %v", topic, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Granted %v permission(s) to role %s on topic %s", actions, role, topicName.String())), nil
}

// handleTopicRevokePermissions revokes topic permissions from a role.
func (b *PulsarAdminTopicToolBuilder) handleTopicRevokePermissions(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topic, err := request.RequireString("topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.revoke-permissions: %v", err)), nil
	}

	role, err := request.RequireString("role")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'role' for topic.revoke-permissions: %v", err)), nil
	}
	role = strings.TrimSpace(role)
	if role == "" {
		return mcp.NewToolResultError("Missing required parameter 'role' for topic.revoke-permissions"), nil
	}

	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	if err := admin.Topics().RevokePermission(*topicName, role); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to revoke permissions for topic '%s': %v", topic, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Revoked all permissions from role %s on topic %s", role, topicName.String())), nil
}

// handleTopicStats gets the stats for an existing topic
func (b *PulsarAdminTopicToolBuilder) handleTopicStats(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := request.RequireString("topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.stats: %v", err)), nil
	}

	// Get optional parameters
	partitioned := request.GetBool("partitioned", false)
	perPartition := request.GetBool("per-partition", false)

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	namespaceName, err := utils.GetNamespaceName(topicName.GetTenant() + "/" + topicName.GetNamespace())
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// List topics to determine if this topic is partitioned
	partitionedTopics, nonPartitionedTopics, err := admin.Topics().List(*namespaceName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list topics in namespace '%s': %v",
			namespaceName, err)), nil
	}

	if slices.Contains(partitionedTopics, topicName.String()) {
		partitioned = true
	}
	if slices.Contains(nonPartitionedTopics, topicName.String()) {
		partitioned = false
	}

	var data interface{}
	if partitioned {
		// Get partitioned topic stats
		stats, err := admin.Topics().GetPartitionedStats(*topicName, perPartition)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get stats for partitioned topic '%s': %v",
				topic, err)), nil
		}
		data = stats
	} else {
		// Get topic stats
		stats, err := admin.Topics().GetStats(*topicName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get stats for topic '%s': %v",
				topic, err)), nil
		}
		data = stats
	}

	return b.marshalResponse(data)
}

// handleTopicLookup looks up the owner broker of a topic
func (b *PulsarAdminTopicToolBuilder) handleTopicLookup(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := request.RequireString("topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.lookup: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Lookup topic
	lookup, err := admin.Topics().Lookup(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to lookup topic '%s': %v",
			topic, err)), nil
	}

	return b.marshalResponse(lookup)
}

// handleTopicCreate creates a topic with the specified number of partitions
func (b *PulsarAdminTopicToolBuilder) handleTopicCreate(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := request.RequireString("topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.create: %v", err)), nil
	}

	partitions, err := request.RequireFloat("partitions")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'partitions' for topic.create: %v", err)), nil
	}

	// Validate partitions
	if partitions < 0 {
		return mcp.NewToolResultError("Invalid partitions number: must be non-negative. Use 0 for a non-partitioned topic."), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Create topic
	err = admin.Topics().Create(*topicName, int(partitions))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create topic '%s' with %d partitions: %v",
			topic, int(partitions), err)), nil
	}

	if int(partitions) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Successfully created non-partitioned topic '%s'",
			topicName.String())), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("Successfully created topic '%s' with %d partitions",
		topicName.String(), int(partitions))), nil
}

// handleTopicDelete deletes a topic
func (b *PulsarAdminTopicToolBuilder) handleTopicDelete(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := request.RequireString("topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.delete: %v", err)), nil
	}

	// Get optional parameters
	force := request.GetBool("force", false)
	nonPartitioned := request.GetBool("non-partitioned", false)

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Delete topic
	err = admin.Topics().Delete(*topicName, force, nonPartitioned)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete topic '%s': %v", topic, err)), nil
	}

	forceStr := ""
	if force {
		forceStr = " forcefully"
	}

	nonPartitionedStr := ""
	if nonPartitioned {
		nonPartitionedStr = " (non-partitioned)"
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully deleted topic '%s'%s%s",
		topicName.String(), forceStr, nonPartitionedStr)), nil
}

// handleTopicUnload unloads a topic
func (b *PulsarAdminTopicToolBuilder) handleTopicUnload(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := request.RequireString("topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.unload: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Unload topic
	err = admin.Topics().Unload(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to unload topic '%s': %v", topic, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully unloaded topic '%s'", topicName.String())), nil
}

// handleTopicTerminate terminates a topic
func (b *PulsarAdminTopicToolBuilder) handleTopicTerminate(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := request.RequireString("topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.terminate: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Terminate topic
	messageID, err := admin.Topics().Terminate(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to terminate topic '%s': %v", topic, err)), nil
	}

	// Convert message ID to string
	msgIDStr := fmt.Sprintf("%d:%d", messageID.LedgerID, messageID.EntryID)

	return mcp.NewToolResultText(fmt.Sprintf("Successfully terminated topic '%s' at message %s. "+
		"No more messages can be published to this topic.",
		topicName.String(), msgIDStr)), nil
}

// handleTopicCompact triggers compaction on a topic
func (b *PulsarAdminTopicToolBuilder) handleTopicCompact(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := request.RequireString("topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.compact: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Compact topic
	err = admin.Topics().Compact(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to trigger compaction for topic '%s': %v", topic, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully triggered compaction for topic '%s'. "+
		"Use operation='compact-status' to check compaction status.", topicName.String())), nil
}

// handleTopicInternalStats gets the internal stats for a topic
func (b *PulsarAdminTopicToolBuilder) handleTopicInternalStats(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := request.RequireString("topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.internal-stats: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Get internal stats
	stats, err := admin.Topics().GetInternalStats(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get internal stats for topic '%s': %v", topic, err)), nil
	}

	return b.marshalResponse(stats)
}

// handleTopicInternalInfo gets the internal info for a topic
func (b *PulsarAdminTopicToolBuilder) handleTopicInternalInfo(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := request.RequireString("topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.internal-info: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Get internal info
	info, err := admin.Topics().GetInternalInfo(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get internal info for topic '%s': %v", topic, err)), nil
	}

	return b.marshalResponse(info)
}

// handleTopicBundleRange gets the bundle range of a topic
func (b *PulsarAdminTopicToolBuilder) handleTopicBundleRange(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := request.RequireString("topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.bundle-range: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Get bundle range
	bundle, err := admin.Topics().GetBundleRange(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get bundle range for topic '%s': %v", topic, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Bundle range for topic '%s': %s", topicName.String(), bundle)), nil
}

// handleTopicLastMessageID gets the last message ID of a topic
func (b *PulsarAdminTopicToolBuilder) handleTopicLastMessageID(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := request.RequireString("topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.last-message-id: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Get last message ID
	messageID, err := admin.Topics().GetLastMessageID(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get last message ID for topic '%s': %v", topic, err)), nil
	}

	return b.marshalResponse(messageID)
}

// handleTopicCompactStatus gets the compaction status of a topic.
func (b *PulsarAdminTopicToolBuilder) handleTopicCompactStatus(ctx context.Context, admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := request.RequireString("topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.compact-status: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	if !topicName.IsPersistent() {
		return mcp.NewToolResultError("Need to provide a persistent topic"), nil
	}

	wait := request.GetBool("wait", false)

	status, err := admin.Topics().CompactStatus(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get compaction status for topic '%s': %v", topic, err)), nil
	}

	err = waitForTopicLongRunningStatus(ctx, wait, topicStatusPollInterval, func() bool {
		return status.Status == utils.RUNNING
	}, func() error {
		status, err = admin.Topics().CompactStatus(*topicName)
		return err
	})
	if err != nil {
		if isTopicStatusWaitInterrupted(err) {
			return mcp.NewToolResultError(fmt.Sprintf(
				"Waiting for compaction status for topic '%s' was interrupted: %v", topic, err,
			)), nil
		}
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get compaction status for topic '%s': %v", topic, err)), nil
	}

	return b.marshalResponse(status)
}

// handleTopicUpdate updates a topic configuration
func (b *PulsarAdminTopicToolBuilder) handleTopicUpdate(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := request.RequireString("topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.update: %v", err)), nil
	}

	partitions, err := request.RequireFloat("partitions")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'partitions' for topic.update: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	err = admin.Topics().Update(*topicName, int(partitions))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to update topic '%s': %v", topic, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully updated topic '%s' partitions to %d",
		topicName.String(), int(partitions))), nil
}

// handleTopicOffload offloads data from a topic to long-term storage
func (b *PulsarAdminTopicToolBuilder) handleTopicOffload(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := request.RequireString("topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.offload: %v", err)), nil
	}

	messageIDStr, err := request.RequireString("messageId")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'messageId' for topic.offload: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Parse message ID from format "ledgerId:entryId"
	var ledgerID, entryID int64
	if _, err := fmt.Sscanf(messageIDStr, "%d:%d", &ledgerID, &entryID); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid message ID format (expected 'ledgerId:entryId'): %v. "+
			"Valid examples: '123:456'", err)), nil
	}

	// Create MessageID object
	messageID := utils.MessageID{
		LedgerID: ledgerID,
		EntryID:  entryID,
	}

	// Offload topic
	err = admin.Topics().Offload(*topicName, messageID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to trigger offload for topic '%s': %v", topic, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully triggered offload for topic '%s' up to message %s. "+
		"Use 'topic.offload-status' to check the offload progress.",
		topicName.String(), messageIDStr)), nil
}

// handleTopicOffloadStatus checks the status of data offloading for a topic
func (b *PulsarAdminTopicToolBuilder) handleTopicOffloadStatus(ctx context.Context, admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := request.RequireString("topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.offload-status: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Get offload status
	status, err := admin.Topics().OffloadStatus(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get offload status for topic '%s': %v", topic, err)), nil
	}

	err = waitForTopicLongRunningStatus(ctx, request.GetBool("wait", false), topicStatusPollInterval, func() bool {
		return status.Status == utils.RUNNING
	}, func() error {
		status, err = admin.Topics().OffloadStatus(*topicName)
		return err
	})
	if err != nil {
		if isTopicStatusWaitInterrupted(err) {
			return mcp.NewToolResultError(fmt.Sprintf(
				"Waiting for offload status for topic '%s' was interrupted: %v", topic, err,
			)), nil
		}
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get offload status for topic '%s': %v", topic, err)), nil
	}

	return b.marshalResponse(status)
}

func waitForTopicLongRunningStatus(
	ctx context.Context,
	wait bool,
	pollInterval time.Duration,
	isRunning func() bool,
	refresh func() error,
) error {
	if !wait || !isRunning() {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if pollInterval <= 0 {
		pollInterval = topicStatusPollInterval
	}

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for isRunning() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := refresh(); err != nil {
				return err
			}
		}
	}

	return nil
}

func isTopicStatusWaitInterrupted(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

func parseTopicActions(actions []string) ([]utils.AuthAction, error) {
	parsed := make([]utils.AuthAction, 0, len(actions))
	for _, action := range actions {
		authAction, err := utils.ParseAuthAction(action)
		if err != nil {
			return nil, err
		}
		parsed = append(parsed, authAction)
	}

	return parsed, nil
}
