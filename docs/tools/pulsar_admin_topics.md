#### pulsar_admin_topic

**Claude connector safety:** Actual MCP tools: `pulsar_admin_topic_read` (`list`, `get`, `get-permissions`, stats/lookup/status operations) and `pulsar_admin_topic_write` (permission grants/revokes and mutating topic lifecycle operations).


Manage Apache Pulsar topics. Topics are the core messaging entities in Pulsar that store and transmit messages. Pulsar supports two types of topics: persistent (durable storage with guaranteed delivery) and non-persistent (in-memory with at-most-once delivery). Topics can be partitioned for parallel processing and higher throughput.

- **topic**
  - **get**: Get metadata for a topic
    - `topic` (string, required): The fully qualified topic name
  - **get-permissions**: Get the current topic permissions for every role
    - `topic` (string, required): The fully qualified topic name
  - **grant-permissions**: Grant topic permissions to a role
    - `topic` (string, required): The fully qualified topic name
    - `role` (string, required): The role to grant permissions to
    - `actions` (array of strings, required): Allowed values are `produce`, `consume`, `sources`, `sinks`, `functions`, `packages`
  - **revoke-permissions**: Revoke all topic permissions from a role
    - `topic` (string, required): The fully qualified topic name
    - `role` (string, required): The role to revoke permissions from
  - **create**: Create a new topic with optional partitions
    - `topic` (string, required): The fully qualified topic name
    - `partitions` (number, required): The number of partitions (0 for non-partitioned)
  - **delete**: Delete a topic
    - `topic` (string, required): The fully qualified topic name
    - `force` (boolean, optional): Force operation even if it disrupts producers/consumers
    - `non-partitioned` (boolean, optional): Delete only the non-partitioned topic with the same name
  - **stats**: Get statistics for a topic
    - `topic` (string, required): The fully qualified topic name
    - `partitioned` (boolean, optional): Get stats for a partitioned topic
    - `per-partition` (boolean, optional): Include per-partition stats
  - **lookup**: Look up the broker serving a topic
    - `topic` (string, required): The fully qualified topic name
  - **internal-stats**: Get internal stats for a topic
    - `topic` (string, required): The fully qualified topic name
  - **internal-info**: Get internal info for a topic
    - `topic` (string, required): The fully qualified topic name
  - **bundle-range**: Get the bundle range of a topic
    - `topic` (string, required): The fully qualified topic name
  - **last-message-id**: Get the last message ID of a topic
    - `topic` (string, required): The fully qualified topic name
  - **compact-status**: Get compaction status for a topic (`status` is still accepted as a legacy alias)
    - `topic` (string, required): The fully qualified topic name
    - `wait` (boolean, optional): Poll until the compaction finishes
  - **unload**: Unload a topic from broker memory
    - `topic` (string, required): The fully qualified topic name
  - **terminate**: Terminate a topic (close all producers and mark as inactive)
    - `topic` (string, required): The fully qualified topic name
  - **compact**: Trigger compaction on a topic
    - `topic` (string, required): The fully qualified topic name
  - **update**: Update the number of partitions for a topic
    - `topic` (string, required): The fully qualified topic name
    - `partitions` (number, required): The new number of partitions
  - **offload**: Offload data from a topic to long-term storage
    - `topic` (string, required): The fully qualified topic name
    - `messageId` (string, required): Message ID up to which to offload (format: ledgerId:entryId)
  - **offload-status**: Check the status of data offloading for a topic
    - `topic` (string, required): The fully qualified topic name
    - `wait` (boolean, optional): Poll until the offload finishes

- **topics**
  - **list**: List all topics in a namespace
    - `namespace` (string, required): The namespace name (format: tenant/namespace)
