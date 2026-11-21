---
slug: 20261121-1826-ci-green-for-v1
lens: Integration
---

# Feature: Green Continuous Integration for the 1.0.0 Release

## 1. Context & Value
The first continuous integration run on the `1.x` line failed in three places: the security job (a false-positive credential warning and twelve known vulnerabilities in the toolchain's standard library that the released binary would ship), one macOS test and 61 Windows tests. Two of the Windows failures are real defects that Windows users hit today (archiving instructions lose their protocol); the rest are tests that do not isolate user state or assume Unix behaviour. The stable 1.0.0 must not be released until every CI job is green on Linux, macOS and Windows.

## 2. Out of Scope (Anti-Goals)
- Do not change any behaviour on Linux or macOS other than the toolchain version.
- Do not normalize the path separator of paths the CLI prints or returns to agents; native separators stay (`\` on Windows).
- Do not add new CI jobs, platforms or test tooling; do not relax the coverage gate or the mutation gate.
- Do not tag, publish or release anything; only commits and pushes to the `1.x` branch are allowed for CI iterations.
- Do not fix unrelated findings (terminal queries on standard output, slash-style wording in skills).

## 3. Acceptance Criteria (BDD)

### [US-1] Secure toolchain and clean security scan
**User Story:** AS A Specforce user, I WANT released binaries built with a toolchain that has no known standard-library vulnerabilities reachable from Specforce, and a security scan without false alarms, SO THAT the 1.0.0 I install is not exposed to already-fixed network and TLS flaws.

**Scenarios:**
1. **[Happy Path]** GIVEN the repository WHEN the vulnerability scan runs THEN it reports no vulnerability reachable from Specforce code.
2. **[Happy Path]** GIVEN the repository WHEN the static security scan runs THEN it reports zero issues, and the one accepted false positive (a pattern for specification identifiers mistaken for a credential) carries a written justification next to it.
3. **[Happy Path]** GIVEN continuous integration and the release run WHEN they build or test THEN they use the same fixed toolchain version as the module definition (1.26.11).
4. **[Edge Case]** GIVEN a developer builds from source with an older toolchain WHEN the build starts THEN the toolchain required by the module definition is used (downloaded automatically when allowed) instead of silently building with the vulnerable one.

**Technical Constraints (NFR):**
- **[Performance]:** No extra CI job; the existing security job must pass in its current time budget.
- **[Safety & Security]:** Exactly one suppression, scoped to the one line, with its reason.
- **[Integrity]:** Module definition, CI and release workflows name the same toolchain version.

### [US-2] Archiving instructions and skill metadata work on Windows
**User Story:** AS A developer using Specforce on Windows, I WANT the archiving instructions and the skill metadata to be read from the embedded kit exactly as on Linux and macOS, SO THAT my agent receives the full archiving protocol.

**Scenarios:**
1. **[Happy Path]** GIVEN a Windows machine WHEN `specforce archive instructions` runs THEN the output contains the embedded archiving protocol followed by the project-specific instructions, identical to the Linux output for the same project.
2. **[Happy Path]** GIVEN a Windows machine WHEN the kit's skills are scanned THEN each skill's name, description and version come from its manifest, as on Linux.
3. **[Edge Case]** GIVEN an instruction category that does not exist in the kit WHEN instructions are requested THEN only the project-specific instructions are returned on every platform (unchanged behaviour).
4. **[Edge Case]** GIVEN the installation of a tool's skills on Windows (the refresh scenario of `specforce init` / tool update) WHEN it runs THEN every expected file is created; if the cause of the current Windows failure is a product defect of the same kind (path separators on the embedded kit) it is fixed here, otherwise the work stops and the cause is reported before any further change.

**Technical Constraints (NFR):**
- **[Performance]:** No change in command latency.
- **[Integrity]:** Every lookup in the embedded kit uses the embedded file system's own path rules, on every platform.

### [US-3] Test suite is portable and isolated on every platform
**User Story:** AS A release maintainer, I WANT the unit tests to pass on Linux, macOS and Windows without touching the real user profile, SO THAT a red CI job always means a real problem.

**Scenarios:**
1. **[Happy Path]** GIVEN any test that touches per-user Specforce state WHEN it runs on any platform THEN it reads and writes only a temporary home (both the Unix and the Windows home variables point to it) and never the real profile.
2. **[Happy Path]** GIVEN tests that compare file paths WHEN they run on Windows or macOS THEN they compare against the platform's native separator and the fully resolved path (for example `/private/var` on macOS, long user-folder names on Windows).
3. **[Happy Path]** GIVEN tests of verification hooks WHEN they run on any platform THEN they use commands that exist on every platform where the tests run.
4. **[Edge Case]** GIVEN a test that checks file permission bits or a read-only directory WHEN it runs on Windows THEN it is skipped with the explicit reason that permission bits do not apply on Windows, and it still runs on Linux and macOS.
5. **[Edge Case]** GIVEN the configuration security test that needs an impossible project root WHEN it runs on any platform THEN it uses a root that cannot be created on that platform and still expects an error.
6. **[Edge Case]** GIVEN the complete unit test suite WHEN continuous integration runs on Linux, macOS and Windows THEN all three test jobs, lint, security, mutation quality and release tooling tests pass.

**Technical Constraints (NFR):**
- **[Performance]:** No test becomes slower than 1 second because of these changes.
- **[Integrity]:** No existing assertion is weakened on Linux or macOS; skips are limited to the permission-bit tests and only on Windows.
- **[Observability]:** Every Windows-only skip prints its reason in verbose test output.

### [US-4] The isolation rule is written down
**User Story:** AS A Specforce maintainer, I WANT the engineering rule on user-state isolation to name both home variables and the permission-bit exception, SO THAT new tests do not repeat this failure.

**Scenarios:**
1. **[Happy Path]** GIVEN the engineering constitution WHEN a contributor reads the user-state isolation rule THEN it requires pointing both `HOME` and `USERPROFILE` at a temporary directory and allows skipping only permission-bit assertions on Windows, with a stated reason.
2. **[Edge Case]** GIVEN a test that skips on Windows for any other reason WHEN it is reviewed THEN it violates the rule.

**Technical Constraints (NFR):**
- **[Performance]:** Documentation only.
- **[Integrity]:** PATCH amendment approved by the human lead (approved on 2026-11-21).

## 4. Business Invariants
- Linux and macOS behaviour of every command is unchanged.
- No test reads or writes the real user profile on any platform.
- The 1.0.0 is not tagged while any CI job on `1.x` is red.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** CI wall time stays within the current order (no new jobs).
- **[Reliability]:** Each fix is proven by the Windows/macOS CI job that failed before it.
- **[Security]:** Toolchain without reachable known vulnerabilities; one justified scanner suppression.
- **[Maintainability]:** Core coverage stays at or above 80%; no new dependencies.
