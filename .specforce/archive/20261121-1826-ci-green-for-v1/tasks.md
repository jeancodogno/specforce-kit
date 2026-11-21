---
slug: 20261121-1826-ci-green-for-v1
lens: Integration
---

# Implementation Roadmap: Green Continuous Integration for the 1.0.0 Release

## 1. Execution Strategy
- **Gravity Order:** toolchain and security scan (Linux-verifiable) -> the two Windows product defects -> portable and isolated tests -> constitution rule. Each batch is committed to `1.x` (no Co-Authored-By trailer) and pushed so the CI run proves the Windows and macOS part (approved by the user on 2026-11-21); no tags.
- **Windows evidence:** Windows and macOS behaviour can only run in CI. The RED state of Phase 2 and Phase 3 is the failing CI run `44947632566` (job `Unit Tests (windows-latest)` with 61 failures, job `Unit Tests (macos-latest)` with 1 failure); those existing tests are the regression tests and MUST NOT be weakened. Local acceptance checks prove Linux behaviour plus a Windows cross-compile of the tests (`GOOS=windows go vet`); the Windows run itself is Black-box Scenario S1.
- **Test pairing:** Phase 1 has no testable behaviour beyond the scanners themselves; Phase 2 GREEN tasks are paired with the already failing Windows tests named in each task; Phase 3 changes only tests; Phase 4 is documentation.
- **Module impact:** none (`delta.impact: none`): no module behaviour changes; Windows now behaves as the module specs already describe.

## 2. Tasks

### Phase 1: Toolchain and security scan

- [x] T1.1: [CONFIG] Pin the Go toolchain to 1.26.11 everywhere
**Target:** `go.mod`, `.github/workflows/ci.yml`, `.github/workflows/release.yml`
**Context:** [US-1]

**Constraints:**
- `go.mod` declares `go 1.26.11`; every `actions/setup-go` step in both workflows uses `go-version: "1.26.11"`. Do not touch any other workflow setting.

**Acceptance Check:**
`grep -q '^go 1.26.11$' go.mod && [ "$(grep -c 'go-version: "1.26.11"' .github/workflows/ci.yml)" -eq "$(grep -c 'setup-go' .github/workflows/ci.yml)" ] && [ "$(grep -c 'go-version: "1.26.11"' .github/workflows/release.yml)" -eq "$(grep -c 'setup-go' .github/workflows/release.yml)" ] && GOTOOLCHAIN=go1.26.11 $(go env GOPATH)/bin/govulncheck ./... | grep -q 'No vulnerabilities found' && echo OK` prints `OK`.

- [x] T1.2: [CODE] Justify the false-positive credential warning
**Target:** `src/internal/spec/delta.go`
**Context:** [US-1]
**Parallel With:** T1.1

**Constraints:**
- A single `#nosec G101` on the `deltaTokenInner` declaration with the reason (a regular expression for specification identifiers, not a credential), following the existing `#nosec <rule> - <reason>` style in the repository. No other change.

**Acceptance Check:**
`$(go env GOPATH)/bin/gosec -fmt=text ./... 2>/dev/null | grep -q 'Issues : 0' && go test ./src/internal/spec/ && echo OK` prints `OK`.

- [x] T1.3: [CODE] Brand asset generator passes the current static security scan
**Target:** `scripts/brandgen/brand.go`, `scripts/brandgen/main.go`
**Context:** [US-1]

**Constraints:**
- Fix, do not suppress: colour channels are masked to 8 bits before conversion (`brand.go:43`), the output folder is created with `0o750` and files are written with `0o600`. The generated asset bytes must stay identical (run the generator into a temporary folder before and after and compare).

**Acceptance Check:**
`go run github.com/securego/gosec/v2/cmd/gosec@v2.30.0 -fmt=text ./... 2>/dev/null | grep -q 'Issues : 0' && (cd scripts/brandgen && go test ./...) && echo OK` prints `OK` (`scripts/brandgen` is its own Go module).

- [x] T1.4: [CONFIG] Pin the scanner versions used by CI
**Target:** `Makefile`
**Context:** [US-1]

**Constraints:**
- `bootstrap` installs `gosec@v2.30.0`, `govulncheck@v1.10.0` and `gremlins@v0.6.0` instead of `@latest`; no other Makefile change.

**Acceptance Check:**
`grep -q 'gosec@v2.30.0' Makefile && grep -q 'govulncheck@v1.10.0' Makefile && grep -q 'gremlins@v0.6.0' Makefile && ! grep -q '@latest' Makefile && echo OK` prints `OK`.

### Phase 2: Windows product defects

- [x] T2.1: [GREEN] Embedded-kit lookups use the embedded file system's path rules
**Target:** `src/internal/agent/instructions.go`, `src/internal/agent/registry.go`
**Context:** [US-2]

**Constraints:**
- Paths passed to `io/fs` functions are joined with forward slashes (`path.Join`), never `filepath.Join`; paths on the real disk keep `filepath`. Paired RED tests (already failing on Windows CI): `TestInstructionMerging`, `TestRegistry_ScanSkills_Metadata` (package `agent`) and `TestHandleArchiveInstructions` (package `cli`).

**Acceptance Check:**
`go test ./src/internal/agent/ ./src/internal/cli/ && GOOS=windows go vet ./src/internal/agent/ ./src/internal/cli/ && ! grep -nE 'fs\.(ReadFile|ReadDir|Stat|Sub|Glob)\([^)]*filepath\.' src/internal/agent/*.go && echo OK` prints `OK` (the Windows proof is S1).

- [x] T2.2: [GREEN] Tools-only installation filter accepts native Windows paths
**Target:** `src/internal/project/refresh_integration_test.go` (diagnostics, done), `src/internal/installer/installer.go`
**Context:** [US-2]

**Constraints:**
- Cause (found in CI run `44952829546`: only `.claude/rules/AGENTS.md` was created): `ShouldInstall` compares the mapped path, which uses `\` on Windows, against `/`-terminated tool prefixes and the `.specforce/` prefix. Normalize the path with `filepath.ToSlash` before the prefix checks; keep the absolute-path rule and the `.specforce/` exclusion. Paired RED test: `TestService_UpdateTools_RefreshesAllCategories` (failing on Windows CI); existing `installer` tests must keep passing.
- The refresh test's assertions must not be weakened.

**Acceptance Check:**
`go test ./src/internal/installer/ ./src/internal/project/ && GOOS=windows go vet ./src/internal/installer/ ./src/internal/project/ && grep -q 'filepath.ToSlash' src/internal/installer/installer.go && echo OK` prints `OK`; the Windows proof is S1.

### Phase 3: Portable and isolated tests

- [x] T3.1: [CODE] Updater tests: isolated home, resolved paths, Windows permission skips
**Target:** `src/internal/upgrade/service_test.go`, `src/internal/upgrade/state_test.go`, `src/internal/upgrade/installer_binary_test.go`
**Context:** [US-3]
**Parallel With:** T3.2, T3.3, T3.4

**Constraints:**
- Every place that sets `HOME` to a temporary directory also sets `USERPROFILE` to the same directory (a small shared test helper in the package is allowed).
- Expected executable paths are resolved the same way the product resolves them (`filepath.EvalSymlinks` of the expected path), which covers `/private/var` on macOS and long names on Windows.
- Tests that rely on `os.Chmod` read-only directories or on file mode bits skip on Windows only, with the reason "permission bits do not apply on Windows"; they keep running elsewhere.

**Acceptance Check:**
`go test -count=1 ./src/internal/upgrade/ && GOOS=windows go vet ./src/internal/upgrade/ && grep -c 'USERPROFILE' src/internal/upgrade/*_test.go | awk -F: '{s+=$2} END {exit !(s>0)}' && echo OK` prints `OK`.

- [x] T3.2: [CODE] Agent, CLI and constitution tests: home isolation and native separators
**Target:** `src/internal/agent/translator_test.go`, `src/internal/cli/cobra/root_notice_test.go`, `src/internal/constitution/registry_test.go`
**Context:** [US-3]
**Parallel With:** T3.1, T3.3, T3.4

**Constraints:**
- Set `USERPROFILE` wherever `HOME` is set; compare expected paths built with `filepath.Join`/`filepath.FromSlash` instead of hard-coded `/` (product output stays native, per the approved decision).

**Acceptance Check:**
`go test -count=1 ./src/internal/agent/ ./src/internal/cli/cobra/ ./src/internal/constitution/ && GOOS=windows go vet ./src/internal/agent/ ./src/internal/cli/cobra/ ./src/internal/constitution/ && echo OK` prints `OK`.

- [x] T3.3: [CODE] Spec tests: portable hooks and Windows permission skips
**Target:** `src/internal/spec/implementation_test.go`, `src/internal/spec/delta_test.go`, `src/internal/spec/qa_test.go`
**Context:** [US-3]
**Parallel With:** T3.1, T3.2, T3.4

**Constraints:**
- Hook commands in tests are portable: the failure case uses `go no-such-command` instead of `false` (approved decision); hooks whose runs are counted through a log file (previously `#!/bin/sh` scripts) re-run the test binary on a helper test (`TestHookHelperProcess`), because `go version` cannot record a run. The helper fails with a clear message when the test binary path contains whitespace, since hook commands are split on spaces.
- File-mode assertions (`0600`) skip on Windows only, with the reason "permission bits do not apply on Windows".

**Acceptance Check:**
`go test -count=1 ./src/internal/spec/ && GOOS=windows go vet ./src/internal/spec/ && ! grep -nE '"(true|false)"' src/internal/spec/implementation_test.go && echo OK` prints `OK`.

- [x] T3.4: [CODE] Config security test uses an impossible root on every platform
**Target:** `src/internal/core/config_test.go`
**Context:** [US-3]
**Parallel With:** T3.1, T3.2, T3.3

**Constraints:**
- Replace `"/nonexistent/root"` with a root that cannot be created anywhere (for example a path below a regular file created in `t.TempDir()`); the test still expects an error.

**Acceptance Check:**
`go test -count=1 ./src/internal/core/ && GOOS=windows go vet ./src/internal/core/ && ! grep -q '/nonexistent/root' src/internal/core/config_test.go && echo OK` prints `OK`.

- [x] T3.5: [CODE] Archive instructions test drains captured output concurrently
**Target:** `src/internal/cli/archive_test.go`
**Context:** [US-3], [US-2]

**Constraints:**
- Found in CI run `44954510860`: once T2.1 loads the full protocol (~8 KB), the test deadlocks on Windows because it writes to an `os.Pipe` that is only read after the handler returns (the Windows pipe buffer is smaller than the output). Read the pipe in a goroutine started before calling the handler; keep both content assertions unchanged.

**Acceptance Check:**
`go test -count=1 -run TestHandleArchiveInstructions ./src/internal/cli/ && GOOS=windows go vet ./src/internal/cli/ && echo OK` prints `OK`; the Windows proof is S1.

### Phase 4: Constitution

- [x] T4.1: [DOCS] User-state isolation rule names both home variables
**Target:** `.specforce/docs/engineering.md`
**Context:** [US-4]

**Constraints:**
- Amend only the "User-State Isolation" bullet: require `HOME` and `USERPROFILE` pointed at the temporary directory, and allow skipping only permission-bit / read-only-directory assertions on Windows with a stated reason (PATCH amendment approved on 2026-11-21).

**Acceptance Check:**
`grep -A1 'User-State Isolation' .specforce/docs/engineering.md | grep -q 'USERPROFILE' && grep -q 'permission bits do not apply on Windows' .specforce/docs/engineering.md && echo OK` prints `OK`.

## 3. Black-box Scenarios

- **S1:** All CI jobs green on the version line
**Context:** [US-1], [US-2], [US-3]
**Command:** `gh run list --branch 1.x --workflow CI --limit 1 --json conclusion,headSha --jq '.[0].conclusion'`
**Expected:** `success` for the run of the last pushed commit (lint, security, unit tests on ubuntu/macos/windows, mutation quality, release tooling tests).

- **S2:** Archiving instructions carry the protocol
**Context:** [US-2]
**Command:** `HOME=$(mktemp -d) go run ./src/cmd/specforce archive instructions </dev/null | grep -c 'Execution Protocol'`
**Expected:** `1` or more (on Linux here; on Windows the same output is proven by `TestHandleArchiveInstructions` in S1).

## 4. Pre-emptive Mitigations
- **Risk:** Each Windows iteration costs a push and a CI run (about 10 minutes). -> **Mitigation:** cross-compile tests with `GOOS=windows go vet` before every push and batch all Phase 3 fixes into one push.
- **Risk:** A test that leaked into the real profile on a developer's Windows machine left a directory at `~/.specforce/state.json`. -> **Mitigation:** CI runners are ephemeral; the isolation fix prevents new leaks.
