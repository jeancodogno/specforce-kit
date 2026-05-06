---
slug: 20260506-1539-shields-and-navigation-badges
lens: UI-heavy
---

# Implementation Roadmap: Shields.io and Multi-Language Navigation

## 1. Execution Strategy
- **Gravity Order:** Update root READMEs (Foundation) -> Update English Docs (Template) -> Update Portuguese Docs -> Update Spanish Docs.
- **Verification:** Every file update will be verified by checking the visual structure and relative link validity.

## 2. Tasks

### Phase 1: Root READMEs Standardization

- [x] T1.1: [DOCS] Update `README.md` with Shields.io badges
**Target:** `README.md`
**Context:** [REQ-1], [REQ-2]

**Action Steps:**
- Insert the Metrics Block (Release, CI, NPM, Go, License) using `flat-square` style below the logo.
- Insert the Language Switcher Block (EN, PT, ES) using colored badges.
- Ensure all links are HTTPS or correct relative paths.

**Verification (TDD):**
- Inspect file to confirm all 11 badges (8 metrics + 3 language switchers) are present and correctly formatted.
- Verify that `README.pt.md` and `README.es.md` links are valid.

- [x] T1.2: [DOCS] Update `README.pt.md` with Shields.io badges
**Target:** `README.pt.md`
**Context:** [REQ-1], [REQ-2]

**Action Steps:**
- Replace existing text-based language navigation with the standardized Language Switcher badges.
- Add the standardized Metrics Block.

**Verification (TDD):**
- Inspect file to confirm visual parity with `README.md` (translated content but same badge structure).

- [x] T1.3: [DOCS] Update `README.es.md` with Shields.io badges
**Target:** `README.es.md`
**Context:** [REQ-1], [REQ-2]

**Action Steps:**
- Replace existing text-based language navigation with the standardized Language Switcher badges.
- Add the standardized Metrics Block.

**Verification (TDD):**
- Inspect file to confirm visual parity with `README.md`.

### Phase 2: English Documentation Navigation

- [x] T2.1: [DOCS] Add navigation header to `docs/en/*.md`
**Target:** `docs/en/*.md`
**Context:** [REQ-3]

**Action Steps:**
- Prepend `[ [English](filename.md) | [Português](../pt/filename.md) | [Español](../es/filename.md) ]` to all 6 files.
- Files: `artifacts.md`, `cli.md`, `configuration.md`, `getting-started.md`, `git-worktrees.md`, `supported-tools.md`.

**Verification (TDD):**
- Confirm the header is at the very top of each file.
- Verify relative paths point correctly to `../pt/` and `../es/`.

### Phase 3: Portuguese Documentation Navigation

- [x] T3.1: [DOCS] Add navigation header to `docs/pt/*.md`
**Target:** `docs/pt/*.md`
**Context:** [REQ-3]

**Action Steps:**
- Prepend `[ [English](../en/filename.md) | [Português](filename.md) | [Español](../es/filename.md) ]` to all 6 files.

**Verification (TDD):**
- Confirm the header is at the very top of each file.

### Phase 4: Spanish Documentation Navigation

- [x] T4.1: [DOCS] Add navigation header to `docs/es/*.md`
**Target:** `docs/es/*.md`
**Context:** [REQ-3]

**Action Steps:**
- Prepend `[ [English](../en/filename.md) | [Português](../pt/filename.md) | [Español](filename.md) ]` to all 6 files.

**Verification (TDD):**
- Confirm the header is at the very top of each file.

## 3. Pre-emptive Mitigations
- **Risk:** Broken relative links in subdirectories -> **Mitigation:** Manually verify a sample of links from each language directory using the GitHub "Preview" feature if available, or by verifying paths via `ls`.
- **Risk:** Inconsistent filenames across languages -> **Mitigation:** Before applying Phase 3 and 4, verify that filenames in `pt/` and `es/` exactly match those in `en/`.
