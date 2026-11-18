---
prefix: UPD
---

# Module: Auto-update

## 1. Domain Scope
Keeps an installed Specforce CLI on the newest release of its own channel without blocking or polluting the user's commands. It checks for new versions in a detached background process, stages verified native binaries, swaps them in on the next run, installs exact versions for npm installations, and tells interactive users what happened (or what they must do by hand) exactly once. It never downgrades, never touches an executable it was not staged for, and never writes to a command's standard output.

## 2. Business Rules & Invariants
- `[BR-UPD-01]` A remote version check runs at most once every 6 hours per user. The check time is persisted whether the check succeeds or fails, so being offline never causes repeated checks.
- `[BR-UPD-02]` The background check is only started by interactive (terminal) commands that are not agent commands; it runs detached and never delays the command.
- `[BR-UPD-03]` The installation method decides the version source: native binary installs follow the latest GitHub release; npm installs follow the npm registry tag of the running version's channel (pre-release identifier such as `alpha` or `beta`, or `latest` for stable versions). An installation is npm only when the executable lives inside the Specforce npm platform package.
- `[BR-UPD-04]` An update is never applied unless it is strictly newer than the running version. Older or equal staged versions are discarded.
- `[BR-UPD-05]` A staged update belongs to the exact executable path it was staged for. Any other executable ignores it and leaves it intact.
- `[BR-UPD-06]` A native binary is staged only when the executable's directory is writable; otherwise the version is recorded as requiring a manual update and nothing is downloaded.
- `[BR-UPD-07]` A staged binary replaces the executable only if it is a regular file, matches the fingerprint (release-verified hash and size) recorded at staging, reports the staged version in an isolated test run, and still matches the fingerprint right after that run. Any doubt discards it.
- `[BR-UPD-08]` A verification failure is retried once (the version is staged again by the next check); a second failure for the same version makes it a manual update and it is not downloaded again.
- `[BR-UPD-09]` A failed swap (the executable cannot be moved or replaced) is terminal for that version: the staged file is deleted, the pending flag cleared and the version recorded as requiring a manual update. It is never retried in a loop.
- `[BR-UPD-10]` If the new executable fails to start after the swap, the previous executable is restored, the version is not reported as applied, and it becomes a manual update.
- `[BR-UPD-11]` npm installs install exactly `@<version>` of the package; a missing channel tag means "no update"; a failed install makes the version a manual update and is not retried.
- `[BR-UPD-12]` Update notices ("updated", "available", "manual update required") are written to standard error only, only for interactive non-agent commands. The manual hint is shown once per version, names the version and the exact manual command, and never claims the update will be applied automatically.
- `[BR-UPD-13]` The swap check runs before every command, agent or not, and writes nothing to standard output.
- `[BR-UPD-14]` An unparseable state file is treated as empty and repaired by the next check; any staged binary without trusted state is deleted. A state path that cannot be read (for example a directory) disables the updater for that run without crashing and is never overwritten.
- `[BR-UPD-15]` Update state and staged binaries are owner-only (`0600` files, `0700` staging directory); temporary downloads are removed on every failure path.

## 3. Canonical Requirements & Use Cases
### [US-UPD-01] Background update of a native binary
- **Scenario:** Newer release, writable install
  - **GIVEN** a native install in a writable directory running version V and a newer release N
  - **WHEN** an interactive command runs more than 6 hours after the last check
  - **THEN** a detached check downloads N, verifies it against the release checksum and stages it for this executable
  - **AND** the next command verifies the staged file, replaces the executable, restarts on N, and the next interactive command shows "(Specforce updated to N)" once on standard error.
- **Scenario:** Offline
  - **GIVEN** the release source is unreachable
  - **WHEN** interactive commands run repeatedly
  - **THEN** only the first one starts a check and later ones within 6 hours start none.

### [US-UPD-02] Install the user cannot write to
- **Scenario:** Root-owned install directory
  - **GIVEN** the executable lives in a directory the user cannot write to and a newer version N exists
  - **WHEN** the background check runs
  - **THEN** nothing is downloaded and N is recorded as requiring a manual update
  - **AND** the next interactive command shows "MANUAL UPDATE REQUIRED" with N and the manual command once, and never again for N.
- **Scenario:** Swap fails anyway
  - **GIVEN** a staged update whose swap fails (e.g. permission denied)
  - **WHEN** any command starts
  - **THEN** the command continues on the current version, the update is cleared, and the manual hint behavior applies; the swap is not attempted again.

### [US-UPD-03] Safety guards at apply time
- **Scenario:** Downgrade attempt
  - **GIVEN** the staged version is older than or equal to the running version
  - **WHEN** any command starts
  - **THEN** the executable is not replaced and the staged file is discarded.
- **Scenario:** Another executable
  - **GIVEN** an update staged for executable A
  - **WHEN** executable B (e.g. a development build) runs
  - **THEN** B is not replaced and A's staged update stays intact.
- **Scenario:** Corrupted or unrunnable staged binary
  - **GIVEN** a staged file that is empty, altered, a symlink, not runnable, reports another version, hangs, or changes during its own test run
  - **WHEN** any command starts
  - **THEN** the executable is untouched and the staged file discarded; the first failure lets the next check re-stage the version, the second makes it a manual update.
- **Scenario:** Restart failure
  - **GIVEN** the executable was replaced and the new one cannot be started
  - **WHEN** the restart fails
  - **THEN** the previous executable is restored to its path and the version becomes a manual update.

### [US-UPD-04] npm channel updates
- **Scenario:** Beta user
  - **GIVEN** an npm install running `1.0.0-beta.2` and the registry tag `beta` pointing to `1.0.0-beta.3`
  - **WHEN** the background check runs
  - **THEN** `npm install -g @jeancodogno/specforce-kit@1.0.0-beta.3` is executed (never the `latest` tag of another channel).
- **Scenario:** Channel without tag
  - **GIVEN** the running channel has no registry tag
  - **WHEN** the check runs
  - **THEN** nothing is installed or shown and the check time is persisted.
- **Scenario:** npm install fails
  - **GIVEN** the global install fails
  - **WHEN** the check finishes
  - **THEN** the version becomes a manual update with the exact npm command, and it is not retried.

### [US-UPD-05] Agent-safe output
- **Scenario:** Agent command during a pending swap
  - **GIVEN** an update is pending
  - **WHEN** an agent runs a command with `--json`
  - **THEN** the update is applied and standard output contains only valid JSON; no notice is written to any stream.

### [US-UPD-06] Damaged state
- **Scenario:** Garbage state file
  - **GIVEN** the state file contains invalid content and an orphan staged binary exists
  - **WHEN** any command or background check runs
  - **THEN** the state is treated as empty, the orphan is deleted, the next check writes a valid state, and no command fails.

## 4. Public Integration Surfaces & Contracts
- **Implicit trigger:** every `specforce` command (swap check before parsing; background check and notices for interactive non-agent commands).
- **Internal flag:** `specforce --internal-upgrade-check` runs one check, stages or installs, and exits (used by the detached process; exit code 1 on failure, silent).
- **Diagnostics:** `DEBUG=1` writes each updater decision and failure reason (throttled, not writable, not newer, other target, verification failure, rollback, state repair) to standard error.
- **Per-user state:** `~/.specforce/state.json` (owner-only JSON; unknown keys ignored, older files load) and `~/.specforce/upgrade/staged/`.
- **External sources (read-only):** GitHub releases (`/releases/latest`, release assets and `specforce-kit_<version>_checksums.txt`) and the npm registry dist-tags of `@jeancodogno/specforce-kit`.
- **Manual update commands shown to users:** native installs — download the release or `sudo make install` from source; npm — `npm install -g @jeancodogno/specforce-kit@<version>`.
- **Cross-module:** agent commands are identified by the `IsAgentCommand` annotation used across the CLI (see `agent-kit`, `spec-management`).

## 5. Operational & Quality Invariants
- The background check finishes within 1 minute; every network request is bound to that budget.
- The pre-swap test run of a staged binary runs only `--version` with an empty temporary home as working directory, is bounded by 5 seconds, and kills its whole process group on timeout (Unix).
- The foreground cost of a pending swap is one state read, one fingerprint check and one test run, once per version.
- The swap never leaves a corrupted executable: it moves the current executable to `<path>.old`, which is removed only after a successful restart (on the next normal command).
- Windows keeps a spawn-and-exit restart; all rules above live in platform-independent logic.
- Known limitation: concurrent first runs racing on the same pending swap can record an inconsistent outcome (no corrupted executable); a cross-process lock is not yet in place.
