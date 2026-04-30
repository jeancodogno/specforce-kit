---
slug: discovery-command
lens: Integration
---

# Feature: Discovery Command

## 1. Context & Value
Developers using Specforce often need an unstructured phase to brainstorm features or investigate bugs before committing to a formal SDD specification. The `discovery` command provides a collaborative AI-driven environment for exploration, technical research, and diagnostic investigation while maintaining system safety through a strict read-only policy.

## 2. Out of Scope (Anti-Goals)
- Do not implement any code changes or file modifications during the discovery phase.
- Do not generate Specforce artifacts (requirements, design, tasks) automatically.
- Do not integrate with external logging or monitoring platforms in this phase.
- Do not support persistent sessions; each discovery session is ephemeral.

## 3. Acceptance Criteria (BDD)

### [REQ-1] Mode Selection: Brainstormer vs. Detective
**User Story:** AS A Developer, I WANT TO choose between brainstorming new ideas or investigating existing bugs, SO THAT the AI agent adopts the appropriate persona and toolset for my specific need.

**Scenarios:**
1. **[Happy Path - Brainstorming]** GIVEN the user starts a discovery session with a feature idea WHEN the agent identifies the intent THEN it adopts the "Senior Product Architect" persona and focuses on clarifying scope, feasibility, and tradeoffs.
2. **[Happy Path - Investigation]** GIVEN the user reports a bug or technical issue WHEN the agent identifies the intent THEN it adopts the "Senior Systems Engineer" persona and focuses on tracing code paths and formulating hypotheses.
3. **[Edge Case - Vague Intent]** GIVEN the user provides a vague or ambiguous prompt WHEN the discovery session begins THEN the agent MUST ask clarifying questions to determine if it should act as a Brainstormer or a Detective.

### [REQ-2] Strict Read-Only Enforcement
**User Story:** AS A Tech Lead, I WANT TO ensure the AI agent cannot modify the codebase during discovery, SO THAT the system integrity is preserved during unstructured exploration.

**Scenarios:**
1. **[Happy Path]** GIVEN a discovery session is active WHEN the agent proposes a solution THEN it MUST NOT attempt to apply the fix or write any implementation code.
2. **[Negative Case]** GIVEN a discovery session is active WHEN the user explicitly asks the agent to "fix this bug" or "create this file" THEN the agent MUST refuse and explain that it is in a read-only discovery mode.

### [REQ-3] SDD Handoff Recommendation
**User Story:** AS A Developer, I WANT TO be guided towards the formal SDD pipeline once my discovery is complete, SO THAT my findings can be converted into actionable and verified tasks.

**Scenarios:**
1. **[Happy Path]** GIVEN a discovery session has reached a conclusion or identified a clear path WHEN the conversation naturally ends THEN the agent MUST recommend running the `/spec` command to start the formal planning phase.

## 4. Business Invariants
- The discovery phase is purely conversational; no state is saved in the Specforce repository artifacts.
- The agent MUST prioritize reading Constitution artifacts (`.specforce/docs/`) during brainstorming to ensure alignment with project principles.
- The command MUST be available across all supported agent platforms (Gemini, Claude, OpenCode, etc.).

## 5. UI/UX Contract
*(Omitted - Lens is Integration)*
