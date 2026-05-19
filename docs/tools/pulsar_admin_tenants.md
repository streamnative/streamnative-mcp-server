#### pulsar_admin_tenant

**Claude connector safety:** Actual MCP tools: `pulsar_admin_tenant_read` (`list`, `get`) and `pulsar_admin_tenant_write` (`create`, `update`, `delete`).


Manage Pulsar tenants, which are the highest level administrative units.

- **list**: List all tenants in the Pulsar instance
- **get**: Get configuration details for a specific tenant
  - `tenant` (string, required): The tenant name
- **create**: Create a new tenant
  - `tenant` (string, required): The tenant name
  - `admin_roles` (array, optional): List of roles with admin permissions
  - `allowed_clusters` (array, required): List of clusters tenant can access
- **update**: Update configuration for an existing tenant
  - `tenant` (string, required): The tenant name
  - `admin_roles` (array, optional): List of roles with admin permissions
  - `allowed_clusters` (array, required): List of clusters tenant can access
- **delete**: Delete a tenant
  - `tenant` (string, required): The tenant name
