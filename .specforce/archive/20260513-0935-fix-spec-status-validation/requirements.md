---
slug: 20260513-0935-fix-spec-status-validation
lens: Bugfix
---

# Bugfix: Spec Status Validation & Redundant Rewriting

## 1. Issue Description
AI agents are redundantly rewriting spec artifacts because `spec status` reports `is_valid: false` prematurely or fails to detect malformed files correctly. Additionally, `tasks.md` validation is inconsistently applied or fails to catch major structural errors, leading to "false positives" where a broken file is reported as valid.

## 2. Evidence & Observations
- **Symptom 1:** `spec status` returns `is_valid: true` even if `tasks.md` uses the wrong header levels (e.g., `## Phase` instead of `### Phase`) or is missing mandatory fields.
- **Symptom 2:** Agents enter an implementation-fix loop because `GetStatus` performs deep validation of `tasks.md` before all other artifacts are generated, causing noise during the batch planning phase.
- **Trace:** 
    - `src/internal/spec/tasks.go`: `ValidateTasks` only checks sequentiality and internal blocks if headers match its regex, but doesn't enforce that *any* header was matched.
    - `src/internal/spec/status.go`: `GetStatus` calls `ValidateTasks` eagerly for any existing `tasks.md` regardless of total spec progress.

## 3. Reproduction Steps
1. Create a `tasks.md` with `## Phase 1` (incorrect level) and `### Task 1.1` (incorrect format).
2. Run `specforce spec status <slug> --json`.
3. **Observed Outcome:** `is_valid: true` and 0 `validation_errors` for the tasks artifact.
4. **Expected Outcome:** `is_valid: false` and error "No valid Phase found".

## 4. Root Cause Analysis (RCA)
1. **Silent Failure in Validator:** `ValidateTasks` doesn't check if any phases were processed. If the regexes don't match, the state remains empty but "valid".
2. **Precocious Validation:** Deep validation in `GetStatus` is performed as soon as a file exists. In SDD, artifacts are often written in a sequence that might trigger validation errors for partially completed specs (e.g., tasks existing but requirements/design missing or being updated).

## 5. Acceptance Criteria (Regression Tests)

### [FIX-1] Robust Task Validation
**Scenario: [Regression]**
GIVEN a `tasks.md` file without any valid `### Phase N:` headers
WHEN `ValidateTasks` is executed
THEN it MUST return an error: "No valid Phase (### Phase N: Name) found in tasks.md".

### [FIX-2] Conditional Deep Validation
**Scenario: [Workflow Improvement]**
GIVEN a specification with 3 total artifacts
WHEN only 1 artifact (`tasks.md`) exists
THEN `GetStatus` MUST NOT perform deep `tasks.md` validation (it should only check existence).
AND WHEN all 3 artifacts exist
THEN `GetStatus` MUST perform deep validation and report any errors.

### [FIX-3] Bug-Type Task Identification
**Scenario: [Consistency]**
GIVEN a spec of `type: bug`
WHEN `spec status` is executed
THEN the `tasks.md` artifact (identified as `bug-tasks` or `tasks`) MUST be correctly validated.

## 6. Technical Constraints (NFR)
- **[Safety]:** Ensure backward compatibility for specs without `spec.yaml` (default to `feature`).
- **[Performance]:** No significant increase in status reporting latency.
