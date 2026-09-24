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

package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/streamnative/streamnative-mcp-server/pkg/auth"
	"github.com/streamnative/streamnative-mcp-server/pkg/auth/store"
	"github.com/stretchr/testify/require"
)

func TestStandaloneKeyFileUsesSessionLocalGrants(t *testing.T) {
	viper.Set("config-dir", t.TempDir())
	t.Cleanup(viper.Reset)
	create := func() *Options {
		o := NewConfigOptions()
		o.BackendOverride = "file"
		o.IssuerEndpoint = "https://issuer.example.test/"
		o.KeyFile = `data://{"type":"sn_service_account","client_id":"subject","client_secret":"fake-secret"}`
		require.NoError(t, o.Complete())
		return o
	}
	a, b := create(), create()
	require.NoError(t, a.SaveGrant("audience", auth.AuthorizationGrant{
		Type: auth.GrantTypeClientCredentials, ClientCredentials: &auth.KeyFile{ClientID: "subject"},
	}))
	_, err := b.LoadGrant("audience")
	require.ErrorIs(t, err, store.ErrNoAuthenticationData, "standalone sessions must not share another process's grants")
}

func TestStandaloneRejectsMalformedInlineCredentialsWithoutLeakingThem(t *testing.T) {
	o := NewConfigOptions()
	o.ConfigDir = t.TempDir()
	o.KeyFile = `data://{"type":"unsupported","client_id":"test","client_secret":"never-echo-me"}`
	err := o.Complete()
	require.Error(t, err)
	require.NotContains(t, err.Error(), "never-echo-me")
}
