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

// PulsarAdminSinksLegacyToolBuilder implements the legacy ToolBuilder interface for Pulsar admin sinks.
// /nolint:revive
type PulsarAdminSinksLegacyToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminSinksLegacyToolBuilder creates a new Pulsar admin sinks legacy tool builder instance.
func NewPulsarAdminSinksLegacyToolBuilder() *PulsarAdminSinksLegacyToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_sinks",
		Version:     "1.0.0",
		Description: "Pulsar admin sinks management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "admin", "sinks"},
	}

	features := []string{
		"pulsar-admin-sinks",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminSinksLegacyToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin sinks legacy tool list.
func (b *PulsarAdminSinksLegacyToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	if err := b.Validate(config); err != nil {
		return nil, err
	}

	tool, err := b.buildSinksTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildSinksHandler(config.ReadOnly)

	return []server.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

func (b *PulsarAdminSinksLegacyToolBuilder) buildSinksTool() (mcp.Tool, error) {
	inputSchema, err := buildPulsarAdminSinksInputSchema()
	if err != nil {
		return mcp.Tool{}, err
	}

	schemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		return mcp.Tool{}, fmt.Errorf("marshal input schema: %w", err)
	}

	toolDesc := "Manage Apache Pulsar Sinks for data movement and integration. " +
		"Pulsar Sinks are connectors that export data from Pulsar topics to external systems such as databases, " +
		"storage services, messaging systems, and third-party applications. " +
		"Sinks consume messages from one or more Pulsar topics, transform the data if needed, " +
		"and write it to external systems in a format compatible with the target destination. " +
		"Built-in sink connectors are available for common systems like Kafka, JDBC, Elasticsearch, and cloud storage. " +
		"Sinks follow the tenant/namespace/name hierarchy for organization and access control, " +
		"can scale through parallelism configuration, and support configurable subscription types. " +
		"This tool provides complete lifecycle management including deployment, configuration, " +
		"monitoring, and runtime control. Sinks require proper permissions to access their input topics."

	return mcp.Tool{
		Name:           "pulsar_admin_sinks",
		Description:    toolDesc,
		RawInputSchema: schemaJSON,
	}, nil
}

func (b *PulsarAdminSinksLegacyToolBuilder) buildSinksHandler(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sdkBuilder := NewPulsarAdminSinksToolBuilder()
	sdkHandler := sdkBuilder.buildSinksHandler(readOnly)

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var input pulsarAdminSinksInput
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
