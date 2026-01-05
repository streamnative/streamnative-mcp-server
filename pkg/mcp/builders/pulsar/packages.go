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

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/adapter"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

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
func (b *PulsarAdminPackagesToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ServerTool, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool := b.buildPackagesTool()
	handler := b.buildPackagesHandler(config.ReadOnly)

	return []builders.ServerTool{
		{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildPackagesTool builds the Pulsar admin packages MCP tool definition
func (b *PulsarAdminPackagesToolBuilder) buildPackagesTool() *mcpsdk.Tool {
	toolDesc := "Manage packages in Apache Pulsar. Support package scheme: `function://`, `source://`, `sink://`" +
		"Allows listing, viewing, updating, downloading and uploading packages. " +
		"Some operations require super-user permissions."

	resourceDesc := "Resource to operate on. Available resources:\n" +
		"- package: A specific package\n" +
		"- packages: All packages of a specific type"

	operationDesc := "Operation to perform. Available operations:\n" +
		"- list: List all packages of a specific type or versions of a package\n" +
		"- get: Get metadata of a package\n" +
		"- update: Update metadata of a package (requires super-user permissions)\n" +
		"- delete: Delete a package (requires super-user permissions)\n" +
		"- download: Download a package (requires super-user permissions)\n" +
		"- upload: Upload a package (requires super-user permissions)"

	return builders.NewTool("pulsar_admin_package",
		builders.WithDescription(toolDesc),
		builders.WithString("resource", builders.Required(),
			builders.Description(resourceDesc),
		),
		builders.WithString("operation", builders.Required(),
			builders.Description(operationDesc),
		),
		builders.WithString("packageName",
			builders.Description("Name of the package to operate on. "+
				"Required for operations on a specific package: get, update, delete, download, upload"),
		),
		builders.WithString("namespace",
			builders.Description("The namespace name. Required for listing packages of a specific type"),
		),
		builders.WithString("type",
			builders.Description("Package type (function, source, sink). Required for listing packages of a specific type"),
		),
		builders.WithString("description",
			builders.Description("Description of the package. Required for update and upload operations"),
		),
		builders.WithString("contact",
			builders.Description("Contact information for the package. Optional for update and upload operations"),
		),
		builders.WithString("path",
			builders.Description("Path to download a package to or upload a package from. Required for download and upload operations"),
		),
		builders.WithObject("properties",
			builders.Description("Additional properties for the package as key-value pairs. Optional for update and upload operations"),
		),
	)
}

// buildPackagesHandler builds the Pulsar admin packages handler function
func (b *PulsarAdminPackagesToolBuilder) buildPackagesHandler(readOnly bool) func(context.Context, *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	return func(ctx context.Context, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
		// Get required parameters
		resource, err := adapter.RequireString(request, "resource")
		if err != nil {
			return adapter.NewErrorResult("Failed to get resource: %v", err), nil
		}

		operation, err := adapter.RequireString(request, "operation")
		if err != nil {
			return adapter.NewErrorResult("Failed to get operation: %v", err), nil
		}

		// Normalize parameters
		resource = strings.ToLower(resource)
		operation = strings.ToLower(operation)

		// Validate write operations in read-only mode
		if readOnly && (operation == "update" || operation == "delete" || operation == "download" || operation == "upload") {
			return adapter.NewErrorResult("Write operations are not allowed in read-only mode"), nil
		}

		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return adapter.NewErrorResult("Pulsar session not found in context"), nil
		}

		client, err := session.GetAdminV3Client()
		if err != nil {
			return adapter.NewErrorResult("Failed to get Pulsar client: %v", err), nil
		}

		// Dispatch based on resource type
		switch resource {
		case "package":
			return b.handlePackageResource(client, operation, request)
		case "packages":
			return b.handlePackagesResource(client, operation, request)
		default:
			return adapter.NewErrorResult("Invalid resource: %s. Available resources: package, packages", resource), nil
		}
	}
}

// Helper functions

// handlePackageResource handles operations on a specific package
func (b *PulsarAdminPackagesToolBuilder) handlePackageResource(client cmdutils.Client, operation string, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	packageName, err := adapter.RequireString(request, "packageName")
	if err != nil {
		return adapter.NewErrorResult("Missing required parameter 'packageName' for package operations: %v", err), nil
	}

	switch operation {
	case "list":
		// Get package versions
		packageVersions, err := client.Packages().ListVersions(packageName)
		if err != nil {
			return adapter.NewErrorResult("Failed to list package versions: %v", err), nil
		}

		// Convert result to JSON string
		packageVersionsJSON, err := json.Marshal(packageVersions)
		if err != nil {
			return adapter.NewErrorResult("Failed to serialize package versions: %v", err), nil
		}

		return adapter.NewTextResult(string(packageVersionsJSON)), nil

	case "get":
		// Get package metadata
		metadata, err := client.Packages().GetMetadata(packageName)
		if err != nil {
			return adapter.NewErrorResult("Failed to get package metadata: %v", err), nil
		}

		// Convert result to JSON string
		metadataJSON, err := json.Marshal(metadata)
		if err != nil {
			return adapter.NewErrorResult("Failed to serialize package metadata: %v", err), nil
		}

		return adapter.NewTextResult(string(metadataJSON)), nil

	case "update":
		description, err := adapter.RequireString(request, "description")
		if err != nil {
			return adapter.NewErrorResult("Missing required parameter 'description' for package.update: %v", err), nil
		}

		contact := adapter.GetString(request, "contact", "")
		args, _ := adapter.GetArgumentsMap(request)
		properties := b.extractProperties(args)

		// Update package metadata
		err = client.Packages().UpdateMetadata(packageName, description, contact, properties)
		if err != nil {
			return adapter.NewErrorResult("Failed to update package metadata: %v", err), nil
		}

		return adapter.NewTextResult(fmt.Sprintf("The metadata of the package '%s' updated successfully", packageName)), nil

	case "delete":
		// Delete package
		err = client.Packages().Delete(packageName)
		if err != nil {
			return adapter.NewErrorResult("Failed to delete package: %v", err), nil
		}

		return adapter.NewTextResult(fmt.Sprintf("The package '%s' deleted successfully", packageName)), nil

	case "download":
		path, err := adapter.RequireString(request, "path")
		if err != nil {
			return adapter.NewErrorResult("Missing required parameter 'path' for package.download: %v", err), nil
		}

		// Download package
		err = client.Packages().Download(packageName, path)
		if err != nil {
			return adapter.NewErrorResult("Failed to download package: %v", err), nil
		}

		return adapter.NewTextResult(
			fmt.Sprintf("The package '%s' downloaded to path '%s' successfully", packageName, path),
		), nil

	case "upload":
		path, err := adapter.RequireString(request, "path")
		if err != nil {
			return adapter.NewErrorResult("Missing required parameter 'path' for package.upload: %v", err), nil
		}

		description, err := adapter.RequireString(request, "description")
		if err != nil {
			return adapter.NewErrorResult("Missing required parameter 'description' for package.upload: %v", err), nil
		}

		contact := adapter.GetString(request, "contact", "")
		args, _ := adapter.GetArgumentsMap(request)
		properties := b.extractProperties(args)

		// Upload package
		err = client.Packages().Upload(packageName, path, description, contact, properties)
		if err != nil {
			return adapter.NewErrorResult("Failed to upload package: %v", err), nil
		}

		return adapter.NewTextResult(
			fmt.Sprintf("The package '%s' uploaded from path '%s' successfully", packageName, path),
		), nil

	default:
		return adapter.NewErrorResult("Invalid operation for resource 'package': %s. Available operations: list, get, update, delete, download, upload", operation), nil
	}
}

// handlePackagesResource handles operations on multiple packages
func (b *PulsarAdminPackagesToolBuilder) handlePackagesResource(client cmdutils.Client, operation string, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	switch operation {
	case "list":
		packageType, err := adapter.RequireString(request, "type")
		if err != nil {
			return adapter.NewErrorResult("Missing required parameter 'type' for packages.list: %v", err), nil
		}

		namespace, err := adapter.RequireString(request, "namespace")
		if err != nil {
			return adapter.NewErrorResult("Missing required parameter 'namespace' for packages.list: %v", err), nil
		}

		// Get package list
		packages, err := client.Packages().List(packageType, namespace)
		if err != nil {
			return adapter.NewErrorResult("Failed to list packages: %v", err), nil
		}

		// Convert result to JSON string
		packagesJSON, err := json.Marshal(packages)
		if err != nil {
			return adapter.NewErrorResult("Failed to serialize package list: %v", err), nil
		}

		return adapter.NewTextResult(string(packagesJSON)), nil

	default:
		return adapter.NewErrorResult("Invalid operation for resource 'packages': %s. Available operations: list", operation), nil
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
