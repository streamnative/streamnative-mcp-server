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

package session

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

func TestExtractBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		authHeader    string
		expectedToken string
	}{
		{
			name:          "valid bearer token",
			authHeader:    "Bearer sample-token-123",
			expectedToken: "sample-token-123",
		},
		{
			name:          "empty header",
			authHeader:    "",
			expectedToken: "",
		},
		{
			name:          "no bearer prefix",
			authHeader:    "Basic dXNlcjpwYXNz",
			expectedToken: "",
		},
		{
			name:          "bearer lowercase",
			authHeader:    "bearer token123",
			expectedToken: "",
		},
		{
			name:          "only bearer prefix",
			authHeader:    "Bearer ",
			expectedToken: "",
		},
		{
			name:          "token with spaces",
			authHeader:    "Bearer token with spaces",
			expectedToken: "token with spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			token := ExtractBearerToken(req)
			if token != tt.expectedToken {
				t.Errorf("ExtractBearerToken() = %q, want %q", token, tt.expectedToken)
			}
		})
	}
}

func TestPulsarSessionManager_HashToken(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	manager := NewPulsarSessionManager(nil, nil, logger)
	defer manager.Stop()

	// Same token should produce same hash
	token := "test-token-123"
	hash1 := manager.hashToken(token)
	hash2 := manager.hashToken(token)

	if hash1 != hash2 {
		t.Errorf("Same token produced different hashes: %s vs %s", hash1, hash2)
	}

	// Different tokens should produce different hashes
	hash3 := manager.hashToken("different-token")
	if hash1 == hash3 {
		t.Errorf("Different tokens produced same hash")
	}

	// Hash should be 64 characters (SHA256 hex)
	if len(hash1) != 64 {
		t.Errorf("Hash length = %d, want 64", len(hash1))
	}
}

func TestPulsarSessionManager_HashTokenForLog(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	manager := NewPulsarSessionManager(nil, nil, logger)
	defer manager.Stop()

	token := "test-token-123"
	shortHash := manager.HashTokenForLog(token)

	// Short hash should be 8 characters
	if len(shortHash) != 8 {
		t.Errorf("HashTokenForLog length = %d, want 8", len(shortHash))
	}

	// Should be prefix of full hash
	fullHash := manager.hashToken(token)
	if shortHash != fullHash[:8] {
		t.Errorf("HashTokenForLog = %s, should be prefix of %s", shortHash, fullHash)
	}
}

func TestPulsarSessionManager_EmptyTokenReturnsGlobalSession(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	globalSession := &pulsar.Session{}
	manager := NewPulsarSessionManager(nil, globalSession, logger)
	defer manager.Stop()

	// Empty token should return global session
	session, err := manager.GetOrCreateSession(context.TODO(), "")
	if err != nil {
		t.Errorf("GetOrCreateSession with empty token returned error: %v", err)
	}
	if session != globalSession {
		t.Errorf("GetOrCreateSession with empty token did not return global session")
	}

	// Session count should be 0 (no cached sessions)
	if count := manager.SessionCount(); count != 0 {
		t.Errorf("SessionCount() = %d, want 0", count)
	}
}

func TestPulsarSessionManager_SessionCount(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	manager := NewPulsarSessionManager(nil, nil, logger)
	defer manager.Stop()

	// Initially empty
	if count := manager.SessionCount(); count != 0 {
		t.Errorf("Initial SessionCount() = %d, want 0", count)
	}
}

func TestDefaultPulsarSessionManagerConfig(t *testing.T) {
	config := DefaultPulsarSessionManagerConfig()

	if config.MaxSessions != 100 {
		t.Errorf("MaxSessions = %d, want 100", config.MaxSessions)
	}

	if config.SessionTTL.Minutes() != 30 {
		t.Errorf("SessionTTL = %v, want 30 minutes", config.SessionTTL)
	}

	if config.CleanupInterval.Minutes() != 5 {
		t.Errorf("CleanupInterval = %v, want 5 minutes", config.CleanupInterval)
	}
}

func TestContextHelpers(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	manager := NewPulsarSessionManager(nil, nil, logger)
	defer manager.Stop()

	// Test WithPulsarSessionManager and GetPulsarSessionManager
	ctx := WithPulsarSessionManager(context.Background(), manager)
	retrievedManager := GetPulsarSessionManager(ctx)
	if retrievedManager != manager {
		t.Error("GetPulsarSessionManager did not return the same manager")
	}

	// Test context without manager
	emptyCtx := context.Background()
	if GetPulsarSessionManager(emptyCtx) != nil {
		t.Error("GetPulsarSessionManager on empty context should return nil")
	}

	// Test WithUserTokenHash and GetUserTokenHash
	ctx = WithUserTokenHash(context.Background(), "abc12345")
	hash := GetUserTokenHash(ctx)
	if hash != "abc12345" {
		t.Errorf("GetUserTokenHash() = %s, want abc12345", hash)
	}

	// Test context without hash
	if GetUserTokenHash(emptyCtx) != "" {
		t.Error("GetUserTokenHash on empty context should return empty string")
	}
}
