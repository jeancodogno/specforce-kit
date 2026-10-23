---
prefix: SPEC
---

# Module: Specification Management & Sizing Lifecycle

## 1. Domain Scope
Governs the lifecycle of feature and bug specifications (`.specforce/specs/<slug>/`), including tiered complexity sizing (`small`, `medium`, `large`, `complex`), dynamic lifecycle resizing (`specforce spec resize`), artifact requirements, task verification rigor, and the per-specification module delta (declaration, sealing against stale bases, and the archive gate that verifies the module reflects the delta).

## 2. Business Rules & Invariants
- `[BR-SPEC-01]` Every specification must have an explicit or defaulted size attribute (`small`, `medium`, `large`, `complex`). When unspecified during initialization, size defaults to `medium`.
- `[BR-SPEC-02]` Artifact matrix and readiness gates dynamically adapt to spec sizing across all commands:
  - `small`: requires only `tasks.md` (or `bug-tasks.md`). `specforce implementation status` evaluates readiness solely based on `tasks.md`, and omitted upstream dependencies do not mark downstream artifacts as blocked.
  - `medium`: requires `requirements.md` and `tasks.md`. Implementation readiness evaluates both documents, bypassing optional architectural design artifacts.
  - `large` & `complex`: requires the full triad (`requirements.md`, `design.md`, `tasks.md`).
- `[BR-SPEC-03]` Dynamic Resizing (`specforce spec resize <slug> --size <size>`):
  - Promotion (`small` -> `large`): newly required upstream artifacts are marked missing (`exists: false`), blocking implementation until completed.
  - Demotion (`large` -> `small`): validation rules relax, but existing files on disk are preserved as reference context (non-destructive).
- `[BR-SPEC-04]` Task verification (`ValidateTasks`) and coherence auditing (`DeterministicCheck`) adjust strictness based on size:
  - `small`: allows single-step tasks and optional `**Context:** [US-X]` mapping. Missing `requirements.md` is not flagged as an error.
  - `medium`, `large`, `complex`: strictly enforces minimum 2 action steps, mandatory `**Context:** [US-X]`, and deterministic requirement-to-task coherence.
- `[BR-SPEC-05]` Every `medium`, `large` or `complex` specification declares its module impact before implementation: either one delta document per affected module (ADDED, MODIFIED and REMOVED sections carrying complete blocks, with provisional `NEW:` keys for additions) or the explicit declaration `delta.impact: none`. `small` specifications are exempt and keep the free-form archive path.
- `[BR-SPEC-06]` Sealing a delta (`specforce spec delta seal <slug>`) records a fingerprint of the current text of every module block that the delta modifies or removes, so that parallel specifications changing the same block are detected: a block matching the recorded fingerprint is pending, a block matching the delta text is applied, and any other text is a conflict that requires redoing the delta against the current module. A modified use case may not lose scenarios, and unknown delta sections are reported as problems.
- `[BR-SPEC-07]` `specforce spec archive <slug>` refuses to archive a `medium`, `large` or `complex` specification unless its delta status is `applied` or `none`; `small` specifications report `exempt`. The override flag for pending tasks never bypasses this check, and a refusal leaves the filesystem untouched.
- `[BR-SPEC-08]` Final identifiers for added blocks are allocated at archive time as the next free number per kind (rules and use cases numbered independently) and recorded in the specification metadata as a key-to-identifier mapping; identifiers of untouched or modified blocks are never renumbered and removed identifiers are never reused.

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

### [US-SPEC-04] Size-Aware Implementation Readiness and Artifact Dependency Evaluation
- **Scenario:** Implementation readiness query evaluates missing artifacts strictly against the active sizing tier
  - **GIVEN** a `small` specification with a valid `tasks.md`
  - **WHEN** executing `specforce implementation status <slug> --json`
  - **THEN** status evaluates to `ready` and `missing_artifacts` is empty, never blocking on unrequired triad artifacts.
- **Scenario:** Spec status dependency check ignores optional upstream artifacts omitted by size tier
  - **GIVEN** a `small` or `medium` specification where `design.md` is not required
  - **WHEN** executing `specforce spec status <slug> --json`
  - **THEN** `tasks.md` is not marked as blocked by missing `design.md`.

### [US-SPEC-05] Delta Declaration and Sealing During Planning
- **Scenario:** Planning a specification that changes module behavior
  - **GIVEN** a `medium` or larger specification that adds, modifies or removes module rules
  - **WHEN** planning writes the delta and runs `specforce spec delta seal <slug>`
  - **THEN** metadata records the base fingerprints and `specforce spec delta check <slug>` reports every item as pending.
- **Scenario:** Planning a specification with no module impact
  - **GIVEN** a `medium` specification that changes no module behavior
  - **WHEN** it declares `delta.impact: none`
  - **THEN** `specforce spec delta check <slug>` reports `none` and sealing is not required.
- **Scenario:** Another specification changed the same block after sealing
  - **GIVEN** a sealed delta whose target block was changed by a different specification
  - **WHEN** `specforce spec delta check <slug>` runs
  - **THEN** the item is reported as a conflict with the instruction to redo the delta and seal again.

### [US-SPEC-06] Archive Gate on Delta State
- **Scenario:** Archive refused while the module does not reflect the delta
  - **GIVEN** a `medium` or larger specification whose delta status is pending, missing or invalid
  - **WHEN** archiving is requested, with or without the override for pending tasks
  - **THEN** archiving is refused, every problem and pending item is listed with the next command to run, and the specification stays active.
- **Scenario:** Archive allowed when the module reflects the delta
  - **GIVEN** a specification whose delta status is applied, none or exempt
  - **WHEN** archiving is requested
  - **THEN** the specification moves to the archive.

## 4. Public Integration Surfaces & Contracts
- **Public CLI Commands:**
  - `specforce spec init <slug> [--type <type>] [--size <size>] [--json]`: Initializes a new feature or bug specification directory with tiered sizing metadata.
  - `specforce spec resize <slug> --size <size> [--json]`: Dynamically resizes an active specification and recalculates required artifact gates.
  - `specforce spec status <slug> [--json]`: Evaluates completion, artifact presence, and task progress based on spec sizing.
  - `specforce spec list [--type <type>] [--status <status>] [--json]`: Lists active or archived specifications across the workspace.
  - `specforce implementation status <slug> [--json]`: Evaluates implementation readiness and missing artifacts strictly matching the active spec size.
  - `specforce spec archive <slug> [--force]`: Finalizes and transitions an active specification to `.specforce/archive/`. For `medium`, `large` and `complex` specifications it requires a delta status of `applied` or `none`; `--force` only overrides the pending-tasks check, never the delta check.
  - `specforce spec delta check <slug> [--json]`: Read-only inspection of the module delta. Reports the overall status (`none`, `exempt`, `missing`, `pending`, `applied`, `invalid`), every item with its state, every problem at once, and the next free identifiers per module. Fails only for `invalid`.
  - `specforce spec delta seal <slug> [--json]`: Records the fingerprints of the module blocks that the delta modifies or removes in the specification metadata; refuses invalid deltas and writes nothing else.
- **Cross-Module Contracts & Dependencies:**
  - Emits specification status and task roadmap data consumed by `agent-kit` implementation workflows.
  - Interacts with `constitution` living spec rules to validate coherence during planning and archival.

## 5. Operational & Quality Invariants
- All specification CLI commands must return within 50ms.
- Metadata operations must never overwrite or delete user markdown artifacts on disk.
- Lifecycle state transitions (active -> archived) must be atomic on the filesystem.
- Delta commands and the archive gate must complete within 50ms and never modify module documents; `delta check` writes nothing and `delta seal` writes only the specification metadata (owner-only permissions).
- Module names in deltas are restricted to lowercase hyphenated names confined to the project's module collection; symbolic links for module or delta files are rejected.

