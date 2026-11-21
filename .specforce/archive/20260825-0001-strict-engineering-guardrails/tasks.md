---
slug: 20260825-0001-strict-engineering-guardrails
lens: Integration
---

# Implementation Roadmap: Strict Engineering Guardrails and Test Invariance

## 1. Execution Strategy
- **Gravity Order:** Blueprint Templates & Manifest Tests (TDD Red-Green) -> Unabridged Blueprint Enhancement -> Skill & Constitution Living Docs Synchronization -> Global Verification Suite.

## 2. Tasks

### Phase 1: Kit Blueprints & Validation Tests

- [x] T1.1: [RED] Write Blueprint Validation Tests for Implementation & Engineering Blueprints
**Target:** `src/internal/agent/kit_manifests_test.go`
**Context:** [US-1], [US-2], [US-3]

**Action Steps:**
- Add `TestImplementBlueprintGuardrails` test verifying that `implement.yaml` contains the complete `NON-NEGOTIABLE WORKER GUARDRAILS` block.
- Assert that `implement.yaml` includes all pre-coding directives (uncertainty, ambiguity, multiple interpretations, anti-sycophancy).
- Assert that `implement.yaml` includes all during-coding directives (no unrequested features, no unrequested refactoring, no impossible tests, no dead code removal, never delete/weaken/skip tests, confirm genuinely broken tests, tests = spec).
- Assert that `implement.yaml` includes post-implementation directives (activate spec skill on divergence).
- Assert that `engineering.yaml` contains all 14 exhaustive AI constraints across Pre-Coding, During Coding, and Post-Implementation.

**Acceptance Check:**
`go test -v ./src/internal/agent -run TestImplementBlueprintGuardrails` (Fails initially due to missing sections)

- [x] T1.2: [GREEN] Implement Full Unabridged Guardrails in Commands and Constitution Blueprints
**Target:** `src/internal/agent/kit/commands/implement.yaml`
**Context:** [US-1], [US-2], [US-3]

**Action Steps:**
- Update `src/internal/agent/kit/commands/implement.yaml` to include the full, unabridged Pre-Coding, During Coding, and Post-Implementation guardrails in both the Orchestrator protocol and the `Mission Brief Envelope` template.
- Update `src/internal/agent/artifacts/constitution/engineering.yaml` with the complete 14-point AI coding constraints and test invariance blueprint.
- Ensure all formatting is strict YAML/Markdown compatible across all tool mappings.

**Acceptance Check:**
`go test -v ./src/internal/agent -run TestImplementBlueprintGuardrails` (Passes with exit code 0)

### Phase 2: Translation Adaptation & Living Constitution Sync

- [x] T2.1: [RED] Add Skill Adaptation and Translation Tests
**Target:** `src/internal/agent/translator_test.go`
**Context:** [US-1], [US-2]

**Action Steps:**
- Add test case verifying that `spf.implement` skill translation into markdown maintains the unabridged `NON-NEGOTIABLE WORKER GUARDRAILS` envelope.
- Verify that multi-tool translation (`claude`, `open-code`, `antigravity`) preserves the complete 14 guardrails directives without truncation.
- Ensure blueprint parsing and variable injections do not alter the integrity of the guardrails directives.

**Acceptance Check:**
`go test -v ./src/internal/agent -run TestSkillAdaptationGuardrails` (Fails initially if translator expectations differ)

- [x] T2.2: [GREEN] Synchronize Adapted Skills and Living Documentation
**Target:** `.agents/skills/spf-implement/SKILL.md`
**Context:** [US-1], [US-2], [US-3]

**Action Steps:**
- Synchronize `.agents/skills/spf-implement/SKILL.md` with the updated blueprint content containing all 14 principles.
- Update `.specforce/docs/engineering.md` to reflect the full 3-phase AI Coding Constraints and Test Invariance standards.
- Update `.specforce/docs/modules/agent-kit.md` with `[BR-KIT-07]` enforcing the full unabridged worker guardrails in execution envelopes.

**Acceptance Check:**
`go test ./...` (All tests across the entire repository pass with exit code 0)
