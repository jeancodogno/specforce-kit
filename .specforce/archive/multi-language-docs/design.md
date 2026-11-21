---
slug: multi-language-docs
lens: Migration
---

# Technical Design: Multi-language Documentation Expansion

## 1. Architecture Blueprint

```mermaid
graph TD
    Root[Project Root] --> R_EN[README.md - EN]
    Root --> R_PT[README.pt.md]
    Root --> R_ES[README.es.md]
    Root --> C_EN[CONTRIBUTING.md - EN]
    Root --> C_PT[CONTRIBUTING.pt.md]
    Root --> C_ES[CONTRIBUTING.es.md]
    Root --> DocsDir[docs/]
    
    DocsDir --> EN[docs/en/]
    DocsDir --> PT[docs/pt/]
    DocsDir --> ES[docs/es/]
    
    EN --> E_Files[artifacts.md, cli.md, configuration.md, ...]
    PT --> P_Files[artifacts.md, cli.md, configuration.md, ...]
    ES --> S_Files[artifacts.md, cli.md, configuration.md, ...]
    
    R_EN -- "links to" --> EN
    R_PT -- "links to" --> PT
    R_ES -- "links to" --> ES
```

## 2. File & Component Inventory

**Documentation Migration:**
- `README.md` -> [Add language selector at top, update relative links to `docs/en/`]
- `README.pt.md` -> [Translated version of README.md, links to `docs/pt/`]
- `README.es.md` -> [Translated version of README.md, links to `docs/es/`]
- `CONTRIBUTING.md` -> [Update relative links to `docs/en/`]
- `CONTRIBUTING.pt.md` -> [Translated version of CONTRIBUTING.md, links to `docs/pt/`]
- `CONTRIBUTING.es.md` -> [Translated version of CONTRIBUTING.md, links to `docs/es/`]

**Directory Restructuring:**
- `docs/en/` -> [Container for English docs]
- `docs/pt/` -> [Container for Portuguese docs]
- `docs/es/` -> [Container for Spanish docs]

**Existing Docs (Moving to `docs/en/` and duplicating/translating to `pt/` and `es/`):**
- `docs/artifacts.md` -> `docs/{en,pt,es}/artifacts.md`
- `docs/cli.md` -> `docs/{en,pt,es}/cli.md`
- `docs/configuration.md` -> `docs/{en,pt,es}/configuration.md`
- `docs/getting-started.md` -> `docs/{en,pt,es}/getting-started.md`
- `docs/git-worktrees.md` -> `docs/{en,pt,es}/git-worktrees.md`
- `docs/supported-tools.md` -> `docs/{en,pt,es}/supported-tools.md`

## 3. Link Update Logic
Every link within the `.md` files must be verified and adjusted:
- Links from `README.md` to `docs/*.md` must now point to `docs/en/*.md`.
- Links within `docs/en/*.md` to other files in the same directory remain the same but must be checked for correctness relative to the new `en/` nesting if they were using complex relative paths.
- Translated files in `docs/pt/` must have all their internal links point to other `docs/pt/` files.
