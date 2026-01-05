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

//go:build e2e

package e2e_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/require"
)

// TestPulsarTopicList tests listing topics in a namespace.
func TestPulsarTopicList(t *testing.T) {
	t.Parallel()

	server := setupPulsarE2EServer(t)
	ctx := context.Background()

	// List topics in public/default namespace
	response, err := server.Session.CallTool(ctx, &mcp.CallToolParams{
		Name: "pulsar_admin_topic",
		Arguments: map[string]any{
			"resource":  "topics",
			"operation": "list",
			"namespace": "public/default",
		},
	})

	require.NoError(t, err, "failed to call pulsar_admin_topic tool")
	require.NotNil(t, response)
	require.False(t, response.IsError, "expected success, got error: %s", getErrorText(response))
	require.Len(t, response.Content, 1, "expected one content item")

	textContent, ok := response.Content[0].(*mcp.TextContent)
	require.True(t, ok, "expected TextContent")

	// Verify response is a JSON array
	var topics []string
	err = json.Unmarshal([]byte(textContent.Text), &topics)
	require.NoError(t, err, "expected JSON array response")
}

// TestPulsarTopicCreateGetDelete tests creating, getting, and deleting a topic.
func TestPulsarTopicCreateGetDelete(t *testing.T) {
	t.Parallel()

	server := setupPulsarE2EServer(t)
	ctx := context.Background()

	testTopic := generateTestTopicName()
	t.Logf("Using test topic: %s", testTopic)

	// Cleanup on test failure
	t.Cleanup(func() {
		cleanupTestTopic(t, server.PulsarHelper, testTopic)
	})

	// 1. Create topic
	t.Run("Create", func(t *testing.T) {
		createResp, err := server.Session.CallTool(ctx, &mcp.CallToolParams{
			Name: "pulsar_admin_topic",
			Arguments: map[string]any{
				"resource":   "topic",
				"operation":  "create",
				"topic":      testTopic,
				"partitions": 0,
			},
		})

		require.NoError(t, err, "failed to call create operation")
		require.NotNil(t, createResp)
		require.False(t, createResp.IsError, "expected success, got error: %s", getErrorText(createResp))
		require.Len(t, createResp.Content, 1, "expected one content item")

		textContent, ok := createResp.Content[0].(*mcp.TextContent)
		require.True(t, ok, "expected TextContent")

		var result map[string]string
		err = json.Unmarshal([]byte(textContent.Text), &result)
		require.NoError(t, err)
		require.Equal(t, "created", result["status"])
		require.Equal(t, testTopic, result["topic"])
	})

	// 2. Get topic metadata
	t.Run("Get", func(t *testing.T) {
		getResp, err := server.Session.CallTool(ctx, &mcp.CallToolParams{
			Name: "pulsar_admin_topic",
			Arguments: map[string]any{
				"resource":  "topic",
				"operation": "get",
				"topic":     testTopic,
			},
		})

		require.NoError(t, err, "failed to call get operation")
		require.NotNil(t, getResp)
		require.False(t, getResp.IsError, "expected success, got error: %s", getErrorText(getResp))
		require.Len(t, getResp.Content, 1, "expected one content item")

		textContent, ok := getResp.Content[0].(*mcp.TextContent)
		require.True(t, ok, "expected TextContent")

		var metadata map[string]interface{}
		err = json.Unmarshal([]byte(textContent.Text), &metadata)
		require.NoError(t, err)
		require.Equal(t, testTopic, metadata["name"])
	})

	// 3. Delete topic
	t.Run("Delete", func(t *testing.T) {
		deleteResp, err := server.Session.CallTool(ctx, &mcp.CallToolParams{
			Name: "pulsar_admin_topic",
			Arguments: map[string]any{
				"resource":  "topic",
				"operation": "delete",
				"topic":     testTopic,
			},
		})

		require.NoError(t, err, "failed to call delete operation")
		require.NotNil(t, deleteResp)
		require.False(t, deleteResp.IsError, "expected success, got error: %s", getErrorText(deleteResp))
		require.Len(t, deleteResp.Content, 1, "expected one content item")

		textContent, ok := deleteResp.Content[0].(*mcp.TextContent)
		require.True(t, ok, "expected TextContent")

		var result map[string]string
		err = json.Unmarshal([]byte(textContent.Text), &result)
		require.NoError(t, err)
		require.Equal(t, "deleted", result["status"])
		require.Equal(t, testTopic, result["topic"])
	})
}

// TestPulsarTopicCreateWithPartitions tests creating a partitioned topic.
func TestPulsarTopicCreateWithPartitions(t *testing.T) {
	t.Parallel()

	server := setupPulsarE2EServer(t)
	ctx := context.Background()

	testTopic := generateTestTopicName()
	t.Logf("Using test topic: %s", testTopic)

	t.Cleanup(func() {
		cleanupTestTopic(t, server.PulsarHelper, testTopic)
	})

	response, err := server.Session.CallTool(ctx, &mcp.CallToolParams{
		Name: "pulsar_admin_topic",
		Arguments: map[string]any{
			"resource":   "topic",
			"operation":  "create",
			"topic":      testTopic,
			"partitions": 3,
		},
	})

	require.NoError(t, err, "failed to call create operation")
	require.NotNil(t, response)
	require.False(t, response.IsError, "expected success, got error: %s", getErrorText(response))
	require.Len(t, response.Content, 1, "expected one content item")

	// Verify the topic was created with partitions
	getResp, err := server.Session.CallTool(ctx, &mcp.CallToolParams{
		Name: "pulsar_admin_topic",
		Arguments: map[string]any{
			"resource":  "topic",
			"operation": "get",
			"topic":     testTopic,
		},
	})

	require.NoError(t, err, "failed to call get operation")
	require.False(t, getResp.IsError, "expected success, got error: %s", getErrorText(getResp))

	textContent, ok := getResp.Content[0].(*mcp.TextContent)
	require.True(t, ok, "expected TextContent")

	var metadata map[string]interface{}
	err = json.Unmarshal([]byte(textContent.Text), &metadata)
	require.NoError(t, err)
	require.Equal(t, testTopic, metadata["name"])

	// Partitions field should be present and > 0
	partitions, ok := metadata["partitions"]
	require.True(t, ok, "expected partitions field")
	require.Equal(t, float64(3), partitions)
}

// TestPulsarTopicErrorCases tests error handling for invalid requests.
func TestPulsarTopicErrorCases(t *testing.T) {
	t.Parallel()

	server := setupPulsarE2EServer(t)
	ctx := context.Background()

	t.Run("MissingRequiredParameter", func(t *testing.T) {
		response, err := server.Session.CallTool(ctx, &mcp.CallToolParams{
			Name: "pulsar_admin_topic",
			Arguments: map[string]any{
				// Missing "resource" and "operation"
			},
		})

		require.NoError(t, err)
		require.True(t, response.IsError, "expected error for missing parameters")
	})

	t.Run("InvalidOperation", func(t *testing.T) {
		response, err := server.Session.CallTool(ctx, &mcp.CallToolParams{
			Name: "pulsar_admin_topic",
			Arguments: map[string]any{
				"resource":  "topic",
				"operation": "invalid_op",
			},
		})

		require.NoError(t, err)
		require.True(t, response.IsError, "expected error for invalid operation")
		require.Contains(t, getErrorText(response), "unsupported operation")
	})

	t.Run("CreateMissingTopic", func(t *testing.T) {
		response, err := server.Session.CallTool(ctx, &mcp.CallToolParams{
			Name: "pulsar_admin_topic",
			Arguments: map[string]any{
				"resource":  "topic",
				"operation": "create",
				// Missing "topic"
			},
		})

		require.NoError(t, err)
		require.True(t, response.IsError, "expected error for missing topic")
		require.Contains(t, getErrorText(response), "topic is required")
	})
}

// getErrorText extracts error text from a CallToolResult.
func getErrorText(response *mcp.CallToolResult) string {
	if len(response.Content) > 0 {
		if textContent, ok := response.Content[0].(*mcp.TextContent); ok {
			return textContent.Text
		}
	}
	return ""
}
