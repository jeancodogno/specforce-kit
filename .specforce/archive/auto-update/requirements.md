---
slug: auto-update
lens: Balanced full-stack
---

# Feature: Auto-Update Capability

## 1. Context & Value
Specforce users need to stay aligned with the latest engine improvements and AI skill updates without manual maintenance overhead. This feature introduces a non-blocking update mechanism that keeps human developers informed of new versions while ensuring that AI-driven workflows remain uninterrupted. By automating the detection and installation process, we reduce version fragmentation and support costs.

## 2. Out of Scope (Anti-Goals)
- **Non-TTY updates:** Automatically triggering update checks or notifications in non-interactive environments (CI/CD, piped commands).
- **Blocking checks:** Delaying the start or execution of any command to perform a network check.
- **Forced updates:** Updating the CLI without explicit user confirmation in the interactive terminal.
- **In-process updates:** Replacing the binary while the primary command is still executing.

## 3. Acceptance Criteria (BDD)
### [REQ-1] Automatic Background Update Check
**User Story:** AS A user, I WANT the CLI to check for updates in the background, SO THAT I am always informed of new features and fixes without manual checking.
**Scenarios:**
1. **[Happy Path]** GIVEN the user has an interactive terminal (TTY) WHEN they run 'init', 'console', or 'install' THEN a background check for updates is triggered asynchronously.
2. **[Throttling]** GIVEN the last update check was performed less than 24 hours ago WHEN a command is run THEN no background check is triggered to preserve bandwidth and performance.
3. **[Silent Failure]** GIVEN the background check encounters a network error or timeout WHEN the main command is running THEN the error is suppressed and the user experience is not affected.

### [REQ-2] Interactive Update Notification
**User Story:** AS A user, I WANT to be notified of a new version in a visually consistent way, SO THAT I can decide whether to upgrade.
**Scenarios:**
1. **[Happy Path]** GIVEN a new version is detected and the terminal is interactive WHEN a 'human-only' command finishes THEN a "New Version Available" message is displayed in the terminal footer.
2. **[Non-Interactive Protection]** GIVEN a new version is detected but the terminal is NOT a TTY WHEN any command finishes THEN NO update message is displayed.
3. **[LLM-Driven Protection]** GIVEN a new version is detected and the terminal is a TTY WHEN an LLM-driven command (e.g., 'spec', 'implement', 'constitution') is run THEN the update check is skipped and no notification is shown to prevent polluting agent logs.

### [REQ-3] Automated Installation
**User Story:** AS A user, I WANT the CLI to handle the upgrade process automatically, SO THAT I don't have to manually run NPM or download binaries.
**Scenarios:**
1. **[NPM Path]** GIVEN the CLI was installed via NPM and the user confirms the upgrade WHEN the upgrade process starts THEN the system executes the global npm install command and verifies the new version.
2. **[Binary Path]** GIVEN the CLI was installed as a standalone binary and the user confirms the upgrade WHEN the upgrade process starts THEN the system downloads the latest release for the current OS/Arch, replaces the active binary, and validates the checksum.
3. **[Rollback]** GIVEN an automated installation fails WHEN the process is interrupted THEN the system ensures the previous version of the binary remains functional or provides clear instructions for manual recovery.

## 4. Business Invariants
- Update checks MUST be non-blocking and execute in a separate process or goroutine to ensure zero latency for the primary command.
- The update notification MUST only appear AFTER the primary command has successfully completed its output.
- The system MUST maintain a global state (e.g., in `~/.specforce/state.json`) to track the timestamp of the last check.

## 5. UI/UX Contract
- **Visual Posture:** Notifications use the "Neon" theme style (e.g., Cyan/Magenta text) consistent with the Specforce brand.
- **Notification Type:** A discreet box or line at the end of the command output, never a popup or intrusive interrupt.
- **Interaction:** Confirmations for upgrades use the standard TUI prompt component.
