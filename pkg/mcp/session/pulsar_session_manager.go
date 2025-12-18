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

package session

import (
	"container/list"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// PulsarSessionManagerConfig holds configuration for the session manager
type PulsarSessionManagerConfig struct {
	// MaxSessions is the maximum number of sessions to keep in the cache
	MaxSessions int
	// SessionTTL is how long a session can remain idle before eviction
	SessionTTL time.Duration
	// CleanupInterval is how often to run the cleanup goroutine
	CleanupInterval time.Duration
	// BaseContext provides default Pulsar connection parameters (without token)
	BaseContext pulsar.PulsarContext
}

// DefaultPulsarSessionManagerConfig returns sensible defaults
func DefaultPulsarSessionManagerConfig() *PulsarSessionManagerConfig {
	return &PulsarSessionManagerConfig{
		MaxSessions:     100,
		SessionTTL:      30 * time.Minute,
		CleanupInterval: 5 * time.Minute,
	}
}

// CachedSession wraps a Pulsar session with metadata
type CachedSession struct {
	Session    *pulsar.Session
	TokenHash  string
	LastAccess time.Time
	CreatedAt  time.Time
	element    *list.Element // for LRU tracking
}

// PulsarSessionManager manages per-user Pulsar sessions with LRU eviction
type PulsarSessionManager struct {
	config        *PulsarSessionManagerConfig
	sessions      map[string]*CachedSession // tokenHash -> session
	lruList       *list.List                // LRU ordering
	mutex         sync.RWMutex
	stopCh        chan struct{}
	logger        *logrus.Logger
	globalSession *pulsar.Session // fallback session when no token provided
}

// NewPulsarSessionManager creates a new session manager
func NewPulsarSessionManager(
	config *PulsarSessionManagerConfig,
	globalSession *pulsar.Session,
	logger *logrus.Logger,
) *PulsarSessionManager {
	if config == nil {
		config = DefaultPulsarSessionManagerConfig()
	}

	m := &PulsarSessionManager{
		config:        config,
		sessions:      make(map[string]*CachedSession),
		lruList:       list.New(),
		stopCh:        make(chan struct{}),
		logger:        logger,
		globalSession: globalSession,
	}

	// Start background cleanup
	go m.cleanupLoop()

	return m
}

// GetOrCreateSession retrieves or creates a Pulsar session for the given token
func (m *PulsarSessionManager) GetOrCreateSession(_ context.Context, token string) (*pulsar.Session, error) {
	if token == "" {
		// Return global session when no token provided
		// If no global session exists (multi-session mode), return error
		if m.globalSession == nil {
			return nil, fmt.Errorf("authentication required: missing Authorization header")
		}
		return m.globalSession, nil
	}

	tokenHash := m.hashToken(token)

	// Use write lock for all operations to avoid stale entry issues during lock upgrade
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if session exists
	if cached, exists := m.sessions[tokenHash]; exists {
		cached.LastAccess = time.Now()
		m.lruList.MoveToFront(cached.element)
		return cached.Session, nil
	}

	// Evict if at capacity
	if len(m.sessions) >= m.config.MaxSessions {
		m.evictOldest()
	}

	// Create new Pulsar session with token
	pulsarCtx := m.config.BaseContext
	pulsarCtx.Token = token
	// Clear AuthPlugin/AuthParams to use token-based auth
	pulsarCtx.AuthPlugin = ""
	pulsarCtx.AuthParams = ""

	session, err := pulsar.NewSession(pulsarCtx)
	if err != nil {
		m.logger.WithError(err).WithField("tokenHash", tokenHash[:8]).Error("Failed to create Pulsar session for token")
		return nil, fmt.Errorf("failed to create Pulsar session: %w", err)
	}

	// Cache the session
	now := time.Now()
	cachedSession := &CachedSession{
		Session:    session,
		TokenHash:  tokenHash,
		LastAccess: now,
		CreatedAt:  now,
	}
	cachedSession.element = m.lruList.PushFront(cachedSession)
	m.sessions[tokenHash] = cachedSession

	m.logger.WithField("tokenHash", tokenHash[:8]).Info("Created new Pulsar session")

	return session, nil
}

// hashToken creates a SHA256 hash of the token for safe storage/comparison
func (m *PulsarSessionManager) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// HashTokenForLog returns a short hash prefix for logging
func (m *PulsarSessionManager) HashTokenForLog(token string) string {
	return m.hashToken(token)[:8]
}

// evictOldest removes the least recently used session (caller must hold write lock)
func (m *PulsarSessionManager) evictOldest() {
	oldest := m.lruList.Back()
	if oldest == nil {
		return
	}

	cached := oldest.Value.(*CachedSession)
	m.lruList.Remove(oldest)
	delete(m.sessions, cached.TokenHash)

	// Close the session's clients
	m.closeSession(cached.Session)

	m.logger.WithField("tokenHash", cached.TokenHash[:8]).Info("Evicted oldest Pulsar session")
}

// closeSession safely closes a Pulsar session's clients
func (m *PulsarSessionManager) closeSession(session *pulsar.Session) {
	if session == nil {
		return
	}
	if session.Client != nil {
		session.Client.Close()
	}
}

// cleanupLoop periodically removes expired sessions
func (m *PulsarSessionManager) cleanupLoop() {
	ticker := time.NewTicker(m.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.cleanupExpired()
		case <-m.stopCh:
			return
		}
	}
}

// cleanupExpired removes sessions that have exceeded TTL
func (m *PulsarSessionManager) cleanupExpired() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	var toRemove []*CachedSession

	for _, cached := range m.sessions {
		if now.Sub(cached.LastAccess) > m.config.SessionTTL {
			toRemove = append(toRemove, cached)
		}
	}

	for _, cached := range toRemove {
		m.lruList.Remove(cached.element)
		delete(m.sessions, cached.TokenHash)
		m.closeSession(cached.Session)
		m.logger.WithField("tokenHash", cached.TokenHash[:8]).Info("Cleaned up expired Pulsar session")
	}
}

// Stop stops the session manager and cleans up all sessions
func (m *PulsarSessionManager) Stop() {
	close(m.stopCh)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, cached := range m.sessions {
		m.closeSession(cached.Session)
	}
	m.sessions = make(map[string]*CachedSession)
	m.lruList.Init()

	m.logger.Info("Pulsar session manager stopped")
}

// SessionCount returns the current number of cached sessions
func (m *PulsarSessionManager) SessionCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.sessions)
}

// ExtractBearerToken extracts JWT token from Authorization header
func ExtractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Expected format: "Bearer <token>"
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return ""
	}

	return strings.TrimPrefix(authHeader, bearerPrefix)
}
