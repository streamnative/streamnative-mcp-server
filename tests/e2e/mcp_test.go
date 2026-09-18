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
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const requestTimeout = 30 * time.Second

// No endpoint or credential overrides: every operation targets this run's new,
// disposable local container. The test must fail (not skip) without Docker.
func TestPulsar(t *testing.T) {
	artifacts := newArtifacts(t)
	webURL, serviceURL := startPulsar(t, artifacts)
	binary := buildServer(t, artifacts)
	for _, mode := range []string{"stdio", "sse", "http"} {
		t.Run(mode, func(t *testing.T) {
			c := startServer(t, binary, mode, artifacts, []string{
				"--use-external-pulsar", "--pulsar-web-service-url", webURL, "--pulsar-service-url", serviceURL,
			})
			initialize(t, c, mode)
			checkTools(t, c, "pulsar_admin_topic_write", "pulsar_admin_topic_read")
			checkTopicLifecycle(t, c, webURL)
		})
	}
}

func newArtifacts(t *testing.T) string {
	t.Helper()
	base := os.Getenv("E2E_ARTIFACTS_DIR")
	if base != "" {
		require.NoError(t, os.MkdirAll(base, 0700)) // #nosec G703 -- explicitly configured artifact output, never a backend/config input.
	}
	artifacts, err := os.MkdirTemp(base, "snmcp-e2e-logs-"+t.Name()+"-")
	require.NoError(t, err)
	t.Logf("diagnostic directory (retained on failure): %s", artifacts)
	t.Cleanup(func() {
		if t.Failed() {
			t.Logf("diagnostics retained: %s", artifacts)
		} else {
			assert.NoError(t, os.RemoveAll(artifacts)) // #nosec G703 -- removes only the unique MkdirTemp child, never the configured parent.
		}
	})

	return artifacts
}

func requireLocalDocker(t *testing.T) {
	t.Helper()
	// Testcontainers can panic during Docker socket discovery. Report a failed
	// prerequisite instead of silently skipping or losing the diagnostic path.
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("local Docker/Testcontainers prerequisite failed: %v", recovered)
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	docker, err := testcontainers.NewDockerClientWithOpts(ctx)
	require.NoError(t, err, "a running local Docker daemon is required")
	defer func() { assert.NoError(t, docker.Close()) }()
	require.True(t, strings.HasPrefix(docker.DaemonHost(), "unix://") || strings.HasPrefix(docker.DaemonHost(), "npipe://"),
		"remote Docker daemons are forbidden: %s", docker.DaemonHost())
	if configured := os.Getenv("DOCKER_HOST"); configured != "" {
		// Testcontainers may fall back to a discovered socket when an explicit
		// host is unavailable. A broken prerequisite must not silently pass.
		require.Equal(t, configured, docker.DaemonHost(), "refusing fallback from explicit DOCKER_HOST; check local Docker configuration")
	}
	_, err = docker.Ping(ctx)
	require.NoError(t, err, "a running local Docker daemon is required")
}

func startPulsar(t *testing.T, artifacts string) (string, string) {
	t.Helper()
	requireLocalDocker(t)
	startup, stop := context.WithTimeout(context.Background(), 4*time.Minute)
	defer stop()
	broker, err := testcontainers.GenericContainer(startup, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "apachepulsar/pulsar:4.1.0",
			ExposedPorts: []string{"8080/tcp", "6650/tcp"},
			Env:          map[string]string{"PULSAR_MEM": "-Xms256m -Xmx512m -XX:MaxDirectMemorySize=256m"},
			Cmd:          []string{"bin/pulsar", "standalone", "--no-functions-worker", "--no-stream-storage"},
			HostConfigModifier: func(config *container.HostConfig) {
				config.PortBindings = nat.PortMap{
					"8080/tcp": {{HostIP: "127.0.0.1"}},
					"6650/tcp": {{HostIP: "127.0.0.1"}},
				}
			},
			WaitingFor: wait.ForHTTP("/admin/v2/brokers/health").WithPort("8080/tcp").WithStartupTimeout(2 * time.Minute),
		},
	})
	// Register cleanup even if creation returned a partially created container.
	cleanupContainer(t, broker, artifacts, "pulsar")
	require.NoError(t, err, "create disposable Pulsar container")
	require.NoError(t, broker.Start(startup), "Pulsar readiness; see pulsar.log")
	host, err := broker.Host(startup)
	require.NoError(t, err)
	require.True(t, host == "localhost" || net.ParseIP(host).IsLoopback(), "only local container endpoints are allowed: %s", host)
	webPort, err := broker.MappedPort(startup, "8080/tcp")
	require.NoError(t, err)
	servicePort, err := broker.MappedPort(startup, "6650/tcp")
	require.NoError(t, err)
	return "http://" + net.JoinHostPort(host, webPort.Port()), "pulsar://" + net.JoinHostPort(host, servicePort.Port())
}

func cleanupContainer(t *testing.T, broker testcontainers.Container, artifacts, backend string) {
	t.Helper()
	if broker != nil {
		t.Cleanup(func() {
			logCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			logs, logErr := broker.Logs(logCtx)
			if assert.NoError(t, logErr, "collect %s logs", backend) {
				data, readErr := io.ReadAll(logs)
				assert.NoError(t, readErr)
				assert.NoError(t, logs.Close())
				assert.NoError(t, os.WriteFile(filepath.Join(artifacts, backend+".log"), data, 0600))
			}
			cancel()
			cleanupCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			// Removing the container and its anonymous volumes also removes all
			// test data, including topics left by failed assertions. No reuse.
			assert.NoError(t, broker.Terminate(cleanupCtx), "remove %s container and data", backend)
		})
	}
}

func buildServer(t *testing.T, artifacts string) string {
	t.Helper()
	root, err := filepath.Abs("../..")
	require.NoError(t, err)
	binary := filepath.Join(t.TempDir(), "snmcp")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "build", "-o", binary, "./cmd/streamnative-mcp-server") // #nosec G204 -- fixed build target and test-owned output path.
	cmd.Dir = root
	output, err := cmd.CombinedOutput()
	require.NoError(t, os.WriteFile(filepath.Join(artifacts, "build.log"), output, 0600))
	require.NoError(t, err, "build real server: %s", output)
	return binary
}

func startServer(t *testing.T, binary, mode, artifacts string, backendArgs []string) *client.Client {
	t.Helper()
	home := t.TempDir()
	// A whitelist, not os.Environ(): do not inherit SNMCP_*, proxy settings,
	// Cloud credentials, or the user's saved backend context.
	env := []string{"HOME=" + home, "PATH=" + os.Getenv("PATH"), "TMPDIR=" + home}
	args := append([]string{mode, "--config-dir", filepath.Join(home, "config")}, backendArgs...)
	logFile, err := os.OpenFile(filepath.Join(artifacts, mode+".log"), os.O_CREATE|os.O_WRONLY, 0600) // #nosec G304 -- test-owned directory and fixed transport names.
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, logFile.Close()) })
	var c *client.Client
	if mode == "stdio" {
		// The NewStdioMCPClient options variant permits a hermetic child
		// environment and full stderr capture instead of the SDK's ring buffer.
		c, err = client.NewStdioMCPClientWithOptions(binary, env, args, transport.WithCommandFunc(
			func(ctx context.Context, command string, env, args []string) (*exec.Cmd, error) {
				cmd := exec.CommandContext(ctx, command, args...) // #nosec G204 -- only the freshly built binary and local fixture arguments.
				cmd.Env, cmd.Dir, cmd.Stderr = env, home, logFile
				return cmd, nil
			}))
	} else {
		// The existing SSE command cannot report its actual port with :0.
		// Ask the OS for a free loopback port, then hand it to the binary.
		listener, listenErr := net.Listen("tcp", "127.0.0.1:0")
		require.NoError(t, listenErr)
		addr := listener.Addr().String()
		require.NoError(t, listener.Close())
		args = append(args, "--http-addr", addr, "--http-path", "/mcp")
		cmd := exec.Command(binary, args...) // #nosec G204 -- only the freshly built binary and local fixture arguments.
		cmd.Env, cmd.Dir, cmd.Stdout, cmd.Stderr = env, home, logFile, logFile
		require.NoError(t, cmd.Start())
		exited := make(chan struct{})
		var exitErr error
		go func() { exitErr = cmd.Wait(); close(exited) }()
		t.Cleanup(func() {
			select {
			case <-exited:
				t.Errorf("%s server exited before cleanup: %v", mode, exitErr)
				return
			default:
			}
			assert.NoError(t, cmd.Process.Signal(os.Interrupt))
			select {
			case <-exited:
				assert.NoError(t, exitErr, "server shutdown")
			case <-time.After(15 * time.Second):
				assert.NoError(t, cmd.Process.Kill())
				<-exited // Reap the child; outputs are files, not blocking pipes.
				t.Error("server did not stop gracefully; killed and reaped")
			}
		})
		endpoint := "http://" + addr + "/mcp"
		httpClient := localHTTPClient(t)
		require.Eventually(t, func() bool {
			select {
			case <-exited:
				return true // Report on the test goroutine, not Eventually's worker.
			default:
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			status, _, err := nativeRequest(ctx, httpClient, endpoint+"/readyz")
			return err == nil && status == http.StatusOK
		}, requestTimeout, 100*time.Millisecond, "%s readiness (see %s.log)", mode, mode)
		select {
		case <-exited:
			t.Fatalf("%s server exited during startup: %v (see %s.log)", mode, exitErr, mode)
		default:
		}
		if mode == "sse" {
			httpClient.Timeout = 2 * time.Minute // Bound the SSE stream, not each RPC.
			c, err = client.NewSSEMCPClient(endpoint+"/sse", client.WithHTTPClient(httpClient))
		} else {
			c, err = client.NewStreamableHttpClient(endpoint, transport.WithHTTPBasicClient(httpClient))
		}
	}
	require.NoError(t, err)
	// SDK stdio Close closes stdin, waits, then escalates TERM/KILL with
	// bounded waits and reaps the binary. HTTP processes are owned above.
	t.Cleanup(func() { assert.NoError(t, c.Close(), "close MCP client") })
	if mode != "stdio" {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		t.Cleanup(cancel) // SSE keeps this context for its connection lifetime.
		require.NoError(t, c.Start(ctx))
	}
	return c
}

func initialize(t *testing.T, c *client.Client, mode string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	version := mcp.LATEST_PROTOCOL_VERSION
	if mode == "sse" {
		version = mcp.LATEST_LEGACY_PROTOCOL_VERSION
	}
	request := mcp.InitializeRequest{}
	request.Params.ProtocolVersion = version
	request.Params.ClientInfo = mcp.Implementation{Name: "snmcp-binary-e2e", Version: "1"}
	// v1.1.0 Initialize sends server/discover for modern revisions, and
	// initialize + notifications/initialized for legacy SSE. Check the exact
	// negotiated version so an unintended fallback cannot pass the test.
	result, err := c.Initialize(ctx, request)
	require.NoError(t, err)
	require.Equal(t, version, result.ProtocolVersion)
	require.Equal(t, version, c.ProtocolVersion())
	require.NotEmpty(t, result.ServerInfo.Name)
	require.NotNil(t, result.Capabilities.Tools)
	if mode != "sse" {
		discovery, err := c.Discover(ctx, mcp.DiscoverRequest{})
		require.NoError(t, err)
		require.Contains(t, discovery.SupportedVersions, version)
	}
}

func checkTools(t *testing.T, c *client.Client, expected ...string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	names := map[string]bool{}
	request := mcp.ListToolsRequest{}
	for page := 0; ; page++ {
		require.Less(t, page, 100, "tools/list pagination must terminate")
		result, err := c.ListTools(ctx, request)
		require.NoError(t, err)
		for _, tool := range result.Tools {
			names[tool.Name] = true
		}
		if result.NextCursor == "" {
			break
		}
		request.Params.Cursor = result.NextCursor
	}
	for _, name := range expected {
		require.Contains(t, names, name)
	}
}

func checkTopicLifecycle(t *testing.T, c *client.Client, webURL string) {
	t.Helper()
	name := "e2e-" + strings.ToLower(rand.Text())
	topic := "persistent://public/default/" + name
	httpClient := localHTTPClient(t)
	metadataURL := webURL + "/admin/v2/persistent/public/default/" + name + "/partitions?checkAllowAutoCreation=false"
	listURL := webURL + "/admin/v2/persistent/public/default/partitioned"
	var topics []string
	nativeJSON(t, httpClient, listURL, &topics)
	require.NotNil(t, topics, "native topic list must be a JSON array")
	require.NotContains(t, topics, topic)

	args := map[string]any{"resource": "topic", "operation": "create", "topic": topic, "partitions": 2}
	created := callTool(t, c, "pulsar_admin_topic_write", args, false)
	require.Equal(t, fmt.Sprintf("Successfully created topic '%s' with 2 partitions", topic), created)

	var metadata struct {
		Partitions int `json:"partitions"`
	}
	nativeJSON(t, httpClient, metadataURL, &metadata)
	require.Equal(t, 2, metadata.Partitions, "native Pulsar API must observe MCP's write")
	nativeJSON(t, httpClient, listURL, &topics)
	require.Contains(t, topics, topic)
	read := callTool(t, c, "pulsar_admin_topic_read", map[string]any{"resource": "topic", "operation": "get", "topic": topic}, false)
	metadata.Partitions = 0 // Missing fields must not reuse a previous response.
	require.NoError(t, json.Unmarshal([]byte(read), &metadata))
	require.Equal(t, 2, metadata.Partitions)

	// A real backend conflict, not schema validation or a transport failure.
	args["partitions"] = 3
	conflict := callTool(t, c, "pulsar_admin_topic_write", args, true)
	require.Contains(t, conflict, "Failed to create topic")
	require.Contains(t, strings.ToLower(conflict), "already exist")
	metadata.Partitions = 0
	nativeJSON(t, httpClient, metadataURL, &metadata)
	require.Equal(t, 2, metadata.Partitions, "rejected create must not alter the existing topic")

	deleted := callTool(t, c, "pulsar_admin_topic_write", map[string]any{"resource": "topic", "operation": "delete", "topic": topic}, false)
	require.Equal(t, fmt.Sprintf("Successfully deleted topic '%s'", topic), deleted)
	nativeJSON(t, httpClient, listURL, &topics)
	require.NotNil(t, topics, "native topic list must be a JSON array")
	require.NotContains(t, topics, topic, "native Pulsar API must confirm deletion")
}

func callTool(t *testing.T, c *client.Client, name string, args map[string]any, isError bool) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	request := mcp.CallToolRequest{}
	request.Params.Name, request.Params.Arguments = name, args
	result, err := c.CallTool(ctx, request)
	require.NoError(t, err, "%s: transport/RPC errors are not tool errors", name)
	require.NotNil(t, result)
	require.Equal(t, isError, result.IsError, "%s: %+v", name, result.Content)
	require.Len(t, result.Content, 1)
	content, ok := result.Content[0].(mcp.TextContent)
	require.True(t, ok, "expected text tool result")
	require.NotEmpty(t, content.Text)
	return content.Text
}

func localHTTPClient(t *testing.T) *http.Client {
	t.Helper()
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.Proxy = nil // Never send local test traffic or credentials through a proxy.
	t.Cleanup(tr.CloseIdleConnections)
	return &http.Client{Transport: tr, Timeout: requestTimeout,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
}

func nativeRequest(ctx context.Context, c *http.Client, endpoint string) (int, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	return resp.StatusCode, body, err
}

func nativeJSON(t *testing.T, c *http.Client, endpoint string, out any) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	status, body, err := nativeRequest(ctx, c, endpoint)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status, "native Pulsar API: %s", body)
	require.NoError(t, json.Unmarshal(body, out))
}
