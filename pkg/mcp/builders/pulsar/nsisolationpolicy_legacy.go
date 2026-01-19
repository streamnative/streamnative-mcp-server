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

// PulsarAdminNsIsolationPolicyLegacyToolBuilder implements the legacy ToolBuilder interface for Pulsar admin namespace isolation policy tools.
// /nolint:revive
type PulsarAdminNsIsolationPolicyLegacyToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminNsIsolationPolicyLegacyToolBuilder creates a new Pulsar admin namespace isolation policy legacy tool builder instance.
func NewPulsarAdminNsIsolationPolicyLegacyToolBuilder() *PulsarAdminNsIsolationPolicyLegacyToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_nsisolationpolicy",
		Version:     "1.0.0",
		Description: "Pulsar admin namespace isolation policy management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "admin", "nsisolationpolicy"},
	}

	features := []string{
		"pulsar-admin-nsisolationpolicy",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminNsIsolationPolicyLegacyToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin namespace isolation policy legacy tool list.
func (b *PulsarAdminNsIsolationPolicyLegacyToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	if err := b.Validate(config); err != nil {
		return nil, err
	}

	tool, err := b.buildNsIsolationPolicyTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildNsIsolationPolicyHandler(config.ReadOnly)

	return []server.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

func (b *PulsarAdminNsIsolationPolicyLegacyToolBuilder) buildNsIsolationPolicyTool() (mcp.Tool, error) {
	inputSchema, err := buildPulsarAdminNsIsolationPolicyInputSchema()
	if err != nil {
		return mcp.Tool{}, err
	}

	schemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		return mcp.Tool{}, fmt.Errorf("marshal input schema: %w", err)
	}

	toolDesc := "Manage namespace isolation policies in a Pulsar cluster. " +
		"Allows viewing, creating, updating, and deleting namespace isolation policies. " +
		"Some operations require super-user permissions."

	return mcp.Tool{
		Name:           "pulsar_admin_nsisolationpolicy",
		Description:    toolDesc,
		RawInputSchema: schemaJSON,
	}, nil
}

func (b *PulsarAdminNsIsolationPolicyLegacyToolBuilder) buildNsIsolationPolicyHandler(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sdkBuilder := NewPulsarAdminNsIsolationPolicyToolBuilder()
	sdkHandler := sdkBuilder.buildNsIsolationPolicyHandler(readOnly)

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var input pulsarAdminNsIsolationPolicyInput
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
