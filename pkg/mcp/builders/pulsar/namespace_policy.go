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

package pulsar

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	pulsarctlutils "github.com/streamnative/pulsarctl/pkg/ctl/utils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

type pulsarAdminNamespacePolicyGetInput struct {
	Namespace string `json:"namespace"`
}

type pulsarAdminNamespacePolicySetInput struct {
	Namespace string   `json:"namespace"`
	Policy    string   `json:"policy"`
	Role      *string  `json:"role,omitempty"`
	Actions   []string `json:"actions,omitempty"`
	Clusters  []string `json:"clusters,omitempty"`
	Roles     []string `json:"roles,omitempty"`
	TTL       *string  `json:"ttl,omitempty"`
	Time      *string  `json:"time,omitempty"`
	Size      *string  `json:"size,omitempty"`
	LimitSize *string  `json:"limit-size,omitempty"`
	LimitTime *string  `json:"limit-time,omitempty"`
	Type      *string  `json:"type,omitempty"`
}

type pulsarAdminNamespacePolicyRemoveInput struct {
	Namespace    string  `json:"namespace"`
	Policy       string  `json:"policy"`
	Role         *string `json:"role,omitempty"`
	Subscription *string `json:"subscription,omitempty"`
	Type         *string `json:"type,omitempty"`
}

const (
	pulsarAdminNamespacePolicyGetToolDesc = "Get the configuration policies of a namespace. " +
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
	pulsarAdminNamespacePolicySetToolDesc = "Set a policy for a namespace. " +
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
	pulsarAdminNamespacePolicyRemoveToolDesc = "Remove a policy from a namespace. " +
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

	pulsarAdminNamespacePolicyGetNamespaceDesc = "The namespace name (tenant/namespace) to get policies for"
	pulsarAdminNamespacePolicySetNamespaceDesc = "The namespace name (tenant/namespace) to set the policy for"
	pulsarAdminNamespacePolicySetPolicyDesc    = "Type of policy to set. Available options: " +
		"message-ttl, retention, permission, replication-clusters, backlog-quota, " +
		"topic-auto-creation, schema-validation, schema-auto-update, auto-update-schema, " +
		"offload-threshold, offload-deletion-lag, compaction-threshold, " +
		"max-producers-per-topic, max-consumers-per-topic, max-consumers-per-subscription, " +
		"anti-affinity-group, persistence, deduplication, encryption-required, " +
		"subscription-auth-mode, subscription-permission, dispatch-rate, " +
		"replicator-dispatch-rate, subscribe-rate, subscription-dispatch-rate, publish-rate"
	pulsarAdminNamespacePolicySetRoleDesc         = "Role name for permission policies"
	pulsarAdminNamespacePolicySetActionsDesc      = "Actions to grant for permission policies (e.g., produce, consume)"
	pulsarAdminNamespacePolicySetClustersDesc     = "List of clusters for replication policies"
	pulsarAdminNamespacePolicySetRolesDesc        = "List of roles for subscription permission policies"
	pulsarAdminNamespacePolicySetTTLDesc          = "Message TTL in seconds (or 0 to disable TTL)"
	pulsarAdminNamespacePolicySetTimeDesc         = "Retention time in minutes, or special values: 0 (no retention) or -1 (infinite retention)"
	pulsarAdminNamespacePolicySetSizeDesc         = "Retention size limit (e.g., 10M, 16G, 3T), or special values: 0 (no retention) or -1 (infinite size retention)"
	pulsarAdminNamespacePolicySetLimitSizeDesc    = "Size limit for backlog quota (e.g., 10M, 16G)"
	pulsarAdminNamespacePolicySetLimitTimeDesc    = "Time limit in seconds for backlog quota. Default is -1 (infinite)"
	pulsarAdminNamespacePolicySetTypeDesc         = "Type of backlog quota to apply"
	pulsarAdminNamespacePolicyRemoveNamespaceDesc = "The namespace name (tenant/namespace) to remove the policy from"
	pulsarAdminNamespacePolicyRemovePolicyDesc    = "Type of policy to remove. Available options: " +
		"backlog-quota, topic-auto-creation, offload-deletion-lag, anti-affinity-group, " +
		"permission, subscription-permission"
	pulsarAdminNamespacePolicyRemoveRoleDesc         = "Role name for permission policies"
	pulsarAdminNamespacePolicyRemoveSubscriptionDesc = "Subscription name for subscription permission policies"
	pulsarAdminNamespacePolicyRemoveTypeDesc         = "Type of backlog quota to remove"
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
func (b *PulsarAdminNamespacePolicyToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	tools := []builders.ToolDefinition{}

	getTool, err := b.buildNamespaceGetPoliciesTool()
	if err != nil {
		return nil, err
	}
	getHandler := b.buildNamespaceGetPoliciesHandler()
	tools = append(tools, builders.ServerTool[pulsarAdminNamespacePolicyGetInput, any]{
		Tool:    getTool,
		Handler: getHandler,
	})

	if !config.ReadOnly {
		setTool, err := b.buildNamespaceSetPolicyTool()
		if err != nil {
			return nil, err
		}
		setHandler := b.buildNamespaceSetPolicyHandler()
		tools = append(tools, builders.ServerTool[pulsarAdminNamespacePolicySetInput, any]{
			Tool:    setTool,
			Handler: setHandler,
		})

		removeTool, err := b.buildNamespaceRemovePolicyTool()
		if err != nil {
			return nil, err
		}
		removeHandler := b.buildNamespaceRemovePolicyHandler()
		tools = append(tools, builders.ServerTool[pulsarAdminNamespacePolicyRemoveInput, any]{
			Tool:    removeTool,
			Handler: removeHandler,
		})
	}

	return tools, nil
}

// buildNamespaceGetPoliciesTool builds the get policies tool
func (b *PulsarAdminNamespacePolicyToolBuilder) buildNamespaceGetPoliciesTool() (*sdk.Tool, error) {
	inputSchema, err := buildPulsarAdminNamespacePolicyGetInputSchema()
	if err != nil {
		return nil, err
	}

	return &sdk.Tool{
		Name:        "pulsar_admin_namespace_policy_get",
		Description: pulsarAdminNamespacePolicyGetToolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildNamespaceSetPolicyTool builds the set policy tool
func (b *PulsarAdminNamespacePolicyToolBuilder) buildNamespaceSetPolicyTool() (*sdk.Tool, error) {
	inputSchema, err := buildPulsarAdminNamespacePolicySetInputSchema()
	if err != nil {
		return nil, err
	}

	return &sdk.Tool{
		Name:        "pulsar_admin_namespace_policy_set",
		Description: pulsarAdminNamespacePolicySetToolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildNamespaceRemovePolicyTool builds the remove policy tool
func (b *PulsarAdminNamespacePolicyToolBuilder) buildNamespaceRemovePolicyTool() (*sdk.Tool, error) {
	inputSchema, err := buildPulsarAdminNamespacePolicyRemoveInputSchema()
	if err != nil {
		return nil, err
	}

	return &sdk.Tool{
		Name:        "pulsar_admin_namespace_policy_remove",
		Description: pulsarAdminNamespacePolicyRemoveToolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildNamespaceGetPoliciesHandler builds the get policies handler
func (b *PulsarAdminNamespacePolicyToolBuilder) buildNamespaceGetPoliciesHandler() builders.ToolHandlerFunc[pulsarAdminNamespacePolicyGetInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input pulsarAdminNamespacePolicyGetInput) (*sdk.CallToolResult, any, error) {
		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return nil, nil, fmt.Errorf("pulsar session not found in context")
		}

		client, err := session.GetAdminClient()
		if err != nil {
			return nil, nil, b.handleError("get admin client", err)
		}

		if input.Namespace == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'namespace'")
		}

		// Get policies
		policies, err := client.Namespaces().GetPolicies(input.Namespace)
		if err != nil {
			return nil, nil, b.handleError("get policies", err)
		}

		result, err := b.marshalResponse(policies)
		return result, nil, err
	}
}

// buildNamespaceSetPolicyHandler builds the set policy handler
func (b *PulsarAdminNamespacePolicyToolBuilder) buildNamespaceSetPolicyHandler() builders.ToolHandlerFunc[pulsarAdminNamespacePolicySetInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input pulsarAdminNamespacePolicySetInput) (*sdk.CallToolResult, any, error) {
		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return nil, nil, fmt.Errorf("pulsar session not found in context")
		}

		client, err := session.GetAdminClient()
		if err != nil {
			return nil, nil, b.handleError("get admin client", err)
		}

		namespace := input.Namespace
		if namespace == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'namespace'")
		}

		policy := input.Policy
		if policy == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'policy'")
		}

		// Handle different policy types
		switch policy {
		case "message-ttl":
			result, handlerErr := b.handleSetMessageTTL(client, namespace, input)
			return result, nil, handlerErr
		case "retention":
			result, handlerErr := b.handleSetRetention(client, namespace, input)
			return result, nil, handlerErr
		case "permission":
			result, handlerErr := b.handleGrantPermission(client, namespace, input)
			return result, nil, handlerErr
		case "replication-clusters":
			result, handlerErr := b.handleSetReplicationClusters(client, namespace, input)
			return result, nil, handlerErr
		case "backlog-quota":
			result, handlerErr := b.handleSetBacklogQuota(client, namespace, input)
			return result, nil, handlerErr
		default:
			return nil, nil, fmt.Errorf("unsupported policy type: %s", policy)
		}
	}
}

// buildNamespaceRemovePolicyHandler builds the remove policy handler
func (b *PulsarAdminNamespacePolicyToolBuilder) buildNamespaceRemovePolicyHandler() builders.ToolHandlerFunc[pulsarAdminNamespacePolicyRemoveInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input pulsarAdminNamespacePolicyRemoveInput) (*sdk.CallToolResult, any, error) {
		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return nil, nil, fmt.Errorf("pulsar session not found in context")
		}

		client, err := session.GetAdminClient()
		if err != nil {
			return nil, nil, b.handleError("get admin client", err)
		}

		namespace := input.Namespace
		if namespace == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'namespace'")
		}

		policy := input.Policy
		if policy == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'policy'")
		}

		// Handle different policy types
		switch policy {
		case "permission":
			result, handlerErr := b.handleRevokePermission(client, namespace, input)
			return result, nil, handlerErr
		case "backlog-quota":
			result, handlerErr := b.handleRemoveBacklogQuota(client, namespace)
			return result, nil, handlerErr
		default:
			return nil, nil, fmt.Errorf("unsupported policy type for removal: %s", policy)
		}
	}
}

// Utility functions
func (b *PulsarAdminNamespacePolicyToolBuilder) handleError(operation string, err error) error {
	return fmt.Errorf("failed to %s: %v", operation, err)
}

func (b *PulsarAdminNamespacePolicyToolBuilder) marshalResponse(data interface{}) (*sdk.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, b.handleError("marshal response", err)
	}
	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: string(jsonBytes)}},
	}, nil
}

// Policy-specific handler functions

// handleSetMessageTTL handles setting message TTL for a namespace
func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetMessageTTL(client cmdutils.Client, namespace string, input pulsarAdminNamespacePolicySetInput) (*sdk.CallToolResult, error) {
	ttlStr, err := requireString(input.TTL, "ttl")
	if err != nil {
		return nil, fmt.Errorf("failed to get TTL: %v", err)
	}

	ttl, err := strconv.ParseInt(ttlStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid TTL value, must be an integer: %v", err)
	}

	// Set message TTL
	err = client.Namespaces().SetNamespaceMessageTTL(namespace, int(ttl))
	if err != nil {
		return nil, fmt.Errorf("failed to set message TTL: %v", err)
	}

	return textResult(fmt.Sprintf("Set message TTL for %s to %d seconds", namespace, ttl)), nil
}

// handleSetRetention handles setting retention for a namespace
func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetRetention(client cmdutils.Client, namespace string, input pulsarAdminNamespacePolicySetInput) (*sdk.CallToolResult, error) {
	timeStr := ""
	if input.Time != nil {
		timeStr = *input.Time
	}

	sizeStr := ""
	if input.Size != nil {
		sizeStr = *input.Size
	}

	if timeStr == "" && sizeStr == "" {
		return nil, fmt.Errorf("at least one of 'time' or 'size' must be specified")
	}

	// Parse retention time
	var retentionTimeInMin int
	if timeStr != "" {
		// Parse relative time in seconds from the input string
		retentionTime, err := pulsarctlutils.ParseRelativeTimeInSeconds(timeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid retention time format: %v", err)
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
				return nil, fmt.Errorf("invalid retention size format: %v", err)
			}

			if sizeInBytes != -1 {
				// Convert bytes to MB
				retentionSizeInMB = int(sizeInBytes / (1024 * 1024))
				if retentionSizeInMB < 1 {
					return nil, fmt.Errorf("retention size must be at least 1MB")
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
		return nil, fmt.Errorf("failed to set retention: %v", err)
	}

	return textResult(fmt.Sprintf("Set retention for %s successfully", namespace)), nil
}

// handleGrantPermission handles granting permissions on a namespace
func (b *PulsarAdminNamespacePolicyToolBuilder) handleGrantPermission(client cmdutils.Client, namespace string, input pulsarAdminNamespacePolicySetInput) (*sdk.CallToolResult, error) {
	role, err := requireString(input.Role, "role")
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %v", err)
	}

	actions, err := requireStringSlice(input.Actions, "actions")
	if err != nil {
		return nil, fmt.Errorf("failed to get actions: %v", err)
	}

	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return nil, fmt.Errorf("invalid namespace name: %v", err)
	}

	a, err := b.parseActions(actions)
	if err != nil {
		return nil, fmt.Errorf("failed to parse actions: %v", err)
	}

	// Grant permissions
	err = client.Namespaces().GrantNamespacePermission(*ns, role, a)
	if err != nil {
		return nil, fmt.Errorf("failed to grant permission: %v", err)
	}

	return textResult(fmt.Sprintf("Granted %v permission(s) to role %s on %s", actions, role, namespace)), nil
}

// handleRevokePermission handles revoking permissions from a namespace
func (b *PulsarAdminNamespacePolicyToolBuilder) handleRevokePermission(client cmdutils.Client, namespace string, input pulsarAdminNamespacePolicyRemoveInput) (*sdk.CallToolResult, error) {
	role, err := requireString(input.Role, "role")
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %v", err)
	}

	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return nil, fmt.Errorf("invalid namespace name: %v", err)
	}

	// Revoke permissions
	err = client.Namespaces().RevokeNamespacePermission(*ns, role)
	if err != nil {
		return nil, fmt.Errorf("failed to revoke permission: %v", err)
	}

	return textResult(fmt.Sprintf("Revoked all permissions from role %s on %s", role, namespace)), nil
}

// handleSetReplicationClusters handles setting replication clusters for a namespace
func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetReplicationClusters(client cmdutils.Client, namespace string, input pulsarAdminNamespacePolicySetInput) (*sdk.CallToolResult, error) {
	clusters, err := requireStringSlice(input.Clusters, "clusters")
	if err != nil {
		return nil, fmt.Errorf("failed to get clusters: %v", err)
	}

	if len(clusters) == 0 {
		return nil, fmt.Errorf("at least one cluster must be specified")
	}

	// Set replication clusters
	err = client.Namespaces().SetNamespaceReplicationClusters(namespace, clusters)
	if err != nil {
		return nil, fmt.Errorf("failed to set replication clusters: %v", err)
	}

	return textResult(fmt.Sprintf("Set replication clusters for %s to %s", namespace, strings.Join(clusters, ", "))), nil
}

// handleSetBacklogQuota handles setting backlog quota for a namespace
func (b *PulsarAdminNamespacePolicyToolBuilder) handleSetBacklogQuota(client cmdutils.Client, namespace string, input pulsarAdminNamespacePolicySetInput) (*sdk.CallToolResult, error) {
	limitSizeStr, err := requireString(input.LimitSize, "limit-size")
	if err != nil {
		return nil, fmt.Errorf("failed to get limit size: %v", err)
	}

	policyStr := input.Policy
	if policyStr == "" {
		return nil, fmt.Errorf("failed to get policy: required argument \"policy\" not found")
	}

	// Parse backlog size limit
	limitSize, err := pulsarctlutils.ValidateSizeString(limitSizeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid limit size format: %v", err)
	}

	// Parse backlog quota policy using the ParseRetentionPolicy function
	policy, err := utils.ParseRetentionPolicy(policyStr)
	if err != nil {
		return nil, fmt.Errorf("invalid backlog quota policy: %s. Valid options: producer_request_hold, producer_exception, consumer_backlog_eviction", policyStr)
	}

	// Get optional time limit
	limitTimeStr := "-1"
	if input.LimitTime != nil {
		limitTimeStr = *input.LimitTime
	}
	limitTime, err := strconv.ParseInt(limitTimeStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid limit time: %v", err)
	}

	// Parse quota type (optional, default to destination_storage)
	quotaTypeStr := "destination_storage"
	if input.Type != nil && *input.Type != "" {
		quotaTypeStr = *input.Type
	}
	quotaType := utils.DestinationStorage
	if quotaTypeStr != "" {
		parsedType, err := utils.ParseBacklogQuotaType(quotaTypeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid backlog quota type: %v", err)
		}
		quotaType = parsedType
	}

	// Create and set backlog quota
	backlogQuota := utils.NewBacklogQuota(limitSize, limitTime, policy)
	err = client.Namespaces().SetBacklogQuota(namespace, backlogQuota, quotaType)
	if err != nil {
		return nil, fmt.Errorf("failed to set backlog quota: %v", err)
	}

	return textResult(fmt.Sprintf("Set backlog quota for %s successfully", namespace)), nil
}

// handleRemoveBacklogQuota handles removing backlog quota for a namespace
func (b *PulsarAdminNamespacePolicyToolBuilder) handleRemoveBacklogQuota(client cmdutils.Client, namespace string) (*sdk.CallToolResult, error) {
	// Remove backlog quota (API doesn't require quota type for removal)
	err := client.Namespaces().RemoveBacklogQuota(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to remove backlog quota: %v", err)
	}

	return textResult(fmt.Sprintf("Removed backlog quota for %s successfully", namespace)), nil
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

func buildPulsarAdminNamespacePolicyGetInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[pulsarAdminNamespacePolicyGetInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}
	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}

	setSchemaDescription(schema, "namespace", pulsarAdminNamespacePolicyGetNamespaceDesc)
	normalizeAdditionalProperties(schema)
	return schema, nil
}

func buildPulsarAdminNamespacePolicySetInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[pulsarAdminNamespacePolicySetInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}
	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}

	setSchemaDescription(schema, "namespace", pulsarAdminNamespacePolicySetNamespaceDesc)
	setSchemaDescription(schema, "policy", pulsarAdminNamespacePolicySetPolicyDesc)
	setSchemaDescription(schema, "role", pulsarAdminNamespacePolicySetRoleDesc)
	setSchemaDescription(schema, "actions", pulsarAdminNamespacePolicySetActionsDesc)
	setSchemaDescription(schema, "clusters", pulsarAdminNamespacePolicySetClustersDesc)
	setSchemaDescription(schema, "roles", pulsarAdminNamespacePolicySetRolesDesc)
	setSchemaDescription(schema, "ttl", pulsarAdminNamespacePolicySetTTLDesc)
	setSchemaDescription(schema, "time", pulsarAdminNamespacePolicySetTimeDesc)
	setSchemaDescription(schema, "size", pulsarAdminNamespacePolicySetSizeDesc)
	setSchemaDescription(schema, "limit-size", pulsarAdminNamespacePolicySetLimitSizeDesc)
	setSchemaDescription(schema, "limit-time", pulsarAdminNamespacePolicySetLimitTimeDesc)
	setSchemaDescription(schema, "type", pulsarAdminNamespacePolicySetTypeDesc)

	if actionsSchema := schema.Properties["actions"]; actionsSchema != nil && actionsSchema.Items != nil {
		actionsSchema.Items.Description = "action"
	}
	if clustersSchema := schema.Properties["clusters"]; clustersSchema != nil && clustersSchema.Items != nil {
		clustersSchema.Items.Description = "cluster"
	}
	if rolesSchema := schema.Properties["roles"]; rolesSchema != nil && rolesSchema.Items != nil {
		rolesSchema.Items.Description = "role"
	}

	normalizeAdditionalProperties(schema)
	return schema, nil
}

func buildPulsarAdminNamespacePolicyRemoveInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[pulsarAdminNamespacePolicyRemoveInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}
	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}

	setSchemaDescription(schema, "namespace", pulsarAdminNamespacePolicyRemoveNamespaceDesc)
	setSchemaDescription(schema, "policy", pulsarAdminNamespacePolicyRemovePolicyDesc)
	setSchemaDescription(schema, "role", pulsarAdminNamespacePolicyRemoveRoleDesc)
	setSchemaDescription(schema, "subscription", pulsarAdminNamespacePolicyRemoveSubscriptionDesc)
	setSchemaDescription(schema, "type", pulsarAdminNamespacePolicyRemoveTypeDesc)

	normalizeAdditionalProperties(schema)
	return schema, nil
}
