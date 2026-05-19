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
	"strconv"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/toolannotations"
)

var pulsarNamespaceWriteOperations = map[string]struct{}{
	"create":        {},
	"delete":        {},
	"clear_backlog": {},
	"unsubscribe":   {},
	"unload":        {},
	"split_bundle":  {},
}

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
func (b *PulsarAdminNamespaceToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
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
			Tool:    b.buildNamespaceTool(toolModeRead),
			Handler: b.buildNamespaceHandler(toolModeRead),
		},
	}
	if !config.ReadOnly {
		tools = append(tools, server.ServerTool{
			Tool:    b.buildNamespaceTool(toolModeWrite),
			Handler: b.buildNamespaceHandler(toolModeWrite),
		})
	}

	return tools, nil
}

// buildNamespaceTool builds the Pulsar Admin Namespace MCP tool definition
// Migrated from the original tool definition logic
func (b *PulsarAdminNamespaceToolBuilder) buildNamespaceTool(mode toolMode) mcp.Tool {
	toolDesc := "Read Apache Pulsar namespaces. " +
		"This read-only tool lists namespaces for a tenant and lists topics in a namespace without changing namespace state."

	operationDesc := "Operation to perform on namespaces. Available operations:\n" +
		"- list: List all namespaces for a tenant\n" +
		"- get_topics: Get all topics within a namespace"

	operationEnum := []string{"list", "get_topics"}
	toolName := "pulsar_admin_namespace_read"
	annotation := toolannotations.ReadOnly("Read Pulsar Namespaces")
	if isToolModeWrite(mode) {
		toolDesc = "Manage Apache Pulsar namespaces. " +
			"This write tool creates, deletes, unloads, splits bundles, clears backlog, or unsubscribes namespace subscriptions."
		operationDesc = "Operation to perform on namespaces. Available operations:\n" +
			"- create: Create a new namespace\n" +
			"- delete: Delete an existing namespace\n" +
			"- clear_backlog: Clear backlog for all topics in a namespace\n" +
			"- unsubscribe: Unsubscribe from a subscription for all topics in a namespace\n" +
			"- unload: Unload a namespace from the current serving broker\n" +
			"- split_bundle: Split a namespace bundle"
		operationEnum = []string{"create", "delete", "clear_backlog", "unsubscribe", "unload", "split_bundle"}
		toolName = "pulsar_admin_namespace_write"
		annotation = toolannotations.Destructive("Manage Pulsar Namespaces")
	}

	tool := mcp.NewTool(toolName,
		mcp.WithDescription(toolDesc),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
			mcp.Enum(operationEnum...),
		),
		mcp.WithString("tenant",
			mcp.Description("The tenant name. Required for 'list' operation."),
		),
		mcp.WithString("namespace",
			mcp.Description("The namespace name in format 'tenant/namespace'. Required for all operations except 'list'."),
		),
		mcp.WithString("bundles",
			mcp.Description("Number of bundles to activate when creating a namespace (default: 0 for default number of bundles). Used with 'create' operation."),
		),
		mcp.WithArray("clusters",
			mcp.Description("List of clusters to assign when creating a namespace. Used with 'create' operation."),
			mcp.Items(
				map[string]interface{}{
					"type":        "string",
					"description": "Cluster name",
				},
			),
		),
		mcp.WithString("subscription",
			mcp.Description("Subscription name. Required for 'unsubscribe' operation, optional for 'clear_backlog'."),
		),
		mcp.WithString("bundle",
			mcp.Description("Bundle name or range. Required for 'split_bundle' operation, optional for 'clear_backlog', 'unsubscribe', and 'unload'."),
		),
		mcp.WithString("force",
			mcp.Description("Force clear backlog (true/false). Used with 'clear_backlog' operation."),
		),
		mcp.WithString("unload",
			mcp.Description("Unload newly split bundles after splitting (true/false). Used with 'split_bundle' operation."),
		),
		annotation,
	)
	if isToolModeWrite(mode) {
		pruneToolInputSchema(&tool, []string{"operation", "namespace", "bundles", "clusters", "subscription", "bundle", "force", "unload"})
	} else {
		pruneToolInputSchema(&tool, []string{"operation", "tenant", "namespace"})
	}
	return tool
}

// buildNamespaceHandler builds the Pulsar Admin Namespace handler function
// Migrated from the original handler logic
func (b *PulsarAdminNamespaceToolBuilder) buildNamespaceHandler(mode toolMode) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get operation parameter
		operation, err := request.RequireString("operation")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get operation: %v", err)), nil
		}

		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return mcp.NewToolResultError("Pulsar session not found in context"), nil
		}

		// Create Pulsar client
		client, err := session.GetAdminClient()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
		}

		if !validateModeOperation(mode, operation, pulsarNamespaceWriteOperations) {
			return mcp.NewToolResultError(fmt.Sprintf("Operation %q is not available in %s mode", operation, mode)), nil
		}

		// Route to appropriate handler based on operation
		switch operation {
		case "list":
			return b.handleNamespaceList(ctx, client, request)
		case "get_topics":
			return b.handleNamespaceGetTopics(ctx, client, request)
		case "create", "delete", "clear_backlog", "unsubscribe", "unload", "split_bundle":
			// Route to appropriate write operation handler
			switch operation {
			case "create":
				return b.handleNamespaceCreate(ctx, client, request)
			case "delete":
				return b.handleNamespaceDelete(ctx, client, request)
			case "clear_backlog":
				return b.handleClearBacklog(ctx, client, request)
			case "unsubscribe":
				return b.handleUnsubscribe(ctx, client, request)
			case "unload":
				return b.handleUnload(ctx, client, request)
			case "split_bundle":
				return b.handleSplitBundle(ctx, client, request)
			}
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Unknown operation: %s. Supported operations: %s", operation,
				modeSupportedOperations(mode,
					[]string{"list", "get_topics"},
					[]string{"create", "delete", "clear_backlog", "unsubscribe", "unload", "split_bundle"}))), nil
		}

		// Should not reach here
		return mcp.NewToolResultError("Unexpected error: operation not handled"), nil
	}
}

// Unified error handling and utility functions

// handleError provides unified error handling
func (b *PulsarAdminNamespaceToolBuilder) handleError(operation string, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("Failed to %s: %v", operation, err))
}

// marshalResponse provides unified JSON serialization for responses
func (b *PulsarAdminNamespaceToolBuilder) marshalResponse(data interface{}) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return b.handleError("marshal response", err), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// Operation handler functions - migrated from the original implementation

// handleNamespaceList handles listing namespaces for a tenant
func (b *PulsarAdminNamespaceToolBuilder) handleNamespaceList(_ context.Context, client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tenant, err := request.RequireString("tenant")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get tenant name: %v", err)), nil
	}

	// Get namespace list
	namespaces, err := client.Namespaces().GetNamespaces(tenant)
	if err != nil {
		return b.handleError("list namespaces", err), nil
	}

	return b.marshalResponse(namespaces)
}

// handleNamespaceGetTopics handles getting topics for a namespace
func (b *PulsarAdminNamespaceToolBuilder) handleNamespaceGetTopics(_ context.Context, client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	namespace, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Get topics list
	topics, err := client.Namespaces().GetTopics(namespace)
	if err != nil {
		return b.handleError("get topics", err), nil
	}

	return b.marshalResponse(topics)
}

// handleNamespaceCreate handles creating a new namespace
func (b *PulsarAdminNamespaceToolBuilder) handleNamespaceCreate(_ context.Context, client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	namespace, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Get optional parameters
	bundlesStr := request.GetString("bundles", "")
	bundles := 0
	if bundlesStr != "" {
		bundlesInt, err := strconv.Atoi(bundlesStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid bundles value, must be an integer: %v", err)), nil
		}
		bundles = bundlesInt
	}

	clusters := request.GetStringSlice("clusters", []string{})

	// Prepare policies
	policies := utils.NewDefaultPolicies()

	// Set bundles if provided
	if bundles > 0 {
		if bundles < 0 || bundles > int(^uint32(0)) { // MaxInt32
			return mcp.NewToolResultError(
				fmt.Sprintf("Invalid number of bundles. Number of bundles has to be in the range of (0, %d].", int(^uint32(0))),
			), nil
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
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	err = client.Namespaces().CreateNsWithPolices(ns.String(), *policies)
	if err != nil {
		return b.handleError("create namespace", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Created %s successfully", namespace)), nil
}

// handleNamespaceDelete handles deleting a namespace
func (b *PulsarAdminNamespaceToolBuilder) handleNamespaceDelete(_ context.Context, client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	namespace, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Delete namespace
	err = client.Namespaces().DeleteNamespace(namespace)
	if err != nil {
		return b.handleError("delete namespace", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Deleted %s successfully", namespace)), nil
}

// handleClearBacklog handles clearing the backlog for all topics in a namespace
func (b *PulsarAdminNamespaceToolBuilder) handleClearBacklog(_ context.Context, client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	namespace, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Get optional parameters
	subscription := request.GetString("subscription", "")
	bundle := request.GetString("bundle", "")
	force := request.GetString("force", "")
	forceFlag := force == "true"

	// If not forced, return an error requiring explicit force flag
	if !forceFlag {
		return mcp.NewToolResultError(
			"Clear backlog operation requires explicit confirmation. Please set force=true to proceed.",
		), nil
	}

	// Get namespace name
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
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
		return b.handleError("clear backlog", clearErr), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Successfully cleared backlog for all topics in namespace %s", namespace),
	), nil
}

// handleUnsubscribe handles unsubscribing the specified subscription for all topics of a namespace
func (b *PulsarAdminNamespaceToolBuilder) handleUnsubscribe(_ context.Context, client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	namespace, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	subscription, err := request.RequireString("subscription")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get subscription name: %v", err)), nil
	}

	// Get optional bundle
	bundle := request.GetString("bundle", "")

	// Get namespace name
	ns, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name: %v", err)), nil
	}

	// Unsubscribe namespace
	var unsubErr error
	if bundle == "" {
		unsubErr = client.Namespaces().UnsubscribeNamespace(*ns, subscription)
	} else {
		unsubErr = client.Namespaces().UnsubscribeNamespaceBundle(*ns, bundle, subscription)
	}

	if unsubErr != nil {
		return b.handleError("unsubscribe", unsubErr), nil
	}

	if bundle == "" {
		return mcp.NewToolResultText(
			fmt.Sprintf("Successfully unsubscribed the subscription %s for all topics of the namespace %s",
				subscription, namespace),
		), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Successfully unsubscribed the subscription %s for all topics of the namespace %s with bundle range %s",
			subscription, namespace, bundle),
	), nil
}

// handleUnload handles unloading a namespace from the current serving broker
func (b *PulsarAdminNamespaceToolBuilder) handleUnload(_ context.Context, client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	namespace, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Get optional bundle
	bundle := request.GetString("bundle", "")

	// Unload namespace
	var unloadErr error
	if bundle == "" {
		unloadErr = client.Namespaces().Unload(namespace)
	} else {
		unloadErr = client.Namespaces().UnloadNamespaceBundle(namespace, bundle)
	}

	if unloadErr != nil {
		return b.handleError("unload namespace", unloadErr), nil
	}

	if bundle == "" {
		return mcp.NewToolResultText(
			fmt.Sprintf("Unloaded namespace %s successfully", namespace),
		), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Unloaded namespace %s with bundle %s successfully", namespace, bundle),
	), nil
}

// handleSplitBundle handles splitting a namespace bundle
func (b *PulsarAdminNamespaceToolBuilder) handleSplitBundle(_ context.Context, client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	namespace, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	bundle, err := request.RequireString("bundle")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get bundle: %v", err)), nil
	}

	// Get optional unload flag
	unload := request.GetString("unload", "") == "true"

	// Split namespace bundle
	err = client.Namespaces().SplitNamespaceBundle(namespace, bundle, unload)
	if err != nil {
		return b.handleError("split namespace bundle", err), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Split namespace bundle %s successfully", bundle),
	), nil
}
