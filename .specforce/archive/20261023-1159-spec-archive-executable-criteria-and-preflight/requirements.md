---
slug: 20261023-1159-spec-archive-executable-criteria-and-preflight
lens: Integration
---

# Feature: Executable Criteria in Spec and Preflight/Closing Improvements in Archive

## 1. Context & Value
The independent QA tester proves every acceptance criterion with evidence, but planning produces criteria that are not always runnable and the planning conversation re-asks decisions already taken during discovery. Archival also trusts a single textual instruction to refuse incomplete work and closes without clickable links or a commit suggestion. This feature makes planned criteria executable by a tester, makes the planning interview lighter, and makes archival stricter at the start and more useful at the end.

## 2. Out of Scope (Anti-Goals)
- Do not change the implementation skill, the QA tester charter or the task status commands.
- Do not make archival run any version-control action; the commit message is only suggested.
- Do not allow archiving an incomplete specification, not even with a user override.
- Do not change the sizing categories or the module delta protocol.
- Do not change the Constitution; the allowed use of read-only exploration helpers is consistent with the existing "Primary Orchestration Only" principle.
- Do not grow the planning or archival skill blueprints beyond their current size discipline; text is folded into existing sections.

## 3. Acceptance Criteria (BDD)

### [US-1] Executable Acceptance Checks
**User Story:** AS A spec author, I WANT every task's acceptance check to be something a tester can run without interaction, SO THAT the independent QA can prove completion with evidence.

**Scenarios:**
1. **[Happy Path]** GIVEN a planned task roadmap WHEN the planner writes a task's acceptance check THEN it is a single non-interactive command or test scenario with an observable expected result.
2. **[Edge Case]** GIVEN a task whose outcome seems to need human judgment (for example a visual result) WHEN the planner writes its acceptance check THEN the planner rewrites it into an observable, scriptable check or records it explicitly as not verifiable by a tester, instead of leaving a vague instruction.

**UI/UX Specifics:**
- **View/Component:** Acceptance Check block of each task in the roadmap.
- **Feedback Logic:** Each check names the command or scenario and the expected observable result in one or two lines.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Performance]:** Each acceptance check is runnable in under 2 minutes on a developer machine.
- **[Integrity]:** Checks never require a prompt, a manual click or a pasted secret.

### [US-2] Black-box Verification Scenarios for Large Specs
**User Story:** AS A spec author, I WANT large and complex specifications to include end-to-end scenarios for the real deliverable, SO THAT the QA tester can exercise it black-box, beyond unit tests.

**Scenarios:**
1. **[Happy Path]** GIVEN a `large` or `complex` specification WHEN the roadmap is produced THEN it ends with a distinct black-box verification phase whose scenarios each trace to a requirement and state a runnable command plus the expected observable result.
2. **[Edge Case]** GIVEN a `small` or `medium` specification WHEN the roadmap is produced THEN no black-box verification phase is required and none is imposed.

**UI/UX Specifics:**
- **View/Component:** Final phase of the roadmap named for black-box verification.
- **Feedback Logic:** Every scenario lists its requirement reference, its command and its expected result.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Performance]:** The whole black-box phase is runnable by a tester in under 10 minutes.
- **[Integrity]:** Requirements remain free of technical detail; the commands appear only in the roadmap.

### [US-3] Bug Specs Start From Reproduction and a Failing Test
**User Story:** AS A maintainer fixing a defect, I WANT a bug specification to record how to reproduce it and its root cause and to start with a failing regression test, SO THAT the fix is proven and the defect cannot silently return.

**Scenarios:**
1. **[Happy Path]** GIVEN a bug specification WHEN the planner drafts it THEN the requirements contain reproduction steps and a root cause, and the first task of the roadmap creates a regression test that fails before the fix.
2. **[Edge Case]** GIVEN a bug whose root cause is not yet established WHEN the planner reaches the root cause section THEN the planner halts and asks the user (or recommends discovery) instead of writing a guess.

**UI/UX Specifics:**
- **View/Component:** Reproduction and root cause sections of the bug requirements; first task of the roadmap.
- **Feedback Logic:** The first task states the expected failing outcome before the fix.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Performance]:** The regression test runs in under 1 minute.
- **[Integrity]:** No fix task precedes the failing regression test task.

### [US-4] Proposal Pre-answers the Planning Interview
**User Story:** AS A spec author, I WANT decisions already taken in a proposal to be restated rather than re-asked, SO THAT the interview focuses on open questions only.

**Scenarios:**
1. **[Happy Path]** GIVEN a proposal listing decisions taken and open questions WHEN the planning interview starts THEN the planner restates the taken decisions for one confirmation and asks only the open questions and any newly discovered gaps.
2. **[Edge Case]** GIVEN a taken decision that contradicts the Constitution or the current code WHEN the planner detects it THEN the planner raises it as a question with a recommendation instead of silently following it.

**UI/UX Specifics:**
- **View/Component:** Interview question waves.
- **Feedback Logic:** Each question shows a recommended option; restated decisions are shown as a single confirmation item.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Performance]:** A medium specification interview needs no more than two question waves when a proposal exists.
- **[Integrity]:** No artifact is written before every wave is confirmed.

### [US-5] Interview in Waves With Recommendations
**User Story:** AS A spec author, I WANT questions grouped in waves, each with a recommendation, and a way to approve all recommendations of a wave at once, SO THAT confirming sound defaults is fast without losing explicit consent.

**Scenarios:**
1. **[Happy Path]** GIVEN a wave of questions with recommendations WHEN the user says to approve all recommendations THEN every question of that wave counts as explicitly confirmed.
2. **[Edge Case]** GIVEN an approval given for one wave WHEN a later wave starts THEN the approval does not carry over and the new wave still requires its own confirmation.

**UI/UX Specifics:**
- **View/Component:** Question wave presented through the interactive question tool.
- **Feedback Logic:** The recommended option is listed first and labeled as recommended.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Performance]:** At most four questions per wave.
- **[Integrity]:** Vague answers are re-asked; the hard gate before writing artifacts is preserved.

### [US-6] Read-only Exploration Helpers During Grounding
**User Story:** AS A spec author planning a large change, I WANT the planner to be able to use read-only exploration helpers while grounding, SO THAT broad code reading does not flood the main conversation.

**Scenarios:**
1. **[Happy Path]** GIVEN a `large` or `complex` specification WHEN the planner grounds itself in the codebase THEN it may delegate read-only exploration and receives findings back, while interview, decisions and artifact writing stay in the main session.
2. **[Edge Case]** GIVEN a `small` or `medium` specification, or a helper that would write or change files WHEN grounding runs THEN no helper is used for planning, and any helper that is used is read-only.

**UI/UX Specifics:**
- **View/Component:** Grounding step of the planning flow.
- **Feedback Logic:** The planner states in one line what it delegated to read-only exploration.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Safety & Security]:** Helpers have no authority to write files, run mutating commands or ask the user questions.
- **[Integrity]:** Planning, interview and artifact writing are never delegated.

### [US-7] Placeholder Ban Without the Literal N/A Conflict
**User Story:** AS A spec author, I WANT sections that genuinely do not apply to be omitted or briefly justified, SO THAT artifacts are free of empty markers yet honest about non-applicable content.

**Scenarios:**
1. **[Happy Path]** GIVEN a template section that does not apply WHEN the planner drafts the artifact THEN the section is omitted or kept with a one-line justification.
2. **[Edge Case]** GIVEN a section left with an unresolved marker such as "to be defined" or a bare "N/A" WHEN the planner audits the artifact THEN the artifact is treated as incomplete and regenerated after clarification.

**UI/UX Specifics:**
- **View/Component:** Artifact sections.
- **Feedback Logic:** Justifications are a single line stating why the section does not apply.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Integrity]:** Unresolved markers (TBD, TODO, bare N/A) remain forbidden in every artifact.

### [US-8] Archive Refuses Incomplete Work Up Front
**User Story:** AS A maintainer, I WANT archival to stop immediately and list what is pending when the implementation is not fully complete, SO THAT no knowledge is harvested from unfinished work.

**Scenarios:**
1. **[Happy Path]** GIVEN a specification with complete progress WHEN archival starts THEN it proceeds to the retrospective.
2. **[Edge Case]** GIVEN a specification with progress below complete, including blocked or unfinished tasks WHEN archival starts THEN it stops before any retrospective, Constitution or module change, lists each pending task with its identifier and title, and directs the user to the implementation flow.
3. **[Edge Case]** GIVEN the user asks to archive anyway WHEN the specification is incomplete THEN archival still refuses.

**UI/UX Specifics:**
- **View/Component:** Archive preflight message.
- **Feedback Logic:** One line per pending task: identifier, title, state; followed by the hand-off command for implementation.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Performance]:** Preflight needs exactly one status query and produces no file changes.
- **[Integrity]:** The refusal happens before any write.

### [US-9] Archive Summary With Links and Commit Suggestion
**User Story:** AS A maintainer, I WANT the closing summary to link the archived artifacts and changed module documents and to suggest a commit message, SO THAT I can review and commit the result quickly.

**Scenarios:**
1. **[Happy Path]** GIVEN a successful archival WHEN the final summary is produced THEN it lists the archived artifacts and every changed module document as clickable links and proposes one conventional commit message.
2. **[Edge Case]** GIVEN a specification with no module changes WHEN the final summary is produced THEN it states that no module documents changed and still proposes the commit message.
3. **[Edge Case]** GIVEN any archival WHEN the summary is produced THEN the agent only suggests the commit message and never commits.

**UI/UX Specifics:**
- **View/Component:** Final archive summary.
- **Feedback Logic:** The commit suggestion appears as a single code block after the links and before the next steps.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Integrity]:** Links point to the archived location, not the removed active location.
- **[Safety & Security]:** No version-control command is executed.

## 4. Business Invariants
- Every roadmap task keeps a non-interactive, observable acceptance check.
- A bug specification's first roadmap task is a failing regression test.
- No artifact is written before the interview is confirmed; approval of one wave never covers another.
- An incomplete specification is never archived.
- Read-only helpers never plan, interview or write.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Archive preflight adds a single status query.
- **[Reliability]:** Refusal and unresolved-marker handling fail fast with a clear message.
- **[Security]:** Archival and planning never run version-control or mutating commands on behalf of read-only helpers.
- **[Maintainability]:** Skill blueprints keep their current size discipline; rules are asserted by the manifest tests and the full test suite passes.
