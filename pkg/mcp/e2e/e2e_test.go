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
	"os"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/testutil"
	"github.com/stretchr/testify/require"
)

// setupPulsarE2EServer sets up a Pulsar E2E test server.
// It uses environment variables for configuration or falls back to defaults.
func setupPulsarE2EServer(t *testing.T) *testutil.PulsarE2ETestServer {
	adminURL := os.Getenv("PULSAR_ADMIN_URL")
	if adminURL == "" {
		adminURL = "http://localhost:8080"
	}

	serviceURL := os.Getenv("PULSAR_SERVICE_URL")
	if serviceURL == "" {
		serviceURL = "pulsar://localhost:6650"
	}

	server := testutil.NewPulsarE2ETestServer(t, adminURL, serviceURL)

	return server
}

// TestE2EPulsarConnection verifies that we can connect to Pulsar.
func TestE2EPulsarConnection(t *testing.T) {
	t.Parallel()
	server := setupPulsarE2EServer(t)
	ctx := context.Background()

	// Verify connection by listing tools
	response, err := server.Session.ListTools(ctx, &mcp.ListToolsParams{})
	require.NoError(t, err, "failed to list tools")
	require.NotNil(t, response)
	require.NotEmpty(t, response.Tools, "expected at least one tool")

	// Verify pulsar_admin_topic tool is registered
	found := false
	for _, tool := range response.Tools {
		if tool.Name == "pulsar_admin_topic" {
			found = true
			break
		}
	}
	require.True(t, found, "pulsar_admin_topic tool not found")
}

// TestE2EPulsarWaitForReady tests the WaitForReady helper.
func TestE2EPulsarWaitForReady(t *testing.T) {
	t.Parallel()
	adminURL := os.Getenv("PULSAR_ADMIN_URL")
	if adminURL == "" {
		adminURL = "http://localhost:8080"
	}

	serviceURL := os.Getenv("PULSAR_SERVICE_URL")
	if serviceURL == "" {
		serviceURL = "pulsar://localhost:6650"
	}

	helper, err := testutil.NewPulsarTestHelper(adminURL, serviceURL)
	require.NoError(t, err, "failed to create pulsar helper")
	t.Cleanup(func() { helper.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	err = helper.WaitForReady(ctx)
	require.NoError(t, err, "pulsar not ready")
}

// cleanupTestTopic is a helper to cleanup a test topic.
func cleanupTestTopic(t *testing.T, helper *testutil.PulsarTestHelper, topic string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := helper.CleanupTopic(ctx, topic); err != nil {
		t.Logf("Warning: failed to cleanup topic %s: %v", topic, err)
	}
}

// generateTestTopicName generates a unique test topic name.
func generateTestTopicName() string {
	return testutil.GenerateTestTopicName("test")
}
