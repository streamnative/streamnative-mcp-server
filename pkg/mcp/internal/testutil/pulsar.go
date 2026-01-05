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

//go:build e2e

package testutil

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/apache/pulsar-client-go/pulsaradmin"
	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/admin/config"
	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
)

// PulsarTestHelper provides helper functions for Pulsar E2E testing.
type PulsarTestHelper struct {
	adminClient pulsaradmin.Client
	adminURL    string
	serviceURL  string
	httpClient  *http.Client
}

// NewPulsarTestHelper creates a new PulsarTestHelper.
func NewPulsarTestHelper(adminURL, serviceURL string) (*PulsarTestHelper, error) {
	if adminURL == "" {
		adminURL = "http://localhost:8080"
	}
	if serviceURL == "" {
		serviceURL = "pulsar://localhost:6650"
	}

	cfg := &config.Config{}
	cfg.WebServiceURL = adminURL

	client, err := pulsaradmin.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create pulsar admin client: %w", err)
	}

	return &PulsarTestHelper{
		adminClient: client,
		adminURL:    adminURL,
		serviceURL:  serviceURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// Close closes the Pulsar admin client's underlying resources.
func (h *PulsarTestHelper) Close() {
	// pulsaradmin.Client interface doesn't have Close method
	// The underlying HTTP client will be garbage collected
}

// WaitForReady waits for Pulsar to be ready by checking the admin API.
func (h *PulsarTestHelper) WaitForReady(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	url := fmt.Sprintf("%s/admin/v2/clusters", h.adminURL)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for pulsar to be ready")
		case <-ticker.C:
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				continue
			}

			resp, err := h.httpClient.Do(req)
			if err == nil && resp.StatusCode == http.StatusOK {
				_ = resp.Body.Close()
				return nil
			}
			if resp != nil {
				_ = resp.Body.Close()
			}
		}
	}
}

// EnsureNamespace ensures the specified namespace exists.
func (h *PulsarTestHelper) EnsureNamespace(ctx context.Context, tenant, namespace string) error {
	fullNamespace := tenant + "/" + namespace

	// Check if namespace exists
	namespaces, err := h.adminClient.Namespaces().GetNamespaces(tenant)
	if err != nil {
		// If we cannot list namespaces, we'll try to create it anyway.
		// This could be due to permission issues or other transient errors.
		// The CreateNamespace call will fail with a proper error if there's a real issue.
	} else {
		// Successfully retrieved namespaces, check if ours exists
		for _, ns := range namespaces {
			if ns == fullNamespace {
				return nil // Already exists
			}
		}
	}

	// Create namespace - this will fail if it already exists or if there are permission/connection issues
	err = h.adminClient.Namespaces().CreateNamespace(fullNamespace)
	if err != nil {
		// Check if the error is because namespace already exists.
		// The Pulsar admin client doesn't provide typed errors, so we check the error message.
		// This is a best-effort approach for test utilities - we check for common patterns:
		// - "already exists" (common HTTP API message)
		// - "alreadyexists" (possible gRPC/protobuf message)
		// - "409" (HTTP Conflict status code)
		// - "conflict" (HTTP status text)
		// Using case-insensitive matching for robustness
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "already exists") ||
			strings.Contains(errMsg, "alreadyexists") ||
			strings.Contains(errMsg, "409") ||
			strings.Contains(errMsg, "conflict") {
			return nil // Namespace already exists, which is what we want
		}
		return fmt.Errorf("failed to create namespace %s: %w", fullNamespace, err)
	}

	return nil
}

// CreateTopic creates a Pulsar topic with the specified name and partition count.
func (h *PulsarTestHelper) CreateTopic(ctx context.Context, topic string, partitions int) error {
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return fmt.Errorf("invalid topic name %s: %w", topic, err)
	}

	// Use dereference as done in existing code (see pkg/mcp/builders/pulsar/topic.go)
	if err := h.adminClient.Topics().Create(*topicName, partitions); err != nil {
		return fmt.Errorf("failed to create topic %s: %w", topic, err)
	}
	return nil
}

// TopicExists checks if a topic exists.
func (h *PulsarTestHelper) TopicExists(ctx context.Context, topic string) (bool, error) {
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return false, nil
	}

	_, err = h.adminClient.Topics().GetMetadata(*topicName)
	if err != nil {
		return false, nil // Assume doesn't exist
	}
	return true, nil
}

// DeleteTopic deletes a Pulsar topic.
func (h *PulsarTestHelper) DeleteTopic(ctx context.Context, topic string) error {
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return fmt.Errorf("invalid topic name %s: %w", topic, err)
	}

	// Delete with force=true, treat as non-partitioned (works for both)
	if err := h.adminClient.Topics().Delete(*topicName, true, true); err != nil {
		// Try as partitioned
		if err2 := h.adminClient.Topics().Delete(*topicName, true, false); err2 != nil {
			return fmt.Errorf("failed to delete topic %s: %w", topic, err)
		}
	}
	return nil
}

// CleanupTopic cleans up a topic if it exists.
func (h *PulsarTestHelper) CleanupTopic(ctx context.Context, topic string) error {
	exists, err := h.TopicExists(ctx, topic)
	if err != nil {
		return err
	}
	if exists {
		return h.DeleteTopic(ctx, topic)
	}
	return nil
}

// ListTopics lists topics in a namespace.
func (h *PulsarTestHelper) ListTopics(namespace string) ([]string, error) {
	nsName, err := utils.GetNamespaceName(namespace)
	if err != nil {
		return nil, fmt.Errorf("invalid namespace name %s: %w", namespace, err)
	}

	topics, _, err := h.adminClient.Topics().List(*nsName)
	if err != nil {
		return nil, fmt.Errorf("failed to list topics: %w", err)
	}
	return topics, nil
}

// GetTopicMetadata gets metadata for a topic.
func (h *PulsarTestHelper) GetTopicMetadata(topic string) (map[string]interface{}, error) {
	topicName, err := utils.GetTopicName(topic)
	if err != nil {
		return nil, fmt.Errorf("invalid topic name %s: %w", topic, err)
	}

	metadata, err := h.adminClient.Topics().GetMetadata(*topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to get topic metadata: %w", err)
	}

	// PartitionedTopicMetadata only has Partitions field
	result := map[string]interface{}{
		"name":       topic,
		"partitions": metadata.Partitions,
	}
	return result, nil
}

// GenerateTestTopicName generates a unique test topic name.
func GenerateTestTopicName(prefix string) string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("persistent://public/default/%s-%d", prefix, timestamp)
}
