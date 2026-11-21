# Implementation Roadmap: Parallel Task Support

## 1. Execution Strategy
- **Gravity Order:** Data Layer -> Parser -> Logic Engine -> Agent Artifacts -> CLI/JSON Output.
- **Verification Strategy:** TDD-focused approach. Unit tests for parser and validation logic will be written BEFORE or alongside the implementation to ensure all ACs (Conflict, Symmetry, Boundaries) are met.

## 2. Tasks

### Phase 1: Go Data Structures & Parser

- [x] T1.1: [MODEL] Update `ImplementationTask` struct
**Target:** `src/internal/spec/implementation.go`
**Context:** [US-1, US-8]

**Action Steps:**
- Add `ParallelWith []string `json:"parallel_with"` ` to `ImplementationTask` struct.
- Add `IsParallel bool `json:"is_parallel"` ` to `ImplementationTask` struct.
- Ensure JSON tags match the required output schema.

**Acceptance Check:**
Write a unit test in `src/internal/spec/implementation_test.go` asserting the struct fields exist and are serializable to JSON.

- [x] T1.2: [PARSER] Implement Parallel Metadata Extraction
**Target:** `src/internal/spec/implementation.go`
**Context:** [US-1, US-2, US-3]

**Action Steps:**
- Define regex: `\*\*Parallel With:\*\* (T[\d.]+(?:,\s*T[\d.]+)*)`.
- In `parseTaskBlock`, use the regex to extract task IDs from the block content.
- Split the matched string by comma/whitespace to populate `ParallelWith`.
- Set `IsParallel = true` if `ParallelWith` is not empty.

**Acceptance Check:**
Unit test in `src/internal/spec/implementation_test.go` using a raw string block containing `**Parallel With: T1.1, T1.2**` and asserting the slice contains `["T1.1", "T1.2"]`.

- [x] T1.3: [TEST] Verify Parser with Multiline Metadata
**Target:** `src/internal/spec/implementation_test.go`
**Context:** [US-2]

**Action Steps:**
- Add a test case to `TestParseTasks` with a complete `tasks.md` content containing parallel tasks.
- Assert that `ParallelWith` is correctly populated regardless of its position in the task description.

**Acceptance Check:**
`go test ./src/internal/spec/ -run TestParseTasks` passes.

### Phase 2: Validation Logic (Integrity & Isolation)

- [x] T2.1: [VALIDATION] Track Parallelism in `taskBlock`
**Target:** `src/internal/spec/tasks.go`
**Context:** [US-4, US-5, US-6, US-7]

**Action Steps:**
- Add `parallelWith []string` and `isParallel bool` to `taskBlock` struct.
- Update `updateTaskBlockState` to recognize the `**Parallel With:**` line and store the peer IDs.

**Acceptance Check:**
Internal validation state correctly tracks the metadata during the line-by-line scanning process.

- [x] T2.2: [VALIDATION] Implement Reference & Phase Boundary Checks
**Target:** `src/internal/spec/tasks.go`
**Context:** [US-6, US-7]

**Action Steps:**
- In `ValidateTasks`, after the main loop, iterate through all collected tasks.
- For each task with `ParallelWith`:
  - Verify each peer ID exists in the global task map.
  - Verify each peer ID belongs to the same `Phase`.
- Append descriptive error messages to `state.errors` on failure.

**Acceptance Check:**
Unit test in `src/internal/spec/tasks_validation_test.go` asserting failure when referencing `T9.9` (missing) or a task from another phase.

- [x] T2.3: [VALIDATION] Implement Symmetry Enforcement
**Target:** `src/internal/spec/tasks.go`
**Context:** [US-5]

**Action Steps:**
- Implement logic to ensure if T1.1 is parallel with T1.2, T1.2 is implicitly parallel with T1.1.
- This can be done by normalizing the relationship in a global adjacency map during validation.

**Acceptance Check:**
Unit test asserting that `spec status` passes even if only one side of the relationship is explicitly declared.

- [x] T2.4: [VALIDATION] Implement Target Isolation (Conflict Detection)
**Target:** `src/internal/spec/tasks.go`
**Context:** [US-4, US-9]

**Action Steps:**
- For every parallel group, compare the `Target` field of each task.
- If any two tasks in the same parallel group share the same `Target` file path, trigger a validation error.
- Error must include: `Parallel conflict: T1.1 and T1.2 both target "path/to/file"`.

**Acceptance Check:**
Unit test in `src/internal/spec/tasks_validation_test.go` with two tasks targeting `src/common.go` marked as parallel, asserting validation failure.

### Phase 3: Implementation Engine (Readiness)

- [x] T3.1: [LOGIC] Update Readiness Logic
**Target:** `src/internal/spec/implementation.go`
**Context:** [US-10] (Implicit Readiness Rule)

**Action Steps:**
- Update `calculateReportStatus` or relevant readiness functions.
- A task is "Ready" if: It's in the active phase AND all non-parallel preceding tasks are `FINISHED`.

**Acceptance Check:**
Unit test asserting that T1.2 and T1.3 (parallel) are both `Ready` as soon as T1.1 (sequential predecessor) is `FINISHED`.

### Phase 4: Agent Artifacts

- [x] T4.1: [ARTIFACTS] Update `tasks.yaml` instructions
**Target:** `src/internal/agent/artifacts/spec/tasks.yaml`
**Context:** [US-1]

**Action Steps:**
- Change Rule 4 ("Strict Sequentiality") to allow parallelism when appropriate.
- Add instruction on how to use `**Parallel With:**` syntax.
- Emphasize the "File Isolation" rule for agents.

**Acceptance Check:**
`cat src/internal/agent/artifacts/spec/tasks.yaml` confirms updated instructions.

- [x] T4.2: [ARTIFACTS] Update `tasks.yaml` template
**Target:** `src/internal/agent/artifacts/spec/tasks.yaml`
**Context:** [US-1]

**Action Steps:**
- Add optional `**Parallel With:** T1.x` field to the task boilerplate in the template.
- Ensure the template example includes a sample parallel task for guidance.

**Acceptance Check:**
Visual verification of the template section in `tasks.yaml`.

### Phase 5: CLI & Reporting

- [x] T5.1: [CLI] Verify JSON Output
**Target:** `src/internal/spec/service.go`
**Context:** [US-8]

**Action Steps:**
- Ensure `SpecService.Status` correctly returns the updated `ImplementationReport` with parallel fields.
- Verify that `specforce spec status --json` output contains `parallel_with` and `is_parallel`.

**Acceptance Check:**
Run `specforce spec status --json` on a test specification and pipe to `jq` to verify fields.

## 3. Pre-emptive Mitigations
- **Risk:** Recursive parallel dependencies (T1.1 parallel with T1.2, T1.2 parallel with T1.3, etc.). -> **Mitigation:** The adjacency map approach will treat the entire cluster as parallel.
- **Risk:** Over-complicated regex for IDs. -> **Mitigation:** Use strict `T\d+.\d+` pattern.
