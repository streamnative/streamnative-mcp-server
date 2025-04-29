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

// PulsarAdminAddTenantTools adds tenant-related tools to the MCP server
func PulsarAdminAddTenantTools(s *server.MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeaturePulsarAdminTenants)) && !slices.Contains(features, string(FeatureAll)) && !slices.Contains(features, string(FeatureAllPulsar)) {
		return
	}

	// Add unified tenant management tool
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

	tenantTool := mcp.NewTool("pulsar_admin_tenant",
		mcp.WithDescription(toolDesc),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description(resourceDesc),
		),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
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
		),
		mcp.WithArray("allowedClusters",
			mcp.Description("List of clusters that this tenant can access. Required for 'create' and 'update' operations. "+
				"Restricts the tenant to only use specified clusters, enabling geographic or infrastructure isolation. "+
				"Format: array of cluster names, e.g., ['us-west', 'us-east']. "+
				"An empty list means no clusters are accessible to this tenant."),
		),
	)
	s.AddTool(tenantTool, handleTenantTool(readOnly))
}

// handleTenantTool returns a function that handles tenant operations
func handleTenantTool(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		if readOnly && (operation == "create" || operation == "update" || operation == "delete") {
			return mcp.NewToolResultError("Write operations are not allowed in read-only mode"), nil
		}

		// Verify resource type
		if resource != "tenant" {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resource: %s. Only 'tenant' is supported", resource)), nil
		}

		// Create the admin client
		admin, err := pulsar.GetAdminClient()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get admin client: %v", err)), nil
		}

		// Dispatch based on operation
		switch operation {
		case "list":
			return handleTenantsList(admin)
		case "get":
			return handleTenantGet(admin, request)
		case "create":
			return handleTenantCreate(admin, request)
		case "update":
			return handleTenantUpdate(admin, request)
		case "delete":
			return handleTenantDelete(admin, request)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid operation: %s. Available operations: list, get, create, update, delete", operation)), nil
		}
	}
}

// handleTenantsList handles listing all tenants
func handleTenantsList(admin cmdutils.Client) (*mcp.CallToolResult, error) {
	// Get tenant list
	tenants, err := admin.Tenants().List()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list tenants: %v", err)), nil
	}

	// Convert result to JSON string
	tenantsJSON, err := json.Marshal(tenants)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to process tenant list: %v", err)), nil
	}

	return mcp.NewToolResultText(string(tenantsJSON)), nil
}

// handleTenantGet handles getting a tenant's configuration
func handleTenantGet(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tenant, err := requiredParam[string](request.Params.Arguments, "tenant")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'tenant' for tenant.get: %v", err)), nil
	}

	// Get tenant info
	tenantInfo, err := admin.Tenants().Get(tenant)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get tenant '%s' configuration: %v", tenant, err)), nil
	}

	// Convert result to JSON string
	tenantInfoJSON, err := json.Marshal(tenantInfo)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to process tenant info: %v", err)), nil
	}

	return mcp.NewToolResultText(string(tenantInfoJSON)), nil
}

// handleTenantCreate handles creating a new tenant
func handleTenantCreate(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tenant, err := requiredParam[string](request.Params.Arguments, "tenant")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'tenant' for tenant.create: %v", err)), nil
	}

	adminRoles, hasAdminRoles := optionalParamArray[string](request.Params.Arguments, "adminRoles")
	if !hasAdminRoles {
		adminRoles = []string{}
	}

	allowedClusters, hasAllowedClusters := optionalParamArray[string](request.Params.Arguments, "allowedClusters")
	if !hasAllowedClusters {
		allowedClusters = []string{""}
	}

	// Create tenant data struct
	tenantData := utils.TenantData{
		Name:            tenant,
		AdminRoles:      adminRoles,
		AllowedClusters: allowedClusters,
	}

	// Create tenant
	err = admin.Tenants().Create(tenantData)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create tenant '%s': %v", tenant, err)), nil
	}

	// Prepare result message with details
	adminRolesStr := "none"
	if len(adminRoles) > 0 {
		adminRolesStr = fmt.Sprintf("'%s'", strings.Join(adminRoles, "', '"))
	}

	allowedClustersStr := "none"
	if len(allowedClusters) > 0 && allowedClusters[0] != "" {
		allowedClustersStr = fmt.Sprintf("'%s'", strings.Join(allowedClusters, "', '"))
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully created tenant '%s' with admin roles: %s, allowed clusters: %s",
		tenant, adminRolesStr, allowedClustersStr)), nil
}

// handleTenantUpdate handles updating an existing tenant
func handleTenantUpdate(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tenant, err := requiredParam[string](request.Params.Arguments, "tenant")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'tenant' for tenant.update: %v", err)), nil
	}

	// Get current tenant info to ensure it exists and to show changes
	currentTenantInfo, err := admin.Tenants().Get(tenant)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to verify tenant '%s' exists: %v", tenant, err)), nil
	}

	// Get update parameters
	adminRoles, hasAdminRoles := optionalParamArray[string](request.Params.Arguments, "adminRoles")
	allowedClusters, hasAllowedClusters := optionalParamArray[string](request.Params.Arguments, "allowedClusters")

	// If parameters not provided, keep existing values
	if !hasAdminRoles {
		adminRoles = currentTenantInfo.AdminRoles
	}

	if !hasAllowedClusters {
		allowedClusters = currentTenantInfo.AllowedClusters
	}

	// Create tenant data struct with updated values
	tenantData := utils.TenantData{
		Name:            tenant,
		AdminRoles:      adminRoles,
		AllowedClusters: allowedClusters,
	}

	// Update tenant
	err = admin.Tenants().Update(tenantData)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to update tenant '%s': %v", tenant, err)), nil
	}

	// Prepare result message with details
	adminRolesStr := "none"
	if len(adminRoles) > 0 {
		adminRolesStr = fmt.Sprintf("'%s'", strings.Join(adminRoles, "', '"))
	}

	allowedClustersStr := "none"
	if len(allowedClusters) > 0 && allowedClusters[0] != "" {
		allowedClustersStr = fmt.Sprintf("'%s'", strings.Join(allowedClusters, "', '"))
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully updated tenant '%s' with admin roles: %s, allowed clusters: %s",
		tenant, adminRolesStr, allowedClustersStr)), nil
}

// handleTenantDelete handles deleting an existing tenant
func handleTenantDelete(admin cmdutils.Client, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tenant, err := requiredParam[string](request.Params.Arguments, "tenant")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'tenant' for tenant.delete: %v", err)), nil
	}

	// Delete tenant
	err = admin.Tenants().Delete(tenant)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete tenant '%s': %v. Note that a tenant cannot be deleted if it has active namespaces.", tenant, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully deleted tenant '%s'", tenant)), nil
}
