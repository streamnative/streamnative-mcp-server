#### kafka-admin-groups

This tool provides access to Kafka consumer group operations including listing, describing, and managing group membership.

- **groups**
  - **list**: List all Kafka Consumer Groups in the cluster
    - _Parameters_: None

- **group**
  - **describe**: Get detailed information about a specific Consumer Group
    - `group` (string, required): The name of the Kafka Consumer Group
  - **remove-members**: Remove specific members from a Consumer Group
    - `group` (string, required): The name of the Kafka Consumer Group
    - `members` (string, required): Comma-separated list of member instance IDs (e.g. "consumer-instance-1,consumer-instance-2")
  - **offsets**: Get offsets for a specific consumer group
    - `group` (string, required): The name of the Kafka Consumer Group
  - **delete-offset**: Delete a specific offset for a consumer group of a topic
    - `group` (string, required): The name of the Kafka Consumer Group
    - `topic` (string, required): The name of the Kafka topic
  - **set-offset**: Set a specific offset for a consumer group's topic-partition
    - `group` (string, required): The name of the Kafka Consumer Group
    - `topic` (string, required): The name of the Kafka topic
    - `partition` (number, required): The partition number
    - `offset` (number, required): The offset value to set (use -1 for earliest, -2 for latest, or a specific value) 