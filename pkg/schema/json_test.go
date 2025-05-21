package schema

import (
	"testing"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/stretchr/testify/assert"

	"github.com/mark3labs/mcp-go/mcp"
)

// newJSONSchemaInfo is a helper to create SchemaInfo for JSON type with a given AVRO schema string.
func newJSONSchemaInfo(avroSchemaString string) *utils.SchemaInfo {
	return &utils.SchemaInfo{
		Name:   "test-json-schema",
		Type:   "JSON",
		Schema: []byte(avroSchemaString),
	}
}

// TestNewJSONConverter tests the NewJSONConverter constructor.
func TestNewJSONConverter(t *testing.T) {
	converter := NewJSONConverter()
	assert.NotNil(t, converter, "NewJSONConverter should not return nil")
	// JSONConverter does not have a ParamName like simpler converters, it relies on Avro structure.
}

// TestJSONConverter_ToMCPToolInputSchemaProperties tests ToMCPToolInputSchemaProperties for JSONConverter.
func TestJSONConverter_ToMCPToolInputSchemaProperties(t *testing.T) {
	converter := NewJSONConverter()

	// Re-define or ensure accessibility of AVRO schema constants from avro_core_test.go
	// For this example, let's use a simple one. Ideally, these would be shared or accessible.
	const localSimpleRecordSchema = `{
		"type": "record",
		"name": "SimpleRecordForJSON",
		"fields": [
			{"name": "name", "type": "string"},
			{"name": "age", "type": "int"}
		]
	}`

	const invalidAvroSchema = `{"type": "invalid}`

	tests := []struct {
		name            string
		schemaInfo      *utils.SchemaInfo
		expectedOptions []mcp.ToolOption // Simplified expectation for brevity
		expectError     bool
		errorContains   string
	}{
		{
			name:       "Valid JSON schema (based on simple Avro record)",
			schemaInfo: newJSONSchemaInfo(localSimpleRecordSchema),
			expectedOptions: []mcp.ToolOption{ // This structure depends on processAvroSchemaStringToMCPToolInput
				mcp.WithString("name", mcp.Description("name (type: string)"), mcp.Required()),
				mcp.WithNumber("age", mcp.Description("age (type: int)"), mcp.Required()),
			},
			expectError: false,
		},
		{
			name:            "SchemaInfo type is not JSON",
			schemaInfo:      &utils.SchemaInfo{Type: "AVRO", Schema: []byte(localSimpleRecordSchema)},
			expectedOptions: nil,
			expectError:     true,
			errorContains:   "expected JSON schema, got AVRO",
		},
		{
			name:            "Invalid underlying Avro schema string",
			schemaInfo:      newJSONSchemaInfo(invalidAvroSchema),
			expectedOptions: nil,
			expectError:     true,
			errorContains:   "unknown type: {\"type\": \"invalid}", // Outer error from JSONConverter
		},
		{
			name:            "Underlying Avro schema is nil",
			schemaInfo:      &utils.SchemaInfo{Type: "JSON", Schema: nil},
			expectedOptions: nil,
			expectError:     true,
			errorContains:   "avro: unknown type: ",
		},
		{
			name:            "Underlying Avro schema is empty string",
			schemaInfo:      newJSONSchemaInfo(""),
			expectedOptions: nil,
			expectError:     true,
			errorContains:   "avro: unknown type: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := converter.ToMCPToolInputSchemaProperties(tt.schemaInfo)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				assert.Nil(t, opts) // Or check for empty slice if appropriate
			} else {
				assert.NoError(t, err)
				var expectedTool = mcp.NewTool("test", tt.expectedOptions...)
				var actualTool = mcp.NewTool("test", opts...)
				expectedToolSchemaJSON, _ := expectedTool.InputSchema.MarshalJSON()
				actualToolSchemaJSON, _ := actualTool.InputSchema.MarshalJSON()
				assert.Equal(t, expectedToolSchemaJSON, actualToolSchemaJSON, "ToolOptions mismatch")
			}
		})
	}
}

// TestJSONConverter_ValidateArguments tests ValidateArguments for JSONConverter.
func TestJSONConverter_ValidateArguments(t *testing.T) {
	converter := NewJSONConverter()

	// Re-use localSimpleRecordSchema and invalidAvroSchema from previous test or ensure they are accessible.
	const localSimpleRecordSchemaForValidation = `{
		"type": "record",
		"name": "SimpleRecordForJSONValidation",
		"fields": [
			{"name": "name", "type": "string"},
			{"name": "age", "type": "int"}
		]
	}`
	const invalidAvroSchemaForValidation = `{"type": "invalid}`

	tests := []struct {
		name          string
		schemaInfo    *utils.SchemaInfo
		args          map[string]any
		expectError   bool
		errorContains string
	}{
		{
			name:       "Valid arguments for JSON schema (simple record)",
			schemaInfo: newJSONSchemaInfo(localSimpleRecordSchemaForValidation),
			args: map[string]any{
				"name": "testUser",
				"age":  30,
			},
			expectError: false,
		},
		{
			name:          "SchemaInfo type is not JSON",
			schemaInfo:    &utils.SchemaInfo{Type: "AVRO", Schema: []byte(localSimpleRecordSchemaForValidation)},
			args:          map[string]any{"name": "testUser", "age": 30},
			expectError:   true,
			errorContains: "expected JSON schema for validation, got AVRO",
		},
		{
			name:          "Invalid underlying Avro schema string",
			schemaInfo:    newJSONSchemaInfo(invalidAvroSchemaForValidation),
			args:          map[string]any{"name": "testUser", "age": 30},
			expectError:   true,
			errorContains: "unknown type: {\"type\": \"invalid}", // Outer error
		},
		{
			name:          "Missing required field (age) in args for JSON schema",
			schemaInfo:    newJSONSchemaInfo(localSimpleRecordSchemaForValidation),
			args:          map[string]any{"name": "testUser"},
			expectError:   true,
			errorContains: "required field 'age' is missing and has no default value", // Error from validateArgumentsAgainstAvroSchemaString
		},
		{
			name:          "Wrong type for field (age as string) in args for JSON schema",
			schemaInfo:    newJSONSchemaInfo(localSimpleRecordSchemaForValidation),
			args:          map[string]any{"name": "testUser", "age": "thirty"},
			expectError:   true,
			errorContains: "field 'age': expected int, got string (value: thirty)", // Error from validateArgumentsAgainstAvroSchemaString
		},
		{
			name:          "Nil arguments map",
			schemaInfo:    newJSONSchemaInfo(localSimpleRecordSchemaForValidation),
			args:          nil, // validateArgumentsAgainstAvroSchemaString treats nil as empty map
			expectError:   true,
			errorContains: "required field 'name' is missing and has no default value",
		},
		{
			name:          "Underlying Avro schema is nil",
			schemaInfo:    &utils.SchemaInfo{Type: "JSON", Schema: nil},
			args:          map[string]any{"name": "testUser", "age": 30},
			expectError:   true,
			errorContains: "failed to parse AVRO schema for validation: avro: unknown type: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := converter.ValidateArguments(tt.args, tt.schemaInfo)
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestJSONConverter_SerializeMCPRequestToPulsarPayload tests SerializeMCPRequestToPulsarPayload for JSONConverter.
func TestJSONConverter_SerializeMCPRequestToPulsarPayload(t *testing.T) {
	converter := NewJSONConverter()

	const localSimpleRecordSchemaForSerialization = `{
		"type": "record",
		"name": "SimpleRecordForJSONSerialization",
		"fields": [
			{"name": "name", "type": "string"},
			{"name": "age", "type": "int"},
			{"name": "city", "type": ["null", "string"], "default": null}
		]
	}`

	tests := []struct {
		name          string
		schemaInfo    *utils.SchemaInfo
		args          map[string]any
		expectedJSON  string // Expected JSON string output
		expectError   bool
		errorContains string
	}{
		{
			name:       "Valid arguments, serialize to JSON",
			schemaInfo: newJSONSchemaInfo(localSimpleRecordSchemaForSerialization),
			args: map[string]any{
				"name": "Alice",
				"age":  30,
				"city": "New York",
			},
			// Note: JSON marshaling of maps doesn't guarantee key order.
			// We will unmarshal and compare maps for robust checking if direct string comparison is flaky.
			// For simplicity here, we assume a common order or use a more robust comparison later if needed.
			expectedJSON: `{"age":30,"city":"New York","name":"Alice"}`,
			expectError:  false,
		},
		{
			name:       "Valid arguments with optional field null, serialize to JSON",
			schemaInfo: newJSONSchemaInfo(localSimpleRecordSchemaForSerialization),
			args: map[string]any{
				"name": "Bob",
				"age":  40,
				"city": nil, // Explicit null for optional field
			},
			expectedJSON: `{"age":40,"city":null,"name":"Bob"}`,
			expectError:  false,
		},
		{
			name:       "Valid arguments with optional field omitted, serialize to JSON",
			schemaInfo: newJSONSchemaInfo(localSimpleRecordSchemaForSerialization),
			args: map[string]any{
				"name": "Charlie",
				"age":  35,
				// city is omitted, should be treated as null by Avro logic but json.Marshal will omit it if not in map
			},
			// If Avro layer adds default null to args map before JSON marshal, then `"city":null` would be here.
			// JSONConverter.SerializeMCPRequestToPulsarPayload directly marshals the provided args map.
			expectedJSON: `{"age":35,"name":"Charlie"}`,
			expectError:  false,
		},
		{
			name:       "Validation error (missing required field name)",
			schemaInfo: newJSONSchemaInfo(localSimpleRecordSchemaForSerialization),
			args: map[string]any{
				"age": 25,
			},
			expectedJSON:  "",
			expectError:   true,
			errorContains: "arguments validation failed", // Outer error from SerializeMCPRequestToPulsarPayload
		},
		{
			name:       "SchemaInfo type is not JSON",
			schemaInfo: &utils.SchemaInfo{Type: "AVRO", Schema: []byte(localSimpleRecordSchemaForSerialization)},
			args: map[string]any{
				"name": "David",
				"age":  28,
			},
			expectedJSON:  "",
			expectError:   true,
			errorContains: "expected JSON schema for validation, got AVRO",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := converter.SerializeMCPRequestToPulsarPayload(tt.args, tt.schemaInfo)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				assert.Nil(t, payload)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, payload)
				// For robust JSON comparison, unmarshal both expected and actual to maps and compare maps,
				// or use a library that does canonical JSON comparison.
				// For now, direct string comparison of compact JSON.
				assert.JSONEq(t, tt.expectedJSON, string(payload), "Serialized JSON does not match expected JSON")
			}
		})
	}
}
