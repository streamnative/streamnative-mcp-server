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
	"sort"
	"sync"
)

// ToolRegistry manages the registration and building of all tool builders
// It provides centralized management for tool builder registration and tool construction
type ToolRegistry struct {
	mu       sync.RWMutex
	builders map[string]ToolBuilder
	metadata map[string]ToolMetadata
}

// NewToolRegistry creates a new tool registry instance
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		builders: make(map[string]ToolBuilder),
		metadata: make(map[string]ToolMetadata),
	}
}

// Register registers a tool builder
// Returns an error if the builder already exists
func (r *ToolRegistry) Register(builder ToolBuilder) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := builder.GetName()
	if _, exists := r.builders[name]; exists {
		return fmt.Errorf("builder %s already registered", name)
	}

	r.builders[name] = builder

	// Try to get metadata if the builder supports it
	// Use interface method instead of type assertion
	if metadataProvider, ok := builder.(interface{ GetMetadata() ToolMetadata }); ok {
		r.metadata[name] = metadataProvider.GetMetadata()
	}

	return nil
}

// MustRegister registers a tool builder and panics if it fails
func (r *ToolRegistry) MustRegister(builder ToolBuilder) {
	if err := r.Register(builder); err != nil {
		panic(fmt.Sprintf("failed to register builder: %v", err))
	}
}

// GetBuilder retrieves the tool builder with the specified name
func (r *ToolRegistry) GetBuilder(name string) (ToolBuilder, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	builder, exists := r.builders[name]
	return builder, exists
}

// ListBuilders returns a list of all registered builder names
func (r *ToolRegistry) ListBuilders() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.builders))
	for name := range r.builders {
		names = append(names, name)
	}

	// Sort to ensure consistent output
	sort.Strings(names)
	return names
}

// GetMetadata retrieves the metadata for the specified builder
func (r *ToolRegistry) GetMetadata(name string) (ToolMetadata, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	metadata, exists := r.metadata[name]
	return metadata, exists
}

// ListMetadata returns metadata for all registered builders
func (r *ToolRegistry) ListMetadata() map[string]ToolMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]ToolMetadata, len(r.metadata))
	for name, metadata := range r.metadata {
		result[name] = metadata
	}
	return result
}

// BuildAll builds tools for all specified configurations
// Returns all successfully built tools and any errors encountered
func (r *ToolRegistry) BuildAll(configs map[string]ToolBuildConfig) ([]ToolDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var allTools []ToolDefinition
	var errors []error

	for name, config := range configs {
		if builder, exists := r.builders[name]; exists {
			if err := builder.Validate(config); err != nil {
				errors = append(errors, fmt.Errorf("validation failed for %s: %w", name, err))
				continue
			}

			tools, err := builder.BuildTools(context.Background(), config)
			if err != nil {
				errors = append(errors, fmt.Errorf("build failed for %s: %w", name, err))
				continue
			}

			allTools = append(allTools, tools...)
		} else {
			errors = append(errors, fmt.Errorf("builder %s not found", name))
		}
	}

	if len(errors) > 0 {
		return allTools, fmt.Errorf("build errors: %v", errors)
	}

	return allTools, nil
}

// BuildSingle builds tools for a single tool builder
func (r *ToolRegistry) BuildSingle(name string, config ToolBuildConfig) ([]ToolDefinition, error) {
	r.mu.RLock()
	builder, exists := r.builders[name]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("builder %s not found", name)
	}

	if err := builder.Validate(config); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return builder.BuildTools(context.Background(), config)
}

// BuildAllWithFeatures builds all relevant tools based on the feature list
// Automatically creates configuration for each builder
func (r *ToolRegistry) BuildAllWithFeatures(readOnly bool, features []string) ([]ToolDefinition, error) {
	r.mu.RLock()
	builders := make(map[string]ToolBuilder, len(r.builders))
	for name, builder := range r.builders {
		builders[name] = builder
	}
	r.mu.RUnlock()

	var allTools []ToolDefinition
	var errors []error

	for name, builder := range builders {
		config := ToolBuildConfig{
			ReadOnly: readOnly,
			Features: features,
		}

		// Check if the builder needs these features
		if !builder.HasAnyRequiredFeature(features) {
			continue // Skip builders that don't need these features
		}

		if err := builder.Validate(config); err != nil {
			// Validation failure is not an error, just skip
			continue
		}

		tools, err := builder.BuildTools(context.Background(), config)
		if err != nil {
			errors = append(errors, fmt.Errorf("build failed for %s: %w", name, err))
			continue
		}

		allTools = append(allTools, tools...)
	}

	if len(errors) > 0 {
		return allTools, fmt.Errorf("build errors: %v", errors)
	}

	return allTools, nil
}

// Count returns the number of registered builders
func (r *ToolRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.builders)
}

// Clear removes all registered builders
func (r *ToolRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.builders = make(map[string]ToolBuilder)
	r.metadata = make(map[string]ToolMetadata)
}

// Unregister removes the specified tool builder
func (r *ToolRegistry) Unregister(name string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.builders[name]; exists {
		delete(r.builders, name)
		delete(r.metadata, name)
		return true
	}
	return false
}
