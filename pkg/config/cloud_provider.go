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

package config

import (
	"context"
	"fmt"
	"net/http"

	"github.com/streamnative/streamnative-mcp-server/pkg/auth/store"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
	"golang.org/x/oauth2"
)

// ResolvedRuntime describes protocol targets, not live clients. A resolver must
// authorize and validate the selected scope before returning these values.
// Kafka is nil for a Pulsar-only target. TokenSources must be concurrency-safe.
type ResolvedRuntime struct {
	Pulsar pulsar.PulsarContext
	Kafka  *kafka.KafkaContext
}

// ClusterResolver is bound to one organization and effective identity. It must
// not mutate the active binding or perform interactive authentication.
type ClusterResolver interface {
	ResolveCluster(context.Context, string, string) (ResolvedRuntime, error)
}

// CloudProvider supplies embedding-owned authentication. It is process-local
// configuration: never serialize it, persist it, or expose it in tool results.
type CloudProvider struct {
	Identity    string
	TokenSource oauth2.TokenSource
	Resolver    ClusterResolver
	// Transport optionally supplies TLS/proxy configuration without authentication.
	Transport http.RoundTripper
}

// IsCloudConfigured includes both standalone and embedding-owned Cloud modes.
func (o *Options) IsCloudConfigured() bool { return o.CloudProvider != nil || o.KeyFile != "" }

// CloudIdentity reports the effective identity, without requiring a grant store
// from an embedding application.
func (o *Options) CloudIdentity() (string, error) {
	if o.CloudProvider != nil {
		return o.CloudProvider.Identity, nil
	}
	if o.Store == nil {
		return "", store.ErrNoAuthenticationData
	}
	return o.WhoAmI(o.Audience)
}

func (o *Options) validateCloudProvider() error {
	p := o.CloudProvider
	if o.KeyFile != "" || o.UseExternalKafka || o.UseExternalPulsar {
		return fmt.Errorf("injected Cloud authentication cannot be combined with key-file or external modes")
	}
	if p.Identity == "" || p.TokenSource == nil || p.Resolver == nil || o.Organization == "" || o.Server == "" {
		return fmt.Errorf("injected Cloud mode requires identity, token source, resolver, organization and server")
	}
	return nil
}
