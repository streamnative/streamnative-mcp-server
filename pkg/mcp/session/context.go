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

// Package session provides context helpers for MCP session management.
package session

import (
	"context"
)

type contextKey string

const (
	// PulsarSessionManagerKey is used to store the session manager in context
	PulsarSessionManagerKey contextKey = "pulsar_session_manager"
	// UserTokenHashKey stores the user's JWT token hash for debugging
	UserTokenHashKey contextKey = "user_token_hash"
)

// WithPulsarSessionManager adds the session manager to the context
func WithPulsarSessionManager(ctx context.Context, manager *PulsarSessionManager) context.Context {
	return context.WithValue(ctx, PulsarSessionManagerKey, manager)
}

// GetPulsarSessionManager retrieves the session manager from context
func GetPulsarSessionManager(ctx context.Context) *PulsarSessionManager {
	if val := ctx.Value(PulsarSessionManagerKey); val != nil {
		if manager, ok := val.(*PulsarSessionManager); ok {
			return manager
		}
	}
	return nil
}

// WithUserTokenHash adds the hashed user token to context
func WithUserTokenHash(ctx context.Context, hash string) context.Context {
	return context.WithValue(ctx, UserTokenHashKey, hash)
}

// GetUserTokenHash retrieves the hashed user token from context
func GetUserTokenHash(ctx context.Context) string {
	if val := ctx.Value(UserTokenHashKey); val != nil {
		if hash, ok := val.(string); ok {
			return hash
		}
	}
	return ""
}
