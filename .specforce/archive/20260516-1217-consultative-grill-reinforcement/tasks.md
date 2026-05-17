---
slug: 20260516-1217-consultative-grill-reinforcement
lens: Integration
---

# Implementation Roadmap: Consultative Grill Reinforcement

## 1. Execution Strategy
- **Gravity Order:** Skill Refactoring -> Instruction Hardening -> Agent Integration -> Verification.

## 2. Tasks

### Phase 1: Skill Refactoring & Identity

- [x] T1.1: [REFACTOR] Rename Skill Directory
**Target:** `Global Scope`
**Context:** [US-1]

**Action Steps:**
- Move `src/internal/agent/kit/skills/spec-clarification-interview` to `src/internal/agent/kit/skills/consultative-grill`.

**Verification (TDD):**
`ls src/internal/agent/kit/skills/consultative-grill/SKILL.yaml` should succeed.

- [x] T1.2: [CONFIG] Update Skill Metadata
**Target:** `src/internal/agent/kit/skills/consultative-grill/SKILL.yaml`
**Context:** [US-1]

**Action Steps:**
- Change `name` to `consultative-grill`.
- Update `description` to: "Proactive consultative gate and adversarial design reviewer. Use to ground specifications in the project's Constitution and Codebase through a 'Grill' interview before artifact generation."

**Verification (TDD):**
`grep "name: consultative-grill" src/internal/agent/kit/skills/consultative-grill/SKILL.yaml` should succeed.

### Phase 2: Skill Hardening (Grill Logic)

- [x] T2.1: [CONFIG] Implement Consultative Triad & Adversarial Mode
**Target:** `src/internal/agent/kit/skills/consultative-grill/SKILL.yaml`
**Context:** [US-3, US-4]

**Action Steps:**
- Update `content` to include the **Consultative Triad** rule: `[📜 CONSTITUTION] + [🔍 CODEBASE] = [💡 RECOMMENDATION]`.
- Add the **Adversarial Grill** section: "If 'grill me' intent is detected, switch to Adversarial Mode to identify technical debt, security risks, and architectural gaps."
- Reinforce the "One question at a time" and "Binary decision" rules.

**Verification (TDD):**
`grep "CONSTITUTION" src/internal/agent/kit/skills/consultative-grill/SKILL.yaml` should succeed.

- [x] T2.2: [CONFIG] Modernize Grill Logic & Remove Legacy Markers
**Target:** `src/internal/agent/kit/skills/consultative-grill/SKILL.yaml`
**Context:** [US-3, US-4]

**Action Steps:**
- Remove the 3-question limit in Adversarial Mode.
- Replace legacy `[UNCLEAR]` and `{PLACEHOLDERS}` references with "Hardened Definitions" and "Decision Persistence" rules.
- Implement "Termination Signal" logic based on risk resolution.

**Verification (TDD):**
`grep "Hardened Definitions" src/internal/agent/kit/skills/consultative-grill/SKILL.yaml` should succeed.

### Phase 3: Orchestrator & Agent Integration

- [x] T3.1: [CONFIG] Harden spec.yaml Orchestrator
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [US-2, US-4]

**Action Steps:**
- Update Step 1.C (The Consultative Grill) to explicitly trigger the `consultative-grill` skill.
- Add a rule to the "Discovery & Intent Clarification" block: "If the user prompt contains 'grill me', you MUST activate the consultative-grill skill in Adversarial Mode immediately."
- Clarify that artifact generation is blocked until core design decisions are confirmed via the Grill.

**Verification (TDD):**
`grep "consultative-grill" src/internal/agent/kit/commands/spec.yaml` should succeed.

- [x] T3.2: [CONFIG] Update Agent Skill Sets
**Target:** `src/internal/agent/kit/agents/`
**Context:** [US-1]

**Action Steps:**
- In `product-analyst.yaml`: Replace `spec-clarification-interview` with `consultative-grill`.
- In `technical-solution-architect.yaml`: Add `consultative-grill` to the `skills` list.
- In `technical-developer.yaml`: Add `consultative-grill` to the `skills` list.

**Verification (TDD):**
`grep "consultative-grill" src/internal/agent/kit/agents/technical-solution-architect.yaml` should succeed.

- [x] T3.3: [CONFIG] Modernize Agent Ambiguity Handling
**Target:** `src/internal/agent/kit/agents/`
**Context:** [US-1, US-3]

**Action Steps:**
- Update `product-analyst.yaml` to replace `[UNCLEAR]` tags with a mandate to trigger the `consultative-grill`.
- Update `technical-solution-architect.yaml` to replace manual flagging with `consultative-grill` triggers.

**Verification (TDD):**
`grep "trigger the consultative-grill" src/internal/agent/kit/agents/product-analyst.yaml` should succeed.

### Phase 4: Final Verification

- [x] T4.1: [VERIFY] Orchestrator behavioral dry-run
**Target:** `Global Scope`
**Context:** [US-2, US-4]

**Action Steps:**
- Execute `specforce spec status 20260516-1217-consultative-grill-reinforcement --json` one final time to ensure the spec itself is valid and all tasks are tracked.

**Verification (TDD):**
The command returns `"progress": 100` and `"is_valid": true`.
