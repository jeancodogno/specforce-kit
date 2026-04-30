---
slug: git-worktree-support
lens: Backend-heavy
---

# Implementation Roadmap: Git Worktree Support

## 1. Execution Strategy
The implementation follows a data-first approach: update the core state model, implement the background discovery logic, integrate it into the scanning loop, and finally update the TUI to reflect the new metadata.

## 2. Tasks

### Phase 1: Domain Model & Git Utility

#### T1.1: [SCAFFOLD] Update StateItem Struct
**State:** `[FINISHED]`
**Target:** `src/internal/spec/scanner.go`
**Context:** [REQ-2]

**Action Steps:**
- Add `Worktree string` field to the `StateItem` struct.
- Update `NewStateTree` or initialization logic if necessary.

**Verification (TDD):**
`go test ./src/internal/spec/...` (Ensure no regressions in existing scanner tests and struct alignment).

#### T1.2: [CODE] Implement Worktree Discovery
**State:** `[FINISHED]`
**Target:** `src/internal/spec/scanner.go`
**Context:** [REQ-1]

**Action Steps:**
- Implement a helper function `discoverWorktrees` that executes `git worktree list --porcelain`.
- Parse the output into a slice of structs containing `Path` and `Branch`.

**Verification (TDD):**
Add a unit test in `src/internal/spec/scanner_test.go` that mocks the git output and verifies parsing logic.

### Phase 2: Scanner Integration

#### T2.1: [CODE] Update ScanProject for Multi-Root Discovery
**State:** `[FINISHED]`
**Target:** `src/internal/spec/scanner.go`
**Context:** [REQ-1]

**Action Steps:**
- Modify `ScanProject` to first call `discoverWorktrees`.
- Iterate through each discovered worktree path.
- If the path is the current project root, execute `scanConstitution`, `scanActiveSpecs`, and `scanArchivedSpecs`.
- If the path is an external worktree, execute `scanActiveSpecs` and `scanArchivedSpecs` (skipping `scanConstitution`).
- Ensure `StateItem` instances are tagged with the `Worktree` (branch name) for external items.

**Verification (TDD):**
`go test -v src/internal/spec/scanner_test.go` with a simulated multi-worktree environment.

### Phase 3: UI Rendering

#### T3.1: [CODE] Display Worktree Labels in Console
**State:** `[FINISHED]`
**Target:** `src/internal/tui/console.go`
**Context:** [REQ-2]

**Action Steps:**
- Update the rendering logic in `renderDashboard` or the component responsible for the spec list.
- If `item.Worktree` is not empty and not the current branch, append a muted label like `[wt:branch-name]` to the display name.

**Verification (TDD):**
Manually run `specforce console` in a repo with worktrees and verify the visual output.
