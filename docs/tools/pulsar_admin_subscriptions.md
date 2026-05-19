#### pulsar_admin_subscription

**Claude connector safety:** Actual MCP tools are split into `pulsar_admin_subscription_read` and `pulsar_admin_subscription_write`. The read tool is read-only and only exposes read operations/parameters. The write tool is destructive and is not registered in read-only mode.

### `pulsar_admin_subscription_read`

Read subscription names and inspect messages without advancing subscription cursors.

- **subscription**
  - **list**: List subscriptions for a topic
    - `topic` (string, required): Fully qualified topic name
  - **peek**: Peek messages for a subscription without advancing the cursor
    - `topic` (string, required): Fully qualified topic name
    - `subscription` (string, required): Subscription name
    - `count` (number, optional): Number of messages to return; default is 1, maximum is 100
  - **get-message-by-id**: Read a message by ledger ID and entry ID
    - `topic` (string, required): Fully qualified topic name
    - `ledgerId` (number, required): Non-negative ledger ID
    - `entryId` (number, required): Non-negative entry ID

Message payloads returned by `peek` and `get-message-by-id` include `payload`, `payloadBase64`, and `payloadHex`.

### `pulsar_admin_subscription_write`

Manage subscription lifecycle and cursor position.

- **subscription**
  - **create**: Create a subscription
    - `topic` (string, required): Fully qualified topic name
    - `subscription` (string, required): Subscription name
    - `messageId` (string, optional): Initial position, such as `latest`, `earliest`, or `ledgerId:entryId`
  - **delete**: Delete a subscription
    - `topic` (string, required): Fully qualified topic name
    - `subscription` (string, required): Subscription name
    - `force` (boolean, optional): Force delete active consumers
  - **skip**: Skip messages for a subscription
    - `topic` (string, required): Fully qualified topic name
    - `subscription` (string, required): Subscription name
    - `count` (number, required): Number of messages to skip
  - **expire**: Expire messages for a subscription
    - `topic` (string, required): Fully qualified topic name
    - `subscription` (string, required): Subscription name
    - `expireTimeInSeconds` (number, required): Expiry time in seconds
  - **reset-cursor**: Reset subscription cursor position
    - `topic` (string, required): Fully qualified topic name
    - `subscription` (string, required): Subscription name
    - `messageId` (string, required): Message ID to reset to
