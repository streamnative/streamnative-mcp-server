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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/adapter"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

// PulsarAdminSchemaToolBuilder implements the ToolBuilder interface for Pulsar Admin Schema tools
// It provides functionality to build Pulsar schema management tools
// /nolint:revive
type PulsarAdminSchemaToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminSchemaToolBuilder creates a new Pulsar Admin Schema tool builder instance
func NewPulsarAdminSchemaToolBuilder() *PulsarAdminSchemaToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_schema",
		Version:     "1.0.0",
		Description: "Pulsar Admin schema management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "schema", "admin"},
	}

	features := []string{
		"pulsar-admin-schemas",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminSchemaToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar Admin Schema tool list
// This is the core method implementing the ToolBuilder interface
func (b *PulsarAdminSchemaToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ServerTool, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool := b.buildSchemaTool()
	handler := b.buildSchemaHandler(config.ReadOnly)

	return []builders.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildSchemaTool builds the Pulsar Admin Schema MCP tool definition
// Migrated from the original tool definition logic
func (b *PulsarAdminSchemaToolBuilder) buildSchemaTool() *mcpsdk.Tool {
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

	return builders.NewTool("pulsar_admin_schema",
		builders.WithDescription(toolDesc),
		builders.WithString("resource", builders.Required(),
			builders.Description(resourceDesc),
		),
		builders.WithString("operation", builders.Required(),
			builders.Description(operationDesc),
		),
		builders.WithString("topic", builders.Required(),
			builders.Description("The fully qualified topic name in the format 'persistent://tenant/namespace/topic'. "+
				"A schema is always associated with a specific topic. The schema will be enforced for all producers "+
				"and consumers of this topic."),
		),
		builders.WithNumber("version",
			builders.Description("The schema version (optional for 'get' operation). "+
				"Pulsar maintains a versioned history of schemas. If not specified, the latest schema version will be returned. "+
				"Use this parameter to retrieve a specific historical version of the schema."),
		),
		builders.WithString("filename",
			builders.Description("The file path of the schema definition (required for 'upload' operation). "+
				"The file should contain a JSON object with 'type', 'schema', and optionally 'properties' fields. "+
				"Supported schema types include: AVRO, JSON, PROTOBUF, PROTOBUF_NATIVE, KEY_VALUE, BYTES, STRING, INT8, INT16, INT32, INT64, FLOAT, DOUBLE, BOOLEAN, NONE."),
		),
	)
}

// buildSchemaHandler builds the Pulsar Admin Schema handler function
// Migrated from the original handler logic
func (b *PulsarAdminSchemaToolBuilder) buildSchemaHandler(readOnly bool) func(context.Context, *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	return func(ctx context.Context, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
		// Get required parameters
		resource, err := adapter.RequireString(request, "resource")
		if err != nil {
			return adapter.NewErrorResult("Failed to get resource: %v", err), nil
		}

		operation, err := adapter.RequireString(request, "operation")
		if err != nil {
			return adapter.NewErrorResult("Failed to get operation: %v", err), nil
		}

		topic, err := adapter.RequireString(request, "topic")
		if err != nil {
			return adapter.NewErrorResult("Missing required parameter 'topic'. Please provide the fully qualified topic name: %v", err), nil
		}

		// Normalize parameters
		resource = strings.ToLower(resource)
		operation = strings.ToLower(operation)

		// Validate write operations in read-only mode
		if readOnly && (operation == "upload" || operation == "delete") {
			return adapter.NewErrorResult("Write operations are not allowed in read-only mode"), nil
		}

		// Verify resource type
		if resource != "schema" {
			return adapter.NewErrorResult("Invalid resource: %s. Only 'schema' is supported", resource), nil
		}

		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return adapter.NewErrorResult("Pulsar session not found in context"), nil
		}

		// Create the admin client
		admin, err := session.GetAdminClient()
		if err != nil {
			return adapter.NewErrorResult("Failed to get admin client: %v", err), nil
		}

		// Dispatch based on operation
		switch operation {
		case "get":
			return b.handleSchemaGet(admin, topic, request)
		case "upload":
			return b.handleSchemaUpload(admin, topic, request)
		case "delete":
			return b.handleSchemaDelete(admin, topic)
		default:
			return adapter.NewErrorResult("Unknown operation: %s", operation), nil
		}
	}
}

// Unified error handling and utility functions

// prettyPrint formats JSON bytes with indentation
func (b *PulsarAdminSchemaToolBuilder) prettyPrint(data []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, data, "", "    ")
	return out.Bytes(), err
}

// Operation handler functions - migrated from the original implementation

// handleSchemaGet handles getting a schema
func (b *PulsarAdminSchemaToolBuilder) handleSchemaGet(admin cmdutils.Client, topic string, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	// Get optional version parameter
	version := adapter.GetFloat(request, "version", 0)

	// Get schema info
	if version != 0 {
		// Get schema by version
		info, err := admin.Schemas().GetSchemaInfoByVersion(topic, int64(version))
		if err != nil {
			return adapter.NewErrorResult(fmt.Sprintf("Failed to get schema version %v for topic '%s': %v",
				version, topic, err)), nil
		}

		jsonBytes, err := json.Marshal(info)
		if err != nil {
			return adapter.NewErrorResult("Failed to process schema information: %v", err), nil
		}

		return adapter.NewTextResult(string(jsonBytes)), nil
	}
	// Get latest schema
	schemaInfoWithVersion, err := admin.Schemas().GetSchemaInfoWithVersion(topic)
	if err != nil {
		return adapter.NewErrorResult(fmt.Sprintf("Failed to get latest schema for topic '%s': %v",
			topic, err)), nil
	}

	// Format the output
	var output bytes.Buffer
	name, err := json.Marshal(schemaInfoWithVersion.SchemaInfo.Name)
	if err != nil {
		return adapter.NewErrorResult("Failed to process schema name: %v", err), nil
	}

	schemaType, err := json.Marshal(schemaInfoWithVersion.SchemaInfo.Type)
	if err != nil {
		return adapter.NewErrorResult("Failed to process schema type: %v", err), nil
	}

	properties, err := json.Marshal(schemaInfoWithVersion.SchemaInfo.Properties)
	if err != nil {
		return adapter.NewErrorResult("Failed to process schema properties: %v", err), nil
	}

	schema, err := b.prettyPrint(schemaInfoWithVersion.SchemaInfo.Schema)
	if err != nil {
		return adapter.NewErrorResult("Failed to format schema definition: %v", err), nil
	}

	fmt.Fprintf(&output, "{\n  name: %s \n  schema: %s\n  type: %s \n  properties: %s\n  version: %d\n}",
		string(name), string(schema), string(schemaType), string(properties), schemaInfoWithVersion.Version)

	return adapter.NewTextResult(output.String()), nil
}

// handleSchemaUpload handles uploading a schema
func (b *PulsarAdminSchemaToolBuilder) handleSchemaUpload(admin cmdutils.Client, topic string, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	filename, err := adapter.RequireString(request, "filename")
	if err != nil {
		return adapter.NewErrorResult("Missing required parameter 'filename' for schema.upload. Please provide the path to the schema definition file: %v", err), nil
	}

	// Read and parse the schema file
	var payload utils.PostSchemaPayload
	file, err := os.ReadFile(filename)
	if err != nil {
		return adapter.NewErrorResult("Failed to read schema file '%s': %v", filename, err), nil
	}

	err = json.Unmarshal(file, &payload)
	if err != nil {
		return adapter.NewErrorResult(fmt.Sprintf("Failed to parse schema file '%s'. The file must contain valid JSON with 'type', 'schema', and optionally 'properties' fields: %v",
			filename, err)), nil
	}

	// Upload the schema
	err = admin.Schemas().CreateSchemaByPayload(topic, payload)
	if err != nil {
		return adapter.NewErrorResult("Failed to upload schema for topic '%s': %v", topic, err), nil
	}

	return adapter.NewTextResult(fmt.Sprintf("Schema uploaded successfully for topic '%s'", topic)), nil
}

// handleSchemaDelete handles deleting a schema
func (b *PulsarAdminSchemaToolBuilder) handleSchemaDelete(admin cmdutils.Client, topic string) (*mcpsdk.CallToolResult, error) {
	// Delete the schema
	err := admin.Schemas().DeleteSchema(topic)
	if err != nil {
		return adapter.NewErrorResult("Failed to delete schema for topic '%s': %v", topic, err), nil
	}

	return adapter.NewTextResult(fmt.Sprintf("Schema deleted successfully for topic '%s'", topic)), nil
}
