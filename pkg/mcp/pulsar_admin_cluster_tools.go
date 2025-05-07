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

package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// PulsarAdminAddClusterTools adds cluster-related tools to the MCP server
func PulsarAdminAddClusterTools(s *server.MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeaturePulsarAdminClusters)) && !slices.Contains(features, string(FeatureAll)) && !slices.Contains(features, string(FeatureAllPulsar)) {
		return
	}
	// Create a single unified cluster tool
	clusterTool := mcp.NewTool("pulsar_admin_cluster",
		mcp.WithDescription("Unified tool for managing Apache Pulsar clusters.\n"+
			"This tool provides access to various cluster resources and operations, including:\n"+
			"1. Manage clusters (resource=cluster): List, get, create, update, delete clusters\n"+
			"2. Manage peer clusters (resource=peer_clusters): Get, update peer clusters\n"+
			"3. Manage failure domains (resource=failure_domain): List, get, create, update, delete failure domains\n\n"+
			"Different functions are accessed by combining resource and operation parameters, with other parameters used selectively based on operation type.\n\n"+
			"Examples:\n"+
			"- {\"resource\": \"cluster\", \"operation\": \"list\"} lists all clusters\n"+
			"- {\"resource\": \"cluster\", \"operation\": \"get\", \"cluster_name\": \"my-cluster\"} gets cluster configuration\n"+
			"- {\"resource\": \"failure_domain\", \"operation\": \"list\", \"cluster_name\": \"my-cluster\"} lists failure domains\n"+
			"This tool requires Pulsar super-user permissions."),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description("Type of cluster resource to access, available options:\n"+
				"- cluster: Pulsar cluster configuration\n"+
				"- peer_clusters: Peer clusters for geo-replication\n"+
				"- failure_domain: Failure domains for fault tolerance"),
		),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description("Operation to perform, available options (depend on resource):\n"+
				"- list: List resources (used with cluster, failure_domain)\n"+
				"- get: Retrieve resource information (used with cluster, peer_clusters, failure_domain)\n"+
				"- create: Create a new resource (used with cluster, failure_domain)\n"+
				"- update: Update an existing resource (used with cluster, peer_clusters, failure_domain)\n"+
				"- delete: Delete a resource (used with cluster, failure_domain)"),
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
		),
		mcp.WithArray("brokers",
			mcp.Description("List of broker names to include in a failure domain, required when resource=failure_domain and operation is create or update"),
		),
	)
	s.AddTool(clusterTool, handleClusterTool(readOnly))
}

// handleClusterTool returns a function to handle cluster tool requests
func handleClusterTool(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := pulsar.GetAdminClient()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
		}

		// Get required parameters
		resource, err := requiredParam[string](request.Params.Arguments, "resource")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required resource parameter. " +
				"Please specify one of: cluster, peer_clusters, failure_domain.")), nil
		}

		operation, err := requiredParam[string](request.Params.Arguments, "operation")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required operation parameter. " +
				"Please specify one of: list, get, create, update, delete based on the resource type.")), nil
		}

		// Validate if the parameter combination is valid
		validCombination, errMsg := validateClusterResourceOperation(resource, operation)
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
			return handleClusterResource(client, operation, request)
		case "peer_clusters":
			return handlePeerClustersResource(client, operation, request)
		case "failure_domain":
			return handleFailureDomainResource(client, operation, request)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Unsupported resource: %s. "+
				"Please use one of: cluster, peer_clusters, failure_domain.", resource)), nil
		}
	}
}

// Validate if the resource and operation combination is valid
func validateClusterResourceOperation(resource, operation string) (bool, string) {
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
func handleClusterResource(client cmdutils.Client, operation string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	switch operation {
	case "list":
		return handleClusterList(client)
	case "get":
		clusterName, err := requiredParam[string](request.Params.Arguments, "cluster_name")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'cluster_name'. " +
				"Please provide the name of the cluster to get information for.")), nil
		}
		return getClusterData(client, clusterName)
	case "create":
		return createCluster(client, request)
	case "update":
		return updateCluster(client, request)
	case "delete":
		clusterName, err := requiredParam[string](request.Params.Arguments, "cluster_name")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'cluster_name'. " +
				"Please provide the name of the cluster to delete.")), nil
		}
		return deleteCluster(client, clusterName)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unsupported operation '%s' for cluster resource. "+
			"Supported operations are: list, get, create, update, delete.", operation)), nil
	}
}

// Handle peer_clusters resource operations
func handlePeerClustersResource(client cmdutils.Client, operation string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clusterName, err := requiredParam[string](request.Params.Arguments, "cluster_name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'cluster_name'. " +
			"Please provide the name of the cluster to operate on.")), nil
	}

	switch operation {
	case "get":
		return getPeerClusters(client, clusterName)
	case "update":
		peerClusters, err := requiredParamArray[string](request.Params.Arguments, "peer_cluster_names")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'peer_cluster_names'. " +
				"Please provide an array of peer cluster names to set.")), nil
		}
		return updatePeerClusters(client, clusterName, peerClusters)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unsupported operation '%s' for peer_clusters resource. "+
			"Supported operations are: get, update.", operation)), nil
	}
}

// Handle failure_domain resource operations
func handleFailureDomainResource(client cmdutils.Client, operation string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clusterName, err := requiredParam[string](request.Params.Arguments, "cluster_name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'cluster_name'. " +
			"Please provide the name of the cluster to operate on.")), nil
	}

	switch operation {
	case "list":
		return listFailureDomains(client, clusterName)
	case "get":
		domainName, err := requiredParam[string](request.Params.Arguments, "domain_name")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'domain_name'. " +
				"Please provide the name of the failure domain to get.")), nil
		}
		return getFailureDomain(client, clusterName, domainName)
	case "create":
		return createFailureDomain(client, request)
	case "update":
		return updateFailureDomain(client, request)
	case "delete":
		domainName, err := requiredParam[string](request.Params.Arguments, "domain_name")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'domain_name'. " +
				"Please provide the name of the failure domain to delete.")), nil
		}
		return deleteFailureDomain(client, clusterName, domainName)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unsupported operation '%s' for failure_domain resource. "+
			"Supported operations are: list, get, create, update, delete.", operation)), nil
	}
}

// handleClusterList handles listing all clusters
func handleClusterList(client cmdutils.Client) (*mcp.CallToolResult, error) {
	// Get cluster list
	clusters, err := client.Clusters().List()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get cluster list: %v", err)), nil
	}

	// Convert result to JSON string
	clustersJSON, err := json.Marshal(clusters)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize cluster list: %v", err)), nil
	}

	return mcp.NewToolResultText(string(clustersJSON)), nil
}

// getClusterData gets the cluster data
func getClusterData(client cmdutils.Client, clusterName string) (*mcp.CallToolResult, error) {
	clusterData, err := client.Clusters().Get(clusterName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get cluster data: %v", err)), nil
	}

	// Convert result to JSON string
	clusterDataJSON, err := json.Marshal(clusterData)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize cluster data: %v", err)), nil
	}

	return mcp.NewToolResultText(string(clusterDataJSON)), nil
}

// createCluster creates a new Pulsar cluster
func createCluster(client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clusterName, err := requiredParam[string](request.Params.Arguments, "cluster_name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'cluster_name'. " +
			"Please provide the name of the cluster to create.")), nil
	}

	clusterData := utils.ClusterData{
		Name: clusterName,
	}

	// Set optional parameters if provided
	if serviceURL, ok := optionalParam[string](request.Params.Arguments, "service_url"); ok {
		clusterData.ServiceURL = serviceURL
	}
	if serviceURLTls, ok := optionalParam[string](request.Params.Arguments, "service_url_tls"); ok {
		clusterData.ServiceURLTls = serviceURLTls
	}
	if brokerServiceURL, ok := optionalParam[string](request.Params.Arguments, "broker_service_url"); ok {
		clusterData.BrokerServiceURL = brokerServiceURL
	}
	if brokerServiceURLTls, ok := optionalParam[string](request.Params.Arguments, "broker_service_url_tls"); ok {
		clusterData.BrokerServiceURLTls = brokerServiceURLTls
	}
	if peerClusters, ok := optionalParamArray[string](request.Params.Arguments, "peer_cluster_names"); ok {
		clusterData.PeerClusterNames = peerClusters
	}

	err = client.Clusters().Create(clusterData)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create cluster: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Cluster %s created successfully", clusterName)), nil
}

// updateCluster updates an existing Pulsar cluster
func updateCluster(client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clusterName, err := requiredParam[string](request.Params.Arguments, "cluster_name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'cluster_name'. " +
			"Please provide the name of the cluster to update.")), nil
	}

	clusterData := utils.ClusterData{
		Name: clusterName,
	}

	// Set optional parameters if provided
	if serviceURL, ok := optionalParam[string](request.Params.Arguments, "service_url"); ok {
		clusterData.ServiceURL = serviceURL
	}
	if serviceURLTls, ok := optionalParam[string](request.Params.Arguments, "service_url_tls"); ok {
		clusterData.ServiceURLTls = serviceURLTls
	}
	if brokerServiceURL, ok := optionalParam[string](request.Params.Arguments, "broker_service_url"); ok {
		clusterData.BrokerServiceURL = brokerServiceURL
	}
	if brokerServiceURLTls, ok := optionalParam[string](request.Params.Arguments, "broker_service_url_tls"); ok {
		clusterData.BrokerServiceURLTls = brokerServiceURLTls
	}
	if peerClusters, ok := optionalParamArray[string](request.Params.Arguments, "peer_cluster_names"); ok {
		clusterData.PeerClusterNames = peerClusters
	}

	err = client.Clusters().Update(clusterData)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to update cluster: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Cluster %s updated successfully", clusterName)), nil
}

// deleteCluster deletes an existing Pulsar cluster
func deleteCluster(client cmdutils.Client, clusterName string) (*mcp.CallToolResult, error) {
	err := client.Clusters().Delete(clusterName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete cluster: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Cluster %s deleted successfully", clusterName)), nil
}

// getPeerClusters gets the peer clusters for a specified cluster
func getPeerClusters(client cmdutils.Client, clusterName string) (*mcp.CallToolResult, error) {
	peerClusters, err := client.Clusters().GetPeerClusters(clusterName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get peer clusters: %v", err)), nil
	}

	// Convert result to JSON string
	peerClustersJSON, err := json.Marshal(peerClusters)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize peer clusters: %v", err)), nil
	}

	return mcp.NewToolResultText(string(peerClustersJSON)), nil
}

// updatePeerClusters updates the peer clusters for a specified cluster
func updatePeerClusters(client cmdutils.Client, clusterName string, peerClusters []string) (*mcp.CallToolResult, error) {
	err := client.Clusters().UpdatePeerClusters(clusterName, peerClusters)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to update peer clusters: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Peer clusters for %s updated successfully", clusterName)), nil
}

// getFailureDomain gets the configuration of a specific failure domain
func getFailureDomain(client cmdutils.Client, clusterName string, domainName string) (*mcp.CallToolResult, error) {
	failureDomain, err := client.Clusters().GetFailureDomain(clusterName, domainName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get failure domain: %v", err)), nil
	}

	// Convert result to JSON string
	failureDomainJSON, err := json.Marshal(failureDomain)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize failure domain: %v", err)), nil
	}

	return mcp.NewToolResultText(string(failureDomainJSON)), nil
}

// listFailureDomains gets the list of failure domains for a specified cluster
func listFailureDomains(client cmdutils.Client, clusterName string) (*mcp.CallToolResult, error) {
	failureDomains, err := client.Clusters().ListFailureDomains(clusterName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list failure domains: %v", err)), nil
	}

	// Convert result to JSON string
	failureDomainsJSON, err := json.Marshal(failureDomains)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize failure domains: %v", err)), nil
	}

	return mcp.NewToolResultText(string(failureDomainsJSON)), nil
}

// createFailureDomain creates a new failure domain in the specified cluster
func createFailureDomain(client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clusterName, err := requiredParam[string](request.Params.Arguments, "cluster_name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'cluster_name'. " +
			"Please provide the name of the cluster.")), nil
	}

	domainName, err := requiredParam[string](request.Params.Arguments, "domain_name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'domain_name'. " +
			"Please provide the name of the failure domain to create.")), nil
	}

	brokers, err := requiredParamArray[string](request.Params.Arguments, "brokers")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'brokers'. " +
			"Please provide an array of broker names to include in this failure domain.")), nil
	}

	failureDomain := utils.FailureDomainData{
		ClusterName: clusterName,
		DomainName:  domainName,
		BrokerList:  brokers,
	}

	err = client.Clusters().CreateFailureDomain(failureDomain)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create failure domain: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Failure domain %s created successfully in cluster %s", domainName, clusterName)), nil
}

// updateFailureDomain updates an existing failure domain in the specified cluster
func updateFailureDomain(client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clusterName, err := requiredParam[string](request.Params.Arguments, "cluster_name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'cluster_name'. " +
			"Please provide the name of the cluster.")), nil
	}

	domainName, err := requiredParam[string](request.Params.Arguments, "domain_name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'domain_name'. " +
			"Please provide the name of the failure domain to update.")), nil
	}

	brokers, err := requiredParamArray[string](request.Params.Arguments, "brokers")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'brokers'. " +
			"Please provide an array of broker names to include in this failure domain.")), nil
	}

	failureDomain := utils.FailureDomainData{
		ClusterName: clusterName,
		DomainName:  domainName,
		BrokerList:  brokers,
	}

	err = client.Clusters().UpdateFailureDomain(failureDomain)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to update failure domain: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Failure domain %s updated successfully in cluster %s", domainName, clusterName),
	), nil
}

// deleteFailureDomain deletes a failure domain from the specified cluster
func deleteFailureDomain(client cmdutils.Client, clusterName string, domainName string) (*mcp.CallToolResult, error) {
	failureDomain := utils.FailureDomainData{
		ClusterName: clusterName,
		DomainName:  domainName,
	}

	err := client.Clusters().DeleteFailureDomain(failureDomain)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete failure domain: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf(
		"Failure domain %s deleted successfully from cluster %s", domainName, clusterName)), nil
}
