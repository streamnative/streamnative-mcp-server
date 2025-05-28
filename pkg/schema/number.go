package schema

import (
	"fmt"
	"math"

	cliutils "github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
)

// NumberConverter handles the conversion for Pulsar numeric schemas (INT8, INT16, INT32, INT64, FLOAT, DOUBLE).
type NumberConverter struct {
	BaseConverter
}

// NewNumberConverter creates a new instance of NumberConverter.
func NewNumberConverter() *NumberConverter {
	return &NumberConverter{
		BaseConverter: BaseConverter{
			ParamName: ParamName,
		},
	}
}

func (c *NumberConverter) ToMCPToolInputSchemaProperties(schemaInfo *cliutils.SchemaInfo) ([]mcp.ToolOption, error) {
	if schemaInfo.Type != "INT8" && schemaInfo.Type != "INT16" && schemaInfo.Type != "INT32" && schemaInfo.Type != "INT64" && schemaInfo.Type != "FLOAT" && schemaInfo.Type != "DOUBLE" {
		return nil, fmt.Errorf("expected INT8, INT16, INT32, INT64, FLOAT, or DOUBLE schema, got %s", schemaInfo.Type)
	}

	return []mcp.ToolOption{
		mcp.WithNumber(c.ParamName, mcp.Description(fmt.Sprintf("The input schema is a %s schema", schemaInfo.Type)), mcp.Required()),
	}, nil
}

func (c *NumberConverter) SerializeMCPRequestToPulsarPayload(arguments map[string]any, targetPulsarSchemaInfo *cliutils.SchemaInfo) ([]byte, error) {
	if err := c.ValidateArguments(arguments, targetPulsarSchemaInfo); err != nil {
		return nil, fmt.Errorf("arguments validation failed: %w", err)
	}

	payload, err := common.RequiredParam[float64](arguments, c.ParamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get payload: %w", err)
	}

	switch targetPulsarSchemaInfo.Type {
	case "INT8":
		return []byte(fmt.Sprintf("%d", int8(payload))), nil
	case "INT16":
		return []byte(fmt.Sprintf("%d", int16(payload))), nil
	case "INT32":
		return []byte(fmt.Sprintf("%d", int32(payload))), nil
	case "INT64":
		return []byte(fmt.Sprintf("%d", int64(payload))), nil
	case "FLOAT":
		return []byte(fmt.Sprintf("%f", payload)), nil
	case "DOUBLE":
		return []byte(fmt.Sprintf("%f", payload)), nil
	default:
		return nil, fmt.Errorf("unsupported schema type: %s", targetPulsarSchemaInfo.Type)
	}
}

func (c *NumberConverter) ValidateArguments(arguments map[string]any, targetPulsarSchemaInfo *cliutils.SchemaInfo) error {
	if targetPulsarSchemaInfo.Type != "INT8" && targetPulsarSchemaInfo.Type != "INT16" && targetPulsarSchemaInfo.Type != "INT32" && targetPulsarSchemaInfo.Type != "INT64" && targetPulsarSchemaInfo.Type != "FLOAT" && targetPulsarSchemaInfo.Type != "DOUBLE" {
		return fmt.Errorf("expected INT8, INT16, INT32, INT64, FLOAT, or DOUBLE schema, got %s", targetPulsarSchemaInfo.Type)
	}

	payload, err := common.RequiredParam[float64](arguments, c.ParamName)
	if err != nil {
		return fmt.Errorf("failed to get payload: %w", err)
	}

	switch targetPulsarSchemaInfo.Type {
	case "INT8":
		if payload < math.MinInt8 || payload > math.MaxInt8 {
			return fmt.Errorf("payload out of range for INT8")
		}
	case "INT16":
		if payload < math.MinInt16 || payload > math.MaxInt16 {
			return fmt.Errorf("payload out of range for INT16")
		}
	case "INT32":
		if payload < math.MinInt32 || payload > math.MaxInt32 {
			return fmt.Errorf("payload out of range for INT32")
		}
	case "INT64":
		if payload < math.MinInt64 || payload > math.MaxInt64 {
			return fmt.Errorf("payload out of range for INT64")
		}
	case "FLOAT":
		if payload < math.SmallestNonzeroFloat32 || payload > math.MaxFloat32 {
			return fmt.Errorf("payload out of range for FLOAT")
		}
	case "DOUBLE":
		if payload < math.SmallestNonzeroFloat64 || payload > math.MaxFloat64 {
			return fmt.Errorf("payload out of range for DOUBLE")
		}
	default:
		return fmt.Errorf("unsupported schema type: %s", targetPulsarSchemaInfo.Type)
	}

	return nil
}
