package mcp

type Feature string

const (
	FeatureAll                          Feature = "all"
	FeatureAllKafka                     Feature = "all-kafka"
	FeatureAllPulsar                    Feature = "all-pulsar"
	FeatureKafkaClient                  Feature = "kafka-client"
	FeatureKafkaAdmin                   Feature = "kafka-admin"
	FeatureKafkaAdminSchemaRegistry     Feature = "kafka-admin-schema-registry"
	FeatureKafkaAdminKafkaConnect       Feature = "kafka-admin-kafka-connect"
	FeaturePulsarAdmin                  Feature = "pulsar-admin"
	FeaturePulsarAdminBrokersStatus     Feature = "pulsar-admin-brokers-status"
	FeaturePulsarAdminBrokers           Feature = "pulsar-admin-brokers"
	FeaturePulsarAdminClusters          Feature = "pulsar-admin-clusters"
	FeaturePulsarAdminFunctions         Feature = "pulsar-admin-functions"
	FeaturePulsarAdminNamespaces        Feature = "pulsar-admin-namespaces"
	FeaturePulsarAdminTenants           Feature = "pulsar-admin-tenants"
	FeaturePulsarAdminTopics            Feature = "pulsar-admin-topics"
	FeaturePulsarAdminFunctionsWorker   Feature = "pulsar-admin-functions-worker"
	FeaturePulsarAdminSinks             Feature = "pulsar-admin-sinks"
	FeaturePulsarAdminSources           Feature = "pulsar-admin-sources"
	FeaturePulsarAdminNamespacePolicy   Feature = "pulsar-admin-namespace-policy"
	FeaturePulsarAdminNsIsolationPolicy Feature = "pulsar-admin-ns-isolation-policy"
	FeaturePulsarAdminPackages          Feature = "pulsar-admin-packages"
	FeaturePulsarAdminResourceQuotas    Feature = "pulsar-admin-resource-quotas"
	FeaturePulsarAdminSchemas           Feature = "pulsar-admin-schemas"
	FeaturePulsarAdminSubscriptions     Feature = "pulsar-admin-subscriptions"
	FeaturePulsarAdminTopicPolicy       Feature = "pulsar-admin-topic-policy"
	FeaturePulsarClient                 Feature = "pulsar-client"
	FeatureStreamNativeCloud            Feature = "streamnative-cloud"
)
