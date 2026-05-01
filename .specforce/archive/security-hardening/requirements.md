---
slug: security-hardening
lens: Backend-heavy
---

# Feature: Security Hardening

## 1. Context & Value
The project has been scanned for security vulnerabilities using `gosec`, revealing critical issues related to goroutine context management, subprocess command sanitization, and file permissions. This feature hardens the application against potential resource leaks, command injection, and unauthorized data access, while also automating security enforcement in the CI pipeline.

## 2. Out of Scope (Anti-Goals)
* Implementing a full-blown Authentication/Authorization system.
* Encrypting project files or git history.
* Rewriting the entire subprocess execution logic (only fixing the identified insecure patterns).
* Adding support for external secret managers (Vault, etc.).

## 3. Acceptance Criteria (BDD)

### [REQ-1] Safe Goroutine Context Propagation
**User Story:** AS A Developer, I WANT TO ensure goroutines don't inherit a cancelable context that could prematurely terminate background tasks, SO THAT long-running operations complete reliably even if the parent request/process context is canceled.

**Scenarios:**
1. **[Happy Path]** GIVEN a background goroutine started from a cancelable context WHEN the parent context is canceled THEN the goroutine continues its execution using a detached context (e.g., via `context.WithoutCancel`).
2. **[Edge Case]** GIVEN a detached context WHEN the process receives a termination signal THEN the background goroutine should still be able to receive a signal for graceful shutdown if the detached context still respects process-level signals.

### [REQ-2] Secure Subprocess Execution Sanitization
**User Story:** AS A Security Auditor, I WANT TO prevent shell injection attacks in command execution, SO THAT variable arguments cannot be used to execute unintended system commands.

**Scenarios:**
1. **[Happy Path]** GIVEN a command executed with variable arguments WHEN the system runs the command THEN it MUST use a base command and an explicit slice of arguments, avoiding raw string concatenation in a shell context.
2. **[Audit Trail]** GIVEN a safe but dynamic command execution (e.g., calling `git` or `go`) WHEN `gosec` flags it as G204 THEN the code MUST include a `#nosec G204` annotation with a clear comment justifying why the specific call is safe.

### [REQ-3] Restrictive File System Permissions
**User Story:** AS A User, I WANT my local state and configuration to be inaccessible to other users on the same system, SO THAT sensitive project data is protected.

**Scenarios:**
1. **[Happy Path]** GIVEN the application creates or manages its internal state/data directories WHEN the directory is created THEN its permissions MUST be restricted to `0700` (Owner: RWX, Group: ---, Others: ---).
2. **[Binaries]** GIVEN the build process generates executable binaries WHEN the binary is created THEN its permissions SHOULD be `0755` (Owner: RWX, Group: R-X, Others: R-X) to allow execution while maintaining security.

### [REQ-4] Automated Security Enforcement in Makefile
**User Story:** AS A DevOps Engineer, I WANT the CI/CD pipeline to fail if security vulnerabilities are detected, SO THAT no insecure code is merged into the main branch.

**Scenarios:**
1. **[Happy Path]** GIVEN a Makefile with a `security` target WHEN `gosec` finds any vulnerability THEN the `make security` command MUST exit with a non-zero status code.
2. **[No Vulnerabilities]** GIVEN a codebase with no security issues WHEN `make security` is run THEN the command MUST exit with status code 0 and provide a clean report.

## 4. Business Invariants
1. ALL file system state directories created by the tool MUST strictly follow the `0700` permission rule.
2. NO dynamic command execution is permitted without either being fully sanitized (base + args) or having an audited and justified `#nosec` annotation.
3. The Build/CI process is CONSIDERED BROKEN if the `security` scan target fails.
