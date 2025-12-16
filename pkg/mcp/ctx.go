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

	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	internalContext "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// WithSNCloudOrganization sets the SNCloud organization in the context
func WithSNCloudOrganization(ctx context.Context, organization string) context.Context {
	return internalContext.WithSNCloudOrganization(ctx, organization)
}

// WithSNCloudInstance sets the SNCloud instance in the context
func WithSNCloudInstance(ctx context.Context, instance string) context.Context {
	return internalContext.WithSNCloudInstance(ctx, instance)
}

// WithSNCloudCluster sets the SNCloud cluster in the context
func WithSNCloudCluster(ctx context.Context, cluster string) context.Context {
	return internalContext.WithSNCloudCluster(ctx, cluster)
}

// WithSNCloudSession sets the SNCloud session in the context
func WithSNCloudSession(ctx context.Context, session *config.Session) context.Context {
	return internalContext.WithSNCloudSession(ctx, session)
}

// WithPulsarSession sets the Pulsar session in the context
func WithPulsarSession(ctx context.Context, session *pulsar.Session) context.Context {
	return internalContext.WithPulsarSession(ctx, session)
}

// WithKafkaSession sets the Kafka session in the context
func WithKafkaSession(ctx context.Context, session *kafka.Session) context.Context {
	return internalContext.WithKafkaSession(ctx, session)
}

// GetSNCloudOrganization gets the SNCloud organization from the context
func GetSNCloudOrganization(ctx context.Context) string {
	return internalContext.GetSNCloudOrganization(ctx)
}

// GetSNCloudInstance gets the SNCloud instance from the context
func GetSNCloudInstance(ctx context.Context) string {
	return internalContext.GetSNCloudInstance(ctx)
}

// GetSNCloudCluster gets the SNCloud cluster from the context
func GetSNCloudCluster(ctx context.Context) string {
	return internalContext.GetSNCloudCluster(ctx)
}

// GetSNCloudSession gets the SNCloud session from the context
func GetSNCloudSession(ctx context.Context) *config.Session {
	return internalContext.GetSNCloudSession(ctx)
}

// GetPulsarSession gets the Pulsar session from the context
func GetPulsarSession(ctx context.Context) *pulsar.Session {
	return internalContext.GetPulsarSession(ctx)
}

// GetKafkaSession gets the Kafka session from the context
func GetKafkaSession(ctx context.Context) *kafka.Session {
	return internalContext.GetKafkaSession(ctx)
}
