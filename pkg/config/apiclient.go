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

// Package config provides configuration loading and client setup helpers.
package config

import (
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"k8s.io/utils/clock"

	"github.com/streamnative/streamnative-mcp-server/pkg/auth"
	"github.com/streamnative/streamnative-mcp-server/pkg/auth/cache"
	"github.com/streamnative/streamnative-mcp-server/pkg/auth/store"
	sncloud "github.com/streamnative/streamnative-mcp-server/sdk/sdk-apiserver"
)

// SNCloudContext represents the configuration context for StreamNative Cloud session
type SNCloudContext struct {
	IssuerURL      string
	Audience       string
	KeyFilePath    string
	JWTToken       string
	APIURL         string
	LogAPIURL      string
	Timeout        time.Duration
	Organization   string
	TokenStore     store.Store
	PulsarInstance string
	PulsarCluster  string
}

// Session represents a StreamNative Cloud session with managed clients
type Session struct {
	Ctx            SNCloudContext
	APIClient      *sncloud.APIClient
	LogClient      *http.Client
	TokenRefresher *OAuth2TokenRefresher
	TokenSource    oauth2.TokenSource
	Configuration  *sncloud.Configuration
	mutex          sync.RWMutex
	apiClientOnce  sync.Once
	logClientOnce  sync.Once
	useJWT         bool
}

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

// JWTTokenSource implements oauth2.TokenSource interface for static JWT tokens
type JWTTokenSource struct {
	token *oauth2.Token
}

// NewJWTTokenSource creates a new token source for static JWT tokens
func NewJWTTokenSource(jwtToken string) *JWTTokenSource {
	return &JWTTokenSource{
		token: &oauth2.Token{
			AccessToken: jwtToken,
			TokenType:   "Bearer",
		},
	}
}

// Token implements the oauth2.TokenSource interface for static JWT tokens
func (j *JWTTokenSource) Token() (*oauth2.Token, error) {
	return j.token, nil
}

// NewSNCloudSession creates a new StreamNative Cloud session with the provided context
func NewSNCloudSession(ctx SNCloudContext) (*Session, error) {
	session := &Session{
		Ctx: ctx,
	}

	// Check if JWT token is provided
	if ctx.JWTToken != "" {
		// Use JWT token directly without refresh mechanism
		session.useJWT = true
		session.TokenSource = NewJWTTokenSource(ctx.JWTToken)
	} else {
		// Initialize the session by setting up the token refresher for OAuth flow
		if err := session.initializeTokenRefresher(); err != nil {
			return nil, errors.Wrap(err, "failed to initialize token refresher")
		}
	}

	return session, nil
}

// NewSNCloudSessionFromOptions creates a new StreamNative Cloud session from configuration options
func NewSNCloudSessionFromOptions(options *Options) (*Session, error) {
	if options == nil {
		return nil, errors.New("options cannot be nil")
	}

	// Create SNCloudContext from options
	ctx := SNCloudContext{
		IssuerURL:      options.IssuerEndpoint,
		Audience:       options.Audience,
		KeyFilePath:    options.KeyFile,
		APIURL:         options.Server,
		LogAPIURL:      options.LogLocation,
		Timeout:        30 * time.Second,
		Organization:   options.Organization,
		TokenStore:     options.Store,
		PulsarInstance: options.PulsarInstance,
		PulsarCluster:  options.PulsarCluster,
	}

	// Create session
	session, err := NewSNCloudSession(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create session from options")
	}

	return session, nil
}

// SetPulsarClusterContext updates the StreamNative Cloud cluster binding for the session.
func (s *Session) SetPulsarClusterContext(instance, cluster string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.Ctx.PulsarInstance = instance
	s.Ctx.PulsarCluster = cluster
}

// ResetPulsarClusterContext clears the StreamNative Cloud cluster binding for the session.
func (s *Session) ResetPulsarClusterContext() {
	s.SetPulsarClusterContext("", "")
}

// GetPulsarClusterContext returns the current StreamNative Cloud cluster binding.
func (s *Session) GetPulsarClusterContext() (string, string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.Ctx.PulsarInstance, s.Ctx.PulsarCluster
}

// initializeTokenRefresher initializes the token refresher for the session
func (s *Session) initializeTokenRefresher() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Create Issuer configuration
	issuerData := auth.Issuer{
		IssuerEndpoint: s.Ctx.IssuerURL,
		Audience:       s.Ctx.Audience,
	}

	// Check if we have an existing grant in the store
	grant, err := s.Ctx.TokenStore.LoadGrant(s.Ctx.Audience)
	if err != nil && err != store.ErrNoAuthenticationData {
		return errors.Wrap(err, "failed to load grant from store")
	}

	// If no grant exists or there was an error, create a new one
	if err == store.ErrNoAuthenticationData || grant == nil {
		// Create OAuth2 client credentials flow
		flow, err := auth.NewDefaultClientCredentialsFlow(issuerData, s.Ctx.KeyFilePath)
		if err != nil {
			return errors.Wrap(err, "failed to create client credentials flow")
		}

		// Get initial authorization
		grant, err = flow.Authorize()
		if err != nil {
			return errors.Wrap(err, "failed to authorize client")
		}

		// Save the grant to the store
		err = s.Ctx.TokenStore.SaveGrant(s.Ctx.Audience, *grant)
		if err != nil {
			return errors.Wrap(err, "failed to save grant to store")
		}
	}

	// Create token refresher
	refresher, err := auth.NewDefaultClientCredentialsGrantRefresher(issuerData, clock.RealClock{})
	if err != nil {
		return errors.Wrap(err, "failed to create token refresher")
	}

	// Create token source with caching
	tokenRefresher, err := NewOAuth2TokenRefresher(s.Ctx.TokenStore, s.Ctx.Audience, refresher)
	if err != nil {
		return errors.Wrap(err, "failed to create token refresher")
	}

	s.TokenRefresher = tokenRefresher

	return nil
}

// GetAPIClient returns the API client for the session, initializing it if necessary
func (s *Session) GetAPIClient() (*sncloud.APIClient, error) {
	var err error
	s.apiClientOnce.Do(func() {
		err = s.initializeAPIClient()
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize API client")
	}

	return s.APIClient, nil
}

// initializeAPIClient initializes the API client for the session
func (s *Session) initializeAPIClient() error {
	var tokenSource oauth2.TokenSource

	if s.useJWT {
		// Use JWT token directly
		tokenSource = s.TokenSource
	} else {
		// Use OAuth token with refresh
		if s.TokenRefresher == nil {
			return errors.New("token refresher not initialized")
		}
		tokenSource = oauth2.ReuseTokenSource(nil, s.TokenRefresher)
	}

	// Create HTTP client with OAuth2 Transport
	httpClient := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
			Base:   http.DefaultTransport,
		},
		Timeout: s.Ctx.Timeout,
	}

	// Create API client configuration
	parsedURL, err := url.Parse(s.Ctx.APIURL)
	if err != nil {
		return errors.Wrap(err, "failed to parse API URL")
	}

	config := sncloud.NewConfiguration()
	config.Host = parsedURL.Host
	config.Scheme = parsedURL.Scheme
	config.HTTPClient = httpClient
	config.UserAgent = "StreamNative-MCP-Server/1.0.0"

	// Create API client
	s.Configuration = config
	s.APIClient = sncloud.NewAPIClient(config)

	return nil
}

// GetLogClient returns the log client for the session, initializing it if necessary
func (s *Session) GetLogClient() (*http.Client, error) {
	var err error
	s.logClientOnce.Do(func() {
		err = s.initializeLogClient()
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize log client")
	}

	return s.LogClient, nil
}

// initializeLogClient initializes the log client for the session
func (s *Session) initializeLogClient() error {
	var tokenSource oauth2.TokenSource

	if s.useJWT {
		// Use JWT token directly
		tokenSource = s.TokenSource
	} else {
		// Use OAuth token with refresh
		if s.TokenRefresher == nil {
			return errors.New("token refresher not initialized")
		}
		tokenSource = oauth2.ReuseTokenSource(nil, s.TokenRefresher)
	}

	// Create HTTP client with OAuth2 Transport
	s.LogClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &oauth2.Transport{
			Source: tokenSource,
			Base:   http.DefaultTransport,
		},
	}

	return nil
}

// Close closes the session and cleans up resources
func (s *Session) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Clear references to clients
	// Note: HTTP clients don't need explicit closing as they use connection pooling
	s.APIClient = nil
	s.LogClient = nil
	s.TokenRefresher = nil
	s.Configuration = nil

	return nil
}
