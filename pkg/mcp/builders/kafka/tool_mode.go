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
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
)

type toolMode = builders.OperationMode

const (
	toolModeRead  = builders.OperationModeRead
	toolModeWrite = builders.OperationModeWrite
)

func isToolModeWrite(mode toolMode) bool {
	return mode == toolModeWrite
}

func validateModeOperation(mode toolMode, operation string, operations builders.OperationRegistry) error {
	return operations.ValidateModeOperation(mode, operation)
}

func pruneToolInputSchema(tool *mcp.Tool, allowedProperties []string) {
	builders.PruneToolInputSchema(tool, allowedProperties)
}
