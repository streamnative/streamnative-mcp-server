# AGENTS.md
Guide for autonomous coding agents (e.g., OpenAI Codex) working in streamnative-mcp-server

---

## 1 | Project snapshot
* What it is - A fast, developer-friendly Model Context Protocol (MCP) server that lets LLM agents talk to Apache Kafka, Apache Pulsar, and StreamNative Cloud through a single interface.
* Key outputs - Go binary `snmcp`; Docker images `streamnative/mcp-server` (preferred) and `streamnative/snmcp` (legacy); Helm chart `charts/snmcp`.
* Transports - stdio and SSE (HTTP).
* Primary language / tooling - Go 1.24.4, make, golangci-lint, goreleaser, license-eye.

---

## 2 | Repo map (orient yourself quickly)

| Path | What lives here | Touch with care? |
|------|-----------------|------------------|
| `cmd/streamnative-mcp-server/` | `main.go` entry point for the CLI/server | yes |
| `pkg/cmd/mcp/` | stdio/SSE transport, auth middleware, health endpoints | yes |
| `pkg/mcp/` | MCP server, sessions, feature flags, tool registration | yes |
| `pkg/mcp/builders/` | Tool builder framework and concrete tool builders | yes |
| `pkg/mcp/session/` | Multi-session Pulsar cache and cleanup | yes |
| `docs/tools/` | One Markdown file per MCP tool surfaced to LLMs at runtime | yes |
| `charts/snmcp/` | Helm chart for Kubernetes deployment | yes |
| `sdk/` | Generated Go SDKs for external APIs | can be re-generated |
| `.github/workflows/` | CI (lint, tests, release) | only if changing CI |
| `Makefile` | Local build helpers | safe |

---

## 3 | Required dev workflow

> Agents MUST follow every step below before committing code or opening a PR.

1. Install deps
   ```bash
   brew install golangci-lint            # or use the official install script
   go install github.com/apache/skywalking-eyes/cmd/license-eye@latest
   ```

2. Build and unit tests
   ```bash
   make build            # invokes `go build` with version metadata
   go test ./...
   ```

3. Static analysis and formatting
   Run `golangci-lint run` and ensure zero issues. Use `gofmt`/`go fmt` and keep import grouping consistent (use goimports if configured).

4. License headers
   ```bash
   make license-fix
   make license-check
   ```

5. Generate artifacts (only if you edited generated SDKs or OpenAPI specs)
   - Follow the README/HOW_TO_BUILD files under `sdk/`.
   - Run `go mod tidy` after regeneration.

6. Commit and conventional message
   Use Conventional Commits (`feat:`, `fix:`, `docs:`). Keep title <= 72 chars and add a body explaining why.

7. Run the full release checks locally (optional but recommended)
   ```bash
   goreleaser release --snapshot --clean
   ```

---

## 4 | How to run the server locally

```bash
# StreamNative Cloud (stdio)
bin/snmcp stdio --organization $ORG --key-file my_sa.json

# External Kafka
bin/snmcp stdio --use-external-kafka --kafka-bootstrap-servers localhost:9092

# External Pulsar
bin/snmcp stdio --use-external-pulsar --pulsar-web-service-url http://localhost:8080
```

```bash
# SSE server
bin/snmcp sse --http-addr :9090 --http-path /mcp --organization $ORG --key-file my_sa.json
```

SSE notes:
- Health endpoints: `/mcp/healthz` and `/mcp/readyz`.
- Multi-session Pulsar (SSE + external Pulsar only): `--multi-session-pulsar`, `--session-cache-size`, `--session-ttl-minutes`.

---

## 5 | Architecture overview

Request flow:
```
Client Request -> MCP Server (pkg/mcp/server.go)
                 -> Transport (pkg/cmd/mcp)
                 -> Tool Handler (builders)
                 -> Context functions (pkg/mcp/ctx.go)
                 -> Service client (Kafka/Pulsar/SNCloud)
```

Core components:
1. Server and sessions - `pkg/mcp/server.go` manages Kafka/Pulsar/SNCloud sessions and MCP server wiring.
2. Tool builders - `pkg/mcp/builders/` implements the builder pattern and feature validation.
3. Tool registration - `pkg/mcp/*_tools.go` adds tools based on feature flags.
4. Functions-as-tools (PFTools) - dynamic Pulsar Functions mapping with a circuit breaker.
5. Session management - `pkg/mcp/session/` implements LRU + TTL cleanup for multi-session Pulsar mode.

Key patterns:
- Builder + registry for tool assembly.
- Context injection for sessions.
- Feature flags defined in `pkg/mcp/features.go`.
- Multi-session mode disables the global Pulsar session and requires per-request tokens.

---

## 6 | Coding conventions

* Package layout - one feature per package; avoid cyclic imports.
* Error handling - wrap with `fmt.Errorf("context: %w", err)`; export sentinel errors when needed.
* Logging - use `logrus` with structured fields.
* Tests - use Go `testing` plus `testify/require`. When adding a tool, add:
  * happy-path invocation
  * typical error path
  * integration stub (may be skipped with `-short`)

---

## 7 | Common tasks for agents

| Task | Checklist |
|------|-----------|
| Add a new MCP tool | 1) create package in `pkg/mcp/builders/...` 2) update `docs/tools/<tool>.md` 3) add to feature flag map 4) tests + go vet |
| Bug fix | reproduce with unit test first -> fix -> ensure lint/test pass |
| Docs | update README and the per-tool doc; regenerate any table of contents |
| Release prep | bump version tag, update changelog, run `goreleaser release --snapshot --clean` |

---

## 8 | Programmatic checks the agent MUST run

```bash
go vet ./...
golangci-lint run
go test ./...
make build
make license-check
```

Codex MUST not finish a task until all checks pass locally. If CI fails, iterate until green.

---

## 9 | Pull-request etiquette

* Open PR against `main`, no new branches needed.
* Include a Summary, Testing plan, and Docs updated checklist.
* Mention which `--features` flags were affected, so reviewers know what to smoke-test.
* If the change affects the public CLI, update `README.md` usage examples.

---

## 10 | AGENTS.md spec compliance

This file follows the AGENTS.md spec described in the Codex system message (scope, precedence, required programmatic checks, etc.).

* Its scope is the entire repository.
* Deeper-level AGENTS.md (not currently present) may override these instructions for their subtree.
* Direct human instructions in a prompt override anything here.

See `CLAUDE.md` for deeper architecture and workflow details.
