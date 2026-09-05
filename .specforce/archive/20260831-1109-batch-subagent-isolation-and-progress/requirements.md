---
slug: 20260831-1109-batch-subagent-isolation-and-progress
lens: Integration
---

# Feature: Isolate Batch Subagents and Implementation Progress Tracking

## 1. Context & Value
During implementation execution, re-using existing subagent worker sessions across different task batches causes context contamination and token exhaustion. By enforcing strict worker isolation (one clean subagent session per task batch) with targeted within-batch feedback loops, and providing real-time progress indicators (Tasks X of Y, Batch M of N), developers gain predictable execution, token efficiency, and clear visibility into implementation progress.

## 2. Out of Scope (Anti-Goals)
- Modifying CLI commands or underlying Go execution engines for task parsing (`specforce implementation status` / `specforce implementation update`).
- Changing the blueprint format or protocol for discovery, planning, or archival workflows.
- Auto-canceling background tasks outside of the standard subagent lifecycle.

## 3. Acceptance Criteria (BDD)

### [US-1] Fresh Subagent Spawning per Task Batch
**User Story:** AS A developer using the implementation orchestrator, I WANT each task batch to be delegated to a freshly initialized subagent worker session, SO THAT previous task context does not contaminate subsequent batches.

**Scenarios:**
1. **[Happy Path - Batch Delegation]** GIVEN a roadmap with multiple sequential or parallel task batches WHEN the orchestrator moves to the next batch THEN it spawns a brand-new subagent worker session with clean context instead of reusing the previous batch's session.
2. **[Edge Case - Multi-Task Batch]** GIVEN a batch containing multiple grouped tasks WHEN delegation occurs THEN all tasks within that batch are handled by that newly spawned worker session.

**UI/UX Specifics:**
- **View/Component:** Orchestrator CLI console output.
- **Feedback Logic:** Clear delegation header showing fresh worker initialization.
- **Keybindings:** Non-interactive console stream.

**Technical Constraints (NFR):**
- **[Performance]:** Zero context accumulation across different batch boundaries.
- **[Integrity]:** Subagent workers must not access conversation memory from unrelated prior batches.
- **[Observability]:** Explicit delegation log emitted on every batch dispatch.

### [US-2] Within-Batch Feedback and Self-Healing Loop
**User Story:** AS A developer, I WANT verification failures to be sent back exclusively to the specific subagent that implemented the failing batch, SO THAT fixes are applied with relevant local context while preventing leaks to other batches.

**Scenarios:**
1. **[Happy Path - Test Failure Correction]** GIVEN a subagent has completed a batch but terminal verification fails WHEN the orchestrator initiates a self-healing iteration THEN it sends error feedback and correction requests to the exact same subagent session that implemented that batch (up to 2 iterations).
2. **[Edge Case - Persistent Failure Escalation]** GIVEN a subagent fails verification after the maximum allowed iterations WHEN resolution fails THEN the orchestrator marks the batch tasks as failed and halts execution to ask the developer for assistance without invoking new workers.

**UI/UX Specifics:**
- **View/Component:** Orchestrator feedback log.
- **Feedback Logic:** Explicit iteration counter on self-healing retry.
- **Keybindings:** Standard CLI flow.

**Technical Constraints (NFR):**
- **[Performance]:** Bounded self-healing retry limit (max 2 iterations) before human escalation.
- **[Integrity]:** Feedback messages restricted strictly to the authoring subagent session of the active batch.
- **[Observability]:** Detailed error message and fix instructions passed in feedback envelope.

### [US-3] Real-Time Progress Metrics and Execution Visibility
**User Story:** AS A developer monitoring feature implementation, I WANT to see explicit progress metrics (current batch out of total batches, current tasks out of total tasks, and completion percentage), SO THAT I can accurately track implementation velocity and remaining effort.

**Scenarios:**
1. **[Happy Path - Delegation Progress]** GIVEN a roadmap with mapped batches and total tasks WHEN a batch is delegated THEN the orchestrator outputs a mandatory log with batch index, total batches, task index range, total tasks, and progress percentage.
2. **[Happy Path - Batch Completion Progress]** GIVEN a batch successfully passes verification WHEN status updates are completed THEN the orchestrator outputs a completion log showing updated completed tasks count, total tasks, percentage completed, and the upcoming batch indicator.
3. **[Edge Case - Single Batch Roadmap]** GIVEN a roadmap containing only 1 batch WHEN execution begins and finishes THEN the progress indicator accurately displays 100% upon batch resolution.

**UI/UX Specifics:**
- **View/Component:** Orchestrator status stream and Mission Brief header.
- **Feedback Logic:** Structured bracket format `[DELEGATING BATCH X/Y] [Tasks A..B of Z | P%]` and `[BATCH X/Y COMPLETED] Progress: [C/Z tasks finished - P%]`.
- **Keybindings:** Streamed output.

**Technical Constraints (NFR):**
- **[Performance]:** O(1) roadmap progress metric computation during session start.
- **[Integrity]:** Consistent task and batch counters across mission brief envelopes and console logs.
- **[Observability]:** Standardized structured log format across all supported agent harnesses.

## 4. Business Invariants
- A subagent session spawned for Batch N MUST NEVER receive task assignments for Batch N+1.
- Message-based feedback loops are strictly confined to the active batch's authoring subagent during verification failure adjustments.
- Total tasks and total batches metrics must remain mathematically consistent from roadmap initialization to completion.

## 5. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Immediate prompt evaluation with zero redundant CLI queries.
- **[Reliability]:** Fail-fast behavior on unresolvable batch verification failures.
- **[Maintainability]:** Clean synchronization between `implement.yaml` blueprint and `.agents/skills/spf-implement/SKILL.md`.
