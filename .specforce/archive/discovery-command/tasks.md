---
slug: discovery-command
lens: Integration
---

# Implementation Roadmap: Discovery Command

## 1. Execution Strategy
- **Gravity Order:** Kit Resource Definition -> Agent Translation Verification.
- Since this is an integration feature affecting the agent's behavior prompts, the focus is on high-fidelity prompt engineering and ensuring the translator correctly propagates the new command to all targets.

## 2. Tasks

### Phase 1: Kit Resource Definition

#### T1.1: [SCAFFOLD] Create Discovery Command YAML
**State:** `[FINISHED]`
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [REQ-1, REQ-2, REQ-3]

**Action Steps:**
- Create the `discovery.yaml` file in the kit commands directory.
- Define the command `name` as `spf.discovery`.
- Add a descriptive `description` for agent discovery.
- Configure `mapping` for `open-code`, `kilo-code`, and `kimi-code`.
- Write the Markdown `content` prompt implementing the Brainstormer/Detective personas, read-only constraints, and handoff rules.

**Verification (TDD):**
- Run `ls src/internal/agent/kit/commands/discovery.yaml` to confirm file existence.
- Review the content to ensure all placeholders and personas from the Requirements are present.

### Phase 2: Translation & Propagation

#### T2.1: [VERIFY] Run Kit Translation Tests
**State:** `[FINISHED]`
**Target:** `src/internal/agent/translator_test.go`
**Context:** [REQ-4]

**Action Steps:**
- Execute the project's test suite for the agent domain to ensure the new YAML is valid and can be translated without errors.
- Run `go test ./src/internal/agent/...`

**Verification (TDD):**
- `go test -v ./src/internal/agent/ -run TestTranslator` (or equivalent test that processes kit manifests).

#### T2.2: [VERIFY] Manual Inspection of Generated Commands
**State:** `[FINISHED]`
**Target:** `Global Scope`
**Context:** [REQ-4]

**Action Steps:**
- Run the kit generator (e.g., `make run` or a specialized build command that updates the `.gemini/`, `.claude/`, etc. folders).
- Inspect `.gemini/commands/spf.discovery.toml` (if generated) or equivalent for other agents.

**Verification (TDD):**
- Verify that the generated file contains the prompt defined in `discovery.yaml`.

### Phase 3: Bootstrap & Documentation

#### T3.1: [BOOTSTRAP] Update AGENTS.md Template
**State:** `[FINISHED]`
**Target:** `src/internal/project/agents_md.go`
**Context:** [REQ-1]

**Action Steps:**
- Locate the template used for `AGENTS.md` generation.
- Add the `Discovery (/discovery)` command to the SDD Protocol section.
- Define its purpose as a read-only brainstorming/diagnostic entry point.

**Verification (TDD):**
- Run `go test ./src/internal/project/...` to ensure template changes don't break generation tests.

#### T3.2: [DOCS] Update Project Documentation
**State:** `[FINISHED]`
**Target:** `README.md`, `docs/cli.md`
**Context:** [REQ-1, REQ-3]

**Action Steps:**
- Update `README.md` to include `spf.discovery` in the command list.
- Update `docs/cli.md` with a detailed explanation of the Discovery mode and its handoff to `/spec`.

**Verification (TDD):**
- Inspect the updated files for clarity and correctness.
