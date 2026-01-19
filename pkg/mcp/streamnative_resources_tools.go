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
	"strings"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
	context2 "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	sncloud "github.com/streamnative/streamnative-mcp-server/sdk/sdk-apiserver"
)

type streamnativeResourcesApplyInput struct {
	JSONContent string `json:"json_content" jsonschema:"The JSON content to apply."`
	DryRun      bool   `json:"dry_run,omitempty" jsonschema:"If true, only validate the resource without applying it to the server."`
}

type streamnativeResourcesDeleteInput struct {
	Name string `json:"name" jsonschema:"The name of the resource to delete."`
	Type string `json:"type" jsonschema:"The type of the resource to delete, it can be PulsarInstance or PulsarCluster."`
}

// StreamNativeAddResourceTools adds StreamNative resources tools
func StreamNativeAddResourceTools(s *sdk.Server, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeatureStreamNativeCloud)) && !slices.Contains(features, string(FeatureAll)) {
		return
	}

	if !readOnly {
		applySchema, err := InputSchema[streamnativeResourcesApplyInput]()
		if err != nil {
			return
		}
		setSchemaPropertyDefault(applySchema, "dry_run", false)
		applyTool := &sdk.Tool{
			Name:        "sncloud_resources_apply",
			Description: "Apply StreamNative Cloud resources from JSON definitions. This tool allows you to apply (create or update) StreamNative Cloud resources such as PulsarInstances and PulsarClusters using JSON definitions. Please give feedback to USER if the resource is applied with error, and ask USER to check the resource definition.",
			InputSchema: applySchema,
			Annotations: &sdk.ToolAnnotations{Title: "Apply StreamNative Cloud Resources"},
		}

		deleteSchema, err := InputSchema[streamnativeResourcesDeleteInput]()
		if err != nil {
			return
		}
		setSchemaPropertyEnum(deleteSchema, "type", []string{"PulsarInstance", "PulsarCluster"})
		deleteTool := &sdk.Tool{
			Name:        "sncloud_resources_delete",
			Description: "Delete StreamNative Cloud resources. This tool allows you to delete StreamNative Cloud resources such as PulsarInstances and PulsarClusters.",
			InputSchema: deleteSchema,
			Annotations: &sdk.ToolAnnotations{
				Title:           "Delete StreamNative Cloud Resources",
				DestructiveHint: boolPtr(true),
			},
		}

		sdk.AddTool(s, applyTool, handleStreamNativeResourcesApply)
		sdk.AddTool(s, deleteTool, handleStreamNativeResourcesDelete)
	}
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

// handleStreamNativeResourcesApply handles the streaming_cloud_resources_apply tool
func handleStreamNativeResourcesApply(ctx context.Context, _ *sdk.CallToolRequest, input streamnativeResourcesApplyInput) (*sdk.CallToolResult, any, error) {
	// Get necessary parameters
	snConfig := common.GetOptions(ctx)
	organization := snConfig.Organization
	if organization == "" {
		return nil, nil, fmt.Errorf("no organization is set. please set the organization using the appropriate context tool")
	}

	// Get JSON content
	jsonContent := strings.TrimSpace(input.JSONContent)
	if jsonContent == "" {
		return nil, nil, fmt.Errorf("no valid resources found in the provided JSON")
	}

	// Get dry_run flag
	dryRun := input.DryRun

	// Get API client from session
	session := context2.GetSNCloudSession(ctx)
	if session == nil {
		return nil, nil, fmt.Errorf("failed to get StreamNative Cloud session")
	}

	apiClient, err := session.GetAPIClient()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get API client: %v", err)
	}

	// Parse JSON document
	var resource Resource
	err = json.Unmarshal([]byte(jsonContent), &resource)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse JSON document: %v", err)
	}

	// Check if resource is valid
	if resource.APIVersion == "" || resource.Kind == "" {
		return nil, nil, fmt.Errorf("invalid resource definition")
	}

	// Set namespace if not specified
	if resource.Metadata.Namespace == "" {
		resource.Metadata.Namespace = organization
	}

	// Apply resource
	result, err := applyResource(ctx, apiClient, resource, jsonContent, organization, dryRun)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to apply resource: %v", err)
	}

	return newToolResultText(result), nil, nil
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
		defer func() { _ = bdy.Body.Close() }()
		if err == nil {
			exists = true
			if existingInstance.Metadata != nil && existingInstance.Metadata.ResourceVersion != nil {
				existingResourceVersion = existingInstance.Metadata.ResourceVersion
			}
		}
	}

	if dryRun {
		return "Resource validation successful (dry run)", nil
	}

	// Set resource version for update
	if exists {
		if instance.Metadata != nil {
			instance.Metadata.ResourceVersion = existingResourceVersion
		}

		// Update existing instance
		_, resp, err := apiClient.CloudStreamnativeIoV1alpha1Api.ReplaceCloudStreamnativeIoV1alpha1NamespacedPulsarInstance(ctx, name, organization).Body(instance).Execute()
		if resp != nil {
			defer func() { _ = resp.Body.Close() }()
		}
		if err != nil {
			return "", fmt.Errorf("failed to update instance: %v", err)
		}
		return fmt.Sprintf("Updated PulsarInstance %s successfully", name), nil
	}

	// Create new instance
	_, resp, err := apiClient.CloudStreamnativeIoV1alpha1Api.CreateCloudStreamnativeIoV1alpha1NamespacedPulsarInstance(ctx, organization).Body(instance).Execute()
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	if err != nil {
		return "", fmt.Errorf("failed to create instance: %v", err)
	}
	return fmt.Sprintf("Created PulsarInstance %s successfully", name), nil
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

	if dryRun {
		return "Resource validation successful (dry run)", nil
	}

	// Set resource version for update
	if exists {
		if cluster.Metadata != nil {
			cluster.Metadata.ResourceVersion = existingResourceVersion
		}

		// Update existing cluster
		_, resp, err := apiClient.CloudStreamnativeIoV1alpha1Api.ReplaceCloudStreamnativeIoV1alpha1NamespacedPulsarCluster(ctx, name, organization).Body(cluster).Execute()
		if resp != nil {
			defer func() { _ = resp.Body.Close() }()
		}
		if err != nil {
			return "", fmt.Errorf("failed to update cluster: %v", err)
		}
		return fmt.Sprintf("Updated PulsarCluster %s successfully", name), nil
	}

	// Create new cluster
	_, resp, err := apiClient.CloudStreamnativeIoV1alpha1Api.CreateCloudStreamnativeIoV1alpha1NamespacedPulsarCluster(ctx, organization).Body(cluster).Execute()
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	if err != nil {
		return "", fmt.Errorf("failed to create cluster: %v", err)
	}
	return fmt.Sprintf("Created PulsarCluster %s successfully", name), nil
}

// handleStreamNativeResourcesDelete handles the streaming_cloud_resources_delete tool
func handleStreamNativeResourcesDelete(ctx context.Context, _ *sdk.CallToolRequest, input streamnativeResourcesDeleteInput) (*sdk.CallToolResult, any, error) {
	// Get necessary parameters
	snConfig := common.GetOptions(ctx)
	organization := snConfig.Organization
	if organization == "" {
		return nil, nil, fmt.Errorf("no organization is set. please set the organization using the appropriate context tool")
	}

	name := input.Name
	if name == "" {
		return nil, nil, fmt.Errorf("missing required parameter 'name'")
	}

	resourceType := input.Type
	if resourceType == "" {
		return nil, nil, fmt.Errorf("missing required parameter 'type'")
	}

	// Get API client from session
	session := context2.GetSNCloudSession(ctx)
	if session == nil {
		return nil, nil, fmt.Errorf("failed to get StreamNative Cloud session")
	}

	apiClient, err := session.GetAPIClient()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get API client: %v", err)
	}

	// Delete resource based on type
	switch resourceType {
	case "PulsarInstance":
		_, resp, err := apiClient.CloudStreamnativeIoV1alpha1Api.DeleteCloudStreamnativeIoV1alpha1NamespacedPulsarInstance(ctx, name, organization).Execute()
		if resp != nil {
			defer func() { _ = resp.Body.Close() }()
		}
		if err != nil {
			return nil, nil, fmt.Errorf("failed to delete PulsarInstance: %v", err)
		}
	case "PulsarCluster":
		_, resp, err := apiClient.CloudStreamnativeIoV1alpha1Api.DeleteCloudStreamnativeIoV1alpha1NamespacedPulsarCluster(ctx, name, organization).Execute()
		if resp != nil {
			defer func() { _ = resp.Body.Close() }()
		}
		if err != nil {
			return nil, nil, fmt.Errorf("failed to delete PulsarCluster: %v", err)
		}
	default:
		return nil, nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	return newToolResultText(fmt.Sprintf("Resource %q %s deleted", name, resourceType)), nil, nil
}

func boolPtr(value bool) *bool {
	return &value
}

// applyResourceRequest sends a JSON request with the provided body
