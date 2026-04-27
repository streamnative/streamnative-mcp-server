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
	"net/url"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
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
	pulsarFunctionsTemplateURI           = "pulsar://admin/v3/functions/{tenant}/{namespace}"
	pulsarFunctionMetadataTemplateURI    = "pulsar://admin/v3/functions/{tenant}/{namespace}/{function}/metadata"
	pulsarFunctionStatusTemplateURI      = "pulsar://admin/v3/functions/{tenant}/{namespace}/{function}/status"
	pulsarFunctionStatsTemplateURI       = "pulsar://admin/v3/functions/{tenant}/{namespace}/{function}/stats"
	pulsarSourcesTemplateURI             = "pulsar://admin/v3/sources/{tenant}/{namespace}"
	pulsarSourceMetadataTemplateURI      = "pulsar://admin/v3/sources/{tenant}/{namespace}/{source}/metadata"
	pulsarSourceStatusTemplateURI        = "pulsar://admin/v3/sources/{tenant}/{namespace}/{source}/status"
	pulsarSinksTemplateURI               = "pulsar://admin/v3/sinks/{tenant}/{namespace}"
	pulsarSinkMetadataTemplateURI        = "pulsar://admin/v3/sinks/{tenant}/{namespace}/{sink}/metadata"
	pulsarSinkStatusTemplateURI          = "pulsar://admin/v3/sinks/{tenant}/{namespace}/{sink}/status"
	pulsarPackagesTemplateURI            = "pulsar://admin/v3/packages/{type}/{tenant}/{namespace}"
	pulsarPackageVersionsTemplateURI     = "pulsar://admin/v3/packages/{type}/{tenant}/{namespace}/{package}/versions"
	pulsarPackageMetadataTemplateURI     = "pulsar://admin/v3/packages/{type}/{tenant}/{namespace}/{package}/{version}/metadata"
	pulsarWorkerClusterResourceURI       = "pulsar://admin/v2/worker/cluster"
	pulsarWorkerLeaderResourceURI        = "pulsar://admin/v2/worker/cluster/leader"
	pulsarWorkerAssignmentsResourceURI   = "pulsar://admin/v2/worker/assignments"
	pulsarWorkerFunctionStatsResourceURI = "pulsar://admin/v2/worker-stats/functionsmetrics"
	pulsarWorkerMetricsResourceURI       = "pulsar://admin/v2/worker-stats/metrics"
	pulsarResourceJSONMIMEType           = "application/json"
	pulsarResourceSummaryStringLimit     = 50
	pulsarResourceSanitizeDepthLimit     = 6
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
	pulsarResourceKindFunctions            pulsarResourceKind = "functions"
	pulsarResourceKindFunctionMetadata     pulsarResourceKind = "functionMetadata"
	pulsarResourceKindFunctionStatus       pulsarResourceKind = "functionStatus"
	pulsarResourceKindFunctionStats        pulsarResourceKind = "functionStats"
	pulsarResourceKindSources              pulsarResourceKind = "sources"
	pulsarResourceKindSourceMetadata       pulsarResourceKind = "sourceMetadata"
	pulsarResourceKindSourceStatus         pulsarResourceKind = "sourceStatus"
	pulsarResourceKindSinks                pulsarResourceKind = "sinks"
	pulsarResourceKindSinkMetadata         pulsarResourceKind = "sinkMetadata"
	pulsarResourceKindSinkStatus           pulsarResourceKind = "sinkStatus"
	pulsarResourceKindPackages             pulsarResourceKind = "packages"
	pulsarResourceKindPackageVersions      pulsarResourceKind = "packageVersions"
	pulsarResourceKindPackageMetadata      pulsarResourceKind = "packageMetadata"
	pulsarResourceKindWorkerCluster        pulsarResourceKind = "workerCluster"
	pulsarResourceKindWorkerLeader         pulsarResourceKind = "workerLeader"
	pulsarResourceKindWorkerAssignments    pulsarResourceKind = "workerAssignments"
	pulsarResourceKindWorkerFunctionStats  pulsarResourceKind = "workerFunctionStats"
	pulsarResourceKindWorkerMetrics        pulsarResourceKind = "workerMetrics"
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
	workload     string
	packageType  string
	packageName  string
	versionName  string
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

type pulsarWorkloadCollectionResource struct {
	Kind      string   `json:"kind"`
	URI       string   `json:"uri"`
	Type      string   `json:"type"`
	Tenant    string   `json:"tenant"`
	Namespace string   `json:"namespace"`
	Names     []string `json:"names"`
	Count     int      `json:"count"`
	Limit     int      `json:"limit"`
	Limited   bool     `json:"limited"`
}

type pulsarFunctionMetadataResource struct {
	Kind      string                      `json:"kind"`
	URI       string                      `json:"uri"`
	Tenant    string                      `json:"tenant"`
	Namespace string                      `json:"namespace"`
	Name      string                      `json:"name"`
	Config    pulsarFunctionConfigSummary `json:"config"`
}

type pulsarFunctionConfigSummary struct {
	Tenant                             string                   `json:"tenant,omitempty"`
	Namespace                          string                   `json:"namespace,omitempty"`
	Name                               string                   `json:"name,omitempty"`
	ClassName                          string                   `json:"className,omitempty"`
	Runtime                            string                   `json:"runtime,omitempty"`
	FunctionTypeConfigured             bool                     `json:"functionTypeConfigured"`
	PackageLocationConfigured          bool                     `json:"packageLocationConfigured"`
	PackageLocationType                string                   `json:"packageLocationType,omitempty"`
	Parallelism                        int                      `json:"parallelism,omitempty"`
	Inputs                             []string                 `json:"inputs,omitempty"`
	InputsCount                        int                      `json:"inputsCount"`
	InputSpecsCount                    int                      `json:"inputSpecsCount"`
	TopicsPattern                      string                   `json:"topicsPattern,omitempty"`
	Output                             string                   `json:"output,omitempty"`
	LogTopic                           string                   `json:"logTopic,omitempty"`
	ProcessingGuarantees               string                   `json:"processingGuarantees,omitempty"`
	Resources                          *utils.Resources         `json:"resources,omitempty"`
	UserConfig                         map[string]any           `json:"userConfig,omitempty"`
	UserConfigCount                    int                      `json:"userConfigCount"`
	SecretsCount                       int                      `json:"secretsCount"`
	ProducerConfigConfigured           bool                     `json:"producerConfigConfigured"`
	CustomSchemaOutputsCount           int                      `json:"customSchemaOutputsCount"`
	CustomSerdeInputsCount             int                      `json:"customSerdeInputsCount"`
	CustomSchemaInputsCount            int                      `json:"customSchemaInputsCount"`
	CustomRuntimeOptionsConfigured     bool                     `json:"customRuntimeOptionsConfigured"`
	DeadLetterTopic                    string                   `json:"deadLetterTopic,omitempty"`
	SubscriptionName                   string                   `json:"subscriptionName,omitempty"`
	SubscriptionPosition               string                   `json:"subscriptionPosition,omitempty"`
	TimeoutMs                          *int64                   `json:"timeoutMs,omitempty"`
	MaxMessageRetries                  *int                     `json:"maxMessageRetries,omitempty"`
	CleanupSubscription                bool                     `json:"cleanupSubscription"`
	RetainOrdering                     bool                     `json:"retainOrdering"`
	RetainKeyOrdering                  bool                     `json:"retainKeyOrdering"`
	AutoAck                            bool                     `json:"autoAck"`
	ForwardSourceMessageProperty       bool                     `json:"forwardSourceMessageProperty"`
	ExposePulsarAdminClientEnabled     bool                     `json:"exposePulsarAdminClientEnabled"`
	SkipToLatest                       bool                     `json:"skipToLatest"`
	MaxPendingAsyncRequests            int                      `json:"maxPendingAsyncRequests,omitempty"`
	WindowConfigConfigured             bool                     `json:"windowConfigConfigured"`
	InputTypeClassName                 string                   `json:"inputTypeClassName,omitempty"`
	OutputTypeClassName                string                   `json:"outputTypeClassName,omitempty"`
	OutputSchemaType                   string                   `json:"outputSchemaType,omitempty"`
	OutputSerdeClassName               string                   `json:"outputSerdeClassName,omitempty"`
	CustomSchemaOutputs                map[string]string        `json:"customSchemaOutputs,omitempty"`
	SanitizedCustomSchemaOutputsCount  int                      `json:"sanitizedCustomSchemaOutputsCount"`
	SanitizedCustomSerdeInputs         map[string]string        `json:"customSerdeInputs,omitempty"`
	SanitizedCustomSerdeInputsCount    int                      `json:"sanitizedCustomSerdeInputsCount"`
	SanitizedCustomSchemaInputs        map[string]string        `json:"customSchemaInputs,omitempty"`
	SanitizedCustomSchemaInputsCount   int                      `json:"sanitizedCustomSchemaInputsCount"`
	SanitizedInputSpecsSchemaSummaries []pulsarInputSpecSummary `json:"inputSpecs,omitempty"`
}

type pulsarInputSpecSummary struct {
	Topic                   string            `json:"topic"`
	SchemaType              string            `json:"schemaType,omitempty"`
	SerdeClassName          string            `json:"serdeClassName,omitempty"`
	RegexPattern            bool              `json:"regexPattern"`
	ReceiverQueueSize       int               `json:"receiverQueueSize,omitempty"`
	SchemaProperties        map[string]string `json:"schemaProperties,omitempty"`
	SchemaPropertiesCount   int               `json:"schemaPropertiesCount"`
	ConsumerProperties      map[string]string `json:"consumerProperties,omitempty"`
	ConsumerPropertiesCount int               `json:"consumerPropertiesCount"`
	CryptoConfigConfigured  bool              `json:"cryptoConfigConfigured"`
	PoolMessages            bool              `json:"poolMessages"`
}

type pulsarSourceMetadataResource struct {
	Kind      string                    `json:"kind"`
	URI       string                    `json:"uri"`
	Tenant    string                    `json:"tenant"`
	Namespace string                    `json:"namespace"`
	Name      string                    `json:"name"`
	Config    pulsarSourceConfigSummary `json:"config"`
}

type pulsarSourceConfigSummary struct {
	Tenant                         string           `json:"tenant,omitempty"`
	Namespace                      string           `json:"namespace,omitempty"`
	Name                           string           `json:"name,omitempty"`
	ClassName                      string           `json:"className,omitempty"`
	TopicName                      string           `json:"topicName,omitempty"`
	SerdeClassName                 string           `json:"serdeClassName,omitempty"`
	SchemaType                     string           `json:"schemaType,omitempty"`
	Parallelism                    int              `json:"parallelism,omitempty"`
	ProcessingGuarantees           string           `json:"processingGuarantees,omitempty"`
	Resources                      *utils.Resources `json:"resources,omitempty"`
	ArchiveConfigured              bool             `json:"archiveConfigured"`
	ProducerConfigConfigured       bool             `json:"producerConfigConfigured"`
	RuntimeFlagsConfigured         bool             `json:"runtimeFlagsConfigured"`
	CustomRuntimeOptionsConfigured bool             `json:"customRuntimeOptionsConfigured"`
	BatchSourceConfigConfigured    bool             `json:"batchSourceConfigConfigured"`
	BatchBuilder                   string           `json:"batchBuilder,omitempty"`
	Configs                        map[string]any   `json:"configs,omitempty"`
	ConfigsCount                   int              `json:"configsCount"`
	SecretsCount                   int              `json:"secretsCount"`
}

type pulsarSinkMetadataResource struct {
	Kind      string                  `json:"kind"`
	URI       string                  `json:"uri"`
	Tenant    string                  `json:"tenant"`
	Namespace string                  `json:"namespace"`
	Name      string                  `json:"name"`
	Config    pulsarSinkConfigSummary `json:"config"`
}

type pulsarSinkConfigSummary struct {
	Tenant                         string           `json:"tenant,omitempty"`
	Namespace                      string           `json:"namespace,omitempty"`
	Name                           string           `json:"name,omitempty"`
	ClassName                      string           `json:"className,omitempty"`
	SinkType                       string           `json:"sinkType,omitempty"`
	Parallelism                    int              `json:"parallelism,omitempty"`
	Inputs                         []string         `json:"inputs,omitempty"`
	InputsCount                    int              `json:"inputsCount"`
	InputSpecsCount                int              `json:"inputSpecsCount"`
	TopicsPattern                  string           `json:"topicsPattern,omitempty"`
	ProcessingGuarantees           string           `json:"processingGuarantees,omitempty"`
	SourceSubscriptionName         string           `json:"sourceSubscriptionName,omitempty"`
	SourceSubscriptionPosition     string           `json:"sourceSubscriptionPosition,omitempty"`
	Resources                      *utils.Resources `json:"resources,omitempty"`
	TimeoutMs                      *int64           `json:"timeoutMs,omitempty"`
	ArchiveConfigured              bool             `json:"archiveConfigured"`
	RuntimeFlagsConfigured         bool             `json:"runtimeFlagsConfigured"`
	CustomRuntimeOptionsConfigured bool             `json:"customRuntimeOptionsConfigured"`
	CleanupSubscription            bool             `json:"cleanupSubscription"`
	RetainOrdering                 bool             `json:"retainOrdering"`
	RetainKeyOrdering              bool             `json:"retainKeyOrdering"`
	AutoAck                        bool             `json:"autoAck"`
	TopicToSerdeClassNameCount     int              `json:"topicToSerdeClassNameCount"`
	TopicToSchemaTypeCount         int              `json:"topicToSchemaTypeCount"`
	TopicToSchemaPropertiesCount   int              `json:"topicToSchemaPropertiesCount"`
	Configs                        map[string]any   `json:"configs,omitempty"`
	ConfigsCount                   int              `json:"configsCount"`
	SecretsCount                   int              `json:"secretsCount"`
	MaxMessageRetries              int              `json:"maxMessageRetries,omitempty"`
	DeadLetterTopic                string           `json:"deadLetterTopic,omitempty"`
	NegativeAckRedeliveryDelayMs   int64            `json:"negativeAckRedeliveryDelayMs,omitempty"`
	TransformFunctionConfigured    bool             `json:"transformFunctionConfigured"`
}

type pulsarFunctionStatusResource struct {
	Kind      string                      `json:"kind"`
	URI       string                      `json:"uri"`
	Tenant    string                      `json:"tenant"`
	Namespace string                      `json:"namespace"`
	Name      string                      `json:"name"`
	Status    pulsarFunctionStatusSummary `json:"status"`
}

type pulsarFunctionStatusSummary struct {
	NumInstances int                                   `json:"numInstances"`
	NumRunning   int                                   `json:"numRunning"`
	Instances    []pulsarFunctionInstanceStatusSummary `json:"instances,omitempty"`
	Limit        int                                   `json:"limit"`
	Limited      bool                                  `json:"limited"`
}

type pulsarFunctionInstanceStatusSummary struct {
	InstanceID                  int     `json:"instanceId"`
	Running                     bool    `json:"running"`
	ErrorPresent                bool    `json:"errorPresent"`
	NumRestarts                 int64   `json:"numRestarts"`
	NumReceived                 int64   `json:"numReceived"`
	NumSuccessfullyProcessed    int64   `json:"numSuccessfullyProcessed"`
	NumUserExceptions           int64   `json:"numUserExceptions"`
	LatestUserExceptionsCount   int     `json:"latestUserExceptionsCount"`
	NumSystemExceptions         int64   `json:"numSystemExceptions"`
	LatestSystemExceptionsCount int     `json:"latestSystemExceptionsCount"`
	AverageLatency              float64 `json:"averageLatency"`
	LastInvocationTime          int64   `json:"lastInvocationTime,omitempty"`
	WorkerID                    string  `json:"workerId,omitempty"`
}

type pulsarFunctionStatsResource struct {
	Kind      string                     `json:"kind"`
	URI       string                     `json:"uri"`
	Tenant    string                     `json:"tenant"`
	Namespace string                     `json:"namespace"`
	Name      string                     `json:"name"`
	Stats     pulsarFunctionStatsSummary `json:"stats"`
}

type pulsarFunctionStatsSummary struct {
	ReceivedTotal              int64                                `json:"receivedTotal"`
	ProcessedSuccessfullyTotal int64                                `json:"processedSuccessfullyTotal"`
	SystemExceptionsTotal      int64                                `json:"systemExceptionsTotal"`
	UserExceptionsTotal        int64                                `json:"userExceptionsTotal"`
	AvgProcessLatency          float64                              `json:"avgProcessLatency"`
	LastInvocation             int64                                `json:"lastInvocation,omitempty"`
	OneMin                     pulsarFunctionStatsDataSummary       `json:"oneMin"`
	InstanceCount              int                                  `json:"instanceCount"`
	Instances                  []pulsarFunctionInstanceStatsSummary `json:"instances,omitempty"`
	Limit                      int                                  `json:"limit"`
	Limited                    bool                                 `json:"limited"`
}

type pulsarFunctionStatsDataSummary struct {
	ReceivedTotal              int64   `json:"receivedTotal"`
	ProcessedSuccessfullyTotal int64   `json:"processedSuccessfullyTotal"`
	SystemExceptionsTotal      int64   `json:"systemExceptionsTotal"`
	UserExceptionsTotal        int64   `json:"userExceptionsTotal"`
	AvgProcessLatency          float64 `json:"avgProcessLatency"`
}

type pulsarFunctionInstanceStatsSummary struct {
	InstanceID                 int64                          `json:"instanceId"`
	ReceivedTotal              int64                          `json:"receivedTotal"`
	ProcessedSuccessfullyTotal int64                          `json:"processedSuccessfullyTotal"`
	SystemExceptionsTotal      int64                          `json:"systemExceptionsTotal"`
	UserExceptionsTotal        int64                          `json:"userExceptionsTotal"`
	AvgProcessLatency          float64                        `json:"avgProcessLatency"`
	LastInvocation             int64                          `json:"lastInvocation,omitempty"`
	OneMin                     pulsarFunctionStatsDataSummary `json:"oneMin"`
	UserMetricNames            []string                       `json:"userMetricNames,omitempty"`
	UserMetricCount            int                            `json:"userMetricCount"`
}

type pulsarSourceStatusResource struct {
	Kind      string                    `json:"kind"`
	URI       string                    `json:"uri"`
	Tenant    string                    `json:"tenant"`
	Namespace string                    `json:"namespace"`
	Name      string                    `json:"name"`
	Status    pulsarSourceStatusSummary `json:"status"`
}

type pulsarSourceStatusSummary struct {
	NumInstances int                                 `json:"numInstances"`
	NumRunning   int                                 `json:"numRunning"`
	Instances    []pulsarSourceInstanceStatusSummary `json:"instances,omitempty"`
	Limit        int                                 `json:"limit"`
	Limited      bool                                `json:"limited"`
}

type pulsarSourceInstanceStatusSummary struct {
	InstanceID                  int    `json:"instanceId"`
	Running                     bool   `json:"running"`
	ErrorPresent                bool   `json:"errorPresent"`
	NumRestarts                 int64  `json:"numRestarts"`
	NumReceivedFromSource       int64  `json:"numReceivedFromSource"`
	NumWritten                  int64  `json:"numWritten"`
	NumSystemExceptions         int64  `json:"numSystemExceptions"`
	LatestSystemExceptionsCount int    `json:"latestSystemExceptionsCount"`
	NumSourceExceptions         int64  `json:"numSourceExceptions"`
	LatestSourceExceptionsCount int    `json:"latestSourceExceptionsCount"`
	LastReceivedTime            int64  `json:"lastReceivedTime,omitempty"`
	WorkerID                    string `json:"workerId,omitempty"`
}

type pulsarSinkStatusResource struct {
	Kind      string                  `json:"kind"`
	URI       string                  `json:"uri"`
	Tenant    string                  `json:"tenant"`
	Namespace string                  `json:"namespace"`
	Name      string                  `json:"name"`
	Status    pulsarSinkStatusSummary `json:"status"`
}

type pulsarSinkStatusSummary struct {
	NumInstances int                               `json:"numInstances"`
	NumRunning   int                               `json:"numRunning"`
	Instances    []pulsarSinkInstanceStatusSummary `json:"instances,omitempty"`
	Limit        int                               `json:"limit"`
	Limited      bool                              `json:"limited"`
}

type pulsarSinkInstanceStatusSummary struct {
	InstanceID                  int    `json:"instanceId"`
	Running                     bool   `json:"running"`
	ErrorPresent                bool   `json:"errorPresent"`
	NumRestarts                 int64  `json:"numRestarts"`
	NumReadFromPulsar           int64  `json:"numReadFromPulsar"`
	NumWrittenToSink            int64  `json:"numWrittenToSink"`
	NumSystemExceptions         int64  `json:"numSystemExceptions"`
	LatestSystemExceptionsCount int    `json:"latestSystemExceptionsCount"`
	NumSinkExceptions           int64  `json:"numSinkExceptions"`
	LatestSinkExceptionsCount   int    `json:"latestSinkExceptionsCount"`
	LastReceivedTime            int64  `json:"lastReceivedTime,omitempty"`
	WorkerID                    string `json:"workerId,omitempty"`
}

type pulsarPackageCollectionResource struct {
	Kind        string   `json:"kind"`
	URI         string   `json:"uri"`
	PackageType string   `json:"packageType"`
	Tenant      string   `json:"tenant"`
	Namespace   string   `json:"namespace"`
	Packages    []string `json:"packages"`
	Count       int      `json:"count"`
	Limit       int      `json:"limit"`
	Limited     bool     `json:"limited"`
}

type pulsarPackageVersionsResource struct {
	Kind        string   `json:"kind"`
	URI         string   `json:"uri"`
	PackageType string   `json:"packageType"`
	Tenant      string   `json:"tenant"`
	Namespace   string   `json:"namespace"`
	Package     string   `json:"package"`
	PackageURL  string   `json:"packageUrl"`
	Versions    []string `json:"versions"`
	Count       int      `json:"count"`
	Limit       int      `json:"limit"`
	Limited     bool     `json:"limited"`
}

type pulsarPackageMetadataResource struct {
	Kind        string                       `json:"kind"`
	URI         string                       `json:"uri"`
	PackageType string                       `json:"packageType"`
	Tenant      string                       `json:"tenant"`
	Namespace   string                       `json:"namespace"`
	Package     string                       `json:"package"`
	Version     string                       `json:"version"`
	PackageURL  string                       `json:"packageUrl"`
	Metadata    pulsarPackageMetadataSummary `json:"metadata"`
}

type pulsarPackageMetadataSummary struct {
	Description      string            `json:"description,omitempty"`
	Contact          string            `json:"contact,omitempty"`
	CreateTime       int64             `json:"createTime,omitempty"`
	ModificationTime int64             `json:"modificationTime,omitempty"`
	Properties       map[string]string `json:"properties,omitempty"`
	PropertiesCount  int               `json:"propertiesCount"`
}

type pulsarWorkerClusterResource struct {
	Kind    string              `json:"kind"`
	URI     string              `json:"uri"`
	Workers []*utils.WorkerInfo `json:"workers"`
	Count   int                 `json:"count"`
	Limit   int                 `json:"limit"`
	Limited bool                `json:"limited"`
}

type pulsarWorkerLeaderResource struct {
	Kind   string            `json:"kind"`
	URI    string            `json:"uri"`
	Leader *utils.WorkerInfo `json:"leader,omitempty"`
}

type pulsarWorkerAssignmentsResource struct {
	Kind            string                          `json:"kind"`
	URI             string                          `json:"uri"`
	Workers         []pulsarWorkerAssignmentSummary `json:"workers"`
	WorkerCount     int                             `json:"workerCount"`
	AssignmentCount int                             `json:"assignmentCount"`
	Limit           int                             `json:"limit"`
	Limited         bool                            `json:"limited"`
}

type pulsarWorkerAssignmentSummary struct {
	WorkerID         string   `json:"workerId"`
	Assignments      []string `json:"assignments,omitempty"`
	AssignmentsCount int      `json:"assignmentsCount"`
	Limited          bool     `json:"limited"`
}

type pulsarWorkerFunctionStatsResource struct {
	Kind      string                                     `json:"kind"`
	URI       string                                     `json:"uri"`
	Functions []pulsarWorkerFunctionInstanceStatsSummary `json:"functions"`
	Count     int                                        `json:"count"`
	Limit     int                                        `json:"limit"`
	Limited   bool                                       `json:"limited"`
}

type pulsarWorkerFunctionInstanceStatsSummary struct {
	Name    string                             `json:"name"`
	Metrics pulsarFunctionInstanceStatsSummary `json:"metrics"`
}

type pulsarWorkerMetricsResource struct {
	Kind              string                         `json:"kind"`
	URI               string                         `json:"uri"`
	MonitoringMetrics pulsarMonitoringMetricsSummary `json:"monitoringMetrics"`
}

// PulsarResourceRegistrations contains read-only Pulsar MCP resources and resource templates.
type PulsarResourceRegistrations struct {
	Resources []server.ServerResource
	Templates []server.ServerResourceTemplate
}

// NewPulsarResourceRegistrations builds the read-only Pulsar MCP resource registrations for the enabled features.
func NewPulsarResourceRegistrations(features []string) PulsarResourceRegistrations {
	resourceRegistrations, templateRegistrations := buildPulsarResourceRegistrations(features)
	if len(resourceRegistrations) == 0 && len(templateRegistrations) == 0 {
		return PulsarResourceRegistrations{}
	}

	catalog := buildPulsarResourceCatalog(resourceRegistrations, templateRegistrations)
	baseResources := []server.ServerResource{
		{
			Resource: mcp.NewResource(pulsarResourceContextURI, "Pulsar Context",
				mcp.WithResourceDescription("Current Pulsar session connection metadata with authentication material redacted."),
				mcp.WithMIMEType(pulsarResourceJSONMIMEType),
			),
			Handler: handlePulsarContextResource,
		},
		{
			Resource: mcp.NewResource(pulsarResourceCatalogURI, "Pulsar Resource Catalog",
				mcp.WithResourceDescription("Stable catalog of Pulsar MCP resource URIs and URI templates."),
				mcp.WithMIMEType(pulsarResourceJSONMIMEType),
			),
			Handler: func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
				return handlePulsarCatalogResource(ctx, request, catalog)
			},
		},
	}

	allResources := make([]server.ServerResource, 0, len(baseResources)+len(resourceRegistrations))
	allResources = append(allResources, baseResources...)
	allResources = append(allResources, resourceRegistrations...)

	return PulsarResourceRegistrations{
		Resources: allResources,
		Templates: templateRegistrations,
	}
}

// PulsarAddResources registers the read-only Pulsar MCP resource surface.
func PulsarAddResources(s *server.MCPServer, features []string) {
	registrations := NewPulsarResourceRegistrations(features)
	if len(registrations.Resources) == 0 && len(registrations.Templates) == 0 {
		return
	}

	s.AddResources(registrations.Resources...)
	if len(registrations.Templates) > 0 {
		s.AddResourceTemplates(registrations.Templates...)
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

	if pulsarResourceFeatureEnabled(features, FeaturePulsarAdminFunctions) {
		templates = append(templates,
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarFunctionsTemplateURI, "Pulsar Functions by Namespace",
					mcp.WithTemplateDescription("List Pulsar Functions for a tenant and namespace."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarFunctionsResource,
			},
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarFunctionMetadataTemplateURI, "Pulsar Function Metadata",
					mcp.WithTemplateDescription("Get sanitized metadata for a Pulsar Function."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarFunctionMetadataResource,
			},
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarFunctionStatusTemplateURI, "Pulsar Function Status",
					mcp.WithTemplateDescription("Get bounded runtime status for a Pulsar Function."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarFunctionStatusResource,
			},
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarFunctionStatsTemplateURI, "Pulsar Function Stats",
					mcp.WithTemplateDescription("Get bounded statistics for a Pulsar Function."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarFunctionStatsResource,
			},
		)
	}

	if pulsarResourceFeatureEnabled(features, FeaturePulsarAdminSources) {
		templates = append(templates,
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarSourcesTemplateURI, "Pulsar Sources by Namespace",
					mcp.WithTemplateDescription("List Pulsar Sources for a tenant and namespace."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarSourcesResource,
			},
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarSourceMetadataTemplateURI, "Pulsar Source Metadata",
					mcp.WithTemplateDescription("Get sanitized metadata for a Pulsar Source."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarSourceMetadataResource,
			},
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarSourceStatusTemplateURI, "Pulsar Source Status",
					mcp.WithTemplateDescription("Get bounded runtime status for a Pulsar Source."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarSourceStatusResource,
			},
		)
	}

	if pulsarResourceFeatureEnabled(features, FeaturePulsarAdminSinks) {
		templates = append(templates,
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarSinksTemplateURI, "Pulsar Sinks by Namespace",
					mcp.WithTemplateDescription("List Pulsar Sinks for a tenant and namespace."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarSinksResource,
			},
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarSinkMetadataTemplateURI, "Pulsar Sink Metadata",
					mcp.WithTemplateDescription("Get sanitized metadata for a Pulsar Sink."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarSinkMetadataResource,
			},
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarSinkStatusTemplateURI, "Pulsar Sink Status",
					mcp.WithTemplateDescription("Get bounded runtime status for a Pulsar Sink."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarSinkStatusResource,
			},
		)
	}

	if pulsarResourceFeatureEnabled(features, FeaturePulsarAdminPackages) {
		templates = append(templates,
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarPackagesTemplateURI, "Pulsar Packages by Namespace",
					mcp.WithTemplateDescription("List Pulsar packages by type, tenant, and namespace."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarPackagesResource,
			},
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarPackageVersionsTemplateURI, "Pulsar Package Versions",
					mcp.WithTemplateDescription("List versions for a Pulsar package."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarPackageVersionsResource,
			},
			server.ServerResourceTemplate{
				Template: mcp.NewResourceTemplate(pulsarPackageMetadataTemplateURI, "Pulsar Package Metadata",
					mcp.WithTemplateDescription("Get sanitized metadata for one Pulsar package version."),
					mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarPackageMetadataResource,
			},
		)
	}

	if pulsarResourceFeatureEnabled(features, FeaturePulsarAdminFunctionsWorker) {
		resources = append(resources,
			server.ServerResource{
				Resource: mcp.NewResource(pulsarWorkerClusterResourceURI, "Pulsar Functions Worker Cluster",
					mcp.WithResourceDescription("Bounded summary of Pulsar Functions workers in the current cluster."),
					mcp.WithMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarWorkerClusterResource,
			},
			server.ServerResource{
				Resource: mcp.NewResource(pulsarWorkerLeaderResourceURI, "Pulsar Functions Worker Leader",
					mcp.WithResourceDescription("Current Pulsar Functions worker leader."),
					mcp.WithMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarWorkerLeaderResource,
			},
			server.ServerResource{
				Resource: mcp.NewResource(pulsarWorkerAssignmentsResourceURI, "Pulsar Functions Worker Assignments",
					mcp.WithResourceDescription("Bounded summary of Pulsar Functions worker assignments."),
					mcp.WithMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarWorkerAssignmentsResource,
			},
			server.ServerResource{
				Resource: mcp.NewResource(pulsarWorkerFunctionStatsResourceURI, "Pulsar Functions Worker Function Stats",
					mcp.WithResourceDescription("Bounded function instance stats reported by the Pulsar Functions worker."),
					mcp.WithMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarWorkerFunctionStatsResource,
			},
			server.ServerResource{
				Resource: mcp.NewResource(pulsarWorkerMetricsResourceURI, "Pulsar Functions Worker Metrics",
					mcp.WithResourceDescription("Bounded summary of Pulsar Functions worker monitoring metrics."),
					mcp.WithMIMEType(pulsarResourceJSONMIMEType),
				),
				Handler: handlePulsarWorkerMetricsResource,
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

type pulsarResourceReadFunc func(*pulsarsession.Session, pulsarResourceURI) (any, error)

type pulsarAdminResourceReadFunc func(cmdutils.Client, pulsarResourceURI) (any, error)

func parsePulsarResourceRequest(request mcp.ReadResourceRequest, want pulsarResourceKind) (pulsarResourceURI, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return pulsarResourceURI{}, err
	}
	if parsed.kind != want {
		return pulsarResourceURI{}, fmt.Errorf(
			"unsupported Pulsar resource URI %q: got kind %q, want kind %q",
			request.Params.URI,
			parsed.kind,
			want,
		)
	}
	return parsed, nil
}

func handlePulsarResource(ctx context.Context, request mcp.ReadResourceRequest, want pulsarResourceKind, fn pulsarResourceReadFunc) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceRequest(request, want)
	if err != nil {
		return nil, err
	}
	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	value, err := fn(session, parsed)
	if err != nil {
		return nil, err
	}
	return newPulsarJSONResourceContents(request.Params.URI, value)
}

func handlePulsarAdminResource(ctx context.Context, request mcp.ReadResourceRequest, want pulsarResourceKind, fn pulsarAdminResourceReadFunc) ([]mcp.ResourceContents, error) {
	return handlePulsarResource(ctx, request, want, func(session *pulsarsession.Session, parsed pulsarResourceURI) (any, error) {
		adminClient, err := getPulsarResourceAdminClient(session)
		if err != nil {
			return nil, err
		}
		return fn(adminClient, parsed)
	})
}

func handlePulsarAdminV3Resource(ctx context.Context, request mcp.ReadResourceRequest, want pulsarResourceKind, fn pulsarAdminResourceReadFunc) ([]mcp.ResourceContents, error) {
	return handlePulsarResource(ctx, request, want, func(session *pulsarsession.Session, parsed pulsarResourceURI) (any, error) {
		adminClient, err := getPulsarResourceAdminV3Client(session)
		if err != nil {
			return nil, err
		}
		return fn(adminClient, parsed)
	})
}

func handlePulsarContextResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarResource(ctx, request, pulsarResourceKindContext, func(session *pulsarsession.Session, _ pulsarResourceURI) (any, error) {
		return buildPulsarContextResource(request.Params.URI, session)
	})
}

func handlePulsarCatalogResource(ctx context.Context, request mcp.ReadResourceRequest, catalog pulsarResourceCatalog) ([]mcp.ResourceContents, error) {
	_, err := parsePulsarResourceRequest(request, pulsarResourceKindCatalog)
	if err != nil {
		return nil, err
	}
	runtimeCatalog, found, err := buildPulsarRegisteredResourceCatalogFromServerContext(ctx)
	if err != nil {
		return nil, err
	}
	if found {
		return newPulsarJSONResourceContents(request.Params.URI, runtimeCatalog)
	}
	return newPulsarJSONResourceContents(request.Params.URI, catalog)
}

func handlePulsarTenantsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindTenants, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		tenants, err := adminClient.Tenants().List()
		if err != nil {
			return nil, fmt.Errorf("failed to list tenants: %w", err)
		}

		return pulsarTenantCollectionResource{
			Kind:    string(parsed.kind),
			URI:     request.Params.URI,
			Tenants: tenants,
			Count:   len(tenants),
		}, nil
	})
}

func handlePulsarTenantResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindTenant, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		tenantData, err := adminClient.Tenants().Get(parsed.tenant)
		if err != nil {
			return nil, fmt.Errorf("failed to get tenant %q: %w", parsed.tenant, err)
		}

		return pulsarTenantResource{
			Kind:   string(parsed.kind),
			URI:    request.Params.URI,
			Tenant: parsed.tenant,
			Data:   tenantData,
		}, nil
	})
}

func handlePulsarNamespacesResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindNamespaces, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		namespaces, err := adminClient.Namespaces().GetNamespaces(parsed.tenant)
		if err != nil {
			return nil, fmt.Errorf("failed to list namespaces for tenant %q: %w", parsed.tenant, err)
		}

		return pulsarNamespaceCollectionResource{
			Kind:       string(parsed.kind),
			URI:        request.Params.URI,
			Tenant:     parsed.tenant,
			Namespaces: namespaces,
			Count:      len(namespaces),
		}, nil
	})
}

func handlePulsarNamespaceResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindNamespace, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		namespaceName := parsed.tenant + "/" + parsed.namespace
		policies, err := adminClient.Namespaces().GetPolicies(namespaceName)
		if err != nil {
			return nil, fmt.Errorf("failed to get policies for namespace %q: %w", namespaceName, err)
		}

		return pulsarNamespaceResource{
			Kind:      string(parsed.kind),
			URI:       request.Params.URI,
			Tenant:    parsed.tenant,
			Namespace: parsed.namespace,
			Policies:  policies,
		}, nil
	})
}

func handlePulsarTopicsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindTopics, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		namespaceName := parsed.tenant + "/" + parsed.namespace
		topics, err := adminClient.Namespaces().GetTopics(namespaceName)
		if err != nil {
			return nil, fmt.Errorf("failed to list topics for namespace %q: %w", namespaceName, err)
		}

		return pulsarTopicCollectionResource{
			Kind:      string(parsed.kind),
			URI:       request.Params.URI,
			Tenant:    parsed.tenant,
			Namespace: parsed.namespace,
			Topics:    topics,
			Count:     len(topics),
		}, nil
	})
}

func handlePulsarDefaultResourceQuotaResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindDefaultResourceQuota, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		quota, err := adminClient.ResourceQuotas().GetDefaultResourceQuota()
		if err != nil {
			return nil, fmt.Errorf("failed to get default resource quota: %w", err)
		}

		return pulsarResourceQuotaResource{
			Kind:  string(parsed.kind),
			URI:   request.Params.URI,
			Scope: "default",
			Quota: quota,
		}, nil
	})
}

func handlePulsarResourceQuotaResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindResourceQuota, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
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

		return pulsarResourceQuotaResource{
			Kind:      string(parsed.kind),
			URI:       request.Params.URI,
			Scope:     "namespaceBundle",
			Tenant:    parsed.tenant,
			Namespace: parsed.namespace,
			Bundle:    parsed.bundle,
			Quota:     quota,
		}, nil
	})
}

func handlePulsarClusterStatusResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarResource(ctx, request, pulsarResourceKindStatus, func(session *pulsarsession.Session, parsed pulsarResourceURI) (any, error) {
		status, err := getPulsarClusterStatus(session)
		if err != nil {
			return nil, err
		}

		return pulsarClusterStatusResource{
			Kind:   string(parsed.kind),
			URI:    request.Params.URI,
			Status: status,
		}, nil
	})
}

func handlePulsarClustersResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindClusters, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		clusters, err := adminClient.Clusters().List()
		if err != nil {
			return nil, fmt.Errorf("failed to list clusters: %w", err)
		}

		return pulsarClusterCollectionResource{
			Kind:     string(parsed.kind),
			URI:      request.Params.URI,
			Clusters: clusters,
			Count:    len(clusters),
		}, nil
	})
}

func handlePulsarClusterResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindCluster, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		clusterData, err := adminClient.Clusters().Get(parsed.cluster)
		if err != nil {
			return nil, fmt.Errorf("failed to get cluster %q: %w", parsed.cluster, err)
		}

		return pulsarClusterResource{
			Kind:    string(parsed.kind),
			URI:     request.Params.URI,
			Cluster: parsed.cluster,
			Data:    sanitizePulsarClusterData(parsed.cluster, clusterData),
		}, nil
	})
}

func handlePulsarBrokersResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindBrokers, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		brokers, err := adminClient.Brokers().GetActiveBrokers(parsed.cluster)
		if err != nil {
			return nil, fmt.Errorf("failed to list brokers for cluster %q: %w", parsed.cluster, err)
		}

		return pulsarBrokerCollectionResource{
			Kind:    string(parsed.kind),
			URI:     request.Params.URI,
			Cluster: parsed.cluster,
			Brokers: brokers,
			Count:   len(brokers),
		}, nil
	})
}

func handlePulsarBrokerStatsSummaryResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindBrokerStatsSummary, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		metrics, err := adminClient.BrokerStats().GetMetrics()
		if err != nil {
			return nil, fmt.Errorf("failed to get broker monitoring metrics: %w", err)
		}
		loadReport, err := adminClient.BrokerStats().GetLoadReport()
		if err != nil {
			return nil, fmt.Errorf("failed to get broker load report: %w", err)
		}

		return pulsarBrokerStatsSummaryResource{
			Kind:              string(parsed.kind),
			URI:               request.Params.URI,
			MonitoringMetrics: summarizePulsarMonitoringMetrics(metrics),
			LoadReport:        summarizePulsarBrokerLoadReport(loadReport),
		}, nil
	})
}

func handlePulsarFailureDomainsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindFailureDomains, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		failureDomains, err := adminClient.Clusters().ListFailureDomains(parsed.cluster)
		if err != nil {
			return nil, fmt.Errorf("failed to list failure domains for cluster %q: %w", parsed.cluster, err)
		}
		summaries := summarizePulsarFailureDomains(failureDomains)

		return pulsarFailureDomainCollectionResource{
			Kind:           string(parsed.kind),
			URI:            request.Params.URI,
			Cluster:        parsed.cluster,
			FailureDomains: summaries,
			Count:          len(summaries),
		}, nil
	})
}

func handlePulsarFailureDomainResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindFailureDomain, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		failureDomain, err := adminClient.Clusters().GetFailureDomain(parsed.cluster, parsed.domain)
		if err != nil {
			return nil, fmt.Errorf("failed to get failure domain %q for cluster %q: %w", parsed.domain, parsed.cluster, err)
		}

		return pulsarFailureDomainResource{
			Kind:          string(parsed.kind),
			URI:           request.Params.URI,
			Cluster:       parsed.cluster,
			FailureDomain: summarizePulsarFailureDomain(parsed.domain, failureDomain),
		}, nil
	})
}

func handlePulsarNamespaceIsolationPoliciesResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindNSIsolationPolicies, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		policies, err := adminClient.NsIsolationPolicy().GetNamespaceIsolationPolicies(parsed.cluster)
		if err != nil {
			return nil, fmt.Errorf("failed to list namespace isolation policies for cluster %q: %w", parsed.cluster, err)
		}
		summaries := summarizePulsarNamespaceIsolationPolicies(policies)

		return pulsarNamespaceIsolationPolicyCollectionResource{
			Kind:     string(parsed.kind),
			URI:      request.Params.URI,
			Cluster:  parsed.cluster,
			Policies: summaries,
			Count:    len(summaries),
		}, nil
	})
}

func handlePulsarNamespaceIsolationPolicyResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindNSIsolationPolicy, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		policy, err := adminClient.NsIsolationPolicy().GetNamespaceIsolationPolicy(parsed.cluster, parsed.policy)
		if err != nil {
			return nil, fmt.Errorf("failed to get namespace isolation policy %q for cluster %q: %w", parsed.policy, parsed.cluster, err)
		}
		if policy == nil {
			return nil, fmt.Errorf("namespace isolation policy %q for cluster %q is empty", parsed.policy, parsed.cluster)
		}

		return pulsarNamespaceIsolationPolicyResource{
			Kind:    string(parsed.kind),
			URI:     request.Params.URI,
			Cluster: parsed.cluster,
			Policy: pulsarNamespaceIsolationPolicy{
				Name: parsed.policy,
				Data: *policy,
			},
		}, nil
	})
}

func handlePulsarTopicMetadataResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindTopicMetadata, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		topicName, err := buildPulsarResourceTopicName(parsed)
		if err != nil {
			return nil, err
		}

		properties, err := adminClient.Topics().GetProperties(*topicName)
		if err != nil {
			return nil, fmt.Errorf("failed to get properties for topic %q: %w", topicName.String(), err)
		}

		return pulsarTopicMetadataResource{
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
		}, nil
	})
}

func handlePulsarTopicStatsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindTopicStats, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
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

		return resource, nil
	})
}

func handlePulsarTopicPartitionMetadataResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindTopicPartitions, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		topicName, err := buildPulsarResourceTopicName(parsed)
		if err != nil {
			return nil, err
		}

		metadata, err := adminClient.Topics().GetMetadata(*topicName)
		if err != nil {
			return nil, fmt.Errorf("failed to get partition metadata for topic %q: %w", topicName.String(), err)
		}

		return pulsarTopicPartitionMetadataResource{
			Kind:        string(parsed.kind),
			URI:         request.Params.URI,
			Topic:       topicName.String(),
			Metadata:    metadata,
			Partitioned: metadata.Partitions > 0,
		}, nil
	})
}

func handlePulsarTopicPolicyResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindTopicPolicy, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		topicName, err := buildPulsarResourceTopicName(parsed)
		if err != nil {
			return nil, err
		}

		value, err := readPulsarTopicPolicyValue(adminClient, *topicName, parsed.policy)
		if err != nil {
			return nil, err
		}

		return pulsarTopicPolicyResource{
			Kind:   string(parsed.kind),
			URI:    request.Params.URI,
			Topic:  topicName.String(),
			Policy: parsed.policy,
			Value:  value,
		}, nil
	})
}

func handlePulsarTopicSchemaResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindTopicSchema, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
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

		return pulsarTopicSchemaResource{
			Kind:    string(parsed.kind),
			URI:     request.Params.URI,
			Topic:   topicName.String(),
			Version: schemaInfo.Version,
			Schema:  summarizePulsarTopicSchema(schemaInfo.SchemaInfo),
		}, nil
	})
}

func handlePulsarTopicSchemaVersionResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindTopicSchemaVersion, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
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

		return pulsarTopicSchemaResource{
			Kind:    string(parsed.kind),
			URI:     request.Params.URI,
			Topic:   topicName.String(),
			Version: parsed.version,
			Schema:  summarizePulsarTopicSchema(schemaInfo),
		}, nil
	})
}

func handlePulsarSubscriptionsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindSubscriptions, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		topicName, err := buildPulsarResourceTopicName(parsed)
		if err != nil {
			return nil, err
		}

		subscriptions, err := adminClient.Subscriptions().List(*topicName)
		if err != nil {
			return nil, fmt.Errorf("failed to list subscriptions for topic %q: %w", topicName.String(), err)
		}

		return pulsarSubscriptionCollectionResource{
			Kind:          string(parsed.kind),
			URI:           request.Params.URI,
			Topic:         topicName.String(),
			Domain:        parsed.topicDomain,
			Tenant:        parsed.tenant,
			Namespace:     parsed.namespace,
			Subscriptions: subscriptions,
			Count:         len(subscriptions),
		}, nil
	})
}

func handlePulsarSubscriptionStatsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindSubscriptionStats, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
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

		return pulsarSubscriptionStatsResource{
			Kind:         string(parsed.kind),
			URI:          request.Params.URI,
			Topic:        topicName.String(),
			Subscription: parsed.subscription,
			Partitioned:  partitioned,
			Stats:        summarizePulsarSubscriptionStats(stats),
		}, nil
	})
}

func handlePulsarSubscriptionBacklogResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindSubscriptionBacklog, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
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

		return pulsarSubscriptionBacklogResource{
			Kind:         string(parsed.kind),
			URI:          request.Params.URI,
			Topic:        topicName.String(),
			Subscription: parsed.subscription,
			Partitioned:  partitioned,
			Backlog:      summarizePulsarSubscriptionBacklog(stats),
		}, nil
	})
}

func handlePulsarSubscriptionCursorResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindSubscriptionCursor, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
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

		return pulsarSubscriptionCursorResource{
			Kind:         string(parsed.kind),
			URI:          request.Params.URI,
			Topic:        topicName.String(),
			Subscription: parsed.subscription,
			Cursor:       summarizePulsarSubscriptionCursor(cursor),
		}, nil
	})
}

func handlePulsarFunctionsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminV3Resource(ctx, request, pulsarResourceKindFunctions, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		functions, err := adminClient.Functions().GetFunctions(parsed.tenant, parsed.namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to list functions for namespace %q/%q: %w", parsed.tenant, parsed.namespace, err)
		}
		names, limited := limitStringSlice(functions, pulsarResourceSummaryStringLimit)

		return pulsarWorkloadCollectionResource{
			Kind:      string(parsed.kind),
			URI:       request.Params.URI,
			Type:      "function",
			Tenant:    parsed.tenant,
			Namespace: parsed.namespace,
			Names:     names,
			Count:     len(functions),
			Limit:     pulsarResourceSummaryStringLimit,
			Limited:   limited,
		}, nil
	})
}

func handlePulsarFunctionMetadataResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminV3Resource(ctx, request, pulsarResourceKindFunctionMetadata, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		config, err := adminClient.Functions().GetFunction(parsed.tenant, parsed.namespace, parsed.workload)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get function %q in namespace %q/%q: %w",
				parsed.workload,
				parsed.tenant,
				parsed.namespace,
				err,
			)
		}

		return pulsarFunctionMetadataResource{
			Kind:      string(parsed.kind),
			URI:       request.Params.URI,
			Tenant:    parsed.tenant,
			Namespace: parsed.namespace,
			Name:      parsed.workload,
			Config:    summarizePulsarFunctionConfig(config),
		}, nil
	})
}

func handlePulsarFunctionStatusResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminV3Resource(ctx, request, pulsarResourceKindFunctionStatus, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		status, err := adminClient.Functions().GetFunctionStatus(parsed.tenant, parsed.namespace, parsed.workload)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get status for function %q in namespace %q/%q: %w",
				parsed.workload,
				parsed.tenant,
				parsed.namespace,
				err,
			)
		}

		return pulsarFunctionStatusResource{
			Kind:      string(parsed.kind),
			URI:       request.Params.URI,
			Tenant:    parsed.tenant,
			Namespace: parsed.namespace,
			Name:      parsed.workload,
			Status:    summarizePulsarFunctionStatus(status),
		}, nil
	})
}

func handlePulsarFunctionStatsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminV3Resource(ctx, request, pulsarResourceKindFunctionStats, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		stats, err := adminClient.Functions().GetFunctionStats(parsed.tenant, parsed.namespace, parsed.workload)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get stats for function %q in namespace %q/%q: %w",
				parsed.workload,
				parsed.tenant,
				parsed.namespace,
				err,
			)
		}

		return pulsarFunctionStatsResource{
			Kind:      string(parsed.kind),
			URI:       request.Params.URI,
			Tenant:    parsed.tenant,
			Namespace: parsed.namespace,
			Name:      parsed.workload,
			Stats:     summarizePulsarFunctionStats(stats),
		}, nil
	})
}

func handlePulsarSourcesResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminV3Resource(ctx, request, pulsarResourceKindSources, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		sources, err := adminClient.Sources().ListSources(parsed.tenant, parsed.namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to list sources for namespace %q/%q: %w", parsed.tenant, parsed.namespace, err)
		}
		names, limited := limitStringSlice(sources, pulsarResourceSummaryStringLimit)

		return pulsarWorkloadCollectionResource{
			Kind:      string(parsed.kind),
			URI:       request.Params.URI,
			Type:      "source",
			Tenant:    parsed.tenant,
			Namespace: parsed.namespace,
			Names:     names,
			Count:     len(sources),
			Limit:     pulsarResourceSummaryStringLimit,
			Limited:   limited,
		}, nil
	})
}

func handlePulsarSourceMetadataResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminV3Resource(ctx, request, pulsarResourceKindSourceMetadata, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		config, err := adminClient.Sources().GetSource(parsed.tenant, parsed.namespace, parsed.workload)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get source %q in namespace %q/%q: %w",
				parsed.workload,
				parsed.tenant,
				parsed.namespace,
				err,
			)
		}

		return pulsarSourceMetadataResource{
			Kind:      string(parsed.kind),
			URI:       request.Params.URI,
			Tenant:    parsed.tenant,
			Namespace: parsed.namespace,
			Name:      parsed.workload,
			Config:    summarizePulsarSourceConfig(config),
		}, nil
	})
}

func handlePulsarSourceStatusResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminV3Resource(ctx, request, pulsarResourceKindSourceStatus, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		status, err := adminClient.Sources().GetSourceStatus(parsed.tenant, parsed.namespace, parsed.workload)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get status for source %q in namespace %q/%q: %w",
				parsed.workload,
				parsed.tenant,
				parsed.namespace,
				err,
			)
		}

		return pulsarSourceStatusResource{
			Kind:      string(parsed.kind),
			URI:       request.Params.URI,
			Tenant:    parsed.tenant,
			Namespace: parsed.namespace,
			Name:      parsed.workload,
			Status:    summarizePulsarSourceStatus(status),
		}, nil
	})
}

func handlePulsarSinksResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminV3Resource(ctx, request, pulsarResourceKindSinks, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		sinks, err := adminClient.Sinks().ListSinks(parsed.tenant, parsed.namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to list sinks for namespace %q/%q: %w", parsed.tenant, parsed.namespace, err)
		}
		names, limited := limitStringSlice(sinks, pulsarResourceSummaryStringLimit)

		return pulsarWorkloadCollectionResource{
			Kind:      string(parsed.kind),
			URI:       request.Params.URI,
			Type:      "sink",
			Tenant:    parsed.tenant,
			Namespace: parsed.namespace,
			Names:     names,
			Count:     len(sinks),
			Limit:     pulsarResourceSummaryStringLimit,
			Limited:   limited,
		}, nil
	})
}

func handlePulsarSinkMetadataResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminV3Resource(ctx, request, pulsarResourceKindSinkMetadata, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		config, err := adminClient.Sinks().GetSink(parsed.tenant, parsed.namespace, parsed.workload)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get sink %q in namespace %q/%q: %w",
				parsed.workload,
				parsed.tenant,
				parsed.namespace,
				err,
			)
		}

		return pulsarSinkMetadataResource{
			Kind:      string(parsed.kind),
			URI:       request.Params.URI,
			Tenant:    parsed.tenant,
			Namespace: parsed.namespace,
			Name:      parsed.workload,
			Config:    summarizePulsarSinkConfig(config),
		}, nil
	})
}

func handlePulsarSinkStatusResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminV3Resource(ctx, request, pulsarResourceKindSinkStatus, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		status, err := adminClient.Sinks().GetSinkStatus(parsed.tenant, parsed.namespace, parsed.workload)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get status for sink %q in namespace %q/%q: %w",
				parsed.workload,
				parsed.tenant,
				parsed.namespace,
				err,
			)
		}

		return pulsarSinkStatusResource{
			Kind:      string(parsed.kind),
			URI:       request.Params.URI,
			Tenant:    parsed.tenant,
			Namespace: parsed.namespace,
			Name:      parsed.workload,
			Status:    summarizePulsarSinkStatus(status),
		}, nil
	})
}

func handlePulsarPackagesResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminV3Resource(ctx, request, pulsarResourceKindPackages, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		namespaceName := parsed.tenant + "/" + parsed.namespace
		packages, err := adminClient.Packages().List(parsed.packageType, namespaceName)
		if err != nil {
			return nil, fmt.Errorf("failed to list %s packages for namespace %q: %w", parsed.packageType, namespaceName, err)
		}
		names, limited := limitStringSlice(packages, pulsarResourceSummaryStringLimit)

		return pulsarPackageCollectionResource{
			Kind:        string(parsed.kind),
			URI:         request.Params.URI,
			PackageType: parsed.packageType,
			Tenant:      parsed.tenant,
			Namespace:   parsed.namespace,
			Packages:    names,
			Count:       len(packages),
			Limit:       pulsarResourceSummaryStringLimit,
			Limited:     limited,
		}, nil
	})
}

func handlePulsarPackageVersionsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminV3Resource(ctx, request, pulsarResourceKindPackageVersions, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		packageURL := buildPulsarPackageURL(parsed, "")
		versions, err := adminClient.Packages().ListVersions(packageURL)
		if err != nil {
			return nil, fmt.Errorf("failed to list versions for package %q: %w", packageURL, err)
		}
		names, limited := limitStringSlice(versions, pulsarResourceSummaryStringLimit)

		return pulsarPackageVersionsResource{
			Kind:        string(parsed.kind),
			URI:         request.Params.URI,
			PackageType: parsed.packageType,
			Tenant:      parsed.tenant,
			Namespace:   parsed.namespace,
			Package:     parsed.packageName,
			PackageURL:  packageURL,
			Versions:    names,
			Count:       len(versions),
			Limit:       pulsarResourceSummaryStringLimit,
			Limited:     limited,
		}, nil
	})
}

func handlePulsarPackageMetadataResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminV3Resource(ctx, request, pulsarResourceKindPackageMetadata, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		packageURL := buildPulsarPackageURL(parsed, parsed.versionName)
		metadata, err := adminClient.Packages().GetMetadata(packageURL)
		if err != nil {
			return nil, fmt.Errorf("failed to get metadata for package %q: %w", packageURL, err)
		}

		return pulsarPackageMetadataResource{
			Kind:        string(parsed.kind),
			URI:         request.Params.URI,
			PackageType: parsed.packageType,
			Tenant:      parsed.tenant,
			Namespace:   parsed.namespace,
			Package:     parsed.packageName,
			Version:     parsed.versionName,
			PackageURL:  packageURL,
			Metadata:    summarizePulsarPackageMetadata(metadata),
		}, nil
	})
}

func handlePulsarWorkerClusterResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindWorkerCluster, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		workers, err := adminClient.FunctionsWorker().GetCluster()
		if err != nil {
			return nil, fmt.Errorf("failed to get functions worker cluster: %w", err)
		}
		limitedWorkers, limited := limitWorkerInfoSlice(workers, pulsarResourceSummaryStringLimit)

		return pulsarWorkerClusterResource{
			Kind:    string(parsed.kind),
			URI:     request.Params.URI,
			Workers: limitedWorkers,
			Count:   len(workers),
			Limit:   pulsarResourceSummaryStringLimit,
			Limited: limited,
		}, nil
	})
}

func handlePulsarWorkerLeaderResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindWorkerLeader, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		leader, err := adminClient.FunctionsWorker().GetClusterLeader()
		if err != nil {
			return nil, fmt.Errorf("failed to get functions worker leader: %w", err)
		}

		return pulsarWorkerLeaderResource{
			Kind:   string(parsed.kind),
			URI:    request.Params.URI,
			Leader: leader,
		}, nil
	})
}

func handlePulsarWorkerAssignmentsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindWorkerAssignments, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		assignments, err := adminClient.FunctionsWorker().GetAssignments()
		if err != nil {
			return nil, fmt.Errorf("failed to get functions worker assignments: %w", err)
		}

		return summarizePulsarWorkerAssignments(request.Params.URI, parsed.kind, assignments), nil
	})
}

func handlePulsarWorkerFunctionStatsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindWorkerFunctionStats, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		stats, err := adminClient.FunctionsWorker().GetFunctionsStats()
		if err != nil {
			return nil, fmt.Errorf("failed to get functions worker stats: %w", err)
		}

		return summarizePulsarWorkerFunctionStats(request.Params.URI, parsed.kind, stats), nil
	})
}

func handlePulsarWorkerMetricsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return handlePulsarAdminResource(ctx, request, pulsarResourceKindWorkerMetrics, func(adminClient cmdutils.Client, parsed pulsarResourceURI) (any, error) {
		metrics, err := adminClient.FunctionsWorker().GetMetrics()
		if err != nil {
			return nil, fmt.Errorf("failed to get functions worker metrics: %w", err)
		}

		return pulsarWorkerMetricsResource{
			Kind:              string(parsed.kind),
			URI:               request.Params.URI,
			MonitoringMetrics: summarizePulsarPointerMonitoringMetrics(metrics),
		}, nil
	})
}

func requirePulsarResourceSession(ctx context.Context) (*pulsarsession.Session, error) {
	session := GetPulsarSession(ctx)
	if session == nil {
		return nil, fmt.Errorf("pulsar session not found in context")
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
		return nil, fmt.Errorf("pulsar admin client not found in session")
	}
	return adminClient, nil
}

func getPulsarResourceAdminV3Client(session *pulsarsession.Session) (cmdutils.Client, error) {
	if _, err := session.GetPulsarCtlConfig(); err != nil {
		return nil, fmt.Errorf("failed to get Pulsar admin configuration: %w", err)
	}
	adminClient, err := session.GetAdminV3Client()
	if err != nil {
		return nil, fmt.Errorf("failed to get Pulsar admin v3 client: %w", err)
	}
	if adminClient == nil {
		return nil, fmt.Errorf("pulsar admin v3 client not found in session")
	}
	return adminClient, nil
}

func getPulsarClusterStatus(session *pulsarsession.Session) (string, error) {
	statusClient, err := session.GetAdminStatusClient()
	if err != nil {
		return "", fmt.Errorf("failed to get Pulsar status client: %w", err)
	}
	data, err := statusClient.GetWithQueryParams("/status.html", nil, nil, false)
	if err != nil {
		return "", fmt.Errorf("failed to check Pulsar status: %w", err)
	}

	return strings.TrimSpace(string(data)), nil
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

func summarizePulsarFunctionConfig(config utils.FunctionConfig) pulsarFunctionConfigSummary {
	inputs, _ := limitStringSlice(config.Inputs, pulsarResourceSummaryStringLimit)
	locationType := pulsarFunctionPackageLocationType(config)
	return pulsarFunctionConfigSummary{
		Tenant:                             config.Tenant,
		Namespace:                          config.Namespace,
		Name:                               config.Name,
		ClassName:                          config.ClassName,
		Runtime:                            config.Runtime,
		FunctionTypeConfigured:             config.FunctionType != nil && strings.TrimSpace(*config.FunctionType) != "",
		PackageLocationConfigured:          locationType != "",
		PackageLocationType:                locationType,
		Parallelism:                        config.Parallelism,
		Inputs:                             inputs,
		InputsCount:                        len(config.Inputs),
		InputSpecsCount:                    len(config.InputSpecs),
		TopicsPattern:                      stringPointerValue(config.TopicsPattern),
		Output:                             config.Output,
		LogTopic:                           config.LogTopic,
		ProcessingGuarantees:               config.ProcessingGuarantees,
		Resources:                          config.Resources,
		UserConfig:                         sanitizePulsarResourceAnyMap(config.UserConfig, pulsarResourceSummaryStringLimit),
		UserConfigCount:                    len(config.UserConfig),
		SecretsCount:                       len(config.Secrets),
		ProducerConfigConfigured:           config.ProducerConfig != nil,
		CustomSchemaOutputsCount:           len(config.CustomSchemaOutputs),
		CustomSerdeInputsCount:             len(config.CustomSerdeInputs),
		CustomSchemaInputsCount:            len(config.CustomSchemaInputs),
		CustomRuntimeOptionsConfigured:     config.CustomRuntimeOptions != "",
		DeadLetterTopic:                    config.DeadLetterTopic,
		SubscriptionName:                   config.SubName,
		SubscriptionPosition:               config.SubscriptionPosition,
		TimeoutMs:                          config.TimeoutMs,
		MaxMessageRetries:                  config.MaxMessageRetries,
		CleanupSubscription:                config.CleanupSubscription,
		RetainOrdering:                     config.RetainOrdering,
		RetainKeyOrdering:                  config.RetainKeyOrdering,
		AutoAck:                            config.AutoAck,
		ForwardSourceMessageProperty:       config.ForwardSourceMessageProperty,
		ExposePulsarAdminClientEnabled:     config.ExposePulsarAdminClientEnabled,
		SkipToLatest:                       config.SkipToLatest,
		MaxPendingAsyncRequests:            config.MaxPendingAsyncRequests,
		WindowConfigConfigured:             config.WindowConfig != nil,
		InputTypeClassName:                 config.InputTypeClassName,
		OutputTypeClassName:                config.OutputTypeClassName,
		OutputSchemaType:                   config.OutputSchemaType,
		OutputSerdeClassName:               config.OutputSerdeClassName,
		CustomSchemaOutputs:                sanitizePulsarResourceStringMap(config.CustomSchemaOutputs, pulsarResourceSummaryStringLimit),
		SanitizedCustomSchemaOutputsCount:  len(config.CustomSchemaOutputs),
		SanitizedCustomSerdeInputs:         sanitizePulsarResourceStringMap(config.CustomSerdeInputs, pulsarResourceSummaryStringLimit),
		SanitizedCustomSerdeInputsCount:    len(config.CustomSerdeInputs),
		SanitizedCustomSchemaInputs:        sanitizePulsarResourceStringMap(config.CustomSchemaInputs, pulsarResourceSummaryStringLimit),
		SanitizedCustomSchemaInputsCount:   len(config.CustomSchemaInputs),
		SanitizedInputSpecsSchemaSummaries: summarizePulsarInputSpecs(config.InputSpecs, pulsarResourceSummaryStringLimit),
	}
}

func pulsarFunctionPackageLocationType(config utils.FunctionConfig) string {
	switch {
	case config.Jar != nil && strings.TrimSpace(*config.Jar) != "":
		return "jar"
	case config.Py != nil && strings.TrimSpace(*config.Py) != "":
		return "py"
	case config.Go != nil && strings.TrimSpace(*config.Go) != "":
		return "go"
	case config.FunctionType != nil && strings.TrimSpace(*config.FunctionType) != "":
		return "builtin"
	default:
		return ""
	}
}

func summarizePulsarInputSpecs(inputSpecs map[string]utils.ConsumerConfig, limit int) []pulsarInputSpecSummary {
	if len(inputSpecs) == 0 || limit <= 0 {
		return nil
	}
	topics := make([]string, 0, len(inputSpecs))
	for topic := range inputSpecs {
		topics = append(topics, topic)
	}
	sort.Strings(topics)
	if len(topics) > limit {
		topics = topics[:limit]
	}
	summaries := make([]pulsarInputSpecSummary, 0, len(topics))
	for _, topic := range topics {
		config := inputSpecs[topic]
		summaries = append(summaries, pulsarInputSpecSummary{
			Topic:                   topic,
			SchemaType:              config.SchemaType,
			SerdeClassName:          config.SerdeClassName,
			RegexPattern:            config.RegexPattern,
			ReceiverQueueSize:       config.ReceiverQueueSize,
			SchemaProperties:        sanitizePulsarResourceStringMap(config.SchemaProperties, pulsarResourceSummaryStringLimit),
			SchemaPropertiesCount:   len(config.SchemaProperties),
			ConsumerProperties:      sanitizePulsarResourceStringMap(config.ConsumerProperties, pulsarResourceSummaryStringLimit),
			ConsumerPropertiesCount: len(config.ConsumerProperties),
			CryptoConfigConfigured:  config.CryptoConfig != nil,
			PoolMessages:            config.PoolMessages,
		})
	}
	return summaries
}

func summarizePulsarSourceConfig(config utils.SourceConfig) pulsarSourceConfigSummary {
	return pulsarSourceConfigSummary{
		Tenant:                         config.Tenant,
		Namespace:                      config.Namespace,
		Name:                           config.Name,
		ClassName:                      config.ClassName,
		TopicName:                      config.TopicName,
		SerdeClassName:                 config.SerdeClassName,
		SchemaType:                     config.SchemaType,
		Parallelism:                    config.Parallelism,
		ProcessingGuarantees:           config.ProcessingGuarantees,
		Resources:                      config.Resources,
		ArchiveConfigured:              config.Archive != "",
		ProducerConfigConfigured:       config.ProducerConfig != nil,
		RuntimeFlagsConfigured:         config.RuntimeFlags != "",
		CustomRuntimeOptionsConfigured: config.CustomRuntimeOptions != "",
		BatchSourceConfigConfigured:    config.BatchSourceConfig != nil,
		BatchBuilder:                   config.BatchBuilder,
		Configs:                        sanitizePulsarResourceAnyMap(config.Configs, pulsarResourceSummaryStringLimit),
		ConfigsCount:                   len(config.Configs),
		SecretsCount:                   len(config.Secrets),
	}
}

func summarizePulsarSinkConfig(config utils.SinkConfig) pulsarSinkConfigSummary {
	inputs, _ := limitStringSlice(config.Inputs, pulsarResourceSummaryStringLimit)
	return pulsarSinkConfigSummary{
		Tenant:                         config.Tenant,
		Namespace:                      config.Namespace,
		Name:                           config.Name,
		ClassName:                      config.ClassName,
		SinkType:                       config.SinkType,
		Parallelism:                    config.Parallelism,
		Inputs:                         inputs,
		InputsCount:                    len(config.Inputs),
		InputSpecsCount:                len(config.InputSpecs),
		TopicsPattern:                  stringPointerValue(config.TopicsPattern),
		ProcessingGuarantees:           config.ProcessingGuarantees,
		SourceSubscriptionName:         config.SourceSubscriptionName,
		SourceSubscriptionPosition:     config.SourceSubscriptionPosition,
		Resources:                      config.Resources,
		TimeoutMs:                      config.TimeoutMs,
		ArchiveConfigured:              config.Archive != "",
		RuntimeFlagsConfigured:         config.RuntimeFlags != "",
		CustomRuntimeOptionsConfigured: config.CustomRuntimeOptions != "",
		CleanupSubscription:            config.CleanupSubscription,
		RetainOrdering:                 config.RetainOrdering,
		RetainKeyOrdering:              config.RetainKeyOrdering,
		AutoAck:                        config.AutoAck,
		TopicToSerdeClassNameCount:     len(config.TopicToSerdeClassName),
		TopicToSchemaTypeCount:         len(config.TopicToSchemaType),
		TopicToSchemaPropertiesCount:   len(config.TopicToSchemaProperties),
		Configs:                        sanitizePulsarResourceAnyMap(config.Configs, pulsarResourceSummaryStringLimit),
		ConfigsCount:                   len(config.Configs),
		SecretsCount:                   len(config.Secrets),
		MaxMessageRetries:              config.MaxMessageRetries,
		DeadLetterTopic:                config.DeadLetterTopic,
		NegativeAckRedeliveryDelayMs:   config.NegativeAckRedeliveryDelayMs,
		TransformFunctionConfigured: config.TransformFunction != "" ||
			config.TransformFunctionClassName != "" ||
			config.TransformFunctionConfig != "",
	}
}

func summarizePulsarFunctionStatus(status utils.FunctionStatus) pulsarFunctionStatusSummary {
	instances := status.Instances
	limited := false
	if len(instances) > pulsarResourceSummaryStringLimit {
		instances = instances[:pulsarResourceSummaryStringLimit]
		limited = true
	}
	summaries := make([]pulsarFunctionInstanceStatusSummary, 0, len(instances))
	for _, instance := range instances {
		summaries = append(summaries, summarizePulsarFunctionInstanceStatus(instance.InstanceID, instance.Status))
	}
	return pulsarFunctionStatusSummary{
		NumInstances: status.NumInstances,
		NumRunning:   status.NumRunning,
		Instances:    summaries,
		Limit:        pulsarResourceSummaryStringLimit,
		Limited:      limited,
	}
}

func summarizePulsarFunctionInstanceStatus(instanceID int, status utils.FunctionInstanceStatusData) pulsarFunctionInstanceStatusSummary {
	return pulsarFunctionInstanceStatusSummary{
		InstanceID:                  instanceID,
		Running:                     status.Running,
		ErrorPresent:                status.Err != "",
		NumRestarts:                 status.NumRestarts,
		NumReceived:                 status.NumReceived,
		NumSuccessfullyProcessed:    status.NumSuccessfullyProcessed,
		NumUserExceptions:           status.NumUserExceptions,
		LatestUserExceptionsCount:   len(status.LatestUserExceptions),
		NumSystemExceptions:         status.NumSystemExceptions,
		LatestSystemExceptionsCount: len(status.LatestSystemExceptions),
		AverageLatency:              status.AverageLatency,
		LastInvocationTime:          status.LastInvocationTime,
		WorkerID:                    status.WorkerID,
	}
}

func summarizePulsarFunctionStats(stats utils.FunctionStats) pulsarFunctionStatsSummary {
	instances := stats.Instances
	limited := false
	if len(instances) > pulsarResourceSummaryStringLimit {
		instances = instances[:pulsarResourceSummaryStringLimit]
		limited = true
	}
	summaries := make([]pulsarFunctionInstanceStatsSummary, 0, len(instances))
	for _, instance := range instances {
		summaries = append(summaries, summarizePulsarFunctionInstanceStats(instance))
	}
	return pulsarFunctionStatsSummary{
		ReceivedTotal:              stats.ReceivedTotal,
		ProcessedSuccessfullyTotal: stats.ProcessedSuccessfullyTotal,
		SystemExceptionsTotal:      stats.SystemExceptionsTotal,
		UserExceptionsTotal:        stats.UserExceptionsTotal,
		AvgProcessLatency:          stats.AvgProcessLatency,
		LastInvocation:             stats.LastInvocation,
		OneMin:                     summarizePulsarFunctionStatsData(stats.OneMin),
		InstanceCount:              len(stats.Instances),
		Instances:                  summaries,
		Limit:                      pulsarResourceSummaryStringLimit,
		Limited:                    limited,
	}
}

func summarizePulsarFunctionInstanceStats(stats utils.FunctionInstanceStats) pulsarFunctionInstanceStatsSummary {
	metricNames := make(map[string]struct{}, len(stats.Metrics.UserMetrics))
	for name := range stats.Metrics.UserMetrics {
		metricNames[name] = struct{}{}
	}
	return pulsarFunctionInstanceStatsSummary{
		InstanceID:                 stats.InstanceID,
		ReceivedTotal:              stats.Metrics.ReceivedTotal,
		ProcessedSuccessfullyTotal: stats.Metrics.ProcessedSuccessfullyTotal,
		SystemExceptionsTotal:      stats.Metrics.SystemExceptionsTotal,
		UserExceptionsTotal:        stats.Metrics.UserExceptionsTotal,
		AvgProcessLatency:          stats.Metrics.AvgProcessLatency,
		LastInvocation:             stats.Metrics.LastInvocation,
		OneMin:                     summarizePulsarFunctionStatsData(stats.Metrics.OneMin),
		UserMetricNames:            sortedLimitedStrings(metricNames, pulsarResourceSummaryStringLimit),
		UserMetricCount:            len(stats.Metrics.UserMetrics),
	}
}

func summarizePulsarFunctionStatsData(stats utils.FunctionInstanceStatsDataBase) pulsarFunctionStatsDataSummary {
	return pulsarFunctionStatsDataSummary{
		ReceivedTotal:              stats.ReceivedTotal,
		ProcessedSuccessfullyTotal: stats.ProcessedSuccessfullyTotal,
		SystemExceptionsTotal:      stats.SystemExceptionsTotal,
		UserExceptionsTotal:        stats.UserExceptionsTotal,
		AvgProcessLatency:          stats.AvgProcessLatency,
	}
}

func summarizePulsarSourceStatus(status utils.SourceStatus) pulsarSourceStatusSummary {
	instances := status.Instances
	limited := false
	if len(instances) > pulsarResourceSummaryStringLimit {
		instances = instances[:pulsarResourceSummaryStringLimit]
		limited = true
	}
	summaries := make([]pulsarSourceInstanceStatusSummary, 0, len(instances))
	for _, instance := range instances {
		if instance == nil {
			continue
		}
		summaries = append(summaries, summarizePulsarSourceInstanceStatus(instance.InstanceID, instance.Status))
	}
	return pulsarSourceStatusSummary{
		NumInstances: status.NumInstances,
		NumRunning:   status.NumRunning,
		Instances:    summaries,
		Limit:        pulsarResourceSummaryStringLimit,
		Limited:      limited,
	}
}

func summarizePulsarSourceInstanceStatus(instanceID int, status utils.SourceInstanceStatusData) pulsarSourceInstanceStatusSummary {
	return pulsarSourceInstanceStatusSummary{
		InstanceID:                  instanceID,
		Running:                     status.Running,
		ErrorPresent:                status.Err != "",
		NumRestarts:                 status.NumRestarts,
		NumReceivedFromSource:       status.NumReceivedFromSource,
		NumWritten:                  status.NumWritten,
		NumSystemExceptions:         status.NumSystemExceptions,
		LatestSystemExceptionsCount: len(status.LatestSystemExceptions),
		NumSourceExceptions:         status.NumSourceExceptions,
		LatestSourceExceptionsCount: len(status.LatestSourceExceptions),
		LastReceivedTime:            status.LastReceivedTime,
		WorkerID:                    status.WorkerID,
	}
}

func summarizePulsarSinkStatus(status utils.SinkStatus) pulsarSinkStatusSummary {
	instances := status.Instances
	limited := false
	if len(instances) > pulsarResourceSummaryStringLimit {
		instances = instances[:pulsarResourceSummaryStringLimit]
		limited = true
	}
	summaries := make([]pulsarSinkInstanceStatusSummary, 0, len(instances))
	for _, instance := range instances {
		if instance == nil {
			continue
		}
		summaries = append(summaries, summarizePulsarSinkInstanceStatus(instance.InstanceID, instance.Status))
	}
	return pulsarSinkStatusSummary{
		NumInstances: status.NumInstances,
		NumRunning:   status.NumRunning,
		Instances:    summaries,
		Limit:        pulsarResourceSummaryStringLimit,
		Limited:      limited,
	}
}

func summarizePulsarSinkInstanceStatus(instanceID int, status utils.SinkInstanceStatusData) pulsarSinkInstanceStatusSummary {
	return pulsarSinkInstanceStatusSummary{
		InstanceID:                  instanceID,
		Running:                     status.Running,
		ErrorPresent:                status.Err != "",
		NumRestarts:                 status.NumRestarts,
		NumReadFromPulsar:           status.NumReadFromPulsar,
		NumWrittenToSink:            status.NumWrittenToSink,
		NumSystemExceptions:         status.NumSystemExceptions,
		LatestSystemExceptionsCount: len(status.LatestSystemExceptions),
		NumSinkExceptions:           status.NumSinkExceptions,
		LatestSinkExceptionsCount:   len(status.LatestSinkExceptions),
		LastReceivedTime:            status.LastReceivedTime,
		WorkerID:                    status.WorkerID,
	}
}

func summarizePulsarPackageMetadata(metadata utils.PackageMetadata) pulsarPackageMetadataSummary {
	return pulsarPackageMetadataSummary{
		Description:      metadata.Description,
		Contact:          metadata.Contact,
		CreateTime:       metadata.CreateTime,
		ModificationTime: metadata.ModificationTime,
		Properties:       sanitizePulsarResourceStringMap(metadata.Properties, pulsarResourceSummaryStringLimit),
		PropertiesCount:  len(metadata.Properties),
	}
}

func summarizePulsarWorkerAssignments(
	uri string,
	kind pulsarResourceKind,
	assignments map[string][]string,
) pulsarWorkerAssignmentsResource {
	workerIDs := make([]string, 0, len(assignments))
	for workerID := range assignments {
		workerIDs = append(workerIDs, workerID)
	}
	sort.Strings(workerIDs)
	limited := false
	if len(workerIDs) > pulsarResourceSummaryStringLimit {
		workerIDs = workerIDs[:pulsarResourceSummaryStringLimit]
		limited = true
	}

	summaries := make([]pulsarWorkerAssignmentSummary, 0, len(workerIDs))
	totalAssignments := 0
	for _, values := range assignments {
		totalAssignments += len(values)
	}
	for _, workerID := range workerIDs {
		values := append([]string(nil), assignments[workerID]...)
		sort.Strings(values)
		limitedValues, valuesLimited := limitStringSlice(values, pulsarResourceSummaryStringLimit)
		if valuesLimited {
			limited = true
		}
		summaries = append(summaries, pulsarWorkerAssignmentSummary{
			WorkerID:         workerID,
			Assignments:      limitedValues,
			AssignmentsCount: len(assignments[workerID]),
			Limited:          valuesLimited,
		})
	}

	return pulsarWorkerAssignmentsResource{
		Kind:            string(kind),
		URI:             uri,
		Workers:         summaries,
		WorkerCount:     len(assignments),
		AssignmentCount: totalAssignments,
		Limit:           pulsarResourceSummaryStringLimit,
		Limited:         limited,
	}
}

func summarizePulsarWorkerFunctionStats(
	uri string,
	kind pulsarResourceKind,
	stats []*utils.WorkerFunctionInstanceStats,
) pulsarWorkerFunctionStatsResource {
	values := append([]*utils.WorkerFunctionInstanceStats(nil), stats...)
	sort.Slice(values, func(i, j int) bool {
		if values[i] == nil {
			return false
		}
		if values[j] == nil {
			return true
		}
		return values[i].Name < values[j].Name
	})
	limited := false
	if len(values) > pulsarResourceSummaryStringLimit {
		values = values[:pulsarResourceSummaryStringLimit]
		limited = true
	}

	summaries := make([]pulsarWorkerFunctionInstanceStatsSummary, 0, len(values))
	for _, value := range values {
		if value == nil {
			continue
		}
		summaries = append(summaries, pulsarWorkerFunctionInstanceStatsSummary{
			Name: value.Name,
			Metrics: summarizePulsarFunctionInstanceStats(utils.FunctionInstanceStats{
				InstanceID: 0,
				Metrics:    value.Metrics,
			}),
		})
	}
	return pulsarWorkerFunctionStatsResource{
		Kind:      string(kind),
		URI:       uri,
		Functions: summaries,
		Count:     len(stats),
		Limit:     pulsarResourceSummaryStringLimit,
		Limited:   limited,
	}
}

func summarizePulsarPointerMonitoringMetrics(metrics []*utils.Metrics) pulsarMonitoringMetricsSummary {
	values := make([]utils.Metrics, 0, len(metrics))
	for _, metric := range metrics {
		if metric == nil {
			continue
		}
		values = append(values, *metric)
	}
	return summarizePulsarMonitoringMetrics(values)
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

func buildPulsarPackageURL(parsed pulsarResourceURI, version string) string {
	packageURL := fmt.Sprintf(
		"%s://%s/%s/%s",
		parsed.packageType,
		parsed.tenant,
		parsed.namespace,
		parsed.packageName,
	)
	if version != "" {
		packageURL += "@" + version
	}
	return packageURL
}

func stringPointerValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
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

func sanitizePulsarResourceAnyMap(values map[string]any, limit int) map[string]any {
	return sanitizePulsarResourceAnyMapWithDepth(values, limit, pulsarResourceSanitizeDepthLimit)
}

func sanitizePulsarResourceAnyMapWithDepth(values map[string]any, limit int, depth int) map[string]any {
	if len(values) == 0 || limit <= 0 {
		return nil
	}
	if depth <= 0 {
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

	sanitized := make(map[string]any, len(keys))
	for _, key := range keys {
		if isSensitivePulsarResourceKey(key) {
			sanitized[key] = pulsarResourceRedactedValue
			continue
		}
		sanitized[key] = sanitizePulsarResourceAnyValue(values[key], depth-1)
	}
	return sanitized
}

func sanitizePulsarResourceAnyValue(value any, depth int) any {
	switch typed := value.(type) {
	case map[string]any:
		if depth <= 0 {
			return pulsarResourceRedactedValue
		}
		return sanitizePulsarResourceAnyMapWithDepth(typed, pulsarResourceSummaryStringLimit, depth)
	case map[string]string:
		if depth <= 0 {
			return pulsarResourceRedactedValue
		}
		return sanitizePulsarResourceStringMap(typed, pulsarResourceSummaryStringLimit)
	case []any:
		if depth <= 0 {
			return pulsarResourceRedactedValue
		}
		values := typed
		if len(values) > pulsarResourceSummaryStringLimit {
			values = values[:pulsarResourceSummaryStringLimit]
		}
		sanitized := make([]any, 0, len(values))
		for _, item := range values {
			sanitized = append(sanitized, sanitizePulsarResourceAnyValue(item, depth-1))
		}
		return sanitized
	case []string:
		if depth <= 0 {
			return pulsarResourceRedactedValue
		}
		values, _ := limitStringSlice(typed, pulsarResourceSummaryStringLimit)
		return values
	default:
		return value
	}
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

func limitStringSlice(values []string, limit int) ([]string, bool) {
	if len(values) == 0 || limit <= 0 {
		return nil, false
	}
	limited := len(values) > limit
	if limited {
		return append([]string(nil), values[:limit]...), true
	}
	return append([]string(nil), values...), false
}

func limitWorkerInfoSlice(values []*utils.WorkerInfo, limit int) ([]*utils.WorkerInfo, bool) {
	if len(values) == 0 || limit <= 0 {
		return nil, false
	}
	limited := len(values) > limit
	if limited {
		return append([]*utils.WorkerInfo(nil), values[:limit]...), true
	}
	return append([]*utils.WorkerInfo(nil), values...), false
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
		return pulsarContextResource{}, fmt.Errorf("pulsar session is not configured")
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

func buildPulsarRegisteredResourceCatalogFromServerContext(ctx context.Context) (pulsarResourceCatalog, bool, error) {
	upstream := server.ServerFromContext(ctx)
	if upstream == nil {
		return pulsarResourceCatalog{}, false, nil
	}

	resources, err := listPulsarRegisteredResources(ctx, upstream)
	if err != nil {
		return pulsarResourceCatalog{}, false, err
	}
	templates, err := listPulsarRegisteredResourceTemplates(ctx, upstream)
	if err != nil {
		return pulsarResourceCatalog{}, false, err
	}

	return buildPulsarCatalogFromRegisteredEntries(resources, templates), true, nil
}

func listPulsarRegisteredResources(ctx context.Context, upstream *server.MCPServer) ([]mcp.Resource, error) {
	response, err := handlePulsarCatalogServerRequest(
		ctx,
		upstream,
		map[string]any{
			"jsonrpc": mcp.JSONRPC_VERSION,
			"id":      "pulsar-resource-catalog-list",
			"method":  "resources/list",
		},
	)
	if err != nil {
		return nil, err
	}

	result, ok := response.Result.(mcp.ListResourcesResult)
	if !ok {
		return nil, fmt.Errorf("unexpected resources/list result type %T", response.Result)
	}

	resources := make([]mcp.Resource, 0, len(result.Resources))
	for _, resource := range result.Resources {
		if !isPulsarCatalogURI(resource.URI) {
			continue
		}
		resources = append(resources, resource)
	}
	return resources, nil
}

func listPulsarRegisteredResourceTemplates(ctx context.Context, upstream *server.MCPServer) ([]mcp.ResourceTemplate, error) {
	response, err := handlePulsarCatalogServerRequest(
		ctx,
		upstream,
		map[string]any{
			"jsonrpc": mcp.JSONRPC_VERSION,
			"id":      "pulsar-resource-catalog-templates",
			"method":  "resources/templates/list",
		},
	)
	if err != nil {
		return nil, err
	}

	result, ok := response.Result.(mcp.ListResourceTemplatesResult)
	if !ok {
		return nil, fmt.Errorf("unexpected resources/templates/list result type %T", response.Result)
	}

	templates := make([]mcp.ResourceTemplate, 0, len(result.ResourceTemplates))
	for _, template := range result.ResourceTemplates {
		if template.URITemplate == nil || !isPulsarCatalogURI(template.URITemplate.Raw()) {
			continue
		}
		templates = append(templates, template)
	}
	return templates, nil
}

func handlePulsarCatalogServerRequest(
	ctx context.Context,
	upstream *server.MCPServer,
	request map[string]any,
) (mcp.JSONRPCResponse, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return mcp.JSONRPCResponse{}, fmt.Errorf("failed to marshal Pulsar catalog request: %w", err)
	}

	message := upstream.HandleMessage(ctx, payload)
	resp, ok := message.(mcp.JSONRPCResponse)
	if ok {
		return resp, nil
	}

	if errResp, ok := message.(mcp.JSONRPCError); ok {
		return mcp.JSONRPCResponse{}, fmt.Errorf("pulsar catalog request failed: %s", errResp.Error.Message)
	}

	return mcp.JSONRPCResponse{}, fmt.Errorf("unexpected Pulsar catalog response type %T", message)
}

func isPulsarCatalogURI(uri string) bool {
	return strings.HasPrefix(uri, "pulsar://")
}

func buildPulsarCatalogFromRegisteredEntries(
	resources []mcp.Resource,
	templates []mcp.ResourceTemplate,
) pulsarResourceCatalog {
	catalog := newPulsarResourceCatalog()

	for _, resource := range resources {
		catalog.Resources = append(catalog.Resources, pulsarCatalogResource{
			URI:         resource.URI,
			Name:        resource.Name,
			Description: resource.Description,
			MIMEType:    resource.MIMEType,
		})
	}
	for _, template := range templates {
		if template.URITemplate == nil {
			continue
		}
		catalog.Templates = append(catalog.Templates, pulsarCatalogTemplate{
			URITemplate: template.URITemplate.Raw(),
			Name:        template.Name,
			Description: template.Description,
			MIMEType:    template.MIMEType,
		})
	}

	return catalog
}

func buildPulsarResourceCatalog(
	resources []server.ServerResource,
	templates []server.ServerResourceTemplate,
) pulsarResourceCatalog {
	catalog := newPulsarResourceCatalog()
	catalog.Resources = append(catalog.Resources,
		pulsarCatalogResource{
			URI:         pulsarResourceContextURI,
			Name:        "Pulsar Context",
			Description: "Current Pulsar session connection metadata with authentication material redacted.",
			MIMEType:    pulsarResourceJSONMIMEType,
		},
		pulsarCatalogResource{
			URI:         pulsarResourceCatalogURI,
			Name:        "Pulsar Resource Catalog",
			Description: "Stable catalog of Pulsar MCP resource URIs and URI templates.",
			MIMEType:    pulsarResourceJSONMIMEType,
		},
	)

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

func newPulsarResourceCatalog() pulsarResourceCatalog {
	return pulsarResourceCatalog{
		Version:   1,
		Scheme:    "pulsar",
		Resources: []pulsarCatalogResource{},
		Templates: []pulsarCatalogTemplate{},
		Notes: []string{
			"Resource handlers are read-only and do not return tokens, auth params, key files, TLS private keys, or secret values.",
		},
	}
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
		return pulsarResourceURI{}, fmt.Errorf("pulsar resource URI %q must not include user info", rawURI)
	}
	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return pulsarResourceURI{}, fmt.Errorf("pulsar resource URI %q must not include query or fragment", rawURI)
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
	case len(parts) == 3 && parts[0] == "v2" && parts[1] == "worker" && parts[2] == "cluster":
		return pulsarResourceURI{kind: pulsarResourceKindWorkerCluster}, nil
	case len(parts) == 4 && parts[0] == "v2" && parts[1] == "worker" && parts[2] == "cluster" && parts[3] == "leader":
		return pulsarResourceURI{kind: pulsarResourceKindWorkerLeader}, nil
	case len(parts) == 3 && parts[0] == "v2" && parts[1] == "worker" && parts[2] == "assignments":
		return pulsarResourceURI{kind: pulsarResourceKindWorkerAssignments}, nil
	case len(parts) == 3 && parts[0] == "v2" && parts[1] == "worker-stats" && parts[2] == "functionsmetrics":
		return pulsarResourceURI{kind: pulsarResourceKindWorkerFunctionStats}, nil
	case len(parts) == 3 && parts[0] == "v2" && parts[1] == "worker-stats" && parts[2] == "metrics":
		return pulsarResourceURI{kind: pulsarResourceKindWorkerMetrics}, nil
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
	case len(parts) >= 4 && parts[0] == "v3" && isPulsarWorkloadResourceType(parts[1]):
		return parsePulsarWorkloadAdminResourceURI(rawURI, parts)
	case len(parts) >= 5 && parts[0] == "v3" && parts[1] == "packages":
		return parsePulsarPackageAdminResourceURI(rawURI, parts)
	case len(parts) >= 6 && parts[0] == "v2" && isPulsarResourceTopicDomain(parts[1]):
		return parsePulsarTopicAdminResourceURI(rawURI, parts)
	default:
		return pulsarResourceURI{}, fmt.Errorf("unsupported Pulsar admin resource URI %q", rawURI)
	}
}

func parsePulsarWorkloadAdminResourceURI(rawURI string, parts []string) (pulsarResourceURI, error) {
	resourceType := parts[1]
	tenant := parts[2]
	namespace := parts[3]
	if err := validatePulsarResourcePathSegment("tenant", tenant); err != nil {
		return pulsarResourceURI{}, err
	}
	if err := validatePulsarResourcePathSegment("namespace", namespace); err != nil {
		return pulsarResourceURI{}, err
	}

	base := pulsarResourceURI{
		tenant:    tenant,
		namespace: namespace,
	}
	if len(parts) == 4 {
		switch resourceType {
		case "functions":
			base.kind = pulsarResourceKindFunctions
		case "sources":
			base.kind = pulsarResourceKindSources
		case "sinks":
			base.kind = pulsarResourceKindSinks
		}
		return base, nil
	}
	if len(parts) != 6 {
		return pulsarResourceURI{}, fmt.Errorf("unsupported Pulsar workload resource URI %q", rawURI)
	}

	workload := parts[4]
	suffix := parts[5]
	if err := validatePulsarResourcePathSegment("workload", workload); err != nil {
		return pulsarResourceURI{}, err
	}
	base.workload = workload

	switch resourceType {
	case "functions":
		switch suffix {
		case "metadata":
			base.kind = pulsarResourceKindFunctionMetadata
		case "status":
			base.kind = pulsarResourceKindFunctionStatus
		case "stats":
			base.kind = pulsarResourceKindFunctionStats
		default:
			return pulsarResourceURI{}, fmt.Errorf("unsupported Pulsar function resource URI %q", rawURI)
		}
	case "sources":
		switch suffix {
		case "metadata":
			base.kind = pulsarResourceKindSourceMetadata
		case "status":
			base.kind = pulsarResourceKindSourceStatus
		default:
			return pulsarResourceURI{}, fmt.Errorf("unsupported Pulsar source resource URI %q", rawURI)
		}
	case "sinks":
		switch suffix {
		case "metadata":
			base.kind = pulsarResourceKindSinkMetadata
		case "status":
			base.kind = pulsarResourceKindSinkStatus
		default:
			return pulsarResourceURI{}, fmt.Errorf("unsupported Pulsar sink resource URI %q", rawURI)
		}
	default:
		return pulsarResourceURI{}, fmt.Errorf("unsupported Pulsar workload resource URI %q", rawURI)
	}
	return base, nil
}

func parsePulsarPackageAdminResourceURI(rawURI string, parts []string) (pulsarResourceURI, error) {
	packageType := parts[2]
	tenant := parts[3]
	namespace := parts[4]
	if err := validatePulsarPackageType(packageType); err != nil {
		return pulsarResourceURI{}, err
	}
	if err := validatePulsarResourcePathSegment("tenant", tenant); err != nil {
		return pulsarResourceURI{}, err
	}
	if err := validatePulsarResourcePathSegment("namespace", namespace); err != nil {
		return pulsarResourceURI{}, err
	}

	base := pulsarResourceURI{
		packageType: packageType,
		tenant:      tenant,
		namespace:   namespace,
	}
	switch {
	case len(parts) == 5:
		base.kind = pulsarResourceKindPackages
		return base, nil
	case len(parts) == 7 && parts[6] == "versions":
		packageName := parts[5]
		if err := validatePulsarResourcePathSegment("package", packageName); err != nil {
			return pulsarResourceURI{}, err
		}
		base.kind = pulsarResourceKindPackageVersions
		base.packageName = packageName
		return base, nil
	case len(parts) == 8 && parts[7] == "metadata":
		packageName := parts[5]
		version := parts[6]
		if err := validatePulsarResourcePathSegment("package", packageName); err != nil {
			return pulsarResourceURI{}, err
		}
		if err := validatePulsarResourcePathSegment("version", version); err != nil {
			return pulsarResourceURI{}, err
		}
		base.kind = pulsarResourceKindPackageMetadata
		base.packageName = packageName
		base.versionName = version
		return base, nil
	default:
		return pulsarResourceURI{}, fmt.Errorf("unsupported Pulsar package resource URI %q", rawURI)
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
				return pulsarResourceURI{}, fmt.Errorf("pulsar subscription cursor resource URI %q only supports persistent topics", rawURI)
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

func isPulsarWorkloadResourceType(value string) bool {
	return value == "functions" || value == "sources" || value == "sinks"
}

func validatePulsarPackageType(value string) error {
	if err := validatePulsarResourcePathSegment("package type", value); err != nil {
		return err
	}
	switch value {
	case "function", "source", "sink":
		return nil
	default:
		return fmt.Errorf("unsupported Pulsar package type %q", value)
	}
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
		return fmt.Errorf("pulsar resource URI %s segment must not be empty", name)
	}
	if strings.Contains(value, "/") {
		return fmt.Errorf("pulsar resource URI %s segment must not contain path separators", name)
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
