---
slug: fix-macos-ci-and-coverage
lens: Technical-Debt
---

# Implementation Roadmap: Fix macOS CI and Improve Coverage

## 1. Execution Strategy
- **Gravity Order:** Fix macOS-specific test failures -> Add missing tests for `upgrade` package -> Verify coverage reaches 80% target.

## 2. Tasks

### Phase 1: Fix macOS Test Failures

#### T1.1: [CODE] Implement path canonicalization in ScanProject
**State:** `[FINISHED]`
**Target:** `src/internal/spec/scanner.go`
**Context:** [Design 1.1](design.md)

**Action Steps:**
- Add `evalPath(path string) string` helper function to `scanner.go` that uses `filepath.EvalSymlinks` and `filepath.Abs`.
- Update `ScanProject` to use `evalPath` for both `projectRoot` and `wt.Path` before comparison in `isMainRoot`.

**Verification (TDD):**
`go test -v -run TestScanProject_WithWorktrees ./src/internal/spec/...`

#### T1.2: [TEST] Make tilde expansion test resilient to platform fallbacks
**State:** `[FINISHED]`
**Target:** `src/internal/agent/translator_test.go`
**Context:** [Design 1.2](design.md)

**Action Steps:**
- Update `TestResolveMapping_TildeExpansionFailure` in `translator_test.go`.
- After unsetting `HOME`, call `os.UserHomeDir()`.
- If `err == nil`, call `t.Skip` with a message explaining that the platform has a fallback for home directory.

**Verification (TDD):**
`go test -v -run TestResolveMapping_TildeExpansionFailure ./src/internal/agent/...`

### Phase 2: Improve upgrade Package Coverage (>80%)

#### T2.1: [TEST] Add tests for missing constructors and simple getters
**State:** `[FINISHED]`
**Target:** `src/internal/upgrade/installer_binary_test.go`
**Context:** [Design 1.3](design.md)

**Action Steps:**
- Add `TestNewBinaryInstaller` to `installer_binary_test.go`.
- Add `TestNewNPMInstaller` to `installer_npm_test.go`.
- Add `TestNewStateManager` to `state_test.go`.

**Verification (TDD):**
`go test -v -cover ./src/internal/upgrade/...`

#### T2.2: [TEST] Implement tests for binary replacement and move fallback
**State:** `[FINISHED]`
**Target:** `src/internal/upgrade/installer_binary_test.go`
**Context:** [Design 1.3](design.md)

**Action Steps:**
- Add `TestBinaryInstaller_Replace` (mocking `os.Executable` or testing via a helper).
- Add `TestBinaryInstaller_MoveFile_Fallback` to test the copy+delete logic in `moveFile`.

**Verification (TDD):**
`go test -v -cover ./src/internal/upgrade/...`

#### T2.3: [TEST] Test Service.PerformUpgrade and error paths
**State:** `[FINISHED]`
**Target:** `src/internal/upgrade/service_test.go`
**Context:** [Design 1.3](design.md)

**Action Steps:**
- Add `TestService_PerformUpgrade` to `service_test.go`.
- Add tests for failure scenarios in `BinaryInstaller.DownloadAndVerify` (e.g., checksum mismatch, download error).

**Verification (TDD):**
`go test -v -cover ./src/internal/upgrade/...`

### Phase 3: Final Verification

#### T3.1: [QA] Verify total project coverage and gate compliance
**State:** `[FINISHED]`
**Target:** `Global Scope`
**Context:** [REQ-3](requirements.md)

**Action Steps:**
- Run the full test suite with coverage.
- Verify that the `upgrade` package coverage is above 80%.

**Verification (TDD):**
`go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | grep "github.com/jeancodogno/specforce-kit/src/internal/upgrade" | awk '{print $NF}'`
