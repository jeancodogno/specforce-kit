---
slug: security-hardening-socket-alerts
lens: Backend-heavy
---

# Feature: Security Hardening (Socket.dev Alerts)

## 1. Context & Value
Address security alerts from socket.dev by removing supply chain risks associated with automated install scripts and hardening the native binary proxy. This improves the project's security posture and user trust by adhering to "Zero Scripts" installation patterns.

## 2. Out of Scope (Anti-Goals)
- Do not implement a new package manager or installation logic.
- Do not modify the Go binary source code or functionality.
- Do not change the overall architecture of the proxy, only its validation and script dependencies.

## 3. Acceptance Criteria (BDD)

### [REQ-1] Zero Install Scripts
**User Story:** AS A user/developer, I WANT the package to have no automated installation scripts, SO THAT I am protected from potential supply chain attacks executed during `npm install`.

**Scenarios:**
1. **[Happy Path]** GIVEN a clean environment WHEN I run `npm install @jeancodogno/specforce-kit` THEN no `postinstall` or `prepare` scripts are executed.
2. **[Fallback Path]** GIVEN the native binary is missing after install WHEN I run `specforce` THEN I receive a clear diagnostic message explaining how to fix it manually.
3. **[Edge Case]** GIVEN a developer environment WHEN I run `npm install` for local development THEN the project still builds successfully via manual commands (e.g., `make build`).

### [REQ-2] Hardened Binary Proxy
**User Story:** AS A security-conscious user, I WANT the Node.js proxy to strictly validate the binary it executes, SO THAT unauthorized or malicious files are not inadvertently run.

**Scenarios:**
1. **[Happy Path]** GIVEN a valid platform-specific binary is present WHEN I run `specforce` THEN the proxy verifies the path is absolute and the file exists before execution.
2. **[Edge Case]** GIVEN a manipulated `binaryPath` that is not absolute or doesn't exist WHEN I run `specforce` THEN the proxy must abort execution and show a diagnostic error.

## 4. Business Invariants
- The package MUST NOT contain any automated lifecycle scripts (`preinstall`, `postinstall`, `prepare`) that execute shell commands.
- The `index.js` MUST NOT use shell-mediated execution (e.g., `sh -c`) for the binary; it must continue to use the argument-slice pattern of `spawnSync`.
- All binary path resolutions MUST result in an absolute path before execution.
