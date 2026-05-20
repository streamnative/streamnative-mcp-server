#### kafka-admin-groups


<!-- generated:operations:start -->
| Tool | Mode | Operations |
|---|---|---|
| `kafka_admin_groups_read` | read | `list`, `describe`, `offsets` |
| `kafka_admin_groups_write` | write | `remove-members`, `delete-offset`, `set-offset` |
<!-- generated:operations:end -->

**Claude connector safety:** Actual MCP tools are split into `kafka_admin_groups_read` and `kafka_admin_groups_write`. The read tool is read-only and only exposes read operations/parameters. The write tool is destructive and is not registered in read-only mode.

### `kafka_admin_groups_read`

Read Kafka consumer group metadata and committed offsets.

- **groups**
  - **list**: List Kafka consumer groups

- **group**
  - **describe**: Get detailed information about a consumer group
    - `group` (string, required): Consumer group name
  - **offsets**: Get committed offsets for a consumer group
    - `group` (string, required): Consumer group name

### `kafka_admin_groups_write`

Change consumer group membership or committed offsets.

- **group**
  - **remove-members**: Remove specific members from a consumer group
    - `group` (string, required): Consumer group name
    - `members` (string, required): Comma-separated member instance IDs
  - **delete-offset**: Delete offsets for a consumer group topic
    - `group` (string, required): Consumer group name
    - `topic` (string, required): Kafka topic name
  - **set-offset**: Set a consumer group offset for one topic partition
    - `group` (string, required): Consumer group name
    - `topic` (string, required): Kafka topic name
    - `partition` (number, required): Partition number
    - `offset` (number, required): Offset value
