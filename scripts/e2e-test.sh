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

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
E2E_DIR="${ROOT_DIR}/charts/snmcp/e2e"

PULSAR_CONTAINER="${PULSAR_CONTAINER:-pulsar-standalone}"
PULSAR_IMAGE="${PULSAR_IMAGE:-apachepulsar/pulsar-all:4.1.0}"
KIND_NETWORK="${KIND_NETWORK:-kind}"
KIND_CLUSTER_NAME="${KIND_CLUSTER_NAME:-kind}"
PULSAR_WEB_PORT="${PULSAR_WEB_PORT:-8080}"
PULSAR_BROKER_PORT="${PULSAR_BROKER_PORT:-6650}"
PULSAR_STARTUP_TIMEOUT="${PULSAR_STARTUP_TIMEOUT:-180}"
PULSAR_STARTUP_INTERVAL="${PULSAR_STARTUP_INTERVAL:-3}"

SNMCP_RELEASE="${SNMCP_RELEASE:-snmcp}"
SNMCP_NAMESPACE="${SNMCP_NAMESPACE:-default}"
SNMCP_CHART_DIR="${SNMCP_CHART_DIR:-${ROOT_DIR}/charts/snmcp}"
SNMCP_FEATURES=""
SNMCP_IMAGE_REPO="${SNMCP_IMAGE_REPO:-}"
SNMCP_IMAGE_TAG="${SNMCP_IMAGE_TAG:-}"
SNMCP_WAIT_TIMEOUT="${SNMCP_WAIT_TIMEOUT:-180s}"
SNMCP_HTTP_PATH="${SNMCP_HTTP_PATH:-/mcp}"
SNMCP_SERVICE_PORT="${SNMCP_SERVICE_PORT:-9090}"
SNMCP_LOCAL_PORT="${SNMCP_LOCAL_PORT:-19090}"
SNMCP_PORT_FORWARD_TIMEOUT="${SNMCP_PORT_FORWARD_TIMEOUT:-60}"
SNMCP_E2E_BIN="${SNMCP_E2E_BIN:-${ROOT_DIR}/bin/snmcp-e2e}"
SNMCP_PORT_FORWARD_PID=""

TOKEN_ENV_FILE="${TOKEN_ENV_FILE:-${E2E_DIR}/test-tokens.env}"
TOKEN_SECRET_FILE="${TOKEN_SECRET_FILE:-${E2E_DIR}/test-secret.key}"

log() {
  echo "[e2e] $*"
}

die() {
  echo "[e2e] $*" >&2
  exit 1
}

collect_logs() {
  log "collecting debug logs"
  if command -v kubectl >/dev/null 2>&1; then
    log "snmcp logs (last 200 lines)"
    kubectl logs "deployment/${SNMCP_RELEASE}" \
      --namespace "$SNMCP_NAMESPACE" \
      --tail=200 \
      || true
  fi

  if command -v docker >/dev/null 2>&1; then
    if docker ps -a --format '{{.Names}}' | grep -qx "$PULSAR_CONTAINER"; then
      log "pulsar container logs (last 200 lines)"
      docker logs --tail 200 "$PULSAR_CONTAINER" || true
    fi
  fi
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "missing command: $1"
}

load_tokens() {
  if [[ -f "$TOKEN_ENV_FILE" ]]; then
    set -a
    # shellcheck disable=SC1090
    source "$TOKEN_ENV_FILE"
    set +a
  fi

  if [[ -z "${ADMIN_TOKEN:-}" || -z "${TEST_USER_TOKEN:-}" ]]; then
    require_cmd python3
    [[ -f "$TOKEN_SECRET_FILE" ]] || die "missing token secret file: $TOKEN_SECRET_FILE"
    generate_test_tokens "$TOKEN_SECRET_FILE"
  fi

  [[ -n "${ADMIN_TOKEN:-}" ]] || die "ADMIN_TOKEN is not set and could not be generated"
  [[ -n "${TEST_USER_TOKEN:-}" ]] || die "TEST_USER_TOKEN is not set and could not be generated"
}

generate_test_tokens() {
  local secret_file="$1"
  local generated
  if ! generated="$(python3 - "$secret_file" <<'PY'
import base64
import hashlib
import hmac
import json
import sys

secret_path = sys.argv[1]
with open(secret_path, "rb") as secret_file:
    secret = secret_file.read()


def b64url(value):
    return base64.urlsafe_b64encode(value).rstrip(b"=").decode("ascii")


def token(subject):
    header = {"alg": "HS256", "typ": "JWT"}
    payload = {"exp": 4102444800, "iat": 1700000000, "sub": subject}
    signing_input = ".".join(
        [
            b64url(json.dumps(header, separators=(",", ":")).encode("utf-8")),
            b64url(json.dumps(payload, separators=(",", ":")).encode("utf-8")),
        ]
    ).encode("ascii")
    signature = hmac.new(secret, signing_input, hashlib.sha256).digest()
    return signing_input.decode("ascii") + "." + b64url(signature)

print("ADMIN_TOKEN=" + token("admin"))
print("TEST_USER_TOKEN=" + token("test-user"))
PY
  )"; then
    die "failed to generate e2e JWT tokens"
  fi
  ADMIN_TOKEN="$(printf '%s\n' "$generated" | awk -F= '$1 == "ADMIN_TOKEN" {print substr($0, index($0, "=") + 1)}')"
  TEST_USER_TOKEN="$(printf '%s\n' "$generated" | awk -F= '$1 == "TEST_USER_TOKEN" {print substr($0, index($0, "=") + 1)}')"
  export ADMIN_TOKEN TEST_USER_TOKEN
}

ensure_kind_network() {
  docker network inspect "$KIND_NETWORK" >/dev/null 2>&1 || die "missing kind network: $KIND_NETWORK"
}

pulsar_ip() {
  docker inspect -f "{{.NetworkSettings.Networks.${KIND_NETWORK}.IPAddress}}" "$PULSAR_CONTAINER"
}

wait_for_http() {
  local url="$1"
  local timeout="$2"
  local deadline=$((SECONDS + timeout))
  while ((SECONDS < deadline)); do
    if curl -fsS "$url" >/dev/null; then
      return 0
    fi
    sleep 2
  done
  return 1
}

setup_pulsar() {
  require_cmd docker
  require_cmd curl
  ensure_kind_network
  load_tokens
  [[ -f "$TOKEN_SECRET_FILE" ]] || die "missing secret key file: $TOKEN_SECRET_FILE"

  if docker ps -a --format '{{.Names}}' | grep -qx "$PULSAR_CONTAINER"; then
    log "removing existing container: $PULSAR_CONTAINER"
    docker rm -f "$PULSAR_CONTAINER" >/dev/null
  fi

  log "starting pulsar container: $PULSAR_CONTAINER"
  docker run -d \
    --name "$PULSAR_CONTAINER" \
    --network "$KIND_NETWORK" \
    -e PULSAR_PREFIX_authenticationEnabled=true \
    -e PULSAR_PREFIX_authenticationProviders=org.apache.pulsar.broker.authentication.AuthenticationProviderToken \
    -e PULSAR_PREFIX_authorizationEnabled=true \
    -e PULSAR_PREFIX_superUserRoles=admin \
    -e PULSAR_PREFIX_tokenSecretKey=file:///pulsarctl/test/auth/token/secret.key \
    -e PULSAR_PREFIX_brokerClientAuthenticationPlugin=org.apache.pulsar.client.impl.auth.AuthenticationToken \
    -e PULSAR_PREFIX_brokerClientAuthenticationParameters="token:${ADMIN_TOKEN}" \
    -v "$TOKEN_SECRET_FILE:/pulsarctl/test/auth/token/secret.key:ro" \
    "$PULSAR_IMAGE" \
    bash -lc 'set -- $(hostname -i); export PULSAR_PREFIX_advertisedAddress=$1; export JAVA_HOME=/opt/jvm; export PATH="$JAVA_HOME/bin:$PATH"; bin/apply-config-from-env.py /pulsar/conf/standalone.conf; exec bin/pulsar standalone' \
    >/dev/null

  log "waiting for pulsar to be ready"
  local deadline=$((SECONDS + PULSAR_STARTUP_TIMEOUT))
  while ((SECONDS < deadline)); do
    local ip
    if ip="$(pulsar_ip 2>/dev/null)" && [[ -n "$ip" ]]; then
      if curl -fsS "http://${ip}:${PULSAR_WEB_PORT}/admin/v2/clusters" \
        -H "Authorization: Bearer ${ADMIN_TOKEN}" >/dev/null; then
        log "pulsar ready at ${ip}"
        return 0
      fi
    fi
    sleep "$PULSAR_STARTUP_INTERVAL"
  done

  die "pulsar not ready after ${PULSAR_STARTUP_TIMEOUT}s"
}

build_image() {
  require_cmd docker
  require_cmd kind

  SNMCP_IMAGE_REPO="${SNMCP_IMAGE_REPO:-snmcp-e2e}"
  SNMCP_IMAGE_TAG="${SNMCP_IMAGE_TAG:-local}"

  local image_ref="${SNMCP_IMAGE_REPO}:${SNMCP_IMAGE_TAG}"
  log "building image ${image_ref}"
  docker build -t "$image_ref" -f "${ROOT_DIR}/Dockerfile" "$ROOT_DIR" >/dev/null

  log "loading image into kind"
  kind load docker-image "$image_ref" --name "$KIND_CLUSTER_NAME" >/dev/null
}

deploy_mcp() {
  require_cmd helm
  require_cmd kubectl
  require_cmd docker
  ensure_kind_network

  local ip
  ip="$(pulsar_ip)"
  [[ -n "$ip" ]] || die "failed to resolve pulsar container IP"

  local helm_args=(
    --namespace "$SNMCP_NAMESPACE"
    --create-namespace
    --set "pulsar.webServiceURL=http://${ip}:${PULSAR_WEB_PORT}"
    --set "pulsar.serviceURL=pulsar://${ip}:${PULSAR_BROKER_PORT}"
    --wait
    --timeout "$SNMCP_WAIT_TIMEOUT"
  )

  if [[ -n "$SNMCP_IMAGE_REPO" ]]; then
    helm_args+=(--set "image.repository=${SNMCP_IMAGE_REPO}")
  fi
  if [[ -n "$SNMCP_IMAGE_TAG" ]]; then
    helm_args+=(--set "image.tag=${SNMCP_IMAGE_TAG}")
  fi

  log "deploying snmcp with helm"
  helm upgrade --install "$SNMCP_RELEASE" "$SNMCP_CHART_DIR" "${helm_args[@]}" >/dev/null

  log "waiting for snmcp deployment rollout"
  kubectl rollout status "deployment/${SNMCP_RELEASE}" \
    --namespace "$SNMCP_NAMESPACE" \
    --timeout "$SNMCP_WAIT_TIMEOUT" \
    >/dev/null

  log "snmcp deployed and ready"
}

run_tests() {
  require_cmd kubectl
  require_cmd go
  require_cmd curl
  load_tokens

  log "building snmcp-e2e binary"
  go build -o "$SNMCP_E2E_BIN" "${ROOT_DIR}/cmd/snmcp-e2e" >/dev/null

  log "starting port-forward for snmcp service"
  kubectl port-forward "svc/${SNMCP_RELEASE}" "${SNMCP_LOCAL_PORT}:${SNMCP_SERVICE_PORT}" \
    --namespace "$SNMCP_NAMESPACE" >/dev/null 2>&1 &
  SNMCP_PORT_FORWARD_PID=$!

  trap 'if [[ -n "${SNMCP_PORT_FORWARD_PID:-}" ]]; then kill "$SNMCP_PORT_FORWARD_PID" >/dev/null 2>&1 || true; fi' RETURN

  local health_url="http://127.0.0.1:${SNMCP_LOCAL_PORT}${SNMCP_HTTP_PATH}/healthz"
  if ! wait_for_http "$health_url" "$SNMCP_PORT_FORWARD_TIMEOUT"; then
    die "port-forward did not become ready within ${SNMCP_PORT_FORWARD_TIMEOUT}s"
  fi

  local http_base="http://127.0.0.1:${SNMCP_LOCAL_PORT}${SNMCP_HTTP_PATH}"
  log "running snmcp-e2e against ${http_base}"
  if ! E2E_HTTP_BASE="$http_base" "$SNMCP_E2E_BIN"; then
    collect_logs
    return 1
  fi
}

cleanup() {
  require_cmd docker

  if command -v helm >/dev/null 2>&1; then
    helm uninstall "$SNMCP_RELEASE" --namespace "$SNMCP_NAMESPACE" >/dev/null 2>&1 || true
  fi

  if docker ps -a --format '{{.Names}}' | grep -qx "$PULSAR_CONTAINER"; then
    docker rm -f "$PULSAR_CONTAINER" >/dev/null
  fi
}

usage() {
  cat <<USAGE
Usage: $0 <command>

Commands:
  setup-pulsar   Start Pulsar standalone with JWT auth on the kind network
  build-image    Build snmcp image and load into kind
  deploy-mcp     Deploy snmcp Helm chart and wait for readiness
  run-tests      Port-forward snmcp service and run E2E client
  cleanup        Remove Pulsar container and uninstall snmcp release
  all            Run setup-pulsar, build-image, deploy-mcp, and run-tests
USAGE
}

main() {
  local cmd="${1:-}"
  case "$cmd" in
    setup-pulsar)
      setup_pulsar
      ;;
    build-image)
      build_image
      ;;
    deploy-mcp)
      deploy_mcp
      ;;
    run-tests)
      run_tests
      ;;
    cleanup)
      cleanup
      ;;
    all)
      setup_pulsar
      build_image
      deploy_mcp
      run_tests
      ;;
    *)
      usage
      exit 1
      ;;
  esac
}

main "$@"
