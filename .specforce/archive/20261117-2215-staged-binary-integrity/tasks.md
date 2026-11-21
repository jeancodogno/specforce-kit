---
slug: 20261117-2215-staged-binary-integrity
lens: Backend-heavy
---

# Implementation Roadmap: Staged Binary Integrity

## 1. Execution Strategy
- **Gravity Order:** fingerprint recorded at staging and checked before the swap -> test run of the staged binary -> retry-once policy for verification failures -> rollback when the restart fails -> repair of unreadable state.
- Builds on the uncommitted implementation of `20261117-1840-autoupdate-reliability` in the working tree (target binding, no-downgrade guard, `ErrSwapSkipped`, `ManualUpdateVersion`, `AppliedVersion`, notices on stderr). Do not change those behaviors.
- Every phase opens with a `[TEST: RED]` task whose tests fail against the current code, followed by its `[CODE: GREEN]` task.

## 2. Tasks

### Phase 1: Fingerprint the Staged Binary

- [x] T1.1: [TEST: RED] Regression tests for corrupted and unfingerprinted staged files
**Target:** `src/internal/upgrade/service_test.go`
**Context:** [FIX-1]

**Constraints:**
- Cover: staging records the release checksum and size of the verified download; staged file emptied, truncated or overwritten after staging → swap does not modify the target, staged file deleted, update no longer ready; ready update with no fingerprint (simulates a state written by an earlier build) → discarded, target untouched; intact staged file → applied (stub the test run of Phase 2 as passing so this phase stays independent).
- Must fail before the fix (today an empty staged file replaces the target). Prefix test names with `TestIntegrity`.

**Acceptance Check:**
`go test ./src/internal/upgrade -run 'TestIntegrity' -v` fails, reporting at least the empty-staged-file case.

- [x] T1.2: [CODE: GREEN] Record and verify the fingerprint
**Target:** `src/internal/upgrade/state.go`, `src/internal/upgrade/service.go`, `src/internal/upgrade/installer_binary.go`
**Context:** [FIX-1]

**Constraints:**
- The fingerprint is the sha256 already verified against the release checksum file plus the byte size; the installer returns it with the downloaded path. Verify the staged file against it right after staging (a bad cross-filesystem copy is discarded immediately) and again before the swap.
- New state fields are optional so earlier state files still load; a ready update without a fingerprint is discarded (FIX-1 "Staged update without a fingerprint").

**Acceptance Check:**
`go test ./src/internal/upgrade -run 'TestIntegrity' -v` passes.

### Phase 2: Test Run of the Staged Binary

- [x] T2.1: [TEST: RED] Tests for the pre-swap test run
**Target:** `src/internal/upgrade/service_test.go`
**Context:** [FIX-1]

**Constraints:**
- Use an injectable runner for the test run. Cover: runner reports the staged version (with or without the leading `v`) → applied; runner error, timeout, or a different version → target untouched and update discarded; the runner receives the staged file path, an empty temporary HOME, and a timeout of a few seconds.
- One case builds a real shell script as staged file (`#!/bin/sh` echoing `specforce version <v>`) to prove the real runner path, skipped on Windows.

**Acceptance Check:**
`go test ./src/internal/upgrade -run 'TestStagedRun' -v` fails on the error, timeout and wrong-version cases.

- [x] T2.2: [CODE: GREEN] Run the staged binary before swapping
**Target:** `src/internal/upgrade/staged_check.go`, `src/internal/upgrade/service.go`
**Context:** [FIX-1]

**Constraints:**
- The test run executes the staged file with `--version` only, with HOME set to a fresh empty temporary directory, so any binary (including older release code) finds no state and cannot check or swap; remove that directory afterwards.
- The staged file is made executable for the owner only (0700) for the test run; it never becomes group/world readable before the swap.
- Compare by normalized version (leading `v` ignored) against the first line `specforce version <v>`.

**Acceptance Check:**
`go test ./src/internal/upgrade -run 'TestStagedRun|TestIntegrity' -v` passes.

### Phase 3: Retry Once, Then Manual

- [x] T3.1: [TEST: RED] Tests for the verification-failure retry policy
**Target:** `src/internal/upgrade/service_test.go`
**Context:** [FIX-2]

**Constraints:**
- Cover: first verification failure for version V → discarded, not manual, and the next background check downloads and stages V again (httptest request count > 0); second failure for V → recorded as manual, the manual notice for V is pending once, and a later background check makes no download request; a swap failure (move/rename error) stays terminal on the first occurrence, as before.

**Acceptance Check:**
`go test ./src/internal/upgrade -run 'TestVerifyRetry' -v` fails on the re-stage and second-failure cases.

- [x] T3.2: [CODE: GREEN] Count verification failures per version
**Target:** `src/internal/upgrade/state.go`, `src/internal/upgrade/service.go`
**Context:** [FIX-2]

**Constraints:**
- The counter belongs to one version; staging a different version resets it.

**Acceptance Check:**
`go test ./src/internal/upgrade -run 'TestVerifyRetry|TestIntegrity|TestStagedRun' -v` passes.

### Phase 4: Roll Back When the Restart Fails

- [x] T4.1: [TEST: RED] Tests for re-exec failure rollback
**Target:** `src/cmd/specforce/main_test.go`, `src/internal/upgrade/service_test.go`
**Context:** [FIX-3]

**Constraints:**
- Entrypoint: with a healthy staged update and an injected exec that returns an error, after the swap check the target file has the original content again, no `.old` remains, the applied version is not recorded and the version is recorded as manual. With an injected exec that succeeds, the applied version is recorded (existing behavior).
- Service: a rollback operation restores `<target>.old` over the target and reports an error when no backup exists.

**Acceptance Check:**
`go test ./src/internal/upgrade ./src/cmd/specforce -run 'Rollback' -v` fails.

- [x] T4.2: [CODE: GREEN] Restore the previous executable when the restart fails
**Target:** `src/internal/upgrade/service.go`, `src/cmd/specforce/main.go`
**Context:** [FIX-3]

**Constraints:**
- The `.old` backup must survive until the restart outcome is known; the existing `.old` cleanup (start of the swap and the CLI pre-run) runs only after a successful start, so do not move it before the exec.
- Windows keeps its spawn-and-exit mechanism: a spawn error is the restart failure.

**Acceptance Check:**
`go test ./src/internal/upgrade ./src/cmd/specforce -run 'Rollback' -v` passes.

### Phase 5: Repair Unreadable State

- [x] T5.1: [TEST: RED] Tests for unreadable state handling
**Target:** `src/internal/upgrade/state_test.go`, `src/internal/upgrade/service_test.go`
**Context:** [FIX-4]

**Constraints:**
- Cover: garbage, truncated and wrong-typed state files load as empty state without error; an orphan `staged/specforce_next` present with an unparseable state is deleted; the background check then writes a valid state (parses as JSON with `last_checked_at`); a directory at the state path → no panic, the updater does nothing and returns.

**Acceptance Check:**
`go test ./src/internal/upgrade -run 'TestStateRepair' -v` fails on the garbage-state and orphan cases.

- [x] T5.2: [CODE: GREEN] Treat unparseable state as empty and drop orphan staged files
**Target:** `src/internal/upgrade/state.go`, `src/internal/upgrade/service.go`
**Context:** [FIX-4]

**Constraints:**
- Only a parse failure is repaired; a read error other than "not found" (e.g. the path is a directory or unreadable) keeps returning an error so nothing is overwritten.

**Acceptance Check:**
`go test ./src/internal/upgrade -run 'TestStateRepair' -v` passes.

## 3. Black-box Scenarios

Run from the repository root with `M=github.com/jeancodogno/specforce-kit/src/internal/core.Version`. S4 needs network access and a published release newer than `1.0.0-beta.2`.

- **S1:** Empty staged file never replaces the executable
**Context:** [FIX-1]
**Command:** `bash -c 'T=$(mktemp -d); go build -ldflags="-X github.com/jeancodogno/specforce-kit/src/internal/core.Version=1.0.0-beta.3" -o $T/w/specforce ./src/cmd/specforce && mkdir -p $T/h/.specforce/upgrade/staged && : > $T/h/.specforce/upgrade/staged/specforce_next && printf "{\"latest_version\":\"v1.0.0-beta.9\",\"staged_version\":\"v1.0.0-beta.9\",\"update_ready\":true,\"target_path\":\"%s\"}" $T/w/specforce > $T/h/.specforce/state.json && HOME=$T/h timeout 15 $T/w/specforce --version && HOME=$T/h timeout 15 $T/w/specforce --version; ls $T/h/.specforce/upgrade/staged | wc -l'`
**Expected:** prints `specforce version 1.0.0-beta.3` twice, then `0`.

- **S2:** Staged file that is not a runnable Specforce is discarded
**Context:** [FIX-1]
**Command:** `bash -c 'T=$(mktemp -d); go build -ldflags="-X github.com/jeancodogno/specforce-kit/src/internal/core.Version=1.0.0-beta.3" -o $T/w/specforce ./src/cmd/specforce && mkdir -p $T/h/.specforce/upgrade/staged && head -c 2048 /dev/urandom > $T/h/.specforce/upgrade/staged/specforce_next && printf "{\"latest_version\":\"v1.0.0-beta.9\",\"staged_version\":\"v1.0.0-beta.9\",\"update_ready\":true,\"target_path\":\"%s\"}" $T/w/specforce > $T/h/.specforce/state.json && HOME=$T/h timeout 15 $T/w/specforce --version; HOME=$T/h timeout 15 $T/w/specforce --version'`
**Expected:** prints `specforce version 1.0.0-beta.3` twice; exit code 0.

- **S3:** Garbage state file is repaired
**Context:** [FIX-4]
**Command:** `bash -c 'T=$(mktemp -d); go build -o $T/specforce ./src/cmd/specforce && mkdir -p $T/h/.specforce && echo "not json" > $T/h/.specforce/state.json && HOME=$T/h HTTPS_PROXY=http://127.0.0.1:9 $T/specforce --internal-upgrade-check; python3 -c "import json,sys; d=json.load(open(sys.argv[1])); print(\"REPAIRED\" if d.get(\"last_checked_at\") else \"NO\")" $T/h/.specforce/state.json'`
**Expected:** prints `REPAIRED`.

- **S4:** A healthy update is still applied end to end
**Context:** [FIX-1], [FIX-3]
**Command:** `bash -c 'T=$(mktemp -d); go build -ldflags="-X github.com/jeancodogno/specforce-kit/src/internal/core.Version=1.0.0-beta.2" -o $T/i/specforce ./src/cmd/specforce && export HOME=$T/h && $T/i/specforce --internal-upgrade-check && $T/i/specforce --version && $T/i/specforce --version'`
**Expected:** prints `specforce version` followed by the newer published version (not `1.0.0-beta.2`) on both lines.

## 4. Pre-emptive Mitigations
- **Risk:** The `on_task_finished` hook runs `make test`, so a `[TEST: RED]` task cannot be finished alone. -> **Mitigation:** finish each RED task together with its GREEN partner in one status update, after the GREEN acceptance check passes.
- **Risk:** Tests or manual runs against the real home directory consume the developer's staged update. -> **Mitigation:** every test sets `HOME` to `t.TempDir()` before creating the state manager; never run a built binary with the real HOME.
- **Risk:** The test run of a staged binary executes downloaded code. -> **Mitigation:** it runs only after the release checksum and fingerprint matched, with `--version` only, an empty temporary HOME and a timeout; this is the same code the swap would run seconds later.
