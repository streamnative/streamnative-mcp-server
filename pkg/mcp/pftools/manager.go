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

package pftools

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	pulsarclient "github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/rest"
	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
	"github.com/streamnative/streamnative-mcp-server/pkg/schema"
)

const (
	// CustomRuntimeOptionsEnvMcpToolNameKey is the env var name for tool names.
	CustomRuntimeOptionsEnvMcpToolNameKey = "MCP_TOOL_NAME"
	// CustomRuntimeOptionsEnvMcpToolDescriptionKey is the env var name for tool descriptions.
	CustomRuntimeOptionsEnvMcpToolDescriptionKey = "MCP_TOOL_DESCRIPTION"
)

// DefaultStringSchemaInfo defines the default schema info for string payloads.
var DefaultStringSchemaInfo = &SchemaInfo{
	Type: "STRING",
	Definition: map[string]interface{}{
		"type": "string",
	},
	PulsarSchemaInfo: &utils.SchemaInfo{
		Type: "STRING",
	},
}

// Server is imported directly to avoid circular dependency
type Server struct {
	MCPServer     *server.MCPServer
	KafkaSession  *kafka.Session
	PulsarSession *pulsar.Session
	Logger        *logrus.Logger
}

// NewPulsarFunctionManager creates a new PulsarFunctionManager
func NewPulsarFunctionManager(snServer *Server, readOnly bool, options *ManagerOptions, sessionID string) (*PulsarFunctionManager, error) {
	// Get Pulsar client and admin client
	if snServer.PulsarSession == nil {
		return nil, fmt.Errorf("pulsar session not found in context")
	}

	// Get Pulsar client from session using type-safe interface
	pulsarClient, err := snServer.PulsarSession.GetPulsarClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get Pulsar client: %w", err)
	}

	adminClient, err := snServer.PulsarSession.GetAdminV3Client()
	if err != nil {
		return nil, fmt.Errorf("failed to get Pulsar admin v3 client: %w", err)
	}

	v2adminClient, err := snServer.PulsarSession.GetAdminClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get Pulsar admin client: %w", err)
	}

	if options == nil {
		options = DefaultManagerOptions()
	}

	logger := snServer.Logger
	if logger == nil {
		logger = logrus.New()
	}

	// Create the manager
	manager := &PulsarFunctionManager{
		adminClient:         adminClient,
		v2adminClient:       v2adminClient,
		pulsarClient:        pulsarClient,
		fnToToolMap:         make(map[string]*FunctionTool),
		failedFunctions:     make(map[string]*functionFailureState),
		mutex:               sync.RWMutex{},
		producerCache:       make(map[string]pulsarclient.Producer),
		producerMutex:       sync.RWMutex{},
		pollInterval:        options.PollInterval,
		stopCh:              make(chan struct{}),
		callInProgressMap:   make(map[string]context.CancelFunc),
		mcpServer:           snServer.MCPServer,
		logger:              logger,
		readOnly:            readOnly,
		defaultTimeout:      options.DefaultTimeout,
		circuitBreakers:     make(map[string]*CircuitBreaker),
		tenantNamespaces:    options.TenantNamespaces,
		strictExport:        options.StrictExport,
		sessionID:           sessionID,
		clusterErrorHandler: options.ClusterErrorHandler,
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

	m.producerMutex.Lock()
	defer m.producerMutex.Unlock()
	for topic, producer := range m.producerCache {
		m.logger.WithField("topic", topic).Info("Closing producer for topic")
		producer.Close()
	}
	m.producerCache = make(map[string]pulsarclient.Producer)
	m.logger.Info("All cached producers closed and cache cleared")
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
		m.logger.WithError(err).Warn("Failed to get functions list")

		// Check if this is a cluster health error and invoke callback if configured
		if (IsClusterUnhealthy(err) || IsAuthError(err) || IsNotFoundError(err)) && m.clusterErrorHandler != nil {
			go m.clusterErrorHandler(m, err)
		}
		return
	}

	// Track which functions we've seen
	seenFunctions := make(map[string]bool)

	// Add or update functions
	for _, fn := range functions {
		fullName := getFunctionFullName(fn.Tenant, fn.Namespace, fn.Name)
		seenFunctions[fullName] = true

		configHash, hashErr := computeFunctionConfigHash(fn)
		if hashErr != nil {
			m.logger.WithError(hashErr).WithField("function", fullName).Warn("Failed to compute config hash")
		}

		// Check if we already have this function
		m.mutex.RLock()
		existingFn, exists := m.fnToToolMap[fullName]
		failureState, hasFailure := m.failedFunctions[fullName]
		m.mutex.RUnlock()

		if hasFailure && configHash != "" && failureState.configHash != configHash {
			m.mutex.Lock()
			delete(m.failedFunctions, fullName)
			m.mutex.Unlock()
			hasFailure = false
			failureState = nil
		}

		changed := false
		if exists {
			// Check if the function has changed
			if !cmp.Equal(*existingFn.Function, *fn) {
				changed = true
			}
			if !existingFn.SchemaFetchSuccess {
				changed = true
			}
			if !changed {
				continue
			}
		}

		if hasFailure && configHash != "" && failureState.configHash == configHash {
			if shouldSkipFailure(failureState, m.pollInterval, time.Now()) {
				continue
			}
		}

		// Convert function to tool
		attemptAt := time.Now()
		fnTool, err := m.convertFunctionToTool(fn)
		if err != nil || (fnTool != nil && !fnTool.SchemaFetchSuccess) {
			failureErr := err
			if failureErr == nil && fnTool != nil && fnTool.SchemaFetchError != nil {
				failureErr = fnTool.SchemaFetchError
			}
			if failureErr == nil {
				failureErr = errors.New("schema fetch failed")
			}

			category := classifyConvertError(failureErr)
			errorMsg := failureErr.Error()
			logNow := shouldLogFailure(failureState, configHash, category, errorMsg)

			if configHash != "" {
				newState := &functionFailureState{
					configHash:    configHash,
					category:      category,
					lastError:     errorMsg,
					lastAttemptAt: attemptAt,
				}
				if logNow {
					newState.lastLoggedAt = time.Now()
				} else if failureState != nil {
					newState.lastLoggedAt = failureState.lastLoggedAt
				}
				m.mutex.Lock()
				m.failedFunctions[fullName] = newState
				m.mutex.Unlock()
			}
			if logNow {
				if err != nil {
					m.logger.WithError(failureErr).WithFields(logrus.Fields{
						"function": fullName,
						"category": category,
					}).Warn("Failed to convert function to tool")
				} else {
					m.logger.WithError(failureErr).WithFields(logrus.Fields{
						"function": fullName,
						"category": category,
					}).Warn("Failed to fetch schema for function, retry later")
				}
			}
			continue
		}

		if hasFailure {
			m.mutex.Lock()
			delete(m.failedFunctions, fullName)
			m.mutex.Unlock()
		}

		if changed {
			if m.sessionID != "" {
				err := m.mcpServer.DeleteSessionTools(m.sessionID, fnTool.Tool.Name)
				if err != nil {
					m.logger.WithError(err).WithFields(logrus.Fields{
						"tool":       fnTool.Tool.Name,
						"session_id": m.sessionID,
					}).Warn("Failed to delete tool from session")
				}
			} else {
				m.mcpServer.DeleteTools(fnTool.Tool.Name)
			}
		}
		if m.sessionID != "" {
			err := m.mcpServer.AddSessionTool(m.sessionID, fnTool.Tool, m.handleToolCall(fnTool))
			if err != nil {
				m.logger.WithError(err).WithFields(logrus.Fields{
					"tool":       fnTool.Tool.Name,
					"session_id": m.sessionID,
				}).Warn("Failed to add tool to session")
			}
		} else {
			m.mcpServer.AddTool(fnTool.Tool, m.handleToolCall(fnTool))
		}

		// Add function to map
		m.mutex.Lock()
		m.fnToToolMap[fullName] = fnTool
		m.mutex.Unlock()

		if changed {
			m.logger.WithFields(logrus.Fields{
				"function": fullName,
				"tool":     fnTool.Tool.Name,
			}).Info("Updated function as MCP tool")
		} else {
			m.logger.WithFields(logrus.Fields{
				"function": fullName,
				"tool":     fnTool.Tool.Name,
			}).Info("Added function as MCP tool")
		}
	}

	// Remove deleted functions
	m.mutex.Lock()
	for fullName, fnTool := range m.fnToToolMap {
		if !seenFunctions[fullName] {
			if m.sessionID != "" {
				err := m.mcpServer.DeleteSessionTools(m.sessionID, fnTool.Tool.Name)
				if err != nil {
					m.logger.WithError(err).WithFields(logrus.Fields{
						"tool":       fnTool.Tool.Name,
						"session_id": m.sessionID,
					}).Warn("Failed to delete tool from session")
				}
			} else {
				m.mcpServer.DeleteTools(fnTool.Tool.Name)
			}
			delete(m.fnToToolMap, fullName)
			delete(m.failedFunctions, fullName)
			m.logger.WithFields(logrus.Fields{
				"function": fullName,
				"tool":     fnTool.Tool.Name,
			}).Info("Removed function from MCP tools")
		}
	}
	m.mutex.Unlock()
}

func computeFunctionConfigHash(fn *utils.FunctionConfig) (string, error) {
	if fn == nil {
		return "", errors.New("function config is nil")
	}
	data, err := json.Marshal(fn)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return fmt.Sprintf("%x", sum[:]), nil
}

func shouldSkipFailure(state *functionFailureState, pollInterval time.Duration, now time.Time) bool {
	if state == nil {
		return false
	}
	switch state.category {
	case failurePermanent:
		return true
	case failureRetryable:
		if state.lastAttemptAt.IsZero() {
			return false
		}
		return now.Sub(state.lastAttemptAt) < pollInterval
	default:
		return true
	}
}

func shouldLogFailure(prev *functionFailureState, configHash string, category failureCategory, errMsg string) bool {
	if prev == nil {
		return true
	}
	if configHash == "" {
		return true
	}
	if prev.configHash != configHash {
		return true
	}
	if prev.category != category {
		return true
	}
	if prev.lastError != errMsg {
		return true
	}
	return false
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
				m.logger.WithField("tenant_namespace", tn).Warn("Invalid tenant/namespace format")
				continue
			}

			tenant := parts[0]
			namespace := parts[1]

			functions, err := m.getFunctionsInNamespace(tenant, namespace)
			if err != nil {
				m.logger.WithError(err).WithFields(logrus.Fields{
					"tenant":    tenant,
					"namespace": namespace,
				}).Warn("Failed to get functions in namespace")
				continue
			}

			allFunctions = append(allFunctions, functions...)
		}
	}

	for _, fn := range allFunctions {
		if m.strictExport &&
			!strings.Contains(fn.CustomRuntimeOptions, CustomRuntimeOptionsEnvMcpToolNameKey) &&
			!strings.Contains(fn.CustomRuntimeOptions, CustomRuntimeOptionsEnvMcpToolDescriptionKey) {
			continue
		}
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
			m.logger.WithField("function_name", name).Warn("Invalid function name format")
			continue
		}

		function, err := m.adminClient.Functions().GetFunction(parts[0], parts[1], parts[2])
		if err != nil {
			m.logger.WithError(err).WithFields(logrus.Fields{
				"tenant":    parts[0],
				"namespace": parts[1],
				"function":  parts[2],
			}).Warn("Failed to get function details")
			continue
		}

		functions = append(functions, &function)
	}

	return functions, nil
}

// convertFunctionToTool converts a Pulsar Function to an MCP Tool
func (m *PulsarFunctionManager) convertFunctionToTool(fn *utils.FunctionConfig) (*FunctionTool, error) {
	schemaFetchSuccess := true
	var schemaFetchErr error
	// Determine input and output topics
	if len(fn.InputSpecs) == 0 {
		return nil, ErrFunctionNoInputTopics
	}

	var inputTopic string
	// Get the first input topic
	for topic := range fn.InputSpecs {
		inputTopic = topic
		break
	}
	if inputTopic == "" {
		return nil, ErrFunctionNoInputTopics
	}

	// Get schema for input topic
	inputSchema, err := GetSchemaFromTopic(m.v2adminClient, inputTopic)
	if err != nil {
		// Continue with a default schema
		inputSchema = DefaultStringSchemaInfo
		if restError, ok := err.(rest.Error); ok {
			if restError.Code != 404 {
				m.logger.WithError(err).WithFields(logrus.Fields{
					"topic":     inputTopic,
					"direction": "input",
				}).Warn("Failed to get schema for topic")
				schemaFetchSuccess = false
				schemaFetchErr = errors.Join(schemaFetchErr, err)
			}
		} else {
			m.logger.WithError(err).WithFields(logrus.Fields{
				"topic":     inputTopic,
				"direction": "input",
			}).Warn("Failed to get schema for topic")
			schemaFetchSuccess = false
			schemaFetchErr = errors.Join(schemaFetchErr, err)
		}
	}

	// Get output topic and schema
	outputTopic := fn.Output
	var outputSchema *SchemaInfo
	if outputTopic != "" {
		outputSchema, err = GetSchemaFromTopic(m.v2adminClient, outputTopic)
		if err != nil {
			// Continue with a default schema
			outputSchema = DefaultStringSchemaInfo
			if restError, ok := err.(rest.Error); ok {
				if restError.Code != 404 {
					m.logger.WithError(err).WithFields(logrus.Fields{
						"topic":     outputTopic,
						"direction": "output",
					}).Warn("Failed to get schema for topic")
					schemaFetchSuccess = false
					schemaFetchErr = errors.Join(schemaFetchErr, err)
				}
			} else {
				m.logger.WithError(err).WithFields(logrus.Fields{
					"topic":     outputTopic,
					"direction": "output",
				}).Warn("Failed to get schema for topic")
				schemaFetchSuccess = false
				schemaFetchErr = errors.Join(schemaFetchErr, err)
			}
		}
	}

	toolName := retrieveToolName(fn)
	// Replace non-alphanumeric characters
	toolName = strings.ReplaceAll(toolName, "-", "_")
	toolName = strings.ReplaceAll(toolName, ".", "_")

	// Create description
	description := retrieveToolDescription(fn)

	schemaConverter, err := schema.ConverterFactory(inputSchema.Type)
	if err != nil {
		return nil, errors.Join(ErrSchemaConversionFailed, err)
	}

	toolInputSchemaProperties, err := schemaConverter.ToMCPToolInputSchemaProperties(inputSchema.PulsarSchemaInfo)
	if err != nil {
		return nil, errors.Join(ErrSchemaConversionFailed, err)
	}

	toolInputSchemaProperties = append(toolInputSchemaProperties, mcp.WithDescription(description))

	// Create the tool
	tool := mcp.NewTool(toolName,
		toolInputSchemaProperties...,
	)

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
		SchemaFetchError:   schemaFetchErr,
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
		invoker := NewFunctionInvoker(m)

		// Create context with timeout
		timeoutCtx, cancel := context.WithTimeout(ctx, m.defaultTimeout)
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
		result, err := invoker.InvokeFunctionAndWait(timeoutCtx, fnTool, request.GetArguments())

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

// GetProducer retrieves a producer from the cache or creates a new one if not found.
func (m *PulsarFunctionManager) GetProducer(topic string) (pulsarclient.Producer, error) {
	m.producerMutex.RLock()
	producer, found := m.producerCache[topic]
	m.producerMutex.RUnlock()

	if found {
		return producer, nil
	}

	m.producerMutex.Lock()
	defer m.producerMutex.Unlock()

	producer, found = m.producerCache[topic]
	if found {
		return producer, nil
	}

	newProducer, err := m.pulsarClient.CreateProducer(pulsarclient.ProducerOptions{
		Topic: topic,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create producer for topic %s: %w", topic, err)
	}

	m.producerCache[topic] = newProducer
	m.logger.WithField("topic", topic).Info("Created and cached producer for topic")
	return newProducer, nil
}
