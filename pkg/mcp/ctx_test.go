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
)

func TestMCPRequestContextHelpers(t *testing.T) {
	base := context.Background()
	if GetMCPRequest(base) != nil {
		t.Fatal("expected nil request")
	}
	if GetMCPRequestExtra(base) != nil {
		t.Fatal("expected nil request extra")
	}
	if GetMCPSession(base) != nil {
		t.Fatal("expected nil session")
	}
	if GetMCPSessionID(base) != "" {
		t.Fatal("expected empty session ID")
	}

	extra := &sdk.RequestExtra{}
	reqStub := &sdk.ServerRequest[*sdk.CallToolParamsRaw]{Extra: extra}
	ctx := WithMCPRequest(base, reqStub)

	req := GetMCPRequest(ctx)
	if req == nil {
		t.Fatal("expected request")
	}
	if got := GetMCPRequestExtra(ctx); got != extra {
		t.Fatal("expected request extra to match")
	}
	if GetMCPSession(ctx) != nil {
		t.Fatal("expected nil session")
	}
	if GetMCPSessionID(ctx) != "" {
		t.Fatal("expected empty session ID")
	}
	if stub, ok := req.(*sdk.ServerRequest[*sdk.CallToolParamsRaw]); !ok || stub.Extra != extra {
		t.Fatal("expected stored request")
	}
}
