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
			mcp.Description("The type of the resource to delete, it can be Instance, PulsarInstance, PulsarCluster, or KafkaCluster."),
			mcp.Enum("Instance", "PulsarInstance", "PulsarCluster", "KafkaCluster"),
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
	case apiVersion == "cloud.streamnative.io/v1alpha1" && kind == "Instance":
		return applyInstance(ctx, apiClient, jsonContent, organization, dryRun)
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

type sncloudResourceApplyAdapter[T any] struct {
	kind string

	ensureMetadata     func(*T)
	getName            func(*T) *string
	getNamespace       func(*T) *string
	setNamespace       func(*T, string)
	getResourceVersion func(*T) *string
	setResourceVersion func(*T, *string)

	readExistingResourceVersion func(context.Context, *sncloud.APIClient, string, string) (*string, *http.Response, error)
	create                      func(context.Context, *sncloud.APIClient, string, T, bool) (*http.Response, error)
	replace                     func(context.Context, *sncloud.APIClient, string, string, T, bool) (*http.Response, error)
}

func applyTypedSNCloudResource[T any](
	ctx context.Context,
	apiClient *sncloud.APIClient,
	jsonContent string,
	organization string,
	dryRun bool,
	adapter sncloudResourceApplyAdapter[T],
) (string, error) {
	var resource T
	if err := json.Unmarshal([]byte(jsonContent), &resource); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON to %s: %v", adapter.kind, err)
	}

	adapter.ensureMetadata(&resource)
	if namespace := adapter.getNamespace(&resource); namespace == nil || *namespace == "" {
		adapter.setNamespace(&resource, organization)
	}

	name := ""
	if resourceName := adapter.getName(&resource); resourceName != nil {
		name = *resourceName
	}

	exists := false
	var existingResourceVersion *string
	if name != "" {
		readResourceVersion, bdy, err := adapter.readExistingResourceVersion(ctx, apiClient, name, organization)
		if bdy != nil && bdy.Body != nil {
			defer func() {
				_ = bdy.Body.Close()
			}()
		}
		if err == nil {
			exists = true
			existingResourceVersion = readResourceVersion
		}
	}

	var (
		bdy  *http.Response
		err  error
		verb string
	)

	if exists {
		verb = "updated"
		if existingResourceVersion != nil && adapter.getResourceVersion(&resource) == nil {
			adapter.setResourceVersion(&resource, existingResourceVersion)
		}
		bdy, err = adapter.replace(ctx, apiClient, name, organization, resource, dryRun)
	} else {
		verb = "created"
		bdy, err = adapter.create(ctx, apiClient, organization, resource, dryRun)
	}
	if bdy != nil && bdy.Body != nil {
		defer func() {
			_ = bdy.Body.Close()
		}()
	}

	if err != nil {
		return "", formatSNCloudApplyError(adapter.kind, verb, bdy, err)
	}

	if dryRun {
		return fmt.Sprintf("%s %q would be %s (dry run)", adapter.kind, name, verb), nil
	}
	return fmt.Sprintf("%s %q %s", adapter.kind, name, verb), nil
}

func applyV1ObjectMetaSNCloudResource[T any](
	ctx context.Context,
	apiClient *sncloud.APIClient,
	jsonContent string,
	organization string,
	dryRun bool,
	kind string,
	ensureMetadata func(*T) *sncloud.V1ObjectMeta,
	readExistingResourceVersion func(context.Context, *sncloud.APIClient, string, string) (*string, *http.Response, error),
	create func(context.Context, *sncloud.APIClient, string, T, bool) (*http.Response, error),
	replace func(context.Context, *sncloud.APIClient, string, string, T, bool) (*http.Response, error),
) (string, error) {
	return applyTypedSNCloudResource(ctx, apiClient, jsonContent, organization, dryRun, sncloudResourceApplyAdapter[T]{
		kind: kind,
		ensureMetadata: func(resource *T) {
			_ = ensureMetadata(resource)
		},
		getName: func(resource *T) *string {
			return ensureMetadata(resource).Name
		},
		getNamespace: func(resource *T) *string {
			return ensureMetadata(resource).Namespace
		},
		setNamespace: func(resource *T, namespace string) {
			ensureMetadata(resource).Namespace = &namespace
		},
		getResourceVersion: func(resource *T) *string {
			return ensureMetadata(resource).ResourceVersion
		},
		setResourceVersion: func(resource *T, resourceVersion *string) {
			ensureMetadata(resource).ResourceVersion = resourceVersion
		},
		readExistingResourceVersion: readExistingResourceVersion,
		create:                      create,
		replace:                     replace,
	})
}

func formatSNCloudApplyError(kind string, verb string, bdy *http.Response, err error) error {
	if bdy == nil || bdy.Body == nil {
		return fmt.Errorf("failed to %s %s: %v", verb, kind, err)
	}

	body, innerErr := io.ReadAll(bdy.Body)
	if innerErr != nil {
		return fmt.Errorf("failed to read body: %v", innerErr)
	}
	return fmt.Errorf("failed to %s %s: %v (%s)", verb, kind, err, string(body))
}

func applyInstance(ctx context.Context, apiClient *sncloud.APIClient, jsonContent string, organization string, dryRun bool) (string, error) {
	return applyTypedSNCloudResource(ctx, apiClient, jsonContent, organization, dryRun, sncloudResourceApplyAdapter[sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Instance]{
		kind: "Instance",
		ensureMetadata: func(resource *sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Instance) {
			if resource.Metadata == nil {
				resource.Metadata = &sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1InstanceMetadata{}
			}
		},
		getName: func(resource *sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Instance) *string {
			if resource.Metadata == nil {
				return nil
			}
			return resource.Metadata.Name
		},
		getNamespace: func(resource *sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Instance) *string {
			if resource.Metadata == nil {
				return nil
			}
			return resource.Metadata.Namespace
		},
		setNamespace: func(resource *sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Instance, namespace string) {
			resource.Metadata.Namespace = &namespace
		},
		getResourceVersion: func(resource *sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Instance) *string {
			if resource.Metadata == nil {
				return nil
			}
			return resource.Metadata.ResourceVersion
		},
		setResourceVersion: func(resource *sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Instance, resourceVersion *string) {
			resource.Metadata.ResourceVersion = resourceVersion
		},
		readExistingResourceVersion: func(ctx context.Context, apiClient *sncloud.APIClient, name string, organization string) (*string, *http.Response, error) {
			existingInstance, bdy, err := apiClient.CloudStreamnativeIoV1alpha1Api.ReadCloudStreamnativeIoV1alpha1NamespacedInstance(ctx, name, organization).Execute()
			if err != nil {
				return nil, bdy, err
			}
			if existingInstance.Metadata == nil {
				return nil, bdy, nil
			}
			return existingInstance.Metadata.ResourceVersion, bdy, nil
		},
		create: func(ctx context.Context, apiClient *sncloud.APIClient, organization string, resource sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Instance, dryRun bool) (*http.Response, error) {
			request := apiClient.CloudStreamnativeIoV1alpha1Api.CreateCloudStreamnativeIoV1alpha1NamespacedInstance(ctx, organization).Body(resource)
			if dryRun {
				request = request.DryRun("All")
			}
			_, bdy, err := request.Execute()
			return bdy, err
		},
		replace: func(ctx context.Context, apiClient *sncloud.APIClient, name string, organization string, resource sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Instance, dryRun bool) (*http.Response, error) {
			request := apiClient.CloudStreamnativeIoV1alpha1Api.ReplaceCloudStreamnativeIoV1alpha1NamespacedInstance(ctx, name, organization).Body(resource)
			if dryRun {
				request = request.DryRun("All")
			}
			_, bdy, err := request.Execute()
			return bdy, err
		},
	})
}

// applyPulsarInstance applies PulsarInstance resource
func applyPulsarInstance(ctx context.Context, apiClient *sncloud.APIClient, jsonContent string, organization string, dryRun bool) (string, error) {
	return applyV1ObjectMetaSNCloudResource(
		ctx,
		apiClient,
		jsonContent,
		organization,
		dryRun,
		"PulsarInstance",
		func(resource *sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstance) *sncloud.V1ObjectMeta {
			if resource.Metadata == nil {
				resource.Metadata = &sncloud.V1ObjectMeta{}
			}
			return resource.Metadata
		},
		func(ctx context.Context, apiClient *sncloud.APIClient, name string, organization string) (*string, *http.Response, error) {
			existingInstance, bdy, err := apiClient.CloudStreamnativeIoV1alpha1Api.ReadCloudStreamnativeIoV1alpha1NamespacedPulsarInstance(ctx, name, organization).Execute()
			if err != nil {
				return nil, bdy, err
			}
			if existingInstance.Metadata == nil {
				return nil, bdy, nil
			}
			return existingInstance.Metadata.ResourceVersion, bdy, nil
		},
		func(ctx context.Context, apiClient *sncloud.APIClient, organization string, resource sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstance, dryRun bool) (*http.Response, error) {
			request := apiClient.CloudStreamnativeIoV1alpha1Api.CreateCloudStreamnativeIoV1alpha1NamespacedPulsarInstance(ctx, organization).Body(resource)
			if dryRun {
				request = request.DryRun("All")
			}
			_, bdy, err := request.Execute()
			return bdy, err
		},
		func(ctx context.Context, apiClient *sncloud.APIClient, name string, organization string, resource sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstance, dryRun bool) (*http.Response, error) {
			request := apiClient.CloudStreamnativeIoV1alpha1Api.ReplaceCloudStreamnativeIoV1alpha1NamespacedPulsarInstance(ctx, name, organization).Body(resource)
			if dryRun {
				request = request.DryRun("All")
			}
			_, bdy, err := request.Execute()
			return bdy, err
		},
	)
}

// applyPulsarCluster applies PulsarCluster resource
func applyPulsarCluster(ctx context.Context, apiClient *sncloud.APIClient, jsonContent string, organization string, dryRun bool) (string, error) {
	return applyV1ObjectMetaSNCloudResource(
		ctx,
		apiClient,
		jsonContent,
		organization,
		dryRun,
		"PulsarCluster",
		func(resource *sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarCluster) *sncloud.V1ObjectMeta {
			if resource.Metadata == nil {
				resource.Metadata = &sncloud.V1ObjectMeta{}
			}
			return resource.Metadata
		},
		func(ctx context.Context, apiClient *sncloud.APIClient, name string, organization string) (*string, *http.Response, error) {
			existingCluster, bdy, err := apiClient.CloudStreamnativeIoV1alpha1Api.ReadCloudStreamnativeIoV1alpha1NamespacedPulsarCluster(ctx, name, organization).Execute()
			if err != nil {
				return nil, bdy, err
			}
			if existingCluster.Metadata == nil {
				return nil, bdy, nil
			}
			return existingCluster.Metadata.ResourceVersion, bdy, nil
		},
		func(ctx context.Context, apiClient *sncloud.APIClient, organization string, resource sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarCluster, dryRun bool) (*http.Response, error) {
			request := apiClient.CloudStreamnativeIoV1alpha1Api.CreateCloudStreamnativeIoV1alpha1NamespacedPulsarCluster(ctx, organization).Body(resource)
			if dryRun {
				request = request.DryRun("All")
			}
			_, bdy, err := request.Execute()
			return bdy, err
		},
		func(ctx context.Context, apiClient *sncloud.APIClient, name string, organization string, resource sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarCluster, dryRun bool) (*http.Response, error) {
			request := apiClient.CloudStreamnativeIoV1alpha1Api.ReplaceCloudStreamnativeIoV1alpha1NamespacedPulsarCluster(ctx, name, organization).Body(resource)
			if dryRun {
				request = request.DryRun("All")
			}
			_, bdy, err := request.Execute()
			return bdy, err
		},
	)
}

func applyKafkaCluster(ctx context.Context, apiClient *sncloud.APIClient, jsonContent string, organization string, dryRun bool) (string, error) {
	return applyV1ObjectMetaSNCloudResource(
		ctx,
		apiClient,
		jsonContent,
		organization,
		dryRun,
		"KafkaCluster",
		func(resource *sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaCluster) *sncloud.V1ObjectMeta {
			if resource.Metadata == nil {
				resource.Metadata = &sncloud.V1ObjectMeta{}
			}
			return resource.Metadata
		},
		func(ctx context.Context, apiClient *sncloud.APIClient, name string, organization string) (*string, *http.Response, error) {
			existingCluster, bdy, err := apiClient.CloudStreamnativeIoV1alpha1Api.ReadCloudStreamnativeIoV1alpha1NamespacedKafkaCluster(ctx, name, organization).Execute()
			if err != nil {
				return nil, bdy, err
			}
			if existingCluster.Metadata == nil {
				return nil, bdy, nil
			}
			return existingCluster.Metadata.ResourceVersion, bdy, nil
		},
		func(ctx context.Context, apiClient *sncloud.APIClient, organization string, resource sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaCluster, dryRun bool) (*http.Response, error) {
			request := apiClient.CloudStreamnativeIoV1alpha1Api.CreateCloudStreamnativeIoV1alpha1NamespacedKafkaCluster(ctx, organization).Body(resource)
			if dryRun {
				request = request.DryRun("All")
			}
			_, bdy, err := request.Execute()
			return bdy, err
		},
		func(ctx context.Context, apiClient *sncloud.APIClient, name string, organization string, resource sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1KafkaCluster, dryRun bool) (*http.Response, error) {
			request := apiClient.CloudStreamnativeIoV1alpha1Api.ReplaceCloudStreamnativeIoV1alpha1NamespacedKafkaCluster(ctx, name, organization).Body(resource)
			if dryRun {
				request = request.DryRun("All")
			}
			_, bdy, err := request.Execute()
			return bdy, err
		},
	)
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
	case "Instance":
		//nolint:bodyclose
		_, _, err = apiClient.CloudStreamnativeIoV1alpha1Api.DeleteCloudStreamnativeIoV1alpha1NamespacedInstance(ctx, name, organization).Execute()
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
