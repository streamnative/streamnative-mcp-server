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
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
)

// BooleanConverter handles the conversion for Pulsar BOOLEAN schemas.
type BooleanConverter struct {
	BaseConverter
}

// NewBooleanConverter creates a new instance of BooleanConverter.
func NewBooleanConverter() *BooleanConverter {
	return &BooleanConverter{
		BaseConverter: BaseConverter{
			ParamName: ParamName,
		},
	}
}

func (c *BooleanConverter) ToMCPToolInputSchemaProperties(schemaInfo *cliutils.SchemaInfo) (*jsonschema.Schema, error) {
	if schemaInfo.Type != "BOOLEAN" {
		return nil, fmt.Errorf("expected BOOLEAN schema, got %s", schemaInfo.Type)
	}

	props := jsonschema.NewProperties()
	props.Set(c.ParamName, &jsonschema.Schema{
		Type:        "boolean",
		Description: fmt.Sprintf("The input schema is a %s schema", schemaInfo.Type),
	})

	return &jsonschema.Schema{
		Type:       "object",
		Properties: props,
		Required:   []string{c.ParamName},
	}, nil
}

func (c *BooleanConverter) SerializeMCPRequestToPulsarPayload(arguments map[string]any, targetPulsarSchemaInfo *cliutils.SchemaInfo) ([]byte, error) {
	if err := c.ValidateArguments(arguments, targetPulsarSchemaInfo); err != nil {
		return nil, fmt.Errorf("arguments validation failed: %w", err)
	}

	payload, err := common.RequiredParam[bool](arguments, c.ParamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get payload: %w", err)
	}

	return []byte(fmt.Sprintf("%t", payload)), nil
}

func (c *BooleanConverter) ValidateArguments(arguments map[string]any, targetPulsarSchemaInfo *cliutils.SchemaInfo) error {
	if targetPulsarSchemaInfo.Type != "BOOLEAN" {
		return fmt.Errorf("expected BOOLEAN schema, got %s", targetPulsarSchemaInfo.Type)
	}

	_, err := common.RequiredParam[bool](arguments, c.ParamName)
	if err != nil {
		return fmt.Errorf("failed to get payload: %w", err)
	}

	return nil
}
