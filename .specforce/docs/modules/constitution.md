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

## 4. Public Integration Surfaces & Contracts
- **Public CLI Commands:**
  - `specforce constitution status [--json]`: Scans and outputs the initialization status of core artifacts (6 standard constitution docs) and discovered module living specs.
  - `specforce constitution artifact <slug> [--json]`: Returns blueprint metadata, instructions, and markdown templates for specific core constitution documents or module living specs.
  - `specforce archive instructions`: Emits the sovereign archival protocol and living spec reconciliation directives.
- **Cross-Module Contracts & Dependencies:**
  - Provides canonical living spec templates and reconciliation rules for `agent-kit` archival workflows.
  - Supplies architectural constraints and domain business rules to active spec planning in `spec-management`.

## 5. Operational & Quality Invariants
- Core Constitution artifacts count is fixed to 6 standard documents (`principles.md`, `architecture.md`, `ui-ux.md`, `security.md`, `engineering.md`, `governance.md`).
- Template rendering must be deterministic with zero leftover template placeholders.
- Module living specs must adhere strictly to the 5 canonical sections without internal code dumps.

