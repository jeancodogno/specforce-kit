---
slug: 20260904-2317-pure-skills-agent-kit
lens: Migration
---

# Implementation Roadmap: Pure Skills Agent Kit & Legacy Decommissioning

## 1. Execution Strategy
- **Gravity Order:** Kit Blueprints & Translator -> Legacy Asset Detection & Gemini Removal -> Governance & E2E Validation.
- **TDD Rigor:** Every phase enforces explicit [RED] unit/integration tests covering Happy Path and Edge Cases before the [GREEN] implementation task.

## 2. Tasks

### Phase 1: Blueprint Reorganization & Pure Skills Mapping

- [x] T1.1: [RED] Write Tests for Pure Skills Blueprint Structure and Embed Integrity
**Target:** `src/internal/agent/kit_manifests_test.go`
**Context:** [US-1]

**Action Steps:**
- Update manifest existence assertions to target `kit/skills/spf-{discovery,spec,constitution,implement,archive}/SKILL.yaml`.
- Add test assertions verifying that no files or folders exist under `kit/commands/`.
- Add validation asserting valid frontmatter and metadata for all 6 skill manifests.

**Acceptance Check:**
`go test -v -run TestKitManifests ./src/internal/agent/...` fails due to missing new skill blueprint files.

- [x] T1.2: [GREEN] Migrate Kit Blueprints to Pure Skills and Update `kit.yaml`
**Target:** `src/internal/agent/kit/kit.yaml`
**Context:** [US-1]

**Action Steps:**
- Relocate and format blueprints from `kit/commands/*.yaml` to `kit/skills/spf-*/SKILL.yaml`.
- Delete the legacy directory `src/internal/agent/kit/commands/`.
- Refactor `kit.yaml` to define standard `skills:` mappings for all agents (`antigravity`, `claude`, `cursor`, `codex`, `opencode`, `kilocode`, `qwen`, `kimi-code`) and eliminate all `commands:` and `workflows:` mapping blocks.

**Acceptance Check:**
`go test -v -run TestKitManifests ./src/internal/agent/...` passes.

- [x] T1.3: [RED] Add Test Cases for Skills-Only Adaptation and TOML Transformer Removal
**Target:** `src/internal/agent/compatibility_test.go`
**Context:** [US-1, US-2]

**Action Steps:**
- Update agent installation assertions to expect exclusively `.agents/skills/spf-*/SKILL.md`, `.claude/skills/spf-*/SKILL.md`, `.cursor/skills/spf-*/SKILL.md`, etc.
- Add negative assertions ensuring that `.agents/workflows/` and `<tool>/commands/` directories are never created.
- Add assertions verifying that Gemini CLI is no longer in the agent registry list.

**Acceptance Check:**
`go test -v -run TestCompatibility ./src/internal/agent/...` fails indicating mismatched target paths.

- [x] T1.4: [GREEN] Update Translator Engine and Strip TOML Transformer
**Target:** `src/internal/agent/translator.go`
**Context:** [US-1, US-2]

**Action Steps:**
- Remove the `".toml"` transformer from the `transformers` map.
- Remove obsolete command-specific frontmatter formatting branches in `injectYAMLHeader` where category is `commands`.
- Adjust blueprint mapping category resolution to support the standardized `kit/skills/spf-*/SKILL.yaml` hierarchy.

**Acceptance Check:**
`go test -v -run TestCompatibility ./src/internal/agent/...` passes.

### Phase 2: Gemini CLI Decommissioning & Legacy Asset Detection

- [x] T2.1: [RED] Write Tests for Workflows, Commands, and Gemini Legacy Asset Detection
**Target:** `src/internal/project/legacy_test.go`
**Context:** [US-2, US-3]

**Action Steps:**
- Add test case verifying detection of `.agents/workflows/spf-*.md` and `.agents/workflows/` directory.
- Add test case verifying detection of `<tool>/commands/` folders across supported tools.
- Add test case verifying detection and interactive deletion of `.gemini/` directory in project root.

**Acceptance Check:**
`go test -v -run TestDetectLegacyAssets ./src/internal/project/...` fails due to unhandled legacy paths.

- [x] T2.2: [GREEN] Implement Broad Legacy Asset Detection and Decommissioning
**Target:** `src/internal/project/legacy.go`
**Context:** [US-2, US-3]

**Action Steps:**
- Remove `".gemini"` from `KnownToolDirs` and add root-level `.gemini/` detection to `DetectLegacyAssets`.
- Add scanning logic in `DetectLegacyAssets` for `.agents/workflows` and `<tool>/commands` across all known tool directories.
- Ensure `CleanupLegacyAssets` safely purges matching folders and files upon confirmation.

**Acceptance Check:**
`go test -v -run TestDetectLegacyAssets ./src/internal/project/...` passes.

- [x] T2.3: [RED] Write Tests for Platform Configs and AGENTS.md Pure Skills Protocol
**Target:** `src/internal/project/agents_md_test.go`
**Context:** [US-1, US-2]

**Action Steps:**
- Add test asserting that `EnsureAgentsMD` never creates `.gemini/` or `.gemini/settings.json`.
- Add test verifying that generated `AGENTS.md` contains the updated Pure Skills protocol and no legacy workflow/command slash directives.
- Add test asserting symlink integrity for supported tools without Gemini.

**Acceptance Check:**
`go test -v -run TestEnsureAgentsMD ./src/internal/project/...` fails against outdated expectations.

- [x] T2.4: [GREEN] Clean Platform Configs, Constants, and Update AGENTS.md Template
**Target:** `src/internal/project/agents_md.go`
**Context:** [US-1, US-2]

**Action Steps:**
- Remove `".gemini/"` from `ToolPrefixes` in `src/internal/core/constants.go`.
- Remove Gemini directory creation and `settings.json` generation logic from `ensurePlatformConfigs`.
- Rewrite Section 1 of `agentsMDTemplate` to mandate Specforce Skills (`spf.discovery`, `spf.spec`, `spf.constitution`, `spf.implement`, `spf.archive`) and decommission legacy command/workflow instructions.

**Acceptance Check:**
`go test -v -run TestEnsureAgentsMD ./src/internal/project/...` passes.

### Phase 3: Verification & Governance Alignment

- [x] T3.1: [TEST] Run Complete Test Suite and End-to-End Synchronization
**Target:** `Global Scope`
**Context:** [US-1, US-2, US-3]

**Action Steps:**
- Run `go test ./...` to verify zero regression across all packages.
- Clean up any residual test fixture references to gemini or command workflows.
- Run `specforce init` on a temporary directory to verify that only skills directories are created.

**Acceptance Check:**
`go test ./...` exits with code 0 and 100% passing tests across all modules.

- [x] T3.2: [DOC] Update Project Governance and Living Spec Documentation
**Target:** `.specforce/docs/modules/agent-kit.md`
**Context:** [US-1, US-2]

**Action Steps:**
- Update `[BR-KIT-01]` in `.specforce/docs/modules/agent-kit.md` to establish the 100% Pure Skills standard.
- Remove Gemini CLI references from integration surfaces and documentation.
- Update `engineering.md` and user-facing docs to reflect skill-only directory layouts.

**Acceptance Check:**
`specforce spec status 20260904-2317-pure-skills-agent-kit --json` reports 100% progress and valid specification.
