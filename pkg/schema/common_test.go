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
	"testing"

	"github.com/apache/pulsar-client-go/pulsar"
)

// Future test functions will be added here.

func TestGetSchemaType(t *testing.T) {
	tests := []struct {
		name       string
		schemaType pulsar.SchemaType
		want       string
	}{
		{name: "AVRO", schemaType: pulsar.AVRO, want: "AVRO"},
		{name: "JSON", schemaType: pulsar.JSON, want: "JSON"},
		{name: "STRING", schemaType: pulsar.STRING, want: "STRING"},
		{name: "INT8", schemaType: pulsar.INT8, want: "INT8"},
		{name: "INT16", schemaType: pulsar.INT16, want: "INT16"},
		{name: "INT32", schemaType: pulsar.INT32, want: "INT32"},
		{name: "INT64", schemaType: pulsar.INT64, want: "INT64"},
		{name: "FLOAT", schemaType: pulsar.FLOAT, want: "FLOAT"},
		{name: "DOUBLE", schemaType: pulsar.DOUBLE, want: "DOUBLE"},
		{name: "BOOLEAN", schemaType: pulsar.BOOLEAN, want: "BOOLEAN"},
		{name: "BYTES", schemaType: pulsar.BYTES, want: "BYTES"},
		{name: "Unknown", schemaType: pulsar.SchemaType(999), want: "Unknown(999)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSchemaType(tt.schemaType); got != tt.want {
				t.Errorf("GetSchemaType() = %v, want %v", got, tt.want)
			}
		})
	}
}
