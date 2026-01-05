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

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/adapter"
	context2 "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
)

var FunctionConnectorList = []string{"sink", "source", "function", "kafka-connect"}

// StreamNativeAddLogTools adds log tools
func StreamNativeAddLogTools(s *MCPServer, _ bool, features []string) {
	if !slices.Contains(features, string(FeatureStreamNativeCloud)) && !slices.Contains(features, string(FeatureAll)) {
		return
	}

	logTool := builders.NewTool("sncloud_logs",
		builders.WithDescription("Display the logs of resources in StreamNative Cloud, including pulsar functions, pulsar source connectors, pulsar sink connectors, and kafka connect connectors logs running along with PulsarInstance and PulsarCluster."+
			"This tool is used to help you debug the issues of resources in StreamNative Cloud. You can use `sncloud_context_use_cluster` to change the context to a specific cluster first, then use this tool to get the logs of resources in the cluster. This tool is suggested to be used with 'pulsar_admin_functions', 'pulsar_admin_sinks', 'pulsar_admin_sources', and 'kafka_admin_connect'"),
		builders.WithString("component", builders.Required(),
			builders.Description("The component to get logs from, including "+strings.Join(FunctionConnectorList, ", ")),
			builders.Enum(FunctionConnectorList...),
		),
		builders.WithString("name", builders.Required(),
			builders.Description("The name of the resource to get logs from."),
		),
		builders.WithString("tenant", builders.Required(),
			builders.Description("The pulsar tenant of the resource to get logs from. This is required for pulsar functions, pulsar source connectors, pulsar sink connectors. For kafka connect connectors, this is optional, and the default value is 'public'."),
			builders.DefaultString("public"),
		),
		builders.WithString("namespace", builders.Required(),
			builders.Description("The pulsar namespace of the resource to get logs from. This is required for pulsar functions, pulsar source connectors, pulsar sink connectors. For kafka connect connectors, this is optional, and the default value is 'default'."),
			builders.DefaultString("default"),
		),
		builders.WithString("size",
			builders.Description("String format of the number of lines to get from the logs, e.g. 10, 100, etc. (default: 20)"),
			builders.DefaultString("20"),
		),
		builders.WithNumber("replica_id",
			builders.Description("The replica index of the resource to get logs from, this is used for multiple replicas of running pulsar functions, pulsar source connectors, pulsar sink connectors, and kafka connect connectors. The value should be a positive integer (like 0, 1, 2, etc.), and the default value is -1, which means all replicas."),
			builders.DefaultNumber(-1),
		),
		builders.WithString("timestamp",
			builders.Description("Start timestamp of logs, for example: 1662430984225"),
		),
		builders.WithString("since",
			builders.Description("Since time of logs, numbers end with s|m|h, for example one hour ago: 1h"),
		),
		builders.WithBoolean("previous_container",
			builders.Description("Return previous terminated container logs, defaults to false."),
			builders.DefaultBool(false),
		),
	)
	s.AddTool(logTool, handleStreamNativeResourcesLog)
}

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

type LogResult struct {
	Total int          `json:"total"`
	Data  []LogContent `json:"data"`
}

type LogContent struct {
	Message  string `json:"message"`
	Position int64  `json:"position"`
	Record   int64  `json:"record"`
}

func handleStreamNativeResourcesLog(ctx context.Context, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	// Get log client from session
	session := context2.GetSNCloudSession(ctx)
	if session == nil {
		return nil, fmt.Errorf("failed to get StreamNative Cloud session")
	}
	instance, cluster, organization := session.Ctx.PulsarInstance, session.Ctx.PulsarCluster, session.Ctx.Organization
	if instance == "" || cluster == "" || organization == "" {
		return adapter.NewErrorResult("No context is set, please use `sncloud_context_use_cluster` to set the context first."), nil
	}

	// Extract required parameters with validation
	component, err := adapter.RequireString(request, "component")
	if err != nil {
		return adapter.NewErrorResult("Failed to get component: %v", err), nil
	}

	name, err := adapter.RequireString(request, "name")
	if err != nil {
		return adapter.NewErrorResult("Failed to get name: %v", err), nil
	}

	tenant := adapter.GetString(request, "tenant", "public")

	namespace := adapter.GetString(request, "namespace", "default")

	size := adapter.GetString(request, "size", "20")

	replicaID := adapter.GetInt(request, "replica_id", -1)
	if replicaID == 0 {
		replicaID = -1
	}

	timestampStr := adapter.GetString(request, "timestamp", "")
	sinceStr := adapter.GetString(request, "since", "")

	previousContainer := adapter.GetBool(request, "previous_container", false)

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
		return adapter.NewErrorResult("Failed to get logging client: %v", err), nil
	}

	results := []string{}
	results, err = logOptions.getLogs(client, 0, 0, results)
	if err != nil {
		return adapter.NewErrorResult("Failed to get logs: %v", err), nil
	}

	if len(results) == 0 {
		return adapter.NewTextResult("No logs found"), nil
	}

	return adapter.NewTextResult(strings.Join(results, "\n")), nil
}

func (o *LogOptions) getLogs(client *http.Client, position int64,
	record int64, results []string) ([]string, error) {
	var err error
	logBrowseMode := "next"
	if o.Previous {
		logBrowseMode = "previous"
	}
	url := fmt.Sprintf("%s/v1alpha1/logs?"+
		"organization=%s"+
		"&instance=%s"+
		"&cluster=%s"+
		"&component=%s"+
		"&name=%s"+
		"&pulsar_tenant=%s"+
		"&pulsar_namespace=%s"+
		"&size=%s"+
		"&log_browse_mode=%s"+
		"&timestamp=%s"+
		"&next_position=%d"+
		"&next_record=%d"+
		"&replica_id=%d",
		o.ServiceURL,
		o.Organization,
		o.Instance,
		o.Cluster,
		o.Component,
		o.Name,
		o.PulsarTenant,
		o.PulsarNamespace,
		o.Size,
		logBrowseMode,
		o.Timestamp,
		position,
		record,
		o.replicaID,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return results, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return results, fmt.Errorf("failed to request logs (%s): %v", url, err)
	}
	defer resp.Body.Close()

	var logResult LogResult
	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return results, fmt.Errorf("failed to read response body: %v", err)
	}
	err = json.Unmarshal(body, &logResult)
	if err != nil {
		return results, fmt.Errorf("failed to decode logs (%s): %v", url, err)
	}

	nextPosition := position
	nextRecord := record
	for _, log := range logResult.Data {
		results = append(results, log.Message)
		nextPosition = log.Position
		nextRecord = log.Record
	}
	if o.Timestamp == "" {
		o.Timestamp = strconv.FormatInt(nextPosition, 10)
		return o.getLogs(
			client, nextPosition, nextRecord, results)
	}

	return results, err
}
