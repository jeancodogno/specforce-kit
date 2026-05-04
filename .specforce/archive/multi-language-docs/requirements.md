---
slug: multi-language-docs
lens: Migration
---

# Feature: Multi-language Documentation Expansion (EN, PT, ES)

## 1. Context & Value
Specforce Kit is used by a global community. Providing documentation in Portuguese and Spanish alongside English lowers the barrier to entry for non-English speakers and reinforces the framework's commitment to accessibility. This migration organizes the documentation into a standardized multi-language structure.

## 2. Out of Scope (Anti-Goals)
- Do not translate the "Project Constitution" (.specforce/docs/*.md) in this spec. It remains in English to ensure precision for AI agents.
- Do not add automated translation tools or CI/CD translation scripts.
- Do not modify the source code of the Specforce CLI/TUI.

## 3. Acceptance Criteria (BDD)

### [REQ-1] Multi-language Folder Structure
**User Story:** AS A maintainer, I WANT TO organize documentation by language folders, SO THAT the structure is scalable and standard.

**Scenarios:**
1. **[Happy Path]** GIVEN the current `docs/` folder WHEN the migration is executed THEN a new structure `docs/{en,pt,es}/` MUST be created, and all existing English docs MUST be moved to `docs/en/`.
2. **[Edge Case]** GIVEN an existing file in `docs/` WHEN it is moved THEN all internal relative links in that file MUST be updated to point to the correct relative path in the new structure.

### [REQ-2] Regional README & CONTRIBUTING Files
**User Story:** AS A visitor from a specific region, I WANT TO read the project overview in my native language, SO THAT I can understand the tool more easily.

**Scenarios:**
1. **[Happy Path]** GIVEN the root directory WHEN the migration is complete THEN `README.pt.md`, `README.es.md`, `CONTRIBUTING.pt.md`, and `CONTRIBUTING.es.md` MUST exist with fully translated content.
2. **[Happy Path]** GIVEN the main `README.md` WHEN a user visits THEN it MUST contain a language selector at the top linking to the regional versions.

### [REQ-3] Content Accuracy & Linking
**User Story:** AS A user, I WANT TO navigate between translated documents, SO THAT I don't get lost or end up in a different language unexpectedly.

**Scenarios:**
1. **[Happy Path]** GIVEN a translated document (e.g., `docs/pt/getting-started.md`) WHEN I click a link to another doc (e.g., `cli.md`) THEN it MUST lead to the translated version in the same language (`docs/pt/cli.md`).
2. **[Edge Case]** GIVEN a missing translation for a specific document WHEN a link is clicked THEN it SHOULD fallback to the English version (`docs/en/`) with a clear notice (Assumed fallback strategy).

## 4. Business Invariants
- The English version (`docs/en/` and root `README.md`) remains the "Source of Truth".
- Every translated document MUST contain a link back to the English version for comparison.
- Filenames MUST remain identical across language folders (e.g., `docs/en/cli.md` vs `docs/pt/cli.md`).
