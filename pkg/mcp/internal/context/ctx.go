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

// Package context provides internal context helpers for MCP sessions.
package context //nolint:revive

import (
	"context"
	"reflect"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
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
	MCPRequestContextKey          contextKey = "mcp_request"
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

// WithMCPRequest sets the MCP request in the context.
func WithMCPRequest(ctx context.Context, request sdk.Request) context.Context {
	return context.WithValue(ctx, MCPRequestContextKey, request)
}

// GetMCPRequest gets the MCP request from the context.
func GetMCPRequest(ctx context.Context) sdk.Request {
	if val := ctx.Value(MCPRequestContextKey); val != nil {
		if request, ok := val.(sdk.Request); ok {
			return request
		}
	}
	return nil
}

// GetMCPSession gets the MCP session from the context.
func GetMCPSession(ctx context.Context) sdk.Session {
	request := GetMCPRequest(ctx)
	if request == nil {
		return nil
	}
	session := request.GetSession()
	if session == nil {
		return nil
	}
	value := reflect.ValueOf(session)
	if value.Kind() == reflect.Ptr && value.IsNil() {
		return nil
	}
	return session
}

// GetMCPSessionID gets the MCP session ID from the context.
func GetMCPSessionID(ctx context.Context) string {
	session := GetMCPSession(ctx)
	if session == nil {
		return ""
	}
	return session.ID()
}

// GetMCPRequestExtra gets the MCP request extra from the context.
func GetMCPRequestExtra(ctx context.Context) *sdk.RequestExtra {
	request := GetMCPRequest(ctx)
	if request == nil {
		return nil
	}
	return request.GetExtra()
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
