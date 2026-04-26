# Pulsar Resources

The Pulsar resource surface exposes read-only MCP resources for lightweight cluster context and representative large collections.

## Static resources

- `pulsar://context`: current Pulsar session connection metadata with authentication material redacted.
- `pulsar://resources`: catalog of the registered Pulsar resource URIs and URI templates.
- `pulsar://admin/v2/tenants`: tenant names known to the current Pulsar admin endpoint.
- `pulsar://admin/v2/resource-quotas`: default resource quota for new namespace bundles.
- `pulsar://admin/v2/status`: broker or proxy status for the current Pulsar admin endpoint.
- `pulsar://admin/v2/clusters`: cluster names known to the current Pulsar admin endpoint.
- `pulsar://admin/v2/broker-stats/summary`: bounded summary of broker monitoring metrics and load report.

All static resources return `application/json`.

## Resource templates

- `pulsar://admin/v2/tenants/{tenant}`: gets tenant configuration.
- `pulsar://admin/v2/tenants/{tenant}/namespaces`: lists namespaces for a tenant.
- `pulsar://admin/v2/namespaces/{tenant}/{namespace}`: gets namespace policies.
- `pulsar://admin/v2/namespaces/{tenant}/{namespace}/topics`: lists topics for a namespace.
- `pulsar://admin/v2/resource-quotas/{tenant}/{namespace}/{bundle}`: gets resource quota for a namespace bundle.
- `pulsar://admin/v2/clusters/{cluster}`: sanitized configuration for a cluster.
- `pulsar://admin/v2/brokers/{cluster}`: lists active brokers for a cluster.
- `pulsar://admin/v2/clusters/{cluster}/failureDomains`: lists failure domains for a cluster.
- `pulsar://admin/v2/clusters/{cluster}/failureDomains/{domain}`: gets a failure domain.
- `pulsar://admin/v2/clusters/{cluster}/namespaceIsolationPolicies`: lists namespace isolation policies for a cluster.
- `pulsar://admin/v2/clusters/{cluster}/namespaceIsolationPolicies/{policy}`: gets a namespace isolation policy.

Template reads return `application/json` and require a Pulsar session in the request context.

The resource list is feature-gated. Tenant and namespace resources require the matching Pulsar admin feature such as `pulsar-admin-tenants`, `pulsar-admin-namespaces`, `pulsar-admin-namespace-policy`, `pulsar-admin-topics`, or `pulsar-admin-resource-quotas`. Cluster resources require the matching Pulsar admin feature such as `pulsar-admin-clusters`, `pulsar-admin-brokers`, `pulsar-admin-broker-stats`, `pulsar-admin-brokers-status`, or `pulsar-admin-ns-isolation-policy`. All resources are also enabled by one of `pulsar-admin`, `all-pulsar`, or `all`.

## Safety

Resource handlers are read-only. They do not consume messages, commit cursors, clear backlog, unload topics, split bundles, delete resources, start workloads, or stop workloads. They also do not return tokens, auth params, key files, TLS private keys, or secret values.
