---
slug: 20260504-1704-reorganize-agents-skills-v1
lens: Backend-heavy
---

# Feature: Reorganized Agents & Skills Structure (v1.x)

## 1. Context & Value
As Specforce moves towards v1.x, the current internal structure for managing AI agents and their skills is fragmented across hardcoded logic and repetitive configuration. This feature standardizes the "Kit" format to enable dynamic discovery, easier extensibility by users, and a more robust instruction-injection pipeline.

## 2. Out of Scope (Anti-Goals)
- This spec does NOT include a GUI for managing agents.
- It does NOT cover the implementation of new AI models (it is model-agnostic).
- It does NOT change the TUI presentation of specifications, only the underlying agent orchestration.

## 3. Acceptance Criteria (BDD)

### [REQ-1] Unified Discovery Registry
**User Story:** AS A developer, I WANT the system to automatically discover agents and skills from both embedded defaults and local overrides, SO THAT I can customize my AI workflow without recompiling.

**Scenarios:**
1. **[Happy Path]** GIVEN a local `.specforce/kit.yaml` exists WHEN the CLI starts THEN the local configuration merges with or overrides the embedded defaults.
2. **[Happy Path]** GIVEN a new skill directory in `.specforce/skills/` WHEN an agent performs skill discovery THEN the new skill is visible and usable.
3. **[Edge Case]** GIVEN a collision between an embedded skill and a local skill with the same name WHEN discovered THEN the local version MUST take precedence.

### [REQ-2] Standardized Kit Schema
**User Story:** AS A Kit maintainer, I WANT a non-repetitive configuration format, SO THAT I can define tool mappings with sensible defaults.

**Scenarios:**
1. **[Happy Path]** GIVEN a kit definition WITHOUT explicit mappings for a category WHEN installed THEN the system applies default paths (e.g., `agents/` -> `.agent/agents/`).
2. **[Edge Case]** GIVEN an invalid YAML structure in `kit.yaml` WHEN the registry initializes THEN the system MUST fail gracefully with a clear error message instead of crashing.

### [REQ-3] Instruction Variable Injection
**User Story:** AS A Lead Architect, I WANT to use variables in my global instructions, SO THAT I can inject project-specific context into agent prompts dynamically.

**Scenarios:**
1. **[Happy Path]** GIVEN an instruction file containing `{{project_name}}` WHEN fetched by an agent THEN the variable is replaced with the value from the project configuration.
2. **[Edge Case]** GIVEN a missing variable in the template WHEN injected THEN the system SHOULD leave the placeholder or use an empty string, preventing a processing failure.

### [REQ-4] Versioned Skill Lifecycle
**User Story:** AS A Developer, I WANT skills to have explicit versions, SO THAT I can ensure compatibility between agents and specific workflow rules.

**Scenarios:**
1. **[Happy Path]** GIVEN a skill definition with `version: 1.1.0` WHEN the agent info is retrieved THEN the version is correctly reported in the metadata.
2. **[Ambiguity]** If multiple versions of the same skill exist locally, the system ASSUMES the one with the highest semver (or the last one discovered) is the default unless specified otherwise.

### [REQ-5] Default Configuration & Documentation
**User Story:** AS A new user, I WANT the initial `config.yaml` to show me how to use context variables AND I WANT updated documentation, SO THAT I can quickly adopt the new structure.

**Scenarios:**
1. **[Happy Path]** GIVEN a project is initialized with `specforce init` WHEN `.specforce/config.yaml` is generated THEN it MUST include commented-out examples of the `context` block.
2. **[Happy Path]** GIVEN the new version is released WHEN I read `docs/en/configuration.md` THEN I see a detailed section explaining "Instruction Variable Injection".

## 4. Business Invariants
- **Slug Uniqueness:** Every agent and skill MUST have a globally unique identifier within the registry.
- **Read-Only Defaults:** Embedded kits MUST NEVER be modified at runtime; all customizations MUST happen via local overrides in the project directory.
- **Security Gating:** Agents discovered locally MUST still respect the `globalEnabledAgents` security list for any operations outside the project root.
