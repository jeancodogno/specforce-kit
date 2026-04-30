---
slug: git-worktree-support
lens: Backend-heavy
---

# Technical Design: Git Worktree Support

## 1. Architecture Blueprint
```mermaid
graph TB
    Scanner[spec.Scanner] -- exec --> Git[Git CLI]
    Git -- list --> Worktrees[Git Worktrees]
    Scanner -- scan --> WorktreeFS[.specforce/specs in Worktrees]
    Scanner -- populate --> StateTree[spec.StateTree]
    StateTree -- provide --> Console[tui.ConsoleModel]
    Console -- render --> User((User))
```

## 4. File & Component Inventory

**Backend:**
- `src/internal/spec/scanner.go` -> Implement `discoverWorktrees` using `git worktree list --porcelain`.
- `src/internal/spec/scanner.go` -> Update `ScanProject` to iterate over detected worktrees, but execute ONLY `scanActiveSpecs` and `scanArchivedSpecs` for external worktrees (skipping `scanConstitution`).
- `src/internal/spec/scanner.go` -> Add `Worktree string` field to `StateItem` struct.

**Frontend:**
- `src/internal/tui/console.go` -> Update `renderItem` to display the `item.Worktree` label next to items not in the primary root.
