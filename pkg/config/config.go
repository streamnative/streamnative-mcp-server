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
	"errors"
	"net/url"

	"github.com/streamnative/streamnative-mcp-server/pkg/auth"
)

type SnConfig struct {
	// the API server endpoint
	Server string `yaml:"server"`
	// CA bundle (base64, PEM)
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
	// indicates whether to skip TLS verification
	InsecureSkipTLSVerify bool `yaml:"insecure-skip-tls-verify"`
	// user auth information
	Auth Auth `yaml:"auth"`
	// settable context
	Context       Context `yaml:"context"`
	ProxyLocation string  `yaml:"proxy-location"`
	KeyFile       string  `yaml:"key-file"`

	ExternalKafka  *ExternalKafka  `yaml:"external-kafka"`
	ExternalPulsar *ExternalPulsar `yaml:"external-pulsar"`
}

type Auth struct {
	// the OAuth 2.0 issuer endpoint
	IssuerEndpoint string `yaml:"issuer"`
	// the audience identifier for the API server (default: server URL)
	Audience string `yaml:"audience"`
	// the client ID to use for authorization grants (note: not used for service accounts)
	ClientID string `yaml:"client-id"`
}

func (a *Auth) Validate() error {
	if !(isValidIssuer(a.IssuerEndpoint) && isValidClientID(a.ClientID) && isValidAudience(a.Audience)) {
		return errors.New("configuration error: auth section is incomplete or invalid")
	}
	return nil
}

func isValidIssuer(iss string) bool {
	u, err := url.Parse(iss)
	return err == nil && iss != "" && u.IsAbs()
}

func isValidClientID(cid string) bool {
	return cid != ""
}

func isValidAudience(aud string) bool {
	return aud != ""
}

func (a *Auth) Issuer() auth.Issuer {
	return auth.Issuer{
		IssuerEndpoint: a.IssuerEndpoint,
		ClientID:       a.ClientID,
		Audience:       a.Audience,
	}
}

type Context struct {
	Organization   string `yaml:"organization,omitempty"`
	PulsarInstance string `yaml:"pulsar-instance,omitempty"`
	PulsarCluster  string `yaml:"pulsar-cluster,omitempty"`
}

type Storage interface {
	// Gets the config directory for configuration files, credentials and caches
	GetConfigDirectory() string

	// LoadConfig loads the raw configuration from storage.
	LoadConfig() (*SnConfig, error)

	// LoadConfigOrDie loads the raw configuration from storage, or dies if unable to.
	LoadConfigOrDie() *SnConfig

	// SaveConfig saves the given configuration to storage, overwriting any previous configuration.
	SaveConfig(config *SnConfig) error
}
