---
slug: 20260724-1509-fix-archival-instruction-gap
lens: Backend-heavy
---

# Implementation Roadmap: Fix Archival Instruction Gap

## 1. Execution Strategy
- **Gravity Order:** Instruction Documentation Update (`archive.md`) -> CLI Output Update (`archive.go`) -> Verification via CLI execution and test suite.

## 2. Tasks

### Phase 1: Instruction & CLI Guidance Enhancements

- [x] T1.1: [DOC] Enhance archive.md with Dual-Step Mandatory Warnings
**Target:** `src/internal/agent/kit/instructions/archive.md`
**Context:** [FIX-1, FIX-2]

**Action Steps:**
- Add a prominent `CRITICAL DUAL-STEP REQUIREMENT` banner to Step 8 of `archive.md`.
- Explicitly state that `specforce archive memorial` and `specforce archive distill` manage global memory, whereas `specforce spec archive <slug>` is required to update spec lifecycle status from `active` to `archived`.
- Add a mandatory Guardrail in `archive.md` stating that completing the archival flow without running `specforce spec archive <slug>` is a protocol violation.

**Acceptance Check:**
- Inspect `src/internal/agent/kit/instructions/archive.md` and verify Step 8 and Guardrails contain explicit dual-step warnings.

- [x] T1.2: [CODE] Update CLI printArchiveInstructions Output
**Target:** `src/internal/cli/archive.go`
**Context:** [FIX-3]

**Action Steps:**
- Update `printArchiveInstructions` function in `src/internal/cli/archive.go`.
- Add a dedicated section outputting `## IMPORTANT: Archival Scopes & Command Separation` to clearly distinguish `specforce archive ...` (Memory) from `specforce spec archive <slug>` (Spec State).
- Ensure output is clean, readable, and properly formatted in terminal logs.

**Acceptance Check:**
- Run `go run main.go archive instructions` (or `specforce archive instructions`) and verify the output contains the scope separation banner.
