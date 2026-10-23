# Delta: agent-kit
## ADDED
- `[NEW:exec-criteria]` Planning blueprints MUST direct every task acceptance check to be a single non-interactive command or test scenario with an observable expected result, rewriting or explicitly flagging criteria that need human judgment as not verifiable by a tester; `large` and `complex` roadmaps MUST end with a black-box verification phase whose scenarios trace to requirements and are runnable by the independent qa-tester.
- `[NEW:bug-red-first]` Bug specification blueprints MUST require concrete reproduction steps and an established root cause (halting to ask or recommend discovery when unknown), and the roadmap MUST start with a regression test task that fails before the fix, with no fix task preceding it.
- `[NEW:grill-waves]` Planning blueprints MUST restate decisions already taken in a proposal for a single confirmation and ask only open questions and newly found gaps; questions go in waves of at most four, each with a recommendation, and an approval of all recommendations confirms that wave only and never carries over; the hard gate before writing artifacts is preserved.
- `[NEW:readonly-explore]` Planning blueprints MAY delegate strictly read-only exploration to helpers during grounding of `large` and `complex` specifications, while interview, decisions and artifact writing remain in the primary session; unresolved markers (`TBD`, `TODO`, bare `N/A`) stay forbidden, and a section that does not apply is omitted or kept with a one-line justification.
- `[NEW:archive-preflight]` Archival instructions MUST refuse, before any retrospective, Constitution or module change, to archive a specification whose progress is below complete, listing each pending task with identifier, title and state and pointing to the implementation flow, even when the user asks to archive anyway.
- `[NEW:archive-closing]` Archival summaries MUST list the archived artifacts and changed module documents as clickable links to their archived or current location and MUST suggest, never execute, one conventional commit message.
### [NEW:us-exec-criteria] Executable Acceptance Checks and Black-box Verification
- **Scenario:** Planning a large specification
  - **GIVEN** a `large` or `complex` specification being planned
  - **WHEN** the roadmap is produced
  - **THEN** every task has a non-interactive acceptance check and the final phase holds black-box scenarios traced to requirements.
- **Scenario:** Criterion needing human judgment
  - **GIVEN** a task whose outcome cannot be scripted
  - **WHEN** its acceptance check is written
  - **THEN** it is rewritten into an observable check or flagged as not verifiable by a tester.
### [NEW:us-bug-red-first] Bug Specification Starts From a Failing Test
- **Scenario:** Planning a bug fix
  - **GIVEN** a bug specification
  - **WHEN** requirements and roadmap are produced
  - **THEN** reproduction and root cause are documented and the first task creates a regression test that fails before the fix.
- **Scenario:** Unknown root cause
  - **GIVEN** a bug whose root cause is not established
  - **WHEN** the planner reaches the root cause
  - **THEN** it halts to ask the user or recommends discovery instead of guessing.
### [NEW:us-grill-waves] Planning Interview in Waves
- **Scenario:** Proposal with decisions taken
  - **GIVEN** a proposal listing taken decisions and open questions
  - **WHEN** the interview starts
  - **THEN** taken decisions are restated for one confirmation and only open questions are asked, in waves with recommendations.
- **Scenario:** Approval does not carry over
  - **GIVEN** the user approved all recommendations of one wave
  - **WHEN** the next wave starts
  - **THEN** the next wave still requires its own confirmation.
### [NEW:us-readonly-explore] Read-only Grounding Helpers
- **Scenario:** Grounding a large specification
  - **GIVEN** a `large` or `complex` specification
  - **WHEN** the planner grounds itself in the codebase
  - **THEN** it may delegate read-only exploration while interview, decisions and writing stay in the primary session.
- **Scenario:** Section that does not apply
  - **GIVEN** a template section that does not apply
  - **WHEN** the artifact is drafted
  - **THEN** the section is omitted or justified in one line, and unresolved markers are rejected.
### [NEW:us-archive-preflight] Archive Refuses Incomplete Work
- **Scenario:** Incomplete implementation
  - **GIVEN** a specification with progress below complete
  - **WHEN** `/spf:archive` runs
  - **THEN** it stops before any change, lists each pending task and points to `/spf:implement`, even if the user asks to archive anyway.
### [NEW:us-archive-closing] Archive Summary With Links and Commit Suggestion
- **Scenario:** Successful archival
  - **GIVEN** a completed and archived specification
  - **WHEN** the final summary is produced
  - **THEN** it links archived artifacts and changed modules and suggests one conventional commit message without running it.
- **Scenario:** No module changes
  - **GIVEN** a specification with no module impact
  - **WHEN** the final summary is produced
  - **THEN** it states no module changed and still suggests the commit message.
