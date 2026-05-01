---
slug: git-worktrees
lens: Balanced full-stack
---

# Git Worktree Support

## Introduction

Spec-Driven Development often involves working on multiple features or bug fixes in parallel across different git branches. Specforce simplifies this workflow by providing native support for **Git Worktrees**. This feature allows the Specforce TUI console to discover and display specifications and progress from linked worktrees, enabling a unified view of the entire project's status without the need to constantly switch branches.

## Discovery Engine

Specforce includes a Unified Discovery Engine that automatically detects linked git worktrees associated with the current repository. This engine scans the `.git/worktrees` directory to identify other active worktrees and their respective paths.

The primary benefit of this integration is the ability to see specifications from multiple branches (e.g., `main`, `feature/X`, `bugfix/Y`) simultaneously within the Specforce TUI console. This provides developers with a global view of all ongoing work across the project, facilitating better coordination and preventing duplicate efforts.

## Constraints

To ensure data integrity and prevent accidental cross-branch modifications, Specforce enforces the following constraints:

- **Read-Only Access:** All specifications discovered from external git worktrees (those other than your current working directory) are strictly **Read-Only**.
- **Blocked Actions:** For these external specifications, all state-modifying actions are blocked. This includes claiming or finishing tasks, transitioning specification states, and any automated file edits.
- **Constitution Isolation:** Specforce ignores the `.specforce/docs/` directory (the Project Constitution) from external worktrees. This ensures that your current implementation always adheres to the architectural and engineering principles defined in your local branch.

## Scanning Scope

The Unified Discovery Engine performs a deep scan of linked worktrees to gather all relevant SDD (Spec-Driven Development) artifacts. The following items are included in the cross-worktree scan:

- **Active Specs:** All files under `.specforce/specs/**` across all discovered worktrees.
- **Implementations:** Any implementation files or roadmaps currently referenced by active specifications.
- **Archived Specs:** Previously completed and archived specifications located in `.specforce/archive/**`.

By including both active and archived work, Specforce provides a historical and current context of the project's evolution across all branches.
