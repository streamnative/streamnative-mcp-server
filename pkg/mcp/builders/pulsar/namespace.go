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
	"strconv"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

type pulsarAdminNamespaceInput struct {
	Operation    string   `json:"operation"`
	Tenant       *string  `json:"tenant,omitempty"`
	Namespace    *string  `json:"namespace,omitempty"`
	Bundles      *string  `json:"bundles,omitempty"`
	Clusters     []string `json:"clusters,omitempty"`
	Subscription *string  `json:"subscription,omitempty"`
	Bundle       *string  `json:"bundle,omitempty"`
	Force        *string  `json:"force,omitempty"`
	Unload       *string  `json:"unload,omitempty"`
}

const (
	pulsarAdminNamespaceOperationDesc = "Operation to perform on namespaces. Available operations:\n" +
		"- list: List all namespaces for a tenant\n" +
		"- get_topics: Get all topics within a namespace\n" +
		"- create: Create a new namespace\n" +
		"- delete: Delete an existing namespace\n" +
		"- clear_backlog: Clear backlog for all topics in a namespace\n" +
		"- unsubscribe: Unsubscribe from a subscription for all topics in a namespace\n" +
		"- unload: Unload a namespace from the current serving broker\n" +
		"- split_bundle: Split a namespace bundle"
	pulsarAdminNamespaceTenantDesc       = "The tenant name. Required for 'list' operation."
	pulsarAdminNamespaceNamespaceDesc    = "The namespace name in format 'tenant/namespace'. Required for all operations except 'list'."
	pulsarAdminNamespaceBundlesDesc      = "Number of bundles to activate when creating a namespace (default: 0 for default number of bundles). Used with 'create' operation."
	pulsarAdminNamespaceClustersDesc     = "List of clusters to assign when creating a namespace. Used with 'create' operation."
	pulsarAdminNamespaceSubscriptionDesc = "Subscription name. Required for 'unsubscribe' operation, optional for 'clear_backlog'."
	pulsarAdminNamespaceBundleDesc       = "Bundle name or range. Required for 'split_bundle' operation, optional for 'clear_backlog', 'unsubscribe', and 'unload'."
	pulsarAdminNamespaceForceDesc        = "Force clear backlog (true/false). Used with 'clear_backlog' operation."
	pulsarAdminNamespaceUnloadDesc       = "Unload newly split bundles after splitting (true/false). Used with 'split_bundle' operation."
)

// PulsarAdminNamespaceToolBuilder implements the ToolBuilder interface for Pulsar Admin Namespace tools
// It provides functionality to build Pulsar namespace management tools
// /nolint:revive
type PulsarAdminNamespaceToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminNamespaceToolBuilder creates a new Pulsar Admin Namespace tool builder instance
func NewPulsarAdminNamespaceToolBuilder() *PulsarAdminNamespaceToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_namespace",
		Version:     "1.0.0",
		Description: "Pulsar Admin namespace management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "namespace", "admin"},
	}

	features := []string{
		"pulsar-admin-namespaces",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminNamespaceToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar Admin Namespace tool list
// This is the core method implementing the ToolBuilder interface
func (b *PulsarAdminNamespaceToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool, err := b.buildNamespaceTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildNamespaceHandler(config.ReadOnly)

	return []builders.ToolDefinition{
		builders.ServerTool[pulsarAdminNamespaceInput, any]{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildNamespaceTool builds the Pulsar Admin Namespace MCP tool definition
// Migrated from the original tool definition logic
func (b *PulsarAdminNamespaceToolBuilder) buildNamespaceTool() (*sdk.Tool, error) {
	inputSchema, err := buildPulsarAdminNamespaceInputSchema()
	if err != nil {
		return nil, err
	}

	toolDesc := "Manage Pulsar namespaces with various operations. " +
		"This tool provides functionality to work with namespaces in Apache Pulsar, " +
		"including listing, creating, deleting, and performing various operations on namespaces."

	return &sdk.Tool{
		Name:        "pulsar_admin_namespace",
		Description: toolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildNamespaceHandler builds the Pulsar Admin Namespace handler function
// Migrated from the original handler logic
func (b *PulsarAdminNamespaceToolBuilder) buildNamespaceHandler(readOnly bool) builders.ToolHandlerFunc[pulsarAdminNamespaceInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input pulsarAdminNamespaceInput) (*sdk.CallToolResult, any, error) {
		operation := input.Operation
		if operation == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'operation'")
		}

		// Validate write operations in read-only mode
		if readOnly && (operation == "create" || operation == "delete" || operation == "clear_backlog" ||
			operation == "unsubscribe" || operation == "unload" || operation == "split_bundle") {
			return nil, nil, fmt.Errorf("write operations are not allowed in read-only mode")
		}

		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return nil, nil, fmt.Errorf("pulsar session not found in context")
		}

		// Create Pulsar client
		client, err := session.GetAdminClient()
		if err != nil {
			return nil, nil, b.handleError("get admin client", err)
		}

		// Route to appropriate handler based on operation
		switch operation {
		case "list":
			result, err := b.handleNamespaceList(client, input)
			return result, nil, err
		case "get_topics":
			result, err := b.handleNamespaceGetTopics(client, input)
			return result, nil, err
		case "create":
			result, err := b.handleNamespaceCreate(client, input)
			return result, nil, err
		case "delete":
			result, err := b.handleNamespaceDelete(client, input)
			return result, nil, err
		case "clear_backlog":
			result, err := b.handleClearBacklog(client, input)
			return result, nil, err
		case "unsubscribe":
			result, err := b.handleUnsubscribe(client, input)
			return result, nil, err
		case "unload":
			result, err := b.handleUnload(client, input)
			return result, nil, err
		case "split_bundle":
			result, err := b.handleSplitBundle(client, input)
			return result, nil, err
		default:
			return nil, nil, fmt.Errorf("unknown operation: %s. supported operations: list, get_topics, create, delete, clear_backlog, unsubscribe, unload, split_bundle", operation)
		}
	}
}

// Unified error handling and utility functions

// handleError provides unified error handling
func (b *PulsarAdminNamespaceToolBuilder) handleError(operation string, err error) error {
	return fmt.Errorf("failed to %s: %v", operation, err)
}

// marshalResponse provides unified JSON serialization for responses
func (b *PulsarAdminNamespaceToolBuilder) marshalResponse(data interface{}) (*sdk.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, b.handleError("marshal response", err)
	}
	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: string(jsonBytes)}},
	}, nil
}

// Operation handler functions - migrated from the original implementation

// handleNamespaceList handles listing namespaces for a tenant
func (b *PulsarAdminNamespaceToolBuilder) handleNamespaceList(client cmdutils.Client, input pulsarAdminNamespaceInput) (*sdk.CallToolResult, error) {
	tenant, err := requireString(input.Tenant, "tenant")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'tenant' for namespace.list: %v", err)
	}

	// Get namespace list
	namespaces, err := client.Namespaces().GetNamespaces(tenant)
	if err != nil {
		return nil, b.handleError("list namespaces", err)
	}

	return b.marshalResponse(namespaces)
}

// handleNamespaceGetTopics handles getting topics for a namespace
func (b *PulsarAdminNamespaceToolBuilder) handleNamespaceGetTopics(client cmdutils.Client, input pulsarAdminNamespaceInput) (*sdk.CallToolResult, error) {
	namespace, err := requireString(input.Namespace, "namespace")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'namespace' for namespace.get_topics: %v", err)
	}

	// Get topics list
	topics, err := client.Namespaces().GetTopics(namespace)
	if err != nil {
		return nil, b.handleError("get topics", err)
	}

	return b.marshalResponse(topics)
}

// handleNamespaceCreate handles creating a new namespace
func (b *PulsarAdminNamespaceToolBuilder) handleNamespaceCreate(client cmdutils.Client, input pulsarAdminNamespaceInput) (*sdk.CallToolResult, error) {
	namespace, err := requireString(input.Namespace, "namespace")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'namespace' for namespace.create: %v", err)
	}

	// Get optional parameters
	bundlesStr := stringValue(input.Bundles)
	bundles := 0
	if bundlesStr != "" {
		bundlesInt, err := strconv.Atoi(bundlesStr)
		if err != nil {
			return nil, fmt.Errorf("invalid bundles value, must be an integer: %v", err)
		}
		bundles = bundlesInt
	}

	clusters := input.Clusters

	// Prepare policies
	policies := utils.NewDefaultPolicies()

	// Set bundles if provided
	if bundles > 0 {
		if bundles < 0 || bundles > int(^uint32(0)) { // MaxInt32
			return nil, fmt.Errorf("invalid number of bundles, number of bundles has to be in the range of (0, %d]", int(^uint32(0)))
		}
		policies.Bundles = utils.NewBundlesDataWithNumBundles(bundles)
	}

	// Set clusters if provided
	if len(clusters) > 0 {
		policies.ReplicationClusters = clusters
	}

	// Create namespace
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return nil, fmt.Errorf("invalid namespace name: %v", err)
	}

	err = client.Namespaces().CreateNsWithPolices(ns.String(), *policies)
	if err != nil {
		return nil, b.handleError("create namespace", err)
	}

	return textResult(fmt.Sprintf("Created %s successfully", namespace)), nil
}

// handleNamespaceDelete handles deleting a namespace
func (b *PulsarAdminNamespaceToolBuilder) handleNamespaceDelete(client cmdutils.Client, input pulsarAdminNamespaceInput) (*sdk.CallToolResult, error) {
	namespace, err := requireString(input.Namespace, "namespace")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'namespace' for namespace.delete: %v", err)
	}

	// Delete namespace
	err = client.Namespaces().DeleteNamespace(namespace)
	if err != nil {
		return nil, b.handleError("delete namespace", err)
	}

	return textResult(fmt.Sprintf("Deleted %s successfully", namespace)), nil
}

// handleClearBacklog handles clearing the backlog for all topics in a namespace
func (b *PulsarAdminNamespaceToolBuilder) handleClearBacklog(client cmdutils.Client, input pulsarAdminNamespaceInput) (*sdk.CallToolResult, error) {
	namespace, err := requireString(input.Namespace, "namespace")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'namespace' for namespace.clear_backlog: %v", err)
	}

	// Get optional parameters
	subscription := stringValue(input.Subscription)
	bundle := stringValue(input.Bundle)
	forceFlag := stringValue(input.Force) == "true"

	// If not forced, return an error requiring explicit force flag
	if !forceFlag {
		return nil, fmt.Errorf("clear backlog operation requires explicit confirmation, set force=true to proceed")
	}

	// Get namespace name
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return nil, fmt.Errorf("invalid namespace name: %v", err)
	}

	// Handle different backlog clearing scenarios
	var clearErr error
	//nolint:gocritic
	if subscription != "" {
		if bundle != "" {
			clearErr = client.Namespaces().ClearNamespaceBundleBacklogForSubscription(*ns, bundle, subscription)
		} else {
			clearErr = client.Namespaces().ClearNamespaceBacklogForSubscription(*ns, subscription)
		}
	} else if bundle != "" {
		clearErr = client.Namespaces().ClearNamespaceBundleBacklog(*ns, bundle)
	} else {
		clearErr = client.Namespaces().ClearNamespaceBacklog(*ns)
	}

	if clearErr != nil {
		return nil, b.handleError("clear backlog", clearErr)
	}

	return textResult(
		fmt.Sprintf("Successfully cleared backlog for all topics in namespace %s", namespace),
	), nil
}

// handleUnsubscribe handles unsubscribing the specified subscription for all topics of a namespace
func (b *PulsarAdminNamespaceToolBuilder) handleUnsubscribe(client cmdutils.Client, input pulsarAdminNamespaceInput) (*sdk.CallToolResult, error) {
	namespace, err := requireString(input.Namespace, "namespace")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'namespace' for namespace.unsubscribe: %v", err)
	}

	subscription, err := requireString(input.Subscription, "subscription")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'subscription' for namespace.unsubscribe: %v", err)
	}

	// Get optional bundle
	bundle := stringValue(input.Bundle)

	// Get namespace name
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return nil, fmt.Errorf("invalid namespace name: %v", err)
	}

	// Unsubscribe namespace
	var unsubErr error
	if bundle == "" {
		unsubErr = client.Namespaces().UnsubscribeNamespace(*ns, subscription)
	} else {
		unsubErr = client.Namespaces().UnsubscribeNamespaceBundle(*ns, bundle, subscription)
	}

	if unsubErr != nil {
		return nil, b.handleError("unsubscribe", unsubErr)
	}

	if bundle == "" {
		return textResult(
			fmt.Sprintf("Successfully unsubscribed the subscription %s for all topics of the namespace %s",
				subscription, namespace),
		), nil
	}

	return textResult(
		fmt.Sprintf("Successfully unsubscribed the subscription %s for all topics of the namespace %s with bundle range %s",
			subscription, namespace, bundle),
	), nil
}

// handleUnload handles unloading a namespace from the current serving broker
func (b *PulsarAdminNamespaceToolBuilder) handleUnload(client cmdutils.Client, input pulsarAdminNamespaceInput) (*sdk.CallToolResult, error) {
	namespace, err := requireString(input.Namespace, "namespace")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'namespace' for namespace.unload: %v", err)
	}

	// Get optional bundle
	bundle := stringValue(input.Bundle)

	// Unload namespace
	var unloadErr error
	if bundle == "" {
		unloadErr = client.Namespaces().Unload(namespace)
	} else {
		unloadErr = client.Namespaces().UnloadNamespaceBundle(namespace, bundle)
	}

	if unloadErr != nil {
		return nil, b.handleError("unload namespace", unloadErr)
	}

	if bundle == "" {
		return textResult(
			fmt.Sprintf("Unloaded namespace %s successfully", namespace),
		), nil
	}

	return textResult(
		fmt.Sprintf("Unloaded namespace %s with bundle %s successfully", namespace, bundle),
	), nil
}

// handleSplitBundle handles splitting a namespace bundle
func (b *PulsarAdminNamespaceToolBuilder) handleSplitBundle(client cmdutils.Client, input pulsarAdminNamespaceInput) (*sdk.CallToolResult, error) {
	namespace, err := requireString(input.Namespace, "namespace")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'namespace' for namespace.split_bundle: %v", err)
	}

	bundle, err := requireString(input.Bundle, "bundle")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'bundle' for namespace.split_bundle: %v", err)
	}

	// Get optional unload flag
	unload := stringValue(input.Unload) == "true"

	// Split namespace bundle
	err = client.Namespaces().SplitNamespaceBundle(namespace, bundle, unload)
	if err != nil {
		return nil, b.handleError("split namespace bundle", err)
	}

	return textResult(
		fmt.Sprintf("Split namespace bundle %s successfully", bundle),
	), nil
}

func buildPulsarAdminNamespaceInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[pulsarAdminNamespaceInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}

	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}

	setSchemaDescription(schema, "operation", pulsarAdminNamespaceOperationDesc)
	setSchemaDescription(schema, "tenant", pulsarAdminNamespaceTenantDesc)
	setSchemaDescription(schema, "namespace", pulsarAdminNamespaceNamespaceDesc)
	setSchemaDescription(schema, "bundles", pulsarAdminNamespaceBundlesDesc)
	setSchemaDescription(schema, "clusters", pulsarAdminNamespaceClustersDesc)
	setSchemaDescription(schema, "subscription", pulsarAdminNamespaceSubscriptionDesc)
	setSchemaDescription(schema, "bundle", pulsarAdminNamespaceBundleDesc)
	setSchemaDescription(schema, "force", pulsarAdminNamespaceForceDesc)
	setSchemaDescription(schema, "unload", pulsarAdminNamespaceUnloadDesc)

	if clustersSchema := schema.Properties["clusters"]; clustersSchema != nil && clustersSchema.Items != nil {
		clustersSchema.Items.Description = "cluster"
	}

	normalizeAdditionalProperties(schema)
	return schema, nil
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
