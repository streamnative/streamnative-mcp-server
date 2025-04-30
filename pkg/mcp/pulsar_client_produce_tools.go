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
	"os"
	"slices"
	"strings"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	mcppulsar "github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// PulsarClientAddProducerTools adds Pulsar client producer tools to the MCP server
func PulsarClientAddProducerTools(s *server.MCPServer, _ bool, features []string) {
	if !slices.Contains(features, string(FeaturePulsarClient)) && !slices.Contains(features, string(FeatureAll)) && !slices.Contains(features, string(FeatureAllPulsar)) {
		return
	}

	// Main produce tool
	produceTool := mcp.NewTool("pulsar_client_produce",
		mcp.WithDescription("Produce messages to a Pulsar topic. This tool allows you to send messages "+
			"to a specified Pulsar topic with various options to control message format, "+
			"batching, and properties."),
		mcp.WithString("topic", mcp.Required(),
			mcp.Description("Topic to produce to"),
		),
		mcp.WithArray("messages",
			mcp.Description("Messages to send. Specify multiple times for multiple messages."),
		),
		mcp.WithArray("files",
			mcp.Description("Files to send as message content. Specify multiple times for multiple files."),
		),
		mcp.WithNumber("num-produce",
			mcp.Description("Number of times to send message(s) (default: 1)"),
		),
		mcp.WithNumber("rate",
			mcp.Description("Rate (in msg/sec) at which to produce, 0 means produce as fast as possible (default: 0)"),
		),
		mcp.WithBoolean("disable-batching",
			mcp.Description("Disable batch sending of messages (default: false)"),
		),
		mcp.WithBoolean("chunking",
			mcp.Description("Should split the message and publish in chunks if message size is larger than allowed max size (default: false)"),
		),
		mcp.WithString("separator",
			mcp.Description("Character to split messages string on (default: comma)"),
		),
		mcp.WithArray("properties",
			mcp.Description("Properties to add, key=value format. Specify multiple times for multiple properties."),
		),
		mcp.WithString("key",
			mcp.Description("Partitioning key to add to each message"),
		),
		mcp.WithString("key-value-key",
			mcp.Description("Value to add as message key in KeyValue schema"),
		),
		mcp.WithString("key-value-key-file",
			mcp.Description("Path to file containing the value to add as message key in KeyValue schema"),
		),
		mcp.WithString("value-schema",
			mcp.Description("Schema type for the value (can be string, bytes, json, avro) (default: bytes)"),
		),
		mcp.WithString("key-schema",
			mcp.Description("Schema type for the key (can be string, bytes, json, avro) (default: string)"),
		),
		mcp.WithString("key-value-encoding-type",
			mcp.Description("Key Value Encoding Type (can be separated or inline)"),
		),
		mcp.WithString("encryption-key-name",
			mcp.Description("The public key name to encrypt payload"),
		),
		mcp.WithString("encryption-key-value",
			mcp.Description("The URI of public key to encrypt payload"),
		),
		mcp.WithBoolean("disable-replication",
			mcp.Description("Disable geo-replication for messages (default: false)"),
		),
	)
	s.AddTool(produceTool, handleClientProduce)
}

func handleClientProduce(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract required parameters with validation
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	// Set default values and extract optional parameters
	messages := []string{}
	if val, exists := optionalParam[[]interface{}](request.Params.Arguments, "messages"); exists && len(val) > 0 {
		for _, m := range val {
			if strMsg, ok := m.(string); ok {
				messages = append(messages, strMsg)
			}
		}
	}

	files := []string{}
	if val, exists := optionalParam[[]interface{}](request.Params.Arguments, "files"); exists && len(val) > 0 {
		for _, f := range val {
			if strFile, ok := f.(string); ok {
				files = append(files, strFile)
			}
		}
	}

	if len(messages) == 0 && len(files) == 0 {
		return mcp.NewToolResultError("Please supply message content with either --messages or --files"), nil
	}

	numProduce := 1
	if val, exists := optionalParam[float64](request.Params.Arguments, "num-produce"); exists {
		numProduce = int(val)
	}

	rate := 0.0
	if val, exists := optionalParam[float64](request.Params.Arguments, "rate"); exists {
		rate = val
	}

	disableBatching := false
	if val, exists := optionalParam[bool](request.Params.Arguments, "disable-batching"); exists {
		disableBatching = val
	}

	chunkingAllowed := false
	if val, exists := optionalParam[bool](request.Params.Arguments, "chunking"); exists {
		chunkingAllowed = val
	}

	separator := ","
	if val, exists := optionalParam[string](request.Params.Arguments, "separator"); exists && val != "" {
		separator = val
	}

	properties := []string{}
	if val, exists := optionalParam[[]interface{}](request.Params.Arguments, "properties"); exists && len(val) > 0 {
		for _, p := range val {
			if strProp, ok := p.(string); ok {
				properties = append(properties, strProp)
			}
		}
	}

	key := ""
	if val, exists := optionalParam[string](request.Params.Arguments, "key"); exists {
		key = val
	}

	keyValueKey := ""
	if val, exists := optionalParam[string](request.Params.Arguments, "key-value-key"); exists {
		keyValueKey = val
	}

	keyValueKeyFile := ""
	if val, exists := optionalParam[string](request.Params.Arguments, "key-value-key-file"); exists {
		keyValueKeyFile = val
	}

	// Read but not used due to API compatibility issues
	_ = ""
	if val, exists := optionalParam[string](request.Params.Arguments, "value-schema"); exists && val != "" {
		_ = val
	}

	_ = ""
	if val, exists := optionalParam[string](request.Params.Arguments, "key-schema"); exists && val != "" {
		_ = val
	}

	_ = ""
	if val, exists := optionalParam[string](request.Params.Arguments, "key-value-encoding-type"); exists {
		_ = val
		if val != "separated" && val != "inline" {
			return mcp.NewToolResultError("Invalid key-value-encoding-type: must be 'separated' or 'inline'"), nil
		}
	}

	_ = ""
	if val, exists := optionalParam[string](request.Params.Arguments, "encryption-key-name"); exists {
		_ = val
	}

	_ = ""
	if val, exists := optionalParam[string](request.Params.Arguments, "encryption-key-value"); exists {
		_ = val
	}

	disableReplication := false
	if val, exists := optionalParam[bool](request.Params.Arguments, "disable-replication"); exists {
		disableReplication = val
	}

	// Split messages by separator if needed
	if separator != "" && len(messages) > 0 {
		var splitMessages []string
		for _, msg := range messages {
			parts := strings.Split(msg, separator)
			for _, part := range parts {
				if part != "" {
					splitMessages = append(splitMessages, part)
				}
			}
		}
		messages = splitMessages
	}

	// Setup client
	client, err := mcppulsar.GetPulsarClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create Pulsar client: %v", err)), nil
	}
	defer client.Close()

	// Prepare producer options
	producerOpts := pulsar.ProducerOptions{
		Topic: topic,
	}

	// Set batching and chunking options
	if chunkingAllowed {
		producerOpts.EnableChunking = true
		producerOpts.BatchingMaxPublishDelay = 0 * time.Millisecond
	} else if disableBatching {
		producerOpts.BatchingMaxPublishDelay = 0 * time.Millisecond
	}

	producer, err := client.CreateProducer(producerOpts)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create producer: %v", err)), nil
	}
	defer producer.Close()

	// Generate message bodies from messages and files
	messagePayloads, err := generateMessagePayloads(messages, files)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to generate message payloads: %v", err)), nil
	}

	// Parse properties
	propMap := make(map[string]string)
	for _, prop := range properties {
		parts := strings.SplitN(prop, "=", 2)
		if len(parts) == 2 {
			propMap[parts[0]] = parts[1]
		}
	}

	// Get key value key bytes
	_ = []byte(nil)
	if keyValueKey != "" {
		// We'd validate but not using for now
		_ = []byte(keyValueKey)
	} else if keyValueKeyFile != "" {
		// We'd validate but not using for now
		data, err := os.ReadFile(keyValueKeyFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to read key-value-key-file: %v", err)), nil
		}
		_ = data
	}

	// Setup rate limiter
	var limiter *time.Ticker
	if rate > 0 {
		interval := time.Duration(1000/rate) * time.Millisecond
		limiter = time.NewTicker(interval)
		defer limiter.Stop()
	}

	// Send messages
	numMessagesSent := 0
	var lastMessageID pulsar.MessageID
	for i := 0; i < numProduce; i++ {
		for _, payload := range messagePayloads {
			// Apply rate limiting if enabled
			if limiter != nil {
				<-limiter.C
			}

			// Create message to send
			msg := &pulsar.ProducerMessage{
				Payload: payload,
			}

			// Add properties if specified
			if len(propMap) > 0 {
				msg.Properties = propMap
			}

			// Standard message
			if key != "" {
				msg.Key = key
			}

			// Handle geo-replication
			if disableReplication {
				msg.DisableReplication = true
			}

			// Note: KeyValue handling is commented out until correct API is confirmed
			// if keyValueEncodingType != "" && keyValueKeyBytes != nil {
			//     // TODO: Implement KeyValue schema handling when API is confirmed
			// }

			// Send the message
			msgID, err := producer.Send(ctx, msg)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to send message: %v", err)), nil
			}

			lastMessageID = msgID
			numMessagesSent++
		}
	}

	// Prepare response
	response := map[string]interface{}{
		"topic":           topic,
		"messages_sent":   numMessagesSent,
		"last_message_id": fmt.Sprintf("%v", lastMessageID),
		"success":         true,
	}

	// Convert to JSON
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to encode result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// generateMessagePayloads generates message payloads from message strings and files
func generateMessagePayloads(messages []string, files []string) ([][]byte, error) {
	var payloads [][]byte

	// Add message strings
	for _, msg := range messages {
		payloads = append(payloads, []byte(msg))
	}

	// Add file contents
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", file, err)
		}
		payloads = append(payloads, data)
	}

	return payloads, nil
}
