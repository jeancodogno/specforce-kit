---
slug: 20260527-2259-standardize-acceptance-check
---

# Tasks: Standardize Acceptance Check

Implementation roadmap for migrating from `**Acceptance Check:**` to the methodology-agnostic `**Acceptance Check:**` label.

## 2. Tasks

### Phase 1: Core Parser & Logic Implementation
Focus on updating the Go internal logic to support the new label and handle deprecation gracefully.

- [x] T1.1: Update `tasks.go` parser state machine
**Target:** `src/internal/spec/tasks.go`
**Context:** [US-1]

**Action Steps:**
- Update `updateTaskBlockState` to recognize `**Acceptance Check:**` (and legacy label for now).
- Update `validateLastTask` to enforce the new label and flag legacy as error.
- Update `findTaskBlock` and other status update logic.

**Acceptance Check:**
Unit tests in `tasks_validation_test.go` pass for both valid new label and invalid old label.

- [x] T1.2: Refactor `implementation.go` data extraction
**Target:** `src/internal/spec/implementation.go`
**Context:** [US-1]

**Action Steps:**
- Update the extraction regex to capture content under `**Acceptance Check:**`.
- Ensure the JSON output for `specforce implementation status` correctly maps the new field.

**Acceptance Check:**
`specforce implementation status --json` correctly populates the `Verification` field.

- [x] T1.3: Update `status.go` validation UI
**Target:** `src/internal/spec/status.go`
**Context:** [US-1]

**Action Steps:**
- Update the `validationGuide` string to show the correct example format using the new label.
- Update UI helper text in the status command output to mention "Acceptance Check".

**Acceptance Check:**
Error messages when validation fails show the new "Acceptance Check" example.

- [x] T1.4: Verify parser logic via TDD (Red-Green)
**Target:** `src/internal/spec/tasks_test.go`
**Context:** [US-1]

**Action Steps:**
- Add test cases for case-insensitivity and whitespace variations of the new label.
- Add regression tests for the legacy label ensuring it triggers a validation error.

**Acceptance Check:**
`go test ./src/internal/spec/...` returns 100% success for task parsing.

### Phase 2: Artifacts, Templates & Agent Instructions
Update all templates and instruction files to ensure all future specs use the standardized label.

- [x] T2.1: Update core Markdown templates
**Target:** `src/internal/agent/artifacts/spec/tasks.yaml`
**Context:** [US-1]

**Action Steps:**
- Modify `tasks.yaml` and `bug-design.yaml` in `src/internal/agent/artifacts/spec/`.
- Ensure all template examples use `**Acceptance Check:**` instead of `**Acceptance Check:**`.

**Acceptance Check:**
`grep -r "**Acceptance Check:**" src/internal/agent/artifacts/` returns zero matches.

- [x] T2.2: Update Agent & Command instructions
**Target:** `src/internal/agent/kit/commands/implement.yaml`
**Context:** [US-1]

**Action Steps:**
- Update `specforce-planner.yaml` and `implement.yaml` command definition.
- Refactor prompt instructions to refer to "Acceptance Check" for task verification.

**Acceptance Check:**
Agent prompts generated via `specforce` use the "Acceptance Check" terminology.

- [x] T2.3: Update Skill definitions
**Target:** `.gemini/skills/task-atomic-decomposition/SKILL.md`
**Context:** [US-1]

**Action Steps:**
- Update `task-atomic-decomposition/SKILL.yaml` to enforce the new label in its instructions.
- Ensure the TDD skill itself refers to task verification as the "Acceptance Check" phase.

**Acceptance Check:**
Skill verification instructions match the new standard.

- [x] T2.4: Update hardcoded `AGENTS.md` logic
**Target:** `src/internal/project/agents_md.go`
**Context:** [US-1]

**Action Steps:**
- Update `src/internal/project/agents_md.go` which handles the `SPECFORCE_AGENTS` marker content.
- Build and run the generation tool to refresh the project's `AGENTS.md`.

**Acceptance Check:**
`make build` and local execution shows correct instructions in `AGENTS.md`.

### Phase 3: Global Repository Migration
Automated migration of existing specifications to maintain consistency across the project history.

- [x] T3.1: Execute batch migration script
**Target:** `Global Scope`
**Context:** [US-2]

**Action Steps:**
- Run `find .specforce -type f -name "*.md" -exec sed -i 's/\*\*Verification (TDD):\*\*/\*\*Acceptance Check:\*\*/g' {} +`.
- Manually check a few files in `.specforce/archive/` to confirm the replacement worked.

**Acceptance Check:**
`grep -r "**Acceptance Check:**" .specforce/` returns zero matches.

- [x] T3.2: Verify migration integrity
**Target:** `Global Scope`
**Context:** [US-2]

**Action Steps:**
- Run `specforce spec status --all` to ensure all archived and active specs remain valid.
- Fix any manual structural issues that might have been hidden by previous parser leniency.

**Acceptance Check:**
All specs report `IsValid: true`.

### Phase 4: Final Validation & UX Verification
Final checks to ensure the transition is seamless and provides the expected developer experience.

- [x] T4.1: E2E New Spec Initialization
**Target:** `CLI Command`
**Context:** [US-1]

**Action Steps:**
- Run `specforce spec init` for a dummy feature and verify the generated `tasks.md`.
- Confirm the new `tasks.md` passes validation out-of-the-box.

**Acceptance Check:**
The generated `tasks.md` contains `**Acceptance Check:**` for all tasks.

- [x] T4.2: Verify Deprecation UX
**Target:** `CLI Command`
**Context:** [US-1]

**Action Steps:**
- Manually change a label back to the old one and run `specforce status`.
- Ensure the error message is helpful and points the user to the migration fix.

**Acceptance Check:**
CLI outputs a clear error in red mentioning the deprecated label and suggesting the new one.
