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
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	"github.com/stretchr/testify/require"
)

func TestKafkaOperationSpecsAreValid(t *testing.T) {
	for name, registry := range kafkaComplianceOperationSpecs() {
		require.NotPanics(t, registry.MustValidate, name)
		require.NotEmpty(t, registry.NamesForMode(builders.OperationModeRead), name)
		require.NotEmpty(t, registry.NamesForMode(builders.OperationModeWrite), name)
		require.ErrorContains(t, registry.ValidateModeOperation(builders.OperationModeRead, "__unknown__"), "unknown operation", name)
	}
}

func TestKafkaToolAnnotationCompliance(t *testing.T) {
	builderList := []builders.ToolBuilder{
		NewKafkaTopicsToolBuilder(),
		NewKafkaGroupsToolBuilder(),
		NewKafkaSchemaRegistryToolBuilder(),
		NewKafkaConnectToolBuilder(),
		NewKafkaPartitionsToolBuilder(),
		NewKafkaProduceToolBuilder(),
		NewKafkaConsumeToolBuilder(),
	}

	for _, builder := range builderList {
		tools, err := builder.BuildTools(context.Background(), builders.ToolBuildConfig{
			Features: []string{"all", "all-kafka", "kafka-admin", "kafka-admin-kafka-connect", "kafka-admin-schema-registry", "kafka-client"},
		})
		require.NoError(t, err)
		for _, serverTool := range tools {
			tool := serverTool.Tool
			require.NotEmpty(t, tool.Annotations.Title, tool.Name)
			require.NotNil(t, tool.Annotations.ReadOnlyHint, tool.Name)
			require.NotNil(t, tool.Annotations.DestructiveHint, tool.Name)
			require.NotNil(t, tool.Annotations.IdempotentHint, tool.Name)
			require.NotNil(t, tool.Annotations.OpenWorldHint, tool.Name)
			require.LessOrEqual(t, len(tool.Name), 64, tool.Name)

			isRead := strings.HasSuffix(tool.Name, "_read")
			isWrite := strings.HasSuffix(tool.Name, "_write") || strings.Contains(tool.Name, "produce") || strings.Contains(tool.Name, "consume")
			if isRead {
				require.True(t, *tool.Annotations.ReadOnlyHint, tool.Name)
				require.False(t, *tool.Annotations.DestructiveHint, tool.Name)
				require.True(t, *tool.Annotations.IdempotentHint, tool.Name)
				require.True(t, *tool.Annotations.OpenWorldHint, tool.Name)
			}
			if isWrite {
				require.False(t, *tool.Annotations.ReadOnlyHint, tool.Name)
				require.True(t, *tool.Annotations.DestructiveHint, tool.Name)
				require.False(t, *tool.Annotations.IdempotentHint, tool.Name)
				require.True(t, *tool.Annotations.OpenWorldHint, tool.Name)
			}
			assertOperationEnumMode(t, tool.Name, tool.InputSchema.Properties["operation"])
		}
	}
}

func toolPropertyNames(tool mcp.Tool) []string {
	names := make([]string, 0, len(tool.InputSchema.Properties))
	for name := range tool.InputSchema.Properties {
		names = append(names, name)
	}
	return names
}

func assertOperationEnumMode(t *testing.T, toolName string, operationSchema any) {
	t.Helper()
	registry, ok := kafkaComplianceOperationSpecs()[toolName]
	if !ok {
		return
	}
	rawEnum, ok := stringEnum(operationSchema)
	if !ok {
		return
	}
	seenRead, seenWrite := false, false
	for _, op := range rawEnum {
		spec, ok := registry.SpecFor(op)
		require.True(t, ok, "%s operation %q must be declared in OperationSpec", toolName, op)
		switch spec.Mode {
		case builders.OperationModeRead:
			seenRead = true
		case builders.OperationModeWrite:
			seenWrite = true
		default:
			require.Failf(t, "invalid operation mode", "%s operation %q has mode %q", toolName, op, spec.Mode)
		}
	}
	require.False(t, seenRead && seenWrite, toolName)
}

func kafkaComplianceOperationSpecs() map[string]builders.OperationRegistry {
	return map[string]builders.OperationRegistry{
		"kafka_admin_topics_read":   kafkaTopicOperationSpecs,
		"kafka_admin_topics_write":  kafkaTopicOperationSpecs,
		"kafka_admin_groups_read":   kafkaGroupOperationSpecs,
		"kafka_admin_groups_write":  kafkaGroupOperationSpecs,
		"kafka_admin_sr_read":       kafkaSchemaRegistryOperationSpecs,
		"kafka_admin_sr_write":      kafkaSchemaRegistryOperationSpecs,
		"kafka_admin_connect_read":  kafkaConnectOperationSpecs,
		"kafka_admin_connect_write": kafkaConnectOperationSpecs,
	}
}

func requireStringEnum(t *testing.T, toolName string, propertySchema any, expected []string) {
	t.Helper()
	actual, ok := stringEnum(propertySchema)
	require.True(t, ok, "%s resource enum must be explicit", toolName)
	require.ElementsMatch(t, expected, actual, toolName)
}

func stringEnum(propertySchema any) ([]string, bool) {
	schema, ok := propertySchema.(map[string]any)
	if !ok {
		return nil, false
	}
	rawEnum, ok := schema["enum"].([]string)
	return rawEnum, ok
}

func TestKafkaSplitToolsExposeModeSpecificParameters(t *testing.T) {
	builderList := []builders.ToolBuilder{
		NewKafkaTopicsToolBuilder(),
		NewKafkaGroupsToolBuilder(),
		NewKafkaSchemaRegistryToolBuilder(),
		NewKafkaConnectToolBuilder(),
	}

	expectedProperties := map[string][]string{
		"kafka_admin_topics_read":   {"resource", "operation", "name", "includeInternal"},
		"kafka_admin_topics_write":  {"resource", "operation", "name", "partitions", "replicationFactor", "configs"},
		"kafka_admin_groups_read":   {"resource", "operation", "group"},
		"kafka_admin_groups_write":  {"resource", "operation", "group", "members", "topic", "partition", "offset"},
		"kafka_admin_sr_read":       {"resource", "operation", "subject", "version"},
		"kafka_admin_sr_write":      {"resource", "operation", "subject", "version", "compatibility", "schemaType", "schema"},
		"kafka_admin_connect_read":  {"resource", "operation", "name"},
		"kafka_admin_connect_write": {"resource", "operation", "name", "config"},
	}
	expectedResourceEnums := map[string][]string{
		"kafka_admin_topics_read":   {"topic", "topics"},
		"kafka_admin_topics_write":  {"topic"},
		"kafka_admin_groups_read":   {"group", "groups"},
		"kafka_admin_groups_write":  {"group"},
		"kafka_admin_sr_read":       {"subjects", "subject", "versions", "version", "compatibility", "types"},
		"kafka_admin_sr_write":      {"subject", "version", "compatibility"},
		"kafka_admin_connect_read":  {"kafka-connect-cluster", "connector", "connectors", "connector-plugins"},
		"kafka_admin_connect_write": {"connector"},
	}

	for _, builder := range builderList {
		tools, err := builder.BuildTools(context.Background(), builders.ToolBuildConfig{
			Features: []string{"all", "all-kafka", "kafka-admin", "kafka-admin-kafka-connect", "kafka-admin-schema-registry"},
		})
		require.NoError(t, err)
		for _, serverTool := range tools {
			tool := serverTool.Tool
			require.ElementsMatch(t, expectedProperties[tool.Name], toolPropertyNames(tool), tool.Name)
			requireStringEnum(t, tool.Name, tool.InputSchema.Properties["resource"], expectedResourceEnums[tool.Name])
		}
	}
}

func TestKafkaReadOnlyBuildsNoWriteTools(t *testing.T) {
	builderList := []builders.ToolBuilder{
		NewKafkaTopicsToolBuilder(),
		NewKafkaGroupsToolBuilder(),
		NewKafkaSchemaRegistryToolBuilder(),
		NewKafkaConnectToolBuilder(),
		NewKafkaPartitionsToolBuilder(),
		NewKafkaProduceToolBuilder(),
		NewKafkaConsumeToolBuilder(),
	}

	for _, builder := range builderList {
		tools, err := builder.BuildTools(context.Background(), builders.ToolBuildConfig{
			ReadOnly: true,
			Features: []string{"all", "all-kafka", "kafka-admin", "kafka-admin-kafka-connect", "kafka-admin-schema-registry", "kafka-client"},
		})
		require.NoError(t, err)
		for _, serverTool := range tools {
			name := serverTool.Tool.Name
			require.NotContains(t, name, "_write")
			require.NotEqual(t, "kafka_admin_partitions_write", name)
			require.NotEqual(t, "kafka_client_produce", name)
		}
	}
}
