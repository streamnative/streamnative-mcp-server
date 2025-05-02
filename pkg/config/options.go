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

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	EnvConfigDir               = "SNMCP_CONFIG_DIR"
	GlobalDefaultIssuer        = "https://auth.streamnative.cloud/"
	GlobalDefaultClientID      = "AJYEdHWi9EFekEaUXkPWA2MqQ3lq1NrI"
	GlobalDefaultAudience      = "https://api.streamnative.cloud"
	GlobalDefaultAPIServer     = "https://api.streamnative.cloud"
	GlobalDefaultProxyLocation = "https://proxy.streamnative.cloud"
	GlobalDefaultLogLocation   = "https://log.streamnative.cloud"
)

// Options represents the common options used throughout the program.
type Options struct {
	AuthOptions
	ConfigDir  string
	ConfigPath string
	Server     string
	// the OAuth 2.0 issuer endpoint
	IssuerEndpoint string
	// the audience identifier for the API server (default: server URL)
	Audience string
	// the client ID to use for authorization grants (note: not used for service accounts)
	ClientID       string
	Organization   string
	PulsarInstance string
	PulsarCluster  string
	ProxyLocation  string
	LogLocation    string
	KeyFile        string

	UseExternalKafka  bool
	UseExternalPulsar bool
	Kafka             ExternalKafka
	Pulsar            ExternalPulsar
}

// NewConfigOptions creates and returns a new Options instance with default values
func NewConfigOptions() *Options {
	return &Options{
		AuthOptions: NewDefaultAuthOptions(),
	}
}

// AddFlags adds command-line flags for Options
func (o *Options) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&o.ConfigDir, "config-dir", "",
		"If present, the config directory to use")
	cmd.PersistentFlags().StringVar(&o.KeyFile, "key-file", "",
		"The key file to use for authentication to StreamNative Cloud")
	cmd.PersistentFlags().StringVar(&o.Server, "server", GlobalDefaultAPIServer,
		"The server to connect to")
	cmd.PersistentFlags().StringVar(&o.IssuerEndpoint, "issuer", GlobalDefaultIssuer,
		"The OAuth 2.0 issuer endpoint")
	cmd.PersistentFlags().StringVar(&o.Audience, "audience", GlobalDefaultAudience,
		"The audience identifier for the API server")
	cmd.PersistentFlags().StringVar(&o.ClientID, "client-id", GlobalDefaultClientID,
		"The client ID to use for authorization grants")
	cmd.PersistentFlags().StringVar(&o.Organization, "organization", "",
		"The organization to use for the API server")
	cmd.PersistentFlags().StringVar(&o.ProxyLocation, "proxy-location", GlobalDefaultProxyLocation,
		"The proxy location to use for the API server")
	cmd.PersistentFlags().StringVar(&o.LogLocation, "log-location", GlobalDefaultLogLocation,
		"The log location to use for the API server")
	_ = cmd.MarkFlagRequired("organization")
	cmd.PersistentFlags().StringVar(&o.PulsarInstance, "pulsar-instance", "",
		"The default instance to use for the API server")
	cmd.PersistentFlags().StringVar(&o.PulsarCluster, "pulsar-cluster", "",
		"The default cluster to use for the API server")
	cmd.PersistentFlags().BoolVar(&o.UseExternalKafka, "use-external-kafka", false,
		"Use external Kafka")
	cmd.PersistentFlags().BoolVar(&o.UseExternalPulsar, "use-external-pulsar", false,
		"Use external Pulsar")
	cmd.PersistentFlags().StringVar(&o.Kafka.BootstrapServers, "kafka-bootstrap-servers", "",
		"The bootstrap servers to use for Kafka")
	cmd.PersistentFlags().StringVar(&o.Kafka.SchemaRegistryURL, "kafka-schema-registry-url", "",
		"The schema registry URL to use for Kafka")
	cmd.PersistentFlags().StringVar(&o.Kafka.AuthType, "kafka-auth-type", "",
		"The auth type to use for Kafka")
	cmd.PersistentFlags().StringVar(&o.Kafka.AuthMechanism, "kafka-auth-mechanism", "",
		"The auth mechanism to use for Kafka")
	cmd.PersistentFlags().StringVar(&o.Kafka.AuthUser, "kafka-auth-user", "",
		"The auth user to use for Kafka")
	cmd.PersistentFlags().StringVar(&o.Kafka.AuthPass, "kafka-auth-pass", "",
		"The auth password to use for Kafka")
	cmd.PersistentFlags().BoolVar(&o.Kafka.UseTLS, "kafka-use-tls", false,
		"Use TLS for Kafka")
	cmd.PersistentFlags().StringVar(&o.Kafka.ClientKeyFile, "kafka-client-key-file", "",
		"The client key file to use for Kafka")
	cmd.PersistentFlags().StringVar(&o.Kafka.ClientCertFile, "kafka-client-cert-file", "",
		"The client certificate file to use for Kafka")
	cmd.PersistentFlags().StringVar(&o.Kafka.CaFile, "kafka-ca-file", "",
		"The CA file to use for Kafka")
	cmd.PersistentFlags().StringVar(&o.Kafka.SchemaRegistryAuthUser, "kafka-schema-registry-auth-user", "",
		"The auth user to use for the schema registry")
	cmd.PersistentFlags().StringVar(&o.Kafka.SchemaRegistryAuthPass, "kafka-schema-registry-auth-pass", "",
		"The auth password to use for the schema registry")
	cmd.PersistentFlags().StringVar(&o.Kafka.SchemaRegistryBearerToken, "kafka-schema-registry-bearer-token", "",
		"The bearer token to use for the schema registry")
	cmd.PersistentFlags().StringVar(&o.Pulsar.WebServiceURL, "pulsar-web-service-url", "",
		"The web service URL to use for Pulsar")
	cmd.PersistentFlags().StringVar(&o.Pulsar.AuthPlugin, "pulsar-auth-plugin", "",
		"The auth plugin to use for Pulsar")
	cmd.PersistentFlags().StringVar(&o.Pulsar.AuthParams, "pulsar-auth-params", "",
		"The auth params to use for Pulsar")
	cmd.PersistentFlags().BoolVar(&o.Pulsar.TLSAllowInsecureConnection, "pulsar-tls-allow-insecure-connection", false,
		"The TLS allow insecure connection to use for Pulsar")
	cmd.PersistentFlags().BoolVar(&o.Pulsar.TLSEnableHostnameVerification, "pulsar-tls-enable-hostname-verification", true,
		"The TLS enable hostname verification to use for Pulsar")
	cmd.PersistentFlags().StringVar(&o.Pulsar.TLSTrustCertsFilePath, "pulsar-tls-trust-certs-file-path", "",
		"The TLS trust certs file path to use for Pulsar")
	cmd.PersistentFlags().StringVar(&o.Pulsar.TLSCertFile, "pulsar-tls-cert-file", "",
		"The TLS cert file to use for Pulsar")
	cmd.PersistentFlags().StringVar(&o.Pulsar.TLSKeyFile, "pulsar-tls-key-file", "",
		"The TLS key file to use for Pulsar")
	cmd.PersistentFlags().StringVar(&o.Pulsar.Token, "pulsar-token", "",
		"The token to use for Pulsar")

	cmd.MarkFlagsMutuallyExclusive("key-file", "use-external-kafka", "use-external-pulsar")
	cmd.MarkFlagsRequiredTogether("pulsar-cluster", "pulsar-instance")
	cmd.MarkFlagsRequiredTogether("use-external-kafka", "kafka-bootstrap-servers")
	cmd.MarkFlagsRequiredTogether("use-external-pulsar", "pulsar-web-service-url")
	o.AuthOptions.AddFlags(cmd)
}

// Complete completes options from the provided values
func (o *Options) Complete() error {
	if o.ConfigDir == "" {
		o.ConfigDir = os.Getenv(EnvConfigDir)
	}
	if o.ConfigDir == "" {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		o.ConfigDir = filepath.Join(home, ".snmcp")
	}
	if _, err := os.Stat(o.ConfigDir); os.IsNotExist(err) {
		err := os.MkdirAll(o.ConfigDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	}
	o.ConfigPath = filepath.Join(o.ConfigDir, "config")

	err := o.AuthOptions.Complete(o)
	if err != nil {
		return err
	}

	return nil
}

func (o *Options) GetConfigDirectory() string {
	return o.ConfigDir
}

func (o *Options) LoadConfig() (*SnConfig, error) {
	config := &SnConfig{
		Server: o.Server,
		Auth: Auth{
			IssuerEndpoint: o.IssuerEndpoint,
			Audience:       o.Audience,
			ClientID:       o.ClientID,
		},
		Context: Context{
			Organization:   o.Organization,
			PulsarInstance: o.PulsarInstance,
			PulsarCluster:  o.PulsarCluster,
		},
		ProxyLocation: o.ProxyLocation,
		LogLocation:   o.LogLocation,
		KeyFile:       o.KeyFile,
	}
	if o.UseExternalKafka {
		config.ExternalKafka = &o.Kafka
	}
	if o.UseExternalPulsar {
		config.ExternalPulsar = &o.Pulsar
	}
	return config, nil
}

func (o *Options) LoadConfigOrDie() *SnConfig {
	cfg, _ := o.LoadConfig()
	return cfg
}

func (o *Options) SaveConfig(config *SnConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = os.WriteFile(o.ConfigPath, data, 0600)
	if err != nil {
		return fmt.Errorf("WriteFile: %v", err)
	}
	return nil
}
