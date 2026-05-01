---
slug: security-hardening
lens: Backend-heavy
---

# Tasks: Security Hardening

## Phase 1: Environment & Tooling

#### T1.1: Integrate `govulncheck` into the local development environment
**State:** [FINISHED]
**Description:** Add a mechanism to install or verify the presence of `govulncheck` to enable vulnerability scanning of dependencies.
**Verification:** Run `govulncheck ./...` and ensure it executes correctly (finding vulnerabilities is expected at this stage).

#### T1.2: Implement the `make security` target in the Makefile
**State:** [FINISHED]
**Description:** Define a `security` target in the project Makefile that runs both `gosec` and `govulncheck`. Ensure it returns a non-zero exit code if any issues are found.
**Verification:** Run `make security`. The command should fail if the vulnerabilities described in the requirements still exist.

## Phase 2: Core Hardening

#### T2.1: Hardening Goroutine Context Propagation (G118)
**State:** [FINISHED]
**Description:** In `src/internal/upgrade/service.go`, update background goroutine spawning to use `context.WithoutCancel(ctx)` instead of `context.Background()`. This ensures background tasks are not prematurely canceled while still carrying context values.
**Verification:** Update/Create `src/internal/upgrade/service_test.go` to assert that background operations continue even after the parent context is canceled, but still have access to initial context values.

#### T2.2: Enforce Restrictive Directory Permissions (G301)
**State:** [FINISHED]
**Description:** In `src/internal/upgrade/state.go`, ensure all directory creation calls (e.g., `os.MkdirAll`) use `0700` permissions.
**Verification:** Run `src/internal/upgrade/state_test.go`. Add a test case that creates a directory and uses `os.Stat` to verify that the mode is exactly `drwx------`.

#### T2.3: Secure Binary File Permissions (G302)
**State:** [FINISHED]
**Description:** In `src/internal/upgrade/installer_binary.go`, ensure that installed binaries are set to `0755` permissions using `os.Chmod`. Add a `// #nosec G302` annotation with a justification for the executable bit.
**Verification:** Run `src/internal/upgrade/installer_binary_test.go`. Verify that the resulting binary file has `0755` permissions after installation.

#### T2.4: Sanitize and Audit Subprocess Execution (G204)
**State:** [FINISHED]
**Description:** In `src/internal/upgrade/installer_npm.go`, audit all `os/exec.Command` calls. Ensure they use a base command and an explicit slice of arguments. Add `// #nosec G204` where the command is dynamic but the arguments are sanitized or controlled by the application.
**Verification:** Run `make security`. The G204 warning for `installer_npm.go` should disappear or be acknowledged via audited annotation. Run `src/internal/upgrade/installer_npm_test.go` to ensure functionality remains intact.

## Phase 3: Automated Enforcement

#### T3.1: Enable Blocking Security Gate in CI/CD
**State:** [FINISHED]
**Description:** Update the GitHub Actions workflow (or relevant CI config) to include the `make security` command as a mandatory step in the pipeline.
**Verification:** Push a temporary branch with a known (safe) security annotation change and verify that the CI runs the security check.

#### T3.2: Final Security Verification Scan
**State:** [FINISHED]
**Description:** Perform a final full project scan using `make security` to ensure all identified vulnerabilities (G118, G204, G301, G302) are resolved or properly audited.
**Verification:** `make security` returns exit code 0 with no unexpected findings.
