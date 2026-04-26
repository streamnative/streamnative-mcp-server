# Base Contract

You are executing one Ralph Loop iteration for Pulsar MCP resources in
`streamnative-mcp-server`.

Implement exactly the selected backlog item unless you discover a correctness
blocker. Make the smallest complete change that satisfies the item acceptance
criteria. Prefer existing repository patterns over new abstractions.

## Required behavior

- Read the relevant code before editing.
- Keep all resource handlers read-only.
- Reuse the existing Pulsar session context helpers.
- Keep feature gates aligned with `pkg/mcp/features.go`.
- Register resources near the existing MCP server wiring.
- Put tests next to the behavior they validate.
- Update docs when the runtime-visible resource surface changes.
- Do not expose tokens, auth params, key files, TLS private keys, or secrets.
- Do not consume Pulsar messages, commit cursors, clear backlog, unload topics,
  split bundles, delete resources, start workloads, or stop workloads from a
  resource handler.
- Do not commit. The runner owns commits.

## Expected implementation shape

- Static resources should describe stable, small top-level context.
- Dynamic resource templates should cover large or parameterized collections.
- URI parsing must reject malformed or unsupported resource URIs.
- Missing Pulsar session should return a clear error.
- JSON resource contents should use `application/json`.
- Text documentation resources, if added, should use `text/markdown`.

## Verification

Run focused tests that prove the selected item is complete. If you cannot run a
required test, explain the blocker in `acceptance_notes` and use `blocked` or
`failed` status instead of `implemented`.

## Final response

Return only JSON matching `ralph/pulsar-resources/result_schema.json`.
