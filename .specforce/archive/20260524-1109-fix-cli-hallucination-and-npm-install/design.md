---
slug: 20260524-1109-fix-cli-hallucination-and-npm-install
lens: Bugfix
---

# Technical Design: CLI Path Hardening & NPM Recovery Guide (Fix Blueprint)

## 1. Code Path Inventory
- `src/internal/core/constants.go` -> Add `Version` constant to centralize the tool version.
- `src/internal/project/agents_md.go` -> Update `agentsMDTemplate` to include explicit CLI execution rules and recovery steps.
- `src/internal/agent/kit/commands/archive.yaml` -> Add `Guardrails` section.
- `src/internal/agent/kit/commands/constitution.yaml` -> Add `Guardrails` section.
- `src/internal/agent/kit/commands/discovery.yaml` -> Add `Guardrails` section.
- `src/internal/agent/kit/commands/implement.yaml` -> Add `Guardrails` section.
- `src/internal/agent/kit/commands/spec.yaml` -> Add `Guardrails` section.

## 2. Regression Strategy (Verification Plan)
- **Unit Tests:**
    - Update `src/internal/project/agents_md_test.go` to assert that the generated `AGENTS.md` contains the strings "CLI Execution & Environment" and "Environment Recovery".
- **Integration Tests:**
    - N/A (Documentation change).
- **Manual Verification:**
    - Generate `AGENTS.md` using a mock setup or by running `specforce init` in a test directory and verify the content.
    - Inspect the source YAML files in `src/internal/agent/kit/commands/` to ensure the `Guardrails` are present.

## 3. Side Effects & Risks
- **Performance:** No impact on runtime performance.
- **Compatibility:** No impact on existing CLI commands or workflows.
- **Context Size:** Slightly increases the context size of the `AGENTS.md` and command instructions, but this is a necessary trade-off for reliability.

## 4. Proposed Fix (Abstract Logic)

### core/constants.go
```go
const Version = "1.0.0-alpha.1"
```

### project/agents_md.go
```markdown
## 2. CLI Execution & Environment
**CRITICAL:** The `specforce` CLI is globally available in your environment path. You MUST execute it directly as a standard command (e.g., `specforce spec list --json`). 
- DO NOT use `npx specforce`.
- DO NOT use relative or absolute paths like `./bin/specforce`.
- If a command fails, DO NOT hallucinate alternative execution paths.

## 3. Environment Recovery
If the `specforce` binary is not found in the environment:
1. Run: `npm i -g @jeancodogno/specforce-kit@1.0.0-alpha.1`
2. Ensure your global npm bin directory is in your PATH.
```

### agent/kit/commands/*.yaml
```yaml
  ## Guardrails
  - **CLI Execution:** The `specforce` CLI is globally available. Execute it directly (e.g., `specforce spec list`). DO NOT use `npx` or relative paths.
```
