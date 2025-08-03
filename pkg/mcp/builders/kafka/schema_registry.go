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

package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	"github.com/twmb/franz-go/pkg/sr"
)

// KafkaSchemaRegistryToolBuilder implements the ToolBuilder interface for Kafka Schema Registry
type KafkaSchemaRegistryToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewKafkaSchemaRegistryToolBuilder creates a new Kafka Schema Registry tool builder instance
func NewKafkaSchemaRegistryToolBuilder() *KafkaSchemaRegistryToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "kafka_schema_registry",
		Version:     "1.0.0",
		Description: "Kafka Schema Registry administration tools",
		Category:    "kafka_admin",
		Tags:        []string{"kafka", "schema", "registry", "admin"},
	}

	features := []string{
		"kafka-admin-schema-registry",
		"all",
		"all-kafka",
		"kafka-admin",
	}

	return &KafkaSchemaRegistryToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Kafka Schema Registry tool list
func (b *KafkaSchemaRegistryToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool := b.buildKafkaSchemaRegistryTool()
	handler := b.buildKafkaSchemaRegistryHandler(config.ReadOnly)

	return []server.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildKafkaSchemaRegistryTool builds the Kafka Schema Registry MCP tool definition
func (b *KafkaSchemaRegistryToolBuilder) buildKafkaSchemaRegistryTool() mcp.Tool {
	resourceDesc := "Resource to operate on. Available resources:\n" +
		"- subjects: Collection of all schema subjects in the Schema Registry\n" +
		"- subject: A specific schema subject (a named schema that can have multiple versions)\n" +
		"- versions: Collection of all versions for a specific subject\n" +
		"- version: A specific version of a subject's schema\n" +
		"- compatibility: Compatibility settings that control schema evolution rules\n" +
		"- types: Supported schema format types (like AVRO, JSON, PROTOBUF)"

	operationDesc := "Operation to perform. Available operations:\n" +
		"- list: List all subjects, versions for a subject, or supported schema types\n" +
		"- get: Get a subject's latest schema, a specific version, or compatibility setting\n" +
		"- set: Set compatibility level for global or subject-specific schema evolution\n" +
		"- create: Register a new schema for a subject\n" +
		"- delete: Delete a schema subject or a specific version"

	toolDesc := "Unified tool for managing Apache Kafka Schema Registry.\n" +
		"Schema Registry provides a centralized repository for managing and validating schemas for Kafka data.\n" +
		"It enables schema evolution while maintaining compatibility between producers and consumers.\n\n" +
		"Key concepts:\n" +
		"- Subject: A named schema that can have multiple versions\n" +
		"- Version: A specific instance of a schema for a subject\n" +
		"- Compatibility: Rules that govern how schemas can evolve\n" +
		"- Schema Types: Format types like AVRO, JSON Schema, and Protocol Buffers\n\n" +
		"Compatibility levels:\n" +
		"- BACKWARD: New schema can read data written with previous schema\n" +
		"- FORWARD: Previous schema can read data written with new schema\n" +
		"- FULL: Both backward and forward compatibility\n" +
		"- NONE: No compatibility enforcement\n\n" +
		"Usage Examples:\n\n" +
		"1. List all schema subjects:\n" +
		"   resource: \"subjects\"\n" +
		"   operation: \"list\"\n\n" +
		"2. Get the latest schema for a subject:\n" +
		"   resource: \"subject\"\n" +
		"   operation: \"get\"\n" +
		"   subject: \"user-events-value\"\n\n" +
		"3. Register a new schema:\n" +
		"   resource: \"subject\"\n" +
		"   operation: \"create\"\n" +
		"   subject: \"user-events-value\"\n" +
		"   schema: {...}\n" +
		"   schemaType: \"AVRO\"\n\n" +
		"4. Set compatibility level:\n" +
		"   resource: \"compatibility\"\n" +
		"   operation: \"set\"\n" +
		"   subject: \"user-events-value\"\n" +
		"   compatibility: \"BACKWARD\"\n\n" +
		"This tool requires appropriate Schema Registry permissions."

	return mcp.NewTool("kafka_admin_sr",
		mcp.WithDescription(toolDesc),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description(resourceDesc),
		),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
		),
		mcp.WithString("subject",
			mcp.Description("The name of the schema subject. "+
				"Required for operations on 'subject', 'versions', 'version', and subject-specific 'compatibility' resources. "+
				"Subject names typically follow the pattern '<topic-name>-key' or '<topic-name>-value'.")),
		mcp.WithString("version",
			mcp.Description("The version number or 'latest' for the most recent version. "+
				"Required for 'version' resource operations.")),
		mcp.WithString("compatibility",
			mcp.Description("The compatibility level to set. "+
				"Valid values: BACKWARD, FORWARD, FULL, NONE. "+
				"Required for 'set' operation on 'compatibility' resource.")),
		mcp.WithString("schemaType",
			mcp.Description("The schema format type. "+
				"Valid values: AVRO, JSON, PROTOBUF. "+
				"Required for 'create' operation on 'subject' resource.")),
		mcp.WithObject("schema",
			mcp.Description("The schema definition as a JSON object. "+
				"Required for 'create' operation on 'subject' resource. "+
				"The structure depends on the schema type (AVRO, JSON Schema, or Protocol Buffers).")),
	)
}

// buildKafkaSchemaRegistryHandler builds the Kafka Schema Registry handler function
func (b *KafkaSchemaRegistryToolBuilder) buildKafkaSchemaRegistryHandler(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		// Validate write operations in read-only mode
		if readOnly && (operation == "create" || operation == "delete" || operation == "set") {
			return mcp.NewToolResultError("Write operations are not allowed in read-only mode"), nil
		}

		// Get Schema Registry client
		session := mcpCtx.GetKafkaSession(ctx)
		if session == nil {
			return b.handleError("get Kafka session not found in context", nil), nil
		}
		client, err := session.GetSchemaRegistryClient()
		if err != nil {
			return b.handleError("get Schema Registry client", err), nil
		}

		// Dispatch based on resource and operation
		switch resource {
		case "subjects":
			switch operation {
			case "list":
				return b.handleSchemaSubjectsList(ctx, client, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'subjects': %s", operation)), nil
			}
		case "subject":
			switch operation {
			case "get":
				return b.handleSchemaSubjectGet(ctx, client, request)
			case "create":
				return b.handleSchemaSubjectCreate(ctx, client, request)
			case "delete":
				return b.handleSchemaSubjectDelete(ctx, client, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'subject': %s", operation)), nil
			}
		case "versions":
			switch operation {
			case "list":
				return b.handleSchemaVersionsList(ctx, client, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'versions': %s", operation)), nil
			}
		case "version":
			switch operation {
			case "get":
				return b.handleSchemaVersionGet(ctx, client, request)
			case "delete":
				return b.handleSchemaVersionDelete(ctx, client, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'version': %s", operation)), nil
			}
		case "compatibility":
			switch operation {
			case "get":
				return b.handleSchemaCompatibilityGet(ctx, client, request)
			case "set":
				return b.handleSchemaCompatibilitySet(ctx, client, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'compatibility': %s", operation)), nil
			}
		case "types":
			switch operation {
			case "list":
				return b.handleSchemaTypesList(ctx, client, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'types': %s", operation)), nil
			}
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resource: %s. Available resources: subjects, subject, versions, version, compatibility, types", resource)), nil
		}
	}
}

// Utility functions

// handleError provides unified error handling
func (b *KafkaSchemaRegistryToolBuilder) handleError(operation string, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("Failed to %s: %v", operation, err))
}

// marshalResponse provides unified JSON serialization for responses
func (b *KafkaSchemaRegistryToolBuilder) marshalResponse(data interface{}) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return b.handleError("marshal response", err), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// Specific operation handler functions

// handleSchemaSubjectsList handles listing all schema subjects
func (b *KafkaSchemaRegistryToolBuilder) handleSchemaSubjectsList(ctx context.Context, client *sr.Client, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subjects, err := client.Subjects(ctx)
	if err != nil {
		return b.handleError("list schema subjects", err), nil
	}
	return b.marshalResponse(subjects)
}

// handleSchemaSubjectGet handles getting the latest schema for a subject
func (b *KafkaSchemaRegistryToolBuilder) handleSchemaSubjectGet(ctx context.Context, client *sr.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subject, err := request.RequireString("subject")
	if err != nil {
		return b.handleError("get subject name", err), nil
	}

	schema, err := client.SchemaByVersion(ctx, subject, -1) // -1 for latest
	if err != nil {
		return b.handleError("get schema for subject", err), nil
	}
	return b.marshalResponse(schema)
}

// handleSchemaSubjectCreate handles registering a new schema for a subject
func (b *KafkaSchemaRegistryToolBuilder) handleSchemaSubjectCreate(ctx context.Context, client *sr.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subject, err := request.RequireString("subject")
	if err != nil {
		return b.handleError("get subject name", err), nil
	}

	schemaTypeStr, err := request.RequireString("schemaType")
	if err != nil {
		return b.handleError("get schema type", err), nil
	}

	schema, err := request.RequireString("schema")
	if err != nil {
		return b.handleError("get schema object", err), nil
	}

	// Parse schema type
	var schemaType sr.SchemaType
	err = schemaType.UnmarshalText([]byte(schemaTypeStr))
	if err != nil {
		return b.handleError("unmarshal schema type", err), nil
	}

	// Create schema
	schemaObj := sr.Schema{
		Type:   schemaType,
		Schema: schema,
	}

	result, err := client.CreateSchema(ctx, subject, schemaObj)
	if err != nil {
		return b.handleError("create schema", err), nil
	}
	return b.marshalResponse(result)
}

// handleSchemaSubjectDelete handles deleting a schema subject
func (b *KafkaSchemaRegistryToolBuilder) handleSchemaSubjectDelete(ctx context.Context, client *sr.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subject, err := request.RequireString("subject")
	if err != nil {
		return b.handleError("get subject name", err), nil
	}

	// Delete subject using correct API signature (soft delete by default)
	versions, err := client.DeleteSubject(ctx, subject, sr.SoftDelete)
	if err != nil {
		return b.handleError("delete schema subject", err), nil
	}
	return b.marshalResponse(versions)
}

// handleSchemaVersionsList handles listing all versions for a subject
func (b *KafkaSchemaRegistryToolBuilder) handleSchemaVersionsList(ctx context.Context, client *sr.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subject, err := request.RequireString("subject")
	if err != nil {
		return b.handleError("get subject name", err), nil
	}

	versions, err := client.SubjectVersions(ctx, subject)
	if err != nil {
		return b.handleError("list schema versions", err), nil
	}
	return b.marshalResponse(versions)
}

// handleSchemaVersionGet handles getting a specific version of a schema
func (b *KafkaSchemaRegistryToolBuilder) handleSchemaVersionGet(ctx context.Context, client *sr.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subject, err := request.RequireString("subject")
	if err != nil {
		return b.handleError("get subject name", err), nil
	}

	versionStr, err := request.RequireString("version")
	if err != nil {
		return b.handleError("get version", err), nil
	}

	var version int
	if versionStr == "latest" {
		version = -1
	} else {
		var parseErr error
		version, parseErr = strconv.Atoi(versionStr)
		if parseErr != nil {
			return b.handleError("parse version number", parseErr), nil
		}
	}

	schema, err := client.SchemaByVersion(ctx, subject, version)
	if err != nil {
		return b.handleError("get schema version", err), nil
	}
	return b.marshalResponse(schema)
}

// handleSchemaVersionDelete handles deleting a specific version of a schema
func (b *KafkaSchemaRegistryToolBuilder) handleSchemaVersionDelete(ctx context.Context, client *sr.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subject, err := request.RequireString("subject")
	if err != nil {
		return b.handleError("get subject name", err), nil
	}

	versionStr, err := request.RequireString("version")
	if err != nil {
		return b.handleError("get version", err), nil
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return b.handleError("parse version number", err), nil
	}

	// Delete schema version using correct API signature (soft delete by default)
	err = client.DeleteSchema(ctx, subject, version, sr.SoftDelete)
	if err != nil {
		return b.handleError("delete schema version", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Schema version %d for subject %s deleted successfully", version, subject)), nil
}

// handleSchemaCompatibilityGet handles getting compatibility setting
func (b *KafkaSchemaRegistryToolBuilder) handleSchemaCompatibilityGet(ctx context.Context, client *sr.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subject := request.GetString("subject", "") // Optional for global compatibility

	var results []sr.CompatibilityResult
	if subject != "" {
		// Get compatibility for specific subject
		results = client.Compatibility(ctx, subject)
	} else {
		// Get global compatibility
		results = client.Compatibility(ctx)
	}

	// Check for errors in results
	for _, result := range results {
		if result.Err != nil {
			return b.handleError("get compatibility setting", result.Err), nil
		}
	}

	// Return the first result (there should only be one)
	if len(results) > 0 {
		return b.marshalResponse(map[string]string{"compatibility": results[0].Level.String()})
	}

	return mcp.NewToolResultError("No compatibility result returned"), nil
}

// handleSchemaCompatibilitySet handles setting compatibility level
func (b *KafkaSchemaRegistryToolBuilder) handleSchemaCompatibilitySet(ctx context.Context, client *sr.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	compatibilityStr, err := request.RequireString("compatibility")
	if err != nil {
		return b.handleError("get compatibility level", err), nil
	}

	subject := request.GetString("subject", "") // Optional for global compatibility

	// Parse compatibility level
	var compatibility sr.CompatibilityLevel
	switch strings.ToUpper(compatibilityStr) {
	case "BACKWARD":
		compatibility = sr.CompatBackward
	case "FORWARD":
		compatibility = sr.CompatForward
	case "FULL":
		compatibility = sr.CompatFull
	case "NONE":
		compatibility = sr.CompatNone
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Invalid compatibility level: %s. Valid levels: BACKWARD, FORWARD, FULL, NONE", compatibilityStr)), nil
	}

	// Create SetCompatibility request
	setCompat := sr.SetCompatibility{
		Level: compatibility,
	}

	var results []sr.CompatibilityResult
	if subject != "" {
		// Set compatibility for specific subject
		results = client.SetCompatibility(ctx, setCompat, subject)
	} else {
		// Set global compatibility
		results = client.SetCompatibility(ctx, setCompat)
	}

	// Check for errors in results
	for _, result := range results {
		if result.Err != nil {
			return b.handleError("set compatibility level", result.Err), nil
		}
	}

	return mcp.NewToolResultText(fmt.Sprintf("Compatibility level set to %s", compatibilityStr)), nil
}

// handleSchemaTypesList handles listing supported schema types
func (b *KafkaSchemaRegistryToolBuilder) handleSchemaTypesList(ctx context.Context, client *sr.Client, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	types := []string{"AVRO", "JSON", "PROTOBUF"}
	return b.marshalResponse(types)
}
