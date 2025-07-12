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

package mcp

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

type Server struct {
	MCPServer      *server.MCPServer
	KafkaSession   *kafka.Session
	PulsarSession  *pulsar.Session
	SNCloudSession *config.Session
	logger         *logrus.Logger
}

func NewServer(name, version string, logger *logrus.Logger, opts ...server.ServerOption) *Server {
	// Create a new MCP server
	opts = addOpts(opts...)
	s := server.NewMCPServer(name, version, opts...)
	mcpserver := createSNCloudMCPServer(s, logger)
	return mcpserver
}

func addOpts(opts ...server.ServerOption) []server.ServerOption {
	defaultOpts := []server.ServerOption{
		server.WithResourceCapabilities(true, true),
		server.WithRecovery(),
		server.WithLogging(),
	}
	opts = append(defaultOpts, opts...)
	return opts
}

func createSNCloudMCPServer(s *server.MCPServer, logger *logrus.Logger) *Server {
	mcpserver := &Server{
		MCPServer: s,
		logger:    logger,
	}

	return mcpserver
}
