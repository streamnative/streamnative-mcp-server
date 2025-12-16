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

	"github.com/mark3labs/mcp-go/server"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp"
	context2 "github.com/streamnative/streamnative-mcp-server/pkg/mcp"
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

	// 4. Set the context
	ctx = context2.WithSNCloudSession(ctx, mcpServer.SNCloudSession)
	ctx = context2.WithPulsarSession(ctx, mcpServer.PulsarSession)
	ctx = context2.WithKafkaSession(ctx, mcpServer.KafkaSession)
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
		pulsarSessionManager = session.NewPulsarSessionManager(managerConfig, mcpServer.PulsarSession, logger)
		logger.Info("Multi-session Pulsar mode enabled")
	}

	sseServer := server.NewSSEServer(
		mcpServer.MCPServer,
		server.WithStaticBasePath(configOpts.HTTPPath),
		server.WithSSEContextFunc(func(ctx context.Context, r *http.Request) context.Context {
			c := context.WithValue(ctx, common.OptionsKey, configOpts.Options)
			c = context2.WithKafkaSession(c, mcpServer.KafkaSession)
			c = context2.WithSNCloudSession(c, mcpServer.SNCloudSession)

			// Handle per-user Pulsar sessions
			if pulsarSessionManager != nil {
				token := session.ExtractBearerToken(r)
				if pulsarSession, err := pulsarSessionManager.GetOrCreateSession(ctx, token); err == nil {
					c = context2.WithPulsarSession(c, pulsarSession)
					if token != "" {
						c = session.WithUserTokenHash(c, pulsarSessionManager.HashTokenForLog(token))
					}
				} else {
					logger.WithError(err).Warn("Failed to get per-user Pulsar session, using global")
					c = context2.WithPulsarSession(c, mcpServer.PulsarSession)
				}
			} else {
				c = context2.WithPulsarSession(c, mcpServer.PulsarSession)
			}

			return c
		}),
	)

	// 4. Expose the full SSE URL to the user
	ssePath := sseServer.CompleteSsePath()
	fmt.Fprintf(os.Stderr, "StreamNative Cloud MCP Server listening on http://%s%s\n",
		configOpts.HTTPAddr, ssePath)

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

	// First try to shut down the SSE server
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
