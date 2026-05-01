# Requirements: Fix macOS CI and Improve Coverage

## 1. Problem Statement
The CI pipeline is failing on `macos-latest` due to two main reasons:
1. **Failing Tests:** `TestScanProject_WithWorktrees` fails because of path canonicalization discrepancies (symlinks like `/var` -> `/private/var`). `TestResolveMapping_TildeExpansionFailure` likely fails because `os.UserHomeDir` behaves differently on macOS.
2. **Coverage Gate:** The `upgrade` package has ~65% coverage, and although the core average is high, the failing tests prevent the CI from completing the coverage check step on macOS.

## 2. Functional Requirements
- **FR-1: Path Resilience:** The project scanner MUST correctly identify the main root even when accessed via symlinked paths (common on macOS `/var`).
- **FR-2: Test Robustness:** Platform-specific tests MUST handle environment differences (like home directory resolution fallbacks) gracefully or skip execution when the test prerequisite (forcing failure) cannot be met.
- **FR-3: Coverage Threshold:** The `upgrade` package MUST reach at least 80% statement coverage to ensure project health and meet internal quality standards.

## 3. Acceptance Criteria
- **AC-1:** `TestScanProject_WithWorktrees` passes on both Linux and macOS.
- **AC-2:** `TestResolveMapping_TildeExpansionFailure` (and similar tests) passes on macOS or is skipped if the failure cannot be reproduced.
- **AC-3:** `go test -cover ./src/internal/upgrade/...` reports at least 80% coverage.
- **AC-4:** The GitHub Action workflow for `macos-latest` completes the `test-coverage` job successfully.
