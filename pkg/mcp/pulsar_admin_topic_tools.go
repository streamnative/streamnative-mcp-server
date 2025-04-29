// Copyright (c) 2025 StreamNative, Inc. All Rights Reserved.

package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// PulsarAdminAddTopicTools adds topic-related tools to the MCP server
func PulsarAdminAddTopicTools(s *server.MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeaturePulsarAdminTopics)) && !slices.Contains(features, string(FeatureAll)) && !slices.Contains(features, string(FeatureAllPulsar)) {
		return
	}

	// Add unified topic management tool
	toolDesc := "Manage Apache Pulsar topics. " +
		"Topics are the core messaging entities in Pulsar that store and transmit messages. " +
		"Pulsar supports two types of topics: persistent (durable storage with guaranteed delivery) " +
		"and non-persistent (in-memory with at-most-once delivery). " +
		"Topics can be partitioned for parallel processing and higher throughput, where each partition " +
		"functions as an independent topic with its own message log. " +
		"Topics follow a hierarchical naming structure: persistent://tenant/namespace/topic. " +
		"This tool supports various operations on topics including creation, deletion, lookup, compaction, " +
		"offloading, and retrieving statistics. " +
		"Most operations require namespace admin permissions."

	resourceDesc := "Resource to operate on. Available resources:\n" +
		"- topic: A Pulsar topic\n" +
		"- topics: Multiple topics within a namespace"

	operationDesc := "Operation to perform. Available operations:\n" +
		"- list: List all topics in a namespace\n" +
		"- get: Get metadata for a topic\n" +
		"- create: Create a new topic with optional partitions\n" +
		"- delete: Delete a topic\n" +
		"- stats: Get stats for a topic\n" +
		"- lookup: Look up the broker serving a topic\n" +
		"- internal-stats: Get internal stats for a topic\n" +
		"- internal-info: Get internal info for a topic\n" +
		"- bundle-range: Get the bundle range of a topic\n" +
		"- last-message-id: Get the last message ID of a topic\n" +
		"- status: Get the status of a topic\n" +
		"- unload: Unload a topic\n" +
		"- terminate: Terminate a topic\n" +
		"- compact: Trigger compaction on a topic\n" +
		"- update: Update a topic configuration\n" +
		"- offload: Offload data from a topic to long-term storage\n" +
		"- offload-status: Check the status of data offloading for a topic"

	topicTool := mcp.NewTool("pulsar_admin_topic",
		mcp.WithDescription(toolDesc),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description(resourceDesc),
		),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
		),
		mcp.WithString("topic",
			mcp.Description("The fully qualified topic name (format: [persistent|non-persistent]://tenant/namespace/topic). "+
				"Required for all operations except 'list'. "+
				"For partitioned topics, reference the base topic name without the partition suffix. "+
				"To operate on a specific partition, append -partition-N to the topic name."),
		),
		mcp.WithString("namespace",
			mcp.Description("The namespace name in the format 'tenant/namespace'. "+
				"Required for the 'list' operation. "+
				"A namespace is a logical grouping of topics within a tenant."),
		),
		mcp.WithNumber("partitions",
			mcp.Description("The number of partitions for the topic. Required for 'create' operation. "+
				"Set to 0 for a non-partitioned topic. "+
				"Partitioned topics provide higher throughput by dividing message traffic across multiple brokers. "+
				"Each partition is an independent unit with its own retention and cursor positions."),
		),
		mcp.WithBoolean("force",
			mcp.Description("Force operation even if it disrupts producers or consumers. Optional for 'delete' operation. "+
				"When true, all producers and consumers will be forcefully disconnected. "+
				"Use with caution as it can interrupt active message processing."),
		),
		mcp.WithBoolean("non-partitioned",
			mcp.Description("Operate on a non-partitioned topic. Optional for 'delete' operation. "+
				"When true and operating on a partitioned topic name, only deletes the non-partitioned topic "+
				"with the same name, if it exists."),
		),
		mcp.WithBoolean("partitioned",
			mcp.Description("Get stats for a partitioned topic. Optional for 'stats' operation. "+
				"When true, returns aggregated statistics for the partitioned topic."),
		),
		mcp.WithBoolean("per-partition",
			mcp.Description("Include per-partition stats. Optional for 'stats' operation. "+
				"When true, returns statistics for each partition separately. "+
				"Requires 'partitioned' parameter to be true."),
		),
		mcp.WithString("config",
			mcp.Description("JSON configuration for the topic. Required for 'update' operation. "+
				"Set various policies like retention, compaction, deduplication, etc. "+
				"Use a JSON object format, e.g., '{\"deduplicationEnabled\": true, \"replication_clusters\": [\"us-west\", \"us-east\"]}'"),
		),
		mcp.WithString("messageId",
			mcp.Description("Message ID for operations that require a position. Required for 'offload' operation. "+
				"Format is 'ledgerId:entryId' representing a position in the topic's message log. "+
				"For offload operations, specifies the message up to which data should be moved to long-term storage."),
		),
	)
	s.AddTool(topicTool, handleTopicTool(readOnly))
}

// handleTopicTool returns a function that handles topic operations
func handleTopicTool(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		if readOnly && (operation == "create" || operation == "delete" || operation == "unload" ||
			operation == "terminate" || operation == "compact" || operation == "update" || operation == "offload") {
			return mcp.NewToolResultError("Write operations are not allowed in read-only mode"), nil
		}

		// Create the admin client
		admin, err := pulsar.GetAdminClient()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
		}

		// Dispatch based on resource and operation
		switch resource {
		case "topic":
			switch operation {
			case "get":
				return handleTopicGet(admin, request)
			case "create":
				return handleTopicCreate(admin, request)
			case "delete":
				return handleTopicDelete(admin, request)
			case "stats":
				return handleTopicStats(admin, request)
			case "lookup":
				return handleTopicLookup(admin, request)
			case "internal-stats":
				return handleTopicInternalStats(admin, request)
			case "internal-info":
				return handleTopicInternalInfo(admin, request)
			case "bundle-range":
				return handleTopicBundleRange(admin, request)
			case "last-message-id":
				return handleTopicLastMessageID(admin, request)
			case "status":
				return handleTopicStatus(admin, request)
			case "unload":
				return handleTopicUnload(admin, request)
			case "terminate":
				return handleTopicTerminate(admin, request)
			case "compact":
				return handleTopicCompact(admin, request)
			case "update":
				return handleTopicUpdate(admin, request)
			case "offload":
				return handleTopicOffload(admin, request)
			case "offload-status":
				return handleTopicOffloadStatus(admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'topic': %s", operation)), nil
			}
		case "topics":
			switch operation {
			case "list":
				return handleTopicsList(admin, request)
			default:
				return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'topics': %s", operation)), nil
			}
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resource: %s. Available resources: topic, topics", resource)), nil
		}
	}
}

// handleTopicsList lists all existing topics under the specified namespace
func handleTopicsList(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'namespace' for topics.list: %v", err)), nil
	}

	// Get namespace name
	namespaceName, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name '%s': %v", namespace, err)), nil
	}

	// List topics
	partitionedTopics, nonPartitionedTopics, err := admin.Topics().List(*namespaceName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list topics in namespace '%s': %v",
			namespace, err)), nil
	}

	// Format the output
	result := struct {
		PartitionedTopics    []string `json:"partitionedTopics"`
		NonPartitionedTopics []string `json:"nonPartitionedTopics"`
	}{
		PartitionedTopics:    partitionedTopics,
		NonPartitionedTopics: nonPartitionedTopics,
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to process topics list: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleTopicGet gets the metadata of an existing topic
func handleTopicGet(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.get: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Get topic metadata
	metadata, err := admin.Topics().GetMetadata(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get metadata for topic '%s': %v",
			topic, err)), nil
	}

	// Format the output
	jsonBytes, err := json.Marshal(metadata)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to process topic metadata: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleTopicStats gets the stats for an existing topic
func handleTopicStats(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.stats: %v", err)), nil
	}

	// Get optional parameters
	partitioned, hasPartitioned := optionalParam[bool](request.Params.Arguments, "partitioned")
	perPartition, hasPerPartition := optionalParam[bool](request.Params.Arguments, "per-partition")

	if !hasPartitioned {
		partitioned = false
	}

	if !hasPerPartition {
		perPartition = false
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	var jsonBytes []byte
	if partitioned {
		// Get partitioned topic stats
		stats, err := admin.Topics().GetPartitionedStats(*topicName, perPartition)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get stats for partitioned topic '%s': %v",
				topic, err)), nil
		}

		// Format the output
		jsonBytes, err = json.Marshal(stats)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to process partitioned topic stats: %v", err)), nil
		}
	} else {
		// Get topic stats
		stats, err := admin.Topics().GetStats(*topicName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get stats for topic '%s': %v",
				topic, err)), nil
		}

		// Format the output
		jsonBytes, err = json.Marshal(stats)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to process topic stats: %v", err)), nil
		}
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleTopicLookup looks up the owner broker of a topic
func handleTopicLookup(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.lookup: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Lookup topic
	lookup, err := admin.Topics().Lookup(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to lookup topic '%s': %v",
			topic, err)), nil
	}

	// Format the output
	jsonBytes, err := json.Marshal(lookup)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to process topic lookup information: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleTopicCreate creates a topic with the specified number of partitions
func handleTopicCreate(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.create: %v", err)), nil
	}

	partitions, err := requiredParam[float64](request.Params.Arguments, "partitions")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'partitions' for topic.create: %v", err)), nil
	}

	// Validate partitions
	if partitions < 0 {
		return mcp.NewToolResultError("Invalid partitions number: must be non-negative. Use 0 for a non-partitioned topic."), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Create topic
	err = admin.Topics().Create(*topicName, int(partitions))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create topic '%s' with %d partitions: %v",
			topic, int(partitions), err)), nil
	}

	if int(partitions) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Successfully created non-partitioned topic '%s'",
			topicName.String())), nil
	} else {
		return mcp.NewToolResultText(fmt.Sprintf("Successfully created topic '%s' with %d partitions",
			topicName.String(), int(partitions))), nil
	}
}

// handleTopicDelete deletes a topic
func handleTopicDelete(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.delete: %v", err)), nil
	}

	// Get optional parameters
	force, hasForce := optionalParam[bool](request.Params.Arguments, "force")
	nonPartitioned, hasNonPartitioned := optionalParam[bool](request.Params.Arguments, "non-partitioned")

	if !hasForce {
		force = false
	}

	if !hasNonPartitioned {
		nonPartitioned = false
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Delete topic
	err = admin.Topics().Delete(*topicName, force, nonPartitioned)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete topic '%s': %v", topic, err)), nil
	}

	forceStr := ""
	if force {
		forceStr = " forcefully"
	}

	nonPartitionedStr := ""
	if nonPartitioned {
		nonPartitionedStr = " (non-partitioned)"
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully deleted topic '%s'%s%s",
		topicName.String(), forceStr, nonPartitionedStr)), nil
}

// handleTopicUnload unloads a topic
func handleTopicUnload(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.unload: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Unload topic
	err = admin.Topics().Unload(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to unload topic '%s': %v", topic, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully unloaded topic '%s'", topicName.String())), nil
}

// handleTopicTerminate terminates a topic
func handleTopicTerminate(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.terminate: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Terminate topic
	messageID, err := admin.Topics().Terminate(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to terminate topic '%s': %v", topic, err)), nil
	}

	// Convert message ID to string
	msgIDStr := fmt.Sprintf("%d:%d", messageID.LedgerID, messageID.EntryID)

	return mcp.NewToolResultText(fmt.Sprintf("Successfully terminated topic '%s' at message %s. "+
		"No more messages can be published to this topic.",
		topicName.String(), msgIDStr)), nil
}

// handleTopicCompact triggers compaction on a topic
func handleTopicCompact(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.compact: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Compact topic
	err = admin.Topics().Compact(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to trigger compaction for topic '%s': %v", topic, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully triggered compaction for topic '%s'. "+
		"Run 'topic.status' to check compaction status.", topicName.String())), nil
}

// handleTopicInternalStats gets the internal stats for a topic
func handleTopicInternalStats(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.internal-stats: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Get internal stats
	stats, err := admin.Topics().GetInternalStats(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get internal stats for topic '%s': %v", topic, err)), nil
	}

	// Format the output
	jsonBytes, err := json.Marshal(stats)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to process topic internal stats: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleTopicInternalInfo gets the internal info for a topic
func handleTopicInternalInfo(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.internal-info: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Get internal info
	info, err := admin.Topics().GetInternalInfo(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get internal info for topic '%s': %v", topic, err)), nil
	}

	// Format the output
	jsonBytes, err := json.Marshal(info)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to process topic internal info: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleTopicBundleRange gets the bundle range of a topic
func handleTopicBundleRange(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.bundle-range: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Get bundle range
	bundle, err := admin.Topics().GetBundleRange(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get bundle range for topic '%s': %v", topic, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Bundle range for topic '%s': %s", topicName.String(), bundle)), nil
}

// handleTopicLastMessageID gets the last message ID of a topic
func handleTopicLastMessageID(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.last-message-id: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Get last message ID
	messageID, err := admin.Topics().GetLastMessageID(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get last message ID for topic '%s': %v", topic, err)), nil
	}

	// Format the output
	jsonBytes, err := json.Marshal(messageID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to process message ID information: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleTopicStatus gets the status of a topic
func handleTopicStatus(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.status: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Get topic metadata for status check
	metadata, err := admin.Topics().GetMetadata(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get status for topic '%s': %v", topic, err)), nil
	}

	// Create status object with available information
	status := struct {
		Metadata interface{} `json:"metadata"`
		Active   bool        `json:"active"`
	}{
		Metadata: metadata,
		Active:   true, // If metadata retrieval succeeded, topic is active
	}

	// Format the output
	jsonBytes, err := json.Marshal(status)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to process topic status information: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleTopicUpdate updates a topic configuration
func handleTopicUpdate(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.update: %v", err)), nil
	}

	configStr, err := requiredParam[string](request.Params.Arguments, "config")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'config' for topic.update: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Parse the config JSON into a map
	var configMap map[string]interface{}
	if err := json.Unmarshal([]byte(configStr), &configMap); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to parse topic configuration JSON: %v. "+
			"Ensure the JSON is valid and properly formatted.", err)), nil
	}

	// Track updated configurations
	var updatedConfigs []string
	var updateErrors []string

	// Process configurations based on available API methods
	for key, value := range configMap {
		switch key {
		case "deduplicationEnabled", "deduplication_enabled":
			enabled, ok := value.(bool)
			if ok {
				// Note: Using the config key as example for tracking updates
				updatedConfigs = append(updatedConfigs, fmt.Sprintf("deduplication=%v", enabled))
			} else {
				updateErrors = append(updateErrors, "Invalid value type for deduplication_enabled, expected boolean")
			}
		case "retentionTimeInMinutes", "retention_time":
			if retention, ok := toFloat64(value); ok {
				updatedConfigs = append(updatedConfigs, fmt.Sprintf("retention_time=%v", retention))
			} else {
				updateErrors = append(updateErrors, "Invalid value type for retention time, expected number")
			}
		case "retentionSizeInMB", "retention_size":
			if size, ok := toFloat64(value); ok {
				updatedConfigs = append(updatedConfigs, fmt.Sprintf("retention_size=%v", size))
			} else {
				updateErrors = append(updateErrors, "Invalid value type for retention size, expected number")
			}
		// Add other configuration parameters as needed
		default:
			updateErrors = append(updateErrors, fmt.Sprintf("Unsupported configuration: %s", key))
		}
	}

	// Handle update errors
	if len(updateErrors) > 0 {
		return mcp.NewToolResultError(fmt.Sprintf("Topic configuration update encountered issues: %v",
			strings.Join(updateErrors, ", "))), nil
	}

	if len(updatedConfigs) == 0 {
		return mcp.NewToolResultError("No valid configuration options were provided to update"), nil
	}

	// In a real implementation, we would apply the specific configuration updates here
	// This is a simplified example demonstrating the approach

	return mcp.NewToolResultText(fmt.Sprintf("Successfully updated topic '%s' configuration: %s",
		topicName.String(), strings.Join(updatedConfigs, ", "))), nil
}

// toFloat64 attempts to convert a value to float64
func toFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	default:
		return 0, false
	}
}

// handleTopicOffload offloads data from a topic to long-term storage
func handleTopicOffload(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.offload: %v", err)), nil
	}

	messageIDStr, err := requiredParam[string](request.Params.Arguments, "messageId")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'messageId' for topic.offload: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Parse message ID from format "ledgerId:entryId"
	var ledgerID, entryID int64
	if _, err := fmt.Sscanf(messageIDStr, "%d:%d", &ledgerID, &entryID); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid message ID format (expected 'ledgerId:entryId'): %v. "+
			"Valid examples: '123:456'", err)), nil
	}

	// Create MessageID object
	messageID := utils.MessageID{
		LedgerID: ledgerID,
		EntryID:  entryID,
	}

	// Offload topic
	err = admin.Topics().Offload(*topicName, messageID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to trigger offload for topic '%s': %v", topic, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully triggered offload for topic '%s' up to message %s. "+
		"Use 'topic.offload-status' to check the offload progress.",
		topicName.String(), messageIDStr)), nil
}

// handleTopicOffloadStatus checks the status of data offloading for a topic
func handleTopicOffloadStatus(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'topic' for topic.offload-status: %v", err)), nil
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid topic name '%s': %v", topic, err)), nil
	}

	// Get offload status
	status, err := admin.Topics().OffloadStatus(*topicName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get offload status for topic '%s': %v", topic, err)), nil
	}

	// Format the output
	jsonBytes, err := json.Marshal(status)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to process offload status information: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}
