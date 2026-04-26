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
- `pulsar://admin/v2/worker/cluster`: bounded summary of Pulsar Functions workers.
- `pulsar://admin/v2/worker/cluster/leader`: current Pulsar Functions worker leader.
- `pulsar://admin/v2/worker/assignments`: bounded summary of Pulsar Functions worker assignments.
- `pulsar://admin/v2/worker-stats/functionsmetrics`: bounded function instance stats reported by the Pulsar Functions worker.
- `pulsar://admin/v2/worker-stats/metrics`: bounded summary of Pulsar Functions worker monitoring metrics.

All static resources return `application/json`.

## Resource templates

- `pulsar://admin/v2/tenants/{tenant}`: gets tenant configuration.
- `pulsar://admin/v2/tenants/{tenant}/namespaces`: lists namespaces for a tenant.
- `pulsar://admin/v2/namespaces/{tenant}/{namespace}`: gets namespace policies.
- `pulsar://admin/v2/namespaces/{tenant}/{namespace}/topics`: lists topics for a namespace.
- `pulsar://admin/v2/resource-quotas/{tenant}/{namespace}/{bundle}`: gets resource quota for a namespace bundle.
- `pulsar://admin/v2/{domain}/{tenant}/{namespace}/{topic}/metadata`: gets parsed topic identity and sanitized topic properties.
- `pulsar://admin/v2/{domain}/{tenant}/{namespace}/{topic}/stats`: gets a bounded topic statistics summary without publisher or consumer details.
- `pulsar://admin/v2/{domain}/{tenant}/{namespace}/{topic}/partitions`: gets topic partition metadata.
- `pulsar://admin/v2/{domain}/{tenant}/{namespace}/{topic}/policies/{policy}`: gets one read-only topic policy value. Supported policies are `retention`, `message-ttl`, `max-producers`, `max-consumers`, `max-unacked-messages-per-consumer`, `max-unacked-messages-per-subscription`, `persistence`, `delayed-delivery`, `dispatch-rate`, `subscription-dispatch-rate`, `deduplication`, `backlog-quotas`, `compaction-threshold`, `publish-rate`, and `inactive-topic-policies`.
- `pulsar://admin/v2/{domain}/{tenant}/{namespace}/{topic}/schema`: gets the latest topic schema and version.
- `pulsar://admin/v2/{domain}/{tenant}/{namespace}/{topic}/schema/{version}`: gets a specific topic schema version.
- `pulsar://admin/v2/{domain}/{tenant}/{namespace}/{topic}/subscriptions`: lists subscriptions for a topic.
- `pulsar://admin/v2/{domain}/{tenant}/{namespace}/{topic}/subscriptions/{subscription}/stats`: gets bounded subscription statistics without consumer details.
- `pulsar://admin/v2/{domain}/{tenant}/{namespace}/{topic}/subscriptions/{subscription}/backlog`: gets subscription backlog counters without changing cursor state.
- `pulsar://admin/v2/persistent/{tenant}/{namespace}/{topic}/subscriptions/{subscription}/cursor`: gets persistent topic cursor positions for a subscription.
- `pulsar://admin/v3/functions/{tenant}/{namespace}`: lists Pulsar Functions for a namespace, bounded to the first 50 names.
- `pulsar://admin/v3/functions/{tenant}/{namespace}/{function}/metadata`: gets sanitized Pulsar Function metadata.
- `pulsar://admin/v3/functions/{tenant}/{namespace}/{function}/status`: gets bounded Pulsar Function runtime status without exception detail strings.
- `pulsar://admin/v3/functions/{tenant}/{namespace}/{function}/stats`: gets bounded Pulsar Function statistics and user metric names.
- `pulsar://admin/v3/sources/{tenant}/{namespace}`: lists Pulsar Sources for a namespace, bounded to the first 50 names.
- `pulsar://admin/v3/sources/{tenant}/{namespace}/{source}/metadata`: gets sanitized Pulsar Source metadata.
- `pulsar://admin/v3/sources/{tenant}/{namespace}/{source}/status`: gets bounded Pulsar Source runtime status without exception detail strings.
- `pulsar://admin/v3/sinks/{tenant}/{namespace}`: lists Pulsar Sinks for a namespace, bounded to the first 50 names.
- `pulsar://admin/v3/sinks/{tenant}/{namespace}/{sink}/metadata`: gets sanitized Pulsar Sink metadata.
- `pulsar://admin/v3/sinks/{tenant}/{namespace}/{sink}/status`: gets bounded Pulsar Sink runtime status without exception detail strings.
- `pulsar://admin/v3/packages/{type}/{tenant}/{namespace}`: lists packages by type and namespace, bounded to the first 50 names. Supported package types are `function`, `source`, and `sink`.
- `pulsar://admin/v3/packages/{type}/{tenant}/{namespace}/{package}/versions`: lists package versions, bounded to the first 50 names.
- `pulsar://admin/v3/packages/{type}/{tenant}/{namespace}/{package}/{version}/metadata`: gets sanitized metadata for one package version.
- `pulsar://admin/v2/clusters/{cluster}`: sanitized configuration for a cluster.
- `pulsar://admin/v2/brokers/{cluster}`: lists active brokers for a cluster.
- `pulsar://admin/v2/clusters/{cluster}/failureDomains`: lists failure domains for a cluster.
- `pulsar://admin/v2/clusters/{cluster}/failureDomains/{domain}`: gets a failure domain.
- `pulsar://admin/v2/clusters/{cluster}/namespaceIsolationPolicies`: lists namespace isolation policies for a cluster.
- `pulsar://admin/v2/clusters/{cluster}/namespaceIsolationPolicies/{policy}`: gets a namespace isolation policy.

Template reads return `application/json` and require a Pulsar session in the request context. Topic templates accept `domain` values of `persistent` or `non-persistent`; `topic` is the local topic name path segment. Subscription cursor resources are persistent-only because they are backed by topic internal stats. Workload and package templates use Pulsar admin v3 APIs; functions worker resources use the current Pulsar admin v2 worker endpoints.

The resource list is feature-gated. Tenant and namespace resources require the matching Pulsar admin feature such as `pulsar-admin-tenants`, `pulsar-admin-namespaces`, `pulsar-admin-namespace-policy`, `pulsar-admin-topics`, or `pulsar-admin-resource-quotas`. Topic metadata, stats, and partition metadata resources require `pulsar-admin-topics`; topic policy resources require `pulsar-admin-topic-policy`; schema resources require `pulsar-admin-schemas`; subscription resources require `pulsar-admin-subscriptions`. Workload resources require the matching feature such as `pulsar-admin-functions`, `pulsar-admin-sources`, `pulsar-admin-sinks`, `pulsar-admin-packages`, or `pulsar-admin-functions-worker`. Cluster resources require the matching Pulsar admin feature such as `pulsar-admin-clusters`, `pulsar-admin-brokers`, `pulsar-admin-broker-stats`, `pulsar-admin-brokers-status`, or `pulsar-admin-ns-isolation-policy`. All resources are also enabled by one of `pulsar-admin`, `all-pulsar`, or `all`.

## Safety

Resource handlers are read-only. They do not consume messages, commit cursors, clear backlog, unload topics, split bundles, delete resources, start workloads, or stop workloads. They also do not return tokens, auth params, key files, TLS private keys, or secret values.
