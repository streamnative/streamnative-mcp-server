#### kafka-admin-topics

**Claude connector safety:** Actual MCP tools are split into `kafka_admin_topics_read` and `kafka_admin_topics_write`. The read tool is read-only and only exposes read operations/parameters. The write tool is destructive and is not registered in read-only mode.

### `kafka_admin_topics_read`

Read Kafka topic lists, configurations, and metadata.

- **topics**
  - **list**: List topics in the Kafka cluster
    - `includeInternal` (boolean, optional): Include internal Kafka topics. Default: false

- **topic**
  - **get**: Get detailed configuration for a topic
    - `name` (string, required): Kafka topic name
  - **metadata**: Get metadata for a topic
    - `name` (string, required): Kafka topic name

### `kafka_admin_topics_write`

Create or delete Kafka topics.

- **topic**
  - **create**: Create a topic
    - `name` (string, required): Kafka topic name
    - `partitions` (number, optional): Number of partitions
    - `replicationFactor` (number, optional): Replication factor
    - `configs` (object, optional): Topic configuration overrides
  - **delete**: Delete a topic
    - `name` (string, required): Kafka topic name
