---
slug: 20260506-1539-shields-and-navigation-badges
lens: UI-heavy
---

# Technical Design: Shields.io and Multi-Language Navigation

## 1. Architecture Blueprint

The documentation navigation system follows a decentralized, static linking model. Every markdown file acts as a node with pre-defined relative links to its counterparts in other languages.

```mermaid
graph TD
    subgraph "Root Level"
        R[README.md] <--> RP[README.pt.md]
        R <--> RE[README.es.md]
        RP <--> RE
    end

    subgraph "Docs Level"
        DE[docs/en/*.md] <--> DP[docs/pt/*.md]
        DE <--> DS[docs/es/*.md]
        DP <--> DS
    end

    R -- "Metric Badges" --> S[Shields.io API]
    RP -- "Status Badges" --> S
    RE -- "Status Badges" --> S
```

## 2. File & Component Inventory

This feature involves updating 21 files (3 READMEs + 18 documentation files).

### Documentation Navigation Pattern
Every file in `docs/{lang}/*.md` will receive the following header (example for `docs/en/cli.md`):
```markdown
# CLI Reference
[ [English](cli.md) | [Português](../pt/cli.md) | [Español](../es/cli.md) ]
```

### README Badge Layout
The top section of `README.md`, `README.pt.md`, and `README.es.md` will be standardized:

**Metrics Block (Primary README only):**
```markdown
<p align="center">
  <a href="https://github.com/jeancodogno/specforce-kit/releases"><img src="https://img.shields.io/github/v/release/jeancodogno/specforce-kit?style=flat-square" alt="Latest Release"></a>
  <a href="https://github.com/jeancodogno/specforce-kit/actions/workflows/ci.yml"><img src="https://img.shields.io/github/actions/workflow/status/jeancodogno/specforce-kit/ci.yml?branch=main&style=flat-square" alt="CI Status"></a>
  <a href="https://goreportcard.com/report/github.com/jeancodogno/specforce-kit"><img src="https://goreportcard.com/badge/github.com/jeancodogno/specforce-kit?style=flat-square" alt="Go Report Card"></a>
  <a href="https://www.npmjs.com/package/@jeancodogno/specforce-kit"><img src="https://img.shields.io/npm/v/@jeancodogno/specforce-kit?style=flat-square" alt="NPM Version"></a>
  <a href="go.mod"><img src="https://img.shields.io/github/go-mod/go-version/jeancodogno/specforce-kit?style=flat-square" alt="Go Version"></a>
  <a href="https://github.com/jeancodogno/specforce-kit/issues"><img src="https://img.shields.io/github/issues/jeancodogno/specforce-kit?style=flat-square" alt="GitHub issues"></a>
  <a href="https://github.com/jeancodogno/specforce-kit/pulls"><img src="https://img.shields.io/github/issues-pr/jeancodogno/specforce-kit?style=flat-square" alt="GitHub pull requests"></a>
  <a href="LICENSE"><img src="https://img.shields.io/github/license/jeancodogno/specforce-kit?style=flat-square" alt="License"></a>
</p>
```

**Language Switcher Block (All READMEs):**
```markdown
<p align="center">
  <a href="README.md"><img src="https://img.shields.io/badge/lang-en-red.svg?style=flat-square" alt="English"></a>
  <a href="README.pt.md"><img src="https://img.shields.io/badge/lang-pt--br-green.svg?style=flat-square" alt="Português"></a>
  <a href="README.es.md"><img src="https://img.shields.io/badge/lang-es-yellow.svg?style=flat-square" alt="Español"></a>
</p>
```

### Modified Files List:
- `README.md`: Add Metrics and Language Switcher badges.
- `README.pt.md`: Replace text navigation with Language Switcher badges. Add Metrics badges.
- `README.es.md`: Replace text navigation with Language Switcher badges. Add Metrics badges.
- `docs/en/*.md` (6 files): Add relative language navigation header.
- `docs/pt/*.md` (6 files): Add relative language navigation header.
- `docs/es/*.md` (6 files): Add relative language navigation header.

## 3. Interaction Constraints
- **Relative Linking:** All links within the `docs/` folder must use relative paths (e.g., `../pt/filename.md`) to ensure they work correctly in both GitHub UI and local clones.
- **Badge Consistency:** All Shields.io badges MUST use `?style=flat-square` to match the project's modern aesthetic.
- **Alt Text:** Every badge must have descriptive alt text for accessibility.
