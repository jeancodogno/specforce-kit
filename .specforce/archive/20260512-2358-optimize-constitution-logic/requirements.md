---
slug: 20260512-2358-optimize-constitution-logic
lens: Balanced full-stack
---

# Feature: Constitution Logic Optimization

## 1. Context & Value
Current Specforce agents exhibit redundant Constitution reads in `spf.spec` and `spf.constitution`, while lacking explicit alignment commands in `spf.discovery`. This feature optimizes token usage by introducing "Short-term Memory" awareness and ensuring consistent architectural alignment across all SDD workflows.

## 2. Out of Scope (Anti-Goals)
- Changing the underlying CLI implementation of `constitution status`.
- Modifying the existing Constitution documents in `.specforce/docs/`.
- Implementing persistent server-side state for agent memory.

## 3. Acceptance Criteria (BDD)

### [US-1] Context-Sensitive Status Reuse
**User Story:** AS AN AI Agent, I WANT TO distinguish between stable and dynamic status commands, SO THAT I reuse context when safe (Constitution) but re-verify when state changes are expected (Spec/Implementation status).

**Scenarios:**
1. **[Happy Path - Stable Reuse]** GIVEN the agent's history contains a recent execution of `specforce constitution status` WHEN a command is activated AND no modifications to `.specforce/docs/` have occurred THEN the agent SHOULD reuse the existing context.
2. **[Happy Path - Dynamic Re-verification]** GIVEN the agent is in an Implementation or Planning loop WHEN starting the phase THEN the agent MUST run the relevant status command ONCE to map work, and ONCE at the end of the batch to verify completion, avoiding intermediate polling.
3. **[Edge Case]** GIVEN the agent is unsure if the context is stale WHEN performing a critical governance check THEN it SHOULD prioritize re-running the command over potentially stale history.
4. **[Efficiency Pattern]** GIVEN a multi-task implementation or planning turn THEN the agent MUST NOT run status commands between individual file writes, relying on its internal plan until the batch is finished.

**UI/UX Specifics:**
- **View/Component:** Command execution logs/Topic updates.
- **Feedback Logic:** Mentioning "Reusing stable Constitution context" vs "Re-verifying dynamic state" in summaries.

**Technical Constraints (NFR):**
- **[Performance]:** Token savings of 10-20% per session by avoiding redundant JSON parsing.
- **[Safety & Security]:** Ensures agent always aligns with the *most recent* context it has seen.
- **[Integrity]:** No modification to local state, only history-based reasoning.
- **[Observability]:** Clearly states in `update_topic` when context is reused.

### [US-2] Explicit Alignment in Discovery
**User Story:** AS A Specforce Scout, I WANT TO explicitly synchronize with the project's Constitution during the discovery phase, SO THAT my brainstorming and technical mapping are grounded in the project's established rules.

**Scenarios:**
1. **[Happy Path]** GIVEN the `spf.discovery` skill is activated WHEN the agent begins research THEN it MUST include an explicit step to check the constitution status and read relevant documents.
2. **[Edge Case]** GIVEN no constitution documents exist WHEN `spf.discovery` is run THEN the agent SHOULD suggest initializing them via `spf.constitution`.

**UI/UX Specifics:**
- **View/Component:** Discovery Intelligence Brief.
- **Feedback Logic:** Citing specific Constitution artifacts (e.g., "Aligned with architecture.md").

**Technical Constraints (NFR):**
- **[Performance]:** Context injection happens early in the turn-cycle.
- **[Safety & Security]:** Prevents proposing architectural changes that violate project principles.
- **[Integrity]:** Enforces the read-only nature of Discovery while maintaining alignment.
- **[Observability]:** Visible inclusion of `constitution status` in the Discovery command logic.

### [US-3] Optimized Instructions Distribution
**User Story:** AS A developer, I WANT global efficiency rules in AGENTS.md and workflow-specific rules in command blueprints, SO THAT agents receive relevant context without global bloat.

**Scenarios:**
1. **[Happy Path - Global Rules]** GIVEN the `AGENTS.md` file WHEN read by an agent THEN it MUST contain instructions for **Stable Context** (Constitution) reuse and surgical read patterns (`grep_search`).
2. **[Happy Path - Localized Rules]** GIVEN a specific workflow (`spf.spec` or `spf.implement`) WHEN activated THEN the agent MUST receive the **Start & End pattern** instructions directly within that command's blueprint.
3. **[Edge Case]** GIVEN an update to `AGENTS.md` THEN it MUST NOT include instructions for Spec or Implementation status polling, as those belong in their specific YAMLs.

**UI/UX Specifics:**
- **View/Component:** `AGENTS.md` and Kit YAML files.
- **Feedback Logic:** "Updating AGENTS.md..." and "Adapting artifacts..." sub-tasks.

**Technical Constraints (NFR):**
- **[Performance]:** Promotes `grep_search` over full-file reads for large documents.
- **[Safety & Security]:** Restricts instructions within the `SPECFORCE_AGENTS_START/END` markers.
- **[Integrity]:** Merging logic preserves user-defined rules outside the markers.

## 4. Business Invariants
- AI Agents MUST ALWAYS align with the Constitution before proposing changes.
- Token efficiency MUST NOT compromise the accuracy or safety of the SDD workflow.

## 5. Global UI/UX Contract (TUI Ghost Protocol)
- **Density Posture:** Compact (80x24).
- **Signature Moves:** Clear progress indicators for artifact updates.
- **Interaction Model:** Automated update flow for existing projects.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Significant reduction in redundant tool calls.
- **[Reliability]:** Synchronous update of all skill blueprints.
- **[Security]:** Strict adherence to 0600 file permissions for generated artifacts.
- **[Maintainability]:** Centralized update in `agents_md.go` and Kit YAMLs.
