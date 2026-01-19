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

// PulsarAdminBrokerStatsLegacyToolBuilder implements the legacy ToolBuilder interface for Pulsar broker stats.
// /nolint:revive
type PulsarAdminBrokerStatsLegacyToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminBrokerStatsLegacyToolBuilder creates a new Pulsar admin broker stats legacy tool builder instance.
func NewPulsarAdminBrokerStatsLegacyToolBuilder() *PulsarAdminBrokerStatsLegacyToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_broker_stats",
		Version:     "1.0.0",
		Description: "Pulsar Broker Statistics administration tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "broker", "stats", "admin", "monitoring"},
	}

	features := []string{
		"pulsar-admin-brokers-status",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminBrokerStatsLegacyToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin broker stats legacy tool list.
func (b *PulsarAdminBrokerStatsLegacyToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	if err := b.Validate(config); err != nil {
		return nil, err
	}

	tool, err := b.buildBrokerStatsTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildBrokerStatsHandler(config.ReadOnly)

	return []server.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

func (b *PulsarAdminBrokerStatsLegacyToolBuilder) buildBrokerStatsTool() (mcp.Tool, error) {
	inputSchema, err := buildPulsarAdminBrokerStatsInputSchema()
	if err != nil {
		return mcp.Tool{}, err
	}

	schemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		return mcp.Tool{}, fmt.Errorf("marshal input schema: %w", err)
	}

	toolDesc := "Unified tool for retrieving Apache Pulsar broker statistics.\n" +
		"This tool provides access to various broker stats resources, including:\n" +
		"1. Monitoring metrics (resource=monitoring_metrics): Metrics for the broker's monitoring system\n" +
		"2. MBean stats (resource=mbeans): JVM MBeans statistics\n" +
		"3. Topics stats (resource=topics): Statistics about all topics managed by the broker\n" +
		"4. Allocator stats (resource=allocator_stats): Memory allocator statistics for specific allocator\n" +
		"5. Load report (resource=load_report): Broker load information, sometimes the load report is not available, so suggest to use other resources to get the broker metrics\n\n" +
		"Example: {\"resource\": \"monitoring_metrics\"} retrieves all monitoring metrics\n" +
		"Example: {\"resource\": \"allocator_stats\", \"allocator_name\": \"default\"} retrieves stats for the default allocator\n" +
		"This tool requires Pulsar super-user permissions."

	return mcp.Tool{
		Name:           "pulsar_admin_broker_stats",
		Description:    toolDesc,
		RawInputSchema: schemaJSON,
	}, nil
}

func (b *PulsarAdminBrokerStatsLegacyToolBuilder) buildBrokerStatsHandler(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sdkBuilder := NewPulsarAdminBrokerStatsToolBuilder()
	sdkHandler := sdkBuilder.buildBrokerStatsHandler(readOnly)

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var input pulsarAdminBrokerStatsInput
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
