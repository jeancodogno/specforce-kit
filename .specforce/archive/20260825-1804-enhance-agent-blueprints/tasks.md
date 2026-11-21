---
slug: 20260825-1804-enhance-agent-blueprints
lens: Integration
---

# Implementation Roadmap: Enhance Agent Command Blueprints

## 1. Execution Strategy
- **Gravity Order:** Test Assertions (RED) -> Discovery Blueprint (GREEN) -> Implement Blueprint (GREEN) -> Archive Blueprint & Instructions (GREEN) -> Full Test Suite Verification.

## 2. Tasks

### Phase 1: Blueprint Enhancements & Guardrail Verification

- [x] T1.1: [RED] Add unit test assertions for new blueprint directives in test suite
**Target:** `src/internal/agent/kit_manifests_test.go`
**Context:** [US-1], [US-2], [US-3]

**Action Steps:**
- Define `TestEnhancedBlueprintsDirectives` in `kit_manifests_test.go` checking `discovery.yaml` for proactive suggestions and strong architectural opinions.
- Add assertions in `TestEnhancedBlueprintsDirectives` verifying `implement.yaml` contains worker reuse and subagent feedback loop directives.
- Add assertions in `TestEnhancedBlueprintsDirectives` checking `archive.yaml` and `archive.md` for suggested next steps and follow-up specs sections.

**Acceptance Check:**
`go test -v ./src/internal/agent -run TestEnhancedBlueprintsDirectives` (Fails until blueprints are updated).

- [x] T1.2: [GREEN] Update Discovery blueprint with proactive suggestions and architectural opinions
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [US-1]

**Action Steps:**
- Update "The Scout's Stance & Constitutional Anchor" section in `discovery.yaml` to include "Proactive Suggestions & Technical Opinions" directive.
- Enrich "Brainstorming & Architecture Ideation" to mandate concrete recommendations ("Strong Opinions, Weakly Held"), pattern improvements, and edge-case anticipation.
- Update "Bug Detective & Root Cause Isolation" with concrete fix strategies and architectural opinions.

**Acceptance Check:**
`go test -v ./src/internal/agent -run TestEnhancedBlueprintsDirectives` passes discovery checks.

- [x] T1.3: [GREEN] Update Implement blueprint with worker reuse, universal harness compatibility, and adaptive effort routing
**Target:** `src/internal/agent/kit/commands/implement.yaml`
**Context:** [US-2]

**Action Steps:**
- Add universal harness compatibility guidelines under "Subagent Orchestration & Sizing Directives", requiring explicit worker role/tool registration across all harnesses (Claude Code, OpenCode, Antigravity, etc.) and avoiding fragile implicit self-inheritance.
- Add adaptive model tier (`fast_lite`, `flash`, `pro`) and reasoning effort routing (`low`, `medium`, `high`) scaled to batch complexity.
- Add worker continuity and anti-churn guidelines mandating message feedback to the same subagent session.
- Update "Verification & Batch Resolution / Self-Healing" protocol to send compiler/test error logs to the same subagent for up to 2 attempts before escalating.
- Ensure all existing mandatory guardrails and test invariance rules are strictly preserved.

**Acceptance Check:**
`go test -v ./src/internal/agent -run TestImplementBlueprintGuardrails` and `go test -v ./src/internal/agent -run TestEnhancedBlueprintsDirectives` pass.

- [x] T1.4: [GREEN] Update Archive blueprint and instructions with conceptual next steps
**Target:** `src/internal/agent/kit/commands/archive.yaml`
**Context:** [US-3]

**Action Steps:**
- Add step 4 in `src/internal/agent/kit/commands/archive.yaml` directing the agent to suggest next steps and follow-up specs after the final summary.
- Update `src/internal/agent/kit/instructions/archive.md` step 7 markdown template to include the "Suggested Next Steps & Follow-up Specs" section.
- Ensure instructions emphasize conceptual next steps and architectural follow-ups without hallucinating unrelated system requirements.

**Acceptance Check:**
`go test -v ./src/internal/agent -run TestEnhancedBlueprintsDirectives` and `go test ./src/internal/cli -run TestHandleArchiveInstructions` pass.
