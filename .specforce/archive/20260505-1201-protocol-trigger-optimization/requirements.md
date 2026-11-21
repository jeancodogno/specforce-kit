---
slug: 20260505-1201-protocol-trigger-optimization
lens: Balanced full-stack
---

# Feature: Protocol Trigger Optimization

## 1. Context & Value
Improve the reliability of the Spec-Driven Development (SDD) protocol by refining how and when agents trigger Specforce workflows (Discovery, Spec, Implement). Currently, agents may bypass these workflows and edit code directly, violating the "Specs First" principle.

## 2. Out of Scope (Anti-Goals)
- Do not implement a server-side validation for these rules (yet).
- Do not change the underlying Specforce CLI command logic.
- Do not add new Specforce skills in this iteration.

## 3. Acceptance Criteria (BDD)

### [REQ-1] Enhanced Proactive Mandate in AGENTS.md
**User Story:** AS A developer, I WANT the AI agent to explicitly know which Specforce skill to activate based on my intent, SO THAT the SDD protocol is always respected.

**Scenarios:**
1. **[Discovery Trigger]** GIVEN the user expresses a vague idea, asks "how to" do something, or reports a bug/issue requiring investigation WHEN the agent processes this intent THEN the agent MUST activate `spf.discovery` to perform root cause analysis or brainstorming.
2. **[Spec Trigger]** GIVEN the user requests a new feature, a structural change, or a confirmed fix strategy WHEN the agent processes this intent THEN the agent MUST activate `spf.spec`.
3. **[Implementation Trigger]** GIVEN an approved spec exists AND the user says "go" or "implement" WHEN the agent processes this intent THEN the agent MUST activate `spf.implement`.
4. **[Direct Edit Prohibition]** GIVEN any intent that requires a codebase change WHEN a valid spec does not exist or is not approved THEN the agent MUST refuse to use `replace` or `write_file` and instead trigger `spf.discovery` (for research) or `spf.spec` (for planning).

### [REQ-2] Refined Skill Metadata for Discovery
**User Story:** AS AN AI agent, I WANT the description of Specforce skills to contain clear trigger keywords, SO THAT I can automatically match my context to the correct skill.

**Scenarios:**
1. **[Discovery Metadata]** GIVEN the `spf.discovery` skill definition WHEN inspected by the agent THEN its description MUST include keywords like "brainstorm", "research", "diagnose", "bug investigation", "root cause analysis", and "vague intent".
2. **[Spec Metadata]** GIVEN the `spf.spec` skill definition WHEN inspected by the agent THEN its description MUST include keywords like "plan", "initialize", "update specs", and "formalize".

### [REQ-3] Clarification of "Primary Orchestration Only" in Engineering Docs
**User Story:** AS A Solutions Architect, I WANT the engineering guidelines to distinguish between "LLM native planning" and "Specforce orchestration", SO THAT agents don't mistakenly skip `spf.spec`.

**Scenarios:**
1. **[Guideline Distinction]** GIVEN the `engineering.md` document WHEN an agent reads the "Agent Orchestration Protocol" THEN it MUST see an explicit instruction that `spf.spec` is the mandatory orchestration tool and ONLY internal "Blackbox" planning modes are prohibited.

## 4. Business Invariants
- No implementation is allowed without an approved `tasks.md` marked for implementation.
- `AGENTS.md` is the "source of truth" for the agent's behavior protocol.
