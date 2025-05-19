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
	adminClient       cmdutils.Client               // Pulsar Admin 客户端
	v2adminClient     cmdutils.Client               // Pulsar Admin 客户端
	pulsarClient      pulsar.Client                 // Pulsar 客户端
	fnToToolMap       map[string]*FunctionTool      // 函数到工具的映射
	mutex             sync.RWMutex                  // 用于保护共享数据
	pollInterval      time.Duration                 // 轮询间隔
	stopCh            chan struct{}                 // 用于停止轮询
	callInProgressMap map[string]context.CancelFunc // 正在处理的调用
	mcpServer         *server.MCPServer             // MCP 服务器引用
	readOnly          bool                          // 是否只读模式
	defaultTimeout    time.Duration                 // 默认超时时间
	circuitBreakers   map[string]*CircuitBreaker    // 断路器映射
	tenantNamespaces  []string                      // 需要监听的租户和命名空间列表
}

// FunctionTool 表示一个函数工具
type FunctionTool struct {
	Name         string                // 工具名称
	Function     *utils.FunctionConfig // Pulsar Function 信息
	InputSchema  *SchemaInfo           // 输入 Schema
	OutputSchema *SchemaInfo           // 输出 Schema
	InputTopic   string                // 输入 Topic
	OutputTopic  string                // 输出 Topic
	Tool         mcp.Tool              // MCP 工具定义
}

// SchemaInfo 用于保存 Schema 信息
type SchemaInfo struct {
	Type       string                 // Schema 类型
	Definition map[string]interface{} // Schema 定义
}

// CircuitBreaker 实现断路器模式
type CircuitBreaker struct {
	failureCount     int
	failureThreshold int
	resetTimeout     time.Duration
	lastFailure      time.Time
	state            CircuitState
	mutex            sync.RWMutex
}

// CircuitState 表示断路器状态
type CircuitState int

const (
	// StateOpen 表示断路器打开状态，不允许调用
	StateOpen CircuitState = iota
	// StateHalfOpen 表示断路器半开状态，允许有限调用
	StateHalfOpen
	// StateClosed 表示断路器关闭状态，允许正常调用
	StateClosed
)

// ManagerOptions 配置PulsarFunctionManager的选项
type ManagerOptions struct {
	PollInterval     time.Duration
	DefaultTimeout   time.Duration
	FailureThreshold int
	ResetTimeout     time.Duration
	TenantNamespaces []string
}

// DefaultManagerOptions 返回默认的管理器选项
func DefaultManagerOptions() *ManagerOptions {
	return &ManagerOptions{
		PollInterval:     30 * time.Second,
		DefaultTimeout:   10 * time.Second,
		FailureThreshold: 5,
		ResetTimeout:     60 * time.Second,
		TenantNamespaces: []string{},
	}
}
