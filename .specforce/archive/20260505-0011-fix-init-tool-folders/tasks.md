---
slug: 20260505-0011-fix-init-tool-folders
lens: Backend-heavy
---

# Implementation Roadmap: Fix tool folder creation during init

## 1. Execution Strategy
- **Gravity Order:** Test Refactoring -> Core Domain Logic Refactoring -> Service Orchestration -> Final Verification.

## 2. Tasks

### Phase 1: Preparation & Test Refactoring

- [x] T1.1: [TEST] Create reproduction test for conditional initialization
**Target:** `src/internal/project/service_test.go`
**Context:** [REQ-1]

**Action Steps:**
- Add `TestService_InitializeProject_ConditionalCreation` to assert that `.gemini/`, `.claude/`, and `.agent/` are NOT created when unselected during project initialization.

**Verification (TDD):**
Run `go test -v -run TestService_InitializeProject_ConditionalCreation src/internal/project/service_test.go`. Should fail because directories are currently created unconditionally.

- [x] T1.2: [TEST] Prepare `agents_md_test.go` for signature changes
**Target:** `src/internal/project/agents_md_test.go`
**Context:** [REQ-1], [REQ-2]

**Action Steps:**
- Update calls to `EnsureAgentsMD` and `ensurePlatformConfigs` in tests to include an empty `[]string{}` as the `selectedAgents` argument.

**Verification (TDD):**
Run `go test -v src/internal/project/agents_md_test.go`. Should fail to compile until Phase 2 tasks are started, or can be done incrementally.

### Phase 2: Core Logic Implementation

- [x] T2.1: [CODE] Update `EnsureAgentsMD` signature to accept selected agents
**Target:** `src/internal/project/agents_md.go`
**Context:** [REQ-1]

**Action Steps:**
- Change signature of `EnsureAgentsMD` to `func EnsureAgentsMD(root string, ui core.UI, selectedAgents []string) error`.
- Update the internal call to `ensurePlatformConfigs` to pass `selectedAgents`.

**Verification (TDD):**
Code should fail to compile at call sites in `bootstrapper.go` and `service.go`.

- [x] T2.2: [CODE] Implement conditional logic in `ensurePlatformConfigs`
**Target:** `src/internal/project/agents_md.go`
**Context:** [REQ-1], [REQ-2]

**Action Steps:**
- Update `ensurePlatformConfigs` to accept `selectedAgents []string`.
- For each platform (.gemini, .claude, .agent), check if the tool is in `selectedAgents` OR if the directory already exists before calling `os.MkdirAll` or creating config files/symlinks.

**Verification (TDD):**
Update `TestEnsurePlatformConfigs` in `agents_md_test.go` with specific selections and verify it passes.

- [x] T2.3: [CODE] Refactor `BootstrapProject` to remove premature `EnsureAgentsMD` call
**Target:** `src/internal/project/bootstrapper.go`
**Context:** [REQ-1]

**Action Steps:**
- Remove the call to `EnsureAgentsMD` from `BootstrapProject`. This ensures `AGENTS.md` and platform configs are not created before user selection is known.

**Verification (TDD):**
Run `go test -v src/internal/project/bootstrapper_test.go`. Mock expectations might need adjustment.

- [x] T2.4: [CODE] Orchestrate `EnsureAgentsMD` in `Service` methods
**Target:** `src/internal/project/service.go`
**Context:** [REQ-1]

**Action Steps:**
- In `InitializeProject`, call `EnsureAgentsMD(root, ui, config.SelectedAgents)` after `BootstrapProject` returns successfully.
- In `UpdateTools`, ensure `EnsureAgentsMD(root, ui, selectedAgents)` is called with the current selection.

**Verification (TDD):**
Run `go test -v src/internal/project/service_test.go`. All initialization tests should pass.

### Phase 3: Final Integration & Verification

- [x] T3.1: [VERIFY] Verify reproduction test success
**Target:** `src/internal/project/service_test.go`
**Context:** [REQ-1]

**Action Steps:**
- Run the test case created in T1.1.

**Verification (TDD):**
`go test -v -run TestService_InitializeProject_ConditionalCreation src/internal/project/service_test.go` passes.

- [x] T3.2: [TEST] Verify existing directory detection (REQ-2 AC-3)
**Target:** `src/internal/project/agents_md_test.go`
**Context:** [REQ-2]

**Action Steps:**
- Add a test case where `selectedAgents` is empty but a tool directory (e.g., `.gemini/`) already exists.
- Assert that the configuration files (e.g., `settings.json`) are still created/updated within that directory.

**Verification (TDD):**
Run the new test case and verify it passes.

- [x] T3.3: [VERIFY] Regression testing
**Target:** `src/internal/project/`
**Context:** [REQ-1], [REQ-2]

**Action Steps:**
- Run all tests in the project package to ensure no side effects on other features.

**Verification (TDD):**
`go test ./src/internal/project/...` returns all PASS.
