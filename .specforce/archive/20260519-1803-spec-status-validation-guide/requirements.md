---
slug: 20260519-1803-spec-status-validation-guide
lens: Backend-heavy
---

# Feature: Validation Guide Golden Model

## 1. Context & Value
When `spec status` fails validation for `tasks.md`, users only receive a list of errors. Adding a static "Golden Model" guide will provide an immediate reference on how a valid phase and task should look.

## 2. Out of Scope (Anti-Goals)
- Do not implement dynamic error suggestions or intelligent hints based on the specific error line.
- Do not modify validation logic for other artifacts besides `tasks.md`.

## 3. Acceptance Criteria (BDD)

### [US-1] Include Golden Model in JSON output
**User Story:** AS A developer, I WANT TO see a valid tasks.md example when validation fails, SO THAT I can quickly fix my spec format.

**Scenarios:**
1. **[Happy Path]** GIVEN a `tasks.md` with validation errors WHEN I run `specforce spec status <slug> --json` THEN the JSON output includes a `validation_guide` field containing a valid "Golden Model" example.
2. **[Edge Case]** GIVEN a valid `tasks.md` WHEN I run `specforce spec status <slug> --json` THEN the `validation_guide` field is either empty or omitted from the JSON output.

**UI/UX Specifics:**
- **View/Component:** JSON Output structure for `spec status`.
- **Feedback Logic:** N/A (JSON format only for now).
- **Keybindings:** N/A.

**Technical Constraints (NFR):**
- **[Performance]:** Adding the static guide must have zero impact on the validation execution time.
- **[Integrity]:** The guide MUST reflect the exact syntax expected by the `ValidateTasks` parser.

## 4. Business Invariants
- The system must fail-fast with clear errors when artifacts break the required structure.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Latency < 50ms for spec status check.
- **[Reliability]:** Valid JSON output even on structural errors.
- **[Maintainability]:** 80% test coverage minimum.