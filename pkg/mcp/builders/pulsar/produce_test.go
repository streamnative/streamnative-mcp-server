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

func TestPulsarClientProduceToolBuilder(t *testing.T) {
	builder := NewPulsarClientProduceToolBuilder()

	t.Run("GetName", func(t *testing.T) {
		assert.Equal(t, "pulsar_client_produce", builder.GetName())
	})

	t.Run("GetRequiredFeatures", func(t *testing.T) {
		features := builder.GetRequiredFeatures()
		assert.NotEmpty(t, features)
		assert.Contains(t, features, "pulsar-client")
	})

	t.Run("BuildTools_Success", func(t *testing.T) {
		config := builders.ToolBuildConfig{
			ReadOnly: false,
			Features: []string{"pulsar-client"},
		}

		tools, err := builder.BuildTools(context.Background(), config)
		require.NoError(t, err)
		assert.Len(t, tools, 1)
		assert.Equal(t, "pulsar_client_produce", tools[0].Definition().Name)
		assert.NotNil(t, tools[0])
	})

	t.Run("BuildTools_ReadOnlyMode", func(t *testing.T) {
		config := builders.ToolBuildConfig{
			ReadOnly: true,
			Features: []string{"pulsar-client"},
		}

		tools, err := builder.BuildTools(context.Background(), config)
		require.NoError(t, err)
		assert.Len(t, tools, 1)
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
			Features: []string{"pulsar-client"},
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

	t.Run("Handler_ReadOnlyRejects", func(t *testing.T) {
		handler := builder.buildProduceHandler(true)

		_, _, err := handler(context.Background(), nil, pulsarClientProduceInput{
			Topic:    "persistent://tenant/ns/topic",
			Messages: []string{"message"},
		})
		require.Error(t, err)
		assert.Equal(t, "message production is not allowed in read-only mode", err.Error())
	})

	t.Run("Handler_MissingMessages", func(t *testing.T) {
		handler := builder.buildProduceHandler(false)

		_, _, err := handler(context.Background(), nil, pulsarClientProduceInput{
			Topic: "persistent://tenant/ns/topic",
		})
		require.Error(t, err)
		assert.Equal(t, "please supply message content with 'messages' parameter", err.Error())
	})
}

func TestPulsarClientProduceToolSchema(t *testing.T) {
	builder := NewPulsarClientProduceToolBuilder()
	tool, err := builder.buildProduceTool()
	require.NoError(t, err)
	assert.Equal(t, "pulsar_client_produce", tool.Name)

	schema, ok := tool.InputSchema.(*jsonschema.Schema)
	require.True(t, ok)
	require.NotNil(t, schema.Properties)

	expectedRequired := []string{"topic"}
	assert.ElementsMatch(t, expectedRequired, schema.Required)

	expectedProps := []string{
		"topic",
		"messages",
		"num-produce",
		"rate",
		"disable-batching",
		"chunking",
		"separator",
		"properties",
		"key",
	}
	assert.ElementsMatch(t, expectedProps, mapStringKeys(schema.Properties))

	topicSchema := schema.Properties["topic"]
	require.NotNil(t, topicSchema)
	assert.Equal(t, pulsarClientProduceTopicDesc, topicSchema.Description)

	messagesSchema := schema.Properties["messages"]
	require.NotNil(t, messagesSchema)
	require.NotNil(t, messagesSchema.Items)
	assert.Equal(t, pulsarClientProduceMessageItemDesc, messagesSchema.Items.Description)

	propertiesSchema := schema.Properties["properties"]
	require.NotNil(t, propertiesSchema)
	require.NotNil(t, propertiesSchema.Items)
	assert.Equal(t, pulsarClientProducePropertyItemDesc, propertiesSchema.Items.Description)
}
