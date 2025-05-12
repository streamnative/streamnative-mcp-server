#### pulsar_admin_topic_policy

Manage Pulsar topic-level policies, which override namespace-level policies.

- **get**: Get policies for a topic
  - `topic` (string, required): The fully qualified topic name
  - `policy_type` (string, required): Type of policy to retrieve
- **set**: Set a policy for a topic
  - `topic` (string, required): The fully qualified topic name
  - `policy_type` (string, required): Type of policy to set
  - `policy_value` (string/object, required): Policy value
- **delete**: Remove a policy from a topic
  - `topic` (string, required): The fully qualified topic name
  - `policy_type` (string, required): Type of policy to remove 