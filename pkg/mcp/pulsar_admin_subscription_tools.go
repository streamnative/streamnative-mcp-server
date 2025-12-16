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
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	pulsarBuilders "github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders/pulsar"
)

// PulsarAdminSubscriptionTools creates Pulsar Admin Subscription tool list using the new builder pattern
func PulsarAdminSubscriptionTools(readOnly bool, features []string) []server.ServerTool {
	builder := pulsarBuilders.NewPulsarAdminSubscriptionToolBuilder()
	config := builders.ToolBuildConfig{
		ReadOnly: readOnly,
		Features: features,
	}

	tools, err := builder.BuildTools(context.Background(), config)
	if err != nil {
		// In production environment, this should use proper logging
		fmt.Printf("Failed to build Pulsar Admin Subscription tools: %v\n", err)
		return nil
	}

	return tools
}

// PulsarAdminAddSubscriptionTools adds subscription-related tools to the MCP server
func PulsarAdminAddSubscriptionTools(s *server.MCPServer, readOnly bool, features []string) {
	tools := PulsarAdminSubscriptionTools(readOnly, features)

	for _, tool := range tools {
		s.AddTool(tool.Tool, tool.Handler)
	}
}
