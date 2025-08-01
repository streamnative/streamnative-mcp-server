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

	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	pulsarbuilders "github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders/pulsar"
)

// PulsarAdminAddFunctionsTools adds a unified function-related tool to the MCP server
func PulsarAdminAddFunctionsTools(s *server.MCPServer, readOnly bool, features []string) {
	// Use the new builder pattern
	builder := pulsarbuilders.NewPulsarAdminFunctionsToolBuilder()
	config := builders.ToolBuildConfig{
		ReadOnly: readOnly,
		Features: features,
	}

	tools, err := builder.BuildTools(context.Background(), config)
	if err != nil {
		// Log error but don't fail - this maintains backward compatibility
		return
	}

	// Add all built tools to the server
	for _, tool := range tools {
		s.AddTool(tool.Tool, tool.Handler)
	}
}