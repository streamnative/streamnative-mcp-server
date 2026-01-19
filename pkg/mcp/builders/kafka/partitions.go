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

package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	"github.com/twmb/franz-go/pkg/kadm"
)

type kafkaPartitionsInput struct {
	Resource  string  `json:"resource"`
	Operation string  `json:"operation"`
	Topic     *string `json:"topic,omitempty"`
	NewTotal  *int    `json:"new-total,omitempty"`
}

const (
	kafkaPartitionsResourceDesc = "Resource to operate on. Available resources:\n" +
		"- partition: A single Kafka Partition of a Kafka topic. Partitions are the basic unit of parallelism and data distribution in Kafka."
	kafkaPartitionsOperationDesc = "Operation to perform. Available operations:\n" +
		"- update: Update the number of partitions for an existing Kafka topic. This operation can only increase the number of partitions, not decrease them."
	kafkaPartitionsTopicDesc = "The name of the Kafka topic to operate on. " +
		"Required for the 'update' operation. " +
		"Must be an existing topic name in the Kafka cluster. " +
		"The topic must be in a healthy state for partition updates to succeed."
	kafkaPartitionsNewTotalDesc = "The new total number of partitions for the Kafka topic. " +
		"Required for the 'update' operation. " +
		"Must be greater than the current number of partitions (cannot decrease partitions). " +
		"A larger number of partitions can help increase parallelism and throughput, but may also " +
		"increase resource utilization on the brokers. " +
		"Consider Kafka cluster capacity when setting this value."
)

// KafkaPartitionsToolBuilder implements the ToolBuilder interface for Kafka Partitions
// /nolint:revive
type KafkaPartitionsToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewKafkaPartitionsToolBuilder creates a new Kafka Partitions tool builder instance
func NewKafkaPartitionsToolBuilder() *KafkaPartitionsToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "kafka_partitions",
		Version:     "1.0.0",
		Description: "Kafka Partitions administration tools",
		Category:    "kafka_admin",
		Tags:        []string{"kafka", "partitions", "admin"},
	}

	features := []string{
		"kafka-admin",
		"all",
		"all-kafka",
	}

	return &KafkaPartitionsToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Kafka Partitions tool list
func (b *KafkaPartitionsToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool, err := b.buildKafkaPartitionsTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildKafkaPartitionsHandler(config.ReadOnly)

	return []builders.ToolDefinition{
		builders.ServerTool[kafkaPartitionsInput, any]{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildKafkaPartitionsTool builds the Kafka Partitions MCP tool definition
func (b *KafkaPartitionsToolBuilder) buildKafkaPartitionsTool() (*sdk.Tool, error) {
	inputSchema, err := buildKafkaPartitionsInputSchema()
	if err != nil {
		return nil, err
	}

	toolDesc := "Unified tool for managing Apache Kafka partitions.\n" +
		"This tool provides access to Kafka partition operations, particularly adding partitions to existing topics.\n" +
		"Kafka partitions are the fundamental unit of parallelism and scalability in Kafka. Each partition is an ordered, " +
		"immutable sequence of records that is continually appended to. Partitions can be distributed across multiple brokers " +
		"to enable parallel processing of a topic.\n\n" +
		"Important notes:\n" +
		"- You can only increase the number of partitions, never decrease them\n" +
		"- Adding partitions may change the mapping of keys to partitions\n" +
		"- Existing messages remain in their original partitions\n" +
		"- Consider the impact on ordering guarantees when adding partitions\n\n" +
		"Usage Examples:\n\n" +
		"1. Increase the partition count for a topic:\n" +
		"   resource: \"partition\"\n" +
		"   operation: \"update\"\n" +
		"   topic: \"user-events\"\n" +
		"   partitions: 6\n\n" +
		"2. Scale up partitions for high-throughput topic:\n" +
		"   resource: \"partition\"\n" +
		"   operation: \"update\"\n" +
		"   topic: \"metrics-data\"\n" +
		"   partitions: 12\n\n" +
		"This tool requires appropriate Kafka permissions for partition management."

	return &sdk.Tool{
		Name:        "kafka_admin_partitions",
		Description: toolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildKafkaPartitionsHandler builds the Kafka Partitions handler function
func (b *KafkaPartitionsToolBuilder) buildKafkaPartitionsHandler(readOnly bool) builders.ToolHandlerFunc[kafkaPartitionsInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input kafkaPartitionsInput) (*sdk.CallToolResult, any, error) {
		// Normalize parameters
		resource := strings.ToLower(input.Resource)
		operation := strings.ToLower(input.Operation)

		// Validate write operations in read-only mode
		if readOnly && operation == "update" {
			return nil, nil, fmt.Errorf("write operations are not allowed in read-only mode")
		}

		// Get Kafka admin client
		session := mcpCtx.GetKafkaSession(ctx)
		if session == nil {
			return nil, nil, b.handleError("get Kafka session not found in context", nil)
		}
		admin, err := session.GetAdminClient()
		if err != nil {
			return nil, nil, b.handleError("get admin client", err)
		}

		// Dispatch based on resource and operation
		switch resource {
		case "partition":
			switch operation {
			case "update":
				result, err := b.handleKafkaPartitionUpdate(ctx, admin, input)
				return result, nil, err
			default:
				return nil, nil, fmt.Errorf("invalid operation for resource 'partition': %s", operation)
			}
		default:
			return nil, nil, fmt.Errorf("invalid resource: %s. available resources: partition", resource)
		}
	}
}

// Utility functions

// handleError provides unified error handling
func (b *KafkaPartitionsToolBuilder) handleError(operation string, err error) error {
	return fmt.Errorf("failed to %s: %v", operation, err)
}

// marshalResponse provides unified JSON serialization for responses
func (b *KafkaPartitionsToolBuilder) marshalResponse(data interface{}) (*sdk.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, b.handleError("marshal response", err)
	}
	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: string(jsonBytes)}},
	}, nil
}

// handleKafkaPartitionUpdate handles updating the number of partitions for a topic
func (b *KafkaPartitionsToolBuilder) handleKafkaPartitionUpdate(ctx context.Context, admin *kadm.Client, input kafkaPartitionsInput) (*sdk.CallToolResult, error) {
	topicName, err := requireString(input.Topic, "topic")
	if err != nil {
		return nil, b.handleError("get topic name", err)
	}

	newTotal, err := requireInt(input.NewTotal, "new-total")
	if err != nil {
		return nil, b.handleError("get new total", err)
	}

	response, err := admin.UpdatePartitions(ctx, newTotal, topicName)
	if err != nil {
		return nil, b.handleError("update Kafka topic partitions", err)
	}

	return b.marshalResponse(response)
}

func buildKafkaPartitionsInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[kafkaPartitionsInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}

	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}
	setSchemaDescription(schema, "resource", kafkaPartitionsResourceDesc)
	setSchemaDescription(schema, "operation", kafkaPartitionsOperationDesc)
	setSchemaDescription(schema, "topic", kafkaPartitionsTopicDesc)
	setSchemaDescription(schema, "new-total", kafkaPartitionsNewTotalDesc)

	normalizeAdditionalProperties(schema)
	return schema, nil
}
