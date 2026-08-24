---
slug: 20260823-2158-batch-task-update-and-subagent-orchestration
lens: Integration
---

# Feature: Batch Task Update & Subagent Orchestration

## 1. Context & Value
Developers and AI orchestrators currently need to execute status update commands one task at a time, creating unnecessary tool invocation overhead and redundant test verification runs. Enabling multi-task updates and standardizing subagent delegation boundaries ensures faster execution cycles, lower operational cost, and strict task boundary safety.

## 2. Out of Scope (Anti-Goals)
- Modifying the markdown syntax of how tasks are formatted inside specifications.
- Altering the archival lifecycle or constitution generation workflows.
- Permitting subagents or workers to mutate specification files directly.

## 3. Acceptance Criteria (BDD)

### [US-1] Multi-Task Status Update via CLI
**User Story:** AS A software engineer or automated orchestrator, I WANT TO update multiple task statuses simultaneously in a single command, SO THAT I reduce execution latency and avoid redundant verification cycles.

**Scenarios:**
1. **[Happy Path - Multi-task Transition]** GIVEN a specification with tasks in pending state WHEN the user updates multiple task identifiers to in-progress or finished in a single command THEN all referenced tasks transition to the requested state atomically and their session time tracking updates.
2. **[Edge Case - Invalid Task Identifier]** GIVEN a specification with tasks WHEN the user submits a batch update containing at least one nonexistent task identifier THEN the operation fails fast with an explicit error and no task states are modified.
3. **[Edge Case - Hook Failure]** GIVEN a specification configured with automated verification hooks WHEN a batch update to finished status encounters a verification hook failure THEN the operation aborts immediately and task states remain unchanged.

**UI/UX Specifics:**
- **View/Component:** Command-line standard output and status reporter.
- **Feedback Logic:** Clear confirmation displaying all updated task identifiers on success, or structured failure summaries on hook or validation errors.
- **Keybindings:** Standard CLI invocation flags.

**Technical Constraints (NFR):**
- **[Performance]:** Batch update execution latency under 150ms (excluding external user hooks duration).
- **[Integrity]:** Atomic file updates preventing partial state writes.
- **[Observability]:** Explicit console status log listing all impacted task identifiers.

### [US-2] Deduplicated Verification Hook Execution
**User Story:** AS A developer with automated test hooks, I WANT verification hooks to run once per batch operation, SO THAT my test suites are not wastefully re-executed for every individual task in the batch.

**Scenarios:**
1. **[Happy Path - Single Hook Run]** GIVEN multiple tasks being marked as finished in a single batch WHEN verification hooks are configured THEN all applicable hooks are executed exactly once before persisting the task state changes.
2. **[Happy Path - Phase and Spec Hooks]** GIVEN a batch that includes the final task of a phase or the entire roadmap WHEN the batch completes successfully THEN phase and project completion hooks trigger appropriately in deduplicated sequence.

**UI/UX Specifics:**
- **View/Component:** Verification progress feedback.
- **Feedback Logic:** Display stdout and stderr streams of failed hooks clearly with non-zero exit code indicators.
- **Keybindings:** Standard terminal stream output.

**Technical Constraints (NFR):**
- **[Performance]:** Zero duplicate execution of identical hook commands within the same batch pass.
- **[Safety & Security]:** Block state transition if any hook command returns a non-zero exit status.

### [US-3] Subagent Orchestration and Boundary Directives in Blueprint
**User Story:** AS AN AI orchestrator executing implementation roadmaps, I WANT precise guidelines on worker allocation and task boundaries, SO THAT subagents are systematically spawned with safe concurrency and never corrupt task state files.

**Scenarios:**
1. **[Happy Path - Worker Sizing Guidance]** GIVEN an implementation phase with varying batch complexities WHEN an orchestrator prepares task delegation THEN it scales subagent allocation dynamically (1 for small tasks, 2-3 for standard tasks, up to 4 for complex tasks).
2. **[Happy Path - Mandatory Subagent Utilization]** GIVEN an environment with subagent invocation capabilities WHEN executing implementation tasks THEN the orchestrator always delegates code modifications to subagents.
3. **[Edge Case - Direct State Mutation Prohibition]** GIVEN spawned subagents modifying codebase targets WHEN they finish their code changes THEN they return control to the primary orchestrator without touching the specification state files directly.

**UI/UX Specifics:**
- **View/Component:** Orchestrator blueprint instructions.
- **Feedback Logic:** Clear structured directives, mission brief schema, and role definitions.
- **Keybindings:** Standard markdown blueprint specification.

**Technical Constraints (NFR):**
- **[Performance]:** Clear concise directives minimizing instruction prompt overhead.
- **[Integrity]:** Absolute prohibition against secondary workers directly editing task files.

## 4. Business Invariants
- An update operation referencing multiple tasks must be all-or-nothing: either all specified tasks update, or none update.
- Secondary workers and subagents must never write directly to specification files (`tasks.md`, `requirements.md`, `design.md`).
- Automated verification hooks for finishing tasks must execute in deduplicated order and pass with zero exit status before any state is finalized.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Sub-150ms processing overhead for CLI updates (excluding user hook runtimes).
- **[Reliability]:** Fail-fast error reporting with zero silent failures or half-applied mutations.
- **[Security]:** Path validation using secure project root boundary restrictions.
- **[Maintainability]:** Clean separation between CLI parsing, core task manipulation, and agent blueprint directives.
