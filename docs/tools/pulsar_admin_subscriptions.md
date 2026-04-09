#### pulsar_admin_subscription

Manage Pulsar topic subscriptions, which represent consumer groups reading from topics.

- **list**: List all subscriptions for a topic
  - `topic` (string, required): The fully qualified topic name
- **create**: Create a new subscription
  - `topic` (string, required): The fully qualified topic name
  - `subscription` (string, required): The subscription name
  - `messageId` (string, optional): Initial position, default is latest
- **delete**: Delete a subscription
  - `topic` (string, required): The fully qualified topic name
  - `subscription` (string, required): The subscription name
  - `force` (boolean, optional): Force delete active consumers
- **skip**: Skip messages for a subscription
  - `topic` (string, required): The fully qualified topic name
  - `subscription` (string, required): The subscription name
  - `count` (number, required): Number of messages to skip
- **expire**: Expire messages for a subscription
  - `topic` (string, required): The fully qualified topic name
  - `subscription` (string, required): The subscription name
  - `expireTimeInSeconds` (number, required): Expiry time in seconds
- **reset-cursor**: Reset subscription position
  - `topic` (string, required): The fully qualified topic name
  - `subscription` (string, required): The subscription name
  - `messageId` (string, required): Message ID to reset to
- **peek**: Peek messages for a subscription without advancing the cursor
  - `topic` (string, required): The fully qualified topic name
  - `subscription` (string, required): The subscription name
  - `count` (number, optional): Number of messages to return, default is 1
- **get-message-by-id**: Read a message by ledger ID and entry ID
  - `topic` (string, required): The fully qualified topic name
  - `ledgerId` (number, required): Ledger ID of the message
  - `entryId` (number, required): Entry ID of the message
