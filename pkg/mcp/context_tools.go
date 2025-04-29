package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/auth/store"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
)

func RegisterContextTools(s *server.MCPServer, features []string) {
	if !slices.Contains(features, string(FeatureStreamNativeCloud)) && !slices.Contains(features, string(FeatureAll)) {
		return
	}
	// Add whoami tool
	whoamiTool := mcp.NewTool("streamnative_cloud_context_whoami",
		mcp.WithDescription("Display the currently logged-in service account. "+
			"Returns the name of the authenticated service account and the organization."),
	)
	s.AddTool(whoamiTool, handleWhoami)

	// Add set-context tool
	setContextTool := mcp.NewTool("streamnative_cloud_context_use_cluster",
		mcp.WithDescription("Set the current context to a specific StreamNative Cloud cluster, once you set the context, you can use pulsar and kafka tools to interact with the cluster. If you encounter ContextNotSetErr, please use `streamnative_cloud_context_available_clusters` to list the available clusters and set the context to a specific cluster."),
		mcp.WithString("instanceName", mcp.Required(),
			mcp.Description("The name of the pulsar instance to use"),
		),
		mcp.WithString("clusterName", mcp.Required(),
			mcp.Description("The name of the pulsar cluster to use"),
		),
	)
	s.AddTool(setContextTool, handleSetContext)

	// Add available-contexts tool
	availableContextsTool := mcp.NewTool("streamnative_cloud_available_contexts",
		mcp.WithDescription("Display the available pulsar clusters for the current organization on StreamNative Cloud. You can use `streamnative_cloud_context_use_cluster` to change the context to a specific cluster. You will need to ask for the USER to confirm the target context cluster if there are multiple clusters."),
	)
	s.AddTool(availableContextsTool, handleAvailableContexts)
}

// handleWhoami handles the whoami tool request
func handleWhoami(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	options := ctx.Value(OptionsKey).(*config.Options)
	issuer := options.LoadConfigOrDie().Auth.Issuer()

	userName, err := options.WhoAmI(issuer.Audience)
	if err != nil {
		if err == store.ErrNoAuthenticationData {
			return mcp.NewToolResultText("Not logged in."), nil
		}
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get user information: %v", err)), nil
	}

	response := struct {
		ServiceAccount string `json:"service_account"`
		Organization   string `json:"organization"`
	}{
		ServiceAccount: userName,
		Organization:   options.Organization,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonResponse)), nil
}

// handleSetContext handles the set-context tool request
func handleSetContext(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	options := ctx.Value(OptionsKey).(*config.Options)

	instanceName, err := requiredParam[string](request.Params.Arguments, "instanceName")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get instance name: %v", err)), nil
	}

	clusterName, err := requiredParam[string](request.Params.Arguments, "clusterName")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get cluster name: %v", err)), nil
	}

	err = SetContext(options, instanceName, clusterName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set context: %v", err)), nil
	}

	return mcp.NewToolResultText("StreamNative Cloud context set successfully"), nil
}

func handleAvailableContexts(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	promptResponse, err := handleListPulsarClusters(ctx, mcp.GetPromptRequest{})
	if err != nil || promptResponse == nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list pulsar clusters: %v", err)), nil
	}

	response := ""
	for _, message := range promptResponse.Messages {
		if textContent, ok := message.Content.(mcp.TextContent); ok {
			response += textContent.Text + "\n"
		}
	}
	response += "Please confirm the target context cluster with USER if there are multiple clusters!"

	return mcp.NewToolResultText(response), nil
}
