---
slug: worktree-documentation
lens: Balanced full-stack
---

# Technical Design: Git Worktree Documentation

## 1. Architecture Blueprint
*This diagram illustrates the relationship between the Specforce TUI and the documentation structure explaining the cross-worktree discovery engine.*

```mermaid
graph TD
    User((User)) --> README[README.md]
    README -- "Links to" --> WT_DOC[docs/git-worktrees.md]
    
    subgraph "Documentation Context"
        WT_DOC -- "Explains" --> Discovery[Unified Discovery Engine]
        WT_DOC -- "Explains" --> Permissions[Read-Only Constraints]
        WT_DOC -- "Explains" --> Scope[Scanning Scope]
    end

    subgraph "System Behavior (Documented)"
        Discovery -- "Scans" --> LocalWT[Local Worktree Specs]
        Discovery -- "Scans" --> ExternalWT[Linked Worktrees Specs]
        Permissions -- "Enforces" --> ReadOnly[External Specs: View Only]
        Scope -- "Includes" --> Active[Active Specs]
        Scope -- "Includes" --> Archived[Archived Specs]
    end
```

## 4. File & Component Inventory
*The exact files that the Developer must create or modify. Map the core responsibility.*

**Documentation:**
- `docs/git-worktrees.md` -> Core documentation file providing a comprehensive guide on Specforce's Git Worktree support, including technical constraints and discovery logic.
- `README.md` -> Main entry point update to include references and links to the Git Worktree documentation, improving feature discoverability.

**Note:** Only write documentation. DO NOT write code.
