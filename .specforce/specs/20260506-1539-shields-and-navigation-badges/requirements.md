---
slug: 20260506-1539-shields-and-navigation-badges
lens: UI-heavy
---

# Feature: Shields.io and Multi-Language Navigation

## 1. Context & Value
Specforce Kit is a global framework for AI-Native Engineering. To reinforce its professionalism and accessibility, we need to integrate Shields.io badges for project metrics and establish a consistent, visual navigation system between different languages (English, Portuguese, Spanish) across all READMEs and technical documentation.

## 2. Out of Scope (Anti-Goals)
- Do not implement automatic translation of existing documentation.
- Do not create a dedicated documentation website (e.g., Docusaurus, MkDocs) in this spec.
- Do not add badges for metrics that are not currently tracked (e.g., code coverage, as it's not yet integrated into CI for all platforms).

## 3. Acceptance Criteria (BDD)

### [REQ-1] Project Metric Badges (Main README)
**User Story:** AS A developer visiting the repository, I WANT TO see the project's health and status at a glance, SO THAT I can trust the tool's maturity.

**Scenarios:**
1. **[Happy Path]** GIVEN the user is on `README.md` WHEN they look at the top section THEN they see badges for Release Version, CI Status, NPM Version, Go Version, Go Report Card, GitHub Issues, GitHub PRs, and MIT License.
2. **[Visual Consistency]** GIVEN the metric badges WHEN rendered THEN they all use the `flat-square` style for a modern and unified look.

**UI/UX Specifics:**
- **View/Component:** Badge header row below the project description.
- **Feedback Logic:** Links on badges point to relevant sources (Releases page, Actions log, NPM page, License file).

**Technical Constraints (NFR):**
- **[Performance]:** Badges must load asynchronously via Shields.io CDN (standard browser behavior).
- **[Safety & Security]:** Links must use HTTPS.

### [REQ-2] Language Switcher Navigation (READMEs)
**User Story:** AS A non-English speaking developer, I WANT TO quickly find the documentation in my preferred language, SO THAT I can start using the tool without language barriers.

**Scenarios:**
1. **[Happy Path]** GIVEN the user is on any of the README files WHEN they click a language badge THEN they are redirected to the corresponding README (e.g., clicking 'Português' badge on `README.md` goes to `README.pt.md`).
2. **[Edge Case]** GIVEN a new language is added in the future WHEN updating the switcher THEN the change must be applied consistently across all README versions to prevent dead ends.

**UI/UX Specifics:**
- **View/Component:** Language badge row (English, Português, Español) using colored badges (Red, Green, Yellow).
- **Feedback Logic:** Active language badge should ideally be first or highlighted.

**Technical Constraints (NFR):**
- **[Maintainability]:** Links between READMEs must be relative paths.
- **[Performance]:** Navigation must be instantaneous (standard Markdown link).

### [REQ-3] Technical Docs Navigation Header
**User Story:** AS A developer reading deep documentation, I WANT TO switch languages for the specific page I am reading, SO THAT I can compare technical terms or read in my native language.

**Scenarios:**
1. **[Happy Path]** GIVEN the user is reading `docs/en/cli.md` WHEN they use the navigation header THEN they can jump to `docs/pt/cli.md` or `docs/es/cli.md`.
2. **[Edge Case]** GIVEN a documentation file exists in English but not yet in Spanish WHEN the user clicks the Spanish link THEN it should point to the correct path (which may be a 404 or a "contribution wanted" if handled later, but the link structure must be correct).

**UI/UX Specifics:**
- **View/Component:** Simple text-based navigation line at the very top of every `.md` file in `docs/`.
- **Interaction Model:** `[ English | Português | Español ]` format.

**Technical Constraints (NFR):**
- **[Maintainability]:** Documentation links must be relative to the file location.
- **[Consistency]:** All 18 documentation files (6 per language) must have this header.

## 4. Business Invariants
- `README.md` is the source of truth for the project's primary entry point.
- All language versions of the documentation must remain structurally synchronized (same filenames across `en/`, `pt/`, `es/`).

## 5. Global UI/UX Contract (Documentation Design)
- **Density Posture:** Standard Markdown rendering.
- **Signature Moves:** Consistent use of `flat-square` badges.
- **Interaction Model:** Standard hyperlinking.
- **State Behavior:** Broken links are strictly forbidden; all relative paths must be verified.

## 6. Global Non-Functional Requirements (NFRs)
- **[Maintainability]:** Documentation structure must follow the existing pattern: `README.{lang}.md` and `docs/{lang}/{topic}.md`.
- **[Accessibility]:** All images and badges must have descriptive `alt` text.
- **[Reliability]:** Shields.io is the standard provider; no custom badge servers allowed.
