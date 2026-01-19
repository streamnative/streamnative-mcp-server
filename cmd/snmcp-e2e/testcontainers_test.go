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

package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoadTestcontainersConfigDefaults(t *testing.T) {
	cfg := loadTestcontainersConfig()
	require.False(t, cfg.Enabled)
	require.Equal(t, defaultPulsarImage, cfg.PulsarImage)
	require.Equal(t, defaultKafkaImage, cfg.KafkaImage)
	require.Equal(t, defaultStartupTimeout, cfg.StartupTimeout)
}

func TestLoadTestcontainersConfigOverrides(t *testing.T) {
	t.Setenv("E2E_USE_TESTCONTAINERS", "true")
	t.Setenv("E2E_PULSAR_IMAGE", "apachepulsar/pulsar:3.2.2")
	t.Setenv("E2E_KAFKA_IMAGE", "confluentinc/confluent-local:7.6.0")
	t.Setenv("E2E_TESTCONTAINERS_TIMEOUT", "2m")

	cfg := loadTestcontainersConfig()
	require.True(t, cfg.Enabled)
	require.Equal(t, "apachepulsar/pulsar:3.2.2", cfg.PulsarImage)
	require.Equal(t, "confluentinc/confluent-local:7.6.0", cfg.KafkaImage)
	require.Equal(t, 2*time.Minute, cfg.StartupTimeout)
}

func TestLoadTestcontainersConfigInvalidTimeout(t *testing.T) {
	t.Setenv("E2E_TESTCONTAINERS_TIMEOUT", "not-a-duration")
	cfg := loadTestcontainersConfig()
	require.Equal(t, defaultStartupTimeout, cfg.StartupTimeout)
}

func TestStartTestcontainersDisabled(t *testing.T) {
	cfg := testcontainersConfig{Enabled: false}
	env, err := startTestcontainers(context.Background(), cfg)
	require.ErrorIs(t, err, errTestcontainersDisabled)
	require.Nil(t, env)
}

func TestTestcontainersEnvTerminateNil(t *testing.T) {
	var env *testcontainersEnv
	require.NoError(t, env.Terminate(context.Background()))
}

func TestValidateTestcontainersConfig(t *testing.T) {
	cfg := testcontainersConfig{
		Enabled:        true,
		PulsarImage:    "apachepulsar/pulsar:latest",
		KafkaImage:     "confluentinc/confluent-local:7.5.0",
		StartupTimeout: time.Minute,
	}
	require.NoError(t, validateTestcontainersConfig(cfg))

	cfg.PulsarImage = ""
	require.Error(t, validateTestcontainersConfig(cfg))

	cfg.PulsarImage = "apachepulsar/pulsar:latest"
	cfg.KafkaImage = ""
	require.Error(t, validateTestcontainersConfig(cfg))

	cfg.KafkaImage = "confluentinc/confluent-local:7.5.0"
	cfg.StartupTimeout = 0
	require.Error(t, validateTestcontainersConfig(cfg))
}
