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
	"math"
	"strconv"
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	ctlutil "github.com/streamnative/pulsarctl/pkg/ctl/utils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

var readOnlyRestrictedTopicPolicyOperations = map[string]struct{}{
	"set-retention":                                {},
	"remove-retention":                             {},
	"set-message-ttl":                              {},
	"remove-message-ttl":                           {},
	"set-max-producers":                            {},
	"remove-max-producers":                         {},
	"set-max-consumers":                            {},
	"remove-max-consumers":                         {},
	"set-max-unacked-messages-per-consumer":        {},
	"remove-max-unacked-messages-per-consumer":     {},
	"set-max-unacked-messages-per-subscription":    {},
	"remove-max-unacked-messages-per-subscription": {},
	"set-persistence":                              {},
	"remove-persistence":                           {},
	"set-delayed-delivery":                         {},
	"remove-delayed-delivery":                      {},
	"set-dispatch-rate":                            {},
	"remove-dispatch-rate":                         {},
	"set-subscription-dispatch-rate":               {},
	"remove-subscription-dispatch-rate":            {},
	"set-deduplication":                            {},
	"remove-deduplication":                         {},
	"set-backlog-quota":                            {},
	"remove-backlog-quota":                         {},
	"set-compaction-threshold":                     {},
	"remove-compaction-threshold":                  {},
	"set-publish-rate":                             {},
	"remove-publish-rate":                          {},
	"set-inactive-topic-policies":                  {},
	"remove-inactive-topic-policies":               {},
	"set-subscription-types":                       {},
	"remove-subscription-types":                    {},
}

var topicPolicyOperationAliases = map[string]string{
	"get_ttl":                     "get-message-ttl",
	"set_ttl":                     "set-message-ttl",
	"remove_ttl":                  "remove-message-ttl",
	"get_compaction":              "get-compaction-threshold",
	"set_compaction":              "set-compaction-threshold",
	"remove_compaction":           "remove-compaction-threshold",
	"get-deduplication-status":    "get-deduplication",
	"set-deduplication-status":    "set-deduplication",
	"remove-deduplication-status": "remove-deduplication",
	"get-backlog-quota":           "get-backlog-quotas",
	"get_backlog_quota":           "get-backlog-quotas",
	"get-inactive-topic":          "get-inactive-topic-policies",
	"set-inactive-topic":          "set-inactive-topic-policies",
	"remove-inactive-topic":       "remove-inactive-topic-policies",
}

// PulsarAdminTopicPolicyToolBuilder implements the ToolBuilder interface for Pulsar admin topic policies
// /nolint:revive
type PulsarAdminTopicPolicyToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminTopicPolicyToolBuilder creates a new Pulsar admin topic policy tool builder instance
func NewPulsarAdminTopicPolicyToolBuilder() *PulsarAdminTopicPolicyToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_topic_policy",
		Version:     "1.0.0",
		Description: "Pulsar admin topic policy management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "admin", "topic_policy"},
	}

	features := []string{
		"pulsar-admin-topic-policy",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminTopicPolicyToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin topic policy tool list
func (b *PulsarAdminTopicPolicyToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	if err := b.Validate(config); err != nil {
		return nil, err
	}

	return []server.ServerTool{
		{
			Tool:    b.buildTopicPolicyTool(),
			Handler: b.buildTopicPolicyHandler(config.ReadOnly),
		},
	}, nil
}

// buildTopicPolicyTool builds the Pulsar Admin Topic Policy MCP tool definition
func (b *PulsarAdminTopicPolicyToolBuilder) buildTopicPolicyTool() mcp.Tool {
	toolDesc := "Manage Pulsar topic-level policies with operation names aligned to pulsarctl topic policy commands. " +
		"This tool covers retention, message TTL, producer and consumer limits, persistence, delayed delivery, " +
		"dispatch throttling, deduplication, backlog quotas, compaction thresholds, publish rates, inactive topic policies, " +
		"and subscription dispatch throttling. Legacy underscore operation aliases from the older MCP implementation remain supported."

	operationDesc := strings.Join([]string{
		"Operation to perform. Available operations:",
		"- get-retention / set-retention / remove-retention: manage topic retention policy",
		"- get-message-ttl / set-message-ttl / remove-message-ttl: manage topic message TTL",
		"- get-max-producers / set-max-producers / remove-max-producers: manage producer limit",
		"- get-max-consumers / set-max-consumers / remove-max-consumers: manage consumer limit",
		"- get-max-unacked-messages-per-consumer / set-max-unacked-messages-per-consumer / remove-max-unacked-messages-per-consumer",
		"- get-max-unacked-messages-per-subscription / set-max-unacked-messages-per-subscription / remove-max-unacked-messages-per-subscription",
		"- get-persistence / set-persistence / remove-persistence: manage topic persistence",
		"- get-delayed-delivery / set-delayed-delivery / remove-delayed-delivery: manage delayed delivery policy",
		"- get-dispatch-rate / set-dispatch-rate / remove-dispatch-rate: manage topic dispatch throttling",
		"- get-subscription-dispatch-rate / set-subscription-dispatch-rate / remove-subscription-dispatch-rate",
		"- get-deduplication / set-deduplication / remove-deduplication: manage deduplication policy",
		"- get-backlog-quotas / set-backlog-quota / remove-backlog-quota: manage backlog quotas",
		"- get-compaction-threshold / set-compaction-threshold / remove-compaction-threshold: manage compaction threshold",
		"- get-publish-rate / set-publish-rate / remove-publish-rate: manage publish rate policy",
		"- get-inactive-topic-policies / set-inactive-topic-policies / remove-inactive-topic-policies",
		"- get-subscription-types / set-subscription-types / remove-subscription-types: additional MCP-only compatibility operations",
	}, "\n")

	return mcp.NewTool("pulsar_admin_topic_policy",
		mcp.WithDescription(toolDesc),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
		),
		mcp.WithString("topic", mcp.Required(),
			mcp.Description("Topic name in the format 'persistent://tenant/namespace/topic' or 'tenant/namespace/topic'."),
		),
		mcp.WithBoolean("applied",
			mcp.Description("Return the effective policy inherited from namespace defaults when supported. Used by get-retention, get-backlog-quotas, get-compaction-threshold, and get-inactive-topic-policies."),
		),
		mcp.WithString("retention-time",
			mcp.Description("Retention time string like '10m', '24h', or '7d'. Used by set-retention."),
		),
		mcp.WithString("retention-size",
			mcp.Description("Retention size string like '100M', '5G', or '1T'. Used by set-retention."),
		),
		mcp.WithNumber("ttl-seconds",
			mcp.Description("Message TTL in seconds. Used by set-message-ttl."),
		),
		mcp.WithNumber("count",
			mcp.Description("Generic count for set-max-* operations. Legacy pulsarctl-style aliases such as max-producers, max-consumers, and maxNum are also accepted."),
		),
		mcp.WithString("threshold",
			mcp.Description("Compaction threshold size string like '100M' or '1G'. Used by set-compaction-threshold."),
		),
		mcp.WithNumber("bookkeeper-ensemble",
			mcp.Description("BookKeeper ensemble size. Used by set-persistence."),
		),
		mcp.WithNumber("bookkeeper-write-quorum",
			mcp.Description("BookKeeper write quorum size. Used by set-persistence."),
		),
		mcp.WithNumber("bookkeeper-ack-quorum",
			mcp.Description("BookKeeper ack quorum size. Used by set-persistence."),
		),
		mcp.WithNumber("ml-mark-delete-max-rate",
			mcp.Description("Managed ledger max mark delete rate. Used by set-persistence."),
		),
		mcp.WithBoolean("enable",
			mcp.Description("Enable a toggle-style policy. Used by set-deduplication and set-delayed-delivery."),
		),
		mcp.WithBoolean("disable",
			mcp.Description("Disable a toggle-style policy. Used by set-deduplication and set-delayed-delivery."),
		),
		mcp.WithString("time",
			mcp.Description("Relative time string like '10s', '1m', or '1h'. Used by set-delayed-delivery."),
		),
		mcp.WithNumber("msg-rate",
			mcp.Description("Message rate for dispatch or publish throttling. Used by set-dispatch-rate, set-subscription-dispatch-rate, and set-publish-rate."),
		),
		mcp.WithNumber("byte-rate",
			mcp.Description("Byte rate for dispatch or publish throttling. Used by set-dispatch-rate, set-subscription-dispatch-rate, and set-publish-rate."),
		),
		mcp.WithNumber("period",
			mcp.Description("Rate period in seconds for dispatch throttling. Used by set-dispatch-rate and set-subscription-dispatch-rate."),
		),
		mcp.WithBoolean("relative-to-publish-rate",
			mcp.Description("Whether dispatch throttling should be applied relative to publish rate. Used by set-dispatch-rate and set-subscription-dispatch-rate."),
		),
		mcp.WithString("limit-size",
			mcp.Description("Backlog quota size limit like '16G'. Used by set-backlog-quota."),
		),
		mcp.WithNumber("limit-time",
			mcp.Description("Backlog quota time limit in seconds. Used by set-backlog-quota."),
		),
		mcp.WithString("policy",
			mcp.Description("Backlog quota retention policy. Valid values: producer_request_hold, producer_exception, consumer_backlog_eviction. Used by set-backlog-quota."),
		),
		mcp.WithString("type",
			mcp.Description("Backlog quota type. Valid values: destination_storage or message_age. Used by get-backlog-quotas, set-backlog-quota, and remove-backlog-quota."),
		),
		mcp.WithBoolean("delete-while-inactive",
			mcp.Description("Whether inactive topics should be deleted. Used by set-inactive-topic-policies."),
		),
		mcp.WithString("max-inactive-duration",
			mcp.Description("Relative duration like '1h' or '7d'. Used by set-inactive-topic-policies."),
		),
		mcp.WithString("delete-mode",
			mcp.Description("Inactive topic delete mode. Valid values: delete_when_no_subscriptions or delete_when_subscriptions_caught_up. Used by set-inactive-topic-policies."),
		),
		mcp.WithArray("subscription-types",
			mcp.Description("Allowed subscription types. Used by set-subscription-types. Valid values: Exclusive, Shared, Failover, Key_Shared."),
			mcp.Items(
				map[string]any{
					"type":        "string",
					"description": "subscription type",
				},
			),
		),
	)
}

// buildTopicPolicyHandler builds the Pulsar Admin Topic Policy handler function
func (b *PulsarAdminTopicPolicyToolBuilder) buildTopicPolicyHandler(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		operation, err := request.RequireString("operation")
		if err != nil {
			return mcp.NewToolResultError("Missing required parameter 'operation'"), nil
		}
		operation = normalizeTopicPolicyOperation(operation)

		topic, err := request.RequireString("topic")
		if err != nil {
			return mcp.NewToolResultError("Missing required parameter 'topic'"), nil
		}

		if readOnly && isReadOnlyRestrictedTopicPolicyOperation(operation) {
			return mcp.NewToolResultError("Write operations not allowed in read-only mode"), nil
		}

		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return mcp.NewToolResultError("Pulsar session not found in context"), nil
		}

		client, err := session.GetAdminClient()
		if err != nil {
			return b.handleError("get admin client", err), nil
		}

		switch operation {
		case "get-retention":
			return b.handleGetTopicRetention(client, topic, request)
		case "set-retention":
			return b.handleSetTopicRetention(client, topic, request)
		case "remove-retention":
			return b.handleRemoveTopicRetention(client, topic)
		case "get-message-ttl":
			return b.handleGetTopicTTL(client, topic)
		case "set-message-ttl":
			return b.handleSetTopicTTL(client, topic, request)
		case "remove-message-ttl":
			return b.handleRemoveTopicTTL(client, topic)
		case "get-max-producers":
			return b.handleGetTopicCountPolicy(topic, client.Topics().GetMaxProducers, "max producers")
		case "set-max-producers":
			return b.handleSetTopicCountPolicy(topic, request, []string{"count", "max-producers"}, client.Topics().SetMaxProducers, "max producers")
		case "remove-max-producers":
			return b.handleRemoveTopicCountPolicy(topic, client.Topics().RemoveMaxProducers, "max producers")
		case "get-max-consumers":
			return b.handleGetTopicCountPolicy(topic, client.Topics().GetMaxConsumers, "max consumers")
		case "set-max-consumers":
			return b.handleSetTopicCountPolicy(topic, request, []string{"count", "max-consumers"}, client.Topics().SetMaxConsumers, "max consumers")
		case "remove-max-consumers":
			return b.handleRemoveTopicCountPolicy(topic, client.Topics().RemoveMaxConsumers, "max consumers")
		case "get-max-unacked-messages-per-consumer":
			return b.handleGetTopicCountPolicy(topic, client.Topics().GetMaxUnackMessagesPerConsumer, "max unacked messages per consumer")
		case "set-max-unacked-messages-per-consumer":
			return b.handleSetTopicCountPolicy(topic, request, []string{"count", "maxNum", "max-num"}, client.Topics().SetMaxUnackMessagesPerConsumer, "max unacked messages per consumer")
		case "remove-max-unacked-messages-per-consumer":
			return b.handleRemoveTopicCountPolicy(topic, client.Topics().RemoveMaxUnackMessagesPerConsumer, "max unacked messages per consumer")
		case "get-max-unacked-messages-per-subscription":
			return b.handleGetTopicCountPolicy(topic, client.Topics().GetMaxUnackMessagesPerSubscription, "max unacked messages per subscription")
		case "set-max-unacked-messages-per-subscription":
			return b.handleSetTopicCountPolicy(topic, request, []string{"count", "maxNum", "max-num"}, client.Topics().SetMaxUnackMessagesPerSubscription, "max unacked messages per subscription")
		case "remove-max-unacked-messages-per-subscription":
			return b.handleRemoveTopicCountPolicy(topic, client.Topics().RemoveMaxUnackMessagesPerSubscription, "max unacked messages per subscription")
		case "get-persistence":
			return b.handleGetTopicPersistence(client, topic)
		case "set-persistence":
			return b.handleSetTopicPersistence(client, topic, request)
		case "remove-persistence":
			return b.handleRemoveTopicPersistence(client, topic)
		case "get-delayed-delivery":
			return b.handleGetTopicDelayedDelivery(client, topic)
		case "set-delayed-delivery":
			return b.handleSetTopicDelayedDelivery(client, topic, request)
		case "remove-delayed-delivery":
			return b.handleRemoveTopicDelayedDelivery(client, topic)
		case "get-dispatch-rate":
			return b.handleGetTopicDispatchRate(client, topic)
		case "set-dispatch-rate":
			return b.handleSetTopicDispatchRate(client, topic, request)
		case "remove-dispatch-rate":
			return b.handleRemoveTopicDispatchRate(client, topic)
		case "get-subscription-dispatch-rate":
			return b.handleGetTopicSubscriptionDispatchRate(client, topic)
		case "set-subscription-dispatch-rate":
			return b.handleSetTopicSubscriptionDispatchRate(client, topic, request)
		case "remove-subscription-dispatch-rate":
			return b.handleRemoveTopicSubscriptionDispatchRate(client, topic)
		case "get-deduplication":
			return b.handleGetTopicDeduplication(client, topic)
		case "set-deduplication":
			return b.handleSetTopicDeduplication(client, topic, request)
		case "remove-deduplication":
			return b.handleRemoveTopicDeduplication(client, topic)
		case "get-backlog-quotas":
			return b.handleGetTopicBacklogQuotas(client, topic, request)
		case "set-backlog-quota":
			return b.handleSetTopicBacklogQuota(client, topic, request)
		case "remove-backlog-quota":
			return b.handleRemoveTopicBacklogQuota(client, topic, request)
		case "get-compaction-threshold":
			return b.handleGetTopicCompaction(client, topic, request)
		case "set-compaction-threshold":
			return b.handleSetTopicCompaction(client, topic, request)
		case "remove-compaction-threshold":
			return b.handleRemoveTopicCompaction(client, topic)
		case "get-publish-rate":
			return b.handleGetTopicPublishRate(client, topic)
		case "set-publish-rate":
			return b.handleSetTopicPublishRate(client, topic, request)
		case "remove-publish-rate":
			return b.handleRemoveTopicPublishRate(client, topic)
		case "get-inactive-topic-policies":
			return b.handleGetTopicInactiveTopicPolicies(client, topic, request)
		case "set-inactive-topic-policies":
			return b.handleSetTopicInactiveTopicPolicies(client, topic, request)
		case "remove-inactive-topic-policies":
			return b.handleRemoveTopicInactiveTopicPolicies(client, topic)
		case "get-subscription-types":
			return b.handleGetTopicSubscriptionTypes(client, topic)
		case "set-subscription-types":
			return b.handleSetTopicSubscriptionTypes(client, topic, request)
		case "remove-subscription-types":
			return b.handleRemoveTopicSubscriptionTypes(client, topic)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Unsupported operation: %s", operation)), nil
		}
	}
}

func normalizeTopicPolicyOperation(operation string) string {
	normalized := strings.ToLower(strings.TrimSpace(operation))
	if alias, ok := topicPolicyOperationAliases[normalized]; ok {
		return alias
	}

	normalized = strings.ReplaceAll(normalized, "_", "-")
	if alias, ok := topicPolicyOperationAliases[normalized]; ok {
		return alias
	}

	return normalized
}

func isReadOnlyRestrictedTopicPolicyOperation(operation string) bool {
	_, ok := readOnlyRestrictedTopicPolicyOperations[normalizeTopicPolicyOperation(operation)]
	return ok
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleError(operation string, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("Failed to %s: %v", operation, err))
}

func (b *PulsarAdminTopicPolicyToolBuilder) marshalResponse(data any) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return b.handleError("marshal response", err), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) parseTopicName(topic string) (*utils.TopicName, error) {
	return utils.GetTopicName(strings.TrimSpace(topic))
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleGetTopicRetention(client cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	retention, err := client.Topics().GetRetention(*topicName, request.GetBool("applied", false))
	if err != nil {
		return b.handleError("get topic retention policy", err), nil
	}

	if retention == nil {
		return mcp.NewToolResultText(fmt.Sprintf("No retention policy found for topic %s", topicName.String())), nil
	}

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

	return mcp.NewToolResultText(fmt.Sprintf("Retention policy for topic %s: %s and %s",
		topicName.String(), retentionTime, retentionSize)), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleSetTopicRetention(client cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	retentionTimeInMinutes := int64(-1)
	retentionSizeInMB := int64(-1)

	if retentionTime := b.getStringArgument(request, "", "retention-time", "retention_time", "time"); retentionTime != "" {
		parsed, parseErr := b.parseRetentionTime(retentionTime)
		if parseErr != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid retention time format: %v", parseErr)), nil
		}
		retentionTimeInMinutes = parsed
	}

	if retentionSize := b.getStringArgument(request, "", "retention-size", "retention_size", "size"); retentionSize != "" {
		parsed, parseErr := b.parseRetentionSize(retentionSize)
		if parseErr != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid retention size format: %v", parseErr)), nil
		}
		retentionSizeInMB = parsed
	}

	retentionPolicy := utils.RetentionPolicies{
		RetentionTimeInMinutes: int(retentionTimeInMinutes),
		RetentionSizeInMB:      retentionSizeInMB,
	}

	if err := client.Topics().SetRetention(*topicName, retentionPolicy); err != nil {
		return b.handleError("set topic retention policy", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Retention policy set for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleRemoveTopicRetention(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	if err := client.Topics().RemoveRetention(*topicName); err != nil {
		return b.handleError("remove topic retention policy", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Retention policy removed for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleGetTopicTTL(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	ttl, err := client.Topics().GetMessageTTL(*topicName)
	if err != nil {
		return b.handleError("get topic message TTL", err), nil
	}

	if ttl == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Message TTL is not configured for topic %s", topicName.String())), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Message TTL for topic %s is %d seconds", topicName.String(), ttl)), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleSetTopicTTL(client cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	ttlSeconds, err := b.requireIntegerArgument(request, "ttl-seconds", "ttl_seconds", "ttl")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if ttlSeconds < 0 {
		return mcp.NewToolResultError("TTL seconds must be non-negative"), nil
	}

	if err := client.Topics().SetMessageTTL(*topicName, ttlSeconds); err != nil {
		return b.handleError("set topic message TTL", err), nil
	}

	if ttlSeconds == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Message TTL disabled for topic %s", topicName.String())), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Message TTL set to %d seconds for topic %s", ttlSeconds, topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleRemoveTopicTTL(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	if err := client.Topics().RemoveMessageTTL(*topicName); err != nil {
		return b.handleError("remove topic message TTL", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Message TTL removed for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleGetTopicCountPolicy(
	topic string,
	getter func(utils.TopicName) (int, error),
	label string,
) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	value, err := getter(*topicName)
	if err != nil {
		return b.handleError(fmt.Sprintf("get %s", label), err), nil
	}

	if value == -1 {
		return mcp.NewToolResultText(fmt.Sprintf("%s is not set for topic %s", label, topicName.String())), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("%s for topic %s is %d", label, topicName.String(), value)), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleSetTopicCountPolicy(
	topic string,
	request mcp.CallToolRequest,
	keys []string,
	setter func(utils.TopicName, int) error,
	label string,
) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	count, err := b.requireIntegerArgument(request, keys...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if count < 0 {
		return mcp.NewToolResultError(fmt.Sprintf("%s must be non-negative", label)), nil
	}

	if err := setter(*topicName, count); err != nil {
		return b.handleError(fmt.Sprintf("set %s", label), err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set %s for topic %s to %d", label, topicName.String(), count)), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleRemoveTopicCountPolicy(
	topic string,
	remover func(utils.TopicName) error,
	label string,
) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	if err := remover(*topicName); err != nil {
		return b.handleError(fmt.Sprintf("remove %s", label), err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Removed %s policy for topic %s", label, topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleGetTopicPersistence(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	persistence, err := client.Topics().GetPersistence(*topicName)
	if err != nil {
		return b.handleError("get topic persistence", err), nil
	}

	return b.marshalResponse(persistence)
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleSetTopicPersistence(client cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	persistence, err := b.buildPersistenceData(request)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build persistence data: %v", err)), nil
	}

	if err := client.Topics().SetPersistence(*topicName, persistence); err != nil {
		return b.handleError("set topic persistence", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set persistence policy for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleRemoveTopicPersistence(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	if err := client.Topics().RemovePersistence(*topicName); err != nil {
		return b.handleError("remove topic persistence", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Removed persistence policy for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleGetTopicDelayedDelivery(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	data, err := client.Topics().GetDelayedDelivery(*topicName)
	if err != nil {
		return b.handleError("get topic delayed delivery", err), nil
	}

	return b.marshalResponse(data)
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleSetTopicDelayedDelivery(client cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	data, err := b.buildDelayedDeliveryData(request)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build delayed delivery policy: %v", err)), nil
	}

	if err := client.Topics().SetDelayedDelivery(*topicName, data); err != nil {
		return b.handleError("set topic delayed delivery", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set delayed delivery policy for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleRemoveTopicDelayedDelivery(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	if err := client.Topics().RemoveDelayedDelivery(*topicName); err != nil {
		return b.handleError("remove topic delayed delivery", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Removed delayed delivery policy for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleGetTopicDispatchRate(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	rate, err := client.Topics().GetDispatchRate(*topicName)
	if err != nil {
		return b.handleError("get topic dispatch rate", err), nil
	}

	return b.marshalResponse(rate)
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleSetTopicDispatchRate(client cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	rate, err := b.buildDispatchRateData(request)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build dispatch rate: %v", err)), nil
	}

	if err := client.Topics().SetDispatchRate(*topicName, rate); err != nil {
		return b.handleError("set topic dispatch rate", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set dispatch rate for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleRemoveTopicDispatchRate(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	if err := client.Topics().RemoveDispatchRate(*topicName); err != nil {
		return b.handleError("remove topic dispatch rate", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Removed dispatch rate policy for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleGetTopicSubscriptionDispatchRate(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	rate, err := client.Topics().GetSubscriptionDispatchRate(*topicName)
	if err != nil {
		return b.handleError("get topic subscription dispatch rate", err), nil
	}

	return b.marshalResponse(rate)
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleSetTopicSubscriptionDispatchRate(client cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	rate, err := b.buildDispatchRateData(request)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build subscription dispatch rate: %v", err)), nil
	}

	if err := client.Topics().SetSubscriptionDispatchRate(*topicName, rate); err != nil {
		return b.handleError("set topic subscription dispatch rate", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set subscription dispatch rate for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleRemoveTopicSubscriptionDispatchRate(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	if err := client.Topics().RemoveSubscriptionDispatchRate(*topicName); err != nil {
		return b.handleError("remove topic subscription dispatch rate", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Removed subscription dispatch rate policy for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleGetTopicDeduplication(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	enabled, err := client.Topics().GetDeduplicationStatus(*topicName)
	if err != nil {
		return b.handleError("get topic deduplication", err), nil
	}

	return b.marshalResponse(map[string]any{
		"topic":   topicName.String(),
		"enabled": enabled,
	})
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleSetTopicDeduplication(client cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	enabled, err := b.resolveToggleArgument(request)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if err := client.Topics().SetDeduplicationStatus(*topicName, enabled); err != nil {
		return b.handleError("set topic deduplication", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set deduplication policy for topic %s to %t", topicName.String(), enabled)), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleRemoveTopicDeduplication(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	if err := client.Topics().RemoveDeduplicationStatus(*topicName); err != nil {
		return b.handleError("remove topic deduplication", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Removed deduplication policy for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleGetTopicBacklogQuotas(client cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	quotas, err := client.Topics().GetBacklogQuotaMap(*topicName, request.GetBool("applied", false))
	if err != nil {
		return b.handleError("get topic backlog quotas", err), nil
	}

	return b.marshalResponse(quotas)
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleSetTopicBacklogQuota(client cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	quota, quotaType, err := b.buildBacklogQuota(request)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build backlog quota: %v", err)), nil
	}

	if err := client.Topics().SetBacklogQuota(*topicName, quota, quotaType); err != nil {
		return b.handleError("set topic backlog quota", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set backlog quota for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleRemoveTopicBacklogQuota(client cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	quotaTypeStr := b.getStringArgument(request, string(utils.DestinationStorage), "type")
	quotaType, err := utils.ParseBacklogQuotaType(quotaTypeStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid backlog quota type: %v", err)), nil
	}

	if err := client.Topics().RemoveBacklogQuota(*topicName, quotaType); err != nil {
		return b.handleError("remove topic backlog quota", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Removed backlog quota for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleGetTopicCompaction(client cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	threshold, err := client.Topics().GetCompactionThreshold(*topicName, request.GetBool("applied", false))
	if err != nil {
		return b.handleError("get topic compaction threshold", err), nil
	}

	if threshold == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Automatic compaction is disabled for topic %s", topicName.String())), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("The compaction threshold of topic %s is %d bytes", topicName.String(), threshold)), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleSetTopicCompaction(client cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	threshold, err := b.buildCompactionThreshold(request)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid compaction threshold: %v", err)), nil
	}
	if threshold < 0 {
		return mcp.NewToolResultError("Compaction threshold must be non-negative"), nil
	}

	if err := client.Topics().SetCompactionThreshold(*topicName, threshold); err != nil {
		return b.handleError("set topic compaction threshold", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Compaction threshold set to %d bytes for topic %s", threshold, topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleRemoveTopicCompaction(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	if err := client.Topics().RemoveCompactionThreshold(*topicName); err != nil {
		return b.handleError("remove topic compaction threshold", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Compaction threshold removed for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleGetTopicPublishRate(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	rate, err := client.Topics().GetPublishRate(*topicName)
	if err != nil {
		return b.handleError("get topic publish rate", err), nil
	}

	return b.marshalResponse(rate)
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleSetTopicPublishRate(client cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	rate, err := b.buildPublishRateData(request)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build publish rate: %v", err)), nil
	}

	if err := client.Topics().SetPublishRate(*topicName, rate); err != nil {
		return b.handleError("set topic publish rate", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set publish rate for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleRemoveTopicPublishRate(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	if err := client.Topics().RemovePublishRate(*topicName); err != nil {
		return b.handleError("remove topic publish rate", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Removed publish rate policy for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleGetTopicInactiveTopicPolicies(client cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	policies, err := client.Topics().GetInactiveTopicPolicies(*topicName, request.GetBool("applied", false))
	if err != nil {
		return b.handleError("get inactive topic policies", err), nil
	}

	return b.marshalResponse(policies)
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleSetTopicInactiveTopicPolicies(client cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	policies, err := b.buildInactiveTopicPolicies(request)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build inactive topic policies: %v", err)), nil
	}

	if err := client.Topics().SetInactiveTopicPolicies(*topicName, policies); err != nil {
		return b.handleError("set inactive topic policies", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set inactive topic policies for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleRemoveTopicInactiveTopicPolicies(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	if err := client.Topics().RemoveInactiveTopicPolicies(*topicName); err != nil {
		return b.handleError("remove inactive topic policies", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Removed inactive topic policies for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleGetTopicSubscriptionTypes(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	type subscriptionTypesGetter interface {
		GetSubscriptionTypesEnabled(utils.TopicName) ([]string, error)
	}

	if getter, ok := client.Topics().(subscriptionTypesGetter); ok {
		subscriptionTypes, err := getter.GetSubscriptionTypesEnabled(*topicName)
		if err != nil {
			return b.handleError("get topic subscription types", err), nil
		}

		if len(subscriptionTypes) == 0 {
			return mcp.NewToolResultText(fmt.Sprintf("No subscription type restrictions configured for topic %s (all types allowed)", topicName.String())), nil
		}

		return b.marshalResponse(map[string]any{
			"topic":             topicName.String(),
			"subscriptionTypes": subscriptionTypes,
		})
	}

	return mcp.NewToolResultError("Subscription types policy management is not available in the current pulsarctl API version. This feature may require a newer version of Pulsar or pulsarctl."), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleSetTopicSubscriptionTypes(client cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	subscriptionTypes, err := b.requireStringSliceArgument(request, "subscription-types", "subscription_types")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	validTypes := map[string]bool{
		"Exclusive":  true,
		"Shared":     true,
		"Failover":   true,
		"Key_Shared": true,
	}

	validatedTypes := make([]string, 0, len(subscriptionTypes))
	for _, subType := range subscriptionTypes {
		if !validTypes[subType] {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid subscription type: %s. Valid types are: Exclusive, Shared, Failover, Key_Shared", subType)), nil
		}
		validatedTypes = append(validatedTypes, subType)
	}

	if len(validatedTypes) == 0 {
		return mcp.NewToolResultError("At least one valid subscription type must be specified"), nil
	}

	type subscriptionTypesSetter interface {
		SetSubscriptionTypesEnabled(utils.TopicName, []string) error
	}

	if setter, ok := client.Topics().(subscriptionTypesSetter); ok {
		if err := setter.SetSubscriptionTypesEnabled(*topicName, validatedTypes); err != nil {
			return b.handleError("set topic subscription types", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Subscription types set for topic %s: %s",
			topicName.String(), strings.Join(validatedTypes, ", "))), nil
	}

	return mcp.NewToolResultError("Subscription types policy management is not available in the current pulsarctl API version. This feature may require a newer version of Pulsar or pulsarctl."), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleRemoveTopicSubscriptionTypes(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	topicName, err := b.parseTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	type subscriptionTypesRemover interface {
		RemoveSubscriptionTypesEnabled(utils.TopicName) error
	}

	if remover, ok := client.Topics().(subscriptionTypesRemover); ok {
		if err := remover.RemoveSubscriptionTypesEnabled(*topicName); err != nil {
			return b.handleError("remove topic subscription types policy", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Subscription types policy removed for topic %s (all types now allowed)", topicName.String())), nil
	}

	return mcp.NewToolResultError("Subscription types policy management is not available in the current pulsarctl API version. This feature may require a newer version of Pulsar or pulsarctl."), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) buildPersistenceData(request mcp.CallToolRequest) (utils.PersistenceData, error) {
	ensemble, err := b.requireInt64Argument(request, "bookkeeper-ensemble", "ensemble-size")
	if err != nil {
		return utils.PersistenceData{}, err
	}
	writeQuorum, err := b.requireInt64Argument(request, "bookkeeper-write-quorum", "write-quorum-size")
	if err != nil {
		return utils.PersistenceData{}, err
	}
	ackQuorum, err := b.requireInt64Argument(request, "bookkeeper-ack-quorum", "ack-quorum-size")
	if err != nil {
		return utils.PersistenceData{}, err
	}

	key, ok := b.firstArgumentKey(request, "ml-mark-delete-max-rate")
	if !ok {
		return utils.PersistenceData{}, fmt.Errorf("missing required parameter 'ml-mark-delete-max-rate'")
	}
	rate, err := request.RequireFloat(key)
	if err != nil {
		return utils.PersistenceData{}, fmt.Errorf("missing required parameter 'ml-mark-delete-max-rate': %w", err)
	}

	return utils.PersistenceData{
		BookkeeperEnsemble:             ensemble,
		BookkeeperWriteQuorum:          writeQuorum,
		BookkeeperAckQuorum:            ackQuorum,
		ManagedLedgerMaxMarkDeleteRate: rate,
	}, nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) buildDelayedDeliveryData(request mcp.CallToolRequest) (utils.DelayedDeliveryData, error) {
	enabled, err := b.resolveToggleArgument(request)
	if err != nil {
		return utils.DelayedDeliveryData{}, err
	}

	data := utils.DelayedDeliveryData{Active: enabled}
	if !enabled {
		return data, nil
	}

	timeValue := b.getStringArgument(request, "1s", "time", "tick-time")
	tickTime, err := ctlutil.ParseRelativeTimeInSeconds(timeValue)
	if err != nil {
		return utils.DelayedDeliveryData{}, err
	}
	data.TickTime = tickTime.Seconds()

	maxDelay, err := b.getInt64Argument(request, 0, "max-delay-ms")
	if err != nil {
		return utils.DelayedDeliveryData{}, err
	}
	data.MaxDelayInMillis = maxDelay

	return data, nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) buildDispatchRateData(request mcp.CallToolRequest) (utils.DispatchRateData, error) {
	msgRate, err := b.getInt64Argument(request, -1, "msg-rate", "msg-dispatch-rate")
	if err != nil {
		return utils.DispatchRateData{}, err
	}
	byteRate, err := b.getInt64Argument(request, -1, "byte-rate", "byte-dispatch-rate")
	if err != nil {
		return utils.DispatchRateData{}, err
	}
	period, err := b.getInt64Argument(request, 1, "period", "dispatch-rate-period")
	if err != nil {
		return utils.DispatchRateData{}, err
	}

	return utils.DispatchRateData{
		DispatchThrottlingRateInMsg:  msgRate,
		DispatchThrottlingRateInByte: byteRate,
		RatePeriodInSecond:           period,
		RelativeToPublishRate:        b.getBoolArgument(request, false, "relative-to-publish-rate"),
	}, nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) buildPublishRateData(request mcp.CallToolRequest) (utils.PublishRateData, error) {
	msgRate, err := b.getInt64Argument(request, -1, "msg-rate", "msg-publish-rate")
	if err != nil {
		return utils.PublishRateData{}, err
	}
	byteRate, err := b.getInt64Argument(request, -1, "byte-rate", "byte-publish-rate")
	if err != nil {
		return utils.PublishRateData{}, err
	}

	return utils.PublishRateData{
		PublishThrottlingRateInMsg:  msgRate,
		PublishThrottlingRateInByte: byteRate,
	}, nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) buildBacklogQuota(request mcp.CallToolRequest) (utils.BacklogQuota, utils.BacklogQuotaType, error) {
	limitSizeStr, err := b.requireStringArgument(request, "limit-size")
	if err != nil {
		return utils.BacklogQuota{}, "", err
	}
	limitSize, err := ctlutil.ValidateSizeString(limitSizeStr)
	if err != nil {
		return utils.BacklogQuota{}, "", err
	}

	policyStr, err := b.requireStringArgument(request, "policy", "backlog-policy")
	if err != nil {
		return utils.BacklogQuota{}, "", err
	}
	policy, err := utils.ParseRetentionPolicy(policyStr)
	if err != nil {
		return utils.BacklogQuota{}, "", err
	}

	limitTime, err := b.getInt64Argument(request, -1, "limit-time")
	if err != nil {
		return utils.BacklogQuota{}, "", err
	}

	quotaTypeStr := b.getStringArgument(request, string(utils.DestinationStorage), "type")
	quotaType, err := utils.ParseBacklogQuotaType(quotaTypeStr)
	if err != nil {
		return utils.BacklogQuota{}, "", err
	}

	return utils.NewBacklogQuota(limitSize, limitTime, policy), quotaType, nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) buildInactiveTopicPolicies(request mcp.CallToolRequest) (utils.InactiveTopicPolicies, error) {
	deleteModeStr, err := b.requireStringArgument(request, "delete-mode")
	if err != nil {
		return utils.InactiveTopicPolicies{}, err
	}
	deleteMode, err := utils.ParseInactiveTopicDeleteMode(deleteModeStr)
	if err != nil {
		return utils.InactiveTopicPolicies{}, err
	}

	maxInactiveDuration, err := b.requireStringArgument(request, "max-inactive-duration")
	if err != nil {
		return utils.InactiveTopicPolicies{}, err
	}
	duration, err := ctlutil.ParseRelativeTimeInSeconds(maxInactiveDuration)
	if err != nil {
		return utils.InactiveTopicPolicies{}, err
	}

	deleteWhileInactive := b.getBoolArgument(request, false, "delete-while-inactive", "enable-delete-while-inactive")

	return utils.NewInactiveTopicPolicies(&deleteMode, int(duration.Seconds()), deleteWhileInactive), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) buildCompactionThreshold(request mcp.CallToolRequest) (int64, error) {
	if threshold := b.getStringArgument(request, "", "threshold"); threshold != "" {
		return ctlutil.ValidateSizeString(threshold)
	}

	return b.requireInt64Argument(request, "compaction_threshold")
}

func (b *PulsarAdminTopicPolicyToolBuilder) resolveToggleArgument(request mcp.CallToolRequest) (bool, error) {
	if key, ok := b.firstArgumentKey(request, "enabled"); ok {
		enabled, err := request.RequireBool(key)
		if err != nil {
			return false, fmt.Errorf("missing required parameter 'enabled': %w", err)
		}
		return enabled, nil
	}

	enable := b.getBoolArgument(request, false, "enable")
	disable := b.getBoolArgument(request, false, "disable")
	if enable == disable {
		return false, fmt.Errorf("need to specify either 'enable' or 'disable'")
	}
	return enable, nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) firstArgumentKey(request mcp.CallToolRequest, keys ...string) (string, bool) {
	args := request.GetArguments()
	for _, key := range keys {
		if _, ok := args[key]; ok {
			return key, true
		}
	}
	return "", false
}

func (b *PulsarAdminTopicPolicyToolBuilder) getStringArgument(request mcp.CallToolRequest, defaultValue string, keys ...string) string {
	if key, ok := b.firstArgumentKey(request, keys...); ok {
		return request.GetString(key, defaultValue)
	}
	return defaultValue
}

func (b *PulsarAdminTopicPolicyToolBuilder) requireStringArgument(request mcp.CallToolRequest, keys ...string) (string, error) {
	key, ok := b.firstArgumentKey(request, keys...)
	if !ok {
		return "", fmt.Errorf("missing required parameter '%s'", keys[0])
	}

	value, err := request.RequireString(key)
	if err != nil {
		return "", fmt.Errorf("missing required parameter '%s': %w", keys[0], err)
	}
	return value, nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) getBoolArgument(request mcp.CallToolRequest, defaultValue bool, keys ...string) bool {
	if key, ok := b.firstArgumentKey(request, keys...); ok {
		return request.GetBool(key, defaultValue)
	}
	return defaultValue
}

func (b *PulsarAdminTopicPolicyToolBuilder) requireIntegerArgument(request mcp.CallToolRequest, keys ...string) (int, error) {
	value, err := b.requireFloatArgument(request, keys...)
	if err != nil {
		return 0, err
	}
	return b.floatToInt(keys[0], value)
}

func (b *PulsarAdminTopicPolicyToolBuilder) requireInt64Argument(request mcp.CallToolRequest, keys ...string) (int64, error) {
	value, err := b.requireFloatArgument(request, keys...)
	if err != nil {
		return 0, err
	}
	if math.Trunc(value) != value {
		return 0, fmt.Errorf("parameter '%s' must be an integer", keys[0])
	}
	return int64(value), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) getInt64Argument(request mcp.CallToolRequest, defaultValue int64, keys ...string) (int64, error) {
	key, ok := b.firstArgumentKey(request, keys...)
	if !ok {
		return defaultValue, nil
	}

	value := request.GetFloat(key, float64(defaultValue))
	if math.Trunc(value) != value {
		return 0, fmt.Errorf("parameter '%s' must be an integer", key)
	}

	return int64(value), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) requireFloatArgument(request mcp.CallToolRequest, keys ...string) (float64, error) {
	key, ok := b.firstArgumentKey(request, keys...)
	if !ok {
		return 0, fmt.Errorf("missing required parameter '%s'", keys[0])
	}

	value, err := request.RequireFloat(key)
	if err != nil {
		return 0, fmt.Errorf("missing required parameter '%s': %w", keys[0], err)
	}
	return value, nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) requireStringSliceArgument(request mcp.CallToolRequest, keys ...string) ([]string, error) {
	key, ok := b.firstArgumentKey(request, keys...)
	if !ok {
		return nil, fmt.Errorf("missing required parameter '%s'", keys[0])
	}

	raw, ok := request.GetArguments()[key]
	if !ok {
		return nil, fmt.Errorf("missing required parameter '%s'", keys[0])
	}

	switch value := raw.(type) {
	case []string:
		return value, nil
	case []any:
		items := make([]string, 0, len(value))
		for _, item := range value {
			text, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("parameter '%s' must be an array of strings", keys[0])
			}
			items = append(items, text)
		}
		return items, nil
	default:
		items, err := request.RequireStringSlice(key)
		if err != nil {
			return nil, fmt.Errorf("missing required parameter '%s': %w", keys[0], err)
		}
		return items, nil
	}
}

func (b *PulsarAdminTopicPolicyToolBuilder) floatToInt(key string, value float64) (int, error) {
	if math.Trunc(value) != value {
		return 0, fmt.Errorf("parameter '%s' must be an integer", key)
	}
	return int(value), nil
}

// Utility functions for parsing retention parameters

// parseRetentionTime parses retention time strings like "1d", "24h", "1440m" and returns minutes
func (b *PulsarAdminTopicPolicyToolBuilder) parseRetentionTime(retentionTime string) (int64, error) {
	if retentionTime == "" {
		return -1, nil
	}

	retentionTime = strings.TrimSpace(retentionTime)
	if retentionTime == "0" {
		return 0, nil
	}
	if retentionTime == "-1" {
		return -1, nil
	}

	if len(retentionTime) < 2 {
		return -1, fmt.Errorf("invalid retention time format: %s", retentionTime)
	}

	valueStr := retentionTime[:len(retentionTime)-1]
	unit := retentionTime[len(retentionTime)-1:]

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return -1, fmt.Errorf("invalid retention time value: %s", valueStr)
	}

	switch unit {
	case "m":
		return int64(value), nil
	case "h":
		return int64(value * 60), nil
	case "d":
		return int64(value * 60 * 24), nil
	default:
		return -1, fmt.Errorf("invalid retention time unit: %s (use m, h, or d)", unit)
	}
}

// parseRetentionSize parses retention size strings like "100MB", "1GB" and returns MB
func (b *PulsarAdminTopicPolicyToolBuilder) parseRetentionSize(retentionSize string) (int64, error) {
	if retentionSize == "" {
		return -1, nil
	}

	retentionSize = strings.TrimSpace(strings.ToUpper(retentionSize))
	if retentionSize == "0" {
		return 0, nil
	}
	if retentionSize == "-1" {
		return -1, nil
	}

	var (
		value float64
		unit  string
		err   error
	)

	// /nolint:gocritic
	if strings.HasSuffix(retentionSize, "TB") || strings.HasSuffix(retentionSize, "T") {
		if strings.HasSuffix(retentionSize, "TB") {
			value, err = strconv.ParseFloat(retentionSize[:len(retentionSize)-2], 64)
			unit = "TB"
		} else {
			value, err = strconv.ParseFloat(retentionSize[:len(retentionSize)-1], 64)
			unit = "T"
		}
	} else if strings.HasSuffix(retentionSize, "GB") || strings.HasSuffix(retentionSize, "G") {
		if strings.HasSuffix(retentionSize, "GB") {
			value, err = strconv.ParseFloat(retentionSize[:len(retentionSize)-2], 64)
			unit = "GB"
		} else {
			value, err = strconv.ParseFloat(retentionSize[:len(retentionSize)-1], 64)
			unit = "G"
		}
	} else if strings.HasSuffix(retentionSize, "MB") || strings.HasSuffix(retentionSize, "M") {
		if strings.HasSuffix(retentionSize, "MB") {
			value, err = strconv.ParseFloat(retentionSize[:len(retentionSize)-2], 64)
			unit = "MB"
		} else {
			value, err = strconv.ParseFloat(retentionSize[:len(retentionSize)-1], 64)
			unit = "M"
		}
	} else {
		return -1, fmt.Errorf("invalid retention size format: %s (use MB, GB, or TB)", retentionSize)
	}

	if err != nil {
		return -1, fmt.Errorf("invalid retention size value: %v", err)
	}

	switch unit {
	case "MB", "M":
		return int64(value), nil
	case "GB", "G":
		return int64(value * 1024), nil
	case "TB", "T":
		return int64(value * 1024 * 1024), nil
	default:
		return -1, fmt.Errorf("invalid retention size unit: %s", unit)
	}
}
