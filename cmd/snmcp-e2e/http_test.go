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

package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func TestHTTPWireWithoutInitialize(t *testing.T) {
	calls := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		require.Equal(t, "/mcp", r.URL.Path)
		require.Equal(t, "POST", r.Method)
		require.Equal(t, "Bearer secret", r.Header.Get("Authorization"))
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		require.Equal(t, "application/json, text/event-stream", r.Header.Get("Accept"))
		require.Equal(t, httpProtocolVersion, r.Header.Get("Mcp-Protocol-Version"))
		require.Empty(t, r.Header.Get("Mcp-Session-Id"))
		require.True(t, r.Close, "each request should exercise a fresh connection")
		var body struct {
			JSONRPC string         `json:"jsonrpc"`
			ID      int            `json:"id"`
			Method  string         `json:"method"`
			Params  map[string]any `json:"params"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, "2.0", body.JSONRPC)
		require.Equal(t, 1, body.ID)
		require.Equal(t, body.Method, r.Header.Get("Mcp-Method"))
		require.NotEqual(t, "initialize", body.Method)
		require.Equal(t, map[string]any{"io.modelcontextprotocol/protocolVersion": httpProtocolVersion, "io.modelcontextprotocol/clientCapabilities": map[string]any{}}, body.Params["_meta"])
		switch body.Method {
		case "tools/call":
			require.Equal(t, "echo", r.Header.Get("Mcp-Name"))
		case "resources/read":
			require.Equal(t, "pulsar://admin/v2/clusters", r.Header.Get("Mcp-Name"))
		default:
			require.Empty(t, r.Header.Get("Mcp-Name"))
		}
		writeHTTPTestRPC(w, 200, 0, map[string]any{})
	}))
	defer ts.Close()
	c := newHTTPE2EClient(ts.URL+"/mcp", "secret")
	defer c.client.CloseIdleConnections()
	params := map[string]any{"name": "echo", "arguments": map[string]any{"value": "hi"}}
	for _, method := range []string{"server/discover", "tools/call", "resources/list", "resources/read"} {
		p := map[string]any(nil)
		switch method {
		case "tools/call":
			p = params
		case "resources/read":
			p = map[string]any{"uri": "pulsar://admin/v2/clusters"}
		}
		var result map[string]any
		require.NoError(t, c.rpc(t.Context(), method, p, &result))
	}
	require.Equal(t, 4, calls)
	require.NotContains(t, params, "_meta", "request builder must not mutate caller arguments")
}

func TestHTTPExchangeRejectsFalsePositives(t *testing.T) {
	for _, tc := range []struct {
		name, body, contentType, session string
		actual, expected, code           int
	}{
		{name: "auth unexpectedly succeeds", actual: 200, expected: 401},
		{name: "missing challenge", actual: 401, expected: 401},
		{name: "missing Allow", actual: 405, expected: 405},
		{name: "backend unavailable is not auth", actual: 503, expected: 401},
		{name: "session", actual: 200, expected: 200, session: "secret"},
		{name: "wrong media", actual: 200, expected: 200, body: `{}`},
		{name: "bad JSON", actual: 200, expected: 200, contentType: "application/json", body: `secret`},
		{name: "missing result", actual: 200, expected: 200, contentType: "application/json", body: `{"jsonrpc":"2.0","id":1}`},
		{name: "null result", actual: 200, expected: 200, contentType: "application/json", body: `{"jsonrpc":"2.0","id":1,"result":null}`},
		{name: "wrong id", actual: 200, expected: 200, contentType: "application/json", body: `{"jsonrpc":"2.0","id":2,"result":{}}`},
		{name: "wrong version", actual: 200, expected: 200, contentType: "application/json", body: `{"jsonrpc":"1.0","id":1,"result":{}}`},
		{name: "unexpected RPC error", actual: 200, expected: 200, contentType: "application/json", body: `{"jsonrpc":"2.0","id":1,"error":{"code":-32603,"message":"secret"}}`},
		{name: "wrong RPC error", actual: 400, expected: 400, code: mcp.HEADER_MISMATCH, contentType: "application/json", body: `{"jsonrpc":"2.0","id":1,"error":{"code":-32603,"message":"secret"}}`},
		{name: "error and result", actual: 400, expected: 400, code: mcp.INVALID_PARAMS, contentType: "application/json", body: `{"jsonrpc":"2.0","id":1,"error":{"code":-32602},"result":{}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", tc.contentType)
				if tc.session != "" {
					w.Header().Set("Mcp-Session-Id", tc.session)
				}
				w.WriteHeader(tc.actual)
				_, _ = io.WriteString(w, tc.body)
			}))
			defer ts.Close()
			c := newHTTPE2EClient(ts.URL, "secret")
			defer c.client.CloseIdleConnections()
			req, err := c.request(t.Context(), "tools/list", nil)
			require.NoError(t, err)
			_, err = c.exchange(req, tc.expected, tc.code)
			require.Error(t, err)
			require.NotContains(t, err.Error(), "secret")
		})
	}
}

type failingHTTPTransport struct{}

func (failingHTTPTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("secret transport error")
}

func TestHTTPNetworkFailureAndRedirectAreNotAuth(t *testing.T) {
	c := newHTTPE2EClient("http://example.invalid/mcp", "secret")
	c.client.Transport = failingHTTPTransport{}
	req, err := c.request(t.Context(), "tools/list", nil)
	require.NoError(t, err)
	_, err = c.exchange(req, 401, 0)
	require.Error(t, err)
	require.NotContains(t, err.Error(), "secret")
	redirects := 0
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { redirects++; w.WriteHeader(401) }))
	defer target.Close()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusTemporaryRedirect)
	}))
	defer ts.Close()
	c = newHTTPE2EClient(ts.URL, "secret")
	defer c.client.CloseIdleConnections()
	req, err = c.request(t.Context(), "tools/list", nil)
	require.NoError(t, err)
	_, err = c.exchange(req, 401, 0)
	require.Error(t, err)
	require.Zero(t, redirects)
}

func TestHTTPToolDenialRequiresBackendAuth(t *testing.T) {
	for _, tc := range []struct {
		name, text                 string
		isError, denied, wantError bool
	}{
		{"authorized", "ok", false, false, false},
		{"denied", "AuthorizationError secret", true, true, false},
		{"unexpected success", "secret", false, true, true},
		{"unexpected failure", "secret", true, false, true},
		{"unrelated failure", "broker unavailable secret", true, true, true},
		{"session construction is not authentication", "session not found secret", true, true, true},
		{"empty", "", false, false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				writeHTTPTestRPC(w, 200, 0, map[string]any{"isError": tc.isError, "content": []any{map[string]any{"type": "text", "text": tc.text}}})
			}))
			defer ts.Close()
			c := newHTTPE2EClient(ts.URL, "secret")
			defer c.client.CloseIdleConnections()
			_, err := c.tool(t.Context(), "pulsar_client_produce", nil, tc.denied)
			if tc.wantError {
				require.Error(t, err)
				require.NotContains(t, err.Error(), "secret")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func writeHTTPTestRPC(w http.ResponseWriter, status, code int, result any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	envelope := map[string]any{"jsonrpc": "2.0", "id": 1}
	if code != 0 {
		envelope["error"] = map[string]any{"code": code}
	} else {
		envelope["result"] = result
	}
	_ = json.NewEncoder(w).Encode(envelope)
}

// Test the runner itself, including deliberately bad servers, so a missing
// assertion cannot silently turn the read-only lane green. This fake is not
// claimed as broker integration coverage.
func TestHTTPReadOnlyRunner(t *testing.T) {
	for _, fault := range []string{"", "discovery", "write catalog", "empty clusters", "write accepted", "side effect"} {
		t.Run("fault="+fault, func(t *testing.T) {
			tenant := ""
			reads, rejectedWrites := 0, 0
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.HasSuffix(r.URL.Path, "healthz") || strings.HasSuffix(r.URL.Path, "readyz") {
					w.WriteHeader(200)
					return
				}
				if r.Method != "POST" {
					w.Header().Set("Allow", "POST")
					w.WriteHeader(405)
					return
				}
				var req struct {
					Method string
					Params map[string]json.RawMessage
				}
				require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
				var name string
				_ = json.Unmarshal(req.Params["name"], &name)
				if req.Method == "resources/read" {
					require.NoError(t, json.Unmarshal(req.Params["uri"], &name))
				}
				if req.Method == "tools/call" {
					var args map[string]any
					require.NoError(t, json.Unmarshal(req.Params["arguments"], &args))
					tenant, _ = args["tenant"].(string)
				}
				if r.Header.Get("Origin") != "" {
					w.WriteHeader(403)
					return
				}
				if r.Header.Get("Authorization") == "" {
					w.Header().Set("WWW-Authenticate", "Bearer")
					w.WriteHeader(401)
					return
				}
				var meta map[string]any
				require.NoError(t, json.Unmarshal(req.Params["_meta"], &meta))
				version, _ := meta["io.modelcontextprotocol/protocolVersion"].(string)
				if r.Header.Get("Mcp-Method") != req.Method || r.Header.Get("Mcp-Name") != name || r.Header.Get("Mcp-Protocol-Version") != version {
					writeHTTPTestRPC(w, 400, mcp.HEADER_MISMATCH, nil)
					return
				}
				if version != httpProtocolVersion {
					writeHTTPTestRPC(w, 400, mcp.UNSUPPORTED_PROTOCOL_VERSION, nil)
					return
				}
				switch req.Method {
				case "initialize":
					writeHTTPTestRPC(w, 404, mcp.METHOD_NOT_FOUND, nil)
				case "server/discover":
					version := httpProtocolVersion
					if fault == "discovery" {
						version = "2025-11-25"
					}
					writeHTTPTestRPC(w, 200, 0, map[string]any{"supportedVersions": []string{version}, "resultType": "complete", "cacheScope": "private", "ttlMs": 0, "_meta": map[string]any{"io.modelcontextprotocol/serverInfo": map[string]any{"name": "test"}}})
				case "tools/list":
					tools := []any{map[string]any{"name": "pulsar_admin_cluster_read", "annotations": map[string]any{"readOnlyHint": true}}}
					if fault == "write catalog" {
						tools = append(tools, map[string]any{"name": "pulsar_client_produce", "annotations": map[string]any{"readOnlyHint": true}})
					}
					writeHTTPTestRPC(w, 200, 0, map[string]any{"tools": tools})
				case "resources/list":
					writeHTTPTestRPC(w, 200, 0, map[string]any{"resources": []any{map[string]any{"uri": "pulsar://admin/v2/clusters"}}})
				case "resources/read":
					reads++
					var uri string
					require.NoError(t, json.Unmarshal(req.Params["uri"], &uri))
					payload := map[string]any{"tenants": []string{"public"}}
					if uri == "pulsar://admin/v2/clusters" {
						payload = map[string]any{"clusters": []string{"standalone"}}
						if fault == "empty clusters" {
							payload["clusters"] = []string{}
						}
					} else if fault == "side effect" {
						payload["tenants"] = []string{"public", tenant}
					}
					raw, err := json.Marshal(payload)
					require.NoError(t, err)
					writeHTTPTestRPC(w, 200, 0, map[string]any{"contents": []any{map[string]any{"uri": uri, "text": string(raw)}}})
				case "tools/call":
					rejectedWrites++
					if fault == "write accepted" {
						writeHTTPTestRPC(w, 200, 0, map[string]any{})
						return
					}
					writeHTTPTestRPC(w, 400, mcp.INVALID_PARAMS, nil)
				default:
					t.Errorf("unexpected RPC method %s", req.Method)
					w.WriteHeader(500)
				}
			}))
			defer ts.Close()
			err := run(t.Context(), config{httpBaseURL: ts.URL + "/mcp", transport: "http", readOnly: true, adminToken: "admin", testUserToken: "test-user"})
			if fault != "" {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, 3, reads)
				require.Equal(t, 1, rejectedWrites)
			}
		})
	}
}

func TestTransportConfig(t *testing.T) {
	t.Setenv("ADMIN_TOKEN", "admin")
	t.Setenv("TEST_USER_TOKEN", "user")
	t.Setenv("E2E_TRANSPORT", "")
	t.Setenv("E2E_READ_ONLY", "false")
	cfg, err := parseConfigArgs(nil)
	require.NoError(t, err)
	require.Equal(t, "sse", cfg.transport)
	require.False(t, cfg.readOnly)
	t.Setenv("E2E_TRANSPORT", "http")
	t.Setenv("E2E_READ_ONLY", "true")
	cfg, err = parseConfigArgs(nil)
	require.NoError(t, err)
	require.Equal(t, "http", cfg.transport)
	require.True(t, cfg.readOnly)
	cfg, err = parseConfigArgs([]string{"--transport=sse", "--read-only=false"})
	require.NoError(t, err)
	require.Equal(t, "sse", cfg.transport)
	require.False(t, cfg.readOnly)
	for _, args := range [][]string{{"--transport=unknown"}, {"--transport=sse"}, {"--read-only=secret"}} {
		_, err := parseConfigArgs(args)
		require.Error(t, err)
		require.NotContains(t, err.Error(), "secret")
	}
	t.Setenv("E2E_TRANSPORT", "unknown")
	_, err = parseConfigArgs(nil)
	require.Error(t, err)
}

func TestHTTPConsumeRequiresExpectedPayload(t *testing.T) {
	require.NoError(t, assertHTTPConsumed(`{"messages_consumed":1,"messages":[{"data":"expected"}]}`, "expected"))
	for _, raw := range []string{
		`{}`, `null`, `secret`,
		`{"messages_consumed":0,"messages":[]}`,
		`{"messages_consumed":1,"messages":[]}`,
		`{"messages_consumed":1,"messages":[{"data":"secret"}]}`,
		`{"messages_consumed":2,"messages":[{"data":"expected"},{"data":"secret"}]}`,
	} {
		err := assertHTTPConsumed(raw, "expected")
		require.Error(t, err)
		require.NotContains(t, err.Error(), "secret")
	}
}
