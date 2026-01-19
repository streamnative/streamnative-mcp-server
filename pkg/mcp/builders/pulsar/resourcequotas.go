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

type pulsarAdminResourceQuotasInput struct {
	Resource     string   `json:"resource"`
	Operation    string   `json:"operation"`
	Namespace    *string  `json:"namespace,omitempty"`
	Bundle       *string  `json:"bundle,omitempty"`
	MsgRateIn    *float64 `json:"msgRateIn,omitempty"`
	MsgRateOut   *float64 `json:"msgRateOut,omitempty"`
	BandwidthIn  *float64 `json:"bandwidthIn,omitempty"`
	BandwidthOut *float64 `json:"bandwidthOut,omitempty"`
	Memory       *float64 `json:"memory,omitempty"`
	Dynamic      *bool    `json:"dynamic,omitempty"`
}

const (
	pulsarAdminResourceQuotasResourceDesc = "Resource to operate on. Available resources:\n" +
		"- quota: The resource quota configuration for a specific namespace bundle or the default quota"
	pulsarAdminResourceQuotasOperationDesc = "Operation to perform. Available operations:\n" +
		"- get: Get the resource quota for a specified namespace bundle or default quota\n" +
		"- set: Set the resource quota for a specified namespace bundle or default quota (requires super-user permissions)\n" +
		"- reset: Reset a namespace bundle's resource quota to default value (requires super-user permissions)"
	pulsarAdminResourceQuotasNamespaceDesc = "The namespace name in the format 'tenant/namespace'. " +
		"Optional for 'get' and 'set' operations (to get/set default quota if omitted). " +
		"Required for 'reset' operation."
	pulsarAdminResourceQuotasBundleDesc = "The bundle range in the format '{start-boundary}_{end-boundary}'. " +
		"Must be specified together with namespace. Bundle is a hash range of the topic names belonging to a namespace."
	pulsarAdminResourceQuotasMsgRateInDesc = "Expected incoming messages per second. Required for 'set' operation. " +
		"This defines the maximum rate of incoming messages allowed for the namespace or bundle."
	pulsarAdminResourceQuotasMsgRateOutDesc = "Expected outgoing messages per second. Required for 'set' operation. " +
		"This defines the maximum rate of outgoing messages allowed for the namespace or bundle."
	pulsarAdminResourceQuotasBandwidthInDesc = "Expected inbound bandwidth in bytes per second. Required for 'set' operation. " +
		"This defines the maximum rate of incoming bytes allowed for the namespace or bundle."
	pulsarAdminResourceQuotasBandwidthOutDesc = "Expected outbound bandwidth in bytes per second. Required for 'set' operation. " +
		"This defines the maximum rate of outgoing bytes allowed for the namespace or bundle."
	pulsarAdminResourceQuotasMemoryDesc = "Expected memory usage in Mbytes. Required for 'set' operation. " +
		"This defines the maximum memory allowed for storing messages for the namespace or bundle."
	pulsarAdminResourceQuotasDynamicDesc = "Whether to allow quota to be dynamically re-calculated. Optional for 'set' operation. " +
		"If true, the broker can dynamically adjust the quota based on the current usage patterns."
)

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
		"pulsar-admin-resourcequotas",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminResourceQuotasToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin resource quotas tool list
func (b *PulsarAdminResourceQuotasToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool, err := b.buildResourceQuotasTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildResourceQuotasHandler(config.ReadOnly)

	return []builders.ToolDefinition{
		builders.ServerTool[pulsarAdminResourceQuotasInput, any]{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildResourceQuotasTool builds the Pulsar admin resource quotas MCP tool definition
func (b *PulsarAdminResourceQuotasToolBuilder) buildResourceQuotasTool() (*sdk.Tool, error) {
	inputSchema, err := buildPulsarAdminResourceQuotasInputSchema()
	if err != nil {
		return nil, err
	}

	toolDesc := "Manage Apache Pulsar resource quotas for brokers, namespaces and bundles. " +
		"Resource quotas define limits for resource usage such as message rates, bandwidth, and memory. " +
		"These quotas help prevent resource abuse and ensure fair resource allocation across the Pulsar cluster. " +
		"Operations include getting, setting, and resetting quotas. " +
		"Requires super-user permissions for all operations."

	return &sdk.Tool{
		Name:        "pulsar_admin_resourcequota",
		Description: toolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildResourceQuotasHandler builds the Pulsar admin resource quotas handler function
func (b *PulsarAdminResourceQuotasToolBuilder) buildResourceQuotasHandler(readOnly bool) builders.ToolHandlerFunc[pulsarAdminResourceQuotasInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input pulsarAdminResourceQuotasInput) (*sdk.CallToolResult, any, error) {
		resource := strings.ToLower(input.Resource)
		if resource == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'resource'; please specify: quota")
		}

		operation := strings.ToLower(input.Operation)
		if operation == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'operation'; please specify: get, set, reset")
		}

		// Validate write operations in read-only mode
		if readOnly && (operation == "set" || operation == "reset") {
			return nil, nil, fmt.Errorf("write operations are not allowed in read-only mode")
		}

		// Verify resource type
		if resource != "quota" {
			return nil, nil, fmt.Errorf("invalid resource: %s. available resources: quota", resource)
		}

		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return nil, nil, fmt.Errorf("pulsar session not found in context")
		}

		admin, err := session.GetAdminClient()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get admin client: %v", err)
		}

		// Dispatch based on operation
		switch operation {
		case "get":
			result, handlerErr := b.handleQuotaGet(admin, input)
			return result, nil, handlerErr
		case "set":
			result, handlerErr := b.handleQuotaSet(admin, input)
			return result, nil, handlerErr
		case "reset":
			result, handlerErr := b.handleQuotaReset(admin, input)
			return result, nil, handlerErr
		default:
			return nil, nil, fmt.Errorf("invalid operation: %s. available operations: get, set, reset", operation)
		}
	}
}

// handleQuotaGet handles getting a resource quota
func (b *PulsarAdminResourceQuotasToolBuilder) handleQuotaGet(admin cmdutils.Client, input pulsarAdminResourceQuotasInput) (*sdk.CallToolResult, error) {
	namespace := ""
	if input.Namespace != nil {
		namespace = *input.Namespace
	}
	bundle := ""
	if input.Bundle != nil {
		bundle = *input.Bundle
	}

	// Check if both namespace and bundle are provided or neither is provided
	if (namespace != "" && bundle == "") || (namespace == "" && bundle != "") {
		return nil, fmt.Errorf("when specifying a namespace, you must also specify a bundle and vice versa")
	}

	var (
		resourceQuotaData *utils.ResourceQuota
		getErr            error
	)

	if namespace == "" && bundle == "" {
		// Get default resource quota
		resourceQuotaData, getErr = admin.ResourceQuotas().GetDefaultResourceQuota()
		if getErr != nil {
			return nil, fmt.Errorf("failed to get default resource quota: %v", getErr)
		}
	} else {
		// Get namespace bundle resource quota
		nsName, getErr := utils.GetNamespaceName(namespace)
		if getErr != nil {
			return nil, fmt.Errorf("invalid namespace name '%s': %v", namespace, getErr)
		}
		resourceQuotaData, getErr = admin.ResourceQuotas().GetNamespaceBundleResourceQuota(nsName.String(), bundle)
		if getErr != nil {
			return nil, fmt.Errorf("failed to get resource quota for namespace '%s' bundle '%s': %v",
				namespace, bundle, getErr)
		}
	}

	// Format the output
	jsonBytes, err := json.Marshal(resourceQuotaData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource quota data: %v", err)
	}

	return textResult(string(jsonBytes)), nil
}

// handleQuotaSet handles setting a resource quota
func (b *PulsarAdminResourceQuotasToolBuilder) handleQuotaSet(admin cmdutils.Client, input pulsarAdminResourceQuotasInput) (*sdk.CallToolResult, error) {
	msgRateIn, err := requireFloat(input.MsgRateIn, "msgRateIn")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'msgRateIn' for quota.set: %v", err)
	}

	msgRateOut, err := requireFloat(input.MsgRateOut, "msgRateOut")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'msgRateOut' for quota.set: %v", err)
	}

	bandwidthIn, err := requireFloat(input.BandwidthIn, "bandwidthIn")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'bandwidthIn' for quota.set: %v", err)
	}

	bandwidthOut, err := requireFloat(input.BandwidthOut, "bandwidthOut")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'bandwidthOut' for quota.set: %v", err)
	}

	memory, err := requireFloat(input.Memory, "memory")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'memory' for quota.set: %v", err)
	}

	// Get optional parameters
	namespace := ""
	if input.Namespace != nil {
		namespace = *input.Namespace
	}
	bundle := ""
	if input.Bundle != nil {
		bundle = *input.Bundle
	}
	dynamic := false
	if input.Dynamic != nil {
		dynamic = *input.Dynamic
	}

	// Check if both namespace and bundle are provided or neither is provided
	if (namespace != "" && bundle == "") || (namespace == "" && bundle != "") {
		return nil, fmt.Errorf("when specifying a namespace, you must also specify a bundle and vice versa")
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
			return nil, fmt.Errorf("failed to set default resource quota: %v", err)
		}
		resultMsg = "Default resource quota set successfully"
	} else {
		// Set namespace bundle resource quota
		err = admin.ResourceQuotas().SetNamespaceBundleResourceQuota(namespace, bundle, *quota)
		if err != nil {
			return nil, fmt.Errorf("failed to set resource quota for namespace '%s' bundle '%s': %v",
				namespace, bundle, err)
		}
		resultMsg = fmt.Sprintf("Resource quota for namespace '%s' bundle '%s' set successfully", namespace, bundle)
	}

	return textResult(resultMsg), nil
}

// handleQuotaReset handles resetting a resource quota
func (b *PulsarAdminResourceQuotasToolBuilder) handleQuotaReset(admin cmdutils.Client, input pulsarAdminResourceQuotasInput) (*sdk.CallToolResult, error) {
	namespace, err := requireString(input.Namespace, "namespace")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'namespace' for quota.reset: %v", err)
	}

	bundle, err := requireString(input.Bundle, "bundle")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'bundle' for quota.reset: %v", err)
	}

	// Parse namespace name
	nsName, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return nil, fmt.Errorf("invalid namespace name '%s': %v", namespace, err)
	}

	// Reset namespace bundle resource quota
	err = admin.ResourceQuotas().ResetNamespaceBundleResourceQuota(nsName.String(), bundle)
	if err != nil {
		return nil, fmt.Errorf("failed to reset resource quota for namespace '%s' bundle '%s': %v",
			namespace, bundle, err)
	}

	return textResult(fmt.Sprintf("Resource quota for namespace '%s' bundle '%s' reset to default successfully",
		namespace, bundle)), nil
}

func buildPulsarAdminResourceQuotasInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[pulsarAdminResourceQuotasInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}
	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}

	setSchemaDescription(schema, "resource", pulsarAdminResourceQuotasResourceDesc)
	setSchemaDescription(schema, "operation", pulsarAdminResourceQuotasOperationDesc)
	setSchemaDescription(schema, "namespace", pulsarAdminResourceQuotasNamespaceDesc)
	setSchemaDescription(schema, "bundle", pulsarAdminResourceQuotasBundleDesc)
	setSchemaDescription(schema, "msgRateIn", pulsarAdminResourceQuotasMsgRateInDesc)
	setSchemaDescription(schema, "msgRateOut", pulsarAdminResourceQuotasMsgRateOutDesc)
	setSchemaDescription(schema, "bandwidthIn", pulsarAdminResourceQuotasBandwidthInDesc)
	setSchemaDescription(schema, "bandwidthOut", pulsarAdminResourceQuotasBandwidthOutDesc)
	setSchemaDescription(schema, "memory", pulsarAdminResourceQuotasMemoryDesc)
	setSchemaDescription(schema, "dynamic", pulsarAdminResourceQuotasDynamicDesc)

	normalizeAdditionalProperties(schema)
	return schema, nil
}

func requireFloat(value *float64, key string) (float64, error) {
	if value == nil {
		return 0, fmt.Errorf("required argument %q not found", key)
	}
	return *value, nil
}
