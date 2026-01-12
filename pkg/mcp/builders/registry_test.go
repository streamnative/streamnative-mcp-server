// Copyright 2025 StreamNative
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package builders

import (
	"context"
	"fmt"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockToolBuilder is a mock builder for testing purposes
type MockToolBuilder struct {
	name     string
	features []string
	tools    []server.ServerTool
	err      error
	metadata ToolMetadata
}

func NewMockToolBuilder(name string, features []string) *MockToolBuilder {
	return &MockToolBuilder{
		name:     name,
		features: features,
		metadata: ToolMetadata{
			Name:        name,
			Version:     "1.0.0",
			Description: fmt.Sprintf("Mock tool builder for %s", name),
			Category:    "test",
		},
		tools: []server.ServerTool{
			{
				Tool: mcp.NewTool(name,
					mcp.WithDescription(fmt.Sprintf("Mock tool %s", name)),
				),
				Handler: func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
					return mcp.NewToolResultText(fmt.Sprintf("Mock response from %s", name)), nil
				},
			},
		},
	}
}

func (m *MockToolBuilder) GetName() string {
	return m.name
}

func (m *MockToolBuilder) GetRequiredFeatures() []string {
	return m.features
}

func (m *MockToolBuilder) BuildTools(_ context.Context, _ ToolBuildConfig) ([]server.ServerTool, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.tools, nil
}

func (m *MockToolBuilder) Validate(config ToolBuildConfig) error {
	// Simple validation: check if any required feature is present
	for _, required := range m.features {
		for _, provided := range config.Features {
			if required == provided {
				return nil
			}
		}
	}
	return fmt.Errorf("no matching features found")
}

func (m *MockToolBuilder) HasAnyRequiredFeature(features []string) bool {
	for _, required := range m.features {
		for _, provided := range features {
			if required == provided {
				return true
			}
		}
	}
	return false
}

func (m *MockToolBuilder) GetMetadata() ToolMetadata {
	return m.metadata
}

func (m *MockToolBuilder) SetError(err error) {
	m.err = err
}

func TestToolRegistry(t *testing.T) {
	t.Run("NewToolRegistry", func(t *testing.T) {
		registry := NewToolRegistry()
		assert.NotNil(t, registry)
		assert.Equal(t, 0, registry.Count())
	})

	t.Run("Register_Success", func(t *testing.T) {
		registry := NewToolRegistry()
		builder := NewMockToolBuilder("test_tool", []string{"test_feature"})

		///nolint:errcheck
		err := registry.Register(builder)
		assert.NoError(t, err)
		assert.Equal(t, 1, registry.Count())
	})

	t.Run("Register_Duplicate", func(t *testing.T) {
		registry := NewToolRegistry()
		builder1 := NewMockToolBuilder("duplicate_tool", []string{"feature1"})
		builder2 := NewMockToolBuilder("duplicate_tool", []string{"feature2"})

		err1 := registry.Register(builder1)
		err2 := registry.Register(builder2)

		assert.NoError(t, err1)
		assert.Error(t, err2)
		assert.Contains(t, err2.Error(), "already registered")
		assert.Equal(t, 1, registry.Count())
	})

	t.Run("MustRegister_Success", func(t *testing.T) {
		registry := NewToolRegistry()
		builder := NewMockToolBuilder("must_register_tool", []string{"feature"})

		assert.NotPanics(t, func() {
			registry.MustRegister(builder)
		})
		assert.Equal(t, 1, registry.Count())
	})

	t.Run("MustRegister_Panic", func(t *testing.T) {
		registry := NewToolRegistry()
		builder := NewMockToolBuilder("panic_tool", []string{"feature"})

		_ = registry.Register(builder) // First registration
		assert.Panics(t, func() {
			registry.MustRegister(builder) // Duplicate registration should panic
		})
	})

	t.Run("GetBuilder_Success", func(t *testing.T) {
		registry := NewToolRegistry()
		builder := NewMockToolBuilder("get_test_tool", []string{"feature"})
		_ = registry.Register(builder)

		retrieved, exists := registry.GetBuilder("get_test_tool")
		assert.True(t, exists)
		assert.Equal(t, builder, retrieved)
	})

	t.Run("GetBuilder_NotFound", func(t *testing.T) {
		registry := NewToolRegistry()

		_, exists := registry.GetBuilder("nonexistent_tool")
		assert.False(t, exists)
	})

	t.Run("ListBuilders", func(t *testing.T) {
		registry := NewToolRegistry()
		builder1 := NewMockToolBuilder("tool1", []string{"feature1"})
		builder2 := NewMockToolBuilder("tool2", []string{"feature2"})

		_ = registry.Register(builder1)
		_ = registry.Register(builder2)

		names := registry.ListBuilders()
		assert.Len(t, names, 2)
		assert.Contains(t, names, "tool1")
		assert.Contains(t, names, "tool2")
		// Should be sorted
		assert.Equal(t, []string{"tool1", "tool2"}, names)
	})

	t.Run("GetMetadata", func(t *testing.T) {
		registry := NewToolRegistry()
		builder := NewMockToolBuilder("metadata_tool", []string{"feature"})
		_ = registry.Register(builder)

		metadata, exists := registry.GetMetadata("metadata_tool")
		assert.True(t, exists)
		assert.Equal(t, "metadata_tool", metadata.Name)
		assert.Equal(t, "1.0.0", metadata.Version)
	})

	t.Run("ListMetadata", func(t *testing.T) {
		registry := NewToolRegistry()
		builder1 := NewMockToolBuilder("tool1", []string{"feature1"})
		builder2 := NewMockToolBuilder("tool2", []string{"feature2"})

		_ = registry.Register(builder1)
		_ = registry.Register(builder2)

		metadata := registry.ListMetadata()
		assert.Len(t, metadata, 2)
		assert.Contains(t, metadata, "tool1")
		assert.Contains(t, metadata, "tool2")
	})

	t.Run("BuildSingle_Success", func(t *testing.T) {
		registry := NewToolRegistry()
		builder := NewMockToolBuilder("single_tool", []string{"test_feature"})
		_ = registry.Register(builder)

		config := ToolBuildConfig{
			Features: []string{"test_feature"},
		}

		tools, err := registry.BuildSingle("single_tool", config)
		require.NoError(t, err)
		assert.Len(t, tools, 1)
		assert.Equal(t, "single_tool", tools[0].Tool.Name)
	})

	t.Run("BuildSingle_NotFound", func(t *testing.T) {
		registry := NewToolRegistry()

		config := ToolBuildConfig{
			Features: []string{"test_feature"},
		}

		_, err := registry.BuildSingle("nonexistent_tool", config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("BuildSingle_ValidationFailed", func(t *testing.T) {
		registry := NewToolRegistry()
		builder := NewMockToolBuilder("validation_tool", []string{"required_feature"})
		_ = registry.Register(builder)

		config := ToolBuildConfig{
			Features: []string{"wrong_feature"},
		}

		_, err := registry.BuildSingle("validation_tool", config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("BuildAll_Success", func(t *testing.T) {
		registry := NewToolRegistry()
		builder1 := NewMockToolBuilder("tool1", []string{"feature1"})
		builder2 := NewMockToolBuilder("tool2", []string{"feature2"})

		_ = registry.Register(builder1)
		_ = registry.Register(builder2)

		configs := map[string]ToolBuildConfig{
			"tool1": {Features: []string{"feature1"}},
			"tool2": {Features: []string{"feature2"}},
		}

		tools, err := registry.BuildAll(configs)
		require.NoError(t, err)
		assert.Len(t, tools, 2)
	})

	t.Run("BuildAll_WithErrors", func(t *testing.T) {
		registry := NewToolRegistry()
		builder1 := NewMockToolBuilder("tool1", []string{"feature1"})
		builder2 := NewMockToolBuilder("tool2", []string{"feature2"})

		builder2.SetError(fmt.Errorf("build error"))

		_ = registry.Register(builder1)
		_ = registry.Register(builder2)

		configs := map[string]ToolBuildConfig{
			"tool1": {Features: []string{"feature1"}},
			"tool2": {Features: []string{"feature2"}},
		}

		tools, err := registry.BuildAll(configs)
		assert.Error(t, err)
		assert.Len(t, tools, 1) // Only tool1 should succeed
	})

	t.Run("BuildAllWithFeatures", func(t *testing.T) {
		registry := NewToolRegistry()
		builder1 := NewMockToolBuilder("tool1", []string{"feature1"})
		builder2 := NewMockToolBuilder("tool2", []string{"feature2"})
		builder3 := NewMockToolBuilder("tool3", []string{"feature3"})

		_ = registry.Register(builder1)
		_ = registry.Register(builder2)
		_ = registry.Register(builder3)

		// Only provide feature1 and feature2
		tools, err := registry.BuildAllWithFeatures(false, []string{"feature1", "feature2"})
		require.NoError(t, err)
		assert.Len(t, tools, 2) // tool3 should be skipped
	})

	t.Run("Clear", func(t *testing.T) {
		registry := NewToolRegistry()
		builder := NewMockToolBuilder("clear_tool", []string{"feature"})
		_ = registry.Register(builder)

		assert.Equal(t, 1, registry.Count())

		registry.Clear()
		assert.Equal(t, 0, registry.Count())
	})

	t.Run("Unregister", func(t *testing.T) {
		registry := NewToolRegistry()
		builder := NewMockToolBuilder("unregister_tool", []string{"feature"})
		_ = registry.Register(builder)

		assert.Equal(t, 1, registry.Count())

		removed := registry.Unregister("unregister_tool")
		assert.True(t, removed)
		assert.Equal(t, 0, registry.Count())

		removed = registry.Unregister("nonexistent_tool")
		assert.False(t, removed)
	})
}
