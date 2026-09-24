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
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	app "github.com/streamnative/streamnative-mcp-server/pkg/mcp"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

type unusedResolver struct{}

func (unusedResolver) ResolveCluster(context.Context, string, string) (config.ResolvedRuntime, error) {
	return config.ResolvedRuntime{}, fmt.Errorf("not expected during bootstrap")
}

func TestInjectedCloudBootstrapDoesNotRequireKeyFileOrGrantStore(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer injected-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	}))
	defer api.Close()
	o := config.NewConfigOptions()
	o.ConfigDir, o.BackendOverride = t.TempDir(), "must-not-open-keyring"
	o.Server, o.LogLocation, o.Organization = api.URL, api.URL, "org"
	o.CloudProvider = &config.CloudProvider{
		Identity: "effective-sa", Resolver: unusedResolver{},
		TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "injected-token", TokenType: "Bearer"}),
	}
	options := NewMcpServerOptions(o)
	require.NoError(t, options.Complete())
	require.Nil(t, o.Store)
	require.Empty(t, o.KeyFile)
	srv, err := newMcpServer(context.Background(), options, logrus.New())
	require.NoError(t, err)
	defer srv.Close()
	client, err := srv.SNCloudSession.GetAPIClient()
	require.NoError(t, err)
	_, response, err := client.CloudStreamnativeIoV1alpha1Api.ListCloudStreamnativeIoV1alpha1NamespacedPulsarCluster(context.Background(), "org").Execute()
	if response != nil {
		defer func() { _ = response.Body.Close() }()
	}
	require.NoError(t, err)
	logs, err := srv.SNCloudSession.GetLogClient()
	require.NoError(t, err)
	response, err = logs.Get(api.URL + "/logs")
	require.NoError(t, err)
	defer func() { _ = response.Body.Close() }()
	require.Equal(t, http.StatusOK, response.StatusCode)
	ctx := context.WithValue(context.Background(), common.OptionsKey, o)
	result, err := srv.MCPServer.GetTool("sncloud_context_whoami").Handler(ctx, mcpgo.CallToolRequest{})
	require.NoError(t, err)
	require.False(t, result.IsError)
	require.Contains(t, fmt.Sprint(result.Content), "effective-sa")
}

func TestStandaloneKeyFileStillResolvesLegacyRuntime(t *testing.T) {
	var mu sync.Mutex
	audiences := make(map[string]int)
	var api *httptest.Server
	api = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/.well-known/openid-configuration"):
			_, _ = fmt.Fprintf(w, `{"token_endpoint":%q}`, api.URL+"/token")
		case r.URL.Path == "/token":
			_ = r.ParseForm()
			mu.Lock()
			audiences[r.Form.Get("audience")]++
			mu.Unlock()
			_, _ = fmt.Fprint(w, `{"access_token":"legacy-token","token_type":"Bearer","expires_in":3600}`)
		case strings.HasSuffix(r.URL.Path, "/pulsarclusters/cluster"):
			_, _ = fmt.Fprint(w, `{"metadata":{"name":"cluster","namespace":"org","uid":"uid"},
			"spec":{"instanceName":"instance","serviceEndpoints":[{"type":"service","dnsName":"localhost"}]},
			"status":{"broker":{"readyReplicas":1},"conditions":[{"type":"Ready","status":"True"}]}}`)
		case strings.HasSuffix(r.URL.Path, "/pulsarinstances/instance"):
			_, _ = fmt.Fprintf(w, `{"metadata":{"name":"instance","namespace":"org"},
			"status":{"auth":{"type":"oauth2","oauth2":{"issuerURL":%q,"audience":"legacy-runtime"}}}}`, api.URL+"/")
		default:
			http.NotFound(w, r)
		}
	}))
	defer api.Close()
	o := config.NewConfigOptions()
	o.ConfigDir, o.Server, o.ProxyLocation = t.TempDir(), api.URL, api.URL
	o.Organization, o.Audience, o.IssuerEndpoint = "org", "control", api.URL+"/"
	o.KeyFile = `data://{"type":"sn_service_account","client_id":"legacy-sa","client_secret":"fake-secret"}`
	options := NewMcpServerOptions(o)
	require.NoError(t, options.Complete())
	srv, err := newMcpServer(context.Background(), options, logrus.New())
	require.NoError(t, err)
	defer srv.Close()
	ctx := app.WithSNCloudSession(context.Background(), srv.SNCloudSession)
	require.NoError(t, app.SetContext(ctx, o, "instance", "cluster"))
	require.NoError(t, app.SetContext(ctx, o, "instance", "cluster"))
	mu.Lock()
	defer mu.Unlock()
	require.Equal(t, map[string]int{"control": 1, "legacy-runtime": 1}, audiences)
}

func TestHTTPRejectsInjectedCloudBeforeAuthentication(t *testing.T) {
	o := config.NewConfigOptions()
	o.CloudProvider = &config.CloudProvider{}
	o.UseExternalKafka = true
	cmd := NewCmdMcpHTTPServer(NewMcpServerOptions(o))
	err := cmd.ExecuteContext(context.Background())
	require.ErrorContains(t, err, "http transport")
	require.Nil(t, o.Store)
}

func TestInjectedInitialClusterIsResolvedBeforeServing(t *testing.T) {
	for _, transport := range []struct {
		name string
		run  func(*ServerOptions) error
	}{{"stdio", runStdioServer}, {"sse", runSseServer}} {
		t.Run(transport.name, func(t *testing.T) {
			o := config.NewConfigOptions()
			o.Server, o.Organization = "http://127.0.0.1:1", "org"
			o.PulsarInstance, o.PulsarCluster = "instance", "cluster"
			o.CloudProvider = &config.CloudProvider{
				Identity: "subject", Resolver: unusedResolver{},
				TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "opaque"}),
			}
			opts := NewMcpServerOptions(o)
			require.NoError(t, opts.Complete())
			require.ErrorContains(t, transport.run(opts), "not expected during bootstrap")
		})
	}
}
