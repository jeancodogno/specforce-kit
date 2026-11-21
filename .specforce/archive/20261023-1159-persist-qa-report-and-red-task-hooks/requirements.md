---
slug: 20261023-1159-persist-qa-report-and-red-task-hooks
lens: Backend-heavy
---

# Feature: Persist QA Report and Scoped Hooks for RED Tasks

## 1. Context & Value
The independent QA verdict (criterion-by-criterion traceability and out-of-scope findings) exists only in the conversation, so archiving from a new session loses the findings and wrongly reports that none exist. Separately, finishing a RED task (tests written before the code they exercise) is blocked by the project's per-task verification hooks, which forces teams to finish all tasks together at the end. This feature makes QA evidence durable and gates archival on it, and lets RED tasks finish through an explicit, visible and narrowly scoped relaxation.

## 2. Out of Scope (Anti-Goals)
- The QA tester does not gain write access to anything; only the orchestrating agent records the report.
- No history of past QA runs is kept; only the latest report is stored.
- No new planning artifact and no change to the required artifact list or dependency order of any specification size.
- No project-wide switch that disables verification hooks, and no per-command override flag for them.
- Phase-level and all-tasks hooks are not relaxed for RED tasks.
- The local tooling problem where a snap-installed linter shadows the user-installed one is not addressed by changing configuration.

## 3. Acceptance Criteria (BDD)

### [US-1] Record the QA Report
**User Story:** AS A maintainer driving a specification through implementation, I WANT the orchestrating agent to save the independent QA verdict, SO THAT the evidence survives beyond the conversation.

**Scenarios:**
1. **[Happy Path]** GIVEN a specification whose QA tester returned a traceability matrix and out-of-scope findings WHEN the orchestrator records the report THEN the specification stores each criterion with its scenario, evidence and verdict, the findings, the overall status and the time of recording.
2. **[Happy Path]** GIVEN a recorded report where every criterion is PASSED or NOT-VERIFIABLE WHEN the overall status is derived THEN the status is PASSED.
3. **[Edge Case]** GIVEN a report containing at least one FAILED criterion WHEN it is recorded THEN the overall status is FAILED.
4. **[Edge Case]** GIVEN a report with an unknown verdict value, no criteria, or an unreadable structure WHEN recording is attempted THEN the recording is refused with a message naming the problem and the previously stored report, if any, is left untouched.
5. **[Edge Case]** GIVEN a specification that already has a stored report WHEN a new report is recorded after a correction cycle THEN the new report fully replaces the previous one and a run counter increases by one.
6. **[Edge Case]** GIVEN a specification that does not exist or is already archived WHEN recording is attempted THEN the recording is refused and nothing is written.

**UI/UX Specifics:**
- **View/Component:** One-line confirmation stating the overall status, the number of criteria per verdict and the number of findings.
- **Feedback Logic:** Refusals are shown as an error naming the offending field or criterion.
- **Keybindings:** Not applicable; non-interactive command.

**Technical Constraints (NFR):**
- **[Performance]:** Recording completes in under 50ms for up to 200 criteria.
- **[Safety & Security]:** The stored metadata keeps owner-only permissions and the write replaces the file atomically.
- **[Integrity]:** Task states, time logs and delta bookkeeping already stored for the specification are preserved unchanged.
- **[Observability]:** The confirmation line is the only output on success.

### [US-2] See the QA State in Implementation Status
**User Story:** AS A maintainer or agent resuming work in a new session, I WANT the implementation status to show the stored QA state, SO THAT I know whether validation happened and what it found.

**Scenarios:**
1. **[Happy Path]** GIVEN a specification with a stored report WHEN the implementation status is requested in machine-readable form THEN it includes the overall status, the verdict counts, the findings with their titles and the recording time.
2. **[Edge Case]** GIVEN a specification without a stored report WHEN the implementation status is requested THEN the QA information is absent rather than reported as passed.
3. **[Edge Case]** GIVEN a stored report whose metadata is partly unreadable WHEN the implementation status is requested THEN the status command still reports the tasks and signals that QA information could not be read.

**UI/UX Specifics:**
- **View/Component:** A QA line in the human-readable status view showing status and counts.
- **Feedback Logic:** FAILED is rendered in the error color, PASSED in the success color, absent as a muted "not recorded".
- **Keybindings:** Unchanged.

**Technical Constraints (NFR):**
- **[Performance]:** No measurable added latency (status stays under 50ms).
- **[Integrity]:** Absence of a report is never presented as success.

### [US-3] Archive Uses Persisted QA Findings
**User Story:** AS A maintainer archiving in a different session than the one that validated, I WANT archival to read the stored QA outcome, SO THAT next-step suggestions and the refusal rule reflect what really happened.

**Scenarios:**
1. **[Happy Path]** GIVEN a specification with a stored PASSED report that has out-of-scope findings and NOT-VERIFIABLE criteria WHEN archival runs in a fresh session THEN the suggested follow-ups cite those stored findings and criteria.
2. **[Edge Case]** GIVEN a specification whose stored report is FAILED WHEN archival is requested, with or without the override for pending tasks, THEN archival is refused, the failed criteria are listed with the next step (fix, re-validate and re-record), and the specification stays active.
3. **[Edge Case]** GIVEN a specification with no stored report WHEN archival is requested THEN a visible warning states that no QA evidence exists and archival proceeds.
4. **[Edge Case]** GIVEN a stored report with no findings and no unverifiable criteria WHEN follow-ups are suggested THEN the suggestions state that none were found and are based on the completed work.

**UI/UX Specifics:**
- **View/Component:** Refusal message listing each FAILED criterion as identifier and scenario, followed by the command to re-record.
- **Feedback Logic:** Warning for missing report; error for FAILED.
- **Keybindings:** Not applicable.

**Technical Constraints (NFR):**
- **[Performance]:** The QA check adds under 10ms to archival.
- **[Integrity]:** The refusal for FAILED is never bypassed by the override that only skips pending tasks.
- **[Safety & Security]:** Archival never modifies the stored report.

### [US-4] Finish RED Tasks Without False Blocking
**User Story:** AS A maintainer practicing test-first pairs, I WANT a task marked as RED to finish without running the per-task verification hooks, SO THAT tests written ahead of their implementation do not block progress.

**Scenarios:**
1. **[Happy Path]** GIVEN a task whose title carries the RED marker and per-task hooks are configured WHEN the task is finished THEN the per-task hooks are not run, the task is finished, and a notice names every skipped hook.
2. **[Happy Path]** GIVEN the GREEN task paired with that RED task WHEN it is finished THEN the per-task hooks run and a failure blocks it as before.
3. **[Edge Case]** GIVEN a batch containing RED and non-RED tasks WHEN the batch is finished THEN the per-task hooks run once because at least one task is not RED.
4. **[Edge Case]** GIVEN a RED task that is also the last of its phase or of the whole roadmap WHEN it is finished THEN the phase-level and all-tasks hooks still run and can block it.
5. **[Edge Case]** GIVEN a task without the RED marker, or a marker written in different letter case or in the description rather than the title WHEN it is finished THEN hooks behave exactly as today.
6. **[Edge Case]** GIVEN a RED task and no configured per-task hooks WHEN it is finished THEN no skip notice is shown.
7. **[Edge Case]** GIVEN a RED task being moved to in-progress or pending WHEN the status is updated THEN behavior is unchanged because hooks only apply to finishing.

**UI/UX Specifics:**
- **View/Component:** A warning line "RED task: skipped per-task hooks: <hook list>" shown before the success confirmation.
- **Feedback Logic:** The warning is shown every time a skip occurs, never silently.
- **Keybindings:** Not applicable.

**Technical Constraints (NFR):**
- **[Performance]:** Skipping removes the hook run time entirely; detection adds under 5ms.
- **[Integrity]:** A task is only finished after verification passes or after an explicit, visible, task-scoped skip.
- **[Observability]:** The skipped hooks are named in the notice.

### [US-5] Teach Agents the New Flow
**User Story:** AS A maintainer relying on the planning, implementation and archival skills, I WANT the agents' instructions to use the stored report and the RED marker, SO THAT the behavior happens without manual reminders.

**Scenarios:**
1. **[Happy Path]** GIVEN the implementation orchestrator finished the independent QA step WHEN it closes validation THEN its instructions require it to record the report before declaring validation complete, and repeat it after every correction cycle.
2. **[Happy Path]** GIVEN the archival instructions WHEN a specification is archived THEN they direct the agent to read the stored QA state, refuse on FAILED, warn on missing, and derive next steps from the stored findings instead of the conversation.
3. **[Happy Path]** GIVEN the planning instructions for roadmaps WHEN test-first pairs are planned THEN RED tasks carry the RED marker in their title so the relaxation applies.
4. **[Edge Case]** GIVEN the QA tester role WHEN it reports THEN its instructions still forbid it from writing the report, which only the orchestrator does.

**UI/UX Specifics:**
- **View/Component:** Not applicable; text blueprints only.
- **Feedback Logic:** Not applicable.
- **Keybindings:** Not applicable.

**Technical Constraints (NFR):**
- **[Performance]:** Not applicable to blueprint text; the kit installs in under 1s as today.
- **[Maintainability]:** Blueprint behavior is covered by automated tests that fail if the instructions drop the record, read or refuse rules.

## 4. Business Invariants
- The QA tester never persists anything; the orchestrator is the only writer of the report.
- Only the latest QA report exists for a specification; recording replaces it.
- A FAILED stored report always blocks archival, whatever other overrides are used.
- A missing QA report never counts as a pass and always produces a visible warning at archival and a visible absence in status.
- The relaxation for RED tasks is limited to per-task hooks, applies only to tasks explicitly marked in their title, and is announced every time it is used.
- Finishing an unmarked or GREEN task always runs the configured verification.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Every added check or write completes in under 50ms.
- **[Reliability]:** Failures are explicit and leave previously stored data untouched; no silent fallbacks.
- **[Security]:** Stored metadata keeps owner-only permissions; specification names are confined to the project's specification directory.
- **[Maintainability]:** New behavior is covered by unit tests for the happy and edge scenarios above, keeping the project coverage gate.
