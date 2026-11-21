---
slug: 20260901-1100-implement-batch-progress-clarity
lens: Backend-heavy
---

# Implementation Roadmap: Implement Batch Progress Clarity

## 1. Execution Strategy
- **Gravity Order:** Unit Test Assertions (RED) -> Blueprint & Manifest Updates (GREEN) -> Module Documentation Alignment

## 2. Tasks

### Phase 1: Test Assertions & Command Blueprint Update

- [x] T1.1: [TEST] Add Unit Test Assertions for Clarified Progress Metrics Format
**Target:** `src/internal/agent/kit_manifests_test.go`
**Context:** [US-KIT-07]

**Action Steps:**
- Update `TestImplementBatchSubagentIsolationAndProgress` in `src/internal/agent/kit_manifests_test.go`.
- Replace legacy assertion strings with mandatory format tokens: `DELEGATING BATCH`, `Done:`, `Batch:`, `Remaining:`, `completed_tasks_count`, `total_tasks_count`, `BATCH`, and `COMPLETED`.
- Execute `go test ./src/internal/agent/...` to verify test failure against existing unupdated blueprint (`implement.yaml`).

**Acceptance Check:**
`go test -v -run TestImplementBatchSubagentIsolationAndProgress ./src/internal/agent/...` fails due to missing new progress format tokens.

- [x] T1.2: [CODE] Update implement.yaml with Option A Explicit Progress Metrics
**Target:** `src/internal/agent/kit/commands/implement.yaml`
**Context:** [US-KIT-07]

**Action Steps:**
- Update the initialization metrics guidance in Step 1 to define `total_tasks_count` and `completed_tasks_count`.
- Update the Mission Brief Envelope progress line to: `Progress: {completed_tasks_count}/{total_tasks_count} completed ({progress_pct}%) | Executing: {batch_task_count} tasks [{task_ids_list}] | Remaining: {remaining_tasks_count}`.
- Update `MANDATORY LOG` directive to: `[DELEGATING BATCH {batch_index}/{total_batches}] [Done: {completed_tasks_count}/{total_tasks_count} ({progress_pct}%) | Batch: {batch_task_count} tasks ({task_ids_list}) | Remaining: {remaining_tasks_count}] -> Persona: {required_persona_role} | Model Tier: {recommended_model_tier} | Effort: {recommended_effort} | Tasks: {task_ids_list}`.
- Update `Log Progress` directive on success to: `[BATCH {batch_index}/{total_batches} COMPLETED] [Done: {completed_tasks_count}/{total_tasks_count} ({progress_pct}%) | Remaining: {remaining_tasks_count}] -> All batch tasks verified`.

**Acceptance Check:**
`go test -v -run TestImplementBatchSubagentIsolationAndProgress ./src/internal/agent/...` passes with exit code 0.

### Phase 2: Living Module Documentation Alignment

- [x] T2.1: [DOCS] Update Agent Kit Living Specification Invariants and User Stories
**Target:** `.specforce/docs/modules/agent-kit.md`
**Context:** [US-KIT-07]

**Action Steps:**
- Update rule `[BR-KIT-09]` in `.specforce/docs/modules/agent-kit.md` to reference the new standardized real-time progress metrics format (`[DELEGATING BATCH X/Y] [Done: C/Z (P%) | Batch: N tasks (IDs) | Remaining: R]` and `[BATCH X/Y COMPLETED] [Done: C/Z (P%) | Remaining: R]`).
- Update user story `[US-KIT-07]` scenario THEN-clause to reflect the explicit quantity counters (`Done`, `Batch`, `Remaining`).
- Run `git status --porcelain` to verify clean changes and consistency across documentation and implementation.

**Acceptance Check:**
`grep -E "Done:.*Batch:.*Remaining:" .specforce/docs/modules/agent-kit.md` returns match.
