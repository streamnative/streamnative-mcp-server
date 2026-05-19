#### pulsar_admin_sinks_read / pulsar_admin_sinks_write

**Claude connector safety:** Actual MCP tools are split into `pulsar_admin_sinks_read` and `pulsar_admin_sinks_write`. The read tool is read-only and only exposes read operations/parameters. The write tool is destructive and is not registered in read-only mode.

Pulsar Sinks export data from Pulsar topics to external systems.

### `pulsar_admin_sinks_read`

Read sink lists, configuration, runtime status, and built-in sink connector types.

- **list**: List sinks in a namespace
  - `tenant` (string, optional): Tenant name; default `public`
  - `namespace` (string, optional): Namespace name; default `default`
- **get**: Get sink configuration
  - `tenant` (string, optional): Tenant name; default `public`
  - `namespace` (string, optional): Namespace name; default `default`
  - `name` (string, required): Sink name
- **status**: Get sink runtime status
  - `tenant` (string, optional): Tenant name; default `public`
  - `namespace` (string, optional): Namespace name; default `default`
  - `name` (string, required): Sink name
- **list-built-in**: List built-in sink connectors

### `pulsar_admin_sinks_write`

Manage sink connector lifecycle and runtime state.

Common identity parameters:

- `tenant` (string, optional): Tenant name; default `public`
- `namespace` (string, optional): Namespace name; default `default`
- `name` (string, required for operations targeting one sink): Sink name

Operations:

- **create**: Deploy a sink connector
  - Common identity parameters
  - Connector/package parameters include `archive` or `sink-type`, input selection (`inputs` or `topics-pattern`), subscription settings, `classname`, processing guarantees, ordering flags, acknowledgement settings, resources (`cpu`, `ram`, `disk`), schema/serde settings, retry/dead-letter settings, secrets, `sink-config-file`, `sink-config`, transform settings, and runtime options
- **update**: Update sink connector configuration
  - Common identity parameters
  - Same configuration parameters as `create`, plus `update-auth-data`
- **delete**: Delete a sink
  - Common identity parameters
- **start**: Start a stopped sink
  - Common identity parameters
- **stop**: Stop a running sink
  - Common identity parameters
- **restart**: Restart a sink
  - Common identity parameters
