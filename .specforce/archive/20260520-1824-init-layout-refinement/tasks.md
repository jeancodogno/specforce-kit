---
slug: 20260520-1824-init-layout-refinement
lens: UI-heavy
---

# Implementation Roadmap: Init Layout Refinement

## 1. Execution Strategy
- **Gravity Order:** TUI Components Foundation -> Infrastructure Bootstrapper Refactoring -> CLI Orchestration & Handoff Integration

## 2. Tasks

### Phase 1: Ghost Protocol Foundation (TUI Components)

- [x] T1.1: [CODE] Refine TUI Theme and Component Styles
**Target:** `src/internal/tui/theme.go` & `src/internal/tui/components.go`
**Context:** [US-1], [US-3]

**Action Steps:**
- Verify `CleanBorder` is defined in `theme.go`. If not, add it using the Ghost Protocol specification (thin borders).
- Define `ActiveStatusStyle` (Cyan) and `ReadyStatusStyle` (Silver/Ice) in `theme.go` for the agent selection badges.
- Update `PrintCompletionBox` in `components.go` to use `CleanBorder` and align the text padding to match the "Ghost Handoff" layout.

**Verification (TDD):**
Run `go test ./src/internal/tui/...` to ensure no existing TUI component tests are broken by style changes. Create a new test case for `PrintCompletionBox` if missing, verifying border application.

- [x] T1.2: [CODE] Implement 'The Arsenal' Multiselect Layout
**Target:** `src/internal/tui/multiselect.go`
**Context:** [US-1]

**Action Steps:**
- Refactor the `viewSelection()` method to wrap the entire list in a `CleanBorder` frame with the title "SELECT AI AGENTS".
- Update the item rendering to include a status badge (e.g., `[ ACTIVE ]` in Cyan when selected, `[ READY  ]` in Silver/Ice when unselected).
- Align the output so the cursor, bullet, agent name, and status badge form a clean, readable column structure.

**Verification (TDD):**
Run `go test ./src/internal/tui/...`. If a test for `multiselect.go` UI output exists, update the expected string output to match the new bounded layout.

### Phase 2: Surgical Infrastructure Pulse

- [x] T2.1: [CODE] Enhance Bootstrapper with Granular SubTask Logging
**Target:** `src/internal/project/bootstrapper.go`
**Context:** [US-2]

**Action Steps:**
- Modify `BootstrapProject` to accept `core.UI` and utilize `ui.LogSubTask` for each directory created in the `dirs` array, instead of a single `StartSpinner`.
- Format the `LogSubTask` string to match the Pulse aesthetic (e.g., `↳ .specforce/docs ................................. OK`).

**Verification (TDD):**
Run `go test ./src/internal/project/...`. Verify that the mock UI in the tests captures the expected sequence of `LogSubTask` calls for each directory.

- [x] T2.2: [CODE] Refine Initialization Service Output
**Target:** `src/internal/project/service.go`
**Context:** [US-2]

**Action Steps:**
- Update `InitializeProject` to output a generic `ui.LogSubTask(" › DEPLOYING INFRASTRUCTURE...")` before calling `BootstrapProject`.
- Remove redundant `ui.StartSpinner` or `ui.SubTask` calls inside `InitializeProject` that might clash with the new granular Pulse feedback from `BootstrapProject` and `agent.AdaptArtifacts`.

**Verification (TDD):**
Run `go test ./src/internal/project/...` to ensure the overall initialization flow remains unbroken and the UI interactions are sequenced correctly.

### Phase 3: CLI Orchestration & Handoff Integration

- [x] T3.1: [CODE] Integrate Ghost Handoff Completion Box
**Target:** `src/internal/cli/cli.go`
**Context:** [US-3]

**Action Steps:**
- Update the `handleNewInitFlow` function to replace the generic `PrintCompletionBox` message with the finalized "Ghost Handoff" text.
- Ensure the output strictly reads:
  "Specforce structure is live.\n\nNEXT: Run '/spf:discovery' to start the SDD cycle."
- Verify the title of the completion box remains "MISSION ACCOMPLISHED".

**Verification (TDD):**
Run `go test ./src/internal/cli/...`. Update any CLI integration tests that assert against the specific output string of a successful initialization to match the new Handoff text.

## 3. Pre-emptive Mitigations
- **Risk:** Refactoring the `viewSelection` string building in `multiselect.go` may break terminal alignment across different character encodings or widths. -> **Mitigation:** Rely on Lipgloss for width calculation and padding, rather than manual string spacing (`fmt.Sprintf("%-20s", ...)`), to ensure robust alignment.