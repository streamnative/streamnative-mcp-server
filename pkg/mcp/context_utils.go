// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
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
	"fmt"

	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
	sncloud "github.com/streamnative/streamnative-mcp-server/sdk/sdk-apiserve"
)

const DefaultKafkaPort = 9093

func SetContext(options *config.Options, instanceName, clusterName string) error {
	snConfig := options.LoadConfigOrDie()
	myselfGrant, err := options.AuthOptions.LoadGrant(snConfig.Auth.Audience)
	ctx := context.Background()
	if err != nil || myselfGrant == nil {
		return fmt.Errorf("failed to auth to StreamNative Cloud: %v", err)
	}

	apiClient, err := config.GetAPIClient()
	if err != nil {
		return fmt.Errorf("failed to get API client: %v", err)
	}

	instances, _, err := apiClient.CloudStreamnativeIoV1alpha1Api.ListCloudStreamnativeIoV1alpha1NamespacedPulsarInstance(ctx, options.Organization).Execute()
	if err != nil {
		return fmt.Errorf("failed to list pulsar instances: %v", err)
	}

	var instance sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstance
	foundInstance := false
	for _, i := range instances.Items {
		if *i.Metadata.Name == instanceName {
			if isInstanceValid(i) {
				instance = i
				foundInstance = true
				break
			}
			return fmt.Errorf("Pulsar instance %s is not valid", instanceName)
		}
	}
	if !foundInstance {
		return fmt.Errorf("Pulsar instance %s not found", instanceName)
	}

	clusters, _, err := apiClient.CloudStreamnativeIoV1alpha1Api.ListCloudStreamnativeIoV1alpha1NamespacedPulsarCluster(ctx, options.Organization).Execute()
	if err != nil {
		return fmt.Errorf("failed to list pulsar clusters: %v", err)
	}

	var cluster sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarCluster
	foundCluster := false
	for _, c := range clusters.Items {
		if *c.Metadata.Name == clusterName && c.Spec.InstanceName == instanceName {
			if isClusterAvailable(c) {
				cluster = c
				foundCluster = true
				break
			}
			return fmt.Errorf("Pulsar cluster %s is not available", clusterName)
		}
	}
	if !foundCluster {
		return fmt.Errorf("Pulsar cluster %s not found", clusterName)
	}

	clusterUID := string(*cluster.Metadata.Uid)
	dnsName := ""
	for _, endpoint := range cluster.Spec.ServiceEndpoints {
		if *endpoint.Type == "service" {
			dnsName = endpoint.DnsName
			break
		}
	}

	if dnsName == "" {
		return fmt.Errorf("no valid service endpoint found for PulsarCluster '%s'", clusterName)
	}

	issuer, err := getIssuer(&instance, snConfig.Auth.Issuer())
	if err != nil {
		return fmt.Errorf("failed to get issuer: %v", err)
	}

	tokenKey := buildTokenKey(options.Organization, clusterUID, myselfGrant.ClientCredentials.ClientEmail)

	accessToken := ""
	refreshToken := true
	cachedGrant, err := options.AuthOptions.LoadGrant(tokenKey)
	if err == nil && cachedGrant != nil {

		cacheHasValidToken, err := hasCachedValidToken(cachedGrant)
		if err != nil {
			cacheHasValidToken = false
		}

		if cacheHasValidToken {
			tokenAboutToExpire, err := isTokenAboutToExpire(cachedGrant, TokenRefreshWindow)
			if err != nil {
				tokenAboutToExpire = true
			}

			if !tokenAboutToExpire {
				refreshToken = false
				accessToken = cachedGrant.Token.AccessToken
			}
		}
	}

	if refreshToken {
		flow, err := getFlow(issuer, myselfGrant)
		if err != nil {
			return fmt.Errorf("failed to get flow: %v", err)
		}

		newGrant, err := flow.Authorize()
		if err != nil {
			return fmt.Errorf("failed to authorize: %v", err)
		}

		if newGrant.Token != nil {
			_ = options.AuthOptions.SaveGrant(tokenKey, *newGrant)
			accessToken = newGrant.Token.AccessToken
		}
	}

	if accessToken == "" {
		return fmt.Errorf("failed to get access token")
	}

	err = pulsar.NewCurrentPulsarContext(pulsar.PulsarContext{
		WebServiceURL: getBasePath(snConfig.ProxyLocation, options.Organization, clusterUID),
		ServiceURL:    getServiceURL(dnsName),
		Token:         accessToken,
	})
	if err != nil {
		return fmt.Errorf("failed to set pulsar context: %v", err)
	}

	err = kafka.NewCurrentKafkaContext(kafka.KafkaContext{
		BootstrapServers:       fmt.Sprintf("%s:%d", dnsName, DefaultKafkaPort),
		SchemaRegistryURL:      fmt.Sprintf("https://%s/kafka", dnsName),
		ConnectURL:             fmt.Sprintf("%s/admin/kafkaconnect/", snConfig.ProxyLocation),
		AuthType:               "sasl_ssl",
		AuthMechanism:          "PLAIN",
		AuthUser:               "public/default",
		AuthPass:               "token:" + accessToken,
		UseTLS:                 true,
		SchemaRegistryAuthUser: "public/default",
		SchemaRegistryAuthPass: accessToken,
		ConnectAuthUser:        "public/default",
		ConnectAuthPass:        accessToken,
	})
	if err != nil {
		return fmt.Errorf("failed to set kafka context: %v", err)
	}

	return nil
}
