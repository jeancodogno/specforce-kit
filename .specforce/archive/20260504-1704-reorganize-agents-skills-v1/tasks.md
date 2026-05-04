---
slug: 20260504-1704-reorganize-agents-skills-v1
lens: Backend-heavy
---

# Implementation Roadmap: Reorganized Agents & Skills Structure (v1.x)

## 1. Execution Strategy
- **Gravity Order:** Core Registry (Discovery) -> Schema Update (kit.yaml) -> Translator Refactor (Data-driven) -> Instruction Manager (Variable Injection).

## 2. Tasks

### Phase 1: Core Registry & Discovery

- [x] T1.1: [CODE] Multi-Source Registry Initialization
**Target:** `src/internal/agent/registry.go`
**Context:** [REQ-1]

**Action Steps:**
- Update `Registry` struct to store `localAgents` and `localSkills`.
- Modify `Initialize` to accept `rootDir` and scan `.specforce/kit.yaml` and `.specforce/skills/`.
- Implement merging logic where local entities override embedded ones with the same ID.

**Verification (TDD):**
- Run `go test ./src/internal/agent/registry_test.go`. Add a test case that initializes with a mock FS containing a local override and asserts the local version is returned.

- [x] T1.2: [CODE] Versioned Skill Metadata
**Target:** `src/internal/agent/manifest.go`
**Context:** [REQ-4]

**Action Steps:**
- Add `Version` field to `SkillMetadata` struct.
- Update parsing logic to extract version from `SKILL.yaml` or `SKILL.md` frontmatter.

**Verification (TDD):**
- Add test in `manifest_test.go` to verify version parsing from a sample YAML/Markdown skill.

### Phase 2: Standardized Schema & Translator

- [x] T2.1: [SCAFFOLD] Update Kit YAML Schema
**Target:** `src/internal/agent/kit/kit.yaml`
**Context:** [REQ-2]

**Action Steps:**
- Refactor the YAML structure to include `security: { global_write: true/false }` for each agent.
- Define default mappings for `agents`, `commands`, and `skills` categories.

**Verification (TDD):**
- Run `go test ./src/internal/agent/kit_manifests_test.go` to ensure the new schema parses correctly.

- [x] T2.2: [CODE] Data-Driven Translator
**Target:** `src/internal/agent/translator.go`
**Context:** [REQ-2]

**Action Steps:**
- Remove `globalEnabledAgents` slice.
- Update `isGlobalEnabled` to check the agent's `security` metadata in `KitConfig`.
- Refactor `resolveMappings` to use default mappings if the category is missing in the agent's specific `mappings`.

**Verification (TDD):**
- Run `go test ./src/internal/agent/translator_test.go`. Assert that an agent marked with `global_write: true` in the mock config can write to global paths.

### Phase 3: Instruction Variable Injection

- [x] T3.1: [CODE] Implement InstructionManager
**Target:** `src/internal/agent/instructions.go`
**Context:** [REQ-3]

**Action Steps:**
- Create `InstructionManager` struct and `GetInstructions` method.
- Implement regex-based replacement for `{{key}}` using a provided context map.

**Verification (TDD):**
- Create `instructions_test.go`. Assert that `GetInstructions` correctly replaces `{{project_name}}` with the expected string.

- [x] T3.2: [CODE] Integrate Instruction Manager with CLI
**Target:** `src/internal/cli/spec.go`
**Context:** [REQ-3]

**Action Steps:**
- Update `HandleSpecArtifact` and similar handlers to use `InstructionManager` before outputting content to the agent.

**Verification (TDD):**
- Run `specforce spec artifact requirements --json` and verify that any injected variables in the template are resolved.

### Phase 4: Final Configuration & Documentation

- [x] T4.1: [SCAFFOLD] Commented-out Context in Default Config
**Target:** `src/internal/core/config.go`
**Context:** [REQ-5]

**Action Steps:**
- Update `DefaultConfigContent` to make the `context` block commented out or clearly marked as an example.

**Verification (TDD):**
- Verify `src/internal/core/config.go` has the updated `DefaultConfigContent`.

- [x] T4.2: [DOCS] Update Configuration Documentation
**Target:** `docs/en/configuration.md`
**Context:** [REQ-5]

**Action Steps:**
- Add a section for "Instruction Variable Injection" explaining how to use `context` and `{{var}}`.

**Verification (TDD):**
- Verify the file exists and contains the new section.
