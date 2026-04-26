# Ralph Loop For Pulsar MCP Resources

This directory defines a repeatable Ralph Loop for completing Pulsar MCP
resource support in small, reviewable iterations. Each successful iteration
must implement one backlog item or fix one resource defect, verify it, update
state, and create exactly one git commit.

## Files

- `backlog.yaml`: ordered implementation backlog for Pulsar resources.
- `state.json`: tracked loop state. It records completed items and can carry
  manually blocked backlog items, but intentionally does not store commit
  hashes.
- `prompts/`: base contract and lane-specific prompts used by the runner.
- `result_schema.json`: JSON schema for the agent's final response.

Runtime reports are written under `tmp/ralph-pulsar-resources/`, which is
ignored by git.

## Usage

Preview the next two iterations without changing tracked files:

```bash
./scripts/ralph-pulsar-resources-loop.sh --dry-run --max-iterations 2
```

Run one real iteration with the default Codex backend:

```bash
./scripts/ralph-pulsar-resources-loop.sh --max-iterations 1
```

Run only a specific lane:

```bash
./scripts/ralph-pulsar-resources-loop.sh --lane protocol --max-iterations 1
```

Use a different Codex binary or wrapper:

```bash
./scripts/ralph-pulsar-resources-loop.sh --agent-cmd "codex" --model gpt-5.4 --max-iterations 1
```

## Iteration Contract

The runner owns commits. The agent must not commit directly. A successful
iteration has this shape:

1. Select one pending backlog item.
2. Generate a fresh prompt from the base contract, lane prompt, current state,
   and selected backlog item.
3. Invoke `codex exec` with the result schema.
4. Require a structured `implemented` result with non-empty `tests_run`.
5. Run the backlog item's focused verification commands.
6. Update `state.json`.
7. Commit all tracked changes with the agent-provided commit subject.

If the agent reports `blocked`, `failed`, or `no_op`, the runner stops without
advancing state or committing. If tracked files changed in that path, the runner
also stops and leaves the worktree for manual inspection.

## Completion Definition

Pulsar MCP resource support is complete when the backlog covers both protocol
surface and Pulsar resource families:

- MCP `resources/list`, `resources/templates/list`, and `resources/read` work
  for Pulsar resources.
- Resource registration honors existing feature gates and current Pulsar
  session context.
- Resource handlers are read-only and do not consume messages, mutate cursors,
  or perform admin write operations.
- Large collections are exposed through templates and lazy reads, not static
  full-cluster dumps.
- Tests cover registration, URI parsing, missing-session behavior, feature-gate
  behavior, and representative handler behavior for each resource family.
- Docs describe URI shapes, safety boundaries, and intended usage.
