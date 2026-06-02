---
slug: 20260602-1650-fix-discovery-intrusiveness
lens: Bugfix
---

# Bugfix: Discovery Mode Intrusiveness & Workflow Friction

## 1. Issue Description
The `spf.discovery` skill, intended as a "purely read-only sanctuary" for exploration and research, has become overly proactive and intrusive. It currently mandates the use of `ask_user` for scope assessments and formalization handoffs. This interrupts the natural flow of technical brainstorming, forces unnecessary wait states, and contradicts the skill's goal of being a low-pressure intellectual sandbox.

## 2. Evidence & Observations
- **Symptom:** The agent stops exploration to ask for permission to recommend a split or to create a proposal, even when the user is still in the middle of a discovery thread.
- **Observations:** In `src/internal/agent/kit/commands/discovery.yaml`, the sections "Scope & Complexity Assessment" and "Formalization Handoff" explicitly require `ask_user` tool calls.
- **Trace:** During discovery sessions, agents often trigger `ask_user` with "Would you like to split this spec?" or "Should I create a proposal.md?", creating a bottleneck in what should be a high-bandwidth conversational phase.

## 3. Reproduction Steps
1. Activate `spf.discovery`.
2. Provide a complex feature idea that spans multiple domains.
3. Observe the agent's behavior: it will hit the "Scope & Complexity Assessment" logic and call `ask_user` to ask if it should split the spec, halting the conversation.
4. Try to reach a point of technical clarity.
5. Observe the "Formalization Handoff": the agent will again call `ask_user` to ask if it should create a `proposal.md`.

## 4. Root Cause Analysis (RCA)
- **Hard-coded Interaction Gates:** The `discovery.yaml` instruction set mandates `ask_user` as a "MUST" for assessment and handoff.
- **Active Mindset Bias:** The current instructions favor an "Active/Assertive" mindset over a "Passive/Advisory" one, which is inappropriate for the discovery phase.
- **Workflow Coupling:** Formalization (creation of `proposal.md`) is gated by mandatory confirmation rather than being offered as an on-demand technical capability.

## 5. Out of Scope (Anti-Goals)
- Do not remove the `proposal.md` creation capability itself (only the mandatory `ask_user` tool gate).
- Do not modify other agent skills (e.g., `spf.spec`, `spf.implement`, `spf.archive`).
- Do not change the `specforce spec init` command logic or its underlying Go implementation.

## 6. Acceptance Criteria (Regression Tests)

### [FIX-1] Passive Intelligence Mindset
**Scenario: [Regression]**
GIVEN the agent is in `spf.discovery` mode
WHEN providing technical insights or architectural sketches
THEN it MUST NOT use the `ask_user` tool to validate its thoughts or guide the conversation.
THEN it MUST use advisory language and open-ended questions in its text output to guide the developer.

### [FIX-2] Advisory Scope Assessment
**Scenario: [Regression]**
GIVEN a brainstormed idea is identified as "Too Large" (spanning multiple domains or >20 requirements)
WHEN the agent performs a scope assessment
THEN it MUST provide an "Informative Recommendation" in the text output explaining the decomposition rationale.
THEN it MUST NOT call the `ask_user` tool to gate the execution or force a decision on splitting.

### [FIX-3] On-Demand Formalization (The Proposal Protocol)
**Scenario: [Regression]**
GIVEN technical clarity is achieved in discovery
WHEN the agent suggests formalization
THEN it MUST present the "Scout Intelligence Brief" and "Decomposition Strategy" (if applicable) as a textual summary.
THEN it MUST inform the user that it *can* execute `specforce spec init` and create `proposal.md` if the user explicitly requests it.
THEN it MUST NOT proactively call the `ask_user` tool to ask "Would you like to formalize?".

### [FIX-4] Elimination of Mandatory ask_user in discovery.yaml
**Scenario: [Regression]**
GIVEN the `src/internal/agent/kit/commands/discovery.yaml` file
WHEN inspected
THEN all mandatory requirements for `ask_user` tool calls in the "Scope & Complexity Assessment" and "Formalization Handoff" sections MUST be removed or converted to advisory text instructions.

## 7. Technical Constraints (NFR)
- **[Mindset]:** The Scout must remain "Inquisitive & Emergent" without being "Obstructionist".
- **[Consistency]:** The "Non-Mutation Covenant" remains absolute: no changes to `src/` or `tests/`.
- **[Usability]:** Suggestions for `proposal.md` must be clear enough that a user can easily say "Go ahead and create the proposal" to trigger the next step.
