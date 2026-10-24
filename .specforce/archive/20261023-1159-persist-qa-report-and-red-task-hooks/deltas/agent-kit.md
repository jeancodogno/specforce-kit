# Delta: agent-kit
## ADDED
- `[NEW:qa-record-flow]` Implementation orchestration blueprints MUST direct the orchestrator, never the independent qa-tester, to persist the QA verdict through the specification metadata after the QA step and again after every correction cycle, before declaring validation complete.
- `[NEW:archive-uses-qa]` Archival blueprints and instructions MUST read the persisted QA state, refuse to archive when it is FAILED, warn visibly when it is absent, and derive next-step suggestions from the persisted findings and NOT-VERIFIABLE criteria instead of the conversation.
- `[NEW:red-title-marker]` Planning blueprints MUST direct test-first roadmaps to mark RED tasks with the RED token in the task title so the scoped hook relaxation applies, and finishing a RED task MUST be announced with a notice naming the skipped per-task hooks.
### [NEW:us-qa-record-flow] Orchestrator Persists the QA Verdict
- **Scenario:** QA step completed
  - **GIVEN** the independent qa-tester returned its traceability matrix and findings
  - **WHEN** the orchestrator closes validation
  - **THEN** it records the report before declaring validation complete and records again after each correction cycle.
- **Scenario:** Tester attempts to persist
  - **GIVEN** the qa-tester role charter
  - **WHEN** the tester reports its results
  - **THEN** it only returns them to the orchestrator and writes nothing.
### [NEW:us-archive-uses-qa] Archival Reads Persisted QA
- **Scenario:** Archival in a fresh session
  - **GIVEN** a specification with a persisted PASSED report that has findings
  - **WHEN** archival produces its next-step suggestions
  - **THEN** the suggestions cite the persisted findings and NOT-VERIFIABLE criteria.
- **Scenario:** Persisted report is FAILED
  - **GIVEN** a specification whose persisted report is FAILED
  - **WHEN** archival starts
  - **THEN** it stops before any change, lists the failed criteria and points to fixing and re-recording.
- **Scenario:** No persisted report
  - **GIVEN** a specification without a persisted report
  - **WHEN** archival starts
  - **THEN** it warns that no QA evidence exists and continues.
### [NEW:us-red-title-marker] RED Tasks Are Marked in Planning
- **Scenario:** Planning a test-first pair
  - **GIVEN** a roadmap with RED and GREEN task pairs
  - **WHEN** the roadmap is written
  - **THEN** every RED task carries the RED token in its title.
