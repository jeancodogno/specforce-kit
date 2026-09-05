---
slug: 20260831-1109-batch-subagent-isolation-and-progress
lens: Integration
---

# Implementation Roadmap: Isolate Batch Subagents and Implementation Progress Tracking

## 1. Execution Strategy
- **Phase 1 (RED & GREEN):** Write test assertions in `kit_manifests_test.go` demanding batch subagent isolation and standardized progress metrics (RED), then update `src/internal/agent/kit/commands/implement.yaml` with clean worker spawning per batch, within-batch feedback loops, and explicit progress formats (GREEN).
- **Phase 2 (Skill Synchronization & Integration Verification):** Synchronize `.agents/skills/spf-implement/SKILL.md` with the updated batch orchestration protocol and run full test suites across the package to ensure blueprint consistency.

## 2. Tasks

### Phase 1: Implement Blueprint & Test Verification

- [x] T1.1: [TEST] [RED] Add unit test assertions for subagent batch isolation and progress metrics tracking
**Target:** `src/internal/agent/kit_manifests_test.go`
**Context:** [US-1], [US-2], [US-3]

**Action Steps:**
- Add `TestImplementBatchSubagentIsolationAndProgress` in `kit_manifests_test.go`.
- Define assertions validating the presence of fresh subagent spawning per batch directives in `implement.yaml`.
- Define assertions validating within-batch feedback loop constraints and structured progress metrics `[DELEGATING BATCH {batch_index}/{total_batches}] [Tasks {task_start}..{task_end} of {total_tasks} | {progress_pct}%]` and `[BATCH {batch_index}/{total_batches} COMPLETED] Progress: [{completed_tasks}/{total_tasks} tasks finished - {progress_pct}%]`.

**Acceptance Check:**
Run `go test -v ./src/internal/agent/ -run TestImplementBatchSubagentIsolationAndProgress` (should fail RED prior to implement.yaml changes).

- [x] T1.2: [CODE] [GREEN] Update implement.yaml with fresh subagent spawning, within-batch feedback, and progress metrics
**Target:** `src/internal/agent/kit/commands/implement.yaml`
**Context:** [US-1], [US-2], [US-3]

**Action Steps:**
- Update Step 1 (Roadmap Mapping) in `implement.yaml` to calculate total tasks count ($N$) and total batches count ($B$).
- Replace the ambiguous worker reuse directive in Section 2 with explicit fresh subagent worker session spawning per batch and strictly confine message reuse to the within-batch verification failure self-healing loop.
- Update the Mission Brief envelope and `MANDATORY LOG` to include `[DELEGATING BATCH {batch_index}/{total_batches}] [Tasks {task_start}..{task_end} of {total_tasks} | {progress_pct}%]` and completion log `[BATCH {batch_index}/{total_batches} COMPLETED] Progress: [{completed_tasks}/{total_tasks} tasks finished - {progress_pct}%]`.

**Acceptance Check:**
Run `go test -v ./src/internal/agent/ -run TestImplementBatchSubagentIsolationAndProgress` (must pass GREEN with exit code 0).

### Phase 2: Skill Synchronization & Integration Verification

- [x] T2.1: [CODE] Synchronize spf-implement skill definition
**Target:** `.agents/skills/spf-implement/SKILL.md`
**Context:** [US-1], [US-2], [US-3]

**Action Steps:**
- Update `.agents/skills/spf-implement/SKILL.md` to reflect the exact subagent isolation rules from `implement.yaml`.
- Update the Mission Brief template and `MANDATORY LOG` formats to mirror the progress counters and batch tracking directives.
- Verify that no conflicting instructions remain regarding cross-batch worker reuse.

**Acceptance Check:**
Run `grep "Fresh Subagent per Batch" .agents/skills/spf-implement/SKILL.md` and verify clean matching.

- [x] T2.2: [VERIFY] Run package test suite and specforce manifest validation
**Target:** `src/internal/agent/...`
**Context:** [US-1], [US-2], [US-3]

**Action Steps:**
- Execute `go test ./src/internal/agent/...` to verify all kit manifest and translator tests pass.
- Verify YAML structural integrity with `git diff src/internal/agent/kit/commands/implement.yaml`.
- Ensure all test suites pass without regression.

**Acceptance Check:**
Run `go test -v ./src/internal/agent/...` (exit code 0).
