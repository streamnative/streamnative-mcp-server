package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// PulsarAdminAddBrokersTools adds broker-related tools to the MCP server
func PulsarAdminAddBrokersTools(s *server.MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeaturePulsarAdminBrokers)) && !slices.Contains(features, string(FeatureAll)) && !slices.Contains(features, string(FeatureAllPulsar)) {
		return
	}
	// Create a single unified broker tool to replace multiple individual tools
	brokersTool := mcp.NewTool("pulsar_admin_brokers",
		mcp.WithDescription("Unified tool for managing Apache Pulsar broker resources. This tool integrates multiple broker management functions, including:\n"+
			"1. List active brokers in a cluster (resource=brokers, operation=list)\n"+
			"2. Check broker health status (resource=health, operation=get)\n"+
			"3. Manage broker configurations (resource=config, operation=get/update/delete)\n"+
			"4. View namespaces owned by a broker (resource=namespaces, operation=get)\n\n"+
			"Different functions are accessed by combining resource and operation parameters, with other parameters used selectively based on operation type.\n"+
			"Example: {\"resource\": \"config\", \"operation\": \"get\", \"configType\": \"dynamic\"} retrieves all dynamic configuration names.\n"+
			"This tool requires Pulsar super-user permissions."),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description("Type of resource to access, available options:\n"+
				"- brokers: Manage broker listings\n"+
				"- health: Check broker health status\n"+
				"- config: Manage broker configurations\n"+
				"- namespaces: Manage namespaces owned by a broker"),
		),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description("Operation to perform, available options:\n"+
				"- list: List resources (used with brokers)\n"+
				"- get: Retrieve resource information (used with health, config, namespaces)\n"+
				"- update: Update a resource (used with config)\n"+
				"- delete: Delete a resource (used with config)"),
		),
		mcp.WithString("clusterName",
			mcp.Description("Pulsar cluster name, required for these operations:\n"+
				"- When resource=brokers, operation=list\n"+
				"- When resource=namespaces, operation=get"),
		),
		mcp.WithString("brokerUrl",
			mcp.Description("Broker URL, such as '127.0.0.1:8080', required for these operations:\n"+
				"- When resource=namespaces, operation=get"),
		),
		mcp.WithString("configType",
			mcp.Description("Configuration type, required when resource=config, operation=get, available options:\n"+
				"- dynamic: Get list of dynamically modifiable configuration names\n"+
				"- runtime: Get all runtime configurations (including static and dynamic configs)\n"+
				"- internal: Get internal configuration information\n"+
				"- all_dynamic: Get all dynamic configurations and their current values"),
		),
		mcp.WithString("configName",
			mcp.Description("Configuration parameter name, required for these operations:\n"+
				"- When resource=config, operation=update\n"+
				"- When resource=config, operation=delete"),
		),
		mcp.WithString("configValue",
			mcp.Description("Configuration parameter value, required for these operations:\n"+
				"- When resource=config, operation=update"),
		),
	)
	s.AddTool(brokersTool, handleBrokerTool(readOnly))
}

// Return a handler function, with readOnly state captured in closure
func handleBrokerTool(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := pulsar.GetAdminClient()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
		}

		// Get required parameters
		resource, err := requiredParam[string](request.Params.Arguments, "resource")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required resource parameter. " +
				"Please specify one of: brokers, health, config, namespaces.")), nil
		}

		operation, err := requiredParam[string](request.Params.Arguments, "operation")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required operation parameter. " +
				"Please specify one of: list, get, update, delete based on the resource type.")), nil
		}

		// Validate if the parameter combination is valid
		validCombination, errMsg := validateResourceOperation(resource, operation)
		if !validCombination {
			return mcp.NewToolResultError(errMsg), nil
		}

		// Process request based on resource type
		switch resource {
		case "brokers":
			return handleBrokersResource(client, operation, request)
		case "health":
			return handleHealthResource(client, operation, request)
		case "config":
			// Check write operation permissions
			if (operation == "update" || operation == "delete") && readOnly {
				return mcp.NewToolResultError("Configuration update/delete operations not allowed in read-only mode. " +
					"Please contact your administrator if you need to modify broker configurations."), nil
			}
			return handleConfigResource(client, operation, request)
		case "namespaces":
			return handleNamespacesResource(client, operation, request)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Unsupported resource: %s. "+
				"Please use one of: brokers, health, config, namespaces.", resource)), nil
		}
	}
}

// Validate if the resource and operation combination is valid
func validateResourceOperation(resource, operation string) (bool, string) {
	validCombinations := map[string][]string{
		"brokers":    {"list"},
		"health":     {"get"},
		"config":     {"get", "update", "delete"},
		"namespaces": {"get"},
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

	return false, fmt.Sprintf("Invalid resource '%s'. Valid resources are: brokers, health, config, namespaces", resource)
}

// Handle brokers resource
func handleBrokersResource(client cmdutils.Client, operation string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	switch operation {
	case "list":
		clusterName, err := requiredParam[string](request.Params.Arguments, "clusterName")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'clusterName'. " +
				"Please provide the name of the Pulsar cluster to list brokers for.")), nil
		}

		brokers, err := client.Brokers().GetActiveBrokers(clusterName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get active brokers: %v. "+
				"Please verify the cluster name and ensure the Pulsar service is running.", err)), nil
		}

		brokersJSON, err := json.Marshal(brokers)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize brokers list: %v", err)), nil
		}

		return mcp.NewToolResultText(string(brokersJSON)), nil
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unsupported operation '%s' for brokers resource. "+
			"The only supported operation is 'list'.", operation)), nil
	}
}

// Handle health resource
func handleHealthResource(client cmdutils.Client, operation string, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	switch operation {
	case "get":
		//nolint:staticcheck
		err := client.Brokers().HealthCheck()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Broker health check failed: %v. "+
				"The broker might be down or experiencing issues.", err)), nil
		}
		return mcp.NewToolResultText("ok"), nil
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unsupported operation '%s' for health resource. "+
			"The only supported operation is 'get'.", operation)), nil
	}
}

// Handle config resource
func handleConfigResource(client cmdutils.Client, operation string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	switch operation {
	case "get":
		configType, err := requiredParam[string](request.Params.Arguments, "configType")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'configType'. " +
				"Please specify one of: dynamic, runtime, internal, all_dynamic.")), nil
		}

		var result interface{}
		var fetchErr error

		switch configType {
		case "dynamic":
			result, fetchErr = client.Brokers().GetDynamicConfigurationNames()
		case "runtime":
			result, fetchErr = client.Brokers().GetRuntimeConfigurations()
		case "internal":
			result, fetchErr = client.Brokers().GetInternalConfigurationData()
		case "all_dynamic":
			result, fetchErr = client.Brokers().GetAllDynamicConfigurations()
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid config type: '%s'. "+
				"Valid types are: dynamic, runtime, internal, all_dynamic.", configType)), nil
		}

		if fetchErr != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get %s configuration: %v", configType, fetchErr)), nil
		}

		resultJSON, err := json.Marshal(result)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize configuration: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil

	case "update":
		configName, err := requiredParam[string](request.Params.Arguments, "configName")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'configName'. " +
				"Please provide the name of the configuration parameter to update.")), nil
		}

		configValue, err := requiredParam[string](request.Params.Arguments, "configValue")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'configValue'. " +
				"Please provide the new value for the configuration parameter.")), nil
		}

		err = client.Brokers().UpdateDynamicConfiguration(configName, configValue)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update configuration: %v. "+
				"Please verify the configuration name is valid and the value is of the correct type.", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Dynamic configuration '%s' updated successfully to '%s'",
			configName, configValue)), nil

	case "delete":
		configName, err := requiredParam[string](request.Params.Arguments, "configName")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'configName'. " +
				"Please provide the name of the configuration parameter to delete.")), nil
		}

		err = client.Brokers().DeleteDynamicConfiguration(configName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete configuration: %v. "+
				"Please verify the configuration name is valid and exists.", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Dynamic configuration '%s' deleted successfully", configName)), nil

	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unsupported operation '%s' for config resource. "+
			"Supported operations are: get, update, delete.", operation)), nil
	}
}

// Handle namespaces resource
func handleNamespacesResource(client cmdutils.Client, operation string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	switch operation {
	case "get":
		clusterName, err := requiredParam[string](request.Params.Arguments, "clusterName")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'clusterName'. " +
				"Please provide the name of the Pulsar cluster.")), nil
		}

		brokerURL, err := requiredParam[string](request.Params.Arguments, "brokerUrl")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'brokerUrl'. " +
				"Please provide the URL of the broker (e.g., '127.0.0.1:8080').")), nil
		}

		namespaces, err := client.Brokers().GetOwnedNamespaces(clusterName, brokerURL)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get owned namespaces: %v. "+
				"Please verify the cluster name and broker URL are correct.", err)), nil
		}

		namespacesJSON, err := json.Marshal(namespaces)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize namespaces: %v", err)), nil
		}

		return mcp.NewToolResultText(string(namespacesJSON)), nil

	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unsupported operation '%s' for namespaces resource. "+
			"The only supported operation is 'get'.", operation)), nil
	}
}
