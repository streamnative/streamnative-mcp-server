package schema

import (
	"fmt"
	"testing"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

// Helper function to create SchemaInfo for boolean tests
func newBoolSchemaInfo(schemaType string) *utils.SchemaInfo {
	return &utils.SchemaInfo{
		Type:   schemaType,
		Schema: []byte{},
	}
}

func TestNewBooleanConverter(t *testing.T) {
	converter := NewBooleanConverter()
	assert.NotNil(t, converter)
	assert.Equal(t, ParamName, converter.ParamName, "ParamName should be initialized to the package constant")
}

func TestBooleanConverter_ToMCPToolInputSchemaProperties(t *testing.T) {
	converter := NewBooleanConverter()

	tests := []struct {
		name       string
		schemaInfo *utils.SchemaInfo
		wantOpts   []mcp.ToolOption
		wantErr    bool
	}{
		{
			name:       "Valid BOOLEAN schema",
			schemaInfo: newBoolSchemaInfo("BOOLEAN"),
			wantOpts: []mcp.ToolOption{
				mcp.WithBoolean(ParamName, mcp.Description(fmt.Sprintf("The input schema is a %s schema", "BOOLEAN")), mcp.Required()),
			},
			wantErr: false,
		},
		{
			name:       "Invalid schema type (STRING)",
			schemaInfo: newBoolSchemaInfo("STRING"),
			wantOpts:   nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOpts, err := converter.ToMCPToolInputSchemaProperties(tt.schemaInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToMCPToolInputSchemaProperties() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var expectedTool, actualTool mcp.Tool
			expectedTool = mcp.NewTool("test", tt.wantOpts...)
			actualTool = mcp.NewTool("test", gotOpts...)
			expectedToolSchemaJSON, _ := expectedTool.InputSchema.MarshalJSON()
			actualToolSchemaJSON, _ := actualTool.InputSchema.MarshalJSON()
			assert.Equal(t, expectedToolSchemaJSON, actualToolSchemaJSON)
			if tt.wantErr && err != nil {
				assert.Contains(t, err.Error(), "expected BOOLEAN schema, got")
			}
		})
	}
}

func TestBooleanConverter_ValidateArguments(t *testing.T) {
	converter := NewBooleanConverter()

	tests := []struct {
		name       string
		schemaInfo *utils.SchemaInfo
		args       map[string]any
		wantErr    bool
		errContain string // Substring to check in error message if wantErr is true
	}{
		{
			name:       "Valid arguments for BOOLEAN schema",
			schemaInfo: newBoolSchemaInfo("BOOLEAN"),
			args:       map[string]any{ParamName: true},
			wantErr:    false,
		},
		{
			name:       "Invalid schema type (STRING)",
			schemaInfo: newBoolSchemaInfo("STRING"),
			args:       map[string]any{ParamName: true},
			wantErr:    true,
			errContain: "expected BOOLEAN schema, got STRING",
		},
		{
			name:       "Missing payload argument",
			schemaInfo: newBoolSchemaInfo("BOOLEAN"),
			args:       map[string]any{},
			wantErr:    true,
			errContain: "missing required parameter: payload",
		},
		{
			name:       "Incorrect payload type (string instead of bool)",
			schemaInfo: newBoolSchemaInfo("BOOLEAN"),
			args:       map[string]any{ParamName: "true"},
			wantErr:    true,
			errContain: "parameter payload is not of type bool",
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

func TestBooleanConverter_SerializeMCPRequestToPulsarPayload(t *testing.T) {
	converter := NewBooleanConverter()

	tests := []struct {
		name       string
		schemaInfo *utils.SchemaInfo
		args       map[string]any
		want       []byte
		wantErr    bool
		errContain string
	}{
		{
			name:       "Serialize true for BOOLEAN schema",
			schemaInfo: newBoolSchemaInfo("BOOLEAN"),
			args:       map[string]any{ParamName: true},
			want:       []byte("true"),
			wantErr:    false,
		},
		{
			name:       "Serialize false for BOOLEAN schema",
			schemaInfo: newBoolSchemaInfo("BOOLEAN"),
			args:       map[string]any{ParamName: false},
			want:       []byte("false"),
			wantErr:    false,
		},
		{
			name:       "Validation error (e.g., missing payload)",
			schemaInfo: newBoolSchemaInfo("BOOLEAN"),
			args:       map[string]any{},
			want:       nil,
			wantErr:    true,
			errContain: "arguments validation failed", // Outer error message from SerializeMCPRequestToPulsarPayload
		},
		{
			name:       "Validation error (e.g., wrong schema type during ValidateArguments)",
			schemaInfo: newBoolSchemaInfo("STRING"), // Invalid schema type for this converter
			args:       map[string]any{ParamName: true},
			want:       nil,
			wantErr:    true,
			errContain: "arguments validation failed", // Outer error message from SerializeMCPRequestToPulsarPayload
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
