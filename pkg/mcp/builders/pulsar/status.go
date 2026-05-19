// Copyright 2026 StreamNative
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
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/toolannotations"
)

// PulsarAdminStatusToolBuilder implements the ToolBuilder interface for Pulsar status checks.
// /nolint:revive
type PulsarAdminStatusToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminStatusToolBuilder creates a new Pulsar admin status tool builder instance.
func NewPulsarAdminStatusToolBuilder() *PulsarAdminStatusToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_status",
		Version:     "1.0.0",
		Description: "Pulsar admin broker/proxy status check tool",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "admin", "status", "health"},
	}

	features := []string{
		"pulsar-admin-brokers-status",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminStatusToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin status tool list.
func (b *PulsarAdminStatusToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	if err := b.Validate(config); err != nil {
		return nil, err
	}

	return []server.ServerTool{
		{
			Tool:    b.buildStatusTool(),
			Handler: b.buildStatusHandler(),
		},
	}, nil
}

func (b *PulsarAdminStatusToolBuilder) buildStatusTool() mcp.Tool {
	return mcp.NewTool("pulsar_admin_status",
		mcp.WithDescription("Check Pulsar broker or proxy service status via the /status.html endpoint. "+
			"This is equivalent to `pulsarctl status check` and requires super-user permissions."),
		toolannotations.ReadOnly("Check Pulsar Status"),
	)
}

func (b *PulsarAdminStatusToolBuilder) buildStatusHandler() func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return mcp.NewToolResultError("Pulsar session not found in context"), nil
		}

		statusClient, err := session.GetAdminStatusClient()
		if err != nil {
			return b.handleError("get Pulsar status client", err), nil
		}

		data, err := statusClient.GetWithQueryParams("/status.html", nil, nil, false)
		if err != nil {
			return b.handleError("check Pulsar status", err), nil
		}

		return mcp.NewToolResultText(strings.TrimSpace(string(data))), nil
	}
}

func (b *PulsarAdminStatusToolBuilder) handleError(operation string, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("Failed to %s: %v", operation, err))
}
