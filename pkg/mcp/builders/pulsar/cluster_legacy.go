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

// PulsarAdminClusterLegacyToolBuilder implements the legacy ToolBuilder interface for Pulsar admin clusters.
// /nolint:revive
type PulsarAdminClusterLegacyToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminClusterLegacyToolBuilder creates a new Pulsar admin cluster legacy tool builder instance.
func NewPulsarAdminClusterLegacyToolBuilder() *PulsarAdminClusterLegacyToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_cluster",
		Version:     "1.0.0",
		Description: "Pulsar Admin cluster management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "cluster", "admin"},
	}

	features := []string{
		"pulsar-admin-clusters",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminClusterLegacyToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin cluster legacy tool list.
func (b *PulsarAdminClusterLegacyToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	if err := b.Validate(config); err != nil {
		return nil, err
	}

	tool, err := b.buildClusterTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildClusterHandler(config.ReadOnly)

	return []server.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

func (b *PulsarAdminClusterLegacyToolBuilder) buildClusterTool() (mcp.Tool, error) {
	inputSchema, err := buildPulsarAdminClusterInputSchema()
	if err != nil {
		return mcp.Tool{}, err
	}

	schemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		return mcp.Tool{}, fmt.Errorf("marshal input schema: %w", err)
	}

	toolDesc := "Unified tool for managing Apache Pulsar clusters.\n" +
		"This tool provides access to various cluster resources and operations, including:\n" +
		"1. Manage clusters (resource=cluster): List, get, create, update, delete clusters\n" +
		"2. Manage peer clusters (resource=peer_clusters): Get, update peer clusters\n" +
		"3. Manage failure domains (resource=failure_domain): List, get, create, update, delete failure domains\n\n" +
		"Different functions are accessed by combining resource and operation parameters, with other parameters used selectively based on operation type.\n\n" +
		"Examples:\n" +
		"- {\"resource\": \"cluster\", \"operation\": \"list\"} lists all clusters\n" +
		"- {\"resource\": \"cluster\", \"operation\": \"get\", \"cluster_name\": \"my-cluster\"} gets cluster configuration\n" +
		"- {\"resource\": \"failure_domain\", \"operation\": \"list\", \"cluster_name\": \"my-cluster\"} lists failure domains\n" +
		"This tool requires Pulsar super-user permissions."

	return mcp.Tool{
		Name:           "pulsar_admin_cluster",
		Description:    toolDesc,
		RawInputSchema: schemaJSON,
	}, nil
}

func (b *PulsarAdminClusterLegacyToolBuilder) buildClusterHandler(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sdkBuilder := NewPulsarAdminClusterToolBuilder()
	sdkHandler := sdkBuilder.buildClusterHandler(readOnly)

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var input pulsarAdminClusterInput
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
