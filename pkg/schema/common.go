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
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"
)

// GetSchemaType returns the string representation of a schema type.
func GetSchemaType(schemaType pulsar.SchemaType) string {
	switch schemaType {
	case pulsar.AVRO:
		return "AVRO"
	case pulsar.JSON:
		return "JSON"
	case pulsar.STRING:
		return "STRING"
	case pulsar.INT8:
		return "INT8"
	case pulsar.INT16:
		return "INT16"
	case pulsar.INT32:
		return "INT32"
	case pulsar.INT64:
		return "INT64"
	case pulsar.FLOAT:
		return "FLOAT"
	case pulsar.DOUBLE:
		return "DOUBLE"
	case pulsar.BOOLEAN:
		return "BOOLEAN"
	case pulsar.BYTES:
		return "BYTES"
	default:
		return fmt.Sprintf("Unknown(%d)", schemaType)
	}
}
