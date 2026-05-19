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

package toolannotations

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func TestAnnotationHelpersSetGranularHints(t *testing.T) {
	tests := []struct {
		name        string
		option      mcp.ToolOption
		title       string
		readOnly    bool
		destructive bool
		idempotent  bool
		openWorld   bool
	}{
		{
			name:        "read only",
			option:      ReadOnly("Read Things"),
			title:       "Read Things",
			readOnly:    true,
			destructive: false,
			idempotent:  true,
			openWorld:   true,
		},
		{
			name:        "external read",
			option:      ExternalRead("Read External Things"),
			title:       "Read External Things",
			readOnly:    true,
			destructive: false,
			idempotent:  true,
			openWorld:   true,
		},
		{
			name:        "local session mutation",
			option:      LocalSessionMutation("Use Context"),
			title:       "Use Context",
			readOnly:    false,
			destructive: false,
			idempotent:  false,
			openWorld:   false,
		},
		{
			name:        "non destructive idempotent mutation",
			option:      Mutating("Update Setting", false, true),
			title:       "Update Setting",
			readOnly:    false,
			destructive: false,
			idempotent:  true,
			openWorld:   true,
		},
		{
			name:        "destructive compatibility",
			option:      Destructive("Delete Things"),
			title:       "Delete Things",
			readOnly:    false,
			destructive: true,
			idempotent:  false,
			openWorld:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := mcp.NewTool("test", tt.option)
			require.Equal(t, tt.title, tool.Annotations.Title)
			require.NotNil(t, tool.Annotations.ReadOnlyHint)
			require.NotNil(t, tool.Annotations.DestructiveHint)
			require.NotNil(t, tool.Annotations.IdempotentHint)
			require.NotNil(t, tool.Annotations.OpenWorldHint)
			require.Equal(t, tt.readOnly, *tool.Annotations.ReadOnlyHint)
			require.Equal(t, tt.destructive, *tool.Annotations.DestructiveHint)
			require.Equal(t, tt.idempotent, *tool.Annotations.IdempotentHint)
			require.Equal(t, tt.openWorld, *tool.Annotations.OpenWorldHint)
		})
	}
}
