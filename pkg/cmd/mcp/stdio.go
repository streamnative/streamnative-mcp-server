// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package mcp

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	stdlog "log"

	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/streamnative/streamnative-mcp-server/pkg/log"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp"
)

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
	ctx = context.WithValue(ctx, mcp.OptionsKey, configOpts.Options)
	stdLogger := stdlog.New(logger.Writer(), "snmcp-server", 0)
	stdioServer := server.NewStdioServer(newStdioServer(configOpts, logger))

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

func initLogger(filePath string) (*logrus.Logger, error) {
	if filePath == "" {
		return logrus.New(), nil
	}

	fd, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(fd)
	return logger, nil
}

func newStdioServer(configOpts *ServerOptions, logrusLogger *logrus.Logger) *server.MCPServer {
	snConfig := configOpts.Options.LoadConfigOrDie()
	var s *server.MCPServer
	switch {
	case snConfig.KeyFile != "":
		{
			issuer := snConfig.Auth.Issuer()
			userName, err := configOpts.Options.WhoAmI(issuer.Audience)
			if err != nil {
				stdlog.Fatalf("failed to get user name: %v", err)
				os.Exit(1)
			}
			// Create a new MCP server
			s = server.NewMCPServer(
				"streamnative-mcp-server",
				"0.0.1",
				server.WithResourceCapabilities(true, true),
				server.WithInstructions(mcp.GetStreamNativeCloudServerInstructions(userName)),
				server.WithLogging())

			mcp.RegisterPrompts(s)
			mcp.RegisterContextTools(s, configOpts.Features)
		}
	case snConfig.ExternalKafka != nil:
		{
			s = server.NewMCPServer(
				"streamnative-mcp-server/kafka",
				"0.0.1",
				server.WithResourceCapabilities(true, true),
				server.WithInstructions(mcp.GetExternalKafkaServerInstructions(snConfig.ExternalKafka.BootstrapServers)),
				server.WithLogging())
		}
	case snConfig.ExternalPulsar != nil:
		{
			s = server.NewMCPServer(
				"streamnative-mcp-server/pulsar",
				"0.0.1",
				server.WithResourceCapabilities(true, true),
				server.WithInstructions(mcp.GetExternalPulsarServerInstructions(snConfig.ExternalPulsar.WebServiceURL)),
				server.WithLogging())
		}
	default:
		{
			stdlog.Fatalf("no valid configuration found")
			os.Exit(1)
		}
	}

	mcp.PulsarAdminAddBrokersTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddBrokerStatsTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddClusterTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddFunctionsWorkerTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddNamespaceTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddNamespacePolicyTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddNsIsolationPolicyTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddPackagesTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddResourceQuotasTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddSchemasTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddSubscriptionTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddTenantTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddTopicTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddSinksTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddFunctionsTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddSourcesTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddTopicPolicyTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarClientAddConsumerTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarClientAddProducerTools(s, configOpts.ReadOnly, configOpts.Features)

	mcp.KafkaAdminAddTopicTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.KafkaAdminAddPartitionsTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.KafkaAdminAddGroupsTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.KafkaAdminAddSchemaRegistryTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.KafkaAdminAddKafkaConnectTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.KafkaClientAddConsumeTools(s, configOpts.ReadOnly, logrusLogger, configOpts.Features)
	mcp.KafkaClientAddProduceTools(s, configOpts.ReadOnly, configOpts.Features)
	return s
}
