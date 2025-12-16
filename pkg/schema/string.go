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

	cliutils "github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
)

// StringConverter handles the conversion for Pulsar STRING schemas.
type StringConverter struct {
	BaseConverter
}

// NewStringConverter creates a new instance of StringConverter.
func NewStringConverter() *StringConverter {
	return &StringConverter{
		BaseConverter: BaseConverter{
			ParamName: ParamName,
		},
	}
}

func (c *StringConverter) ToMCPToolInputSchemaProperties(schemaInfo *cliutils.SchemaInfo) ([]mcp.ToolOption, error) {
	if schemaInfo.Type != "STRING" && schemaInfo.Type != "BYTES" {
		return nil, fmt.Errorf("expected STRING or BYTES schema, got %s", schemaInfo.Type)
	}

	return []mcp.ToolOption{
		mcp.WithString(c.ParamName, mcp.Description(fmt.Sprintf("The input schema is a %s schema", schemaInfo.Type)), mcp.Required()),
	}, nil
}

func (c *StringConverter) SerializeMCPRequestToPulsarPayload(arguments map[string]any, targetPulsarSchemaInfo *cliutils.SchemaInfo) ([]byte, error) {
	if err := c.ValidateArguments(arguments, targetPulsarSchemaInfo); err != nil {
		return nil, fmt.Errorf("arguments validation failed: %w", err)
	}

	payload, err := common.RequiredParam[string](arguments, c.ParamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get payload: %w", err)
	}

	return []byte(payload), nil
}

func (c *StringConverter) ValidateArguments(arguments map[string]any, targetPulsarSchemaInfo *cliutils.SchemaInfo) error {
	if targetPulsarSchemaInfo.Type != "STRING" && targetPulsarSchemaInfo.Type != "BYTES" {
		return fmt.Errorf("expected STRING or BYTES schema, got %s", targetPulsarSchemaInfo.Type)
	}

	payload, err := common.RequiredParam[string](arguments, c.ParamName)
	if err != nil {
		return fmt.Errorf("failed to get payload: %w", err)
	}

	if payload == "" {
		return fmt.Errorf("payload cannot be empty")
	}

	return nil
}
