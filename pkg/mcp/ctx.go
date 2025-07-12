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
	"context"

	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

type contextKey string

const (
	SNCloudOrganizationContextKey contextKey = "sncloud_organization"
	SNCloudInstanceContextKey     contextKey = "sncloud_instance"
	SNCloudClusterContextKey      contextKey = "sncloud_cluster"
	SNCloudSessionContextKey      contextKey = "sncloud_session"
	PulsarSessionContextKey       contextKey = "pulsar_session"
	KafkaSessionContextKey        contextKey = "kafka_session"
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

// GetSNCloudOrganization gets the SNCloud organization from the context
func GetSNCloudOrganization(ctx context.Context) string {
	return ctx.Value(SNCloudOrganizationContextKey).(string)
}

// GetSNCloudInstance gets the SNCloud instance from the context
func GetSNCloudInstance(ctx context.Context) string {
	return ctx.Value(SNCloudInstanceContextKey).(string)
}

// GetSNCloudCluster gets the SNCloud cluster from the context
func GetSNCloudCluster(ctx context.Context) string {
	return ctx.Value(SNCloudClusterContextKey).(string)
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
	session, ok := ctx.Value(PulsarSessionContextKey).(*pulsar.Session)
	if !ok {
		return nil
	}
	return session
}

// GetKafkaSession gets the Kafka session from the context
func GetKafkaSession(ctx context.Context) *kafka.Session {
	session, ok := ctx.Value(KafkaSessionContextKey).(*kafka.Session)
	if !ok {
		return nil
	}
	return session
}
