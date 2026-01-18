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
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	stdlog "log"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
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
	mcpServer, err := newMcpServer(ctx, configOpts, logger)
	if err != nil {
		return fmt.Errorf("failed to create MCP server: %w", err)
	}

	var transport sdk.Transport = &sdk.StdioTransport{}
	if configOpts.LogCommands {
		transport = &sdk.LoggingTransport{
			Transport: transport,
			Writer:    logger.Writer(),
		}
	}

	// Start listening for messages
	errC := make(chan error, 1)
	go func() {
		errC <- mcpServer.Run(ctx, transport, stdLogger)
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
