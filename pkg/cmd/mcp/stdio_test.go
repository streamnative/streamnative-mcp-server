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
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync/atomic"
	"testing"
	"time"

	protocol "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	app "github.com/streamnative/streamnative-mcp-server/pkg/mcp"
	"github.com/stretchr/testify/require"
)

func TestStdioProtocolProfiles(t *testing.T) {
	for _, tc := range []struct {
		name         string
		options      config.Options
		multiSession bool
		modern       bool
	}{
		{name: "Cloud", options: config.Options{KeyFile: "not-read.json"}},
		{name: "fixed Cloud", options: config.Options{KeyFile: "not-read.json", PulsarInstance: "instance", PulsarCluster: "cluster"}},
		{name: "unconfigured"},
		{name: "ambiguous backends", options: config.Options{UseExternalKafka: true, UseExternalPulsar: true}},
		{name: "Cloud takes precedence", options: config.Options{KeyFile: "not-read.json", UseExternalPulsar: true}},
		{name: "multi-session", options: config.Options{UseExternalPulsar: true}, multiSession: true},
		{name: "fixed Kafka", options: config.Options{UseExternalKafka: true}, modern: true},
		{name: "fixed Pulsar", options: config.Options{UseExternalPulsar: true}, modern: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			opts := NewMcpServerOptions(&tc.options)
			opts.MultiSessionPulsar = tc.multiSession
			s := app.NewServer("profile-test", "1", logrus.New(), stdioServerOptions(opts)...)
			var calls atomic.Int32
			var reads atomic.Int32
			s.MCPServer.AddResource(protocol.NewResource("test://value", "test"), func(context.Context, protocol.ReadResourceRequest) ([]protocol.ResourceContents, error) {
				reads.Add(1)
				return []protocol.ResourceContents{protocol.TextResourceContents{URI: "test://value", Text: "value"}}, nil
			})
			s.MCPServer.AddTool(protocol.NewTool("mutate"), func(context.Context, protocol.CallToolRequest) (*protocol.CallToolResult, error) {
				calls.Add(1)
				return protocol.NewToolResultText("done"), nil
			})
			call := stdioWireCall(t, s.MCPServer)
			params := map[string]any{"name": "mutate"}
			// Direct calls must be gated before discovery or initialization, too.
			direct := call("tools/call", params, true)
			if tc.modern {
				require.NotContains(t, direct, "error")
				require.EqualValues(t, 1, calls.Load())
			} else {
				require.Contains(t, direct, "error")
				require.Contains(t, direct["error"].(map[string]any)["message"], "legacy")
				require.Zero(t, calls.Load())
			}
			resource := call("resources/read", map[string]any{"uri": "test://value"}, true)
			if tc.modern {
				require.NotContains(t, resource, "error")
				require.EqualValues(t, 1, reads.Load())
			} else {
				require.Contains(t, resource, "error")
				require.Zero(t, reads.Load())
			}
			discover := call("server/discover", nil, true)
			require.NotContains(t, discover, "error")
			versions := discover["result"].(map[string]any)["supportedVersions"]
			if tc.modern {
				require.Contains(t, versions, protocol.ProtocolVersion20260728)
			} else {
				require.NotContains(t, versions, protocol.ProtocolVersion20260728)
			}
			resources := discover["result"].(map[string]any)["capabilities"].(map[string]any)["resources"].(map[string]any)
			if tc.modern {
				require.NotEqual(t, true, resources["subscribe"])
				require.NotEqual(t, true, resources["listChanged"])
			}
			subscriptions := call("subscriptions/listen", map[string]any{"notifications": map[string]any{
				"resourcesListChanged": true, "resourceSubscriptions": []string{"pulsar://context"},
			}}, true)
			if tc.modern {
				require.NotContains(t, subscriptions, "error")
			} else {
				require.Contains(t, subscriptions, "error")
			}
			initialize := call("initialize", map[string]any{
				"protocolVersion": protocol.LATEST_LEGACY_PROTOCOL_VERSION,
				"clientInfo":      map[string]any{"name": "test", "version": "1"}, "capabilities": map[string]any{},
			}, false)
			require.NotContains(t, initialize, "error")
			require.Equal(t, protocol.LATEST_LEGACY_PROTOCOL_VERSION, initialize["result"].(map[string]any)["protocolVersion"])
			legacyResources := initialize["result"].(map[string]any)["capabilities"].(map[string]any)["resources"].(map[string]any)
			require.Equal(t, true, legacyResources["subscribe"])
			require.Equal(t, true, legacyResources["listChanged"])
			before := calls.Load()
			require.NotContains(t, call("tools/call", params, false), "error")
			require.Equal(t, before+1, calls.Load(), "legacy remains usable after modern rejection")
			afterLegacy := calls.Load()
			again := call("tools/call", params, true)
			if tc.modern {
				require.NotContains(t, again, "error")
				require.Equal(t, afterLegacy+1, calls.Load())
			} else {
				require.Contains(t, again, "error")
				require.Equal(t, afterLegacy, calls.Load(), "legacy initialization must not bypass the modern gate")
			}
		})
	}
}

// Use actual stdio framing and workers, not just direct HandleMessage calls.
func stdioWireCall(t *testing.T, s *server.MCPServer) func(string, map[string]any, bool) map[string]any {
	t.Helper()
	input, writer, err := os.Pipe()
	require.NoError(t, err)
	reader, output, err := os.Pipe()
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan error, 1)
	go func() { done <- server.NewStdioServer(s).Listen(ctx, input, output) }()
	t.Cleanup(func() {
		cancel()
		_ = writer.Close()
		select {
		case err := <-done:
			if !errors.Is(err, context.Canceled) {
				require.NoError(t, err)
			}
		case <-time.After(5 * time.Second):
			t.Error("stdio server did not stop")
		}
		_ = input.Close()
		_ = output.Close()
		_ = reader.Close()
	})
	encoder, decoder := json.NewEncoder(writer), json.NewDecoder(reader)
	id := 0
	return func(method string, params map[string]any, modern bool) map[string]any {
		t.Helper()
		id++
		p := make(map[string]any, len(params)+1)
		for k, v := range params {
			p[k] = v
		}
		if modern {
			p["_meta"] = map[string]any{protocol.MetaKeyProtocolVersion: protocol.ProtocolVersion20260728, protocol.MetaKeyClientCapabilities: map[string]any{}}
		}
		require.NoError(t, writer.SetWriteDeadline(time.Now().Add(5*time.Second)))
		require.NoError(t, reader.SetReadDeadline(time.Now().Add(5*time.Second)))
		require.NoError(t, encoder.Encode(map[string]any{"jsonrpc": "2.0", "id": id, "method": method, "params": p}))
		var response map[string]any
		for {
			response = nil
			require.NoError(t, decoder.Decode(&response))
			if response["method"] != "notifications/subscriptions/acknowledged" {
				break
			}
			require.Empty(t, response["params"].(map[string]any)["notifications"], "unsupported subscriptions must not be established")
		}
		require.EqualValues(t, id, response["id"])
		return response
	}
}
