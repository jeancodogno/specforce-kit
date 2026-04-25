# Specforce Configuration

Configuring Specforce's internal behavior is done through the `.specforce/config.yaml` file. With it, you add layers of intelligence, security, and standardization that apply invisibly to your AI agents.

## 1. Context-Aware Instructions

Instead of repeating rules ("Remember to use BDD") in every prompt, you can configure them to be injected automatically at the exact moment the agent works on a specific artifact.

In `config.yaml`:

```yaml
instructions:
  requirements:
    - "Always use BDD GIVEN/WHEN/THEN syntax for acceptance criteria."
    - "Ensure accessibility is always mentioned if UI is involved."
  design:
    - "Always evaluate security impacts for new API routes."
  implementation:
    - "Do not bypass TypeScript typings. Use Type Guards."
  archive:
    - "Always update the project memorial with lessons learned."
    - "Summarize any technical debt introduced during implementation."
```

## 2. Validation Hooks (Quality Gates)

Hooks are external commands executed by Specforce before or after a state transition. Their main use is to act as **Quality Gates** to prevent the agent from finishing tasks if the code is inadequate.

In `config.yaml`:

```yaml
hooks:
  on_task_finished:
    - "npm run lint"
    - "npm run test"
```

### How does it work?
1. The agent attempts to mark a task from its specification as finished.
2. Specforce pauses the agent and executes `npm run lint` and `npm run test`.
3. **If successful (exit 0):** The task is marked as finished.
4. **If failed:** The transition is blocked, and the error log is returned to the agent so it can fix the bug before trying to proceed.

### Essential Tips for Hooks:
- **Be Fast and Silent:** Use fast and silent commands (`--silent`, `--quiet`).
- **Filter the Output:** Do not dump thousands of log lines into the agent's context. Use focused tools or filters like `grep` to send only specific errors. Example: `golangci-lint run --out-format=line-number` or `npm run test --silent | grep -A 10 "FAIL"`.

---

By combining *Context-Aware Instructions* and *Validation Hooks*, you create an unbreakable "guardrail" for your team and your AI agents.