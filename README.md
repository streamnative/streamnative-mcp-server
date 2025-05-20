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

### Homebrew (macOS and Linux)

The easiest way to install streamnative-mcp-server is using Homebrew:

```bash
# Add the tap repository
brew tap streamnative/streamnative

# Install streamnative-mcp-server
brew install streamnative/streamnative/snmcp
```

### Docker Image

StreamNative MCP Server releases the Docker Image to [streamnative/snmcp](https://hub.docker.com/r/streamnative/snmcp), and it can be used to run both stdio server and sse server via docker command.

```bash
# Pull image from Docker Hub
docker pull streamnative/snmcp 
```

### From Github Release

Visit https://github.com/streamnative/streamnative-mcp-server/releases to get the latest binary of StreamNative MCP Server.

### From Source

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

#### Using stdio Server

```bash
# Start MCP server with StreamNative Cloud authentication
bin/snmcp stdio --organization my-org --key-file /path/to/key-file.json

# Start MCP server with external Kafka
bin/snmcp stdio --use-external-kafka --kafka-bootstrap-servers localhost:9092 --kafka-auth-type SASL_SSL --kafka-auth-mechanism PLAIN --kafka-auth-user user --kafka-auth-pass pass --kafka-use-tls --kafka-schema-registry-url https://sr.local --kafka-schema-registry-auth-user user --kafka-schema-registry-auth-pass pass

# Start MCP server with external Pulsar
snmcp stdio --use-external-pulsar --pulsar-web-service-url http://pulsar.example.com:8080
bin/snmcp stdio --use-external-pulsar --pulsar-web-service-url http://pulsar.example.com:8080  --pulsar-token "xxx"

# Start MCP server with SSE by docker with StreamNative Cloud authentication
docker run -i --rm -e SNMCP_ORGANIZATION=my-org -e SNMCP_KEY_FILE=/key.json -v /path/to/key-file.json:/key.json -p 9090:9090 streamnative/snmcp stdio
```

#### Using SSE (Server-Sent Events) Server

```bash
# Start MCP server with SSE and StreamNative Cloud authentication
snmcp sse --http-addr :9090 --http-path /mcp --organization my-org --key-file /path/to/key-file.json

# Start MCP server with SSE and external Kafka
snmcp sse --http-addr :9090 --http-path /mcp --use-external-kafka --kafka-bootstrap-servers localhost:9092

# Start MCP server with SSE and external Pulsar
snmcp sse --http-addr :9090 --http-path /mcp --use-external-pulsar --pulsar-web-service-url http://pulsar.example.com:8080

# Start MCP server with SSE by docker with StreamNative Cloud authentication
docker run -i --rm -e SNMCP_ORGANIZATION=my-org -e SNMCP_KEY_FILE=/key.json -v /path/to/key-file.json:/key.json -p 9090:9090 streamnative/snmcp sse
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
      --http-addr string                       HTTP server address (default ":9090")
      --http-path string                       HTTP server path for SSE endpoint (default "/mcp")
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

| Feature                   | Description                                      | Docs |
|--------------------------|--------------------------------------------------|------|
| `all-kafka`   | Enables all Kafka admin and client tools, without Apache Pulsar and StreamNative Cloud tools                                    |
| `kafka-admin`             | Kafka administrative operations (all admin tools) | |
| `kafka-client`            | Kafka client operations (produce/consume)         |[kafka_client_consume.md](docs/tools/kafka_client_consume.md), [kafka_client_produce.md](docs/tools/kafka_client_produce.md)  |
| `kafka-admin-topics`      | Manage Kafka topics                              | [kafka_admin_topics.md](docs/tools/kafka_admin_topics.md) |
| `kafka-admin-partitions`  | Manage Kafka partitions                          | [kafka_admin_partitions.md](docs/tools/kafka_admin_partitions.md) |
| `kafka-admin-groups`      | Manage Kafka consumer groups                     | [kafka_admin_groups.md](docs/tools/kafka_admin_groups.md) |
| `kafka-admin-schema-registry` | Interact with Kafka Schema Registry          | [kafka_admin_schema_registry.md](docs/tools/kafka_admin_schema_registry.md) |
| `kafka-admin-connect`     | Manage Kafka Connect connectors                  | [kafka_admin_connect.md](docs/tools/kafka_admin_connect.md) |

---

#### Pulsar Features

| Feature                   | Description                                      | Docs |
|--------------------------|--------------------------------------------------|------|
| `all-pulsar`  | Enables all Pulsar admin and client tools, without Apache Kafka and StreamNative Cloud tools                                    | |
| `pulsar-admin`                 | Pulsar administrative operations (all admin tools)|  |
| `pulsar-client`                | Pulsar client operations (produce/consume)        | [pulsar_client_consume.md](docs/tools/pulsar_client_consume.md), [pulsar_client_produce.md](docs/tools/pulsar_client_produce.md) |
| `pulsar-admin-brokers`         | Manage Pulsar brokers                            | [pulsar_admin_brokers.md](docs/tools/pulsar_admin_brokers.md) |
| `pulsar-admin-broker-stats`    | Access Pulsar broker statistics                   | [pulsar_admin_broker_stats.md](docs/tools/pulsar_admin_broker_stats.md) |
| `pulsar-admin-clusters`        | Manage Pulsar clusters                            | [pulsar_admin_clusters.md](docs/tools/pulsar_admin_clusters.md) |
| `pulsar-admin-functions-worker`| Manage Pulsar Function workers                    | [pulsar_admin_functions_worker.md](docs/tools/pulsar_admin_functions_worker.md) |
| `pulsar-admin-namespaces`      | Manage Pulsar namespaces                          | [pulsar_admin_namespaces.md](docs/tools/pulsar_admin_namespaces.md) |
| `pulsar-admin-namespace-policy`| Configure Pulsar namespace policies               | [pulsar_admin_namespace_policy.md](docs/tools/pulsar_admin_namespace_policy.md) |
| `pulsar-admin-isolation-policy`| Manage namespace isolation policies               | [pulsar_admin_nsisolationpolicy.md](docs/tools/pulsar_admin_nsisolationpolicy.md) |
| `pulsar-admin-packages`        | Manage Pulsar packages                            | [pulsar_admin_packages.md](docs/tools/pulsar_admin_packages.md) |
| `pulsar-admin-resource-quotas` | Configure resource quotas                         | [pulsar_admin_resource_quotas.md](docs/tools/pulsar_admin_resource_quotas.md) |
| `pulsar-admin-schemas`         | Manage Pulsar schemas                             | [pulsar_admin_schemas.md](docs/tools/pulsar_admin_schemas.md) |
| `pulsar-admin-subscriptions`   | Manage Pulsar subscriptions                       | [pulsar_admin_subscriptions.md](docs/tools/pulsar_admin_subscriptions.md) |
| `pulsar-admin-tenants`         | Manage Pulsar tenants                             | [pulsar_admin_tenants.md](docs/tools/pulsar_admin_tenants.md) |
| `pulsar-admin-topics`          | Manage Pulsar topics                              | [pulsar_admin_topics.md](docs/tools/pulsar_admin_topics.md) |
| `pulsar-admin-sinks`           | Manage Pulsar IO sinks                            | [pulsar_admin_sinks.md](docs/tools/pulsar_admin_sinks.md) |
| `pulsar-admin-functions`       | Manage Pulsar Functions                           | [pulsar_admin_functions.md](docs/tools/pulsar_admin_functions.md) |
| `pulsar-admin-sources`         | Manage Pulsar Sources                           | [pulsar_admin_sources.md](docs/tools/pulsar_admin_sources.md) |
| `pulsar-admin-topic-policy`    | Configure Pulsar topic policies                   | [pulsar_admin_topic_policy.md](docs/tools/pulsar_admin_topic_policy.md) |

---

#### StreamNative Cloud Features

| Feature              | Description                                                      | Docs |
|---------------------|------------------------------------------------------------------|------|
| `streamnative-cloud`| Manage StreamNative Cloud context and check resource logs         | [streamnative_cloud.md](docs/tools/streamnative_cloud.md) |
| `functions-as-tools`     | Dynamically exposes deployed Pulsar Functions as invokable MCP tools, with automatic input/output schema handling. | [functions_as_tools.md](docs/tools/functions_as_tools.md)  |

You can combine these features as needed using the `--features` flag. For example, to enable only Pulsar client features:
```bash
# Enable only Pulsar client features
bin/snmcp stdio --organization my-org --key-file /path/to/key-file.json --features pulsar-client
```

## Inspecting the MCP Server

You can use the [@modelcontextprotocol/inspector](https://www.npmjs.com/package/@modelcontextprotocol/inspector) tool to inspect and test your MCP server. This is particularly useful for debugging and verifying your server's configuration.

### Installation

```bash
npm install -g @modelcontextprotocol/inspector
```

### Usage

```bash
# Inspect a stdio server
mcp-inspector stdio --command "snmcp stdio --organization my-org --key-file /path/to/key-file.json"

# Inspect an SSE server
mcp-inspector sse --url "http://localhost:9090/mcp"
```

The inspector provides a web interface where you can:
- View available tools and their schemas
- Test tool invocations
- Monitor server responses
- Debug connection issues

## Integration with MCP Clients

This server can be used with any MCP-compatible client, such as:

- Claude Desktop
- Other AI assistants supporting the MCP protocol
- Custom applications built with MCP client libraries

> ⚠️ Reminder: Please ensure you have an active paid plan with your LLM provider to fully utilize the MCP server.
Without it, you may encounter the error: `message will exceed the length limit for this chat`.


### Usage with Claude Desktop

#### Using stdio Server

```json
{
  "mcpServers": {
    "mcp-streamnative": {
      "command": "${PATH_TO_SNMCP}/bin/snmcp",
      "args": [
        "stdio",
        "--organization",
        "${STREAMNATIVE_CLOUD_ORGANIZATION_ID}",
        "--key-file",
        "${STREAMNATIVE_CLOUD_KEY_FILE}"
      ]
    }
  }
}
```

Please remember to replace `${PATH_TO_SNMCP}` with the actual path to the `snmcp` binary and `${STREAMNATIVE_CLOUD_ORGANIZATION_ID}` and `${STREAMNATIVE_CLOUD_KEY_FILE}` with your StreamNative Cloud organization ID and key file path, respectively.

Optionally, you can use docker image to start the stdio server if you have [Docker](https://www.docker.com/) installed.

```json
{
  "mcpServers": {
    "mcp-streamnative": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "-e",
        "SNMCP_ORGANIZATION",
        "-e",
        "SNMCP_KEY_FILE",
        "-v",
        "${STREAMNATIVE_CLOUD_KEY_FILE}:/key.json",
        "streamnative/snmcp",
        "stdio"
      ],
      "env": {
        "SNMCP_ORGANIZATION": "${STREAMNATIVE_CLOUD_ORGANIZATION_ID}",
        "SNMCP_KEY_FILE": "/key.json"
      }
    }
  }
}
```

#### Using SSE Server

First, install the mcp-proxy tool:

```bash
pip install mcp-proxy
```

Then configure Claude Desktop to use the SSE server:

```json
{
  "mcpServers": {
    "mcp-streamnative-proxy": {
      "command": "mcp-proxy",
      "args": [
        "http://localhost:9090/mcp/sse"
      ]
    }
  }
}
```

Note: If mcp-proxy is not in your system PATH, you'll need to provide the full path to the executable. For example:
- On macOS: `/Library/Frameworks/Python.framework/Versions/3.11/bin/mcp-proxy`
- On Linux: `/usr/local/bin/mcp-proxy`
- On Windows: `C:\Python311\Scripts\mcp-proxy.exe`

Please remember to replace `http://localhost:9090/mcp/sse` with the right URL.

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


