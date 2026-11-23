# Delta: spec-management
## ADDED
- `[NEW:audit-all-links]` The coverage audit MUST treat a requirement as covered when any task's requirement link cites it, whatever its position in the link and whether scenario letters follow it. A requirement link is a line that starts with the Context label, ignoring indentation; mentions anywhere else (titles, constraints, bulleted or inline Context labels) do not count, and requirements listed under Out of Scope stay excluded.
### [NEW:us-audit-multi-links] Coverage Audit Recognizes Every Linked Requirement
- **Scenario:** Task covering several requirements
  - **GIVEN** a task whose requirement link cites US-1 and US-2
  - **WHEN** the coverage audit runs
  - **THEN** neither requirement is reported as missing.
- **Scenario:** Link with scenario letters
  - **GIVEN** a requirement link that cites a requirement followed by scenario letters
  - **WHEN** the coverage audit runs
  - **THEN** the requirement is recognized as covered.
- **Scenario:** Uncovered requirement
  - **GIVEN** a requirement that no requirement link cites, even if a task title or constraint mentions it
  - **WHEN** the coverage audit runs
  - **THEN** it is reported as missing.
## MODIFIED
## REMOVED
