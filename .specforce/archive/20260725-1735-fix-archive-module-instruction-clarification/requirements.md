# Requirements Specification: Fix Archive Module Instruction Clarification

## Problem Statement
During feature archival, when `specforce constitution status` returns an empty `modules: []` list, AI agents incorrectly assume that module consolidation is not applicable (`N/A`), skipping the creation of new domain module documents (`.specforce/docs/modules/<module>.md`) as mandated by Rule 4/Step 5 of `archive.md`.

## User Stories & Acceptance Criteria

### User Story 1
As a Specforce agent executing feature archival, I need explicit, deterministic instructions for module harvesting so that I never incorrectly mark module consolidation as `N/A` when an affected domain exists (even if no module files exist yet).

#### Acceptance Criteria
1. `archive.md` Step 5 MUST clearly separate module identification, checking existing modules, and requesting creation of new module documents.
2. Step 5 MUST explicitly clarify that an empty `modules` list in Constitution status DOES NOT mean module consolidation is `N/A`.
3. Step 5 MUST restrict `N/A` labeling ONLY to features that are strictly cross-cutting infrastructure or tooling without any domain business logic context.
