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

// PulsarAdminFunctionsWorkerLegacyToolBuilder implements the legacy ToolBuilder interface for Pulsar functions worker tools.
// /nolint:revive
type PulsarAdminFunctionsWorkerLegacyToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminFunctionsWorkerLegacyToolBuilder creates a new Pulsar admin functions worker legacy tool builder instance.
func NewPulsarAdminFunctionsWorkerLegacyToolBuilder() *PulsarAdminFunctionsWorkerLegacyToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_functions_worker",
		Version:     "1.0.0",
		Description: "Pulsar Admin functions worker management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "functions", "worker", "admin", "monitoring"},
	}

	features := []string{
		"pulsar-admin-functions-worker",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminFunctionsWorkerLegacyToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin functions worker legacy tool list.
func (b *PulsarAdminFunctionsWorkerLegacyToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	if err := b.Validate(config); err != nil {
		return nil, err
	}

	tool, err := b.buildFunctionsWorkerTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildFunctionsWorkerHandler(config.ReadOnly)

	return []server.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

func (b *PulsarAdminFunctionsWorkerLegacyToolBuilder) buildFunctionsWorkerTool() (mcp.Tool, error) {
	inputSchema, err := buildPulsarAdminFunctionsWorkerInputSchema()
	if err != nil {
		return mcp.Tool{}, err
	}

	schemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		return mcp.Tool{}, fmt.Errorf("marshal input schema: %w", err)
	}

	return mcp.Tool{
		Name:           "pulsar_admin_functions_worker",
		Description:    pulsarAdminFunctionsWorkerToolDesc,
		RawInputSchema: schemaJSON,
	}, nil
}

func (b *PulsarAdminFunctionsWorkerLegacyToolBuilder) buildFunctionsWorkerHandler(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sdkBuilder := NewPulsarAdminFunctionsWorkerToolBuilder()
	sdkHandler := sdkBuilder.buildFunctionsWorkerHandler(readOnly)

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var input pulsarAdminFunctionsWorkerInput
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
