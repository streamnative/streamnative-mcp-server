#### pulsar_admin_topic_read / pulsar_admin_topic_write


<!-- generated:operations:start -->
| Tool | Mode | Operations |
|---|---|---|
| `pulsar_admin_topic_read` | read | `list`, `get`, `get-permissions`, `stats`, `lookup`, `internal-stats`, `internal-info`, `bundle-range`, `last-message-id`, `compact-status`, `offload-status` |
| `pulsar_admin_topic_write` | write | `grant-permissions`, `revoke-permissions`, `create`, `delete`, `unload`, `terminate`, `compact`, `update`, `offload` |
<!-- generated:operations:end -->

**Claude connector safety:** Actual MCP tools are split into `pulsar_admin_topic_read` and `pulsar_admin_topic_write`. The read tool is read-only and only exposes read operations/parameters. The write tool is destructive and is not registered in read-only mode.

### `pulsar_admin_topic_read`

Read topic metadata, permissions, statistics, lookup details, and long-running operation status.

- **topics**
  - **list**: List topics in a namespace
    - `namespace` (string, required): Namespace name in `tenant/namespace` format

- **topic**
  - **get**: Get topic metadata
    - `topic` (string, required): Fully qualified topic name
  - **get-permissions**: Get topic permissions for all roles
    - `topic` (string, required): Fully qualified topic name
  - **stats**: Get topic statistics
    - `topic` (string, required): Fully qualified topic name
    - `partitioned` (boolean, optional): Treat topic as partitioned
    - `per-partition` (boolean, optional): Include per-partition stats
  - **lookup**: Look up the broker serving a topic
    - `topic` (string, required): Fully qualified topic name
  - **internal-stats**: Get topic internal stats
    - `topic` (string, required): Fully qualified topic name
  - **internal-info**: Get topic internal info
    - `topic` (string, required): Fully qualified topic name
  - **bundle-range**: Get topic bundle range
    - `topic` (string, required): Fully qualified topic name
  - **last-message-id**: Get the last message ID of a topic
    - `topic` (string, required): Fully qualified topic name
  - **compact-status**: Get topic compaction status (`status` is accepted as a legacy alias)
    - `topic` (string, required): Fully qualified topic name
    - `wait` (boolean, optional): Poll until compaction status is final
  - **offload-status**: Get topic offload status
    - `topic` (string, required): Fully qualified topic name
    - `wait` (boolean, optional): Poll until offload status is final

### `pulsar_admin_topic_write`

Manage topic lifecycle, permissions, partitioning, compaction, and offload state.

- **topic**
  - **grant-permissions**: Grant topic permissions to a role
    - `topic` (string, required): Fully qualified topic name
    - `role` (string, required): Role to grant permissions to
    - `actions` (array, required): Allowed values include `produce`, `consume`, `sources`, `sinks`, `functions`, `packages`
  - **revoke-permissions**: Revoke all topic permissions from a role
    - `topic` (string, required): Fully qualified topic name
    - `role` (string, required): Role to revoke permissions from
  - **create**: Create a topic
    - `topic` (string, required): Fully qualified topic name
    - `partitions` (number, required): Number of partitions; use `0` for non-partitioned
  - **delete**: Delete a topic
    - `topic` (string, required): Fully qualified topic name
    - `force` (boolean, optional): Force operation even if producers or consumers are active
    - `non-partitioned` (boolean, optional): Delete only a non-partitioned topic with this name
  - **unload**: Unload a topic from broker memory
    - `topic` (string, required): Fully qualified topic name
  - **terminate**: Terminate a topic
    - `topic` (string, required): Fully qualified topic name
  - **compact**: Trigger topic compaction
    - `topic` (string, required): Fully qualified topic name
  - **update**: Update topic partitions or configuration
    - `topic` (string, required): Fully qualified topic name
    - `partitions` (number, optional): New partition count
    - `config` (string, optional): JSON topic configuration
  - **offload**: Offload topic data to long-term storage
    - `topic` (string, required): Fully qualified topic name
    - `messageId` (string, required): Message ID up to which data should be offloaded
