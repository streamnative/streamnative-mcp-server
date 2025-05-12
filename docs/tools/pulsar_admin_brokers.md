#### pulsar_admin_brokers

Unified tool for managing Apache Pulsar broker resources.

- **brokers**
  - **list**: List all active brokers in a cluster
    - `clusterName` (string, required): The cluster name
    
- **health**
  - **get**: Check the health status of a broker
    
- **config**
  - **get**: Get broker configuration
    - `configType` (string, required): Configuration type, available options:
      - `dynamic`: Get list of dynamically modifiable configuration names
      - `runtime`: Get all runtime configurations (including static and dynamic configs)
      - `internal`: Get internal configuration information
      - `all_dynamic`: Get all dynamic configurations and their current values
  - **update**: Update broker configuration
    - `configName` (string, required): Configuration parameter name
    - `configValue` (string, required): Configuration parameter value
  - **delete**: Delete broker configuration
    - `configName` (string, required): Configuration parameter name
    
- **namespaces**
  - **get**: Get namespaces managed by a broker
    - `brokerUrl` (string, required): Broker URL, e.g., '127.0.0.1:8080' 