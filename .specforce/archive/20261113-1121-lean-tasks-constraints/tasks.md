---
slug: 20261113-1121-lean-tasks-constraints
lens: Backend-heavy
---

# Implementation Roadmap: Lean Tasks, Design-Anchored Decisions

## 1. Execution Strategy
- **Gravity Order:** Task validation core -> status and implementation reports and CLI rendering -> planning/design/implement blueprints -> constitution amendment (owner approval) -> black-box verification. Every phase pairs [RED] tests with their [GREEN] change. The Action Steps lines below stay short only because the current validator still demands them; they are the last use of that rule.

## 2. Tasks

### Phase 1: Task Validation Core

- [x] T1.1: [RED] Tests for lean tasks, legacy step lists and pseudo-code warnings
**Target:** `src/internal/spec/tasks_validation_test.go`
**Context:** [US-1], [US-2], [US-3]

**Action Steps:**
- Replace the expectations that demand a step list and a minimum item count with cases: no step list accepted, constraints heading accepted, legacy heading accepted, empty block accepted, mixed headings accepted.
- Keep the rejection cases for a missing target, requirement link (medium and larger) and acceptance check.
- Add warning cases: fenced code block flagged, four distinct paths flagged, three not flagged, and warnings absent from the error list.

**Acceptance Check:**
Run `go test ./src/internal/spec/ -run 'TestValidateTasks'` and observe the new cases failing while the old density rule is still active.

- [x] T1.2: [GREEN] Make step lists optional and collect warnings
**Target:** `src/internal/spec/tasks.go`
**Context:** [US-1], [US-2], [US-3]

**Action Steps:**
- Treat the constraints heading and the legacy step-list heading as one optional list and remove the missing-header and density errors.
- Collect the two warning signals per task in the same pass and expose them through a report that carries errors and warnings, keeping the existing error-only entry point as a wrapper.
- Follow the grammar and path-counting rule recorded in the design contract.

**Acceptance Check:**
Run `go test ./src/internal/spec/ -run 'TestValidateTasks'` and verify all cases pass.

### Phase 2: Status, Implementation Report and CLI Rendering

- [x] T2.1: [RED] Tests for warnings in the specification status
**Target:** `src/internal/spec/status_test.go`
**Context:** [US-3]
**Parallel With:** T2.2, T2.3

**Action Steps:**
- Assert a flagged task yields a warnings list on the tasks artifact while `is_valid` and `progress` stay unchanged.
- Assert the list is omitted when there are no warnings and that the validation guide no longer shows a step list.

**Acceptance Check:**
Run `go test ./src/internal/spec/ -run 'TestGetStatus'` and observe the new assertions failing.

- [x] T2.2: [RED] Tests for constraints in the implementation report
**Target:** `src/internal/spec/implementation_test.go`
**Context:** [US-1], [US-2]
**Parallel With:** T2.1, T2.3

**Action Steps:**
- Assert constraints are filled from either heading and are empty when a task has none.
- Assert a task without a step list never makes the report fail.

**Acceptance Check:**
Run `go test ./src/internal/spec/ -run 'TestImplementation'` and observe the new assertions failing.

- [x] T2.3: [RED] Tests for CLI output of warnings and constraints
**Target:** `src/internal/cli/integration_test.go`
**Context:** [US-1], [US-3]
**Parallel With:** T2.1, T2.2

**Action Steps:**
- Assert the human status output prints warnings after errors and marked as warnings, and that the exit outcome does not change.
- Assert the machine-readable outputs expose the warnings list and the constraints field.

**Acceptance Check:**
Run `go test ./src/internal/cli/ -run 'Integration'` and observe the new assertions failing.

- [x] T2.4: [GREEN] Expose warnings in the specification status
**Target:** `src/internal/spec/status.go`
**Context:** [US-3]
**Parallel With:** T2.5

**Action Steps:**
- Carry the warnings of the tasks report into the artifact status as an optional list that never affects validity or progress.
- Update the validation guide example to the lean form.

**Acceptance Check:**
Run `go test ./src/internal/spec/ -run 'TestGetStatus'` and verify it passes.

- [x] T2.5: [GREEN] Expose constraints in the implementation report
**Target:** `src/internal/spec/implementation.go`
**Context:** [US-1], [US-2]
**Parallel With:** T2.4

**Action Steps:**
- Replace the step-list field with constraints, filled from either heading and empty when absent.
- Keep the report building when a task has no constraint block.

**Acceptance Check:**
Run `go test ./src/internal/spec/ -run 'TestImplementation'` and verify it passes.

- [x] T2.6: [GREEN] Render warnings in the CLI
**Target:** `src/internal/tui/spec.go`
**Context:** [US-3]

**Action Steps:**
- Print the warnings list after the errors, visually distinct, in the human output.
- Leave the exit outcome and the validity flag untouched.

**Acceptance Check:**
Run `go test ./src/internal/cli/` and verify it passes.

### Phase 3: Blueprints (Planning, Design, Implementation Brief)

- [x] T3.1: [RED] Manifest tests for the lean blueprints
**Target:** `src/internal/agent/kit_manifests_test.go`
**Context:** [US-1], [US-4], [US-5]

**Action Steps:**
- Assert the tasks blueprint no longer mandates a step density, forbids signatures, constants and snippets, and keeps the acceptance-check rule.
- Assert the design blueprint allows a one-line `Why:` without a new section and keeps the ban on debating alternatives.
- Assert the implement brief lists constraints (omitted when empty) and carries the verbatim design-excerpt directive with the no-entry fallback.

**Acceptance Check:**
Run `go test ./src/internal/agent/ -run 'Manifest'` and observe the new assertions failing.

- [x] T3.2: [GREEN] Lean task blueprint
**Target:** `src/internal/agent/artifacts/spec/tasks.yaml`
**Context:** [US-1]
**Parallel With:** T3.3, T3.4

**Action Steps:**
- Replace the directive and density rules with the lean rules: constraints optional and only for non-deducible items, design decisions referenced and not repeated, no signatures, constants or snippets.
- Reword the determinism rule to no placeholders and no scope ambiguity, and update the template and the verification-phase example.

**Acceptance Check:**
Run `go test ./src/internal/agent/ -run 'Manifest'` and verify the tasks assertions pass.

- [x] T3.3: [GREEN] Design blueprint with inline rationale
**Target:** `src/internal/agent/artifacts/spec/design.yaml`
**Context:** [US-5]
**Parallel With:** T3.2, T3.4

**Action Steps:**
- Relax the rule against rationale to allow one `Why:` line per non-obvious decision inside contracts and inventory entries, still forbidding alternatives debate and new sections.
- State that the design is the only place for code-symbol anchors and add a `Why:` example line to the template.

**Acceptance Check:**
Run `go test ./src/internal/agent/ -run 'Manifest'` and verify the design assertions pass.

- [x] T3.4: [GREEN] Worker brief with constraints and design excerpt
**Target:** `src/internal/agent/kit/skills/spf-implement/SKILL.yaml`
**Context:** [US-4]
**Parallel With:** T3.2, T3.3

**Action Steps:**
- Replace the step-list line of the brief with the constraints line, omitted when empty.
- Add the directive to copy verbatim, with source, the design entries for the task's target, or to state that none applies.

**Acceptance Check:**
Run `go test ./src/internal/agent/ -run 'Manifest'` and verify the brief assertions pass.

- [x] T3.5: [GREEN] Planning skill alignment
**Target:** `src/internal/agent/kit/skills/spf-spec/SKILL.yaml`
**Context:** [US-1], [US-5]

**Action Steps:**
- Remove wording that implies a step density and point planning to the lean rules and the `Why:` guidance.
- Keep the delta step and the quality gate unchanged.

**Acceptance Check:**
Run `go test ./src/internal/agent/` and verify it passes.

- [x] T3.6: [GREEN] Refresh installed skill copies
**Target:** `.claude/skills`
**Context:** [US-1], [US-4]

**Action Steps:**
- Regenerate the installed copies of the changed skills with the repository's existing tool sync, never by hand.
- Confirm only the planning and implementation skills changed.

**Acceptance Check:**
Run `grep -c "Constraints" .claude/skills/spf-implement/SKILL.md` and expect a count of at least 1.

### Phase 4: Constitution Amendment

- [x] T4.1: [GREEN] Amend the self-containment rule after owner approval
**Target:** `.specforce/docs/engineering.md`
**Context:** [US-6]

**Action Steps:**
- Present the proposed wording to the owner and wait for explicit approval; do not edit before it.
- After approval, replace the self-containment rule so it requires self-containment in business rules, intent and decisions and limits symbol-level anchors to the design and to non-obvious constraints.

**Acceptance Check:**
Run `grep -c "symbol-level" .specforce/docs/engineering.md` and expect the rule to mention symbol-level anchors only in the design-scoped wording, with the owner's approval recorded in the session.

### Phase 5: Black-box Verification

- [x] T5.1: [VERIFY] Lean roadmap accepted by the compiled CLI
**Target:** `make build && ./specforce spec status 20261113-1121-lean-tasks-constraints --json`
**Context:** [US-1], [US-2]

**Action Steps:**
- Build the binary and run the status on this very specification, whose tasks use the previous step-list heading.
- Compare against the expected result: `is_valid` true and no validation errors.
- Record the output.

**Acceptance Check:**
`./specforce spec status 20261113-1121-lean-tasks-constraints --json` exits 0 with `"is_valid": true`.

- [x] T5.2: [VERIFY] Warnings visible and non-blocking in the compiled CLI
**Target:** `./specforce spec status <fixture-slug> --json`
**Context:** [US-3]

**Action Steps:**
- Create a disposable fixture specification under the scratch area whose tasks contain a fenced code block and cite four paths.
- Run the status and compare against the expected result: two warnings listed, `is_valid` true, exit 0.
- Record the output and remove the fixture.

**Acceptance Check:**
`./specforce spec status <fixture-slug> --json` exits 0, lists the warnings and keeps `"is_valid": true`.

- [x] T5.3: [VERIFY] Constraints in the implementation status and design excerpt in the brief
**Target:** `./specforce implementation status <fixture-slug> --json`
**Context:** [US-4], [US-5]

**Action Steps:**
- Run the implementation status on the fixture with one task having constraints and one without and compare the constraints fields.
- Read the changed implement blueprint and confirm the verbatim design-excerpt directive and its fallback are present.
- Record the output.

**Acceptance Check:**
`./specforce implementation status <fixture-slug> --json` exits 0 with a populated constraints array for the first task and an empty one for the second.

- [x] T5.4: [VERIFY] Constitution amendment gated by approval
**Target:** `Global Scope`
**Context:** [US-6]

**Action Steps:**
- Confirm from the session record that the owner approved the wording before the document changed.
- Run the full suite and compare against the expected result: all green.

**Acceptance Check:**
`go test ./...` exits 0.

## 3. Pre-emptive Mitigations
- **Risk:** Removing the density rule leaves old validator tests and guides assuming a step list -> **Mitigation:** T1.1, T2.1 and T2.3 rewrite those expectations before any production change.
- **Risk:** The constitution is changed without approval -> **Mitigation:** T4.1 stops for explicit owner approval and T5.4 checks the record.
