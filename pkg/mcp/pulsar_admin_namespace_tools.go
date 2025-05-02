// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// PulsarAdminAddNamespaceTools adds namespace-related tools to the MCP server
func PulsarAdminAddNamespaceTools(s *server.MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeaturePulsarAdminNamespaces)) && !slices.Contains(features, string(FeatureAll)) && !slices.Contains(features, string(FeatureAllPulsar)) {
		return
	}

	// Create a unified namespace tool with resource+operation parameters
	namespaceTool := mcp.NewTool("pulsar_admin_namespace",
		mcp.WithDescription("Manage Pulsar namespaces with various operations. "+
			"This tool provides functionality to work with namespaces in Apache Pulsar, "+
			"including listing, creating, deleting, and performing various operations on namespaces."),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description("Operation to perform on namespaces. Available operations:\n"+
				"- list: List all namespaces for a tenant\n"+
				"- get_topics: Get all topics within a namespace\n"+
				"- create: Create a new namespace\n"+
				"- delete: Delete an existing namespace\n"+
				"- clear_backlog: Clear backlog for all topics in a namespace\n"+
				"- unsubscribe: Unsubscribe from a subscription for all topics in a namespace\n"+
				"- unload: Unload a namespace from the current serving broker\n"+
				"- split_bundle: Split a namespace bundle"),
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
	)

	// Add the unified tool with the handler closure that captures readOnly status
	s.AddTool(namespaceTool, handleNamespace(readOnly))
}

// handleNamespace returns a handler function for the namespace tool that captures the readOnly status
func handleNamespace(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get operation parameter
		operation, err := requiredParam[string](request.Params.Arguments, "operation")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get operation: %v", err)), nil
		}

		// Create Pulsar client
		client, err := pulsar.GetAdminClient()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
		}

		// Route to appropriate handler based on operation
		switch operation {
		case "list":
			return handleNamespaceList(ctx, client, request)
		case "get_topics":
			return handleNamespaceGetTopics(ctx, client, request)
		case "create", "delete", "clear_backlog", "unsubscribe", "unload", "split_bundle":
			// Check if write operations are allowed
			if readOnly {
				return mcp.NewToolResultError(fmt.Sprintf("Operation '%s' not allowed in read-only mode", operation)), nil
			}

			// Route to appropriate write operation handler
			switch operation {
			case "create":
				return handleNamespaceCreate(ctx, client, request)
			case "delete":
				return handleNamespaceDelete(ctx, client, request)
			case "clear_backlog":
				return handleClearBacklog(ctx, client, request)
			case "unsubscribe":
				return handleUnsubscribe(ctx, client, request)
			case "unload":
				return handleUnload(ctx, client, request)
			case "split_bundle":
				return handleSplitBundle(ctx, client, request)
			}
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Unknown operation: %s. Supported operations: list, get_topics, create, delete, clear_backlog, unsubscribe, unload, split_bundle", operation)), nil
		}

		// Should not reach here
		return mcp.NewToolResultError("Unexpected error: operation not handled"), nil
	}
}

// handleNamespaceList handles listing namespaces for a tenant
func handleNamespaceList(_ context.Context, client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tenant, err := requiredParam[string](request.Params.Arguments, "tenant")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get tenant name: %v", err)), nil
	}

	// Get namespace list
	namespaces, err := client.Namespaces().GetNamespaces(tenant)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list namespaces: %v", err)), nil
	}

	// Convert result to JSON string
	namespacesJSON, err := json.Marshal(namespaces)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize namespace list: %v", err)), nil
	}

	return mcp.NewToolResultText(string(namespacesJSON)), nil
}

// handleNamespaceGetTopics handles getting topics for a namespace
func handleNamespaceGetTopics(_ context.Context, client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Get topics list
	topics, err := client.Namespaces().GetTopics(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get topics: %v", err)), nil
	}

	// Convert result to JSON string
	topicsJSON, err := json.Marshal(topics)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize topics list: %v", err)), nil
	}

	return mcp.NewToolResultText(string(topicsJSON)), nil
}

// handleNamespaceCreate handles creating a new namespace
func handleNamespaceCreate(_ context.Context, client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Get optional parameters
	bundlesStr, hasBundles := optionalParam[string](request.Params.Arguments, "bundles")
	bundles := 0
	if hasBundles && bundlesStr != "" {
		bundlesInt, err := strconv.Atoi(bundlesStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid bundles value, must be an integer: %v", err)), nil
		}
		bundles = bundlesInt
	}

	clusters, _ := optionalParamArray[string](request.Params.Arguments, "clusters")

	// Prepare policies
	policies := utils.NewDefaultPolicies()

	// Set bundles if provided
	if hasBundles && bundles > 0 {
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
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create namespace: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Created %s successfully", namespace)), nil
}

// handleNamespaceDelete handles deleting a namespace
func handleNamespaceDelete(_ context.Context, client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Delete namespace
	err = client.Namespaces().DeleteNamespace(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete namespace: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Deleted %s successfully", namespace)), nil
}

// handleClearBacklog handles clearing the backlog for all topics in a namespace
func handleClearBacklog(_ context.Context, client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Get optional parameters
	subscription, _ := optionalParam[string](request.Params.Arguments, "subscription")
	bundle, _ := optionalParam[string](request.Params.Arguments, "bundle")
	force, _ := optionalParam[string](request.Params.Arguments, "force")
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
		return mcp.NewToolResultError(fmt.Sprintf("Failed to clear backlog: %v", clearErr)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Successfully cleared backlog for all topics in namespace %s", namespace),
	), nil
}

// handleUnsubscribe handles unsubscribing the specified subscription for all topics of a namespace
func handleUnsubscribe(_ context.Context, client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	subscription, err := requiredParam[string](request.Params.Arguments, "subscription")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get subscription name: %v", err)), nil
	}

	// Get optional bundle
	bundle, _ := optionalParam[string](request.Params.Arguments, "bundle")

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
		return mcp.NewToolResultError(fmt.Sprintf("Failed to unsubscribe: %v", unsubErr)), nil
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
func handleUnload(_ context.Context, client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	// Get optional bundle
	bundle, _ := optionalParam[string](request.Params.Arguments, "bundle")

	// Unload namespace
	var unloadErr error
	if bundle == "" {
		unloadErr = client.Namespaces().Unload(namespace)
	} else {
		unloadErr = client.Namespaces().UnloadNamespaceBundle(namespace, bundle)
	}

	if unloadErr != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to unload namespace: %v", unloadErr)), nil
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
func handleSplitBundle(_ context.Context, client cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get namespace name: %v", err)), nil
	}

	bundle, err := requiredParam[string](request.Params.Arguments, "bundle")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get bundle: %v", err)), nil
	}

	// Get optional unload flag
	unloadStr, _ := optionalParam[string](request.Params.Arguments, "unload")
	unload := unloadStr == "true"

	// Split namespace bundle
	err = client.Namespaces().SplitNamespaceBundle(namespace, bundle, unload)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to split namespace bundle: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Split namespace bundle %s successfully", bundle),
	), nil
}
