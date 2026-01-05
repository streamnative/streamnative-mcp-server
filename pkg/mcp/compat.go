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
)

// MCPServer is a type alias for go-sdk's Server.
// This provides backward compatibility with existing tool registration code.
type MCPServer = mcpsdk.Server

// AddTool is a compatibility wrapper for adding tools to the server.
// This matches the signature used by mark3labs/mcp-go server.
func (s *Server) AddToolCompat(tool interface{}, handler interface{}) {
	s.AddTool(tool, handler)
}
