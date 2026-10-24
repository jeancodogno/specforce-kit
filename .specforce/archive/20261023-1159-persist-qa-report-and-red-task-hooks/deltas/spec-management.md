# Delta: spec-management
## ADDED
- `[NEW:qa-persist]` A specification stores at most one QA report in its metadata: criteria with scenario, evidence and verdict (PASSED, FAILED or NOT-VERIFIABLE), out-of-scope findings, a run counter and a recording time. The overall status is derived by the system (FAILED when any criterion FAILED, otherwise PASSED); invalid input is refused without altering the stored report, and recording replaces the previous report atomically while preserving all other metadata.
- `[NEW:qa-status-exposed]` The implementation status exposes the stored QA state (status, verdict counts, failed and not-verifiable criteria, findings) and omits it when no report exists; an unreadable QA block is reported as an error without hiding task data.
- `[NEW:qa-archive-gate]` Archiving refuses a specification whose stored QA status is FAILED, listing the failed criteria; the override for pending tasks never bypasses this refusal. A specification without a stored report archives with a visible warning.
- `[NEW:red-hook-scope]` Finishing a batch in which every task is marked RED in its title skips only the per-task verification hooks and prints a notice naming each skipped hook; phase-level and all-tasks hooks still run, and any batch containing an unmarked task runs all hooks as before.
### [NEW:us-qa-persist] Record and Expose the QA Report
- **Scenario:** Recording a valid report
  - **GIVEN** a QA matrix with valid verdicts for an active specification
  - **WHEN** the orchestrator records it
  - **THEN** the metadata stores criteria, findings, derived status, run counter and time, and other metadata is unchanged.
- **Scenario:** Recording an invalid report
  - **GIVEN** a report with an unknown verdict or no criteria
  - **WHEN** recording is attempted
  - **THEN** it is refused naming the problem and the stored report stays untouched.
- **Scenario:** Re-recording after a correction
  - **GIVEN** a stored report
  - **WHEN** a new report is recorded
  - **THEN** it replaces the previous one and the run counter increases.
- **Scenario:** Reading status
  - **GIVEN** a specification with and without a stored report
  - **WHEN** implementation status is requested
  - **THEN** the stored one shows its QA state and the other shows none.
### [NEW:us-qa-archive-gate] Archive Gate on QA State
- **Scenario:** FAILED report blocks archival
  - **GIVEN** a stored FAILED report
  - **WHEN** archiving is requested, with or without the override for pending tasks
  - **THEN** archiving is refused, failed criteria are listed and the specification stays active.
- **Scenario:** Missing report warns
  - **GIVEN** no stored report
  - **WHEN** archiving is requested
  - **THEN** a warning is shown and archiving continues.
### [NEW:us-red-hook-scope] Scoped Hook Relaxation for RED Tasks
- **Scenario:** Finishing a RED task
  - **GIVEN** a task marked RED in its title and configured per-task hooks
  - **WHEN** the task is finished
  - **THEN** per-task hooks are skipped with a notice naming them and the task is finished.
- **Scenario:** Finishing a GREEN or unmarked task
  - **GIVEN** a task not marked RED
  - **WHEN** it is finished
  - **THEN** all configured hooks run and a failure blocks it.
- **Scenario:** RED task closing a phase
  - **GIVEN** a RED task that is the last in its phase
  - **WHEN** it is finished
  - **THEN** the phase-level hooks still run.
