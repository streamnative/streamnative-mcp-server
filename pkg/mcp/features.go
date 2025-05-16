// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

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
	FeaturePulsarFunctions              Feature = "pulsar-functions"
)
