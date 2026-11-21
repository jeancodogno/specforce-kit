# Delta: spec-management
## ADDED
- `[NEW:task-warnings]` Task validation produces non-blocking warnings, listed apart from errors in the human and machine-readable specification status, for a task that contains a fenced code block and for a task whose target and constraints cite more than three distinct file paths; warnings never change validity, progress or the outcome of any command.
- `[NEW:task-constraints]` The step list of a task is optional: a task may omit it, may use the constraints heading or the legacy step-list heading, and may leave the block empty, with no minimum or maximum number of items; the implementation status exposes the items of either heading as the task's constraints, empty when absent, and every roadmap valid before this change stays valid.
### [NEW:us-lean-tasks] Lean Task Roadmaps
- **Scenario:** Task without a step list
  - **GIVEN** a task with a target, a requirement link and an acceptance check but no step or constraint list
  - **WHEN** the roadmap is validated
  - **THEN** it is accepted and no message about constraints is produced.
- **Scenario:** Legacy step list
  - **GIVEN** a roadmap whose tasks use the legacy step-list heading
  - **WHEN** it is validated and its tasks are read for implementation
  - **THEN** it is accepted and the items are delivered as the task's constraints.
- **Scenario:** Pseudo-code drift
  - **GIVEN** a task containing a fenced code block or citing more than three distinct file paths
  - **WHEN** the specification status is requested
  - **THEN** a warning names the task and the reason while the specification is still reported as valid.
## MODIFIED
- `[BR-SPEC-04]` Task verification (`ValidateTasks`) and coherence auditing (`DeterministicCheck`) adjust strictness based on size:
  - `small`: optional `**Context:** [US-X]` mapping. Missing `requirements.md` is not flagged as an error.
  - `medium`, `large`, `complex`: mandatory `**Context:** [US-X]` and deterministic requirement-to-task coherence.
  - Every size: a task requires a target and an acceptance check, and enforces no minimum number of action steps or constraints.
### [US-SPEC-03] Adaptive Task and Coherence Auditing
- **Scenario:** Verifying task density and requirement mappings
  - **GIVEN** a `small` specification with simple tasks
  - **WHEN** task validation and coherence audit run
  - **THEN** the validation suite passes without forcing multi-step decomposition or artificial requirement tags.
- **Scenario:** Verifying lean tasks of a larger specification
  - **GIVEN** a `medium`, `large` or `complex` specification whose tasks carry no step list
  - **WHEN** task validation and coherence audit run
  - **THEN** tasks are accepted as long as each has a target, a requirement link and an acceptance check.
