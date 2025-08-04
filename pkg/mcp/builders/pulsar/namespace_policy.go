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

package pulsar

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	pulsarctlutils "github.com/streamnative/pulsarctl/pkg/ctl/utils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

// PulsarAdminNamespacePolicyToolBuilder implements the ToolBuilder interface for Pulsar admin namespace policies
// /nolint:revive
type PulsarAdminNamespacePolicyToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminNamespacePolicyToolBuilder creates a new Pulsar admin namespace policy tool builder instance
func NewPulsarAdminNamespacePolicyToolBuilder() *PulsarAdminNamespacePolicyToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_namespace_policy",
		Version:     "1.0.0",
		Description: "Pulsar admin namespace policy management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "admin", "namespace_policy"},
	}

	features := []string{
		"pulsar-admin-namespace-policy",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminNamespacePolicyToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin namespace policy tool list
func (b *PulsarAdminNamespacePolicyToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tools := []server.ServerTool{}

	// Always add get policies tool
	getTool := b.buildNamespaceGetPoliciesTool()
	getHandler := b.buildNamespaceGetPoliciesHandler()
	tools = append(tools, server.ServerTool{
		Tool:    getTool,
		Handler: getHandler,
	})

	// Add write operations if not in read-only mode
	if !config.ReadOnly {
		// Add set policy tool
		setTool := b.buildNamespaceSetPolicyTool()
		setHandler := b.buildNamespaceSetPolicyHandler()
		tools = append(tools, server.ServerTool{
			Tool:    setTool,
			Handler: setHandler,
		})

		// Add remove policy tool
		removeTool := b.buildNamespaceRemovePolicyTool()
		removeHandler := b.buildNamespaceRemovePolicyHandler()
		tools = append(tools, server.ServerTool{
			Tool:    removeTool,
			Handler: removeHandler,
		})
	}

	return tools, nil
}

// buildNamespaceGetPoliciesTool builds the get policies tool
func (b *PulsarAdminNamespacePolicyToolBuilder) buildNamespaceGetPoliciesTool() mcp.Tool {
	toolDesc := "Get the configuration policies of a namespace. " +
		"Returns a comprehensive view of all policies applied to the namespace. " +
		"The response includes the following fields:" +
		"\n* bundles: Namespace bundle configuration, including boundaries and number of bundles" +
		"\n* persistence: Message persistence configuration" +
		"\n* retention_policies: Message retention policies defining how long messages are retained" +
		"\n* schema_validation_enforced: Whether schema validation is enforced" +
		"\n* deduplicationEnabled: Whether message deduplication is enabled" +
		"\n* deleted: Whether the namespace is marked as deleted" +
		"\n* encryption_required: Whether message encryption is required" +
		"\n* message_ttl_in_seconds: Time-to-live for messages in seconds, after which unacknowledged messages are deleted" +
		"\n* max_producers_per_topic: Maximum number of producers allowed per topic" +
		"\n* max_consumers_per_topic: Maximum number of consumers allowed per topic" +
		"\n* max_consumers_per_subscription: Maximum number of consumers allowed per subscription" +
		"\n* compaction_threshold: Threshold for topic compaction in bytes" +
		"\n* offload_threshold: Threshold for offloading data to tiered storage in bytes" +
		"\n* offload_deletion_lag_ms: Time lag in milliseconds for retaining offloaded data in hot storage" +
		"\n* antiAffinityGroup: Anti-affinity group for distribution across brokers" +
		"\n* replication_clusters: List of clusters to which data is replicated" +
		"\n* latency_stats_sample_rate: Rate at which latency statistics are collected" +
		"\n* backlog_quota_map: Backlog quota settings controlling behavior when quotas are exceeded" +
		"\n* topicDispatchRate: Rate limiting for topic message dispatch" +
		"\n* subscriptionDispatchRate: Rate limiting for subscription message dispatch" +
		"\n* replicatorDispatchRate: Rate limiting for replicator message dispatch" +
		"\n* publishMaxMessageRate: Maximum rate for publishing messages" +
		"\n* clusterSubscribeRate: Rate limiting for subscriptions per cluster" +
		"\n* autoTopicCreationOverride: Settings for automatic topic creation" +
		"\n* schema_auto_update_compatibility_strategy: Strategy for schema auto-updates" +
		"\n* auth_policies: Authorization policies for the namespace" +
		"\n* subscription_auth_mode: Authentication mode for subscriptions" +
		"\n* is_allow_auto_update_schema: Whether automatic schema updates are allowed" +
		"\nRequires tenant admin permissions."

	return mcp.NewTool("pulsar_admin_namespace_policy_get",
		mcp.WithDescription(toolDesc),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("The namespace name (tenant/namespace) to get policies for"),
		),
	)
}

// buildNamespaceSetPolicyTool builds the set policy tool
func (b *PulsarAdminNamespacePolicyToolBuilder) buildNamespaceSetPolicyTool() mcp.Tool {
	toolDesc := "Set a policy for a namespace. " +
		"This is a unified tool for setting different types of policies on a namespace. " +
		"The policy type determines which specific policy will be set, and the required parameters " +
		"vary based on the policy type. " +
		"Requires appropriate admin permissions based on the policy being modified.\n\n" +
		"Available policy types and their required parameters:\n" +
		"1. message-ttl: Sets time to live for messages\n" +
		"   - Required: namespace, ttl (in seconds)\n\n" +
		"2. retention: Sets retention policy for messages\n" +
		"   - Required: namespace\n" +
		"   - Optional: time (retention time in minutes, 0=no retention, -1=infinite), " +
		"size (e.g., 10M, 16G, 3T, 0=no retention, -1=infinite)\n\n" +
		"3. permission: Grants permissions to a role\n" +
		"   - Required: namespace, role, actions (array of permissions: produce, consume)\n\n" +
		"4. replication-clusters: Sets clusters to replicate to\n" +
		"   - Required: namespace, clusters (array of cluster names)\n\n" +
		"5. backlog-quota: Sets backlog quota policy\n" +
		"   - Required: namespace, limit-size (e.g., 10M, 16G), policy (producer_request_hold, producer_exception, consumer_backlog_eviction)\n" +
		"   - Optional: limit-time (seconds, -1=infinite), type (destination_storage, message_age)\n\n" +
		"Additional policy types: topic-auto-creation, schema-validation, schema-auto-update, auto-update-schema, " +
		"offload-threshold, offload-deletion-lag, compaction-threshold, max-producers-per-topic, max-consumers-per-topic, " +
		"max-consumers-per-subscription, anti-affinity-group, persistence, deduplication, encryption-required, " +
		"subscription-auth-mode, subscription-permission, dispatch-rate, replicator-dispatch-rate, subscribe-rate, " +
		"subscription-dispatch-rate, publish-rate"

	return mcp.NewTool("pulsar_admin_namespace_policy_set",
		mcp.WithDescription(toolDesc),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("The namespace name (tenant/namespace) to set the policy for"),
		),
		mcp.WithString("policy", mcp.Required(),
			mcp.Description("Type of policy to set. Available options: "+
				"message-ttl, retention, permission, replication-clusters, backlog-quota, "+
				"topic-auto-creation, schema-validation, schema-auto-update, auto-update-schema, "+
				"offload-threshold, offload-deletion-lag, compaction-threshold, "+
				"max-producers-per-topic, max-consumers-per-topic, max-consumers-per-subscription, "+
				"anti-affinity-group, persistence, deduplication, encryption-required, "+
				"subscription-auth-mode, subscription-permission, dispatch-rate, "+
				"replicator-dispatch-rate, subscribe-rate, subscription-dispatch-rate, publish-rate"),
		),
		// Generic policy parameters - specific ones will be used based on the policy type
		mcp.WithString("role",
			mcp.Description("Role name for permission policies"),
		),
		mcp.WithArray("actions",
			mcp.Description("Actions to grant for permission policies (e.g., produce, consume)"),
			mcp.Items(
				map[string]interface{}{
					"type":        "string",
					"description": "action",
				},
			),
		),
		mcp.WithArray("clusters",
			mcp.Description("List of clusters for replication policies"),
			mcp.Items(
				map[string]interface{}{
					"type":        "string",
					"description": "cluster",
				},
			),
		),
		mcp.WithArray("roles",
			mcp.Description("List of roles for subscription permission policies"),
			mcp.Items(
				map[string]interface{}{
					"type":        "string",
					"description": "role",
				},
			),
		),
		mcp.WithString("ttl",
			mcp.Description("Message TTL in seconds (or 0 to disable TTL)"),
		),
		mcp.WithString("time",
			mcp.Description("Retention time in minutes, or special values: 0 (no retention) or -1 (infinite retention)"),
		),
		mcp.WithString("size",
			mcp.Description("Retention size limit (e.g., 10M, 16G, 3T), or special values: 0 (no retention) or -1 (infinite size retention)"),
		),
		mcp.WithString("limit-size",
			mcp.Description("Size limit for backlog quota (e.g., 10M, 16G)"),
		),
		mcp.WithString("limit-time",
			mcp.Description("Time limit in seconds for backlog quota. Default is -1 (infinite)"),
		),
		mcp.WithString("policy",
			mcp.Description("Retention policy for backlog quota (valid options: producer_request_hold, producer_exception, consumer_backlog_eviction)"),
		),
		// Add more parameters as needed
	)
}

// buildNamespaceRemovePolicyTool builds the remove policy tool
func (b *PulsarAdminNamespacePolicyToolBuilder) buildNamespaceRemovePolicyTool() mcp.Tool {
	toolDesc := "Remove a policy from a namespace. " +
		"This is a unified tool for removing different types of policies from a namespace. " +
		"The policy type determines which specific policy will be removed. " +
		"Requires appropriate admin permissions based on the policy being modified.\n\n" +
		"Available policy types to remove and their required parameters:\n" +
		"1. backlog-quota: Removes the backlog quota policies\n" +
		"   - Required: namespace\n" +
		"   - Optional: type (destination_storage, message_age)\n\n" +
		"2. topic-auto-creation: Removes topic auto-creation config\n" +
		"   - Required: namespace\n\n" +
		"3. offload-deletion-lag: Clears the offload deletion lag configuration\n" +
		"   - Required: namespace\n\n" +
		"4. anti-affinity-group: Removes the namespace from its anti-affinity group\n" +
		"   - Required: namespace\n\n" +
		"5. permission: Revokes all permissions from a role\n" +
		"   - Required: namespace, role\n\n" +
		"6. subscription-permission: Revokes permission from a role to access a subscription\n" +
		"   - Required: namespace, subscription, role"

	return mcp.NewTool("pulsar_admin_namespace_policy_remove",
		mcp.WithDescription(toolDesc),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("The namespace name (tenant/namespace) to remove the policy from"),
		),
		mcp.WithString("policy", mcp.Required(),
			mcp.Description("Type of policy to remove. Available options: "+
				"backlog-quota, topic-auto-creation, offload-deletion-lag, anti-affinity-group, "+
				"permission, subscription-permission"),
		),
		mcp.WithString("role",
			mcp.Description("Role name for permission policies"),
		),
		mcp.WithString("subscription",
			mcp.Description("Subscription name for subscription permission policies"),
		),
		mcp.WithString("type",
			mcp.Description("Type of backlog quota to remove"),
		),
	)
}

// buildNamespaceGetPoliciesHandler builds the get policies handler
func (b *PulsarAdminNamespacePolicyToolBuilder) buildNamespaceGetPoliciesHandler() func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return mcp.NewToolResultError("Pulsar session not found in context"), nil
		}

		client, err := session.GetAdminClient()
		if err != nil {
			return b.handleError("get admin client", err), nil
		}

		namespace, err := request.RequireString("namespace")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
		}

		// Get policies
		policies, err := client.Namespaces().GetPolicies(namespace)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get policies: %v", err)), nil
		}

		return b.marshalResponse(policies)
	}
}

// buildNamespaceSetPolicyHandler builds the set policy handler
func (b *PulsarAdminNamespacePolicyToolBuilder) buildNamespaceSetPolicyHandler() func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return mcp.NewToolResultError("Pulsar session not found in context"), nil
		}

		client, err := session.GetAdminClient()
		if err != nil {
			return b.handleError("get admin client", err), nil
		}

		namespace, err := request.RequireString("namespace")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
		}

		policy, err := request.RequireString("policy")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get policy type: %v", err)), nil
		}

		// Handle different policy types
		switch policy {
		case "message-ttl":
			return b.handleSetMessageTTL(ctx, client, namespace, request)
		case "retention":
			return b.handleSetRetention(ctx, client, namespace, request)
		case "permission":
			return b.handleGrantPermission(ctx, client, namespace, request)
		case "replication-clusters":
			return b.handleSetReplicationClusters(ctx, client, namespace, request)
		case "backlog-quota":
			return b.handleSetBacklogQuota(ctx, client, namespace, request)
		// Add more policy types as needed
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Unsupported policy type: %s", policy)), nil
		}
	}
}

// buildNamespaceRemovePolicyHandler builds the remove policy handler
func (b *PulsarAdminNamespacePolicyToolBuilder) buildNamespaceRemovePolicyHandler() func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return mcp.NewToolResultError("Pulsar session not found in context"), nil
		}

		client, err := session.GetAdminClient()
		if err != nil {
			return b.handleError("get admin client", err), nil
		}

		namespace, err := request.RequireString("namespace")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
		}

		policy, err := request.RequireString("policy")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get policy type: %v", err)), nil
		}

		// Handle different policy types
		switch policy {
		case "permission":
			return b.handleRevokePermission(ctx, client, namespace, request)
		case "backlog-quota":
			return b.handleRemoveBacklogQuota(ctx, client, namespace, request)
		// Add more policy types as needed
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Unsupported policy type for removal: %s", policy)), nil
		}
	}
}

// Utility functions
func (b *PulsarAdminNamespacePolicyToolBuilder) handleError(operation string, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("Failed to %s: %v", operation, err))
}

func (b *PulsarAdminNamespacePolicyToolBuilder) marshalResponse(data interface{}) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return b.handleError("marshal response", err), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// Policy-specific handler functions

// handleSetMessageTTL handles setting message TTL for a namespace
func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetMessageTTL(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ttlStr, err := request.RequireString("ttl")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get TTL: %v", err)), nil
	}

	ttl, err := strconv.ParseInt(ttlStr, 10, 64)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid TTL value, must be an integer: %v", err)), nil
	}

	// Set message TTL
	err = client.Namespaces().SetNamespaceMessageTTL(namespace, int(ttl))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set message TTL: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set message TTL for %s to %d seconds", namespace, ttl)), nil
}

// handleSetRetention handles setting retention for a namespace
func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetRetention(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	timeStr := request.GetString("time", "")
	sizeStr := request.GetString("size", "")

	if timeStr == "" && sizeStr == "" {
		return mcp.NewToolResultError("At least one of 'time' or 'size' must be specified"), nil
	}

	// Parse retention time
	var retentionTimeInMin int
	if timeStr != "" {
		// Parse relative time in seconds from the input string
		retentionTime, err := pulsarctlutils.ParseRelativeTimeInSeconds(timeStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid retention time format: %v", err)), nil
		}

		if retentionTime != -1 {
			// Convert seconds to minutes
			retentionTimeInMin = int(retentionTime.Minutes())
		} else {
			retentionTimeInMin = -1 // Infinite time retention
		}
	} else {
		retentionTimeInMin = -1 // Default to infinite time retention
	}

	// Parse retention size
	var retentionSizeInMB int
	if sizeStr != "" {
		if sizeStr == "-1" {
			retentionSizeInMB = -1 // Infinite size retention
		} else {
			// Parse size string (e.g., "10M", "16G", "3T")
			sizeInBytes, err := pulsarctlutils.ValidateSizeString(sizeStr)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Invalid retention size format: %v", err)), nil
			}

			if sizeInBytes != -1 {
				// Convert bytes to MB
				retentionSizeInMB = int(sizeInBytes / (1024 * 1024))
				if retentionSizeInMB < 1 {
					return mcp.NewToolResultError("Retention size must be at least 1MB"), nil
				}
			} else {
				retentionSizeInMB = -1 // Infinite size retention
			}
		}
	} else {
		retentionSizeInMB = -1 // Default to infinite size retention
	}

	// Create retention policy
	retention := utils.NewRetentionPolicies(retentionTimeInMin, retentionSizeInMB)

	// Set retention
	err := client.Namespaces().SetRetention(namespace, retention)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set retention: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set retention for %s successfully", namespace)), nil
}

// handleGrantPermission handles granting permissions on a namespace
func (b *PulsarAdminNamespacePolicyToolBuilder) handleGrantPermission(ctx context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	role, err := request.RequireString("role")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get role: %v", err)), nil
	}

	actions, err := request.RequireStringSlice("actions")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get actions: %v", err)), nil
	}

	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	a, err := b.parseActions(actions)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to parse actions: %v", err)), nil
	}

	// Grant permissions
	err = client.Namespaces().GrantNamespacePermission(*ns, role, a)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to grant permission: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Granted %v permission(s) to role %s on %s", actions, role, namespace)), nil
}

// handleRevokePermission handles revoking permissions from a namespace
func (b *PulsarAdminNamespacePolicyToolBuilder) handleRevokePermission(ctx context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	role, err := request.RequireString("role")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get role: %v", err)), nil
	}

	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Revoke permissions
	err = client.Namespaces().RevokeNamespacePermission(*ns, role)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to revoke permission: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Revoked all permissions from role %s on %s", role, namespace)), nil
}

// handleSetReplicationClusters handles setting replication clusters for a namespace
func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetReplicationClusters(ctx context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clusters, err := request.RequireStringSlice("clusters")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get clusters: %v", err)), nil
	}

	if len(clusters) == 0 {
		return mcp.NewToolResultError("At least one cluster must be specified"), nil
	}

	// Set replication clusters
	err = client.Namespaces().SetNamespaceReplicationClusters(namespace, clusters)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set replication clusters: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set replication clusters for %s to %s", namespace, strings.Join(clusters, ", "))), nil
}

// handleSetBacklogQuota handles setting backlog quota for a namespace
func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetBacklogQuota(ctx context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	limitSizeStr, err := request.RequireString("limit-size")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get limit size: %v", err)), nil
	}

	policyStr, err := request.RequireString("policy")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get policy: %v", err)), nil
	}

	// Parse backlog size limit
	limitSize, err := pulsarctlutils.ValidateSizeString(limitSizeStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid limit size format: %v", err)), nil
	}

	// Parse backlog quota policy using the ParseRetentionPolicy function
	policy, err := utils.ParseRetentionPolicy(policyStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid backlog quota policy: %s. Valid options: producer_request_hold, producer_exception, consumer_backlog_eviction", policyStr)), nil
	}

	// Get optional time limit
	limitTimeStr := request.GetString("limit-time", "-1")
	limitTime, err := strconv.ParseInt(limitTimeStr, 10, 64)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid limit time: %v", err)), nil
	}

	// Parse quota type (optional, default to destination_storage)
	quotaTypeStr := request.GetString("type", "destination_storage")
	quotaType := utils.DestinationStorage // Default
	if quotaTypeStr != "" {
		parsedType, err := utils.ParseBacklogQuotaType(quotaTypeStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid backlog quota type: %v", err)), nil
		}
		quotaType = parsedType
	}

	// Create and set backlog quota
	backlogQuota := utils.NewBacklogQuota(limitSize, limitTime, policy)
	err = client.Namespaces().SetBacklogQuota(namespace, backlogQuota, quotaType)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set backlog quota: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set backlog quota for %s successfully", namespace)), nil
}

// handleRemoveBacklogQuota handles removing backlog quota for a namespace
func (b *PulsarAdminNamespacePolicyToolBuilder) handleRemoveBacklogQuota(ctx context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Remove backlog quota (API doesn't require quota type for removal)
	err := client.Namespaces().RemoveBacklogQuota(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to remove backlog quota: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Removed backlog quota for %s successfully", namespace)), nil
}

// parseActions parses action strings into AuthAction enums
func (b *PulsarAdminNamespacePolicyToolBuilder) parseActions(actions []string) ([]utils.AuthAction, error) {
	r := make([]utils.AuthAction, 0)
	for _, v := range actions {
		a, err := utils.ParseAuthAction(v)
		if err != nil {
			return nil, err
		}
		r = append(r, a)
	}
	return r, nil
}
