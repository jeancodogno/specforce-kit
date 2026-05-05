---
slug: 20260505-0926-tasks-md-validation-in-status
lens: Backend-heavy
---

# Feature: Tasks Markdown Validation in Status

## 1. Context & Value
The `spec status` command currently checks for the existence of the `tasks.md` file but does not verify its structural integrity. This feature adds a validation layer to ensure that `tasks.md` adheres to the mandatory Specforce hierarchical format (Phases as H3, Tasks as H4 with specific IDs). This prevents malformed task lists from breaking implementation workflows or reporting inaccurate progress.

## 2. Out of Scope (Anti-Goals)
- Automatic fixing or refactoring of malformed `tasks.md` files.
- Validation of the content/logic within the tasks themselves.
- Validation of other specification artifacts (requirements.md, design.md) in this specific feature.
- Changing the existing CLI output format for `spec status` beyond adding error messages.

## 3. Acceptance Criteria (BDD)

### [REQ-1] Hierarchical Structure Validation
**User Story:** AS A developer, I WANT TO be alerted if my task list does not follow the H3/H4 hierarchy or if task internal blocks are malformed, SO THAT the orchestration tools can correctly parse the implementation phases and task details.

**Scenarios:**
1. **[Happy Path]** GIVEN a `tasks.md` file with Phase headers (###), Task headers (####), and all mandatory fields (Target, Context, Action Steps, Verification) WHEN I run `specforce spec status` THEN the command succeeds.
2. **[Edge Case]** GIVEN a `tasks.md` file where a Task header (`####`) appears before any Phase header (`###`) WHEN I run `specforce spec status` THEN the command reports a validation error: "Task found before any Phase definition".
3. **[Edge Case]** GIVEN a task block missing the mandatory `**Verification (TDD):**` section WHEN I run `specforce spec status` THEN the command reports a validation error: "Task {ID} is missing mandatory Verification (TDD) section".
4. **[Edge Case]** GIVEN a task block where the `**Target:**` is empty or not formatted correctly WHEN I run `specforce spec status` THEN the command reports a validation error for that task.

### [REQ-2] Task ID Naming and Sequentiality
**User Story:** AS A developer, I WANT TO ensure my tasks are named using the `T{Phase}.{Task}` format and are strictly sequential, SO THAT they can be uniquely identified and tracked.

**Scenarios:**
1. **[Happy Path]** GIVEN sequential tasks `T1.1`, `T1.2` under `Phase 1` WHEN I run `specforce spec status` THEN no errors are reported.
2. **[Edge Case]** GIVEN a gap in sequentiality (e.g., `T1.1` followed by `T1.3`) WHEN I run `specforce spec status` THEN the command reports: "Task sequence gap: expected T1.2, found T1.3".
3. **[Edge Case]** GIVEN a task defined as `#### T2.1: My Task` under `### Phase 1: My Phase` WHEN I run `specforce spec status` THEN the command reports: "Task ID T2.1 does not match the parent Phase 1".

### [REQ-3] Status Reporting
**User Story:** AS A developer, I WANT TO see a list of all formatting errors in `spec status`, SO THAT I can fix them before starting implementation.

**Scenarios:**
1. **[Happy Path]** GIVEN a malformed `tasks.md` WHEN I run `specforce spec status` THEN the output includes a section listing all identified formatting errors with line numbers or context.
2. **[Edge Case]** GIVEN multiple formatting errors in different parts of the file WHEN I run `specforce spec status` THEN all errors are reported at once (batch validation).

## 4. Business Invariants
- A task list MUST NOT be considered "valid" if it contains any Phase without at least one Task.
- Task IDs MUST be sequential within their Phase (e.g., T1.1, T1.2).
- The `tasks.md` file MUST NOT use H1 or H2 headers for anything other than the main feature title (if present).
