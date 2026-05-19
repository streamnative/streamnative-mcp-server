#### pulsar_admin_topic_policy

**Claude connector safety:** Actual MCP tools: `pulsar_admin_topic_policy_read` (get policy operations) and `pulsar_admin_topic_policy_write` (set/remove policy operations).


Manage Pulsar topic-level policies with operation names aligned to `pulsarctl topics`.

- **Core operations**
  - Retention: `get-retention`, `set-retention`, `remove-retention`
  - Message TTL: `get-message-ttl`, `set-message-ttl`, `remove-message-ttl`
  - Producer and consumer limits:
    - `get-max-producers`, `set-max-producers`, `remove-max-producers`
    - `get-max-consumers`, `set-max-consumers`, `remove-max-consumers`
    - `get-max-unacked-messages-per-consumer`, `set-max-unacked-messages-per-consumer`, `remove-max-unacked-messages-per-consumer`
    - `get-max-unacked-messages-per-subscription`, `set-max-unacked-messages-per-subscription`, `remove-max-unacked-messages-per-subscription`
  - Persistence: `get-persistence`, `set-persistence`, `remove-persistence`
  - Delayed delivery: `get-delayed-delivery`, `set-delayed-delivery`, `remove-delayed-delivery`
  - Dispatch throttling:
    - `get-dispatch-rate`, `set-dispatch-rate`, `remove-dispatch-rate`
    - `get-subscription-dispatch-rate`, `set-subscription-dispatch-rate`, `remove-subscription-dispatch-rate`
    - `get-publish-rate`, `set-publish-rate`, `remove-publish-rate`
  - Topic-level toggles: `get-deduplication`, `set-deduplication`, `remove-deduplication`
  - Storage and cleanup:
    - `get-backlog-quotas`, `set-backlog-quota`, `remove-backlog-quota`
    - `get-compaction-threshold`, `set-compaction-threshold`, `remove-compaction-threshold`
    - `get-inactive-topic-policies`, `set-inactive-topic-policies`, `remove-inactive-topic-policies`
  - Additional MCP-only compatibility operations: `get-subscription-types`, `set-subscription-types`, `remove-subscription-types`

- **Shared parameters**
  - `topic` (string, required): Fully qualified topic name
  - `applied` (boolean, optional): Return the effective inherited policy for `get-retention`, `get-backlog-quotas`, `get-compaction-threshold`, and `get-inactive-topic-policies`

- **Operation-specific parameters**
  - Retention: `retention-time`, `retention-size`
  - Message TTL: `ttl-seconds`
  - Max limits: `count`
  - Persistence: `bookkeeper-ensemble`, `bookkeeper-write-quorum`, `bookkeeper-ack-quorum`, `ml-mark-delete-max-rate`
  - Delayed delivery: `enable`, `disable`, `time`
  - Dispatch throttling: `msg-rate`, `byte-rate`, `period`, `relative-to-publish-rate`
  - Backlog quota: `limit-size`, `limit-time`, `policy`, `type`
  - Inactive topic policies: `delete-while-inactive`, `max-inactive-duration`, `delete-mode`
  - Subscription type restriction: `subscription-types`

Legacy underscore operation aliases from the older MCP implementation are still accepted for backward compatibility.
