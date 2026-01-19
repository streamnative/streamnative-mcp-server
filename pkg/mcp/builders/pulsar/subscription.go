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
	"github.com/google/jsonschema-go/jsonschema"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

type pulsarAdminSubscriptionInput struct {
	Resource            string   `json:"resource"`
	Operation           string   `json:"operation"`
	Topic               string   `json:"topic"`
	Subscription        *string  `json:"subscription,omitempty"`
	MessageID           *string  `json:"messageId,omitempty"`
	Count               *float64 `json:"count,omitempty"`
	ExpireTimeInSeconds *float64 `json:"expireTimeInSeconds,omitempty"`
	Force               *bool    `json:"force,omitempty"`
}

const (
	pulsarAdminSubscriptionResourceDesc = "Resource to operate on. Available resources:\n" +
		"- subscription: A subscription on a topic representing a consumer group"
	pulsarAdminSubscriptionOperationDesc = "Operation to perform. Available operations:\n" +
		"- list: List all subscriptions for a topic\n" +
		"- create: Create a new subscription on a topic\n" +
		"- delete: Delete a subscription from a topic\n" +
		"- skip: Skip a specified number of messages for a subscription\n" +
		"- expire: Expire messages older than specified time for a subscription\n" +
		"- reset-cursor: Reset the cursor position for a subscription to a specific message ID"
	pulsarAdminSubscriptionTopicDesc = "The fully qualified topic name in the format 'persistent://tenant/namespace/topic'. " +
		"For partitioned topics, you can either specify the base topic name (to apply the operation across all partitions) " +
		"or a specific partition in the format 'topicName-partition-N'."
	pulsarAdminSubscriptionNameDesc = "The subscription name. Required for all operations except 'list'. " +
		"A subscription name is a logical identifier for a durable position in a topic. " +
		"Multiple consumers can attach to the same subscription to implement different messaging patterns."
	pulsarAdminSubscriptionMessageIDDesc = "Message ID for positioning the subscription cursor. Used in 'create' and 'reset-cursor' operations. " +
		"Values can be:\n" +
		"- 'latest': Position at the latest (most recent) message\n" +
		"- 'earliest': Position at the earliest (oldest available) message\n" +
		"- specific position in 'ledgerId:entryId' format for precise positioning"
	pulsarAdminSubscriptionCountDesc = "The number of messages to skip (required for 'skip' operation). " +
		"This moves the subscription cursor forward by the specified number of messages without processing them."
	pulsarAdminSubscriptionExpireDesc = "Expire messages older than the specified seconds (required for 'expire' operation). " +
		"This moves the subscription cursor to skip all messages published before the specified time."
	pulsarAdminSubscriptionForceDesc = "Force deletion of subscription (optional for 'delete' operation). " +
		"When true, all consumers will be forcefully disconnected and the subscription will be deleted. " +
		"Use with caution as it can interrupt active message processing."
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
func (b *PulsarAdminSubscriptionToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool, err := b.buildSubscriptionTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildSubscriptionHandler(config.ReadOnly)

	return []builders.ToolDefinition{
		builders.ServerTool[pulsarAdminSubscriptionInput, any]{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildSubscriptionTool builds the Pulsar Admin Subscription MCP tool definition
// Migrated from the original tool definition logic
func (b *PulsarAdminSubscriptionToolBuilder) buildSubscriptionTool() (*sdk.Tool, error) {
	inputSchema, err := buildPulsarAdminSubscriptionInputSchema()
	if err != nil {
		return nil, err
	}

	toolDesc := "Manage Apache Pulsar subscriptions on topics. " +
		"Subscriptions are named entities representing consumer groups that maintain their position in a topic. " +
		"Pulsar supports multiple subscription modes (Exclusive, Shared, Failover, Key_Shared) to accommodate different messaging patterns. " +
		"Each subscription tracks message acknowledgments independently, allowing multiple consumers to process messages at their own pace. " +
		"Subscriptions persist even when all consumers disconnect, maintaining state and preventing message loss. " +
		"Operations include listing, creating, deleting, and manipulating message cursors within subscriptions. " +
		"Most operations require namespace admin permissions plus produce/consume permissions on the topic."

	return &sdk.Tool{
		Name:        "pulsar_admin_subscription",
		Description: toolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildSubscriptionHandler builds the Pulsar Admin Subscription handler function
// Migrated from the original handler logic
func (b *PulsarAdminSubscriptionToolBuilder) buildSubscriptionHandler(readOnly bool) builders.ToolHandlerFunc[pulsarAdminSubscriptionInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input pulsarAdminSubscriptionInput) (*sdk.CallToolResult, any, error) {
		resource := input.Resource
		if resource == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'resource'")
		}

		operation := input.Operation
		if operation == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'operation'")
		}

		topic := input.Topic
		if topic == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'topic'. please provide the fully qualified topic name")
		}

		// Normalize parameters
		resource = strings.ToLower(resource)
		operation = strings.ToLower(operation)

		// Validate write operations in read-only mode
		if readOnly && (operation != "list") {
			return nil, nil, fmt.Errorf("write operations are not allowed in read-only mode")
		}

		// Verify resource type
		if resource != "subscription" {
			return nil, nil, fmt.Errorf("invalid resource: %s. only 'subscription' is supported", resource)
		}

		// Parse topic name
		topicName, err := utils.GetTopicName(topic)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid topic name '%s': %v", topic, err)
		}

		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return nil, nil, fmt.Errorf("pulsar session not found in context")
		}

		// Create the admin client
		admin, err := session.GetAdminClient()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get admin client: %v", err)
		}

		// Dispatch based on operation
		switch operation {
		case "list":
			result, err := b.handleSubsList(admin, topicName)
			return result, nil, err
		case "create":
			result, err := b.handleSubsCreate(admin, topicName, input)
			return result, nil, err
		case "delete":
			result, err := b.handleSubsDelete(admin, topicName, input)
			return result, nil, err
		case "skip":
			result, err := b.handleSubsSkip(admin, topicName, input)
			return result, nil, err
		case "expire":
			result, err := b.handleSubsExpire(admin, topicName, input)
			return result, nil, err
		case "reset-cursor":
			result, err := b.handleSubsResetCursor(admin, topicName, input)
			return result, nil, err
		default:
			return nil, nil, fmt.Errorf("unknown operation: %s", operation)
		}
	}
}

// Unified error handling and utility functions

// handleError provides unified error handling
func (b *PulsarAdminSubscriptionToolBuilder) handleError(operation string, err error) error {
	return fmt.Errorf("failed to %s: %v", operation, err)
}

// marshalResponse provides unified JSON serialization for responses
func (b *PulsarAdminSubscriptionToolBuilder) marshalResponse(data interface{}) (*sdk.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, b.handleError("marshal response", err)
	}
	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: string(jsonBytes)}},
	}, nil
}

// Operation handler functions - migrated from the original implementation

// handleSubsList handles listing all subscriptions for a topic
func (b *PulsarAdminSubscriptionToolBuilder) handleSubsList(admin cmdutils.Client, topicName *utils.TopicName) (*sdk.CallToolResult, error) {
	// List subscriptions
	subscriptions, err := admin.Subscriptions().List(*topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions for topic '%s': %v", topicName.String(), err)
	}

	return b.marshalResponse(subscriptions)
}

// handleSubsCreate handles creating a new subscription
func (b *PulsarAdminSubscriptionToolBuilder) handleSubsCreate(admin cmdutils.Client, topicName *utils.TopicName, input pulsarAdminSubscriptionInput) (*sdk.CallToolResult, error) {
	// Get required parameter
	subscription, err := requireString(input.Subscription, "subscription")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'subscription' for subscription.create: %v", err)
	}

	// Get optional messageID parameter (default is "latest")
	messageID := stringValue(input.MessageID)
	if messageID == "" {
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
			return nil, fmt.Errorf("invalid messageId format: %s. use 'latest', 'earliest', or 'ledgerId:entryId' format", messageID)
		}
		msgID, err := utils.ParseMessageID(messageID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse messageId '%s': %v", messageID, err)
		}
		messageIDObj = *msgID
	}

	// Create subscription
	err = admin.Subscriptions().Create(*topicName, subscription, messageIDObj)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription '%s' on topic '%s': %v", subscription, topicName.String(), err)
	}

	return textResult(fmt.Sprintf("Created subscription '%s' on topic '%s' from position '%s' successfully", subscription, topicName.String(), messageID)), nil
}

// handleSubsDelete handles deleting a subscription
func (b *PulsarAdminSubscriptionToolBuilder) handleSubsDelete(admin cmdutils.Client, topicName *utils.TopicName, input pulsarAdminSubscriptionInput) (*sdk.CallToolResult, error) {
	// Get required parameter
	subscription, err := requireString(input.Subscription, "subscription")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'subscription' for subscription.delete: %v", err)
	}

	// Get optional force parameter (default is false)
	force := input.Force != nil && *input.Force

	// Delete subscription
	if force {
		err = admin.Subscriptions().ForceDelete(*topicName, subscription)
		if err != nil {
			return nil, fmt.Errorf("failed to forcefully delete subscription '%s' from topic '%s': %v", subscription, topicName.String(), err)
		}
	} else {
		err = admin.Subscriptions().Delete(*topicName, subscription)
		if err != nil {
			return nil, fmt.Errorf("failed to delete subscription '%s' from topic '%s': %v", subscription, topicName.String(), err)
		}
	}

	forceStr := ""
	if force {
		forceStr = " forcefully"
	}
	return textResult(fmt.Sprintf("Deleted subscription '%s' from topic '%s'%s successfully", subscription, topicName.String(), forceStr)), nil
}

// handleSubsSkip handles skipping messages for a subscription
func (b *PulsarAdminSubscriptionToolBuilder) handleSubsSkip(admin cmdutils.Client, topicName *utils.TopicName, input pulsarAdminSubscriptionInput) (*sdk.CallToolResult, error) {
	// Get required parameters
	subscription, err := requireString(input.Subscription, "subscription")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'subscription' for subscription.skip: %v", err)
	}

	if input.Count == nil {
		return nil, fmt.Errorf("missing required parameter 'count' for subscription.skip")
	}

	count := *input.Count

	// Skip messages
	err = admin.Subscriptions().SkipMessages(*topicName, subscription, int64(count))
	if err != nil {
		return nil, fmt.Errorf("failed to skip messages for subscription '%s' on topic '%s': %v", subscription, topicName.String(), err)
	}

	return textResult(fmt.Sprintf("Skipped %d messages for subscription '%s' on topic '%s' successfully", int(count), subscription, topicName.String())), nil
}

// handleSubsExpire handles expiring messages for a subscription
func (b *PulsarAdminSubscriptionToolBuilder) handleSubsExpire(admin cmdutils.Client, topicName *utils.TopicName, input pulsarAdminSubscriptionInput) (*sdk.CallToolResult, error) {
	// Get required parameters
	subscription, err := requireString(input.Subscription, "subscription")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'subscription' for subscription.expire: %v", err)
	}

	if input.ExpireTimeInSeconds == nil {
		return nil, fmt.Errorf("missing required parameter 'expireTimeInSeconds' for subscription.expire")
	}

	expireTime := *input.ExpireTimeInSeconds

	// Expire messages
	err = admin.Subscriptions().ExpireMessages(*topicName, subscription, int64(expireTime))
	if err != nil {
		return nil, fmt.Errorf("failed to expire messages for subscription '%s' on topic '%s': %v", subscription, topicName.String(), err)
	}

	return textResult(
		fmt.Sprintf("Expired messages older than %d seconds for subscription '%s' on topic '%s' successfully", int(expireTime), subscription, topicName.String()),
	), nil
}

// handleSubsResetCursor handles resetting a subscription cursor
func (b *PulsarAdminSubscriptionToolBuilder) handleSubsResetCursor(admin cmdutils.Client, topicName *utils.TopicName, input pulsarAdminSubscriptionInput) (*sdk.CallToolResult, error) {
	// Get required parameters
	subscription, err := requireString(input.Subscription, "subscription")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'subscription' for subscription.reset-cursor: %v", err)
	}

	messageID, err := requireString(input.MessageID, "messageId")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'messageId' for subscription.reset-cursor: %v", err)
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
			return nil, fmt.Errorf("invalid messageId format: %s. use 'latest', 'earliest', or 'ledgerId:entryId' format", messageID)
		}
		msgID, err := utils.ParseMessageID(messageID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse messageId '%s': %v", messageID, err)
		}
		messageIDObj = *msgID
	}

	// Reset cursor
	err = admin.Subscriptions().ResetCursorToMessageID(*topicName, subscription, messageIDObj)
	if err != nil {
		return nil, fmt.Errorf("failed to reset cursor for subscription '%s' on topic '%s': %v", subscription, topicName.String(), err)
	}

	return textResult(fmt.Sprintf("Reset cursor for subscription '%s' on topic '%s' to position '%s' successfully", subscription, topicName.String(), messageID)), nil
}

func buildPulsarAdminSubscriptionInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[pulsarAdminSubscriptionInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}

	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}

	setSchemaDescription(schema, "resource", pulsarAdminSubscriptionResourceDesc)
	setSchemaDescription(schema, "operation", pulsarAdminSubscriptionOperationDesc)
	setSchemaDescription(schema, "topic", pulsarAdminSubscriptionTopicDesc)
	setSchemaDescription(schema, "subscription", pulsarAdminSubscriptionNameDesc)
	setSchemaDescription(schema, "messageId", pulsarAdminSubscriptionMessageIDDesc)
	setSchemaDescription(schema, "count", pulsarAdminSubscriptionCountDesc)
	setSchemaDescription(schema, "expireTimeInSeconds", pulsarAdminSubscriptionExpireDesc)
	setSchemaDescription(schema, "force", pulsarAdminSubscriptionForceDesc)

	normalizeAdditionalProperties(schema)
	return schema, nil
}
