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
	"net/http"
	"slices"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
	context2 "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	sncloud "github.com/streamnative/streamnative-mcp-server/sdk/sdk-apiserver"
	"k8s.io/utils/ptr"
)

// ServerlessPoolMember describes a serverless pool option.
type ServerlessPoolMember struct {
	Provider  string
	Namespace string
	Pool      string
	Location  string
}

var (
	// ServerlessPoolMemberList defines the supported serverless pools.
	ServerlessPoolMemberList = []ServerlessPoolMember{
		{
			Provider:  "azure",
			Namespace: "streamnative",
			Pool:      "shared-azure",
			Location:  "eastus",
		},
		{
			Provider:  "aws",
			Namespace: "streamnative",
			Pool:      "shared-aws",
			Location:  "us-east-2",
		},
		// {
		// 	Provider:  "gcloud",
		// 	Namespace: "streamnative",
		// 	Pool:      "shared-gcp",
		// 	Location:  "us-central1",
		// },
	}
	// AvailableProviders lists supported cloud providers.
	AvailableProviders = []string{"azure", "aws", "gcloud"}
)

// RegisterPrompts registers prompt handlers on the server.
func RegisterPrompts(s *server.MCPServer) {
	s.AddPrompt(NewListSNCloudClustersPrompt(), HandleListSNCloudClusters)
	s.AddPrompt(NewReadSNCloudClusterPrompt(), HandleReadSNCloudCluster)
	s.AddPrompt(
		NewBuildSNCloudServerlessClusterPrompt(),
		HandleBuildSNCloudServerlessCluster,
	)
}

// NewListSNCloudClustersPrompt creates the reusable StreamNative Cloud cluster list prompt definition.
func NewListSNCloudClustersPrompt() mcp.Prompt {
	return mcp.NewPrompt("list-sncloud-clusters",
		mcp.WithPromptDescription("List clusters available in the current StreamNative Cloud session."),
	)
}

// NewReadSNCloudClusterPrompt creates the reusable StreamNative Cloud cluster read prompt definition.
func NewReadSNCloudClusterPrompt() mcp.Prompt {
	return mcp.NewPrompt("read-sncloud-cluster",
		mcp.WithPromptDescription("Read details for a StreamNative Cloud cluster."),
		mcp.WithArgument("name", mcp.RequiredArgument(), mcp.ArgumentDescription("The name of the cluster")),
	)
}

// NewBuildSNCloudServerlessClusterPrompt creates the reusable serverless cluster build prompt definition.
func NewBuildSNCloudServerlessClusterPrompt() mcp.Prompt {
	return mcp.NewPrompt("build-sncloud-serverless-cluster",
		mcp.WithPromptDescription("Build a serverless cluster in StreamNative Cloud."),
		mcp.WithArgument("instance-name", mcp.RequiredArgument(), mcp.ArgumentDescription("The name of the Pulsar instance, cannot reuse the name of existing instance.")),
		mcp.WithArgument("cluster-name", mcp.RequiredArgument(), mcp.ArgumentDescription("The name of the Pulsar cluster, cannot reuse the name of existing cluster.")),
		mcp.WithArgument("provider", mcp.ArgumentDescription("The cloud provider, could be `aws`, `gcp`, `azure`. If the selected provider do not serve serverless cluster, the prompt will return an error. If not specified, the system will use a random provider depending on the availability.")),
	)
}

// HandleListSNCloudClusters handles listing StreamNative Cloud clusters.
func HandleListSNCloudClusters(ctx context.Context, _ mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	// Get API client from session
	session := context2.GetSNCloudSession(ctx)
	if session == nil {
		return nil, fmt.Errorf("failed to get StreamNative Cloud session")
	}

	apiClient, err := session.GetAPIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get API client: %v", err)
	}

	clusters, clustersBody, err := apiClient.CloudStreamnativeIoV1alpha1Api.ListCloudStreamnativeIoV1alpha1NamespacedPulsarCluster(ctx, session.Ctx.Organization).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list pulsar clusters: %v", err)
	}
	defer func() { _ = clustersBody.Body.Close() }()

	var messages = make(
		[]mcp.PromptMessage,
		len(clusters.Items)+1,
	)

	messages[0] = mcp.PromptMessage{
		Content: mcp.TextContent{
			Type: "text",
			Text: fmt.Sprintf(
				"There are %d Pulsar clusters in the StreamNative Cloud from organization %s:",
				len(clusters.Items),
				session.Ctx.Organization,
			),
		},
		Role: mcp.RoleUser,
	}

	for i, cluster := range clusters.Items {
		instanceName := cluster.Spec.InstanceName
		displayName := cluster.Spec.DisplayName
		if displayName == nil || *displayName == "" {
			displayName = cluster.Metadata.Name
		}

		status := "Not Ready"
		if common.IsClusterAvailable(cluster) {
			status = "Ready"
		}

		engineType := common.GetEngineType(cluster)

		messages[i+1] = mcp.PromptMessage{
			Content: mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf(
					"Instance Name: %s\nCluster Name: %s\nCluster Display Name: %s\nCluster Status: %s\nCluster Engine Type: %s",
					instanceName,
					*cluster.Metadata.Name,
					*displayName,
					status,
					engineType,
				),
			},
			Role: mcp.RoleUser,
		}
	}

	return &mcp.GetPromptResult{
		Description: fmt.Sprintf("Clusters from StreamNative Cloud organization %s. Use `sncloud_context_use_cluster` to bind the current session to a cluster before running cluster-specific tools.", session.Ctx.Organization),
		Messages:    messages,
	}, nil
}

// HandleReadSNCloudCluster handles reading a StreamNative Cloud cluster.
func HandleReadSNCloudCluster(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	// Get API client from session
	session := context2.GetSNCloudSession(ctx)
	if session == nil {
		return nil, fmt.Errorf("failed to get StreamNative Cloud session")
	}

	apiClient, err := session.GetAPIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get API client: %v", err)
	}

	name, err := common.RequiredParam[string](common.ConvertToMapInterface(request.Params.Arguments), "name")
	if err != nil {
		return nil, fmt.Errorf("failed to get name: %v", err)
	}

	cluster, clusterBody, err := apiClient.CloudStreamnativeIoV1alpha1Api.ReadCloudStreamnativeIoV1alpha1NamespacedPulsarCluster(ctx, name, session.Ctx.Organization).Execute()
	defer func() {
		if clusterBody != nil && clusterBody.Body != nil {
			_ = clusterBody.Body.Close()
		}
	}()

	var (
		clusterJSON []byte
	)

	if err == nil {
		if cluster.Metadata == nil {
			return nil, fmt.Errorf("failed to find cluster: %s", name)
		}
		if len(cluster.Metadata.ManagedFields) > 0 {
			cluster.Metadata.ManagedFields = nil
		}
		clusterJSON, err = json.Marshal(cluster)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal pulsar cluster: %v", err)
		}
	} else {
		if clusterBody == nil || clusterBody.StatusCode != http.StatusNotFound {
			return nil, fmt.Errorf("failed to read pulsar cluster: %v", err)
		}

		kafkaCluster, bdy, err := apiClient.CloudStreamnativeIoV1alpha1Api.ReadCloudStreamnativeIoV1alpha1NamespacedKafkaCluster(ctx, name, session.Ctx.Organization).Execute()
		defer func() {
			if bdy != nil && bdy.Body != nil {
				_ = bdy.Body.Close()
			}
		}()
		if err != nil {
			if bdy != nil && bdy.StatusCode == http.StatusNotFound {
				return nil, fmt.Errorf("failed to find cluster: %s", name)
			}
			return nil, fmt.Errorf("failed to read kafka cluster: %v", err)
		}

		if kafkaCluster.Metadata == nil {
			return nil, fmt.Errorf("failed to find cluster: %s", name)
		}
		if len(kafkaCluster.Metadata.ManagedFields) > 0 {
			kafkaCluster.Metadata.ManagedFields = nil
		}
		clusterJSON, err = json.Marshal(kafkaCluster)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal kafka cluster: %v", err)
		}
	}

	var messages = make(
		[]mcp.PromptMessage,
		1,
	)

	messages[0] = mcp.PromptMessage{
		Content: mcp.TextContent{
			Type: "text",
			Text: string(clusterJSON),
		},
		Role: mcp.RoleUser,
	}

	return &mcp.GetPromptResult{
		Description: fmt.Sprintf("Detailed information for StreamNative Cloud cluster %s. Use `sncloud_context_use_cluster` to bind the current session to this cluster before running cluster-specific tools.", name),
		Messages:    messages,
	}, nil
}

// HandleBuildSNCloudServerlessCluster handles building a serverless StreamNative Cloud cluster definition.
func HandleBuildSNCloudServerlessCluster(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	// Get API client from session
	session := context2.GetSNCloudSession(ctx)
	if session == nil {
		return nil, fmt.Errorf("failed to get StreamNative Cloud session")
	}

	apiClient, err := session.GetAPIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get API client: %v", err)
	}
	arguments := common.ConvertToMapInterface(request.Params.Arguments)

	instanceName, err := common.RequiredParam[string](arguments, "instance-name")
	if err != nil {
		return nil, fmt.Errorf("failed to get instance name: %v", err)
	}

	clusterName, err := common.RequiredParam[string](arguments, "cluster-name")
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster name: %v", err)
	}

	provider, hasProvider := common.OptionalParam[string](arguments, "provider")
	if !hasProvider {
		provider = ""
	}
	if provider != "" {
		if !slices.Contains(AvailableProviders, provider) {
			return nil, fmt.Errorf("invalid provider: %s, available providers: %v", provider, AvailableProviders)
		}
	}

	poolOptions, poolOptionsBody, err := apiClient.CloudStreamnativeIoV1alpha1Api.ListCloudStreamnativeIoV1alpha1NamespacedPoolOption(ctx, session.Ctx.Organization).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list pool options: %v", err)
	}
	defer func() { _ = poolOptionsBody.Body.Close() }()
	if poolOptions == nil {
		return nil, fmt.Errorf("no pool options found")
	}

	var poolRef *sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PoolRef
	var selectedLocation *string

	for _, poolOpt := range poolOptions.Items {
		if pr, ok := poolOpt.Spec.GetPoolRefOk(); ok {
			for _, poolMember := range ServerlessPoolMemberList {
				if provider != "" && poolOpt.Spec.CloudType != provider {
					continue
				}
				if pr.Name == poolMember.Pool && pr.Namespace == poolMember.Namespace {
					for _, location := range poolOpt.Spec.Locations {
						if location.Location == poolMember.Location {
							poolRef = pr
							selectedLocation = &location.Location
							break
						}
					}
				}
			}
		}
	}

	if poolRef == nil || selectedLocation == nil {
		return nil, fmt.Errorf("no available pool")
	}

	inst := sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstance{}
	clus := sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarCluster{}

	inst.ApiVersion = ptr.To("cloud.streamnative.io/v1alpha1")
	inst.Kind = ptr.To("PulsarInstance")
	inst.Metadata = &sncloud.V1ObjectMeta{
		Name:      &instanceName,
		Namespace: &session.Ctx.Organization,
		Labels: &map[string]string{
			"managed-by": "streamnative-mcp",
		},
	}

	inst.Spec = &sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstanceSpec{
		AvailabilityMode: "zonal",
		PoolRef:          poolRef,
		Type:             ptr.To("serverless"),
	}

	clus.ApiVersion = ptr.To("cloud.streamnative.io/v1alpha1")
	clus.Kind = ptr.To("PulsarCluster")
	clus.Metadata = &sncloud.V1ObjectMeta{
		Name:      ptr.To(""),
		Namespace: &session.Ctx.Organization,
		Labels: &map[string]string{
			"managed-by": "streamnative-mcp",
		},
	}

	clus.Spec = &sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarClusterSpec{
		Broker: sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1Broker{
			Replicas: 2,
			Resources: &sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1DefaultNodeResource{
				Cpu:    "1000m",
				Memory: "4294967296",
			},
		},
		DisplayName:    ptr.To(clusterName),
		InstanceName:   instanceName,
		Location:       *selectedLocation,
		ReleaseChannel: ptr.To("rapid"),
	}

	instJSON, err := json.Marshal(inst)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal instance: %v", err)
	}
	clusJSON, err := json.Marshal(clus)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal cluster: %v", err)
	}

	messages := []mcp.PromptMessage{
		{
			Content: mcp.TextContent{
				Type: "text",
				Text: "The following is the Pulsar instance JSON definition and the Pulsar cluster JSON definition, you can use the `sncloud_resources_apply` tool to apply the resources to the StreamNative Cloud. Please directly use the JSON content and not modify the content. The PulsarCluster name is required to be empty. You will need to apply PulsarInstance first, then apply PulsarCluster.",
			},
			Role: mcp.RoleUser,
		},
		{
			Content: mcp.TextContent{
				Type: "text",
				Text: string(instJSON),
			},
			Role: mcp.RoleUser,
		},
		{
			Content: mcp.TextContent{
				Type: "text",
				Text: string(clusJSON),
			},
			Role: mcp.RoleUser,
		},
	}

	return &mcp.GetPromptResult{
		Description: fmt.Sprintf("Create a new Serverless Pulsar cluster %s's related resources that can be applied to the StreamNative Cloud.", clusterName),
		Messages:    messages,
	}, nil
}
