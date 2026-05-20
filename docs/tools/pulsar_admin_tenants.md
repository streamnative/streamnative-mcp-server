#### pulsar_admin_tenant_read / pulsar_admin_tenant_write


<!-- generated:operations:start -->
| Tool | Mode | Operations |
|---|---|---|
| `pulsar_admin_tenant_read` | read | `list`, `get` |
| `pulsar_admin_tenant_write` | write | `create`, `update`, `delete` |
<!-- generated:operations:end -->

**Claude connector safety:** Actual MCP tools are split into `pulsar_admin_tenant_read` and `pulsar_admin_tenant_write`. The read tool is read-only and only exposes read operations/parameters. The write tool is destructive and is not registered in read-only mode.

### `pulsar_admin_tenant_read`

Read tenant names and tenant configuration.

- **tenant**
  - **list**: List tenants
  - **get**: Get tenant configuration
    - `tenant` (string, required): Tenant name

### `pulsar_admin_tenant_write`

Manage tenant lifecycle and configuration.

- **tenant**
  - **create**: Create a tenant
    - `tenant` (string, required): Tenant name
    - `adminRoles` (array, optional): Admin roles
    - `allowedClusters` (array, required): Clusters the tenant can access
  - **update**: Update tenant configuration
    - `tenant` (string, required): Tenant name
    - `adminRoles` (array, optional): Admin roles
    - `allowedClusters` (array, required): Clusters the tenant can access
  - **delete**: Delete a tenant
    - `tenant` (string, required): Tenant name
