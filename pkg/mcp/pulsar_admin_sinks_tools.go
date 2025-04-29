package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/admin/config"
	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
)

// PulsarAdminAddSinksTools adds a unified sink-related tool to the MCP server
func PulsarAdminAddSinksTools(s *server.MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeaturePulsarAdminSinks)) && !slices.Contains(features, string(FeatureAll)) && !slices.Contains(features, string(FeatureAllPulsar)) {
		return
	}

	// Create a single unified tool for all sink operations
	toolDesc := "Manage Apache Pulsar Sinks for data movement and integration. " +
		"Pulsar Sinks are connectors that export data from Pulsar topics to external systems such as databases, " +
		"storage services, messaging systems, and third-party applications. " +
		"Sinks consume messages from one or more Pulsar topics, transform the data if needed, " +
		"and write it to external systems in a format compatible with the target destination. " +
		"Built-in sink connectors are available for common systems like Kafka, JDBC, Elasticsearch, and cloud storage. " +
		"Sinks follow the tenant/namespace/name hierarchy for organization and access control, " +
		"can scale through parallelism configuration, and support configurable subscription types. " +
		"This tool provides complete lifecycle management including deployment, configuration, " +
		"monitoring, and runtime control. Sinks require proper permissions to access their input topics."

	operationDesc := "Operation to perform. Available operations:\n" +
		"- list: List all sinks under a specific tenant and namespace\n" +
		"- get: Get the configuration of a sink\n" +
		"- status: Get the runtime status of a sink (instances, metrics)\n" +
		"- create: Deploy a new sink with specified parameters\n" +
		"- update: Update the configuration of an existing sink\n" +
		"- delete: Delete a sink\n" +
		"- start: Start a stopped sink\n" +
		"- stop: Stop a running sink\n" +
		"- restart: Restart a sink\n" +
		"- list-built-in: List all built-in sink connectors available in the system"

	sinkTool := mcp.NewTool("pulsar_admin_sinks",
		mcp.WithDescription(toolDesc),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc)),
		mcp.WithString("tenant",
			mcp.Description("The tenant name. Tenants are the primary organizational unit in Pulsar, "+
				"providing multi-tenancy and resource isolation. Sinks deployed within a tenant "+
				"inherit its permissions and resource quotas. "+
				"Required for all operations except 'list-built-in'.")),
		mcp.WithString("namespace",
			mcp.Description("The namespace name. Namespaces are logical groupings of topics and sinks "+
				"within a tenant. They encapsulate configuration policies and access control. "+
				"Sinks in a namespace typically process topics within the same namespace. "+
				"Required for all operations except 'list-built-in'.")),
		mcp.WithString("name",
			mcp.Description("The sink name. Required for all operations except 'list' and 'list-built-in'. "+
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
				"Either inputs or topics-pattern must be specified.")),
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
		mcp.WithNumber("parallelism",
			mcp.Description("The parallelism factor of the sink. Optional for 'create' and 'update' operations. "+
				"Determines how many instances of the sink will run concurrently. "+
				"Higher values improve throughput but require more resources. "+
				"Default is 1 (single instance). Recommended to align with topic partition count "+
				"when consuming from partitioned topics.")),
		mcp.WithObject("sink-config",
			mcp.Description("User-defined sink config key/values. Optional for 'create' and 'update' operations. "+
				"Provides configuration parameters specific to the sink connector being used. "+
				"For example, JDBC connection strings, Elasticsearch indices, S3 bucket details, etc. "+
				"Specify as a JSON object with configuration properties required by the specific sink type. "+
				"Example: {\"jdbcUrl\": \"jdbc:postgresql://localhost:5432/database\", \"tableName\": \"events\"}")),
	)

	// Add the tool with the handler that captures readOnly state
	s.AddTool(sinkTool, handleSinksTool(readOnly))
}

// handleSinksTool returns a handler function for the unified sinks tool
func handleSinksTool(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract and validate operation parameter
		operation, err := requiredParam[string](request.Params.Arguments, "operation")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'operation': %v", err)), nil
		}

		// Check if the operation is valid
		validOperations := map[string]bool{
			"list": true, "get": true, "status": true, "create": true, "update": true,
			"delete": true, "start": true, "stop": true, "restart": true, "list-built-in": true,
		}

		if !validOperations[operation] {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid operation: '%s'. Supported operations: list, get, status, create, update, delete, start, stop, restart, list-built-in", operation)), nil
		}

		// Check write permissions for write operations
		writeOperations := map[string]bool{
			"create": true, "update": true, "delete": true, "start": true,
			"stop": true, "restart": true,
		}

		if readOnly && writeOperations[operation] {
			return mcp.NewToolResultError(fmt.Sprintf("Operation '%s' not allowed in read-only mode. Read-only mode restricts modifications to Pulsar Sinks.", operation)), nil
		}

		// Create admin client
		admin := cmdutils.NewPulsarClientWithAPIVersion(config.V3)

		// List built-in sinks doesn't require tenant, namespace or name
		if operation == "list-built-in" {
			return handleListBuiltInSinks(ctx, admin)
		}

		// Extract common parameters (all operations except list-built-in require tenant and namespace)
		tenant, err := requiredParam[string](request.Params.Arguments, "tenant")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'tenant': %v. A tenant is required for operation '%s'.", err, operation)), nil
		}

		namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'namespace': %v. A namespace is required for operation '%s'.", err, operation)), nil
		}

		// For all operations except 'list', name is required
		var name string
		if operation != "list" {
			name, err = requiredParam[string](request.Params.Arguments, "name")
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'name' for operation '%s': %v. The sink name must be specified for this operation.", operation, err)), nil
			}
		}

		// Handle operations
		switch operation {
		case "list":
			return handleSinkList(ctx, admin, tenant, namespace)
		case "get":
			return handleSinkGet(ctx, admin, tenant, namespace, name)
		case "status":
			return handleSinkStatus(ctx, admin, tenant, namespace, name)
		case "create":
			return handleSinkCreate(ctx, admin, request.Params.Arguments)
		case "update":
			return handleSinkUpdate(ctx, admin, request.Params.Arguments)
		case "delete":
			return handleSinkDelete(ctx, admin, tenant, namespace, name)
		case "start":
			return handleSinkStart(ctx, admin, tenant, namespace, name)
		case "stop":
			return handleSinkStop(ctx, admin, tenant, namespace, name)
		case "restart":
			return handleSinkRestart(ctx, admin, tenant, namespace, name)
		default:
			// This should never happen due to the valid operations check above
			return mcp.NewToolResultError(fmt.Sprintf("Unsupported operation: %s", operation)), nil
		}
	}
}

// handleSinkList handles listing all sinks under a namespace
func handleSinkList(_ context.Context, admin cmdutils.Client, tenant, namespace string) (*mcp.CallToolResult, error) {
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
func handleSinkGet(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
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
func handleSinkStatus(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
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
func handleSinkCreate(_ context.Context, admin cmdutils.Client, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	tenant, err := requiredParam[string](arguments, "tenant")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get tenant: %v", err)), nil
	}

	namespace, err := requiredParam[string](arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace: %v", err)), nil
	}

	name, err := requiredParam[string](arguments, "name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get name: %v", err)), nil
	}

	// Create a new SinkData object
	sinkData := &utils.SinkData{
		Tenant:    tenant,
		Namespace: namespace,
		Name:      name,
		SinkConf:  &utils.SinkConfig{},
	}

	// Get optional parameters
	archive, hasArchive := optionalParam[string](arguments, "archive")
	if hasArchive && archive != "" {
		sinkData.Archive = archive
	}

	sinkType, hasSinkType := optionalParam[string](arguments, "sink-type")
	if hasSinkType && sinkType != "" {
		sinkData.SinkType = sinkType
	}

	inputsArray, hasInputs := optionalParamArray[string](arguments, "inputs")
	if hasInputs && len(inputsArray) > 0 {
		sinkData.Inputs = strings.Join(inputsArray, ",")
	}

	topicsPattern, hasTopicsPattern := optionalParam[string](arguments, "topics-pattern")
	if hasTopicsPattern && topicsPattern != "" {
		sinkData.TopicsPattern = topicsPattern
	}

	subsName, hasSubsName := optionalParam[string](arguments, "subs-name")
	if hasSubsName && subsName != "" {
		sinkData.SubsName = subsName
	}

	parallelismFloat, hasParallelism := optionalParam[float64](arguments, "parallelism")
	if hasParallelism {
		sinkData.Parallelism = int(parallelismFloat)
	}

	// Get sink config if available
	var sinkConfigMap map[string]interface{}
	sinkConfigObj, ok := arguments["sink-config"]
	if ok && sinkConfigObj != nil {
		if configMap, isMap := sinkConfigObj.(map[string]interface{}); isMap {
			sinkConfigMap = configMap
			// Convert to JSON string
			sinkConfigJSON, err := json.Marshal(sinkConfigMap)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal sink-config: %v. Ensure the sink configuration is a valid JSON object.", err)), nil
			}
			sinkData.SinkConfigString = string(sinkConfigJSON)
		}
	}

	// Validate inputs
	if sinkData.Archive == "" && sinkData.SinkType == "" {
		return mcp.NewToolResultError("Missing required parameter: Either 'archive' or 'sink-type' must be specified for sink creation. Use 'archive' for custom connectors or 'sink-type' for built-in connectors."), nil
	}

	if sinkData.Archive != "" && sinkData.SinkType != "" {
		return mcp.NewToolResultError("Invalid parameters: Cannot specify both 'archive' and 'sink-type'. Use only one of these parameters based on your connector type."), nil
	}

	if sinkData.Inputs == "" && sinkData.TopicsPattern == "" {
		return mcp.NewToolResultError("Missing required parameter: Either 'inputs' or 'topics-pattern' must be specified. The sink needs a source of data to consume from Pulsar."), nil
	}

	// Process the arguments
	err = processArguments(sinkData)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to process arguments: %v", err)), nil
	}

	// Create the sink
	if sinkData.Archive != "" && isPackageURLSupported(sinkData.Archive) {
		err = admin.Sinks().CreateSinkWithURL(sinkData.SinkConf, sinkData.Archive)
	} else {
		err = admin.Sinks().CreateSink(sinkData.SinkConf, sinkData.Archive)
	}

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create sink '%s' in tenant '%s' namespace '%s': %v. Verify all parameters are correct and required resources exist.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Created sink '%s' successfully in tenant '%s' namespace '%s'. The sink will start consuming from its input topics and writing to the configured destination.",
		name, tenant, namespace)), nil
}

// handleSinkUpdate handles updating an existing sink
func handleSinkUpdate(_ context.Context, admin cmdutils.Client, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	tenant, err := requiredParam[string](arguments, "tenant")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get tenant: %v", err)), nil
	}

	namespace, err := requiredParam[string](arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace: %v", err)), nil
	}

	name, err := requiredParam[string](arguments, "name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get name: %v", err)), nil
	}

	// Create a new SinkData object
	sinkData := &utils.SinkData{
		Tenant:    tenant,
		Namespace: namespace,
		Name:      name,
		SinkConf:  &utils.SinkConfig{},
	}

	// Get optional parameters
	archive, hasArchive := optionalParam[string](arguments, "archive")
	if hasArchive && archive != "" {
		sinkData.Archive = archive
	}

	sinkType, hasSinkType := optionalParam[string](arguments, "sink-type")
	if hasSinkType && sinkType != "" {
		sinkData.SinkType = sinkType
	}

	inputsArray, hasInputs := optionalParamArray[string](arguments, "inputs")
	if hasInputs && len(inputsArray) > 0 {
		sinkData.Inputs = strings.Join(inputsArray, ",")
	}

	topicsPattern, hasTopicsPattern := optionalParam[string](arguments, "topics-pattern")
	if hasTopicsPattern && topicsPattern != "" {
		sinkData.TopicsPattern = topicsPattern
	}

	subsName, hasSubsName := optionalParam[string](arguments, "subs-name")
	if hasSubsName && subsName != "" {
		sinkData.SubsName = subsName
	}

	parallelismFloat, hasParallelism := optionalParam[float64](arguments, "parallelism")
	if hasParallelism {
		sinkData.Parallelism = int(parallelismFloat)
	}

	// Get sink config if available
	var sinkConfigMap map[string]interface{}
	sinkConfigObj, ok := arguments["sink-config"]
	if ok && sinkConfigObj != nil {
		if configMap, isMap := sinkConfigObj.(map[string]interface{}); isMap {
			sinkConfigMap = configMap
			// Convert to JSON string
			sinkConfigJSON, err := json.Marshal(sinkConfigMap)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal sink-config: %v. Ensure the sink configuration is a valid JSON object.", err)), nil
			}
			sinkData.SinkConfigString = string(sinkConfigJSON)
		}
	}

	// Validate inputs if both are specified
	if sinkData.Archive != "" && sinkData.SinkType != "" {
		return mcp.NewToolResultError("Invalid parameters: Cannot specify both 'archive' and 'sink-type'. Use only one of these parameters based on your connector type."), nil
	}

	// Process the arguments
	err = processArguments(sinkData)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to process arguments: %v", err)), nil
	}

	// Create update options
	updateOptions := &utils.UpdateOptions{
		UpdateAuthData: true,
	}

	// Update the sink
	if sinkData.Archive != "" && isPackageURLSupported(sinkData.Archive) {
		err = admin.Sinks().UpdateSinkWithURL(sinkData.SinkConf, sinkData.Archive, updateOptions)
	} else {
		err = admin.Sinks().UpdateSink(sinkData.SinkConf, sinkData.Archive, updateOptions)
	}

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to update sink '%s' in tenant '%s' namespace '%s': %v. Verify the sink exists and all parameters are valid.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Updated sink '%s' successfully in tenant '%s' namespace '%s'. The sink may need to be restarted to apply all changes.",
		name, tenant, namespace)), nil
}

// handleSinkDelete handles deleting a sink
func handleSinkDelete(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	err := admin.Sinks().DeleteSink(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete sink '%s' in tenant '%s' namespace '%s': %v. Verify the sink exists and you have deletion permissions.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Deleted sink '%s' successfully from tenant '%s' namespace '%s'. All running instances have been terminated.",
		name, tenant, namespace)), nil
}

// handleSinkStart handles starting a sink
func handleSinkStart(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	err := admin.Sinks().StartSink(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to start sink '%s' in tenant '%s' namespace '%s': %v. Verify the sink exists and is not already running.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Started sink '%s' successfully in tenant '%s' namespace '%s'. The sink will begin consuming from its input topics.",
		name, tenant, namespace)), nil
}

// handleSinkStop handles stopping a sink
func handleSinkStop(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	err := admin.Sinks().StopSink(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to stop sink '%s' in tenant '%s' namespace '%s': %v. Verify the sink exists and is currently running.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Stopped sink '%s' successfully in tenant '%s' namespace '%s'. The sink will no longer consume messages until restarted.",
		name, tenant, namespace)), nil
}

// handleSinkRestart handles restarting a sink
func handleSinkRestart(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	err := admin.Sinks().RestartSink(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to restart sink '%s' in tenant '%s' namespace '%s': %v. Verify the sink exists and is properly deployed.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Restarted sink '%s' successfully in tenant '%s' namespace '%s'. All sink instances have been restarted.",
		name, tenant, namespace)), nil
}

// handleListBuiltInSinks handles listing all built-in sink connectors
func handleListBuiltInSinks(_ context.Context, admin cmdutils.Client) (*mcp.CallToolResult, error) {
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

// processArguments is a simplified version of the pulsarctl function to process sink arguments
func processArguments(sinkData *utils.SinkData) error {
	// Initialize config if needed
	if sinkData.SinkConf == nil {
		sinkData.SinkConf = new(utils.SinkConfig)
	}

	// Set basic config values
	sinkData.SinkConf.Tenant = sinkData.Tenant
	sinkData.SinkConf.Namespace = sinkData.Namespace
	sinkData.SinkConf.Name = sinkData.Name

	// Set inputs if provided
	if sinkData.Inputs != "" {
		inputTopics := strings.Split(sinkData.Inputs, ",")
		sinkData.SinkConf.Inputs = inputTopics
	}

	// Set topics pattern if provided
	if sinkData.TopicsPattern != "" {
		sinkData.SinkConf.TopicsPattern = &sinkData.TopicsPattern
	}

	// Set subscription name if provided
	if sinkData.SubsName != "" {
		sinkData.SinkConf.SourceSubscriptionName = sinkData.SubsName
	}

	// Set parallelism if provided
	if sinkData.Parallelism != 0 {
		sinkData.SinkConf.Parallelism = sinkData.Parallelism
	} else if sinkData.SinkConf.Parallelism <= 0 {
		sinkData.SinkConf.Parallelism = 1
	}

	// Handle archive and sink-type
	if sinkData.Archive != "" && sinkData.SinkType != "" {
		return fmt.Errorf("cannot specify both archive and sink-type")
	}

	if sinkData.Archive != "" {
		sinkData.SinkConf.Archive = sinkData.Archive
	}

	if sinkData.SinkType != "" {
		// In a real implementation, we would validate the sink type here
		sinkData.SinkConf.Archive = sinkData.SinkType
	}

	// Parse sink config if provided
	if sinkData.SinkConfigString != "" {
		var configs map[string]interface{}
		if err := json.Unmarshal([]byte(sinkData.SinkConfigString), &configs); err != nil {
			return fmt.Errorf("failed to parse sink config: %v", err)
		}
		sinkData.SinkConf.Configs = configs
	}

	return nil
}
