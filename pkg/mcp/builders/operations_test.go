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

package builders

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func TestOperationRegistryModeEnumsAndValidation(t *testing.T) {
	registry := OperationRegistry{
		{Name: "list", Mode: OperationModeRead},
		{Name: "get", Mode: OperationModeRead},
		{Name: "create", Mode: OperationModeWrite, Destructive: true},
	}

	require.Equal(t, []string{"list", "get"}, registry.NamesForMode(OperationModeRead))
	require.Equal(t, []string{"create"}, registry.NamesForMode(OperationModeWrite))
	require.NoError(t, registry.ValidateModeOperation(OperationModeRead, "LIST"))
	require.NoError(t, registry.ValidateModeOperation(OperationModeWrite, "create"))
	require.ErrorContains(t, registry.ValidateModeOperation(OperationModeRead, "create"), "not available in read mode")
	require.ErrorContains(t, registry.ValidateModeOperation(OperationModeWrite, "unknown"), "unknown operation")
}

func TestToolAnnotationForMode(t *testing.T) {
	registry := OperationRegistry{
		{Name: "list", Mode: OperationModeRead},
		{Name: "update", Mode: OperationModeWrite, Destructive: false, Idempotent: true},
	}

	read := ToolAnnotationForMode(OperationModeRead, "Read Things", "Manage Things", registry)
	write := ToolAnnotationForMode(OperationModeWrite, "Read Things", "Manage Things", registry)

	readTool := mcp.NewTool("read", read)
	writeTool := mcp.NewTool("write", write)

	require.Equal(t, "Read Things", readTool.Annotations.Title)
	require.NotNil(t, readTool.Annotations.ReadOnlyHint)
	require.NotNil(t, readTool.Annotations.DestructiveHint)
	require.NotNil(t, readTool.Annotations.IdempotentHint)
	require.NotNil(t, readTool.Annotations.OpenWorldHint)
	require.True(t, *readTool.Annotations.ReadOnlyHint)
	require.False(t, *readTool.Annotations.DestructiveHint)
	require.True(t, *readTool.Annotations.IdempotentHint)
	require.True(t, *readTool.Annotations.OpenWorldHint)

	require.Equal(t, "Manage Things", writeTool.Annotations.Title)
	require.NotNil(t, writeTool.Annotations.ReadOnlyHint)
	require.NotNil(t, writeTool.Annotations.DestructiveHint)
	require.NotNil(t, writeTool.Annotations.IdempotentHint)
	require.NotNil(t, writeTool.Annotations.OpenWorldHint)
	require.False(t, *writeTool.Annotations.ReadOnlyHint)
	require.False(t, *writeTool.Annotations.DestructiveHint)
	require.True(t, *writeTool.Annotations.IdempotentHint)
	require.True(t, *writeTool.Annotations.OpenWorldHint)
}
