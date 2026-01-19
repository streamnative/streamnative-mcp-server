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

package pulsar

import (
	"context"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPulsarAdminTopicToolBuilder(t *testing.T) {
	builder := NewPulsarAdminTopicToolBuilder()

	t.Run("GetName", func(t *testing.T) {
		assert.Equal(t, "pulsar_admin_topic", builder.GetName())
	})

	t.Run("GetRequiredFeatures", func(t *testing.T) {
		features := builder.GetRequiredFeatures()
		assert.NotEmpty(t, features)
		assert.Contains(t, features, "pulsar-admin-topics")
	})

	t.Run("BuildTools_Success", func(t *testing.T) {
		config := builders.ToolBuildConfig{
			ReadOnly: false,
			Features: []string{"pulsar-admin-topics"},
		}

		tools, err := builder.BuildTools(context.Background(), config)
		require.NoError(t, err)
		assert.Len(t, tools, 1)
		assert.Equal(t, "pulsar_admin_topic", tools[0].Definition().Name)
		assert.NotNil(t, tools[0])
	})

	t.Run("BuildTools_ReadOnly", func(t *testing.T) {
		config := builders.ToolBuildConfig{
			ReadOnly: true,
			Features: []string{"pulsar-admin-topics"},
		}

		tools, err := builder.BuildTools(context.Background(), config)
		require.NoError(t, err)
		assert.Len(t, tools, 1)
		assert.Equal(t, "pulsar_admin_topic", tools[0].Definition().Name)
	})

	t.Run("BuildTools_NoFeatures", func(t *testing.T) {
		config := builders.ToolBuildConfig{
			ReadOnly: false,
			Features: []string{"unrelated-feature"},
		}

		tools, err := builder.BuildTools(context.Background(), config)
		require.NoError(t, err)
		assert.Len(t, tools, 0)
	})

	t.Run("Validate_Success", func(t *testing.T) {
		config := builders.ToolBuildConfig{
			Features: []string{"pulsar-admin-topics"},
		}

		err := builder.Validate(config)
		assert.NoError(t, err)
	})

	t.Run("Validate_MissingFeatures", func(t *testing.T) {
		config := builders.ToolBuildConfig{
			Features: []string{"unrelated-feature"},
		}

		err := builder.Validate(config)
		assert.Error(t, err)
	})
}

func TestPulsarAdminTopicToolSchema(t *testing.T) {
	builder := NewPulsarAdminTopicToolBuilder()
	tool, err := builder.buildTopicTool()
	require.NoError(t, err)
	assert.Equal(t, "pulsar_admin_topic", tool.Name)

	schema, ok := tool.InputSchema.(*jsonschema.Schema)
	require.True(t, ok)
	require.NotNil(t, schema.Properties)

	expectedRequired := []string{"resource", "operation"}
	assert.ElementsMatch(t, expectedRequired, schema.Required)

	expectedProps := []string{
		"resource",
		"operation",
		"topic",
		"namespace",
		"partitions",
		"force",
		"non-partitioned",
		"partitioned",
		"per-partition",
		"config",
		"messageId",
	}
	assert.ElementsMatch(t, expectedProps, mapStringKeys(schema.Properties))

	resourceSchema := schema.Properties["resource"]
	require.NotNil(t, resourceSchema)
	assert.Equal(t, pulsarAdminTopicResourceDesc, resourceSchema.Description)

	operationSchema := schema.Properties["operation"]
	require.NotNil(t, operationSchema)
	assert.Equal(t, pulsarAdminTopicOperationDesc, operationSchema.Description)
}

func TestPulsarAdminTopicToolBuilder_ReadOnlyRejectsWrite(t *testing.T) {
	builder := NewPulsarAdminTopicToolBuilder()
	handler := builder.buildTopicHandler(true)

	_, _, err := handler(context.Background(), nil, pulsarAdminTopicInput{
		Resource:  "topic",
		Operation: "create",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "read-only")
}

func mapStringKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
