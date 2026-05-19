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
			require.LessOrEqual(t, len(tool.Name), 64, tool.Name)

			isRead := strings.HasSuffix(tool.Name, "_read")
			isWrite := strings.HasSuffix(tool.Name, "_write") || strings.Contains(tool.Name, "produce") || strings.Contains(tool.Name, "consume")
			if isRead {
				require.True(t, *tool.Annotations.ReadOnlyHint, tool.Name)
				require.False(t, *tool.Annotations.DestructiveHint, tool.Name)
			}
			if isWrite {
				require.False(t, *tool.Annotations.ReadOnlyHint, tool.Name)
				require.True(t, *tool.Annotations.DestructiveHint, tool.Name)
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
	schema, ok := operationSchema.(map[string]any)
	if !ok {
		return
	}
	rawEnum, ok := schema["enum"].([]string)
	if !ok {
		return
	}
	writeOperations := map[string]struct{}{
		"create": {}, "delete": {}, "set": {}, "update": {}, "restart": {}, "pause": {}, "resume": {},
		"remove-members": {}, "delete-offset": {}, "set-offset": {},
	}
	seenRead, seenWrite := false, false
	for _, op := range rawEnum {
		if _, ok := writeOperations[op]; ok {
			seenWrite = true
		} else {
			seenRead = true
		}
	}
	require.False(t, seenRead && seenWrite, toolName)
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

	for _, builder := range builderList {
		tools, err := builder.BuildTools(context.Background(), builders.ToolBuildConfig{
			Features: []string{"all", "all-kafka", "kafka-admin", "kafka-admin-kafka-connect", "kafka-admin-schema-registry"},
		})
		require.NoError(t, err)
		for _, serverTool := range tools {
			tool := serverTool.Tool
			require.ElementsMatch(t, expectedProperties[tool.Name], toolPropertyNames(tool), tool.Name)
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
