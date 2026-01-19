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
	"github.com/sirupsen/logrus"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKafkaConsumeToolBuilder(t *testing.T) {
	builder := NewKafkaConsumeToolBuilder()

	t.Run("GetName", func(t *testing.T) {
		assert.Equal(t, "kafka_consume", builder.GetName())
	})

	t.Run("GetRequiredFeatures", func(t *testing.T) {
		features := builder.GetRequiredFeatures()
		assert.NotEmpty(t, features)
		assert.Contains(t, features, "kafka-client")
	})

	t.Run("BuildTools_Success", func(t *testing.T) {
		config := builders.ToolBuildConfig{
			ReadOnly: false,
			Features: []string{"kafka-client"},
		}

		tools, err := builder.BuildTools(context.Background(), config)
		require.NoError(t, err)
		assert.Len(t, tools, 1)
		assert.Equal(t, "kafka_client_consume", tools[0].Definition().Name)
		assert.NotNil(t, tools[0])
	})

	t.Run("BuildTools_WithLogger", func(t *testing.T) {
		logger := logrus.New()
		config := builders.ToolBuildConfig{
			ReadOnly: false,
			Features: []string{"kafka-client"},
			Options: map[string]interface{}{
				"logger": logger,
			},
		}

		tools, err := builder.BuildTools(context.Background(), config)
		require.NoError(t, err)
		assert.Len(t, tools, 1)
		assert.Equal(t, logger, builder.logger)
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
			Features: []string{"kafka-client"},
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

func TestKafkaConsumeToolSchema(t *testing.T) {
	builder := NewKafkaConsumeToolBuilder()
	tool, err := builder.buildKafkaConsumeTool()
	require.NoError(t, err)
	assert.Equal(t, "kafka_client_consume", tool.Name)

	schema, ok := tool.InputSchema.(*jsonschema.Schema)
	require.True(t, ok)
	require.NotNil(t, schema.Properties)

	expectedRequired := []string{"topic"}
	assert.ElementsMatch(t, expectedRequired, schema.Required)

	expectedProps := []string{
		"topic",
		"group",
		"offset",
		"max-messages",
		"timeout",
	}
	assert.ElementsMatch(t, expectedProps, mapStringKeys(schema.Properties))

	topicSchema := schema.Properties["topic"]
	require.NotNil(t, topicSchema)
	assert.Equal(t, kafkaConsumeTopicDesc, topicSchema.Description)
}
