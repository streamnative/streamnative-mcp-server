# StreamNative MCP Server

A Model Context Protocol (MCP) server for integrating AI agents with StreamNative Cloud resources and Apache Kafka/Pulsar messaging systems.

## Overview

StreamNative MCP Server provides a standard interface for LLMs (Large Language Models) and AI agents to interact with StreamNative Cloud services, Apache Kafka, and Apache Pulsar. This implementation follows the [Model Context Protocol](https://modelcontextprotocol.io/introduction) specification, enabling AI applications to access messaging services through a standardized interface.

## Features

- **StreamNative Cloud Integration**: 
  - Connect to StreamNative Cloud resources with authentication
  - Switch to clusters available in your organization
  - Describe the status of clusters resources
- **Apache Kafka Support**: Interact with Apache Kafka resources including:
  - Kafka Admin operations (topics, partitions, consumer groups)
  - Schema Registry operations
  - Kafka Connect operations (*)
  - Kafka Client operations (producers, consumers)
- **Apache Pulsar Support**: Interact with Apache Pulsar resources including:
  - Pulsar Admin operations (topics, namespaces, tenants, schemas, etc.)
  - Pulsar Client operations (producers, consumers)
  - Functions, Sources, and Sinks management
- **Multiple Connection Options**:
  - Connect to StreamNative Cloud with service account authentication
  - Connect directly to external Apache Kafka clusters
  - Connect directly to external Apache Pulsar clusters

> *: The Kafka Connect operations are only tested and verfied on StreamNative Cloud. 

## Installation

```bash
# Clone the repository
git clone https://github.com/streamnative/streamnative-mcp-server.git
cd streamnative-mcp-server

go mod tidy
go mod download

# Build the binary
make
```

## Usage

### Prerequisites 

If you want to access to your StreamNative Cloud, you will need to have following resources ready:

1. Access to [StreamNative Cloud](https://console.streamnative.cloud/?defaultMethod=signup).
2. StreamNative Cloud Organization
3. StreamNative Cloud instance and cluster
4. Service Account with admin role
5. Download the Service Account Key file

### Start the MCP Server

```bash
# Start MCP server with StreamNative Cloud authentication
snmcp stdio --organization my-org --key-file /path/to/key-file.json

# Start MCP server with external Kafka
bin/snmcp stdio --use-external-kafka --kafka-bootstrap-servers localhost:9092 --kafka-auth-type SASL_SSL --kafka-auth-mechanism PLAIN --kafka-auth-user user --kafka-auth-pass pass --kafka-use-tls --kafka-schema-registry-url https://sr.local --kafka-schema-registry-auth-user user --kafka-schema-registry-auth-pass pass

# Start MCP server with external Pulsar
snmcp stdio --use-external-pulsar --pulsar-web-service-url http://pulsar.example.com:8080
bin/snmcp stdio --use-external-pulsar --pulsar-web-service-url http://pulsar.example.com:8080  --pulsar-token "xxx"
```

### Command-line Options

```
Usage:
  snmcp [command]

Available Commands:
  stdio       Start stdio server
  help        Help about any command

Flags:
      --audience string                        The audience identifier for the API server (default "https://api.streamnative.cloud")
      --client-id string                       The client ID to use for authorization grants (default "AJYEdHWi9EFekEaUXkPWA2MqQ3lq1NrI")
      --config-dir string                      If present, the config directory to use
      --enable-command-logging                 When enabled, the server will log all command requests and responses to the log file
      --features strings                       Features to enable, defaults to `all`
  -h, --help                                   help for snmcp
      --issuer string                          The OAuth 2.0 issuer endpoint (default "https://auth.streamnative.cloud/")
      --kafka-auth-mechanism string            The auth mechanism to use for Kafka
      --kafka-auth-pass string                 The auth password to use for Kafka
      --kafka-auth-type string                 The auth type to use for Kafka
      --kafka-auth-user string                 The auth user to use for Kafka
      --kafka-bootstrap-servers string         The bootstrap servers to use for Kafka
      --kafka-ca-file string                   The CA file to use for Kafka
      --kafka-client-cert-file string          The client certificate file to use for Kafka
      --kafka-client-key-file string           The client key file to use for Kafka
      --kafka-schema-registry-auth-pass string The auth password to use for the schema registry
      --kafka-schema-registry-auth-user string The auth user to use for the schema registry
      --kafka-schema-registry-bearer-token string The bearer token to use for the schema registry
      --kafka-schema-registry-url string       The schema registry URL to use for Kafka
      --key-file string                        The key file to use for authentication to StreamNative Cloud
      --log-file string                        Path to log file
      --organization string                    The organization to use for the API server
      --proxy-location string                  The proxy location to use for the API server (default "https://proxy.streamnative.cloud")
      --pulsar-auth-params string              The auth params to use for Pulsar
      --pulsar-auth-plugin string              The auth plugin to use for Pulsar
      --pulsar-token string                    The token to use for Pulsar
      --pulsar-cluster string                  The default cluster to use for the API server
      --pulsar-instance string                 The default instance to use for the API server
      --pulsar-tls-allow-insecure-connection   The TLS allow insecure connection to use for Pulsar
      --pulsar-tls-cert-file string            The TLS cert file to use for Pulsar
      --pulsar-tls-enable-hostname-verification The TLS enable hostname verification to use for Pulsar (default true)
      --pulsar-tls-key-file string             The TLS key file to use for Pulsar
      --pulsar-tls-trust-certs-file-path string The TLS trust certs file path to use for Pulsar
      --pulsar-web-service-url string          The web service URL to use for Pulsar
  -r, --read-only                              Read-only mode
      --server string                          The server to connect to (default "https://api.streamnative.cloud")
      --use-external-kafka                     Use external Kafka
      --use-external-pulsar                    Use external Pulsar
  -v, --version                                version for snmcp
```

## Tool Configuration

The StreamNative MCP Server supports enabling or disabling specific groups of functionalities via the `--features` flag. This allows you to control which MCP tools are available to your AI tools. Enabling only the toolsets that you need can help the LLM with tool choice and reduce the context size.

### Available Features

The StreamNative MCP Server allows you to enable or disable specific groups of features using the `--features` flag. This helps you control which tools are available to your AI agents and can reduce context size for LLMs.

#### Combination Feature Sets

| Feature        | Description                                                                 |
|---------------|-----------------------------------------------------------------------------|
| `all`         | Enables all features: StreamNative Cloud, Pulsar, and Kafka tools            |

---

#### Kafka Features

| Feature                   | Description                                      |
|--------------------------|--------------------------------------------------|
| `all-kafka`   | Enables all Kafka admin and client tools, without Apache Pulsar and StreamNative Cloud tools                                    |
| `kafka-admin`             | Kafka administrative operations (all admin tools) |
| `kafka-client`            | Kafka client operations (produce/consume)         |
| `kafka-admin-topics`      | Manage Kafka topics                              |
| `kafka-admin-partitions`  | Manage Kafka partitions                          |
| `kafka-admin-groups`      | Manage Kafka consumer groups                     |
| `kafka-admin-schema-registry` | Interact with Kafka Schema Registry           |
| `kafka-admin-connect`     | Manage Kafka Connect connectors                  |

---

#### Pulsar Features

| Feature                        | Description                                      |
|-------------------------------|--------------------------------------------------|
| `all-pulsar`  | Enables all Pulsar admin and client tools, without Apache Kafka and StreamNative Cloud tools                                    |
| `pulsar-admin`                 | Pulsar administrative operations (all admin tools)|
| `pulsar-client`                | Pulsar client operations (produce/consume)        |
| `pulsar-admin-brokers`         | Manage Pulsar brokers                            |
| `pulsar-admin-broker-stats`    | Access Pulsar broker statistics                   |
| `pulsar-admin-clusters`        | Manage Pulsar clusters                            |
| `pulsar-admin-functions-worker`| Manage Pulsar Function workers                    |
| `pulsar-admin-namespaces`      | Manage Pulsar namespaces                          |
| `pulsar-admin-namespace-policy`| Configure Pulsar namespace policies               |
| `pulsar-admin-isolation-policy`| Manage namespace isolation policies               |
| `pulsar-admin-packages`        | Manage Pulsar packages                            |
| `pulsar-admin-resource-quotas` | Configure resource quotas                         |
| `pulsar-admin-schemas`         | Manage Pulsar schemas                             |
| `pulsar-admin-subscriptions`   | Manage Pulsar subscriptions                       |
| `pulsar-admin-tenants`         | Manage Pulsar tenants                             |
| `pulsar-admin-topics`          | Manage Pulsar topics                              |
| `pulsar-admin-sinks`           | Manage Pulsar IO sinks                            |
| `pulsar-admin-functions`       | Manage Pulsar Functions                           |
| `pulsar-admin-sources`         | Manage Pulsar IO sources                          |
| `pulsar-admin-topic-policy`    | Configure Pulsar topic policies                   |

---

#### StreamNative Cloud Features

| Feature              | Description                                                      |
|---------------------|------------------------------------------------------------------|
| `streamnative-cloud`| Manage StreamNative Cloud context and check resource logs         |

You can combine these features as needed using the `--features` flag. For example, to enable only Pulsar client features:

```bash
snmcp stdio --organization my-org --key-file /path/to/key-file.json --features pulsar-client
```

### MCP Tools per feature

#### kafka-admin-topics

This tool provides access to various Kafka topic operations, including creation, deletion, listing, and configuration retrieval.

- **topics**
  - **list**: List all topics in the Kafka cluster
    - `include-internal` (boolean, optional): Whether to include internal Kafka topics (those starting with an underscore). Default: false

- **topic**
  - **get**: Get detailed configuration for a specific topic
    - `name` (string, required): The name of the Kafka topic
  - **create**: Create a new topic
    - `name` (string, required): The name of the Kafka topic
    - `partitions` (number, optional): Number of partitions. Default: 1
    - `replication-factor` (number, optional): Replication factor. Default: 1
    - `configs` (array of string, optional): Topic configuration overrides as key-value strings, e.g. ["cleanup.policy=compact", "retention.ms=604800000"]
  - **delete**: Delete an existing topic
    - `name` (string, required): The name of the Kafka topic
  - **metadata**: Get metadata for a specific topic
    - `name` (string, required): The name of the Kafka topic

---

#### kafka-admin-partitions

This tool provides access to Kafka partition operations, particularly adding partitions to existing topics.

- **partition**
  - **update**: Update the number of partitions for an existing Kafka topic (can only increase, not decrease)
    - `topic` (string, required): The name of the Kafka topic
    - `new-total` (number, required): The new total number of partitions (must be greater than current)

---

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

---

#### kafka-admin-schema-registry

This tool provides access to Kafka Schema Registry operations, including managing subjects, versions, and compatibility settings.

- **subjects**
  - **list**: List all schema subjects in the Schema Registry
    - _Parameters_: None

- **subject**
  - **get**: Get the latest schema for a subject
    - `subject` (string, required): The subject name
  - **create**: Register a new schema for a subject
    - `subject` (string, required): The subject name
    - `schema` (string, required): The schema definition (in AVRO/JSON/PROTOBUF, etc.)
    - `type` (string, optional): The schema type (e.g. AVRO, JSON, PROTOBUF)
  - **delete**: Delete a schema subject
    - `subject` (string, required): The subject name

- **versions**
  - **list**: List all versions for a specific subject
    - `subject` (string, required): The subject name
  - **get**: Get a specific version of a subject's schema
    - `subject` (string, required): The subject name
    - `version` (number, required): The version number
  - **delete**: Delete a specific version of a subject's schema
    - `subject` (string, required): The subject name
    - `version` (number, required): The version number

- **compatibility**
  - **get**: Get compatibility setting for a subject
    - `subject` (string, required): The subject name
  - **set**: Set compatibility level for a subject
    - `subject` (string, required): The subject name
    - `level` (string, required): The compatibility level (e.g. BACKWARD, FORWARD, FULL, NONE)

- **types**
  - **list**: List supported schema types (e.g. AVRO, JSON, PROTOBUF)
    - _Parameters_: None

---

#### kafka-admin-connect

Kafka Connect is a framework for integrating Kafka with external systems. The following resources and operations are supported:

- **kafka-connect-cluster**
  - **get**: Get information about the Kafka Connect cluster
    - _Parameters_: None

- **connectors**
  - **list**: List all connectors in the cluster
    - _Parameters_: None

- **connector**
  - **get**: Get details of a specific connector
    - `name` (string, required): The connector name
  - **create**: Create a new connector
    - `name` (string, required): The connector name
    - `config` (object, required): Connector configuration
      - Must include at least `connector.class` and other required fields for the connector type
  - **update**: Update an existing connector
    - `name` (string, required): The connector name
    - `config` (object, required): Updated configuration
  - **delete**: Delete a connector
    - `name` (string, required): The connector name
  - **restart**: Restart a connector
    - `name` (string, required): The connector name
  - **pause**: Pause a connector
    - `name` (string, required): The connector name
  - **resume**: Resume a paused connector
    - `name` (string, required): The connector name

- **connector-plugins**
  - **list**: List all available connector plugins
    - _Parameters_: None

#### kafka-client-consume

Consume messages from a Kafka topic. This tool allows you to read messages from Kafka topics with various consumption options.

- **kafka_client_consume**
  - **Description**: Read messages from a Kafka topic, with support for consumer groups, offset control, and timeouts. If schema registry integration enabled, and the topic have schema with `topicName-value`, the consume tool will try to use the schema to decode the messages.
  - **Parameters**:
    - `topic` (string, required): The name of the Kafka topic to consume messages from.
    - `group` (string, optional): The consumer group ID to use. If provided, offsets are tracked and committed; otherwise, a random group is used and offsets are not committed.
    - `offset` (string, optional): The offset position to start consuming from. One of:
      - 'atstart': Begin from the earliest available message (default)
      - 'atend': Begin from the next message after the consumer starts
      - 'atcommitted': Begin from the last committed offset (only works with specified 'group')
    - `max-messages` (number, optional): Maximum number of messages to consume in this request. Default: 10
    - `timeout` (number, optional): Maximum time in seconds to wait for messages. Default: 10

---

#### kafka-client-produce

Produce messages to a Kafka topic. This tool allows you to send single or multiple messages with various options.

- **kafka_client_produce**
  - **Description**: Send messages to a Kafka topic, supporting keys, headers, partitions, batching, and file-based payloads.
  - **Parameters**:
    - `topic` (string, required): The name of the Kafka topic to produce messages to.
    - `key` (string, optional): The key for the message. Used for partition assignment and ordering.
    - `value` (string, required if 'messages' is not provided): The value/content of the message to send.
    - `headers` (array, optional): Message headers in the format of [{"key": "header-key", "value": "header-value"}].
    - `sync` (boolean, optional): Whether to wait for server acknowledgment before returning. Default: true.
    

## Integration with MCP Clients

This server can be used with any MCP-compatible client, such as:

- Claude Desktop
- Other AI assistants supporting the MCP protocol
- Custom applications built with MCP client libraries

### Usage with Claude Desktop

```json
{
  "mcpServers": {
    "github": {
      "command": "snmcp",
      "args": [
        "stdio",
        "--organization",
        "my-org",
        "--key-file",
        "/path/to/key-file.json"
      ],
    }
  }
}
```

## About Model Context Protocol (MCP)

The Model Context Protocol (MCP) is an open protocol that standardizes how applications provide context to LLMs. MCP helps build agents and complex workflows on top of LLMs by providing:

- A growing list of pre-built integrations that your LLM can directly plug into
- The flexibility to switch between LLM providers and vendors
- Best practices for securing your data within your infrastructure

For more information, visit [modelcontextprotocol.io](https://modelcontextprotocol.io/introduction).

## Release

_This section describes how to release a new version of `snmcp`._

1. Generate a tag for the new version (see Versioning, below):

```
git tag -a v0.0.1 -m "v0.0.1"
```

5. Push the tag to the git repository:

```
git push origin refs/tags/v0.0.1
```

The release workflow will:

- build Go binaries for supported platforms
- archive the binaries
- publish a release to the github repository ([ref](https://github.com/streamnative/streamnative-mcp-server/releases))

### Versioning

This project uses [semver](https://semver.org/) semantics.

- Stable: `vX.Y.Z`
- Pre-release: `vX.Y.Z-rc.W`
- Snapshot: `vX.Y.Z-SNAPSHOT-commit`

## License

Licensed under the Apache License Version 2.0: http://www.apache.org/licenses/LICENSE-2.0
