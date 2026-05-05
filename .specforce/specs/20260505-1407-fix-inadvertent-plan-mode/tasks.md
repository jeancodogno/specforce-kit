---
slug: 20260505-1407-fix-inadvertent-plan-mode
lens: Integration
---

# Implementation Roadmap: Fix Inadvertent Plan Mode Activation

## 1. Execution Strategy
- **Gravity Order:** Update Master Templates (YAML) in `src/` -> Rebuild/Verify.

## 2. Tasks

### Phase 1: Skill Orchestration Hardening (Source)

- [x] T1.1: [CODE] Update spec.yaml Master Template
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [REQ-1]

**Action Steps:**
- Rewrite Rule #3 in the `content` block to prohibit platform-native planning modes.
- Add an "Internalized Planning" guardrail in the `content` block.
- Update Step 2.B (Agent Discovery) in the `content` block to include a "Context-Isolation" directive.

**Verification (TDD):**
- Verify the file content contains the new instructions.

### Phase 2: Subagent Constraint Hardening (Source)

- [x] T2.1: [CODE] Update Technical Solution Architect Master Template
**Target:** `src/internal/agent/kit/agents/technical-solution-architect.yaml`
**Context:** [REQ-2]

**Action Steps:**
- Add an "Atomic Execution" guardrail to the `content` block.

**Verification (TDD):**
- Verify the file content.

- [x] T2.2: [CODE] Update Product Analyst Master Template
**Target:** `src/internal/agent/kit/agents/product-analyst.yaml`
**Context:** [REQ-2]

**Action Steps:**
- Add an "Atomic Execution" guardrail to the `content` block.

**Verification (TDD):**
- Verify the file content.

- [x] T2.3: [CODE] Update Technical Project Planner Master Template
**Target:** `src/internal/agent/kit/agents/technical-project-planner.yaml`
**Context:** [REQ-2]

**Action Steps:**
- Add an "Atomic Execution" guardrail to the `content` block.

**Verification (TDD):**
- Verify the file content.

- [x] T2.4: [CODE] Update Technical QA Engineer Master Template
**Target:** `src/internal/agent/kit/agents/technical-qa-engineer.yaml`
**Context:** [REQ-2]

**Action Steps:**
- Add an "Atomic Execution" guardrail to the `content` block.

**Verification (TDD):**
- Verify the file content.
