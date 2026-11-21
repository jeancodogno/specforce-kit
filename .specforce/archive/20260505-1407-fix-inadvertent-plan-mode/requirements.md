---
slug: 20260505-1407-fix-inadvertent-plan-mode
lens: Integration
---

# Feature: Fix Inadvertent Plan Mode Activation

## 1. Context & Value
AI agents and subagents occasionally activate platform-native "plan modes" (like Gemini CLI's `enter_plan_mode`) during the Specforce planning phase. This causes redundant design artifacts and breaks the SDD workflow by initiating external processes that conflict with the sovereign Specforce loop.

## 2. Out of Scope (Anti-Goals)
- Modifying the core Gemini CLI binary or logic.
- Adding dependencies on specific CLI tools in the shared skills/agents.
- Changing how `specforce` binary itself operates.

## 3. Acceptance Criteria (BDD)

### [REQ-1] Inhibit External Planning in spf.spec
**User Story:** AS A Developer, I WANT the `spf.spec` skill to prevent the LLM from triggering external planning modes, SO THAT the Specforce workflow remains the sole source of truth.

**Scenarios:**
1. **[Happy Path]** GIVEN the `spf.spec` skill is active WHEN a complex feature is being planned THEN the LLM remains within the orchestration loop without calling `enter_plan_mode`.
2. **[Edge Case]** GIVEN a vague or massive feature request WHEN `spf.spec` is processing it THEN the LLM uses `ask_user` for clarification instead of attempting to enter a native plan mode for "research".

### [REQ-2] Inhibit External Planning in Subagents
**User Story:** AS an Orchestrator, I WANT subagents to act as stateless content generators, SO THAT they don't trigger their own planning workflows when delegated a task.

**Scenarios:**
1. **[Happy Path]** GIVEN the `Technical Solutions Architect` is invoked by `spf.spec` WHEN asked to design a schema THEN it produces the markdown content directly without initiating a native planning session.
2. **[Constraint]** GIVEN any subagent is invoked THEN it must adhere to a "Documentation-Only" and "No-External-Planning" mandate.

## 4. Business Invariants
- The Specforce SDD workflow is sovereign: no external planning artifacts should be created during a spec session.
- Instructions must remain platform-agnostic (no mention of specific tool names like `enter_plan_mode`).

## 5. UI/UX Contract
- N/A (Backend/Logic Heavy).
