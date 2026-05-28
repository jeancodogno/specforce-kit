---
slug: 20260520-1434-fix-config-init-logic
lens: Backend-heavy
---

# Implementation Roadmap: fix config initialization logic

## 1. Execution Strategy
- **Gravity Order:** Test Reproduction -> CLI Implementation -> Verification

## 2. Tasks

### Phase 1: Test Reproduction (Red)

- [x] T1.1: [TEST] Add failing regression test for config creation in update flow
**Target:** `src/internal/cli/cli_test.go`
**Context:** [FIX-1, FIX-2]

**Action Steps:**
- Add `TestHandleInit_EnsuresConfigExistsInUpdateFlow` function.
- Mock a temporary directory with `.specforce/` existing but no `config.yaml`.
- Call `executor.HandleInit`.
- Assert that `config.yaml` exists after the call.

**Acceptance Check:**
Run `go test -v src/internal/cli/cli_test.go src/internal/cli/cli.go ...` and verify the new test fails.

### Phase 2: Implementation (Green)

- [x] T2.1: [CODE] Move config guarantee to a shared point in HandleInit
**Target:** `src/internal/cli/cli.go`
**Context:** [FIX-1, FIX-2]

**Action Steps:**
- Identify the common exit/branch point in `HandleInit`.
- Ensure `core.EnsureConfigExists(".")` is called before entering either `handleUpdateFlow` or `handleNewInitFlow`.
- Remove the redundant call from `handleNewInitFlow`.

**Acceptance Check:**
Run the newly created test case and verify it passes.

### Phase 3: Validation & Cleanup

- [x] T3.1: [CLI] Comprehensive test suite run
**Target:** `Global Scope`
**Context:** [FIX-1, FIX-2]

**Action Steps:**
- Run all tests in `src/internal/cli/` and `src/internal/core/` to ensure no regressions.

**Acceptance Check:**
`go test ./src/internal/cli/... ./src/internal/core/...`
