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

func TestPulsarAdminTopicPolicyToolBuilder(t *testing.T) {
	builder := NewPulsarAdminTopicPolicyToolBuilder()

	t.Run("GetName", func(t *testing.T) {
		assert.Equal(t, "pulsar_admin_topic_policy", builder.GetName())
	})

	t.Run("GetRequiredFeatures", func(t *testing.T) {
		features := builder.GetRequiredFeatures()
		assert.NotEmpty(t, features)
		assert.Contains(t, features, "pulsar-admin-topic-policy")
	})

	t.Run("BuildTools_Success", func(t *testing.T) {
		config := builders.ToolBuildConfig{
			ReadOnly: false,
			Features: []string{"pulsar-admin-topic-policy"},
		}

		tools, err := builder.BuildTools(context.Background(), config)
		require.NoError(t, err)
		assert.Len(t, tools, 1)
		assert.Equal(t, "pulsar_admin_topic_policy", tools[0].Definition().Name)
	})

	t.Run("BuildTools_ReadOnly", func(t *testing.T) {
		config := builders.ToolBuildConfig{
			ReadOnly: true,
			Features: []string{"pulsar-admin-topic-policy"},
		}

		tools, err := builder.BuildTools(context.Background(), config)
		require.NoError(t, err)
		assert.Len(t, tools, 1)
		assert.Equal(t, "pulsar_admin_topic_policy", tools[0].Definition().Name)
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
			Features: []string{"pulsar-admin-topic-policy"},
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

func TestPulsarAdminTopicPolicyToolSchema(t *testing.T) {
	builder := NewPulsarAdminTopicPolicyToolBuilder()

	tool, err := builder.buildTopicPolicyTool()
	require.NoError(t, err)
	assert.Equal(t, "pulsar_admin_topic_policy", tool.Name)

	schema, ok := tool.InputSchema.(*jsonschema.Schema)
	require.True(t, ok)
	require.NotNil(t, schema.Properties)

	assert.ElementsMatch(t, []string{"operation", "topic"}, schema.Required)
	assert.ElementsMatch(t, []string{
		"operation",
		"topic",
		"retention_size",
		"retention_time",
		"ttl_seconds",
		"compaction_threshold",
		"subscription_types",
	}, mapStringKeys(schema.Properties))

	assert.Equal(t, pulsarAdminTopicPolicyOperationDesc, schema.Properties["operation"].Description)
	assert.Equal(t, pulsarAdminTopicPolicyTopicDesc, schema.Properties["topic"].Description)
	assert.Equal(t, pulsarAdminTopicPolicyRetentionTimeDesc, schema.Properties["retention_time"].Description)
	assert.Equal(t, pulsarAdminTopicPolicyRetentionSizeDesc, schema.Properties["retention_size"].Description)
}
