#### pulsar_client_consume

Consume messages from a Pulsar topic. This tool allows you to consume messages from a specified Pulsar topic with various options to control the subscription behavior, message processing, and display format.

- **pulsar_client_consume**
  - `topic` (string, required): Topic to consume from
  - `subscription-name` (string, required): Subscription name
  - `subscription-type` (string, optional): Subscription type (default: exclusive)
    - Options: exclusive, shared, failover, key_shared
  - `subscription-mode` (string, optional): Subscription mode (default: durable)
    - Options: durable, non-durable
  - `initial-position` (string, optional): Initial position (default: latest)
    - Options: latest (consume from the latest message), earliest (consume from the earliest message)
  - `num-messages` (number, optional): Number of messages to consume (0 for unlimited, default: 0)
  - `timeout` (number, optional): Timeout for consuming messages in seconds (default: 30)
  - `show-properties` (boolean, optional): Show message properties (default: false)
  - `hide-payload` (boolean, optional): Hide message payload (default: false) 