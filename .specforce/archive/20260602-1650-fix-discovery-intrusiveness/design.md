---
slug: 20260602-1650-fix-discovery-intrusiveness
lens: Bugfix
---

# Bug Design: Discovery Mode Intrusiveness & Workflow Friction

## 1. Technical Strategy
The core of this fix is a shift from an **Active/Gatekeeper** mindset to a **Passive/Advisory** intelligence posture within the `spf.discovery` skill. The discovery phase is intended to be a "low-pressure intellectual sandbox," but current instructions mandate tool-based interruptions (`ask_user`) that break the developer's focus.

### 1.1. Passive Intelligence Transition
- **From:** Gating execution with mandatory `ask_user` tool calls.
- **To:** Providing "Informative Recommendations" via text output, allowing the user to drive the transition when ready.

### 1.2. Protocol Refinement
- **Advisory Scope Assessment:** Instead of halting to ask "Should I split this?", the agent will now describe the complexity and *recommend* a split in its text response, continuing the conversation unless the user explicitly asks to pivot.
- **On-Demand Formalization:** The agent will present the "Scout Intelligence Brief" and "Decomposition Strategy" as a conclusion to a discovery thread and inform the user that it is ready to execute `specforce spec init` upon request, rather than proactively triggering a confirmation prompt.

## 2. Affected Components

| Component | Path | Impact |
|--- |--- |--- |
| **Discovery Skill** | `src/internal/agent/kit/commands/discovery.yaml` | Primary prompt instruction set for the Scout agent. |

## 3. Implementation Plan

### 3.1. Instruction Refactoring (`discovery.yaml`)
1.  **Covenant Update:** Clarify that the "Proposal Exception" is triggered by explicit user request in the chat, not by a proactive tool-gated prompt.
2.  **Scope Assessment Refactor:**
    - Remove the `Wait for Confirmation` requirement.
    - Remove the mandatory `ask_user` instruction.
    - Replace with instructions to provide an "Informative Recommendation" explaining the rationale for decomposition.
3.  **Formalization Handoff Refactor:**
    - Remove the instruction to use `ask_user` for formalization.
    - Instruct the agent to present the brief and state that it *can* formalize if requested.
    - Standardize the "Scout Intelligence Brief" as the final discovery artifact before user-led transition.

## 4. Verification Strategy

### 4.1. Regression Testing (Manual)
Since this is a prompt-engineering fix, verification relies on observing agent behavior under specific triggers:

| Test Case | Trigger | Expected Behavior (PASS) |
|--- |--- |--- |
| **Large Scope Detection** | Describe a feature that spans 3+ domains and >20 requirements. | Agent suggests a split in text but does **NOT** call `ask_user`. |
| **Technical Clarity Reached** | Ask the agent to summarize a complex architecture. | Agent provides the Intelligence Brief and mentions it can create `proposal.md` without calling `ask_user`. |
| **User-Led Formalization** | User says: "Create the proposals for these 3 sub-features." | Agent executes `specforce spec init` and `write_file` for `proposal.md` as expected. |

### 4.2. Prompt Integrity Check
- Validate that the "Non-Mutation Covenant" (read-only) remains intact.
- Ensure no implementation code or binary logic is modified, preserving the structural integrity of the `specforce` Go codebase.
