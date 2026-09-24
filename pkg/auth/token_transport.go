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

package auth

import (
	"net/http"

	"golang.org/x/oauth2"
)

// TokenTransport supplies current credentials for each request without retrying
// application operations. An empty Username selects Bearer authentication.
type TokenTransport struct {
	Source   oauth2.TokenSource
	Username string
	Base     http.RoundTripper
}

// RoundTrip clones the request so authentication cannot mutate a caller's headers.
func (t *TokenTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	token, err := t.Source.Token()
	if err != nil {
		return nil, err
	}
	cloned := request.Clone(request.Context())
	if t.Username == "" {
		cloned.Header.Set("Authorization", "Bearer "+token.AccessToken)
	} else {
		cloned.SetBasicAuth(t.Username, token.AccessToken)
	}
	return t.Transport().RoundTrip(cloned)
}

// Transport returns the underlying TLS-aware HTTP transport.
func (t *TokenTransport) Transport() http.RoundTripper {
	if t.Base == nil {
		return http.DefaultTransport
	}
	return t.Base
}

// WithTransport implements the Pulsar admin authentication provider contract.
func (t *TokenTransport) WithTransport(base http.RoundTripper) { t.Base = base }
