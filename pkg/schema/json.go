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

package schema

import (
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	cliutils "github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
)

// JSONConverter handles the conversion for Pulsar JSON schemas.
// It relies on the underlying AVRO schema definition for structure and validation,
// but serializes to a standard JSON text payload.
type JSONConverter struct {
	BaseConverter
}

// NewJSONConverter creates a new instance of JSONConverter.
func NewJSONConverter() *JSONConverter {
	return &JSONConverter{}
}

// ToMCPToolInputSchemaProperties converts the Pulsar JSON SchemaInfo (which is AVRO based)
// to jsonschema.Schema for tool input.
func (c *JSONConverter) ToMCPToolInputSchemaProperties(schemaInfo *cliutils.SchemaInfo) (*jsonschema.Schema, error) {
	if schemaInfo.Type != "JSON" {
		return nil, fmt.Errorf("expected JSON schema, got %s", schemaInfo.Type)
	}
	// Delegate to the core AVRO processing function from avro_core.go.
	return processAvroSchemaStringToMCPToolInput(string(schemaInfo.Schema))
}

// SerializeMCPRequestToPulsarPayload validates arguments against the underlying AVRO schema definition
// and then serializes them to a JSON text payload for Pulsar.
func (c *JSONConverter) SerializeMCPRequestToPulsarPayload(arguments map[string]any, targetPulsarSchemaInfo *cliutils.SchemaInfo) ([]byte, error) {
	if err := c.ValidateArguments(arguments, targetPulsarSchemaInfo); err != nil {
		return nil, fmt.Errorf("arguments validation failed: %w", err)
	}

	// Serialize arguments to standard JSON []byte.
	jsonData, err := json.Marshal(arguments)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments to JSON: %w", err)
	}

	return jsonData, nil
}

// ValidateArguments validates the given arguments against the Pulsar JSON SchemaInfo's
// underlying AVRO schema definition.
func (c *JSONConverter) ValidateArguments(arguments map[string]any, targetPulsarSchemaInfo *cliutils.SchemaInfo) error {
	if targetPulsarSchemaInfo.Type != "JSON" {
		return fmt.Errorf("expected JSON schema for validation, got %s", targetPulsarSchemaInfo.Type)
	}
	// The schemaInfo.Schema for JSON type is the AVRO schema string definition.
	// Delegate to the core AVRO validation function from avro_core.go.
	return validateArgumentsAgainstAvroSchemaString(arguments, string(targetPulsarSchemaInfo.Schema))
}
