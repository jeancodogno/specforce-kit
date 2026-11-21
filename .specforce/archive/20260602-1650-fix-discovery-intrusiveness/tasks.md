---
slug: 20260602-1650-fix-discovery-intrusiveness
lens: Backend-heavy
---

# Implementation Roadmap: Fix Discovery Command Intrusiveness

## 1. Execution Strategy
The implementation follows a surgical refactoring of the `spf.discovery` skill instructions to shift from an intrusive, tool-gated behavior to a passive, advisory intelligence model. The work is divided into three phases: auditing and preparing the target file, applying behavioral changes to the core instructions, and performing manual regression tests to ensure the new "Passive Intelligence" mindset is respected.

## 2. Tasks

### Phase 1: Preparation & Scaffolding

- [x] T1.1: [AUDIT] Identify exact line ranges in `discovery.yaml` for behavioral modification
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** Map current `ask_user` triggers for Scope Assessment and Formalization.

**Action Steps:**
- Identify blocks in `discovery.yaml` containing `ask_user` and "Wait for Confirmation".
- Record the exact line numbers for surgical replacement.

**Acceptance Check:**
List of line ranges for blocks needing refactoring is identified.

### Phase 2: Behavioral Refactoring

- [x] T2.1: [REFRACTOR] Update the "Non-Mutation Covenant" to reflect passive triggering
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [FIX-1]

**Action Steps:**
- Locate the "Non-Mutation Covenant" section in `discovery.yaml`.
- Replace "after obtaining explicit user confirmation" with "upon explicit user request".

**Acceptance Check:**
Phrase "upon explicit user request" is present in the Covenant section.

- [x] T2.2: [REFRACTOR] Convert Scope Assessment from a tool-gate to an advisory recommendation
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [FIX-2]

**Action Steps:**
- Remove "Wait for Confirmation" and "MUST use the ask_user tool" instructions.
- Replace with instructions to provide an "Informative Recommendation" explaining the rationale for decomposition in the text output.

**Acceptance Check:**
`ask_user` is no longer mentioned in the Scope & Complexity Assessment section.

- [x] T2.3: [REFRACTOR] Refactor Formalization Handoff to be user-led and on-demand
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [FIX-3]

**Action Steps:**
- Remove the requirement to use `ask_user`.
- Instruct the agent to state it *can* execute `specforce spec init` if requested by the user.

**Acceptance Check:**
Mandatory interaction tools for formalization are removed.

### Phase 3: Final Verification

- [x] T3.1: [VERIFY] Regression test: Large Scope Advisory (Manual)
**Target:** Agent Session (`spf.discovery`)
**Context:** [FIX-2, FIX-4]

**Action Steps:**
- Trigger a "Too Large" feature discussion in a simulated `spf.discovery` session.
- Verify the agent recommends a split in text but does **NOT** call `ask_user`.

**Acceptance Check:**
Agent provides decomposition rationale without an interrupting tool call.

- [x] T3.2: [VERIFY] Regression test: On-Demand Formalization (Manual)
**Target:** Agent Session (`spf.discovery`)
**Context:** [FIX-3, FIX-4]

**Action Steps:**
- Reach technical clarity on a topic in a simulated `spf.discovery` session.
- Verify the agent presents the "Scout Intelligence Brief" and mentions its capability to create `proposal.md` without calling `ask_user`.

**Acceptance Check:**
Brief is presented; no `ask_user` tool is triggered proactively.

- [x] T3.3: [VERIFY] Regression test: User-Led Execution (Manual)
**Target:** Agent Session (`spf.discovery`)
**Context:** [FIX-3]

**Action Steps:**
- Reach a point of technical clarity in a simulated `spf.discovery` session.
- Explicitly ask the agent: "Create a proposal for this" and verify it executes the correct tools.

**Acceptance Check:**
Agent successfully executes `specforce spec init` and `write_file` for the `proposal.md`.
