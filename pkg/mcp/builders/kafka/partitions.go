// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	"github.com/twmb/franz-go/pkg/kadm"
)

// KafkaPartitionsToolBuilder implements the ToolBuilder interface for Kafka Partitions
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
func (b *KafkaPartitionsToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool := b.buildKafkaPartitionsTool()
	handler := b.buildKafkaPartitionsHandler(config.ReadOnly)

	return []server.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildKafkaPartitionsTool builds the Kafka Partitions MCP tool definition
func (b *KafkaPartitionsToolBuilder) buildKafkaPartitionsTool() mcp.Tool {
	resourceDesc := "Resource to operate on. Available resources:\n" +
		"- partition: A single Kafka Partition of a Kafka topic. Partitions are the basic unit of parallelism and data distribution in Kafka."

	operationDesc := "Operation to perform. Available operations:\n" +
		"- update: Update the number of partitions for an existing Kafka topic. This operation can only increase the number of partitions, not decrease them."

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

	return mcp.NewTool("kafka_admin_partitions",
		mcp.WithDescription(toolDesc),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description(resourceDesc),
		),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
		),
		mcp.WithString("topic",
			mcp.Description("The name of the Kafka topic to operate on. "+
				"Required for the 'update' operation. "+
				"Must be an existing topic name in the Kafka cluster. "+
				"The topic must be in a healthy state for partition updates to succeed.")),
		mcp.WithNumber("new-total",
			mcp.Description("The new total number of partitions for the Kafka topic. "+
				"Required for the 'update' operation. "+
				"Must be greater than the current number of partitions (cannot decrease partitions). "+
				"A larger number of partitions can help increase parallelism and throughput, but may also "+
				"increase resource utilization on the brokers. "+
				"Consider Kafka cluster capacity when setting this value.")),
	)
}

// buildKafkaPartitionsHandler builds the Kafka Partitions handler function
func (b *KafkaPartitionsToolBuilder) buildKafkaPartitionsHandler(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get required parameters
		resource, err := request.RequireString("resource")
		if err != nil {
			return b.handleError("get resource", err), nil
		}

		operation, err := request.RequireString("operation")
		if err != nil {
			return b.handleError("get operation", err), nil
		}

		// Normalize parameters
		resource = strings.ToLower(resource)
		operation = strings.ToLower(operation)

		// Validate write operations in read-only mode
		if readOnly && operation == "update" {
			return mcp.NewToolResultError("Write operations are not allowed in read-only mode"), nil
		}

		// Get Kafka admin client
		session := mcpCtx.GetKafkaSession(ctx)
		if session == nil {
			return b.handleError("get Kafka session not found in context", nil), nil
		}
		admin, err := session.GetAdminClient()
		if err != nil {
			return b.handleError("get admin client", err), nil
		}

		// Dispatch based on resource and operation
		switch resource {
		case "partition":
			switch operation {
			case "update":
				return b.handleKafkaPartitionUpdate(ctx, admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'partition': %s", operation)), nil
			}
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resource: %s. Available resources: partition", resource)), nil
		}
	}
}

// Utility functions

// handleError provides unified error handling
func (b *KafkaPartitionsToolBuilder) handleError(operation string, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("Failed to %s: %v", operation, err))
}

// marshalResponse provides unified JSON serialization for responses
func (b *KafkaPartitionsToolBuilder) marshalResponse(data interface{}) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return b.handleError("marshal response", err), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleKafkaPartitionUpdate handles updating the number of partitions for a topic
func (b *KafkaPartitionsToolBuilder) handleKafkaPartitionUpdate(ctx context.Context, admin *kadm.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := request.RequireString("topic")
	if err != nil {
		return b.handleError("get topic name", err), nil
	}

	newTotal, err := request.RequireInt("new-total")
	if err != nil {
		return b.handleError("get new total", err), nil
	}

	response, err := admin.UpdatePartitions(ctx, newTotal, topicName)
	if err != nil {
		return b.handleError("update Kafka topic partitions", err), nil
	}

	return b.marshalResponse(response)
}
