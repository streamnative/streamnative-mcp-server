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

package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	"github.com/twmb/franz-go/pkg/sr"
)

func KafkaAdminAddSchemaRegistryTools(s *server.MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeatureKafkaAdminSchemaRegistry)) && !slices.Contains(features, string(FeatureAll)) && !slices.Contains(features, string(FeatureAllKafka)) && !slices.Contains(features, string(FeatureKafkaAdmin)) {
		return
	}

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
		"Schema Registry is a component that stores and manages schemas for Kafka message data. " +
		"It enforces schema compatibility during evolution, enabling safe, backward-compatible changes to data structures. " +
		"Key concepts in Schema Registry:\n\n" +
		"- Subject: A named schema, typically following the pattern '<topic-name>-key' or '<topic-name>-value'\n" +
		"- Version: Each schema evolution creates a new numbered version\n" +
		"- Schema: The actual data structure definition (in AVRO, JSON Schema, or Protobuf format)\n" +
		"- Compatibility: Rules that control how schemas can evolve (BACKWARD, FORWARD, FULL, NONE)\n\n" +
		"Schema Registry helps ensure data consistency, enables schema evolution governance, " +
		"and reduces payload sizes by transmitting schema IDs instead of full schemas.\n\n" +
		"Usage Examples:\n\n" +
		"1. List all schema subjects in the registry:\n" +
		"   resource: \"subjects\"\n" +
		"   operation: \"list\"\n\n" +
		"2. List all versions for a specific subject:\n" +
		"   resource: \"versions\"\n" +
		"   operation: \"list\"\n" +
		"   subject: \"my-topic-value\"\n\n" +
		"3. List all supported schema types:\n" +
		"   resource: \"types\"\n" +
		"   operation: \"list\"\n\n" +
		"4. Get the latest schema for a subject:\n" +
		"   resource: \"subject\"\n" +
		"   operation: \"get\"\n" +
		"   subject: \"my-topic-value\"\n\n" +
		"5. Get a specific schema version:\n" +
		"   resource: \"version\"\n" +
		"   operation: \"get\"\n" +
		"   subject: \"my-topic-value\"\n" +
		"   version: 1\n\n" +
		"6. Get compatibility level for a subject:\n" +
		"   resource: \"compatibility\"\n" +
		"   operation: \"get\"\n" +
		"   subject: \"my-topic-value\"\n\n" +
		"7. Set compatibility level for a subject:\n" +
		"   resource: \"compatibility\"\n" +
		"   operation: \"set\"\n" +
		"   subject: \"my-topic-value\"\n" +
		"   compatibility: \"BACKWARD\"\n\n" +
		"8. Register a new AVRO schema for a subject:\n" +
		"   resource: \"subject\"\n" +
		"   operation: \"create\"\n" +
		"   subject: \"my-topic-value\"\n" +
		"   type: \"AVRO\"\n" +
		"   schema: \"{\\\"type\\\":\\\"record\\\",\\\"name\\\":\\\"MyRecord\\\",\\\"fields\\\":[{\\\"name\\\":\\\"id\\\",\\\"type\\\":\\\"int\\\"},{\\\"name\\\":\\\"name\\\",\\\"type\\\":\\\"string\\\"}]}\"\n\n" +
		"9. Delete a subject (all versions):\n" +
		"   resource: \"subject\"\n" +
		"   operation: \"delete\"\n" +
		"   subject: \"my-topic-value\"\n\n" +
		"This tool requires appropriate Schema Registry permissions."

	kafkaSchemaRegistryTool := mcp.NewTool("kafka_admin_sr",
		mcp.WithDescription(toolDesc),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description(resourceDesc),
		),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
		),
		mcp.WithString("subject",
			mcp.Description("The name of the Schema Registry subject to operate on. "+
				"Required for operations on specific subjects, versions, and compatibility settings. "+
				"Typically follows the pattern '<topic-name>-key' for key schemas or '<topic-name>-value' for value schemas. "+
				"Case-sensitive and must match exactly with the subject in Schema Registry.")),
		mcp.WithString("compatibility",
			mcp.Description("The compatibility level to set for schema evolution. "+
				"Required for the 'set' operation on 'compatibility' resource. "+
				"Valid values are:\n"+
				"- BACKWARD: New schema can read data written with previous schema (default)\n"+
				"- FORWARD: Previous schema can read data written with new schema\n"+
				"- FULL: Both backward and forward compatibility\n"+
				"- NONE: No compatibility checks enforced")),
		mcp.WithString("schema",
			mcp.Description("The schema definition to register. "+
				"Required for the 'create' operation on 'subject' resource. "+
				"Must be a valid schema according to the specified 'type' (AVRO, JSON, or PROTOBUF). "+
				"For AVRO, use JSON format with escaped quotes. "+
				"Example: '{\\\"type\\\":\\\"record\\\",\\\"fields\\\":[{\\\"name\\\":\\\"id\\\",\\\"type\\\":\\\"int\\\"}]}'")),
		mcp.WithString("version",
			mcp.Description("The schema version to retrieve or delete. "+
				"Required for 'get' operation on 'version' resource and optional for 'delete' operation on 'subject' resource. "+
				"An integer specifying the version number, where 1 is the first version. "+
				"Use 'latest' or -1 to reference the most recent version.")),
		mcp.WithString("type",
			mcp.Description("The schema format type. "+
				"Required for the 'create' operation on 'subject' resource. "+
				"Must be one of:\n"+
				"- AVRO: Apache Avro schema format\n"+
				"- JSON: JSON Schema format\n"+
				"- PROTOBUF: Protocol Buffers schema format\n"+
				"The format must match the provided schema definition.")),
	)

	s.AddTool(kafkaSchemaRegistryTool, handleKafkaSchemaRegistryTool(readOnly))
}

func handleKafkaSchemaRegistryTool(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		if readOnly && (operation == "set" || operation == "create" || operation == "delete") {
			return mcp.NewToolResultError("Write operations are not allowed in read-only mode"), nil
		}

		// Create the admin client
		admin, err := kafka.GetKafkaSchemaRegistryClient()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
		}

		// Dispatch based on resource and operation
		switch resource {
		case "subjects":
			switch operation {
			case "list":
				return handleKafkaSubjectsList(ctx, admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'subjects': %s", operation)), nil
			}
		case "versions":
			switch operation {
			case "list":
				return handleKafkaVersionsList(ctx, admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'versions': %s", operation)), nil
			}
		case "types":
			switch operation {
			case "list":
				return handleKafkaTypesList(ctx, admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'types': %s", operation)), nil
			}
		case "subject":
			switch operation {
			case "get":
				return handleKafkaSubjectGet(ctx, admin, request)
			case "create":
				return handleKafkaSubjectCreate(ctx, admin, request)
			case "delete":
				return handleKafkaSubjectDelete(ctx, admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'subject': %s", operation)), nil
			}
		case "version":
			switch operation {
			case "get":
				return handleKafkaVersionGet(ctx, admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'version': %s", operation)), nil
			}
		case "compatibility":
			switch operation {
			case "get":
				return handleKafkaCompatibilityGet(ctx, admin, request)
			case "set":
				return handleKafkaCompatibilitySet(ctx, admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'compatibility': %s", operation)), nil
			}
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resource: %s. Available resources: subjects, subject, versions, version, compatibility, types", resource)), nil
		}
	}
}

func handleKafkaSubjectsList(ctx context.Context, admin *sr.Client, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subjects, err := admin.Subjects(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get subjects: %v", err)), nil
	}

	jsonBytes, err := json.Marshal(subjects)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func handleKafkaVersionsList(ctx context.Context, admin *sr.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subject, err := requiredParam[string](request.Params.Arguments, "subject")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get subject: %v", err)), nil
	}

	versions, err := admin.SubjectVersions(ctx, subject)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get versions: %v", err)), nil
	}

	jsonBytes, err := json.Marshal(versions)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func handleKafkaTypesList(ctx context.Context, admin *sr.Client, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	types, err := admin.SupportedTypes(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get types: %v", err)), nil
	}

	jsonBytes, err := json.Marshal(types)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func handleKafkaSubjectGet(ctx context.Context, admin *sr.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subject, err := requiredParam[string](request.Params.Arguments, "subject")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get subject: %v", err)), nil
	}

	schemas, err := admin.Schemas(ctx, subject)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get schemas: %v", err)), nil
	}

	jsonBytes, err := json.Marshal(schemas)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func handleKafkaSubjectCreate(ctx context.Context, admin *sr.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subject, err := requiredParam[string](request.Params.Arguments, "subject")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get subject: %v", err)), nil
	}

	schema, err := requiredParam[string](request.Params.Arguments, "schema")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get schema: %v", err)), nil
	}

	typeString, err := requiredParam[string](request.Params.Arguments, "type")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get schema type: %v", err)), nil
	}

	var schemaType sr.SchemaType
	err = schemaType.UnmarshalText([]byte(typeString))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to unmarshal schema type: %v", err)), nil
	}

	schemaData := sr.Schema{
		Schema: schema,
		Type:   schemaType,
	}

	createdSchema, err := admin.CreateSchema(ctx, subject, schemaData)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create schema: %v", err)), nil
	}

	jsonBytes, err := json.Marshal(createdSchema)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func handleKafkaSubjectDelete(ctx context.Context, admin *sr.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subject, err := requiredParam[string](request.Params.Arguments, "subject")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get subject: %v", err)), nil
	}

	version, err := requiredParam[int](request.Params.Arguments, "version")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get version: %v", err)), nil
	}

	err = admin.DeleteSchema(ctx, subject, version, sr.SoftDelete)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete schema: %v", err)), nil
	}

	return mcp.NewToolResultText("Schema deleted successfully"), nil
}

func handleKafkaVersionGet(ctx context.Context, admin *sr.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subject, err := requiredParam[string](request.Params.Arguments, "subject")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get subject: %v", err)), nil
	}

	version, err := requiredParam[int](request.Params.Arguments, "version")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get version: %v", err)), nil
	}

	schema, err := admin.SchemaByVersion(ctx, subject, version)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get schema: %v", err)), nil
	}

	jsonBytes, err := json.Marshal(schema)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func handleKafkaCompatibilityGet(ctx context.Context, admin *sr.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subject, err := requiredParam[string](request.Params.Arguments, "subject")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get subject: %v", err)), nil
	}

	compatibility := admin.Compatibility(ctx, subject)

	jsonBytes, err := json.Marshal(compatibility)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func handleKafkaCompatibilitySet(ctx context.Context, admin *sr.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subject, err := requiredParam[string](request.Params.Arguments, "subject")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get subject: %v", err)), nil
	}

	compatibility, err := requiredParam[string](request.Params.Arguments, "compatibility")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get compatibility: %v", err)), nil
	}

	var compatibilityLevel sr.CompatibilityLevel
	err = compatibilityLevel.UnmarshalText([]byte(compatibility))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to unmarshal compatibility: %v", err)), nil
	}

	response := admin.SetCompatibility(ctx, sr.SetCompatibility{
		Level: compatibilityLevel,
	}, subject)

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}
