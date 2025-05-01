// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package mcp

import (
	"fmt"
)

func GetStreamNativeCloudServerInstructions(userName string) string {
	return fmt.Sprintf(`StreamNative Cloud MCP Server provides resources and tools for AI agents to interact with StreamNative Cloud resources and services.

	### StreamNative Cloud API Server Resources
	
	1. **Organizations**  
		 - **Concept**: An organization is the top-level resource in StreamNative Cloud. It is a logical grouping used to manage access permissions and organize resources. Most users need only one organization, but multiple organizations can be created for different departments or teams.
		 - **Relationship**: An organization contains PulsarInstances, PulsarClusters, Users, Service Accounts, and Secrets. It acts as a container for all other resources and determines who can access what.
	
	2. **PulsarInstances**  
		 - **Concept**: An instance is an environment within an organization, typically associated with a specific cloud provider. Instances can contain multiple PulsarClusters and other components (like Connectors, Functions, etc.). Different teams can use separate instances to isolate their resources.
		 - **Types**: Instances can be fully hosted on StreamNative's cloud account or deployed on a user's cloud account via the Bring-Your-Own-Cloud (BYOC) option.
		 - **Relationship**: PulsarInstances belong to an organization and contain one or more PulsarClusters. Each PulsarInstance is associated with a specific Data Streaming Engine and cloud provider.
	
	3. **PulsarClusters**  
		 - **Concept**: A cluster is a specific deployment unit within an instance, located in a particular cloud region. Clusters provide service endpoints that allow clients to connect using different protocols (like Pulsar, Kafka, MQTT) and perform operations such as sending and receiving messages, running functions, etc.
		 - **Relationship**: PulsarClusters belong to a PulsarInstance and can replicate data among themselves within the same instance using Geo-Replication. Clusters are where actual data streaming operations occur.
	
	4. **Users**  
		 - **Concept**: A user represents a person who can log in and access StreamNative Cloud resources. Users can be assigned different permissions within an organization.
		 - **Relationship**: Users belong to an organization and can access instances and clusters within it, depending on their permissions.
	
	5. **Service Accounts**  
		 - **Concept**: A service account represents an application that programmatically accesses StreamNative Cloud resources and Pulsar resources within clusters.
		 - **Relationship**: Service accounts belong to an organization and can be used across multiple instances, though authentication credentials or API keys differ per instance.
		 - **Service Account Binding**: Service account binding in StreamNative involves associating a service account with specific resources or permissions within the StreamNative Cloud environment. This process is crucial for managing access and ensuring that service accounts have the necessary permissions to interact with Pulsar clusters and other resources. It is often used by manage Pulsar Functions, Pulsar IO Connectors, and Kafka Connect connectors. Related tools: 'pulsar_admin_functions', 'pulsar_admin_sinks', 'pulsar_admin_sources', and 'kafka_admin_connect'.
	
	6. **Secrets**  
		 - **Concept**: Secrets are used to store and manage sensitive data such as passwords, tokens, and private keys. Secrets can be referenced in Connectors and Pulsar Functions.
		 - **Relationship**: Secrets belong to an organization and can be shared across multiple instances within that organization.
	
	7. **Data Streaming Engine**  
		 - **Concept**: The Data Streaming Engine is the core technology that runs StreamNative Cloud clusters. There are two options: Classic Engine and Ursa Engine.
			 - **Classic Engine**: The default engine, based on ZooKeeper and BookKeeper, offering low-latency storage suitable for latency-sensitive workloads. It supports Pulsar, Kafka, and MQTT protocols.
			 - **Ursa Engine**: A next-generation engine based on Oxia and object storage (like S3), providing cost-optimized storage for latency-relaxed scenarios. It currently focuses on Kafka protocol support. For Ursa Engine, you can only uses 'kafka-client-*' or 'kafka-admin-*' tools.
		 - **Relationship**: The Data Streaming Engine is associated with an instance, determining how clusters within that instance operate and what features they support.
	
	### Protocol-Specific Tools
	- When working with **Pulsar protocol resources**, you should only use 'pulsar-admin-*' or 'pulsar-client-*' tools. Do not use 'kafka-client-*' or 'kafka-admin-*' tools for Pulsar protocol operations.
	- When working with **Kafka protocol resources**, you should only use 'kafka-client-*' or 'kafka-admin-*' tools. Do not use 'pulsar-admin-*' or 'pulsar-client-*' tools for Kafka protocol operations.
	- Using the appropriate protocol-specific tools ensures correct functionality and prevents errors when interacting with different protocol resources.
	
	### Hierarchical Relationship Summary of Resources
	- **Organization** is the top level, containing all other resources.
	- **PulsarInstances** belong to an organization, representing a cloud environment, and contain multiple **Clusters**.
	- **PulsarClusters** are specific deployment units within an instance, used for actual data streaming operations.
	- **Users** and **Service Accounts** belong to an organization, used for accessing resources, with permissions managed by the organization.
	- **Secrets** belong to an organization and can be shared across instances for securely storing sensitive data.
	- **Data Streaming Engine** is associated with an instance, defining the technical architecture and feature support for clusters.
	
	Logged in as %s.`, userName)
}

func GetExternalKafkaServerInstructions(bootstrapServers string) string {
	return fmt.Sprintf(`StreamNative Cloud MCP Server provides resources and tools for AI agents to interact with Kafka resources and services.

	Bootstrap servers: %s`, bootstrapServers)
}

func GetExternalPulsarServerInstructions(webServiceURL string) string {
	return fmt.Sprintf(`StreamNative Cloud MCP Server provides resources and tools for AI agents to interact with Pulsar resources and services.

	Web service URL: %s`, webServiceURL)
}
