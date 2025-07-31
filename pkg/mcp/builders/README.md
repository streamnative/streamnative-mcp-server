# MCP Tools Builder Framework

This is the tool builder framework for StreamNative MCP Server, designed to provide unified management and construction of MCP tools.

## Overview

The tool builder framework separates tool creation logic from server registration logic, providing an extensible and testable architecture for managing MCP tools.

## Core Components

### 1. ToolBuilder Interface

The core interface that all tool builders must implement:

```go
type ToolBuilder interface {
    GetName() string
    GetRequiredFeatures() []string
    BuildTools(ctx context.Context, config ToolBuildConfig) ([]server.ServerTool, error)
    Validate(config ToolBuildConfig) error
}
```

### 2. BaseToolBuilder

An abstract builder that provides basic implementation:

```go
type BaseToolBuilder struct {
    metadata         ToolMetadata
    requiredFeatures []string
}
```

### 3. ToolRegistry

A registry that manages all tool builders:

```go
type ToolRegistry struct {
    mu       sync.RWMutex
    builders map[string]ToolBuilder
    metadata map[string]ToolMetadata
}
```

### 4. Configuration Management

Supports flexible configuration management:

```go
type ToolBuildConfig struct {
    ReadOnly bool
    Features []string
    Options  map[string]interface{}
}
```

## Usage Examples

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/mark3labs/mcp-go/server"
    "github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
)

func main() {
    // Create registry
    registry := builders.NewToolRegistry()
    
    // Register builders (actual builder implementation needed here)
    // registry.Register(kafkaBuilders.NewKafkaConnectToolBuilder())
    
    // Configure tools
    configs := map[string]builders.ToolBuildConfig{
        "kafka_connect": {
            ReadOnly: false,
            Features: []string{"kafka_admin_kafka_connect"},
        },
    }
    
    // Build all tools
    tools, err := registry.BuildAll(configs)
    if err != nil {
        panic(err)
    }
    
    // Create MCP server and add tools
    mcpServer := server.NewMCPServer("example", "1.0.0")
    for _, tool := range tools {
        mcpServer.AddTool(tool.Tool, tool.Handler)
    }
    
    fmt.Printf("Successfully added %d tools\n", len(tools))
}
```

### Configuration-Driven Usage

```go
// Use configuration file
config, err := builders.LoadToolsConfig("tools.yaml")
if err != nil {
    panic(err)
}

// Convert to build configurations
buildConfigs := config.ToToolBuildConfigs()

// Build tools
tools, err := registry.BuildAll(buildConfigs)
if err != nil {
    panic(err)
}
```

### Configuration File Example

```yaml
# tools.yaml
tools:
  kafka_connect:
    enabled: true
    readOnly: false
    features:
      - "kafka_admin_kafka_connect"
    options:
      timeout: "30s"
      
  pulsar_functions:
    enabled: true
    readOnly: true
    features:
      - "pulsar_admin_functions"
    options:
      maxRetries: 3
```

## Creating Custom Builders

### 1. Implement the ToolBuilder Interface

```go
type MyToolBuilder struct {
    *builders.BaseToolBuilder
}

func NewMyToolBuilder() *MyToolBuilder {
    metadata := builders.ToolMetadata{
        Name:        "my_tool",
        Version:     "1.0.0",
        Description: "My custom tool",
        Category:    "custom",
        Tags:        []string{"custom", "example"},
    }
    
    features := []string{"my_feature"}
    
    return &MyToolBuilder{
        BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
    }
}
```

### 2. Implement the BuildTools Method

```go
func (b *MyToolBuilder) BuildTools(ctx context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
    // Validate configuration
    if err := b.Validate(config); err != nil {
        return nil, err
    }
    
    // Check features
    if !b.HasAnyRequiredFeature(config.Features) {
        return nil, nil // Return empty list
    }
    
    // Build tools
    tool := b.buildMyTool()
    handler := b.buildMyHandler(config.ReadOnly)
    
    return []server.ServerTool{
        {
            Tool:    tool,
            Handler: handler,
        },
    }, nil
}
```

### 3. Register the Builder

```go
registry := builders.NewToolRegistry()
registry.Register(NewMyToolBuilder())
```

## Testing

The framework provides complete test coverage:

```bash
cd pkg/mcp/builders
go test -v
```

## Configuration Options

The framework provides various configuration options:

```go
config := builders.NewToolBuildConfig(
    builders.WithReadOnly(true),
    builders.WithFeatures("feature1", "feature2"),
    builders.WithTimeout(30*time.Second),
    builders.WithMaxRetries(3),
    builders.WithOption("custom", "value"),
)
```

## Best Practices

1. **Keep it Simple**: Each builder should only be responsible for one tool or a group of closely related tools
2. **Error Handling**: Provide clear error messages for easy debugging
3. **Configuration Validation**: Validate all required configurations before building
4. **Thread Safety**: Ensure builders can be used safely concurrently
5. **Test Coverage**: Write comprehensive unit tests for each builder

## Architecture Benefits

- **Separation of Concerns**: Tool creation and server registration logic are separated
- **Reusability**: Builders can be reused in different contexts
- **Testability**: Each builder can be tested independently
- **Extensibility**: New tools only need to implement the interface
- **Configuration-Driven**: Supports managing tools through configuration files

## Next Steps

1. Implement specific tool builders (such as Kafka Connect)
2. Integrate into existing tool registration functions
3. Add more configuration options and features
4. Improve documentation and examples

---

**Version**: 1.0.0  
**Created**: 2025-08-01  