---
slug: 20260528-0055-reinforce-agent-orchestration
lens: Backend-heavy
---

# 1. Executive Summary

## Problem Statement
Currently, the implementation orchestration (`implement.yaml`) is less structured than the planning orchestration (`spec.yaml`). It lacks a formal "Mission Brief" envelope and high-priority project rules, leading to potential context leaks and inconsistent behavior by implementation subagents. Additionally, there is no explicit gate to ensure coherence between planning artifacts during the specification phase.

## Objective
Standardize the agent orchestration and delegation patterns across the Specforce SDD pipeline. This involves reinforcing both `spec.yaml` and `implement.yaml` with structured mission briefs, priority rules, and a coherence gate to ensure architectural integrity and specialized agent execution.

## Success Metrics
- **Business Metric:** 100% of implementation tasks executed through `spf.implement` utilize a structured MISSION BRIEF with PROJECT RULES.
- **Performance Target:** Orchestration logic and prompt synthesis overhead remains < 100ms per agent delegation.
- **UX Efficiency:** Zero manual intervention required by Lead Architects to enforce global project rules during implementation sessions.

# 2. Key Entities & Domain Integrity

## Agent Orchestration Pattern
The standardized protocol for a "Primary" orchestrator to discover, brief, and delegate work to "Secondary" subagents.

## Mission Brief Envelope
A high-fidelity prompt container used to communicate a specific task to a subagent. It MUST include:
1. **Mandated Skills:** Specific skills/heuristics to activate.
2. **Project Rules (Maximum Priority):** Non-negotiable global constraints.
3. **Shared Context:** Feature-specific design decisions and upstream dependencies.
4. **Base Instructions:** The specific atomic task to perform.

## Coherence Gate
A logical checkpoint within the planning workflow that verifies the consistency and alignment between `requirements.md`, `design.md`, and `tasks.md` before a specification is considered "Ready for Implementation".

# 3. Golden Rules (Business Invariants)
- **Rule 1: Priority Dominance:** "PROJECT RULES (MAXIMUM PRIORITY)" MUST always be positioned above "BASE INSTRUCTIONS" in any delegated prompt.
- **Rule 2: Specialization over Genericity:** Implementation tasks MUST be delegated to specialized roles (`specforce-developer`) and verification to QA roles (`specforce-qa`).
- **Rule 3: Artifact Synchronicity:** No specification session can terminate successfully if there is a detected drift between requirements and downstream artifacts (Design/Tasks).

# 4. Functional Requirements

## [US-1] Standardized Implementation Mission Brief
As a Lead Architect, I want the implementation engine to use a formal Mission Brief envelope so that implementation subagents operate under strict project-wide constraints.

### Acceptance Criteria
- **GIVEN** a pending implementation task
- **WHEN** the orchestrator delegates the task to a subagent
- **THEN** it MUST construct a MISSION BRIEF envelope containing sections for Mandated Skills, Project Rules (Maximum Priority), Shared Context, and Base Instructions.
- **GIVEN** a Mission Brief is being constructed
- **WHEN** project-specific instructions exist in the project configuration or status
- **THEN** they MUST be populated into the `## 2. PROJECT RULES (MAXIMUM PRIORITY)` section.

**Technical Constraints (NFR):**
- **[Performance]:** Envelope generation MUST NOT increase total prompt latency by more than 50ms.
- **[Reliability]:** 100% task coverage for Mission Brief envelope.

## [US-2] Spec Coherence Gate & Verification
As a Lead Architect, I want the planning engine to verify artifact coherence so that implementation starts from a consistent source of truth.

### Acceptance Criteria
- **GIVEN** a planning session is nearing completion
- **WHEN** the orchestrator reaches the "Verification & Handoff" phase
- **THEN** it MUST execute a "Coherence Gate" check that compares requirements, design, and tasks for alignment.
- **GIVEN** a drift is detected between `requirements.md` and `tasks.md` (e.g., a requirement has no corresponding task)
- **WHEN** the Coherence Gate is active
- **THEN** the orchestrator MUST block the session conclusion and prompt for alignment.

**Technical Constraints (NFR):**
- **[Performance]:** Coherence check execution SHOULD NOT exceed 1s of processing time.

## [US-3] Specialized Agent Delegation
As a Lead Architect, I want implementation and QA to be handled by specialized agent personas so that outcomes are professionally verified.

### Acceptance Criteria
- **GIVEN** a task requires code modification
- **WHEN** the implement orchestrator performs agent discovery
- **THEN** it MUST prioritize the `specforce-developer` persona with `spf.implement` skills.
- **GIVEN** all implementation tasks are finished
- **WHEN** the "Final QA Validation" phase starts
- **THEN** it MUST explicitly switch to the `specforce-qa` persona and activate BDD-focused verification skills.

**Technical Constraints (NFR):**
- **[Architecture]:** Delegation MUST be metadata-driven (matching agent roles to task types).
- **[Performance]:** Persona switching latency < 100ms.

# 5. Global Non-Functional Requirements (NFRs)
- **[Traceability]:** All delegated Mission Briefs MUST be visible in the agent's output log for auditability.
- **[Self-Containment]:** The orchestration logic MUST NOT depend on the presence of specific external subagents, failing gracefully if a specialized role is unavailable.
- **[Idempotency]:** Re-running an orchestration loop with the same state MUST produce a consistent Mission Brief envelope.

# 6. Out of Scope
- Modification of individual agent skill logic (e.g., changing how `tdd` works).
- Addition of new sub-commands to the `specforce` CLI binary.
- UI/UX changes to the terminal-based progress bars or TUI views.
- Automating the "Approval" of specs (remains a Human-Governed action).
