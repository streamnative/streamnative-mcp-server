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

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
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
func RegisterPrompts(s *sdk.Server) {
	s.AddPrompt(&sdk.Prompt{
		Name:        "list-sncloud-clusters",
		Description: "List all clusters from the StreamNative Cloud",
	}, HandleListPulsarClusters)
	s.AddPrompt(&sdk.Prompt{
		Name:        "read-sncloud-cluster",
		Description: "Read a cluster from the StreamNative Cloud",
		Arguments: []*sdk.PromptArgument{
			{
				Name:        "name",
				Description: "The name of the cluster",
				Required:    true,
			},
		},
	}, handleReadPulsarCluster)
	s.AddPrompt(&sdk.Prompt{
		Name:        "build-sncloud-serverless-cluster",
		Description: "Build a Serverless cluster in the StreamNative Cloud",
		Arguments: []*sdk.PromptArgument{
			{
				Name:        "instance-name",
				Description: "The name of the Pulsar instance, cannot reuse the name of existing instance.",
				Required:    true,
			},
			{
				Name:        "cluster-name",
				Description: "The name of the Pulsar cluster, cannot reuse the name of existing cluster.",
				Required:    true,
			},
			{
				Name:        "provider",
				Description: "The cloud provider, could be `aws`, `gcp`, `azure`. If the selected provider do not serve serverless cluster, the prompt will return an error. If not specified, the system will use a random provider depending on the availability.",
			},
		},
	}, handleBuildServerlessPulsarCluster)
}

// HandleListPulsarClusters handles listing StreamNative Cloud clusters.
func HandleListPulsarClusters(ctx context.Context, _ *sdk.GetPromptRequest) (*sdk.GetPromptResult, error) {
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

	messages := make([]*sdk.PromptMessage, 0, len(clusters.Items)+1)

	messages = append(messages, &sdk.PromptMessage{
		Content: &sdk.TextContent{
			Text: fmt.Sprintf(
				"There are %d Pulsar clusters in the StreamNative Cloud from organization %s:",
				len(clusters.Items),
				session.Ctx.Organization,
			),
		},
		Role: "user",
	})

	for _, cluster := range clusters.Items {
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

		messages = append(messages, &sdk.PromptMessage{
			Content: &sdk.TextContent{
				Text: fmt.Sprintf(
					"Instance Name: %s\nCluster Name: %s\nCluster Display Name: %s\nCluster Status: %s\nCluster Engine Type: %s",
					instanceName,
					*cluster.Metadata.Name,
					*displayName,
					status,
					engineType,
				),
			},
			Role: "user",
		})
	}

	return &sdk.GetPromptResult{
		Description: fmt.Sprintf("Pulsar clusters from StreamNative Cloud organization %s, you can use `sncloud_context_use_cluster` tool to switch to selected cluster, and use pulsar and kafka tools to interact with the cluster.", session.Ctx.Organization),
		Messages:    messages,
	}, nil
}

func handleReadPulsarCluster(ctx context.Context, request *sdk.GetPromptRequest) (*sdk.GetPromptResult, error) {
	// Get API client from session
	session := context2.GetSNCloudSession(ctx)
	if session == nil {
		return nil, fmt.Errorf("failed to get StreamNative Cloud session")
	}

	apiClient, err := session.GetAPIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get API client: %v", err)
	}

	name := ""
	if request != nil && request.Params != nil {
		name = request.Params.Arguments["name"]
	}
	if name == "" {
		return nil, fmt.Errorf("failed to get name: missing required parameter 'name'")
	}

	clusters, clustersBody, err := apiClient.CloudStreamnativeIoV1alpha1Api.ListCloudStreamnativeIoV1alpha1NamespacedPulsarCluster(ctx, session.Ctx.Organization).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list pulsar clusters: %v", err)
	}
	defer func() { _ = clustersBody.Body.Close() }()
	var cluster sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarCluster
	for _, c := range clusters.Items {
		if *c.Metadata.Name == name {
			cluster = c
			break
		}
	}

	if cluster.Metadata == nil {
		return nil, fmt.Errorf("failed to find pulsar cluster: %s", name)
	}

	if cluster.Metadata != nil && len(cluster.Metadata.ManagedFields) > 0 {
		cluster.Metadata.ManagedFields = nil
	}

	contextJSON, err := json.Marshal(cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal cluster: %v", err)
	}

	messages := []*sdk.PromptMessage{
		{
			Content: &sdk.TextContent{Text: string(contextJSON)},
			Role:    "user",
		},
	}

	return &sdk.GetPromptResult{
		Description: fmt.Sprintf("Detailed information of Pulsar cluster %s, you can use `sncloud_context_use_cluster` tool to switch to this cluster, and use pulsar and kafka tools to interact with the cluster.", name),
		Messages:    messages,
	}, nil
}

func handleBuildServerlessPulsarCluster(ctx context.Context, request *sdk.GetPromptRequest) (*sdk.GetPromptResult, error) {
	// Get API client from session
	session := context2.GetSNCloudSession(ctx)
	if session == nil {
		return nil, fmt.Errorf("failed to get StreamNative Cloud session")
	}

	apiClient, err := session.GetAPIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get API client: %v", err)
	}

	arguments := map[string]string{}
	if request != nil && request.Params != nil && request.Params.Arguments != nil {
		arguments = request.Params.Arguments
	}

	instanceName := arguments["instance-name"]
	if instanceName == "" {
		return nil, fmt.Errorf("failed to get instance name: missing required parameter 'instance-name'")
	}

	clusterName := arguments["cluster-name"]
	if clusterName == "" {
		return nil, fmt.Errorf("failed to get cluster name: missing required parameter 'cluster-name'")
	}

	provider := arguments["provider"]
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

	messages := []*sdk.PromptMessage{
		{
			Content: &sdk.TextContent{
				Text: "The following is the Pulsar instance JSON definition and the Pulsar cluster JSON definition, you can use the `sncloud_resources_apply` tool to apply the resources to the StreamNative Cloud. Please directly use the JSON content and not modify the content. The PulsarCluster name is required to be empty. You will need to apply PulsarInstance first, then apply PulsarCluster.",
			},
			Role: "user",
		},
		{
			Content: &sdk.TextContent{Text: string(instJSON)},
			Role:    "user",
		},
		{
			Content: &sdk.TextContent{Text: string(clusJSON)},
			Role:    "user",
		},
	}

	return &sdk.GetPromptResult{
		Description: fmt.Sprintf("Create a new Serverless Pulsar cluster %s's related resources that can be applied to the StreamNative Cloud.", clusterName),
		Messages:    messages,
	}, nil
}
