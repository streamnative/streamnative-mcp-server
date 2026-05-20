#### pulsar_admin_nsisolationpolicy_read / pulsar_admin_nsisolationpolicy_write


<!-- generated:operations:start -->
| Tool | Mode | Operations |
|---|---|---|
| `pulsar_admin_nsisolationpolicy_read` | read | `get`, `list` |
| `pulsar_admin_nsisolationpolicy_write` | write | `set`, `delete` |
<!-- generated:operations:end -->

**Claude connector safety:** Actual MCP tools are split into `pulsar_admin_nsisolationpolicy_read` and `pulsar_admin_nsisolationpolicy_write`. The read tool is read-only and only exposes read operations/parameters. The write tool is destructive and is not registered in read-only mode.

Namespace isolation policies control which brokers specific namespaces can use.

### `pulsar_admin_nsisolationpolicy_read`

Read namespace isolation policies and related broker assignments.

- **policy**
  - **get**: Get an isolation policy
    - `cluster` (string, required): Cluster name
    - `name` (string, required): Isolation policy name
  - **list**: List isolation policies in a cluster
    - `cluster` (string, required): Cluster name
- **broker**
  - **get**: Get a broker with its isolation policies
    - `cluster` (string, required): Cluster name
    - `name` (string, required): Broker name
- **brokers**
  - **list**: List brokers with isolation policies
    - `cluster` (string, required): Cluster name

### `pulsar_admin_nsisolationpolicy_write`

Create, update, or delete namespace isolation policies.

- **policy**
  - **set**: Create or update an isolation policy
    - `cluster` (string, required): Cluster name
    - `name` (string, required): Isolation policy name
    - `namespaces` (array, required): Namespaces to apply the policy to
    - `primary` (array, required): Primary brokers
    - `secondary` (array, optional): Secondary brokers
    - `autoFailoverPolicyType` (string, optional): Auto failover policy type
    - `autoFailoverPolicyParams` (object, optional): Auto failover policy parameters
  - **delete**: Delete an isolation policy
    - `cluster` (string, required): Cluster name
    - `name` (string, required): Isolation policy name
