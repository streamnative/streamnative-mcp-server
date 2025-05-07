// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
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
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// PulsarAdminAddResourceQuotasTools adds resource quotas-related tools to the MCP server
func PulsarAdminAddResourceQuotasTools(s *server.MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeaturePulsarAdminResourceQuotas)) && !slices.Contains(features, string(FeatureAll)) && !slices.Contains(features, string(FeatureAllPulsar)) {
		return
	}
	// Add unified resource quotas management tool
	toolDesc := "Manage Apache Pulsar resource quotas for brokers, namespaces and bundles. " +
		"Resource quotas define limits for resource usage such as message rates, bandwidth, and memory. " +
		"These quotas help prevent resource abuse and ensure fair resource allocation across the Pulsar cluster. " +
		"Operations include getting, setting, and resetting quotas. " +
		"Requires super-user permissions for all operations."

	resourceDesc := "Resource to operate on. Available resources:\n" +
		"- quota: The resource quota configuration for a specific namespace bundle or the default quota"

	operationDesc := "Operation to perform. Available operations:\n" +
		"- get: Get the resource quota for a specified namespace bundle or default quota\n" +
		"- set: Set the resource quota for a specified namespace bundle or default quota (requires super-user permissions)\n" +
		"- reset: Reset a namespace bundle's resource quota to default value (requires super-user permissions)"

	resourceQuotasTool := mcp.NewTool("pulsar_admin_resourcequota",
		mcp.WithDescription(toolDesc),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description(resourceDesc),
		),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
		),
		mcp.WithString("namespace",
			mcp.Description("The namespace name in the format 'tenant/namespace'. "+
				"Optional for 'get' and 'set' operations (to get/set default quota if omitted). "+
				"Required for 'reset' operation."),
		),
		mcp.WithString("bundle",
			mcp.Description("The bundle range in the format '{start-boundary}_{end-boundary}'. "+
				"Must be specified together with namespace. Bundle is a hash range of the topic names belonging to a namespace."),
		),
		mcp.WithNumber("msgRateIn",
			mcp.Description("Expected incoming messages per second. Required for 'set' operation. "+
				"This defines the maximum rate of incoming messages allowed for the namespace or bundle."),
		),
		mcp.WithNumber("msgRateOut",
			mcp.Description("Expected outgoing messages per second. Required for 'set' operation. "+
				"This defines the maximum rate of outgoing messages allowed for the namespace or bundle."),
		),
		mcp.WithNumber("bandwidthIn",
			mcp.Description("Expected inbound bandwidth in bytes per second. Required for 'set' operation. "+
				"This defines the maximum rate of incoming bytes allowed for the namespace or bundle."),
		),
		mcp.WithNumber("bandwidthOut",
			mcp.Description("Expected outbound bandwidth in bytes per second. Required for 'set' operation. "+
				"This defines the maximum rate of outgoing bytes allowed for the namespace or bundle."),
		),
		mcp.WithNumber("memory",
			mcp.Description("Expected memory usage in Mbytes. Required for 'set' operation. "+
				"This defines the maximum memory allowed for storing messages for the namespace or bundle."),
		),
		mcp.WithBoolean("dynamic",
			mcp.Description("Whether to allow quota to be dynamically re-calculated. Optional for 'set' operation. "+
				"If true, the broker can dynamically adjust the quota based on the current usage patterns."),
		),
	)
	s.AddTool(resourceQuotasTool, handleResourceQuotaTool(readOnly))
}

// handleResourceQuotaTool returns a function that handles resource quota operations
func handleResourceQuotaTool(readOnly bool) func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get required parameters
		resource, err := requiredParam[string](request.Params.Arguments, "resource")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get resource: %v", err)), nil
		}

		operation, err := requiredParam[string](request.Params.Arguments, "operation")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get operation: %v", err)), nil
		}

		// Normalize parameters
		resource = strings.ToLower(resource)
		operation = strings.ToLower(operation)

		// Validate write operations in read-only mode
		if readOnly && (operation == "set" || operation == "reset") {
			return mcp.NewToolResultError("Write operations are not allowed in read-only mode"), nil
		}

		// Verify resource type
		if resource != "quota" {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resource: %s. Only 'quota' is supported", resource)), nil
		}

		// Create the admin client
		admin, err := pulsar.GetAdminClient()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
		}

		// Dispatch based on operation
		switch operation {
		case "get":
			return handleQuotaGet(admin, request)
		case "set":
			return handleQuotaSet(admin, request)
		case "reset":
			return handleQuotaReset(admin, request)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid operation: %s. Available operations: get, set, reset", operation)), nil
		}
	}
}

// handleQuotaGet handles getting a resource quota
func handleQuotaGet(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get optional parameters
	namespace, hasNamespace := optionalParam[string](request.Params.Arguments, "namespace")
	bundle, hasBundle := optionalParam[string](request.Params.Arguments, "bundle")

	// Check if both namespace and bundle are provided or neither is provided
	if (hasNamespace && !hasBundle) || (!hasNamespace && hasBundle) {
		return mcp.NewToolResultError("When specifying a namespace, you must also specify a bundle and vice versa."), nil
	}

	var (
		resourceQuotaData *utils.ResourceQuota
		getErr            error
	)

	if !hasNamespace && !hasBundle {
		// Get default resource quota
		resourceQuotaData, getErr = admin.ResourceQuotas().GetDefaultResourceQuota()
		if getErr != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get default resource quota: %v", getErr)), nil
		}
	} else {
		// Get namespace bundle resource quota
		nsName, getErr := utils.GetNamespaceName(namespace)
		if getErr != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name '%s': %v", namespace, getErr)), nil
		}
		resourceQuotaData, getErr = admin.ResourceQuotas().GetNamespaceBundleResourceQuota(nsName.String(), bundle)
		if getErr != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get resource quota for namespace '%s' bundle '%s': %v",
				namespace, bundle, getErr)), nil
		}
	}

	// Format the output
	jsonBytes, err := json.Marshal(resourceQuotaData)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal resource quota data: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleQuotaSet handles setting a resource quota
func handleQuotaSet(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters for set operation
	msgRateIn, err := requiredParam[float64](request.Params.Arguments, "msgRateIn")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'msgRateIn' for quota.set: %v", err)), nil
	}

	msgRateOut, err := requiredParam[float64](request.Params.Arguments, "msgRateOut")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'msgRateOut' for quota.set: %v", err)), nil
	}

	bandwidthIn, err := requiredParam[float64](request.Params.Arguments, "bandwidthIn")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'bandwidthIn' for quota.set: %v", err)), nil
	}

	bandwidthOut, err := requiredParam[float64](request.Params.Arguments, "bandwidthOut")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'bandwidthOut' for quota.set: %v", err)), nil
	}

	memory, err := requiredParam[float64](request.Params.Arguments, "memory")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'memory' for quota.set: %v", err)), nil
	}

	// Get optional parameters
	namespace, hasNamespace := optionalParam[string](request.Params.Arguments, "namespace")
	bundle, hasBundle := optionalParam[string](request.Params.Arguments, "bundle")
	dynamic, hasDynamic := optionalParam[bool](request.Params.Arguments, "dynamic")

	if !hasDynamic {
		dynamic = false
	}

	// Check if both namespace and bundle are provided or neither is provided
	if (hasNamespace && !hasBundle) || (!hasNamespace && hasBundle) {
		return mcp.NewToolResultError("When specifying a namespace, you must also specify a bundle and vice versa."), nil
	}

	// Create resource quota object
	quota := utils.NewResourceQuota()
	quota.MsgRateIn = msgRateIn
	quota.MsgRateOut = msgRateOut
	quota.BandwidthIn = bandwidthIn
	quota.BandwidthOut = bandwidthOut
	quota.Memory = memory
	quota.Dynamic = dynamic

	var resultMsg string
	if !hasNamespace && !hasBundle {
		// Set default resource quota
		err = admin.ResourceQuotas().SetDefaultResourceQuota(*quota)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to set default resource quota: %v", err)), nil
		}
		resultMsg = "Default resource quota set successfully"
	} else {
		// Set namespace bundle resource quota
		err = admin.ResourceQuotas().SetNamespaceBundleResourceQuota(namespace, bundle, *quota)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to set resource quota for namespace '%s' bundle '%s': %v",
				namespace, bundle, err)), nil
		}
		resultMsg = fmt.Sprintf("Resource quota for namespace '%s' bundle '%s' set successfully", namespace, bundle)
	}

	return mcp.NewToolResultText(resultMsg), nil
}

// handleQuotaReset handles resetting a resource quota
func handleQuotaReset(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters for reset operation
	namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'namespace' for quota.reset: %v", err)), nil
	}

	bundle, err := requiredParam[string](request.Params.Arguments, "bundle")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'bundle' for quota.reset: %v", err)), nil
	}

	// Parse namespace name
	nsName, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid namespace name '%s': %v", namespace, err)), nil
	}

	// Reset namespace bundle resource quota
	err = admin.ResourceQuotas().ResetNamespaceBundleResourceQuota(nsName.String(), bundle)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to reset resource quota for namespace '%s' bundle '%s': %v",
			namespace, bundle, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Resource quota for namespace '%s' bundle '%s' reset to default successfully",
		namespace, bundle)), nil
}
