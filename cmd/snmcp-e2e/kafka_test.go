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

package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/require"
)

func TestKafkaE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}
	if !getenvBool("E2E_USE_TESTCONTAINERS", false) {
		t.Skip("set E2E_USE_TESTCONTAINERS=true to run e2e tests")
	}

	cfg := loadTestcontainersConfig()
	overallTimeout := cfg.StartupTimeout + 3*time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), overallTimeout)
	defer cancel()

	kafkaContainer, kafkaBrokers, err := startKafkaContainer(ctx, cfg)
	require.NoError(t, err)
	env := &testcontainersEnv{
		Kafka:        kafkaContainer,
		KafkaBrokers: kafkaBrokers,
	}
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cleanupCancel()
		_ = env.Terminate(cleanupCtx)
	})

	snmcpBaseURL, stopServer, err := startSNMCPServerWithKafka(t, kafkaBrokers)
	require.NoError(t, err)
	t.Cleanup(stopServer)

	sseURL := snmcpBaseURL + "/sse"
	client, err := newAuthedClient(ctx, sseURL, "", "snmcp-e2e-kafka")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = client.Close()
	})

	suffix := time.Now().UnixNano()
	topic := fmt.Sprintf("e2e-topic-%d", suffix)
	group := fmt.Sprintf("e2e-group-%d", suffix)
	message := fmt.Sprintf("message-%d", suffix)

	result, err := callTool(ctx, client, "kafka_admin_topics", map[string]any{
		"resource":          "topic",
		"operation":         "create",
		"name":              topic,
		"partitions":        float64(1),
		"replicationFactor": float64(1),
	})
	require.NoError(t, requireToolOK(result, err, "kafka_admin_topics create"))

	result, err = callTool(ctx, client, "kafka_admin_topics", map[string]any{
		"resource":  "topics",
		"operation": "list",
	})
	require.NoError(t, requireToolOK(result, err, "kafka_admin_topics list"))
	topics, err := parseKafkaTopicNames(firstText(result))
	require.NoError(t, err)
	require.Contains(t, topics, topic)

	result, err = callTool(ctx, client, "kafka_client_produce", map[string]any{
		"topic": topic,
		"value": message,
	})
	require.NoError(t, requireToolOK(result, err, "kafka_client_produce"))

	result, err = callTool(ctx, client, "kafka_client_consume", map[string]any{
		"topic":        topic,
		"group":        group,
		"offset":       "atstart",
		"max-messages": float64(1),
		"timeout":      float64(20),
	})
	require.NoError(t, requireToolOK(result, err, "kafka_client_consume"))
	messages, err := decodeKafkaConsumeValues(firstText(result))
	require.NoError(t, err)
	require.Contains(t, messages, message)

	require.NoError(t, waitForKafkaGroup(ctx, client, group))

	result, err = callTool(ctx, client, "kafka_admin_groups", map[string]any{
		"resource":  "group",
		"operation": "describe",
		"group":     group,
	})
	require.NoError(t, requireToolOK(result, err, "kafka_admin_groups describe"))

	result, err = callTool(ctx, client, "kafka_admin_groups", map[string]any{
		"resource":  "group",
		"operation": "offsets",
		"group":     group,
	})
	require.NoError(t, requireToolOK(result, err, "kafka_admin_groups offsets"))

	result, err = callTool(ctx, client, "kafka_admin_partitions", map[string]any{
		"resource":  "partition",
		"operation": "update",
		"topic":     topic,
		"new-total": float64(2),
	})
	require.NoError(t, requireToolOK(result, err, "kafka_admin_partitions update"))
	require.NoError(t, waitForKafkaPartitions(ctx, client, topic, 2))
}

func startSNMCPServerWithKafka(t *testing.T, kafkaBrokers string) (string, func(), error) {
	t.Helper()

	if strings.TrimSpace(kafkaBrokers) == "" {
		return "", nil, errors.New("kafka brokers are required")
	}

	repoRoot := findRepoRoot(t)
	binaryPath := buildSNMCPBinary(t, repoRoot)

	addr := reserveLocalAddr(t)
	baseURL := fmt.Sprintf("http://%s/mcp", addr)

	//nolint:gosec // test binary path and arguments are controlled
	cmd := exec.Command(binaryPath,
		"sse",
		"--http-addr", addr,
		"--http-path", "/mcp",
		"--use-external-kafka",
		"--kafka-bootstrap-servers", kafkaBrokers,
	)
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(), "SNMCP_CONFIG_DIR="+t.TempDir())

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Start(); err != nil {
		return "", nil, fmt.Errorf("start snmcp: %w", err)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- cmd.Wait()
	}()

	readyCtx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	if err := waitForHTTPStatus(readyCtx, baseURL+"/healthz", 200, errCh); err != nil {
		_ = cmd.Process.Signal(os.Interrupt)
		<-errCh
		return "", nil, fmt.Errorf("snmcp not ready: %w\n%s", err, strings.TrimSpace(output.String()))
	}

	cleanup := func() {
		if cmd.Process == nil {
			return
		}
		_ = cmd.Process.Signal(os.Interrupt)
		select {
		case <-errCh:
		case <-time.After(10 * time.Second):
			_ = cmd.Process.Kill()
			<-errCh
		}
	}

	return baseURL, cleanup, nil
}

func parseKafkaTopicNames(raw string) ([]string, error) {
	if raw == "" {
		return nil, errors.New("empty topics response")
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return nil, fmt.Errorf("parse topics response: %w", err)
	}
	names := make([]string, 0, len(payload))
	for name := range payload {
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}

func parseKafkaGroupNames(raw string) ([]string, error) {
	if raw == "" {
		return nil, errors.New("empty groups response")
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return nil, fmt.Errorf("parse groups response: %w", err)
	}
	names := make([]string, 0, len(payload))
	for name := range payload {
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}

func decodeKafkaConsumeValues(raw string) ([]string, error) {
	if raw == "" {
		return nil, errors.New("empty consume response")
	}
	var payload []any
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return nil, fmt.Errorf("parse consume response: %w", err)
	}
	values := make([]string, 0, len(payload))
	for _, entry := range payload {
		switch value := entry.(type) {
		case string:
			decoded, err := base64.StdEncoding.DecodeString(value)
			if err != nil {
				values = append(values, value)
				continue
			}
			values = append(values, string(decoded))
		default:
			encoded, err := json.Marshal(value)
			if err != nil {
				values = append(values, fmt.Sprintf("%v", value))
				continue
			}
			values = append(values, string(encoded))
		}
	}
	return values, nil
}

func waitForKafkaGroup(ctx context.Context, client *mcp.ClientSession, group string) error {
	waitCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		result, err := callTool(waitCtx, client, "kafka_admin_groups", map[string]any{
			"resource":  "groups",
			"operation": "list",
		})
		if err == nil && result != nil && !result.IsError {
			names, parseErr := parseKafkaGroupNames(firstText(result))
			if parseErr == nil && slices.Contains(names, group) {
				return nil
			}
		}

		select {
		case <-waitCtx.Done():
			return fmt.Errorf("timeout waiting for kafka group %s", group)
		case <-ticker.C:
		}
	}
}

func waitForKafkaPartitions(ctx context.Context, client *mcp.ClientSession, topic string, expected int) error {
	waitCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		result, err := callTool(waitCtx, client, "kafka_admin_topics", map[string]any{
			"resource":  "topic",
			"operation": "metadata",
			"name":      topic,
		})
		if err == nil && result != nil && !result.IsError {
			count, parseErr := extractKafkaPartitionCount(firstText(result), topic)
			if parseErr == nil && count >= expected {
				return nil
			}
		}

		select {
		case <-waitCtx.Done():
			return fmt.Errorf("timeout waiting for kafka partitions on %s", topic)
		case <-ticker.C:
		}
	}
}

func extractKafkaPartitionCount(raw, topic string) (int, error) {
	if raw == "" {
		return 0, errors.New("empty metadata response")
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return 0, fmt.Errorf("parse metadata response: %w", err)
	}
	topicsValue, ok := payload["Topics"]
	if !ok {
		topicsValue = payload["topics"]
	}
	topics, ok := topicsValue.(map[string]any)
	if !ok {
		return 0, errors.New("unexpected topics payload")
	}
	topicValue, ok := topics[topic]
	if !ok {
		return 0, fmt.Errorf("topic %s not found in metadata", topic)
	}
	topicMap, ok := topicValue.(map[string]any)
	if !ok {
		return 0, errors.New("unexpected topic metadata payload")
	}
	partitionsValue, ok := topicMap["Partitions"]
	if !ok {
		partitionsValue = topicMap["partitions"]
	}
	partitions, ok := partitionsValue.(map[string]any)
	if !ok {
		return 0, errors.New("unexpected partitions payload")
	}
	return len(partitions), nil
}
