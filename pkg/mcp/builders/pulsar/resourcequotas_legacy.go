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

// PulsarAdminResourceQuotasLegacyToolBuilder implements the legacy ToolBuilder interface for Pulsar admin resource quotas.
// /nolint:revive
type PulsarAdminResourceQuotasLegacyToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminResourceQuotasLegacyToolBuilder creates a new Pulsar admin resource quotas legacy tool builder instance.
func NewPulsarAdminResourceQuotasLegacyToolBuilder() *PulsarAdminResourceQuotasLegacyToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_resourcequotas",
		Version:     "1.0.0",
		Description: "Pulsar admin resource quotas management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "admin", "resourcequotas"},
	}

	features := []string{
		"pulsar-admin-resourcequotas",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminResourceQuotasLegacyToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin resource quotas legacy tool list.
func (b *PulsarAdminResourceQuotasLegacyToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	if err := b.Validate(config); err != nil {
		return nil, err
	}

	tool, err := b.buildResourceQuotasTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildResourceQuotasHandler(config.ReadOnly)

	return []server.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

func (b *PulsarAdminResourceQuotasLegacyToolBuilder) buildResourceQuotasTool() (mcp.Tool, error) {
	inputSchema, err := buildPulsarAdminResourceQuotasInputSchema()
	if err != nil {
		return mcp.Tool{}, err
	}

	schemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		return mcp.Tool{}, fmt.Errorf("marshal input schema: %w", err)
	}

	toolDesc := "Manage Apache Pulsar resource quotas for brokers, namespaces and bundles. " +
		"Resource quotas define limits for resource usage such as message rates, bandwidth, and memory. " +
		"These quotas help prevent resource abuse and ensure fair resource allocation across the Pulsar cluster. " +
		"Operations include getting, setting, and resetting quotas. " +
		"Requires super-user permissions for all operations."

	return mcp.Tool{
		Name:           "pulsar_admin_resourcequota",
		Description:    toolDesc,
		RawInputSchema: schemaJSON,
	}, nil
}

func (b *PulsarAdminResourceQuotasLegacyToolBuilder) buildResourceQuotasHandler(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sdkBuilder := NewPulsarAdminResourceQuotasToolBuilder()
	sdkHandler := sdkBuilder.buildResourceQuotasHandler(readOnly)

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var input pulsarAdminResourceQuotasInput
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
