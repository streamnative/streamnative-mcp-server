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

	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	pulsarBuilders "github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders/pulsar"
)

// PulsarAdminSchemaTools creates Pulsar Admin Schema tool list using the new builder pattern
func PulsarAdminSchemaTools(readOnly bool, features []string) []builders.ServerTool {
	builder := pulsarBuilders.NewPulsarAdminSchemaToolBuilder()
	config := builders.ToolBuildConfig{
		ReadOnly: readOnly,
		Features: features,
	}

	tools, err := builder.BuildTools(context.Background(), config)
	if err != nil {
		// In production environment, this should use proper logging
		fmt.Printf("Failed to build Pulsar Admin Schema tools: %v\n", err)
		return nil
	}

	return tools
}

// PulsarAdminAddSchemasTools adds schema-related tools to the MCP server
func PulsarAdminAddSchemasTools(s *MCPServer, readOnly bool, features []string) {
	tools := PulsarAdminSchemaTools(readOnly, features)

	for _, tool := range tools {
		s.AddTool(tool.Tool, tool.Handler)
	}
}
