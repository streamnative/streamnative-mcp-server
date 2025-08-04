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
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/streamnative/streamnative-mcp-server/pkg/auth"
	"github.com/streamnative/streamnative-mcp-server/pkg/auth/store"
	sncloud "github.com/streamnative/streamnative-mcp-server/sdk/sdk-apiserver"
)

// Mock implementations for testing
type mockStore struct {
	mock.Mock
}

func (m *mockStore) LoadGrant(audience string) (*auth.AuthorizationGrant, error) {
	args := m.Called(audience)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.AuthorizationGrant), args.Error(1)
}

func (m *mockStore) SaveGrant(audience string, grant auth.AuthorizationGrant) error {
	args := m.Called(audience, grant)
	return args.Error(0)
}

func (m *mockStore) WhoAmI(audience string) (string, error) {
	args := m.Called(audience)
	return args.String(0), args.Error(1)
}

func (m *mockStore) Logout() error {
	args := m.Called()
	return args.Error(0)
}

type mockCachingTokenSource struct {
	mock.Mock
}

func (m *mockCachingTokenSource) Token() (*oauth2.Token, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func (m *mockCachingTokenSource) InvalidateToken() error {
	args := m.Called()
	return args.Error(0)
}

func TestSession_GetAPIClient_JWT_LazyInitialization(t *testing.T) {
	// Test lazy initialization with JWT token
	session := &Session{
		Ctx: SNCloudContext{
			APIURL:  "https://api.example.com",
			Timeout: 30 * time.Second,
		},
		useJWT:      true,
		TokenSource: NewJWTTokenSource("test-jwt-token"),
	}

	// First call should initialize the client
	client1, err := session.GetAPIClient()
	require.NoError(t, err)
	require.NotNil(t, client1)
	require.NotNil(t, session.APIClient)
	require.NotNil(t, session.Configuration)

	// Verify configuration
	assert.Equal(t, "api.example.com", session.Configuration.Host)
	assert.Equal(t, "https", session.Configuration.Scheme)
	assert.Equal(t, "StreamNative-MCP-Server/1.0.0", session.Configuration.UserAgent)

	// Second call should return the same client (cached)
	client2, err := session.GetAPIClient()
	require.NoError(t, err)
	assert.Same(t, client1, client2)
}

func TestSession_GetAPIClient_OAuth_LazyInitialization(t *testing.T) {
	// Create mock token refresher
	mockRefresher := &mockCachingTokenSource{}
	mockRefresher.On("Token").Return(&oauth2.Token{
		AccessToken: "test-access-token",
		TokenType:   "Bearer",
	}, nil)

	session := &Session{
		Ctx: SNCloudContext{
			APIURL:  "https://api.example.com",
			Timeout: 30 * time.Second,
		},
		useJWT:         false,
		TokenRefresher: &OAuth2TokenRefresher{source: mockRefresher},
	}

	// First call should initialize the client
	client1, err := session.GetAPIClient()
	require.NoError(t, err)
	require.NotNil(t, client1)
	require.NotNil(t, session.APIClient)

	// Second call should return the same client (cached)
	client2, err := session.GetAPIClient()
	require.NoError(t, err)
	assert.Same(t, client1, client2)
}

func TestSession_GetAPIClient_ErrorPaths(t *testing.T) {
	tests := []struct {
		name        string
		session     *Session
		expectError string
	}{
		{
			name: "OAuth without token refresher",
			session: &Session{
				Ctx: SNCloudContext{
					APIURL:  "https://api.example.com",
					Timeout: 30 * time.Second,
				},
				useJWT:         false,
				TokenRefresher: nil,
			},
			expectError: "token refresher not initialized",
		},
		{
			name: "Invalid API URL",
			session: &Session{
				Ctx: SNCloudContext{
					APIURL:  "://invalid-url",
					Timeout: 30 * time.Second,
				},
				useJWT:      true,
				TokenSource: NewJWTTokenSource("test-jwt-token"),
			},
			expectError: "failed to parse API URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := tt.session.GetAPIClient()
			assert.Nil(t, client)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectError)
		})
	}
}

func TestSession_GetAPIClient_ConcurrentAccess(t *testing.T) {
	// Test thread safety of lazy initialization
	session := &Session{
		Ctx: SNCloudContext{
			APIURL:  "https://api.example.com",
			Timeout: 30 * time.Second,
		},
		useJWT:      true,
		TokenSource: NewJWTTokenSource("test-jwt-token"),
	}

	const numGoroutines = 10
	clients := make([]*sncloud.APIClient, numGoroutines)
	errors := make([]error, numGoroutines)
	var wg sync.WaitGroup

	// Launch multiple goroutines to call GetAPIClient concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			clients[index], errors[index] = session.GetAPIClient()
		}(i)
	}

	wg.Wait()

	// All calls should succeed and return the same client instance
	for i := 0; i < numGoroutines; i++ {
		require.NoError(t, errors[i])
		require.NotNil(t, clients[i])
		assert.Same(t, clients[0], clients[i])
	}
}

func TestSession_GetLogClient_JWT_LazyInitialization(t *testing.T) {
	// Test lazy initialization with JWT token
	session := &Session{
		Ctx: SNCloudContext{
			LogAPIURL: "https://logs.example.com",
			Timeout:   30 * time.Second,
		},
		useJWT:      true,
		TokenSource: NewJWTTokenSource("test-jwt-token"),
	}

	// First call should initialize the client
	client1, err := session.GetLogClient()
	require.NoError(t, err)
	require.NotNil(t, client1)
	require.NotNil(t, session.LogClient)

	// Verify client configuration
	assert.Equal(t, 10*time.Second, client1.Timeout)

	// Second call should return the same client (cached)
	client2, err := session.GetLogClient()
	require.NoError(t, err)
	assert.Same(t, client1, client2)
}

func TestSession_GetLogClient_OAuth_LazyInitialization(t *testing.T) {
	// Create mock token refresher
	mockRefresher := &mockCachingTokenSource{}
	mockRefresher.On("Token").Return(&oauth2.Token{
		AccessToken: "test-access-token",
		TokenType:   "Bearer",
	}, nil)

	session := &Session{
		Ctx: SNCloudContext{
			LogAPIURL: "https://logs.example.com",
			Timeout:   30 * time.Second,
		},
		useJWT:         false,
		TokenRefresher: &OAuth2TokenRefresher{source: mockRefresher},
	}

	// First call should initialize the client
	client1, err := session.GetLogClient()
	require.NoError(t, err)
	require.NotNil(t, client1)

	// Second call should return the same client (cached)
	client2, err := session.GetLogClient()
	require.NoError(t, err)
	assert.Same(t, client1, client2)
}

func TestSession_GetLogClient_ErrorPaths(t *testing.T) {
	tests := []struct {
		name        string
		session     *Session
		expectError string
	}{
		{
			name: "OAuth without token refresher",
			session: &Session{
				Ctx: SNCloudContext{
					LogAPIURL: "https://logs.example.com",
					Timeout:   30 * time.Second,
				},
				useJWT:         false,
				TokenRefresher: nil,
			},
			expectError: "token refresher not initialized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := tt.session.GetLogClient()
			assert.Nil(t, client)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectError)
		})
	}
}

func TestSession_GetLogClient_ConcurrentAccess(t *testing.T) {
	// Test thread safety of lazy initialization
	session := &Session{
		Ctx: SNCloudContext{
			LogAPIURL: "https://logs.example.com",
			Timeout:   30 * time.Second,
		},
		useJWT:      true,
		TokenSource: NewJWTTokenSource("test-jwt-token"),
	}

	const numGoroutines = 10
	clients := make([]*http.Client, numGoroutines)
	errors := make([]error, numGoroutines)
	var wg sync.WaitGroup

	// Launch multiple goroutines to call GetLogClient concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			clients[idx], errors[idx] = session.GetLogClient()
		}(i)
	}

	wg.Wait()

	// All calls should succeed and return the same client instance
	for i := 0; i < numGoroutines; i++ {
		require.NoError(t, errors[i])
		require.NotNil(t, clients[i])
		assert.Same(t, clients[0], clients[i])
	}
}

func TestSession_TokenRefreshing_Scenarios(t *testing.T) {
	// Test token refreshing scenarios
	t.Run("successful token refresh", func(t *testing.T) {
		mockRefresher := &mockCachingTokenSource{}
		// Token may be called during initialization, but we don't need to guarantee it
		mockRefresher.On("Token").Return(&oauth2.Token{
			AccessToken: "fresh-access-token",
			TokenType:   "Bearer",
			Expiry:      time.Now().Add(time.Hour),
		}, nil).Maybe()

		session := &Session{
			Ctx: SNCloudContext{
				APIURL:  "https://api.example.com",
				Timeout: 30 * time.Second,
			},
			useJWT:         false,
			TokenRefresher: &OAuth2TokenRefresher{source: mockRefresher},
		}

		client, err := session.GetAPIClient()
		require.NoError(t, err)
		require.NotNil(t, client)

		// The token source is available and configured correctly
		assert.NotNil(t, session.TokenRefresher)
	})

	t.Run("token refresh failure", func(t *testing.T) {
		mockRefresher := &mockCachingTokenSource{}
		mockRefresher.On("Token").Return(nil, errors.New("token refresh failed")).Maybe()

		session := &Session{
			Ctx: SNCloudContext{
				APIURL:  "https://api.example.com",
				Timeout: 30 * time.Second,
			},
			useJWT:         false,
			TokenRefresher: &OAuth2TokenRefresher{source: mockRefresher},
		}

		// The client should still be created even if token refresh fails initially
		// because oauth2.Transport handles token errors during actual HTTP requests
		client, err := session.GetAPIClient()
		require.NoError(t, err)
		require.NotNil(t, client)
	})
}

func TestSession_InitializationOnlyOnce(t *testing.T) {
	// Test that initialization happens only once even with multiple calls
	session := &Session{
		Ctx: SNCloudContext{
			APIURL:  "https://api.example.com",
			Timeout: 30 * time.Second,
		},
		useJWT:      true,
		TokenSource: NewJWTTokenSource("test-jwt-token"),
	}

	// Multiple calls to GetAPIClient should only initialize once
	const numCalls = 5
	clients := make([]*sncloud.APIClient, numCalls)
	errors := make([]error, numCalls)

	for i := 0; i < numCalls; i++ {
		clients[i], errors[i] = session.GetAPIClient()
	}

	// All calls should succeed and return the same client instance
	for i := 0; i < numCalls; i++ {
		require.NoError(t, errors[i])
		require.NotNil(t, clients[i])
		assert.Same(t, clients[0], clients[i])
	}

	// Verify that configuration was set (indicating initialization occurred)
	assert.NotNil(t, session.Configuration)
	assert.Equal(t, "api.example.com", session.Configuration.Host)
}

func TestJWTTokenSource(t *testing.T) {
	// Test JWT token source implementation
	//nolint:gosec
	jwtToken := "test.jwt.token"
	source := NewJWTTokenSource(jwtToken)

	token, err := source.Token()
	require.NoError(t, err)
	require.NotNil(t, token)
	assert.Equal(t, jwtToken, token.AccessToken)
	assert.Equal(t, "Bearer", token.TokenType)
}

func TestNewSNCloudSession_JWT(t *testing.T) {
	// Test session creation with JWT token
	ctx := SNCloudContext{
		APIURL:    "https://api.example.com",
		LogAPIURL: "https://logs.example.com",
		//nolint:gosec
		JWTToken:    "test.jwt.token",
		Timeout:     30 * time.Second,
		Audience:    "test-audience",
		KeyFilePath: "/path/to/key.json",
	}

	session, err := NewSNCloudSession(ctx)
	require.NoError(t, err)
	require.NotNil(t, session)

	assert.True(t, session.useJWT)
	assert.NotNil(t, session.TokenSource)
	assert.Nil(t, session.TokenRefresher)

	// Test that JWT token source works
	token, err := session.TokenSource.Token()
	require.NoError(t, err)
	assert.Equal(t, "test.jwt.token", token.AccessToken)
}

func TestNewSNCloudSession_OAuth_InitializationError(t *testing.T) {
	// Test session creation with OAuth that fails during initialization
	mockStore := &mockStore{}
	mockStore.On("LoadGrant", "test-audience").Return(nil, store.ErrNoAuthenticationData)

	ctx := SNCloudContext{
		APIURL:      "https://api.example.com",
		LogAPIURL:   "https://logs.example.com",
		Timeout:     30 * time.Second,
		Audience:    "test-audience",
		KeyFilePath: "/invalid/path/to/key.json", // This will cause initialization to fail
		TokenStore:  mockStore,
	}

	session, err := NewSNCloudSession(ctx)
	assert.Nil(t, session)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to initialize token refresher")
}
