package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// PulsarAdminAddNsIsolationPolicyTools adds namespace isolation policy related tools to the MCP server
func PulsarAdminAddNsIsolationPolicyTools(s *server.MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeaturePulsarAdminNsIsolationPolicy)) && !slices.Contains(features, string(FeatureAll)) && !slices.Contains(features, string(FeatureAllPulsar)) {
		return
	}

	// Add unified namespace isolation policy tool
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

	nsIsolationPolicyTool := mcp.NewTool("pulsar_admin_nsisolationpolicy",
		mcp.WithDescription(toolDesc),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description(resourceDesc),
		),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
		),
		mcp.WithString("cluster", mcp.Required(),
			mcp.Description("Cluster name"),
		),
		mcp.WithString("name",
			mcp.Description("Name of the policy or broker to operate on, based on resource type.\n"+
				"Required for: policy.get, policy.delete, policy.set, broker.get"),
		),
		mcp.WithArray("namespaces",
			mcp.Description("List of namespaces to apply the isolation policy. Required for policy.set"),
		),
		mcp.WithArray("primary",
			mcp.Description("List of primary brokers for the namespaces. Required for policy.set"),
		),
		mcp.WithArray("secondary",
			mcp.Description("List of secondary brokers for the namespaces. Optional for policy.set"),
		),
		mcp.WithString("autoFailoverPolicyType",
			mcp.Description("Auto failover policy type (e.g., min_available). Optional for policy.set"),
		),
		mcp.WithObject("autoFailoverPolicyParams",
			mcp.Description("Auto failover policy parameters as an object (e.g., {'min_limit': '1', 'usage_threshold': '100'}). Optional for policy.set"),
		),
	)
	s.AddTool(nsIsolationPolicyTool, handleNsIsolationPolicy(readOnly))
}

// handleNsIsolationPolicy returns a function that handles namespace isolation policy operations
func handleNsIsolationPolicy(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Pulsar client
		client, err := pulsar.GetAdminClient()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
		}

		// Get required parameters
		resource, err := requiredParam[string](request.Params.Arguments, "resource")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get resource: %v", err)), nil
		}

		operation, err := requiredParam[string](request.Params.Arguments, "operation")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get operation: %v", err)), nil
		}

		cluster, err := requiredParam[string](request.Params.Arguments, "cluster")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get cluster name: %v", err)), nil
		}

		// Normalize parameters
		resource = strings.ToLower(resource)
		operation = strings.ToLower(operation)

		// Validate write operations in read-only mode
		if readOnly && (operation == "set" || operation == "delete") {
			return mcp.NewToolResultError("Write operations are not allowed in read-only mode"), nil
		}

		// Dispatch based on resource type
		switch resource {
		case "policy":
			return handlePolicyResource(client, operation, cluster, request)
		case "broker":
			return handleBrokerResource(client, operation, cluster, request)
		case "brokers":
			return handleNsIsolationBrokersResource(client, operation, cluster, request)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resource: %s. Available resources: policy, broker, brokers", resource)), nil
		}
	}
}

// handlePolicyResource handles operations on the "policy" resource
func handlePolicyResource(client cmdutils.Client, operation, cluster string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	switch operation {
	case "get":
		name, err := requiredParam[string](request.Params.Arguments, "name")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'name' for policy.get: %v", err)), nil
		}

		// Get namespace isolation policy
		policyInfo, err := client.NsIsolationPolicy().GetNamespaceIsolationPolicy(cluster, name)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace isolation policy: %v", err)), nil
		}

		// Convert result to JSON string
		policyInfoJSON, err := json.Marshal(policyInfo)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize namespace isolation policy: %v", err)), nil
		}

		return mcp.NewToolResultText(string(policyInfoJSON)), nil

	case "list":
		// Get namespace isolation policies
		policies, err := client.NsIsolationPolicy().GetNamespaceIsolationPolicies(cluster)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list namespace isolation policies: %v", err)), nil
		}

		// Convert result to JSON string
		policiesJSON, err := json.Marshal(policies)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize namespace isolation policies: %v", err)), nil
		}

		return mcp.NewToolResultText(string(policiesJSON)), nil

	case "delete":
		name, err := requiredParam[string](request.Params.Arguments, "name")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'name' for policy.delete: %v", err)), nil
		}

		// Delete namespace isolation policy
		err = client.NsIsolationPolicy().DeleteNamespaceIsolationPolicy(cluster, name)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete namespace isolation policy: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Delete namespace isolation policy %s successfully", name)), nil

	case "set":
		name, err := requiredParam[string](request.Params.Arguments, "name")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'name' for policy.set: %v", err)), nil
		}

		namespaces, err := getRequiredParamArray[string](request.Params.Arguments, "namespaces")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'namespaces' for policy.set: %v", err)), nil
		}

		primary, err := getRequiredParamArray[string](request.Params.Arguments, "primary")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'primary' for policy.set: %v", err)), nil
		}

		secondary, _ := optionalParamArray[string](request.Params.Arguments, "secondary")
		autoFailoverPolicyType, _ := optionalParam[string](request.Params.Arguments, "autoFailoverPolicyType")

		// Parse autoFailoverPolicyParams as a map
		autoFailoverPolicyParamsRaw, _ := optionalParam[map[string]interface{}](request.Params.Arguments, "autoFailoverPolicyParams")
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
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create namespace isolation data: %v", err)), nil
		}

		// Create/update namespace isolation policy
		err = client.NsIsolationPolicy().CreateNamespaceIsolationPolicy(cluster, name, *nsIsolationData)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create/update namespace isolation policy: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Create/Update namespace isolation policy %s successfully", name)), nil

	default:
		return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'policy': %s. Available operations: get, list, delete, set", operation)), nil
	}
}

// handleBrokerResource handles operations on the "broker" resource
func handleBrokerResource(client cmdutils.Client, operation, cluster string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	switch operation {
	case "get":
		name, err := requiredParam[string](request.Params.Arguments, "name")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'name' for broker.get: %v", err)), nil
		}

		// Get broker with policies
		brokerInfo, err := client.NsIsolationPolicy().GetBrokerWithNamespaceIsolationPolicy(cluster, name)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get broker with namespace isolation policy: %v", err)), nil
		}

		// Convert result to JSON string
		brokerInfoJSON, err := json.Marshal(brokerInfo)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize broker information: %v", err)), nil
		}

		return mcp.NewToolResultText(string(brokerInfoJSON)), nil

	default:
		return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'broker': %s. Available operations: get", operation)), nil
	}
}

// handleNsIsolationBrokersResource handles operations on the "brokers" resource for namespace isolation policies
func handleNsIsolationBrokersResource(client cmdutils.Client, operation, cluster string, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	switch operation {
	case "list":
		// Get all brokers with policies
		brokersInfo, err := client.NsIsolationPolicy().GetBrokersWithNamespaceIsolationPolicy(cluster)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get brokers with namespace isolation policy: %v", err)), nil
		}

		// Convert result to JSON string
		brokersInfoJSON, err := json.Marshal(brokersInfo)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize brokers information: %v", err)), nil
		}

		return mcp.NewToolResultText(string(brokersInfoJSON)), nil

	default:
		return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'brokers': %s. Available operations: list", operation)), nil
	}
}

// getRequiredParamArray is a helper function to get a required array parameter with a specific type
func getRequiredParamArray[T any](args map[string]interface{}, paramName string) ([]T, error) {
	value, ok := args[paramName]
	if !ok {
		return nil, fmt.Errorf("missing required parameter: %s", paramName)
	}

	arrayValue, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("parameter %s is not an array", paramName)
	}

	result := make([]T, 0, len(arrayValue))
	for _, v := range arrayValue {
		typedValue, ok := v.(T)
		if !ok {
			return nil, fmt.Errorf("array element in %s has incorrect type", paramName)
		}
		result = append(result, typedValue)
	}

	return result, nil
}
