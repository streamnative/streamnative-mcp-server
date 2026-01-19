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

package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	context2 "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

// FunctionConnectorList lists supported log components.
var FunctionConnectorList = []string{"sink", "source", "function", "kafka-connect"}

type streamnativeLogsInput struct {
	Component         string `json:"component" jsonschema:"The component to get logs from, including sink, source, function, kafka-connect"`
	Name              string `json:"name" jsonschema:"The name of the resource to get logs from."`
	Tenant            string `json:"tenant,omitempty" jsonschema:"The pulsar tenant of the resource to get logs from. This is required for pulsar functions, pulsar source connectors, pulsar sink connectors. For kafka connect connectors, this is optional, and the default value is 'public'."`
	Namespace         string `json:"namespace,omitempty" jsonschema:"The pulsar namespace of the resource to get logs from. This is required for pulsar functions, pulsar source connectors, pulsar sink connectors. For kafka connect connectors, this is optional, and the default value is 'default'."`
	Size              string `json:"size,omitempty" jsonschema:"String format of the number of lines to get from the logs, e.g. 10, 100, etc. (default: 20)"`
	ReplicaID         int    `json:"replica_id,omitempty" jsonschema:"The replica index of the resource to get logs from, this is used for multiple replicas of running pulsar functions, pulsar source connectors, pulsar sink connectors, and kafka connect connectors. The value should be a positive integer (like 0, 1, 2, etc.), and the default value is -1, which means all replicas."`
	Timestamp         string `json:"timestamp,omitempty" jsonschema:"Start timestamp of logs, for example: 1662430984225"`
	Since             string `json:"since,omitempty" jsonschema:"Since time of logs, numbers end with s|m|h, for example one hour ago: 1h"`
	PreviousContainer bool   `json:"previous_container,omitempty" jsonschema:"Return previous terminated container logs, defaults to false."`
}

// StreamNativeAddLogTools adds log tools
func StreamNativeAddLogTools(s *sdk.Server, _ bool, features []string) {
	if !slices.Contains(features, string(FeatureStreamNativeCloud)) && !slices.Contains(features, string(FeatureAll)) {
		return
	}

	inputSchema, err := InputSchema[streamnativeLogsInput]()
	if err != nil {
		return
	}
	setSchemaPropertyEnum(inputSchema, "component", FunctionConnectorList)
	setSchemaPropertyDefault(inputSchema, "tenant", "public")
	setSchemaPropertyDefault(inputSchema, "namespace", "default")
	setSchemaPropertyDefault(inputSchema, "size", "20")
	setSchemaPropertyDefault(inputSchema, "replica_id", -1)
	setSchemaPropertyDefault(inputSchema, "previous_container", false)

	logTool := &sdk.Tool{
		Name:        "sncloud_logs",
		Description: "Display the logs of resources in StreamNative Cloud, including pulsar functions, pulsar source connectors, pulsar sink connectors, and kafka connect connectors logs running along with PulsarInstance and PulsarCluster. This tool is used to help you debug the issues of resources in StreamNative Cloud. You can use `sncloud_context_use_cluster` to change the context to a specific cluster first, then use this tool to get the logs of resources in the cluster. This tool is suggested to be used with 'pulsar_admin_functions', 'pulsar_admin_sinks', 'pulsar_admin_sources', and 'kafka_admin_connect'",
		InputSchema: inputSchema,
	}
	sdk.AddTool(s, logTool, handleStreamNativeResourcesLog)
}

// LogOptions captures parameters for StreamNative log queries.
type LogOptions struct {
	ServiceURL                   string
	Organization                 string
	Instance                     string
	Cluster                      string
	Component                    string
	Name                         string
	PulsarTenant                 string
	PulsarNamespace              string
	Size                         string
	Since                        string
	Timestamp                    string
	Follow                       bool
	replicaID                    int
	Previous                     bool
	InsecureSkipTLSVerifyBackend bool
}

// LogResult represents a log query response.
type LogResult struct {
	Total int          `json:"total"`
	Data  []LogContent `json:"data"`
}

// LogContent represents a single log entry.
type LogContent struct {
	Message  string `json:"message"`
	Position int64  `json:"position"`
	Record   int64  `json:"record"`
}

func handleStreamNativeResourcesLog(ctx context.Context, _ *sdk.CallToolRequest, input streamnativeLogsInput) (*sdk.CallToolResult, any, error) {
	// Get log client from session
	session := context2.GetSNCloudSession(ctx)
	if session == nil {
		return nil, nil, fmt.Errorf("failed to get StreamNative Cloud session")
	}
	instance, cluster, organization := session.Ctx.PulsarInstance, session.Ctx.PulsarCluster, session.Ctx.Organization
	if instance == "" || cluster == "" || organization == "" {
		return nil, nil, fmt.Errorf("no context is set, please use `sncloud_context_use_cluster` to set the context first")
	}

	// Extract required parameters with validation
	component := input.Component
	if component == "" {
		return nil, nil, fmt.Errorf("missing required parameter 'component'")
	}

	name := input.Name
	if name == "" {
		return nil, nil, fmt.Errorf("missing required parameter 'name'")
	}

	tenant := input.Tenant
	if tenant == "" {
		tenant = "public"
	}

	namespace := input.Namespace
	if namespace == "" {
		namespace = "default"
	}

	size := input.Size
	if size == "" {
		size = "20"
	}

	replicaID := input.ReplicaID
	if replicaID == 0 {
		replicaID = -1
	}

	timestampStr := input.Timestamp
	sinceStr := input.Since

	previousContainer := input.PreviousContainer

	if sinceStr != "" {
		sinceStr = "-" + sinceStr
	}
	t := time.Now()
	r1 := regexp.MustCompile(`^-(\d+)(s|m|h)$`)
	if timestampStr == "" {
		if r1.MatchString(sinceStr) {
			ago, _ := time.ParseDuration(sinceStr)
			t = t.Add(ago)
			timestampStr = strconv.FormatInt(t.UnixNano()/1e6, 10)
		}
	}

	logOptions := &LogOptions{
		ServiceURL:                   session.Ctx.LogAPIURL,
		Organization:                 organization,
		Instance:                     instance,
		Cluster:                      cluster,
		Component:                    component,
		Name:                         name,
		PulsarTenant:                 tenant,
		PulsarNamespace:              namespace,
		Size:                         size,
		Follow:                       false, // we do not support follow as streaming in MCP yet.
		replicaID:                    replicaID,
		Previous:                     previousContainer,
		InsecureSkipTLSVerifyBackend: false,
		Since:                        sinceStr,
		Timestamp:                    timestampStr,
	}

	client, err := session.GetLogClient()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get logging client: %v", err)
	}

	results := []string{}
	results, err = logOptions.getLogs(client, 0, 0, results)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get logs: %v", err)
	}

	if len(results) == 0 {
		return newToolResultText("No logs found"), nil, nil
	}

	return newToolResultText(strings.Join(results, "\n")), nil, nil
}

// Helper function to send HTTP GET request to fetch logs
func (opts *LogOptions) sendLogRequest(client *http.Client, offset int, _ int) (*http.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, opts.ServiceURL, nil)
	if err != nil {
		return nil, err
	}

	// Set query parameters
	query := req.URL.Query()
	query.Add("organization", opts.Organization)
	query.Add("instance", opts.Instance)
	query.Add("cluster", opts.Cluster)
	query.Add("component", opts.Component)
	query.Add("name", opts.Name)
	query.Add("namespace", opts.PulsarNamespace)
	query.Add("tenant", opts.PulsarTenant)
	query.Add("size", opts.Size)
	query.Add("offset", strconv.Itoa(offset))
	query.Add("follow", "false")
	if opts.replicaID != -1 {
		query.Add("replica_id", strconv.Itoa(opts.replicaID))
	}
	if opts.Since != "" {
		query.Add("since", opts.Since)
	}
	if opts.Timestamp != "" {
		query.Add("timestamp", opts.Timestamp)
	}
	if opts.Previous {
		query.Add("previous_container", "true")
	}
	req.URL.RawQuery = query.Encode()

	// Send request
	return client.Do(req)
}

// getLogs fetches logs from StreamNative log API.
func (opts *LogOptions) getLogs(client *http.Client, offset int, lastRecord int, results []string) ([]string, error) {
	response, err := opts.sendLogRequest(client, offset, lastRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var logResult LogResult
	if err := json.Unmarshal(body, &logResult); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	// Process logs
	for _, log := range logResult.Data {
		results = append(results, log.Message)
		offset++
		lastRecord = int(log.Record)
	}

	if offset < logResult.Total {
		return opts.getLogs(client, offset, lastRecord, results)
	}

	return results, nil
}
