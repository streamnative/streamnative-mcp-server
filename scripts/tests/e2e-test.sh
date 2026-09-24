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

# Native command-seam tests: no Docker daemon, Kind cluster, or Go build needed.
# Run with: bash scripts/tests/e2e-test.sh
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
HARNESS="${ROOT_DIR}/scripts/e2e-test.sh"
TEST_ROOT="$(mktemp -d "${TMPDIR:-/tmp}/snmcp-harness-tests.XXXXXX")"
STUB_BIN="${TEST_ROOT}/bin"
CASE_DIR=""
CASE_STATUS=0
PASSED=0
FAILED=0
mkdir -p "$STUB_BIN"

cleanup() {
  local pid_file pid
  for pid_file in "$TEST_ROOT"/*/port-forward.pid; do
    [[ -f "$pid_file" ]] || continue
    read -r pid < "$pid_file"
    kill "$pid" 2>/dev/null || true
  done
  rm -rf "$TEST_ROOT"
}
trap cleanup EXIT

cat > "${STUB_BIN}/stub" <<'STUB'
#!/usr/bin/env bash
set -euo pipefail
name="${0##*/}"
# Assemble first, then append in one write so concurrent curl/port-forward calls
# cannot interleave their argv records.
record="$name"
for argument in "$@"; do record+="$(printf '\t%s' "$argument")"; done
printf '%s\n' "$record" >> "${CASE_DIR}/commands"
case "$name" in
  docker)
    case "${1:-}" in
      inspect) echo 172.18.0.2 ;;
      ps) echo pulsar-standalone ;;
      logs) echo stub-pulsar-diagnostic ;;
    esac
    ;;
  kubectl)
    while [[ "${1:-}" == --* ]]; do
      case "$1" in
        --*=*) shift ;;
        *) shift 2 ;;
      esac
    done
    case "${1:-}" in
      port-forward)
        echo "$$" > "${CASE_DIR}/port-forward.pid"
        exec sleep 60
        ;;
      logs) echo stub-snmcp-diagnostic ;;
      get|describe) echo stub-kubernetes-diagnostic ;;
    esac
    ;;
  go)
    [[ "${1:-}" == build && "${2:-}" == -o ]] || exit 91
    cp "${STUB_BIN}/client" "$3"
    ;;
  snmcp-e2e|client)
    printf 'transport=%s\nread_only=%s\nbase=%s\n' \
      "${E2E_TRANSPORT-<unset>}" "${E2E_READ_ONLY-<unset>}" \
      "${E2E_HTTP_BASE-<unset>}" > "${CASE_DIR}/client-env"
    exit "${STUB_CLIENT_EXIT:-0}"
    ;;
  curl)
    # Synchronize with the background stub before the harness can terminate it.
    for ((attempt = 0; attempt < 100; attempt++)); do
      [[ -f "${CASE_DIR}/port-forward.pid" ]] && break
      sleep 0.01
    done
    exit "${STUB_CURL_EXIT:-0}"
    ;;
  helm|kind) ;;
  *) echo "unexpected external command: $name" >&2; exit 92 ;;
esac
STUB
chmod +x "${STUB_BIN}/stub"
for name in docker kubectl helm kind go curl python3 client; do
  ln -s stub "${STUB_BIN}/${name}"
done

new_case() {
  CASE_DIR="${TEST_ROOT}/$1"
  mkdir -p "$CASE_DIR"
  : > "${CASE_DIR}/commands"
}

invoke() {
  local command="$1"
  shift
  CASE_STATUS=0
  # Clear caller configuration; an omitted setting must exercise the real default.
  env -i PATH="${STUB_BIN}:$PATH" HOME="$TEST_ROOT" \
    STUB_BIN="$STUB_BIN" CASE_DIR="$CASE_DIR" \
    ADMIN_TOKEN=stub-admin TEST_USER_TOKEN=stub-user \
    TOKEN_ENV_FILE="${CASE_DIR}/no-token-file" \
    SNMCP_E2E_BIN="${CASE_DIR}/snmcp-e2e" \
    SNMCP_LOCAL_PORT=19090 SNMCP_PORT_FORWARD_TIMEOUT=1 \
    E2E_ARTIFACTS_DIR="${CASE_DIR}/artifacts" \
    "$@" bash "$HARNESS" "$command" > "${CASE_DIR}/output" 2>&1 || CASE_STATUS=$?
}

has_line() {
  grep -Fxq -- "$2" "$1" || { echo "missing line: $2 ($1)"; return 1; }
}

has_arg() {
  local command="$1" argument="$2"
  awk -F '\t' -v command="$command" -v argument="$argument" '
    $1 == command { for (i = 2; i <= NF; i++) if ($i == argument) found = 1 }
    END { exit !found }
  ' "${CASE_DIR}/commands" || { echo "missing $command argument: $argument"; return 1; }
}

port_forward_stopped() {
  local pid attempt
  [[ -f "${CASE_DIR}/port-forward.pid" ]] || { echo 'port-forward never started'; return 1; }
  read -r pid < "${CASE_DIR}/port-forward.pid"
  for ((attempt = 0; attempt < 40; attempt++)); do
    kill -0 "$pid" 2>/dev/null || return 0
    sleep 0.05
  done
  echo "port-forward survived harness exit: $pid"
  return 1
}

deploy_contract() {
  local transport="$1" read_only="$2"
  shift 2
  invoke deploy-mcp "$@" SNMCP_HTTP_PATH=/custom/mcp
  [[ "$CASE_STATUS" == 0 ]] || return 1
  has_arg helm "server.transport=${transport}" &&
    has_arg helm "server.readOnly=${read_only}" &&
    has_arg helm server.httpPath=/custom/mcp
}

run_contract() {
  local transport="$1" read_only="$2"
  shift 2
  invoke run-tests "$@" SNMCP_HTTP_PATH=/custom/mcp
  [[ "$CASE_STATUS" == 0 ]] || return 1
  has_line "${CASE_DIR}/client-env" "transport=${transport}" &&
    has_line "${CASE_DIR}/client-env" "read_only=${read_only}" &&
    has_line "${CASE_DIR}/client-env" base=http://127.0.0.1:19090/custom/mcp &&
    has_arg curl http://127.0.0.1:19090/custom/mcp/healthz &&
    port_forward_stopped
}

invalid_contract() {
  invoke "$1" "$2"
  [[ "$CASE_STATUS" != 0 ]] || { echo 'invalid configuration accepted'; return 1; }
  [[ ! -s "${CASE_DIR}/commands" ]] || { echo 'external command ran before validation'; return 1; }
}

bounded_wait_contract() {
  invoke run-tests
  [[ "$CASE_STATUS" == 0 ]] || return 1
  # Accept curl's long/short and --option=value spellings, require positive bounds.
  awk -F '\t' '
    $1 == "curl" {
      seen = 1; connect = 0; total = 0
      for (i = 2; i <= NF; i++) {
        if ($i == "--connect-timeout" && $(i+1) + 0 > 0) connect = 1
        if (($i == "--max-time" || $i == "-m") && $(i+1) + 0 > 0) total = 1
        if ($i ~ /^--connect-timeout=/) { split($i, a, "="); if (a[2] + 0 > 0) connect = 1 }
        if ($i ~ /^--max-time=/) { split($i, a, "="); if (a[2] + 0 > 0) total = 1 }
      }
      if (!connect || !total) bad = 1
    }
    END { exit (!seen || bad) }
  ' "${CASE_DIR}/commands" || { echo 'curl readiness calls need connect and total timeouts'; return 1; }
  port_forward_stopped
}

failure_cleanup_contract() {
  invoke run-tests "$1"
  [[ "$CASE_STATUS" != 0 ]] || { echo 'expected run-tests failure'; return 1; }
  port_forward_stopped
}

diagnostics_contract() {
  invoke "$1" STUB_CLIENT_EXIT=17
  if [[ "$1" == logs ]]; then
    [[ "$CASE_STATUS" == 0 ]] || return 1
  else
    [[ "$CASE_STATUS" != 0 ]] || return 1
    port_forward_stopped || return 1
  fi
  [[ -d "${CASE_DIR}/artifacts" ]] || { echo 'artifacts directory missing'; return 1; }
  grep -Rq stub-snmcp-diagnostic "${CASE_DIR}/artifacts" &&
    grep -Rq stub-pulsar-diagnostic "${CASE_DIR}/artifacts"
}

test_case() {
  local name="$1"
  shift
  new_case "$name"
  if "$@"; then
    echo "PASS $name"
    PASSED=$((PASSED + 1))
  else
    echo "FAIL $name"
    [[ ! -f "${CASE_DIR}/output" ]] || cat "${CASE_DIR}/output"
    FAILED=$((FAILED + 1))
  fi
}

test_case deploy-defaults deploy_contract sse false
test_case run-defaults run_contract sse false
for transport in sse http; do
  for read_only in false true; do
    test_case "deploy-${transport}-${read_only}" deploy_contract "$transport" "$read_only" \
      "SNMCP_TRANSPORT=$transport" "SNMCP_READ_ONLY=$read_only"
    test_case "run-${transport}-${read_only}" run_contract "$transport" "$read_only" \
      "SNMCP_TRANSPORT=$transport" "SNMCP_READ_ONLY=$read_only"
  done
done
for command in setup-pulsar build-image deploy-mcp run-tests cleanup all; do
  for setting in SNMCP_TRANSPORT=stdio SNMCP_TRANSPORT=HTTP SNMCP_READ_ONLY=1 SNMCP_READ_ONLY=TRUE; do
    test_case "invalid-${command}-${setting}" invalid_contract "$command" "$setting"
  done
done
test_case bounded-readiness bounded_wait_contract
test_case cleanup-client-failure failure_cleanup_contract STUB_CLIENT_EXIT=17
test_case cleanup-readiness-failure failure_cleanup_contract STUB_CURL_EXIT=7
test_case explicit-diagnostics diagnostics_contract logs
test_case failure-diagnostics diagnostics_contract run-tests
printf '\n%d passed, %d failed\n' "$PASSED" "$FAILED"
[[ "$FAILED" == 0 ]]
