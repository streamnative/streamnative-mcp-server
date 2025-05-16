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

package pftools

import (
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
)

// GetSchemaFromTopic retrieves schema information from a topic
func GetSchemaFromTopic(admin cmdutils.Client, topic string) (*SchemaInfo, error) {
	if admin == nil {
		return nil, fmt.Errorf("admin client is nil")
	}

	// Get schema info from topic
	schemaInfoWithVersion, err := admin.Schemas().GetSchemaInfoWithVersion(topic)
	if err != nil {
		return nil, fmt.Errorf("failed to get schema for topic '%s': %w", topic, err)
	}

	if schemaInfoWithVersion == nil || schemaInfoWithVersion.SchemaInfo == nil {
		return nil, fmt.Errorf("no schema found for topic '%s'", topic)
	}

	// Parse schema definition
	var definition map[string]interface{}
	if schemaInfoWithVersion.SchemaInfo.Schema != nil {
		if err := json.Unmarshal(schemaInfoWithVersion.SchemaInfo.Schema, &definition); err != nil {
			// If it's not a valid JSON, just create a string type schema
			definition = map[string]interface{}{
				"type": "string",
			}
		}
	} else {
		// Default to string type if no schema is provided
		definition = map[string]interface{}{
			"type": "string",
		}
	}

	return &SchemaInfo{
		Type:       string(schemaInfoWithVersion.SchemaInfo.Type),
		Definition: definition,
	}, nil
}

// ConvertSchemaToToolInput converts a schema to MCP tool input schema
func ConvertSchemaToToolInput(schemaInfo *SchemaInfo) (*mcp.ToolInputSchema, error) {
	if schemaInfo == nil {
		// Default to object with any fields if no schema is provided
		return &mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"data": map[string]interface{}{
					"type": "string",
				},
			},
		}, nil
	}

	// Handle different schema types
	switch schemaInfo.Type {
	case "STRING":
		return &mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"data": map[string]interface{}{
					"type": "string",
				},
			},
		}, nil
	case "JSON", "AVRO":
		// For complex types, try to extract the JSON schema
		return convertComplexSchemaToToolInput(schemaInfo)
	default:
		// For other types, default to string
		return &mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"data": map[string]interface{}{
					"type": "string",
				},
			},
		}, nil
	}
}

// convertComplexSchemaToToolInput handles conversion of complex schema types
func convertComplexSchemaToToolInput(schemaInfo *SchemaInfo) (*mcp.ToolInputSchema, error) {
	if schemaInfo.Definition == nil {
		return &mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"data": map[string]interface{}{
					"type": "string",
				},
			},
		}, nil
	}

	// For AVRO schemas, extract fields from the definition
	if schemaInfo.Type == "AVRO" {
		return convertAvroSchemaToToolInput(schemaInfo.Definition)
	}

	// For JSON schemas, use the definition directly
	return &mcp.ToolInputSchema{
		Type:       "object",
		Properties: schemaInfo.Definition,
	}, nil
}

// convertAvroSchemaToToolInput converts Avro schema to MCP tool input schema
func convertAvroSchemaToToolInput(definition map[string]interface{}) (*mcp.ToolInputSchema, error) {
	fields, ok := definition["fields"]
	if !ok {
		return &mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"data": map[string]interface{}{
					"type": "string",
				},
			},
		}, nil
	}

	// Extract fields from Avro schema
	fieldsList, ok := fields.([]interface{})
	if !ok {
		return &mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"data": map[string]interface{}{
					"type": "string",
				},
			},
		}, nil
	}

	// Convert Avro fields to JSON schema properties
	properties := make(map[string]interface{})
	for _, field := range fieldsList {
		fieldMap, ok := field.(map[string]interface{})
		if !ok {
			continue
		}

		name, ok := fieldMap["name"].(string)
		if !ok {
			continue
		}

		// Convert Avro type to JSON schema type
		fieldType, ok := fieldMap["type"]
		if !ok {
			continue
		}

		properties[name] = convertAvroTypeToJsonSchema(fieldType)
	}

	return &mcp.ToolInputSchema{
		Type:       "object",
		Properties: properties,
	}, nil
}

// convertAvroTypeToJsonSchema converts Avro type to JSON schema type
func convertAvroTypeToJsonSchema(avroType interface{}) map[string]interface{} {
	switch t := avroType.(type) {
	case string:
		// Handle primitive types
		switch t {
		case "string":
			return map[string]interface{}{
				"type": "string",
			}
		case "int", "long":
			return map[string]interface{}{
				"type": "integer",
			}
		case "float", "double":
			return map[string]interface{}{
				"type": "number",
			}
		case "boolean":
			return map[string]interface{}{
				"type": "boolean",
			}
		default:
			return map[string]interface{}{
				"type": "string",
			}
		}
	case map[string]interface{}:
		// Handle complex types
		typeStr, ok := t["type"].(string)
		if !ok {
			return map[string]interface{}{
				"type": "object",
			}
		}

		switch typeStr {
		case "record":
			fields, ok := t["fields"].([]interface{})
			if !ok {
				return map[string]interface{}{
					"type": "object",
				}
			}

			properties := make(map[string]interface{})
			for _, field := range fields {
				fieldMap, ok := field.(map[string]interface{})
				if !ok {
					continue
				}

				name, ok := fieldMap["name"].(string)
				if !ok {
					continue
				}

				fieldType, ok := fieldMap["type"]
				if !ok {
					continue
				}

				properties[name] = convertAvroTypeToJsonSchema(fieldType)
			}

			return map[string]interface{}{
				"type":       "object",
				"properties": properties,
			}
		case "array":
			items, ok := t["items"]
			if !ok {
				return map[string]interface{}{
					"type":  "array",
					"items": map[string]interface{}{"type": "string"},
				}
			}

			return map[string]interface{}{
				"type":  "array",
				"items": convertAvroTypeToJsonSchema(items),
			}
		default:
			return map[string]interface{}{
				"type": "string",
			}
		}
	default:
		return map[string]interface{}{
			"type": "string",
		}
	}
}
