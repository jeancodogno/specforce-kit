---
slug: 20260527-0014-standardize-agents-naming-and-prompts
lens: Integration
---

# Implementation Roadmap: Standardize Agents Naming and Prompts

## 1. Execution Strategy
- **Gravity Order:** Agent Renaming -> Prompt Hardening (safety) -> Command Orchestration Updates -> Validation.

## 2. Tasks

### Phase 1: Agent Naming Standardization

- [x] T1.1: [FS] Rename Agent Files to `specforce-` prefix
**Target:** `src/internal/agent/kit/agents/`
**Context:** [US-1]

**Action Steps:**
- Rename `product-analyst.yaml` to `specforce-product-analyst.yaml`.
- Rename `technical-solution-architect.yaml` to `specforce-architect.yaml`.
- Rename `technical-project-planner.yaml` to `specforce-planner.yaml`.
- Rename `technical-developer.yaml` to `specforce-developer.yaml`.
- Rename `technical-qa-engineer.yaml` to `specforce-qa.yaml`.

**Acceptance Check:**
`ls src/internal/agent/kit/agents/ | grep "specforce-"` (Should show 5 files).

- [x] T1.2: [CODE] Update Internal Agent Names
**Target:** `src/internal/agent/kit/agents/specforce-*.yaml`
**Context:** [US-1]

**Action Steps:**
- Update `name` field in each file to match the new filename (without .yaml).
- Ensure `description` remains descriptive but concise.

**Acceptance Check:**
`grep "name: specforce-" src/internal/agent/kit/agents/*.yaml` (Should return 5 matches).

### Phase 2: Prompt Hardening & Portability

- [x] T2.1: [FS] Remove Mapping Blocks from Agents
**Target:** `src/internal/agent/kit/agents/specforce-*.yaml`
**Context:** [US-3]

**Action Steps:**
- Remove the `mapping` block from each agent YAML.
- Ensure the YAML structure remains valid after removal.

**Acceptance Check:**
`grep -v "mapping:" src/internal/agent/kit/agents/specforce-*.yaml` (Should confirm removal).

- [x] T2.2: [CODE] Inject Safety & Boundary Rules into Prompts
**Target:** `src/internal/agent/kit/agents/specforce-*.yaml`
**Context:** [US-3, US-4]

**Action Steps:**
- Add the "Environment Awareness & Safety" block to the `content` of all agents.
- Add strict "No Implementation Code" rules to `specforce-product-analyst` and `specforce-planner`.
- Update `specforce-qa` handoff message to point to `/spf:archive` and `/spf:implement`.

**Acceptance Check:**
`grep "Non-Recursive Mandate" src/internal/agent/kit/agents/*.yaml` (Should return 5 matches).

### Phase 3: Command Orchestration Updates

- [x] T3.1: [CODE] Update `spf.spec` Delegation & Names
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [US-2]

**Action Steps:**
- Update all agent references to use `specforce-` prefix.
- Change `tasks` artifact delegation from `specforce-developer` (formerly `technical-developer`) to `specforce-planner`.

**Acceptance Check:**
`grep "specforce-planner" src/internal/agent/kit/commands/spec.yaml` (Should confirm delegation change).

- [x] T3.2: [CODE] Update `spf.implement` and others
**Target:** `src/internal/agent/kit/commands/*.yaml`
**Context:** [US-2]

**Action Steps:**
- Update `implement.yaml` agent names to `specforce-developer` and `specforce-qa`.
- Update `discovery.yaml` prompt to use `specforce-scout` (internal reference).
- Update `constitution.yaml` to use the new `specforce-architect` or related specialists.

**Acceptance Check:**
`grep "specforce-" src/internal/agent/kit/commands/*.yaml` (Verify all commands use the new prefix).

### Phase 4: Fix Tests

- [x] T4.1: [CODE] Update Compatibility Tests
**Target:** `src/internal/agent/compatibility_test.go`
**Context:** [US-1]

**Action Steps:**
- Revert or adjust `compatibility_test.go` to handle the removal of the mapping block.
- Ensure tests point to the new `specforce-developer.yaml`.

**Acceptance Check:**
`make test` (Should pass all tests).

### Phase 5: Global Validation

- [x] T5.1: [TEST] Final Specforce Status Validation
**Target:** `Global Scope`
**Context:** [US-1, US-2]

**Action Steps:**
- Run `specforce spec status 20260527-0014-standardize-agents-naming-and-prompts`.
- Verify that the CLI correctly parses all YAMLs without errors.

**Acceptance Check:**
`specforce spec status 20260527-0014-standardize-agents-naming-and-prompts --json` (Should return `is_valid: true`).
