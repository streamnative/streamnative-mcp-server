#### pulsar_admin_resourcequota

**Claude connector safety:** Actual MCP tools are split into `pulsar_admin_resourcequota_read` and `pulsar_admin_resourcequota_write`. The read tool is read-only and only exposes read operations/parameters. The write tool is destructive and is not registered in read-only mode.

### `pulsar_admin_resourcequota_read`

Read resource quota configuration for default quotas or namespace bundles.

- **quota**
  - **get**: Get a resource quota
    - `namespace` (string, optional): Namespace name in `tenant/namespace` format
    - `bundle` (string, optional): Bundle range in `{start-boundary}_{end-boundary}` format

If `namespace` and `bundle` are omitted, the default quota is returned. If one of `namespace` or `bundle` is specified, the other must also be specified.

### `pulsar_admin_resourcequota_write`

Set or reset resource quota configuration.

- **quota**
  - **set**: Set a resource quota
    - `namespace` (string, optional): Namespace name in `tenant/namespace` format
    - `bundle` (string, optional): Bundle range in `{start-boundary}_{end-boundary}` format
    - `msgRateIn` (number, required): Incoming messages per second
    - `msgRateOut` (number, required): Outgoing messages per second
    - `bandwidthIn` (number, required): Inbound bandwidth in bytes per second
    - `bandwidthOut` (number, required): Outbound bandwidth in bytes per second
    - `memory` (number, required): Memory usage in Mbytes
    - `dynamic` (boolean, optional): Allow dynamic recalculation
  - **reset**: Reset a namespace bundle resource quota to the default
    - `namespace` (string, required): Namespace name in `tenant/namespace` format
    - `bundle` (string, required): Bundle range in `{start-boundary}_{end-boundary}` format
