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

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

const httpProtocolVersion = "2026-07-28"

// This intentionally does not use the SDK client: the modern lane must never
// initialize or depend on a protocol session. Disable keep-alives to exercise
// repeated connections, not to claim coverage of multiple server replicas.
type httpE2EClient struct {
	endpoint string
	token    string
	client   *http.Client
}

func newHTTPE2EClient(endpoint, token string) *httpE2EClient {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DisableKeepAlives = true
	return &httpE2EClient{endpoint: endpoint, token: token, client: &http.Client{
		Transport: transport, Timeout: 45 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}}
}

func (c *httpE2EClient) request(ctx context.Context, method string, params map[string]any) (*http.Request, error) {
	p := make(map[string]any, len(params)+1)
	for k, v := range params {
		p[k] = v
	}
	p["_meta"] = map[string]any{
		"io.modelcontextprotocol/protocolVersion":    httpProtocolVersion,
		"io.modelcontextprotocol/clientCapabilities": map[string]any{},
	}
	body, err := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": 1, "method": method, "params": p})
	if err != nil {
		return nil, errors.New("cannot encode HTTP RPC request")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, errors.New("cannot construct HTTP RPC request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	req.Header.Set("Mcp-Protocol-Version", httpProtocolVersion)
	req.Header.Set("Mcp-Method", method)
	nameKey := "name"
	if method == "resources/read" {
		nameKey = "uri"
	}
	if name, ok := p[nameKey].(string); ok {
		req.Header.Set("Mcp-Name", name)
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	return req, nil
}

type httpRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *struct {
		Code int `json:"code"`
	} `json:"error"`
}

// Never include response bodies or transport errors in diagnostics: either can
// contain a backend-echoed token. Negative tests must validate HTTP and RPC
// errors separately; a network failure is not an authorization success.
func (c *httpE2EClient) exchange(req *http.Request, status int, rpcCode int) (json.RawMessage, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.New("HTTP RPC transport failed")
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != status {
		return nil, fmt.Errorf("HTTP RPC status: got %d, want %d", resp.StatusCode, status)
	}
	if len(resp.Header.Values("Mcp-Session-Id")) != 0 {
		return nil, errors.New("unexpected protocol session ID")
	}
	if status == http.StatusUnauthorized || status == http.StatusForbidden || status == http.StatusMethodNotAllowed {
		if status == http.StatusUnauthorized && resp.Header.Get("WWW-Authenticate") != "Bearer" {
			return nil, errors.New("missing Bearer challenge")
		}
		if status == http.StatusMethodNotAllowed && resp.Header.Get("Allow") != "POST" {
			return nil, errors.New("missing POST Allow header")
		}
		return nil, nil
	}
	media, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil || media != "application/json" {
		return nil, errors.New("expected JSON HTTP RPC response")
	}
	const maxResponse = 4 << 20
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponse+1))
	if err != nil || len(body) > maxResponse {
		return nil, errors.New("cannot read bounded HTTP RPC response")
	}
	var msg httpRPCResponse
	if json.Unmarshal(body, &msg) != nil || msg.JSONRPC != "2.0" || string(msg.ID) != "1" {
		return nil, errors.New("invalid HTTP RPC envelope or response ID")
	}
	if rpcCode != 0 {
		if msg.Error == nil || msg.Error.Code != rpcCode || len(msg.Result) != 0 {
			return nil, fmt.Errorf("expected RPC error code %d", rpcCode)
		}
		return nil, nil
	}
	if msg.Error != nil {
		return nil, fmt.Errorf("unexpected RPC error code %d", msg.Error.Code)
	}
	if len(msg.Result) == 0 || string(msg.Result) == "null" {
		return nil, errors.New("missing RPC result")
	}
	return msg.Result, nil
}

func (c *httpE2EClient) rpc(ctx context.Context, method string, params map[string]any, out any) error {
	req, err := c.request(ctx, method, params)
	if err != nil {
		return err
	}
	raw, err := c.exchange(req, http.StatusOK, 0)
	if err != nil {
		return fmt.Errorf("%s: %w", method, err)
	}
	if json.Unmarshal(raw, out) != nil {
		return fmt.Errorf("%s: invalid result shape", method)
	}
	return nil
}

type httpToolResult struct {
	IsError bool `json:"isError"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}

func (c *httpE2EClient) tool(ctx context.Context, name string, args map[string]any, denied bool) (string, error) {
	var result httpToolResult
	if err := c.rpc(ctx, "tools/call", map[string]any{"name": name, "arguments": args}, &result); err != nil {
		return "", err
	}
	if len(result.Content) == 0 || result.Content[0].Type != "text" || result.Content[0].Text == "" {
		return "", errors.New("empty tool result")
	}
	if result.IsError != denied {
		return "", fmt.Errorf("%s: unexpected tool isError=%t", name, result.IsError)
	}
	if denied && !isHTTPBackendAuthText(result.Content[0].Text) {
		return "", fmt.Errorf("%s: expected backend authorization denial", name)
	}
	return result.Content[0].Text, nil
}

func isHTTPBackendAuthText(text string) bool {
	lower := strings.ToLower(text)
	for _, marker := range []string{"unauthorized", "not authorized", "authorization", "authentication", "permission denied", "access denied", "status code: 401", "status code: 403", " 401", " 403"} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

func (c *httpE2EClient) resource(ctx context.Context, uri string, out any) error {
	var result struct {
		Contents []struct {
			URI  string `json:"uri"`
			Text string `json:"text"`
		} `json:"contents"`
	}
	if err := c.rpc(ctx, "resources/read", map[string]any{"uri": uri}, &result); err != nil {
		return err
	}
	if len(result.Contents) != 1 || result.Contents[0].URI != uri || json.Unmarshal([]byte(result.Contents[0].Text), out) != nil {
		return errors.New("invalid broker resource contents")
	}
	return nil
}

func runHTTP(ctx context.Context, cfg config) error {
	if err := checkHealth(ctx, cfg.httpBaseURL); err != nil {
		return errors.New("HTTP health probes failed")
	}
	admin := newHTTPE2EClient(cfg.httpBaseURL, cfg.adminToken)
	user := newHTTPE2EClient(cfg.httpBaseURL, cfg.testUserToken)
	defer admin.client.CloseIdleConnections()
	defer user.client.CloseIdleConnections()

	var discovery struct {
		SupportedVersions []string                   `json:"supportedVersions"`
		ResultType        string                     `json:"resultType"`
		CacheScope        string                     `json:"cacheScope"`
		TTL               *int                       `json:"ttlMs"`
		Meta              map[string]json.RawMessage `json:"_meta"`
	}
	if err := admin.rpc(ctx, "server/discover", nil, &discovery); err != nil {
		return err
	}
	if len(discovery.SupportedVersions) != 1 || discovery.SupportedVersions[0] != httpProtocolVersion || discovery.ResultType != "complete" || discovery.CacheScope != "private" || discovery.TTL == nil || *discovery.TTL != 0 || len(discovery.Meta["io.modelcontextprotocol/serverInfo"]) == 0 {
		return errors.New("unexpected modern discovery contract")
	}

	var catalog struct {
		Tools []struct {
			Name        string `json:"name"`
			Annotations struct {
				ReadOnly *bool `json:"readOnlyHint"`
			} `json:"annotations"`
		} `json:"tools"`
		NextCursor string `json:"nextCursor"`
	}
	names := make(map[string]bool)
	params := map[string]any{}
	for page := 0; ; page++ {
		if page == 100 {
			return errors.New("tools pagination did not terminate")
		}
		catalog.Tools = nil
		if err := admin.rpc(ctx, "tools/list", params, &catalog); err != nil {
			return err
		}
		if catalog.Tools == nil {
			return errors.New("missing tools array")
		}
		for _, tool := range catalog.Tools {
			names[tool.Name] = true
			if cfg.readOnly && (tool.Annotations.ReadOnly == nil || !*tool.Annotations.ReadOnly) {
				return errors.New("read-only catalog contains non-read-only tool")
			}
		}
		if catalog.NextCursor == "" {
			break
		}
		params["cursor"] = catalog.NextCursor
		catalog.NextCursor = ""
	}
	if !names["pulsar_admin_cluster_read"] {
		return errors.New("missing Pulsar read tool")
	}
	if cfg.readOnly && (names["pulsar_client_produce"] || names["pulsar_admin_tenant_write"]) {
		return errors.New("write tool leaked into read-only catalog")
	}

	var resources struct {
		Resources []struct {
			URI string `json:"uri"`
		} `json:"resources"`
		NextCursor string `json:"nextCursor"`
	}
	found := false
	params = map[string]any{}
	for page := 0; ; page++ {
		if page == 100 {
			return errors.New("resources pagination did not terminate")
		}
		resources.Resources = nil
		if err := admin.rpc(ctx, "resources/list", params, &resources); err != nil {
			return err
		}
		if resources.Resources == nil {
			return errors.New("missing resources array")
		}
		for _, resource := range resources.Resources {
			if resource.URI == "pulsar://admin/v2/clusters" {
				found = true
			}
		}
		if resources.NextCursor == "" {
			break
		}
		params["cursor"] = resources.NextCursor
		resources.NextCursor = ""
	}
	if !found {
		return errors.New("broker clusters resource not advertised")
	}
	var clusters struct {
		Clusters []string `json:"clusters"`
	}
	if err := admin.resource(ctx, "pulsar://admin/v2/clusters", &clusters); err != nil {
		return err
	}
	if len(clusters.Clusters) == 0 {
		return errors.New("empty broker cluster resource")
	}

	tenant := fmt.Sprintf("e2e-http-%d", time.Now().UnixNano())
	createArgs := map[string]any{"resource": "tenant", "operation": "create", "tenant": tenant, "adminRoles": []string{"admin"}, "allowedClusters": clusters.Clusters[:1]}
	if err := httpTransportChecks(ctx, admin, createArgs); err != nil {
		return err
	}
	// Real backend read proves all malformed write probes above had no effect.
	var tenants struct {
		Tenants []string `json:"tenants"`
	}
	if err := admin.resource(ctx, "pulsar://admin/v2/tenants", &tenants); err != nil {
		return err
	}
	if tenants.Tenants == nil {
		return errors.New("missing tenants in broker resource")
	}
	for _, existing := range tenants.Tenants {
		if existing == tenant {
			return errors.New("rejected HTTP request created a tenant")
		}
	}

	if cfg.readOnly {
		req, err := admin.request(ctx, "tools/call", map[string]any{"name": "pulsar_admin_tenant_write", "arguments": createArgs})
		if err != nil {
			return err
		}
		if _, err := admin.exchange(req, http.StatusBadRequest, mcp.INVALID_PARAMS); err != nil {
			return fmt.Errorf("read-only write rejection: %w", err)
		}
		tenants.Tenants = nil
		if err := admin.resource(ctx, "pulsar://admin/v2/tenants", &tenants); err != nil {
			return err
		}
		if tenants.Tenants == nil {
			return errors.New("missing tenants after rejected read-only write")
		}
		for _, existing := range tenants.Tenants {
			if existing == tenant {
				return errors.New("read-only call created a tenant")
			}
		}
		logf(cfg.verbose, "modern HTTP read-only catalog and broker resource checks passed")
		return nil
	}
	for _, name := range []string{"pulsar_admin_tenant_write", "pulsar_admin_namespace_write", "pulsar_admin_namespace_policy_set", "pulsar_admin_topic_write", "pulsar_client_produce", "pulsar_client_consume"} {
		if !names[name] {
			return fmt.Errorf("missing required tool %s", name)
		}
	}
	if _, err := admin.tool(ctx, "pulsar_admin_tenant_write", createArgs, false); err != nil {
		return err
	}
	namespace := tenant + "/ns"
	topic := "persistent://" + namespace + "/messages"
	for _, step := range []struct {
		name string
		args map[string]any
	}{
		{"pulsar_admin_namespace_write", map[string]any{"operation": "create", "namespace": namespace, "clusters": clusters.Clusters[:1]}},
		{"pulsar_admin_namespace_policy_set", map[string]any{"namespace": namespace, "policy": "permission", "role": "test-user", "actions": []string{"consume"}}},
		{"pulsar_admin_topic_write", map[string]any{"resource": "topic", "operation": "create", "topic": topic, "partitions": 0}},
	} {
		if _, err := admin.tool(ctx, step.name, step.args, false); err != nil {
			return err
		}
	}

	// Each round races both identities over fresh connections. Verify payloads,
	// not merely a successful (possibly empty) consume result.
	for round := 0; round < 3; round++ {
		payload := fmt.Sprintf("admin-message-%d", round)
		args := map[string]any{"topic": topic, "messages": []string{payload}}
		var wg sync.WaitGroup
		errs := make(chan error, 2)
		for _, identity := range []struct {
			client *httpE2EClient
			denied bool
		}{{admin, false}, {user, true}} {
			wg.Add(1)
			go func(c *httpE2EClient, denied bool) {
				defer wg.Done()
				_, err := c.tool(ctx, "pulsar_client_produce", args, denied)
				errs <- err
			}(identity.client, identity.denied)
		}
		wg.Wait()
		close(errs)
		for err := range errs {
			if err != nil {
				return err
			}
		}
		raw, err := user.tool(ctx, "pulsar_client_consume", map[string]any{
			"topic": topic, "subscription-name": "e2e-http", "initial-position": "earliest",
			"num-messages": 1, "timeout": 15, "subscription-type": "exclusive",
			"subscription-mode": "durable", "hide-payload": false,
		}, false)
		if err != nil {
			return err
		}
		if err := assertHTTPConsumed(raw, payload); err != nil {
			return err
		}
	}
	// Invalid JWTs may construct a client successfully. Only a broker-backed
	// operation can prove authentication failure.
	invalid := newHTTPE2EClient(cfg.httpBaseURL, "invalid-token")
	defer invalid.client.CloseIdleConnections()
	if _, err := invalid.tool(ctx, "pulsar_admin_cluster_read", map[string]any{"resource": "cluster", "operation": "list"}, true); err != nil {
		return err
	}
	logf(cfg.verbose, "modern HTTP broker authorization and repeated-connection identity checks passed")
	return nil
}

func assertHTTPConsumed(raw, payload string) error {
	var consumed consumeResponse
	if json.Unmarshal([]byte(raw), &consumed) != nil || consumed.MessagesConsumed != 1 || len(consumed.Messages) != 1 || consumed.Messages[0].Data != payload {
		return errors.New("consume did not return the expected admin payload")
	}
	return nil
}

func httpTransportChecks(ctx context.Context, c *httpE2EClient, createArgs map[string]any) error {
	for _, tc := range []struct {
		label, header, value string
		status, code         int
	}{
		{"missing credentials", "Authorization", "", 401, 0},
		{"invalid Origin", "Origin", "https://untrusted.invalid", 403, 0},
		{"missing method", "Mcp-Method", "", 400, mcp.HEADER_MISMATCH},
		{"mismatched method", "Mcp-Method", "tools/list", 400, mcp.HEADER_MISMATCH},
		{"missing name", "Mcp-Name", "", 400, mcp.HEADER_MISMATCH},
		{"mismatched name", "Mcp-Name", "pulsar_client_produce", 400, mcp.HEADER_MISMATCH},
		{"missing version", "Mcp-Protocol-Version", "", 400, mcp.HEADER_MISMATCH},
		{"mismatched version", "Mcp-Protocol-Version", "2025-11-25", 400, mcp.HEADER_MISMATCH},
	} {
		req, err := c.request(ctx, "tools/call", map[string]any{"name": "pulsar_admin_tenant_write", "arguments": createArgs})
		if err != nil {
			return err
		}
		if tc.value == "" {
			req.Header.Del(tc.header)
		} else {
			req.Header.Set(tc.header, tc.value)
		}
		if _, err := c.exchange(req, tc.status, tc.code); err != nil {
			return fmt.Errorf("%s: %w", tc.label, err)
		}
	}
	for _, method := range []string{http.MethodGet, http.MethodDelete} {
		req, err := c.request(ctx, "tools/list", nil)
		if err != nil {
			return err
		}
		req.Method, req.Body, req.ContentLength = method, nil, 0
		if _, err := c.exchange(req, http.StatusMethodNotAllowed, 0); err != nil {
			return err
		}
	}
	req, err := c.request(ctx, "initialize", nil)
	if err != nil {
		return err
	}
	if _, err := c.exchange(req, http.StatusNotFound, mcp.METHOD_NOT_FOUND); err != nil {
		return fmt.Errorf("removed initialize: %w", err)
	}
	// Unsupported metadata and header agree: this is not a header mismatch.
	req, err = c.request(ctx, "tools/call", map[string]any{"name": "pulsar_admin_tenant_write", "arguments": createArgs})
	if err != nil {
		return err
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return errors.New("cannot build unsupported-version probe")
	}
	body = bytes.ReplaceAll(body, []byte(httpProtocolVersion), []byte("2099-01-01"))
	req.Body = io.NopCloser(bytes.NewReader(body))
	req.ContentLength = int64(len(body))
	req.Header.Set("Mcp-Protocol-Version", "2099-01-01")
	if _, err := c.exchange(req, http.StatusBadRequest, mcp.UNSUPPORTED_PROTOCOL_VERSION); err != nil {
		return err
	}
	return nil
}
