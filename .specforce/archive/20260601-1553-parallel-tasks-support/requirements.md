# Requirements: Parallel Task Support

This document defines the functional requirements and business rules for enabling parallel execution of tasks within a Specforce roadmap.

## 1. Persona & Value
- **Persona:** Senior Developer / AI Implementation Agent (`spf.implement`).
- **User:** Developers seeking to optimize implementation speed by leveraging multiple concurrent agents or parallelizable tool calls.
- **Value:** Reduces the critical path of a feature implementation by allowing independent, non-conflicting tasks to run simultaneously.

## 2. Success Metrics
- **Business Metric:** Up to 50% reduction in total implementation time for highly parallelizable specifications.
- **Performance Target:** CLI validation for task file conflicts SHALL execute in < 1s for roadmaps with up to 100 tasks.
- **UX Efficiency:** Developers SHALL be able to mark parallel dependencies using a single, readable line of metadata in `tasks.md`.

## 3. Functional Requirements

### Parser Enhancement
- **[US-1] Parallel Metadata Recognition:** The `tasks.md` parser SHALL support the identification of parallel tasks using the syntax: `**Parallel With: T<Phase>.<Index>, ...**`.
- **[US-2] Block-Level Placement:** The metadata MUST be recognized when placed anywhere within the task description block (between the task header and the next task/phase header).
- **[US-3] Reference Resolution:** The parser MUST resolve the provided task IDs (e.g., `T1.2`) to their corresponding task objects within the same roadmap.

### Validation Logic (State Machine & Integrity)
- **[US-4] Conflict Detection (File Isolation):** The CLI MUST perform strict validation to ensure that any two tasks marked as parallel do NOT have overlapping `Target` files. 
    - *Example:* If T1.1 targets `src/a.go` and T1.2 targets `src/a.go`, they CANNOT be parallel.
- **[US-5] Symmetry Enforcement:** If T1.1 is marked as "Parallel With: T1.2", the CLI SHALL treat them as mutually parallel regardless of whether T1.2 explicitly mentions T1.1.
- **[US-6] Reference Integrity:** The CLI MUST reject roadmaps where a task references a non-existent task ID in its `Parallel With` metadata.
- **[US-7] Phase Boundary Adherence:** Tasks marked as parallel MUST belong to the same Phase. Parallelism across phases is prohibited to maintain structural order.

### CLI Output & Integration
- **[US-8] JSON Schema Update:** The output of `specforce spec status --json` MUST include a `parallel_with` array field for each task containing the IDs of its parallel peers.
- **[US-9] Explicit Error Reporting:** Validation errors related to parallel conflicts MUST identify the specific tasks and the conflicting `Target` files.

## 4. Business Rules & Invariants
- **Strict File Isolation:** Parallel tasks MUST operate on disjoint sets of files to prevent race conditions or merge conflicts during implementation.
- **Sequential Default:** In the absence of `Parallel With` metadata, all tasks MUST be treated as strictly sequential (default behavior).
- **Readiness Rule:** A task is considered "Ready for Parallel Execution" ONLY when all its preceding sequential dependencies (all tasks appearing earlier in the `tasks.md` that are NOT marked as parallel with it) are `FINISHED`.
- **Backward Compatibility:** Existing `tasks.md` files without parallel metadata MUST remain valid and functional.

## 5. Acceptance Criteria

### Scenario: Valid Parallel Task Definition
- **GIVEN** a phase with two tasks:
    - Task T1.1: Target `src/auth.go`, Metadata `**Parallel With: T1.2**`
    - Task T1.2: Target `src/user.go`
- **WHEN** running `specforce spec status`
- **THEN** validation MUST pass.
- **AND** the JSON output for T1.1 MUST include `"parallel_with": ["T1.2"]`.

### Scenario: Conflict Detection (Target Overlap)
- **GIVEN** a phase with two tasks:
    - Task T1.1: Target `src/common.go`, Metadata `**Parallel With: T1.2**`
    - Task T1.2: Target `src/common.go`
- **WHEN** running `specforce spec status`
- **THEN** validation MUST fail with an error: `Parallel conflict: T1.1 and T1.2 both target "src/common.go"`.

### Scenario: Invalid Reference
- **GIVEN** a task T1.1 with Metadata `**Parallel With: T9.9**` (non-existent)
- **WHEN** running `specforce spec status`
- **THEN** validation MUST fail with an error: `Task T1.1 references non-existent task T9.9`.

### Scenario: Cross-Phase Parallelism Prohibition
- **GIVEN** Task T1.1 in Phase 1 and Task T2.1 in Phase 2
- **WHEN** T1.1 is marked as `**Parallel With: T2.1**`
- **THEN** validation MUST fail with an error: `Parallel tasks must belong to the same phase (T1.1 and T2.1)`.

### Scenario: Execution Readiness
- **GIVEN** a sequence T1.1 (Seq) -> T1.2 (Parallel with T1.3) -> T1.3 (Parallel with T1.2)
- **WHEN** T1.1 is `PENDING`
- **THEN** T1.2 and T1.3 MUST NOT be reported as "Ready" in the implementation status.
- **WHEN** T1.1 is marked as `FINISHED`
- **THEN** both T1.2 and T1.3 MUST be reported as "Ready".
