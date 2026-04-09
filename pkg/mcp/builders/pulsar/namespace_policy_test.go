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

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func TestBuildNamespaceSetPolicyToolIncludesBacklogPolicyArgument(t *testing.T) {
	builder := NewPulsarAdminNamespacePolicyToolBuilder()

	tool := builder.buildNamespaceSetPolicyTool()

	require.Contains(t, tool.InputSchema.Properties, "policy")
	require.Contains(t, tool.InputSchema.Properties, "backlog-policy")
	require.Contains(t, tool.InputSchema.Properties, "enabled")
	require.Contains(t, tool.InputSchema.Required, "policy")
}

func TestBuildTopicAutoCreationConfigEnabled(t *testing.T) {
	builder := NewPulsarAdminNamespacePolicyToolBuilder()
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"enabled":    true,
		"topic-type": "partitioned",
		"partitions": 3.0,
	}}}

	config, err := builder.buildTopicAutoCreationConfig(req)

	require.NoError(t, err)
	require.True(t, config.Allow)
	require.Equal(t, utils.Partitioned, config.Type)
	require.NotNil(t, config.Partitions)
	require.Equal(t, 3, *config.Partitions)
}

func TestBuildTopicAutoCreationConfigDisabled(t *testing.T) {
	builder := NewPulsarAdminNamespacePolicyToolBuilder()
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"enabled": false,
	}}}

	config, err := builder.buildTopicAutoCreationConfig(req)

	require.NoError(t, err)
	require.False(t, config.Allow)
	require.Nil(t, config.Partitions)
}

func TestBuildTopicAutoCreationConfigRequiresTypeWhenEnabled(t *testing.T) {
	builder := NewPulsarAdminNamespacePolicyToolBuilder()
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"enabled": true,
	}}}

	_, err := builder.buildTopicAutoCreationConfig(req)

	require.Error(t, err)
}

func TestBuildPersistencePolicy(t *testing.T) {
	builder := NewPulsarAdminNamespacePolicyToolBuilder()
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"ensemble-size":           3.0,
		"write-quorum-size":       2.0,
		"ack-quorum-size":         2.0,
		"ml-mark-delete-max-rate": 1.5,
	}}}

	policy, err := builder.buildPersistencePolicy(req)

	require.NoError(t, err)
	require.Equal(t, 3, policy.BookkeeperEnsemble)
	require.Equal(t, 2, policy.BookkeeperWriteQuorum)
	require.Equal(t, 2, policy.BookkeeperAckQuorum)
	require.Equal(t, 1.5, policy.ManagedLedgerMaxMarkDeleteRate)
}

func TestBuildDispatchRateDefaults(t *testing.T) {
	builder := NewPulsarAdminNamespacePolicyToolBuilder()
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{}}}

	rate, err := builder.buildDispatchRate(req)

	require.NoError(t, err)
	require.Equal(t, -1, rate.DispatchThrottlingRateInMsg)
	require.EqualValues(t, -1, rate.DispatchThrottlingRateInByte)
	require.Equal(t, 1, rate.RatePeriodInSecond)
}

func TestBuildPublishRateDefaults(t *testing.T) {
	builder := NewPulsarAdminNamespacePolicyToolBuilder()
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{}}}

	rate, err := builder.buildPublishRate(req)

	require.NoError(t, err)
	require.Equal(t, -1, rate.PublishThrottlingRateInMsg)
	require.EqualValues(t, -1, rate.PublishThrottlingRateInByte)
}

func TestBuildSubscribeRateDefaults(t *testing.T) {
	builder := NewPulsarAdminNamespacePolicyToolBuilder()
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{}}}

	rate, err := builder.buildSubscribeRate(req)

	require.NoError(t, err)
	require.Equal(t, -1, rate.SubscribeThrottlingRatePerConsumer)
	require.Equal(t, 30, rate.RatePeriodInSecond)
}

func TestRequireBooleanArgument(t *testing.T) {
	builder := NewPulsarAdminNamespacePolicyToolBuilder()
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"enabled": true,
	}}}

	value, err := builder.requireBooleanArgument(req, "enabled")

	require.NoError(t, err)
	require.True(t, value)
}

func TestRequireBooleanArgumentRequiresPresence(t *testing.T) {
	builder := NewPulsarAdminNamespacePolicyToolBuilder()
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{}}}

	_, err := builder.requireBooleanArgument(req, "enabled")

	require.Error(t, err)
}

func TestRequireIntegerArgumentRejectsFractions(t *testing.T) {
	builder := NewPulsarAdminNamespacePolicyToolBuilder()
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"count": 1.5,
	}}}

	_, err := builder.requireIntegerArgument(req, "count")

	require.Error(t, err)
}
