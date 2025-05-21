package schema

import (
	"testing"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/stretchr/testify/assert"

	"github.com/mark3labs/mcp-go/mcp"
)

// complexRecordSchemaString is used by TestAvroConverter_ValidateArguments
// and is not defined in avro_core_test.go, so we define it here.
const complexRecordSchemaString = `{
    "type": "record",
    "name": "ComplexRecord",
    "fields": [
        {"name": "fieldString", "type": "string"},
        {"name": "fieldInt", "type": "int"},
        {"name": "fieldLong", "type": "long"},
        {"name": "fieldDouble", "type": "double"},
        {"name": "fieldFloat", "type": "float"},
        {"name": "fieldBool", "type": "boolean"},
        {"name": "fieldBytes", "type": "bytes"},
        {"name": "fieldFixed", "type": {"type": "fixed", "name": "MyFixed", "size": 8}},
        {"name": "fieldEnum", "type": {"type": "enum", "name": "MyEnum", "symbols": ["A", "B", "C"]}},
        {"name": "fieldArray", "type": {"type": "array", "items": "int"}},
        {"name": "fieldMap", "type": {"type": "map", "values": "long"}},
        {"name": "fieldRecord", "type": {
            "type": "record", "name": "SubRecord", "fields": [
                {"name": "subField", "type": "string"}
            ]}
        },
        {"name": "fieldUnion", "type": ["null", "int", "string"]}
    ]
}`

// newAvroSchemaInfo is a helper to create SchemaInfo for AVRO type with a given AVRO schema string.
func newAvroSchemaInfo(avroSchemaString string) *utils.SchemaInfo {
	return &utils.SchemaInfo{
		Name:   "test-avro-schema",
		Type:   "AVRO", // Pulsar schema type is string
		Schema: []byte(avroSchemaString),
	}
}

// TestNewAvroConverter tests the NewAvroConverter constructor.
func TestNewAvroConverter(t *testing.T) {
	converter := NewAvroConverter()
	assert.NotNil(t, converter, "NewAvroConverter should not return nil")
	// AvroConverter also relies on Avro structure primarily, similar to JSONConverter.
}

// TestAvroConverter_ToMCPToolInputSchemaProperties tests ToMCPToolInputSchemaProperties for AvroConverter.
func TestAvroConverter_ToMCPToolInputSchemaProperties(t *testing.T) {
	converter := NewAvroConverter()

	const localSimpleRecordSchemaForAvro = `{
		"type": "record",
		"name": "SimpleRecordForAvro",
		"fields": [
			{"name": "id", "type": "long"},
			{"name": "data", "type": "string"}
		]
	}`

	const invalidAvroSchemaForAvro = `{"type": "invalidAvro}`

	tests := []struct {
		name            string
		schemaInfo      *utils.SchemaInfo
		expectedOptions []mcp.ToolOption
		expectError     bool
		errorContains   string
	}{
		{
			name:       "Valid AVRO schema",
			schemaInfo: newAvroSchemaInfo(localSimpleRecordSchemaForAvro),
			expectedOptions: []mcp.ToolOption{
				mcp.WithNumber("id", mcp.Description("id (type: long)"), mcp.Required()),
				mcp.WithString("data", mcp.Description("data (type: string)"), mcp.Required()),
			},
			expectError: false,
		},
		{
			name:            "SchemaInfo type is not AVRO (e.g., JSON)",
			schemaInfo:      &utils.SchemaInfo{Type: "JSON", Schema: []byte(localSimpleRecordSchemaForAvro)},
			expectedOptions: nil,
			expectError:     true,
			errorContains:   "expected AVRO schema, got JSON",
		},
		{
			name:            "Invalid underlying Avro schema string",
			schemaInfo:      newAvroSchemaInfo(invalidAvroSchemaForAvro),
			expectedOptions: nil,
			expectError:     true,
			errorContains:   "unknown type: {\"type\": \"invalidAvro}",
		},
		{
			name:            "Underlying Avro schema is nil",
			schemaInfo:      &utils.SchemaInfo{Type: "AVRO", Schema: nil},
			expectedOptions: nil,
			expectError:     true,
			errorContains:   "unknown type: ",
		},
		{
			name:            "Underlying Avro schema is empty string",
			schemaInfo:      newAvroSchemaInfo(""),
			expectedOptions: nil,
			expectError:     true,
			errorContains:   "unknown type: ",
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
				assert.Nil(t, opts)
			} else {
				assert.NoError(t, err)
				var expectedTool, actualTool mcp.Tool
				expectedTool = mcp.NewTool("test", tt.expectedOptions...)
				actualTool = mcp.NewTool("test", opts...)
				expectedToolSchemaJSON, _ := expectedTool.InputSchema.MarshalJSON()
				actualToolSchemaJSON, _ := actualTool.InputSchema.MarshalJSON()
				assert.Equal(t, expectedToolSchemaJSON, actualToolSchemaJSON, "Tool mismatch")
			}
		})
	}
}

func TestAvroConverter_ValidateArguments(t *testing.T) {
	tests := []struct {
		name       string
		schemaInfo utils.SchemaInfo
		args       map[string]interface{}
		wantErr    bool
	}{
		{
			name: "valid arguments for simple record",
			schemaInfo: utils.SchemaInfo{
				Name:   "SimpleAvro",
				Schema: []byte(simpleRecordSchema),
				Type:   "AVRO",
			},
			args: map[string]interface{}{
				"id":   int64(123),
				"name": "TestName",
			},
			wantErr: false,
		},
		{
			name: "invalid type for field in simple record",
			schemaInfo: utils.SchemaInfo{
				Name:   "SimpleAvroInvalidType",
				Schema: []byte(simpleRecordSchema),
				Type:   "AVRO",
			},
			args: map[string]interface{}{
				"id":   "not_a_long",
				"name": "TestName",
			},
			wantErr: true,
		},
		{
			name: "missing required field in simple record",
			schemaInfo: utils.SchemaInfo{
				Name:   "SimpleAvroMissingField",
				Schema: []byte(simpleRecordSchema),
				Type:   "AVRO",
			},
			args: map[string]interface{}{
				"name": "TestName",
			},
			wantErr: true,
		},
		{
			name: "valid arguments for complex record",
			schemaInfo: utils.SchemaInfo{
				Name:   "ComplexAvroValid",
				Schema: []byte(complexRecordSchemaString),
				Type:   "AVRO",
			},
			args: map[string]interface{}{
				"fieldString": "test string",
				"fieldInt":    int32(123),
				"fieldLong":   int64(456),
				"fieldDouble": float64(12.34),
				"fieldFloat":  float32(56.78),
				"fieldBool":   true,
				"fieldBytes":  []byte("test bytes"),
				"fieldFixed":  uint64(0x1234567890123456),
				"fieldEnum":   "A",
				"fieldArray":  []interface{}{int32(1), int32(2)},
				"fieldMap":    map[string]interface{}{"key1": int64(100)},
				"fieldRecord": map[string]interface{}{"subField": "sub value"},
				"fieldUnion":  int32(99),
			},
			wantErr: false,
		},
		{
			name: "invalid enum value for complex record",
			schemaInfo: utils.SchemaInfo{
				Name:   "ComplexAvroInvalidEnum",
				Schema: []byte(complexRecordSchemaString),
				Type:   "AVRO",
			},
			args: map[string]interface{}{
				"fieldString": "test string",
				"fieldInt":    int32(123),
				"fieldLong":   int64(456),
				"fieldDouble": float64(12.34),
				"fieldFloat":  float32(56.78),
				"fieldBool":   true,
				"fieldBytes":  []byte("test bytes"),
				"fieldFixed":  uint64(0x1234567890123456),
				"fieldEnum":   "X", // Invalid enum
				"fieldArray":  []interface{}{int32(1), int32(2)},
				"fieldMap":    map[string]interface{}{"key1": int64(100)},
				"fieldRecord": map[string]interface{}{"subField": "sub value"},
				"fieldUnion":  int32(99),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conv := NewAvroConverter()

			err := conv.ValidateArguments(tt.args, &tt.schemaInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("AvroConverter.ValidateArguments() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAvroConverter_SerializeMCPRequestToPulsarPayload(t *testing.T) {
	// Assumes simpleRecordSchema and complexRecordSchemaString are available
	// and getExpectedAvroBinary is accessible from avro_core_test.go (same package)

	tests := []struct {
		name            string
		schemaInfo      utils.SchemaInfo
		args            map[string]interface{}
		expectedPayload []byte // Can be nil if wantErr is true
		wantErr         bool
		assertPayload   bool // True if we want to compare payload, false if only error matters
		errorContains   string
	}{
		{
			name: "valid serialization for simple record",
			schemaInfo: utils.SchemaInfo{
				Name:   "SimpleAvroSerialize",
				Schema: []byte(simpleRecordSchema), // Uses simpleRecordSchema from avro_core_test.go
				Type:   "AVRO",
			},
			args: map[string]interface{}{
				"id":   int64(42), // Note: field names in avro_core_test.go's simpleRecordSchema are lowercase
				"name": "Test Payload",
			},
			assertPayload: true,
			wantErr:       false,
		},
		{
			name: "serialization with mismatched argument type for simple record",
			schemaInfo: utils.SchemaInfo{
				Name:   "SimpleAvroSerializeInvalidArg",
				Schema: []byte(simpleRecordSchema),
				Type:   "AVRO",
			},
			args: map[string]interface{}{
				"id":   "not_a_long", // Invalid type
				"name": "Test Payload",
			},
			assertPayload: false,
			wantErr:       true,
			errorContains: "arguments validation failed",
		},
		{
			name: "valid serialization for complex record",
			schemaInfo: utils.SchemaInfo{
				Name:   "ComplexAvroSerialize",
				Schema: []byte(complexRecordSchemaString), // Uses complexRecordSchemaString defined in avro_test.go
				Type:   "AVRO",
			},
			args: map[string]interface{}{
				"fieldString": "hello avro",
				"fieldInt":    int32(101),
				"fieldLong":   int64(202),
				"fieldDouble": float64(30.3),
				"fieldFloat":  float32(4.04),
				"fieldBool":   true,
				"fieldBytes":  []byte("avro bytes"),
				"fieldFixed":  uint64(0x1234567890123456), // Must be 16 bytes, will be corrected in loop
				"fieldEnum":   "B",
				"fieldArray":  []interface{}{int32(11), int32(22)},
				"fieldMap":    map[string]interface{}{"mKey": int64(333)},
				"fieldRecord": map[string]interface{}{"subField": "sub data"},
				"fieldUnion":  "union string val",
			},
			assertPayload: true, // We'll generate expected payload for this
			wantErr:       false,
		},
		{
			name: "serialization with invalid schema type",
			schemaInfo: utils.SchemaInfo{
				Name:   "InvalidSchemaTypeSerialize",
				Schema: []byte(simpleRecordSchema),
				Type:   "JSON", // Invalid type for AvroConverter
			},
			args: map[string]interface{}{
				"id":   int64(1),
				"name": "Dummy",
			},
			assertPayload: false,
			wantErr:       true,
			errorContains: "arguments validation failed", // ValidateArguments checks schema type first
		},
	}

	for i := range tests {
		tt := &tests[i] // Use pointer to allow modification for expectedPayload

		// Pre-calculate expected payload for valid cases
		if tt.assertPayload && !tt.wantErr {
			var schemaToUse string
			var argsToMarshal map[string]interface{}

			if tt.schemaInfo.Name == "SimpleAvroSerialize" {
				schemaToUse = simpleRecordSchema
				argsToMarshal = tt.args
			} else if tt.schemaInfo.Name == "ComplexAvroSerialize" {
				schemaToUse = complexRecordSchemaString
				complexArgsCopy := make(map[string]interface{})
				for k, v := range tt.args {
					complexArgsCopy[k] = v
				}
				complexArgsCopy["fieldFixed"] = uint64(0x0123456789abcdef) // 16 bytes
				argsToMarshal = complexArgsCopy
				tt.args = complexArgsCopy // Update tt.args with corrected fixed field
			}

			if schemaToUse != "" {
				tt.expectedPayload = getExpectedAvroBinary(t, schemaToUse, argsToMarshal)
			} else {
				t.Fatalf("Test setup error for %s: schemaToUse not set for payload generation", tt.name)
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			conv := NewAvroConverter()
			payload, err := conv.SerializeMCPRequestToPulsarPayload(tt.args, &tt.schemaInfo)

			if tt.wantErr {
				assert.Error(t, err, "Expected an error for test case: %s", tt.name)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains, "Error message mismatch for test case: %s", tt.name)
				}
			} else {
				assert.NoError(t, err, "Did not expect an error for test case: %s", tt.name)
				if tt.assertPayload {
					assert.Equal(t, tt.expectedPayload, payload, "Payload mismatch for test case: %s", tt.name)
				}
			}
		})
	}
}
