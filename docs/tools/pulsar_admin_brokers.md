#### pulsar_admin_brokers

**Claude connector safety:** Actual MCP tools are split into `pulsar_admin_brokers_read` and `pulsar_admin_brokers_write`. The read tool is read-only and only exposes read operations/parameters. The write tool is destructive and is not registered in read-only mode.

### `pulsar_admin_brokers_read`

Read broker lists, health, configurations, and namespace ownership.

- **brokers**
  - **list**: List active brokers in a cluster
    - `clusterName` (string, required): Cluster name
- **health**
  - **get**: Check broker health status
- **config**
  - **get**: Get broker configuration
    - `configType` (string, required): `dynamic`, `runtime`, `internal`, or `all_dynamic`
- **namespaces**
  - **get**: Get namespaces managed by a broker
    - `clusterName` (string, required): Cluster name
    - `brokerUrl` (string, required): Broker URL, such as `127.0.0.1:8080`

### `pulsar_admin_brokers_write`

Manage broker dynamic configuration values.

- **config**
  - **update**: Update broker configuration
    - `configName` (string, required): Configuration parameter name
    - `configValue` (string, required): Configuration parameter value
  - **delete**: Delete broker configuration
    - `configName` (string, required): Configuration parameter name
