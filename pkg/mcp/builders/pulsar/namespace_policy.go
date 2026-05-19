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
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	pulsarctlutils "github.com/streamnative/pulsarctl/pkg/ctl/utils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/toolannotations"
)

var supportedNamespaceSetPolicies = []string{
	"message-ttl",
	"retention",
	"permission",
	"replication-clusters",
	"backlog-quota",
	"topic-auto-creation",
	"schema-validation",
	"schema-auto-update",
	"auto-update-schema",
	"offload-threshold",
	"offload-deletion-lag",
	"compaction-threshold",
	"max-producers-per-topic",
	"max-consumers-per-topic",
	"max-consumers-per-subscription",
	"anti-affinity-group",
	"persistence",
	"deduplication",
	"encryption-required",
	"subscription-auth-mode",
	"subscription-permission",
	"dispatch-rate",
	"replicator-dispatch-rate",
	"subscribe-rate",
	"subscription-dispatch-rate",
	"publish-rate",
}

var supportedNamespaceRemovePolicies = []string{
	"backlog-quota",
	"topic-auto-creation",
	"offload-deletion-lag",
	"anti-affinity-group",
	"permission",
	"subscription-permission",
}

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

	antiAffinityNamespacesTool := b.buildNamespaceGetAntiAffinityNamespacesTool()
	antiAffinityNamespacesHandler := b.buildNamespaceGetAntiAffinityNamespacesHandler()
	tools = append(tools, server.ServerTool{
		Tool:    antiAffinityNamespacesTool,
		Handler: antiAffinityNamespacesHandler,
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
		toolannotations.ReadOnly("Get Pulsar Namespace Policies"),
	)
}

// buildNamespaceSetPolicyTool builds the set policy tool
func (b *PulsarAdminNamespacePolicyToolBuilder) buildNamespaceSetPolicyTool() mcp.Tool {
	toolDesc := "Set a namespace policy in Pulsar. Supported policy types: " +
		strings.Join(supportedNamespaceSetPolicies, ", ") + ". " +
		"Provide the parameters that match the selected policy type. " +
		"For example: ttl for message-ttl; time/size for retention; role/actions for permission; " +
		"clusters for replication-clusters; limit-size/backlog-policy for backlog-quota; enabled for boolean policies; " +
		"topic-type/partitions for topic-auto-creation; compatibility for schema-auto-update; " +
		"lag for offload-deletion-lag; count for max-* policies; group for anti-affinity-group; " +
		"subscription/roles for subscription-permission; persistence and rate fields for persistence and rate policies."

	return mcp.NewTool("pulsar_admin_namespace_policy_set",
		mcp.WithDescription(toolDesc),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("The namespace name (tenant/namespace) to set the policy for"),
		),
		mcp.WithString("policy", mcp.Required(),
			mcp.Description("Type of policy to set. Available options: "+strings.Join(supportedNamespaceSetPolicies, ", ")),
		),
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
			mcp.Description("Retention time value (for example 10m, 2h, 1d, 0, -1)"),
		),
		mcp.WithString("size",
			mcp.Description("Size value used by retention, offload-threshold, or compaction-threshold (for example 10M, 16G, 3T, 0, -1)"),
		),
		mcp.WithString("limit-size",
			mcp.Description("Size limit for backlog quota (e.g., 10M, 16G)"),
		),
		mcp.WithString("limit-time",
			mcp.Description("Time limit in seconds for backlog quota. Default is -1 (infinite)"),
		),
		mcp.WithString("backlog-policy",
			mcp.Description("Retention policy for backlog quota (producer_request_hold, producer_exception, consumer_backlog_eviction)"),
		),
		mcp.WithBoolean("enabled",
			mcp.Description("Boolean value used by schema-validation, auto-update-schema, deduplication, encryption-required, and topic-auto-creation"),
		),
		mcp.WithString("compatibility",
			mcp.Description("Schema auto-update compatibility strategy (AutoUpdateDisabled, Backward, Forward, Full, AlwaysCompatible, BackwardTransitive, ForwardTransitive, FullTransitive)"),
		),
		mcp.WithString("lag",
			mcp.Description("Duration used by offload-deletion-lag (for example 1m, 1h, 24h)"),
		),
		mcp.WithNumber("count",
			mcp.Description("Integer count used by max-producers-per-topic, max-consumers-per-topic, and max-consumers-per-subscription"),
		),
		mcp.WithString("group",
			mcp.Description("Anti-affinity group name used by anti-affinity-group"),
		),
		mcp.WithString("mode",
			mcp.Description("Subscription auth mode used by subscription-auth-mode (for example None or Prefix)"),
		),
		mcp.WithString("subscription",
			mcp.Description("Subscription name used by subscription-permission"),
		),
		mcp.WithString("topic-type",
			mcp.Description("Topic type used by topic-auto-creation (partitioned or non-partitioned)"),
		),
		mcp.WithNumber("partitions",
			mcp.Description("Default partition count used by topic-auto-creation"),
		),
		mcp.WithNumber("ensemble-size",
			mcp.Description("BookKeeper ensemble size used by persistence"),
		),
		mcp.WithNumber("write-quorum-size",
			mcp.Description("BookKeeper write quorum size used by persistence"),
		),
		mcp.WithNumber("ack-quorum-size",
			mcp.Description("BookKeeper ack quorum size used by persistence"),
		),
		mcp.WithNumber("ml-mark-delete-max-rate",
			mcp.Description("Managed ledger mark delete max rate used by persistence"),
		),
		mcp.WithNumber("msg-rate",
			mcp.Description("Message rate used by dispatch-rate, replicator-dispatch-rate, subscription-dispatch-rate, or publish-rate"),
		),
		mcp.WithNumber("byte-rate",
			mcp.Description("Byte rate used by dispatch-rate, replicator-dispatch-rate, subscription-dispatch-rate, or publish-rate"),
		),
		mcp.WithNumber("period",
			mcp.Description("Period in seconds used by dispatch-rate, replicator-dispatch-rate, subscription-dispatch-rate, or subscribe-rate"),
		),
		mcp.WithNumber("subscribe-rate",
			mcp.Description("Subscribe rate per consumer used by subscribe-rate"),
		),
		toolannotations.Destructive("Set Pulsar Namespace Policies"),
	)
}

// buildNamespaceRemovePolicyTool builds the remove policy tool
func (b *PulsarAdminNamespacePolicyToolBuilder) buildNamespaceRemovePolicyTool() mcp.Tool {
	toolDesc := "Remove a namespace policy in Pulsar. Supported policy types: " +
		strings.Join(supportedNamespaceRemovePolicies, ", ") + "."

	return mcp.NewTool("pulsar_admin_namespace_policy_remove",
		mcp.WithDescription(toolDesc),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("The namespace name (tenant/namespace) to remove the policy from"),
		),
		mcp.WithString("policy", mcp.Required(),
			mcp.Description("Type of policy to remove. Available options: "+strings.Join(supportedNamespaceRemovePolicies, ", ")),
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
		toolannotations.Destructive("Remove Pulsar Namespace Policies"),
	)
}

// buildNamespaceGetAntiAffinityNamespacesTool builds the anti-affinity namespace lookup tool.
func (b *PulsarAdminNamespacePolicyToolBuilder) buildNamespaceGetAntiAffinityNamespacesTool() mcp.Tool {
	return mcp.NewTool("pulsar_admin_namespace_policy_get_anti_affinity_namespaces",
		mcp.WithDescription("Get the list of namespaces that belong to the same anti-affinity group. "+
			"This matches `pulsarctl namespaces get-anti-affinity-namespaces` and requires tenant admin permissions."),
		mcp.WithString("group", mcp.Required(),
			mcp.Description("Anti-affinity group name to query."),
		),
		mcp.WithString("cluster",
			mcp.Description("Cluster name to scope the lookup. Optional when the broker can infer it."),
		),
		mcp.WithString("tenant",
			mcp.Description("Tenant name used for authorization. Optional, but recommended when the caller administers multiple tenants."),
		),
		toolannotations.ReadOnly("Get Pulsar Namespace Anti-Affinity Namespaces"),
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

// buildNamespaceGetAntiAffinityNamespacesHandler builds the anti-affinity namespace lookup handler.
func (b *PulsarAdminNamespacePolicyToolBuilder) buildNamespaceGetAntiAffinityNamespacesHandler() func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		group, err := request.RequireString("group")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get anti-affinity group: %v", err)), nil
		}
		group = strings.TrimSpace(group)
		if group == "" {
			return mcp.NewToolResultError("Failed to get anti-affinity group: value cannot be empty"), nil
		}

		tenant := strings.TrimSpace(request.GetString("tenant", ""))
		cluster := strings.TrimSpace(request.GetString("cluster", ""))

		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return mcp.NewToolResultError("Pulsar session not found in context"), nil
		}

		client, err := session.GetAdminClient()
		if err != nil {
			return b.handleError("get admin client", err), nil
		}

		namespaces, err := client.Namespaces().GetAntiAffinityNamespaces(tenant, cluster, group)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get anti-affinity namespaces: %v", err)), nil
		}

		return b.marshalResponse(namespaces)
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
		case "topic-auto-creation":
			return b.handleSetTopicAutoCreation(ctx, client, namespace, request)
		case "schema-validation":
			return b.handleSetSchemaValidation(ctx, client, namespace, request)
		case "schema-auto-update":
			return b.handleSetSchemaAutoUpdate(ctx, client, namespace, request)
		case "auto-update-schema":
			return b.handleSetAutoUpdateSchema(ctx, client, namespace, request)
		case "offload-threshold":
			return b.handleSetOffloadThreshold(ctx, client, namespace, request)
		case "offload-deletion-lag":
			return b.handleSetOffloadDeletionLag(ctx, client, namespace, request)
		case "compaction-threshold":
			return b.handleSetCompactionThreshold(ctx, client, namespace, request)
		case "max-producers-per-topic":
			return b.handleSetMaxProducersPerTopic(ctx, client, namespace, request)
		case "max-consumers-per-topic":
			return b.handleSetMaxConsumersPerTopic(ctx, client, namespace, request)
		case "max-consumers-per-subscription":
			return b.handleSetMaxConsumersPerSubscription(ctx, client, namespace, request)
		case "anti-affinity-group":
			return b.handleSetAntiAffinityGroup(ctx, client, namespace, request)
		case "persistence":
			return b.handleSetPersistence(ctx, client, namespace, request)
		case "deduplication":
			return b.handleSetDeduplication(ctx, client, namespace, request)
		case "encryption-required":
			return b.handleSetEncryptionRequired(ctx, client, namespace, request)
		case "subscription-auth-mode":
			return b.handleSetSubscriptionAuthMode(ctx, client, namespace, request)
		case "subscription-permission":
			return b.handleGrantSubscriptionPermission(ctx, client, namespace, request)
		case "dispatch-rate":
			return b.handleSetDispatchRate(ctx, client, namespace, request)
		case "replicator-dispatch-rate":
			return b.handleSetReplicatorDispatchRate(ctx, client, namespace, request)
		case "subscribe-rate":
			return b.handleSetSubscribeRate(ctx, client, namespace, request)
		case "subscription-dispatch-rate":
			return b.handleSetSubscriptionDispatchRate(ctx, client, namespace, request)
		case "publish-rate":
			return b.handleSetPublishRate(ctx, client, namespace, request)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Unsupported policy type: %s. Supported policies: %s",
				policy, strings.Join(supportedNamespaceSetPolicies, ", "))), nil
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
		case "topic-auto-creation":
			return b.handleRemoveTopicAutoCreation(ctx, client, namespace)
		case "offload-deletion-lag":
			return b.handleRemoveOffloadDeletionLag(ctx, client, namespace)
		case "anti-affinity-group":
			return b.handleRemoveAntiAffinityGroup(ctx, client, namespace)
		case "subscription-permission":
			return b.handleRevokeSubscriptionPermission(ctx, client, namespace, request)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Unsupported policy type for removal: %s. Supported policies: %s",
				policy, strings.Join(supportedNamespaceRemovePolicies, ", "))), nil
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
func (b *PulsarAdminNamespacePolicyToolBuilder) handleGrantPermission(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
func (b *PulsarAdminNamespacePolicyToolBuilder) handleRevokePermission(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetReplicationClusters(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetBacklogQuota(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	limitSizeStr, err := request.RequireString("limit-size")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get limit size: %v", err)), nil
	}

	policyStr, err := request.RequireString("backlog-policy")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get backlog policy: %v", err)), nil
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

func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetTopicAutoCreation(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := b.parseNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	config, err := b.buildTopicAutoCreationConfig(request)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build topic auto-creation config: %v", err)), nil
	}

	if err := client.Namespaces().SetTopicAutoCreation(*ns, config); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set topic auto-creation: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set topic auto-creation policy for %s successfully", namespace)), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetSchemaValidation(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := b.parseNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	enabled, err := b.requireBooleanArgument(request, "enabled")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if err := client.Namespaces().SetSchemaValidationEnforced(*ns, enabled); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set schema validation policy: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set schema validation policy for %s to %t", namespace, enabled)), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetSchemaAutoUpdate(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := b.parseNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	strategy := utils.AutoUpdateDisabled
	if compatibility := request.GetString("compatibility", ""); compatibility != "" {
		strategy, err = utils.ParseSchemaAutoUpdateCompatibilityStrategy(compatibility)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid schema auto-update compatibility: %v", err)), nil
		}
	}

	if err := client.Namespaces().SetSchemaAutoUpdateCompatibilityStrategy(*ns, strategy); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set schema auto-update policy: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set schema auto-update policy for %s to %s", namespace, strategy.String())), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetAutoUpdateSchema(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := b.parseNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	enabled, err := b.requireBooleanArgument(request, "enabled")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if err := client.Namespaces().SetIsAllowAutoUpdateSchema(*ns, enabled); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set auto-update-schema policy: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set auto-update-schema policy for %s to %t", namespace, enabled)), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetOffloadThreshold(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := b.parseNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	sizeStr, err := request.RequireString("size")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get size: %v", err)), nil
	}

	size, err := pulsarctlutils.ValidateSizeString(sizeStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid offload threshold size: %v", err)), nil
	}

	if err := client.Namespaces().SetOffloadThreshold(*ns, size); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set offload threshold: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set offload threshold for %s to %s", namespace, sizeStr)), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetOffloadDeletionLag(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := b.parseNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	lagStr, err := request.RequireString("lag")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get lag duration: %v", err)), nil
	}

	lag, err := time.ParseDuration(lagStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid lag duration: %v", err)), nil
	}

	if err := client.Namespaces().SetOffloadDeleteLag(*ns, lag.Nanoseconds()/1e6); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set offload deletion lag: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set offload deletion lag for %s to %s", namespace, lagStr)), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetCompactionThreshold(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := b.parseNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	sizeStr, err := request.RequireString("size")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get size: %v", err)), nil
	}

	size, err := pulsarctlutils.ValidateSizeString(sizeStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid compaction threshold size: %v", err)), nil
	}

	if err := client.Namespaces().SetCompactionThreshold(*ns, size); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set compaction threshold: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set compaction threshold for %s to %s", namespace, sizeStr)), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetMaxProducersPerTopic(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return b.handleNamespaceCountPolicy(
		namespace,
		request,
		func(ns utils.NameSpaceName, count int) error {
			return client.Namespaces().SetMaxProducersPerTopic(ns, count)
		},
		"max producers per topic",
	)
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetMaxConsumersPerTopic(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return b.handleNamespaceCountPolicy(
		namespace,
		request,
		func(ns utils.NameSpaceName, count int) error {
			return client.Namespaces().SetMaxConsumersPerTopic(ns, count)
		},
		"max consumers per topic",
	)
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetMaxConsumersPerSubscription(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return b.handleNamespaceCountPolicy(
		namespace,
		request,
		func(ns utils.NameSpaceName, count int) error {
			return client.Namespaces().SetMaxConsumersPerSubscription(ns, count)
		},
		"max consumers per subscription",
	)
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetAntiAffinityGroup(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	group, err := request.RequireString("group")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get anti-affinity group: %v", err)), nil
	}

	if err := client.Namespaces().SetNamespaceAntiAffinityGroup(namespace, group); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set anti-affinity group: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set anti-affinity group for %s to %s", namespace, group)), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetPersistence(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	policy, err := b.buildPersistencePolicy(request)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build persistence policy: %v", err)), nil
	}

	if err := client.Namespaces().SetPersistence(namespace, policy); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set persistence policy: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set persistence policy for %s successfully", namespace)), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetDeduplication(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	enabled, err := b.requireBooleanArgument(request, "enabled")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if err := client.Namespaces().SetDeduplicationStatus(namespace, enabled); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set deduplication policy: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set deduplication policy for %s to %t", namespace, enabled)), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetEncryptionRequired(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := b.parseNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	enabled, err := b.requireBooleanArgument(request, "enabled")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if err := client.Namespaces().SetEncryptionRequiredStatus(*ns, enabled); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set encryption-required policy: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set encryption-required policy for %s to %t", namespace, enabled)), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetSubscriptionAuthMode(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := b.parseNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	modeStr, err := request.RequireString("mode")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get subscription auth mode: %v", err)), nil
	}

	mode, err := utils.ParseSubscriptionAuthMode(modeStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid subscription auth mode: %v", err)), nil
	}

	if err := client.Namespaces().SetSubscriptionAuthMode(*ns, mode); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set subscription auth mode: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set subscription auth mode for %s to %s", namespace, mode.String())), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleGrantSubscriptionPermission(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := b.parseNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	subscription, err := request.RequireString("subscription")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get subscription: %v", err)), nil
	}

	roles, err := request.RequireStringSlice("roles")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get roles: %v", err)), nil
	}

	if err := client.Namespaces().GrantSubPermission(*ns, subscription, roles); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to grant subscription permission: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Granted subscription permission for %s/%s to %s",
		namespace, subscription, strings.Join(roles, ", "))), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetDispatchRate(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return b.handleNamespaceDispatchPolicy(
		namespace,
		request,
		b.buildDispatchRate,
		func(ns utils.NameSpaceName, rate utils.DispatchRate) error {
			return client.Namespaces().SetDispatchRate(ns, rate)
		},
		"dispatch rate",
	)
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetReplicatorDispatchRate(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return b.handleNamespaceDispatchPolicy(
		namespace,
		request,
		b.buildDispatchRate,
		func(ns utils.NameSpaceName, rate utils.DispatchRate) error {
			return client.Namespaces().SetReplicatorDispatchRate(ns, rate)
		},
		"replicator dispatch rate",
	)
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetSubscriptionDispatchRate(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return b.handleNamespaceDispatchPolicy(
		namespace,
		request,
		b.buildDispatchRate,
		func(ns utils.NameSpaceName, rate utils.DispatchRate) error {
			return client.Namespaces().SetSubscriptionDispatchRate(ns, rate)
		},
		"subscription dispatch rate",
	)
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetSubscribeRate(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := b.parseNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	rate, err := b.buildSubscribeRate(request)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build subscribe rate: %v", err)), nil
	}

	if err := client.Namespaces().SetSubscribeRate(*ns, rate); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set subscribe rate: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set subscribe rate for %s successfully", namespace)), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetPublishRate(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := b.parseNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	rate, err := b.buildPublishRate(request)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build publish rate: %v", err)), nil
	}

	if err := client.Namespaces().SetPublishRate(*ns, rate); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set publish rate: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set publish rate for %s successfully", namespace)), nil
}

// handleRemoveBacklogQuota handles removing backlog quota for a namespace
func (b *PulsarAdminNamespacePolicyToolBuilder) handleRemoveBacklogQuota(_ context.Context, client cmdutils.Client, namespace string, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if err := client.Namespaces().RemoveBacklogQuota(namespace); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to remove backlog quota: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Removed backlog quota for %s successfully", namespace)), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleRemoveTopicAutoCreation(_ context.Context, client cmdutils.Client, namespace string) (*mcp.CallToolResult, error) {
	ns, err := b.parseNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	if err := client.Namespaces().RemoveTopicAutoCreation(*ns); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to remove topic auto-creation policy: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Removed topic auto-creation policy for %s successfully", namespace)), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleRemoveOffloadDeletionLag(_ context.Context, client cmdutils.Client, namespace string) (*mcp.CallToolResult, error) {
	ns, err := b.parseNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	if err := client.Namespaces().ClearOffloadDeleteLag(*ns); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to remove offload deletion lag policy: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Removed offload deletion lag policy for %s successfully", namespace)), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleRemoveAntiAffinityGroup(_ context.Context, client cmdutils.Client, namespace string) (*mcp.CallToolResult, error) {
	if err := client.Namespaces().DeleteNamespaceAntiAffinityGroup(namespace); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to remove anti-affinity group policy: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Removed anti-affinity group policy for %s successfully", namespace)), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleRevokeSubscriptionPermission(_ context.Context, client cmdutils.Client, namespace string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := b.parseNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	subscription, err := request.RequireString("subscription")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get subscription: %v", err)), nil
	}

	role, err := request.RequireString("role")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get role: %v", err)), nil
	}

	if err := client.Namespaces().RevokeSubPermission(*ns, subscription, role); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to revoke subscription permission: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Revoked subscription permission for role %s on %s/%s",
		role, namespace, subscription)), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleNamespaceCountPolicy(
	namespace string,
	request mcp.CallToolRequest,
	setter func(utils.NameSpaceName, int) error,
	label string,
) (*mcp.CallToolResult, error) {
	ns, err := b.parseNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	count, err := b.requireIntegerArgument(request, "count")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if count < 0 {
		return mcp.NewToolResultError(fmt.Sprintf("%s must be non-negative", label)), nil
	}

	if err := setter(*ns, count); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set %s: %v", label, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set %s for %s to %d", label, namespace, count)), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) handleNamespaceDispatchPolicy(
	namespace string,
	request mcp.CallToolRequest,
	builder func(mcp.CallToolRequest) (utils.DispatchRate, error),
	setter func(utils.NameSpaceName, utils.DispatchRate) error,
	label string,
) (*mcp.CallToolResult, error) {
	ns, err := b.parseNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	rate, err := builder(request)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build %s: %v", label, err)), nil
	}

	if err := setter(*ns, rate); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set %s: %v", label, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set %s for %s successfully", label, namespace)), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) buildTopicAutoCreationConfig(request mcp.CallToolRequest) (utils.TopicAutoCreationConfig, error) {
	enabled, err := b.requireBooleanArgument(request, "enabled")
	if err != nil {
		return utils.TopicAutoCreationConfig{}, err
	}

	config := utils.TopicAutoCreationConfig{Allow: enabled}
	if !enabled {
		return config, nil
	}

	topicType, err := request.RequireString("topic-type")
	if err != nil {
		return utils.TopicAutoCreationConfig{}, fmt.Errorf("missing required parameter 'topic-type': %w", err)
	}

	parsedTopicType, err := utils.ParseTopicType(topicType)
	if err != nil {
		return utils.TopicAutoCreationConfig{}, err
	}

	partitions, err := b.getIntegerArgument(request, "partitions", 0)
	if err != nil {
		return utils.TopicAutoCreationConfig{}, err
	}

	config.Type = parsedTopicType
	config.Partitions = &partitions
	return config, nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) buildPersistencePolicy(request mcp.CallToolRequest) (utils.PersistencePolicies, error) {
	ensembleSize, err := b.requireIntegerArgument(request, "ensemble-size")
	if err != nil {
		return utils.PersistencePolicies{}, err
	}
	writeQuorumSize, err := b.requireIntegerArgument(request, "write-quorum-size")
	if err != nil {
		return utils.PersistencePolicies{}, err
	}
	ackQuorumSize, err := b.requireIntegerArgument(request, "ack-quorum-size")
	if err != nil {
		return utils.PersistencePolicies{}, err
	}
	mlMarkDeleteMaxRate, err := request.RequireFloat("ml-mark-delete-max-rate")
	if err != nil {
		return utils.PersistencePolicies{}, fmt.Errorf("missing required parameter 'ml-mark-delete-max-rate': %w", err)
	}

	return utils.NewPersistencePolicies(
		ensembleSize,
		writeQuorumSize,
		ackQuorumSize,
		mlMarkDeleteMaxRate,
	), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) buildDispatchRate(request mcp.CallToolRequest) (utils.DispatchRate, error) {
	msgRate, err := b.getIntegerArgument(request, "msg-rate", -1)
	if err != nil {
		return utils.DispatchRate{}, err
	}
	byteRate, err := b.getInt64Argument(request, "byte-rate", -1)
	if err != nil {
		return utils.DispatchRate{}, err
	}
	period, err := b.getIntegerArgument(request, "period", 1)
	if err != nil {
		return utils.DispatchRate{}, err
	}

	return utils.DispatchRate{
		DispatchThrottlingRateInMsg:  msgRate,
		DispatchThrottlingRateInByte: byteRate,
		RatePeriodInSecond:           period,
	}, nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) buildSubscribeRate(request mcp.CallToolRequest) (utils.SubscribeRate, error) {
	subscribeRate, err := b.getIntegerArgument(request, "subscribe-rate", -1)
	if err != nil {
		return utils.SubscribeRate{}, err
	}
	period, err := b.getIntegerArgument(request, "period", 30)
	if err != nil {
		return utils.SubscribeRate{}, err
	}

	return utils.SubscribeRate{
		SubscribeThrottlingRatePerConsumer: subscribeRate,
		RatePeriodInSecond:                 period,
	}, nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) buildPublishRate(request mcp.CallToolRequest) (utils.PublishRate, error) {
	msgRate, err := b.getIntegerArgument(request, "msg-rate", -1)
	if err != nil {
		return utils.PublishRate{}, err
	}
	byteRate, err := b.getInt64Argument(request, "byte-rate", -1)
	if err != nil {
		return utils.PublishRate{}, err
	}

	return utils.PublishRate{
		PublishThrottlingRateInMsg:  msgRate,
		PublishThrottlingRateInByte: byteRate,
	}, nil
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

func (b *PulsarAdminNamespacePolicyToolBuilder) parseNamespaceName(namespace string) (*utils.NameSpaceName, error) {
	return utils.GetNamespaceName(namespace)
}

func (b *PulsarAdminNamespacePolicyToolBuilder) requireBooleanArgument(request mcp.CallToolRequest, key string) (bool, error) {
	value, err := request.RequireBool(key)
	if err != nil {
		return false, fmt.Errorf("missing required parameter '%s': %w", key, err)
	}

	return value, nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) requireIntegerArgument(request mcp.CallToolRequest, key string) (int, error) {
	value, err := request.RequireFloat(key)
	if err != nil {
		return 0, fmt.Errorf("missing required parameter '%s': %w", key, err)
	}
	return b.floatToInt(key, value)
}

func (b *PulsarAdminNamespacePolicyToolBuilder) getIntegerArgument(request mcp.CallToolRequest, key string, defaultValue int) (int, error) {
	value := request.GetFloat(key, float64(defaultValue))
	return b.floatToInt(key, value)
}

func (b *PulsarAdminNamespacePolicyToolBuilder) getInt64Argument(request mcp.CallToolRequest, key string, defaultValue int64) (int64, error) {
	value := request.GetFloat(key, float64(defaultValue))
	if math.Trunc(value) != value {
		return 0, fmt.Errorf("parameter '%s' must be an integer", key)
	}

	return int64(value), nil
}

func (b *PulsarAdminNamespacePolicyToolBuilder) floatToInt(key string, value float64) (int, error) {
	if math.Trunc(value) != value {
		return 0, fmt.Errorf("parameter '%s' must be an integer", key)
	}
	return int(value), nil
}
