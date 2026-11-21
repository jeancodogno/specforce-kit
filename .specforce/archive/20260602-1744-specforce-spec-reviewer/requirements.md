---
slug: 20260602-1744-specforce-spec-reviewer
lens: Backend-heavy
---

# Feature: Specforce Spec Reviewer

## 1. Context & Value
The goal of this feature is to enforce a "Quality-First" gate in the Spec-Driven Development (SDD) pipeline. Currently, spec artifacts (`requirements.md`, `design.md`, and `tasks.md`) are generated sequentially, which can lead to logical drift (e.g., a technical decision in the design is not reflected in the tasks, or a requirement is missing an implementation step). 

By introducing a specialized `specforce-spec-reviewer` agent and an automated refinement loop within the `spf.spec` command, the system will automatically detect and attempt to fix these inconsistencies before the user is asked to approve the specification for implementation.

## 2. Out of Scope (Anti-Goals)
- Manual user-triggered review during the loop (the loop is automatic).
- Code implementation or code review (this is strictly for planning artifacts).
- Automated correction of the Constitution (Reviewer audits against it, but does not modify it).
- Reviewing artifacts outside of the `.specforce/specs/{slug}/` directory.

## 3. Acceptance Criteria (BDD)

### [US-1] Coherence Auditing Agent (`specforce-spec-reviewer`)
The system MUST include a dedicated agent profile capable of performing cross-artifact logical audits.

**Acceptance Criteria:**
- **GIVEN** a specification set containing `requirements.md`, `design.md`, and `tasks.md`.
- **WHEN** the `specforce-spec-reviewer` is triggered.
- **THEN** it MUST verify that every Functional Requirement `[US-X]` in `requirements.md` has at least one corresponding execution step in `tasks.md`.
- **AND** it MUST verify that all high-level technical decisions (e.g., "Use SQLite") from `design.md` are accounted for in the task roadmap.
- **AND** it MUST output any detected gaps as standardized `[COHERENCE_ERROR]` blocks.

**Edge Case:**
- **GIVEN** a requirement `[US-99]` that is explicitly marked as "Out of Scope".
- **WHEN** audited by the Reviewer.
- **THEN** the agent SHOULD NOT flag it as missing from the `tasks.md`.

**Technical Constraints (NFR):**
- **[Performance]:** Audit execution for a standard feature spec SHALL complete in < 8 seconds.
- **[Reliability]:** The agent MUST use a deterministic parsing strategy to identify `[US-X]` tags.

---

### [US-2] Automated Refinement Loop
The `spf.spec` command pipeline MUST integrate a Phase 4 "Refinement Loop" that executes after artifact generation.

**Acceptance Criteria:**
- **GIVEN** that the initial artifact generation (Phase 2) is complete.
- **WHEN** the `specforce-spec-reviewer` identifies one or more `[COHERENCE_ERROR]` markers.
- **THEN** the system MUST automatically trigger a "Refinement Loop".
- **AND** the system MUST re-invoke the specific agent (Analyst, Architect, or Planner) responsible for the artifact containing the error.

**Edge Case:**
- **GIVEN** an error that involves a contradiction between Requirements and Design.
- **WHEN** the loop starts.
- **THEN** the system SHOULD prioritize fixing the "Source of Truth" (Requirements) or prompt for clarification if the logic is irreconcilable.

**Technical Constraints (NFR):**
- **[Performance]:** The loop overhead (checking for errors when none exist) SHALL be < 2 seconds.
- **[State Management]:** The `iteration_count` MUST be stored in the specification's `spec.yaml` to prevent infinite cycles.

---

### [US-3] Surgical Correction Protocol
When re-invoking agents during a refinement loop, the system MUST provide localized context to ensure surgical fixes.

**Acceptance Criteria:**
- **GIVEN** a `[COHERENCE_ERROR]` detected in `tasks.md` regarding a missing step for `[US-3]`.
- **WHEN** the `specforce-planner` is re-activated for correction.
- **THEN** the prompt MUST include the specific error text and the full content of the requirement `[US-3]` it needs to satisfy.

**Edge Case:**
- **GIVEN** a correction that creates a new inconsistency in a previously verified file.
- **WHEN** the next audit runs.
- **THEN** the system MUST detect the new error and attempt a fix in the subsequent loop iteration.

**Technical Constraints (NFR):**
- **[Performance]:** Surgical correction prompts MUST be token-optimized, sending only relevant artifact snippets plus the target error.

---

### [US-4] Loop Termination & Escalation
The system MUST gracefully handle scenarios where automatic refinement fails to achieve coherence.

**Acceptance Criteria:**
- **GIVEN** that the Refinement Loop has executed 3 times.
- **WHEN** the `specforce-spec-reviewer` still reports `[COHERENCE_ERROR]` markers.
- **THEN** the system MUST terminate the automatic loop.
- **AND** the system MUST present all remaining errors to the human user for manual intervention.
- **AND** the specification status MUST be marked as `is_valid: false`.

**Edge Case:**
- **GIVEN** the loop achieves 100% coherence on the 2nd iteration.
- **WHEN** the audit passes.
- **THEN** the system MUST immediately exit the loop and proceed to the final verification/handoff summary.

**Technical Constraints (NFR):**
- **[Safety]:** The hard limit of 3 iterations SHALL NOT be bypassable by the agents.
- **[UX]:** Escalation messages MUST clearly identify which artifact is blocking the "Valid" state.

---
### [US-5] CLI Audit Command
The system MUST provide a CLI command to persist and clear coherence errors in the `spec.yaml` file.

**Acceptance Criteria:**
- **GIVEN** a specification slug.
- **WHEN** the command `specforce spec audit <slug> --error "Message"` is executed.
- **THEN** the message MUST be appended to the `refinement.errors` list in `spec.yaml`.
- **AND** if executed with `--clear`, the `errors` list MUST be emptied.

**Technical Constraints (NFR):**
- **[Safety]:** The command MUST validate that the spec directory exists before writing.
- **[Integrity]:** Concurrent writes to `spec.yaml` MUST be handled gracefully (single-process CLI assumption).

## 4. Global Non-Functional Requirements (NFRs)
- **[Performance]:** The entire `spf.spec` pipeline, including up to 3 refinement loops, SHOULD complete in < 2 minutes for a medium-sized feature.
- **[Observability]:** Every iteration of the loop MUST be logged in the terminal/TUI with an "Iteration X/3" indicator.
- **[Reliability]:** The refinement loop MUST NOT result in data loss or accidental deletion of approved sections of the artifacts.
