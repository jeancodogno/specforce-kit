---
slug: 20260505-1759-zero-ambiguity-workflow
lens: Backend-heavy
---

# Implementation Roadmap: Zero-Ambiguity SDD Workflow

This roadmap implements the "Zero-Doubt" policy across the Specforce SDD workflow. It transitions the system from a permissive model that allowed documented ambiguities to a strictly deterministic model that halts for clarification before artifact generation.

## 1. Execution Strategy
- **Gravity Order:** Template Refactoring -> Orchestration & Skill Hardening -> Integration Verification. We must clean the "source of truth" (templates) before hardening the logic that uses them.

## 2. Tasks

### Phase 1: Artifact Template Refactoring
Purge existing YAML templates of ambiguity-related sections and inject mandatory non-assumption instructions to enforce determinism at the source.

- [x] T1.1: [CONFIG] Refactor `requirements.yaml` for strict determinism
**Target:** `src/internal/agent/artifacts/spec/requirements.yaml`
**Context:** [REQ-1]

**Action Steps:**
- Remove the `3. [Ambiguity]` scenario from the `template` section.
- Prepend the rule: "ASSUMPTIONS PROHIBITED: You MUST NOT make assumptions about business logic. If a rule is missing, you MUST halt and request clarification via the orchestrator." to the `instruction` field.

**Verification (TDD):**
`grep -i "ambiguity" src/internal/agent/artifacts/spec/requirements.yaml` should return 0 matches in the template section.

- [x] T1.2: [CONFIG] Inject Determinism Mandate into `design.yaml`
**Target:** `src/internal/agent/artifacts/spec/design.yaml`
**Context:** [REQ-1]

**Action Steps:**
- Add rule: "DETERMINISM MANDATE: Every technical decision and implementation step MUST be concrete. Use of 'TBD', 'To be defined', or placeholders is strictly forbidden." to the `instruction` field.

**Verification (TDD):**
`cat src/internal/agent/artifacts/spec/design.yaml | grep "DETERMINISM MANDATE"` should succeed.

- [x] T1.3: [CONFIG] Inject Determinism Mandate into `tasks.yaml`
**Target:** `src/internal/agent/artifacts/spec/tasks.yaml`
**Context:** [REQ-1]

**Action Steps:**
- Add rule: "DETERMINISM MANDATE: Every technical decision and implementation step MUST be concrete. Use of 'TBD', 'To be defined', or placeholders is strictly forbidden." to the `instruction` field.

**Verification (TDD):**
`cat src/internal/agent/artifacts/spec/tasks.yaml | grep "DETERMINISM MANDATE"` should succeed.

### Phase 2: Orchestration & Skill Hardening
Update the core command logic and specialized skills to transform ambiguity detection from a warning into a hard blocking gate.

- [x] T2.1: [CONFIG] Harden `spf.spec` Orchestrator Logic
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [REQ-2]

**Action Steps:**
- Change "If necessary, use your user-interaction tool to clarify" to "You MUST use your user-interaction tool to clarify ALL technical and business ambiguities before proceeding."
- Add a "Zero-Doubt Policy" guardrail: "Any artifact containing a placeholder or ambiguity marker is considered a failure and MUST be regenerated after clarification."

**Verification (TDD):**
`grep "Zero-Doubt Policy" src/internal/agent/kit/commands/spec.yaml` should return the new guardrail.

- [x] T2.2: [CONFIG] Update `pragmatic-product-owner` Skill
**Target:** `src/internal/agent/kit/skills/pragmatic-product-owner/SKILL.yaml`
**Context:** [REQ-2]

**Action Steps:**
- Add a "Zero-Ambiguity Mandate" section requiring binary, non-ambiguous ACs.

**Verification (TDD):**
`grep "Zero-Ambiguity Mandate" src/internal/agent/kit/skills/pragmatic-product-owner/SKILL.yaml` should succeed.

- [x] T2.3: [CONFIG] Redefine `spec-clarification-interview` as a Blocking Gate
**Target:** `src/internal/agent/kit/skills/spec-clarification-interview/SKILL.yaml`
**Context:** [REQ-2]

**Action Steps:**
- Update description to emphasize its role as a mandatory gate for resolving unknowns.

**Verification (TDD):**
`grep "mandatory gate" src/internal/agent/kit/skills/spec-clarification-interview/SKILL.yaml` should succeed.

### Phase 3: Integration Verification
Validate that the entire pipeline correctly halts and triggers resolution when presented with ambiguous input.

- [x] T3.1: [TEST] End-to-End "Halt & Resolve" Validation
**Target:** `Global Scope`
**Context:** [REQ-2]

**Action Steps:**
- Attempt to generate a spec using a deliberately vague prompt (e.g., "Add a feature that does something with users"). 
- Assert that the agent halts generation and invokes the `spec-clarification-interview` tool instead of creating a `requirements.md` with placeholders.

**Verification (TDD):**
Manual verification of the agent behavior in a new session or subagent delegation.
