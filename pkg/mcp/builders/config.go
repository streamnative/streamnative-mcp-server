// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package builders

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// ToolsConfig represents the structure of the tools configuration file
type ToolsConfig struct {
	Tools map[string]ToolConfig `yaml:"tools"`
}

// ToolConfig represents the configuration for a single tool
type ToolConfig struct {
	Enabled  bool                   `yaml:"enabled"`
	ReadOnly bool                   `yaml:"readOnly"`
	Features []string               `yaml:"features"`
	Options  map[string]interface{} `yaml:"options,omitempty"`
}

// LoadToolsConfig loads tool configuration from a file
func LoadToolsConfig(filename string) (*ToolsConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", filename, err)
	}

	var config ToolsConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", filename, err)
	}

	return &config, nil
}

// SaveToolsConfig saves tool configuration to a file
func SaveToolsConfig(config *ToolsConfig, filename string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", filename, err)
	}

	return nil
}

// ToToolBuildConfigs converts tool configuration to build configurations
func (tc *ToolsConfig) ToToolBuildConfigs() map[string]ToolBuildConfig {
	configs := make(map[string]ToolBuildConfig)

	for name, toolConfig := range tc.Tools {
		if !toolConfig.Enabled {
			continue
		}

		configs[name] = ToolBuildConfig{
			ReadOnly: toolConfig.ReadOnly,
			Features: toolConfig.Features,
			Options:  toolConfig.Options,
		}
	}

	return configs
}

// GetEnabledTools returns the names of all enabled tools
func (tc *ToolsConfig) GetEnabledTools() []string {
	var enabled []string
	for name, config := range tc.Tools {
		if config.Enabled {
			enabled = append(enabled, name)
		}
	}
	return enabled
}

// IsToolEnabled checks if the specified tool is enabled
func (tc *ToolsConfig) IsToolEnabled(name string) bool {
	if config, exists := tc.Tools[name]; exists {
		return config.Enabled
	}
	return false
}

// SetToolEnabled sets the enabled status of a tool
func (tc *ToolsConfig) SetToolEnabled(name string, enabled bool) {
	if tc.Tools == nil {
		tc.Tools = make(map[string]ToolConfig)
	}

	config := tc.Tools[name]
	config.Enabled = enabled
	tc.Tools[name] = config
}

// DefaultToolsConfig creates a default tool configuration
func DefaultToolsConfig() *ToolsConfig {
	return &ToolsConfig{
		Tools: map[string]ToolConfig{
			"kafka_connect": {
				Enabled:  true,
				ReadOnly: false,
				Features: []string{"kafka_admin_kafka_connect"},
				Options: map[string]interface{}{
					"timeout": "30s",
				},
			},
			"pulsar_functions": {
				Enabled:  true,
				ReadOnly: false,
				Features: []string{"pulsar_admin_functions"},
				Options: map[string]interface{}{
					"maxRetries": 3,
				},
			},
		},
	}
}

// ConfigOption defines the function type for configuration options
type ConfigOption func(*ToolBuildConfig)

// WithReadOnly sets the read-only mode
func WithReadOnly(readOnly bool) ConfigOption {
	return func(config *ToolBuildConfig) {
		config.ReadOnly = readOnly
	}
}

// WithFeatures sets the feature list
func WithFeatures(features ...string) ConfigOption {
	return func(config *ToolBuildConfig) {
		config.Features = features
	}
}

// WithOption sets a configuration option
func WithOption(key string, value interface{}) ConfigOption {
	return func(config *ToolBuildConfig) {
		if config.Options == nil {
			config.Options = make(map[string]interface{})
		}
		config.Options[key] = value
	}
}

// WithTimeout sets the timeout option
func WithTimeout(timeout time.Duration) ConfigOption {
	return WithOption("timeout", timeout.String())
}

// WithMaxRetries sets the maximum retry count option
func WithMaxRetries(maxRetries int) ConfigOption {
	return WithOption("maxRetries", maxRetries)
}

// NewToolBuildConfig creates a new tool build configuration
func NewToolBuildConfig(options ...ConfigOption) ToolBuildConfig {
	config := ToolBuildConfig{
		ReadOnly: false,
		Features: []string{},
		Options:  make(map[string]interface{}),
	}

	for _, option := range options {
		option(&config)
	}

	return config
}

// Validate validates the tool configuration
func (tc *ToolConfig) Validate() error {
	if len(tc.Features) == 0 {
		return fmt.Errorf("no features specified")
	}
	return nil
}

// Validate validates the entire tools configuration file
func (tc *ToolsConfig) Validate() error {
	if len(tc.Tools) == 0 {
		return fmt.Errorf("no tools configured")
	}

	for name, config := range tc.Tools {
		if config.Enabled {
			if err := config.Validate(); err != nil {
				return fmt.Errorf("invalid config for tool %s: %w", name, err)
			}
		}
	}

	return nil
}
