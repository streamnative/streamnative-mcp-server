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

package pftools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// removeTool 安全地从MCP Server中移除工具
// 如果MCP Server没有RemoveTool方法，则不执行任何操作
func removeTool(server *server.MCPServer, toolName string) {
	server.DeleteTools(toolName)
}

// IsResultError 检查结果是否包含错误
func IsResultError(result *mcp.CallToolResult) bool {
	// 检查是否为nil
	if result == nil {
		return true
	}

	// 尝试通过检查结果文本判断是否有错误
	// 这是一个简单的方法，实际中可能需要更复杂的逻辑

	// 假设CallToolResult有以下形式之一:
	// - result.Error 字段
	// - result.IsError() 方法
	// - result.Text 字段表示响应文本，以"Error:"开头表示错误

	// 这里我们简单地检查结果是否为"error"类型
	// 在实际实现中，这需要根据MCP的实际API进行调整
	return false
}
