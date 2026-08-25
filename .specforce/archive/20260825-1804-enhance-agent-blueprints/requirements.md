---
slug: 20260825-1804-enhance-agent-blueprints
lens: Integration
---

# Feature: Enhance Agent Command Blueprints

## 1. Context & Value
Agent command blueprints define the behavioral instructions for AI agents across different phases of Spec-Driven Development. By enriching Discovery with proactive recommendations, establishing worker reuse feedback loops in Implementation, and appending conceptual next steps in Archival, the agent workflow becomes more collaborative, cost-effective, and actionable.

## 2. Out of Scope (Anti-Goals)
- Do not modify CLI command execution flags or parser logic in Go binaries.
- Do not generate automated shell script execution blocks in the archival summary.
- Do not remove existing mandatory worker guardrails or test invariance directives from the implementation blueprint.
- Do not add interactive prompt questions to non-interactive CLI modes.

## 3. Acceptance Criteria (BDD)

### [US-1] Proactive Suggestions and Technical Opinions in Discovery
**User Story:** AS A developer exploring a feature or architecture, I WANT the discovery workflow to provide proactive opinions, best-practice recommendations, and concrete trade-offs, SO THAT I can make informed architectural decisions rather than answering purely passive queries.

**Scenarios:**
1. **[Happy Path]** GIVEN an active discovery session exploring a technical pattern WHEN the developer presents a problem or idea THEN the blueprint instructs the agent to adopt a consultative stance ("Strong Opinions, Weakly Held"), recommending the optimal architectural approach with concrete pros, cons, and potential pitfalls.
2. **[Edge Case]** GIVEN a proposed feature that challenges existing architectural boundaries WHEN the agent evaluates the design THEN the blueprint directs the agent to proactively suggest cleaner alternative patterns and preventative guardrails instead of silently accepting suboptimal designs.

**Technical Constraints (NFR):**
- **[Performance]:** Blueprint prompt parsing overhead < 5ms.
- **[Safety & Security]:** Discovery remains strictly read-only for application code.
- **[Integrity]:** Adherence to constitutional governance in `.specforce/docs/`.
- **[Observability]:** Clear structured trade-off recommendations in markdown output.

### [US-2] Worker Feedback Loop, Universal Harness Compatibility, and Adaptive Effort Routing in Implementation
**User Story:** AS A lead orchestrator executing an implementation roadmap, I WANT failed task verifications to be routed back to the same subagent session, worker spawning to be compatible across any harness (Claude Code, OpenCode, Antigravity, etc.), and model tier/effort to adapt dynamically to task complexity, SO THAT execution is robust, cost-effective, and fast without wasteful agent thrashing.

**Scenarios:**
1. **[Happy Path]** GIVEN a spawned worker whose implementation produces verification failures, missing code, or compilation errors WHEN the orchestrator detects the failure THEN the blueprint instructs the orchestrator to send the exact failure logs back to the same worker conversation for adjustment (up to 2 iterations) rather than spawning a new subagent.
2. **[Edge Case]** GIVEN an implementation batch with varying task complexity WHEN delegating to a worker across any harness THEN the blueprint guides the orchestrator to select explicit worker roles/tools (avoiding fragile `self` inheritance) and route the appropriate model tier and reasoning effort (`low`, `medium`, `high`) according to task risk.

**Technical Constraints (NFR):**
- **[Performance]:** Eliminates redundant agent initialization and optimizes token usage via tier/effort routing.
- **[Safety & Security]:** Subagents remain strictly forbidden from altering `tasks.md` or task states.
- **[Integrity]:** Preservation of all existing test invariance and non-negotiable worker guardrails.
- **[Observability]:** Transparent logging of persona, model tier, effort level, and verification feedback iterations.

### [US-3] Conceptual Next Steps and Spec Suggestions in Archival
**User Story:** AS A developer closing a feature lifecycle, I WANT the archival workflow to suggest logical next steps and follow-up specifications, SO THAT project momentum and architectural continuity are maintained.

**Scenarios:**
1. **[Happy Path]** GIVEN a completed feature undergoing archival WHEN the spec is successfully archived and module living specs are consolidated THEN the blueprint instructs the agent to append 2-3 conceptual suggestions for next steps, follow-up specs, or technical debt resolutions.
2. **[Edge Case]** GIVEN an archived feature with minimal cross-cutting impact WHEN next steps are generated THEN the blueprint directs the agent to provide focused domain-specific follow-ups without hallucinating unrelated system requirements or emitting executable CLI commands.

**Technical Constraints (NFR):**
- **[Performance]:** Archival output generation < 100ms after CLI command completion.
- **[Safety & Security]:** Archival execution remains mediated by `specforce spec archive <slug>`.
- **[Integrity]:** Consistency between command blueprint (`archive.yaml`) and instruction text (`archive.md`).
- **[Observability]:** Structured markdown section for suggested next steps.

## 4. Business Invariants
- All prompt updates must maintain backward compatibility with existing multi-agent mappings (Claude, Kimi, Kilo, OpenCode, Codex, Antigravity).
- Mandatory guardrails in `implement.yaml` (test invariance, scope containment, no direct task state editing by subagents) must remain intact.
- Archival must always execute `specforce spec archive <slug>` as the terminal action before printing handoff notes.

## 5. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Zero runtime regression in kit embedding or blueprint translation.
- **[Reliability]:** All unit and integration tests for blueprint guardrails and kit manifests must pass with 100% success.
- **[Maintainability]:** Clean YAML syntax, consistent markdown formatting, and clear developer-facing instructions.
