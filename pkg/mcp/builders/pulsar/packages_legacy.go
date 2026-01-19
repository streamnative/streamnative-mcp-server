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

// PulsarAdminPackagesLegacyToolBuilder implements the legacy ToolBuilder interface for Pulsar admin packages.
// /nolint:revive
type PulsarAdminPackagesLegacyToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminPackagesLegacyToolBuilder creates a new Pulsar admin packages legacy tool builder instance.
func NewPulsarAdminPackagesLegacyToolBuilder() *PulsarAdminPackagesLegacyToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_packages",
		Version:     "1.0.0",
		Description: "Pulsar admin packages management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "admin", "packages"},
	}

	features := []string{
		"pulsar-admin-packages",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminPackagesLegacyToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin packages legacy tool list.
func (b *PulsarAdminPackagesLegacyToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	if err := b.Validate(config); err != nil {
		return nil, err
	}

	tool, err := b.buildPackagesTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildPackagesHandler(config.ReadOnly)

	return []server.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

func (b *PulsarAdminPackagesLegacyToolBuilder) buildPackagesTool() (mcp.Tool, error) {
	inputSchema, err := buildPulsarAdminPackagesInputSchema()
	if err != nil {
		return mcp.Tool{}, err
	}

	schemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		return mcp.Tool{}, fmt.Errorf("marshal input schema: %w", err)
	}

	toolDesc := "Manage packages in Apache Pulsar. Support package scheme: `function://`, `source://`, `sink://`" +
		"Allows listing, viewing, updating, downloading and uploading packages. " +
		"Some operations require super-user permissions."

	return mcp.Tool{
		Name:           "pulsar_admin_package",
		Description:    toolDesc,
		RawInputSchema: schemaJSON,
	}, nil
}

func (b *PulsarAdminPackagesLegacyToolBuilder) buildPackagesHandler(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sdkBuilder := NewPulsarAdminPackagesToolBuilder()
	sdkHandler := sdkBuilder.buildPackagesHandler(readOnly)

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var input pulsarAdminPackagesInput
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
