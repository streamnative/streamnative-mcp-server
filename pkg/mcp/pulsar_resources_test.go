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
		assert.ElementsMatch(t, []string{
			pulsarResourceContextURI,
			pulsarResourceCatalogURI,
			pulsarClusterStatusResourceURI,
			pulsarClustersResourceURI,
			pulsarBrokerStatsSummaryResourceURI,
		}, resourceURIs)

		templates := listPulsarTestResourceTemplates(t, s)
		templateURIs := make([]string, 0, len(templates))
		for _, template := range templates {
			templateURIs = append(templateURIs, template.URITemplate.Raw())
			assert.Equal(t, pulsarResourceJSONMIMEType, template.MIMEType)
		}
		assert.ElementsMatch(t, []string{
			pulsarNamespacesResourceTemplateURI,
			pulsarTopicsResourceTemplateURI,
			pulsarClusterResourceTemplateURI,
			pulsarBrokersResourceTemplateURI,
			pulsarFailureDomainsTemplateURI,
			pulsarFailureDomainTemplateURI,
			pulsarNSIsolationPoliciesTemplateURI,
			pulsarNSIsolationPolicyTemplateURI,
		}, templateURIs)
	})

	t.Run("cluster feature registers cluster resource family only", func(t *testing.T) {
		s := server.NewMCPServer("test", "0.0.1", server.WithResourceCapabilities(true, true))
		PulsarAddResources(s, []string{string(FeaturePulsarAdminClusters)})

		resources := listPulsarTestResources(t, s)
		resourceURIs := make([]string, 0, len(resources))
		for _, resource := range resources {
			resourceURIs = append(resourceURIs, resource.URI)
		}
		assert.ElementsMatch(t, []string{
			pulsarResourceContextURI,
			pulsarResourceCatalogURI,
			pulsarClustersResourceURI,
		}, resourceURIs)

		templates := listPulsarTestResourceTemplates(t, s)
		templateURIs := make([]string, 0, len(templates))
		for _, template := range templates {
			templateURIs = append(templateURIs, template.URITemplate.Raw())
		}
		assert.ElementsMatch(t, []string{
			pulsarClusterResourceTemplateURI,
			pulsarFailureDomainsTemplateURI,
			pulsarFailureDomainTemplateURI,
		}, templateURIs)
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

func TestPulsarClusterResourceFamilyRead(t *testing.T) {
	t.Parallel()

	adminServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
			return
		}

		switch r.URL.Path {
		case "/status.html":
			w.Header().Set("Content-Type", "text/plain")
			_, _ = w.Write([]byte("OK\n"))
		case "/admin/v2/clusters":
			writePulsarTestJSON(w, `["use"]`)
		case "/admin/v2/clusters/use":
			writePulsarTestJSON(w, `{
				"serviceUrl": "http://admin.example",
				"serviceUrlTls": "https://admin.example",
				"brokerServiceUrl": "pulsar://broker.example:6650",
				"brokerServiceUrlTls": "pulsar+ssl://broker.example:6651",
				"peerClusterNames": ["peer-a"],
				"authenticationPlugin": "org.example.SecretAuth",
				"authenticationParameters": "secret-auth-params",
				"brokerClientTrustCertsFilePath": "/tmp/ca.pem",
				"brokerClientTlsEnabled": true
			}`)
		case "/admin/v2/brokers/use":
			writePulsarTestJSON(w, `["broker-a:8080","broker-b:8080"]`)
		case "/admin/v2/broker-stats/metrics":
			writePulsarTestJSON(w, `[{
				"metrics": {
					"brk_msg_rate_in": 12.5,
					"brk_msg_rate_out": 7.5
				},
				"dimensions": {
					"cluster": "use",
					"broker": "broker-a"
				}
			}]`)
		case "/admin/v2/broker-stats/load-report":
			writePulsarTestJSON(w, `{
				"webServiceUrl": "http://broker-a:8080",
				"pulsarServiceUrl": "pulsar://broker-a:6650",
				"persistentTopicsEnabled": true,
				"nonPersistentTopicsEnabled": true,
				"cpu": {"usage": 1, "limit": 4},
				"memory": {"usage": 256, "limit": 1024},
				"msgRateIn": 10,
				"msgRateOut": 8,
				"numTopics": 2,
				"numBundles": 1,
				"numConsumers": 3,
				"numProducers": 4,
				"bundles": ["bundle-a", "bundle-b"],
				"lastBundleGains": ["bundle-a"],
				"lastBundleLosses": ["bundle-z"],
				"brokerVersionString": "4.1.0",
				"loadReportType": "LocalBrokerData",
				"protocols": {"kafka": "9092"}
			}`)
		case "/admin/v2/clusters/use/failureDomains":
			writePulsarTestJSON(w, `{
				"zone-a": {"brokers": ["broker-a:8080"]}
			}`)
		case "/admin/v2/clusters/use/failureDomains/zone-a":
			writePulsarTestJSON(w, `{"brokers": ["broker-a:8080"]}`)
		case "/admin/v2/clusters/use/namespaceIsolationPolicies":
			writePulsarTestJSON(w, `{
				"policy-a": {
					"namespaces": ["public/.*"],
					"primary": ["broker-a.*"],
					"secondary": ["broker-b.*"],
					"auto_failover_policy": {
						"policy_type": "min_available",
						"parameters": {
							"min_limit": "1",
							"usage_threshold": "80"
						}
					}
				}
			}`)
		case "/admin/v2/clusters/use/namespaceIsolationPolicies/policy-a":
			writePulsarTestJSON(w, `{
				"namespaces": ["public/.*"],
				"primary": ["broker-a.*"],
				"secondary": ["broker-b.*"],
				"auto_failover_policy": {
					"policy_type": "min_available",
					"parameters": {
						"min_limit": "1",
						"usage_threshold": "80"
					}
				}
			}`)
		default:
			http.Error(w, "unexpected path "+r.URL.Path, http.StatusNotFound)
		}
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

	statusContent := readPulsarTestResource(t, ctx, s, pulsarClusterStatusResourceURI)
	var statusPayload pulsarClusterStatusResource
	require.NoError(t, json.Unmarshal([]byte(statusContent.Text), &statusPayload))
	assert.Equal(t, "OK", statusPayload.Status)

	clustersContent := readPulsarTestResource(t, ctx, s, pulsarClustersResourceURI)
	var clustersPayload pulsarClusterCollectionResource
	require.NoError(t, json.Unmarshal([]byte(clustersContent.Text), &clustersPayload))
	assert.Equal(t, []string{"use"}, clustersPayload.Clusters)
	assert.Equal(t, 1, clustersPayload.Count)

	clusterContent := readPulsarTestResource(t, ctx, s, "pulsar://admin/v2/clusters/use")
	assert.NotContains(t, clusterContent.Text, "secret-auth-params")
	assert.NotContains(t, clusterContent.Text, "/tmp/ca.pem")
	var clusterPayload pulsarClusterResource
	require.NoError(t, json.Unmarshal([]byte(clusterContent.Text), &clusterPayload))
	assert.Equal(t, "use", clusterPayload.Data.Name)
	assert.True(t, clusterPayload.Data.AuthenticationParametersConfigured)
	assert.True(t, clusterPayload.Data.BrokerClientTrustCertsFileConfigured)
	assert.True(t, clusterPayload.Data.BrokerClientTLSEnabled)

	brokersContent := readPulsarTestResource(t, ctx, s, "pulsar://admin/v2/brokers/use")
	var brokersPayload pulsarBrokerCollectionResource
	require.NoError(t, json.Unmarshal([]byte(brokersContent.Text), &brokersPayload))
	assert.ElementsMatch(t, []string{"broker-a:8080", "broker-b:8080"}, brokersPayload.Brokers)
	assert.Equal(t, 2, brokersPayload.Count)

	statsContent := readPulsarTestResource(t, ctx, s, pulsarBrokerStatsSummaryResourceURI)
	var statsPayload pulsarBrokerStatsSummaryResource
	require.NoError(t, json.Unmarshal([]byte(statsContent.Text), &statsPayload))
	assert.Equal(t, 1, statsPayload.MonitoringMetrics.Count)
	assert.ElementsMatch(t, []string{"brk_msg_rate_in", "brk_msg_rate_out"}, statsPayload.MonitoringMetrics.MetricNames)
	assert.ElementsMatch(t, []string{"broker", "cluster"}, statsPayload.MonitoringMetrics.DimensionKeys)
	assert.True(t, statsPayload.LoadReport.Available)
	assert.Equal(t, 2, statsPayload.LoadReport.BundleCount)
	assert.Equal(t, 1, statsPayload.LoadReport.ProtocolCount)

	failureDomainsContent := readPulsarTestResource(t, ctx, s, "pulsar://admin/v2/clusters/use/failureDomains")
	var failureDomainsPayload pulsarFailureDomainCollectionResource
	require.NoError(t, json.Unmarshal([]byte(failureDomainsContent.Text), &failureDomainsPayload))
	require.Len(t, failureDomainsPayload.FailureDomains, 1)
	assert.Equal(t, "zone-a", failureDomainsPayload.FailureDomains[0].Name)
	assert.Equal(t, []string{"broker-a:8080"}, failureDomainsPayload.FailureDomains[0].Brokers)

	failureDomainContent := readPulsarTestResource(t, ctx, s, "pulsar://admin/v2/clusters/use/failureDomains/zone-a")
	var failureDomainPayload pulsarFailureDomainResource
	require.NoError(t, json.Unmarshal([]byte(failureDomainContent.Text), &failureDomainPayload))
	assert.Equal(t, "zone-a", failureDomainPayload.FailureDomain.Name)

	policiesContent := readPulsarTestResource(t, ctx, s, "pulsar://admin/v2/clusters/use/namespaceIsolationPolicies")
	var policiesPayload pulsarNamespaceIsolationPolicyCollectionResource
	require.NoError(t, json.Unmarshal([]byte(policiesContent.Text), &policiesPayload))
	require.Len(t, policiesPayload.Policies, 1)
	assert.Equal(t, "policy-a", policiesPayload.Policies[0].Name)
	assert.Equal(t, 1, policiesPayload.Policies[0].NamespacesCount)
	assert.Equal(t, "min_available", policiesPayload.Policies[0].AutoFailoverPolicyType)

	policyContent := readPulsarTestResource(t, ctx, s, "pulsar://admin/v2/clusters/use/namespaceIsolationPolicies/policy-a")
	var policyPayload pulsarNamespaceIsolationPolicyResource
	require.NoError(t, json.Unmarshal([]byte(policyContent.Text), &policyPayload))
	assert.Equal(t, "policy-a", policyPayload.Policy.Name)
	assert.Equal(t, []string{"public/.*"}, policyPayload.Policy.Data.Namespaces)
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
		"pulsar://admin/v2/broker-stats/topics",
		"pulsar://admin/v2/brokers",
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

func writePulsarTestJSON(w http.ResponseWriter, body string) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(body))
}
