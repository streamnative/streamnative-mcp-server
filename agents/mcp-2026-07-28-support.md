# MCP 2026-07-28: research and staged support plan

Research date: **2026-09-07**. Repository baseline: `77134bc`. The research below
distinguishes upstream support from application conformance. The implementation
snapshot records the accompanying changes; it does not certify full support.

## Implementation snapshot

- Upgraded `go.mod`/`go.sum` from mcp-go v0.58.0 to v1.0.0.
- Expanded `pkg/mcp/server_test.go` to cover all four legacy versions, empty
  and unknown-version fallback, and the legacy ceiling when initialize requests
  the modern version. Pinned the SSE E2E client to
  `LATEST_LEGACY_PROTOCOL_VERSION` rather than the moving latest constant.
- Added `pkg/mcp/protocol_20260728_test.go`: raw JSON-RPC over real stdio pipes,
  using the application constructor, read-only Pulsar tenant registrations and
  an HTTP test backend. Covers discovery without initialization, tool/resource
  listing and reads, per-request capabilities, server identity, result type,
  private zero-TTL cache hints, removed-method rejection, unsupported versions,
  rejection before backend calls, and recovery after invalid requests.
- Stage 0 code is implemented but its full E2E gate remains blocked. Stage 1
  has an initial fixed-context read-only Pulsar tenant profile, not complete
  coverage of all backends, stateful tools, notifications, or cancellation.
  Stages 2–4 below remain planned.
- Passed with a consistent Go 1.25.14 toolchain: `go mod verify`,
  `go mod download`, `go fmt ./...`, `go mod tidy`,
  `golangci-lint run --timeout=3m` (0 issues), `go test -race ./...`
  (557 tests), `make build`, `make license-check`, and `git diff --check`.
  The shell initially mixed a Go 1.27.1 executable with a Go 1.25.14 GOROOT;
  validation used `PATH="$GOROOT/bin:$PATH"` to align them.
- `./scripts/e2e-test.sh all` reached Pulsar readiness, then remained blocked
  in Docker image construction for approximately 31 minutes before termination.
  Helm deployment and authenticated E2E assertions did not run. Re-run this
  gate in a working Docker environment before release. The test Pulsar
  container may remain; no broad environment cleanup was attempted.

## Executive conclusion

- `mark3labs/mcp-go` **v1.0.0 is a published, non-prerelease release**, published
  on 2026-09-02. Its release notes explicitly announce 2026-07-28 support.
  The release target commit is `51930cffc30e4708c60e9763a23ad0ed6ad06387` [S1].
- The official dated specification publishes substantial **breaking wire-protocol
  changes**, not simply a new initialization version string [S3–S5].
- The SDK contains modern discovery, request metadata, response decoration,
  subscriptions, MRTR, and Streamable HTTP handling. This is implementation
  evidence, **not proof of complete SDK or application conformance** [S2, S8–S11].
- The repository currently exposes **stdio and legacy HTTP+SSE**, not a
  Streamable HTTP endpoint. Upgrading the module does not add that endpoint,
  remove application state, implement OAuth discovery, or migrate dynamic tools.
- Recommend a dependency-focused compatibility release first, followed by an
  explicitly tested modern stdio/static-tool profile and a separate Streamable
  HTTP implementation. Keep state-dependent profiles legacy until migrated.
- Do not advertise the redesigned Tasks extension based on v1.0.0's existing
  task APIs: the inspected release still exposes the older task model and lacks
  `tasks/update` [S6, S12]. Optional-extension absence is not itself a failure
  to support the core protocol.

## Primary sources and evidence boundaries

All external references below are upstream release, source, or specification
material. SDK source links are pinned to the release tag; specification URLs
were read on the research date and may receive editorial updates. In case of
disagreement, normative dated specification/schema requirements take precedence
over SDK documentation or older repository proposals. No third-party compliance
claim is used as evidence.

| ID | Primary source | What it establishes |
| --- | --- | --- |
| S1 | [Release v1.0.0](https://github.com/mark3labs/mcp-go/releases/tag/v1.0.0), [release API](https://api.github.com/repos/mark3labs/mcp-go/releases/tags/v1.0.0), [v0.58.0 comparison](https://github.com/mark3labs/mcp-go/compare/v0.58.0...v1.0.0) | Exact release status, date, changes, and comparison baseline. |
| S2 | [Specification implementation PR #951](https://github.com/mark3labs/mcp-go/pull/951), [release-pinned SDK protocol guide](https://github.com/mark3labs/mcp-go/blob/v1.0.0/www/docs/pages/protocol-2026-07-28.mdx) | SDK implementation scope, compatibility helpers, and migration guidance. |
| S3 | [2026-07-28 specification](https://modelcontextprotocol.io/specification/2026-07-28), [key changes](https://modelcontextprotocol.io/specification/2026-07-28/changelog), [authoritative schema](https://github.com/modelcontextprotocol/modelcontextprotocol/blob/main/schema/2026-07-28/schema.ts) | Normative revision and changes from 2025-11-25; not an app compliance certificate. |
| S4 | [Versioning and compatibility](https://modelcontextprotocol.io/specification/2026-07-28/basic/lifecycle), [server/discover](https://modelcontextprotocol.io/specification/2026-07-28/server/discover) | Stateless request metadata and mandatory server discovery. |
| S5 | [Streamable HTTP](https://modelcontextprotocol.io/specification/2026-07-28/basic/transports/streamable-http), [subscriptions](https://modelcontextprotocol.io/specification/2026-07-28/basic/patterns/subscriptions), [MRTR](https://modelcontextprotocol.io/specification/2026-07-28/basic/patterns/mrtr) | HTTP headers, transport removal/migration rules, and new interaction patterns. |
| S6 | [Tasks overview](https://modelcontextprotocol.io/extensions/tasks/overview), [official Tasks repository](https://github.com/modelcontextprotocol/ext-tasks) | Optional extension negotiation, durable task handles, polling and mid-flight input. |
| S7 | [Deprecated features](https://modelcontextprotocol.io/specification/2026-07-28/deprecated), [feature lifecycle](https://modelcontextprotocol.io/community/feature-lifecycle), [authorization](https://modelcontextprotocol.io/specification/2026-07-28/basic/authorization) | Distinction between deprecated features and removed mechanisms; authorization obligations. |
| S8 | SDK [version.go](https://github.com/mark3labs/mcp-go/blob/v1.0.0/mcp/version.go), [protocol.go](https://github.com/mark3labs/mcp-go/blob/v1.0.0/server/protocol.go), [discover.go](https://github.com/mark3labs/mcp-go/blob/v1.0.0/server/discover.go), [response.go](https://github.com/mark3labs/mcp-go/blob/v1.0.0/server/response.go) | Modern version constant, metadata parsing, discovery and result/cache metadata. |
| S9 | SDK [streamable_http.go](https://github.com/mark3labs/mcp-go/blob/v1.0.0/server/streamable_http.go), [modern transport tests](https://github.com/mark3labs/mcp-go/blob/v1.0.0/server/streamable_http_modern_test.go) | Modern HTTP implementation and upstream test cases, not repository integration results. |
| S10 | SDK [client.go](https://github.com/mark3labs/mcp-go/blob/v1.0.0/client/client.go), [client protocol.go](https://github.com/mark3labs/mcp-go/blob/v1.0.0/client/protocol.go) | `Initialize` discovery/fallback behavior and version selection. |
| S11 | SDK [server.go](https://github.com/mark3labs/mcp-go/blob/v1.0.0/server/server.go), [tools.go](https://github.com/mark3labs/mcp-go/blob/v1.0.0/mcp/tools.go), [types.go](https://github.com/mark3labs/mcp-go/blob/v1.0.0/mcp/types.go) | Registration, schemas/results, error codes, capabilities and retained legacy APIs. |
| S12 | SDK [tasks.go](https://github.com/mark3labs/mcp-go/blob/v1.0.0/mcp/tasks.go), [types.go](https://github.com/mark3labs/mcp-go/blob/v1.0.0/mcp/types.go), [server.go](https://github.com/mark3labs/mcp-go/blob/v1.0.0/server/server.go) | Retained `tasks/list`/`tasks/result` model; extension-name support alone does not implement redesigned Tasks. |
| S13 | SDK [servertest/sse.go](https://github.com/mark3labs/mcp-go/blob/v1.0.0/server/servertest/sse.go), [servertest/streamable.go](https://github.com/mark3labs/mcp-go/blob/v1.0.0/server/servertest/streamable.go), [go.mod](https://github.com/mark3labs/mcp-go/blob/v1.0.0/go.mod) | Test-helper relocation and SDK Go 1.25.5 minimum. |

## Protocol changes and actual support boundaries

| Area | 2026-07-28 requirement/change | SDK v1.0.0 evidence | Application implication |
| --- | --- | --- | --- |
| Lifecycle | Remove `initialize`/`notifications/initialized` for modern requests. Every request carries `io.modelcontextprotocol/protocolVersion` and `io.modelcontextprotocol/clientCapabilities` in `params._meta`; client identity is SHOULD, as is server identity in result `_meta`. Servers MUST implement `server/discover`. | Modern parsing/discovery and client discovery-first `Initialize` facade are present [S8, S10]. | No custom handshake is needed, but test the actual transport. A successful Go `Initialize` call does not prove which wire era was used. |
| Protocol sessions | Remove `Mcp-Session-Id` from modern Streamable HTTP. Catalogs must not depend on connection state. Cross-call state uses explicit server-minted handles as ordinary arguments. | Modern HTTP path avoids protocol sessions; version restriction helper supports legacy-only deployments [S2, S9]. | Audit Cloud context mutation and session-specific Functions-as-tools before enabling modern support for those profiles. Backend connection caches are not inherently prohibited. |
| HTTP | Single POST endpoint; JSON or request-scoped SSE responses. No modern standalone GET stream, session DELETE, or replay with SSE event IDs/`Last-Event-ID`. Broken streams require reissuing requests with new request IDs. | Implemented in `NewStreamableHTTPServer`; modern GET/DELETE handling differs from legacy [S9]. | Current `NewSSEServer` endpoints do not become Streamable HTTP after an upgrade. Do not automatically retry destructive operations without application idempotency protection. |
| Routing headers | `MCP-Protocol-Version` must match body metadata; `Mcp-Method` is required for requests; `Mcp-Name` is required for `tools/call`, `resources/read`, and `prompts/get`. Header values must match body values, including prescribed Base64 decoding. | SDK has modern header handling and tests [S9–S11]. | Validate through the real ingress/middleware, including missing, mismatched, Unicode and sentinel-encoded values. No hard-coded claim that `Mcp-Name` is required on every method. |
| Custom headers | Servers MAY annotate primitive tool parameters with `x-mcp-header`; HTTP clients MUST support valid annotations and reject invalid tool definitions as specified. | SDK implementation must be checked against the full annotation/encoding rules, not merely type presence. | Not a requirement to add annotations to existing tools. If introduced, test schema validity, header forwarding and injection resistance. |
| Results/cache | Required `resultType`; ordinary results are `complete`, MRTR interim results are `input_required`. Missing legacy result type means `complete`. `ttlMs` and `cacheScope` are required on five list/read methods named below. | Response decoration provides modern result metadata and default `ttlMs: 0`, `cacheScope: private` [S8]. | Validate serialized responses, especially custom handlers. Keep private/revalidate defaults until visibility and invalidation semantics are proven. |
| Notifications | `subscriptions/listen` replaces standalone GET and `resources/subscribe`/`unsubscribe`; opt-in filters and subscription IDs scope change notifications. Progress/log messages remain on their originating request stream. | SDK guide and implementation expose subscription support [S2, S9]. | Existing resource capability flags are not proof that application changes generate correct notifications. Test supported filters and delivery with actual resources/dynamic tools. |
| MRTR | Independent server-initiated requests are replaced by `InputRequiredResult`, `inputRequests`, retry `inputResponses`, and opaque `requestState`. | SDK builder/client helpers exist; direct `RequestSampling`, `RequestElicitation`, `RequestRoots` reject modern requests [S2]. | No such direct calls were found in inspected application `pkg`/`cmd` Go sources. Do not invent a feature dependency; if later used, validate continuation integrity, authorization and side-effect idempotency. |
| Removed/deprecated | Remove `ping`, `logging/setLevel`, `notifications/roots/list_changed`, `notifications/elicitation/complete`, and URL elicitation `elicitationId`. Logging level moves to per-request `_meta`; no requested level means no `notifications/message`. Roots, Sampling, Logging and HTTP+SSE are deprecated, not all removed features. | Legacy APIs remain for compatibility; SDK log helper enforces modern request-level filtering [S2]. | Retain OS-level health probes; do not substitute protocol ping. Review `WithLogging` and resource capability advertising against actual behavior. Logrus operational logging is distinct from MCP log notifications. |
| Tasks | Experimental core tasks move to opt-in `io.modelcontextprotocol/tasks`. New `resultType: task`, `tasks/get` includes terminal result/error; `tasks/update` supplies input; remove `tasks/list` and blocking `tasks/result`. | Release still contains old task APIs. Source search found no `tasks/update` implementation; the extension name occurs as an example, not an implementation [S12]. | Leave unadvertised. A task-polling helper in the release notes is not evidence of new extension support. |
| Schemas/errors | Allow JSON Schema 2020-12 keywords and any JSON structured content; add reference/composition safety requirements. Resource-not-found becomes `-32602`; header mismatch `-32020`, missing capability `-32021`, unsupported version `-32022`. | Raw schema paths, `structuredContent` as `any`, and revised constants are present [S11]. | Retest schema serialization and error envelopes. Type flexibility does not establish bounded `$ref` resolution or full schema-validator conformance. |
| Authorization/telemetry | Updated issuer validation and issuer-bound client credentials; DCR `application_type`; DCR is deprecated in favor of Client ID Metadata Documents. Trace context conventions use `_meta` `traceparent`, `tracestate`, `baggage`. | Not established by the release's protocol-support headline. | Backend OAuth and bearer-token forwarding are not automatically an MCP authorization resource-server implementation. Trace propagation requires application integration if adopted. |

The five cacheable methods are `tools/list`, `prompts/list`, `resources/list`,
`resources/read`, and `resources/templates/list`. Deterministic tool ordering is
a SHOULD; SDK sorting is present, but dynamic registration and pagination need
application-level stability tests. Authentication-based visibility must not be
confused with connection-local catalog mutation or an authorization-free public
cache. Do not enable public caching for identity-sensitive results.

## Repository findings

The inspected baseline has `mcp-go v0.58.0` in committed `go.mod`. During research,
the shared working tree already had `go.mod`/`go.sum` changes selecting v1.0.0,
owned by the upgrade task. This document does not interpret those edits as a
completed upgrade or a passed test suite. SDK inspection used the downloaded
`github.com/mark3labs/mcp-go@v1.0.0` module and release-pinned upstream sources.

- `pkg/cmd/mcp/stdio.go` uses `server.NewStdioServer` and supplies application
  sessions through context. It can inherit modern SDK dispatch, but no separate
  modern wire-conformance result was established by this research.
- `pkg/cmd/mcp/sse.go` uses `server.NewSSEServer`, registering
  `<httpPath>/sse`, `<httpPath>/message`, `<httpPath>/healthz`, and
  `<httpPath>/readyz`. It does **not** register a Streamable HTTP POST handler at
  `<httpPath>`. SSE framing inside modern HTTP is not the old HTTP+SSE transport.
- Multi-session external Pulsar mode checks bearer tokens and resolves a
  backend session through `pkg/mcp/session/`; health endpoints are deliberately
  outside that middleware. This cache can remain an implementation detail if
  authorization and routing are reconstructed on every modern request.
- `pkg/mcp/server.go` defaults to `WithResourceCapabilities(true, true)`,
  `WithRecovery`, and `WithLogging`. Capability claims require behavioral tests,
  not only successful server construction.
- `pkg/mcp/sncontext_tools.go` calls `SetContext` in
  `pkg/mcp/sncontext_utils.go`, mutating in-memory connection context. A preceding
  context-switch call cannot become a hidden routing dependency for modern
  stateless requests. Prefer fixed startup context for the first modern profile;
  later introduce explicit, authorization-bound context handles if needed.
- `pkg/mcp/pulsar_functions_as_tools.go` and `pkg/mcp/pftools/manager.go` depend
  on session IDs and `AddSessionTool`/`DeleteSessionTools`. SSE currently passes
  `FIXED_SESSION_ID`. Do not assume removing protocol sessions preserves this
  dynamic catalog. Either keep the profile legacy or redesign registration and
  call routing independently of connection identity.
- `cmd/snmcp-e2e/main.go` constructs `NewSSEMCPClient` and uses
  `LATEST_PROTOCOL_VERSION` in `Initialize`. In v1.0.0 that constant changes to
  `2026-07-28`; the SDK facade can discover or fall back. Existing SSE success
  is not evidence of modern Streamable HTTP support. Assert the returned version
  and inspect wire behavior in separate modern and legacy cases.
- Several older `agents/mcp-compliant-*` documents, `mcp-compliance-guide.md`,
  and `documentation-update-summary.md` describe `/mcp` Streamable HTTP endpoints
  as present. The inspected router contradicts those claims. Treat them as
  historical designs, not deployed implementation evidence; reconcile them in
  the transport implementation task. Only this new document is changed here.

## Breaking migration requirements

There are two different compatibility boundaries:

1. **Go source/build compatibility.** v1.0.0 moves
   `server.NewTestServer` and `server.NewTestStreamableHTTPServer` to
   `server/servertest` [S13]. Update callers if any; ordinary `net/http/httptest`
   imports are unrelated. SDK Go 1.25.5 is below this repository's Go 1.25.12.
   Compile/test the whole repository rather than assuming the SDK guide's
   “existing code keeps working unchanged” covers all exported APIs.
2. **Wire and application semantics.** Modern peers cannot depend on the old
   handshake, transport session ID, session tool catalogs, standalone GET
   notifications, stream replay, or independent server RPCs. Compatibility
   shims in the SDK do not migrate application state or infrastructure.

For a future Streamable HTTP endpoint that cannot yet serve modern requests,
the SDK documents
`server.WithStreamableHTTPProtocolVersions(mcp.LegacyProtocolVersions()...)`.
That option applies to Streamable HTTP, **not** the existing SSE constructor.
For stdio/SSE legacy restrictions, verify supported-version context/dispatch
behavior and tests rather than copying an option from another transport.
Do not rely solely on client pinning to prevent an unsupported server profile
from advertising modern support.

## Staged implementation and acceptance gates

### Stage 0 — dependency upgrade and honest compatibility baseline

Upgrade only the dependency and any necessary API callers. Preserve deployed
stdio/SSE CLI behavior. Record exact SDK and negotiated protocol versions.
Keep release wording limited to “updated mcp-go to v1.0.0” until later gates pass.

Acceptance: repository dependency verification, formatting/tidy, lint, race
tests, build, and representative existing stdio/SSE authentication checks pass.
Test a pinned legacy client as well as the new SDK client; identify accidental
modern advertisement for state-dependent profiles. Do not add broad modern
support claims based on dependency or compile success.

### Stage 1 — bounded modern stdio support

Start with a fixed backend context and connection-independent static tool set
(for example external Kafka/Pulsar). Test raw JSON-RPC, not only same-SDK calls:

- `server/discover`, modern `tools/list` and representative `tools/call` without
  initialization; missing/malformed metadata and unsupported-version errors.
- Required result metadata and cache fields for every applicable method; legacy
  result compatibility; deterministic ordering and pagination.
- Read-only feature gating, identity/cluster isolation, cancellation/progress
  where supported, and no unsolicited protocol logging without requested level.

Acceptance: publish exactly the tested transport/features/protocol versions.
Do not include context switching or Functions-as-tools until Stage 3.

### Stage 2 — modern Streamable HTTP integration

Add a deliberate HTTP command/endpoint/configuration design using
`NewStreamableHTTPServer`, while retaining documented legacy SSE endpoints for
existing users. Reuse typed context helpers and per-request authorization;
retain explicit unauthenticated health probes. Review auth requirements
separately from backend token validation.

Acceptance: modern JSON and SSE response paths, discovery/direct first calls,
header/body validation and error-to-HTTP mapping, no modern protocol session
header, appropriate GET/DELETE behavior, no stream replay, and independently
pinned legacy interoperability. Verify Origin validation (invalid present
Origin returns 403), localhost defaults where applicable, proxy header
forwarding, TLS/auth boundaries, and rejection before tool side effects.
Exercise two identities and multiple replicas without sticky sessions. Update
README, chart values/templates, E2E clients/scripts, and stale design documents
as part of that implementation, not as an assumption in this research.

### Stage 3 — state, dynamic catalogs and notifications

Define explicit server-minted handles for any cross-call state, bound to
principal/tenant/cluster with expiry and authorization revalidation. Alternatively
keep fixed-context deployments as the documented supported profile. Migrate
Functions-as-tools away from connection-keyed catalogs; preserve read-only
filtering and stable discovery. Internal backend caches must be disposable and
must not be the only source of authorization/routing truth.

Acceptance: connection/replica changes do not change catalog or target context
unexpectedly; concurrent users cannot affect each other's state; expired or
foreign handles are rejected. Validate `subscriptions/listen` filters, IDs,
cancellation and resource/tool changes. Keep cache hints private with zero TTL
until invalidation behavior is proven. Define safe retry/idempotency behavior
for mutations after stream loss; a fresh JSON-RPC ID does not deduplicate work.

### Stage 4 — optional extensions and broader feature coverage

Add MRTR only for a concrete interactive operation, with bounded retries and
authenticated continuation state. Add Tasks only after the SDK or a deliberate
adapter implements the redesigned extension and durable execution semantics;
explicitly negotiate the extension and test restart/reconnect, authorization,
polling, `tasks/update`, cancellation and terminal results. Do not advertise
deprecated Roots/Sampling/Logging as new product capabilities. Evaluate trace
context and authorization additions against actual deployment needs.

Acceptance: a versioned feature matrix and raw-wire/independent-client tests
cover every advertised feature. Unsupported optional extensions stay absent.

## Validation ownership and remaining uncertainty

This research inspected sources and official pages; it did **not** run builds,
unit tests, network conformance tests, or a full SDK audit. The main implementation
task owns `go mod verify`, `go mod download`, `go fmt ./...`, `go mod tidy`,
`golangci-lint run --timeout=3m`, `go test -race ./...`, and `make build`. Chart
or E2E changes additionally require `./scripts/e2e-test.sh all`; release packaging
changes require the repository's GoReleaser checks.

In particular, full authorization conformance, JSON Schema resolution limits,
all custom-header edge cases, subscriptions over each exposed transport, and
stateful tool behavior under modern negotiation remain unverified. Resolve them
with targeted tests and explicit support restrictions, not a blanket claim of
“MCP 2026-07-28 compliant.”
