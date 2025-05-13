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

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	// Environment variable prefix
	EnvPrefix = "SNMCP"
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
	// Setup viper with environment variables
	viper.SetEnvPrefix(EnvPrefix)
	viper.AutomaticEnv()
	// Replace dots and dashes in env var names with underscores
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	cmd.PersistentFlags().StringVar(&o.ConfigDir, "config-dir", "",
		"If present, the config directory to use [env: SNMCP_CONFIG_DIR]")
	cmd.PersistentFlags().StringVar(&o.KeyFile, "key-file", "",
		"The key file to use for authentication to StreamNative Cloud [env: SNMCP_KEY_FILE]")
	cmd.PersistentFlags().StringVar(&o.Server, "server", GlobalDefaultAPIServer,
		"The server to connect to [env: SNMCP_SERVER]")
	cmd.PersistentFlags().StringVar(&o.IssuerEndpoint, "issuer", GlobalDefaultIssuer,
		"The OAuth 2.0 issuer endpoint [env: SNMCP_ISSUER]")
	cmd.PersistentFlags().StringVar(&o.Audience, "audience", GlobalDefaultAudience,
		"The audience identifier for the API server [env: SNMCP_AUDIENCE]")
	cmd.PersistentFlags().StringVar(&o.ClientID, "client-id", GlobalDefaultClientID,
		"The client ID to use for authorization grants [env: SNMCP_CLIENT_ID]")
	cmd.PersistentFlags().StringVar(&o.Organization, "organization", "",
		"The organization to use for the API server [env: SNMCP_ORGANIZATION]")
	cmd.PersistentFlags().StringVar(&o.ProxyLocation, "proxy-location", GlobalDefaultProxyLocation,
		"The proxy location to use for the API server [env: SNMCP_PROXY_LOCATION]")
	cmd.PersistentFlags().StringVar(&o.LogLocation, "log-location", GlobalDefaultLogLocation,
		"The log location to use for the API server [env: SNMCP_LOG_LOCATION]")
	_ = cmd.MarkFlagRequired("organization")
	cmd.PersistentFlags().StringVar(&o.PulsarInstance, "pulsar-instance", "",
		"The default instance to use for the API server [env: SNMCP_PULSAR_INSTANCE]")
	cmd.PersistentFlags().StringVar(&o.PulsarCluster, "pulsar-cluster", "",
		"The default cluster to use for the API server [env: SNMCP_PULSAR_CLUSTER]")
	cmd.PersistentFlags().BoolVar(&o.UseExternalKafka, "use-external-kafka", false,
		"Use external Kafka [env: SNMCP_USE_EXTERNAL_KAFKA]")
	cmd.PersistentFlags().BoolVar(&o.UseExternalPulsar, "use-external-pulsar", false,
		"Use external Pulsar [env: SNMCP_USE_EXTERNAL_PULSAR]")
	cmd.PersistentFlags().StringVar(&o.Kafka.BootstrapServers, "kafka-bootstrap-servers", "",
		"The bootstrap servers to use for Kafka [env: SNMCP_KAFKA_BOOTSTRAP_SERVERS]")
	cmd.PersistentFlags().StringVar(&o.Kafka.SchemaRegistryURL, "kafka-schema-registry-url", "",
		"The schema registry URL to use for Kafka [env: SNMCP_KAFKA_SCHEMA_REGISTRY_URL]")
	cmd.PersistentFlags().StringVar(&o.Kafka.AuthType, "kafka-auth-type", "",
		"The auth type to use for Kafka [env: SNMCP_KAFKA_AUTH_TYPE]")
	cmd.PersistentFlags().StringVar(&o.Kafka.AuthMechanism, "kafka-auth-mechanism", "",
		"The auth mechanism to use for Kafka [env: SNMCP_KAFKA_AUTH_MECHANISM]")
	cmd.PersistentFlags().StringVar(&o.Kafka.AuthUser, "kafka-auth-user", "",
		"The auth user to use for Kafka [env: SNMCP_KAFKA_AUTH_USER]")
	cmd.PersistentFlags().StringVar(&o.Kafka.AuthPass, "kafka-auth-pass", "",
		"The auth password to use for Kafka [env: SNMCP_KAFKA_AUTH_PASS]")
	cmd.PersistentFlags().BoolVar(&o.Kafka.UseTLS, "kafka-use-tls", false,
		"Use TLS for Kafka [env: SNMCP_KAFKA_USE_TLS]")
	cmd.PersistentFlags().StringVar(&o.Kafka.ClientKeyFile, "kafka-client-key-file", "",
		"The client key file to use for Kafka [env: SNMCP_KAFKA_CLIENT_KEY_FILE]")
	cmd.PersistentFlags().StringVar(&o.Kafka.ClientCertFile, "kafka-client-cert-file", "",
		"The client certificate file to use for Kafka [env: SNMCP_KAFKA_CLIENT_CERT_FILE]")
	cmd.PersistentFlags().StringVar(&o.Kafka.CaFile, "kafka-ca-file", "",
		"The CA file to use for Kafka [env: SNMCP_KAFKA_CA_FILE]")
	cmd.PersistentFlags().StringVar(&o.Kafka.SchemaRegistryAuthUser, "kafka-schema-registry-auth-user", "",
		"The auth user to use for the schema registry [env: SNMCP_KAFKA_SCHEMA_REGISTRY_AUTH_USER]")
	cmd.PersistentFlags().StringVar(&o.Kafka.SchemaRegistryAuthPass, "kafka-schema-registry-auth-pass", "",
		"The auth password to use for the schema registry [env: SNMCP_KAFKA_SCHEMA_REGISTRY_AUTH_PASS]")
	cmd.PersistentFlags().StringVar(&o.Kafka.SchemaRegistryBearerToken, "kafka-schema-registry-bearer-token", "",
		"The bearer token to use for the schema registry [env: SNMCP_KAFKA_SCHEMA_REGISTRY_BEARER_TOKEN]")
	cmd.PersistentFlags().StringVar(&o.Pulsar.WebServiceURL, "pulsar-web-service-url", "",
		"The web service URL to use for Pulsar [env: SNMCP_PULSAR_WEB_SERVICE_URL]")
	cmd.PersistentFlags().StringVar(&o.Pulsar.AuthPlugin, "pulsar-auth-plugin", "",
		"The auth plugin to use for Pulsar [env: SNMCP_PULSAR_AUTH_PLUGIN]")
	cmd.PersistentFlags().StringVar(&o.Pulsar.AuthParams, "pulsar-auth-params", "",
		"The auth params to use for Pulsar [env: SNMCP_PULSAR_AUTH_PARAMS]")
	cmd.PersistentFlags().BoolVar(&o.Pulsar.TLSAllowInsecureConnection, "pulsar-tls-allow-insecure-connection", false,
		"The TLS allow insecure connection to use for Pulsar [env: SNMCP_PULSAR_TLS_ALLOW_INSECURE_CONNECTION]")
	cmd.PersistentFlags().BoolVar(&o.Pulsar.TLSEnableHostnameVerification, "pulsar-tls-enable-hostname-verification", true,
		"The TLS enable hostname verification to use for Pulsar [env: SNMCP_PULSAR_TLS_ENABLE_HOSTNAME_VERIFICATION]")
	cmd.PersistentFlags().StringVar(&o.Pulsar.TLSTrustCertsFilePath, "pulsar-tls-trust-certs-file-path", "",
		"The TLS trust certs file path to use for Pulsar [env: SNMCP_PULSAR_TLS_TRUST_CERTS_FILE_PATH]")
	cmd.PersistentFlags().StringVar(&o.Pulsar.TLSCertFile, "pulsar-tls-cert-file", "",
		"The TLS cert file to use for Pulsar [env: SNMCP_PULSAR_TLS_CERT_FILE]")
	cmd.PersistentFlags().StringVar(&o.Pulsar.TLSKeyFile, "pulsar-tls-key-file", "",
		"The TLS key file to use for Pulsar [env: SNMCP_PULSAR_TLS_KEY_FILE]")
	cmd.PersistentFlags().StringVar(&o.Pulsar.Token, "pulsar-token", "",
		"The token to use for Pulsar [env: SNMCP_PULSAR_TOKEN]")

	cmd.MarkFlagsMutuallyExclusive("key-file", "use-external-kafka", "use-external-pulsar")
	cmd.MarkFlagsRequiredTogether("pulsar-cluster", "pulsar-instance")
	cmd.MarkFlagsRequiredTogether("use-external-kafka", "kafka-bootstrap-servers")
	cmd.MarkFlagsRequiredTogether("use-external-pulsar", "pulsar-web-service-url")
	o.AuthOptions.AddFlags(cmd)

	// Bind command line flags to viper
	viper.BindPFlag("config-dir", cmd.PersistentFlags().Lookup("config-dir"))
	viper.BindPFlag("key-file", cmd.PersistentFlags().Lookup("key-file"))
	viper.BindPFlag("server", cmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("issuer", cmd.PersistentFlags().Lookup("issuer"))
	viper.BindPFlag("audience", cmd.PersistentFlags().Lookup("audience"))
	viper.BindPFlag("client-id", cmd.PersistentFlags().Lookup("client-id"))
	viper.BindPFlag("organization", cmd.PersistentFlags().Lookup("organization"))
	viper.BindPFlag("proxy-location", cmd.PersistentFlags().Lookup("proxy-location"))
	viper.BindPFlag("log-location", cmd.PersistentFlags().Lookup("log-location"))
	viper.BindPFlag("pulsar-instance", cmd.PersistentFlags().Lookup("pulsar-instance"))
	viper.BindPFlag("pulsar-cluster", cmd.PersistentFlags().Lookup("pulsar-cluster"))
	viper.BindPFlag("use-external-kafka", cmd.PersistentFlags().Lookup("use-external-kafka"))
	viper.BindPFlag("use-external-pulsar", cmd.PersistentFlags().Lookup("use-external-pulsar"))
	viper.BindPFlag("kafka-bootstrap-servers", cmd.PersistentFlags().Lookup("kafka-bootstrap-servers"))
	viper.BindPFlag("kafka-schema-registry-url", cmd.PersistentFlags().Lookup("kafka-schema-registry-url"))
	viper.BindPFlag("kafka-auth-type", cmd.PersistentFlags().Lookup("kafka-auth-type"))
	viper.BindPFlag("kafka-auth-mechanism", cmd.PersistentFlags().Lookup("kafka-auth-mechanism"))
	viper.BindPFlag("kafka-auth-user", cmd.PersistentFlags().Lookup("kafka-auth-user"))
	viper.BindPFlag("kafka-auth-pass", cmd.PersistentFlags().Lookup("kafka-auth-pass"))
	viper.BindPFlag("kafka-use-tls", cmd.PersistentFlags().Lookup("kafka-use-tls"))
	viper.BindPFlag("kafka-client-key-file", cmd.PersistentFlags().Lookup("kafka-client-key-file"))
	viper.BindPFlag("kafka-client-cert-file", cmd.PersistentFlags().Lookup("kafka-client-cert-file"))
	viper.BindPFlag("kafka-ca-file", cmd.PersistentFlags().Lookup("kafka-ca-file"))
	viper.BindPFlag("kafka-schema-registry-auth-user", cmd.PersistentFlags().Lookup("kafka-schema-registry-auth-user"))
	viper.BindPFlag("kafka-schema-registry-auth-pass", cmd.PersistentFlags().Lookup("kafka-schema-registry-auth-pass"))
	viper.BindPFlag("kafka-schema-registry-bearer-token", cmd.PersistentFlags().Lookup("kafka-schema-registry-bearer-token"))
	viper.BindPFlag("pulsar-web-service-url", cmd.PersistentFlags().Lookup("pulsar-web-service-url"))
	viper.BindPFlag("pulsar-auth-plugin", cmd.PersistentFlags().Lookup("pulsar-auth-plugin"))
	viper.BindPFlag("pulsar-auth-params", cmd.PersistentFlags().Lookup("pulsar-auth-params"))
	viper.BindPFlag("pulsar-tls-allow-insecure-connection", cmd.PersistentFlags().Lookup("pulsar-tls-allow-insecure-connection"))
	viper.BindPFlag("pulsar-tls-enable-hostname-verification", cmd.PersistentFlags().Lookup("pulsar-tls-enable-hostname-verification"))
	viper.BindPFlag("pulsar-tls-trust-certs-file-path", cmd.PersistentFlags().Lookup("pulsar-tls-trust-certs-file-path"))
	viper.BindPFlag("pulsar-tls-cert-file", cmd.PersistentFlags().Lookup("pulsar-tls-cert-file"))
	viper.BindPFlag("pulsar-tls-key-file", cmd.PersistentFlags().Lookup("pulsar-tls-key-file"))
	viper.BindPFlag("pulsar-token", cmd.PersistentFlags().Lookup("pulsar-token"))
}

// Complete completes options from the provided values
func (o *Options) Complete() error {
	// First try to get config directory from environment variables
	o.ConfigDir = viper.GetString("config-dir")
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

	// Read configurations from viper (environment variables)
	if o.Server == GlobalDefaultAPIServer {
		if server := viper.GetString("server"); server != "" {
			o.Server = server
		}
	}
	if o.IssuerEndpoint == GlobalDefaultIssuer {
		if issuer := viper.GetString("issuer"); issuer != "" {
			o.IssuerEndpoint = issuer
		}
	}
	if o.Audience == GlobalDefaultAudience {
		if audience := viper.GetString("audience"); audience != "" {
			o.Audience = audience
		}
	}
	if o.ClientID == GlobalDefaultClientID {
		if clientID := viper.GetString("client-id"); clientID != "" {
			o.ClientID = clientID
		}
	}
	if o.Organization == "" {
		if org := viper.GetString("organization"); org != "" {
			o.Organization = org
		}
	}
	if o.ProxyLocation == GlobalDefaultProxyLocation {
		if proxy := viper.GetString("proxy-location"); proxy != "" {
			o.ProxyLocation = proxy
		}
	}
	if o.LogLocation == GlobalDefaultLogLocation {
		if log := viper.GetString("log-location"); log != "" {
			o.LogLocation = log
		}
	}
	if o.KeyFile == "" {
		if keyFile := viper.GetString("key-file"); keyFile != "" {
			o.KeyFile = keyFile
		}
	}
	if o.PulsarInstance == "" {
		if instance := viper.GetString("pulsar-instance"); instance != "" {
			o.PulsarInstance = instance
		}
	}
	if o.PulsarCluster == "" {
		if cluster := viper.GetString("pulsar-cluster"); cluster != "" {
			o.PulsarCluster = cluster
		}
	}

	// Boolean values
	if !o.UseExternalKafka {
		o.UseExternalKafka = viper.GetBool("use-external-kafka")
	}
	if !o.UseExternalPulsar {
		o.UseExternalPulsar = viper.GetBool("use-external-pulsar")
	}

	// Kafka configuration
	if o.Kafka.BootstrapServers == "" {
		if servers := viper.GetString("kafka-bootstrap-servers"); servers != "" {
			o.Kafka.BootstrapServers = servers
		}
	}
	if o.Kafka.SchemaRegistryURL == "" {
		if url := viper.GetString("kafka-schema-registry-url"); url != "" {
			o.Kafka.SchemaRegistryURL = url
		}
	}
	if o.Kafka.AuthType == "" {
		if authType := viper.GetString("kafka-auth-type"); authType != "" {
			o.Kafka.AuthType = authType
		}
	}
	if o.Kafka.AuthMechanism == "" {
		if authMech := viper.GetString("kafka-auth-mechanism"); authMech != "" {
			o.Kafka.AuthMechanism = authMech
		}
	}
	if o.Kafka.AuthUser == "" {
		if authUser := viper.GetString("kafka-auth-user"); authUser != "" {
			o.Kafka.AuthUser = authUser
		}
	}
	if o.Kafka.AuthPass == "" {
		if authPass := viper.GetString("kafka-auth-pass"); authPass != "" {
			o.Kafka.AuthPass = authPass
		}
	}
	if !o.Kafka.UseTLS {
		o.Kafka.UseTLS = viper.GetBool("kafka-use-tls")
	}
	if o.Kafka.ClientKeyFile == "" {
		if keyFile := viper.GetString("kafka-client-key-file"); keyFile != "" {
			o.Kafka.ClientKeyFile = keyFile
		}
	}
	if o.Kafka.ClientCertFile == "" {
		if certFile := viper.GetString("kafka-client-cert-file"); certFile != "" {
			o.Kafka.ClientCertFile = certFile
		}
	}
	if o.Kafka.CaFile == "" {
		if caFile := viper.GetString("kafka-ca-file"); caFile != "" {
			o.Kafka.CaFile = caFile
		}
	}
	if o.Kafka.SchemaRegistryAuthUser == "" {
		if srUser := viper.GetString("kafka-schema-registry-auth-user"); srUser != "" {
			o.Kafka.SchemaRegistryAuthUser = srUser
		}
	}
	if o.Kafka.SchemaRegistryAuthPass == "" {
		if srPass := viper.GetString("kafka-schema-registry-auth-pass"); srPass != "" {
			o.Kafka.SchemaRegistryAuthPass = srPass
		}
	}
	if o.Kafka.SchemaRegistryBearerToken == "" {
		if srToken := viper.GetString("kafka-schema-registry-bearer-token"); srToken != "" {
			o.Kafka.SchemaRegistryBearerToken = srToken
		}
	}

	// Pulsar configuration
	if o.Pulsar.WebServiceURL == "" {
		if wsURL := viper.GetString("pulsar-web-service-url"); wsURL != "" {
			o.Pulsar.WebServiceURL = wsURL
		}
	}
	if o.Pulsar.AuthPlugin == "" {
		if authPlugin := viper.GetString("pulsar-auth-plugin"); authPlugin != "" {
			o.Pulsar.AuthPlugin = authPlugin
		}
	}
	if o.Pulsar.AuthParams == "" {
		if authParams := viper.GetString("pulsar-auth-params"); authParams != "" {
			o.Pulsar.AuthParams = authParams
		}
	}
	if !o.Pulsar.TLSAllowInsecureConnection {
		o.Pulsar.TLSAllowInsecureConnection = viper.GetBool("pulsar-tls-allow-insecure-connection")
	}
	if o.Pulsar.TLSEnableHostnameVerification {
		// Only override if explicitly set to false in env var
		if !viper.GetBool("pulsar-tls-enable-hostname-verification") {
			o.Pulsar.TLSEnableHostnameVerification = false
		}
	}
	if o.Pulsar.TLSTrustCertsFilePath == "" {
		if trustCerts := viper.GetString("pulsar-tls-trust-certs-file-path"); trustCerts != "" {
			o.Pulsar.TLSTrustCertsFilePath = trustCerts
		}
	}
	if o.Pulsar.TLSCertFile == "" {
		if certFile := viper.GetString("pulsar-tls-cert-file"); certFile != "" {
			o.Pulsar.TLSCertFile = certFile
		}
	}
	if o.Pulsar.TLSKeyFile == "" {
		if keyFile := viper.GetString("pulsar-tls-key-file"); keyFile != "" {
			o.Pulsar.TLSKeyFile = keyFile
		}
	}
	if o.Pulsar.Token == "" {
		if token := viper.GetString("pulsar-token"); token != "" {
			o.Pulsar.Token = token
		}
	}

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
