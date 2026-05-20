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

package pulsar

import (
	oslib "os"
	path "path/filepath"
	"testing"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func TestParseFunctionIdentity(t *testing.T) {
	builder := NewPulsarAdminFunctionsToolBuilder()

	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{}}}
	identity, err := builder.parseFunctionIdentity(req, "list")
	require.NoError(t, err)
	require.Equal(t, defaultTenant, identity.Tenant)
	require.Equal(t, defaultNamespace, identity.Namespace)

	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"fqfn":      "t/ns/name",
		"tenant":    "t",
		"namespace": "ns",
	}}}
	_, err = builder.parseFunctionIdentity(req, "get")
	require.Error(t, err)

	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"fqfn": "tenant//name",
	}}}
	_, err = builder.parseFunctionIdentity(req, "get")
	require.Error(t, err)

	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"fqfn": "//name",
	}}}
	_, err = builder.parseFunctionIdentity(req, "get")
	require.Error(t, err)

	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"tenant": "t",
	}}}
	_, err = builder.parseFunctionIdentity(req, "get")
	require.Error(t, err)

	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"tenant":    "t",
		"namespace": "ns",
	}}}
	identity, err = builder.parseFunctionIdentity(req, "update")
	require.NoError(t, err)
	require.Equal(t, "t", identity.Tenant)
	require.Equal(t, "ns", identity.Namespace)
}

func TestFunctionOperationCoverageIncludesFileTransfer(t *testing.T) {
	require.True(t, isSupportedFunctionOperation("download"))
	require.True(t, isSupportedFunctionOperation("upload"))
	require.False(t, isReadOnlyRestrictedFunctionOperation("download"))
	require.True(t, isReadOnlyRestrictedFunctionOperation("upload"))
}

func TestParseFunctionDownloadTarget(t *testing.T) {
	builder := NewPulsarAdminFunctionsToolBuilder()

	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"destinationFile": "/tmp/function.nar",
		"path":            "public/default/function.nar",
	}}}
	target, err := builder.parseFunctionDownloadTarget(req)
	require.NoError(t, err)
	require.True(t, target.UsePath)
	require.Equal(t, "/tmp/function.nar", target.DestinationFile)
	require.Equal(t, "public/default/function.nar", target.Path)

	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"destinationFile": "/tmp/function.nar",
		"fqfn":            "tenant/ns/name",
	}}}
	target, err = builder.parseFunctionDownloadTarget(req)
	require.NoError(t, err)
	require.False(t, target.UsePath)
	require.Equal(t, "/tmp/function.nar", target.DestinationFile)
	require.Equal(t, "tenant", target.Identity.Tenant)
	require.Equal(t, "ns", target.Identity.Namespace)
	require.Equal(t, "name", target.Identity.Name)
}

func TestParseFunctionDownloadTargetErrors(t *testing.T) {
	builder := NewPulsarAdminFunctionsToolBuilder()

	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"path": "public/default/function.nar",
	}}}
	_, err := builder.parseFunctionDownloadTarget(req)
	require.Error(t, err)

	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"destinationFile": "/tmp/function.nar",
	}}}
	_, err = builder.parseFunctionDownloadTarget(req)
	require.Error(t, err)
}

func TestValidateFunctionConfigs(t *testing.T) {
	jar := "builtin://identity"
	config := &utils.FunctionConfig{
		Jar:       &jar,
		ClassName: "org.example.Identity",
	}
	require.NoError(t, validateFunctionConfigs(config))

	require.Error(t, validateFunctionConfigs(&utils.FunctionConfig{}))

	pyDir := t.TempDir()
	pyFile := path.Join(pyDir, "echo.py")
	require.NoError(t, oslib.WriteFile(pyFile, []byte("print('ok')"), 0o600))
	config = &utils.FunctionConfig{
		Py:        &pyFile,
		ClassName: "echo.EchoFunction",
	}
	require.NoError(t, validateFunctionConfigs(config))

	py := "a.py"
	config = &utils.FunctionConfig{
		Jar: &jar,
		Py:  &py,
	}
	require.Error(t, validateFunctionConfigs(config))
}

func TestCheckArgsForUpdate(t *testing.T) {
	config := &utils.FunctionConfig{
		ClassName: "org.example.Example",
	}
	require.NoError(t, checkArgsForUpdate(config))
	require.NotEmpty(t, config.Name)
	require.Equal(t, defaultTenant, config.Tenant)
	require.Equal(t, defaultNamespace, config.Namespace)
}

func TestBuildFunctionConfigMutualExclusion(t *testing.T) {
	builder := NewPulsarAdminFunctionsToolBuilder()

	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"functionType": "identity",
		"jar":          "builtin://identity",
	}}}
	_, err := builder.buildFunctionConfig("public", "default", "name", req)
	require.Error(t, err)

	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"functionType": "identity",
		"py":           "echo.py",
	}}}
	_, err = builder.buildFunctionConfig("public", "default", "name", req)
	require.Error(t, err)

	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"functionType": "identity",
		"go":           "echo",
	}}}
	_, err = builder.buildFunctionConfig("public", "default", "name", req)
	require.Error(t, err)

	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{
		"jar": "builtin://identity",
		"py":  "echo.py",
	}}}
	_, err = builder.buildFunctionConfig("public", "default", "name", req)
	require.Error(t, err)
}
