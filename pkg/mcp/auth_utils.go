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
	"fmt"

	"github.com/streamnative/streamnative-mcp-server/pkg/auth"
	sncloud "github.com/streamnative/streamnative-mcp-server/sdk/sdk-apiserver"
)

const (
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

func getIssuer(instance *sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstance, configIssuer auth.Issuer) (*auth.Issuer, error) {
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
