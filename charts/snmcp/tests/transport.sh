#!/usr/bin/env bash
# Copyright 2026 StreamNative
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Rendering-only checks; no Kubernetes cluster or running backend required.
set -euo pipefail
chart=$(cd "$(dirname "$0")/.." && pwd)
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
common=(--set pulsar.webServiceURL=http://pulsar:8080)

for transport in default sse http; do
  options=("${common[@]}")
  if [[ $transport != default ]]; then
    options+=(--set "server.transport=$transport")
  fi
  helm lint "$chart" "${options[@]}"
  for path in /mcp /mcp/ /; do
    helm template test "$chart" "${options[@]}" \
      --set "server.httpPath=$path" > "$tmp/rendered"
    grep -Fq -- "path: ${path%/}/healthz" "$tmp/rendered"
    grep -Fq -- "path: ${path%/}/readyz" "$tmp/rendered"
    grep -Fq -- '- --multi-session-pulsar' "$tmp/rendered"
    if [[ $transport == http ]]; then
      grep -Fxq '            - http' "$tmp/rendered"
      grep -Fq -- '- --http-addr' "$tmp/rendered"
      grep -Fq -- '- ":9090"' "$tmp/rendered"
      if grep -Eq -- '--session-(cache-size|ttl-minutes)' "$tmp/rendered"; then
        echo 'HTTP must not pass legacy session cache flags' >&2
        exit 1
      fi
    else
      grep -Fxq '            - sse' "$tmp/rendered"
      if grep -Fq -- '--http-addr' "$tmp/rendered"; then
        echo 'SSE must retain its existing bind defaults' >&2
        exit 1
      fi
    fi
    # Client-only dry-run includes NOTES, unlike helm template.
    helm install test "$chart" --dry-run=client "${options[@]}" \
      --set "server.httpPath=$path" > "$tmp/install"
    sed -n '/^NOTES:/,$p' "$tmp/install" > "$tmp/notes"
    endpoint=${path%/}
    if [[ $transport == http ]]; then
      endpoint=${endpoint:-/}
      grep -Fq 'Configure a Streamable HTTP MCP client' "$tmp/notes"
      if grep -Fq 'curl -H' "$tmp/notes"; then
        echo 'HTTP NOTES must not recommend a GET curl request' >&2
        exit 1
      fi
    else
      endpoint=$endpoint/sse
      grep -Fq 'curl -H' "$tmp/notes"
    fi
    grep -Fq "http://localhost:9090$endpoint\"" "$tmp/notes"
  done
done

# Old releases reused with --reuse-values may have no transport key/default.
cp -R "$chart" "$tmp/legacy-chart"
sed '/^  transport: /d' "$chart/values.yaml" > "$tmp/legacy-chart/values.yaml"
helm template legacy "$tmp/legacy-chart" "${common[@]}" > "$tmp/legacy"
grep -Fxq '            - sse' "$tmp/legacy"
helm install legacy "$tmp/legacy-chart" --dry-run=client "${common[@]}" > "$tmp/legacy-notes"
grep -Fq 'http://localhost:9090/mcp/sse' "$tmp/legacy-notes"

for transport in invalid HTTP ''; do
  for operation in lint template; do
    if helm "$operation" "$chart" "${common[@]}" \
      --set-string "server.transport=$transport" > "$tmp/invalid" 2>&1; then
      # Helm lint treats the template fail function as a diagnostic, not
      # necessarily an error. Actual rendering must always reject it.
      if [[ $operation == template ]]; then
        echo "$operation accepted invalid transport '$transport'" >&2
        exit 1
      fi
    fi
    grep -Fq 'server.transport must be one of: sse, http' "$tmp/invalid"
  done
done
echo 'Transport rendering checks passed'
