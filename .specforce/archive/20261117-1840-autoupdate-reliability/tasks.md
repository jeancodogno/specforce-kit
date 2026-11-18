---
slug: 20261117-1840-autoupdate-reliability
lens: Backend-heavy
---

# Implementation Roadmap: Auto-update Reliability

## 1. Execution Strategy
- **Gravity Order:** Apply-time guards on the swap (state + service) -> background check robustness (throttle, writability) -> notice selection and output channel (service, TUI, CLI entrypoints) -> npm channel-aware provider and exact-version install -> binary download hygiene -> removal of unused update features.
- Every phase opens with a `[TEST: RED]` task whose tests fail against the current code, followed by the `[CODE: GREEN]` task that makes them pass. Tests use a temporary home directory and a temporary install directory; they never touch the real `~/.specforce` (see Pre-emptive Mitigations).

## 2. Tasks

### Phase 1: Apply-time Guards (swap safety)

- [x] T1.1: [TEST: RED] Regression tests for swap decisions
**Target:** `src/internal/upgrade/service_test.go`
**Context:** [FIX-2], [FIX-3], [FIX-4]

**Constraints:**
- Cover: staged version older than running → not applied, staged file deleted, update no longer ready; staged version equal → same; staged newer for this exact path in a writable dir → applied; update staged for another executable path → this executable is untouched and the staged update stays intact; swap failure (target dir made read-only) → running binary unchanged, update no longer ready, staged file deleted, version recorded as requiring manual update, and a second call does not attempt the swap.
- These tests MUST fail before the fix (today an older staged version is applied, other paths consume the update and a failed swap leaves the update ready).
- Drive the swap through the existing target-path entry point so no real executable is replaced.

**Acceptance Check:**
`go test ./src/internal/upgrade -run 'TestSwap' -v` fails, reporting at least the older-version, other-target and failed-swap cases.

- [x] T1.2: [CODE: GREEN] Bind staged updates to their target and guard the swap
**Target:** `src/internal/upgrade/state.go`, `src/internal/upgrade/service.go`
**Context:** [FIX-2], [FIX-3], [FIX-4]

**Constraints:**
- The state records the resolved executable path the update was staged for, and the version that requires a manual update; new fields are optional so state files from earlier versions still load.
- Apply only when all hold: same target path, staged version strictly newer than the running version, update marked ready. Older/equal staged versions are discarded (file deleted, ready cleared). Another target path is left untouched.
- Any swap error is terminal for that version: roll back the binary as today, then clear ready, delete the staged file, record the manual-update version, and persist the state.
- With `DEBUG=1`, each skip/failure reason is written to the error stream.

**Acceptance Check:**
`go test ./src/internal/upgrade -run 'TestSwap' -v` passes.

### Phase 2: Background Check Robustness

- [x] T2.1: [TEST: RED] Tests for throttle on failure and non-writable install dir
**Target:** `src/internal/upgrade/service_test.go`
**Context:** [FIX-6], [FIX-1]

**Constraints:**
- Cover: provider returns an error → check time is persisted and the 6-hour throttle then skips the next check; newer version available but install dir not writable → nothing downloaded or staged and the version is recorded as requiring manual update; writable dir → staging proceeds as today (use the existing mock provider and an `httptest` server for the download).
- The executable path used by the background check must be injectable for tests (same approach as the existing injectable clock).

**Acceptance Check:**
`go test ./src/internal/upgrade -run 'TestBackgroundUpdate' -v` fails on the provider-error and non-writable cases.

- [x] T2.2: [CODE: GREEN] Persist check time on failure and pre-flight writability
**Target:** `src/internal/upgrade/service.go`
**Context:** [FIX-6], [FIX-1]

**Constraints:**
- Writability is tested by creating and removing a temporary file in the executable's directory, not by inspecting permission bits (ownership and ACLs make bits unreliable).
- When staging for the binary path, the recorded target path is the executable path resolved at check time (T1.2 field).

**Acceptance Check:**
`go test ./src/internal/upgrade -run 'TestBackgroundUpdate' -v` passes.

### Phase 3: Notices and Output Channel

- [x] T3.1: [TEST: RED] Tests for notice selection, once-per-version hint and clean stdout
**Target:** `src/internal/upgrade/service_test.go`, `src/internal/tui/upgrade_notification_test.go`, `src/cmd/specforce/main_test.go`
**Context:** [FIX-5], [FIX-1]

**Constraints:**
- Service: after a successful swap, an "applied" notice for that version is pending exactly once; a manual-update version produces a "manual" notice exactly once per version (marking it shown persists it); an available-and-stageable version produces the "available" notice; nothing is pending otherwise.
- TUI: the manual notice names the version and the exact manual command (reinstall for binary installs, `npm install -g` of the exact version for npm installs) and never contains "applied automatically"; the available notice no longer promises automatic application when the version is recorded as manual.
- Entrypoint: the pre-command swap writes nothing to standard output (capture stdout around the swap check with a staged newer binary in a temp install dir).

**Acceptance Check:**
`go test ./src/internal/upgrade ./src/internal/tui ./src/cmd/specforce -run 'Notice|Notification|Swap' -v` fails on the once-per-version, manual-text and stdout cases.

- [x] T3.2: [CODE: GREEN] Route all updater notices through the interactive post-run on stderr
**Target:** `src/internal/upgrade/service.go`, `src/internal/tui/upgrade_notification.go`, `src/cmd/specforce/main.go`
**Context:** [FIX-5], [FIX-1]

**Constraints:**
- The pre-command swap in `main.go` stays in place (it runs for every command, agent or not, per the confirmed decision) but prints nothing; it only records the applied version in state.
- Notice display is decided after command parsing, where TTY and the agent-command annotation are known (T3.3).

**Acceptance Check:**
`go test ./src/internal/upgrade ./src/internal/tui ./src/cmd/specforce -run 'Notice|Notification|Swap' -v` passes.

- [x] T3.3: [CODE: GREEN] Show pending notice on stderr only for interactive non-agent commands
**Target:** `src/internal/cli/cobra/root.go`
**Context:** [FIX-5], [FIX-1]

**Constraints:**
- Write to the command's error stream, never stdout; no notice at all when stdout is not a terminal or the command (or a parent) is an agent command.
- Mark the manual notice as shown only after it is actually written.

**Acceptance Check:**
`go test ./src/internal/cli/cobra -run 'TestRootUpdateNotice' -v` passes, covering: agent command → no output on either stream; non-TTY → no output; interactive non-agent → notice on the error stream and empty stdout addition. (The test is added in this task, using an injectable TTY check.)

### Phase 4: npm Channel-aware Updates

- [x] T4.1: [TEST: RED] Tests for npm detection, channel tag and exact-version install
**Target:** `src/internal/upgrade/provider_test.go`, `src/internal/upgrade/installer_npm_test.go`, `src/internal/upgrade/service_test.go`
**Context:** [FIX-7]

**Constraints:**
- Detection: a path inside the Specforce npm platform package (`node_modules/@jeancodogno/specforce-kit-<os>-<arch>/`) is npm; a path merely containing "node" or "npm" elsewhere is binary.
- Channel: `1.0.0-beta.2` → `beta`, `1.0.0-alpha.3` → `alpha`, `1.0.0` → `latest`; registry served by `httptest` with a dist-tags document; missing channel tag → "no update", check time persisted, nothing recorded as manual.
- Installer: invoked with the package name suffixed by `@<exact version>` (assert with the existing command-executor mock).
- npm install error → version recorded as requiring manual update, not retried on the next check.

**Acceptance Check:**
`go test ./src/internal/upgrade -run 'NPM|Channel|Detect' -v` fails on the detection, channel and exact-version cases.

- [x] T4.2: [CODE: GREEN] Channel-aware npm provider, exact install and specific detection
**Target:** `src/internal/upgrade/provider_npm.go`, `src/internal/upgrade/installer_npm.go`, `src/internal/upgrade/service.go`
**Context:** [FIX-7]

**Constraints:**
- The npm provider reads the registry dist-tags of `@jeancodogno/specforce-kit` and returns the version for the running version's channel; a missing tag is a distinct "no update" outcome, not an error that could be confused with network failure.
- npm install output is captured and, with `DEBUG=1`, written to the error stream on failure.

**Acceptance Check:**
`go test ./src/internal/upgrade -run 'NPM|Channel|Detect' -v` passes.

- [x] T4.3: [CODE: GREEN] Select provider by installation method in both entrypoints
**Target:** `src/internal/cli/cobra/root.go`, `src/cmd/specforce/main.go`
**Context:** [FIX-7]

**Constraints:**
- Replace both hard-coded GitHub provider choices (the two open `TODO`s) with one shared selection exposed by the upgrade package, so the foreground and background processes always agree.

**Acceptance Check:**
`go test ./src/internal/upgrade -run 'TestSelectProvider' -v` passes (test added in this task: npm package path → npm provider; any other path → GitHub provider), and `grep -c "TODO" src/cmd/specforce/main.go src/internal/cli/cobra/root.go` reports 0 for both files.

### Phase 5: Binary Download Hygiene

- [x] T5.1: [TEST: RED] Tests for temp cleanup, bounded checksum fetch and private staged file
**Target:** `src/internal/upgrade/installer_binary_test.go`, `src/internal/upgrade/service_test.go`
**Context:** [FIX-8]

**Constraints:**
- Point the temp directory at a test directory (`TMPDIR`) and assert it is empty after: binary 404, checksum 404, checksum entry missing, checksum mismatch.
- Checksum server that never responds + a context with a short deadline → the call returns before a few seconds.
- After staging, the staged file mode is `0600`.

**Acceptance Check:**
`go test ./src/internal/upgrade -run 'Hygiene|TempCleanup|ChecksumTimeout|StagedPerm' -v` fails on cleanup, timeout and permission cases.

- [x] T5.2: [CODE: GREEN] Clean temp files on every failure, use context for checksum, stage as 0600
**Target:** `src/internal/upgrade/installer_binary.go`, `src/internal/upgrade/service.go`
**Context:** [FIX-8]

**Constraints:**
- The executable bit is still set only at swap time (existing behavior), so the staged file stays owner-only until applied.

**Acceptance Check:**
`go test ./src/internal/upgrade -run 'Hygiene|TempCleanup|ChecksumTimeout|StagedPerm' -v` passes.

### Phase 6: Remove Unused Update Features

- [x] T6.1: [TEST: RED] Test that legacy state with an ignored version has no effect
**Target:** `src/internal/upgrade/state_test.go`
**Context:** [FIX-9]

**Constraints:**
- Load a state file containing `ignored_version` equal to the latest version: it loads without error, and the service still reports the update as available and stages it.
- This fails before the fix because the ignored version currently suppresses the update.

**Acceptance Check:**
`go test ./src/internal/upgrade -run 'TestLegacyState' -v` fails on the "still available" assertion.

- [x] T6.2: [CODE: GREEN] Remove the ignored-version feature and the unused upgrade path
**Target:** `src/internal/upgrade/state.go`, `src/internal/upgrade/service.go`, `src/internal/upgrade/service_test.go`
**Context:** [FIX-9]

**Constraints:**
- Remove the field, every read of it, the unused direct-upgrade method, and tests that exist only for them; unknown JSON keys are ignored on load, so no migration is needed.

**Acceptance Check:**
`go test ./src/internal/upgrade -run 'TestLegacyState' -v` passes and `grep -rn "IgnoredVersion\|PerformUpgrade" src/` prints nothing.

## 3. Black-box Scenarios

Run from the repository root. S1, S2 and S4 need network access and a published release newer than `1.0.0-beta.2`; S3 needs `script` and `strace`.

- **S1:** Non-writable install shows a single honest hint and never loops
**Context:** [FIX-1], [FIX-2], [FIX-5]
**Command:** `bash -c 'T=$(mktemp -d); go build -ldflags="-X github.com/jeancodogno/specforce-kit/src/internal/core.Version=1.0.0-beta.2" -o $T/i/specforce ./src/cmd/specforce && chmod 555 $T/i && export HOME=$T/h && $T/i/specforce --internal-upgrade-check; ls -A $T/h/.specforce/upgrade/staged 2>/dev/null | wc -l; script -qec "$T/i/specforce completion bash" /dev/null | grep -c "1.0.0-beta"; script -qec "$T/i/specforce completion bash" /dev/null | grep -c "1.0.0-beta"; chmod 755 $T/i'`
**Expected:** prints `0` (nothing staged), then `1` or more (hint shown on the first interactive run, naming the version), then `0` (not shown again); no output line contains "applied automatically".

- **S2:** Agent JSON stays valid while an update is applied
**Context:** [FIX-5], [FIX-3]
**Command:** `bash -c 'T=$(mktemp -d); go build -ldflags="-X github.com/jeancodogno/specforce-kit/src/internal/core.Version=1.0.0-beta.2" -o $T/i/specforce ./src/cmd/specforce && export HOME=$T/h && $T/i/specforce --internal-upgrade-check && $T/i/specforce spec list --json 2>/dev/null | python3 -m json.tool >/dev/null && echo JSON_OK && $T/i/specforce --version'`
**Expected:** prints `JSON_OK`, then `specforce version` followed by the newer published version (not `1.0.0-beta.2`).

- **S3:** Offline runs respect the 6-hour throttle
**Context:** [FIX-6]
**Command:** `bash -c 'T=$(mktemp -d); go build -o $T/specforce ./src/cmd/specforce && for i in 1 2 3; do HOME=$T/h HTTPS_PROXY=http://127.0.0.1:9 script -qec "strace -f -e trace=execve -o $T/tr.$i $T/specforce completion bash" /dev/null >/dev/null; sleep 12; grep -c internal-upgrade-check $T/tr.$i; done'`
**Expected:** prints `1`, `0`, `0`.

- **S4:** npm installs install the exact version of their own channel
**Context:** [FIX-7]
**Command:** `bash -c 'T=$(mktemp -d); P=$T/p/lib/node_modules/@jeancodogno/specforce-kit-linux-x64/bin; mkdir -p $P $T/f; go build -ldflags="-X github.com/jeancodogno/specforce-kit/src/internal/core.Version=1.0.0-beta.2" -o $P/specforce ./src/cmd/specforce && printf "#!/bin/sh\necho \"\$*\" >> $T/npm.log\n" > $T/f/npm && chmod +x $T/f/npm && HOME=$T/h PATH=$T/f:$PATH $P/specforce --internal-upgrade-check; cat $T/npm.log'`
**Expected:** prints `install -g @jeancodogno/specforce-kit@<version>` where `<version>` equals `npm view @jeancodogno/specforce-kit dist-tags.beta` (never the bare package name, never the `latest` tag version).

## 4. Pre-emptive Mitigations
- **Risk:** The `on_task_finished` hook runs `make test`, so a `[TEST: RED]` task cannot be finished alone while its tests fail. -> **Mitigation:** finish each RED task together with its GREEN partner in one batch status update, after the GREEN acceptance check passes.
- **Risk:** Tests or manual runs that use the real home directory consume or overwrite the developer's real staged update (observed during discovery). -> **Mitigation:** every test sets `HOME` to `t.TempDir()` (via `t.Setenv`) before creating the state manager; never run a locally built binary against the real home while an update is staged.
- **Risk:** Read-only directories created by tests block `t.TempDir()` cleanup. -> **Mitigation:** restore write permission in a `t.Cleanup` registered right after making a directory read-only; skip those cases when running as root (root ignores directory permissions).
