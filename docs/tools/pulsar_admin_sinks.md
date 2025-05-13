#### pulsar_admin_sinks

Manage Apache Pulsar Sinks for data movement and integration. Pulsar Sinks are connectors that export data from Pulsar topics to external systems such as databases, storage services, messaging systems, and third-party applications. Sinks consume messages from one or more Pulsar topics, transform the data if needed, and write it to external systems in a format compatible with the target destination.

This tool provides complete lifecycle management for sink connectors:

- **list**: List all sinks in a namespace
  - `tenant` (string, required): The tenant name
  - `namespace` (string, required): The namespace name

- **get**: Get sink configuration
  - `tenant` (string, required): The tenant name
  - `namespace` (string, required): The namespace name
  - `name` (string, required): The sink name

- **status**: Get runtime status of a sink (instances, metrics)
  - `tenant` (string, required): The tenant name
  - `namespace` (string, required): The namespace name
  - `name` (string, required): The sink name

- **create**: Deploy a new sink connector
  - `tenant` (string, required): The tenant name
  - `namespace` (string, required): The namespace name
  - `name` (string, required): The sink name
  - Either `archive` or `sink-type` must be specified (but not both):
    - `archive` (string): Path to the archive file containing sink code
    - `sink-type` (string): Built-in connector type to use (e.g., 'jdbc', 'elastic-search', 'kafka')
  - Either `inputs` or `topics-pattern` must be specified:
    - `inputs` (array): The sink's input topics (array of strings)
    - `topics-pattern` (string): TopicsPattern to consume from topics matching the pattern (regex)
  - `subs-name` (string, optional): Pulsar subscription name for input topic consumer
  - `parallelism` (number, optional): Number of instances to run concurrently (default: 1)
  - `sink-config` (object, optional): Connector-specific configuration parameters

- **update**: Update an existing sink connector
  - `tenant` (string, required): The tenant name
  - `namespace` (string, required): The namespace name
  - `name` (string, required): The sink name
  - Parameters similar to `create` operation (all optional during update)

- **delete**: Delete a sink
  - `tenant` (string, required): The tenant name
  - `namespace` (string, required): The namespace name
  - `name` (string, required): The sink name

- **start**: Start a stopped sink
  - `tenant` (string, required): The tenant name
  - `namespace` (string, required): The namespace name
  - `name` (string, required): The sink name

- **stop**: Stop a running sink
  - `tenant` (string, required): The tenant name
  - `namespace` (string, required): The namespace name
  - `name` (string, required): The sink name

- **restart**: Restart a sink
  - `tenant` (string, required): The tenant name
  - `namespace` (string, required): The namespace name
  - `name` (string, required): The sink name

- **list-built-in**: List all built-in sink connectors available in the system
  - No parameters required

Built-in sink connectors are available for common systems like Kafka, JDBC, Elasticsearch, and cloud storage. Sinks follow the tenant/namespace/name hierarchy for organization and access control, can scale through parallelism configuration, and support configurable subscription types. Sinks require proper permissions to access their input topics. 