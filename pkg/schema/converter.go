// Copyright 2026 StreamNative
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
)

const (
	// ParamName is the default parameter name for payload arguments.
	ParamName = "payload"
)

// Converter defines schema conversion behaviors for MCP tools.
type Converter interface {
	ToMCPToolInputSchemaProperties(pulsarSchemaInfo *cliutils.SchemaInfo) ([]mcp.ToolOption, error)

	SerializeMCPRequestToPulsarPayload(arguments map[string]any, targetPulsarSchemaInfo *cliutils.SchemaInfo) ([]byte, error)

	ValidateArguments(arguments map[string]any, targetPulsarSchemaInfo *cliutils.SchemaInfo) error
}

// ConverterFactory returns a converter for the given schema type.
func ConverterFactory(schemaType string) (Converter, error) {
	switch schemaType {
	case "AVRO":
		return NewAvroConverter(), nil
	case "JSON":
		return NewJSONConverter(), nil
	case "STRING", "BYTES":
		return NewStringConverter(), nil
	case "INT8", "INT16", "INT32", "INT64", "FLOAT", "DOUBLE":
		return NewNumberConverter(), nil
	case "BOOLEAN":
		return NewBooleanConverter(), nil
	default:
		return nil, fmt.Errorf("unsupported schema type: %v", schemaType)
	}
}

// BaseConverter provides shared fields for schema converters.
type BaseConverter struct {
	ParamName string
}
