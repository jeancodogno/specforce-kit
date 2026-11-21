---
slug: 20260527-0014-standardize-agents-naming-and-prompts
lens: Integration
---

# Feature: Standardize Agents Naming and Prompts

## 1. Context & Value
The current agent ecosystem has inconsistent naming conventions and orchestration bugs (e.g., `spf.spec` delegating to the developer instead of the planner). Standardizing agents with the `specforce-` prefix and improving prompts will ensure a cohesive, reliable, and portable multi-agent experience across different LLM environments.

## 2. Out of Scope (Anti-Goals)
- Implementing new agents beyond the existing ones.
- Modifying the core Go logic of the CLI (except for command YAML definitions).
- Changing the Skill/Heuristic system logic.

## 3. Acceptance Criteria (BDD)

### [US-1] Agent Naming Standardization
**User Story:** AS A Specforce maintainer, I WANT TO use a consistent naming convention for all agents, SO THAT they are easily distinguishable from orchestration commands.

**Scenarios:**
1. **[Happy Path]** GIVEN the agents directory WHEN I rename agents to use the `specforce-` prefix THEN all internal name fields and filenames MUST match.
2. **[Ref Integrity]** GIVEN renamed agents WHEN I check command YAMLs THEN all `delegate` or `activate_agent` references MUST be updated to the new names.

**Technical Constraints (NFR):**
- **[Integrity]:** Filename and internal `name` field MUST be identical (prefix included).
- **[Performance]:** Zero impact on agent loading time.

### [US-2] Correct Orchestration Delegation
**User Story:** AS A developer using `/spec`, I WANT the orchestrator to delegate tasks to the appropriate specialist, SO THAT the implementation roadmap is high-quality.

**Scenarios:**
1. **[Planner Fix]** GIVEN the `spf.spec` command WHEN generating the `tasks` artifact THEN it MUST delegate to `specforce-planner` instead of `specforce-developer`.
2. **[QA Handoff]** GIVEN the `specforce-qa` agent WHEN completing validation THEN the handoff message MUST point to `/spf:archive` or `/spf:implement` instead of non-existent agents.

**Technical Constraints (NFR):**
- **[Reliability]:** Delegation MUST fail-fast if the target agent is not found.
- **[Performance]:** Latency for delegation < 100ms.

### [US-3] Environment Portability & Sub-Agent Awareness
**User Story:** AS A user on different AI platforms, I WANT agents to work consistently across providers, SO THAT my workflow is not tied to a single model.

**Scenarios:**
1. **[Sub-Agent Awareness]** GIVEN an agent prompt WHEN defined THEN it MUST include instructions on how to handle environments that DO NOT support sub-agent spawning (no recursive calls).

**Technical Constraints (NFR):**
- **[Safety & Security]:** Agents MUST NOT attempt to spawn sub-agents if they detect they are already in a restricted sub-agent context.

### [US-4] Strict Prompt Hardening
**User Story:** AS A Specforce user, I WANT agents to stay within their functional boundaries, SO THAT they don't produce hallucinated code during planning phases.

**Scenarios:**
1. **[No Code in Planning]** GIVEN the `specforce-product-analyst` or `specforce-planner` agents WHEN generating artifacts THEN they MUST NOT output implementation code (Go, TS, etc.).
2. **[Interactive Consultation]** GIVEN an ambiguity WHEN an agent is processing THEN it MUST use the environment's interaction tool (`ask_user`, `ask`) before proceeding.

**Technical Constraints (NFR):**
- **[Reliability]:** Mandatory use of the "Mission Brief Envelope" in command prompts.
- **[Performance]:** Prompt tokens optimized for 60-90% savings where possible.

## 4. Business Invariants
- No agent file can exist without a `specforce-` prefix.
- The `technical-developer` (now `specforce-developer`) must never be the target for `tasks.md` generation.

## 5. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Instantaneous agent discovery and mapping.
- **[Reliability]:** Zero-Doubt Policy: Agents must halt and ask if context is insufficient.
- **[Maintainability]:** YAML files must follow standard formatting and avoid redundant logic.
- **[Portability]:** Agents must work in CLI, IDE extensions, and web interfaces where Specforce is supported.
