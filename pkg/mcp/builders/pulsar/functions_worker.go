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
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

// PulsarAdminFunctionsWorkerToolBuilder implements the ToolBuilder interface for Pulsar Admin Functions Worker tools
// It provides functionality to build Pulsar functions worker monitoring and management tools
// /nolint:revive
type PulsarAdminFunctionsWorkerToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminFunctionsWorkerToolBuilder creates a new Pulsar Admin Functions Worker tool builder instance
func NewPulsarAdminFunctionsWorkerToolBuilder() *PulsarAdminFunctionsWorkerToolBuilder {
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

	return &PulsarAdminFunctionsWorkerToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar Admin Functions Worker tool list
// This is the core method implementing the ToolBuilder interface
func (b *PulsarAdminFunctionsWorkerToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool := b.buildFunctionsWorkerTool()
	handler := b.buildFunctionsWorkerHandler(config.ReadOnly)

	return []server.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildFunctionsWorkerTool builds the Pulsar Admin Functions Worker MCP tool definition
// Migrated from the original tool definition logic
func (b *PulsarAdminFunctionsWorkerToolBuilder) buildFunctionsWorkerTool() mcp.Tool {
	toolDesc := "Unified tool for managing Apache Pulsar Functions Worker resources. " +
		"Pulsar Functions is a serverless compute framework that allows you to process messages in a streaming fashion. " +
		"The Functions Worker is the runtime environment that executes and manages Pulsar Functions. " +
		"This tool provides comprehensive access to functions worker resources including function statistics, " +
		"monitoring metrics, cluster information, leader election status, and function assignments across the cluster. " +
		"Functions workers can be deployed in multiple modes (standalone, cluster) and this tool helps monitor " +
		"and manage the worker cluster state, performance metrics, and function distribution. " +
		"Most operations require Pulsar super-user permissions for security reasons."

	resourceDesc := "Type of functions worker resource to access. Available resources:\\n" +
		"- function_stats: Statistics for all functions running on the functions worker, including processing rates, error counts, and resource usage\\n" +
		"- monitoring_metrics: Comprehensive metrics for monitoring function workers, including JVM metrics, resource utilization, and performance indicators\\n" +
		"- cluster: Information about all workers in the functions worker cluster, including their status, capabilities, and workload distribution\\n" +
		"- cluster_leader: Information about the leader of the functions worker cluster, essential for understanding cluster coordination\\n" +
		"- function_assignments: Current assignments of functions across the functions worker cluster, showing which functions are running on which workers"

	return mcp.NewTool("pulsar_admin_functions_worker",
		mcp.WithDescription(toolDesc),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description(resourceDesc),
		),
	)
}

// buildFunctionsWorkerHandler builds the Pulsar Admin Functions Worker handler function
// Migrated from the original handler logic
func (b *PulsarAdminFunctionsWorkerToolBuilder) buildFunctionsWorkerHandler(_ bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return mcp.NewToolResultError("Pulsar session not found in context"), nil
		}

		// Create the admin client
		admin, err := session.GetAdminClient()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
		}

		// Get required resource parameter
		resource, err := request.RequireString("resource")
		if err != nil {
			return mcp.NewToolResultError("Missing required parameter 'resource'. " +
				"Please specify one of: function_stats, monitoring_metrics, cluster, cluster_leader, function_assignments"), nil
		}

		// Process request based on resource type
		switch resource {
		case "function_stats":
			return b.handleFunctionsWorkerFunctionStats(admin)
		case "monitoring_metrics":
			return b.handleFunctionsWorkerMonitoringMetrics(admin)
		case "cluster":
			return b.handleFunctionsWorkerGetCluster(admin)
		case "cluster_leader":
			return b.handleFunctionsWorkerGetClusterLeader(admin)
		case "function_assignments":
			return b.handleFunctionsWorkerGetFunctionAssignments(admin)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Unsupported resource: %s. "+
				"Please use one of: function_stats, monitoring_metrics, cluster, cluster_leader, function_assignments", resource)), nil
		}
	}
}

// Unified error handling and utility functions

// handleError provides unified error handling
func (b *PulsarAdminFunctionsWorkerToolBuilder) handleError(operation string, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("Failed to %s: %v", operation, err))
}

// marshalResponse provides unified JSON serialization for responses
func (b *PulsarAdminFunctionsWorkerToolBuilder) marshalResponse(data interface{}) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return b.handleError("marshal response", err), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// Operation handler functions - migrated from the original implementation

// handleFunctionsWorkerFunctionStats handles retrieving function statistics
func (b *PulsarAdminFunctionsWorkerToolBuilder) handleFunctionsWorkerFunctionStats(admin cmdutils.Client) (*mcp.CallToolResult, error) {
	// Get function stats
	stats, err := admin.FunctionsWorker().GetFunctionsStats()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get functions stats: %v", err)), nil
	}

	return b.marshalResponse(stats)
}

// handleFunctionsWorkerMonitoringMetrics handles retrieving monitoring metrics
func (b *PulsarAdminFunctionsWorkerToolBuilder) handleFunctionsWorkerMonitoringMetrics(admin cmdutils.Client) (*mcp.CallToolResult, error) {
	// Get monitoring metrics
	metrics, err := admin.FunctionsWorker().GetMetrics()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get monitoring metrics: %v", err)), nil
	}

	return b.marshalResponse(metrics)
}

// handleFunctionsWorkerGetCluster handles retrieving cluster information
func (b *PulsarAdminFunctionsWorkerToolBuilder) handleFunctionsWorkerGetCluster(admin cmdutils.Client) (*mcp.CallToolResult, error) {
	// Get cluster info
	cluster, err := admin.FunctionsWorker().GetCluster()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get worker cluster: %v", err)), nil
	}

	return b.marshalResponse(cluster)
}

// handleFunctionsWorkerGetClusterLeader handles retrieving cluster leader information
func (b *PulsarAdminFunctionsWorkerToolBuilder) handleFunctionsWorkerGetClusterLeader(admin cmdutils.Client) (*mcp.CallToolResult, error) {
	// Get cluster leader
	leader, err := admin.FunctionsWorker().GetClusterLeader()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get worker cluster leader: %v", err)), nil
	}

	return b.marshalResponse(leader)
}

// handleFunctionsWorkerGetFunctionAssignments handles retrieving function assignments
func (b *PulsarAdminFunctionsWorkerToolBuilder) handleFunctionsWorkerGetFunctionAssignments(admin cmdutils.Client) (*mcp.CallToolResult, error) {
	// Get function assignments
	assignments, err := admin.FunctionsWorker().GetAssignments()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get function assignments: %v", err)), nil
	}

	return b.marshalResponse(assignments)
}
