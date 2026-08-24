---
slug: 20260823-2158-batch-task-update-and-subagent-orchestration
lens: Integration
---

# Implementation Roadmap: Batch Task Update & Subagent Orchestration

## 1. Execution Strategy
- **Gravity Order:** Core Spec Engine & Multi-Task Mutation (TDD Red-Green) -> Spec Service & Deduplicated Hook Collection (TDD Red-Green) -> CLI Parsing & Cobra Layer -> Agent Kit Implementation Blueprint Updates.

## 2. Tasks

### Phase 1: Core Spec Engine & Hook Deduplication

- [x] T1.1: [RED] Add unit tests for batch task status updates and atomic validation
**Target:** `src/internal/spec/tasks_test.go`
**Context:** [US-1]

**Action Steps:**
- Add `TestUpdateTaskStatusesFile_Success` covering multiple tasks transitioning to `IN-PROGRESS` and `FINISHED` simultaneously.
- Add `TestUpdateTaskStatusesFile_AtomicFailure` verifying that if any task ID is nonexistent, the file remains unchanged and an error is returned.
- Add `TestUpdateTaskStatusesFile_MetadataTiming` verifying start and end session timing logs for all tasks in the batch.

**Acceptance Check:**
```bash
go test -v ./src/internal/spec -run "TestUpdateTaskStatusesFile"
```

- [x] T1.2: [GREEN] Implement batch task status updates and atomic file mutation
**Target:** `src/internal/spec/tasks.go`
**Context:** [US-1]

**Action Steps:**
- Implement `updateTaskStatusesFile(projectRoot, slug string, taskIDs []string, newStatus string) error` validating that all `taskIDs` exist before mutation.
- Update checklist markdown brackets (`- [x]`, `- [/]`, `- [ ]`) and `**State:** [FINISHED]` markers for all matched task blocks in a single file write.
- Update `spec.yaml` metadata session timers for all supplied task IDs upon transition.

**Acceptance Check:**
```bash
go test -v ./src/internal/spec -run "TestUpdateTaskStatusesFile"
```

- [x] T1.3: [RED] Add unit tests for deduplicated batch hook execution
**Target:** `src/internal/spec/implementation_test.go`
**Context:** [US-2]

**Action Steps:**
- Add `TestService_UpdateTaskStatus_BatchHooks` verifying `OnTaskFinished` executes exactly once for a multi-task batch.
- Add `TestService_UpdateTaskStatus_PhaseAndSpecHooks` verifying `OnPhaseFinished` and `OnAllTasksFinished` triggers when the batch includes final tasks.
- Add `TestService_UpdateTaskStatus_HookFailureRollback` verifying task file is not modified when a hook fails with exit code 1.

**Acceptance Check:**
```bash
go test -v ./src/internal/spec -run "TestService_UpdateTaskStatus_Batch"
```

- [x] T1.4: [GREEN] Implement batch hook collection and execution in Spec Service
**Target:** `src/internal/spec/service.go`
**Context:** [US-2]

**Action Steps:**
- Update `UpdateTaskStatus(ctx context.Context, projectRoot, slug string, taskIDs []string, status string) error` signature.
- Implement `collectHooksForBatch(ctx context.Context, projectRoot, slug string, taskIDs []string, config *core.ProjectConfig) []string` deduplicating hook commands across all task IDs.
- Execute deduplicated hooks via `core.ExecuteHooks` before delegating to `updateTaskStatusesFile`.

**Acceptance Check:**
```bash
go test -v ./src/internal/spec -run "TestService_UpdateTaskStatus"
```

### Phase 2: CLI Layer Integration & Agent Kit Directives

- [x] T2.1: [RED] Add CLI integration tests for multi-task flag parsing and execution
**Target:** `src/internal/cli/implementation_test.go`
**Context:** [US-1]

**Action Steps:**
- Add test case verifying `--task T1.1,T1.2` comma-separated flag parsing.
- Add test case verifying repeated flags `--task T1.1 --task T1.2`.
- Add test case verifying consolidated stdout feedback formatting (`[OK] Tasks T1.1, T1.2 status updated to finished in <slug>`).

**Acceptance Check:**
```bash
go test -v ./src/internal/cli -run "TestImplementationUpdate_Batch"
```

- [x] T2.2: [GREEN] Update Cobra commands and CLI executor for batch task updates
**Target:** `src/internal/cli/cobra/implementation.go`
**Context:** [US-1]

**Action Steps:**
- Replace `taskId string` with `taskIDs []string` using `StringSliceVar` on `implementationUpdateCmd`.
- Update `HandleImplementationUpdate` in `src/internal/cli/implementation.go` to accept `taskIDs []string` and format consolidated output.
- Forward batch task slice from Cobra handler to `executor.HandleImplementationUpdate`.

**Acceptance Check:**
```bash
go test -v ./src/internal/cli/...
```

- [x] T2.3: [GREEN] Update Agent Kit implement blueprint with subagent rules and batch CLI syntax
**Target:** `src/internal/agent/kit/commands/implement.yaml`
**Context:** [US-3]

**Action Steps:**
- Update CLI execution snippets in `implement.yaml` to showcase batch updates (`specforce implementation update <slug> --task T1.1,T1.2 --status in-progress` and `--status finished`).
- Add subagent sizing directives (1 subagent for small tasks, 2-3 sweet spot for standard tasks, up to 4 for complex tasks) and mandate subagent usage whenever available.
- Add explicit safety prohibition forbidding subagents from directly editing `tasks.md` or altering task states.

**Acceptance Check:**
```bash
go test -v ./src/internal/agent/...
```

## 3. Pre-emptive Mitigations
- **Risk:** Backward compatibility breakage if single string `--task T1.1` is passed -> **Mitigation:** Cobra's `StringSliceVar` natively accepts single values, comma-separated lists, and repeated flags without breaking existing single-task invocations.
