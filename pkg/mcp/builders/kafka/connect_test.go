// Copyright 2026 StreamNative
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
	require.Len(t, tools, 2)

	tool := tools[0]
	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "kafka_admin_connect_read", tool.Tool.Name)
	assert.Equal(t, "kafka_admin_connect_write", tools[1].Tool.Name)
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
	assert.Equal(t, "kafka_admin_connect_read", tool.Tool.Name)
}

func TestKafkaConnectToolBuilder_BuildTools_WithAllKafkaFeature(t *testing.T) {
	builder := NewKafkaConnectToolBuilder()

	config := builders.ToolBuildConfig{
		ReadOnly: false,
		Features: []string{"all-kafka"},
	}

	tools, err := builder.BuildTools(context.Background(), config)

	require.NoError(t, err)
	require.Len(t, tools, 2)
}

func TestKafkaConnectToolBuilder_BuildTools_WithKafkaAdminFeature(t *testing.T) {
	builder := NewKafkaConnectToolBuilder()

	config := builders.ToolBuildConfig{
		ReadOnly: false,
		Features: []string{"kafka-admin"},
	}

	tools, err := builder.BuildTools(context.Background(), config)

	require.NoError(t, err)
	require.Len(t, tools, 2)
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

	tool := builder.buildKafkaConnectTool(toolModeRead)

	assert.Equal(t, "kafka_admin_connect_read", tool.Name)
	assert.NotEmpty(t, tool.Description)

	// Verify required parameters
	assert.Contains(t, tool.InputSchema.Properties, "resource")
	assert.Contains(t, tool.InputSchema.Properties, "operation")
	assert.Contains(t, tool.InputSchema.Properties, "name")
	assert.NotContains(t, tool.InputSchema.Properties, "config")

	writeTool := builder.buildKafkaConnectTool(toolModeWrite)
	assert.Equal(t, "kafka_admin_connect_write", writeTool.Name)
	assert.Contains(t, writeTool.InputSchema.Properties, "config")

	// Verify required fields
	assert.Contains(t, tool.InputSchema.Required, "resource")
	assert.Contains(t, tool.InputSchema.Required, "operation")
}
