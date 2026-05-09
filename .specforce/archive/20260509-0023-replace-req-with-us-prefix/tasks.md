---
slug: 20260509-0023-replace-req-with-us-prefix
lens: Backend-heavy
---

# Implementation Roadmap: Replace REQ with US Prefix

## 1. Execution Strategy
The implementation follows a "Templates-First" approach. We will first update the embedded YAML templates for requirements and tasks to use the new `[US-x]` (User Story) prefix. This change aligns the generated documentation with industry standards and the explicit "User Story" field already present in our templates. Once the templates are updated, we will align the unit tests in the Go core to use the new prefix in their test cases to ensure future-proof consistency, even though the validator itself remains prefix-agnostic.

## 2. Tasks

### Phase 1: Template Standardization

- [x] T1.1: [CODE] Update Requirements Template Prefix
**Target:** `src/internal/agent/artifacts/spec/requirements.yaml`
**Context:** [US-1]

**Action Steps:**
- Replace the `[REQ-x]` pattern with `[US-x]` in the `instruction` section (Rule 5: "Every functional requirement `[US-x]` block...").
- Update the `template` section to replace `### [REQ-1]` with `### [US-1]`.
- Update the `template` section to replace `### [REQ-2]` with `### [US-2]`.

**Verification (TDD):**
Run `grep "\[US-" src/internal/agent/artifacts/spec/requirements.yaml` and verify it returns matches in the template and instruction sections. Run `grep "\[REQ-" src/internal/agent/artifacts/spec/requirements.yaml` and verify it returns no matches in those specific sections.

- [x] T1.2: [CODE] Update Tasks Template Prefix
**Target:** `src/internal/agent/artifacts/spec/tasks.yaml`
**Context:** [US-2]

**Action Steps:**
- Navigate to the `template` section in `tasks.yaml`.
- Replace all occurrences of `**Context:** [REQ-X]` with `**Context:** [US-X]`.

**Verification (TDD):**
Run `grep "\[US-X\]" src/internal/agent/artifacts/spec/tasks.yaml` and verify that the placeholder in the template now uses the US prefix.

### Phase 2: Test Alignment

- [x] T2.1: [CODE] Align Task Validation Tests with New Prefix
**Target:** `src/internal/spec/tasks_validation_test.go`
**Context:** [US-1], [US-2]

**Action Steps:**
- In `getHappyPathCases`, update the `**Context:**` field from `REQ-1` to `US-1`.
- In `getHierarchyErrorCases`, update all `**Context:**` fields to use `US-1`.
- In `getFieldAndPhaseErrorCases`, update all `**Context:**` fields to use `US-1`.

**Verification (TDD):**
Run `go test ./src/internal/spec/...` and verify all tests pass.

- [x] T2.2: [CLI] Global Regression Suite
**Target:** `Global Scope`
**Context:** [US-1], [US-2]

**Action Steps:**
- Execute the full project test suite to ensure that changing the prefix in templates and tests did not break any unexpected dependencies.

**Verification (TDD):**
Run `go test ./...` and ensure the output ends with `PASS` and a 0 exit code.

### Phase 3: Verification of Exclusion

- [x] T3.1: [CLI] Verify Bugfix Template Integrity
**Target:** `src/internal/agent/artifacts/spec/bug-requirements.yaml`
**Context:** [US-1]

**Action Steps:**
- Inspect the bugfix requirements template to ensure it still uses the `[FIX-x]` prefix.

**Verification (TDD):**
Run `grep "\[FIX-" src/internal/agent/artifacts/spec/bug-requirements.yaml` and verify it still returns matches, and `grep "\[US-"` returns no matches in that file.
