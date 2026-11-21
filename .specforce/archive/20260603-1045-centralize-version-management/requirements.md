# Requirements: Centralized Version Management

**Slug:** 20260603-1045-centralize-version-management
**Feature Name:** Centralized Version Management
**Status:** DRAFT
**Persona:** Maintainers & Users

## 1. Problem Statement & Value

**Problem:** The version of the Specforce CLI is currently manually maintained in multiple disparate files (`main.go`, `registry.go`, `constants.go`, `logo.go`, `package.json`). This duplication leads to "version drift," where the binary reports a version different from the distribution package, causing confusion for users and errors during the release process.

**Value:** Establishing a Single Source of Truth (SSoT) ensures that the version reported by the CLI is always consistent with the NPM package, reduces manual maintenance overhead, and eliminates release-blocking inconsistencies.

---

## 2. Success Metrics

*   **[Business Metric]:** 100% reduction in "mismatched version" release failures.
*   **[Performance]:** Automated version extraction and propagation SHALL take less than 500ms during the release process.
*   **[UX Efficiency]:** The CLI `--version` output must match the `package.json` version string exactly, with zero manual intervention required in Go source files after the initial configuration.

---

## 3. Key Entities & Domain Logic

*   **Single Source of Truth (SSoT):** The `version` field in the root `package.json` file.
*   **Version Constant:** The internal Go variable (`core.Version`) that holds the version string used for display and telemetry.
*   **Synchronization Script:** The mechanism that propagates the SSoT to the Go source code.
*   **Build Injection:** The process of overriding the Version Constant during the compilation phase using build flags.

---

## 4. Functional Requirements

### US.1: Authoritative Versioning via Package.json
The version defined in `package.json` must be the authoritative version for the entire project.

*   **AC 1:** GIVEN a version update in `package.json`, WHEN the build process is triggered, THEN the resulting binary must report that exact version.
*   **AC 2:** GIVEN multiple files containing version strings, WHEN a new release is initiated, THEN only `package.json` should require a manual update.
*   **Edge Case:** GIVEN an invalid semver string in `package.json`, WHEN the synchronization script runs, THEN the process must fail with a descriptive error and prevent the build.
*   **UI/UX Specifics:** None (CLI Backend logic).
*   **Technical Constraints (NFR):**
    *   **[Consistency]:** The version string must follow SemVer 2.0.0 standards.
    *   **[Performance]:** Automated version extraction must not add perceptible delay to the local development build cycle.

### US.2: Automated Version Propagation to Go Source
The version information must be automatically synchronized from `package.json` to the Go source code.

*   **AC 1:** GIVEN a change in `package.json` version, WHEN the standard package manager version hook is triggered, THEN a synchronization script must update the internal Go constants.
*   **AC 2:** GIVEN the synchronization script runs, THEN it must only modify the specific version constant and not touch any other logic in the target Go files.
*   **Edge Case:** GIVEN the target Go file is missing or read-only, WHEN the synchronization script runs, THEN it must exit with a non-zero status and alert the maintainer.
*   **UI/UX Specifics:** None.
*   **Technical Constraints (NFR):**
    *   **[Reliability]:** The synchronization must be idempotent; running it multiple times with the same version should result in no further changes.
    *   **[Performance]:** File I/O operations for synchronization must be atomic to prevent partial writes.

### US.3: Build-time Version Injection Consistency
The build system must support injecting the version information during the compilation phase to ensure consistency.

*   **AC 1:** GIVEN the `Makefile` is used to build the project, WHEN the build command is executed, THEN it must extract the version from `package.json` and inject it into the binary using compiler flags.
*   **AC 2:** GIVEN a binary built via the `Makefile`, WHEN running `specforce --version`, THEN it must display the version injected at build time.
*   **Edge Case:** GIVEN a direct `go build` execution (bypassing the Makefile), WHEN the binary is run, THEN it should fallback to the last synchronized constant value in the source code.
*   **UI/UX Specifics:** The version output format must remain `specforce-kit version x.y.z`.
*   **Technical Constraints (NFR):**
    *   **[Performance]:** Build-time injection must not increase compilation time by more than 1%.

---

## 5. Scope Containment (Anti-Goals)

*   **NO** auto-update mechanism for the CLI itself.
*   **NO** new CLI command to manage versions.
*   **NO** modification of how external dependencies are versioned.
*   **NO** automation of the Git tagging or pushing process (this remains a manual or separate CI step).
*   **NO** change to the user-facing version display format beyond ensuring data consistency.
