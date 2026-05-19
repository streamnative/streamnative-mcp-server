// Copyright 2026 StreamNative
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

// Package toolannotations provides helpers for setting MCP tool safety annotations.
package toolannotations

import "github.com/mark3labs/mcp-go/mcp"

// ReadOnly annotates a non-mutating read tool. It is safe to call repeatedly and
// may interact with external services to retrieve data.
func ReadOnly(title string) mcp.ToolOption {
	return compose(
		mcp.WithTitleAnnotation(title),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(true),
	)
}

// ExternalRead is an explicit alias for read tools whose purpose is to retrieve
// data from external systems.
func ExternalRead(title string) mcp.ToolOption {
	return ReadOnly(title)
}

// LocalSessionMutation annotates a tool that mutates only MCP server/session
// state, not external StreamNative, Kafka, or Pulsar resources.
func LocalSessionMutation(title string) mcp.ToolOption {
	return compose(
		mcp.WithTitleAnnotation(title),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	)
}

// Mutating annotates a tool that may perform external side effects. Callers must
// classify whether the operation may be destructive and whether repeated calls
// with identical arguments are idempotent.
func Mutating(title string, destructive bool, idempotent bool) mcp.ToolOption {
	return compose(
		mcp.WithTitleAnnotation(title),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(destructive),
		mcp.WithIdempotentHintAnnotation(idempotent),
		mcp.WithOpenWorldHintAnnotation(true),
	)
}

// Destructive annotates a non-idempotent tool that may modify external resources.
func Destructive(title string) mcp.ToolOption {
	return Mutating(title, true, false)
}

func compose(opts ...mcp.ToolOption) mcp.ToolOption {
	return func(t *mcp.Tool) {
		for _, opt := range opts {
			opt(t)
		}
	}
}
