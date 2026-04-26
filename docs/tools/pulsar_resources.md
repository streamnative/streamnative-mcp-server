# Pulsar Resources

The Pulsar resource surface exposes read-only MCP resources for lightweight cluster context and representative large collections.

## Static resources

- `pulsar://context`: current Pulsar session connection metadata with authentication material redacted.
- `pulsar://resources`: catalog of the registered Pulsar resource URIs and URI templates.

Both resources return `application/json`.

## Resource templates

- `pulsar://admin/v2/tenants/{tenant}/namespaces`: lists namespaces for a tenant.
- `pulsar://admin/v2/namespaces/{tenant}/{namespace}/topics`: lists topics for a namespace.

Template reads return `application/json` and require a Pulsar session in the request context.

## Safety

Resource handlers are read-only. They do not consume messages, commit cursors, clear backlog, unload topics, split bundles, delete resources, start workloads, or stop workloads. They also do not return tokens, auth params, key files, TLS private keys, or secret values.
