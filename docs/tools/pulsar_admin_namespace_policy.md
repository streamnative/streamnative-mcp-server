#### pulsar_admin_namespace_policy

Tools for managing Pulsar namespace policies.

- **pulsar_admin_namespace_policy_get**: Get configuration policies of a namespace
  - `namespace` (string, required): The namespace name (format: tenant/namespace)
  
- **pulsar_admin_namespace_policy_set**: Set a policy for a namespace
  - `namespace` (string, required): The namespace name (format: tenant/namespace)
  - `policy` (string, required): Type of policy to set, options include:
    - `message-ttl`: Sets time to live for messages
    - `retention`: Sets retention policy for messages
    - `permission`: Grants permissions to a role
    - `replication-clusters`: Sets clusters to replicate to
    - `backlog-quota`: Sets backlog quota policy
    - `topic-auto-creation`: Configures automatic topic creation
    - `schema-validation`: Sets schema validation enforcement
    - `schema-auto-update`: Sets schema auto-update strategy
    - `auto-update-schema`: Controls if schemas can be automatically updated
    - `offload-threshold`: Sets threshold for data offloading
    - `offload-deletion-lag`: Sets time to wait before deleting offloaded data
    - `compaction-threshold`: Sets threshold for topic compaction
    - `max-producers-per-topic`: Limits producers per topic
    - `max-consumers-per-topic`: Limits consumers per topic
    - `max-consumers-per-subscription`: Limits consumers per subscription
    - `anti-affinity-group`: Sets anti-affinity group for isolation
    - `persistence`: Sets persistence configuration
    - `deduplication`: Controls message deduplication
    - `encryption-required`: Controls message encryption
    - `subscription-auth-mode`: Sets subscription auth mode
    - `subscription-permission`: Grants permissions to access a subscription
    - `dispatch-rate`: Sets message dispatch rate
    - `replicator-dispatch-rate`: Sets replicator dispatch rate
    - `subscribe-rate`: Sets subscribe rate per consumer
    - `subscription-dispatch-rate`: Sets subscription dispatch rate
    - `publish-rate`: Sets maximum message publish rate
  - Additional parameters vary based on the policy type

- **pulsar_admin_namespace_policy_remove**: Remove a policy from a namespace
  - `namespace` (string, required): The namespace name (format: tenant/namespace)
  - `policy` (string, required): Type of policy to remove 