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

func TestPulsarClientConsumeToolBuilder(t *testing.T) {
	builder := NewPulsarClientConsumeToolBuilder()

	t.Run("GetName", func(t *testing.T) {
		assert.Equal(t, "pulsar_client_consume", builder.GetName())
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
		assert.Equal(t, "pulsar_client_consume", tools[0].Definition().Name)
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

	t.Run("Handler_MissingTopic", func(t *testing.T) {
		handler := builder.buildConsumeHandler(false)

		_, _, err := handler(context.Background(), nil, pulsarClientConsumeInput{
			SubscriptionName: "sub",
		})
		require.Error(t, err)
		assert.Equal(t, "failed to get topic: topic is required", err.Error())
	})

	t.Run("Handler_MissingSubscriptionName", func(t *testing.T) {
		handler := builder.buildConsumeHandler(false)

		_, _, err := handler(context.Background(), nil, pulsarClientConsumeInput{
			Topic: "persistent://tenant/ns/topic",
		})
		require.Error(t, err)
		assert.Equal(t, "failed to get subscription name: subscription-name is required", err.Error())
	})

	t.Run("Handler_InvalidSubscriptionTypeMissingSession", func(t *testing.T) {
		handler := builder.buildConsumeHandler(false)
		invalidType := "unknown"

		_, _, err := handler(context.Background(), nil, pulsarClientConsumeInput{
			Topic:            "persistent://tenant/ns/topic",
			SubscriptionName: "sub",
			SubscriptionType: &invalidType,
		})
		require.Error(t, err)
		assert.Equal(t, "pulsar session not found in context", err.Error())
	})

	t.Run("Handler_SessionMissing", func(t *testing.T) {
		handler := builder.buildConsumeHandler(false)

		_, _, err := handler(context.Background(), nil, pulsarClientConsumeInput{
			Topic:            "persistent://tenant/ns/topic",
			SubscriptionName: "sub",
		})
		require.Error(t, err)
		assert.Equal(t, "pulsar session not found in context", err.Error())
	})
}

func TestPulsarClientConsumeToolSchema(t *testing.T) {
	builder := NewPulsarClientConsumeToolBuilder()
	tool, err := builder.buildConsumeTool()
	require.NoError(t, err)
	assert.Equal(t, "pulsar_client_consume", tool.Name)

	schema, ok := tool.InputSchema.(*jsonschema.Schema)
	require.True(t, ok)
	require.NotNil(t, schema.Properties)

	expectedRequired := []string{"topic", "subscription-name"}
	assert.ElementsMatch(t, expectedRequired, schema.Required)

	expectedProps := []string{
		"topic",
		"subscription-name",
		"subscription-type",
		"subscription-mode",
		"initial-position",
		"num-messages",
		"timeout",
		"show-properties",
		"hide-payload",
	}
	assert.ElementsMatch(t, expectedProps, mapStringKeys(schema.Properties))

	topicSchema := schema.Properties["topic"]
	require.NotNil(t, topicSchema)
	assert.Equal(t, pulsarClientConsumeTopicDesc, topicSchema.Description)

	subscriptionNameSchema := schema.Properties["subscription-name"]
	require.NotNil(t, subscriptionNameSchema)
	assert.Equal(t, pulsarClientConsumeSubscriptionNameDesc, subscriptionNameSchema.Description)

	subscriptionTypeSchema := schema.Properties["subscription-type"]
	require.NotNil(t, subscriptionTypeSchema)
	assert.Equal(t, pulsarClientConsumeSubscriptionTypeDesc, subscriptionTypeSchema.Description)
}
