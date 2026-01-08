# StreamNative MCP Server Helm Chart

This Helm chart deploys the StreamNative MCP Server on Kubernetes with **Multi-Session Pulsar** mode.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- External Pulsar cluster accessible from the Kubernetes cluster

## Installation

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
| `server.readOnly` | `false` | Enable read-only mode |
| `server.features` | `[]` | Features to enable (default: all-pulsar) |
| `server.httpAddr` | `:9090` | HTTP server address |
| `server.httpPath` | `/mcp` | SSE endpoint path |

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
helm install snmcp ./charts/snmcp \
  --set pulsar.webServiceURL=http://pulsar:8080
```

### With Ingress

```bash
helm install snmcp ./charts/snmcp \
  --set pulsar.webServiceURL=http://pulsar:8080 \
  --set ingress.enabled=true \
  --set ingress.hosts[0].host=mcp.example.com \
  --set ingress.hosts[0].paths[0].path=/ \
  --set ingress.hosts[0].paths[0].pathType=Prefix
```

### Read-Only Mode with Limited Features

```bash
helm install snmcp ./charts/snmcp \
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

helm install snmcp ./charts/snmcp \
  --set pulsar.webServiceURL=https://pulsar:8443 \
  --set pulsar.tls.enabled=true \
  --set pulsar.tls.secretName=pulsar-tls \
  --set pulsar.tls.trustCertsFilePath=/etc/snmcp/tls/ca.crt
```

## Authentication

This chart runs MCP Server in Multi-Session Pulsar mode. Each client request must include a valid Pulsar JWT token:

```bash
curl -H "Authorization: Bearer <PULSAR_TOKEN>" http://localhost:9090/mcp/sse
```

## License

Apache License 2.0
