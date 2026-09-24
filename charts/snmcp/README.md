# StreamNative MCP Server Helm Chart

This Helm chart deploys the StreamNative MCP Server on Kubernetes with **Multi-Session Pulsar** mode.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- External Pulsar cluster accessible from the Kubernetes cluster

## Installation

Install from the StreamNative Helm repository:

```bash
helm repo add streamnative https://charts.streamnative.io
helm repo update

helm install snmcp streamnative/snmcp \
  --set pulsar.webServiceURL=http://pulsar.example.com:8080
```

Optional: check available chart versions

```bash
helm search repo streamnative/snmcp
```

Install from local chart (for development):

```bash
helm install snmcp ./charts/snmcp \
  --set pulsar.webServiceURL=http://pulsar.example.com:8080
```

## Configuration

### Required Parameters

| Parameter | Description |
|-----------|-------------|
| `pulsar.webServiceURL` | Pulsar web service URL (required) |

### Server Configuration

| Parameter | Default | Description |
|-----------|---------|-------------|
| `server.transport` | `sse` | Transport: legacy `sse` or opt-in Streamable HTTP `http`; other values fail rendering |
| `server.readOnly` | `false` | Enable read-only mode |
| `server.features` | `[]` | Features to enable (default: all-pulsar) |
| `server.httpPath` | `/mcp` | HTTP MCP endpoint, or base path for SSE/message endpoints; also the health endpoint base path |
| `server.configDir` | `/var/lib/snmcp` | Config directory for snmcp state (must be writable) |

The container listens on port 9090. Use `service.port` or Ingress to expose a different port.
For `http`, the chart explicitly passes `--http-addr :9090` so Kubernetes Services and probes can reach it, overriding the CLI's loopback default. SSE arguments and bind behavior remain unchanged.

### Session Configuration

| Parameter | Default | Description |
|-----------|---------|-------------|
| `session.cacheSize` | `100` | Max cached sessions |
| `session.ttlMinutes` | `30` | Session TTL before eviction |

### Pulsar TLS Configuration

| Parameter | Default | Description |
|-----------|---------|-------------|
| `pulsar.tls.enabled` | `false` | Enable TLS for Pulsar |
| `pulsar.tls.secretName` | `""` | Secret containing TLS certs |
| `pulsar.tls.allowInsecureConnection` | `false` | Allow insecure TLS |
| `pulsar.tls.enableHostnameVerification` | `true` | Verify hostname |

### Kubernetes Resources

| Parameter | Default | Description |
|-----------|---------|-------------|
| `replicaCount` | `1` | Number of replicas |
| `service.type` | `ClusterIP` | Service type |
| `service.port` | `9090` | Service port |
| `ingress.enabled` | `false` | Enable Ingress |
| `serviceAccount.create` | `true` | Create ServiceAccount |

## Examples

### Basic Installation

```bash
helm install snmcp streamnative/snmcp \
  --set pulsar.webServiceURL=http://pulsar:8080
```

### Streamable HTTP (Opt-In)

Use an image that includes the `http` command; the chart's version and default image tag are unchanged.

```bash
helm install snmcp ./charts/snmcp \
  --set pulsar.webServiceURL=http://pulsar:8080 \
  --set server.transport=http \
  --set image.tag="<HTTP_CAPABLE_IMAGE_TAG>"
```

Configure a compatible Streamable HTTP MCP client with `http://localhost:9090/mcp` after port-forwarding and the `Authorization: Bearer <PULSAR_TOKEN>` header. The HTTP transport targets MCP `2026-07-28`. Use per-request protocol metadata and `server/discover`, not `initialize` or a plain GET `curl` request. The chart remains in Multi-Session Pulsar mode; changing transport does not enable Kafka or StreamNative Cloud mode.

Legacy SSE remains the default and uses `/mcp/sse` (with `/mcp/message` for messages). For either transport, health probes remain `/mcp/healthz` and `/mcp/readyz`. Setting `server.httpPath=/` uses `/` for HTTP MCP, `/sse` for legacy SSE, and `/healthz` and `/readyz` for probes.

### With Ingress

```bash
helm install snmcp streamnative/snmcp \
  --set pulsar.webServiceURL=http://pulsar:8080 \
  --set ingress.enabled=true \
  --set ingress.hosts[0].host=mcp.example.com \
  --set ingress.hosts[0].paths[0].path=/ \
  --set ingress.hosts[0].paths[0].pathType=Prefix
```

### Read-Only Mode with Limited Features

```bash
helm install snmcp streamnative/snmcp \
  --set pulsar.webServiceURL=http://pulsar:8080 \
  --set server.readOnly=true \
  --set server.features="{pulsar-admin,pulsar-client}"
```

### With Pulsar TLS

```bash
# First create a secret with TLS certs
kubectl create secret generic pulsar-tls \
  --from-file=ca.crt=./ca.crt \
  --from-file=tls.crt=./tls.crt \
  --from-file=tls.key=./tls.key

helm install snmcp streamnative/snmcp \
  --set pulsar.webServiceURL=https://pulsar:8443 \
  --set pulsar.tls.enabled=true \
  --set pulsar.tls.secretName=pulsar-tls \
  --set pulsar.tls.trustCertsFilePath=/etc/snmcp/tls/ca.crt
```

## Authentication

This chart runs MCP Server in Multi-Session Pulsar mode. Each client request must include a valid Pulsar JWT token:

For legacy SSE only:

```bash
curl -H "Authorization: Bearer <PULSAR_TOKEN>" http://localhost:9090/mcp/sse
```

For HTTP, set the same header in a Streamable HTTP MCP client. HTTP creates request-scoped clients and closes them after each response; session cache size/TTL values apply only to SSE. Client creation does not validate a credential; Pulsar authorizes backend operations, and discovery can succeed with an unauthorized token. This is backend token forwarding, **not an MCP OAuth resource server**; there is no MCP audience/issuer validation or OAuth discovery. Multi-session mode alone is not a public authentication boundary. Use an authenticating TLS gateway and network restrictions. HTTP does not inherit fixed tokens, auth plugins or client certificates into request-scoped clients.

The HTTP transport rejects requests containing an `Origin` header by default. Non-browser MCP clients without that header work; this chart does not configure browser-origin access. Do not strip browser Origin headers at a proxy to bypass this policy.

Before exposing either transport beyond a trusted network, configure a trusted TLS-terminating proxy/Ingress and restrict direct backend access with network controls (for example, Kubernetes NetworkPolicy). Forward the `Authorization` header to the backend and avoid logging tokens. The Ingress example above does not configure TLS or network restrictions by itself.

Health endpoints do not require authentication and can be used for liveness/readiness:

- `GET http://localhost:9090/mcp/healthz`
- `GET http://localhost:9090/mcp/readyz`

## Chart Validation

Run `bash charts/snmcp/tests/transport.sh` from the repository root with Helm 3.13+ (or Helm 4); the test uses client-only dry-run. This checks Helm lint/template for both transports, invalid transport rejection, and default/root/trailing-slash health paths without requiring a cluster.

### Pull request E2E checks

The `E2E Tests` workflow runs on every pull request (including SDK, Go backend,
and transport changes), pushes to `main`, and manual dispatch. Its isolated
Kind jobs build the checked-out revision into a local Docker image, load it
into Kind, and deploy this chart against a real JWT-enabled Pulsar standalone:

| Check | Transport | Profile | Endpoint |
| --- | --- | --- | --- |
| `E2E (legacy-sse)` | Legacy SSE | Existing integration suite | `/mcp/sse` |
| `E2E (modern-http)` | MCP 2026-07-28 HTTP | Protocol and backend authorization | `/mcp` |
| `E2E (modern-http-read-only)` | MCP 2026-07-28 HTTP | Read-only tool filtering and reads | `/custom/mcp` |

These are ordinary `pull_request` runs on hosted runners with read-only
repository permissions, no production secrets, and JWTs generated from the
committed **test-only** signing key. Do not reuse the test key in deployments.
Fork pull requests may need a maintainer to approve their first workflow run.
Repository administrators can select the three check names above as required
branch-protection checks after the first run; workflow code does not configure
branch protection itself. There are no path filters that leave a required check
pending on unrelated changes.

Each job uploads `e2e-<profile>-<attempt>` diagnostics (run output, pod/events,
current/previous MCP logs, Pulsar logs and port-forward output) before cleanup,
with seven-day retention. A failed assertion exits nonzero; jobs do not hide
test failures with `continue-on-error`. This validates the documented backend
credential-forwarding profile, not MCP OAuth compliance.

To reproduce on a Linux host with Docker, Go, Helm, kubectl and Kind:

```bash
kind create cluster --name kind
# Choose one lane; each all run recreates the test Pulsar container.
SNMCP_TRANSPORT=http ./scripts/e2e-test.sh all
SNMCP_TRANSPORT=http SNMCP_READ_ONLY=true SNMCP_HTTP_PATH=/custom/mcp \
  ./scripts/e2e-test.sh all
./scripts/e2e-test.sh logs
./scripts/e2e-test.sh cleanup
```

Omit `SNMCP_TRANSPORT` to retain the legacy SSE suite. `SNMCP_READ_ONLY` defaults
to `false`. `E2E_ARTIFACTS_DIR` defaults to `tmp/e2e`. The harness connects to
the Pulsar container IP on the Kind bridge, so Docker Desktop hosts may need
additional routing; the authoritative full E2E gate is the Linux PR workflow.
Rendering/stub checks run without Docker via `bash scripts/tests/e2e-test.sh`
and `bash charts/snmcp/tests/transport.sh`.

## License

Apache License 2.0
