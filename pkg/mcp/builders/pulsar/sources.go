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

// PulsarAdminSourcesToolBuilder implements the ToolBuilder interface for Pulsar admin sources
// /nolint:revive
type PulsarAdminSourcesToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminSourcesToolBuilder creates a new Pulsar admin sources tool builder instance
func NewPulsarAdminSourcesToolBuilder() *PulsarAdminSourcesToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_sources",
		Version:     "1.0.0",
		Description: "Pulsar admin sources management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "admin", "sources"},
	}

	features := []string{
		"pulsar-admin-sources",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminSourcesToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin sources tool list
func (b *PulsarAdminSourcesToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool := b.buildSourcesTool()
	handler := b.buildSourcesHandler(config.ReadOnly)

	return []server.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildSourcesTool builds the Pulsar admin sources MCP tool definition
func (b *PulsarAdminSourcesToolBuilder) buildSourcesTool() mcp.Tool {
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

	return mcp.NewTool("pulsar_admin_sources",
		mcp.WithDescription(toolDesc),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc)),
		mcp.WithString("tenant",
			mcp.Description("The tenant name. Tenants are the primary organizational unit in Pulsar, "+
				"providing multi-tenancy and resource isolation. Sources deployed within a tenant "+
				"inherit its permissions and resource quotas. "+
				"Defaults to 'public' if not provided.")),
		mcp.WithString("namespace",
			mcp.Description("The namespace name. Namespaces are logical groupings of topics and sources "+
				"within a tenant. They encapsulate configuration policies and access control. "+
				"Sources in a namespace typically publish to topics within the same namespace. "+
				"Defaults to 'default' if not provided.")),
		mcp.WithString("name",
			mcp.Description("The source name. Required for all operations except 'list' and 'list-built-in'. "+
				"Can be provided via source-config-file for create/update. "+
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
		mcp.WithString("source-config-file",
			mcp.Description("Path to a YAML source configuration file. Optional for 'create' and 'update'. "+
				"When provided, the file is loaded before applying explicit parameters.")),
		mcp.WithString("destination-topic-name",
			mcp.Description("The Pulsar topic to which data is published. Required for 'create' operation, optional for 'update'. "+
				"Specified in the format 'persistent://tenant/namespace/topic'. "+
				"This is the topic where the source will send the data it extracts from the external system. "+
				"The topic will be automatically created if it doesn't exist.")),
		mcp.WithObject("producer-config",
			mcp.Description("Custom producer configuration as JSON object. Optional for 'create' and 'update'.")),
		mcp.WithString("batch-builder",
			mcp.Description("Batch builder type (DEFAULT or KEY_BASED). Optional for 'create' and 'update'.")),
		mcp.WithObject("batch-source-config",
			mcp.Description("Batch source configuration as JSON object. Optional for 'create' and 'update'.")),
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
		mcp.WithNumber("cpu",
			mcp.Description("CPU cores allocated per source instance. Optional for 'create' and 'update'. "+
				"Applicable to process and container runtimes.")),
		mcp.WithNumber("ram",
			mcp.Description("RAM bytes allocated per source instance. Optional for 'create' and 'update'. "+
				"Applicable to process and container runtimes.")),
		mcp.WithNumber("disk",
			mcp.Description("Disk bytes allocated per source instance. Optional for 'create' and 'update'. "+
				"Applicable to process and container runtimes.")),
		mcp.WithString("custom-runtime-options",
			mcp.Description("Runtime customization options. Optional for 'create' and 'update'.")),
		mcp.WithObject("secrets",
			mcp.Description("Secrets configuration map. Optional for 'create' and 'update'. "+
				"Specify as a JSON object where values describe how secrets are fetched.")),
		mcp.WithObject("source-config",
			mcp.Description("User-defined source config key/values. Optional for 'create' and 'update' operations. "+
				"Provides configuration parameters specific to the source connector being used. "+
				"For example, database connection details, Kafka bootstrap servers, credentials, etc. "+
				"Specify as a JSON object with configuration properties required by the specific source type. "+
				"Example: {\"topic\": \"external-kafka-topic\", \"bootstrapServers\": \"kafka:9092\"}")),
		mcp.WithBoolean("update-auth-data",
			mcp.Description("Whether to update authentication data during source update. Optional for 'update' only.")),
	)
}

// buildSourcesHandler builds the Pulsar admin sources handler function
func (b *PulsarAdminSourcesToolBuilder) buildSourcesHandler(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return mcp.NewToolResultError("Pulsar session not found in context"), nil
		}

		admin, err := session.GetAdminV3Client()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get Pulsar client: %v", err)), nil
		}

		// List built-in sources doesn't require tenant, namespace or name
		if operation == "list-built-in" {
			return b.handleListBuiltInSources(ctx, admin)
		}

		args := request.GetArguments()

		if operation == "create" || operation == "update" {
			name, _ := getStringArg(args, "name")
			tenant, _ := getStringArg(args, "tenant")
			namespace, _ := getStringArg(args, "namespace")
			if operation == "create" {
				return b.handleSourceCreate(ctx, admin, tenant, namespace, name, request)
			}
			return b.handleSourceUpdate(ctx, admin, tenant, namespace, name, request)
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
				return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'name' for operation '%s': %v. The source name must be specified for this operation.", operation, err)), nil
			}
		}

		// Handle operations
		switch operation {
		case "list":
			return b.handleSourceList(ctx, admin, tenant, namespace)
		case "get":
			return b.handleSourceGet(ctx, admin, tenant, namespace, name)
		case "status":
			return b.handleSourceStatus(ctx, admin, tenant, namespace, name)
		case "delete":
			return b.handleSourceDelete(ctx, admin, tenant, namespace, name)
		case "start":
			return b.handleSourceStart(ctx, admin, tenant, namespace, name)
		case "stop":
			return b.handleSourceStop(ctx, admin, tenant, namespace, name)
		case "restart":
			return b.handleSourceRestart(ctx, admin, tenant, namespace, name)
		default:
			// This should never happen due to the valid operations check above
			return mcp.NewToolResultError(fmt.Sprintf("Unsupported operation: %s", operation)), nil
		}
	}
}

// Helper functions

// handleSourceList handles listing all sources under a namespace
func (b *PulsarAdminSourcesToolBuilder) handleSourceList(_ context.Context, admin cmdutils.Client, tenant, namespace string) (*mcp.CallToolResult, error) {
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
func (b *PulsarAdminSourcesToolBuilder) handleSourceGet(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
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
func (b *PulsarAdminSourcesToolBuilder) handleSourceStatus(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
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
func (b *PulsarAdminSourcesToolBuilder) handleSourceCreate(_ context.Context, admin cmdutils.Client, tenant, namespace, name string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if _, err := request.RequireString("destination-topic-name"); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'destination-topic-name': %v. This is the Pulsar topic where the source will publish data.", err)), nil
	}

	config, archiveArg, sourceTypeArg, err := b.buildSourceConfig(tenant, namespace, name, request)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build source configuration for '%s': %v. Please verify all required parameters are provided correctly.", name, err)), nil
	}

	if err := validateSourceArchiveArgs(archiveArg, sourceTypeArg); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid parameters: %v", err)), nil
	}

	uploadArchive, err := b.resolveSourceArchive(admin, config, archiveArg, sourceTypeArg)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to resolve source archive: %v", err)), nil
	}

	b.applySourceDefaults(config)
	if config.Name == "" {
		return mcp.NewToolResultError("Source name not specified. Provide the 'name' parameter or set it in the source-config-file."), nil
	}
	if config.TopicName == "" {
		return mcp.NewToolResultError("Missing required parameter: 'destination-topic-name' must be specified. This is the Pulsar topic where the source will publish data."), nil
	}
	if err := validateSourceConfig(config); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to validate source configuration for '%s': %v.", config.Name, err)), nil
	}

	if uploadArchive != "" && b.isPackageURLSupported(uploadArchive) {
		err = admin.Sources().CreateSourceWithURL(config, uploadArchive)
	} else {
		err = admin.Sources().CreateSource(config, uploadArchive)
	}

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create source '%s' in tenant '%s' namespace '%s': %v. Verify all parameters are correct and required resources exist.",
			config.Name, config.Tenant, config.Namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Created source '%s' successfully in tenant '%s' namespace '%s'. The source will start pulling data from the external system and publishing to the destination topic.",
		config.Name, config.Tenant, config.Namespace)), nil
}

// handleSourceUpdate handles updating an existing source
func (b *PulsarAdminSourcesToolBuilder) handleSourceUpdate(_ context.Context, admin cmdutils.Client, tenant, namespace, name string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	config, archiveArg, sourceTypeArg, err := b.buildSourceConfig(tenant, namespace, name, request)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build source configuration for '%s': %v. Please verify all parameters are provided correctly.", name, err)), nil
	}

	if err := validateSourceArchiveArgs(archiveArg, sourceTypeArg); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid parameters: %v", err)), nil
	}

	uploadArchive, err := b.resolveSourceArchive(admin, config, archiveArg, sourceTypeArg)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to resolve source archive: %v", err)), nil
	}

	b.applySourceDefaults(config)
	if config.Name == "" {
		return mcp.NewToolResultError("Source name not specified. Provide the 'name' parameter or set it in the source-config-file."), nil
	}

	updateOptions := utils.NewUpdateOptions()
	if updateAuthData, ok := getBoolArg(request.GetArguments(), "update-auth-data"); ok {
		updateOptions.UpdateAuthData = updateAuthData
	}

	if uploadArchive != "" && b.isPackageURLSupported(uploadArchive) {
		err = admin.Sources().UpdateSourceWithURL(config, uploadArchive, updateOptions)
	} else {
		err = admin.Sources().UpdateSource(config, uploadArchive, updateOptions)
	}

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to update source '%s' in tenant '%s' namespace '%s': %v. Verify the source exists and all parameters are valid.",
			config.Name, config.Tenant, config.Namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Updated source '%s' successfully in tenant '%s' namespace '%s'. The source may need to be restarted to apply all changes.",
		config.Name, config.Tenant, config.Namespace)), nil
}

// handleSourceDelete handles deleting a source
func (b *PulsarAdminSourcesToolBuilder) handleSourceDelete(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	err := admin.Sources().DeleteSource(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete source '%s' in tenant '%s' namespace '%s': %v. Verify the source exists and you have deletion permissions.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Deleted source '%s' successfully from tenant '%s' namespace '%s'. All running instances have been terminated.",
		name, tenant, namespace)), nil
}

// handleSourceStart handles starting a source
func (b *PulsarAdminSourcesToolBuilder) handleSourceStart(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	err := admin.Sources().StartSource(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to start source '%s' in tenant '%s' namespace '%s': %v. Verify the source exists and is not already running.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Started source '%s' successfully in tenant '%s' namespace '%s'. The source will begin pulling data from the external system.",
		name, tenant, namespace)), nil
}

// handleSourceStop handles stopping a source
func (b *PulsarAdminSourcesToolBuilder) handleSourceStop(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	err := admin.Sources().StopSource(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to stop source '%s' in tenant '%s' namespace '%s': %v. Verify the source exists and is currently running.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Stopped source '%s' successfully in tenant '%s' namespace '%s'. The source will no longer pull data until restarted.",
		name, tenant, namespace)), nil
}

// handleSourceRestart handles restarting a source
func (b *PulsarAdminSourcesToolBuilder) handleSourceRestart(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*mcp.CallToolResult, error) {
	err := admin.Sources().RestartSource(tenant, namespace, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to restart source '%s' in tenant '%s' namespace '%s': %v. Verify the source exists and is properly deployed.",
			name, tenant, namespace, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Restarted source '%s' successfully in tenant '%s' namespace '%s'. All source instances have been restarted.",
		name, tenant, namespace)), nil
}

// handleListBuiltInSources handles listing all built-in source connectors
func (b *PulsarAdminSourcesToolBuilder) handleListBuiltInSources(_ context.Context, admin cmdutils.Client) (*mcp.CallToolResult, error) {
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

func (b *PulsarAdminSourcesToolBuilder) buildSourceConfig(tenant, namespace, name string, request mcp.CallToolRequest) (*utils.SourceConfig, string, string, error) {
	config := &utils.SourceConfig{}
	args := request.GetArguments()

	if configFile, ok := getStringArg(args, "source-config-file"); ok && configFile != "" {
		//nolint:gosec
		data, err := os.ReadFile(configFile)
		if err != nil {
			return nil, "", "", fmt.Errorf("load source config file failed: %w", err)
		}
		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, "", "", fmt.Errorf("unmarshal source config file error: %w", err)
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

	if destinationTopic, ok := getStringArg(args, "destination-topic-name"); ok && destinationTopic != "" {
		config.TopicName = destinationTopic
	}

	if deserializationClassName, ok := getStringArg(args, "deserialization-classname"); ok && deserializationClassName != "" {
		config.SerdeClassName = deserializationClassName
	}

	if schemaType, ok := getStringArg(args, "schema-type"); ok && schemaType != "" {
		config.SchemaType = schemaType
	}

	if className, ok := getStringArg(args, "classname"); ok && className != "" {
		config.ClassName = className
	}

	if processingGuarantees, ok := getStringArg(args, "processing-guarantees"); ok && processingGuarantees != "" {
		config.ProcessingGuarantees = processingGuarantees
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

	if sourceConfigValue, exists := args["source-config"]; exists && sourceConfigValue != nil {
		sourceConfig, err := decodeInterfaceMap(sourceConfigValue)
		if err != nil {
			return nil, "", "", fmt.Errorf("invalid source-config: %w", err)
		}
		config.Configs = sourceConfig
	}

	if producerConfigValue, exists := args["producer-config"]; exists && producerConfigValue != nil {
		producerConfig, err := decodeProducerConfig(producerConfigValue)
		if err != nil {
			return nil, "", "", fmt.Errorf("invalid producer-config: %w", err)
		}
		config.ProducerConfig = producerConfig
	}

	if batchBuilder, ok := getStringArg(args, "batch-builder"); ok && batchBuilder != "" {
		config.BatchBuilder = batchBuilder
	}

	if batchSourceConfigValue, exists := args["batch-source-config"]; exists && batchSourceConfigValue != nil {
		batchSourceConfig, err := decodeBatchSourceConfig(batchSourceConfigValue)
		if err != nil {
			return nil, "", "", fmt.Errorf("invalid batch-source-config: %w", err)
		}
		config.BatchSourceConfig = batchSourceConfig
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

	b.normalizeSourceConfigMaps(config)

	archiveArg, _ := getStringArg(args, "archive")
	sourceTypeArg, _ := getStringArg(args, "source-type")

	return config, archiveArg, sourceTypeArg, nil
}

func (b *PulsarAdminSourcesToolBuilder) normalizeSourceConfigMaps(config *utils.SourceConfig) {
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

func (b *PulsarAdminSourcesToolBuilder) applySourceDefaults(config *utils.SourceConfig) {
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

func validateSourceConfig(config *utils.SourceConfig) error {
	if config.Archive == "" {
		return fmt.Errorf("source archive not specified")
	}
	if config.Name == "" {
		return fmt.Errorf("source name not specified")
	}
	return nil
}

func validateSourceArchiveArgs(archive, sourceType string) error {
	if archive != "" && sourceType != "" {
		return fmt.Errorf("cannot specify both 'archive' and 'source-type'")
	}
	return nil
}

func (b *PulsarAdminSourcesToolBuilder) resolveSourceArchive(admin cmdutils.Client, config *utils.SourceConfig, archive, sourceType string) (string, error) {
	if sourceType != "" {
		resolved, err := b.validateSourceType(admin, sourceType)
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

	if config.Archive != "" {
		if strings.HasPrefix(config.Archive, "builtin://") {
			return "", nil
		}
		return config.Archive, nil
	}

	return "", nil
}

func (b *PulsarAdminSourcesToolBuilder) validateSourceType(admin cmdutils.Client, sourceType string) (string, error) {
	builtins, err := admin.Sources().GetBuiltInSources()
	if err != nil {
		return "", fmt.Errorf("failed to list built-in sources: %w", err)
	}

	names := make([]string, 0, len(builtins))
	for _, builtin := range builtins {
		names = append(names, builtin.Name)
		if builtin.Name == sourceType {
			return "builtin://" + sourceType, nil
		}
	}

	sort.Strings(names)
	return "", fmt.Errorf("invalid source-type %q. Available sources: %s", sourceType, strings.Join(names, ", "))
}

func decodeProducerConfig(value interface{}) (*utils.ProducerConfig, error) {
	var config utils.ProducerConfig
	if err := decodeInto(value, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func decodeBatchSourceConfig(value interface{}) (*utils.BatchSourceConfig, error) {
	var config utils.BatchSourceConfig
	if err := decodeInto(value, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// isPackageURLSupported checks if the package URL is supported
// Validates URLs for Pulsar source packages
func (b *PulsarAdminSourcesToolBuilder) isPackageURLSupported(archive string) bool {
	return archive != "" && (strings.HasPrefix(archive, "http") ||
		strings.HasPrefix(archive, "file") ||
		strings.HasPrefix(archive, "function") ||
		strings.HasPrefix(archive, "sink") ||
		strings.HasPrefix(archive, "source"))
}
