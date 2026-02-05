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
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func TestBuildSinkConfigDefaults(t *testing.T) {
	builder := NewPulsarAdminSinksToolBuilder()

	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{}}}
	config, _, _, err := builder.buildSinkConfig("", "", "name", req)
	require.NoError(t, err)

	builder.applySinkDefaults(config)
	require.Equal(t, defaultTenant, config.Tenant)
	require.Equal(t, defaultNamespace, config.Namespace)
	require.Equal(t, 1, config.Parallelism)
	require.True(t, config.CleanupSubscription)
}

func TestBuildSinkConfigInputsParsing(t *testing.T) {
	builder := NewPulsarAdminSinksToolBuilder()

	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"inputs": "t1,t2",
	}}}
	config, _, _, err := builder.buildSinkConfig("tenant", "namespace", "name", req)
	require.NoError(t, err)
	require.Equal(t, []string{"t1", "t2"}, config.Inputs)
}

func TestBuildSinkConfigMaps(t *testing.T) {
	builder := NewPulsarAdminSinksToolBuilder()

	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"custom-serde-inputs": map[string]any{
			"t1": "serde",
		},
		"custom-schema-inputs": map[string]any{
			"t1": "schema",
		},
		"input-specs": map[string]any{
			"t1": map[string]any{
				"schemaType": "string",
			},
		},
		"secrets": map[string]any{
			"secret": map[string]any{
				"value": "x",
			},
		},
		"sink-config": map[string]any{
			"path": "/tmp/out",
		},
	}}}
	config, _, _, err := builder.buildSinkConfig("tenant", "namespace", "name", req)
	require.NoError(t, err)
	require.Equal(t, "serde", config.TopicToSerdeClassName["t1"])
	require.Equal(t, "schema", config.TopicToSchemaType["t1"])
	require.Equal(t, "string", config.InputSpecs["t1"].SchemaType)
	require.Equal(t, "/tmp/out", config.Configs["path"])
	require.NotNil(t, config.Secrets)
}

func TestValidateSinkArchiveArgs(t *testing.T) {
	require.NoError(t, validateSinkArchiveArgs("", ""))
	require.NoError(t, validateSinkArchiveArgs("archive.nar", ""))
	require.NoError(t, validateSinkArchiveArgs("", "file"))
	require.Error(t, validateSinkArchiveArgs("archive.nar", "file"))
}
