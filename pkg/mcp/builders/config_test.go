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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToolsConfig(t *testing.T) {
	t.Run("DefaultToolsConfig", func(t *testing.T) {
		config := DefaultToolsConfig()
		require.NotNil(t, config)
		assert.NotEmpty(t, config.Tools)
		assert.Contains(t, config.Tools, "kafka_connect")
		assert.Contains(t, config.Tools, "pulsar_functions")
	})

	t.Run("ToToolBuildConfigs", func(t *testing.T) {
		config := &ToolsConfig{
			Tools: map[string]ToolConfig{
				"tool1": {
					Enabled:  true,
					ReadOnly: false,
					Features: []string{"feature1"},
					Options:  map[string]interface{}{"key": "value"},
				},
				"tool2": {
					Enabled:  false, // Should be excluded
					ReadOnly: true,
					Features: []string{"feature2"},
				},
				"tool3": {
					Enabled:  true,
					ReadOnly: true,
					Features: []string{"feature3"},
				},
			},
		}

		buildConfigs := config.ToToolBuildConfigs()
		assert.Len(t, buildConfigs, 2) // Only enabled tools
		assert.Contains(t, buildConfigs, "tool1")
		assert.Contains(t, buildConfigs, "tool3")
		assert.NotContains(t, buildConfigs, "tool2") // Disabled tool excluded

		assert.False(t, buildConfigs["tool1"].ReadOnly)
		assert.True(t, buildConfigs["tool3"].ReadOnly)
	})

	t.Run("GetEnabledTools", func(t *testing.T) {
		config := &ToolsConfig{
			Tools: map[string]ToolConfig{
				"enabled1":  {Enabled: true},
				"disabled1": {Enabled: false},
				"enabled2":  {Enabled: true},
			},
		}

		enabled := config.GetEnabledTools()
		assert.Len(t, enabled, 2)
		assert.Contains(t, enabled, "enabled1")
		assert.Contains(t, enabled, "enabled2")
		assert.NotContains(t, enabled, "disabled1")
	})

	t.Run("IsToolEnabled", func(t *testing.T) {
		config := &ToolsConfig{
			Tools: map[string]ToolConfig{
				"enabled_tool":  {Enabled: true},
				"disabled_tool": {Enabled: false},
			},
		}

		assert.True(t, config.IsToolEnabled("enabled_tool"))
		assert.False(t, config.IsToolEnabled("disabled_tool"))
		assert.False(t, config.IsToolEnabled("nonexistent_tool"))
	})

	t.Run("SetToolEnabled", func(t *testing.T) {
		config := &ToolsConfig{}

		config.SetToolEnabled("new_tool", true)
		assert.True(t, config.IsToolEnabled("new_tool"))

		config.SetToolEnabled("new_tool", false)
		assert.False(t, config.IsToolEnabled("new_tool"))
	})

	t.Run("Validate_Success", func(t *testing.T) {
		config := &ToolsConfig{
			Tools: map[string]ToolConfig{
				"valid_tool": {
					Enabled:  true,
					Features: []string{"feature1"},
				},
				"disabled_tool": {
					Enabled:  false,
					Features: []string{}, // Empty features OK for disabled tools
				},
			},
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("Validate_EmptyConfig", func(t *testing.T) {
		config := &ToolsConfig{}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no tools configured")
	})

	t.Run("Validate_InvalidTool", func(t *testing.T) {
		config := &ToolsConfig{
			Tools: map[string]ToolConfig{
				"invalid_tool": {
					Enabled:  true,
					Features: []string{}, // Empty features for enabled tool
				},
			},
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no features specified")
	})
}

func TestToolConfig(t *testing.T) {
	t.Run("Validate_Success", func(t *testing.T) {
		config := ToolConfig{
			Features: []string{"feature1", "feature2"},
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("Validate_NoFeatures", func(t *testing.T) {
		config := ToolConfig{
			Features: []string{},
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no features specified")
	})
}

func TestConfigOptions(t *testing.T) {
	t.Run("WithReadOnly", func(t *testing.T) {
		config := NewToolBuildConfig(WithReadOnly(true))
		assert.True(t, config.ReadOnly)

		config = NewToolBuildConfig(WithReadOnly(false))
		assert.False(t, config.ReadOnly)
	})

	t.Run("WithFeatures", func(t *testing.T) {
		config := NewToolBuildConfig(WithFeatures("feature1", "feature2", "feature3"))
		assert.Equal(t, []string{"feature1", "feature2", "feature3"}, config.Features)
	})

	t.Run("WithOption", func(t *testing.T) {
		config := NewToolBuildConfig(
			WithOption("key1", "value1"),
			WithOption("key2", 42),
		)

		assert.Equal(t, "value1", config.Options["key1"])
		assert.Equal(t, 42, config.Options["key2"])
	})

	t.Run("WithTimeout", func(t *testing.T) {
		timeout := 30 * time.Second
		config := NewToolBuildConfig(WithTimeout(timeout))

		assert.Equal(t, timeout.String(), config.Options["timeout"])
	})

	t.Run("WithMaxRetries", func(t *testing.T) {
		config := NewToolBuildConfig(WithMaxRetries(5))
		assert.Equal(t, 5, config.Options["maxRetries"])
	})

	t.Run("CombinedOptions", func(t *testing.T) {
		config := NewToolBuildConfig(
			WithReadOnly(true),
			WithFeatures("feature1", "feature2"),
			WithTimeout(45*time.Second),
			WithMaxRetries(3),
			WithOption("custom", "value"),
		)

		assert.True(t, config.ReadOnly)
		assert.Equal(t, []string{"feature1", "feature2"}, config.Features)
		assert.Equal(t, (45 * time.Second).String(), config.Options["timeout"])
		assert.Equal(t, 3, config.Options["maxRetries"])
		assert.Equal(t, "value", config.Options["custom"])
	})
}
