#### kafka-admin-partitions

This tool provides access to Kafka partition operations, particularly adding partitions to existing topics.

- **partition**
  - **update**: Update the number of partitions for an existing Kafka topic (can only increase, not decrease)
    - `topic` (string, required): The name of the Kafka topic
    - `new-total` (number, required): The new total number of partitions (must be greater than current) 