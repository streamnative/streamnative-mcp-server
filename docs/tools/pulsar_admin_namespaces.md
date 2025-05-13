#### pulsar_admin_namespace

Manage Pulsar namespaces with various operations.

- **list**: List all namespaces for a tenant
  - `tenant` (string, required): The tenant name
  
- **get_topics**: Get all topics within a namespace
  - `namespace` (string, required): The namespace name (format: tenant/namespace)
  
- **create**: Create a new namespace
  - `namespace` (string, required): The namespace name (format: tenant/namespace)
  - `bundles` (string, optional): Number of bundles to activate
  - `clusters` (array, optional): List of clusters to assign
  
- **delete**: Delete a namespace
  - `namespace` (string, required): The namespace name (format: tenant/namespace)
  
- **clear_backlog**: Clear backlog for all topics in a namespace
  - `namespace` (string, required): The namespace name (format: tenant/namespace)
  - `subscription` (string, optional): Subscription name
  - `bundle` (string, optional): Bundle name or range
  - `force` (string, optional): Force clear backlog (true/false)
  
- **unsubscribe**: Unsubscribe from a subscription for all topics in a namespace
  - `namespace` (string, required): The namespace name (format: tenant/namespace)
  - `subscription` (string, required): Subscription name
  - `bundle` (string, optional): Bundle name or range
  
- **unload**: Unload a namespace from the current serving broker
  - `namespace` (string, required): The namespace name (format: tenant/namespace)
  - `bundle` (string, optional): Bundle name or range
  
- **split_bundle**: Split a namespace bundle
  - `namespace` (string, required): The namespace name (format: tenant/namespace)
  - `bundle` (string, required): Bundle name or range
  - `unload` (string, optional): Unload newly split bundles (true/false) 