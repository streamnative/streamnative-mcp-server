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
	"context"
	"testing"

	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPulsarAdminNsIsolationPolicyToolBuilder_FeatureGate(t *testing.T) {
	builder := NewPulsarAdminNsIsolationPolicyToolBuilder()

	assert.Contains(t, builder.GetRequiredFeatures(), "pulsar-admin-ns-isolation-policy")

	tools, err := builder.BuildTools(context.Background(), builders.ToolBuildConfig{
		Features: []string{"pulsar-admin-ns-isolation-policy"},
	})
	require.NoError(t, err)
	require.Len(t, tools, 1)
	assert.Equal(t, "pulsar_admin_nsisolationpolicy", tools[0].Tool.Name)
}

func TestPulsarAdminResourceQuotasToolBuilder_FeatureGate(t *testing.T) {
	builder := NewPulsarAdminResourceQuotasToolBuilder()

	assert.Contains(t, builder.GetRequiredFeatures(), "pulsar-admin-resource-quotas")

	tools, err := builder.BuildTools(context.Background(), builders.ToolBuildConfig{
		Features: []string{"pulsar-admin-resource-quotas"},
	})
	require.NoError(t, err)
	require.Len(t, tools, 1)
	assert.Equal(t, "pulsar_admin_resourcequota", tools[0].Tool.Name)
}
