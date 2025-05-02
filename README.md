# StreamNative MCP Server

A Model Context Protocol (MCP) server for integrating AI agents with StreamNative Cloud resources and Apache Pulsar/Kafka messaging systems.

## Overview

StreamNative MCP Server provides a standard interface for LLMs (Large Language Models) and AI agents to interact with StreamNative Cloud services, Apache Pulsar, and Apache Kafka. This implementation follows the [Model Context Protocol](https://modelcontextprotocol.io/introduction) specification, enabling AI applications to access messaging services through a standardized interface.

## Features

- **StreamNative Cloud Integration**: 
  - Connect to StreamNative Cloud resources with authentication
  - Switch to clusters available in your organization
  - Describe the status of clusters resources
- **Apache Pulsar Support**: Interact with Pulsar resources including:
  - Pulsar Admin operations (topics, namespaces, tenants, schemas, etc.)
  - Pulsar Client operations (producers, consumers)
  - Functions, Sources, and Sinks management
- **Apache Kafka Support**: Interact with Kafka resources including:
  - Kafka Admin operations (topics, partitions, consumer groups)
  - Schema Registry operations
  - Kafka Connect operations
  - Kafka Client operations (producers, consumers)
- **Multiple Connection Options**:
  - Connect to StreamNative Cloud with service account authentication
  - Connect directly to external Pulsar clusters
  - Connect directly to external Kafka clusters

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

The following sets of tools are available (all available by default on StreamNative Cloud)

| Features | Description |
| ------|-------|
| `all` | All tools, including StreamNative Cloud tools, Pulsar tools and Kafka tools |
| `all-pulsar` | All Pulsar admin and Pulsar client tools |
| `all-kafka` | All Kafka admin and Kafka client tools |
| `pulsar-admin` | Pulsar administrative operations (including all `pulsar-admin-*`) |
| `pulsar-client` | Pulsar client operations (produce and consume messages) |
| `pulsar-admin-brokers` | Manage Pulsar brokers |
| `pulsar-admin-broker-stats` | Access Pulsar broker statistics |
| `pulsar-admin-clusters` | Manage Pulsar clusters |
| `pulsar-admin-functions-worker` | Manage Pulsar Function workers |
| `pulsar-admin-namespaces` | Manage Pulsar namespaces |
| `pulsar-admin-namespace-policy` | Configure Pulsar namespace policies |
| `pulsar-admin-isolation-policy` | Manage namespace isolation policies |
| `pulsar-admin-packages` | Manage Pulsar packages |
| `pulsar-admin-resource-quotas` | Configure resource quotas |
| `pulsar-admin-schemas` | Manage Pulsar schemas |
| `pulsar-admin-subscriptions` | Manage Pulsar subscriptions |
| `pulsar-admin-tenants` | Manage Pulsar tenants |
| `pulsar-admin-topics` | Manage Pulsar topics |
| `pulsar-admin-sinks` | Manage Pulsar IO sinks |
| `pulsar-admin-functions` | Manage Pulsar Functions |
| `pulsar-admin-sources` | Manage Pulsar IO sources |
| `pulsar-admin-topic-policy` | Configure Pulsar topic policies |
| `kafka-admin` | Kafka administrative operations (including all `kafka-admin-*`) |
| `kafka-client` | Kafka client operations (produce and consume messages) |
| `kafka-admin-topics` | Manage Kafka partitions |
| `kafka-admin-partitions` | Manage Kafka partitions |
| `kafka-admin-groups` | Manage Kafka consumer groups |
| `kafka-admin-schema-registry` | Interact with Kafka Schema Registry |
| `kafka-admin-connect` | Manage Kafka Connect connectors |
| `streamnative-cloud` | Manage the context, check resources logs of StreamNative Cloud |

### Usage Examples

To enable only specific feature sets:

```bash
# Enable only Pulsar client features
snmcp stdio --organization my-org --key-file /path/to/key-file.json --features pulsar-client
```

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

## License

Licensed under the Apache License Version 2.0: http://www.apache.org/licenses/LICENSE-2.0
