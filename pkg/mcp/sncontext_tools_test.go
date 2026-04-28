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

package mcp

import (
	"context"
	"testing"

	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	pulsarsession "github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
	"github.com/stretchr/testify/require"
)

func TestHandleResetContextClearsSNCloudClusterSessions(t *testing.T) {
	snSession, err := config.NewSNCloudSession(config.SNCloudContext{
		JWTToken:       "token",
		APIURL:         "https://api.example.com",
		LogAPIURL:      "https://logs.example.com",
		Organization:   "org",
		PulsarInstance: "instance-a",
		PulsarCluster:  "cluster-a",
	})
	require.NoError(t, err)

	pulsarSession := &pulsarsession.Session{
		Ctx: pulsarsession.PulsarContext{
			ServiceURL:    "pulsar://pulsar.example.com:6650",
			WebServiceURL: "https://pulsar.example.com",
			Token:         "token",
		},
		PulsarCtlConfig: &cmdutils.ClusterConfig{WebServiceURL: "https://pulsar.example.com"},
	}
	kafkaSession := &kafka.Session{
		Ctx: kafka.KafkaContext{
			BootstrapServers:  "kafka.example.com:9093",
			SchemaRegistryURL: "https://kafka.example.com/kafka",
			ConnectURL:        "https://api.example.com/admin/kafkaconnect/",
		},
	}

	ctx := context.Background()
	ctx = WithSNCloudSession(ctx, snSession)
	ctx = WithPulsarSession(ctx, pulsarSession)
	ctx = WithKafkaSession(ctx, kafkaSession)

	result, err := handleResetContext(ctx, mcpgo.CallToolRequest{})
	require.NoError(t, err)
	require.False(t, result.IsError)

	require.Empty(t, snSession.Ctx.PulsarInstance)
	require.Empty(t, snSession.Ctx.PulsarCluster)
	require.Equal(t, pulsarsession.PulsarContext{}, pulsarSession.Ctx)
	require.Equal(t, kafka.KafkaContext{}, kafkaSession.Ctx)

	_, err = pulsarSession.GetAdminClient()
	require.EqualError(t, err, "err: ContextNotSetErr: Please set the cluster context first")

	_, err = kafkaSession.GetAdminClient()
	require.EqualError(t, err, "err: ContextNotSetErr: Please set the cluster context first")
}
