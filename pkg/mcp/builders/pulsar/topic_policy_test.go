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
	"testing"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func TestBuildTopicPolicyToolIncludesPulsarctlParityOperations(t *testing.T) {
	builder := NewPulsarAdminTopicPolicyToolBuilder()

	tool := builder.buildTopicPolicyTool(toolModeRead)
	operationSchema, ok := tool.InputSchema.Properties["operation"].(map[string]any)
	require.True(t, ok)

	description, ok := operationSchema["description"].(string)
	require.True(t, ok)
	require.Contains(t, description, "get-message-ttl")
	require.Contains(t, description, "get-max-producers")
	require.Contains(t, description, "get-backlog-quotas")
	require.Contains(t, description, "get-subscription-dispatch-rate")
	require.Contains(t, description, "get-inactive-topic-policies")
}

func TestBuildTopicPolicyToolIncludesModeSpecificParameters(t *testing.T) {
	builder := NewPulsarAdminTopicPolicyToolBuilder()

	readTool := builder.buildTopicPolicyTool(toolModeRead)
	require.Contains(t, readTool.InputSchema.Properties, "applied")
	require.NotContains(t, readTool.InputSchema.Properties, "count")
	require.NotContains(t, readTool.InputSchema.Properties, "limit-size")
	require.NotContains(t, readTool.InputSchema.Properties, "delete-mode")
	require.NotContains(t, readTool.InputSchema.Properties, "subscription-types")

	writeTool := builder.buildTopicPolicyTool(toolModeWrite)
	require.NotContains(t, writeTool.InputSchema.Properties, "applied")
	require.Contains(t, writeTool.InputSchema.Properties, "count")
	require.Contains(t, writeTool.InputSchema.Properties, "limit-size")
	require.Contains(t, writeTool.InputSchema.Properties, "delete-mode")
	require.Contains(t, writeTool.InputSchema.Properties, "subscription-types")
}

func TestNormalizeTopicPolicyOperationSupportsLegacyAliases(t *testing.T) {
	testCases := map[string]string{
		"get_ttl":                  "get-message-ttl",
		"set_compaction":           "set-compaction-threshold",
		"get-deduplication-status": "get-deduplication",
		"get-backlog-quota":        "get-backlog-quotas",
		"GET_SUBSCRIPTION_TYPES":   "get-subscription-types",
	}

	for input, expected := range testCases {
		require.Equal(t, expected, normalizeTopicPolicyOperation(input))
	}
}

func TestTopicPolicyWriteOperationsRespectReadOnly(t *testing.T) {
	require.True(t, isReadOnlyRestrictedTopicPolicyOperation("set-message-ttl"))
	require.True(t, isReadOnlyRestrictedTopicPolicyOperation("set-backlog-quota"))
	require.True(t, isReadOnlyRestrictedTopicPolicyOperation("set_deduplication"))
	require.True(t, isReadOnlyRestrictedTopicPolicyOperation("remove-compaction-threshold"))
	require.False(t, isReadOnlyRestrictedTopicPolicyOperation("get-message-ttl"))
	require.False(t, isReadOnlyRestrictedTopicPolicyOperation("get-subscription-types"))
}

func TestTopicPolicyHandlerBlocksWriteBeforeSessionLookup(t *testing.T) {
	builder := NewPulsarAdminTopicPolicyToolBuilder()
	handler := builder.buildTopicPolicyHandler(toolModeRead)

	result, err := handler(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"operation":   "set-message-ttl",
				"topic":       "persistent://public/default/example",
				"ttl-seconds": 10,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.True(t, result.IsError)

	text, ok := result.Content[0].(mcp.TextContent)
	require.True(t, ok)
	require.Contains(t, text.Text, "not available in read mode")
}

func TestBuildDelayedDeliveryDataUsesPulsarctlStyleArguments(t *testing.T) {
	builder := NewPulsarAdminTopicPolicyToolBuilder()

	data, err := builder.buildDelayedDeliveryData(mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"enable": true,
				"time":   "10s",
			},
		},
	})
	require.NoError(t, err)
	require.True(t, data.Active)
	require.Equal(t, 10.0, data.TickTime)
}

func TestBuildInactiveTopicPoliciesUsesPulsarctlStyleArguments(t *testing.T) {
	builder := NewPulsarAdminTopicPolicyToolBuilder()

	policies, err := builder.buildInactiveTopicPolicies(mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"delete-while-inactive": true,
				"max-inactive-duration": "1h",
				"delete-mode":           "delete_when_no_subscriptions",
			},
		},
	})
	require.NoError(t, err)
	require.True(t, policies.DeleteWhileInactive)
	require.Equal(t, 3600, policies.MaxInactiveDurationSeconds)
	require.NotNil(t, policies.InactiveTopicDeleteMode)
	require.Equal(t, utils.DeleteWhenNoSubscriptions, *policies.InactiveTopicDeleteMode)
}
