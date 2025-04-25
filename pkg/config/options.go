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
	ClientID      string
	Organization  string
	ProxyLocation string
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
	cmd.MarkFlagRequired("organization")
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
			Organization: o.Organization,
		},
		ProxyLocation: o.ProxyLocation,
	}
	return config, nil
}

func (o *Options) LoadConfigOrDie() *SnConfig {
	config := &SnConfig{
		Server: o.Server,
		Auth: Auth{
			IssuerEndpoint: o.IssuerEndpoint,
			Audience:       o.Audience,
			ClientID:       o.ClientID,
		},
		Context: Context{
			Organization: o.Organization,
		},
		ProxyLocation: o.ProxyLocation,
	}
	return config
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
