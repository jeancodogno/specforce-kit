---
slug: 20260910-1415-fix-implementation-status-tiered-sizing
lens: Bugfix
---

# Bugfix: Tiered Sizing Invariance in Implementation Status and Artifact Dependencies

## 1. Issue Description
When a specification is initialized with a lightweight sizing tier such as `small` (which requires only task breakdown) or `medium` (which requires requirements and tasks), the implementation readiness check incorrectly reports missing artifacts from the three-tier documentation triad (`requirements.md`, `design.md`, `tasks.md`). As a result, the implementation pipeline is falsely marked as blocked. Autonomous agents attempting to execute the implementation protocol react by promoting the specification to `large` via `spec resize` and synthesizing unnecessary documentation artifacts, breaking the tiered sizing contract established during planning. Furthermore, specification status tracking incorrectly flags artifacts as blocked when their upstream dependencies are not mandated by the active specification size.

## 2. Evidence & Observations
- **Symptom:** Executing the implementation status query on a specification configured with size `small` returns status `blocked` and enumerates `requirements.md` and `design.md` under missing artifacts.
- **Agent Reaction:** AI agents reading the blocked status invoke `specforce spec resize <slug> --size large` and generate full triad documents, overwriting the user-selected `size: small` in specification metadata.
- **Trace:**
  - Artifact completeness check during implementation inspection evaluates a fixed list of triad filenames rather than the dynamic artifact matrix defined by the specification size.
  - Specification artifact status evaluates upstream dependency existence against global registry definitions rather than the filtered active artifact set for the specification size.

## 3. Reproduction Steps
1. Initialize a new specification with type `bug` and size `small`:
   `specforce spec init mondial-service-type-mapping --type bug --size small`
2. Create the mandated `tasks.md` roadmap according to the `small` sizing tier rules.
3. Verify specification progress via status query:
   `specforce spec status mondial-service-type-mapping --json`
   (Returns 100% progress and valid specification state).
4. Query implementation readiness:
   `specforce implementation status mondial-service-type-mapping --json`
5. Observe output:
   Status is reported as `blocked` with `missing_artifacts: ["requirements.md", "design.md"]`, and tasks in status `blocked` in the specification artifact view.

## 4. Root Cause Analysis (RCA)
- The implementation verification routine unconditionally checks for the existence of all three standard triad documents (`requirements.md`, `design.md`, `tasks.md`), regardless of the specification's configured size. It fails to consult the size-aware artifact registry to determine which artifacts are actually mandatory.
- The artifact status evaluation checks dependency existence across all registered artifact templates. When an optional upstream artifact (such as `design.md` for `small` or `medium` tiers) is not part of the active size requirement, the downstream artifact is mistakenly marked as blocked.

## 5. Acceptance Criteria (Regression Tests)

### [FIX-1] Size-Aware Implementation Readiness Verification
**Scenario: [Regression - Small Spec Implementation Status]**
GIVEN a specification initialized or configured with size `small`
AND the specification contains a valid `tasks.md` artifact
WHEN the implementation readiness query is executed
THEN the status evaluates to `ready`
AND `missing_artifacts` is empty
AND the specification is not blocked by absence of `requirements.md` or `design.md`.

**Scenario: [Regression - Medium Spec Implementation Status]**
GIVEN a specification initialized or configured with size `medium`
AND the specification contains valid `requirements.md` and `tasks.md` artifacts
WHEN the implementation readiness query is executed
THEN the status evaluates to `ready`
AND `missing_artifacts` does not include `design.md`.

**Scenario: [Regression - Large Spec Implementation Status Preserves Triad Enforcement]**
GIVEN a specification configured with size `large` or `complex`
AND `design.md` or `requirements.md` is absent
WHEN the implementation readiness query is executed
THEN the status evaluates to `blocked`
AND `missing_artifacts` accurately reports the absent triad artifacts.

### [FIX-2] Size-Aware Artifact Dependency Evaluation in Spec Status
**Scenario: [Regression - Small Spec Task Dependency]**
GIVEN a specification configured with size `small`
AND `design.md` is not present on disk
WHEN the specification status query is executed
THEN the `tasks.md` artifact is not marked as blocked by missing `design.md`
AND the overall specification progress evaluates to 100% upon completing `tasks.md`.

**Scenario: [Regression - Medium Spec Task Dependency]**
GIVEN a specification configured with size `medium`
AND `requirements.md` is present but `design.md` is absent
WHEN the specification status query is executed
THEN the `tasks.md` artifact is not blocked by missing `design.md`
AND `tasks.md` is only blocked if `requirements.md` is absent.

## 6. Technical Constraints (NFR)
- **[Performance]:** Implementation readiness and specification status queries must evaluate in under 30ms.
- **[Safety]:** Verification must remain purely non-destructive and must not alter metadata or files on disk.
- **[Observability]:** JSON outputs must accurately reflect required missing artifacts matching the active specification size.
