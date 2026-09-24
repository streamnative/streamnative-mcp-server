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

// Package mcp contains core MCP server integrations and tools.
package mcp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/streamnative/streamnative-mcp-server/pkg/auth"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
	sncloud "github.com/streamnative/streamnative-mcp-server/sdk/sdk-apiserver"
)

const (
	// DefaultPulsarPort is the default Pulsar protocol port.
	DefaultPulsarPort = 6651
)

// getFlow creates the appropriate OAuth2 flow based on the grant type
func getFlow(issuer *auth.Issuer, grant *auth.AuthorizationGrant) (auth.Flow, error) {
	switch grant.Type {
	case auth.GrantTypeClientCredentials:
		// Use client credentials flow for service accounts
		if grant.ClientCredentials == nil {
			return nil, fmt.Errorf("client credentials grant missing required credentials")
		}
		return auth.NewDefaultClientCredentialsFlowWithKeyFileStruct(*issuer, grant.ClientCredentials)
	default:
		return nil, fmt.Errorf("unsupported grant type: %s", grant.Type)
	}
}

// getBasePath constructs the base path for Pulsar admin API
func getBasePath(proxyLocation, org, clusterID string) string {
	// Ensure proxyLocation doesn't end with a slash to prevent double slashes
	if proxyLocation != "" && proxyLocation[len(proxyLocation)-1] == '/' {
		proxyLocation = proxyLocation[:len(proxyLocation)-1]
	}

	return fmt.Sprintf("%s/pulsar-admin/%s/pulsarcluster-%s", proxyLocation, org, clusterID)
}

// getServiceURL constructs the service URL for Pulsar protocol
func getServiceURL(dnsName string) string {
	return fmt.Sprintf("pulsar+ssl://%s:%d", dnsName, DefaultPulsarPort)
}

func getLegacyIssuer(instance *sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstance,
	configIssuer auth.Issuer,
) (*auth.Issuer, error) {
	if instance.Status == nil {
		return nil, fmt.Errorf("PulsarInstance '%s' has no auth configuration", *instance.Metadata.Name)
	}

	if instance.Status.Auth.Type != "oauth2" && instance.Status.Auth.Type != "apikey" {
		return nil, fmt.Errorf("PulsarInstance '%s' has unsupported auth type: %s",
			*instance.Metadata.Name, instance.Status.Auth.Type)
	}

	if instance.Status.Auth.Oauth2.Audience == "" || instance.Status.Auth.Oauth2.IssuerURL == "" {
		return nil, fmt.Errorf("PulsarInstance '%s' has no OAuth2 configuration", *instance.Metadata.Name)
	}

	// Construct issuer information from instance and config
	return &auth.Issuer{
		IssuerEndpoint: instance.Status.Auth.Oauth2.IssuerURL,
		ClientID:       configIssuer.ClientID,
		Audience:       instance.Status.Auth.Oauth2.Audience,
	}, nil
}

// resolveLegacyPulsarCluster implements the standalone status-based contract.
// More specialized resource/authentication policies belong to injected resolvers.
func resolveLegacyPulsarCluster(ctx context.Context, client *sncloud.APIClient, organization, instanceName, clusterName string,
	issuer auth.Issuer,
) (*sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarCluster, *auth.Issuer, error) {
	if organization == "" {
		return nil, nil, fmt.Errorf("organization is required")
	}
	api := client.CloudStreamnativeIoV1alpha1Api
	cluster, response, err := api.ReadCloudStreamnativeIoV1alpha1NamespacedPulsarCluster(ctx, clusterName, organization).Execute()
	closeAuthResponse(response)
	if err != nil {
		return nil, nil, fmt.Errorf("get PulsarCluster: %w", err)
	}
	if cluster == nil || cluster.Spec == nil {
		return nil, nil, fmt.Errorf("PulsarCluster has no spec")
	}
	if err := validateLegacyMetadata(cluster.Metadata, clusterName, organization); err != nil {
		return nil, nil, err
	}
	if cluster.Spec.InstanceName != instanceName {
		return nil, nil, fmt.Errorf("PulsarCluster belongs to a different instance")
	}
	if cluster.Metadata.Uid == nil || *cluster.Metadata.Uid == "" || cluster.Status == nil || !common.IsClusterAvailable(*cluster) {
		return nil, nil, fmt.Errorf("PulsarCluster is not available")
	}
	instance, response, err := api.ReadCloudStreamnativeIoV1alpha1NamespacedPulsarInstance(ctx, instanceName, organization).Execute()
	closeAuthResponse(response)
	if err != nil {
		return nil, nil, fmt.Errorf("get PulsarInstance: %w", err)
	}
	if instance == nil {
		return nil, nil, fmt.Errorf("empty PulsarInstance response")
	}
	if err := validateLegacyMetadata(instance.Metadata, instanceName, organization); err != nil {
		return nil, nil, err
	}
	runtimeIssuer, err := getLegacyIssuer(instance, issuer)
	return cluster, runtimeIssuer, err
}

func validateLegacyMetadata(metadata *sncloud.V1ObjectMeta, name, organization string) error {
	if metadata == nil || metadata.Name == nil || *metadata.Name != name ||
		metadata.Namespace == nil || *metadata.Namespace != organization {
		return fmt.Errorf("resource metadata does not match %s/%s", organization, name)
	}
	return nil
}

func closeAuthResponse(response *http.Response) {
	if response != nil && response.Body != nil {
		_ = response.Body.Close()
	}
}
