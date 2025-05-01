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

func KafkaAdminAddTopicTools(s *server.MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeatureKafkaAdmin)) && !slices.Contains(features, string(FeatureAll)) && !slices.Contains(features, string(FeatureAllKafka)) {
		return
	}

	resourceDesc := "Resource to operate on. Available resources:\n" +
		"- topic: A single Kafka topic for operations on individual topics (create, get, delete)\n" +
		"- topics: Collection of Kafka topics for bulk operations (list)"

	operationDesc := "Operation to perform. Available operations:\n" +
		"- list: List all topics in the Kafka cluster, optionally including internal topics\n" +
		"- get: Get detailed configuration for a specific topic\n" +
		"- create: Create a new topic with specified partitions, replication factor, and optional configs\n" +
		"- delete: Delete an existing topic\n" +
		"- metadata: Get metadata for a specific topic\n"

	toolDesc := "Unified tool for managing Apache Kafka topics.\n" +
		"This tool provides access to various Kafka topic operations, including creation, deletion, listing, and configuration retrieval.\n" +
		"Kafka topics are distributed append-only logs that store ordered, immutable records. Topics can be partitioned for parallel processing.\n\n" +
		"Do not use this tool for Pulsar protocol operations. Use 'pulsar_admin_topics' instead." +
		"Usage Examples:\n\n" +
		"1. List all topics:\n" +
		"   resource: \"topics\"\n" +
		"   operation: \"list\"\n" +
		"   include-internal: false\n\n" +
		"2. Get topic configuration:\n" +
		"   resource: \"topic\"\n" +
		"   operation: \"get\"\n" +
		"   name: \"my-topic\"\n\n" +
		"3. Create a new topic with 3 partitions and replication factor 2:\n" +
		"   resource: \"topic\"\n" +
		"   operation: \"create\"\n" +
		"   name: \"new-topic\"\n" +
		"   partitions: 3\n" +
		"   replication-factor: 2\n" +
		"   configs: [\"cleanup.policy=compact\", \"retention.ms=604800000\"]\n\n" +
		"4. Delete a topic:\n" +
		"   resource: \"topic\"\n" +
		"   operation: \"delete\"\n" +
		"   name: \"topic-to-delete\"\n\n" +
		"5. Get full metadata including topic and broker information for a specific topic:\n" +
		"   resource: \"topic\"\n" +
		"   operation: \"metadata\"\n" +
		"   name: \"my-topic\"\n\n" +
		"This tool requires Kafka super-user permissions."

	kafkaTopicTool := mcp.NewTool("kafka_admin_topics",
		mcp.WithDescription(toolDesc),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description(resourceDesc),
		),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
		),
		mcp.WithBoolean("include-internal",
			mcp.Description("Whether to include internal Kafka topics (those starting with an underscore) in the results. "+
				"Only applicable for the 'list' operation on the 'topics' resource. "+
				"Default: false."),
		),
		mcp.WithString("name",
			mcp.Description("The name of the Kafka topic to operate on. "+
				"Required for 'get', 'create', and 'delete' operations on the 'topic' resource. "+
				"Must be a valid Kafka topic name."),
		),
		mcp.WithNumber("partitions",
			mcp.Description("The number of partitions to create for the topic. "+
				"Only applicable for the 'create' operation on the 'topic' resource. "+
				"Should be a positive integer. More partitions enable higher throughput but require more brokers. "+
				"Default: 1."),
		),
		mcp.WithNumber("replication-factor",
			mcp.Description("The replication factor for the topic, determining how many copies of data are maintained across brokers. "+
				"Only applicable for the 'create' operation on the 'topic' resource. "+
				"Should be a positive integer, not exceeding the number of available brokers. "+
				"Higher values increase fault tolerance but consume more storage. "+
				"Default: 1."),
		),
		mcp.WithArray("configs",
			mcp.Description("Topic configuration overrides as an array of key-value strings. "+
				"Only applicable for the 'create' operation on the 'topic' resource. "+
				"Each entry should be in the format 'key=value'. "+
				"Example configs include 'cleanup.policy', 'retention.ms', 'compression.type'. "+
				"Optional."),
		),
	)
	s.AddTool(kafkaTopicTool, handleKafkaTopicTool(readOnly))
}

func handleKafkaTopicTool(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		if readOnly && (operation == "create" || operation == "delete") {
			return mcp.NewToolResultError("Write operations are not allowed in read-only mode"), nil
		}

		// Create the admin client
		admin, err := kafka.GetKafkaAdminClient()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
		}

		// Dispatch based on resource and operation
		switch resource {
		case "topic":
			switch operation {
			case "get":
				return handleKafkaTopicGet(ctx, admin, request)
			case "create":
				return handleKafkaTopicCreate(ctx, admin, request)
			case "delete":
				return handleKafkaTopicDelete(ctx, admin, request)
			case "metadata":
				return handleKafkaTopicMetadata(ctx, admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'topic': %s", operation)), nil
			}
		case "topics":
			switch operation {
			case "list":
				return handleKafkaTopicsList(ctx, admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'topics': %s", operation)), nil
			}
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resource: %s. Available resources: topic, topics", resource)), nil
		}
	}
}

func handleKafkaTopicGet(ctx context.Context, admin *kadm.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topicName, err := requiredParam[string](request.Params.Arguments, "name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic name: %v", err)), nil
	}

	// Get topic configs
	configs, err := admin.DescribeTopicConfigs(ctx, topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic configs: %v", err)), nil
	}

	// Format the output
	jsonBytes, err := json.Marshal(configs)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal topic configs: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func handleKafkaTopicCreate(ctx context.Context, admin *kadm.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topicName, err := requiredParam[string](request.Params.Arguments, "name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic name: %v", err)), nil
	}

	// Get optional parameters
	partitions, ok := optionalParam[float64](request.Params.Arguments, "partitions")
	if !ok {
		partitions = 1 // Default to 1 partition
	}

	replicationFactor, ok := optionalParam[float64](request.Params.Arguments, "replication-factor")
	if !ok {
		replicationFactor = 1 // Default to replication factor 1
	}

	// Get configs if provided
	var configEntries map[string]*string
	configsArray, ok := optionalParamConfigs(request.Params.Arguments, "configs")
	if ok {
		configEntries, err = parseMessageConfigs(configsArray)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to parse configs: %v", err)), nil
		}
	}

	// Create the topic
	_, err = admin.CreateTopic(ctx, int32(partitions), int16(replicationFactor), configEntries, topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create topic: %v", err)), nil
	}

	return mcp.NewToolResultText("Topic created successfully"), nil
}

func handleKafkaTopicDelete(ctx context.Context, admin *kadm.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topicName, err := requiredParam[string](request.Params.Arguments, "name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic name: %v", err)), nil
	}

	// Delete the topic
	_, err = admin.DeleteTopic(ctx, topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete topic: %v", err)), nil
	}

	return mcp.NewToolResultText("Topic deleted successfully"), nil
}

func handleKafkaTopicsList(ctx context.Context, admin *kadm.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	includeInternal, ok := optionalParam[bool](request.Params.Arguments, "include-internal")
	if !ok {
		includeInternal = false
	}

	var topics kadm.TopicDetails
	var err error
	// List topics
	if !includeInternal {
		topics, err = admin.ListTopics(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list topics: %v", err)), nil
		}
	} else {
		topics, err = admin.ListTopicsWithInternal(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list topics: %v", err)), nil
		}
	}
	topicNames := []string{}
	for topicName := range topics {
		topicNames = append(topicNames, topicName)
	}

	// Format the output
	jsonBytes, err := json.Marshal(topicNames)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal topic names: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func handleKafkaTopicMetadata(ctx context.Context, admin *kadm.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topicName, err := requiredParam[string](request.Params.Arguments, "name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic name: %v", err)), nil
	}

	topicDetails, err := admin.Metadata(ctx, topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic metadata: %v", err)), nil
	}

	jsonBytes, err := json.Marshal(topicDetails)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal topic metadata: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}
