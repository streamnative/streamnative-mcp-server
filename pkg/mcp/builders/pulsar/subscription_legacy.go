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

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
)

// PulsarAdminSubscriptionLegacyToolBuilder implements the legacy ToolBuilder interface for Pulsar admin subscriptions.
// /nolint:revive
type PulsarAdminSubscriptionLegacyToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminSubscriptionLegacyToolBuilder creates a new Pulsar admin subscription legacy tool builder instance.
func NewPulsarAdminSubscriptionLegacyToolBuilder() *PulsarAdminSubscriptionLegacyToolBuilder {
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

	return &PulsarAdminSubscriptionLegacyToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin subscription legacy tool list.
func (b *PulsarAdminSubscriptionLegacyToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	if err := b.Validate(config); err != nil {
		return nil, err
	}

	tool, err := b.buildSubscriptionTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildSubscriptionHandler(config.ReadOnly)

	return []server.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

func (b *PulsarAdminSubscriptionLegacyToolBuilder) buildSubscriptionTool() (mcp.Tool, error) {
	inputSchema, err := buildPulsarAdminSubscriptionInputSchema()
	if err != nil {
		return mcp.Tool{}, err
	}

	schemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		return mcp.Tool{}, fmt.Errorf("marshal input schema: %w", err)
	}

	toolDesc := "Manage Apache Pulsar subscriptions on topics. " +
		"Subscriptions are named entities representing consumer groups that maintain their position in a topic. " +
		"Pulsar supports multiple subscription modes (Exclusive, Shared, Failover, Key_Shared) to accommodate different messaging patterns. " +
		"Each subscription tracks message acknowledgments independently, allowing multiple consumers to process messages at their own pace. " +
		"Subscriptions persist even when all consumers disconnect, maintaining state and preventing message loss. " +
		"Operations include listing, creating, deleting, and manipulating message cursors within subscriptions. " +
		"Most operations require namespace admin permissions plus produce/consume permissions on the topic."

	return mcp.Tool{
		Name:           "pulsar_admin_subscription",
		Description:    toolDesc,
		RawInputSchema: schemaJSON,
	}, nil
}

func (b *PulsarAdminSubscriptionLegacyToolBuilder) buildSubscriptionHandler(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sdkBuilder := NewPulsarAdminSubscriptionToolBuilder()
	sdkHandler := sdkBuilder.buildSubscriptionHandler(readOnly)

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var input pulsarAdminSubscriptionInput
		if err := request.BindArguments(&input); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to parse arguments: %v", err)), nil
		}

		result, _, err := sdkHandler(ctx, nil, input)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return legacyToolResultFromSDK(result), nil
	}
}
