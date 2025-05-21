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
	"sync"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
)

// PulsarFunctionManager manages the lifecycle of Pulsar Functions as MCP tools
type PulsarFunctionManager struct {
	adminClient       cmdutils.Client
	v2adminClient     cmdutils.Client
	pulsarClient      pulsar.Client
	fnToToolMap       map[string]*FunctionTool
	mutex             sync.RWMutex
	pollInterval      time.Duration
	stopCh            chan struct{}
	callInProgressMap map[string]context.CancelFunc
	mcpServer         *server.MCPServer
	readOnly          bool
	defaultTimeout    time.Duration
	circuitBreakers   map[string]*CircuitBreaker
	tenantNamespaces  []string
	strictExport      bool
}

type FunctionTool struct {
	Name               string
	Function           *utils.FunctionConfig
	InputSchema        *SchemaInfo
	OutputSchema       *SchemaInfo
	InputTopic         string
	OutputTopic        string
	Tool               mcp.Tool
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

type ManagerOptions struct {
	PollInterval     time.Duration
	DefaultTimeout   time.Duration
	FailureThreshold int
	ResetTimeout     time.Duration
	TenantNamespaces []string
	StrictExport     bool
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
