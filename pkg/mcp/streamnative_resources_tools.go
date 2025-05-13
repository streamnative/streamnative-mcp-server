package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	sncloud "github.com/streamnative/streamnative-mcp-server/sdk/sdk-apiserver"
)

// StreamNativeAddResourceTools adds StreamNative resources tools
func StreamNativeAddResourceTools(s *server.MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeatureStreamNativeCloud)) && !slices.Contains(features, string(FeatureAll)) {
		return
	}

	if !readOnly {
		// Add Apply tool
		applyTool := mcp.NewTool("sncloud_resources_apply",
			mcp.WithDescription("Apply StreamNative Cloud resources from JSON definitions. This tool allows you to apply (create or update) StreamNative Cloud resources such as PulsarInstances and PulsarClusters using JSON definitions. Please give feedback to USER if the resource is applied with error, and ask USER to check the resource definition."),
			mcp.WithString("json_content", mcp.Required(),
				mcp.Description("The JSON content to apply."),
			),
			mcp.WithBoolean("dry_run",
				mcp.Description("If true, only validate the resource without applying it to the server."),
				mcp.DefaultBool(false),
			),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title: "Apply StreamNative Cloud Resources",
			}),
		)
		// Add delete tool
		deleteTool := mcp.NewTool("sncloud_resources_delete",
			mcp.WithDescription("Delete StreamNative Cloud resources. This tool allows you to delete StreamNative Cloud resources such as PulsarInstances and PulsarClusters."),
			mcp.WithString("name", mcp.Required(),
				mcp.Description("The name of the resource to delete."),
			),
			mcp.WithString("type", mcp.Required(),
				mcp.Description("The type of the resource to delete, it can be PulsarInstance or PulsarCluster."),
				mcp.Enum("PulsarInstance", "PulsarCluster"),
			),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:           "Delete StreamNative Cloud Resources",
				DestructiveHint: &[]bool{true}[0],
			}),
		)
		s.AddTool(applyTool, handleStreamNativeResourcesApply)
		s.AddTool(deleteTool, handleStreamNativeResourcesDelete)
	}
}

// Define simple resource structure for parsing YAML documents
type Resource struct {
	APIVersion string                 `json:"apiVersion" yaml:"apiVersion"`
	Kind       string                 `json:"kind" yaml:"kind"`
	Metadata   Metadata               `json:"metadata" yaml:"metadata"`
	Spec       map[string]interface{} `json:"spec" yaml:"spec"`
}

type Metadata struct {
	Name      string            `json:"name" yaml:"name"`
	Namespace string            `json:"namespace" yaml:"namespace"`
	Labels    map[string]string `json:"labels" yaml:"labels"`
}

// handleStreamNativeResourcesApply handles the streaming_cloud_resources_apply tool
func handleStreamNativeResourcesApply(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get necessary parameters
	snConfig := getOptions(ctx)
	organization := snConfig.Organization
	if organization == "" {
		return mcp.NewToolResultError("No organization is set. Please set the organization using the appropriate context tool."), nil
	}

	// Get YAML content
	jsonContent, err := requiredParam[string](request.Params.Arguments, "json_content")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get json_content: %v", err)), nil
	}

	// Get dry_run flag
	dryRun, hasDryRun := optionalParam[bool](request.Params.Arguments, "dry_run")
	if !hasDryRun {
		dryRun = false
	}

	// Get API client
	apiClient, err := config.GetAPIClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get API client: %v", err)), nil
	}

	jsonContent = strings.TrimSpace(jsonContent)
	if jsonContent == "" {
		return mcp.NewToolResultError("No valid resources found in the provided JSON."), nil
	}

	// Parse YAML document
	var resource Resource
	err = json.Unmarshal([]byte(jsonContent), &resource)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to parse JSON document: %v", err)), nil
	}

	// Check if resource is valid
	if resource.APIVersion == "" || resource.Kind == "" {
		return mcp.NewToolResultError("Invalid resource definition."), nil
	}

	// Set namespace if not specified
	if resource.Metadata.Namespace == "" {
		resource.Metadata.Namespace = organization
	}

	// Apply resource
	result, err := applyResource(ctx, apiClient, resource, jsonContent, organization, dryRun)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to apply resource: %v", err)), nil
	}

	return mcp.NewToolResultText(result), nil
}

// applyResource applies the resource based on its type
func applyResource(ctx context.Context, apiClient *sncloud.APIClient, resource Resource, jsonContent string, organization string, dryRun bool) (string, error) {
	apiVersion := resource.APIVersion
	kind := resource.Kind

	// Call different APIs based on resource type
	switch {
	case apiVersion == "cloud.streamnative.io/v1alpha1" && kind == "PulsarInstance":
		return applyPulsarInstance(ctx, apiClient, jsonContent, organization, dryRun)
	case apiVersion == "cloud.streamnative.io/v1alpha1" && kind == "PulsarCluster":
		return applyPulsarCluster(ctx, apiClient, jsonContent, organization, dryRun)
	// Can add handling for more resource types
	default:
		return "", fmt.Errorf("unsupported resource type: %s/%s", apiVersion, kind)
	}
}

// applyPulsarInstance applies PulsarInstance resource
func applyPulsarInstance(ctx context.Context, apiClient *sncloud.APIClient, jsonContent string, organization string, dryRun bool) (string, error) {
	var instance sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstance
	if err := json.Unmarshal([]byte(jsonContent), &instance); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON to PulsarInstance: %v", err)
	}

	// Ensure namespace is set correctly
	if instance.Metadata == nil {
		instance.Metadata = &sncloud.V1ObjectMeta{}
	}
	if instance.Metadata.Namespace == nil || *instance.Metadata.Namespace == "" {
		ns := organization
		instance.Metadata.Namespace = &ns
	}

	name := ""
	if instance.Metadata.Name != nil {
		name = *instance.Metadata.Name
	}

	// Check if resource already exists
	exists := false
	var existingResourceVersion *string

	if name != "" {
		// Try to get existing resource
		existingInstance, bdy, err := apiClient.CloudStreamnativeIoV1alpha1Api.ReadCloudStreamnativeIoV1alpha1NamespacedPulsarInstance(ctx, name, organization).Execute()
		defer bdy.Body.Close()
		if err == nil {
			exists = true
			if existingInstance.Metadata != nil && existingInstance.Metadata.ResourceVersion != nil {
				existingResourceVersion = existingInstance.Metadata.ResourceVersion
			}
		}
	}

	var verb string

	// Convert dryRun bool to string parameter required by API
	dryRunStr := "All"

	// Create or update based on whether resource exists
	var bdy *http.Response
	var err error
	if exists {
		verb = "updated"
		// Make sure resourceVersion is set to support updates
		if existingResourceVersion != nil {
			if instance.Metadata.ResourceVersion == nil {
				instance.Metadata.ResourceVersion = existingResourceVersion
			}
		}

		// Use Replace method to update resource
		request := apiClient.CloudStreamnativeIoV1alpha1Api.ReplaceCloudStreamnativeIoV1alpha1NamespacedPulsarInstance(
			ctx, name, organization).Body(instance)
		if dryRun {
			request = request.DryRun(dryRunStr)
		}
		_, bdy, err = request.Execute()
		defer bdy.Body.Close()
	} else {
		verb = "created"
		// Create new resource
		request := apiClient.CloudStreamnativeIoV1alpha1Api.CreateCloudStreamnativeIoV1alpha1NamespacedPulsarInstance(
			ctx, organization).Body(instance)
		if dryRun {
			request = request.DryRun(dryRunStr)
		}
		_, bdy, err = request.Execute()
		defer bdy.Body.Close()
	}

	if err != nil {
		body, innerErr := io.ReadAll(bdy.Body)
		if innerErr != nil {
			return "", fmt.Errorf("failed to read body: %v", innerErr)
		}
		return "", fmt.Errorf("failed to %s PulsarInstance: %v (%s)", verb, err, string(body))
	}

	if dryRun {
		return fmt.Sprintf("PulsarInstance %q would be %s (dry run)", name, verb), nil
	}
	return fmt.Sprintf("PulsarInstance %q %s", name, verb), nil
}

// applyPulsarCluster applies PulsarCluster resource
func applyPulsarCluster(ctx context.Context, apiClient *sncloud.APIClient, jsonContent string, organization string, dryRun bool) (string, error) {
	var cluster sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarCluster
	if err := json.Unmarshal([]byte(jsonContent), &cluster); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON to PulsarCluster: %v", err)
	}

	// Ensure namespace is set correctly
	if cluster.Metadata == nil {
		cluster.Metadata = &sncloud.V1ObjectMeta{}
	}
	if cluster.Metadata.Namespace == nil || *cluster.Metadata.Namespace == "" {
		ns := organization
		cluster.Metadata.Namespace = &ns
	}

	name := ""
	if cluster.Metadata.Name != nil {
		name = *cluster.Metadata.Name
	}
	// Check if resource already exists
	exists := false
	var existingResourceVersion *string

	if name != "" {
		// Try to get existing resource
		existingCluster, bdy, err := apiClient.CloudStreamnativeIoV1alpha1Api.ReadCloudStreamnativeIoV1alpha1NamespacedPulsarCluster(ctx, name, organization).Execute()
		defer bdy.Body.Close()
		if err == nil {
			exists = true
			if existingCluster.Metadata != nil && existingCluster.Metadata.ResourceVersion != nil {
				existingResourceVersion = existingCluster.Metadata.ResourceVersion
			}
		}
	}

	var verb string

	// Convert dryRun bool to string parameter required by API
	dryRunStr := "All"

	// Create or update based on whether resource exists
	var bdy *http.Response
	var err error
	if exists {
		verb = "updated"
		// Make sure resourceVersion is set to support updates
		if existingResourceVersion != nil {
			if cluster.Metadata.ResourceVersion == nil {
				cluster.Metadata.ResourceVersion = existingResourceVersion
			}
		}

		// Use Replace method to update resource
		request := apiClient.CloudStreamnativeIoV1alpha1Api.ReplaceCloudStreamnativeIoV1alpha1NamespacedPulsarCluster(
			ctx, name, organization).Body(cluster)
		if dryRun {
			request = request.DryRun(dryRunStr)
		}

		_, bdy, err = request.Execute()
		defer bdy.Body.Close()
	} else {
		verb = "created"
		// Create new resource
		request := apiClient.CloudStreamnativeIoV1alpha1Api.CreateCloudStreamnativeIoV1alpha1NamespacedPulsarCluster(
			ctx, organization).Body(cluster)
		if dryRun {
			request = request.DryRun(dryRunStr)
		}
		_, bdy, err = request.Execute()
		defer bdy.Body.Close()
	}

	if err != nil {
		body, innerErr := io.ReadAll(bdy.Body)
		if innerErr != nil {
			return "", fmt.Errorf("failed to read body: %v", innerErr)
		}
		return "", fmt.Errorf("failed to %s PulsarCluster: %v (%s)", verb, err, string(body))
	}

	if dryRun {
		return fmt.Sprintf("PulsarCluster %q would be %s (dry run)", name, verb), nil
	}
	return fmt.Sprintf("PulsarCluster %q %s", name, verb), nil
}

// handleStreamNativeResourcesDelete handles the streaming_cloud_resources_delete tool
func handleStreamNativeResourcesDelete(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	snConfig := getOptions(ctx)
	organization := snConfig.Organization

	name, err := requiredParam[string](request.Params.Arguments, "name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get name: %v", err)), nil
	}

	resourceType, err := requiredParam[string](request.Params.Arguments, "type")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get type: %v", err)), nil
	}

	apiClient, err := config.GetAPIClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get API client: %v", err)), nil
	}

	switch resourceType {
	case "PulsarInstance":
		//nolint:bodyclose
		_, _, err = apiClient.CloudStreamnativeIoV1alpha1Api.DeleteCloudStreamnativeIoV1alpha1NamespacedPulsarInstance(ctx, name, organization).Execute()
	case "PulsarCluster":
		//nolint:bodyclose
		_, _, err = apiClient.CloudStreamnativeIoV1alpha1Api.DeleteCloudStreamnativeIoV1alpha1NamespacedPulsarCluster(ctx, name, organization).Execute()
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unsupported resource type: %s", resourceType)), nil
	}

	// the delete operation will return a V1Status object, which is not handled by the SDK
	if err != nil && !strings.Contains(err.Error(), "json: cannot unmarshal") {
		return mcp.NewToolResultError(fmt.Sprintf("failed to delete resource: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Resource %q %s deleted", name, resourceType)), nil
}
