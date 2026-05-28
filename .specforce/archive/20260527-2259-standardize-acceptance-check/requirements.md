---
slug: 20260527-2259-standardize-acceptance-check
lens: Migration
---

# Feature: Standardize Acceptance Check

## 1. Context & Value
Developers currently use "Verification (TDD)" as a label for task verification in `tasks.md`. This implies a full Test-Driven Development methodology which might not always be applicable, leading to methodology confusion. Standardizing on "Acceptance Check" provides a clear, methodology-agnostic label for verifying that a task meets its definition of done.

## 2. Out of Scope (Anti-Goals)
- Supporting multiple verification labels (e.g., keeping "Verification (TDD)" as an alias).
- Modifying the verification logic itself (only the label/keyword is changing).
- Updating non-spec files that might mention TDD in a different context (e.g., general testing strategy docs).

## 3. Acceptance Criteria (BDD)

### [US-1] Standardized Verification Label
**User Story:** AS A developer, I WANT TO use 'Acceptance Check' as the only valid label for task verification, SO THAT the intent of the verification step is clear and unambiguous.

**Scenarios:**
1. **[Happy Path] Successful task validation with new label**
   GIVEN a `tasks.md` file containing a task with an `Acceptance Check:` section
   WHEN the `specforce status` command is executed
   THEN the parser correctly identifies the verification steps and marks the task as valid.

2. **[Edge Case] Failed task validation with old label**
   GIVEN a `tasks.md` file containing a task with a `Verification (TDD):` section
   WHEN the `specforce status` command is executed
   THEN the parser fails to recognize the section and returns a validation error indicating that 'Acceptance Check' is required.

**UI/UX Specifics:**
- **View/Component:** CLI Output / TUI Status View.
- **Feedback Logic:** Display a clear error message in **Error Red (#FF5F5F)** when the old label is used, suggesting the correct one in **Mint Green (#00FA9A)**.
- **Keybindings:** N/A (CLI/TUI purely informational).

**Technical Constraints (NFR):**
- **[Performance]:** Parser validation MUST be < 50ms.
- **[Safety & Security]:** Migration of existing files MUST be idempotent and safe (no data loss).
- **[Integrity]:** The parser MUST strictly enforce the 'Acceptance Check' string (case-insensitive for flexibility but preferring the standardized casing).
- **[Observability]:** Log a warning or error when the old label is encountered via the `core.UI` interface.

### [US-2] Template and Codebase Update
**User Story:** AS A framework maintainer, I WANT TO update all internal templates and parser logic, SO THAT all new and existing specs follow the new standard.

**Scenarios:**
1. **[Happy Path] New spec generation uses correct label**
   GIVEN a developer initializes a new spec via `specforce spec init`
   WHEN the `tasks.md` template is rendered
   THEN the verification section is labeled as `Acceptance Check:`.

**Technical Constraints (NFR):**
- **[Performance]:** Template rendering MUST be instantaneous (< 10ms).
- **[Maintainability]:** The `src/internal/agent/artifacts/spec/tasks.yaml` template MUST be updated to reflect the change.

## 4. Business Invariants
- The string "Acceptance Check:" is the only recognized prefix for task verification steps in `tasks.md`.
- A task is considered invalid for implementation if it lacks an `Acceptance Check:` section or uses the deprecated `Verification (TDD):` label.

## 5. Global UI/UX Contract (TUI Ghost Protocol)
- **Density Posture:** Compact (Optimized for 80x24).
- **Signature Moves:** Mint Green highlights for valid states, Error Red for deprecated label detection.
- **State Behavior:** Validation errors must be prominent and instructive, guiding the user toward the fix.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Total parsing and validation time for a `tasks.md` file MUST be < 100ms.
- **[Reliability]:** The system MUST fail-fast with clear, actionable errors when the legacy label is detected.
- **[Maintainability]:** Update `src/internal/spec/tasks.go` to enforce the new label and ensure 80% test coverage on the updated logic.
