// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
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
	stdlog "log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp"
)

func newMcpServer(ctx context.Context, configOpts *ServerOptions, logrusLogger *logrus.Logger) (*mcp.Server, error) {
	snConfig := configOpts.Options.LoadConfigOrDie()
	var s *server.MCPServer
	var mcpServer *mcp.Server
	switch {
	case snConfig.KeyFile != "":
		{
			issuer := snConfig.Auth.Issuer()
			userName, err := configOpts.Options.WhoAmI(issuer.Audience)
			if err != nil {
				stdlog.Fatalf("failed to get user name: %v", err)
				os.Exit(1)
			}
			// sncloudClient, err := config.GetAPIClient()
			// if err != nil {
			// 	stdlog.Fatalf("failed to get SNCloud client: %v", err)
			// 	os.Exit(1)
			// }
			// sncloudLogClient, err := config.GetSNCloudLogClient()
			// if err != nil {
			// 	stdlog.Fatalf("failed to get SNCloud log client: %v", err)
			// 	os.Exit(1)
			// }
			mcpServer = mcp.NewServer("streamnative-mcp-server", "0.0.1", logrusLogger, server.WithInstructions(mcp.GetStreamNativeCloudServerInstructions(userName, snConfig)))
			s = mcpServer.MCPServer
			mcp.RegisterPrompts(s)
			mcp.RegisterContextTools(s, configOpts.Features)
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
			mcpServer = mcp.NewServer("streamnative-mcp-server", "0.0.1", logrusLogger, server.WithInstructions(mcp.GetExternalKafkaServerInstructions(snConfig.ExternalKafka.BootstrapServers)))
			mcpServer.KafkaSession = ksession
			s = mcpServer.MCPServer
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
	return mcpServer, nil
}
