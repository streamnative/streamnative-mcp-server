package schema

import (
	"testing"

	"github.com/hamba/avro/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mark3labs/mcp-go/mcp"
)

// AVRO Schema Definitions for Testing

const simpleRecordSchema = `{
	"type": "record",
	"name": "SimpleRecord",
	"fields": [
		{"name": "id", "type": "long"},
		{"name": "name", "type": "string"}
	]
}`

const schemaWithAllPrimitives = `{
	"type": "record",
	"name": "AllPrimitives",
	"fields": [
		{"name": "boolField", "type": "boolean"},
		{"name": "intField", "type": "int"},
		{"name": "longField", "type": "long"},
		{"name": "floatField", "type": "float"},
		{"name": "doubleField", "type": "double"},
		{"name": "bytesField", "type": "bytes"},
		{"name": "stringField", "type": "string"}
	]
}`

const schemaWithOptionalField = `{
	"type": "record",
	"name": "OptionalFieldRecord",
	"fields": [
		{"name": "requiredField", "type": "string"},
		{"name": "optionalField", "type": ["null", "string"], "default": null}
	]
}`

const schemaWithDefaultValue = `{
	"type": "record",
	"name": "DefaultValueRecord",
	"fields": [
		{"name": "name", "type": "string", "default": "DefaultName"},
		{"name": "age", "type": "int", "default": 30}
	]
}`

const nestedRecordSchema = `{
	"type": "record",
	"name": "OuterRecord",
	"fields": [
		{"name": "id", "type": "string"},
		{
			"name": "inner",
			"type": {
				"type": "record",
				"name": "InnerRecord",
				"fields": [
					{"name": "value", "type": "int"},
					{"name": "description", "type": ["null", "string"], "default": null}
				]
			}
		}
	]
}`

const arraySchemaPrimitive = `{
	"type": "record",
	"name": "ArrayPrimitiveRecord",
	"fields": [
		{"name": "stringArray", "type": {"type": "array", "items": "string"}},
		{"name": "optionalIntArray", "type": ["null", {"type": "array", "items": "int"}], "default": null}
	]
}`

const arraySchemaRecord = `{
	"type": "record",
	"name": "ArrayRecordContainer",
	"fields": [
		{
			"name": "records", 
			"type": {"type": "array", "items": {
				"type": "record",
				"name": "ContainedRecord",
				"fields": [
					{"name": "key", "type": "string"},
					{"name": "val", "type": "long"}
				]
			}}
		}
	]
}`

const mapSchemaPrimitive = `{
	"type": "record",
	"name": "MapPrimitiveRecord",
	"fields": [
		{"name": "stringMap", "type": {"type": "map", "values": "string"}},
		{"name": "optionalIntMap", "type": ["null", {"type": "map", "values": "int"}], "default": null}
	]
}`

const mapSchemaRecord = `{
	"type": "record",
	"name": "MapRecordContainer",
	"fields": [
		{
			"name": "recordsMap", 
			"type": {"type": "map", "values": {
				"type": "record",
				"name": "MappedRecord",
				"fields": [
					{"name": "id", "type": "int"},
					{"name": "status", "type": "string"}
				]
			}}
		}
	]
}`

const enumSchema = `{
	"type": "record",
	"name": "EnumRecord",
	"fields": [
		{
			"name": "suit", 
			"type": { "type": "enum", "name": "Suit", "symbols" : ["SPADES", "HEARTS", "DIAMONDS", "CLUBS"] }
		}
	]
}`

const unionSchemaSimple = `{
	"type": "record",
	"name": "UnionRecord",
	"fields": [
		{"name": "stringOrInt", "type": ["string", "int"]},
		{"name": "nullableStringOrInt", "type": ["null", "string", "int"], "default": null}
	]
}`

// Helper function to get expected AVRO binary for test comparisons
func getExpectedAvroBinary(t *testing.T, schemaStr string, data map[string]any) []byte {
	avroSchema, err := avro.Parse(schemaStr)
	require.NoError(t, err, "Failed to parse AVRO schema in helper")

	// The hamba/avro/v2 Marshal function can take a map[string]any directly
	// if its structure is compatible with the schema.
	bin, err := avro.Marshal(avroSchema, data)
	require.NoError(t, err, "Failed to marshal to AVRO binary in helper")
	return bin
}

func TestProcessAvroSchemaStringToMCPToolInput(t *testing.T) {
	tests := []struct {
		name             string
		schemaStr        string
		expectedOptions  []mcp.ToolOption
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:      "Simple Record",
			schemaStr: simpleRecordSchema,
			expectedOptions: []mcp.ToolOption{
				mcp.WithNumber("id", mcp.Description("id (type: long)"), mcp.Required()),
				mcp.WithString("name", mcp.Description("name (type: string)"), mcp.Required()),
			},
			expectError: false,
		},
		{
			name:      "Schema With All Primitives",
			schemaStr: schemaWithAllPrimitives,
			expectedOptions: []mcp.ToolOption{
				mcp.WithBoolean("boolField", mcp.Description("boolField (type: boolean)"), mcp.Required()),
				mcp.WithNumber("intField", mcp.Description("intField (type: int)"), mcp.Required()),
				mcp.WithNumber("longField", mcp.Description("longField (type: long)"), mcp.Required()),
				mcp.WithNumber("floatField", mcp.Description("floatField (type: float)"), mcp.Required()),
				mcp.WithNumber("doubleField", mcp.Description("doubleField (type: double)"), mcp.Required()),
				mcp.WithString("bytesField", mcp.Description("bytesField (type: bytes)"), mcp.Required()),
				mcp.WithString("stringField", mcp.Description("stringField (type: string)"), mcp.Required()),
			},
			expectError: false,
		},
		{
			name:             "Invalid AVRO schema string",
			schemaStr:        `{"type": "invalid"`,
			expectedOptions:  nil,
			expectError:      true,
			expectedErrorMsg: "failed to parse AVRO schema",
		},
		{
			name:             "Top-level non-record (string)",
			schemaStr:        `"string"`,
			expectedOptions:  nil,
			expectError:      true,
			expectedErrorMsg: "expected AVRO record schema at the top level, got *avro.PrimitiveSchema",
		},
		{
			name:      "Schema With Optional Field (string)",
			schemaStr: schemaWithOptionalField,
			expectedOptions: []mcp.ToolOption{
				mcp.WithString("requiredField", mcp.Description("requiredField (type: string)"), mcp.Required()),
				mcp.WithString("optionalField", mcp.Description("optionalField (type: [null,string])")),
			},
			expectError: false,
		},
		{
			name:      "Schema With Default Value (string and int)",
			schemaStr: schemaWithDefaultValue,
			expectedOptions: []mcp.ToolOption{
				mcp.WithString("name", mcp.Description("name (type: string)"), mcp.Required(), mcp.DefaultString("DefaultName")),
				mcp.WithNumber("age", mcp.Description("age (type: int)"), mcp.Required(), mcp.DefaultNumber(30)),
			},
			expectError: false,
		},
		{
			name:      "Nested Record",
			schemaStr: nestedRecordSchema,
			expectedOptions: []mcp.ToolOption{
				mcp.WithString("id", mcp.Description("id (type: string)"), mcp.Required()),
				mcp.WithObject("inner",
					mcp.Description("inner (type: {name:InnerRecord,type:record,fields:[{name:value,type:int},{name:description,type:[null,string]}]})"),
					mcp.Required(),
					mcp.Properties(map[string]any{
						"value": map[string]any{
							"type":        "number",
							"description": "value (type: int)",
						},
						"description": map[string]any{
							"type":        "string",
							"description": "description (type: [null,string]) (nullable)",
						},
					}),
				),
			},
			expectError: false,
		},
		{
			name:      "Array of Primitives (stringArray, optionalIntArray)",
			schemaStr: arraySchemaPrimitive,
			expectedOptions: []mcp.ToolOption{
				mcp.WithArray("stringArray",
					mcp.Description("stringArray (type: {type:array,items:string})"),
					mcp.Required(),
					mcp.Items(map[string]any{
						"type":        "string",
						"description": "Array items",
					}),
				),
				mcp.WithArray("optionalIntArray",
					mcp.Description("optionalIntArray (type: [null,{type:array,items:int}])"),
					mcp.Items(map[string]any{
						"type":        "number",
						"description": "Array items",
					}),
				),
			},
			expectError: false,
		},
		{
			name:      "Array of Records",
			schemaStr: arraySchemaRecord,
			expectedOptions: []mcp.ToolOption{
				mcp.WithArray("records",
					mcp.Description("records (type: {type:array,items:{name:ContainedRecord,type:record,fields:[{name:key,type:string},{name:val,type:long}]}})"),
					mcp.Required(),
					mcp.Items(map[string]any{
						"type":        "object",
						"description": "Array items",
						"properties": map[string]any{
							"key": map[string]any{
								"type":        "string",
								"description": "key (type: string)",
							},
							"val": map[string]any{
								"type":        "number",
								"description": "val (type: long)",
							},
						},
					}),
				),
			},
			expectError: false,
		},
		{
			name:      "Map of Primitives (stringMap, optionalIntMap)",
			schemaStr: mapSchemaPrimitive,
			expectedOptions: []mcp.ToolOption{
				mcp.WithObject("stringMap", // Avro map becomes MCP object
					mcp.Description("stringMap (type: {type:map,values:string})"),
					mcp.Required(),
					mcp.Properties(map[string]any{ // Based on avroFieldToMcpOption logic for map
						"*": map[string]any{
							"type":        "string",
							"description": "Map values",
						},
					}),
				),
				mcp.WithObject("optionalIntMap",
					mcp.Description("optionalIntMap (type: [null,{type:map,values:int}])"),
					// Not required due to ["null", map] union
					mcp.Properties(map[string]any{
						"*": map[string]any{
							"type":        "number",
							"description": "Map values",
						},
					}),
					// Avro default: null for the union handled by not being required.
				),
			},
			expectError: false,
		},
		{
			name:      "Map of Records",
			schemaStr: mapSchemaRecord,
			expectedOptions: []mcp.ToolOption{
				mcp.WithObject("recordsMap",
					mcp.Description("recordsMap (type: {type:map,values:{name:MappedRecord,type:record,fields:[{name:id,type:int},{name:status,type:string}]}})"),
					mcp.Required(),
					mcp.Properties(map[string]any{
						"*": map[string]any{
							"type":        "object",
							"description": "Map values",
							"properties": map[string]any{
								"id": map[string]any{
									"type":        "number",
									"description": "id (type: int)",
								},
								"status": map[string]any{
									"type":        "string",
									"description": "status (type: string)",
								},
							},
						},
					}),
				),
			},
			expectError: false,
		},
		{
			name:      "Enum Field",
			schemaStr: enumSchema,
			expectedOptions: []mcp.ToolOption{
				mcp.WithString("suit",
					mcp.Description("suit (type: {name:Suit,type:enum,symbols:[SPADES,HEARTS,DIAMONDS,CLUBS]})"),
					mcp.Required(),
					mcp.Enum("SPADES", "HEARTS", "DIAMONDS", "CLUBS"),
				),
			},
			expectError: false,
		},
		{
			name:      "Simple Union Field (stringOrInt)",
			schemaStr: unionSchemaSimple,
			expectedOptions: []mcp.ToolOption{
				// Based on avroFieldToMcpOption, a complex union ["string", "int"] becomes mcp.WithString
				// with a description indicating it's a union. It's marked as required by default.
				mcp.WithString("stringOrInt",
					mcp.Description("stringOrInt (type: [string,int]) (union type: [string,int])"),
					mcp.Required(),
				),
				// For ["null", "string", "int"], it's also a complex union but not required.
				mcp.WithString("nullableStringOrInt",
					mcp.Description("nullableStringOrInt (type: [null,string,int]) (union type: [null,string,int])"),
					// Not mcp.Required() due to presence of "null" in union.
					// Default is Avro null, so no mcp.DefaultString is added.
				),
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := processAvroSchemaStringToMCPToolInput(tt.schemaStr)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedErrorMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrorMsg)
				}
				assert.Nil(t, opts)
			} else {
				assert.NoError(t, err)
				require.Equal(t, len(tt.expectedOptions), len(opts), "Number of options should match")
				var actualTool, expectedTool mcp.Tool
				actualTool = mcp.NewTool("test", opts...)
				expectedTool = mcp.NewTool("test", tt.expectedOptions...)
				actualToolInputSchemaJSON, _ := actualTool.InputSchema.MarshalJSON()
				expectedToolInputSchemaJSON, _ := expectedTool.InputSchema.MarshalJSON()
				assert.Equal(t, string(expectedToolInputSchemaJSON), string(actualToolInputSchemaJSON), "ToolOption did not produce the same property configuration. Expected: %+v, Got: %+v", expectedTool, actualTool)
			}
		})
	}
}

// TestValidateArgumentsAgainstAvroSchemaString tests the validateArgumentsAgainstAvroSchemaString function.
func TestValidateArgumentsAgainstAvroSchemaString(t *testing.T) {
	tests := []struct {
		name      string
		schemaStr string
		args      map[string]any
		wantErr   bool
		errText   string // Optional: specific error text to check for
	}{
		{
			name:      "Valid: Simple Record - all fields present",
			schemaStr: simpleRecordSchema,
			args: map[string]any{
				"id":   int64(123),
				"name": "test name",
			},
			wantErr: false,
		},
		{
			name:      "Valid: Schema With All Primitives - all fields present",
			schemaStr: schemaWithAllPrimitives,
			args: map[string]any{
				"boolField":   true,
				"intField":    int32(100),
				"longField":   int64(2000),
				"floatField":  float32(3.14),
				"doubleField": float64(6.28),
				"bytesField":  []byte("bytesdata"),
				"stringField": "stringdata",
			},
			wantErr: false,
		},
		{
			name:      "Invalid: Simple Record - missing required field 'name'",
			schemaStr: simpleRecordSchema,
			args: map[string]any{
				"id": int64(123),
			},
			wantErr: true,
			errText: "required field 'name' is missing and has no default value",
		},
		{
			name:      "Invalid: Simple Record - wrong type for 'id' (string instead of long)",
			schemaStr: simpleRecordSchema,
			args: map[string]any{
				"id":   "not-a-long",
				"name": "A Name",
			},
			wantErr: true,
			errText: "field 'id': expected long, got string",
		},
		{
			name:      "Invalid: Simple Record - extra field 'extra'",
			schemaStr: simpleRecordSchema,
			args: map[string]any{
				"id":    int64(123),
				"name":  "A Name",
				"extra": "value",
			},
			wantErr: true,
			errText: "unknown field 'extra' provided in arguments",
		},
		{
			name:      "Valid: Optional Field - present",
			schemaStr: schemaWithOptionalField,
			args: map[string]any{
				"requiredField": "req",
				"optionalField": "opt",
			},
			wantErr: false,
		},
		{
			name:      "Valid: Optional Field - absent (should use default null)",
			schemaStr: schemaWithOptionalField,
			args: map[string]any{
				"requiredField": "req",
			},
			wantErr: false, // AVRO validation itself passes if default is null and field is omitted
		},
		{
			name:      "Valid: Default Value - fields omitted",
			schemaStr: schemaWithDefaultValue,
			args:      map[string]any{},
			wantErr:   false, // Default values are used
		},
		{
			name:      "Valid: Default Value - fields provided",
			schemaStr: schemaWithDefaultValue,
			args: map[string]any{
				"name": "ProvidedName",
				"age":  int32(40),
			},
			wantErr: false,
		},
		{
			name:      "Valid: Nested Record",
			schemaStr: nestedRecordSchema,
			args: map[string]any{
				"id": "outerID",
				"inner": map[string]any{
					"value":       int32(101),
					"description": "inner desc",
				},
			},
			wantErr: false,
		},
		{
			name:      "Invalid: Nested Record - missing field in inner record",
			schemaStr: nestedRecordSchema,
			args: map[string]any{
				"id":    "outerID",
				"inner": map[string]any{"name": "Inner Name"}, // Missing inner.value
			},
			wantErr: true,
			errText: "field 'inner.value' is required but missing",
		},
		{
			name:      "Valid: Array of Primitives",
			schemaStr: arraySchemaPrimitive,
			args: map[string]any{
				"stringArray":      []any{"a", "b", "c"},
				"optionalIntArray": []any{int32(1), int32(2)},
			},
			wantErr: false,
		},
		{
			name:      "Invalid: Array of Primitives - wrong item type",
			schemaStr: arraySchemaPrimitive,
			args: map[string]any{
				"stringArray": []any{"hello", 1, "world"}, // int is not string
			},
			wantErr: true,
			errText: "field 'stringArray[1]': expected string, got int",
		},
		{
			name:      "Valid: Array of Records",
			schemaStr: arraySchemaRecord,
			args: map[string]any{
				"records": []any{
					map[string]any{"key": "k1", "val": int64(1)},
					map[string]any{"key": "k2", "val": int64(2)},
				},
			},
			wantErr: false,
		},
		{
			name:      "Valid: Map of Primitives",
			schemaStr: mapSchemaPrimitive,
			args: map[string]any{
				"stringMap":      map[string]any{"key1": "val1", "key2": "val2"},
				"optionalIntMap": map[string]any{"opt1": int32(10)},
			},
			wantErr: false,
		},
		{
			name:      "Invalid: Map of Primitives - wrong value type",
			schemaStr: mapSchemaPrimitive,
			args: map[string]any{
				"stringMap": map[string]any{"key1": 123, "key2": "val2"}, // 123 is not string
			},
			wantErr: true,
			errText: "field 'stringMap.key1': expected string, got int",
		},
		{
			name:      "Valid: Map of Records",
			schemaStr: mapSchemaRecord,
			args: map[string]any{
				"recordsMap": map[string]any{
					"recA": map[string]any{"id": int32(1), "status": "active"},
					"recB": map[string]any{"id": int32(2), "status": "inactive"},
				},
			},
			wantErr: false,
		},
		{
			name:      "Valid: Enum",
			schemaStr: enumSchema,
			args: map[string]any{
				"suit": "SPADES",
			},
			wantErr: false,
		},
		{
			name:      "Invalid: Enum - invalid symbol",
			schemaStr: enumSchema,
			args: map[string]any{
				"suit": "INVALID_SUIT",
			},
			wantErr: true,
			errText: "value 'INVALID_SUIT' is not a valid symbol for enum Suit",
		},
		{
			name:      "Valid: Union (stringOrInt) - string",
			schemaStr: unionSchemaSimple,
			args: map[string]any{
				"stringOrInt":         "hello",
				"nullableStringOrInt": "world",
			},
			wantErr: false,
		},
		{
			name:      "Valid: Union (stringOrInt) - int",
			schemaStr: unionSchemaSimple,
			args: map[string]any{
				"stringOrInt":         int32(123),
				"nullableStringOrInt": int32(456),
			},
			wantErr: false,
		},
		{
			name:      "Invalid: Union (stringOrInt) - boolean (not in union)",
			schemaStr: unionSchemaSimple,
			args: map[string]any{
				"stringOrInt": true,
			},
			wantErr: true,
			errText: "does not match any type in union schema",
		},
		{
			name:      "Invalid: Schema string is empty",
			schemaStr: "",
			args:      map[string]any{"foo": "bar"},
			wantErr:   true,
			errText:   "failed to parse AVRO schema for validation: avro: unknown type: ",
		},
		{
			name:      "Invalid: Schema string is not valid AVRO json",
			schemaStr: "{invalid json",
			args:      map[string]any{"foo": "bar"},
			wantErr:   true,
			errText:   "failed to parse AVRO schema",
		},
		{
			name:      "Valid: schemaWithAllPrimitives - bytes field with string input (should be accepted as per current code)",
			schemaStr: schemaWithAllPrimitives,
			args: map[string]any{
				"boolField":   true,
				"intField":    int32(100),
				"longField":   int64(2000),
				"floatField":  float32(3.14),
				"doubleField": float64(6.28),
				"bytesField":  "stringtobytes", // This is the key part for this test
				"stringField": "stringdata",
			},
			wantErr: false, // Current validateValueAgainstAvroType for bytes accepts string
		},
		{
			name:      "Invalid: schemaWithAllPrimitives - bytes field with int input",
			schemaStr: schemaWithAllPrimitives,
			args: map[string]any{
				"boolField":   true,
				"intField":    int32(100),
				"longField":   int64(2000),
				"floatField":  float32(3.14),
				"doubleField": float64(6.28),
				"bytesField":  123, // int is not convertible to bytes in this path
				"stringField": "stringdata",
			},
			wantErr: true,
			errText: "field 'bytesField': expected string or []byte for bytes, got int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateArgumentsAgainstAvroSchemaString(tt.args, tt.schemaStr)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errText != "" {
					assert.Contains(t, err.Error(), tt.errText)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestSerializeArgumentsToAvroBinary tests the serializeArgumentsToAvroBinary function.
func TestSerializeArgumentsToAvroBinary(t *testing.T) {
	tests := []struct {
		name           string
		schemaStr      string
		args           map[string]any
		expectError    bool
		errorContains  string
		expectedBinary []byte // Optional: if nil, don't check binary equality, just no error
	}{
		{
			name:      "Valid: Simple Record",
			schemaStr: simpleRecordSchema,
			args: map[string]any{
				"id":   int64(123),
				"name": "test name",
			},
			expectError: false,
		},
		{
			name:      "Valid: Schema With All Primitives (matching getExpectedAvroBinary)",
			schemaStr: schemaWithAllPrimitives,
			args: map[string]any{
				"boolField":   true,
				"intField":    int32(100),
				"longField":   int64(2000),
				"floatField":  float32(3.14),
				"doubleField": float64(6.28),
				"bytesField":  []byte("bytesdata"),
				"stringField": "stringdata",
			},
			expectError: false,
		},
		{
			name:      "Valid: Nested Record",
			schemaStr: nestedRecordSchema,
			args: map[string]any{
				"id": "outerID",
				"inner": map[string]any{
					"value":       int32(101),
					"description": "inner desc",
				},
			},
			expectError: false,
		},
		{
			name:      "Valid: Array of Primitives",
			schemaStr: arraySchemaPrimitive,
			args: map[string]any{
				"stringArray":      []any{"a", "b", "c"},
				"optionalIntArray": []any{int32(1), int32(2)},
			},
			expectError: false,
		},
		{
			name:      "Valid: Map of Records",
			schemaStr: mapSchemaRecord,
			args: map[string]any{
				"recordsMap": map[string]any{
					"recA": map[string]any{"id": int32(1), "status": "active"},
				},
			},
			expectError: false,
		},
		{
			name:      "Valid: Enum",
			schemaStr: enumSchema,
			args: map[string]any{
				"suit": "SPADES",
			},
			expectError: false,
		},
		{
			name:      "Valid: Union - int type",
			schemaStr: unionSchemaSimple,
			args: map[string]any{
				"stringOrInt": int32(123),
			},
			expectError: false,
		},
		{
			name:      "Invalid: Serialization fails due to validation (missing required field)",
			schemaStr: simpleRecordSchema,
			args: map[string]any{
				"id": int64(123), // "name" is missing
			},
			expectError:   true,
			errorContains: "required field 'name' is missing and has no default value",
		},
		{
			name:      "Invalid: Serialization fails due to validation (wrong type)",
			schemaStr: simpleRecordSchema,
			args: map[string]any{
				"id":   "not-a-long",
				"name": "Test Name",
			},
			expectError:   true,
			errorContains: "field 'id': expected long, got string (value: not-a-long)",
		},
		{
			name:          "Invalid: Empty schema string",
			schemaStr:     "",
			args:          map[string]any{"id": int64(1)},
			expectError:   true,
			errorContains: "failed to parse AVRO schema for validation: avro: unknown type: ",
		},
		{
			name:          "Invalid: Malformed schema string",
			schemaStr:     `{"type": "record"`,
			args:          map[string]any{"id": 123},
			expectError:   true,
			errorContains: "failed to parse AVRO schema for validation: avro: unknown type: {\"type\": \"record\"",
		},
		{
			name:      "Valid: schemaWithAllPrimitives - bytes field with string input for serialization",
			schemaStr: schemaWithAllPrimitives,
			args: map[string]any{
				"boolField":   true,
				"intField":    int32(100),
				"longField":   int64(2000),
				"floatField":  float32(3.14),
				"doubleField": float64(6.28),
				"bytesField":  []byte("stringtobytes"),
				"stringField": "stringdata",
			},
			expectError: false, // Should serialize correctly as current code handles string to []byte for bytes type
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualBinary, err := serializeArgumentsToAvroBinary(tt.args, tt.schemaStr)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				assert.Nil(t, actualBinary)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, actualBinary)

				// To validate the binary output, we use the helper.
				// This ensures that our serialization matches a known-good Avro library's output.
				expectedBinary := getExpectedAvroBinary(t, tt.schemaStr, tt.args)
				assert.Equal(t, expectedBinary, actualBinary, "Serialized binary output does not match expected binary.")
			}
		})
	}
}
