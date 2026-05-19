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

// Package kafka provides Kafka MCP tool builders.
package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

var kafkaConnectOperationSpecs = builders.OperationRegistry{
	{Name: "list", Mode: builders.OperationModeRead},
	{Name: "get", Mode: builders.OperationModeRead},
	{Name: "create", Mode: builders.OperationModeWrite, Destructive: true},
	{Name: "update", Mode: builders.OperationModeWrite, Destructive: true},
	{Name: "delete", Mode: builders.OperationModeWrite, Destructive: true},
	{Name: "restart", Mode: builders.OperationModeWrite, Destructive: true},
	{Name: "pause", Mode: builders.OperationModeWrite, Destructive: true},
	{Name: "resume", Mode: builders.OperationModeWrite, Destructive: true},
}

// KafkaConnectToolBuilder implements the ToolBuilder interface for Kafka Connect
// It provides functionality to build Kafka Connect administration tools
// /nolint:revive
type KafkaConnectToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewKafkaConnectToolBuilder creates a new Kafka Connect tool builder instance
func NewKafkaConnectToolBuilder() *KafkaConnectToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "kafka_connect",
		Version:     "1.0.0",
		Description: "Kafka Connect administration tools",
		Category:    "kafka_admin",
		Tags:        []string{"kafka", "connect", "admin"},
	}

	features := []string{
		"kafka-admin-kafka-connect",
		"all",
		"all-kafka",
		"kafka-admin",
	}

	return &KafkaConnectToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Kafka Connect tool list
// This is the core method implementing the ToolBuilder interface
func (b *KafkaConnectToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	tools := []server.ServerTool{
		{
			Tool:    b.buildKafkaConnectTool(toolModeRead),
			Handler: b.buildKafkaConnectHandler(toolModeRead),
		},
	}
	if !config.ReadOnly {
		tools = append(tools, server.ServerTool{
			Tool:    b.buildKafkaConnectTool(toolModeWrite),
			Handler: b.buildKafkaConnectHandler(toolModeWrite),
		})
	}

	return tools, nil
}

// buildKafkaConnectTool builds the Kafka Connect MCP tool definition
// Migrated from the original tool definition logic
func (b *KafkaConnectToolBuilder) buildKafkaConnectTool(mode toolMode) mcp.Tool {
	resourceDesc := "Resource to operate on. Available resources:\n" +
		"- kafka-connect-cluster: A single Kafka Connect cluster that manages connectors and tasks.\n" +
		"- connector: A single Kafka Connect connector instance that moves data between Kafka and external systems.\n" +
		"- connectors: Collection of all Kafka Connect connectors in a cluster.\n" +
		"- connector-plugins: Collection of all Kafka Connect connector plugins, StreamNative Cloud provides a set of built-in connectors via this resource."
	resourceEnum := []string{"kafka-connect-cluster", "connector", "connectors", "connector-plugins"}

	operationDesc := "Operation to perform. Available operations:\n" +
		"- list: List all connectors or connector plugins in a cluster.\n" +
		"- get: Retrieve detailed information about a Kafka Connect cluster or specific connector."
	operationEnum := kafkaConnectOperationSpecs.NamesForMode(mode)
	toolName := "kafka_admin_connect_read"
	annotation := builders.ToolAnnotationForMode(mode, "Read Kafka Connect", "Manage Kafka Connect", kafkaConnectOperationSpecs)
	if isToolModeWrite(mode) {
		resourceDesc = "Resource to operate on. Available resources:\n" +
			"- connector: A single Kafka Connect connector instance that moves data between Kafka and external systems."
		resourceEnum = []string{"connector"}
		operationDesc = "Operation to perform. Available operations:\n" +
			"- create: Create a new connector with specified configuration.\n" +
			"- update: Modify an existing connector's configuration.\n" +
			"- delete: Remove a connector from the Kafka Connect cluster.\n" +
			"- restart: Restart a running connector (useful after failures or configuration changes).\n" +
			"- pause: Temporarily stop a connector from processing data.\n" +
			"- resume: Continue processing with a previously paused connector."
		operationEnum = kafkaConnectOperationSpecs.NamesForMode(mode)
		toolName = "kafka_admin_connect_write"
	}

	toolDesc := "Read Apache Kafka Connect cluster, connector, and plugin information.\n" +
		"Kafka Connect is a framework for connecting Kafka with external systems through reusable connectors.\n" +
		"This read-only tool lists connectors and plugins and retrieves cluster or connector details.\n\n" +
		"Usage Examples:\n\n" +
		"1. List all connectors in the Kafka Connect cluster:\n" +
		"   resource: \"connectors\"\n" +
		"   operation: \"list\"\n\n" +
		"2. Get information about a specific connector:\n" +
		"   resource: \"connector\"\n" +
		"   operation: \"get\"\n" +
		"   name: \"my-jdbc-source\"\n\n" +
		"3. List all available connector plugins:\n" +
		"   resource: \"connector-plugins\"\n" +
		"   operation: \"list\"\n\n" +
		"4. Get information about the Kafka Connect cluster:\n" +
		"   resource: \"kafka-connect-cluster\"\n" +
		"   operation: \"get\"\n\n" +
		"This tool requires appropriate Kafka Connect read permissions."
	if isToolModeWrite(mode) {
		toolDesc = "Manage Apache Kafka Connect connectors.\n" +
			"This write tool creates, updates, deletes, restarts, pauses, or resumes connectors.\n\n" +
			"Usage Examples:\n\n" +
			"1. Create a new connector:\n" +
			"   resource: \"connector\"\n" +
			"   operation: \"create\"\n" +
			"   name: \"my-jdbc-source\"\n" +
			"   config: {...}\n\n" +
			"2. Update an existing connector's configuration:\n" +
			"   resource: \"connector\"\n" +
			"   operation: \"update\"\n" +
			"   name: \"my-jdbc-source\"\n" +
			"   config: {...}\n\n" +
			"3. Restart a connector:\n" +
			"   resource: \"connector\"\n" +
			"   operation: \"restart\"\n" +
			"   name: \"my-jdbc-source\"\n\n" +
			"This tool requires appropriate Kafka Connect write permissions."
	}

	tool := mcp.NewTool(toolName,
		mcp.WithDescription(toolDesc),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description(resourceDesc),
			mcp.Enum(resourceEnum...),
		),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
			mcp.Enum(operationEnum...),
		),
		mcp.WithString("name",
			mcp.Description("The name of the Kafka Connect connector to operate on. Required for operations that target one connector. "+
				"Must be unique within the Kafka Connect cluster. "+
				"Should be descriptive of the connector's purpose, such as 'mysql-inventory-source' or 'elasticsearch-logs-sink'.")),
		mcp.WithObject("config",
			mcp.Description("The configuration settings for the connector. "+
				"Required for 'create' and 'update' operations on the 'connector' resource. "+
				"Must include 'connector.class' which specifies the connector implementation. "+
				"Common configurations include:\n"+
				"- connector.class: The Java class implementing the connector\n"+
				"- tasks.max: Maximum number of tasks to use for this connector\n"+
				"- topics/topic.regex/topic.prefix: Topic specification (varies by connector)\n"+
				"- key.converter/value.converter: Data format converters\n"+
				"- transforms: Optional transformations to apply to data\n"+
				"Additional fields depend on the specific connector type being used.")),
		annotation,
	)
	if isToolModeWrite(mode) {
		pruneToolInputSchema(&tool, []string{"resource", "operation", "name", "config"})
	} else {
		pruneToolInputSchema(&tool, []string{"resource", "operation", "name"})
	}
	return tool
}

// buildKafkaConnectHandler builds the Kafka Connect handler function
// Migrated from the original handler logic
func (b *KafkaConnectToolBuilder) buildKafkaConnectHandler(mode toolMode) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get required parameters
		resource, err := request.RequireString("resource")
		if err != nil {
			return b.handleError("get resource", err), nil
		}

		operation, err := request.RequireString("operation")
		if err != nil {
			return b.handleError("get operation", err), nil
		}

		// Normalize parameters
		resource = strings.ToLower(resource)
		operation = strings.ToLower(operation)

		if err := validateModeOperation(mode, operation, kafkaConnectOperationSpecs); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get Kafka Connect client
		session := mcpCtx.GetKafkaSession(ctx)
		if session == nil {
			return mcp.NewToolResultError("Kafka session not found in context"), nil
		}
		admin, err := session.GetConnectClient()
		if err != nil {
			return b.handleError("get admin client", err), nil
		}

		// Dispatch based on resource and operation
		switch resource {
		case "kafka-connect-cluster":
			switch operation {
			case "get":
				return b.handleKafkaConnectClusterGet(ctx, admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'kafka-connect-cluster': %s", operation)), nil
			}
		case "connectors":
			switch operation {
			case "list":
				return b.handleKafkaConnectorsList(ctx, admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'connectors': %s", operation)), nil
			}
		case "connector":
			switch operation {
			case "get":
				return b.handleKafkaConnectorGet(ctx, admin, request)
			case "create":
				return b.handleKafkaConnectorCreate(ctx, admin, request)
			case "update":
				return b.handleKafkaConnectorUpdate(ctx, admin, request)
			case "delete":
				return b.handleKafkaConnectorDelete(ctx, admin, request)
			case "restart":
				return b.handleKafkaConnectorRestart(ctx, admin, request)
			case "pause":
				return b.handleKafkaConnectorPause(ctx, admin, request)
			case "resume":
				return b.handleKafkaConnectorResume(ctx, admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'connector': %s", operation)), nil
			}
		case "connector-plugins":
			switch operation {
			case "list":
				return b.handleKafkaConnectorPluginsList(ctx, admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'connector-plugins': %s", operation)), nil
			}
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resource: %s. Available resources: kafka-connect-cluster, connectors, connector, connector-plugins", resource)), nil
		}
	}
}

// Unified error handling and utility functions

// handleError provides unified error handling
func (b *KafkaConnectToolBuilder) handleError(operation string, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("Failed to %s: %v", operation, err))
}

// marshalResponse provides unified JSON serialization for responses
func (b *KafkaConnectToolBuilder) marshalResponse(data interface{}) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return b.handleError("marshal response", err), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// Specific operation handler functions - migrated from original implementation

// handleKafkaConnectClusterGet handles retrieving Kafka Connect cluster information
func (b *KafkaConnectToolBuilder) handleKafkaConnectClusterGet(ctx context.Context, admin kafka.Connect, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cluster, err := admin.GetInfo(ctx)
	if err != nil {
		return b.handleError("get Kafka Connect cluster", err), nil
	}
	return b.marshalResponse(cluster)
}

// handleKafkaConnectorsList handles listing all connectors
func (b *KafkaConnectToolBuilder) handleKafkaConnectorsList(ctx context.Context, admin kafka.Connect, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	connectors, err := admin.ListConnectors(ctx)
	if err != nil {
		return b.handleError("get Kafka Connect connectors", err), nil
	}
	return b.marshalResponse(connectors)
}

// handleKafkaConnectorGet handles retrieving specific connector information
func (b *KafkaConnectToolBuilder) handleKafkaConnectorGet(ctx context.Context, admin kafka.Connect, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return b.handleError("get connector name", err), nil
	}

	connector, err := admin.GetConnector(ctx, name)
	if err != nil {
		return b.handleError("get Kafka Connect connector", err), nil
	}
	return b.marshalResponse(connector)
}

// handleKafkaConnectorCreate handles creating a new connector
func (b *KafkaConnectToolBuilder) handleKafkaConnectorCreate(ctx context.Context, admin kafka.Connect, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return b.handleError("get connector name", err), nil
	}

	configMap, err := common.RequiredParamObject(request.GetArguments(), "config")
	if err != nil {
		return b.handleError("get config", err), nil
	}

	config := common.ConvertToMapString(configMap)
	config["name"] = name

	connector, err := admin.CreateConnector(ctx, name, config)
	if err != nil {
		return b.handleError("create Kafka Connect connector", err), nil
	}
	return b.marshalResponse(connector)
}

// handleKafkaConnectorUpdate handles updating a connector
func (b *KafkaConnectToolBuilder) handleKafkaConnectorUpdate(ctx context.Context, admin kafka.Connect, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return b.handleError("get connector name", err), nil
	}

	configMap, err := common.RequiredParamObject(request.GetArguments(), "config")
	if err != nil {
		return b.handleError("get config", err), nil
	}

	config := common.ConvertToMapString(configMap)
	config["name"] = name

	connector, err := admin.UpdateConnector(ctx, name, config)
	if err != nil {
		return b.handleError("update Kafka Connect connector", err), nil
	}
	return b.marshalResponse(connector)
}

// handleKafkaConnectorDelete handles deleting a connector
func (b *KafkaConnectToolBuilder) handleKafkaConnectorDelete(ctx context.Context, admin kafka.Connect, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return b.handleError("get connector name", err), nil
	}

	err = admin.DeleteConnector(ctx, name)
	if err != nil {
		return b.handleError("delete Kafka Connect connector", err), nil
	}

	return mcp.NewToolResultText("Connector deleted successfully"), nil
}

// handleKafkaConnectorRestart handles restarting a connector
func (b *KafkaConnectToolBuilder) handleKafkaConnectorRestart(ctx context.Context, admin kafka.Connect, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return b.handleError("get connector name", err), nil
	}

	err = admin.RestartConnector(ctx, name)
	if err != nil {
		return b.handleError("restart Kafka Connect connector", err), nil
	}

	return mcp.NewToolResultText("Connector restarted successfully"), nil
}

// handleKafkaConnectorPause handles pausing a connector
func (b *KafkaConnectToolBuilder) handleKafkaConnectorPause(ctx context.Context, admin kafka.Connect, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return b.handleError("get connector name", err), nil
	}

	err = admin.PauseConnector(ctx, name)
	if err != nil {
		return b.handleError("pause Kafka Connect connector", err), nil
	}

	return mcp.NewToolResultText("Connector paused successfully"), nil
}

// handleKafkaConnectorResume handles resuming a connector
func (b *KafkaConnectToolBuilder) handleKafkaConnectorResume(ctx context.Context, admin kafka.Connect, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return b.handleError("get connector name", err), nil
	}

	err = admin.ResumeConnector(ctx, name)
	if err != nil {
		return b.handleError("resume Kafka Connect connector", err), nil
	}

	return mcp.NewToolResultText("Connector resumed successfully"), nil
}

// handleKafkaConnectorPluginsList handles listing connector plugins
func (b *KafkaConnectToolBuilder) handleKafkaConnectorPluginsList(ctx context.Context, admin kafka.Connect, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	plugins, err := admin.ListPlugins(ctx)
	if err != nil {
		return b.handleError("get Kafka Connect connector plugins", err), nil
	}
	return b.marshalResponse(plugins)
}
