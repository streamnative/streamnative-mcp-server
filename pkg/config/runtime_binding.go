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
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// RuntimeBinding is a prepared generation. Once published its sessions are not
// reconfigured in place. Requests hold a lease until all use of these clients ends.
type RuntimeBinding struct {
	Instance string
	Cluster  string
	Pulsar   *pulsar.Session
	Kafka    *kafka.Session
}

// Close releases clients belonging to an unpublished or retired generation.
func (b *RuntimeBinding) Close() {
	if b == nil {
		return
	}
	if b.Pulsar != nil {
		b.Pulsar.ResetPulsarContext()
	}
	if b.Kafka != nil {
		b.Kafka.ResetKafkaContext()
	}
}

// LockRuntimeMutation serializes preparation, publication and reset operations.
func (s *Session) LockRuntimeMutation() func() {
	s.runtimeMutationMu.Lock()
	return s.runtimeMutationMu.Unlock
}

// AcquireRuntimeBinding pins a generation for an entire MCP request. Callers
// must use the returned snapshot rather than reacquire the read lock.
func (s *Session) AcquireRuntimeBinding() (*RuntimeBinding, func()) {
	s.runtimeBindingMu.RLock()
	return s.runtimeBinding, s.runtimeBindingMu.RUnlock
}

// CurrentRuntimeBinding is for setup/inspection outside a request. Concurrent
// request handlers must use AcquireRuntimeBinding instead.
func (s *Session) CurrentRuntimeBinding() *RuntimeBinding {
	s.runtimeBindingMu.RLock()
	defer s.runtimeBindingMu.RUnlock()
	return s.runtimeBinding
}

// PublishRuntimeBinding waits for in-flight requests, then replaces the whole
// generation. No fallible client initialization is allowed in this step.
// The caller must hold LockRuntimeMutation and pass a new, non-nil binding.
func (s *Session) PublishRuntimeBinding(next *RuntimeBinding) {
	s.runtimeBindingMu.Lock()
	old := s.runtimeBinding
	s.runtimeBinding = next
	s.SetPulsarClusterContext(next.Instance, next.Cluster)
	s.runtimeBindingMu.Unlock()
	old.Close()
}
