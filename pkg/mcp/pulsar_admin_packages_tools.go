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

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// PulsarAdminAddPackagesTools adds package-related tools to the MCP server
func PulsarAdminAddPackagesTools(s *server.MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeaturePulsarAdminPackages)) && !slices.Contains(features, string(FeatureAll)) && !slices.Contains(features, string(FeatureAllPulsar)) {
		return
	}

	// Add unified package management tool
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

	packageTool := mcp.NewTool("pulsar_admin_package",
		mcp.WithDescription(toolDesc),
		mcp.WithString("resource", mcp.Required(),
			mcp.Description(resourceDesc),
		),
		mcp.WithString("operation", mcp.Required(),
			mcp.Description(operationDesc),
		),
		mcp.WithString("packageName",
			mcp.Description("Name of the package to operate on. "+
				"Required for operations on a specific package: get, update, delete, download, upload"),
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
			mcp.Description("Path to download a package to or upload a package from. Required for download and upload operations"),
		),
		mcp.WithObject("properties",
			mcp.Description("Additional properties for the package as key-value pairs. Optional for update and upload operations"),
		),
	)
	s.AddTool(packageTool, handlePackageTool(readOnly))
}

// handlePackageTool returns a function that handles package-related operations
func handlePackageTool(readOnly bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		if readOnly && (operation == "update" || operation == "delete" || operation == "download" || operation == "upload") {
			return mcp.NewToolResultError("Write operations are not allowed in read-only mode"), nil
		}

		// Create Pulsar client with API version V3
		client, err := pulsar.GetAdminV3Client()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get Pulsar client: %v", err)), nil
		}

		// Dispatch based on resource type
		switch resource {
		case "package":
			return handlePackageResource(client, operation, request)
		case "packages":
			return handlePackagesResource(client, operation, request)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resource: %s. Available resources: package, packages", resource)), nil
		}
	}
}

// handlePackageResource handles operations on a specific package
func handlePackageResource(client cmdutils.Client, operation string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	packageName, err := requiredParam[string](request.Params.Arguments, "packageName")
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
		description, err := requiredParam[string](request.Params.Arguments, "description")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'description' for package.update: %v", err)), nil
		}

		contact, _ := optionalParam[string](request.Params.Arguments, "contact")
		properties := extractProperties(request.Params.Arguments)

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
		path, err := requiredParam[string](request.Params.Arguments, "path")
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
		path, err := requiredParam[string](request.Params.Arguments, "path")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'path' for package.upload: %v", err)), nil
		}

		description, err := requiredParam[string](request.Params.Arguments, "description")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'description' for package.upload: %v", err)), nil
		}

		contact, _ := optionalParam[string](request.Params.Arguments, "contact")
		properties := extractProperties(request.Params.Arguments)

		// Upload package
		err = client.Packages().Upload(packageName, path, description, contact, properties)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to upload package: %v", err)), nil
		}

		return mcp.NewToolResultText(
			fmt.Sprintf("The package '%s' uploaded from path '%s' successfully", packageName, path),
		), nil

	default:
		return mcp.NewToolResultError(fmt.Sprintf("Invalid operation for resource 'package': %s. Available operations: list, get, update, delete, download, upload", operation)), nil
	}
}

// handlePackagesResource handles operations on multiple packages
func handlePackagesResource(client cmdutils.Client, operation string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	switch operation {
	case "list":
		packageType, err := requiredParam[string](request.Params.Arguments, "type")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'type' for packages.list: %v", err)), nil
		}

		namespace, err := requiredParam[string](request.Params.Arguments, "namespace")
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
func extractProperties(args map[string]interface{}) map[string]string {
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
