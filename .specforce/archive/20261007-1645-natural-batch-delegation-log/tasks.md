---
slug: 20261007-1645-natural-batch-delegation-log
lens: Backend-heavy
---

# Implementation Roadmap: Natural Batch Delegation Log

## 1. Execution Strategy
- **Gravity Order:** Update Unit Tests (RED) -> Update `spf-implement` embedded template & agent kit rules (GREEN) -> Verify manifest tests and embedded bundle integrity.

## 2. Tasks

### Phase 1: Natural Log Directive & Test Alignment

- [x] T1.1: [TEST] [RED] Update Agent Kit Manifest Tests for Natural Batch Delegation Log
**Target:** `src/internal/agent/kit_manifests_test.go`
**Context:** [US-KIT-01]

**Action Steps:**
- Update `TestImplementBatchSubagentIsolationAndProgress` to assert the new natural delegation log tokens (`Iniciando Lote` or `Delegating Batch`, `concluídas` or `completed`, and summary scope).
- Ensure assertions validate the removal of repetitive raw bracket chains in favor of conversational, informative delegation logging.
- Run `go test ./src/internal/agent -run TestImplementBatchSubagentIsolationAndProgress` to verify the test fails or reflects the expected contract.

**Acceptance Check:**
`go test ./src/internal/agent -run TestImplementBatchSubagentIsolationAndProgress`

- [x] T1.2: [CODE] [GREEN] Update spf-implement Skill and Constitution Business Rules
**Target:** `src/internal/agent/kit/skills/spf-implement/SKILL.yaml`
**Context:** [US-KIT-01]

**Action Steps:**
- Replace the rigid bracketed `MANDATORY LOG` directive in `src/internal/agent/kit/skills/spf-implement/SKILL.yaml` with the natural, contextual Option 1 format:
  ```markdown
  🚀 **Iniciando Lote {batch_index}/{total_batches}** ({completed_tasks_count}/{total_tasks_count} concluídas · {progress_pct}%)
  {short_conversational_scope_summary}
  Delegando para perfil de *{required_persona_role}* (tier {recommended_model_tier}, esforço {recommended_effort}).
  ```
- Update `.specforce/docs/modules/agent-kit.md` rule `[BR-KIT-09]` to reflect the natural conversational batch delegation log standard.
- Sync active workspace skill file `.agents/skills/spf-implement/SKILL.md` (or run sync/tests) to ensure the current environment reflects the new template.
- Execute unit and integration tests across the agent package.

**Acceptance Check:**
`go test ./src/internal/agent/...`
