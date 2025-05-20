package schema

import (
	"fmt"

	cliutils "github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
)

const (
	ParamName = "payload"
)

type Converter interface {
	ToMCPToolInputSchemaProperties(pulsarSchemaInfo *cliutils.SchemaInfo) ([]mcp.ToolOption, error)

	SerializeMCPRequestToPulsarPayload(arguments map[string]any, targetPulsarSchemaInfo *cliutils.SchemaInfo) ([]byte, error)

	ValidateArguments(arguments map[string]any, targetPulsarSchemaInfo *cliutils.SchemaInfo) error
}

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

type BaseConverter struct {
	ParamName string
}
