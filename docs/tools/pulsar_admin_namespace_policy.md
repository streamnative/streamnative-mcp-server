#### pulsar_admin_namespace_policy

Tools for managing Pulsar namespace policies.

- **pulsar_admin_namespace_policy_get**: Get the full policy set for a namespace
  - `namespace` (string, required): Namespace in `tenant/namespace` format

- **pulsar_admin_namespace_policy_set**: Set a namespace policy
  - `namespace` (string, required): Namespace in `tenant/namespace` format
  - `policy` (string, required): One of:
    - `message-ttl`
    - `retention`
    - `permission`
    - `replication-clusters`
    - `backlog-quota`
    - `topic-auto-creation`
    - `schema-validation`
    - `schema-auto-update`
    - `auto-update-schema`
    - `offload-threshold`
    - `offload-deletion-lag`
    - `compaction-threshold`
    - `max-producers-per-topic`
    - `max-consumers-per-topic`
    - `max-consumers-per-subscription`
    - `anti-affinity-group`
    - `persistence`
    - `deduplication`
    - `encryption-required`
    - `subscription-auth-mode`
    - `subscription-permission`
    - `dispatch-rate`
    - `replicator-dispatch-rate`
    - `subscribe-rate`
    - `subscription-dispatch-rate`
    - `publish-rate`
  - Common policy-specific parameters:
    - `ttl` (string): For `message-ttl`
    - `time`, `size` (string): For `retention`
    - `role`, `actions` (array): For `permission`
    - `clusters` (array): For `replication-clusters`
    - `limit-size`, `limit-time`, `backlog-policy`, `type`: For `backlog-quota`
    - `enabled`, `topic-type`, `partitions`: For `topic-auto-creation`
    - `enabled`: For `schema-validation`, `auto-update-schema`, `deduplication`, `encryption-required`
    - `compatibility`: For `schema-auto-update`
    - `lag`: For `offload-deletion-lag`
    - `count`: For `max-*` policies
    - `group`: For `anti-affinity-group`
    - `ensemble-size`, `write-quorum-size`, `ack-quorum-size`, `ml-mark-delete-max-rate`: For `persistence`
    - `mode`: For `subscription-auth-mode`
    - `subscription`, `roles`: For `subscription-permission`
    - `msg-rate`, `byte-rate`, `period`: For dispatch and publish rate policies
    - `subscribe-rate`, `period`: For `subscribe-rate`

- **pulsar_admin_namespace_policy_remove**: Remove a namespace policy override
  - `namespace` (string, required): Namespace in `tenant/namespace` format
  - `policy` (string, required): One of `backlog-quota`, `topic-auto-creation`, `offload-deletion-lag`, `anti-affinity-group`, `permission`, `subscription-permission`
  - `role` (string): Required for `permission` and `subscription-permission`
  - `subscription` (string): Required for `subscription-permission`
