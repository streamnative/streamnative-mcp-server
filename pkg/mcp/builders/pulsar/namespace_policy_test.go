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

func TestPulsarAdminNamespacePolicyToolBuilder(t *testing.T) {
	builder := NewPulsarAdminNamespacePolicyToolBuilder()

	t.Run("GetName", func(t *testing.T) {
		assert.Equal(t, "pulsar_admin_namespace_policy", builder.GetName())
	})

	t.Run("GetRequiredFeatures", func(t *testing.T) {
		features := builder.GetRequiredFeatures()
		assert.NotEmpty(t, features)
		assert.Contains(t, features, "pulsar-admin-namespace-policy")
	})

	t.Run("BuildTools_Success", func(t *testing.T) {
		config := builders.ToolBuildConfig{
			ReadOnly: false,
			Features: []string{"pulsar-admin-namespace-policy"},
		}

		tools, err := builder.BuildTools(context.Background(), config)
		require.NoError(t, err)
		assert.Len(t, tools, 3)

		names := make([]string, 0, len(tools))
		for _, tool := range tools {
			names = append(names, tool.Definition().Name)
		}

		assert.ElementsMatch(t, []string{
			"pulsar_admin_namespace_policy_get",
			"pulsar_admin_namespace_policy_set",
			"pulsar_admin_namespace_policy_remove",
		}, names)
	})

	t.Run("BuildTools_ReadOnly", func(t *testing.T) {
		config := builders.ToolBuildConfig{
			ReadOnly: true,
			Features: []string{"pulsar-admin-namespace-policy"},
		}

		tools, err := builder.BuildTools(context.Background(), config)
		require.NoError(t, err)
		assert.Len(t, tools, 1)
		assert.Equal(t, "pulsar_admin_namespace_policy_get", tools[0].Definition().Name)
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
			Features: []string{"pulsar-admin-namespace-policy"},
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

func TestPulsarAdminNamespacePolicyToolSchema(t *testing.T) {
	builder := NewPulsarAdminNamespacePolicyToolBuilder()

	getTool, err := builder.buildNamespaceGetPoliciesTool()
	require.NoError(t, err)
	assert.Equal(t, "pulsar_admin_namespace_policy_get", getTool.Name)

	getSchema, ok := getTool.InputSchema.(*jsonschema.Schema)
	require.True(t, ok)
	require.NotNil(t, getSchema.Properties)

	assert.ElementsMatch(t, []string{"namespace"}, getSchema.Required)
	assert.ElementsMatch(t, []string{"namespace"}, mapStringKeys(getSchema.Properties))
	assert.Equal(t, pulsarAdminNamespacePolicyGetNamespaceDesc, getSchema.Properties["namespace"].Description)

	setTool, err := builder.buildNamespaceSetPolicyTool()
	require.NoError(t, err)
	assert.Equal(t, "pulsar_admin_namespace_policy_set", setTool.Name)

	setSchema, ok := setTool.InputSchema.(*jsonschema.Schema)
	require.True(t, ok)
	require.NotNil(t, setSchema.Properties)

	assert.ElementsMatch(t, []string{"namespace", "policy"}, setSchema.Required)
	assert.ElementsMatch(t, []string{
		"namespace",
		"policy",
		"role",
		"actions",
		"clusters",
		"roles",
		"ttl",
		"time",
		"size",
		"limit-size",
		"limit-time",
		"type",
	}, mapStringKeys(setSchema.Properties))

	assert.Equal(t, pulsarAdminNamespacePolicySetNamespaceDesc, setSchema.Properties["namespace"].Description)
	assert.Equal(t, pulsarAdminNamespacePolicySetPolicyDesc, setSchema.Properties["policy"].Description)
	assert.Equal(t, pulsarAdminNamespacePolicySetTypeDesc, setSchema.Properties["type"].Description)

	removeTool, err := builder.buildNamespaceRemovePolicyTool()
	require.NoError(t, err)
	assert.Equal(t, "pulsar_admin_namespace_policy_remove", removeTool.Name)

	removeSchema, ok := removeTool.InputSchema.(*jsonschema.Schema)
	require.True(t, ok)
	require.NotNil(t, removeSchema.Properties)

	assert.ElementsMatch(t, []string{"namespace", "policy"}, removeSchema.Required)
	assert.ElementsMatch(t, []string{
		"namespace",
		"policy",
		"role",
		"subscription",
		"type",
	}, mapStringKeys(removeSchema.Properties))

	assert.Equal(t, pulsarAdminNamespacePolicyRemoveNamespaceDesc, removeSchema.Properties["namespace"].Description)
	assert.Equal(t, pulsarAdminNamespacePolicyRemovePolicyDesc, removeSchema.Properties["policy"].Description)
}
