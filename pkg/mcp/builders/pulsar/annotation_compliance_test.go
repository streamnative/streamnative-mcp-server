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

	"github.com/mark3labs/mcp-go/mcp"
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
	rawEnum, ok := pulsarStringEnum(operationSchema)
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

func requirePulsarStringEnum(t *testing.T, toolName string, propertySchema any, expected []string) {
	t.Helper()
	actual, ok := pulsarStringEnum(propertySchema)
	require.True(t, ok, "%s resource enum must be explicit", toolName)
	require.ElementsMatch(t, expected, actual, toolName)
}

func pulsarStringEnum(propertySchema any) ([]string, bool) {
	schema, ok := propertySchema.(map[string]any)
	if !ok {
		return nil, false
	}
	rawEnum, ok := schema["enum"].([]string)
	return rawEnum, ok
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

func TestPulsarSplitToolsExposeModeSpecificParameters(t *testing.T) {
	expectedProperties := map[string][]string{
		"pulsar_admin_brokers_read":            {"resource", "operation", "clusterName", "brokerUrl", "configType"},
		"pulsar_admin_brokers_write":           {"resource", "operation", "configName", "configValue"},
		"pulsar_admin_cluster_read":            {"resource", "operation", "cluster_name", "domain_name"},
		"pulsar_admin_cluster_write":           {"resource", "operation", "cluster_name", "domain_name", "service_url", "service_url_tls", "broker_service_url", "broker_service_url_tls", "peer_cluster_names", "brokers"},
		"pulsar_admin_functions_read":          {"operation", "fqfn", "tenant", "namespace", "name", "instanceId", "key", "path", "destinationFile"},
		"pulsar_admin_namespace_read":          {"operation", "tenant", "namespace"},
		"pulsar_admin_nsisolationpolicy_read":  {"resource", "operation", "cluster", "name"},
		"pulsar_admin_nsisolationpolicy_write": {"resource", "operation", "cluster", "name", "namespaces", "primary", "secondary", "autoFailoverPolicyType", "autoFailoverPolicyParams"},
		"pulsar_admin_package_read":            {"resource", "operation", "packageName", "namespace", "type", "path"},
		"pulsar_admin_package_write":           {"resource", "operation", "packageName", "description", "contact", "path", "properties"},
		"pulsar_admin_resourcequota_read":      {"resource", "operation", "namespace", "bundle"},
		"pulsar_admin_schema_read":             {"resource", "operation", "topic", "version"},
		"pulsar_admin_sinks_read":              {"operation", "tenant", "namespace", "name"},
		"pulsar_admin_sources_read":            {"operation", "tenant", "namespace", "name"},
		"pulsar_admin_subscription_read":       {"resource", "operation", "topic", "subscription", "ledgerId", "entryId", "count"},
		"pulsar_admin_tenant_read":             {"resource", "operation", "tenant"},
		"pulsar_admin_topic_policy_read":       {"operation", "topic", "applied", "type"},
		"pulsar_admin_topic_read":              {"resource", "operation", "topic", "namespace", "partitioned", "per-partition", "wait"},
		"pulsar_admin_topic_write":             {"resource", "operation", "topic", "partitions", "force", "non-partitioned", "config", "messageId", "role", "actions"},
	}
	expectedResourceEnums := map[string][]string{
		"pulsar_admin_brokers_read":            {"brokers", "health", "config", "namespaces"},
		"pulsar_admin_brokers_write":           {"config"},
		"pulsar_admin_cluster_read":            {"cluster", "peer_clusters", "failure_domain"},
		"pulsar_admin_cluster_write":           {"cluster", "peer_clusters", "failure_domain"},
		"pulsar_admin_nsisolationpolicy_read":  {"policy", "broker", "brokers"},
		"pulsar_admin_nsisolationpolicy_write": {"policy"},
		"pulsar_admin_package_read":            {"package", "packages"},
		"pulsar_admin_package_write":           {"package"},
		"pulsar_admin_topic_read":              {"topic", "topics"},
		"pulsar_admin_topic_write":             {"topic"},
	}

	for _, builder := range allPulsarComplianceBuilders() {
		tools, err := builder.BuildTools(context.Background(), builders.ToolBuildConfig{
			Features: []string{"all", "all-pulsar", "pulsar-admin", "pulsar-client"},
		})
		require.NoError(t, err)
		for _, serverTool := range tools {
			tool := serverTool.Tool
			expected, ok := expectedProperties[tool.Name]
			if !ok {
				continue
			}
			require.ElementsMatch(t, expected, pulsarToolPropertyNames(tool), tool.Name)
			if expectedEnum, ok := expectedResourceEnums[tool.Name]; ok {
				requirePulsarStringEnum(t, tool.Name, tool.InputSchema.Properties["resource"], expectedEnum)
			}
		}
	}
}

func pulsarToolPropertyNames(tool mcp.Tool) []string {
	names := make([]string, 0, len(tool.InputSchema.Properties))
	for name := range tool.InputSchema.Properties {
		names = append(names, name)
	}
	return names
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
