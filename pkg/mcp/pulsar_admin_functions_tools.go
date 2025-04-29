package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/admin/config"
	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
)

// PulsarAdminAddFunctionsTools adds a unified function-related tool to the MCP server
func PulsarAdminAddFunctionsTools(s *server.MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeaturePulsarAdminFunctions)) && !slices.Contains(features, string(FeatureAll)) && !slices.Contains(features, string(FeatureAllPulsar)) {
		return
	}
	// Create a single unified tool for all function operations
	toolDesc := "Manage Apache Pulsar Functions for stream processing. " +
		"Pulsar Functions are lightweight compute processes that can consume messages from one or more Pulsar topics, " +
		"apply user-defined processing logic, and produce results to another topic. " +
		"Functions support Java, Python, and Go runtimes, enabling complex event processing, " +
		"data transformations, filtering, and integration with external systems. " +
		"Functions follow the tenant/namespace/name hierarchy for organization, " +
		"can maintain state, and can scale through parallelism configuration. " +
		"This tool provides complete lifecycle management including deployment, monitoring, scaling, " +
		"state management, and triggering. Functions require proper permissions to access their topics."

	operationDesc := "Operation to perform. Available operations:\n" +
		"- list: List all functions under a specific tenant and namespace\n" +
		"- get: Get the configuration of a function\n" +
		"- status: Get the runtime status of a function (instances, metrics)\n" +
		"- stats: Get detailed statistics of a function (throughput, processing latency)\n" +
		"- querystate: Query state stored by a stateful function for a specific key\n" +
		"- create: Deploy a new function with specified parameters\n" +
		"- update: Update the configuration of an existing function\n" +
		"- delete: Delete a function\n" +
		"- start: Start a stopped function\n" +
		"- stop: Stop a running function\n" +
		"- restart: Restart a function\n" +
		"- putstate: Store state in a function's state store\n" +
		"- trigger: Manually trigger a function with a specific value"

	functionsTool := mcp.NewTool("pulsar_admin_functions",
		mcp.WithDescription(toolDesc),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc)),
		mcp.WithString("tenant", mcp.Required(),
			mcp.Description("The tenant name. Tenants are the primary organizational unit in Pulsar, "+
				"providing multi-tenancy and resource isolation. Functions deployed within a tenant "+
				"inherit its permissions and resource quotas.")),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("The namespace name. Namespaces are logical groupings of topics and functions "+
				"within a tenant. They encapsulate configuration policies and access control. "+
				"Functions in a namespace typically process topics within the same namespace.")),
		mcp.WithString("name",
			mcp.Description("The function name. Required for all operations except 'list'. "+
				"Names should be descriptive of the function's purpose and must be unique within a namespace. "+
				"Function names are used in metrics, logs, and when addressing the function via APIs.")),
		// Additional parameters for specific operations
		mcp.WithString("classname",
			mcp.Description("The fully qualified class name implementing the function. Required for 'create' operation, optional for 'update'. "+
				"For Java functions, this should be the class that implements pulsar io interfaces. "+
				"For Python, this is typically the file and function name (e.g., 'example.py:process'). "+
				"Go functions should specify the 'main' function of the binary.")),
		mcp.WithArray("inputs",
			mcp.Description("The input topics for the function (array of strings). Optional for 'create' and 'update' operations. "+
				"Topics must be specified in the format 'persistent://tenant/namespace/topic'. "+
				"Functions can consume from multiple topics, each with potentially different serialization types. "+
				"All input topics should exist before the function is created.")),
		mcp.WithString("output",
			mcp.Description("The output topic for the function results. Optional for 'create' and 'update' operations. "+
				"Specified in the format 'persistent://tenant/namespace/topic'. "+
				"If not set, the function will not produce any output to topics. "+
				"The output topic will be automatically created if it doesn't exist.")),
		mcp.WithString("jar",
			mcp.Description("Path to the JAR file containing the function code. Optional for 'create' and 'update' operations. "+
				"Can be a local path or a URL accessible to the Pulsar broker. "+
				"For Java functions, this should contain all dependencies for the function. "+
				"The jar file must be compatible with the Pulsar Functions API.")),
		mcp.WithNumber("parallelism",
			mcp.Description("The parallelism factor of the function. Optional for 'create' and 'update' operations. "+
				"Determines how many instances of the function will run concurrently. "+
				"Higher values improve throughput but require more resources. "+
				"For stateful functions, consider how parallelism affects state consistency. "+
				"Default is 1 (single instance).")),
		mcp.WithObject("userConfig",
			mcp.Description("User-defined config key/values. Optional for 'create' and 'update' operations. "+
				"Provides configuration parameters accessible to the function at runtime. "+
				"Specify as a JSON object with string, number, or boolean values. "+
				"Common configs include connection parameters, batch sizes, or feature toggles. "+
				"Example: {\"maxBatchSize\": 100, \"connectionString\": \"host:port\", \"debugMode\": true}")),
		mcp.WithString("key",
			mcp.Description("The state key. Required for 'querystate' and 'putstate' operations. "+
				"Keys are used to identify values in the function's state store. "+
				"They should be reasonable in length and follow a consistent pattern. "+
				"State keys are typically limited to 128 characters.")),
		mcp.WithString("value",
			mcp.Description("The state value. Required for 'putstate' operation. "+
				"Values are stored in the function's state system. "+
				"For simple values, specify as a string. For complex objects, use JSON-serialized strings. "+
				"State values are typically limited to 1MB in size.")),
		mcp.WithString("topic",
			mcp.Description("The specific topic name that the function should consume from. Optional for 'trigger' operation. "+
				"Specified in the format 'persistent://tenant/namespace/topic'. "+
				"Used when triggering a function that consumes from multiple topics. "+
				"If not provided, the first input topic will be used.")),
		mcp.WithString("triggerValue",
			mcp.Description("The value with which to trigger the function. Optional for 'trigger' operation. "+
				"This value will be passed to the function as if it were a message from the input topic. "+
				"String values are sent as is; for typed values, ensure proper formatting based on function expectations. "+
				"The function processes this value just like a normal message.")),
	)

	// Add the tool with the handler that captures readOnly state
	s.AddTool(functionsTool, handleFunctionsTool(readOnly))
}

// handleFunctionsTool returns a handler function for the unified functions tool
func handleFunctionsTool(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Pulsar client with API version V3
		client := cmdutils.NewPulsarClientWithAPIVersion(config.V3)

		// Extract and validate operation parameter
		operation, err := requiredParam[string](request.Params.Arguments, "operation")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'operation': %v", err)), nil
		}

		// Check if the operation is valid
		validOperations := map[string]bool{
			"list": true, "get": true, "status": true, "stats": true, "querystate": true,
			"create": true, "update": true, "delete": true, "start": true, "stop": true,
			"restart": true, "putstate": true, "trigger": true,
		}

		if !validOperations[operation] {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid operation: '%s'. Supported operations: list, get, status, stats, querystate, create, update, delete, start, stop, restart, putstate, trigger", operation)), nil
		}

		// Check write permissions for write operations
		writeOperations := map[string]bool{
			"create": true, "update": true, "delete": true, "start": true,
			"stop": true, "restart": true, "putstate": true, "trigger": true,
		}

		if readOnly && writeOperations[operation] {
			return mcp.NewToolResultError(fmt.Sprintf("Operation '%s' not allowed in read-only mode. Read-only mode restricts modifications to Pulsar Functions.", operation)), nil
		}

		// Extract common parameters
		tenant, err := requiredParam[string](request.Params.Arguments, "tenant")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'tenant': %v. A tenant is required for all Pulsar Functions operations.", err)), nil
		}

		namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'namespace': %v. A namespace is required for all Pulsar Functions operations.", err)), nil
		}

		// For all operations except 'list', name is required
		var name string
		if operation != "list" {
			name, err = requiredParam[string](request.Params.Arguments, "name")
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'name' for operation '%s': %v. The function name must be specified for this operation.", operation, err)), nil
			}
		}

		// Handle operation
		switch operation {
		case "list":
			return handleFunctionList(ctx, client, tenant, namespace)
		case "get":
			return handleFunctionGet(ctx, client, tenant, namespace, name)
		case "status":
			return handleFunctionStatus(ctx, client, tenant, namespace, name)
		case "stats":
			return handleFunctionStats(ctx, client, tenant, namespace, name)
		case "querystate":
			key, err := requiredParam[string](request.Params.Arguments, "key")
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'key' for operation 'querystate': %v. A key is required to look up state in the function's state store.", err)), nil
			}
			return handleFunctionQuerystate(ctx, client, tenant, namespace, name, key)
		case "create":
			return handleFunctionCreate(ctx, client, tenant, namespace, name, request.Params.Arguments)
		case "update":
			return handleFunctionUpdate(ctx, client, tenant, namespace, name, request.Params.Arguments)
		case "delete":
			return handleFunctionDelete(ctx, client, tenant, namespace, name)
		case "start":
			return handleFunctionStart(ctx, client, tenant, namespace, name)
		case "stop":
			return handleFunctionStop(ctx, client, tenant, namespace, name)
		case "restart":
			return handleFunctionRestart(ctx, client, tenant, namespace, name)
		case "putstate":
			key, err := requiredParam[string](request.Params.Arguments, "key")
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'key' for operation 'putstate': %v. A key is required to store state in the function's state store.", err)), nil
			}
			value, err := requiredParam[string](request.Params.Arguments, "value")
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'value' for operation 'putstate': %v. A value is required to store state in the function's state store.", err)), nil
			}
			return handleFunctionPutstate(ctx, client, tenant, namespace, name, key, value)
		case "trigger":
			topic, _ := optionalParam[string](request.Params.Arguments, "topic")
			triggerValue, _ := optionalParam[string](request.Params.Arguments, "triggerValue")
			return handleFunctionTrigger(ctx, client, tenant, namespace, name, topic, triggerValue)
		default:
			// This should never happen due to the valid operations check above
			return mcp.NewToolResultError(fmt.Sprintf("Unsupported operation: %s", operation)), nil
		}
	}
}

// handleFunctionList handles listing all functions under a namespace
func handleFunctionList(ctx context.Context, client cmdutils.Client, tenant, namespace string) (*mcp.CallToolResult, error) {
	// Get function list
	functions, err := client.Functions().GetFunctions(tenant, namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list functions in tenant '%s' namespace '%s': %v. Check that the tenant and namespace exist and you have proper permissions.",
			tenant, namespace, err)), nil
	}

	// Convert result to JSON string
	functionsJSON, err := json.Marshal(functions)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize function list: %v", err)), nil
	}

	return mcp.NewToolResultText(string(functionsJSON)), nil
}

// handleFunctionGet handles getting information about a function
func handleFunctionGet(ctx context.Context, client cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	// Get function info
	functionConfig, err := client.Functions().GetFunction(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get function '%s' in tenant '%s' namespace '%s': %v. Verify the function exists and you have proper permissions.",
			name, tenant, namespace, err)), nil
	}

	// Convert result to JSON string
	functionJSON, err := json.Marshal(functionConfig)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize function info: %v", err)), nil
	}

	return mcp.NewToolResultText(string(functionJSON)), nil
}

// handleFunctionStatus handles getting the status of a function
func handleFunctionStatus(ctx context.Context, client cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	// Get function status
	status, err := client.Functions().GetFunctionStatus(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get status for function '%s' in tenant '%s' namespace '%s': %v. Verify the function exists and is properly deployed.",
			name, tenant, namespace, err)), nil
	}

	// Convert result to JSON string
	statusJSON, err := json.Marshal(status)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize function status: %v", err)), nil
	}

	return mcp.NewToolResultText(string(statusJSON)), nil
}

// handleFunctionStats handles getting the statistics of a function
func handleFunctionStats(ctx context.Context, client cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	// Get function stats
	stats, err := client.Functions().GetFunctionStats(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get stats for function '%s' in tenant '%s' namespace '%s': %v. The function must be running to retrieve statistics.",
			name, tenant, namespace, err)), nil
	}

	// Convert result to JSON string
	statsJSON, err := json.Marshal(stats)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize function stats: %v", err)), nil
	}

	return mcp.NewToolResultText(string(statsJSON)), nil
}

// handleFunctionQuerystate handles querying the state of a function
func handleFunctionQuerystate(ctx context.Context, client cmdutils.Client, tenant, namespace, name, key string) (*mcp.CallToolResult, error) {
	// Query function state
	state, err := client.Functions().GetFunctionState(tenant, namespace, name, key)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to query state for function '%s' in tenant '%s' namespace '%s' with key '%s': %v. Verify the function exists and has state storage enabled.",
			name, tenant, namespace, key, err)), nil
	}

	// Convert to JSON for display
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize state: %v", err)), nil
	}

	return mcp.NewToolResultText(string(stateJSON)), nil
}

// handleFunctionCreate handles creating a new function
func handleFunctionCreate(ctx context.Context, client cmdutils.Client, tenant, namespace, name string, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	// Extract required parameters
	classname, err := requiredParam[string](arguments, "classname")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get classname (required for function creation): %v. The classname must specify the fully qualified class implementing the function.", err)), nil
	}

	// Extract optional parameters
	inputTopics, _ := optionalParamArray[string](arguments, "inputs")
	output, _ := optionalParam[string](arguments, "output")
	jar, _ := optionalParam[string](arguments, "jar")
	parallelismFloat, _ := optionalParam[float64](arguments, "parallelism")
	parallelism := int(parallelismFloat)

	// Get user config if available
	var userConfigMap map[string]interface{}
	userConfigObj, ok := arguments["userConfig"]
	if ok && userConfigObj != nil {
		if configMap, isMap := userConfigObj.(map[string]interface{}); isMap {
			userConfigMap = configMap
		}
	}

	// Create function config
	functionConfig := &utils.FunctionConfig{
		Tenant:      tenant,
		Namespace:   namespace,
		Name:        name,
		ClassName:   classname,
		Output:      output,
		Parallelism: parallelism,
	}

	// Set jar if provided
	if jar != "" {
		jarPtr := jar
		functionConfig.Jar = &jarPtr
	}

	// Set inputs
	if len(inputTopics) > 0 {
		functionConfig.Inputs = inputTopics
	} else {
		return mcp.NewToolResultError("At least one input topic should be specified for function creation. Use the 'inputs' parameter with an array of topic names."), nil
	}

	// Set user config if available
	if userConfigMap != nil {
		functionConfig.UserConfig = userConfigMap
	}

	// Create function
	err = client.Functions().CreateFunc(functionConfig, "")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create function '%s' in tenant '%s' namespace '%s': %v. Verify all parameters are correct and required resources exist.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Created function '%s' successfully in tenant '%s' namespace '%s'. The function will start processing messages from input topics.",
		name, tenant, namespace)), nil
}

// handleFunctionUpdate handles updating an existing function
func handleFunctionUpdate(ctx context.Context, client cmdutils.Client, tenant, namespace, name string, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	// Extract optional parameters
	classname, _ := optionalParam[string](arguments, "classname")
	inputTopics, _ := optionalParamArray[string](arguments, "inputs")
	output, _ := optionalParam[string](arguments, "output")
	jar, _ := optionalParam[string](arguments, "jar")
	parallelismFloat, _ := optionalParam[float64](arguments, "parallelism")
	parallelism := int(parallelismFloat)

	// Get user config if available
	var userConfigMap map[string]interface{}
	userConfigObj, ok := arguments["userConfig"]
	if ok && userConfigObj != nil {
		if configMap, isMap := userConfigObj.(map[string]interface{}); isMap {
			userConfigMap = configMap
		}
	}

	// Create function config
	functionConfig := &utils.FunctionConfig{
		Tenant:      tenant,
		Namespace:   namespace,
		Name:        name,
		ClassName:   classname,
		Output:      output,
		Parallelism: parallelism,
	}

	// Set jar if provided
	if jar != "" {
		jarPtr := jar
		functionConfig.Jar = &jarPtr
	}

	// Set inputs
	if len(inputTopics) > 0 {
		functionConfig.Inputs = inputTopics
	}

	// Set user config if available
	if userConfigMap != nil {
		functionConfig.UserConfig = userConfigMap
	}

	// Create update options
	updateOptions := &utils.UpdateOptions{
		UpdateAuthData: true,
	}

	// Update function
	err := client.Functions().UpdateFunction(functionConfig, "", updateOptions)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to update function '%s' in tenant '%s' namespace '%s': %v. Verify the function exists and all parameters are valid.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Updated function '%s' successfully in tenant '%s' namespace '%s'. The function may need to be restarted to apply all changes.",
		name, tenant, namespace)), nil
}

// handleFunctionDelete handles deleting a function
func handleFunctionDelete(ctx context.Context, client cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	// Delete function
	err := client.Functions().DeleteFunction(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete function '%s' in tenant '%s' namespace '%s': %v. Verify the function exists and you have deletion permissions.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Deleted function '%s' successfully from tenant '%s' namespace '%s'. All running instances have been terminated.",
		name, tenant, namespace)), nil
}

// handleFunctionStart handles starting a function
func handleFunctionStart(ctx context.Context, client cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	// Start function
	err := client.Functions().StartFunction(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to start function '%s' in tenant '%s' namespace '%s': %v. Verify the function exists and is not already running.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Started function '%s' successfully in tenant '%s' namespace '%s'. The function will begin processing messages from its input topics.",
		name, tenant, namespace)), nil
}

// handleFunctionStop handles stopping a function
func handleFunctionStop(ctx context.Context, client cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	// Stop function
	err := client.Functions().StopFunction(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to stop function '%s' in tenant '%s' namespace '%s': %v. Verify the function exists and is currently running.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Stopped function '%s' successfully in tenant '%s' namespace '%s'. The function will no longer process messages until restarted.",
		name, tenant, namespace)), nil
}

// handleFunctionRestart handles restarting a function
func handleFunctionRestart(ctx context.Context, client cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	// Restart function
	err := client.Functions().RestartFunction(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to restart function '%s' in tenant '%s' namespace '%s': %v. Verify the function exists and is properly deployed.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Restarted function '%s' successfully in tenant '%s' namespace '%s'. All function instances have been restarted.",
		name, tenant, namespace)), nil
}

// handleFunctionPutstate handles putting state into the state store of a function
func handleFunctionPutstate(ctx context.Context, client cmdutils.Client, tenant, namespace, name, key, value string) (*mcp.CallToolResult, error) {
	// Create function state
	state := utils.FunctionState{
		Key:         key,
		StringValue: value,
	}

	// Put function state
	err := client.Functions().PutFunctionState(tenant, namespace, name, state)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to put state for function '%s' in tenant '%s' namespace '%s' with key '%s': %v. Verify the function exists and has state storage enabled.",
			name, tenant, namespace, key, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Put state for function '%s' with key '%s' successfully in tenant '%s' namespace '%s'. The state is now stored and available to the function instances.",
		name, key, tenant, namespace)), nil
}

// handleFunctionTrigger handles triggering a function
func handleFunctionTrigger(ctx context.Context, client cmdutils.Client, tenant, namespace, name, topic, triggerValue string) (*mcp.CallToolResult, error) {
	// Trigger function
	result, err := client.Functions().TriggerFunction(tenant, namespace, name, topic, triggerValue, "")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to trigger function '%s' in tenant '%s' namespace '%s': %v. Verify the function exists and is running.",
			name, tenant, namespace, err)), nil
	}

	if result == "" {
		return mcp.NewToolResultText(fmt.Sprintf("Triggered function '%s' successfully in tenant '%s' namespace '%s', but it did not return any output. The function was executed with the provided value.",
			name, tenant, namespace)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Triggered function '%s' successfully in tenant '%s' namespace '%s'. Function output: %s",
		name, tenant, namespace, result)), nil
}
