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
	"net/url"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/admin"
	pulsaradminauth "github.com/apache/pulsar-client-go/pulsaradmin/pkg/admin/auth"
	pulsaradminconfig "github.com/apache/pulsar-client-go/pulsaradmin/pkg/admin/config"
	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/rest"
	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	context2 "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	pulsarsession "github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

const (
	pulsarResourceContextURI             = "pulsar://context"
	pulsarResourceCatalogURI             = "pulsar://resources"
	pulsarTenantsResourceURI             = "pulsar://admin/v2/tenants"
	pulsarTenantResourceTemplateURI      = "pulsar://admin/v2/tenants/{tenant}"
	pulsarNamespacesResourceTemplateURI  = "pulsar://admin/v2/tenants/{tenant}/namespaces"
	pulsarNamespaceResourceTemplateURI   = "pulsar://admin/v2/namespaces/{tenant}/{namespace}"
	pulsarTopicsResourceTemplateURI      = "pulsar://admin/v2/namespaces/{tenant}/{namespace}/topics"
	pulsarDefaultResourceQuotaURI        = "pulsar://admin/v2/resource-quotas"
	pulsarResourceQuotaTemplateURI       = "pulsar://admin/v2/resource-quotas/{tenant}/{namespace}/{bundle}"
	pulsarClusterStatusResourceURI       = "pulsar://admin/v2/status"
	pulsarClustersResourceURI            = "pulsar://admin/v2/clusters"
	pulsarBrokerStatsSummaryResourceURI  = "pulsar://admin/v2/broker-stats/summary"
	pulsarClusterResourceTemplateURI     = "pulsar://admin/v2/clusters/{cluster}"
	pulsarBrokersResourceTemplateURI     = "pulsar://admin/v2/brokers/{cluster}"
	pulsarFailureDomainsTemplateURI      = "pulsar://admin/v2/clusters/{cluster}/failureDomains"
	pulsarFailureDomainTemplateURI       = "pulsar://admin/v2/clusters/{cluster}/failureDomains/{domain}"
	pulsarNSIsolationPoliciesTemplateURI = "pulsar://admin/v2/clusters/{cluster}/namespaceIsolationPolicies"
	pulsarNSIsolationPolicyTemplateURI   = "pulsar://admin/v2/clusters/{cluster}/namespaceIsolationPolicies/{policy}"
	pulsarTopicMetadataTemplateURI       = "pulsar://admin/v2/{domain}/{tenant}/{namespace}/{topic}/metadata"
	pulsarTopicStatsTemplateURI          = "pulsar://admin/v2/{domain}/{tenant}/{namespace}/{topic}/stats"
	pulsarTopicPartitionMetadataURI      = "pulsar://admin/v2/{domain}/{tenant}/{namespace}/{topic}/partitions"
	pulsarTopicPolicyTemplateURI         = "pulsar://admin/v2/{domain}/{tenant}/{namespace}/{topic}/policies/{policy}"
	pulsarTopicSchemaTemplateURI         = "pulsar://admin/v2/{domain}/{tenant}/{namespace}/{topic}/schema"
	pulsarTopicSchemaVersionTemplateURI  = "pulsar://admin/v2/{domain}/{tenant}/{namespace}/{topic}/schema/{version}"
	pulsarSubscriptionsTemplateURI       = "pulsar://admin/v2/{domain}/{tenant}/{namespace}/{topic}/subscriptions"
	pulsarSubscriptionStatsTemplateURI   = "pulsar://admin/v2/{domain}/{tenant}/{namespace}/{topic}/subscriptions/{subscription}/stats"
	pulsarSubscriptionBacklogTemplateURI = "pulsar://admin/v2/{domain}/{tenant}/{namespace}/{topic}/subscriptions/{subscription}/backlog"
	pulsarSubscriptionCursorTemplateURI  = "pulsar://admin/v2/persistent/{tenant}/{namespace}/{topic}/subscriptions/{subscription}/cursor"
	pulsarResourceJSONMIMEType           = "application/json"
	pulsarResourceSummaryStringLimit     = 50
	pulsarResourceRedactedValue          = "<redacted>"
)

type pulsarResourceKind string

const (
	pulsarResourceKindContext              pulsarResourceKind = "context"
	pulsarResourceKindCatalog              pulsarResourceKind = "catalog"
	pulsarResourceKindTenants              pulsarResourceKind = "tenants"
	pulsarResourceKindTenant               pulsarResourceKind = "tenant"
	pulsarResourceKindNamespaces           pulsarResourceKind = "namespaces"
	pulsarResourceKindNamespace            pulsarResourceKind = "namespace"
	pulsarResourceKindTopics               pulsarResourceKind = "topics"
	pulsarResourceKindDefaultResourceQuota pulsarResourceKind = "defaultResourceQuota"
	pulsarResourceKindResourceQuota        pulsarResourceKind = "resourceQuota"
	pulsarResourceKindStatus               pulsarResourceKind = "status"
	pulsarResourceKindClusters             pulsarResourceKind = "clusters"
	pulsarResourceKindCluster              pulsarResourceKind = "cluster"
	pulsarResourceKindBrokers              pulsarResourceKind = "brokers"
	pulsarResourceKindBrokerStatsSummary   pulsarResourceKind = "brokerStatsSummary"
	pulsarResourceKindFailureDomains       pulsarResourceKind = "failureDomains"
	pulsarResourceKindFailureDomain        pulsarResourceKind = "failureDomain"
	pulsarResourceKindNSIsolationPolicies  pulsarResourceKind = "namespaceIsolationPolicies"
	pulsarResourceKindNSIsolationPolicy    pulsarResourceKind = "namespaceIsolationPolicy"
	pulsarResourceKindTopicMetadata        pulsarResourceKind = "topicMetadata"
	pulsarResourceKindTopicStats           pulsarResourceKind = "topicStats"
	pulsarResourceKindTopicPartitions      pulsarResourceKind = "topicPartitions"
	pulsarResourceKindTopicPolicy          pulsarResourceKind = "topicPolicy"
	pulsarResourceKindTopicSchema          pulsarResourceKind = "topicSchema"
	pulsarResourceKindTopicSchemaVersion   pulsarResourceKind = "topicSchemaVersion"
	pulsarResourceKindSubscriptions        pulsarResourceKind = "subscriptions"
	pulsarResourceKindSubscriptionStats    pulsarResourceKind = "subscriptionStats"
	pulsarResourceKindSubscriptionBacklog  pulsarResourceKind = "subscriptionBacklog"
	pulsarResourceKindSubscriptionCursor   pulsarResourceKind = "subscriptionCursor"
)

type pulsarResourceURI struct {
	kind         pulsarResourceKind
	tenant       string
	namespace    string
	topic        string
	topicDomain  string
	cluster      string
	domain       string
	policy       string
	bundle       string
	subscription string
	version      int64
}

type pulsarResourceCatalog struct {
	Version   int                     `json:"version"`
	Scheme    string                  `json:"scheme"`
	Resources []pulsarCatalogResource `json:"resources"`
	Templates []pulsarCatalogTemplate `json:"templates"`
	Notes     []string                `json:"notes,omitempty"`
}

type pulsarCatalogResource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MIMEType    string `json:"mimeType"`
}

type pulsarCatalogTemplate struct {
	URITemplate string `json:"uriTemplate"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MIMEType    string `json:"mimeType"`
}

type pulsarContextResource struct {
	Kind           string                      `json:"kind"`
	URI            string                      `json:"uri"`
	ServiceURL     string                      `json:"serviceUrl,omitempty"`
	WebServiceURL  string                      `json:"webServiceUrl,omitempty"`
	Authentication pulsarAuthenticationSummary `json:"authentication"`
	TLS            pulsarTLSSummary            `json:"tls"`
}

type pulsarAuthenticationSummary struct {
	Configured bool   `json:"configured"`
	Method     string `json:"method"`
	Plugin     string `json:"plugin,omitempty"`
}

type pulsarTLSSummary struct {
	AllowInsecureConnection    bool `json:"allowInsecureConnection"`
	EnableHostnameVerification bool `json:"enableHostnameVerification"`
	TrustCertsFileConfigured   bool `json:"trustCertsFileConfigured"`
	ClientCertFileConfigured   bool `json:"clientCertFileConfigured"`
	ClientKeyFileConfigured    bool `json:"clientKeyFileConfigured"`
}

type pulsarTenantCollectionResource struct {
	Kind    string   `json:"kind"`
	URI     string   `json:"uri"`
	Tenants []string `json:"tenants"`
	Count   int      `json:"count"`
}

type pulsarTenantResource struct {
	Kind   string           `json:"kind"`
	URI    string           `json:"uri"`
	Tenant string           `json:"tenant"`
	Data   utils.TenantData `json:"data"`
}

type pulsarNamespaceCollectionResource struct {
	Kind       string   `json:"kind"`
	URI        string   `json:"uri"`
	Tenant     string   `json:"tenant"`
	Namespaces []string `json:"namespaces"`
	Count      int      `json:"count"`
}

type pulsarNamespaceResource struct {
	Kind      string          `json:"kind"`
	URI       string          `json:"uri"`
	Tenant    string          `json:"tenant"`
	Namespace string          `json:"namespace"`
	Policies  *utils.Policies `json:"policies"`
}

type pulsarTopicCollectionResource struct {
	Kind      string   `json:"kind"`
	URI       string   `json:"uri"`
	Tenant    string   `json:"tenant"`
	Namespace string   `json:"namespace"`
	Topics    []string `json:"topics"`
	Count     int      `json:"count"`
}

type pulsarResourceQuotaResource struct {
	Kind      string               `json:"kind"`
	URI       string               `json:"uri"`
	Scope     string               `json:"scope"`
	Tenant    string               `json:"tenant,omitempty"`
	Namespace string               `json:"namespace,omitempty"`
	Bundle    string               `json:"bundle,omitempty"`
	Quota     *utils.ResourceQuota `json:"quota"`
}

type pulsarClusterStatusResource struct {
	Kind   string `json:"kind"`
	URI    string `json:"uri"`
	Status string `json:"status"`
}

type pulsarClusterCollectionResource struct {
	Kind     string   `json:"kind"`
	URI      string   `json:"uri"`
	Clusters []string `json:"clusters"`
	Count    int      `json:"count"`
}

type pulsarClusterResource struct {
	Kind    string                   `json:"kind"`
	URI     string                   `json:"uri"`
	Cluster string                   `json:"cluster"`
	Data    pulsarClusterDataSummary `json:"data"`
}

type pulsarClusterDataSummary struct {
	Name                                 string   `json:"name"`
	ServiceURL                           string   `json:"serviceUrl,omitempty"`
	ServiceURLTLS                        string   `json:"serviceUrlTls,omitempty"`
	BrokerServiceURL                     string   `json:"brokerServiceUrl,omitempty"`
	BrokerServiceURLTLS                  string   `json:"brokerServiceUrlTls,omitempty"`
	PeerClusterNames                     []string `json:"peerClusterNames,omitempty"`
	AuthenticationPlugin                 string   `json:"authenticationPlugin,omitempty"`
	AuthenticationParametersConfigured   bool     `json:"authenticationParametersConfigured"`
	BrokerClientTrustCertsFileConfigured bool     `json:"brokerClientTrustCertsFileConfigured"`
	BrokerClientTLSEnabled               bool     `json:"brokerClientTlsEnabled"`
}

type pulsarBrokerCollectionResource struct {
	Kind    string   `json:"kind"`
	URI     string   `json:"uri"`
	Cluster string   `json:"cluster"`
	Brokers []string `json:"brokers"`
	Count   int      `json:"count"`
}

type pulsarBrokerStatsSummaryResource struct {
	Kind              string                         `json:"kind"`
	URI               string                         `json:"uri"`
	MonitoringMetrics pulsarMonitoringMetricsSummary `json:"monitoringMetrics"`
	LoadReport        pulsarBrokerLoadReportSummary  `json:"loadReport"`
}

type pulsarMonitoringMetricsSummary struct {
	Count         int      `json:"count"`
	MetricNames   []string `json:"metricNames,omitempty"`
	DimensionKeys []string `json:"dimensionKeys,omitempty"`
}

type pulsarBrokerLoadReportSummary struct {
	Available                  bool                       `json:"available"`
	WebServiceURL              string                     `json:"webServiceUrl,omitempty"`
	WebServiceURLTLS           string                     `json:"webServiceUrlTls,omitempty"`
	PulsarServiceURL           string                     `json:"pulsarServiceUrl,omitempty"`
	PulsarServiceURLTLS        string                     `json:"pulsarServiceUrlTls,omitempty"`
	PersistentTopicsEnabled    bool                       `json:"persistentTopicsEnabled"`
	NonPersistentTopicsEnabled bool                       `json:"nonPersistentTopicsEnabled"`
	CPU                        pulsarResourceUsageSummary `json:"cpu"`
	Memory                     pulsarResourceUsageSummary `json:"memory"`
	DirectMemory               pulsarResourceUsageSummary `json:"directMemory"`
	BandwidthIn                pulsarResourceUsageSummary `json:"bandwidthIn"`
	BandwidthOut               pulsarResourceUsageSummary `json:"bandwidthOut"`
	MsgThroughputIn            float64                    `json:"msgThroughputIn"`
	MsgThroughputOut           float64                    `json:"msgThroughputOut"`
	MsgRateIn                  float64                    `json:"msgRateIn"`
	MsgRateOut                 float64                    `json:"msgRateOut"`
	LastUpdate                 int64                      `json:"lastUpdate,omitempty"`
	NumTopics                  int                        `json:"numTopics"`
	NumBundles                 int                        `json:"numBundles"`
	NumConsumers               int                        `json:"numConsumers"`
	NumProducers               int                        `json:"numProducers"`
	BundleCount                int                        `json:"bundleCount"`
	LastBundleGainsCount       int                        `json:"lastBundleGainsCount"`
	LastBundleLossesCount      int                        `json:"lastBundleLossesCount"`
	BrokerVersionString        string                     `json:"brokerVersionString,omitempty"`
	LoadReportType             string                     `json:"loadReportType,omitempty"`
	ProtocolCount              int                        `json:"protocolCount"`
}

type pulsarResourceUsageSummary struct {
	Usage        float64 `json:"usage"`
	Limit        float64 `json:"limit"`
	PercentUsage float32 `json:"percentUsage"`
}

type pulsarFailureDomainCollectionResource struct {
	Kind           string                       `json:"kind"`
	URI            string                       `json:"uri"`
	Cluster        string                       `json:"cluster"`
	FailureDomains []pulsarFailureDomainSummary `json:"failureDomains"`
	Count          int                          `json:"count"`
}

type pulsarFailureDomainResource struct {
	Kind          string                     `json:"kind"`
	URI           string                     `json:"uri"`
	Cluster       string                     `json:"cluster"`
	FailureDomain pulsarFailureDomainSummary `json:"failureDomain"`
}

type pulsarFailureDomainSummary struct {
	Name    string   `json:"name"`
	Brokers []string `json:"brokers"`
}

type pulsarNamespaceIsolationPolicyCollectionResource struct {
	Kind     string                                  `json:"kind"`
	URI      string                                  `json:"uri"`
	Cluster  string                                  `json:"cluster"`
	Policies []pulsarNamespaceIsolationPolicySummary `json:"policies"`
	Count    int                                     `json:"count"`
}

type pulsarNamespaceIsolationPolicyResource struct {
	Kind    string                         `json:"kind"`
	URI     string                         `json:"uri"`
	Cluster string                         `json:"cluster"`
	Policy  pulsarNamespaceIsolationPolicy `json:"policy"`
}

type pulsarNamespaceIsolationPolicySummary struct {
	Name                   string `json:"name"`
	NamespacesCount        int    `json:"namespacesCount"`
	PrimaryBrokersCount    int    `json:"primaryBrokersCount"`
	SecondaryBrokersCount  int    `json:"secondaryBrokersCount"`
	AutoFailoverPolicyType string `json:"autoFailoverPolicyType,omitempty"`
}

type pulsarNamespaceIsolationPolicy struct {
	Name string                       `json:"name"`
	Data utils.NamespaceIsolationData `json:"data"`
}

type pulsarTopicMetadataResource struct {
	Kind            string            `json:"kind"`
	URI             string            `json:"uri"`
	Topic           string            `json:"topic"`
	Domain          string            `json:"domain"`
	Tenant          string            `json:"tenant"`
	Namespace       string            `json:"namespace"`
	Name            string            `json:"name"`
	PartitionIndex  int               `json:"partitionIndex"`
	Properties      map[string]string `json:"properties,omitempty"`
	PropertiesCount int               `json:"propertiesCount"`
}

type pulsarTopicPartitionMetadataResource struct {
	Kind        string                         `json:"kind"`
	URI         string                         `json:"uri"`
	Topic       string                         `json:"topic"`
	Metadata    utils.PartitionedTopicMetadata `json:"metadata"`
	Partitioned bool                           `json:"partitioned"`
}

type pulsarTopicStatsResource struct {
	Kind                string                  `json:"kind"`
	URI                 string                  `json:"uri"`
	Topic               string                  `json:"topic"`
	Partitioned         bool                    `json:"partitioned"`
	PartitionCount      int                     `json:"partitionCount"`
	PartitionStatsCount int                     `json:"partitionStatsCount,omitempty"`
	Stats               pulsarTopicStatsSummary `json:"stats"`
}

type pulsarTopicStatsSummary struct {
	BacklogSize          int64   `json:"backlogSize,omitempty"`
	MsgCounterIn         int64   `json:"msgInCounter,omitempty"`
	MsgCounterOut        int64   `json:"msgOutCounter,omitempty"`
	MsgRateIn            float64 `json:"msgRateIn"`
	MsgRateOut           float64 `json:"msgRateOut"`
	MsgThroughputIn      float64 `json:"msgThroughputIn"`
	MsgThroughputOut     float64 `json:"msgThroughputOut"`
	AverageMsgSize       float64 `json:"averageMsgSize"`
	StorageSize          int64   `json:"storageSize"`
	PublisherCount       int     `json:"publisherCount"`
	SubscriptionCount    int     `json:"subscriptionCount"`
	ReplicationCount     int     `json:"replicationCount"`
	DeduplicationStatus  string  `json:"deduplicationStatus,omitempty"`
	TopicCreationTime    int64   `json:"topicCreationTimeStamp,omitempty"`
	LastPublishTimestamp int64   `json:"lastPublishTimestamp,omitempty"`
}

type pulsarTopicPolicyResource struct {
	Kind   string `json:"kind"`
	URI    string `json:"uri"`
	Topic  string `json:"topic"`
	Policy string `json:"policy"`
	Value  any    `json:"value"`
}

type pulsarTopicSchemaResource struct {
	Kind    string                   `json:"kind"`
	URI     string                   `json:"uri"`
	Topic   string                   `json:"topic"`
	Version int64                    `json:"version"`
	Schema  pulsarTopicSchemaSummary `json:"schema"`
}

type pulsarTopicSchemaSummary struct {
	Name            string            `json:"name"`
	Type            string            `json:"type"`
	Schema          string            `json:"schema"`
	Properties      map[string]string `json:"properties,omitempty"`
	PropertiesCount int               `json:"propertiesCount"`
	Timestamp       int64             `json:"timestamp,omitempty"`
}

type pulsarSubscriptionCollectionResource struct {
	Kind          string   `json:"kind"`
	URI           string   `json:"uri"`
	Topic         string   `json:"topic"`
	Domain        string   `json:"domain"`
	Tenant        string   `json:"tenant"`
	Namespace     string   `json:"namespace"`
	Subscriptions []string `json:"subscriptions"`
	Count         int      `json:"count"`
}

type pulsarSubscriptionStatsResource struct {
	Kind         string                         `json:"kind"`
	URI          string                         `json:"uri"`
	Topic        string                         `json:"topic"`
	Subscription string                         `json:"subscription"`
	Partitioned  bool                           `json:"partitioned"`
	Stats        pulsarSubscriptionStatsSummary `json:"stats"`
}

type pulsarSubscriptionStatsSummary struct {
	Type                                      string            `json:"type,omitempty"`
	Durable                                   bool              `json:"durable"`
	Replicated                                bool              `json:"replicated"`
	BlockedOnUnackedMessages                  bool              `json:"blockedOnUnackedMessages"`
	MsgRateOut                                float64           `json:"msgRateOut"`
	MsgThroughputOut                          float64           `json:"msgThroughputOut"`
	MsgRateRedeliver                          float64           `json:"msgRateRedeliver"`
	MsgRateExpired                            float64           `json:"msgRateExpired"`
	MsgBacklog                                int64             `json:"msgBacklog"`
	MsgBacklogNoDelayed                       int64             `json:"msgBacklogNoDelayed"`
	MsgDelayed                                int64             `json:"msgDelayed"`
	UnackedMessages                           int64             `json:"unackedMessages"`
	BytesOutCounter                           int64             `json:"bytesOutCounter"`
	MsgOutCounter                             int64             `json:"msgOutCounter"`
	MessageAckRate                            float64           `json:"messageAckRate"`
	ChunkedMessageRate                        float64           `json:"chunkedMessageRate"`
	BacklogSize                               int64             `json:"backlogSize"`
	EarliestMsgPublishTimeInBacklog           int64             `json:"earliestMsgPublishTimeInBacklog,omitempty"`
	TotalMsgExpired                           int64             `json:"totalMsgExpired"`
	LastExpireTimestamp                       int64             `json:"lastExpireTimestamp,omitempty"`
	LastConsumedFlowTimestamp                 int64             `json:"lastConsumedFlowTimestamp,omitempty"`
	LastConsumedTimestamp                     int64             `json:"lastConsumedTimestamp,omitempty"`
	LastAckedTimestamp                        int64             `json:"lastAckedTimestamp,omitempty"`
	LastMarkDeleteAdvancedTimestamp           int64             `json:"lastMarkDeleteAdvancedTimestamp,omitempty"`
	AllowOutOfOrderDelivery                   bool              `json:"allowOutOfOrderDelivery"`
	NonContiguousDeletedMessagesRanges        int               `json:"nonContiguousDeletedMessagesRanges"`
	NonContiguousDeletedMessagesRangesSrzSize int               `json:"nonContiguousDeletedMessagesRangesSerializedSize"`
	DelayedMessageIndexSizeInBytes            int64             `json:"delayedMessageIndexSizeInBytes"`
	FilterProcessedMsgCount                   int64             `json:"filterProcessedMsgCount"`
	FilterAcceptedMsgCount                    int64             `json:"filterAcceptedMsgCount"`
	FilterRejectedMsgCount                    int64             `json:"filterRejectedMsgCount"`
	FilterRescheduledMsgCount                 int64             `json:"filterRescheduledMsgCount"`
	SubscriptionProperties                    map[string]string `json:"subscriptionProperties,omitempty"`
	SubscriptionPropertiesCount               int               `json:"subscriptionPropertiesCount"`
}

type pulsarSubscriptionBacklogResource struct {
	Kind         string                           `json:"kind"`
	URI          string                           `json:"uri"`
	Topic        string                           `json:"topic"`
	Subscription string                           `json:"subscription"`
	Partitioned  bool                             `json:"partitioned"`
	Backlog      pulsarSubscriptionBacklogSummary `json:"backlog"`
}

type pulsarSubscriptionBacklogSummary struct {
	MsgBacklog                      int64 `json:"msgBacklog"`
	MsgBacklogNoDelayed             int64 `json:"msgBacklogNoDelayed"`
	BacklogSize                     int64 `json:"backlogSize"`
	MsgDelayed                      int64 `json:"msgDelayed"`
	UnackedMessages                 int64 `json:"unackedMessages"`
	EarliestMsgPublishTimeInBacklog int64 `json:"earliestMsgPublishTimeInBacklog,omitempty"`
	DelayedMessageIndexSizeInBytes  int64 `json:"delayedMessageIndexSizeInBytes"`
}

type pulsarSubscriptionCursorResource struct {
	Kind         string                          `json:"kind"`
	URI          string                          `json:"uri"`
	Topic        string                          `json:"topic"`
	Subscription string                          `json:"subscription"`
	Cursor       pulsarSubscriptionCursorSummary `json:"cursor"`
}

type pulsarSubscriptionCursorSummary struct {
	MarkDeletePosition                       string           `json:"markDeletePosition,omitempty"`
	ReadPosition                             string           `json:"readPosition,omitempty"`
	WaitingReadOp                            bool             `json:"waitingReadOp"`
	PendingReadOps                           int              `json:"pendingReadOps"`
	MessagesConsumedCounter                  int64            `json:"messagesConsumedCounter"`
	CursorLedger                             int64            `json:"cursorLedger"`
	CursorLedgerLastEntry                    int64            `json:"cursorLedgerLastEntry"`
	LastLedgerSwitchTimestamp                string           `json:"lastLedgerSwitchTimestamp,omitempty"`
	State                                    string           `json:"state,omitempty"`
	NumberOfEntriesSinceFirstNotAckedMessage int64            `json:"numberOfEntriesSinceFirstNotAckedMessage"`
	TotalNonContiguousDeletedMessagesRange   int              `json:"totalNonContiguousDeletedMessagesRange"`
	Properties                               map[string]int64 `json:"properties,omitempty"`
	PropertiesCount                          int              `json:"propertiesCount"`
}

// PulsarAddResources registers the read-only Pulsar MCP resource surface.
func PulsarAddResources(s *server.MCPServer, features []string) {
	resourceRegistrations, templateRegistrations := buildPulsarResourceRegistrations(features)
	if len(resourceRegistrations) == 0 && len(templateRegistrations) == 0 {
		return
	}

	catalog := buildPulsarResourceCatalog(resourceRegistrations, templateRegistrations)
	baseResources := []server.ServerResource{
		server.ServerResource{
			Resource: mcp.NewResource(pulsarResourceContextURI, "Pulsar Context",
				mcp.WithResourceDescription("Current Pulsar session connection metadata with authentication material redacted."),
				mcp.WithMIMEType(pulsarResourceJSONMIMEType),
			),
			Handler: handlePulsarContextResource,
		},
		server.ServerResource{
			Resource: mcp.NewResource(pulsarResourceCatalogURI, "Pulsar Resource Catalog",
				mcp.WithResourceDescription("Stable catalog of Pulsar MCP resource URIs and URI templates."),
				mcp.WithMIMEType(pulsarResourceJSONMIMEType),
			),
			Handler: func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
				return handlePulsarCatalogResource(ctx, request, catalog)
			},
		},
	}

	allResources := append(baseResources, resourceRegistrations...)
	s.AddResources(allResources...)
	if len(templateRegistrations) > 0 {
		s.AddResourceTemplates(templateRegistrations...)
	}
}

func buildPulsarResourceRegistrations(features []string) ([]server.ServerResource, []server.ServerResourceTemplate) {
	var resources []server.ServerResource
	var templates []server.ServerResourceTemplate

	if pulsarResourceFeatureEnabled(features, FeaturePulsarAdminBrokersStatus) {
		resources = append(resources, server.ServerResource{
			Resource: mcp.NewResource(pulsarClusterStatusResourceURI, "Pulsar Cluster Status",
				mcp.WithResourceDescription("Read the Pulsar broker or proxy status endpoint for the current session."),
				mcp.WithMIMEType(pulsarResourceJSONMIMEType),
			),
			Handler: handlePulsarClusterStatusResource,
		})
	}

	if pulsarResourceFeatureEnabled(features, FeaturePulsarAdminClusters) {
		resources = append(resources, server.ServerResource{
			Resource: mcp.NewResource(pulsarClustersResourceURI, "Pulsar Clusters",
				mcp.WithResourceDescription("List Pulsar clusters known to the current admin endpoint."),
				mcp.WithMIMEType(pulsarResourceJSONMIMEType),
			),
			Handler: handlePulsarClustersResource,
		})
		templates = append(templates,
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarClusterResourceTemplateURI, "Pulsar Cluster",
					mcp.WithTemplateDescription("Get sanitized configuration for a Pulsar cluster."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarClusterResource,
			},
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarFailureDomainsTemplateURI, "Pulsar Failure Domains",
					mcp.WithTemplateDescription("List failure domains for a Pulsar cluster."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarFailureDomainsResource,
			},
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarFailureDomainTemplateURI, "Pulsar Failure Domain",
					mcp.WithTemplateDescription("Get a failure domain for a Pulsar cluster."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarFailureDomainResource,
			},
		)
	}

	if pulsarResourceFeatureEnabled(features, FeaturePulsarAdminBrokers) {
		templates = append(templates, server.ServerResourceTemplate{
			Template: mcp.NewResourceTemplate(pulsarBrokersResourceTemplateURI, "Pulsar Brokers by Cluster",
				mcp.WithTemplateDescription("List active brokers for a Pulsar cluster."),
				mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
			),
			Handler: handlePulsarBrokersResource,
		})
	}

	if pulsarResourceFeatureEnabled(features, FeaturePulsarAdminBrokerStats) {
		resources = append(resources, server.ServerResource{
			Resource: mcp.NewResource(pulsarBrokerStatsSummaryResourceURI, "Pulsar Broker Stats Summary",
				mcp.WithResourceDescription("Bounded summary of broker monitoring metrics and load report for the current session."),
				mcp.WithMIMEType(pulsarResourceJSONMIMEType),
			),
			Handler: handlePulsarBrokerStatsSummaryResource,
		})
	}

	if pulsarResourceFeatureEnabled(features, FeaturePulsarAdminNsIsolationPolicy) {
		templates = append(templates,
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarNSIsolationPoliciesTemplateURI, "Pulsar Namespace Isolation Policies",
					mcp.WithTemplateDescription("List namespace isolation policies for a Pulsar cluster."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarNamespaceIsolationPoliciesResource,
			},
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarNSIsolationPolicyTemplateURI, "Pulsar Namespace Isolation Policy",
					mcp.WithTemplateDescription("Get a namespace isolation policy for a Pulsar cluster."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarNamespaceIsolationPolicyResource,
			},
		)
	}

	if pulsarResourceFeatureEnabled(features, FeaturePulsarAdminTenants) {
		resources = append(resources, server.ServerResource{
			Resource: mcp.NewResource(pulsarTenantsResourceURI, "Pulsar Tenants",
				mcp.WithResourceDescription("List tenants known to the current Pulsar admin endpoint."),
				mcp.WithMIMEType(pulsarResourceJSONMIMEType),
			),
			Handler: handlePulsarTenantsResource,
		})
		templates = append(templates, server.ServerResourceTemplate{
			Template: mcp.NewResourceTemplate(pulsarTenantResourceTemplateURI, "Pulsar Tenant",
				mcp.WithTemplateDescription("Get configuration for a Pulsar tenant."),
				mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
			),
			Handler: handlePulsarTenantResource,
		})
	}

	if pulsarResourceFeatureEnabled(features, FeaturePulsarAdminNamespaces) {
		templates = append(templates, server.ServerResourceTemplate{
			Template: mcp.NewResourceTemplate(pulsarNamespacesResourceTemplateURI, "Pulsar Namespaces by Tenant",
				mcp.WithTemplateDescription("List namespaces for a Pulsar tenant."),
				mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
			),
			Handler: handlePulsarNamespacesResource,
		})
	}

	if pulsarResourceFeatureEnabled(features, FeaturePulsarAdminNamespacePolicy) {
		templates = append(templates, server.ServerResourceTemplate{
			Template: mcp.NewResourceTemplate(pulsarNamespaceResourceTemplateURI, "Pulsar Namespace Policies",
				mcp.WithTemplateDescription("Get policies for a Pulsar namespace."),
				mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
			),
			Handler: handlePulsarNamespaceResource,
		})
	}

	if pulsarResourceFeatureEnabled(features, FeaturePulsarAdminTopics) {
		templates = append(templates,
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarTopicsResourceTemplateURI, "Pulsar Topics by Namespace",
					mcp.WithTemplateDescription("List topics for a Pulsar namespace."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarTopicsResource,
			},
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarTopicMetadataTemplateURI, "Pulsar Topic Metadata",
					mcp.WithTemplateDescription("Get parsed topic identity and sanitized topic properties."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarTopicMetadataResource,
			},
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarTopicStatsTemplateURI, "Pulsar Topic Stats Summary",
					mcp.WithTemplateDescription("Get a bounded topic statistics summary without publisher or consumer details."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarTopicStatsResource,
			},
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarTopicPartitionMetadataURI, "Pulsar Topic Partition Metadata",
					mcp.WithTemplateDescription("Get partition metadata for a Pulsar topic."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarTopicPartitionMetadataResource,
			},
		)
	}

	if pulsarResourceFeatureEnabled(features, FeaturePulsarAdminTopicPolicy) {
		templates = append(templates, server.ServerResourceTemplate{
			Template: mcp.NewResourceTemplate(pulsarTopicPolicyTemplateURI, "Pulsar Topic Policy",
				mcp.WithTemplateDescription("Get a read-only topic policy value for a Pulsar topic."),
				mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
			),
			Handler: handlePulsarTopicPolicyResource,
		})
	}

	if pulsarResourceFeatureEnabled(features, FeaturePulsarAdminSchemas) {
		templates = append(templates,
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarTopicSchemaTemplateURI, "Pulsar Topic Latest Schema",
					mcp.WithTemplateDescription("Get the latest schema and version for a Pulsar topic."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarTopicSchemaResource,
			},
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarTopicSchemaVersionTemplateURI, "Pulsar Topic Schema Version",
					mcp.WithTemplateDescription("Get a specific schema version for a Pulsar topic."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarTopicSchemaVersionResource,
			},
		)
	}

	if pulsarResourceFeatureEnabled(features, FeaturePulsarAdminSubscriptions) {
		templates = append(templates,
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarSubscriptionsTemplateURI, "Pulsar Topic Subscriptions",
					mcp.WithTemplateDescription("List subscriptions for a Pulsar topic."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarSubscriptionsResource,
			},
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarSubscriptionStatsTemplateURI, "Pulsar Subscription Stats Summary",
					mcp.WithTemplateDescription("Get bounded statistics for one Pulsar subscription without consumer details."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarSubscriptionStatsResource,
			},
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarSubscriptionBacklogTemplateURI, "Pulsar Subscription Backlog Summary",
					mcp.WithTemplateDescription("Get backlog counters for one Pulsar subscription without changing cursor state."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarSubscriptionBacklogResource,
			},
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarSubscriptionCursorTemplateURI, "Pulsar Subscription Cursor Summary",
					mcp.WithTemplateDescription("Get persistent topic cursor positions for one Pulsar subscription."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarSubscriptionCursorResource,
			},
		)
	}

	if pulsarResourceFeatureEnabled(features, FeaturePulsarAdminResourceQuotas) {
		resources = append(resources, server.ServerResource{
			Resource: mcp.NewResource(pulsarDefaultResourceQuotaURI, "Pulsar Default Resource Quota",
				mcp.WithResourceDescription("Get the default resource quota for new namespace bundles."),
				mcp.WithMIMEType(pulsarResourceJSONMIMEType),
			),
			Handler: handlePulsarDefaultResourceQuotaResource,
		})
		templates = append(templates, server.ServerResourceTemplate{
			Template: mcp.NewResourceTemplate(pulsarResourceQuotaTemplateURI, "Pulsar Namespace Bundle Resource Quota",
				mcp.WithTemplateDescription("Get the resource quota for a Pulsar namespace bundle."),
				mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
			),
			Handler: handlePulsarResourceQuotaResource,
		})
	}

	return resources, templates
}

func pulsarResourceFeatureEnabled(features []string, resourceFeatures ...Feature) bool {
	requiredFeatures := append([]Feature{
		FeatureAll,
		FeatureAllPulsar,
		FeaturePulsarAdmin,
	}, resourceFeatures...)
	for _, feature := range requiredFeatures {
		if slices.Contains(features, string(feature)) {
			return true
		}
	}
	return false
}

func handlePulsarContextResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindContext {
		return nil, fmt.Errorf("unsupported Pulsar context resource URI %q", request.Params.URI)
	}
	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}

	resource, err := buildPulsarContextResource(request.Params.URI, session)
	if err != nil {
		return nil, err
	}
	return newPulsarJSONResourceContents(request.Params.URI, resource)
}

func handlePulsarCatalogResource(_ context.Context, request mcp.ReadResourceRequest, catalog pulsarResourceCatalog) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindCatalog {
		return nil, fmt.Errorf("unsupported Pulsar resource catalog URI %q", request.Params.URI)
	}
	return newPulsarJSONResourceContents(request.Params.URI, catalog)
}

func handlePulsarTenantsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindTenants {
		return nil, fmt.Errorf("unsupported Pulsar tenants resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}

	tenants, err := adminClient.Tenants().List()
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", err)
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarTenantCollectionResource{
		Kind:    string(parsed.kind),
		URI:     request.Params.URI,
		Tenants: tenants,
		Count:   len(tenants),
	})
}

func handlePulsarTenantResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindTenant {
		return nil, fmt.Errorf("unsupported Pulsar tenant resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}

	tenantData, err := adminClient.Tenants().Get(parsed.tenant)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant %q: %w", parsed.tenant, err)
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarTenantResource{
		Kind:   string(parsed.kind),
		URI:    request.Params.URI,
		Tenant: parsed.tenant,
		Data:   tenantData,
	})
}

func handlePulsarNamespacesResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindNamespaces {
		return nil, fmt.Errorf("unsupported Pulsar namespaces resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}

	namespaces, err := adminClient.Namespaces().GetNamespaces(parsed.tenant)
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces for tenant %q: %w", parsed.tenant, err)
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarNamespaceCollectionResource{
		Kind:       string(parsed.kind),
		URI:        request.Params.URI,
		Tenant:     parsed.tenant,
		Namespaces: namespaces,
		Count:      len(namespaces),
	})
}

func handlePulsarNamespaceResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindNamespace {
		return nil, fmt.Errorf("unsupported Pulsar namespace resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}

	namespaceName := parsed.tenant + "/" + parsed.namespace
	policies, err := adminClient.Namespaces().GetPolicies(namespaceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get policies for namespace %q: %w", namespaceName, err)
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarNamespaceResource{
		Kind:      string(parsed.kind),
		URI:       request.Params.URI,
		Tenant:    parsed.tenant,
		Namespace: parsed.namespace,
		Policies:  policies,
	})
}

func handlePulsarTopicsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindTopics {
		return nil, fmt.Errorf("unsupported Pulsar topics resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}

	namespaceName := parsed.tenant + "/" + parsed.namespace
	topics, err := adminClient.Namespaces().GetTopics(namespaceName)
	if err != nil {
		return nil, fmt.Errorf("failed to list topics for namespace %q: %w", namespaceName, err)
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarTopicCollectionResource{
		Kind:      string(parsed.kind),
		URI:       request.Params.URI,
		Tenant:    parsed.tenant,
		Namespace: parsed.namespace,
		Topics:    topics,
		Count:     len(topics),
	})
}

func handlePulsarDefaultResourceQuotaResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindDefaultResourceQuota {
		return nil, fmt.Errorf("unsupported Pulsar default resource quota URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}

	quota, err := adminClient.ResourceQuotas().GetDefaultResourceQuota()
	if err != nil {
		return nil, fmt.Errorf("failed to get default resource quota: %w", err)
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarResourceQuotaResource{
		Kind:  string(parsed.kind),
		URI:   request.Params.URI,
		Scope: "default",
		Quota: quota,
	})
}

func handlePulsarResourceQuotaResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindResourceQuota {
		return nil, fmt.Errorf("unsupported Pulsar resource quota URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}

	namespaceName := parsed.tenant + "/" + parsed.namespace
	quota, err := adminClient.ResourceQuotas().GetNamespaceBundleResourceQuota(namespaceName, parsed.bundle)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get resource quota for namespace %q bundle %q: %w",
			namespaceName,
			parsed.bundle,
			err,
		)
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarResourceQuotaResource{
		Kind:      string(parsed.kind),
		URI:       request.Params.URI,
		Scope:     "namespaceBundle",
		Tenant:    parsed.tenant,
		Namespace: parsed.namespace,
		Bundle:    parsed.bundle,
		Quota:     quota,
	})
}

func handlePulsarClusterStatusResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindStatus {
		return nil, fmt.Errorf("unsupported Pulsar status resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	status, err := getPulsarClusterStatus(session)
	if err != nil {
		return nil, err
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarClusterStatusResource{
		Kind:   string(parsed.kind),
		URI:    request.Params.URI,
		Status: status,
	})
}

func handlePulsarClustersResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindClusters {
		return nil, fmt.Errorf("unsupported Pulsar clusters resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}

	clusters, err := adminClient.Clusters().List()
	if err != nil {
		return nil, fmt.Errorf("failed to list clusters: %w", err)
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarClusterCollectionResource{
		Kind:     string(parsed.kind),
		URI:      request.Params.URI,
		Clusters: clusters,
		Count:    len(clusters),
	})
}

func handlePulsarClusterResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindCluster {
		return nil, fmt.Errorf("unsupported Pulsar cluster resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}

	clusterData, err := adminClient.Clusters().Get(parsed.cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster %q: %w", parsed.cluster, err)
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarClusterResource{
		Kind:    string(parsed.kind),
		URI:     request.Params.URI,
		Cluster: parsed.cluster,
		Data:    sanitizePulsarClusterData(parsed.cluster, clusterData),
	})
}

func handlePulsarBrokersResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindBrokers {
		return nil, fmt.Errorf("unsupported Pulsar brokers resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}

	brokers, err := adminClient.Brokers().GetActiveBrokers(parsed.cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to list brokers for cluster %q: %w", parsed.cluster, err)
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarBrokerCollectionResource{
		Kind:    string(parsed.kind),
		URI:     request.Params.URI,
		Cluster: parsed.cluster,
		Brokers: brokers,
		Count:   len(brokers),
	})
}

func handlePulsarBrokerStatsSummaryResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindBrokerStatsSummary {
		return nil, fmt.Errorf("unsupported Pulsar broker stats summary resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}

	metrics, err := adminClient.BrokerStats().GetMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to get broker monitoring metrics: %w", err)
	}
	loadReport, err := adminClient.BrokerStats().GetLoadReport()
	if err != nil {
		return nil, fmt.Errorf("failed to get broker load report: %w", err)
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarBrokerStatsSummaryResource{
		Kind:              string(parsed.kind),
		URI:               request.Params.URI,
		MonitoringMetrics: summarizePulsarMonitoringMetrics(metrics),
		LoadReport:        summarizePulsarBrokerLoadReport(loadReport),
	})
}

func handlePulsarFailureDomainsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindFailureDomains {
		return nil, fmt.Errorf("unsupported Pulsar failure domains resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}

	failureDomains, err := adminClient.Clusters().ListFailureDomains(parsed.cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to list failure domains for cluster %q: %w", parsed.cluster, err)
	}
	summaries := summarizePulsarFailureDomains(failureDomains)

	return newPulsarJSONResourceContents(request.Params.URI, pulsarFailureDomainCollectionResource{
		Kind:           string(parsed.kind),
		URI:            request.Params.URI,
		Cluster:        parsed.cluster,
		FailureDomains: summaries,
		Count:          len(summaries),
	})
}

func handlePulsarFailureDomainResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindFailureDomain {
		return nil, fmt.Errorf("unsupported Pulsar failure domain resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}

	failureDomain, err := adminClient.Clusters().GetFailureDomain(parsed.cluster, parsed.domain)
	if err != nil {
		return nil, fmt.Errorf("failed to get failure domain %q for cluster %q: %w", parsed.domain, parsed.cluster, err)
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarFailureDomainResource{
		Kind:          string(parsed.kind),
		URI:           request.Params.URI,
		Cluster:       parsed.cluster,
		FailureDomain: summarizePulsarFailureDomain(parsed.domain, failureDomain),
	})
}

func handlePulsarNamespaceIsolationPoliciesResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindNSIsolationPolicies {
		return nil, fmt.Errorf("unsupported Pulsar namespace isolation policies resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}

	policies, err := adminClient.NsIsolationPolicy().GetNamespaceIsolationPolicies(parsed.cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to list namespace isolation policies for cluster %q: %w", parsed.cluster, err)
	}
	summaries := summarizePulsarNamespaceIsolationPolicies(policies)

	return newPulsarJSONResourceContents(request.Params.URI, pulsarNamespaceIsolationPolicyCollectionResource{
		Kind:     string(parsed.kind),
		URI:      request.Params.URI,
		Cluster:  parsed.cluster,
		Policies: summaries,
		Count:    len(summaries),
	})
}

func handlePulsarNamespaceIsolationPolicyResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindNSIsolationPolicy {
		return nil, fmt.Errorf("unsupported Pulsar namespace isolation policy resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}

	policy, err := adminClient.NsIsolationPolicy().GetNamespaceIsolationPolicy(parsed.cluster, parsed.policy)
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace isolation policy %q for cluster %q: %w", parsed.policy, parsed.cluster, err)
	}
	if policy == nil {
		return nil, fmt.Errorf("namespace isolation policy %q for cluster %q is empty", parsed.policy, parsed.cluster)
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarNamespaceIsolationPolicyResource{
		Kind:    string(parsed.kind),
		URI:     request.Params.URI,
		Cluster: parsed.cluster,
		Policy: pulsarNamespaceIsolationPolicy{
			Name: parsed.policy,
			Data: *policy,
		},
	})
}

func handlePulsarTopicMetadataResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindTopicMetadata {
		return nil, fmt.Errorf("unsupported Pulsar topic metadata resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}
	topicName, err := buildPulsarResourceTopicName(parsed)
	if err != nil {
		return nil, err
	}

	properties, err := adminClient.Topics().GetProperties(*topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to get properties for topic %q: %w", topicName.String(), err)
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarTopicMetadataResource{
		Kind:            string(parsed.kind),
		URI:             request.Params.URI,
		Topic:           topicName.String(),
		Domain:          parsed.topicDomain,
		Tenant:          parsed.tenant,
		Namespace:       parsed.namespace,
		Name:            parsed.topic,
		PartitionIndex:  topicName.GetPartitionIndex(),
		Properties:      sanitizePulsarResourceStringMap(properties, pulsarResourceSummaryStringLimit),
		PropertiesCount: len(properties),
	})
}

func handlePulsarTopicStatsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindTopicStats {
		return nil, fmt.Errorf("unsupported Pulsar topic stats resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}
	topicName, err := buildPulsarResourceTopicName(parsed)
	if err != nil {
		return nil, err
	}

	metadata, err := adminClient.Topics().GetMetadata(*topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to get partition metadata for topic %q: %w", topicName.String(), err)
	}

	resource := pulsarTopicStatsResource{
		Kind:           string(parsed.kind),
		URI:            request.Params.URI,
		Topic:          topicName.String(),
		Partitioned:    metadata.Partitions > 0,
		PartitionCount: metadata.Partitions,
	}

	statsOptions := utils.GetStatsOptions{
		ExcludePublishers: true,
		ExcludeConsumers:  true,
	}
	if metadata.Partitions > 0 {
		stats, err := adminClient.Topics().GetPartitionedStatsWithOption(*topicName, false, statsOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to get partitioned stats for topic %q: %w", topicName.String(), err)
		}
		resource.PartitionStatsCount = len(stats.Partitions)
		resource.Stats = summarizePulsarPartitionedTopicStats(stats)
	} else {
		stats, err := adminClient.Topics().GetStatsWithOption(*topicName, statsOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to get stats for topic %q: %w", topicName.String(), err)
		}
		resource.Stats = summarizePulsarTopicStats(stats)
	}

	return newPulsarJSONResourceContents(request.Params.URI, resource)
}

func handlePulsarTopicPartitionMetadataResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindTopicPartitions {
		return nil, fmt.Errorf("unsupported Pulsar topic partition metadata resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}
	topicName, err := buildPulsarResourceTopicName(parsed)
	if err != nil {
		return nil, err
	}

	metadata, err := adminClient.Topics().GetMetadata(*topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to get partition metadata for topic %q: %w", topicName.String(), err)
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarTopicPartitionMetadataResource{
		Kind:        string(parsed.kind),
		URI:         request.Params.URI,
		Topic:       topicName.String(),
		Metadata:    metadata,
		Partitioned: metadata.Partitions > 0,
	})
}

func handlePulsarTopicPolicyResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindTopicPolicy {
		return nil, fmt.Errorf("unsupported Pulsar topic policy resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}
	topicName, err := buildPulsarResourceTopicName(parsed)
	if err != nil {
		return nil, err
	}

	value, err := readPulsarTopicPolicyValue(adminClient, *topicName, parsed.policy)
	if err != nil {
		return nil, err
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarTopicPolicyResource{
		Kind:   string(parsed.kind),
		URI:    request.Params.URI,
		Topic:  topicName.String(),
		Policy: parsed.policy,
		Value:  value,
	})
}

func handlePulsarTopicSchemaResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindTopicSchema {
		return nil, fmt.Errorf("unsupported Pulsar topic schema resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}
	topicName, err := buildPulsarResourceTopicName(parsed)
	if err != nil {
		return nil, err
	}

	schemaInfo, err := adminClient.Schemas().GetSchemaInfoWithVersion(topicName.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get latest schema for topic %q: %w", topicName.String(), err)
	}
	if schemaInfo == nil || schemaInfo.SchemaInfo == nil {
		return nil, fmt.Errorf("latest schema for topic %q is empty", topicName.String())
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarTopicSchemaResource{
		Kind:    string(parsed.kind),
		URI:     request.Params.URI,
		Topic:   topicName.String(),
		Version: schemaInfo.Version,
		Schema:  summarizePulsarTopicSchema(schemaInfo.SchemaInfo),
	})
}

func handlePulsarTopicSchemaVersionResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindTopicSchemaVersion {
		return nil, fmt.Errorf("unsupported Pulsar topic schema version resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}
	topicName, err := buildPulsarResourceTopicName(parsed)
	if err != nil {
		return nil, err
	}

	schemaInfo, err := adminClient.Schemas().GetSchemaInfoByVersion(topicName.String(), parsed.version)
	if err != nil {
		return nil, fmt.Errorf("failed to get schema version %d for topic %q: %w", parsed.version, topicName.String(), err)
	}
	if schemaInfo == nil {
		return nil, fmt.Errorf("schema version %d for topic %q is empty", parsed.version, topicName.String())
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarTopicSchemaResource{
		Kind:    string(parsed.kind),
		URI:     request.Params.URI,
		Topic:   topicName.String(),
		Version: parsed.version,
		Schema:  summarizePulsarTopicSchema(schemaInfo),
	})
}

func handlePulsarSubscriptionsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindSubscriptions {
		return nil, fmt.Errorf("unsupported Pulsar subscriptions resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}
	topicName, err := buildPulsarResourceTopicName(parsed)
	if err != nil {
		return nil, err
	}

	subscriptions, err := adminClient.Subscriptions().List(*topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions for topic %q: %w", topicName.String(), err)
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarSubscriptionCollectionResource{
		Kind:          string(parsed.kind),
		URI:           request.Params.URI,
		Topic:         topicName.String(),
		Domain:        parsed.topicDomain,
		Tenant:        parsed.tenant,
		Namespace:     parsed.namespace,
		Subscriptions: subscriptions,
		Count:         len(subscriptions),
	})
}

func handlePulsarSubscriptionStatsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindSubscriptionStats {
		return nil, fmt.Errorf("unsupported Pulsar subscription stats resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}
	topicName, err := buildPulsarResourceTopicName(parsed)
	if err != nil {
		return nil, err
	}

	stats, partitioned, err := readPulsarSubscriptionStats(adminClient, *topicName, parsed.subscription, utils.GetStatsOptions{
		ExcludePublishers: true,
		ExcludeConsumers:  true,
	})
	if err != nil {
		return nil, err
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarSubscriptionStatsResource{
		Kind:         string(parsed.kind),
		URI:          request.Params.URI,
		Topic:        topicName.String(),
		Subscription: parsed.subscription,
		Partitioned:  partitioned,
		Stats:        summarizePulsarSubscriptionStats(stats),
	})
}

func handlePulsarSubscriptionBacklogResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindSubscriptionBacklog {
		return nil, fmt.Errorf("unsupported Pulsar subscription backlog resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}
	topicName, err := buildPulsarResourceTopicName(parsed)
	if err != nil {
		return nil, err
	}

	stats, partitioned, err := readPulsarSubscriptionStats(adminClient, *topicName, parsed.subscription, utils.GetStatsOptions{
		SubscriptionBacklogSize:  true,
		GetEarliestTimeInBacklog: true,
		ExcludePublishers:        true,
		ExcludeConsumers:         true,
	})
	if err != nil {
		return nil, err
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarSubscriptionBacklogResource{
		Kind:         string(parsed.kind),
		URI:          request.Params.URI,
		Topic:        topicName.String(),
		Subscription: parsed.subscription,
		Partitioned:  partitioned,
		Backlog:      summarizePulsarSubscriptionBacklog(stats),
	})
}

func handlePulsarSubscriptionCursorResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindSubscriptionCursor {
		return nil, fmt.Errorf("unsupported Pulsar subscription cursor resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}
	topicName, err := buildPulsarResourceTopicName(parsed)
	if err != nil {
		return nil, err
	}

	internalStats, err := adminClient.Topics().GetInternalStats(*topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to get internal stats for topic %q: %w", topicName.String(), err)
	}
	cursor, ok := internalStats.Cursors[parsed.subscription]
	if !ok {
		return nil, fmt.Errorf("subscription %q was not found in cursor stats for topic %q", parsed.subscription, topicName.String())
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarSubscriptionCursorResource{
		Kind:         string(parsed.kind),
		URI:          request.Params.URI,
		Topic:        topicName.String(),
		Subscription: parsed.subscription,
		Cursor:       summarizePulsarSubscriptionCursor(cursor),
	})
}

func requirePulsarResourceSession(ctx context.Context) (*pulsarsession.Session, error) {
	session := context2.GetPulsarSession(ctx)
	if session == nil {
		return nil, fmt.Errorf("Pulsar session not found in context")
	}
	return session, nil
}

func getPulsarResourceAdminClient(session *pulsarsession.Session) (cmdutils.Client, error) {
	if _, err := session.GetPulsarCtlConfig(); err != nil {
		return nil, fmt.Errorf("failed to get Pulsar admin configuration: %w", err)
	}
	adminClient, err := session.GetAdminClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get Pulsar admin client: %w", err)
	}
	if adminClient == nil {
		return nil, fmt.Errorf("Pulsar admin client not found in session")
	}
	return adminClient, nil
}

func getPulsarClusterStatus(session *pulsarsession.Session) (string, error) {
	cfg, err := session.GetPulsarCtlConfig()
	if err != nil {
		return "", fmt.Errorf("failed to get Pulsar admin configuration: %w", err)
	}

	authProvider, err := pulsaradminauth.GetAuthProvider((*pulsaradminconfig.Config)(cfg))
	if err != nil {
		return "", fmt.Errorf("failed to build status auth provider: %w", err)
	}

	statusClient := &rest.Client{
		ServiceURL:  cfg.WebServiceURL,
		VersionInfo: admin.ReleaseVersion,
		HTTPClient: &http.Client{
			Timeout:   admin.DefaultHTTPTimeOutDuration,
			Transport: authProvider,
		},
	}
	data, err := statusClient.GetWithQueryParams("/status.html", nil, nil, false)
	if err != nil {
		return "", fmt.Errorf("failed to check Pulsar status: %w", err)
	}

	status := strings.TrimSpace(string(data))
	if status == "" {
		status = string(data)
	}
	return status, nil
}

func sanitizePulsarClusterData(clusterName string, data utils.ClusterData) pulsarClusterDataSummary {
	name := data.Name
	if name == "" {
		name = clusterName
	}
	return pulsarClusterDataSummary{
		Name:                                 name,
		ServiceURL:                           data.ServiceURL,
		ServiceURLTLS:                        data.ServiceURLTls,
		BrokerServiceURL:                     data.BrokerServiceURL,
		BrokerServiceURLTLS:                  data.BrokerServiceURLTls,
		PeerClusterNames:                     data.PeerClusterNames,
		AuthenticationPlugin:                 data.AuthenticationPlugin,
		AuthenticationParametersConfigured:   data.AuthenticationParameters != "",
		BrokerClientTrustCertsFileConfigured: data.BrokerClientTrustCertsFilePath != "",
		BrokerClientTLSEnabled:               data.BrokerClientTLSEnabled,
	}
}

func summarizePulsarMonitoringMetrics(metrics []utils.Metrics) pulsarMonitoringMetricsSummary {
	metricNames := make(map[string]struct{})
	dimensionKeys := make(map[string]struct{})
	for _, metric := range metrics {
		for name := range metric.Metrics {
			metricNames[name] = struct{}{}
		}
		for name := range metric.Dimensions {
			dimensionKeys[name] = struct{}{}
		}
	}

	return pulsarMonitoringMetricsSummary{
		Count:         len(metrics),
		MetricNames:   sortedLimitedStrings(metricNames, pulsarResourceSummaryStringLimit),
		DimensionKeys: sortedLimitedStrings(dimensionKeys, pulsarResourceSummaryStringLimit),
	}
}

func summarizePulsarBrokerLoadReport(loadReport *utils.LocalBrokerData) pulsarBrokerLoadReportSummary {
	if loadReport == nil {
		return pulsarBrokerLoadReportSummary{Available: false}
	}
	return pulsarBrokerLoadReportSummary{
		Available:                  true,
		WebServiceURL:              loadReport.WebServiceURL,
		WebServiceURLTLS:           loadReport.WebServiceURLTLS,
		PulsarServiceURL:           loadReport.PulsarServiceURL,
		PulsarServiceURLTLS:        loadReport.PulsarServiceURLTLS,
		PersistentTopicsEnabled:    loadReport.PersistentTopicsEnabled,
		NonPersistentTopicsEnabled: loadReport.NonPersistentTopicsEnabled,
		CPU:                        summarizePulsarResourceUsage(loadReport.CPU),
		Memory:                     summarizePulsarResourceUsage(loadReport.Memory),
		DirectMemory:               summarizePulsarResourceUsage(loadReport.DirectMemory),
		BandwidthIn:                summarizePulsarResourceUsage(loadReport.BandwidthIn),
		BandwidthOut:               summarizePulsarResourceUsage(loadReport.BandwidthOut),
		MsgThroughputIn:            loadReport.MsgThroughputIn,
		MsgThroughputOut:           loadReport.MsgThroughputOut,
		MsgRateIn:                  loadReport.MsgRateIn,
		MsgRateOut:                 loadReport.MsgRateOut,
		LastUpdate:                 loadReport.LastUpdate,
		NumTopics:                  loadReport.NumTopics,
		NumBundles:                 loadReport.NumBundles,
		NumConsumers:               loadReport.NumConsumers,
		NumProducers:               loadReport.NumProducers,
		BundleCount:                len(loadReport.Bundles),
		LastBundleGainsCount:       len(loadReport.LastBundleGains),
		LastBundleLossesCount:      len(loadReport.LastBundleLosses),
		BrokerVersionString:        loadReport.BrokerVersionString,
		LoadReportType:             loadReport.LoadReportType,
		ProtocolCount:              len(loadReport.Protocols),
	}
}

func summarizePulsarResourceUsage(usage utils.ResourceUsage) pulsarResourceUsageSummary {
	return pulsarResourceUsageSummary{
		Usage:        usage.Usage,
		Limit:        usage.Limit,
		PercentUsage: usage.PercentUsage(),
	}
}

func summarizePulsarFailureDomains(failureDomains utils.FailureDomainMap) []pulsarFailureDomainSummary {
	names := make([]string, 0, len(failureDomains))
	for name := range failureDomains {
		names = append(names, name)
	}
	sort.Strings(names)

	summaries := make([]pulsarFailureDomainSummary, 0, len(names))
	for _, name := range names {
		summaries = append(summaries, summarizePulsarFailureDomain(name, failureDomains[name]))
	}
	return summaries
}

func summarizePulsarFailureDomain(name string, failureDomain utils.FailureDomainData) pulsarFailureDomainSummary {
	return pulsarFailureDomainSummary{
		Name:    name,
		Brokers: failureDomain.BrokerList,
	}
}

func summarizePulsarNamespaceIsolationPolicies(
	policies map[string]utils.NamespaceIsolationData,
) []pulsarNamespaceIsolationPolicySummary {
	names := make([]string, 0, len(policies))
	for name := range policies {
		names = append(names, name)
	}
	sort.Strings(names)

	summaries := make([]pulsarNamespaceIsolationPolicySummary, 0, len(names))
	for _, name := range names {
		summaries = append(summaries, summarizePulsarNamespaceIsolationPolicy(name, policies[name]))
	}
	return summaries
}

func summarizePulsarNamespaceIsolationPolicy(
	name string,
	policy utils.NamespaceIsolationData,
) pulsarNamespaceIsolationPolicySummary {
	return pulsarNamespaceIsolationPolicySummary{
		Name:                   name,
		NamespacesCount:        len(policy.Namespaces),
		PrimaryBrokersCount:    len(policy.Primary),
		SecondaryBrokersCount:  len(policy.Secondary),
		AutoFailoverPolicyType: string(policy.AutoFailoverPolicy.PolicyType),
	}
}

func summarizePulsarTopicStats(stats utils.TopicStats) pulsarTopicStatsSummary {
	return pulsarTopicStatsSummary{
		BacklogSize:          stats.BacklogSize,
		MsgCounterIn:         stats.MsgCounterIn,
		MsgCounterOut:        stats.MsgCounterOut,
		MsgRateIn:            stats.MsgRateIn,
		MsgRateOut:           stats.MsgRateOut,
		MsgThroughputIn:      stats.MsgThroughputIn,
		MsgThroughputOut:     stats.MsgThroughputOut,
		AverageMsgSize:       stats.AverageMsgSize,
		StorageSize:          stats.StorageSize,
		PublisherCount:       len(stats.Publishers),
		SubscriptionCount:    len(stats.Subscriptions),
		ReplicationCount:     len(stats.Replication),
		DeduplicationStatus:  stats.DeDuplicationStatus,
		TopicCreationTime:    stats.TopicCreationTimeStamp,
		LastPublishTimestamp: stats.LastPublishTimestamp,
	}
}

func summarizePulsarPartitionedTopicStats(stats utils.PartitionedTopicStats) pulsarTopicStatsSummary {
	return pulsarTopicStatsSummary{
		MsgRateIn:            stats.MsgRateIn,
		MsgRateOut:           stats.MsgRateOut,
		MsgThroughputIn:      stats.MsgThroughputIn,
		MsgThroughputOut:     stats.MsgThroughputOut,
		AverageMsgSize:       stats.AverageMsgSize,
		StorageSize:          stats.StorageSize,
		PublisherCount:       len(stats.Publishers),
		SubscriptionCount:    len(stats.Subscriptions),
		ReplicationCount:     len(stats.Replication),
		DeduplicationStatus:  stats.DeDuplicationStatus,
		TopicCreationTime:    stats.TopicCreationTimeStamp,
		LastPublishTimestamp: stats.LastPublishTimestamp,
	}
}

func readPulsarSubscriptionStats(
	adminClient cmdutils.Client,
	topicName utils.TopicName,
	subscription string,
	options utils.GetStatsOptions,
) (utils.SubscriptionStats, bool, error) {
	metadata, err := adminClient.Topics().GetMetadata(topicName)
	if err != nil {
		return utils.SubscriptionStats{}, false, fmt.Errorf("failed to get partition metadata for topic %q: %w", topicName.String(), err)
	}

	var subscriptions map[string]utils.SubscriptionStats
	partitioned := metadata.Partitions > 0
	if partitioned {
		stats, err := adminClient.Topics().GetPartitionedStatsWithOption(topicName, false, options)
		if err != nil {
			return utils.SubscriptionStats{}, false, fmt.Errorf("failed to get partitioned stats for topic %q: %w", topicName.String(), err)
		}
		subscriptions = stats.Subscriptions
	} else {
		stats, err := adminClient.Topics().GetStatsWithOption(topicName, options)
		if err != nil {
			return utils.SubscriptionStats{}, false, fmt.Errorf("failed to get stats for topic %q: %w", topicName.String(), err)
		}
		subscriptions = stats.Subscriptions
	}

	stats, ok := subscriptions[subscription]
	if !ok {
		return utils.SubscriptionStats{}, partitioned, fmt.Errorf("subscription %q was not found in stats for topic %q", subscription, topicName.String())
	}
	return stats, partitioned, nil
}

func summarizePulsarSubscriptionStats(stats utils.SubscriptionStats) pulsarSubscriptionStatsSummary {
	return pulsarSubscriptionStatsSummary{
		Type:                               stats.SubType,
		Durable:                            stats.IsDurable,
		Replicated:                         stats.IsReplicated,
		BlockedOnUnackedMessages:           stats.BlockedSubscriptionOnUnackedMsgs,
		MsgRateOut:                         stats.MsgRateOut,
		MsgThroughputOut:                   stats.MsgThroughputOut,
		MsgRateRedeliver:                   stats.MsgRateRedeliver,
		MsgRateExpired:                     stats.MsgRateExpired,
		MsgBacklog:                         stats.MsgBacklog,
		MsgBacklogNoDelayed:                stats.MsgBacklogNoDelayed,
		MsgDelayed:                         stats.MsgDelayed,
		UnackedMessages:                    stats.UnAckedMessages,
		BytesOutCounter:                    stats.BytesOutCounter,
		MsgOutCounter:                      stats.MsgOutCounter,
		MessageAckRate:                     stats.MessageAckRate,
		ChunkedMessageRate:                 stats.ChunkedMessageRate,
		BacklogSize:                        stats.BacklogSize,
		EarliestMsgPublishTimeInBacklog:    stats.EarliestMsgPublishTimeInBacklog,
		TotalMsgExpired:                    stats.TotalMsgExpired,
		LastExpireTimestamp:                stats.LastExpireTimestamp,
		LastConsumedFlowTimestamp:          stats.LastConsumedFlowTimestamp,
		LastConsumedTimestamp:              stats.LastConsumedTimestamp,
		LastAckedTimestamp:                 stats.LastAckedTimestamp,
		LastMarkDeleteAdvancedTimestamp:    stats.LastMarkDeleteAdvancedTimestamp,
		AllowOutOfOrderDelivery:            stats.AllowOutOfOrderDelivery,
		NonContiguousDeletedMessagesRanges: stats.NonContiguousDeletedMessagesRanges,
		NonContiguousDeletedMessagesRangesSrzSize: stats.NonContiguousDeletedMessagesRangesSrzSize,
		DelayedMessageIndexSizeInBytes:            stats.DelayedMessageIndexSizeInBytes,
		FilterProcessedMsgCount:                   stats.FilterProcessedMsgCount,
		FilterAcceptedMsgCount:                    stats.FilterAcceptedMsgCount,
		FilterRejectedMsgCount:                    stats.FilterRejectedMsgCount,
		FilterRescheduledMsgCount:                 stats.FilterRescheduledMsgCount,
		SubscriptionProperties: sanitizePulsarResourceStringMap(
			stats.SubscriptionProperties,
			pulsarResourceSummaryStringLimit,
		),
		SubscriptionPropertiesCount: len(stats.SubscriptionProperties),
	}
}

func summarizePulsarSubscriptionBacklog(stats utils.SubscriptionStats) pulsarSubscriptionBacklogSummary {
	return pulsarSubscriptionBacklogSummary{
		MsgBacklog:                      stats.MsgBacklog,
		MsgBacklogNoDelayed:             stats.MsgBacklogNoDelayed,
		BacklogSize:                     stats.BacklogSize,
		MsgDelayed:                      stats.MsgDelayed,
		UnackedMessages:                 stats.UnAckedMessages,
		EarliestMsgPublishTimeInBacklog: stats.EarliestMsgPublishTimeInBacklog,
		DelayedMessageIndexSizeInBytes:  stats.DelayedMessageIndexSizeInBytes,
	}
}

func summarizePulsarSubscriptionCursor(stats utils.CursorStats) pulsarSubscriptionCursorSummary {
	return pulsarSubscriptionCursorSummary{
		MarkDeletePosition:                       stats.MarkDeletePosition,
		ReadPosition:                             stats.ReadPosition,
		WaitingReadOp:                            stats.WaitingReadOp,
		PendingReadOps:                           stats.PendingReadOps,
		MessagesConsumedCounter:                  stats.MessagesConsumedCounter,
		CursorLedger:                             stats.CursorLedger,
		CursorLedgerLastEntry:                    stats.CursorLedgerLastEntry,
		LastLedgerSwitchTimestamp:                stats.LastLedgerWitchTimestamp,
		State:                                    stats.State,
		NumberOfEntriesSinceFirstNotAckedMessage: stats.NumberOfEntriesSinceFirstNotAckedMessage,
		TotalNonContiguousDeletedMessagesRange:   stats.TotalNonContiguousDeletedMessagesRange,
		Properties:                               sanitizePulsarResourceInt64Map(stats.Properties, pulsarResourceSummaryStringLimit),
		PropertiesCount:                          len(stats.Properties),
	}
}

func summarizePulsarTopicSchema(schemaInfo *utils.SchemaInfo) pulsarTopicSchemaSummary {
	if schemaInfo == nil {
		return pulsarTopicSchemaSummary{}
	}
	return pulsarTopicSchemaSummary{
		Name:            schemaInfo.Name,
		Type:            schemaInfo.Type,
		Schema:          string(schemaInfo.Schema),
		Properties:      sanitizePulsarResourceStringMap(schemaInfo.Properties, pulsarResourceSummaryStringLimit),
		PropertiesCount: len(schemaInfo.Properties),
		Timestamp:       schemaInfo.Timestamp,
	}
}

func readPulsarTopicPolicyValue(adminClient cmdutils.Client, topicName utils.TopicName, policy string) (any, error) {
	switch policy {
	case "retention":
		value, err := adminClient.Topics().GetRetention(topicName, false)
		return value, wrapPulsarTopicPolicyReadError(topicName, policy, err)
	case "message-ttl":
		value, err := adminClient.Topics().GetMessageTTL(topicName)
		return value, wrapPulsarTopicPolicyReadError(topicName, policy, err)
	case "max-producers":
		value, err := adminClient.Topics().GetMaxProducers(topicName)
		return value, wrapPulsarTopicPolicyReadError(topicName, policy, err)
	case "max-consumers":
		value, err := adminClient.Topics().GetMaxConsumers(topicName)
		return value, wrapPulsarTopicPolicyReadError(topicName, policy, err)
	case "max-unacked-messages-per-consumer":
		value, err := adminClient.Topics().GetMaxUnackMessagesPerConsumer(topicName)
		return value, wrapPulsarTopicPolicyReadError(topicName, policy, err)
	case "max-unacked-messages-per-subscription":
		value, err := adminClient.Topics().GetMaxUnackMessagesPerSubscription(topicName)
		return value, wrapPulsarTopicPolicyReadError(topicName, policy, err)
	case "persistence":
		value, err := adminClient.Topics().GetPersistence(topicName)
		return value, wrapPulsarTopicPolicyReadError(topicName, policy, err)
	case "delayed-delivery":
		value, err := adminClient.Topics().GetDelayedDelivery(topicName)
		return value, wrapPulsarTopicPolicyReadError(topicName, policy, err)
	case "dispatch-rate":
		value, err := adminClient.Topics().GetDispatchRate(topicName)
		return value, wrapPulsarTopicPolicyReadError(topicName, policy, err)
	case "subscription-dispatch-rate":
		value, err := adminClient.Topics().GetSubscriptionDispatchRate(topicName)
		return value, wrapPulsarTopicPolicyReadError(topicName, policy, err)
	case "deduplication":
		value, err := adminClient.Topics().GetDeduplicationStatus(topicName)
		return value, wrapPulsarTopicPolicyReadError(topicName, policy, err)
	case "backlog-quotas":
		value, err := adminClient.Topics().GetBacklogQuotaMap(topicName, false)
		return value, wrapPulsarTopicPolicyReadError(topicName, policy, err)
	case "compaction-threshold":
		value, err := adminClient.Topics().GetCompactionThreshold(topicName, false)
		return value, wrapPulsarTopicPolicyReadError(topicName, policy, err)
	case "publish-rate":
		value, err := adminClient.Topics().GetPublishRate(topicName)
		return value, wrapPulsarTopicPolicyReadError(topicName, policy, err)
	case "inactive-topic-policies":
		value, err := adminClient.Topics().GetInactiveTopicPolicies(topicName, false)
		return value, wrapPulsarTopicPolicyReadError(topicName, policy, err)
	default:
		return nil, fmt.Errorf("unsupported Pulsar topic policy %q", policy)
	}
}

func wrapPulsarTopicPolicyReadError(topicName utils.TopicName, policy string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("failed to get topic policy %q for topic %q: %w", policy, topicName.String(), err)
}

func buildPulsarResourceTopicName(parsed pulsarResourceURI) (*utils.TopicName, error) {
	topicName := fmt.Sprintf("%s://%s/%s/%s", parsed.topicDomain, parsed.tenant, parsed.namespace, parsed.topic)
	return utils.GetTopicName(topicName)
}

func sanitizePulsarResourceStringMap(values map[string]string, limit int) map[string]string {
	if len(values) == 0 || limit <= 0 {
		return nil
	}

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	if len(keys) > limit {
		keys = keys[:limit]
	}

	sanitized := make(map[string]string, len(keys))
	for _, key := range keys {
		if isSensitivePulsarResourceKey(key) {
			sanitized[key] = pulsarResourceRedactedValue
			continue
		}
		sanitized[key] = values[key]
	}
	return sanitized
}

func sanitizePulsarResourceInt64Map(values map[string]int64, limit int) map[string]int64 {
	if len(values) == 0 || limit <= 0 {
		return nil
	}

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	if len(keys) > limit {
		keys = keys[:limit]
	}

	sanitized := make(map[string]int64, len(keys))
	for _, key := range keys {
		if isSensitivePulsarResourceKey(key) {
			continue
		}
		sanitized[key] = values[key]
	}
	if len(sanitized) == 0 {
		return nil
	}
	return sanitized
}

func isSensitivePulsarResourceKey(key string) bool {
	normalized := strings.NewReplacer("-", "", "_", "", ".", "").Replace(strings.ToLower(key))
	sensitiveFragments := []string{
		"token",
		"secret",
		"password",
		"passwd",
		"credential",
		"privatekey",
		"clientkey",
		"tlskey",
		"authparams",
		"authenticationparameters",
		"keyfile",
	}
	for _, fragment := range sensitiveFragments {
		if strings.Contains(normalized, fragment) {
			return true
		}
	}
	return false
}

func sortedLimitedStrings(values map[string]struct{}, limit int) []string {
	if len(values) == 0 || limit <= 0 {
		return nil
	}
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	sort.Strings(result)
	if len(result) > limit {
		return result[:limit]
	}
	return result
}

func buildPulsarContextResource(uri string, session *pulsarsession.Session) (pulsarContextResource, error) {
	ctx := session.Ctx
	if ctx.WebServiceURL == "" {
		if cfg, err := session.GetPulsarCtlConfig(); err == nil && cfg != nil {
			ctx.WebServiceURL = cfg.WebServiceURL
			ctx.TLSAllowInsecureConnection = cfg.TLSAllowInsecureConnection
			ctx.TLSEnableHostnameVerification = cfg.TLSEnableHostnameVerification
			ctx.TLSTrustCertsFilePath = cfg.TLSTrustCertsFilePath
			ctx.TLSCertFile = cfg.TLSCertFile
			ctx.TLSKeyFile = cfg.TLSKeyFile
		}
	}
	if ctx.ServiceURL == "" && ctx.WebServiceURL == "" {
		return pulsarContextResource{}, fmt.Errorf("Pulsar session is not configured")
	}

	return pulsarContextResource{
		Kind:           string(pulsarResourceKindContext),
		URI:            uri,
		ServiceURL:     ctx.ServiceURL,
		WebServiceURL:  ctx.WebServiceURL,
		Authentication: summarizePulsarAuthentication(ctx),
		TLS: pulsarTLSSummary{
			AllowInsecureConnection:    ctx.TLSAllowInsecureConnection,
			EnableHostnameVerification: ctx.TLSEnableHostnameVerification,
			TrustCertsFileConfigured:   ctx.TLSTrustCertsFilePath != "",
			ClientCertFileConfigured:   ctx.TLSCertFile != "",
			ClientKeyFileConfigured:    ctx.TLSKeyFile != "",
		},
	}, nil
}

func summarizePulsarAuthentication(ctx pulsarsession.PulsarContext) pulsarAuthenticationSummary {
	switch {
	case ctx.Token != "":
		return pulsarAuthenticationSummary{
			Configured: true,
			Method:     "token",
		}
	case ctx.AuthPlugin != "" || ctx.AuthParams != "":
		return pulsarAuthenticationSummary{
			Configured: true,
			Method:     "authPlugin",
			Plugin:     ctx.AuthPlugin,
		}
	default:
		return pulsarAuthenticationSummary{
			Configured: false,
			Method:     "none",
		}
	}
}

func buildPulsarResourceCatalog(
	resources []server.ServerResource,
	templates []server.ServerResourceTemplate,
) pulsarResourceCatalog {
	catalog := pulsarResourceCatalog{
		Version: 1,
		Scheme:  "pulsar",
		Resources: []pulsarCatalogResource{
			{
				URI:         pulsarResourceContextURI,
				Name:        "Pulsar Context",
				Description: "Current Pulsar session connection metadata with authentication material redacted.",
				MIMEType:    pulsarResourceJSONMIMEType,
			},
			{
				URI:         pulsarResourceCatalogURI,
				Name:        "Pulsar Resource Catalog",
				Description: "Stable catalog of Pulsar MCP resource URIs and URI templates.",
				MIMEType:    pulsarResourceJSONMIMEType,
			},
		},
		Templates: []pulsarCatalogTemplate{},
		Notes: []string{
			"Resource handlers are read-only and do not return tokens, auth params, key files, TLS private keys, or secret values.",
		},
	}

	for _, resource := range resources {
		catalog.Resources = append(catalog.Resources, pulsarCatalogResource{
			URI:         resource.Resource.URI,
			Name:        resource.Resource.Name,
			Description: resource.Resource.Description,
			MIMEType:    resource.Resource.MIMEType,
		})
	}
	for _, template := range templates {
		catalog.Templates = append(catalog.Templates, pulsarCatalogTemplate{
			URITemplate: template.Template.URITemplate.Raw(),
			Name:        template.Template.Name,
			Description: template.Template.Description,
			MIMEType:    template.Template.MIMEType,
		})
	}

	return catalog
}

func parsePulsarResourceURI(rawURI string) (pulsarResourceURI, error) {
	parsed, err := url.Parse(rawURI)
	if err != nil {
		return pulsarResourceURI{}, fmt.Errorf("malformed Pulsar resource URI %q: %w", rawURI, err)
	}
	if parsed.Scheme != "pulsar" {
		return pulsarResourceURI{}, fmt.Errorf("unsupported Pulsar resource URI scheme %q", parsed.Scheme)
	}
	if parsed.User != nil {
		return pulsarResourceURI{}, fmt.Errorf("Pulsar resource URI %q must not include user info", rawURI)
	}
	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return pulsarResourceURI{}, fmt.Errorf("Pulsar resource URI %q must not include query or fragment", rawURI)
	}

	switch parsed.Host {
	case "context":
		if parsed.Path != "" {
			return pulsarResourceURI{}, fmt.Errorf("unsupported Pulsar context resource URI %q", rawURI)
		}
		return pulsarResourceURI{kind: pulsarResourceKindContext}, nil
	case "resources":
		if parsed.Path != "" {
			return pulsarResourceURI{}, fmt.Errorf("unsupported Pulsar resource catalog URI %q", rawURI)
		}
		return pulsarResourceURI{kind: pulsarResourceKindCatalog}, nil
	case "admin":
		return parsePulsarAdminResourceURI(rawURI, parsed.Path)
	default:
		return pulsarResourceURI{}, fmt.Errorf("unsupported Pulsar resource URI host %q", parsed.Host)
	}
}

func parsePulsarAdminResourceURI(rawURI, path string) (pulsarResourceURI, error) {
	parts := splitPulsarResourcePath(path)
	switch {
	case len(parts) == 2 && parts[0] == "v2" && parts[1] == "tenants":
		return pulsarResourceURI{kind: pulsarResourceKindTenants}, nil
	case len(parts) == 3 && parts[0] == "v2" && parts[1] == "tenants":
		tenant := parts[2]
		if err := validatePulsarResourcePathSegment("tenant", tenant); err != nil {
			return pulsarResourceURI{}, err
		}
		return pulsarResourceURI{
			kind:   pulsarResourceKindTenant,
			tenant: tenant,
		}, nil
	case len(parts) == 4 && parts[0] == "v2" && parts[1] == "tenants" && parts[3] == "namespaces":
		tenant := parts[2]
		if err := validatePulsarResourcePathSegment("tenant", tenant); err != nil {
			return pulsarResourceURI{}, err
		}
		return pulsarResourceURI{
			kind:   pulsarResourceKindNamespaces,
			tenant: tenant,
		}, nil
	case len(parts) == 4 && parts[0] == "v2" && parts[1] == "namespaces":
		tenant := parts[2]
		namespace := parts[3]
		if err := validatePulsarResourcePathSegment("tenant", tenant); err != nil {
			return pulsarResourceURI{}, err
		}
		if err := validatePulsarResourcePathSegment("namespace", namespace); err != nil {
			return pulsarResourceURI{}, err
		}
		return pulsarResourceURI{
			kind:      pulsarResourceKindNamespace,
			tenant:    tenant,
			namespace: namespace,
		}, nil
	case len(parts) == 5 && parts[0] == "v2" && parts[1] == "namespaces" && parts[4] == "topics":
		tenant := parts[2]
		namespace := parts[3]
		if err := validatePulsarResourcePathSegment("tenant", tenant); err != nil {
			return pulsarResourceURI{}, err
		}
		if err := validatePulsarResourcePathSegment("namespace", namespace); err != nil {
			return pulsarResourceURI{}, err
		}
		return pulsarResourceURI{
			kind:      pulsarResourceKindTopics,
			tenant:    tenant,
			namespace: namespace,
		}, nil
	case len(parts) == 2 && parts[0] == "v2" && parts[1] == "resource-quotas":
		return pulsarResourceURI{kind: pulsarResourceKindDefaultResourceQuota}, nil
	case len(parts) == 5 && parts[0] == "v2" && parts[1] == "resource-quotas":
		tenant := parts[2]
		namespace := parts[3]
		bundle := parts[4]
		if err := validatePulsarResourcePathSegment("tenant", tenant); err != nil {
			return pulsarResourceURI{}, err
		}
		if err := validatePulsarResourcePathSegment("namespace", namespace); err != nil {
			return pulsarResourceURI{}, err
		}
		if err := validatePulsarResourcePathSegment("bundle", bundle); err != nil {
			return pulsarResourceURI{}, err
		}
		return pulsarResourceURI{
			kind:      pulsarResourceKindResourceQuota,
			tenant:    tenant,
			namespace: namespace,
			bundle:    bundle,
		}, nil
	case len(parts) == 2 && parts[0] == "v2" && parts[1] == "status":
		return pulsarResourceURI{kind: pulsarResourceKindStatus}, nil
	case len(parts) == 2 && parts[0] == "v2" && parts[1] == "clusters":
		return pulsarResourceURI{kind: pulsarResourceKindClusters}, nil
	case len(parts) == 3 && parts[0] == "v2" && parts[1] == "clusters":
		cluster := parts[2]
		if err := validatePulsarResourcePathSegment("cluster", cluster); err != nil {
			return pulsarResourceURI{}, err
		}
		return pulsarResourceURI{
			kind:    pulsarResourceKindCluster,
			cluster: cluster,
		}, nil
	case len(parts) == 3 && parts[0] == "v2" && parts[1] == "brokers":
		cluster := parts[2]
		if err := validatePulsarResourcePathSegment("cluster", cluster); err != nil {
			return pulsarResourceURI{}, err
		}
		return pulsarResourceURI{
			kind:    pulsarResourceKindBrokers,
			cluster: cluster,
		}, nil
	case len(parts) == 3 && parts[0] == "v2" && parts[1] == "broker-stats" && parts[2] == "summary":
		return pulsarResourceURI{kind: pulsarResourceKindBrokerStatsSummary}, nil
	case len(parts) == 4 && parts[0] == "v2" && parts[1] == "clusters" && parts[3] == "failureDomains":
		cluster := parts[2]
		if err := validatePulsarResourcePathSegment("cluster", cluster); err != nil {
			return pulsarResourceURI{}, err
		}
		return pulsarResourceURI{
			kind:    pulsarResourceKindFailureDomains,
			cluster: cluster,
		}, nil
	case len(parts) == 5 && parts[0] == "v2" && parts[1] == "clusters" && parts[3] == "failureDomains":
		cluster := parts[2]
		domain := parts[4]
		if err := validatePulsarResourcePathSegment("cluster", cluster); err != nil {
			return pulsarResourceURI{}, err
		}
		if err := validatePulsarResourcePathSegment("domain", domain); err != nil {
			return pulsarResourceURI{}, err
		}
		return pulsarResourceURI{
			kind:    pulsarResourceKindFailureDomain,
			cluster: cluster,
			domain:  domain,
		}, nil
	case len(parts) == 4 && parts[0] == "v2" && parts[1] == "clusters" && parts[3] == "namespaceIsolationPolicies":
		cluster := parts[2]
		if err := validatePulsarResourcePathSegment("cluster", cluster); err != nil {
			return pulsarResourceURI{}, err
		}
		return pulsarResourceURI{
			kind:    pulsarResourceKindNSIsolationPolicies,
			cluster: cluster,
		}, nil
	case len(parts) == 5 && parts[0] == "v2" && parts[1] == "clusters" && parts[3] == "namespaceIsolationPolicies":
		cluster := parts[2]
		policy := parts[4]
		if err := validatePulsarResourcePathSegment("cluster", cluster); err != nil {
			return pulsarResourceURI{}, err
		}
		if err := validatePulsarResourcePathSegment("policy", policy); err != nil {
			return pulsarResourceURI{}, err
		}
		return pulsarResourceURI{
			kind:    pulsarResourceKindNSIsolationPolicy,
			cluster: cluster,
			policy:  policy,
		}, nil
	case len(parts) >= 6 && parts[0] == "v2" && isPulsarResourceTopicDomain(parts[1]):
		return parsePulsarTopicAdminResourceURI(rawURI, parts)
	default:
		return pulsarResourceURI{}, fmt.Errorf("unsupported Pulsar admin resource URI %q", rawURI)
	}
}

func parsePulsarTopicAdminResourceURI(rawURI string, parts []string) (pulsarResourceURI, error) {
	topicDomain := parts[1]
	tenant := parts[2]
	namespace := parts[3]
	topic := parts[4]
	if err := validatePulsarResourcePathSegment("domain", topicDomain); err != nil {
		return pulsarResourceURI{}, err
	}
	if err := validatePulsarResourcePathSegment("tenant", tenant); err != nil {
		return pulsarResourceURI{}, err
	}
	if err := validatePulsarResourcePathSegment("namespace", namespace); err != nil {
		return pulsarResourceURI{}, err
	}
	if err := validatePulsarResourcePathSegment("topic", topic); err != nil {
		return pulsarResourceURI{}, err
	}

	base := pulsarResourceURI{
		tenant:      tenant,
		namespace:   namespace,
		topic:       topic,
		topicDomain: topicDomain,
	}
	if _, err := buildPulsarResourceTopicName(base); err != nil {
		return pulsarResourceURI{}, fmt.Errorf("invalid Pulsar topic resource URI %q: %w", rawURI, err)
	}

	switch {
	case len(parts) == 6 && parts[5] == "metadata":
		base.kind = pulsarResourceKindTopicMetadata
		return base, nil
	case len(parts) == 6 && parts[5] == "stats":
		base.kind = pulsarResourceKindTopicStats
		return base, nil
	case len(parts) == 6 && parts[5] == "partitions":
		base.kind = pulsarResourceKindTopicPartitions
		return base, nil
	case len(parts) == 7 && parts[5] == "policies":
		policy := parts[6]
		if err := validatePulsarResourcePathSegment("policy", policy); err != nil {
			return pulsarResourceURI{}, err
		}
		if !isSupportedPulsarTopicPolicy(policy) {
			return pulsarResourceURI{}, fmt.Errorf("unsupported Pulsar topic policy %q in resource URI %q", policy, rawURI)
		}
		base.kind = pulsarResourceKindTopicPolicy
		base.policy = policy
		return base, nil
	case len(parts) == 6 && parts[5] == "schema":
		base.kind = pulsarResourceKindTopicSchema
		return base, nil
	case len(parts) == 7 && parts[5] == "schema":
		versionSegment := parts[6]
		if err := validatePulsarResourcePathSegment("version", versionSegment); err != nil {
			return pulsarResourceURI{}, err
		}
		version, err := strconv.ParseInt(versionSegment, 10, 64)
		if err != nil || version < 0 {
			return pulsarResourceURI{}, fmt.Errorf("invalid Pulsar schema version %q in resource URI %q", versionSegment, rawURI)
		}
		base.kind = pulsarResourceKindTopicSchemaVersion
		base.version = version
		return base, nil
	case len(parts) == 6 && parts[5] == "subscriptions":
		base.kind = pulsarResourceKindSubscriptions
		return base, nil
	case len(parts) == 8 && parts[5] == "subscriptions":
		subscription := parts[6]
		if err := validatePulsarResourcePathSegment("subscription", subscription); err != nil {
			return pulsarResourceURI{}, err
		}
		base.subscription = subscription
		switch parts[7] {
		case "stats":
			base.kind = pulsarResourceKindSubscriptionStats
			return base, nil
		case "backlog":
			base.kind = pulsarResourceKindSubscriptionBacklog
			return base, nil
		case "cursor":
			if topicDomain != "persistent" {
				return pulsarResourceURI{}, fmt.Errorf("Pulsar subscription cursor resource URI %q only supports persistent topics", rawURI)
			}
			base.kind = pulsarResourceKindSubscriptionCursor
			return base, nil
		default:
			return pulsarResourceURI{}, fmt.Errorf("unsupported Pulsar subscription resource URI %q", rawURI)
		}
	default:
		return pulsarResourceURI{}, fmt.Errorf("unsupported Pulsar topic admin resource URI %q", rawURI)
	}
}

func isPulsarResourceTopicDomain(value string) bool {
	return value == "persistent" || value == "non-persistent"
}

func isSupportedPulsarTopicPolicy(policy string) bool {
	switch policy {
	case "retention",
		"message-ttl",
		"max-producers",
		"max-consumers",
		"max-unacked-messages-per-consumer",
		"max-unacked-messages-per-subscription",
		"persistence",
		"delayed-delivery",
		"dispatch-rate",
		"subscription-dispatch-rate",
		"deduplication",
		"backlog-quotas",
		"compaction-threshold",
		"publish-rate",
		"inactive-topic-policies":
		return true
	default:
		return false
	}
}

func splitPulsarResourcePath(path string) []string {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "/")
}

func validatePulsarResourcePathSegment(name, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("Pulsar resource URI %s segment must not be empty", name)
	}
	if strings.Contains(value, "/") {
		return fmt.Errorf("Pulsar resource URI %s segment must not contain path separators", name)
	}
	return nil
}

func newPulsarJSONResourceContents(uri string, value any) ([]mcp.ResourceContents, error) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Pulsar resource %q: %w", uri, err)
	}
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      uri,
			MIMEType: pulsarResourceJSONMIMEType,
			Text:     string(data),
		},
	}, nil
}
