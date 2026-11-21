# Proposal: Persist QA Report and Hooks for RED Tasks

**Suggested type/size:** feature / large (touches Go CLI + blueprints). Run `/spf:discovery` first if the design is still open.
**Origin:** follow-up 3 of archived spec `.specforce/archive/20261020-2222-independent-qa-and-free-batch-log` (commit 9c6eab8 on `1.x`).

## Problem
1. **QA evidence is ephemeral.** The independent qa-tester returns a traceability matrix (criterion, scenario, evidence, verdict PASSED/FAILED/NOT-VERIFIABLE) and out-of-scope findings only in the conversation. `/spf:archive` derives next steps from "findings present in the conversation"; if archival runs in a new session the findings are gone and the fallback wrongly states that none were reported.
2. **The `on_task_finished` hook blocks RED tasks.** `.specforce/config.yaml` runs `make lint` and `make test` when finishing tasks. In a TDD pair, task T1.1 (RED) adds failing tests whose passing depends on a later batch's files. `specforce implementation update <slug> --task T1.1 --status finished` fails with "task update blocked by hook failures" until the later batch is done, so tasks had to be finished all together at the end. Observed in practice during the previous implementation.
3. Environment note (separate from the product): `make lint` failed because `/snap/bin/golangci-lint` shadows `~/go/bin/golangci-lint` (snap-confine permissions); the workaround was `PATH="$HOME/go/bin:$PATH"` on the command. Do not change config to hide it.

## Findings (verified)
- Spec metadata lives in `spec.yaml` per spec; machine-readable state is meant to go there (engineering.md "Machine-Readable Metadata Standard"); Markdown artifacts stay human-readable.
- Hooks run once per update command, deduplicated, in parallel (`[BR-KIT-06]`, architecture principle 8, `core.HookResult`).
- Archival instructions come from the CLI (`specforce archive instructions`, template in `src/internal/agent/kit/instructions/archive.md`).

## Ideas to evaluate (not decided)
- QA report persistence: a `qa-report.md` (or a section in `spec.yaml`) written by the orchestrator (never by the tester, which is read-only) at the end of validation; `specforce archive instructions` and the archive skill read it.
- Hook semantics for RED tasks: a `[RED]` tag skips test hooks, or a per-task `--allow-red` flag, or deferring hooks to phase-level (`on_phase_finished`) for RED/GREEN pairs. Preserve the rule that a task is `finished` only when verification passes.
- Whether `specforce implementation status --json` should expose the QA report path/state.

## Decisions already taken
- The qa-tester stays read-only on production behavior; only the orchestrator persists the report.
- No weakening of verification: any hook relaxation must be explicit and scoped, never a silent skip.

## Open questions
1. Report as a new artifact (affects the artifact registry/topological order and `spec status`) or as data in `spec.yaml`?
2. Does the RED relaxation belong in config (`.specforce/config.yaml`) or in the task model (a state/tag)?
3. Should archive refuse to run when QA status is FAILED?

## Acceptance hints
After implementation, a new session running `/spf:archive` shows next steps citing persisted findings; finishing a RED task whose tests legitimately fail no longer blocks, while finishing a GREEN task still runs lint+tests; unit tests cover both.

## How to use this proposal in a new conversation
Start with `/spf:discovery` ("explore `.specforce/specs/20261023-1159-persist-qa-report-and-red-task-hooks/proposal.md`"), settle the open questions, then `/spf:spec`. The implementation will exercise the full new flow including the qa-tester and correction cycle.
