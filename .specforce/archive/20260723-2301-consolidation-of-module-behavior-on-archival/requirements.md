---
slug: 20260723-2301-consolidation-of-module-behavior-on-archival
lens: Backend-heavy
---

# Feature: Consolidation of Module Behavior on Archival

## 1. Context & Value
During feature archival via `/spf.archive`, the system must consolidate the domain behavior and invariants from the completed feature's design and requirements into `.specforce/docs/modules/<slug>.md`. This aligns Specforce with canonical SDD practices (similar to OpenSpec), ensuring that `.specforce/docs/modules/` functions as the living single source of truth for all module behaviors across feature lifecycles.

## 2. Out of Scope (Anti-Goals)
- Do not modify Go CLI binary code in `src/internal/spec/archive.go` or `src/internal/cli/archive.go` in this feature; focus strictly on instruction-driven agent workflows in `archive.md` and `archive.yaml`.
- Do not create custom automated merge parsers or script hooks for module Markdown consolidation.
- Do not bypass user confirmation when updating or creating module files.

## 3. Acceptance Criteria (BDD)

### [US-1] Mandatory Module Consolidation Step in Archival Instructions
**User Story:** AS AN AI agent executing feature archival (`/spf.archive`), I WANT TO be explicitly instructed to identify and consolidate the domain behavior of affected modules, SO THAT the project's living module documentation reflects the post-archival truth.

**Scenarios:**
1. **[Happy Path]** GIVEN a completed feature spec and active module documentation in `.specforce/docs/modules/<module>.md`, WHEN `/spf.archive` executes instruction generation, THEN `archive.md` instructs the agent as a mandatory step to scan the feature's `requirements.md` and `design.md`, extract new domain invariants and behaviors, and merge them into `.specforce/docs/modules/<module>.md`.
2. **[Edge Case - Non-Existent Module]** GIVEN a feature affecting a domain that does not yet have a file in `.specforce/docs/modules/`, WHEN the archival protocol runs, THEN the instructions direct the agent to create `.specforce/docs/modules/<module>.md` with the new domain invariants upon user confirmation.

**Technical Constraints (NFR):**
- **[Performance]:** Instruction file retrieval via `specforce archive instructions` remains instantaneous (< 10ms).
- **[Safety & Security]:** File write/update operations for module files require explicit user approval.
- **[Integrity]:** Zero loss of pre-existing module documentation during consolidation updates.
- **[Observability]:** The final archival summary output must explicitly report module consolidation status.

### [US-2] Synchronization of Agent Command Wrappers
**User Story:** AS A developer invoking `/spf.archive` across different supported AI coding agents, I WANT `archive.yaml` and `archive.md` to be fully aligned, SO THAT all agent clients execute the updated mandatory module consolidation protocol seamlessly.

**Scenarios:**
1. **[Happy Path]** GIVEN `archive.yaml` defines the `spf.archive` command content, WHEN an agent executes `/spf.archive`, THEN the execution protocol enforces step-by-step compliance with the module behavior consolidation mandate in `archive.md`.
2. **[Edge Case - No Module Affinity]** GIVEN a feature that has no domain affinity with any existing or new module, WHEN archival occurs, THEN the agent explicitly logs "Module Consolidation: N/A (No module affinity identified)" in the final summary output.

**Technical Constraints (NFR):**
- **[Performance]:** Agent prompt overhead added by module consolidation instructions is minimized (< 300 words).
- **[Maintainability]:** Clean separation of concerns between CLI binary execution and agent kit instructions.

## 4. Business Invariants
- Archival MUST NOT proceed to `specforce spec archive <slug>` until module behavior consolidation is evaluated and executed or explicitly waived.
- Existing module invariants in `.specforce/docs/modules/` MUST NEVER be deleted or overwritten without explicit user approval.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Instruction generation latency < 50ms.
- **[Reliability]:** Fail-safe execution; if module update fails or is rejected, archival halts cleanly before directory relocation.
- **[Maintainability]:** 100% compliance with Specforce Instruction-Driven Agent Pattern and Engineering constitution.
