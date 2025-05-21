package schema

import (
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"
)

// GetSchemaType 返回Schema类型的字符串表示
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
