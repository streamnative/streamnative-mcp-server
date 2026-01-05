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
	"io"
	"net/http"
	"slices"
	"strings"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/adapter"
	context2 "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	sncloud "github.com/streamnative/streamnative-mcp-server/sdk/sdk-apiserver"
)

// StreamNativeAddResourceTools adds StreamNative resources tools
func StreamNativeAddResourceTools(s *MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeatureStreamNativeCloud)) && !slices.Contains(features, string(FeatureAll)) {
		return
	}

	if !readOnly {
		// Add Apply tool
		applyTool := builders.NewTool("sncloud_resources_apply",
			builders.WithDescription("Apply StreamNative Cloud resources from JSON definitions. This tool allows you to apply (create or update) StreamNative Cloud resources such as PulsarInstances and PulsarClusters using JSON definitions. Please give feedback to USER if the resource is applied with error, and ask USER to check the resource definition."),
			builders.WithString("json_content", builders.Required(),
				builders.Description("The JSON content to apply."),
			),
			builders.WithBoolean("dry_run",
				builders.Description("If true, only validate the resource without applying it to the server."),
				builders.DefaultBool(false),
			),
			builders.WithToolAnnotation(builders.ToolAnnotation{
				Title: "Apply StreamNative Cloud Resources",
			}),
		)
		// Add delete tool
		deleteTool := builders.NewTool("sncloud_resources_delete",
			builders.WithDescription("Delete StreamNative Cloud resources. This tool allows you to delete StreamNative Cloud resources such as PulsarInstances and PulsarClusters."),
			builders.WithString("name", builders.Required(),
				builders.Description("The name of the resource to delete."),
			),
			builders.WithString("type", builders.Required(),
				builders.Description("The type of the resource to delete, it can be PulsarInstance or PulsarCluster."),
				builders.Enum("PulsarInstance", "PulsarCluster"),
			),
			builders.WithToolAnnotation(builders.ToolAnnotation{
				Title:           "Delete StreamNative Cloud Resources",
				DestructiveHint: func() *bool { b := true; return &b }(),
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
func handleStreamNativeResourcesApply(ctx context.Context, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	// Get necessary parameters
	snConfig := common.GetOptions(ctx)
	organization := snConfig.Organization
	if organization == "" {
		return adapter.NewErrorResult("No organization is set. Please set the organization using the appropriate context tool."), nil
	}

	// Get YAML content
	jsonContent, err := adapter.RequireString(request, "json_content")
	if err != nil {
		return adapter.NewErrorResult("Failed to get json_content: %v", err), nil
	}

	// Get dry_run flag
	dryRun := adapter.GetBool(request, "dry_run", false)

	// Get API client from session
	session := context2.GetSNCloudSession(ctx)
	if session == nil {
		return nil, fmt.Errorf("failed to get StreamNative Cloud session")
	}

	apiClient, err := session.GetAPIClient()
	if err != nil {
		return adapter.NewErrorResult("Failed to get API client: %v", err), nil
	}

	jsonContent = strings.TrimSpace(jsonContent)
	if jsonContent == "" {
		return adapter.NewErrorResult("No valid resources found in the provided JSON."), nil
	}

	// Parse YAML document
	var resource Resource
	err = json.Unmarshal([]byte(jsonContent), &resource)
	if err != nil {
		return adapter.NewErrorResult("Failed to parse JSON document: %v", err), nil
	}

	// Check if resource is valid
	if resource.APIVersion == "" || resource.Kind == "" {
		return adapter.NewErrorResult("Invalid resource definition."), nil
	}

	// Set namespace if not specified
	if resource.Metadata.Namespace == "" {
		resource.Metadata.Namespace = organization
	}

	// Apply resource
	result, err := applyResource(ctx, apiClient, resource, jsonContent, organization, dryRun)
	if err != nil {
		return adapter.NewErrorResult("Failed to apply resource: %v", err), nil
	}

	return adapter.NewTextResult(result), nil
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
func handleStreamNativeResourcesDelete(ctx context.Context, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	snConfig := common.GetOptions(ctx)
	organization := snConfig.Organization

	name, err := adapter.RequireString(request, "name")
	if err != nil {
		return adapter.NewErrorResult("Failed to get name: %v", err), nil
	}

	resourceType, err := adapter.RequireString(request, "type")
	if err != nil {
		return adapter.NewErrorResult("Failed to get type: %v", err), nil
	}

	// Get API client from session
	session := context2.GetSNCloudSession(ctx)
	if session == nil {
		return nil, fmt.Errorf("failed to get StreamNative Cloud session")
	}

	apiClient, err := session.GetAPIClient()
	if err != nil {
		return adapter.NewErrorResult("Failed to get API client: %v", err), nil
	}

	switch resourceType {
	case "PulsarInstance":
		//nolint:bodyclose
		_, _, err = apiClient.CloudStreamnativeIoV1alpha1Api.DeleteCloudStreamnativeIoV1alpha1NamespacedPulsarInstance(ctx, name, organization).Execute()
	case "PulsarCluster":
		//nolint:bodyclose
		_, _, err = apiClient.CloudStreamnativeIoV1alpha1Api.DeleteCloudStreamnativeIoV1alpha1NamespacedPulsarCluster(ctx, name, organization).Execute()
	default:
		return adapter.NewErrorResult("Unsupported resource type: %s", resourceType), nil
	}

	// the delete operation will return a V1Status object, which is not handled by the SDK
	if err != nil && !strings.Contains(err.Error(), "json: cannot unmarshal") {
		return adapter.NewErrorResult("failed to delete resource: %v", err), nil
	}

	return adapter.NewTextResult(fmt.Sprintf("Resource %q %s deleted", name, resourceType)), nil
}
