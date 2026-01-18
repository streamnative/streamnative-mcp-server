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

package mcp

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// LegacyServer wraps MCP server state and StreamNative sessions for mark3labs/mcp-go.
type LegacyServer struct {
	MCPServer      *server.MCPServer
	KafkaSession   *kafka.Session
	PulsarSession  *pulsar.Session
	SNCloudSession *config.Session
	logger         *logrus.Logger
}

// NewLegacyServer creates a new MCP server with StreamNative integrations.
func NewLegacyServer(name, version string, logger *logrus.Logger, opts ...server.ServerOption) *LegacyServer {
	// Create a new MCP server
	opts = addLegacyOpts(opts...)
	s := server.NewMCPServer(name, version, opts...)
	mcpserver := createSNCloudLegacyServer(s, logger)
	return mcpserver
}

// addLegacyOpts merges default server options with custom options.
func addLegacyOpts(opts ...server.ServerOption) []server.ServerOption {
	defaultOpts := []server.ServerOption{
		server.WithResourceCapabilities(true, true),
		server.WithRecovery(),
		server.WithLogging(),
	}
	opts = append(defaultOpts, opts...)
	return opts
}

// createSNCloudLegacyServer constructs a LegacyServer wrapper for StreamNative Cloud.
func createSNCloudLegacyServer(s *server.MCPServer, logger *logrus.Logger) *LegacyServer {
	mcpserver := &LegacyServer{
		MCPServer:      s,
		logger:         logger,
		SNCloudSession: &config.Session{},
		KafkaSession:   &kafka.Session{},
		PulsarSession:  &pulsar.Session{},
	}

	return mcpserver
}
