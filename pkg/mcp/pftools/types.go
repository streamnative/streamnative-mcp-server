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
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCPServerInterface defines the interface for an MCP server.
// This is used by pftools to avoid circular dependency on the mcp package.
type MCPServerInterface interface {
	AddTool(tool interface{}, handler interface{}) error
	AddSessionTool(sessionID string, tool interface{}, handler interface{}) error
	DeleteTools(name string) error
	DeleteSessionTools(sessionID, name string) error
}

// PulsarFunctionManager manages the lifecycle of Pulsar Functions as MCP tools
type PulsarFunctionManager struct {
	adminClient         cmdutils.Client
	v2adminClient       cmdutils.Client
	pulsarClient        pulsar.Client
	fnToToolMap         map[string]*FunctionTool
	mutex               sync.RWMutex
	producerCache       map[string]pulsar.Producer
	producerMutex       sync.RWMutex
	pollInterval        time.Duration
	stopCh              chan struct{}
	callInProgressMap   map[string]context.CancelFunc
	mcpServer           MCPServerInterface
	readOnly            bool
	defaultTimeout      time.Duration
	circuitBreakers     map[string]*CircuitBreaker
	tenantNamespaces    []string
	strictExport        bool
	sessionID           string
	clusterErrorHandler ClusterErrorHandler
}

type FunctionTool struct {
	Name               string
	Function           *utils.FunctionConfig
	InputSchema        *SchemaInfo
	OutputSchema       *SchemaInfo
	InputTopic         string
	OutputTopic        string
	Tool               mcpsdk.Tool
	SchemaFetchSuccess bool
}

type SchemaInfo struct {
	Type             string
	Definition       map[string]interface{}
	PulsarSchemaInfo *utils.SchemaInfo
}

type CircuitBreaker struct {
	failureCount     int
	failureThreshold int
	resetTimeout     time.Duration
	lastFailure      time.Time
	state            CircuitState
	mutex            sync.RWMutex
}

type CircuitState int

const (
	StateOpen CircuitState = iota
	StateHalfOpen
	StateClosed
)

type ClusterErrorHandler func(*PulsarFunctionManager, error)

type ManagerOptions struct {
	PollInterval        time.Duration
	DefaultTimeout      time.Duration
	FailureThreshold    int
	ResetTimeout        time.Duration
	TenantNamespaces    []string
	StrictExport        bool
	ClusterErrorHandler ClusterErrorHandler
}

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

// ToolInputSchema represents the input schema for a tool.
// This is a local type used within pftools to avoid external dependencies.
type ToolInputSchema struct {
	Type       string
	Properties map[string]interface{}
}
