---
slug: auto-update
lens: Balanced full-stack
---

# Implementation Roadmap: Auto-Update Capability

## 1. Execution Strategy
- **Gravity Order:** Internal Service Logic -> CLI Hook Integration -> UI/UX Components -> Installation Execution -> End-to-End Verification.
- **Resilience:** Focus on non-blocking background operations and atomic file updates for the state and binary.

## 2. Tasks

### Phase 1: Core Upgrade Infrastructure

#### T1.1: [CODE] Define Global State Management
**State:** `[FINISHED]`
**Target:** `src/internal/upgrade/state.go`
**Context:** [REQ-1]
**Action Steps:**
- Create `State` struct with `LastCheckAt` (time.Time), `LatestVersion` (string), and `IgnoredVersion` (string).
- Implement `LoadState()` and `SaveState(state)` using `core.ExpandPath("~/.specforce/state.json")`.
- Implement atomic write (write to `.tmp` + rename) for `SaveState` to prevent file corruption.
**Verification (TDD):**
- Write unit test in `state_test.go` asserting that state is correctly persisted and reloaded from a temporary directory.

#### T1.2: [CODE] Define Provider Interface & Mock
**State:** `[FINISHED]`
**Target:** `src/internal/upgrade/provider.go`
**Context:** [REQ-1]
**Action Steps:**
- Define `Provider` interface with `GetLatestVersion(ctx context.Context) (string, error)`.
- Implement `MockProvider` for testing purposes.
**Verification (TDD):**
- Unit test ensuring `MockProvider` returns a hardcoded version string.

#### T1.3: [CODE] Implement GitHub Release Provider
**State:** `[FINISHED]`
**Target:** `src/internal/upgrade/provider_github.go`
**Context:** [REQ-1]
**Action Steps:**
- Implement `GitHubProvider` targeting `api.github.com/repos/jeancodogno/specforce-kit/releases/latest`.
- Parse the `tag_name` from the JSON response.
- Add a 2-second timeout to the HTTP client.
**Verification (TDD):**
- Unit test using `httptest.NewServer` to simulate GitHub API response.

#### T1.4: [CODE] Implement NPM Registry Provider
**State:** `[FINISHED]`
**Target:** `src/internal/upgrade/provider_npm.go`
**Context:** [REQ-1]
**Action Steps:**
- Implement `NPMProvider` targeting `registry.npmjs.org/@specforce/cli/latest`.
- Parse the `version` field from the JSON response.
**Verification (TDD):**
- Unit test using `httptest.NewServer` to simulate NPM registry response.

#### T1.5: [CODE] Implement SemVer Comparison Utility
**State:** `[FINISHED]`
**Target:** `src/internal/upgrade/semver.go`
**Context:** [REQ-1]
**Action Steps:**
- Add helper to compare two version strings (e.g., `v1.2.3` vs `1.2.4`).
- Support both `v`-prefixed and bare versions.
**Verification (TDD):**
- Unit tests covering "greater than", "equal", and "less than" scenarios for semver strings.

#### T1.6: [CODE] Implement Orchestration Service
**State:** `[FINISHED]`
**Target:** `src/internal/upgrade/service.go`
**Context:** [REQ-1]
**Action Steps:**
- Create `Service` struct that coordinates `State` and `Provider`.
- Implement `CheckForUpdate(ctx)`: checks if 24h have passed since `LastCheckAt`, then runs provider fetch in a background goroutine.
- Implement `IsUpdateAvailable()`: returns true if `LatestVersion` > current version.
**Verification (TDD):**
- Unit test using a mock clock and mock provider to verify that background check is skipped if throttled or already checked recently.

### Phase 2: CLI Integration

#### T2.1: [CODE] Mark Agent Commands with Annotations
**State:** `[FINISHED]`
**Target:** `src/internal/cli/cobra/*.go`
**Context:** [REQ-2]
**Action Steps:**
- Add `Annotations: map[string]string{"IsAgentCommand": "true"}` to `specCmd`, `implementCmd`, and `constitutionCmd`.
**Verification (TDD):**
- Unit test in `root_test.go` asserting these specific commands have the required annotation.

#### T2.2: [CODE] Integrate PersistentPreRun Hook
**State:** `[FINISHED]`
**Target:** `src/internal/cli/cobra/root.go`
**Context:** [REQ-1]
**Action Steps:**
- In `rootCmd.PersistentPreRunE`, initialize `upgrade.Service`.
- Call `service.CheckForUpdate(ctx)` ONLY if `tui.IsTTY()` is true AND command annotation `IsAgentCommand` is not "true".
**Verification (TDD):**
- Manual verification using `DEBUG=1` logs to ensure background check is triggered/skipped based on context.

#### T2.3: [CODE] Integrate PersistentPostRun Hook
**State:** `[FINISHED]`
**Target:** `src/internal/cli/cobra/root.go`
**Context:** [REQ-2]
**Action Steps:**
- In `rootCmd.PersistentPostRun`, check `service.IsUpdateAvailable()`.
- If true and TTY is active, trigger the notification display.
**Verification (TDD):**
- Manual verification: run a standard command and see the notification appear AFTER the main output.

### Phase 3: UI/UX Components

#### T3.1: [CODE] Create Lipgloss Notification Box
**State:** `[FINISHED]`
**Target:** `src/internal/tui/upgrade_notification.go`
**Context:** [REQ-2]
**Action Steps:**
- Implement `RenderUpdateNotification(current, latest string)` using Lipgloss.
- Apply "Neon" theme colors (Cyan/Magenta) to match brand guidelines.
**Verification (TDD):**
- Visual verification by running a standalone TUI test script that renders the box to stdout.

#### T3.2: [CODE] Create Interactive Upgrade Prompt
**State:** `[FINISHED]`
**Target:** `src/internal/tui/upgrade_prompt.go`
**Context:** [REQ-3]
**Action Steps:**
- Create a Bubbletea model for a simple "Update now? [y/N]" interactive prompt.
**Verification (TDD):**
- Manual test in an interactive terminal verifying selection logic and default behavior.

#### T3.3: [CODE] Create Upgrade Progress View
**State:** `[FINISHED]`
**Target:** `src/internal/tui/upgrade_progress.go`
**Context:** [REQ-3]
**Action Steps:**
- Implement a Bubbletea model with a progress bar (using existing `tui.ProgressBar`) to show download status.
**Verification (TDD):**
- Manual test with a simulated 3-second delay to verify progress bar rendering and completion message.

### Phase 4: Installation Logic

#### T4.1: [CODE] Implement NPM Upgrade Strategy
**State:** `[FINISHED]`
**Target:** `src/internal/upgrade/installer_npm.go`
**Context:** [REQ-3]
**Action Steps:**
- Implement logic to run `npm install -g @jeancodogno/specforce-kit` via `os/exec`.
- Capture and log stderr on failure for user feedback.
**Verification (TDD):**
- Unit test mocking `exec.Command` to ensure correct arguments and environment are used.

#### T4.2: [CODE] Implement Binary Download & Checksum
**State:** `[FINISHED]`
**Target:** `src/internal/upgrade/installer_binary.go`
**Context:** [REQ-3]
**Action Steps:**
- Implement download logic for the binary based on `runtime.GOOS` and `runtime.GOARCH`.
- Download corresponding `.sha256` file and verify hash before installation.
**Verification (TDD):**
- Unit test with `httptest` serving a dummy file and verifying correct/incorrect hash detection.

#### T4.3: [CODE] Implement Atomic Binary Replacement
**State:** `[FINISHED]`
**Target:** `src/internal/upgrade/installer_binary.go`
**Context:** [REQ-3]
**Action Steps:**
- Use `os.Executable()` to find current binary path.
- Implement "Move to Backup -> Move New -> Remove Backup" atomic flow.
- Ensure rollback to the backup if the replacement fails.
**Verification (TDD):**
- Unit test simulating a restricted permission error during move to verify original binary remains functional.

#### T4.4: [CODE] Wire Service to PerformUpgrade
**State:** `[FINISHED]`
**Target:** `src/internal/upgrade/service.go`
**Context:** [REQ-3]
**Action Steps:**
- Implement `PerformUpgrade(ctx, version)`.
- Detect installation type (checking execution path or metadata).
- Invoke the appropriate installer strategy (NPM vs Binary).
**Verification (TDD):**
- Integration test ensuring `PerformUpgrade` triggers the expected installer based on environment.

### Phase 5: Final Verification

#### T5.1: [TEST] Full End-to-End Simulation
**State:** `[FINISHED]`
**Target:** `src/internal/upgrade/integration_test.go`
**Context:** [REQ-1, REQ-2, REQ-3]
**Action Steps:**
- Set up a full CLI execution flow in a controlled test environment.
- Mock network to return a new version.
- Simulate user "Y" input to the upgrade prompt.
- Verify that the dummy binary is correctly replaced and version updated.
**Verification (TDD):**
- Execute `go test ./src/internal/upgrade/...` and ensure all integration tests pass.

#### T5.2: [TEST] Verify Silence in Non-TTY & Agent Modes
**State:** `[FINISHED]`
**Target:** `src/internal/upgrade/integration_test.go`
**Context:** [REQ-2]
**Action Steps:**
- Run commands with redirected stdout (non-TTY).
- Run `spec` command in a TTY.
- Assert that no update check is performed and no notification is rendered.
**Verification (TDD):**
- Integration test checking logs/output for absence of upgrade indicators in restricted modes.

#### T5.3: [TEST] Cross-Platform Path Verification
**State:** `[FINISHED]`
**Target:** `src/internal/upgrade/state_test.go`
**Context:** [REQ-1]
**Action Steps:**
- Run tests on Linux and macOS (via GitHub Actions) to ensure `~/.specforce` expansion and file permissions work correctly.
**Verification (TDD):**
- CI Pipeline green on all target platforms for the upgrade package.
