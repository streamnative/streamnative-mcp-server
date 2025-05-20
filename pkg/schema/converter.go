package schema

import (
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/mark3labs/mcp-go/mcp"
)

type SchemaConverter interface {
	// ToMCPToolInputSchemaProperties 将给定的Pulsar SchemaInfo转换为
	// 一系列mcp.ToolOption，用于构建mcp.Tool的输入模式
	ToMCPToolInputSchemaProperties(pulsarSchemaInfo *pulsar.SchemaInfo) ([]mcp.ToolOption, error)

	// SerializeMCPRequestToPulsarPayload 将来自MCP工具调用的参数(map[string]any)
	// 根据目标Pulsar SchemaInfo序列化为字节切片，作为Pulsar消息的负载
	SerializeMCPRequestToPulsarPayload(arguments map[string]any, targetPulsarSchemaInfo *pulsar.SchemaInfo) ([]byte, error)

	// ValidateArguments 验证参数是否符合目标Pulsar模式要求
	ValidateArguments(arguments map[string]any, targetPulsarSchemaInfo *pulsar.SchemaInfo) error
}

func ConverterFactory(schemaType pulsar.SchemaType) (SchemaConverter, error) {
	switch schemaType {
	case pulsar.AVRO:
		return NewAvroConverter(), nil
	case pulsar.JSON:
		return NewJSONConverter(), nil
	case pulsar.STRING:
		return NewStringConverter(), nil
	case pulsar.INT8, pulsar.INT16, pulsar.INT32, pulsar.INT64, pulsar.FLOAT, pulsar.DOUBLE:
		return NewNumberConverter(), nil
	case pulsar.BOOLEAN:
		return NewBooleanConverter(), nil
	case pulsar.BYTES:
		return NewBytesConverter(), nil
	default:
		return nil, fmt.Errorf("unsupported schema type: %v", schemaType)
	}
}

type BaseConverter struct{}
