---
slug: 20260516-1217-consultative-grill-reinforcement
lens: Integration
---

# Feature: Consultative Grill Reinforcement

## 1. Context & Value
Reinforce the SDD workflow by transforming the passive "clarification interview" into a proactive "Consultative Grill". This ensures all specifications are grounded in the project's Constitution and Codebase before artifact generation, reducing hallucinations and architectural drift.

## 2. Out of Scope (Anti-Goals)
- Do not implement a separate TUI component for the grill; it remains an LLM-driven conversational pattern using existing tools (`ask_user`).
- Do not modify the `specforce` binary code for command execution logic.

## 3. Acceptance Criteria (BDD)

### [US-1] Skill Renaming & Identity Update
**User Story:** AS A Developer, I WANT the clarification skill to be identified as "Consultative Grill", SO THAT its proactive and design-focused nature is clear.

**Scenarios:**
1. **[Happy Path]** GIVEN the skill directory `src/internal/agent/kit/skills/spec-clarification-interview` WHEN the system is updated THEN it MUST be moved to `src/internal/agent/kit/skills/consultative-grill` AND the `name` field in `SKILL.yaml` MUST be `consultative-grill`.
2. **[Edge Case]** GIVEN existing agents referencing the old skill name WHEN they are updated THEN they MUST correctly reference `consultative-grill`.

**Technical Constraints (NFR):**
- **[Performance]:** Skill discovery latency must remain < 50ms.
- **[Safety & Security]:** Ensure no broken references in agent configuration files.
- **[Integrity]:** Directory name and YAML `name` field must be identical.

### [US-2] Mandatory Consultative Pre-Flight
**User Story:** AS A Product Owner, I WANT the orchestrator to mandate the Grill phase, SO THAT no specification is generated without architectural alignment.

**Scenarios:**
1. **[Happy Path]** GIVEN the `spec.yaml` orchestrator WHEN initializing a new spec THEN it MUST explicitly state that the "Consultative Grill" (Layer 3) is a mandatory gate that requires user confirmation of recommendations.
2. **[Edge Case]** GIVEN a user who tries to skip the discovery phase WHEN the orchestrator is running THEN it MUST halt and insist on the Consultative Grill if core branches (Security, Architecture) are not resolved.

**Technical Constraints (NFR):**
- **[Performance]:** Orchestration loop overhead must be negligible.
- **[Integrity]:** The workflow MUST remain sequential: Constitution -> Code -> Grill.

### [US-3] Consultative Decision Pattern
**User Story:** AS A Developer, I WANT the agent to provide recommendations instead of just asking questions, SO THAT I can benefit from its analysis of the project's state.

**Scenarios:**
1. **[Happy Path]** GIVEN the `consultative-grill` skill is active WHEN proposing a design decision THEN it MUST present the triad: **[Constitution] + [Codebase] = [Recommendation]**.
2. **[Edge Case]** GIVEN a conflict between the Constitution and the user's request WHEN the grill is active THEN the agent MUST flag the violation and recommend a Constitution-compliant path.

**Technical Constraints (NFR):**
- **[Performance]:** Decision synthesis must occur within the LLM's standard response time.
- **[Observability]:** Decisions must be clearly logged in the conversation history.

### [US-4] Intent-Triggered "Grill Me" Mode
**User Story:** AS A Lead Architect, I WANT to challenge my own ideas by asking the agent to "Grill me", SO THAT we find potential failures before implementation.

**Scenarios:**
1. **[Happy Path]** GIVEN the user prompt contains "grill me" or similar intent WHEN the orchestrator triggers the grill THEN the agent MUST adopt an adversarial persona focused on finding edge cases, security gaps, and technical debt.
2. **[Edge Case]** GIVEN the "Grill Me" mode is active WHEN the user provides a weak justification THEN the agent MUST push back and ask for a more robust technical rationale.

**Technical Constraints (NFR):**
- **[Performance]:** Transition to adversarial mode must be instantaneous.
- **[Safety & Security]:** The adversarial mode must not violate any safety guardrails.

## 4. Business Invariants
- Artifact generation MUST NOT start until the "Consultative Grill" phase is marked as complete in the orchestrator's logic.
- Recommendations MUST always prioritize the project's Constitution over user "wants" unless explicitly overridden with a technical justification.

## 5. Global Non-Functional Requirements (NFRs)
- **[Reliability]:** Renaming must be atomic across all configuration files in `src/`.
- **[Security]:** The Grill must always include a security review layer (referencing `security.md`).
- **[Maintainability]:** Instruction updates in `spec.yaml` must follow the established multi-layered structure.
