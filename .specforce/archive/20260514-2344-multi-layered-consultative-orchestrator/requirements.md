---
slug: 20260514-2344-multi-layered-consultative-orchestrator
lens: Balanced full-stack
---

# Feature: Multi-Layered Consultative Orchestrator

## 1. Context & Value
The current Specforce orchestrator (`spf.spec`) is focused on speed and intentionality but lacks a deep consultative phase. This feature introduces a three-layer "Pre-Flight" analysis (Constitution -> Codebase -> Consultative Grill) to ensure that every specification is grounded in the project's architecture and current state before artifacts are generated. This transforms the agent from a passive documentation writer into a high-fidelity technical design partner.

## 2. Out of Scope (Anti-Goals)
- Implementing a separate "grill" TUI component (the interview happens via standard agent interaction tools).
- Automating the fix of constitutional violations (the agent must only flag and discuss them).
- Modifying the core `specforce` binary's command structure (this is an instruction-level enhancement in `spec.yaml`).

## 3. Acceptance Criteria (BDD)

### [US-1] Constitutional Alignment Layer
**User Story:** AS A Developer, I WANT the orchestrator to check the project constitution first, SO THAT I don't propose features that violate our core principles or architecture.

**Scenarios:**
1. **[Happy Path]** GIVEN an active project with constitution docs WHEN I run `/spec` THEN the agent MUST read `.specforce/docs/` and summarize any relevant constraints to me before asking questions.
2. **[Conflict State]** GIVEN a feature request that violates a principle in `principles.md` WHEN the agent performs the check THEN it MUST cite the specific rule and ask how to reconcile the design.

**Technical Constraints (NFR):**
- **[Performance]:** Constitution check MUST reuse context if already loaded in the session.
- **[Safety & Security]:** Agent MUST NOT proceed if a critical security violation is identified without user confirmation.

### [US-2] Empirical Codebase Scan Layer
**User Story:** AS A Developer, I WANT the orchestrator to scan the existing code related to my request, SO THAT it can provide recommendations based on actual implementation patterns.

**Scenarios:**
1. **[Pattern Discovery]** GIVEN a request to modify a specific module WHEN the agent identifies the intent THEN it MUST use `grep_search` to find existing patterns in that module before the interview.
2. **[Zero-Context Fallback]** GIVEN a request for a completely new module WHEN no code patterns are found THEN the agent MUST inform me it found no existing precedents and will rely on Constitutional defaults.

**Technical Constraints (NFR):**
- **[Performance]:** Scans MUST be targeted and surgical (max 5-10 relevant files).
- **[Integrity]:** Agent MUST cite the files it read when making recommendations.

### [US-3] Consultative "Grill" Interview
**User Story:** AS A Developer, I WANT the orchestrator to "grill" me with technical recommendations and trade-offs, SO THAT we reach a shared understanding of the design tree.

**Scenarios:**
1. **[High-Fidelity Interview]** GIVEN the constitution and code context are loaded WHEN the agent starts the interview THEN it MUST provide a recommended answer for every question it asks.
2. **[Branching Decision]** GIVEN a technical choice with trade-offs WHEN the user makes a decision THEN the agent MUST explore the next logical branch of the design tree based on that choice.

**Technical Constraints (NFR):**
- **[Observability]:** The agent MUST use `ask_user` (or the discovered interaction tool) to maintain a structured dialogue.
- **[Performance]:** Decision loops MUST be resolved before generating `requirements.md`.

## 4. Business Invariants
- The agent MUST NOT generate any specification artifacts until the "Discovery & Intent Clarification" phase is marked as complete in its internal logic.
- The agent MUST cite at least one Constitutional document during the initialization of a new spec.

## 5. Global Non-Functional Requirements (NFRs)
- **[Performance]:** The multi-layered check should not add more than 2-3 turns to the initialization process.
- **[Reliability]:** The flow must be deterministic: Constitution -> Code -> Grill.
- **[Maintainability]:** Instructions must be modular enough to be updated as the kit's core logic evolves.
