#### pulsar_admin_topic_policy

**Claude connector safety:** Actual MCP tools are split into `pulsar_admin_topic_policy_read` and `pulsar_admin_topic_policy_write`. The read tool is read-only and only exposes get operations/parameters. The write tool is destructive and is not registered in read-only mode.

### `pulsar_admin_topic_policy_read`

Read topic-level policies.

Read operations:

- `get-retention`
- `get-message-ttl`
- `get-max-producers`
- `get-max-consumers`
- `get-max-unacked-messages-per-consumer`
- `get-max-unacked-messages-per-subscription`
- `get-persistence`
- `get-delayed-delivery`
- `get-dispatch-rate`
- `get-subscription-dispatch-rate`
- `get-deduplication`
- `get-backlog-quotas`
- `get-compaction-threshold`
- `get-publish-rate`
- `get-inactive-topic-policies`
- `get-subscription-types`

Read parameters:

- `topic` (string, required): Fully qualified topic name
- `applied` (boolean, optional): Return effective inherited policy where supported
- `type` (string, optional): Backlog quota type for backlog quota reads

### `pulsar_admin_topic_policy_write`

Set or remove topic-level policies.

Write operations:

- `set-retention`, `remove-retention`
- `set-message-ttl`, `remove-message-ttl`
- `set-max-producers`, `remove-max-producers`
- `set-max-consumers`, `remove-max-consumers`
- `set-max-unacked-messages-per-consumer`, `remove-max-unacked-messages-per-consumer`
- `set-max-unacked-messages-per-subscription`, `remove-max-unacked-messages-per-subscription`
- `set-persistence`, `remove-persistence`
- `set-delayed-delivery`, `remove-delayed-delivery`
- `set-dispatch-rate`, `remove-dispatch-rate`
- `set-subscription-dispatch-rate`, `remove-subscription-dispatch-rate`
- `set-deduplication`, `remove-deduplication`
- `set-backlog-quota`, `remove-backlog-quota`
- `set-compaction-threshold`, `remove-compaction-threshold`
- `set-publish-rate`, `remove-publish-rate`
- `set-inactive-topic-policies`, `remove-inactive-topic-policies`
- `set-subscription-types`, `remove-subscription-types`

Write parameters:

- `topic` (string, required): Fully qualified topic name
- Retention: `retention-time`, `retention-size`
- Message TTL: `ttl-seconds`
- Max limits: `count`
- Persistence: `bookkeeper-ensemble`, `bookkeeper-write-quorum`, `bookkeeper-ack-quorum`, `ml-mark-delete-max-rate`
- Delayed delivery: `enable`, `disable`, `time`
- Dispatch/publish throttling: `msg-rate`, `byte-rate`, `period`, `relative-to-publish-rate`
- Backlog quota: `limit-size`, `limit-time`, `policy`, `type`
- Inactive topic policies: `delete-while-inactive`, `max-inactive-duration`, `delete-mode`
- Subscription type restriction: `subscription-types`

Legacy underscore operation aliases from the older MCP implementation are still accepted by the handlers.
