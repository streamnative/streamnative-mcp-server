# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
make build                  # Build server binary to bin/snmcp
go test -race ./...         # Run all tests with race detection
go test -race ./pkg/mcp/builders/...  # Run specific package tests
go test -v -run TestName ./pkg/...    # Run a single test
make license-check          # Check license headers
make license-fix            # Fix license headers
make docker-build           # Build local Docker image
make docker-build-push      # Build and push multi-platform image
```

## Architecture Overview

The StreamNative MCP Server implements the Model Context Protocol using the `mark3labs/mcp-go` library to enable AI agents to interact with Apache Kafka, Apache Pulsar, and StreamNative Cloud resources.

### Request Flow

```
Client Request → MCP Server (pkg/mcp/server.go)
                    ↓
              Tool Handler (from builders)
                    ↓
              Session Context (pkg/mcp/ctx.go)
                    ↓
              Service Client (Kafka/Pulsar/SNCloud)
```

### Core Components

1. **Server & Sessions** (`pkg/mcp/server.go`)
   - `Server` struct holds `MCPServer`, `KafkaSession`, `PulsarSession`, and `SNCloudSession`
   - Sessions provide lazy-initialized clients for each service
   - Context functions (`pkg/mcp/ctx.go`) inject/retrieve sessions from request context

2. **Tool Builders** (`pkg/mcp/builders/`)
   - `ToolBuilder` interface: `GetName()`, `GetRequiredFeatures()`, `BuildTools()`, `Validate()`
   - `BaseToolBuilder` provides common feature validation logic
   - Each builder creates `[]server.ServerTool` with tool definitions and handlers
   - Builders in `builders/kafka/` and `builders/pulsar/` implement service-specific tools

3. **Tool Registration** (`pkg/mcp/*_tools.go`)
   - Each `*_tools.go` file creates a builder, builds tools, and adds them to the server
   - Tools are conditionally registered based on `--features` flag
   - Feature constants defined in `pkg/mcp/features.go`

4. **PFTools - Functions as Tools** (`pkg/mcp/pftools/`)
   - `PulsarFunctionManager` dynamically converts Pulsar Functions to MCP tools
   - Polls for function changes and auto-registers/unregisters tools
   - Circuit breaker pattern (`circuit_breaker.go`) for fault tolerance
   - Schema conversion (`schema.go`) for input/output handling

### Key Design Patterns

- **Builder Pattern**: Tool builders create tools based on features and read-only mode
- **Context Injection**: Sessions passed via `context.Context` using typed keys
- **Feature Flags**: Tools enabled/disabled via string feature identifiers
- **Circuit Breaker**: PFTools uses failure thresholds to prevent cascading failures

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
   session := mcpCtx.GetKafkaSession(ctx)  // or GetPulsarSession
   if session == nil {
       return mcp.NewToolResultError("session not found"), nil
   }
   admin, err := session.GetAdminClient()
   ```

## Session Context Access

Handlers receive sessions via context (see `pkg/mcp/internal/context/ctx.go`):
- `mcpCtx.GetKafkaSession(ctx)` → `*kafka.Session`
- `mcpCtx.GetPulsarSession(ctx)` → `*pulsar.Session`
- `mcpCtx.GetSNCloudSession(ctx)` → `*config.Session`

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
- `pkg/cmd/mcp/sse.go` - Auth middleware wraps SSEHandler()/MessageHandler()
- `pkg/mcp/session/pulsar_session_manager.go` - LRU session cache with TTL cleanup
- `pkg/cmd/mcp/server.go` - Skips global PulsarSession when multi-session enabled

## Error Handling

- Wrap errors: `fmt.Errorf("failed to X: %w", err)`
- Return tool errors: `mcp.NewToolResultError("message")`
- Check session nil before operations
- For PFTools, use circuit breaker to handle repeated failures