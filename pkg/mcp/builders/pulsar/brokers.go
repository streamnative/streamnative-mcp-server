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

// Package pulsar provides MCP tool builders for Pulsar admin operations.
package pulsar

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

type pulsarAdminBrokersInput struct {
	Resource    string  `json:"resource"`
	Operation   string  `json:"operation"`
	ClusterName *string `json:"clusterName,omitempty"`
	BrokerURL   *string `json:"brokerUrl,omitempty"`
	ConfigType  *string `json:"configType,omitempty"`
	ConfigName  *string `json:"configName,omitempty"`
	ConfigValue *string `json:"configValue,omitempty"`
}

const (
	pulsarAdminBrokersResourceDesc = "Type of resource to access, available options:\n" +
		"- brokers: Manage broker listings\n" +
		"- health: Check broker health status\n" +
		"- config: Manage broker configurations\n" +
		"- namespaces: Manage namespaces owned by a broker"
	pulsarAdminBrokersOperationDesc = "Operation to perform, available options:\n" +
		"- list: List resources (used with brokers)\n" +
		"- get: Retrieve resource information (used with health, config, namespaces)\n" +
		"- update: Update a resource (used with config)\n" +
		"- delete: Delete a resource (used with config)"
	pulsarAdminBrokersClusterNameDesc = "Pulsar cluster name, required for these operations:\n" +
		"- When resource=brokers, operation=list\n" +
		"- When resource=namespaces, operation=get"
	pulsarAdminBrokersBrokerURLDesc = "Broker URL, such as '127.0.0.1:8080', required for these operations:\n" +
		"- When resource=namespaces, operation=get"
	pulsarAdminBrokersConfigTypeDesc = "Configuration type, required when resource=config, operation=get, available options:\n" +
		"- dynamic: Get list of dynamically modifiable configuration names\n" +
		"- runtime: Get all runtime configurations (including static and dynamic configs)\n" +
		"- internal: Get internal configuration information\n" +
		"- all_dynamic: Get all dynamic configurations and their current values"
	pulsarAdminBrokersConfigNameDesc = "Configuration parameter name, required for these operations:\n" +
		"- When resource=config, operation=update\n" +
		"- When resource=config, operation=delete"
	pulsarAdminBrokersConfigValueDesc = "Configuration parameter value, required for these operations:\n" +
		"- When resource=config, operation=update"
)

// PulsarAdminBrokersToolBuilder implements the ToolBuilder interface for Pulsar admin brokers
// /nolint:revive
type PulsarAdminBrokersToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminBrokersToolBuilder creates a new Pulsar admin brokers tool builder instance
func NewPulsarAdminBrokersToolBuilder() *PulsarAdminBrokersToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_brokers",
		Version:     "1.0.0",
		Description: "Pulsar admin brokers management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "admin", "brokers"},
	}

	features := []string{
		"pulsar-admin-brokers",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminBrokersToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin brokers tool list
func (b *PulsarAdminBrokersToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool, err := b.buildPulsarAdminBrokersTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildPulsarAdminBrokersHandler(config.ReadOnly)

	return []builders.ToolDefinition{
		builders.ServerTool[pulsarAdminBrokersInput, any]{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildPulsarAdminBrokersTool builds the Pulsar admin brokers MCP tool definition
func (b *PulsarAdminBrokersToolBuilder) buildPulsarAdminBrokersTool() (*sdk.Tool, error) {
	inputSchema, err := buildPulsarAdminBrokersInputSchema()
	if err != nil {
		return nil, err
	}

	toolDesc := "Unified tool for managing Apache Pulsar broker resources. This tool integrates multiple broker management functions, including:\n" +
		"1. List active brokers in a cluster (resource=brokers, operation=list)\n" +
		"2. Check broker health status (resource=health, operation=get)\n" +
		"3. Manage broker configurations (resource=config, operation=get/update/delete)\n" +
		"4. View namespaces owned by a broker (resource=namespaces, operation=get)\n\n" +
		"Different functions are accessed by combining resource and operation parameters, with other parameters used selectively based on operation type.\n" +
		"Example: {\"resource\": \"config\", \"operation\": \"get\", \"configType\": \"dynamic\"} retrieves all dynamic configuration names.\n" +
		"This tool requires Pulsar super-user permissions."

	return &sdk.Tool{
		Name:        "pulsar_admin_brokers",
		Description: toolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildPulsarAdminBrokersHandler builds the Pulsar admin brokers handler function
func (b *PulsarAdminBrokersToolBuilder) buildPulsarAdminBrokersHandler(readOnly bool) builders.ToolHandlerFunc[pulsarAdminBrokersInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input pulsarAdminBrokersInput) (*sdk.CallToolResult, any, error) {
		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return nil, nil, fmt.Errorf("pulsar session not found in context")
		}

		// Get admin client
		client, err := session.GetAdminClient()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get admin client: %v", err)
		}

		// Get required parameters
		resource := input.Resource
		if resource == "" {
			return nil, nil, fmt.Errorf("missing required resource parameter. please specify one of: brokers, health, config, namespaces")
		}

		operation := input.Operation
		if operation == "" {
			return nil, nil, fmt.Errorf("missing required operation parameter. please specify one of: list, get, update, delete based on the resource type")
		}

		// Validate if the parameter combination is valid
		validCombination, errMsg := b.validateResourceOperation(resource, operation)
		if !validCombination {
			return nil, nil, fmt.Errorf("%s", errMsg)
		}

		// Process request based on resource type
		switch resource {
		case "brokers":
			result, err := b.handleBrokersResource(client, operation, input)
			return result, nil, err
		case "health":
			result, err := b.handleHealthResource(client, operation)
			return result, nil, err
		case "config":
			// Check write operation permissions
			if (operation == "update" || operation == "delete") && readOnly {
				return nil, nil, fmt.Errorf("configuration update/delete operations not allowed in read-only mode. please contact your administrator if you need to modify broker configurations")
			}
			result, err := b.handleConfigResource(client, operation, input)
			return result, nil, err
		case "namespaces":
			result, err := b.handleNamespacesResource(client, operation, input)
			return result, nil, err
		default:
			return nil, nil, fmt.Errorf("unsupported resource: %s. please use one of: brokers, health, config, namespaces", resource)
		}
	}
}

// Helper functions

// validateResourceOperation validates if the resource and operation combination is valid
func (b *PulsarAdminBrokersToolBuilder) validateResourceOperation(resource, operation string) (bool, string) {
	validCombinations := map[string][]string{
		"brokers":    {"list"},
		"health":     {"get"},
		"config":     {"get", "update", "delete"},
		"namespaces": {"get"},
	}

	if operations, ok := validCombinations[resource]; ok {
		for _, op := range operations {
			if op == operation {
				return true, ""
			}
		}
		return false, fmt.Sprintf("Invalid operation '%s' for resource '%s'. Valid operations are: %v",
			operation, resource, validCombinations[resource])
	}

	return false, fmt.Sprintf("Invalid resource '%s'. Valid resources are: brokers, health, config, namespaces", resource)
}

// handleBrokersResource handles brokers resource
func (b *PulsarAdminBrokersToolBuilder) handleBrokersResource(client cmdutils.Client, operation string, input pulsarAdminBrokersInput) (*sdk.CallToolResult, error) {
	switch operation {
	case "list":
		clusterName, err := requireString(input.ClusterName, "clusterName")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'clusterName'. please provide the name of the Pulsar cluster to list brokers for")
		}

		brokers, err := client.Brokers().GetActiveBrokers(clusterName)
		if err != nil {
			return nil, fmt.Errorf("failed to get active brokers: %v. please verify the cluster name and ensure the Pulsar service is running", err)
		}

		return b.marshalResponse(brokers)
	default:
		return nil, fmt.Errorf("unsupported operation '%s' for brokers resource. the only supported operation is 'list'", operation)
	}
}

// handleHealthResource handles health resource
func (b *PulsarAdminBrokersToolBuilder) handleHealthResource(client cmdutils.Client, operation string) (*sdk.CallToolResult, error) {
	switch operation {
	case "get":
		//nolint:staticcheck
		err := client.Brokers().HealthCheck()
		if err != nil {
			return nil, fmt.Errorf("broker health check failed: %v. the broker might be down or experiencing issues", err)
		}
		return textResult("ok"), nil
	default:
		return nil, fmt.Errorf("unsupported operation '%s' for health resource. the only supported operation is 'get'", operation)
	}
}

// handleConfigResource handles config resource
func (b *PulsarAdminBrokersToolBuilder) handleConfigResource(client cmdutils.Client, operation string, input pulsarAdminBrokersInput) (*sdk.CallToolResult, error) {
	switch operation {
	case "get":
		configType, err := requireString(input.ConfigType, "configType")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'configType'. please specify one of: dynamic, runtime, internal, all_dynamic")
		}

		var result interface{}
		var fetchErr error

		switch configType {
		case "dynamic":
			result, fetchErr = client.Brokers().GetDynamicConfigurationNames()
		case "runtime":
			result, fetchErr = client.Brokers().GetRuntimeConfigurations()
		case "internal":
			result, fetchErr = client.Brokers().GetInternalConfigurationData()
		case "all_dynamic":
			result, fetchErr = client.Brokers().GetAllDynamicConfigurations()
		default:
			return nil, fmt.Errorf("invalid config type: '%s'. valid types are: dynamic, runtime, internal, all_dynamic", configType)
		}

		if fetchErr != nil {
			return nil, fmt.Errorf("failed to get %s configuration: %v", configType, fetchErr)
		}

		return b.marshalResponse(result)

	case "update":
		configName, err := requireString(input.ConfigName, "configName")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'configName'. please provide the name of the configuration parameter to update")
		}

		configValue, err := requireString(input.ConfigValue, "configValue")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'configValue'. please provide the new value for the configuration parameter")
		}

		err = client.Brokers().UpdateDynamicConfiguration(configName, configValue)
		if err != nil {
			return nil, fmt.Errorf("failed to update configuration: %v. please verify the configuration name is valid and the value is of the correct type", err)
		}

		return textResult(fmt.Sprintf("Dynamic configuration '%s' updated successfully to '%s'", configName, configValue)), nil

	case "delete":
		configName, err := requireString(input.ConfigName, "configName")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'configName'. please provide the name of the configuration parameter to delete")
		}

		err = client.Brokers().DeleteDynamicConfiguration(configName)
		if err != nil {
			return nil, fmt.Errorf("failed to delete configuration: %v. please verify the configuration name is valid and exists", err)
		}

		return textResult(fmt.Sprintf("Dynamic configuration '%s' deleted successfully", configName)), nil

	default:
		return nil, fmt.Errorf("unsupported operation '%s' for config resource. supported operations are: get, update, delete", operation)
	}
}

// handleNamespacesResource handles namespaces resource
func (b *PulsarAdminBrokersToolBuilder) handleNamespacesResource(client cmdutils.Client, operation string, input pulsarAdminBrokersInput) (*sdk.CallToolResult, error) {
	switch operation {
	case "get":
		clusterName, err := requireString(input.ClusterName, "clusterName")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'clusterName'. please provide the name of the Pulsar cluster")
		}

		brokerURL, err := requireString(input.BrokerURL, "brokerUrl")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'brokerUrl'. please provide the URL of the broker (e.g., '127.0.0.1:8080')")
		}

		namespaces, err := client.Brokers().GetOwnedNamespaces(clusterName, brokerURL)
		if err != nil {
			return nil, fmt.Errorf("failed to get owned namespaces: %v. please verify the cluster name and broker URL are correct", err)
		}

		return b.marshalResponse(namespaces)

	default:
		return nil, fmt.Errorf("unsupported operation '%s' for namespaces resource. the only supported operation is 'get'", operation)
	}
}

func (b *PulsarAdminBrokersToolBuilder) marshalResponse(data interface{}) (*sdk.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %v", err)
	}
	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: string(jsonBytes)}},
	}, nil
}

func buildPulsarAdminBrokersInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[pulsarAdminBrokersInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}

	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}

	setSchemaDescription(schema, "resource", pulsarAdminBrokersResourceDesc)
	setSchemaDescription(schema, "operation", pulsarAdminBrokersOperationDesc)
	setSchemaDescription(schema, "clusterName", pulsarAdminBrokersClusterNameDesc)
	setSchemaDescription(schema, "brokerUrl", pulsarAdminBrokersBrokerURLDesc)
	setSchemaDescription(schema, "configType", pulsarAdminBrokersConfigTypeDesc)
	setSchemaDescription(schema, "configName", pulsarAdminBrokersConfigNameDesc)
	setSchemaDescription(schema, "configValue", pulsarAdminBrokersConfigValueDesc)

	normalizeAdditionalProperties(schema)
	return schema, nil
}
