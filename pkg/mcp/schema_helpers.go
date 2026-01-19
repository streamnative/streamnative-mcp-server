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
	"encoding/json"

	"github.com/google/jsonschema-go/jsonschema"
)

func setSchemaPropertyDefault(schema *jsonschema.Schema, name string, value any) {
	if schema == nil || schema.Properties == nil {
		return
	}
	property, ok := schema.Properties[name]
	if !ok || property == nil {
		return
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return
	}
	property.Default = raw
}

func setSchemaPropertyEnum(schema *jsonschema.Schema, name string, values []string) {
	if schema == nil || schema.Properties == nil {
		return
	}
	property, ok := schema.Properties[name]
	if !ok || property == nil {
		return
	}
	enumValues := make([]any, 0, len(values))
	for _, value := range values {
		enumValues = append(enumValues, value)
	}
	property.Enum = enumValues
}
