[ [English](workflows.md) | [Português](../pt/workflows.md) | [Español](../es/workflows.md) ]

# Workflows & Tips

This page explains some Specforce workflows and tips for using it.

Skills are referred to by name (for example [`spf.spec`](supported-tools.md#invoking-skills)). How to call them explicitly in your tool is shown in [Invoking skills](supported-tools.md#invoking-skills); a plain-language request that names the skill works everywhere.

**How to read the workflows**

- 💻 **Terminal:** a command **you** type in your own terminal. Only a handful of commands are ever shown here: `specforce init`, `specforce console`, `specforce spec list` and a couple of `git` commands. Everything else (status, delta checks, QA results) the agent checks and reports to you in chat.
- 💬 **Agent chat:** a request **you** type in your AI agent's chat. Replace the parts in `<angle brackets>`.
- **Behind the scenes** (collapsed): what the agent runs on its own. You never need to run those commands; they are listed so you can follow what happens. All of them are described in the [CLI Reference](cli.md).
- `<slug>` is the spec name, for example `csv-export`. `specforce spec list` shows it.

## Workflows

### 1. Prototype in discovery before implementing

**When to use:** the idea is still vague, or you want proof (a quick test, a script, a screenshot) before committing to a plan.

1. 💬 **Agent chat:**
   ```text
   Use spf.discovery to explore: in our invoicing SaaS, let customers export all their invoices to CSV from the dashboard. Big accounts have 200k invoices, so prove with a disposable test that streaming the export keeps memory flat before we plan anything.
   ```
2. 💬 **Agent chat:** once you agree on the approach:
   ```text
   Formalize this into a proposal.
   ```
3. 💬 **Agent chat:**
   ```text
   Use spf.spec to plan csv-export from the proposal.
   ```

<details><summary>Behind the scenes</summary>

- Discovery never edits, renames or deletes your existing files. Evidence (tests, scripts, logs, screenshots) goes to `.specforce/scratch/`, which ignores itself in version control. If a hypothesis can only be proven inside the source tree, the agent creates a **new** file whose name contains `spf-probe`, records it and removes it before the session ends, checking `git status`.
- Formalizing runs `specforce spec init <slug>` and writes `proposal.md` (problem, findings, decisions taken, rejected alternatives, open questions, suggested type and size, evidence paths).
- `spf.spec` restates the decisions already in the proposal for a single confirmation and asks only the open questions.

</details>

**If it goes wrong:**
- Discovery will not push you to a spec on its own; ask for the proposal explicitly when you are ready.
- A leftover `spf-probe` file shows up in `git status`: ask the agent to remove it.
- Scratch evidence from an earlier session is stale; ask the agent to re-run it before relying on it.

### 2. Fix a bug with a failing test first

**When to use:** something is broken and you want the fix proven by a regression test.

1. 💬 **Agent chat:**
   ```text
   Use spf.discovery to investigate this bug: in the invoices dashboard, filtering by "last month" leaves out the invoices issued on the last day of the month for customers in UTC-3. No error is shown, the list just comes back incomplete.
   ```
2. 💬 **Agent chat:** once the root cause is confirmed:
   ```text
   Formalize this into a bug spec.
   ```
3. 💬 **Agent chat:** review and approve the spec, then:
   ```text
   Use spf.implement to implement invoice-date-filter.
   ```

<details><summary>Behind the scenes</summary>

- The bug spec is created with `specforce spec init <slug> --type bug`. Its requirements hold the reproduction steps, the evidence and the root cause.
- The roadmap starts with a `[RED]` task that adds a regression test failing before the fix; no fix task comes before it.

</details>

**If it goes wrong:**
- `spf.spec` stops when the root cause is not established: go back to step 1 instead of letting it guess.
- The `[RED]` task finishes while the tests fail: that is expected; only the per-task hooks are skipped (see [Hook failure blocks a task](#hook-failure-blocks-a-task)).

### 3. Make a small change

**When to use:** a contained change (a fix in one place, a copy tweak, a config option) that does not need requirements or a design.

1. 💬 **Agent chat:**
   ```text
   Use spf.spec to plan: add a "Due date" column to the invoices table in the dashboard.
   ```
2. 💬 **Agent chat:** review the task list, then:
   ```text
   Use spf.implement to implement invoice-due-date-column.
   ```

<details><summary>Behind the scenes</summary>

- You do not need to say the change is small: the agent sizes it while planning and runs `specforce spec init <slug> --size small`: only `tasks.md` is required and the module delta is `exempt` (see [Sizes](artifacts.md#3-sizes)).

</details>

**If it goes wrong:**
- The change turns out bigger than expected: follow [workflow 4](#4-the-scope-grew-resize-the-spec).

### 4. The scope grew: resize the spec

**When to use:** while planning or implementing, the change turns out to need more (or less) than first planned.

1. 💬 **Agent chat:** just describe the new scope:
   ```text
   Use spf.spec to update invoice-due-date-column: besides the column, email the customer a reminder 3 days before the due date.
   ```
2. 💬 **Agent chat:** after approving the new artifacts:
   ```text
   Use spf.implement to continue invoice-due-date-column.
   ```

<details><summary>Behind the scenes</summary>

- You do not need to ask for a resize: the agent notices the scope changed and runs `specforce spec resize <slug> --size <new size>` on its own. Promoting marks the newly required artifacts as missing and blocks implementation until they exist; demoting relaxes validation and keeps existing files as reference.
- The agent then writes the artifacts the new size requires and, from `medium` up, writes and seals the module delta, or declares `delta.impact: none`.

</details>

**If it goes wrong:**
- Implementation refuses to continue: see [Resize during implementation](#resize-during-implementation).

### 5. Change after implementing, before archiving

**When to use:** the implementation is done (or almost), and you want something different before closing the spec.

1. 💬 **Agent chat:**
   ```text
   Use spf.spec to update invoice-due-date-column: send the reminder 5 days before the due date instead of 3. Keep requirements, design, delta and tasks consistent.
   ```
2. 💬 **Agent chat:** review and approve the updated spec, then:
   ```text
   Use spf.implement to implement only the new or changed tasks of invoice-due-date-column.
   ```

<details><summary>Behind the scenes</summary>

- An approved spec is a contract: the agent updates every affected artifact together, never only `tasks.md`, and adds new tasks instead of rewriting finished ones.
- `spf.implement` runs a fresh independent QA step afterwards and records a new verdict.

</details>

**If it goes wrong:**
- Do not edit the code or `tasks.md` by hand to "just fix it": see [Edit an approved spec](#edit-an-approved-spec).

### 6. Change after archiving: follow-up spec

**When to use:** the spec is already archived and you want to change or extend what it delivered.

1. 💬 **Agent chat:**
   ```text
   Use spf.spec to plan a follow-up of csv-export: let customers pick a date range before exporting.
   ```
2. 💬 **Agent chat:** review and approve, then use the name the agent gave the new spec:
   ```text
   Use spf.implement to implement csv-export-date-range.
   ```

<details><summary>Behind the scenes</summary>

- Archived specs are history and are never reopened. The new spec has its own folder.
- When the change alters rules already in a module living spec, its delta carries them in `## MODIFIED` (or `## REMOVED`) with the complete new block (see [Module deltas](artifacts.md#6-module-deltas)).
- The summary written by `spf.archive` lists suggested next steps and follow-up specs: a good starting point for the request.

</details>

**If it goes wrong:**
- The delta check reports a conflict: see [Delta base conflict](#delta-base-conflict).

### 7. QA failed

**When to use:** the independent QA step reported `FAILED` criteria and the agent escalated to you.

1. 💬 **Agent chat:** read the agent's report: it already lists each failed criterion with its evidence. If the code is wrong:
   ```text
   Use spf.implement to fix the failed QA criteria of invoice-due-date-column and run QA again.
   ```
   If the criterion itself is wrong, follow [workflow 5](#5-change-after-implementing-before-archiving) instead.

<details><summary>Behind the scenes</summary>

- The agent already ran up to two correction cycles (an implementer fixes, a fresh tester re-checks) before escalating, so what reaches you is only what it could not fix on its own.
- The verdict is stored in `spec.yaml` with `specforce implementation qa record <slug>` after every QA run.

</details>

**If it goes wrong:**
- Archiving is refused while the stored verdict is `FAILED`, even with `--force`: see [Archive refused: QA FAILED](#archive-refused-qa-failed).

### 8. Parallel specs across worktrees

**When to use:** you want two agents working on two specs at the same time, each on its own branch.

1. 💻 **Terminal:** create a worktree for the second spec:
   ```bash
   git worktree add ../myapp-csv-export -b feature/csv-export
   ```
2. 💬 **Agent chat:** start a second agent session inside `../myapp-csv-export`, then:
   ```text
   Use spf.spec to plan: let customers export all their invoices to CSV from the dashboard.
   ```
3. 💻 **Terminal:** follow both from your main checkout:
   ```bash
   specforce console
   ```

<details><summary>Behind the scenes</summary>

- Specs from other worktrees are read-only in your checkout: their tasks can only be changed from their own worktree. The Constitution always comes from your current branch (see [Git Worktree Support](git-worktrees.md)).
- Each spec seals its delta. When both change the same module block, the second one to archive sees a conflict and must redo its delta.

</details>

**If it goes wrong:**
- The agent in the new worktree does not see the skills: if your tool folders are not committed, run `specforce init` inside that worktree.
- Two specs touch the same module rule: see [Delta base conflict](#delta-base-conflict).

### 9. Adopt Specforce in an existing project

**When to use:** the project already has code and you want Specforce to document it before the first spec.

1. 💻 **Terminal:** in the project root:
   ```bash
   specforce init
   ```
2. 💬 **Agent chat:**
   ```text
   Use spf.constitution to document this existing project: read the code and draft the constitution from what is already there.
   ```
3. 💬 **Agent chat:** review the documents in `.specforce/docs/`, correct what is wrong, then plan the first change:
   ```text
   Use spf.spec to plan: add a search box to the invoices table that filters by customer name.
   ```

<details><summary>Behind the scenes</summary>

- `specforce init` writes its rules inside a marked block of the root `AGENTS.md`; your existing content outside the block is preserved.
- `spf.constitution` reads the code, asks about decisions it cannot infer and writes the six documents plus `modules/`.

</details>

**If it goes wrong:**
- The constitution describes the code as it is, including habits you want to drop: say so in chat; it records the target rule, not the accident.

### 10. Resume an interrupted implementation

**When to use:** the session crashed, you closed the terminal or the agent stopped in the middle of a roadmap.

1. 💬 **Agent chat:**
   ```text
   Use spf.implement to check the implementation status of invoice-due-date-column and resume it.
   ```

<details><summary>Behind the scenes</summary>

- The agent reads the task states and the repository. Tasks left `in-progress` were interrupted: it inspects the affected code and runs them again. Finished tasks are not redone, and uncommitted changes are kept.

</details>

**If it goes wrong:**
- Do not discard the uncommitted changes before resuming: the agent uses them to continue safely.

### 11. Switch AI tools mid-project

**When to use:** you move to another AI tool (or add one) in the middle of a project or a spec.

1. 💻 **Terminal:** list every tool you want to keep, using the keys from [Supported Tools](supported-tools.md):
   ```bash
   specforce init claude cursor
   ```
   ⚠️ Every installed tool you leave out is uninstalled and **its whole folder is deleted** (for example `.claude/`, with your settings and rules), without asking for confirmation. Commit or back up that folder first if it holds anything of yours.
2. 💬 **Agent chat:** in the new tool:
   ```text
   Use spf.implement to check the implementation status of invoice-due-date-column and resume it.
   ```

<details><summary>Behind the scenes</summary>

- The Constitution, the specs and their progress live in `.specforce/` and do not depend on any tool, so nothing is lost.
- `specforce init` installs the skills of the listed tools and deletes the folder of every installed tool you left out. Without arguments, the interactive picker shows a warning before deleting anything.

</details>

**If it goes wrong:**
- A tool's folder disappeared: you left its key out. Restore the folder from version control if it held your own files, then run `specforce init` again with its key, or without arguments to use the picker.

### 12. Enforce TDD through configuration

**When to use:** you want every roadmap to be test-first and every task to pass your checks before it is marked finished.

1. 💬 **Agent chat:** (or edit the file yourself)
   ```text
   In .specforce/config.yaml, add a tasks instruction requiring [RED] test tasks followed by [GREEN] tasks in the same phase, and on_task_finished hooks "make lint" and "make test".
   ```
   The result looks like this:
   ```yaml
   instructions:
     tasks:
       - "Each phase must contain [RED] tasks (tests for happy path and edge cases) followed by [GREEN] tasks (the implementation)."
   hooks:
     on_task_finished:
       - "make lint"
       - "make test"
   ```
2. 💬 **Agent chat:**
   ```text
   Use spf.spec to plan: let customers download each invoice as a PDF from the dashboard.
   ```

<details><summary>Behind the scenes</summary>

- Tasks get the exact `[RED]` token in their title. When every task finished in a batch is `[RED]`, only `on_task_finished` hooks are skipped (with a notice); phase and all-tasks hooks still run.

</details>

**If it goes wrong:**
- A hook entry with `|` or `&&` never works: see [Hook failure blocks a task](#hook-failure-blocks-a-task).

## Common situations

### Archive refused: tasks are pending

**Symptom:** archiving stops and lists tasks that are not finished.

**Why:** a spec is archived only when every task is finished. The archive skill refuses even if you ask it to archive anyway; `--force` only skips this check and never the delta or QA checks.

**What to do:** the agent already lists the unfinished tasks. Ask `Use spf.implement to finish <slug>`. If a task is no longer needed, ask `spf.spec` to remove it from the spec instead.

### Archive refused: delta missing or pending

**Symptom:** archiving is refused with a delta state of `missing`, `pending` or `invalid`.

**Why:** from `medium` up, the module living specs must reflect the spec before archiving: the delta must be `applied` (or declared `none`). `missing` means no delta was written; `pending` means it was not applied to the modules yet.

**What to do:** the agent reports the delta state. For `missing`, ask `Use spf.spec to write the module delta for <slug>, or declare no module impact`. For `pending`, ask `Use spf.archive to apply the delta and archive <slug>`. For `invalid`, ask `spf.spec` to fix the reported problems.

### Delta base conflict

**Symptom:** the agent reports a delta conflict on a module block, and archiving leaves the module untouched.

**Why:** after your delta was sealed, another spec (often from a parallel branch) changed the same module block. Applying yours would silently overwrite that change.

**What to do:** ask `Use spf.spec to redo the delta of <slug> against the current module and seal it again`. The agent checks the delta again before archiving.

### Archive refused: QA FAILED

**Symptom:** archiving stops and lists failed QA criteria.

**Why:** the stored QA verdict is `FAILED`. Archiving never proceeds over a failed verdict, not even with `--force`.

**What to do:** read the failed criteria the agent lists, then follow [workflow 7](#7-qa-failed).

### Archive warning: no QA report

**Symptom:** archiving continues but warns that no QA report exists.

**Why:** no verdict was recorded for the spec, usually because the implementation was done or finished outside `spf.implement`, or the session ended before QA.

**What to do:** if you want the evidence, stop and ask `Use spf.implement to run the independent QA step for <slug>`, then archive again. Otherwise the warning is informational.

### Hook failure blocks a task

**Symptom:** the agent cannot mark a task finished and shows a hook's output.

**Why:** a hook from `.specforce/config.yaml` exited with a non-zero code. Hooks run without a shell (each entry is split on spaces), in parallel, so `|`, `&&` and redirections are never interpreted. Only an all-`[RED]` batch skips the `on_task_finished` hooks.

**What to do:** ask `Fix the cause of the failing hook; do not skip or remove it`. If the hook entry itself is wrong, edit `.specforce/config.yaml` so it calls a script or a `make` target (see [Configuration](configuration.md)).

### Agent wants to code without a spec

**Symptom:** you ask for a change and the agent refuses to edit code, or stops mid-implementation asking for a spec update.

**Why:** this is intended. Code is changed only through an approved spec, and a change of scope or architecture must go through planning first.

**What to do:** ask `Use spf.spec to plan: <the change>`, or `Use spf.discovery to explore: <the idea>` when it is not clear yet.

### Worker returns a BLOCKED question

**Symptom:** during implementation the agent stops and asks you a question with options and a recommendation.

**Why:** the worker doing a batch cannot ask you directly. When it needs a decision, it stops with a BLOCKED report and the main session relays it to you.

**What to do:** answer in chat, picking an option or the recommendation. If your answer changes the scope, expect the agent to update the spec with `spf.spec` first.

### Resize during implementation

**Symptom:** after a resize, implementation refuses to continue and lists missing artifacts.

**Why:** promoting a spec (for example to `large`) marks the newly required artifacts as missing, and implementation stays blocked until they exist. Demoting never deletes files.

**What to do:** ask `Use spf.spec to write the missing artifacts of <slug>`. After approving them, ask `Use spf.implement to continue <slug>`.

### Edit an approved spec

**Symptom:** you want to change requirements or tasks of a spec you already approved.

**Why:** an approved spec is a contract. Changing it means updating every related artifact together (requirements, design, delta, tasks) and approving it again; editing only `tasks.md` leaves the spec inconsistent.

**What to do:** ask `Use spf.spec to update <slug>: <the change>. Keep all artifacts consistent`, review, approve, then follow [workflow 5](#5-change-after-implementing-before-archiving).
