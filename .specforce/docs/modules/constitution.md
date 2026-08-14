# Module: Constitution & Living Specs

## 1. Domain Scope
Manages the project's foundational guidelines, core architectural documents (`.specforce/docs/`), and domain-specific living specifications (`.specforce/docs/modules/`).

## 2. Business Rules & Invariants
- `[BR-CONST-01]` All core constitution documents MUST reside in `.specforce/docs/` (`principles.md`, `architecture.md`, `ui-ux.md`, `security.md`, `engineering.md`, `governance.md`).
- `[BR-CONST-02]` Domain-specific module living specifications MUST reside in `.specforce/docs/modules/<slug>.md`.
- `[BR-CONST-03]` The Memorial subsystem is decommissioned. Historical lessons and as-built business rules MUST be consolidated directly into the respective domain module living spec during archival.
- `[BR-CONST-04]` Modules MUST NOT be dumped in bulk into global prompt context; they must be lazy-loaded on demand when domain affinity is detected.

## 3. Canonical Requirements & Use Cases
### [US-CONST-01] Module Manifest Discovery & Lazy Loading
- **Scenario:** Agent works on a feature within a specific domain
  - **GIVEN** a module specification exists at `.specforce/docs/modules/<domain>.md`
  - **WHEN** the agent identifies affinity via working directory, feature slug, or user prompt
  - **THEN** the agent loads and adheres to the module living spec without loading unneeded domain manifests.

### [US-CONST-02] As-Built Archival Consolidation
- **Scenario:** Feature lifecycle closure and reconciliation
  - **GIVEN** a completed feature specification
  - **WHEN** `/spf.archive` is triggered
  - **THEN** the agent compares original requirements with the actual implemented code/tests and non-destructively merges the updated rules and use cases into the corresponding domain module.

## 4. Technical Contracts & Integration Points
- **CLI Commands:**
  - `specforce constitution status [--json]`: Scans and outputs status of core artifacts (Total: 6) and discovered module slugs.
  - `specforce constitution artifact <slug> [--json]`: Returns metadata, instructions, and template for a specific artifact.
  - `specforce archive instructions`: Outputs the complete Archival Protocol with As-Built Living Spec reconciliation rules.
  - `specforce spec archive <slug>`: Finalizes and moves the specification to `.specforce/archive/`.

## 5. Operational Invariants
- Core Constitution artifacts count is fixed to 6.
- Template rendering must be deterministic with zero leftover template placeholders.
