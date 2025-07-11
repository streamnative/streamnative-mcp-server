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

package pulsar

import (
	"fmt"
	"sync"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	pulsaradminconfig "github.com/apache/pulsar-client-go/pulsaradmin/pkg/admin/config"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/streamnative/streamnative-mcp-server/pkg/auth"
	"github.com/streamnative/streamnative-mcp-server/pkg/auth/store"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
)

const (
	DefaultClientTimeout = 30 * time.Second
)

//nolint:revive
type PulsarContext struct {
	ServiceURL                    string
	WebServiceURL                 string
	Token                         string
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
	Ctx           PulsarContext
	Client        pulsar.Client
	AdminClient   cmdutils.Client
	AdminV3Client cmdutils.Client
	ClientOptions pulsar.ClientOptions
	mutex         sync.RWMutex
}

func init() {
	cmdutils.PulsarCtlConfig = &cmdutils.ClusterConfig{}
}

// NewSession creates a new Pulsar session with the given context
// This function dynamically constructs clients without relying on global state
func NewSession(ctx PulsarContext, issuer *auth.Issuer, tokenStore *store.Store) (*Session, error) {
	session := &Session{
		Ctx: ctx,
	}

	if issuer != nil && tokenStore != nil {
		_ = config.InitSNCloudLogClient(*issuer, *tokenStore)
	} else {
		config.ResetSNCloudLogClient()
	}

	if err := session.SetPulsarContext(); err != nil {
		return nil, fmt.Errorf("failed to set pulsar context: %w", err)
	}

	return session, nil
}

func (s *Session) ChangeContext(ctx PulsarContext, issuer *auth.Issuer, tokenStore *store.Store) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.Ctx = ctx
	if issuer != nil && tokenStore != nil {
		_ = config.InitSNCloudLogClient(*issuer, *tokenStore)
	} else {
		config.ResetSNCloudLogClient()
	}
	return s.SetPulsarContext()
}

func (s *Session) SetPulsarContext() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	pc := &s.Ctx
	var err error
	// Configure pulsarctl with the token
	switch {
	case pc.Token != "":
		cmdutils.PulsarCtlConfig = &cmdutils.ClusterConfig{
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
		cmdutils.PulsarCtlConfig = &cmdutils.ClusterConfig{
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
		cmdutils.PulsarCtlConfig = &cmdutils.ClusterConfig{
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

	s.AdminClient = cmdutils.NewPulsarClient()
	s.AdminV3Client = cmdutils.NewPulsarClientWithAPIVersion(pulsaradminconfig.V3)

	s.Client, err = pulsar.NewClient(s.ClientOptions)
	if err != nil {
		return fmt.Errorf("failed to create pulsar client: %w", err)
	}

	return nil
}

func (s *Session) GetAdminClient() (cmdutils.Client, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if cmdutils.PulsarCtlConfig.WebServiceURL == "" {
		return nil, fmt.Errorf("err: ContextNotSetErr: Please set the cluster context first")
	}
	return s.AdminClient, nil
}

func (s *Session) GetAdminV3Client() (cmdutils.Client, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if cmdutils.PulsarCtlConfig.WebServiceURL == "" {
		return nil, fmt.Errorf("err: ContextNotSetErr: Please set the cluster context first")
	}
	return s.AdminV3Client, nil
}

func (s *Session) GetPulsarClient() (pulsar.Client, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.ClientOptions.URL == "" {
		return nil, fmt.Errorf("err: ContextNotSetErr: Please set the cluster context first")
	}
	return s.Client, nil
}
