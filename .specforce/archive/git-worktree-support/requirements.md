---
slug: git-worktree-support
lens: Backend-heavy
---

# Feature: Git Worktree Support in Console

## 1. Context & Value
Developers using Git worktrees often have multiple active branches and specifications across different directories. This feature enables the Specforce Console to automatically discover and display specifications from all active git worktrees, providing a unified view of project progress regardless of the current checkout.

## 2. Out of Scope (Anti-Goals)
- No cross-worktree task modification: The console will be read-only for specifications not in the current worktree (to avoid accidental state corruption).
- No automatic worktree creation or management.
- No support for non-git version control systems.
- **RESTRICTION:** Discovery in external worktrees is limited ONLY to Active Specifications, Active Implementations, and Recently Archived specs. Constitution documents from external worktrees must be ignored.
- **RESTRICTION:** For now, the implementation will prioritize stability of the main root scan.

## 3. Acceptance Criteria (BDD)

### [REQ-1] Worktree Discovery
**User Story:** AS A developer, I WANT the console to identify all active git worktrees, SO THAT I can see specifications from other branches.

**Scenarios:**
1. **[Happy Path]** GIVEN a project with multiple git worktrees WHEN the console starts THEN it must list all specifications from all detected worktrees.
2. **[Edge Case]** GIVEN a project where a worktree has been deleted from disk but not from git metadata WHEN the console scans THEN it must gracefully ignore the missing directory and continue.
3. **[Ambiguity]** It is assumed that worktrees are located on the same local filesystem and accessible by the user.

### [REQ-2] Unified Specification View
**User Story:** AS A developer, I WANT a unified list of specifications, SO THAT I don't have to manually switch contexts to check progress.

**Scenarios:**
1. **[Happy Path]** GIVEN specifications with the same slug in different worktrees WHEN displayed in the console THEN they must be distinguished by their worktree/branch name.
2. **[Edge Case]** GIVEN a worktree that is not initialized with Specforce WHEN scanned THEN it must be ignored without error.

## 4. Business Invariants
- The single source of truth for the *active* worktree's state remains its local `.specforce` directory.
- Specifications from non-current worktrees must be clearly labeled as external to prevent confusion.
- Discovery failure in one worktree must not block the display of specifications from other worktrees.
