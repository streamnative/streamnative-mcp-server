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
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer_InitializeCompatibility(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		protocolVersion string
		expectedVersion string
	}{
		{
			name:            "explicit latest protocol",
			protocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			expectedVersion: mcp.LATEST_PROTOCOL_VERSION,
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
