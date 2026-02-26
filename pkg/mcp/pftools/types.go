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
	"sync"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
)

// PulsarFunctionManager manages the lifecycle of Pulsar Functions as MCP tools
type PulsarFunctionManager struct {
	adminClient         cmdutils.Client
	v2adminClient       cmdutils.Client
	pulsarClient        pulsar.Client
	fnToToolMap         map[string]*FunctionTool
	failedFunctions     map[string]*functionFailureState
	mutex               sync.RWMutex
	producerCache       map[string]pulsar.Producer
	producerMutex       sync.RWMutex
	pollInterval        time.Duration
	stopCh              chan struct{}
	callInProgressMap   map[string]context.CancelFunc
	mcpServer           *server.MCPServer
	logger              *logrus.Logger
	readOnly            bool
	defaultTimeout      time.Duration
	circuitBreakers     map[string]*CircuitBreaker
	tenantNamespaces    []string
	strictExport        bool
	sessionID           string
	clusterErrorHandler ClusterErrorHandler
}

// FunctionTool represents a Pulsar function exposed as an MCP tool.
type FunctionTool struct {
	Name               string
	Function           *utils.FunctionConfig
	InputSchema        *SchemaInfo
	OutputSchema       *SchemaInfo
	InputTopic         string
	OutputTopic        string
	Tool               mcp.Tool
	SchemaFetchSuccess bool
	SchemaFetchError   error
}

// SchemaInfo represents schema metadata for Pulsar functions.
type SchemaInfo struct {
	Type             string
	Definition       map[string]interface{}
	PulsarSchemaInfo *utils.SchemaInfo
}

type failureCategory string

const (
	failurePermanent failureCategory = "permanent"
	failureRetryable failureCategory = "retryable"
	failureUnknown   failureCategory = "unknown"
)

type functionFailureState struct {
	configHash    string
	category      failureCategory
	lastError     string
	lastLoggedAt  time.Time
	lastAttemptAt time.Time
}

// CircuitBreaker guards function invocations to prevent repeated failures.
type CircuitBreaker struct {
	failureCount     int
	failureThreshold int
	resetTimeout     time.Duration
	lastFailure      time.Time
	state            CircuitState
	mutex            sync.RWMutex
}

// CircuitState represents the circuit breaker state.
type CircuitState int

// Circuit breaker states.
const (
	StateOpen CircuitState = iota
	StateHalfOpen
	StateClosed
)

// ClusterErrorHandler handles cluster errors for Pulsar function managers.
type ClusterErrorHandler func(*PulsarFunctionManager, error)

// ManagerOptions configures PulsarFunctionManager behavior.
type ManagerOptions struct {
	PollInterval        time.Duration
	DefaultTimeout      time.Duration
	FailureThreshold    int
	ResetTimeout        time.Duration
	TenantNamespaces    []string
	StrictExport        bool
	ClusterErrorHandler ClusterErrorHandler
}

// DefaultManagerOptions returns default manager options.
func DefaultManagerOptions() *ManagerOptions {
	return &ManagerOptions{
		PollInterval:     30 * time.Second,
		DefaultTimeout:   10 * time.Second,
		FailureThreshold: 5,
		ResetTimeout:     60 * time.Second,
		TenantNamespaces: []string{},
		StrictExport:     false,
	}
}
