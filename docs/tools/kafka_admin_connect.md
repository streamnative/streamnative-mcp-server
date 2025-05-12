#### kafka-admin-connect

Kafka Connect is a framework for integrating Kafka with external systems. The following resources and operations are supported:

- **kafka-connect-cluster**
  - **get**: Get information about the Kafka Connect cluster
    - _Parameters_: None

- **connectors**
  - **list**: List all connectors in the cluster
    - _Parameters_: None

- **connector**
  - **get**: Get details of a specific connector
    - `name` (string, required): The connector name
  - **create**: Create a new connector
    - `name` (string, required): The connector name
    - `config` (object, required): Connector configuration
      - Must include at least `connector.class` and other required fields for the connector type
  - **update**: Update an existing connector
    - `name` (string, required): The connector name
    - `config` (object, required): Updated configuration
  - **delete**: Delete a connector
    - `name` (string, required): The connector name
  - **restart**: Restart a connector
    - `name` (string, required): The connector name
  - **pause**: Pause a connector
    - `name` (string, required): The connector name
  - **resume**: Resume a paused connector
    - `name` (string, required): The connector name

- **connector-plugins**
  - **list**: List all available connector plugins
    - _Parameters_: None 