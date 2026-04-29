---
slug: autonomous-sdd-workflow
lens: Integration
---

# Implementation Roadmap: Autonomous SDD Workflow via Skills

## 1. Execution Strategy
- **Gravity Order:** Update Command Descriptions -> Update AGENTS.md -> Verify Agent Behavior

## 2. Tasks

### Phase 1: Dual-Installation Mechanism

#### T1.1: [CODE] Modify Kit Mapping for Dual Installation
**State:** `[FINISHED]`
**Target:** `src/internal/agent/kit/kit.yaml`
**Context:** [REQ-3]

**Action Steps:**
- Update the `kit.yaml` mappings for supported agents to ensure `commands/*.yaml` are mapped to both their native command path and the agent's skills directory.
- If `kit.yaml` cannot express multiple destinations for a single source, update the Go `translator.go` or `installer.go` logic to duplicate the output.

**Verification (TDD):**
Run `go test ./src/internal/agent/...` and ensure `commands` exist in both target folders during kit installation tests.

### Phase 2: Command Configuration Refinement

#### T2.1: [CONFIG] Update spec.yaml description
**State:** `[FINISHED]`
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [REQ-3]

**Action Steps:**
- Modify the `description` field of the `spec` command to clearly state its role as an orchestrator and define the exact scenarios when the LLM should invoke it autonomously.

**Verification (TDD):**
`cat src/internal/agent/kit/commands/spec.yaml | grep "description"` should reflect the new intent-driven description.

#### T2.2: [CONFIG] Update implement.yaml description
**State:** `[FINISHED]`
**Target:** `src/internal/agent/kit/commands/implement.yaml`
**Context:** [REQ-3]

**Action Steps:**
- Modify the `description` field of the `implement` command to clearly state its role as the TDD execution engine and when it should be triggered.

**Verification (TDD):**
`cat src/internal/agent/kit/commands/implement.yaml | grep "description"` should reflect the new intent-driven description.

### Phase 3: Agent Instruction Updates

#### T3.1: [DOCS] Update AGENTS.md Protocol
**State:** `[FINISHED]`
**Target:** `AGENTS.md`
**Context:** [REQ-1][REQ-2]

**Action Steps:**
- Rewrite the "Spec-Driven Development (SDD) Protocol" section.
- Explicitly instruct the agent to operate exclusively through the Specforce workflow commands (`spf.spec`, `spf.implement`, `spf.constitution`).
- Add the "Proactive Mandate" instructing the agent not to wait for explicit slash commands but to invoke them when the intent matches the phase.

**Verification (TDD):**
`cat AGENTS.md | grep "Proactive Mandate"` should exist.
