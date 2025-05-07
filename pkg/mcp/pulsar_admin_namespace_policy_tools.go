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

package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	pulsarctlutils "github.com/streamnative/pulsarctl/pkg/ctl/utils"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// PulsarAdminAddNamespacePolicyTools adds namespace policy-related tools to the MCP server
func PulsarAdminAddNamespacePolicyTools(s *server.MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeaturePulsarAdminNamespacePolicy)) && !slices.Contains(features, string(FeatureAll)) && !slices.Contains(features, string(FeatureAllPulsar)) {
		return
	}
	namespaceGetPoliciesTool := mcp.NewTool("pulsar_admin_namespace_policy_get",
		mcp.WithDescription("Get the configuration policies of a namespace. "+
			"Returns a comprehensive view of all policies applied to the namespace. "+
			"The response includes the following fields:"+
			"\n* bundles: Namespace bundle configuration, including boundaries and number of bundles"+
			"\n* persistence: Message persistence configuration"+
			"\n* retention_policies: Message retention policies defining how long messages are retained"+
			"\n* schema_validation_enforced: Whether schema validation is enforced"+
			"\n* deduplicationEnabled: Whether message deduplication is enabled"+
			"\n* deleted: Whether the namespace is marked as deleted"+
			"\n* encryption_required: Whether message encryption is required"+
			"\n* message_ttl_in_seconds: Time-to-live for messages in seconds, after which unacknowledged messages are deleted"+
			"\n* max_producers_per_topic: Maximum number of producers allowed per topic"+
			"\n* max_consumers_per_topic: Maximum number of consumers allowed per topic"+
			"\n* max_consumers_per_subscription: Maximum number of consumers allowed per subscription"+
			"\n* compaction_threshold: Threshold for topic compaction in bytes"+
			"\n* offload_threshold: Threshold for offloading data to tiered storage in bytes"+
			"\n* offload_deletion_lag_ms: Time lag in milliseconds for retaining offloaded data in hot storage"+
			"\n* antiAffinityGroup: Anti-affinity group for distribution across brokers"+
			"\n* replication_clusters: List of clusters to which data is replicated"+
			"\n* latency_stats_sample_rate: Rate at which latency statistics are collected"+
			"\n* backlog_quota_map: Backlog quota settings controlling behavior when quotas are exceeded"+
			"\n* topicDispatchRate: Rate limiting for topic message dispatch"+
			"\n* subscriptionDispatchRate: Rate limiting for subscription message dispatch"+
			"\n* replicatorDispatchRate: Rate limiting for replicator message dispatch"+
			"\n* publishMaxMessageRate: Maximum rate for publishing messages"+
			"\n* clusterSubscribeRate: Rate limiting for subscriptions per cluster"+
			"\n* autoTopicCreationOverride: Settings for automatic topic creation"+
			"\n* schema_auto_update_compatibility_strategy: Strategy for schema auto-updates"+
			"\n* auth_policies: Authorization policies for the namespace"+
			"\n* subscription_auth_mode: Authentication mode for subscriptions"+
			"\n* is_allow_auto_update_schema: Whether automatic schema updates are allowed"+
			"\nRequires tenant admin permissions."),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("The namespace name (tenant/namespace) to get policies for"),
		),
	)
	s.AddTool(namespaceGetPoliciesTool, handleNamespaceGetPolicies)

	// Only add write operations if not in read-only mode
	if !readOnly {
		// Add a unified namespace set policy tool
		namespaceSetPolicyTool := mcp.NewTool("pulsar_admin_namespace_policy_set",
			mcp.WithDescription("Set a policy for a namespace. "+
				"This is a unified tool for setting different types of policies on a namespace. "+
				"The policy type determines which specific policy will be set, and the required parameters "+
				"vary based on the policy type. "+
				"Requires appropriate admin permissions based on the policy being modified.\n\n"+
				"Available policy types and their required parameters:\n"+
				"1. message-ttl: Sets time to live for messages\n"+
				"   - Required: namespace, ttl (in seconds)\n\n"+
				"2. retention: Sets retention policy for messages\n"+
				"   - Required: namespace\n"+
				"   - Optional: time (retention time in minutes, 0=no retention, -1=infinite), "+
				"size (e.g., 10M, 16G, 3T, 0=no retention, -1=infinite)\n\n"+
				"3. permission: Grants permissions to a role\n"+
				"   - Required: namespace, role, actions (array of permissions: produce, consume)\n\n"+
				"4. replication-clusters: Sets clusters to replicate to\n"+
				"   - Required: namespace, clusters (array of cluster names)\n\n"+
				"5. backlog-quota: Sets backlog quota policy\n"+
				"   - Required: namespace, limit-size (e.g., 10M, 16G), policy (producer_request_hold, producer_exception, consumer_backlog_eviction)\n"+
				"   - Optional: limit-time (seconds, -1=infinite), type (destination_storage, message_age)\n\n"+
				"6. topic-auto-creation: Configures automatic topic creation\n"+
				"   - Required: namespace\n"+
				"   - Optional: disable (true/false), type (partitioned, non-partitioned), partitions (integer)\n\n"+
				"7. schema-validation: Sets schema validation enforcement\n"+
				"   - Required: namespace\n"+
				"   - Optional: disable (true/false)\n\n"+
				"8. schema-auto-update: Sets schema auto-update strategy\n"+
				"   - Required: namespace, compatibility (e.g., Backward, Forward, Full, AlwaysCompatible)\n\n"+
				"9. auto-update-schema: Controls if schemas can be automatically updated\n"+
				"   - Required: namespace\n"+
				"   - Optional: enable (true/false), disable (true/false)\n\n"+
				"10. offload-threshold: Sets threshold for data offloading\n"+
				"   - Required: namespace, threshold (e.g., 10M, 16G, 3T, -1 to disable)\n\n"+
				"11. offload-deletion-lag: Sets time to wait before deleting offloaded data\n"+
				"   - Required: namespace, lag (e.g., 1s, 1m, 1h)\n\n"+
				"12. compaction-threshold: Sets threshold for topic compaction\n"+
				"   - Required: namespace, threshold (e.g., 10M, 16G, 3T, 0 to disable)\n\n"+
				"13. max-producers-per-topic: Limits producers per topic\n"+
				"   - Required: namespace, max (integer > 0)\n\n"+
				"14. max-consumers-per-topic: Limits consumers per topic\n"+
				"   - Required: namespace, max (integer > 0)\n\n"+
				"15. max-consumers-per-subscription: Limits consumers per subscription\n"+
				"   - Required: namespace, max (integer > 0)\n\n"+
				"16. anti-affinity-group: Sets anti-affinity group for isolation\n"+
				"   - Required: namespace, group (group name)\n\n"+
				"17. persistence: Sets persistence configuration\n"+
				"   - Required: namespace, ensemble-size, write-quorum-size, ack-quorum-size\n"+
				"   - Optional: ml-mark-delete-max-rate\n\n"+
				"18. deduplication: Controls message deduplication\n"+
				"   - Required: namespace\n"+
				"   - Optional: enable (true/false), disable (true/false)\n\n"+
				"19. encryption-required: Controls message encryption\n"+
				"   - Required: namespace\n"+
				"   - Optional: disable (true/false)\n\n"+
				"20. subscription-auth-mode: Sets subscription auth mode\n"+
				"   - Required: namespace, mode (e.g., None, Prefix)\n\n"+
				"21. subscription-permission: Grants permissions to access a subscription\n"+
				"   - Required: namespace, subscription, roles (array of role names)\n\n"+
				"22. dispatch-rate: Sets message dispatch rate\n"+
				"   - Required: namespace\n"+
				"   - Optional: dispatchThrottlingRateInMsg, dispatchThrottlingRateInByte, ratePeriodInSecond\n\n"+
				"23. replicator-dispatch-rate: Sets replicator dispatch rate\n"+
				"   - Required: namespace\n"+
				"   - Optional: dispatchThrottlingRateInMsg, dispatchThrottlingRateInByte, ratePeriodInSecond\n\n"+
				"24. subscribe-rate: Sets subscribe rate per consumer\n"+
				"   - Required: namespace\n"+
				"   - Optional: subscribeThrottlingRatePerConsumer, ratePeriodInSecond\n\n"+
				"25. subscription-dispatch-rate: Sets subscription dispatch rate\n"+
				"   - Required: namespace\n"+
				"   - Optional: dispatchThrottlingRateInMsg, dispatchThrottlingRateInByte, ratePeriodInSecond\n\n"+
				"26. publish-rate: Sets maximum message publish rate\n"+
				"   - Required: namespace\n"+
				"   - Optional: publishThrottlingRateInMsg, publishThrottlingRateInByte"),
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
			),
			mcp.WithArray("clusters",
				mcp.Description("List of clusters for replication policies"),
			),
			mcp.WithArray("roles",
				mcp.Description("List of roles for subscription permission policies"),
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
			mcp.WithString("type",
				mcp.Description("Type parameter whose meaning depends on the policy type being set"),
			),
			mcp.WithString("disable",
				mcp.Description("Disable flag for various policies (true/false)"),
			),
			mcp.WithString("enable",
				mcp.Description("Enable flag for various policies (true/false)"),
			),
			mcp.WithString("partitions",
				mcp.Description("Number of partitions for topic auto-creation policy"),
			),
			mcp.WithString("compatibility",
				mcp.Description("Compatibility level for schema auto-update (e.g., BACKWARD, FORWARD, FULL)"),
			),
			mcp.WithString("threshold",
				mcp.Description("Threshold value for compaction or offload policies"),
			),
			mcp.WithString("lag",
				mcp.Description("Lag duration for offload deletion (e.g., 1h, 2d)"),
			),
			mcp.WithString("max",
				mcp.Description("Maximum value for producer/consumer limit policies"),
			),
			mcp.WithString("group",
				mcp.Description("Group name for anti-affinity policies"),
			),
			mcp.WithString("ensemble-size",
				mcp.Description("Number of bookies to use for a topic in persistence policy"),
			),
			mcp.WithString("write-quorum-size",
				mcp.Description("How many writes to make of each entry in persistence policy"),
			),
			mcp.WithString("ack-quorum-size",
				mcp.Description("Number of acks to wait for each entry in persistence policy"),
			),
			mcp.WithString("ml-mark-delete-max-rate",
				mcp.Description("Throttling rate of mark-delete operation in persistence policy"),
			),
			mcp.WithString("mode",
				mcp.Description("Mode for subscription auth policy (e.g., None, Prefix)"),
			),
			mcp.WithString("subscription",
				mcp.Description("Subscription name for subscription permission policies"),
			),
			mcp.WithString("dispatchThrottlingRateInMsg",
				mcp.Description("Message dispatch/throttling rate in messages per second"),
			),
			mcp.WithString("dispatchThrottlingRateInByte",
				mcp.Description("Message dispatch/throttling rate in bytes per second"),
			),
			mcp.WithString("ratePeriodInSecond",
				mcp.Description("Rate period in seconds for rate limiting policies"),
			),
			mcp.WithString("subscribeThrottlingRatePerConsumer",
				mcp.Description("Subscribe rate in subscriptions per consumer per second"),
			),
			mcp.WithString("publishThrottlingRateInMsg",
				mcp.Description("Message publish rate in messages per second"),
			),
			mcp.WithString("publishThrottlingRateInByte",
				mcp.Description("Message publish rate in bytes per second"),
			),
		)
		s.AddTool(namespaceSetPolicyTool, handleNamespaceSetPolicy)

		// Add a unified namespace remove policy tool
		namespaceRemovePolicyTool := mcp.NewTool("pulsar_admin_namespace_policy_remove",
			mcp.WithDescription("Remove a policy from a namespace. "+
				"This is a unified tool for removing different types of policies from a namespace. "+
				"The policy type determines which specific policy will be removed. "+
				"Requires appropriate admin permissions based on the policy being modified.\n\n"+
				"Available policy types to remove and their required parameters:\n"+
				"1. backlog-quota: Removes the backlog quota policies\n"+
				"   - Required: namespace\n"+
				"   - Optional: type (destination_storage, message_age)\n\n"+
				"2. topic-auto-creation: Removes topic auto-creation config\n"+
				"   - Required: namespace\n\n"+
				"3. offload-deletion-lag: Clears the offload deletion lag configuration\n"+
				"   - Required: namespace\n\n"+
				"4. anti-affinity-group: Removes the namespace from its anti-affinity group\n"+
				"   - Required: namespace\n\n"+
				"5. permission: Revokes all permissions from a role\n"+
				"   - Required: namespace, role\n\n"+
				"6. subscription-permission: Revokes permission from a role to access a subscription\n"+
				"   - Required: namespace, subscription, role"),
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
		s.AddTool(namespaceRemovePolicyTool, handleNamespaceRemovePolicy)
	}
}

// handleNamespaceGetPolicies handles getting policies for a namespace
func handleNamespaceGetPolicies(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Get policies
	policies, err := client.Namespaces().GetPolicies(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get policies: %v", err)), nil
	}

	// Convert result to JSON string
	policiesJSON, err := json.Marshal(policies)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize policies: %v", err)), nil
	}

	return mcp.NewToolResultText(string(policiesJSON)), nil
}

// handleSetMessageTTL handles setting message TTL for a namespace
func handleSetMessageTTL(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	ttlStr, err := requiredParam[string](request.Params.Arguments, "ttl")
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
func handleSetRetention(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	timeStr, hasTime := optionalParam[string](request.Params.Arguments, "time")
	sizeStr, hasSize := optionalParam[string](request.Params.Arguments, "size")

	if !hasTime && !hasSize {
		return mcp.NewToolResultError("At least one of 'time' or 'size' must be specified"), nil
	}

	// Parse retention time
	var retentionTimeInMin int
	if hasTime {
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
	if hasSize {
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
	err = client.Namespaces().SetRetention(namespace, retention)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set retention: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Set retention for %s successfully. ", namespace),
	), nil
}

// handleGrantPermission handles granting permissions on a namespace
func handleGrantPermission(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	role, err := requiredParam[string](request.Params.Arguments, "role")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get role: %v", err)), nil
	}

	actions, err := requiredParamArray[string](request.Params.Arguments, "actions")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get actions: %v", err)), nil
	}

	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	a, err := parseActions(actions)
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

func parseActions(actions []string) ([]utils.AuthAction, error) {
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

// handleRevokePermission handles revoking permissions from a namespace
func handleRevokePermission(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	role, err := requiredParam[string](request.Params.Arguments, "role")
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
func handleSetReplicationClusters(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	clusters, err := requiredParamArray[string](request.Params.Arguments, "clusters")
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

	return mcp.NewToolResultText(
		fmt.Sprintf("Set replication clusters for %s to %s", namespace, strings.Join(clusters, ", ")),
	), nil
}

// handleSetBacklogQuota handles setting backlog quota for a namespace
func handleSetBacklogQuota(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	limitSizeStr, err := requiredParam[string](request.Params.Arguments, "limit-size")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get limit size: %v", err)), nil
	}

	policyStr, err := requiredParam[string](request.Params.Arguments, "policy")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get policy: %v", err)), nil
	}

	// Parse size limit
	sizeLimit, err := pulsarctlutils.ValidateSizeString(limitSizeStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid size limit format: %v", err)), nil
	}

	// Parse time limit (optional)
	limitTimeStr, hasLimitTime := optionalParam[string](request.Params.Arguments, "limit-time")
	var limitTime int64 = -1 // Default to -1 (infinite)
	if hasLimitTime && limitTimeStr != "" {
		limitTimeVal, err := strconv.ParseInt(limitTimeStr, 10, 64)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid time limit format: %v", err)), nil
		}
		limitTime = limitTimeVal
	}

	// Parse policy
	var policy utils.RetentionPolicy
	switch policyStr {
	case "producer_request_hold":
		policy = utils.ProducerRequestHold
	case "producer_exception":
		policy = utils.ProducerException
	case "consumer_backlog_eviction":
		policy = utils.ConsumerBacklogEviction
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Invalid retention policy type: %v", policyStr)), nil
	}

	// Parse quota type (optional)
	quotaTypeStr, hasQuotaType := optionalParam[string](request.Params.Arguments, "type")
	quotaType := utils.DestinationStorage // Default
	if hasQuotaType && quotaTypeStr != "" {
		parsedType, err := utils.ParseBacklogQuotaType(quotaTypeStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid backlog quota type: %v", err)), nil
		}
		quotaType = parsedType
	}

	// Set backlog quota
	backlogQuota := utils.NewBacklogQuota(sizeLimit, limitTime, policy)
	err = client.Namespaces().SetBacklogQuota(namespace, backlogQuota, quotaType)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set backlog quota: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set backlog quota successfully for [%s]", namespace)), nil
}

// handleRemoveBacklogQuota handles removing backlog quota from a namespace
func handleRemoveBacklogQuota(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Remove backlog quota
	err = client.Namespaces().RemoveBacklogQuota(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to remove backlog quota: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Remove backlog quota successfully for [%s]", namespace)), nil
}

// handleSetTopicAutoCreation handles setting topic auto-creation for a namespace
func handleSetTopicAutoCreation(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Get namespace
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Check if disabled
	disableStr, hasDisable := optionalParam[string](request.Params.Arguments, "disable")
	disable := false
	if hasDisable && disableStr == "true" {
		disable = true
	}

	// Create config
	config := utils.TopicAutoCreationConfig{
		Allow: !disable,
	}

	// Only set topic type and partitions if not disabled
	if !disable {
		// Get topic type (optional)
		topicTypeStr, hasType := optionalParam[string](request.Params.Arguments, "type")
		if hasType && topicTypeStr != "" {
			parsedType, err := utils.ParseTopicType(topicTypeStr)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Invalid topic type: %v", err)), nil
			}
			config.Type = parsedType
		}

		// Get partitions (optional)
		partitionsStr, hasPartitions := optionalParam[string](request.Params.Arguments, "partitions")
		if hasPartitions && partitionsStr != "" {
			partitions, err := strconv.Atoi(partitionsStr)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Invalid partitions value: %v", err)), nil
			}
			config.Partitions = &partitions
		}
	}

	// Set topic auto-creation
	err = client.Namespaces().SetTopicAutoCreation(*ns, config)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set topic auto-creation: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set topic auto-creation config successfully for [%s]", namespace)), nil
}

// handleRemoveTopicAutoCreation handles removing topic auto-creation config from a namespace
func handleRemoveTopicAutoCreation(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Get namespace
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Remove topic auto-creation
	err = client.Namespaces().RemoveTopicAutoCreation(*ns)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to remove topic auto-creation: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Removed topic auto-creation config successfully for [%s]", namespace)), nil
}

// handleSetSchemaValidationEnforced handles setting schema validation enforced status for a namespace
func handleSetSchemaValidationEnforced(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Get namespace name
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Check if disabled
	disableStr, hasDisable := optionalParam[string](request.Params.Arguments, "disable")
	disable := false
	if hasDisable && disableStr == "true" {
		disable = true
	}

	// Set schema validation enforced
	err = client.Namespaces().SetSchemaValidationEnforced(*ns, !disable)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set schema validation enforced: %v", err)), nil
	}

	action := "enabled"
	if disable {
		action = "disabled"
	}
	return mcp.NewToolResultText(fmt.Sprintf("Schema validation enforced is %s for [%s]", action, namespace)), nil
}

// handleSetSchemaAutoUpdateStrategy handles setting schema auto-update strategy for a namespace
func handleSetSchemaAutoUpdateStrategy(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	strategyStr, err := requiredParam[string](request.Params.Arguments, "compatibility")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get compatibility strategy: %v", err)), nil
	}

	// Get namespace name
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Parse strategy
	strategy, err := utils.ParseSchemaAutoUpdateCompatibilityStrategy(strategyStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid schema auto-update strategy: %v", err)), nil
	}

	// Set schema auto-update strategy
	err = client.Namespaces().SetSchemaAutoUpdateCompatibilityStrategy(*ns, strategy)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set schema auto-update strategy: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Set schema auto-update strategy to %s for [%s]", strategy.String(), namespace),
	), nil
}

// handleSetIsAllowAutoUpdateSchema handles setting auto update schema allowance for a namespace
func handleSetIsAllowAutoUpdateSchema(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Get namespace name
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Check if enabled or disabled
	enableStr, hasEnable := optionalParam[string](request.Params.Arguments, "enable")
	disableStr, hasDisable := optionalParam[string](request.Params.Arguments, "disable")

	if (hasEnable && enableStr == "true") && (hasDisable && disableStr == "true") {
		return mcp.NewToolResultError("Specify only one of 'enable' or 'disable'"), nil
	}

	var isAllowUpdateSchema bool
	//nolint:gocritic
	if hasEnable && enableStr == "true" {
		isAllowUpdateSchema = true
	} else if hasDisable && disableStr == "true" {
		isAllowUpdateSchema = false
	} else if !hasEnable && !hasDisable {
		return mcp.NewToolResultError("You must specify either 'enable=true' or 'disable=true'"), nil
	}

	// Set is allow auto update schema
	err = client.Namespaces().SetIsAllowAutoUpdateSchema(*ns, isAllowUpdateSchema)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set auto update schema allowance: %v", err)), nil
	}

	action := "enabled"
	if !isAllowUpdateSchema {
		action = "disabled"
	}
	return mcp.NewToolResultText(fmt.Sprintf("Auto update schema is %s for [%s]", action, namespace)), nil
}

// handleSetOffloadThreshold handles setting offload threshold for a namespace
func handleSetOffloadThreshold(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	thresholdStr, err := requiredParam[string](request.Params.Arguments, "threshold")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get threshold: %v", err)), nil
	}

	// Get namespace name
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Parse size threshold
	threshold, err := pulsarctlutils.ValidateSizeString(thresholdStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid threshold format: %v", err)), nil
	}

	// Set offload threshold
	err = client.Namespaces().SetOffloadThreshold(*ns, threshold)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set offload threshold: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Successfully set offload threshold for %s to %s", namespace, thresholdStr),
	), nil
}

// handleSetOffloadDeletionLag handles setting offload deletion lag for a namespace
func handleSetOffloadDeletionLag(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	lagStr, err := requiredParam[string](request.Params.Arguments, "lag")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get lag: %v", err)), nil
	}

	// Get namespace name
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Parse duration string
	duration, err := time.ParseDuration(lagStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid duration format: %v", err)), nil
	}

	// Convert duration to milliseconds
	lagMs := duration.Milliseconds()

	// Set offload deletion lag
	err = client.Namespaces().SetOffloadDeleteLag(*ns, lagMs)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set offload deletion lag: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Successfully set offload deletion lag for %s to %s", namespace, lagStr),
	), nil
}

// handleClearOffloadDeletionLag handles clearing offload deletion lag for a namespace
func handleClearOffloadDeletionLag(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Get namespace name
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Clear offload deletion lag
	err = client.Namespaces().ClearOffloadDeleteLag(*ns)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to clear offload deletion lag: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Successfully cleared offload deletion lag for %s", namespace),
	), nil
}

// handleSetCompactionThreshold handles setting compaction threshold for a namespace
func handleSetCompactionThreshold(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	thresholdStr, err := requiredParam[string](request.Params.Arguments, "threshold")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get threshold: %v", err)), nil
	}

	// Get namespace name
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Parse size threshold
	threshold, err := pulsarctlutils.ValidateSizeString(thresholdStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid threshold format: %v", err)), nil
	}

	// Set compaction threshold
	err = client.Namespaces().SetCompactionThreshold(*ns, threshold)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set compaction threshold: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Successfully set compaction threshold for %s to %s", namespace, thresholdStr),
	), nil
}

// handleSetMaxProducersPerTopic handles setting max producers per topic for a namespace
func handleSetMaxProducersPerTopic(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	maxStr, err := requiredParam[string](request.Params.Arguments, "max")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get max value: %v", err)), nil
	}

	// Parse maxVal value
	maxVal, err := strconv.Atoi(maxStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid max value, must be an integer: %v", err)), nil
	}

	if maxVal <= 0 {
		return mcp.NewToolResultError("Max producers per topic must be greater than 0"), nil
	}

	// Get namespace name
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Set max producers per topic
	err = client.Namespaces().SetMaxProducersPerTopic(*ns, maxVal)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set max producers per topic: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Successfully set max producers per topic for %s to %d", namespace, maxVal),
	), nil
}

// handleSetMaxConsumersPerTopic handles setting max consumers per topic for a namespace
func handleSetMaxConsumersPerTopic(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	maxStr, err := requiredParam[string](request.Params.Arguments, "max")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get max value: %v", err)), nil
	}

	// Parse maxVal value
	maxVal, err := strconv.Atoi(maxStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid max value, must be an integer: %v", err)), nil
	}

	if maxVal <= 0 {
		return mcp.NewToolResultError("Max consumers per topic must be greater than 0"), nil
	}

	// Get namespace name
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Set max consumers per topic
	err = client.Namespaces().SetMaxConsumersPerTopic(*ns, maxVal)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set max consumers per topic: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Successfully set max consumers per topic for %s to %d", namespace, maxVal),
	), nil
}

// handleSetMaxConsumersPerSubscription handles setting max consumers per subscription for a namespace
func handleSetMaxConsumersPerSubscription(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	maxStr, err := requiredParam[string](request.Params.Arguments, "max")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get max value: %v", err)), nil
	}

	// Parse maxVal value
	maxVal, err := strconv.Atoi(maxStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid max value, must be an integer: %v", err)), nil
	}

	if maxVal <= 0 {
		return mcp.NewToolResultError("Max consumers per subscription must be greater than 0"), nil
	}

	// Get namespace name
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Set max consumers per subscription
	err = client.Namespaces().SetMaxConsumersPerSubscription(*ns, maxVal)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set max consumers per subscription: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Successfully set max consumers per subscription for %s to %d", namespace, maxVal),
	), nil
}

// handleSetAntiAffinityGroup handles setting the anti-affinity group for a namespace
func handleSetAntiAffinityGroup(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	group, err := requiredParam[string](request.Params.Arguments, "group")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get anti-affinity group: %v", err)), nil
	}

	// Set anti-affinity group
	err = client.Namespaces().SetNamespaceAntiAffinityGroup(namespace, group)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set anti-affinity group: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Successfully set anti-affinity group '%s' for namespace %s", group, namespace),
	), nil
}

// handleDeleteAntiAffinityGroup handles deleting the anti-affinity group of a namespace
func handleDeleteAntiAffinityGroup(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Delete anti-affinity group
	err = client.Namespaces().DeleteNamespaceAntiAffinityGroup(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete anti-affinity group: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Successfully deleted anti-affinity group for namespace %s", namespace),
	), nil
}

// handleSetPersistence handles setting persistence policy for a namespace
func handleSetPersistence(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	ensembleSizeStr, err := requiredParam[string](request.Params.Arguments, "ensemble-size")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get ensemble size: %v", err)), nil
	}

	writeQuorumSizeStr, err := requiredParam[string](request.Params.Arguments, "write-quorum-size")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get write quorum size: %v", err)), nil
	}

	ackQuorumSizeStr, err := requiredParam[string](request.Params.Arguments, "ack-quorum-size")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get ack quorum size: %v", err)), nil
	}

	// Parse integer parameters
	ensembleSize, err := strconv.Atoi(ensembleSizeStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid ensemble size, must be an integer: %v", err)), nil
	}

	writeQuorumSize, err := strconv.Atoi(writeQuorumSizeStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid write quorum size, must be an integer: %v", err)), nil
	}

	ackQuorumSize, err := strconv.Atoi(ackQuorumSizeStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid ack quorum size, must be an integer: %v", err)), nil
	}

	// Check Bookkeeper constraints
	if ensembleSize < writeQuorumSize || writeQuorumSize < ackQuorumSize {
		return mcp.NewToolResultError("Bookkeeper ensemble size must be >= write quorum size >= ack quorum size"), nil
	}

	// Parse optional rate parameter
	markDeleteMaxRate := 0.0
	markDeleteMaxRateStr, hasRate := optionalParam[string](request.Params.Arguments, "ml-mark-delete-max-rate")
	if hasRate && markDeleteMaxRateStr != "" {
		rate, err := strconv.ParseFloat(markDeleteMaxRateStr, 64)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid mark delete max rate, must be a number: %v", err)), nil
		}
		markDeleteMaxRate = rate
	}

	// Create persistence policy
	persistencePolicies := utils.NewPersistencePolicies(ensembleSize, writeQuorumSize, ackQuorumSize, markDeleteMaxRate)

	// Set persistence policy
	err = client.Namespaces().SetPersistence(namespace, persistencePolicies)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set persistence policy: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Successfully set persistence policy for %s", namespace),
	), nil
}

// handleSetDeduplication handles setting the deduplication status for a namespace
func handleSetDeduplication(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Get enable flag, default to false if not specified
	enable, _ := optionalParam[string](request.Params.Arguments, "enable")
	enableBool := false
	if enable == "true" {
		enableBool = true
	}

	// Set deduplication status
	err = client.Namespaces().SetDeduplicationStatus(namespace, enableBool)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set deduplication status: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Successfully set deduplication status to [%v] for %s", enableBool, namespace),
	), nil
}

// handleSetEncryptionRequired handles setting whether encryption is required for a namespace
func handleSetEncryptionRequired(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Get disable flag, default to false if not specified
	disable, _ := optionalParam[string](request.Params.Arguments, "disable")
	disableFlag := disable == "true"

	// Get namespace name
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Set encryption required status (note: API expects !disable)
	err = client.Namespaces().SetEncryptionRequiredStatus(*ns, !disableFlag)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set encryption required status: %v", err)), nil
	}

	status := "Enable"
	if disableFlag {
		status = "Disable"
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("%s messages encryption for namespace %s", status, namespace),
	), nil
}

// handleSetSubscriptionAuthMode handles setting the default subscription auth mode for a namespace
func handleSetSubscriptionAuthMode(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	mode, err := requiredParam[string](request.Params.Arguments, "mode")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get subscription auth mode: %v", err)), nil
	}

	// Get namespace name
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Parse subscription auth mode
	authMode, err := utils.ParseSubscriptionAuthMode(mode)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid subscription auth mode: %v", err)), nil
	}

	// Set subscription auth mode
	err = client.Namespaces().SetSubscriptionAuthMode(*ns, authMode)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set subscription auth mode: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Successfully set the default subscription auth mode of namespace %s to %s",
			namespace, authMode.String()),
	), nil
}

// handleGrantSubscriptionPermission handles granting subscription permissions to roles
func handleGrantSubscriptionPermission(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	subscription, err := requiredParam[string](request.Params.Arguments, "subscription")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get subscription name: %v", err)), nil
	}

	roles, err := requiredParamArray[string](request.Params.Arguments, "roles")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get roles: %v", err)), nil
	}

	if len(roles) == 0 {
		return mcp.NewToolResultError("At least one role must be specified"), nil
	}

	// Get namespace name
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Grant subscription permission
	err = client.Namespaces().GrantSubPermission(*ns, subscription, roles)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to grant subscription permission: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Granted client roles %v to access the subscription %s of the namespace %s",
			roles, subscription, namespace),
	), nil
}

// handleRevokeSubscriptionPermission handles revoking subscription permissions from a role
func handleRevokeSubscriptionPermission(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	subscription, err := requiredParam[string](request.Params.Arguments, "subscription")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get subscription name: %v", err)), nil
	}

	role, err := requiredParam[string](request.Params.Arguments, "role")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get role: %v", err)), nil
	}

	// Get namespace name
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Revoke subscription permission
	err = client.Namespaces().RevokeSubPermission(*ns, subscription, role)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to revoke subscription permission: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Revoked client role %s permissions from subscription %s of namespace %s",
			role, subscription, namespace),
	), nil
}

// handleSetDispatchRate handles setting the default message dispatch rate for a namespace
func handleSetDispatchRate(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Get rate parameters
	msgRateStr, _ := optionalParam[string](request.Params.Arguments, "dispatchThrottlingRateInMsg")
	msgRate := -1 // Default value
	if msgRateStr != "" {
		parsedMsgRate, err := strconv.Atoi(msgRateStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid message rate parameter: %v", err)), nil
		}
		msgRate = parsedMsgRate
	}

	byteRateStr, _ := optionalParam[string](request.Params.Arguments, "dispatchThrottlingRateInByte")
	byteRate := int64(-1) // Default value
	if byteRateStr != "" {
		parsedByteRate, err := strconv.ParseInt(byteRateStr, 10, 64)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid byte rate parameter: %v", err)), nil
		}
		byteRate = parsedByteRate
	}

	ratePeriodStr, _ := optionalParam[string](request.Params.Arguments, "ratePeriodInSecond")
	ratePeriod := 1 // Default value
	if ratePeriodStr != "" {
		parsedRatePeriod, err := strconv.Atoi(ratePeriodStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid rate period parameter: %v", err)), nil
		}
		ratePeriod = parsedRatePeriod
	}

	// Get namespace name
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Create dispatch rate object
	rate := utils.DispatchRate{
		DispatchThrottlingRateInMsg:  msgRate,
		DispatchThrottlingRateInByte: byteRate,
		RatePeriodInSecond:           ratePeriod,
	}

	// Set dispatch rate
	err = client.Namespaces().SetDispatchRate(*ns, rate)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set dispatch rate: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Successfully set the default message dispatch rate of namespace %s to %+v", namespace, rate),
	), nil
}

// handleSetReplicatorDispatchRate handles setting the default replicator message dispatch rate for a namespace
func handleSetReplicatorDispatchRate(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Get rate parameters
	msgRateStr, _ := optionalParam[string](request.Params.Arguments, "dispatchThrottlingRateInMsg")
	msgRate := -1 // Default value
	if msgRateStr != "" {
		parsedMsgRate, err := strconv.Atoi(msgRateStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid message rate parameter: %v", err)), nil
		}
		msgRate = parsedMsgRate
	}

	byteRateStr, _ := optionalParam[string](request.Params.Arguments, "dispatchThrottlingRateInByte")
	byteRate := int64(-1) // Default value
	if byteRateStr != "" {
		parsedByteRate, err := strconv.ParseInt(byteRateStr, 10, 64)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid byte rate parameter: %v", err)), nil
		}
		byteRate = parsedByteRate
	}

	ratePeriodStr, _ := optionalParam[string](request.Params.Arguments, "ratePeriodInSecond")
	ratePeriod := 1 // Default value
	if ratePeriodStr != "" {
		parsedRatePeriod, err := strconv.Atoi(ratePeriodStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid rate period parameter: %v", err)), nil
		}
		ratePeriod = parsedRatePeriod
	}

	// Get namespace name
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Create dispatch rate object
	rate := utils.DispatchRate{
		DispatchThrottlingRateInMsg:  msgRate,
		DispatchThrottlingRateInByte: byteRate,
		RatePeriodInSecond:           ratePeriod,
	}

	// Set replicator dispatch rate
	err = client.Namespaces().SetReplicatorDispatchRate(*ns, rate)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set replicator dispatch rate: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Successfully set the default replicator message dispatch rate of namespace %s to %+v", namespace, rate),
	), nil
}

// handleSetSubscribeRate handles setting the default subscribe rate per consumer for a namespace
func handleSetSubscribeRate(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Get rate parameters
	subRateStr, _ := optionalParam[string](request.Params.Arguments, "subscribeThrottlingRatePerConsumer")
	subRate := -1 // Default value
	if subRateStr != "" {
		parsedSubRate, err := strconv.Atoi(subRateStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid subscribe rate value: %v", err)), nil
		}
		subRate = parsedSubRate
	}

	periodStr, _ := optionalParam[string](request.Params.Arguments, "ratePeriodInSecond")
	period := 30 // Default value
	if periodStr != "" {
		parsedPeriod, err := strconv.Atoi(periodStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid rate period value: %v", err)), nil
		}
		period = parsedPeriod
	}

	// Get namespace name
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Create subscription rate object
	rate := utils.SubscribeRate{
		SubscribeThrottlingRatePerConsumer: subRate,
		RatePeriodInSecond:                 period,
	}

	// Set subscribe rate
	err = client.Namespaces().SetSubscribeRate(*ns, rate)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set subscribe rate: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully set subscribe rate for namespace %s", namespace)), nil
}

// handleSetSubscriptionDispatchRate handles setting the default subscription message dispatch rate for a namespace
func handleSetSubscriptionDispatchRate(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Get rate parameters
	msgRateStr, _ := optionalParam[string](request.Params.Arguments, "dispatchThrottlingRateInMsg")
	msgRate := -1 // Default value
	if msgRateStr != "" {
		parsedMsgRate, err := strconv.Atoi(msgRateStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid message rate value: %v", err)), nil
		}
		msgRate = parsedMsgRate
	}

	byteRateStr, _ := optionalParam[string](request.Params.Arguments, "dispatchThrottlingRateInByte")
	byteRate := int64(-1) // Default value
	if byteRateStr != "" {
		parsedByteRate, err := strconv.ParseInt(byteRateStr, 10, 64)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid byte rate value: %v", err)), nil
		}
		byteRate = parsedByteRate
	}

	periodStr, _ := optionalParam[string](request.Params.Arguments, "ratePeriodInSecond")
	period := 1 // Default value
	if periodStr != "" {
		parsedPeriod, err := strconv.Atoi(periodStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid rate period value: %v", err)), nil
		}
		period = parsedPeriod
	}

	// Get namespace name
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Create dispatch rate object
	rate := utils.DispatchRate{
		DispatchThrottlingRateInMsg:  msgRate,
		DispatchThrottlingRateInByte: byteRate,
		RatePeriodInSecond:           period,
	}

	// Set subscription dispatch rate
	err = client.Namespaces().SetSubscriptionDispatchRate(*ns, rate)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set subscription dispatch rate: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Successfully set subscription dispatch rate for namespace %s", namespace),
	), nil
}

// handleSetPublishRate handles setting the default message publish rate of a namespace
func handleSetPublishRate(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Pulsar client
	client, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Parse namespace
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Initialize publish rate with defaults
	rate := utils.PublishRate{
		PublishThrottlingRateInMsg:  -1,
		PublishThrottlingRateInByte: -1,
	}

	// Get message rate if provided
	msgRateStr, hasMsgRate := optionalParam[string](request.Params.Arguments, "publishThrottlingRateInMsg")
	if hasMsgRate && msgRateStr != "" {
		msgRate, err := strconv.Atoi(msgRateStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid message rate: %v", err)), nil
		}
		rate.PublishThrottlingRateInMsg = msgRate
	}

	// Get byte rate if provided
	byteRateStr, hasByteRate := optionalParam[string](request.Params.Arguments, "publishThrottlingRateInByte")
	if hasByteRate && byteRateStr != "" {
		byteRate, err := strconv.ParseInt(byteRateStr, 10, 64)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid byte rate: %v", err)), nil
		}
		rate.PublishThrottlingRateInByte = byteRate
	}

	// Set publish rate
	err = client.Namespaces().SetPublishRate(*ns, rate)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set publish rate: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Successfully set the default message publish rate for namespace %s", namespace),
	), nil
}

// handleNamespaceSetPolicy handles setting policies for a namespace using the unified tool
func handleNamespaceSetPolicy(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	_, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	policyType, err := requiredParam[string](request.Params.Arguments, "policy")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get policy type: %v", err)), nil
	}

	// Route to the appropriate handler based on the policy type
	switch policyType {
	case "message-ttl":
		return handleSetMessageTTL(ctx, request)
	case "retention":
		return handleSetRetention(ctx, request)
	case "permission":
		return handleGrantPermission(ctx, request)
	case "replication-clusters":
		return handleSetReplicationClusters(ctx, request)
	case "backlog-quota":
		return handleSetBacklogQuota(ctx, request)
	case "topic-auto-creation":
		return handleSetTopicAutoCreation(ctx, request)
	case "schema-validation":
		return handleSetSchemaValidationEnforced(ctx, request)
	case "schema-auto-update":
		return handleSetSchemaAutoUpdateStrategy(ctx, request)
	case "auto-update-schema":
		return handleSetIsAllowAutoUpdateSchema(ctx, request)
	case "offload-threshold":
		return handleSetOffloadThreshold(ctx, request)
	case "offload-deletion-lag":
		return handleSetOffloadDeletionLag(ctx, request)
	case "compaction-threshold":
		return handleSetCompactionThreshold(ctx, request)
	case "max-producers-per-topic":
		return handleSetMaxProducersPerTopic(ctx, request)
	case "max-consumers-per-topic":
		return handleSetMaxConsumersPerTopic(ctx, request)
	case "max-consumers-per-subscription":
		return handleSetMaxConsumersPerSubscription(ctx, request)
	case "anti-affinity-group":
		return handleSetAntiAffinityGroup(ctx, request)
	case "persistence":
		return handleSetPersistence(ctx, request)
	case "deduplication":
		return handleSetDeduplication(ctx, request)
	case "encryption-required":
		return handleSetEncryptionRequired(ctx, request)
	case "subscription-auth-mode":
		return handleSetSubscriptionAuthMode(ctx, request)
	case "subscription-permission":
		return handleGrantSubscriptionPermission(ctx, request)
	case "dispatch-rate":
		return handleSetDispatchRate(ctx, request)
	case "replicator-dispatch-rate":
		return handleSetReplicatorDispatchRate(ctx, request)
	case "subscribe-rate":
		return handleSetSubscribeRate(ctx, request)
	case "subscription-dispatch-rate":
		return handleSetSubscriptionDispatchRate(ctx, request)
	case "publish-rate":
		return handleSetPublishRate(ctx, request)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unknown policy type: %s", policyType)), nil
	}
}

// handleNamespaceRemovePolicy handles removing policies from a namespace using the unified tool
func handleNamespaceRemovePolicy(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	_, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	policyType, err := requiredParam[string](request.Params.Arguments, "policy")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get policy type: %v", err)), nil
	}

	// Route to the appropriate handler based on the policy type
	switch policyType {
	case "backlog-quota":
		return handleRemoveBacklogQuota(ctx, request)
	case "topic-auto-creation":
		return handleRemoveTopicAutoCreation(ctx, request)
	case "offload-deletion-lag":
		return handleClearOffloadDeletionLag(ctx, request)
	case "anti-affinity-group":
		return handleDeleteAntiAffinityGroup(ctx, request)
	case "permission":
		return handleRevokePermission(ctx, request)
	case "subscription-permission":
		return handleRevokeSubscriptionPermission(ctx, request)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unknown or unsupported removable policy type: %s", policyType)), nil
	}
}
