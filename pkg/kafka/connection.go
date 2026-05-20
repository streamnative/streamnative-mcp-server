// Copyright 2026 StreamNative
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

// Package kafka provides Kafka connection and client helpers.
package kafka

import (
	"crypto/tls"
	"fmt"
	"strings"
	"sync"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kversion"
	"github.com/twmb/franz-go/pkg/sasl/plain"
	"github.com/twmb/franz-go/pkg/sasl/scram"
	"github.com/twmb/franz-go/pkg/sr"
	"github.com/twmb/tlscfg"
)

//nolint:revive
type KafkaContext struct {
	BootstrapServers  string
	AuthType          string
	AuthMechanism     string
	AuthUser          string
	AuthPass          string
	UseTLS            bool
	ClientKeyFile     string
	ClientCertFile    string
	CaFile            string
	SchemaRegistryURL string
	ConnectURL        string

	SchemaRegistryAuthUser    string
	SchemaRegistryAuthPass    string
	SchemaRegistryBearerToken string

	ConnectAuthUser string
	ConnectAuthPass string
}

// Session represents a Kafka session
type Session struct {
	Ctx                  KafkaContext
	Client               *kgo.Client
	AdminClient          *kadm.Client
	SchemaRegistryClient *sr.Client
	ConnectClient        Connect
	Options              []kgo.Opt
	mutex                sync.RWMutex
}

// NewSession creates a new Kafka session with the given context
// This function dynamically constructs clients without relying on global state
func NewSession(ctx KafkaContext) (*Session, error) {
	if ctx.BootstrapServers == "" {
		return nil, fmt.Errorf("bootstrap servers are required")
	}

	session := &Session{
		Ctx: ctx,
	}

	if err := session.SetKafkaContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to set kafka context: %w", err)
	}

	return session, nil
}

// SASLConfig holds SASL authentication configuration.
type SASLConfig struct {
	Mechanism string
	Username  string
	Password  string
}

// TLSConfig holds TLS configuration for Kafka connections.
type TLSConfig struct {
	Enabled        bool
	ClientKeyFile  string
	ClientCertFile string
	CaFile         string
}

// Initializes the necessary TLS configuration options
func tlsOpt(config *TLSConfig, opts []kgo.Opt) ([]kgo.Opt, error) {
	if config.Enabled {
		if config.CaFile != "" || config.ClientCertFile != "" || config.ClientKeyFile != "" {
			tc, err := tlscfg.New(
				tlscfg.MaybeWithDiskCA(config.CaFile, tlscfg.ForClient),
				tlscfg.MaybeWithDiskKeyPair(config.ClientCertFile, config.ClientKeyFile),
			)
			if err != nil {
				return nil, fmt.Errorf("unable to create TLS config: %v", err)
			}
			opts = append(opts, kgo.DialTLSConfig(tc))
		} else {
			opts = append(opts, kgo.DialTLSConfig(new(tls.Config)))
		}
	}
	return opts, nil
}

// Initializes the necessary SASL configuration options
func saslOpt(config *SASLConfig, opts []kgo.Opt) ([]kgo.Opt, error) {
	if config.Mechanism != "" || config.Username != "" || config.Password != "" {
		if config.Mechanism == "" || config.Username == "" || config.Password == "" {
			return nil, fmt.Errorf("all of Mechanism, Username, and Password must be specified if any are")
		}
		method := strings.ToLower(config.Mechanism)
		method = strings.ReplaceAll(method, "-", "")
		method = strings.ReplaceAll(method, "_", "")
		switch method {
		case "plain":
			opts = append(opts, kgo.SASL(plain.Auth{
				User: config.Username,
				Pass: config.Password,
			}.AsMechanism()))
		case "scramsha256":
			opts = append(opts, kgo.SASL(scram.Auth{
				User: config.Username,
				Pass: config.Password,
			}.AsSha256Mechanism()))
		case "scramsha512":
			opts = append(opts, kgo.SASL(scram.Auth{
				User: config.Username,
				Pass: config.Password,
			}.AsSha512Mechanism()))
		default:
			return nil, fmt.Errorf("unrecognized SASL method: %s", config.Mechanism)
		}
	}
	return opts, nil
}

// SetKafkaContext initializes Kafka clients using the provided context.
func (s *Session) SetKafkaContext(ctx KafkaContext) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.Ctx = ctx
	kc := &s.Ctx
	var err error
	s.Options = []kgo.Opt{}
	s.Options = append(s.Options, kgo.SeedBrokers(strings.Split(kc.BootstrapServers, ",")...))
	tlsConfig := &TLSConfig{
		Enabled:        kc.UseTLS,
		ClientKeyFile:  kc.ClientKeyFile,
		ClientCertFile: kc.ClientCertFile,
		CaFile:         kc.CaFile,
	}

	saslConfig := &SASLConfig{
		Mechanism: kc.AuthMechanism,
		Username:  kc.AuthUser,
		Password:  kc.AuthPass,
	}

	s.Options, err = tlsOpt(tlsConfig, s.Options)
	if err != nil {
		return fmt.Errorf("failed to create TLS config: %w", err)
	}
	s.Options, err = saslOpt(saslConfig, s.Options)
	if err != nil {
		return fmt.Errorf("failed to create SASL config: %w", err)
	}
	s.Options = append(s.Options, kgo.MaxVersions(kversion.V2_8_0()))

	s.Client, err = kgo.NewClient(
		s.Options...,
	)
	if err != nil {
		return fmt.Errorf("failed to create kafka client: %w", err)
	}

	s.AdminClient = kadm.NewClient(s.Client)
	if kc.SchemaRegistryURL != "" {
		SrOpts := []sr.ClientOpt{}
		SrOpts = append(SrOpts, sr.URLs(kc.SchemaRegistryURL))
		if kc.SchemaRegistryAuthUser != "" && kc.SchemaRegistryAuthPass != "" {
			SrOpts = append(SrOpts, sr.BasicAuth(kc.SchemaRegistryAuthUser, kc.SchemaRegistryAuthPass))
		} else if kc.SchemaRegistryBearerToken != "" {
			SrOpts = append(SrOpts, sr.BearerToken(kc.SchemaRegistryBearerToken))
		}
		SrOpts = append(SrOpts, sr.UserAgent("streamnative-mcp-server"))
		s.SchemaRegistryClient, err = sr.NewClient(SrOpts...)
		if err != nil {
			return fmt.Errorf("failed to create kafka schema registry client: %w", err)
		}
	}

	if kc.ConnectURL != "" {
		s.ConnectClient, err = NewConnect(kc)
		if err != nil {
			return fmt.Errorf("failed to create kafka connect client: %w", err)
		}
	}
	return nil
}

// ResetKafkaContext clears the current Kafka context and closes the data client.
func (s *Session) ResetKafkaContext() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.Client != nil {
		s.Client.Close()
	}

	s.Ctx = KafkaContext{}
	s.Client = nil
	s.AdminClient = nil
	s.SchemaRegistryClient = nil
	s.ConnectClient = nil
	s.Options = nil
}

// GetClient returns a Kafka client with optional overrides.
func (s *Session) GetClient(opts ...kgo.Opt) (*kgo.Client, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.Ctx.BootstrapServers == "" {
		return nil, fmt.Errorf("err: ContextNotSetErr: Please set the cluster context first")
	}

	if len(opts) > 0 {
		//nolint:gocritic
		clientOpts := append(s.Options, opts...)
		cli, err := kgo.NewClient(clientOpts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create kafka client with custom options: %w", err)
		}
		return cli, nil
	}

	if s.Client == nil {
		var err error
		s.Client, err = kgo.NewClient(s.Options...)
		if err != nil {
			return nil, fmt.Errorf("failed to create kafka client: %w", err)
		}
	}

	return s.Client, nil
}

// GetAdminClient returns the Kafka admin client.
func (s *Session) GetAdminClient() (*kadm.Client, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.Ctx.BootstrapServers == "" {
		return nil, fmt.Errorf("err: ContextNotSetErr: Please set the cluster context first")
	}

	if s.AdminClient == nil {
		if s.Client == nil {
			var err error
			s.Client, err = kgo.NewClient(s.Options...)
			if err != nil {
				return nil, fmt.Errorf("failed to create kafka client for admin: %w", err)
			}
		}
		s.AdminClient = kadm.NewClient(s.Client)
	}

	return s.AdminClient, nil
}

// GetSchemaRegistryClient returns the schema registry client.
func (s *Session) GetSchemaRegistryClient() (*sr.Client, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.Ctx.BootstrapServers == "" {
		return nil, fmt.Errorf("err: ContextNotSetErr: Please set the cluster context first")
	}
	if s.Ctx.SchemaRegistryURL == "" {
		return nil, fmt.Errorf("schema registry not enabled on the current context")
	}

	if s.SchemaRegistryClient == nil {
		SrOpts := []sr.ClientOpt{}
		SrOpts = append(SrOpts, sr.URLs(s.Ctx.SchemaRegistryURL))
		if s.Ctx.SchemaRegistryAuthUser != "" && s.Ctx.SchemaRegistryAuthPass != "" {
			SrOpts = append(SrOpts, sr.BasicAuth(s.Ctx.SchemaRegistryAuthUser, s.Ctx.SchemaRegistryAuthPass))
		} else if s.Ctx.SchemaRegistryBearerToken != "" {
			SrOpts = append(SrOpts, sr.BearerToken(s.Ctx.SchemaRegistryBearerToken))
		}
		SrOpts = append(SrOpts, sr.UserAgent("streamnative-mcp-server"))

		var err error
		s.SchemaRegistryClient, err = sr.NewClient(SrOpts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create kafka schema registry client: %w", err)
		}
	}

	return s.SchemaRegistryClient, nil
}

// GetConnectClient returns the Kafka Connect client.
func (s *Session) GetConnectClient() (Connect, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.Ctx.BootstrapServers == "" {
		return nil, fmt.Errorf("err: ContextNotSetErr: Please set the cluster context first")
	}
	if s.Ctx.ConnectURL == "" {
		return nil, fmt.Errorf("kafka connect not enabled on the current context")
	}

	if s.ConnectClient == nil {
		var err error
		s.ConnectClient, err = NewConnect(&s.Ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create kafka connect client: %w", err)
		}
	}

	return s.ConnectClient, nil
}
