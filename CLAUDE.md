# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
make build                  # Build server binary to bin/snmcp
make docker-build           # Build local Docker image (both streamnative/mcp-server and streamnative/snmcp tags)
make docker-build-push      # Build and push multi-platform image (linux/amd64,linux/arm64)
make docker-build-multiplatform  # Build multi-platform image locally
make docker-buildx-setup    # Setup Docker buildx for multi-platform builds
make license-check          # Check license headers
make license-fix            # Fix license headers
go test -race ./...         # Run all tests with race detection
go test -race ./pkg/mcp/builders/...  # Run specific package tests
go test -v -run TestName ./pkg/...    # Run a single test
```

## Architecture Overview

The StreamNative MCP Server implements the Model Context Protocol using the `mark3labs/mcp-go` library to enable AI agents to interact with Apache Kafka, Apache Pulsar, and StreamNative Cloud resources.

### Request Flow

```
Client Request → MCP Server (pkg/mcp/server.go)
                    ↓
              SSE/stdio transport layer (pkg/cmd/mcp/)
                    ↓
              Tool Handler (from builders)
                    ↓
              Context Functions (pkg/mcp/ctx.go)
                    ↓
              Service Client (Kafka/Pulsar/SNCloud)
```

### Core Components

1. **Server & Sessions** (`pkg/mcp/server.go`)
   - `Server` struct holds `MCPServer`, `KafkaSession`, `PulsarSession`, and `SNCloudSession`
   - Sessions provide lazy-initialized clients for each service
   - Context functions (`pkg/mcp/ctx.go`) inject/retrieve sessions from request context

2. **Tool Builders Framework** (`pkg/mcp/builders/`)
   - `ToolBuilder` interface: `GetName()`, `GetRequiredFeatures()`, `BuildTools()`, `Validate()`
   - `BaseToolBuilder` provides common feature validation logic
   - `ToolRegistry` manages all tool builders with concurrent-safe registration
   - `ToolBuildConfig` specifies build parameters (ReadOnly, Features, Options)
   - `ToolMetadata` describes tool information (Name, Version, Description, Category, Tags)

3. **Tool Builders Organization**
   - `builders/kafka/` - Kafka-specific tool builders (connect, consume, groups, partitions, produce, schema_registry, topics)
   - `builders/pulsar/` - Pulsar-specific tool builders (brokers, brokers_stats, cluster, functions, functions_worker, namespace, namespace_policy, nsisolationpolicy, packages, resourcequotas, schema, sinks, sources, subscription, tenant, topic, topic_policy)
   - `builders/streamnative/` - StreamNative Cloud tool builders

4. **Tool Registration** (`pkg/mcp/*_tools.go`)
   - Each `*_tools.go` file creates a builder, builds tools, and adds them to the server
   - Tools are conditionally registered based on `--features` flag
   - Feature constants defined in `pkg/mcp/features.go`

5. **PFTools - Functions as Tools** (`pkg/mcp/pftools/`)
   - `PulsarFunctionManager` dynamically converts Pulsar Functions to MCP tools
   - Polls for function changes and auto-registers/unregisters tools
   - Circuit breaker pattern (`circuit_breaker.go`) for fault tolerance
   - Schema conversion (`schema.go`) for input/output handling

6. **Session Management** (`pkg/mcp/session/`)
   - `pulsar_session_manager.go` - LRU session cache with TTL cleanup for multi-session mode

7. **Transport Layer** (`pkg/cmd/mcp/`)
   - `sse.go` - SSE transport with health endpoints (`/healthz`, `/readyz`) and auth middleware
   - `server.go` - Stdio transport and common server initialization

### Key Design Patterns

- **Builder Pattern**: Tool builders create tools based on features and read-only mode
- **Registry Pattern**: ToolRegistry provides centralized management of all builders
- **Context Injection**: Sessions passed via `context.Context` using typed keys
- **Feature Flags**: Tools enabled/disabled via string feature identifiers
- **Circuit Breaker**: PFTools uses failure thresholds to prevent cascading failures
- **Multi-Session Pattern**: Per-user Pulsar sessions with LRU caching for SSE mode

## Adding New Tools

1. **Create Builder** in `pkg/mcp/builders/kafka/` or `pkg/mcp/builders/pulsar/`:
   ```go
   type MyToolBuilder struct {
       *builders.BaseToolBuilder
   }

   func NewMyToolBuilder() *MyToolBuilder {
       metadata := builders.ToolMetadata{
           Name:        "my_tool",
           Description: "Tool description",
           Category:    "kafka_admin",
       }
       features := []string{"kafka-admin", "all", "all-kafka"}
       return &MyToolBuilder{
           BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
       }
   }

   func (b *MyToolBuilder) BuildTools(ctx context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
       if !b.HasAnyRequiredFeature(config.Features) {
           return nil, nil
       }
       // Build and return tools
   }
   ```

2. **Add Feature Constant** in `pkg/mcp/features.go` if needed

3. **Create Registration File** `pkg/mcp/my_tools.go`:
   ```go
   func AddMyTools(s *server.MCPServer, readOnly bool, features []string) {
       builder := kafkabuilders.NewMyToolBuilder()
       config := builders.ToolBuildConfig{ReadOnly: readOnly, Features: features}
       tools, _ := builder.BuildTools(context.Background(), config)
       for _, tool := range tools {
           s.AddTool(tool.Tool, tool.Handler)
       }
   }
   ```

4. **Get Session in Handler**:
   ```go
   session := mcp.GetKafkaSession(ctx)  // or GetPulsarSession
   if session == nil {
       return mcp.NewToolResultError("session not found"), nil
   }
   admin, err := session.GetAdminClient()
   ```

## Session Context Access

Handlers receive sessions via context (see `pkg/mcp/ctx.go`):
- `mcp.GetKafkaSession(ctx)` → `*kafka.Session`
- `mcp.GetPulsarSession(ctx)` → `*pulsar.Session`
- `mcp.GetSNCloudSession(ctx)` → `*config.Session`
- `mcp.GetSNCloudOrganization(ctx)` → organization string
- `mcp.GetSNCloudInstance(ctx)` → instance string
- `mcp.GetSNCloudCluster(ctx)` → cluster string

From sessions:
- `session.GetAdminClient()` / `session.GetAdminV3Client()` for Pulsar admin
- `session.GetPulsarClient()` for Pulsar messaging
- `session.GetAdminClient()` for Kafka admin (via franz-go/kadm)

## Configuration Modes

1. **StreamNative Cloud**: `--organization` + `--key-file`
2. **External Kafka**: `--use-external-kafka` + Kafka params
3. **External Pulsar**: `--use-external-pulsar` + Pulsar params
4. **Multi-Session Pulsar** (SSE only): `--use-external-pulsar` + `--multi-session-pulsar`

Pre-configured context: `--pulsar-instance` + `--pulsar-cluster` disables context management tools.

### Multi-Session Pulsar Mode

When `--multi-session-pulsar` is enabled (SSE server with external Pulsar only):

- **No global PulsarSession**: Each request must provide its own token via `Authorization: Bearer <token>` header
- **HTTP 401 on auth failure**: Requests without valid tokens are rejected with HTTP 401 Unauthorized
- **Per-user session caching**: Sessions are cached using LRU with configurable size and TTL
- **Session management**: See `pkg/mcp/session/pulsar_session_manager.go`

Key files:
- `pkg/cmd/mcp/sse.go` - Auth middleware wraps SSEHandler()/MessageHandler(), health endpoints
- `pkg/mcp/session/pulsar_session_manager.go` - LRU session cache with TTL cleanup
- `pkg/cmd/mcp/server.go` - Skips global PulsarSession when multi-session enabled

### Health Endpoints

SSE server exposes health check endpoints:
- `GET /mcp/healthz` - Liveness probe (always returns "ok")
- `GET /mcp/readyz` - Readiness probe (always returns "ready")

## Feature Flags

Available feature flags (defined in `pkg/mcp/features.go`):

| Feature | Description |
|---------|-------------|
| `all` | Enable all features |
| `all-kafka` | All Kafka features |
| `all-pulsar` | All Pulsar features |
| `kafka-client` | Kafka produce/consume |
| `kafka-admin` | Kafka admin operations (all admin tools) |
| `kafka-admin-schema-registry` | Schema Registry |
| `kafka-admin-kafka-connect` | Kafka Connect |
| `kafka-admin-topics` | Manage Kafka topics |
| `kafka-admin-partitions` | Manage Kafka partitions |
| `kafka-admin-groups` | Manage Kafka consumer groups |
| `pulsar-admin` | Pulsar admin operations (all admin tools) |
| `pulsar-client` | Pulsar produce/consume |
| `pulsar-admin-brokers` | Manage Pulsar brokers |
| `pulsar-admin-brokers-status` | Pulsar broker status |
| `pulsar-admin-broker-stats` | Access Pulsar broker statistics |
| `pulsar-admin-clusters` | Manage Pulsar clusters |
| `pulsar-admin-functions` | Manage Pulsar Functions |
| `pulsar-admin-functions-worker` | Manage Pulsar Function workers |
| `pulsar-admin-namespaces` | Manage Pulsar namespaces |
| `pulsar-admin-namespace-policy` | Configure namespace policies |
| `pulsar-admin-ns-isolation-policy` | Manage namespace isolation policies |
| `pulsar-admin-packages` | Manage Pulsar packages |
| `pulsar-admin-resource-quotas` | Configure resource quotas |
| `pulsar-admin-schemas` | Manage Pulsar schemas |
| `pulsar-admin-subscriptions` | Manage Pulsar subscriptions |
| `pulsar-admin-tenants` | Manage Pulsar tenants |
| `pulsar-admin-topics` | Manage Pulsar topics |
| `pulsar-admin-sinks` | Manage Pulsar IO sinks |
| `pulsar-admin-sources` | Manage Pulsar Sources |
| `pulsar-admin-topic-policy` | Configure topic policies |
| `streamnative-cloud` | StreamNative Cloud context management |
| `functions-as-tools` | Dynamic Pulsar Functions as MCP tools |

## Helm Chart

The project includes a Helm chart for Kubernetes deployment at `charts/snmcp/`:

```bash
# Basic installation
helm install snmcp ./charts/snmcp \
  --set pulsar.webServiceURL=http://pulsar.example.com:8080

# With TLS
helm install snmcp ./charts/snmcp \
  --set pulsar.webServiceURL=https://pulsar:8443 \
  --set pulsar.tls.enabled=true \
  --set pulsar.tls.secretName=pulsar-tls
```

The chart runs MCP Server in Multi-Session Pulsar mode with authentication via `Authorization: Bearer <token>` header.

## SDK Packages

The project includes generated SDK packages:
- `sdk/sdk-apiserver/` - StreamNative Cloud API server client
- `sdk/sdk-kafkaconnect/` - Kafka Connect client

## Error Handling

- Wrap errors: `fmt.Errorf("failed to X: %w", err)`
- Return tool errors: `mcp.NewToolResultError("message")`
- Check session nil before operations
- For PFTools, use circuit breaker to handle repeated failures
