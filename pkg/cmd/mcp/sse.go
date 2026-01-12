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
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	stdlog "log"

	"github.com/mark3labs/mcp-go/server"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
	mcpctx "github.com/streamnative/streamnative-mcp-server/pkg/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/session"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
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
		stdlog.Fatal("Failed to initialize logger:", err)
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
	// SSE is not support session-based tools, so we pass an fixed sessionId
	mcpServer.PulsarFunctionManagedMcpTools(configOpts.ReadOnly, configOpts.Features, "FIXED_SESSION_ID")

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
	sseContextFunc := func(ctx context.Context, r *http.Request) context.Context {
		c := context.WithValue(ctx, common.OptionsKey, configOpts.Options)
		c = mcpctx.WithKafkaSession(c, mcpServer.KafkaSession)
		c = mcpctx.WithSNCloudSession(c, mcpServer.SNCloudSession)

		// Handle per-user Pulsar sessions
		if pulsarSessionManager != nil {
			token := session.ExtractBearerToken(r)
			// Token is already validated in auth middleware, this should always succeed
			if pulsarSession, err := pulsarSessionManager.GetOrCreateSession(ctx, token); err == nil {
				c = mcpctx.WithPulsarSession(c, pulsarSession)
				if token != "" {
					c = session.WithUserTokenHash(c, pulsarSessionManager.HashTokenForLog(token))
				}
			} else {
				// Should not happen since middleware validates token first
				logger.WithError(err).Error("Unexpected auth error after middleware validation")
				// Don't set PulsarSession - tool handlers will fail gracefully with "session not found"
			}
		} else {
			c = mcpctx.WithPulsarSession(c, mcpServer.PulsarSession)
		}

		return c
	}
	sseServer := server.NewSSEServer(
		mcpServer.MCPServer,
		server.WithHTTPServer(httpServer),
		server.WithStaticBasePath(configOpts.HTTPPath),
		server.WithSSEContextFunc(sseContextFunc),
	)

	// 4. Expose the full SSE URL to the user
	ssePath := sseServer.CompleteSsePath()
	msgPath := sseServer.CompleteMessagePath()
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
					logger.WithError(err).Warn("Authentication failed")
					w.Header().Set("Content-Type", "application/json")
					http.Error(w, `{"error":"authentication failed"}`, http.StatusUnauthorized)
					return
				}
				next.ServeHTTP(w, r)
			})
		}
		logger.Info("SSE server started with authentication middleware")
	}

	mux.Handle(ssePath, authMiddleware(sseServer.SSEHandler()))
	mux.Handle(msgPath, authMiddleware(sseServer.MessageHandler()))
	mux.HandleFunc(healthPath, healthHandler("ok"))
	mux.HandleFunc(readyPath, healthHandler("ready"))

	// 5. Run the HTTP listener in a goroutine
	errCh := make(chan error, 1)
	go func() {
		if err := sseServer.Start(configOpts.HTTPAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
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

	// Shut down the SSE server (also closes the underlying HTTP server)
	if err := sseServer.Shutdown(shCtx); err != nil {
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

func joinHTTPPath(basePath string, suffix string) string {
	joined := path.Join(basePath, suffix)
	if joined == "" {
		return "/" + suffix
	}
	if joined[0] != '/' {
		return "/" + joined
	}
	return joined
}

func healthHandler(status string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(status))
	}
}
