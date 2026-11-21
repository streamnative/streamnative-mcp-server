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

// Package pftools provides Pulsar Functions tooling for MCP.
package pftools

import (
	"fmt"
	"time"
)

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(threshold int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		failureCount:     0,
		failureThreshold: threshold,
		resetTimeout:     resetTimeout,
		state:            StateClosed,
	}
}

// RecordSuccess records a successful operation
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	// Reset the failure count
	cb.failureCount = 0

	// If we're half-open, close the circuit
	if cb.state == StateHalfOpen {
		cb.state = StateClosed
	}
}

// RecordFailure records a failed operation
func (cb *CircuitBreaker) RecordFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	// Record the failure
	cb.lastFailure = time.Now()

	// If we're already open, do nothing
	if cb.state == StateOpen {
		return
	}

	// If we're half-open, open the circuit immediately
	if cb.state == StateHalfOpen {
		cb.state = StateOpen
		return
	}

	// Increment the failure count
	cb.failureCount++

	// If we've exceeded the threshold, open the circuit
	if cb.failureCount >= cb.failureThreshold {
		cb.state = StateOpen
	}
}

// AllowRequest determines if a request should be allowed
func (cb *CircuitBreaker) AllowRequest() bool {
	// This can change state, so it needs the write lock.
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	switch cb.state {
	case StateClosed:
		// Always allow if closed
		return true
	case StateOpen:
		// If open, check if timeout has expired
		if time.Since(cb.lastFailure) > cb.resetTimeout {
			// Move to half-open and let this request through as the probe
			cb.state = StateHalfOpen
			cb.probeStartedAt = time.Now()
			return true
		}
		// Still open and timeout hasn't expired
		return false
	case StateHalfOpen:
		// A probe is already in flight, so hold everything else back until it
		// reports success or failure. If it never reports, allow a new probe
		// after another reset timeout so the breaker can't get stuck.
		if time.Since(cb.probeStartedAt) > cb.resetTimeout {
			cb.probeStartedAt = time.Now()
			return true
		}
		return false
	default:
		return false
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// ForceOpen forces the circuit breaker to open
func (cb *CircuitBreaker) ForceOpen() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	cb.state = StateOpen
	cb.lastFailure = time.Now()
}

// ForceClose forces the circuit breaker to close
func (cb *CircuitBreaker) ForceClose() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	cb.state = StateClosed
	cb.failureCount = 0
}

// Reset resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	cb.state = StateClosed
	cb.failureCount = 0
}

// GetStateString returns a string representation of the circuit breaker state
func (cb *CircuitBreaker) GetStateString() string {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	switch cb.state {
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF-OPEN"
	case StateClosed:
		return "CLOSED"
	default:
		return "UNKNOWN"
	}
}

// String returns a string representation of the circuit breaker
func (cb *CircuitBreaker) String() string {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	return fmt.Sprintf("CircuitBreaker{state=%s, failures=%d/%d, lastFailure=%v}",
		cb.GetStateString(), cb.failureCount, cb.failureThreshold, cb.lastFailure)
}
