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
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/pftools"
)

var (
	functionManagers     = make(map[string]*pftools.PulsarFunctionManager)
	functionManagersLock sync.RWMutex
)

func StopAllPulsarFunctionManagers() {
	functionManagersLock.Lock()
	defer functionManagersLock.Unlock()

	for id, manager := range functionManagers {
		log.Printf("Stopping Pulsar Function manager: %s", id)
		manager.Stop()
		delete(functionManagers, id)
	}

	if len(functionManagers) > 0 {
		time.Sleep(500 * time.Millisecond)
	}

	log.Println("All Pulsar Function managers stopped")
}

func PulsarFunctionManagedMcpTools(s *server.MCPServer, readOnly bool, features []string) {
	if !slices.Contains(features, string(FeatureAll)) &&
		!slices.Contains(features, string(FeatureFunctionsAsTools)) &&
		!slices.Contains(features, string(FeatureStreamNativeCloud)) {
		return
	}

	options := pftools.DefaultManagerOptions()

	if pollIntervalStr := os.Getenv("FUNCTIONS_AS_TOOLS_POLL_INTERVAL"); pollIntervalStr != "" {
		if seconds, err := strconv.Atoi(pollIntervalStr); err == nil && seconds > 0 {
			options.PollInterval = time.Duration(seconds) * time.Second
			log.Printf("Setting Pulsar Functions poll interval to %v", options.PollInterval)
		}
	}

	if timeoutStr := os.Getenv("FUNCTIONS_AS_TOOLS_TIMEOUT"); timeoutStr != "" {
		if seconds, err := strconv.Atoi(timeoutStr); err == nil && seconds > 0 {
			options.DefaultTimeout = time.Duration(seconds) * time.Second
			log.Printf("Setting Pulsar Functions default timeout to %v", options.DefaultTimeout)
		}
	}

	if failureThresholdStr := os.Getenv("FUNCTIONS_AS_TOOLS_FAILURE_THRESHOLD"); failureThresholdStr != "" {
		if threshold, err := strconv.Atoi(failureThresholdStr); err == nil && threshold > 0 {
			options.FailureThreshold = threshold
			log.Printf("Setting Pulsar Functions failure threshold to %d", options.FailureThreshold)
		}
	}

	if resetTimeoutStr := os.Getenv("FUNCTIONS_AS_TOOLS_RESET_TIMEOUT"); resetTimeoutStr != "" {
		if seconds, err := strconv.Atoi(resetTimeoutStr); err == nil && seconds > 0 {
			options.ResetTimeout = time.Duration(seconds) * time.Second
			log.Printf("Setting Pulsar Functions reset timeout to %v", options.ResetTimeout)
		}
	}

	if tenantNamespacesStr := os.Getenv("FUNCTIONS_AS_TOOLS_TENANT_NAMESPACES"); tenantNamespacesStr != "" {
		options.TenantNamespaces = strings.Split(tenantNamespacesStr, ",")
		log.Printf("Setting Pulsar Functions tenant namespaces to %v", options.TenantNamespaces)
	}

	if strictExportStr := os.Getenv("FUNCTIONS_AS_TOOLS_STRICT_EXPORT"); strictExportStr != "" {
		options.StrictExport = strictExportStr == "true"
		log.Printf("Setting Pulsar Functions strict export to %v", options.StrictExport)
	}

	manager, err := pftools.NewPulsarFunctionManager(s, readOnly, options)
	if err != nil {
		log.Printf("Failed to create Pulsar Function manager: %v", err)
		return
	}

	manager.Start()

	managerID := "FUNCTIONS_AS_TOOLS_manager_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	functionManagersLock.Lock()
	functionManagers[managerID] = manager
	functionManagersLock.Unlock()

	log.Printf("Registered Pulsar Function Manager with ID: %s", managerID)
	log.Printf("Pulsar Functions as MCP Tools service started. Functions will be dynamically converted to MCP tools.")
}
