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
	"sort"
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	"gopkg.in/yaml.v2"
)

var pulsarSinkOperationSpecs = builders.OperationRegistry{
	{Name: "list", Mode: builders.OperationModeRead},
	{Name: "get", Mode: builders.OperationModeRead},
	{Name: "status", Mode: builders.OperationModeRead},
	{Name: "list-built-in", Mode: builders.OperationModeRead},
	{Name: "create", Mode: builders.OperationModeWrite, Destructive: true},
	{Name: "update", Mode: builders.OperationModeWrite, Destructive: true},
	{Name: "delete", Mode: builders.OperationModeWrite, Destructive: true},
	{Name: "start", Mode: builders.OperationModeWrite, Destructive: true},
	{Name: "stop", Mode: builders.OperationModeWrite, Destructive: true},
	{Name: "restart", Mode: builders.OperationModeWrite, Destructive: true},
}

// PulsarAdminSinksToolBuilder implements the ToolBuilder interface for Pulsar admin sinks
// /nolint:revive
type PulsarAdminSinksToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminSinksToolBuilder creates a new Pulsar admin sinks tool builder instance
func NewPulsarAdminSinksToolBuilder() *PulsarAdminSinksToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_sinks",
		Version:     "1.0.0",
		Description: "Pulsar admin sinks management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "admin", "sinks"},
	}

	features := []string{
		"pulsar-admin-sinks",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminSinksToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin sinks tool list
func (b *PulsarAdminSinksToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
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
			Tool:    b.buildSinksTool(toolModeRead),
			Handler: b.buildSinksHandler(toolModeRead),
		},
	}
	if !config.ReadOnly {
		tools = append(tools, server.ServerTool{
			Tool:    b.buildSinksTool(toolModeWrite),
			Handler: b.buildSinksHandler(toolModeWrite),
		})
	}

	return tools, nil
}

// buildSinksTool builds the Pulsar admin sinks MCP tool definition
func (b *PulsarAdminSinksToolBuilder) buildSinksTool(mode toolMode) mcp.Tool {
	toolDesc := "Read Apache Pulsar Sinks for data movement and integration. " +
		"Pulsar Sinks are connectors that export data from Pulsar topics to external systems. " +
		"This read-only tool lists sinks and built-in sink connectors and retrieves sink configuration or runtime status."

	operationDesc := "Operation to perform. Available operations:\n" +
		"- list: List all sinks under a specific tenant and namespace\n" +
		"- get: Get the configuration of a sink\n" +
		"- status: Get the runtime status of a sink (instances, metrics)\n" +
		"- list-built-in: List all built-in sink connectors available in the system"

	operationEnum := pulsarSinkOperationSpecs.NamesForMode(mode)
	toolName := "pulsar_admin_sinks_read"
	annotation := builders.ToolAnnotationForMode(mode, "Read Pulsar Sinks", "Manage Pulsar Sinks", pulsarSinkOperationSpecs)
	if isToolModeWrite(mode) {
		toolDesc = "Manage Apache Pulsar Sinks for data movement and integration. " +
			"This write tool deploys, updates, deletes, starts, stops, or restarts sinks."
		operationDesc = "Operation to perform. Available operations:\n" +
			"- create: Deploy a new sink with specified parameters\n" +
			"- update: Update the configuration of an existing sink\n" +
			"- delete: Delete a sink\n" +
			"- start: Start a stopped sink\n" +
			"- stop: Stop a running sink\n" +
			"- restart: Restart a sink"
		operationEnum = pulsarSinkOperationSpecs.NamesForMode(mode)
		toolName = "pulsar_admin_sinks_write"
	}

	tool := mcp.NewTool(toolName,
		mcp.WithDescription(toolDesc),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
			mcp.Enum(operationEnum...)),
		mcp.WithString("tenant",
			mcp.Description("The tenant name. Tenants are the primary organizational unit in Pulsar, "+
				"providing multi-tenancy and resource isolation. Sinks deployed within a tenant "+
				"inherit its permissions and resource quotas. "+
				"Defaults to 'public' if not provided.")),
		mcp.WithString("namespace",
			mcp.Description("The namespace name. Namespaces are logical groupings of topics and sinks "+
				"within a tenant. They encapsulate configuration policies and access control. "+
				"Sinks in a namespace typically process topics within the same namespace. "+
				"Defaults to 'default' if not provided.")),
		mcp.WithString("name",
			mcp.Description("The sink name. Required for operations that target one sink. "+
				"Names should be descriptive of the sink's purpose and must be unique within a namespace. "+
				"Sink names are used in metrics, logs, and when addressing the sink via APIs.")),
		mcp.WithString("archive",
			mcp.Description("Path to the archive file containing the sink code. Optional for 'create' and 'update' operations. "+
				"Can be a local path, NAR file, or a URL accessible to the Pulsar broker. "+
				"The archive should contain all dependencies for the sink connector. "+
				"Either archive or sink-type must be specified, but not both.")),
		mcp.WithString("sink-type",
			mcp.Description("The built-in sink connector type to use. Optional for 'create' and 'update' operations. "+
				"Specifies which built-in connector to use, such as 'jdbc', 'elastic-search', 'kafka', etc. "+
				"Use 'list-built-in' operation to see available sink types. "+
				"Either sink-type or archive must be specified, but not both.")),
		mcp.WithArray("inputs",
			mcp.Description("The sink's input topics (array of strings). Optional for 'create' and 'update' operations. "+
				"Topics must be specified in the format 'persistent://tenant/namespace/topic'. "+
				"Sinks can consume from multiple topics, but they should have compatible schemas. "+
				"All input topics should exist before the sink is created. "+
				"Either inputs or topics-pattern must be specified."),
			mcp.Items(
				map[string]interface{}{
					"type":        "string",
					"description": "input topic",
				},
			),
		),
		mcp.WithString("topics-pattern",
			mcp.Description("TopicsPattern to consume from list of topics that match the pattern. Optional for 'create' and 'update' operations. "+
				"Specified as a regular expression, e.g., 'persistent://tenant/namespace/prefix.*'. "+
				"This allows the sink to automatically consume from topics that match the pattern, "+
				"including topics created after the sink is deployed. "+
				"Either topics-pattern or inputs must be specified.")),
		mcp.WithString("subs-name",
			mcp.Description("Pulsar subscription name for input topic consumer. Optional for 'create' and 'update' operations. "+
				"Defines the subscription name used by the sink to consume from input topics. "+
				"If not specified, a default subscription name will be generated. "+
				"The subscription type used is Shared by default.")),
		mcp.WithString("subs-position",
			mcp.Description("Pulsar source subscription position for input-topic consumers. Optional for 'create' and 'update' operations. "+
				"Possible values: 'Latest' or 'Earliest'. If not specified, the sink uses the default subscription position.")),
		mcp.WithString("classname",
			mcp.Description("The sink class name when using a custom archive. Optional for 'create' and 'update'. "+
				"Specifies the fully qualified class name implementing the sink connector.")),
		mcp.WithString("processing-guarantees",
			mcp.Description("The processing guarantees (delivery semantics) for the sink. Optional for 'create' and 'update'. "+
				"Available options: 'atleast_once', 'atmost_once', 'effectively_once'.")),
		mcp.WithBoolean("retain-ordering",
			mcp.Description("Whether the sink preserves message ordering. Optional for 'create' and 'update'.")),
		mcp.WithBoolean("retain-key-ordering",
			mcp.Description("Whether the sink preserves key ordering. Optional for 'create' and 'update'.")),
		mcp.WithBoolean("auto-ack",
			mcp.Description("Whether the sink automatically acknowledges messages. Optional for 'create' and 'update'.")),
		mcp.WithBoolean("cleanup-subscription",
			mcp.Description("Whether the subscription used by the sink should be deleted when the sink is deleted. Optional for 'create' and 'update'. Default: true.")),
		mcp.WithNumber("parallelism",
			mcp.Description("The parallelism factor of the sink. Optional for 'create' and 'update' operations. "+
				"Determines how many instances of the sink will run concurrently. "+
				"Higher values improve throughput but require more resources. "+
				"Default is 1 (single instance). Recommended to align with topic partition count "+
				"when consuming from partitioned topics.")),
		mcp.WithNumber("cpu",
			mcp.Description("CPU cores allocated per sink instance. Optional for 'create' and 'update'. "+
				"Applicable to process and container runtimes.")),
		mcp.WithNumber("ram",
			mcp.Description("RAM bytes allocated per sink instance. Optional for 'create' and 'update'. "+
				"Applicable to process and container runtimes.")),
		mcp.WithNumber("disk",
			mcp.Description("Disk bytes allocated per sink instance. Optional for 'create' and 'update'. "+
				"Applicable to process and container runtimes.")),
		mcp.WithObject("custom-serde-inputs",
			mcp.Description("Map of input topics to SerDe class names. Optional for 'create' and 'update'. "+
				"Specify as a JSON object of topic to SerDe class.")),
		mcp.WithObject("custom-schema-inputs",
			mcp.Description("Map of input topics to Schema types or class names. Optional for 'create' and 'update'. "+
				"Specify as a JSON object of topic to schema type/class.")),
		mcp.WithObject("input-specs",
			mcp.Description("Map of input topics to consumer configuration. Optional for 'create' and 'update'. "+
				"Specify as a JSON object of topic to consumer config.")),
		mcp.WithNumber("max-redeliver-count",
			mcp.Description("Maximum number of retries before sending a message to the dead letter topic. Optional for 'create' and 'update'.")),
		mcp.WithString("dead-letter-topic",
			mcp.Description("Dead letter topic for messages that exceed retry attempts. Optional for 'create' and 'update'.")),
		mcp.WithNumber("timeout-ms",
			mcp.Description("Message processing timeout in milliseconds. Optional for 'create' and 'update'.")),
		mcp.WithNumber("negative-ack-redelivery-delay-ms",
			mcp.Description("Delay in milliseconds before redelivering negatively acknowledged messages. Optional for 'create' and 'update'.")),
		mcp.WithString("custom-runtime-options",
			mcp.Description("Runtime customization options. Optional for 'create' and 'update'.")),
		mcp.WithObject("secrets",
			mcp.Description("Secrets configuration map. Optional for 'create' and 'update'. "+
				"Specify as a JSON object where values describe how secrets are fetched.")),
		mcp.WithString("sink-config-file",
			mcp.Description("Path to a YAML sink configuration file. Optional for 'create' and 'update'. "+
				"When provided, the file is loaded before applying explicit parameters.")),
		mcp.WithObject("sink-config",
			mcp.Description("User-defined sink config key/values. Optional for 'create' and 'update' operations. "+
				"Provides configuration parameters specific to the sink connector being used. "+
				"For example, JDBC connection strings, Elasticsearch indices, S3 bucket details, etc. "+
				"Specify as a JSON object with configuration properties required by the specific sink type. "+
				"Example: {\"jdbcUrl\": \"jdbc:postgresql://localhost:5432/database\", \"tableName\": \"events\"}")),
		mcp.WithString("transform-function",
			mcp.Description("Transform function applied before the sink. Optional for 'create' and 'update'.")),
		mcp.WithString("transform-function-classname",
			mcp.Description("Transform function class name. Optional for 'create' and 'update'.")),
		mcp.WithString("transform-function-config",
			mcp.Description("Transform function configuration. Optional for 'create' and 'update'.")),
		mcp.WithBoolean("update-auth-data",
			mcp.Description("Whether to update authentication data during sink update. Optional for 'update' only.")),
		annotation,
	)
	if !isToolModeWrite(mode) {
		pruneToolInputSchema(&tool, []string{"operation", "tenant", "namespace", "name"})
	}
	return tool
}

// buildSinksHandler builds the Pulsar admin sinks handler function
func (b *PulsarAdminSinksToolBuilder) buildSinksHandler(mode toolMode) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract and validate operation parameter
		operation, err := request.RequireString("operation")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'operation': %v", err)), nil
		}

		// Check if the operation is valid
		validOperations := map[string]bool{
			"list": true, "get": true, "status": true, "create": true, "update": true,
			"delete": true, "start": true, "stop": true, "restart": true, "list-built-in": true,
		}

		if !validOperations[operation] {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid operation: '%s'. Supported operations: %s", operation,
				modeSupportedOperations(mode, pulsarSinkOperationSpecs))), nil
		}

		if err := validateModeOperation(mode, operation, pulsarSinkOperationSpecs); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return mcp.NewToolResultError("Pulsar session not found in context"), nil
		}

		admin, err := session.GetAdminV3Client()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get Pulsar client: %v", err)), nil
		}

		// List built-in sinks doesn't require tenant, namespace or name
		if operation == "list-built-in" {
			return b.handleListBuiltInSinks(ctx, admin)
		}
		args := request.GetArguments()

		if operation == "create" || operation == "update" {
			name, _ := getStringArg(args, "name")
			tenant, _ := getStringArg(args, "tenant")
			namespace, _ := getStringArg(args, "namespace")
			if operation == "create" {
				return b.handleSinkCreate(ctx, admin, tenant, namespace, name, request)
			}
			return b.handleSinkUpdate(ctx, admin, tenant, namespace, name, request)
		}

		tenant, _ := getStringArg(args, "tenant")
		namespace, _ := getStringArg(args, "namespace")
		if tenant == "" {
			tenant = defaultTenant
		}
		if namespace == "" {
			namespace = defaultNamespace
		}

		var name string
		if operation != "list" {
			name, err = request.RequireString("name")
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'name' for operation '%s': %v. The sink name must be specified for this operation.", operation, err)), nil
			}
		}

		// Handle operations
		switch operation {
		case "list":
			return b.handleSinkList(ctx, admin, tenant, namespace)
		case "get":
			return b.handleSinkGet(ctx, admin, tenant, namespace, name)
		case "status":
			return b.handleSinkStatus(ctx, admin, tenant, namespace, name)
		case "delete":
			return b.handleSinkDelete(ctx, admin, tenant, namespace, name)
		case "start":
			return b.handleSinkStart(ctx, admin, tenant, namespace, name)
		case "stop":
			return b.handleSinkStop(ctx, admin, tenant, namespace, name)
		case "restart":
			return b.handleSinkRestart(ctx, admin, tenant, namespace, name)
		default:
			// This should never happen due to the valid operations check above
			return mcp.NewToolResultError(fmt.Sprintf("Unsupported operation: %s", operation)), nil
		}
	}
}

// Helper functions

// handleSinkList handles listing all sinks under a namespace
func (b *PulsarAdminSinksToolBuilder) handleSinkList(_ context.Context, admin cmdutils.Client, tenant, namespace string) (*mcp.CallToolResult, error) {
	sinks, err := admin.Sinks().ListSinks(tenant, namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list sinks in tenant '%s' namespace '%s': %v. Check that the tenant and namespace exist and you have proper permissions.",
			tenant, namespace, err)), nil
	}

	// Convert result to JSON string
	sinksJSON, err := json.Marshal(sinks)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize sink list: %v", err)), nil
	}

	return mcp.NewToolResultText(string(sinksJSON)), nil
}

// handleSinkGet handles getting information about a sink
func (b *PulsarAdminSinksToolBuilder) handleSinkGet(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	sink, err := admin.Sinks().GetSink(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get sink '%s' in tenant '%s' namespace '%s': %v. Verify the sink exists and you have proper permissions.",
			name, tenant, namespace, err)), nil
	}

	// Convert result to JSON string
	sinkJSON, err := json.Marshal(sink)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize sink info: %v", err)), nil
	}

	return mcp.NewToolResultText(string(sinkJSON)), nil
}

// handleSinkStatus handles getting the status of a sink
func (b *PulsarAdminSinksToolBuilder) handleSinkStatus(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	status, err := admin.Sinks().GetSinkStatus(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get status for sink '%s' in tenant '%s' namespace '%s': %v. Verify the sink exists and is properly deployed.",
			name, tenant, namespace, err)), nil
	}

	// Convert result to JSON string
	statusJSON, err := json.Marshal(status)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize sink status: %v", err)), nil
	}

	return mcp.NewToolResultText(string(statusJSON)), nil
}

// handleSinkCreate handles creating a new sink
func (b *PulsarAdminSinksToolBuilder) handleSinkCreate(_ context.Context, admin cmdutils.Client, tenant, namespace, name string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	config, archiveArg, sinkTypeArg, err := b.buildSinkConfig(tenant, namespace, name, request)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build sink configuration for '%s': %v. Please verify all required parameters are provided correctly.", name, err)), nil
	}

	if err := validateSinkArchiveArgs(archiveArg, sinkTypeArg); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid parameters: %v", err)), nil
	}

	uploadArchive, err := b.resolveSinkArchive(admin, config, archiveArg, sinkTypeArg)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to resolve sink archive: %v", err)), nil
	}

	b.applySinkDefaults(config)
	if err := validateSinkConfig(config); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to validate sink configuration for '%s': %v.", config.Name, err)), nil
	}

	if len(config.Inputs) == 0 && (config.TopicsPattern == nil || strings.TrimSpace(*config.TopicsPattern) == "") {
		return mcp.NewToolResultError("Missing required parameter: Either 'inputs' or 'topics-pattern' must be specified. The sink needs a source of data to consume from Pulsar."), nil
	}

	if uploadArchive != "" && b.isPackageURLSupported(uploadArchive) {
		err = admin.Sinks().CreateSinkWithURL(config, uploadArchive)
	} else {
		err = admin.Sinks().CreateSink(config, uploadArchive)
	}

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create sink '%s' in tenant '%s' namespace '%s': %v. Verify all parameters are correct and required resources exist.",
			config.Name, config.Tenant, config.Namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Created sink '%s' successfully in tenant '%s' namespace '%s'. The sink will start consuming from its input topics and writing to the configured destination.",
		config.Name, config.Tenant, config.Namespace)), nil
}

// handleSinkUpdate handles updating an existing sink
func (b *PulsarAdminSinksToolBuilder) handleSinkUpdate(_ context.Context, admin cmdutils.Client, tenant, namespace, name string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	config, archiveArg, sinkTypeArg, err := b.buildSinkConfig(tenant, namespace, name, request)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build sink configuration for '%s': %v. Please verify all parameters are provided correctly.", name, err)), nil
	}

	if err := validateSinkArchiveArgs(archiveArg, sinkTypeArg); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid parameters: %v", err)), nil
	}

	uploadArchive, err := b.resolveSinkArchive(admin, config, archiveArg, sinkTypeArg)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to resolve sink archive: %v", err)), nil
	}

	b.applySinkDefaults(config)

	if config.Name == "" {
		return mcp.NewToolResultError("Sink name not specified. Provide the 'name' parameter or set it in the sink-config-file."), nil
	}

	latestConfig, err := admin.Sinks().GetSink(config.Tenant, config.Namespace, config.Name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get sink '%s' in tenant '%s' namespace '%s': %v. Verify the sink exists and you have proper permissions.",
			config.Name, config.Tenant, config.Namespace, err)), nil
	}

	config.AutoAck = latestConfig.AutoAck
	config.RetainOrdering = latestConfig.RetainOrdering
	config.ProcessingGuarantees = latestConfig.ProcessingGuarantees

	updateOptions := utils.NewUpdateOptions()
	if updateAuthData, ok := getBoolArg(request.GetArguments(), "update-auth-data"); ok {
		updateOptions.UpdateAuthData = updateAuthData
	}

	if uploadArchive != "" && b.isPackageURLSupported(uploadArchive) {
		err = admin.Sinks().UpdateSinkWithURL(config, uploadArchive, updateOptions)
	} else {
		err = admin.Sinks().UpdateSink(config, uploadArchive, updateOptions)
	}

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to update sink '%s' in tenant '%s' namespace '%s': %v. Verify the sink exists and all parameters are valid.",
			config.Name, config.Tenant, config.Namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Updated sink '%s' successfully in tenant '%s' namespace '%s'. The sink may need to be restarted to apply all changes.",
		config.Name, config.Tenant, config.Namespace)), nil
}

// handleSinkDelete handles deleting a sink
func (b *PulsarAdminSinksToolBuilder) handleSinkDelete(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	err := admin.Sinks().DeleteSink(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete sink '%s' in tenant '%s' namespace '%s': %v. Verify the sink exists and you have deletion permissions.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Deleted sink '%s' successfully from tenant '%s' namespace '%s'. All running instances have been terminated.",
		name, tenant, namespace)), nil
}

// handleSinkStart handles starting a sink
func (b *PulsarAdminSinksToolBuilder) handleSinkStart(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	err := admin.Sinks().StartSink(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to start sink '%s' in tenant '%s' namespace '%s': %v. Verify the sink exists and is not already running.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Started sink '%s' successfully in tenant '%s' namespace '%s'. The sink will begin consuming from its input topics.",
		name, tenant, namespace)), nil
}

// handleSinkStop handles stopping a sink
func (b *PulsarAdminSinksToolBuilder) handleSinkStop(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	err := admin.Sinks().StopSink(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to stop sink '%s' in tenant '%s' namespace '%s': %v. Verify the sink exists and is currently running.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Stopped sink '%s' successfully in tenant '%s' namespace '%s'. The sink will no longer consume messages until restarted.",
		name, tenant, namespace)), nil
}

// handleSinkRestart handles restarting a sink
func (b *PulsarAdminSinksToolBuilder) handleSinkRestart(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	err := admin.Sinks().RestartSink(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to restart sink '%s' in tenant '%s' namespace '%s': %v. Verify the sink exists and is properly deployed.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Restarted sink '%s' successfully in tenant '%s' namespace '%s'. All sink instances have been restarted.",
		name, tenant, namespace)), nil
}

// handleListBuiltInSinks handles listing all built-in sink connectors
func (b *PulsarAdminSinksToolBuilder) handleListBuiltInSinks(_ context.Context, admin cmdutils.Client) (*mcp.CallToolResult, error) {
	sinks, err := admin.Sinks().GetBuiltInSinks()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list built-in sinks: %v. There might be an issue connecting to the Pulsar cluster.", err)), nil
	}

	// Convert result to JSON string
	sinksJSON, err := json.Marshal(sinks)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize built-in sinks: %v", err)), nil
	}

	return mcp.NewToolResultText(string(sinksJSON)), nil
}

func (b *PulsarAdminSinksToolBuilder) buildSinkConfig(tenant, namespace, name string, request mcp.CallToolRequest) (*utils.SinkConfig, string, string, error) {
	config := &utils.SinkConfig{}
	args := request.GetArguments()

	if configFile, ok := getStringArg(args, "sink-config-file"); ok && configFile != "" {
		//nolint:gosec
		data, err := os.ReadFile(configFile)
		if err != nil {
			return nil, "", "", fmt.Errorf("load sink config file failed: %w", err)
		}
		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, "", "", fmt.Errorf("unmarshal sink config file error: %w", err)
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

	if className, ok := getStringArg(args, "classname"); ok && className != "" {
		config.ClassName = className
	}

	if processingGuarantees, ok := getStringArg(args, "processing-guarantees"); ok && processingGuarantees != "" {
		config.ProcessingGuarantees = processingGuarantees
	}

	if retainOrdering, ok := getBoolArg(args, "retain-ordering"); ok {
		config.RetainOrdering = retainOrdering
	}

	if retainKeyOrdering, ok := getBoolArg(args, "retain-key-ordering"); ok {
		config.RetainKeyOrdering = retainKeyOrdering
	}

	if autoAck, ok := getBoolArg(args, "auto-ack"); ok {
		config.AutoAck = autoAck
	}

	if inputsValue, exists := args["inputs"]; exists && inputsValue != nil {
		inputs, err := parseStringSlice(inputsValue)
		if err != nil {
			return nil, "", "", fmt.Errorf("invalid inputs: %w", err)
		}
		config.Inputs = inputs
	}

	if topicsPattern, ok := getStringArg(args, "topics-pattern"); ok && topicsPattern != "" {
		config.TopicsPattern = &topicsPattern
	}

	if subsName, ok := getStringArg(args, "subs-name"); ok && subsName != "" {
		config.SourceSubscriptionName = subsName
	}

	if subsPosition, ok := getStringArg(args, "subs-position"); ok && subsPosition != "" {
		config.SourceSubscriptionPosition = subsPosition
	}

	if parallelism, ok, err := parseOptionalIntArg(args, "parallelism"); err != nil {
		return nil, "", "", err
	} else if ok {
		config.Parallelism = parallelism
	}

	if cpu, ok, err := parseOptionalFloatArg(args, "cpu"); err != nil {
		return nil, "", "", err
	} else if ok {
		config.Resources = ensureResources(config.Resources)
		config.Resources.CPU = cpu
	}

	if ram, ok, err := parseOptionalInt64Arg(args, "ram"); err != nil {
		return nil, "", "", err
	} else if ok {
		config.Resources = ensureResources(config.Resources)
		config.Resources.RAM = ram
	}

	if disk, ok, err := parseOptionalInt64Arg(args, "disk"); err != nil {
		return nil, "", "", err
	} else if ok {
		config.Resources = ensureResources(config.Resources)
		config.Resources.Disk = disk
	}

	if customSerdeInputsValue, exists := args["custom-serde-inputs"]; exists && customSerdeInputsValue != nil {
		customSerdeInputs, err := decodeStringMap(customSerdeInputsValue)
		if err != nil {
			return nil, "", "", fmt.Errorf("invalid custom-serde-inputs: %w", err)
		}
		config.TopicToSerdeClassName = customSerdeInputs
	}

	if customSchemaInputsValue, exists := args["custom-schema-inputs"]; exists && customSchemaInputsValue != nil {
		customSchemaInputs, err := decodeStringMap(customSchemaInputsValue)
		if err != nil {
			return nil, "", "", fmt.Errorf("invalid custom-schema-inputs: %w", err)
		}
		config.TopicToSchemaType = customSchemaInputs
	}

	if inputSpecsValue, exists := args["input-specs"]; exists && inputSpecsValue != nil {
		inputSpecs, err := decodeConsumerConfigMap(inputSpecsValue)
		if err != nil {
			return nil, "", "", fmt.Errorf("invalid input-specs: %w", err)
		}
		config.InputSpecs = inputSpecs
	}

	if maxMessageRetries, ok, err := parseOptionalIntArg(args, "max-redeliver-count"); err != nil {
		return nil, "", "", err
	} else if ok {
		config.MaxMessageRetries = maxMessageRetries
	} else {
		config.MaxMessageRetries = 0
	}

	if deadLetterTopic, ok := getStringArg(args, "dead-letter-topic"); ok && deadLetterTopic != "" {
		config.DeadLetterTopic = deadLetterTopic
	}

	if timeoutMs, ok, err := parseOptionalInt64Arg(args, "timeout-ms"); err != nil {
		return nil, "", "", err
	} else if ok {
		config.TimeoutMs = &timeoutMs
	}

	if negativeAckRedeliveryDelayMs, ok, err := parseOptionalInt64Arg(args, "negative-ack-redelivery-delay-ms"); err != nil {
		return nil, "", "", err
	} else if ok {
		config.NegativeAckRedeliveryDelayMs = negativeAckRedeliveryDelayMs
	} else {
		config.NegativeAckRedeliveryDelayMs = 0
	}

	if customRuntimeOptions, ok := getStringArg(args, "custom-runtime-options"); ok && customRuntimeOptions != "" {
		config.CustomRuntimeOptions = customRuntimeOptions
	}

	if secretsValue, exists := args["secrets"]; exists && secretsValue != nil {
		secrets, err := decodeInterfaceMap(secretsValue)
		if err != nil {
			return nil, "", "", fmt.Errorf("invalid secrets: %w", err)
		}
		config.Secrets = secrets
	}

	if cleanupSubscription, ok := getBoolArg(args, "cleanup-subscription"); ok {
		config.CleanupSubscription = cleanupSubscription
	} else {
		config.CleanupSubscription = true
	}

	if sinkConfigValue, exists := args["sink-config"]; exists && sinkConfigValue != nil {
		sinkConfigs, err := decodeInterfaceMap(sinkConfigValue)
		if err != nil {
			return nil, "", "", fmt.Errorf("invalid sink-config: %w", err)
		}
		config.Configs = sinkConfigs
	}

	if transformFunction, ok := getStringArg(args, "transform-function"); ok && transformFunction != "" {
		config.TransformFunction = transformFunction
	}

	if transformFunctionClassName, ok := getStringArg(args, "transform-function-classname"); ok && transformFunctionClassName != "" {
		config.TransformFunctionClassName = transformFunctionClassName
	}

	if transformFunctionConfig, ok := getStringArg(args, "transform-function-config"); ok && transformFunctionConfig != "" {
		config.TransformFunctionConfig = transformFunctionConfig
	}

	b.normalizeSinkConfigMaps(config)

	archiveArg, _ := getStringArg(args, "archive")
	sinkTypeArg, _ := getStringArg(args, "sink-type")

	return config, archiveArg, sinkTypeArg, nil
}

func (b *PulsarAdminSinksToolBuilder) normalizeSinkConfigMaps(config *utils.SinkConfig) {
	if config.Configs != nil {
		if converted, ok := convertMap(config.Configs).(map[string]interface{}); ok {
			config.Configs = converted
		}
	}

	if config.Secrets != nil {
		if converted, ok := convertMap(config.Secrets).(map[string]interface{}); ok {
			config.Secrets = converted
		}
	}

	if config.Secrets == nil {
		config.Secrets = make(map[string]interface{})
	}
}

func (b *PulsarAdminSinksToolBuilder) applySinkDefaults(config *utils.SinkConfig) {
	if config.Tenant == "" {
		config.Tenant = defaultTenant
	}
	if config.Namespace == "" {
		config.Namespace = defaultNamespace
	}
	if config.Parallelism <= 0 {
		config.Parallelism = 1
	}
}

func validateSinkConfig(config *utils.SinkConfig) error {
	if config.Archive == "" {
		return fmt.Errorf("sink archive not specified")
	}
	if config.Name == "" {
		return fmt.Errorf("sink name not specified")
	}
	return nil
}

func validateSinkArchiveArgs(archive, sinkType string) error {
	if archive != "" && sinkType != "" {
		return fmt.Errorf("cannot specify both 'archive' and 'sink-type'")
	}
	return nil
}

func (b *PulsarAdminSinksToolBuilder) resolveSinkArchive(admin cmdutils.Client, config *utils.SinkConfig, archive, sinkType string) (string, error) {
	if sinkType != "" {
		resolved, err := b.validateSinkType(admin, sinkType)
		if err != nil {
			return "", err
		}
		config.Archive = resolved
		return "", nil
	}

	if archive != "" {
		config.Archive = archive
		return archive, nil
	}

	return "", nil
}

func (b *PulsarAdminSinksToolBuilder) validateSinkType(admin cmdutils.Client, sinkType string) (string, error) {
	builtins, err := admin.Sinks().GetBuiltInSinks()
	if err != nil {
		return "", fmt.Errorf("failed to list built-in sinks: %w", err)
	}

	names := make([]string, 0, len(builtins))
	for _, builtin := range builtins {
		names = append(names, builtin.Name)
		if builtin.Name == sinkType {
			return "builtin://" + sinkType, nil
		}
	}

	sort.Strings(names)
	return "", fmt.Errorf("invalid sink-type %q. Available sinks: %s", sinkType, strings.Join(names, ", "))
}

// isPackageURLSupported checks if the package URL is supported
// Validates URLs for Pulsar sink packages
func (b *PulsarAdminSinksToolBuilder) isPackageURLSupported(archive string) bool {
	return archive != "" && (strings.HasPrefix(archive, "http") ||
		strings.HasPrefix(archive, "file") ||
		strings.HasPrefix(archive, "function") ||
		strings.HasPrefix(archive, "sink") ||
		strings.HasPrefix(archive, "source"))
}
