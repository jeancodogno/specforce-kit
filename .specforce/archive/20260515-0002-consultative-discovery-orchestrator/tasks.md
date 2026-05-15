---
slug: 20260515-0002-consultative-discovery-orchestrator
lens: Integration
---

# Implementation Roadmap: Consultative Discovery Orchestrator

## 1. Execution Strategy
- **Gravity Order:** Update the `discovery.yaml` command instructions -> Verify behavior via a dry-run.

## 2. Tasks

### Phase 1: Instruction Refactoring

- [x] T1.1: [DOCS] Refactor `spf.discovery` Instructions
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [US-1, US-2, US-3]

**Action Steps:**
- Rewrite the `Ecosystem Contextualization` and `Diagnostic Workflow` sections into a cohesive `Consultative Exploration Funnel`.
- Add explicit steps for `Layer 1: Constitutional Anchor` and `Layer 2: Codebase Archeology`.
- Add explicit instructions for `Layer 3: Consultative Brainstorming` detailing the requirement to present 1-3 paths with trade-offs.
- Ensure the Non-Mutation Covenant remains prominently displayed.

**Verification (TDD):**
Run `/spf:discovery` on a dummy vague request (e.g., "I want to add a caching layer") and verify the agent:
1. Reads the constitution.
2. Checks the codebase for existing caching.
3. Proposes multiple caching options with trade-offs instead of just demanding a decision.