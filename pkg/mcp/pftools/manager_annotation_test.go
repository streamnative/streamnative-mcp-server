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

package pftools

import (
	"sync"
	"testing"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestConvertFunctionToToolSetsDestructiveAnnotation(t *testing.T) {
	manager := &PulsarFunctionManager{
		logger:          logrus.New(),
		circuitBreakers: make(map[string]*CircuitBreaker),
		mutex:           sync.RWMutex{},
	}

	fnTool, err := manager.convertFunctionToTool(&utils.FunctionConfig{
		Tenant:    "public",
		Namespace: "default",
		Name:      "echo-fn",
		InputSpecs: map[string]utils.ConsumerConfig{
			"persistent://public/default/input": {},
		},
	})
	require.NoError(t, err)

	annotations := fnTool.Tool.Annotations
	require.NotEmpty(t, annotations.Title)
	require.NotNil(t, annotations.ReadOnlyHint)
	require.NotNil(t, annotations.DestructiveHint)
	require.False(t, *annotations.ReadOnlyHint)
	require.True(t, *annotations.DestructiveHint)
}
