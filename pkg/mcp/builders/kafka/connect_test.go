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

func TestNewKafkaConnectToolBuilder(t *testing.T) {
	builder := NewKafkaConnectToolBuilder()

	assert.NotNil(t, builder)
	assert.Equal(t, "kafka_connect", builder.GetName())

	expectedFeatures := []string{
		"kafka-admin-kafka-connect",
		"all",
		"all-kafka",
		"kafka-admin",
	}
	assert.Equal(t, expectedFeatures, builder.GetRequiredFeatures())

	metadata := builder.GetMetadata()
	assert.Equal(t, "kafka_connect", metadata.Name)
	assert.Equal(t, "1.0.0", metadata.Version)
	assert.Equal(t, "Kafka Connect administration tools", metadata.Description)
	assert.Equal(t, "kafka_admin", metadata.Category)
	assert.Equal(t, []string{"kafka", "connect", "admin"}, metadata.Tags)
}

func TestKafkaConnectToolBuilder_BuildTools_Success(t *testing.T) {
	builder := NewKafkaConnectToolBuilder()

	config := builders.ToolBuildConfig{
		ReadOnly: false,
		Features: []string{"kafka-admin-kafka-connect"},
	}

	tools, err := builder.BuildTools(context.Background(), config)

	require.NoError(t, err)
	require.Len(t, tools, 1)

	tool := tools[0]
	assert.NotNil(t, tool)
	assert.NotNil(t, tool.Definition())
	assert.Equal(t, "kafka_admin_connect", tool.Definition().Name)
}

func TestKafkaConnectToolBuilder_BuildTools_WithAllFeature(t *testing.T) {
	builder := NewKafkaConnectToolBuilder()

	config := builders.ToolBuildConfig{
		ReadOnly: true,
		Features: []string{"all"},
	}

	tools, err := builder.BuildTools(context.Background(), config)

	require.NoError(t, err)
	require.Len(t, tools, 1)

	tool := tools[0]
	assert.Equal(t, "kafka_admin_connect", tool.Definition().Name)
}

func TestKafkaConnectToolBuilder_BuildTools_WithAllKafkaFeature(t *testing.T) {
	builder := NewKafkaConnectToolBuilder()

	config := builders.ToolBuildConfig{
		ReadOnly: false,
		Features: []string{"all-kafka"},
	}

	tools, err := builder.BuildTools(context.Background(), config)

	require.NoError(t, err)
	require.Len(t, tools, 1)
}

func TestKafkaConnectToolBuilder_BuildTools_WithKafkaAdminFeature(t *testing.T) {
	builder := NewKafkaConnectToolBuilder()

	config := builders.ToolBuildConfig{
		ReadOnly: false,
		Features: []string{"kafka-admin"},
	}

	tools, err := builder.BuildTools(context.Background(), config)

	require.NoError(t, err)
	require.Len(t, tools, 1)
}

func TestKafkaConnectToolBuilder_BuildTools_NoMatchingFeatures(t *testing.T) {
	builder := NewKafkaConnectToolBuilder()

	config := builders.ToolBuildConfig{
		ReadOnly: false,
		Features: []string{"pulsar-admin", "other-feature"},
	}

	tools, err := builder.BuildTools(context.Background(), config)

	require.NoError(t, err)
	assert.Len(t, tools, 0) // Should return empty slice, not error
}

func TestKafkaConnectToolBuilder_BuildTools_EmptyFeatures(t *testing.T) {
	builder := NewKafkaConnectToolBuilder()

	config := builders.ToolBuildConfig{
		ReadOnly: false,
		Features: []string{},
	}

	tools, err := builder.BuildTools(context.Background(), config)

	require.NoError(t, err)
	assert.Len(t, tools, 0)
}

func TestKafkaConnectToolBuilder_Validate_Success(t *testing.T) {
	builder := NewKafkaConnectToolBuilder()

	config := builders.ToolBuildConfig{
		ReadOnly: false,
		Features: []string{"kafka-admin-kafka-connect"},
	}

	err := builder.Validate(config)
	assert.NoError(t, err)
}

func TestKafkaConnectToolBuilder_Validate_NoRequiredFeatures(t *testing.T) {
	builder := NewKafkaConnectToolBuilder()

	config := builders.ToolBuildConfig{
		ReadOnly: false,
		Features: []string{"other-feature"},
	}

	err := builder.Validate(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "none of required features found")
}

func TestKafkaConnectToolBuilder_ToolDefinition(t *testing.T) {
	builder := NewKafkaConnectToolBuilder()

	tool, err := builder.buildKafkaConnectTool()
	require.NoError(t, err)
	assert.Equal(t, "kafka_admin_connect", tool.Name)
	assert.NotEmpty(t, tool.Description)

	schema, ok := tool.InputSchema.(*jsonschema.Schema)
	require.True(t, ok)
	require.NotNil(t, schema.Properties)

	expectedRequired := []string{"resource", "operation"}
	assert.ElementsMatch(t, expectedRequired, schema.Required)

	expectedProps := []string{"resource", "operation", "name", "config"}
	assert.ElementsMatch(t, expectedProps, mapStringKeys(schema.Properties))

	resourceSchema := schema.Properties["resource"]
	require.NotNil(t, resourceSchema)
	assert.Equal(t, kafkaConnectResourceDesc, resourceSchema.Description)

	operationSchema := schema.Properties["operation"]
	require.NotNil(t, operationSchema)
	assert.Equal(t, kafkaConnectOperationDesc, operationSchema.Description)
}

func TestKafkaConnectToolBuilder_ReadOnlyRejectsWrite(t *testing.T) {
	builder := NewKafkaConnectToolBuilder()
	handler := builder.buildKafkaConnectHandler(true)

	_, _, err := handler(context.Background(), nil, kafkaConnectInput{
		Resource:  "connector",
		Operation: "create",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "read-only")
}
