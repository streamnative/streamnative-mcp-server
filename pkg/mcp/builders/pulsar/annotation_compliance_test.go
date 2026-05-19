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

package pulsar

import (
	"context"
	"strings"
	"testing"

	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	"github.com/stretchr/testify/require"
)

func TestPulsarToolAnnotationCompliance(t *testing.T) {
	builderList := allPulsarComplianceBuilders()

	for _, builder := range builderList {
		tools, err := builder.BuildTools(context.Background(), builders.ToolBuildConfig{
			Features: []string{"all", "all-pulsar", "pulsar-admin", "pulsar-client"},
		})
		require.NoError(t, err)
		for _, serverTool := range tools {
			tool := serverTool.Tool
			require.NotEmpty(t, tool.Annotations.Title, tool.Name)
			require.NotNil(t, tool.Annotations.ReadOnlyHint, tool.Name)
			require.NotNil(t, tool.Annotations.DestructiveHint, tool.Name)
			require.LessOrEqual(t, len(tool.Name), 64, tool.Name)

			isRead := strings.HasSuffix(tool.Name, "_read") || strings.HasPrefix(tool.Name, "pulsar_admin_namespace_policy_get") || tool.Name == "pulsar_admin_status" || tool.Name == "pulsar_admin_broker_stats" || tool.Name == "pulsar_admin_functions_worker"
			isWrite := strings.HasSuffix(tool.Name, "_write") || strings.Contains(tool.Name, "_set") || strings.Contains(tool.Name, "_remove") || strings.Contains(tool.Name, "produce") || strings.Contains(tool.Name, "consume")
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
	writeOperations := pulsarComplianceWriteOperations()
	seenRead, seenWrite := false, false
	for _, op := range rawEnum {
		if _, isWrite := writeOperations[op]; isWrite {
			seenWrite = true
		} else {
			seenRead = true
		}
	}
	require.False(t, seenRead && seenWrite, toolName)
}

func pulsarComplianceWriteOperations() map[string]struct{} {
	writeOperations := map[string]struct{}{
		"create": {},
		"update": {},
		"delete": {},
		"upload": {},
	}
	for _, source := range []map[string]struct{}{
		readOnlyRestrictedTopicOperations,
		readOnlyRestrictedSubscriptionOperations,
		pulsarNamespaceWriteOperations,
		readOnlyRestrictedTopicPolicyOperations,
		pulsarBrokerWriteOperations,
		pulsarClusterWriteOperations,
		readOnlyRestrictedFunctionOperations,
		pulsarSinkWriteOperations,
		pulsarSourceWriteOperations,
		pulsarPackageWriteOperations,
		pulsarSchemaWriteOperations,
		pulsarTenantWriteOperations,
		pulsarNsIsolationPolicyWriteOperations,
		pulsarResourceQuotaWriteOperations,
	} {
		for op := range source {
			writeOperations[op] = struct{}{}
		}
	}
	return writeOperations
}

func TestPulsarReadOnlyBuildsNoWriteTools(t *testing.T) {
	for _, builder := range allPulsarComplianceBuilders() {
		tools, err := builder.BuildTools(context.Background(), builders.ToolBuildConfig{
			ReadOnly: true,
			Features: []string{"all", "all-pulsar", "pulsar-admin", "pulsar-client"},
		})
		require.NoError(t, err)
		for _, serverTool := range tools {
			name := serverTool.Tool.Name
			require.NotContains(t, name, "_write")
			require.NotContains(t, name, "_set")
			require.NotContains(t, name, "_remove")
			require.NotEqual(t, "pulsar_client_produce", name)
			require.NotEqual(t, "pulsar_client_consume", name)
		}
	}
}

func allPulsarComplianceBuilders() []builders.ToolBuilder {
	return []builders.ToolBuilder{
		NewPulsarAdminBrokerStatsToolBuilder(),
		NewPulsarAdminBrokersToolBuilder(),
		NewPulsarAdminClusterToolBuilder(),
		NewPulsarAdminFunctionsToolBuilder(),
		NewPulsarAdminFunctionsWorkerToolBuilder(),
		NewPulsarAdminNamespaceToolBuilder(),
		NewPulsarAdminNamespacePolicyToolBuilder(),
		NewPulsarAdminNsIsolationPolicyToolBuilder(),
		NewPulsarAdminPackagesToolBuilder(),
		NewPulsarAdminResourceQuotasToolBuilder(),
		NewPulsarAdminSchemaToolBuilder(),
		NewPulsarAdminSinksToolBuilder(),
		NewPulsarAdminSourcesToolBuilder(),
		NewPulsarAdminStatusToolBuilder(),
		NewPulsarAdminSubscriptionToolBuilder(),
		NewPulsarAdminTenantToolBuilder(),
		NewPulsarAdminTopicPolicyToolBuilder(),
		NewPulsarAdminTopicToolBuilder(),
		NewPulsarClientConsumeToolBuilder(),
		NewPulsarClientProduceToolBuilder(),
	}
}
