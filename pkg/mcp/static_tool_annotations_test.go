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

package mcp

import (
	"testing"

	mcpgotypes "github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/require"
)

func TestStreamNativeStaticToolAnnotations(t *testing.T) {
	tools := []struct {
		name        string
		tool        mcpgotypes.Tool
		readOnly    bool
		destructive bool
		idempotent  bool
		openWorld   bool
	}{
		{name: "sncloud_logs", tool: NewSNCloudLogsTool(), readOnly: true, idempotent: true, openWorld: true},
		{name: "sncloud_resources_apply", tool: NewSNCloudResourcesApplyTool(), destructive: true, openWorld: true},
		{name: "sncloud_resources_delete", tool: NewSNCloudResourcesDeleteTool(), destructive: true, openWorld: true},
	}

	for _, tt := range tools {
		requireToolAnnotation(t, tt.name, tt.tool, tt.readOnly, tt.destructive, tt.idempotent, tt.openWorld)
	}
}

func TestStreamNativeContextToolAnnotations(t *testing.T) {
	server := mcpserver.NewMCPServer("test", "test")
	RegisterContextTools(server, []string{string(FeatureStreamNativeCloud)}, false, false)

	expectations := map[string]struct {
		readOnly    bool
		destructive bool
		idempotent  bool
		openWorld   bool
	}{
		"sncloud_context_whoami":             {readOnly: true, idempotent: true, openWorld: true},
		"sncloud_context_available_clusters": {readOnly: true, idempotent: true, openWorld: true},
		"sncloud_context_use_cluster":        {},
		"sncloud_context_reset":              {},
	}

	for name, expected := range expectations {
		serverTool := server.GetTool(name)
		require.NotNil(t, serverTool, name)
		tool := serverTool.Tool
		requireToolAnnotation(t, name, tool, expected.readOnly, expected.destructive, expected.idempotent, expected.openWorld)
	}
}

func requireToolAnnotation(t *testing.T, name string, tool mcpgotypes.Tool, readOnly, destructive, idempotent, openWorld bool) {
	t.Helper()
	require.NotEmpty(t, tool.Annotations.Title, name)
	require.NotNil(t, tool.Annotations.ReadOnlyHint, name)
	require.NotNil(t, tool.Annotations.DestructiveHint, name)
	require.NotNil(t, tool.Annotations.IdempotentHint, name)
	require.NotNil(t, tool.Annotations.OpenWorldHint, name)
	require.Equal(t, readOnly, *tool.Annotations.ReadOnlyHint, name)
	require.Equal(t, destructive, *tool.Annotations.DestructiveHint, name)
	require.Equal(t, idempotent, *tool.Annotations.IdempotentHint, name)
	require.Equal(t, openWorld, *tool.Annotations.OpenWorldHint, name)
}
