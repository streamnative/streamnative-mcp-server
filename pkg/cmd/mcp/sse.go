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
	"syscall"
	"time"

	stdlog "log"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/session"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

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

	// 4. Set SNCloud context if needed
	if configOpts.Options.KeyFile != "" {
		if configOpts.Options.PulsarInstance != "" && configOpts.Options.PulsarCluster != "" {
			err = mcp.SetContext(ctx, configOpts.Options, configOpts.Options.PulsarInstance, configOpts.Options.PulsarCluster)
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
	snConfig := configOpts.Options.LoadConfigOrDie()
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

	// Create go-sdk StreamableHTTPHandler
	// The getServer function returns the server instance for each request
	streamableHandler := mcpsdk.NewStreamableHTTPHandler(
		func(r *http.Request) *mcpsdk.Server {
			return mcpServer.Server
		},
		&mcpsdk.StreamableHTTPOptions{
			SessionTimeout: 5 * time.Minute,
		},
	)

	// Build the full path
	ssePath := configOpts.HTTPPath
	fmt.Fprintf(os.Stderr, "StreamNative Cloud MCP Server listening on http://%s%s\n",
		configOpts.HTTPAddr, ssePath)

	// Create HTTP server
	mux := http.NewServeMux()

	if pulsarSessionManager != nil {
		// Multi-session mode: wrap with auth middleware
		authMiddleware := func(next http.Handler) http.Handler {
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

				// Inject sessions into request context
				// Use server middleware to pass context to handlers
				next.ServeHTTP(w, r)
			})
		}

		// Add receiving middleware to inject sessions into MCP handler context
		mcpServer.Server.AddReceivingMiddleware(func(next mcpsdk.MethodHandler) mcpsdk.MethodHandler {
			return func(ctx context.Context, method string, req mcpsdk.Request) (mcpsdk.Result, error) {
				// Extract sessions from context (injected by HTTP middleware)
				// Note: go-sdk doesn't directly pass HTTP context to MCP handlers
				// We use the base sessions for non-multi-session, or per-request sessions
				return next(ctx, method, req)
			}
		})

		mux.Handle(ssePath, authMiddleware(streamableHandler))
	} else {
		// Non-multi-session mode: direct handler
		mux.Handle(ssePath, streamableHandler)
	}

	// Start HTTP server
	httpServer := &http.Server{
		Addr:              configOpts.HTTPAddr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second, // Prevent Slowloris attacks
	}

	errCh := make(chan error, 1)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	// Wait for shutdown signal
	select {
	case <-ctx.Done():
		fmt.Fprintln(os.Stderr, "Received shutdown signal, stopping server...")
	case err := <-errCh:
		return fmt.Errorf("sse server error: %w", err)
	}

	// Graceful shutdown
	shCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Stop Pulsar session manager first
	if pulsarSessionManager != nil {
		pulsarSessionManager.Stop()
	}

	// Shut down the HTTP server
	if err := httpServer.Shutdown(shCtx); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("Error shutting down HTTP server: %v", err)
		}
	}

	fmt.Fprintln(os.Stderr, "SSE server stopped gracefully")
	return nil
}
