# Delta: agent-kit
## ADDED
- `[NEW:lean-planning]` Planning blueprints MUST direct task roadmaps to state constraints only when they carry something the implementer cannot deduce from the design and the code (a pitfall, a deviation from repository convention, or a recorded design decision referenced rather than repeated), MUST forbid method signatures, constants and code snippets in tasks, and MUST direct the design to be the only artifact where code-symbol anchors are written.
- `[NEW:design-why]` The design blueprint MUST allow a one-line `Why:` on a non-obvious technical decision inside the contract or inventory entry it affects, MUST forbid debating alternatives and adding a dedicated rationale section, and MUST NOT require it for decisions with no real alternative.
- `[NEW:brief-design-excerpt]` Implementation orchestration blueprints MUST direct the orchestrator to put into each worker brief the task's constraints (omitting the line when there are none) and, copied verbatim and naming their source, the design inventory entries, related contracts and `Why:` lines that concern the task's target, or to state that none applies when the design is absent or has no matching entry.
### [NEW:us-lean-blueprints] Lean Planning and Design-Anchored Worker Briefs
- **Scenario:** Writing a roadmap
  - **GIVEN** a planning session with a design that records a decision
  - **WHEN** the roadmap is written
  - **THEN** tasks reference the decision instead of repeating it and contain no signatures, constants or code snippets.
- **Scenario:** Preparing a worker brief
  - **GIVEN** a task whose target has an inventory entry with a rationale line
  - **WHEN** the orchestrator prepares the brief
  - **THEN** the brief contains the entry and rationale verbatim with their source and the task's constraints.
- **Scenario:** No design entry
  - **GIVEN** a task whose target has no matching design entry
  - **WHEN** the orchestrator prepares the brief
  - **THEN** the brief states that no design excerpt applies and invents none.
