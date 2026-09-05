---
slug: 20260904-2317-pure-skills-agent-kit
lens: Migration
---

# Feature: Pure Skills Agent Kit & Legacy Decommissioning

## 1. Context & Value
Specforce currently maintains redundant configuration layers by generating both command/workflow files and agent skills across supported AI tools. By standardizing 100% on the modern Agent Skill specification and retiring deprecated tools like Gemini CLI, Specforce simplifies the developer workspace, reduces token context waste, and guarantees consistent agent behavior across all supported environments.

## 2. Out of Scope (Anti-Goals)
- Modifying the core CLI execution lifecycle commands (`specforce spec`, `specforce implementation`, `specforce constitution`).
- Altering the textual content or functional behavior of the core SDD prompts beyond directory references and skill invocation directives.
- Removing or breaking support for other existing agent harnesses (Claude, Antigravity, Cursor, Codex, OpenCode, KiloCode, Qwen, Kimi).

## 3. Acceptance Criteria (BDD)

### [US-1] Pure Skills Installation for Supported Agents
**User Story:** AS A software engineer using Specforce with AI coding assistants, I WANT TO have all Specforce capabilities installed strictly as standardized Agent Skills, SO THAT my AI tools have a single, unified method of discovering and executing Specforce workflows without duplicate configuration files.

**Scenarios:**
1. **[Happy Path]** GIVEN an engineer initializing or updating Specforce for supported AI tools WHEN the installation finishes THEN all Specforce blueprints are installed exclusively within each tool's designated skills directory under standardized subdirectories, and no workflow or command files are created.
2. **[Edge Case - Fallback Tool Mapping]** GIVEN a project where a supported tool has no custom skill overrides WHEN Specforce updates tool configurations THEN the system falls back to default skill placement without creating empty command or workflow directories.

**UI/UX Specifics:**
- **View/Component:** TUI Tool Synchronization Progress.
- **Feedback Logic:** Checkmark indicators next to each adapted skill, spinner during filesystem writes.
- **Keybindings:** Standard CLI non-blocking execution.

**Technical Constraints (NFR):**
- **[Performance]:** Skill adaptation and writing completed in < 200ms per agent harness.
- **[Safety & Security]:** Atomic file writing with restricted 0600 file permissions for generated skills.
- **[Integrity]:** Idempotent installation; re-running yields identical file content without duplicate directories.
- **[Observability]:** Progress and subtasks logged to the terminal UI during synchronization.

### [US-2] Decommissioning of Gemini CLI Integration
**User Story:** AS A development team maintaining project configurations, I WANT TO eliminate legacy Gemini CLI integration, SO THAT our project repository remains free from unused tool settings and proprietary configuration formats.

**Scenarios:**
1. **[Happy Path]** GIVEN an engineer configuring Specforce in a repository WHEN tool selection or synchronization occurs THEN Gemini CLI is not offered as a supported option and no tool-specific configuration settings or directories are created for it.
2. **[Edge Case - Pre-existing Unselected Directory]** GIVEN a repository where an unselected Gemini CLI directory already exists WHEN project initialization runs THEN no new configuration files or settings are injected into that directory.

**UI/UX Specifics:**
- **View/Component:** Multi-select Agent Selection Prompt.
- **Feedback Logic:** Gemini CLI does not appear in the available tool list.
- **Keybindings:** Arrow keys for tool navigation, space to select, enter to confirm.

**Technical Constraints (NFR):**
- **[Performance]:** Zero processing overhead for Gemini configurations during project setup.
- **[Safety & Security]:** No unintended modifications to foreign project files.
- **[Integrity]:** Clean decoupling without leaving orphaned references in core routing.
- **[Observability]:** Clean logging of only active, supported agents.

### [US-3] Legacy Asset Detection and Automated Removal
**User Story:** AS A software engineer maintaining an existing codebase, I WANT TO be alerted to and guided through the cleanup of obsolete workflow definitions, legacy command files, and deprecated tool integrations, SO THAT my project workspace stays tidy and AI agents avoid reading conflicting instructions.

**Scenarios:**
1. **[Happy Path]** GIVEN an existing project containing obsolete Specforce workflows, command files, or deprecated tool folders WHEN project initialization executes THEN the system identifies all legacy assets, presents them to the user for interactive confirmation, and safely removes them upon approval.
2. **[Declined Cleanup]** GIVEN an existing project containing obsolete Specforce assets WHEN project initialization runs AND the user declines the cleanup prompt THEN existing legacy files are left completely untouched and no data is lost.

**UI/UX Specifics:**
- **View/Component:** Interactive Cleanup Confirmation Prompt.
- **Feedback Logic:** Clear list of detected obsolete paths, red warning banner, confirmation prompt (y/N).
- **Keybindings:** 'y' or 'enter' for confirmation, 'n' or 'esc' to decline.

**Technical Constraints (NFR):**
- **[Performance]:** Detection scans complete in < 50ms across standard project roots.
- **[Safety & Security]:** Strict path validation preventing traversal outside the project root before deletion.
- **[Integrity]:** Preserves user-created custom skills and non-Specforce assets.
- **[Observability]:** Log each safely deleted path with dot-leader formatting in the terminal output.

## 4. Business Invariants
- All Specforce capabilities must be packaged exclusively as Agent Skills in `<tool-dir>/skills/<skill-slug>/SKILL.md`.
- No workflow files or command files may be emitted during initialization or update of any tool.
- Gemini CLI is decommissioned from tool discovery, installation, and project configuration management.
- Destructive cleanup of legacy assets requires explicit user confirmation or automated flag bypass.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Full initialization or update lifecycle finishes in < 500ms.
- **[Reliability]:** Zero-failure fallback to default skill routes for all supported agents.
- **[Security]:** All path generation rigorously sandboxed to prevent arbitrary directory traversal.
- **[Maintainability]:** Comprehensive automated test coverage (>85%) across agent adaptation and legacy detection modules.
