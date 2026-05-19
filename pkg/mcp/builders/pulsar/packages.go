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

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

var pulsarPackageOperationSpecs = builders.OperationRegistry{
	{Name: "list", Mode: builders.OperationModeRead},
	{Name: "get", Mode: builders.OperationModeRead},
	{Name: "download", Mode: builders.OperationModeRead},
	{Name: "update", Mode: builders.OperationModeWrite, Destructive: true},
	{Name: "delete", Mode: builders.OperationModeWrite, Destructive: true},
	{Name: "upload", Mode: builders.OperationModeWrite, Destructive: true},
}

// PulsarAdminPackagesToolBuilder implements the ToolBuilder interface for Pulsar admin packages
// /nolint:revive
type PulsarAdminPackagesToolBuilder struct {
	*builders.BaseToolBuilder
}

// NewPulsarAdminPackagesToolBuilder creates a new Pulsar admin packages tool builder instance
func NewPulsarAdminPackagesToolBuilder() *PulsarAdminPackagesToolBuilder {
	metadata := builders.ToolMetadata{
		Name:        "pulsar_admin_packages",
		Version:     "1.0.0",
		Description: "Pulsar admin packages management tools",
		Category:    "pulsar_admin",
		Tags:        []string{"pulsar", "admin", "packages"},
	}

	features := []string{
		"pulsar-admin-packages",
		"all",
		"all-pulsar",
		"pulsar-admin",
	}

	return &PulsarAdminPackagesToolBuilder{
		BaseToolBuilder: builders.NewBaseToolBuilder(metadata, features),
	}
}

// BuildTools builds the Pulsar admin packages tool list
func (b *PulsarAdminPackagesToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]server.ServerTool, error) {
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
			Tool:    b.buildPackagesTool(toolModeRead),
			Handler: b.buildPackagesHandler(toolModeRead),
		},
	}
	if !config.ReadOnly {
		tools = append(tools, server.ServerTool{
			Tool:    b.buildPackagesTool(toolModeWrite),
			Handler: b.buildPackagesHandler(toolModeWrite),
		})
	}

	return tools, nil
}

// buildPackagesTool builds the Pulsar admin packages MCP tool definition
func (b *PulsarAdminPackagesToolBuilder) buildPackagesTool(mode toolMode) mcp.Tool {
	toolDesc := "Read packages in Apache Pulsar. Support package schemes: `function://`, `source://`, `sink://`. " +
		"Allows listing, viewing, and downloading packages."

	resourceDesc := "Resource to operate on. Available resources:\n" +
		"- package: A specific package\n" +
		"- packages: All packages of a specific type"
	resourceEnum := []string{"package", "packages"}

	operationDesc := "Operation to perform. Available operations:\n" +
		"- list: List all packages of a specific type or versions of a package\n" +
		"- get: Get metadata of a package\n" +
		"- download: Download a package"

	operationEnum := pulsarPackageOperationSpecs.NamesForMode(mode)
	toolName := "pulsar_admin_package_read"
	annotation := builders.ToolAnnotationForMode(mode, "Read Pulsar Packages", "Manage Pulsar Packages")
	if isToolModeWrite(mode) {
		toolDesc = "Manage packages in Apache Pulsar. Support package schemes: `function://`, `source://`, `sink://`. " +
			"This write tool updates metadata, deletes packages, or uploads package contents."
		resourceDesc = "Resource to operate on. Available resources:\n" +
			"- package: A specific package"
		resourceEnum = []string{"package"}
		operationDesc = "Operation to perform. Available operations:\n" +
			"- update: Update metadata of a package (requires super-user permissions)\n" +
			"- delete: Delete a package (requires super-user permissions)\n" +
			"- upload: Upload a package (requires super-user permissions)"
		operationEnum = pulsarPackageOperationSpecs.NamesForMode(mode)
		toolName = "pulsar_admin_package_write"
	}

	tool := mcp.NewTool(toolName,
		mcp.WithDescription(toolDesc),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description(resourceDesc),
			mcp.Enum(resourceEnum...),
		),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
			mcp.Enum(operationEnum...),
		),
		mcp.WithString("packageName",
			mcp.Description("Name of the package to operate on. Required for operations that target one specific package."),
		),
		mcp.WithString("namespace",
			mcp.Description("The namespace name. Required for listing packages of a specific type"),
		),
		mcp.WithString("type",
			mcp.Description("Package type (function, source, sink). Required for listing packages of a specific type"),
		),
		mcp.WithString("description",
			mcp.Description("Description of the package. Required for update and upload operations"),
		),
		mcp.WithString("contact",
			mcp.Description("Contact information for the package. Optional for update and upload operations"),
		),
		mcp.WithString("path",
			mcp.Description("Filesystem path used by package transfer operations. For downloads, this is the destination path."),
		),
		mcp.WithObject("properties",
			mcp.Description("Additional properties for the package as key-value pairs. Optional for update and upload operations"),
		),
		annotation,
	)
	if isToolModeWrite(mode) {
		pruneToolInputSchema(&tool, []string{"resource", "operation", "packageName", "description", "contact", "path", "properties"})
	} else {
		pruneToolInputSchema(&tool, []string{"resource", "operation", "packageName", "namespace", "type", "path"})
	}
	return tool
}

// buildPackagesHandler builds the Pulsar admin packages handler function
func (b *PulsarAdminPackagesToolBuilder) buildPackagesHandler(mode toolMode) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		if err := validateModeOperation(mode, operation, pulsarPackageOperationSpecs); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return mcp.NewToolResultError("Pulsar session not found in context"), nil
		}

		client, err := session.GetAdminV3Client()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get Pulsar client: %v", err)), nil
		}

		// Dispatch based on resource type
		switch resource {
		case "package":
			return b.handlePackageResource(client, operation, request, mode)
		case "packages":
			return b.handlePackagesResource(client, operation, request)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resource: %s. Available resources: package, packages", resource)), nil
		}
	}
}

// Helper functions

// handlePackageResource handles operations on a specific package
func (b *PulsarAdminPackagesToolBuilder) handlePackageResource(client cmdutils.Client, operation string, request mcp.CallToolRequest, mode toolMode) (*mcp.CallToolResult, error) {
	packageName, err := request.RequireString("packageName")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'packageName' for package operations: %v", err)), nil
	}

	switch operation {
	case "list":
		// Get package versions
		packageVersions, err := client.Packages().ListVersions(packageName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list package versions: %v", err)), nil
		}

		// Convert result to JSON string
		packageVersionsJSON, err := json.Marshal(packageVersions)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize package versions: %v", err)), nil
		}

		return mcp.NewToolResultText(string(packageVersionsJSON)), nil

	case "get":
		// Get package metadata
		metadata, err := client.Packages().GetMetadata(packageName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get package metadata: %v", err)), nil
		}

		// Convert result to JSON string
		metadataJSON, err := json.Marshal(metadata)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize package metadata: %v", err)), nil
		}

		return mcp.NewToolResultText(string(metadataJSON)), nil

	case "update":
		description, err := request.RequireString("description")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'description' for package.update: %v", err)), nil
		}

		contact := request.GetString("contact", "")
		properties := b.extractProperties(request.GetArguments())

		// Update package metadata
		err = client.Packages().UpdateMetadata(packageName, description, contact, properties)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update package metadata: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("The metadata of the package '%s' updated successfully", packageName)), nil

	case "delete":
		// Delete package
		err = client.Packages().Delete(packageName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete package: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("The package '%s' deleted successfully", packageName)), nil

	case "download":
		path, err := request.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'path' for package.download: %v", err)), nil
		}

		// Download package
		err = client.Packages().Download(packageName, path)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to download package: %v", err)), nil
		}

		return mcp.NewToolResultText(
			fmt.Sprintf("The package '%s' downloaded to path '%s' successfully", packageName, path),
		), nil

	case "upload":
		path, err := request.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'path' for package.upload: %v", err)), nil
		}

		description, err := request.RequireString("description")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'description' for package.upload: %v", err)), nil
		}

		contact := request.GetString("contact", "")
		properties := b.extractProperties(request.GetArguments())

		// Upload package
		err = client.Packages().Upload(packageName, path, description, contact, properties)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to upload package: %v", err)), nil
		}

		return mcp.NewToolResultText(
			fmt.Sprintf("The package '%s' uploaded from path '%s' successfully", packageName, path),
		), nil

	default:
		return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'package': %s. Available operations: %s", operation,
			modeSupportedOperations(mode, pulsarPackageOperationSpecs))), nil
	}
}

// handlePackagesResource handles operations on multiple packages
func (b *PulsarAdminPackagesToolBuilder) handlePackagesResource(client cmdutils.Client, operation string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	switch operation {
	case "list":
		packageType, err := request.RequireString("type")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'type' for packages.list: %v", err)), nil
		}

		namespace, err := request.RequireString("namespace")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'namespace' for packages.list: %v", err)), nil
		}

		// Get package list
		packages, err := client.Packages().List(packageType, namespace)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list packages: %v", err)), nil
		}

		// Convert result to JSON string
		packagesJSON, err := json.Marshal(packages)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize package list: %v", err)), nil
		}

		return mcp.NewToolResultText(string(packagesJSON)), nil

	default:
		return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'packages': %s. Available operations: list", operation)), nil
	}
}

// extractProperties extracts properties from request arguments
func (b *PulsarAdminPackagesToolBuilder) extractProperties(args map[string]interface{}) map[string]string {
	var properties map[string]string
	propsObj, ok := args["properties"]
	if ok && propsObj != nil {
		// Convert to map[string]string
		if propsMap, isMap := propsObj.(map[string]interface{}); isMap {
			properties = make(map[string]string)
			for k, v := range propsMap {
				if strVal, ok := v.(string); ok {
					properties[k] = strVal
				} else {
					properties[k] = fmt.Sprintf("%v", v)
				}
			}
		}
	}
	return properties
}
