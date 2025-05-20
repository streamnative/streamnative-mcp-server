package schema

import (
	"fmt"
	"reflect"

	"github.com/hamba/avro/v2"
	"github.com/mark3labs/mcp-go/mcp"
)

// processAvroSchemaStringToMCPToolInput takes an AVRO schema string, parses it,
// and converts it to MCP tool input schema properties.
func processAvroSchemaStringToMCPToolInput(avroSchemaString string) ([]mcp.ToolOption, error) {
	schema, err := avro.Parse(avroSchemaString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AVRO schema: %w", err)
	}

	recordSchema, ok := schema.(*avro.RecordSchema)
	if !ok {
		// If it's not a record, perhaps it's a simpler type that can't be directly mapped to tool options,
		// or we need a different handling strategy. For now, strict record check.
		return nil, fmt.Errorf("expected AVRO record schema at the top level, got %s", reflect.TypeOf(schema).String())
	}

	var opts []mcp.ToolOption
	for _, field := range recordSchema.Fields() {
		fieldOption, err := avroFieldToMcpOption(field)
		if err != nil {
			return nil, fmt.Errorf("failed to convert field '%s': %w", field.Name(), err)
		}
		opts = append(opts, fieldOption)
	}
	return opts, nil
}

// avroFieldToMcpOption converts a single AVRO field to an mcp.ToolOption.
func avroFieldToMcpOption(field *avro.Field) (mcp.ToolOption, error) {
	fieldType := field.Type()
	fieldName := field.Name()

	var description string
	if field.Doc() != "" {
		description = field.Doc()
	} else {
		description = fmt.Sprintf("%s (type: %s)", fieldName, fieldType.String()) // Default description
	}

	isRequired := true
	var underlyingTypeForDefault avro.Schema = fieldType // Used to check default value against non-union type

	if unionSchema, ok := fieldType.(*avro.UnionSchema); ok {
		isNullAble := false
		var nonNullTypes []avro.Schema
		for _, t := range unionSchema.Types() {
			if t.Type() == avro.Null {
				isNullAble = true
			} else {
				nonNullTypes = append(nonNullTypes, t)
			}
		}
		isRequired = !isNullAble

		// If it's a nullable union with one other type (e.g., ["null", "string"]),
		// treat it as that other type for default value and MCP type mapping.
		//nolint:gocritic // This is a valid use of len(nonNullTypes) == 1
		if isNullAble && len(nonNullTypes) == 1 {
			underlyingTypeForDefault = nonNullTypes[0]
		} else if len(nonNullTypes) == 1 {
			// Not nullable, but still a union with one type (should ideally not happen, but handle)
			underlyingTypeForDefault = nonNullTypes[0]
		} else if len(nonNullTypes) > 1 {
			// Complex union (e.g., ["string", "int"]), for MCP, describe as string and mention union nature.
			// Default values for complex unions are tricky with current MCP options.
			// MCP schema might need to be a string with a description of the union.
			props := []mcp.PropertyOption{mcp.Description(description + " (union type: " + fieldType.String() + ")")}
			if isRequired {
				props = append(props, mcp.Required())
			}
			// Default value for complex union is not straightforward to map to a single MCP type's default.
			// We will skip setting mcp.Default... for complex unions for now.
			return mcp.WithString(fieldName, props...), nil
		}
		// If only "null" type was in union, or empty nonNullTypes (invalid schema), this will be caught by later type switch.
	}

	var prop []mcp.PropertyOption
	if isRequired {
		prop = append(prop, mcp.Required())
	}
	prop = append(prop, mcp.Description(description))

	var opt mcp.ToolOption

	// Use underlyingTypeForDefault for determining MCP type and handling default values
	// This handles cases like ["null", "string"] by treating it as "string" for MCP mapping.
	effectiveType := underlyingTypeForDefault.Type()

	switch effectiveType {
	case avro.String:
		if field.HasDefault() {
			if defaultVal, ok := field.Default().(string); ok {
				prop = append(prop, mcp.DefaultString(defaultVal))
			}
		}
		opt = mcp.WithString(fieldName, prop...)
	case avro.Int, avro.Long: // MCP 'number' can represent both
		if field.HasDefault() {
			// Avro library parses numeric defaults as float64 for int/long/float/double from JSON representation
			if defaultVal, ok := field.Default().(float64); ok {
				prop = append(prop, mcp.DefaultNumber(defaultVal))
			} else if defaultIntVal, ok := field.Default().(int); ok { // direct int
				prop = append(prop, mcp.DefaultNumber(float64(defaultIntVal)))
			} else if defaultInt32Val, ok := field.Default().(int32); ok {
				prop = append(prop, mcp.DefaultNumber(float64(defaultInt32Val)))
			} else if defaultInt64Val, ok := field.Default().(int64); ok {
				prop = append(prop, mcp.DefaultNumber(float64(defaultInt64Val)))
			}
		}
		opt = mcp.WithNumber(fieldName, prop...)
	case avro.Float, avro.Double: // MCP 'number' can represent both
		if field.HasDefault() {
			if defaultVal, ok := field.Default().(float64); ok {
				prop = append(prop, mcp.DefaultNumber(defaultVal))
			}
		}
		opt = mcp.WithNumber(fieldName, prop...)
	case avro.Boolean:
		if field.HasDefault() {
			if defaultVal, ok := field.Default().(bool); ok {
				prop = append(prop, mcp.DefaultBool(defaultVal))
			}
		}
		opt = mcp.WithBoolean(fieldName, prop...)
	case avro.Bytes: // MCP represents bytes typically as a base64 encoded string or just string
		if field.HasDefault() {
			// Default for bytes is also string in Avro JSON schema
			if defaultVal, ok := field.Default().(string); ok {
				prop = append(prop, mcp.DefaultString(defaultVal))
			} else if defaultBytes, ok := field.Default().([]byte); ok {
				prop = append(prop, mcp.DefaultString(string(defaultBytes))) // Or base64, TBD by MCP spec
			}
		}
		opt = mcp.WithString(fieldName, prop...) // Or a specific mcp.WithBytes if available/required
	case avro.Array:
		arraySchema, _ := underlyingTypeForDefault.(*avro.ArraySchema) // Safe due to switch
		itemsDef, err := avroSchemaDefinitionToMcpProperties("item", arraySchema.Items(), "Array items")
		if err != nil {
			return nil, fmt.Errorf("failed to convert array items for field '%s': %w", fieldName, err)
		}
		prop = append(prop, mcp.Items(itemsDef))
		// Default for array is not directly supported by mcp.DefaultArray like mcp.DefaultString etc.
		// It would require converting []any to a string representation or specific handling.
		opt = mcp.WithArray(fieldName, prop...)
	case avro.Map:
		mapSchema, _ := underlyingTypeForDefault.(*avro.MapSchema) // Safe due to switch
		// For MCP, map values are usually represented as an object where keys are arbitrary strings
		// and all values conform to a single schema.
		// mcp.Properties expects a map[string]any defining named properties.
		// This is a slight mismatch. MCP's WithObject with mcp.Properties(map[string]any{"*": valuesDef})
		// is one way, or a more flexible mcp.WithMap that takes a value schema.
		// For now, let's assume map values translate to a generic object property definition.
		valuesDef, err := avroSchemaDefinitionToMcpProperties("value", mapSchema.Values(), "Map values")
		if err != nil {
			return nil, fmt.Errorf("failed to convert map values for field '%s': %w", fieldName, err)
		}
		// This isn't a perfect fit for mcp.Properties which expects fixed keys.
		// A better MCP representation for a map might be WithObject where AdditionalProperties is set.
		// For now, we describe it as an object and the value schema applies to all its properties.
		// A more accurate MCP representation might be needed if maps are used extensively.
		// Let's use a single property "*" to denote the schema for all values.
		prop = append(prop, mcp.Properties(map[string]any{"*": valuesDef}))
		opt = mcp.WithObject(fieldName, prop...) // Representing map as a generic object.
	case avro.Record:
		recordSchema, _ := underlyingTypeForDefault.(*avro.RecordSchema) // Safe due to switch
		subProps := make(map[string]any)
		for _, subField := range recordSchema.Fields() {
			// Recursively call avroFieldToMcpOption to get the ToolOption, then extract its schema part?
			// No, avroSchemaDefinitionToMcpProperties is for defining schema of items/values, not named fields.
			// We need to build the properties map for mcp.WithObject.
			// Each subField needs its schema definition.
			var subFieldDescription string
			if subField.Doc() != "" {
				subFieldDescription = subField.Doc()
			} else {
				subFieldDescription = fmt.Sprintf("%s (type: %s)", subField.Name(), subField.Type().String()) // Default description
			}
			subFieldDef, err := avroSchemaDefinitionToMcpProperties(subField.Name(), subField.Type(), subFieldDescription)
			if err != nil {
				return nil, fmt.Errorf("failed to convert sub-field '%s' of record '%s': %w", subField.Name(), fieldName, err)
			}
			subProps[subField.Name()] = subFieldDef
		}
		prop = append(prop, mcp.Properties(subProps))
		opt = mcp.WithObject(fieldName, prop...)
	case avro.Enum:
		enumSchema, _ := underlyingTypeForDefault.(*avro.EnumSchema) // Safe due to switch
		prop = append(prop, mcp.Enum(enumSchema.Symbols()...))
		if field.HasDefault() {
			if defaultVal, ok := field.Default().(string); ok { // Enum default is string
				prop = append(prop, mcp.DefaultString(defaultVal))
			}
		}
		opt = mcp.WithString(fieldName, prop...) // Enum is represented as string in MCP
	case avro.Fixed: // Fixed is like bytes but with a fixed size
		fixedSchema, _ := underlyingTypeForDefault.(*avro.FixedSchema)
		prop = append(prop, mcp.Description(description+fmt.Sprintf(" (fixed size: %d bytes)", fixedSchema.Size())))
		if field.HasDefault() {
			// Default for fixed is also string in Avro JSON schema
			if defaultVal, ok := field.Default().(string); ok {
				prop = append(prop, mcp.DefaultString(defaultVal))
			} else if defaultBytes, ok := field.Default().([]byte); ok {
				prop = append(prop, mcp.DefaultString(string(defaultBytes)))
			}
		}
		opt = mcp.WithString(fieldName, prop...)
	case avro.Null:
		// This case should ideally be handled by the union logic making the field not required.
		// If a field is solely "null", it's an odd schema. For MCP, maybe a string with note.
		// If isRequired is true here (meaning it wasn't a union with null), it's a non-optional null field.
		// This is unusual. Let's represent as a non-required string.
		if isRequired { // Should not happen if only type is null and handled by union logic
			prop = []mcp.PropertyOption{mcp.Description(description + " (null type)")} // remove mcp.Required
		} else {
			prop = append(prop, mcp.Description(description+" (null type)"))
		}
		opt = mcp.WithString(fieldName, prop...) // Or handle as a special ignorable field?
	default:
		// For unknown or unsupported AVRO types, represent as a string in MCP with a description.
		var defaultCaseProps []mcp.PropertyOption
		defaultCaseProps = append(defaultCaseProps, mcp.Description(description+" (unsupported AVRO type: "+string(effectiveType)+")"))
		if isRequired {
			defaultCaseProps = append(defaultCaseProps, mcp.Required())
		}
		opt = mcp.WithString(fieldName, defaultCaseProps...)
	}
	return opt, nil
}

// avroSchemaDefinitionToMcpProperties converts an AVRO schema definition (like for array items or map values)
// into a map[string]any structure suitable for mcp.PropertyOption's Items or Properties.
func avroSchemaDefinitionToMcpProperties(name string, avroDef avro.Schema, description string) (map[string]any, error) {
	props := make(map[string]any)
	if description == "" {
		props["description"] = fmt.Sprintf("Schema for %s", name)
	} else {
		props["description"] = description
	}

	// Handle unions for nested types as well
	var effectiveSchema = avroDef
	if unionSchema, ok := avroDef.(*avro.UnionSchema); ok {
		var nonNullTypes []avro.Schema
		for _, t := range unionSchema.Types() {
			if t.Type() != avro.Null {
				nonNullTypes = append(nonNullTypes, t)
			}
		}
		//nolint:gocritic // This is a valid use of len(nonNullTypes) == 1
		if len(nonNullTypes) == 1 {
			effectiveSchema = nonNullTypes[0]
			props["description"] = props["description"].(string) + " (nullable)"
		} else if len(nonNullTypes) > 1 {
			props["type"] = "string" // Represent complex union as string
			props["description"] = props["description"].(string) + " (union type: " + avroDef.String() + ")"
			return props, nil
		} else { // Only null in union or empty union (invalid)
			props["type"] = "string" // Fallback for null type
			props["description"] = props["description"].(string) + " (effectively null type)"
			return props, nil
		}
	}

	switch effectiveSchema.Type() {
	case avro.String:
		props["type"] = "string"
	case avro.Int, avro.Long:
		props["type"] = "number"
	case avro.Float, avro.Double:
		props["type"] = "number"
	case avro.Boolean:
		props["type"] = "boolean"
	case avro.Bytes, avro.Fixed: // Fixed size bytes
		props["type"] = "string" // Bytes/Fixed represented as string in MCP JSON schema
	case avro.Array:
		arraySchema, _ := effectiveSchema.(*avro.ArraySchema)
		itemsProps, err := avroSchemaDefinitionToMcpProperties("item", arraySchema.Items(), "Array items")
		if err != nil {
			return nil, err
		}
		props["type"] = "array"
		props["items"] = itemsProps
	case avro.Map:
		mapSchema, _ := effectiveSchema.(*avro.MapSchema)
		// MCP object properties are named. Avro map keys are strings, values are of a single schema.
		// We represent this as an object where all properties conform to the map's value schema.
		// The key "*" can signify this pattern.
		valuesProps, err := avroSchemaDefinitionToMcpProperties("value", mapSchema.Values(), "Map values schema")
		if err != nil {
			return nil, err
		}
		props["type"] = "object"
		// To represent an Avro map (string keys, common value schema) in JSON schema properties:
		// we can use "additionalProperties" with the schema of map values.
		// Or, for mcp.Properties, we might define a placeholder like "*".
		// For now, let's return a structure that indicates it's an object,
		// and the `valuesProps` describes the schema for any property within this object.
		// This is a common pattern for map-like structures in JSON Schema if not using additionalProperties.
		props["properties"] = map[string]any{"*": valuesProps} // Indicating all values have this schema.
	case avro.Record:
		recordSchema, _ := effectiveSchema.(*avro.RecordSchema)
		subProps := make(map[string]any)
		for _, field := range recordSchema.Fields() {
			var fieldDescription string
			if field.Doc() != "" {
				fieldDescription = field.Doc()
			} else {
				fieldDescription = fmt.Sprintf("%s (type: %s)", field.Name(), field.Type().String()) // Default description
			}
			fieldProp, err := avroSchemaDefinitionToMcpProperties(field.Name(), field.Type(), fieldDescription)
			if err != nil {
				return nil, err
			}
			subProps[field.Name()] = fieldProp
		}
		props["type"] = "object"
		props["properties"] = subProps
	case avro.Enum:
		enumSchema, _ := effectiveSchema.(*avro.EnumSchema)
		props["type"] = "string"
		props["enum"] = enumSchema.Symbols()
	case avro.Null: // Should be handled by union logic primarily
		props["type"] = "string" // Fallback for a standalone null type.
		props["description"] = props["description"].(string) + " (null type)"
	default:
		props["type"] = "string" // Fallback for unknown types
		props["description"] = props["description"].(string) + " (unknown AVRO type: " + string(effectiveSchema.Type()) + ")"
	}
	return props, nil
}

// validateArgumentsAgainstAvroSchemaString validates arguments against an AVRO schema string.
func validateArgumentsAgainstAvroSchemaString(arguments map[string]any, avroSchemaString string) error {
	schema, err := avro.Parse(avroSchemaString)
	if err != nil {
		return fmt.Errorf("failed to parse AVRO schema for validation: %w", err)
	}

	// Expecting a record schema at the top level for arguments map
	recordSchema, ok := schema.(*avro.RecordSchema)
	if !ok {
		// If the schema is not a record, but arguments are a map, it's a mismatch unless the schema is a map itself.
		// However, tool inputs are typically records/objects.
		// If schema is a single type (e.g. string), arguments shouldn't be a map. This needs clarification on Pulsar schema use.
		// For now, assume top-level schema for arguments is a record.
		return fmt.Errorf("expected AVRO record schema for validating arguments map, got %s", reflect.TypeOf(schema).String())
	}

	// Check for missing required fields
	for _, field := range recordSchema.Fields() {
		fieldName := field.Name()
		isRequired := true // Assume required unless it's a union with null

		if unionSchema, okF := field.Type().(*avro.UnionSchema); okF {
			isNullable := false
			for _, t := range unionSchema.Types() {
				if t.Type() == avro.Null {
					isNullable = true
					break
				}
			}
			if isNullable {
				isRequired = false
			}
		}

		if _, exists := arguments[fieldName]; !exists && isRequired {
			return fmt.Errorf("required field '%s' is missing", fieldName)
		}
	}

	// Validate provided arguments against their schema types
	for argName, argValue := range arguments {
		var field *avro.Field
		for _, f := range recordSchema.Fields() {
			if f.Name() == argName {
				field = f
				break
			}
		}
		if field == nil {
			// Allow extra fields? Or be strict? For now, let's be strict as per common schema validation.
			// Depending on Pulsar's behavior, this might need to be more lenient.
			return fmt.Errorf("unknown field '%s' provided in arguments", argName)
		}
		if err := validateValueAgainstAvroType(argValue, field.Type(), argName); err != nil {
			return err
		}
	}

	return nil
}

// validateValueAgainstAvroType validates a single value against a given AVRO schema type.
// path is for constructing helpful error messages.
func validateValueAgainstAvroType(value any, avroDef avro.Schema, path string) error {
	if value == nil {
		// If value is nil, check if avroDef allows null
		if avroDef.Type() == avro.Null {
			return nil // Explicitly null type allows nil
		}
		if unionSchema, ok := avroDef.(*avro.UnionSchema); ok {
			for _, t := range unionSchema.Types() {
				if t.Type() == avro.Null {
					return nil // Union includes null type
				}
			}
		}
		return fmt.Errorf("field '%s' is null, but schema type '%s' does not permit null", path, avroDef.Type())
	}

	// If avroDef is a union, try to validate against each type in the union.
	if unionSchema, ok := avroDef.(*avro.UnionSchema); ok {
		var lastErr error
		for _, schemaTypeInUnion := range unionSchema.Types() {
			// Skip null type here as we've handled nil value above. If value is not nil, null type won't match.
			if schemaTypeInUnion.Type() == avro.Null {
				continue
			}
			err := validateValueAgainstAvroType(value, schemaTypeInUnion, path)
			if err == nil {
				return nil // Valid against one of the types in the union
			}
			lastErr = err // Keep the last error for context if none match
		}
		if lastErr != nil {
			return fmt.Errorf("field '%s' (value: %v, type: %T) does not match any type in union schema '%s': last error: %w", path, value, value, unionSchema.String(), lastErr)
		}
		// If union was only ["null"] and value is not nil, this will be an error.
		return fmt.Errorf("field '%s' (value: %v) of type %T does not match union schema '%s' (no non-null types matched or union is only null)", path, value, value, unionSchema.String())
	}

	// Non-union type validation
	switch avroDef.Type() {
	case avro.String:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("field '%s': expected string, got %T (value: %v)", path, value, value)
		}
	case avro.Int:
		// JSON unmarshals numbers into float64. Avro library might expect specific int types or handle conversion.
		// We should be flexible with input types that can reasonably convert to int.
		switch value.(type) {
		case int, int8, int16, int32, int64, float32, float64: // float32/64 if they are whole numbers
			// Further check if float64 is a whole number
			if fVal, ok := value.(float64); ok && fVal != float64(int64(fVal)) {
				return fmt.Errorf("field '%s': expected int, got float64 with fractional part (value: %v)", path, value)
			}
			if fVal, ok := value.(float32); ok && fVal != float32(int32(fVal)) {
				return fmt.Errorf("field '%s': expected int, got float32 with fractional part (value: %v)", path, value)
			}
			return nil // Acceptable numeric types
		default:
			return fmt.Errorf("field '%s': expected int, got %T (value: %v)", path, value, value)
		}
	case avro.Long:
		switch value.(type) {
		case int, int8, int16, int32, int64, float32, float64: // float32/64 if they are whole numbers
			if fVal, ok := value.(float64); ok && fVal != float64(int64(fVal)) {
				return fmt.Errorf("field '%s': expected long, got float64 with fractional part (value: %v)", path, value)
			}
			if fVal, ok := value.(float32); ok && fVal != float32(int64(fVal)) { // Note: int64 cast might be too large for float32 precision
				return fmt.Errorf("field '%s': expected long, got float32 with fractional part (value: %v)", path, value)
			}
			return nil
		default:
			return fmt.Errorf("field '%s': expected long, got %T (value: %v)", path, value, value)
		}
	case avro.Float:
		switch value.(type) {
		case float32, float64, int, int8, int16, int32, int64: // int types are fine
			return nil
		default:
			return fmt.Errorf("field '%s': expected float, got %T (value: %v)", path, value, value)
		}
	case avro.Double:
		switch value.(type) {
		case float32, float64, int, int8, int16, int32, int64: // int types are fine
			return nil
		default:
			return fmt.Errorf("field '%s': expected double, got %T (value: %v)", path, value, value)
		}
	case avro.Boolean:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("field '%s': expected boolean, got %T (value: %v)", path, value, value)
		}
	case avro.Bytes, avro.Fixed: // Bytes and Fixed are []byte in Avro Go, but string might come from JSON
		if _, okStr := value.(string); okStr { // Allow string (e.g. base64 encoded) for bytes if that's the convention
			return nil // Assume string is an acceptable representation for bytes here.
		}
		// Native []byte is not directly produced by json.Unmarshal into map[string]any
		// This needs alignment on how bytes are passed in `arguments`.
		// If arguments come from JSON, bytes are likely base64 strings.
		// For now, we primarily expect string.
		// if _, okBytes := value.([]byte); !okBytes {
		//	 return fmt.Errorf("field '%s': expected bytes (string or []byte), got %T (value: %v)", path, value, value)
		// }
		return fmt.Errorf("field '%s': expected string for bytes/fixed, got %T (value: %v)", path, value, value)

	case avro.Array:
		arrSchema, _ := avroDef.(*avro.ArraySchema)
		sliceVal, ok := value.([]any) // JSON unmarshals to []any
		if !ok {
			// Check if it's a typed slice, e.g. []string, []map[string]any, etc.
			// This requires more reflection if we want to support e.g. []string directly.
			// For map[string]any from JSON, []any is standard.
			return fmt.Errorf("field '%s': expected array (slice of any), got %T (value: %v)", path, value, value)
		}
		for i, item := range sliceVal {
			if err := validateValueAgainstAvroType(item, arrSchema.Items(), fmt.Sprintf("%s[%d]", path, i)); err != nil {
				return err
			}
		}
	case avro.Map:
		mapSchema, _ := avroDef.(*avro.MapSchema)
		mapVal, ok := value.(map[string]any) // JSON unmarshals to map[string]any
		if !ok {
			return fmt.Errorf("field '%s': expected map (map[string]any), got %T (value: %v)", path, value, value)
		}
		for k, v := range mapVal {
			if err := validateValueAgainstAvroType(v, mapSchema.Values(), fmt.Sprintf("%s.%s", path, k)); err != nil {
				return err
			}
		}
	case avro.Record:
		recSchema, _ := avroDef.(*avro.RecordSchema)
		mapVal, ok := value.(map[string]any) // JSON unmarshals to map[string]any
		if !ok {
			return fmt.Errorf("field '%s': expected object (map[string]any) for record, got %T (value: %v)", path, value, value)
		}
		// Check required fields within the record
		for _, f := range recSchema.Fields() {
			isFieldRequired := true
			if unionF, okF := f.Type().(*avro.UnionSchema); okF {
				isNullableF := false
				for _, t := range unionF.Types() {
					if t.Type() == avro.Null {
						isNullableF = true
						break
					}
				}
				if isNullableF {
					isFieldRequired = false
				}
			}
			if _, exists := mapVal[f.Name()]; !exists && isFieldRequired {
				return fmt.Errorf("field '%s.%s' is required but missing", path, f.Name())
			}
		}
		// Validate present fields
		for k, v := range mapVal {
			var recField *avro.Field
			for _, f := range recSchema.Fields() {
				if f.Name() == k {
					recField = f
					break
				}
			}
			if recField == nil {
				return fmt.Errorf("field '%s.%s' is not defined in record schema", path, k)
			}
			if err := validateValueAgainstAvroType(v, recField.Type(), fmt.Sprintf("%s.%s", path, k)); err != nil {
				return err
			}
		}
	case avro.Enum:
		enumSchema, _ := avroDef.(*avro.EnumSchema)
		strVal, ok := value.(string)
		if !ok {
			return fmt.Errorf("field '%s': expected string for enum, got %T (value: %v)", path, value, value)
		}
		isValidSymbol := false
		for _, s := range enumSchema.Symbols() {
			if s == strVal {
				isValidSymbol = true
				break
			}
		}
		if !isValidSymbol {
			return fmt.Errorf("field '%s': value '%s' is not a valid symbol for enum %s. Valid symbols: %v", path, strVal, enumSchema.Name(), enumSchema.Symbols())
		}
	case avro.Null: // This case should be hit if a non-union field is explicitly "null" type.
		// If value is not nil, it's an error. Nil value handled at the start of the function.
		// This means value is non-nil here.
		return fmt.Errorf("field '%s': schema type is explicitly 'null' but received non-nil value %T (value: %v)", path, value, value)

	default:
		return fmt.Errorf("field '%s': unsupported AVRO type '%s' for validation", path, avroDef.Type())
	}
	return nil
}

// serializeArgumentsToAvroBinary validates arguments against an AVRO schema string
// and then serializes them to AVRO binary format.
func serializeArgumentsToAvroBinary(arguments map[string]any, avroSchemaString string) ([]byte, error) {
	// First, validate arguments.
	// The validation logic already parses the schema string.
	if err := validateArgumentsAgainstAvroSchemaString(arguments, avroSchemaString); err != nil {
		return nil, fmt.Errorf("arguments validation failed before AVRO serialization: %w", err)
	}

	// Parse schema again for marshaling (or pass parsed schema from validation if we refactor to return it)
	schema, err := avro.Parse(avroSchemaString)
	if err != nil {
		// This error should ideally not happen if validation passed, as it also parses.
		return nil, fmt.Errorf("failed to parse AVRO schema for serialization (should have been caught by validation): %w", err)
	}

	// Marshal to AVRO binary
	data, err := avro.Marshal(schema, arguments)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments to AVRO binary: %w", err)
	}
	return data, nil
}
