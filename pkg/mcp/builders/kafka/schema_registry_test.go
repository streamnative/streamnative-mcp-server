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

package kafka

import (
	"context"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKafkaSchemaRegistryToolBuilder(t *testing.T) {
	builder := NewKafkaSchemaRegistryToolBuilder()

	t.Run("GetName", func(t *testing.T) {
		assert.Equal(t, "kafka_schema_registry", builder.GetName())
	})

	t.Run("GetRequiredFeatures", func(t *testing.T) {
		features := builder.GetRequiredFeatures()
		assert.NotEmpty(t, features)
		assert.Contains(t, features, "kafka-admin-schema-registry")
	})

	t.Run("BuildTools_Success", func(t *testing.T) {
		config := builders.ToolBuildConfig{
			ReadOnly: false,
			Features: []string{"kafka-admin-schema-registry"},
		}

		tools, err := builder.BuildTools(context.Background(), config)
		require.NoError(t, err)
		assert.Len(t, tools, 1)
		assert.Equal(t, "kafka_admin_sr", tools[0].Definition().Name)
		assert.NotNil(t, tools[0])
	})

	t.Run("BuildTools_ReadOnly", func(t *testing.T) {
		config := builders.ToolBuildConfig{
			ReadOnly: true,
			Features: []string{"kafka-admin-schema-registry"},
		}

		tools, err := builder.BuildTools(context.Background(), config)
		require.NoError(t, err)
		assert.Len(t, tools, 1)
		assert.Equal(t, "kafka_admin_sr", tools[0].Definition().Name)
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
			Features: []string{"kafka-admin-schema-registry"},
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

func TestKafkaSchemaRegistryToolSchema(t *testing.T) {
	builder := NewKafkaSchemaRegistryToolBuilder()
	tool, err := builder.buildKafkaSchemaRegistryTool()
	require.NoError(t, err)
	assert.Equal(t, "kafka_admin_sr", tool.Name)

	schema, ok := tool.InputSchema.(*jsonschema.Schema)
	require.True(t, ok)
	require.NotNil(t, schema.Properties)

	expectedRequired := []string{"resource", "operation"}
	assert.ElementsMatch(t, expectedRequired, schema.Required)

	expectedProps := []string{
		"resource",
		"operation",
		"subject",
		"version",
		"compatibility",
		"schemaType",
		"schema",
	}
	assert.ElementsMatch(t, expectedProps, mapStringKeys(schema.Properties))

	resourceSchema := schema.Properties["resource"]
	require.NotNil(t, resourceSchema)
	assert.Equal(t, kafkaSchemaRegistryResourceDesc, resourceSchema.Description)

	operationSchema := schema.Properties["operation"]
	require.NotNil(t, operationSchema)
	assert.Equal(t, kafkaSchemaRegistryOperationDesc, operationSchema.Description)

	schemaSchema := schema.Properties["schema"]
	require.NotNil(t, schemaSchema)
	assert.Equal(t, kafkaSchemaRegistrySchemaDesc, schemaSchema.Description)
	assert.Contains(t, schemaSchema.Types, "string")
}

func TestKafkaSchemaRegistryToolBuilder_ReadOnlyRejectsWrite(t *testing.T) {
	builder := NewKafkaSchemaRegistryToolBuilder()
	handler := builder.buildKafkaSchemaRegistryHandler(true)

	_, _, err := handler(context.Background(), nil, kafkaSchemaRegistryInput{
		Resource:  "subject",
		Operation: "create",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "read-only")
}
