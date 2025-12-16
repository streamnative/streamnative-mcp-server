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
	"slices"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// FeatureChecker defines the interface for checking feature requirements
// It provides methods to determine if required features are available
type FeatureChecker interface {
	// HasAnyRequiredFeature checks if any of the required features are present in the provided list
	HasAnyRequiredFeature(features []string) bool
}

// ToolBuilder defines the interface that all tool builders must implement
// It specifies the methods required for building and managing MCP tools
type ToolBuilder interface {
	// GetName returns the builder name
	GetName() string

	// GetRequiredFeatures returns the list of required features
	GetRequiredFeatures() []string

	// BuildTools builds and returns a list of server tools
	BuildTools(ctx context.Context, config ToolBuildConfig) ([]server.ServerTool, error)

	// Validate validates the builder configuration
	Validate(config ToolBuildConfig) error

	// Embed FeatureChecker interface
	FeatureChecker
}

// ToolBuildConfig contains all configuration information needed to build tools
// It specifies build parameters such as read-only mode, features, and options
type ToolBuildConfig struct {
	ReadOnly bool                   `json:"readOnly"`
	Features []string               `json:"features"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

// ToolMetadata describes the basic information and attributes of a tool
// It contains descriptive information about the tool builder
type ToolMetadata struct {
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Description  string   `json:"description"`
	Category     string   `json:"category"`
	Dependencies []string `json:"dependencies,omitempty"`
	Tags         []string `json:"tags,omitempty"`
}

// BaseToolBuilder provides common functionality implementation for all builders
// It serves as the foundation for concrete tool builder implementations
type BaseToolBuilder struct {
	metadata         ToolMetadata
	requiredFeatures []string
}

// NewBaseToolBuilder creates a new base tool builder instance
func NewBaseToolBuilder(metadata ToolMetadata, features []string) *BaseToolBuilder {
	return &BaseToolBuilder{
		metadata:         metadata,
		requiredFeatures: features,
	}
}

// GetName returns the builder name
func (b *BaseToolBuilder) GetName() string {
	return b.metadata.Name
}

// GetRequiredFeatures returns the list of required features
func (b *BaseToolBuilder) GetRequiredFeatures() []string {
	return b.requiredFeatures
}

// GetMetadata returns the tool metadata
func (b *BaseToolBuilder) GetMetadata() ToolMetadata {
	return b.metadata
}

// Validate validates the builder configuration
// It checks if the configuration contains at least one required feature
func (b *BaseToolBuilder) Validate(config ToolBuildConfig) error {
	return b.validateFeatures(config.Features)
}

// validateFeatures validates the feature configuration
func (b *BaseToolBuilder) validateFeatures(features []string) error {
	hasAny := false
	for _, required := range b.requiredFeatures {
		if slices.Contains(features, required) {
			hasAny = true
			break
		}
	}
	if !hasAny {
		return fmt.Errorf("none of required features found: %v", b.requiredFeatures)
	}
	return nil
}

// HasAnyRequiredFeature checks if any required feature is present
func (b *BaseToolBuilder) HasAnyRequiredFeature(features []string) bool {
	for _, required := range b.requiredFeatures {
		if slices.Contains(features, required) {
			return true
		}
	}
	return false
}

// ToolHandlerFunc defines the tool handler function type
// It maintains consistency with server.ToolHandlerFunc
type ToolHandlerFunc func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
