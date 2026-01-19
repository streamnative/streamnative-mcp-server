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
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	pulsarBrokerPort      = "6650/tcp"
	pulsarWebServicePort  = "8080/tcp"
	defaultPulsarImage    = "apachepulsar/pulsar:latest"
	defaultKafkaImage     = "confluentinc/confluent-local:7.5.0"
	defaultStartupTimeout = 5 * time.Minute
)

var errTestcontainersDisabled = errors.New("testcontainers disabled")

type testcontainersConfig struct {
	Enabled        bool
	PulsarImage    string
	KafkaImage     string
	StartupTimeout time.Duration
}

type testcontainersEnv struct {
	Pulsar              testcontainers.Container
	Kafka               *kafka.KafkaContainer
	PulsarWebServiceURL string
	PulsarBrokerURL     string
	KafkaBrokers        string
}

func loadTestcontainersConfig() testcontainersConfig {
	return testcontainersConfig{
		Enabled:        getenvBool("E2E_USE_TESTCONTAINERS", false),
		PulsarImage:    getenv("E2E_PULSAR_IMAGE", defaultPulsarImage),
		KafkaImage:     getenv("E2E_KAFKA_IMAGE", defaultKafkaImage),
		StartupTimeout: getenvDuration("E2E_TESTCONTAINERS_TIMEOUT", defaultStartupTimeout),
	}
}

func startTestcontainers(ctx context.Context, cfg testcontainersConfig) (*testcontainersEnv, error) {
	if !cfg.Enabled {
		return nil, errTestcontainersDisabled
	}
	if err := validateTestcontainersConfig(cfg); err != nil {
		return nil, err
	}

	startupCtx, cancel := context.WithTimeout(ctx, cfg.StartupTimeout)
	defer cancel()

	pulsarContainer, err := startPulsarContainer(startupCtx, cfg)
	if err != nil {
		return nil, err
	}

	pulsarWebURL, pulsarBrokerURL, err := resolvePulsarEndpoints(startupCtx, pulsarContainer)
	if err != nil {
		_ = pulsarContainer.Terminate(startupCtx)
		return nil, err
	}

	kafkaContainer, kafkaBrokers, err := startKafkaContainer(startupCtx, cfg)
	if err != nil {
		_ = pulsarContainer.Terminate(startupCtx)
		return nil, err
	}

	env := &testcontainersEnv{
		Pulsar:              pulsarContainer,
		Kafka:               kafkaContainer,
		PulsarWebServiceURL: pulsarWebURL,
		PulsarBrokerURL:     pulsarBrokerURL,
		KafkaBrokers:        kafkaBrokers,
	}
	return env, nil
}

func (env *testcontainersEnv) Terminate(ctx context.Context) error {
	if env == nil {
		return nil
	}

	var err error
	if env.Kafka != nil {
		err = errors.Join(err, env.Kafka.Terminate(ctx))
	}
	if env.Pulsar != nil {
		err = errors.Join(err, env.Pulsar.Terminate(ctx))
	}
	return err
}

func validateTestcontainersConfig(cfg testcontainersConfig) error {
	if strings.TrimSpace(cfg.PulsarImage) == "" {
		return errors.New("pulsar image is required")
	}
	if strings.TrimSpace(cfg.KafkaImage) == "" {
		return errors.New("kafka image is required")
	}
	if cfg.StartupTimeout <= 0 {
		return errors.New("startup timeout must be positive")
	}
	return nil
}

func startPulsarContainer(ctx context.Context, cfg testcontainersConfig) (testcontainers.Container, error) {
	request := testcontainers.ContainerRequest{
		Image:        cfg.PulsarImage,
		ExposedPorts: []string{pulsarBrokerPort, pulsarWebServicePort},
		Cmd:          []string{"standalone"},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(pulsarBrokerPort),
			wait.ForListeningPort(pulsarWebServicePort),
		).WithDeadline(cfg.StartupTimeout),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("start pulsar container: %w", err)
	}
	return container, nil
}

func resolvePulsarEndpoints(ctx context.Context, container testcontainers.Container) (string, string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return "", "", fmt.Errorf("pulsar host: %w", err)
	}
	webPort, err := container.MappedPort(ctx, pulsarWebServicePort)
	if err != nil {
		return "", "", fmt.Errorf("pulsar web service port: %w", err)
	}
	brokerPort, err := container.MappedPort(ctx, pulsarBrokerPort)
	if err != nil {
		return "", "", fmt.Errorf("pulsar broker port: %w", err)
	}

	webURL := fmt.Sprintf("http://%s:%s", host, webPort.Port())
	brokerURL := fmt.Sprintf("pulsar://%s:%s", host, brokerPort.Port())
	return webURL, brokerURL, nil
}

func startKafkaContainer(ctx context.Context, cfg testcontainersConfig) (*kafka.KafkaContainer, string, error) {
	container, err := kafka.Run(ctx, cfg.KafkaImage)
	if err != nil {
		return nil, "", fmt.Errorf("start kafka container: %w", err)
	}

	brokers, err := container.Brokers(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, "", fmt.Errorf("kafka brokers: %w", err)
	}
	return container, strings.Join(brokers, ","), nil
}

func getenvDuration(key string, fallback time.Duration) time.Duration {
	raw := os.Getenv(key)
	if strings.TrimSpace(raw) == "" {
		return fallback
	}
	value, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}
	return value
}
