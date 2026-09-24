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

package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

type testRuntimeResolver struct {
	value      config.ResolvedRuntime
	failure    error
	selections [][2]string
}

func (r *testRuntimeResolver) ResolveCluster(_ context.Context, instance, cluster string) (config.ResolvedRuntime, error) {
	r.selections = append(r.selections, [2]string{instance, cluster})
	return r.value, r.failure
}

func TestInjectedResolverOwnsAuthenticationAndResolution(t *testing.T) {
	//nolint:gosec // Synthetic token; no real issuer or credential is used by this test.
	source := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "opaque-runtime-token"})
	resolver := &testRuntimeResolver{value: config.ResolvedRuntime{Pulsar: pulsar.PulsarContext{
		ServiceURL: "pulsar://localhost:6650", WebServiceURL: "http://localhost:8080", TokenSource: source,
	}}}
	o := config.NewConfigOptions()
	o.CloudProvider = &config.CloudProvider{Identity: "subject", TokenSource: source, Resolver: resolver}
	o.Server, o.Organization = "http://127.0.0.1:1", "org"
	session, err := config.NewSNCloudSessionFromOptions(o)
	require.NoError(t, err)
	defer func() { _ = session.Close() }()
	ctx := WithSNCloudSession(context.Background(), session)
	require.NoError(t, SetContext(ctx, o, "one", "cluster-one"))
	previous := GetPulsarSession(ctx)
	require.NotNil(t, previous.Client)
	resolver.failure = fmt.Errorf("selection denied")
	require.ErrorContains(t, SetContext(ctx, o, "two", "denied"), "selection denied")
	require.Same(t, previous, GetPulsarSession(ctx))
	resolver.failure = nil
	resolver.value.Kafka = &kafka.KafkaContext{BootstrapServers: "localhost:9092", ConnectURL: "http://invalid/%"}
	require.Error(t, SetContext(ctx, o, "two", "invalid"))
	require.Same(t, previous, GetPulsarSession(ctx))
	resolver.value.Kafka = nil
	require.NoError(t, SetContext(ctx, o, "two", "cluster-two"))
	require.Nil(t, previous.Client, "successful publication closes the retired generation")
	require.Equal(t, [][2]string{{"one", "cluster-one"}, {"two", "denied"}, {"two", "invalid"}, {"two", "cluster-two"}}, resolver.selections)
	require.NoError(t, ResetContext(ctx))
	_, err = GetPulsarSession(ctx).GetPulsarClient()
	require.ErrorContains(t, err, "ContextNotSetErr")
	require.Nil(t, o.Store, "the runtime must not require caller credentials or a grant store")
}

func newInjectedRuntime(t *testing.T) (*testRuntimeResolver, context.Context, *config.Options) {
	t.Helper()
	source := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "opaque"})
	resolver := &testRuntimeResolver{value: config.ResolvedRuntime{
		Pulsar: pulsar.PulsarContext{ServiceURL: "pulsar://localhost:6650", WebServiceURL: "http://localhost:8080", TokenSource: source},
		Kafka:  &kafka.KafkaContext{BootstrapServers: "localhost:9092"},
	}}
	options := config.NewConfigOptions()
	options.Server, options.Organization = "http://127.0.0.1:1", "org"
	options.CloudProvider = &config.CloudProvider{Identity: "test-sa", TokenSource: source, Resolver: resolver}
	session, err := config.NewSNCloudSessionFromOptions(options)
	require.NoError(t, err)
	t.Cleanup(func() { _ = session.Close() })
	return resolver, WithSNCloudSession(context.Background(), session), options
}

func TestInjectedPulsarOnlyBindingDropsKafkaClients(t *testing.T) {
	resolver, ctx, options := newInjectedRuntime(t)
	require.NoError(t, SetContext(ctx, options, "instance", "cluster"))
	previous := GetKafkaSession(ctx)
	require.NotNil(t, previous.Client)
	resolver.value.Kafka = nil
	require.NoError(t, SetContext(ctx, options, "instance", "another"))
	_, err := GetKafkaSession(ctx).GetClient()
	require.ErrorContains(t, err, "ContextNotSetErr")
	require.Nil(t, previous.Client)
}

func TestResetWaitsForInFlightMcpRequest(t *testing.T) {
	_, ctx, options := newInjectedRuntime(t)
	require.NoError(t, SetContext(ctx, options, "instance-a", "cluster-a"))
	old := GetPulsarSession(ctx)
	entered, release, completed := make(chan struct{}), make(chan struct{}), make(chan struct{})
	srv := NewServer("test", "test", logrus.New())
	srv.MCPServer.AddTool(mcpgo.NewTool("hold-binding"), func(callCtx context.Context, _ mcpgo.CallToolRequest) (*mcpgo.CallToolResult, error) {
		close(entered)
		<-release
		if GetPulsarSession(callCtx) != old || old.Client == nil {
			return nil, fmt.Errorf("in-flight binding was replaced")
		}
		return mcpgo.NewToolResultText("preserved"), nil
	})
	response := make(chan mcpgo.JSONRPCMessage, 1)
	go func() {
		response <- srv.MCPServer.HandleMessage(ctx, json.RawMessage(`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"hold-binding","arguments":{}}}`))
	}()
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("request did not enter handler")
	}
	go func() { _ = ResetContext(ctx); close(completed) }()
	select {
	case <-completed:
		close(release)
		<-response
		t.Fatal("reset closed a binding still used by a request")
	case <-time.After(25 * time.Millisecond):
	}
	close(release)
	result := <-response
	encoded, err := json.Marshal(result)
	require.NoError(t, err)
	require.Contains(t, string(encoded), "preserved")
	select {
	case <-completed:
	case <-time.After(time.Second):
		t.Fatal("reset did not finish after request released binding")
	}
	require.Nil(t, old.Client, "retired clients must be closed")
}

func TestServerCloseReleasesSelectedRuntimeClients(t *testing.T) {
	_, ctx, options := newInjectedRuntime(t)
	require.NoError(t, SetContext(ctx, options, "instance-a", "cluster-a"))
	ps, ks := GetPulsarSession(ctx), GetKafkaSession(ctx)
	srv := NewServer("test", "test", logrus.New())
	srv.SNCloudSession = GetSNCloudSession(ctx)
	srv.Close()
	require.Nil(t, ps.Client)
	require.Nil(t, ks.Client)
}
