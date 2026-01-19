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

package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	"github.com/twmb/franz-go/pkg/sr"
)

type kafkaSchemaRegistryInput struct {
	Resource      string  `json:"resource"`
	Operation     string  `json:"operation"`
	Subject       *string `json:"subject,omitempty"`
	Version       *string `json:"version,omitempty"`
	Compatibility *string `json:"compatibility,omitempty"`
	SchemaType    *string `json:"schemaType,omitempty"`
	Schema        *string `json:"schema,omitempty"`
}

const (
	kafkaSchemaRegistryResourceDesc = "Resource to operate on. Available resources:\n" +
		"- subjects: Collection of all schema subjects in the Schema Registry\n" +
		"- subject: A specific schema subject (a named schema that can have multiple versions)\n" +
		"- versions: Collection of all versions for a specific subject\n" +
		"- version: A specific version of a subject's schema\n" +
		"- compatibility: Compatibility settings that control schema evolution rules\n" +
		"- types: Supported schema format types (like AVRO, JSON, PROTOBUF)"
	kafkaSchemaRegistryOperationDesc = "Operation to perform. Available operations:\n" +
		"- list: List all subjects, versions for a subject, or supported schema types\n" +
		"- get: Get a subject's latest schema, a specific version, or compatibility setting\n" +
		"- set: Set compatibility level for global or subject-specific schema evolution\n" +
		"- create: Register a new schema for a subject\n" +
		"- delete: Delete a schema subject or a specific version"
	kafkaSchemaRegistrySubjectDesc = "The name of the schema subject. " +
		"Required for operations on 'subject', 'versions', 'version', and subject-specific 'compatibility' resources. " +
		"Subject names typically follow the pattern '<topic-name>-key' or '<topic-name>-value'."
	kafkaSchemaRegistryVersionDesc = "The version number or 'latest' for the most recent version. " +
		"Required for 'version' resource operations."
	kafkaSchemaRegistryCompatibilityDesc = "The compatibility level to set. " +
		"Valid values: BACKWARD, FORWARD, FULL, NONE. " +
		"Required for 'set' operation on 'compatibility' resource."
	kafkaSchemaRegistrySchemaTypeDesc = "The schema format type. " +
		"Valid values: AVRO, JSON, PROTOBUF. " +
		"Required for 'create' operation on 'subject' resource."
	kafkaSchemaRegistrySchemaDesc = "The schema definition as a JSON string. " +
		"Required for 'create' operation on 'subject' resource. " +
		"The structure depends on the schema type (AVRO, JSON Schema, or Protocol Buffers)."
)

// KafkaSchemaRegistryToolBuilder implements the ToolBuilder interface for Kafka Schema Registry
// /nolint:revive
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
func (b *KafkaSchemaRegistryToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool, err := b.buildKafkaSchemaRegistryTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildKafkaSchemaRegistryHandler(config.ReadOnly)

	return []builders.ToolDefinition{
		builders.ServerTool[kafkaSchemaRegistryInput, any]{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildKafkaSchemaRegistryTool builds the Kafka Schema Registry MCP tool definition
func (b *KafkaSchemaRegistryToolBuilder) buildKafkaSchemaRegistryTool() (*sdk.Tool, error) {
	inputSchema, err := buildKafkaSchemaRegistryInputSchema()
	if err != nil {
		return nil, err
	}

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
		"   schema: \"{...}\"\n" +
		"   schemaType: \"AVRO\"\n\n" +
		"4. Set compatibility level:\n" +
		"   resource: \"compatibility\"\n" +
		"   operation: \"set\"\n" +
		"   subject: \"user-events-value\"\n" +
		"   compatibility: \"BACKWARD\"\n\n" +
		"This tool requires appropriate Schema Registry permissions."

	return &sdk.Tool{
		Name:        "kafka_admin_sr",
		Description: toolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildKafkaSchemaRegistryHandler builds the Kafka Schema Registry handler function
func (b *KafkaSchemaRegistryToolBuilder) buildKafkaSchemaRegistryHandler(readOnly bool) builders.ToolHandlerFunc[kafkaSchemaRegistryInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input kafkaSchemaRegistryInput) (*sdk.CallToolResult, any, error) {
		// Normalize parameters
		resource := strings.ToLower(input.Resource)
		operation := strings.ToLower(input.Operation)

		// Validate write operations in read-only mode
		if readOnly && (operation == "create" || operation == "delete" || operation == "set") {
			return nil, nil, fmt.Errorf("write operations are not allowed in read-only mode")
		}

		// Get Schema Registry client
		session := mcpCtx.GetKafkaSession(ctx)
		if session == nil {
			return nil, nil, b.handleError("get Kafka session not found in context", nil)
		}
		client, err := session.GetSchemaRegistryClient()
		if err != nil {
			return nil, nil, b.handleError("get Schema Registry client", err)
		}

		// Dispatch based on resource and operation
		switch resource {
		case "subjects":
			switch operation {
			case "list":
				result, err := b.handleSchemaSubjectsList(ctx, client)
				return result, nil, err
			default:
				return nil, nil, fmt.Errorf("invalid operation for resource 'subjects': %s", operation)
			}
		case "subject":
			switch operation {
			case "get":
				result, err := b.handleSchemaSubjectGet(ctx, client, input)
				return result, nil, err
			case "create":
				result, err := b.handleSchemaSubjectCreate(ctx, client, input)
				return result, nil, err
			case "delete":
				result, err := b.handleSchemaSubjectDelete(ctx, client, input)
				return result, nil, err
			default:
				return nil, nil, fmt.Errorf("invalid operation for resource 'subject': %s", operation)
			}
		case "versions":
			switch operation {
			case "list":
				result, err := b.handleSchemaVersionsList(ctx, client, input)
				return result, nil, err
			default:
				return nil, nil, fmt.Errorf("invalid operation for resource 'versions': %s", operation)
			}
		case "version":
			switch operation {
			case "get":
				result, err := b.handleSchemaVersionGet(ctx, client, input)
				return result, nil, err
			case "delete":
				result, err := b.handleSchemaVersionDelete(ctx, client, input)
				return result, nil, err
			default:
				return nil, nil, fmt.Errorf("invalid operation for resource 'version': %s", operation)
			}
		case "compatibility":
			switch operation {
			case "get":
				result, err := b.handleSchemaCompatibilityGet(ctx, client, input)
				return result, nil, err
			case "set":
				result, err := b.handleSchemaCompatibilitySet(ctx, client, input)
				return result, nil, err
			default:
				return nil, nil, fmt.Errorf("invalid operation for resource 'compatibility': %s", operation)
			}
		case "types":
			switch operation {
			case "list":
				result, err := b.handleSchemaTypesList()
				return result, nil, err
			default:
				return nil, nil, fmt.Errorf("invalid operation for resource 'types': %s", operation)
			}
		default:
			return nil, nil, fmt.Errorf("invalid resource: %s. available resources: subjects, subject, versions, version, compatibility, types", resource)
		}
	}
}

// Utility functions

// handleError provides unified error handling
func (b *KafkaSchemaRegistryToolBuilder) handleError(operation string, err error) error {
	return fmt.Errorf("failed to %s: %v", operation, err)
}

// marshalResponse provides unified JSON serialization for responses
func (b *KafkaSchemaRegistryToolBuilder) marshalResponse(data interface{}) (*sdk.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, b.handleError("marshal response", err)
	}
	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: string(jsonBytes)}},
	}, nil
}

// Specific operation handler functions

// handleSchemaSubjectsList handles listing all schema subjects
func (b *KafkaSchemaRegistryToolBuilder) handleSchemaSubjectsList(ctx context.Context, client *sr.Client) (*sdk.CallToolResult, error) {
	subjects, err := client.Subjects(ctx)
	if err != nil {
		return nil, b.handleError("list schema subjects", err)
	}
	return b.marshalResponse(subjects)
}

// handleSchemaSubjectGet handles getting the latest schema for a subject
func (b *KafkaSchemaRegistryToolBuilder) handleSchemaSubjectGet(ctx context.Context, client *sr.Client, input kafkaSchemaRegistryInput) (*sdk.CallToolResult, error) {
	subject, err := requireString(input.Subject, "subject")
	if err != nil {
		return nil, b.handleError("get subject name", err)
	}

	schema, err := client.SchemaByVersion(ctx, subject, -1) // -1 for latest
	if err != nil {
		return nil, b.handleError("get schema for subject", err)
	}
	return b.marshalResponse(schema)
}

// handleSchemaSubjectCreate handles registering a new schema for a subject
func (b *KafkaSchemaRegistryToolBuilder) handleSchemaSubjectCreate(ctx context.Context, client *sr.Client, input kafkaSchemaRegistryInput) (*sdk.CallToolResult, error) {
	subject, err := requireString(input.Subject, "subject")
	if err != nil {
		return nil, b.handleError("get subject name", err)
	}

	schemaTypeStr, err := requireString(input.SchemaType, "schemaType")
	if err != nil {
		return nil, b.handleError("get schema type", err)
	}

	schema, err := requireString(input.Schema, "schema")
	if err != nil {
		return nil, b.handleError("get schema object", err)
	}

	// Parse schema type
	var schemaType sr.SchemaType
	err = schemaType.UnmarshalText([]byte(schemaTypeStr))
	if err != nil {
		return nil, b.handleError("unmarshal schema type", err)
	}

	// Create schema
	schemaObj := sr.Schema{
		Type:   schemaType,
		Schema: schema,
	}

	result, err := client.CreateSchema(ctx, subject, schemaObj)
	if err != nil {
		return nil, b.handleError("create schema", err)
	}
	return b.marshalResponse(result)
}

// handleSchemaSubjectDelete handles deleting a schema subject
func (b *KafkaSchemaRegistryToolBuilder) handleSchemaSubjectDelete(ctx context.Context, client *sr.Client, input kafkaSchemaRegistryInput) (*sdk.CallToolResult, error) {
	subject, err := requireString(input.Subject, "subject")
	if err != nil {
		return nil, b.handleError("get subject name", err)
	}

	// Delete subject using correct API signature (soft delete by default)
	versions, err := client.DeleteSubject(ctx, subject, sr.SoftDelete)
	if err != nil {
		return nil, b.handleError("delete schema subject", err)
	}
	return b.marshalResponse(versions)
}

// handleSchemaVersionsList handles listing all versions for a subject
func (b *KafkaSchemaRegistryToolBuilder) handleSchemaVersionsList(ctx context.Context, client *sr.Client, input kafkaSchemaRegistryInput) (*sdk.CallToolResult, error) {
	subject, err := requireString(input.Subject, "subject")
	if err != nil {
		return nil, b.handleError("get subject name", err)
	}

	versions, err := client.SubjectVersions(ctx, subject)
	if err != nil {
		return nil, b.handleError("list schema versions", err)
	}
	return b.marshalResponse(versions)
}

// handleSchemaVersionGet handles getting a specific version of a schema
func (b *KafkaSchemaRegistryToolBuilder) handleSchemaVersionGet(ctx context.Context, client *sr.Client, input kafkaSchemaRegistryInput) (*sdk.CallToolResult, error) {
	subject, err := requireString(input.Subject, "subject")
	if err != nil {
		return nil, b.handleError("get subject name", err)
	}

	versionStr, err := requireString(input.Version, "version")
	if err != nil {
		return nil, b.handleError("get version", err)
	}

	var version int
	if versionStr == "latest" {
		version = -1
	} else {
		var parseErr error
		version, parseErr = strconv.Atoi(versionStr)
		if parseErr != nil {
			return nil, b.handleError("parse version number", parseErr)
		}
	}

	schema, err := client.SchemaByVersion(ctx, subject, version)
	if err != nil {
		return nil, b.handleError("get schema version", err)
	}
	return b.marshalResponse(schema)
}

// handleSchemaVersionDelete handles deleting a specific version of a schema
func (b *KafkaSchemaRegistryToolBuilder) handleSchemaVersionDelete(ctx context.Context, client *sr.Client, input kafkaSchemaRegistryInput) (*sdk.CallToolResult, error) {
	subject, err := requireString(input.Subject, "subject")
	if err != nil {
		return nil, b.handleError("get subject name", err)
	}

	versionStr, err := requireString(input.Version, "version")
	if err != nil {
		return nil, b.handleError("get version", err)
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return nil, b.handleError("parse version number", err)
	}

	// Delete schema version using correct API signature (soft delete by default)
	err = client.DeleteSchema(ctx, subject, version, sr.SoftDelete)
	if err != nil {
		return nil, b.handleError("delete schema version", err)
	}

	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: fmt.Sprintf("Schema version %d for subject %s deleted successfully", version, subject)}},
	}, nil
}

// handleSchemaCompatibilityGet handles getting compatibility setting
func (b *KafkaSchemaRegistryToolBuilder) handleSchemaCompatibilityGet(ctx context.Context, client *sr.Client, input kafkaSchemaRegistryInput) (*sdk.CallToolResult, error) {
	subject := ""
	if input.Subject != nil {
		subject = *input.Subject
	}

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
			return nil, b.handleError("get compatibility setting", result.Err)
		}
	}

	// Return the first result (there should only be one)
	if len(results) > 0 {
		return b.marshalResponse(map[string]string{"compatibility": results[0].Level.String()})
	}

	return nil, fmt.Errorf("no compatibility result returned")
}

// handleSchemaCompatibilitySet handles setting compatibility level
func (b *KafkaSchemaRegistryToolBuilder) handleSchemaCompatibilitySet(ctx context.Context, client *sr.Client, input kafkaSchemaRegistryInput) (*sdk.CallToolResult, error) {
	compatibilityStr, err := requireString(input.Compatibility, "compatibility")
	if err != nil {
		return nil, b.handleError("get compatibility level", err)
	}

	subject := ""
	if input.Subject != nil {
		subject = *input.Subject
	}

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
		return nil, fmt.Errorf("invalid compatibility level: %s. valid levels: BACKWARD, FORWARD, FULL, NONE", compatibilityStr)
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
			return nil, b.handleError("set compatibility level", result.Err)
		}
	}

	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: fmt.Sprintf("Compatibility level set to %s", compatibilityStr)}},
	}, nil
}

// handleSchemaTypesList handles listing supported schema types
func (b *KafkaSchemaRegistryToolBuilder) handleSchemaTypesList() (*sdk.CallToolResult, error) {
	types := []string{"AVRO", "JSON", "PROTOBUF"}
	return b.marshalResponse(types)
}

func buildKafkaSchemaRegistryInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[kafkaSchemaRegistryInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}

	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}
	setSchemaDescription(schema, "resource", kafkaSchemaRegistryResourceDesc)
	setSchemaDescription(schema, "operation", kafkaSchemaRegistryOperationDesc)
	setSchemaDescription(schema, "subject", kafkaSchemaRegistrySubjectDesc)
	setSchemaDescription(schema, "version", kafkaSchemaRegistryVersionDesc)
	setSchemaDescription(schema, "compatibility", kafkaSchemaRegistryCompatibilityDesc)
	setSchemaDescription(schema, "schemaType", kafkaSchemaRegistrySchemaTypeDesc)
	setSchemaDescription(schema, "schema", kafkaSchemaRegistrySchemaDesc)

	normalizeAdditionalProperties(schema)
	return schema, nil
}
