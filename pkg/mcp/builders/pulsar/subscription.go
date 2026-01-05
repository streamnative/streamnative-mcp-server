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
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/adapter"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

// PulsarAdminSubscriptionToolBuilder implements the ToolBuilder interface for Pulsar Admin Subscription tools
// It provides functionality to build Pulsar subscription management tools
// /nolint:revive
type PulsarAdminSubscriptionToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminSubscriptionToolBuilder creates a new Pulsar Admin Subscription tool builder instance
func NewPulsarAdminSubscriptionToolBuilder() *PulsarAdminSubscriptionToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_subscription",
		Version:     "1.0.0",
		Description: "Pulsar Admin subscription management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "subscription", "admin"},
	}

	features := []string{
		"pulsar-admin-subscriptions",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminSubscriptionToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar Admin Subscription tool list
// This is the core method implementing the ToolBuilder interface
func (b *PulsarAdminSubscriptionToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ServerTool, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool := b.buildSubscriptionTool()
	handler := b.buildSubscriptionHandler(config.ReadOnly)

	return []builders.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildSubscriptionTool builds the Pulsar Admin Subscription MCP tool definition
// Migrated from the original tool definition logic
func (b *PulsarAdminSubscriptionToolBuilder) buildSubscriptionTool() *mcpsdk.Tool {
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

	return builders.NewTool("pulsar_admin_subscription",
		builders.WithDescription(toolDesc),
		builders.WithString("resource", builders.Required(),
			builders.Description(resourceDesc),
		),
		builders.WithString("operation", builders.Required(),
			builders.Description(operationDesc),
		),
		builders.WithString("topic", builders.Required(),
			builders.Description("The fully qualified topic name in the format 'persistent://tenant/namespace/topic'. "+
				"For partitioned topics, you can either specify the base topic name (to apply the operation across all partitions) "+
				"or a specific partition in the format 'topicName-partition-N'."),
		),
		builders.WithString("subscription",
			builders.Description("The subscription name. Required for all operations except 'list'. "+
				"A subscription name is a logical identifier for a durable position in a topic. "+
				"Multiple consumers can attach to the same subscription to implement different messaging patterns."),
		),
		builders.WithString("messageId",
			builders.Description("Message ID for positioning the subscription cursor. Used in 'create' and 'reset-cursor' operations. "+
				"Values can be:\n"+
				"- 'latest': Position at the latest (most recent) message\n"+
				"- 'earliest': Position at the earliest (oldest available) message\n"+
				"- specific position in 'ledgerId:entryId' format for precise positioning"),
		),
		builders.WithNumber("count",
			builders.Description("The number of messages to skip (required for 'skip' operation). "+
				"This moves the subscription cursor forward by the specified number of messages without processing them."),
		),
		builders.WithNumber("expireTimeInSeconds",
			builders.Description("Expire messages older than the specified seconds (required for 'expire' operation). "+
				"This moves the subscription cursor to skip all messages published before the specified time."),
		),
		builders.WithBoolean("force",
			builders.Description("Force deletion of subscription (optional for 'delete' operation). "+
				"When true, all consumers will be forcefully disconnected and the subscription will be deleted. "+
				"Use with caution as it can interrupt active message processing."),
		),
	)
}

// buildSubscriptionHandler builds the Pulsar Admin Subscription handler function
// Migrated from the original handler logic
func (b *PulsarAdminSubscriptionToolBuilder) buildSubscriptionHandler(readOnly bool) func(context.Context, *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	return func(ctx context.Context, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
		// Get required parameters
		resource, err := adapter.RequireString(request, "resource")
		if err != nil {
			return adapter.NewErrorResult("Failed to get resource: %v", err), nil
		}

		operation, err := adapter.RequireString(request, "operation")
		if err != nil {
			return adapter.NewErrorResult("Failed to get operation: %v", err), nil
		}

		topic, err := adapter.RequireString(request, "topic")
		if err != nil {
			return adapter.NewErrorResult("Missing required parameter 'topic'. Please provide the fully qualified topic name: %v", err), nil
		}

		// Normalize parameters
		resource = strings.ToLower(resource)
		operation = strings.ToLower(operation)

		// Validate write operations in read-only mode
		if readOnly && (operation != "list") {
			return adapter.NewErrorResult("Write operations are not allowed in read-only mode"), nil
		}

		// Verify resource type
		if resource != "subscription" {
			return adapter.NewErrorResult("Invalid resource: %s. Only 'subscription' is supported", resource), nil
		}

		// Parse topic name
		topicName, err := utils.GetTopicName(topic)
		if err != nil {
			return adapter.NewErrorResult("Invalid topic name '%s': %v", topic, err), nil
		}

		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return adapter.NewErrorResult("Pulsar session not found in context"), nil
		}

		// Create the admin client
		admin, err := session.GetAdminClient()
		if err != nil {
			return adapter.NewErrorResult("Failed to get admin client: %v", err), nil
		}

		// Dispatch based on operation
		switch operation {
		case "list":
			return b.handleSubsList(admin, topicName)
		case "create":
			return b.handleSubsCreate(admin, topicName, request)
		case "delete":
			return b.handleSubsDelete(admin, topicName, request)
		case "skip":
			return b.handleSubsSkip(admin, topicName, request)
		case "expire":
			return b.handleSubsExpire(admin, topicName, request)
		case "reset-cursor":
			return b.handleSubsResetCursor(admin, topicName, request)
		default:
			return adapter.NewErrorResult("Unknown operation: %s", operation), nil
		}
	}
}

// Unified error handling and utility functions

// handleError provides unified error handling
func (b *PulsarAdminSubscriptionToolBuilder) handleError(operation string, err error) *mcpsdk.CallToolResult {
	return adapter.NewErrorResult("Failed to %s: %v", operation, err)
}

// marshalResponse provides unified JSON serialization for responses
func (b *PulsarAdminSubscriptionToolBuilder) marshalResponse(data interface{}) (*mcpsdk.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return b.handleError("marshal response", err), nil
	}
	return adapter.NewTextResult(string(jsonBytes)), nil
}

// Operation handler functions - migrated from the original implementation

// handleSubsList handles listing all subscriptions for a topic
func (b *PulsarAdminSubscriptionToolBuilder) handleSubsList(admin cmdutils.Client, topicName *utils.TopicName) (*mcpsdk.CallToolResult, error) {
	// List subscriptions
	subscriptions, err := admin.Subscriptions().List(*topicName)
	if err != nil {
		return adapter.NewErrorResult(fmt.Sprintf("Failed to list subscriptions for topic '%s': %v",
			topicName.String(), err)), nil
	}

	return b.marshalResponse(subscriptions)
}

// handleSubsCreate handles creating a new subscription
func (b *PulsarAdminSubscriptionToolBuilder) handleSubsCreate(admin cmdutils.Client, topicName *utils.TopicName, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	// Get required parameter
	subscription, err := adapter.RequireString(request, "subscription")
	if err != nil {
		return adapter.NewErrorResult("Missing required parameter 'subscription' for subscription.create: %v", err), nil
	}

	// Get optional messageID parameter (default is "latest")
	messageID := adapter.GetString(request, "messageId",  "latest")

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
			return adapter.NewErrorResult(fmt.Sprintf(
				"Invalid messageId format: %s. Use 'latest', 'earliest', or 'ledgerId:entryId' format", messageID)), nil
		}
		msgID, err := utils.ParseMessageID(messageID)
		if err != nil {
			return adapter.NewErrorResult("Failed to parse messageId '%s': %v", messageID, err), nil
		}
		messageIDObj = *msgID
	}

	// Create subscription
	err = admin.Subscriptions().Create(*topicName, subscription, messageIDObj)
	if err != nil {
		return adapter.NewErrorResult(fmt.Sprintf("Failed to create subscription '%s' on topic '%s': %v",
			subscription, topicName.String(), err)), nil
	}

	return adapter.NewTextResult(fmt.Sprintf("Created subscription '%s' on topic '%s' from position '%s' successfully",
		subscription, topicName.String(), messageID)), nil
}

// handleSubsDelete handles deleting a subscription
func (b *PulsarAdminSubscriptionToolBuilder) handleSubsDelete(admin cmdutils.Client, topicName *utils.TopicName, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	// Get required parameter
	subscription, err := adapter.RequireString(request, "subscription")
	if err != nil {
		return adapter.NewErrorResult("Missing required parameter 'subscription' for subscription.delete: %v", err), nil
	}

	// Get optional force parameter (default is false)
	force := adapter.GetBool(request, "force",  false)

	// Delete subscription
	if force {
		err = admin.Subscriptions().ForceDelete(*topicName, subscription)
		if err != nil {
			return adapter.NewErrorResult(fmt.Sprintf("Failed to forcefully delete subscription '%s' from topic '%s': %v",
				subscription, topicName.String(), err)), nil
		}
	} else {
		err = admin.Subscriptions().Delete(*topicName, subscription)
		if err != nil {
			return adapter.NewErrorResult(fmt.Sprintf("Failed to delete subscription '%s' from topic '%s': %v",
				subscription, topicName.String(), err)), nil
		}
	}

	forceStr := ""
	if force {
		forceStr = " forcefully"
	}
	return adapter.NewTextResult(fmt.Sprintf("Deleted subscription '%s' from topic '%s'%s successfully",
		subscription, topicName.String(), forceStr)), nil
}

// handleSubsSkip handles skipping messages for a subscription
func (b *PulsarAdminSubscriptionToolBuilder) handleSubsSkip(admin cmdutils.Client, topicName *utils.TopicName, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	// Get required parameters
	subscription, err := adapter.RequireString(request, "subscription")
	if err != nil {
		return adapter.NewErrorResult("Missing required parameter 'subscription' for subscription.skip: %v", err), nil
	}

	count, err := adapter.RequireFloat(request, "count")
	if err != nil {
		return adapter.NewErrorResult("Missing required parameter 'count' for subscription.skip: %v", err), nil
	}

	// Skip messages
	err = admin.Subscriptions().SkipMessages(*topicName, subscription, int64(count))
	if err != nil {
		return adapter.NewErrorResult(fmt.Sprintf("Failed to skip messages for subscription '%s' on topic '%s': %v",
			subscription, topicName.String(), err)), nil
	}

	return adapter.NewTextResult(fmt.Sprintf("Skipped %d messages for subscription '%s' on topic '%s' successfully",
		int(count), subscription, topicName.String())), nil
}

// handleSubsExpire handles expiring messages for a subscription
func (b *PulsarAdminSubscriptionToolBuilder) handleSubsExpire(admin cmdutils.Client, topicName *utils.TopicName, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	// Get required parameters
	subscription, err := adapter.RequireString(request, "subscription")
	if err != nil {
		return adapter.NewErrorResult("Missing required parameter 'subscription' for subscription.expire: %v", err), nil
	}

	expireTime, err := adapter.RequireFloat(request, "expireTimeInSeconds")
	if err != nil {
		return adapter.NewErrorResult("Missing required parameter 'expireTimeInSeconds' for subscription.expire: %v", err), nil
	}

	// Expire messages
	err = admin.Subscriptions().ExpireMessages(*topicName, subscription, int64(expireTime))
	if err != nil {
		return adapter.NewErrorResult(fmt.Sprintf("Failed to expire messages for subscription '%s' on topic '%s': %v",
			subscription, topicName.String(), err)), nil
	}

	return adapter.NewTextResult(
		fmt.Sprintf("Expired messages older than %d seconds for subscription '%s' on topic '%s' successfully",
			int(expireTime), subscription, topicName.String()),
	), nil
}

// handleSubsResetCursor handles resetting a subscription cursor
func (b *PulsarAdminSubscriptionToolBuilder) handleSubsResetCursor(admin cmdutils.Client, topicName *utils.TopicName, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	// Get required parameters
	subscription, err := adapter.RequireString(request, "subscription")
	if err != nil {
		return adapter.NewErrorResult("Missing required parameter 'subscription' for subscription.reset-cursor: %v", err), nil
	}

	messageID, err := adapter.RequireString(request, "messageId")
	if err != nil {
		return adapter.NewErrorResult("Missing required parameter 'messageId' for subscription.reset-cursor: %v", err), nil
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
			return adapter.NewErrorResult(fmt.Sprintf(
				"Invalid messageId format: %s. Use 'latest', 'earliest', or 'ledgerId:entryId' format", messageID)), nil
		}
		msgID, err := utils.ParseMessageID(messageID)
		if err != nil {
			return adapter.NewErrorResult("Failed to parse messageId '%s': %v", messageID, err), nil
		}
		messageIDObj = *msgID
	}

	// Reset cursor
	err = admin.Subscriptions().ResetCursorToMessageID(*topicName, subscription, messageIDObj)
	if err != nil {
		return adapter.NewErrorResult(fmt.Sprintf("Failed to reset cursor for subscription '%s' on topic '%s': %v",
			subscription, topicName.String(), err)), nil
	}

	return adapter.NewTextResult(
		fmt.Sprintf("Reset cursor for subscription '%s' on topic '%s' to position '%s' successfully",
			subscription, topicName.String(), messageID),
	), nil
}
