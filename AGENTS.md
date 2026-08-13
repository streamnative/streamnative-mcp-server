# AGENTS.md
Guide for autonomous coding agents working in `streamnative-mcp-server`.

---

## 1 | Project Snapshot

- **What it is**: A Go-based Model Context Protocol (MCP) server that exposes StreamNative Cloud, Apache Kafka, and Apache Pulsar operations through one MCP surface.
- **Primary binary**: `bin/snmcp`, built from `cmd/streamnative-mcp-server/main.go`.
- **Server modes**: `stdio` and `sse`.
- **Current stack**: Go `1.25.12`, `cobra`, `mark3labs/mcp-go`, Docker, Helm, GoReleaser.
- **Published container names**: `streamnative/mcp-server` and `streamnative/snmcp`.

---

## 2 | Repository Map

| Path | Purpose | Notes |
|------|---------|-------|
| `cmd/streamnative-mcp-server/` | Root CLI entrypoint | Owns the `snmcp` command and top-level Cobra wiring |
| `cmd/snmcp-e2e/` | Go E2E client | Used by chart and SSE authentication end-to-end tests |
| `pkg/cmd/mcp/` | Transport/bootstrap layer | Completes config, creates sessions, starts `stdio` or `sse` servers |
| `pkg/mcp/` | Core MCP server logic | Feature flags, prompts, context tools, wrapper registration, runtime instructions |
| `pkg/mcp/internal/context/` | Session/context plumbing | Reuse these typed context helpers instead of introducing parallel context-key patterns |
| `pkg/mcp/builders/kafka/` | Kafka tool builders | Preferred home for Kafka MCP tool schemas and handlers |
| `pkg/mcp/builders/pulsar/` | Pulsar tool builders | Preferred home for Pulsar MCP tool schemas and handlers |
| `pkg/mcp/session/` | Multi-session Pulsar support | SSE bearer-token session cache and eviction logic |
| `pkg/config/` | CLI/config resolution | StreamNative Cloud, external Kafka, external Pulsar options |
| `pkg/kafka/`, `pkg/pulsar/`, `pkg/auth/`, `pkg/schema/` | Service/session/auth/schema helpers | Touch carefully; these support the tool layer |
| `docs/tools/` | Tool-facing Markdown docs | These docs are surfaced to MCP clients at runtime |
| `agents/` | Design notes and MCP proposals | Review when changing protocol, auth, compliance, or multi-session behavior |
| `charts/snmcp/` | Helm chart | Deploys SSE mode with multi-session Pulsar support |
| `charts/snmcp/e2e/` | Chart test fixtures | Test tokens, local values, and auth fixtures for E2E |
| `scripts/e2e-test.sh` | Local chart E2E harness | Uses Docker, Kind, Helm, kubectl, and `cmd/snmcp-e2e` |
| `sdk/sdk-apiserver/` | Local generated SDK module | Referenced via `replace` in root `go.mod` |
| `sdk/sdk-kafkaconnect/` | Local generated SDK module | OpenAPI-generated code; avoid casual hand-edits |
| `.github/workflows/` | CI pipelines | Unit tests, lint, GoReleaser, chart E2E |

---

## 3 | Current Architecture Notes

- `cmd/streamnative-mcp-server/main.go` creates the root `snmcp` command, adds shared config flags, and attaches `stdio` and `sse` subcommands.
- `pkg/cmd/mcp/mcp.go` selects one runtime mode:
  - `--key-file` => StreamNative Cloud mode. This defaults to `all` features, or forces `streamnative-cloud` into the selected feature set.
  - `--use-external-kafka` => Kafka-only mode. Extra feature flags are rejected and Kafka admin/client features are inferred automatically.
  - `--use-external-pulsar` => Pulsar-only mode. Extra feature flags are rejected and `all-pulsar` is inferred automatically.
- Most CLI flags are also exposed through `SNMCP_*` environment variables. The default local state/config directory is `~/.snmcp` unless `--config-dir` or `SNMCP_CONFIG_DIR` overrides it.
- `pkg/mcp/server.go` wraps the underlying `mark3labs/mcp-go` server and applies the default MCP capabilities, recovery, and logging options.
- `pkg/cmd/mcp/server.go` is the current tool registration hub. It wires sessions, builds the MCP server, and registers each tool family.
- Request-scoped Kafka, Pulsar, and StreamNative Cloud sessions flow through `pkg/mcp/internal/context/ctx.go`. Reuse these helpers instead of creating new ad hoc context keys.
- The `pkg/mcp/*_tools.go` files are mostly registration wrappers. The actual tool schemas and handlers live in `pkg/mcp/builders/...`.
- New Kafka and Pulsar MCP tool work should follow the builder pattern under `pkg/mcp/builders/{kafka,pulsar}`. Do not create a new `pkg/tools/...` tree.
- StreamNative Cloud resource and log tools are still registered from `pkg/mcp/streamnative_*_tools.go`; `pkg/mcp/builders/streamnative/` exists but is not the active implementation path today.
- `pkg/cmd/mcp/sse.go` exposes the SSE endpoint, message endpoint, and unauthenticated health probes:
  - `<httpPath>/sse`
  - `<httpPath>/message`
  - `<httpPath>/healthz`
  - `<httpPath>/readyz`
- Multi-session Pulsar mode only exists for `sse` with external Pulsar. It uses bearer tokens and the cache in `pkg/mcp/session/`.
- `pkg/mcp/pftools/` and `pkg/mcp/pulsar_functions_as_tools.go` dynamically expose deployed Pulsar Functions as MCP tools. Changes here affect runtime-generated tool surfaces and should be reviewed carefully.
- Functions-as-tools runtime behavior is also shaped by `FUNCTIONS_AS_TOOLS_*` environment variables. Treat changes there as runtime-surface changes, not just internal refactors.
- Logging currently uses `logrus`, not `zap`.

---

## 4 | Required Local Workflow

Agents should match the current repository workflow, not a generic Go template.

1. **Dependencies**
   ```bash
   go mod verify
   go mod download
   ```

2. **Formatting and module hygiene**
   ```bash
   go fmt ./...
   go mod tidy
   ```
   Both commands are enforced by CI and must leave no diff behind.

3. **Lint**
   ```bash
   golangci-lint run --timeout=3m
   ```
   CI currently installs `golangci-lint` `v2.7.2`.

4. **Tests**
   ```bash
   go test -race ./...
   ```

5. **Build**
   ```bash
   make build
   ```
   This builds `bin/snmcp` with version, commit, and build-date ldflags. It is a host build only; cross-platform packaging is handled by GoReleaser.

6. **License headers**
   ```bash
   make license-check
   make license-fix
   ```
   Use these when you add files or touch headers. The repo uses `license-eye`, not `go-licenser`.

7. **Chart / E2E changes**
   If you touch `charts/**`, `cmd/snmcp-e2e/**`, or `scripts/e2e-test.sh`, also run:
   ```bash
   ./scripts/e2e-test.sh all
   ```
   This requires Docker, Kind, Helm, and kubectl.

8. **Release-related changes**
   For changes that affect packaging, Docker images, or release metadata, also run:
   ```bash
   goreleaser release --snapshot --clean
   ```
   GoReleaser also runs `go mod tidy` and `go generate ./...` in `before.hooks`, so run `go generate ./...` locally when you touch generated inputs or release plumbing.

---

## 5 | Local Runbook

```bash
# Build
make build

# StreamNative Cloud over stdio
bin/snmcp stdio --organization "$ORG" --key-file /path/to/key.json

# External Kafka over stdio
bin/snmcp stdio --use-external-kafka --kafka-bootstrap-servers localhost:9092

# External Pulsar over stdio
bin/snmcp stdio --use-external-pulsar --pulsar-web-service-url http://localhost:8080

# SSE mode
bin/snmcp sse --http-addr :9090 --http-path /mcp --organization "$ORG" --key-file /path/to/key.json

# SSE multi-session Pulsar mode
bin/snmcp sse \
  --http-addr :9090 \
  --http-path /mcp \
  --use-external-pulsar \
  --pulsar-web-service-url http://localhost:8080 \
  --multi-session-pulsar
```

For chart-based local development:

```bash
helm install snmcp ./charts/snmcp \
  --set pulsar.webServiceURL=http://pulsar.example.com:8080
```

---

## 6 | Coding Conventions For This Repo

- Prefer small, reviewable changes that preserve the current builder-and-wrapper architecture.
- Wrap errors with context using `%w`.
- Keep feature gating aligned with `pkg/mcp/features.go`.
- Respect `--read-only` behavior. Write-capable tools must not leak into read-only mode.
- Put tests next to the code they validate. Builder packages already contain unit and parity tests; extend those patterns instead of inventing a separate layout.
- If you change public CLI flags, startup behavior, feature names, or connection examples, update `README.md`.
- If you change a runtime-visible tool, update the matching file in `docs/tools/`.
- If you change Helm values, update `charts/snmcp/README.md` and template assumptions together.
- If you change MCP protocol handling, authentication flows, or multi-session behavior, update the relevant design notes under `agents/` when the current design docs would become stale.
- Generated SDK directories are local modules. Prefer regenerating from source specs/config rather than manually editing generated files unless you are making a deliberate patch.
- `.golangci.yml` enables `gofmt` and `goimports`, but CI explicitly checks `go fmt ./...`, `go mod tidy`, and `golangci-lint run`. Keep your local workflow aligned with those concrete commands.

---

## 7 | Common Tasks

| Task | Expected steps |
|------|----------------|
| **Add or extend an MCP tool** | Implement or update the builder in `pkg/mcp/builders/{kafka,pulsar}`; wire or update the wrapper in `pkg/mcp/*_tools.go` if needed; update `pkg/mcp/features.go` when feature gates change; add or update tests; update `docs/tools/<tool>.md` |
| **Change CLI or configuration flow** | Update `pkg/config/`, `pkg/cmd/mcp/`, and any affected README examples; verify mode-selection rules still make sense |
| **Change SSE or auth/session behavior** | Review `pkg/cmd/mcp/sse.go`, `pkg/mcp/session/`, `pkg/mcp/internal/context/`, `cmd/snmcp-e2e/`, `scripts/e2e-test.sh`, `charts/snmcp/templates/`, and the relevant design docs in `agents/` together |
| **Change Helm chart behavior** | Update `values.yaml`, templates, and `charts/snmcp/README.md` together; rerun the E2E harness |
| **Change generated SDKs** | Update the relevant source spec/config, regenerate, keep the root `replace` directives valid, and avoid partial manual edits |
| **Change protocol/auth design docs** | Keep the matching documents in `agents/` aligned with the implemented behavior; do not leave architecture or compliance notes describing superseded flows |
| **Release prep** | Check `.goreleaser.yaml`, `Dockerfile.goreleaser`, Docker image naming, Homebrew handoff, and run `go generate ./...` plus `goreleaser release --snapshot --clean` |

---

## 8 | Programmatic Checks To Run Before Finishing

```bash
go mod verify
go mod download
go fmt ./...
go mod tidy
golangci-lint run --timeout=3m
go test -race ./...
make build
```

When the change touches the chart or E2E harness:

```bash
./scripts/e2e-test.sh all
```

Do not mark work complete while these checks are still failing.

---

## 9 | Commit And PR Expectations

- Use Conventional Commits such as `feat:`, `fix:`, `docs:`, or `refactor:`.
- Open PRs against `main` unless a release branch is explicitly requested.
- Include a short summary, a concrete testing section, and a docs-updated note.
- Mention any affected feature flags, CLI flags, SSE paths, or Helm values so reviewers know what to smoke-test.

---

## 10 | Scope

- This file applies to the entire repository.
- If a deeper `AGENTS.md` is added later, it overrides this file for its subtree.
- Direct user instructions still take precedence.
