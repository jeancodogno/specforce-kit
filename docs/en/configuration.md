[ [English](configuration.md) | [Português](../pt/configuration.md) | [Español](../es/configuration.md) ]

# Specforce Configuration

Specforce's behavior is configured in the `.specforce/config.yaml` file (created by `specforce init`). With it, you add project rules, shared variables and quality gates that apply automatically to your AI agents.

## 1. Context-Aware Instructions

Instead of repeating rules ("Remember to use BDD") in every prompt, you can configure them to be injected automatically at the exact moment the agent works on a specific artifact or phase.

In `config.yaml`:

```yaml
instructions:
  requirements:
    - "Always use BDD GIVEN/WHEN/THEN syntax for acceptance criteria."
    - "Ensure accessibility is always mentioned if UI is involved."
  bug-requirements:
    - "Always include the exact error message in the evidence section."
  design:
    - "Always evaluate security impacts for new API routes."
  tasks:
    - "Structure every phase as [RED] test tasks followed by [GREEN] implementation tasks."
  implementation:
    - "Do not bypass TypeScript typings. Use Type Guards."
  archive:
    - "Always consolidate module living specs with lessons learned."
    - "Summarize any technical debt introduced during implementation."
```

### Accepted instruction keys

| Key | Applied when |
|---|---|
| `requirements` | the agent writes `requirements.md` (feature and bug specs) |
| `bug-requirements` | the agent writes `requirements.md` of a bug spec |
| `design` | the agent writes `design.md` (feature and bug specs) |
| `bug-design` | the agent writes `design.md` of a bug spec |
| `tasks` | the agent writes `tasks.md` |
| `implementation` | the agent loads the roadmap to implement it |
| `archive` | the agent archives a spec |

**Generic and specific keys are combined.** When the agent works on a variant artifact such as `bug-requirements`, it receives the instructions of the generic key (`requirements`) first, followed by those of the specific key (`bug-requirements`); a rule listed under both appears only once. The same applies to `design` and `bug-design`.

## 2. Instruction Variable Injection

To make instructions even more powerful, you can define project-level variables in the `context` block. These variables can then be injected into any instruction or artifact template using the `{{variable_name}}` syntax.

In `config.yaml`:

```yaml
# Global context variables
context:
  project_name: "Specforce Kit"
  tech_stack: "Go, React, PostgreSQL"
  primary_branch: "main"

instructions:
  tasks:
    - "Remember: the code will be merged into the {{primary_branch}} branch."
    - "Every new feature in {{project_name}} must respect the {{tech_stack}} stack."
```

### Why use variables?
- **Single Source of Truth:** Update the project name or stack once, and all agent-facing instructions stay in sync.
- **Reusable Templates:** Use the same instruction set across different projects by only changing the `context` block.
- **Consistency:** Ensures agents always have correct, non-hallucinated project data.

Never put secrets (tokens, passwords, keys) in `context` or `instructions`: these values are shown to the agent. Keep secrets in environment files.

## 3. Validation Hooks (Quality Gates)

Hooks are external commands that Specforce runs when the agent marks tasks as finished. Their main use is to act as **Quality Gates** that prevent the agent from finishing tasks while the code is inadequate.

In `config.yaml`:

```yaml
hooks:
  on_task_finished:
    - "npm run lint"
    - "npm run test"
  on_phase_finished:
    - "npm run integration-tests"
  on_all_tasks_finished:
    - "npm run e2e-tests"
```

### Supported Hooks:
- **`on_task_finished`**: Triggered whenever tasks are marked as finished.
- **`on_phase_finished`**: Triggered when a finished task is the last one of its Phase (H3 heading in `tasks.md`).
- **`on_all_tasks_finished`**: Triggered when a finished task is the last one of the roadmap.

### How does it work?
1. The agent marks one or more tasks as finished (`specforce implementation update <slug> --task <ids> --status finished`).
2. Before changing any state, Specforce runs the applicable hooks once for the whole batch, without duplicates (in the example above, `npm run lint` and `npm run test`).
3. **If every hook exits with code 0:** the tasks are marked as finished.
4. **If any hook exits with a non-zero code:** the transition is blocked and the hook output is returned to the agent so it can fix the problem before trying again.

Test-first tasks: when every task of the batch carries `[RED]` in its title, only the `on_task_finished` hooks are skipped (with a notice naming them); phase and all-tasks hooks still run. See [Artifacts](artifacts.md#8-test-first-red-tasks).

### Essential Tips for Hooks:
- **One command per entry:** each entry is split on spaces and run directly, not through a shell, so pipes (`|`), redirections and `&&` are not interpreted. Put pipelines in a script or a `make` target and call that instead (for example `make test-quiet`).
- **Hooks of one batch run at the same time:** do not make one hook depend on the result of another.
- **Be Fast and Silent:** Use fast and silent commands (`--silent`, `--quiet`).
- **Filter the Output:** Do not dump thousands of log lines into the agent's context. Inside your script, send only the relevant errors, e.g. `golangci-lint run --out-format=line-number` or `npm run test --silent | grep -A 10 "FAIL"`.

---

By combining *Context-Aware Instructions* and *Validation Hooks*, you create an unbreakable "guardrail" for your team and your AI agents.
