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
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sirupsen/logrus"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/adapter"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

type Server struct {
	Server         *mcpsdk.Server
	KafkaSession   *kafka.Session
	PulsarSession  *pulsar.Session
	SNCloudSession *config.Session
	logger         *logrus.Logger
}

// NewServer creates a new MCP server using go-sdk.
// The instructions parameter provides server description for clients.
func NewServer(name, version string, logger *logrus.Logger, instructions string) *Server {
	sdkServer := mcpsdk.NewServer(&mcpsdk.Implementation{
		Name:    name,
		Version: version,
	}, &mcpsdk.ServerOptions{
		Instructions: instructions,
	})

	mcpserver := &Server{
		Server:         sdkServer,
		logger:         logger,
		SNCloudSession: &config.Session{},
		KafkaSession:   &kafka.Session{},
		PulsarSession:  &pulsar.Session{},
	}

	return mcpserver
}

// AddTool adds a tool to the server with backward compatibility.
// This method accepts both mark3labs-style and go-sdk-style handlers.
func (s *Server) AddTool(tool interface{}, handler interface{}) error {
	// Adapt handler to go-sdk signature
	adaptedHandler := adapter.AdaptHandlerV1ToV2(handler)

	// Convert tool to *mcpsdk.Tool if it's not already
	var mcpTool *mcpsdk.Tool
	switch t := tool.(type) {
	case *mcpsdk.Tool:
		mcpTool = t
	default:
		// For any other type, try to use it as-is (should be *mcpsdk.Tool from builders)
		mcpTool = t.(*mcpsdk.Tool)
	}

	s.Server.AddTool(mcpTool, adaptedHandler)
	return nil
}

// AddSessionTool adds a session-scoped tool (used by PFTools).
func (s *Server) AddSessionTool(sessionID string, tool interface{}, handler interface{}) error {
	// For go-sdk, session-scoped tools are handled differently
	// The actual implementation will be in Phase 2 when we migrate PFTools
	adaptedHandler := adapter.AdaptHandlerV1ToV2(handler)

	// Convert tool to *mcpsdk.Tool if it's not already
	var mcpTool *mcpsdk.Tool
	switch t := tool.(type) {
	case *mcpsdk.Tool:
		mcpTool = t
	default:
		// For any other type, try to use it as-is (should be *mcpsdk.Tool from builders)
		mcpTool = t.(*mcpsdk.Tool)
	}

	s.Server.AddTool(mcpTool, adaptedHandler)
	return nil
}

// DeleteTools removes tools by name prefix.
func (s *Server) DeleteTools(name string) error {
	// go-sdk doesn't have a direct DeleteTools method
	// This will be implemented in Phase 2 when we migrate PFTools
	return nil
}

// DeleteSessionTools removes session-scoped tools.
func (s *Server) DeleteSessionTools(sessionID, name string) error {
	// go-sdk doesn't have a direct DeleteSessionTools method
	// This will be implemented in Phase 2 when we migrate PFTools
	return nil
}
