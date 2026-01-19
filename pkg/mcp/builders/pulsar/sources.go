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
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

type pulsarAdminSourcesInput struct {
	Operation            string         `json:"operation"`
	Tenant               *string        `json:"tenant,omitempty"`
	Namespace            *string        `json:"namespace,omitempty"`
	Name                 *string        `json:"name,omitempty"`
	Archive              *string        `json:"archive,omitempty"`
	SourceType           *string        `json:"source-type,omitempty"`
	DestinationTopicName *string        `json:"destination-topic-name,omitempty"`
	DeserializationClass *string        `json:"deserialization-classname,omitempty"`
	SchemaType           *string        `json:"schema-type,omitempty"`
	ClassName            *string        `json:"classname,omitempty"`
	ProcessingGuarantees *string        `json:"processing-guarantees,omitempty"`
	Parallelism          *float64       `json:"parallelism,omitempty"`
	SourceConfig         map[string]any `json:"source-config,omitempty"`
}

const (
	pulsarAdminSourcesOperationDesc = "Operation to perform. Available operations:\n" +
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
	pulsarAdminSourcesTenantDesc = "The tenant name. Tenants are the primary organizational unit in Pulsar, " +
		"providing multi-tenancy and resource isolation. Sources deployed within a tenant " +
		"inherit its permissions and resource quotas. " +
		"Required for all operations except 'list-built-in'."
	pulsarAdminSourcesNamespaceDesc = "The namespace name. Namespaces are logical groupings of topics and sources " +
		"within a tenant. They encapsulate configuration policies and access control. " +
		"Sources in a namespace typically publish to topics within the same namespace. " +
		"Required for all operations except 'list-built-in'."
	pulsarAdminSourcesNameDesc = "The source name. Required for all operations except 'list' and 'list-built-in'. " +
		"Names should be descriptive of the source's purpose and must be unique within a namespace. " +
		"Source names are used in metrics, logs, and when addressing the source via APIs."
	pulsarAdminSourcesArchiveDesc = "Path to the archive file containing the source code. Optional for 'create' and 'update' operations. " +
		"Can be a local path, NAR file, or a URL accessible to the Pulsar broker. " +
		"The archive should contain all dependencies for the source connector. " +
		"Either archive or source-type must be specified, but not both."
	pulsarAdminSourcesSourceTypeDesc = "The built-in source connector type to use. Optional for 'create' and 'update' operations. " +
		"Specifies which built-in connector to use, such as 'kafka', 'jdbc', 'file', etc. " +
		"Use 'list-built-in' operation to see available source types. " +
		"Either source-type or archive must be specified, but not both."
	pulsarAdminSourcesDestinationTopicDesc = "The Pulsar topic to which data is published. Required for 'create' operation, optional for 'update'. " +
		"Specified in the format 'persistent://tenant/namespace/topic'. " +
		"This is the topic where the source will send the data it extracts from the external system. " +
		"The topic will be automatically created if it doesn't exist."
	pulsarAdminSourcesDeserializationClassDesc = "The SerDe (Serialization/Deserialization) classname for the source. Optional for 'create' and 'update'. " +
		"Specifies how to convert data from the external system into Pulsar messages. " +
		"Common SerDe classes include AvroSchema, JsonSchema, StringSchema, etc. " +
		"If not specified, the source will use the default SerDe for the connector type."
	pulsarAdminSourcesSchemaTypeDesc = "The schema type to be used to encode messages emitted from the source. Optional for 'create' and 'update'. " +
		"Available schema types include: 'avro', 'json', 'protobuf', 'string', etc. " +
		"Schema types ensure data compatibility and enable schema evolution. " +
		"The schema type should match the format of data being ingested."
	pulsarAdminSourcesClassNameDesc = "The source's class name if archive is a file-url-path (file://...). Optional for 'create' and 'update'. " +
		"This specifies the fully qualified class name that implements the source connector. " +
		"Only needed when using a custom source implementation in a JAR file. " +
		"Built-in connectors don't require this parameter."
	pulsarAdminSourcesProcessingGuaranteesDesc = "The processing guarantees (delivery semantics) applied to the source. Optional for 'create' and 'update'. " +
		"Available options: 'atleast_once', 'atmost_once', 'effectively_once'. " +
		"Controls how data is delivered in failure scenarios. " +
		"'atleast_once' is the most common and ensures no data loss but may have duplicates. " +
		"Default is 'atleast_once'."
	pulsarAdminSourcesParallelismDesc = "The parallelism factor of the source. Optional for 'create' and 'update' operations. " +
		"Determines how many instances of the source will run concurrently. " +
		"Higher values improve throughput but require more resources. " +
		"Default is 1 (single instance). Recommended to align with both source capacity " +
		"and destination topic partition count."
	pulsarAdminSourcesConfigDesc = "User-defined source config key/values. Optional for 'create' and 'update' operations. " +
		"Provides configuration parameters specific to the source connector being used. " +
		"For example, database connection details, Kafka bootstrap servers, credentials, etc. " +
		"Specify as a JSON object with configuration properties required by the specific source type. " +
		"Example: {\"topic\": \"external-kafka-topic\", \"bootstrapServers\": \"kafka:9092\"}"
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
func (b *PulsarAdminSourcesToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool, err := b.buildSourcesTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildSourcesHandler(config.ReadOnly)

	return []builders.ToolDefinition{
		builders.ServerTool[pulsarAdminSourcesInput, any]{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildSourcesTool builds the Pulsar admin sources MCP tool definition
func (b *PulsarAdminSourcesToolBuilder) buildSourcesTool() (*sdk.Tool, error) {
	inputSchema, err := buildPulsarAdminSourcesInputSchema()
	if err != nil {
		return nil, err
	}

	toolDesc := "Manage Apache Pulsar Sources for data ingestion and integration. " +
		"Pulsar Sources are connectors that import data from external systems into Pulsar topics. " +
		"Sources connect to external systems such as databases, messaging platforms, storage services, " +
		"and real-time data streams to pull data and publish it to Pulsar topics. " +
		"Built-in source connectors are available for common systems like Kafka, JDBC, AWS services, and more. " +
		"Sources follow the tenant/namespace/name hierarchy for organization and access control, " +
		"can scale through parallelism configuration, and support various processing guarantees. " +
		"This tool provides complete lifecycle management including deployment, configuration, " +
		"monitoring, and runtime control. Sources use schema types to ensure data compatibility."

	return &sdk.Tool{
		Name:        "pulsar_admin_sources",
		Description: toolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildSourcesHandler builds the Pulsar admin sources handler function
func (b *PulsarAdminSourcesToolBuilder) buildSourcesHandler(readOnly bool) builders.ToolHandlerFunc[pulsarAdminSourcesInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input pulsarAdminSourcesInput) (*sdk.CallToolResult, any, error) {
		// Extract and validate operation parameter
		operation := input.Operation
		if operation == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'operation'")
		}

		// Check if the operation is valid
		validOperations := map[string]bool{
			"list": true, "get": true, "status": true, "create": true, "update": true,
			"delete": true, "start": true, "stop": true, "restart": true, "list-built-in": true,
		}

		if !validOperations[operation] {
			return nil, nil, fmt.Errorf("invalid operation: '%s'. supported operations: list, get, status, create, update, delete, start, stop, restart, list-built-in", operation)
		}

		// Check write permissions for write operations
		writeOperations := map[string]bool{
			"create": true, "update": true, "delete": true, "start": true,
			"stop": true, "restart": true,
		}

		if readOnly && writeOperations[operation] {
			return nil, nil, fmt.Errorf("operation '%s' not allowed in read-only mode. read-only mode restricts modifications to Pulsar Sources", operation)
		}

		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return nil, nil, fmt.Errorf("pulsar session not found in context")
		}

		admin, err := session.GetAdminV3Client()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get Pulsar client: %v", err)
		}

		// List built-in sources doesn't require tenant, namespace or name
		if operation == "list-built-in" {
			result, err := b.handleListBuiltInSources(ctx, admin)
			return result, nil, err
		}

		// Extract common parameters (all operations except list-built-in require tenant and namespace)
		tenant, err := requireString(input.Tenant, "tenant")
		if err != nil {
			return nil, nil, fmt.Errorf("missing required parameter 'tenant': %v", err)
		}

		namespace, err := requireString(input.Namespace, "namespace")
		if err != nil {
			return nil, nil, fmt.Errorf("missing required parameter 'namespace': %v", err)
		}

		// name is required for all operations except list and list-built-in
		name := ""
		if operation != "list" {
			name, err = requireString(input.Name, "name")
			if err != nil {
				return nil, nil, fmt.Errorf("missing required parameter 'name': %v", err)
			}
		}

		// Dispatch based on operation
		switch operation {
		case "list":
			result, err := b.handleSourceList(ctx, admin, tenant, namespace)
			return result, nil, err
		case "get":
			result, err := b.handleSourceGet(ctx, admin, tenant, namespace, name)
			return result, nil, err
		case "status":
			result, err := b.handleSourceStatus(ctx, admin, tenant, namespace, name)
			return result, nil, err
		case "create":
			result, err := b.handleSourceCreate(ctx, admin, input, tenant, namespace, name)
			return result, nil, err
		case "update":
			result, err := b.handleSourceUpdate(ctx, admin, input, tenant, namespace, name)
			return result, nil, err
		case "delete":
			result, err := b.handleSourceDelete(ctx, admin, tenant, namespace, name)
			return result, nil, err
		case "start":
			result, err := b.handleSourceStart(ctx, admin, tenant, namespace, name)
			return result, nil, err
		case "stop":
			result, err := b.handleSourceStop(ctx, admin, tenant, namespace, name)
			return result, nil, err
		case "restart":
			result, err := b.handleSourceRestart(ctx, admin, tenant, namespace, name)
			return result, nil, err
		default:
			return nil, nil, fmt.Errorf("unsupported operation: %s", operation)
		}
	}
}

// handleSourceList handles listing sources
func (b *PulsarAdminSourcesToolBuilder) handleSourceList(_ context.Context, admin cmdutils.Client, tenant, namespace string) (*sdk.CallToolResult, error) {
	sources, err := admin.Sources().ListSources(tenant, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to list sources in tenant '%s' namespace '%s': %v", tenant, namespace, err)
	}

	// Convert result to JSON string
	sourcesJSON, err := json.Marshal(sources)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize sources list: %v", err)
	}

	return textResult(string(sourcesJSON)), nil
}

// handleSourceGet handles getting a source's details
func (b *PulsarAdminSourcesToolBuilder) handleSourceGet(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	source, err := admin.Sources().GetSource(tenant, namespace, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get source '%s' in tenant '%s' namespace '%s': %v. verify the source exists and you have the correct permissions", name, tenant, namespace, err)
	}

	// Convert result to JSON string
	sourceJSON, err := json.Marshal(source)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize source details: %v", err)
	}

	return textResult(string(sourceJSON)), nil
}

// handleSourceStatus handles getting the status of a source
func (b *PulsarAdminSourcesToolBuilder) handleSourceStatus(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	status, err := admin.Sources().GetSourceStatus(tenant, namespace, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get status for source '%s' in tenant '%s' namespace '%s': %v. verify the source exists and is properly deployed", name, tenant, namespace, err)
	}

	// Convert result to JSON string
	statusJSON, err := json.Marshal(status)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize source status: %v", err)
	}

	return textResult(string(statusJSON)), nil
}

// handleSourceCreate handles creating a new source
func (b *PulsarAdminSourcesToolBuilder) handleSourceCreate(_ context.Context, admin cmdutils.Client, input pulsarAdminSourcesInput, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	// Create a new SourceData object
	sourceData := &utils.SourceData{
		Tenant:     tenant,
		Namespace:  namespace,
		Name:       name,
		SourceConf: &utils.SourceConfig{},
	}

	// Get optional parameters
	if archive := stringValue(input.Archive); archive != "" {
		sourceData.Archive = archive
	}

	if sourceType := stringValue(input.SourceType); sourceType != "" {
		sourceData.SourceType = sourceType
	}

	if destTopic := stringValue(input.DestinationTopicName); destTopic != "" {
		sourceData.DestinationTopicName = destTopic
	}

	if deserializationClassName := stringValue(input.DeserializationClass); deserializationClassName != "" {
		sourceData.DeserializationClassName = deserializationClassName
	}

	if schemaType := stringValue(input.SchemaType); schemaType != "" {
		sourceData.SchemaType = schemaType
	}

	if className := stringValue(input.ClassName); className != "" {
		sourceData.ClassName = className
	}

	if processingGuarantees := stringValue(input.ProcessingGuarantees); processingGuarantees != "" {
		sourceData.ProcessingGuarantees = processingGuarantees
	}

	if input.Parallelism != nil && *input.Parallelism >= 0 {
		sourceData.Parallelism = int(*input.Parallelism)
	}

	// Get source config if available
	if input.SourceConfig != nil {
		sourceConfigJSON, err := json.Marshal(input.SourceConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal source-config: %v. ensure the source configuration is a valid JSON object", err)
		}
		sourceData.SourceConfigString = string(sourceConfigJSON)
	}

	// Validate inputs
	if sourceData.Archive == "" && sourceData.SourceType == "" {
		return nil, fmt.Errorf("missing required parameter: either 'archive' or 'source-type' must be specified for source creation. use 'archive' for custom connectors or 'source-type' for built-in connectors")
	}

	if sourceData.Archive != "" && sourceData.SourceType != "" {
		return nil, fmt.Errorf("invalid parameters: cannot specify both 'archive' and 'source-type'. use only one of these parameters based on your connector type")
	}

	if sourceData.DestinationTopicName == "" {
		return nil, fmt.Errorf("missing required parameter: 'destination-topic-name' must be specified. this is the Pulsar topic where the source will publish data")
	}

	// Process the arguments
	err := b.processSourceArguments(sourceData)
	if err != nil {
		return nil, fmt.Errorf("failed to process arguments: %v", err)
	}

	// Create the source
	if sourceData.Archive != "" && b.isPackageURLSupported(sourceData.Archive) {
		err = admin.Sources().CreateSourceWithURL(sourceData.SourceConf, sourceData.Archive)
	} else {
		err = admin.Sources().CreateSource(sourceData.SourceConf, sourceData.Archive)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create source '%s' in tenant '%s' namespace '%s': %v. verify all parameters are correct and required resources exist", name, tenant, namespace, err)
	}

	return textResult(fmt.Sprintf("Created source '%s' successfully in tenant '%s' namespace '%s'. The source will start pulling data from the external system and publishing to the destination topic.", name, tenant, namespace)), nil
}

// handleSourceUpdate handles updating an existing source
func (b *PulsarAdminSourcesToolBuilder) handleSourceUpdate(_ context.Context, admin cmdutils.Client, input pulsarAdminSourcesInput, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	// Create a new SourceData object
	sourceData := &utils.SourceData{
		Tenant:     tenant,
		Namespace:  namespace,
		Name:       name,
		SourceConf: &utils.SourceConfig{},
	}

	// Get optional parameters
	if archive := stringValue(input.Archive); archive != "" {
		sourceData.Archive = archive
	}

	if sourceType := stringValue(input.SourceType); sourceType != "" {
		sourceData.SourceType = sourceType
	}

	if destTopic := stringValue(input.DestinationTopicName); destTopic != "" {
		sourceData.DestinationTopicName = destTopic
	}

	if deserializationClassName := stringValue(input.DeserializationClass); deserializationClassName != "" {
		sourceData.DeserializationClassName = deserializationClassName
	}

	if schemaType := stringValue(input.SchemaType); schemaType != "" {
		sourceData.SchemaType = schemaType
	}

	if className := stringValue(input.ClassName); className != "" {
		sourceData.ClassName = className
	}

	if processingGuarantees := stringValue(input.ProcessingGuarantees); processingGuarantees != "" {
		sourceData.ProcessingGuarantees = processingGuarantees
	}

	if input.Parallelism != nil && *input.Parallelism >= 0 {
		sourceData.Parallelism = int(*input.Parallelism)
	}

	// Get source config if available
	if input.SourceConfig != nil {
		sourceConfigJSON, err := json.Marshal(input.SourceConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal source-config: %v. ensure the source configuration is a valid JSON object", err)
		}
		sourceData.SourceConfigString = string(sourceConfigJSON)
	}

	// Validate inputs if both are specified
	if sourceData.Archive != "" && sourceData.SourceType != "" {
		return nil, fmt.Errorf("invalid parameters: cannot specify both 'archive' and 'source-type'. use only one of these parameters based on your connector type")
	}

	// Process the arguments
	err := b.processSourceArguments(sourceData)
	if err != nil {
		return nil, fmt.Errorf("failed to process arguments: %v", err)
	}

	// Create update options
	updateOptions := &utils.UpdateOptions{
		UpdateAuthData: true,
	}

	// Update the source
	if sourceData.Archive != "" && b.isPackageURLSupported(sourceData.Archive) {
		err = admin.Sources().UpdateSourceWithURL(sourceData.SourceConf, sourceData.Archive, updateOptions)
	} else {
		err = admin.Sources().UpdateSource(sourceData.SourceConf, sourceData.Archive, updateOptions)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to update source '%s' in tenant '%s' namespace '%s': %v. verify the source exists and all parameters are valid", name, tenant, namespace, err)
	}

	return textResult(fmt.Sprintf("Updated source '%s' successfully in tenant '%s' namespace '%s'. The source may need to be restarted to apply all changes.", name, tenant, namespace)), nil
}

// handleSourceDelete handles deleting a source
func (b *PulsarAdminSourcesToolBuilder) handleSourceDelete(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	err := admin.Sources().DeleteSource(tenant, namespace, name)
	if err != nil {
		return nil, fmt.Errorf("failed to delete source '%s' in tenant '%s' namespace '%s': %v. verify the source exists and you have deletion permissions", name, tenant, namespace, err)
	}

	return textResult(fmt.Sprintf("Deleted source '%s' successfully from tenant '%s' namespace '%s'. All running instances have been terminated.", name, tenant, namespace)), nil
}

// handleSourceStart handles starting a source
func (b *PulsarAdminSourcesToolBuilder) handleSourceStart(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	err := admin.Sources().StartSource(tenant, namespace, name)
	if err != nil {
		return nil, fmt.Errorf("failed to start source '%s' in tenant '%s' namespace '%s': %v. verify the source exists and is not already running", name, tenant, namespace, err)
	}

	return textResult(fmt.Sprintf("Started source '%s' successfully in tenant '%s' namespace '%s'. The source will begin pulling data from the external system.", name, tenant, namespace)), nil
}

// handleSourceStop handles stopping a source
func (b *PulsarAdminSourcesToolBuilder) handleSourceStop(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	err := admin.Sources().StopSource(tenant, namespace, name)
	if err != nil {
		return nil, fmt.Errorf("failed to stop source '%s' in tenant '%s' namespace '%s': %v. verify the source exists and is currently running", name, tenant, namespace, err)
	}

	return textResult(fmt.Sprintf("Stopped source '%s' successfully in tenant '%s' namespace '%s'. The source will no longer pull data until restarted.", name, tenant, namespace)), nil
}

// handleSourceRestart handles restarting a source
func (b *PulsarAdminSourcesToolBuilder) handleSourceRestart(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	err := admin.Sources().RestartSource(tenant, namespace, name)
	if err != nil {
		return nil, fmt.Errorf("failed to restart source '%s' in tenant '%s' namespace '%s': %v. verify the source exists and is properly deployed", name, tenant, namespace, err)
	}

	return textResult(fmt.Sprintf("Restarted source '%s' successfully in tenant '%s' namespace '%s'. All source instances have been restarted.", name, tenant, namespace)), nil
}

// handleListBuiltInSources handles listing all built-in source connectors
func (b *PulsarAdminSourcesToolBuilder) handleListBuiltInSources(_ context.Context, admin cmdutils.Client) (*sdk.CallToolResult, error) {
	sources, err := admin.Sources().GetBuiltInSources()
	if err != nil {
		return nil, fmt.Errorf("failed to list built-in sources: %v. there might be an issue connecting to the Pulsar cluster", err)
	}

	// Convert result to JSON string
	sourcesJSON, err := json.Marshal(sources)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize built-in sources: %v", err)
	}

	return textResult(string(sourcesJSON)), nil
}

// processSourceArguments is a simplified version of the pulsarctl function to process source arguments
func (b *PulsarAdminSourcesToolBuilder) processSourceArguments(sourceData *utils.SourceData) error {
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

// isPackageURLSupported checks if the package URL is supported
// Validates URLs for Pulsar source packages
func (b *PulsarAdminSourcesToolBuilder) isPackageURLSupported(archive string) bool {
	if archive == "" {
		return false
	}

	// Check for supported URL schemes for Pulsar source packages
	supportedSchemes := []string{
		"http://",
		"https://",
		"file://",
		"function://", // Pulsar function package URL
		"source://",   // Pulsar source package URL
	}

	for _, scheme := range supportedSchemes {
		if strings.HasPrefix(archive, scheme) {
			return true
		}
	}

	// Also check if it's a local file path (not a URL)
	return !strings.Contains(archive, "://")
}

func buildPulsarAdminSourcesInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[pulsarAdminSourcesInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}

	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}

	setSchemaDescription(schema, "operation", pulsarAdminSourcesOperationDesc)
	setSchemaDescription(schema, "tenant", pulsarAdminSourcesTenantDesc)
	setSchemaDescription(schema, "namespace", pulsarAdminSourcesNamespaceDesc)
	setSchemaDescription(schema, "name", pulsarAdminSourcesNameDesc)
	setSchemaDescription(schema, "archive", pulsarAdminSourcesArchiveDesc)
	setSchemaDescription(schema, "source-type", pulsarAdminSourcesSourceTypeDesc)
	setSchemaDescription(schema, "destination-topic-name", pulsarAdminSourcesDestinationTopicDesc)
	setSchemaDescription(schema, "deserialization-classname", pulsarAdminSourcesDeserializationClassDesc)
	setSchemaDescription(schema, "schema-type", pulsarAdminSourcesSchemaTypeDesc)
	setSchemaDescription(schema, "classname", pulsarAdminSourcesClassNameDesc)
	setSchemaDescription(schema, "processing-guarantees", pulsarAdminSourcesProcessingGuaranteesDesc)
	setSchemaDescription(schema, "parallelism", pulsarAdminSourcesParallelismDesc)
	setSchemaDescription(schema, "source-config", pulsarAdminSourcesConfigDesc)

	normalizeAdditionalProperties(schema)
	return schema, nil
}
