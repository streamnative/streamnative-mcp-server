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

package pftools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/admin/config"
	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

const (
	CustomRuntimeOptionsEnvMcpToolNameKey        = "MCP_TOOL_NAME"
	CustomRuntimeOptionsEnvMcpToolDescriptionKey = "MCP_TOOL_DESCRIPTION"
)

// NewPulsarFunctionManager creates a new PulsarFunctionManager
func NewPulsarFunctionManager(mcpServer *server.MCPServer, readOnly bool, options *ManagerOptions) (*PulsarFunctionManager, error) {
	// Get Pulsar client and admin client
	pulsarClient, err := pulsar.GetPulsarClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get Pulsar client: %w", err)
	}

	adminClient := cmdutils.NewPulsarClientWithAPIVersion(config.V3)
	v2adminClient := cmdutils.NewPulsarClientWithAPIVersion(config.V2)
	if options == nil {
		options = DefaultManagerOptions()
	}

	// Create the manager
	manager := &PulsarFunctionManager{
		adminClient:       adminClient,
		v2adminClient:     v2adminClient,
		pulsarClient:      pulsarClient,
		fnToToolMap:       make(map[string]*FunctionTool),
		mutex:             sync.RWMutex{},
		pollInterval:      options.PollInterval,
		stopCh:            make(chan struct{}),
		callInProgressMap: make(map[string]context.CancelFunc),
		mcpServer:         mcpServer,
		readOnly:          readOnly,
		defaultTimeout:    options.DefaultTimeout,
		circuitBreakers:   make(map[string]*CircuitBreaker),
		tenantNamespaces:  options.TenantNamespaces,
	}

	return manager, nil
}

// Start starts polling for functions
func (m *PulsarFunctionManager) Start() {
	go m.pollFunctions()
}

// Stop stops polling for functions
func (m *PulsarFunctionManager) Stop() {
	close(m.stopCh)
}

// pollFunctions polls for functions periodically
func (m *PulsarFunctionManager) pollFunctions() {
	ticker := time.NewTicker(m.pollInterval)
	defer ticker.Stop()

	// Initial update
	m.updateFunctions()

	for {
		select {
		case <-ticker.C:
			m.updateFunctions()
		case <-m.stopCh:
			return
		}
	}
}

// updateFunctions updates the function tool mappings
func (m *PulsarFunctionManager) updateFunctions() {
	// Get all functions
	functions, err := m.getFunctionsList()
	if err != nil {
		log.Printf("Failed to get functions list: %v", err)
		return
	}

	// Track which functions we've seen
	seenFunctions := make(map[string]bool)

	// Add or update functions
	for _, fn := range functions {
		fullName := getFunctionFullName(fn.Tenant, fn.Namespace, fn.Name)
		seenFunctions[fullName] = true

		// Check if we already have this function
		m.mutex.RLock()
		_, exists := m.fnToToolMap[fullName]
		m.mutex.RUnlock()

		changed := false
		if exists {
			// Check if the function has changed
			existingFn, exists := m.fnToToolMap[fullName]
			if exists {
				if !cmp.Equal(*existingFn.Function, *fn) {
					changed = true
				}
				if !existingFn.SchemaFetchSuccess {
					changed = true
				}
			}
			if !changed {
				continue
			}
		}

		// Convert function to tool
		fnTool, err := m.convertFunctionToTool(fn)
		if err != nil || !fnTool.SchemaFetchSuccess {
			if err != nil {
				log.Printf("Failed to convert function %s to tool: %v", fullName, err)
			} else {
				log.Printf("Failed to fetch schema for function %s, retry later...", fullName)
			}
			continue
		}

		if changed {
			m.mcpServer.DeleteTools(fnTool.Tool.Name)
		}
		m.mcpServer.AddTool(fnTool.Tool, m.handleToolCall(fnTool))

		// Add function to map
		m.mutex.Lock()
		m.fnToToolMap[fullName] = fnTool
		m.mutex.Unlock()

		if changed {
			log.Printf("Updated function %s as MCP tool [%s]", fullName, fnTool.Tool.Name)
		} else {
			log.Printf("Added function %s as MCP tool [%s]", fullName, fnTool.Tool.Name)
		}
	}

	// Remove deleted functions
	m.mutex.Lock()
	for fullName, fnTool := range m.fnToToolMap {
		if !seenFunctions[fullName] {
			m.mcpServer.DeleteTools(fnTool.Tool.Name)
			delete(m.fnToToolMap, fullName)
			log.Printf("Removed function %s from MCP tools [%s]", fullName, fnTool.Tool.Name)
		}
	}
	m.mutex.Unlock()
}

// getFunctionsList retrieves all functions from the specified tenants/namespaces
func (m *PulsarFunctionManager) getFunctionsList() ([]*utils.FunctionConfig, error) {
	var allFunctions []*utils.FunctionConfig
	var runningFunctions []*utils.FunctionConfig

	if len(m.tenantNamespaces) == 0 {
		// This is StreamNative supported way to get all functions when using Function Mesh
		functions, err := m.getFunctionsInNamespace("tenant@all", "namespace@all")
		if err != nil {
			return nil, fmt.Errorf("failed to get functions in all namespaces: %w", err)
		}

		allFunctions = append(allFunctions, functions...)
	} else {
		// Get functions from specified tenant/namespaces
		for _, tn := range m.tenantNamespaces {
			parts := strings.Split(tn, "/")
			if len(parts) != 2 {
				log.Printf("Invalid tenant/namespace format: %s", tn)
				continue
			}

			tenant := parts[0]
			namespace := parts[1]

			functions, err := m.getFunctionsInNamespace(tenant, namespace)
			if err != nil {
				log.Printf("Failed to get functions in namespace %s/%s: %v", tenant, namespace, err)
				continue
			}

			allFunctions = append(allFunctions, functions...)
		}
	}

	for _, fn := range allFunctions {
		status, err := m.adminClient.Functions().GetFunctionStatus(fn.Tenant, fn.Namespace, fn.Name)
		if err != nil {
			continue
		}
		if status.NumRunning <= 0 {
			continue
		}
		running := false
		for _, instance := range status.Instances {
			if instance.Status.Err != "" {
				continue
			}
			if instance.Status.Running {
				running = true
				break
			}
		}
		if !running {
			continue
		}
		runningFunctions = append(runningFunctions, fn)
	}

	return runningFunctions, nil
}

// getFunctionsInNamespace retrieves all functions in a namespace
func (m *PulsarFunctionManager) getFunctionsInNamespace(tenant, namespace string) ([]*utils.FunctionConfig, error) {
	var functions []*utils.FunctionConfig

	// Get function names
	functionNames, err := m.adminClient.Functions().GetFunctions(tenant, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get function names: %w", err)
	}

	// Get details for each function
	for _, name := range functionNames {
		parts := strings.Split(name, "/")
		if len(parts) != 3 {
			log.Printf("Invalid function name format: %s", name)
			continue
		}

		function, err := m.adminClient.Functions().GetFunction(parts[0], parts[1], parts[2])
		if err != nil {
			log.Printf("Failed to get function details for %s/%s/%s: %v", parts[0], parts[1], parts[2], err)
			continue
		}

		functions = append(functions, &function)
	}

	return functions, nil
}

// convertFunctionToTool converts a Pulsar Function to an MCP Tool
func (m *PulsarFunctionManager) convertFunctionToTool(fn *utils.FunctionConfig) (*FunctionTool, error) {
	schemaFetchSuccess := true
	// Determine input and output topics
	if len(fn.InputSpecs) == 0 {
		return nil, fmt.Errorf("function has no input topics")
	}

	var inputTopic string
	// Get the first input topic
	for topic := range fn.InputSpecs {
		inputTopic = topic
		break
	}
	if inputTopic == "" {
		return nil, fmt.Errorf("function has no input topics")
	}

	// Get schema for input topic
	inputSchema, err := GetSchemaFromTopic(m.v2adminClient, inputTopic)
	if err != nil {
		log.Printf("Failed to get schema for input topic %s: %v", inputTopic, err)
		// Continue with a default schema
		inputSchema = &SchemaInfo{
			Type: "STRING",
			Definition: map[string]interface{}{
				"type": "string",
			},
		}
		schemaFetchSuccess = false
	}

	// Get output topic and schema
	outputTopic := fn.Output
	var outputSchema *SchemaInfo
	if outputTopic != "" {
		outputSchema, err = GetSchemaFromTopic(m.v2adminClient, outputTopic)
		if err != nil {
			log.Printf("Failed to get schema for output topic %s: %v", outputTopic, err)
			// Continue with a default schema
			outputSchema = &SchemaInfo{
				Type: "STRING",
				Definition: map[string]interface{}{
					"type": "string",
				},
			}
			schemaFetchSuccess = false
		}
	}

	// Generate tool input schema from input topic schema
	toolInputSchema, err := ConvertSchemaToToolInput(inputSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to convert input schema: %w", err)
	}

	toolName := retrieveToolName(fn)
	// Replace non-alphanumeric characters
	toolName = strings.ReplaceAll(toolName, "-", "_")
	toolName = strings.ReplaceAll(toolName, ".", "_")

	// Create description
	description := retrieveToolDescription(fn)

	// Create the tool
	tool := mcp.NewTool(toolName,
		mcp.WithDescription(description),
	)
	tool.InputSchema = *toolInputSchema

	// Create circuit breaker for this function
	circuitBreaker := NewCircuitBreaker(5, 60*time.Second)

	// Store in map
	m.mutex.Lock()
	m.circuitBreakers[toolName] = circuitBreaker
	m.mutex.Unlock()

	return &FunctionTool{
		Name:               toolName,
		Function:           fn,
		InputSchema:        inputSchema,
		OutputSchema:       outputSchema,
		InputTopic:         inputTopic,
		OutputTopic:        outputTopic,
		Tool:               tool,
		SchemaFetchSuccess: schemaFetchSuccess,
	}, nil
}

// handleToolCall returns a handler function for a specific function tool
func (m *PulsarFunctionManager) handleToolCall(fnTool *FunctionTool) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get the circuit breaker
		m.mutex.RLock()
		cb, exists := m.circuitBreakers[fnTool.Name]
		m.mutex.RUnlock()

		if !exists {
			cb = NewCircuitBreaker(5, 60*time.Second)
			m.mutex.Lock()
			m.circuitBreakers[fnTool.Name] = cb
			m.mutex.Unlock()
		}

		// Check if the circuit breaker allows the request
		if !cb.AllowRequest() {
			return mcp.NewToolResultError(fmt.Sprintf("Circuit breaker is open for function %s. Too many failures, please try again later.", fnTool.Name)), nil
		}

		// Create function invoker
		invoker := NewFunctionInvoker(m.pulsarClient)

		// Create context with timeout
		timeoutCtx, cancel := context.WithTimeout(ctx, m.defaultTimeout)
		timeoutCtx = context.WithValue(timeoutCtx, "timeout", m.defaultTimeout)
		defer cancel()

		// Register call
		m.mutex.Lock()
		m.callInProgressMap[fnTool.Name] = cancel
		m.mutex.Unlock()
		defer func() {
			m.mutex.Lock()
			delete(m.callInProgressMap, fnTool.Name)
			m.mutex.Unlock()
		}()

		// Invoke function and wait for result
		result, err := invoker.InvokeFunctionAndWait(timeoutCtx, fnTool, request.Params.Arguments)

		// Record success or failure
		if err != nil {
			cb.RecordFailure()
		} else {
			cb.RecordSuccess()
		}

		return result, err
	}
}

// getFunctionFullName returns the full name of a function
func getFunctionFullName(tenant, namespace, name string) string {
	return fmt.Sprintf("%s/%s/%s", tenant, namespace, name)
}

// retrieveToolName retrieves the tool name from a function
func retrieveToolName(fn *utils.FunctionConfig) string {
	if fn == nil {
		return ""
	}
	fallbackName := fmt.Sprintf("pulsar_function_%s_%s_%s", fn.Tenant, fn.Namespace, fn.Name)
	if fn.CustomRuntimeOptions != "" {
		option := make(map[string]interface{})
		if err := json.Unmarshal([]byte(fn.CustomRuntimeOptions), &option); err != nil {
			return fallbackName
		}
		if envs, ok := option["env"]; ok {
			if envsMap, ok := envs.(map[string]interface{}); ok {
				if name, ok := envsMap[CustomRuntimeOptionsEnvMcpToolNameKey]; ok {
					return name.(string)
				}
			}
		}
	}
	return fallbackName
}

// retrieveToolDescription retrieves the tool description from a function
func retrieveToolDescription(fn *utils.FunctionConfig) string {
	if fn == nil {
		return ""
	}
	fallbackDescription := fmt.Sprintf("Linked to Pulsar Function: %s/%s/%s", fn.Tenant, fn.Namespace, fn.Name)
	if fn.CustomRuntimeOptions != "" {
		option := make(map[string]interface{})
		if err := json.Unmarshal([]byte(fn.CustomRuntimeOptions), &option); err != nil {
			return fallbackDescription
		}
		if envs, ok := option["env"]; ok {
			if envsMap, ok := envs.(map[string]interface{}); ok {
				if description, ok := envsMap[CustomRuntimeOptionsEnvMcpToolDescriptionKey]; ok {
					return description.(string) + " " + fallbackDescription
				}
			}
		}
	}
	return fallbackDescription
}
