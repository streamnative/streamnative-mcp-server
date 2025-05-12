#### kafka-client-produce

Produce messages to a Kafka topic. This tool allows you to send single or multiple messages with various options.

- **kafka_client_produce**
  - **Description**: Send messages to a Kafka topic, supporting keys, headers, partitions, batching, and file-based payloads.
  - **Parameters**:
    - `topic` (string, required): The name of the Kafka topic to produce messages to.
    - `key` (string, optional): The key for the message. Used for partition assignment and ordering.
    - `value` (string, required if 'messages' is not provided): The value/content of the message to send.
    - `headers` (array, optional): Message headers in the format of [{"key": "header-key", "value": "header-value"}].
    - `sync` (boolean, optional): Whether to wait for server acknowledgment before returning. Default: true. 