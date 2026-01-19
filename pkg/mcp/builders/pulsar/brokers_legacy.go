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

// PulsarAdminBrokersLegacyToolBuilder implements the legacy ToolBuilder interface for Pulsar admin brokers.
// /nolint:revive
type PulsarAdminBrokersLegacyToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminBrokersLegacyToolBuilder creates a new Pulsar admin brokers legacy tool builder instance.
func NewPulsarAdminBrokersLegacyToolBuilder() *PulsarAdminBrokersLegacyToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_brokers",
		Version:     "1.0.0",
		Description: "Pulsar admin brokers management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "admin", "brokers"},
	}

	features := []string{
		"pulsar-admin-brokers",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminBrokersLegacyToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin brokers legacy tool list.
func (b *PulsarAdminBrokersLegacyToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	if err := b.Validate(config); err != nil {
		return nil, err
	}

	tool, err := b.buildPulsarAdminBrokersTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildPulsarAdminBrokersHandler(config.ReadOnly)

	return []server.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

func (b *PulsarAdminBrokersLegacyToolBuilder) buildPulsarAdminBrokersTool() (mcp.Tool, error) {
	inputSchema, err := buildPulsarAdminBrokersInputSchema()
	if err != nil {
		return mcp.Tool{}, err
	}

	schemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		return mcp.Tool{}, fmt.Errorf("marshal input schema: %w", err)
	}

	toolDesc := "Unified tool for managing Apache Pulsar broker resources. This tool integrates multiple broker management functions, including:\n" +
		"1. List active brokers in a cluster (resource=brokers, operation=list)\n" +
		"2. Check broker health status (resource=health, operation=get)\n" +
		"3. Manage broker configurations (resource=config, operation=get/update/delete)\n" +
		"4. View namespaces owned by a broker (resource=namespaces, operation=get)\n\n" +
		"Different functions are accessed by combining resource and operation parameters, with other parameters used selectively based on operation type.\n" +
		"Example: {\"resource\": \"config\", \"operation\": \"get\", \"configType\": \"dynamic\"} retrieves all dynamic configuration names.\n" +
		"This tool requires Pulsar super-user permissions."

	return mcp.Tool{
		Name:           "pulsar_admin_brokers",
		Description:    toolDesc,
		RawInputSchema: schemaJSON,
	}, nil
}

func (b *PulsarAdminBrokersLegacyToolBuilder) buildPulsarAdminBrokersHandler(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sdkBuilder := NewPulsarAdminBrokersToolBuilder()
	sdkHandler := sdkBuilder.buildPulsarAdminBrokersHandler(readOnly)

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var input pulsarAdminBrokersInput
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
