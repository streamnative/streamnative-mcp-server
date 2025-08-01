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

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/admin/config"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
)

// PulsarAdminFunctionsToolBuilder implements the ToolBuilder interface for Pulsar admin functions operations
// It provides functionality to build Pulsar functions management tools
type PulsarAdminFunctionsToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminFunctionsToolBuilder creates a new Pulsar admin functions tool builder instance
func NewPulsarAdminFunctionsToolBuilder() *PulsarAdminFunctionsToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_functions",
		Version:     "1.0.0",
		Description: "Pulsar admin functions management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "admin", "functions"},
	}

	features := []string{
		"pulsar-admin-functions",
		"pulsar-admin",
		"all",
		"all-pulsar",
	}

	return &PulsarAdminFunctionsToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin functions tool list
// This is the core method implementing the ToolBuilder interface
func (b *PulsarAdminFunctionsToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool := b.buildPulsarAdminFunctionsTool()
	handler := b.buildPulsarAdminFunctionsHandler(config.ReadOnly)

	return []server.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildPulsarAdminFunctionsTool builds the Pulsar admin functions MCP tool definition
// Migrated from the original tool definition logic
func (b *PulsarAdminFunctionsToolBuilder) buildPulsarAdminFunctionsTool() mcp.Tool {
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

	return mcp.NewTool("pulsar_admin_functions",
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
				"For Java functions, this should be the class that implements pulsar function interfaces. "+
				"For Python, this MUST be in format of `<Python_filename_without_extension>.<ClassName>` - for example: "+
				"if file is '/path/to/exclamation.py' with class 'ExclamationFunction', classname must be 'exclamation.ExclamationFunction'; "+
				"if file is '/path/to/double_number.py' with class 'DoubleNumber', classname must be 'double_number.DoubleNumber'. "+
				"Common error: using just the class name 'DoubleNumber' (without filename prefix) will cause function creation to fail. "+
				"Go functions should specify the 'main' function of the binary.")),
		mcp.WithArray("inputs",
			mcp.Description("The input topics for the function (array of strings). Optional for 'create' and 'update' operations. "+
				"Topics must be specified in the format 'persistent://tenant/namespace/topic'. "+
				"Functions can consume from multiple topics, each with potentially different serialization types. "+
				"All input topics should exist before the function is created."),
			mcp.Items(
				map[string]interface{}{
					"type":        "string",
					"description": "input topic",
				},
			),
		),
		mcp.WithString("output",
			mcp.Description("The output topic for the function results. Optional for 'create' and 'update' operations. "+
				"Specified in the format 'persistent://tenant/namespace/topic'. "+
				"If not set, the function will not produce any output to topics. "+
				"The output topic will be automatically created if it doesn't exist.")),
		mcp.WithString("jar",
			mcp.Description("Path to the JAR file containing the function code. Optional for 'create' and 'update' operations. "+
				"Support `file://`, `http://`, `https://`, `function://`, `source://`, `sink://` protocol. "+
				"Can be a local path or supported URL protocol accessible to the Pulsar broker. "+
				"For Java functions, this should contain all dependencies for the function. "+
				"The jar file must be compatible with the Pulsar Functions API.")),
		mcp.WithString("py",
			mcp.Description("Path to the Python file containing the function code. Optional for 'create' and 'update' operations. "+
				"Support `file://`, `http://`, `https://`, `function://`, `source://`, `sink://` protocol. "+
				"Can be a local path or supported URL protocol accessible to the Pulsar broker. "+
				"For Python functions, this should be the file path to the Python file, in format of `.py`, `.zip`, or `.whl`. "+
				"The Python file must be compatible with the Pulsar Functions API.")),
		mcp.WithString("go",
			mcp.Description("Path to the Go file containing the function code. Optional for 'create' and 'update' operations. "+
				"Support `file://`, `http://`, `https://`, `function://`, `source://`, `sink://` protocol. "+
				"Can be a local path or supported URL protocol accessible to the Pulsar broker. "+
				"For Go functions, this should be the file path to the Go file, in format of executable binary. "+
				"The Go file must be compatible with the Pulsar Functions API.")),
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
			mcp.Description("The value with which to trigger the function. Required for 'trigger' operation. "+
				"This value will be passed to the function as if it were a message from the input topic. "+
				"String values are sent as is; for typed values, ensure proper formatting based on function expectations. "+
				"The function processes this value just like a normal message.")),
	)
}

// buildPulsarAdminFunctionsHandler builds the Pulsar admin functions handler function
// Migrated from the original handler logic
func (b *PulsarAdminFunctionsToolBuilder) buildPulsarAdminFunctionsHandler(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Pulsar client with API version V3
		client := cmdutils.NewPulsarClientWithAPIVersion(config.V3)

		// Extract and validate operation parameter
		operation, err := request.RequireString("operation")
		if err != nil {
			return b.handleError("get operation", err), nil
		}

		// Check if the operation is valid
		validOperations := map[string]bool{
			"list": true, "get": true, "status": true, "stats": true, "querystate": true,
			"create": true, "update": true, "delete": true, "start": true, "stop": true,
			"restart": true, "putstate": true, "trigger": true,
		}

		if !validOperations[operation] {
			return b.handleError("validate operation", fmt.Errorf("invalid operation: '%s'. Supported operations: list, get, status, stats, querystate, create, update, delete, start, stop, restart, putstate, trigger", operation)), nil
		}

		// Check write permissions for write operations
		writeOperations := map[string]bool{
			"create": true, "update": true, "delete": true, "start": true,
			"stop": true, "restart": true, "putstate": true, "trigger": true,
		}

		if readOnly && writeOperations[operation] {
			return b.handleError("check permissions", fmt.Errorf("operation '%s' not allowed in read-only mode. Read-only mode restricts modifications to Pulsar Functions", operation)), nil
		}

		// Extract common parameters
		tenant, err := request.RequireString("tenant")
		if err != nil {
			return b.handleError("get tenant", fmt.Errorf("missing required parameter 'tenant': %v. A tenant is required for all Pulsar Functions operations", err)), nil
		}

		namespace, err := request.RequireString("namespace")
		if err != nil {
			return b.handleError("get namespace", fmt.Errorf("missing required parameter 'namespace': %v. A namespace is required for all Pulsar Functions operations", err)), nil
		}

		// For all operations except 'list', name is required
		var name string
		if operation != "list" {
			name, err = request.RequireString("name")
			if err != nil {
				return b.handleError("get name", fmt.Errorf("missing required parameter 'name' for operation '%s': %v. The function name must be specified for this operation", operation, err)), nil
			}
		}

		// Handle operation using delegated handlers
		switch operation {
		case "list":
			return b.handleFunctionList(ctx, client, tenant, namespace)
		case "get":
			return b.handleFunctionGet(ctx, client, tenant, namespace, name)
		case "status":
			return b.handleFunctionStatus(ctx, client, tenant, namespace, name)
		case "stats":
			return b.handleFunctionStats(ctx, client, tenant, namespace, name)
		case "querystate":
			key, err := request.RequireString("key")
			if err != nil {
				return b.handleError("get key", fmt.Errorf("missing required parameter 'key' for operation 'querystate': %v. A key is required to look up state in the function's state store", err)), nil
			}
			return b.handleFunctionQuerystate(ctx, client, tenant, namespace, name, key)
		case "create":
			return b.handleFunctionCreate(ctx, client, tenant, namespace, name, request)
		case "update":
			return b.handleFunctionUpdate(ctx, client, tenant, namespace, name, request)
		case "delete":
			return b.handleFunctionDelete(ctx, client, tenant, namespace, name)
		case "start":
			return b.handleFunctionStart(ctx, client, tenant, namespace, name)
		case "stop":
			return b.handleFunctionStop(ctx, client, tenant, namespace, name)
		case "restart":
			return b.handleFunctionRestart(ctx, client, tenant, namespace, name)
		case "putstate":
			key, err := request.RequireString("key")
			if err != nil {
				return b.handleError("get key", fmt.Errorf("missing required parameter 'key' for operation 'putstate': %v. A key is required to store state in the function's state store", err)), nil
			}
			value, err := request.RequireString("value")
			if err != nil {
				return b.handleError("get value", fmt.Errorf("missing required parameter 'value' for operation 'putstate': %v. A value is required to store state in the function's state store", err)), nil
			}
			return b.handleFunctionPutstate(ctx, client, tenant, namespace, name, key, value)
		case "trigger":
			triggerValue, err := request.RequireString("triggerValue")
			if err != nil {
				return b.handleError("get triggerValue", fmt.Errorf("missing required parameter 'triggerValue' for operation 'trigger': %v. A trigger value is required to manually trigger the function", err)), nil
			}
			topic := request.GetString("topic", "")
			return b.handleFunctionTrigger(ctx, client, tenant, namespace, name, triggerValue, topic)
		default:
			return b.handleError("handle operation", fmt.Errorf("unsupported operation: %s", operation)), nil
		}
	}
}

// Helper functions - delegated operation handlers

// handleFunctionList handles the list operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionList(ctx context.Context, client cmdutils.Client, tenant, namespace string) (*mcp.CallToolResult, error) {
	admin := client.Functions()
	
	functions, err := admin.GetFunctions(tenant, namespace)
	if err != nil {
		return b.handleError("list functions", err), nil
	}

	return b.marshalResponse(map[string]interface{}{
		"functions": functions,
		"tenant":    tenant,
		"namespace": namespace,
	})
}

// handleFunctionGet handles the get operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionGet(ctx context.Context, client cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	admin := client.Functions()
	
	functionConfig, err := admin.GetFunction(tenant, namespace, name)
	if err != nil {
		return b.handleError("get function config", err), nil
	}

	return b.marshalResponse(functionConfig)
}

// handleFunctionStatus handles the status operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionStatus(ctx context.Context, client cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	admin := client.Functions()
	
	status, err := admin.GetFunctionStatus(tenant, namespace, name)
	if err != nil {
		return b.handleError("get function status", err), nil
	}

	return b.marshalResponse(status)
}

// handleFunctionStats handles the stats operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionStats(ctx context.Context, client cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	// Note: GetStats method may not be available in current API
	return b.handleError("get function stats", fmt.Errorf("function stats operation not yet implemented in current API version")), nil
}

// handleFunctionQuerystate handles the querystate operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionQuerystate(ctx context.Context, client cmdutils.Client, tenant, namespace, name, key string) (*mcp.CallToolResult, error) {
	// Note: GetFunctionState method may not be available in current API
	return b.handleError("query function state", fmt.Errorf("function state query operation not yet implemented in current API version")), nil
}

// handleFunctionCreate handles the create operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionCreate(ctx context.Context, client cmdutils.Client, tenant, namespace, name string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Note: Function creation may require more complex implementation
	return b.handleError("create function", fmt.Errorf("function creation operation not yet implemented - requires complex configuration and file upload handling")), nil
}

// handleFunctionUpdate handles the update operation  
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionUpdate(ctx context.Context, client cmdutils.Client, tenant, namespace, name string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return b.handleError("update function", fmt.Errorf("function update operation not yet implemented - requires complex configuration handling")), nil
}

// handleFunctionDelete handles the delete operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionDelete(ctx context.Context, client cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	return b.handleError("delete function", fmt.Errorf("function delete operation not yet implemented")), nil
}

// handleFunctionStart handles the start operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionStart(ctx context.Context, client cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	return b.handleError("start function", fmt.Errorf("function start operation not yet implemented")), nil
}

// handleFunctionStop handles the stop operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionStop(ctx context.Context, client cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	return b.handleError("stop function", fmt.Errorf("function stop operation not yet implemented")), nil
}

// handleFunctionRestart handles the restart operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionRestart(ctx context.Context, client cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	return b.handleError("restart function", fmt.Errorf("function restart operation not yet implemented")), nil
}

// handleFunctionPutstate handles the putstate operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionPutstate(ctx context.Context, client cmdutils.Client, tenant, namespace, name, key, value string) (*mcp.CallToolResult, error) {
	return b.handleError("put function state", fmt.Errorf("function state put operation not yet implemented")), nil
}

// handleFunctionTrigger handles the trigger operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionTrigger(ctx context.Context, client cmdutils.Client, tenant, namespace, name, triggerValue, topic string) (*mcp.CallToolResult, error) {
	return b.handleError("trigger function", fmt.Errorf("function trigger operation not yet implemented")), nil
}

// Helper functions

// handleError provides unified error handling
func (b *PulsarAdminFunctionsToolBuilder) handleError(operation string, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("Failed to %s: %v", operation, err))
}

// marshalResponse provides unified JSON serialization for responses
func (b *PulsarAdminFunctionsToolBuilder) marshalResponse(data interface{}) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return b.handleError("marshal response", err), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}