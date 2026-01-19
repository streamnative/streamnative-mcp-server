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

func TestPulsarAdminBrokersToolBuilder(t *testing.T) {
	builder := NewPulsarAdminBrokersToolBuilder()

	config := builders.ToolBuildConfig{
		ReadOnly: false,
		Features: []string{"pulsar-admin-brokers"},
	}

	tools, err := builder.BuildTools(context.Background(), config)
	require.NoError(t, err)
	assert.Len(t, tools, 1)
	assert.Equal(t, "pulsar_admin_brokers", tools[0].Definition().Name)

	config.Features = []string{"unrelated-feature"}
	tools, err = builder.BuildTools(context.Background(), config)
	require.NoError(t, err)
	assert.Empty(t, tools)
}

func TestPulsarAdminBrokersToolSchema(t *testing.T) {
	builder := NewPulsarAdminBrokersToolBuilder()
	tool, err := builder.buildPulsarAdminBrokersTool()
	require.NoError(t, err)
	assert.Equal(t, "pulsar_admin_brokers", tool.Name)

	schema, ok := tool.InputSchema.(*jsonschema.Schema)
	require.True(t, ok)
	require.NotNil(t, schema.Properties)

	expectedRequired := []string{"resource", "operation"}
	assert.ElementsMatch(t, expectedRequired, schema.Required)

	expectedProps := []string{
		"resource",
		"operation",
		"clusterName",
		"brokerUrl",
		"configType",
		"configName",
		"configValue",
	}
	assert.ElementsMatch(t, expectedProps, mapStringKeys(schema.Properties))

	resourceSchema := schema.Properties["resource"]
	require.NotNil(t, resourceSchema)
	assert.Equal(t, pulsarAdminBrokersResourceDesc, resourceSchema.Description)
}

func TestPulsarAdminBrokersToolBuilder_RequiresSession(t *testing.T) {
	builder := NewPulsarAdminBrokersToolBuilder()
	handler := builder.buildPulsarAdminBrokersHandler(true)

	_, _, err := handler(context.Background(), nil, pulsarAdminBrokersInput{
		Resource:  "config",
		Operation: "update",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "pulsar session")
}

func TestPulsarAdminClusterToolBuilder(t *testing.T) {
	builder := NewPulsarAdminClusterToolBuilder()

	config := builders.ToolBuildConfig{
		ReadOnly: false,
		Features: []string{"pulsar-admin-clusters"},
	}

	tools, err := builder.BuildTools(context.Background(), config)
	require.NoError(t, err)
	assert.Len(t, tools, 1)
	assert.Equal(t, "pulsar_admin_cluster", tools[0].Definition().Name)
}

func TestPulsarAdminClusterToolSchema(t *testing.T) {
	builder := NewPulsarAdminClusterToolBuilder()
	tool, err := builder.buildClusterTool()
	require.NoError(t, err)
	assert.Equal(t, "pulsar_admin_cluster", tool.Name)

	schema, ok := tool.InputSchema.(*jsonschema.Schema)
	require.True(t, ok)
	require.NotNil(t, schema.Properties)

	expectedRequired := []string{"resource", "operation"}
	assert.ElementsMatch(t, expectedRequired, schema.Required)

	expectedProps := []string{
		"resource",
		"operation",
		"cluster_name",
		"domain_name",
		"service_url",
		"service_url_tls",
		"broker_service_url",
		"broker_service_url_tls",
		"peer_cluster_names",
		"brokers",
	}
	assert.ElementsMatch(t, expectedProps, mapStringKeys(schema.Properties))
}

func TestPulsarAdminClusterToolBuilder_RequiresSession(t *testing.T) {
	builder := NewPulsarAdminClusterToolBuilder()
	handler := builder.buildClusterHandler(true)

	_, _, err := handler(context.Background(), nil, pulsarAdminClusterInput{
		Resource:  "cluster",
		Operation: "create",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "pulsar session")
}

func TestPulsarAdminSourcesToolBuilder(t *testing.T) {
	builder := NewPulsarAdminSourcesToolBuilder()

	config := builders.ToolBuildConfig{
		ReadOnly: false,
		Features: []string{"pulsar-admin-sources"},
	}

	tools, err := builder.BuildTools(context.Background(), config)
	require.NoError(t, err)
	assert.Len(t, tools, 1)
	assert.Equal(t, "pulsar_admin_sources", tools[0].Definition().Name)
}

func TestPulsarAdminSourcesToolSchema(t *testing.T) {
	builder := NewPulsarAdminSourcesToolBuilder()
	tool, err := builder.buildSourcesTool()
	require.NoError(t, err)
	assert.Equal(t, "pulsar_admin_sources", tool.Name)

	schema, ok := tool.InputSchema.(*jsonschema.Schema)
	require.True(t, ok)
	require.NotNil(t, schema.Properties)

	expectedRequired := []string{"operation"}
	assert.ElementsMatch(t, expectedRequired, schema.Required)

	expectedProps := []string{
		"operation",
		"tenant",
		"namespace",
		"name",
		"archive",
		"source-type",
		"destination-topic-name",
		"deserialization-classname",
		"schema-type",
		"classname",
		"processing-guarantees",
		"parallelism",
		"source-config",
	}
	assert.ElementsMatch(t, expectedProps, mapStringKeys(schema.Properties))
}

func TestPulsarAdminSourcesToolBuilder_ReadOnlyRejectsCreate(t *testing.T) {
	builder := NewPulsarAdminSourcesToolBuilder()
	handler := builder.buildSourcesHandler(true)

	_, _, err := handler(context.Background(), nil, pulsarAdminSourcesInput{
		Operation: "create",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "read-only")
}

func TestPulsarAdminSinksToolBuilder(t *testing.T) {
	builder := NewPulsarAdminSinksToolBuilder()

	config := builders.ToolBuildConfig{
		ReadOnly: false,
		Features: []string{"pulsar-admin-sinks"},
	}

	tools, err := builder.BuildTools(context.Background(), config)
	require.NoError(t, err)
	assert.Len(t, tools, 1)
	assert.Equal(t, "pulsar_admin_sinks", tools[0].Definition().Name)
}

func TestPulsarAdminSinksToolSchema(t *testing.T) {
	builder := NewPulsarAdminSinksToolBuilder()
	tool, err := builder.buildSinksTool()
	require.NoError(t, err)
	assert.Equal(t, "pulsar_admin_sinks", tool.Name)

	schema, ok := tool.InputSchema.(*jsonschema.Schema)
	require.True(t, ok)
	require.NotNil(t, schema.Properties)

	expectedRequired := []string{"operation"}
	assert.ElementsMatch(t, expectedRequired, schema.Required)

	expectedProps := []string{
		"operation",
		"tenant",
		"namespace",
		"name",
		"archive",
		"sink-type",
		"inputs",
		"topics-pattern",
		"subs-name",
		"parallelism",
		"sink-config",
	}
	assert.ElementsMatch(t, expectedProps, mapStringKeys(schema.Properties))
}

func TestPulsarAdminSinksToolBuilder_ReadOnlyRejectsCreate(t *testing.T) {
	builder := NewPulsarAdminSinksToolBuilder()
	handler := builder.buildSinksHandler(true)

	_, _, err := handler(context.Background(), nil, pulsarAdminSinksInput{
		Operation: "create",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "read-only")
}

func TestPulsarAdminPackagesToolBuilder(t *testing.T) {
	builder := NewPulsarAdminPackagesToolBuilder()

	config := builders.ToolBuildConfig{
		ReadOnly: false,
		Features: []string{"pulsar-admin-packages"},
	}

	tools, err := builder.BuildTools(context.Background(), config)
	require.NoError(t, err)
	assert.Len(t, tools, 1)
	assert.Equal(t, "pulsar_admin_package", tools[0].Definition().Name)

	config.Features = []string{"unrelated-feature"}
	tools, err = builder.BuildTools(context.Background(), config)
	require.NoError(t, err)
	assert.Empty(t, tools)
}

func TestPulsarAdminPackagesToolSchema(t *testing.T) {
	builder := NewPulsarAdminPackagesToolBuilder()
	tool, err := builder.buildPackagesTool()
	require.NoError(t, err)
	assert.Equal(t, "pulsar_admin_package", tool.Name)

	schema, ok := tool.InputSchema.(*jsonschema.Schema)
	require.True(t, ok)
	require.NotNil(t, schema.Properties)

	expectedRequired := []string{"resource", "operation"}
	assert.ElementsMatch(t, expectedRequired, schema.Required)

	expectedProps := []string{
		"resource",
		"operation",
		"packageName",
		"namespace",
		"type",
		"description",
		"contact",
		"path",
		"properties",
	}
	assert.ElementsMatch(t, expectedProps, mapStringKeys(schema.Properties))

	resourceSchema := schema.Properties["resource"]
	require.NotNil(t, resourceSchema)
	assert.Equal(t, pulsarAdminPackagesResourceDesc, resourceSchema.Description)
}

func TestPulsarAdminPackagesToolBuilder_ReadOnlyRejectsUpload(t *testing.T) {
	builder := NewPulsarAdminPackagesToolBuilder()
	handler := builder.buildPackagesHandler(true)

	_, _, err := handler(context.Background(), nil, pulsarAdminPackagesInput{
		Resource:  "package",
		Operation: "upload",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "read-only")
}

func TestPulsarAdminSubscriptionToolBuilder(t *testing.T) {
	builder := NewPulsarAdminSubscriptionToolBuilder()

	config := builders.ToolBuildConfig{
		ReadOnly: false,
		Features: []string{"pulsar-admin-subscriptions"},
	}

	tools, err := builder.BuildTools(context.Background(), config)
	require.NoError(t, err)
	assert.Len(t, tools, 1)
	assert.Equal(t, "pulsar_admin_subscription", tools[0].Definition().Name)
}

func TestPulsarAdminSubscriptionToolSchema(t *testing.T) {
	builder := NewPulsarAdminSubscriptionToolBuilder()
	tool, err := builder.buildSubscriptionTool()
	require.NoError(t, err)
	assert.Equal(t, "pulsar_admin_subscription", tool.Name)

	schema, ok := tool.InputSchema.(*jsonschema.Schema)
	require.True(t, ok)
	require.NotNil(t, schema.Properties)

	expectedRequired := []string{"resource", "operation", "topic"}
	assert.ElementsMatch(t, expectedRequired, schema.Required)

	expectedProps := []string{
		"resource",
		"operation",
		"topic",
		"subscription",
		"messageId",
		"count",
		"expireTimeInSeconds",
		"force",
	}
	assert.ElementsMatch(t, expectedProps, mapStringKeys(schema.Properties))
}

func TestPulsarAdminSubscriptionToolBuilder_ReadOnlyRejectsDelete(t *testing.T) {
	builder := NewPulsarAdminSubscriptionToolBuilder()
	handler := builder.buildSubscriptionHandler(true)

	_, _, err := handler(context.Background(), nil, pulsarAdminSubscriptionInput{
		Resource:  "subscription",
		Operation: "delete",
		Topic:     "persistent://public/default/test",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "read-only")
}
