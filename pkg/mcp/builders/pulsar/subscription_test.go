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
	"strings"
	"testing"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionOperationCoverageIncludesPeekAndGetMessageByID(t *testing.T) {
	require.True(t, isSupportedSubscriptionOperation("peek"))
	require.True(t, isSupportedSubscriptionOperation("get-message-by-id"))
	require.False(t, isReadOnlyRestrictedSubscriptionOperation("peek"))
	require.False(t, isReadOnlyRestrictedSubscriptionOperation("get-message-by-id"))
}

func TestParseSubscriptionMessageID(t *testing.T) {
	messageID, err := parseSubscriptionMessageID("latest")
	require.NoError(t, err)
	require.Equal(t, utils.Latest.String(), messageID.String())

	messageID, err = parseSubscriptionMessageID("earliest")
	require.NoError(t, err)
	require.Equal(t, utils.Earliest.String(), messageID.String())

	messageID, err = parseSubscriptionMessageID("1:2")
	require.NoError(t, err)
	require.Equal(t, "1:2:-1:-1", messageID.String())

	_, err = parseSubscriptionMessageID("invalid")
	require.Error(t, err)
}

func TestRequireInt64Arg(t *testing.T) {
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"count": 3.0,
	}}}
	value, err := requireInt64Arg(req, "count")
	require.NoError(t, err)
	require.Equal(t, int64(3), value)

	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"count": 1.5,
	}}}
	_, err = requireInt64Arg(req, "count")
	require.Error(t, err)

	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{}}}
	_, err = requireInt64Arg(req, "count")
	require.Error(t, err)
}

func TestGetInt64Arg(t *testing.T) {
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{}}}
	value, err := getInt64Arg(req, "count", 1)
	require.NoError(t, err)
	require.Equal(t, int64(1), value)

	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"count": 2.5,
	}}}
	_, err = getInt64Arg(req, "count", 1)
	require.Error(t, err)
}

func TestBuildSubscriptionMessageData(t *testing.T) {
	message := utils.NewMessage(
		"persistent://public/default/example",
		utils.MessageID{LedgerID: 7, EntryID: 9, PartitionIndex: -1, BatchIndex: -1},
		[]byte("hello"),
		map[string]string{"region": "us-west"},
	)

	data := newSubscriptionMessageData(message)
	require.Equal(t, "persistent://public/default/example", data.Topic)
	require.Equal(t, "7:9:-1:-1", data.MessageID)
	require.Equal(t, "hello", data.Payload)
	require.Equal(t, map[string]string{"region": "us-west"}, data.Properties)
	require.True(t, strings.Contains(data.PayloadHex, "68 65 6c 6c 6f"))
}
