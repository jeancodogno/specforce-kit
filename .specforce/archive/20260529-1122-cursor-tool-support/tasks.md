---
slug: 20260529-1122-cursor-tool-support
lens: Backend-heavy (Integration focus)
---

# Implementation Roadmap: Cursor Tool Support

## 1. Execution Strategy
- **Gravity Order:** Core Constants & Types -> Translation Engine (Dynamic Paths) -> Kit Configuration -> Documentation.
- **Verification:** TDD-focused. Each logic change must be verified with a Go unit test before integration.

## 2. Tasks

### Phase 1: Core Configuration & Schema

- [x] T1.1: [CODE] Register `.cursor/` in ToolPrefixes whitelist
**Target:** `src/internal/core/constants.go`
**Context:** [US-1]

**Action Steps:**
- Add `".cursor/"` to the `ToolPrefixes` slice.
- Ensure the slice remains sorted or logically grouped with other tool prefixes.
- Verify that the constant is accessible from the internal/agent package.

**Acceptance Check:**
Run `go test ./src/internal/core/...` and verify that the `ToolPrefixes` slice includes the new directory.

- [x] T1.2: [CODE] Extend `MappingConfig` with `UseSubdir` field
**Target:** `src/internal/core/blueprint.go`
**Context:** [US-3]

**Action Steps:**
- Locate the `MappingConfig` struct definition in `src/internal/core/blueprint.go`.
- Add the `UseSubdir` boolean field with the `yaml:"use_subdir,omitempty"` tag.
- Ensure the struct still implements any necessary unmarshaling logic.

**Acceptance Check:**
Compile the project: `go build ./src/...` and verify no syntax errors.

### Phase 2: Translation Engine

- [x] T2.1: [CODE] Implement Dynamic Path Resolution Logic
**Target:** `src/internal/agent/translator.go`
**Context:** [US-3]

**Action Steps:**
- Modify `resolveMappings` in `src/internal/agent/translator.go` to intercept the `UseSubdir` flag.
- Implement logic to inject the artifact slug as a subdirectory name when the flag is true.
- Ensure the final path joins the base path, the subdirectory, and the filename correctly for all artifacts.

**Acceptance Check:**
Add a test case to `src/internal/agent/translator_test.go` that verifies the subdirectory creation for skills.

- [x] T2.2: [CODE] Bypass Symlinks for Cursor
**Target:** `src/internal/project/agents_md.go`
**Context:** [US-4]

**Action Steps:**
- Locate the `ensurePlatformConfigs` function in `src/internal/project/agents_md.go`.
- Identify the `agentMappings` map used for symlink generation.
- Ensure that the Cursor tool (targeting `.cursor`) is explicitly excluded from this map or its processing loop.

**Acceptance Check:**
Run `go test ./src/internal/project/...` and verify that no symlink is created for a tool targeting `.cursor/`.

### Phase 3: Kit Integration

- [x] T3.1: [CONFIG] Define Cursor Tool Mapping in kit.yaml
**Target:** `src/internal/agent/kit/kit.yaml`
**Context:** [US-2, US-3]

**Action Steps:**
- Add the `cursor` tool configuration block to `src/internal/agent/kit/kit.yaml`.
- Map `agents` to `agents/*.md`.
- Map `commands` to BOTH `commands/*.md` (standalone) and `skills/*/SKILL.md` (skill-style with `use_subdir: true`).
- Map `skills` to `skills/*/SKILL.md` using `use_subdir: true`.

**Acceptance Check:**
Run `specforce agent list --json` and verify the `cursor` mapping is correctly loaded.

- [x] T3.2: [CLI] Verify Full Sync for Cursor
**Target:** `Global Scope`
**Context:** [US-2, US-3]

**Action Steps:**
- Initialize a test project and run `specforce init cursor`.
- Check that `.cursor/agents/`, `.cursor/commands/`, and `.cursor/skills/` are created.
- Verify that skills contain their primary `SKILL.md` inside a named subdirectory.

**Acceptance Check:**
Inspect the generated directory structure and verify it matches the design (using plural directories).

### Phase 4: Documentation

- [x] T4.1: [DOCS] Update Supported Tools Documentation
**Target:** `docs/en/supported-tools.md`
**Context:** [US-1, US-4]

**Action Steps:**
- Add Cursor to the list of supported agents in `docs/en/supported-tools.md`.
- Briefly explain the structured `.cursor/` directory integration (agents, commands, skills) and the root `AGENTS.md` support.
- Repeat the updates for `docs/pt/supported-tools.md` and `docs/es/supported-tools.md`.

**Acceptance Check:**
Verify the documentation renders correctly and reflects the new integration across all languages.

## 3. Pre-emptive Mitigations
- **Risk:** Deeply nested skills causing path issues. -> **Mitigation:** Sanitize slug/directory names before creation.
- **Risk:** Overwriting user-defined Cursor files. -> **Mitigation:** Ensure Specforce only manages its own subdirectories within `.cursor/`.