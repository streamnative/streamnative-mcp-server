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

	"github.com/stretchr/testify/assert"
)

func TestBaseToolBuilder(t *testing.T) {
	metadata := ToolMetadata{
		Name:        "test_tool",
		Version:     "1.0.0",
		Description: "Test tool for unit testing",
		Category:    "test",
		Tags:        []string{"test", "unit"},
	}

	features := []string{"feature1", "feature2", "feature3"}
	builder := NewBaseToolBuilder(metadata, features)

	t.Run("GetName", func(t *testing.T) {
		assert.Equal(t, "test_tool", builder.GetName())
	})

	t.Run("GetRequiredFeatures", func(t *testing.T) {
		assert.Equal(t, features, builder.GetRequiredFeatures())
	})

	t.Run("GetMetadata", func(t *testing.T) {
		assert.Equal(t, metadata, builder.GetMetadata())
	})

	t.Run("Validate_Success", func(t *testing.T) {
		config := ToolBuildConfig{
			Features: []string{"feature1", "other_feature"},
		}

		err := builder.Validate(config)
		assert.NoError(t, err)
	})

	t.Run("Validate_MissingFeatures", func(t *testing.T) {
		config := ToolBuildConfig{
			Features: []string{"other_feature", "another_feature"},
		}

		err := builder.Validate(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "none of required features found")
	})

	t.Run("HasAnyRequiredFeature_True", func(t *testing.T) {
		features := []string{"feature1", "other_feature"}
		assert.True(t, builder.HasAnyRequiredFeature(features))
	})

	t.Run("HasAnyRequiredFeature_False", func(t *testing.T) {
		features := []string{"other_feature", "another_feature"}
		assert.False(t, builder.HasAnyRequiredFeature(features))
	})

	t.Run("HasAnyRequiredFeature_Empty", func(t *testing.T) {
		features := []string{}
		assert.False(t, builder.HasAnyRequiredFeature(features))
	})
}

func TestToolBuildConfig(t *testing.T) {
	t.Run("NewToolBuildConfig_Default", func(t *testing.T) {
		config := NewToolBuildConfig()

		assert.False(t, config.ReadOnly)
		assert.Empty(t, config.Features)
		assert.NotNil(t, config.Options)
	})

	t.Run("NewToolBuildConfig_WithOptions", func(t *testing.T) {
		config := NewToolBuildConfig(
			WithReadOnly(true),
			WithFeatures("feature1", "feature2"),
			WithOption("key1", "value1"),
			WithTimeout(30),
		)

		assert.True(t, config.ReadOnly)
		assert.Equal(t, []string{"feature1", "feature2"}, config.Features)
		assert.Equal(t, "value1", config.Options["key1"])
		assert.NotNil(t, config.Options["timeout"])
	})
}

func TestToolMetadata(t *testing.T) {
	t.Run("CompleteMetadata", func(t *testing.T) {
		metadata := ToolMetadata{
			Name:         "complete_tool",
			Version:      "2.1.0",
			Description:  "A complete tool with all metadata",
			Category:     "admin",
			Dependencies: []string{"dep1", "dep2"},
			Tags:         []string{"tag1", "tag2", "tag3"},
		}

		assert.Equal(t, "complete_tool", metadata.Name)
		assert.Equal(t, "2.1.0", metadata.Version)
		assert.Equal(t, "A complete tool with all metadata", metadata.Description)
		assert.Equal(t, "admin", metadata.Category)
		assert.Equal(t, []string{"dep1", "dep2"}, metadata.Dependencies)
		assert.Equal(t, []string{"tag1", "tag2", "tag3"}, metadata.Tags)
	})

	t.Run("MinimalMetadata", func(t *testing.T) {
		metadata := ToolMetadata{
			Name:        "minimal_tool",
			Version:     "1.0.0",
			Description: "A minimal tool",
		}

		assert.Equal(t, "minimal_tool", metadata.Name)
		assert.Empty(t, metadata.Dependencies)
		assert.Empty(t, metadata.Tags)
	})
}
