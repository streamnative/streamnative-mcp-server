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
	"strings"
	"time"

	"github.com/streamnative/streamnative-mcp-server/pkg/auth"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	sncloud "github.com/streamnative/streamnative-mcp-server/sdk/sdk-apiserve"
)

type ContextKey string

const (
	OptionsKey                        ContextKey = "snmcp-options"
	AnnotationStreamNativeCloudEngine            = "cloud.streamnative.io/engine"
	KeyPrefix                                    = "snmcp-token"
	TokenRefreshWindow                           = 5 * time.Minute
)

// requiredParam is a helper function that can be used to fetch a requested parameter from the request.
// It does the following checks:
// 1. Checks if the parameter is present in the request.
// 2. Checks if the parameter is of the expected type.
// 3. Checks if the parameter is not empty, i.e: non-zero value
func requiredParam[T comparable](arguments map[string]interface{}, p string) (T, error) {
	var zero T

	// Check if the parameter is present in the request
	if _, ok := arguments[p]; !ok {
		return zero, fmt.Errorf("missing required parameter: %s", p)
	}

	// Check if the parameter is of the expected type
	if _, ok := arguments[p].(T); !ok {
		return zero, fmt.Errorf("parameter %s is not of type %T", p, zero)
	}

	if arguments[p].(T) == zero {
		return zero, fmt.Errorf("missing required parameter: %s", p)

	}

	return arguments[p].(T), nil
}

// Helper function to get an optional parameter from the request
func optionalParam[T any](arguments map[string]interface{}, paramName string) (T, bool) {
	var empty T
	param, ok := arguments[paramName]
	if !ok {
		return empty, false
	}

	value, ok := param.(T)
	if !ok {
		return empty, false
	}

	return value, true
}

// Helper function to get a required array parameter from the request
func requiredParamArray[T any](arguments map[string]interface{}, paramName string) ([]T, error) {
	var empty []T
	param, ok := arguments[paramName]
	if !ok {
		return empty, fmt.Errorf("required parameter %s is missing", paramName)
	}

	paramArray, ok := param.([]interface{})
	if !ok {
		return empty, fmt.Errorf("parameter %s is not an array", paramName)
	}

	result := make([]T, 0, len(paramArray))
	for _, item := range paramArray {
		value, ok := item.(T)
		if !ok {
			return empty, fmt.Errorf("parameter %s contains items of invalid type", paramName)
		}
		result = append(result, value)
	}

	return result, nil
}

func optionalParamArray[T any](arguments map[string]interface{}, paramName string) ([]T, bool) {
	var empty []T
	param, ok := arguments[paramName]
	if !ok {
		return empty, false
	}

	paramArray, ok := param.([]interface{})
	if !ok {
		return empty, false
	}

	result := make([]T, 0, len(paramArray))
	for _, item := range paramArray {
		value, ok := item.(T)
		if !ok {
			return empty, false
		}
		result = append(result, value)
	}

	return result, true
}

// optionalParamConfigs gets an optional parameter as a list of key=value strings
func optionalParamConfigs(arguments map[string]interface{}, paramName string) ([]string, bool) {

	param, ok := arguments[paramName]
	if !ok {
		return nil, false
	}

	paramArray, ok := param.([]interface{})
	if !ok {
		return nil, false
	}

	result := make([]string, 0, len(paramArray))
	for _, item := range paramArray {
		if strValue, ok := item.(string); ok {
			result = append(result, strValue)
		}
	}

	return result, true
}

// requiredParamObject gets a required object parameter from the request
func requiredParamObject(arguments map[string]interface{}, name string) (map[string]interface{}, error) {
	// Get the parameter value
	paramValue, found := arguments[name]
	if !found || paramValue == nil {
		return nil, fmt.Errorf("%s parameter is required", name)
	}

	// Convert to map
	if mapVal, ok := paramValue.(map[string]interface{}); ok {
		return mapVal, nil
	}

	return nil, fmt.Errorf("%s parameter must be an object", name)
}

func getOptions(ctx context.Context) *config.Options {
	return ctx.Value(OptionsKey).(*config.Options)
}

// isClusterAvailable checks if a PulsarCluster is available (ready)
func isClusterAvailable(cluster sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarCluster) bool {
	// Check if broker has ready replicas
	if cluster.Status.Broker == nil || cluster.Status.Broker.ReadyReplicas == nil || *cluster.Status.Broker.ReadyReplicas == 0 {
		return false
	}

	// Check if conditions include Ready=True
	for _, condition := range cluster.Status.Conditions {
		if condition.Type == "Ready" && condition.Status == "True" {
			return true
		}
	}
	return false
}

// getEngineType returns the Pulsar cluster is an Ursa engine or a Classic engine
func getEngineType(cluster sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarCluster) string {
	if cluster.Metadata.Annotations != nil {
		if v, has := (*cluster.Metadata.Annotations)[AnnotationStreamNativeCloudEngine]; has && v == "ursa" {
			return "ursa"
		}
	}
	return "classic"
}

func convertToMapInterface(m map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		result[k] = v
	}
	return result
}

func convertToMapString(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result
}

// isInstanceValid checks if PulsarInstance has valid OAuth2 authentication configuration
func isInstanceValid(instance sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstance) bool {
	return instance.Status != nil &&
		(instance.Status.Auth.Type == "oauth2" || instance.Status.Auth.Type == "apikey") &&
		instance.Status.Auth.Oauth2.IssuerURL != "" &&
		instance.Status.Auth.Oauth2.Audience != ""
}

// buildTokenKey creates a unique key for storing the token in the cache
func buildTokenKey(namespace, clusterUID, myselfGrant string) string {
	return fmt.Sprintf("%s-%s-%s-%s", KeyPrefix, namespace, clusterUID, myselfGrant)
}

func hasCachedValidToken(cachedGrant *auth.AuthorizationGrant) (bool, error) {
	if cachedGrant == nil || cachedGrant.Token == nil {
		return false, nil
	}

	// Check if token is valid and not expired
	return cachedGrant.Token.Valid(), nil
}

func isTokenAboutToExpire(cachedGrant *auth.AuthorizationGrant, window time.Duration) (bool, error) {
	if cachedGrant == nil || cachedGrant.Token == nil {
		return true, nil
	}

	// Check if token will expire within the window
	expiry := cachedGrant.Token.Expiry
	if expiry.IsZero() {
		// If we can't determine expiry, assume it will expire soon
		return true, nil
	}

	timeUntilExpiry := time.Until(expiry)
	return timeUntilExpiry <= window, nil
}

// isPackageURLSupported checks if the provided URL protocol is supported for package download
func isPackageURLSupported(packageURL string) bool {
	// Check if the URL has a supported protocol: http, https, file
	supportedProtocols := []string{"http://", "https://", "file://"}
	for _, protocol := range supportedProtocols {
		if strings.HasPrefix(packageURL, protocol) {
			return true
		}
	}
	return false
}

// parseMessageConfigs parses a list of key=value strings into a map
func parseMessageConfigs(configs []string) (map[string]*string, error) {
	result := make(map[string]*string)

	for _, config := range configs {
		parts := strings.SplitN(config, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid config format: %s (expected key=value)", config)
		}

		key := parts[0]
		value := parts[1]
		result[key] = &value
	}

	return result, nil
}
