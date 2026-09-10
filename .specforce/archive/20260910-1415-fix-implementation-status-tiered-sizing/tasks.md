---
slug: 20260910-1415-fix-implementation-status-tiered-sizing
lens: Balanced full-stack
---

# Implementation Roadmap: Tiered Sizing Invariance in Implementation Status and Artifact Dependencies

## 1. Execution Strategy
- **Gravity Order:** Service/Implementation Status Layer (`service.go`, `implementation.go`) -> Spec Status Dependency Layer (`status.go`) -> CLI Integration Tests (`integration_test.go`) -> Full Suite Verification.

## 2. Tasks

### Phase 1: Size-Aware Implementation Readiness Check

- [x] T1.1: [TEST: RED] Add failing unit tests for GetImplementationStatus with tiered sizing
**Target:** `src/internal/spec/service_test.go`
**Context:** [FIX-1]

**Action Steps:**
- Add `TestGetImplementationStatus_TieredSizing` covering `small` spec with only `tasks.md`, asserting `status == "ready"` and `len(missing_artifacts) == 0`.
- Add test case covering `medium` spec with `requirements.md` and `tasks.md`, asserting `status == "ready"` and `missing_artifacts` does not include `design.md`.
- Add test case covering `large` and `complex` specs missing `design.md`, asserting `status == "blocked"` and `missing_artifacts` contains `design.md`.

**Acceptance Check:**
`go test ./src/internal/spec -run TestGetImplementationStatus_TieredSizing -v` fails showing status is blocked for small and medium specs.

- [x] T1.2: [CODE: GREEN] Implement size-aware artifact checking in GetImplementationStatus
**Target:** `src/internal/spec/service.go`
**Context:** [FIX-1]

**Action Steps:**
- In `GetImplementationStatus`, inspect `meta.Type` and `meta.Size` to retrieve required artifacts via `s.registry.ListForTypeAndSize(meta.Type, meta.Size)`.
- Replace hardcoded triad artifact check with dynamic verification of required artifacts for the active size.
- Populate `report.MissingArtifacts` only with absent required artifacts, setting `report.Status = "blocked"` only when required artifacts are missing.

**Acceptance Check:**
`go test ./src/internal/spec -run TestGetImplementationStatus_TieredSizing -v` passes with 100% success.

### Phase 2: Size-Aware Artifact Dependency Evaluation in Spec Status

- [x] T2.1: [TEST: RED] Add failing unit tests in status_test.go for artifact dependencies under tiered sizing
**Target:** `src/internal/spec/status_test.go`
**Context:** [FIX-2]

**Action Steps:**
- Add `TestSpecStatus_Dependency_TieredSizing` initializing a `small` spec containing only `tasks.md`.
- Assert that `tasks.md` artifact status reports `Blocked: false` even though `design.md` does not exist.
- Add test case for `medium` spec with `requirements.md` and `tasks.md`, asserting `Blocked: false`.
- Assert that an artifact remains `Blocked: true` if an active upstream dependency (e.g. `requirements.md` in `medium`) is absent.

**Acceptance Check:**
`go test ./src/internal/spec -run TestSpecStatus_Dependency_TieredSizing -v` fails showing `tasks.md` is marked blocked: true.

- [x] T2.2: [CODE: GREEN] Update processArtifactStatus to respect allowed artifacts for the spec size
**Target:** `src/internal/spec/status.go`
**Context:** [FIX-2]

**Action Steps:**
- Pass the allowed/required artifact set for the active size into `processArtifactStatus` (or filter dependencies against allowed artifacts).
- In `processArtifactStatus`, evaluate `blocked = true` only if `art.Dependency` is present in the allowed artifact set for the spec size and missing from disk.
- If `art.Dependency` is not required for the spec size (e.g. `design` in `small` or `medium`), keep `blocked = false`.

**Acceptance Check:**
`go test ./src/internal/spec -run TestSpecStatus_Dependency_TieredSizing -v` passes with 100% success.

### Phase 3: End-to-End CLI Integration and Regression Verification

- [x] T3.1: [TEST: RED] Add CLI integration test for implementation status on small and medium specs
**Target:** `src/internal/cli/integration_test.go`
**Context:** [FIX-1], [FIX-2]

**Action Steps:**
- Add `TestIntegration_TieredSizing_ImplementationStatus` initializing a `small` bug spec via CLI.
- Write valid `tasks.md` and execute `specforce implementation status <slug> --json`.
- Assert JSON output returns `"status": "ready"` and `"missing_artifacts": []`.
- Execute `specforce spec status <slug> --json` and assert progress 100% with no blocked artifacts.

**Acceptance Check:**
`go test ./src/internal/cli -run TestIntegration_TieredSizing_ImplementationStatus -v` executes and confirms CLI behavior.

- [x] T3.2: [CODE: GREEN] Verify full test suite and clean execution
**Target:** `src/internal/cli/spec_test.go`
**Context:** [FIX-1], [FIX-2]

**Action Steps:**
- Execute full test suite across the repository (`go test ./... -count=1`).
- Verify backward compatibility for `large` and `complex` specs without regressions.
- Verify that `specforce spec audit` passes with zero coherence errors.

**Acceptance Check:**
`go test ./... -count=1` passes with 0 failures across all packages.
