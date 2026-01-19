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
	stdlog "log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

func newMcpServer(_ context.Context, configOpts *ServerOptions, logrusLogger *logrus.Logger) (*mcp.LegacyServer, error) {
	snConfig := configOpts.LoadConfigOrDie()
	var s *server.MCPServer
	var mcpServer *mcp.LegacyServer
	switch {
	case snConfig.KeyFile != "":
		{
			issuer := snConfig.Auth.Issuer()
			userName, err := configOpts.WhoAmI(issuer.Audience)
			if err != nil {
				stdlog.Fatalf("failed to get user name: %v", err)
				os.Exit(1)
			}
			// Create StreamNative Cloud session and set as default
			session, err := config.NewSNCloudSessionFromOptions(configOpts.Options)
			if err != nil {
				return nil, errors.Wrap(err, "failed to create StreamNative Cloud session")
			}
			mcpServer = mcp.NewLegacyServer("streamnative-mcp-server", "0.0.1", logrusLogger, server.WithInstructions(mcp.GetStreamNativeCloudServerInstructions(userName, snConfig)))
			mcpServer.SNCloudSession = session

			s = mcpServer.MCPServer
			mcp.RegisterPrompts(s)
			// Skip context tools if pulsar instance and cluster are provided via CLI
			skipContextTools := snConfig.Context.PulsarInstance != "" && snConfig.Context.PulsarCluster != ""
			mcp.RegisterContextTools(s, configOpts.Features, skipContextTools)
			mcp.StreamNativeAddLogTools(s, configOpts.ReadOnly, configOpts.Features)
			mcp.StreamNativeAddResourceTools(s, configOpts.ReadOnly, configOpts.Features)
		}
	case snConfig.ExternalKafka != nil:
		{
			ksession, err := kafka.NewSession(kafka.KafkaContext{
				BootstrapServers:          snConfig.ExternalKafka.BootstrapServers,
				AuthType:                  snConfig.ExternalKafka.AuthType,
				AuthMechanism:             snConfig.ExternalKafka.AuthMechanism,
				AuthUser:                  snConfig.ExternalKafka.AuthUser,
				AuthPass:                  snConfig.ExternalKafka.AuthPass,
				UseTLS:                    snConfig.ExternalKafka.UseTLS,
				ClientKeyFile:             snConfig.ExternalKafka.ClientKeyFile,
				ClientCertFile:            snConfig.ExternalKafka.ClientCertFile,
				CaFile:                    snConfig.ExternalKafka.CaFile,
				SchemaRegistryURL:         snConfig.ExternalKafka.SchemaRegistryURL,
				SchemaRegistryAuthUser:    snConfig.ExternalKafka.SchemaRegistryAuthUser,
				SchemaRegistryAuthPass:    snConfig.ExternalKafka.SchemaRegistryAuthPass,
				SchemaRegistryBearerToken: snConfig.ExternalKafka.SchemaRegistryBearerToken,
			})
			if err != nil {
				return nil, errors.Wrap(err, "failed to set external Kafka context")
			}
			mcpServer = mcp.NewLegacyServer("streamnative-mcp-server", "0.0.1", logrusLogger, server.WithInstructions(mcp.GetExternalKafkaServerInstructions(snConfig.ExternalKafka.BootstrapServers)))
			mcpServer.KafkaSession = ksession
			s = mcpServer.MCPServer
		}
	case snConfig.ExternalPulsar != nil:
		{
			mcpServer = mcp.NewLegacyServer("streamnative-mcp-server", "0.0.1", logrusLogger, server.WithInstructions(mcp.GetExternalPulsarServerInstructions(snConfig.ExternalPulsar.WebServiceURL)))
			s = mcpServer.MCPServer

			// Only create global PulsarSession if not in multi-session mode
			// In multi-session mode, each request must provide its own token via Authorization header
			if !configOpts.MultiSessionPulsar {
				psession, err := pulsar.NewSession(pulsar.PulsarContext{
					ServiceURL:                    snConfig.ExternalPulsar.ServiceURL,
					WebServiceURL:                 snConfig.ExternalPulsar.WebServiceURL,
					AuthPlugin:                    snConfig.ExternalPulsar.AuthPlugin,
					AuthParams:                    snConfig.ExternalPulsar.AuthParams,
					Token:                         snConfig.ExternalPulsar.Token,
					TLSAllowInsecureConnection:    snConfig.ExternalPulsar.TLSAllowInsecureConnection,
					TLSEnableHostnameVerification: snConfig.ExternalPulsar.TLSEnableHostnameVerification,
					TLSTrustCertsFilePath:         snConfig.ExternalPulsar.TLSTrustCertsFilePath,
					TLSCertFile:                   snConfig.ExternalPulsar.TLSCertFile,
					TLSKeyFile:                    snConfig.ExternalPulsar.TLSKeyFile,
				})
				if err != nil {
					return nil, errors.Wrap(err, "failed to set external Pulsar context")
				}
				mcpServer.PulsarSession = psession
			}
		}
	default:
		{
			stdlog.Fatalf("no valid configuration found")
			os.Exit(1)
		}
	}

	mcp.PulsarAdminAddBrokersToolsLegacy(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddBrokerStatsTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddClusterToolsLegacy(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddFunctionsWorkerTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddNamespaceToolsLegacy(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddNamespacePolicyTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddNsIsolationPolicyTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddPackagesTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddResourceQuotasTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddSchemasToolsLegacy(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddSubscriptionToolsLegacy(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddTenantToolsLegacy(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddTopicToolsLegacy(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddSinksToolsLegacy(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddFunctionsToolsLegacy(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddSourcesToolsLegacy(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddTopicPolicyTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarClientAddConsumerTools(s, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarClientAddProducerTools(s, configOpts.ReadOnly, configOpts.Features)

	mcp.KafkaAdminAddTopicToolsLegacy(s, configOpts.ReadOnly, configOpts.Features)
	mcp.KafkaAdminAddPartitionsToolsLegacy(s, configOpts.ReadOnly, configOpts.Features)
	mcp.KafkaAdminAddGroupsToolsLegacy(s, configOpts.ReadOnly, configOpts.Features)
	mcp.KafkaAdminAddSchemaRegistryToolsLegacy(s, configOpts.ReadOnly, configOpts.Features)
	mcp.KafkaAdminAddKafkaConnectToolsLegacy(s, configOpts.ReadOnly, configOpts.Features)
	mcp.KafkaClientAddConsumeToolsLegacy(s, configOpts.ReadOnly, logrusLogger, configOpts.Features)
	mcp.KafkaClientAddProduceToolsLegacy(s, configOpts.ReadOnly, configOpts.Features)
	return mcpServer, nil
}
