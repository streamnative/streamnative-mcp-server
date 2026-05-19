// Copyright 2026 StreamNative
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
	"os"
	"strconv"
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/toolannotations"
	"gopkg.in/yaml.v2"
)

// PulsarAdminFunctionsToolBuilder implements the ToolBuilder interface for Pulsar admin functions operations
// It provides functionality to build Pulsar functions management tools
// /nolint:revive
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

	tools := []server.ServerTool{
		{
			Tool:    b.buildPulsarAdminFunctionsTool(toolModeRead),
			Handler: b.buildPulsarAdminFunctionsHandler(toolModeRead),
		},
	}
	if !config.ReadOnly {
		tools = append(tools, server.ServerTool{
			Tool:    b.buildPulsarAdminFunctionsTool(toolModeWrite),
			Handler: b.buildPulsarAdminFunctionsHandler(toolModeWrite),
		})
	}

	return tools, nil
}

// buildPulsarAdminFunctionsTool builds the Pulsar admin functions MCP tool definition
// Migrated from the original tool definition logic
func (b *PulsarAdminFunctionsToolBuilder) buildPulsarAdminFunctionsTool(mode toolMode) mcp.Tool {
	toolDesc := "Read Apache Pulsar Functions for stream processing. " +
		"Pulsar Functions are lightweight compute processes that can consume messages from one or more Pulsar topics, apply user-defined processing logic, and produce results to another topic. " +
		"This read-only tool lists functions and retrieves configuration, status, statistics, state, or package data. Functions require proper permissions to access their topics."

	operationDesc := "Operation to perform. Available operations:\n" +
		"- list: List all functions under a specific tenant and namespace\n" +
		"- get: Get the configuration of a function\n" +
		"- status: Get the runtime status of a function (instances, metrics)\n" +
		"- stats: Get detailed statistics of a function (throughput, processing latency)\n" +
		"- querystate: Query state stored by a stateful function for a specific key\n" +
		"- download: Download function package data from Pulsar to a local file"

	operationEnum := []string{"list", "get", "status", "stats", "querystate", "download"}
	toolName := "pulsar_admin_functions_read"
	annotation := toolannotations.ReadOnly("Read Pulsar Functions")
	if isToolModeWrite(mode) {
		toolDesc = "Manage Apache Pulsar Functions for stream processing. " +
			"This write tool deploys, updates, deletes, starts, stops, restarts, stores state for, triggers, or uploads packages for functions."
		operationDesc = "Operation to perform. Available operations:\n" +
			"- create: Deploy a new function with specified parameters\n" +
			"- update: Update the configuration of an existing function\n" +
			"- delete: Delete a function\n" +
			"- start: Start a stopped function\n" +
			"- stop: Stop a running function\n" +
			"- restart: Restart a function\n" +
			"- putstate: Store state in a function's state store\n" +
			"- trigger: Manually trigger a function with a specific value\n" +
			"- upload: Upload a local file into Pulsar function package storage"
		operationEnum = []string{"create", "update", "delete", "start", "stop", "restart", "putstate", "trigger", "upload"}
		toolName = "pulsar_admin_functions_write"
		annotation = toolannotations.Destructive("Manage Pulsar Functions")
	}

	tool := mcp.NewTool(toolName,
		mcp.WithDescription(toolDesc),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
			mcp.Enum(operationEnum...)),
		mcp.WithString("fqfn",
			mcp.Description("The Fully Qualified Function Name in the form tenant/namespace/name. "+
				"Mutually exclusive with tenant, namespace, and name parameters.")),
		mcp.WithString("tenant",
			mcp.Description("The tenant name. Tenants are the primary organizational unit in Pulsar, "+
				"providing multi-tenancy and resource isolation. Functions deployed within a tenant "+
				"inherit its permissions and resource quotas. Defaults to 'public' if not provided.")),
		mcp.WithString("namespace",
			mcp.Description("The namespace name. Namespaces are logical groupings of topics and functions "+
				"within a tenant. They encapsulate configuration policies and access control. "+
				"Functions in a namespace typically process topics within the same namespace. "+
				"Defaults to 'default' if not provided.")),
		mcp.WithString("name",
			mcp.Description("The function name. Required for operations that target one function. "+
				"Names should be descriptive of the function's purpose and must be unique within a namespace. "+
				"Function names are used in metrics, logs, and when addressing the function via APIs.")),
		mcp.WithString("instanceId",
			mcp.Description("Function instance ID for status/stats operations. "+
				"If not set, the aggregated status/stats of all instances will be returned.")),
		// Additional parameters for specific operations
		mcp.WithString("classname",
			mcp.Description("The fully qualified class name implementing the function. Required for 'create' operation, optional for 'update'. "+
				"For Java functions, this should be the class that implements pulsar function interfaces. "+
				"For Python, this MUST be in format of `<Python_filename_without_extension>.<ClassName>` - for example: "+
				"if file is '/path/to/exclamation.py' with class 'ExclamationFunction', classname must be 'exclamation.ExclamationFunction'; "+
				"if file is '/path/to/double_number.py' with class 'DoubleNumber', classname must be 'double_number.DoubleNumber'. "+
				"Common error: using just the class name 'DoubleNumber' (without filename prefix) will cause function creation to fail. "+
				"Go functions should specify the 'main' function of the binary.")),
		mcp.WithString("functionType",
			mcp.Description("The built-in Pulsar Function type. When set, it is translated to a builtin:// package URL. "+
				"Mutually exclusive with jar/py/go.")),
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
		mcp.WithString("topicsPattern",
			mcp.Description("Topic pattern to consume from. Mutually exclusive with inputs for typical use cases. "+
				"Example: persistent://tenant/namespace/topicPattern*")),
		mcp.WithObject("inputSpecs",
			mcp.Description("Map of input topics to consumer configuration (JSON object).")),
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
		mcp.WithNumber("cpu",
			mcp.Description("CPU cores allocated per function instance (docker runtime only).")),
		mcp.WithNumber("ram",
			mcp.Description("RAM bytes allocated per function instance (process/docker runtime only).")),
		mcp.WithNumber("disk",
			mcp.Description("Disk bytes allocated per function instance (docker runtime only).")),
		mcp.WithObject("userConfig",
			mcp.Description("User-defined config key/values. Optional for 'create' and 'update' operations. "+
				"Provides configuration parameters accessible to the function at runtime. "+
				"Specify as a JSON object with string, number, or boolean values. "+
				"Common configs include connection parameters, batch sizes, or feature toggles. "+
				"Example: {\"maxBatchSize\": 100, \"connectionString\": \"host:port\", \"debugMode\": true}")),
		mcp.WithObject("producerConfig",
			mcp.Description("Custom producer configuration as JSON object.")),
		mcp.WithString("logTopic",
			mcp.Description("Topic where function logs are written.")),
		mcp.WithString("schemaType",
			mcp.Description("Output schema type or schema class name.")),
		mcp.WithString("outputSerdeClassName",
			mcp.Description("SerDe class for output messages.")),
		mcp.WithObject("customSerdeInputs",
			mcp.Description("Map of input topics to SerDe class names (JSON object).")),
		mcp.WithObject("customSchemaInputs",
			mcp.Description("Map of input topics to Schema class names (JSON object).")),
		mcp.WithObject("customSchemaOutputs",
			mcp.Description("Map of output topics to Schema properties (JSON object).")),
		mcp.WithString("inputTypeClassName",
			mcp.Description("Input type class name.")),
		mcp.WithString("outputTypeClassName",
			mcp.Description("Output type class name.")),
		mcp.WithString("processingGuarantees",
			mcp.Description("Processing guarantees (delivery semantics).")),
		mcp.WithBoolean("retainOrdering",
			mcp.Description("Process messages in order.")),
		mcp.WithBoolean("retainKeyOrdering",
			mcp.Description("Process messages in key order.")),
		mcp.WithString("batchBuilder",
			mcp.Description("Batch builder type (DEFAULT or KEY_BASED).")),
		mcp.WithBoolean("forwardSourceMessageProperty",
			mcp.Description("Forward input message properties to output topic.")),
		mcp.WithString("subsName",
			mcp.Description("Subscription name for input-topic consumer.")),
		mcp.WithString("subsPosition",
			mcp.Description("Subscription position for input-topic consumer.")),
		mcp.WithBoolean("skipToLatest",
			mcp.Description("Skip to latest message on function instance restart.")),
		mcp.WithBoolean("autoAck",
			mcp.Description("Whether the framework acknowledges messages automatically.")),
		mcp.WithNumber("timeoutMs",
			mcp.Description("Message timeout in milliseconds.")),
		mcp.WithNumber("maxMessageRetries",
			mcp.Description("Number of message retries before giving up.")),
		mcp.WithString("deadLetterTopic",
			mcp.Description("Topic for messages that are not processed successfully.")),
		mcp.WithString("customRuntimeOptions",
			mcp.Description("Custom runtime options for the function.")),
		mcp.WithObject("secrets",
			mcp.Description("Secrets map for the function (JSON object).")),
		mcp.WithBoolean("cleanupSubscription",
			mcp.Description("Whether to delete the subscription when the function is deleted.")),
		mcp.WithNumber("windowLengthCount",
			mcp.Description("Number of messages per window.")),
		mcp.WithNumber("windowLengthDurationMs",
			mcp.Description("Window length in milliseconds.")),
		mcp.WithNumber("slidingIntervalCount",
			mcp.Description("Number of messages after which the window slides.")),
		mcp.WithNumber("slidingIntervalDurationMs",
			mcp.Description("Window sliding interval in milliseconds.")),
		mcp.WithString("functionConfigFile",
			mcp.Description("Path to a YAML config file that specifies the function configuration.")),
		mcp.WithString("sourceFile",
			mcp.Description("Path to the local file that should be uploaded into Pulsar function package storage. "+
				"Required for the 'upload' operation.")),
		mcp.WithString("path",
			mcp.Description("Pulsar package storage path used by package transfer operations. For downloads, provide this to read directly from package storage.")),
		mcp.WithString("destinationFile",
			mcp.Description("Local file path where downloaded function package data should be written. "+
				"Required for the 'download' operation.")),
		mcp.WithBoolean("updateAuthData",
			mcp.Description("Whether to update auth data on update operations.")),
		mcp.WithString("key",
			mcp.Description("The state key. Required for operations that access one value in the function's state store. "+
				"Keys should be reasonable in length and follow a consistent pattern. "+
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
			mcp.Description("The value with which to trigger the function. Required for 'trigger' operation unless triggerFile is set. "+
				"This value will be passed to the function as if it were a message from the input topic. "+
				"String values are sent as is; for typed values, ensure proper formatting based on function expectations. "+
				"The function processes this value just like a normal message.")),
		mcp.WithString("triggerFile",
			mcp.Description("Path to a file containing the trigger value. Required for 'trigger' operation unless triggerValue is set.")),
		annotation,
	)
	if isToolModeWrite(mode) {
		removeToolInputSchemaProperties(&tool, []string{"instanceId", "destinationFile"})
	} else {
		pruneToolInputSchema(&tool, []string{"operation", "fqfn", "tenant", "namespace", "name", "instanceId", "key", "path", "destinationFile"})
	}
	return tool
}

// buildPulsarAdminFunctionsHandler builds the Pulsar admin functions handler function
// Migrated from the original handler logic
func (b *PulsarAdminFunctionsToolBuilder) buildPulsarAdminFunctionsHandler(mode toolMode) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return mcp.NewToolResultError("Pulsar session not found in context"), nil
		}

		client, err := session.GetAdminV3Client()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get Pulsar client: %v", err)), nil
		}

		// Extract and validate operation parameter
		operation, err := request.RequireString("operation")
		if err != nil {
			return b.handleError("get operation", err), nil
		}

		// Check if the operation is valid
		if !isSupportedFunctionOperation(operation) {
			return b.handleError("validate operation", fmt.Errorf("invalid operation: '%s'. Supported operations: list, get, status, stats, querystate, create, update, delete, download, start, stop, restart, putstate, trigger, upload", operation)), nil
		}

		if !validateModeOperation(mode, operation, readOnlyRestrictedFunctionOperations) {
			return b.handleError("check permissions", fmt.Errorf("operation %q is not available in %s mode", operation, mode)), nil
		}

		var identity functionIdentity
		if operation != "download" && operation != "upload" {
			identity, err = b.parseFunctionIdentity(request, operation)
			if err != nil {
				return b.handleError("get function identity", err), nil
			}
		}

		// Handle operation using delegated handlers
		switch operation {
		case "list":
			return b.handleFunctionList(ctx, client, identity.Tenant, identity.Namespace)
		case "get":
			return b.handleFunctionGet(ctx, client, identity.Tenant, identity.Namespace, identity.Name)
		case "status":
			instanceID, ok, err := parseOptionalIntArg(request.GetArguments(), "instanceId")
			if err != nil {
				return b.handleError("get instanceId", err), nil
			}
			if ok {
				return b.handleFunctionStatusWithInstance(ctx, client, identity.Tenant, identity.Namespace, identity.Name, instanceID)
			}
			return b.handleFunctionStatus(ctx, client, identity.Tenant, identity.Namespace, identity.Name)
		case "stats":
			instanceID, ok, err := parseOptionalIntArg(request.GetArguments(), "instanceId")
			if err != nil {
				return b.handleError("get instanceId", err), nil
			}
			if ok {
				return b.handleFunctionStatsWithInstance(ctx, client, identity.Tenant, identity.Namespace, identity.Name, instanceID)
			}
			return b.handleFunctionStats(ctx, client, identity.Tenant, identity.Namespace, identity.Name)
		case "querystate":
			key, err := request.RequireString("key")
			if err != nil {
				return b.handleError("get key", fmt.Errorf("missing required parameter 'key' for operation 'querystate': %v. A key is required to look up state in the function's state store", err)), nil
			}
			return b.handleFunctionQuerystate(ctx, client, identity.Tenant, identity.Namespace, identity.Name, key)
		case "create":
			return b.handleFunctionCreate(ctx, client, identity.Tenant, identity.Namespace, identity.Name, request)
		case "update":
			return b.handleFunctionUpdate(ctx, client, identity.Tenant, identity.Namespace, identity.Name, request)
		case "delete":
			return b.handleFunctionDelete(ctx, client, identity.Tenant, identity.Namespace, identity.Name)
		case "download":
			return b.handleFunctionDownload(ctx, client, request)
		case "start":
			return b.handleFunctionStart(ctx, client, identity.Tenant, identity.Namespace, identity.Name)
		case "stop":
			return b.handleFunctionStop(ctx, client, identity.Tenant, identity.Namespace, identity.Name)
		case "restart":
			return b.handleFunctionRestart(ctx, client, identity.Tenant, identity.Namespace, identity.Name)
		case "putstate":
			key, err := request.RequireString("key")
			if err != nil {
				return b.handleError("get key", fmt.Errorf("missing required parameter 'key' for operation 'putstate': %v. A key is required to store state in the function's state store", err)), nil
			}
			value, err := request.RequireString("value")
			if err != nil {
				return b.handleError("get value", fmt.Errorf("missing required parameter 'value' for operation 'putstate': %v. A value is required to store state in the function's state store", err)), nil
			}
			return b.handleFunctionPutstate(ctx, client, identity.Tenant, identity.Namespace, identity.Name, key, value)
		case "trigger":
			args := request.GetArguments()
			triggerValue, hasValue := getStringArg(args, "triggerValue")
			triggerFile, hasFile := getStringArg(args, "triggerFile")
			if hasValue && triggerValue == "" {
				return b.handleError("get triggerValue", fmt.Errorf("triggerValue cannot be empty")), nil
			}
			if hasFile && triggerFile == "" {
				return b.handleError("get triggerFile", fmt.Errorf("triggerFile cannot be empty")), nil
			}
			if hasValue == hasFile {
				return b.handleError("validate trigger arguments", fmt.Errorf("exactly one of triggerValue or triggerFile must be provided")), nil
			}
			topic := request.GetString("topic", "")
			return b.handleFunctionTrigger(ctx, client, identity.Tenant, identity.Namespace, identity.Name, triggerValue, triggerFile, topic)
		case "upload":
			return b.handleFunctionUpload(ctx, client, request)
		default:
			return b.handleError("handle operation", fmt.Errorf("unsupported operation: %s", operation)), nil
		}
	}
}

// Helper functions - delegated operation handlers

// handleFunctionList handles the list operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionList(_ context.Context, client cmdutils.Client, tenant, namespace string) (*mcp.CallToolResult, error) {
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
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionGet(_ context.Context, client cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	admin := client.Functions()

	functionConfig, err := admin.GetFunction(tenant, namespace, name)
	if err != nil {
		return b.handleError("get function config", err), nil
	}

	return b.marshalResponse(functionConfig)
}

// handleFunctionStatus handles the status operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionStatus(_ context.Context, client cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	admin := client.Functions()

	status, err := admin.GetFunctionStatus(tenant, namespace, name)
	if err != nil {
		return b.handleError("get function status", err), nil
	}

	return b.marshalResponse(status)
}

// handleFunctionStatusWithInstance handles the status operation for a single instance
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionStatusWithInstance(_ context.Context, client cmdutils.Client, tenant, namespace, name string, instanceID int) (*mcp.CallToolResult, error) {
	admin := client.Functions()

	status, err := admin.GetFunctionStatusWithInstanceID(tenant, namespace, name, instanceID)
	if err != nil {
		return b.handleError("get function status", err), nil
	}

	return b.marshalResponse(status)
}

// handleFunctionStats handles the stats operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionStats(_ context.Context, client cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	admin := client.Functions()

	stats, err := admin.GetFunctionStats(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get stats for function '%s' in tenant '%s' namespace '%s': %v. Verify the function exists and is running.",
			name, tenant, namespace, err)), nil
	}

	return b.marshalResponse(stats)
}

// handleFunctionStatsWithInstance handles the stats operation for a single instance
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionStatsWithInstance(_ context.Context, client cmdutils.Client, tenant, namespace, name string, instanceID int) (*mcp.CallToolResult, error) {
	admin := client.Functions()

	stats, err := admin.GetFunctionStatsWithInstanceID(tenant, namespace, name, instanceID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get stats for function '%s' instance '%d' in tenant '%s' namespace '%s': %v. Verify the function exists and is running.",
			name, instanceID, tenant, namespace, err)), nil
	}

	return b.marshalResponse(stats)
}

// handleFunctionQuerystate handles the querystate operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionQuerystate(_ context.Context, client cmdutils.Client, tenant, namespace, name, key string) (*mcp.CallToolResult, error) {
	admin := client.Functions()

	state, err := admin.GetFunctionState(tenant, namespace, name, key)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to query state for key '%s' in function '%s' (tenant '%s' namespace '%s'): %v. Verify the function exists and has state enabled.",
			key, name, tenant, namespace, err)), nil
	}

	return b.marshalResponse(map[string]interface{}{
		"key":   key,
		"value": state,
		"function": map[string]string{
			"tenant":    tenant,
			"namespace": namespace,
			"name":      name,
		},
	})
}

// handleFunctionCreate handles the create operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionCreate(_ context.Context, client cmdutils.Client, tenant, namespace, name string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Build function configuration from request parameters to validate
	functionConfig, err := b.buildFunctionConfig(tenant, namespace, name, request)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build function configuration for '%s' in tenant '%s' namespace '%s': %v. Please verify all required parameters are provided correctly.",
			name, tenant, namespace, err)), nil
	}
	if err := validateFunctionConfigs(functionConfig); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to validate function configuration for '%s' in tenant '%s' namespace '%s': %v.",
			functionConfig.Name, functionConfig.Tenant, functionConfig.Namespace, err)), nil
	}

	admin := client.Functions()
	packagePath := resolvePackagePath(functionConfig)
	if packagePath == "" {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create function '%s': missing function package path", functionConfig.Name)), nil
	}

	if isPackageURLSupported(packagePath) {
		err = admin.CreateFuncWithURL(functionConfig, packagePath)
	} else {
		err = admin.CreateFunc(functionConfig, packagePath)
	}
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create function '%s' in tenant '%s' namespace '%s': %v. Verify the function configuration is valid.",
			functionConfig.Name, functionConfig.Tenant, functionConfig.Namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Created function '%s' successfully in tenant '%s' namespace '%s'. The function configuration has been created.",
		functionConfig.Name, functionConfig.Tenant, functionConfig.Namespace)), nil
}

// handleFunctionUpdate handles the update operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionUpdate(_ context.Context, client cmdutils.Client, tenant, namespace, name string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	admin := client.Functions()

	// Build function configuration from request parameters
	config, err := b.buildFunctionConfig(tenant, namespace, name, request)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build function configuration for '%s' in tenant '%s' namespace '%s': %v. Please verify all parameters are provided correctly.",
			name, tenant, namespace, err)), nil
	}
	if err := checkArgsForUpdate(config); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to validate update configuration for '%s' in tenant '%s' namespace '%s': %v.",
			config.Name, config.Tenant, config.Namespace, err)), nil
	}

	// Update the function
	updateOptions := utils.NewUpdateOptions()
	if updateAuthData, ok := getBoolArg(request.GetArguments(), "updateAuthData"); ok {
		updateOptions.UpdateAuthData = updateAuthData
	}

	packagePath := resolvePackagePath(config)
	if packagePath != "" && isPackageURLSupported(packagePath) {
		err = admin.UpdateFunctionWithURL(config, packagePath, updateOptions)
	} else {
		err = admin.UpdateFunction(config, packagePath, updateOptions)
	}
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to update function '%s' in tenant '%s' namespace '%s': %v. Verify the function exists and the configuration is valid.",
			config.Name, config.Tenant, config.Namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Updated function '%s' successfully in tenant '%s' namespace '%s'. The function configuration has been modified.",
		config.Name, config.Tenant, config.Namespace)), nil
}

// handleFunctionDelete handles the delete operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionDelete(_ context.Context, client cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	admin := client.Functions()

	err := admin.DeleteFunction(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete function '%s' in tenant '%s' namespace '%s': %v. Verify the function exists and you have deletion permissions.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Deleted function '%s' successfully from tenant '%s' namespace '%s'. All running instances have been terminated.",
		name, tenant, namespace)), nil
}

// handleFunctionDownload handles downloading function package data from Pulsar.
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionDownload(_ context.Context, client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	admin := client.Functions()

	target, err := b.parseFunctionDownloadTarget(request)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to parse function download request: %v", err)), nil
	}

	if target.UsePath {
		err = admin.DownloadFunction(target.Path, target.DestinationFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to download function package from Pulsar path '%s' to '%s': %v.",
				target.Path, target.DestinationFile, err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Downloaded function package from Pulsar path '%s' to '%s' successfully.",
			target.Path, target.DestinationFile)), nil
	}

	err = admin.DownloadFunctionByNs(target.DestinationFile, target.Identity.Tenant, target.Identity.Namespace, target.Identity.Name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to download function '%s' in tenant '%s' namespace '%s' to '%s': %v.",
			target.Identity.Name, target.Identity.Tenant, target.Identity.Namespace, target.DestinationFile, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Downloaded function '%s' from tenant '%s' namespace '%s' to '%s' successfully.",
		target.Identity.Name, target.Identity.Tenant, target.Identity.Namespace, target.DestinationFile)), nil
}

// handleFunctionUpload handles uploading a local file to Pulsar function package storage.
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionUpload(_ context.Context, client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	admin := client.Functions()

	sourceFile, err := request.RequireString("sourceFile")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'sourceFile' for function.upload: %v", err)), nil
	}
	sourceFile = strings.TrimSpace(sourceFile)
	if sourceFile == "" {
		return mcp.NewToolResultError("Parameter 'sourceFile' for function.upload cannot be empty"), nil
	}

	path, err := request.RequireString("path")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'path' for function.upload: %v", err)), nil
	}
	path = strings.TrimSpace(path)
	if path == "" {
		return mcp.NewToolResultError("Parameter 'path' for function.upload cannot be empty"), nil
	}

	err = admin.Upload(sourceFile, path)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to upload function package '%s' to Pulsar path '%s': %v.",
			sourceFile, path, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Uploaded function package '%s' to Pulsar path '%s' successfully.",
		sourceFile, path)), nil
}

// handleFunctionStart handles the start operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionStart(_ context.Context, client cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	admin := client.Functions()

	err := admin.StartFunction(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to start function '%s' in tenant '%s' namespace '%s': %v. Verify the function exists and is not already running.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Started function '%s' successfully in tenant '%s' namespace '%s'. The function instances are now processing messages.",
		name, tenant, namespace)), nil
}

// handleFunctionStop handles the stop operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionStop(_ context.Context, client cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	admin := client.Functions()

	err := admin.StopFunction(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to stop function '%s' in tenant '%s' namespace '%s': %v. Verify the function exists and is currently running.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Stopped function '%s' successfully in tenant '%s' namespace '%s'. The function will no longer process messages until restarted.",
		name, tenant, namespace)), nil
}

// handleFunctionRestart handles the restart operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionRestart(_ context.Context, client cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	admin := client.Functions()

	err := admin.RestartFunction(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to restart function '%s' in tenant '%s' namespace '%s': %v. Verify the function exists and is properly deployed.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Restarted function '%s' successfully in tenant '%s' namespace '%s'. All function instances have been restarted.",
		name, tenant, namespace)), nil
}

// handleFunctionPutstate handles the putstate operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionPutstate(_ context.Context, client cmdutils.Client, tenant, namespace, name, key, value string) (*mcp.CallToolResult, error) {
	admin := client.Functions()

	err := admin.PutFunctionState(tenant, namespace, name, utils.FunctionState{
		Key:         key,
		StringValue: value,
	})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to put state for key '%s' in function '%s' (tenant '%s' namespace '%s'): %v. Verify the function exists and has state enabled.",
			key, name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully stored state for key '%s' in function '%s' (tenant '%s' namespace '%s'). State value has been updated.",
		key, name, tenant, namespace)), nil
}

// handleFunctionTrigger handles the trigger operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionTrigger(_ context.Context, client cmdutils.Client, tenant, namespace, name, triggerValue, triggerFile, topic string) (*mcp.CallToolResult, error) {
	admin := client.Functions()

	var err error
	var result string
	if topic != "" {
		// Trigger with specific topic
		result, err = admin.TriggerFunction(tenant, namespace, name, topic, triggerValue, triggerFile)
	} else {
		// Trigger without specific topic (uses first input topic)
		result, err = admin.TriggerFunction(tenant, namespace, name, "", triggerValue, triggerFile)
	}

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to trigger function '%s' in tenant '%s' namespace '%s': %v. Verify the function exists and is running.",
			name, tenant, namespace, err)), nil
	}

	var message string
	if topic != "" {
		message = fmt.Sprintf("Successfully triggered function '%s' in tenant '%s' namespace '%s' with topic '%s'. Result: %s",
			name, tenant, namespace, topic, result)
	} else {
		message = fmt.Sprintf("Successfully triggered function '%s' in tenant '%s' namespace '%s'. Result: %s",
			name, tenant, namespace, result)
	}

	return mcp.NewToolResultText(message), nil
}

// Helper functions

type functionIdentity struct {
	Tenant    string
	Namespace string
	Name      string
}

type functionDownloadTarget struct {
	DestinationFile string
	Path            string
	Identity        functionIdentity
	UsePath         bool
}

const (
	defaultTenant    = "public"
	defaultNamespace = "default"
)

var supportedFunctionOperations = map[string]struct{}{
	"list":       {},
	"get":        {},
	"status":     {},
	"stats":      {},
	"querystate": {},
	"create":     {},
	"update":     {},
	"delete":     {},
	"download":   {},
	"start":      {},
	"stop":       {},
	"restart":    {},
	"putstate":   {},
	"trigger":    {},
	"upload":     {},
}

var readOnlyRestrictedFunctionOperations = map[string]struct{}{
	"create":   {},
	"update":   {},
	"delete":   {},
	"start":    {},
	"stop":     {},
	"restart":  {},
	"putstate": {},
	"trigger":  {},
	"upload":   {},
}

func isSupportedFunctionOperation(operation string) bool {
	_, ok := supportedFunctionOperations[operation]
	return ok
}

func isReadOnlyRestrictedFunctionOperation(operation string) bool {
	_, ok := readOnlyRestrictedFunctionOperations[operation]
	return ok
}

func (b *PulsarAdminFunctionsToolBuilder) parseFunctionIdentity(request mcp.CallToolRequest, operation string) (functionIdentity, error) {
	args := request.GetArguments()
	fqfn, _ := getStringArg(args, "fqfn")
	tenant, _ := getStringArg(args, "tenant")
	namespace, _ := getStringArg(args, "namespace")
	name, _ := getStringArg(args, "name")

	if fqfn != "" {
		if operation == "list" {
			return functionIdentity{}, fmt.Errorf("fqfn is not supported for operation 'list'")
		}
		if tenant != "" || namespace != "" || name != "" {
			return functionIdentity{}, fmt.Errorf("fqfn cannot be combined with tenant, namespace, or name")
		}
		parsed, err := parseFullyQualifiedFunctionName(fqfn)
		if err != nil {
			return functionIdentity{}, err
		}
		tenant, namespace, name = parsed.Tenant, parsed.Namespace, parsed.Name
	}

	if tenant == "" {
		tenant = defaultTenant
	}
	if namespace == "" {
		namespace = defaultNamespace
	}

	requiresName := operation != "list" && operation != "create" && operation != "update"
	if requiresName && name == "" {
		return functionIdentity{}, fmt.Errorf("you must specify a name for the function or a fully qualified function name (fqfn)")
	}

	return functionIdentity{
		Tenant:    tenant,
		Namespace: namespace,
		Name:      name,
	}, nil
}

func (b *PulsarAdminFunctionsToolBuilder) parseFunctionDownloadTarget(request mcp.CallToolRequest) (functionDownloadTarget, error) {
	destinationFile, err := request.RequireString("destinationFile")
	if err != nil {
		return functionDownloadTarget{}, fmt.Errorf("missing required parameter 'destinationFile' for operation 'download': %w", err)
	}
	destinationFile = strings.TrimSpace(destinationFile)
	if destinationFile == "" {
		return functionDownloadTarget{}, fmt.Errorf("parameter 'destinationFile' for operation 'download' cannot be empty")
	}

	path, _ := getStringArg(request.GetArguments(), "path")
	path = strings.TrimSpace(path)
	if path != "" {
		return functionDownloadTarget{
			DestinationFile: destinationFile,
			Path:            path,
			UsePath:         true,
		}, nil
	}

	identity, err := b.parseFunctionIdentity(request, "download")
	if err != nil {
		return functionDownloadTarget{}, err
	}

	return functionDownloadTarget{
		DestinationFile: destinationFile,
		Identity:        identity,
	}, nil
}

// buildFunctionConfig builds a Pulsar Function configuration from MCP request parameters
func (b *PulsarAdminFunctionsToolBuilder) buildFunctionConfig(tenant, namespace, name string, request mcp.CallToolRequest) (*utils.FunctionConfig, error) {
	config := &utils.FunctionConfig{}
	args := request.GetArguments()

	if configFile, ok := getStringArg(args, "functionConfigFile"); ok && configFile != "" {
		//nolint:gosec
		data, err := os.ReadFile(configFile)
		if err == nil {
			if err := yaml.Unmarshal(data, config); err != nil {
				return nil, fmt.Errorf("unmarshal function config file error: %w", err)
			}
		} else if !os.IsNotExist(err) {
			return nil, fmt.Errorf("load function config file failed: %w", err)
		}
	}

	if tenant != "" {
		config.Tenant = tenant
	}
	if namespace != "" {
		config.Namespace = namespace
	}
	if name != "" {
		config.Name = name
	}

	if fqfn, ok := getStringArg(args, "fqfn"); ok && fqfn != "" {
		parsed, err := parseFullyQualifiedFunctionName(fqfn)
		if err != nil {
			return nil, err
		}
		config.Tenant = parsed.Tenant
		config.Namespace = parsed.Namespace
		config.Name = parsed.Name
	}

	if classname, ok := getStringArg(args, "classname"); ok && classname != "" {
		config.ClassName = classname
	}

	if inputsValue, exists := args["inputs"]; exists && inputsValue != nil {
		inputs, err := parseStringSlice(inputsValue)
		if err != nil {
			return nil, fmt.Errorf("invalid inputs: %w", err)
		}
		config.Inputs = inputs
	}

	if topicsPattern, ok := getStringArg(args, "topicsPattern"); ok && topicsPattern != "" {
		config.TopicsPattern = &topicsPattern
	}

	if inputSpecsValue, exists := args["inputSpecs"]; exists && inputSpecsValue != nil {
		inputSpecs, err := decodeConsumerConfigMap(inputSpecsValue)
		if err != nil {
			return nil, fmt.Errorf("invalid inputSpecs: %w", err)
		}
		config.InputSpecs = inputSpecs
	}

	if output, ok := getStringArg(args, "output"); ok && output != "" {
		config.Output = output
	}

	if parallelism, ok, err := parseOptionalIntArg(args, "parallelism"); err != nil {
		return nil, err
	} else if ok {
		config.Parallelism = parallelism
	}

	if cpu, ok, err := parseOptionalFloatArg(args, "cpu"); err != nil {
		return nil, err
	} else if ok {
		config.Resources = ensureResources(config.Resources)
		config.Resources.CPU = cpu
	}

	if ram, ok, err := parseOptionalInt64Arg(args, "ram"); err != nil {
		return nil, err
	} else if ok {
		config.Resources = ensureResources(config.Resources)
		config.Resources.RAM = ram
	}

	if disk, ok, err := parseOptionalInt64Arg(args, "disk"); err != nil {
		return nil, err
	} else if ok {
		config.Resources = ensureResources(config.Resources)
		config.Resources.Disk = disk
	}

	if userConfigValue, exists := args["userConfig"]; exists && userConfigValue != nil {
		userConfig, err := decodeInterfaceMap(userConfigValue)
		if err != nil {
			return nil, fmt.Errorf("invalid userConfig: %w", err)
		}
		config.UserConfig = userConfig
	}

	if producerConfigValue, exists := args["producerConfig"]; exists && producerConfigValue != nil {
		producerConfig := &utils.ProducerConfig{}
		if err := decodeInto(producerConfigValue, producerConfig); err != nil {
			return nil, fmt.Errorf("invalid producerConfig: %w", err)
		}
		config.ProducerConfig = producerConfig
	}

	if logTopic, ok := getStringArg(args, "logTopic"); ok && logTopic != "" {
		config.LogTopic = logTopic
	}

	if schemaType, ok := getStringArg(args, "schemaType"); ok && schemaType != "" {
		config.OutputSchemaType = schemaType
	}

	if outputSerdeClassName, ok := getStringArg(args, "outputSerdeClassName"); ok && outputSerdeClassName != "" {
		config.OutputSerdeClassName = outputSerdeClassName
	}

	if customSerdeInputsValue, exists := args["customSerdeInputs"]; exists && customSerdeInputsValue != nil {
		customSerdeInputs, err := decodeStringMap(customSerdeInputsValue)
		if err != nil {
			return nil, fmt.Errorf("invalid customSerdeInputs: %w", err)
		}
		config.CustomSerdeInputs = customSerdeInputs
	}

	if customSchemaInputsValue, exists := args["customSchemaInputs"]; exists && customSchemaInputsValue != nil {
		customSchemaInputs, err := decodeStringMap(customSchemaInputsValue)
		if err != nil {
			return nil, fmt.Errorf("invalid customSchemaInputs: %w", err)
		}
		config.CustomSchemaInputs = customSchemaInputs
	}

	if customSchemaOutputsValue, exists := args["customSchemaOutputs"]; exists && customSchemaOutputsValue != nil {
		customSchemaOutputs, err := decodeStringMap(customSchemaOutputsValue)
		if err != nil {
			return nil, fmt.Errorf("invalid customSchemaOutputs: %w", err)
		}
		config.CustomSchemaOutputs = customSchemaOutputs
	}

	if inputTypeClassName, ok := getStringArg(args, "inputTypeClassName"); ok && inputTypeClassName != "" {
		config.InputTypeClassName = inputTypeClassName
	}

	if outputTypeClassName, ok := getStringArg(args, "outputTypeClassName"); ok && outputTypeClassName != "" {
		config.OutputTypeClassName = outputTypeClassName
	}

	if processingGuarantees, ok := getStringArg(args, "processingGuarantees"); ok && processingGuarantees != "" {
		config.ProcessingGuarantees = processingGuarantees
	}

	if retainOrdering, ok := getBoolArg(args, "retainOrdering"); ok {
		config.RetainOrdering = retainOrdering
	}

	if retainKeyOrdering, ok := getBoolArg(args, "retainKeyOrdering"); ok {
		config.RetainKeyOrdering = retainKeyOrdering
	}

	if batchBuilder, ok := getStringArg(args, "batchBuilder"); ok && batchBuilder != "" {
		config.BatchBuilder = batchBuilder
	}

	if forwardSourceMessageProperty, ok := getBoolArg(args, "forwardSourceMessageProperty"); ok {
		config.ForwardSourceMessageProperty = forwardSourceMessageProperty
	} else {
		config.ForwardSourceMessageProperty = true
	}

	if subsName, ok := getStringArg(args, "subsName"); ok && subsName != "" {
		config.SubName = subsName
	}

	if subsPosition, ok := getStringArg(args, "subsPosition"); ok && subsPosition != "" {
		config.SubscriptionPosition = subsPosition
	}

	if skipToLatest, ok := getBoolArg(args, "skipToLatest"); ok {
		config.SkipToLatest = skipToLatest
	}

	if timeoutMs, ok, err := parseOptionalInt64Arg(args, "timeoutMs"); err != nil {
		return nil, err
	} else if ok {
		config.TimeoutMs = &timeoutMs
	}

	if maxMessageRetries, ok, err := parseOptionalIntArg(args, "maxMessageRetries"); err != nil {
		return nil, err
	} else if ok {
		config.MaxMessageRetries = &maxMessageRetries
	}

	if deadLetterTopic, ok := getStringArg(args, "deadLetterTopic"); ok && deadLetterTopic != "" {
		config.DeadLetterTopic = deadLetterTopic
	}

	if customRuntimeOptions, ok := getStringArg(args, "customRuntimeOptions"); ok && customRuntimeOptions != "" {
		config.CustomRuntimeOptions = customRuntimeOptions
	}

	if secretsValue, exists := args["secrets"]; exists && secretsValue != nil {
		secrets, err := decodeInterfaceMap(secretsValue)
		if err != nil {
			return nil, fmt.Errorf("invalid secrets: %w", err)
		}
		config.Secrets = secrets
	}

	if cleanupSubscription, ok := getBoolArg(args, "cleanupSubscription"); ok {
		config.CleanupSubscription = cleanupSubscription
	} else {
		config.CleanupSubscription = true
	}

	if autoAck, ok := getBoolArg(args, "autoAck"); ok {
		config.AutoAck = autoAck
	} else {
		config.AutoAck = true
	}

	if windowLengthCount, ok, err := parseOptionalIntArg(args, "windowLengthCount"); err != nil {
		return nil, err
	} else if ok {
		config.WindowConfig = ensureWindowConfig(config.WindowConfig)
		value := windowLengthCount
		config.WindowConfig.WindowLengthCount = &value
	}

	if windowLengthDurationMs, ok, err := parseOptionalInt64Arg(args, "windowLengthDurationMs"); err != nil {
		return nil, err
	} else if ok {
		config.WindowConfig = ensureWindowConfig(config.WindowConfig)
		value := windowLengthDurationMs
		config.WindowConfig.WindowLengthDurationMs = &value
	}

	if slidingIntervalCount, ok, err := parseOptionalIntArg(args, "slidingIntervalCount"); err != nil {
		return nil, err
	} else if ok {
		config.WindowConfig = ensureWindowConfig(config.WindowConfig)
		value := slidingIntervalCount
		config.WindowConfig.SlidingIntervalCount = &value
	}

	if slidingIntervalDurationMs, ok, err := parseOptionalInt64Arg(args, "slidingIntervalDurationMs"); err != nil {
		return nil, err
	} else if ok {
		config.WindowConfig = ensureWindowConfig(config.WindowConfig)
		value := slidingIntervalDurationMs
		config.WindowConfig.SlidingIntervalDurationMs = &value
	}

	functionTypeArg, _ := getStringArg(args, "functionType")
	jarArg, _ := getStringArg(args, "jar")
	pyArg, _ := getStringArg(args, "py")
	goArg, _ := getStringArg(args, "go")

	functionTypeValue := functionTypeArg
	if functionTypeValue == "" && config.FunctionType != nil {
		functionTypeValue = *config.FunctionType
	}
	jarValue := jarArg
	if jarValue == "" && config.Jar != nil {
		jarValue = *config.Jar
	}
	pyValue := pyArg
	if pyValue == "" && config.Py != nil {
		pyValue = *config.Py
	}
	goValue := goArg
	if goValue == "" && config.Go != nil {
		goValue = *config.Go
	}

	if functionTypeValue != "" && (jarValue != "" || pyValue != "" || goValue != "") {
		return nil, fmt.Errorf("functionType is mutually exclusive with jar, py, and go")
	}
	providedPackages := 0
	if jarValue != "" {
		providedPackages++
	}
	if pyValue != "" {
		providedPackages++
	}
	if goValue != "" {
		providedPackages++
	}
	if providedPackages > 1 {
		return nil, fmt.Errorf("jar, py, and go are mutually exclusive")
	}

	if functionTypeValue != "" {
		jar := fmt.Sprintf("builtin://%s", functionTypeValue)
		config.Jar = &jar
	}

	if jarValue != "" {
		config.Jar = &jarValue
	}

	if pyValue != "" {
		config.Py = &pyValue
	}

	if goValue != "" {
		config.Go = &goValue
	}

	if config.Go != nil {
		config.Runtime = utils.GoRuntime
	}
	if config.Py != nil {
		config.Runtime = utils.PythonRuntime
	}
	if config.Jar != nil {
		config.Runtime = utils.JavaRuntime
	}

	if config.Parallelism <= 0 {
		config.Parallelism = 1
	}

	if config.UserConfig == nil {
		config.UserConfig = make(map[string]interface{})
	}
	if config.Secrets == nil {
		config.Secrets = make(map[string]interface{})
	}

	formatFunctionConfig(config)

	return config, nil
}

func parseFullyQualifiedFunctionName(fqfn string) (functionIdentity, error) {
	parts := strings.Split(fqfn, "/")
	if len(parts) != 3 {
		return functionIdentity{}, fmt.Errorf("fully qualified function names must be of the form tenant/namespace/name")
	}
	if parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return functionIdentity{}, fmt.Errorf("fully qualified function names must be of the form tenant/namespace/name")
	}
	return functionIdentity{
		Tenant:    parts[0],
		Namespace: parts[1],
		Name:      parts[2],
	}, nil
}

func validateFunctionConfigs(functionConfig *utils.FunctionConfig) error {
	if functionConfig.Name == "" {
		inferMissingFunctionName(functionConfig)
	}
	if functionConfig.Tenant == "" {
		inferMissingTenant(functionConfig)
	}
	if functionConfig.Namespace == "" {
		inferMissingNamespace(functionConfig)
	}

	switch numProvidedStrings(functionConfig.Jar, functionConfig.Py, functionConfig.Go) {
	case 0:
		return fmt.Errorf("either a Java jar or a Python file or a Go executable binary needs to be specified for the function")
	case 1:
	default:
		return fmt.Errorf("either a Java jar or a Python file or a Go executable binary needs to be specified for the function, cannot specify more than one")
	}

	if functionConfig.Jar != nil && !strings.HasPrefix(*functionConfig.Jar, "builtin://") &&
		!isPackageURLSupported(*functionConfig.Jar) &&
		!fileExists(*functionConfig.Jar) {
		return fmt.Errorf("the specified jar file does not exist")
	}
	if functionConfig.Py != nil && !isPackageURLSupported(*functionConfig.Py) && !fileExists(*functionConfig.Py) {
		return fmt.Errorf("the specified py file does not exist")
	}
	if functionConfig.Go != nil && !isPackageURLSupported(*functionConfig.Go) && !fileExists(*functionConfig.Go) {
		return fmt.Errorf("the specified go file does not exist")
	}

	if functionConfig.Go != nil {
		functionConfig.Runtime = utils.GoRuntime
	}
	if functionConfig.Py != nil {
		functionConfig.Runtime = utils.PythonRuntime
	}
	if functionConfig.Jar != nil {
		functionConfig.Runtime = utils.JavaRuntime
	}

	if functionConfig.Runtime == utils.JavaRuntime || functionConfig.Runtime == utils.PythonRuntime {
		if functionConfig.ClassName == "" {
			return fmt.Errorf("no function classname specified")
		}
	}
	return nil
}

func checkArgsForUpdate(functionConfig *utils.FunctionConfig) error {
	if functionConfig.ClassName == "" {
		if functionConfig.Name == "" {
			return fmt.Errorf("function name not provided")
		}
	} else if functionConfig.Name == "" {
		inferMissingFunctionName(functionConfig)
	}

	if functionConfig.Tenant == "" {
		inferMissingTenant(functionConfig)
	}
	if functionConfig.Namespace == "" {
		inferMissingNamespace(functionConfig)
	}
	return nil
}

func inferMissingFunctionName(funcConf *utils.FunctionConfig) {
	className := funcConf.ClassName
	domains := strings.Split(className, ".")
	if len(domains) == 0 {
		funcConf.Name = funcConf.ClassName
		return
	}
	funcConf.Name = domains[len(domains)-1]
}

func inferMissingTenant(funcConf *utils.FunctionConfig) {
	funcConf.Tenant = defaultTenant
}

func inferMissingNamespace(funcConf *utils.FunctionConfig) {
	funcConf.Namespace = defaultNamespace
}

func formatFunctionConfig(funcConf *utils.FunctionConfig) {
	if funcConf == nil {
		return
	}
	for key, value := range funcConf.UserConfig {
		funcConf.UserConfig[key] = convertMap(value)
	}
	for key, value := range funcConf.Secrets {
		funcConf.Secrets[key] = convertMap(value)
	}
}

func resolvePackagePath(config *utils.FunctionConfig) string {
	if config.Jar != nil {
		return *config.Jar
	}
	if config.Py != nil {
		return *config.Py
	}
	if config.Go != nil {
		return *config.Go
	}
	return ""
}

func ensureResources(resources *utils.Resources) *utils.Resources {
	if resources == nil {
		return utils.NewDefaultResources()
	}
	return resources
}

func ensureWindowConfig(conf *utils.WindowConfig) *utils.WindowConfig {
	if conf == nil {
		return utils.NewDefaultWindowConfing()
	}
	return conf
}

func parseStringSlice(value interface{}) ([]string, error) {
	switch v := value.(type) {
	case []string:
		return v, nil
	case []interface{}:
		result := make([]string, 0, len(v))
		for _, item := range v {
			str, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("inputs must be strings")
			}
			if strings.TrimSpace(str) != "" {
				result = append(result, str)
			}
		}
		return result, nil
	case string:
		parts := strings.Split(v, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			item := strings.TrimSpace(part)
			if item != "" {
				result = append(result, item)
			}
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unsupported inputs type")
	}
}

func decodeStringMap(value interface{}) (map[string]string, error) {
	result := map[string]string{}
	if err := decodeInto(value, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func decodeInterfaceMap(value interface{}) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	if err := decodeInto(value, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func decodeConsumerConfigMap(value interface{}) (map[string]utils.ConsumerConfig, error) {
	result := map[string]utils.ConsumerConfig{}
	if err := decodeInto(value, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func decodeInto(value interface{}, out interface{}) error {
	switch v := value.(type) {
	case string:
		if err := json.Unmarshal([]byte(v), out); err != nil {
			return err
		}
		return nil
	default:
		raw, err := json.Marshal(v)
		if err != nil {
			return err
		}
		return json.Unmarshal(raw, out)
	}
}

func parseOptionalIntArg(args map[string]interface{}, key string) (int, bool, error) {
	value, ok := args[key]
	if !ok || value == nil {
		return 0, false, nil
	}
	switch v := value.(type) {
	case int:
		return v, true, nil
	case int32:
		return int(v), true, nil
	case int64:
		return int(v), true, nil
	case float64:
		return int(v), true, nil
	case float32:
		return int(v), true, nil
	case json.Number:
		parsed, err := v.Int64()
		return int(parsed), true, err
	case string:
		parsed, err := strconv.Atoi(v)
		return parsed, true, err
	default:
		return 0, true, fmt.Errorf("invalid value for %s", key)
	}
}

func parseOptionalInt64Arg(args map[string]interface{}, key string) (int64, bool, error) {
	value, ok := args[key]
	if !ok || value == nil {
		return 0, false, nil
	}
	switch v := value.(type) {
	case int:
		return int64(v), true, nil
	case int32:
		return int64(v), true, nil
	case int64:
		return v, true, nil
	case float64:
		return int64(v), true, nil
	case float32:
		return int64(v), true, nil
	case json.Number:
		parsed, err := v.Int64()
		return parsed, true, err
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		return parsed, true, err
	default:
		return 0, true, fmt.Errorf("invalid value for %s", key)
	}
}

func parseOptionalFloatArg(args map[string]interface{}, key string) (float64, bool, error) {
	value, ok := args[key]
	if !ok || value == nil {
		return 0, false, nil
	}
	switch v := value.(type) {
	case float64:
		return v, true, nil
	case float32:
		return float64(v), true, nil
	case int:
		return float64(v), true, nil
	case int64:
		return float64(v), true, nil
	case json.Number:
		parsed, err := v.Float64()
		return parsed, true, err
	case string:
		parsed, err := strconv.ParseFloat(v, 64)
		return parsed, true, err
	default:
		return 0, true, fmt.Errorf("invalid value for %s", key)
	}
}

func getStringArg(args map[string]interface{}, key string) (string, bool) {
	value, ok := args[key]
	if !ok || value == nil {
		return "", false
	}
	switch v := value.(type) {
	case string:
		return v, true
	default:
		return fmt.Sprintf("%v", v), true
	}
}

func getBoolArg(args map[string]interface{}, key string) (bool, bool) {
	value, ok := args[key]
	if !ok || value == nil {
		return false, false
	}
	switch v := value.(type) {
	case bool:
		return v, true
	case string:
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			return false, true
		}
		return parsed, true
	case float64:
		return v != 0, true
	case int:
		return v != 0, true
	default:
		return false, true
	}
}

func numProvidedStrings(values ...*string) int {
	count := 0
	for _, value := range values {
		if value != nil && *value != "" {
			count++
		}
	}
	return count
}

func convertMap(value interface{}) interface{} {
	switch v := value.(type) {
	case map[interface{}]interface{}:
		converted := make(map[string]interface{}, len(v))
		for key, item := range v {
			converted[fmt.Sprintf("%v", key)] = convertMap(item)
		}
		return converted
	case map[string]interface{}:
		converted := make(map[string]interface{}, len(v))
		for key, item := range v {
			converted[key] = convertMap(item)
		}
		return converted
	case []interface{}:
		converted := make([]interface{}, 0, len(v))
		for _, item := range v {
			converted = append(converted, convertMap(item))
		}
		return converted
	default:
		return v
	}
}

func isPackageURLSupported(functionPkgURL string) bool {
	return functionPkgURL != "" && (strings.HasPrefix(functionPkgURL, "http") ||
		strings.HasPrefix(functionPkgURL, "file") ||
		strings.HasPrefix(functionPkgURL, "function") ||
		strings.HasPrefix(functionPkgURL, "sink") ||
		strings.HasPrefix(functionPkgURL, "source"))
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

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
