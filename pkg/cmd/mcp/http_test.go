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
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	protocol "github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	app "github.com/streamnative/streamnative-mcp-server/pkg/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
	"github.com/stretchr/testify/require"
)

func TestHTTPRequestScopedSSEHoldsSessionUntilToolCompletes(t *testing.T) {
	s := app.NewServer("http-test", "1", logrus.New())
	finish := make(chan struct{})
	defer close(finish)
	var released atomic.Bool
	s.MCPServer.AddTool(protocol.NewTool("stream"), func(ctx context.Context, _ protocol.CallToolRequest) (*protocol.CallToolResult, error) {
		if err := s.MCPServer.SendNotificationToClient(ctx, "notifications/progress", map[string]any{"progressToken": "p", "progress": 1}); err != nil {
			return nil, err
		}
		select {
		case <-finish:
		case <-ctx.Done():
		}
		return protocol.NewToolResultText("done"), nil
	})
	handler, err := newHTTPHandler(s, nil, "/mcp", nil, func(context.Context, string) (*pulsar.Session, func(), error) {
		return &pulsar.Session{}, func() { released.Store(true) }, nil
	})
	require.NoError(t, err)
	ts := httptest.NewServer(handler)
	defer ts.Close()
	req := modernHTTPRequest(t, ts.URL, "tools/call")
	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	body, err := io.ReadAll(req.Body)
	require.NoError(t, err)
	body = []byte(strings.Replace(string(body), `"params":{`, `"params":{"name":"stream",`, 1))
	req.Body = io.NopCloser(strings.NewReader(string(body)))
	req.ContentLength = int64(len(body))
	req.Header.Set("Mcp-Name", "stream")
	req.Header.Set("Authorization", "Bearer alice")
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	require.Equal(t, "text/event-stream", resp.Header.Get("Content-Type"))
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		require.NoError(t, err)
		if strings.Contains(line, "notifications/progress") {
			break
		}
	}
	require.False(t, released.Load(), "session must survive the streaming response")
	finish <- struct{}{}
	data, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Contains(t, string(data), "done")
	require.True(t, released.Load())
}

func TestHTTPRunPropagatesListenerErrorsAndStopsOnContext(t *testing.T) {
	cfg := config.NewConfigOptions()
	cfg.UseExternalPulsar = true
	cfg.Pulsar.WebServiceURL = "http://127.0.0.1:8080"
	opts := NewMcpServerOptions(cfg)
	opts.MultiSessionPulsar = true
	opts.Features = []string{string(app.FeatureAllPulsar)}
	opts.HTTPPath = "/mcp"
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer func() { _ = listener.Close() }()
	opts.HTTPAddr = listener.Addr().String()
	require.ErrorContains(t, runHTTPServer(t.Context(), opts, nil), "listen HTTP")
	opts.HTTPAddr = "127.0.0.1:0"
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	require.NoError(t, runHTTPServer(ctx, opts, nil))
}

func TestHTTPRejectsLegacyAndRemovedRequests(t *testing.T) {
	s := app.NewServer("http-test", "1", logrus.New())
	handler, err := newHTTPHandler(s, nil, "/mcp", nil, nil)
	require.NoError(t, err)
	ts := httptest.NewServer(handler)
	defer ts.Close()
	for _, tc := range []struct {
		version, method string
		status, code    int
	}{
		{"2025-11-25", "initialize", 400, protocol.UNSUPPORTED_PROTOCOL_VERSION},
		{"2099-01-01", "tools/list", 400, protocol.UNSUPPORTED_PROTOCOL_VERSION},
		{"2026-07-28", "initialize", 404, protocol.METHOD_NOT_FOUND},
		{"2026-07-28", "ping", 404, protocol.METHOD_NOT_FOUND},
	} {
		req := modernHTTPRequest(t, ts.URL, tc.method)
		body, err := io.ReadAll(req.Body)
		require.NoError(t, err)
		body = []byte(strings.ReplaceAll(string(body), "2026-07-28", tc.version))
		req.Body = io.NopCloser(strings.NewReader(string(body)))
		req.ContentLength = int64(len(body))
		req.Header.Set("Mcp-Protocol-Version", tc.version)
		resp, err := ts.Client().Do(req)
		require.NoError(t, err)
		require.Equal(t, tc.status, resp.StatusCode)
		var msg struct {
			Error struct {
				Code int
				Data protocol.UnsupportedProtocolVersionData
			}
		}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&msg))
		require.NoError(t, resp.Body.Close())
		require.Equal(t, tc.code, msg.Error.Code)
		if tc.code == protocol.UNSUPPORTED_PROTOCOL_VERSION {
			require.Equal(t, []string{"2026-07-28"}, msg.Error.Data.Supported)
			require.Equal(t, tc.version, msg.Error.Data.Requested)
		}
	}
}

func TestHTTPDirectToolCallsIsolateCredentialsAndReleaseAfterResponse(t *testing.T) {
	s := app.NewServer("http-test", "1", logrus.New())
	var released atomic.Int32
	s.MCPServer.AddTool(protocol.NewTool("identity"), func(ctx context.Context, _ protocol.CallToolRequest) (*protocol.CallToolResult, error) {
		ps := app.GetPulsarSession(ctx)
		if ps == nil {
			return nil, errors.New("missing session")
		}
		return protocol.NewToolResultText(ps.Ctx.Token), nil
	})
	resolve := func(_ context.Context, token string) (*pulsar.Session, func(), error) {
		return &pulsar.Session{Ctx: pulsar.PulsarContext{Token: token}}, func() { released.Add(1) }, nil
	}
	handler, err := newHTTPHandler(s, nil, "/mcp", nil, resolve)
	require.NoError(t, err)
	ts := httptest.NewServer(handler)
	defer ts.Close()
	for _, token := range []string{"alice", "bob", "alice"} {
		req := modernHTTPRequest(t, ts.URL, "tools/call")
		body, err := io.ReadAll(req.Body)
		require.NoError(t, err)
		body = []byte(strings.Replace(string(body), `"params":{`, `"params":{"name":"identity",`, 1))
		req.Body = io.NopCloser(strings.NewReader(string(body)))
		req.ContentLength = int64(len(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Mcp-Name", "identity")
		resp, err := ts.Client().Do(req)
		require.NoError(t, err)
		data, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.NoError(t, resp.Body.Close())
		require.Equal(t, 200, resp.StatusCode, "%s", data)
		require.Contains(t, string(data), `"text":"`+token+`"`)
		require.Empty(t, resp.Header.Get("Mcp-Session-Id"))
	}
	require.Equal(t, int32(3), released.Load())
}

func TestHTTPWireValidationBeforeSessionResolution(t *testing.T) {
	s := app.NewServer("http-test", "1", logrus.New())
	calls := 0
	handler, err := newHTTPHandler(s, nil, "/mcp", nil, func(context.Context, string) (*pulsar.Session, func(), error) {
		calls++
		return &pulsar.Session{}, func() {}, nil
	})
	require.NoError(t, err)
	ts := httptest.NewServer(handler)
	defer ts.Close()
	for _, tc := range []struct {
		header, value string
		code          int
	}{
		{"Mcp-Protocol-Version", "2025-11-25", protocol.HEADER_MISMATCH},
		{"Mcp-Protocol-Version", "", protocol.HEADER_MISMATCH},
		{"Mcp-Method", "tools/call", protocol.HEADER_MISMATCH},
	} {
		req := modernHTTPRequest(t, ts.URL, "tools/list")
		req.Header.Set("Authorization", "Bearer a")
		req.Header.Set(tc.header, tc.value)
		resp, err := ts.Client().Do(req)
		require.NoError(t, err)
		require.Equal(t, 400, resp.StatusCode)
		var msg struct{ Error struct{ Code int } }
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&msg))
		require.NoError(t, resp.Body.Close())
		require.Equal(t, tc.code, msg.Error.Code)
	}
	require.Zero(t, calls)
}

func TestHTTPMissingCapabilitiesDoesNotAcquireBackend(t *testing.T) {
	s := app.NewServer("http-test", "1", logrus.New())
	calls := 0
	handler, err := newHTTPHandler(s, nil, "/mcp", nil, func(context.Context, string) (*pulsar.Session, func(), error) {
		calls++
		return &pulsar.Session{}, func() {}, nil
	})
	require.NoError(t, err)
	ts := httptest.NewServer(handler)
	defer ts.Close()
	req := modernHTTPRequest(t, ts.URL, "tools/list")
	body, err := io.ReadAll(req.Body)
	require.NoError(t, err)
	body = []byte(strings.Replace(string(body), `,"io.modelcontextprotocol/clientCapabilities":{}`, "", 1))
	req.Body = io.NopCloser(strings.NewReader(string(body)))
	req.ContentLength = int64(len(body))
	req.Header.Set("Authorization", "Bearer alice")
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	var message struct{ Error struct{ Code int } }
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&message))
	require.Equal(t, protocol.MISSING_REQUIRED_CLIENT_CAPABILITY, message.Error.Code)
	require.NoError(t, resp.Body.Close())
	require.Equal(t, 400, resp.StatusCode)
	require.Zero(t, calls)
}

func TestHTTPMalformedCapabilities(t *testing.T) {
	s := app.NewServer("http-test", "1", logrus.New())
	handler, err := newHTTPHandler(s, nil, "/mcp", nil, nil)
	require.NoError(t, err)
	ts := httptest.NewServer(handler)
	defer ts.Close()
	req := modernHTTPRequest(t, ts.URL, "tools/list")
	body, err := io.ReadAll(req.Body)
	require.NoError(t, err)
	body = []byte(strings.Replace(string(body), `"io.modelcontextprotocol/clientCapabilities":{}`, `"io.modelcontextprotocol/clientCapabilities":"invalid"`, 1))
	req.Body = io.NopCloser(strings.NewReader(string(body)))
	req.ContentLength = int64(len(body))
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	var message struct{ Error struct{ Code int } }
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&message))
	require.Equal(t, protocol.HEADER_MISMATCH, message.Error.Code)
}

func TestHTTPCommandRejectsUnsafeConfigurationBeforeAuth(t *testing.T) {
	for _, args := range [][]string{
		{"--key-file", "/nonexistent"},
		{"--use-external-kafka", "--kafka-bootstrap-servers", "localhost:9092", "--http-addr", ":9090"},
		{"--use-external-kafka", "--kafka-bootstrap-servers", "localhost:9092", "--multi-session-pulsar"},
	} {
		viper.Reset()
		cfg := config.NewConfigOptions()
		root := &cobra.Command{Use: "test", SilenceErrors: true, SilenceUsage: true}
		cfg.AddFlags(root)
		opts := NewMcpServerOptions(cfg)
		opts.AddFlags(root)
		root.AddCommand(NewCmdMcpHTTPServer(opts))
		root.SetArgs(append([]string{"http"}, args...))
		err := root.ExecuteContext(t.Context())
		require.Error(t, err)
		require.Contains(t, err.Error(), "http transport")
		require.Nil(t, cfg.Store, "must reject before auth/keyring setup")
	}
	viper.Reset()
}

func TestHTTPRejectsLegacyMethodsAndOversizedBody(t *testing.T) {
	s := app.NewServer("http-test", "1", logrus.New())
	handler, err := newHTTPHandler(s, nil, "/mcp", nil, nil)
	require.NoError(t, err)
	ts := httptest.NewServer(handler)
	defer ts.Close()
	for _, method := range []string{http.MethodGet, http.MethodDelete} {
		req, err := http.NewRequestWithContext(t.Context(), method, ts.URL+"/mcp", nil)
		require.NoError(t, err)
		resp, err := ts.Client().Do(req)
		require.NoError(t, err)
		require.Equal(t, 405, resp.StatusCode)
		require.Equal(t, "POST", resp.Header.Get("Allow"))
		require.NoError(t, resp.Body.Close())
	}
	req := modernHTTPRequest(t, ts.URL, "tools/list")
	req.Body = io.NopCloser(strings.NewReader(strings.Repeat("x", 4*1024*1024+1)))
	req.ContentLength = -1
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	require.Equal(t, 413, resp.StatusCode)
	require.NoError(t, resp.Body.Close())
}

func TestHTTPOriginAndCredentialsFailClosed(t *testing.T) {
	s := app.NewServer("http-test", "1", logrus.New())
	calls := 0
	resolve := func(context.Context, string) (*pulsar.Session, func(), error) {
		calls++
		return nil, nil, errors.New("secret backend error")
	}
	handler, err := newHTTPHandler(s, nil, "/mcp", []string{"https://trusted.example"}, resolve)
	require.NoError(t, err)
	ts := httptest.NewServer(handler)
	defer ts.Close()
	for _, tc := range []struct {
		origin, auth string
		status       int
	}{
		{"null", "Bearer a", 403}, {"https://evil.example", "Bearer a", 403},
		{"https://trusted.example/", "Bearer a", 403}, {"", "", 401},
		{"", "Bearer ", 401}, {"", "Basic a", 401},
		{"https://trusted.example", "Bearer a", 503},
	} {
		req := modernHTTPRequest(t, ts.URL, "server/discover")
		if tc.origin != "" {
			req.Header.Set("Origin", tc.origin)
		}
		req.Header.Set("Authorization", tc.auth)
		resp, err := ts.Client().Do(req)
		require.NoError(t, err)
		require.Equal(t, tc.status, resp.StatusCode, "%+v", tc)
		require.NoError(t, resp.Body.Close())
	}
	require.Equal(t, 1, calls)
	for _, path := range []string{"healthz", "readyz"} {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, ts.URL+"/mcp/"+path, nil)
		require.NoError(t, err)
		resp, err := ts.Client().Do(req)
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)
		require.NoError(t, resp.Body.Close())
	}
	require.Equal(t, 1, calls)
}

func TestHTTPDiscoverWithoutInitialization(t *testing.T) {
	s := app.NewServer("http-test", "1", logrus.New())
	handler, err := newHTTPHandler(s, nil, "/mcp", nil, nil)
	require.NoError(t, err)
	ts := httptest.NewServer(handler)
	defer ts.Close()
	req := modernHTTPRequest(t, ts.URL, "server/discover")
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Empty(t, resp.Header.Get("Mcp-Session-Id"))
	var body map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	result := body["result"].(map[string]any)
	require.Equal(t, []any{"2026-07-28"}, result["supportedVersions"])
	require.Equal(t, "complete", result["resultType"])
	require.Equal(t, "private", result["cacheScope"])
	require.Equal(t, float64(0), result["ttlMs"])
	require.Contains(t, result["_meta"], "io.modelcontextprotocol/serverInfo")
	resources := result["capabilities"].(map[string]any)["resources"].(map[string]any)
	require.NotEqual(t, true, resources["subscribe"])
	require.NotEqual(t, true, resources["listChanged"])
}

func TestHTTPDeclinesUnsupportedResourceSubscriptions(t *testing.T) {
	s := app.NewServer("http-test", "1", logrus.New())
	handler, err := newHTTPHandler(s, nil, "/mcp", nil, nil)
	require.NoError(t, err)
	req := modernHTTPRequest(t, "http://localhost", "subscriptions/listen")
	var body map[string]any
	require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
	body["params"].(map[string]any)["notifications"] = map[string]any{
		"resourcesListChanged": true, "resourceSubscriptions": []string{"pulsar://context"},
	}
	data, err := json.Marshal(body)
	require.NoError(t, err)
	req.Body = io.NopCloser(strings.NewReader(string(data)))
	req.ContentLength = int64(len(data))
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req.WithContext(ctx))
	require.Equal(t, http.StatusOK, recorder.Code)
	require.NoError(t, ctx.Err(), "unsupported subscriptions must complete without waiting for cancellation")
	require.Contains(t, recorder.Body.String(), `"resultType":"complete"`)
	require.NotContains(t, recorder.Body.String(), `"resourcesListChanged":true`)
	require.NotContains(t, recorder.Body.String(), "pulsar://context")
}

func TestHTTPPulsarSessionResolver(t *testing.T) {
	cfg := config.ExternalPulsar{
		ServiceURL: "pulsar://127.0.0.1:6650", WebServiceURL: "http://127.0.0.1:8080",
		Token: "global-token", AuthPlugin: "unused-plugin", AuthParams: "unused-params",
		TLSCertFile: "/must-not-read-client-cert", TLSKeyFile: "/must-not-read-client-key",
		TLSAllowInsecureConnection: true, TLSEnableHostnameVerification: true,
	}
	resolve := newHTTPPulsarSessionResolver(cfg)
	first, releaseFirst, err := resolve(t.Context(), "alice")
	require.NoError(t, err)
	t.Cleanup(releaseFirst)
	second, releaseSecond, err := resolve(t.Context(), "bob")
	require.NoError(t, err)
	t.Cleanup(releaseSecond)
	third, releaseThird, err := resolve(t.Context(), "alice")
	require.NoError(t, err)
	t.Cleanup(releaseThird)
	require.NotSame(t, first, second)
	require.NotSame(t, first, third, "even identical credentials have independent request lifetimes")
	require.Equal(t, pulsar.PulsarContext{
		ServiceURL: cfg.ServiceURL, WebServiceURL: cfg.WebServiceURL, Token: "alice",
		TLSAllowInsecureConnection: true, TLSEnableHostnameVerification: true,
	}, first.Ctx)
	require.Equal(t, "bob", second.Ctx.Token)
	require.NotNil(t, first.Client)
	releaseFirst()
	releaseFirst() // Idempotent release, including the deferred cleanup.
	require.Nil(t, first.Client)
	require.Nil(t, first.AdminClient)
	require.Empty(t, first.Ctx.Token)
	require.NotNil(t, second.Client, "releasing one request must not close another")
	require.NotNil(t, third.Client)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	s, release, err := resolve(ctx, "alice")
	require.ErrorIs(t, err, context.Canceled)
	require.Nil(t, s)
	require.Nil(t, release)
	s, release, err = resolve(t.Context(), "")
	require.Error(t, err)
	require.Nil(t, s)
	require.Nil(t, release)
	invalid := cfg
	invalid.ServiceURL = "://invalid"
	s, release, err = newHTTPPulsarSessionResolver(invalid)(t.Context(), "alice")
	require.Error(t, err)
	require.Nil(t, s)
	require.Nil(t, release)
}

func TestHTTPReleaseAfterToolFailureOrCancellation(t *testing.T) {
	for _, cancelRequest := range []bool{false, true} {
		t.Run(fmt.Sprintf("cancel=%t", cancelRequest), func(t *testing.T) {
			s := app.NewServer("http-test", "1", logrus.New())
			var released atomic.Int32
			started := make(chan struct{})
			s.MCPServer.AddTool(protocol.NewTool("fail"), func(ctx context.Context, _ protocol.CallToolRequest) (*protocol.CallToolResult, error) {
				close(started)
				if cancelRequest {
					<-ctx.Done()
				}
				if released.Load() != 0 {
					t.Error("the tool must retain its backend until it returns")
				}
				return nil, errors.New("tool failed")
			})
			handler, err := newHTTPHandler(s, nil, "/mcp", nil, func(context.Context, string) (*pulsar.Session, func(), error) {
				return &pulsar.Session{}, func() { released.Add(1) }, nil
			})
			require.NoError(t, err)
			req := modernHTTPRequest(t, "http://localhost", "tools/call")
			var body map[string]any
			require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
			body["params"].(map[string]any)["name"] = "fail"
			data, err := json.Marshal(body)
			require.NoError(t, err)
			req.Body = io.NopCloser(strings.NewReader(string(data)))
			req.ContentLength = int64(len(data))
			req.Header.Set("Mcp-Name", "fail")
			req.Header.Set("Authorization", "Bearer alice")
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			done := make(chan struct{})
			go func() {
				defer close(done)
				handler.ServeHTTP(httptest.NewRecorder(), req.WithContext(ctx))
			}()
			select {
			case <-started:
			case <-time.After(5 * time.Second):
				t.Fatal("tool did not start")
			}
			if cancelRequest {
				cancel()
			}
			select {
			case <-done:
			case <-time.After(5 * time.Second):
				t.Fatal("handler did not finish")
			}
			require.EqualValues(t, 1, released.Load())
		})
	}
}

func modernHTTPRequest(t *testing.T, base, method string) *http.Request {
	t.Helper()
	body := `{"jsonrpc":"2.0","id":1,"method":"` + method + `","params":{"_meta":{"io.modelcontextprotocol/protocolVersion":"2026-07-28","io.modelcontextprotocol/clientCapabilities":{}}}}`
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, base+"/mcp", strings.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	req.Header.Set("Mcp-Protocol-Version", "2026-07-28")
	req.Header.Set("Mcp-Method", method)
	return req
}
