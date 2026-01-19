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

type pulsarAdminBrokerStatsInput struct {
	Resource      string  `json:"resource"`
	AllocatorName *string `json:"allocator_name,omitempty"`
}

const (
	pulsarAdminBrokerStatsResourceDesc = "Type of broker stats resource to access, available options:\n" +
		"- monitoring_metrics: Metrics for the broker's monitoring system\n" +
		"- mbeans: JVM MBeans statistics\n" +
		"- topics: Statistics about all topics managed by the broker\n" +
		"- allocator_stats: Memory allocator statistics (requires allocator_name parameter)\n" +
		"- load_report: Broker load information"
	pulsarAdminBrokerStatsAllocatorNameDesc = "The name of the allocator to get statistics for. Required only when resource=allocator_stats"
)

// PulsarAdminBrokerStatsToolBuilder implements the ToolBuilder interface for Pulsar Broker Statistics
// /nolint:revive
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
func (b *PulsarAdminBrokerStatsToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool, err := b.buildBrokerStatsTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildBrokerStatsHandler(config.ReadOnly)

	return []builders.ToolDefinition{
		builders.ServerTool[pulsarAdminBrokerStatsInput, any]{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildBrokerStatsTool builds the Pulsar Admin Broker Stats MCP tool definition
func (b *PulsarAdminBrokerStatsToolBuilder) buildBrokerStatsTool() (*sdk.Tool, error) {
	inputSchema, err := buildPulsarAdminBrokerStatsInputSchema()
	if err != nil {
		return nil, err
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

	return &sdk.Tool{
		Name:        "pulsar_admin_broker_stats",
		Description: toolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildBrokerStatsHandler builds the Pulsar Admin Broker Stats handler function
func (b *PulsarAdminBrokerStatsToolBuilder) buildBrokerStatsHandler(_ bool) builders.ToolHandlerFunc[pulsarAdminBrokerStatsInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input pulsarAdminBrokerStatsInput) (*sdk.CallToolResult, any, error) {
		// Get Pulsar admin client
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return nil, nil, fmt.Errorf("pulsar session not found in context")
		}
		client, err := session.GetAdminClient()
		if err != nil {
			return nil, nil, b.handleError("get admin client", err)
		}

		// Get required resource parameter
		resource := input.Resource
		if resource == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'resource'; please specify one of: monitoring_metrics, mbeans, topics, allocator_stats, load_report")
		}

		// Process request based on resource type
		switch resource {
		case "monitoring_metrics":
			result, handlerErr := b.handleMonitoringMetrics(client)
			return result, nil, handlerErr
		case "mbeans":
			result, handlerErr := b.handleMBeans(client)
			return result, nil, handlerErr
		case "topics":
			result, handlerErr := b.handleTopics(client)
			return result, nil, handlerErr
		case "allocator_stats":
			allocatorName := ""
			if input.AllocatorName != nil {
				allocatorName = *input.AllocatorName
			}
			if allocatorName == "" {
				return nil, nil, fmt.Errorf("missing required parameter 'allocator_name' for allocator_stats resource; please provide the name of the allocator to get statistics for")
			}
			result, handlerErr := b.handleAllocatorStats(client, allocatorName)
			return result, nil, handlerErr
		case "load_report":
			result, handlerErr := b.handleLoadReport(client)
			return result, nil, handlerErr
		default:
			return nil, nil, fmt.Errorf("unsupported resource: %s. please use one of: monitoring_metrics, mbeans, topics, allocator_stats, load_report", resource)
		}
	}
}

// Utility functions

// handleError provides unified error handling
func (b *PulsarAdminBrokerStatsToolBuilder) handleError(operation string, err error) error {
	return fmt.Errorf("failed to %s: %v", operation, err)
}

// marshalResponse provides unified JSON serialization for responses
func (b *PulsarAdminBrokerStatsToolBuilder) marshalResponse(data interface{}) (*sdk.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, b.handleError("marshal response", err)
	}
	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: string(jsonBytes)}},
	}, nil
}

// Specific operation handler functions

// handleMonitoringMetrics handles retrieving monitoring metrics
func (b *PulsarAdminBrokerStatsToolBuilder) handleMonitoringMetrics(client cmdutils.Client) (*sdk.CallToolResult, error) {
	stats, err := client.BrokerStats().GetMetrics()
	if err != nil {
		return nil, b.handleError("get monitoring metrics", err)
	}
	return b.marshalResponse(stats)
}

// handleMBeans handles retrieving MBeans statistics
func (b *PulsarAdminBrokerStatsToolBuilder) handleMBeans(client cmdutils.Client) (*sdk.CallToolResult, error) {
	stats, err := client.BrokerStats().GetMBeans()
	if err != nil {
		return nil, b.handleError("get MBeans", err)
	}
	return b.marshalResponse(stats)
}

// handleTopics handles retrieving topics statistics
func (b *PulsarAdminBrokerStatsToolBuilder) handleTopics(client cmdutils.Client) (*sdk.CallToolResult, error) {
	stats, err := client.BrokerStats().GetTopics()
	if err != nil {
		return nil, b.handleError("get topics stats", err)
	}
	return b.marshalResponse(stats)
}

// handleAllocatorStats handles retrieving allocator statistics
func (b *PulsarAdminBrokerStatsToolBuilder) handleAllocatorStats(client cmdutils.Client, allocatorName string) (*sdk.CallToolResult, error) {
	stats, err := client.BrokerStats().GetAllocatorStats(allocatorName)
	if err != nil {
		return nil, b.handleError("get allocator stats", err)
	}
	return b.marshalResponse(stats)
}

// handleLoadReport handles retrieving load report
func (b *PulsarAdminBrokerStatsToolBuilder) handleLoadReport(client cmdutils.Client) (*sdk.CallToolResult, error) {
	stats, err := client.BrokerStats().GetLoadReport()
	if err != nil {
		return nil, b.handleError("get load report", err)
	}
	return b.marshalResponse(stats)
}

func buildPulsarAdminBrokerStatsInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[pulsarAdminBrokerStatsInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}

	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}

	setSchemaDescription(schema, "resource", pulsarAdminBrokerStatsResourceDesc)
	setSchemaDescription(schema, "allocator_name", pulsarAdminBrokerStatsAllocatorNameDesc)

	normalizeAdditionalProperties(schema)
	return schema, nil
}
