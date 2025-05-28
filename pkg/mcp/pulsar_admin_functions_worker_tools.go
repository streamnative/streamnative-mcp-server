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

package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// PulsarAdminAddFunctionsWorkerTools adds functions worker-related tools to the MCP server
func PulsarAdminAddFunctionsWorkerTools(s *server.MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeaturePulsarAdminFunctionsWorker)) && !slices.Contains(features, string(FeatureAll)) && !slices.Contains(features, string(FeatureAllPulsar)) && !slices.Contains(features, string(FeaturePulsarAdmin)) {
		return
	}
	// Create a single unified functions worker tool
	functionsWorkerTool := mcp.NewTool("pulsar_admin_functions_worker",
		mcp.WithDescription("Unified tool for managing Apache Pulsar Functions Worker resources.\n"+
			"This tool provides access to various functions worker resources, including:\n"+
			"1. Function statistics (resource=function_stats): Get statistics for all functions running on this functions worker\n"+
			"2. Monitoring metrics (resource=monitoring_metrics): Get metrics for monitoring function workers\n"+
			"3. Cluster information (resource=cluster): Get information about the function worker cluster\n"+
			"4. Cluster leader (resource=cluster_leader): Get the leader of the worker cluster\n"+
			"5. Function assignments (resource=function_assignments): Get the assignments of functions across the worker cluster\n\n"+
			"Examples:\n"+
			"- {\"resource\": \"function_stats\"} returns statistics for all functions\n"+
			"- {\"resource\": \"cluster\"} returns all workers in the cluster\n"+
			"This tool requires Pulsar super-user permissions."),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description("Type of functions worker resource to access, available options:\n"+
				"- function_stats: Statistics for all functions running on the functions worker\n"+
				"- monitoring_metrics: Metrics for monitoring function workers\n"+
				"- cluster: Information about all workers in the functions worker cluster\n"+
				"- cluster_leader: Information about the leader of the functions worker cluster\n"+
				"- function_assignments: Assignments of functions across the functions worker cluster"),
		),
	)
	s.AddTool(functionsWorkerTool, handleFunctionsWorkerTool(readOnly))
}

// handleFunctionsWorkerTool returns a function to handle functions worker tool requests
func handleFunctionsWorkerTool(_ bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create the admin client
		admin := pulsar.AdminClient

		// Get required resource parameter
		resource, err := common.RequiredParam[string](request.Params.Arguments, "resource")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'resource'. " +
				"Please specify one of: function_stats, monitoring_metrics, cluster, cluster_leader, function_assignments.")), nil
		}

		// Process request based on resource type
		switch resource {
		case "function_stats":
			return handleFunctionsWorkerFunctionStats(admin)
		case "monitoring_metrics":
			return handleFunctionsWorkerMonitoringMetrics(admin)
		case "cluster":
			return handleFunctionsWorkerGetCluster(admin)
		case "cluster_leader":
			return handleFunctionsWorkerGetClusterLeader(admin)
		case "function_assignments":
			return handleFunctionsWorkerGetFunctionAssignments(admin)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Unsupported resource: %s. "+
				"Please use one of: function_stats, monitoring_metrics, cluster, cluster_leader, function_assignments.", resource)), nil
		}
	}
}

// handleFunctionsWorkerFunctionStats handles retrieving function statistics
func handleFunctionsWorkerFunctionStats(admin cmdutils.Client) (*mcp.CallToolResult, error) {
	// Get function stats
	stats, err := admin.FunctionsWorker().GetFunctionsStats()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get functions stats: %v", err)), nil
	}

	// Format the output
	jsonBytes, err := json.Marshal(stats)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal functions stats: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleFunctionsWorkerMonitoringMetrics handles retrieving monitoring metrics
func handleFunctionsWorkerMonitoringMetrics(admin cmdutils.Client) (*mcp.CallToolResult, error) {
	// Get monitoring metrics
	metrics, err := admin.FunctionsWorker().GetMetrics()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get monitoring metrics: %v", err)), nil
	}

	// Format the output
	jsonBytes, err := json.Marshal(metrics)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal monitoring metrics: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleFunctionsWorkerGetCluster handles retrieving cluster information
func handleFunctionsWorkerGetCluster(admin cmdutils.Client) (*mcp.CallToolResult, error) {
	// Get cluster info
	cluster, err := admin.FunctionsWorker().GetCluster()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get worker cluster: %v", err)), nil
	}

	// Format the output
	jsonBytes, err := json.Marshal(cluster)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal worker cluster info: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleFunctionsWorkerGetClusterLeader handles retrieving cluster leader information
func handleFunctionsWorkerGetClusterLeader(admin cmdutils.Client) (*mcp.CallToolResult, error) {
	// Get cluster leader
	leader, err := admin.FunctionsWorker().GetClusterLeader()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get worker cluster leader: %v", err)), nil
	}

	// Format the output
	jsonBytes, err := json.Marshal(leader)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal worker cluster leader info: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleFunctionsWorkerGetFunctionAssignments handles retrieving function assignments
func handleFunctionsWorkerGetFunctionAssignments(admin cmdutils.Client) (*mcp.CallToolResult, error) {
	// Get function assignments
	assignments, err := admin.FunctionsWorker().GetAssignments()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get function assignments: %v", err)), nil
	}

	// Format the output
	jsonBytes, err := json.Marshal(assignments)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal function assignments: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}
