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

package mcp

import (
	"context"
	"fmt"

	"github.com/streamnative/streamnative-mcp-server/pkg/auth"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	context2 "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
	"golang.org/x/oauth2"
)

// DefaultKafkaPort is the default Kafka port for standalone Cloud clusters.
const DefaultKafkaPort = 9093

// SetContext resolves a selection through its provider, prepares protocol clients,
// and atomically publishes them. It does not interpret injected credentials.
func SetContext(ctx context.Context, options *config.Options, instanceName, clusterName string) error {
	if options == nil || instanceName == "" || clusterName == "" {
		return fmt.Errorf("options, instance and cluster are required")
	}
	session := context2.GetSNCloudSession(ctx)
	if session == nil {
		return fmt.Errorf("failed to get StreamNative Cloud session")
	}
	unlock := session.LockRuntimeMutation()
	defer unlock()
	var resolver config.ClusterResolver = standaloneClusterResolver{options: options, session: session}
	if options.CloudProvider != nil {
		if options.CloudProvider.Resolver == nil {
			return fmt.Errorf("cloud resolver is not configured")
		}
		resolver = options.CloudProvider.Resolver
	}
	resolved, err := resolver.ResolveCluster(ctx, instanceName, clusterName)
	if err != nil {
		return err
	}
	if resolved.Pulsar.ServiceURL == "" || resolved.Pulsar.WebServiceURL == "" {
		return fmt.Errorf("resolver returned incomplete Pulsar endpoints")
	}
	candidate := &config.RuntimeBinding{
		Instance: instanceName, Cluster: clusterName, Pulsar: &pulsar.Session{}, Kafka: &kafka.Session{},
	}
	published := false
	defer func() {
		if !published {
			candidate.Close()
		}
	}()
	if err := candidate.Pulsar.SetPulsarContext(resolved.Pulsar); err != nil {
		return fmt.Errorf("prepare Pulsar context: %w", err)
	}
	if resolved.Kafka != nil {
		if err := candidate.Kafka.SetKafkaContext(*resolved.Kafka); err != nil {
			return fmt.Errorf("prepare Kafka context: %w", err)
		}
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	session.PublishRuntimeBinding(candidate)
	published = true
	return nil
}

type standaloneClusterResolver struct {
	options *config.Options
	session *config.Session
}

func (r standaloneClusterResolver) ResolveCluster(ctx context.Context,
	instanceName, clusterName string,
) (config.ResolvedRuntime, error) {
	options, session := r.options, r.session
	var resolved config.ResolvedRuntime
	if options.Store == nil {
		return resolved, fmt.Errorf("standalone Cloud authentication is not configured")
	}
	snConfig := options.LoadConfigOrDie()
	grant, err := options.LoadGrant(snConfig.Auth.Audience)
	if err != nil {
		return resolved, fmt.Errorf("load Cloud authentication: %w", err)
	}
	if grant == nil {
		return resolved, fmt.Errorf("cloud authentication is missing")
	}
	client, err := session.GetAPIClient()
	if err != nil {
		return resolved, err
	}
	cluster, issuer, err := resolveLegacyPulsarCluster(ctx, client, options.Organization, instanceName, clusterName, snConfig.Auth.Issuer())
	if err != nil {
		return resolved, err
	}
	dnsName := ""
	for _, endpoint := range cluster.Spec.ServiceEndpoints {
		if endpoint.Type != nil && *endpoint.Type == "service" {
			dnsName = endpoint.DnsName
			break
		}
	}
	if dnsName == "" {
		return resolved, fmt.Errorf("no valid service endpoint found for PulsarCluster %q", clusterName)
	}
	if grant.Type != auth.GrantTypeClientCredentials || grant.ClientCredentials == nil {
		return resolved, fmt.Errorf("standalone Cloud mode requires service account credentials")
	}
	// The standalone session has one immutable identity and no persistent
	// cross-session runtime cache. No resource-specific audience policy is applied.
	key := fmt.Sprintf("%q:%q", issuer.IssuerEndpoint, issuer.Audience)
	source := session.RuntimeTokenSource(key, func() oauth2.TokenSource {
		return oauth2.ReuseTokenSourceWithExpiry(nil, &legacyTokenRefresher{
			issuer: *issuer, credentials: *grant.ClientCredentials,
		}, common.TokenRefreshWindow)
	})
	token, err := source.Token()
	if err != nil {
		return resolved, err
	}
	resolved.Pulsar = pulsar.PulsarContext{
		WebServiceURL: getBasePath(snConfig.ProxyLocation, options.Organization, *cluster.Metadata.Uid),
		ServiceURL:    getServiceURL(dnsName), Token: token.AccessToken, TokenSource: source,
	}
	if cluster.Spec.Config != nil && cluster.Spec.Config.Protocols != nil && cluster.Spec.Config.Protocols.Kafka != nil {
		resolved.Kafka = &kafka.KafkaContext{
			BootstrapServers:  fmt.Sprintf("%s:%d", dnsName, DefaultKafkaPort),
			SchemaRegistryURL: fmt.Sprintf("https://%s/kafka", dnsName),
			ConnectURL:        fmt.Sprintf("%s/admin/kafkaconnect/", snConfig.ProxyLocation),
			AuthType:          "sasl_ssl", AuthMechanism: "PLAIN", AuthUser: "public/default", AuthPass: "token:" + token.AccessToken,
			UseTLS: true, SchemaRegistryAuthUser: "public/default", SchemaRegistryAuthPass: token.AccessToken,
			ConnectAuthUser: "public/default", ConnectAuthPass: token.AccessToken,
			TokenSource: source, TokenPrefix: "token:",
		}
	}
	return resolved, nil
}

type legacyTokenRefresher struct {
	issuer      auth.Issuer
	credentials auth.KeyFile
}

func (s *legacyTokenRefresher) Token() (*oauth2.Token, error) {
	flow, err := getFlow(&s.issuer, &auth.AuthorizationGrant{
		Type: auth.GrantTypeClientCredentials, ClientCredentials: &s.credentials,
	})
	if err != nil {
		return nil, err
	}
	grant, err := flow.Authorize()
	if err != nil {
		return nil, err
	}
	if grant == nil || grant.Token == nil || !grant.Token.Valid() || grant.Token.Expiry.IsZero() {
		return nil, fmt.Errorf("authorization returned no valid expiring access token")
	}
	return grant.Token, nil
}

// ResetContext clears StreamNative Cloud cluster bindings and protocol sessions.
func ResetContext(ctx context.Context) error {
	session := context2.GetSNCloudSession(ctx)
	if session == nil {
		return fmt.Errorf("failed to get StreamNative Cloud session")
	}
	unlock := session.LockRuntimeMutation()
	defer unlock()
	var initial *config.RuntimeBinding
	if session.CurrentRuntimeBinding() == nil {
		initial = &config.RuntimeBinding{Pulsar: context2.GetPulsarSession(ctx), Kafka: context2.GetKafkaSession(ctx)}
	}
	session.PublishRuntimeBinding(&config.RuntimeBinding{Pulsar: &pulsar.Session{}, Kafka: &kafka.Session{}})
	initial.Close()

	return nil
}
