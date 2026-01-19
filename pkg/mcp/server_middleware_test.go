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
	"testing"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/require"
)

func TestRecoveryMiddleware_ToolPanicUsesToolName(t *testing.T) {
	middleware := recoveryMiddleware(nil)
	handler := middleware(func(context.Context, string, sdk.Request) (sdk.Result, error) {
		panic("boom")
	})

	req := &sdk.CallToolRequest{
		Params: &sdk.CallToolParamsRaw{
			Name: "kafka.topics.create",
		},
	}
	result, err := handler(context.Background(), "tools/call", req)
	require.Nil(t, result)
	require.EqualError(t, err, "panic recovered in kafka.topics.create tool handler: boom")
}

func TestRecoveryMiddleware_MethodPanicUsesMethodName(t *testing.T) {
	middleware := recoveryMiddleware(nil)
	handler := middleware(func(context.Context, string, sdk.Request) (sdk.Result, error) {
		panic("boom")
	})

	req := &sdk.ListResourcesRequest{}
	result, err := handler(context.Background(), "resources/list", req)
	require.Nil(t, result)
	require.EqualError(t, err, "panic recovered in resources/list request: boom")
}
