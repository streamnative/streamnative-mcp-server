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
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

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
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool := b.buildTopicPolicyTool()
	handler := b.buildTopicPolicyHandler(config.ReadOnly)

	return []server.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildTopicPolicyTool builds the Pulsar Admin Topic Policy MCP tool definition
func (b *PulsarAdminTopicPolicyToolBuilder) buildTopicPolicyTool() mcp.Tool {
	toolDesc := "Manage Pulsar topic policies including retention, TTL, compaction, and subscription policies. " +
		"This tool provides functionality to get, set, and remove various topic-level policies in Apache Pulsar."

	operationDesc := "Operation to perform on topic policies. Available operations:\n" +
		"- get_retention: Get message retention policy for a topic\n" +
		"- set_retention: Set message retention policy for a topic\n" +
		"- remove_retention: Remove message retention policy for a topic\n" +
		"- get_ttl: Get message TTL policy for a topic\n" +
		"- set_ttl: Set message TTL policy for a topic\n" +
		"- remove_ttl: Remove message TTL policy for a topic\n" +
		"- get_compaction: Get compaction policy for a topic\n" +
		"- set_compaction: Set compaction policy for a topic\n" +
		"- remove_compaction: Remove compaction policy for a topic\n" +
		"- get_subscription_types: Get allowed subscription types for a topic\n" +
		"- set_subscription_types: Set allowed subscription types for a topic\n" +
		"- remove_subscription_types: Remove subscription types restriction for a topic"

	return mcp.NewTool("pulsar_admin_topic_policy",
		mcp.WithDescription(toolDesc),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
		),
		mcp.WithString("topic", mcp.Required(),
			mcp.Description("Topic name in format 'persistent://tenant/namespace/topic' or 'tenant/namespace/topic'"),
		),
		mcp.WithString("retention_size",
			mcp.Description("Retention size policy (e.g., '100MB', '1GB') - used with retention operations"),
		),
		mcp.WithString("retention_time",
			mcp.Description("Retention time policy (e.g., '1d', '24h', '1440m') - used with retention operations"),
		),
		mcp.WithNumber("ttl_seconds",
			mcp.Description("TTL in seconds - used with TTL operations"),
		),
		mcp.WithNumber("compaction_threshold",
			mcp.Description("Compaction threshold in bytes - used with compaction operations"),
		),
		mcp.WithArray("subscription_types",
			mcp.Description("List of allowed subscription types - used with subscription type operations"),
			mcp.Items(
				map[string]interface{}{
					"type":        "string",
					"description": "subscription type: Exclusive, Shared, Failover, Key_Shared",
				},
			),
		),
	)
}

// buildTopicPolicyHandler builds the Pulsar Admin Topic Policy handler function
func (b *PulsarAdminTopicPolicyToolBuilder) buildTopicPolicyHandler(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		// Get required parameters
		operation, err := request.RequireString("operation")
		if err != nil {
			return mcp.NewToolResultError("Missing required operation parameter"), nil
		}

		topic, err := request.RequireString("topic")
		if err != nil {
			return mcp.NewToolResultError("Missing required topic parameter"), nil
		}

		// Check write operation permissions
		writeOps := []string{"set_retention", "remove_retention", "set_ttl", "remove_ttl",
			"set_compaction", "remove_compaction", "set_subscription_types", "remove_subscription_types"}
		isWriteOp := false
		for _, op := range writeOps {
			if operation == op {
				isWriteOp = true
				break
			}
		}

		if isWriteOp && readOnly {
			return mcp.NewToolResultError("Write operations not allowed in read-only mode"), nil
		}

		// Handle operations
		switch operation {
		case "get_retention":
			return b.handleGetTopicRetention(client, topic)
		case "set_retention":
			return b.handleSetTopicRetention(client, topic, request)
		case "remove_retention":
			return b.handleRemoveTopicRetention(client, topic)
		case "get_ttl":
			return b.handleGetTopicTTL(client, topic)
		case "set_ttl":
			return b.handleSetTopicTTL(client, topic, request)
		case "remove_ttl":
			return b.handleRemoveTopicTTL(client, topic)
		case "get_compaction":
			return b.handleGetTopicCompaction(client, topic)
		case "set_compaction":
			return b.handleSetTopicCompaction(client, topic, request)
		case "remove_compaction":
			return b.handleRemoveTopicCompaction(client, topic)
		case "get_subscription_types":
			return b.handleGetTopicSubscriptionTypes(client, topic)
		case "set_subscription_types":
			return b.handleSetTopicSubscriptionTypes(client, topic, request)
		case "remove_subscription_types":
			return b.handleRemoveTopicSubscriptionTypes(client, topic)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Unsupported operation: %s", operation)), nil
		}
	}
}

// Utility functions
func (b *PulsarAdminTopicPolicyToolBuilder) handleError(operation string, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("Failed to %s: %v", operation, err))
}

func (b *PulsarAdminTopicPolicyToolBuilder) marshalResponse(data interface{}) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return b.handleError("marshal response", err), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// Topic policy operation handlers
func (b *PulsarAdminTopicPolicyToolBuilder) handleGetTopicRetention(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	// Get retention policy
	retention, err := client.Topics().GetRetention(*topicName, false)
	if err != nil {
		return b.handleError("get topic retention policy", err), nil
	}

	// If no retention policy is defined
	if retention == nil {
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

	return mcp.NewToolResultText(fmt.Sprintf("Retention policy for topic %s: %s and %s",
		topicName.String(), retentionTime, retentionSize)), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleSetTopicRetention(client cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	// Parse retention parameters
	var retentionTimeInMinutes int64 = -1
	var retentionSizeInMB int64 = -1

	// /nolint:revive
	if retentionTime := request.GetString("retention_time", ""); retentionTime != "" {
		if parsed, err := b.parseRetentionTime(retentionTime); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid retention time format: %v", err)), nil
		} else {
			retentionTimeInMinutes = parsed
		}
	}

	// /nolint:revive
	if retentionSize := request.GetString("retention_size", ""); retentionSize != "" {
		if parsed, err := b.parseRetentionSize(retentionSize); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid retention size format: %v", err)), nil
		} else {
			retentionSizeInMB = parsed
		}
	}

	// Create retention policy
	retentionPolicy := utils.RetentionPolicies{
		RetentionTimeInMinutes: int(retentionTimeInMinutes),
		RetentionSizeInMB:      int64(retentionSizeInMB),
	}

	// Set retention policy
	err = client.Topics().SetRetention(*topicName, retentionPolicy)
	if err != nil {
		return b.handleError("set topic retention policy", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Retention policy set for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleRemoveTopicRetention(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	// Remove retention policy
	err = client.Topics().RemoveRetention(*topicName)
	if err != nil {
		return b.handleError("remove topic retention policy", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Retention policy removed for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleGetTopicTTL(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	// Get message TTL
	ttl, err := client.Topics().GetMessageTTL(*topicName)
	if err != nil {
		return b.handleError("get topic message TTL", err), nil
	}

	// Check if TTL is set
	if ttl == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Message TTL is not configured for topic %s", topicName.String())), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Message TTL for topic %s is %d seconds", topicName.String(), ttl)), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleSetTopicTTL(client cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	// Get TTL seconds parameter
	ttlSeconds, err := request.RequireFloat("ttl_seconds")
	if err != nil {
		return mcp.NewToolResultError("Missing required parameter 'ttl_seconds'"), nil
	}

	if ttlSeconds < 0 {
		return mcp.NewToolResultError("TTL seconds must be non-negative"), nil
	}

	// Set message TTL
	err = client.Topics().SetMessageTTL(*topicName, int(ttlSeconds))
	if err != nil {
		return b.handleError("set topic message TTL", err), nil
	}

	if ttlSeconds == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Message TTL disabled for topic %s", topicName.String())), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Message TTL set to %d seconds for topic %s", int(ttlSeconds), topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleRemoveTopicTTL(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	// Remove message TTL
	err = client.Topics().RemoveMessageTTL(*topicName)
	if err != nil {
		return b.handleError("remove topic message TTL", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Message TTL removed for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleGetTopicCompaction(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	// Get compaction threshold
	threshold, err := client.Topics().GetCompactionThreshold(*topicName, false)
	if err != nil {
		return b.handleError("get topic compaction threshold", err), nil
	}

	// Format the result
	if threshold == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Automatic compaction is disabled for topic %s", topicName.String())), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("The compaction threshold of the topic %s is %d byte(s)",
		topicName.String(), threshold)), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleSetTopicCompaction(client cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	// Get compaction threshold parameter
	thresholdNum, err := request.RequireFloat("compaction_threshold")
	if err != nil {
		return mcp.NewToolResultError("Missing required parameter 'compaction_threshold'"), nil
	}

	if thresholdNum < 0 {
		return mcp.NewToolResultError("Compaction threshold must be non-negative"), nil
	}

	threshold := int64(thresholdNum)

	// Set compaction threshold
	err = client.Topics().SetCompactionThreshold(*topicName, threshold)
	if err != nil {
		return b.handleError("set topic compaction threshold", err), nil
	}

	if threshold == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Automatic compaction disabled for topic %s", topicName.String())), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Compaction threshold set to %d bytes for topic %s", threshold, topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleRemoveTopicCompaction(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	// Remove compaction threshold
	err = client.Topics().RemoveCompactionThreshold(*topicName)
	if err != nil {
		return b.handleError("remove topic compaction threshold", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Compaction threshold removed for topic %s", topicName.String())), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleGetTopicSubscriptionTypes(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	// Check if the API supports subscription types management
	// Some versions of pulsarctl may not have this functionality
	topicsClient := client.Topics()

	// Try to use reflection or type assertion to check if the method exists
	type SubscriptionTypesGetter interface {
		GetSubscriptionTypesEnabled(utils.TopicName) ([]string, error)
	}

	if getter, ok := topicsClient.(SubscriptionTypesGetter); ok {
		subscriptionTypes, err := getter.GetSubscriptionTypesEnabled(*topicName)
		if err != nil {
			return b.handleError("get topic subscription types", err), nil
		}

		if len(subscriptionTypes) == 0 {
			return mcp.NewToolResultText(fmt.Sprintf("No subscription type restrictions configured for topic %s (all types allowed)", topicName.String())), nil
		}

		return b.marshalResponse(map[string]interface{}{
			"topic":             topicName.String(),
			"subscriptionTypes": subscriptionTypes,
		})
	}

	// Fallback: API not available in current version
	return mcp.NewToolResultError("Subscription types policy management is not available in the current pulsarctl API version. " +
		"This feature may require a newer version of Pulsar or pulsarctl."), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleSetTopicSubscriptionTypes(client cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	// Get subscription types parameter
	subscriptionTypes, err := request.RequireStringSlice("subscription_types")
	if err != nil {
		return mcp.NewToolResultError("Missing required parameter 'subscription_types'. " +
			"Please provide an array of subscription types: Exclusive, Shared, Failover, Key_Shared"), nil
	}

	// Validate subscription types
	validTypes := map[string]bool{
		"Exclusive":  true,
		"Shared":     true,
		"Failover":   true,
		"Key_Shared": true,
	}

	var validatedTypes []string
	for _, subType := range subscriptionTypes {
		if !validTypes[subType] {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid subscription type: %s. Valid types are: Exclusive, Shared, Failover, Key_Shared", subType)), nil
		}
		validatedTypes = append(validatedTypes, subType)
	}

	if len(validatedTypes) == 0 {
		return mcp.NewToolResultError("At least one valid subscription type must be specified"), nil
	}

	// Check if the API supports subscription types management
	topicsClient := client.Topics()

	// Try to use reflection or type assertion to check if the method exists
	type SubscriptionTypesSetter interface {
		SetSubscriptionTypesEnabled(utils.TopicName, []string) error
	}

	if setter, ok := topicsClient.(SubscriptionTypesSetter); ok {
		err := setter.SetSubscriptionTypesEnabled(*topicName, validatedTypes)
		if err != nil {
			return b.handleError("set topic subscription types", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Subscription types set for topic %s: %s",
			topicName.String(), strings.Join(validatedTypes, ", "))), nil
	}

	// Fallback: API not available in current version
	return mcp.NewToolResultError("Subscription types policy management is not available in the current pulsarctl API version. " +
		"This feature may require a newer version of Pulsar or pulsarctl."), nil
}

func (b *PulsarAdminTopicPolicyToolBuilder) handleRemoveTopicSubscriptionTypes(client cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return b.handleError("parse topic name", err), nil
	}

	// Check if the API supports subscription types management
	topicsClient := client.Topics()

	// Try to use reflection or type assertion to check if the method exists
	type SubscriptionTypesRemover interface {
		RemoveSubscriptionTypesEnabled(utils.TopicName) error
	}

	if remover, ok := topicsClient.(SubscriptionTypesRemover); ok {
		err := remover.RemoveSubscriptionTypesEnabled(*topicName)
		if err != nil {
			return b.handleError("remove topic subscription types policy", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Subscription types policy removed for topic %s (all types now allowed)", topicName.String())), nil
	}

	// Fallback: API not available in current version
	return mcp.NewToolResultError("Subscription types policy management is not available in the current pulsarctl API version. " +
		"This feature may require a newer version of Pulsar or pulsarctl."), nil
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

	// Parse time unit
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

	// Parse size unit
	var value float64
	var unit string
	var err error
	// /nolint:gocritic
	if strings.HasSuffix(retentionSize, "TB") || strings.HasSuffix(retentionSize, "T") {
		if strings.HasSuffix(retentionSize, "TB") {
			valueStr := retentionSize[:len(retentionSize)-2]
			value, err = strconv.ParseFloat(valueStr, 64)
			unit = "TB"
		} else {
			valueStr := retentionSize[:len(retentionSize)-1]
			value, err = strconv.ParseFloat(valueStr, 64)
			unit = "T"
		}
	} else if strings.HasSuffix(retentionSize, "GB") || strings.HasSuffix(retentionSize, "G") {
		if strings.HasSuffix(retentionSize, "GB") {
			valueStr := retentionSize[:len(retentionSize)-2]
			value, err = strconv.ParseFloat(valueStr, 64)
			unit = "GB"
		} else {
			valueStr := retentionSize[:len(retentionSize)-1]
			value, err = strconv.ParseFloat(valueStr, 64)
			unit = "G"
		}
	} else if strings.HasSuffix(retentionSize, "MB") || strings.HasSuffix(retentionSize, "M") {
		if strings.HasSuffix(retentionSize, "MB") {
			valueStr := retentionSize[:len(retentionSize)-2]
			value, err = strconv.ParseFloat(valueStr, 64)
			unit = "MB"
		} else {
			valueStr := retentionSize[:len(retentionSize)-1]
			value, err = strconv.ParseFloat(valueStr, 64)
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
