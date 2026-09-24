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

	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	context2 "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// Server wraps MCP server state and StreamNative sessions.
type Server struct {
	MCPServer      *server.MCPServer
	KafkaSession   *kafka.Session
	PulsarSession  *pulsar.Session
	SNCloudSession *config.Session
	logger         *logrus.Logger
}

// Close releases Cloud runtime generations and external protocol sessions.
func (s *Server) Close() {
	if s.SNCloudSession != nil {
		_ = s.SNCloudSession.Close()
	}
	if s.PulsarSession != nil {
		s.PulsarSession.ResetPulsarContext()
	}
	if s.KafkaSession != nil {
		s.KafkaSession.ResetKafkaContext()
	}
}

// NewServer creates a new MCP server with StreamNative integrations.
func NewServer(name, version string, logger *logrus.Logger, opts ...server.ServerOption) *Server {
	// Create a new MCP server
	opts = AddOpts(opts...)
	s := server.NewMCPServer(name, version, opts...)
	mcpserver := CreateSNCloudMCPServer(s, logger)
	return mcpserver
}

// AddOpts merges default server options with custom options.
func AddOpts(opts ...server.ServerOption) []server.ServerOption {
	defaultOpts := []server.ServerOption{
		server.WithResourceCapabilities(true, true),
		server.WithRecovery(),
		server.WithLogging(),
		server.WithToolHandlerMiddleware(func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
			return func(ctx context.Context, request mcpgo.CallToolRequest) (*mcpgo.CallToolResult, error) {
				if request.Params.Name == "sncloud_context_use_cluster" || request.Params.Name == "sncloud_context_reset" {
					return next(ctx, request)
				}
				ctx, release := context2.AcquireRuntimeContext(ctx)
				defer release()
				return next(ctx, request)
			}
		}),
		server.WithResourceHandlerMiddleware(func(next server.ResourceHandlerFunc) server.ResourceHandlerFunc {
			return func(ctx context.Context, request mcpgo.ReadResourceRequest) ([]mcpgo.ResourceContents, error) {
				ctx, release := context2.AcquireRuntimeContext(ctx)
				defer release()
				return next(ctx, request)
			}
		}),
		server.WithPromptHandlerMiddleware(func(next server.PromptHandlerFunc) server.PromptHandlerFunc {
			return func(ctx context.Context, request mcpgo.GetPromptRequest) (*mcpgo.GetPromptResult, error) {
				ctx, release := context2.AcquireRuntimeContext(ctx)
				defer release()
				return next(ctx, request)
			}
		}),
	}
	opts = append(defaultOpts, opts...)
	return opts
}

// CreateSNCloudMCPServer constructs a Server wrapper for StreamNative Cloud.
func CreateSNCloudMCPServer(s *server.MCPServer, logger *logrus.Logger) *Server {
	mcpserver := &Server{
		MCPServer:      s,
		logger:         logger,
		SNCloudSession: &config.Session{},
		KafkaSession:   &kafka.Session{},
		PulsarSession:  &pulsar.Session{},
	}

	return mcpserver
}
