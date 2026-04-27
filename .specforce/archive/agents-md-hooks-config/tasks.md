---
slug: agents-md-hooks-config
lens: Backend-heavy
---

# Implementation Roadmap: AGENTS.md Hooks and Platform Configuration

## 1. Execution Strategy
- **Gravity Order:** Update Internal Template -> Implement Platform Config Logic -> Update Tests -> Verify -> Update Public Documentation.

## 2. Tasks

### Phase 1: Template and Core Logic

#### T1.1: [CODE] Update AGENTS.md Template Hooks
**State:** `[FINISHED]`
**Target:** `src/internal/project/agents_md.go`
**Context:** [REQ-1]

**Action Steps:**
- Update `agentsMDTemplate` to replace outdated hook names with `on_task_finished`, `on_phase_finished`, `on_all_tasks_finished`.

**Verification (TDD):**
Run `go test ./src/internal/project/...`.

#### T1.2: [CODE] Implement ensurePlatformConfigs
**State:** `[FINISHED]`
**Target:** `src/internal/project/agents_md.go`
**Context:** [REQ-2]

**Action Steps:**
- Create `ensurePlatformConfigs(root string) error` helper.
- Implement Gemini settings creation using an array for `fileName`.
- Implement Antigravity and Claude Code symlink creation.
- Call `ensurePlatformConfigs` from `EnsureAgentsMD`.

**Verification (TDD):**
Run `go test ./src/internal/project/...`.

### Phase 2: Testing and Verification

#### T2.1: [TEST] Add Platform Config Tests
**State:** `[FINISHED]`
**Target:** `src/internal/project/agents_md_test.go`
**Context:** [REQ-2]

**Action Steps:**
- Add test cases to verify `.gemini/settings.json` content.
- Add test cases to verify symlink existence and targets for `.agent` and `.claude`.

**Verification (TDD):**
Run `go test ./src/internal/project/...` and ensure all tests pass.

#### T2.2: [CLI] Verify Final Generation
**State:** `[FINISHED]`
**Target:** `Global Scope`
**Context:** [REQ-1, REQ-2]

**Action Steps:**
- Run `go run src/cmd/specforce/main.go project refresh`.
- Verify files and symlinks exist in the project root.

**Verification (TDD):**
`ls -l .agent/rules/AGENTS.md` and `cat .gemini/settings.json`.

#### T2.3: [DOCS] Update documentation files
**State:** `[FINISHED]`
**Target:** `docs/`
**Context:** [REQ-1, REQ-2]

**Action Steps:**
- Update `docs/configuration.md` with all supported hook names.
- Update `docs/supported-tools.md` to mention automated platform configuration.

**Verification (TDD):**
Inspect markdown files.
