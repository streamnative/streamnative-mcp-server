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

func TestPulsarAdminResourceQuotasToolBuilder(t *testing.T) {
	builder := NewPulsarAdminResourceQuotasToolBuilder()

	t.Run("GetName", func(t *testing.T) {
		assert.Equal(t, "pulsar_admin_resourcequotas", builder.GetName())
	})

	t.Run("GetRequiredFeatures", func(t *testing.T) {
		features := builder.GetRequiredFeatures()
		assert.NotEmpty(t, features)
		assert.Contains(t, features, "pulsar-admin-resourcequotas")
	})

	t.Run("BuildTools_Success", func(t *testing.T) {
		config := builders.ToolBuildConfig{
			ReadOnly: false,
			Features: []string{"pulsar-admin-resourcequotas"},
		}

		tools, err := builder.BuildTools(context.Background(), config)
		require.NoError(t, err)
		assert.Len(t, tools, 1)
		assert.Equal(t, "pulsar_admin_resourcequota", tools[0].Definition().Name)
	})

	t.Run("BuildTools_ReadOnly", func(t *testing.T) {
		config := builders.ToolBuildConfig{
			ReadOnly: true,
			Features: []string{"pulsar-admin-resourcequotas"},
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
			Features: []string{"pulsar-admin-resourcequotas"},
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

func TestPulsarAdminResourceQuotasToolSchema(t *testing.T) {
	builder := NewPulsarAdminResourceQuotasToolBuilder()
	tool, err := builder.buildResourceQuotasTool()
	require.NoError(t, err)
	assert.Equal(t, "pulsar_admin_resourcequota", tool.Name)

	schema, ok := tool.InputSchema.(*jsonschema.Schema)
	require.True(t, ok)
	require.NotNil(t, schema.Properties)

	expectedRequired := []string{"resource", "operation"}
	assert.ElementsMatch(t, expectedRequired, schema.Required)

	expectedProps := []string{
		"resource",
		"operation",
		"namespace",
		"bundle",
		"msgRateIn",
		"msgRateOut",
		"bandwidthIn",
		"bandwidthOut",
		"memory",
		"dynamic",
	}
	assert.ElementsMatch(t, expectedProps, mapStringKeys(schema.Properties))

	resourceSchema := schema.Properties["resource"]
	require.NotNil(t, resourceSchema)
	assert.Equal(t, pulsarAdminResourceQuotasResourceDesc, resourceSchema.Description)

	operationSchema := schema.Properties["operation"]
	require.NotNil(t, operationSchema)
	assert.Equal(t, pulsarAdminResourceQuotasOperationDesc, operationSchema.Description)

	namespaceSchema := schema.Properties["namespace"]
	require.NotNil(t, namespaceSchema)
	assert.Equal(t, pulsarAdminResourceQuotasNamespaceDesc, namespaceSchema.Description)

	bundleSchema := schema.Properties["bundle"]
	require.NotNil(t, bundleSchema)
	assert.Equal(t, pulsarAdminResourceQuotasBundleDesc, bundleSchema.Description)

	msgRateInSchema := schema.Properties["msgRateIn"]
	require.NotNil(t, msgRateInSchema)
	assert.Equal(t, pulsarAdminResourceQuotasMsgRateInDesc, msgRateInSchema.Description)

	msgRateOutSchema := schema.Properties["msgRateOut"]
	require.NotNil(t, msgRateOutSchema)
	assert.Equal(t, pulsarAdminResourceQuotasMsgRateOutDesc, msgRateOutSchema.Description)

	bandwidthInSchema := schema.Properties["bandwidthIn"]
	require.NotNil(t, bandwidthInSchema)
	assert.Equal(t, pulsarAdminResourceQuotasBandwidthInDesc, bandwidthInSchema.Description)

	bandwidthOutSchema := schema.Properties["bandwidthOut"]
	require.NotNil(t, bandwidthOutSchema)
	assert.Equal(t, pulsarAdminResourceQuotasBandwidthOutDesc, bandwidthOutSchema.Description)

	memorySchema := schema.Properties["memory"]
	require.NotNil(t, memorySchema)
	assert.Equal(t, pulsarAdminResourceQuotasMemoryDesc, memorySchema.Description)

	dynamicSchema := schema.Properties["dynamic"]
	require.NotNil(t, dynamicSchema)
	assert.Equal(t, pulsarAdminResourceQuotasDynamicDesc, dynamicSchema.Description)
}

func TestPulsarAdminResourceQuotasToolBuilder_ReadOnlyRejectsWrite(t *testing.T) {
	builder := NewPulsarAdminResourceQuotasToolBuilder()
	handler := builder.buildResourceQuotasHandler(true)

	_, _, err := handler(context.Background(), nil, pulsarAdminResourceQuotasInput{
		Resource:  "quota",
		Operation: "set",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "read-only")
}
