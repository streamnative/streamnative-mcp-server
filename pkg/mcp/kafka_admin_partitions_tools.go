// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	"github.com/twmb/franz-go/pkg/kadm"
)

func KafkaAdminAddPartitionsTools(s *server.MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeatureKafkaAdmin)) && !slices.Contains(features, string(FeatureAll)) && !slices.Contains(features, string(FeatureAllKafka)) {
		return
	}

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
		"- Maximum number of partitions is limited by broker configuration\n\n" +
		"Usage Examples:\n\n" +
		"1. Update the number of partitions for an existing Kafka topic from 1 to 3:\n" +
		"   resource: \"partition\"\n" +
		"   operation: \"update\"\n" +
		"   topic: \"my-topic\"\n" +
		"   new-total: 3\n\n" +
		"This tool requires Kafka super-user permissions.\n" +
		"If you want to list or describe partitions of a topic, please use the `kafka_admin_topics` tool's `metadata` operation instead."

	kafkaPartitionsTool := mcp.NewTool("kafka_admin_partitions",
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

	s.AddTool(kafkaPartitionsTool, handleKafkaPartitionsTool(readOnly))
}

func handleKafkaPartitionsTool(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get required parameters
		resource, err := requiredParam[string](request.Params.Arguments, "resource")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get resource: %v", err)), nil
		}

		operation, err := requiredParam[string](request.Params.Arguments, "operation")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get operation: %v", err)), nil
		}

		// Normalize parameters
		resource = strings.ToLower(resource)
		operation = strings.ToLower(operation)

		// Validate write operations in read-only mode
		if readOnly && (operation == "update") {
			return mcp.NewToolResultError("Write operations are not allowed in read-only mode"), nil
		}

		// Create the admin client
		admin, err := kafka.GetKafkaAdminClient()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
		}

		// Dispatch based on resource and operation
		switch resource {
		case "partition":
			switch operation {
			case "update":
				return handleKafkaPartitionUpdate(ctx, admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'partition': %s", operation)), nil
			}
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resource: %s. Available resources: partition", resource)), nil
		}
	}
}

func handleKafkaPartitionUpdate(ctx context.Context, admin *kadm.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topicName, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic name: %v", err)), nil
	}

	newTotal, err := requiredParam[int](request.Params.Arguments, "new-total")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get new total: %v", err)), nil
	}

	response, err := admin.UpdatePartitions(ctx, newTotal, topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to update partitions: %v", err)), nil
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}
