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
	"fmt"
	"sort"
	"strings"
)

// OperationMode classifies whether an MCP operation is safe to expose through a
// read-only tool or must be exposed through a write/destructive tool.
type OperationMode string

const (
	// OperationModeRead marks an operation as safe for read-only tools.
	OperationModeRead OperationMode = "read"
	// OperationModeWrite marks an operation as requiring a write/destructive tool.
	OperationModeWrite OperationMode = "write"
)

// ParamSpec describes an operation parameter. Builders can use this metadata for
// generated documentation and schema checks without coupling it to mcp-go types.
type ParamSpec struct {
	Name        string
	Description string
	Required    bool
}

// OperationHandler is a package-neutral handler placeholder for builders that
// later choose to dispatch directly from operation metadata.
type OperationHandler any

// OperationSpec is the single source of truth for an operation's safety mode and
// documentation metadata.
type OperationSpec struct {
	Name        string
	Mode        OperationMode
	Description string
	Destructive bool
	Idempotent  bool
	Resources   []string
	Params      []ParamSpec
	Handler     OperationHandler
}

// OperationRegistry stores a tool family's operation specs in schema order.
type OperationRegistry []OperationSpec

// Names returns operation names in registry order.
func (r OperationRegistry) Names() []string {
	names := make([]string, 0, len(r))
	for _, spec := range r {
		names = append(names, spec.Name)
	}
	return names
}

// NamesForMode returns operation names for one mode in registry order.
func (r OperationRegistry) NamesForMode(mode OperationMode) []string {
	names := make([]string, 0, len(r))
	for _, spec := range r {
		if spec.Mode == mode {
			names = append(names, spec.Name)
		}
	}
	return names
}

// ReadNames returns read operation names in registry order.
func (r OperationRegistry) ReadNames() []string {
	return r.NamesForMode(OperationModeRead)
}

// WriteNames returns write operation names in registry order.
func (r OperationRegistry) WriteNames() []string {
	return r.NamesForMode(OperationModeWrite)
}

// WriteSet returns write operation names as a lookup set.
func (r OperationRegistry) WriteSet() map[string]struct{} {
	set := make(map[string]struct{})
	for _, spec := range r {
		if spec.Mode == OperationModeWrite {
			set[strings.ToLower(spec.Name)] = struct{}{}
		}
	}
	return set
}

// SpecFor returns the metadata for operation. Matching is case-insensitive.
func (r OperationRegistry) SpecFor(operation string) (OperationSpec, bool) {
	operation = strings.ToLower(operation)
	for _, spec := range r {
		if strings.ToLower(spec.Name) == operation {
			return spec, true
		}
	}
	return OperationSpec{}, false
}

// ValidateModeOperation rejects unknown operations and operations exposed through
// the wrong read/write tool mode.
func (r OperationRegistry) ValidateModeOperation(mode OperationMode, operation string) error {
	spec, ok := r.SpecFor(operation)
	if !ok {
		return fmt.Errorf("unknown operation %q. Supported operations: %s", operation, strings.Join(r.NamesForMode(mode), ", "))
	}
	if spec.Mode != mode {
		return fmt.Errorf("operation %q is not available in %s mode", operation, mode)
	}
	return nil
}

// DescriptionsForMode returns markdown bullet rows for operations with
// descriptions. Specs without descriptions are rendered as just the operation name.
func (r OperationRegistry) DescriptionsForMode(mode OperationMode) []string {
	rows := make([]string, 0, len(r))
	for _, spec := range r {
		if spec.Mode != mode {
			continue
		}
		if spec.Description == "" {
			rows = append(rows, "- "+spec.Name)
			continue
		}
		rows = append(rows, fmt.Sprintf("- %s: %s", spec.Name, spec.Description))
	}
	return rows
}

// MustValidate checks that the registry has unique names and valid modes. It is
// intended for tests and init-time assertions in generated registries.
func (r OperationRegistry) MustValidate() {
	seen := make(map[string]struct{}, len(r))
	for _, spec := range r {
		name := strings.ToLower(spec.Name)
		if name == "" {
			panic("operation spec has empty name")
		}
		if _, ok := seen[name]; ok {
			panic(fmt.Sprintf("duplicate operation spec %q", spec.Name))
		}
		seen[name] = struct{}{}
		switch spec.Mode {
		case OperationModeRead, OperationModeWrite:
		default:
			panic(fmt.Sprintf("operation %q has invalid mode %q", spec.Name, spec.Mode))
		}
	}
}

// SortedNames returns operation names sorted lexicographically for stable tests.
func (r OperationRegistry) SortedNames() []string {
	names := r.Names()
	sort.Strings(names)
	return names
}
