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

// PulsarAdminSourcesLegacyToolBuilder implements the legacy ToolBuilder interface for Pulsar admin sources.
// /nolint:revive
type PulsarAdminSourcesLegacyToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminSourcesLegacyToolBuilder creates a new Pulsar admin sources legacy tool builder instance.
func NewPulsarAdminSourcesLegacyToolBuilder() *PulsarAdminSourcesLegacyToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_sources",
		Version:     "1.0.0",
		Description: "Pulsar admin sources management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "admin", "sources"},
	}

	features := []string{
		"pulsar-admin-sources",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminSourcesLegacyToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin sources legacy tool list.
func (b *PulsarAdminSourcesLegacyToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	if err := b.Validate(config); err != nil {
		return nil, err
	}

	tool, err := b.buildSourcesTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildSourcesHandler(config.ReadOnly)

	return []server.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

func (b *PulsarAdminSourcesLegacyToolBuilder) buildSourcesTool() (mcp.Tool, error) {
	inputSchema, err := buildPulsarAdminSourcesInputSchema()
	if err != nil {
		return mcp.Tool{}, err
	}

	schemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		return mcp.Tool{}, fmt.Errorf("marshal input schema: %w", err)
	}

	toolDesc := "Manage Apache Pulsar Sources for data ingestion and integration. " +
		"Pulsar Sources are connectors that import data from external systems into Pulsar topics. " +
		"Sources connect to external systems such as databases, messaging platforms, storage services, " +
		"and real-time data streams to pull data and publish it to Pulsar topics. " +
		"Built-in source connectors are available for common systems like Kafka, JDBC, AWS services, and more. " +
		"Sources follow the tenant/namespace/name hierarchy for organization and access control, " +
		"can scale through parallelism configuration, and support various processing guarantees. " +
		"This tool provides complete lifecycle management including deployment, configuration, " +
		"monitoring, and runtime control. Sources use schema types to ensure data compatibility."

	return mcp.Tool{
		Name:           "pulsar_admin_sources",
		Description:    toolDesc,
		RawInputSchema: schemaJSON,
	}, nil
}

func (b *PulsarAdminSourcesLegacyToolBuilder) buildSourcesHandler(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sdkBuilder := NewPulsarAdminSourcesToolBuilder()
	sdkHandler := sdkBuilder.buildSourcesHandler(readOnly)

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var input pulsarAdminSourcesInput
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
