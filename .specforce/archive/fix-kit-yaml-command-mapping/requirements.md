---
slug: fix-kit-yaml-command-mapping
lens: Backend-heavy
---

# Feature: Unique Command Names in Skill Headers

## 1. Context & Value
Currently, when mapping commands to "skills" (e.g., for Kimi Code), `kit.yaml` uses `name: "SKILL"` to ensure the filename is `SKILL.md`. However, this causes the generated YAML frontmatter inside these files to also use `name: SKILL`, resulting in name collisions and poor discovery. This feature ensures each command retains its unique identity in the header while maintaining the required directory/file structure.

## 2. Out of Scope (Anti-Goals)
- Do not change the existing `kit.yaml` schema unless strictly necessary for the fix.
- Do not refactor the entire `translator.go` logic; keep changes surgical.
- Do not modify the source `.yaml` files in `src/internal/agent/kit/commands/` unless fallback logic is insufficient.

## 3. Acceptance Criteria (BDD)

### [REQ-1] Unique Name in Frontmatter
**User Story:** AS AN AI agent developer, I WANT each command skill to have a unique name in its header, SO THAT the agent can correctly identify and distinguish between different commands.

**Scenarios:**
1. **[Happy Path]** GIVEN a command mapped with `name: "SKILL"` WHEN the artifact is generated THEN the `name` in the YAML frontmatter MUST be the unique command name (e.g., `spf.spec`) instead of `"SKILL"`.
2. **[Edge Case]** GIVEN a command that already defines a `name` in its own metadata WHEN generated THEN the metadata name MUST take precedence over any mapping-derived name.

### [REQ-2] Filename Preservation
**User Story:** AS A system architect, I WANT the physical filename to remain `SKILL.md` for agents that require it, SO THAT compatibility with those agents is maintained.

**Scenarios:**
1. **[Happy Path]** GIVEN a mapping with `name: "SKILL"` and `ext: ".md"` WHEN the file is written to disk THEN the filename MUST still be `SKILL.md`.

### [REQ-3] Global Name Uniqueness
**User Story:** AS A system administrator, I WANT the framework to prevent name collisions in headers, SO THAT agent discovery remains deterministic and error-free.

**Scenarios:**
1. **[Happy Path]** GIVEN multiple commands or skills WHEN generated THEN every resulting YAML frontmatter `name` MUST be unique across the entire project.
2. **[Failure State]** GIVEN a configuration that would result in two artifacts having the same header `name` WHEN processed THEN the system MUST return an error instead of generating conflicting files.

## 4. Business Invariants
- No two commands or skills should end up with the same `name` property in their generated headers.
- The `name` in the header must be suitable for agent discovery (usually prefixed with `spf.` or `spf-`).
- The system MUST validate header name uniqueness during the artifact generation/installation phase.
