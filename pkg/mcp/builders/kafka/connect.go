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

// Package kafka provides Kafka MCP tool builders.
package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

type kafkaConnectInput struct {
	Resource  string                 `json:"resource"`
	Operation string                 `json:"operation"`
	Name      *string                `json:"name,omitempty"`
	Config    map[string]interface{} `json:"config,omitempty"`
}

const (
	kafkaConnectResourceDesc = "Resource to operate on. Available resources:\n" +
		"- kafka-connect-cluster: A single Kafka Connect cluster that manages connectors and tasks.\n" +
		"- connector: A single Kafka Connect connector instance that moves data between Kafka and external systems.\n" +
		"- connectors: Collection of all Kafka Connect connectors in a cluster.\n" +
		"- connector-plugins: Collection of all Kafka Connect connector plugins, StreamNative Cloud provides a set of built-in connectors via this resource."
	kafkaConnectOperationDesc = "Operation to perform. Available operations:\n" +
		"- list: List all connectors or connector plugins in a cluster.\n" +
		"- get: Retrieve detailed information about a Kafka Connect cluster or specific connector.\n" +
		"- create: Create a new connector with specified configuration.\n" +
		"- update: Modify an existing connector's configuration.\n" +
		"- delete: Remove a connector from the Kafka Connect cluster.\n" +
		"- restart: Restart a running connector (useful after failures or configuration changes).\n" +
		"- pause: Temporarily stop a connector from processing data.\n" +
		"- resume: Continue processing with a previously paused connector."
	kafkaConnectNameDesc = "The name of the Kafka Connect connector to operate on. " +
		"Required for 'get', 'create', 'update', 'delete', 'restart', 'pause', and 'resume' operations on the 'connector' resource. " +
		"Must be unique within the Kafka Connect cluster. " +
		"Should be descriptive of the connector's purpose, such as 'mysql-inventory-source' or 'elasticsearch-logs-sink'."
	kafkaConnectConfigDesc = "The configuration settings for the connector. " +
		"Required for 'create' and 'update' operations on the 'connector' resource. " +
		"Must include 'connector.class' which specifies the connector implementation. " +
		"Common configurations include:\n" +
		"- connector.class: The Java class implementing the connector\n" +
		"- tasks.max: Maximum number of tasks to use for this connector\n" +
		"- topics/topic.regex/topic.prefix: Topic specification (varies by connector)\n" +
		"- key.converter/value.converter: Data format converters\n" +
		"- transforms: Optional transformations to apply to data\n" +
		"Additional fields depend on the specific connector type being used."

	kafkaConnectToolDesc = `Unified tool for managing Apache Kafka Connect.
Kafka Connect is a framework for connecting Kafka with external systems such as databases, key-value stores, search indexes, and file systems. It provides a standardized way to stream data in and out of Kafka, without requiring custom integration code.

Key concepts in Kafka Connect:

- Connectors: The high-level abstraction that coordinates data streaming by managing tasks
- Tasks: The implementation of how data is copied to or from Kafka
- Workers: The running processes that execute connectors and tasks
- Plugins: Reusable connector implementations for specific external systems
- Source Connectors: Import data from external systems into Kafka topics
- Sink Connectors: Export data from Kafka topics to external systems

Kafka Connect simplifies data integration, enables scalable and reliable streaming pipelines, and reduces the operational burden of managing data flows.

Usage Examples:

1. List all connectors in the Kafka Connect cluster:
   resource: "connectors"
   operation: "list"

2. Get information about a specific connector:
   resource: "connector"
   operation: "get"
   name: "my-jdbc-source"

3. Create a new JDBC source connector:
   resource: "connector"
   operation: "create"
   name: "my-jdbc-source"
   config: {
     "connector.class": "io.confluent.connect.jdbc.JdbcSourceConnector",
     "connection.url": "jdbc:mysql://mysql:3306/mydb",
     "connection.user": "user",
     "connection.password": "password",
     "topic.prefix": "mysql-",
     "table.whitelist": "users,orders",
     "mode": "incrementing",
     "incrementing.column.name": "id",
     "tasks.max": "1"
   }

4. Update an existing connector's configuration:
   resource: "connector"
   operation: "update"
   name: "my-jdbc-source"
   config: {
     "connector.class": "io.confluent.connect.jdbc.JdbcSourceConnector",
     "tasks.max": "2",
     "table.whitelist": "users,orders,products"
   }

5. Delete a connector:
   resource: "connector"
   operation: "delete"
   name: "my-jdbc-source"

6. Restart a connector after configuration changes or errors:
   resource: "connector"
   operation: "restart"
   name: "my-jdbc-source"

7. Pause a connector temporarily:
   resource: "connector"
   operation: "pause"
   name: "my-jdbc-source"

8. Resume a paused connector:
   resource: "connector"
   operation: "resume"
   name: "my-jdbc-source"

9. List all available connector plugins:
   resource: "connector-plugins"
   operation: "list"

10. Get information about the Kafka Connect cluster:
    resource: "kafka-connect-cluster"
    operation: "get"

This tool requires appropriate Kafka Connect permissions.`
)

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
func (b *KafkaConnectToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool, err := b.buildKafkaConnectTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildKafkaConnectHandler(config.ReadOnly)

	return []builders.ToolDefinition{
		builders.ServerTool[kafkaConnectInput, any]{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildKafkaConnectTool builds the Kafka Connect MCP tool definition
// Migrated from the original tool definition logic
func (b *KafkaConnectToolBuilder) buildKafkaConnectTool() (*sdk.Tool, error) {
	inputSchema, err := buildKafkaConnectInputSchema()
	if err != nil {
		return nil, err
	}

	return &sdk.Tool{
		Name:        "kafka_admin_connect",
		Description: kafkaConnectToolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildKafkaConnectHandler builds the Kafka Connect handler function
// Migrated from the original handler logic
func (b *KafkaConnectToolBuilder) buildKafkaConnectHandler(readOnly bool) builders.ToolHandlerFunc[kafkaConnectInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input kafkaConnectInput) (*sdk.CallToolResult, any, error) {
		// Normalize parameters
		resource := strings.ToLower(input.Resource)
		operation := strings.ToLower(input.Operation)

		// Validate write operations in read-only mode
		if readOnly && (operation == "create" || operation == "update" || operation == "delete" || operation == "restart" || operation == "pause" || operation == "resume") {
			return nil, nil, fmt.Errorf("write operations are not allowed in read-only mode")
		}

		// Get Kafka Connect client
		session := mcpCtx.GetKafkaSession(ctx)
		if session == nil {
			return nil, nil, b.handleError("get Kafka session not found in context", nil)
		}
		admin, err := session.GetConnectClient()
		if err != nil {
			return nil, nil, b.handleError("get admin client", err)
		}

		// Dispatch based on resource and operation
		switch resource {
		case "kafka-connect-cluster":
			switch operation {
			case "get":
				result, err := b.handleKafkaConnectClusterGet(ctx, admin)
				return result, nil, err
			default:
				return nil, nil, fmt.Errorf("invalid operation for resource 'kafka-connect-cluster': %s", operation)
			}
		case "connectors":
			switch operation {
			case "list":
				result, err := b.handleKafkaConnectorsList(ctx, admin)
				return result, nil, err
			default:
				return nil, nil, fmt.Errorf("invalid operation for resource 'connectors': %s", operation)
			}
		case "connector":
			switch operation {
			case "get":
				result, err := b.handleKafkaConnectorGet(ctx, admin, input)
				return result, nil, err
			case "create":
				result, err := b.handleKafkaConnectorCreate(ctx, admin, input)
				return result, nil, err
			case "update":
				result, err := b.handleKafkaConnectorUpdate(ctx, admin, input)
				return result, nil, err
			case "delete":
				result, err := b.handleKafkaConnectorDelete(ctx, admin, input)
				return result, nil, err
			case "restart":
				result, err := b.handleKafkaConnectorRestart(ctx, admin, input)
				return result, nil, err
			case "pause":
				result, err := b.handleKafkaConnectorPause(ctx, admin, input)
				return result, nil, err
			case "resume":
				result, err := b.handleKafkaConnectorResume(ctx, admin, input)
				return result, nil, err
			default:
				return nil, nil, fmt.Errorf("invalid operation for resource 'connector': %s", operation)
			}
		case "connector-plugins":
			switch operation {
			case "list":
				result, err := b.handleKafkaConnectorPluginsList(ctx, admin)
				return result, nil, err
			default:
				return nil, nil, fmt.Errorf("invalid operation for resource 'connector-plugins': %s", operation)
			}
		default:
			return nil, nil, fmt.Errorf("invalid resource: %s. available resources: kafka-connect-cluster, connectors, connector, connector-plugins", resource)
		}
	}
}

// Utility functions

// handleError provides unified error handling
func (b *KafkaConnectToolBuilder) handleError(operation string, err error) error {
	return fmt.Errorf("failed to %s: %v", operation, err)
}

// marshalResponse provides unified JSON serialization for responses
func (b *KafkaConnectToolBuilder) marshalResponse(data interface{}) (*sdk.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, b.handleError("marshal response", err)
	}
	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: string(jsonBytes)}},
	}, nil
}

func requireMap(value map[string]interface{}, key string) (map[string]interface{}, error) {
	if value == nil {
		return nil, fmt.Errorf("required argument %q not found", key)
	}
	return value, nil
}

func buildKafkaConnectInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[kafkaConnectInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}

	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}
	setSchemaDescription(schema, "resource", kafkaConnectResourceDesc)
	setSchemaDescription(schema, "operation", kafkaConnectOperationDesc)
	setSchemaDescription(schema, "name", kafkaConnectNameDesc)
	setSchemaDescription(schema, "config", kafkaConnectConfigDesc)

	normalizeAdditionalProperties(schema)
	return schema, nil
}

// Specific operation handler functions - migrated from original implementation

// handleKafkaConnectClusterGet handles retrieving Kafka Connect cluster information
func (b *KafkaConnectToolBuilder) handleKafkaConnectClusterGet(ctx context.Context, admin kafka.Connect) (*sdk.CallToolResult, error) {
	cluster, err := admin.GetInfo(ctx)
	if err != nil {
		return nil, b.handleError("get Kafka Connect cluster", err)
	}
	return b.marshalResponse(cluster)
}

// handleKafkaConnectorsList handles listing all connectors
func (b *KafkaConnectToolBuilder) handleKafkaConnectorsList(ctx context.Context, admin kafka.Connect) (*sdk.CallToolResult, error) {
	connectors, err := admin.ListConnectors(ctx)
	if err != nil {
		return nil, b.handleError("get Kafka Connect connectors", err)
	}
	return b.marshalResponse(connectors)
}

// handleKafkaConnectorGet handles retrieving specific connector information
func (b *KafkaConnectToolBuilder) handleKafkaConnectorGet(ctx context.Context, admin kafka.Connect, input kafkaConnectInput) (*sdk.CallToolResult, error) {
	name, err := requireString(input.Name, "name")
	if err != nil {
		return nil, b.handleError("get connector name", err)
	}

	connector, err := admin.GetConnector(ctx, name)
	if err != nil {
		return nil, b.handleError("get Kafka Connect connector", err)
	}
	return b.marshalResponse(connector)
}

// handleKafkaConnectorCreate handles creating a new connector
func (b *KafkaConnectToolBuilder) handleKafkaConnectorCreate(ctx context.Context, admin kafka.Connect, input kafkaConnectInput) (*sdk.CallToolResult, error) {
	name, err := requireString(input.Name, "name")
	if err != nil {
		return nil, b.handleError("get connector name", err)
	}

	configMap, err := requireMap(input.Config, "config")
	if err != nil {
		return nil, b.handleError("get config", err)
	}

	config := common.ConvertToMapString(configMap)
	config["name"] = name

	connector, err := admin.CreateConnector(ctx, name, config)
	if err != nil {
		return nil, b.handleError("create Kafka Connect connector", err)
	}
	return b.marshalResponse(connector)
}

// handleKafkaConnectorUpdate handles updating a connector
func (b *KafkaConnectToolBuilder) handleKafkaConnectorUpdate(ctx context.Context, admin kafka.Connect, input kafkaConnectInput) (*sdk.CallToolResult, error) {
	name, err := requireString(input.Name, "name")
	if err != nil {
		return nil, b.handleError("get connector name", err)
	}

	configMap, err := requireMap(input.Config, "config")
	if err != nil {
		return nil, b.handleError("get config", err)
	}

	config := common.ConvertToMapString(configMap)
	config["name"] = name

	connector, err := admin.UpdateConnector(ctx, name, config)
	if err != nil {
		return nil, b.handleError("update Kafka Connect connector", err)
	}
	return b.marshalResponse(connector)
}

// handleKafkaConnectorDelete handles deleting a connector
func (b *KafkaConnectToolBuilder) handleKafkaConnectorDelete(ctx context.Context, admin kafka.Connect, input kafkaConnectInput) (*sdk.CallToolResult, error) {
	name, err := requireString(input.Name, "name")
	if err != nil {
		return nil, b.handleError("get connector name", err)
	}

	if err := admin.DeleteConnector(ctx, name); err != nil {
		return nil, b.handleError("delete Kafka Connect connector", err)
	}

	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: "Connector deleted successfully"}},
	}, nil
}

// handleKafkaConnectorRestart handles restarting a connector
func (b *KafkaConnectToolBuilder) handleKafkaConnectorRestart(ctx context.Context, admin kafka.Connect, input kafkaConnectInput) (*sdk.CallToolResult, error) {
	name, err := requireString(input.Name, "name")
	if err != nil {
		return nil, b.handleError("get connector name", err)
	}

	if err := admin.RestartConnector(ctx, name); err != nil {
		return nil, b.handleError("restart Kafka Connect connector", err)
	}

	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: "Connector restarted successfully"}},
	}, nil
}

// handleKafkaConnectorPause handles pausing a connector
func (b *KafkaConnectToolBuilder) handleKafkaConnectorPause(ctx context.Context, admin kafka.Connect, input kafkaConnectInput) (*sdk.CallToolResult, error) {
	name, err := requireString(input.Name, "name")
	if err != nil {
		return nil, b.handleError("get connector name", err)
	}

	if err := admin.PauseConnector(ctx, name); err != nil {
		return nil, b.handleError("pause Kafka Connect connector", err)
	}

	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: "Connector paused successfully"}},
	}, nil
}

// handleKafkaConnectorResume handles resuming a connector
func (b *KafkaConnectToolBuilder) handleKafkaConnectorResume(ctx context.Context, admin kafka.Connect, input kafkaConnectInput) (*sdk.CallToolResult, error) {
	name, err := requireString(input.Name, "name")
	if err != nil {
		return nil, b.handleError("get connector name", err)
	}

	if err := admin.ResumeConnector(ctx, name); err != nil {
		return nil, b.handleError("resume Kafka Connect connector", err)
	}

	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: "Connector resumed successfully"}},
	}, nil
}

// handleKafkaConnectorPluginsList handles listing connector plugins
func (b *KafkaConnectToolBuilder) handleKafkaConnectorPluginsList(ctx context.Context, admin kafka.Connect) (*sdk.CallToolResult, error) {
	plugins, err := admin.ListPlugins(ctx)
	if err != nil {
		return nil, b.handleError("get Kafka Connect connector plugins", err)
	}
	return b.marshalResponse(plugins)
}
