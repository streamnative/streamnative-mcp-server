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

package pftools

import (
	"encoding/json"
	"testing"

	legacy "github.com/mark3labs/mcp-go/mcp"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/require"
)

func TestLegacyToolToSDKUsesStructuredSchema(t *testing.T) {
	legacyTool := legacy.NewTool("test_tool", legacy.WithDescription("desc"))

	sdkTool := legacyToolToSDK(legacyTool)
	require.NotNil(t, sdkTool)
	require.Equal(t, "test_tool", sdkTool.Name)
	require.Equal(t, "desc", sdkTool.Description)

	schema, ok := sdkTool.InputSchema.(legacy.ToolInputSchema)
	require.True(t, ok)
	require.Equal(t, "object", schema.Type)
}

func TestSDKToolToLegacyHandlesRawSchema(t *testing.T) {
	raw := json.RawMessage(`{"type":"object","properties":{"foo":{"type":"string"}}}`)
	sdkTool := &sdk.Tool{
		Name:        "test_tool",
		Description: "desc",
		InputSchema: raw,
	}

	legacyTool := sdkToolToLegacy(sdkTool)
	require.Equal(t, "test_tool", legacyTool.Name)
	require.Equal(t, "desc", legacyTool.Description)
	require.Equal(t, raw, legacyTool.RawInputSchema)
	require.Equal(t, "", legacyTool.InputSchema.Type)
}

func TestLegacyResultToSDKTextContent(t *testing.T) {
	legacyResult := legacy.NewToolResultText("ok")

	sdkResult := legacyResultToSDK(legacyResult)
	require.NotNil(t, sdkResult)
	require.Len(t, sdkResult.Content, 1)

	textContent, ok := sdkResult.Content[0].(*sdk.TextContent)
	require.True(t, ok)
	require.Equal(t, "ok", textContent.Text)
}
