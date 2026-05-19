#### pulsar_admin_cluster_read / pulsar_admin_cluster_write

**Claude connector safety:** Actual MCP tools are split into `pulsar_admin_cluster_read` and `pulsar_admin_cluster_write`. The read tool is read-only and only exposes read operations/parameters. The write tool is destructive and is not registered in read-only mode.

### `pulsar_admin_cluster_read`

Read Pulsar cluster, peer cluster, and failure domain configuration.

- **cluster**
  - **list**: List clusters
  - **get**: Get cluster configuration
    - `cluster_name` (string, required): Cluster name
- **peer_clusters**
  - **get**: Get peer clusters
    - `cluster_name` (string, required): Cluster name
- **failure_domain**
  - **list**: List failure domains in a cluster
    - `cluster_name` (string, required): Cluster name
  - **get**: Get failure domain configuration
    - `cluster_name` (string, required): Cluster name
    - `domain_name` (string, required): Failure domain name

### `pulsar_admin_cluster_write`

Manage clusters, peer clusters, and failure domains.

- **cluster**
  - **create**: Create a cluster
    - `cluster_name` (string, required): Cluster name
    - `service_url` (string, optional): Web service URL
    - `service_url_tls` (string, optional): TLS web service URL
    - `broker_service_url` (string, optional): Broker service URL
    - `broker_service_url_tls` (string, optional): TLS broker service URL
    - `peer_cluster_names` (array, optional): Peer clusters
  - **update**: Update a cluster
    - Same parameters as `create`
  - **delete**: Delete a cluster
    - `cluster_name` (string, required): Cluster name
- **peer_clusters**
  - **update**: Update peer clusters
    - `cluster_name` (string, required): Cluster name
    - `peer_cluster_names` (array, required): Peer cluster names
- **failure_domain**
  - **create**: Create a failure domain
    - `cluster_name` (string, required): Cluster name
    - `domain_name` (string, required): Failure domain name
    - `brokers` (array, required): Brokers in the domain
  - **update**: Update a failure domain
    - Same parameters as `create`
  - **delete**: Delete a failure domain
    - `cluster_name` (string, required): Cluster name
    - `domain_name` (string, required): Failure domain name
