---
slug: worktree-documentation
lens: Balanced full-stack
---

# Implementation Roadmap: Git Worktree Documentation

## 1. Execution Strategy
- **Gravity Order:** Documentation Base -> Detailed Content -> Entry Point References.
- **Phases:** Scaffolding, Implementation, Integration.
- **Constraints:** Documentation only, align with existing `git-worktree-support` logic.

## 2. Tasks

### Phase 1: Documentation Scaffolding

#### T1.1: [SCAFFOLD] Initialize Git Worktrees Guide
**State:** `[FINISHED]`
**Target:** `docs/git-worktrees.md`
**Context:** [REQ-1]

**Action Steps:**
- Create the file `docs/git-worktrees.md`.
- Add frontmatter (slug, lens).
- Add the main title: `# Git Worktree Support`.
- Add placeholder headers for: Introduction, Discovery Engine, Constraints, and Scanning Scope.

**Verification (TDD):**
`test -f docs/git-worktrees.md && grep -q "# Git Worktree Support" docs/git-worktrees.md`

### Phase 2: Documentation Implementation

#### T2.1: [DOC] Explain Unified Discovery Engine
**State:** `[FINISHED]`
**Target:** `docs/git-worktrees.md`
**Context:** [REQ-1]

**Action Steps:**
- Describe how Specforce automatically detects linked git worktrees.
- Explain the benefit of seeing specs from multiple branches (e.g., `main`, `feature/X`) simultaneously in the TUI console.

**Verification (TDD):**
`grep -q "Discovery Engine" docs/git-worktrees.md && grep -i "git worktree" docs/git-worktrees.md`

#### T2.2: [DOC] Detail Permissions and Read-Only Constraints
**State:** `[FINISHED]`
**Target:** `docs/git-worktrees.md`
**Context:** [REQ-1]

**Action Steps:**
- Explicitly state that specifications discovered from external worktrees (other than the current one) are **Read-Only**.
- Clarify that task modification, state transitions, and file edits are blocked for external specs.
- Note that `.specforce/docs/` (Constitution) from external worktrees are ignored to maintain local architectural integrity.

**Verification (TDD):**
`grep -q "Read-Only" docs/git-worktrees.md && grep -q "Constitution" docs/git-worktrees.md`

#### T2.3: [DOC] Define Scanning Scope
**State:** `[FINISHED]`
**Target:** `docs/git-worktrees.md`
**Context:** [REQ-1]

**Action Steps:**
- Explain what artifacts are included in the cross-worktree scan.
- List: Active Specs (`.specforce/specs/**`), Implementations (referenced in specs), and Archived Specs (`.specforce/archive/**`).

**Verification (TDD):**
`grep -q "Active Specs" docs/git-worktrees.md && grep -q "Archived Specs" docs/git-worktrees.md`

### Phase 3: README Integration

#### T3.1: [DOC] Reference Worktree Support in README
**State:** `[FINISHED]`
**Target:** `README.md`
**Context:** [REQ-2]

**Action Steps:**
- Locate the "Features" or "Documentation" section in `README.md`.
- Add a bullet point or link describing "Git Worktree Support".
- Link directly to `docs/git-worktrees.md`.

**Verification (TDD):**
`grep -q "docs/git-worktrees.md" README.md`
