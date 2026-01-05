// Copyright 2025 StreamNative
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package schema

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"

	"github.com/hamba/avro/v2"
	"github.com/invopop/jsonschema"
)

// processAvroSchemaStringToMCPToolInput takes an AVRO schema string, parses it,
// and converts it to jsonschema.Schema for tool input.
func processAvroSchemaStringToMCPToolInput(avroSchemaString string) (*jsonschema.Schema, error) {
	schema, err := avro.Parse(avroSchemaString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AVRO schema: %w", err)
	}

	recordSchema, ok := schema.(*avro.RecordSchema)
	if !ok {
		return nil, fmt.Errorf("expected AVRO record schema at the top level, got %s", reflect.TypeOf(schema).String())
	}

	props := jsonschema.NewProperties()
	required := make([]string, 0)

	for _, field := range recordSchema.Fields() {
		fieldSchema, isRequired, err := avroFieldToJsonSchema(field)
		if err != nil {
			return nil, fmt.Errorf("failed to convert field '%s': %w", field.Name(), err)
		}
		props.Set(field.Name(), fieldSchema)
		if isRequired {
			required = append(required, field.Name())
		}
	}

	result := &jsonschema.Schema{
		Type:       "object",
		Properties: props,
	}
	if len(required) > 0 {
		result.Required = required
	}

	return result, nil
}

// avroFieldToJsonSchema converts a single AVRO field to jsonschema.Schema.
func avroFieldToJsonSchema(field *avro.Field) (*jsonschema.Schema, bool, error) {
	fieldType := field.Type()
	fieldName := field.Name()

	var description string
	if field.Doc() != "" {
		description = field.Doc()
	} else {
		description = fmt.Sprintf("%s (type: %s)", fieldName, strings.ReplaceAll(fieldType.String(), "\"", ""))
	}

	isRequired := true
	var underlyingTypeForDefault avro.Schema = fieldType

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

		if isNullAble && len(nonNullTypes) == 1 {
			underlyingTypeForDefault = nonNullTypes[0]
		} else if len(nonNullTypes) == 1 {
			underlyingTypeForDefault = nonNullTypes[0]
		} else if len(nonNullTypes) > 1 {
			// Complex union - represent as string with description
			schema := &jsonschema.Schema{
				Type:        "string",
				Description: description + " (union type: " + strings.ReplaceAll(fieldType.String(), "\"", "") + ")",
			}
			return schema, isRequired, nil
		}
	}

	schema := &jsonschema.Schema{
		Description: description,
	}

	// Use underlyingTypeForDefault for determining type and handling default values
	effectiveType := underlyingTypeForDefault.Type()

	switch effectiveType {
	case avro.String:
		schema.Type = "string"
		if field.HasDefault() {
			if defaultVal, ok := field.Default().(string); ok {
				schema.Default = defaultVal
			}
		}
	case avro.Int, avro.Long:
		schema.Type = "number"
		if field.HasDefault() {
			if defaultVal, ok := field.Default().(float64); ok {
				schema.Default = defaultVal
			} else if defaultIntVal, ok := field.Default().(int); ok {
				schema.Default = float64(defaultIntVal)
			} else if defaultInt32Val, ok := field.Default().(int32); ok {
				schema.Default = float64(defaultInt32Val)
			} else if defaultInt64Val, ok := field.Default().(int64); ok {
				schema.Default = float64(defaultInt64Val)
			}
		}
	case avro.Float, avro.Double:
		schema.Type = "number"
		if field.HasDefault() {
			if defaultVal, ok := field.Default().(float64); ok {
				schema.Default = defaultVal
			}
		}
	case avro.Boolean:
		schema.Type = "boolean"
		if field.HasDefault() {
			if defaultVal, ok := field.Default().(bool); ok {
				schema.Default = defaultVal
			}
		}
	case avro.Bytes, avro.Fixed:
		schema.Type = "string"
		if field.HasDefault() {
			if defaultVal, ok := field.Default().(string); ok {
				schema.Default = defaultVal
			} else if defaultBytes, ok := field.Default().([]byte); ok {
				schema.Default = string(defaultBytes)
			}
		}
		if fixedSchema, ok := underlyingTypeForDefault.(*avro.FixedSchema); ok {
			schema.Description = fmt.Sprintf("%s (fixed size: %d bytes)", description, fixedSchema.Size())
		}
	case avro.Array:
		arraySchema, _ := underlyingTypeForDefault.(*avro.ArraySchema)
		itemsSchema, err := avroSchemaDefinitionToJsonSchema(arraySchema.Items(), "Array items")
		if err != nil {
			return nil, false, fmt.Errorf("failed to convert array items for field '%s': %w", fieldName, err)
		}
		schema.Type = "array"
		schema.Items = itemsSchema
	case avro.Map:
		mapSchema, _ := underlyingTypeForDefault.(*avro.MapSchema)
		valuesSchema, err := avroSchemaDefinitionToJsonSchema(mapSchema.Values(), "Map values")
		if err != nil {
			return nil, false, fmt.Errorf("failed to convert map values for field '%s': %w", fieldName, err)
		}
		schema.Type = "object"
		schema.AdditionalProperties = valuesSchema
	case avro.Record:
		recordSchema, _ := underlyingTypeForDefault.(*avro.RecordSchema)
		subProps := jsonschema.NewProperties()
		subRequired := make([]string, 0)
		for _, subField := range recordSchema.Fields() {
			subFieldSchema, subIsRequired, err := avroFieldToJsonSchema(subField)
			if err != nil {
				return nil, false, fmt.Errorf("failed to convert sub-field '%s' of record '%s': %w", subField.Name(), fieldName, err)
			}
			subProps.Set(subField.Name(), subFieldSchema)
			if subIsRequired {
				subRequired = append(subRequired, subField.Name())
			}
		}
		schema.Type = "object"
		schema.Properties = subProps
		if len(subRequired) > 0 {
			schema.Required = subRequired
		}
	case avro.Enum:
		enumSchema, _ := underlyingTypeForDefault.(*avro.EnumSchema)
		schema.Type = "string"
		schema.Enum = make([]interface{}, len(enumSchema.Symbols()))
		for i, symbol := range enumSchema.Symbols() {
			schema.Enum[i] = symbol
		}
		if field.HasDefault() {
			if defaultVal, ok := field.Default().(string); ok {
				schema.Default = defaultVal
			}
		}
	case avro.Null:
		schema.Type = "string"
		schema.Description = description + " (null type)"
		if isRequired {
			isRequired = false
		}
	default:
		schema.Type = "string"
		schema.Description = description + " (unsupported AVRO type: " + string(effectiveType) + ")"
	}

	return schema, isRequired, nil
}

// avroSchemaDefinitionToJsonSchema converts an AVRO schema definition (like for array items or map values)
// into jsonschema.Schema.
func avroSchemaDefinitionToJsonSchema(avroDef avro.Schema, description string) (*jsonschema.Schema, error) {
	schema := &jsonschema.Schema{
		Description: description,
	}

	if description == "" {
		schema.Description = fmt.Sprintf("Schema for type")
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
		if len(nonNullTypes) == 1 {
			effectiveSchema = nonNullTypes[0]
			schema.Description = schema.Description + " (nullable)"
		} else if len(nonNullTypes) > 1 {
			schema.Type = "string"
			schema.Description = schema.Description + " (union type: " + strings.ReplaceAll(avroDef.String(), "\"", "") + ")"
			return schema, nil
		} else {
			schema.Type = "string"
			schema.Description = schema.Description + " (effectively null type)"
			return schema, nil
		}
	}

	switch effectiveSchema.Type() {
	case avro.String:
		schema.Type = "string"
	case avro.Int, avro.Long:
		schema.Type = "number"
	case avro.Float, avro.Double:
		schema.Type = "number"
	case avro.Boolean:
		schema.Type = "boolean"
	case avro.Bytes, avro.Fixed:
		schema.Type = "string"
	case avro.Array:
		arraySchema, _ := effectiveSchema.(*avro.ArraySchema)
		itemsSchema, err := avroSchemaDefinitionToJsonSchema(arraySchema.Items(), "Array items")
		if err != nil {
			return nil, err
		}
		schema.Type = "array"
		schema.Items = itemsSchema
	case avro.Map:
		mapSchema, _ := effectiveSchema.(*avro.MapSchema)
		valuesSchema, err := avroSchemaDefinitionToJsonSchema(mapSchema.Values(), "Map values schema")
		if err != nil {
			return nil, err
		}
		schema.Type = "object"
		schema.AdditionalProperties = valuesSchema
	case avro.Record:
		recordSchema, _ := effectiveSchema.(*avro.RecordSchema)
		subProps := jsonschema.NewProperties()
		subRequired := make([]string, 0)
		for _, field := range recordSchema.Fields() {
			fieldSchema, fieldRequired, err := avroFieldToJsonSchema(field)
			if err != nil {
				return nil, err
			}
			subProps.Set(field.Name(), fieldSchema)
			if fieldRequired {
				subRequired = append(subRequired, field.Name())
			}
		}
		schema.Type = "object"
		schema.Properties = subProps
		if len(subRequired) > 0 {
			schema.Required = subRequired
		}
	case avro.Enum:
		enumSchema, _ := effectiveSchema.(*avro.EnumSchema)
		schema.Type = "string"
		schema.Enum = make([]interface{}, len(enumSchema.Symbols()))
		for i, symbol := range enumSchema.Symbols() {
			schema.Enum[i] = symbol
		}
	case avro.Null:
		schema.Type = "string"
		schema.Description = schema.Description + " (null type)"
	default:
		schema.Type = "string"
		schema.Description = schema.Description + " (unknown AVRO type: " + string(effectiveSchema.Type()) + ")"
	}

	return schema, nil
}

// validateArgumentsAgainstAvroSchemaString validates arguments against an AVRO schema string.
func validateArgumentsAgainstAvroSchemaString(arguments map[string]any, avroSchemaString string) error {
	schema, err := avro.Parse(avroSchemaString)
	if err != nil {
		return fmt.Errorf("failed to parse AVRO schema for validation: %w", err)
	}

	recordSchema, ok := schema.(*avro.RecordSchema)
	if !ok {
		return fmt.Errorf("expected AVRO record schema for validating arguments map, got %s", reflect.TypeOf(schema).String())
	}

	// Check for missing required fields
	for _, field := range recordSchema.Fields() {
		fieldName := field.Name()
		isReq := true
		if unionSchemaVal, ok := field.Type().(*avro.UnionSchema); ok {
			isNullableInUnion := false
			for _, t := range unionSchemaVal.Types() {
				if t.Type() == avro.Null {
					isNullableInUnion = true
					break
				}
			}
			isReq = !isNullableInUnion
		}

		value, valueOk := arguments[fieldName]

		if !valueOk {
			if isReq && !field.HasDefault() {
				return fmt.Errorf("required field '%s' is missing and has no default value", fieldName)
			}
			continue
		}

		if err := validateValueAgainstAvroType(value, field.Type(), fieldName); err != nil {
			return err
		}
	}

	// Check for extra fields
	for argName := range arguments {
		foundInSchema := false
		for _, field := range recordSchema.Fields() {
			if field.Name() == argName {
				foundInSchema = true
				break
			}
		}
		if !foundInSchema {
			return fmt.Errorf("unknown field '%s' provided in arguments", argName)
		}
	}

	return nil
}

// validateValueAgainstAvroType validates a single value against a given AVRO schema type.
func validateValueAgainstAvroType(value any, avroDef avro.Schema, path string) error {
	if value == nil {
		if avroDef.Type() == avro.Null {
			return nil
		}
		if unionSchema, ok := avroDef.(*avro.UnionSchema); ok {
			for _, t := range unionSchema.Types() {
				if t.Type() == avro.Null {
					return nil
				}
			}
		}
		return fmt.Errorf("field '%s' is null, but schema type '%s' does not permit null", path, avroDef.Type())
	}

	if unionSchema, ok := avroDef.(*avro.UnionSchema); ok {
		var lastErr error
		for _, schemaTypeInUnion := range unionSchema.Types() {
			if schemaTypeInUnion.Type() == avro.Null {
				continue
			}
			err := validateValueAgainstAvroType(value, schemaTypeInUnion, path)
			if err == nil {
				return nil
			}
			lastErr = err
		}
		if lastErr != nil {
			return fmt.Errorf("field '%s' (value: %v, type: %T) does not match any type in union schema '%s': last error: %w", path, value, value, unionSchema.String(), lastErr)
		}
		return fmt.Errorf("field '%s' (value: %v) of type %T does not match union schema '%s'", path, value, value, unionSchema.String())
	}

	switch avroDef.Type() {
	case avro.String:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("field '%s': expected string, got %T (value: %v)", path, value, value)
		}
	case avro.Int:
		switch value.(type) {
		case int, int8, int16, int32, int64, float32, float64:
			if fVal, ok := value.(float64); ok && fVal != float64(int64(fVal)) {
				return fmt.Errorf("field '%s': expected int, got float64 with fractional part (value: %v)", path, value)
			}
			if fVal, ok := value.(float32); ok && fVal != float32(int32(fVal)) {
				return fmt.Errorf("field '%s': expected int, got float32 with fractional part (value: %v)", path, value)
			}
			return nil
		default:
			return fmt.Errorf("field '%s': expected int, got %T (value: %v)", path, value, value)
		}
	case avro.Long:
		switch value.(type) {
		case int, int8, int16, int32, int64, float32, float64:
			if fVal, ok := value.(float64); ok && fVal != float64(int64(fVal)) {
				return fmt.Errorf("field '%s': expected long, got float64 with fractional part (value: %v)", path, value)
			}
			return nil
		default:
			return fmt.Errorf("field '%s': expected long, got %T (value: %v)", path, value, value)
		}
	case avro.Float:
		switch value.(type) {
		case float32, float64, int, int8, int16, int32, int64:
			return nil
		default:
			return fmt.Errorf("field '%s': expected float, got %T (value: %v)", path, value, value)
		}
	case avro.Double:
		switch value.(type) {
		case float32, float64, int, int8, int16, int32, int64:
			return nil
		default:
			return fmt.Errorf("field '%s': expected double, got %T (value: %v)", path, value, value)
		}
	case avro.Boolean:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("field '%s': expected boolean, got %T (value: %v)", path, value, value)
		}
	case avro.Bytes:
		if _, okStr := value.(string); okStr {
			return nil
		}
		if _, okBytes := value.([]byte); okBytes {
			return nil
		}
		return fmt.Errorf("field '%s': expected string or []byte for bytes, got %T (value: %v)", path, value, value)
	case avro.Fixed:
		if _, ok := value.(uint64); ok {
			return nil
		}
		return fmt.Errorf("field '%s': expected uint64 for fixed, got %T (value: %v)", path, value, value)
	case avro.Array:
		arrSchema, _ := avroDef.(*avro.ArraySchema)
		sliceVal, ok := value.([]any)
		if !ok {
			return fmt.Errorf("field '%s': expected array (slice of any), got %T (value: %v)", path, value, value)
		}
		for i, item := range sliceVal {
			if err := validateValueAgainstAvroType(item, arrSchema.Items(), fmt.Sprintf("%s[%d]", path, i)); err != nil {
				return err
			}
		}
	case avro.Map:
		mapSchema, _ := avroDef.(*avro.MapSchema)
		mapVal, ok := value.(map[string]any)
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
		mapVal, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("field '%s': expected object (map[string]any) for record, got %T (value: %v)", path, value, value)
		}
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
			return fmt.Errorf("field '%s': value '%s' is not a valid symbol for enum %s. Valid symbols: %v", path, strVal, enumSchema.FullName(), enumSchema.Symbols())
		}
	case avro.Null:
		if value != nil {
			return fmt.Errorf("field '%s': schema type is explicitly 'null' but received non-nil value %T (value: %v)", path, value, value)
		}
	default:
		return fmt.Errorf("field '%s': unsupported AVRO type '%s' for validation", path, avroDef.Type())
	}
	return nil
}

// serializeArgumentsToAvroBinary validates arguments against an AVRO schema string
// and then serializes them to AVRO binary format.
func serializeArgumentsToAvroBinary(arguments map[string]any, avroSchemaString string) ([]byte, error) {
	if err := validateArgumentsAgainstAvroSchemaString(arguments, avroSchemaString); err != nil {
		return nil, fmt.Errorf("arguments validation failed before AVRO serialization: %w", err)
	}

	schema, err := avro.Parse(avroSchemaString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AVRO schema for serialization: %w", err)
	}

	coercedArgs := make(map[string]any, len(arguments))
	for k, v := range arguments {
		coercedArgs[k] = v
	}

	recordSchema, ok := schema.(*avro.RecordSchema)
	if !ok {
		return nil, fmt.Errorf("parsed schema is not a record schema, cannot prepare arguments for serialization")
	}

	for _, field := range recordSchema.Fields() {
		fieldName := field.Name()
		val, argExists := arguments[fieldName]
		if !argExists {
			continue
		}

		fieldType := field.Type().Type()
		if unionSchema, isUnion := field.Type().(*avro.UnionSchema); isUnion {
			for _, unionMemberType := range unionSchema.Types() {
				if unionMemberType.Type() == avro.Bytes || unionMemberType.Type() == avro.Fixed {
					fieldType = unionMemberType.Type()
					break
				}
			}
		}

		if fieldType == avro.Bytes {
			if strVal, isStr := val.(string); isStr {
				decodedBytes, err := base64.StdEncoding.DecodeString(strVal)
				if err == nil {
					coercedArgs[fieldName] = decodedBytes
				} else {
					coercedArgs[fieldName] = []byte(strVal)
				}
			}
		} else if fieldType == avro.Fixed {
			if strVal, isStr := val.(string); isStr {
				fixedSchema, _ := field.Type().(*avro.FixedSchema)
				if actualUnionFieldSchema, okUnion := field.Type().(*avro.UnionSchema); okUnion {
					for _, ut := range actualUnionFieldSchema.Types() {
						if fs, okUFS := ut.(*avro.FixedSchema); okUFS {
							fixedSchema = fs
							break
						}
					}
				}
				if fixedSchema != nil {
					decodedBytes, err := base64.StdEncoding.DecodeString(strVal)
					if err == nil {
						if len(decodedBytes) == fixedSchema.Size() {
							fixedArray := reflect.New(reflect.ArrayOf(fixedSchema.Size(), reflect.TypeOf(byte(0)))).Elem()
							reflect.Copy(fixedArray, reflect.ValueOf(decodedBytes))
							coercedArgs[fieldName] = fixedArray.Interface()
						} else {
							return nil, fmt.Errorf("field '%s' (fixed[%d]): base64 decoded string has length %d, expected %d", fieldName, fixedSchema.Size(), len(decodedBytes), fixedSchema.Size())
						}
					}
				}
			} else if byteSlice, isSlice := val.([]byte); isSlice {
				fixedSchema, _ := field.Type().(*avro.FixedSchema)
				if actualUnionFieldSchema, okUnion := field.Type().(*avro.UnionSchema); okUnion {
					for _, ut := range actualUnionFieldSchema.Types() {
						if fs, okUFS := ut.(*avro.FixedSchema); okUFS {
							fixedSchema = fs
							break
						}
					}
				}
				if fixedSchema != nil && len(byteSlice) == fixedSchema.Size() {
					fixedArray := reflect.New(reflect.ArrayOf(fixedSchema.Size(), reflect.TypeOf(byte(0)))).Elem()
					reflect.Copy(fixedArray, reflect.ValueOf(byteSlice))
					coercedArgs[fieldName] = fixedArray.Interface()
				} else if fixedSchema != nil && len(byteSlice) != fixedSchema.Size() {
					return nil, fmt.Errorf("field '%s' (fixed[%d]): provided []byte has length %d, expected %d", fieldName, fixedSchema.Size(), len(byteSlice), fixedSchema.Size())
				}
			}
		}
	}

	return avro.Marshal(schema, coercedArgs)
}
