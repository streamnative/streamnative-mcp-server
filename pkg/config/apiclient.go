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

package config

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"k8s.io/utils/clock"

	"github.com/streamnative/streamnative-mcp-server/pkg/auth"
	"github.com/streamnative/streamnative-mcp-server/pkg/auth/cache"
	"github.com/streamnative/streamnative-mcp-server/pkg/auth/store"
	sncloud "github.com/streamnative/streamnative-mcp-server/sdk/sdk-apiserver"
)

var SNCloudClientConfiguration *sncloud.Configuration
var SNCloudClient *sncloud.APIClient
var SNCloudLogClient *http.Client

// OAuth2TokenRefresher implements oauth2.TokenSource interface for refreshing OAuth2 tokens
// This is now a wrapper around the cache.CachingTokenSource to leverage the existing token caching
type OAuth2TokenRefresher struct {
	source cache.CachingTokenSource
}

// NewOAuth2TokenRefresher creates a new token refresher that uses the stored token cache
func NewOAuth2TokenRefresher(tokenStore store.Store, audience string, refresher auth.AuthorizationGrantRefresher) (*OAuth2TokenRefresher, error) {
	// Create a token cache that will automatically use the store for persistence
	source, err := cache.NewDefaultTokenCache(tokenStore, audience, refresher)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create token cache")
	}

	return &OAuth2TokenRefresher{
		source: source,
	}, nil
}

// Token implements the oauth2.TokenSource interface, leveraging the cached token
func (t *OAuth2TokenRefresher) Token() (*oauth2.Token, error) {
	// The source already handles caching logic, token validation, and refreshing
	return t.source.Token()
}

// InitSNCloudClient initializes the StreamNative Cloud API client
// Parameters:
//   - issuerURL: OAuth2 authorization server URL
//   - audience: API service audience identifier
//   - keyFilePath: Client credentials key file path
//   - apiURL: API server URL
//   - timeout: HTTP client timeout
//   - tokenStore: Store for caching tokens
func InitSNCloudClient(issuerURL, audience, keyFilePath, apiURL string, timeout time.Duration, tokenStore store.Store) error {
	// 1. Create Issuer configuration
	issuerData := auth.Issuer{
		IssuerEndpoint: issuerURL,
		Audience:       audience,
	}

	// 2. Check if we have an existing grant in the store
	grant, err := tokenStore.LoadGrant(audience)
	if err != nil && err != store.ErrNoAuthenticationData {
		return errors.Wrap(err, "failed to load grant from store")
	}

	// 3. If no grant exists or there was an error, create a new one
	if err == store.ErrNoAuthenticationData || grant == nil {
		// Create OAuth2 client credentials flow
		flow, err := auth.NewDefaultClientCredentialsFlow(issuerData, keyFilePath)
		if err != nil {
			return errors.Wrap(err, "failed to create client credentials flow")
		}

		// Get initial authorization
		grant, err = flow.Authorize()
		if err != nil {
			return errors.Wrap(err, "failed to authorize client")
		}

		// Save the grant to the store
		err = tokenStore.SaveGrant(audience, *grant)
		if err != nil {
			return errors.Wrap(err, "failed to save grant to store")
		}
	}

	// 4. Create token refresher
	refresher, err := auth.NewDefaultClientCredentialsGrantRefresher(issuerData, clock.RealClock{})
	if err != nil {
		return errors.Wrap(err, "failed to create token refresher")
	}

	// 5. Create token source with caching
	tokenRefresher, err := NewOAuth2TokenRefresher(tokenStore, audience, refresher)
	if err != nil {
		return errors.Wrap(err, "failed to create token refresher")
	}
	tokenSource := oauth2.ReuseTokenSource(nil, tokenRefresher)

	// 6. Create HTTP client with OAuth2 Transport
	httpClient := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
			Base:   http.DefaultTransport,
		},
		Timeout: timeout,
	}

	// 7. Create API client configuration
	parsedURL, err := url.Parse(apiURL)
	if err != nil {
		return errors.Wrap(err, "failed to parse API URL")
	}
	config := sncloud.NewConfiguration()
	config.Host = parsedURL.Host
	config.Scheme = parsedURL.Scheme
	config.HTTPClient = httpClient
	config.UserAgent = "StreamNative-MCP-Server/1.0.0"

	// 8. Create API client
	SNCloudClientConfiguration = config
	SNCloudClient = sncloud.NewAPIClient(config)

	return nil
}

// GetAPIClient returns the initialized API client or an error if not initialized
func GetAPIClient() (*sncloud.APIClient, error) {
	if SNCloudClient == nil {
		return nil, errors.New("API client not initialized, call InitSNCloudClient first")
	}
	return SNCloudClient, nil
}

// CreateContextWithTokenSource creates a context with the TokenSource
// This might be useful in special scenarios, but is usually not needed as we've already configured the default HTTP client
func CreateContextWithTokenSource(ctx context.Context) (context.Context, error) {
	client, err := GetAPIClient()
	if err != nil {
		return nil, err
	}

	// Extract TokenSource from HTTP client
	transport, ok := client.GetConfig().HTTPClient.Transport.(*oauth2.Transport)
	if !ok {
		return nil, errors.New("HTTP client transport is not an OAuth2 transport")
	}

	// Create context with TokenSource
	return context.WithValue(ctx, sncloud.ContextOAuth2, transport.Source), nil
}

// TokenRefreshed is called when a token is refreshed to persist the updated token
func TokenRefreshed(audience string, token *oauth2.Token, tokenStore store.Store) error {
	// Load existing grant
	grant, err := tokenStore.LoadGrant(audience)
	if err != nil {
		return errors.Wrap(err, "failed to load grant for token update")
	}

	// Update the token
	grant.Token = token

	// Save back to store
	return tokenStore.SaveGrant(audience, *grant)
}

func InitSNCloudLogClient(issuerData auth.Issuer, tokenStore store.Store) error {
	refresher, err := auth.NewDefaultClientCredentialsGrantRefresher(issuerData, clock.RealClock{})
	if err != nil {
		return errors.Wrap(err, "failed to create token refresher")
	}

	tokenRefresher, err := NewOAuth2TokenRefresher(tokenStore, issuerData.Audience, refresher)
	if err != nil {
		return errors.Wrap(err, "failed to create token refresher")
	}

	tokenSource := oauth2.ReuseTokenSource(nil, tokenRefresher)
	SNCloudLogClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &oauth2.Transport{
			Source: tokenSource,
			Base:   http.DefaultTransport,
		},
	}

	return nil
}

func ResetSNCloudLogClient() {
	SNCloudLogClient = nil
}

func GetSNCloudLogClient() (*http.Client, error) {
	if SNCloudLogClient == nil {
		return nil, errors.New("log tools are for StreamNative Cloud only, please check your context")
	}
	return SNCloudLogClient, nil
}
