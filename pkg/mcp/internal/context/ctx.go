// Copyright 2026 StreamNative
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

// Package context provides internal context helpers for MCP sessions.
package context //nolint:revive

import (
	"context"

	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

type contextKey string

// Context keys for StreamNative sessions and identifiers.
const (
	SNCloudOrganizationContextKey contextKey = "sncloud_organization"
	SNCloudInstanceContextKey     contextKey = "sncloud_instance"
	SNCloudClusterContextKey      contextKey = "sncloud_cluster"
	SNCloudSessionContextKey      contextKey = "sncloud_session"
	PulsarSessionContextKey       contextKey = "pulsar_session"
	KafkaSessionContextKey        contextKey = "kafka_session"
	runtimeSnapshotKey            contextKey = "runtime_snapshot"
)

// WithSNCloudOrganization sets the SNCloud organization in the context
func WithSNCloudOrganization(ctx context.Context, organization string) context.Context {
	return context.WithValue(ctx, SNCloudOrganizationContextKey, organization)
}

// WithSNCloudInstance sets the SNCloud instance in the context
func WithSNCloudInstance(ctx context.Context, instance string) context.Context {
	return context.WithValue(ctx, SNCloudInstanceContextKey, instance)
}

// WithSNCloudCluster sets the SNCloud cluster in the context
func WithSNCloudCluster(ctx context.Context, cluster string) context.Context {
	return context.WithValue(ctx, SNCloudClusterContextKey, cluster)
}

// WithSNCloudSession sets the SNCloud session in the context
func WithSNCloudSession(ctx context.Context, session *config.Session) context.Context {
	return context.WithValue(ctx, SNCloudSessionContextKey, session)
}

// WithPulsarSession sets the Pulsar session in the context
func WithPulsarSession(ctx context.Context, session *pulsar.Session) context.Context {
	return context.WithValue(ctx, PulsarSessionContextKey, session)
}

// WithKafkaSession sets the Kafka session in the context
func WithKafkaSession(ctx context.Context, session *kafka.Session) context.Context {
	return context.WithValue(ctx, KafkaSessionContextKey, session)
}

type runtimeSnapshot struct{ binding *config.RuntimeBinding }

// AcquireRuntimeContext pins the selected generation for one request. Mutation
// handlers must not acquire this lease, since publishing requires a write lock.
func AcquireRuntimeContext(ctx context.Context) (context.Context, func()) {
	session := GetSNCloudSession(ctx)
	if session == nil {
		return ctx, func() {}
	}
	binding, release := session.AcquireRuntimeBinding()
	if binding == nil && session.Ctx.Organization == "" {
		// External modes have their own client lifetime and no Cloud binding.
		release()
		return ctx, func() {}
	}
	return context.WithValue(ctx, runtimeSnapshotKey, runtimeSnapshot{binding}), release
}

func runtimeBinding(ctx context.Context) *config.RuntimeBinding {
	if snapshot, ok := ctx.Value(runtimeSnapshotKey).(runtimeSnapshot); ok {
		return snapshot.binding
	}
	if session := GetSNCloudSession(ctx); session != nil {
		return session.CurrentRuntimeBinding()
	}
	return nil
}

// GetSNCloudOrganization gets the SNCloud organization from the context
func GetSNCloudOrganization(ctx context.Context) string {
	if val := ctx.Value(SNCloudOrganizationContextKey); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// GetSNCloudInstance gets the SNCloud instance from the context
func GetSNCloudInstance(ctx context.Context) string {
	if val := ctx.Value(SNCloudInstanceContextKey); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// GetSNCloudCluster gets the SNCloud cluster from the context
func GetSNCloudCluster(ctx context.Context) string {
	if val := ctx.Value(SNCloudClusterContextKey); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// GetSNCloudSession gets the SNCloud session from the context
func GetSNCloudSession(ctx context.Context) *config.Session {
	session, ok := ctx.Value(SNCloudSessionContextKey).(*config.Session)
	if !ok {
		return nil
	}
	return session
}

// GetPulsarSession gets the Pulsar session from the context
func GetPulsarSession(ctx context.Context) *pulsar.Session {
	if binding := runtimeBinding(ctx); binding != nil {
		return binding.Pulsar
	}
	session, ok := ctx.Value(PulsarSessionContextKey).(*pulsar.Session)
	if !ok {
		return nil
	}
	return session
}

// GetKafkaSession gets the Kafka session from the context
func GetKafkaSession(ctx context.Context) *kafka.Session {
	if binding := runtimeBinding(ctx); binding != nil {
		return binding.Kafka
	}
	session, ok := ctx.Value(KafkaSessionContextKey).(*kafka.Session)
	if !ok {
		return nil
	}
	return session
}
