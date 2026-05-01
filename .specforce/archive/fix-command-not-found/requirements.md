---
slug: fix-command-not-found
lens: Backend-heavy
---

# Feature: Fix Command Not Found

## 1. Context & Value
Problem: Users often encounter a "command not found" error after installing `specforce` globally via NPM or from source. This creates a friction-filled first impression and blocks adoption.
Value: Automating the build process and providing precise diagnostics reduces support overhead and ensures a "plug-and-play" experience for all installation methods.

## 2. Out of Scope (Anti-Goals)
* **Automatic PATH modification:** The system MUST NOT attempt to modify the user's shell profile (`.bashrc`, `.zshrc`, etc.) automatically to avoid security risks and side effects.
* **Dependency Management:** This feature does NOT cover installing system dependencies like `make` or `go`. It assumes they are present if building from source.

## 3. Acceptance Criteria (BDD)

### [REQ-1] Diagnostic Error Messages
**User Story:** AS A user, I WANT TO see a helpful diagnostic message when the command fails, SO THAT I can quickly fix my environment.

**Scenarios:**
1. **[Happy Path] PATH Diagnostic:** GIVEN the `specforce` command is executed but the proxy fails to locate the native binary WHEN the NPM global bin directory is not in the system's `PATH` THEN the diagnostic output MUST display the exact shell command required to add the directory to the `PATH`.
2. **[Happy Path] Build Diagnostic:** GIVEN the `specforce` command is executed but the binary is missing WHEN the `PATH` is correct but the binary was never built (e.g., manual Git clone) THEN the diagnostic output MUST suggest running `make build`.
3. **[Edge Case] Cross-Platform Support:** GIVEN a Windows environment WHEN a PATH issue is detected THEN the diagnostic MUST provide the equivalent Windows command (`setx PATH ...`) or instructions for the System Environment Variables UI.

### [REQ-2] Automated Build on Installation
**User Story:** AS A developer installing from source/Git, I WANT the native binary to be built automatically, SO THAT the command works immediately.

**Scenarios:**
1. **[Happy Path] NPM Install Trigger:** GIVEN a local Git clone WHEN the user runs `npm install` THEN a `postinstall` script MUST be triggered to run `make build`.
2. **[Negative Path] Graceful Build Failure:** GIVEN the automated build fails (e.g., `make` or `go` missing) WHEN running `npm install` THEN the installation MUST NOT crash; instead, it MUST print a clear warning informing the user that the binary could not be built automatically and providing manual instructions.

### [REQ-3] Manual Fix Documentation
**User Story:** AS A user whose automated install failed, I WANT clear documentation for manual recovery, SO THAT I am not stuck.

**Scenarios:**
1. **[Happy Path] Troubleshooting Guide:** GIVEN the project documentation (README or dedicated troubleshooting file) WHEN a user searches for installation issues THEN there MUST be a section covering "Command Not Found" with step-by-step instructions for PATH verification and manual building.

## 4. Business Invariants
* The proxy script (`index.js`) must remain lightweight and avoid adding heavy dependencies.
* The diagnostic messages must be localized to the user's OS where possible (Unix vs. Windows).
* The `postinstall` script must check for the presence of build tools before attempting to execute them to avoid noisy errors.
