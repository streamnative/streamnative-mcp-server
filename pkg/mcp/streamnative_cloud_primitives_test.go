package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	sncloud "github.com/streamnative/streamnative-mcp-server/sdk/sdk-apiserver"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newTestSNCloudAPIClient(t *testing.T, transport http.RoundTripper) *sncloud.APIClient {
	t.Helper()

	cfg := sncloud.NewConfiguration()
	cfg.HTTPClient = &http.Client{Transport: transport}
	cfg.Servers = sncloud.ServerConfigurations{
		{URL: "http://sncloud.test"},
	}

	return sncloud.NewAPIClient(cfg)
}

func TestSNCloudToolConstructors(t *testing.T) {
	t.Parallel()

	logTool := NewSNCloudLogsTool()
	if logTool.Name != "sncloud_logs" {
		t.Fatalf("expected log tool name sncloud_logs, got %q", logTool.Name)
	}

	applyTool := NewSNCloudResourcesApplyTool()
	if applyTool.Name != "sncloud_resources_apply" {
		t.Fatalf("expected apply tool name sncloud_resources_apply, got %q", applyTool.Name)
	}

	deleteTool := NewSNCloudResourcesDeleteTool()
	if deleteTool.Name != "sncloud_resources_delete" {
		t.Fatalf("expected delete tool name sncloud_resources_delete, got %q", deleteTool.Name)
	}
}

func TestSNCloudPromptConstructors(t *testing.T) {
	t.Parallel()

	listPrompt := NewListSNCloudClustersPrompt()
	if listPrompt.Name != "list-sncloud-clusters" {
		t.Fatalf("expected list prompt name list-sncloud-clusters, got %q", listPrompt.Name)
	}
	if len(listPrompt.Arguments) != 0 {
		t.Fatalf("expected list prompt to have no arguments, got %d", len(listPrompt.Arguments))
	}

	readPrompt := NewReadSNCloudClusterPrompt()
	if readPrompt.Name != "read-sncloud-cluster" {
		t.Fatalf("expected read prompt name read-sncloud-cluster, got %q", readPrompt.Name)
	}
	if len(readPrompt.Arguments) != 1 || readPrompt.Arguments[0].Name != "name" || !readPrompt.Arguments[0].Required {
		t.Fatalf("unexpected read prompt arguments: %+v", readPrompt.Arguments)
	}

	buildPrompt := NewBuildSNCloudServerlessClusterPrompt()
	if buildPrompt.Name != "build-sncloud-serverless-cluster" {
		t.Fatalf("expected build prompt name build-sncloud-serverless-cluster, got %q", buildPrompt.Name)
	}
	if len(buildPrompt.Arguments) != 3 {
		t.Fatalf("expected build prompt to have 3 arguments, got %d", len(buildPrompt.Arguments))
	}
	if buildPrompt.Arguments[0].Name != "instance-name" || !buildPrompt.Arguments[0].Required {
		t.Fatalf("unexpected first build prompt argument: %+v", buildPrompt.Arguments[0])
	}
	if buildPrompt.Arguments[1].Name != "cluster-name" || !buildPrompt.Arguments[1].Required {
		t.Fatalf("unexpected second build prompt argument: %+v", buildPrompt.Arguments[1])
	}
	if buildPrompt.Arguments[2].Name != "provider" || buildPrompt.Arguments[2].Required {
		t.Fatalf("unexpected third build prompt argument: %+v", buildPrompt.Arguments[2])
	}
}

func TestResolveSNCloudOrganizationPrefersSession(t *testing.T) {
	t.Parallel()

	session, err := config.NewSNCloudSession(config.SNCloudContext{
		JWTToken:     "token",
		APIURL:       "https://api.example.com",
		LogAPIURL:    "https://logs.example.com",
		Organization: "session-org",
	})
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	ctx := context.Background()
	ctx = WithSNCloudSession(ctx, session)
	ctx = context.WithValue(ctx, common.OptionsKey, &config.Options{Organization: "options-org"})

	if got := resolveSNCloudOrganization(ctx); got != "session-org" {
		t.Fatalf("expected session organization, got %q", got)
	}
}

func TestResolveSNCloudOrganizationFallsBackToOptions(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), common.OptionsKey, &config.Options{Organization: "options-org"})
	if got := resolveSNCloudOrganization(ctx); got != "options-org" {
		t.Fatalf("expected options organization, got %q", got)
	}
}

func TestHandleSNCloudResourcesApplyUsesSessionOrganization(t *testing.T) {
	t.Parallel()

	var requestedPaths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPaths = append(requestedPaths, r.URL.Path)

		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"metadata":{"name":"pi-test","resourceVersion":"1"}}`))
		case http.MethodPut:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"metadata":{"name":"pi-test","resourceVersion":"2"}}`))
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	}))
	defer server.Close()

	session, err := config.NewSNCloudSession(config.SNCloudContext{
		JWTToken:     "token",
		APIURL:       server.URL,
		LogAPIURL:    server.URL,
		Organization: "session-org",
	})
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	ctx := context.Background()
	ctx = WithSNCloudSession(ctx, session)
	ctx = context.WithValue(ctx, common.OptionsKey, &config.Options{Organization: "options-org"})

	result, err := HandleSNCloudResourcesApply(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"json_content": `{"apiVersion":"cloud.streamnative.io/v1alpha1","kind":"PulsarInstance","metadata":{"name":"pi-test"}}`,
			},
		},
	})
	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if result.IsError {
		t.Fatalf("expected successful result, got error result: %+v", result)
	}

	if len(requestedPaths) != 2 {
		t.Fatalf("expected 2 upstream requests, got %d", len(requestedPaths))
	}
	for _, path := range requestedPaths {
		if !strings.Contains(path, "/session-org/") {
			t.Fatalf("expected session organization in path, got %q", path)
		}
		if strings.Contains(path, "/options-org/") {
			t.Fatalf("expected options organization to be ignored, got %q", path)
		}
	}
}

func TestHandleSNCloudResourcesApplySupportsInstanceCreate(t *testing.T) {
	t.Parallel()

	var (
		postBody  map[string]any
		postQuery string
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/apis/cloud.streamnative.io/v1alpha1/namespaces/session-org/instances/inst-create":
			http.NotFound(w, r)
		case r.Method == http.MethodPost && r.URL.Path == "/apis/cloud.streamnative.io/v1alpha1/namespaces/session-org/instances":
			postQuery = r.URL.RawQuery
			defer func() { _ = r.Body.Close() }()
			if err := json.NewDecoder(r.Body).Decode(&postBody); err != nil {
				t.Fatalf("failed to decode POST body: %v", err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"metadata":{"name":"inst-create"}}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	}))
	defer server.Close()

	session, err := config.NewSNCloudSession(config.SNCloudContext{
		JWTToken:     "token",
		APIURL:       server.URL,
		LogAPIURL:    server.URL,
		Organization: "session-org",
	})
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	ctx := context.Background()
	ctx = WithSNCloudSession(ctx, session)

	result, err := HandleSNCloudResourcesApply(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"json_content": `{"apiVersion":"cloud.streamnative.io/v1alpha1","kind":"Instance","metadata":{"name":"inst-create"},"spec":{"type":"serverless"}}`,
			},
		},
	})
	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if result.IsError {
		t.Fatalf("expected successful result, got error result: %+v", result)
	}
	if postQuery != "" {
		t.Fatalf("expected create without dryRun query, got %q", postQuery)
	}

	metadata, ok := postBody["metadata"].(map[string]any)
	if !ok {
		t.Fatalf("expected metadata map in create body, got %+v", postBody["metadata"])
	}
	if got := metadata["namespace"]; got != "session-org" {
		t.Fatalf("expected session organization namespace, got %#v", got)
	}
}

func TestHandleSNCloudResourcesApplySupportsKafkaClusterUpdate(t *testing.T) {
	t.Parallel()

	var (
		putBody  map[string]any
		putQuery string
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/apis/cloud.streamnative.io/v1alpha1/namespaces/session-org/kafkaclusters/kc-test":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"apiVersion":"cloud.streamnative.io/v1alpha1","kind":"KafkaCluster","metadata":{"name":"kc-test","resourceVersion":"7"}}`))
		case r.Method == http.MethodPut && r.URL.Path == "/apis/cloud.streamnative.io/v1alpha1/namespaces/session-org/kafkaclusters/kc-test":
			putQuery = r.URL.RawQuery
			defer func() { _ = r.Body.Close() }()
			if err := json.NewDecoder(r.Body).Decode(&putBody); err != nil {
				t.Fatalf("failed to decode PUT body: %v", err)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"metadata":{"name":"kc-test","resourceVersion":"8"}}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	}))
	defer server.Close()

	session, err := config.NewSNCloudSession(config.SNCloudContext{
		JWTToken:     "token",
		APIURL:       server.URL,
		LogAPIURL:    server.URL,
		Organization: "session-org",
	})
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	ctx := context.Background()
	ctx = WithSNCloudSession(ctx, session)

	result, err := HandleSNCloudResourcesApply(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"json_content": `{"apiVersion":"cloud.streamnative.io/v1alpha1","kind":"KafkaCluster","metadata":{"name":"kc-test"},"spec":{"instanceName":"inst-1"}}`,
			},
		},
	})
	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if result.IsError {
		t.Fatalf("expected successful result, got error result: %+v", result)
	}

	if putQuery != "" {
		t.Fatalf("expected update without dryRun query, got %q", putQuery)
	}

	metadata, ok := putBody["metadata"].(map[string]any)
	if !ok {
		t.Fatalf("expected metadata map in update body, got %+v", putBody["metadata"])
	}
	if got := metadata["resourceVersion"]; got != "7" {
		t.Fatalf("expected propagated resourceVersion 7, got %#v", got)
	}
	if got := metadata["namespace"]; got != "session-org" {
		t.Fatalf("expected session organization namespace, got %#v", got)
	}
}

func TestHandleSNCloudResourcesApplySupportsKafkaClusterDryRunCreate(t *testing.T) {
	t.Parallel()

	var (
		postBody  map[string]any
		postQuery string
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/apis/cloud.streamnative.io/v1alpha1/namespaces/session-org/kafkaclusters/kc-dry-run":
			http.NotFound(w, r)
		case r.Method == http.MethodPost && r.URL.Path == "/apis/cloud.streamnative.io/v1alpha1/namespaces/session-org/kafkaclusters":
			postQuery = r.URL.RawQuery
			defer func() { _ = r.Body.Close() }()
			if err := json.NewDecoder(r.Body).Decode(&postBody); err != nil {
				t.Fatalf("failed to decode POST body: %v", err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"metadata":{"name":"kc-dry-run"}}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	}))
	defer server.Close()

	session, err := config.NewSNCloudSession(config.SNCloudContext{
		JWTToken:     "token",
		APIURL:       server.URL,
		LogAPIURL:    server.URL,
		Organization: "session-org",
	})
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	ctx := context.Background()
	ctx = WithSNCloudSession(ctx, session)

	result, err := HandleSNCloudResourcesApply(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"json_content": `{"apiVersion":"cloud.streamnative.io/v1alpha1","kind":"KafkaCluster","metadata":{"name":"kc-dry-run"},"spec":{"instanceName":"inst-1"}}`,
				"dry_run":      true,
			},
		},
	})
	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if result.IsError {
		t.Fatalf("expected successful result, got error result: %+v", result)
	}
	if !strings.Contains(postQuery, "dryRun=All") {
		t.Fatalf("expected dryRun query parameter, got %q", postQuery)
	}

	metadata, ok := postBody["metadata"].(map[string]any)
	if !ok {
		t.Fatalf("expected metadata map in create body, got %+v", postBody["metadata"])
	}
	if got := metadata["namespace"]; got != "session-org" {
		t.Fatalf("expected session organization namespace, got %#v", got)
	}
}

func TestApplyPulsarInstanceHandlesNilResponseError(t *testing.T) {
	t.Parallel()

	requestCount := 0
	apiClient := newTestSNCloudAPIClient(t, roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		requestCount++
		return nil, errors.New("upstream transport failure")
	}))

	_, err := applyPulsarInstance(
		context.Background(),
		apiClient,
		`{"apiVersion":"cloud.streamnative.io/v1alpha1","kind":"PulsarInstance","metadata":{"name":"pi-nil-response"}}`,
		"session-org",
		false,
	)
	if err == nil {
		t.Fatal("expected nil-response transport error")
	}
	if !strings.Contains(err.Error(), "failed to created PulsarInstance:") {
		t.Fatalf("expected wrapped PulsarInstance error prefix, got %v", err)
	}
	if !strings.Contains(err.Error(), "upstream transport failure") {
		t.Fatalf("expected wrapped PulsarInstance error, got %v", err)
	}
	if requestCount != 2 {
		t.Fatalf("expected 2 upstream requests, got %d", requestCount)
	}
}

func TestApplyPulsarClusterHandlesNilResponseError(t *testing.T) {
	t.Parallel()

	requestCount := 0
	apiClient := newTestSNCloudAPIClient(t, roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		requestCount++
		return nil, errors.New("upstream transport failure")
	}))

	_, err := applyPulsarCluster(
		context.Background(),
		apiClient,
		`{"apiVersion":"cloud.streamnative.io/v1alpha1","kind":"PulsarCluster","metadata":{"name":"pc-nil-response"}}`,
		"session-org",
		false,
	)
	if err == nil {
		t.Fatal("expected nil-response transport error")
	}
	if !strings.Contains(err.Error(), "failed to created PulsarCluster:") {
		t.Fatalf("expected wrapped PulsarCluster error prefix, got %v", err)
	}
	if !strings.Contains(err.Error(), "upstream transport failure") {
		t.Fatalf("expected wrapped PulsarCluster error, got %v", err)
	}
	if requestCount != 2 {
		t.Fatalf("expected 2 upstream requests, got %d", requestCount)
	}
}

func TestHandleSNCloudResourcesDeleteSupportsKafkaCluster(t *testing.T) {
	t.Parallel()

	var deletePath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method %s", r.Method)
		}
		deletePath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"Success"}`))
	}))
	defer server.Close()

	session, err := config.NewSNCloudSession(config.SNCloudContext{
		JWTToken:     "token",
		APIURL:       server.URL,
		LogAPIURL:    server.URL,
		Organization: "session-org",
	})
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	ctx := context.Background()
	ctx = WithSNCloudSession(ctx, session)

	result, err := HandleSNCloudResourcesDelete(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"name": "kc-test",
				"type": "KafkaCluster",
			},
		},
	})
	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if result.IsError {
		t.Fatalf("expected successful result, got error result: %+v", result)
	}

	expectedPath := "/apis/cloud.streamnative.io/v1alpha1/namespaces/session-org/kafkaclusters/kc-test"
	if deletePath != expectedPath {
		t.Fatalf("expected delete path %q, got %q", expectedPath, deletePath)
	}
}

func TestHandleSNCloudResourcesDeleteSupportsInstance(t *testing.T) {
	t.Parallel()

	var deletePath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method %s", r.Method)
		}
		deletePath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"Success"}`))
	}))
	defer server.Close()

	session, err := config.NewSNCloudSession(config.SNCloudContext{
		JWTToken:     "token",
		APIURL:       server.URL,
		LogAPIURL:    server.URL,
		Organization: "session-org",
	})
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	ctx := context.Background()
	ctx = WithSNCloudSession(ctx, session)

	result, err := HandleSNCloudResourcesDelete(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"name": "inst-delete",
				"type": "Instance",
			},
		},
	})
	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if result.IsError {
		t.Fatalf("expected successful result, got error result: %+v", result)
	}

	expectedPath := "/apis/cloud.streamnative.io/v1alpha1/namespaces/session-org/instances/inst-delete"
	if deletePath != expectedPath {
		t.Fatalf("expected delete path %q, got %q", expectedPath, deletePath)
	}
}

func TestHandleReadSNCloudClusterFallsBackToKafkaCluster(t *testing.T) {
	t.Parallel()

	var pulsarReadRequests, kafkaReadRequests int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/apis/cloud.streamnative.io/v1alpha1/namespaces/session-org/pulsarclusters/kc-test":
			pulsarReadRequests++
			http.Error(w, `{"message":"not found"}`, http.StatusNotFound)
		case "/apis/cloud.streamnative.io/v1alpha1/namespaces/session-org/kafkaclusters/kc-test":
			kafkaReadRequests++
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"apiVersion":"cloud.streamnative.io/v1alpha1","kind":"KafkaCluster","metadata":{"name":"kc-test","managedFields":[{"manager":"controller"}]},"spec":{"instanceName":"inst-1"}}`))
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	session, err := config.NewSNCloudSession(config.SNCloudContext{
		JWTToken:     "token",
		APIURL:       server.URL,
		LogAPIURL:    server.URL,
		Organization: "session-org",
	})
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	ctx := context.Background()
	ctx = WithSNCloudSession(ctx, session)

	result, err := HandleReadSNCloudCluster(ctx, mcp.GetPromptRequest{
		Params: mcp.GetPromptParams{
			Arguments: map[string]string{"name": "kc-test"},
		},
	})
	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if len(result.Messages) != 1 {
		t.Fatalf("expected one prompt message, got %d", len(result.Messages))
	}

	text, ok := result.Messages[0].Content.(mcp.TextContent)
	if !ok {
		t.Fatalf("expected text content, got %T", result.Messages[0].Content)
	}
	if !strings.Contains(text.Text, `"kind":"KafkaCluster"`) {
		t.Fatalf("expected KafkaCluster payload, got %q", text.Text)
	}
	if strings.Contains(text.Text, "managedFields") {
		t.Fatalf("expected managedFields to be removed, got %q", text.Text)
	}
	if pulsarReadRequests != 1 {
		t.Fatalf("expected one pulsar read request, got %d", pulsarReadRequests)
	}
	if kafkaReadRequests != 1 {
		t.Fatalf("expected one kafka read request, got %d", kafkaReadRequests)
	}
}

func TestHandleReadSNCloudClusterDoesNotFallbackOnPulsarReadError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/apis/cloud.streamnative.io/v1alpha1/namespaces/session-org/pulsarclusters/pc-test":
			http.Error(w, `{"message":"upstream error"}`, http.StatusInternalServerError)
		case "/apis/cloud.streamnative.io/v1alpha1/namespaces/session-org/kafkaclusters/pc-test":
			t.Fatalf("unexpected kafka fallback for pulsar read error")
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	session, err := config.NewSNCloudSession(config.SNCloudContext{
		JWTToken:     "token",
		APIURL:       server.URL,
		LogAPIURL:    server.URL,
		Organization: "session-org",
	})
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	ctx := context.Background()
	ctx = WithSNCloudSession(ctx, session)

	_, err = HandleReadSNCloudCluster(ctx, mcp.GetPromptRequest{
		Params: mcp.GetPromptParams{
			Arguments: map[string]string{"name": "pc-test"},
		},
	})
	if err == nil {
		t.Fatal("expected handler error")
	}
	if !strings.Contains(err.Error(), "failed to read pulsar cluster") {
		t.Fatalf("expected pulsar read error, got %v", err)
	}
}

func TestHandleReadSNCloudClusterReturnsNotFoundWhenClusterDoesNotExist(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/apis/cloud.streamnative.io/v1alpha1/namespaces/session-org/pulsarclusters/missing",
			"/apis/cloud.streamnative.io/v1alpha1/namespaces/session-org/kafkaclusters/missing":
			http.Error(w, `{"message":"not found"}`, http.StatusNotFound)
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	session, err := config.NewSNCloudSession(config.SNCloudContext{
		JWTToken:     "token",
		APIURL:       server.URL,
		LogAPIURL:    server.URL,
		Organization: "session-org",
	})
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	ctx := context.Background()
	ctx = WithSNCloudSession(ctx, session)

	_, err = HandleReadSNCloudCluster(ctx, mcp.GetPromptRequest{
		Params: mcp.GetPromptParams{
			Arguments: map[string]string{"name": "missing"},
		},
	})
	if err == nil {
		t.Fatal("expected handler error")
	}
	if err.Error() != "failed to find cluster: missing" {
		t.Fatalf("expected not found error, got %v", err)
	}
}
