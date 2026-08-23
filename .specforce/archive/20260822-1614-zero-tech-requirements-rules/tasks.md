---
slug: 20260822-1614-zero-tech-requirements-rules
lens: Backend-heavy
---

# Implementation Roadmap: Strict Business Requirements & Agent Kit Consolidation with Legacy Cleanup

## 1. Execution Strategy
- **Gravity Order:** Update prompt constraints (`spec.yaml`, `requirements.yaml`) -> Purge legacy skills & agents from kit (`kit/`) -> Implement legacy detector & cleanup service in `src/internal/project/` -> Integrate interactive confirmation in `specforce init` -> Comprehensive test verification.

## 2. Tasks

### Phase 1: Agent Kit Consolidation & Prompt Rules

- [x] T1.1: [PROMPT/CONFIG] Add strict zero-technical-specification policy to orchestrator command
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [US-1]
**Parallel With:** T1.2

**Action Steps:**
- Add Rule 6 to `**CRITICAL RULES:**` declaring the `STRICT ZERO-TECHNICAL-SPECIFICATION POLICY`.
- Mandate business personas for User Stories and domain-only outcomes for BDD scenarios.
- Reinforce that all technical implementation details belong strictly in `design.md`.

**Acceptance Check:**
`go test ./src/internal/agent/...`

- [x] T1.2: [PROMPT/CONFIG] Update requirements artifact instruction and template guidance
**Target:** `src/internal/agent/artifacts/spec/requirements.yaml`
**Context:** [US-1]
**Parallel With:** T1.1

**Action Steps:**
- Refine rule 2 in `instruction` to declare `Zero Technical Specification Policy`.
- Explicitly forbid endpoints, HTTP status codes, SQL queries, database table names, and JSON keys in requirements.md.
- Ensure instructions mandate business personas for User Stories.

**Acceptance Check:**
`go test ./src/internal/agent/...`

- [x] T1.3: [REFACTOR/DELETE] Remove legacy skills and agents from embedded kit
**Target:** `src/internal/agent/kit`
**Context:** [US-2]

**Action Steps:**
- Delete directories `src/internal/agent/kit/skills/opportunity-framing`, `src/internal/agent/kit/skills/pragmatic-product-owner`, and `src/internal/agent/kit/skills/task-atomic-decomposition`.
- Delete all agent YAML files in `src/internal/agent/kit/agents/`.
- Update `src/internal/agent/kit/kit.yaml` to clean up deprecated agent mappings.
- Update `src/internal/agent/compatibility_test.go` to test active skills (`consultative-grill`) and commands.

**Acceptance Check:**
`go test ./src/internal/agent/...`

### Phase 2: Legacy Asset Detection & Interactive Cleanup

- [x] T2.1: [CODE] Implement legacy asset detector and cleanup module
**Target:** `src/internal/project/legacy.go`
**Context:** [US-3]

**Action Steps:**
- Define slices for `LegacyAgents` and `LegacySkills`.
- Implement `DetectLegacyAssets(root string) ([]string, error)` scanning tool directories (`.agents`, `.cursor`, `.claude`, `.gemini`, `.qwen`, `.opencode`, `.kilocode`, `.codex`).
- Implement `CleanupLegacyAssets(root string, assetPaths []string, ui core.UI) error` removing confirmed assets safely.

**Acceptance Check:**
`go test ./src/internal/project/...`

- [x] T2.2: [CODE] Integrate legacy detection and confirmation prompt into project service
**Target:** `src/internal/project/service.go`
**Context:** [US-3]

**Action Steps:**
- In `InitializeProject` and `UpdateTools`, call `DetectLegacyAssets(projectRoot)`.
- If legacy assets are found, prompt user with `ui.Confirm("Legacy Specforce agents/skills detected. Do you want to remove them?")`.
- If user confirms, call `CleanupLegacyAssets` to remove them; otherwise log and continue without removing.

**Acceptance Check:**
`go test ./src/internal/project/...` and `go test ./src/internal/cli/...`

### Phase 3: Integration Verification & Tests

- [x] T3.1: [TEST] Add comprehensive unit tests for legacy cleanup and run full test suite
**Target:** `src/internal/project/legacy_test.go`
**Context:** [US-3]

**Action Steps:**
- Write test cases for detection when legacy files are present vs absent.
- Write test cases verifying confirmation acceptance and rejection behavior.
- Run full test suite across all packages.

**Acceptance Check:**
`go test ./...`
