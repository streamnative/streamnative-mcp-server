// Copyright (c) 2025 StreamNative, Inc. All Rights Reserved.

package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// PulsarAdminAddTopicPolicyTools adds topic policy-related tools to the MCP server
func PulsarAdminAddTopicPolicyTools(s *server.MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeaturePulsarAdminTopicPolicy)) && !slices.Contains(features, string(FeatureAll)) && !slices.Contains(features, string(FeatureAllPulsar)) {
		return
	}

	// Add unified topic policy get tool
	topicGetPolicyTool := mcp.NewTool("pulsar_admin_topic_policy_get",
		mcp.WithDescription("Get a policy configuration for a topic. "+
			"This is a unified tool for retrieving different types of policies on a topic. "+
			"The policy type determines which specific policy will be retrieved. "+
			"Required permissions vary based on the policy being queried.\n\n"+
			"Available policy types and their parameters:\n"+
			"1. publish-rate: Gets the publish rate configuration limiting the maximum rate at which producers can publish messages\n"+
			"   - Required: topic, policy\n"+
			"   - Returns: publishThrottlingRateInMsg (maximum message rate per second), publishThrottlingRateInByte (maximum byte rate per second)\n\n"+
			"2. permissions: Gets permissions on a topic showing which roles have which actions (produce/consume) allowed\n"+
			"   - Required: topic, policy\n"+
			"   - Returns: Map of roles to their allowed actions\n\n"+
			"3. message-ttl: Gets time to live for messages in seconds, after which unacknowledged messages are deleted\n"+
			"   - Required: topic, policy\n"+
			"   - Returns: TTL value in seconds, or null if not set\n\n"+
			"4. max-producers: Gets maximum number of producers allowed on this topic\n"+
			"   - Required: topic, policy\n"+
			"   - Returns: Maximum number of producers allowed, or -1 if unlimited\n\n"+
			"5. max-consumers: Gets maximum number of consumers allowed on this topic\n"+
			"   - Required: topic, policy\n"+
			"   - Returns: Maximum number of consumers allowed, or -1 if unlimited\n\n"+
			"6. max-unack-messages-per-consumer: Gets maximum number of unacknowledged messages allowed per consumer\n"+
			"   - Required: topic, policy\n"+
			"   - Returns: Maximum number of unacked messages per consumer, or -1 if unlimited\n\n"+
			"7. max-unack-messages-per-subscription: Gets maximum number of unacknowledged messages allowed per subscription\n"+
			"   - Required: topic, policy\n"+
			"   - Returns: Maximum number of unacked messages per subscription, or -1 if unlimited\n\n"+
			"8. persistence: Gets persistence policy defining how messages are persisted on this topic\n"+
			"   - Required: topic, policy\n"+
			"   - Returns: ensemble size, write quorum size, ack quorum size, and optional mark-delete rate\n\n"+
			"9. delayed-delivery: Gets delayed delivery policy for the topic\n"+
			"   - Required: topic, policy\n"+
			"   - Returns: tickTime (delay precision in ms) and active status\n\n"+
			"10. dispatch-rate: Gets the message dispatch rate limiting policy\n"+
			"    - Required: topic, policy\n"+
			"    - Returns: dispatchThrottlingRateInMsg (maximum message rate), dispatchThrottlingRateInByte (maximum byte rate), ratePeriodInSecond\n\n"+
			"11. deduplication-status: Gets whether producer-side message deduplication is enabled\n"+
			"    - Required: topic, policy\n"+
			"    - Returns: enabled/disabled status\n\n"+
			"12. retention: Gets message retention policy defining how long messages are retained\n"+
			"    - Required: topic, policy\n"+
			"    - Returns: retention size limit and retention time limit\n\n"+
			"13. backlog-quota: Gets backlog quota policies for this topic\n"+
			"    - Required: topic, policy\n"+
			"    - Returns: size limits, time limits, and enforcement policy\n\n"+
			"14. compaction-threshold: Gets the threshold in bytes for when compaction is triggered\n"+
			"    - Required: topic, policy\n"+
			"    - Returns: threshold size in bytes, or 0 if disabled\n\n"+
			"15. inactive-topic: Gets inactive topic policy for automatic topic cleanup\n"+
			"    - Required: topic, policy\n"+
			"    - Returns: inactivity duration, deletion mode (delete_when_no_subscriptions/delete_when_subscriptions_caught_up)\n\n"+
			"This tool requires appropriate topic admin or read permissions. The response format is JSON."),
		mcp.WithString("topic", mcp.Required(),
			mcp.Description("The topic name to get policy for (persistent://tenant/namespace/topic)"),
		),
		mcp.WithString("policy", mcp.Required(),
			mcp.Description("Type of policy to get. Available options: "+
				"publish-rate, permissions, message-ttl, max-producers, max-consumers, "+
				"max-unack-messages-per-consumer, max-unack-messages-per-subscription, persistence, "+
				"delayed-delivery, dispatch-rate, deduplication-status, retention, backlog-quota, "+
				"compaction-threshold, inactive-topic"),
		),
	)
	s.AddTool(topicGetPolicyTool, handleTopicGetPolicy)

	// Only add write operations if not in read-only mode
	if !readOnly {
		// Add unified topic set policy tool
		topicSetPolicyTool := mcp.NewTool("pulsar_admin_topic_policy_set",
			mcp.WithDescription("Set a policy for a topic. "+
				"This is a unified tool for setting different types of policies on a topic. "+
				"The policy type determines which specific policy will be set, and the required parameters "+
				"vary based on the policy type. "+
				"Requires appropriate admin permissions based on the policy being modified.\n\n"+
				"Available policy types and their required parameters:\n"+
				"1. publish-rate: Sets maximum rate at which producers can publish messages\n"+
				"   - Required: topic, policy\n"+
				"   - Optional: publishThrottlingRateInMsg (messages per second, -1 for unlimited), "+
				"publishThrottlingRateInByte (bytes per second, -1 for unlimited)\n\n"+
				"2. permissions: Grants permissions to a role\n"+
				"   - Required: topic, policy, role, actions (array of permissions: produce, consume)\n"+
				"   - Example: actions=[\"produce\", \"consume\"] to allow both producing and consuming\n\n"+
				"3. message-ttl: Sets time to live for messages\n"+
				"   - Required: topic, policy, ttl (in seconds, 0 to disable)\n"+
				"   - After TTL expires, unacknowledged messages are deleted\n\n"+
				"4. max-producers: Sets maximum number of producers allowed\n"+
				"   - Required: topic, policy, max (integer > 0, or -1 for unlimited)\n"+
				"   - New producers will be rejected once limit is reached\n\n"+
				"5. max-consumers: Sets maximum number of consumers allowed\n"+
				"   - Required: topic, policy, max (integer > 0, or -1 for unlimited)\n"+
				"   - New consumers will be rejected once limit is reached\n\n"+
				"6. max-unack-messages-per-consumer: Sets maximum unacked messages per consumer\n"+
				"   - Required: topic, policy, max (integer > 0, or -1 for unlimited)\n"+
				"   - Consumer will be blocked from receiving more messages until acknowledgment\n\n"+
				"7. max-unack-messages-per-subscription: Sets maximum unacked messages per subscription\n"+
				"   - Required: topic, policy, max (integer > 0, or -1 for unlimited)\n"+
				"   - All consumers on the subscription will be blocked from receiving more messages until acknowledgment\n\n"+
				"8. persistence: Sets persistence configuration\n"+
				"   - Required: topic, policy, ensemble-size, write-quorum-size, ack-quorum-size\n"+
				"   - Optional: ml-mark-delete-max-rate\n"+
				"   - Ensemble size must be >= write quorum size >= ack quorum size\n\n"+
				"9. delayed-delivery: Configures delayed message delivery\n"+
				"   - Required: topic, policy\n"+
				"   - Optional: tickTime (precision in ms), delay-enabled (true/false)\n\n"+
				"10. dispatch-rate: Sets message dispatch rate\n"+
				"    - Required: topic, policy\n"+
				"    - Optional: dispatchThrottlingRateInMsg (messages per second), dispatchThrottlingRateInByte (bytes per second), ratePeriodInSecond\n\n"+
				"11. deduplication-status: Sets message deduplication status\n"+
				"    - Required: topic, policy\n"+
				"    - Optional: enable (true/false)\n\n"+
				"12. retention: Sets retention policy for messages\n"+
				"    - Required: topic, policy\n"+
				"    - Optional: time (retention time in minutes, 0=no retention, -1=infinite), "+
				"size (e.g., 10M, 16G, 3T, 0=no retention, -1=infinite)\n\n"+
				"13. backlog-quota: Sets backlog quota policy\n"+
				"    - Required: topic, policy, limit-size (e.g., 10M, 16G), policy (producer_request_hold, producer_exception, consumer_backlog_eviction)\n"+
				"    - Optional: limit-time (seconds, -1=infinite), type (destination_storage, message_age)\n\n"+
				"14. compaction-threshold: Sets threshold for topic compaction\n"+
				"    - Required: topic, policy, threshold (e.g., 10M, 16G, 3T, 0 to disable)\n"+
				"    - Compaction runs when the topic reaches this size\n\n"+
				"15. inactive-topic: Sets inactive topic policy\n"+
				"    - Required: topic, policy\n"+
				"    - Optional: inactive-time-seconds, deletion-mode (delete_when_no_subscriptions/delete_when_subscriptions_caught_up)\n"+
				"    - Controls automatic cleanup of unused topics\n\n"+
				"This tool requires appropriate topic admin permissions. The response will confirm successful application of the policy."),
			mcp.WithString("topic", mcp.Required(),
				mcp.Description("The topic name to set policy for (persistent://tenant/namespace/topic)"),
			),
			mcp.WithString("policy", mcp.Required(),
				mcp.Description("Type of policy to set. Available options: "+
					"publish-rate, permissions, message-ttl, max-producers, max-consumers, "+
					"max-unack-messages-per-consumer, max-unack-messages-per-subscription, persistence, "+
					"delayed-delivery, dispatch-rate, deduplication-status, retention, backlog-quota, "+
					"compaction-threshold, inactive-topic"),
			),
			// Generic policy parameters - specific ones will be used based on the policy type
			mcp.WithString("role",
				mcp.Description("Role name for permission policies (e.g., 'admin', 'client', 'producer', 'consumer')"),
			),
			mcp.WithArray("actions",
				mcp.Description("Actions to grant for permission policies. Valid values are: produce, consume"),
			),
			mcp.WithString("ttl",
				mcp.Description("Message TTL in seconds (or 0 to disable TTL)"),
			),
			mcp.WithString("max",
				mcp.Description("Maximum value for producer/consumer/unacked message limit policies. Use -1 for unlimited."),
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
				mcp.Description("Type parameter whose meaning depends on the policy type being set (e.g., backlog quota type)"),
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
			mcp.WithString("dispatchThrottlingRateInMsg",
				mcp.Description("Message dispatch/throttling rate in messages per second"),
			),
			mcp.WithString("dispatchThrottlingRateInByte",
				mcp.Description("Message dispatch/throttling rate in bytes per second"),
			),
			mcp.WithString("ratePeriodInSecond",
				mcp.Description("Rate period in seconds for rate limiting policies"),
			),
			mcp.WithString("publishThrottlingRateInMsg",
				mcp.Description("Message publish rate in messages per second"),
			),
			mcp.WithString("publishThrottlingRateInByte",
				mcp.Description("Message publish rate in bytes per second"),
			),
			mcp.WithString("tickTime",
				mcp.Description("Tick time for delayed message delivery precision in milliseconds"),
			),
			mcp.WithString("delay-enabled",
				mcp.Description("Whether delayed delivery is enabled (true/false)"),
			),
			mcp.WithString("enable",
				mcp.Description("Enable flag for various policies (true/false)"),
			),
			mcp.WithString("threshold",
				mcp.Description("Threshold value for compaction in bytes (e.g., 10M, 16G, 3T, 0 to disable)"),
			),
			mcp.WithString("inactive-time-seconds",
				mcp.Description("Time in seconds after which an inactive topic will be cleaned up"),
			),
			mcp.WithString("deletion-mode",
				mcp.Description("Mode for inactive topic deletion (delete_when_no_subscriptions or delete_when_subscriptions_caught_up)"),
			),
		)
		s.AddTool(topicSetPolicyTool, handleTopicSetPolicy)

		// Add unified topic remove policy tool
		topicRemovePolicyTool := mcp.NewTool("pulsar_admin_topic_policy_remove",
			mcp.WithDescription("Remove a policy from a topic. "+
				"This is a unified tool for removing different types of policies from a topic. "+
				"The policy type determines which specific policy will be removed. "+
				"Requires appropriate admin permissions based on the policy being removed.\n\n"+
				"Available policy types to remove and their required parameters:\n"+
				"1. publish-rate: Removes the publish rate limiting policy\n"+
				"   - Required: topic, policy\n"+
				"   - Topic will inherit namespace-level publish rate limits if configured\n\n"+
				"2. permissions: Revokes all permissions from a role\n"+
				"   - Required: topic, policy, role\n"+
				"   - The specified role will no longer have any access to the topic\n\n"+
				"3. message-ttl: Removes the message TTL policy\n"+
				"   - Required: topic, policy\n"+
				"   - Topic will inherit namespace-level TTL policy if configured\n\n"+
				"4. max-producers: Removes the maximum producers limit\n"+
				"   - Required: topic, policy\n"+
				"   - Topic will inherit namespace-level limit if configured\n\n"+
				"5. max-consumers: Removes the maximum consumers limit\n"+
				"   - Required: topic, policy\n"+
				"   - Topic will inherit namespace-level limit if configured\n\n"+
				"6. max-unack-messages-per-consumer: Removes the unacked messages per consumer limit\n"+
				"   - Required: topic, policy\n"+
				"   - Topic will inherit namespace-level limit if configured\n\n"+
				"7. max-unack-messages-per-subscription: Removes the unacked messages per subscription limit\n"+
				"   - Required: topic, policy\n"+
				"   - Topic will inherit namespace-level limit if configured\n\n"+
				"8. persistence: Removes custom persistence policy\n"+
				"   - Required: topic, policy\n"+
				"   - Topic will inherit namespace-level persistence policy\n\n"+
				"9. delayed-delivery: Removes delayed delivery policy\n"+
				"   - Required: topic, policy\n"+
				"   - Topic will inherit namespace-level delayed delivery policy if configured\n\n"+
				"10. dispatch-rate: Removes the dispatch rate limiting policy\n"+
				"    - Required: topic, policy\n"+
				"    - Topic will inherit namespace-level dispatch rate limits if configured\n\n"+
				"11. deduplication-status: Removes the deduplication status policy\n"+
				"    - Required: topic, policy\n"+
				"    - Topic will inherit namespace-level deduplication policy if configured\n\n"+
				"12. retention: Removes the retention policy\n"+
				"    - Required: topic, policy\n"+
				"    - Topic will inherit namespace-level retention policy if configured\n\n"+
				"13. backlog-quota: Removes the backlog quota policy\n"+
				"    - Required: topic, policy\n"+
				"    - Optional: type (destination_storage, message_age)\n"+
				"    - Topic will inherit namespace-level backlog quota policy if configured\n\n"+
				"14. compaction-threshold: Removes the compaction threshold policy\n"+
				"    - Required: topic, policy\n"+
				"    - Topic will inherit namespace-level compaction threshold if configured\n\n"+
				"15. inactive-topic: Removes the inactive topic policy\n"+
				"    - Required: topic, policy\n"+
				"    - Topic will inherit namespace-level inactive topic policy if configured\n\n"+
				"After removing a topic-level policy, the topic will inherit the corresponding namespace-level policy (if any). "+
				"This tool requires appropriate topic admin permissions. The response will confirm successful removal of the policy."),
			mcp.WithString("topic", mcp.Required(),
				mcp.Description("The topic name to remove the policy from (persistent://tenant/namespace/topic)"),
			),
			mcp.WithString("policy", mcp.Required(),
				mcp.Description("Type of policy to remove. Available options: "+
					"publish-rate, permissions, message-ttl, max-producers, max-consumers, "+
					"max-unack-messages-per-consumer, max-unack-messages-per-subscription, persistence, "+
					"delayed-delivery, dispatch-rate, deduplication-status, retention, backlog-quota, "+
					"compaction-threshold, inactive-topic"),
			),
			mcp.WithString("role",
				mcp.Description("Role name for revoking permissions"),
			),
			mcp.WithString("type",
				mcp.Description("Type of backlog quota to remove (destination_storage or message_age)"),
			),
		)
		s.AddTool(topicRemovePolicyTool, handleTopicRemovePolicy)
	}
}

// handleTopicsGetPublishRate gets the publish rate for a topic
func handleTopicsGetPublishRate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Get publish rate
	publishRate, err := admin.Topics().GetPublishRate(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic publish rate: %v", err)), nil
	}

	// Format the output
	jsonBytes, err := json.Marshal(publishRate)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal topic publish rate: %v", err)), nil
	}

	// Check if publish rate is disabled/unlimited
	if publishRate.PublishThrottlingRateInMsg < 0 && publishRate.PublishThrottlingRateInByte < 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Publish rate is unlimited for topic %s", topicName.String())), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleTopicsSetPublishRate sets the publish rate for a topic
func handleTopicsSetPublishRate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Default values for publish rate
	publishRateInMsg := int64(-1)  // unlimited
	publishRateInByte := int64(-1) // unlimited

	// Get publish rate in messages if provided
	msgRateParam, hasMsgRate := optionalParam[float64](request.Params.Arguments, "publishThrottlingRateInMsg")
	if hasMsgRate {
		publishRateInMsg = int64(msgRateParam)
	}

	// Get publish rate in bytes if provided
	byteRateParam, hasByteRate := optionalParam[float64](request.Params.Arguments, "publishThrottlingRateInByte")
	if hasByteRate {
		publishRateInByte = int64(byteRateParam)
	}

	// Create publish rate policy
	publishRate := utils.PublishRateData{
		PublishThrottlingRateInMsg:  publishRateInMsg,
		PublishThrottlingRateInByte: publishRateInByte,
	}

	// Set publish rate
	err = admin.Topics().SetPublishRate(*topicName, publishRate)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set topic publish rate: %v", err)), nil
	}

	// Format the response message
	if publishRateInMsg < 0 && publishRateInByte < 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Publish rate set to unlimited for topic %s", topicName.String())), nil
	}

	var msgText string
	if publishRateInMsg < 0 {
		msgText = "unlimited messages"
	} else {
		msgText = fmt.Sprintf("%d messages", publishRateInMsg)
	}

	var byteText string
	if publishRateInByte < 0 {
		byteText = "unlimited bytes"
	} else {
		byteText = fmt.Sprintf("%d bytes", publishRateInByte)
	}

	return mcp.NewToolResultText(fmt.Sprintf("Publish rate set for topic %s: %s and %s per second",
		topicName.String(), msgText, byteText)), nil
}

// handleTopicsRemovePublishRate removes the publish rate for a topic
func handleTopicsRemovePublishRate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Remove publish rate
	err = admin.Topics().RemovePublishRate(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to remove topic publish rate: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Publish rate removed for topic %s", topicName.String())), nil
}

// handleTopicsGetPermissions gets the permissions on a topic
func handleTopicsGetPermissions(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Get permissions
	permissions, err := admin.Topics().GetPermissions(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic permissions: %v", err)), nil
	}

	// Format the output
	jsonBytes, err := json.Marshal(permissions)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal topic permissions: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleTopicsGrantPermissions grants a new permission to a role on a topic
func handleTopicsGrantPermissions(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	role, err := requiredParam[string](request.Params.Arguments, "role")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get role: %v", err)), nil
	}

	actions, err := requiredParamArray[string](request.Params.Arguments, "actions")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get actions: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	var validActions []string
	for _, action := range actions {
		actionStr := strings.ToLower(action)
		switch actionStr {
		case "produce", "consume":
			validActions = append(validActions, actionStr)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid action: %s (must be 'produce' or 'consume')", actionStr)), nil
		}
	}

	if len(validActions) == 0 {
		return mcp.NewToolResultError("No valid actions specified"), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Grant permissions
	err = grantTopicPermission(admin, *topicName, role, validActions)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to grant permissions: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully granted %s permissions on topic %s to role %s",
		strings.Join(validActions, ", "), topicName.String(), role)), nil
}

// handleTopicsRevokePermissions revokes all permissions for a role on a topic
func handleTopicsRevokePermissions(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	role, err := requiredParam[string](request.Params.Arguments, "role")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get role: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Revoke permissions
	err = admin.Topics().RevokePermission(*topicName, role)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to revoke permissions: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully revoked all permissions on topic %s from role %s",
		topicName.String(), role)), nil
}

func grantTopicPermission(admin interface{}, topicName utils.TopicName, role string, actions []string) error {
	authActions := make([]utils.AuthAction, len(actions))
	for i, action := range actions {
		authActions[i] = utils.AuthAction(action)
	}

	type TopicsProvider interface {
		Topics() interface {
			GrantPermission(utils.TopicName, string, []utils.AuthAction) error
		}
	}

	if provider, ok := admin.(TopicsProvider); ok {
		return provider.Topics().GrantPermission(topicName, role, authActions)
	}

	return fmt.Errorf("admin client does not implement the required interface")
}

// handleTopicsGetMessageTTL gets the message TTL for a topic
func handleTopicsGetMessageTTL(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Get message TTL
	ttl, err := admin.Topics().GetMessageTTL(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic message TTL: %v", err)), nil
	}

	// Format the result
	if ttl == 0 {
		return mcp.NewToolResultText("Message TTL is not set for topic " + topicName.String()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Message TTL for topic %s is %d seconds", topicName.String(), ttl)), nil
}

// handleTopicsSetMessageTTL sets the message TTL for a topic
func handleTopicsSetMessageTTL(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	ttl, err := requiredParam[float64](request.Params.Arguments, "ttl")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get TTL: %v", err)), nil
	}

	// Validate TTL
	if ttl < 0 {
		return mcp.NewToolResultError("TTL must be a non-negative value"), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Set message TTL
	err = admin.Topics().SetMessageTTL(*topicName, int(ttl))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set topic message TTL: %v", err)), nil
	}

	// Format the result based on whether TTL was enabled or disabled
	if ttl == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Message TTL disabled for topic %s", topicName.String())), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Message TTL for topic %s set to %d seconds", topicName.String(), int(ttl)),
	), nil
}

// handleTopicsRemoveMessageTTL removes the message TTL for a topic
func handleTopicsRemoveMessageTTL(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Remove message TTL
	err = admin.Topics().RemoveMessageTTL(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to remove topic message TTL: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Message TTL removed for topic %s", topicName.String())), nil
}

// handleTopicsGetMaxProducers gets the maximum number of producers allowed for a topic
func handleTopicsGetMaxProducers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Get max producers
	maxProducers, err := admin.Topics().GetMaxProducers(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic max producers: %v", err)), nil
	}

	// Format the result
	if maxProducers == 0 {
		return mcp.NewToolResultText("Maximum number of producers is unlimited for topic " + topicName.String()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Maximum number of producers allowed for topic %s is %d",
		topicName.String(), maxProducers)), nil
}

// handleTopicsSetMaxProducers sets the maximum number of producers allowed for a topic
func handleTopicsSetMaxProducers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	maxProducers, err := requiredParam[float64](request.Params.Arguments, "maxProducers")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get maxProducers: %v", err)), nil
	}

	// Validate maxProducers
	if maxProducers < 0 {
		return mcp.NewToolResultError("Maximum number of producers must be a non-negative value"), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Set max producers
	err = admin.Topics().SetMaxProducers(*topicName, int(maxProducers))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set topic max producers: %v", err)), nil
	}

	// Format the result based on whether the limit was set or removed
	if maxProducers == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Maximum number of producers set to unlimited for topic %s",
			topicName.String())), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Maximum number of producers allowed for topic %s set to %d",
		topicName.String(), int(maxProducers))), nil
}

// handleTopicsRemoveMaxProducers removes the maximum producers limit for a topic
func handleTopicsRemoveMaxProducers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Remove max producers limit
	err = admin.Topics().RemoveMaxProducers(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to remove topic max producers limit: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Maximum producers limit removed for topic %s", topicName.String())), nil
}

// handleTopicsGetMaxConsumers gets the maximum number of consumers allowed for a topic
func handleTopicsGetMaxConsumers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Get max consumers
	maxConsumers, err := admin.Topics().GetMaxConsumers(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic max consumers: %v", err)), nil
	}

	// Format the result
	if maxConsumers == 0 {
		return mcp.NewToolResultText("Maximum number of consumers is unlimited for topic " + topicName.String()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Maximum number of consumers allowed for topic %s is %d",
		topicName.String(), maxConsumers)), nil
}

// handleTopicsSetMaxConsumers sets the maximum number of consumers allowed for a topic
func handleTopicsSetMaxConsumers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	maxConsumers, err := requiredParam[float64](request.Params.Arguments, "maxConsumers")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get maxConsumers: %v", err)), nil
	}

	// Validate maxConsumers
	if maxConsumers < 0 {
		return mcp.NewToolResultError("Maximum number of consumers must be a non-negative value"), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Set max consumers
	err = admin.Topics().SetMaxConsumers(*topicName, int(maxConsumers))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set topic max consumers: %v", err)), nil
	}

	// Format the result based on whether the limit was set or removed
	if maxConsumers == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Maximum number of consumers set to unlimited for topic %s",
			topicName.String())), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Maximum number of consumers allowed for topic %s set to %d",
		topicName.String(), int(maxConsumers))), nil
}

// handleTopicsRemoveMaxConsumers removes the maximum consumers limit for a topic
func handleTopicsRemoveMaxConsumers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Remove max consumers limit
	err = admin.Topics().RemoveMaxConsumers(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to remove topic max consumers limit: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Maximum consumers limit removed for topic %s", topicName.String())), nil
}

// handleTopicsGetMaxUnackMessagesPerConsumer
// gets the maximum number of unacknowledged messages allowed for a consumer on a topic
func handleTopicsGetMaxUnackMessagesPerConsumer(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Get max unack messages per consumer
	maxUnack, err := admin.Topics().GetMaxUnackMessagesPerConsumer(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic max unack messages per consumer: %v", err)), nil
	}

	// Format the result
	if maxUnack == 0 {
		return mcp.NewToolResultText(
			"Maximum unacknowledged messages per consumer is unlimited for topic " + topicName.String(),
		), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Maximum unacknowledged messages per consumer for topic %s is %d",
		topicName.String(), maxUnack)), nil
}

// handleTopicsSetMaxUnackMessagesPerConsumer
// sets the maximum number of unacknowledged messages allowed for a consumer on a topic
func handleTopicsSetMaxUnackMessagesPerConsumer(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	maxUnack, err := requiredParam[float64](request.Params.Arguments, "maxUnackMessagesPerConsumer")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get maxUnackMessagesPerConsumer: %v", err)), nil
	}

	// Validate maxUnack
	if maxUnack < 0 {
		return mcp.NewToolResultError("Maximum unacknowledged messages per consumer must be a non-negative value"), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Set max unack messages per consumer
	err = admin.Topics().SetMaxUnackMessagesPerConsumer(*topicName, int(maxUnack))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set topic max unack messages per consumer: %v", err)), nil
	}

	// Format the result based on whether the limit was set or unlimited
	if maxUnack == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Maximum unacknowledged messages per consumer set to unlimited for topic %s",
			topicName.String())), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Maximum unacknowledged messages per consumer for topic %s set to %d",
		topicName.String(), int(maxUnack))), nil
}

// handleTopicsRemoveMaxUnackMessagesPerConsumer
// removes the maximum unacknowledged messages per consumer limit for a topic
func handleTopicsRemoveMaxUnackMessagesPerConsumer(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Remove max unack messages per consumer limit
	err = admin.Topics().RemoveMaxUnackMessagesPerConsumer(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf(
			"Failed to remove topic max unacknowledged messages per consumer limit: %v", err),
		), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf(
		"Maximum unacknowledged messages per consumer limit removed for topic %s",
		topicName.String(),
	)), nil
}

// handleTopicsGetMaxUnackMessagesPerSubscription
// gets the maximum number of unacknowledged messages allowed for a subscription on a topic
func handleTopicsGetMaxUnackMessagesPerSubscription(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Get max unack messages per subscription
	maxUnack, err := admin.Topics().GetMaxUnackMessagesPerSubscription(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic max unack messages per subscription: %v", err)), nil
	}

	// Format the result
	if maxUnack == 0 {
		return mcp.NewToolResultText(
			"Maximum unacknowledged messages per subscription is unlimited for topic " + topicName.String(),
		), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf(
		"Maximum unacknowledged messages per subscription for topic %s is %d",
		topicName.String(), maxUnack,
	)), nil
}

// handleTopicsSetMaxUnackMessagesPerSubscription
// sets the maximum number of unacknowledged messages allowed for a subscription on a topic
func handleTopicsSetMaxUnackMessagesPerSubscription(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	maxUnack, err := requiredParam[float64](request.Params.Arguments, "maxUnackMessagesPerSubscription")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get maxUnackMessagesPerSubscription: %v", err)), nil
	}

	// Validate maxUnack
	if maxUnack < 0 {
		return mcp.NewToolResultError("Maximum unacknowledged messages per subscription must be a non-negative value"), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Set max unack messages per subscription
	err = admin.Topics().SetMaxUnackMessagesPerSubscription(*topicName, int(maxUnack))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set topic max unack messages per subscription: %v", err)), nil
	}

	// Format the result based on whether the limit was set or unlimited
	if maxUnack == 0 {
		return mcp.NewToolResultText(fmt.Sprintf(
			"Maximum unacknowledged messages per subscription set to unlimited for topic %s",
			topicName.String(),
		)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf(
		"Maximum unacknowledged messages per subscription for topic %s set to %d",
		topicName.String(), int(maxUnack),
	)), nil
}

// handleTopicsRemoveMaxUnackMessagesPerSubscription
// removes the maximum unacknowledged messages per subscription limit for a topic
func handleTopicsRemoveMaxUnackMessagesPerSubscription(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Remove max unack messages per subscription limit
	err = admin.Topics().RemoveMaxUnackMessagesPerSubscription(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf(
			"Failed to remove topic max unacknowledged messages per subscription limit: %v", err),
		), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf(
		"Maximum unacknowledged messages per subscription limit removed for topic %s",
		topicName.String(),
	)), nil
}

// handleTopicsGetPersistence gets the persistence policy for a topic
func handleTopicsGetPersistence(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Get persistence policy
	persistence, err := admin.Topics().GetPersistence(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic persistence policy: %v", err)), nil
	}

	// Format the output
	jsonBytes, err := json.Marshal(persistence)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal topic persistence policy: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleTopicsSetPersistence sets the persistence policy for a topic
func handleTopicsSetPersistence(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	ensembleSize, err := requiredParam[float64](request.Params.Arguments, "ensembleSize")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get ensembleSize: %v", err)), nil
	}

	writeQuorum, err := requiredParam[float64](request.Params.Arguments, "writeQuorum")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get writeQuorum: %v", err)), nil
	}

	ackQuorum, err := requiredParam[float64](request.Params.Arguments, "ackQuorum")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get ackQuorum: %v", err)), nil
	}

	// Validate values
	if ensembleSize <= 0 || writeQuorum <= 0 || ackQuorum <= 0 {
		return mcp.NewToolResultError("All persistence values must be positive"), nil
	}

	// Validate quorum relationships
	if writeQuorum > ensembleSize {
		return mcp.NewToolResultError("Write quorum cannot be greater than ensemble size"), nil
	}

	if ackQuorum > writeQuorum {
		return mcp.NewToolResultError("Ack quorum cannot be greater than write quorum"), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Create persistence policy as PersistenceData (required by SetPersistence)
	persistence := utils.PersistenceData{
		BookkeeperEnsemble:             int64(ensembleSize),
		BookkeeperWriteQuorum:          int64(writeQuorum),
		BookkeeperAckQuorum:            int64(ackQuorum),
		ManagedLedgerMaxMarkDeleteRate: 0.0,
	}

	// Set persistence policy
	err = admin.Topics().SetPersistence(*topicName, persistence)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set topic persistence policy: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Set persistence policy for topic %s with ensemble size: %d, write quorum: %d, ack quorum: %d",
			topicName.String(), int(ensembleSize), int(writeQuorum), int(ackQuorum))), nil
}

// handleTopicsRemovePersistence removes the persistence policy for a topic
func handleTopicsRemovePersistence(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Remove persistence policy
	err = admin.Topics().RemovePersistence(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to remove topic persistence policy: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Persistence policy removed for topic %s", topicName.String())), nil
}

// handleTopicsGetDelayedDelivery gets the delayed delivery policy for a topic
func handleTopicsGetDelayedDelivery(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Get delayed delivery policy
	delayedDelivery, err := admin.Topics().GetDelayedDelivery(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic delayed delivery policy: %v", err)), nil
	}

	// Format the output
	if delayedDelivery.Active {
		return mcp.NewToolResultText(fmt.Sprintf("Delayed delivery is enabled for topic %s with tick time of %d ms",
			topicName.String(), int(delayedDelivery.TickTime))), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Delayed delivery is disabled for topic %s", topicName.String())), nil
}

// handleTopicsSetDelayedDelivery sets the delayed delivery policy for a topic
func handleTopicsSetDelayedDelivery(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	delayInMillis, err := requiredParam[float64](request.Params.Arguments, "delayInMillis")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get delayInMillis: %v", err)), nil
	}

	// Validate delay
	if delayInMillis < 0 {
		return mcp.NewToolResultError("Delay time must be a non-negative value"), nil
	}

	// Default tick time is 1 second (1000ms)
	tickTime := 1000.0
	tickTimeParam, hasTickTime := optionalParam[float64](request.Params.Arguments, "tickTime")
	if hasTickTime && tickTimeParam > 0 {
		tickTime = tickTimeParam
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Create delayed delivery policy
	delayedDelivery := utils.DelayedDeliveryData{
		TickTime: tickTime,
		Active:   delayInMillis > 0,
	}

	// Set delayed delivery policy
	err = admin.Topics().SetDelayedDelivery(*topicName, delayedDelivery)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set topic delayed delivery policy: %v", err)), nil
	}

	// Format the result based on whether delayed delivery was enabled or disabled
	if delayInMillis == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Delayed delivery disabled for topic %s", topicName.String())), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Delayed delivery enabled for topic %s with tick time of %d ms",
		topicName.String(), int(tickTime))), nil
}

// handleTopicsRemoveDelayedDelivery removes the delayed delivery policy for a topic
func handleTopicsRemoveDelayedDelivery(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Remove delayed delivery policy
	err = admin.Topics().RemoveDelayedDelivery(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to remove topic delayed delivery policy: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Delayed delivery policy removed for topic %s", topicName.String())), nil
}

// handleTopicsGetDispatchRate gets the message dispatch rate for a topic
func handleTopicsGetDispatchRate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Get dispatch rate
	dispatchRate, err := admin.Topics().GetDispatchRate(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic dispatch rate: %v", err)), nil
	}

	// Format the output
	jsonBytes, err := json.Marshal(dispatchRate)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal topic dispatch rate: %v", err)), nil
	}

	// Check if dispatch rate is disabled/unlimited
	if dispatchRate.DispatchThrottlingRateInMsg < 0 && dispatchRate.DispatchThrottlingRateInByte < 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Dispatch rate is unlimited for topic %s", topicName.String())), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleTopicsSetDispatchRate sets the message dispatch rate for a topic
func handleTopicsSetDispatchRate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Default values for dispatch rate
	dispatchRateInMsg := int64(-1)  // unlimited
	dispatchRateInByte := int64(-1) // unlimited
	ratePeriodInSecond := int64(1)  // default 1 second

	// Get dispatch rate in messages if provided
	msgRateParam, hasMsgRate := optionalParam[float64](request.Params.Arguments, "dispatchThrottlingRateInMsg")
	if hasMsgRate {
		dispatchRateInMsg = int64(msgRateParam)
	}

	// Get dispatch rate in bytes if provided
	byteRateParam, hasByteRate := optionalParam[float64](request.Params.Arguments, "dispatchThrottlingRateInByte")
	if hasByteRate {
		dispatchRateInByte = int64(byteRateParam)
	}

	// Get rate period if provided
	ratePeriodParam, hasRatePeriod := optionalParam[float64](request.Params.Arguments, "ratePeriodInSecond")
	if hasRatePeriod && ratePeriodParam > 0 {
		ratePeriodInSecond = int64(ratePeriodParam)
	} else if hasRatePeriod && ratePeriodParam <= 0 {
		return mcp.NewToolResultError("Rate period must be positive"), nil
	}

	// Create dispatch rate policy
	dispatchRate := utils.DispatchRateData{
		DispatchThrottlingRateInMsg:  dispatchRateInMsg,
		DispatchThrottlingRateInByte: dispatchRateInByte,
		RatePeriodInSecond:           ratePeriodInSecond,
		RelativeToPublishRate:        false,
	}

	// Set dispatch rate
	err = admin.Topics().SetDispatchRate(*topicName, dispatchRate)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set topic dispatch rate: %v", err)), nil
	}

	// Format the response message
	if dispatchRateInMsg < 0 && dispatchRateInByte < 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Dispatch rate set to unlimited for topic %s", topicName.String())), nil
	}

	var msgText string
	if dispatchRateInMsg < 0 {
		msgText = "unlimited messages"
	} else {
		msgText = fmt.Sprintf("%d messages", dispatchRateInMsg)
	}

	var byteText string
	if dispatchRateInByte < 0 {
		byteText = "unlimited bytes"
	} else {
		byteText = fmt.Sprintf("%d bytes", dispatchRateInByte)
	}

	return mcp.NewToolResultText(fmt.Sprintf("Dispatch rate set for topic %s: %s and %s per %d second(s)",
		topicName.String(), msgText, byteText, ratePeriodInSecond)), nil
}

// handleTopicsRemoveDispatchRate removes the message dispatch rate for a topic
func handleTopicsRemoveDispatchRate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Remove dispatch rate
	err = admin.Topics().RemoveDispatchRate(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to remove topic dispatch rate: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Dispatch rate removed for topic %s", topicName.String())), nil
}

// handleTopicsGetDeduplicationStatus gets the deduplication status for a topic
func handleTopicsGetDeduplicationStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Get deduplication status
	enabled, err := admin.Topics().GetDeduplicationStatus(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic deduplication status: %v", err)), nil
	}

	// Format the result
	if enabled {
		return mcp.NewToolResultText(fmt.Sprintf("Deduplication is enabled for topic %s", topicName.String())), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Deduplication is disabled for topic %s", topicName.String())), nil
}

// handleTopicsSetDeduplicationStatus sets the deduplication status for a topic
func handleTopicsSetDeduplicationStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	enabled, err := requiredParam[bool](request.Params.Arguments, "enabled")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get enabled parameter: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Set deduplication status
	err = admin.Topics().SetDeduplicationStatus(*topicName, enabled)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set topic deduplication status: %v", err)), nil
	}

	// Format the result
	if enabled {
		return mcp.NewToolResultText(fmt.Sprintf("Deduplication enabled for topic %s", topicName.String())), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Deduplication disabled for topic %s", topicName.String())), nil
}

// handleTopicsRemoveDeduplicationStatus removes the deduplication status for a topic
func handleTopicsRemoveDeduplicationStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Remove deduplication status
	err = admin.Topics().RemoveDeduplicationStatus(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to remove topic deduplication status: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Deduplication status removed for topic %s", topicName.String())), nil
}

// handleTopicsGetRetention gets the retention policy for a topic
func handleTopicsGetRetention(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Check if applied policies should be included
	applied := false
	appliedParam, hasApplied := optionalParam[bool](request.Params.Arguments, "applied")
	if hasApplied {
		applied = appliedParam
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Get retention policy
	retention, err := admin.Topics().GetRetention(*topicName, applied)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic retention policy: %v", err)), nil
	}

	// If no retention policy is defined
	if retention == nil {
		if applied {
			return mcp.NewToolResultText(fmt.Sprintf("No retention policy found for topic %s (including namespace level)",
				topicName.String())), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("No retention policy found for topic %s", topicName.String())), nil
	}

	// Format the output
	var retentionTime string
	if retention.RetentionTimeInMinutes <= 0 {
		retentionTime = "infinite time"
	} else {
		retentionTime = fmt.Sprintf("%d minutes", retention.RetentionTimeInMinutes)
	}

	var retentionSize string
	if retention.RetentionSizeInMB <= 0 {
		retentionSize = "infinite size"
	} else {
		retentionSize = fmt.Sprintf("%d MB", retention.RetentionSizeInMB)
	}

	var source string
	if applied {
		source = " (may include namespace level policy)"
	}

	return mcp.NewToolResultText(fmt.Sprintf("Retention policy for topic %s%s: %s and %s",
		topicName.String(), source, retentionTime, retentionSize)), nil
}

// handleTopicsSetRetention sets the retention policy for a topic
func handleTopicsSetRetention(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	retentionTimeInMinutes, err := requiredParam[float64](request.Params.Arguments, "retentionTimeInMinutes")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get retentionTimeInMinutes: %v", err)), nil
	}

	retentionSizeInMB, err := requiredParam[float64](request.Params.Arguments, "retentionSizeInMB")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get retentionSizeInMB: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Create retention policy
	retention := utils.RetentionPolicies{
		RetentionTimeInMinutes: int(retentionTimeInMinutes),
		RetentionSizeInMB:      int64(retentionSizeInMB),
	}

	// Set retention policy
	err = admin.Topics().SetRetention(*topicName, retention)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set topic retention policy: %v", err)), nil
	}

	// Format the result
	var timeText, sizeText string

	//nolint:gocritic
	if retention.RetentionTimeInMinutes < 0 {
		timeText = "disabled time retention"
	} else if retention.RetentionTimeInMinutes == 0 {
		timeText = "infinite time retention"
	} else {
		timeText = fmt.Sprintf("%d minutes time retention", retention.RetentionTimeInMinutes)
	}

	//nolint:gocritic
	if retention.RetentionSizeInMB < 0 {
		sizeText = "disabled size retention"
	} else if retention.RetentionSizeInMB == 0 {
		sizeText = "infinite size retention"
	} else {
		sizeText = fmt.Sprintf("%d MB size retention", retention.RetentionSizeInMB)
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set retention policy for topic %s with %s and %s",
		topicName.String(), timeText, sizeText)), nil
}

// handleTopicsRemoveRetention removes the retention policy for a topic
func handleTopicsRemoveRetention(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Remove retention policy
	err = admin.Topics().RemoveRetention(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to remove topic retention policy: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Retention policy removed for topic %s", topicName.String())), nil
}

// handleTopicsGetBacklogQuota gets the backlog quota policy for a topic
func handleTopicsGetBacklogQuota(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Check if applied policies should be included
	applied := false
	appliedParam, hasApplied := optionalParam[bool](request.Params.Arguments, "applied")
	if hasApplied {
		applied = appliedParam
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Get backlog quota
	backlogQuotaMap, err := admin.Topics().GetBacklogQuotaMap(*topicName, applied)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic backlog quota: %v", err)), nil
	}

	// Format the output
	jsonBytes, err := json.Marshal(backlogQuotaMap)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal topic backlog quota: %v", err)), nil
	}

	if len(backlogQuotaMap) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No backlog quota policy found for topic %s", topicName.String())), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleTopicsSetBacklogQuota sets the backlog quota policy for a topic
func handleTopicsSetBacklogQuota(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	limitSize, err := requiredParam[float64](request.Params.Arguments, "limitSize")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get limitSize: %v", err)), nil
	}

	policy, err := requiredParam[string](request.Params.Arguments, "policy")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get policy: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Default values for optional parameters
	limitTime := int64(-1) // unlimited by default

	// Get limit time if provided
	limitTimeParam, hasLimitTime := optionalParam[float64](request.Params.Arguments, "limitTime")
	if hasLimitTime {
		limitTime = int64(limitTimeParam)
	}

	// Parse retention policy
	retentionPolicy, err := utils.ParseRetentionPolicy(policy)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to parse policy: %v", err)), nil
	}

	// Default backlog quota type is destination_storage
	backlogQuotaType := utils.DestinationStorage

	// Get quota type if provided
	quotaTypeStr, hasQuotaType := optionalParam[string](request.Params.Arguments, "type")
	if hasQuotaType {
		parsedType, err := utils.ParseBacklogQuotaType(quotaTypeStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to parse quota type: %v", err)), nil
		}
		backlogQuotaType = parsedType
	}

	// Create backlog quota
	backlogQuota := utils.NewBacklogQuota(int64(limitSize), limitTime, retentionPolicy)

	// Set backlog quota
	err = admin.Topics().SetBacklogQuota(*topicName, backlogQuota, backlogQuotaType)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set topic backlog quota: %v", err)), nil
	}

	// Format the result
	return mcp.NewToolResultText(fmt.Sprintf("Successfully set backlog quota for topic %s", topicName.String())), nil
}

// handleTopicsRemoveBacklogQuota removes the backlog quota policy from a topic
func handleTopicsRemoveBacklogQuota(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Default backlog quota type is destination_storage
	backlogQuotaType := utils.DestinationStorage

	// Get quota type if provided
	quotaTypeStr, hasQuotaType := optionalParam[string](request.Params.Arguments, "type")
	if hasQuotaType {
		parsedType, err := utils.ParseBacklogQuotaType(quotaTypeStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to parse quota type: %v", err)), nil
		}
		backlogQuotaType = parsedType
	}

	// Remove backlog quota
	err = admin.Topics().RemoveBacklogQuota(*topicName, backlogQuotaType)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to remove topic backlog quota: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully removed backlog quota from topic %s", topicName.String())), nil
}

// handleTopicsGetCompactionThreshold gets the compaction threshold for a topic
func handleTopicsGetCompactionThreshold(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Check if applied policies should be included
	applied := false
	appliedParam, hasApplied := optionalParam[bool](request.Params.Arguments, "applied")
	if hasApplied {
		applied = appliedParam
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Get compaction threshold
	threshold, err := admin.Topics().GetCompactionThreshold(*topicName, applied)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic compaction threshold: %v", err)), nil
	}

	// Format the result
	if threshold == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Automatic compaction is disabled for topic %s", topicName.String())), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("The compaction threshold of the topic %s is %d byte(s)",
		topicName.String(), threshold)), nil
}

// handleTopicsSetCompactionThreshold sets the compaction threshold for a topic
func handleTopicsSetCompactionThreshold(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	threshold, err := requiredParam[float64](request.Params.Arguments, "threshold")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get threshold: %v", err)), nil
	}

	// Validate threshold
	if threshold < 0 {
		return mcp.NewToolResultError("Threshold must be a non-negative value"), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Set compaction threshold
	err = admin.Topics().SetCompactionThreshold(*topicName, int64(threshold))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set topic compaction threshold: %v", err)), nil
	}

	// Format the result based on whether automatic compaction was enabled or disabled
	if threshold == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Successfully disabled automatic compaction for topic %s",
			topicName.String())), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully set compaction threshold to %d byte(s) for topic %s",
		int64(threshold), topicName.String())), nil
}

// handleTopicsRemoveCompactionThreshold removes the compaction threshold for a topic
func handleTopicsRemoveCompactionThreshold(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Remove compaction threshold
	err = admin.Topics().RemoveCompactionThreshold(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to remove topic compaction threshold: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully removed compaction threshold for topic %s",
		topicName.String())), nil
}

// handleTopicsGetInactiveTopic gets the inactive topic policies for a topic
func handleTopicsGetInactiveTopic(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Check if applied policies should be included
	applied := false
	appliedParam, hasApplied := optionalParam[bool](request.Params.Arguments, "applied")
	if hasApplied {
		applied = appliedParam
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Get inactive topic policies
	policies, err := admin.Topics().GetInactiveTopicPolicies(*topicName, applied)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic inactive policies: %v", err)), nil
	}

	// Format the output
	jsonBytes, err := json.Marshal(policies)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal topic inactive policies: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleTopicsSetInactiveTopic sets the inactive topic policies for a topic
func handleTopicsSetInactiveTopic(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	enableDelete, err := requiredParam[bool](request.Params.Arguments, "enableDeleteWhileInactive")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get enableDeleteWhileInactive: %v", err)), nil
	}

	maxInactiveDuration, err := requiredParam[float64](request.Params.Arguments, "maxInactiveDurationSeconds")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get maxInactiveDurationSeconds: %v", err)), nil
	}

	deleteModeStr, err := requiredParam[string](request.Params.Arguments, "deleteMode")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get deleteMode: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Validate maxInactiveDuration
	if maxInactiveDuration <= 0 {
		return mcp.NewToolResultError("Max inactive duration must be positive"), nil
	}

	// Parse delete mode
	deleteMode, err := utils.ParseInactiveTopicDeleteMode(deleteModeStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid delete mode: %v", err)), nil
	}

	// Create inactive topic policies
	inactiveTopicPolicies := utils.NewInactiveTopicPolicies(
		&deleteMode,
		int(maxInactiveDuration),
		enableDelete,
	)

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Set inactive topic policies
	err = admin.Topics().SetInactiveTopicPolicies(*topicName, inactiveTopicPolicies)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set topic inactive policies: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully set inactive topic policies for topic %s",
		topicName.String())), nil
}

// handleTopicsRemoveInactiveTopic removes the inactive topic policies from a topic
func handleTopicsRemoveInactiveTopic(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name: %v", err)), nil
	}

	// Create the admin client
	admin, err := pulsar.GetAdminClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
	}

	// Remove inactive topic policies
	err = admin.Topics().RemoveInactiveTopicPolicies(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to remove topic inactive policies: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully removed inactive topic policies from topic %s",
		topicName.String())), nil
}

// handleTopicGetPolicy handles getting policies for a topic using the unified tool
func handleTopicGetPolicy(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	_, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic name: %v", err)), nil
	}

	policyType, err := requiredParam[string](request.Params.Arguments, "policy")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get policy type: %v", err)), nil
	}

	// Route to the appropriate handler based on the policy type
	switch policyType {
	case "publish-rate":
		return handleTopicsGetPublishRate(ctx, request)
	case "permissions":
		return handleTopicsGetPermissions(ctx, request)
	case "message-ttl":
		return handleTopicsGetMessageTTL(ctx, request)
	case "max-producers":
		return handleTopicsGetMaxProducers(ctx, request)
	case "max-consumers":
		return handleTopicsGetMaxConsumers(ctx, request)
	case "max-unack-messages-per-consumer":
		return handleTopicsGetMaxUnackMessagesPerConsumer(ctx, request)
	case "max-unack-messages-per-subscription":
		return handleTopicsGetMaxUnackMessagesPerSubscription(ctx, request)
	case "persistence":
		return handleTopicsGetPersistence(ctx, request)
	case "delayed-delivery":
		return handleTopicsGetDelayedDelivery(ctx, request)
	case "dispatch-rate":
		return handleTopicsGetDispatchRate(ctx, request)
	case "deduplication":
		return handleTopicsGetDeduplicationStatus(ctx, request)
	case "retention":
		return handleTopicsGetRetention(ctx, request)
	case "backlog-quota":
		return handleTopicsGetBacklogQuota(ctx, request)
	case "compaction-threshold":
		return handleTopicsGetCompactionThreshold(ctx, request)
	case "inactive-topic":
		return handleTopicsGetInactiveTopic(ctx, request)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unknown policy type: %s", policyType)), nil
	}
}

// handleTopicSetPolicy handles setting policies for a topic using the unified tool
func handleTopicSetPolicy(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	_, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic name: %v", err)), nil
	}

	policyType, err := requiredParam[string](request.Params.Arguments, "policy")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get policy type: %v", err)), nil
	}

	// Route to the appropriate handler based on the policy type
	switch policyType {
	case "publish-rate":
		return handleTopicsSetPublishRate(ctx, request)
	case "permissions":
		return handleTopicsGrantPermissions(ctx, request)
	case "message-ttl":
		return handleTopicsSetMessageTTL(ctx, request)
	case "max-producers":
		return handleTopicsSetMaxProducers(ctx, request)
	case "max-consumers":
		return handleTopicsSetMaxConsumers(ctx, request)
	case "max-unack-messages-per-consumer":
		return handleTopicsSetMaxUnackMessagesPerConsumer(ctx, request)
	case "max-unack-messages-per-subscription":
		return handleTopicsSetMaxUnackMessagesPerSubscription(ctx, request)
	case "persistence":
		return handleTopicsSetPersistence(ctx, request)
	case "delayed-delivery":
		return handleTopicsSetDelayedDelivery(ctx, request)
	case "dispatch-rate":
		return handleTopicsSetDispatchRate(ctx, request)
	case "deduplication":
		return handleTopicsSetDeduplicationStatus(ctx, request)
	case "retention":
		return handleTopicsSetRetention(ctx, request)
	case "backlog-quota":
		return handleTopicsSetBacklogQuota(ctx, request)
	case "compaction-threshold":
		return handleTopicsSetCompactionThreshold(ctx, request)
	case "inactive-topic":
		return handleTopicsSetInactiveTopic(ctx, request)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unknown policy type: %s", policyType)), nil
	}
}

// handleTopicRemovePolicy handles removing policies for a topic using the unified tool
func handleTopicRemovePolicy(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	_, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic name: %v", err)), nil
	}

	policyType, err := requiredParam[string](request.Params.Arguments, "policy")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get policy type: %v", err)), nil
	}

	// Route to the appropriate handler based on the policy type
	switch policyType {
	case "publish-rate":
		return handleTopicsRemovePublishRate(ctx, request)
	case "permissions":
		return handleTopicsRevokePermissions(ctx, request)
	case "message-ttl":
		return handleTopicsRemoveMessageTTL(ctx, request)
	case "max-producers":
		return handleTopicsRemoveMaxProducers(ctx, request)
	case "max-consumers":
		return handleTopicsRemoveMaxConsumers(ctx, request)
	case "max-unack-messages-per-consumer":
		return handleTopicsRemoveMaxUnackMessagesPerConsumer(ctx, request)
	case "max-unack-messages-per-subscription":
		return handleTopicsRemoveMaxUnackMessagesPerSubscription(ctx, request)
	case "persistence":
		return handleTopicsRemovePersistence(ctx, request)
	case "delayed-delivery":
		return handleTopicsRemoveDelayedDelivery(ctx, request)
	case "dispatch-rate":
		return handleTopicsRemoveDispatchRate(ctx, request)
	case "deduplication":
		return handleTopicsRemoveDeduplicationStatus(ctx, request)
	case "retention":
		return handleTopicsRemoveRetention(ctx, request)
	case "backlog-quota":
		return handleTopicsRemoveBacklogQuota(ctx, request)
	case "compaction-threshold":
		return handleTopicsRemoveCompactionThreshold(ctx, request)
	case "inactive-topic":
		return handleTopicsRemoveInactiveTopic(ctx, request)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unknown policy type: %s", policyType)), nil
	}
}
