---
slug: 20260506-1701-automated-background-updates
lens: Backend-heavy
---

# Feature: Automated Background Updates

## 1. Context & Value
Currently, Specforce Kit requires manual updates via `specforce install`, which leads to version fragmentation and users running outdated SDD protocols. This feature automates the update lifecycle—performing background checks, staging downloads, and swapping binaries on the next run—ensuring all agents and developers are always using the most secure and feature-complete version of the framework without manual intervention.

## 2. Out of Scope (Anti-Goals)
- Automatic updates for the project's own `go.mod` or `package.json` dependencies.
- Self-update of the specialized agent "skills" (managed via separate SDD workflows).
- Real-time binary hot-swapping while a command is executing (must wait for next execution).
- Support for updating Specforce if installed via OS-level package managers (e.g., `brew`, `apt`) unless they follow the standalone binary pattern.

## 3. Acceptance Criteria (BDD)

### [REQ-1] Background Version Checking & Throttling
**User Story:** AS A developer, I WANT the tool to silently check for updates when I run commands, SO THAT I don't have to manually track releases.

**Scenarios:**
1. **[Happy Path]** GIVEN the user executes a `specforce` command WHEN the execution starts THEN the system spawns a background process to check the latest version from the source (NPM or GitHub) AND the current command continues without delay.
2. **[Throttling]** GIVEN a successful version check was performed in the last 6 hours WHEN a command is executed THEN the system skips the remote check to minimize network overhead.
3. **[Error Tolerance]** GIVEN the user is offline or the registry is down WHEN the version check fails THEN the error is logged silently to the debug log AND the user's current command is unaffected.

**Technical Constraints (NFR):**
- **[Performance]:** The check trigger must add < 10ms to command startup.
- **[Reliability]:** Background check must be "fire-and-forget" with no risk of blocking the main thread.
- **[Integrity]:** The check must respect the installation method (Binary vs NPM) to query the correct source.

### [REQ-2] Background Staging (Next-Run Preparation)
**User Story:** AS A developer, I WANT updates to be downloaded and prepared while I work, SO THAT they are ready to use immediately on my next interaction.

**Scenarios:**
1. **[Happy Path]** GIVEN a new version is detected WHEN the background process is running THEN it downloads the appropriate package AND stages the new binary in a secure temporary location (e.g., `.specforce/bin/next/`).
2. **[Integrity Check]** GIVEN a download is complete WHEN staging the binary THEN the system verifies the checksum/signature of the new version AND only marks it as "Ready" if valid.

**Technical Constraints (NFR):**
- **[Safety]:** Use atomic file operations to move the downloaded file to the staging area.
- **[Security]:** Staged binaries must be set with `0755` permissions.
- **[Efficiency]:** Download must be throttled or low-priority to avoid consuming all available bandwidth.

### [REQ-3] Atomic Next-Run Binary Swap
**User Story:** AS A developer, I WANT the tool to automatically upgrade itself, SO THAT I always have the latest fixes without running an installer.

**Scenarios:**
1. **[Happy Path]** GIVEN a verified binary is staged in the "next" directory WHEN a new `specforce` command is executed THEN the system replaces the active binary with the staged one AND continues execution with the new version.
2. **[Permissions Edge Case]** GIVEN the system lacks write permissions to the install directory WHEN an upgrade is attempted THEN it logs the failure AND continues with the old version AND notifies the user that a manual `specforce install` is required.

**Technical Constraints (NFR):**
- **[Safety]:** The swap MUST be atomic (e.g., `os.Rename` or similar) to prevent binary corruption.
- **[Observability]:** Log the upgrade event with old and new version strings in the `memorial.md`.

### [REQ-4] Multi-Provider Update Strategy (NPM & Binary)
**User Story:** AS A user, I WANT the update to work regardless of how I installed the tool, SO THAT the experience is consistent across environments.

**Scenarios:**
1. **[NPM Install]** GIVEN Specforce was installed via `npm -g` WHEN an update is available THEN the background process triggers `npm install -g @specforce/kit` (respecting user permissions).
2. **[Binary Install]** GIVEN Specforce was installed as a standalone binary WHEN an update is available THEN the background process downloads the specific architecture-compatible binary from GitHub Releases.

**Technical Constraints (NFR):**
- **[Integrity]:** Implementation must correctly detect the active installation method at runtime.

### [REQ-5] Subtle Update Notification
**User Story:** AS A developer, I WANT to be notified that an update occurred, SO THAT I am aware of new capabilities or fixed bugs.

**Scenarios:**
1. **[Happy Path]** GIVEN a command is running on a version that was just swapped WHEN the command output is finished THEN the system prints a subtle message: `(Specforce updated to vX.Y.Z)`.

**UI/UX Specifics:**
- **Feedback Logic:** Notification appears as the very last line of output.
- **Style:** Dimmed or color-coded (e.g., light blue) to distinguish it from command-specific results.

## 4. Business Invariants
- **Zero-Block Policy:** No update-related operation (check, download, staging) shall ever block or delay a user-initiated command.
- **Atomic Swap:** A binary update must either succeed completely or leave the original binary untouched; partial writes are forbidden.
- **Throttling:** Remote version checks must be performed at most once every 6 hours per machine.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Background download must use a low-priority thread/process.
- **[Reliability]:** Fail-safe to the current version on any network or filesystem error.
- **[Security]:** All downloads must occur over HTTPS; checksum verification is mandatory for standalone binaries.
- **[Maintainability]:** The update logic must be encapsulated in the `src/internal/upgrade` package.
