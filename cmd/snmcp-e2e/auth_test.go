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
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}
	if !getenvBool("E2E_USE_TESTCONTAINERS", false) {
		t.Skip("set E2E_USE_TESTCONTAINERS=true to run e2e tests")
	}

	tokens := loadAuthTokens(t)
	secretKeyPath := loadSecretKeyPath(t)
	cfg := loadTestcontainersConfig()

	overallTimeout := cfg.StartupTimeout + 2*time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), overallTimeout)
	defer cancel()

	env, err := startPulsarContainerWithAuth(ctx, cfg, tokens.AdminToken, secretKeyPath)
	require.NoError(t, err)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cleanupCancel()
		_ = env.Terminate(cleanupCtx)
	})

	snmcpBaseURL, stopServer, err := startSNMCPServer(t, env.PulsarWebServiceURL, env.PulsarBrokerURL)
	require.NoError(t, err)
	t.Cleanup(stopServer)

	sseURL := snmcpBaseURL + "/sse"
	require.NoError(t, expectUnauthorized(ctx, sseURL, "", false))
	require.NoError(t, expectUnauthorized(ctx, sseURL, "invalid-token", false))

	adminClient, err := newAuthedClient(ctx, sseURL, tokens.AdminToken, "snmcp-e2e-auth-admin")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = adminClient.Close()
	})

	testClient, err := newAuthedClient(ctx, sseURL, tokens.TestUserToken, "snmcp-e2e-auth-test-user")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = testClient.Close()
	})

	clusters, err := listClusters(ctx, adminClient)
	require.NoError(t, err)
	require.NotEmpty(t, clusters)
	cluster := clusters[0]

	suffix := time.Now().UnixNano()
	tenant := fmt.Sprintf("auth-e2e-%d", suffix)

	result, err := callTool(ctx, adminClient, "pulsar_admin_tenant", map[string]any{
		"resource":        "tenant",
		"operation":       "create",
		"tenant":          tenant,
		"adminRoles":      []string{"admin"},
		"allowedClusters": []string{cluster},
	})
	require.NoError(t, requireToolOK(result, err, "pulsar_admin_tenant create"))

	result, err = callTool(ctx, testClient, "pulsar_admin_tenant", map[string]any{
		"resource":        "tenant",
		"operation":       "create",
		"tenant":          tenant + "-unauthorized",
		"adminRoles":      []string{"test-user"},
		"allowedClusters": []string{cluster},
	})
	require.NoError(t, requireToolError(result, err, "pulsar_admin_tenant unauthorized create"))
}
