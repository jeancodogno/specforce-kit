---
slug: 20260514-2344-multi-layered-consultative-orchestrator
lens: Integration
---

# Implementation Roadmap: Multi-Layered Consultative Orchestrator

## 1. Execution Strategy
The implementation focuses on refactoring the `spec.yaml` instruction file. Since this is an instruction change, "verification" involves a dry-run of the orchestrator's behavior to ensure it follows the new multi-layered protocol.

## 2. Tasks

### Phase 1: Instruction Refactoring

- [x] T1.1: [DOCS] Refactor `spf.spec` Discovery Section
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [US-1, US-2, US-3]

**Action Steps:**
- Rewrite the "Discovery & Intent Clarification" block.
- Explicitly define the "Constitutional Anchor" step.
- Explicitly define the "Empirical Grounding" step.
- Explicitly define the "Consultative Grill" step with the recommendation mandate.

**Verification (TDD):**
Run `/spf:spec` on a dummy request and verify that the agent:
1. Mentions reading `.specforce/docs/`.
2. Performs a `grep` or search for existing code.
3. Provides a recommendation before asking a design question.

## 3. Pre-emptive Mitigations
- **Risk:** The agent might become too verbose or "loop" in the interview.
- **Mitigation:** Include instructions to "consolidate design decisions" and "move to artifact generation once core branches are resolved".
