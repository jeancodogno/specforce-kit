[ [English](artifacts.md) | [Português](../pt/artifacts.md) | [Español](../es/artifacts.md) ]

# Specforce Artifacts

To orchestrate AI agents without piling up unnecessary context, Specforce splits the project's knowledge into focused documents called **Artifacts**. This page explains each of them and the lifecycle rules that decide what the agent produces and why a step can be blocked. Skills are referred to by name (for example [`spf.spec`](supported-tools.md#invoking-skills)).

## 1. Constitution (`.specforce/docs/`)

The **Constitution** lives in your project's `.specforce/docs/` folder. It holds the durable rules of the system and is created and maintained with `spf.constitution`:

- **`principles.md`**: core values, philosophy and technical axioms.
- **`architecture.md`**: system boundaries, dependency direction and persistence topology.
- **`ui-ux.md`**: visual direction, interaction patterns and aesthetic rules (when applicable).
- **`security.md`**: authorization, roles, permissions and data protection.
- **`engineering.md`**: coding standards, testing strategy and refactoring guidelines.
- **`governance.md`**: lifecycle rules, ownership and AI boundaries.
- **`modules/<module>.md`**: living specs of each domain (business rules, use cases, public surfaces), updated at archive time from the spec deltas.

The documents are not loaded all at once: the agent reads the ones relevant to the task at hand.

## 2. The specification package (`.specforce/specs/<timestamp>-<slug>/`)

Every change starts with a specification created by `spf.spec`. Its folder can contain:

- **`spec.yaml`**: metadata managed by the CLI: type, size, delta declaration and seal, refinement audit, task time logs and the QA report.
- **`proposal.md`** *(optional)*: the formalized result of a `spf.discovery` session (problem, verified findings, decisions, open questions, suggested type and size). Planning reuses it instead of asking again.
- **`requirements.md`**: the "what" and "why": user stories, BDD acceptance criteria (GIVEN/WHEN/THEN) and business constraints. No technical details.
- **`design.md`**: the "how": architecture, contracts and the file and component inventory. No final code.
- **`tasks.md`**: the atomic execution roadmap (see [Task format](#5-task-format)).
- **`deltas/<module>.md`**: the proposed changes to each affected module living spec (see [Module deltas](#6-module-deltas)).

Which of these files are required depends on the spec size.

## 3. Sizes

Every spec has a size, chosen at creation (`specforce spec init <slug> --size <size>`) and changeable later with `specforce spec resize`. The default size is `medium`.

| Size | Required artifacts | Module delta |
|---|---|---|
| `small` | `tasks.md` | Exempt |
| `medium` (default) | `requirements.md`, `tasks.md` | Required (or `delta.impact: none`) |
| `large` | `requirements.md`, `design.md`, `tasks.md` | Required (or `delta.impact: none`) |
| `complex` | `requirements.md`, `design.md`, `tasks.md` | Required (or `delta.impact: none`) |

- Implementation stays blocked until every artifact required by the size exists.
- Promoting a spec (for example `small` to `large`) marks the new artifacts as missing; demoting it relaxes validation and keeps existing files as reference.
- In `small` specs, tasks may omit the `**Context:** [US-X]` link; from `medium` up, every task must link to a requirement.
- `large` and `complex` roadmaps must include a Black-box Scenarios section (see below).

## 4. Bug specs

A spec created with `--type bug` uses the bug variants of the planning artifacts, from the templates `bug-requirements` and `bug-design`:

- **Requirements (`bug-requirements`)**: issue description, evidence (logs, errors), concrete reproduction steps, an evidence-based root cause analysis and regression acceptance criteria (`[FIX-x]`). When the root cause is not established, the agent stops to ask or recommends `spf.discovery` instead of guessing.
- **Design (`bug-design`)**: the code paths to change, the regression strategy and side effects.
- **Roadmap**: the first task is a regression test that fails before the fix (a `[RED]` task). No fix task may come before it.

The files keep the usual names (`requirements.md`, `design.md`, `tasks.md`).

## 5. Task format

Each task in `tasks.md` is a checklist item such as `- [ ] T1.1: [CODE] Title`, followed by:

- **`**Target:**`**: the exact file path or CLI command the task changes.
- **`**Context:**`**: the requirement it implements (`[US-X]`), mandatory from `medium` up.
- **`**Constraints:**`** *(optional)*: only what the implementer cannot deduce from the design and the code (a pitfall, a deviation from convention, a referenced design decision). No signatures, constants or code snippets.
- **`**Acceptance Check:**`**: a single non-interactive command or test scenario with an observable expected result. Criteria needing human judgment are rewritten or flagged as not verifiable.

Rules:
- Every task changes at least one product file or artifact. Verification-only tasks (run the full suite, a smoke run, a build check) are forbidden: that belongs to the independent QA step and the `on_all_tasks_finished` hooks.
- Black-box checks go in a non-executable `## 3. Black-box Scenarios` section after the task phases, never as tasks. Each scenario (`S1`, `S2`...) links to a requirement and states one command and its expected result. The section is mandatory for `large` and `complex`, optional for `medium` specs that produce a runnable artifact, and omitted for `small`.
- `specforce spec status` shows non-blocking warnings for tasks containing code blocks or citing more than three file paths.

## 6. Module deltas

From `medium` up, a spec must declare how it changes the module living specs in `.specforce/docs/modules/` before implementation:

- **Delta documents**: one `deltas/<module>.md` per affected module, with `## ADDED`, `## MODIFIED` and `## REMOVED` sections carrying complete blocks. New rules and use cases use provisional `NEW:<key>` keys; final identifiers are allocated at archive time.
- **No impact**: when no module changes, the spec declares `delta.impact: none` in `spec.yaml` instead.
- **Seal**: `specforce spec delta seal <slug>` records fingerprints of the blocks the delta modifies or removes. If another spec changes the same block afterwards, the delta check reports a conflict and the delta must be redone against the current module and sealed again.
- **States** reported by `specforce spec delta check <slug>`:
  - `none`: `delta.impact: none` is declared.
  - `exempt`: a `small` spec without delta.
  - `missing`: a delta is required but not written.
  - `pending`: the delta is not yet reflected in the modules.
  - `applied`: the modules match the delta.
  - `invalid`: the delta is malformed or conflicts with the current modules.
- **Archive gate**: `spf.archive` applies the delta to the modules; a spec is archived only when its state is `applied`, `none` or `exempt`.

## 7. QA report

After implementation, `spf.implement` runs an independent QA step: a tester who did not implement any task checks every acceptance criterion (and every delta item and black-box scenario) and gives each one a verdict:

- **`PASSED`**: verified with evidence.
- **`FAILED`**: the criterion is not met. Failed criteria trigger at most two correction cycles before the agent escalates to you.
- **`NOT-VERIFIABLE`**: the criterion cannot be executed by a tester; the reason is stated.

The report is stored in `spec.yaml` (one report per spec, replaced on each new run). Its overall status is `FAILED` when any criterion failed, otherwise `PASSED`. `specforce implementation status` shows it.

Archive gate: a `FAILED` report blocks archiving (even with `--force`); a missing report only produces a warning.

## 8. Test-first (RED) tasks

A task whose title carries the exact token `[RED]` (case-sensitive) is a test-first task: it writes a test that is expected to fail until a later `[GREEN]` task implements the behavior.

- When every task finished in a batch is `[RED]`, Specforce skips only the per-task `on_task_finished` hooks and prints a notice naming them, since a failing test suite is expected at that point.
- `on_phase_finished` and `on_all_tasks_finished` hooks still run, even when a `[RED]` task closes the phase.
- A batch containing any task not marked `[RED]` runs all hooks as usual.

To make roadmaps test-first, add an instruction to `tasks` in `config.yaml` (see [Configuration](configuration.md)).

## 9. The `.specforce/` folder

The root `.specforce/` folder holds the orchestration state. Apart from the configuration file, it is not meant to be edited by hand:

- **`config.yaml`**: project instructions, context variables and validation hooks (see [Configuration](configuration.md)).
- **`docs/`**: the Constitution and module living specs.
- **`specs/`**: active specifications.
- **`archive/`**: specifications closed by `spf.archive`, the history of the project's technical decisions.
- **`scratch/`**: disposable evidence created by `spf.discovery` (tests, scripts, logs, screenshots). It ignores itself in version control and is treated as stale in later sessions.

---

Maintaining the quality of these artifacts is the best way to ensure the AI produces excellent code.
