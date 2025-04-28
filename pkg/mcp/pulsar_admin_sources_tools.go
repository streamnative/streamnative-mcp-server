package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/admin/config"
	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
)

// PulsarAdminAddSourcesTools adds a unified source-related tool to the MCP server
func PulsarAdminAddSourcesTools(s *server.MCPServer, readOnly bool) {
	// Create a single unified tool for all source operations
	toolDesc := "Manage Apache Pulsar Sources for data ingestion and integration. " +
		"Pulsar Sources are connectors that import data from external systems into Pulsar topics. " +
		"Sources connect to external systems such as databases, messaging platforms, storage services, " +
		"and real-time data streams to pull data and publish it to Pulsar topics. " +
		"Built-in source connectors are available for common systems like Kafka, JDBC, AWS services, and more. " +
		"Sources follow the tenant/namespace/name hierarchy for organization and access control, " +
		"can scale through parallelism configuration, and support various processing guarantees. " +
		"This tool provides complete lifecycle management including deployment, configuration, " +
		"monitoring, and runtime control. Sources use schema types to ensure data compatibility."

	operationDesc := "Operation to perform. Available operations:\n" +
		"- list: List all sources under a specific tenant and namespace\n" +
		"- get: Get the configuration of a source\n" +
		"- status: Get the runtime status of a source (instances, metrics)\n" +
		"- create: Deploy a new source with specified parameters\n" +
		"- update: Update the configuration of an existing source\n" +
		"- delete: Delete a source\n" +
		"- start: Start a stopped source\n" +
		"- stop: Stop a running source\n" +
		"- restart: Restart a source\n" +
		"- list-built-in: List all built-in source connectors available in the system"

	sourceTool := mcp.NewTool("pulsar_admin_sources",
		mcp.WithDescription(toolDesc),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc)),
		mcp.WithString("tenant",
			mcp.Description("The tenant name. Tenants are the primary organizational unit in Pulsar, "+
				"providing multi-tenancy and resource isolation. Sources deployed within a tenant "+
				"inherit its permissions and resource quotas. "+
				"Required for all operations except 'list-built-in'.")),
		mcp.WithString("namespace",
			mcp.Description("The namespace name. Namespaces are logical groupings of topics and sources "+
				"within a tenant. They encapsulate configuration policies and access control. "+
				"Sources in a namespace typically publish to topics within the same namespace. "+
				"Required for all operations except 'list-built-in'.")),
		mcp.WithString("name",
			mcp.Description("The source name. Required for all operations except 'list' and 'list-built-in'. "+
				"Names should be descriptive of the source's purpose and must be unique within a namespace. "+
				"Source names are used in metrics, logs, and when addressing the source via APIs.")),
		mcp.WithString("archive",
			mcp.Description("Path to the archive file containing the source code. Optional for 'create' and 'update' operations. "+
				"Can be a local path, NAR file, or a URL accessible to the Pulsar broker. "+
				"The archive should contain all dependencies for the source connector. "+
				"Either archive or source-type must be specified, but not both.")),
		mcp.WithString("source-type",
			mcp.Description("The built-in source connector type to use. Optional for 'create' and 'update' operations. "+
				"Specifies which built-in connector to use, such as 'kafka', 'jdbc', 'file', etc. "+
				"Use 'list-built-in' operation to see available source types. "+
				"Either source-type or archive must be specified, but not both.")),
		mcp.WithString("destination-topic-name",
			mcp.Description("The Pulsar topic to which data is published. Required for 'create' operation, optional for 'update'. "+
				"Specified in the format 'persistent://tenant/namespace/topic'. "+
				"This is the topic where the source will send the data it extracts from the external system. "+
				"The topic will be automatically created if it doesn't exist.")),
		mcp.WithString("deserialization-classname",
			mcp.Description("The SerDe (Serialization/Deserialization) classname for the source. Optional for 'create' and 'update'. "+
				"Specifies how to convert data from the external system into Pulsar messages. "+
				"Common SerDe classes include AvroSchema, JsonSchema, StringSchema, etc. "+
				"If not specified, the source will use the default SerDe for the connector type.")),
		mcp.WithString("schema-type",
			mcp.Description("The schema type to be used to encode messages emitted from the source. Optional for 'create' and 'update'. "+
				"Available schema types include: 'avro', 'json', 'protobuf', 'string', etc. "+
				"Schema types ensure data compatibility and enable schema evolution. "+
				"The schema type should match the format of data being ingested.")),
		mcp.WithString("classname",
			mcp.Description("The source's class name if archive is a file-url-path (file://...). Optional for 'create' and 'update'. "+
				"This specifies the fully qualified class name that implements the source connector. "+
				"Only needed when using a custom source implementation in a JAR file. "+
				"Built-in connectors don't require this parameter.")),
		mcp.WithString("processing-guarantees",
			mcp.Description("The processing guarantees (delivery semantics) applied to the source. Optional for 'create' and 'update'. "+
				"Available options: 'atleast_once', 'atmost_once', 'effectively_once'. "+
				"Controls how data is delivered in failure scenarios. "+
				"'atleast_once' is the most common and ensures no data loss but may have duplicates. "+
				"Default is 'atleast_once'.")),
		mcp.WithNumber("parallelism",
			mcp.Description("The parallelism factor of the source. Optional for 'create' and 'update' operations. "+
				"Determines how many instances of the source will run concurrently. "+
				"Higher values improve throughput but require more resources. "+
				"Default is 1 (single instance). Recommended to align with both source capacity "+
				"and destination topic partition count.")),
		mcp.WithObject("source-config",
			mcp.Description("User-defined source config key/values. Optional for 'create' and 'update' operations. "+
				"Provides configuration parameters specific to the source connector being used. "+
				"For example, database connection details, Kafka bootstrap servers, credentials, etc. "+
				"Specify as a JSON object with configuration properties required by the specific source type. "+
				"Example: {\"topic\": \"external-kafka-topic\", \"bootstrapServers\": \"kafka:9092\"}")),
	)

	// Add the tool with the handler that captures readOnly state
	s.AddTool(sourceTool, handleSourcesTool(readOnly))
}

// handleSourcesTool returns a handler function for the unified sources tool
func handleSourcesTool(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
			return mcp.NewToolResultError(fmt.Sprintf("Operation '%s' not allowed in read-only mode. Read-only mode restricts modifications to Pulsar Sources.", operation)), nil
		}

		// Create admin client
		admin := cmdutils.NewPulsarClientWithAPIVersion(config.V3)

		// List built-in sources doesn't require tenant, namespace or name
		if operation == "list-built-in" {
			return handleListBuiltInSources(ctx, admin)
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
				return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'name' for operation '%s': %v. The source name must be specified for this operation.", operation, err)), nil
			}
		}

		// Handle operations
		switch operation {
		case "list":
			return handleSourceList(ctx, admin, tenant, namespace)
		case "get":
			return handleSourceGet(ctx, admin, tenant, namespace, name)
		case "status":
			return handleSourceStatus(ctx, admin, tenant, namespace, name)
		case "create":
			return handleSourceCreate(ctx, admin, request.Params.Arguments)
		case "update":
			return handleSourceUpdate(ctx, admin, request.Params.Arguments)
		case "delete":
			return handleSourceDelete(ctx, admin, tenant, namespace, name)
		case "start":
			return handleSourceStart(ctx, admin, tenant, namespace, name)
		case "stop":
			return handleSourceStop(ctx, admin, tenant, namespace, name)
		case "restart":
			return handleSourceRestart(ctx, admin, tenant, namespace, name)
		default:
			// This should never happen due to the valid operations check above
			return mcp.NewToolResultError(fmt.Sprintf("Unsupported operation: %s", operation)), nil
		}
	}
}

// handleSourceList handles listing all sources under a namespace
func handleSourceList(ctx context.Context, admin cmdutils.Client, tenant, namespace string) (*mcp.CallToolResult, error) {
	sources, err := admin.Sources().ListSources(tenant, namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list sources in tenant '%s' namespace '%s': %v. Check that the tenant and namespace exist and you have proper permissions.",
			tenant, namespace, err)), nil
	}

	// Convert result to JSON string
	sourcesJSON, err := json.Marshal(sources)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize source list: %v", err)), nil
	}

	return mcp.NewToolResultText(string(sourcesJSON)), nil
}

// handleSourceGet handles getting information about a source
func handleSourceGet(ctx context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	source, err := admin.Sources().GetSource(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get source '%s' in tenant '%s' namespace '%s': %v. Verify the source exists and you have proper permissions.",
			name, tenant, namespace, err)), nil
	}

	// Convert result to JSON string
	sourceJSON, err := json.Marshal(source)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize source info: %v", err)), nil
	}

	return mcp.NewToolResultText(string(sourceJSON)), nil
}

// handleSourceStatus handles getting the status of a source
func handleSourceStatus(ctx context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	status, err := admin.Sources().GetSourceStatus(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get status for source '%s' in tenant '%s' namespace '%s': %v. Verify the source exists and is properly deployed.",
			name, tenant, namespace, err)), nil
	}

	// Convert result to JSON string
	statusJSON, err := json.Marshal(status)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize source status: %v", err)), nil
	}

	return mcp.NewToolResultText(string(statusJSON)), nil
}

// handleSourceCreate handles creating a new source
func handleSourceCreate(ctx context.Context, admin cmdutils.Client, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
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

	// Create a new SourceData object
	sourceData := &utils.SourceData{
		Tenant:     tenant,
		Namespace:  namespace,
		Name:       name,
		SourceConf: &utils.SourceConfig{},
	}

	// Get optional parameters
	archive, hasArchive := optionalParam[string](arguments, "archive")
	if hasArchive && archive != "" {
		sourceData.Archive = archive
	}

	sourceType, hasSourceType := optionalParam[string](arguments, "source-type")
	if hasSourceType && sourceType != "" {
		sourceData.SourceType = sourceType
	}

	destTopic, hasDestTopic := optionalParam[string](arguments, "destination-topic-name")
	if hasDestTopic && destTopic != "" {
		sourceData.DestinationTopicName = destTopic
	}

	deserializationClassName, hasDeserialization := optionalParam[string](arguments, "deserialization-classname")
	if hasDeserialization && deserializationClassName != "" {
		sourceData.DeserializationClassName = deserializationClassName
	}

	schemaType, hasSchemaType := optionalParam[string](arguments, "schema-type")
	if hasSchemaType && schemaType != "" {
		sourceData.SchemaType = schemaType
	}

	className, hasClassName := optionalParam[string](arguments, "classname")
	if hasClassName && className != "" {
		sourceData.ClassName = className
	}

	processingGuarantees, hasProcessingGuarantees := optionalParam[string](arguments, "processing-guarantees")
	if hasProcessingGuarantees && processingGuarantees != "" {
		sourceData.ProcessingGuarantees = processingGuarantees
	}

	parallelismFloat, hasParallelism := optionalParam[float64](arguments, "parallelism")
	if hasParallelism {
		sourceData.Parallelism = int(parallelismFloat)
	}

	// Get source config if available
	var sourceConfigMap map[string]interface{}
	sourceConfigObj, ok := arguments["source-config"]
	if ok && sourceConfigObj != nil {
		if configMap, isMap := sourceConfigObj.(map[string]interface{}); isMap {
			sourceConfigMap = configMap
			// Convert to JSON string
			sourceConfigJSON, err := json.Marshal(sourceConfigMap)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal source-config: %v. Ensure the source configuration is a valid JSON object.", err)), nil
			}
			sourceData.SourceConfigString = string(sourceConfigJSON)
		}
	}

	// Validate inputs
	if sourceData.Archive == "" && sourceData.SourceType == "" {
		return mcp.NewToolResultError("Missing required parameter: Either 'archive' or 'source-type' must be specified for source creation. Use 'archive' for custom connectors or 'source-type' for built-in connectors."), nil
	}

	if sourceData.Archive != "" && sourceData.SourceType != "" {
		return mcp.NewToolResultError("Invalid parameters: Cannot specify both 'archive' and 'source-type'. Use only one of these parameters based on your connector type."), nil
	}

	if sourceData.DestinationTopicName == "" {
		return mcp.NewToolResultError("Missing required parameter: 'destination-topic-name' must be specified. This is the Pulsar topic where the source will publish data."), nil
	}

	// Process the arguments
	err = processSourceArguments(sourceData)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to process arguments: %v", err)), nil
	}

	// Create the source
	if sourceData.Archive != "" && isPackageURLSupported(sourceData.Archive) {
		err = admin.Sources().CreateSourceWithURL(sourceData.SourceConf, sourceData.Archive)
	} else {
		err = admin.Sources().CreateSource(sourceData.SourceConf, sourceData.Archive)
	}

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create source '%s' in tenant '%s' namespace '%s': %v. Verify all parameters are correct and required resources exist.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Created source '%s' successfully in tenant '%s' namespace '%s'. The source will start pulling data from the external system and publishing to the destination topic.",
		name, tenant, namespace)), nil
}

// handleSourceUpdate handles updating an existing source
func handleSourceUpdate(ctx context.Context, admin cmdutils.Client, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
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

	// Create a new SourceData object
	sourceData := &utils.SourceData{
		Tenant:     tenant,
		Namespace:  namespace,
		Name:       name,
		SourceConf: &utils.SourceConfig{},
	}

	// Get optional parameters
	archive, hasArchive := optionalParam[string](arguments, "archive")
	if hasArchive && archive != "" {
		sourceData.Archive = archive
	}

	sourceType, hasSourceType := optionalParam[string](arguments, "source-type")
	if hasSourceType && sourceType != "" {
		sourceData.SourceType = sourceType
	}

	destTopic, hasDestTopic := optionalParam[string](arguments, "destination-topic-name")
	if hasDestTopic && destTopic != "" {
		sourceData.DestinationTopicName = destTopic
	}

	deserializationClassName, hasDeserialization := optionalParam[string](arguments, "deserialization-classname")
	if hasDeserialization && deserializationClassName != "" {
		sourceData.DeserializationClassName = deserializationClassName
	}

	schemaType, hasSchemaType := optionalParam[string](arguments, "schema-type")
	if hasSchemaType && schemaType != "" {
		sourceData.SchemaType = schemaType
	}

	className, hasClassName := optionalParam[string](arguments, "classname")
	if hasClassName && className != "" {
		sourceData.ClassName = className
	}

	processingGuarantees, hasProcessingGuarantees := optionalParam[string](arguments, "processing-guarantees")
	if hasProcessingGuarantees && processingGuarantees != "" {
		sourceData.ProcessingGuarantees = processingGuarantees
	}

	parallelismFloat, hasParallelism := optionalParam[float64](arguments, "parallelism")
	if hasParallelism {
		sourceData.Parallelism = int(parallelismFloat)
	}

	// Get source config if available
	var sourceConfigMap map[string]interface{}
	sourceConfigObj, ok := arguments["source-config"]
	if ok && sourceConfigObj != nil {
		if configMap, isMap := sourceConfigObj.(map[string]interface{}); isMap {
			sourceConfigMap = configMap
			// Convert to JSON string
			sourceConfigJSON, err := json.Marshal(sourceConfigMap)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal source-config: %v. Ensure the source configuration is a valid JSON object.", err)), nil
			}
			sourceData.SourceConfigString = string(sourceConfigJSON)
		}
	}

	// Validate inputs if both are specified
	if sourceData.Archive != "" && sourceData.SourceType != "" {
		return mcp.NewToolResultError("Invalid parameters: Cannot specify both 'archive' and 'source-type'. Use only one of these parameters based on your connector type."), nil
	}

	// Process the arguments
	err = processSourceArguments(sourceData)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to process arguments: %v", err)), nil
	}

	// Create update options
	updateOptions := &utils.UpdateOptions{
		UpdateAuthData: true,
	}

	// Update the source
	if sourceData.Archive != "" && isPackageURLSupported(sourceData.Archive) {
		err = admin.Sources().UpdateSourceWithURL(sourceData.SourceConf, sourceData.Archive, updateOptions)
	} else {
		err = admin.Sources().UpdateSource(sourceData.SourceConf, sourceData.Archive, updateOptions)
	}

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to update source '%s' in tenant '%s' namespace '%s': %v. Verify the source exists and all parameters are valid.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Updated source '%s' successfully in tenant '%s' namespace '%s'. The source may need to be restarted to apply all changes.",
		name, tenant, namespace)), nil
}

// handleSourceDelete handles deleting a source
func handleSourceDelete(ctx context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	err := admin.Sources().DeleteSource(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete source '%s' in tenant '%s' namespace '%s': %v. Verify the source exists and you have deletion permissions.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Deleted source '%s' successfully from tenant '%s' namespace '%s'. All running instances have been terminated.",
		name, tenant, namespace)), nil
}

// handleSourceStart handles starting a source
func handleSourceStart(ctx context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	err := admin.Sources().StartSource(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to start source '%s' in tenant '%s' namespace '%s': %v. Verify the source exists and is not already running.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Started source '%s' successfully in tenant '%s' namespace '%s'. The source will begin pulling data from the external system.",
		name, tenant, namespace)), nil
}

// handleSourceStop handles stopping a source
func handleSourceStop(ctx context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	err := admin.Sources().StopSource(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to stop source '%s' in tenant '%s' namespace '%s': %v. Verify the source exists and is currently running.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Stopped source '%s' successfully in tenant '%s' namespace '%s'. The source will no longer pull data until restarted.",
		name, tenant, namespace)), nil
}

// handleSourceRestart handles restarting a source
func handleSourceRestart(ctx context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	err := admin.Sources().RestartSource(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to restart source '%s' in tenant '%s' namespace '%s': %v. Verify the source exists and is properly deployed.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Restarted source '%s' successfully in tenant '%s' namespace '%s'. All source instances have been restarted.",
		name, tenant, namespace)), nil
}

// handleListBuiltInSources handles listing all built-in source connectors
func handleListBuiltInSources(ctx context.Context, admin cmdutils.Client) (*mcp.CallToolResult, error) {
	sources, err := admin.Sources().GetBuiltInSources()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list built-in sources: %v. There might be an issue connecting to the Pulsar cluster.", err)), nil
	}

	// Convert result to JSON string
	sourcesJSON, err := json.Marshal(sources)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize built-in sources: %v", err)), nil
	}

	return mcp.NewToolResultText(string(sourcesJSON)), nil
}

// processSourceArguments is a simplified version of the pulsarctl function to process source arguments
func processSourceArguments(sourceData *utils.SourceData) error {
	// Initialize config if needed
	if sourceData.SourceConf == nil {
		sourceData.SourceConf = new(utils.SourceConfig)
	}

	// Set basic config values
	sourceData.SourceConf.Tenant = sourceData.Tenant
	sourceData.SourceConf.Namespace = sourceData.Namespace
	sourceData.SourceConf.Name = sourceData.Name

	// Set destination topic if provided
	if sourceData.DestinationTopicName != "" {
		sourceData.SourceConf.TopicName = sourceData.DestinationTopicName
	}

	// Set deserialization class name if provided
	if sourceData.DeserializationClassName != "" {
		sourceData.SourceConf.SerdeClassName = sourceData.DeserializationClassName
	}

	// Set schema type if provided
	if sourceData.SchemaType != "" {
		sourceData.SourceConf.SchemaType = sourceData.SchemaType
	}

	// Set class name if provided
	if sourceData.ClassName != "" {
		sourceData.SourceConf.ClassName = sourceData.ClassName
	}

	// Set processing guarantees if provided
	if sourceData.ProcessingGuarantees != "" {
		sourceData.SourceConf.ProcessingGuarantees = sourceData.ProcessingGuarantees
	}

	// Set parallelism if provided
	if sourceData.Parallelism != 0 {
		sourceData.SourceConf.Parallelism = sourceData.Parallelism
	} else if sourceData.SourceConf.Parallelism <= 0 {
		sourceData.SourceConf.Parallelism = 1
	}

	// Handle archive and source-type
	if sourceData.Archive != "" && sourceData.SourceType != "" {
		return fmt.Errorf("cannot specify both archive and source-type")
	}

	if sourceData.Archive != "" {
		sourceData.SourceConf.Archive = sourceData.Archive
	}

	if sourceData.SourceType != "" {
		// In a real implementation, we would validate the source type here
		sourceData.SourceConf.Archive = sourceData.SourceType
	}

	// Parse source config if provided
	if sourceData.SourceConfigString != "" {
		var configs map[string]interface{}
		if err := json.Unmarshal([]byte(sourceData.SourceConfigString), &configs); err != nil {
			return fmt.Errorf("failed to parse source config: %v", err)
		}
		sourceData.SourceConf.Configs = configs
	}

	return nil
}
