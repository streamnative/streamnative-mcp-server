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
	stdlog "log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp"
)

func newMcpServer(configOpts *ServerOptions, logrusLogger *logrus.Logger) *server.MCPServer {
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
				server.WithInstructions(mcp.GetStreamNativeCloudServerInstructions(userName, snConfig)),
				server.WithLogging())

			mcp.RegisterPrompts(s)
			mcp.RegisterContextTools(s, configOpts.Features)
			mcp.StreamNativeAddLogTools(s, configOpts.ReadOnly, configOpts.Features)
			mcp.StreamNativeAddResourceTools(s, configOpts.ReadOnly, configOpts.Features)
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
