package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
)

func KafkaAdminAddKafkaConnectTools(s *server.MCPServer, readOnly bool) {
	resourceDesc := "Resource to operate on. Available resources:\n" +
		"- kafka-connect-cluster: A single Kafka Connect cluster that manages connectors and tasks.\n" +
		"- connector: A single Kafka Connect connector instance that moves data between Kafka and external systems.\n" +
		"- connectors: Collection of all Kafka Connect connectors in a cluster.\n" +
		"- connector-plugins: Collection of all Kafka Connect connector plugins, StreamNative Cloud provides a set of built-in connectors via this resource."

	operationDesc := "Operation to perform. Available operations:\n" +
		"- list: List all connectors or connector plugins in a cluster.\n" +
		"- get: Retrieve detailed information about a Kafka Connect cluster or specific connector.\n" +
		"- create: Create a new connector with specified configuration.\n" +
		"- update: Modify an existing connector's configuration.\n" +
		"- delete: Remove a connector from the Kafka Connect cluster.\n" +
		"- restart: Restart a running connector (useful after failures or configuration changes).\n" +
		"- pause: Temporarily stop a connector from processing data.\n" +
		"- resume: Continue processing with a previously paused connector."

	toolDesc := "Unified tool for managing Apache Kafka Connect.\n" +
		"Kafka Connect is a framework for connecting Kafka with external systems such as databases, key-value stores, search indexes, and file systems. " +
		"It provides a standardized way to stream data in and out of Kafka, without requiring custom integration code.\n\n" +
		"Key concepts in Kafka Connect:\n\n" +
		"- Connectors: The high-level abstraction that coordinates data streaming by managing tasks\n" +
		"- Tasks: The implementation of how data is copied to or from Kafka\n" +
		"- Workers: The running processes that execute connectors and tasks\n" +
		"- Plugins: Reusable connector implementations for specific external systems\n" +
		"- Source Connectors: Import data from external systems into Kafka topics\n" +
		"- Sink Connectors: Export data from Kafka topics to external systems\n\n" +
		"Kafka Connect simplifies data integration, enables scalable and reliable streaming pipelines, " +
		"and reduces the operational burden of managing data flows.\n\n" +
		"Usage Examples:\n\n" +
		"1. List all connectors in the Kafka Connect cluster:\n" +
		"   resource: \"connectors\"\n" +
		"   operation: \"list\"\n\n" +
		"2. Get information about a specific connector:\n" +
		"   resource: \"connector\"\n" +
		"   operation: \"get\"\n" +
		"   name: \"my-jdbc-source\"\n\n" +
		"3. Create a new JDBC source connector:\n" +
		"   resource: \"connector\"\n" +
		"   operation: \"create\"\n" +
		"   name: \"my-jdbc-source\"\n" +
		"   config: {\n" +
		"     \"connector.class\": \"io.confluent.connect.jdbc.JdbcSourceConnector\",\n" +
		"     \"connection.url\": \"jdbc:mysql://mysql:3306/mydb\",\n" +
		"     \"connection.user\": \"user\",\n" +
		"     \"connection.password\": \"password\",\n" +
		"     \"topic.prefix\": \"mysql-\",\n" +
		"     \"table.whitelist\": \"users,orders\",\n" +
		"     \"mode\": \"incrementing\",\n" +
		"     \"incrementing.column.name\": \"id\",\n" +
		"     \"tasks.max\": \"1\"\n" +
		"   }\n\n" +
		"4. Update an existing connector's configuration:\n" +
		"   resource: \"connector\"\n" +
		"   operation: \"update\"\n" +
		"   name: \"my-jdbc-source\"\n" +
		"   config: {\n" +
		"     \"connector.class\": \"io.confluent.connect.jdbc.JdbcSourceConnector\",\n" +
		"     \"tasks.max\": \"2\",\n" +
		"     \"table.whitelist\": \"users,orders,products\"\n" +
		"   }\n\n" +
		"5. Delete a connector:\n" +
		"   resource: \"connector\"\n" +
		"   operation: \"delete\"\n" +
		"   name: \"my-jdbc-source\"\n\n" +
		"6. Restart a connector after configuration changes or errors:\n" +
		"   resource: \"connector\"\n" +
		"   operation: \"restart\"\n" +
		"   name: \"my-jdbc-source\"\n\n" +
		"7. Pause a connector temporarily:\n" +
		"   resource: \"connector\"\n" +
		"   operation: \"pause\"\n" +
		"   name: \"my-jdbc-source\"\n\n" +
		"8. Resume a paused connector:\n" +
		"   resource: \"connector\"\n" +
		"   operation: \"resume\"\n" +
		"   name: \"my-jdbc-source\"\n\n" +
		"9. List all available connector plugins:\n" +
		"   resource: \"connector-plugins\"\n" +
		"   operation: \"list\"\n\n" +
		"10. Get information about the Kafka Connect cluster:\n" +
		"    resource: \"kafka-connect-cluster\"\n" +
		"    operation: \"get\"\n\n" +
		"This tool requires appropriate Kafka Connect permissions."

	kafkaConnectTool := mcp.NewTool("kafka_admin_connect",
		mcp.WithDescription(toolDesc),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description(resourceDesc),
		),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
		),
		mcp.WithString("name",
			mcp.Description("The name of the Kafka Connect connector to operate on. "+
				"Required for 'get', 'create', 'update', 'delete', 'restart', 'pause', and 'resume' operations on the 'connector' resource. "+
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
	)

	s.AddTool(kafkaConnectTool, handleKafkaConnectTool(readOnly))
}

func handleKafkaConnectTool(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get required parameters
		resource, err := requiredParam[string](request.Params.Arguments, "resource")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get resource: %v", err)), nil
		}

		operation, err := requiredParam[string](request.Params.Arguments, "operation")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get operation: %v", err)), nil
		}

		// Normalize parameters
		resource = strings.ToLower(resource)
		operation = strings.ToLower(operation)

		// Validate write operations in read-only mode
		if readOnly && (operation == "create" || operation == "update" || operation == "delete" || operation == "restart" || operation == "pause" || operation == "resume") {
			return mcp.NewToolResultError("Write operations are not allowed in read-only mode"), nil
		}

		// Create the admin client
		admin, err := kafka.GetKafkaConnectClient()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
		}

		// Dispatch based on resource and operation
		switch resource {
		case "kafka-connect-cluster":
			switch operation {
			case "get":
				return handleKafkaConnectClusterGet(ctx, admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'kafka-connect-cluster': %s", operation)), nil
			}
		case "connectors":
			switch operation {
			case "list":
				return handleKafkaConnectorsList(ctx, admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'connectors': %s", operation)), nil
			}
		case "connector":
			switch operation {
			case "get":
				return handleKafkaConnectorGet(ctx, admin, request)
			case "create":
				return handleKafkaConnectorCreate(ctx, admin, request)
			case "update":
				return handleKafkaConnectorUpdate(ctx, admin, request)
			case "delete":
				return handleKafkaConnectorDelete(ctx, admin, request)
			case "restart":
				return handleKafkaConnectorRestart(ctx, admin, request)
			case "pause":
				return handleKafkaConnectorPause(ctx, admin, request)
			case "resume":
				return handleKafkaConnectorResume(ctx, admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'connector': %s", operation)), nil
			}
		case "connector-plugins":
			switch operation {
			case "list":
				return handleKafkaConnectorPluginsList(ctx, admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'connector-plugins': %s", operation)), nil
			}
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resource: %s. Available resources: kafka-connect-cluster, connectors, connector, connector-plugins", resource)), nil
		}
	}
}

func handleKafkaConnectClusterGet(ctx context.Context, admin kafka.Connect, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get the Kafka Connect cluster
	cluster, err := admin.GetInfo(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get Kafka Connect cluster: %v", err)), nil
	}

	jsonBytes, err := json.Marshal(cluster)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal Kafka Connect cluster: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func handleKafkaConnectorsList(ctx context.Context, admin kafka.Connect, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get all connectors
	connectors, err := admin.ListConnectors(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get Kafka Connect connectors: %v", err)), nil
	}

	jsonBytes, err := json.Marshal(connectors)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal Kafka Connect connectors: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func handleKafkaConnectorGet(ctx context.Context, admin kafka.Connect, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get a specific connector
	name, err := requiredParam[string](request.Params.Arguments, "name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get connector name: %v", err)), nil
	}

	connector, err := admin.GetConnector(ctx, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get Kafka Connect connector: %v", err)), nil
	}

	jsonBytes, err := json.Marshal(connector)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal Kafka Connect connector: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func handleKafkaConnectorCreate(ctx context.Context, admin kafka.Connect, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create a new connector
	name, err := requiredParam[string](request.Params.Arguments, "name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get connector name: %v", err)), nil
	}

	configMap, err := requiredParamObject(request.Params.Arguments, "config")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get config: %v", err)), nil
	}

	config := convertToMapString(configMap)

	config["name"] = name

	connector, err := admin.CreateConnector(ctx, name, config)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create Kafka Connect connector: %v", err)), nil
	}

	jsonBytes, err := json.Marshal(connector)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal Kafka Connect connector: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func handleKafkaConnectorUpdate(ctx context.Context, admin kafka.Connect, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Update a connector
	name, err := requiredParam[string](request.Params.Arguments, "name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get connector name: %v", err)), nil
	}

	configMap, err := requiredParamObject(request.Params.Arguments, "config")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get config: %v", err)), nil
	}

	config := convertToMapString(configMap)

	config["name"] = name

	connector, err := admin.UpdateConnector(ctx, name, config)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to update Kafka Connect connector: %v", err)), nil
	}

	jsonBytes, err := json.Marshal(connector)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal Kafka Connect connector: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func handleKafkaConnectorDelete(ctx context.Context, admin kafka.Connect, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Delete a connector
	name, err := requiredParam[string](request.Params.Arguments, "name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get connector name: %v", err)), nil
	}

	err = admin.DeleteConnector(ctx, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete Kafka Connect connector: %v", err)), nil
	}

	return mcp.NewToolResultText("Connector deleted successfully"), nil
}

func handleKafkaConnectorRestart(ctx context.Context, admin kafka.Connect, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Restart a connector
	name, err := requiredParam[string](request.Params.Arguments, "name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get connector name: %v", err)), nil
	}

	err = admin.RestartConnector(ctx, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to restart Kafka Connect connector: %v", err)), nil
	}

	return mcp.NewToolResultText("Connector restarted successfully"), nil
}

func handleKafkaConnectorPause(ctx context.Context, admin kafka.Connect, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Pause a connector
	name, err := requiredParam[string](request.Params.Arguments, "name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get connector name: %v", err)), nil
	}

	err = admin.PauseConnector(ctx, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to pause Kafka Connect connector: %v", err)), nil
	}

	return mcp.NewToolResultText("Connector paused successfully"), nil
}

func handleKafkaConnectorResume(ctx context.Context, admin kafka.Connect, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Resume a connector
	name, err := requiredParam[string](request.Params.Arguments, "name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get connector name: %v", err)), nil
	}

	err = admin.ResumeConnector(ctx, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to resume Kafka Connect connector: %v", err)), nil
	}

	return mcp.NewToolResultText("Connector resumed successfully"), nil
}

func handleKafkaConnectorPluginsList(ctx context.Context, admin kafka.Connect, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// List all connector plugins
	plugins, err := admin.ListPlugins(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get Kafka Connect connector plugins: %v", err)), nil
	}

	jsonBytes, err := json.Marshal(plugins)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal Kafka Connect connector plugins: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}
