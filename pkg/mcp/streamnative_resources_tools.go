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

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	context2 "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	sncloud "github.com/streamnative/streamnative-mcp-server/sdk/sdk-apiserver"
)

// StreamNativeAddResourceTools adds StreamNative resources tools
func StreamNativeAddResourceTools(s *server.MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeatureStreamNativeCloud)) && !slices.Contains(features, string(FeatureAll)) {
		return
	}

	if !readOnly {
		s.AddTool(NewSNCloudResourcesApplyTool(), HandleSNCloudResourcesApply)
		s.AddTool(NewSNCloudResourcesDeleteTool(), HandleSNCloudResourcesDelete)
	}
}

// NewSNCloudResourcesApplyTool creates the reusable StreamNative Cloud apply tool definition.
func NewSNCloudResourcesApplyTool() mcp.Tool {
	return mcp.NewTool("sncloud_resources_apply",
		mcp.WithDescription("Apply StreamNative Cloud resources from JSON definitions in the organization bound to the current session."),
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
}

// NewSNCloudResourcesDeleteTool creates the reusable StreamNative Cloud delete tool definition.
func NewSNCloudResourcesDeleteTool() mcp.Tool {
	return mcp.NewTool("sncloud_resources_delete",
		mcp.WithDescription("Delete StreamNative Cloud resources in the organization bound to the current session."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("The name of the resource to delete."),
		),
		mcp.WithString("type", mcp.Required(),
			mcp.Description("The type of the resource to delete, it can be PulsarInstance, PulsarCluster, or KafkaCluster."),
			mcp.Enum("PulsarInstance", "PulsarCluster", "KafkaCluster"),
		),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Delete StreamNative Cloud Resources",
			DestructiveHint: &[]bool{true}[0],
		}),
	)
}

// Resource represents a StreamNative resource manifest.
type Resource struct {
	APIVersion string                 `json:"apiVersion" yaml:"apiVersion"`
	Kind       string                 `json:"kind" yaml:"kind"`
	Metadata   Metadata               `json:"metadata" yaml:"metadata"`
	Spec       map[string]interface{} `json:"spec" yaml:"spec"`
}

// Metadata holds standard resource metadata.
type Metadata struct {
	Name      string            `json:"name" yaml:"name"`
	Namespace string            `json:"namespace" yaml:"namespace"`
	Labels    map[string]string `json:"labels" yaml:"labels"`
}

// HandleSNCloudResourcesApply handles the StreamNative Cloud resource apply tool.
func HandleSNCloudResourcesApply(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	organization := resolveSNCloudOrganization(ctx)
	if organization == "" {
		return mcp.NewToolResultError("No organization is set. Please set the organization using the appropriate context tool."), nil
	}

	// Get YAML content
	jsonContent, err := request.RequireString("json_content")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get json_content: %v", err)), nil
	}

	// Get dry_run flag
	dryRun := request.GetBool("dry_run", false)

	// Get API client from session
	session := context2.GetSNCloudSession(ctx)
	if session == nil {
		return nil, fmt.Errorf("failed to get StreamNative Cloud session")
	}

	apiClient, err := session.GetAPIClient()
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
	case apiVersion == "cloud.streamnative.io/v1alpha1" && kind == "KafkaCluster":
		return applyKafkaCluster(ctx, apiClient, jsonContent, organization, dryRun)
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
		defer func() { _ = bdy.Body.Close() }()
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
		defer func() { _ = bdy.Body.Close() }()
	} else {
		verb = "created"
		// Create new resource
		request := apiClient.CloudStreamnativeIoV1alpha1Api.CreateCloudStreamnativeIoV1alpha1NamespacedPulsarInstance(
			ctx, organization).Body(instance)
		if dryRun {
			request = request.DryRun(dryRunStr)
		}
		_, bdy, err = request.Execute()
		defer func() { _ = bdy.Body.Close() }()
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
		defer func() { _ = bdy.Body.Close() }()
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
		defer func() { _ = bdy.Body.Close() }()
	} else {
		verb = "created"
		// Create new resource
		request := apiClient.CloudStreamnativeIoV1alpha1Api.CreateCloudStreamnativeIoV1alpha1NamespacedPulsarCluster(
			ctx, organization).Body(cluster)
		if dryRun {
			request = request.DryRun(dryRunStr)
		}
		_, bdy, err = request.Execute()
		defer func() { _ = bdy.Body.Close() }()
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

func applyKafkaCluster(ctx context.Context, apiClient *sncloud.APIClient, jsonContent string, organization string, dryRun bool) (string, error) {
	var cluster sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaCluster
	if err := json.Unmarshal([]byte(jsonContent), &cluster); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON to KafkaCluster: %v", err)
	}

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

	exists := false
	var existingResourceVersion *string
	if name != "" {
		existingCluster, bdy, err := apiClient.CloudStreamnativeIoV1alpha1Api.ReadCloudStreamnativeIoV1alpha1NamespacedKafkaCluster(ctx, name, organization).Execute()
		defer func() {
			if bdy != nil && bdy.Body != nil {
				_ = bdy.Body.Close()
			}
		}()
		if err == nil {
			exists = true
			if existingCluster.Metadata != nil && existingCluster.Metadata.ResourceVersion != nil {
				existingResourceVersion = existingCluster.Metadata.ResourceVersion
			}
		}
	}

	var verb string
	if exists {
		verb = "updated"
		if existingResourceVersion != nil && cluster.Metadata.ResourceVersion == nil {
			cluster.Metadata.ResourceVersion = existingResourceVersion
		}

		request := apiClient.CloudStreamnativeIoV1alpha1Api.ReplaceCloudStreamnativeIoV1alpha1NamespacedKafkaCluster(
			ctx, name, organization).Body(cluster)
		if dryRun {
			request = request.DryRun("All")
		}
		_, bdy, err := request.Execute()
		defer func() {
			if bdy != nil && bdy.Body != nil {
				_ = bdy.Body.Close()
			}
		}()
		if err != nil {
			if bdy == nil || bdy.Body == nil {
				return "", fmt.Errorf("failed to %s KafkaCluster: %v", verb, err)
			}
			body, innerErr := io.ReadAll(bdy.Body)
			if innerErr != nil {
				return "", fmt.Errorf("failed to read body: %v", innerErr)
			}
			return "", fmt.Errorf("failed to %s KafkaCluster: %v (%s)", verb, err, string(body))
		}
	} else {
		verb = "created"

		request := apiClient.CloudStreamnativeIoV1alpha1Api.CreateCloudStreamnativeIoV1alpha1NamespacedKafkaCluster(
			ctx, organization).Body(cluster)
		if dryRun {
			request = request.DryRun("All")
		}
		_, bdy, err := request.Execute()
		defer func() {
			if bdy != nil && bdy.Body != nil {
				_ = bdy.Body.Close()
			}
		}()
		if err != nil {
			if bdy == nil || bdy.Body == nil {
				return "", fmt.Errorf("failed to %s KafkaCluster: %v", verb, err)
			}
			body, innerErr := io.ReadAll(bdy.Body)
			if innerErr != nil {
				return "", fmt.Errorf("failed to read body: %v", innerErr)
			}
			return "", fmt.Errorf("failed to %s KafkaCluster: %v (%s)", verb, err, string(body))
		}
	}

	if dryRun {
		return fmt.Sprintf("KafkaCluster %q would be %s (dry run)", name, verb), nil
	}
	return fmt.Sprintf("KafkaCluster %q %s", name, verb), nil
}

// HandleSNCloudResourcesDelete handles the StreamNative Cloud resource delete tool.
func HandleSNCloudResourcesDelete(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	organization := resolveSNCloudOrganization(ctx)
	if organization == "" {
		return mcp.NewToolResultError("No organization is set. Please set the organization using the appropriate context tool."), nil
	}

	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get name: %v", err)), nil
	}

	resourceType, err := request.RequireString("type")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get type: %v", err)), nil
	}

	// Get API client from session
	session := context2.GetSNCloudSession(ctx)
	if session == nil {
		return nil, fmt.Errorf("failed to get StreamNative Cloud session")
	}

	apiClient, err := session.GetAPIClient()
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
	case "KafkaCluster":
		//nolint:bodyclose
		_, _, err = apiClient.CloudStreamnativeIoV1alpha1Api.DeleteCloudStreamnativeIoV1alpha1NamespacedKafkaCluster(ctx, name, organization).Execute()
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unsupported resource type: %s", resourceType)), nil
	}

	// the delete operation will return a V1Status object, which is not handled by the SDK
	if err != nil && !strings.Contains(err.Error(), "json: cannot unmarshal") {
		return mcp.NewToolResultError(fmt.Sprintf("failed to delete resource: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Resource %q %s deleted", name, resourceType)), nil
}

func resolveSNCloudOrganization(ctx context.Context) string {
	if session := context2.GetSNCloudSession(ctx); session != nil {
		if organization := strings.TrimSpace(session.Ctx.Organization); organization != "" {
			return organization
		}
	}

	if organization := strings.TrimSpace(context2.GetSNCloudOrganization(ctx)); organization != "" {
		return organization
	}

	options, ok := ctx.Value(common.OptionsKey).(*config.Options)
	if ok && options != nil {
		return strings.TrimSpace(options.Organization)
	}

	return ""
}
