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

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	pulsarsession "github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPulsarAdminStatusToolBuilder(t *testing.T) {
	builder := NewPulsarAdminStatusToolBuilder()

	t.Run("BuildTools_Success", func(t *testing.T) {
		tools, err := builder.BuildTools(context.Background(), builders.ToolBuildConfig{
			Features: []string{"pulsar-admin-brokers-status"},
		})
		require.NoError(t, err)
		require.Len(t, tools, 1)
		assert.Equal(t, "pulsar_admin_status", tools[0].Tool.Name)
	})

	t.Run("BuildTools_NoFeatures", func(t *testing.T) {
		tools, err := builder.BuildTools(context.Background(), builders.ToolBuildConfig{
			Features: []string{"unrelated-feature"},
		})
		require.NoError(t, err)
		assert.Len(t, tools, 0)
	})

	t.Run("Handler_MissingSession", func(t *testing.T) {
		tools, err := builder.BuildTools(context.Background(), builders.ToolBuildConfig{
			Features: []string{"pulsar-admin-brokers-status"},
		})
		require.NoError(t, err)
		require.Len(t, tools, 1)

		result, callErr := tools[0].Handler(context.Background(), mcp.CallToolRequest{})
		require.NoError(t, callErr)
		require.NotNil(t, result)
		assert.True(t, result.IsError)
	})

	t.Run("Handler_EmptyWebServiceURL", func(t *testing.T) {
		tools, err := builder.BuildTools(context.Background(), builders.ToolBuildConfig{
			Features: []string{"pulsar-admin-brokers-status"},
		})
		require.NoError(t, err)
		require.Len(t, tools, 1)

		ctx := mcpCtx.WithPulsarSession(context.Background(), &pulsarsession.Session{
			PulsarCtlConfig: &cmdutils.ClusterConfig{},
		})

		result, callErr := tools[0].Handler(ctx, mcp.CallToolRequest{})
		require.NoError(t, callErr)
		require.NotNil(t, result)
		assert.True(t, result.IsError)
		require.Len(t, result.Content, 1)

		text, ok := result.Content[0].(mcp.TextContent)
		require.True(t, ok)
		assert.Contains(t, text.Text, "Please set the cluster context first")
	})
}
