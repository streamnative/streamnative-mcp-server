#### kafka-admin-connect

**Claude connector safety:** Actual MCP tools are split into `kafka_admin_connect_read` and `kafka_admin_connect_write`. The read tool is read-only and only exposes read operations/parameters. The write tool is destructive and is not registered in read-only mode.

### `kafka_admin_connect_read`

Read Kafka Connect cluster, connector, and plugin information.

- **kafka-connect-cluster**
  - **get**: Get Kafka Connect cluster information

- **connectors**
  - **list**: List connectors in the cluster

- **connector**
  - **get**: Get connector details
    - `name` (string, required): Connector name

- **connector-plugins**
  - **list**: List available connector plugins

### `kafka_admin_connect_write`

Manage Kafka Connect connector lifecycle and configuration.

- **connector**
  - **create**: Create a connector
    - `name` (string, required): Connector name
    - `config` (object, required): Connector configuration
  - **update**: Update connector configuration
    - `name` (string, required): Connector name
    - `config` (object, required): Updated connector configuration
  - **delete**: Delete a connector
    - `name` (string, required): Connector name
  - **restart**: Restart a connector
    - `name` (string, required): Connector name
  - **pause**: Pause a connector
    - `name` (string, required): Connector name
  - **resume**: Resume a connector
    - `name` (string, required): Connector name
