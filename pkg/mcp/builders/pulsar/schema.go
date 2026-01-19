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
	"path/filepath"
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

type pulsarAdminSchemaInput struct {
	Resource  string   `json:"resource"`
	Operation string   `json:"operation"`
	Topic     string   `json:"topic"`
	Version   *float64 `json:"version,omitempty"`
	Filename  *string  `json:"filename,omitempty"`
}

const (
	pulsarAdminSchemaResourceDesc = "Resource to operate on. Available resources:\n" +
		"- schema: The schema configuration for a specific topic"
	pulsarAdminSchemaOperationDesc = "Operation to perform. Available operations:\n" +
		"- get: Get the schema for a topic (optionally by version)\n" +
		"- upload: Upload a new schema for a topic (requires namespace admin permissions)\n" +
		"- delete: Delete the schema for a topic (requires namespace admin permissions)"
	pulsarAdminSchemaTopicDesc = "The fully qualified topic name in the format 'persistent://tenant/namespace/topic'. " +
		"A schema is always associated with a specific topic. The schema will be enforced for all producers " +
		"and consumers of this topic."
	pulsarAdminSchemaVersionDesc = "The schema version (optional for 'get' operation). " +
		"Pulsar maintains a versioned history of schemas. If not specified, the latest schema version will be returned. " +
		"Use this parameter to retrieve a specific historical version of the schema."
	pulsarAdminSchemaFilenameDesc = "The file path of the schema definition (required for 'upload' operation). " +
		"The file should contain a JSON object with 'type', 'schema', and optionally 'properties' fields. " +
		"Supported schema types include: AVRO, JSON, PROTOBUF, PROTOBUF_NATIVE, KEY_VALUE, BYTES, STRING, " +
		"INT8, INT16, INT32, INT64, FLOAT, DOUBLE, BOOLEAN, NONE."
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
func (b *PulsarAdminSchemaToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool, err := b.buildSchemaTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildSchemaHandler(config.ReadOnly)

	return []builders.ToolDefinition{
		builders.ServerTool[pulsarAdminSchemaInput, any]{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildSchemaTool builds the Pulsar Admin Schema MCP tool definition
// Migrated from the original tool definition logic
func (b *PulsarAdminSchemaToolBuilder) buildSchemaTool() (*sdk.Tool, error) {
	inputSchema, err := buildPulsarAdminSchemaInputSchema()
	if err != nil {
		return nil, err
	}

	toolDesc := "Manage Apache Pulsar schemas for topics. " +
		"Schemas in Pulsar define the structure of message data, enabling data validation, evolution, and interoperability. " +
		"Pulsar supports multiple schema types including AVRO, JSON, PROTOBUF, etc., allowing strong typing of message content. " +
		"Schema versioning ensures backward/forward compatibility as data structures evolve over time. " +
		"Operations include getting, uploading, and deleting schemas. " +
		"Requires namespace admin permissions for all operations."

	return &sdk.Tool{
		Name:        "pulsar_admin_schema",
		Description: toolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildSchemaHandler builds the Pulsar Admin Schema handler function
// Migrated from the original handler logic
func (b *PulsarAdminSchemaToolBuilder) buildSchemaHandler(readOnly bool) builders.ToolHandlerFunc[pulsarAdminSchemaInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input pulsarAdminSchemaInput) (*sdk.CallToolResult, any, error) {
		resource := strings.ToLower(input.Resource)
		if resource == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'resource'")
		}

		operation := strings.ToLower(input.Operation)
		if operation == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'operation'")
		}

		topic := input.Topic
		if topic == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'topic'. please provide the fully qualified topic name")
		}

		// Validate write operations in read-only mode
		if readOnly && (operation == "upload" || operation == "delete") {
			return nil, nil, fmt.Errorf("write operations are not allowed in read-only mode")
		}

		// Verify resource type
		if resource != "schema" {
			return nil, nil, fmt.Errorf("invalid resource: %s. only 'schema' is supported", resource)
		}

		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return nil, nil, fmt.Errorf("pulsar session not found in context")
		}

		// Create the admin client
		admin, err := session.GetAdminClient()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get admin client: %v", err)
		}

		// Dispatch based on operation
		switch operation {
		case "get":
			result, err := b.handleSchemaGet(admin, topic, input)
			return result, nil, err
		case "upload":
			result, err := b.handleSchemaUpload(admin, topic, input)
			return result, nil, err
		case "delete":
			result, err := b.handleSchemaDelete(admin, topic)
			return result, nil, err
		default:
			return nil, nil, fmt.Errorf("unknown operation: %s", operation)
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
func (b *PulsarAdminSchemaToolBuilder) handleSchemaGet(admin cmdutils.Client, topic string, input pulsarAdminSchemaInput) (*sdk.CallToolResult, error) {
	// Get optional version parameter
	var version float64
	if input.Version != nil {
		version = *input.Version
	}

	// Get schema info
	if version != 0 {
		// Get schema by version
		info, err := admin.Schemas().GetSchemaInfoByVersion(topic, int64(version))
		if err != nil {
			return nil, fmt.Errorf("failed to get schema version %v for topic '%s': %v", version, topic, err)
		}

		jsonBytes, err := json.Marshal(info)
		if err != nil {
			return nil, fmt.Errorf("failed to process schema information: %v", err)
		}

		return textResult(string(jsonBytes)), nil
	}
	// Get latest schema
	schemaInfoWithVersion, err := admin.Schemas().GetSchemaInfoWithVersion(topic)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest schema for topic '%s': %v", topic, err)
	}

	// Format the output
	var output bytes.Buffer
	name, err := json.Marshal(schemaInfoWithVersion.SchemaInfo.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to process schema name: %v", err)
	}

	schemaType, err := json.Marshal(schemaInfoWithVersion.SchemaInfo.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to process schema type: %v", err)
	}

	properties, err := json.Marshal(schemaInfoWithVersion.SchemaInfo.Properties)
	if err != nil {
		return nil, fmt.Errorf("failed to process schema properties: %v", err)
	}

	schema, err := b.prettyPrint(schemaInfoWithVersion.SchemaInfo.Schema)
	if err != nil {
		return nil, fmt.Errorf("failed to format schema definition: %v", err)
	}

	fmt.Fprintf(&output, "{\n  name: %s \n  schema: %s\n  type: %s \n  properties: %s\n  version: %d\n}",
		string(name), string(schema), string(schemaType), string(properties), schemaInfoWithVersion.Version)

	return textResult(output.String()), nil
}

// handleSchemaUpload handles uploading a schema
func (b *PulsarAdminSchemaToolBuilder) handleSchemaUpload(admin cmdutils.Client, topic string, input pulsarAdminSchemaInput) (*sdk.CallToolResult, error) {
	if input.Filename == nil || *input.Filename == "" {
		return nil, fmt.Errorf("missing required parameter 'filename' for schema.upload. please provide the path to the schema definition file")
	}
	filename := *input.Filename

	// Read and parse the schema file
	var payload utils.PostSchemaPayload
	file, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file '%s': %v", filename, err)
	}

	err = json.Unmarshal(file, &payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse schema file '%s'. the file must contain valid JSON with 'type', 'schema', and optionally 'properties' fields: %v",
			filename, err)
	}

	// Upload the schema
	err = admin.Schemas().CreateSchemaByPayload(topic, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to upload schema for topic '%s': %v", topic, err)
	}

	return textResult(fmt.Sprintf("Schema uploaded successfully for topic '%s'", topic)), nil
}

// handleSchemaDelete handles deleting a schema
func (b *PulsarAdminSchemaToolBuilder) handleSchemaDelete(admin cmdutils.Client, topic string) (*sdk.CallToolResult, error) {
	// Delete the schema
	err := admin.Schemas().DeleteSchema(topic)
	if err != nil {
		return nil, fmt.Errorf("failed to delete schema for topic '%s': %v", topic, err)
	}

	return textResult(fmt.Sprintf("Schema deleted successfully for topic '%s'", topic)), nil
}

func buildPulsarAdminSchemaInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[pulsarAdminSchemaInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}

	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}

	setSchemaDescription(schema, "resource", pulsarAdminSchemaResourceDesc)
	setSchemaDescription(schema, "operation", pulsarAdminSchemaOperationDesc)
	setSchemaDescription(schema, "topic", pulsarAdminSchemaTopicDesc)
	setSchemaDescription(schema, "version", pulsarAdminSchemaVersionDesc)
	setSchemaDescription(schema, "filename", pulsarAdminSchemaFilenameDesc)

	normalizeAdditionalProperties(schema)
	return schema, nil
}
