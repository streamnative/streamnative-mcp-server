#### pulsar_admin_schema_read / pulsar_admin_schema_write


<!-- generated:operations:start -->
| Tool | Mode | Operations |
|---|---|---|
| `pulsar_admin_schema_read` | read | `get` |
| `pulsar_admin_schema_write` | write | `upload`, `delete` |
<!-- generated:operations:end -->

**Claude connector safety:** Actual MCP tools are split into `pulsar_admin_schema_read` and `pulsar_admin_schema_write`. The read tool is read-only and only exposes read operations/parameters. The write tool is destructive and is not registered in read-only mode.

### `pulsar_admin_schema_read`

Read schema information for a topic.

- **schema**
  - **get**: Get topic schema
    - `topic` (string, required): Fully qualified topic name
    - `version` (number, optional): Schema version

### `pulsar_admin_schema_write`

Upload or delete topic schemas.

- **schema**
  - **upload**: Upload a schema for a topic
    - `topic` (string, required): Fully qualified topic name
    - `filename` (string, required): Path to the schema definition file
  - **delete**: Delete a topic schema
    - `topic` (string, required): Fully qualified topic name
