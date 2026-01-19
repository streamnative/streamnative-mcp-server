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

type pulsarAdminNsIsolationPolicyInput struct {
	Resource                 string         `json:"resource"`
	Operation                string         `json:"operation"`
	Cluster                  string         `json:"cluster"`
	Name                     *string        `json:"name,omitempty"`
	Namespaces               []string       `json:"namespaces,omitempty"`
	Primary                  []string       `json:"primary,omitempty"`
	Secondary                []string       `json:"secondary,omitempty"`
	AutoFailoverPolicyType   *string        `json:"autoFailoverPolicyType,omitempty"`
	AutoFailoverPolicyParams map[string]any `json:"autoFailoverPolicyParams,omitempty"`
}

const (
	pulsarAdminNsIsolationPolicyResourceDesc = "Resource to operate on. Available resources:\n" +
		"- policy: Namespace isolation policy\n" +
		"- broker: Broker with namespace isolation policies\n" +
		"- brokers: All brokers with namespace isolation policies"
	pulsarAdminNsIsolationPolicyOperationDesc = "Operation to perform. Available operations:\n" +
		"- get: Get resource details\n" +
		"- list: List all instances of the resource\n" +
		"- set: Create or update a resource (requires super-user permissions)\n" +
		"- delete: Delete a resource (requires super-user permissions)"
	pulsarAdminNsIsolationPolicyClusterDesc = "Cluster name"
	pulsarAdminNsIsolationPolicyNameDesc    = "Name of the policy or broker to operate on, based on resource type.\n" +
		"Required for: policy.get, policy.delete, policy.set, broker.get"
	pulsarAdminNsIsolationPolicyNamespacesDesc = "List of namespaces to apply the isolation policy. Required for policy.set"
	pulsarAdminNsIsolationPolicyPrimaryDesc    = "List of primary brokers for the namespaces. Required for policy.set"
	pulsarAdminNsIsolationPolicySecondaryDesc  = "List of secondary brokers for the namespaces. Optional for policy.set"
	pulsarAdminNsIsolationPolicyTypeDesc       = "Auto failover policy type (e.g., min_available). Optional for policy.set"
	pulsarAdminNsIsolationPolicyParamsDesc     = "Auto failover policy parameters as an object (e.g., {'min_limit': '1', 'usage_threshold': '100'}). Optional for policy.set"
)

// PulsarAdminNsIsolationPolicyToolBuilder implements the ToolBuilder interface for Pulsar admin namespace isolation policies
// /nolint:revive
type PulsarAdminNsIsolationPolicyToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminNsIsolationPolicyToolBuilder creates a new Pulsar admin namespace isolation policy tool builder instance
func NewPulsarAdminNsIsolationPolicyToolBuilder() *PulsarAdminNsIsolationPolicyToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_nsisolationpolicy",
		Version:     "1.0.0",
		Description: "Pulsar admin namespace isolation policy management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "admin", "nsisolationpolicy"},
	}

	features := []string{
		"pulsar-admin-nsisolationpolicy",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminNsIsolationPolicyToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin namespace isolation policy tool list
func (b *PulsarAdminNsIsolationPolicyToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	tool, err := b.buildNsIsolationPolicyTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildNsIsolationPolicyHandler(config.ReadOnly)

	return []builders.ToolDefinition{
		builders.ServerTool[pulsarAdminNsIsolationPolicyInput, any]{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildNsIsolationPolicyTool builds the Pulsar admin namespace isolation policy MCP tool definition
func (b *PulsarAdminNsIsolationPolicyToolBuilder) buildNsIsolationPolicyTool() (*sdk.Tool, error) {
	inputSchema, err := buildPulsarAdminNsIsolationPolicyInputSchema()
	if err != nil {
		return nil, err
	}

	toolDesc := "Manage namespace isolation policies in a Pulsar cluster. " +
		"Allows viewing, creating, updating, and deleting namespace isolation policies. " +
		"Some operations require super-user permissions."

	return &sdk.Tool{
		Name:        "pulsar_admin_nsisolationpolicy",
		Description: toolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildNsIsolationPolicyHandler builds the Pulsar admin namespace isolation policy handler function
func (b *PulsarAdminNsIsolationPolicyToolBuilder) buildNsIsolationPolicyHandler(readOnly bool) builders.ToolHandlerFunc[pulsarAdminNsIsolationPolicyInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input pulsarAdminNsIsolationPolicyInput) (*sdk.CallToolResult, any, error) {
		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return nil, nil, fmt.Errorf("pulsar session not found in context")
		}

		client, err := session.GetAdminClient()
		if err != nil {
			return nil, nil, b.handleError("get admin client", err)
		}

		resource := strings.ToLower(input.Resource)
		if resource == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'resource'; please specify one of: policy, broker, brokers")
		}

		operation := strings.ToLower(input.Operation)
		if operation == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'operation'; please specify one of: get, list, set, delete")
		}

		cluster := input.Cluster
		if cluster == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'cluster'")
		}

		// Validate write operations in read-only mode
		if readOnly && (operation == "set" || operation == "delete") {
			return nil, nil, fmt.Errorf("write operations are not allowed in read-only mode")
		}

		// Dispatch based on resource type
		switch resource {
		case "policy":
			result, err := b.handlePolicyResource(client, operation, cluster, input)
			return result, nil, err
		case "broker":
			result, err := b.handleBrokerResource(client, operation, cluster, input)
			return result, nil, err
		case "brokers":
			result, err := b.handleNsIsolationBrokersResource(client, operation, cluster)
			return result, nil, err
		default:
			return nil, nil, fmt.Errorf("invalid resource: %s. available resources: policy, broker, brokers", resource)
		}
	}
}

// Helper functions

// handlePolicyResource handles operations on the "policy" resource
func (b *PulsarAdminNsIsolationPolicyToolBuilder) handlePolicyResource(client cmdutils.Client, operation, cluster string, input pulsarAdminNsIsolationPolicyInput) (*sdk.CallToolResult, error) {
	switch operation {
	case "get":
		name, err := requireString(input.Name, "name")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'name' for policy.get: %v", err)
		}

		policyInfo, err := client.NsIsolationPolicy().GetNamespaceIsolationPolicy(cluster, name)
		if err != nil {
			return nil, b.handleError("get namespace isolation policy", err)
		}

		return b.marshalResponse(policyInfo)
	case "list":
		policies, err := client.NsIsolationPolicy().GetNamespaceIsolationPolicies(cluster)
		if err != nil {
			return nil, b.handleError("list namespace isolation policies", err)
		}

		return b.marshalResponse(policies)
	case "delete":
		name, err := requireString(input.Name, "name")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'name' for policy.delete: %v", err)
		}

		err = client.NsIsolationPolicy().DeleteNamespaceIsolationPolicy(cluster, name)
		if err != nil {
			return nil, b.handleError("delete namespace isolation policy", err)
		}

		return textResult(fmt.Sprintf("Delete namespace isolation policy %s successfully", name)), nil
	case "set":
		name, err := requireString(input.Name, "name")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'name' for policy.set: %v", err)
		}

		namespaces, err := requireStringSlice(input.Namespaces, "namespaces")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'namespaces' for policy.set: %v", err)
		}

		primary, err := requireStringSlice(input.Primary, "primary")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'primary' for policy.set: %v", err)
		}

		secondary := input.Secondary
		autoFailoverPolicyType := ""
		if input.AutoFailoverPolicyType != nil {
			autoFailoverPolicyType = *input.AutoFailoverPolicyType
		}

		autoFailoverPolicyParams, err := b.extractAutoFailoverPolicyParams(input.AutoFailoverPolicyParams)
		if err != nil {
			return nil, err
		}

		nsIsolationData, err := utils.CreateNamespaceIsolationData(namespaces, primary, secondary,
			autoFailoverPolicyType, autoFailoverPolicyParams)
		if err != nil {
			return nil, b.handleError("create namespace isolation data", err)
		}

		err = client.NsIsolationPolicy().CreateNamespaceIsolationPolicy(cluster, name, *nsIsolationData)
		if err != nil {
			return nil, b.handleError("create/update namespace isolation policy", err)
		}

		return textResult(fmt.Sprintf("Create/Update namespace isolation policy %s successfully", name)), nil
	default:
		return nil, fmt.Errorf("invalid operation for resource 'policy': %s. available operations: get, list, delete, set", operation)
	}
}

// handleBrokerResource handles operations on the "broker" resource
func (b *PulsarAdminNsIsolationPolicyToolBuilder) handleBrokerResource(client cmdutils.Client, operation, cluster string, input pulsarAdminNsIsolationPolicyInput) (*sdk.CallToolResult, error) {
	switch operation {
	case "get":
		name, err := requireString(input.Name, "name")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'name' for broker.get: %v", err)
		}

		brokerInfo, err := client.NsIsolationPolicy().GetBrokerWithNamespaceIsolationPolicy(cluster, name)
		if err != nil {
			return nil, b.handleError("get broker with namespace isolation policy", err)
		}

		return b.marshalResponse(brokerInfo)
	default:
		return nil, fmt.Errorf("invalid operation for resource 'broker': %s. available operations: get", operation)
	}
}

// handleNsIsolationBrokersResource handles operations on the "brokers" resource for namespace isolation policies
func (b *PulsarAdminNsIsolationPolicyToolBuilder) handleNsIsolationBrokersResource(client cmdutils.Client, operation, cluster string) (*sdk.CallToolResult, error) {
	switch operation {
	case "list":
		brokersInfo, err := client.NsIsolationPolicy().GetBrokersWithNamespaceIsolationPolicy(cluster)
		if err != nil {
			return nil, b.handleError("get brokers with namespace isolation policy", err)
		}

		return b.marshalResponse(brokersInfo)
	default:
		return nil, fmt.Errorf("invalid operation for resource 'brokers': %s. available operations: list", operation)
	}
}

func (b *PulsarAdminNsIsolationPolicyToolBuilder) extractAutoFailoverPolicyParams(params map[string]any) (map[string]string, error) {
	if params == nil {
		return nil, fmt.Errorf("failed to extract autoFailoverPolicyParams")
	}

	converted := make(map[string]string)
	for key, value := range params {
		if strValue, ok := value.(string); ok {
			converted[key] = strValue
		}
	}
	return converted, nil
}

// Utility functions

func (b *PulsarAdminNsIsolationPolicyToolBuilder) handleError(operation string, err error) error {
	return fmt.Errorf("failed to %s: %v", operation, err)
}

func (b *PulsarAdminNsIsolationPolicyToolBuilder) marshalResponse(data interface{}) (*sdk.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, b.handleError("marshal response", err)
	}
	return textResult(string(jsonBytes)), nil
}

func buildPulsarAdminNsIsolationPolicyInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[pulsarAdminNsIsolationPolicyInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}
	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}

	setSchemaDescription(schema, "resource", pulsarAdminNsIsolationPolicyResourceDesc)
	setSchemaDescription(schema, "operation", pulsarAdminNsIsolationPolicyOperationDesc)
	setSchemaDescription(schema, "cluster", pulsarAdminNsIsolationPolicyClusterDesc)
	setSchemaDescription(schema, "name", pulsarAdminNsIsolationPolicyNameDesc)
	setSchemaDescription(schema, "namespaces", pulsarAdminNsIsolationPolicyNamespacesDesc)
	setSchemaDescription(schema, "primary", pulsarAdminNsIsolationPolicyPrimaryDesc)
	setSchemaDescription(schema, "secondary", pulsarAdminNsIsolationPolicySecondaryDesc)
	setSchemaDescription(schema, "autoFailoverPolicyType", pulsarAdminNsIsolationPolicyTypeDesc)
	setSchemaDescription(schema, "autoFailoverPolicyParams", pulsarAdminNsIsolationPolicyParamsDesc)

	if namespacesSchema := schema.Properties["namespaces"]; namespacesSchema != nil && namespacesSchema.Items != nil {
		namespacesSchema.Items.Description = "namespace"
	}
	if primarySchema := schema.Properties["primary"]; primarySchema != nil && primarySchema.Items != nil {
		primarySchema.Items.Description = "primary broker"
	}
	if secondarySchema := schema.Properties["secondary"]; secondarySchema != nil && secondarySchema.Items != nil {
		secondarySchema.Items.Description = "secondary broker"
	}

	normalizeAdditionalProperties(schema)
	return schema, nil
}
