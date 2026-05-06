---
slug: 20260506-1024-granular-requirements-templates
lens: Balanced full-stack
---

# Implementation Roadmap: Granular UI/UX and NFR Sections in Requirements Template

## 1. Execution Strategy
- **Gravity Order:** Update Instructions -> Refactor Template Structure -> Final Verification. We first define the rules for the AI, then update the template it uses.

## 2. Tasks

### Phase 1: Update Logic and Instructions

- [x] T1.1: [CONFIG] Update Instructions for Localized Context
**Target:** `src/internal/agent/artifacts/spec/requirements.yaml`
**Context:** [REQ-1]

**Action Steps:**
- Modify the `instruction` section in `src/internal/agent/artifacts/spec/requirements.yaml` to mandate that every `[REQ-x]` block MUST include localized `UI/UX Specifics` and `Technical Constraints (NFR)`. 
- Mandate a flexible `Attribute: Value` format for NFRs, including mandatory `Performance`.

**Verification (TDD):**
`grep "Performance" src/internal/agent/artifacts/spec/requirements.yaml | grep "instruction"` should succeed.

- [x] T1.2: [CONFIG] Update Instructions for Global Sections
**Target:** `src/internal/agent/artifacts/spec/requirements.yaml`
**Context:** [REQ-2]

**Action Steps:**
- Update the `instruction` section to require the presence of `Global UI/UX Contract` and `Global Non-Functional Requirements (NFRs)` sections.

**Verification (TDD):**
`grep "Global UI/UX Contract" src/internal/agent/artifacts/spec/requirements.yaml | grep "instruction"` should succeed.

- [x] T1.3: [CONFIG] Refactor Conditional Omission Rule
**Target:** `src/internal/agent/artifacts/spec/requirements.yaml`
**Context:** [REQ-2]

**Action Steps:**
- Update the `Conditional Omission` rule in `instruction` to apply to the new `## 5. Global UI/UX Contract` section when lens is "Backend-heavy" or "Integration".

**Verification (TDD):**
Check that the rule now refers to the correct section name: `grep "Global UI/UX Contract" src/internal/agent/artifacts/spec/requirements.yaml` in the omission part.

### Phase 2: Refactor Template Structure

- [x] T2.1: [CONFIG] Add Localized UI/UX Block to Requirement Template
**Target:** `src/internal/agent/artifacts/spec/requirements.yaml`
**Context:** [REQ-1]

**Action Steps:**
- In the `template` section of `src/internal/agent/artifacts/spec/requirements.yaml`, insert the `**UI/UX Specifics:**` block (View/Component, Feedback Logic, Keybindings) into the `[REQ-x]` structure.

**Verification (TDD):**
`grep "UI/UX Specifics" src/internal/agent/artifacts/spec/requirements.yaml` in the template section should succeed.

- [x] T2.2: [CONFIG] Add Localized NFR Block to Requirement Template
**Target:** `src/internal/agent/artifacts/spec/requirements.yaml`
**Context:** [REQ-1]

**Action Steps:**
- In the `template` section, insert the `**Technical Constraints (NFR):**` block with flexible `Attribute: Value` pairs (Safety, Validation, Performance, Observability) into the `[REQ-x]` structure.

**Verification (TDD):**
`grep "Performance" src/internal/agent/artifacts/spec/requirements.yaml` in the template section should succeed.

- [x] T2.3: [CONFIG] Add Global UI/UX Contract Section
**Target:** `src/internal/agent/artifacts/spec/requirements.yaml`
**Context:** [REQ-2]

**Action Steps:**
- Add the `## 5. Global UI/UX Contract (TUI Ghost Protocol)` section to the bottom of the `template` in `src/internal/agent/artifacts/spec/requirements.yaml`.

**Verification (TDD):**
`grep "TUI Ghost Protocol" src/internal/agent/artifacts/spec/requirements.yaml` should succeed.

- [x] T2.4: [CONFIG] Add Global NFRs Section
**Target:** `src/internal/agent/artifacts/spec/requirements.yaml`
**Context:** [REQ-2]

**Action Steps:**
- Add the `## 6. Global Non-Functional Requirements (NFRs)` section with flexible `Attribute: Value` pairs (including Performance) to the bottom of the `template` in `src/internal/agent/artifacts/spec/requirements.yaml`.

**Verification (TDD):**
`grep "Performance" src/internal/agent/artifacts/spec/requirements.yaml` in the global template section should succeed.

- [x] T2.5: [CONFIG] Remove Legacy UI/UX Section
**Target:** `src/internal/agent/artifacts/spec/requirements.yaml`
**Context:** [REQ-2]

**Action Steps:**
- Delete the old `## 5. UI/UX Contract` section from the `template`.

**Verification (TDD):**
`grep "## 5." src/internal/agent/artifacts/spec/requirements.yaml | wc -l` should be 1.

### Phase 3: Final Verification

- [x] T3.1: [TEST] YAML Validity Check
**Target:** `src/internal/agent/artifacts/spec/requirements.yaml`
**Context:** [Global Scope]

**Action Steps:**
- Ensure the modified `src/internal/agent/artifacts/spec/requirements.yaml` is a valid YAML file.

**Verification (TDD):**
`python3 -c 'import yaml, sys; yaml.safe_load(sys.stdin)' < src/internal/agent/artifacts/spec/requirements.yaml`

- [x] T3.2: [TEST] Content Integrity Audit
**Target:** `Global Scope`
**Context:** [REQ-1, REQ-2]

**Action Steps:**
- Verify that all sub-sections (Safety, Validation, Performance, Observability, etc.) defined in the requirements spec are present in the template.

**Verification (TDD):**
Manual comparison between `requirements.md` and the final `requirements.yaml`.
