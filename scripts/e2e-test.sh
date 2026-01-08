#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
E2E_DIR="${ROOT_DIR}/charts/snmcp/e2e"

PULSAR_CONTAINER="${PULSAR_CONTAINER:-pulsar-standalone}"
PULSAR_IMAGE="${PULSAR_IMAGE:-snstage/pulsar-all:4.1.0.10}"
KIND_NETWORK="${KIND_NETWORK:-kind}"
PULSAR_WEB_PORT="${PULSAR_WEB_PORT:-8080}"
PULSAR_BROKER_PORT="${PULSAR_BROKER_PORT:-6650}"
PULSAR_STARTUP_TIMEOUT="${PULSAR_STARTUP_TIMEOUT:-180}"
PULSAR_STARTUP_INTERVAL="${PULSAR_STARTUP_INTERVAL:-3}"

SNMCP_RELEASE="${SNMCP_RELEASE:-snmcp}"
SNMCP_NAMESPACE="${SNMCP_NAMESPACE:-default}"
SNMCP_CHART_DIR="${SNMCP_CHART_DIR:-${ROOT_DIR}/charts/snmcp}"
SNMCP_FEATURES="${SNMCP_FEATURES:-pulsar-admin,pulsar-client}"
SNMCP_IMAGE_REPO="${SNMCP_IMAGE_REPO:-}"
SNMCP_IMAGE_TAG="${SNMCP_IMAGE_TAG:-}"
SNMCP_WAIT_TIMEOUT="${SNMCP_WAIT_TIMEOUT:-180s}"

TOKEN_ENV_FILE="${TOKEN_ENV_FILE:-${E2E_DIR}/test-tokens.env}"
TOKEN_SECRET_FILE="${TOKEN_SECRET_FILE:-${E2E_DIR}/test-secret.key}"

log() {
  echo "[e2e] $*"
}

die() {
  echo "[e2e] $*" >&2
  exit 1
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "missing command: $1"
}

load_tokens() {
  [[ -f "$TOKEN_ENV_FILE" ]] || die "missing token env file: $TOKEN_ENV_FILE"
  set -a
  # shellcheck disable=SC1090
  source "$TOKEN_ENV_FILE"
  set +a
  [[ -n "${ADMIN_TOKEN:-}" ]] || die "ADMIN_TOKEN not set in $TOKEN_ENV_FILE"
}

ensure_kind_network() {
  docker network inspect "$KIND_NETWORK" >/dev/null 2>&1 || die "missing kind network: $KIND_NETWORK"
}

pulsar_ip() {
  docker inspect -f "{{.NetworkSettings.Networks.${KIND_NETWORK}.IPAddress}}" "$PULSAR_CONTAINER"
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
    bash -lc 'set -- $(hostname -i); export PULSAR_PREFIX_advertisedAddress=$1; exec bin/pulsar standalone' \
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
    --set "server.features={${SNMCP_FEATURES}}"
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
  deploy-mcp     Deploy snmcp Helm chart and wait for readiness
  cleanup        Remove Pulsar container and uninstall snmcp release
  all            Run setup-pulsar then deploy-mcp
USAGE
}

main() {
  local cmd="${1:-}"
  case "$cmd" in
    setup-pulsar)
      setup_pulsar
      ;;
    deploy-mcp)
      deploy_mcp
      ;;
    cleanup)
      cleanup
      ;;
    all)
      setup_pulsar
      deploy_mcp
      ;;
    *)
      usage
      exit 1
      ;;
  esac
}

main "$@"
