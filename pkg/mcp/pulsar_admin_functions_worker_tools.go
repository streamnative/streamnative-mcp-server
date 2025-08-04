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
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	pulsarBuilders "github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders/pulsar"
)

// PulsarAdminFunctionsWorkerTools creates Pulsar Admin Functions Worker tool list using the new builder pattern
func PulsarAdminFunctionsWorkerTools(readOnly bool, features []string) []server.ServerTool {
	builder := pulsarBuilders.NewPulsarAdminFunctionsWorkerToolBuilder()
	config := builders.ToolBuildConfig{
		ReadOnly: readOnly,
		Features: features,
	}

	tools, err := builder.BuildTools(context.Background(), config)
	if err != nil {
		// In production environment, this should use proper logging
		fmt.Printf("Failed to build Pulsar Admin Functions Worker tools: %v\n", err)
		return nil
	}

	return tools
}

// PulsarAdminAddFunctionsWorkerTools adds functions worker-related tools to the MCP server
func PulsarAdminAddFunctionsWorkerTools(s *server.MCPServer, readOnly bool, features []string) {
	tools := PulsarAdminFunctionsWorkerTools(readOnly, features)

	for _, tool := range tools {
		s.AddTool(tool.Tool, tool.Handler)
	}
}

// handleFunctionsWorkerTool returns a function to handle functions worker tool requests
