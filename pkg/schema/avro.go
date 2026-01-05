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
	"fmt"

	"github.com/invopop/jsonschema"
	cliutils "github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
)

type AvroConverter struct {
	BaseConverter
}

func NewAvroConverter() *AvroConverter {
	return &AvroConverter{}
}

func (c *AvroConverter) ToMCPToolInputSchemaProperties(schemaInfo *cliutils.SchemaInfo) (*jsonschema.Schema, error) {
	if schemaInfo.Type != "AVRO" {
		return nil, fmt.Errorf("expected AVRO schema, got %s", schemaInfo.Type)
	}
	return processAvroSchemaStringToMCPToolInput(string(schemaInfo.Schema))
}

func (c *AvroConverter) SerializeMCPRequestToPulsarPayload(arguments map[string]any, targetPulsarSchemaInfo *cliutils.SchemaInfo) ([]byte, error) {
	if err := c.ValidateArguments(arguments, targetPulsarSchemaInfo); err != nil {
		return nil, fmt.Errorf("arguments validation failed: %w", err)
	}
	return serializeArgumentsToAvroBinary(arguments, string(targetPulsarSchemaInfo.Schema))
}

func (c *AvroConverter) ValidateArguments(arguments map[string]any, targetPulsarSchemaInfo *cliutils.SchemaInfo) error {
	if targetPulsarSchemaInfo.Type != "AVRO" {
		return fmt.Errorf("expected AVRO schema for validation, got %s", targetPulsarSchemaInfo.Type)
	}
	return validateArgumentsAgainstAvroSchemaString(arguments, string(targetPulsarSchemaInfo.Schema))
}
