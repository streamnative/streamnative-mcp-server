# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

### Build the server binary
```bash
make build
# Or directly:
go build -o bin/snmcp cmd/streamnative-mcp-server/main.go
```

### Run tests
```bash
go test -race ./...
# Run a specific test:
go test -race ./pkg/mcp/builders/...
```

### Docker operations
```bash
make docker-build          # Build local Docker image
make docker-build-push      # Build and push multi-platform image
```

### License checking
```bash
make license-check         # Check license headers
make license-fix          # Fix license headers
```

## Architecture Overview

The StreamNative MCP Server implements the Model Context Protocol to enable AI agents to interact with Apache Kafka, Apache Pulsar, and StreamNative Cloud resources.

### Core Components

1. **Session Management** (`pkg/config/`, `pkg/kafka/`, `pkg/pulsar/`)
   - Three types of sessions: SNCloudSession, KafkaSession, PulsarSession
   - Sessions manage client connections and authentication
   - Sessions are created and configured based on command-line flags

2. **MCP Server** (`pkg/mcp/`)
   - Central server implementation using `mark3labs/mcp-go` library
   - Handles tool registration and request routing
   - Features can be enabled/disabled via `--features` flag

3. **Tool Builders** (`pkg/mcp/builders/`)
   - Registry pattern for registering MCP tools
   - Separate builders for Kafka and Pulsar operations
   - Each builder creates tools with specific operations (admin, client, etc.)

4. **PFTools** (`pkg/mcp/pftools/`)
   - Abstraction layer for Pulsar Functions as MCP tools
   - Dynamic tool generation from deployed Pulsar Functions
   - Circuit breaker pattern for resilience
   - Schema handling for input/output validation

### Key Design Patterns

1. **Builder Pattern**: Tool builders (`pkg/mcp/builders/`) register tools dynamically based on enabled features
2. **Session Pattern**: Separate sessions for different services (Kafka, Pulsar, SNCloud) with lazy initialization
3. **Registry Pattern**: Central registry (`pkg/mcp/builders/registry.go`) manages all tool builders
4. **Circuit Breaker**: Used in PFTools for handling function invocation failures

## Development Guidelines

### Adding New Tools

1. Create a new builder in `pkg/mcp/builders/kafka/` or `pkg/mcp/builders/pulsar/`
2. Implement the `Builder` interface with `Build()` method
3. Register the builder in the appropriate tools file (e.g., `kafka_admin_*_tools.go`)
4. Add feature flag support in `pkg/mcp/features.go`

### Session Context

The server maintains session context that gets passed to tools via the context:
- Pulsar admin client retrieval: Use `session.GetAdminClient()` or `session.GetAdminV3Client()`
- Kafka client retrieval: Use `session.GetKafkaSession()` methods
- Always check for nil sessions before use

### Error Handling

- Use wrapped errors with context: `fmt.Errorf("failed to X: %w", err)`
- Check session availability before operations
- Return meaningful error messages for AI agent consumption

## Testing

Tests follow standard Go testing patterns:
- Unit tests alongside source files (`*_test.go`)
- Use `testify` for assertions
- Mock external dependencies where appropriate

## Important Files

- `cmd/streamnative-mcp-server/main.go` - Entry point
- `pkg/cmd/mcp/server.go` - Server setup
- `pkg/mcp/server.go` - MCP server implementation
- `pkg/mcp/builders/registry.go` - Tool registration
- `pkg/mcp/pftools/manager.go` - Pulsar Functions management
- `pkg/config/config.go` - Configuration structures

## Configuration

The server supports three modes:
1. **StreamNative Cloud**: Requires `--organization` and `--key-file`
2. **External Kafka**: Use `--use-external-kafka` with Kafka connection parameters
3. **External Pulsar**: Use `--use-external-pulsar` with Pulsar connection parameters

When `--pulsar-instance` and `--pulsar-cluster` are provided together, context management tools are disabled as the context is pre-configured.