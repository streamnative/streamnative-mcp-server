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
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	runtimeauth "github.com/streamnative/streamnative-mcp-server/pkg/auth"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kversion"
	"github.com/twmb/franz-go/pkg/sasl/plain"
	"github.com/twmb/franz-go/pkg/sasl/scram"
	"github.com/twmb/franz-go/pkg/sr"
	"github.com/twmb/tlscfg"
	"golang.org/x/oauth2"
)

//nolint:revive
type KafkaContext struct {
	TokenSource       oauth2.TokenSource `json:"-"`
	TokenPrefix       string
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
	if kc.TokenSource == nil {
		s.Options, err = saslOpt(saslConfig, s.Options)
		if err != nil {
			return fmt.Errorf("failed to create SASL config: %w", err)
		}
	} else {
		if !strings.EqualFold(kc.AuthMechanism, "PLAIN") {
			return fmt.Errorf("runtime token authentication requires SASL PLAIN")
		}
		source, username, prefix := kc.TokenSource, kc.AuthUser, kc.TokenPrefix
		s.Options = append(s.Options, kgo.SASL(plain.Plain(func(context.Context) (plain.Auth, error) {
			token, err := source.Token()
			if err != nil {
				return plain.Auth{}, err
			}
			return plain.Auth{User: username, Pass: prefix + token.AccessToken}, nil
		})))
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
		s.SchemaRegistryClient, err = newSchemaRegistryClient(*kc)
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
		var err error
		s.SchemaRegistryClient, err = newSchemaRegistryClient(s.Ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create kafka schema registry client: %w", err)
		}
	}

	return s.SchemaRegistryClient, nil
}

func newSchemaRegistryClient(kc KafkaContext) (*sr.Client, error) {
	opts := []sr.ClientOpt{sr.URLs(kc.SchemaRegistryURL), sr.UserAgent("streamnative-mcp-server")}
	switch {
	case kc.TokenSource != nil:
		// Preserve franz-go's five-second default HTTP timeout.
		opts = append(opts, sr.HTTPClient(&http.Client{Timeout: 5 * time.Second, Transport: &runtimeauth.TokenTransport{
			Source: kc.TokenSource, Username: kc.SchemaRegistryAuthUser,
		}}))
	case kc.SchemaRegistryAuthUser != "" && kc.SchemaRegistryAuthPass != "":
		opts = append(opts, sr.BasicAuth(kc.SchemaRegistryAuthUser, kc.SchemaRegistryAuthPass))
	case kc.SchemaRegistryBearerToken != "":
		opts = append(opts, sr.BearerToken(kc.SchemaRegistryBearerToken))
	}
	return sr.NewClient(opts...)
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
