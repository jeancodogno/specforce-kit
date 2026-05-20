---
slug: 20260520-1619-fix-spec-list-json-null
lens: Bugfix
---

# Implementation Roadmap: Fix Spec List JSON Null Output

## 1. Execution Strategy
- **Gravity Order:** Unit Test (Red) -> Implementation -> Unit Test (Green) -> Integration Verification.

## 2. Tasks

### Phase 1: Verification & Regression

- [x] T1.1: [TEST] Add regression test for empty spec list JSON
**Target:** `src/internal/spec/discovery_test.go`
**Context:** [FIX-1]

**Action Steps:**
- Create `src/internal/spec/discovery_test.go` if it doesn't exist.
- Add `TestListActiveSpecs_Empty` test case.
- Mock an empty `.specforce/specs` directory.
- Assert that `ListActiveSpecs` returns a non-nil slice and `len(specs) == 0`.

**Verification (TDD):**
`go test ./src/internal/spec/...` (Should fail or return nil slice depending on current state).

### Phase 2: Implementation

- [x] T2.1: [CODE] Fix nil slice initialization in ListActiveSpecs
**Target:** `src/internal/spec/discovery.go`
**Context:** [FIX-1], [FIX-2]

**Action Steps:**
- Locate `ListActiveSpecs` function.
- Replace `var specs []SpecInfo` with `specs := []SpecInfo{}`.

**Verification (TDD):**
`go test ./src/internal/spec/...` (Should pass with non-nil empty slice).

### Phase 3: CLI Integration Verification

- [x] T3.1: [TEST] Update CLI test for empty JSON output
**Target:** `src/internal/cli/spec_test.go`
**Context:** [FIX-2]

**Action Steps:**
- Update `TestHandleSpecList` to verify JSON output is `[]` when no specs are found.

**Verification (TDD):**
`go test ./src/internal/cli/spec_test.go`

- [x] T3.2: [VERIFY] Manual CLI verification
**Target:** `CLI Command`
**Context:** [FIX-1], [FIX-2]

**Action Steps:**
- Ensure current environment has no active specs.
- Run `go run src/cmd/specforce/main.go spec list --json`.

**Verification (TDD):**
Verify output is exactly `[]`.
