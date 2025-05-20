package schema

import (
	"fmt"
	"reflect"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/hamba/avro/v2"
	"github.com/mark3labs/mcp-go/mcp"
)

type AvroConverter struct {
	BaseConverter
}

func NewAvroConverter() *AvroConverter {
	return &AvroConverter{}
}

func (c *AvroConverter) ToMCPToolInputSchemaProperties(schemaInfo *pulsar.SchemaInfo) ([]mcp.ToolOption, error) {
	if schemaInfo.Type != pulsar.AVRO {
		return nil, fmt.Errorf("expected AVRO schema, got %s", GetSchemaType(schemaInfo.Type))
	}

	// 解析AVRO模式
	schema, err := avro.Parse(schemaInfo.Schema)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AVRO schema: %w", err)
	}

	// 检查是否是记录类型
	recordSchema, ok := schema.(*avro.RecordSchema)
	if !ok {
		return nil, fmt.Errorf("expected AVRO record schema, got %s", reflect.TypeOf(schema).String())
	}

	var opts []mcp.ToolOption

	// 遍历AVRO记录的字段
	for _, field := range recordSchema.Fields() {
		fieldOption, err := c.avroFieldToPropertyOption(field)
		if err != nil {
			return nil, fmt.Errorf("failed to convert field '%s': %w", field.Name(), err)
		}
		opts = append(opts, fieldOption)
	}

	return opts, nil
}

// avroFieldToPropertyOption 将AVRO字段转换为MCP属性选项
func (c *AvroConverter) avroFieldToPropertyOption(field *avro.Field) (mcp.ToolOption, error) {
	fieldType := field.Type()
	fieldName := field.Name()

	// 获取字段的文档作为描述
	var description string
	if field.Doc() != "" {
		description = field.Doc()
	} else {
		description = fmt.Sprintf("%s field", fieldName)
	}

	// 检查字段是否必需
	// 在AVRO中，如果字段类型是与null的联合类型，则字段是可选的
	isRequired := true
	if unionType, ok := fieldType.(*avro.UnionSchema); ok {
		for _, t := range unionType.Types() {
			if t.Type() == "null" {
				isRequired = false
				break
			}
		}
	}

	// 根据AVRO类型创建相应的MCP属性选项
	var prop []mcp.PropertyOption
	var opt mcp.ToolOption

	if isRequired {
		prop = append(prop, mcp.Required())
	}
	prop = append(prop, mcp.Description(description))

	switch fieldType.Type() {
	case "string":
		if field.HasDefault() {
			defaultValue := field.Default()
			prop = append(prop, mcp.DefaultString(defaultValue.(string)))
		}
		opt = mcp.WithString(fieldName, prop...)
	case "int", "long":
		if field.HasDefault() {
			defaultValue := field.Default()
			prop = append(prop, mcp.DefaultNumber(defaultValue.(float64)))
		}
		opt = mcp.WithNumber(fieldName, prop...)
	case "float", "double":
		if field.HasDefault() {
			defaultValue := field.Default()
			prop = append(prop, mcp.DefaultNumber(defaultValue.(float64)))
		}
		opt = mcp.WithNumber(fieldName, prop...)
	case "boolean":
		if field.HasDefault() {
			defaultValue := field.Default()
			prop = append(prop, mcp.DefaultBool(defaultValue.(bool)))
		}
		opt = mcp.WithBoolean(fieldName, prop...)
	case "bytes":
		if field.HasDefault() {
			defaultValue := field.Default()
			prop = append(prop, mcp.DefaultString(string(defaultValue.([]byte))))
		}
		opt = mcp.WithString(fieldName, prop...)
	case "array":
		// 处理数组类型
		arrayType, ok := fieldType.(*avro.ArraySchema)
		if !ok {
			return nil, fmt.Errorf("expected array schema for field '%s'", fieldName)
		}

		itemsProp, err := c.avroTypeToPropertyOption("item", arrayType.Items(), description)
		if err != nil {
			return nil, fmt.Errorf("failed to convert array items for field '%s': %w", fieldName, err)
		}

		prop = append(prop, mcp.Items(itemsProp))

		opt = mcp.WithArray(fieldName, prop...)
	case "map":
		// 处理映射类型
		mapType, ok := fieldType.(*avro.MapSchema)
		if !ok {
			return nil, fmt.Errorf("expected map schema for field '%s'", fieldName)
		}

		valuesProp, err := c.avroTypeToPropertyOption("value", mapType.Values(), description)
		if err != nil {
			return nil, fmt.Errorf("failed to convert map values for field '%s': %w", fieldName, err)
		}

		prop = append(prop, mcp.Properties(valuesProp))

		opt = mcp.WithObject(fieldName, prop...)
	case "record":
		// 处理嵌套记录类型
		recordType, ok := fieldType.(*avro.RecordSchema)
		if !ok {
			return nil, fmt.Errorf("expected record schema for field '%s'", fieldName)
		}

		fieldProps := make(map[string]any)
		for _, subField := range recordType.Fields() {
			fieldProp, err := c.avroTypeToPropertyOption(subField.Name(), subField.Type(), subField.Doc())
			if err != nil {
				return nil, fmt.Errorf("failed to convert subfield '%s.%s': %w", fieldName, subField.Name(), err)
			}

			fieldProps[subField.Name()] = fieldProp
		}
		prop = append(prop, mcp.Properties(fieldProps))

		opt = mcp.WithObject(fieldName, prop...)
	case "enum":
		// 处理枚举类型
		enumType, ok := fieldType.(*avro.EnumSchema)
		if !ok {
			return nil, fmt.Errorf("expected enum schema for field '%s'", fieldName)
		}

		symbols := enumType.Symbols()
		prop = append(prop, mcp.Enum(symbols...))
		opt = mcp.WithString(fieldName, prop...)
	case "union":
		// 处理联合类型
		unionType, ok := fieldType.(*avro.UnionSchema)
		if !ok {
			return nil, fmt.Errorf("expected union schema for field '%s'", fieldName)
		}

		opt = mcp.WithString(fieldName, mcp.Description(description+" (union type): "+unionType.String()))
	default:
		// 对于未知或不支持的类型，默认使用字符串表示
		opt = mcp.WithString(fieldName, mcp.Description(description+" ("+string(fieldType.Type())+" type)"))
	}
	return opt, nil
}

// avroTypeToPropertyOption 将AVRO类型转换为MCP属性选项
// 用于处理嵌套结构中的类型，如数组元素和映射值
func (c *AvroConverter) avroTypeToPropertyOption(name string, avroType avro.Schema, description string) (map[string]any, error) {
	var props map[string]any = make(map[string]any)
	props["description"] = description
	switch avroType.Type() {
	case "string":
		props["type"] = "string"
		return props, nil
	case "int", "long":
		props["type"] = "number"
		return props, nil
	case "float", "double":
		props["type"] = "number"
		return props, nil
	case "boolean":
		props["type"] = "boolean"
		return props, nil
	case "bytes":
		props["type"] = "string"
		return props, nil
	case "array":
		// 处理嵌套数组
		arrayType, ok := avroType.(*avro.ArraySchema)
		if !ok {
			return nil, fmt.Errorf("expected array schema for '%s'", name)
		}

		itemsProp, err := c.avroTypeToPropertyOption("item", arrayType.Items(), description)
		if err != nil {
			return nil, err
		}

		props["type"] = "array"
		props["items"] = itemsProp
		return props, nil
	case "map":
		// 处理嵌套映射
		mapType, ok := avroType.(*avro.MapSchema)
		if !ok {
			return nil, fmt.Errorf("expected map schema for '%s'", name)
		}

		valuesProp, err := c.avroTypeToPropertyOption("value", mapType.Values(), description)
		if err != nil {
			return nil, err
		}

		props["type"] = "object"
		props["properties"] = valuesProp
		return props, nil
	case "record":
		// 处理嵌套记录
		recordType, ok := avroType.(*avro.RecordSchema)
		if !ok {
			return nil, fmt.Errorf("expected record schema for '%s'", name)
		}

		var fieldProps map[string]any
		for _, field := range recordType.Fields() {
			fieldProp, err := c.avroTypeToPropertyOption(field.Name(), field.Type(), field.Doc())
			if err != nil {
				return nil, err
			}
			fieldProps[field.Name()] = fieldProp
		}
		props["type"] = "object"
		props["properties"] = fieldProps

		return props, nil
	case "enum":
		// 处理枚举
		enumType, ok := avroType.(*avro.EnumSchema)
		if !ok {
			return nil, fmt.Errorf("expected enum schema for '%s'", name)
		}

		symbols := enumType.Symbols()
		props["type"] = "string"
		props["enum"] = symbols
		return props, nil
	case "union":
		// 简化处理联合类型
		props["type"] = "string"
		props["description"] = description + " (union type)"
		return props, nil
	default:
		// 默认为字符串
		props["type"] = "string"
		props["description"] = description + " (" + string(avroType.Type()) + " type)"
		return props, nil
	}
}

func (c *AvroConverter) SerializeMCPRequestToPulsarPayload(arguments map[string]any, targetPulsarSchemaInfo *pulsar.SchemaInfo) ([]byte, error) {
	if targetPulsarSchemaInfo.Type != pulsar.AVRO {
		return nil, fmt.Errorf("expected AVRO schema, got %s", GetSchemaType(targetPulsarSchemaInfo.Type))
	}

	// 首先验证参数
	if err := c.ValidateArguments(arguments, targetPulsarSchemaInfo); err != nil {
		return nil, fmt.Errorf("arguments validation failed: %w", err)
	}

	// 解析AVRO模式
	schema, err := avro.Parse(targetPulsarSchemaInfo.Schema)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AVRO schema: %w", err)
	}

	// 序列化为AVRO二进制数据
	data, err := avro.Marshal(schema, arguments)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments to AVRO: %w", err)
	}

	return data, nil
}

func (c *AvroConverter) ValidateArguments(arguments map[string]any, targetPulsarSchemaInfo *pulsar.SchemaInfo) error {
	if targetPulsarSchemaInfo.Type != pulsar.AVRO {
		return fmt.Errorf("expected AVRO schema, got %s", GetSchemaType(targetPulsarSchemaInfo.Type))
	}

	// 解析AVRO模式
	schema, err := avro.Parse(targetPulsarSchemaInfo.Schema)
	if err != nil {
		return fmt.Errorf("failed to parse AVRO schema: %w", err)
	}

	// 检查是否是记录类型
	recordSchema, ok := schema.(*avro.RecordSchema)
	if !ok {
		return fmt.Errorf("expected AVRO record schema, got %s", reflect.TypeOf(schema).String())
	}

	// 检查必需字段是否存在
	for _, field := range recordSchema.Fields() {
		fieldName := field.Name()
		fieldType := field.Type()

		// 检查字段是否必需
		// 在AVRO中，如果字段不是与null的联合类型，则通常认为它是必需的
		isRequired := true
		if unionType, ok := fieldType.(*avro.UnionSchema); ok {
			for _, t := range unionType.Types() {
				if t.Type() == "null" {
					isRequired = false
					break
				}
			}
		}

		// 如果字段是必需的但未提供，返回错误
		if isRequired {
			if _, exists := arguments[fieldName]; !exists {
				return fmt.Errorf("required field '%s' is missing", fieldName)
			}
		}

		// 如果提供了值，进行类型检查
		if value, exists := arguments[fieldName]; exists {
			if err := c.validateValueAgainstType(value, fieldType, fieldName); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateValueAgainstType 验证值是否符合AVRO类型要求
func (c *AvroConverter) validateValueAgainstType(value any, avroType avro.Schema, path string) error {
	// 处理null值
	if value == nil {
		if avroType.Type() == "null" {
			return nil
		}

		// 检查是否是包含null的联合类型
		if unionType, ok := avroType.(*avro.UnionSchema); ok {
			for _, t := range unionType.Types() {
				if t.Type() == "null" {
					return nil
				}
			}
		}

		return fmt.Errorf("field '%s' does not accept null value", path)
	}

	// 根据AVRO类型验证值
	switch avroType.Type() {
	case "null":
		if value != nil {
			return fmt.Errorf("expected null for field '%s', got %T", path, value)
		}
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("expected string for field '%s', got %T", path, value)
		}
	case "int":
		// 检查各种可以转换为int的类型
		switch v := value.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, float32, float64:
			// 这些类型都可以安全地转换为int
			return nil
		default:
			return fmt.Errorf("expected integer for field '%s', got %T: %v", path, v, v)
		}
	case "long":
		// 检查各种可以转换为long (int64)的类型
		switch v := value.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, float32, float64:
			// 这些类型都可以安全地转换为int64
			return nil
		default:
			return fmt.Errorf("expected long for field '%s', got %T: %v", path, v, v)
		}
	case "float", "double":
		// 检查各种可以转换为float/double的类型
		switch v := value.(type) {
		case float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32:
			// 这些类型都可以安全地转换为float/double
			return nil
		default:
			return fmt.Errorf("expected float/double for field '%s', got %T: %v", path, v, v)
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("expected boolean for field '%s', got %T", path, value)
		}
	case "bytes":
		// bytes可以用[]byte或base64编码的字符串表示
		switch v := value.(type) {
		case []byte, string:
			return nil
		default:
			return fmt.Errorf("expected bytes for field '%s', got %T", path, v)
		}
	case "array":
		// 检查值是否是切片或数组
		v := reflect.ValueOf(value)
		if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
			return fmt.Errorf("expected array for field '%s', got %T", path, value)
		}

		// 获取数组模式
		arrayType, ok := avroType.(*avro.ArraySchema)
		if !ok {
			return fmt.Errorf("expected array schema for field '%s'", path)
		}

		// 遍历数组元素进行验证
		itemType := arrayType.Items()
		for i := 0; i < v.Len(); i++ {
			if err := c.validateValueAgainstType(v.Index(i).Interface(), itemType, fmt.Sprintf("%s[%d]", path, i)); err != nil {
				return err
			}
		}
	case "map":
		// 检查值是否是映射
		v := reflect.ValueOf(value)
		if v.Kind() != reflect.Map {
			return fmt.Errorf("expected map for field '%s', got %T", path, value)
		}

		// 获取映射模式
		mapType, ok := avroType.(*avro.MapSchema)
		if !ok {
			return fmt.Errorf("expected map schema for field '%s'", path)
		}

		// 遍历映射键值对进行验证
		valueType := mapType.Values()
		for _, key := range v.MapKeys() {
			if key.Kind() != reflect.String {
				return fmt.Errorf("map keys must be strings for field '%s', got %s", path, key.Kind())
			}

			if err := c.validateValueAgainstType(v.MapIndex(key).Interface(), valueType, fmt.Sprintf("%s.%s", path, key.String())); err != nil {
				return err
			}
		}
	case "record":
		// 检查值是否是映射
		valueMap, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("expected map for record field '%s', got %T", path, value)
		}

		// 获取记录模式
		recordType, ok := avroType.(*avro.RecordSchema)
		if !ok {
			return fmt.Errorf("expected record schema for field '%s'", path)
		}

		// 对记录中的每个字段进行验证
		for _, field := range recordType.Fields() {
			fieldName := field.Name()
			fieldType := field.Type()

			// 检查字段是否必需
			isRequired := true
			if unionType, ok := fieldType.(*avro.UnionSchema); ok {
				for _, t := range unionType.Types() {
					if t.Type() == "null" {
						isRequired = false
						break
					}
				}
			}

			// 如果字段是必需的但未提供，返回错误
			if isRequired {
				if _, exists := valueMap[fieldName]; !exists {
					return fmt.Errorf("required field '%s.%s' is missing", path, fieldName)
				}
			}

			// 如果提供了值，进行类型检查
			if fieldValue, exists := valueMap[fieldName]; exists {
				if err := c.validateValueAgainstType(fieldValue, fieldType, fmt.Sprintf("%s.%s", path, fieldName)); err != nil {
					return err
				}
			}
		}
	case "enum":
		// 检查值是否是字符串
		strValue, ok := value.(string)
		if !ok {
			return fmt.Errorf("expected string for enum field '%s', got %T", path, value)
		}

		// 获取枚举模式
		enumType, ok := avroType.(*avro.EnumSchema)
		if !ok {
			return fmt.Errorf("expected enum schema for field '%s'", path)
		}

		// 检查值是否是有效的枚举符号
		symbols := enumType.Symbols()
		for _, symbol := range symbols {
			if symbol == strValue {
				return nil
			}
		}

		return fmt.Errorf("invalid enum value '%s' for field '%s', expected one of %v", strValue, path, symbols)
	case "union":
		// 处理联合类型
		unionType, ok := avroType.(*avro.UnionSchema)
		if !ok {
			return fmt.Errorf("expected union schema for field '%s'", path)
		}

		// 尝试用联合类型中的任一类型验证值
		var errors []error
		for _, t := range unionType.Types() {
			err := c.validateValueAgainstType(value, t, path)
			if err == nil {
				// 如果任一类型验证成功，返回nil
				return nil
			}
			errors = append(errors, err)
		}

		// 如果所有类型都验证失败，返回错误
		return fmt.Errorf("value does not match any type in union for field '%s'", path)
	default:
		// 对于未知或不支持的类型，假设可以接受任何值
		return nil
	}

	return nil
}
