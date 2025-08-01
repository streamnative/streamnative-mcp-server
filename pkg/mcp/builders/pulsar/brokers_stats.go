// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package pulsar

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// PulsarAdminBrokerStatsToolBuilder implements the ToolBuilder interface for Pulsar Broker Statistics
type PulsarAdminBrokerStatsToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminBrokerStatsToolBuilder creates a new Pulsar Admin Broker Stats tool builder instance
func NewPulsarAdminBrokerStatsToolBuilder() *PulsarAdminBrokerStatsToolBuilder {
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

	return &PulsarAdminBrokerStatsToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar Admin Broker Stats tool list
func (b *PulsarAdminBrokerStatsToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool := b.buildBrokerStatsTool()
	handler := b.buildBrokerStatsHandler(config.ReadOnly)

	return []server.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildBrokerStatsTool builds the Pulsar Admin Broker Stats MCP tool definition
func (b *PulsarAdminBrokerStatsToolBuilder) buildBrokerStatsTool() mcp.Tool {
	resourceDesc := "Type of broker stats resource to access, available options:\n" +
		"- monitoring_metrics: Metrics for the broker's monitoring system\n" +
		"- mbeans: JVM MBeans statistics\n" +
		"- topics: Statistics about all topics managed by the broker\n" +
		"- allocator_stats: Memory allocator statistics (requires allocator_name parameter)\n" +
		"- load_report: Broker load information"

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

	return mcp.NewTool("pulsar_admin_broker_stats",
		mcp.WithDescription(toolDesc),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description(resourceDesc),
		),
		mcp.WithString("allocator_name",
			mcp.Description("The name of the allocator to get statistics for. Required only when resource=allocator_stats"),
		),
	)
}

// buildBrokerStatsHandler builds the Pulsar Admin Broker Stats handler function
func (b *PulsarAdminBrokerStatsToolBuilder) buildBrokerStatsHandler(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get Pulsar admin client
		client, err := b.getPulsarAdminClient(ctx)
		if err != nil {
			return b.handleError("get admin client", err), nil
		}

		// Get required resource parameter
		resource, err := request.RequireString("resource")
		if err != nil {
			return mcp.NewToolResultError("Missing required parameter 'resource'. " +
				"Please specify one of: monitoring_metrics, mbeans, topics, allocator_stats, load_report."), nil
		}

		// Process request based on resource type
		switch resource {
		case "monitoring_metrics":
			return b.handleMonitoringMetrics(client)
		case "mbeans":
			return b.handleMBeans(client)
		case "topics":
			return b.handleTopics(client)
		case "allocator_stats":
			allocatorName, err := request.RequireString("allocator_name")
			if err != nil {
				return mcp.NewToolResultError("Missing required parameter 'allocator_name' for allocator_stats resource. " +
					"Please provide the name of the allocator to get statistics for."), nil
			}
			return b.handleAllocatorStats(client, allocatorName)
		case "load_report":
			return b.handleLoadReport(client)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Unsupported resource: %s. "+
				"Please use one of: monitoring_metrics, mbeans, topics, allocator_stats, load_report.", resource)), nil
		}
	}
}

// Utility functions

// handleError provides unified error handling
func (b *PulsarAdminBrokerStatsToolBuilder) handleError(operation string, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("Failed to %s: %v", operation, err))
}

// marshalResponse provides unified JSON serialization for responses
func (b *PulsarAdminBrokerStatsToolBuilder) marshalResponse(data interface{}) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return b.handleError("marshal response", err), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// getPulsarAdminClient retrieves the Pulsar admin client from context
func (b *PulsarAdminBrokerStatsToolBuilder) getPulsarAdminClient(ctx context.Context) (cmdutils.Client, error) {
	// Use the same context key as in the original implementation (pkg/mcp/ctx.go)
	// This maintains consistency with the original approach
	type contextKey string
	const pulsarSessionContextKey contextKey = "pulsar_session"
	
	session, ok := ctx.Value(pulsarSessionContextKey).(*pulsar.Session)
	if !ok || session == nil {
		return nil, fmt.Errorf("Pulsar session not found in context")
	}
	return session.GetAdminClient()
}

// Specific operation handler functions

// handleMonitoringMetrics handles retrieving monitoring metrics
func (b *PulsarAdminBrokerStatsToolBuilder) handleMonitoringMetrics(client cmdutils.Client) (*mcp.CallToolResult, error) {
	stats, err := client.BrokerStats().GetMetrics()
	if err != nil {
		return b.handleError("get monitoring metrics", err), nil
	}
	return b.marshalResponse(stats)
}

// handleMBeans handles retrieving MBeans statistics
func (b *PulsarAdminBrokerStatsToolBuilder) handleMBeans(client cmdutils.Client) (*mcp.CallToolResult, error) {
	stats, err := client.BrokerStats().GetMBeans()
	if err != nil {
		return b.handleError("get MBeans", err), nil
	}
	return b.marshalResponse(stats)
}

// handleTopics handles retrieving topics statistics
func (b *PulsarAdminBrokerStatsToolBuilder) handleTopics(client cmdutils.Client) (*mcp.CallToolResult, error) {
	stats, err := client.BrokerStats().GetTopics()
	if err != nil {
		return b.handleError("get topics stats", err), nil
	}
	return b.marshalResponse(stats)
}

// handleAllocatorStats handles retrieving allocator statistics
func (b *PulsarAdminBrokerStatsToolBuilder) handleAllocatorStats(client cmdutils.Client, allocatorName string) (*mcp.CallToolResult, error) {
	stats, err := client.BrokerStats().GetAllocatorStats(allocatorName)
	if err != nil {
		return b.handleError("get allocator stats", err), nil
	}
	return b.marshalResponse(stats)
}

// handleLoadReport handles retrieving load report
func (b *PulsarAdminBrokerStatsToolBuilder) handleLoadReport(client cmdutils.Client) (*mcp.CallToolResult, error) {
	stats, err := client.BrokerStats().GetLoadReport()
	if err != nil {
		return b.handleError("get load report", err), nil
	}
	return b.marshalResponse(stats)
}