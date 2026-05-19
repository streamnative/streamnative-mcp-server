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

package kafka

import (
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

type toolMode string

const (
	toolModeRead  toolMode = "read"
	toolModeWrite toolMode = "write"
)

func isToolModeWrite(mode toolMode) bool {
	return mode == toolModeWrite
}

func validateModeOperation(mode toolMode, operation string, writeOperations map[string]struct{}) bool {
	_, isWrite := writeOperations[strings.ToLower(operation)]
	return (mode == toolModeWrite) == isWrite
}

func pruneToolInputSchema(tool *mcp.Tool, allowedProperties []string) {
	allowed := make(map[string]struct{}, len(allowedProperties))
	for _, property := range allowedProperties {
		allowed[property] = struct{}{}
	}

	for property := range tool.InputSchema.Properties {
		if _, ok := allowed[property]; !ok {
			delete(tool.InputSchema.Properties, property)
		}
	}

	if len(tool.InputSchema.Required) == 0 {
		return
	}
	required := tool.InputSchema.Required[:0]
	for _, property := range tool.InputSchema.Required {
		if _, ok := allowed[property]; ok {
			required = append(required, property)
		}
	}
	tool.InputSchema.Required = required
}
