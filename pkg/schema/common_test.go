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
