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

	"github.com/google/jsonschema-go/jsonschema"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/auth/store"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
)

type sncloudContextEmptyInput struct{}

type sncloudContextUseClusterInput struct {
	InstanceName string `json:"instanceName" jsonschema:"The name of the pulsar instance to use"`
	ClusterName  string `json:"clusterName" jsonschema:"The name of the pulsar cluster to use"`
}

// RegisterContextTools registers context-related tools on the server.
func RegisterContextTools(s *sdk.Server, features []string, skipContextTools bool) {
	if !slices.Contains(features, string(FeatureStreamNativeCloud)) && !slices.Contains(features, string(FeatureAll)) {
		return
	}

	emptySchema := &jsonschema.Schema{Type: "object"}

	// Add whoami tool
	whoamiTool := &sdk.Tool{
		Name:        "sncloud_context_whoami",
		Description: "Display the currently logged-in service account. Returns the name of the authenticated service account and the organization.",
		InputSchema: emptySchema,
	}
	sdk.AddTool(s, whoamiTool, handleWhoami)

	// Add set-context tool
	setContextSchema, err := InputSchema[sncloudContextUseClusterInput]()
	if err != nil {
		return
	}
	setContextTool := &sdk.Tool{
		Name:        "sncloud_context_use_cluster",
		Description: "Set the current context to a specific StreamNative Cloud cluster, once you set the context, you can use pulsar and kafka tools to interact with the cluster. If you encounter ContextNotSetErr, please use `sncloud_context_available_clusters` to list the available clusters and set the context to a specific cluster.",
		InputSchema: setContextSchema,
	}
	// Skip registering context tools if context is already provided
	if !skipContextTools {
		sdk.AddTool(s, setContextTool, handleSetContext)
	}

	// Add available-contexts tool
	availableContextsTool := &sdk.Tool{
		Name:        "sncloud_context_available_clusters",
		Description: "Display the available pulsar clusters for the current organization on StreamNative Cloud. You can use `sncloud_context_use_cluster` to change the context to a specific cluster. You will need to ask for the USER to confirm the target context cluster if there are multiple clusters.",
		InputSchema: emptySchema,
	}
	sdk.AddTool(s, availableContextsTool, handleAvailableContexts)
}

// handleWhoami handles the whoami tool request
func handleWhoami(ctx context.Context, _ *sdk.CallToolRequest, _ sncloudContextEmptyInput) (*sdk.CallToolResult, any, error) {
	options := common.GetOptions(ctx)
	issuer := options.LoadConfigOrDie().Auth.Issuer()

	userName, err := options.WhoAmI(issuer.Audience)
	if err != nil {
		if err == store.ErrNoAuthenticationData {
			return newToolResultText("Not logged in."), nil, nil
		}
		return nil, nil, fmt.Errorf("failed to get user information: %w", err)
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
		return nil, nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return newToolResultText(string(jsonResponse)), nil, nil
}

// handleSetContext handles the set-context tool request
func handleSetContext(ctx context.Context, _ *sdk.CallToolRequest, input sncloudContextUseClusterInput) (*sdk.CallToolResult, any, error) {
	options := common.GetOptions(ctx)

	instanceName := input.InstanceName
	if instanceName == "" {
		return nil, nil, fmt.Errorf("missing required parameter 'instanceName'")
	}

	clusterName := input.ClusterName
	if clusterName == "" {
		return nil, nil, fmt.Errorf("missing required parameter 'clusterName'")
	}

	err := SetContext(ctx, options, instanceName, clusterName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to set context: %w", err)
	}

	return newToolResultText("StreamNative Cloud context set successfully"), nil, nil
}

func handleAvailableContexts(ctx context.Context, _ *sdk.CallToolRequest, _ sncloudContextEmptyInput) (*sdk.CallToolResult, any, error) {
	promptResponse, err := HandleListPulsarClusters(ctx, &sdk.GetPromptRequest{})
	if err != nil || promptResponse == nil {
		return nil, nil, fmt.Errorf("failed to list pulsar clusters: %w", err)
	}

	response := ""
	for _, message := range promptResponse.Messages {
		if textContent, ok := message.Content.(*sdk.TextContent); ok {
			response += textContent.Text + "\n"
		}
	}
	response += "Please confirm the target context cluster with USER if there are multiple clusters!"

	return newToolResultText(response), nil, nil
}
