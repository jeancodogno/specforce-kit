---
slug: 20260825-1410-tiered-spec-sizing
lens: Balanced full-stack
---

# Feature: Tiered Spec Sizing & Dynamic Lifecycle Transitions

## 1. Context & Value
Developers and AI agents face unnecessary overhead when creating trivial features and bug fixes because the specification workflow enforces a heavy three-tier document pipeline for all changes. By introducing tiered sizing (`small`, `medium`, `large`, `complex`) and dynamic lifecycle transitions (`spec resize`), teams can match specification overhead to actual complexity while allowing seamless upgrade or downgrade when scopes change.

## 2. Out of Scope (Anti-Goals)
- Automatic source code modification or task execution based on spec size.
- Automatic deletion or purging of artifact files from disk during a spec downgrade.
- Altering the archival format or `.specforce/archive/` structure.

## 3. Acceptance Criteria (BDD)

### [US-1] Explicit Spec Sizing at Initialization
**User Story:** AS A Developer or AI Agent, I WANT TO specify a complexity size when initializing a specification, SO THAT the CLI configures the expected artifact matrix and validation rules appropriately.

**Scenarios:**
1. **[Happy Path]** GIVEN a new feature idea WHEN the developer initializes the specification with a valid size parameter (`small`, `medium`, `large`, `complex`) THEN the system records the size metadata and configures the required artifacts list accordingly.
2. **[Edge Case]** GIVEN an unsupported size value WHEN the developer attempts to initialize a specification THEN the system rejects the operation with a descriptive validation error listing the supported size options (`small`, `medium`, `large`, `complex`).

**UI/UX Specifics:**
- **View/Component:** Command line output and JSON response status.
- **Feedback Logic:** Clear confirmation displaying the assigned size category and next actionable steps.
- **Keybindings:** Standard CLI invocation flags.

**Technical Constraints (NFR):**
- **[Performance]:** Initialization completed in under 50ms.
- **[Safety & Security]:** Metadata written atomically with valid directory permissions (0750 / 0600).
- **[Integrity]:** Default size defaults to `medium` if omitted for backwards compatibility.
- **[Observability]:** Output status logged to terminal and structured JSON stream.

### [US-2] Dynamic Specification Resizing (Promotion and Demotion)
**User Story:** AS A Developer or AI Agent, I WANT TO resize an active specification when scope changes, SO THAT missing artifacts are required on promotion or relaxed on demotion without losing existing work.

**Scenarios:**
1. **[Happy Path - Promotion]** GIVEN an active specification with size `small` WHEN the scope expands and the specification is resized to `large` THEN the completion progress reflects the newly required artifacts (`requirements.md`, `design.md`) and blocks implementation until they are created.
2. **[Happy Path - Demotion]** GIVEN an active specification with size `large` containing generated artifacts WHEN the scope is reduced and the specification is resized to `small` THEN validation rules relax immediately, existing files are preserved on disk as reference context, and completion progress reaches 100% based on `tasks.md`.
3. **[Edge Case]** GIVEN a specification that does not exist WHEN a resize command is executed THEN the system returns an informative error without mutating any workspace files.

**UI/UX Specifics:**
- **View/Component:** Specforce Console status table and CLI progress cards.
- **Feedback Logic:** Immediate update of progress percentage and artifact checklist in the status view.
- **Keybindings:** Navigation in TUI status overview.

**Technical Constraints (NFR):**
- **[Performance]:** Metadata update and status recalculation executed in under 20ms.
- **[Safety & Security]:** Non-destructive operation; existing artifacts are never automatically deleted.
- **[Integrity]:** In-memory and on-disk state synchronized across all worktrees.
- **[Observability]:** Resize event recorded with timestamp in specification metadata.

### [US-3] Sizing-Aware Task and Coherence Validations
**User Story:** AS A Developer or AI Agent, I WANT the verification engines to adjust their strictness based on the specification size, SO THAT small tasks are not blocked by heavy ceremony while large architectural changes retain complete rigor.

**Scenarios:**
1. **[Happy Path - Small Spec Validation]** GIVEN a `small` specification containing a single-step task in `tasks.md` without requirement ID tags WHEN the validation suite runs THEN the check passes with zero errors.
2. **[Happy Path - Medium/Large/Complex Spec Validation]** GIVEN a `medium`, `large`, or `complex` specification WHEN validation runs THEN the engine strictly enforces mandatory requirement ID tagging (`[US-X]`), minimum action item density (at least 2 action steps), and acceptance checks.
3. **[Edge Case]** GIVEN a `small` specification that is later promoted to `large` WHEN the validation suite runs before updating `tasks.md` THEN the engine flags missing requirement linkages as actionable blocking errors.

**UI/UX Specifics:**
- **View/Component:** TUI Validation Error Drawer and CLI error summary.
- **Feedback Logic:** Red highlights on missing requirement mappings for `medium+` specs; green checkmark for compliant `small` specs.
- **Keybindings:** Standard error navigation keys.

**Technical Constraints (NFR):**
- **[Performance]:** Deterministic task scan and coherence check in under 30ms.
- **[Safety & Security]:** Read-only inspection during validation runs.
- **[Integrity]:** Consistent rule evaluation between CLI command and agent-driven audit.
- **[Observability]:** Structured validation messages with line numbers and expected guide snippets.

### [US-4] Agent Pre-Flight and Tiered Discovery Orchestration
**User Story:** AS A Developer collaborating with an AI Agent, I WANT the agent planning workflow to adapt its interrogation depth to the spec size, SO THAT quick tasks proceed immediately while complex architectural shifts undergo exhaustive multi-dimensional grilling.

**Scenarios:**
1. **[Happy Path]** GIVEN a user proposing a minor adjustment WHEN the agent detects a `small` scope THEN the agent confirms intent in a single turn, skips the 5-dimension grill gate, and drafts the implementation roadmap directly.
2. **[Happy Path - Complex Discovery]** GIVEN a major structural feature WHEN the agent identifies a `complex` scope THEN the agent initiates an adversarial consultative interview across architecture, security, data contracts, resilience, and UI before drafting artifacts.
3. **[Edge Case]** GIVEN mid-session discovery revealing hidden architectural complexity during a `small` spec planning WHEN new dependencies are uncovered THEN the agent suggests and triggers a promotion to `large` before drafting downstream technical blueprints.

**UI/UX Specifics:**
- **View/Component:** Interactive conversation prompts and question modals.
- **Feedback Logic:** Transparent explanation of sizing rationale and required document depth.
- **Keybindings:** Native interaction controls.

**Technical Constraints (NFR):**
- **[Performance]:** Immediate prompt adaptation with zero extraneous polling.
- **[Safety & Security]:** Zero prompt leakage in generated artifacts.
- **[Integrity]:** Compliance with project constitution and module boundaries.
- **[Observability]:** Explicit confirmation summaries before file generation.

## 4. Business Invariants
- A specification must always possess a valid size attribute (`small`, `medium`, `large`, or `complex`).
- Implementation mode (`specforce implement`) must never proceed if mandatory artifacts for the active size are missing or have failing validations.
- Downgrading a specification must never delete or truncate user-authored documents on disk.

## 5. Global UI/UX Contract (TUI Ghost Protocol)
- **Density Posture:** Standard (80x40) with responsive compact layout support.
- **Signature Moves:** Mint Green highlights on valid statuses, Warning Amber on promoted specs needing artifacts, Red borders on coherence errors.
- **Interaction Model:** Direct CLI flags (`--size`) and interactive selection in TUI spec manager.
- **State Behavior:** Clear indication of current spec size badge (`[SMALL]`, `[MED]`, `[LRG]`, `[CMPLX]`) in status overviews.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** All CLI metadata inspections and validations execute in < 50ms.
- **[Reliability]:** Fail-fast with clear error messages on malformed metadata or unrecognized sizes.
- **[Security]:** Strict 0600 file permissions and 0750 directory permissions.
- **[Maintainability]:** Clean separation of concerns across metadata, registry, validator, and CLI layers with comprehensive unit test coverage.
