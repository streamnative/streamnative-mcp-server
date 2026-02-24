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
	"errors"
	"testing"
	"time"

	pulsarsdk "github.com/apache/pulsar-client-go/pulsar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscribeConsumerWithRetryConfig_SuccessAfterRetryableErrors(t *testing.T) {
	attempts := 0
	subscribe := func(pulsarsdk.ConsumerOptions) (pulsarsdk.Consumer, error) {
		attempts++
		if attempts < 3 {
			return nil, errors.New("ServiceNotReady: Please redo the lookup")
		}
		return nil, nil
	}

	consumer, usedAttempts, err := subscribeConsumerWithRetryConfig(
		context.Background(),
		subscribe,
		pulsarsdk.ConsumerOptions{},
		3,
		func(int) time.Duration { return 0 },
	)
	require.NoError(t, err)
	assert.Nil(t, consumer)
	assert.Equal(t, 3, attempts)
	assert.Equal(t, 3, usedAttempts)
}

func TestSubscribeConsumerWithRetryConfig_NoRetryForNonRetryableError(t *testing.T) {
	attempts := 0
	subscribe := func(pulsarsdk.ConsumerOptions) (pulsarsdk.Consumer, error) {
		attempts++
		return nil, errors.New("AuthenticationError: invalid token")
	}

	consumer, usedAttempts, err := subscribeConsumerWithRetryConfig(
		context.Background(),
		subscribe,
		pulsarsdk.ConsumerOptions{},
		3,
		func(int) time.Duration { return 0 },
	)
	require.Error(t, err)
	assert.Nil(t, consumer)
	assert.Equal(t, 1, attempts)
	assert.Equal(t, 1, usedAttempts)
	assert.Contains(t, err.Error(), "after 1 attempt(s)")
}

func TestSubscribeConsumerWithRetryConfig_InterruptedDuringBackoff(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	attempts := 0
	subscribe := func(pulsarsdk.ConsumerOptions) (pulsarsdk.Consumer, error) {
		attempts++
		if attempts == 1 {
			cancel()
		}
		return nil, errors.New("not served by this instance")
	}

	consumer, usedAttempts, err := subscribeConsumerWithRetryConfig(
		ctx,
		subscribe,
		pulsarsdk.ConsumerOptions{},
		3,
		func(int) time.Duration { return time.Second },
	)
	require.Error(t, err)
	assert.Nil(t, consumer)
	assert.Equal(t, 1, attempts)
	assert.Equal(t, 1, usedAttempts)
	assert.Contains(t, err.Error(), "interrupted")
}

func TestIsLookupRetryableSubscribeError(t *testing.T) {
	assert.True(t, isLookupRetryableSubscribeError(errors.New("ServiceNotReady: anything")))
	assert.True(t, isLookupRetryableSubscribeError(errors.New("Please redo the lookup now")))
	assert.True(t, isLookupRetryableSubscribeError(errors.New("topic not served by this instance")))
	assert.False(t, isLookupRetryableSubscribeError(errors.New("AuthenticationError")))
	assert.False(t, isLookupRetryableSubscribeError(nil))
}

func TestConsumeSubscribeBackoff(t *testing.T) {
	assert.Equal(t, consumeSubscribeBackoffFirst, consumeSubscribeBackoff(1))
	assert.Equal(t, consumeSubscribeBackoffStep2, consumeSubscribeBackoff(2))
	assert.Equal(t, consumeSubscribeBackoffThird, consumeSubscribeBackoff(3))
	assert.Equal(t, consumeSubscribeBackoffThird, consumeSubscribeBackoff(100))
}
