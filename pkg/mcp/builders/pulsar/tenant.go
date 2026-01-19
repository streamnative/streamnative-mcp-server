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

type pulsarAdminTenantInput struct {
	Resource        string   `json:"resource"`
	Operation       string   `json:"operation"`
	Tenant          *string  `json:"tenant,omitempty"`
	AdminRoles      []string `json:"adminRoles,omitempty"`
	AllowedClusters []string `json:"allowedClusters,omitempty"`
}

const (
	pulsarAdminTenantResourceDesc = "Resource to operate on. Available resources:\n" +
		"- tenant: A tenant in the Pulsar instance"
	pulsarAdminTenantOperationDesc = "Operation to perform. Available operations:\n" +
		"- list: List all tenants in the Pulsar instance\n" +
		"- get: Get configuration details for a specific tenant\n" +
		"- create: Create a new tenant with specified configuration\n" +
		"- update: Update configuration for an existing tenant\n" +
		"- delete: Delete an existing tenant (must not have any active namespaces)"
	pulsarAdminTenantNameDesc = "The tenant name to operate on. Required for all operations except 'list'. " +
		"Tenant names are unique identifiers and form the root of the topic naming hierarchy. " +
		"A valid tenant name must be comprised of alphanumeric characters and/or the following special characters: " +
		"'-', '_', '.', ':'. Ensure the tenant name follows your organization's naming conventions."
	pulsarAdminTenantAdminRolesDesc = "List of auth principals (users or roles) allowed to administrate the tenant. " +
		"Required for 'create' and 'update' operations. These roles can create, update, or delete any " +
		"namespaces within the tenant, and can manage topic configurations. " +
		"Format: array of role strings, e.g., ['admin1', 'orgAdmin']. " +
		"Use empty array [] to remove all admin roles."
	pulsarAdminTenantAllowedClustersDesc = "List of clusters that this tenant can access. Required for 'create' and 'update' operations. " +
		"Restricts the tenant to only use specified clusters, enabling geographic or infrastructure isolation. " +
		"Format: array of cluster names, e.g., ['us-west', 'us-east']. " +
		"An empty list means no clusters are accessible to this tenant."
)

// PulsarAdminTenantToolBuilder implements the ToolBuilder interface for Pulsar Admin Tenant tools
// It provides functionality to build Pulsar tenant management tools
// /nolint:revive
type PulsarAdminTenantToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminTenantToolBuilder creates a new Pulsar Admin Tenant tool builder instance
func NewPulsarAdminTenantToolBuilder() *PulsarAdminTenantToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_tenant",
		Version:     "1.0.0",
		Description: "Pulsar Admin tenant management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "tenant", "admin"},
	}

	features := []string{
		"pulsar-admin-tenants",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminTenantToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar Admin Tenant tool list
// This is the core method implementing the ToolBuilder interface
func (b *PulsarAdminTenantToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool, err := b.buildTenantTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildTenantHandler(config.ReadOnly)

	return []builders.ToolDefinition{
		builders.ServerTool[pulsarAdminTenantInput, any]{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildTenantTool builds the Pulsar Admin Tenant MCP tool definition
// Migrated from the original tool definition logic
func (b *PulsarAdminTenantToolBuilder) buildTenantTool() (*sdk.Tool, error) {
	inputSchema, err := buildPulsarAdminTenantInputSchema()
	if err != nil {
		return nil, err
	}

	toolDesc := "Manage Apache Pulsar tenants. " +
		"Tenants are the highest level administrative unit in Pulsar's multi-tenancy hierarchy. " +
		"Each tenant can contain multiple namespaces, allowing for logical isolation of applications. " +
		"Tenant configuration controls admin access and cluster availability across organizations. " +
		"Tenants provide isolation boundaries for topics, security policies, and resource quotas. " +
		"Proper tenant management is essential for multi-tenant Pulsar deployments to ensure data isolation, " +
		"appropriate access controls, and effective resource sharing. " +
		"All tenant operations require super-user permissions."

	return &sdk.Tool{
		Name:        "pulsar_admin_tenant",
		Description: toolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildTenantHandler builds the Pulsar Admin Tenant handler function
// Migrated from the original handler logic
func (b *PulsarAdminTenantToolBuilder) buildTenantHandler(readOnly bool) builders.ToolHandlerFunc[pulsarAdminTenantInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input pulsarAdminTenantInput) (*sdk.CallToolResult, any, error) {
		// Normalize parameters
		resource := strings.ToLower(input.Resource)
		operation := strings.ToLower(input.Operation)

		// Validate resource
		if resource != "tenant" {
			return nil, nil, fmt.Errorf("invalid resource: %s. only 'tenant' is supported", resource)
		}

		// Validate write operations in read-only mode
		if readOnly && (operation == "create" || operation == "update" || operation == "delete") {
			return nil, nil, fmt.Errorf("write operations are not allowed in read-only mode")
		}

		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return nil, nil, fmt.Errorf("pulsar session not found in context")
		}

		// Create the admin client
		admin, err := session.GetAdminClient()
		if err != nil {
			return nil, nil, b.handleError("get admin client", err)
		}

		// Dispatch based on operation
		switch operation {
		case "list":
			result, err := b.handleTenantsList(admin)
			return result, nil, err
		case "get":
			result, err := b.handleTenantGet(admin, input)
			return result, nil, err
		case "create":
			result, err := b.handleTenantCreate(admin, input)
			return result, nil, err
		case "update":
			result, err := b.handleTenantUpdate(admin, input)
			return result, nil, err
		case "delete":
			result, err := b.handleTenantDelete(admin, input)
			return result, nil, err
		default:
			return nil, nil, fmt.Errorf("unknown operation: %s", operation)
		}
	}
}

// Unified error handling and utility functions

// handleError provides unified error handling
func (b *PulsarAdminTenantToolBuilder) handleError(operation string, err error) error {
	return fmt.Errorf("failed to %s: %v", operation, err)
}

// marshalResponse provides unified JSON serialization for responses
func (b *PulsarAdminTenantToolBuilder) marshalResponse(data interface{}) (*sdk.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, b.handleError("marshal response", err)
	}
	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: string(jsonBytes)}},
	}, nil
}

// Operation handler functions - migrated from the original implementation

// handleTenantsList handles listing all tenants
func (b *PulsarAdminTenantToolBuilder) handleTenantsList(admin cmdutils.Client) (*sdk.CallToolResult, error) {
	// Get tenants list
	tenants, err := admin.Tenants().List()
	if err != nil {
		return nil, b.handleError("list tenants", err)
	}

	return b.marshalResponse(tenants)
}

// handleTenantGet handles getting tenant configuration
func (b *PulsarAdminTenantToolBuilder) handleTenantGet(admin cmdutils.Client, input pulsarAdminTenantInput) (*sdk.CallToolResult, error) {
	tenant, err := requireString(input.Tenant, "tenant")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'tenant' for tenant.get: %v", err)
	}

	// Get tenant info
	tenantInfo, err := admin.Tenants().Get(tenant)
	if err != nil {
		return nil, b.handleError("get tenant", err)
	}

	return b.marshalResponse(tenantInfo)
}

// handleTenantCreate handles creating a new tenant
func (b *PulsarAdminTenantToolBuilder) handleTenantCreate(admin cmdutils.Client, input pulsarAdminTenantInput) (*sdk.CallToolResult, error) {
	tenant, err := requireString(input.Tenant, "tenant")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'tenant' for tenant.create: %v", err)
	}

	adminRoles, err := requireStringSlice(input.AdminRoles, "adminRoles")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'adminRoles' for tenant.create: %v", err)
	}

	allowedClusters, err := requireStringSlice(input.AllowedClusters, "allowedClusters")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'allowedClusters' for tenant.create: %v", err)
	}

	// Create tenant data
	tenantData := utils.TenantData{
		Name:            tenant,
		AdminRoles:      adminRoles,
		AllowedClusters: allowedClusters,
	}

	// Create tenant
	err = admin.Tenants().Create(tenantData)
	if err != nil {
		return nil, b.handleError("create tenant", err)
	}

	return textResult(fmt.Sprintf("Tenant %s created successfully", tenant)), nil
}

// handleTenantUpdate handles updating tenant configuration
func (b *PulsarAdminTenantToolBuilder) handleTenantUpdate(admin cmdutils.Client, input pulsarAdminTenantInput) (*sdk.CallToolResult, error) {
	tenant, err := requireString(input.Tenant, "tenant")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'tenant' for tenant.update: %v", err)
	}

	adminRoles, err := requireStringSlice(input.AdminRoles, "adminRoles")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'adminRoles' for tenant.update: %v", err)
	}

	allowedClusters, err := requireStringSlice(input.AllowedClusters, "allowedClusters")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'allowedClusters' for tenant.update: %v", err)
	}

	// Create tenant data
	tenantData := utils.TenantData{
		Name:            tenant,
		AdminRoles:      adminRoles,
		AllowedClusters: allowedClusters,
	}

	// Update tenant
	err = admin.Tenants().Update(tenantData)
	if err != nil {
		return nil, b.handleError("update tenant", err)
	}

	return textResult(fmt.Sprintf("Tenant %s updated successfully", tenant)), nil
}

// handleTenantDelete handles deleting a tenant
func (b *PulsarAdminTenantToolBuilder) handleTenantDelete(admin cmdutils.Client, input pulsarAdminTenantInput) (*sdk.CallToolResult, error) {
	tenant, err := requireString(input.Tenant, "tenant")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'tenant' for tenant.delete: %v", err)
	}

	// Delete tenant
	err = admin.Tenants().Delete(tenant)
	if err != nil {
		return nil, b.handleError("delete tenant", err)
	}

	return textResult(fmt.Sprintf("Tenant %s deleted successfully", tenant)), nil
}

func buildPulsarAdminTenantInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[pulsarAdminTenantInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}

	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}
	setSchemaDescription(schema, "resource", pulsarAdminTenantResourceDesc)
	setSchemaDescription(schema, "operation", pulsarAdminTenantOperationDesc)
	setSchemaDescription(schema, "tenant", pulsarAdminTenantNameDesc)
	setSchemaDescription(schema, "adminRoles", pulsarAdminTenantAdminRolesDesc)
	setSchemaDescription(schema, "allowedClusters", pulsarAdminTenantAllowedClustersDesc)

	if adminRolesSchema := schema.Properties["adminRoles"]; adminRolesSchema != nil && adminRolesSchema.Items != nil {
		adminRolesSchema.Items.Description = "role"
	}
	if allowedClustersSchema := schema.Properties["allowedClusters"]; allowedClustersSchema != nil && allowedClustersSchema.Items != nil {
		allowedClustersSchema.Items.Description = "cluster"
	}

	normalizeAdditionalProperties(schema)
	return schema, nil
}

func requireStringSlice(values []string, key string) ([]string, error) {
	if values == nil {
		return nil, fmt.Errorf("required argument %q not found", key)
	}
	return values, nil
}
