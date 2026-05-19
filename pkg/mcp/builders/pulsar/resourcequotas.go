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
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

var pulsarResourceQuotaOperationSpecs = builders.OperationRegistry{
	{Name: "get", Mode: builders.OperationModeRead},
	{Name: "set", Mode: builders.OperationModeWrite, Destructive: true},
	{Name: "reset", Mode: builders.OperationModeWrite, Destructive: true},
}

// PulsarAdminResourceQuotasToolBuilder implements the ToolBuilder interface for Pulsar admin resource quotas
// /nolint:revive
type PulsarAdminResourceQuotasToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminResourceQuotasToolBuilder creates a new Pulsar admin resource quotas tool builder instance
func NewPulsarAdminResourceQuotasToolBuilder() *PulsarAdminResourceQuotasToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_resourcequotas",
		Version:     "1.0.0",
		Description: "Pulsar admin resource quotas management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "admin", "resourcequotas"},
	}

	features := []string{
		"pulsar-admin-resource-quotas",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminResourceQuotasToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin resource quotas tool list
func (b *PulsarAdminResourceQuotasToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
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
			Tool:    b.buildResourceQuotasTool(toolModeRead),
			Handler: b.buildResourceQuotasHandler(toolModeRead),
		},
	}
	if !config.ReadOnly {
		tools = append(tools, server.ServerTool{
			Tool:    b.buildResourceQuotasTool(toolModeWrite),
			Handler: b.buildResourceQuotasHandler(toolModeWrite),
		})
	}

	return tools, nil
}

// buildResourceQuotasTool builds the Pulsar admin resource quotas MCP tool definition
func (b *PulsarAdminResourceQuotasToolBuilder) buildResourceQuotasTool(mode toolMode) mcp.Tool {
	toolDesc := "Read Apache Pulsar resource quotas for brokers, namespaces and bundles. " +
		"Resource quotas define limits for resource usage such as message rates, bandwidth, and memory. " +
		"This read-only tool retrieves quota configuration without changing it."

	resourceDesc := "Resource to operate on. Available resources:\n" +
		"- quota: The resource quota configuration for a specific namespace bundle or the default quota"

	operationDesc := "Operation to perform. Available operations:\n" +
		"- get: Get the resource quota for a specified namespace bundle or default quota"

	operationEnum := pulsarResourceQuotaOperationSpecs.NamesForMode(mode)
	toolName := "pulsar_admin_resourcequota_read"
	annotation := builders.ToolAnnotationForMode(mode, "Read Pulsar Resource Quotas", "Manage Pulsar Resource Quotas", pulsarResourceQuotaOperationSpecs)
	if isToolModeWrite(mode) {
		toolDesc = "Manage Apache Pulsar resource quotas for brokers, namespaces and bundles. " +
			"This write tool sets or resets quota configuration."
		operationDesc = "Operation to perform. Available operations:\n" +
			"- set: Set the resource quota for a specified namespace bundle or default quota (requires super-user permissions)\n" +
			"- reset: Reset a namespace bundle's resource quota to default value (requires super-user permissions)"
		operationEnum = pulsarResourceQuotaOperationSpecs.NamesForMode(mode)
		toolName = "pulsar_admin_resourcequota_write"
	}

	tool := mcp.NewTool(toolName,
		mcp.WithDescription(toolDesc),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description(resourceDesc),
		),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
			mcp.Enum(operationEnum...),
		),
		mcp.WithString("namespace",
			mcp.Description("The namespace name in the format 'tenant/namespace'. Optional when targeting the default quota."),
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
		annotation,
	)
	if isToolModeWrite(mode) {
		pruneToolInputSchema(&tool, []string{"resource", "operation", "namespace", "bundle", "msgRateIn", "msgRateOut", "bandwidthIn", "bandwidthOut", "memory", "dynamic"})
	} else {
		pruneToolInputSchema(&tool, []string{"resource", "operation", "namespace", "bundle"})
	}
	return tool
}

// buildResourceQuotasHandler builds the Pulsar admin resource quotas handler function
func (b *PulsarAdminResourceQuotasToolBuilder) buildResourceQuotasHandler(mode toolMode) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get required parameters
		resource, err := request.RequireString("resource")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get resource: %v", err)), nil
		}

		operation, err := request.RequireString("operation")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get operation: %v", err)), nil
		}

		// Normalize parameters
		resource = strings.ToLower(resource)
		operation = strings.ToLower(operation)

		if err := validateModeOperation(mode, operation, pulsarResourceQuotaOperationSpecs); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Verify resource type
		if resource != "quota" {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resource: %s. Only 'quota' is supported", resource)), nil
		}

		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return mcp.NewToolResultError("Pulsar session not found in context"), nil
		}

		admin, err := session.GetAdminClient()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
		}

		// Dispatch based on operation
		switch operation {
		case "get":
			return b.handleQuotaGet(admin, request)
		case "set":
			return b.handleQuotaSet(admin, request)
		case "reset":
			return b.handleQuotaReset(admin, request)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid operation: %s. Available operations: %s", operation,
				modeSupportedOperations(mode, pulsarResourceQuotaOperationSpecs))), nil
		}
	}
}

// Helper functions

// handleQuotaGet handles getting a resource quota
func (b *PulsarAdminResourceQuotasToolBuilder) handleQuotaGet(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get optional parameters
	namespace := request.GetString("namespace", "")
	bundle := request.GetString("bundle", "")

	// Check if both namespace and bundle are provided or neither is provided
	if (namespace != "" && bundle == "") || (namespace == "" && bundle != "") {
		return mcp.NewToolResultError("When specifying a namespace, you must also specify a bundle and vice versa."), nil
	}

	var (
		resourceQuotaData *utils.ResourceQuota
		getErr            error
	)

	if namespace == "" && bundle == "" {
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
func (b *PulsarAdminResourceQuotasToolBuilder) handleQuotaSet(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters for set operation
	msgRateIn, err := request.RequireFloat("msgRateIn")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'msgRateIn' for quota.set: %v", err)), nil
	}

	msgRateOut, err := request.RequireFloat("msgRateOut")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'msgRateOut' for quota.set: %v", err)), nil
	}

	bandwidthIn, err := request.RequireFloat("bandwidthIn")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'bandwidthIn' for quota.set: %v", err)), nil
	}

	bandwidthOut, err := request.RequireFloat("bandwidthOut")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'bandwidthOut' for quota.set: %v", err)), nil
	}

	memory, err := request.RequireFloat("memory")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'memory' for quota.set: %v", err)), nil
	}

	// Get optional parameters
	namespace := request.GetString("namespace", "")
	bundle := request.GetString("bundle", "")
	dynamic := request.GetBool("dynamic", false)

	// Check if both namespace and bundle are provided or neither is provided
	if (namespace != "" && bundle == "") || (namespace == "" && bundle != "") {
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
	if namespace == "" && bundle == "" {
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
func (b *PulsarAdminResourceQuotasToolBuilder) handleQuotaReset(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get required parameters for reset operation
	namespace, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'namespace' for quota.reset: %v", err)), nil
	}

	bundle, err := request.RequireString("bundle")
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
