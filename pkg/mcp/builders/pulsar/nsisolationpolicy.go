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
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/adapter"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

// PulsarAdminNsIsolationPolicyToolBuilder implements the ToolBuilder interface for Pulsar admin namespace isolation policies
// /nolint:revive
type PulsarAdminNsIsolationPolicyToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminNsIsolationPolicyToolBuilder creates a new Pulsar admin namespace isolation policy tool builder instance
func NewPulsarAdminNsIsolationPolicyToolBuilder() *PulsarAdminNsIsolationPolicyToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_nsisolationpolicy",
		Version:     "1.0.0",
		Description: "Pulsar admin namespace isolation policy management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "admin", "nsisolationpolicy"},
	}

	features := []string{
		"pulsar-admin-nsisolationpolicy",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminNsIsolationPolicyToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin namespace isolation policy tool list
func (b *PulsarAdminNsIsolationPolicyToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ServerTool, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool := b.buildNsIsolationPolicyTool()
	handler := b.buildNsIsolationPolicyHandler(config.ReadOnly)

	return []builders.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildNsIsolationPolicyTool builds the Pulsar admin namespace isolation policy MCP tool definition
func (b *PulsarAdminNsIsolationPolicyToolBuilder) buildNsIsolationPolicyTool() *mcpsdk.Tool {
	toolDesc := "Manage namespace isolation policies in a Pulsar cluster. " +
		"Allows viewing, creating, updating, and deleting namespace isolation policies. " +
		"Some operations require super-user permissions."

	resourceDesc := "Resource to operate on. Available resources:\n" +
		"- policy: Namespace isolation policy\n" +
		"- broker: Broker with namespace isolation policies\n" +
		"- brokers: All brokers with namespace isolation policies"

	operationDesc := "Operation to perform. Available operations:\n" +
		"- get: Get resource details\n" +
		"- list: List all instances of the resource\n" +
		"- set: Create or update a resource (requires super-user permissions)\n" +
		"- delete: Delete a resource (requires super-user permissions)"

	return builders.NewTool("pulsar_admin_nsisolationpolicy",
		builders.WithDescription(toolDesc),
		builders.WithString("resource", builders.Required(),
			builders.Description(resourceDesc),
		),
		builders.WithString("operation", builders.Required(),
			builders.Description(operationDesc),
		),
		builders.WithString("cluster", builders.Required(),
			builders.Description("Cluster name"),
		),
		builders.WithString("name",
			builders.Description("Name of the policy or broker to operate on, based on resource type.\n"+
				"Required for: policy.get, policy.delete, policy.set, broker.get"),
		),
		builders.WithArray("namespaces",
			builders.Description("List of namespaces to apply the isolation policy. Required for policy.set"),
			builders.Items(
				map[string]interface{}{
					"type":        "string",
					"description": "namespace",
				},
			),
		),
		builders.WithArray("primary",
			builders.Description("List of primary brokers for the namespaces. Required for policy.set"),
			builders.Items(
				map[string]interface{}{
					"type":        "string",
					"description": "primary broker",
				},
			),
		),
		builders.WithArray("secondary",
			builders.Description("List of secondary brokers for the namespaces. Optional for policy.set"),
			builders.Items(
				map[string]interface{}{
					"type":        "string",
					"description": "secondary broker",
				},
			),
		),
		builders.WithString("autoFailoverPolicyType",
			builders.Description("Auto failover policy type (e.g., min_available). Optional for policy.set"),
		),
		builders.WithObject("autoFailoverPolicyParams",
			builders.Description("Auto failover policy parameters as an object (e.g., {'min_limit': '1', 'usage_threshold': '100'}). Optional for policy.set"),
		),
	)
}

// buildNsIsolationPolicyHandler builds the Pulsar admin namespace isolation policy handler function
func (b *PulsarAdminNsIsolationPolicyToolBuilder) buildNsIsolationPolicyHandler(readOnly bool) func(context.Context, *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	return func(ctx context.Context, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return adapter.NewErrorResult("Pulsar session not found in context"), nil
		}

		client, err := session.GetAdminClient()
		if err != nil {
			return adapter.NewErrorResult("Failed to get admin client: %v", err), nil
		}

		// Get required parameters
		resource, err := adapter.RequireString(request, "resource")
		if err != nil {
			return adapter.NewErrorResult("Failed to get resource: %v", err), nil
		}

		operation, err := adapter.RequireString(request, "operation")
		if err != nil {
			return adapter.NewErrorResult("Failed to get operation: %v", err), nil
		}

		cluster, err := adapter.RequireString(request, "cluster")
		if err != nil {
			return adapter.NewErrorResult("Failed to get cluster name: %v", err), nil
		}

		// Normalize parameters
		resource = strings.ToLower(resource)
		operation = strings.ToLower(operation)

		// Validate write operations in read-only mode
		if readOnly && (operation == "set" || operation == "delete") {
			return adapter.NewErrorResult("Write operations are not allowed in read-only mode"), nil
		}

		// Dispatch based on resource type
		switch resource {
		case "policy":
			return b.handlePolicyResource(client, operation, cluster, request)
		case "broker":
			return b.handleBrokerResource(client, operation, cluster, request)
		case "brokers":
			return b.handleNsIsolationBrokersResource(client, operation, cluster, request)
		default:
			return adapter.NewErrorResult("Invalid resource: %s. Available resources: policy, broker, brokers", resource), nil
		}
	}
}

// Helper functions

// handlePolicyResource handles operations on the "policy" resource
func (b *PulsarAdminNsIsolationPolicyToolBuilder) handlePolicyResource(client cmdutils.Client, operation, cluster string, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	switch operation {
	case "get":
		name, err := adapter.RequireString(request, "name")
		if err != nil {
			return adapter.NewErrorResult("Missing required parameter 'name' for policy.get: %v", err), nil
		}

		// Get namespace isolation policy
		policyInfo, err := client.NsIsolationPolicy().GetNamespaceIsolationPolicy(cluster, name)
		if err != nil {
			return adapter.NewErrorResult("Failed to get namespace isolation policy: %v", err), nil
		}

		// Convert result to JSON string
		policyInfoJSON, err := json.Marshal(policyInfo)
		if err != nil {
			return adapter.NewErrorResult("Failed to serialize namespace isolation policy: %v", err), nil
		}

		return adapter.NewTextResult(string(policyInfoJSON)), nil

	case "list":
		// Get namespace isolation policies
		policies, err := client.NsIsolationPolicy().GetNamespaceIsolationPolicies(cluster)
		if err != nil {
			return adapter.NewErrorResult("Failed to list namespace isolation policies: %v", err), nil
		}

		// Convert result to JSON string
		policiesJSON, err := json.Marshal(policies)
		if err != nil {
			return adapter.NewErrorResult("Failed to serialize namespace isolation policies: %v", err), nil
		}

		return adapter.NewTextResult(string(policiesJSON)), nil

	case "delete":
		name, err := adapter.RequireString(request, "name")
		if err != nil {
			return adapter.NewErrorResult("Missing required parameter 'name' for policy.delete: %v", err), nil
		}

		// Delete namespace isolation policy
		err = client.NsIsolationPolicy().DeleteNamespaceIsolationPolicy(cluster, name)
		if err != nil {
			return adapter.NewErrorResult("Failed to delete namespace isolation policy: %v", err), nil
		}

		return adapter.NewTextResult(fmt.Sprintf("Delete namespace isolation policy %s successfully", name)), nil

	case "set":
		name, err := adapter.RequireString(request, "name")
		if err != nil {
			return adapter.NewErrorResult("Missing required parameter 'name' for policy.set: %v", err), nil
		}

		namespaces, err := adapter.RequireStringSlice(request, "namespaces")
		if err != nil {
			return adapter.NewErrorResult("Missing required parameter 'namespaces' for policy.set: %v", err), nil
		}

		primary, err := adapter.RequireStringSlice(request, "primary")
		if err != nil {
			return adapter.NewErrorResult("Missing required parameter 'primary' for policy.set: %v", err), nil
		}

		secondary := adapter.GetStringSlice(request, "secondary", []string{})
		autoFailoverPolicyType := adapter.GetString(request, "autoFailoverPolicyType", "")

		// Parse autoFailoverPolicyParams as a map
		autoFailoverPolicyParamsRaw := adapter.GetObject(request, "autoFailoverPolicyParams")
		if autoFailoverPolicyParamsRaw == nil {
			return adapter.NewErrorResult("Failed to extract autoFailoverPolicyParams"), nil
		}

		autoFailoverPolicyParams := make(map[string]string)
		for k, v := range autoFailoverPolicyParamsRaw {
			if strVal, ok := v.(string); ok {
				autoFailoverPolicyParams[k] = strVal
			}
		}

		// Create namespace isolation policy data
		nsIsolationData, err := utils.CreateNamespaceIsolationData(namespaces, primary, secondary,
			autoFailoverPolicyType, autoFailoverPolicyParams)
		if err != nil {
			return adapter.NewErrorResult("Failed to create namespace isolation data: %v", err), nil
		}

		// Create/update namespace isolation policy
		err = client.NsIsolationPolicy().CreateNamespaceIsolationPolicy(cluster, name, *nsIsolationData)
		if err != nil {
			return adapter.NewErrorResult("Failed to create/update namespace isolation policy: %v", err), nil
		}

		return adapter.NewTextResult(fmt.Sprintf("Create/Update namespace isolation policy %s successfully", name)), nil

	default:
		return adapter.NewErrorResult("Invalid operation for resource 'policy': %s. Available operations: get, list, delete, set", operation), nil
	}
}

// handleBrokerResource handles operations on the "broker" resource
func (b *PulsarAdminNsIsolationPolicyToolBuilder) handleBrokerResource(client cmdutils.Client, operation, cluster string, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	switch operation {
	case "get":
		name, err := adapter.RequireString(request, "name")
		if err != nil {
			return adapter.NewErrorResult("Missing required parameter 'name' for broker.get: %v", err), nil
		}

		// Get broker with policies
		brokerInfo, err := client.NsIsolationPolicy().GetBrokerWithNamespaceIsolationPolicy(cluster, name)
		if err != nil {
			return adapter.NewErrorResult("Failed to get broker with namespace isolation policy: %v", err), nil
		}

		// Convert result to JSON string
		brokerInfoJSON, err := json.Marshal(brokerInfo)
		if err != nil {
			return adapter.NewErrorResult("Failed to serialize broker information: %v", err), nil
		}

		return adapter.NewTextResult(string(brokerInfoJSON)), nil

	default:
		return adapter.NewErrorResult("Invalid operation for resource 'broker': %s. Available operations: get", operation), nil
	}
}

// handleNsIsolationBrokersResource handles operations on the "brokers" resource for namespace isolation policies
func (b *PulsarAdminNsIsolationPolicyToolBuilder) handleNsIsolationBrokersResource(client cmdutils.Client, operation, cluster string, _ *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	switch operation {
	case "list":
		// Get all brokers with policies
		brokersInfo, err := client.NsIsolationPolicy().GetBrokersWithNamespaceIsolationPolicy(cluster)
		if err != nil {
			return adapter.NewErrorResult("Failed to get brokers with namespace isolation policy: %v", err), nil
		}

		// Convert result to JSON string
		brokersInfoJSON, err := json.Marshal(brokersInfo)
		if err != nil {
			return adapter.NewErrorResult("Failed to serialize brokers information: %v", err), nil
		}

		return adapter.NewTextResult(string(brokersInfoJSON)), nil

	default:
		return adapter.NewErrorResult("Invalid operation for resource 'brokers': %s. Available operations: list", operation), nil
	}
}
