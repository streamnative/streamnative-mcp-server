// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
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
	"strings"

	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	pulsarbuilders "github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders/pulsar"
)

const (
	HTTP    = "http"
	FILE    = "file"
	BUILTIN = "builtin"

	FUNCTION = "function"
	SINK     = "sink"
	SOURCE   = "source"

	PublicTenant     = "public"
	DefaultNamespace = "default"
)

func IsPackageURLSupported(functionPkgURL string) bool {
	return functionPkgURL != "" && (strings.HasPrefix(functionPkgURL, HTTP) ||
		strings.HasPrefix(functionPkgURL, FILE) ||
		strings.HasPrefix(functionPkgURL, FUNCTION) ||
		strings.HasPrefix(functionPkgURL, SINK) ||
		strings.HasPrefix(functionPkgURL, SOURCE))
}

// PulsarAdminAddPackagesTools adds package-related tools to the MCP server
func PulsarAdminAddPackagesTools(s *server.MCPServer, readOnly bool, features []string) {
	// Use the new builder pattern
	builder := pulsarbuilders.NewPulsarAdminPackagesToolBuilder()
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
