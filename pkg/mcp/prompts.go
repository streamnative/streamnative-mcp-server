// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	sncloud "github.com/streamnative/streamnative-mcp-server/sdk/sdk-apiserver"
)

func RegisterPrompts(s *server.MCPServer) {
	s.AddPrompt(mcp.NewPrompt("list-streamnative-cloud-pulsar-clusters",
		mcp.WithPromptDescription("List all Pulsar clusters in the StreamNative Cloud"),
	), handleListPulsarClusters)
	s.AddPrompt(mcp.NewPrompt("read-streamnative-cloud-pulsar-cluster",
		mcp.WithPromptDescription("Read a Pulsar cluster in the StreamNative Cloud"),
		mcp.WithArgument("name", mcp.RequiredArgument(), mcp.ArgumentDescription("The name of the Pulsar cluster")),
	), handleReadPulsarCluster)
}

func handleListPulsarClusters(ctx context.Context, _ mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	options := getOptions(ctx)
	apiClient, err := config.GetAPIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get API client: %v", err)
	}

	clusters, clustersBody, err := apiClient.CloudStreamnativeIoV1alpha1Api.ListCloudStreamnativeIoV1alpha1NamespacedPulsarCluster(ctx, options.Organization).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list pulsar clusters: %v", err)
	}
	defer clustersBody.Body.Close()

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
				options.Organization,
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
		if isClusterAvailable(cluster) {
			status = "Ready"
		}

		engineType := getEngineType(cluster)

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
		Description: fmt.Sprintf("Pulsar clusters from StreamNative Cloud organization %s, you can use `streamnative_cloud_context_use_cluster` tool to switch to selected cluster, and use pulsar and kafka tools to interact with the cluster.", options.Organization),
		Messages:    messages,
	}, nil
}

func handleReadPulsarCluster(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	options := getOptions(ctx)
	apiClient, err := config.GetAPIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get API client: %v", err)
	}

	name, err := requiredParam[string](convertToMapInterface(request.Params.Arguments), "name")
	if err != nil {
		return nil, fmt.Errorf("failed to get name: %v", err)
	}

	clusters, clustersBody, err := apiClient.CloudStreamnativeIoV1alpha1Api.ListCloudStreamnativeIoV1alpha1NamespacedPulsarCluster(ctx, options.Organization).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list pulsar clusters: %v", err)
	}
	defer clustersBody.Body.Close()
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

	context, err := json.Marshal(cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal cluster: %v", err)
	}

	var messages = make(
		[]mcp.PromptMessage,
		1,
	)

	messages[0] = mcp.PromptMessage{
		Content: mcp.TextContent{
			Type: "text",
			Text: string(context),
		},
		Role: mcp.RoleUser,
	}

	return &mcp.GetPromptResult{
		Description: fmt.Sprintf("Detailed information of Pulsar cluster %s, you can use `streamnative_cloud_context_use_cluster` tool to switch to this cluster, and use pulsar and kafka tools to interact with the cluster.", name),
		Messages:    messages,
	}, nil
}
