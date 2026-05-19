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
