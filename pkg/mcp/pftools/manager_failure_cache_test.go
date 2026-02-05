package pftools

import (
	"testing"
	"time"
)

func TestShouldSkipFailure(t *testing.T) {
	now := time.Date(2026, 2, 5, 12, 0, 0, 0, time.UTC)
	pollInterval := 30 * time.Second

	t.Run("permanent skips", func(t *testing.T) {
		state := &functionFailureState{category: failurePermanent}
		if !shouldSkipFailure(state, pollInterval, now) {
			t.Fatalf("expected permanent failure to skip")
		}
	})

	t.Run("retryable respects interval", func(t *testing.T) {
		state := &functionFailureState{
			category:      failureRetryable,
			lastAttemptAt: now.Add(-10 * time.Second),
		}
		if !shouldSkipFailure(state, pollInterval, now) {
			t.Fatalf("expected retryable failure within interval to skip")
		}
	})

	t.Run("retryable after interval", func(t *testing.T) {
		state := &functionFailureState{
			category:      failureRetryable,
			lastAttemptAt: now.Add(-40 * time.Second),
		}
		if shouldSkipFailure(state, pollInterval, now) {
			t.Fatalf("expected retryable failure after interval to retry")
		}
	})

	t.Run("unknown skips", func(t *testing.T) {
		state := &functionFailureState{category: failureUnknown}
		if !shouldSkipFailure(state, pollInterval, now) {
			t.Fatalf("expected unknown failure to skip")
		}
	})
}

func TestShouldLogFailure(t *testing.T) {
	t.Run("first failure logs", func(t *testing.T) {
		if !shouldLogFailure(nil, "hash", failureRetryable, "boom") {
			t.Fatalf("expected first failure to log")
		}
	})

	t.Run("same state does not log", func(t *testing.T) {
		prev := &functionFailureState{
			configHash: "hash",
			category:   failureRetryable,
			lastError:  "boom",
		}
		if shouldLogFailure(prev, "hash", failureRetryable, "boom") {
			t.Fatalf("expected identical failure to not log")
		}
	})

	t.Run("category change logs", func(t *testing.T) {
		prev := &functionFailureState{
			configHash: "hash",
			category:   failureRetryable,
			lastError:  "boom",
		}
		if !shouldLogFailure(prev, "hash", failurePermanent, "boom") {
			t.Fatalf("expected category change to log")
		}
	})

	t.Run("error change logs", func(t *testing.T) {
		prev := &functionFailureState{
			configHash: "hash",
			category:   failureRetryable,
			lastError:  "boom",
		}
		if !shouldLogFailure(prev, "hash", failureRetryable, "different") {
			t.Fatalf("expected error change to log")
		}
	})

	t.Run("config change logs", func(t *testing.T) {
		prev := &functionFailureState{
			configHash: "hash",
			category:   failureRetryable,
			lastError:  "boom",
		}
		if !shouldLogFailure(prev, "newhash", failureRetryable, "boom") {
			t.Fatalf("expected config change to log")
		}
	})
}
