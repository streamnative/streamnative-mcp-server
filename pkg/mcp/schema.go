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

package mcp

import (
	"fmt"
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"
)

// InputSchema infers a JSON Schema for tool inputs.
// JSON tags control property names and required fields (omitempty/omitzero).
// jsonschema tags provide property descriptions (jsonschema:"...").
func InputSchema[T any]() (*jsonschema.Schema, error) {
	if isAnyType[T]() {
		return &jsonschema.Schema{Type: "object"}, nil
	}
	return buildSchema[T]("input")
}

// OutputSchema infers a JSON Schema for tool outputs.
// When the output type is any, no schema is returned to avoid constraining output.
func OutputSchema[T any]() (*jsonschema.Schema, error) {
	if isAnyType[T]() {
		return nil, nil
	}
	return buildSchema[T]("output")
}

func buildSchema[T any](label string) (*jsonschema.Schema, error) {
	schema, err := schemaForType[T]()
	if err != nil {
		return nil, fmt.Errorf("%s schema: %w", label, err)
	}
	normalizeAdditionalProperties(schema)
	if !isObjectSchema(schema) {
		return nil, fmt.Errorf("%s schema must have type \"object\"", label)
	}
	return schema, nil
}

func schemaForType[T any]() (*jsonschema.Schema, error) {
	t := reflect.TypeFor[T]()
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return jsonschema.ForType(t, &jsonschema.ForOptions{})
}

func isObjectSchema(schema *jsonschema.Schema) bool {
	if schema == nil {
		return false
	}
	return schema.Type == "object"
}

func isAnyType[T any]() bool {
	t := reflect.TypeFor[T]()
	return t.Kind() == reflect.Interface && t.NumMethod() == 0
}

func normalizeAdditionalProperties(schema *jsonschema.Schema) {
	visited := map[*jsonschema.Schema]bool{}
	var walk func(*jsonschema.Schema)
	walk = func(s *jsonschema.Schema) {
		if s == nil || visited[s] {
			return
		}
		visited[s] = true

		if s.Type == "object" && s.Properties != nil && isFalseSchema(s.AdditionalProperties) {
			s.AdditionalProperties = nil
		}

		for _, prop := range s.Properties {
			walk(prop)
		}
		for _, prop := range s.PatternProperties {
			walk(prop)
		}
		for _, def := range s.Defs {
			walk(def)
		}
		for _, def := range s.Definitions {
			walk(def)
		}
		if s.AdditionalProperties != nil && !isFalseSchema(s.AdditionalProperties) {
			walk(s.AdditionalProperties)
		}
		if s.Items != nil {
			walk(s.Items)
		}
		for _, item := range s.PrefixItems {
			walk(item)
		}
		if s.AdditionalItems != nil {
			walk(s.AdditionalItems)
		}
		if s.UnevaluatedItems != nil {
			walk(s.UnevaluatedItems)
		}
		if s.UnevaluatedProperties != nil {
			walk(s.UnevaluatedProperties)
		}
		if s.PropertyNames != nil {
			walk(s.PropertyNames)
		}
		if s.Contains != nil {
			walk(s.Contains)
		}
		for _, subschema := range s.AllOf {
			walk(subschema)
		}
		for _, subschema := range s.AnyOf {
			walk(subschema)
		}
		for _, subschema := range s.OneOf {
			walk(subschema)
		}
		if s.Not != nil {
			walk(s.Not)
		}
		if s.If != nil {
			walk(s.If)
		}
		if s.Then != nil {
			walk(s.Then)
		}
		if s.Else != nil {
			walk(s.Else)
		}
		for _, subschema := range s.DependentSchemas {
			walk(subschema)
		}
	}
	walk(schema)
}

func isFalseSchema(schema *jsonschema.Schema) bool {
	if schema == nil || schema.Not == nil {
		return false
	}
	if !reflect.ValueOf(*schema.Not).IsZero() {
		return false
	}
	clone := *schema
	clone.Not = nil
	return reflect.ValueOf(clone).IsZero()
}
