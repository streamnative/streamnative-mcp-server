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

package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
	mcpctx "github.com/streamnative/streamnative-mcp-server/pkg/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/session"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

const (
	jsonrpcVersion        = "2.0"
	jsonrpcInvalidRequest = -32600
	jsonrpcInvalidParams  = -32602
)

// NewCmdMcpSseServer builds the SSE server command.
func NewCmdMcpSseServer(configOpts *ServerOptions) *cobra.Command {
	sseCmd := &cobra.Command{
		Use:   "sse",
		Short: "Start SSE server",
		Long:  `Start a server that communicates via HTTP with Server-Sent Events (SSE) for streaming MCP messages.`,
		Run: func(_ *cobra.Command, _ []string) {
			if err := runSseServer(configOpts); err != nil {
				fmt.Fprintf(os.Stderr, "failed to run SSE server: %v\n", err)
			}
		},
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return configOpts.Complete()
		},
	}

	// Add SSE server specific flags
	sseCmd.Flags().StringVar(&configOpts.HTTPAddr, "http-addr", ":9090", "HTTP server address")
	sseCmd.Flags().StringVar(&configOpts.HTTPPath, "http-path", "/mcp", "HTTP server path for SSE endpoint")

	return sseCmd
}

func runSseServer(configOpts *ServerOptions) error {
	// 1. Create a cancellable context that fires on SIGINT / SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 2. Initialize logger if log file specified
	logger, err := initLogger(configOpts.LogFile)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	// 3. Create a new MCP server
	ctx = context.WithValue(ctx, common.OptionsKey, configOpts.Options)
	mcpServer, err := newMcpServer(ctx, configOpts, logger)
	if err != nil {
		return fmt.Errorf("failed to create MCP server: %w", err)
	}

	// 4. Set the context
	ctx = mcpctx.WithSNCloudSession(ctx, mcpServer.SNCloudSession)
	ctx = mcpctx.WithPulsarSession(ctx, mcpServer.PulsarSession)
	ctx = mcpctx.WithKafkaSession(ctx, mcpServer.KafkaSession)
	if configOpts.KeyFile != "" {
		if configOpts.PulsarInstance != "" && configOpts.PulsarCluster != "" {
			err = mcpctx.SetContext(ctx, configOpts.Options, configOpts.PulsarInstance, configOpts.PulsarCluster)
			if err != nil {
				return errors.Wrap(err, "failed to set StreamNative Cloud context")
			}
		}
	}

	// add Pulsar Functions as MCP tools
	mcpServer.PulsarFunctionManagedMcpTools(configOpts.ReadOnly, configOpts.Features)

	// Create Pulsar session manager for multi-session support (only for external Pulsar mode)
	var pulsarSessionManager *session.PulsarSessionManager
	snConfig := configOpts.LoadConfigOrDie()
	if snConfig.ExternalPulsar != nil && configOpts.MultiSessionPulsar {
		managerConfig := &session.PulsarSessionManagerConfig{
			MaxSessions:     configOpts.SessionCacheSize,
			SessionTTL:      time.Duration(configOpts.SessionTTLMinutes) * time.Minute,
			CleanupInterval: 5 * time.Minute,
			BaseContext: pulsar.PulsarContext{
				ServiceURL:                    snConfig.ExternalPulsar.ServiceURL,
				WebServiceURL:                 snConfig.ExternalPulsar.WebServiceURL,
				TLSAllowInsecureConnection:    snConfig.ExternalPulsar.TLSAllowInsecureConnection,
				TLSEnableHostnameVerification: snConfig.ExternalPulsar.TLSEnableHostnameVerification,
				TLSTrustCertsFilePath:         snConfig.ExternalPulsar.TLSTrustCertsFilePath,
				TLSCertFile:                   snConfig.ExternalPulsar.TLSCertFile,
				TLSKeyFile:                    snConfig.ExternalPulsar.TLSKeyFile,
			},
		}
		// Pass nil as globalSession - in multi-session mode, every request must have a valid token
		pulsarSessionManager = session.NewPulsarSessionManager(managerConfig, nil, logger)
		logger.Info("Multi-session Pulsar mode enabled")
	}

	mux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:              configOpts.HTTPAddr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second, // Prevent Slowloris attacks
	}

	sessionContext := func(r *http.Request) context.Context {
		c := context.WithValue(r.Context(), common.OptionsKey, configOpts.Options)
		c = mcpctx.WithKafkaSession(c, mcpServer.KafkaSession)
		c = mcpctx.WithSNCloudSession(c, mcpServer.SNCloudSession)

		// Handle per-user Pulsar sessions
		if pulsarSessionManager != nil {
			token := session.ExtractBearerToken(r)
			// Token is already validated in auth middleware, this should always succeed
			pulsarSession, err := pulsarSessionManager.GetOrCreateSession(c, token)
			if err != nil {
				// Should not happen since middleware validates token first
				if logger != nil {
					logger.WithError(err).Error("Unexpected auth error after middleware validation")
				}
				// Don't set PulsarSession - tool handlers will fail gracefully with "session not found"
			} else {
				c = mcpctx.WithPulsarSession(c, pulsarSession)
				if token != "" {
					c = session.WithUserTokenHash(c, pulsarSessionManager.HashTokenForLog(token))
				}
			}
		} else {
			c = mcpctx.WithPulsarSession(c, mcpServer.PulsarSession)
		}

		return c
	}

	sessions := newSSESessionStore()

	// 4. Expose the full SSE URL to the user
	ssePath := joinHTTPPath(configOpts.HTTPPath, "sse")
	msgPath := joinHTTPPath(configOpts.HTTPPath, "message")
	fmt.Fprintf(os.Stderr, "StreamNative Cloud MCP Server listening on http://%s%s\n",
		configOpts.HTTPAddr, ssePath)

	healthPath := joinHTTPPath(configOpts.HTTPPath, "healthz")
	readyPath := joinHTTPPath(configOpts.HTTPPath, "readyz")

	authMiddleware := func(next http.Handler) http.Handler {
		return next
	}
	if pulsarSessionManager != nil {
		// Multi-session mode: validate token before processing SSE/message requests
		authMiddleware = func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				token := session.ExtractBearerToken(r)
				if token == "" {
					w.Header().Set("Content-Type", "application/json")
					http.Error(w, `{"error":"missing Authorization header"}`, http.StatusUnauthorized)
					return
				}
				// Pre-validate token by attempting to get/create session
				if _, err := pulsarSessionManager.GetOrCreateSession(r.Context(), token); err != nil {
					if logger != nil {
						logger.WithError(err).Warn("Authentication failed")
					}
					w.Header().Set("Content-Type", "application/json")
					http.Error(w, `{"error":"authentication failed"}`, http.StatusUnauthorized)
					return
				}
				next.ServeHTTP(w, r)
			})
		}
		if logger != nil {
			logger.Info("SSE server started with authentication middleware")
		}
	}

	sseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if _, ok := w.(http.Flusher); !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		sessionID := uuid.NewString()
		endpoint := fmt.Sprintf("%s?sessionid=%s", msgPath, sessionID)
		transport := &sdk.SSEServerTransport{
			Endpoint: endpoint,
			Response: w,
		}

		sessionTransport := &sseSessionTransport{transport: transport, sessionID: sessionID}
		serverSession, err := mcpServer.MCPServer.Connect(sessionContext(r), sessionTransport, nil)
		if err != nil {
			http.Error(w, "Failed to start SSE transport", http.StatusInternalServerError)
			return
		}

		state := &sseSessionState{sessionID: sessionID, transport: transport, session: serverSession}
		sessions.Store(sessionID, state)
		defer sessions.Delete(sessionID)
		defer func() {
			_ = serverSession.Close()
		}()

		waitCh := make(chan struct{})
		go func() {
			_ = serverSession.Wait()
			close(waitCh)
		}()

		select {
		case <-r.Context().Done():
		case <-waitCh:
		}
	})

	messageHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSONRPCError(w, jsonrpcInvalidRequest, "Method not allowed")
			return
		}

		sessionID := r.URL.Query().Get("sessionId")
		if sessionID == "" {
			sessionID = r.URL.Query().Get("sessionid")
		}
		if sessionID == "" {
			writeJSONRPCError(w, jsonrpcInvalidParams, "Missing sessionId")
			return
		}

		state, ok := sessions.Load(sessionID)
		if !ok || state.transport == nil {
			writeJSONRPCError(w, jsonrpcInvalidParams, "Invalid session ID")
			return
		}

		state.transport.ServeHTTP(w, r)
	})

	mux.Handle(ssePath, authMiddleware(sseHandler))
	mux.Handle(msgPath, authMiddleware(messageHandler))
	mux.HandleFunc(healthPath, healthHandler("ok"))
	mux.HandleFunc(readyPath, healthHandler("ready"))

	// 5. Run the HTTP listener in a goroutine
	errCh := make(chan error, 1)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err // bubble up real crashes
		}
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// 6. Block until Ctrl-C or an internal error
	select {
	case <-ctx.Done():
		// user hit Ctrl-C
		fmt.Fprintln(os.Stderr, "Received shutdown signal, stopping server...")
	case err := <-errCh:
		// HTTP server crashed
		return fmt.Errorf("sse server error: %w", err)
	}

	// 7. Graceful shutdown
	shCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Stop Pulsar session manager first
	if pulsarSessionManager != nil {
		pulsarSessionManager.Stop()
	}

	sessions.CloseAll()

	// Shut down the HTTP server (also closes the SSE connections)
	if err := httpServer.Shutdown(shCtx); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("Error shutting down SSE server: %v", err)
		}
	}

	// Wait for any remaining operations to complete
	select {
	case <-shCtx.Done():
		return fmt.Errorf("shutdown timed out")
	case <-time.After(100 * time.Millisecond):
		// Give a small grace period for cleanup
	}

	fmt.Fprintln(os.Stderr, "SSE server stopped gracefully")
	return nil
}

type sseSessionState struct {
	sessionID string
	transport *sdk.SSEServerTransport
	session   *sdk.ServerSession
}

type sseSessionStore struct {
	sessions sync.Map
}

func newSSESessionStore() *sseSessionStore {
	return &sseSessionStore{}
}

func (s *sseSessionStore) Store(id string, state *sseSessionState) {
	s.sessions.Store(id, state)
}

func (s *sseSessionStore) Load(id string) (*sseSessionState, bool) {
	value, ok := s.sessions.Load(id)
	if !ok {
		return nil, false
	}
	state, ok := value.(*sseSessionState)
	return state, ok
}

func (s *sseSessionStore) Delete(id string) {
	s.sessions.Delete(id)
}

func (s *sseSessionStore) CloseAll() {
	s.sessions.Range(func(key, value any) bool {
		if state, ok := value.(*sseSessionState); ok {
			if state.session != nil {
				_ = state.session.Close()
			}
		}
		s.sessions.Delete(key)
		return true
	})
}

type sseSessionTransport struct {
	transport *sdk.SSEServerTransport
	sessionID string
}

func (t *sseSessionTransport) Connect(ctx context.Context) (sdk.Connection, error) {
	conn, err := t.transport.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return &sseSessionConn{Connection: conn, sessionID: t.sessionID}, nil
}

type sseSessionConn struct {
	sdk.Connection
	sessionID string
}

func (c *sseSessionConn) SessionID() string {
	return c.sessionID
}

type jsonrpcErrorResponse struct {
	JSONRPC string             `json:"jsonrpc"`
	ID      any                `json:"id"`
	Error   jsonrpcErrorDetail `json:"error"`
}

type jsonrpcErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func writeJSONRPCError(w http.ResponseWriter, code int, message string) {
	response := jsonrpcErrorResponse{
		JSONRPC: jsonrpcVersion,
		ID:      nil,
		Error: jsonrpcErrorDetail{
			Code:    code,
			Message: message,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}

func joinHTTPPath(basePath string, suffix string) string {
	cleanedBase := path.Clean(basePath)
	if cleanedBase == "." || cleanedBase == "/" {
		cleanedBase = ""
	}
	if cleanedBase == "" {
		return "/" + suffix
	}
	return cleanedBase + "/" + suffix
}

func healthHandler(status string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(status))
	}
}
