---
slug: 20261023-1159-persist-qa-report-and-red-task-hooks
lens: Backend-heavy
---

# Implementation Roadmap: Persist QA Report and Scoped Hooks for RED Tasks

## 1. Execution Strategy
- **Gravity Order:** metadata and QA domain model -> status exposure and archive gate in the spec package -> RED hook scope in the spec package -> CLI wiring -> kit blueprints -> black-box verification. Every command runs from the repository root; Go commands run inside `src`. Each phase pairs RED tests with GREEN code.
- **Module delta coverage:** `agent-kit` items `[NEW:qa-record-flow]`, `[NEW:archive-uses-qa]`, `[NEW:red-title-marker]`, `[NEW:us-qa-record-flow]`, `[NEW:us-archive-uses-qa]`, `[NEW:us-red-title-marker]` are covered by Phase 5 and Phase 6; `spec-management` items `[NEW:qa-persist]`, `[NEW:qa-status-exposed]`, `[NEW:qa-archive-gate]`, `[NEW:red-hook-scope]`, `[NEW:us-qa-persist]`, `[NEW:us-qa-archive-gate]`, `[NEW:us-red-hook-scope]` are covered by Phases 1 to 4 and Phase 6.
- **Hook note:** while this roadmap is implemented the RED relaxation does not exist yet, so run `PATH="$HOME/go/bin:$PATH"` with the project hooks and finish each RED task together with its GREEN pair in the same update command.

## 2. Tasks

### Phase 1: QA Domain Model and Atomic Metadata

- [x] T1.1: [RED] Add failing tests for QA input validation, recording and atomic save
**Target:** `src/internal/spec/qa_test.go`
**Context:** [US-1]

**Action Steps:**
- Create `qa_test.go` with table tests for `ParseQAInput` and `ValidateQAInput` covering a valid matrix, unknown verdict (error names the criterion `id`), zero criteria, 201 criteria, empty `scenario`, a field over 2000 characters, an unknown JSON field and empty input.
- Add `TestRecordQA_DerivesStatus` asserting FAILED when any criterion is FAILED and PASSED when all are PASSED or NOT-VERIFIABLE, plus `run` starting at 1 and increasing to 2 on a second call that fully replaces the first report.
- Add `TestRecordQA_PreservesMetadata` seeding `time_logs`, `refinement` and `delta` in `spec.yaml` and asserting they are identical after recording, `TestRecordQA_RejectsMissingAndArchivedSpec` asserting an error and no file written, and `TestSaveMetadata_AtomicMode` asserting mode `0600` and no leftover temp file.

**Acceptance Check:**
Run `cd src && go test ./internal/spec/ -run 'TestRecordQA|TestSaveMetadata_AtomicMode|TestParseQAInput|TestValidateQAInput'` and confirm it FAILS to compile or fails because `RecordQA`, `ParseQAInput` and `ValidateQAInput` do not exist yet.

- [x] T1.2: [GREEN] Implement the QA domain model and atomic metadata write
**Target:** `src/internal/spec/qa.go`
**Context:** [US-1]

**Action Steps:**
- Create `qa.go` defining `QAMetadata`, `QACriterion`, `QAFinding`, `QAInput`, the status and verdict constants, `ParseQAInput` (strict decoding with `DisallowUnknownFields`, 1 MiB limit), `ValidateQAInput` (limits from the design) and `RecordQA` (rejects missing or archived specs, derives status, increments `run`, sets UTC `recorded_at`, saves through `SaveMetadata`).
- Add `QA *QAMetadata` with `json:"qa,omitempty" yaml:"qa,omitempty"` to `Metadata` in `src/internal/spec/metadata.go`.
- Change `SaveMetadata` in `src/internal/spec/metadata.go` to write a temp file in the same directory with mode `0600` and `os.Rename` it over `spec.yaml`, removing the temp file on failure.

**Acceptance Check:**
Run `cd src && go test ./internal/spec/` and confirm all tests in the package pass, including the tests from T1.1.

### Phase 2: Status Exposure and Archive Gate

- [x] T2.1: [RED] Add failing tests for QA status summary and the archive QA gate
**Target:** `src/internal/spec/archive_test.go`
**Context:** [US-2], [US-3]

**Action Steps:**
- Add `TestCheckQAGate` in `archive_test.go` with subtests `Missing` (no `qa` block returns `missing`), `Passed`, `Failed` (returns the failed criteria `id` and `scenario`) and `Unreadable` (malformed `qa` block returns an error).
- Add `TestGetImplementationStatus_QA` in `src/internal/spec/service_test.go` asserting the report carries `QA` with counts, failed and not-verifiable criteria and findings when a report is stored, `QA == nil` when none is stored, and `QAError` set while tasks are still returned when the block is malformed.
- Add `TestSummarizeQA` in `qa_test.go` asserting the counts and lists produced from a stored report.

**Acceptance Check:**
Run `cd src && go test ./internal/spec/ -run 'TestCheckQAGate|TestGetImplementationStatus_QA|TestSummarizeQA'` and confirm it FAILS because `CheckQAGate`, `SummarizeQA` and the report fields do not exist yet.

- [x] T2.2: [GREEN] Expose the QA summary in the implementation report
**Target:** `src/internal/spec/implementation.go`
**Context:** [US-2]

**Action Steps:**
- Add `QA *QASummary` with `json:"qa,omitempty"` and `QAError string` with `json:"qa_error,omitempty"` to `ImplementationReport`.
- Implement `QASummary`, `SummarizeQA` in `src/internal/spec/qa.go` producing status, run, `recorded_at`, counts, `failed_criteria`, `not_verifiable_criteria` and `findings`.
- In `src/internal/spec/service.go` `GetImplementationStatus`, load metadata after parsing tasks and fill `QA`, or set `QAError` when the `qa` block cannot be parsed, without failing the call.

**Acceptance Check:**
Run `cd src && go test ./internal/spec/ -run 'TestGetImplementationStatus_QA|TestSummarizeQA'` and confirm it passes.

- [x] T2.3: [GREEN] Implement the archive QA gate check
**Target:** `src/internal/spec/archive.go`
**Context:** [US-3]

**Action Steps:**
- Add `QAGateResult` with fields `State` (`missing`, `passed`, `failed`) and `Failed []QACriterion`.
- Implement `CheckQAGate(projectRoot, slug string) (QAGateResult, error)` using `LoadMetadata` after `ResolveSlug`, returning an error when the `qa` block is unreadable.
- Leave `ArchiveSpec` unchanged so the gate stays in the CLI layer like the delta gate.

**Acceptance Check:**
Run `cd src && go test ./internal/spec/ -run TestCheckQAGate` and confirm it passes.

### Phase 3: Scoped Hook Relaxation for RED Tasks

- [x] T3.1: [RED] Add failing tests for RED task hook skipping
**Target:** `src/internal/spec/service_test.go`
**Context:** [US-4]

**Action Steps:**
- Add `TestIsRedTask` covering `[RED] Add tests` (true), `Add tests [RED]` (true), `[red] x` (false), `[GREEN] x` (false) and an empty title (false).
- Add `TestUpdateTaskStatusDetailed_RedHooks` using a recording hook provider with `on_task_finished`, `on_phase_finished` and `on_all_tasks_finished` entries, asserting: a RED-only batch skips and returns the `on_task_finished` hooks; a GREEN task runs them; a mixed batch runs them once and skips nothing; a RED task last in its phase still runs the phase hook; a RED task with no `on_task_finished` configured returns an empty skipped list; a RED task moved to in-progress runs no hooks.
- Assert `UpdateTaskStatus` still returns only an error and behaves like the detailed method.

**Acceptance Check:**
Run `cd src && go test ./internal/spec/ -run 'TestIsRedTask|TestUpdateTaskStatusDetailed_RedHooks'` and confirm it FAILS because `IsRedTask` and `UpdateTaskStatusDetailed` do not exist yet.

- [x] T3.2: [GREEN] Implement RED detection and scoped hook skipping
**Target:** `src/internal/spec/service.go`
**Context:** [US-4]

**Action Steps:**
- Add `IsRedTask(title string) bool` in `src/internal/spec/implementation.go` returning true when the title contains the exact case-sensitive token `[RED]`.
- Extend `collectHooksForBatch` in `src/internal/spec/service.go` to look up each task title from `ParseTasks`, and when every task in the batch is RED and `on_task_finished` is non-empty, omit those hooks and return them as the skipped list; keep phase and all-tasks collection and deduplication unchanged.
- Add `UpdateTaskStatusDetailed(ctx, projectRoot, slug, taskIDs, status) ([]string, error)` holding the existing flow and turn `UpdateTaskStatus` into a wrapper that discards the skipped list.

**Acceptance Check:**
Run `cd src && go test ./internal/spec/` and confirm all tests pass, including those from T3.1.

### Phase 4: CLI Wiring

- [x] T4.1: [RED] Add failing CLI tests for the QA command, the RED notice and the archive QA gate
**Target:** `src/internal/cli/implementation_test.go`
**Context:** [US-1], [US-3], [US-4]

**Action Steps:**
- Add `TestHandleImplementationQARecord` in `implementation_test.go` feeding JSON through a `strings.Reader`, asserting the success line format, `spec.yaml` content, refusal messages for invalid input and that an invalid call leaves a stored report untouched.
- Add `TestHandleImplementationUpdate_RedNotice` asserting the warning `RED task: skipped per-task hooks: <hooks>` is emitted for a RED-only finish and not for a GREEN finish.
- Add `TestHandleSpecArchive_QAGate` in `src/internal/cli/spec_test.go` asserting FAILED refuses with and without `force`, `missing` warns and archives, `passed` archives, and add a subcommand registration assertion for `implementation qa record` in `src/internal/cli/cobra/root_test.go`.

**Acceptance Check:**
Run `cd src && go test ./internal/cli/... -run 'TestHandleImplementationQARecord|TestHandleImplementationUpdate_RedNotice|TestHandleSpecArchive_QAGate|TestRootCommands'` and confirm it FAILS because the handlers and command do not exist yet.

- [x] T4.2: [GREEN] Add the qa record command and the RED notice
**Target:** `src/internal/cli/implementation.go`
**Context:** [US-1], [US-4]

**Action Steps:**
- Add `HandleImplementationQA` (subcommand switch with `record`) and `HandleImplementationQARecord(ctx, ui, slug, in io.Reader)` that parses stdin with `spec.ParseQAInput`, calls `spec.RecordQA` and prints the `[OK] QA report recorded for ...` line.
- Change `HandleImplementationUpdate` to call `UpdateTaskStatusDetailed` and, when the skipped list is non-empty, call `ui.Warn("RED task: skipped per-task hooks: " + strings.Join(skipped, "; "))` before the `[OK]` line.
- Register `implementationQACmd` with the `record [slug]` subcommand in `src/internal/cli/cobra/implementation.go` reading `os.Stdin`.

**Acceptance Check:**
Run `cd src && go test ./internal/cli/... -run 'TestHandleImplementationQARecord|TestHandleImplementationUpdate_RedNotice|TestRootCommands'` and confirm it passes.

- [x] T4.3: [GREEN] Enforce the QA gate on archive
**Target:** `src/internal/cli/spec.go`
**Context:** [US-3]

**Action Steps:**
- Add `enforceQAGate(ctx, ui, projectRoot, slug)` calling `spec.CheckQAGate`, returning an error that lists each failed criterion as `<id> - <scenario>` with the re-record instruction for `failed`, a `ui.Warn("no QA report recorded for <slug>; archiving without QA evidence")` for `missing`, and an error for an unreadable block.
- Call `enforceQAGate` in `HandleSpecArchive` after `enforceDeltaGate` and before `spec.ArchiveSpec`, independent of the `force` flag.
- Keep the error text free of the word "force" so it never suggests an override.

**Acceptance Check:**
Run `cd src && go test ./internal/cli/... -run TestHandleSpecArchive_QAGate` and confirm it passes.

- [x] T4.4: [GREEN] Show the QA line in the status view
**Target:** `src/internal/tui/implementation.go`
**Context:** [US-2]

**Action Steps:**
- Update `RenderImplementationStatus` to print `QA: <STATUS> (P passed, F failed, N not-verifiable, K findings)` using the error color for FAILED and the success color for PASSED.
- Print a muted `QA: not recorded` when `report.QA` is nil and no `QAError` exists, and print `QA: unreadable (<QAError>)` when `QAError` is set.
- Keep existing task rendering untouched.

**Acceptance Check:**
Run `cd src && go build ./... && go test ./internal/tui/...` and confirm both succeed.

### Phase 5: Kit Blueprints

- [x] T5.1: [RED] Add failing manifest tests for the new blueprint rules
**Target:** `src/internal/agent/kit_manifests_test.go`
**Context:** [US-5]

**Action Steps:**
- Add `TestKitManifests_QAPersistenceAndRedMarker` with subtests `Implement`, `Archive`, `ArchiveSkill` and `Spec`, reusing the existing `assertFileContains` helper.
- `Implement`: assert `kit/skills/spf-implement/SKILL.yaml` contains `specforce implementation qa record`, the phrase that the qa-tester never writes the report and the phrase about repeating after every correction cycle.
- `Archive`, `ArchiveSkill` and `Spec`: assert `kit/instructions/archive.md` contains `implementation status`, `FAILED`, `no QA report`; assert `kit/skills/spf-archive/SKILL.yaml` contains `persisted` and no longer contains `present in the conversation`; assert `kit/skills/spf-spec/SKILL.yaml` contains `[RED]` in a title rule.

**Acceptance Check:**
Run `cd src && go test ./internal/agent/ -run TestKitManifests_QAPersistenceAndRedMarker` and confirm it FAILS with missing-directive errors for every subtest.

- [x] T5.2: [GREEN] Make archival read the persisted QA state
**Target:** `src/internal/agent/kit/instructions/archive.md`
**Context:** [US-3], [US-5]

**Action Steps:**
- In the verification step, add reading the `qa` object from `specforce implementation status <slug> --json`: when `qa.status` is `FAILED`, stop before any file change, list each failed criterion and point to fixing, re-validating and re-recording; when `qa` is absent, state `no QA report recorded` and continue.
- Add the instruction that next-step suggestions come from `qa.findings` and `qa.not_verifiable_criteria`, never from the conversation.
- Update `src/internal/agent/kit/skills/spf-archive/SKILL.yaml` step 4 to derive suggestions from the persisted QA state, replacing the wording about findings present in the conversation, keeping "Suggest Next Steps & Follow-up Specs" and the `/spf:discovery` and `/spf:spec` phrasing.

**Acceptance Check:**
Run `cd src && go test ./internal/agent/ -run 'TestKitManifests_QAPersistenceAndRedMarker/(Archive|ArchiveSkill)'` and confirm it passes.

- [x] T5.3: [GREEN] Direct the orchestrator to record QA and planners to tag RED tasks
**Target:** `src/internal/agent/kit/skills/spf-implement/SKILL.yaml`
**Context:** [US-5]

**Action Steps:**
- After the independent qa-tester step, add the directive that the orchestrator pipes the report as JSON to `specforce implementation qa record <slug>` before declaring validation complete, repeats it after every correction cycle, and that the qa-tester never writes the report.
- Add the note that tasks titled `[RED]` skip only per-task hooks with a visible notice and that the GREEN pair proves lint and tests.
- In `src/internal/agent/kit/skills/spf-spec/SKILL.yaml` add the roadmap rule that RED tasks carry the `[RED]` token in their title, keeping the file within its existing line-count assertion.

**Acceptance Check:**
Run `cd src && go test ./internal/agent/` and confirm all tests in the package pass, including `TestKitManifests_QAPersistenceAndRedMarker`.

### Phase 6: Black-box Verification

- [x] T6.1: [VERIFY] Record a QA report and read it back (US-1, US-2)
**Target:** `/tmp/spf-bb implementation qa record demo`
**Context:** [US-1], [US-2]

**Action Steps:**
- Build the binary with `cd src && go build -o /tmp/spf-bb ./cmd/specforce` and create a throwaway project containing `.specforce/specs/demo/spec.yaml` and a minimal `tasks.md`.
- Run `echo '{"criteria":[{"id":"US-1.1","scenario":"s","evidence":"e","verdict":"PASSED"},{"id":"US-1.2","scenario":"s2","evidence":"e2","verdict":"NOT-VERIFIABLE"}],"findings":[{"title":"f"}]}' | /tmp/spf-bb implementation qa record demo` inside it.
- Run `/tmp/spf-bb implementation status demo --json` and compare the `qa` object.

**Acceptance Check:**
Run `cd /tmp/spf-bb-proj && echo '{"criteria":[{"id":"US-1.1","scenario":"s","evidence":"e","verdict":"PASSED"}],"findings":[]}' | /tmp/spf-bb implementation qa record demo && /tmp/spf-bb implementation status demo --json | grep -q '"status": "PASSED"'` and confirm exit code 0.

- [x] T6.2: [VERIFY] Invalid report is refused and a stored one survives (US-1)
**Target:** `/tmp/spf-bb implementation qa record demo`
**Context:** [US-1]

**Action Steps:**
- With the report from T6.1 stored, pipe a report containing the verdict `MAYBE` to the record command.
- Confirm the command exits non-zero and names the criterion.
- Run `implementation status demo --json` and confirm the earlier PASSED report is unchanged.

**Acceptance Check:**
Run `cd /tmp/spf-bb-proj && ! (echo '{"criteria":[{"id":"US-9","scenario":"s","evidence":"e","verdict":"MAYBE"}],"findings":[]}' | /tmp/spf-bb implementation qa record demo) && /tmp/spf-bb implementation status demo --json | grep -q '"status": "PASSED"'` and confirm exit code 0.

- [x] T6.3: [VERIFY] FAILED QA blocks archive even with force (US-3)
**Target:** `/tmp/spf-bb spec archive demo --force`
**Context:** [US-3]

**Action Steps:**
- Record a report with one FAILED criterion for `demo`.
- Run `spec archive demo --force` and confirm it exits non-zero, lists the criterion and leaves `.specforce/specs/demo` in place.
- Record a PASSED report and confirm the QA gate no longer refuses.

**Acceptance Check:**
Run `cd /tmp/spf-bb-proj && echo '{"criteria":[{"id":"US-1.1","scenario":"s","evidence":"e","verdict":"FAILED"}],"findings":[]}' | /tmp/spf-bb implementation qa record demo && ! /tmp/spf-bb spec archive demo --force && test -d .specforce/specs/demo` and confirm exit code 0.

- [x] T6.4: [VERIFY] RED task finishes with a visible notice while GREEN still runs hooks (US-4)
**Target:** `/tmp/spf-bb implementation update demo --task T1.1 --status finished`
**Context:** [US-4]

**Action Steps:**
- In the throwaway project configure `on_task_finished` with a command that always fails (`false`) and a `tasks.md` with `T1.1` titled `[RED] tests` and `T1.2` titled `[GREEN] code`.
- Finish `T1.1` and confirm exit code 0 with the skip notice naming `false`.
- Finish `T1.2` and confirm it is blocked by the hook.

**Acceptance Check:**
Run `cd /tmp/spf-bb-proj && /tmp/spf-bb implementation update demo --task T1.1 --status finished 2>&1 | grep -q 'RED task: skipped per-task hooks: false' && ! /tmp/spf-bb implementation update demo --task T1.2 --status finished` and confirm exit code 0.

- [x] T6.5: [VERIFY] Blueprints carry the new flow and the full suite is green (US-5)
**Target:** `cd src && go test ./...`
**Context:** [US-5]

**Action Steps:**
- Run the full Go test suite from `src`.
- Run `PATH="$HOME/go/bin:$PATH" golangci-lint run ./...` and confirm no new findings.
- Run `specforce archive instructions | grep -q 'no QA report'` with the freshly built binary.

**Acceptance Check:**
Run `cd src && go test ./... && /tmp/spf-bb archive instructions | grep -q 'no QA report'` and confirm exit code 0.

## 3. Pre-emptive Mitigations
- **Risk:** Atomic rename in `SaveMetadata` changes behavior for every metadata writer -> **Mitigation:** T1.1 asserts mode `0600`, no temp leftovers and unchanged `time_logs`, `refinement` and `delta`; the full package suite in T1.2 must pass.
- **Risk:** Black-box commands leave files in `/tmp` -> **Mitigation:** every command targets `/tmp/spf-bb-proj` and `/tmp/spf-bb` only, recreated by T6.1.
