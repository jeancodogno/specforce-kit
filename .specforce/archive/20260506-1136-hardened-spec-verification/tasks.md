---
slug: 20260506-1136-hardened-spec-verification
lens: Backend-heavy
---

# Implementation Roadmap: Hardened Spec Verification

## 1. Execution Strategy
The implementation follows a single-phase approach to update the orchestration instructions in the `spec.yaml` kit command.

## 2. Tasks

### Phase 1: Hardening the Orchestration Workflow

- [x] T1.1: [CONFIG] Update spec.yaml Step 3: Verification & Handoff
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [REQ-1]

**Action Steps:**
- Update the section `### 3. Verification & Handoff` to make the `status --json` command a mandatory terminal step.
- Clarify that the summary MUST only be output if verification is successful.

**Verification (TDD):**
`cat src/internal/agent/kit/commands/spec.yaml | grep "execute specforce spec status <slug> --json ONE final time"`

- [x] T1.2: [CONFIG] Harden spec.yaml Guardrails
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [REQ-2]

**Action Steps:**
- Add a new guardrail entry enforcing the final verification check for all sessions.

**Verification (TDD):**
`cat src/internal/agent/kit/commands/spec.yaml | grep "MANDATORY FINAL CHECK"`
