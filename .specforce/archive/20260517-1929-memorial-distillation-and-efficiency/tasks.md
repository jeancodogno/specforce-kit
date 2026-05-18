---
slug: 20260517-1929-memorial-distillation-and-efficiency
lens: Backend-heavy
---

# Implementation Roadmap: Memorial Distillation and Efficiency

## 1. Execution Strategy
- **Gravity Order:** Memorial Service logic -> CLI Playback -> Distillation Command -> Legacy Cleanup -> Instructions.

## 2. Tasks

### Phase 1: Memorial Service & Playback

- [x] T1.1: [CODE] Activate Memorial Playback in Instructions
**Target:** `src/internal/cli/archive.go`
**Context:** [US-1]

**Action Steps:**
- In `HandleArchiveInstructions`, replace the call to `e.getMemorialList()` with `memSvc.Consolidate(ctx, 10)`.
- Update `config.Context["MEMORIAL_FRAGMENTS"]` to use this consolidated string.

**Verification (TDD):**
Run `specforce archive instructions` and verify that the output contains the actual markdown content of the fragments instead of just a list of filenames.

- [x] T1.2: [CODE] Implement Distill method in MemorialService
**Target:** `src/internal/project/memorial.go`
**Context:** [US-2]

**Action Steps:**
- Add `Distill(ctx context.Context, slugs []string, summary string, author string) error` to the `MemorialService` interface and its implementation.
- The method should:
  1. Open/Create `.specforce/memorial/DISTILLED.md`.
  2. Append the summary with headers and timestamps.
  3. Delete the individual fragment files matching the provided slugs.

**Verification (TDD):**
Create a unit test in `src/internal/project/memorial_test.go` that records 3 fragments, calls `Distill` for 2 of them, and verifies that `DISTILLED.md` contains the summary and the 2 files are gone.

### Phase 2: Distillation CLI

- [x] T2.1: [CODE] Route Distillation Command
**Target:** `src/internal/cli/archive.go`
**Context:** [US-2]

**Action Steps:**
- Add a `distill` case to the `HandleArchive` switch statement.
- Implement `HandleArchiveDistill(ctx, ui, slugs, summary, author)`.

**Verification (TDD):**
Run `go build ./src/cmd/specforce` and verify the code compiles without errors.

- [x] T2.2: [CLI] Define 'archive distill' command and flags
**Target:** `src/internal/cli/cobra/archive.go`
**Context:** [US-2]

**Action Steps:**
- Define `archiveDistillCmd` with flags `--slug` (StringSlice) and `--summary` (String).
- Mark both flags as required.
- In `RunE`, call `executor.HandleArchiveDistill`.

**Verification (TDD):**
Run `specforce archive distill --help` and verify all flags are listed.

### Phase 3: Cleanup & Instructions

- [x] T3.1: [CODE] Automatic Legacy Cleanup
**Target:** `src/internal/cli/cli.go`
**Context:** [US-3]

**Action Steps:**
- In `handleNewInitFlow`, add a step to check for the existence of `.specforce/docs/memorial.md` and delete it.
- Also add this check to `HandleArchiveInstructions` to ensure immediate cleanup for existing users.

**Verification (TDD):**
Create a dummy `.specforce/docs/memorial.md`, run `specforce archive instructions`, and verify the file is deleted.

- [x] T3.2: [DOCS] Update Agent Instructions & Commands
**Target:** `Global Scope`
**Context:** [US-2]

**Action Steps:**
- Update `src/internal/agent/kit/instructions/archive.md` to include Step 9: Memory Distillation.
- Update `src/internal/agent/kit/commands/archive.yaml` to include the `distill` command mapping.

**Verification (TDD):**
Run `specforce archive instructions` and verify the new step appears in the protocol.

## 3. Pre-emptive Mitigations
- **Risk:** Deleting fragments before successfully writing to `DISTILLED.md` results in data loss. -> **Mitigation:** Ensure the file append operation is checked for errors before calling `os.Remove`.
