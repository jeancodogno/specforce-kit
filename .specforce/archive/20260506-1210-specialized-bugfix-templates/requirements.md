---
slug: 20260506-1210-specialized-bugfix-templates
lens: Backend-heavy
---

# Feature: Specialized Bugfix Templates

## 1. Context & Value
Currently, Specforce uses a single set of requirements, design, and task templates for all technical activities. While effective for new features, this "one-size-fits-all" approach forces developers to use feature-centric structures (like User Stories and UI Contracts) for bugfixes where root cause analysis and reproduction steps are more relevant. This feature introduces a specialized metadata layer and category-aware templates to provide high-signal documentation for every type of work.

## 2. Out of Scope (Anti-Goals)
- Do not implement a `status` field in the `spec.yaml` metadata (handled by archival logic).
- Do not modify the existing archival directory structure.
- Do not implement a UI for selecting templates; this remains CLI-driven.

## 3. Acceptance Criteria (BDD)

### [REQ-1] YAML-Based Specification Metadata
**User Story:** AS A Developer, I WANT TO have a dedicated metadata file for each spec, SO THAT I can store machine-readable attributes like the task type.

**Scenarios:**
1. **[Happy Path]** GIVEN a new specification is initialized WHEN the directory is created THEN a `spec.yaml` file MUST be generated in that directory containing the `slug`, `name`, and `type`.
2. **[Edge Case]** GIVEN an existing `spec.yaml` WHEN the user manually edits the `type` field THEN subsequent Specforce commands MUST respect the new type for artifact discovery.

**Technical Constraints (NFR):**
- **[Performance]:** YAML parsing must take < 10ms.
- **[Integrity]:** The file MUST NOT contain a `status` field to prevent drift with the archival state.

### [REQ-2] Specialized Bugfix Artifacts
**User Story:** AS A Debugging Specialist, I WANT TO use templates focused on root causes and reproduction, SO THAT I don't waste time filling out irrelevant feature requirements.

**Scenarios:**
1. **[Happy Path]** GIVEN a specification with `type: bug` WHEN the agent generates requirements THEN it MUST use the `bug-requirements` template instead of the standard one.
2. **[Negative Case]** GIVEN a registry with missing specialized artifacts WHEN a bug spec is processed THEN the system MUST fallback to the standard artifacts and log a warning.

**Technical Constraints (NFR):**
- **[Performance]:** Template selection logic must not introduce visible latency in the CLI.
- **[Reliability]:** Fail-safe fallback to standard templates if specialized ones are missing.

### [REQ-3] Type-Aware CLI Commands
**User Story:** AS A CLI User, I WANT TO specify the type of work during initialization, SO THAT the system is configured correctly from the start.

**Scenarios:**
1. **[Happy Path - Init]** GIVEN the command `specforce spec init <slug> --type bug` WHEN executed THEN the generated `spec.yaml` MUST have `type: bug`.
2. **[Happy Path - Status]** GIVEN an active spec WHEN `specforce spec status <slug>` is run THEN the completeness report MUST only show artifacts relevant to the spec's `type`.

**Technical Constraints (NFR):**
- **[Observability]:** The CLI output MUST clearly indicate the active specification type.

## 4. Business Invariants
- A specification MUST always have a valid `type` (defaulting to `feature` if unspecified).
- The `slug` in `spec.yaml` MUST match the directory name.

## 5. Global Non-Functional Requirements (NFRs)
- **[Reliability]:** The system must remain backward compatible with specs that lack a `spec.yaml` file (treating them as `type: feature`).
- **[Security]:** Path sanitization for `spec.yaml` access to prevent directory traversal.
- **[Maintainability]:** Centralize artifact selection logic in the `Registry` or `Service` layers.
