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

package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	pulsaradminconfig "github.com/apache/pulsar-client-go/pulsaradmin/pkg/admin/config"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	pulsarsession "github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPulsarAddResourcesFeatureGate(t *testing.T) {
	t.Parallel()

	t.Run("disabled for unrelated features", func(t *testing.T) {
		s := server.NewMCPServer("test", "0.0.1", server.WithResourceCapabilities(true, true))
		PulsarAddResources(s, []string{string(FeatureKafkaAdmin)})

		resources := listPulsarTestResources(t, s)
		assert.Empty(t, resources)

		templates := listPulsarTestResourceTemplates(t, s)
		assert.Empty(t, templates)
	})

	t.Run("enabled for all-pulsar", func(t *testing.T) {
		s := server.NewMCPServer("test", "0.0.1", server.WithResourceCapabilities(true, true))
		PulsarAddResources(s, []string{string(FeatureAllPulsar)})

		resources := listPulsarTestResources(t, s)
		resourceURIs := make([]string, 0, len(resources))
		for _, resource := range resources {
			resourceURIs = append(resourceURIs, resource.URI)
			assert.Equal(t, pulsarResourceJSONMIMEType, resource.MIMEType)
		}
		assert.ElementsMatch(t, []string{pulsarResourceContextURI, pulsarResourceCatalogURI}, resourceURIs)

		templates := listPulsarTestResourceTemplates(t, s)
		templateURIs := make([]string, 0, len(templates))
		for _, template := range templates {
			templateURIs = append(templateURIs, template.URITemplate.Raw())
			assert.Equal(t, pulsarResourceJSONMIMEType, template.MIMEType)
		}
		assert.ElementsMatch(t, []string{pulsarNamespacesResourceTemplateURI, pulsarTopicsResourceTemplateURI}, templateURIs)
	})
}

func TestPulsarContextResourceReadRedactsSecrets(t *testing.T) {
	t.Parallel()

	s := server.NewMCPServer("test", "0.0.1", server.WithResourceCapabilities(true, true))
	PulsarAddResources(s, []string{string(FeatureAllPulsar)})

	ctx := WithPulsarSession(context.Background(), &pulsarsession.Session{
		Ctx: pulsarsession.PulsarContext{
			ServiceURL:                    "pulsar+ssl://broker.example:6651",
			WebServiceURL:                 "https://admin.example",
			Token:                         "secret-token",
			AuthPlugin:                    "org.example.SecretAuth",
			AuthParams:                    "secret-auth-params",
			TLSAllowInsecureConnection:    true,
			TLSEnableHostnameVerification: true,
			TLSTrustCertsFilePath:         "/tmp/ca.pem",
			TLSCertFile:                   "/tmp/client-cert.pem",
			TLSKeyFile:                    "/tmp/private-key.pem",
		},
	})

	content := readPulsarTestResource(t, ctx, s, pulsarResourceContextURI)
	assert.Equal(t, pulsarResourceContextURI, content.URI)
	assert.Equal(t, pulsarResourceJSONMIMEType, content.MIMEType)
	assert.Contains(t, content.Text, "https://admin.example")
	assert.Contains(t, content.Text, "pulsar+ssl://broker.example:6651")
	assert.Contains(t, content.Text, `"method": "token"`)
	assert.NotContains(t, content.Text, "secret-token")
	assert.NotContains(t, content.Text, "secret-auth-params")
	assert.NotContains(t, content.Text, "/tmp/private-key.pem")
}

func TestPulsarResourceReadMissingSession(t *testing.T) {
	t.Parallel()

	s := server.NewMCPServer("test", "0.0.1", server.WithResourceCapabilities(true, true))
	PulsarAddResources(s, []string{string(FeatureAllPulsar)})

	response := s.HandleMessage(context.Background(), []byte(`{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "resources/read",
		"params": {
			"uri": "pulsar://context"
		}
	}`))

	errResponse, ok := response.(mcp.JSONRPCError)
	require.True(t, ok, "expected JSONRPCError, got %T", response)
	assert.Equal(t, mcp.INTERNAL_ERROR, errResponse.Error.Code)
	assert.Contains(t, errResponse.Error.Message, "Pulsar session not found in context")
}

func TestPulsarNamespaceTemplateRead(t *testing.T) {
	t.Parallel()

	requestedPath := make(chan string, 1)
	adminServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath <- r.URL.Path
		if r.URL.Path != "/admin/v2/namespaces/public" {
			http.Error(w, "unexpected path", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`["public/default","public/functions"]`))
	}))
	defer adminServer.Close()

	cfg := &cmdutils.ClusterConfig{WebServiceURL: adminServer.URL}
	ctx := WithPulsarSession(context.Background(), &pulsarsession.Session{
		Ctx:             pulsarsession.PulsarContext{WebServiceURL: adminServer.URL},
		PulsarCtlConfig: cfg,
		AdminClient:     cfg.Client(pulsaradminconfig.V2),
	})

	s := server.NewMCPServer("test", "0.0.1", server.WithResourceCapabilities(true, true))
	PulsarAddResources(s, []string{string(FeatureAllPulsar)})

	content := readPulsarTestResource(t, ctx, s, "pulsar://admin/v2/tenants/public/namespaces")
	assert.Equal(t, "pulsar://admin/v2/tenants/public/namespaces", content.URI)
	assert.Equal(t, pulsarResourceJSONMIMEType, content.MIMEType)

	var payload pulsarNamespaceCollectionResource
	require.NoError(t, json.Unmarshal([]byte(content.Text), &payload))
	assert.Equal(t, "namespaces", payload.Kind)
	assert.Equal(t, "public", payload.Tenant)
	assert.ElementsMatch(t, []string{"public/default", "public/functions"}, payload.Namespaces)
	assert.Equal(t, 2, payload.Count)
	assert.Equal(t, "/admin/v2/namespaces/public", <-requestedPath)
}

func TestParsePulsarResourceURIRejectsMalformedOrUnsupportedURI(t *testing.T) {
	t.Parallel()

	tests := []string{
		"",
		"http://admin/v2/tenants/public/namespaces",
		"pulsar://user:pass@context",
		"pulsar://context/extra",
		"pulsar://resources?token=secret",
		"pulsar://admin/v3/tenants/public/namespaces",
		"pulsar://admin/v2/tenants/public",
		"pulsar://admin/v2/namespaces/public/topics",
	}

	for _, rawURI := range tests {
		t.Run(strings.ReplaceAll(rawURI, "/", "_"), func(t *testing.T) {
			_, err := parsePulsarResourceURI(rawURI)
			require.Error(t, err)
		})
	}
}

func listPulsarTestResources(t *testing.T, s *server.MCPServer) []mcp.Resource {
	t.Helper()

	response := s.HandleMessage(context.Background(), []byte(`{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "resources/list"
	}`))

	resp, ok := response.(mcp.JSONRPCResponse)
	require.True(t, ok, "expected JSONRPCResponse, got %T", response)
	result, ok := resp.Result.(mcp.ListResourcesResult)
	require.True(t, ok, "expected ListResourcesResult, got %T", resp.Result)
	return result.Resources
}

func listPulsarTestResourceTemplates(t *testing.T, s *server.MCPServer) []mcp.ResourceTemplate {
	t.Helper()

	response := s.HandleMessage(context.Background(), []byte(`{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "resources/templates/list"
	}`))

	resp, ok := response.(mcp.JSONRPCResponse)
	require.True(t, ok, "expected JSONRPCResponse, got %T", response)
	result, ok := resp.Result.(mcp.ListResourceTemplatesResult)
	require.True(t, ok, "expected ListResourceTemplatesResult, got %T", resp.Result)
	return result.ResourceTemplates
}

func readPulsarTestResource(t *testing.T, ctx context.Context, s *server.MCPServer, uri string) mcp.TextResourceContents {
	t.Helper()

	request := fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "resources/read",
		"params": {
			"uri": %q
		}
	}`, uri)
	response := s.HandleMessage(ctx, []byte(request))

	resp, ok := response.(mcp.JSONRPCResponse)
	require.True(t, ok, "expected JSONRPCResponse, got %T", response)
	result, ok := resp.Result.(mcp.ReadResourceResult)
	require.True(t, ok, "expected ReadResourceResult, got %T", resp.Result)
	require.Len(t, result.Contents, 1)
	content, ok := result.Contents[0].(mcp.TextResourceContents)
	require.True(t, ok, "expected TextResourceContents, got %T", result.Contents[0])
	return content
}
