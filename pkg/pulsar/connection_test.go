// Copyright 2025 StreamNative
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

package pulsar

import (
	"testing"

	pulsarclient "github.com/apache/pulsar-client-go/pulsar"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	"github.com/stretchr/testify/require"
)

const contextNotSetErr = "err: ContextNotSetErr: Please set the cluster context first"

func TestSessionGetPulsarCtlConfigRequiresWebServiceURL(t *testing.T) {
	session := &Session{}

	_, err := session.GetPulsarCtlConfig()
	require.EqualError(t, err, contextNotSetErr)

	session.PulsarCtlConfig = &cmdutils.ClusterConfig{}

	_, err = session.GetPulsarCtlConfig()
	require.EqualError(t, err, contextNotSetErr)
}

func TestSessionGetPulsarCtlConfigReturnsCopy(t *testing.T) {
	session := &Session{
		PulsarCtlConfig: &cmdutils.ClusterConfig{
			WebServiceURL: "http://pulsar.example.com:8080",
		},
	}

	cfg, err := session.GetPulsarCtlConfig()
	require.NoError(t, err)
	require.Equal(t, "http://pulsar.example.com:8080", cfg.WebServiceURL)

	cfg.WebServiceURL = "http://mutated.example.com:8080"

	require.Equal(t, "http://pulsar.example.com:8080", session.PulsarCtlConfig.WebServiceURL)
}

func TestSessionGetAdminStatusClientRequiresWebServiceURL(t *testing.T) {
	session := &Session{}

	_, err := session.GetAdminStatusClient()
	require.EqualError(t, err, contextNotSetErr)

	session.PulsarCtlConfig = &cmdutils.ClusterConfig{}

	_, err = session.GetAdminStatusClient()
	require.EqualError(t, err, contextNotSetErr)
}

func TestSessionGetAdminStatusClientReturnsCachedClient(t *testing.T) {
	session := &Session{
		PulsarCtlConfig: &cmdutils.ClusterConfig{
			WebServiceURL: "http://pulsar.example.com:8080",
		},
	}

	firstClient, err := session.GetAdminStatusClient()
	require.NoError(t, err)
	require.Equal(t, "http://pulsar.example.com:8080", firstClient.ServiceURL)

	secondClient, err := session.GetAdminStatusClient()
	require.NoError(t, err)
	require.Same(t, firstClient, secondClient)
}

func TestSessionGetAdminClientsRequireWebServiceURL(t *testing.T) {
	session := &Session{}

	_, err := session.GetAdminClient()
	require.EqualError(t, err, contextNotSetErr)

	_, err = session.GetAdminV3Client()
	require.EqualError(t, err, contextNotSetErr)

	session.PulsarCtlConfig = &cmdutils.ClusterConfig{}

	_, err = session.GetAdminClient()
	require.EqualError(t, err, contextNotSetErr)

	_, err = session.GetAdminV3Client()
	require.EqualError(t, err, contextNotSetErr)
}

func TestSessionResetPulsarContextClearsContextAndClients(t *testing.T) {
	session := &Session{
		Ctx: PulsarContext{
			ServiceURL:    "pulsar://pulsar.example.com:6650",
			WebServiceURL: "http://pulsar.example.com:8080",
			Token:         "token",
		},
		PulsarCtlConfig: &cmdutils.ClusterConfig{
			WebServiceURL: "http://pulsar.example.com:8080",
		},
		ClientOptions: pulsarclient.ClientOptions{URL: "pulsar://pulsar.example.com:6650"},
	}

	session.ResetPulsarContext()

	require.Equal(t, PulsarContext{}, session.Ctx)
	require.Nil(t, session.Client)
	require.Nil(t, session.AdminClient)
	require.Nil(t, session.AdminV3Client)
	require.Equal(t, pulsarclient.ClientOptions{}, session.ClientOptions)
	require.Nil(t, session.PulsarCtlConfig)
	require.Nil(t, session.adminStatusREST)

	_, err := session.GetAdminClient()
	require.EqualError(t, err, contextNotSetErr)

	_, err = session.GetPulsarClient()
	require.EqualError(t, err, contextNotSetErr)
}

func TestSessionSetPulsarContextResetsAdminStatusClient(t *testing.T) {
	session := &Session{
		PulsarCtlConfig: &cmdutils.ClusterConfig{
			WebServiceURL: "http://pulsar.example.com:8080",
		},
	}

	firstClient, err := session.GetAdminStatusClient()
	require.NoError(t, err)

	err = session.SetPulsarContext(PulsarContext{
		ServiceURL:    "pulsar://localhost:6650",
		WebServiceURL: "http://pulsar-next.example.com:8080",
	})
	require.NoError(t, err)
	if session.Client != nil {
		defer session.Client.Close()
	}
	require.Nil(t, session.adminStatusREST)

	secondClient, err := session.GetAdminStatusClient()
	require.NoError(t, err)
	require.NotSame(t, firstClient, secondClient)
	require.Equal(t, "http://pulsar-next.example.com:8080", secondClient.ServiceURL)
}
