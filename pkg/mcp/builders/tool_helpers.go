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

package builders

import (
	"github.com/invopop/jsonschema"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// ToolOption is a function that modifies a Tool being built.
// This provides compatibility with the mark3labs builder pattern.
type ToolOption func(*mcpsdk.Tool)

// NewTool creates a new Tool with the given name and options.
func NewTool(name string, opts ...ToolOption) *mcpsdk.Tool {
	tool := &mcpsdk.Tool{
		Name: name,
	}
	for _, opt := range opts {
		opt(tool)
	}
	return tool
}

// WithDescription sets the tool description.
func WithDescription(desc string) ToolOption {
	return func(t *mcpsdk.Tool) {
		t.Description = desc
	}
}

// WithString adds a string property to the input schema.
func WithString(name string, opts ...PropertyOption) ToolOption {
	return func(t *mcpsdk.Tool) {
		prop := &jsonschema.Schema{
			Type: "string",
		}
		isRequired := false
		for _, opt := range opts {
			if opt.required {
				isRequired = true
			}
			if opt.description != "" {
				prop.Description = opt.description
			}
			if len(opt.enumValues) > 0 {
				// Convert []string to []any for jsonschema.Schema.Enum
				enum := make([]any, len(opt.enumValues))
				for i, v := range opt.enumValues {
					enum[i] = v
				}
				prop.Enum = enum
			}
			if opt.defaultString != nil {
				prop.Default = *opt.defaultString
			}
		}
		addPropertyToTool(t, name, prop, isRequired)
	}
}

// WithNumber adds a number property to the input schema.
func WithNumber(name string, opts ...PropertyOption) ToolOption {
	return func(t *mcpsdk.Tool) {
		prop := &jsonschema.Schema{
			Type: "number",
		}
		isRequired := false
		for _, opt := range opts {
			if opt.required {
				isRequired = true
			}
			if opt.description != "" {
				prop.Description = opt.description
			}
			if opt.defaultNumber != nil {
				prop.Default = *opt.defaultNumber
			}
		}
		addPropertyToTool(t, name, prop, isRequired)
	}
}

// WithBoolean adds a boolean property to the input schema.
func WithBoolean(name string, opts ...PropertyOption) ToolOption {
	return func(t *mcpsdk.Tool) {
		prop := &jsonschema.Schema{
			Type: "boolean",
		}
		isRequired := false
		for _, opt := range opts {
			if opt.required {
				isRequired = true
			}
			if opt.description != "" {
				prop.Description = opt.description
			}
			if opt.defaultBool != nil {
				prop.Default = *opt.defaultBool
			}
		}
		addPropertyToTool(t, name, prop, isRequired)
	}
}

// WithObject adds an object property to the input schema.
func WithObject(name string, opts ...PropertyOption) ToolOption {
	return func(t *mcpsdk.Tool) {
		prop := &jsonschema.Schema{
			Type: "object",
		}
		isRequired := false
		for _, opt := range opts {
			if opt.required {
				isRequired = true
			}
			if opt.description != "" {
				prop.Description = opt.description
			}
		}
		addPropertyToTool(t, name, prop, isRequired)
	}
}

// WithArray adds an array property to the input schema.
func WithArray(name string, opts ...PropertyOption) ToolOption {
	return func(t *mcpsdk.Tool) {
		prop := &jsonschema.Schema{
			Type: "array",
		}
		isRequired := false
		for _, opt := range opts {
			if opt.required {
				isRequired = true
			}
			if opt.description != "" {
				prop.Description = opt.description
			}
			if opt.items != nil {
				prop.Items = opt.items
			}
		}
		addPropertyToTool(t, name, prop, isRequired)
	}
}

// Items returns the items schema for array properties.
// This is used to specify the schema of array elements.
// The items parameter can be a *jsonschema.Schema or a map[string]any
// that will be converted to a schema.
func Items(items interface{}) PropertyOption {
	var itemsSchema *jsonschema.Schema
	switch v := items.(type) {
	case *jsonschema.Schema:
		itemsSchema = v
	case map[string]any:
		// Convert map to jsonschema.Schema
		itemsSchema = mapToSchema(v)
	default:
		// Handle map[string]interface{} which is the same as map[string]any
		if m, ok := items.(map[string]interface{}); ok {
			// Convert each value to any type
			converted := make(map[string]any, len(m))
			for k, v := range m {
				converted[k] = v
			}
			itemsSchema = mapToSchema(converted)
		}
	}
	return PropertyOption{items: itemsSchema}
}

// PropertyOption is an option for modifying a property schema.
type PropertyOption struct {
	description  string
	required     bool
	items        *jsonschema.Schema
	enumValues   []string
	defaultBool  *bool
	defaultString *string
	defaultNumber *float64
}

// mapToSchema converts a map[string]any to a jsonschema.Schema.
func mapToSchema(m map[string]any) *jsonschema.Schema {
	schema := &jsonschema.Schema{}
	if typ, ok := m["type"].(string); ok {
		schema.Type = typ
	}
	if desc, ok := m["description"].(string); ok {
		schema.Description = desc
	}
	return schema
}

// Description sets the property description.
func Description(desc string) PropertyOption {
	return PropertyOption{description: desc}
}

// Required marks the property as required.
func Required() PropertyOption {
	return PropertyOption{required: true}
}

// addPropertyToTool adds a property to the tool's input schema.
func addPropertyToTool(tool *mcpsdk.Tool, name string, prop *jsonschema.Schema, isRequired bool) {
	// Initialize InputSchema if needed
	if tool.InputSchema == nil {
		props := jsonschema.NewProperties()
		required := []string{}
		if isRequired {
			required = []string{name}
		}
		tool.InputSchema = &jsonschema.Schema{
			Type:       "object",
			Properties: props,
			Required:   required,
		}
		props.Set(name, prop)
		return
	}

	// Handle *jsonschema.Schema case
	if schema, ok := tool.InputSchema.(*jsonschema.Schema); ok && schema != nil {
		if schema.Properties == nil {
			schema.Properties = jsonschema.NewProperties()
		}
		schema.Properties.Set(name, prop)
		if isRequired {
			if schema.Required == nil {
				schema.Required = []string{name}
			} else {
				schema.Required = append(schema.Required, name)
			}
		}
	}
}

// Enum adds enum values to a PropertyOption.
// This is used to specify a set of allowed values for a property.
func Enum(values ...string) PropertyOption {
	return PropertyOption{enumValues: values}
}

// DefaultBool adds a default boolean value to a PropertyOption.
// This is used to specify the default value for a boolean property.
func DefaultBool(defaultValue bool) PropertyOption {
	return PropertyOption{defaultBool: &defaultValue}
}

// WithToolAnnotation adds metadata annotations to a Tool.
// This is used to specify title, destructive hints, and other UI-related metadata.
func WithToolAnnotation(annotation ToolAnnotation) ToolOption {
	return func(t *mcpsdk.Tool) {
		annotations := &mcpsdk.ToolAnnotations{}
		if annotation.Title != "" {
			annotations.Title = annotation.Title
		}
		if annotation.DestructiveHint != nil {
			annotations.DestructiveHint = annotation.DestructiveHint
		}
		t.Annotations = annotations
	}
}

// ToolAnnotation represents tool metadata for UI purposes.
type ToolAnnotation struct {
	Title           string
	DestructiveHint *bool
}

// DefaultString adds a default string value to a PropertyOption.
// This is used to specify the default value for a string property.
func DefaultString(defaultValue string) PropertyOption {
	return PropertyOption{defaultString: &defaultValue}
}

// DefaultNumber adds a default number value to a PropertyOption.
// This is used to specify the default value for a number property.
func DefaultNumber(defaultValue float64) PropertyOption {
	return PropertyOption{defaultNumber: &defaultValue}
}
