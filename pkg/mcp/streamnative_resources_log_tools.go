// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

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

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
)

var FunctionConnectorList = []string{"sink", "source", "function", "kafka-connect"}

// StreamNativeAddLogTools adds log tools
func StreamNativeAddLogTools(s *server.MCPServer, _ bool, features []string) {
	if !slices.Contains(features, string(FeatureStreamNativeCloud)) && !slices.Contains(features, string(FeatureAll)) {
		return
	}

	logTool := mcp.NewTool("streamnative_resources_log",
		mcp.WithDescription("Display the logs of resources in StreamNative Cloud, including pulsar functions, pulsar source connectors, pulsar sink connectors, and kafka connect connectors logs running along with PulsarInstance and PulsarCluster."+
			"This tool is used to help you debug the issues of resources in StreamNative Cloud. You can use `streamnative_cloud_context_use_cluster` to change the context to a specific cluster first, then use this tool to get the logs of resources in the cluster. This tool is suggested to be used with 'pulsar_admin_functions', 'pulsar_admin_sinks', 'pulsar_admin_sources', and 'kafka_admin_connect'"),
		mcp.WithString("component", mcp.Required(),
			mcp.Description("The component to get logs from, including "+strings.Join(FunctionConnectorList, ", ")),
			mcp.Enum(FunctionConnectorList...),
		),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("The name of the resource to get logs from."),
		),
		mcp.WithString("tenant", mcp.Required(),
			mcp.Description("The pulsar tenant of the resource to get logs from. This is required for pulsar functions, pulsar source connectors, pulsar sink connectors. For kafka connect connectors, this is optional, and the default value is 'public'."),
			mcp.DefaultString("public"),
		),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("The pulsar namespace of the resource to get logs from. This is required for pulsar functions, pulsar source connectors, pulsar sink connectors. For kafka connect connectors, this is optional, and the default value is 'default'."),
			mcp.DefaultString("default"),
		),
		mcp.WithString("size",
			mcp.Description("String format of the number of lines to get from the logs, e.g. 10, 100, etc. (default: 20)"),
			mcp.DefaultString("20"),
		),
		mcp.WithNumber("replica_id",
			mcp.Description("The replica index of the resource to get logs from, this is used for multiple replicas of running pulsar functions, pulsar source connectors, pulsar sink connectors, and kafka connect connectors. The value should be a positive integer (like 0, 1, 2, etc.), and the default value is -1, which means all replicas."),
			mcp.DefaultNumber(-1),
		),
		mcp.WithString("timestamp",
			mcp.Description("Start timestamp of logs, for example: 1662430984225"),
		),
		mcp.WithString("since",
			mcp.Description("Since time of logs, numbers end with s|m|h, for example one hour ago: 1h"),
		),
		mcp.WithBoolean("previous_container",
			mcp.Description("Return previous terminated container logs, defaults to false."),
			mcp.DefaultBool(false),
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

func handleStreamNativeResourcesLog(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	snConfig := getOptions(ctx)
	instance, cluster, organization := GetMcpContext()
	if instance == "" || cluster == "" || organization == "" {
		return mcp.NewToolResultError("No context is set, please use `streamnative_cloud_context_use_cluster` to set the context first."), nil
	}

	// Extract required parameters with validation
	component, err := requiredParam[string](request.Params.Arguments, "component")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get component: %v", err)), nil
	}

	name, err := requiredParam[string](request.Params.Arguments, "name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get name: %v", err)), nil
	}

	tenant, hasTenant := optionalParam[string](request.Params.Arguments, "tenant")
	if !hasTenant {
		tenant = "public"
	}

	namespace, hasNamespace := optionalParam[string](request.Params.Arguments, "namespace")
	if !hasNamespace {
		namespace = "default"
	}

	size, hasSize := optionalParam[string](request.Params.Arguments, "size")
	if !hasSize {
		size = "20"
	}

	replicaID, hasreplicaID := optionalParam[int](request.Params.Arguments, "replica_id")
	if !hasreplicaID {
		replicaID = -1
	}

	timestampStr, hasTimestamp := optionalParam[string](request.Params.Arguments, "timestamp")
	if !hasTimestamp {
		timestampStr = ""
	}

	sinceStr, hasSince := optionalParam[string](request.Params.Arguments, "since")
	if !hasSince {
		sinceStr = ""
	}

	previousContainer, hasPreviousContainer := optionalParam[bool](request.Params.Arguments, "previous_container")
	if !hasPreviousContainer {
		previousContainer = false
	}

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
		ServiceURL:                   snConfig.LogLocation,
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

	client, err := config.GetSNCloudLogClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get logging client: %v", err)), nil
	}

	results := []string{}
	results, err = logOptions.getLogs(client, 0, 0, results)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get logs: %v", err)), nil
	}

	if len(results) == 0 {
		return mcp.NewToolResultText("No logs found"), nil
	}

	return mcp.NewToolResultText(strings.Join(results, "\n")), nil
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
