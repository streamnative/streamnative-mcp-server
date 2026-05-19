#### pulsar_admin_namespace_read / pulsar_admin_namespace_write

**Claude connector safety:** Actual MCP tools are split into `pulsar_admin_namespace_read` and `pulsar_admin_namespace_write`. The read tool is read-only and only exposes read operations/parameters. The write tool is destructive and is not registered in read-only mode.

### `pulsar_admin_namespace_read`

Read namespace lists and namespace topic lists.

- **list**: List namespaces for a tenant
  - `tenant` (string, required): Tenant name
- **get_topics**: List topics in a namespace
  - `namespace` (string, required): Namespace name in `tenant/namespace` format

### `pulsar_admin_namespace_write`

Manage namespace lifecycle, backlog, subscriptions, unloads, and bundles.

- **create**: Create a namespace
  - `namespace` (string, required): Namespace name in `tenant/namespace` format
  - `bundles` (string, optional): Number of bundles to activate
  - `clusters` (array, optional): Clusters to assign
- **delete**: Delete a namespace
  - `namespace` (string, required): Namespace name in `tenant/namespace` format
- **clear_backlog**: Clear backlog for topics in a namespace
  - `namespace` (string, required): Namespace name in `tenant/namespace` format
  - `subscription` (string, optional): Subscription name
  - `bundle` (string, optional): Bundle name or range
  - `force` (string, optional): Force clear backlog (`true`/`false`)
- **unsubscribe**: Unsubscribe a subscription from topics in a namespace
  - `namespace` (string, required): Namespace name in `tenant/namespace` format
  - `subscription` (string, required): Subscription name
  - `bundle` (string, optional): Bundle name or range
- **unload**: Unload a namespace from the current serving broker
  - `namespace` (string, required): Namespace name in `tenant/namespace` format
  - `bundle` (string, optional): Bundle name or range
- **split_bundle**: Split a namespace bundle
  - `namespace` (string, required): Namespace name in `tenant/namespace` format
  - `bundle` (string, required): Bundle name or range
  - `unload` (string, optional): Unload newly split bundles (`true`/`false`)
