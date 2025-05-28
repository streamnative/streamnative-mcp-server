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

// PulsarAdminAddBrokerStatsTools adds broker-stats related tools to the MCP server
func PulsarAdminAddBrokerStatsTools(s *server.MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeaturePulsarAdminBrokersStatus)) && !slices.Contains(features, string(FeatureAll)) && !slices.Contains(features, string(FeatureAllPulsar)) && !slices.Contains(features, string(FeaturePulsarAdmin)) {
		return
	}
	// Unified broker stats tool
	brokerStatsTool := mcp.NewTool("pulsar_admin_broker_stats",
		mcp.WithDescription("Unified tool for retrieving Apache Pulsar broker statistics.\n"+
			"This tool provides access to various broker stats resources, including:\n"+
			"1. Monitoring metrics (resource=monitoring_metrics): Metrics for the broker's monitoring system\n"+
			"2. MBean stats (resource=mbeans): JVM MBeans statistics\n"+
			"3. Topics stats (resource=topics): Statistics about all topics managed by the broker\n"+
			"4. Allocator stats (resource=allocator_stats): Memory allocator statistics for specific allocator\n"+
			"5. Load report (resource=load_report): Broker load information, sometimes the load report is not available, so suggest to use other resources to get the broker metrics\n\n"+
			"Example: {\"resource\": \"monitoring_metrics\"} retrieves all monitoring metrics\n"+
			"Example: {\"resource\": \"allocator_stats\", \"allocator_name\": \"default\"} retrieves stats for the default allocator\n"+
			"This tool requires Pulsar super-user permissions."),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description("Type of broker stats resource to access, available options:\n"+
				"- monitoring_metrics: Metrics for the broker's monitoring system\n"+
				"- mbeans: JVM MBeans statistics\n"+
				"- topics: Statistics about all topics managed by the broker\n"+
				"- allocator_stats: Memory allocator statistics (requires allocator_name parameter)\n"+
				"- load_report: Broker load information"),
		),
		mcp.WithString("allocator_name",
			mcp.Description("The name of the allocator to get statistics for. Required only when resource=allocator_stats"),
		),
	)
	s.AddTool(brokerStatsTool, handleBrokerStats(readOnly))
}

// handleBrokerStats returns a function to handle broker stats requests
func handleBrokerStats(_ bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := pulsar.GetAdminClient()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
		}

		// Get required resource parameter
		resource, err := common.RequiredParam[string](request.Params.Arguments, "resource")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'resource'. " +
				"Please specify one of: monitoring_metrics, mbeans, topics, allocator_stats, load_report.")), nil
		}

		// Process request based on resource type
		switch resource {
		case "monitoring_metrics":
			return handleMonitoringMetrics(client)
		case "mbeans":
			return handleMBeans(client)
		case "topics":
			return handleTopics(client)
		case "allocator_stats":
			allocatorName, err := common.RequiredParam[string](request.Params.Arguments, "allocator_name")
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'allocator_name' for allocator_stats resource. " +
					"Please provide the name of the allocator to get statistics for.")), nil
			}
			return handleAllocatorStats(client, allocatorName)
		case "load_report":
			return handleLoadReport(client)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Unsupported resource: %s. "+
				"Please use one of: monitoring_metrics, mbeans, topics, allocator_stats, load_report.", resource)), nil
		}
	}
}

// handleMonitoringMetrics handles retrieving monitoring metrics
func handleMonitoringMetrics(client cmdutils.Client) (*mcp.CallToolResult, error) {
	// Get monitoring metrics
	metrics, err := client.BrokerStats().GetMetrics()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get monitoring metrics: %v", err)), nil
	}

	// Convert result to JSON string
	metricsJSON, err := json.Marshal(metrics)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize monitoring metrics: %v", err)), nil
	}

	return mcp.NewToolResultText(string(metricsJSON)), nil
}

// handleMBeans handles retrieving mbean stats
func handleMBeans(client cmdutils.Client) (*mcp.CallToolResult, error) {
	// Get mbean stats
	mbeans, err := client.BrokerStats().GetMBeans()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get mbean stats: %v", err)), nil
	}

	// Convert result to JSON string
	mbeansJSON, err := json.Marshal(mbeans)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize mbean stats: %v", err)), nil
	}

	return mcp.NewToolResultText(string(mbeansJSON)), nil
}

// handleTopics handles retrieving topics stats
func handleTopics(client cmdutils.Client) (*mcp.CallToolResult, error) {
	// Get topics stats
	topics, err := client.BrokerStats().GetTopics()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topics stats: %v", err)), nil
	}

	// Convert result to JSON string
	topicsJSON, err := json.Marshal(topics)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize topics stats: %v", err)), nil
	}

	return mcp.NewToolResultText(string(topicsJSON)), nil
}

// handleAllocatorStats handles retrieving allocator stats
func handleAllocatorStats(client cmdutils.Client, allocatorName string) (*mcp.CallToolResult, error) {
	// Get allocator stats
	stats, err := client.BrokerStats().GetAllocatorStats(allocatorName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get allocator stats: %v", err)), nil
	}

	// Convert result to JSON string
	statsJSON, err := json.Marshal(stats)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize allocator stats: %v", err)), nil
	}

	return mcp.NewToolResultText(string(statsJSON)), nil
}

// handleLoadReport handles retrieving broker load report
func handleLoadReport(client cmdutils.Client) (*mcp.CallToolResult, error) {
	// Get load report
	loadReport, err := client.BrokerStats().GetLoadReport()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get broker load report: %v", err)), nil
	}

	// Convert result to JSON string
	loadReportJSON, err := json.Marshal(loadReport)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize broker load report: %v", err)), nil
	}

	return mcp.NewToolResultText(string(loadReportJSON)), nil
}
