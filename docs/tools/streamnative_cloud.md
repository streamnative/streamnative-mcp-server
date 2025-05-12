#### sncloud_context_available_clusters

Display the available Pulsar clusters for the current organization on StreamNative Cloud. This information helps you select the appropriate cluster when setting your context.

- **sncloud_context_available_clusters**
  - No parameters required

You can use `sncloud_context_use_cluster` to change the context to a specific cluster. You will need to ask for user confirmation of the target context cluster if there are multiple clusters available. 

---

#### sncloud_context_use_cluster

Set the current context to a specific StreamNative Cloud cluster. Once you set the context, you can use Pulsar and Kafka tools to interact with the cluster.

- **sncloud_context_use_cluster**
  - `instanceName` (string, required): The name of the Pulsar instance to use
  - `clusterName` (string, required): The name of the Pulsar cluster to use

If you encounter `ContextNotSetErr`, please use `sncloud_context_available_clusters` to list the available clusters and set the context to a specific cluster. 

---

#### sncloud_context_whoami

Display the currently logged-in service account. Returns the name of the authenticated service account and the organization.

- **sncloud_context_whoami**
  - No parameters required

This tool returns a JSON object containing the service account name and organization. 

---

#### sncloud_logs

Display logs of resources in StreamNative Cloud, including Pulsar functions, source connectors, sink connectors, and Kafka Connect connectors. This tool helps debug issues with resources in StreamNative Cloud and works with the current context cluster.

- **sncloud_logs**
  - `component` (string, required): The component type to get logs from
    - Options: sink, source, function, kafka-connect
  - `name` (string, required): The name of the resource to get logs from
  - `tenant` (string, required): The Pulsar tenant of the resource (default: "public")
    - Required for Pulsar functions, sources, and sinks
    - Optional for Kafka Connect connectors
  - `namespace` (string, required): The Pulsar namespace of the resource (default: "default")
    - Required for Pulsar functions, sources, and sinks
    - Optional for Kafka Connect connectors
  - `size` (string, optional): Number of log lines to retrieve (default: "20")
  - `replica_id` (number, optional): The replica index for resources with multiple replicas (default: -1, which means all replicas)
  - `timestamp` (string, optional): Start timestamp of logs in milliseconds (e.g., "1662430984225")
  - `since` (string, optional): Retrieve logs from a relative time in the past (e.g., "1h" for one hour ago)
  - `previous_container` (boolean, optional): Return logs from previously terminated container (default: false)

---

#### sncloud_resources_apply

Apply (create or update) StreamNative Cloud resources from JSON definitions. This tool can be used to manage infrastructure resources such as PulsarInstances and PulsarClusters in your StreamNative Cloud organization.

- **sncloud_resources_apply**
  - `json_content` (string, required): The JSON content to apply, defining the resource according to the StreamNative Cloud API schema
  - `dry_run` (boolean, optional): If true, only validate the resource without applying it to the server (default: false)

Supported resource types:
- PulsarInstance (apiVersion: cloud.streamnative.io/v1alpha1)
- PulsarCluster (apiVersion: cloud.streamnative.io/v1alpha1)

---

#### sncloud_resources_delete

Delete StreamNative Cloud resources. This tool removes resources from your StreamNative Cloud organization.

- **sncloud_resources_delete**
  - `name` (string, required): The name of the resource to delete
  - `type` (string, required): The type of the resource to delete
    - Options: PulsarInstance, PulsarCluster

This is a destructive operation that cannot be undone. Use with caution. 