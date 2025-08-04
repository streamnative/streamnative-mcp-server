// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package pulsar

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

// PulsarAdminClusterToolBuilder implements the ToolBuilder interface for Pulsar Admin Cluster tools
// It provides functionality to build Pulsar cluster management tools
// /nolint:revive
type PulsarAdminClusterToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminClusterToolBuilder creates a new Pulsar Admin Cluster tool builder instance
func NewPulsarAdminClusterToolBuilder() *PulsarAdminClusterToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_cluster",
		Version:     "1.0.0",
		Description: "Pulsar Admin cluster management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "cluster", "admin"},
	}

	features := []string{
		"pulsar-admin-clusters",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminClusterToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar Admin Cluster tool list
// This is the core method implementing the ToolBuilder interface
func (b *PulsarAdminClusterToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool := b.buildClusterTool()
	handler := b.buildClusterHandler(config.ReadOnly)

	return []server.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildClusterTool builds the Pulsar Admin Cluster MCP tool definition
// Migrated from the original tool definition logic
func (b *PulsarAdminClusterToolBuilder) buildClusterTool() mcp.Tool {
	toolDesc := "Unified tool for managing Apache Pulsar clusters.\n" +
		"This tool provides access to various cluster resources and operations, including:\n" +
		"1. Manage clusters (resource=cluster): List, get, create, update, delete clusters\n" +
		"2. Manage peer clusters (resource=peer_clusters): Get, update peer clusters\n" +
		"3. Manage failure domains (resource=failure_domain): List, get, create, update, delete failure domains\n\n" +
		"Different functions are accessed by combining resource and operation parameters, with other parameters used selectively based on operation type.\n\n" +
		"Examples:\n" +
		"- {\"resource\": \"cluster\", \"operation\": \"list\"} lists all clusters\n" +
		"- {\"resource\": \"cluster\", \"operation\": \"get\", \"cluster_name\": \"my-cluster\"} gets cluster configuration\n" +
		"- {\"resource\": \"failure_domain\", \"operation\": \"list\", \"cluster_name\": \"my-cluster\"} lists failure domains\n" +
		"This tool requires Pulsar super-user permissions."

	resourceDesc := "Type of cluster resource to access, available options:\n" +
		"- cluster: Pulsar cluster configuration\n" +
		"- peer_clusters: Peer clusters for geo-replication\n" +
		"- failure_domain: Failure domains for fault tolerance"

	operationDesc := "Operation to perform, available options (depend on resource):\n" +
		"- list: List resources (used with cluster, failure_domain)\n" +
		"- get: Retrieve resource information (used with cluster, peer_clusters, failure_domain)\n" +
		"- create: Create a new resource (used with cluster, failure_domain)\n" +
		"- update: Update an existing resource (used with cluster, peer_clusters, failure_domain)\n" +
		"- delete: Delete a resource (used with cluster, failure_domain)"

	return mcp.NewTool("pulsar_admin_cluster",
		mcp.WithDescription(toolDesc),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description(resourceDesc),
		),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
		),
		mcp.WithString("cluster_name",
			mcp.Description("Name of the Pulsar cluster, required for all operations except 'list' with resource=cluster"),
		),
		mcp.WithString("domain_name",
			mcp.Description("Name of the failure domain, required when resource=failure_domain and operation is get, create, update, or delete"),
		),
		mcp.WithString("service_url",
			mcp.Description("Pulsar cluster web service URL (e.g., http://example.pulsar.io:8080), used when resource=cluster and operation is create or update"),
		),
		mcp.WithString("service_url_tls",
			mcp.Description("Pulsar cluster TLS secured web service URL (e.g., https://example.pulsar.io:8443), used when resource=cluster and operation is create or update"),
		),
		mcp.WithString("broker_service_url",
			mcp.Description("Pulsar cluster broker service URL (e.g., pulsar://example.pulsar.io:6650), used when resource=cluster and operation is create or update"),
		),
		mcp.WithString("broker_service_url_tls",
			mcp.Description("Pulsar cluster TLS secured broker service URL (e.g., pulsar+ssl://example.pulsar.io:6651), used when resource=cluster and operation is create or update"),
		),
		mcp.WithArray("peer_cluster_names",
			mcp.Description("List of clusters to be registered as peer-clusters, used when:\n"+
				"- resource=cluster and operation is create or update\n"+
				"- resource=peer_clusters and operation is update"),
			mcp.Items(
				map[string]interface{}{
					"type":        "string",
					"description": "peer cluster name",
				},
			),
		),
		mcp.WithArray("brokers",
			mcp.Description("List of broker names to include in a failure domain, required when resource=failure_domain and operation is create or update"),
			mcp.Items(
				map[string]interface{}{
					"type":        "string",
					"description": "broker",
				},
			),
		),
	)
}

// buildClusterHandler builds the Pulsar Admin Cluster handler function
// Migrated from the original handler logic
func (b *PulsarAdminClusterToolBuilder) buildClusterHandler(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return mcp.NewToolResultError("Pulsar session not found in context"), nil
		}

		client, err := session.GetAdminClient()
		if err != nil {
			return b.handleError("get admin client", err), nil
		}

		// Get required parameters
		resource, err := request.RequireString("resource")
		if err != nil {
			return mcp.NewToolResultError("Missing required resource parameter. " +
				"Please specify one of: cluster, peer_clusters, failure_domain."), nil
		}

		operation, err := request.RequireString("operation")
		if err != nil {
			return mcp.NewToolResultError("Missing required operation parameter. " +
				"Please specify one of: list, get, create, update, delete based on the resource type."), nil
		}

		// Validate if the parameter combination is valid
		validCombination, errMsg := b.validateClusterResourceOperation(resource, operation)
		if !validCombination {
			return mcp.NewToolResultError(errMsg), nil
		}

		// Check write operation permissions
		if (operation == "create" || operation == "update" || operation == "delete") && readOnly {
			return mcp.NewToolResultError("Create/update/delete operations not allowed in read-only mode. " +
				"Please contact your administrator if you need to modify cluster resources."), nil
		}

		// Process request based on resource type
		switch resource {
		case "cluster":
			return b.handleClusterResource(client, operation, request)
		case "peer_clusters":
			return b.handlePeerClustersResource(client, operation, request)
		case "failure_domain":
			return b.handleFailureDomainResource(client, operation, request)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Unsupported resource: %s. "+
				"Please use one of: cluster, peer_clusters, failure_domain.", resource)), nil
		}
	}
}

// Unified error handling and utility functions

// handleError provides unified error handling
func (b *PulsarAdminClusterToolBuilder) handleError(operation string, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("Failed to %s: %v", operation, err))
}

// marshalResponse provides unified JSON serialization for responses
func (b *PulsarAdminClusterToolBuilder) marshalResponse(data interface{}) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return b.handleError("marshal response", err), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// Validate if the resource and operation combination is valid
func (b *PulsarAdminClusterToolBuilder) validateClusterResourceOperation(resource, operation string) (bool, string) {
	validCombinations := map[string][]string{
		"cluster":        {"list", "get", "create", "update", "delete"},
		"peer_clusters":  {"get", "update"},
		"failure_domain": {"list", "get", "create", "update", "delete"},
	}

	if operations, ok := validCombinations[resource]; ok {
		for _, op := range operations {
			if op == operation {
				return true, ""
			}
		}
		return false, fmt.Sprintf("Invalid operation '%s' for resource '%s'. Valid operations are: %v",
			operation, resource, validCombinations[resource])
	}

	return false, fmt.Sprintf("Invalid resource '%s'. Valid resources are: cluster, peer_clusters, failure_domain", resource)
}

// Handle cluster resource operations
func (b *PulsarAdminClusterToolBuilder) handleClusterResource(client cmdutils.Client, operation string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	switch operation {
	case "list":
		return b.handleClusterList(client)
	case "get":
		clusterName, err := request.RequireString("cluster_name")
		if err != nil {
			return mcp.NewToolResultError("Missing required parameter 'cluster_name'. " +
				"Please provide the name of the cluster to get information for."), nil
		}
		return b.getClusterData(client, clusterName)
	case "create":
		return b.createCluster(client, request)
	case "update":
		return b.updateCluster(client, request)
	case "delete":
		clusterName, err := request.RequireString("cluster_name")
		if err != nil {
			return mcp.NewToolResultError("Missing required parameter 'cluster_name'. " +
				"Please provide the name of the cluster to delete."), nil
		}
		return b.deleteCluster(client, clusterName)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unsupported cluster operation: %s", operation)), nil
	}
}

// Handle peer clusters resource operations
func (b *PulsarAdminClusterToolBuilder) handlePeerClustersResource(client cmdutils.Client, operation string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clusterName, err := request.RequireString("cluster_name")
	if err != nil {
		return mcp.NewToolResultError("Missing required parameter 'cluster_name'. " +
			"Please provide the name of the cluster for peer clusters operation."), nil
	}

	switch operation {
	case "get":
		return b.getPeerClusters(client, clusterName)
	case "update":
		return b.updatePeerClusters(client, clusterName, request)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unsupported peer_clusters operation: %s", operation)), nil
	}
}

// Handle failure domain resource operations
func (b *PulsarAdminClusterToolBuilder) handleFailureDomainResource(client cmdutils.Client, operation string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clusterName, err := request.RequireString("cluster_name")
	if err != nil {
		return mcp.NewToolResultError("Missing required parameter 'cluster_name'. " +
			"Please provide the name of the cluster for failure domain operation."), nil
	}

	switch operation {
	case "list":
		return b.listFailureDomains(client, clusterName)
	case "get":
		domainName, err := request.RequireString("domain_name")
		if err != nil {
			return mcp.NewToolResultError("Missing required parameter 'domain_name'. " +
				"Please provide the name of the failure domain."), nil
		}
		return b.getFailureDomain(client, clusterName, domainName)
	case "create":
		return b.createFailureDomain(client, clusterName, request)
	case "update":
		return b.updateFailureDomain(client, clusterName, request)
	case "delete":
		domainName, err := request.RequireString("domain_name")
		if err != nil {
			return mcp.NewToolResultError("Missing required parameter 'domain_name'. " +
				"Please provide the name of the failure domain to delete."), nil
		}
		return b.deleteFailureDomain(client, clusterName, domainName)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unsupported failure_domain operation: %s", operation)), nil
	}
}

func (b *PulsarAdminClusterToolBuilder) handleClusterList(client cmdutils.Client) (*mcp.CallToolResult, error) {
	// Get cluster list
	clusters, err := client.Clusters().List()
	if err != nil {
		return b.handleError("get cluster list", err), nil
	}

	return b.marshalResponse(clusters)
}

func (b *PulsarAdminClusterToolBuilder) getClusterData(client cmdutils.Client, clusterName string) (*mcp.CallToolResult, error) {
	// Get cluster data
	clusterData, err := client.Clusters().Get(clusterName)
	if err != nil {
		return b.handleError("get cluster data", err), nil
	}

	return b.marshalResponse(clusterData)
}

func (b *PulsarAdminClusterToolBuilder) createCluster(client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clusterName, err := request.RequireString("cluster_name")
	if err != nil {
		return mcp.NewToolResultError("Missing required parameter 'cluster_name'. " +
			"Please provide the name of the cluster to create."), nil
	}

	// Initialize cluster data
	clusterData := utils.ClusterData{
		Name: clusterName,
	}

	// Set optional parameters if provided
	if serviceURL := request.GetString("service_url", ""); serviceURL != "" {
		clusterData.ServiceURL = serviceURL
	}
	if serviceURLTls := request.GetString("service_url_tls", ""); serviceURLTls != "" {
		clusterData.ServiceURLTls = serviceURLTls
	}
	if brokerServiceURL := request.GetString("broker_service_url", ""); brokerServiceURL != "" {
		clusterData.BrokerServiceURL = brokerServiceURL
	}
	if brokerServiceURLTls := request.GetString("broker_service_url_tls", ""); brokerServiceURLTls != "" {
		clusterData.BrokerServiceURLTls = brokerServiceURLTls
	}
	if peerClusters := request.GetStringSlice("peer_cluster_names", []string{}); len(peerClusters) > 0 {
		clusterData.PeerClusterNames = peerClusters
	}

	// Create cluster
	err = client.Clusters().Create(clusterData)
	if err != nil {
		return b.handleError("create cluster", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Cluster %s created successfully", clusterName)), nil
}

func (b *PulsarAdminClusterToolBuilder) updateCluster(client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clusterName, err := request.RequireString("cluster_name")
	if err != nil {
		return mcp.NewToolResultError("Missing required parameter 'cluster_name'. " +
			"Please provide the name of the cluster to update."), nil
	}

	// Initialize cluster data
	clusterData := utils.ClusterData{
		Name: clusterName,
	}

	// Set optional parameters if provided
	if serviceURL := request.GetString("service_url", ""); serviceURL != "" {
		clusterData.ServiceURL = serviceURL
	}
	if serviceURLTls := request.GetString("service_url_tls", ""); serviceURLTls != "" {
		clusterData.ServiceURLTls = serviceURLTls
	}
	if brokerServiceURL := request.GetString("broker_service_url", ""); brokerServiceURL != "" {
		clusterData.BrokerServiceURL = brokerServiceURL
	}
	if brokerServiceURLTls := request.GetString("broker_service_url_tls", ""); brokerServiceURLTls != "" {
		clusterData.BrokerServiceURLTls = brokerServiceURLTls
	}
	if peerClusters := request.GetStringSlice("peer_cluster_names", []string{}); len(peerClusters) > 0 {
		clusterData.PeerClusterNames = peerClusters
	}

	// Update cluster
	err = client.Clusters().Update(clusterData)
	if err != nil {
		return b.handleError("update cluster", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Cluster %s updated successfully", clusterName)), nil
}

func (b *PulsarAdminClusterToolBuilder) deleteCluster(client cmdutils.Client, clusterName string) (*mcp.CallToolResult, error) {
	// Delete cluster
	err := client.Clusters().Delete(clusterName)
	if err != nil {
		return b.handleError("delete cluster", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Cluster %s deleted successfully", clusterName)), nil
}

func (b *PulsarAdminClusterToolBuilder) getPeerClusters(client cmdutils.Client, clusterName string) (*mcp.CallToolResult, error) {
	// Get peer clusters
	peerClusters, err := client.Clusters().GetPeerClusters(clusterName)
	if err != nil {
		return b.handleError("get peer clusters", err), nil
	}

	return b.marshalResponse(peerClusters)
}

func (b *PulsarAdminClusterToolBuilder) updatePeerClusters(client cmdutils.Client, clusterName string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	peerClusters, err := request.RequireStringSlice("peer_cluster_names")
	if err != nil {
		return mcp.NewToolResultError("Missing required parameter 'peer_cluster_names'. " +
			"Please provide an array of peer cluster names to set."), nil
	}

	// Update peer clusters
	err = client.Clusters().UpdatePeerClusters(clusterName, peerClusters)
	if err != nil {
		return b.handleError("update peer clusters", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Peer clusters for %s updated successfully", clusterName)), nil
}

func (b *PulsarAdminClusterToolBuilder) listFailureDomains(client cmdutils.Client, clusterName string) (*mcp.CallToolResult, error) {
	// Get failure domains list
	failureDomains, err := client.Clusters().ListFailureDomains(clusterName)
	if err != nil {
		return b.handleError("list failure domains", err), nil
	}

	return b.marshalResponse(failureDomains)
}

func (b *PulsarAdminClusterToolBuilder) getFailureDomain(client cmdutils.Client, clusterName, domainName string) (*mcp.CallToolResult, error) {
	// Get failure domain
	failureDomain, err := client.Clusters().GetFailureDomain(clusterName, domainName)
	if err != nil {
		return b.handleError("get failure domain", err), nil
	}

	return b.marshalResponse(failureDomain)
}

func (b *PulsarAdminClusterToolBuilder) createFailureDomain(client cmdutils.Client, clusterName string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	domainName, err := request.RequireString("domain_name")
	if err != nil {
		return mcp.NewToolResultError("Missing required parameter 'domain_name'. " +
			"Please provide the name of the failure domain to create."), nil
	}

	brokers, err := request.RequireStringSlice("brokers")
	if err != nil {
		return mcp.NewToolResultError("Missing required parameter 'brokers'. " +
			"Please provide an array of broker names to include in this failure domain."), nil
	}

	// Create failure domain data
	failureDomainData := utils.FailureDomainData{
		ClusterName: clusterName,
		DomainName:  domainName,
		BrokerList:  brokers,
	}

	// Create failure domain
	err = client.Clusters().CreateFailureDomain(failureDomainData)
	if err != nil {
		return b.handleError("create failure domain", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Failure domain %s created successfully in cluster %s", domainName, clusterName)), nil
}

func (b *PulsarAdminClusterToolBuilder) updateFailureDomain(client cmdutils.Client, clusterName string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	domainName, err := request.RequireString("domain_name")
	if err != nil {
		return mcp.NewToolResultError("Missing required parameter 'domain_name'. " +
			"Please provide the name of the failure domain to update."), nil
	}

	brokers, err := request.RequireStringSlice("brokers")
	if err != nil {
		return mcp.NewToolResultError("Missing required parameter 'brokers'. " +
			"Please provide an array of broker names to include in this failure domain."), nil
	}

	// Create failure domain data
	failureDomainData := utils.FailureDomainData{
		ClusterName: clusterName,
		DomainName:  domainName,
		BrokerList:  brokers,
	}

	// Update failure domain
	err = client.Clusters().UpdateFailureDomain(failureDomainData)
	if err != nil {
		return b.handleError("update failure domain", err), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Failure domain %s updated successfully in cluster %s", domainName, clusterName)), nil
}

func (b *PulsarAdminClusterToolBuilder) deleteFailureDomain(client cmdutils.Client, clusterName, domainName string) (*mcp.CallToolResult, error) {
	// Create failure domain data for deletion
	failureDomainData := utils.FailureDomainData{
		ClusterName: clusterName,
		DomainName:  domainName,
	}

	// Delete failure domain
	err := client.Clusters().DeleteFailureDomain(failureDomainData)
	if err != nil {
		return b.handleError("delete failure domain", err), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Failure domain %s deleted successfully from cluster %s", domainName, clusterName)), nil
}
