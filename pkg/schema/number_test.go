package schema

import (
	"fmt"
	"math"
	"testing"

	"github.com/apache/pulsar-client-go/pulsar"
	cliutils "github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

// Helper function to create SchemaInfo for number tests.
// Based on persistent linter errors, cliutils.SchemaInfo.Type is likely a string.
func newNumberSchemaInfo(schemaType pulsar.SchemaType) *cliutils.SchemaInfo {
	return &cliutils.SchemaInfo{
		Type:   GetSchemaType(schemaType), // Convert pulsar.SchemaType to string for cliutils.SchemaInfo.Type
		Schema: []byte{},
	}
}

func TestNewNumberConverter(t *testing.T) {
	converter := NewNumberConverter()
	assert.Equal(t, "payload", converter.ParamName)
}

func TestNumberConverter_ToMCPToolInputSchemaProperties(t *testing.T) {
	converter := NewNumberConverter()
	tests := []struct {
		name             string
		schemaInfo       *cliutils.SchemaInfo // This is schema.SchemaInfo which has Type as string
		expectedProps    []mcp.ToolOption
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name: "INT8 type",
			// newNumberSchemaInfo now correctly sets schemaInfo.Type as string (e.g., "INT8")
			schemaInfo:    newNumberSchemaInfo(pulsar.INT8),
			expectedProps: []mcp.ToolOption{mcp.WithNumber("payload", mcp.Description("The input schema is a INT8 schema"), mcp.Required())},
			expectError:   false,
		},
		{
			name:          "INT16 type",
			schemaInfo:    newNumberSchemaInfo(pulsar.INT16),
			expectedProps: []mcp.ToolOption{mcp.WithNumber("payload", mcp.Description("The input schema is a INT16 schema"), mcp.Required())},
			expectError:   false,
		},
		{
			name:          "INT32 type",
			schemaInfo:    newNumberSchemaInfo(pulsar.INT32),
			expectedProps: []mcp.ToolOption{mcp.WithNumber("payload", mcp.Description("The input schema is a INT32 schema"), mcp.Required())},
			expectError:   false,
		},
		{
			name:          "INT64 type",
			schemaInfo:    newNumberSchemaInfo(pulsar.INT64),
			expectedProps: []mcp.ToolOption{mcp.WithNumber("payload", mcp.Description("The input schema is a INT64 schema"), mcp.Required())},
			expectError:   false,
		},
		{
			name:          "FLOAT type",
			schemaInfo:    newNumberSchemaInfo(pulsar.FLOAT),
			expectedProps: []mcp.ToolOption{mcp.WithNumber("payload", mcp.Description("The input schema is a FLOAT schema"), mcp.Required())},
			expectError:   false,
		},
		{
			name:          "DOUBLE type",
			schemaInfo:    newNumberSchemaInfo(pulsar.DOUBLE),
			expectedProps: []mcp.ToolOption{mcp.WithNumber("payload", mcp.Description("The input schema is a DOUBLE schema"), mcp.Required())},
			expectError:   false,
		},
		{
			name:             "Unsupported type STRING",
			schemaInfo:       newNumberSchemaInfo(pulsar.STRING),
			expectedProps:    nil,
			expectError:      true,
			expectedErrorMsg: "expected INT8, INT16, INT32, INT64, FLOAT, or DOUBLE schema, got STRING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			props, err := converter.ToMCPToolInputSchemaProperties(tt.schemaInfo)
			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedErrorMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrorMsg)
				}
				assert.Nil(t, props)
			} else {
				assert.NoError(t, err)
				var expectedTool = mcp.NewTool("test", tt.expectedProps...)
				var actualTool = mcp.NewTool("test", props...)
				expectedToolSchemaJSON, _ := expectedTool.InputSchema.MarshalJSON()
				actualToolSchemaJSON, _ := actualTool.InputSchema.MarshalJSON()
				assert.Equal(t, string(expectedToolSchemaJSON), string(actualToolSchemaJSON), "ToolOptions mismatch")
			}
		})
	}
}

func TestNumberConverter_ValidateArguments(t *testing.T) {
	converter := NewNumberConverter()
	type testArgs struct {
		name             string
		schemaInfo       *cliutils.SchemaInfo
		args             map[string]any
		expectError      bool
		expectedErrorMsg string
	}

	tests := []testArgs{}

	numericTypes := []struct {
		pulsarType pulsar.SchemaType
		minVal     float64
		maxVal     float64
	}{
		{pulsar.INT8, -128, 127},
		{pulsar.INT16, -32768, 32767},
		{pulsar.INT32, -2147483648, 2147483647},
		{pulsar.INT64, -9007199254740991, 9007199254740991},
		{pulsar.FLOAT, math.SmallestNonzeroFloat32, math.MaxFloat32},
		{pulsar.DOUBLE, math.SmallestNonzeroFloat64, math.MaxFloat64},
	}

	for _, nt := range numericTypes {
		schemaTypeName := GetSchemaType(nt.pulsarType)
		tests = append(tests, []testArgs{
			{
				name:        fmt.Sprintf("%s - valid value", schemaTypeName),
				schemaInfo:  newNumberSchemaInfo(nt.pulsarType),
				args:        map[string]any{"payload": (nt.minVal + nt.maxVal) / 2},
				expectError: false,
			},
			{
				name:        fmt.Sprintf("%s - min value", schemaTypeName),
				schemaInfo:  newNumberSchemaInfo(nt.pulsarType),
				args:        map[string]any{"payload": nt.minVal},
				expectError: false,
			},
			{
				name:        fmt.Sprintf("%s - max value", schemaTypeName),
				schemaInfo:  newNumberSchemaInfo(nt.pulsarType),
				args:        map[string]any{"payload": nt.maxVal},
				expectError: false,
			},
		}...)
		if nt.pulsarType != pulsar.INT64 && nt.pulsarType != pulsar.FLOAT && nt.pulsarType != pulsar.DOUBLE {
			tests = append(tests, []testArgs{
				{
					name:             fmt.Sprintf("%s - value below min", schemaTypeName),
					schemaInfo:       newNumberSchemaInfo(nt.pulsarType),
					args:             map[string]any{"payload": nt.minVal - 1},
					expectError:      true,
					expectedErrorMsg: fmt.Sprintf("payload out of range for %s", schemaTypeName),
				},
				{
					name:             fmt.Sprintf("%s - value above max", schemaTypeName),
					schemaInfo:       newNumberSchemaInfo(nt.pulsarType),
					args:             map[string]any{"payload": nt.maxVal + 1},
					expectError:      true,
					expectedErrorMsg: fmt.Sprintf("payload out of range for %s", schemaTypeName),
				},
			}...)
		}
	}

	// Common error cases for one type (e.g., INT32) - these should behave similarly for others
	tests = append(tests, []testArgs{
		{
			name:             "INT32 - missing payload",
			schemaInfo:       newNumberSchemaInfo(pulsar.INT32),
			args:             map[string]any{},
			expectError:      true,
			expectedErrorMsg: "failed to get payload: missing required parameter: payload",
		},
		{
			name:             "INT32 - wrong payload type (string)",
			schemaInfo:       newNumberSchemaInfo(pulsar.INT32),
			args:             map[string]any{"payload": "not a number"},
			expectError:      true,
			expectedErrorMsg: "failed to get payload: parameter payload is not of type float64",
		},
		{
			name:             "INT32 - wrong schemaInfo.Type (e.g. STRING)",
			schemaInfo:       newNumberSchemaInfo(pulsar.STRING),
			args:             map[string]any{"payload": 123},
			expectError:      true,
			expectedErrorMsg: "expected INT8, INT16, INT32, INT64, FLOAT, or DOUBLE schema, got STRING",
		},
	}...)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := converter.ValidateArguments(tt.args, tt.schemaInfo)
			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedErrorMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNumberConverter_SerializeMCPRequestToPulsarPayload(t *testing.T) {
	converter := NewNumberConverter()
	type testArgs struct {
		name             string
		schemaInfo       *cliutils.SchemaInfo
		args             map[string]any
		expectedPayload  []byte
		expectError      bool
		expectedErrorMsg string
	}

	tests := []testArgs{}

	numericTypeCases := []struct {
		pulsarType    pulsar.SchemaType
		validPayload  any // Use any because direct float64 might lose precision for int64 for ex.
		expectedBytes []byte
		// Add specific out-of-range or invalid format cases if serialization handles them distinctly from validation
	}{
		{pulsar.INT8, float64(127), []byte("127")},
		{pulsar.INT8, float64(-128), []byte("-128")},
		{pulsar.INT16, float64(32767), []byte("32767")},
		{pulsar.INT32, float64(2147483647), []byte("2147483647")},
		// For INT64, using a number representable by float64 without precision loss
		{pulsar.INT64, float64(9007199254740991), []byte("9007199254740991")},
		{pulsar.FLOAT, float64(123.456), []byte(fmt.Sprintf("%f", float64(123.456)))},              // Note: float formatting can be tricky, use simple case
		{pulsar.DOUBLE, float64(1.23456789e10), []byte(fmt.Sprintf("%f", float64(1.23456789e10)))}, // Again, simple case for float formatting
	}

	for _, tc := range numericTypeCases {
		schemaTypeName := GetSchemaType(tc.pulsarType)
		tests = append(tests, testArgs{
			name:            fmt.Sprintf("%s - valid serialization", schemaTypeName),
			schemaInfo:      newNumberSchemaInfo(tc.pulsarType),
			args:            map[string]any{"payload": tc.validPayload},
			expectedPayload: tc.expectedBytes,
			expectError:     false,
		})
	}

	// Error cases (mostly delegation to ValidateArguments, so these check that path)
	tests = append(tests, []testArgs{
		{
			name:             "Error - INT32 - missing payload",
			schemaInfo:       newNumberSchemaInfo(pulsar.INT32),
			args:             map[string]any{},
			expectError:      true,
			expectedErrorMsg: "arguments validation failed: failed to get payload: missing required parameter",
		},
		{
			name:             "Error - INT32 - wrong payload type",
			schemaInfo:       newNumberSchemaInfo(pulsar.INT32),
			args:             map[string]any{"payload": "not a number"},
			expectError:      true,
			expectedErrorMsg: "arguments validation failed: failed to get payload: parameter payload is not of type float64",
		},
		{
			name:             "Error - INT32 - value out of range",
			schemaInfo:       newNumberSchemaInfo(pulsar.INT32),
			args:             map[string]any{"payload": float64(21474836480)}, // Clearly out of INT32 range
			expectError:      true,
			expectedErrorMsg: "arguments validation failed: payload out of range for INT32",
		},
		{
			name:             "Error - Unsupported Schema Type (STRING)",
			schemaInfo:       newNumberSchemaInfo(pulsar.STRING),
			args:             map[string]any{"payload": 123},
			expectError:      true,
			expectedErrorMsg: "arguments validation failed: expected INT8, INT16, INT32, INT64, FLOAT, or DOUBLE schema, got STRING",
		},
	}...)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := converter.SerializeMCPRequestToPulsarPayload(tt.args, tt.schemaInfo)
			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedErrorMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrorMsg)
				}
				assert.Nil(t, payload)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPayload, payload)
			}
		})
	}
}

// Future test functions will be added here.
