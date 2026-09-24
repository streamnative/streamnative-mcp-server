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

package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	stdlog "log"

	protocol "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
	"github.com/streamnative/streamnative-mcp-server/pkg/log"
	mcpctx "github.com/streamnative/streamnative-mcp-server/pkg/mcp"
)

// NewCmdMcpStdioServer builds the stdio server command.
func NewCmdMcpStdioServer(configOpts *ServerOptions) *cobra.Command {
	stdioCmd := &cobra.Command{
		Use:   "stdio",
		Short: "Start stdio server",
		Long:  `Start a server that communicates via standard input/output streams using JSON-RPC messages.`,
		Run: func(_ *cobra.Command, _ []string) {
			if err := runStdioServer(configOpts); err != nil {
				fmt.Fprintf(os.Stderr, "failed to run stdio server: %v\n", err)
			}
		},
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return configOpts.Complete()
		},
	}

	return stdioCmd
}

func runStdioServer(configOpts *ServerOptions) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize logger if log file specified
	logger, err := initLogger(configOpts.LogFile)
	if err != nil {
		stdlog.Fatal("Failed to initialize logger:", err)
	}

	// Create a new MCP server
	ctx = context.WithValue(ctx, common.OptionsKey, configOpts.Options)
	stdLogger := stdlog.New(logger.Writer(), "snmcp-server", 0)
	mcpServer, err := newMcpServer(ctx, configOpts, logger, stdioServerOptions(configOpts)...)
	if err != nil {
		return fmt.Errorf("failed to create MCP server: %w", err)
	}
	defer mcpServer.Close()

	ctx = mcpctx.WithSNCloudSession(ctx, mcpServer.SNCloudSession)
	ctx = mcpctx.WithPulsarSession(ctx, mcpServer.PulsarSession)
	ctx = mcpctx.WithKafkaSession(ctx, mcpServer.KafkaSession)
	if configOpts.IsCloudConfigured() && configOpts.PulsarInstance != "" && configOpts.PulsarCluster != "" {
		if err := mcpctx.SetContext(ctx, configOpts.Options, configOpts.PulsarInstance, configOpts.PulsarCluster); err != nil {
			return fmt.Errorf("failed to set StreamNative Cloud context: %w", err)
		}
	}

	stdioServer := server.NewStdioServer(mcpServer.MCPServer)
	stdioServer.SetErrorLogger(stdLogger)

	// Start listening for messages
	errC := make(chan error, 1)
	go func() {
		in, out := io.Reader(os.Stdin), io.Writer(os.Stdout)

		if configOpts.LogCommands {
			// If command logging is enabled, wrap the IO with a logger
			loggedIO := log.NewIOLogger(in, out, logger)
			in, out = loggedIO, loggedIO
		}

		errC <- stdioServer.Listen(ctx, in, out)
	}()

	_, _ = fmt.Fprintf(os.Stderr, "StreamNative Cloud MCP Server running on stdio\n")

	// Wait for shutdown signal
	select {
	case <-ctx.Done():
		fmt.Fprintf(os.Stderr, "shutting down server...\n")
		if logger != nil {
			logger.Info("Shutting down server...")
		}
	case err := <-errC:
		if err != nil {
			if logger != nil {
				logger.Errorf("Error running server: %v", err)
			}
			return fmt.Errorf("error running server: %w", err)
		}
	}

	return nil
}

// stdioServerOptions keeps state-dependent profiles on legacy revisions. The
// SDK's supported-version context only changes discovery on stdio; it does not
// enforce dispatch, so reject modern requests before any method handler runs.
func stdioServerOptions(opts *ServerOptions) []server.ServerOption {
	modern := !opts.IsCloudConfigured() && !opts.MultiSessionPulsar && opts.UseExternalKafka != opts.UseExternalPulsar
	hooks := &server.Hooks{}
	hooks.AddOnRequestInitialization(func(_ context.Context, _ any, message any) error {
		var request struct {
			Method protocol.MCPMethod
			Params struct {
				Meta *protocol.Meta `json:"_meta"`
			}
		}
		raw, ok := message.(json.RawMessage)
		if !ok || json.Unmarshal(raw, &request) != nil || request.Params.Meta == nil {
			return nil // Leave malformed envelopes to the SDK.
		}
		if !modern && protocol.IsModernProtocol(request.Params.Meta.ProtocolVersion()) && request.Method != protocol.MethodServerDiscover {
			// v1.1.0 maps initialization-hook errors to INVALID_REQUEST. Avoid
			// forking the SDK just to customize that code; the request is refused.
			return fmt.Errorf("this stdio profile requires a legacy protocol version; modern requests require a fixed external backend")
		}
		return nil
	})
	hooks.AddAfterDiscover(func(ctx context.Context, _ any, _ *protocol.DiscoverRequest, result *protocol.DiscoverResult) {
		if !modern {
			result.SupportedVersions = protocol.LegacyProtocolVersions()
			return
		}
		if info := server.RequestProtocolInfoFromContext(ctx); info != nil && info.Modern && result.Capabilities.Resources != nil {
			// This is a per-response snapshot, not the shared legacy capability.
			result.Capabilities.Resources.Subscribe = false
			result.Capabilities.Resources.ListChanged = false
		}
	})
	hooks.AddBeforeSubscriptionsListen(func(ctx context.Context, _ any, request *protocol.SubscriptionsListenRequest) {
		if info := server.RequestProtocolInfoFromContext(ctx); info != nil && info.Modern {
			// Decline unsupported resource subscriptions using the SDK's normal
			// acknowledged-filter response. Never mutate shared capabilities.
			request.Params.Notifications.ResourceSubscriptions = nil
			request.Params.Notifications.ResourcesListChanged = false
		}
	})
	return []server.ServerOption{server.WithHooks(hooks)}
}

func initLogger(filePath string) (*logrus.Logger, error) {
	if filePath == "" {
		return logrus.New(), nil
	}

	fd, err := os.OpenFile(filepath.Clean(filePath), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(fd)
	return logger, nil
}
