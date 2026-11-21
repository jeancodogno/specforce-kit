---
slug: 20260505-0011-fix-init-tool-folders
lens: Backend-heavy
---

# Feature: Fix tool folder creation during init

## 1. Context & Value
Currently, `specforce init` unconditionally creates `.gemini/`, `.claude/`, and `.agent/` directories regardless of whether the user selected the corresponding AI agents. This leads to project clutter and unused configuration artifacts. This feature ensures that platform-specific directories and their internal configurations (settings, symlinks) are only created for the tools explicitly selected by the user.

## 2. Out of Scope (Anti-Goals)
* Do not modify the core `AdaptArtifacts` logic (it already respects agent selection).
* Do not change the interactive TUI flow or agent selection mechanism.
* Do not modify the content of `AGENTS.md` itself.
* Do not remove existing directories if they already exist (only avoid creating them if they shouldn't be there).

## 3. Acceptance Criteria (BDD)

### [REQ-1] Conditional Platform Directory Creation
**User Story:** AS A developer, I WANT Specforce to only create platform-specific directories if I've selected the corresponding tool, SO THAT my project root remains clean and relevant to my tech stack.

**Scenarios:**
1. **[Only Selected Tools Created]** GIVEN a project being initialized with only "gemini-cli" selected WHEN `specforce init` completes THEN the `.gemini/` directory MUST exist AND the `.claude/` and `.agent/` directories MUST NOT exist.
2. **[Symlinks Only for Selected Tools]** GIVEN a project being initialized with "claude" and "antigravity" selected WHEN `specforce init` completes THEN the `.claude/rules/AGENTS.md` and `.agent/rules/AGENTS.md` symlinks MUST exist AND the `.gemini/` directory MUST NOT exist.
3. **[Update Flow Respects Selection]** GIVEN an existing project initialized with "gemini-cli" WHEN `UpdateTools` is executed (via `specforce init` on an existing project) THEN it MUST NOT create `.claude/` or `.agent/` if they were not selected in the TUI.

### [REQ-2] Configuration File Integrity
**User Story:** AS A developer, I WANT configuration files to be created only within existing or intended tool directories, SO THAT I don't have dangling config files for unused tools.

**Scenarios:**
1. **[Settings JSON Guarded]** GIVEN a project where "gemini-cli" is NOT selected WHEN the platform configuration logic runs THEN `.gemini/settings.json` MUST NOT be created AND the `.gemini/` directory MUST NOT be created by this specific logic.
2. **[Symlink Guarded]** GIVEN a project where "claude" is NOT selected WHEN the platform configuration logic runs THEN the `.claude/rules/AGENTS.md` symlink MUST NOT be created.
3. **[Existing Directory Detection]** GIVEN a project where a directory (e.g., `.claude/`) already exists manually WHEN platform configuration runs THEN Specforce SHOULD ensure the relevant config/symlink exists within it, even if the tool wasn't selected in the current run (to maintain consistency for partially initialized projects).

## 4. Business Invariants
* `AGENTS.md` in the project root must always be created or updated regardless of specific tool selection.
* Platform configuration logic MUST NOT use `os.MkdirAll` for tool directories unless that tool is confirmed for initialization.
* The mapping between agents and their target directories MUST remain consistent with `kit.yaml` (e.g., `gemini-cli` -> `.gemini`, `claude` -> `.claude`, `antigravity` -> `.agent`).
