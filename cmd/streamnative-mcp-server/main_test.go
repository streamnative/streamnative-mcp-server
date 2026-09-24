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

package main

import (
	"testing"

	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransportCommands(t *testing.T) {
	root := newRootCommand(config.NewConfigOptions())
	for _, name := range []string{"stdio", "sse", "http"} {
		cmd, args, err := root.Find([]string{name})
		require.NoError(t, err)
		require.Empty(t, args)
		require.Equal(t, name, cmd.Name())
	}

	httpCmd, _, err := root.Find([]string{"http"})
	require.NoError(t, err)
	addr := httpCmd.Flags().Lookup("http-addr")
	require.NotNil(t, addr)
	assert.Equal(t, "127.0.0.1:9090", addr.DefValue)
	assert.Equal(t, "127.0.0.1:9090", addr.Value.String())
	path := httpCmd.Flags().Lookup("http-path")
	require.NotNil(t, path)
	assert.Equal(t, "/mcp", path.DefValue)

	// Registering the new command must not silently change legacy SSE defaults.
	sseCmd, _, err := root.Find([]string{"sse"})
	require.NoError(t, err)
	sseAddr := sseCmd.Flags().Lookup("http-addr")
	require.NotNil(t, sseAddr)
	assert.Equal(t, ":9090", sseAddr.Value.String())
	require.NoError(t, httpCmd.Flags().Set("http-addr", "127.0.0.1:19090"))
	assert.Equal(t, ":9090", sseAddr.Value.String())
}
