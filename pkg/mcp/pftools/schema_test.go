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

package pftools

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertSchemaToToolInputDefaultsToString(t *testing.T) {
	schema, err := ConvertSchemaToToolInput(nil)
	require.NoError(t, err)
	require.NotNil(t, schema)
	require.Equal(t, "object", schema.Type)

	payload, ok := schema.Properties["payload"]
	require.True(t, ok)
	require.Equal(t, "string", payload.Type)
	require.Contains(t, payload.Description, "plain text")
}

func TestConvertSchemaToToolInputJSONSchema(t *testing.T) {
	fields := []any{map[string]any{"name": "value", "type": "string"}}
	definition := map[string]interface{}{"fields": fields}

	schema, err := ConvertSchemaToToolInput(&SchemaInfo{Type: "JSON", Definition: definition})
	require.NoError(t, err)
	require.NotNil(t, schema)

	payload, ok := schema.Properties["payload"]
	require.True(t, ok)
	require.Equal(t, "string", payload.Type)

	definitionJSON, err := json.Marshal(fields)
	require.NoError(t, err)
	expectedDescription := "The payload of the message, in JSON String format, the schema of the payload in AVRO format is: " + string(definitionJSON)
	require.Equal(t, expectedDescription, payload.Description)
}

func TestConvertSchemaToToolInputRejectsAvro(t *testing.T) {
	schema, err := ConvertSchemaToToolInput(&SchemaInfo{Type: "AVRO", Definition: map[string]interface{}{}})
	require.Error(t, err)
	require.Nil(t, schema)
}
