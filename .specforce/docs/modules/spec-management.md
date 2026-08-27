# Module: Specification Management & Sizing Lifecycle

## 1. Domain Scope
Governs the lifecycle of feature and bug specifications (`.specforce/specs/<slug>/`), including tiered complexity sizing (`small`, `medium`, `large`, `complex`), dynamic lifecycle resizing (`specforce spec resize`), artifact requirements, and task verification rigor.

## 2. Business Rules & Invariants
- `[BR-SPEC-01]` Every specification must have an explicit or defaulted size attribute (`small`, `medium`, `large`, `complex`). When unspecified during initialization, size defaults to `medium`.
- `[BR-SPEC-02]` Artifact matrix dynamically adapts to spec sizing:
  - `small`: requires only `tasks.md` (or `bug-tasks.md`).
  - `medium`: requires `requirements.md` and `tasks.md`.
  - `large` & `complex`: requires the full triad (`requirements.md`, `design.md`, `tasks.md`).
- `[BR-SPEC-03]` Dynamic Resizing (`specforce spec resize <slug> --size <size>`):
  - Promotion (`small` -> `large`): newly required upstream artifacts are marked missing (`exists: false`), blocking implementation until completed.
  - Demotion (`large` -> `small`): validation rules relax, but existing files on disk are preserved as reference context (non-destructive).
- `[BR-SPEC-04]` Task verification (`ValidateTasks`) and coherence auditing (`DeterministicCheck`) adjust strictness based on size:
  - `small`: allows single-step tasks and optional `**Context:** [US-X]` mapping. Missing `requirements.md` is not flagged as an error.
  - `medium`, `large`, `complex`: strictly enforces minimum 2 action steps, mandatory `**Context:** [US-X]`, and deterministic requirement-to-task coherence.

## 3. Canonical Requirements & Use Cases
### [US-SPEC-01] Tiered Sizing at Initialization
- **Scenario:** Initializing a specification with explicit size
  - **GIVEN** a new feature or bug slug
  - **WHEN** running `specforce spec init <slug> --size small`
  - **THEN** metadata records `size: small`, and `specforce spec status <slug>` evaluates progress based solely on `tasks.md`.

### [US-SPEC-02] Dynamic Specification Resizing
- **Scenario:** Scope changes during specification or implementation
  - **GIVEN** an active specification
  - **WHEN** executing `specforce spec resize <slug> --size large`
  - **THEN** metadata updates atomically, status recalculates required artifacts, and implementation is gated until required artifacts exist.

### [US-SPEC-03] Adaptive Task and Coherence Auditing
- **Scenario:** Verifying task density and requirement mappings
  - **GIVEN** a `small` specification with simple tasks
  - **WHEN** task validation and coherence audit run
  - **THEN** the validation suite passes without forcing multi-step decomposition or artificial requirement tags.

## 4. Public Integration Surfaces & Contracts
- **Public CLI Commands:**
  - `specforce spec init <slug> [--type <type>] [--size <size>] [--json]`: Initializes a new feature or bug specification directory with tiered sizing metadata.
  - `specforce spec resize <slug> --size <size> [--json]`: Dynamically resizes an active specification and recalculates required artifact gates.
  - `specforce spec status <slug> [--json]`: Evaluates completion, artifact presence, and task progress based on spec sizing.
  - `specforce spec list [--type <type>] [--status <status>] [--json]`: Lists active or archived specifications across the workspace.
  - `specforce spec archive <slug>`: Finalizes and transitions an active specification to `.specforce/archive/`.
- **Cross-Module Contracts & Dependencies:**
  - Emits specification status and task roadmap data consumed by `agent-kit` implementation workflows.
  - Interacts with `constitution` living spec rules to validate coherence during planning and archival.

## 5. Operational & Quality Invariants
- All specification CLI commands must return within 50ms.
- Metadata operations must never overwrite or delete user markdown artifacts on disk.
- Lifecycle state transitions (active -> archived) must be atomic on the filesystem.

