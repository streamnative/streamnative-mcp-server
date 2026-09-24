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

package kafka

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl"
	"golang.org/x/oauth2"
)

const contextNotSetErr = "err: ContextNotSetErr: Please set the cluster context first"

type rotatingTokenSource struct{ token atomic.Value }

func (s *rotatingTokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{AccessToken: s.token.Load().(string)}, nil
}

func TestKafkaClientsReadCurrentTokenAtAuthenticationTime(t *testing.T) {
	source := &rotatingTokenSource{}
	source.token.Store("old-token")
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, password, ok := r.BasicAuth()
		if !ok || user != "public/default" || password != "new-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/" {
			_, _ = fmt.Fprint(w, "{}")
			return
		}
		_, _ = fmt.Fprint(w, "[]")
	}))
	defer api.Close()
	session, err := NewSession(KafkaContext{
		BootstrapServers: "localhost:9093", AuthMechanism: "PLAIN", AuthUser: "public/default", AuthPass: "token:old-token",
		SchemaRegistryURL: api.URL, SchemaRegistryAuthUser: "public/default", SchemaRegistryAuthPass: "old-token",
		ConnectURL: api.URL, ConnectAuthUser: "public/default", ConnectAuthPass: "old-token",
		TokenSource: source, TokenPrefix: "token:",
	})
	require.NoError(t, err)
	defer session.ResetKafkaContext()
	source.token.Store("new-token")
	client, err := session.GetClient()
	require.NoError(t, err)
	mechanisms := client.OptValue(kgo.SASL).([]sasl.Mechanism)
	_, data, err := mechanisms[0].Authenticate(context.Background(), "localhost")
	require.NoError(t, err)
	require.Equal(t, "\x00public/default\x00token:new-token", string(data))
	schema, err := session.GetSchemaRegistryClient()
	require.NoError(t, err)
	_, err = schema.Subjects(context.Background())
	require.NoError(t, err)
	connect, err := session.GetConnectClient()
	require.NoError(t, err)
	_, err = connect.GetInfo(context.Background())
	require.NoError(t, err)
}

func TestSessionGetKafkaClientsRequireContext(t *testing.T) {
	session := &Session{}

	_, err := session.GetClient()
	require.EqualError(t, err, contextNotSetErr)

	_, err = session.GetAdminClient()
	require.EqualError(t, err, contextNotSetErr)

	_, err = session.GetSchemaRegistryClient()
	require.EqualError(t, err, contextNotSetErr)

	_, err = session.GetConnectClient()
	require.EqualError(t, err, contextNotSetErr)
}

func TestSessionResetKafkaContextClearsContextAndClients(t *testing.T) {
	session := &Session{
		Ctx: KafkaContext{
			BootstrapServers:  "localhost:9092",
			SchemaRegistryURL: "http://localhost:8081",
			ConnectURL:        "http://localhost:8083",
		},
	}

	session.ResetKafkaContext()

	require.Equal(t, KafkaContext{}, session.Ctx)
	require.Nil(t, session.Client)
	require.Nil(t, session.AdminClient)
	require.Nil(t, session.SchemaRegistryClient)
	require.Nil(t, session.ConnectClient)
	require.Nil(t, session.Options)

	_, err := session.GetAdminClient()
	require.EqualError(t, err, contextNotSetErr)
}
