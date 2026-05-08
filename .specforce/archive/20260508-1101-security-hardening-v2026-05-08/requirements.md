---
slug: 20260508-1101-security-hardening-v2026-05-08
lens: Backend-heavy
---

# Feature: Security Hardening (v2026-05-08)

## 1. Context & Value
The project requires urgent security hardening to address critical vulnerabilities in the Go standard library and improve the overall security posture of the framework's internal file management. Discovery has identified that Go 1.26.2 contains critical vulnerabilities in `net` and `net/http` packages (CVEs addressed in 1.26.3). Additionally, sensitive metadata files are currently being created with overly permissive access (`0644`), which must be restricted to owner-only (`0600`) per the project's Engineering Constitution. Finally, global security scan exclusions for path traversal (G304) must be replaced with granular source-level justifications to ensure continuous validation while avoiding false positives.

## 2. Out of Scope (Anti-Goals)
- Implementation of new Authentication (AuthN) or Authorization (AuthZ) systems.
- Major refactoring of the file system logic beyond permission tightening.
- Upgrading Go beyond 1.26.3.
- Modifying the `core.SecurePath` logic itself (only its application and scanning).

## 3. Acceptance Criteria (BDD)

### [REQ-1] Go Toolchain Upgrade to 1.26.3
**User Story:** AS A security-conscious developer, I WANT the project to use a patched Go version, SO THAT critical network-related vulnerabilities are mitigated.

**Scenarios:**
1. **[Happy Path]** GIVEN the current `go.mod` and CI configuration WHEN the update is applied THEN the `go.mod` file MUST specify `go 1.26.3` AND all CI workflows in `.github/workflows/` MUST execute using Go 1.26.3.
2. **[Edge Case]** GIVEN a local environment running a Go version lower than 1.26.3 WHEN `go build` or `make build` is executed THEN the build MUST fail or emit a warning if the toolchain version does not meet the minimum security baseline.

**Technical Constraints (NFR):**
- **[Performance]:** Build times must not increase by more than 5% due to the toolchain upgrade.
- **[Reliability]:** 100% of existing unit and integration tests must pass on the new Go version.

### [REQ-2] Strict File Permissions for Sensitive Metadata
**User Story:** AS A system administrator, I WANT sensitive application files to have restricted access, SO THAT unauthorized users on the same system cannot read upgrade states or specification metadata.

**Scenarios:**
1. **[Happy Path]** GIVEN a new specification is initialized or an upgrade state is saved WHEN the file is created on disk THEN its permissions MUST be set to `0600` (Read/Write for owner only).
2. **[Edge Case]** GIVEN an existing file with `0644` permissions WHEN the system performs an update or access to that file (e.g., during `specforce upgrade` or `specforce spec status`) THEN it MUST automatically correct the permissions to `0600`.

**Technical Constraints (NFR):**
- **[Performance]:** Permission checks and enforcement MUST add less than 10ms overhead to file operations.
- **[Security]:** All directory creations MUST use `0750` as per `engineering.md` hardening rules.

### [REQ-3] Granular Gosec G304 Resolution
**User Story:** AS A developer, I WANT security scan warnings to be handled at the source level, SO THAT we maintain high visibility into potential path traversal risks while justifying safe usage of `core.SecurePath`.

**Scenarios:**
1. **[Happy Path]** GIVEN the use of `os.Open` or similar on a path validated by `core.SecurePath` WHEN running `gosec` THEN the scan MUST pass without global exclusions for G304 in the `Makefile` OR `.golangci.yml` AND each specific site MUST contain a `#nosec G304` comment with a clear justification.
2. **[Edge Case]** GIVEN a new instance of file opening WITHOUT `core.SecurePath` validation WHEN running `gosec` THEN the scan MUST fail, ensuring no unvalidated paths are introduced.

**Technical Constraints (NFR):**
- **[Performance]:** Refined scanning MUST NOT increase CI linting time by more than 10 seconds.
- **[Maintainability]:** Justifications in `#nosec` comments MUST be clear and link to the project's security principles.

## 4. Business Invariants
- **Security by Default:** All new file creations by the framework MUST default to the most restrictive permissions necessary (`0600` for files, `0750` for directories).
- **Zero Regression:** Security hardening MUST NOT break existing CLI commands, TUI functionality, or cross-agent artifact mapping.
- **Auditability:** Security scan results from `gosec` and `govulncheck` MUST be clean (zero critical/high findings) and reproducible in CI.

## 5. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Security enforcement MUST be optimized to ensure O(1) or O(N) complexity for file batch operations, avoiding redundant disk syncs.
- **[Security]:** 100% resolution of critical and high vulnerabilities reported by security scanning tools.
- **[Maintainability]:** Use of standard Go library features for permission management to ensure consistent behavior across Unix-like platforms.
