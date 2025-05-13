#### pulsar_admin_resourcequota

Manage Apache Pulsar resource quotas for brokers, namespaces and bundles. Resource quotas define limits for resource usage such as message rates, bandwidth, and memory. These quotas help prevent resource abuse and ensure fair resource allocation across the Pulsar cluster.

This tool provides operations on the following resource:

- **quota** (Resource quota configuration):
  - **get**: Get resource quota for a namespace bundle or the default quota
    - `namespace` (string, optional): The namespace name in format 'tenant/namespace'
    - `bundle` (string, optional): The bundle range in format '{start-boundary}_{end-boundary}'
      - Note: If namespace and bundle are both omitted, returns the default quota
      - Note: If namespace is specified, bundle must also be specified and vice versa
  - **set**: Set resource quota for a namespace bundle or the default quota
    - `namespace` (string, optional): The namespace name in format 'tenant/namespace'
    - `bundle` (string, optional): The bundle range in format '{start-boundary}_{end-boundary}'
    - `msgRateIn` (number, required): Maximum incoming messages per second
    - `msgRateOut` (number, required): Maximum outgoing messages per second
    - `bandwidthIn` (number, required): Maximum inbound bandwidth in bytes per second
    - `bandwidthOut` (number, required): Maximum outbound bandwidth in bytes per second
    - `memory` (number, required): Maximum memory usage in Mbytes
    - `dynamic` (boolean, optional): Whether to allow quota to be dynamically re-calculated
  - **reset**: Reset a namespace bundle's resource quota to default value
    - `namespace` (string, required): The namespace name in format 'tenant/namespace'
    - `bundle` (string, required): The bundle range in format '{start-boundary}_{end-boundary}' 