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

package pulsar

import (
	legacy "github.com/mark3labs/mcp-go/mcp"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func legacyToolResultFromSDK(result *sdk.CallToolResult) *legacy.CallToolResult {
	if result == nil {
		return legacy.NewToolResultText("")
	}

	text := ""
	for _, content := range result.Content {
		if textContent, ok := content.(*sdk.TextContent); ok {
			text = textContent.Text
			break
		}
	}

	if result.IsError {
		if text == "" {
			text = "tool call failed"
		}
		return legacy.NewToolResultError(text)
	}

	return legacy.NewToolResultText(text)
}
