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

package pulsar

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

type pulsarAdminClusterInput struct {
	Resource            string   `json:"resource"`
	Operation           string   `json:"operation"`
	ClusterName         *string  `json:"cluster_name,omitempty"`
	DomainName          *string  `json:"domain_name,omitempty"`
	ServiceURL          *string  `json:"service_url,omitempty"`
	ServiceURLTLS       *string  `json:"service_url_tls,omitempty"`
	BrokerServiceURL    *string  `json:"broker_service_url,omitempty"`
	BrokerServiceURLTLS *string  `json:"broker_service_url_tls,omitempty"`
	PeerClusterNames    []string `json:"peer_cluster_names,omitempty"`
	Brokers             []string `json:"brokers,omitempty"`
}

const (
	pulsarAdminClusterResourceDesc = "Type of cluster resource to access, available options:\n" +
		"- cluster: Pulsar cluster configuration\n" +
		"- peer_clusters: Peer clusters for geo-replication\n" +
		"- failure_domain: Failure domains for fault tolerance"
	pulsarAdminClusterOperationDesc = "Operation to perform, available options (depend on resource):\n" +
		"- list: List resources (used with cluster, failure_domain)\n" +
		"- get: Retrieve resource information (used with cluster, peer_clusters, failure_domain)\n" +
		"- create: Create a new resource (used with cluster, failure_domain)\n" +
		"- update: Update an existing resource (used with cluster, peer_clusters, failure_domain)\n" +
		"- delete: Delete a resource (used with cluster, failure_domain)"
	pulsarAdminClusterNameDesc                = "Name of the Pulsar cluster, required for all operations except 'list' with resource=cluster"
	pulsarAdminClusterDomainNameDesc          = "Name of the failure domain, required when resource=failure_domain and operation is get, create, update, or delete"
	pulsarAdminClusterServiceURLDesc          = "Pulsar cluster web service URL (e.g., http://example.pulsar.io:8080), used when resource=cluster and operation is create or update"
	pulsarAdminClusterServiceURLTLSDesc       = "Pulsar cluster TLS secured web service URL (e.g., https://example.pulsar.io:8443), used when resource=cluster and operation is create or update"
	pulsarAdminClusterBrokerServiceURLDesc    = "Pulsar cluster broker service URL (e.g., pulsar://example.pulsar.io:6650), used when resource=cluster and operation is create or update"
	pulsarAdminClusterBrokerServiceURLTLSDesc = "Pulsar cluster TLS secured broker service URL (e.g., pulsar+ssl://example.pulsar.io:6651), used when resource=cluster and operation is create or update"
	pulsarAdminClusterPeerClusterNamesDesc    = "List of clusters to be registered as peer-clusters, used when:\n" +
		"- resource=cluster and operation is create or update\n" +
		"- resource=peer_clusters and operation is update"
	pulsarAdminClusterBrokersDesc = "List of broker names to include in a failure domain, required when resource=failure_domain and operation is create or update"
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
func (b *PulsarAdminClusterToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool, err := b.buildClusterTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildClusterHandler(config.ReadOnly)

	return []builders.ToolDefinition{
		builders.ServerTool[pulsarAdminClusterInput, any]{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildClusterTool builds the Pulsar Admin Cluster MCP tool definition
// Migrated from the original tool definition logic
func (b *PulsarAdminClusterToolBuilder) buildClusterTool() (*sdk.Tool, error) {
	inputSchema, err := buildPulsarAdminClusterInputSchema()
	if err != nil {
		return nil, err
	}

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

	return &sdk.Tool{
		Name:        "pulsar_admin_cluster",
		Description: toolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildClusterHandler builds the Pulsar Admin Cluster handler function
// Migrated from the original handler logic
func (b *PulsarAdminClusterToolBuilder) buildClusterHandler(readOnly bool) builders.ToolHandlerFunc[pulsarAdminClusterInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input pulsarAdminClusterInput) (*sdk.CallToolResult, any, error) {
		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return nil, nil, fmt.Errorf("pulsar session not found in context")
		}

		client, err := session.GetAdminClient()
		if err != nil {
			return nil, nil, b.handleError("get admin client", err)
		}

		resource := input.Resource
		if resource == "" {
			return nil, nil, fmt.Errorf("missing required resource parameter. please specify one of: cluster, peer_clusters, failure_domain")
		}

		operation := input.Operation
		if operation == "" {
			return nil, nil, fmt.Errorf("missing required operation parameter. please specify one of: list, get, create, update, delete based on the resource type")
		}

		// Validate if the parameter combination is valid
		validCombination, errMsg := b.validateClusterResourceOperation(resource, operation)
		if !validCombination {
			return nil, nil, fmt.Errorf("%s", errMsg)
		}

		// Check write operation permissions
		if (operation == "create" || operation == "update" || operation == "delete") && readOnly {
			return nil, nil, fmt.Errorf("create/update/delete operations not allowed in read-only mode. please contact your administrator if you need to modify cluster resources")
		}

		// Process request based on resource type
		switch resource {
		case "cluster":
			result, err := b.handleClusterResource(client, operation, input)
			return result, nil, err
		case "peer_clusters":
			result, err := b.handlePeerClustersResource(client, operation, input)
			return result, nil, err
		case "failure_domain":
			result, err := b.handleFailureDomainResource(client, operation, input)
			return result, nil, err
		default:
			return nil, nil, fmt.Errorf("unsupported resource: %s. please use one of: cluster, peer_clusters, failure_domain", resource)
		}
	}
}

// Unified error handling and utility functions

// handleError provides unified error handling
func (b *PulsarAdminClusterToolBuilder) handleError(operation string, err error) error {
	return fmt.Errorf("failed to %s: %v", operation, err)
}

// marshalResponse provides unified JSON serialization for responses
func (b *PulsarAdminClusterToolBuilder) marshalResponse(data interface{}) (*sdk.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, b.handleError("marshal response", err)
	}
	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: string(jsonBytes)}},
	}, nil
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
func (b *PulsarAdminClusterToolBuilder) handleClusterResource(client cmdutils.Client, operation string, input pulsarAdminClusterInput) (*sdk.CallToolResult, error) {
	switch operation {
	case "list":
		return b.handleClusterList(client)
	case "get":
		clusterName, err := requireString(input.ClusterName, "cluster_name")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'cluster_name'. please provide the name of the cluster to get information for")
		}
		return b.getClusterData(client, clusterName)
	case "create":
		return b.createCluster(client, input)
	case "update":
		return b.updateCluster(client, input)
	case "delete":
		clusterName, err := requireString(input.ClusterName, "cluster_name")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'cluster_name'. please provide the name of the cluster to delete")
		}
		return b.deleteCluster(client, clusterName)
	default:
		return nil, fmt.Errorf("unsupported cluster operation: %s", operation)
	}
}

// Handle peer clusters resource operations
func (b *PulsarAdminClusterToolBuilder) handlePeerClustersResource(client cmdutils.Client, operation string, input pulsarAdminClusterInput) (*sdk.CallToolResult, error) {
	clusterName, err := requireString(input.ClusterName, "cluster_name")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'cluster_name'. please provide the name of the cluster for peer clusters operation")
	}

	switch operation {
	case "get":
		return b.getPeerClusters(client, clusterName)
	case "update":
		return b.updatePeerClusters(client, clusterName, input)
	default:
		return nil, fmt.Errorf("unsupported peer_clusters operation: %s", operation)
	}
}

// Handle failure domain resource operations
func (b *PulsarAdminClusterToolBuilder) handleFailureDomainResource(client cmdutils.Client, operation string, input pulsarAdminClusterInput) (*sdk.CallToolResult, error) {
	clusterName, err := requireString(input.ClusterName, "cluster_name")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'cluster_name'. please provide the name of the cluster for failure domain operation")
	}

	switch operation {
	case "list":
		return b.listFailureDomains(client, clusterName)
	case "get":
		domainName, err := requireString(input.DomainName, "domain_name")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'domain_name'. please provide the name of the failure domain")
		}
		return b.getFailureDomain(client, clusterName, domainName)
	case "create":
		return b.createFailureDomain(client, clusterName, input)
	case "update":
		return b.updateFailureDomain(client, clusterName, input)
	case "delete":
		domainName, err := requireString(input.DomainName, "domain_name")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'domain_name'. please provide the name of the failure domain to delete")
		}
		return b.deleteFailureDomain(client, clusterName, domainName)
	default:
		return nil, fmt.Errorf("unsupported failure_domain operation: %s", operation)
	}
}

func (b *PulsarAdminClusterToolBuilder) handleClusterList(client cmdutils.Client) (*sdk.CallToolResult, error) {
	// Get cluster list
	clusters, err := client.Clusters().List()
	if err != nil {
		return nil, b.handleError("get cluster list", err)
	}

	return b.marshalResponse(clusters)
}

func (b *PulsarAdminClusterToolBuilder) getClusterData(client cmdutils.Client, clusterName string) (*sdk.CallToolResult, error) {
	// Get cluster data
	clusterData, err := client.Clusters().Get(clusterName)
	if err != nil {
		return nil, b.handleError("get cluster data", err)
	}

	return b.marshalResponse(clusterData)
}

func (b *PulsarAdminClusterToolBuilder) createCluster(client cmdutils.Client, input pulsarAdminClusterInput) (*sdk.CallToolResult, error) {
	clusterName, err := requireString(input.ClusterName, "cluster_name")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'cluster_name'. please provide the name of the cluster to create")
	}

	// Initialize cluster data
	clusterData := utils.ClusterData{
		Name: clusterName,
	}

	// Set optional parameters if provided
	if serviceURL := stringValue(input.ServiceURL); serviceURL != "" {
		clusterData.ServiceURL = serviceURL
	}
	if serviceURLTLS := stringValue(input.ServiceURLTLS); serviceURLTLS != "" {
		clusterData.ServiceURLTls = serviceURLTLS
	}
	if brokerServiceURL := stringValue(input.BrokerServiceURL); brokerServiceURL != "" {
		clusterData.BrokerServiceURL = brokerServiceURL
	}
	if brokerServiceURLTLS := stringValue(input.BrokerServiceURLTLS); brokerServiceURLTLS != "" {
		clusterData.BrokerServiceURLTls = brokerServiceURLTLS
	}
	if len(input.PeerClusterNames) > 0 {
		clusterData.PeerClusterNames = input.PeerClusterNames
	}

	// Create cluster
	err = client.Clusters().Create(clusterData)
	if err != nil {
		return nil, b.handleError("create cluster", err)
	}

	return textResult(fmt.Sprintf("Cluster %s created successfully", clusterName)), nil
}

func (b *PulsarAdminClusterToolBuilder) updateCluster(client cmdutils.Client, input pulsarAdminClusterInput) (*sdk.CallToolResult, error) {
	clusterName, err := requireString(input.ClusterName, "cluster_name")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'cluster_name'. please provide the name of the cluster to update")
	}

	// Initialize cluster data
	clusterData := utils.ClusterData{
		Name: clusterName,
	}

	// Set optional parameters if provided
	if serviceURL := stringValue(input.ServiceURL); serviceURL != "" {
		clusterData.ServiceURL = serviceURL
	}
	if serviceURLTLS := stringValue(input.ServiceURLTLS); serviceURLTLS != "" {
		clusterData.ServiceURLTls = serviceURLTLS
	}
	if brokerServiceURL := stringValue(input.BrokerServiceURL); brokerServiceURL != "" {
		clusterData.BrokerServiceURL = brokerServiceURL
	}
	if brokerServiceURLTLS := stringValue(input.BrokerServiceURLTLS); brokerServiceURLTLS != "" {
		clusterData.BrokerServiceURLTls = brokerServiceURLTLS
	}
	if len(input.PeerClusterNames) > 0 {
		clusterData.PeerClusterNames = input.PeerClusterNames
	}

	// Update cluster
	err = client.Clusters().Update(clusterData)
	if err != nil {
		return nil, b.handleError("update cluster", err)
	}

	return textResult(fmt.Sprintf("Cluster %s updated successfully", clusterName)), nil
}

func (b *PulsarAdminClusterToolBuilder) deleteCluster(client cmdutils.Client, clusterName string) (*sdk.CallToolResult, error) {
	// Delete cluster
	err := client.Clusters().Delete(clusterName)
	if err != nil {
		return nil, b.handleError("delete cluster", err)
	}

	return textResult(fmt.Sprintf("Cluster %s deleted successfully", clusterName)), nil
}

func (b *PulsarAdminClusterToolBuilder) getPeerClusters(client cmdutils.Client, clusterName string) (*sdk.CallToolResult, error) {
	// Get peer clusters
	peerClusters, err := client.Clusters().GetPeerClusters(clusterName)
	if err != nil {
		return nil, b.handleError("get peer clusters", err)
	}

	return b.marshalResponse(peerClusters)
}

func (b *PulsarAdminClusterToolBuilder) updatePeerClusters(client cmdutils.Client, clusterName string, input pulsarAdminClusterInput) (*sdk.CallToolResult, error) {
	peerClusters, err := requireStringSlice(input.PeerClusterNames, "peer_cluster_names")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'peer_cluster_names'. please provide an array of peer cluster names to set")
	}

	// Update peer clusters
	err = client.Clusters().UpdatePeerClusters(clusterName, peerClusters)
	if err != nil {
		return nil, b.handleError("update peer clusters", err)
	}

	return textResult(fmt.Sprintf("Peer clusters for %s updated successfully", clusterName)), nil
}

func (b *PulsarAdminClusterToolBuilder) listFailureDomains(client cmdutils.Client, clusterName string) (*sdk.CallToolResult, error) {
	// Get failure domains list
	failureDomains, err := client.Clusters().ListFailureDomains(clusterName)
	if err != nil {
		return nil, b.handleError("list failure domains", err)
	}

	return b.marshalResponse(failureDomains)
}

func (b *PulsarAdminClusterToolBuilder) getFailureDomain(client cmdutils.Client, clusterName, domainName string) (*sdk.CallToolResult, error) {
	// Get failure domain
	failureDomain, err := client.Clusters().GetFailureDomain(clusterName, domainName)
	if err != nil {
		return nil, b.handleError("get failure domain", err)
	}

	return b.marshalResponse(failureDomain)
}

func (b *PulsarAdminClusterToolBuilder) createFailureDomain(client cmdutils.Client, clusterName string, input pulsarAdminClusterInput) (*sdk.CallToolResult, error) {
	domainName, err := requireString(input.DomainName, "domain_name")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'domain_name'. please provide the name of the failure domain to create")
	}

	brokers, err := requireStringSlice(input.Brokers, "brokers")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'brokers'. please provide an array of broker names to include in this failure domain")
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
		return nil, b.handleError("create failure domain", err)
	}

	return textResult(fmt.Sprintf("Failure domain %s created successfully in cluster %s", domainName, clusterName)), nil
}

func (b *PulsarAdminClusterToolBuilder) updateFailureDomain(client cmdutils.Client, clusterName string, input pulsarAdminClusterInput) (*sdk.CallToolResult, error) {
	domainName, err := requireString(input.DomainName, "domain_name")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'domain_name'. please provide the name of the failure domain to update")
	}

	brokers, err := requireStringSlice(input.Brokers, "brokers")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'brokers'. please provide an array of broker names to include in this failure domain")
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
		return nil, b.handleError("update failure domain", err)
	}

	return textResult(fmt.Sprintf("Failure domain %s updated successfully in cluster %s", domainName, clusterName)), nil
}

func (b *PulsarAdminClusterToolBuilder) deleteFailureDomain(client cmdutils.Client, clusterName, domainName string) (*sdk.CallToolResult, error) {
	// Create failure domain data for deletion
	failureDomainData := utils.FailureDomainData{
		ClusterName: clusterName,
		DomainName:  domainName,
	}

	// Delete failure domain
	err := client.Clusters().DeleteFailureDomain(failureDomainData)
	if err != nil {
		return nil, b.handleError("delete failure domain", err)
	}

	return textResult(fmt.Sprintf("Failure domain %s deleted successfully from cluster %s", domainName, clusterName)), nil
}

func buildPulsarAdminClusterInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[pulsarAdminClusterInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}

	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}

	setSchemaDescription(schema, "resource", pulsarAdminClusterResourceDesc)
	setSchemaDescription(schema, "operation", pulsarAdminClusterOperationDesc)
	setSchemaDescription(schema, "cluster_name", pulsarAdminClusterNameDesc)
	setSchemaDescription(schema, "domain_name", pulsarAdminClusterDomainNameDesc)
	setSchemaDescription(schema, "service_url", pulsarAdminClusterServiceURLDesc)
	setSchemaDescription(schema, "service_url_tls", pulsarAdminClusterServiceURLTLSDesc)
	setSchemaDescription(schema, "broker_service_url", pulsarAdminClusterBrokerServiceURLDesc)
	setSchemaDescription(schema, "broker_service_url_tls", pulsarAdminClusterBrokerServiceURLTLSDesc)
	setSchemaDescription(schema, "peer_cluster_names", pulsarAdminClusterPeerClusterNamesDesc)
	setSchemaDescription(schema, "brokers", pulsarAdminClusterBrokersDesc)

	if peersSchema := schema.Properties["peer_cluster_names"]; peersSchema != nil && peersSchema.Items != nil {
		peersSchema.Items.Description = "peer cluster name"
	}
	if brokersSchema := schema.Properties["brokers"]; brokersSchema != nil && brokersSchema.Items != nil {
		brokersSchema.Items.Description = "broker"
	}

	normalizeAdditionalProperties(schema)
	return schema, nil
}
