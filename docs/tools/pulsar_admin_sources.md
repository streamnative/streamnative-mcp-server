#### pulsar_admin_sources

Manage Apache Pulsar Sources for data ingestion and integration. Pulsar Sources are connectors that import data from external systems into Pulsar topics. Sources connect to external systems such as databases, messaging platforms, storage services, and real-time data streams to pull data and publish it to Pulsar topics.

This tool provides complete lifecycle management for source connectors:

- **list**: List all sources in a namespace
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)

- **get**: Get source configuration
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The source name

- **status**: Get runtime status of a source (instances, metrics)
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The source name

- **create**: Deploy a new source connector
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The source name (can be provided via `source-config-file`)
  - `destination-topic-name` (string, required): Topic where data will be written
  - Either `archive` or `source-type` must be specified (but not both):
    - `archive` (string): Path to the archive file containing source code
    - `source-type` (string): Built-in connector type to use (e.g., 'kafka', 'jdbc')
  - `source-config-file` (string, optional): YAML file with source configuration
  - `deserialization-classname` (string, optional): SerDe class for the source
  - `schema-type` (string, optional): Schema type for encoding messages (e.g., 'avro', 'json')
  - `classname` (string, optional): Source class name if using custom implementation
  - `processing-guarantees` (string, optional): Delivery semantics ('atleast_once', 'atmost_once', 'effectively_once')
  - `parallelism` (number, optional): Number of instances to run concurrently (default: 1)
  - `cpu` (number, optional): CPU cores per instance
  - `ram` (number, optional): RAM bytes per instance
  - `disk` (number, optional): Disk bytes per instance
  - `source-config` (object, optional): Connector-specific configuration parameters
  - `producer-config` (object, optional): Producer configuration
  - `batch-builder` (string, optional): Batch builder type
  - `batch-source-config` (object, optional): Batch source configuration
  - `custom-runtime-options` (string, optional): Runtime customization options
  - `secrets` (object, optional): Secrets configuration map

- **update**: Update an existing source connector
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The source name (can be provided via `source-config-file`)
  - Parameters similar to `create` operation (all optional during update)
  - `update-auth-data` (boolean, optional): Whether to update auth data

- **delete**: Delete a source
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The source name

- **start**: Start a stopped source
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The source name

- **stop**: Stop a running source
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The source name

- **restart**: Restart a source
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The source name

- **list-built-in**: List all built-in source connectors available in the system
  - No parameters required

Built-in source connectors are available for common systems like Kafka, JDBC, AWS services, and more. Sources follow the tenant/namespace/name hierarchy for organization and access control, can scale through parallelism configuration, and support various processing guarantees. 
