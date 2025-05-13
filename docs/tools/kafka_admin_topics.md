#### kafka-admin-topics

This tool provides access to various Kafka topic operations, including creation, deletion, listing, and configuration retrieval.

- **topics**
  - **list**: List all topics in the Kafka cluster
    - `include-internal` (boolean, optional): Whether to include internal Kafka topics (those starting with an underscore). Default: false

- **topic**
  - **get**: Get detailed configuration for a specific topic
    - `name` (string, required): The name of the Kafka topic
  - **create**: Create a new topic
    - `name` (string, required): The name of the Kafka topic
    - `partitions` (number, optional): Number of partitions. Default: 1
    - `replication-factor` (number, optional): Replication factor. Default: 1
    - `configs` (array of string, optional): Topic configuration overrides as key-value strings, e.g. ["cleanup.policy=compact", "retention.ms=604800000"]
  - **delete**: Delete an existing topic
    - `name` (string, required): The name of the Kafka topic 