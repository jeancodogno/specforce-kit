---
slug: 20260529-1122-cursor-tool-support
lens: Backend-heavy (Integration focus)
---

# Feature: Cursor Tool Support

## 1. Context & Value
Cursor is a leading AI code editor with its own native rule system. To ensure Specforce Kit is fully compatible with Cursor's "Composer" and agentic features, we must implement a native integration that transforms Specforce kits into Cursor-compatible `.md` files organized in a structured `.cursor/` directory. This allows Cursor to automatically discover and apply Specforce's expertise, architectural rules, and SDD protocols through its standard context-loading mechanisms, significantly reducing friction for Cursor users.

## 2. Out of Scope (Anti-Goals)
- **Cursor Settings Management:** This feature does not modify Cursor's internal `settings.json` or other editor-level preferences.
- **Legacy Symlink Logic:** We will not use the legacy symlink pattern (e.g., `AGENTS.md` symlinks) for Cursor, as it natively reads the root `AGENTS.md`.
- **Legacy .mdc Format:** We will not generate `.mdc` files or use the `.cursor/rules/` directory for these artifacts.
- **Dynamic Glob Configuration:** The system will not attempt to guess or dynamically calculate file globs.

## 3. Acceptance Criteria (BDD)

### [US-1] Native Hidden Directory Registration
AS A developer, I WANT Cursor to be recognized as a first-class tool in Specforce, SO THAT its rules are correctly managed within the `.cursor/` directory.

**Acceptance Criteria:**
- **GIVEN** the Specforce core constants.
- **WHEN** the system initializes tool prefixes.
- **THEN** the `.cursor/` directory SHALL be included in the managed `ToolPrefixes` whitelist.
- **AND** the system SHALL allow I/O operations within this directory for Cursor-specific artifacts.

**UI/UX Specifics:**
- N/A (Internal registration).

**Technical Constraints (NFR):**
- The registration MUST be hardcoded in the core constants to prevent accidental omission during upgrades.

### [US-2] Structured Cursor Artifacts (.md)
AS AN AI agent in Cursor, I WANT to receive rules in a structured directory format using standard `.md` files, SO THAT I can correctly interpret the SDD protocols.

**Acceptance Criteria:**
- **GIVEN** a Specforce kit artifact (Agent, Command, or Skill).
- **WHEN** the `specforce init` or agents synchronization is executed for the Cursor tool.
- **THEN** the system SHALL generate corresponding files in the following directory structure:
    - Agents: `.cursor/agents/*.md`
    - Commands (Standalone): `.cursor/commands/spf-*.md`
    - Commands (Skill-style): `.cursor/skills/spf-*/SKILL.md`
- **AND** the files MUST use the `.md` extension.
- **AND** the body of the file MUST contain the artifact's instructions.

**UI/UX Specifics:**
- Errors during generation (e.g., missing kit data) MUST be reported as actionable warnings.

**Technical Constraints (NFR):**
- The system MUST ensure the target directories exist before writing.

### [US-3] Skill Directory Organization
AS A developer, I WANT skills to be organized into their own subdirectories within `.cursor/skills/`, SO THAT all related skill files (primary and secondary) are kept together.

**Acceptance Criteria:**
- **GIVEN** a Specforce kit skill.
- **WHEN** artifacts are translated for Cursor.
- **THEN** the system SHALL create a directory for each skill at `.cursor/skills/<skill-name>/`.
- **AND** the primary skill file SHALL be written as `.cursor/skills/<skill-name>/SKILL.md`.
- **AND** all secondary support files (e.g., references, templates) SHALL be written as standard `.md` files in the same skill folder.

**UI/UX Specifics:**
- N/A (Internal organization).

**Technical Constraints (NFR):**
- The system MUST maintain the relative directory structure for secondary files within the skill's root folder.

### [US-4] Root Context Integration (AGENTS.md)
AS A Cursor user, I WANT the system to rely on the root `AGENTS.md` for global rules, SO THAT I don't have redundant symlinks cluttering my rules directory.

**Acceptance Criteria:**
- **GIVEN** a project with a root `AGENTS.md` file.
- **WHEN** the Cursor tool integration is initialized.
- **THEN** the system SHALL NOT create a symlink or copy of `AGENTS.md` inside `.cursor/`.
- **AND** the system SHALL verify that Cursor has direct access to the root `AGENTS.md` (by virtue of it being in the root).

**UI/UX Specifics:**
- N/A (Cleanup/Omission operation).

**Technical Constraints (NFR):**
- Any legacy logic that automatically creates symlinks for other tools MUST be bypassed for the Cursor tool.

## 4. Business Invariants
- **Discovery Parity:** A kit item that appears as an agent in other tools (like Claude) MUST appear as a corresponding `.md` file in the correct `.cursor/` sub-directories to maintain SDD protocol consistency.
- **No Format Pollution:** No tool-specific frontmatter (like `.mdc` YAML) SHALL be added to these `.md` files.

## 5. Success Metrics
- **Business Metric:** 100% of Cursor-based SDD sessions correctly load Specforce rules via the `.cursor/` structured directory.
- **Performance Target:** Generation of 10+ artifacts during `init` SHALL complete in < 200ms.
- **UX Efficiency:** Zero manual steps for the user to "install" Specforce rules into Cursor; they appear automatically after `specforce init`.

## 6. Global Non-Functional Requirements (NFRs)
- **Path Portability:** All paths within `.cursor/` MUST be relative to the project root.
- **Overwrite Safety:** The system SHALL overwrite existing files in the managed `.cursor/` sub-directories during synchronization to ensure they reflect the latest kit state.
- **Atomic Writes:** Each file transformation MUST be atomic.
