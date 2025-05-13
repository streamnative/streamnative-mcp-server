#### pulsar_admin_nsisolationpolicy

Manage namespace isolation policies in a Pulsar cluster. Namespace isolation policies enable physical isolation of namespaces by controlling which brokers specific namespaces can use. This helps provide predictable performance and resource isolation, especially in multi-tenant environments.

This tool provides operations across three resource types:

- **policy** (Namespace isolation policy):
  - **get**: Get details of a specific isolation policy
    - `cluster` (string, required): The cluster name
    - `name` (string, required): Name of the isolation policy
  - **list**: List all isolation policies in a cluster
    - `cluster` (string, required): The cluster name
  - **set**: Create or update an isolation policy
    - `cluster` (string, required): The cluster name
    - `name` (string, required): Name of the isolation policy
    - `namespaces` (array, required): List of namespaces to apply the isolation policy
    - `primary` (array, required): List of primary brokers for the namespaces
    - `secondary` (array, optional): List of secondary brokers for the namespaces
    - `autoFailoverPolicyType` (string, optional): Auto failover policy type (e.g., min_available)
    - `autoFailoverPolicyParams` (object, optional): Auto failover policy parameters (e.g., {'min_limit': '1', 'usage_threshold': '100'})
  - **delete**: Delete an isolation policy
    - `cluster` (string, required): The cluster name
    - `name` (string, required): Name of the isolation policy

- **broker** (Broker with isolation policies):
  - **get**: Get details of a specific broker with its isolation policies
    - `cluster` (string, required): The cluster name
    - `name` (string, required): Name of the broker

- **brokers** (All brokers with isolation policies):
  - **list**: List all brokers with their isolation policies
    - `cluster` (string, required): The cluster name 