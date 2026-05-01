---
slug: worktree-documentation
lens: Balanced full-stack
---

# Feature: Git Worktree Documentation

## 1. Context & Value
Specforce already supports cross-worktree specification discovery, allowing developers to see active and recently archived specs across multiple git branches simultaneously in the TUI console. The business value of this feature is to provide a comprehensive guide on how to configure and utilize this capability, ensuring users understand how to manage parallel specifications safely and efficiently.

## 2. Out of Scope (Anti-Goals)
- Do not implement any new code features or modifications to the existing `git-worktree-support` engine.
- Do not change how Specforce interacts with git metadata or the `.specforce` directories.
- Do not refactor the TUI console.

## 3. Acceptance Criteria (BDD)

### [REQ-1] Create Worktree Documentation Guide
**User Story:** AS A developer, I WANT a dedicated documentation page explaining git worktree support, SO THAT I can understand its capabilities and limitations.

**Scenarios:**
1. **[Happy Path]** GIVEN a user opens the documentation WHEN they navigate to `docs/git-worktrees.md` (or similar) THEN they must see explanations of unified discovery, read-only constraints for external specs, and what is scanned (Active Specs, Implementations, Archived).
2. **[Happy Path]** GIVEN the documentation guide is created WHEN the user reads the constraints THEN it must clearly state that cross-worktree task modification is blocked and external Constitution documents are ignored.

### [REQ-2] Update README References
**User Story:** AS A developer evaluating Specforce, I WANT to see worktree support mentioned in the main documentation, SO THAT I know this advanced workflow is supported.

**Scenarios:**
1. **[Happy Path]** GIVEN a user reading the `README.md` WHEN they look at the "Documentation" or "Features" sections THEN they must find a reference or link to the new Git Worktree guide.

## 4. Business Invariants
- All documentation additions must align with the current technical implementation of `git-worktree-support` (i.e., cross-worktree specs are strictly read-only, and constitution artifacts from other worktrees are ignored).
- The instructions must remain purely informative and not prescribe non-existent CLI commands for worktree management.