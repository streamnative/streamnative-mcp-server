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
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func TestBuildSourceConfigDefaults(t *testing.T) {
	builder := NewPulsarAdminSourcesToolBuilder()

	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{}}}
	config, _, _, err := builder.buildSourceConfig("", "", "name", req)
	require.NoError(t, err)

	builder.applySourceDefaults(config)
	require.Equal(t, defaultTenant, config.Tenant)
	require.Equal(t, defaultNamespace, config.Namespace)
	require.Equal(t, 1, config.Parallelism)
}

func TestBuildSourceConfigMaps(t *testing.T) {
	builder := NewPulsarAdminSourcesToolBuilder()

	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"source-config": map[string]any{
			"path": "/tmp/in",
		},
		"secrets": map[string]any{
			"secret": map[string]any{
				"value": "x",
			},
		},
		"producer-config": map[string]any{
			"maxPendingMessages": 10,
		},
		"batch-source-config": map[string]any{
			"discoveryTriggererClassName": "example",
		},
		"batch-builder": "DEFAULT",
	}}}
	config, _, _, err := builder.buildSourceConfig("tenant", "namespace", "name", req)
	require.NoError(t, err)
	require.Equal(t, "/tmp/in", config.Configs["path"])
	require.NotNil(t, config.Secrets)
	require.NotNil(t, config.ProducerConfig)
	require.Equal(t, 10, config.ProducerConfig.MaxPendingMessages)
	require.NotNil(t, config.BatchSourceConfig)
	require.Equal(t, "example", config.BatchSourceConfig.DiscoveryTriggererClassName)
	require.Equal(t, "DEFAULT", config.BatchBuilder)
}

func TestValidateSourceArchiveArgs(t *testing.T) {
	require.NoError(t, validateSourceArchiveArgs("", ""))
	require.NoError(t, validateSourceArchiveArgs("archive.nar", ""))
	require.NoError(t, validateSourceArchiveArgs("", "file"))
	require.Error(t, validateSourceArchiveArgs("archive.nar", "file"))
}
