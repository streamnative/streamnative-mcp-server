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
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	defaultAdminRole   = "admin"
	defaultTestUser    = "test-user"
	authSecretFilePath = "/pulsarctl/test/auth/token/secret.key" //nolint:gosec // static test container path
)

type authTokens struct {
	AdminToken    string
	TestUserToken string
}

func TestPulsarAdminE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}
	if !getenvBool("E2E_USE_TESTCONTAINERS", false) {
		t.Skip("set E2E_USE_TESTCONTAINERS=true to run e2e tests")
	}

	tokens := loadAuthTokens(t)
	secretKeyPath := loadSecretKeyPath(t)
	cfg := loadTestcontainersConfig()

	overallTimeout := cfg.StartupTimeout + 3*time.Minute
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

	adminClient, err := newAuthedClient(ctx, sseURL, tokens.AdminToken, "snmcp-e2e-admin")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = adminClient.Close()
	})

	testClient, err := newAuthedClient(ctx, sseURL, tokens.TestUserToken, "snmcp-e2e-test-user")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = testClient.Close()
	})

	clusters, err := listClusters(ctx, adminClient)
	require.NoError(t, err)
	require.NotEmpty(t, clusters)

	suffix := time.Now().UnixNano()
	tenant := fmt.Sprintf("e2e-%d", suffix)
	namespace := fmt.Sprintf("%s/ns-%d", tenant, suffix)
	topic := fmt.Sprintf("persistent://%s/topic-%d", namespace, suffix)

	result, err := callTool(ctx, adminClient, "pulsar_admin_tenant", map[string]any{
		"resource":        "tenant",
		"operation":       "create",
		"tenant":          tenant,
		"adminRoles":      []string{defaultAdminRole},
		"allowedClusters": []string{clusters[0]},
	})
	require.NoError(t, requireToolOK(result, err, "pulsar_admin_tenant create"))

	result, err = callTool(ctx, adminClient, "pulsar_admin_namespace", map[string]any{
		"operation": "create",
		"namespace": namespace,
		"clusters":  []string{clusters[0]},
	})
	require.NoError(t, requireToolOK(result, err, "pulsar_admin_namespace create"))

	result, err = callTool(ctx, adminClient, "pulsar_admin_namespace_policy_set", map[string]any{
		"namespace": namespace,
		"policy":    "permission",
		"role":      defaultTestUser,
		"actions":   []string{"consume"},
	})
	require.NoError(t, requireToolOK(result, err, "pulsar_admin_namespace_policy_set permission"))

	result, err = callTool(ctx, adminClient, "pulsar_admin_topic", map[string]any{
		"resource":   "topic",
		"operation":  "create",
		"topic":      topic,
		"partitions": float64(0),
	})
	require.NoError(t, requireToolOK(result, err, "pulsar_admin_topic create"))

	result, err = callTool(ctx, testClient, "pulsar_admin_tenant", map[string]any{
		"resource":        "tenant",
		"operation":       "create",
		"tenant":          tenant + "-unauthorized",
		"adminRoles":      []string{defaultTestUser},
		"allowedClusters": []string{clusters[0]},
	})
	require.NoError(t, requireToolError(result, err, "pulsar_admin_tenant unauthorized create"))
}

func startPulsarContainerWithAuth(ctx context.Context, cfg testcontainersConfig, adminToken, secretKeyPath string) (*testcontainersEnv, error) {
	if !cfg.Enabled {
		return nil, errTestcontainersDisabled
	}
	if strings.TrimSpace(adminToken) == "" {
		return nil, fmt.Errorf("admin token is required")
	}
	if strings.TrimSpace(secretKeyPath) == "" {
		return nil, fmt.Errorf("secret key path is required")
	}

	request := testcontainers.ContainerRequest{
		Image:        cfg.PulsarImage,
		ExposedPorts: []string{pulsarBrokerPort, pulsarWebServicePort},
		Env: map[string]string{
			"PULSAR_PREFIX_authenticationEnabled":                "true",
			"PULSAR_PREFIX_authenticationProviders":              "org.apache.pulsar.broker.authentication.AuthenticationProviderToken",
			"PULSAR_PREFIX_authorizationEnabled":                 "true",
			"PULSAR_PREFIX_superUserRoles":                       defaultAdminRole,
			"PULSAR_PREFIX_tokenSecretKey":                       "file://" + authSecretFilePath,
			"PULSAR_PREFIX_brokerClientAuthenticationPlugin":     "org.apache.pulsar.client.impl.auth.AuthenticationToken",
			"PULSAR_PREFIX_brokerClientAuthenticationParameters": "token:" + adminToken,
		},
		Cmd: []string{"bash", "-lc", "set -- $(hostname -i); export PULSAR_PREFIX_advertisedAddress=$1; bin/apply-config-from-env.py /pulsar/conf/standalone.conf; exec bin/pulsar standalone"},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      secretKeyPath,
				ContainerFilePath: authSecretFilePath,
				FileMode:          0o600,
			},
		},
		WaitingFor: wait.ForAll(
			wait.ForHTTP("/admin/v2/clusters").
				WithPort(pulsarWebServicePort).
				WithHeaders(map[string]string{"Authorization": "Bearer " + adminToken}),
			wait.ForListeningPort(pulsarBrokerPort),
		).WithDeadline(cfg.StartupTimeout),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("start pulsar container: %w", err)
	}

	pulsarWebURL, pulsarBrokerURL, err := resolvePulsarEndpoints(ctx, container)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, err
	}

	return &testcontainersEnv{
		Pulsar:              container,
		PulsarWebServiceURL: pulsarWebURL,
		PulsarBrokerURL:     pulsarBrokerURL,
	}, nil
}

func startSNMCPServer(t *testing.T, pulsarWebURL, pulsarBrokerURL string) (string, func(), error) {
	t.Helper()

	repoRoot := findRepoRoot(t)
	binaryPath := buildSNMCPBinary(t, repoRoot)

	addr := reserveLocalAddr(t)
	baseURL := fmt.Sprintf("http://%s/mcp", addr)

	//nolint:gosec // test binary path and arguments are controlled
	cmd := exec.Command(binaryPath,
		"sse",
		"--http-addr", addr,
		"--http-path", "/mcp",
		"--use-external-pulsar",
		"--pulsar-web-service-url", pulsarWebURL,
		"--pulsar-service-url", pulsarBrokerURL,
		"--multi-session-pulsar",
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
	if err := waitForHTTPStatus(readyCtx, baseURL+"/healthz", http.StatusOK, errCh); err != nil {
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

func waitForHTTPStatus(ctx context.Context, target string, status int, errCh <-chan error) error {
	for {
		select {
		case err := <-errCh:
			return fmt.Errorf("snmcp exited early: %w", err)
		default:
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
		if err != nil {
			return err
		}
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == status {
				return nil
			}
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for %s", target)
		case <-time.After(250 * time.Millisecond):
		}
	}
}

func buildSNMCPBinary(t *testing.T, repoRoot string) string {
	t.Helper()
	outputPath := filepath.Join(t.TempDir(), "snmcp")

	//nolint:gosec // go build command is static and controlled in tests
	cmd := exec.Command("go", "build", "-o", outputPath, "./cmd/streamnative-mcp-server")
	cmd.Dir = repoRoot
	cmd.Env = os.Environ()

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	require.NoErrorf(t, cmd.Run(), "failed to build snmcp: %s", strings.TrimSpace(output.String()))
	return outputPath
}

func reserveLocalAddr(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := listener.Addr().String()
	require.NoError(t, listener.Close())
	return addr
}

func loadAuthTokens(t *testing.T) authTokens {
	t.Helper()

	tokens := authTokens{
		AdminToken:    strings.TrimSpace(os.Getenv("ADMIN_TOKEN")),
		TestUserToken: strings.TrimSpace(os.Getenv("TEST_USER_TOKEN")),
	}

	if tokens.AdminToken != "" && tokens.TestUserToken != "" {
		return tokens
	}

	repoRoot := findRepoRoot(t)
	path := filepath.Join(repoRoot, "charts", "snmcp", "e2e", "test-tokens.env")
	//nolint:gosec // path is resolved from repo root in tests
	file, err := os.Open(path)
	require.NoError(t, err)
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "ADMIN_TOKEN":
			if tokens.AdminToken == "" {
				tokens.AdminToken = value
			}
		case "TEST_USER_TOKEN":
			if tokens.TestUserToken == "" {
				tokens.TestUserToken = value
			}
		}
	}
	require.NoError(t, scanner.Err())
	require.NotEmpty(t, tokens.AdminToken)
	require.NotEmpty(t, tokens.TestUserToken)
	return tokens
}

func loadSecretKeyPath(t *testing.T) string {
	t.Helper()
	repoRoot := findRepoRoot(t)
	path := filepath.Join(repoRoot, "charts", "snmcp", "e2e", "test-secret.key")
	_, err := os.Stat(path)
	require.NoError(t, err)
	return path
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	require.NoError(t, err)

	for {
		if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	require.Fail(t, "go.mod not found in parent directories")
	return ""
}
