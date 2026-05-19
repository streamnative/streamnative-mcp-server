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
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/toolannotations"
)

var pulsarTenantWriteOperations = map[string]struct{}{
	"create": {},
	"update": {},
	"delete": {},
}

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
func (b *PulsarAdminTenantToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
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
			Tool:    b.buildTenantTool(toolModeRead),
			Handler: b.buildTenantHandler(toolModeRead),
		},
	}
	if !config.ReadOnly {
		tools = append(tools, server.ServerTool{
			Tool:    b.buildTenantTool(toolModeWrite),
			Handler: b.buildTenantHandler(toolModeWrite),
		})
	}

	return tools, nil
}

// buildTenantTool builds the Pulsar Admin Tenant MCP tool definition
// Migrated from the original tool definition logic
func (b *PulsarAdminTenantToolBuilder) buildTenantTool(mode toolMode) mcp.Tool {
	toolDesc := "Manage Apache Pulsar tenants. " +
		"Tenants are the highest level administrative unit in Pulsar's multi-tenancy hierarchy. " +
		"Each tenant can contain multiple namespaces, allowing for logical isolation of applications. " +
		"Tenant configuration controls admin access and cluster availability across organizations. " +
		"Tenants provide isolation boundaries for topics, security policies, and resource quotas. " +
		"Proper tenant management is essential for multi-tenant Pulsar deployments to ensure data isolation, " +
		"appropriate access controls, and effective resource sharing. " +
		"All tenant operations require super-user permissions."

	resourceDesc := "Resource to operate on. Available resources:\n" +
		"- tenant: A tenant in the Pulsar instance"

	operationDesc := "Operation to perform. Available operations:\n" +
		"- list: List all tenants in the Pulsar instance\n" +
		"- get: Get configuration details for a specific tenant\n" +
		"- create: Create a new tenant with specified configuration\n" +
		"- update: Update configuration for an existing tenant\n" +
		"- delete: Delete an existing tenant (must not have any active namespaces)"

	operationEnum := []string{"list", "get"}
	toolName := "pulsar_admin_tenant_read"
	annotation := toolannotations.ReadOnly("Read Pulsar Tenants")
	if isToolModeWrite(mode) {
		operationEnum = []string{"create", "update", "delete"}
		toolName = "pulsar_admin_tenant_write"
		annotation = toolannotations.Destructive("Manage Pulsar Tenants")
	}

	return mcp.NewTool(toolName,
		mcp.WithDescription(toolDesc),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description(resourceDesc),
		),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
			mcp.Enum(operationEnum...),
		),
		mcp.WithString("tenant",
			mcp.Description("The tenant name to operate on. Required for all operations except 'list'. "+
				"Tenant names are unique identifiers and form the root of the topic naming hierarchy. "+
				"A valid tenant name must be comprised of alphanumeric characters and/or the following special characters: "+
				"'-', '_', '.', ':'. Ensure the tenant name follows your organization's naming conventions."),
		),
		mcp.WithArray("adminRoles",
			mcp.Description("List of auth principals (users or roles) allowed to administrate the tenant. "+
				"Required for 'create' and 'update' operations. These roles can create, update, or delete any "+
				"namespaces within the tenant, and can manage topic configurations. "+
				"Format: array of role strings, e.g., ['admin1', 'orgAdmin']. "+
				"Use empty array [] to remove all admin roles."),
			mcp.Items(
				map[string]interface{}{
					"type":        "string",
					"description": "role",
				},
			),
		),
		mcp.WithArray("allowedClusters",
			mcp.Description("List of clusters that this tenant can access. Required for 'create' and 'update' operations. "+
				"Restricts the tenant to only use specified clusters, enabling geographic or infrastructure isolation. "+
				"Format: array of cluster names, e.g., ['us-west', 'us-east']. "+
				"An empty list means no clusters are accessible to this tenant."),
			mcp.Items(
				map[string]interface{}{
					"type":        "string",
					"description": "cluster",
				},
			),
		),
		annotation,
	)
}

// buildTenantHandler builds the Pulsar Admin Tenant handler function
// Migrated from the original handler logic
func (b *PulsarAdminTenantToolBuilder) buildTenantHandler(mode toolMode) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		// Validate resource
		if resource != "tenant" {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resource: %s. Only 'tenant' is supported.", resource)), nil
		}

		if !validateModeOperation(mode, operation, pulsarTenantWriteOperations) {
			return mcp.NewToolResultError(fmt.Sprintf("Operation %q is not available in %s mode", operation, mode)), nil
		}

		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return mcp.NewToolResultError("Pulsar session not found in context"), nil
		}

		// Create the admin client
		admin, err := session.GetAdminClient()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
		}

		// Dispatch based on operation
		switch operation {
		case "list":
			return b.handleTenantsList(admin)
		case "get":
			return b.handleTenantGet(admin, request)
		case "create":
			return b.handleTenantCreate(admin, request)
		case "update":
			return b.handleTenantUpdate(admin, request)
		case "delete":
			return b.handleTenantDelete(admin, request)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Unknown operation: %s", operation)), nil
		}
	}
}

// Unified error handling and utility functions

// handleError provides unified error handling
func (b *PulsarAdminTenantToolBuilder) handleError(operation string, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("Failed to %s: %v", operation, err))
}

// marshalResponse provides unified JSON serialization for responses
func (b *PulsarAdminTenantToolBuilder) marshalResponse(data interface{}) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return b.handleError("marshal response", err), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// Operation handler functions - migrated from the original implementation

// handleTenantsList handles listing all tenants
func (b *PulsarAdminTenantToolBuilder) handleTenantsList(admin cmdutils.Client) (*mcp.CallToolResult, error) {
	// Get tenants list
	tenants, err := admin.Tenants().List()
	if err != nil {
		return b.handleError("list tenants", err), nil
	}

	return b.marshalResponse(tenants)
}

// handleTenantGet handles getting tenant configuration
func (b *PulsarAdminTenantToolBuilder) handleTenantGet(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tenant, err := request.RequireString("tenant")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get tenant name: %v", err)), nil
	}

	// Get tenant info
	tenantInfo, err := admin.Tenants().Get(tenant)
	if err != nil {
		return b.handleError("get tenant", err), nil
	}

	return b.marshalResponse(tenantInfo)
}

// handleTenantCreate handles creating a new tenant
func (b *PulsarAdminTenantToolBuilder) handleTenantCreate(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tenant, err := request.RequireString("tenant")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get tenant name: %v", err)), nil
	}

	adminRoles, err := request.RequireStringSlice("adminRoles")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin roles: %v", err)), nil
	}

	allowedClusters, err := request.RequireStringSlice("allowedClusters")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get allowed clusters: %v", err)), nil
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
		return b.handleError("create tenant", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Tenant %s created successfully", tenant)), nil
}

// handleTenantUpdate handles updating tenant configuration
func (b *PulsarAdminTenantToolBuilder) handleTenantUpdate(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tenant, err := request.RequireString("tenant")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get tenant name: %v", err)), nil
	}

	adminRoles, err := request.RequireStringSlice("adminRoles")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin roles: %v", err)), nil
	}

	allowedClusters, err := request.RequireStringSlice("allowedClusters")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get allowed clusters: %v", err)), nil
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
		return b.handleError("update tenant", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Tenant %s updated successfully", tenant)), nil
}

// handleTenantDelete handles deleting a tenant
func (b *PulsarAdminTenantToolBuilder) handleTenantDelete(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tenant, err := request.RequireString("tenant")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get tenant name: %v", err)), nil
	}

	// Delete tenant
	err = admin.Tenants().Delete(tenant)
	if err != nil {
		return b.handleError("delete tenant", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Tenant %s deleted successfully", tenant)), nil
}
