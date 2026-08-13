# Copilot instructions

Use the repo's concrete workflow and architecture, not generic Go defaults. The codebase targets Go `1.25.12` from `go.mod`.

## Build, test, and lint

```bash
go mod verify
go mod download
go fmt ./...
go mod tidy
golangci-lint run --timeout=3m
go test -race ./...
go test -race ./pkg/mcp/builders/...
go test -v -run TestName ./pkg/...
make build
```

- CI expects `go fmt ./...` and `go mod tidy` to leave no diff behind.
- `make build` writes `bin/snmcp` from `cmd/streamnative-mcp-server/main.go`.
- If you add files or touch license headers, run `make license-check` and `make license-fix`.
- If you touch `charts/**`, `cmd/snmcp-e2e/**`, or `scripts/e2e-test.sh`, also run `./scripts/e2e-test.sh all`.
- If you change release or packaging flow, also run `go generate ./...` and `goreleaser release --snapshot --clean`.

## High-level architecture

- `cmd/streamnative-mcp-server/main.go` defines the `snmcp` Cobra CLI, adds shared config flags, and attaches the `stdio` and `sse` subcommands.
- `pkg/cmd/mcp/mcp.go` resolves config and selects one runtime mode:
  - `--key-file` enables StreamNative Cloud mode, defaults to `all` features, and ensures `streamnative-cloud` is present if features were set explicitly.
  - `--use-external-kafka` enables Kafka-only mode, rejects extra feature flags, and infers Kafka admin, client, and schema registry features.
  - `--use-external-pulsar` enables Pulsar-only mode, rejects extra feature flags, and infers `all-pulsar`.
- `pkg/cmd/mcp/server.go` is the server bootstrap and tool-registration hub. It creates backend sessions and wires every Kafka, Pulsar, and StreamNative tool family into one MCP server.
- `pkg/mcp/server.go` is the thin wrapper around `mark3labs/mcp-go`; it applies default resource capabilities, recovery, and logging.
- Request-scoped sessions are passed through `pkg/mcp/internal/context/ctx.go`. Tool handlers are expected to retrieve Kafka, Pulsar, and StreamNative Cloud sessions from there.
- `pkg/mcp/*_tools.go` files are mostly registration wrappers. The actual Kafka and Pulsar tool schemas and handler logic belong in `pkg/mcp/builders/{kafka,pulsar}`.
- StreamNative Cloud context, resource, and log tools still register from `pkg/mcp/streamnative_*_tools.go`; `pkg/mcp/builders/streamnative/` is not the primary implementation path today.
- `pkg/cmd/mcp/sse.go` owns SSE transport, auth middleware, and health endpoints at `<httpPath>/sse`, `<httpPath>/message`, `<httpPath>/healthz`, and `<httpPath>/readyz`.
- Multi-session Pulsar exists only in SSE plus external Pulsar mode. `pkg/mcp/session/` caches per-token sessions, and requests must send `Authorization: Bearer <token>`.
- `pkg/mcp/pftools/` plus `pkg/mcp/pulsar_functions_as_tools.go` dynamically expose deployed Pulsar Functions as MCP tools, so changes there alter the runtime tool surface.
- `docs/tools/` content is surfaced to MCP clients at runtime, so docs changes can be behavior-visible.

## Key conventions

- Preserve the current builder-and-wrapper split: add Kafka and Pulsar tool logic under `pkg/mcp/builders/{kafka,pulsar}` and keep `pkg/mcp/*_tools.go` focused on wiring.
- Reuse the typed helpers in `pkg/mcp/internal/context/ctx.go`; do not introduce parallel ad hoc context-key patterns.
- Keep feature gating aligned with `pkg/mcp/features.go`. External Kafka and external Pulsar modes infer their feature sets instead of accepting arbitrary extras.
- Respect `--read-only`: write-capable tools must not register or leak into read-only mode.
- Wrap errors with `%w`. In tool handlers, check for missing sessions and return `mcp.NewToolResultError(...)` instead of silently succeeding.
- Keep tests next to the code they validate. Builder packages already contain unit and parity tests; extend those patterns instead of inventing a separate layout.
- If you change public CLI flags, startup behavior, feature names, or connection examples, update `README.md`.
- If you change a runtime-visible tool, update the matching file in `docs/tools/`.
- If you change Helm values, update `charts/snmcp/README.md` and the chart templates together.
- If you change SSE auth, multi-session Pulsar behavior, or protocol/auth flows, review `pkg/cmd/mcp/sse.go`, `pkg/mcp/session/`, `cmd/snmcp-e2e/`, `scripts/e2e-test.sh`, `charts/snmcp/`, and the relevant `agents/` design docs together.
- Treat `sdk/sdk-apiserver/` and `sdk/sdk-kafkaconnect/` as generated local modules; prefer regeneration over hand-editing.
- Logging uses `logrus`, not `zap`.
- Most CLI flags also map to `SNMCP_*` environment variables, and the default config/state directory is `~/.snmcp` unless overridden.
