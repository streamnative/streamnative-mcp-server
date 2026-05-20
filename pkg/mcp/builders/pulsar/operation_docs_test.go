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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	"github.com/stretchr/testify/require"
)

type pulsarOperationDocSpec struct {
	File      string
	ReadTool  string
	WriteTool string
	Registry  builders.OperationRegistry
}

func TestPulsarOperationDocsMatchSpecs(t *testing.T) {
	for _, spec := range pulsarOperationDocSpecs() {
		docPath := filepath.Join("..", "..", "..", "..", "docs", "tools", spec.File)
		// #nosec G304 -- docPath is built from fixed test fixtures in pulsarOperationDocSpecs.
		content, err := os.ReadFile(docPath)
		require.NoError(t, err, spec.File)
		require.Equal(t, expectedPulsarOperationBlock(spec), extractGeneratedOperationBlock(t, string(content)), spec.File)
	}
}

func pulsarOperationDocSpecs() []pulsarOperationDocSpec {
	return []pulsarOperationDocSpec{
		{File: "pulsar_admin_brokers.md", ReadTool: "pulsar_admin_brokers_read", WriteTool: "pulsar_admin_brokers_write", Registry: pulsarBrokerOperationSpecs},
		{File: "pulsar_admin_clusters.md", ReadTool: "pulsar_admin_cluster_read", WriteTool: "pulsar_admin_cluster_write", Registry: pulsarClusterOperationSpecs},
		{File: "pulsar_admin_functions.md", ReadTool: "pulsar_admin_functions_read", WriteTool: "pulsar_admin_functions_write", Registry: pulsarFunctionOperationSpecs},
		{File: "pulsar_admin_namespaces.md", ReadTool: "pulsar_admin_namespace_read", WriteTool: "pulsar_admin_namespace_write", Registry: pulsarNamespaceOperationSpecs},
		{File: "pulsar_admin_nsisolationpolicy.md", ReadTool: "pulsar_admin_nsisolationpolicy_read", WriteTool: "pulsar_admin_nsisolationpolicy_write", Registry: pulsarNsIsolationPolicyOperationSpecs},
		{File: "pulsar_admin_packages.md", ReadTool: "pulsar_admin_package_read", WriteTool: "pulsar_admin_package_write", Registry: pulsarPackageOperationSpecs},
		{File: "pulsar_admin_resource_quotas.md", ReadTool: "pulsar_admin_resourcequota_read", WriteTool: "pulsar_admin_resourcequota_write", Registry: pulsarResourceQuotaOperationSpecs},
		{File: "pulsar_admin_schemas.md", ReadTool: "pulsar_admin_schema_read", WriteTool: "pulsar_admin_schema_write", Registry: pulsarSchemaOperationSpecs},
		{File: "pulsar_admin_sinks.md", ReadTool: "pulsar_admin_sinks_read", WriteTool: "pulsar_admin_sinks_write", Registry: pulsarSinkOperationSpecs},
		{File: "pulsar_admin_sources.md", ReadTool: "pulsar_admin_sources_read", WriteTool: "pulsar_admin_sources_write", Registry: pulsarSourceOperationSpecs},
		{File: "pulsar_admin_subscriptions.md", ReadTool: "pulsar_admin_subscription_read", WriteTool: "pulsar_admin_subscription_write", Registry: pulsarSubscriptionOperationSpecs},
		{File: "pulsar_admin_tenants.md", ReadTool: "pulsar_admin_tenant_read", WriteTool: "pulsar_admin_tenant_write", Registry: pulsarTenantOperationSpecs},
		{File: "pulsar_admin_topic_policy.md", ReadTool: "pulsar_admin_topic_policy_read", WriteTool: "pulsar_admin_topic_policy_write", Registry: pulsarTopicPolicyOperationSpecs},
		{File: "pulsar_admin_topics.md", ReadTool: "pulsar_admin_topic_read", WriteTool: "pulsar_admin_topic_write", Registry: pulsarTopicOperationSpecs},
	}
}

func expectedPulsarOperationBlock(spec pulsarOperationDocSpec) string {
	return strings.Join([]string{
		"<!-- generated:operations:start -->",
		"| Tool | Mode | Operations |",
		"|---|---|---|",
		fmt.Sprintf("| `%s` | read | %s |", spec.ReadTool, formatOperationNames(spec.Registry.NamesForMode(builders.OperationModeRead))),
		fmt.Sprintf("| `%s` | write | %s |", spec.WriteTool, formatOperationNames(spec.Registry.NamesForMode(builders.OperationModeWrite))),
		"<!-- generated:operations:end -->",
	}, "\n")
}

func extractGeneratedOperationBlock(t *testing.T, content string) string {
	t.Helper()
	startMarker := "<!-- generated:operations:start -->"
	endMarker := "<!-- generated:operations:end -->"
	start := strings.Index(content, startMarker)
	require.NotEqual(t, -1, start, "missing generated operation start marker")
	end := strings.Index(content[start:], endMarker)
	require.NotEqual(t, -1, end, "missing generated operation end marker")
	end += start + len(endMarker)
	block := strings.TrimSpace(content[start:end])
	return strings.ReplaceAll(block, "\r\n", "\n")
}

func formatOperationNames(operations []string) string {
	quoted := make([]string, 0, len(operations))
	for _, operation := range operations {
		quoted = append(quoted, "`"+operation+"`")
	}
	return strings.Join(quoted, ", ")
}
