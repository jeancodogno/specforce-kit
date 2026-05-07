---
slug: 20260506-1701-automated-background-updates
lens: Backend-heavy
---

# Implementation Roadmap: Automated Background Updates

## 1. Execution Strategy
- **Gravity Order**: State/Directory Management -> Background Engine -> Atomic Swapper -> CLI Integration -> Cleanup.
- **Verification**: Strictly TDD. Every task must pass its specific unit or integration test before proceeding. Background operations must be verified through state-file markers and process inspection.

## 2. Tasks

### Phase 1: Scaffolding & State

- [x] T1.1: Update Upgrade State Definition
**Target:** `src/internal/upgrade/state.go`
**Context:** [REQ-1]

**Action Steps:**
- Add `StagedVersion` (string), `LastCheckedAt` (time.Time), and `UpdateReady` (bool) fields to the `State` struct.
- Update JSON tags for persistence.

**Verification (TDD):**
Run `go test ./src/internal/upgrade/state_test.go` asserting that the state can be serialized to and deserialized from JSON with the new fields preserved.

- [x] T1.2: Implement Staging Directory Management
**Target:** `src/internal/upgrade/state.go`
**Context:** [REQ-2]

**Action Steps:**
- Implement `EnsureStagedDir()` to manage the creation of `~/.specforce/upgrade/staged/`.
- Handle cross-platform path resolution.

**Verification (TDD):**
Unit test in `src/internal/upgrade/state_test.go` asserting that the directory is created with correct permissions (0755) and its path is correctly resolved.

### Phase 2: Background Engine

- [x] T2.1: Implement Detached Background Launcher
**Target:** `src/internal/upgrade/service.go`
**Context:** [REQ-1]

**Action Steps:**
- Create `LaunchBackgroundCheck()` using `os.StartProcess` with detached process attributes.
- Ensure it spawns `specforce --internal-upgrade-check`.

**Verification (TDD):**
Test in `src/internal/upgrade/service_test.go` asserting the parent process returns immediately while the child process ID is registered in the OS.

- [x] T2.2: Implement Internal Check CLI Command
**Target:** `src/cmd/specforce/main.go`
**Context:** [REQ-1]

**Action Steps:**
- Add a hidden flag `--internal-upgrade-check` (or hidden command).
- Execute version check and download logic silently.

**Verification (TDD):**
Run `specforce --internal-upgrade-check` manually and verify `LastCheckedAt` is updated in `~/.specforce/upgrade/state.json`.

- [x] T2.3: Implement Version Comparison & Asset Download
**Target:** `src/internal/upgrade/service.go`
**Context:** [REQ-2, REQ-4]

**Action Steps:**
- Implement `DownloadLatestIfNewer()`.
- Integrate with GitHub/NPM providers.
- Download and verify checksums/signatures.

**Verification (TDD):**
Integration test using a mock HTTP server in `src/internal/upgrade/integration_test.go` asserting the binary is downloaded and the `UpdateReady` flag is set to true.

### Phase 3: Atomic Swapper

- [x] T3.1: Implement Rename-to-Old Logic
**Target:** `src/internal/upgrade/service.go`
**Context:** [REQ-3]

**Action Steps:**
- Implement `PerformAtomicSwap()`.
- Handle `os.Rename` for the current binary to `.old`.
- Move staged binary to target path and set permissions.

**Verification (TDD):**
Unit test asserting that if the swap fails at any step, the original binary is restored and the environment remains functional.

- [x] T3.2: Implement syscall.Exec Wrapper
**Target:** `src/internal/upgrade/service.go`
**Context:** [REQ-3]

**Action Steps:**
- Implement `ExecuteNewBinary()` using `syscall.Exec` on Unix systems.
- Handle Windows fallback (detached restart).

**Verification (TDD):**
Integration test verifying that the process effectively restarts with the new binary version without returning to the caller.

### Phase 4: CLI Integration

- [x] T4.1: Hook Background Check into PersistentPreRunE
**Target:** `src/internal/cli/cobra/root.go`
**Context:** [REQ-1]

**Action Steps:**
- Update `PersistentPreRunE` to call `LaunchBackgroundCheck()` based on throttling logic (24h or 6h as per REQ).

**Verification (TDD):**
Run `specforce version` and verify via logs or state-file timestamps that the check was triggered in the background.

- [x] T4.2: Implement Post-Command Swap Trigger
**Target:** `src/cmd/specforce/main.go`
**Context:** [REQ-3, REQ-5]

**Action Steps:**
- After `rootCmd.Execute()` returns, check `upgrade.State.UpdateReady`.
- If true, invoke `PerformAtomicSwap()` and `ExecuteNewBinary()`.
- Print the subtle notification: `(Specforce updated to vX.Y.Z)`.

**Verification (TDD):**
Manually set `UpdateReady: true` in state.json and run any command. Verify the command completes, notification is printed, and binary is swapped.
**Manual Testing Note:** Verified that capturing `os.Executable()` *before* the swap is critical for `syscall.Exec` to restart the correct process.

### Phase 5: Cleanup & Migration

- [x] T5.1: Remove Legacy Install Command
**Target:** `src/internal/cli/cobra/install.go`
**Context:** [MIGRATION]

**Action Steps:**
- Delete `install.go` and remove its registration from `rootCmd`.
- Remove unused installer logic in `src/internal/installer/`.

**Verification (TDD):**
Run `go build` and assert that `specforce install` results in an "unknown command" error.

- [x] T5.2: Implement Automatic .old Cleanup
**Target:** `src/internal/upgrade/service.go`
**Context:** [CLEANUP]

**Action Steps:**
- Add logic to delete any `specforce.old` files from previous successful upgrades.

**Verification (TDD):**
Run a full upgrade cycle and verify that the `specforce.old` file is removed upon the next command execution.
**Manual Testing Note:** Logic moved to `PersistentPreRunE` in `root.go` to ensure cleanup happens early and reliably.
