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
	"slices"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/streamnative/streamnative-mcp-server/pkg/auth"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// ServerOptions is the options for the MCP server commands
type ServerOptions struct {
	ReadOnly    bool
	LogFile     string
	LogCommands bool
	Features    []string
	*config.Options
}

func NewMcpServerOptions(configOpts *config.Options) *ServerOptions {
	return &ServerOptions{
		Options: configOpts,
	}
}

func (o *ServerOptions) Complete() error {
	if err := o.Options.Complete(); err != nil {
		return err
	}

	snConfig := o.Options.LoadConfigOrDie()

	// If the key file is provided, use it to authenticate to StreamNative Cloud
	switch {
	case snConfig.KeyFile != "":
		{
			issuer := snConfig.Auth.Issuer()

			// authorize
			flow, err := o.newClientCredentialsFlow(issuer, o.KeyFile)
			if err != nil {
				return errors.Wrap(err, "configuration error: unable to use client credential flow")
			}
			grant, err := flow.Authorize()
			if err != nil {
				return errors.Wrap(err, "activation failed")
			}

			// persist the authorization data
			if err = o.Options.SaveGrant(issuer.Audience, *grant); err != nil {
				return errors.Wrap(err, "Unable to store the authorization data")
			}

			err = config.InitSNCloudClient(
				issuer.IssuerEndpoint, issuer.Audience, o.KeyFile, o.Options.Server, 30*time.Second, o.Options.Store)
			if err != nil {
				return errors.Wrap(err, "failed to initialize StreamNative Cloud client")
			}

			if o.Options.PulsarInstance != "" && o.Options.PulsarCluster != "" {
				err = mcp.SetContext(o.Options, o.Options.PulsarInstance, o.Options.PulsarCluster)
				if err != nil {
					mcp.ResetMcpContext()
					return errors.Wrap(err, "failed to set StreamNative Cloud context")
				}
			}

			if len(o.Features) != 0 {
				requiredFeatures := []mcp.Feature{
					mcp.FeatureStreamNativeCloud,
				}
				for _, feature := range requiredFeatures {
					if !slices.Contains(o.Features, string(feature)) {
						o.Features = append(o.Features, string(feature))
					}
				}
			} else {
				o.Features = []string{string(mcp.FeatureAll)}
			}
		}
	case snConfig.ExternalKafka != nil:
		{
			if len(o.Features) != 0 {
				return errors.New("kafka-only mode does not support additional features")
			}
			o.Features = []string{string(mcp.FeatureKafkaClient), string(mcp.FeatureKafkaAdmin), string(mcp.FeatureKafkaAdminSchemaRegistry)}
			err := kafka.NewCurrentKafkaContext(kafka.KafkaContext{
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
				return errors.Wrap(err, "failed to set external Kafka context")
			}
		}
	case snConfig.ExternalPulsar != nil:
		{
			if len(o.Features) != 0 {
				return errors.New("pulsar-only mode does not support additional features")
			}
			o.Features = []string{string(mcp.FeatureAllPulsar)}
			err := pulsar.NewCurrentPulsarContext(pulsar.PulsarContext{
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
			}, nil, nil)
			if err != nil {
				return errors.Wrap(err, "failed to set external Pulsar context")
			}
		}
	default:
		{
			return errors.New("no valid configuration found")
		}
	}
	return nil
}

func (o *ServerOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&o.ReadOnly, "read-only", "r", false, "Read-only mode")
	cmd.PersistentFlags().StringVar(&o.LogFile, "log-file", "", "Path to log file")
	cmd.PersistentFlags().BoolVar(&o.LogCommands, "enable-command-logging", false, "When enabled, the server will log all command requests and responses to the log file")
	cmd.PersistentFlags().StringSliceVar(&o.Features, "features", []string{}, "Features to enable, defaults to `all`")
}

func (o *ServerOptions) newClientCredentialsFlow(
	issuer auth.Issuer, keyFile string) (*auth.ClientCredentialsFlow, error) {
	flow, err := auth.NewDefaultClientCredentialsFlow(issuer, keyFile)
	if err != nil {
		return nil, err
	}
	return flow, nil
}
