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
	"testing"

	"github.com/stretchr/testify/require"
)

const contextNotSetErr = "err: ContextNotSetErr: Please set the cluster context first"

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
