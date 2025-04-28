// Copyright (c) 2025 StreamNative, Inc.. All Rights Reserved.

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
	"github.com/streamnative/streamnative-mcp-server/pkg/cmd/mcp/log"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp"
)

func NewCmdMcpStdioServer(configOpts *ServerOptions) *cobra.Command {
	stdioCmd := &cobra.Command{
		Use:   "stdio",
		Short: "Start stdio server",
		Long:  `Start a server that communicates via standard input/output streams using JSON-RPC messages.`,
		Run: func(cmd *cobra.Command, args []string) {
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
	var logger *logrus.Logger
	var err error

	if configOpts.LogFile != "" || configOpts.LogCommands {
		logger, err = log.InitLogger(configOpts.LogFile)
		if err != nil {
			return fmt.Errorf("failed to initialize logger: %w", err)
		}

		if configOpts.LogFile != "" {
			logger.Infof("Logging to file: %s", configOpts.LogFile)
		}

		if configOpts.LogCommands {
			logger.Info("Command logging enabled")
		}
	}

	// Create a new MCP server
	ctx = context.WithValue(ctx, mcp.OptionsKey, configOpts.Options)
	stdioServer := server.NewStdioServer(newStdioServer(configOpts))

	// Start listening for messages
	errC := make(chan error, 1)
	go func() {
		in, out := io.Reader(os.Stdin), io.Writer(os.Stdout)

		// If command logging is enabled, wrap the IO with a logger
		if configOpts.LogCommands && logger != nil {
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

func newStdioServer(configOpts *ServerOptions) *server.MCPServer {
	issuer := configOpts.Options.LoadConfigOrDie().Auth.Issuer()
	userName, err := configOpts.Options.WhoAmI(issuer.Audience)
	if err != nil {
		stdlog.Fatalf("failed to get user name: %v", err)
		os.Exit(1)
	}
	// Create a new MCP server
	s := server.NewMCPServer(
		"streamnative-mcp-server",
		"0.0.1",
		server.WithResourceCapabilities(true, true),
		server.WithInstructions(fmt.Sprintf(`StreamNative Cloud MCP Server provides resources and tools for AI agents to interact with StreamNative Cloud resources and services.

### StreamNative Cloud API Server Resources

1. **Organizations**  
   - **Concept**: An organization is the top-level resource in StreamNative Cloud. It is a logical grouping used to manage access permissions and organize resources. Most users need only one organization, but multiple organizations can be created for different departments or teams.
   - **Relationship**: An organization contains PulsarInstances, PulsarClusters, Users, Service Accounts, and Secrets. It acts as a container for all other resources and determines who can access what.

2. **PulsarInstances**  
   - **Concept**: An instance is an environment within an organization, typically associated with a specific cloud provider. Instances can contain multiple PulsarClusters and other components (like Connectors, Functions, etc.). Different teams can use separate instances to isolate their resources.
   - **Types**: Instances can be fully hosted on StreamNative's cloud account or deployed on a user's cloud account via the Bring-Your-Own-Cloud (BYOC) option.
   - **Relationship**: PulsarInstances belong to an organization and contain one or more PulsarClusters. Each PulsarInstance is associated with a specific Data Streaming Engine and cloud provider.

3. **PulsarClusters**  
   - **Concept**: A cluster is a specific deployment unit within an instance, located in a particular cloud region. Clusters provide service endpoints that allow clients to connect using different protocols (like Pulsar, Kafka, MQTT) and perform operations such as sending and receiving messages, running functions, etc.
   - **Relationship**: PulsarClusters belong to a PulsarInstance and can replicate data among themselves within the same instance using Geo-Replication. Clusters are where actual data streaming operations occur.

4. **Users**  
   - **Concept**: A user represents a person who can log in and access StreamNative Cloud resources. Users can be assigned different permissions within an organization.
   - **Relationship**: Users belong to an organization and can access instances and clusters within it, depending on their permissions.

5. **Service Accounts**  
   - **Concept**: A service account represents an application that programmatically accesses StreamNative Cloud resources and Pulsar resources within clusters.
   - **Relationship**: Service accounts belong to an organization and can be used across multiple instances, though authentication credentials or API keys differ per instance.

6. **Secrets**  
   - **Concept**: Secrets are used to store and manage sensitive data such as passwords, tokens, and private keys. Secrets can be referenced in Connectors and Pulsar Functions.
   - **Relationship**: Secrets belong to an organization and can be shared across multiple instances within that organization.

7. **Data Streaming Engine**  
   - **Concept**: The Data Streaming Engine is the core technology that runs StreamNative Cloud clusters. There are two options: Classic Engine and Ursa Engine.
     - **Classic Engine**: The default engine, based on ZooKeeper and BookKeeper, offering low-latency storage suitable for latency-sensitive workloads. It supports Pulsar, Kafka, and MQTT protocols.
     - **Ursa Engine**: A next-generation engine based on Oxia and object storage (like S3), providing cost-optimized storage for latency-relaxed scenarios. It currently focuses on Kafka protocol support.
   - **Relationship**: The Data Streaming Engine is associated with an instance, determining how clusters within that instance operate and what features they support.

### Hierarchical Relationship Summary of Resources
- **Organization** is the top level, containing all other resources.
- **PulsarInstances** belong to an organization, representing a cloud environment, and contain multiple **Clusters**.
- **PulsarClusters** are specific deployment units within an instance, used for actual data streaming operations.
- **Users** and **Service Accounts** belong to an organization, used for accessing resources, with permissions managed by the organization.
- **Secrets** belong to an organization and can be shared across instances for securely storing sensitive data.
- **Data Streaming Engine** is associated with an instance, defining the technical architecture and feature support for clusters.

Logged in as %s.`, userName)),
		server.WithLogging())

	// Set up API server resources and tools here - these will need to be implemented
	mcp.RegisterPrompts(s)

	// Tools
	mcp.RegisterContextTools(s)

	mcp.PulsarAdminAddBrokersTools(s, configOpts.ReadOnly)
	mcp.PulsarAdminAddBrokerStatsTools(s, configOpts.ReadOnly)
	mcp.PulsarAdminAddClusterTools(s, configOpts.ReadOnly)
	mcp.PulsarAdminAddFunctionsWorkerTools(s, configOpts.ReadOnly)
	mcp.PulsarAdminAddNamespaceTools(s, configOpts.ReadOnly)
	mcp.PulsarAdminAddNamespacePolicyTools(s, configOpts.ReadOnly)
	mcp.PulsarAdminAddNsIsolationPolicyTools(s, configOpts.ReadOnly)
	mcp.PulsarAdminAddPackagesTools(s, configOpts.ReadOnly)
	mcp.PulsarAdminAddResourceQuotasTools(s, configOpts.ReadOnly)
	mcp.PulsarAdminAddSchemasTools(s, configOpts.ReadOnly)
	mcp.PulsarAdminAddSubscriptionTools(s, configOpts.ReadOnly)
	mcp.PulsarAdminAddTenantTools(s, configOpts.ReadOnly)
	mcp.PulsarAdminAddTopicTools(s, configOpts.ReadOnly)
	mcp.PulsarAdminAddSinksTools(s, configOpts.ReadOnly)
	mcp.PulsarAdminAddFunctionsTools(s, configOpts.ReadOnly)
	mcp.PulsarAdminAddSourcesTools(s, configOpts.ReadOnly)
	mcp.PulsarAdminAddTopicPolicyTools(s, configOpts.ReadOnly)
	mcp.PulsarClientAddConsumerTools(s, configOpts.ReadOnly)
	mcp.PulsarClientAddProducerTools(s, configOpts.ReadOnly)

	mcp.KafkaAdminAddTopicTools(s, configOpts.ReadOnly)
	mcp.KafkaAdminAddPartitionsTools(s, configOpts.ReadOnly)
	mcp.KafkaAdminAddGroupsTools(s, configOpts.ReadOnly)
	mcp.KafkaAdminAddSchemaRegistryTools(s, configOpts.ReadOnly)
	mcp.KafkaAdminAddKafkaConnectTools(s, configOpts.ReadOnly)
	mcp.KafkaClientAddConsumeTools(s, configOpts.ReadOnly)
	mcp.KafkaClientAddProduceTools(s, configOpts.ReadOnly)
	return s
}
