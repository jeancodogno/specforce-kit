---
slug: 20261117-1840-autoupdate-reliability
lens: Bugfix
---

# Bugfix: Auto-update Reliability

## 1. Issue Description
The background auto-update detects and stages new releases correctly, but applying them is broken or unsafe:
- Users whose binary lives in a directory they cannot write to (the default `make install` location, owned by root) are never updated, are told the update "will be applied automatically", and the CLI retries the failing swap on every invocation forever.
- npm users would be moved to the stable `latest` npm tag (`0.2.2`) instead of the next beta: a silent downgrade.
- A staged version older than the running one is applied (downgrade).
- The post-swap notice is written to standard output, corrupting machine-readable output consumed by AI agents.
- When the version check fails (offline), the 6-hour throttle never takes effect and every interactive command spawns a new background check.
- Update state is shared by every executable of the user, so any other build (e.g. a development binary) consumes the staged update and overwrites itself with the release binary.

This regresses rules already approved in the archived spec `20260506-1701-automated-background-updates` (REQ-1 "the check must respect the installation method"; REQ-3 "lacking write permission notifies the user that a manual update is required").

## 2. Evidence & Observations
All observations were reproduced on 2026-11-17 with two local builds (versions `1.0.0-beta.2` and `1.0.0-beta.3`) and an isolated home directory, except where stated.
- **Symptom (real machine):** `~/.specforce/state.json` stuck at `update_ready: true`, `staged_version: v1.0.0-beta.3`; tracing `specforce --version` shows `renameat("/usr/local/bin/specforce", "/usr/local/bin/specforce.old") = -1 EACCES`, then the command continues on the old binary with no message.
- **Symptom (non-writable dir):** 3 interactive runs with the install directory set read-only: version stays `1.0.0-beta.2`, `update_ready` stays `true`, and the box "It will be downloaded in the background and applied automatically" is printed every run.
- **Symptom (npm):** with the binary placed under `.../node_modules/@jeancodogno/specforce-kit-linux-x64/bin/` and a recording fake `npm`, the background check ran `npm install -g @jeancodogno/specforce-kit` (no version). Registry dist-tags: `latest=0.2.2`, `beta=1.0.0-beta.3`, `alpha=1.0.0-alpha.3`. The next run printed `(Specforce updated to v1.0.0-beta.3)` while `--version` still reported `1.0.0-beta.2`.
- **Symptom (downgrade):** installed `1.0.0-beta.3`, staged `1.0.0-beta.2` → after one run `--version` reports `1.0.0-beta.2`.
- **Symptom (agent output):** `specforce spec list --json 2>/dev/null` during a pending swap → stdout starts with `\n\x1b[38;5;245m(Specforce updated to v1.0.0-beta.2)\x1b[0m`; parsing it as JSON fails with `JSONDecodeError`.
- **Symptom (throttle):** 3 interactive runs, background checks spawned: online `1, 0, 0`; offline (unreachable proxy) `1, 1, 1`, and no state file is ever written.
- **Symptom (cross-executable):** a locally built binary run without an isolated home consumed the staged update meant for `/usr/local/bin/specforce` and replaced itself; its sha256 `33e26bc1…` equals the official `specforce-kit_linux_amd64` release asset.
- **Symptom (temp leak, code reading + artifact):** a 0-byte `/tmp/specforce-update-2412556037` left behind by a failed download.
- **Trace:** `src/cmd/specforce/main.go:38-69` (`checkForSwap`), `src/internal/upgrade/service.go:83-111` (`PerformBackgroundUpdate`), `service.go:144-200` (`PerformAtomicSwapAt`), `service.go:285-294` (`isNPM`), `src/internal/upgrade/installer_npm.go:36`, `src/internal/upgrade/installer_binary.go:29-86`, `src/internal/cli/cobra/root.go:28-86`, `src/internal/tui/upgrade_notification.go:31`.

## 3. Reproduction Steps
All steps use a temporary directory `$T` and `HOME=$T/home` so the real user state is never touched.
1. Build two versions from the repository root: `go build -ldflags="-X 'github.com/jeancodogno/specforce-kit/src/internal/core.Version=1.0.0-beta.2'" -o $T/sf-b2 src/cmd/specforce/main.go`, and the same with `1.0.0-beta.3` into `$T/sf-b3`.
2. **Non-writable install:** copy `$T/sf-b2` to `$T/install/specforce`; copy `$T/sf-b3` to `$T/home/.specforce/upgrade/staged/specforce_next`; write `$T/home/.specforce/state.json` with `update_ready: true`, `staged_version` and `latest_version` `v1.0.0-beta.3`; run `chmod 555 $T/install`; run `HOME=$T/home script -qec "$T/install/specforce completion bash" /dev/null` three times. Expected: no retry loop and an honest manual-update hint. Actual: version unchanged, `update_ready` stays `true`, "applied automatically" box every run.
3. **npm path:** copy `$T/sf-b2` to `$T/prefix/lib/node_modules/@jeancodogno/specforce-kit-linux-x64/bin/specforce`; put a script named `npm` that appends its arguments to a log first on `PATH`; run that binary with `--internal-upgrade-check`. Expected: install of `@jeancodogno/specforce-kit@1.0.0-beta.3`. Actual: `npm install -g @jeancodogno/specforce-kit`.
4. **Downgrade + agent output:** copy `$T/sf-b3` to a writable `$T/w/specforce`; stage `$T/sf-b2` with `staged_version: v1.0.0-beta.2`, `update_ready: true`; run `HOME=$T/home $T/w/specforce spec list --json 2>/dev/null` and parse stdout as JSON. Expected: valid JSON and version stays `1.0.0-beta.3`. Actual: JSON parse error and version becomes `1.0.0-beta.2`.
5. **Throttle offline:** with an empty `$T/home`, run three times `HOME=$T/home HTTPS_PROXY=http://127.0.0.1:9 script -qec "strace -f -e trace=execve -o $T/trace.N $T/install/specforce completion bash" /dev/null` and count `internal-upgrade-check` in each trace. Expected: `1, 0, 0`. Actual: `1, 1, 1`.
6. **Cross-executable:** with a valid staged update for executable A, run a different copy B with the same `HOME`. Expected: B ignores the update. Actual: B replaces itself with the staged binary.

## 4. Root Cause Analysis (RCA)
- **Swap errors are swallowed and never resolved.** `checkForSwap` ignores the error from the swap and leaves the "update ready" flag set, so the same failing swap is retried on every invocation. Nothing checks, before downloading, whether the install directory is writable, and the notification text is static ("applied automatically").
- **Installation-method detection and source are wrong for npm.** npm installs are detected by any occurrence of `node`/`npm` in the executable path, the version source is hard-coded to GitHub releases (two open `TODO`s), and the npm installer installs the package without a version, so npm resolves its `latest` dist-tag, which is the old stable line.
- **No version guard at apply time.** The swap compares neither the staged version with the running version, so an older or equal staged binary is applied.
- **Output channel ignores the execution context.** The swap notice is printed to standard output from `main`, before command parsing, so neither the TTY check nor the agent-command annotation applies to it.
- **Throttle timestamp is only saved on success.** When fetching the latest version fails, the function returns before persisting the check time, so the 6-hour window never starts.
- **State is not bound to an executable.** State and staged binary live per user, and the swap targets whichever executable is currently running.
- **Download hygiene gaps.** The temporary download is removed only on checksum mismatch; the checksum request has no context or timeout; the staged file is created with group/world-readable permissions, contrary to the security rule that upgrade-internal files use `0600`.
- **Dead code.** The "ignored version" field is never written by any command, and `PerformUpgrade` is unused.

## 5. Acceptance Criteria (Regression Tests)

### [FIX-1] Non-writable install location is handled honestly
**Scenario: [Regression] No download, no retry loop**
GIVEN the running executable's directory is not writable by the current user
AND a newer version is available
WHEN the background check runs
THEN no binary is downloaded or staged
AND the state records that this version requires a manual update.

**Scenario: [Regression] Manual hint shown once per version**
GIVEN a newer version was recorded as requiring a manual update
WHEN the user runs an interactive, non-agent command for the first time after that
THEN a single hint naming the new version and the manual update command is shown on the error stream
AND subsequent commands do not show the hint again for that same version
AND the message never claims the update will be applied automatically.

### [FIX-2] Swap failure ends the attempt for that version
**Scenario: [Regression] Failed swap clears pending state**
GIVEN an update is staged for the running executable
AND replacing the executable fails
WHEN any command starts
THEN the running version is unchanged and the command proceeds normally
AND the update is no longer marked as ready
AND the staged binary is deleted
AND the version is recorded as requiring a manual update (hint behavior of FIX-1 applies)
AND the next command does not attempt the swap again.

### [FIX-3] Never downgrade or reapply the same version
**Scenario: [Regression] Older or equal staged version is discarded**
GIVEN the staged version is older than or equal to the running version
WHEN any command starts
THEN the executable is not replaced
AND the staged binary is deleted and the update is no longer marked as ready.

**Scenario: Newer staged version is applied**
GIVEN the staged version is newer than the running version and the install directory is writable
WHEN any command starts
THEN the executable is replaced by the staged version and the command runs on the new version.

### [FIX-4] Staged updates belong to one executable
**Scenario: [Regression] Other executables ignore the staged update**
GIVEN an update was staged for executable A
WHEN executable B (a different path) runs any command
THEN B is not replaced and the staged update for A remains intact.

**Scenario: Target no longer exists**
GIVEN an update was staged for an executable path that no longer exists
WHEN any other executable starts
THEN the staged update for that path is left untouched by this executable.

### [FIX-5] Machine-readable output stays clean; notices only in interactive sessions
**Scenario: [Regression] Agent JSON output during a swap**
GIVEN a newer update is staged for the running executable
WHEN an agent command is run with JSON output and its standard output is captured
THEN the update is applied
AND standard output contains only the command's own output and parses as valid JSON.

**Scenario: Notices are interactive-only**
GIVEN any update notice (applied, available, or manual-update hint) is pending
WHEN a command runs without a terminal on standard output, or the command is an agent command
THEN no notice is written to any stream
AND WHEN an interactive, non-agent command runs
THEN the notice is written to the error stream only.

### [FIX-6] Throttle holds when the version check fails
**Scenario: [Regression] Offline runs do not spawn repeated checks**
GIVEN the remote version source is unreachable
WHEN the background check runs and fails
THEN the check time is still persisted
AND interactive commands within the next 6 hours do not start a new background check.

### [FIX-7] npm installations follow their own release channel
**Scenario: [Regression] Exact version from the matching channel**
GIVEN Specforce was installed through the npm package
AND the running version is a pre-release of channel C (e.g. `1.0.0-beta.2` → `beta`), or a stable release (channel `latest`)
WHEN the background check runs
THEN the latest version is read from the npm registry tag of channel C
AND when it is newer, exactly that version of the package is installed globally.

**Scenario: Channel tag missing**
GIVEN the npm registry has no tag for the running version's channel
WHEN the background check runs
THEN no update is performed, the check time is persisted, and nothing is shown to the user.

**Scenario: npm install fails**
GIVEN the global npm install of the new version fails (e.g. no permission on the global prefix)
WHEN the background check finishes
THEN the version is recorded as requiring a manual update (hint behavior of FIX-1 applies, with the npm install command for that exact version)
AND the install is not retried for that version.

**Scenario: Detection is specific to the Specforce npm package**
GIVEN the executable path merely contains the words "node" or "npm" but is not inside the Specforce npm platform package
WHEN the installation method is determined
THEN the installation is treated as a native binary installation.

### [FIX-8] Download hygiene
**Scenario: [Regression] No temporary files left on failure**
GIVEN the binary or checksum download fails at any step
WHEN the background check finishes
THEN no temporary download file remains on disk.

**Scenario: Bounded download time**
GIVEN the checksum server never answers
WHEN the background check runs
THEN the check ends within the background time budget (1 minute).

**Scenario: Private staged binary**
GIVEN a binary update is staged
THEN the staged file is readable and writable only by its owner until it is applied.

### [FIX-9] Unused update features removed without breaking existing state
**Scenario: Legacy state file still loads**
GIVEN an existing state file containing an "ignored version" entry
WHEN any command starts
THEN the state loads without error and the entry has no effect on update decisions.

## 6. Technical Constraints (NFR)
- **[Performance]:** No new work on the foreground path beyond reading the state file and, when a swap is pending, comparing versions and paths; the version check stays in the detached background process.
- **[Safety]:** The running executable is only replaced when all hold: the update targets this exact executable path, the staged version is strictly newer, and the staged binary's checksum was verified at staging time. On any doubt the update is discarded, never applied.
- **[Security]:** Upgrade-internal files (state, staged binary) use owner-only permissions, per `security.md` and `engineering.md`.
- **[Observability]:** With `DEBUG=1`, every updater decision (skipped: throttled / not writable / older / other target; failed: reason) is written to the error stream; without it the updater is silent except for the interactive notices above.
- **[Compatibility]:** State files written by earlier versions load without error. New rules are implemented in platform-independent code; the Windows-specific swap mechanism is kept as is and covered only through the shared logic.
- **[Scope]:** Native binary installs keep GitHub releases as their version source; release signing is out of scope.
