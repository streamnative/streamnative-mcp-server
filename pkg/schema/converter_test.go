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

package schema

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Future test functions will be added here.

func TestConverterFactory(t *testing.T) {
	tests := []struct {
		name       string
		schemaType string
		wantType   reflect.Type
		wantErr    bool
	}{
		{name: "AVRO", schemaType: "AVRO", wantType: reflect.TypeOf(&AvroConverter{}), wantErr: false},
		{name: "JSON", schemaType: "JSON", wantType: reflect.TypeOf(&JSONConverter{}), wantErr: false},
		{name: "STRING", schemaType: "STRING", wantType: reflect.TypeOf(&StringConverter{}), wantErr: false},
		{name: "BYTES", schemaType: "BYTES", wantType: reflect.TypeOf(&StringConverter{}), wantErr: false},
		{name: "INT8", schemaType: "INT8", wantType: reflect.TypeOf(&NumberConverter{}), wantErr: false},
		{name: "INT16", schemaType: "INT16", wantType: reflect.TypeOf(&NumberConverter{}), wantErr: false},
		{name: "INT32", schemaType: "INT32", wantType: reflect.TypeOf(&NumberConverter{}), wantErr: false},
		{name: "INT64", schemaType: "INT64", wantType: reflect.TypeOf(&NumberConverter{}), wantErr: false},
		{name: "FLOAT", schemaType: "FLOAT", wantType: reflect.TypeOf(&NumberConverter{}), wantErr: false},
		{name: "DOUBLE", schemaType: "DOUBLE", wantType: reflect.TypeOf(&NumberConverter{}), wantErr: false},
		{name: "BOOLEAN", schemaType: "BOOLEAN", wantType: reflect.TypeOf(&BooleanConverter{}), wantErr: false},
		{name: "avro_lowercase", schemaType: "avro", wantType: nil, wantErr: true},
		{name: "UNKNOWN_TYPE", schemaType: "UNKNOWN_TYPE", wantType: nil, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConverterFactory(tt.schemaType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConverterFactory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && reflect.TypeOf(got) != tt.wantType {
				t.Errorf("ConverterFactory() got = %v, want %v", reflect.TypeOf(got), tt.wantType)
			}
			// For error cases, we might also want to assert the error message if it's specific.
			if tt.wantErr && err != nil {
				assert.Contains(t, err.Error(), "unsupported schema type")
			}
		})
	}
}
