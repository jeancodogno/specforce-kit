---
slug: 20260822-1614-zero-tech-requirements-rules
lens: Backend-heavy
---

# Feature: Strict Business Requirements & Agent Kit Consolidation with Legacy Cleanup

## 1. Context & Value
Specforce Kit is streamlining its agent architecture and strengthening specification quality gates. We are removing redundant legacy agents and skills from the embedded kit (retaining exclusively the sovereign orchestrators and `consultative-grill`), while simultaneously enforcing a strict Zero-Technical-Specification policy in `requirements.md`. Additionally, `specforce init` will proactively detect legacy agents and skills in existing workspaces and prompt the user for permission to clean them up.

## 2. Out of Scope (Anti-Goals)
- Removing core CLI commands (`spec`, `constitution`, `discovery`, `implement`, `archive`).
- Deleting user-created custom skills that do not collide with deprecated built-in Specforce kit items.
- Auto-deleting files without explicit user confirmation when running interactively.

## 3. Acceptance Criteria (BDD)

### [US-1] Strict Zero-Technical-Specification Policy in Orchestrator Prompts
**User Story:** AS A Solutions Architect or Product Lead, I WANT the Specforce feature orchestrator to strictly enforce business-only content during `requirements.md` generation, SO THAT AI planners never pollute functional specs with implementation mechanics.

**Scenarios:**
1. **[Happy Path]** GIVEN the orchestrator is drafting `requirements.md` WHEN generating user stories and BDD scenarios THEN it uses only end-user/business personas and business-level domain outcomes without technical transport, schemas, or status codes.
2. **[Edge Case - Technical Detail Identified]** GIVEN the user or codebase exploration reveals specific technical endpoints, database schemas, or code libraries WHEN the orchestrator compiles `requirements.md` THEN it excludes these technical details from `requirements.md` and defers them entirely to `design.md`.

**UI/UX Specifics:**
- **View/Component:** Command generation prompt.
- **Feedback Logic:** Clear directive instructions presented to the agent.
- **Keybindings:** N/A.

**Technical Constraints (NFR):**
- **[Performance]:** Zero added latency during command execution or template embedding.
- **[Safety & Security]:** Idempotent prompt declarations.
- **[Integrity]:** Consistent rule definitions across all agent kit command templates.

### [US-2] Agent Kit Consolidation (Single Core Skill & Deprecated Agents Removal)
**User Story:** AS A Developer or AI Agent, I WANT the embedded Specforce kit to contain only the `consultative-grill` skill and no standalone legacy agent definitions, SO THAT the agent architecture is lean, unified, and free of obsolete sub-agent roles.

**Scenarios:**
1. **[Happy Path]** GIVEN the agent kit embedded resources WHEN querying embedded skills THEN only `consultative-grill` is present (legacy skills `pragmatic-product-owner`, `opportunity-framing`, and `task-atomic-decomposition` are removed).
2. **[Happy Path - Agents]** GIVEN the agent kit embedded resources WHEN inspecting the embedded `agents/` directory THEN all legacy subagents (`specforce-architect`, `specforce-developer`, `specforce-planner`, `specforce-product-analyst`, `specforce-qa`, `specforce-spec-reviewer`) are removed.

**UI/UX Specifics:**
- **View/Component:** Embedded kit filesystem.
- **Feedback Logic:** Lean skill and agent asset tree.
- **Keybindings:** N/A.

**Technical Constraints (NFR):**
- **[Performance]:** Reduced binary size and faster asset translation.
- **[Maintainability]:** Elimination of unmaintained skill duplicates.

### [US-3] Interactive Legacy Agent & Skill Cleanup on Init
**User Story:** AS A Developer running `specforce init`, I WANT Specforce to detect any legacy agents and skills in my project and ask if I want to delete them, SO THAT my workspace stays clean without risking accidental deletion of files I might want to keep.

**Scenarios:**
1. **[Happy Path - User Confirms Cleanup]** GIVEN an existing project containing legacy agents or skills WHEN `specforce init` executes and detects them THEN the system prompts for confirmation to delete legacy assets, and upon user confirmation ("yes"), deletes the legacy agents and skills before syncing active tools.
2. **[Alternative Path - User Declines Cleanup]** GIVEN an existing project containing legacy agents or skills WHEN `specforce init` executes and prompts for confirmation THEN upon user decline ("no"), the system retains the legacy files untouched and updates the active tools normally.
3. **[Happy Path - Clean Workspace]** GIVEN a project with no legacy agents or skills WHEN `specforce init` runs THEN it initializes or updates tools directly without prompting for legacy cleanup.

**UI/UX Specifics:**
- **View/Component:** CLI / TUI confirmation prompt.
- **Feedback Logic:** Prompt: "Legacy Specforce agents/skills detected. Would you like to remove them? [y/N]" with styled status logging.
- **Keybindings:** y/n key confirmation.

**Technical Constraints (NFR):**
- **[Safety & Security]:** Safe path resolution and no unconfirmed file deletions.
- **[Reliability]:** Non-blocking execution in non-interactive / test environments.

## 4. Business Invariants
- `requirements.md` must never contain code blocks, HTTP endpoints, SQL schemas, or programming language constructs.
- Legacy asset deletion must never occur without user consent when running interactively.
- The `consultative-grill` skill must remain fully functional.

## 6. Global Non-Functional Requirements (NFRs)
- **[Maintainability]:** 100% adherence to Specforce standard kit directory structure and YAML formats.
- **[Reliability]:** All unit and integration tests across all packages must pass cleanly.
