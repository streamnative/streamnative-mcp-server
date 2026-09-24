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

package mcp

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/99designs/keyring"
	"github.com/sirupsen/logrus"
	"github.com/streamnative/streamnative-mcp-server/pkg/auth"
	"github.com/streamnative/streamnative-mcp-server/pkg/auth/store"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestInitialClusterDoesNotImplicitlyLockContext(t *testing.T) {
	issuer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"token_endpoint":"http://127.0.0.1:1/token"}`)
	}))
	defer issuer.Close()
	for _, scenario := range []string{"initial", "locked", "read-only"} {
		t.Run(scenario, func(t *testing.T) {
			o := config.NewConfigOptions()
			o.KeyFile, o.Audience, o.IssuerEndpoint = "unused", "cloud-api", issuer.URL+"/"
			o.Organization, o.PulsarInstance, o.PulsarCluster = "org", "instance", "cluster"
			var err error
			o.Store, err = store.NewKeyringStore(keyring.NewArrayKeyring(nil))
			require.NoError(t, err)
			require.NoError(t, o.SaveGrant("cloud-api", auth.AuthorizationGrant{
				Type: auth.GrantTypeClientCredentials, ClientCredentials: &auth.KeyFile{ClientID: "sa", ClientEmail: "sa@example.test"},
			}))
			srv, err := newMcpServer(context.Background(), &ServerOptions{
				Options: o, Features: []string{"streamnative-cloud"},
				ReadOnly: scenario == "read-only", LockClusterContext: scenario == "locked",
			}, logrus.New())
			require.NoError(t, err)
			defer func() { _ = srv.SNCloudSession.Close() }()
			for _, name := range []string{"sncloud_context_use_cluster", "sncloud_context_reset"} {
				if scenario == "initial" {
					require.NotNil(t, srv.MCPServer.GetTool(name))
				} else {
					require.Nil(t, srv.MCPServer.GetTool(name))
				}
			}
		})
	}
}
