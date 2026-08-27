---
slug: 20260827-1515-living-module-specs-refinement
lens: Integration
---

# Feature: Standardize Module Living Specifications as Behavioral BDD Specs

## 1. Context & Value
Specforce module documentation (`.specforce/docs/modules/*.md`) represents the canonical living truth of domain subsystems. Currently, module synthesis during feature archival tends to capture low-level implementation details and internal code references, causing maintenance overhead and documentation drift. This feature standardizes module specifications as high-density Behavioral Living Specifications centered on domain boundaries, business invariants, accumulated BDD use cases, and public integration contracts.

## 2. Out of Scope (Anti-Goals)
- Do not modify the CLI command syntax for `specforce archive` or `specforce spec archive`.
- Do not build automated code parsing or static analysis to generate module files automatically.
- Do not alter the structure of temporary feature specifications (`requirements.md`, `design.md`, `tasks.md`).

## 3. Acceptance Criteria (BDD)

### [US-1] Behavioral Living Specification Guidelines in Archival Lifecycle
**User Story:** AS A project architect, I WANT the lifecycle archival workflow to enforce behavioral BDD synthesis in module documentation, SO THAT module docs remain an enduring behavioral source of truth without rotting internal code references.

**Scenarios:**
1. **[Happy Path]** GIVEN an agent completing the feature archival lifecycle WHEN synthesizing domain documentation into `.specforce/docs/modules/<domain>.md` THEN it documents business invariants, BDD user scenarios (GIVEN/WHEN/THEN for happy and failure paths), and public integration surfaces, strictly omitting internal file paths, internal struct names, and private helpers.
2. **[Edge Case]** GIVEN a feature that modifies existing domain behavior WHEN merging into an existing module living spec THEN updated BDD scenarios replace outdated behavior while untouched historic BDD scenarios are strictly preserved.

**Technical Constraints (NFR):**
- **[Performance]:** Instruction parsing and dynamic kit assembly completes with zero noticeable latency overhead (< 10ms).
- **[Safety & Security]:** Archival instructions prevent destructive overwrite of unrelated module living spec sections.
- **[Integrity]:** Template structure strictly matches the 5 standard canonical sections.

### [US-2] Canonical Module Living Spec Blueprint Alignment
**User Story:** AS A developer creating or updating a domain module, I WANT the module blueprint template to clearly guide me towards domain rules, BDD use cases, and public integration surfaces, SO THAT all modules across the project maintain structural and behavioral consistency.

**Scenarios:**
1. **[Happy Path]** GIVEN the constitution module artifact blueprint WHEN rendered for domain documentation THEN it provides sections for Domain Scope, Business Rules & Invariants (`[BR-xx]`), Canonical Use Cases with BDD scenarios (`[US-xx]`), Public Integration Surfaces, and Operational Invariants.
2. **[Edge Case]** GIVEN a developer inspecting the technical contracts section WHEN defining module interfaces THEN the template guides defining public CLI commands, API contracts, events, and cross-domain dependencies instead of internal implementation files.

**Technical Constraints (NFR):**
- **[Performance]:** Artifact blueprint generation is instantaneous (< 5ms).
- **[Maintainability]:** Clean YAML frontmatter and template instructions aligned with Specforce core conventions.
- **[Integrity]:** Consistent section numbering and tag formats across all module templates.

### [US-3] Canonical Sanitization of Existing Project Module Documents
**User Story:** AS A developer consulting existing module documentation in `.specforce/docs/modules/`, I WANT existing module files to reflect the pure behavioral living spec format, SO THAT they serve as pristine reference examples for future specifications and AI agents.

**Scenarios:**
1. **[Happy Path]** GIVEN existing module files in `.specforce/docs/modules/` containing internal code references WHEN sanitized THEN internal code file paths and package trees are replaced with public CLI interfaces and architectural integration contracts.
2. **[Edge Case]** GIVEN a module file with historical use cases WHEN sanitizing technical sections THEN all existing business invariants (`[BR-...]`) and BDD user stories (`[US-...]`) are preserved without loss of domain rules.

**Technical Constraints (NFR):**
- **[Performance]:** Clean file writes with standard formatting.
- **[Integrity]:** Zero loss of domain invariants or user stories during cleanup.

### [US-4] Opportunistic Legacy Module Migration during Reconciliation
**User Story:** AS A project architect, I WANT the lifecycle archival workflow to opportunistically upgrade legacy module files to the current 5-part canonical standard during reconciliation, SO THAT technical debt in older module files is progressively eliminated during feature merges.

**Scenarios:**
1. **[Happy Path]** GIVEN an agent consolidating domain documentation into an existing `.specforce/docs/modules/<domain>.md` that uses a legacy structure or contains internal code paths WHEN performing reconciliation THEN it opportunistically reformats the document into the 5 canonical sections, removing legacy code dumps while preserving all accumulated domain rules and use cases.
2. **[Edge Case]** GIVEN an existing module file with malformed headers or missing canonical sections WHEN migrating to current standard THEN the agent restructures the file to include all 5 canonical sections with valid headers.

**Technical Constraints (NFR):**
- **[Performance]:** In-memory markdown migration with zero added latency.
- **[Integrity]:** Complete preservation of existing domain invariants and historical BDD scenarios.

## 4. Business Invariants
- `[BR-MOD-01]` Module documentation files (`.specforce/docs/modules/*.md`) MUST represent behavioral living specifications (Domain Scope, Invariants, BDD Scenarios, Public Integration Surfaces, Operational Invariants) and MUST NOT contain internal source code file paths, internal package listings, or private implementation artifacts.
- `[BR-MOD-02]` Public Integration Surfaces in module documents are strictly limited to external consumer boundaries: public CLI commands, public API endpoints, emitted/consumed events, and cross-module contracts.
- `[BR-MOD-03]` Non-destructive synthesis: Archiving a feature MUST accumulate new BDD scenarios and update modified invariants without deleting untouched domain history.
- `[BR-MOD-04]` Opportunistic Migration: If a pre-existing module document does not follow the 5 canonical sections or contains legacy implementation code references, the archival process MUST upgrade and reformat it to the current canonical standard during reconciliation.

## 5. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Template loading and instruction generation executed instantaneously (< 20ms).
- **[Reliability]:** All embedded instructions and blueprints pass unit and integration test suites without regressions.
- **[Maintainability]:** High density, zero fluff, and zero conversational placeholders across all generated documents.
