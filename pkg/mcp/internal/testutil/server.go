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

package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/require"
)

// PulsarE2ETestServer encapsulates the E2E test components.
type PulsarE2ETestServer struct {
	Server        *mcp.Server
	Session       *mcp.ClientSession
	ServerCtx     context.Context
	ClientCtx     context.Context
	ServerCancel  context.CancelFunc
	ClientCancel  context.CancelFunc
	PulsarHelper  *PulsarTestHelper
	TestNamespace string
}

// NewPulsarE2ETestServer creates a new E2E test server with in-memory transport.
func NewPulsarE2ETestServer(t testing.TB, adminURL, serviceURL string) *PulsarE2ETestServer {
	// Create Pulsar helper
	pulsarHelper, err := NewPulsarTestHelper(adminURL, serviceURL)
	require.NoError(t, err, "failed to create pulsar helper")

	// Wait for Pulsar to be ready
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	require.NoError(t, pulsarHelper.WaitForReady(ctx), "pulsar not ready")

	// Ensure test namespace exists
	require.NoError(t, pulsarHelper.EnsureNamespace(context.Background(), "public", "default"))

	// Create in-memory transports
	serverTransport, clientTransport := mcp.NewInMemoryTransports()

	// Create go-sdk server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "test-server",
		Version: "0.0.1",
	}, &mcp.ServerOptions{})

	// Register test tools
	registerPulsarTopicTools(t, server, pulsarHelper)

	// Start server
	serverCtx, serverCancel := context.WithCancel(context.Background())
	go func() {
		_ = server.Run(serverCtx, serverTransport)
	}()

	// Create client
	clientCtx, clientCancel := context.WithCancel(context.Background())
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "test-client",
		Version: "0.0.1",
	}, nil)

	session, err := client.Connect(clientCtx, clientTransport, nil)
	require.NoError(t, err, "failed to connect client")

	t.Cleanup(func() {
		_ = session.Close()
		serverCancel()
		clientCancel()
		pulsarHelper.Close()
	})

	return &PulsarE2ETestServer{
		Server:        server,
		Session:       session,
		ServerCtx:     serverCtx,
		ClientCtx:     clientCtx,
		ServerCancel:  serverCancel,
		ClientCancel:  clientCancel,
		PulsarHelper:  pulsarHelper,
		TestNamespace: "public/default",
	}
}

// registerPulsarTopicTools registers the pulsar_admin_topic tool for E2E testing.
// This is a simplified version that directly calls Pulsar admin API.
func registerPulsarTopicTools(t testing.TB, server *mcp.Server, helper *PulsarTestHelper) {
	tool := &mcp.Tool{
		Name:        "pulsar_admin_topic",
		Description: "Manage Pulsar topics. Supports list, get, create, and delete operations.",
		InputSchema: json.RawMessage(`{
			"$schema": "http://json-schema.org/draft-07/schema#",
			"type": "object",
			"properties": {
				"resource": {
					"type": "string",
					"description": "Resource to operate on: 'topic' or 'topics'"
				},
				"operation": {
					"type": "string",
					"description": "Operation: 'list', 'get', 'create', 'delete'"
				},
				"topic": {
					"type": "string",
					"description": "Fully qualified topic name (e.g., persistent://public/default/test)"
				},
				"namespace": {
					"type": "string",
					"description": "Namespace name (e.g., public/default)"
				},
				"partitions": {
					"type": "integer",
					"description": "Number of partitions (0 for non-partitioned)"
				}
			},
			"required": ["resource", "operation"]
		}`),
	}

	// Handler using go-sdk signature: func(context.Context, *CallToolRequest) (*CallToolResult, error)
	// req.Params.Arguments is json.RawMessage that needs unmarshaling
	handler := func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var args map[string]any
		_ = json.Unmarshal(req.Params.Arguments, &args)
		resource, _ := args["resource"].(string)
		operation, _ := args["operation"].(string)

		switch operation {
		case "list":
			if resource != "topics" {
				return newToolResultError("list operation requires 'topics' resource"), nil
			}
			namespace, _ := args["namespace"].(string)
			if namespace == "" {
				namespace = "public/default"
			}

			topics, err := helper.ListTopics(namespace)
			if err != nil {
				return newToolResultError(fmt.Sprintf("failed to list topics: %v", err)), nil
			}

			result, _ := json.Marshal(topics)
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: string(result)},
				},
			}, nil

		case "get":
			topic, _ := args["topic"].(string)
			if topic == "" {
				return newToolResultError("topic is required for get operation"), nil
			}

			metadata, err := helper.GetTopicMetadata(topic)
			if err != nil {
				return newToolResultError(fmt.Sprintf("failed to get topic: %v", err)), nil
			}

			result, _ := json.Marshal(metadata)
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: string(result)},
				},
			}, nil

		case "create":
			topic, _ := args["topic"].(string)
			if topic == "" {
				return newToolResultError("topic is required for create operation"), nil
			}
			partitionsVal, _ := args["partitions"].(float64)
			partitions := int(partitionsVal)

			if err := helper.CreateTopic(ctx, topic, partitions); err != nil {
				return newToolResultError(fmt.Sprintf("failed to create topic: %v", err)), nil
			}

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf(`{"status": "created", "topic": "%s"}`, topic)},
				},
			}, nil

		case "delete":
			topic, _ := args["topic"].(string)
			if topic == "" {
				return newToolResultError("topic is required for delete operation"), nil
			}

			if err := helper.DeleteTopic(ctx, topic); err != nil {
				return newToolResultError(fmt.Sprintf("failed to delete topic: %v", err)), nil
			}

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf(`{"status": "deleted", "topic": "%s"}`, topic)},
				},
			}, nil

		default:
			return newToolResultError(fmt.Sprintf("unsupported operation: %s", operation)), nil
		}
	}

	server.AddTool(tool, handler)
}

// newToolResultError creates an error tool result.
func newToolResultError(message string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: message,
			},
		},
		IsError: true,
	}
}
