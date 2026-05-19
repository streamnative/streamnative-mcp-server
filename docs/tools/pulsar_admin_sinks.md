#### pulsar_admin_sinks

**Claude connector safety:** Actual MCP tools: `pulsar_admin_sinks_read` (`list`, `get`, `status`, `list-built-in`) and `pulsar_admin_sinks_write` (`create`, `update`, `delete`, `start`, `stop`, `restart`).


Manage Apache Pulsar Sinks for data movement and integration. Pulsar Sinks are connectors that export data from Pulsar topics to external systems such as databases, storage services, messaging systems, and third-party applications. Sinks consume messages from one or more Pulsar topics, transform the data if needed, and write it to external systems in a format compatible with the target destination.

This tool provides complete lifecycle management for sink connectors:

- **list**: List all sinks in a namespace
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)

- **get**: Get sink configuration
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The sink name

- **status**: Get runtime status of a sink (instances, metrics)
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The sink name

- **create**: Deploy a new sink connector
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The sink name (can be provided via `sink-config-file`)
  - Either `archive` or `sink-type` must be specified (but not both):
    - `archive` (string): Path to the archive file containing sink code
    - `sink-type` (string): Built-in connector type to use (e.g., 'jdbc', 'elastic-search', 'kafka')
  - Either `inputs` or `topics-pattern` must be specified:
    - `inputs` (array): The sink's input topics (array of strings)
    - `topics-pattern` (string): TopicsPattern to consume from topics matching the pattern (regex)
  - `subs-name` (string, optional): Pulsar subscription name for input topic consumer
  - `subs-position` (string, optional): Subscription position (`Latest` or `Earliest`)
  - `classname` (string, optional): Sink class name for custom archives
  - `processing-guarantees` (string, optional): Delivery semantics
  - `retain-ordering` (boolean, optional): Preserve message ordering
  - `retain-key-ordering` (boolean, optional): Preserve key ordering
  - `auto-ack` (boolean, optional): Auto-ack messages
  - `cleanup-subscription` (boolean, optional): Delete subscription on sink delete (default: true)
  - `parallelism` (number, optional): Number of instances to run concurrently (default: 1)
  - `cpu` / `ram` / `disk` (number, optional): Resource allocation per instance
  - `custom-serde-inputs` (object, optional): Map of input topics to SerDe class names
  - `custom-schema-inputs` (object, optional): Map of input topics to schema type/class
  - `input-specs` (object, optional): Map of input topics to consumer config
  - `max-redeliver-count` (number, optional): Max redeliver attempts
  - `dead-letter-topic` (string, optional): Dead letter topic
  - `timeout-ms` (number, optional): Processing timeout in milliseconds
  - `negative-ack-redelivery-delay-ms` (number, optional): Negative ack redelivery delay
  - `custom-runtime-options` (string, optional): Runtime customization options
  - `secrets` (object, optional): Secrets configuration map
  - `sink-config-file` (string, optional): Path to YAML sink config file
  - `sink-config` (object, optional): Connector-specific configuration parameters
  - `transform-function` (string, optional): Transform function
  - `transform-function-classname` (string, optional): Transform class name
  - `transform-function-config` (string, optional): Transform config

- **update**: Update an existing sink connector
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The sink name (can be provided via `sink-config-file`)
  - Parameters similar to `create` operation (all optional during update)
  - `update-auth-data` (boolean, optional): Update auth data during update

- **delete**: Delete a sink
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The sink name

- **start**: Start a stopped sink
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The sink name

- **stop**: Stop a running sink
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The sink name

- **restart**: Restart a sink
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The sink name

- **list-built-in**: List all built-in sink connectors available in the system
  - No parameters required

Built-in sink connectors are available for common systems like Kafka, JDBC, Elasticsearch, and cloud storage. Sinks follow the tenant/namespace/name hierarchy for organization and access control, can scale through parallelism configuration, and support configurable subscription types. Sinks require proper permissions to access their input topics.
