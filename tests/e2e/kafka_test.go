//go:build e2e

// Copyright 2026 StreamNative
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

package e2e_test

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/mark3labs/mcp-go/client"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
)

func TestKafka(t *testing.T) {
	artifacts := newArtifacts(t)
	bootstrap := startKafka(t, artifacts)
	binary := buildServer(t, artifacts)
	for _, mode := range []string{"stdio", "sse", "http"} {
		t.Run(mode, func(t *testing.T) {
			c := startServer(t, binary, mode, artifacts, []string{
				"--use-external-kafka", "--kafka-bootstrap-servers", bootstrap,
			})
			initialize(t, c, mode)
			checkTools(t, c, "kafka_admin_topics_write", "kafka_admin_topics_read", "kafka_client_produce", "kafka_client_consume")
			checkKafkaLifecycle(t, c, bootstrap)
		})
	}
}

func startKafka(t *testing.T, artifacts string) string {
	t.Helper()
	requireLocalDocker(t)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()
	// Defer the stock Kafka entrypoint until Docker has allocated the external
	// port. Both bootstrap and advertised listeners must reach this container,
	// including on Docker Desktop; never reserve a fixed host port.
	broker, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "apache/kafka:3.9.1",
			ExposedPorts: []string{"9093/tcp"},
			Entrypoint:   []string{"/bin/bash", "-c"},
			Cmd:          []string{"while [ ! -f /tmp/e2e-start.sh ]; do sleep 0.1; done; exec /bin/bash /tmp/e2e-start.sh"},
			Env: map[string]string{
				"CLUSTER_ID":                                     "MkU3OEVBNTcwNTJENDM2Qk",
				"KAFKA_NODE_ID":                                  "1",
				"KAFKA_PROCESS_ROLES":                            "broker,controller",
				"KAFKA_LISTENERS":                                "INTERNAL://:9092,EXTERNAL://:9093,CONTROLLER://:9094",
				"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP":           "INTERNAL:PLAINTEXT,EXTERNAL:PLAINTEXT,CONTROLLER:PLAINTEXT",
				"KAFKA_INTER_BROKER_LISTENER_NAME":               "INTERNAL",
				"KAFKA_CONTROLLER_LISTENER_NAMES":                "CONTROLLER",
				"KAFKA_CONTROLLER_QUORUM_VOTERS":                 "1@127.0.0.1:9094",
				"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR":         "1",
				"KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR": "1",
				"KAFKA_TRANSACTION_STATE_LOG_MIN_ISR":            "1",
				"KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS":         "0",
				"KAFKA_AUTO_CREATE_TOPICS_ENABLE":                "false",
				"KAFKA_HEAP_OPTS":                                "-Xms256m -Xmx512m",
			},
			HostConfigModifier: func(config *container.HostConfig) {
				config.PortBindings = nat.PortMap{"9093/tcp": {{HostIP: "127.0.0.1"}}}
			},
		},
	})
	cleanupContainer(t, broker, artifacts, "kafka")
	require.NoError(t, err, "create disposable Kafka container")
	require.NoError(t, broker.Start(ctx))
	host, err := broker.Host(ctx)
	require.NoError(t, err)
	require.True(t, host == "localhost" || net.ParseIP(host).IsLoopback(), "only local container endpoints are allowed: %s", host)
	port, err := broker.MappedPort(ctx, "9093/tcp")
	require.NoError(t, err)
	bootstrap := net.JoinHostPort("127.0.0.1", port.Port())
	script := fmt.Sprintf("#!/bin/bash\nexport KAFKA_ADVERTISED_LISTENERS=INTERNAL://127.0.0.1:9092,EXTERNAL://%s\nexec /etc/kafka/docker/run\n", bootstrap)
	// Docker copies as root; the image runs Kafka as a non-root user. This
	// script contains only local listener addresses, no credentials.
	require.NoError(t, broker.CopyToContainer(ctx, []byte(script), "/tmp/e2e-start.sh", 0644))
	native := newKafkaClient(t, bootstrap)
	require.Eventually(t, func() bool {
		probe, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		_, err := kadm.NewClient(native).ListTopics(probe)
		return err == nil
	}, 2*time.Minute, 200*time.Millisecond, "Kafka API readiness; see kafka.log")
	return bootstrap
}

func newKafkaClient(t *testing.T, bootstrap string, opts ...kgo.Opt) *kgo.Client {
	t.Helper()
	opts = append(opts, kgo.SeedBrokers(bootstrap), kgo.DialTimeout(3*time.Second), kgo.RequestTimeoutOverhead(5*time.Second))
	c, err := kgo.NewClient(opts...)
	require.NoError(t, err)
	t.Cleanup(c.Close)
	return c
}

func kafkaTopics(t *testing.T, c *kgo.Client) kadm.TopicDetails {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	topics, err := kadm.NewClient(c).ListTopics(ctx)
	require.NoError(t, err)
	return topics
}

func checkKafkaLifecycle(t *testing.T, c *client.Client, bootstrap string) {
	t.Helper()
	topic := "e2e-" + strings.ToLower(rand.Text())
	native := newKafkaClient(t, bootstrap)
	require.NotContains(t, kafkaTopics(t, native), topic)
	args := map[string]any{"resource": "topic", "operation": "create", "name": topic, "partitions": 1, "replicationFactor": 1}
	created := callTool(t, c, "kafka_admin_topics_write", args, false)
	checkKafkaTopicReply(t, created, topic, nil)
	topics := kafkaTopics(t, native)
	require.Contains(t, topics, topic)
	require.NoError(t, topics[topic].Err)
	require.Len(t, topics[topic].Partitions, 1)

	value, key := "payload-"+rand.Text(), "key-"+rand.Text()
	produced := callTool(t, c, "kafka_client_produce", map[string]any{
		"topic": topic, "value": value, "key": key, "headers": []string{"source=e2e"}, "sync": true,
	}, false)
	var result struct{ Status, Topic, Key string }
	require.NoError(t, json.Unmarshal([]byte(produced), &result))
	require.Equal(t, "success", result.Status)
	require.Equal(t, topic, result.Topic)
	require.Equal(t, key, result.Key)

	// Independent Kafka protocol fetch: no MCP handler/client/session reuse.
	reader := newKafkaClient(t, bootstrap, kgo.ConsumePartitions(map[string]map[int32]kgo.Offset{
		topic: {0: kgo.NewOffset().AtStart()},
	}))
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	fetches := reader.PollRecords(ctx, 1)
	require.Empty(t, fetches.Errors())
	records := fetches.Records()
	require.Len(t, records, 1)
	require.Equal(t, topic, records[0].Topic)
	require.Equal(t, value, string(records[0].Value))
	require.Equal(t, key, string(records[0].Key))
	require.Equal(t, []kgo.RecordHeader{{Key: "source", Value: []byte("e2e")}}, records[0].Headers)
	reader.Close() // Stop fetching before deletion; Close is idempotent.

	consumed := callTool(t, c, "kafka_client_consume", map[string]any{
		"topic": topic, "offset": "atstart", "max-messages": 1, "timeout": 10,
	}, false)
	var payloads [][]byte // The tool JSON-encodes raw bytes as base64 strings.
	require.NoError(t, json.Unmarshal([]byte(consumed), &payloads))
	require.Equal(t, [][]byte{[]byte(value)}, payloads)

	// Kafka topic admin tools currently return per-topic errors inside the
	// JSON result (not IsError). Assert the actual broker code, not success.
	args["partitions"] = 2
	conflict := callTool(t, c, "kafka_admin_topics_write", args, false)
	checkKafkaTopicReply(t, conflict, topic, kerr.TopicAlreadyExists)
	topics = kafkaTopics(t, native)
	require.Contains(t, topics, topic)
	require.NoError(t, topics[topic].Err)
	require.Len(t, topics[topic].Partitions, 1, "rejected create must not change partitions")

	deleted := callTool(t, c, "kafka_admin_topics_write", map[string]any{"resource": "topic", "operation": "delete", "name": topic}, false)
	checkKafkaTopicReply(t, deleted, topic, nil)
	require.Eventually(t, func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		topics, err := kadm.NewClient(native).ListTopics(ctx)
		_, exists := topics[topic]
		return err == nil && !exists
	}, requestTimeout, 100*time.Millisecond, "native Kafka API must confirm deletion")
}

func checkKafkaTopicReply(t *testing.T, raw, topic string, expected *kerr.Error) {
	t.Helper()
	var results map[string]struct {
		Topic string
		Err   *kerr.Error
	}
	require.NoError(t, json.Unmarshal([]byte(raw), &results))
	require.Contains(t, results, topic)
	require.Equal(t, topic, results[topic].Topic)
	if expected == nil {
		require.Nil(t, results[topic].Err, "Kafka per-topic error: %s", raw)
	} else {
		require.NotNil(t, results[topic].Err)
		require.Equal(t, expected.Code, results[topic].Err.Code)
		require.Equal(t, expected.Message, results[topic].Err.Message)
	}
}
