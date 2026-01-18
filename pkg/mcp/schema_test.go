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
	"testing"

	"github.com/stretchr/testify/require"
)

type sampleInput struct {
	Name  string `json:"name" jsonschema:"resource name"`
	Count int    `json:"count,omitempty" jsonschema:"optional count"`
}

type nestedInput struct {
	Meta sampleInput `json:"meta"`
}

func TestInputSchemaAllowsAdditionalProperties(t *testing.T) {
	schema, err := InputSchema[sampleInput]()
	require.NoError(t, err)
	require.Equal(t, "object", schema.Type)
	require.NotNil(t, schema.Properties["name"])
	require.NotNil(t, schema.Properties["count"])
	require.Contains(t, schema.Required, "name")
	require.NotContains(t, schema.Required, "count")
	require.Nil(t, schema.AdditionalProperties)
}

func TestInputSchemaAllowsNestedAdditionalProperties(t *testing.T) {
	schema, err := InputSchema[nestedInput]()
	require.NoError(t, err)
	metaSchema := schema.Properties["meta"]
	require.NotNil(t, metaSchema)
	require.Nil(t, metaSchema.AdditionalProperties)
}

func TestInputSchemaMapKeepsAdditionalProperties(t *testing.T) {
	schema, err := InputSchema[map[string]string]()
	require.NoError(t, err)
	require.Equal(t, "object", schema.Type)
	require.NotNil(t, schema.AdditionalProperties)
	require.Equal(t, "string", schema.AdditionalProperties.Type)
}

func TestOutputSchemaAnyIsNil(t *testing.T) {
	schema, err := OutputSchema[any]()
	require.NoError(t, err)
	require.Nil(t, schema)
}
