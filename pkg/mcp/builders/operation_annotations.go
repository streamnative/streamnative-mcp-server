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

package builders

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/toolannotations"
)

// ToolAnnotationForMode selects Claude connector safety annotations from the
// operation mode so schemas, validation, and annotations share the same mode
// vocabulary. If operation metadata is supplied, write-tool destructive and
// idempotent hints are derived from the write specs instead of from mode alone.
func ToolAnnotationForMode(mode OperationMode, readTitle, writeTitle string, registries ...OperationRegistry) mcp.ToolOption {
	if mode != OperationModeWrite {
		return toolannotations.ExternalRead(readTitle)
	}

	destructive, idempotent := writeAnnotationHints(registries...)
	return toolannotations.Mutating(writeTitle, destructive, idempotent)
}

func writeAnnotationHints(registries ...OperationRegistry) (destructive bool, idempotent bool) {
	if len(registries) == 0 {
		return true, false
	}

	foundWriteSpec := false
	allIdempotent := true
	for _, registry := range registries {
		for _, spec := range registry.SpecsForMode(OperationModeWrite) {
			foundWriteSpec = true
			destructive = destructive || spec.Destructive
			allIdempotent = allIdempotent && spec.Idempotent
		}
	}
	if !foundWriteSpec {
		return true, false
	}

	// A multi-operation tool is idempotent only when every exposed write
	// operation is idempotent.
	return destructive, allIdempotent
}
