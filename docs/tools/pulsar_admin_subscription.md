#### pulsar_admin_subscription

Manage Pulsar topic subscriptions, which represent consumer groups reading from topics.

- **list**: List all subscriptions for a topic
  - `topic` (string, required): The fully qualified topic name
- **create**: Create a new subscription
  - `topic` (string, required): The fully qualified topic name
  - `subscription` (string, required): The subscription name
  - `message_id` (string, optional): Initial position, default is latest
- **delete**: Delete a subscription
  - `topic` (string, required): The fully qualified topic name
  - `subscription` (string, required): The subscription name
- **skip**: Skip messages for a subscription
  - `topic` (string, required): The fully qualified topic name
  - `subscription` (string, required): The subscription name
  - `count` (number, required): Number of messages to skip
- **expire**: Expire messages for a subscription
  - `topic` (string, required): The fully qualified topic name
  - `subscription` (string, required): The subscription name
  - `expiry_time` (string, required): Expiry time in seconds
- **reset-cursor**: Reset subscription position
  - `topic` (string, required): The fully qualified topic name
  - `subscription` (string, required): The subscription name
  - `message_id` (string, required): Message ID to reset to 