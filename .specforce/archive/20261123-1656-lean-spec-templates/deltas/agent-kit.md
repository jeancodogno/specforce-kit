# Delta: agent-kit
## ADDED
- `[NEW:review-brief]` The feature and bug requirements blueprints MUST open the requirements with a Review Brief of about 8 lines holding what changes, the decisions taken and what to review closely, each with at most 3 short bullets; a feature with no decision beyond the request states so in one line.
- `[NEW:compact-scenarios]` The requirements blueprints MUST write each requirement as a one-line user story followed by lettered one-line scenarios in the form "situation → result", citable as the requirement identifier plus the letter; scenarios with the same result are merged, and a failure scenario is written only when a real failure case exists.
- `[NEW:general-rules]` The feature requirements blueprint MUST state feature-wide domain, interface and quality rules once in a single General Rules section, MUST omit the interface rules for features without a user interface, MUST allow a per-requirement constraint line only for a deviation from the general rules, MUST NOT require a per-requirement quality tag, and MUST state a quality line only for a real, measurable target of the feature.
- `[NEW:task-layout]` The tasks blueprint MUST direct task titles to start with an action verb (Create, Modify, Remove, Move or Rename) followed by the target and the outcome, phase names to state the outcome they deliver, task fields to be separate paragraphs indented under the task, consecutive tasks of a phase to be separated by a horizontal rule preceded by a blank line, and requirement links to cite the scenario letters they cover; it MUST forbid bold task titles.
- `[NEW:execution-strategy-anchors]` The tasks blueprint MUST direct the execution strategy to hold the technical anchors and the decisions that apply to several tasks, stated once and not repeated in tasks, MUST NOT require a narrative of the task order, and MUST direct pre-emptive mitigations to be written only for risks not already handled by a task constraint.
- `[NEW:test-first-neutral]` The artifact blueprints MUST NOT split feature tasks into test and code pairs; when project instructions ask for test-first pairs, the test task title keeps the action verb and ends with the `[RED]` token.
### [NEW:us-lean-reviewable-specs] Lean, Reviewable Specification Blueprints
- **Scenario:** Writing requirements
  - **GIVEN** a planning session for a feature
  - **WHEN** the requirements are written
  - **THEN** they open with a Review Brief, each requirement has a one-line story and lettered one-line scenarios, and feature-wide rules appear once in General Rules.
- **Scenario:** Writing a roadmap
  - **GIVEN** a planning session with requirements
  - **WHEN** the roadmap is written
  - **THEN** task titles start with an action verb, fields are indented paragraphs, consecutive tasks are separated by a horizontal rule, and requirement links cite scenario letters.
- **Scenario:** Project without test-first instructions
  - **GIVEN** project instructions that do not ask for test-first pairs
  - **WHEN** a feature roadmap is written
  - **THEN** tasks are not split into test and code pairs.
## MODIFIED
- `[BR-KIT-35]` Planning blueprints MUST direct task roadmaps to state constraints only when they carry something the implementer cannot deduce from the design and the code (a pitfall, a deviation from repository convention, or a recorded design decision referenced rather than repeated), each stating what goes wrong if it is ignored; MUST exclude from constraints anything that repeats a requirement, repeats another task's constraint, describes the obvious implementation or defines a data contract; MUST forbid method signatures, constants and code snippets in tasks; and MUST direct the design to be the only artifact where code-symbol anchors are written.
### [US-KIT-25] RED Tasks Are Marked in Planning
- **Scenario:** Planning a test-first pair
  - **GIVEN** a roadmap with RED and GREEN task pairs
  - **WHEN** the roadmap is written
  - **THEN** every RED task title starts with its action verb and ends with the RED token.
## REMOVED
