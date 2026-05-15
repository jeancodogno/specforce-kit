---
slug: 20260513-1023-optimize-archive-memorial
lens: Backend-heavy
---

# Implementation Roadmap: Optimize Archive Memorial

## 1. Execution Strategy
- **Gravity Order:** Domain Extensions -> CLI Executor Logic -> Cobra Command Registration -> Kit Instruction Updates -> Verification.

## 2. Tasks

### Phase 1: Domain & CLI Handler Expansion

- [x] T1.1: [DOMAIN] Add 'Context' fragment type
**Target:** `src/internal/project/memorial.go`
**Context:** [US-1]

**Action Steps:**
- Add `FragmentContext FragmentType = "Context"` to the `FragmentType` constants.
- Update the `Record` method to support this type if necessary (currently it handles string conversion, so it should be transparent).

**Verification (TDD):**
Run `go test ./src/internal/project/...` and verify that a memorial fragment with type 'Context' is correctly serialized.

- [x] T1.2: [CODE] Implement HandleArchiveMemorial in Executor
**Target:** `src/internal/cli/archive.go`
**Context:** [US-1]

**Action Steps:**
- Implement `HandleArchiveMemorial(ctx context.Context, ui core.UI, slug string, ftype string, title string, content string, author string) error`.
- Initialize `NewMemorialService(e.projectRoot)`.
- Create `project.Fragment` with the provided data and `time.Now()`.
- Call `memSvc.Record(ctx, fragment)`.
- Use `ui.Success` to report the created file path.

**Verification (TDD):**
Verify the method compiles and a manual test call (via temporary test or CLI) successfully creates the file.

- [x] T1.3: [CODE] Route 'memorial' subcommand in HandleArchive
**Target:** `src/internal/cli/archive.go`
**Context:** [US-1]

**Action Steps:**
- Update `HandleArchive` switch statement to include a `memorial` case.
- Parse basic args for positional use if needed (though primary use is via Cobra flags).

**Verification (TDD):**
Assert that calling the executor with "memorial" routes to the new handler.

### Phase 2: Cobra Integration

- [x] T2.1: [CLI] Define 'archive memorial' command and flags
**Target:** `src/internal/cli/cobra/archive.go`
**Context:** [US-1]

**Action Steps:**
- Define `archiveMemorialCmd` with positional `<slug>` arg.
- Add required string flags: `--type`, `--title`, `--content`.
- Add optional string flag: `--author`.
- In `RunE`, extract flags and call `executor.HandleArchiveMemorial`.

**Verification (TDD):**
Run `specforce archive memorial --help` and verify all flags and requirements are listed.

### Phase 3: Agent Protocol Optimization

- [x] T3.1: [CONFIG] Update Archival Blueprint Instructions
**Target:** `src/internal/agent/kit/instructions/archive.md`
**Context:** [US-2]

**Action Steps:**
- Replace manual file creation instructions in "Step 5: Knowledge Harvesting" with the `specforce archive memorial` command.
- Remove redundant instructions for `date` and `ls`.

**Verification (TDD):**
`specforce archive instructions` output should show the new CLI-based workflow.

- [x] T3.2: [CONFIG] Update Archival Command Blueprint
**Target:** `src/internal/agent/kit/commands/archive.yaml`
**Context:** [US-2]

**Action Steps:**
- Update the blueprint description to reflect the automated memorial capability.

**Verification (TDD):**
Check the file content to ensure it aligns with the new workflow.

### Phase 4: Verification & Cleanup

- [x] T4.1: [TEST] End-to-End Verification
**Target:** `Global Scope`
**Context:** [US-1]

**Action Steps:**
- Build and run: `specforce archive memorial "test-slug" --type lesson --title "Test" --content "Content"`.
- Verify file exists in `.specforce/memorial/`.

**Verification (TDD):**
`ls .specforce/memorial/*test-slug.md` should find the file.

- [x] T4.2: [TEST] Final Integration Check
**Target:** `Global Scope`
**Context:** [US-2]

**Action Steps:**
- Run `specforce archive instructions` and verify the prompt contains the new command.

**Verification (TDD):**
Verify the rendered text in the terminal matches the updated blueprint.

