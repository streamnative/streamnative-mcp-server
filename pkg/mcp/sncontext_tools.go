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
	"encoding/json"
	"fmt"
	"slices"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/auth/store"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/adapter"
)

func RegisterContextTools(s *MCPServer, features []string, skipContextTools bool) {
	if !slices.Contains(features, string(FeatureStreamNativeCloud)) && !slices.Contains(features, string(FeatureAll)) {
		return
	}

	// Add whoami tool
	whoamiTool := builders.NewTool("sncloud_context_whoami",
		builders.WithDescription("Display the currently logged-in service account. "+
			"Returns the name of the authenticated service account and the organization."),
	)
	s.AddTool(whoamiTool, handleWhoami)

	// Add set-context tool
	setContextTool := builders.NewTool("sncloud_context_use_cluster",
		builders.WithDescription("Set the current context to a specific StreamNative Cloud cluster, once you set the context, you can use pulsar and kafka tools to interact with the cluster. If you encounter ContextNotSetErr, please use `sncloud_context_available_clusters` to list the available clusters and set the context to a specific cluster."),
		builders.WithString("instanceName", builders.Required(),
			builders.Description("The name of the pulsar instance to use"),
		),
		builders.WithString("clusterName", builders.Required(),
			builders.Description("The name of the pulsar cluster to use"),
		),
	)
	// Skip registering context tools if context is already provided
	if !skipContextTools {
		s.AddTool(setContextTool, handleSetContext)
	}

	// Add available-contexts tool
	availableContextsTool := builders.NewTool("sncloud_context_available_clusters",
		builders.WithDescription("Display the available pulsar clusters for the current organization on StreamNative Cloud. You can use `sncloud_context_use_cluster` to change the context to a specific cluster. You will need to ask for the USER to confirm the target context cluster if there are multiple clusters."),
	)
	s.AddTool(availableContextsTool, handleAvailableContexts)
}

// handleWhoami handles the whoami tool request
func handleWhoami(ctx context.Context, _ *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	options := common.GetOptions(ctx)
	issuer := options.LoadConfigOrDie().Auth.Issuer()

	userName, err := options.WhoAmI(issuer.Audience)
	if err != nil {
		if err == store.ErrNoAuthenticationData {
			return adapter.NewTextResult("Not logged in."), nil
		}
		return adapter.NewErrorResult("Failed to get user information: %v", err), nil
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
		return adapter.NewErrorResult("Failed to marshal response: %v", err), nil
	}

	return adapter.NewTextResult(string(jsonResponse)), nil
}

// handleSetContext handles the set-context tool request
func handleSetContext(ctx context.Context, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	options := common.GetOptions(ctx)

	instanceName, err := adapter.RequireString(request, "instanceName")
	if err != nil {
		return adapter.NewErrorResult("Failed to get instance name: %v", err), nil
	}

	clusterName, err := adapter.RequireString(request, "clusterName")
	if err != nil {
		return adapter.NewErrorResult("Failed to get cluster name: %v", err), nil
	}

	err = SetContext(ctx, options, instanceName, clusterName)
	if err != nil {
		return adapter.NewErrorResult("Failed to set context: %v", err), nil
	}

	return adapter.NewTextResult("StreamNative Cloud context set successfully"), nil
}

func handleAvailableContexts(ctx context.Context, _ *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	promptResponse, err := HandleListPulsarClusters(ctx, &mcpsdk.GetPromptRequest{})
	if err != nil || promptResponse == nil {
		return adapter.NewErrorResult("Failed to list pulsar clusters: %v", err), nil
	}

	response := ""
	for _, message := range promptResponse.Messages {
		if textContent, ok := message.Content.(*mcpsdk.TextContent); ok {
			response += textContent.Text + "\n"
		}
	}
	response += "Please confirm the target context cluster with USER if there are multiple clusters!"

	return adapter.NewTextResult(response), nil
}

// ContextNotSetErr is returned when the context is not set
var ContextNotSetErr = fmt.Errorf("context not set. Please use sncloud_context_use_cluster to set the context first")
