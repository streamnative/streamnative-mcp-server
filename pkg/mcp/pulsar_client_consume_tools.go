package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	mcppulsar "github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// PulsarClientAddConsumerTools adds Pulsar client consumer tools to the MCP server
func PulsarClientAddConsumerTools(s *server.MCPServer, _ bool, features []string) {
	if !slices.Contains(features, string(FeaturePulsarClient)) && !slices.Contains(features, string(FeatureAll)) && !slices.Contains(features, string(FeatureAllPulsar)) {
		return
	}

	// Main consume tool
	consumeTool := mcp.NewTool("pulsar_client_consume",
		mcp.WithDescription("Consume messages from a Pulsar topic. This tool allows you to consume messages "+
			"from a specified Pulsar topic with various options to control the subscription behavior, "+
			"message processing, and display format."),
		mcp.WithString("topic", mcp.Required(),
			mcp.Description("Topic to consume from"),
		),
		mcp.WithString("subscription-name", mcp.Required(),
			mcp.Description("Subscription name"),
		),
		mcp.WithString("subscription-type",
			mcp.Description("Subscription type: exclusive, shared, failover, key_shared (default: exclusive)"),
		),
		mcp.WithString("subscription-mode",
			mcp.Description("Subscription mode: durable, non-durable (default: durable)"),
		),
		mcp.WithString("initial-position",
			mcp.Description("Initial position: latest (consume from the latest message)"+
				", earliest (consume from the earliest message) (default: latest)"),
		),
		mcp.WithNumber("num-messages",
			mcp.Description("Number of messages to consume (0 for unlimited, default: 0)"),
		),
		mcp.WithNumber("timeout",
			mcp.Description("Timeout for consuming messages in seconds (default: 30)"),
		),
		mcp.WithBoolean("regex",
			mcp.Description("Treat the topic as a regex pattern (default: false)"),
		),
		mcp.WithString("schema",
			mcp.Description("Schema type: string, json, avro, protobuf"),
		),
		mcp.WithNumber("max-redeliver-count",
			mcp.Description("Maximum redelivery count for dead letter queue (0 to disable, default: 0)"),
		),
		mcp.WithString("dlq-topic",
			mcp.Description("Dead letter queue topic"),
		),
		mcp.WithBoolean("show-properties",
			mcp.Description("Show message properties (default: false)"),
		),
		mcp.WithBoolean("hide-payload",
			mcp.Description("Hide message payload (default: false)"),
		),
		mcp.WithBoolean("read-compacted",
			mcp.Description("Read compacted topic (default: false)"),
		),
		mcp.WithNumber("receiver-queue-size",
			mcp.Description("Size of the consumer receive queue (default: 1000)"),
		),
	)
	s.AddTool(consumeTool, handleClientConsume)
}

func handleClientConsume(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract required parameters with validation
	topic, err := requiredParam[string](request.Params.Arguments, "topic")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topic: %v", err)), nil
	}

	subscriptionName, err := requiredParam[string](request.Params.Arguments, "subscription-name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get subscription name: %v", err)), nil
	}

	// Set default values and extract optional parameters
	subscriptionType := "exclusive"
	if val, exists := optionalParam[string](request.Params.Arguments, "subscription-type"); exists && val != "" {
		subscriptionType = val
	}

	subscriptionMode := "durable"
	if val, exists := optionalParam[string](request.Params.Arguments, "subscription-mode"); exists && val != "" {
		subscriptionMode = val
	}

	initialPosition := "latest"
	if val, exists := optionalParam[string](request.Params.Arguments, "initial-position"); exists && val != "" {
		initialPosition = val
	}

	numMessages := 0
	if val, exists := optionalParam[float64](request.Params.Arguments, "num-messages"); exists {
		numMessages = int(val)
	}

	timeout := 30
	if val, exists := optionalParam[float64](request.Params.Arguments, "timeout"); exists {
		timeout = int(val)
	}

	regex := false
	if val, exists := optionalParam[bool](request.Params.Arguments, "regex"); exists {
		regex = val
	}

	var maxRedeliverCount int32
	if val, exists := optionalParam[float64](request.Params.Arguments, "max-redeliver-count"); exists {
		maxRedeliverCount = int32(val)
	}

	dlqTopic := ""
	if val, exists := optionalParam[string](request.Params.Arguments, "dlq-topic"); exists {
		dlqTopic = val
	}

	showProperties := false
	if val, exists := optionalParam[bool](request.Params.Arguments, "show-properties"); exists {
		showProperties = val
	}

	hidePayload := false
	if val, exists := optionalParam[bool](request.Params.Arguments, "hide-payload"); exists {
		hidePayload = val
	}

	readCompacted := false
	if val, exists := optionalParam[bool](request.Params.Arguments, "read-compacted"); exists {
		readCompacted = val
	}

	receiverQueueSize := 1000
	if val, exists := optionalParam[float64](request.Params.Arguments, "receiver-queue-size"); exists {
		receiverQueueSize = int(val)
	}

	// Setup client
	client, err := mcppulsar.GetPulsarClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create Pulsar client: %v", err)), nil
	}
	defer client.Close()

	// Prepare consumer options
	consumerOpts := pulsar.ConsumerOptions{
		Name:              "snmcp-consumer",
		ReceiverQueueSize: receiverQueueSize,
		ReadCompacted:     readCompacted,
		SubscriptionName:  subscriptionName,
	}

	// Set topic or topics pattern
	if regex {
		consumerOpts.TopicsPattern = topic
	} else {
		consumerOpts.Topic = topic
	}

	// Set subscription type
	switch strings.ToLower(subscriptionType) {
	case "exclusive":
		consumerOpts.Type = pulsar.Exclusive
	case "shared":
		consumerOpts.Type = pulsar.Shared
	case "failover":
		consumerOpts.Type = pulsar.Failover
	case "key_shared":
		consumerOpts.Type = pulsar.KeyShared
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Invalid subscription type: %s", subscriptionType)), nil
	}

	// Set subscription mode
	switch strings.ToLower(subscriptionMode) {
	case "durable":
		consumerOpts.SubscriptionMode = pulsar.Durable
	case "non-durable":
		consumerOpts.SubscriptionMode = pulsar.NonDurable
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Invalid subscription mode: %s", subscriptionMode)), nil
	}

	// Set initial position
	switch strings.ToLower(initialPosition) {
	case "latest":
		consumerOpts.SubscriptionInitialPosition = pulsar.SubscriptionPositionLatest
	case "earliest":
		consumerOpts.SubscriptionInitialPosition = pulsar.SubscriptionPositionEarliest
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Invalid initial position: %s", initialPosition)), nil
	}

	// Set DLQ policy if specified
	if maxRedeliverCount > 0 {
		if dlqTopic == "" {
			return mcp.NewToolResultError("DLQ topic is required when max-redeliver-count is specified"), nil
		}
		consumerOpts.DLQ = &pulsar.DLQPolicy{
			//nolint:gosec
			MaxDeliveries:    uint32(maxRedeliverCount),
			DeadLetterTopic:  dlqTopic,
			RetryLetterTopic: fmt.Sprintf("%s-retry", dlqTopic),
		}
	}

	// Create consumer
	consumer, err := client.Subscribe(consumerOpts)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create consumer: %v", err)), nil
	}
	defer consumer.Close()

	// Set up timeout context
	timeoutDuration := time.Duration(timeout) * time.Second
	consumeCtx, cancelConsume := context.WithTimeout(ctx, timeoutDuration)
	defer cancelConsume()

	// Container for messages
	type MessageData struct {
		ID           string            `json:"id"`
		PublishTime  string            `json:"publish_time"`
		Properties   map[string]string `json:"properties,omitempty"`
		Key          string            `json:"key,omitempty"`
		Data         string            `json:"data,omitempty"`
		MessageCount int               `json:"message_count"`
	}

	messages := []MessageData{}
	messageCount := 0

	// Consume messages
	for {
		// Check if we've consumed the requested number of messages
		if numMessages > 0 && messageCount >= numMessages {
			break
		}

		// Receive message with timeout
		msg, err := consumer.Receive(consumeCtx)
		if err != nil {
			if err == context.DeadlineExceeded {
				break
			}
			if err == context.Canceled {
				break
			}
			return mcp.NewToolResultError(fmt.Sprintf("Error receiving message: %v", err)), nil
		}

		// Process the message
		messageCount++

		// Create message data
		messageData := MessageData{
			ID:           msg.ID().String(),
			PublishTime:  msg.PublishTime().Format(time.RFC3339),
			MessageCount: messageCount,
		}

		// Add properties if requested
		if showProperties {
			messageData.Properties = msg.Properties()
		}

		// Add key if present
		if msg.Key() != "" {
			messageData.Key = msg.Key()
		}

		// Add payload unless hidden
		if !hidePayload {
			messageData.Data = string(msg.Payload())
		}

		messages = append(messages, messageData)

		// Acknowledge the message
		_ = consumer.Ack(msg)
	}

	// Prepare response
	response := map[string]interface{}{
		"topic":             topic,
		"subscription_name": subscriptionName,
		"messages_consumed": messageCount,
		"messages":          messages,
	}

	// Convert to JSON
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to encode result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}
