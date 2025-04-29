package mcp

type McpFeature string

const (
	FeatureAll                          McpFeature = "all"
	FeatureAllKafka                     McpFeature = "all-kafka"
	FeatureAllPulsar                    McpFeature = "all-pulsar"
	FeatureKafkaClient                  McpFeature = "kafka-client"
	FeatureKafkaAdmin                   McpFeature = "kafka-admin"
	FeatureKafkaAdminSchemaRegistry     McpFeature = "kafka-admin-schema-registry"
	FeatureKafkaAdminKafkaConnect       McpFeature = "kafka-admin-kafka-connect"
	FeaturePulsarAdmin                  McpFeature = "pulsar-admin"
	FeaturePulsarAdminBrokersStatus     McpFeature = "pulsar-admin-brokers-status"
	FeaturePulsarAdminBrokers           McpFeature = "pulsar-admin-brokers"
	FeaturePulsarAdminClusters          McpFeature = "pulsar-admin-clusters"
	FeaturePulsarAdminFunctions         McpFeature = "pulsar-admin-functions"
	FeaturePulsarAdminNamespaces        McpFeature = "pulsar-admin-namespaces"
	FeaturePulsarAdminTenants           McpFeature = "pulsar-admin-tenants"
	FeaturePulsarAdminTopics            McpFeature = "pulsar-admin-topics"
	FeaturePulsarAdminFunctionsWorker   McpFeature = "pulsar-admin-functions-worker"
	FeaturePulsarAdminSinks             McpFeature = "pulsar-admin-sinks"
	FeaturePulsarAdminSources           McpFeature = "pulsar-admin-sources"
	FeaturePulsarAdminNamespacePolicy   McpFeature = "pulsar-admin-namespace-policy"
	FeaturePulsarAdminNsIsolationPolicy McpFeature = "pulsar-admin-ns-isolation-policy"
	FeaturePulsarAdminPackages          McpFeature = "pulsar-admin-packages"
	FeaturePulsarAdminResourceQuotas    McpFeature = "pulsar-admin-resource-quotas"
	FeaturePulsarAdminSchemas           McpFeature = "pulsar-admin-schemas"
	FeaturePulsarAdminSubscriptions     McpFeature = "pulsar-admin-subscriptions"
	FeaturePulsarAdminTopicPolicy       McpFeature = "pulsar-admin-topic-policy"
	FeaturePulsarClient                 McpFeature = "pulsar-client"
	FeatureStreamNativeCloud            McpFeature = "streamnative-cloud"
)
