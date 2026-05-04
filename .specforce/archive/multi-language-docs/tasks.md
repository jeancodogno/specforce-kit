---
slug: multi-language-docs
lens: Migration
---

# Implementation Roadmap: Multi-language Documentation Expansion

## 1. Execution Strategy
- **Gravity Order:** Folder Restructuring -> English Link Updates -> Root Document Localization (PT/ES) -> Folder Documentation Localization (PT/ES) -> Verification.

## 2. Tasks

### Phase 1: Folder Restructuring & English Baseline

#### T1.1: [SCAFFOLD] Create Language Directories
**State:** `[FINISHED]`
**Target:** `Global Scope`
**Context:** [REQ-1]

**Action Steps:**
- Create directories: `docs/en`, `docs/pt`, `docs/es`.

**Verification (TDD):**
`ls -d docs/en docs/pt docs/es`

#### T1.2: [MIGRATE] Move English Documentation
**State:** `[FINISHED]`
**Target:** `docs/en/`
**Context:** [REQ-1]

**Action Steps:**
- Move all `.md` files currently in `docs/` to `docs/en/`.

**Verification (TDD):**
`ls docs/en` (should show artifacts.md, cli.md, etc.) and `ls docs/*.md` (should be empty).

#### T1.3: [DOCS] Update English Root Links
**State:** `[FINISHED]`
**Target:** `README.md`
**Context:** [REQ-1]

**Action Steps:**
- Update all links pointing to `docs/*.md` to `docs/en/*.md`.
- Add language selector at the top: `🌎 **Idiomas / Languages:** [English](README.md) | [Português](README.pt.md) | [Español](README.es.md)`.

**Verification (TDD):**
Manual link check in README.md.

#### T1.4: [DOCS] Update English Contributing Links
**State:** `[FINISHED]`
**Target:** `CONTRIBUTING.md`
**Context:** [REQ-1]

**Action Steps:**
- Update any links pointing to `docs/*.md` to `docs/en/*.md`.

**Verification (TDD):**
Manual link check in CONTRIBUTING.md.

### Phase 2: Root Localization (PT/ES)

#### T2.1: [DOCS] Localize README (Português)
**State:** `[FINISHED]`
**Target:** `README.pt.md`
**Context:** [REQ-2]

**Action Steps:**
- Translate `README.md` to Portuguese.
- Ensure all links point to `docs/pt/*.md`.

**Verification (TDD):**
Manual review of `README.pt.md`.

#### T2.2: [DOCS] Localize README (Español)
**State:** `[FINISHED]`
**Target:** `README.es.md`
**Context:** [REQ-2]

**Action Steps:**
- Translate `README.md` to Spanish.
- Ensure all links point to `docs/es/*.md`.

**Verification (TDD):**
Manual review of `README.es.md`.

#### T2.3: [DOCS] Localize CONTRIBUTING (PT & ES)
**State:** `[FINISHED]`
**Target:** `CONTRIBUTING.{pt,es}.md`
**Context:** [REQ-2]

**Action Steps:**
- Translate `CONTRIBUTING.md` to Portuguese and Spanish.
- Ensure links point to corresponding language folders.

**Verification (TDD):**
Manual review.

### Phase 3: Folder Documentation Localization (PT/ES)

#### T3.1: [DOCS] Localize docs/pt/
**State:** `[FINISHED]`
**Target:** `docs/pt/`
**Context:** [REQ-3]

**Action Steps:**
- Translate all files in `docs/en/` to `docs/pt/`.
- Ensure all internal links within `docs/pt/` stay within the `pt/` folder.

**Verification (TDD):**
Manual link verification across all files in `docs/pt/`.

#### T3.2: [DOCS] Localize docs/es/
**State:** `[FINISHED]`
**Target:** `docs/es/`
**Context:** [REQ-3]

**Action Steps:**
- Translate all files in `docs/en/` to `docs/es/`.
- Ensure all internal links within `docs/es/` stay within the `es/` folder.

**Verification (TDD):**
Manual link verification across all files in `docs/es/`.

## 3. Pre-emptive Mitigations
- **Risk:** Link rot during migration -> **Mitigation:** Use a recursive grep to find all markdown links and verify their relative path validity after moves.
- **Risk:** Inconsistent translations -> **Mitigation:** Maintain a glossary of Specforce-specific terms (e.g., "Constitution", "Spec", "Implementation") to ensure they are translated consistently across all files.
