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
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// PulsarAdminAddSubscriptionTools adds subscription-related tools to the MCP server
func PulsarAdminAddSubscriptionTools(s *server.MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeaturePulsarAdminSubscriptions)) && !slices.Contains(features, string(FeatureAll)) && !slices.Contains(features, string(FeatureAllPulsar)) {
		return
	}

	// Add unified subscription management tool
	toolDesc := "Manage Apache Pulsar subscriptions on topics. " +
		"Subscriptions are named entities representing consumer groups that maintain their position in a topic. " +
		"Pulsar supports multiple subscription modes (Exclusive, Shared, Failover, Key_Shared) to accommodate different messaging patterns. " +
		"Each subscription tracks message acknowledgments independently, allowing multiple consumers to process messages at their own pace. " +
		"Subscriptions persist even when all consumers disconnect, maintaining state and preventing message loss. " +
		"Operations include listing, creating, deleting, and manipulating message cursors within subscriptions. " +
		"Most operations require namespace admin permissions plus produce/consume permissions on the topic."

	resourceDesc := "Resource to operate on. Available resources:\n" +
		"- subscription: A subscription on a topic representing a consumer group"

	operationDesc := "Operation to perform. Available operations:\n" +
		"- list: List all subscriptions for a topic\n" +
		"- create: Create a new subscription on a topic\n" +
		"- delete: Delete a subscription from a topic\n" +
		"- skip: Skip a specified number of messages for a subscription\n" +
		"- expire: Expire messages older than specified time for a subscription\n" +
		"- reset-cursor: Reset the cursor position for a subscription to a specific message ID"

	subscriptionTool := mcp.NewTool("pulsar_admin_subscription",
		mcp.WithDescription(toolDesc),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description(resourceDesc),
		),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
		),
		mcp.WithString("topic", mcp.Required(),
			mcp.Description("The fully qualified topic name in the format 'persistent://tenant/namespace/topic'. "+
				"For partitioned topics, you can either specify the base topic name (to apply the operation across all partitions) "+
				"or a specific partition in the format 'topicName-partition-N'."),
		),
		mcp.WithString("subscription",
			mcp.Description("The subscription name. Required for all operations except 'list'. "+
				"A subscription name is a logical identifier for a durable position in a topic. "+
				"Multiple consumers can attach to the same subscription to implement different messaging patterns."),
		),
		mcp.WithString("messageId",
			mcp.Description("Message ID for positioning the subscription cursor. Used in 'create' and 'reset-cursor' operations. "+
				"Values can be:\n"+
				"- 'latest': Position at the latest (most recent) message\n"+
				"- 'earliest': Position at the earliest (oldest available) message\n"+
				"- specific position in 'ledgerId:entryId' format for precise positioning"),
		),
		mcp.WithNumber("count",
			mcp.Description("The number of messages to skip (required for 'skip' operation). "+
				"This moves the subscription cursor forward by the specified number of messages without processing them."),
		),
		mcp.WithNumber("expireTimeInSeconds",
			mcp.Description("Expire messages older than the specified seconds (required for 'expire' operation). "+
				"This moves the subscription cursor to skip all messages published before the specified time."),
		),
		mcp.WithBoolean("force",
			mcp.Description("Force deletion of subscription (optional for 'delete' operation). "+
				"When true, all consumers will be forcefully disconnected and the subscription will be deleted. "+
				"Use with caution as it can interrupt active message processing."),
		),
	)
	s.AddTool(subscriptionTool, handleSubscriptionTool(readOnly))
}

// handleSubscriptionTool returns a function that handles subscription operations
func handleSubscriptionTool(readOnly bool) func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get required parameters
		resource, err := requiredParam[string](request.Params.Arguments, "resource")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get resource: %v", err)), nil
		}

		operation, err := requiredParam[string](request.Params.Arguments, "operation")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get operation: %v", err)), nil
		}

		topic, err := requiredParam[string](request.Params.Arguments, "topic")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic'. Please provide the fully qualified topic name: %v", err)), nil
		}

		// Normalize parameters
		resource = strings.ToLower(resource)
		operation = strings.ToLower(operation)

		// Validate write operations in read-only mode
		if readOnly && (operation != "list") {
			return mcp.NewToolResultError("Write operations are not allowed in read-only mode"), nil
		}

		// Verify resource type
		if resource != "subscription" {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resource: %s. Only 'subscription' is supported", resource)), nil
		}

		// Parse topic name
		topicName, err := utils.GetTopicName(topic)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
		}

		// Create the admin client
		admin, err := pulsar.GetAdminClient()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
		}

		// Dispatch based on operation
		switch operation {
		case "list":
			return handleSubsList(admin, topicName)
		case "create":
			return handleSubsCreate(admin, topicName, request)
		case "delete":
			return handleSubsDelete(admin, topicName, request)
		case "skip":
			return handleSubsSkip(admin, topicName, request)
		case "expire":
			return handleSubsExpire(admin, topicName, request)
		case "reset-cursor":
			return handleSubsResetCursor(admin, topicName, request)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid operation: %s. Available operations: list, create, delete, skip, expire, reset-cursor", operation)), nil
		}
	}
}

// handleSubsList handles listing all subscriptions for a topic
func handleSubsList(admin cmdutils.Client, topicName *utils.TopicName) (*mcp.CallToolResult, error) {
	// List subscriptions
	subscriptions, err := admin.Subscriptions().List(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list subscriptions for topic '%s': %v",
			topicName.String(), err)), nil
	}

	// Format the output
	jsonBytes, err := json.Marshal(subscriptions)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to process subscription data: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleSubsCreate handles creating a new subscription
func handleSubsCreate(admin cmdutils.Client, topicName *utils.TopicName, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameter
	subscription, err := requiredParam[string](request.Params.Arguments, "subscription")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'subscription' for subscription.create: %v", err)), nil
	}

	// Get optional messageID parameter (default is "latest")
	messageID, hasMessageID := optionalParam[string](request.Params.Arguments, "messageId")
	if !hasMessageID {
		messageID = "latest"
	}

	// Parse messageId
	var messageIDObj utils.MessageID
	switch messageID {
	case "latest":
		messageIDObj = utils.Latest
	case "earliest":
		messageIDObj = utils.Earliest
	default:
		s := strings.Split(messageID, ":")
		if len(s) != 2 {
			return mcp.NewToolResultError(fmt.Sprintf(
				"Invalid messageId format: %s. Use 'latest', 'earliest', or 'ledgerId:entryId' format", messageID)), nil
		}
		msgID, err := utils.ParseMessageID(messageID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to parse messageId '%s': %v", messageID, err)), nil
		}
		messageIDObj = *msgID
	}

	// Create subscription
	err = admin.Subscriptions().Create(*topicName, subscription, messageIDObj)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create subscription '%s' on topic '%s': %v",
			subscription, topicName.String(), err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Created subscription '%s' on topic '%s' from position '%s' successfully",
		subscription, topicName.String(), messageID)), nil
}

// handleSubsDelete handles deleting a subscription
func handleSubsDelete(admin cmdutils.Client, topicName *utils.TopicName, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameter
	subscription, err := requiredParam[string](request.Params.Arguments, "subscription")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'subscription' for subscription.delete: %v", err)), nil
	}

	// Get optional force parameter (default is false)
	force, hasForce := optionalParam[bool](request.Params.Arguments, "force")
	if !hasForce {
		force = false
	}

	// Delete subscription
	if force {
		err = admin.Subscriptions().ForceDelete(*topicName, subscription)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to forcefully delete subscription '%s' from topic '%s': %v",
				subscription, topicName.String(), err)), nil
		}
	} else {
		err = admin.Subscriptions().Delete(*topicName, subscription)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete subscription '%s' from topic '%s': %v",
				subscription, topicName.String(), err)), nil
		}
	}

	forceStr := ""
	if force {
		forceStr = " forcefully"
	}
	return mcp.NewToolResultText(fmt.Sprintf("Deleted subscription '%s' from topic '%s'%s successfully",
		subscription, topicName.String(), forceStr)), nil
}

// handleSubsSkip handles skipping messages for a subscription
func handleSubsSkip(admin cmdutils.Client, topicName *utils.TopicName, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	subscription, err := requiredParam[string](request.Params.Arguments, "subscription")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'subscription' for subscription.skip: %v", err)), nil
	}

	count, err := requiredParam[float64](request.Params.Arguments, "count")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'count' for subscription.skip: %v", err)), nil
	}

	// Skip messages
	err = admin.Subscriptions().SkipMessages(*topicName, subscription, int64(count))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to skip messages for subscription '%s' on topic '%s': %v",
			subscription, topicName.String(), err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Skipped %d messages for subscription '%s' on topic '%s' successfully",
		int(count), subscription, topicName.String())), nil
}

// handleSubsExpire handles expiring messages for a subscription
func handleSubsExpire(admin cmdutils.Client, topicName *utils.TopicName, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	subscription, err := requiredParam[string](request.Params.Arguments, "subscription")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'subscription' for subscription.expire: %v", err)), nil
	}

	expireTime, err := requiredParam[float64](request.Params.Arguments, "expireTimeInSeconds")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'expireTimeInSeconds' for subscription.expire: %v", err)), nil
	}

	// Expire messages
	err = admin.Subscriptions().ExpireMessages(*topicName, subscription, int64(expireTime))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to expire messages for subscription '%s' on topic '%s': %v",
			subscription, topicName.String(), err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Expired messages older than %d seconds for subscription '%s' on topic '%s' successfully",
			int(expireTime), subscription, topicName.String()),
	), nil
}

// handleSubsResetCursor handles resetting a subscription cursor
func handleSubsResetCursor(admin cmdutils.Client, topicName *utils.TopicName, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	subscription, err := requiredParam[string](request.Params.Arguments, "subscription")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'subscription' for subscription.reset-cursor: %v", err)), nil
	}

	messageID, err := requiredParam[string](request.Params.Arguments, "messageId")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'messageId' for subscription.reset-cursor: %v", err)), nil
	}

	// Parse messageId
	var messageIDObj utils.MessageID
	switch messageID {
	case "latest":
		messageIDObj = utils.Latest
	case "earliest":
		messageIDObj = utils.Earliest
	default:
		s := strings.Split(messageID, ":")
		if len(s) != 2 {
			return mcp.NewToolResultError(fmt.Sprintf(
				"Invalid messageId format: %s. Use 'latest', 'earliest', or 'ledgerId:entryId' format", messageID)), nil
		}
		msgID, err := utils.ParseMessageID(messageID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to parse messageId '%s': %v", messageID, err)), nil
		}
		messageIDObj = *msgID
	}

	// Reset cursor
	err = admin.Subscriptions().ResetCursorToMessageID(*topicName, subscription, messageIDObj)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to reset cursor for subscription '%s' on topic '%s': %v",
			subscription, topicName.String(), err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Reset cursor for subscription '%s' on topic '%s' to position '%s' successfully",
			subscription, topicName.String(), messageID),
	), nil
}
