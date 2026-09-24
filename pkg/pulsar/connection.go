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

// Package pulsar provides Pulsar connection helpers.
package pulsar

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/admin"
	pulsaradminauth "github.com/apache/pulsar-client-go/pulsaradmin/pkg/admin/auth"
	pulsaradminconfig "github.com/apache/pulsar-client-go/pulsaradmin/pkg/admin/config"
	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/rest"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	runtimeauth "github.com/streamnative/streamnative-mcp-server/pkg/auth"
	"golang.org/x/oauth2"
)

const (
	// DefaultClientTimeout is the default timeout for Pulsar client operations.
	DefaultClientTimeout = 30 * time.Second
)

// PulsarContext holds configuration for connecting to a Pulsar cluster.
type PulsarContext struct { //nolint:revive
	ServiceURL                    string
	WebServiceURL                 string
	Token                         string
	TokenSource                   oauth2.TokenSource `json:"-"`
	AuthPlugin                    string
	AuthParams                    string
	TLSAllowInsecureConnection    bool
	TLSEnableHostnameVerification bool
	TLSTrustCertsFilePath         string
	TLSCertFile                   string
	TLSKeyFile                    string
}

// Session represents a Pulsar session
type Session struct {
	Ctx             PulsarContext
	Client          pulsar.Client
	AdminClient     cmdutils.Client
	AdminV3Client   cmdutils.Client
	ClientOptions   pulsar.ClientOptions
	PulsarCtlConfig *cmdutils.ClusterConfig
	mutex           sync.RWMutex
	adminStatusREST *rest.Client
}

// NewSession creates a new Pulsar session with the given context
// This function dynamically constructs clients without relying on global state
func NewSession(ctx PulsarContext) (*Session, error) {
	session := &Session{
		Ctx: ctx,
	}

	if err := session.SetPulsarContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to set pulsar context: %w", err)
	}

	return session, nil
}

// SetPulsarContext initializes Pulsar clients using the provided context.
func (s *Session) SetPulsarContext(ctx PulsarContext) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Ctx = ctx
	s.adminStatusREST = nil
	pc := &s.Ctx
	var err error
	// Configure pulsarctl with the token
	switch {
	case pc.Token != "":
		s.PulsarCtlConfig = &cmdutils.ClusterConfig{
			WebServiceURL:                 pc.WebServiceURL,
			AuthPlugin:                    "org.apache.pulsar.client.impl.auth.AuthenticationToken",
			AuthParams:                    fmt.Sprintf("token:%s", pc.Token),
			TLSAllowInsecureConnection:    pc.TLSAllowInsecureConnection,
			TLSEnableHostnameVerification: pc.TLSEnableHostnameVerification,
			TLSTrustCertsFilePath:         pc.TLSTrustCertsFilePath,
			TLSCertFile:                   pc.TLSCertFile,
			TLSKeyFile:                    pc.TLSKeyFile,
		}

		// Set the client options
		s.ClientOptions = pulsar.ClientOptions{
			URL:                        pc.ServiceURL,
			Authentication:             pulsar.NewAuthenticationToken(pc.Token),
			OperationTimeout:           DefaultClientTimeout,
			ConnectionTimeout:          DefaultClientTimeout,
			TLSAllowInsecureConnection: pc.TLSAllowInsecureConnection,
			TLSValidateHostname:        pc.TLSEnableHostnameVerification,
			TLSTrustCertsFilePath:      pc.TLSTrustCertsFilePath,
			TLSCertificateFile:         pc.TLSCertFile,
			TLSKeyFilePath:             pc.TLSKeyFile,
		}
	case pc.AuthPlugin != "" && pc.AuthParams != "":
		s.PulsarCtlConfig = &cmdutils.ClusterConfig{
			WebServiceURL:                 pc.WebServiceURL,
			AuthPlugin:                    pc.AuthPlugin,
			AuthParams:                    pc.AuthParams,
			TLSAllowInsecureConnection:    pc.TLSAllowInsecureConnection,
			TLSEnableHostnameVerification: pc.TLSEnableHostnameVerification,
			TLSTrustCertsFilePath:         pc.TLSTrustCertsFilePath,
			TLSCertFile:                   pc.TLSCertFile,
			TLSKeyFile:                    pc.TLSKeyFile,
		}

		authProvider, err := pulsar.NewAuthentication(pc.AuthPlugin, pc.AuthParams)
		if err != nil {
			return fmt.Errorf("failed to create authentication provider: %w", err)
		}
		s.ClientOptions = pulsar.ClientOptions{
			URL:                        pc.ServiceURL,
			Authentication:             authProvider,
			OperationTimeout:           DefaultClientTimeout,
			ConnectionTimeout:          DefaultClientTimeout,
			TLSAllowInsecureConnection: pc.TLSAllowInsecureConnection,
			TLSValidateHostname:        pc.TLSEnableHostnameVerification,
			TLSTrustCertsFilePath:      pc.TLSTrustCertsFilePath,
			TLSCertificateFile:         pc.TLSCertFile,
			TLSKeyFilePath:             pc.TLSKeyFile,
		}
	default:
		// No authentication provided
		s.PulsarCtlConfig = &cmdutils.ClusterConfig{
			WebServiceURL:                 pc.WebServiceURL,
			TLSAllowInsecureConnection:    pc.TLSAllowInsecureConnection,
			TLSEnableHostnameVerification: pc.TLSEnableHostnameVerification,
			TLSTrustCertsFilePath:         pc.TLSTrustCertsFilePath,
			TLSCertFile:                   pc.TLSCertFile,
			TLSKeyFile:                    pc.TLSKeyFile,
		}

		// Set the client options without authentication
		s.ClientOptions = pulsar.ClientOptions{
			URL:                        pc.ServiceURL,
			OperationTimeout:           DefaultClientTimeout,
			ConnectionTimeout:          DefaultClientTimeout,
			TLSAllowInsecureConnection: pc.TLSAllowInsecureConnection,
			TLSValidateHostname:        pc.TLSEnableHostnameVerification,
			TLSTrustCertsFilePath:      pc.TLSTrustCertsFilePath,
			TLSCertificateFile:         pc.TLSCertFile,
			TLSKeyFilePath:             pc.TLSKeyFile,
		}
	}

	if pc.TokenSource != nil {
		source := pc.TokenSource
		s.ClientOptions.Authentication = pulsar.NewAuthenticationTokenFromSupplier(func() (string, error) {
			token, err := source.Token()
			if err != nil {
				return "", err
			}
			return token.AccessToken, nil
		})
	}
	s.AdminClient = s.PulsarCtlConfig.Client(pulsaradminconfig.V2)
	s.AdminV3Client = s.PulsarCtlConfig.Client(pulsaradminconfig.V3)
	if pc.TokenSource != nil {
		s.AdminClient, err = newRuntimeAdminClient(s.PulsarCtlConfig, pc.TokenSource, pulsaradminconfig.V2, s.AdminClient)
		if err != nil {
			return err
		}
		s.AdminV3Client, err = newRuntimeAdminClient(s.PulsarCtlConfig, pc.TokenSource, pulsaradminconfig.V3, s.AdminV3Client)
		if err != nil {
			return err
		}
	}

	s.Client, err = pulsar.NewClient(s.ClientOptions)
	if err != nil {
		return fmt.Errorf("failed to create pulsar client: %w", err)
	}

	return nil
}

// ResetPulsarContext clears the current Pulsar context and closes the data client.
func (s *Session) ResetPulsarContext() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.Client != nil {
		s.Client.Close()
	}

	s.Ctx = PulsarContext{}
	s.Client = nil
	s.AdminClient = nil
	s.AdminV3Client = nil
	s.ClientOptions = pulsar.ClientOptions{}
	s.PulsarCtlConfig = nil
	s.adminStatusREST = nil
}

// GetAdminClient returns the Pulsar admin v2 client.
func (s *Session) GetAdminClient() (cmdutils.Client, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.PulsarCtlConfig == nil || s.PulsarCtlConfig.WebServiceURL == "" {
		return nil, fmt.Errorf("err: ContextNotSetErr: Please set the cluster context first")
	}
	return s.AdminClient, nil
}

// GetAdminV3Client returns the Pulsar admin v3 client.
func (s *Session) GetAdminV3Client() (cmdutils.Client, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.PulsarCtlConfig == nil || s.PulsarCtlConfig.WebServiceURL == "" {
		return nil, fmt.Errorf("err: ContextNotSetErr: Please set the cluster context first")
	}
	return s.AdminV3Client, nil
}

// GetPulsarCtlConfig returns a copy of the underlying pulsarctl configuration.
func (s *Session) GetPulsarCtlConfig() (*cmdutils.ClusterConfig, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.PulsarCtlConfig == nil || s.PulsarCtlConfig.WebServiceURL == "" {
		return nil, fmt.Errorf("err: ContextNotSetErr: Please set the cluster context first")
	}

	cfg := *s.PulsarCtlConfig
	return &cfg, nil
}

// GetAdminStatusClient returns a cached REST client for the Pulsar status endpoint.
func (s *Session) GetAdminStatusClient() (*rest.Client, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.PulsarCtlConfig == nil || s.PulsarCtlConfig.WebServiceURL == "" {
		return nil, fmt.Errorf("err: ContextNotSetErr: Please set the cluster context first")
	}
	if s.adminStatusREST != nil {
		return s.adminStatusREST, nil
	}

	cfg := *s.PulsarCtlConfig
	authProvider, err := pulsaradminauth.GetAuthProvider((*pulsaradminconfig.Config)(&cfg))
	if err != nil {
		return nil, fmt.Errorf("failed to build status auth provider: %w", err)
	}
	if s.Ctx.TokenSource != nil {
		authProvider = &runtimeauth.TokenTransport{Source: s.Ctx.TokenSource, Base: authProvider.Transport()}
	}

	s.adminStatusREST = &rest.Client{
		ServiceURL:  cfg.WebServiceURL,
		VersionInfo: admin.ReleaseVersion,
		HTTPClient: &http.Client{
			Timeout:   admin.DefaultHTTPTimeOutDuration,
			Transport: authProvider,
		},
	}
	return s.adminStatusREST, nil
}

// GetPulsarClient returns the Pulsar data client.
func (s *Session) GetPulsarClient() (pulsar.Client, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.ClientOptions.URL == "" {
		return nil, fmt.Errorf("err: ContextNotSetErr: Please set the cluster context first")
	}
	return s.Client, nil
}

type runtimeAdminClient struct {
	admin.Client
	tokenClient cmdutils.Client
}

func (c *runtimeAdminClient) Token() cmdutils.Token { return c.tokenClient.Token() }

func newRuntimeAdminClient(cfg *cmdutils.ClusterConfig, source oauth2.TokenSource,
	version pulsaradminconfig.APIVersion, tokenClient cmdutils.Client,
) (cmdutils.Client, error) {
	configuration := pulsaradminconfig.Config(*cfg)
	configuration.PulsarAPIVersion = version
	base, err := pulsaradminauth.NewDefaultTransport(&configuration)
	if err != nil {
		return nil, err
	}
	client, err := admin.NewPulsarClientWithAuthProvider(&configuration, &runtimeauth.TokenTransport{Source: source, Base: base})
	if err != nil {
		return nil, err
	}
	return &runtimeAdminClient{Client: client, tokenClient: tokenClient}, nil
}
