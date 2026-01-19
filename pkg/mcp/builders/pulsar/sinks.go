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

type pulsarAdminSinksInput struct {
	Operation     string         `json:"operation"`
	Tenant        *string        `json:"tenant,omitempty"`
	Namespace     *string        `json:"namespace,omitempty"`
	Name          *string        `json:"name,omitempty"`
	Archive       *string        `json:"archive,omitempty"`
	SinkType      *string        `json:"sink-type,omitempty"`
	Inputs        []string       `json:"inputs,omitempty"`
	TopicsPattern *string        `json:"topics-pattern,omitempty"`
	SubsName      *string        `json:"subs-name,omitempty"`
	Parallelism   *float64       `json:"parallelism,omitempty"`
	SinkConfig    map[string]any `json:"sink-config,omitempty"`
}

const (
	pulsarAdminSinksOperationDesc = "Operation to perform. Available operations:\n" +
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
	pulsarAdminSinksTenantDesc = "The tenant name. Tenants are the primary organizational unit in Pulsar, " +
		"providing multi-tenancy and resource isolation. Sinks deployed within a tenant " +
		"inherit its permissions and resource quotas. " +
		"Required for all operations except 'list-built-in'."
	pulsarAdminSinksNamespaceDesc = "The namespace name. Namespaces are logical groupings of topics and sinks " +
		"within a tenant. They encapsulate configuration policies and access control. " +
		"Sinks in a namespace typically process topics within the same namespace. " +
		"Required for all operations except 'list-built-in'."
	pulsarAdminSinksNameDesc = "The sink name. Required for all operations except 'list' and 'list-built-in'. " +
		"Names should be descriptive of the sink's purpose and must be unique within a namespace. " +
		"Sink names are used in metrics, logs, and when addressing the sink via APIs."
	pulsarAdminSinksArchiveDesc = "Path to the archive file containing the sink code. Optional for 'create' and 'update' operations. " +
		"Can be a local path, NAR file, or a URL accessible to the Pulsar broker. " +
		"The archive should contain all dependencies for the sink connector. " +
		"Either archive or sink-type must be specified, but not both."
	pulsarAdminSinksSinkTypeDesc = "The built-in sink connector type to use. Optional for 'create' and 'update' operations. " +
		"Specifies which built-in connector to use, such as 'jdbc', 'elastic-search', 'kafka', etc. " +
		"Use 'list-built-in' operation to see available sink types. " +
		"Either sink-type or archive must be specified, but not both."
	pulsarAdminSinksInputsDesc = "The sink's input topics (array of strings). Optional for 'create' and 'update' operations. " +
		"Topics must be specified in the format 'persistent://tenant/namespace/topic'. " +
		"Sinks can consume from multiple topics, but they should have compatible schemas. " +
		"All input topics should exist before the sink is created. " +
		"Either inputs or topics-pattern must be specified."
	pulsarAdminSinksTopicsPatternDesc = "TopicsPattern to consume from list of topics that match the pattern. Optional for 'create' and 'update' operations. " +
		"Specified as a regular expression, e.g., 'persistent://tenant/namespace/prefix.*'. " +
		"This allows the sink to automatically consume from topics that match the pattern, " +
		"including topics created after the sink is deployed. " +
		"Either topics-pattern or inputs must be specified."
	pulsarAdminSinksSubsNameDesc = "Pulsar subscription name for input topic consumer. Optional for 'create' and 'update' operations. " +
		"Defines the subscription name used by the sink to consume from input topics. " +
		"If not specified, a default subscription name will be generated. " +
		"The subscription type used is Shared by default."
	pulsarAdminSinksParallelismDesc = "The parallelism factor of the sink. Optional for 'create' and 'update' operations. " +
		"Determines how many instances of the sink will run concurrently. " +
		"Higher values improve throughput but require more resources. " +
		"Default is 1 (single instance). Recommended to align with topic partition count " +
		"when consuming from partitioned topics."
	pulsarAdminSinksConfigDesc = "User-defined sink config key/values. Optional for 'create' and 'update' operations. " +
		"Provides configuration parameters specific to the sink connector being used. " +
		"For example, JDBC connection strings, Elasticsearch indices, S3 bucket details, etc. " +
		"Specify as a JSON object with configuration properties required by the specific sink type. " +
		"Example: {\"jdbcUrl\": \"jdbc:postgresql://localhost:5432/database\", \"tableName\": \"events\"}"
)

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
func (b *PulsarAdminSinksToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool, err := b.buildSinksTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildSinksHandler(config.ReadOnly)

	return []builders.ToolDefinition{
		builders.ServerTool[pulsarAdminSinksInput, any]{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildSinksTool builds the Pulsar admin sinks MCP tool definition
func (b *PulsarAdminSinksToolBuilder) buildSinksTool() (*sdk.Tool, error) {
	inputSchema, err := buildPulsarAdminSinksInputSchema()
	if err != nil {
		return nil, err
	}

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

	return &sdk.Tool{
		Name:        "pulsar_admin_sinks",
		Description: toolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildSinksHandler builds the Pulsar admin sinks handler function
func (b *PulsarAdminSinksToolBuilder) buildSinksHandler(readOnly bool) builders.ToolHandlerFunc[pulsarAdminSinksInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input pulsarAdminSinksInput) (*sdk.CallToolResult, any, error) {
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
			return nil, nil, fmt.Errorf("operation '%s' not allowed in read-only mode. read-only mode restricts modifications to Pulsar Sinks", operation)
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

		// List built-in sinks doesn't require tenant, namespace or name
		if operation == "list-built-in" {
			result, err := b.handleListBuiltInSinks(ctx, admin)
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
			result, err := b.handleSinkList(ctx, admin, tenant, namespace)
			return result, nil, err
		case "get":
			result, err := b.handleSinkGet(ctx, admin, tenant, namespace, name)
			return result, nil, err
		case "status":
			result, err := b.handleSinkStatus(ctx, admin, tenant, namespace, name)
			return result, nil, err
		case "create":
			result, err := b.handleSinkCreate(ctx, admin, input, tenant, namespace, name)
			return result, nil, err
		case "update":
			result, err := b.handleSinkUpdate(ctx, admin, input, tenant, namespace, name)
			return result, nil, err
		case "delete":
			result, err := b.handleSinkDelete(ctx, admin, tenant, namespace, name)
			return result, nil, err
		case "start":
			result, err := b.handleSinkStart(ctx, admin, tenant, namespace, name)
			return result, nil, err
		case "stop":
			result, err := b.handleSinkStop(ctx, admin, tenant, namespace, name)
			return result, nil, err
		case "restart":
			result, err := b.handleSinkRestart(ctx, admin, tenant, namespace, name)
			return result, nil, err
		default:
			return nil, nil, fmt.Errorf("unsupported operation: %s", operation)
		}
	}
}

// handleSinkList handles listing sinks
func (b *PulsarAdminSinksToolBuilder) handleSinkList(_ context.Context, admin cmdutils.Client, tenant, namespace string) (*sdk.CallToolResult, error) {
	sinks, err := admin.Sinks().ListSinks(tenant, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to list sinks in tenant '%s' namespace '%s': %v", tenant, namespace, err)
	}

	// Convert result to JSON string
	sinksJSON, err := json.Marshal(sinks)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize sinks list: %v", err)
	}

	return textResult(string(sinksJSON)), nil
}

// handleSinkGet handles getting a sink's details
func (b *PulsarAdminSinksToolBuilder) handleSinkGet(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	sink, err := admin.Sinks().GetSink(tenant, namespace, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get sink '%s' in tenant '%s' namespace '%s': %v. verify the sink exists and you have the correct permissions", name, tenant, namespace, err)
	}

	// Convert result to JSON string
	sinkJSON, err := json.Marshal(sink)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize sink details: %v", err)
	}

	return textResult(string(sinkJSON)), nil
}

// handleSinkStatus handles getting the status of a sink
func (b *PulsarAdminSinksToolBuilder) handleSinkStatus(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	status, err := admin.Sinks().GetSinkStatus(tenant, namespace, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get status for sink '%s' in tenant '%s' namespace '%s': %v. verify the sink exists and is properly deployed", name, tenant, namespace, err)
	}

	// Convert result to JSON string
	statusJSON, err := json.Marshal(status)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize sink status: %v", err)
	}

	return textResult(string(statusJSON)), nil
}

// handleSinkCreate handles creating a new sink
func (b *PulsarAdminSinksToolBuilder) handleSinkCreate(_ context.Context, admin cmdutils.Client, input pulsarAdminSinksInput, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	// Create a new SinkData object
	sinkData := &utils.SinkData{
		Tenant:    tenant,
		Namespace: namespace,
		Name:      name,
		SinkConf:  &utils.SinkConfig{},
	}

	// Get optional parameters
	if archive := stringValue(input.Archive); archive != "" {
		sinkData.Archive = archive
	}

	if sinkType := stringValue(input.SinkType); sinkType != "" {
		sinkData.SinkType = sinkType
	}

	if len(input.Inputs) > 0 {
		sinkData.Inputs = strings.Join(input.Inputs, ",")
	}

	if topicsPattern := stringValue(input.TopicsPattern); topicsPattern != "" {
		sinkData.TopicsPattern = topicsPattern
	}

	if subsName := stringValue(input.SubsName); subsName != "" {
		sinkData.SubsName = subsName
	}

	if input.Parallelism != nil && *input.Parallelism >= 0 {
		sinkData.Parallelism = int(*input.Parallelism)
	}

	// Get sink config if available
	if input.SinkConfig != nil {
		sinkConfigJSON, err := json.Marshal(input.SinkConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal sink-config: %v. ensure the sink configuration is a valid JSON object", err)
		}
		sinkData.SinkConfigString = string(sinkConfigJSON)
	}

	// Validate inputs
	if sinkData.Archive == "" && sinkData.SinkType == "" {
		return nil, fmt.Errorf("missing required parameter: either 'archive' or 'sink-type' must be specified for sink creation. use 'archive' for custom connectors or 'sink-type' for built-in connectors")
	}

	if sinkData.Archive != "" && sinkData.SinkType != "" {
		return nil, fmt.Errorf("invalid parameters: cannot specify both 'archive' and 'sink-type'. use only one of these parameters based on your connector type")
	}

	if sinkData.Inputs == "" && sinkData.TopicsPattern == "" {
		return nil, fmt.Errorf("missing required parameter: either 'inputs' or 'topics-pattern' must be specified. the sink needs a source of data to consume from Pulsar")
	}

	// Process the arguments
	err := b.processArguments(sinkData)
	if err != nil {
		return nil, fmt.Errorf("failed to process arguments: %v", err)
	}

	// Create the sink
	if sinkData.Archive != "" && b.isPackageURLSupported(sinkData.Archive) {
		err = admin.Sinks().CreateSinkWithURL(sinkData.SinkConf, sinkData.Archive)
	} else {
		err = admin.Sinks().CreateSink(sinkData.SinkConf, sinkData.Archive)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create sink '%s' in tenant '%s' namespace '%s': %v. verify all parameters are correct and required resources exist", name, tenant, namespace, err)
	}

	return textResult(fmt.Sprintf("Created sink '%s' successfully in tenant '%s' namespace '%s'. The sink will start consuming from its input topics and writing to the configured destination.", name, tenant, namespace)), nil
}

// handleSinkUpdate handles updating an existing sink
func (b *PulsarAdminSinksToolBuilder) handleSinkUpdate(_ context.Context, admin cmdutils.Client, input pulsarAdminSinksInput, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	// Create a new SinkData object
	sinkData := &utils.SinkData{
		Tenant:    tenant,
		Namespace: namespace,
		Name:      name,
		SinkConf:  &utils.SinkConfig{},
	}

	// Get optional parameters
	if archive := stringValue(input.Archive); archive != "" {
		sinkData.Archive = archive
	}

	if sinkType := stringValue(input.SinkType); sinkType != "" {
		sinkData.SinkType = sinkType
	}

	if len(input.Inputs) > 0 {
		sinkData.Inputs = strings.Join(input.Inputs, ",")
	}

	if topicsPattern := stringValue(input.TopicsPattern); topicsPattern != "" {
		sinkData.TopicsPattern = topicsPattern
	}

	if subsName := stringValue(input.SubsName); subsName != "" {
		sinkData.SubsName = subsName
	}

	if input.Parallelism != nil && *input.Parallelism >= 0 {
		sinkData.Parallelism = int(*input.Parallelism)
	}

	// Get sink config if available
	if input.SinkConfig != nil {
		sinkConfigJSON, err := json.Marshal(input.SinkConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal sink-config: %v. ensure the sink configuration is a valid JSON object", err)
		}
		sinkData.SinkConfigString = string(sinkConfigJSON)
	}

	// Validate inputs if both are specified
	if sinkData.Archive != "" && sinkData.SinkType != "" {
		return nil, fmt.Errorf("invalid parameters: cannot specify both 'archive' and 'sink-type'. use only one of these parameters based on your connector type")
	}

	// Process the arguments
	err := b.processArguments(sinkData)
	if err != nil {
		return nil, fmt.Errorf("failed to process arguments: %v", err)
	}

	// Create update options
	updateOptions := &utils.UpdateOptions{
		UpdateAuthData: true,
	}

	// Update the sink
	if sinkData.Archive != "" && b.isPackageURLSupported(sinkData.Archive) {
		err = admin.Sinks().UpdateSinkWithURL(sinkData.SinkConf, sinkData.Archive, updateOptions)
	} else {
		err = admin.Sinks().UpdateSink(sinkData.SinkConf, sinkData.Archive, updateOptions)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to update sink '%s' in tenant '%s' namespace '%s': %v. verify the sink exists and all parameters are valid", name, tenant, namespace, err)
	}

	return textResult(fmt.Sprintf("Updated sink '%s' successfully in tenant '%s' namespace '%s'. The sink may need to be restarted to apply all changes.", name, tenant, namespace)), nil
}

// handleSinkDelete handles deleting a sink
func (b *PulsarAdminSinksToolBuilder) handleSinkDelete(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	err := admin.Sinks().DeleteSink(tenant, namespace, name)
	if err != nil {
		return nil, fmt.Errorf("failed to delete sink '%s' in tenant '%s' namespace '%s': %v. verify the sink exists and you have deletion permissions", name, tenant, namespace, err)
	}

	return textResult(fmt.Sprintf("Deleted sink '%s' successfully from tenant '%s' namespace '%s'. All running instances have been terminated.", name, tenant, namespace)), nil
}

// handleSinkStart handles starting a sink
func (b *PulsarAdminSinksToolBuilder) handleSinkStart(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	err := admin.Sinks().StartSink(tenant, namespace, name)
	if err != nil {
		return nil, fmt.Errorf("failed to start sink '%s' in tenant '%s' namespace '%s': %v. verify the sink exists and is not already running", name, tenant, namespace, err)
	}

	return textResult(fmt.Sprintf("Started sink '%s' successfully in tenant '%s' namespace '%s'. The sink will begin consuming from its input topics.", name, tenant, namespace)), nil
}

// handleSinkStop handles stopping a sink
func (b *PulsarAdminSinksToolBuilder) handleSinkStop(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	err := admin.Sinks().StopSink(tenant, namespace, name)
	if err != nil {
		return nil, fmt.Errorf("failed to stop sink '%s' in tenant '%s' namespace '%s': %v. verify the sink exists and is currently running", name, tenant, namespace, err)
	}

	return textResult(fmt.Sprintf("Stopped sink '%s' successfully in tenant '%s' namespace '%s'. The sink will no longer consume data until restarted.", name, tenant, namespace)), nil
}

// handleSinkRestart handles restarting a sink
func (b *PulsarAdminSinksToolBuilder) handleSinkRestart(_ context.Context, admin cmdutils.Client, tenant, namespace, name string) (*sdk.CallToolResult, error) {
	err := admin.Sinks().RestartSink(tenant, namespace, name)
	if err != nil {
		return nil, fmt.Errorf("failed to restart sink '%s' in tenant '%s' namespace '%s': %v. verify the sink exists and is properly deployed", name, tenant, namespace, err)
	}

	return textResult(fmt.Sprintf("Restarted sink '%s' successfully in tenant '%s' namespace '%s'. All sink instances have been restarted.", name, tenant, namespace)), nil
}

// handleListBuiltInSinks handles listing all built-in sink connectors
func (b *PulsarAdminSinksToolBuilder) handleListBuiltInSinks(_ context.Context, admin cmdutils.Client) (*sdk.CallToolResult, error) {
	sinks, err := admin.Sinks().GetBuiltInSinks()
	if err != nil {
		return nil, fmt.Errorf("failed to list built-in sinks: %v. there might be an issue connecting to the Pulsar cluster", err)
	}

	// Convert result to JSON string
	sinksJSON, err := json.Marshal(sinks)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize built-in sinks: %v", err)
	}

	return textResult(string(sinksJSON)), nil
}

// processArguments is a simplified version of the pulsarctl function to process sink arguments
func (b *PulsarAdminSinksToolBuilder) processArguments(sinkData *utils.SinkData) error {
	// Initialize config if needed
	if sinkData.SinkConf == nil {
		sinkData.SinkConf = new(utils.SinkConfig)
	}

	// Set basic config values
	sinkData.SinkConf.Tenant = sinkData.Tenant
	sinkData.SinkConf.Namespace = sinkData.Namespace
	sinkData.SinkConf.Name = sinkData.Name
	if sinkData.Inputs != "" {
		sinkData.SinkConf.Inputs = strings.Split(sinkData.Inputs, ",")
	}

	if sinkData.SubsName != "" {
		sinkData.SinkConf.SourceSubscriptionName = sinkData.SubsName
	}

	if sinkData.TopicsPattern != "" {
		sinkData.SinkConf.TopicsPattern = &sinkData.TopicsPattern
	}

	if sinkData.Parallelism != 0 {
		sinkData.SinkConf.Parallelism = sinkData.Parallelism
	}

	if sinkData.Archive != "" {
		sinkData.SinkConf.Archive = sinkData.Archive
	}

	if sinkData.SinkType != "" {
		sinkData.SinkConf.Archive = sinkData.SinkType
	}

	if sinkData.SinkConfigString != "" {
		var configs map[string]interface{}
		if err := json.Unmarshal([]byte(sinkData.SinkConfigString), &configs); err != nil {
			return fmt.Errorf("failed to parse sink config: %v", err)
		}
		sinkData.SinkConf.Configs = configs
	}

	return nil
}

// isPackageURLSupported checks if the package URL is supported
// Validates URLs for Pulsar sink packages
func (b *PulsarAdminSinksToolBuilder) isPackageURLSupported(archive string) bool {
	if archive == "" {
		return false
	}

	// Check for supported URL schemes for Pulsar sink packages
	supportedSchemes := []string{
		"http://",
		"https://",
		"file://",
		"sink://", // Pulsar sink package URL
		"function://",
	}

	for _, scheme := range supportedSchemes {
		if strings.HasPrefix(archive, scheme) {
			return true
		}
	}

	// Also check if it's a local file path (not a URL)
	return !strings.Contains(archive, "://")
}

func buildPulsarAdminSinksInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[pulsarAdminSinksInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}

	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}

	setSchemaDescription(schema, "operation", pulsarAdminSinksOperationDesc)
	setSchemaDescription(schema, "tenant", pulsarAdminSinksTenantDesc)
	setSchemaDescription(schema, "namespace", pulsarAdminSinksNamespaceDesc)
	setSchemaDescription(schema, "name", pulsarAdminSinksNameDesc)
	setSchemaDescription(schema, "archive", pulsarAdminSinksArchiveDesc)
	setSchemaDescription(schema, "sink-type", pulsarAdminSinksSinkTypeDesc)
	setSchemaDescription(schema, "inputs", pulsarAdminSinksInputsDesc)
	setSchemaDescription(schema, "topics-pattern", pulsarAdminSinksTopicsPatternDesc)
	setSchemaDescription(schema, "subs-name", pulsarAdminSinksSubsNameDesc)
	setSchemaDescription(schema, "parallelism", pulsarAdminSinksParallelismDesc)
	setSchemaDescription(schema, "sink-config", pulsarAdminSinksConfigDesc)

	if inputsSchema := schema.Properties["inputs"]; inputsSchema != nil && inputsSchema.Items != nil {
		inputsSchema.Items.Description = "input topic"
	}

	normalizeAdditionalProperties(schema)
	return schema, nil
}
