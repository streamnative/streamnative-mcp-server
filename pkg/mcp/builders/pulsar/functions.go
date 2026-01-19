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

type pulsarAdminFunctionsInput struct {
	Operation    string         `json:"operation"`
	Tenant       string         `json:"tenant"`
	Namespace    string         `json:"namespace"`
	Name         *string        `json:"name,omitempty"`
	ClassName    *string        `json:"classname,omitempty"`
	Inputs       []string       `json:"inputs,omitempty"`
	Output       *string        `json:"output,omitempty"`
	Jar          *string        `json:"jar,omitempty"`
	Py           *string        `json:"py,omitempty"`
	GoFile       *string        `json:"go,omitempty"`
	Parallelism  *int           `json:"parallelism,omitempty"`
	UserConfig   map[string]any `json:"userConfig,omitempty"`
	Key          *string        `json:"key,omitempty"`
	Value        *string        `json:"value,omitempty"`
	Topic        *string        `json:"topic,omitempty"`
	TriggerValue *string        `json:"triggerValue,omitempty"`
}

const (
	pulsarAdminFunctionsOperationDesc = "Operation to perform. Available operations:\n" +
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
	pulsarAdminFunctionsTenantDesc = "The tenant name. Tenants are the primary organizational unit in Pulsar, " +
		"providing multi-tenancy and resource isolation. Functions deployed within a tenant " +
		"inherit its permissions and resource quotas."
	pulsarAdminFunctionsNamespaceDesc = "The namespace name. Namespaces are logical groupings of topics and functions " +
		"within a tenant. They encapsulate configuration policies and access control. " +
		"Functions in a namespace typically process topics within the same namespace."
	pulsarAdminFunctionsNameDesc = "The function name. Required for all operations except 'list'. " +
		"Names should be descriptive of the function's purpose and must be unique within a namespace. " +
		"Function names are used in metrics, logs, and when addressing the function via APIs."
	pulsarAdminFunctionsClassNameDesc = "The fully qualified class name implementing the function. Required for 'create' operation, optional for 'update'. " +
		"For Java functions, this should be the class that implements pulsar function interfaces. " +
		"For Python, this MUST be in format of `<Python_filename_without_extension>.<ClassName>` - for example: " +
		"if file is '/path/to/exclamation.py' with class 'ExclamationFunction', classname must be 'exclamation.ExclamationFunction'; " +
		"if file is '/path/to/double_number.py' with class 'DoubleNumber', classname must be 'double_number.DoubleNumber'. " +
		"Common error: using just the class name 'DoubleNumber' (without filename prefix) will cause function creation to fail. " +
		"Go functions should specify the 'main' function of the binary."
	pulsarAdminFunctionsInputsDesc = "The input topics for the function (array of strings). Optional for 'create' and 'update' operations. " +
		"Topics must be specified in the format 'persistent://tenant/namespace/topic'. " +
		"Functions can consume from multiple topics, each with potentially different serialization types. " +
		"All input topics should exist before the function is created."
	pulsarAdminFunctionsOutputDesc = "The output topic for the function results. Optional for 'create' and 'update' operations. " +
		"Specified in the format 'persistent://tenant/namespace/topic'. " +
		"If not set, the function will not produce any output to topics. " +
		"The output topic will be automatically created if it doesn't exist."
	pulsarAdminFunctionsJarDesc = "Path to the JAR file containing the function code. Optional for 'create' and 'update' operations. " +
		"Support `file://`, `http://`, `https://`, `function://`, `source://`, `sink://` protocol. " +
		"Can be a local path or supported URL protocol accessible to the Pulsar broker. " +
		"For Java functions, this should contain all dependencies for the function. " +
		"The jar file must be compatible with the Pulsar Functions API."
	pulsarAdminFunctionsPyDesc = "Path to the Python file containing the function code. Optional for 'create' and 'update' operations. " +
		"Support `file://`, `http://`, `https://`, `function://`, `source://`, `sink://` protocol. " +
		"Can be a local path or supported URL protocol accessible to the Pulsar broker. " +
		"For Python functions, this should be the file path to the Python file, in format of `.py`, `.zip`, or `.whl`. " +
		"The Python file must be compatible with the Pulsar Functions API."
	pulsarAdminFunctionsGoDesc = "Path to the Go file containing the function code. Optional for 'create' and 'update' operations. " +
		"Support `file://`, `http://`, `https://`, `function://`, `source://`, `sink://` protocol. " +
		"Can be a local path or supported URL protocol accessible to the Pulsar broker. " +
		"For Go functions, this should be the file path to the Go file, in format of executable binary. " +
		"The Go file must be compatible with the Pulsar Functions API."
	pulsarAdminFunctionsParallelismDesc = "The parallelism factor of the function. Optional for 'create' and 'update' operations. " +
		"Determines how many instances of the function will run concurrently. " +
		"Higher values improve throughput but require more resources. " +
		"For stateful functions, consider how parallelism affects state consistency. " +
		"Default is 1 (single instance)."
	pulsarAdminFunctionsUserConfigDesc = "User-defined config key/values. Optional for 'create' and 'update' operations. " +
		"Provides configuration parameters accessible to the function at runtime. " +
		"Specify as a JSON object with string, number, or boolean values. " +
		"Common configs include connection parameters, batch sizes, or feature toggles. " +
		"Example: {\"maxBatchSize\": 100, \"connectionString\": \"host:port\", \"debugMode\": true}"
	pulsarAdminFunctionsKeyDesc = "The state key. Required for 'querystate' and 'putstate' operations. " +
		"Keys are used to identify values in the function's state store. " +
		"They should be reasonable in length and follow a consistent pattern. " +
		"State keys are typically limited to 128 characters."
	pulsarAdminFunctionsValueDesc = "The state value. Required for 'putstate' operation. " +
		"Values are stored in the function's state system. " +
		"For simple values, specify as a string. For complex objects, use JSON-serialized strings. " +
		"State values are typically limited to 1MB in size."
	pulsarAdminFunctionsTopicDesc = "The specific topic name that the function should consume from. Optional for 'trigger' operation. " +
		"Specified in the format 'persistent://tenant/namespace/topic'. " +
		"Used when triggering a function that consumes from multiple topics. " +
		"If not provided, the first input topic will be used."
	pulsarAdminFunctionsTriggerValueDesc = "The value with which to trigger the function. Required for 'trigger' operation. " +
		"This value will be passed to the function as if it were a message from the input topic. " +
		"String values are sent as is; for typed values, ensure proper formatting based on function expectations. " +
		"The function processes this value just like a normal message."
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
func (b *PulsarAdminFunctionsToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool, err := b.buildPulsarAdminFunctionsTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildPulsarAdminFunctionsHandler(config.ReadOnly)

	return []builders.ToolDefinition{
		builders.ServerTool[pulsarAdminFunctionsInput, any]{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildPulsarAdminFunctionsTool builds the Pulsar admin functions MCP tool definition
// Migrated from the original tool definition logic
func (b *PulsarAdminFunctionsToolBuilder) buildPulsarAdminFunctionsTool() (*sdk.Tool, error) {
	inputSchema, err := buildPulsarAdminFunctionsInputSchema()
	if err != nil {
		return nil, err
	}

	toolDesc := "Manage Apache Pulsar Functions for stream processing. " +
		"Pulsar Functions are lightweight compute processes that can consume messages from one or more Pulsar topics, " +
		"apply user-defined processing logic, and produce results to another topic. " +
		"Functions support Java, Python, and Go runtimes, enabling complex event processing, " +
		"data transformations, filtering, and integration with external systems. " +
		"Functions follow the tenant/namespace/name hierarchy for organization, " +
		"can maintain state, and can scale through parallelism configuration. " +
		"This tool provides complete lifecycle management including deployment, monitoring, scaling, " +
		"state management, and triggering. Functions require proper permissions to access their topics."

	return &sdk.Tool{
		Name:        "pulsar_admin_functions",
		Description: toolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildPulsarAdminFunctionsHandler builds the Pulsar admin functions handler function
// Migrated from the original handler logic
func (b *PulsarAdminFunctionsToolBuilder) buildPulsarAdminFunctionsHandler(readOnly bool) builders.ToolHandlerFunc[pulsarAdminFunctionsInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input pulsarAdminFunctionsInput) (*sdk.CallToolResult, any, error) {
		// Extract and validate operation parameter
		operation := input.Operation
		if operation == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'operation'")
		}

		// Check if the operation is valid
		validOperations := map[string]bool{
			"list": true, "get": true, "status": true, "stats": true, "querystate": true,
			"create": true, "update": true, "delete": true, "start": true, "stop": true,
			"restart": true, "putstate": true, "trigger": true,
		}

		if !validOperations[operation] {
			return nil, nil, fmt.Errorf("invalid operation: '%s'. Supported operations: list, get, status, stats, querystate, create, update, delete, start, stop, restart, putstate, trigger", operation)
		}

		// Check write permissions for write operations
		writeOperations := map[string]bool{
			"create": true, "update": true, "delete": true, "start": true,
			"stop": true, "restart": true, "putstate": true, "trigger": true,
		}

		if readOnly && writeOperations[operation] {
			return nil, nil, fmt.Errorf("operation '%s' not allowed in read-only mode. Read-only mode restricts modifications to Pulsar Functions", operation)
		}

		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return nil, nil, fmt.Errorf("pulsar session not found in context")
		}

		client, err := session.GetAdminV3Client()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get pulsar client: %v", err)
		}

		// Extract common parameters
		tenant, err := requireNonEmpty(input.Tenant, "tenant")
		if err != nil {
			return nil, nil, fmt.Errorf("missing required parameter 'tenant': %v. A tenant is required for all Pulsar Functions operations", err)
		}

		namespace, err := requireNonEmpty(input.Namespace, "namespace")
		if err != nil {
			return nil, nil, fmt.Errorf("missing required parameter 'namespace': %v. A namespace is required for all Pulsar Functions operations", err)
		}

		// For all operations except 'list', name is required
		var name string
		if operation != "list" {
			name, err = requireString(input.Name, "name")
			if err != nil {
				return nil, nil, fmt.Errorf("missing required parameter 'name' for operation '%s': %v. The function name must be specified for this operation", operation, err)
			}
		}

		// Handle operation using delegated handlers
		switch operation {
		case "list":
			result, err := b.handleFunctionList(ctx, client, tenant, namespace)
			return result, nil, err
		case "get":
			result, err := b.handleFunctionGet(ctx, client, tenant, namespace, name)
			return result, nil, err
		case "status":
			result, err := b.handleFunctionStatus(ctx, client, tenant, namespace, name)
			return result, nil, err
		case "stats":
			result, err := b.handleFunctionStats(ctx, client, tenant, namespace, name)
			return result, nil, err
		case "querystate":
			key, err := requireString(input.Key, "key")
			if err != nil {
				return nil, nil, fmt.Errorf("missing required parameter 'key' for operation 'querystate': %v. A key is required to look up state in the function's state store", err)
			}
			result, err := b.handleFunctionQuerystate(ctx, client, tenant, namespace, name, key)
			return result, nil, err
		case "create":
			result, err := b.handleFunctionCreate(ctx, client, tenant, namespace, name, input)
			return result, nil, err
		case "update":
			result, err := b.handleFunctionUpdate(ctx, client, tenant, namespace, name, input)
			return result, nil, err
		case "delete":
			result, err := b.handleFunctionDelete(ctx, client, tenant, namespace, name)
			return result, nil, err
		case "start":
			result, err := b.handleFunctionStart(ctx, client, tenant, namespace, name)
			return result, nil, err
		case "stop":
			result, err := b.handleFunctionStop(ctx, client, tenant, namespace, name)
			return result, nil, err
		case "restart":
			result, err := b.handleFunctionRestart(ctx, client, tenant, namespace, name)
			return result, nil, err
		case "putstate":
			key, err := requireString(input.Key, "key")
			if err != nil {
				return nil, nil, fmt.Errorf("missing required parameter 'key' for operation 'putstate': %v. A key is required to store state in the function's state store", err)
			}
			value, err := requireString(input.Value, "value")
			if err != nil {
				return nil, nil, fmt.Errorf("missing required parameter 'value' for operation 'putstate': %v. A value is required to store state in the function's state store", err)
			}
			result, err := b.handleFunctionPutstate(ctx, client, tenant, namespace, name, key, value)
			return result, nil, err
		case "trigger":
			triggerValue, err := requireString(input.TriggerValue, "triggerValue")
			if err != nil {
				return nil, nil, fmt.Errorf("missing required parameter 'triggerValue' for operation 'trigger': %v. A trigger value is required to manually trigger the function", err)
			}
			topic := ""
			if input.Topic != nil {
				topic = *input.Topic
			}
			result, err := b.handleFunctionTrigger(ctx, client, tenant, namespace, name, triggerValue, topic)
			return result, nil, err
		default:
			return nil, nil, fmt.Errorf("unsupported operation: %s", operation)
		}
	}
}

// Helper functions - delegated operation handlers

// handleFunctionList handles the list operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionList(_ context.Context, client cmdutils.Client, tenant, namespace string) (*sdk.CallToolResult, error) {
	admin := client.Functions()

	functions, err := admin.GetFunctions(tenant, namespace)
	if err != nil {
		return nil, b.handleError("list functions", err)
	}

	return b.marshalResponse(map[string]interface{}{
		"functions": functions,
		"tenant":    tenant,
		"namespace": namespace,
	})
}

// handleFunctionGet handles the get operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionGet(_ context.Context, client cmdutils.Client, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	admin := client.Functions()

	functionConfig, err := admin.GetFunction(tenant, namespace, name)
	if err != nil {
		return nil, b.handleError("get function config", err)
	}

	return b.marshalResponse(functionConfig)
}

// handleFunctionStatus handles the status operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionStatus(_ context.Context, client cmdutils.Client, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	admin := client.Functions()

	status, err := admin.GetFunctionStatus(tenant, namespace, name)
	if err != nil {
		return nil, b.handleError("get function status", err)
	}

	return b.marshalResponse(status)
}

// handleFunctionStats handles the stats operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionStats(_ context.Context, client cmdutils.Client, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	admin := client.Functions()

	stats, err := admin.GetFunctionStats(tenant, namespace, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats for function '%s' in tenant '%s' namespace '%s': %v; verify the function exists and is running",
			name, tenant, namespace, err)
	}

	return b.marshalResponse(stats)
}

// handleFunctionQuerystate handles the querystate operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionQuerystate(_ context.Context, client cmdutils.Client, tenant, namespace, name, key string) (*sdk.CallToolResult, error) {
	admin := client.Functions()

	state, err := admin.GetFunctionState(tenant, namespace, name, key)
	if err != nil {
		return nil, fmt.Errorf("failed to query state for key '%s' in function '%s' (tenant '%s' namespace '%s'): %v; verify the function exists and has state enabled",
			key, name, tenant, namespace, err)
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
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionCreate(_ context.Context, client cmdutils.Client, tenant, namespace, name string, input pulsarAdminFunctionsInput) (*sdk.CallToolResult, error) {
	// Build function configuration from request parameters to validate
	functionConfig, err := b.buildFunctionConfig(tenant, namespace, name, input, false)
	if err != nil {
		return nil, fmt.Errorf("failed to build function configuration for '%s' in tenant '%s' namespace '%s': %v; please verify all required parameters are provided correctly",
			name, tenant, namespace, err)
	}

	admin := client.Functions()
	packagePath := ""
	//nolint:gocritic
	if functionConfig.Jar != nil {
		packagePath = *functionConfig.Jar
	} else if functionConfig.Py != nil {
		packagePath = *functionConfig.Py
	} else if functionConfig.Go != nil {
		packagePath = *functionConfig.Go
	}

	err = admin.CreateFuncWithURL(functionConfig, packagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create function '%s' in tenant '%s' namespace '%s': %v; verify the function configuration is valid",
			name, tenant, namespace, err)
	}

	return textResult(fmt.Sprintf("Created function '%s' successfully in tenant '%s' namespace '%s'. The function configuration has been created.",
		name, tenant, namespace)), nil
}

// handleFunctionUpdate handles the update operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionUpdate(_ context.Context, client cmdutils.Client, tenant, namespace, name string, input pulsarAdminFunctionsInput) (*sdk.CallToolResult, error) {
	admin := client.Functions()

	// Build function configuration from request parameters
	config, err := b.buildFunctionConfig(tenant, namespace, name, input, true)
	if err != nil {
		return nil, fmt.Errorf("failed to build function configuration for '%s' in tenant '%s' namespace '%s': %v; please verify all parameters are provided correctly",
			name, tenant, namespace, err)
	}

	// Update the function
	updateOptions := &utils.UpdateOptions{
		UpdateAuthData: true,
	}
	err = admin.UpdateFunction(config, "", updateOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to update function '%s' in tenant '%s' namespace '%s': %v; verify the function exists and the configuration is valid",
			name, tenant, namespace, err)
	}

	return textResult(fmt.Sprintf("Updated function '%s' successfully in tenant '%s' namespace '%s'. The function configuration has been modified.",
		name, tenant, namespace)), nil
}

// handleFunctionDelete handles the delete operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionDelete(_ context.Context, client cmdutils.Client, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	admin := client.Functions()

	err := admin.DeleteFunction(tenant, namespace, name)
	if err != nil {
		return nil, fmt.Errorf("failed to delete function '%s' in tenant '%s' namespace '%s': %v; verify the function exists and you have deletion permissions",
			name, tenant, namespace, err)
	}

	return textResult(fmt.Sprintf("Deleted function '%s' successfully from tenant '%s' namespace '%s'. All running instances have been terminated.",
		name, tenant, namespace)), nil
}

// handleFunctionStart handles the start operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionStart(_ context.Context, client cmdutils.Client, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	admin := client.Functions()

	err := admin.StartFunction(tenant, namespace, name)
	if err != nil {
		return nil, fmt.Errorf("failed to start function '%s' in tenant '%s' namespace '%s': %v; verify the function exists and is not already running",
			name, tenant, namespace, err)
	}

	return textResult(fmt.Sprintf("Started function '%s' successfully in tenant '%s' namespace '%s'. The function instances are now processing messages.",
		name, tenant, namespace)), nil
}

// handleFunctionStop handles the stop operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionStop(_ context.Context, client cmdutils.Client, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	admin := client.Functions()

	err := admin.StopFunction(tenant, namespace, name)
	if err != nil {
		return nil, fmt.Errorf("failed to stop function '%s' in tenant '%s' namespace '%s': %v; verify the function exists and is currently running",
			name, tenant, namespace, err)
	}

	return textResult(fmt.Sprintf("Stopped function '%s' successfully in tenant '%s' namespace '%s'. The function will no longer process messages until restarted.",
		name, tenant, namespace)), nil
}

// handleFunctionRestart handles the restart operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionRestart(_ context.Context, client cmdutils.Client, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	admin := client.Functions()

	err := admin.RestartFunction(tenant, namespace, name)
	if err != nil {
		return nil, fmt.Errorf("failed to restart function '%s' in tenant '%s' namespace '%s': %v; verify the function exists and is properly deployed",
			name, tenant, namespace, err)
	}

	return textResult(fmt.Sprintf("Restarted function '%s' successfully in tenant '%s' namespace '%s'. All function instances have been restarted.",
		name, tenant, namespace)), nil
}

// handleFunctionPutstate handles the putstate operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionPutstate(_ context.Context, client cmdutils.Client, tenant, namespace, name, key, value string) (*sdk.CallToolResult, error) {
	admin := client.Functions()

	err := admin.PutFunctionState(tenant, namespace, name, utils.FunctionState{
		Key:         key,
		StringValue: value,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to put state for key '%s' in function '%s' (tenant '%s' namespace '%s'): %v; verify the function exists and has state enabled",
			key, name, tenant, namespace, err)
	}

	return textResult(fmt.Sprintf("Successfully stored state for key '%s' in function '%s' (tenant '%s' namespace '%s'). State value has been updated.",
		key, name, tenant, namespace)), nil
}

// handleFunctionTrigger handles the trigger operation
func (b *PulsarAdminFunctionsToolBuilder) handleFunctionTrigger(_ context.Context, client cmdutils.Client, tenant, namespace, name, triggerValue, topic string) (*sdk.CallToolResult, error) {
	admin := client.Functions()

	var err error
	var result string
	if topic != "" {
		// Trigger with specific topic
		result, err = admin.TriggerFunction(tenant, namespace, name, topic, triggerValue, "")
	} else {
		// Trigger without specific topic (uses first input topic)
		result, err = admin.TriggerFunction(tenant, namespace, name, "", triggerValue, "")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to trigger function '%s' in tenant '%s' namespace '%s': %v; verify the function exists and is running",
			name, tenant, namespace, err)
	}

	var message string
	if topic != "" {
		message = fmt.Sprintf("Successfully triggered function '%s' in tenant '%s' namespace '%s' with topic '%s'. Result: %s",
			name, tenant, namespace, topic, result)
	} else {
		message = fmt.Sprintf("Successfully triggered function '%s' in tenant '%s' namespace '%s'. Result: %s",
			name, tenant, namespace, result)
	}

	return textResult(message), nil
}

// Helper functions

// buildFunctionConfig builds a Pulsar Function configuration from MCP request parameters
func (b *PulsarAdminFunctionsToolBuilder) buildFunctionConfig(tenant, namespace, name string, input pulsarAdminFunctionsInput, isUpdate bool) (*utils.FunctionConfig, error) {
	config := &utils.FunctionConfig{
		Tenant:    tenant,
		Namespace: namespace,
		Name:      name,
	}

	// Get required classname parameter (for create operations)
	if !isUpdate {
		classname, err := requireString(input.ClassName, "classname")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'classname': %v", err)
		}
		config.ClassName = classname
	} else if input.ClassName != nil && *input.ClassName != "" {
		// For update, classname is optional
		config.ClassName = *input.ClassName
	}

	// Get inputs parameter (array of strings)
	if len(input.Inputs) > 0 {
		inputSpecs := make(map[string]utils.ConsumerConfig)
		for _, inputTopic := range input.Inputs {
			inputSpecs[inputTopic] = utils.ConsumerConfig{
				SerdeClassName: "",
				SchemaType:     "",
			}
		}
		if len(inputSpecs) > 0 {
			config.InputSpecs = inputSpecs
		}
	}

	// Get optional output parameter
	if input.Output != nil && *input.Output != "" {
		config.Output = *input.Output
	}

	// Get optional parallelism parameter
	if input.Parallelism != nil {
		config.Parallelism = *input.Parallelism
	}

	// Set default parallelism if not specified
	if config.Parallelism <= 0 {
		config.Parallelism = 1
	}

	// Get optional jar parameter
	if input.Jar != nil && *input.Jar != "" {
		jar := *input.Jar
		config.Jar = &jar
	}

	// Get optional py parameter
	if input.Py != nil && *input.Py != "" {
		py := *input.Py
		config.Py = &py
	}

	// Get optional go parameter
	if input.GoFile != nil && *input.GoFile != "" {
		goFile := *input.GoFile
		config.Go = &goFile
	}

	// Get optional userConfig parameter (JSON object)
	if input.UserConfig != nil {
		config.UserConfig = input.UserConfig
	}

	return config, nil
}

// handleError provides unified error handling
func (b *PulsarAdminFunctionsToolBuilder) handleError(operation string, err error) error {
	return fmt.Errorf("failed to %s: %v", operation, err)
}

// marshalResponse provides unified JSON serialization for responses
func (b *PulsarAdminFunctionsToolBuilder) marshalResponse(data interface{}) (*sdk.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, b.handleError("marshal response", err)
	}
	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: string(jsonBytes)}},
	}, nil
}

func requireNonEmpty(value string, key string) (string, error) {
	if value == "" {
		return "", fmt.Errorf("required argument %q not found", key)
	}
	return value, nil
}

func buildPulsarAdminFunctionsInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[pulsarAdminFunctionsInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}

	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}
	setSchemaDescription(schema, "operation", pulsarAdminFunctionsOperationDesc)
	setSchemaDescription(schema, "tenant", pulsarAdminFunctionsTenantDesc)
	setSchemaDescription(schema, "namespace", pulsarAdminFunctionsNamespaceDesc)
	setSchemaDescription(schema, "name", pulsarAdminFunctionsNameDesc)
	setSchemaDescription(schema, "classname", pulsarAdminFunctionsClassNameDesc)
	setSchemaDescription(schema, "inputs", pulsarAdminFunctionsInputsDesc)
	setSchemaDescription(schema, "output", pulsarAdminFunctionsOutputDesc)
	setSchemaDescription(schema, "jar", pulsarAdminFunctionsJarDesc)
	setSchemaDescription(schema, "py", pulsarAdminFunctionsPyDesc)
	setSchemaDescription(schema, "go", pulsarAdminFunctionsGoDesc)
	setSchemaDescription(schema, "parallelism", pulsarAdminFunctionsParallelismDesc)
	setSchemaDescription(schema, "userConfig", pulsarAdminFunctionsUserConfigDesc)
	setSchemaDescription(schema, "key", pulsarAdminFunctionsKeyDesc)
	setSchemaDescription(schema, "value", pulsarAdminFunctionsValueDesc)
	setSchemaDescription(schema, "topic", pulsarAdminFunctionsTopicDesc)
	setSchemaDescription(schema, "triggerValue", pulsarAdminFunctionsTriggerValueDesc)

	normalizeAdditionalProperties(schema)
	return schema, nil
}
