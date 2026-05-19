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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	"github.com/stretchr/testify/require"
)

type kafkaOperationDocSpec struct {
	File      string
	ReadTool  string
	WriteTool string
	Registry  builders.OperationRegistry
}

func TestKafkaOperationDocsMatchSpecs(t *testing.T) {
	for _, spec := range kafkaOperationDocSpecs() {
		docPath := filepath.Join("..", "..", "..", "..", "docs", "tools", spec.File)
		// #nosec G304 -- docPath is built from fixed test fixtures in kafkaOperationDocSpecs.
		content, err := os.ReadFile(docPath)
		require.NoError(t, err, spec.File)
		require.Equal(t, expectedKafkaOperationBlock(spec), extractGeneratedOperationBlock(t, string(content)), spec.File)
	}
}

func kafkaOperationDocSpecs() []kafkaOperationDocSpec {
	return []kafkaOperationDocSpec{
		{File: "kafka_admin_topics.md", ReadTool: "kafka_admin_topics_read", WriteTool: "kafka_admin_topics_write", Registry: kafkaTopicOperationSpecs},
		{File: "kafka_admin_groups.md", ReadTool: "kafka_admin_groups_read", WriteTool: "kafka_admin_groups_write", Registry: kafkaGroupOperationSpecs},
		{File: "kafka_admin_schema_registry.md", ReadTool: "kafka_admin_sr_read", WriteTool: "kafka_admin_sr_write", Registry: kafkaSchemaRegistryOperationSpecs},
		{File: "kafka_admin_connect.md", ReadTool: "kafka_admin_connect_read", WriteTool: "kafka_admin_connect_write", Registry: kafkaConnectOperationSpecs},
	}
}

func expectedKafkaOperationBlock(spec kafkaOperationDocSpec) string {
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
	return strings.TrimSpace(content[start:end])
}

func formatOperationNames(operations []string) string {
	quoted := make([]string, 0, len(operations))
	for _, operation := range operations {
		quoted = append(quoted, "`"+operation+"`")
	}
	return strings.Join(quoted, ", ")
}
