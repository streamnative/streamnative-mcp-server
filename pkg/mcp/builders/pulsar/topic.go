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

package pulsar

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

type pulsarAdminTopicInput struct {
	Resource       string  `json:"resource"`
	Operation      string  `json:"operation"`
	Topic          *string `json:"topic,omitempty"`
	Namespace      *string `json:"namespace,omitempty"`
	Partitions     *int    `json:"partitions,omitempty"`
	Force          bool    `json:"force,omitempty"`
	NonPartitioned bool    `json:"non-partitioned,omitempty"`
	Partitioned    bool    `json:"partitioned,omitempty"`
	PerPartition   bool    `json:"per-partition,omitempty"`
	Config         *string `json:"config,omitempty"`
	MessageID      *string `json:"messageId,omitempty"`
}

const (
	pulsarAdminTopicResourceDesc = "Resource to operate on. Available resources:\n" +
		"- topic: A Pulsar topic\n" +
		"- topics: Multiple topics within a namespace"
	pulsarAdminTopicOperationDesc = "Operation to perform. Available operations:\n" +
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
		"- update: Update a topic partitions\n" +
		"- offload: Offload data from a topic to long-term storage\n" +
		"- offload-status: Check the status of data offloading for a topic"
	pulsarAdminTopicTopicDesc = "The fully qualified topic name (format: [persistent|non-persistent]://tenant/namespace/topic). " +
		"Required for all operations except 'list'. " +
		"For partitioned topics, reference the base topic name without the partition suffix. " +
		"To operate on a specific partition, append -partition-N to the topic name."
	pulsarAdminTopicNamespaceDesc = "The namespace name in the format 'tenant/namespace'. " +
		"Required for the 'list' operation. " +
		"A namespace is a logical grouping of topics within a tenant."
	pulsarAdminTopicPartitionsDesc = "The number of partitions for the topic. Required for 'create' and 'update' operations. " +
		"Set to 0 for a non-partitioned topic. " +
		"Partitioned topics provide higher throughput by dividing message traffic across multiple brokers. " +
		"Each partition is an independent unit with its own retention and cursor positions."
	pulsarAdminTopicForceDesc = "Force operation even if it disrupts producers or consumers. Optional for 'delete' operation. " +
		"When true, all producers and consumers will be forcefully disconnected. " +
		"Use with caution as it can interrupt active message processing."
	pulsarAdminTopicNonPartitionedDesc = "Operate on a non-partitioned topic. Optional for 'delete' operation. " +
		"When true and operating on a partitioned topic name, only deletes the non-partitioned topic " +
		"with the same name, if it exists."
	pulsarAdminTopicPartitionedDesc = "Get stats for a partitioned topic. Optional for 'stats' operation. " +
		"It has to be true if the topic is partitioned. Leave it empty or false for non-partitioned topic."
	pulsarAdminTopicPerPartitionDesc = "Include per-partition stats. Optional for 'stats' operation. " +
		"When true, returns statistics for each partition separately. " +
		"Requires 'partitioned' parameter to be true."
	pulsarAdminTopicConfigDesc = "JSON configuration for the topic. Required for 'update' operation. " +
		"Set various policies like retention, compaction, deduplication, etc. " +
		"Use a JSON object format, e.g., '{\"deduplicationEnabled\": true, \"replication_clusters\": [\"us-west\", \"us-east\"]}'"
	pulsarAdminTopicMessageIDDesc = "Message ID for operations that require a position. Required for 'offload' operation. " +
		"Format is 'ledgerId:entryId' representing a position in the topic's message log. " +
		"For offload operations, specifies the message up to which data should be moved to long-term storage."
)

// PulsarAdminTopicToolBuilder implements the ToolBuilder interface for Pulsar Admin Topic tools
// It provides functionality to build Pulsar topic management tools
// /nolint:revive
type PulsarAdminTopicToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminTopicToolBuilder creates a new Pulsar Admin Topic tool builder instance
func NewPulsarAdminTopicToolBuilder() *PulsarAdminTopicToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_topic",
		Version:     "1.0.0",
		Description: "Pulsar Admin topic management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "topic", "admin"},
	}

	features := []string{
		"pulsar-admin-topics",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminTopicToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar Admin Topic tool list
// This is the core method implementing the ToolBuilder interface
func (b *PulsarAdminTopicToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool, err := b.buildTopicTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildTopicHandler(config.ReadOnly)

	return []builders.ToolDefinition{
		builders.ServerTool[pulsarAdminTopicInput, any]{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildTopicTool builds the Pulsar Admin Topic MCP tool definition
// Migrated from the original tool definition logic
func (b *PulsarAdminTopicToolBuilder) buildTopicTool() (*sdk.Tool, error) {
	inputSchema, err := buildPulsarAdminTopicInputSchema()
	if err != nil {
		return nil, err
	}

	toolDesc := "Manage Apache Pulsar topics. " +
		"Topics are the core messaging entities in Pulsar that store and transmit messages. " +
		"Pulsar supports two types of topics: persistent (durable storage with guaranteed delivery) " +
		"and non-persistent (in-memory with at-most-once delivery). " +
		"Topics can be partitioned for parallel processing and higher throughput, where each partition " +
		"functions as an independent topic with its own message log. " +
		"Topics follow a hierarchical naming structure: persistent://tenant/namespace/topic. " +
		"This tool supports various operations on topics including creation, deletion, lookup, compaction, " +
		"offloading, and retrieving statistics. " +
		"Do not use this tool for Kafka protocol operations. Use 'kafka_admin_topics' instead." +
		"Most operations require namespace admin permissions."

	return &sdk.Tool{
		Name:        "pulsar_admin_topic",
		Description: toolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildTopicHandler builds the Pulsar Admin Topic handler function
// Migrated from the original handler logic
func (b *PulsarAdminTopicToolBuilder) buildTopicHandler(readOnly bool) builders.ToolHandlerFunc[pulsarAdminTopicInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input pulsarAdminTopicInput) (*sdk.CallToolResult, any, error) {
		// Normalize parameters
		resource := strings.ToLower(input.Resource)
		operation := strings.ToLower(input.Operation)

		// Validate write operations in read-only mode
		if readOnly && (operation == "create" || operation == "delete" || operation == "unload" ||
			operation == "terminate" || operation == "compact" || operation == "update" || operation == "offload") {
			return nil, nil, fmt.Errorf("write operations are not allowed in read-only mode")
		}

		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return nil, nil, fmt.Errorf("pulsar session not found in context")
		}

		// Create the admin client
		admin, err := session.GetAdminClient()
		if err != nil {
			return nil, nil, b.handleError("get admin client", err)
		}

		// Dispatch based on resource and operation
		switch resource {
		case "topic":
			switch operation {
			case "get":
				result, err := b.handleTopicGet(admin, input)
				return result, nil, err
			case "create":
				result, err := b.handleTopicCreate(admin, input)
				return result, nil, err
			case "delete":
				result, err := b.handleTopicDelete(admin, input)
				return result, nil, err
			case "stats":
				result, err := b.handleTopicStats(admin, input)
				return result, nil, err
			case "lookup":
				result, err := b.handleTopicLookup(admin, input)
				return result, nil, err
			case "internal-stats":
				result, err := b.handleTopicInternalStats(admin, input)
				return result, nil, err
			case "internal-info":
				result, err := b.handleTopicInternalInfo(admin, input)
				return result, nil, err
			case "bundle-range":
				result, err := b.handleTopicBundleRange(admin, input)
				return result, nil, err
			case "last-message-id":
				result, err := b.handleTopicLastMessageID(admin, input)
				return result, nil, err
			case "status":
				result, err := b.handleTopicStatus(admin, input)
				return result, nil, err
			case "unload":
				result, err := b.handleTopicUnload(admin, input)
				return result, nil, err
			case "terminate":
				result, err := b.handleTopicTerminate(admin, input)
				return result, nil, err
			case "compact":
				result, err := b.handleTopicCompact(admin, input)
				return result, nil, err
			case "update":
				result, err := b.handleTopicUpdate(admin, input)
				return result, nil, err
			case "offload":
				result, err := b.handleTopicOffload(admin, input)
				return result, nil, err
			case "offload-status":
				result, err := b.handleTopicOffloadStatus(admin, input)
				return result, nil, err
			default:
				return nil, nil, fmt.Errorf("unknown topic operation: %s", operation)
			}
		case "topics":
			switch operation {
			case "list":
				result, err := b.handleTopicsList(admin, input)
				return result, nil, err
			default:
				return nil, nil, fmt.Errorf("unknown topics operation: %s", operation)
			}
		default:
			return nil, nil, fmt.Errorf("unknown resource: %s", resource)
		}
	}
}

func buildPulsarAdminTopicInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[pulsarAdminTopicInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}

	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}
	setSchemaDescription(schema, "resource", pulsarAdminTopicResourceDesc)
	setSchemaDescription(schema, "operation", pulsarAdminTopicOperationDesc)
	setSchemaDescription(schema, "topic", pulsarAdminTopicTopicDesc)
	setSchemaDescription(schema, "namespace", pulsarAdminTopicNamespaceDesc)
	setSchemaDescription(schema, "partitions", pulsarAdminTopicPartitionsDesc)
	setSchemaDescription(schema, "force", pulsarAdminTopicForceDesc)
	setSchemaDescription(schema, "non-partitioned", pulsarAdminTopicNonPartitionedDesc)
	setSchemaDescription(schema, "partitioned", pulsarAdminTopicPartitionedDesc)
	setSchemaDescription(schema, "per-partition", pulsarAdminTopicPerPartitionDesc)
	setSchemaDescription(schema, "config", pulsarAdminTopicConfigDesc)
	setSchemaDescription(schema, "messageId", pulsarAdminTopicMessageIDDesc)

	normalizeAdditionalProperties(schema)
	return schema, nil
}

// Unified error handling and utility functions

// handleError provides unified error handling
func (b *PulsarAdminTopicToolBuilder) handleError(operation string, err error) error {
	return fmt.Errorf("failed to %s: %v", operation, err)
}

// marshalResponse provides unified JSON serialization for responses
func (b *PulsarAdminTopicToolBuilder) marshalResponse(data interface{}) (*sdk.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, b.handleError("marshal response", err)
	}
	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: string(jsonBytes)}},
	}, nil
}

func textResult(message string) *sdk.CallToolResult {
	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: message}},
	}
}

func requireString(value *string, key string) (string, error) {
	if value == nil {
		return "", fmt.Errorf("required argument %q not found", key)
	}
	return *value, nil
}

func requireInt(value *int, key string) (int, error) {
	if value == nil {
		return 0, fmt.Errorf("required argument %q not found", key)
	}
	return *value, nil
}

// handleTopicsList lists all existing topics under the specified namespace
func (b *PulsarAdminTopicToolBuilder) handleTopicsList(admin cmdutils.Client, input pulsarAdminTopicInput) (*sdk.CallToolResult, error) {
	// Get required parameters
	namespace, err := requireString(input.Namespace, "namespace")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'namespace' for topics.list: %v", err)
	}

	// Get namespace name
	namespaceName, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return nil, fmt.Errorf("invalid namespace name '%s': %v", namespace, err)
	}

	// List topics
	partitionedTopics, nonPartitionedTopics, err := admin.Topics().List(*namespaceName)
	if err != nil {
		return nil, fmt.Errorf("failed to list topics in namespace '%s': %v",
			namespace, err)
	}

	// Format the output
	result := struct {
		PartitionedTopics    []string `json:"partitionedTopics"`
		NonPartitionedTopics []string `json:"nonPartitionedTopics"`
	}{
		PartitionedTopics:    partitionedTopics,
		NonPartitionedTopics: nonPartitionedTopics,
	}

	return b.marshalResponse(result)
}

// handleTopicGet gets the metadata of an existing topic
func (b *PulsarAdminTopicToolBuilder) handleTopicGet(admin cmdutils.Client, input pulsarAdminTopicInput) (*sdk.CallToolResult, error) {
	// Get required parameters
	topic, err := requireString(input.Topic, "topic")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'topic' for topic.get: %v", err)
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return nil, fmt.Errorf("invalid topic name '%s': %v", topic, err)
	}

	// Get topic metadata
	metadata, err := admin.Topics().GetMetadata(*topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata for topic '%s': %v",
			topic, err)
	}

	return b.marshalResponse(metadata)
}

// handleTopicStats gets the stats for an existing topic
func (b *PulsarAdminTopicToolBuilder) handleTopicStats(admin cmdutils.Client, input pulsarAdminTopicInput) (*sdk.CallToolResult, error) {
	// Get required parameters
	topic, err := requireString(input.Topic, "topic")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'topic' for topic.stats: %v", err)
	}

	// Get optional parameters
	partitioned := input.Partitioned
	perPartition := input.PerPartition

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return nil, fmt.Errorf("invalid topic name '%s': %v", topic, err)
	}

	namespaceName, err := utils.GetNamespaceName(topicName.GetTenant() + "/" + topicName.GetNamespace())
	if err != nil {
		return nil, fmt.Errorf("invalid namespace name: %v", err)
	}

	// List topics to determine if this topic is partitioned
	partitionedTopics, nonPartitionedTopics, err := admin.Topics().List(*namespaceName)
	if err != nil {
		return nil, fmt.Errorf("failed to list topics in namespace '%s': %v",
			namespaceName, err)
	}

	if slices.Contains(partitionedTopics, topicName.String()) {
		partitioned = true
	}
	if slices.Contains(nonPartitionedTopics, topicName.String()) {
		partitioned = false
	}

	var data interface{}
	if partitioned {
		// Get partitioned topic stats
		stats, err := admin.Topics().GetPartitionedStats(*topicName, perPartition)
		if err != nil {
			return nil, fmt.Errorf("failed to get stats for partitioned topic '%s': %v",
				topic, err)
		}
		data = stats
	} else {
		// Get topic stats
		stats, err := admin.Topics().GetStats(*topicName)
		if err != nil {
			return nil, fmt.Errorf("failed to get stats for topic '%s': %v",
				topic, err)
		}
		data = stats
	}

	return b.marshalResponse(data)
}

// handleTopicLookup looks up the owner broker of a topic
func (b *PulsarAdminTopicToolBuilder) handleTopicLookup(admin cmdutils.Client, input pulsarAdminTopicInput) (*sdk.CallToolResult, error) {
	// Get required parameters
	topic, err := requireString(input.Topic, "topic")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'topic' for topic.lookup: %v", err)
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return nil, fmt.Errorf("invalid topic name '%s': %v", topic, err)
	}

	// Lookup topic
	lookup, err := admin.Topics().Lookup(*topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup topic '%s': %v",
			topic, err)
	}

	return b.marshalResponse(lookup)
}

// handleTopicCreate creates a topic with the specified number of partitions
func (b *PulsarAdminTopicToolBuilder) handleTopicCreate(admin cmdutils.Client, input pulsarAdminTopicInput) (*sdk.CallToolResult, error) {
	// Get required parameters
	topic, err := requireString(input.Topic, "topic")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'topic' for topic.create: %v", err)
	}

	partitions, err := requireInt(input.Partitions, "partitions")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'partitions' for topic.create: %v", err)
	}

	// Validate partitions
	if partitions < 0 {
		return nil, fmt.Errorf("invalid partitions number: must be non-negative; use 0 for a non-partitioned topic")
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return nil, fmt.Errorf("invalid topic name '%s': %v", topic, err)
	}

	// Create topic
	err = admin.Topics().Create(*topicName, partitions)
	if err != nil {
		return nil, fmt.Errorf("failed to create topic '%s' with %d partitions: %v",
			topic, partitions, err)
	}

	if partitions == 0 {
		return textResult(fmt.Sprintf("Successfully created non-partitioned topic '%s'",
			topicName.String())), nil
	}
	return textResult(fmt.Sprintf("Successfully created topic '%s' with %d partitions",
		topicName.String(), partitions)), nil
}

// handleTopicDelete deletes a topic
func (b *PulsarAdminTopicToolBuilder) handleTopicDelete(admin cmdutils.Client, input pulsarAdminTopicInput) (*sdk.CallToolResult, error) {
	// Get required parameters
	topic, err := requireString(input.Topic, "topic")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'topic' for topic.delete: %v", err)
	}

	// Get optional parameters
	force := input.Force
	nonPartitioned := input.NonPartitioned

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return nil, fmt.Errorf("invalid topic name '%s': %v", topic, err)
	}

	// Delete topic
	err = admin.Topics().Delete(*topicName, force, nonPartitioned)
	if err != nil {
		return nil, fmt.Errorf("failed to delete topic '%s': %v", topic, err)
	}

	forceStr := ""
	if force {
		forceStr = " forcefully"
	}

	nonPartitionedStr := ""
	if nonPartitioned {
		nonPartitionedStr = " (non-partitioned)"
	}

	return textResult(fmt.Sprintf("Successfully deleted topic '%s'%s%s",
		topicName.String(), forceStr, nonPartitionedStr)), nil
}

// handleTopicUnload unloads a topic
func (b *PulsarAdminTopicToolBuilder) handleTopicUnload(admin cmdutils.Client, input pulsarAdminTopicInput) (*sdk.CallToolResult, error) {
	// Get required parameters
	topic, err := requireString(input.Topic, "topic")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'topic' for topic.unload: %v", err)
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return nil, fmt.Errorf("invalid topic name '%s': %v", topic, err)
	}

	// Unload topic
	err = admin.Topics().Unload(*topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to unload topic '%s': %v", topic, err)
	}

	return textResult(fmt.Sprintf("Successfully unloaded topic '%s'", topicName.String())), nil
}

// handleTopicTerminate terminates a topic
func (b *PulsarAdminTopicToolBuilder) handleTopicTerminate(admin cmdutils.Client, input pulsarAdminTopicInput) (*sdk.CallToolResult, error) {
	// Get required parameters
	topic, err := requireString(input.Topic, "topic")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'topic' for topic.terminate: %v", err)
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return nil, fmt.Errorf("invalid topic name '%s': %v", topic, err)
	}

	// Terminate topic
	messageID, err := admin.Topics().Terminate(*topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to terminate topic '%s': %v", topic, err)
	}

	// Convert message ID to string
	msgIDStr := fmt.Sprintf("%d:%d", messageID.LedgerID, messageID.EntryID)

	return textResult(fmt.Sprintf("Successfully terminated topic '%s' at message %s. "+
		"No more messages can be published to this topic.",
		topicName.String(), msgIDStr)), nil
}

// handleTopicCompact triggers compaction on a topic
func (b *PulsarAdminTopicToolBuilder) handleTopicCompact(admin cmdutils.Client, input pulsarAdminTopicInput) (*sdk.CallToolResult, error) {
	// Get required parameters
	topic, err := requireString(input.Topic, "topic")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'topic' for topic.compact: %v", err)
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return nil, fmt.Errorf("invalid topic name '%s': %v", topic, err)
	}

	// Compact topic
	err = admin.Topics().Compact(*topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to trigger compaction for topic '%s': %v", topic, err)
	}

	return textResult(fmt.Sprintf("Successfully triggered compaction for topic '%s'. "+
		"Run 'topic.status' to check compaction status.", topicName.String())), nil
}

// handleTopicInternalStats gets the internal stats for a topic
func (b *PulsarAdminTopicToolBuilder) handleTopicInternalStats(admin cmdutils.Client, input pulsarAdminTopicInput) (*sdk.CallToolResult, error) {
	// Get required parameters
	topic, err := requireString(input.Topic, "topic")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'topic' for topic.internal-stats: %v", err)
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return nil, fmt.Errorf("invalid topic name '%s': %v", topic, err)
	}

	// Get internal stats
	stats, err := admin.Topics().GetInternalStats(*topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to get internal stats for topic '%s': %v", topic, err)
	}

	return b.marshalResponse(stats)
}

// handleTopicInternalInfo gets the internal info for a topic
func (b *PulsarAdminTopicToolBuilder) handleTopicInternalInfo(admin cmdutils.Client, input pulsarAdminTopicInput) (*sdk.CallToolResult, error) {
	// Get required parameters
	topic, err := requireString(input.Topic, "topic")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'topic' for topic.internal-info: %v", err)
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return nil, fmt.Errorf("invalid topic name '%s': %v", topic, err)
	}

	// Get internal info
	info, err := admin.Topics().GetInternalInfo(*topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to get internal info for topic '%s': %v", topic, err)
	}

	return b.marshalResponse(info)
}

// handleTopicBundleRange gets the bundle range of a topic
func (b *PulsarAdminTopicToolBuilder) handleTopicBundleRange(admin cmdutils.Client, input pulsarAdminTopicInput) (*sdk.CallToolResult, error) {
	// Get required parameters
	topic, err := requireString(input.Topic, "topic")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'topic' for topic.bundle-range: %v", err)
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return nil, fmt.Errorf("invalid topic name '%s': %v", topic, err)
	}

	// Get bundle range
	bundle, err := admin.Topics().GetBundleRange(*topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to get bundle range for topic '%s': %v", topic, err)
	}

	return textResult(fmt.Sprintf("Bundle range for topic '%s': %s", topicName.String(), bundle)), nil
}

// handleTopicLastMessageID gets the last message ID of a topic
func (b *PulsarAdminTopicToolBuilder) handleTopicLastMessageID(admin cmdutils.Client, input pulsarAdminTopicInput) (*sdk.CallToolResult, error) {
	// Get required parameters
	topic, err := requireString(input.Topic, "topic")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'topic' for topic.last-message-id: %v", err)
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return nil, fmt.Errorf("invalid topic name '%s': %v", topic, err)
	}

	// Get last message ID
	messageID, err := admin.Topics().GetLastMessageID(*topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to get last message ID for topic '%s': %v", topic, err)
	}

	return b.marshalResponse(messageID)
}

// handleTopicStatus gets the status of a topic
func (b *PulsarAdminTopicToolBuilder) handleTopicStatus(admin cmdutils.Client, input pulsarAdminTopicInput) (*sdk.CallToolResult, error) {
	// Get required parameters
	topic, err := requireString(input.Topic, "topic")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'topic' for topic.status: %v", err)
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return nil, fmt.Errorf("invalid topic name '%s': %v", topic, err)
	}

	// Get topic metadata for status check
	metadata, err := admin.Topics().GetMetadata(*topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to get status for topic '%s': %v", topic, err)
	}

	// Create status object with available information
	status := struct {
		Metadata interface{} `json:"metadata"`
		Active   bool        `json:"active"`
	}{
		Metadata: metadata,
		Active:   true,
	}

	return b.marshalResponse(status)
}

// handleTopicUpdate updates a topic configuration
func (b *PulsarAdminTopicToolBuilder) handleTopicUpdate(admin cmdutils.Client, input pulsarAdminTopicInput) (*sdk.CallToolResult, error) {
	// Get required parameters
	topic, err := requireString(input.Topic, "topic")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'topic' for topic.update: %v", err)
	}

	partitions, err := requireInt(input.Partitions, "partitions")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'partitions' for topic.update: %v", err)
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return nil, fmt.Errorf("invalid topic name '%s': %v", topic, err)
	}

	err = admin.Topics().Update(*topicName, partitions)
	if err != nil {
		return nil, fmt.Errorf("failed to update topic '%s': %v", topic, err)
	}

	return textResult(fmt.Sprintf("Successfully updated topic '%s' partitions to %d",
		topicName.String(), partitions)), nil
}

// handleTopicOffload offloads data from a topic to long-term storage
func (b *PulsarAdminTopicToolBuilder) handleTopicOffload(admin cmdutils.Client, input pulsarAdminTopicInput) (*sdk.CallToolResult, error) {
	// Get required parameters
	topic, err := requireString(input.Topic, "topic")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'topic' for topic.offload: %v", err)
	}

	messageIDStr, err := requireString(input.MessageID, "messageId")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'messageId' for topic.offload: %v", err)
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return nil, fmt.Errorf("invalid topic name '%s': %v", topic, err)
	}

	// Parse message ID from format "ledgerId:entryId"
	var ledgerID, entryID int64
	if _, err := fmt.Sscanf(messageIDStr, "%d:%d", &ledgerID, &entryID); err != nil {
		return nil, fmt.Errorf("invalid message ID format (expected 'ledgerId:entryId'): %v. "+
			"Valid examples: '123:456'", err)
	}

	// Create MessageID object
	messageID := utils.MessageID{
		LedgerID: ledgerID,
		EntryID:  entryID,
	}

	// Offload topic
	err = admin.Topics().Offload(*topicName, messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to trigger offload for topic '%s': %v", topic, err)
	}

	return textResult(fmt.Sprintf("Successfully triggered offload for topic '%s' up to message %s. "+
		"Use 'topic.offload-status' to check the offload progress.",
		topicName.String(), messageIDStr)), nil
}

// handleTopicOffloadStatus checks the status of data offloading for a topic
func (b *PulsarAdminTopicToolBuilder) handleTopicOffloadStatus(admin cmdutils.Client, input pulsarAdminTopicInput) (*sdk.CallToolResult, error) {
	// Get required parameters
	topic, err := requireString(input.Topic, "topic")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'topic' for topic.offload-status: %v", err)
	}

	// Get topic name
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return nil, fmt.Errorf("invalid topic name '%s': %v", topic, err)
	}

	// Get offload status
	status, err := admin.Topics().OffloadStatus(*topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to get offload status for topic '%s': %v", topic, err)
	}

	return b.marshalResponse(status)
}

func setSchemaDescription(schema *jsonschema.Schema, name, desc string) {
	if schema == nil {
		return
	}
	prop, ok := schema.Properties[name]
	if !ok || prop == nil {
		return
	}
	prop.Description = desc
}

func normalizeAdditionalProperties(schema *jsonschema.Schema) {
	visited := map[*jsonschema.Schema]bool{}
	var walk func(*jsonschema.Schema)
	walk = func(s *jsonschema.Schema) {
		if s == nil || visited[s] {
			return
		}
		visited[s] = true

		if s.Type == "object" && s.Properties != nil && isFalseSchema(s.AdditionalProperties) {
			s.AdditionalProperties = nil
		}

		for _, prop := range s.Properties {
			walk(prop)
		}
		for _, prop := range s.PatternProperties {
			walk(prop)
		}
		for _, def := range s.Defs {
			walk(def)
		}
		for _, def := range s.Definitions {
			walk(def)
		}
		if s.AdditionalProperties != nil && !isFalseSchema(s.AdditionalProperties) {
			walk(s.AdditionalProperties)
		}
		if s.Items != nil {
			walk(s.Items)
		}
		for _, item := range s.PrefixItems {
			walk(item)
		}
		if s.AdditionalItems != nil {
			walk(s.AdditionalItems)
		}
		if s.UnevaluatedItems != nil {
			walk(s.UnevaluatedItems)
		}
		if s.UnevaluatedProperties != nil {
			walk(s.UnevaluatedProperties)
		}
		if s.PropertyNames != nil {
			walk(s.PropertyNames)
		}
		if s.Contains != nil {
			walk(s.Contains)
		}
		for _, subschema := range s.AllOf {
			walk(subschema)
		}
		for _, subschema := range s.AnyOf {
			walk(subschema)
		}
		for _, subschema := range s.OneOf {
			walk(subschema)
		}
		if s.Not != nil {
			walk(s.Not)
		}
		if s.If != nil {
			walk(s.If)
		}
		if s.Then != nil {
			walk(s.Then)
		}
		if s.Else != nil {
			walk(s.Else)
		}
		for _, subschema := range s.DependentSchemas {
			walk(subschema)
		}
	}
	walk(schema)
}

func isFalseSchema(schema *jsonschema.Schema) bool {
	if schema == nil || schema.Not == nil {
		return false
	}
	if !reflect.ValueOf(*schema.Not).IsZero() {
		return false
	}
	clone := *schema
	clone.Not = nil
	return reflect.ValueOf(clone).IsZero()
}
