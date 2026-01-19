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

	"github.com/google/jsonschema-go/jsonschema"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

type pulsarAdminFunctionsWorkerInput struct {
	Resource string `json:"resource"`
}

const (
	pulsarAdminFunctionsWorkerToolDesc = "Unified tool for managing Apache Pulsar Functions Worker resources. " +
		"Pulsar Functions is a serverless compute framework that allows you to process messages in a streaming fashion. " +
		"The Functions Worker is the runtime environment that executes and manages Pulsar Functions. " +
		"This tool provides comprehensive access to functions worker resources including function statistics, " +
		"monitoring metrics, cluster information, leader election status, and function assignments across the cluster. " +
		"Functions workers can be deployed in multiple modes (standalone, cluster) and this tool helps monitor " +
		"and manage the worker cluster state, performance metrics, and function distribution. " +
		"Most operations require Pulsar super-user permissions for security reasons."
	pulsarAdminFunctionsWorkerResourceDesc = "Type of functions worker resource to access. Available resources:\n" +
		"- function_stats: Statistics for all functions running on the functions worker, including processing rates, error counts, and resource usage\n" +
		"- monitoring_metrics: Comprehensive metrics for monitoring function workers, including JVM metrics, resource utilization, and performance indicators\n" +
		"- cluster: Information about all workers in the functions worker cluster, including their status, capabilities, and workload distribution\n" +
		"- cluster_leader: Information about the leader of the functions worker cluster, essential for understanding cluster coordination\n" +
		"- function_assignments: Current assignments of functions across the functions worker cluster, showing which functions are running on which workers"
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
func (b *PulsarAdminFunctionsWorkerToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool, err := b.buildFunctionsWorkerTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildFunctionsWorkerHandler(config.ReadOnly)

	return []builders.ToolDefinition{
		builders.ServerTool[pulsarAdminFunctionsWorkerInput, any]{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildFunctionsWorkerTool builds the Pulsar Admin Functions Worker MCP tool definition
// Migrated from the original tool definition logic
func (b *PulsarAdminFunctionsWorkerToolBuilder) buildFunctionsWorkerTool() (*sdk.Tool, error) {
	inputSchema, err := buildPulsarAdminFunctionsWorkerInputSchema()
	if err != nil {
		return nil, err
	}

	return &sdk.Tool{
		Name:        "pulsar_admin_functions_worker",
		Description: pulsarAdminFunctionsWorkerToolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildFunctionsWorkerHandler builds the Pulsar Admin Functions Worker handler function
// Migrated from the original handler logic
func (b *PulsarAdminFunctionsWorkerToolBuilder) buildFunctionsWorkerHandler(_ bool) builders.ToolHandlerFunc[pulsarAdminFunctionsWorkerInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input pulsarAdminFunctionsWorkerInput) (*sdk.CallToolResult, any, error) {
		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return nil, nil, fmt.Errorf("pulsar session not found in context")
		}

		// Create the admin client
		admin, err := session.GetAdminClient()
		if err != nil {
			return nil, nil, b.handleError("get admin client", err)
		}

		// Get required resource parameter
		resource := input.Resource
		if resource == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'resource'; please specify one of: function_stats, monitoring_metrics, cluster, cluster_leader, function_assignments")
		}

		// Process request based on resource type
		switch resource {
		case "function_stats":
			result, handlerErr := b.handleFunctionsWorkerFunctionStats(admin)
			return result, nil, handlerErr
		case "monitoring_metrics":
			result, handlerErr := b.handleFunctionsWorkerMonitoringMetrics(admin)
			return result, nil, handlerErr
		case "cluster":
			result, handlerErr := b.handleFunctionsWorkerGetCluster(admin)
			return result, nil, handlerErr
		case "cluster_leader":
			result, handlerErr := b.handleFunctionsWorkerGetClusterLeader(admin)
			return result, nil, handlerErr
		case "function_assignments":
			result, handlerErr := b.handleFunctionsWorkerGetFunctionAssignments(admin)
			return result, nil, handlerErr
		default:
			return nil, nil, fmt.Errorf("unsupported resource: %s. please use one of: function_stats, monitoring_metrics, cluster, cluster_leader, function_assignments", resource)
		}
	}
}

// Unified error handling and utility functions

// handleError provides unified error handling
func (b *PulsarAdminFunctionsWorkerToolBuilder) handleError(operation string, err error) error {
	return fmt.Errorf("failed to %s: %v", operation, err)
}

// marshalResponse provides unified JSON serialization for responses
func (b *PulsarAdminFunctionsWorkerToolBuilder) marshalResponse(data interface{}) (*sdk.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, b.handleError("marshal response", err)
	}
	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: string(jsonBytes)}},
	}, nil
}

// Operation handler functions - migrated from the original implementation

// handleFunctionsWorkerFunctionStats handles retrieving function statistics
func (b *PulsarAdminFunctionsWorkerToolBuilder) handleFunctionsWorkerFunctionStats(admin cmdutils.Client) (*sdk.CallToolResult, error) {
	// Get function stats
	stats, err := admin.FunctionsWorker().GetFunctionsStats()
	if err != nil {
		return nil, b.handleError("get functions stats", err)
	}

	return b.marshalResponse(stats)
}

// handleFunctionsWorkerMonitoringMetrics handles retrieving monitoring metrics
func (b *PulsarAdminFunctionsWorkerToolBuilder) handleFunctionsWorkerMonitoringMetrics(admin cmdutils.Client) (*sdk.CallToolResult, error) {
	// Get monitoring metrics
	metrics, err := admin.FunctionsWorker().GetMetrics()
	if err != nil {
		return nil, b.handleError("get monitoring metrics", err)
	}

	return b.marshalResponse(metrics)
}

// handleFunctionsWorkerGetCluster handles retrieving cluster information
func (b *PulsarAdminFunctionsWorkerToolBuilder) handleFunctionsWorkerGetCluster(admin cmdutils.Client) (*sdk.CallToolResult, error) {
	// Get cluster info
	cluster, err := admin.FunctionsWorker().GetCluster()
	if err != nil {
		return nil, b.handleError("get worker cluster", err)
	}

	return b.marshalResponse(cluster)
}

// handleFunctionsWorkerGetClusterLeader handles retrieving cluster leader information
func (b *PulsarAdminFunctionsWorkerToolBuilder) handleFunctionsWorkerGetClusterLeader(admin cmdutils.Client) (*sdk.CallToolResult, error) {
	// Get cluster leader
	leader, err := admin.FunctionsWorker().GetClusterLeader()
	if err != nil {
		return nil, b.handleError("get worker cluster leader", err)
	}

	return b.marshalResponse(leader)
}

// handleFunctionsWorkerGetFunctionAssignments handles retrieving function assignments
func (b *PulsarAdminFunctionsWorkerToolBuilder) handleFunctionsWorkerGetFunctionAssignments(admin cmdutils.Client) (*sdk.CallToolResult, error) {
	// Get function assignments
	assignments, err := admin.FunctionsWorker().GetAssignments()
	if err != nil {
		return nil, b.handleError("get function assignments", err)
	}

	return b.marshalResponse(assignments)
}

func buildPulsarAdminFunctionsWorkerInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[pulsarAdminFunctionsWorkerInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}

	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}

	setSchemaDescription(schema, "resource", pulsarAdminFunctionsWorkerResourceDesc)
	normalizeAdditionalProperties(schema)
	return schema, nil
}
