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

// PulsarAdminNamespacePolicyLegacyToolBuilder implements the legacy ToolBuilder interface for Pulsar admin namespace policies.
// /nolint:revive
type PulsarAdminNamespacePolicyLegacyToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminNamespacePolicyLegacyToolBuilder creates a new Pulsar admin namespace policy legacy tool builder instance.
func NewPulsarAdminNamespacePolicyLegacyToolBuilder() *PulsarAdminNamespacePolicyLegacyToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_namespace_policy",
		Version:     "1.0.0",
		Description: "Pulsar admin namespace policy management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "admin", "namespace_policy"},
	}

	features := []string{
		"pulsar-admin-namespace-policy",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminNamespacePolicyLegacyToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin namespace policy legacy tool list.
func (b *PulsarAdminNamespacePolicyLegacyToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	if err := b.Validate(config); err != nil {
		return nil, err
	}

	tools := []server.ServerTool{}

	getTool, err := b.buildNamespaceGetPoliciesTool()
	if err != nil {
		return nil, err
	}
	getHandler := b.buildNamespaceGetPoliciesHandler()
	tools = append(tools, server.ServerTool{
		Tool:    getTool,
		Handler: getHandler,
	})

	if !config.ReadOnly {
		setTool, err := b.buildNamespaceSetPolicyTool()
		if err != nil {
			return nil, err
		}
		setHandler := b.buildNamespaceSetPolicyHandler()
		tools = append(tools, server.ServerTool{
			Tool:    setTool,
			Handler: setHandler,
		})

		removeTool, err := b.buildNamespaceRemovePolicyTool()
		if err != nil {
			return nil, err
		}
		removeHandler := b.buildNamespaceRemovePolicyHandler()
		tools = append(tools, server.ServerTool{
			Tool:    removeTool,
			Handler: removeHandler,
		})
	}

	return tools, nil
}

func (b *PulsarAdminNamespacePolicyLegacyToolBuilder) buildNamespaceGetPoliciesTool() (mcp.Tool, error) {
	inputSchema, err := buildPulsarAdminNamespacePolicyGetInputSchema()
	if err != nil {
		return mcp.Tool{}, err
	}

	schemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		return mcp.Tool{}, fmt.Errorf("marshal input schema: %w", err)
	}

	return mcp.Tool{
		Name:           "pulsar_admin_namespace_policy_get",
		Description:    pulsarAdminNamespacePolicyGetToolDesc,
		RawInputSchema: schemaJSON,
	}, nil
}

func (b *PulsarAdminNamespacePolicyLegacyToolBuilder) buildNamespaceSetPolicyTool() (mcp.Tool, error) {
	inputSchema, err := buildPulsarAdminNamespacePolicySetInputSchema()
	if err != nil {
		return mcp.Tool{}, err
	}

	schemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		return mcp.Tool{}, fmt.Errorf("marshal input schema: %w", err)
	}

	return mcp.Tool{
		Name:           "pulsar_admin_namespace_policy_set",
		Description:    pulsarAdminNamespacePolicySetToolDesc,
		RawInputSchema: schemaJSON,
	}, nil
}

func (b *PulsarAdminNamespacePolicyLegacyToolBuilder) buildNamespaceRemovePolicyTool() (mcp.Tool, error) {
	inputSchema, err := buildPulsarAdminNamespacePolicyRemoveInputSchema()
	if err != nil {
		return mcp.Tool{}, err
	}

	schemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		return mcp.Tool{}, fmt.Errorf("marshal input schema: %w", err)
	}

	return mcp.Tool{
		Name:           "pulsar_admin_namespace_policy_remove",
		Description:    pulsarAdminNamespacePolicyRemoveToolDesc,
		RawInputSchema: schemaJSON,
	}, nil
}

func (b *PulsarAdminNamespacePolicyLegacyToolBuilder) buildNamespaceGetPoliciesHandler() func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sdkBuilder := NewPulsarAdminNamespacePolicyToolBuilder()
	sdkHandler := sdkBuilder.buildNamespaceGetPoliciesHandler()

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var input pulsarAdminNamespacePolicyGetInput
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

func (b *PulsarAdminNamespacePolicyLegacyToolBuilder) buildNamespaceSetPolicyHandler() func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sdkBuilder := NewPulsarAdminNamespacePolicyToolBuilder()
	sdkHandler := sdkBuilder.buildNamespaceSetPolicyHandler()

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var input pulsarAdminNamespacePolicySetInput
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

func (b *PulsarAdminNamespacePolicyLegacyToolBuilder) buildNamespaceRemovePolicyHandler() func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sdkBuilder := NewPulsarAdminNamespacePolicyToolBuilder()
	sdkHandler := sdkBuilder.buildNamespaceRemovePolicyHandler()

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var input pulsarAdminNamespacePolicyRemoveInput
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
