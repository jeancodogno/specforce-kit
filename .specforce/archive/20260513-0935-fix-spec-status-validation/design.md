---
slug: 20260513-0935-fix-spec-status-validation
lens: Bugfix
---

# Technical Design: Spec Status Validation & Redundant Rewriting (Fix Blueprint)

## 1. Code Path Inventory

- `src/internal/spec/tasks.go` -> Modify `ValidateTasks` and `taskValidationState` to enforce at least one phase.
- `src/internal/spec/status.go` -> Modify `GetStatus` to only call `processArtifactStatus` with validation intent when `Found == Total`.
- `src/internal/agent/kit/commands/spec.yaml` -> Refine orchestrator instructions to enforce batch processing and clarify "verification" vs "discovery" status checks.

## 2. Regression Strategy (Verification Plan)

- **Unit Tests:**
    - `src/internal/spec/tasks_test.go`: Add test case for malformed headers returning error.
    - `src/internal/spec/status_test.go`: Add test case verifying that `ValidationErrors` are empty if `Progress < 100`.
    - `src/internal/spec/bug_status_test.go`: Verify tasks validation works for bug specs.
- **Manual Verification:**
    - Create a malformed `tasks.md` in a dummy spec and run `specforce spec status dummy --json`.
    - Verify that status remains "valid" (existence only) while requirements/design are missing.
    - Verify that status becomes "invalid" (deep validation) once all artifacts are present.

## 3. Side Effects & Risks
- **Agent Confusion:** If an agent relies on `spec status` to fix a malformed `tasks.md` *before* finishing other files, it won't see the errors until the end.
    - **Mitigation:** This is the intended behavior to avoid redundant rewrites. Agents should focus on finishing the batch first.
- **Complexity:** Adding conditional logic to `GetStatus` might make the code harder to reason about.
    - **Mitigation:** Implement a two-pass approach or a clearly documented flag in `processArtifactStatus`.

## 4. Proposed Fix (Abstract Logic)

### Task Validation (tasks.go)
```go
func (s *taskValidationState) finalize() {
    if s.currentPhase == 0 {
        s.errors = append(s.errors, "No valid Phase found")
    }
}
```

### Conditional Validation (status.go)
```go
// Pass 1: Scan existence
existsMap, foundCount := scanExistence(...)

// Pass 2: Process details
for _, art := range artifacts {
    shouldValidate := (foundCount == totalCount)
    artStatus := processArtifactStatus(..., shouldValidate)
}
```

### Orchestrator Instructions (spec.yaml)
- Explicitly state: "You MUST NOT check status to verify individual artifact correctness. Only check status ONCE at the start to plan, and ONCE at the end to verify the full set."
