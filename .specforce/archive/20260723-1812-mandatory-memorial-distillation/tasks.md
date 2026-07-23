---
slug: 20260723-1812-mandatory-memorial-distillation
lens: Backend-heavy
---

# Implementation Roadmap: Mandatory Memorial Distillation in Archival Lifecycle

## 1. Execution Strategy
- **Gravity Order:** Kit Instruction Updates -> Go CLI Implementation & Metrics -> Unit & Integration Verification.

## 2. Tasks

### Phase 1: Instruction Protocol Reordering

- [x] T1.1: [DOCS] Reorder Archive Instruction Protocol
**Target:** `src/internal/agent/kit/instructions/archive.md`
**Context:** [US-1]

**Action Steps:**
- Update section headers in `archive.md` to place Knowledge Harvesting & Distillation as Step 5 before Archival Execution (Step 7).
- Change memory distillation from "Optional but Recommended" to a mandatory check step.
- Update guardrails to enforce distillation before spec archiving.

**Acceptance Check:**
Verify line ordering in `src/internal/agent/kit/instructions/archive.md` ensuring distillation precedes `specforce spec archive`.

### Phase 2: CLI Engine Metrics & Guidance

- [x] T2.1: [CODE] Add Memorial Fragment Counter to Archive Instructions
**Target:** `src/internal/cli/archive.go`
**Context:** [US-2]

**Action Steps:**
- Count total active fragments in `.specforce/memorial/` inside `HandleArchiveInstructions`.
- Inject active fragment count into CLI output in `printArchiveInstructions`.
- Add explicit warning log when active fragments exceed threshold (e.g. >= 5).

**Acceptance Check:**
Run `go test ./src/internal/cli/...` and run `specforce archive instructions` to inspect output metrics.

### Phase 3: Verification & Integration

- [x] T3.1: [TEST] Add Unit Test for Memorial Service Distill & Instructions Metrics
**Target:** `src/internal/project/memorial_test.go`
**Context:** [US-1, US-2]

**Action Steps:**
- Add unit test case verifying fragment consolidation when distillation is triggered.
- Verify zero-fragment edge case behavior.

**Acceptance Check:**
Run `go test ./src/internal/project/...` to ensure all tests pass cleanly.
