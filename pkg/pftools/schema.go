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

var DefaultStringSchema = &mcp.ToolInputSchema{
	Type: "object",
	Properties: map[string]interface{}{
		"payload": map[string]interface{}{
			"type":        "string",
			"description": "The payload of the message, in plain text format",
		},
	},
}

// GetSchemaFromTopic retrieves schema information from a topic
func GetSchemaFromTopic(admin cmdutils.Client, topic string) (*SchemaInfo, error) {
	if admin == nil {
		return nil, fmt.Errorf("failed to get schema from topic '%s': mcp server is not initialized", topic)
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
		return DefaultStringSchema, nil
	}

	// Handle different schema types
	switch schemaInfo.Type {
	case "JSON":
		return convertComplexSchemaToToolInput(schemaInfo)
	case "AVRO", "PROTOBUF", "PROTOBUF_NATIVE":
		return nil, fmt.Errorf("AVRO, PROTOBUF and PROTOBUF_NATIVE schema is not supported")
	default:
		return DefaultStringSchema, nil
	}
}

// convertComplexSchemaToToolInput handles conversion of complex schema types
func convertComplexSchemaToToolInput(schemaInfo *SchemaInfo) (*mcp.ToolInputSchema, error) {
	if schemaInfo.Definition == nil {
		return DefaultStringSchema, nil
	}

	definitionString, err := json.Marshal(schemaInfo.Definition)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal schema definition: %w", err)
	}

	// For JSON schemas, use the definition directly
	return &mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"payload": map[string]interface{}{
				"type":        "string",
				"description": "The payload of the message, in JSON String format, the schema of the payload is: " + string(definitionString),
			},
		},
	}, nil
}
