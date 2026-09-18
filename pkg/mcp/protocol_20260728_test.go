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
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"

	pulsaradminconfig "github.com/apache/pulsar-client-go/pulsaradmin/pkg/admin/config"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	pulsarsession "github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Exercise the application constructor and registrations over actual newline-
// delimited stdio, not SDK fixture handlers or typed responses that bypass JSON.
// Do not parallelize: the SDK stdio transport has a process-wide session.
func TestProtocol20260728Stdio(t *testing.T) {
	var reads atomic.Int32
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/admin/v2/tenants" {
			t.Errorf("unexpected Pulsar request: %s %s", r.Method, r.URL.Path)
			http.Error(w, "unexpected request", http.StatusBadRequest)
			return
		}
		reads.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`["public","system"]`))
	}))
	t.Cleanup(backend.Close)
	cfg := &cmdutils.ClusterConfig{WebServiceURL: backend.URL}
	session := &pulsarsession.Session{
		Ctx:             pulsarsession.PulsarContext{WebServiceURL: backend.URL},
		PulsarCtlConfig: cfg,
		AdminClient:     cfg.Client(pulsaradminconfig.V2),
	}
	srv := NewServer("streamnative-mcp-server", "protocol-test", logrus.New(),
		server.WithInstructions("Read-only Pulsar tenant inspection"))
	features := []string{string(FeaturePulsarAdminTenants)}
	PulsarAddResources(srv.MCPServer, features)
	PulsarAdminAddTenantTools(srv.MCPServer, true, features)
	call := protocol20260728Pipe(t, srv.MCPServer, session)

	t.Run("discover without initialize", func(t *testing.T) {
		result := protocol20260728Result(t, call(t, "server/discover", nil, protocol20260728Meta()))
		versions, ok := result["supportedVersions"].([]any)
		require.True(t, ok)
		require.NotEmpty(t, versions)
		assert.Equal(t, "2026-07-28", versions[0])
		assert.Equal(t, "Read-only Pulsar tenant inspection", result["instructions"])
		caps := protocol20260728Object(t, result["capabilities"])
		assert.Contains(t, caps, "resources")
		assert.Contains(t, caps, "tools")
		assert.Contains(t, caps, "logging")
		assert.NotContains(t, caps, "tasks")
		assert.NotContains(t, caps, "completions")
		assert.NotContains(t, result, "protocolVersion", "discovery is not legacy negotiation")
	})

	t.Run("registered read-only surface", func(t *testing.T) {
		result := protocol20260728Result(t, call(t, "resources/list", nil, protocol20260728Meta()))
		var listed struct {
			Resources []mcp.Resource `json:"resources"`
		}
		protocol20260728Decode(t, result, &listed)
		uris := make([]string, 0, len(listed.Resources))
		for _, resource := range listed.Resources {
			uris = append(uris, resource.URI)
		}
		assert.ElementsMatch(t, []string{"pulsar://context", "pulsar://resources", "pulsar://admin/v2/tenants"}, uris)

		result = protocol20260728Result(t, call(t, "resources/templates/list", nil, protocol20260728Meta()))
		var templates struct {
			Templates []struct {
				URI string `json:"uriTemplate"`
			} `json:"resourceTemplates"`
		}
		protocol20260728Decode(t, result, &templates)
		require.Len(t, templates.Templates, 1)
		assert.Equal(t, "pulsar://admin/v2/tenants/{tenant}", templates.Templates[0].URI)

		result = protocol20260728Result(t, call(t, "tools/list", nil, protocol20260728Meta()))
		tools, ok := result["tools"].([]any)
		require.True(t, ok)
		require.Len(t, tools, 1)
		tool := protocol20260728Object(t, tools[0])
		assert.Equal(t, "pulsar_admin_tenant_read", tool["name"])
		assert.Equal(t, true, protocol20260728Object(t, tool["annotations"])["readOnlyHint"])
		properties := protocol20260728Object(t, protocol20260728Object(t, tool["inputSchema"])["properties"])
		assert.ElementsMatch(t, []any{"list", "get"}, protocol20260728Object(t, properties["operation"])["enum"])
	})

	resourceParams := map[string]any{"uri": "pulsar://admin/v2/tenants"}
	toolParams := map[string]any{"name": "pulsar_admin_tenant_read", "arguments": map[string]any{"resource": "tenant", "operation": "list"}}
	t.Run("real resource read and tool call", func(t *testing.T) {
		result := protocol20260728Result(t, call(t, "resources/read", resourceParams, protocol20260728Meta()))
		var resource struct {
			Contents []mcp.TextResourceContents `json:"contents"`
		}
		protocol20260728Decode(t, result, &resource)
		require.Len(t, resource.Contents, 1)
		content := resource.Contents[0]
		assert.Equal(t, "pulsar://admin/v2/tenants", content.URI)
		assert.Equal(t, "application/json", content.MIMEType)
		var payload pulsarTenantCollectionResource
		require.NoError(t, json.Unmarshal([]byte(content.Text), &payload))
		assert.ElementsMatch(t, []string{"public", "system"}, payload.Tenants)
		assert.Equal(t, 2, payload.Count)

		result = protocol20260728Result(t, call(t, "tools/call", toolParams, protocol20260728Meta()))
		assert.NotEqual(t, true, result["isError"])
		contents, ok := result["content"].([]any)
		require.True(t, ok)
		require.Len(t, contents, 1)
		text := protocol20260728Object(t, contents[0])
		assert.Equal(t, "text", text["type"])
		require.IsType(t, "", text["text"])
		assert.JSONEq(t, `["public","system"]`, text["text"].(string))
		assert.EqualValues(t, 2, reads.Load())
	})

	t.Run("metadata is required on every request", func(t *testing.T) {
		before := reads.Load()
		for _, request := range []struct {
			method string
			params map[string]any
		}{
			{"resources/read", resourceParams}, {"tools/call", toolParams},
		} {
			meta := protocol20260728Meta()
			delete(meta, mcp.MetaKeyClientCapabilities)
			response := call(t, request.method, request.params, meta)
			protocol20260728Error(t, response, mcp.MISSING_REQUIRED_CLIENT_CAPABILITY)
		}
		assert.Equal(t, before, reads.Load(), "invalid requests must not reach Pulsar")
	})

	t.Run("removed methods", func(t *testing.T) {
		for _, method := range []string{"initialize", "ping", "logging/setLevel", "resources/subscribe", "resources/unsubscribe"} {
			t.Run(method, func(t *testing.T) {
				response := call(t, method, map[string]any{"uri": "pulsar://admin/v2/tenants", "level": "info"}, protocol20260728Meta())
				protocol20260728Error(t, response, mcp.METHOD_NOT_FOUND)
				assert.Contains(t, protocol20260728Object(t, response["error"])["message"], "removed in protocol version")
			})
		}
	})

	t.Run("unknown version fails without executing tool", func(t *testing.T) {
		before := reads.Load()
		meta := protocol20260728Meta()
		meta[mcp.MetaKeyProtocolVersion] = "2099-01-01"
		protocol20260728Error(t, call(t, "tools/call", toolParams, meta), mcp.UNSUPPORTED_PROTOCOL_VERSION)
		assert.Equal(t, before, reads.Load())
		// A rejected request must not poison subsequent traffic on the same pipe.
		protocol20260728Result(t, call(t, "tools/call", toolParams, protocol20260728Meta()))
		assert.Equal(t, before+1, reads.Load())
	})
}

func protocol20260728Meta() map[string]any {
	return map[string]any{
		mcp.MetaKeyProtocolVersion:    "2026-07-28",
		mcp.MetaKeyClientInfo:         map[string]any{"name": "app-protocol-test", "version": "1"},
		mcp.MetaKeyClientCapabilities: map[string]any{},
	}
}

func protocol20260728Object(t *testing.T, value any) map[string]any {
	t.Helper()
	object, ok := value.(map[string]any)
	require.True(t, ok, "expected JSON object, got %T: %v", value, value)
	return object
}

func protocol20260728Result(t *testing.T, response map[string]any) map[string]any {
	t.Helper()
	require.NotContains(t, response, "error")
	result := protocol20260728Object(t, response["result"])
	assert.Equal(t, "complete", result["resultType"])
	meta := protocol20260728Object(t, result["_meta"])
	info := protocol20260728Object(t, meta[mcp.MetaKeyServerInfo])
	assert.Equal(t, "streamnative-mcp-server", info["name"])
	assert.Equal(t, "protocol-test", info["version"])
	assert.NotContains(t, result, "serverInfo", "modern identity belongs in _meta")
	return result
}

func protocol20260728Error(t *testing.T, response map[string]any, code int) {
	t.Helper()
	require.NotContains(t, response, "result")
	assert.EqualValues(t, code, protocol20260728Object(t, response["error"])["code"])
}

func protocol20260728Decode(t *testing.T, value, destination any) {
	t.Helper()
	data, err := json.Marshal(value)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(data, destination))
}

func protocol20260728Pipe(t *testing.T, srv *server.MCPServer, session *pulsarsession.Session) func(*testing.T, string, map[string]any, map[string]any) map[string]any {
	t.Helper()
	input, writer, err := os.Pipe()
	require.NoError(t, err)
	t.Cleanup(func() { _ = input.Close(); _ = writer.Close() })
	reader, output, err := os.Pipe()
	require.NoError(t, err)
	t.Cleanup(func() { _ = reader.Close(); _ = output.Close() })
	ctx, cancel := context.WithCancel(t.Context())
	transport := server.NewStdioServer(srv)
	transport.SetContextFunc(func(ctx context.Context) context.Context {
		return mcpCtx.WithPulsarSession(ctx, session)
	})
	done := make(chan error, 1)
	go func() { done <- transport.Listen(ctx, input, output) }()
	t.Cleanup(func() {
		_ = writer.Close()
		cancel()
		select {
		case err := <-done:
			if !errors.Is(err, context.Canceled) {
				assert.NoError(t, err)
			}
		case <-time.After(5 * time.Second):
			t.Error("stdio server did not stop")
		}
	})
	encoder, decoder := json.NewEncoder(writer), json.NewDecoder(reader)
	id := 0
	return func(t *testing.T, method string, params, meta map[string]any) map[string]any {
		t.Helper()
		id++
		wireParams := make(map[string]any, len(params)+1)
		for key, value := range params {
			wireParams[key] = value
		}
		wireParams["_meta"] = meta
		require.NoError(t, writer.SetWriteDeadline(time.Now().Add(10*time.Second)))
		require.NoError(t, reader.SetReadDeadline(time.Now().Add(10*time.Second)))
		require.NoError(t, encoder.Encode(map[string]any{"jsonrpc": "2.0", "id": id, "method": method, "params": wireParams}))
		var response map[string]any
		require.NoError(t, decoder.Decode(&response), "method %s", method)
		require.Equal(t, "2.0", response["jsonrpc"])
		require.EqualValues(t, id, response["id"])
		if result, ok := response["result"].(map[string]any); ok {
			switch method {
			case "server/discover", "tools/list", "resources/list", "resources/templates/list", "resources/read":
				assert.EqualValues(t, 0, result["ttlMs"], "private backend data must be revalidated")
				assert.Equal(t, "private", result["cacheScope"])
			}
		}
		return response
	}
}
