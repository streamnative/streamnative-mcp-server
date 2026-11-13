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

package mcp

import "testing"

// Admin URLs must point at the cluster endpoint, never at the StreamNative proxy,
// because the proxy rejects the service account tokens that the MCP server uses.
func TestClusterEndpointURLs(t *testing.T) {
	const dns = "pc-1234.gcp-shared-usc1.g.sn2.dev"
	if got, want := getWebServiceURL(dns), "https://pc-1234.gcp-shared-usc1.g.sn2.dev"; got != want {
		t.Errorf("getWebServiceURL() = %q, want %q", got, want)
	}
	if got, want := getKafkaConnectURL(dns), "https://pc-1234.gcp-shared-usc1.g.sn2.dev/admin/kafkaconnect/"; got != want {
		t.Errorf("getKafkaConnectURL() = %q, want %q", got, want)
	}
}
