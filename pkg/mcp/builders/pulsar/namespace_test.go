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

func TestPulsarAdminNamespaceToolBuilder(t *testing.T) {
	builder := NewPulsarAdminNamespaceToolBuilder()

	t.Run("GetName", func(t *testing.T) {
		assert.Equal(t, "pulsar_admin_namespace", builder.GetName())
	})

	t.Run("GetRequiredFeatures", func(t *testing.T) {
		features := builder.GetRequiredFeatures()
		assert.NotEmpty(t, features)
		assert.Contains(t, features, "pulsar-admin-namespaces")
	})

	t.Run("BuildTools_Success", func(t *testing.T) {
		config := builders.ToolBuildConfig{
			ReadOnly: false,
			Features: []string{"pulsar-admin-namespaces"},
		}

		tools, err := builder.BuildTools(context.Background(), config)
		require.NoError(t, err)
		assert.Len(t, tools, 1)
		assert.Equal(t, "pulsar_admin_namespace", tools[0].Definition().Name)
		assert.NotNil(t, tools[0])
	})

	t.Run("BuildTools_ReadOnly", func(t *testing.T) {
		config := builders.ToolBuildConfig{
			ReadOnly: true,
			Features: []string{"pulsar-admin-namespaces"},
		}

		tools, err := builder.BuildTools(context.Background(), config)
		require.NoError(t, err)
		assert.Len(t, tools, 1)
		assert.Equal(t, "pulsar_admin_namespace", tools[0].Definition().Name)
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
			Features: []string{"pulsar-admin-namespaces"},
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

func TestPulsarAdminNamespaceToolSchema(t *testing.T) {
	builder := NewPulsarAdminNamespaceToolBuilder()
	tool, err := builder.buildNamespaceTool()
	require.NoError(t, err)
	assert.Equal(t, "pulsar_admin_namespace", tool.Name)

	schema, ok := tool.InputSchema.(*jsonschema.Schema)
	require.True(t, ok)
	require.NotNil(t, schema.Properties)

	expectedRequired := []string{"operation"}
	assert.ElementsMatch(t, expectedRequired, schema.Required)

	expectedProps := []string{
		"operation",
		"tenant",
		"namespace",
		"bundles",
		"clusters",
		"subscription",
		"bundle",
		"force",
		"unload",
	}
	assert.ElementsMatch(t, expectedProps, mapStringKeys(schema.Properties))

	operationSchema := schema.Properties["operation"]
	require.NotNil(t, operationSchema)
	assert.Equal(t, pulsarAdminNamespaceOperationDesc, operationSchema.Description)

	tenantSchema := schema.Properties["tenant"]
	require.NotNil(t, tenantSchema)
	assert.Equal(t, pulsarAdminNamespaceTenantDesc, tenantSchema.Description)
}

func TestPulsarAdminNamespaceToolBuilder_ReadOnlyRejectsWrite(t *testing.T) {
	builder := NewPulsarAdminNamespaceToolBuilder()
	handler := builder.buildNamespaceHandler(true)

	_, _, err := handler(context.Background(), nil, pulsarAdminNamespaceInput{
		Operation: "create",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "read-only")
}
