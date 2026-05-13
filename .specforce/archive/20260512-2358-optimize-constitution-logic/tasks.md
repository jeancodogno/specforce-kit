---
slug: 20260512-2358-optimize-constitution-logic
lens: Balanced full-stack
---

# Implementation Roadmap: Constitution Logic Optimization

## 1. Execution Strategy
- **Gravity Order:** Skill Blueprints -> Go Template Update -> Integration Testing.
- We update the Kit source YAMLs first to ensure the instructions are available for all agents. Then we update the Go code responsible for generating the `AGENTS.md` file.

## 2. Tasks

### Phase 1: Skill Blueprint Updates

- [x] T1.1: [CODE] Add Ecosystem Contextualization to Discovery
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [US-2]

**Action Steps:**
- Locate the `## Ecosystem Contextualization` section.
- Prepend the instructions to explicitly run `specforce spec list --json` and `specforce constitution status --json` if the context is not already available in the session history.
- Ensure the instruction emphasizes "synchronizing with the project's active state".

**Verification (TDD):**
Manually inspect the file to ensure the new instructions are correctly formatted within the YAML block.

- [x] T1.2: [CODE] Add Start & End Verification to Planning (Spec)
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [US-1, US-3]

**Action Steps:**
- Update Step 1 "Discovery & Intent Clarification".
- Add the instruction: "Execute `specforce spec status <slug> --json` ONCE at the start to map out your work and ONCE at the end to verify. DO NOT poll status between individual artifact generations."

**Verification (TDD):**
Manually inspect the file to ensure the Start & End pattern is present in the `spf.spec` pipeline.

- [x] T1.3: [CODE] Add Start & End Verification to Implementation
**Target:** `src/internal/agent/kit/commands/implement.yaml`
**Context:** [US-1, US-3]

**Action Steps:**
- Update Phase 1 "Initialization & Roadmap Mapping".
- Ensure the instruction explicitly states: "Execute status ONCE at the start and ONCE at the end. DO NOT run the status command again during this session."

**Verification (TDD):**
Manually inspect the file.

- [x] T1.4: [CODE] Add Memory-Aware Clause to Governance (Constitution)
**Target:** `src/internal/agent/kit/commands/constitution.yaml`
**Context:** [US-1]

**Action Steps:**
- Update Step 1 "Scope & Status Discovery".
- Add instruction: "Check your history for existing constitution status. REUSE it if available to save tokens."

**Verification (TDD):**
Manually inspect the file.

### Phase 2: AGENTS.md Template Update

- [x] T2.1: [CODE] Update AGENTS.md Template with Global Efficiency Guidelines
**Target:** `src/internal/project/agents_md.go`
**Context:** [US-3]

**Action Steps:**
- Update the `agentsMDTemplate` constant.
- Add a new section `## 4. Efficiency & Token Optimization`.
- Include rules for **Constitution Context** reuse.
- Include rules for "Surgical Reads" (preferring `grep_search`).
- Include rules for "Parallelism".
- **CRITICAL:** Ensure NO Spec or Implementation status rules are added here.

**Verification (TDD):**
Run `go test ./src/internal/project/agents_md_test.go` to ensure no regressions in merging logic.

### Phase 3: Integration Verification

- [x] T3.1: [CLI] Verify Updated Project Initialization
**Target:** `Global Scope`
**Context:** [US-1, US-2, US-3]

**Action Steps:**
- Build the project: `go build -o specforce src/cmd/specforce/main.go`.
- Create a temporary test directory.
- Run `./specforce init gemini-cli`.
- Verify that the generated `AGENTS.md` contains the new efficiency section.
- Verify that the generated `.gemini/rules/spf-discovery.toml` (or equivalent adapted artifact) contains the new instructions.

**Verification (TDD):**
The presence of the new instructions in the generated files confirms successful propagation of the Kit changes.
