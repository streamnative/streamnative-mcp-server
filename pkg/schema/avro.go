package schema

import (
	"fmt"
	// "reflect" // No longer needed here

	cliutils "github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
)

type AvroConverter struct {
	BaseConverter
}

func NewAvroConverter() *AvroConverter {
	return &AvroConverter{}
}

func (c *AvroConverter) ToMCPToolInputSchemaProperties(schemaInfo *cliutils.SchemaInfo) ([]mcp.ToolOption, error) {
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
