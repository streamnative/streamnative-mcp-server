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
	"context"
	"testing"
	"time"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func TestBuildTopicToolIncludesPermissionOperations(t *testing.T) {
	builder := NewPulsarAdminTopicToolBuilder()

	tool := builder.buildTopicTool()

	require.Contains(t, tool.InputSchema.Properties, "role")
	require.Contains(t, tool.InputSchema.Properties, "actions")
	require.Contains(t, tool.InputSchema.Properties, "wait")

	operationSchema, ok := tool.InputSchema.Properties["operation"].(map[string]any)
	require.True(t, ok)

	description, ok := operationSchema["description"].(string)
	require.True(t, ok)
	require.Contains(t, description, "get-permissions")
	require.Contains(t, description, "grant-permissions")
	require.Contains(t, description, "revoke-permissions")
	require.Contains(t, description, "compact-status")
}

func TestParseTopicActions(t *testing.T) {
	actions, err := parseTopicActions([]string{"produce", "consume"})
	require.NoError(t, err)
	require.Equal(t, []utils.AuthAction{"produce", "consume"}, actions)

	_, err = parseTopicActions([]string{"invalid"})
	require.Error(t, err)
}

func TestTopicPermissionWriteOperationsRespectReadOnly(t *testing.T) {
	require.True(t, isReadOnlyRestrictedTopicOperation("grant-permissions"))
	require.True(t, isReadOnlyRestrictedTopicOperation("revoke-permissions"))
	require.False(t, isReadOnlyRestrictedTopicOperation("get-permissions"))
}

func TestNormalizeTopicOperationSupportsLegacyAliases(t *testing.T) {
	require.Equal(t, "compact-status", normalizeTopicOperation("status"))
	require.Equal(t, "compact-status", normalizeTopicOperation("COMPACT-STATUS"))
	require.Equal(t, "get", normalizeTopicOperation(" get "))
}

func TestTopicGrantPermissionsBlockedInReadOnlyMode(t *testing.T) {
	builder := NewPulsarAdminTopicToolBuilder()
	handler := builder.buildTopicHandler(true)

	result, err := handler(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"resource":  "topic",
				"operation": "grant-permissions",
				"topic":     "persistent://public/default/example",
				"role":      "admin",
				"actions":   []string{"produce"},
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.True(t, result.IsError)

	text, ok := result.Content[0].(mcp.TextContent)
	require.True(t, ok)
	require.Contains(t, text.Text, "read-only mode")
}

func TestWaitForTopicLongRunningStatusStopsOnContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := waitForTopicLongRunningStatus(ctx, true, time.Millisecond, func() bool {
		return true
	}, func() error {
		return nil
	})

	require.ErrorIs(t, err, context.Canceled)
}

func TestWaitForTopicLongRunningStatusPollsUntilComplete(t *testing.T) {
	ctx := context.Background()
	attempts := 0
	running := true

	err := waitForTopicLongRunningStatus(ctx, true, time.Millisecond, func() bool {
		return running
	}, func() error {
		attempts++
		if attempts == 2 {
			running = false
		}
		return nil
	})

	require.NoError(t, err)
	require.Equal(t, 2, attempts)
}

func TestWaitForTopicLongRunningStatusReturnsRefreshError(t *testing.T) {
	ctx := context.Background()
	expectedErr := context.DeadlineExceeded

	err := waitForTopicLongRunningStatus(ctx, true, time.Millisecond, func() bool {
		return true
	}, func() error {
		return expectedErr
	})

	require.ErrorIs(t, err, expectedErr)
}
