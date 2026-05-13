---
slug: 20260513-0935-fix-spec-status-validation
lens: Backend-heavy
---

# Implementation Roadmap: Spec Status Validation & Redundant Rewriting

## 1. Execution Strategy
- **Gravity Order:** Domain Logic (tasks.go) -> Service Logic (status.go) -> Orchestrator Config (spec.yaml) -> Verification.

## 2. Tasks

### Phase 1: Robust Task Validation

- [x] T1.1: [CODE] Enforce Phase detection in ValidateTasks
**Target:** `src/internal/spec/tasks.go`
**Context:** [FIX-1]

**Action Steps:**
- Modify `ValidateTasks` to check `state.currentPhase == 0` after the loop.
- Append error "No valid Phase (### Phase N: Name) found in tasks.md" if zero phases were found.

**Verification (TDD):**
Run `go test -v src/internal/spec/tasks_validation_repro_test.go` (created during discovery) and verify it passes.

### Phase 2: Conditional Deep Validation

- [x] T2.1: [CODE] Modify GetStatus for conditional validation
**Target:** `src/internal/spec/status.go`
**Context:** [FIX-2]

**Action Steps:**
- Update `scanArtifactExistence` to return the count of found artifacts.
- Update `processArtifactStatus` signature to accept `validate bool`.
- In `GetStatus`, set `shouldValidate := (foundCount == totalCount)`.
- Pass `shouldValidate` to `processArtifactStatus`.
- In `processArtifactStatus`, only call `ValidateTasks` if `validate` is true.

**Verification (TDD):**
Add a test case in `src/internal/spec/status_test.go` where `found < total` and `tasks.md` is malformed, verifying `ValidationErrors` is empty.

### Phase 3: Orchestrator Hardening

- [x] T3.1: [DOCS] Refine spec.yaml instructions
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [FIX-2]

**Action Steps:**
- Update `### 2. Artifact Processing (Batch Execution)` to explicitly forbid status polling between artifacts.
- Strengthen the `### 3. Verification & Handoff` section to explain that validation errors only appear when progress is 100%.

**Verification (TDD):**
Inspect the file content and ensure the instructions are clear and unambiguous.

### Phase 4: Final Verification

- [x] T4.1: [VERIFY] Bug spec status consistency
**Target:** `Global Scope`
**Context:** [FIX-3]

**Action Steps:**
- Run `specforce spec status teste --json` (with a broken tasks.md).
- Verify that it reports errors correctly if it's the only artifact or if the spec is complete.

**Verification (TDD):**
Manual CLI execution.

## 3. Pre-emptive Mitigations
- **Risk:** Backward compatibility with existing tests that might rely on eager validation.
- **Mitigation:** Update `status_test.go` to ensure tests set up all required files for deep validation checks.
