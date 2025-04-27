// Copyright (c) 2025 StreamNative, Inc.. All Rights Reserved.

package mcp

import (
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/streamnative/streamnative-mcp-server/pkg/auth"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
)

// ServerOptions is the options for the MCP server commands
type ServerOptions struct {
	ReadOnly    bool
	KeyFile     string
	LogFile     string
	LogCommands bool
	*config.Options
}

func NewMcpServerOptions(configOpts *config.Options) *ServerOptions {
	return &ServerOptions{
		Options: configOpts,
	}
}

func (o *ServerOptions) Complete() error {
	if o.KeyFile == "" {
		log.Fatalf("key-file is required")
		os.Exit(1)
	}

	if err := o.Options.Complete(); err != nil {
		return err
	}

	issuer := o.Options.LoadConfigOrDie().Auth.Issuer()

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

	return nil
}

func (o *ServerOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&o.KeyFile, "key-file", "", "Path to the key file")
	cmd.PersistentFlags().BoolVarP(&o.ReadOnly, "read-only", "r", false, "Read-only mode")
	cmd.PersistentFlags().StringVar(&o.LogFile, "log-file", "", "Path to log file")
	cmd.PersistentFlags().BoolVar(&o.LogCommands, "enable-command-logging", false, "When enabled, the server will log all command requests and responses to the log file")
}

func (o *ServerOptions) newClientCredentialsFlow(
	issuer auth.Issuer, keyFile string) (*auth.ClientCredentialsFlow, error) {
	flow, err := auth.NewDefaultClientCredentialsFlow(issuer, keyFile)
	if err != nil {
		return nil, err
	}
	return flow, nil
}
