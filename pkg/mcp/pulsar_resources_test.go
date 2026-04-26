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
			pulsarTenantsResourceURI,
			pulsarDefaultResourceQuotaURI,
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
			pulsarTenantResourceTemplateURI,
			pulsarNamespacesResourceTemplateURI,
			pulsarNamespaceResourceTemplateURI,
			pulsarTopicsResourceTemplateURI,
			pulsarTopicMetadataTemplateURI,
			pulsarTopicStatsTemplateURI,
			pulsarTopicPartitionMetadataURI,
			pulsarTopicPolicyTemplateURI,
			pulsarTopicSchemaTemplateURI,
			pulsarTopicSchemaVersionTemplateURI,
			pulsarSubscriptionsTemplateURI,
			pulsarSubscriptionStatsTemplateURI,
			pulsarSubscriptionBacklogTemplateURI,
			pulsarSubscriptionCursorTemplateURI,
			pulsarResourceQuotaTemplateURI,
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

	t.Run("tenant namespace feature gates register selected resource families", func(t *testing.T) {
		s := server.NewMCPServer("test", "0.0.1", server.WithResourceCapabilities(true, true))
		PulsarAddResources(s, []string{
			string(FeaturePulsarAdminTenants),
			string(FeaturePulsarAdminNamespacePolicy),
			string(FeaturePulsarAdminResourceQuotas),
		})

		resources := listPulsarTestResources(t, s)
		resourceURIs := make([]string, 0, len(resources))
		for _, resource := range resources {
			resourceURIs = append(resourceURIs, resource.URI)
		}
		assert.ElementsMatch(t, []string{
			pulsarResourceContextURI,
			pulsarResourceCatalogURI,
			pulsarTenantsResourceURI,
			pulsarDefaultResourceQuotaURI,
		}, resourceURIs)

		templates := listPulsarTestResourceTemplates(t, s)
		templateURIs := make([]string, 0, len(templates))
		for _, template := range templates {
			templateURIs = append(templateURIs, template.URITemplate.Raw())
		}
		assert.ElementsMatch(t, []string{
			pulsarTenantResourceTemplateURI,
			pulsarNamespaceResourceTemplateURI,
			pulsarResourceQuotaTemplateURI,
		}, templateURIs)
	})

	t.Run("topic schema feature gates register selected resource families", func(t *testing.T) {
		s := server.NewMCPServer("test", "0.0.1", server.WithResourceCapabilities(true, true))
		PulsarAddResources(s, []string{
			string(FeaturePulsarAdminTopics),
			string(FeaturePulsarAdminTopicPolicy),
			string(FeaturePulsarAdminSchemas),
		})

		resources := listPulsarTestResources(t, s)
		resourceURIs := make([]string, 0, len(resources))
		for _, resource := range resources {
			resourceURIs = append(resourceURIs, resource.URI)
		}
		assert.ElementsMatch(t, []string{
			pulsarResourceContextURI,
			pulsarResourceCatalogURI,
		}, resourceURIs)

		templates := listPulsarTestResourceTemplates(t, s)
		templateURIs := make([]string, 0, len(templates))
		for _, template := range templates {
			templateURIs = append(templateURIs, template.URITemplate.Raw())
		}
		assert.ElementsMatch(t, []string{
			pulsarTopicsResourceTemplateURI,
			pulsarTopicMetadataTemplateURI,
			pulsarTopicStatsTemplateURI,
			pulsarTopicPartitionMetadataURI,
			pulsarTopicPolicyTemplateURI,
			pulsarTopicSchemaTemplateURI,
			pulsarTopicSchemaVersionTemplateURI,
		}, templateURIs)
	})

	t.Run("subscription feature registers subscription resource family only", func(t *testing.T) {
		s := server.NewMCPServer("test", "0.0.1", server.WithResourceCapabilities(true, true))
		PulsarAddResources(s, []string{string(FeaturePulsarAdminSubscriptions)})

		resources := listPulsarTestResources(t, s)
		resourceURIs := make([]string, 0, len(resources))
		for _, resource := range resources {
			resourceURIs = append(resourceURIs, resource.URI)
		}
		assert.ElementsMatch(t, []string{
			pulsarResourceContextURI,
			pulsarResourceCatalogURI,
		}, resourceURIs)

		templates := listPulsarTestResourceTemplates(t, s)
		templateURIs := make([]string, 0, len(templates))
		for _, template := range templates {
			templateURIs = append(templateURIs, template.URITemplate.Raw())
		}
		assert.ElementsMatch(t, []string{
			pulsarSubscriptionsTemplateURI,
			pulsarSubscriptionStatsTemplateURI,
			pulsarSubscriptionBacklogTemplateURI,
			pulsarSubscriptionCursorTemplateURI,
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

func TestPulsarTenantNamespaceResourceFamilyRead(t *testing.T) {
	t.Parallel()

	adminServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
			return
		}

		switch r.URL.Path {
		case "/admin/v2/tenants":
			writePulsarTestJSON(w, `["public","system"]`)
		case "/admin/v2/tenants/public":
			writePulsarTestJSON(w, `{
				"adminRoles": ["tenant-admin"],
				"allowedClusters": ["use"]
			}`)
		case "/admin/v2/namespaces/public":
			writePulsarTestJSON(w, `["public/default","public/functions"]`)
		case "/admin/v2/namespaces/public/default":
			writePulsarTestJSON(w, `{
				"replication_clusters": ["use"],
				"message_ttl_in_seconds": 60,
				"schema_validation_enforced": true
			}`)
		case "/admin/v2/namespaces/public/default/topics":
			writePulsarTestJSON(w, `["persistent://public/default/events"]`)
		case "/admin/v2/resource-quotas":
			writePulsarTestJSON(w, `{
				"msgRateIn": 100,
				"msgRateOut": 200,
				"bandwidthIn": 1024,
				"bandwidthOut": 2048,
				"memory": 4096,
				"dynamic": true
			}`)
		case "/admin/v2/resource-quotas/public/default/0x00000000_0xffffffff":
			writePulsarTestJSON(w, `{
				"msgRateIn": 10,
				"msgRateOut": 20,
				"bandwidthIn": 128,
				"bandwidthOut": 256,
				"memory": 512,
				"dynamic": false
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

	tenantsContent := readPulsarTestResource(t, ctx, s, pulsarTenantsResourceURI)
	var tenantsPayload pulsarTenantCollectionResource
	require.NoError(t, json.Unmarshal([]byte(tenantsContent.Text), &tenantsPayload))
	assert.ElementsMatch(t, []string{"public", "system"}, tenantsPayload.Tenants)
	assert.Equal(t, 2, tenantsPayload.Count)

	tenantContent := readPulsarTestResource(t, ctx, s, "pulsar://admin/v2/tenants/public")
	var tenantPayload pulsarTenantResource
	require.NoError(t, json.Unmarshal([]byte(tenantContent.Text), &tenantPayload))
	assert.Equal(t, "public", tenantPayload.Tenant)
	assert.Equal(t, []string{"tenant-admin"}, tenantPayload.Data.AdminRoles)
	assert.Equal(t, []string{"use"}, tenantPayload.Data.AllowedClusters)

	namespacesContent := readPulsarTestResource(t, ctx, s, "pulsar://admin/v2/tenants/public/namespaces")
	var namespacesPayload pulsarNamespaceCollectionResource
	require.NoError(t, json.Unmarshal([]byte(namespacesContent.Text), &namespacesPayload))
	assert.ElementsMatch(t, []string{"public/default", "public/functions"}, namespacesPayload.Namespaces)
	assert.Equal(t, 2, namespacesPayload.Count)

	namespaceContent := readPulsarTestResource(t, ctx, s, "pulsar://admin/v2/namespaces/public/default")
	var namespacePayload pulsarNamespaceResource
	require.NoError(t, json.Unmarshal([]byte(namespaceContent.Text), &namespacePayload))
	require.NotNil(t, namespacePayload.Policies)
	assert.Equal(t, []string{"use"}, namespacePayload.Policies.ReplicationClusters)
	require.NotNil(t, namespacePayload.Policies.MessageTTLInSeconds)
	assert.Equal(t, 60, *namespacePayload.Policies.MessageTTLInSeconds)
	assert.True(t, namespacePayload.Policies.SchemaValidationEnforced)

	topicsContent := readPulsarTestResource(t, ctx, s, "pulsar://admin/v2/namespaces/public/default/topics")
	var topicsPayload pulsarTopicCollectionResource
	require.NoError(t, json.Unmarshal([]byte(topicsContent.Text), &topicsPayload))
	assert.Equal(t, []string{"persistent://public/default/events"}, topicsPayload.Topics)
	assert.Equal(t, 1, topicsPayload.Count)

	defaultQuotaContent := readPulsarTestResource(t, ctx, s, pulsarDefaultResourceQuotaURI)
	var defaultQuotaPayload pulsarResourceQuotaResource
	require.NoError(t, json.Unmarshal([]byte(defaultQuotaContent.Text), &defaultQuotaPayload))
	assert.Equal(t, "default", defaultQuotaPayload.Scope)
	require.NotNil(t, defaultQuotaPayload.Quota)
	assert.Equal(t, 100.0, defaultQuotaPayload.Quota.MsgRateIn)
	assert.True(t, defaultQuotaPayload.Quota.Dynamic)

	quotaContent := readPulsarTestResource(
		t,
		ctx,
		s,
		"pulsar://admin/v2/resource-quotas/public/default/0x00000000_0xffffffff",
	)
	var quotaPayload pulsarResourceQuotaResource
	require.NoError(t, json.Unmarshal([]byte(quotaContent.Text), &quotaPayload))
	assert.Equal(t, "namespaceBundle", quotaPayload.Scope)
	assert.Equal(t, "public", quotaPayload.Tenant)
	assert.Equal(t, "default", quotaPayload.Namespace)
	assert.Equal(t, "0x00000000_0xffffffff", quotaPayload.Bundle)
	require.NotNil(t, quotaPayload.Quota)
	assert.Equal(t, 10.0, quotaPayload.Quota.MsgRateIn)
	assert.False(t, quotaPayload.Quota.Dynamic)
}

func TestPulsarTopicSchemaResourceFamilyRead(t *testing.T) {
	t.Parallel()

	adminServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
			return
		}

		switch r.URL.Path {
		case "/admin/v2/persistent/public/default/events/properties":
			writePulsarTestJSON(w, `{
				"owner": "team-a",
				"token": "secret-token"
			}`)
		case "/admin/v2/persistent/public/default/events/partitions":
			writePulsarTestJSON(w, `{"partitions": 2}`)
		case "/admin/v2/persistent/public/default/events/partitioned-stats":
			if r.URL.Query().Get("perPartition") != "false" ||
				r.URL.Query().Get("excludePublishers") != "true" ||
				r.URL.Query().Get("excludeConsumers") != "true" {
				http.Error(w, "unexpected stats query "+r.URL.RawQuery, http.StatusBadRequest)
				return
			}
			writePulsarTestJSON(w, `{
				"msgRateIn": 12.5,
				"msgRateOut": 4.5,
				"msgThroughputIn": 1024,
				"msgThroughputOut": 512,
				"averageMsgSize": 128,
				"storageSize": 4096,
				"publishers": [{"producerName": "hidden-producer"}],
				"subscriptions": {
					"sub-a": {"msgBacklog": 5}
				},
				"replication": {
					"use": {"connected": true}
				},
				"deduplicationStatus": "Enabled",
				"metadata": {"partitions": 2},
				"partitions": {
					"persistent://public/default/events-partition-0": {
						"msgRateIn": 1
					}
				},
				"topicCreationTimeStamp": 11,
				"lastPublishTimestamp": 22
			}`)
		case "/admin/v2/persistent/public/default/events/retention":
			if r.URL.Query().Get("applied") != "false" {
				http.Error(w, "unexpected retention query "+r.URL.RawQuery, http.StatusBadRequest)
				return
			}
			writePulsarTestJSON(w, `{
				"retentionTimeInMinutes": 60,
				"retentionSizeInMB": 1024
			}`)
		case "/admin/v2/schemas/public/default/events/schema":
			writePulsarTestJSON(w, `{
				"version": 8,
				"type": "JSON",
				"timestamp": 12345,
				"data": "{\"type\":\"record\",\"name\":\"Event\",\"fields\":[]}",
				"properties": {
					"owner": "team-a",
					"password": "secret-password"
				}
			}`)
		case "/admin/v2/schemas/public/default/events/schema/7":
			writePulsarTestJSON(w, `{
				"version": 7,
				"type": "STRING",
				"timestamp": 12340,
				"data": "string-schema",
				"properties": {
					"owner": "team-b"
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

	metadataContent := readPulsarTestResource(t, ctx, s, "pulsar://admin/v2/persistent/public/default/events/metadata")
	assert.NotContains(t, metadataContent.Text, "secret-token")
	var metadataPayload pulsarTopicMetadataResource
	require.NoError(t, json.Unmarshal([]byte(metadataContent.Text), &metadataPayload))
	assert.Equal(t, "topicMetadata", metadataPayload.Kind)
	assert.Equal(t, "persistent://public/default/events", metadataPayload.Topic)
	assert.Equal(t, "persistent", metadataPayload.Domain)
	assert.Equal(t, "events", metadataPayload.Name)
	assert.Equal(t, 2, metadataPayload.PropertiesCount)
	assert.Equal(t, "team-a", metadataPayload.Properties["owner"])
	assert.Equal(t, pulsarResourceRedactedValue, metadataPayload.Properties["token"])

	partitionsContent := readPulsarTestResource(t, ctx, s, "pulsar://admin/v2/persistent/public/default/events/partitions")
	var partitionsPayload pulsarTopicPartitionMetadataResource
	require.NoError(t, json.Unmarshal([]byte(partitionsContent.Text), &partitionsPayload))
	assert.True(t, partitionsPayload.Partitioned)
	assert.Equal(t, 2, partitionsPayload.Metadata.Partitions)

	statsContent := readPulsarTestResource(t, ctx, s, "pulsar://admin/v2/persistent/public/default/events/stats")
	assert.NotContains(t, statsContent.Text, "hidden-producer")
	var statsPayload pulsarTopicStatsResource
	require.NoError(t, json.Unmarshal([]byte(statsContent.Text), &statsPayload))
	assert.True(t, statsPayload.Partitioned)
	assert.Equal(t, 2, statsPayload.PartitionCount)
	assert.Equal(t, 1, statsPayload.PartitionStatsCount)
	assert.Equal(t, 12.5, statsPayload.Stats.MsgRateIn)
	assert.Equal(t, 1, statsPayload.Stats.PublisherCount)
	assert.Equal(t, 1, statsPayload.Stats.SubscriptionCount)

	policyContent := readPulsarTestResource(
		t,
		ctx,
		s,
		"pulsar://admin/v2/persistent/public/default/events/policies/retention",
	)
	var policyPayload pulsarTopicPolicyResource
	require.NoError(t, json.Unmarshal([]byte(policyContent.Text), &policyPayload))
	assert.Equal(t, "retention", policyPayload.Policy)
	policyValue, ok := policyPayload.Value.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, float64(60), policyValue["retentionTimeInMinutes"])

	latestSchemaContent := readPulsarTestResource(t, ctx, s, "pulsar://admin/v2/persistent/public/default/events/schema")
	assert.NotContains(t, latestSchemaContent.Text, "secret-password")
	var latestSchemaPayload pulsarTopicSchemaResource
	require.NoError(t, json.Unmarshal([]byte(latestSchemaContent.Text), &latestSchemaPayload))
	assert.Equal(t, int64(8), latestSchemaPayload.Version)
	assert.Equal(t, "JSON", latestSchemaPayload.Schema.Type)
	assert.Contains(t, latestSchemaPayload.Schema.Schema, `"name":"Event"`)
	assert.Equal(t, pulsarResourceRedactedValue, latestSchemaPayload.Schema.Properties["password"])

	versionSchemaContent := readPulsarTestResource(t, ctx, s, "pulsar://admin/v2/persistent/public/default/events/schema/7")
	var versionSchemaPayload pulsarTopicSchemaResource
	require.NoError(t, json.Unmarshal([]byte(versionSchemaContent.Text), &versionSchemaPayload))
	assert.Equal(t, int64(7), versionSchemaPayload.Version)
	assert.Equal(t, "STRING", versionSchemaPayload.Schema.Type)
	assert.Equal(t, "string-schema", versionSchemaPayload.Schema.Schema)
}

func TestPulsarSubscriptionResourceFamilyRead(t *testing.T) {
	t.Parallel()

	adminServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
			return
		}

		switch r.URL.Path {
		case "/admin/v2/persistent/public/default/events/subscriptions":
			writePulsarTestJSON(w, `["sub-a","sub-b"]`)
		case "/admin/v2/persistent/public/default/events/partitions":
			writePulsarTestJSON(w, `{"partitions": 0}`)
		case "/admin/v2/persistent/public/default/events/stats":
			if r.URL.Query().Get("excludePublishers") != "true" ||
				r.URL.Query().Get("excludeConsumers") != "true" {
				http.Error(w, "unexpected stats query "+r.URL.RawQuery, http.StatusBadRequest)
				return
			}
			if r.URL.Query().Get("subscriptionBacklogSize") == "true" &&
				r.URL.Query().Get("getEarliestTimeInBacklog") != "true" {
				http.Error(w, "unexpected backlog query "+r.URL.RawQuery, http.StatusBadRequest)
				return
			}
			writePulsarTestJSON(w, `{
				"subscriptions": {
					"sub-a": {
						"type": "Shared",
						"isDurable": true,
						"isReplicated": false,
						"blockedSubscriptionOnUnackedMsgs": true,
						"msgRateOut": 7.5,
						"msgThroughputOut": 2048,
						"msgRateRedeliver": 1.5,
						"msgRateExpired": 0.5,
						"msgBacklog": 42,
						"msgBacklogNoDelayed": 40,
						"msgDelayed": 2,
						"unackedMessages": 3,
						"bytesOutCounter": 1000,
						"msgOutCounter": 10,
						"messageAckRate": 6.5,
						"chunkedMessageRate": 0.25,
						"backlogSize": 4096,
						"earliestMsgPublishTimeInBacklog": 12345,
						"totalMsgExpired": 4,
						"lastExpireTimestamp": 55,
						"lastConsumedFlowTimestamp": 66,
						"lastConsumedTimestamp": 77,
						"lastAckedTimestamp": 88,
						"lastMarkDeleteAdvancedTimestamp": 99,
						"allowOutOfOrderDelivery": true,
						"nonContiguousDeletedMessagesRanges": 2,
						"nonContiguousDeletedMessagesRangesSerializedSize": 16,
						"delayedMessageIndexSizeInBytes": 128,
						"subscriptionProperties": {
							"owner": "team-a",
							"token": "secret-token"
						},
						"filterProcessedMsgCount": 11,
						"filterAcceptedMsgCount": 12,
						"filterRejectedMsgCount": 13,
						"filterRescheduledMsgCount": 14,
						"consumers": [{
							"consumerName": "hidden-consumer",
							"metadata": {"password": "secret-password"}
						}]
					}
				}
			}`)
		case "/admin/v2/persistent/public/default/events/internalStats":
			writePulsarTestJSON(w, `{
				"cursors": {
					"sub-a": {
						"markDeletePosition": "1:2",
						"readPosition": "1:3",
						"waitingReadOp": false,
						"pendingReadOps": 1,
						"messagesConsumedCounter": 10,
						"cursorLedger": 9,
						"cursorLedgerLastEntry": 8,
						"lastLedgerWitchTimestamp": "2026-04-26T00:00:00Z",
						"state": "Open",
						"numberOfEntriesSinceFirstNotAckedMessage": 7,
						"totalNonContiguousDeletedMessagesRange": 2,
						"individuallyDeletedMessages": "hidden-large-range",
						"properties": {
							"safe": 1,
							"token": 2
						}
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

	subscriptionsContent := readPulsarTestResource(t, ctx, s, "pulsar://admin/v2/persistent/public/default/events/subscriptions")
	var subscriptionsPayload pulsarSubscriptionCollectionResource
	require.NoError(t, json.Unmarshal([]byte(subscriptionsContent.Text), &subscriptionsPayload))
	assert.Equal(t, "subscriptions", subscriptionsPayload.Kind)
	assert.Equal(t, "persistent://public/default/events", subscriptionsPayload.Topic)
	assert.ElementsMatch(t, []string{"sub-a", "sub-b"}, subscriptionsPayload.Subscriptions)
	assert.Equal(t, 2, subscriptionsPayload.Count)

	statsContent := readPulsarTestResource(
		t,
		ctx,
		s,
		"pulsar://admin/v2/persistent/public/default/events/subscriptions/sub-a/stats",
	)
	assert.NotContains(t, statsContent.Text, "hidden-consumer")
	assert.NotContains(t, statsContent.Text, "secret-password")
	assert.NotContains(t, statsContent.Text, "secret-token")
	var statsPayload pulsarSubscriptionStatsResource
	require.NoError(t, json.Unmarshal([]byte(statsContent.Text), &statsPayload))
	assert.Equal(t, "subscriptionStats", statsPayload.Kind)
	assert.Equal(t, "sub-a", statsPayload.Subscription)
	assert.False(t, statsPayload.Partitioned)
	assert.Equal(t, "Shared", statsPayload.Stats.Type)
	assert.True(t, statsPayload.Stats.Durable)
	assert.True(t, statsPayload.Stats.BlockedOnUnackedMessages)
	assert.Equal(t, int64(42), statsPayload.Stats.MsgBacklog)
	assert.Equal(t, "team-a", statsPayload.Stats.SubscriptionProperties["owner"])
	assert.Equal(t, pulsarResourceRedactedValue, statsPayload.Stats.SubscriptionProperties["token"])

	backlogContent := readPulsarTestResource(
		t,
		ctx,
		s,
		"pulsar://admin/v2/persistent/public/default/events/subscriptions/sub-a/backlog",
	)
	var backlogPayload pulsarSubscriptionBacklogResource
	require.NoError(t, json.Unmarshal([]byte(backlogContent.Text), &backlogPayload))
	assert.Equal(t, "subscriptionBacklog", backlogPayload.Kind)
	assert.Equal(t, int64(42), backlogPayload.Backlog.MsgBacklog)
	assert.Equal(t, int64(4096), backlogPayload.Backlog.BacklogSize)
	assert.Equal(t, int64(12345), backlogPayload.Backlog.EarliestMsgPublishTimeInBacklog)

	cursorContent := readPulsarTestResource(
		t,
		ctx,
		s,
		"pulsar://admin/v2/persistent/public/default/events/subscriptions/sub-a/cursor",
	)
	assert.NotContains(t, cursorContent.Text, "hidden-large-range")
	var cursorPayload pulsarSubscriptionCursorResource
	require.NoError(t, json.Unmarshal([]byte(cursorContent.Text), &cursorPayload))
	assert.Equal(t, "subscriptionCursor", cursorPayload.Kind)
	assert.Equal(t, "1:2", cursorPayload.Cursor.MarkDeletePosition)
	assert.Equal(t, "1:3", cursorPayload.Cursor.ReadPosition)
	assert.Equal(t, int64(10), cursorPayload.Cursor.MessagesConsumedCounter)
	assert.Equal(t, map[string]int64{"safe": 1}, cursorPayload.Cursor.Properties)
	assert.Equal(t, 2, cursorPayload.Cursor.PropertiesCount)
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
		"pulsar://admin/v2/tenants/public/namespaces/extra",
		"pulsar://admin/v2/namespaces/public/default/topics/extra",
		"pulsar://admin/v2/resource-quotas/default/extra",
		"pulsar://admin/v2/resource-quotas/public/default",
		"pulsar://admin/v2/broker-stats/topics",
		"pulsar://admin/v2/brokers",
		"pulsar://admin/v2/persistent/public/default/events",
		"pulsar://admin/v2/invalid/public/default/events/stats",
		"pulsar://admin/v2/persistent/public/default/events/policies/unsupported",
		"pulsar://admin/v2/persistent/public/default/events/schema/not-a-number",
		"pulsar://admin/v2/persistent/public/default/events/schema/-1",
		"pulsar://admin/v2/persistent/public/default/events/subscriptions/sub-a",
		"pulsar://admin/v2/persistent/public/default/events/subscriptions//stats",
		"pulsar://admin/v2/persistent/public/default/events/subscriptions/sub-a/unsupported",
		"pulsar://admin/v2/non-persistent/public/default/events/subscriptions/sub-a/cursor",
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
