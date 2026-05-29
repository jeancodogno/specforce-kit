---
slug: 20260528-2301-antigravity-new-agent-format
lens: Integration
---

# Feature: Antigravity New Agent Format

## 1. Context & Value
Antigravity has updated its agent discovery and configuration standards. To maintain compatibility and leverage the latest performance optimizations, Specforce Kit must transition its agent management from the legacy `.agent/` hidden directory to the new `.agents/` standard. Additionally, the migration replaces manual context symlinking with a structured `agent.json` profile, enabling native discovery by the Antigravity CLI and reducing project clutter.

## 2. Out of Scope (Anti-Goals)
- **UI for Agent Management:** This feature does not include adding new TUI screens for managing agents; it focuses on the underlying storage and discovery format.
- **Legacy Support:** The system will not maintain dual support for `.agent/` after a project is migrated; `.agents/` becomes the new mandatory standard.
- **Third-Party Internal Modification:** We do not modify how Antigravity processes the data, only how Specforce provides it.

## 3. Acceptance Criteria (BDD)

### [US-1] Native Hidden Directory Migration
AS A developer, I WANT the system to use the standardized '.agents/' directory, SO THAT it aligns with the updated Antigravity CLI standards.

**Acceptance Criteria:**
- **GIVEN** a project initialized with the legacy `.agent/` directory structure.
- **WHEN** the `specforce init` or a migration command is executed.
- **THEN** the system SHALL rename or move the `.agent/` directory to `.agents/`.
- **AND** all subsequent Specforce operations MUST resolve agent resources from the new `.agents/` path.
- **AND** the system SHALL NOT crash if `.agent/` is missing but `.agents/` is present.

**UI/UX Specifics:**
- The CLI MUST provide a clear status message indicating the migration from `.agent/` to `.agents/` is complete.

**Technical Constraints (NFR):**
- Directory migration MUST be performed using atomic filesystem operations where possible to prevent data loss.

### [US-2] Agent Profile Generation (agent.json)
AS AN AI agent, I WANT my profile to be generated as a structured JSON file, SO THAT I can be correctly discovered and configured by the Antigravity CLI.

**Acceptance Criteria:**
- **GIVEN** an active Specforce agent kit (e.g., `spf.spec`, `spf.implement`).
- **WHEN** the project is initialized or the agent is synchronized.
- **THEN** an `agent.json` file SHALL be created within the respective agent's subdirectory in `.agents/`.
- **AND** the `agent.json` SHALL contain the standard Antigravity fields: `name`, `description`, `instructions`, and `tools`.
- **AND** the JSON content MUST be automatically synchronized with any changes in the Specforce kit metadata.

**UI/UX Specifics:**
- Errors during JSON generation (e.g., permission denied) MUST be displayed as actionable alerts in the TUI/CLI.

**Technical Constraints (NFR):**
- The `agent.json` file MUST be valid JSON and use 2-space indentation for human readability.

### [US-3] Automatic Context Discovery & Symlink Removal
AS A developer, I WANT the system to stop creating manual rule symlinks for Antigravity, SO THAT the project remains clean and relies on native CLI discovery.

**Acceptance Criteria:**
- **GIVEN** a migration to the new Antigravity format.
- **WHEN** the agent initialization or synchronization occurs.
- **THEN** the system SHALL NOT create symlinks for `AGENTS.md` (or other project-wide context files) inside the agent's rules directory.
- **AND** any existing legacy symlinks for `AGENTS.md` found within `.agent/` or `.agents/` MUST be removed.

**UI/UX Specifics:**
- N/A (Cleanup operation).

**Technical Constraints (NFR):**
- The removal of symlinks MUST NOT delete the source `AGENTS.md` file in the project root.

## 4. Business Invariants
- **Namespace Exclusivity:** A project cannot simultaneously use `.agent/` and `.agents/` for Specforce operations; one must be chosen as the source of truth.
- **Identity Consistency:** Moving to the new format must not change the `agent_id` or functional role of the agent as perceived by the developer.

## 5. Success Metrics
- **Business Metric:** 100% of new `specforce init` projects use the `.agents/` standard.
- **Performance Target:** Migration of an existing `.agent/` folder to `.agents/` SHALL complete in < 100ms.
- **UX Efficiency:** Zero manual intervention required for users to adopt the new Antigravity format during a routine `init` or update.

## 6. Global Non-Functional Requirements (NFRs)
- **Cross-Platform Compatibility:** Directory renaming and symlink removal MUST work consistently across POSIX and Windows filesystems.
- **Error Resilience:** If the `.agent/` to `.agents/` move fails (e.g., directory locked), the system SHALL roll back and inform the user without corrupting the existing agent configuration.
- **Idempotency:** Re-running the migration or `init` command on an already migrated project SHALL have no side effects and return success.
