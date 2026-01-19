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

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

func newMcpServer(_ context.Context, configOpts *ServerOptions, logrusLogger *logrus.Logger) (*mcp.Server, error) {
	snConfig := configOpts.LoadConfigOrDie()
	var mcpServer *mcp.Server
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
			mcpServer = mcp.NewServer("streamnative-mcp-server", "0.0.1", logrusLogger, mcp.WithInstructions(mcp.GetStreamNativeCloudServerInstructions(userName, snConfig)))
			mcpServer.SNCloudSession = session

			serverInstance := mcpServer.MCPServer
			mcp.RegisterPrompts(serverInstance)
			// Skip context tools if pulsar instance and cluster are provided via CLI
			skipContextTools := snConfig.Context.PulsarInstance != "" && snConfig.Context.PulsarCluster != ""
			mcp.RegisterContextTools(serverInstance, configOpts.Features, skipContextTools)
			mcp.StreamNativeAddLogTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
			mcp.StreamNativeAddResourceTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
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
			mcpServer = mcp.NewServer("streamnative-mcp-server", "0.0.1", logrusLogger, mcp.WithInstructions(mcp.GetExternalKafkaServerInstructions(snConfig.ExternalKafka.BootstrapServers)))
			mcpServer.KafkaSession = ksession
		}
	case snConfig.ExternalPulsar != nil:
		{
			mcpServer = mcp.NewServer("streamnative-mcp-server", "0.0.1", logrusLogger, mcp.WithInstructions(mcp.GetExternalPulsarServerInstructions(snConfig.ExternalPulsar.WebServiceURL)))

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

	serverInstance := mcpServer.MCPServer
	mcp.PulsarAdminAddBrokersTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddBrokerStatsTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddClusterTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddFunctionsWorkerTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddNamespaceTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddNamespacePolicyTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddNsIsolationPolicyTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddPackagesTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddResourceQuotasTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddSchemasTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddSubscriptionTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddTenantTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddTopicTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddSinksTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddFunctionsTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddSourcesTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarAdminAddTopicPolicyTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarClientAddConsumerTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.PulsarClientAddProducerTools(serverInstance, configOpts.ReadOnly, configOpts.Features)

	mcp.KafkaAdminAddTopicTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.KafkaAdminAddPartitionsTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.KafkaAdminAddGroupsTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.KafkaAdminAddSchemaRegistryTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.KafkaAdminAddKafkaConnectTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.KafkaClientAddConsumeTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	mcp.KafkaClientAddProduceTools(serverInstance, configOpts.ReadOnly, configOpts.Features)
	return mcpServer, nil
}
