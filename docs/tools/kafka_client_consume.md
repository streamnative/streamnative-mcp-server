#### kafka-client-consume

Consume messages from a Kafka topic. This tool allows you to read messages from Kafka topics with various consumption options.

- **kafka_client_consume**
  - **Description**: Read messages from a Kafka topic, with support for consumer groups, offset control, and timeouts. If schema registry integration enabled, and the topic have schema with `topicName-value`, the consume tool will try to use the schema to decode the messages.
  - **Parameters**:
    - `topic` (string, required): The name of the Kafka topic to consume messages from.
    - `group` (string, optional): The consumer group ID to use. If provided, offsets are tracked and committed; otherwise, a random group is used and offsets are not committed.
    - `offset` (string, optional): The offset position to start consuming from. One of:
      - 'atstart': Begin from the earliest available message (default)
      - 'atend': Begin from the next message after the consumer starts
      - 'atcommitted': Begin from the last committed offset (only works with specified 'group')
    - `max-messages` (number, optional): Maximum number of messages to consume in this request. Default: 10
    - `timeout` (number, optional): Maximum time in seconds to wait for messages. Default: 10 