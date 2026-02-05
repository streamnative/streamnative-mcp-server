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

package pftools

import (
	"errors"
	"net"
	"net/url"
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/rest"
)

var (
	// ErrFunctionNotFound indicates the function was not found.
	ErrFunctionNotFound = errors.New("function not found")
	// ErrNotOurMessage indicates a message that should be ignored.
	ErrNotOurMessage = errors.New("not our message")
	// ErrFunctionNoInputTopics indicates the function has no input topics.
	ErrFunctionNoInputTopics = errors.New("function has no input topics")
	// ErrSchemaConversionFailed indicates the schema conversion failed.
	ErrSchemaConversionFailed = errors.New("schema conversion failed")
)

// IsClusterUnhealthy checks if an error indicates cluster health issues
func IsClusterUnhealthy(err error) bool {
	if err == nil {
		return false
	}

	if restErr, ok := err.(rest.Error); ok {
		if restErr.Code == 503 && strings.Contains(strings.ToLower(restErr.Reason), "no healthy upstream") {
			return true
		}
	}

	return IsNetworkError(err)
}

// IsAuthError reports whether the error is an authorization error.
func IsAuthError(err error) bool {
	if err == nil {
		return false
	}
	if restErr, ok := err.(rest.Error); ok {
		return restErr.Code == 401 || restErr.Code == 403
	}
	return isAuthErrorText(err.Error())
}

// IsNotFoundError reports whether the error is a not found error.
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	if restErr, ok := err.(rest.Error); ok {
		return restErr.Code == 404
	}
	return isNotFoundText(err.Error())
}

// IsNetworkError reports whether the error indicates a network failure.
func IsNetworkError(err error) bool {
	if err == nil {
		return false
	}
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return true
	}

	errStr := strings.ToLower(err.Error())
	networkErrorPatterns := []string{
		"connection reset",
		"connection refused",
		"broken pipe",
		"tls handshake timeout",
		"i/o timeout",
		"context deadline exceeded",
		"timeout",
		"eof",
		"network is unreachable",
		"no route to host",
	}
	for _, pattern := range networkErrorPatterns {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}
	return false
}

func isAuthErrorText(text string) bool {
	lowered := strings.ToLower(text)
	authErrorPatterns := []string{
		"unauthorized",
		"forbidden",
		"token expired",
		"expired token",
		"invalid token",
		"401",
		"403",
	}
	for _, pattern := range authErrorPatterns {
		if strings.Contains(lowered, pattern) {
			return true
		}
	}
	return false
}

func isNotFoundText(text string) bool {
	lowered := strings.ToLower(text)
	if strings.Contains(lowered, "404") && strings.Contains(lowered, "not found") {
		return true
	}
	return false
}

// classifyConvertError reports whether a conversion failure is retryable.
func classifyConvertError(err error) failureCategory {
	if err == nil {
		return failureUnknown
	}
	if errors.Is(err, ErrFunctionNoInputTopics) || errors.Is(err, ErrSchemaConversionFailed) {
		return failurePermanent
	}
	if IsClusterUnhealthy(err) || IsAuthError(err) || IsNetworkError(err) {
		return failureRetryable
	}
	return failureUnknown
}
