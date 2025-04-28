package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/auth/store"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
	sncloud "github.com/streamnative/streamnative-mcp-server/sdk/sdk-apiserve"
)

func RegisterContextTools(s *server.MCPServer) {
	// Add whoami tool
	whoamiTool := mcp.NewTool("streamnative_cloud_context_whoami",
		mcp.WithDescription("Display the currently logged-in service account. "+
			"Returns the name of the authenticated service account and the organization."),
	)
	s.AddTool(whoamiTool, handleWhoami)

	// Add set-context tool
	setContextTool := mcp.NewTool("streamnative_cloud_use_pulsar_cluster",
		mcp.WithDescription("Set the current context to a specific pulsar cluster, you can use pulsar and kafka tools to interact with the cluster."),
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
		mcp.WithDescription("Display the available pulsar clusters for the current organization on StreamNative Cloud. You will need to ask for the USER to confirm the target context cluster if there are multiple clusters."),
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
	snConfig := options.LoadConfigOrDie()
	myselfGrant, err := options.AuthOptions.LoadGrant(snConfig.Auth.Audience)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to load myself grant: %v", err)), nil
	}

	instanceName, err := requiredParam[string](request.Params.Arguments, "instanceName")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get instance name: %v", err)), nil
	}

	clusterName, err := requiredParam[string](request.Params.Arguments, "clusterName")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get cluster name: %v", err)), nil
	}

	apiClient, err := config.GetAPIClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get API client: %v", err)), nil
	}

	instances, _, err := apiClient.CloudStreamnativeIoV1alpha1Api.ListCloudStreamnativeIoV1alpha1NamespacedPulsarInstance(ctx, options.Organization).Execute()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list pulsar instances: %v", err)), nil
	}

	var instance sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstance
	foundInstance := false
	for _, i := range instances.Items {
		if *i.Metadata.Name == instanceName {
			if isInstanceValid(i) {
				instance = i
				foundInstance = true
				break
			} else {
				return mcp.NewToolResultError(fmt.Sprintf("Pulsar instance %s is not valid", instanceName)), nil
			}
		}
	}
	if !foundInstance {
		return mcp.NewToolResultError(fmt.Sprintf("Pulsar instance %s not found", instanceName)), nil
	}

	clusters, _, err := apiClient.CloudStreamnativeIoV1alpha1Api.ListCloudStreamnativeIoV1alpha1NamespacedPulsarCluster(ctx, options.Organization).Execute()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list pulsar clusters: %v", err)), nil
	}

	var cluster sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarCluster
	foundCluster := false
	for _, c := range clusters.Items {
		if *c.Metadata.Name == clusterName && c.Spec.InstanceName == instanceName {
			if isClusterAvailable(c) {
				cluster = c
				foundCluster = true
				break
			} else {
				return mcp.NewToolResultError(fmt.Sprintf("Pulsar cluster %s is not available", clusterName)), nil
			}
		}
	}
	if !foundCluster {
		return mcp.NewToolResultError(fmt.Sprintf("Pulsar cluster %s not found", clusterName)), nil
	}

	clusterUID := string(*cluster.Metadata.Uid)
	dnsName := ""
	for _, endpoint := range cluster.Spec.ServiceEndpoints {
		if *endpoint.Type == "service" {
			dnsName = endpoint.DnsName
			break
		}
	}

	if dnsName == "" {
		return mcp.NewToolResultError(fmt.Sprintf("no valid service endpoint found for PulsarCluster '%s'", clusterName)), nil
	}

	issuer, err := getIssuer(&instance, snConfig.Auth.Issuer())
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get issuer: %v", err)), nil
	}

	tokenKey := buildTokenKey(options.Organization, clusterUID)

	accessToken := ""
	refreshToken := true
	cachedGrant, err := options.AuthOptions.LoadGrant(tokenKey)
	if err == nil && cachedGrant != nil {

		cacheHasValidToken, err := hasCachedValidToken(cachedGrant)
		if err != nil {
			cacheHasValidToken = false
		}

		if cacheHasValidToken {
			tokenAboutToExpire, err := isTokenAboutToExpire(cachedGrant, TokenRefreshWindow)
			if err != nil {
				tokenAboutToExpire = true
			}

			if !tokenAboutToExpire {
				refreshToken = false
				accessToken = cachedGrant.Token.AccessToken
			}
		}
	}

	if refreshToken {
		flow, err := getFlow(issuer, myselfGrant)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get flow: %v", err)), nil
		}

		newGrant, err := flow.Authorize()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to authorize: %v", err)), nil
		}

		if newGrant.Token != nil {
			_ = options.AuthOptions.SaveGrant(tokenKey, *newGrant)
			accessToken = newGrant.Token.AccessToken
		}
	}

	if accessToken == "" {
		return mcp.NewToolResultError("Failed to get access token"), nil
	}

	err = pulsar.NewCurrentPulsarContext(pulsar.PulsarContext{
		WebServiceURL: getBasePath(snConfig.ProxyLocation, options.Organization, clusterUID),
		Token:         accessToken,
	})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set pulsar context: %v", err)), nil
	}

	return mcp.NewToolResultText("Pulsar context set successfully"), nil
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
