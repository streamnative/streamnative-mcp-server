#!/usr/bin/env python3
"""Run the Ralph Loop for Pulsar MCP resources."""

from __future__ import annotations

import argparse
import datetime as dt
import json
import shlex
import subprocess
import sys
from pathlib import Path
from typing import Any


VALID_STATUSES = {"implemented", "blocked", "failed", "no_op"}
STATE_REL_PATH = "ralph/pulsar-resources/state.json"


class LoopError(RuntimeError):
    """Raised for expected loop failures."""


def now_utc() -> str:
    return dt.datetime.now(dt.timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z")


def run_text(args: list[str], cwd: Path, check: bool = True) -> str:
    proc = subprocess.run(args, cwd=cwd, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    if check and proc.returncode != 0:
        command = " ".join(shlex.quote(part) for part in args)
        raise LoopError(f"command failed: {command}\n{proc.stderr.strip()}")
    return proc.stdout.strip()


def git(repo: Path, *args: str, check: bool = True) -> str:
    return run_text(["git", *args], repo, check=check)


def resolve_repo(path: str | None) -> Path:
    start = Path(path).resolve() if path else Path.cwd().resolve()
    root = run_text(["git", "rev-parse", "--show-toplevel"], start)
    return Path(root).resolve()


def write_json(path: Path, data: Any) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(data, indent=2, sort_keys=False) + "\n", encoding="utf-8")


def read_json(path: Path) -> Any:
    return json.loads(path.read_text(encoding="utf-8"))


def parse_scalar(value: str) -> Any:
    value = value.strip()
    if value == "":
        return ""
    if value.startswith('"') and value.endswith('"'):
        return value[1:-1]
    if value.startswith("'") and value.endswith("'"):
        return value[1:-1]
    if value == "true":
        return True
    if value == "false":
        return False
    if value.startswith("[") and value.endswith("]"):
        inner = value[1:-1].strip()
        if not inner:
            return []
        return [parse_scalar(part.strip()) for part in inner.split(",")]
    if value.isdigit():
        return int(value)
    return value


def load_backlog(path: Path) -> dict[str, Any]:
    items: list[dict[str, Any]] = []
    current: dict[str, Any] | None = None
    in_items = False

    for raw_line in path.read_text(encoding="utf-8").splitlines():
        line = raw_line.strip()
        if not line or line.startswith("#"):
            continue
        if line == "items:":
            in_items = True
            continue
        if not in_items:
            continue
        if line.startswith("- "):
            if current is not None:
                items.append(current)
            current = {}
            remainder = line[2:].strip()
            if remainder:
                key, value = split_key_value(remainder)
                current[key] = parse_scalar(value)
            continue
        if current is None:
            continue
        key, value = split_key_value(line)
        current[key] = parse_scalar(value)

    if current is not None:
        items.append(current)
    if not items:
        raise LoopError(f"no backlog items found in {path}")
    return {"items": items}


def split_key_value(line: str) -> tuple[str, str]:
    if ":" not in line:
        raise LoopError(f"invalid backlog line: {line}")
    key, value = line.split(":", 1)
    return key.strip(), value.strip()


def load_state(path: Path) -> dict[str, Any]:
    if not path.exists():
        raise LoopError(f"state file does not exist: {path}")
    state = read_json(path)
    for key in ["completed_items", "blocked_items", "lane_history"]:
        state.setdefault(key, [])
    state.setdefault("current_iteration", 0)
    return state


def select_item(
    backlog: dict[str, Any],
    state: dict[str, Any],
    lane: str,
    planned_ids: set[str],
) -> dict[str, Any] | None:
    completed = set(state.get("completed_items", [])) | planned_ids
    blocked = set(state.get("blocked_items", []))
    candidates = []
    for item in backlog["items"]:
        item_id = str(item.get("id", ""))
        item_lane = str(item.get("lane", ""))
        item_status = str(item.get("status", "pending"))
        if not item_id or item_id in completed or item_id in blocked:
            continue
        if item_status not in {"pending", "ready"}:
            continue
        if lane != "auto" and item_lane != lane:
            continue
        candidates.append(item)
    candidates.sort(key=lambda item: int(item.get("priority", 1000)))
    return candidates[0] if candidates else None


def read_prompt(path: Path) -> str:
    if not path.exists():
        raise LoopError(f"prompt file does not exist: {path}")
    return path.read_text(encoding="utf-8").strip()


def build_prompt(
    repo: Path,
    prompt_dir: Path,
    state: dict[str, Any],
    item: dict[str, Any],
    iteration_number: int,
    baseline: dict[str, Any],
) -> str:
    lane = str(item.get("lane", "family"))
    lane_prompt = prompt_dir / f"lane_{lane}.md"
    if not lane_prompt.exists():
        lane_prompt = prompt_dir / "lane_family.md"

    base = read_prompt(prompt_dir / "base_contract.md")
    lane_text = read_prompt(lane_prompt)

    context = {
        "iteration": iteration_number,
        "repo": str(repo),
        "selected_item": item,
        "state": state,
        "baseline": baseline,
    }

    return (
        f"{base}\n\n"
        f"{lane_text}\n\n"
        "## Iteration Context\n\n"
        "Use this JSON as the source of truth for the selected work item:\n\n"
        "```json\n"
        f"{json.dumps(context, indent=2, sort_keys=False)}\n"
        "```\n\n"
        "Complete only `selected_item`. Return JSON only.\n"
    )


def unique_run_dir(report_dir: Path) -> Path:
    timestamp = dt.datetime.now(dt.timezone.utc).strftime("%Y%m%dT%H%M%SZ")
    base = report_dir / timestamp
    candidate = base
    suffix = 1
    while candidate.exists():
        suffix += 1
        candidate = Path(f"{base}-{suffix}")
    candidate.mkdir(parents=True, exist_ok=False)
    return candidate


def codex_command(args: argparse.Namespace, repo: Path, schema_file: Path, result_file: Path) -> list[str]:
    command = shlex.split(args.agent_cmd)
    if not command:
        raise LoopError("--agent-cmd must not be empty")
    command.extend(
        [
            "exec",
            "--cd",
            str(repo),
            "--output-schema",
            str(schema_file),
            "--output-last-message",
            str(result_file),
            "--json",
        ]
    )
    if args.model:
        command.extend(["--model", args.model])
    if args.dangerously_bypass:
        command.append("--dangerously-bypass-approvals-and-sandbox")
    else:
        command.extend(["--sandbox", args.sandbox])
    command.append("-")
    return command


def run_agent(
    args: argparse.Namespace,
    repo: Path,
    schema_file: Path,
    prompt: str,
    iteration_dir: Path,
) -> dict[str, Any]:
    result_file = iteration_dir / "agent_result.json"
    stdout_file = iteration_dir / "agent_events.jsonl"
    stderr_file = iteration_dir / "agent_stderr.log"
    command = codex_command(args, repo, schema_file, result_file)
    write_json(iteration_dir / "agent_command.json", {"command": command})

    with stdout_file.open("w", encoding="utf-8") as stdout, stderr_file.open("w", encoding="utf-8") as stderr:
        proc = subprocess.run(command, cwd=repo, input=prompt, text=True, stdout=stdout, stderr=stderr)

    if proc.returncode != 0:
        raise LoopError(f"agent command failed with exit code {proc.returncode}; see {stderr_file}")
    if not result_file.exists():
        raise LoopError(f"agent did not write result file: {result_file}")
    return parse_agent_result(result_file)


def parse_agent_result(path: Path) -> dict[str, Any]:
    text = path.read_text(encoding="utf-8").strip()
    try:
        value = json.loads(text)
    except json.JSONDecodeError:
        start = text.find("{")
        end = text.rfind("}")
        if start < 0 or end <= start:
            raise LoopError(f"agent result is not JSON: {path}") from None
        value = json.loads(text[start : end + 1])
    if not isinstance(value, dict):
        raise LoopError("agent result must be a JSON object")
    return value


def validate_result(result: dict[str, Any], item: dict[str, Any]) -> None:
    status = result.get("status")
    if status not in VALID_STATUSES:
        raise LoopError(f"invalid result status: {status}")
    if result.get("item_id") != item.get("id"):
        raise LoopError(f"result item_id {result.get('item_id')} does not match selected item {item.get('id')}")
    if result.get("lane") != item.get("lane"):
        raise LoopError(f"result lane {result.get('lane')} does not match selected lane {item.get('lane')}")
    if status == "implemented" and not result.get("tests_run"):
        raise LoopError("implemented result must include at least one test command")
    subject = str(result.get("commit_subject", "")).strip()
    if "\n" in subject or not subject:
        raise LoopError("commit_subject must be a non-empty single line")
    result["commit_subject"] = subject


def run_verification(repo: Path, item: dict[str, Any], iteration_dir: Path) -> list[str]:
    commands = item.get("focused_tests", [])
    if not isinstance(commands, list) or not commands:
        raise LoopError(f"backlog item {item.get('id')} has no focused_tests")

    log_path = iteration_dir / "runner_verification.log"
    completed: list[str] = []
    with log_path.open("w", encoding="utf-8") as log:
        for command in commands:
            command_text = str(command)
            log.write(f"$ {command_text}\n")
            log.flush()
            proc = subprocess.run(command_text, cwd=repo, text=True, shell=True, stdout=log, stderr=subprocess.STDOUT)
            if proc.returncode != 0:
                raise LoopError(f"verification failed: {command_text}; see {log_path}")
            completed.append(command_text)
            log.write("\n")
    return completed


def update_state(
    state_path: Path,
    state: dict[str, Any],
    item: dict[str, Any],
    result: dict[str, Any],
    verification_commands: list[str],
) -> dict[str, Any]:
    item_id = str(item["id"])
    state["current_iteration"] = int(state.get("current_iteration", 0)) + 1
    state["updated_at"] = now_utc()
    completed = state.setdefault("completed_items", [])
    if item_id not in completed:
        completed.append(item_id)
    if item_id in state.get("blocked_items", []):
        state["blocked_items"] = [blocked for blocked in state["blocked_items"] if blocked != item_id]

    entry = {
        "iteration": state["current_iteration"],
        "item_id": item_id,
        "lane": item["lane"],
        "status": result["status"],
        "summary": result["summary"],
        "tests_run": result["tests_run"],
        "runner_verification": verification_commands,
        "completed_at": state["updated_at"],
    }
    state.setdefault("lane_history", []).append(entry)
    state["last_result"] = entry
    write_json(state_path, state)
    return state


def ensure_implemented_changes(before_status: str, after_status: str, result: dict[str, Any]) -> None:
    if before_status == after_status and not before_status:
        raise LoopError("agent reported implemented but did not change tracked or untracked files")
    changed = [str(path) for path in result.get("files_changed", [])]
    if not changed:
        raise LoopError("implemented result must list files_changed")


def ensure_no_changes_for_non_success(before_status: str, after_status: str) -> None:
    if before_status != after_status:
        raise LoopError("agent returned non-success status but changed the worktree; inspect before continuing")


def create_commit(repo: Path, subject: str) -> str:
    git(repo, "add", "-A")
    cached_paths = git(repo, "diff", "--cached", "--name-only")
    if not cached_paths.strip():
        raise LoopError("no staged changes to commit")
    non_state_paths = [path for path in cached_paths.splitlines() if path != STATE_REL_PATH]
    if not non_state_paths:
        raise LoopError("refusing to commit state-only iteration")
    git(repo, "commit", "-m", subject)
    return git(repo, "rev-parse", "HEAD")


def write_iteration_report(iteration_dir: Path, report: dict[str, Any]) -> None:
    write_json(iteration_dir / "round_result.json", report)


def run_real_iteration(
    args: argparse.Namespace,
    repo: Path,
    item: dict[str, Any],
    state: dict[str, Any],
    state_path: Path,
    prompt_dir: Path,
    schema_file: Path,
    iteration_dir: Path,
    baseline: dict[str, Any],
) -> dict[str, Any]:
    iteration_number = int(state.get("current_iteration", 0)) + 1
    prompt = build_prompt(repo, prompt_dir, state, item, iteration_number, baseline)
    (iteration_dir / "prompt.md").write_text(prompt, encoding="utf-8")
    write_json(iteration_dir / "selected_item.json", item)

    before_head = git(repo, "rev-parse", "HEAD")
    before_status = git(repo, "status", "--porcelain")
    result = run_agent(args, repo, schema_file, prompt, iteration_dir)
    validate_result(result, item)
    after_agent_status = git(repo, "status", "--porcelain")

    if result["status"] != "implemented":
        ensure_no_changes_for_non_success(before_status, after_agent_status)
        report = {
            "status": result["status"],
            "item_id": item["id"],
            "before_head": before_head,
            "after_head": git(repo, "rev-parse", "HEAD"),
            "result": result,
            "finished_at": now_utc(),
        }
        write_iteration_report(iteration_dir, report)
        return report

    ensure_implemented_changes(before_status, after_agent_status, result)
    verification_commands = run_verification(repo, item, iteration_dir)
    updated_state = update_state(state_path, state, item, result, verification_commands)
    commit_hash = create_commit(repo, result["commit_subject"])
    report = {
        "status": "committed",
        "item_id": item["id"],
        "lane": item["lane"],
        "before_head": before_head,
        "after_head": commit_hash,
        "result": result,
        "state": updated_state,
        "finished_at": now_utc(),
    }
    write_iteration_report(iteration_dir, report)
    return report


def parse_args(argv: list[str]) -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Run the Pulsar resources Ralph Loop.")
    parser.add_argument("--repo", default=None, help="Repository root. Defaults to current git root.")
    parser.add_argument("--max-iterations", type=int, default=1, help="Maximum iterations to run.")
    parser.add_argument("--lane", default="auto", choices=["auto", "protocol", "family", "fix", "docs"])
    parser.add_argument("--agent-cmd", default="codex", help="Codex-compatible command name or path.")
    parser.add_argument("--model", default="", help="Optional model passed to codex exec.")
    parser.add_argument("--sandbox", default="danger-full-access", choices=["read-only", "workspace-write", "danger-full-access"])
    parser.add_argument("--dangerously-bypass", action="store_true", help="Pass Codex's bypass approvals and sandbox flag.")
    parser.add_argument("--dry-run", action="store_true", help="Write prompts and plan without invoking the agent.")
    parser.add_argument("--require-clean", action="store_true", help="Fail if tracked files are dirty at startup.")
    parser.add_argument("--state-file", default=STATE_REL_PATH)
    parser.add_argument("--backlog-file", default="ralph/pulsar-resources/backlog.yaml")
    parser.add_argument("--prompt-dir", default="ralph/pulsar-resources/prompts")
    parser.add_argument("--schema-file", default="ralph/pulsar-resources/result_schema.json")
    parser.add_argument("--report-dir", default="tmp/ralph-pulsar-resources")
    args = parser.parse_args(argv)
    if args.max_iterations < 1:
        raise LoopError("--max-iterations must be at least 1")
    return args


def main(argv: list[str]) -> int:
    try:
        args = parse_args(argv)
        repo = resolve_repo(args.repo)
        state_path = (repo / args.state_file).resolve()
        backlog_path = (repo / args.backlog_file).resolve()
        prompt_dir = (repo / args.prompt_dir).resolve()
        schema_file = (repo / args.schema_file).resolve()
        report_dir = (repo / args.report_dir).resolve()

        tracked_status = git(repo, "status", "--porcelain", "--untracked-files=no")
        if args.require_clean and tracked_status.strip():
            raise LoopError("tracked worktree is dirty; rerun without --require-clean to record baseline")

        baseline = {
            "branch": git(repo, "branch", "--show-current"),
            "head": git(repo, "rev-parse", "HEAD"),
            "status": git(repo, "status", "--short"),
            "tracked_status": tracked_status,
            "started_at": now_utc(),
        }

        run_dir = unique_run_dir(report_dir)
        summary: dict[str, Any] = {
            "status": "running",
            "dry_run": args.dry_run,
            "repo": str(repo),
            "baseline": baseline,
            "iterations": [],
        }
        write_json(run_dir / "run_summary.json", summary)

        planned_ids: set[str] = set()
        for dry_index in range(1, args.max_iterations + 1):
            state = load_state(state_path)
            backlog = load_backlog(backlog_path)
            item = select_item(backlog, state, args.lane, planned_ids)
            if item is None:
                summary["status"] = "complete"
                summary["message"] = "no pending backlog item matched the selection"
                break

            iteration_number = dry_index if args.dry_run else int(state.get("current_iteration", 0)) + 1
            iteration_dir = run_dir / f"iteration-{iteration_number:03d}-{item['id']}"
            iteration_dir.mkdir(parents=True, exist_ok=False)

            if args.dry_run:
                prompt = build_prompt(repo, prompt_dir, state, item, iteration_number, baseline)
                (iteration_dir / "prompt.md").write_text(prompt, encoding="utf-8")
                write_json(iteration_dir / "selected_item.json", item)
                planned_ids.add(str(item["id"]))
                summary["iterations"].append(
                    {
                        "status": "planned",
                        "iteration": iteration_number,
                        "item_id": item["id"],
                        "lane": item["lane"],
                        "prompt": str(iteration_dir / "prompt.md"),
                    }
                )
                write_json(run_dir / "run_summary.json", summary)
                continue

            report = run_real_iteration(
                args,
                repo,
                item,
                state,
                state_path,
                prompt_dir,
                schema_file,
                iteration_dir,
                baseline,
            )
            summary["iterations"].append(report)
            write_json(run_dir / "run_summary.json", summary)
            if report["status"] != "committed":
                summary["status"] = report["status"]
                break
        else:
            summary["status"] = "planned" if args.dry_run else "ok"

        summary["finished_at"] = now_utc()
        write_json(run_dir / "run_summary.json", summary)
        print(f"Ralph Loop report: {run_dir}")
        return 0 if summary["status"] in {"ok", "planned", "complete"} else 1
    except LoopError as exc:
        print(f"error: {exc}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
