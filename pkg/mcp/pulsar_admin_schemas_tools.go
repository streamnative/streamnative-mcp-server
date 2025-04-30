package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// PulsarAdminAddSchemasTools adds schema-related tools to the MCP server
func PulsarAdminAddSchemasTools(s *server.MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeaturePulsarAdminSchemas)) && !slices.Contains(features, string(FeatureAll)) && !slices.Contains(features, string(FeatureAllPulsar)) {
		return
	}

	// Add unified schema management tool
	toolDesc := "Manage Apache Pulsar schemas for topics. " +
		"Schemas in Pulsar define the structure of message data, enabling data validation, evolution, and interoperability. " +
		"Pulsar supports multiple schema types including AVRO, JSON, PROTOBUF, etc., allowing strong typing of message content. " +
		"Schema versioning ensures backward/forward compatibility as data structures evolve over time. " +
		"Operations include getting, uploading, and deleting schemas. " +
		"Requires namespace admin permissions for all operations."

	resourceDesc := "Resource to operate on. Available resources:\n" +
		"- schema: The schema configuration for a specific topic"

	operationDesc := "Operation to perform. Available operations:\n" +
		"- get: Get the schema for a topic (optionally by version)\n" +
		"- upload: Upload a new schema for a topic (requires namespace admin permissions)\n" +
		"- delete: Delete the schema for a topic (requires namespace admin permissions)"

	schemasTool := mcp.NewTool("pulsar_admin_schema",
		mcp.WithDescription(toolDesc),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description(resourceDesc),
		),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
		),
		mcp.WithString("topic", mcp.Required(),
			mcp.Description("The fully qualified topic name in the format 'persistent://tenant/namespace/topic'. "+
				"A schema is always associated with a specific topic. The schema will be enforced for all producers "+
				"and consumers of this topic."),
		),
		mcp.WithNumber("version",
			mcp.Description("The schema version (optional for 'get' operation). "+
				"Pulsar maintains a versioned history of schemas. If not specified, the latest schema version will be returned. "+
				"Use this parameter to retrieve a specific historical version of the schema."),
		),
		mcp.WithString("filename",
			mcp.Description("The file path of the schema definition (required for 'upload' operation). "+
				"The file should contain a JSON object with 'type', 'schema', and optionally 'properties' fields. "+
				"Supported schema types include: AVRO, JSON, PROTOBUF, PROTOBUF_NATIVE, KEY_VALUE, BYTES, STRING, INT8, INT16, INT32, INT64, FLOAT, DOUBLE, BOOLEAN, NONE."),
		),
	)
	s.AddTool(schemasTool, handleSchemaTool(readOnly))
}

// handleSchemaTool returns a function that handles schema operations
func handleSchemaTool(readOnly bool) func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get required parameters
		resource, err := requiredParam[string](request.Params.Arguments, "resource")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get resource: %v", err)), nil
		}

		operation, err := requiredParam[string](request.Params.Arguments, "operation")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get operation: %v", err)), nil
		}

		topic, err := requiredParam[string](request.Params.Arguments, "topic")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic'. Please provide the fully qualified topic name: %v", err)), nil
		}

		// Normalize parameters
		resource = strings.ToLower(resource)
		operation = strings.ToLower(operation)

		// Validate write operations in read-only mode
		if readOnly && (operation == "upload" || operation == "delete") {
			return mcp.NewToolResultError("Write operations are not allowed in read-only mode"), nil
		}

		// Verify resource type
		if resource != "schema" {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resource: %s. Only 'schema' is supported", resource)), nil
		}

		// Create the admin client
		admin, err := pulsar.GetAdminClient()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
		}

		// Dispatch based on operation
		switch operation {
		case "get":
			return handleSchemaGet(admin, topic, request)
		case "upload":
			return handleSchemaUpload(admin, topic, request)
		case "delete":
			return handleSchemaDelete(admin, topic)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid operation: %s. Available operations: get, upload, delete", operation)), nil
		}
	}
}

// handleSchemaGet handles getting a schema
func handleSchemaGet(admin cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get optional version parameter
	version, hasVersion := optionalParam[float64](request.Params.Arguments, "version")

	// Get schema info
	if hasVersion {
		// Get schema by version
		info, err := admin.Schemas().GetSchemaInfoByVersion(topic, int64(version))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get schema version %v for topic '%s': %v",
				version, topic, err)), nil
		}

		jsonBytes, err := json.Marshal(info)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to process schema information: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonBytes)), nil
	}
	// Get latest schema
	schemaInfoWithVersion, err := admin.Schemas().GetSchemaInfoWithVersion(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get latest schema for topic '%s': %v",
			topic, err)), nil
	}

	// Format the output
	var output bytes.Buffer
	name, err := json.Marshal(schemaInfoWithVersion.SchemaInfo.Name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to process schema name: %v", err)), nil
	}

	schemaType, err := json.Marshal(schemaInfoWithVersion.SchemaInfo.Type)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to process schema type: %v", err)), nil
	}

	properties, err := json.Marshal(schemaInfoWithVersion.SchemaInfo.Properties)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to process schema properties: %v", err)), nil
	}

	schema, err := prettyPrint(schemaInfoWithVersion.SchemaInfo.Schema)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format schema definition: %v", err)), nil
	}

	fmt.Fprintf(&output, "{\n  name: %s \n  schema: %s\n  type: %s \n  properties: %s\n  version: %d\n}",
		string(name), string(schema), string(schemaType), string(properties), schemaInfoWithVersion.Version)

	return mcp.NewToolResultText(output.String()), nil

}

// handleSchemaUpload handles uploading a schema
func handleSchemaUpload(admin cmdutils.Client, topic string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filename, err := requiredParam[string](request.Params.Arguments, "filename")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'filename' for schema.upload. Please provide the path to the schema definition file: %v", err)), nil
	}

	// Read and parse the schema file
	var payload utils.PostSchemaPayload
	file, err := os.ReadFile(filename)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to read schema file '%s': %v", filename, err)), nil
	}

	err = json.Unmarshal(file, &payload)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to parse schema file '%s'. The file must contain valid JSON with 'type', 'schema', and optionally 'properties' fields: %v",
			filename, err)), nil
	}

	// Upload the schema
	err = admin.Schemas().CreateSchemaByPayload(topic, payload)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to upload schema for topic '%s': %v", topic, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Schema uploaded successfully for topic '%s'", topic)), nil
}

// handleSchemaDelete handles deleting a schema
func handleSchemaDelete(admin cmdutils.Client, topic string) (*mcp.CallToolResult, error) {
	// Delete the schema
	err := admin.Schemas().DeleteSchema(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete schema for topic '%s': %v", topic, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Schema deleted successfully for topic '%s'", topic)), nil
}

func prettyPrint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "    ")
	return out.Bytes(), err
}
