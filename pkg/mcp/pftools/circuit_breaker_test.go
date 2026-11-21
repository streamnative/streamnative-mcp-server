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
	"sync/atomic"
	"testing"
	"time"
)

const testResetTimeout = 30 * time.Millisecond

// openBreaker returns a breaker that is open and whose reset timeout has passed.
func openBreaker(t *testing.T) *CircuitBreaker {
	t.Helper()
	cb := NewCircuitBreaker(1, testResetTimeout)
	cb.RecordFailure()
	if cb.GetState() != StateOpen {
		t.Fatalf("expected open after reaching the failure threshold, got %s", cb.GetStateString())
	}
	if cb.AllowRequest() {
		t.Fatal("expected requests to be rejected while open")
	}
	time.Sleep(2 * testResetTimeout)
	return cb
}

func TestCircuitBreakerHalfOpenAllowsSingleProbe(t *testing.T) {
	cb := openBreaker(t)

	var allowed int32
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if cb.AllowRequest() {
				atomic.AddInt32(&allowed, 1)
			}
		}()
	}
	wg.Wait()

	if got := atomic.LoadInt32(&allowed); got != 1 {
		t.Fatalf("expected exactly one probe request to be allowed, got %d", got)
	}
	if cb.GetState() != StateHalfOpen {
		t.Fatalf("expected half-open while the probe is in flight, got %s", cb.GetStateString())
	}
}

func TestCircuitBreakerProbeSuccessCloses(t *testing.T) {
	cb := openBreaker(t)

	if !cb.AllowRequest() {
		t.Fatal("expected the probe request to be allowed")
	}
	cb.RecordSuccess()

	if cb.GetState() != StateClosed {
		t.Fatalf("expected closed after a successful probe, got %s", cb.GetStateString())
	}
	if !cb.AllowRequest() || !cb.AllowRequest() {
		t.Fatal("expected all requests to be allowed once closed")
	}
}

func TestCircuitBreakerProbeFailureReopens(t *testing.T) {
	cb := openBreaker(t)

	if !cb.AllowRequest() {
		t.Fatal("expected the probe request to be allowed")
	}
	cb.RecordFailure()

	if cb.GetState() != StateOpen {
		t.Fatalf("expected open after a failed probe, got %s", cb.GetStateString())
	}
	if cb.AllowRequest() {
		t.Fatal("expected requests to be rejected after a failed probe")
	}
}

func TestCircuitBreakerSuccessIsNotOverwrittenByLateTransition(t *testing.T) {
	// Run it a few times: the old implementation flipped the state to half-open
	// from a separate goroutine, which could land after RecordSuccess.
	for i := 0; i < 20; i++ {
		cb := openBreaker(t)
		if !cb.AllowRequest() {
			t.Fatal("expected the probe request to be allowed")
		}
		cb.RecordSuccess()
		time.Sleep(5 * time.Millisecond)

		if cb.GetState() != StateClosed {
			t.Fatalf("run %d: expected closed after a successful probe, got %s", i, cb.GetStateString())
		}
	}
}

func TestCircuitBreakerAllowsNewProbeIfFirstNeverReports(t *testing.T) {
	cb := openBreaker(t)

	if !cb.AllowRequest() {
		t.Fatal("expected the first probe to be allowed")
	}
	if cb.AllowRequest() {
		t.Fatal("expected a second request to be rejected while the probe is in flight")
	}

	// The probe never calls RecordSuccess/RecordFailure (for example it panicked).
	time.Sleep(2 * testResetTimeout)

	if !cb.AllowRequest() {
		t.Fatal("expected a new probe to be allowed after the reset timeout")
	}
}
