package schema

import (
	"fmt"

	cliutils "github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
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

func (c *BooleanConverter) ToMCPToolInputSchemaProperties(schemaInfo *cliutils.SchemaInfo) ([]mcp.ToolOption, error) {
	if schemaInfo.Type != "BOOLEAN" {
		return nil, fmt.Errorf("expected BOOLEAN schema, got %s", schemaInfo.Type)
	}

	return []mcp.ToolOption{
		mcp.WithBoolean(c.ParamName, mcp.Description(fmt.Sprintf("The input schema is a %s schema", schemaInfo.Type)), mcp.Required()),
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
