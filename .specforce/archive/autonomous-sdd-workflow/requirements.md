---
slug: autonomous-sdd-workflow
lens: Integration
---

# Feature: Autonomous SDD Workflow via Skills

## 1. Context & Value
Currently, developers must manually type Specforce commands (e.g., `/spec`, `/implement`) to orchestrate the Spec-Driven Development (SDD) pipeline. This feature enables the AI agent to autonomously trigger these workflows when a specific intent (like a new feature request or a readiness to code) is detected. By conceptually treating these commands as "skills" instructed within `AGENTS.md`, we transform the agent from a passive tool into a proactive partner.

## 2. Out of Scope (Anti-Goals)
- Do not implement natural language parsing within the CLI; the agent's LLM handles intent detection.

## 3. Acceptance Criteria (BDD)

### [REQ-1] Agent Proactive Planning
**User Story:** AS A developer, I WANT the AI agent to automatically initiate the planning workflow when I request a new feature, SO THAT I don't have to manually type `/spec`.

**Scenarios:**
1. **[Happy Path]** GIVEN the agent is configured with the updated `AGENTS.md` WHEN the user says "Let's build a login page" THEN the agent automatically invokes the `/spec` command logic to start discovery and specification generation.
2. **[Edge Case]** GIVEN the user is already in a planning loop WHEN the user provides vague details THEN the agent uses the configured `ask_user` tool (as defined in the `spec` command) to clarify before proceeding, instead of hallucinating.

### [REQ-2] Agent Proactive Implementation
**User Story:** AS A developer, I WANT the AI agent to automatically start implementation once the spec is fully documented, SO THAT the workflow flows seamlessly from planning to execution.

**Scenarios:**
1. **[Happy Path]** GIVEN a feature specification is 100% complete and validated WHEN the user says "Looks good, let's build it" THEN the agent automatically invokes the `/implement` command logic to begin the TDD cycle.
2. **[Edge Case]** GIVEN the spec is only partially complete WHEN the user asks to implement THEN the agent refuses and prompts to finish the planning phase first.

### [REQ-3] Command Configuration Optimization & Dual Installation
**User Story:** AS A system architect, I WANT the workflows (commands) to be installed as both executable commands and agent skills during initialization, SO THAT the LLM can use them seamlessly in either context.

**Scenarios:**
1. **[Happy Path]** GIVEN a new project is being initialized with `specforce init` WHEN the kit is installed THEN the `commands/*.yaml` files are mapped and installed into both the agent's `commands/` directory and `skills/` directory.
2. **[Alternative Path]** GIVEN the agent evaluates its available tools WHEN it reads the `description` of the `spec` command (or skill) THEN it understands it is the "SDD Orchestrator" and should be activated for planning or altering specifications.

## 4. Business Invariants
- The agent MUST NEVER write implementation code without a fully documented and approved specification (`requirements.md`, `design.md`, `tasks.md`).
- The source of truth for the workflow instructions MUST remain the existing command `.yaml` files (e.g., `spec.yaml`), not duplicated elsewhere.
