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

package schema

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// Helper function to create SchemaInfo for boolean tests
func newBoolSchemaInfo(schemaType string) *utils.SchemaInfo {
	return &utils.SchemaInfo{
		Type:   schemaType,
		Schema: []byte{},
	}
}

func TestNewBooleanConverter(t *testing.T) {
	converter := NewBooleanConverter()
	assert.NotNil(t, converter)
	assert.Equal(t, ParamName, converter.ParamName, "ParamName should be initialized to the package constant")
}

func TestBooleanConverter_ToToolInputSchema(t *testing.T) {
	converter := NewBooleanConverter()

	tests := []struct {
		name       string
		schemaInfo *utils.SchemaInfo
		wantSchema string
		wantErr    bool
	}{
		{
			name:       "Valid BOOLEAN schema",
			schemaInfo: newBoolSchemaInfo("BOOLEAN"),
			wantSchema: mustSchemaJSON(newPayloadSchema(ParamName, fmt.Sprintf("The input schema is a %s schema", "BOOLEAN"), "boolean")),
			wantErr:    false,
		},
		{
			name:       "Invalid schema type (STRING)",
			schemaInfo: newBoolSchemaInfo("STRING"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSchema, err := converter.ToToolInputSchema(tt.schemaInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToToolInputSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.Nil(t, gotSchema)
				return
			}
			assert.NotNil(t, gotSchema)
			assert.JSONEq(t, tt.wantSchema, mustSchemaJSON(gotSchema))
		})
	}
}

func TestBooleanConverter_ValidateArguments(t *testing.T) {
	converter := NewBooleanConverter()

	tests := []struct {
		name       string
		schemaInfo *utils.SchemaInfo
		args       map[string]any
		wantErr    bool
		errContain string // Substring to check in error message if wantErr is true
	}{
		{
			name:       "Valid arguments for BOOLEAN schema",
			schemaInfo: newBoolSchemaInfo("BOOLEAN"),
			args:       map[string]any{ParamName: true},
			wantErr:    false,
		},
		{
			name:       "Invalid schema type (STRING)",
			schemaInfo: newBoolSchemaInfo("STRING"),
			args:       map[string]any{ParamName: true},
			wantErr:    true,
			errContain: "expected BOOLEAN schema, got STRING",
		},
		{
			name:       "Missing payload argument",
			schemaInfo: newBoolSchemaInfo("BOOLEAN"),
			args:       map[string]any{},
			wantErr:    true,
			errContain: "missing required parameter: payload",
		},
		{
			name:       "Incorrect payload type (string instead of bool)",
			schemaInfo: newBoolSchemaInfo("BOOLEAN"),
			args:       map[string]any{ParamName: "true"},
			wantErr:    true,
			errContain: "parameter payload is not of type bool",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := converter.ValidateArguments(tt.args, tt.schemaInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateArguments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				assert.Contains(t, err.Error(), tt.errContain)
			}
		})
	}
}

func TestBooleanConverter_SerializeMCPRequestToPulsarPayload(t *testing.T) {
	converter := NewBooleanConverter()

	tests := []struct {
		name       string
		schemaInfo *utils.SchemaInfo
		args       map[string]any
		want       []byte
		wantErr    bool
		errContain string
	}{
		{
			name:       "Serialize true for BOOLEAN schema",
			schemaInfo: newBoolSchemaInfo("BOOLEAN"),
			args:       map[string]any{ParamName: true},
			want:       []byte("true"),
			wantErr:    false,
		},
		{
			name:       "Serialize false for BOOLEAN schema",
			schemaInfo: newBoolSchemaInfo("BOOLEAN"),
			args:       map[string]any{ParamName: false},
			want:       []byte("false"),
			wantErr:    false,
		},
		{
			name:       "Validation error (e.g., missing payload)",
			schemaInfo: newBoolSchemaInfo("BOOLEAN"),
			args:       map[string]any{},
			want:       nil,
			wantErr:    true,
			errContain: "arguments validation failed",
		},
		{
			name:       "Validation error (e.g., wrong schema type during ValidateArguments)",
			schemaInfo: newBoolSchemaInfo("STRING"),
			args:       map[string]any{ParamName: true},
			want:       nil,
			wantErr:    true,
			errContain: "arguments validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.SerializeMCPRequestToPulsarPayload(tt.args, tt.schemaInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeMCPRequestToPulsarPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
			if tt.wantErr && err != nil {
				assert.Contains(t, err.Error(), tt.errContain)
			}
		})
	}
}

func mustSchemaJSON(schema any) string {
	data, err := json.Marshal(schema)
	if err != nil {
		return ""
	}
	return string(data)
}
