---
slug: agents-md-hooks-config
lens: Backend-heavy
---

# Feature: AGENTS.md Hooks and Platform Configuration

## 1. Context & Value
The `AGENTS.md` file serves as the primary entry point for AI agents to understand the project's Specforce-driven rules. Currently, the hook examples are outdated. Additionally, Specforce needs to automatically configure environment-specific settings for Gemini, Antigravity, and Claude Code to ensure they can correctly discover and respect the project's rules.

## 2. Out of Scope (Anti-Goals)
- Changing the core hook execution logic in the Specforce binary.
- Adding support for other AI agents not specified in this request.

## 3. Acceptance Criteria (BDD)

### [REQ-1] Correct Hook Configuration Examples
**User Story:** AS A developer, I WANT TO see accurate hook configuration examples in `AGENTS.md`, SO THAT I can correctly configure my project's CI/CD and verification gates.

**Scenarios:**
1. **[Happy Path]** GIVEN the `AGENTS.md` template WHEN it is generated THEN it MUST include `on_task_finished`, `on_phase_finished`, and `on_all_tasks_finished` as hook examples.
2. **[Outdated Names]** GIVEN the current template WHEN updating THEN `on_task_finish` and `on_spec_finish` MUST be removed or replaced with the correct names.

### [REQ-2] Automatic Platform-Specific Configuration
**User Story:** AS A developer initializing or refreshing a project, I WANT Specforce to automatically set up agent-specific configuration files and links, SO THAT I don't have to perform manual setup for each agent.

**Scenarios:**
1. **[Gemini Config]** GIVEN a project refresh WHEN `EnsureAgentsMD` is called THEN Specforce MUST ensure `.gemini/settings.json` exists with `{ "context": { "fileName": ["AGENTS.md", "GEMINI.md"] } }`.
2. **[Antigravity Link]** GIVEN a project refresh WHEN `EnsureAgentsMD` is called THEN Specforce MUST ensure a symbolic link exists at `.agent/rules/AGENTS.md` pointing to `../../AGENTS.md`.
3. **[Claude Code Link]** GIVEN a project refresh WHEN `EnsureAgentsMD` is called THEN Specforce MUST ensure a symbolic link exists at `.claude/rules/AGENTS.md` pointing to `../../AGENTS.md`.

## 4. Business Invariants
- The generated `AGENTS.md` MUST always use the `<!-- SPECFORCE_AGENTS_START -->` and `<!-- SPECFORCE_AGENTS_END -->` markers.
- Hook names in the template MUST exactly match the YAML tags defined in `src/internal/core/config.go`.
- Symlinks MUST be relative to ensure portability across different environments.
