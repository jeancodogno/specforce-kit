---
slug: natural-task-format
lens: Backend-heavy
---

# Implementation Roadmap: Natural LLM Task Format

## 1. Execution Strategy
- **Gravity Order:** Update tests (Implementation & Tasks) -> Refactor `implementation.go` (Scanner support) -> Refactor `tasks.go` (Update support) -> Update templates -> Clean up skills.

## 2. Tasks

### Phase 1: Parser Refactoring (Hybrid Support)

#### T1.1: [TEST] Add hybrid test cases to implementation_test.go
**State:** `[FINISHED]`
**Target:** `src/internal/spec/implementation_test.go`
**Context:** [REQ-4]

**Action Steps:**
- Add `TestParseTasks_HybridFormat` to `implementation_test.go`.
- Include a sample `tasks.md` with mixed `####` and `- [ ]` tasks.
- Verify that `ParseTasks` returns the correct count and states for both formats.

**Verification (TDD):**
`go test -v -run TestParseTasks ./src/internal/spec/...`

#### T1.2: [CODE] Refactor ParseTasks regex in implementation.go
**State:** `[FINISHED]`
**Target:** `src/internal/spec/implementation.go`
**Context:** [REQ-1], [REQ-4]

**Action Steps:**
- Update `taskHeaderRegex` in `extractTasksFromContent` to support both `####` and `- [ ]` prefixes.
- Update `parseTaskBlock` to correctly map the checkbox character (`x`, `/`, ` `) to the `State` field if the `**State:**` tag is missing.

**Verification (TDD):**
`go test -v -run TestParseTasks ./src/internal/spec/...`

#### T1.3: [TEST] Add checklist test cases to tasks_test.go
**State:** `[FINISHED]`
**Target:** `src/internal/spec/tasks_test.go`
**Context:** [REQ-2]

**Action Steps:**
- Add `TestUpdateTaskStatusFile_WithChecklists` in `tasks_test.go`.
- Assert that updating a task with a checklist correctly mutates the `[ ]` to `[x]`.

**Verification (TDD):**
`go test -v -run TestUpdateTaskStatusFile ./src/internal/spec/...`

#### T1.4: [CODE] Implement hybrid logic in tasks.go
**State:** `[FINISHED]`
**Target:** `src/internal/spec/tasks.go`
**Context:** [REQ-1], [REQ-2]

**Action Steps:**
- Update `findTaskBlock` and `updateTaskStatusFile` to support the natural format.

**Verification (TDD):**
`go test -v -run TestUpdateTaskStatusFile ./src/internal/spec/...`

### Phase 2: Template & Skill Alignment

#### T2.1: [SCAFFOLD] Update tasks.yaml template and instructions
**State:** `[FINISHED]`
**Target:** `src/internal/agent/artifacts/spec/tasks.yaml`
**Context:** [REQ-3]

**Action Steps:**
- Set `- [ ]` as the default in the template.
- Simplify instructions to remove strict `####` constraints.

**Verification (TDD):**
`specforce spec artifact tasks` (Visual check).

#### T2.2: [SCAFFOLD] Update task-atomic-decomposition skill
**State:** `[FINISHED]`
**Target:** `.gemini/skills/task-atomic-decomposition/SKILL.md`
**Context:** [REQ-3]

**Action Steps:**
- Clean up any Specforce-specific formatting rules.

**Verification (TDD):**
`cat .gemini/skills/task-atomic-decomposition/SKILL.md`
