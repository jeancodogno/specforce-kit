---
slug: 20260524-2126-spec-hardening-and-self-containment
lens: Backend-heavy
---

# Implementation Roadmap: Spec Hardening & Self-Containment

## 1. Execution Strategy
- **Gravity Order:** Template Hardening (YAML) -> Logic Hardening (Go) -> Validation & Regression.
- **TDD Protocol:** Each logic change in Phase 2 MUST be preceded by a failing test case in `src/internal/spec/tasks_validation_test.go`.

## 2. Tasks

### Phase 1: Template Hardening

- [x] T1.1: [CODE] Harden Requirements Template
**Target:** `src/internal/agent/artifacts/spec/requirements.yaml`
**Context:** [US-1, US-2]

**Action Steps:**
- Update `instruction` to mandate at least one edge case per US.
- Update `instruction` to enforce localized `Technical Constraints (NFR)` for every US.
- Add mandatory `**[Performance]**` tag requirement to instructions.
- Add logic to instructions for conditional removal of UI/UX sections for non-UI lenses.

**Acceptance Check:**
- Run `specforce spec artifact feature-requirements --json` and verify the new instructions are present.

- [x] T1.2: [CODE] Harden Design Template
**Target:** `src/internal/agent/artifacts/spec/design.yaml`
**Context:** [US-1]

**Action Steps:**
- Update `instruction` to explicitly prohibit `TBD` and placeholders.
- Mandate Mermaid diagrams for architecture/data flow.
- Enforce explicit file paths in the `File & Component Inventory`.

**Acceptance Check:**
- Run `specforce spec artifact feature-design --json` and verify the new instructions are present.

- [x] T1.3: [CODE] Harden Tasks Template
**Target:** `src/internal/agent/artifacts/spec/tasks.yaml`
**Context:** [US-1, US-2]

**Action Steps:**
- Update `instruction` to mandate high-density action steps (minimum 3 per task).
- Enforce mandatory `**Context:**` linking to requirements.
- Mandate concrete technical directives instead of passive descriptions.

**Acceptance Check:**
- Run `specforce spec artifact feature-tasks --json` and verify the new instructions are present.

### Phase 2: Logic Hardening

- [x] T2.1: [TEST] Create regression test for task density
**Target:** `src/internal/spec/tasks_validation_test.go`
**Context:** [US-3]

**Action Steps:**
- Add a test case that provides a `tasks.md` with a task containing only 1 action step.
- Assert that `ValidateTasks` returns an error regarding action density.

**Acceptance Check:**
- Run `go test ./src/internal/spec/...` and verify the test fails (RED).

- [x] T2.2: [CODE] Implement task density and sequence validation
**Target:** `src/internal/spec/tasks.go`
**Context:** [US-3]

**Action Steps:**
- Update `taskBlock` struct with `actionItemsCount int`.
- Update `updateTaskBlockState` to increment `actionItemsCount` for each checklist item in action steps.
- Update `validateLastTask` to return an error if `actionItemsCount < 2`.
- Add sequence validation for Phase IDs and Task IDs.

**Acceptance Check:**
- Run `go test ./src/internal/spec/...` and verify the density test passes (GREEN).

### Phase 3: Final Verification

- [x] T3.1: [TEST] Holistic Verification
**Target:** `Global Scope`
**Context:** [US-1, US-2, US-3]

**Action Steps:**
- Initialize a dummy spec and generate all artifacts using the new templates.
- Run `specforce spec status <dummy-slug>` and verify it passes structural and density validation.

**Acceptance Check:**
- `specforce spec status` returns `is_valid: true` and `progress: 100` for a high-fidelity spec.
