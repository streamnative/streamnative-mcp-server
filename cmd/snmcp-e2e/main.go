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

// Package main is the entry point for the StreamNative MCP E2E test.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

type config struct {
	httpBaseURL   string
	adminToken    string
	testUserToken string
	timeout       time.Duration
	verbose       bool
}

func main() {
	cfg, err := parseConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	err = run(ctx, cfg)
	cancel()
	if err != nil {
		fmt.Fprintf(os.Stderr, "e2e failed: %v\n", err)
		os.Exit(1)
	}
	_, _ = fmt.Fprintln(os.Stdout, "e2e succeeded")
}

func parseConfig() (config, error) {
	var cfg config
	flag.StringVar(&cfg.httpBaseURL, "http-base", getenv("E2E_HTTP_BASE", "http://127.0.0.1:9090/mcp"), "HTTP base URL for MCP endpoints")
	flag.StringVar(&cfg.adminToken, "admin-token", getenv("ADMIN_TOKEN", ""), "Admin JWT token")
	flag.StringVar(&cfg.testUserToken, "test-user-token", getenv("TEST_USER_TOKEN", ""), "Test user JWT token")
	flag.DurationVar(&cfg.timeout, "timeout", 3*time.Minute, "Overall timeout for the E2E run")
	flag.BoolVar(&cfg.verbose, "verbose", getenvBool("E2E_VERBOSE", false), "Enable verbose logging")
	flag.Parse()

	if cfg.adminToken == "" {
		return config{}, errors.New("admin token is required")
	}
	if cfg.testUserToken == "" {
		return config{}, errors.New("test-user token is required")
	}

	normalized, err := normalizeBaseURL(cfg.httpBaseURL)
	if err != nil {
		return config{}, err
	}
	cfg.httpBaseURL = normalized
	return cfg, nil
}

func normalizeBaseURL(raw string) (string, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("invalid http-base URL: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("invalid http-base URL: %s", raw)
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/")
	return parsed.String(), nil
}

func run(ctx context.Context, cfg config) error {
	logf(cfg.verbose, "http base: %s", cfg.httpBaseURL)
	if err := checkHealth(ctx, cfg.httpBaseURL); err != nil {
		return err
	}

	sseURL := cfg.httpBaseURL + "/sse"
	if err := expectUnauthorized(ctx, sseURL, "", cfg.verbose); err != nil {
		return err
	}
	if err := expectUnauthorized(ctx, sseURL, "invalid-token", cfg.verbose); err != nil {
		return err
	}

	adminClient, err := newAuthedClient(ctx, sseURL, cfg.adminToken, "snmcp-e2e-admin")
	if err != nil {
		return err
	}
	defer func() {
		_ = adminClient.Close()
	}()

	testClient, err := newAuthedClient(ctx, sseURL, cfg.testUserToken, "snmcp-e2e-test-user")
	if err != nil {
		return err
	}
	defer func() {
		_ = testClient.Close()
	}()

	clusters, err := listClusters(ctx, adminClient)
	if err != nil {
		return err
	}
	if len(clusters) == 0 {
		return errors.New("no clusters returned from pulsar_admin_cluster")
	}
	cluster := clusters[0]

	suffix := time.Now().UnixNano()
	tenant := fmt.Sprintf("e2e-%d", suffix)
	namespaceName := fmt.Sprintf("ns-%d", suffix)
	namespace := fmt.Sprintf("%s/%s", tenant, namespaceName)
	topic := fmt.Sprintf("persistent://%s/topic-%d", namespace, suffix)
	concurrentTopic := fmt.Sprintf("persistent://%s/topic-concurrent-%d", namespace, suffix)
	functionInputTopic := fmt.Sprintf("persistent://%s/function-input-%d", namespace, suffix)
	functionOutputTopic := fmt.Sprintf("persistent://%s/function-output-%d", namespace, suffix)
	functionName := fmt.Sprintf("echo-%d", suffix)

	result, err := callTool(ctx, adminClient, "pulsar_admin_tenant", map[string]any{
		"resource":        "tenant",
		"operation":       "create",
		"tenant":          tenant,
		"adminRoles":      []string{"admin"},
		"allowedClusters": []string{cluster},
	})
	if err := requireToolOK(result, err, "pulsar_admin_tenant create"); err != nil {
		return err
	}

	result, err = callTool(ctx, adminClient, "pulsar_admin_namespace", map[string]any{
		"operation": "create",
		"namespace": namespace,
		"clusters":  []string{cluster},
	})
	if err := requireToolOK(result, err, "pulsar_admin_namespace create"); err != nil {
		return err
	}

	result, err = callTool(ctx, adminClient, "pulsar_admin_namespace_policy_set", map[string]any{
		"namespace": namespace,
		"policy":    "permission",
		"role":      "test-user",
		"actions":   []string{"consume"},
	})
	if err := requireToolOK(result, err, "pulsar_admin_namespace_policy_set permission"); err != nil {
		return err
	}

	result, err = callTool(ctx, adminClient, "pulsar_admin_topic", map[string]any{
		"resource":   "topic",
		"operation":  "create",
		"topic":      topic,
		"partitions": float64(0),
	})
	if err := requireToolOK(result, err, "pulsar_admin_topic create"); err != nil {
		return err
	}

	result, err = callTool(ctx, adminClient, "pulsar_client_produce", map[string]any{
		"topic":    topic,
		"messages": []string{"admin-message"},
	})
	if err := requireToolOK(result, err, "pulsar_client_produce admin"); err != nil {
		return err
	}

	result, err = callTool(ctx, testClient, "pulsar_client_consume", map[string]any{
		"topic":             topic,
		"subscription-name": fmt.Sprintf("sub-%d", suffix),
		"initial-position":  "earliest",
		"num-messages":      float64(1),
		"timeout":           float64(15),
		"subscription-type": "exclusive",
		"subscription-mode": "durable",
		"show-properties":   false,
		"hide-payload":      false,
	})
	if err := requireToolOK(result, err, "pulsar_client_consume test-user"); err != nil {
		return err
	}

	result, err = callTool(ctx, testClient, "pulsar_client_produce", map[string]any{
		"topic":    topic,
		"messages": []string{"unauthorized-message"},
	})
	if err := requireToolError(result, err, "pulsar_client_produce test-user"); err != nil {
		return err
	}

	result, err = callTool(ctx, adminClient, "pulsar_admin_topic", map[string]any{
		"resource":   "topic",
		"operation":  "create",
		"topic":      functionInputTopic,
		"partitions": float64(0),
	})
	if err := requireToolOK(result, err, "pulsar_admin_topic create function input"); err != nil {
		return err
	}

	result, err = callTool(ctx, adminClient, "pulsar_admin_topic", map[string]any{
		"resource":   "topic",
		"operation":  "create",
		"topic":      functionOutputTopic,
		"partitions": float64(0),
	})
	if err := requireToolOK(result, err, "pulsar_admin_topic create function output"); err != nil {
		return err
	}

	logf(cfg.verbose, "creating function: tenant=%s namespace=%s name=%s inputs=%v output=%s py=%s classname=%s",
		tenant,
		namespaceName,
		functionName,
		[]string{functionInputTopic},
		functionOutputTopic,
		"/server/e2e/functions/echo.py",
		"echo.EchoFunction",
	)
	result, err = callTool(ctx, adminClient, "pulsar_admin_functions", map[string]any{
		"operation": "create",
		"tenant":    tenant,
		"namespace": namespaceName,
		"name":      functionName,
		"classname": "echo.EchoFunction",
		"inputs":    []string{functionInputTopic},
		"output":    functionOutputTopic,
		"py":        "/server/e2e/functions/echo.py",
	})
	if err := requireToolOK(result, err, "pulsar_admin_functions create"); err != nil {
		return err
	}

	if err := waitForFunctionRunning(ctx, adminClient, tenant, namespaceName, functionName, 60*time.Second); err != nil {
		return err
	}

	result, err = callTool(ctx, adminClient, "pulsar_admin_functions", map[string]any{
		"operation":  "stats",
		"tenant":     tenant,
		"namespace":  namespaceName,
		"name":       functionName,
		"instanceId": float64(0),
	})
	if err := requireToolOK(result, err, "pulsar_admin_functions stats"); err != nil {
		return err
	}

	triggerValue := "e2e-trigger"
	result, err = callTool(ctx, adminClient, "pulsar_admin_functions", map[string]any{
		"operation":    "trigger",
		"tenant":       tenant,
		"namespace":    namespaceName,
		"name":         functionName,
		"topic":        functionInputTopic,
		"triggerValue": triggerValue,
	})
	if err := requireToolOK(result, err, "pulsar_admin_functions trigger"); err != nil {
		return err
	}
	if !strings.Contains(firstText(result), triggerValue) {
		return fmt.Errorf("unexpected trigger result: %s", firstText(result))
	}

	result, err = callTool(ctx, adminClient, "pulsar_admin_functions", map[string]any{
		"operation":  "update",
		"tenant":     tenant,
		"namespace":  namespaceName,
		"name":       functionName,
		"userConfig": map[string]any{"updated": true},
	})
	if err := requireToolOK(result, err, "pulsar_admin_functions update"); err != nil {
		return err
	}

	result, err = callTool(ctx, adminClient, "pulsar_admin_functions", map[string]any{
		"operation": "get",
		"tenant":    tenant,
		"namespace": namespaceName,
		"name":      functionName,
	})
	if err := requireToolOK(result, err, "pulsar_admin_functions get"); err != nil {
		return err
	}
	if err := assertFunctionUserConfig(firstText(result), "updated"); err != nil {
		return err
	}

	result, err = callTool(ctx, adminClient, "pulsar_admin_functions", map[string]any{
		"operation": "delete",
		"tenant":    tenant,
		"namespace": namespaceName,
		"name":      functionName,
	})
	if err := requireToolOK(result, err, "pulsar_admin_functions delete"); err != nil {
		return err
	}

	result, err = callTool(ctx, adminClient, "pulsar_admin_topic", map[string]any{
		"resource":   "topic",
		"operation":  "create",
		"topic":      concurrentTopic,
		"partitions": float64(0),
	})
	if err := requireToolOK(result, err, "pulsar_admin_topic create concurrent"); err != nil {
		return err
	}

	if err := runConcurrent(ctx, adminClient, testClient, concurrentTopic, fmt.Sprintf("sub-concurrent-%d", suffix)); err != nil {
		return err
	}

	return nil
}

func checkHealth(ctx context.Context, httpBaseURL string) error {
	healthURL := httpBaseURL + "/healthz"
	readyURL := httpBaseURL + "/readyz"

	if err := expectStatusOK(ctx, healthURL); err != nil {
		return fmt.Errorf("healthz check failed: %w", err)
	}
	if err := expectStatusOK(ctx, readyURL); err != nil {
		return fmt.Errorf("readyz check failed: %w", err)
	}
	return nil
}

func expectStatusOK(ctx context.Context, target string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func expectUnauthorized(ctx context.Context, sseURL, token string, verbose bool) error {
	logf(verbose, "checking unauthorized: url=%s token_present=%t", sseURL, token != "")
	status, err := probeSSEStatus(ctx, sseURL, token)
	if err != nil {
		logf(verbose, "probe SSE failed: %v", err)
		return err
	}
	logf(verbose, "probe SSE status: %d", status)
	if status == http.StatusUnauthorized || status == http.StatusForbidden {
		return nil
	}
	if token == "" {
		return fmt.Errorf("expected unauthorized status for %s, got %d", sseURL, status)
	}

	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}
	c, err := client.NewSSEMCPClient(sseURL, client.WithHeaders(headers))
	if err != nil {
		return err
	}
	defer func() {
		_ = c.Close()
	}()

	if err := c.Start(ctx); err != nil {
		logf(verbose, "sse start error: %v", err)
		if isAuthError(err) {
			return nil
		}
		return fmt.Errorf("expected auth error for %s, got %v", sseURL, err)
	}

	if err := initializeClient(ctx, c, "snmcp-e2e-unauthorized"); err != nil {
		logf(verbose, "initialize error: %v", err)
		if isAuthError(err) {
			return nil
		}
		return fmt.Errorf("expected auth error during initialize for %s, got %v", sseURL, err)
	}

	result, err := callTool(ctx, c, "pulsar_admin_cluster", map[string]any{
		"resource":  "cluster",
		"operation": "list",
	})
	if err != nil {
		logf(verbose, "tool call error: %v", err)
		if isAuthError(err) {
			return nil
		}
		return fmt.Errorf("expected auth error for %s, got %v", sseURL, err)
	}
	if result == nil || !result.IsError {
		logf(verbose, "tool call result: %#v", result)
		return fmt.Errorf("expected auth error for %s", sseURL)
	}
	if !isAuthText(firstText(result)) {
		logf(verbose, "tool call error text: %s", firstText(result))
		return fmt.Errorf("expected auth error for %s, got %s", sseURL, firstText(result))
	}
	return nil
}

func newAuthedClient(ctx context.Context, sseURL, token, clientName string) (*client.Client, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}
	c, err := client.NewSSEMCPClient(sseURL, client.WithHeaders(headers))
	if err != nil {
		return nil, err
	}
	if err := c.Start(ctx); err != nil {
		return nil, err
	}
	if err := initializeClient(ctx, c, clientName); err != nil {
		return nil, err
	}
	return c, nil
}

func initializeClient(ctx context.Context, c *client.Client, name string) error {
	req := mcp.InitializeRequest{}
	req.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	req.Params.ClientInfo = mcp.Implementation{
		Name:    name,
		Version: "1.0.0",
	}
	_, err := c.Initialize(ctx, req)
	return err
}

func callTool(ctx context.Context, c *client.Client, name string, args map[string]any) (*mcp.CallToolResult, error) {
	request := mcp.CallToolRequest{}
	request.Params.Name = name
	request.Params.Arguments = args
	return c.CallTool(ctx, request)
}

func requireToolOK(result *mcp.CallToolResult, err error, label string) error {
	if err != nil {
		return fmt.Errorf("%s failed: %w", label, err)
	}
	if result.IsError {
		return fmt.Errorf("%s returned error: %s", label, firstText(result))
	}
	return nil
}

func requireToolError(result *mcp.CallToolResult, err error, label string) error {
	if err != nil {
		return fmt.Errorf("%s failed: %w", label, err)
	}
	if !result.IsError {
		return fmt.Errorf("%s expected error, got success: %s", label, firstText(result))
	}
	return nil
}

func firstText(result *mcp.CallToolResult) string {
	if result == nil || len(result.Content) == 0 {
		return ""
	}
	if text, ok := result.Content[0].(mcp.TextContent); ok {
		return text.Text
	}
	return ""
}

func probeSSEStatus(ctx context.Context, sseURL, token string) (int, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, sseURL, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Accept", "text/event-stream")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	httpClient := &http.Client{Timeout: 5 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return resp.StatusCode, nil
}

func isAuthError(err error) bool {
	if err == nil {
		return false
	}
	return isAuthText(err.Error())
}

func isAuthText(text string) bool {
	if text == "" {
		return false
	}
	lower := strings.ToLower(text)
	if strings.Contains(lower, "status code: 401") || strings.Contains(lower, "status code: 403") {
		return true
	}
	if strings.Contains(lower, " 401") || strings.Contains(lower, " 403") {
		return true
	}
	if strings.Contains(lower, "unauthorized") {
		return true
	}
	if strings.Contains(lower, "authentication") || strings.Contains(lower, "authorization") {
		return true
	}
	if strings.Contains(lower, "permission") || strings.Contains(lower, "access denied") {
		return true
	}
	if strings.Contains(lower, "missing authorization") || strings.Contains(lower, "session not found") {
		return true
	}
	return false
}

func listClusters(ctx context.Context, c *client.Client) ([]string, error) {
	result, err := callTool(ctx, c, "pulsar_admin_cluster", map[string]any{
		"resource":  "cluster",
		"operation": "list",
	})
	if err := requireToolOK(result, err, "pulsar_admin_cluster list"); err != nil {
		return nil, err
	}
	raw := firstText(result)
	if raw == "" {
		return nil, errors.New("empty cluster list result")
	}
	var clusters []string
	if err := json.Unmarshal([]byte(raw), &clusters); err != nil {
		return nil, fmt.Errorf("failed to parse cluster list: %w", err)
	}
	return clusters, nil
}

type functionStatus struct {
	NumRunning int `json:"numRunning"`
}

func waitForFunctionRunning(ctx context.Context, c *client.Client, tenant, namespace, name string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		status, err := getFunctionStatus(ctx, c, tenant, namespace, name)
		if err == nil && status.NumRunning > 0 {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
	return fmt.Errorf("function %s did not reach running state within %s", name, timeout.String())
}

func getFunctionStatus(ctx context.Context, c *client.Client, tenant, namespace, name string) (functionStatus, error) {
	result, err := callTool(ctx, c, "pulsar_admin_functions", map[string]any{
		"operation": "status",
		"tenant":    tenant,
		"namespace": namespace,
		"name":      name,
	})
	if err := requireToolOK(result, err, "pulsar_admin_functions status"); err != nil {
		return functionStatus{}, err
	}
	raw := firstText(result)
	if raw == "" {
		return functionStatus{}, errors.New("empty function status result")
	}
	var status functionStatus
	if err := json.Unmarshal([]byte(raw), &status); err != nil {
		return functionStatus{}, fmt.Errorf("failed to parse function status: %w", err)
	}
	return status, nil
}

func assertFunctionUserConfig(raw string, key string) error {
	if raw == "" {
		return errors.New("empty function config result")
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return fmt.Errorf("failed to parse function config: %w", err)
	}
	userConfig, ok := payload["userConfig"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("missing userConfig in function config")
	}
	if _, ok := userConfig[key]; !ok {
		return fmt.Errorf("userConfig missing key: %s", key)
	}
	return nil
}

func runConcurrent(ctx context.Context, adminClient, testClient *client.Client, topic, subscription string) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		result, err := callTool(ctx, adminClient, "pulsar_client_produce", map[string]any{
			"topic":    topic,
			"messages": []string{"concurrent-message"},
		})
		err = requireToolOK(result, err, "pulsar_client_produce concurrent admin")
		if err != nil {
			errCh <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		result, err := callTool(ctx, testClient, "pulsar_client_consume", map[string]any{
			"topic":             topic,
			"subscription-name": subscription,
			"initial-position":  "earliest",
			"num-messages":      float64(1),
			"timeout":           float64(15),
		})
		err = requireToolOK(result, err, "pulsar_client_consume concurrent test-user")
		if err != nil {
			errCh <- err
		}
	}()

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}

func getenv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getenvBool(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return fallback
	}
}

func logf(enabled bool, format string, args ...any) {
	if !enabled {
		return
	}
	_, _ = fmt.Fprintf(os.Stderr, "[e2e] "+format+"\n", args...)
}
