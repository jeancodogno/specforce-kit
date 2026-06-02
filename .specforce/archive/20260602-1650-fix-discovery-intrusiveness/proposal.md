# Proposal: Fix Discovery Command Intrusiveness

## Problem Statement
The current `spf.discovery` command (defined in `src/internal/agent/kit/commands/discovery.yaml`) is too proactive in pushing users toward formal specification initialization. It frequently interrupts the exploratory phase with `ask_user` prompts and mandatory "Formalization Handoff" steps, contradicting its goal of being a "purely read-only mode for exploration and research."

## Research Findings
- **Intrusive Instructions:** The `Scope & Complexity Assessment` section mandates the use of `ask_user` to confirm decomposition strategies.
- **Forced Handover:** The `Formalization Handoff` section mandates an `ask_user` call to confirm `proposal.md` creation.
- **Workflow Bias:** The prompt structure assumes a linear path from research to planning, whereas discovery is often non-linear or purely informative.

## Proposed Strategy
1.  **Refactor Mindset:** Update the "TASK" and "Scout's Principles" to emphasize open-ended exploration and a "passive intelligence" model.
2.  **Advisory Scope Assessment:** Change the "Scope & Complexity Assessment" to be informative only. If a topic is large, the agent should note it in the brief but NOT trigger an `ask_user` to split it.
3.  **On-Demand Formalization:** Update the "Formalization Handoff" to provide the "Scout Intelligence Brief" and *inform* the user that they can request a proposal or run `/spf:spec`, instead of proactively asking via `ask_user`.
4.  **Preserve Sandbox Integrity:** Ensure the agent stays in "Discovery Mode" until the user explicitly directs a pivot to planning or execution.

## Expected Outcome
A more relaxed and truly exploratory Discovery mode that provides high-value architectural and technical insights without forcing the user into the SDD state machine prematurely.
