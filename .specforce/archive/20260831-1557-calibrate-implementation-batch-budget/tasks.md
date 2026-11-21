---
slug: 20260831-1557-calibrate-implementation-batch-budget
lens: Integration
---

# Implementation Roadmap: Calibrate Implementation Batch Budget and Global Subagent Sizing

## 1. Execution Strategy
- **Phase 1 (RED & GREEN):** Write test assertions in `kit_manifests_test.go` validating the global implementation batch budget and subagent sizing guidance (RED), then update `src/internal/agent/kit/commands/implement.yaml` to replace the micro-batching rule with a global 2-3 batches budget (max 4) across the entire implementation roadmap (GREEN).
- **Phase 2 (Integration & Regression Verification):** Run full test suites across the package and project to ensure blueprint consistency and no regressions.

## 2. Tasks

### Phase 1: Blueprint Calibration & Test Verification

- [x] T1.1: [TEST] [RED] Add unit test assertions for global batch budget and implementation worker sizing
**Target:** `src/internal/agent/kit_manifests_test.go`
**Context:** [US-1], [US-2]

**Action Steps:**
- Add `TestImplementGlobalBatchBudgetAndSizing` in `kit_manifests_test.go`.
- Define assertions checking for global batch budget directives (e.g., "Global Batch Budget", "2 to 3 Batches (Recommended Sweet Spot)", "Maximum of 4 Batches").
- Define assertions validating that worker sizing refers to the entire implementation lifecycle (2 to 3 implementation subagents, max 4, plus 1 QA specialist).

**Acceptance Check:**
Run `go test -v ./src/internal/agent/ -run TestImplementGlobalBatchBudgetAndSizing` (fails RED prior to implement.yaml changes).

- [x] T1.2: [CODE] [GREEN] Update implement.yaml with global batch budget and revised sizing directives
**Target:** `src/internal/agent/kit/commands/implement.yaml`
**Context:** [US-1], [US-2]

**Action Steps:**
- Replace the "Group up to 3 sequential tasks" rule in Step 1 with global roadmap partitioning guidelines (1 batch for small roadmaps, 2-3 batches sweet spot, max 4 batches for complex features).
- Update Subagent Sizing Guidance to clearly state that the 2 to 3 subagents guideline applies to the entire implementation lifecycle (1 dedicated subagent per batch + 1 final QA subagent), rather than per individual batch.
- Retain the final Step 4 Quality Assurance Specialist persona subagent for global validation.

**Acceptance Check:**
Run `go test -v ./src/internal/agent/ -run TestImplementGlobalBatchBudgetAndSizing` (must pass GREEN with exit code 0).

### Phase 2: Integration & Regression Verification

- [x] T2.1: [VERIFY] Run package test suite and full project validation
**Target:** `src/internal/agent/...`
**Context:** [US-1], [US-2]

**Action Steps:**
- Execute `go test ./src/internal/agent/...` to verify all kit manifest and translator tests pass.
- Execute `go test ./...` across the entire codebase.
- Check git diff for cleanliness.

**Acceptance Check:**
Run `go test -v ./src/internal/agent/...` (exit code 0).
