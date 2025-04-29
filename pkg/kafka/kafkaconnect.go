// Copyright (c) 2025 StreamNative, Inc.. All Rights Reserved.

package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	kafkaconnect "github.com/streamnative/streamnative-mcp-server/sdk/sdk-kafkaconnect"
)

// ConnectorState represents the state of a connector
type ConnectorState string

const (
	// ConnectorStateRunning means the connector is running
	ConnectorStateRunning ConnectorState = "RUNNING"
	// ConnectorStatePaused means the connector is paused
	ConnectorStatePaused ConnectorState = "PAUSED"
	// ConnectorStateFailed means the connector has failed
	ConnectorStateFailed ConnectorState = "FAILED"
	// ConnectorStateStopped means the connector is stopped
	ConnectorStateStopped ConnectorState = "STOPPED"
)

// ConnectorInfo contains information about a connector
type ConnectorInfo struct {
	// Name is the connector name
	Name string
	// Config is the connector configuration
	Config map[string]string
	// Tasks is the list of tasks
	Tasks []TaskInfo
	// Type is the connector type (source or sink)
	Type string
	// State is the connector state
	State ConnectorState
}

// TaskInfo contains information about a connector task
type TaskInfo struct {
	// ID is the task ID
	ID int
	// State is the task state
	State ConnectorState
	// WorkerID is the worker ID
	WorkerID string
	// Trace is the error trace (if any)
	Trace string
}

// PluginInfo contains information about a connector plugin
type PluginInfo struct {
	// Class is the plugin class
	Class string
	// Type is the plugin type (source or sink)
	Type string
	// Version is the plugin version
	Version string
}

// Connect is the interface for Kafka Connect operations
type Connect interface {
	// GetInfo gets information about the Kafka Connect cluster
	GetInfo(ctx context.Context) (map[string]interface{}, error)
	// ListConnectors lists all connectors
	ListConnectors(ctx context.Context) ([]string, error)
	// GetConnector gets information about a connector
	GetConnector(ctx context.Context, name string) (*ConnectorInfo, error)
	// CreateConnector creates a new connector
	CreateConnector(ctx context.Context, name string, config map[string]string) (*ConnectorInfo, error)
	// UpdateConnector updates a connector
	UpdateConnector(ctx context.Context, name string, config map[string]string) (*ConnectorInfo, error)
	// DeleteConnector deletes a connector
	DeleteConnector(ctx context.Context, name string) error
	// PauseConnector pauses a connector
	PauseConnector(ctx context.Context, name string) error
	// ResumeConnector resumes a connector
	ResumeConnector(ctx context.Context, name string) error
	// RestartConnector restarts a connector
	RestartConnector(ctx context.Context, name string) error
	// GetConnectorConfig gets the configuration of a connector
	GetConnectorConfig(ctx context.Context, name string) (map[string]string, error)
	// GetConnectorStatus gets the status of a connector
	GetConnectorStatus(ctx context.Context, name string) (*ConnectorInfo, error)
	// GetConnectorTasks gets the tasks of a connector
	GetConnectorTasks(ctx context.Context, name string) ([]TaskInfo, error)
	// ListPlugins lists all connector plugins
	ListPlugins(ctx context.Context) ([]PluginInfo, error)
	// ValidateConfig validates a connector configuration
	ValidateConfig(ctx context.Context, pluginClass string, config map[string]string) (map[string]interface{}, error)
}

// connectImpl is the implementation of the Connect interface
type connectImpl struct {
	client  *kafkaconnect.APIClient
	ctx     context.Context
	baseURL string
}

// NewConnect creates a new Kafka Connect client
func NewConnect(kc *KafkaContext) (Connect, error) {
	// Create a new configuration for the SDK client
	cfg := kafkaconnect.NewConfiguration()

	// Parse and set the base URL from config
	baseURL := kc.ConnectURL
	baseURL = strings.TrimSuffix(baseURL, "/") // Remove trailing slash if present

	// Parse the URL to extract scheme, host, and path
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid Kafka ConnectURL format: %v", err)
	}

	// Set the scheme and host for the SDK client
	cfg.Host = parsedURL.Host
	cfg.Scheme = parsedURL.Scheme

	// Set the server URL directly in the Servers configuration
	// This allows us to include the path component in the URL
	serverURL := fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, parsedURL.Path)
	cfg.Servers = []kafkaconnect.ServerConfiguration{
		{
			URL:         serverURL,
			Description: "Kafka Connect server",
		},
	}
	ctx := context.Background()

	if kc.ConnectAuthUser != "" && kc.ConnectAuthPass != "" {
		cred := kafkaconnect.BasicAuth{
			UserName: kc.ConnectAuthUser,
			Password: kc.ConnectAuthPass,
		}
		ctx = context.WithValue(ctx, kafkaconnect.ContextBasicAuth, cred)
	}

	// Create API client and context
	apiClient := kafkaconnect.NewAPIClient(cfg)

	return &connectImpl{
		client:  apiClient,
		ctx:     ctx,
		baseURL: baseURL,
	}, nil
}

// GetInfo gets information about the Kafka Connect cluster
func (c *connectImpl) GetInfo(_ context.Context) (map[string]interface{}, error) {
	// Make request
	serverInfo, _, err := c.client.DefaultAPI.ServerInfo(c.ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get Kafka Connect server info: %w", err)
	}

	// Convert to map
	info := make(map[string]interface{})
	if serverInfo.HasVersion() {
		info["version"] = serverInfo.GetVersion()
	}
	if serverInfo.HasKafkaClusterId() {
		info["kafka_cluster_id"] = serverInfo.GetKafkaClusterId()
	}
	if serverInfo.HasCommit() {
		info["commit"] = serverInfo.GetCommit()
	}

	return info, nil
}

// ListConnectors lists all connectors
func (c *connectImpl) ListConnectors(_ context.Context) ([]string, error) {
	// Make request
	resp, err := c.client.DefaultAPI.ListConnectors(c.ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list connectors: %w", err)
	}
	defer resp.Body.Close()

	// Parse response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Unmarshal JSON
	var connectors []string
	if err := json.Unmarshal(body, &connectors); err != nil {
		return nil, fmt.Errorf("failed to parse connectors: %w", err)
	}

	return connectors, nil
}

// GetConnector gets information about a connector
func (c *connectImpl) GetConnector(ctx context.Context, name string) (*ConnectorInfo, error) {
	// Make request
	info, _, err := c.client.DefaultAPI.GetConnector(c.ctx, name).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get connector: %w", err)
	}

	// Get status
	status, err := c.GetConnectorStatus(ctx, name)
	if err != nil {
		return nil, err
	}

	// Create connector info
	connector := &ConnectorInfo{
		Name:   name,
		Config: info.GetConfig(),
		Type:   info.GetType(),
		State:  status.State,
		Tasks:  status.Tasks,
	}

	return connector, nil
}

// CreateConnector creates a new connector
func (c *connectImpl) CreateConnector(_ context.Context, name string, config map[string]string) (*ConnectorInfo, error) {
	// Create request payload
	payload := *kafkaconnect.NewCreateConnectorRequest()
	payload.SetName(name)
	payload.SetConfig(config)

	// Make request
	resp, err := c.client.DefaultAPI.CreateConnector(c.ctx).CreateConnectorRequest(payload).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to create connector: %w", err)
	}
	defer resp.Body.Close()

	// Parse response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Unmarshal JSON
	var connectorInfo struct {
		Name   string            `json:"name"`
		Config map[string]string `json:"config"`
		Tasks  []struct {
			Connector string `json:"connector"`
			Task      int    `json:"task"`
		} `json:"tasks"`
		Type string `json:"type"`
	}

	if err := json.Unmarshal(body, &connectorInfo); err != nil {
		return nil, fmt.Errorf("failed to parse connector info: %w", err)
	}

	// Create connector info
	connector := &ConnectorInfo{
		Name:   name,
		Config: connectorInfo.Config,
		Type:   connectorInfo.Type,
	}

	return connector, nil
}

// UpdateConnector updates a connector
func (c *connectImpl) UpdateConnector(_ context.Context, name string, config map[string]string) (*ConnectorInfo, error) {
	// Create request payload
	// Convert config to map[string]string as required by the SDK
	stringConfig := make(map[string]string)
	stringConfig["connector.class"] = config["connector.class"]
	stringConfig["tasks.max"] = config["tasks.max"]
	for k, v := range config {
		if k != "name" && k != "connector.class" && k != "tasks.max" {
			stringConfig[k] = v
		}
	}

	// Make request
	resp, err := c.client.DefaultAPI.PutConnectorConfig(c.ctx, name).RequestBody(stringConfig).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to update connector: %w", err)
	}
	defer resp.Body.Close()

	// Parse response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Unmarshal JSON
	var connectorInfo struct {
		Name   string            `json:"name"`
		Config map[string]string `json:"config"`
		Tasks  []struct {
			Connector string `json:"connector"`
			Task      int    `json:"task"`
		} `json:"tasks"`
		Type string `json:"type"`
	}

	if err := json.Unmarshal(body, &connectorInfo); err != nil {
		return nil, fmt.Errorf("failed to parse connector info: %w", err)
	}

	// Create connector info
	connector := &ConnectorInfo{
		Name:   name,
		Config: connectorInfo.Config,
		Type:   connectorInfo.Type,
	}

	return connector, nil
}

// DeleteConnector deletes a connector
func (c *connectImpl) DeleteConnector(_ context.Context, name string) error {
	// Make request
	_, err := c.client.DefaultAPI.DestroyConnector(c.ctx, name).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete connector: %w", err)
	}

	return nil
}

// PauseConnector pauses a connector
func (c *connectImpl) PauseConnector(_ context.Context, name string) error {
	// Make request
	_, err := c.client.DefaultAPI.PauseConnector(c.ctx, name).Execute()
	if err != nil {
		return fmt.Errorf("failed to pause connector: %w", err)
	}

	return nil
}

// ResumeConnector resumes a connector
func (c *connectImpl) ResumeConnector(_ context.Context, name string) error {
	// Make request
	_, err := c.client.DefaultAPI.ResumeConnector(c.ctx, name).Execute()
	if err != nil {
		return fmt.Errorf("failed to resume connector: %w", err)
	}

	return nil
}

// RestartConnector restarts a connector
func (c *connectImpl) RestartConnector(_ context.Context, name string) error {
	// Make request
	_, err := c.client.DefaultAPI.RestartConnector(c.ctx, name).Execute()
	if err != nil {
		return fmt.Errorf("failed to restart connector: %w", err)
	}

	return nil
}

// GetConnectorConfig gets the configuration of a connector
func (c *connectImpl) GetConnectorConfig(_ context.Context, name string) (map[string]string, error) {
	// Make request
	config, _, err := c.client.DefaultAPI.GetConnectorConfig(c.ctx, name).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get connector config: %w", err)
	}

	return config, nil
}

// GetConnectorStatus gets the status of a connector
func (c *connectImpl) GetConnectorStatus(_ context.Context, name string) (*ConnectorInfo, error) {
	// Make request
	_, resp, err := c.client.DefaultAPI.GetConnectorStatus(c.ctx, name).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get connector status: %w", err)
	}
	defer resp.Body.Close()

	// Parse response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Unmarshal JSON
	var statusInfo struct {
		Name      string `json:"name"`
		Connector struct {
			State    string `json:"state"`
			WorkerID string `json:"worker_id"`
		} `json:"connector"`
		Tasks []struct {
			ID       int    `json:"id"`
			State    string `json:"state"`
			WorkerID string `json:"worker_id"`
			Trace    string `json:"trace,omitempty"`
		} `json:"tasks"`
	}

	if err := json.Unmarshal(body, &statusInfo); err != nil {
		return nil, fmt.Errorf("failed to parse connector status: %w", err)
	}

	// Convert tasks
	tasks := make([]TaskInfo, len(statusInfo.Tasks))
	for i, task := range statusInfo.Tasks {
		tasks[i] = TaskInfo{
			ID:       task.ID,
			State:    ConnectorState(task.State),
			WorkerID: task.WorkerID,
			Trace:    task.Trace,
		}
	}

	// Create connector info
	connector := &ConnectorInfo{
		Name:  name,
		State: ConnectorState(statusInfo.Connector.State),
		Tasks: tasks,
	}

	return connector, nil
}

// GetConnectorTasks gets the tasks of a connector
func (c *connectImpl) GetConnectorTasks(_ context.Context, name string) ([]TaskInfo, error) {
	// Make request
	_, resp, err := c.client.DefaultAPI.GetTaskConfigs(c.ctx, name).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get connector tasks: %w", err)
	}
	defer resp.Body.Close()

	// Parse response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Unmarshal JSON
	var tasks []struct {
		ID       int    `json:"id"`
		State    string `json:"state,omitempty"`
		WorkerID string `json:"worker_id,omitempty"`
		Trace    string `json:"trace,omitempty"`
	}

	if err := json.Unmarshal(body, &tasks); err != nil {
		return nil, fmt.Errorf("failed to parse connector tasks: %w", err)
	}

	// Convert tasks
	taskInfos := make([]TaskInfo, len(tasks))
	for i, task := range tasks {
		taskInfos[i] = TaskInfo{
			ID:       task.ID,
			State:    ConnectorState(task.State),
			WorkerID: task.WorkerID,
			Trace:    task.Trace,
		}
	}

	return taskInfos, nil
}

// ListPlugins lists all connector plugins
func (c *connectImpl) ListPlugins(_ context.Context) ([]PluginInfo, error) {
	// Make request
	_, resp, err := c.client.DefaultAPI.ListConnectorPlugins(c.ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list connector plugins: %w", err)
	}
	defer resp.Body.Close()

	// Parse response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Unmarshal JSON
	var plugins []struct {
		Class   string `json:"class"`
		Type    string `json:"type"`
		Version string `json:"version"`
	}

	if err := json.Unmarshal(body, &plugins); err != nil {
		return nil, fmt.Errorf("failed to parse plugins: %w", err)
	}

	// Convert plugins
	pluginInfos := make([]PluginInfo, len(plugins))
	for i, plugin := range plugins {
		pluginInfos[i] = PluginInfo{
			Class:   plugin.Class,
			Type:    plugin.Type,
			Version: plugin.Version,
		}
	}

	return pluginInfos, nil
}

// ValidateConfig validates a connector configuration
func (c *connectImpl) ValidateConfig(_ context.Context, pluginClass string, config map[string]string) (map[string]interface{}, error) {
	// Create request payload
	// Convert config to map[string]string as required by the SDK
	stringConfig := make(map[string]string)
	stringConfig["connector.class"] = pluginClass
	for k, v := range config {
		stringConfig[k] = v
	}

	// Make request
	_, resp, err := c.client.DefaultAPI.ValidateConfigs(c.ctx, pluginClass).RequestBody(stringConfig).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to validate connector config: %w", err)
	}
	defer resp.Body.Close()

	// Parse response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Unmarshal JSON
	var validation struct {
		Name       string `json:"name"`
		ErrorCount int    `json:"error_count"`
		Configs    []struct {
			Definition struct {
				Name          string      `json:"name"`
				Type          string      `json:"type"`
				Required      bool        `json:"required"`
				DefaultValue  interface{} `json:"default_value"`
				Importance    string      `json:"importance"`
				Documentation string      `json:"documentation"`
			} `json:"definition"`
			Value struct {
				Name        string      `json:"name"`
				Value       interface{} `json:"value"`
				Recommended interface{} `json:"recommended_values"`
				Errors      []string    `json:"errors"`
				Visible     bool        `json:"visible"`
			} `json:"value"`
		} `json:"configs"`
	}

	if err := json.Unmarshal(body, &validation); err != nil {
		return nil, fmt.Errorf("failed to parse validation result: %w", err)
	}

	// Convert validation result
	result := make(map[string]interface{})
	result["name"] = validation.Name
	result["error_count"] = validation.ErrorCount
	result["configs"] = validation.Configs

	return result, nil
}
