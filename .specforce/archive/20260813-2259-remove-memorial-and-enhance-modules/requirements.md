---
slug: 20260813-2259-remove-memorial-and-enhance-modules
lens: Backend-heavy
---

# Feature: Remove Memorial System and Enhance Module Living Specs

## 1. Context & Value
The legacy distributed memorial system (`.specforce/memorial/`) introduced cognitive and operational overhead through transient memory fragments and distillation steps that duplicated project knowledge. This feature completely eliminates the memorial and memory subsystem across the Specforce CLI, Agent Kit, and Constitution, while elevating `.specforce/docs/modules/` into OpenSpec-style Canonical Living Specifications that capture business rules, use cases, and as-built implementations.

## 2. Out of Scope (Anti-Goals)
- Do not modify existing active spec schemas (`requirements.yaml`, `design.yaml`, `tasks.yaml`) outside of module and archive templates.
- Do not build automated git commit hooks for module syncing; reconciliation happens natively inside the `/spf.archive` agentic workflow.
- Do not remove the specification lifecycle archive command `specforce spec archive <slug>`.

## 3. Acceptance Criteria (BDD)

### [US-1] Complete Decommissioning of Memorial Subsystem
**User Story:** AS A developer and AI agent, I WANT TO use Specforce without any memorial or memory fragment commands, SO THAT all canonical knowledge resides solely in Constitution and Module documents.

**Scenarios:**
1. **[Happy Path - CLI Execution]** GIVEN a Specforce project WHEN the user executes `specforce archive` THEN only `instructions` is available as a sub-command, and executing `memorial` or `distill` returns an unknown command error.
2. **[Happy Path - Bootstrap & Init]** GIVEN `specforce init` is executed in a new project WHEN bootstrapping directory structure THEN `.specforce/memorial/` is NOT created and `MemorialService` is not initialized.
3. **[Happy Path - Constitution Status]** GIVEN `specforce constitution status --json` is executed WHEN querying artifacts THEN `memorial` (`.specforce/memorial/ROUTING.md`) is NOT listed as a constitution artifact.
4. **[Edge Case - Legacy Memorial Cleanup]** GIVEN a project containing an existing `.specforce/memorial/` directory or `memorial.yaml` artifact WHEN `specforce archive instructions` or status scans run THEN legacy memory references are safely ignored without crashes.

**Technical Constraints (NFR):**
- **[Performance]:** Zero disk I/O allocated to memory fragment scanning or consolidation.
- **[Safety & Security]:** Clean deletion of obsolete Go packages and CLI handlers without leaving dangling references or broken builds.
- **[Integrity]:** `AGENTS.md` and `README.md` must not contain any reference to `.specforce/memorial/`.

---

### [US-2] Canonical Living Specs Template for Modules (`module.yaml`)
**User Story:** AS A system architect and developer, I WANT module manifests in `.specforce/docs/modules/<slug>.md` to serve as Canonical Living Specs, SO THAT business rules, BDD use cases, and technical contracts are preserved across feature lifecycles.

**Scenarios:**
1. **[Happy Path - Template Generation]** GIVEN an agent initializing or updating a module manifest WHEN reading `src/internal/agent/artifacts/constitution/module.yaml` THEN the template provides sections for Domain Scope, Business Rules & Invariants, Canonical Requirements & Use Cases (BDD Happy Path & Edge Cases), Technical Contracts & Integration Points, and Operational Invariants.
2. **[Edge Case - Non-Domain Infrastructure]** GIVEN a feature that is purely cross-cutting developer tooling without business domain logic WHEN evaluating module affinity THEN the agent is permitted to mark module consolidation as N/A with explicit justification.

**Technical Constraints (NFR):**
- **[Performance]:** Module template must remain high-density Markdown formatted for rapid AI ingestion (< 200 lines per module).
- **[Maintainability]:** Follow standard Markdown heading structure (`# Module: {{MODULE_NAME}}`, `## 1. Business Rules & Invariants`, `## 2. Canonical Requirements & Use Cases`, `## 3. Technical Contracts & Integration Points`).

---

### [US-3] As-Built Post-Implementation Module Consolidation in Archival
**User Story:** AS AN AI agent executing `/spf.archive`, I WANT TO synthesize actual implemented code, tests, and user chat refinements into the domain module, SO THAT `.specforce/docs/modules/` always reflects the true as-built system behavior.

**Scenarios:**
1. **[Happy Path - As-Built Synthesis]** GIVEN a completed feature implementation with code changes and tests WHEN `/spf.archive` executes instruction Step 5 THEN the agent compares the original spec with actual code diffs and merges new canonical use cases and invariants into `.specforce/docs/modules/<domain>.md`.
2. **[Happy Path - Existing Module Update]** GIVEN an existing module document `.specforce/docs/modules/<domain>.md` WHEN archiving a feature touching that domain THEN new use cases and modified business rules are merged non-destructively without erasing prior domain use cases.
3. **[Edge Case - Divergence Between Spec and Code]** GIVEN a feature where requirements were altered interactively during `/spf.implement` without updating the spec's `requirements.md` WHEN archiving occurs THEN the agent extracts the actual implemented behavior and records the verified reality into the module documentation.

**Technical Constraints (NFR):**
- **[Performance]:** Lazy load only domain-relevant module files during planning and archival.
- **[Integrity]:** Archival handoff summary must explicitly report `**Module Behavior Consolidated:** [Yes/No/NA - List updated or created module files in .specforce/docs/modules/]`.

---

### [US-4] Context Conservation Terminology Alignment across Commands
**User Story:** AS AN AI agent reading agent workflow commands (`constitution.yaml`, `implement.yaml`, `spec.yaml`), I WANT clear token optimization instructions without confusing references to "Memory", SO THAT context reuse is unambiguous.

**Scenarios:**
1. **[Happy Path - Terminology Update]** GIVEN the command definitions for `constitution.yaml`, `implement.yaml`, and `spec.yaml` WHEN inspecting pre-flight instructions THEN clauses previously titled "Memory Check" are renamed to "Context Reuse Check" (referring to conversation history cache).

**Technical Constraints (NFR):**
- **[Observability]:** Unambiguous instructions preventing agents from searching for non-existent memory files.

## 4. Business Invariants
- **[INV-1]** All project domain invariants, canonical use cases, and business rules MUST reside in `.specforce/docs/modules/<slug>.md`.
- **[INV-2]** Global architectural, security, and engineering rules MUST reside in `.specforce/docs/` (`architecture.md`, `engineering.md`, `security.md`, `principles.md`, `governance.md`, `ui-ux.md`).
- **[INV-3]** No tool, command, or workflow shall attempt to read or write to `.specforce/memorial/`.
- **[INV-4]** Spec lifecycle closure MUST exclusively execute `specforce spec archive <slug>`.

## 5. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Zero CPU/IO overhead on memory operations; sub-50ms CLI execution time for `specforce archive instructions`.
- **[Reliability]:** 100% test pass rate across all Go unit tests after removing `memorial.go` and updating test suites.
- **[Security]:** Strict file path validation for all `.specforce/docs/modules/` paths via `core.SecurePath`.
- **[Maintainability]:** Clean codebase without dead code, deprecated flags, or orphaned template files.
