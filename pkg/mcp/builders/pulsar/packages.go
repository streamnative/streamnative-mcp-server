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

	"github.com/google/jsonschema-go/jsonschema"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	mcpCtx "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

type pulsarAdminPackagesInput struct {
	Resource    string         `json:"resource"`
	Operation   string         `json:"operation"`
	PackageName *string        `json:"packageName,omitempty"`
	Namespace   *string        `json:"namespace,omitempty"`
	PackageType *string        `json:"type,omitempty"`
	Description *string        `json:"description,omitempty"`
	Contact     *string        `json:"contact,omitempty"`
	Path        *string        `json:"path,omitempty"`
	Properties  map[string]any `json:"properties,omitempty"`
}

const (
	pulsarAdminPackagesResourceDesc = "Resource to operate on. Available resources:\n" +
		"- package: A specific package\n" +
		"- packages: All packages of a specific type"
	pulsarAdminPackagesOperationDesc = "Operation to perform. Available operations:\n" +
		"- list: List all packages of a specific type or versions of a package\n" +
		"- get: Get metadata of a package\n" +
		"- update: Update metadata of a package (requires super-user permissions)\n" +
		"- delete: Delete a package (requires super-user permissions)\n" +
		"- download: Download a package (requires super-user permissions)\n" +
		"- upload: Upload a package (requires super-user permissions)"
	pulsarAdminPackagesPackageNameDesc = "Name of the package to operate on. " +
		"Required for operations on a specific package: get, update, delete, download, upload"
	pulsarAdminPackagesNamespaceDesc = "The namespace name. Required for listing packages of a specific type"
	pulsarAdminPackagesTypeDesc      = "Package type (function, source, sink). Required for listing packages of a specific type"
	pulsarAdminPackagesDescription   = "Description of the package. Required for update and upload operations"
	pulsarAdminPackagesContactDesc   = "Contact information for the package. Optional for update and upload operations"
	pulsarAdminPackagesPathDesc      = "Path to download a package to or upload a package from. Required for download and upload operations"
	pulsarAdminPackagesProperties    = "Additional properties for the package as key-value pairs. Optional for update and upload operations"
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
func (b *PulsarAdminPackagesToolBuilder) BuildTools(_ context.Context, config builders.ToolBuildConfig) ([]builders.ToolDefinition, error) {
	// Check features - return empty list if no required features are present
	if !b.HasAnyRequiredFeature(config.Features) {
		return nil, nil
	}

	// Validate configuration (only validate when matching features are present)
	if err := b.Validate(config); err != nil {
		return nil, err
	}

	// Build tools
	tool, err := b.buildPackagesTool()
	if err != nil {
		return nil, err
	}
	handler := b.buildPackagesHandler(config.ReadOnly)

	return []builders.ToolDefinition{
		builders.ServerTool[pulsarAdminPackagesInput, any]{
			Tool:    tool,
			Handler: handler,
		},
	}, nil
}

// buildPackagesTool builds the Pulsar admin packages MCP tool definition
func (b *PulsarAdminPackagesToolBuilder) buildPackagesTool() (*sdk.Tool, error) {
	inputSchema, err := buildPulsarAdminPackagesInputSchema()
	if err != nil {
		return nil, err
	}

	toolDesc := "Manage packages in Apache Pulsar. Support package scheme: `function://`, `source://`, `sink://`" +
		"Allows listing, viewing, updating, downloading and uploading packages. " +
		"Some operations require super-user permissions."

	return &sdk.Tool{
		Name:        "pulsar_admin_package",
		Description: toolDesc,
		InputSchema: inputSchema,
	}, nil
}

// buildPackagesHandler builds the Pulsar admin packages handler function
func (b *PulsarAdminPackagesToolBuilder) buildPackagesHandler(readOnly bool) builders.ToolHandlerFunc[pulsarAdminPackagesInput, any] {
	return func(ctx context.Context, _ *sdk.CallToolRequest, input pulsarAdminPackagesInput) (*sdk.CallToolResult, any, error) {
		resource := strings.ToLower(strings.TrimSpace(input.Resource))
		if resource == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'resource'; please specify one of: package, packages")
		}

		operation := strings.ToLower(strings.TrimSpace(input.Operation))
		if operation == "" {
			return nil, nil, fmt.Errorf("missing required parameter 'operation'; please specify one of: list, get, update, delete, download, upload")
		}

		// Validate write operations in read-only mode
		if readOnly && (operation == "update" || operation == "delete" || operation == "download" || operation == "upload") {
			return nil, nil, fmt.Errorf("write operations are not allowed in read-only mode")
		}

		// Get Pulsar session from context
		session := mcpCtx.GetPulsarSession(ctx)
		if session == nil {
			return nil, nil, fmt.Errorf("pulsar session not found in context")
		}

		client, err := session.GetAdminV3Client()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get Pulsar client: %v", err)
		}

		// Dispatch based on resource type
		switch resource {
		case "package":
			result, err := b.handlePackageResource(client, operation, input)
			return result, nil, err
		case "packages":
			result, err := b.handlePackagesResource(client, operation, input)
			return result, nil, err
		default:
			return nil, nil, fmt.Errorf("invalid resource: %s. Available resources: package, packages", resource)
		}
	}
}

// handlePackageResource handles operations on a specific package
func (b *PulsarAdminPackagesToolBuilder) handlePackageResource(client cmdutils.Client, operation string, input pulsarAdminPackagesInput) (*sdk.CallToolResult, error) {
	packageName, err := requireString(input.PackageName, "packageName")
	if err != nil {
		return nil, fmt.Errorf("missing required parameter 'packageName' for package operations: %v", err)
	}

	switch operation {
	case "list":
		// Get package versions
		packageVersions, err := client.Packages().ListVersions(packageName)
		if err != nil {
			return nil, fmt.Errorf("failed to list package versions: %v", err)
		}

		// Convert result to JSON string
		packageVersionsJSON, err := json.Marshal(packageVersions)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize package versions: %v", err)
		}

		return textResult(string(packageVersionsJSON)), nil

	case "get":
		// Get package metadata
		metadata, err := client.Packages().GetMetadata(packageName)
		if err != nil {
			return nil, fmt.Errorf("failed to get package metadata: %v", err)
		}

		// Convert result to JSON string
		metadataJSON, err := json.Marshal(metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize package metadata: %v", err)
		}

		return textResult(string(metadataJSON)), nil

	case "update":
		description, err := requireString(input.Description, "description")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'description' for package.update: %v", err)
		}

		contact := ""
		if input.Contact != nil {
			contact = *input.Contact
		}
		properties := b.extractProperties(input.Properties)

		// Update package metadata
		err = client.Packages().UpdateMetadata(packageName, description, contact, properties)
		if err != nil {
			return nil, fmt.Errorf("failed to update package metadata: %v", err)
		}

		return textResult(fmt.Sprintf("The metadata of the package '%s' updated successfully", packageName)), nil

	case "delete":
		// Delete package
		err = client.Packages().Delete(packageName)
		if err != nil {
			return nil, fmt.Errorf("failed to delete package: %v", err)
		}

		return textResult(fmt.Sprintf("The package '%s' deleted successfully", packageName)), nil

	case "download":
		path, err := requireString(input.Path, "path")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'path' for package.download: %v", err)
		}

		// Download package
		err = client.Packages().Download(packageName, path)
		if err != nil {
			return nil, fmt.Errorf("failed to download package: %v", err)
		}

		return textResult(
			fmt.Sprintf("The package '%s' downloaded to path '%s' successfully", packageName, path),
		), nil

	case "upload":
		path, err := requireString(input.Path, "path")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'path' for package.upload: %v", err)
		}

		description, err := requireString(input.Description, "description")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'description' for package.upload: %v", err)
		}

		contact := ""
		if input.Contact != nil {
			contact = *input.Contact
		}
		properties := b.extractProperties(input.Properties)

		// Upload package
		err = client.Packages().Upload(packageName, path, description, contact, properties)
		if err != nil {
			return nil, fmt.Errorf("failed to upload package: %v", err)
		}

		return textResult(
			fmt.Sprintf("The package '%s' uploaded from path '%s' successfully", packageName, path),
		), nil

	default:
		return nil, fmt.Errorf("invalid operation for resource 'package': %s. Available operations: list, get, update, delete, download, upload", operation)
	}
}

// handlePackagesResource handles operations on multiple packages
func (b *PulsarAdminPackagesToolBuilder) handlePackagesResource(client cmdutils.Client, operation string, input pulsarAdminPackagesInput) (*sdk.CallToolResult, error) {
	switch operation {
	case "list":
		packageType, err := requireString(input.PackageType, "type")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'type' for packages.list: %v", err)
		}

		namespace, err := requireString(input.Namespace, "namespace")
		if err != nil {
			return nil, fmt.Errorf("missing required parameter 'namespace' for packages.list: %v", err)
		}

		// Get package list
		packages, err := client.Packages().List(packageType, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to list packages: %v", err)
		}

		// Convert result to JSON string
		packagesJSON, err := json.Marshal(packages)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize package list: %v", err)
		}

		return textResult(string(packagesJSON)), nil

	default:
		return nil, fmt.Errorf("invalid operation for resource 'packages': %s. Available operations: list", operation)
	}
}

// extractProperties extracts properties from input arguments
func (b *PulsarAdminPackagesToolBuilder) extractProperties(props map[string]any) map[string]string {
	if len(props) == 0 {
		return nil
	}

	properties := make(map[string]string, len(props))
	for key, value := range props {
		switch typed := value.(type) {
		case string:
			properties[key] = typed
		default:
			properties[key] = fmt.Sprintf("%v", value)
		}
	}

	return properties
}

func buildPulsarAdminPackagesInputSchema() (*jsonschema.Schema, error) {
	schema, err := jsonschema.For[pulsarAdminPackagesInput](nil)
	if err != nil {
		return nil, fmt.Errorf("input schema: %w", err)
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("input schema must have type \"object\"")
	}
	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}

	setSchemaDescription(schema, "resource", pulsarAdminPackagesResourceDesc)
	setSchemaDescription(schema, "operation", pulsarAdminPackagesOperationDesc)
	setSchemaDescription(schema, "packageName", pulsarAdminPackagesPackageNameDesc)
	setSchemaDescription(schema, "namespace", pulsarAdminPackagesNamespaceDesc)
	setSchemaDescription(schema, "type", pulsarAdminPackagesTypeDesc)
	setSchemaDescription(schema, "description", pulsarAdminPackagesDescription)
	setSchemaDescription(schema, "contact", pulsarAdminPackagesContactDesc)
	setSchemaDescription(schema, "path", pulsarAdminPackagesPathDesc)
	setSchemaDescription(schema, "properties", pulsarAdminPackagesProperties)

	normalizeAdditionalProperties(schema)
	return schema, nil
}
