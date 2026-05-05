---
slug: 20260505-1201-protocol-trigger-optimization
lens: Balanced full-stack
---

# Implementation Roadmap: Protocol Trigger Optimization

## 1. Execution Strategy
- **Gravity Order:** Update documentation rules (AGENTS.md, engineering.md) -> Update Go source templates -> Update skill metadata.

## 2. Tasks

### Phase 1: Update behavioral rules and guidelines

- [x] T1.1: [DOCS] Update AGENTS.md Protocol Triggers
**Target:** `AGENTS.md`
**Context:** [REQ-1]

**Action Steps:**
- Add explicit trigger keywords to `spf.discovery` (brainstorm, vague, research).
- Add explicit trigger keywords to `spf.spec` (plan, formalize, new feature).
- Add a "Direct Edit Prohibition" rule: Agents MUST NOT use `replace` or `write_file` if no approved spec exists.

**Verification (TDD):**
`cat AGENTS.md | grep "Proactive Mandate"` should show the new rules.

- [x] T1.2: [DOCS] Clarify Orchestration Policy in engineering.md
**Target:** `.specforce/docs/engineering.md`
**Context:** [REQ-3]

**Action Steps:**
- Update "Agent Orchestration Protocol" to explicitly permit `spf.spec` while forbidding "Blackbox" LLM planning.

**Verification (TDD):**
`cat .specforce/docs/engineering.md | grep "Primary Orchestration Only"` should show the distinction.

### Phase 2: Persist changes in Go template

- [x] T2.1: [CODE] Update Go Template for AGENTS.md
**Target:** `src/internal/project/agents_md.go`
**Context:** [REQ-1]

**Action Steps:**
- Update `agentsMDTemplate` constant with the same changes made to `AGENTS.md` in T1.1.

**Verification (TDD):**
`grep "Proactive Mandate" src/internal/project/agents_md.go` should show the updated string.

### Phase 3: Update Source Kit Metadata

- [x] T3.1: [DOCS] Refine Discovery Metadata in Kit
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [REQ-2]

**Action Steps:**
- Update `description` to include "brainstorm", "research", "diagnose", "bug investigation", "root cause analysis", and "vague intent".

**Verification (TDD):**
`cat src/internal/agent/kit/commands/discovery.yaml | grep "bug investigation"` should show the update.

- [x] T3.2: [DOCS] Refine Spec Metadata in Kit
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [REQ-2]

**Action Steps:**
- Update `description` to include "plan", "initialize", "update specs", and "formalize".

**Verification (TDD):**
`cat src/internal/agent/kit/commands/spec.yaml | grep "formalize"` should show the update.

