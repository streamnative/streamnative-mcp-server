// Copyright 2026 StreamNative
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

package mcp

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExternalServerDoesNotAcquireCloudBindingLease(t *testing.T) {
	srv := NewServer("external", "test", logrus.New())
	ctx := WithSNCloudSession(context.Background(), srv.SNCloudSession)
	entered, release := make(chan struct{}), make(chan struct{})
	releaseRequest := sync.OnceFunc(func() { close(release) })
	t.Cleanup(releaseRequest)
	srv.MCPServer.AddTool(mcp.NewTool("hold"), func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		close(entered)
		<-release
		return mcp.NewToolResultText("done"), nil
	})
	requestDone := make(chan struct{})
	go func() {
		srv.MCPServer.HandleMessage(ctx, json.RawMessage(`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"hold"}}`))
		close(requestDone)
	}()
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("handler did not start")
	}
	closed := make(chan struct{})
	go func() { srv.Close(); close(closed) }()
	select {
	case <-closed:
	case <-time.After(100 * time.Millisecond):
		releaseRequest()
		<-requestDone
		<-closed
		t.Fatal("external mode must not wait for a Cloud binding lease")
	}
	releaseRequest()
	<-requestDone
}

func TestNewServer_InitializeCompatibility(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		protocolVersion string
		expectedVersion string
	}{
		{
			name:            "latest legacy protocol",
			protocolVersion: "2025-11-25",
			expectedVersion: "2025-11-25",
		},
		{
			name:            "modern version cannot be negotiated through initialize",
			protocolVersion: "2026-07-28",
			expectedVersion: "2025-11-25",
		},
		{
			name:            "June 2025 protocol",
			protocolVersion: "2025-06-18",
			expectedVersion: "2025-06-18",
		},
		{
			name:            "March 2025 protocol",
			protocolVersion: "2025-03-26",
			expectedVersion: "2025-03-26",
		},
		{
			name:            "original SSE protocol",
			protocolVersion: "2024-11-05",
			expectedVersion: "2024-11-05",
		},
		{
			name:            "unknown version falls back to latest legacy",
			protocolVersion: "2099-01-01",
			expectedVersion: "2025-11-25",
		},
		{
			name:            "empty protocol keeps backward compatible fallback",
			protocolVersion: "",
			expectedVersion: "2025-03-26",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := NewServer(
				"streamnative-mcp-server",
				"test-version",
				logrus.New(),
				server.WithInstructions("test instructions"),
			)

			result := initializeServer(t, srv.MCPServer, tt.protocolVersion)

			assert.Equal(t, tt.expectedVersion, result.ProtocolVersion)
			assert.Equal(t, "streamnative-mcp-server", result.ServerInfo.Name)
			assert.Equal(t, "test-version", result.ServerInfo.Version)
			assert.Equal(t, "test instructions", result.Instructions)

			require.NotNil(t, result.Capabilities.Resources)
			assert.True(t, result.Capabilities.Resources.Subscribe)
			assert.True(t, result.Capabilities.Resources.ListChanged)

			require.NotNil(t, result.Capabilities.Logging)
			assert.Nil(t, result.Capabilities.Tasks)
			assert.Nil(t, result.Capabilities.Completions)
		})
	}
}

func initializeServer(t *testing.T, srv *server.MCPServer, protocolVersion string) mcp.InitializeResult {
	t.Helper()

	initReq := mcp.JSONRPCRequest{
		JSONRPC: mcp.JSONRPC_VERSION,
		ID:      mcp.NewRequestId(int64(1)),
		Request: mcp.Request{
			Method: string(mcp.MethodInitialize),
		},
		Params: struct {
			ProtocolVersion string                 `json:"protocolVersion"`
			ClientInfo      mcp.Implementation     `json:"clientInfo"`
			Capabilities    mcp.ClientCapabilities `json:"capabilities"`
		}{
			ProtocolVersion: protocolVersion,
			ClientInfo: mcp.Implementation{
				Name:    "test-client",
				Version: "1.0.0",
			},
		},
	}

	messageBytes, err := json.Marshal(initReq)
	require.NoError(t, err)

	response := srv.HandleMessage(context.Background(), messageBytes)
	require.NotNil(t, response)

	jsonRPCResponse, ok := response.(mcp.JSONRPCResponse)
	require.True(t, ok, "expected JSONRPCResponse, got %T", response)

	result, ok := jsonRPCResponse.Result.(mcp.InitializeResult)
	require.True(t, ok, "expected InitializeResult, got %T", jsonRPCResponse.Result)

	return result
}
