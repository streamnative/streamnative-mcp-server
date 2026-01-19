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

func TestPulsarAdminNsIsolationPolicyToolBuilder(t *testing.T) {
	builder := NewPulsarAdminNsIsolationPolicyToolBuilder()

	t.Run("GetName", func(t *testing.T) {
		assert.Equal(t, "pulsar_admin_nsisolationpolicy", builder.GetName())
	})

	t.Run("GetRequiredFeatures", func(t *testing.T) {
		features := builder.GetRequiredFeatures()
		assert.NotEmpty(t, features)
		assert.Contains(t, features, "pulsar-admin-nsisolationpolicy")
	})

	t.Run("BuildTools_Success", func(t *testing.T) {
		config := builders.ToolBuildConfig{
			ReadOnly: false,
			Features: []string{"pulsar-admin-nsisolationpolicy"},
		}

		tools, err := builder.BuildTools(context.Background(), config)
		require.NoError(t, err)
		assert.Len(t, tools, 1)
		assert.Equal(t, "pulsar_admin_nsisolationpolicy", tools[0].Definition().Name)
	})

	t.Run("BuildTools_ReadOnly", func(t *testing.T) {
		config := builders.ToolBuildConfig{
			ReadOnly: true,
			Features: []string{"pulsar-admin-nsisolationpolicy"},
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
			Features: []string{"pulsar-admin-nsisolationpolicy"},
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

func TestPulsarAdminNsIsolationPolicyToolSchema(t *testing.T) {
	builder := NewPulsarAdminNsIsolationPolicyToolBuilder()
	tool, err := builder.buildNsIsolationPolicyTool()
	require.NoError(t, err)
	assert.Equal(t, "pulsar_admin_nsisolationpolicy", tool.Name)

	schema, ok := tool.InputSchema.(*jsonschema.Schema)
	require.True(t, ok)
	require.NotNil(t, schema.Properties)

	expectedRequired := []string{"resource", "operation", "cluster"}
	assert.ElementsMatch(t, expectedRequired, schema.Required)

	expectedProps := []string{
		"resource",
		"operation",
		"cluster",
		"name",
		"namespaces",
		"primary",
		"secondary",
		"autoFailoverPolicyType",
		"autoFailoverPolicyParams",
	}
	assert.ElementsMatch(t, expectedProps, mapStringKeys(schema.Properties))

	resourceSchema := schema.Properties["resource"]
	require.NotNil(t, resourceSchema)
	assert.Equal(t, pulsarAdminNsIsolationPolicyResourceDesc, resourceSchema.Description)

	operationSchema := schema.Properties["operation"]
	require.NotNil(t, operationSchema)
	assert.Equal(t, pulsarAdminNsIsolationPolicyOperationDesc, operationSchema.Description)

	clusterSchema := schema.Properties["cluster"]
	require.NotNil(t, clusterSchema)
	assert.Equal(t, pulsarAdminNsIsolationPolicyClusterDesc, clusterSchema.Description)

	nameSchema := schema.Properties["name"]
	require.NotNil(t, nameSchema)
	assert.Equal(t, pulsarAdminNsIsolationPolicyNameDesc, nameSchema.Description)

	namespacesSchema := schema.Properties["namespaces"]
	require.NotNil(t, namespacesSchema)
	assert.Equal(t, pulsarAdminNsIsolationPolicyNamespacesDesc, namespacesSchema.Description)
	require.NotNil(t, namespacesSchema.Items)
	assert.Equal(t, "namespace", namespacesSchema.Items.Description)

	primarySchema := schema.Properties["primary"]
	require.NotNil(t, primarySchema)
	assert.Equal(t, pulsarAdminNsIsolationPolicyPrimaryDesc, primarySchema.Description)
	require.NotNil(t, primarySchema.Items)
	assert.Equal(t, "primary broker", primarySchema.Items.Description)

	secondarySchema := schema.Properties["secondary"]
	require.NotNil(t, secondarySchema)
	assert.Equal(t, pulsarAdminNsIsolationPolicySecondaryDesc, secondarySchema.Description)
	require.NotNil(t, secondarySchema.Items)
	assert.Equal(t, "secondary broker", secondarySchema.Items.Description)

	autoFailoverTypeSchema := schema.Properties["autoFailoverPolicyType"]
	require.NotNil(t, autoFailoverTypeSchema)
	assert.Equal(t, pulsarAdminNsIsolationPolicyTypeDesc, autoFailoverTypeSchema.Description)

	autoFailoverParamsSchema := schema.Properties["autoFailoverPolicyParams"]
	require.NotNil(t, autoFailoverParamsSchema)
	assert.Equal(t, pulsarAdminNsIsolationPolicyParamsDesc, autoFailoverParamsSchema.Description)
}
