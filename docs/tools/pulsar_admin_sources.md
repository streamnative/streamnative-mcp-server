#### pulsar_admin_sources_read / pulsar_admin_sources_write


<!-- generated:operations:start -->
| Tool | Mode | Operations |
|---|---|---|
| `pulsar_admin_sources_read` | read | `list`, `get`, `status`, `list-built-in` |
| `pulsar_admin_sources_write` | write | `create`, `update`, `delete`, `start`, `stop`, `restart` |
<!-- generated:operations:end -->

**Claude connector safety:** Actual MCP tools are split into `pulsar_admin_sources_read` and `pulsar_admin_sources_write`. The read tool is read-only and only exposes read operations/parameters. The write tool is destructive and is not registered in read-only mode.

Pulsar Sources import data from external systems into Pulsar topics.

### `pulsar_admin_sources_read`

Read source lists, configuration, runtime status, and built-in source connector types.

- **list**: List sources in a namespace
  - `tenant` (string, optional): Tenant name; default `public`
  - `namespace` (string, optional): Namespace name; default `default`
- **get**: Get source configuration
  - `tenant` (string, optional): Tenant name; default `public`
  - `namespace` (string, optional): Namespace name; default `default`
  - `name` (string, required): Source name
- **status**: Get source runtime status
  - `tenant` (string, optional): Tenant name; default `public`
  - `namespace` (string, optional): Namespace name; default `default`
  - `name` (string, required): Source name
- **list-built-in**: List built-in source connectors

### `pulsar_admin_sources_write`

Manage source connector lifecycle and runtime state.

Common identity parameters:

- `tenant` (string, optional): Tenant name; default `public`
- `namespace` (string, optional): Namespace name; default `default`
- `name` (string, required for operations targeting one source): Source name

Operations:

- **create**: Deploy a source connector
  - Common identity parameters
  - Connector/package parameters include `destination-topic-name`, `archive` or `source-type`, `source-config-file`, serialization settings, `classname`, processing guarantees, `parallelism`, resources (`cpu`, `ram`, `disk`), `source-config`, producer config, batch config, secrets, and runtime options
- **update**: Update source connector configuration
  - Common identity parameters
  - Same configuration parameters as `create`, plus `update-auth-data`
- **delete**: Delete a source
  - Common identity parameters
- **start**: Start a stopped source
  - Common identity parameters
- **stop**: Stop a running source
  - Common identity parameters
- **restart**: Restart a source
  - Common identity parameters
