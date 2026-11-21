# Delta: agent-kit
## ADDED
- `[NEW:qa-runs-scenarios]` The independent qa-tester's black-box tier MUST execute every scenario of the roadmap's black-box scenario section and add one traceability matrix row per scenario, identified by its scenario identifier, using the existing verdicts; a scenario that cannot be executed is `NOT-VERIFIABLE` with a stated reason, and roadmaps without the section derive black-box checks from the requirements' acceptance scenarios.
### [NEW:us-qa-scenarios] Independent QA Runs Declared Scenarios Once
- **Scenario:** Roadmap with scenarios
  - **GIVEN** a roadmap that declares black-box scenarios
  - **WHEN** the qa-tester runs the black-box tier
  - **THEN** every scenario command is executed and recorded as its own matrix row next to requirement and delta rows.
- **Scenario:** Roadmap without scenarios
  - **GIVEN** a roadmap with no scenario section
  - **WHEN** the qa-tester runs the black-box tier
  - **THEN** black-box checks are derived from the requirements' acceptance scenarios.
- **Scenario:** Scenario cannot run
  - **GIVEN** a scenario whose command cannot be executed in the environment
  - **WHEN** the matrix is built
  - **THEN** its row is `NOT-VERIFIABLE` with a reason and is never reported as passed.
## MODIFIED
- `[BR-KIT-20]` Planning blueprints MUST direct every task acceptance check to be a single non-interactive command or test scenario with an observable expected result, rewriting or explicitly flagging criteria that need human judgment as not verifiable by a tester. Every task MUST change at least one product file or artifact; tasks whose only purpose is running the full test suite, a smoke run, a build check or repository-wide searches are forbidden, since that verification belongs to the independent QA step and the all-tasks hooks. Black-box scenarios are declared in a non-executable scenario section placed after the task phases (never as tasks, checkboxes or task identifiers), each linked to at least one requirement and stating one non-interactive command and its expected observable result; the section is mandatory for `large` and `complex` roadmaps, optional for `medium` roadmaps that produce a runnable artifact, and omitted for `small` roadmaps. Existing roadmaps with a verification phase of tasks remain valid and are executed as ordinary tasks.
### [US-KIT-17] Executable Acceptance Checks and Black-box Scenarios
- **Scenario:** Planning a large specification
  - **GIVEN** a `large` or `complex` specification being planned
  - **WHEN** the roadmap is produced
  - **THEN** every task has a non-interactive acceptance check, no task exists only to run verification, and a black-box scenario section traced to requirements follows the task phases.
- **Scenario:** Planning a small specification
  - **GIVEN** a `small` specification being planned
  - **WHEN** the roadmap is produced
  - **THEN** it has no scenario section and no verification-only task.
- **Scenario:** Criterion needing human judgment
  - **GIVEN** a task whose outcome cannot be scripted
  - **WHEN** its acceptance check is written
  - **THEN** it is rewritten into an observable check or flagged as not verifiable by a tester.
## REMOVED
