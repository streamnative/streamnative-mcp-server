package mcp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
)

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
