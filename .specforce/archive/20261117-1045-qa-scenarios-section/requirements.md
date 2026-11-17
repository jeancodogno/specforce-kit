---
slug: 20261117-1045-qa-scenarios-section
lens: Backend-heavy
---

# Feature: Black-box Scenarios as QA Input

## 1. Context & Value
Roadmaps for `large` and `complex` specifications end with a mandatory verification phase made of tasks, and roadmaps of every size tend to close with a "run the full suite" task. Implementers execute those tasks and the verification hooks fire, and then the independent qa-tester repeats the same suite, smoke runs and acceptance scenarios in the final QA step, so quality assurance runs twice per implementation. This feature moves black-box scenarios out of the task list into a non-executable scenario section that feeds the qa-tester, and forbids verification-only tasks, so verification happens once, in the independent QA step, while keeping scenarios executable and traced to requirements.

## 2. Out of Scope (Anti-Goals)
*Explicitly state what MUST NOT be built in this specific feature to prevent scope creep.*
- Do not change the roadmap validator or the task parser; no new warning or error is emitted for verification tasks. The only status-report change allowed is the mitigations-section recognition required by [US-5].
- Do not rewrite or migrate archived or in-flight specifications; roadmaps that already contain a verification phase with verification tasks keep working unchanged.
- Do not add a coverage gate to QA recording (checking that every requirement, delta item or scenario has a matrix row); recording keeps accepting any criterion identifier.
- Do not change the shape of the specification metadata or of the persisted QA report.
- Do not change the QA tiers, the correction cycle limits, batch budgets, the per-task acceptance check rule or the bug-first regression rule.
- Do not change the requirements or design templates.

## 3. Acceptance Criteria (BDD)
*Every requirement MUST be testable. Define success and failure.*

### [US-1] Black-box Scenarios Live Outside the Task List
**User Story:** AS A spec author, I WANT TO declare black-box scenarios in a dedicated, non-executable section of the roadmap, SO THAT the independent QA runs them once instead of implementers running them as tasks.

**Scenarios:**
1. **[Happy Path]** GIVEN a `large` or `complex` specification being planned WHEN the roadmap is produced THEN it contains a "Black-box Scenarios" section placed after the task phases, where each scenario has a scenario identifier distinct from task identifiers, a link to at least one requirement, one non-interactive command and its expected observable result, and no phase named for verification exists.
2. **[Happy Path]** GIVEN a `medium` specification whose change produces something runnable (a command, a service or a screen) WHEN the roadmap is produced THEN the scenario section MAY be included; GIVEN a `medium` specification without such an artifact THEN it is omitted.
3. **[Edge Case]** GIVEN a `small` specification WHEN the roadmap is produced THEN no scenario section is rendered.
4. **[Edge Case]** GIVEN a roadmap with a scenario section WHEN it is validated and its progress is computed THEN the scenarios are not counted as tasks and the task count is identical to the same roadmap without the section.

**UI/UX Specifics:**
- **View/Component:** Implementation roadmap document and the specification status report.
- **Feedback Logic:** Scenarios read as a plain list (no checkboxes); progress shown in the status report reflects only implementation tasks.
- **Keybindings:** Not applicable.

**Technical Constraints (NFR):**
- **[Performance]:** Roadmap validation keeps a single pass with no added cost relative to today.
- **[Integrity]:** A scenario never carries a checkbox or a task identifier, so no tooling can mark, time or hook it as a task.

### [US-2] No Verification-Only Tasks
**User Story:** AS A maintainer running an implementation, I WANT TO roadmaps to contain only tasks that change the product, SO THAT suites, smoke runs and repository-wide checks are executed once by the final QA and the all-tasks hooks rather than by an implementer.

**Scenarios:**
1. **[Happy Path]** GIVEN any specification size WHEN the roadmap is produced THEN every task changes at least one file or artifact of the product, and no task exists whose only purpose is to run the full test suite, a smoke run, a build check or repository-wide searches.
2. **[Happy Path]** GIVEN a task that changes a file WHEN its acceptance check is written THEN the check stays scoped to proving that task (one non-interactive command with an observable result), as today.
3. **[Edge Case]** GIVEN a planner tempted to close the roadmap with a final "full verification" task WHEN the planning coherence audit runs THEN that task is removed and its checks, when they describe runnable end-to-end behavior, are expressed as scenarios instead (for sizes that allow the section).

**UI/UX Specifics:**
- **View/Component:** Planning guidance for roadmaps and the planning coherence audit.
- **Feedback Logic:** The audit names verification-only tasks it removes or converts in the planning summary.
- **Keybindings:** Not applicable.

**Technical Constraints (NFR):**
- **[Performance]:** The rule removes at least one batch-level verification run per implementation compared with today's roadmaps that end in a verification phase.
- **[Integrity]:** The rule is expressed in the roadmap planning guidance and the planning skill; it is not enforced by the validator.

### [US-3] The Independent QA Executes the Scenarios
**User Story:** AS A maintainer, I WANT TO the independent qa-tester to run every declared scenario and record it in the traceability matrix, SO THAT black-box evidence is produced exactly once and is persisted with the QA verdict.

**Scenarios:**
1. **[Happy Path]** GIVEN a roadmap with a scenario section WHEN the qa-tester runs the black-box tier THEN it executes every scenario's command and adds one matrix row per scenario, identified by the scenario identifier, with the evidence and a verdict `PASSED`, `FAILED` or `NOT-VERIFIABLE`, alongside the existing requirement and delta rows.
2. **[Happy Path]** GIVEN a roadmap without a scenario section WHEN the qa-tester runs the black-box tier THEN it derives the black-box checks from the requirements' acceptance scenarios, as it does today.
3. **[Edge Case]** GIVEN a scenario whose command cannot be executed in the environment WHEN the qa-tester evaluates it THEN its row is `NOT-VERIFIABLE` with a stated reason and is never reported as passed.
4. **[Edge Case]** GIVEN a scenario whose observed result differs from the expected one WHEN the matrix is built THEN its row is `FAILED` and it enters the existing correction cycle like any failed criterion.

**UI/UX Specifics:**
- **View/Component:** Traceability matrix in the QA report and the persisted QA verdict.
- **Feedback Logic:** Scenario rows appear with their own identifiers next to requirement and delta rows in the final completion output.
- **Keybindings:** Not applicable.

**Technical Constraints (NFR):**
- **[Performance]:** Each scenario is a single command runnable in under 2 minutes, matching the acceptance check limit.
- **[Integrity]:** Scenario rows are persisted through the existing QA recording without any change to the report shape.

### [US-4] Existing Roadmaps Keep Working
**User Story:** AS A maintainer with archived and in-flight specifications, I WANT TO keep implementing roadmaps that still have a verification phase with verification tasks, SO THAT nothing already planned breaks.

**Scenarios:**
1. **[Happy Path]** GIVEN an existing roadmap whose last phase holds verification tasks WHEN it is validated and implemented THEN it is accepted and its tasks are executed as ordinary tasks.
2. **[Edge Case]** GIVEN an existing roadmap with verification tasks WHEN its status is reported THEN no new warning or error about verification tasks appears.

**UI/UX Specifics:**
- **View/Component:** Specification status report.
- **Feedback Logic:** Output for existing roadmaps is unchanged.
- **Keybindings:** Not applicable.

**Technical Constraints (NFR):**
- **[Performance]:** No change in validation cost.
- **[Integrity]:** Backward compatibility is total: no validator, parser or metadata change.

### [US-5] Pre-emptive Mitigations Survive Section Renumbering
**User Story:** AS A maintainer reading the implementation status, I WANT TO see a roadmap's pre-emptive mitigations whatever number its section carries, SO THAT adding the scenario section before it does not silently hide execution risks.

**Scenarios:**
1. **[Happy Path]** GIVEN a roadmap whose mitigations section is numbered 4 (after a scenario section) WHEN the implementation status is requested THEN the mitigations text is shown.
2. **[Happy Path]** GIVEN an existing roadmap whose mitigations section is numbered 3 WHEN the implementation status is requested THEN the mitigations text is shown exactly as before.
3. **[Edge Case]** GIVEN a roadmap without a mitigations section, or with a scenario section following the mitigations WHEN the implementation status is requested THEN the mitigations are empty in the first case and, in the second, contain only the mitigations text without any scenario text.

**UI/UX Specifics:**
- **View/Component:** Implementation status report (human and machine-readable).
- **Feedback Logic:** The mitigations block appears whenever the roadmap declares one, regardless of its section number.
- **Keybindings:** Not applicable.

**Technical Constraints (NFR):**
- **[Performance]:** Recognition stays a single pass over the roadmap with no added cost.
- **[Integrity]:** The section is identified by its title, not its number; the extracted text stops at the next top-level section.

## 4. Business Invariants
*Strict domain rules that the system must never violate, regardless of the user story. These must be technically testable.*
- A black-box scenario is never a task: it has no checkbox, no task identifier, no time log and triggers no task hook.
- Every scenario links to at least one requirement and states one non-interactive command with an expected observable result.
- Full-suite, smoke and repository-wide verification belongs to the independent QA step and the all-tasks hooks, never to a roadmap task.
- `large` and `complex` roadmaps always declare scenarios; `small` roadmaps never do.

## 6. Global Non-Functional Requirements (NFRs)
*Technical constraints and quality gates. Aligned with .specforce/docs/engineering.md*
- **[Performance]:** Planning and validation cost does not increase; implementation runs one fewer verification pass for roadmaps that previously ended in a verification phase.
- **[Reliability]:** Kit manifest tests guard the new planning and QA directives so they cannot silently regress.
- **[Security]:** No new commands, inputs or data are introduced.
- **[Maintainability]:** The installed skill copies in the repository match the kit sources after the change.
