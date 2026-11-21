---
slug: 20261117-2215-staged-binary-integrity
lens: Bugfix
---

# Bugfix: Staged Binary Integrity

## 1. Issue Description
Follow-up to `20261117-1840-autoupdate-reliability` (implemented, QA PASSED, findings recorded). The auto-update applies the staged binary without checking it again: if the staged file is corrupted after it was verified at download time (interrupted cross-filesystem copy, full disk, crash, tampering), it replaces the installed executable and leaves Specforce unusable, with only a hidden `.old` file as a manual way back. Separately, an unreadable state file silently disables the auto-update forever.

Impact: one corrupted staged file turns every later `specforce` invocation into a no-op or a shell error, including the commands AI agents rely on. The rule "on any doubt the update is discarded, never applied" from the previous spec is only met at download time.

## 2. Evidence & Observations
Observed on 2026-11-17 with local builds and an isolated home directory (independent QA, re-verified by the orchestrator).
- **Symptom (empty staged file):** installed `1.0.0-beta.3`, a 0-byte `~/.specforce/upgrade/staged/specforce_next`, state with a newer staged version and the correct target path. `specforce --version` printed `specforce version 1.0.0-beta.3` and exited 0, and the installed executable was then 0 bytes with mode 755. Every later run silently does nothing (the shell runs an empty file as an empty script).
- **Symptom (random staged file):** a 2 KB random file was applied the same way; later runs fail with shell errors.
- **Symptom (false success):** in both cases the state recorded the version as applied although running the new binary failed.
- **Symptom (unreadable state):** with a garbage, truncated, wrong-typed state file, or a directory in its place, commands work but no update check ever runs again and the file is never repaired.
- **Trace:** `src/internal/upgrade/service.go` `swapBinary` (only checks that the staged file exists), `PerformAtomicSwapAt` (records the applied version as soon as the move succeeds), `stageUpdate` (no fingerprint recorded); `src/cmd/specforce/main.go` `checkForSwap` (ignores the error returned by the re-exec: "If it fails to restart, we just continue with the old binary in memory"); `src/internal/upgrade/state.go` `Load` (returns the parse error, so callers return early).

## 3. Reproduction Steps
All steps use a temporary directory `$T` and `HOME=$T/h`.
1. Build the current tree with version `1.0.0-beta.3` into `$T/w/specforce`: `go build -ldflags="-X 'github.com/jeancodogno/specforce-kit/src/internal/core.Version=1.0.0-beta.3'" -o $T/w/specforce ./src/cmd/specforce`.
2. Create an empty `$T/h/.specforce/upgrade/staged/specforce_next` and write `$T/h/.specforce/state.json` with `latest_version` and `staged_version` `v1.0.0-beta.9`, `update_ready: true`, `target_path` set to `$T/w/specforce`.
3. Run `HOME=$T/h $T/w/specforce --version`, then `stat -c %s $T/w/specforce`. Expected: executable untouched (non-zero size), staged file discarded. Actual: size `0`.
4. Unreadable state: write `not json` into `$T/h/.specforce/state.json` and run `HOME=$T/h DEBUG=1 script -qec "$T/w/specforce completion bash" /dev/null` twice. Expected: the state file is replaced by a valid one and update checks resume. Actual: the file stays `not json` and no check is ever started.

## 4. Root Cause Analysis (RCA)
- **No fingerprint survives staging.** The checksum is verified on the temporary download, but nothing about the verified file (hash, size) is recorded, so the swap cannot tell whether the staged file is still the verified one.
- **No proof the staged binary runs.** The swap never executes the staged binary before replacing the installed one, so a file that is intact but not runnable on this machine would also be applied.
- **Success is declared before the new binary runs.** The applied version is recorded as soon as the file move succeeds, and a failed re-exec is ignored by the entrypoint, so the broken executable stays in place and the user is told it was updated.
- **State parse errors are fatal to the updater but never repaired.** `Load` returns the parse error; every caller treats it as "do nothing", and nothing rewrites the file.

## 5. Acceptance Criteria (Regression Tests)

### [FIX-1] Staged binary is verified before it can replace the executable
**Scenario: [Regression] Corrupted staged file is never applied**
GIVEN a staged update for the running executable whose staged file no longer matches the fingerprint recorded when it was staged (different content, size, or empty)
WHEN any command starts
THEN the installed executable is not modified
AND the staged file is deleted and the update is no longer marked as ready
AND the command runs normally on the current version.

**Scenario: [Regression] Staged binary that does not run is never applied**
GIVEN a staged update whose file matches its fingerprint
AND running it to report its version fails, times out, or reports a version different from the staged version
WHEN any command starts
THEN the installed executable is not modified and the update is discarded as in the previous scenario.

**Scenario: Staged update without a fingerprint**
GIVEN a staged update recorded without a fingerprint (staged by an earlier build)
WHEN any command starts
THEN it is discarded without being applied.

**Scenario: Healthy staged binary is applied**
GIVEN a staged update that matches its fingerprint and reports the staged version when run
WHEN any command starts
THEN the executable is replaced and the command runs on the new version.

### [FIX-2] A failed verification is retried once, then becomes a manual update
**Scenario: First verification failure re-stages**
GIVEN a version whose staged binary failed verification for the first time
WHEN the next background check runs
THEN that version is downloaded, verified and staged again.

**Scenario: Second verification failure for the same version**
GIVEN the same version fails verification a second time
WHEN the update is discarded
THEN the version is recorded as requiring a manual update
AND the manual-update hint is shown once on the next interactive, non-agent command
AND the version is not downloaded again.

### [FIX-3] A failed restart restores the previous executable
**Scenario: [Regression] Re-exec failure rolls back**
GIVEN the executable was replaced by the staged binary
AND starting the new executable fails
WHEN the command continues on the old binary in memory
THEN the previous executable is restored to its path
AND the version is not recorded as applied
AND the version is recorded as requiring a manual update (hint shown once, as in FIX-2).

**Scenario: Successful restart keeps the new executable**
GIVEN the executable was replaced and the new executable starts
THEN the version is recorded as applied and the "updated" notice behaves as before.

### [FIX-4] An unreadable state file is repaired instead of disabling updates
**Scenario: [Regression] Garbage state file**
GIVEN the state file exists but cannot be parsed (garbage, truncated, wrong types)
WHEN any command starts or the background check runs
THEN the state is treated as empty
AND any staged binary is deleted (it has no trusted version, target or fingerprint)
AND the next background check runs and writes a valid state file
AND no command fails and no notice is shown because of it.

**Scenario: State path is not a regular file**
GIVEN the state path is a directory or otherwise unreadable
WHEN any command starts
THEN the command runs normally and the updater does nothing, without crashing or looping.

## 6. Technical Constraints (NFR)
- **[Safety]:** The installed executable is replaced only after the staged file matches its recorded fingerprint AND a test run of the staged file reports the staged version. Any doubt discards the update.
- **[Performance]:** The test run happens only when an update is about to be applied (at most once per version per executable), is bounded by a timeout of a few seconds, and must not trigger another update check or swap in the staged binary.
- **[Isolation]:** The test run of the staged binary must not read or modify the user's projects; it only asks the binary for its version.
- **[Observability]:** With `DEBUG=1`, each verification failure, rollback and state repair is written to the error stream with its reason.
- **[Compatibility]:** Builds the previous spec's behavior (target binding, no downgrade, notices on stderr only, once-per-version manual hint, npm flow) unchanged. Platform-independent logic; the Windows exec mechanism is kept and covered only through shared logic.
- **[Security]:** The fingerprint lives in the owner-only state file; the staged file stays owner-only until applied.
