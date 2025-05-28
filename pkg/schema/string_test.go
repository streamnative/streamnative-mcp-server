package schema

import (
	"fmt"
	"testing"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

// Helper function to create SchemaInfo for string tests
func newStringSchemaInfo(schemaType string) *utils.SchemaInfo {
	return &utils.SchemaInfo{
		Type:   schemaType,
		Schema: []byte{},
	}
}

func TestNewStringConverter(t *testing.T) {
	converter := NewStringConverter()
	assert.NotNil(t, converter)
	assert.Equal(t, ParamName, converter.ParamName, "ParamName should be initialized to the package constant")
}

func TestStringConverter_ToMCPToolInputSchemaProperties(t *testing.T) {
	converter := NewStringConverter()

	tests := []struct {
		name       string
		schemaInfo *utils.SchemaInfo
		wantOpts   []mcp.ToolOption
		wantErr    bool
		errContain string
	}{
		{
			name:       "Valid STRING schema",
			schemaInfo: newStringSchemaInfo("STRING"),
			wantOpts: []mcp.ToolOption{
				mcp.WithString(ParamName, mcp.Description(fmt.Sprintf("The input schema is a %s schema", "STRING")), mcp.Required()),
			},
			wantErr: false,
		},
		{
			name:       "Valid BYTES schema",
			schemaInfo: newStringSchemaInfo("BYTES"),
			wantOpts: []mcp.ToolOption{
				mcp.WithString(ParamName, mcp.Description(fmt.Sprintf("The input schema is a %s schema", "BYTES")), mcp.Required()),
			},
			wantErr: false,
		},
		{
			name:       "Invalid schema type (JSON)",
			schemaInfo: newStringSchemaInfo("JSON"),
			wantOpts:   nil,
			wantErr:    true,
			errContain: "expected STRING or BYTES schema, got JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOpts, err := converter.ToMCPToolInputSchemaProperties(tt.schemaInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToMCPToolInputSchemaProperties() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var expectedTool = mcp.NewTool("test", gotOpts...)
			var actualTool = mcp.NewTool("test", tt.wantOpts...)
			expectedToolInputSchemaJSON, _ := expectedTool.InputSchema.MarshalJSON()
			actualToolInputSchemaJSON, _ := actualTool.InputSchema.MarshalJSON()
			assert.Equal(t, string(expectedToolInputSchemaJSON), string(actualToolInputSchemaJSON))
			if tt.wantErr && err != nil {
				assert.Contains(t, err.Error(), tt.errContain)
			}
		})
	}
}

func TestStringConverter_ValidateArguments(t *testing.T) {
	converter := NewStringConverter()

	tests := []struct {
		name       string
		schemaInfo *utils.SchemaInfo
		args       map[string]any
		wantErr    bool
		errContain string
	}{
		{
			name:       "Valid arguments for STRING schema",
			schemaInfo: newStringSchemaInfo("STRING"),
			args:       map[string]any{ParamName: "hello world"},
			wantErr:    false,
		},
		{
			name:       "Valid arguments for BYTES schema",
			schemaInfo: newStringSchemaInfo("BYTES"),
			args:       map[string]any{ParamName: "bytes content"},
			wantErr:    false,
		},
		{
			name:       "Invalid schema type (JSON)",
			schemaInfo: newStringSchemaInfo("JSON"),
			args:       map[string]any{ParamName: "test"},
			wantErr:    true,
			errContain: "expected STRING or BYTES schema, got JSON",
		},
		{
			name:       "Missing payload argument",
			schemaInfo: newStringSchemaInfo("STRING"),
			args:       map[string]any{},
			wantErr:    true,
			errContain: "failed to get payload: missing required parameter: payload",
		},
		{
			name:       "Incorrect payload type (int instead of string)",
			schemaInfo: newStringSchemaInfo("STRING"),
			args:       map[string]any{ParamName: 123},
			wantErr:    true,
			errContain: "failed to get payload: parameter payload is not of type string",
		},
		{
			name:       "Empty string payload",
			schemaInfo: newStringSchemaInfo("STRING"),
			args:       map[string]any{ParamName: ""},
			wantErr:    true,
			errContain: "failed to get payload: missing required parameter: payload",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := converter.ValidateArguments(tt.args, tt.schemaInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateArguments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				assert.Contains(t, err.Error(), tt.errContain)
			}
		})
	}
}

func TestStringConverter_SerializeMCPRequestToPulsarPayload(t *testing.T) {
	converter := NewStringConverter()

	tests := []struct {
		name       string
		schemaInfo *utils.SchemaInfo
		args       map[string]any
		want       []byte
		wantErr    bool
		errContain string
	}{
		{
			name:       "Serialize 'hello' for STRING schema",
			schemaInfo: newStringSchemaInfo("STRING"),
			args:       map[string]any{ParamName: "hello"},
			want:       []byte("hello"),
			wantErr:    false,
		},
		{
			name:       "Serialize 'bytes content' for BYTES schema",
			schemaInfo: newStringSchemaInfo("BYTES"),
			args:       map[string]any{ParamName: "bytes content"},
			want:       []byte("bytes content"),
			wantErr:    false,
		},
		{
			name:       "Validation error (e.g., empty string)",
			schemaInfo: newStringSchemaInfo("STRING"),
			args:       map[string]any{ParamName: ""},
			want:       nil,
			wantErr:    true,
			errContain: "arguments validation failed",
		},
		{
			name:       "Validation error (e.g., missing payload)",
			schemaInfo: newStringSchemaInfo("STRING"),
			args:       map[string]any{},
			want:       nil,
			wantErr:    true,
			errContain: "arguments validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.SerializeMCPRequestToPulsarPayload(tt.args, tt.schemaInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeMCPRequestToPulsarPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
			if tt.wantErr && err != nil {
				assert.Contains(t, err.Error(), tt.errContain)
			}
		})
	}
}

// Future test functions will be added here.
