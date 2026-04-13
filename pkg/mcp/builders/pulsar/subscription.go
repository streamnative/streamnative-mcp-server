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
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

var supportedSubscriptionOperations = map[string]struct{}{
	"list":              {},
	"create":            {},
	"delete":            {},
	"skip":              {},
	"expire":            {},
	"reset-cursor":      {},
	"peek":              {},
	"get-message-by-id": {},
}

var readOnlyRestrictedSubscriptionOperations = map[string]struct{}{
	"create":       {},
	"delete":       {},
	"skip":         {},
	"expire":       {},
	"reset-cursor": {},
}

const maxSubscriptionPeekCount int64 = 100

type subscriptionMessageData struct {
	Topic         string            `json:"topic,omitempty"`
	MessageID     string            `json:"messageId"`
	Properties    map[string]string `json:"properties,omitempty"`
	Payload       string            `json:"payload"`
	PayloadBase64 string            `json:"payloadBase64"`
	PayloadHex    string            `json:"payloadHex"`
}

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
func (b *PulsarAdminSubscriptionToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
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

	return []server.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildSubscriptionTool builds the Pulsar Admin Subscription MCP tool definition
// Migrated from the original tool definition logic
func (b *PulsarAdminSubscriptionToolBuilder) buildSubscriptionTool() mcp.Tool {
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
		"- reset-cursor: Reset the cursor position for a subscription to a specific message ID\n" +
		"- peek: Peek one or more messages for a subscription without advancing the cursor\n" +
		"- get-message-by-id: Read a message by ledger ID and entry ID for topic-level debugging"

	return mcp.NewTool("pulsar_admin_subscription",
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
			mcp.Description("The subscription name. Required for all operations except 'list' and 'get-message-by-id'. "+
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
			mcp.Description("The number of messages to skip (required for 'skip' operation) or peek (optional for 'peek', default 1). "+
				"For 'skip', this moves the subscription cursor forward by the specified number of messages without processing them. "+
				"For 'peek', this limits how many messages are returned without moving the cursor. Maximum: 100."),
		),
		mcp.WithNumber("expireTimeInSeconds",
			mcp.Description("Expire messages older than the specified seconds (required for 'expire' operation). "+
				"This moves the subscription cursor to skip all messages published before the specified time."),
		),
		mcp.WithNumber("ledgerId",
			mcp.Description("Ledger ID used by 'get-message-by-id'. Must be a non-negative integer."),
		),
		mcp.WithNumber("entryId",
			mcp.Description("Entry ID used by 'get-message-by-id'. Must be a non-negative integer."),
		),
		mcp.WithBoolean("force",
			mcp.Description("Force deletion of subscription (optional for 'delete' operation). "+
				"When true, all consumers will be forcefully disconnected and the subscription will be deleted. "+
				"Use with caution as it can interrupt active message processing."),
		),
	)
}

// buildSubscriptionHandler builds the Pulsar Admin Subscription handler function
// Migrated from the original handler logic
func (b *PulsarAdminSubscriptionToolBuilder) buildSubscriptionHandler(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		topic, err := request.RequireString("topic")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic'. Please provide the fully qualified topic name: %v", err)), nil
		}

		// Normalize parameters
		resource = strings.ToLower(resource)
		operation = strings.ToLower(operation)

		if !isSupportedSubscriptionOperation(operation) {
			return mcp.NewToolResultError(fmt.Sprintf("Unknown operation: %s", operation)), nil
		}

		// Validate write operations in read-only mode
		if readOnly && isReadOnlyRestrictedSubscriptionOperation(operation) {
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
		case "peek":
			return b.handleSubsPeek(admin, topicName, request)
		case "get-message-by-id":
			return b.handleSubsGetMessageByID(admin, topicName, request)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Unknown operation: %s", operation)), nil
		}
	}
}

// Unified error handling and utility functions

// handleError provides unified error handling
func (b *PulsarAdminSubscriptionToolBuilder) handleError(operation string, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("Failed to %s: %v", operation, err))
}

// marshalResponse provides unified JSON serialization for responses
func (b *PulsarAdminSubscriptionToolBuilder) marshalResponse(data interface{}) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return b.handleError("marshal response", err), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// Operation handler functions - migrated from the original implementation

// handleSubsList handles listing all subscriptions for a topic
func (b *PulsarAdminSubscriptionToolBuilder) handleSubsList(admin cmdutils.Client, topicName *utils.TopicName) (*mcp.CallToolResult, error) {
	// List subscriptions
	subscriptions, err := admin.Subscriptions().List(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list subscriptions for topic '%s': %v",
			topicName.String(), err)), nil
	}

	return b.marshalResponse(subscriptions)
}

// handleSubsCreate handles creating a new subscription
func (b *PulsarAdminSubscriptionToolBuilder) handleSubsCreate(admin cmdutils.Client, topicName *utils.TopicName, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameter
	subscription, err := request.RequireString("subscription")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'subscription' for subscription.create: %v", err)), nil
	}

	// Get optional messageID parameter (default is "latest")
	messageID := request.GetString("messageId", "latest")

	// Parse messageId
	messageIDObj, err := parseSubscriptionMessageID(messageID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
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
func (b *PulsarAdminSubscriptionToolBuilder) handleSubsDelete(admin cmdutils.Client, topicName *utils.TopicName, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameter
	subscription, err := request.RequireString("subscription")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'subscription' for subscription.delete: %v", err)), nil
	}

	// Get optional force parameter (default is false)
	force := request.GetBool("force", false)

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
func (b *PulsarAdminSubscriptionToolBuilder) handleSubsSkip(admin cmdutils.Client, topicName *utils.TopicName, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	subscription, err := request.RequireString("subscription")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'subscription' for subscription.skip: %v", err)), nil
	}

	count, err := requireInt64Arg(request, "count")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing or invalid required parameter 'count' for subscription.skip: %v", err)), nil
	}

	// Skip messages
	err = admin.Subscriptions().SkipMessages(*topicName, subscription, count)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to skip messages for subscription '%s' on topic '%s': %v",
			subscription, topicName.String(), err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Skipped %d messages for subscription '%s' on topic '%s' successfully",
		count, subscription, topicName.String())), nil
}

// handleSubsExpire handles expiring messages for a subscription
func (b *PulsarAdminSubscriptionToolBuilder) handleSubsExpire(admin cmdutils.Client, topicName *utils.TopicName, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	subscription, err := request.RequireString("subscription")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'subscription' for subscription.expire: %v", err)), nil
	}

	expireTime, err := requireInt64Arg(request, "expireTimeInSeconds")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing or invalid required parameter 'expireTimeInSeconds' for subscription.expire: %v", err)), nil
	}

	// Expire messages
	err = admin.Subscriptions().ExpireMessages(*topicName, subscription, expireTime)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to expire messages for subscription '%s' on topic '%s': %v",
			subscription, topicName.String(), err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Expired messages older than %d seconds for subscription '%s' on topic '%s' successfully",
			expireTime, subscription, topicName.String()),
	), nil
}

// handleSubsResetCursor handles resetting a subscription cursor
func (b *PulsarAdminSubscriptionToolBuilder) handleSubsResetCursor(admin cmdutils.Client, topicName *utils.TopicName, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	subscription, err := request.RequireString("subscription")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'subscription' for subscription.reset-cursor: %v", err)), nil
	}

	messageID, err := request.RequireString("messageId")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'messageId' for subscription.reset-cursor: %v", err)), nil
	}

	// Parse messageId
	messageIDObj, err := parseSubscriptionMessageID(messageID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
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

// handleSubsPeek handles peeking messages for a subscription without advancing the cursor.
func (b *PulsarAdminSubscriptionToolBuilder) handleSubsPeek(admin cmdutils.Client, topicName *utils.TopicName, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subscription, err := request.RequireString("subscription")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'subscription' for subscription.peek: %v", err)), nil
	}

	if topicName.GetDomain().String() != "persistent" {
		return mcp.NewToolResultError("The specified topic name is not a persistent topic"), nil
	}

	count, err := getInt64Arg(request, "count", 1)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid parameter 'count' for subscription.peek: %v", err)), nil
	}
	if err := validateSubscriptionPeekCount(count); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	messages, err := admin.Subscriptions().PeekMessages(*topicName, subscription, int(count))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to peek messages for subscription '%s' on topic '%s': %v",
			subscription, topicName.String(), err)), nil
	}

	result := struct {
		Topic        string                    `json:"topic"`
		Subscription string                    `json:"subscription"`
		Count        int                       `json:"count"`
		Messages     []subscriptionMessageData `json:"messages"`
	}{
		Topic:        topicName.String(),
		Subscription: subscription,
		Count:        len(messages),
		Messages:     buildSubscriptionMessageData(messages),
	}

	return b.marshalResponse(result)
}

// handleSubsGetMessageByID handles reading a message by ledger and entry ID.
func (b *PulsarAdminSubscriptionToolBuilder) handleSubsGetMessageByID(admin cmdutils.Client, topicName *utils.TopicName, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ledgerID, err := requireInt64Arg(request, "ledgerId")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing or invalid required parameter 'ledgerId' for subscription.get-message-by-id: %v", err)), nil
	}

	entryID, err := requireInt64Arg(request, "entryId")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing or invalid required parameter 'entryId' for subscription.get-message-by-id: %v", err)), nil
	}
	if err := validateSubscriptionMessageLookupIDs(ledgerID, entryID); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	messages, err := admin.Subscriptions().GetMessagesByID(*topicName, ledgerID, entryID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get message by ledger ID %d and entry ID %d on topic '%s': %v",
			ledgerID, entryID, topicName.String(), err)), nil
	}
	if len(messages) == 0 {
		return mcp.NewToolResultError(fmt.Sprintf(
			"No message found for ledger ID %d and entry ID %d on topic '%s'", ledgerID, entryID, topicName.String(),
		)), nil
	}

	result := struct {
		Topic    string                  `json:"topic"`
		LedgerID int64                   `json:"ledgerId"`
		EntryID  int64                   `json:"entryId"`
		Message  subscriptionMessageData `json:"message"`
	}{
		Topic:    topicName.String(),
		LedgerID: ledgerID,
		EntryID:  entryID,
		Message:  newSubscriptionMessageData(messages[0]),
	}

	return b.marshalResponse(result)
}

func isSupportedSubscriptionOperation(operation string) bool {
	_, ok := supportedSubscriptionOperations[strings.ToLower(operation)]
	return ok
}

func isReadOnlyRestrictedSubscriptionOperation(operation string) bool {
	_, ok := readOnlyRestrictedSubscriptionOperations[strings.ToLower(operation)]
	return ok
}

func parseSubscriptionMessageID(messageID string) (utils.MessageID, error) {
	switch messageID {
	case "latest":
		return utils.Latest, nil
	case "earliest":
		return utils.Earliest, nil
	default:
		parts := strings.Split(messageID, ":")
		if len(parts) != 2 {
			return utils.MessageID{}, fmt.Errorf(
				"invalid messageId format: %s. Use 'latest', 'earliest', or 'ledgerId:entryId' format", messageID,
			)
		}

		msgID, err := utils.ParseMessageID(messageID)
		if err != nil {
			return utils.MessageID{}, fmt.Errorf("failed to parse messageId '%s': %w", messageID, err)
		}
		return *msgID, nil
	}
}

func requireInt64Arg(request mcp.CallToolRequest, key string) (int64, error) {
	value, err := request.RequireFloat(key)
	if err != nil {
		return 0, err
	}
	if math.Trunc(value) != value {
		return 0, fmt.Errorf("parameter '%s' must be an integer", key)
	}
	return int64(value), nil
}

func getInt64Arg(request mcp.CallToolRequest, key string, defaultValue int64) (int64, error) {
	value := request.GetFloat(key, float64(defaultValue))
	if math.Trunc(value) != value {
		return 0, fmt.Errorf("parameter '%s' must be an integer", key)
	}
	return int64(value), nil
}

func validateSubscriptionPeekCount(count int64) error {
	if count <= 0 {
		return fmt.Errorf("parameter 'count' for subscription.peek must be greater than 0")
	}
	if count > int64(maxPlatformInt()) {
		return fmt.Errorf("parameter 'count' for subscription.peek exceeds the platform integer limit")
	}
	if count > maxSubscriptionPeekCount {
		return fmt.Errorf("parameter 'count' for subscription.peek must be less than or equal to %d", maxSubscriptionPeekCount)
	}
	return nil
}

func validateSubscriptionMessageLookupIDs(ledgerID, entryID int64) error {
	if ledgerID < 0 {
		return fmt.Errorf("parameter 'ledgerId' for subscription.get-message-by-id must be greater than or equal to 0")
	}
	if entryID < 0 {
		return fmt.Errorf("parameter 'entryId' for subscription.get-message-by-id must be greater than or equal to 0")
	}
	return nil
}

func maxPlatformInt() int {
	return int(^uint(0) >> 1)
}

func buildSubscriptionMessageData(messages []*utils.Message) []subscriptionMessageData {
	result := make([]subscriptionMessageData, 0, len(messages))
	for _, message := range messages {
		result = append(result, newSubscriptionMessageData(message))
	}
	return result
}

func newSubscriptionMessageData(message *utils.Message) subscriptionMessageData {
	if message == nil {
		return subscriptionMessageData{}
	}

	return subscriptionMessageData{
		Topic:         message.Topic,
		MessageID:     message.GetMessageID().String(),
		Properties:    message.GetProperties(),
		Payload:       string(message.GetPayload()),
		PayloadBase64: base64.StdEncoding.EncodeToString(message.GetPayload()),
		PayloadHex:    hex.Dump(message.GetPayload()),
	}
}
