# Plan: Claude connector tool split + annotations

## Goal

Prepare StreamNative MCP Server for Claude connector submission and review.

## Current review follow-up: annotation helper granularity

Reviewer feedback is valid: `toolannotations.ReadOnly` and `toolannotations.Destructive` currently call `mcp.WithToolAnnotation(...)`, replacing the whole annotation struct and dropping mcp-go defaults for `idempotentHint` and `openWorldHint`. The helper also collapses all side effects into `Destructive`, so local session mutations and reversible lifecycle operations look equivalent to delete/apply operations.

Implementation status: implemented in current work.

1. Replaced whole-annotation assignment helpers with composed field setters so unspecified hints keep mcp-go defaults unless explicitly changed.
2. Exposed helper categories:
   - `ReadOnly(title)` for non-mutating external reads: read-only true, destructive false, idempotent true, open-world true.
   - `ExternalRead(title)` as explicit alias for external read tools when caller intent should be obvious.
   - `LocalSessionMutation(title)` for session/context state changes: read-only false, destructive false, idempotent false by default, open-world false.
   - `Mutating(title, destructive, idempotent)` for external writes/side effects: read-only false, caller-controlled destructive/idempotent, open-world true.
3. Kept `Destructive(title)` as compatibility wrapper over `Mutating(title, true, false)` while migrating call sites to precise helpers.
4. Updated context tools to use `LocalSessionMutation`; updated operation-mode helper to derive destructive/idempotent from `OperationSpec` where available instead of treating all writes as destructive.
5. Extended annotation tests to assert `idempotentHint` and `openWorldHint` for static/context/builder tools, not only read/destructive.

Risks:

- Tool annotation surface changes are runtime-visible to MCP clients and Claude review.
- Incorrect idempotency classification can mislead clients; delete/remove/reset may be idempotent only depending backend behavior.
- Broad migration across all builders risks noisy diff; recommended first patch helper + high-signal call sites, then operation registry cleanup separately.

Recommended validation:

```bash
go test -race ./pkg/mcp/... ./pkg/mcp/builders/... ./pkg/mcp/pftools/...
go test -race ./...
go fmt ./...
go mod tidy
golangci-lint run --timeout=3m
make build
```

## Current review follow-up: operation spec registry

Reviewer feedback is valid: current branch split read/write tool names, but duplicated `toolMode` helpers and per-builder write-operation maps still create scatter-shot maintenance. Adding one operation can require enum updates, write map updates, handler switch updates, docs, and compliance-test classification. `validateModeOperation` also classifies any operation missing from the write map as read, so an unclassified future write can pass read-mode validation until the handler switch rejects or mishandles it.

Recommended next design: make each tool family declare operation metadata once, then derive mode-specific tool schemas, annotations, validation, and tests from that registry.

```go
type OperationMode string

const (
    OperationModeRead  OperationMode = "read"
    OperationModeWrite OperationMode = "write"
)

type OperationSpec struct {
    Name        string
    Mode        OperationMode
    Destructive bool
    Idempotent  bool
    Resources   []string
    Params      []ParamSpec
    Handler     OperationHandler
}
```

Scope is incremental but complete: add shared spec helpers, migrate `kafka/topics.go` and `pulsar/namespace.go` as reference implementations, then migrate every remaining split builder in batches. Keep current read/write tool names unchanged. Use docs generated operation-table blocks first; do not rewrite full Markdown documents.

Hard requirements from Claude docs:

- every MCP tool has non-empty `annotations.title`
- every MCP tool has explicit applicable `annotations.readOnlyHint` or `annotations.destructiveHint`
- read and write operations must be separate tools; no mixed `operation` catch-all that contains both safe and unsafe operations

Source docs checked:

- https://claude.com/docs/connectors/building/submission
- https://claude.com/docs/connectors/building/review-criteria

Reference implementation checked:

- `/Users/rui/playground/sn/mcp-auth0-proxy/internal/hooks/org_session_tools.go`
- Existing pattern there:
  - separate tool names: e.g. `sncloud_byoc_read` and `sncloud_byoc_write`
  - shared builder with mode enum: `controlPlaneToolModeRead` / `controlPlaneToolModeWrite`
  - read tool operation enum: `list`, `get`
  - write tool operation enum: `apply`, `delete`
  - annotation set from mode:
    - read: `readOnlyHint=true`, `destructiveHint=false`
    - write: `readOnlyHint=false`, `destructiveHint=true`
  - shared handler still validates operation against mode

## Current findings

Current follow-up focus: previous read/write split only separated tool names and operation enums in some builders. It still leaves mixed mode descriptions and write-only schema fields visible on read tools. Examples: `pkg/mcp/builders/kafka/topics.go` and `pkg/mcp/builders/pulsar/namespace.go`; same class can exist in other split builders. Connector review can still treat this as a mixed surface because `tools/list` exposes write verbs/examples/parameters through read tools.

Static `mcp.NewTool(...)` definitions found under `pkg/`: 36 tool definitions plus dynamic Pulsar Functions-as-Tools.

Current gaps:

- Most static tools have no explicit title.
- Most static tools rely on `mcp-go` defaults (`readOnlyHint=false`, `destructiveHint=true`, `openWorldHint=true`), which marks read tools as destructive.
- Only `sncloud_resources_apply` and `sncloud_resources_delete` currently set `WithToolAnnotation`; `apply` sets title only.
- Dynamic Pulsar Functions-as-Tools in `pkg/mcp/pftools/manager.go` create tools without title/read-only/destructive annotations.
- Many admin tools multiplex read and write operations through one `operation` parameter. Claude review criteria says mixed read/write catch-all tools can be rejected even if description documents safe/unsafe operations.
- Some already-split tools still have mixed descriptions and schemas. Mode-specific enum is not enough; read tools must not expose write operations, examples, or write-only parameters.

## Proposed design

### Design principle

Follow `mcp-auth0-proxy` pattern:

- split mixed tools into separate read and write tool names
- keep shared internal implementation where practical
- make operation enum, description, examples, and parameters mode-specific, so tool schema itself prevents mixed use
- keep read-only runtime mode simple: register only read tools
- in read-write runtime mode: register read tools and write tools as separate entries
- do not expose legacy mixed tools in Claude-submitted surface

### Naming convention

For mixed tools, replace one tool with two tools:

- `<old_name>_read`
- `<old_name>_write`

Examples:

- `kafka_admin_topics` -> `kafka_admin_topics_read`, `kafka_admin_topics_write`
- `pulsar_admin_topic` -> `pulsar_admin_topic_read`, `pulsar_admin_topic_write`
- `pulsar_admin_namespace_policy` already has partial split; align names and annotations instead of forcing one exact pattern when current names are already narrow.

Pure read tools can keep current names if no write side effects exist.
Pure write/side-effect tools can keep current names if description and annotation are clear.

Compatibility policy:

- Remove mixed legacy tool registration from default surface.
- Do not add opt-in legacy aliases or compatibility flags for old mixed tool names.
- Do not keep mixed legacy tools visible in submitted connector, even with destructive annotation.

### Shared helper APIs

Add a small helper package, likely `pkg/mcp/toolannotations`, to avoid duplicated pointer boilerplate and import cycles:

- `ReadOnly(title string) mcp.ToolOption` -> title, `readOnlyHint=true`, `destructiveHint=false`
- `Destructive(title string) mcp.ToolOption` -> title, `readOnlyHint=false`, `destructiveHint=true`
- optional `NonDestructiveWrite(title string)` only if a tool changes local session state without modifying external service; use sparingly because Claude requirement names `destructiveHint` for modifying/deleting tools.

Add builder-local mode types where useful:

```go
type toolMode string

const (
    toolModeRead  toolMode = "read"
    toolModeWrite toolMode = "write"
)
```

Build functions should accept mode:

- `buildTool(mode toolMode)`
- `buildHandler(mode toolMode, readOnly bool)`
- `validateOperation(mode, operation)`
- `isWriteOperation(operation)`

Mode-specific tool builders should also split:

- `toolDesc`
- `resourceDesc` when resource meanings differ by mode
- `operationDesc`
- `operationEnum`
- parameter set (`mcp.WithString`, `mcp.WithNumber`, `mcp.WithObject`, etc.)

Read tools must not expose write-only fields. Write tools should not expose read-only-only fields unless a write operation genuinely needs them.

## Split inventory

### Kafka builders

#### `kafka_admin_topics`

Split:

- `kafka_admin_topics_read`
  - operations: `list`, `get`, `metadata`
  - annotation: read-only
- `kafka_admin_topics_write`
  - operations: `create`, `delete`
  - annotation: destructive

Read-only runtime: register read only.
Read-write runtime: register both.

#### `kafka_admin_groups`

Split:

- `kafka_admin_groups_read`
  - operations: `list`, `describe`, `offsets`
  - annotation: read-only
- `kafka_admin_groups_write`
  - operations: `remove-members`, `delete-offset`, `set-offset`
  - annotation: destructive

#### `kafka_admin_sr`

Split:

- `kafka_admin_sr_read`
  - operations: `list`, `get`, plus schema type/capability read operations
  - annotation: read-only
- `kafka_admin_sr_write`
  - operations: `set`, `create`, `delete`
  - annotation: destructive

#### `kafka_admin_connect`

Split:

- `kafka_admin_connect_read`
  - read operations: cluster info, connector list/get/status/config, connector plugins, transforms
  - annotation: read-only
- `kafka_admin_connect_write`
  - write operations: `create`, `update`, `delete`, `restart`, `pause`, `resume`
  - annotation: destructive

#### `kafka_admin_partitions`

Current tool appears write-only (`update`). Options:

- keep `kafka_admin_partitions` as destructive write-only; or
- rename to `kafka_admin_partitions_write` for consistency.

Recommendation: rename to `kafka_admin_partitions_write` if no read operations exist, and update docs. If preserving name matters, keep current name but annotate destructive.

#### `kafka_client_produce`

Write/side-effect tool. Keep current name, annotate destructive.

#### `kafka_client_consume`

Ambiguous:

- description says no offset commit unless `group` parameter is explicitly specified
- with `group`, consumer group state may change

Recommendation for review safety:

- split into `kafka_client_consume_read` without `group` / no offset commit
- optional `kafka_client_consume_group` or `kafka_client_consume_write` for group-based consumption that may affect offsets/state, annotated destructive

If implementation never commits offsets, keep single read tool after code verification and adjust description/schema to remove side-effect ambiguity.

### Pulsar builders

#### `pulsar_admin_topic`

Split:

- `pulsar_admin_topic_read`
  - operations: `list`, `get`, `get-permissions`, `stats`, `lookup`, `internal-stats`, `internal-info`, `bundle-range`, `last-message-id`, `compact-status`, `offload-status`
  - annotation: read-only
- `pulsar_admin_topic_write`
  - operations: `grant-permissions`, `revoke-permissions`, `create`, `delete`, `unload`, `terminate`, `compact`, `update`, `offload`
  - annotation: destructive

#### `pulsar_admin_subscription`

Split:

- `pulsar_admin_subscription_read`
  - operations: `list`, `peek`, `get-message-by-id`
  - annotation: read-only
- `pulsar_admin_subscription_write`
  - operations: `create`, `delete`, `skip`, `expire`, `reset-cursor`
  - annotation: destructive

#### `pulsar_admin_namespace`

Split:

- `pulsar_admin_namespace_read`
  - operations: `list`, `get_topics`
  - annotation: read-only
- `pulsar_admin_namespace_write`
  - operations: `create`, `delete`, `clear_backlog`, `unsubscribe`, `unload`, `split_bundle`
  - annotation: destructive

#### `pulsar_admin_namespace_policy*`

Already partly separated:

- `pulsar_admin_namespace_policy_get` -> read-only
- `pulsar_admin_namespace_policy_get_anti_affinity_namespaces` -> read-only
- `pulsar_admin_namespace_policy_set` -> destructive
- `pulsar_admin_namespace_policy_remove` -> destructive

Keep split; add titles/annotations and ensure no tool mixes set/remove/get.

#### `pulsar_admin_topic_policy`

Likely mixed get/set/remove operations. Split:

- `pulsar_admin_topic_policy_read`
- `pulsar_admin_topic_policy_write`

Use same operation partitioning as handler supports.

#### `pulsar_admin_brokers`

Split:

- `pulsar_admin_brokers_read`
  - list/get/health/config/namespaces/runtime/internal/all_dynamic reads
  - annotation: read-only
- `pulsar_admin_brokers_write`
  - dynamic config update/delete or any mutable broker operation
  - annotation: destructive

#### `pulsar_admin_cluster`

Split:

- `pulsar_admin_cluster_read`
  - `list`, `get`, read peer/failure-domain operations
  - annotation: read-only
- `pulsar_admin_cluster_write`
  - `create`, `update`, `delete`, write peer/failure-domain operations
  - annotation: destructive

#### `pulsar_admin_functions`

Split:

- `pulsar_admin_functions_read`
  - `list`, `get`, `status`, `stats`, `querystate`, `download`
  - annotation: read-only
- `pulsar_admin_functions_write`
  - `create`, `update`, `delete`, `start`, `stop`, `restart`, `putstate`, `trigger`, `upload`
  - annotation: destructive

#### `pulsar_admin_sinks` / `pulsar_admin_sources`

Split each:

- `*_read`
  - `list`, `get`, `status`, `list-built-in`
  - annotation: read-only
- `*_write`
  - `create`, `update`, `delete`, `start`, `stop`, `restart`
  - annotation: destructive

#### `pulsar_admin_packages`

Split:

- `pulsar_admin_package_read`
  - `list`, `get`, `download`
  - annotation: read-only
- `pulsar_admin_package_write`
  - `update`, `delete`, `upload`
  - annotation: destructive

#### `pulsar_admin_schema`

Split:

- `pulsar_admin_schema_read`
  - `get`
  - annotation: read-only
- `pulsar_admin_schema_write`
  - `upload`, `delete`
  - annotation: destructive

#### `pulsar_admin_tenant`

Split:

- `pulsar_admin_tenant_read`
  - `list`, `get`
  - annotation: read-only
- `pulsar_admin_tenant_write`
  - `create`, `update`, `delete`
  - annotation: destructive

#### `pulsar_admin_nsisolationpolicy`

Split:

- `pulsar_admin_nsisolationpolicy_read`
  - `get`, `list`, broker read operations
  - annotation: read-only
- `pulsar_admin_nsisolationpolicy_write`
  - `set`, `delete`
  - annotation: destructive

#### `pulsar_admin_resourcequota`

Split:

- `pulsar_admin_resourcequota_read`
  - `get`
  - annotation: read-only
- `pulsar_admin_resourcequota_write`
  - `set`, `reset`
  - annotation: destructive

#### Pure read tools

Keep current names, add read-only annotation:

- `pulsar_admin_status`
- `pulsar_admin_broker_stats`
- `pulsar_admin_functions_worker`
- any MCP resources/templates that are not tools stay out of tool annotation scope

#### Pulsar client tools

- `pulsar_client_produce`: keep current name, destructive annotation.
- `pulsar_client_consume`: likely side-effectful because subscriptions/cursors can be created/advanced. Recommendation: annotate destructive unless implementation is changed to provide a non-mutating peek/read variant.

Possible future split:

- `pulsar_client_peek_read` for non-destructive peeking if supported by admin APIs
- `pulsar_client_consume` remains destructive

### StreamNative Cloud tools

#### Existing resource tools

Already split by action:

- `sncloud_resources_apply`: destructive; include title and `destructiveHint=true`
- `sncloud_resources_delete`: destructive; title already present, ensure readOnlyHint false too

No read counterpart currently. If resource list/get is added, use `sncloud_resources_read` rather than adding list/get to apply/delete tools.

#### Context tools

- `sncloud_context_whoami`: read-only
- `sncloud_context_available_clusters`: read-only
- `sncloud_context_use_cluster`: session/context mutation; annotate destructive or non-read-only. For Claude safety, use destructive unless we explicitly add `NonDestructiveWrite` and verify review accepts it.
- `sncloud_context_reset`: session/context mutation; annotate destructive or non-read-only. For Claude safety, use destructive.

#### Logs

- `sncloud_logs`: read-only

### Dynamic Functions-as-Tools

`pkg/mcp/pftools/manager.go` dynamic tools invoke deployed Pulsar Functions and can produce messages / trigger external effects.

Plan:

- keep dynamic tool name from function metadata
- add human-readable title from function metadata/tool name
- annotate `destructiveHint=true`
- if read-only mode should not expose dynamic invocation tools, verify registration path and add test
- do not mark dynamic function tools read-only unless function metadata explicitly supports safe read-only classification in future

## Implementation phases

### Phase 1: shared annotation + mode helpers

Status: implemented on current branch.

- Add `pkg/mcp/toolannotations` helper.
- Add local read/write mode helpers in builders with mixed operations.
- Add reusable operation validation helpers where a builder already has operation maps.

### Phase 2: split Kafka tools completely

Status: implemented on current branch; follow-up refactor still needed to remove duplicated operation classification.

- Update all Kafka builders to build mode-specific tools.
- Ensure read/write tools have mode-specific descriptions, examples, operation enums, and parameter schemas.
- Read-only config returns only read tools.
- Read-write config returns read + write tools, except pure write tools remain write-only.
- Remove old mixed tool surface; no compatibility alias.
- Update wrapper tests/docs.

### Phase 3: split Pulsar tools completely

Status: implemented on current branch; follow-up refactor still needed to remove duplicated operation classification.

- Update all Pulsar builders to build mode-specific tools.
- Ensure read/write tools have mode-specific descriptions, examples, operation enums, and parameter schemas.
- Preserve existing read-only behavior by not registering write tools in read-only config.
- Ensure mode-specific operation enums and validation errors.
- Remove old mixed tool surface; no compatibility alias.
- Add/extend parity tests for operation coverage.

### Phase 4: StreamNative Cloud/static tool annotations

Status: implemented on current branch.

- Add annotations to context/log/resource tools.
- Keep already split apply/delete tools.
- Ensure no new mixed resource tool appears.

### Phase 5: dynamic tools

Status: implemented on current branch.

- Add annotations to Functions-as-Tools.
- Validate read-only exposure behavior.

### Phase 6: runtime-visible docs

Status: implemented on current branch, with generated operation-table guard tests for migrated split builders.

Update runtime-visible docs together with schema changes:

- Keep current read/write tool names unchanged.
- `README.md` feature/tool examples only if behavior changes.
- `docs/tools/*.md` matching current split tools.
- Split docs into explicit read/write sections where a family has both tool modes.
- Ensure read docs do not mention write-only operations or parameters.
- Ensure write docs do not rely on old mixed tool names.
- Any design notes under `agents/` if tool surface changes are architectural.

### Phase 7: shared operation spec registry follow-up

Status: implemented on current branch.

Goal: make operation metadata single-source-of-truth and remove copy-paste `toolMode` + `writeOperations` maps.

1. Add shared operation metadata API, likely under `pkg/mcp/builders/operations.go`:
   - `OperationModeRead` / `OperationModeWrite`
   - `OperationSpec`
   - `ParamSpec`
   - `OperationHandler` type alias or adapter
   - `OperationRegistry` or helper functions over `[]OperationSpec`
2. Shared helpers must derive:
   - read operation enum
   - write operation enum
   - operation description fragments/table rows
   - mode-specific validation
   - unknown-operation rejection
   - read/write annotation selection
   - compliance-test classification
3. Validation semantics:
   - operation absent from spec => reject, never default to read
   - operation present with wrong mode => reject with mode-specific error
   - operation present with matching mode => dispatch allowed
4. Migrate two reference builders first:
   - `pkg/mcp/builders/kafka/topics.go`
   - `pkg/mcp/builders/pulsar/namespace.go`
5. Reference builder acceptance criteria:
   - no local `xxxWriteOperations` map
   - operation enum generated from specs
   - handler validation uses specs
   - unknown operation test added
   - read-only build still excludes write tools
   - tool names unchanged
   - docs operation table generated or checked from specs
6. Migrate remaining Kafka split builders:
   - `connect.go`
   - `groups.go`
   - `schema_registry.go`
   - `topics.go`
7. Migrate remaining Pulsar split builders:
   - `brokers.go`
   - `cluster.go`
   - `functions.go`
   - `namespace.go`
   - `nsisolationpolicy.go`
   - `packages.go`
   - `resourcequotas.go`
   - `schema.go`
   - `sinks.go`
   - `sources.go`
   - `subscription.go`
   - `tenant.go`
   - `topic.go`
   - `topic_policy.go`
8. Keep pure read/write client tools out of forced registry migration unless helper reuse is cheap:
   - `kafka_client_produce`
   - `kafka_client_consume`
   - `pulsar_client_produce`
   - `pulsar_client_consume`
   - pure read Pulsar admin status/stats/worker tools
9. Replace compliance-test manual classification:
   - derive write/read operation sets from specs
   - assert no enum mixes read/write specs
   - assert every enum value exists in specs
   - assert specs cover every handler switch operation
10. Docs generation/check:
   - do not rewrite whole Markdown documents
   - add generated operation table blocks only:
     `<!-- generated:operations:start -->` / `<!-- generated:operations:end -->`
   - add `go generate` or test helper to refresh/check blocks
11. Cleanup after migration:
   - delete or shrink duplicated Kafka/Pulsar `tool_mode.go`
   - move shared `pruneToolInputSchema`/required filtering helper if still duplicated
   - remove obsolete per-builder write maps
   - remove obsolete compliance-test write maps

## Tests / compliance guard

Add focused tests:

- For every builder under `pkg/mcp/builders/kafka` and `pkg/mcp/builders/pulsar`:
  - no returned tool mixes read and write operations
  - tool name length <= 64
  - title non-empty
  - read-only or destructive hint explicit
  - read tools: `ReadOnlyHint=true`, `DestructiveHint=false`
  - write tools: `DestructiveHint=true`, `ReadOnlyHint=false`
  - read-only config returns no write tools
  - read tools do not expose known write-only parameters
  - read tool descriptions do not mention known write-only operations, examples, or destructive verbs for that family
  - write tool schemas do not expose read-only-only parameters unless genuinely shared
- StreamNative Cloud/context/log/resource tools have valid annotations.
- PFTools dynamic tool creation has valid annotation.
- Operation validation rejects read operations on write tools and write operations on read tools.
- Operation validation rejects unknown operations instead of treating them as read.
- Operation enum values are derived from or checked against `OperationSpec`.
- Compliance tests derive read/write classification from `OperationSpec`, not hand-maintained write maps.
- Docs generated operation table blocks match `OperationSpec`.

Static guard:

- Build all feature sets and assert no `operation` enum contains both read and write verbs in one tool.
- For split tool families, assert mode-specific schema/description purity with family-specific allow/deny lists.
- Assert every handler switch operation has a matching `OperationSpec` entry.

## Risks

- Tool split is runtime-visible and likely breaking for clients/prompts that call old names.
- Operation spec refactor can accidentally change schema ordering, descriptions, or enum content even when tool names stay unchanged.
- Docs under `docs/tools/` can drift if generated operation blocks are not checked in tests or `go generate`.
- Some operations are ambiguous (`consume`, `trigger`, cursor operations, context reset). Conservative destructive annotation may add confirmations but avoids unsafe auto-run.
- Some current tools may have read-only-mode logic embedded in handlers; after split, registration and handler validation must both enforce mode to prevent write leakage.
- `OperationSpec.Handler` can over-couple schema metadata and dispatch if introduced too early. Prefer enum/validation/test/doc generation first, then dispatch consolidation.
- `mcp-go` default annotations are unsafe for compliance because title empty and destructive default true.

## Confirmed decisions

- Fix all current mixed read/write surfaces, not only `kafka/topics.go` and `pulsar/namespace.go`.
- Implement operation-spec follow-up across all affected split builders, not only the two reference files.
- Start with shared spec + two reference migrations: `kafka/topics.go` and `pulsar/namespace.go`.
- Keep current read/write tool names unchanged.
- Do not preserve old mixed tool names or old mixed builder/schema patterns.
- Runtime-visible docs must be updated with read/write split and must avoid mixed read/write wording.
- For docs generation, start with generated operation table blocks only; do not generate whole Markdown documents.
- Conservative safety annotations are acceptable for ambiguous side-effect tools unless implementation proves true read-only behavior.

## Recommended validation

Fast local:

```bash
go test ./pkg/mcp/... ./pkg/schema/...
go test -race ./pkg/mcp/builders/...
go fmt ./...
go mod tidy
```

Full repo before PR:

```bash
go mod verify
go mod download
golangci-lint run --timeout=3m
go test -race ./...
make build
make license-check
```

Connector-specific manual check:

```bash
bin/snmcp stdio --use-external-pulsar --pulsar-web-service-url http://localhost:8080
# then inspect tools/list with MCP Inspector and verify every tool annotation and no mixed read/write operation enum
```

Chart/E2E only needed if chart, SSE auth, or e2e harness touched:

```bash
./scripts/e2e-test.sh all
```
