[ [English](cli.md) | [Português](../pt/cli.md) | [Español](../es/cli.md) ]

# Command Line Interface (CLI)

This page lists every command and flag of the `specforce` CLI, grouped by area. Many of these commands are run by the AI agent while it follows a skill (see [Invoking skills](supported-tools.md#invoking-skills)); they are documented here so you know what exists and what each one does.

Commands that take a `<slug>` accept the spec folder name (for example `20261121-0035-csv-export`) or just its readable part (`csv-export`). Flags marked `--json` print machine-readable JSON intended for automation and agents.

## Global

### `specforce init [agents...]`
Initializes Specforce in the current directory: creates the `.specforce/` folder (with a default `config.yaml`), the root `AGENTS.md` managed block and the skills of the selected tools. Without arguments, an interactive picker lets you choose the tools. With arguments, the given tool keys (`claude`, `qwen`, `kimi-code`, `open-code`, `kilo-code`, `codex`, `antigravity`, `cursor`) are configured immediately. Run it again in an initialized project to add or remove tools and refresh the installed skills. With arguments, every installed tool you leave out is removed by deleting its whole folder (for example `.claude/`, including your own settings), without confirmation; the interactive picker warns before deleting. See [Supported Tools](supported-tools.md).

### `specforce console`
Launches the **Specforce Console**, an interactive terminal interface to browse specifications, the Constitution and implementation tasks, and to follow progress in real time (including specs from linked [Git worktrees](git-worktrees.md)). Leave it running in a separate terminal pane while your agent works.

---

## Global flags & environment

- **`-v`, `--version`**: prints the installed version (`specforce version <x.y.z>`).
- **`-h`, `--help`**: shows the help of any command or group.
- **`--dev`**: uses the local `kit/` folder instead of the kit embedded in the binary. Intended for kit contributors only.
- **`DEBUG=1`**: environment variable that writes diagnostics to standard error, including every auto-update decision (see [Updates](getting-started.md#updates)).

### `specforce completion <shell>`
Generates the shell autocompletion script for `bash`, `zsh`, `fish` or `powershell`. Run `specforce completion <shell> --help` for the installation steps of each shell.

---

## Constitution

### `specforce constitution status [--json]`
Shows the completeness status of the project Constitution documents in `.specforce/docs/` (principles, architecture, ui-ux, security, engineering, governance and modules).

### `specforce constitution artifact <name> [--json]`
Shows the details of one Constitution artifact: its purpose, template and instructions. Agents use it to fetch the right rules when they write or update that document.

---

## Specification

### `specforce spec init <slug> [--type feature|bug] [--size small|medium|large|complex] [--json]`
Creates a new specification folder `.specforce/specs/<timestamp>-<slug>/` with its `spec.yaml` metadata.
- `--type`: `feature` (default) or `bug`. Bug specs use the bug variants of requirements and design.
- `--size`: `small`, `medium` (default), `large` or `complex`. The size decides which artifacts are required (see [Artifacts](artifacts.md)).

### `specforce spec list [--json]`
Lists all active specifications.

### `specforce spec status <slug> [--json]`
Shows the completeness of a specification: which artifacts its size requires, which exist, validation errors and non-blocking task warnings.

### `specforce spec artifact <name> [--json]`
Shows the instructions and template of a spec artifact (for example `requirements`, `design`, `tasks`, `bug-requirements`, `bug-design`), including the project instructions from [`config.yaml`](configuration.md). Agents use it before writing each artifact.

### `specforce spec resize <slug> --size <size> [--json]`
Changes the size of an existing specification. Promoting it marks newly required artifacts as missing and blocks implementation until they exist; demoting it relaxes validation and keeps existing files on disk.

### `specforce spec audit <slug> [flags]`
Updates the refinement audit state stored in `spec.yaml`, used by the planning skill while it checks coherence between artifacts.
- `--valid` / `--invalid`: marks the spec as valid or invalid.
- `--error <msg>`: appends a coherence error message.
- `--clear`: clears all coherence errors.
- `--iteration <n>`: sets the current refinement iteration count.
- `--json`: machine-readable output.

### `specforce spec archive <slug> [--force]`
Finalizes a completed specification and moves its folder to `.specforce/archive/`. It refuses when:
- tasks are still pending;
- a `medium`, `large` or `complex` spec has a delta that is not `applied` (or declared `none`); `small` specs report `exempt`;
- the stored QA report is `FAILED` (the failed criteria are listed).

A spec without a QA report archives with a visible warning. `--force` overrides only the pending-tasks check; it never bypasses the delta check nor a `FAILED` QA verdict. A refusal leaves the filesystem untouched.

---

## Delta

The delta describes the changes a specification makes to the module living specs in `.specforce/docs/modules/`. Archiving requires the delta to be applied, none or exempt.

### `specforce spec delta check <slug> [--json]`
Reports whether the spec's delta (the proposed changes to the module living specs) has been applied. Read-only. The report state is one of:
- `none`: the spec declares `delta.impact: none` (no module impact).
- `exempt`: a `small` spec with no delta; it needs none.
- `missing`: a delta is required but has not been written.
- `pending`: a delta exists but its changes are not yet reflected in the modules.
- `applied`: the modules already match the delta.
- `invalid`: the delta is malformed or conflicts with the current modules.

### `specforce spec delta seal <slug> [--json]`
Seals the spec's delta by recording the fingerprints of the module blocks it is based on, so later hand edits to those blocks are detected as base conflicts.

---

## Implementation & QA

### `specforce implementation status <slug> [--json]`
Shows the implementation roadmap of a specification (from `tasks.md`): readiness, missing artifacts for its size, each task with its state and constraints, the roadmap mitigations, the project `implementation` instructions and, when one exists, the stored QA state.

### `specforce implementation update <slug> --task <id>[,<id>...] --status <status>`
Updates the state of one or more tasks. `--task` accepts comma-separated IDs or can be repeated (`--task T1.1 --task T1.2`). `--status` is `pending`, `in-progress` or `finished`.

When the status is `finished`, Specforce first runs the configured [validation hooks](configuration.md#3-validation-hooks-quality-gates), once per batch and without duplicates: `on_task_finished`, plus `on_phase_finished` when a task is the last of its phase and `on_all_tasks_finished` when it is the last of the roadmap. A failing hook blocks the update. When every task of the batch carries the `[RED]` token in its title, only the `on_task_finished` hooks are skipped, with a notice naming them; phase and all-tasks hooks still run.

### `specforce implementation qa record <slug>`
Records the independent QA report of a specification. It reads the report as JSON on standard input (criteria with scenario, evidence and verdict `PASSED`, `FAILED` or `NOT-VERIFIABLE`, plus out-of-scope findings), validates it, derives the overall status (`FAILED` when any criterion failed, otherwise `PASSED`) and replaces the previous report in `spec.yaml`. Invalid input is refused without changing the stored report. The implementation skill uses it after the QA step.

---

## Archive

### `specforce archive instructions`
Shows the full archiving protocol: the description of every Constitution document, the core archiving rules (applying the delta to the module living specs, updating the Constitution) and the project `archive` instructions from `config.yaml`. The archiving skill follows it before running `specforce spec archive`.
