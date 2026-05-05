---
slug: 20260505-0926-tasks-md-validation-in-status
lens: Backend-heavy
---

# Implementation Roadmap: Tasks Markdown Validation in Status

## 1. Execution Strategy
The implementation follows a logical flow from internal domain logic to the user interface:
1. **Domain Model**: Extend the status structs to support validation error reporting.
2. **Validation Logic**: Implement the core `ValidateTasks` logic in the spec domain.
3. **Service Integration**: Wire the validator into the existing `GetStatus` service.
4. **CLI/TUI Enhancement**: Update the rendering logic to display errors to the user.

## 2. Tasks

### Phase 1: Domain & Validation Logic

- [x] T1.1: [SCAFFOLD] Update status data models
**Target:** `src/internal/spec/status.go`
**Context:** [REQ-3]

**Action Steps:**
- Add `ValidationErrors []string` to `ArtifactStatus` struct.
- Add `IsValid bool` to `SpecStatus` struct.

**Verification (TDD):**
- Verify that `go build ./src/internal/spec/...` succeeds and the new fields are available in the JSON tags.

- [x] T1.2: [CODE] Implement `ValidateTasks` in `tasks.go`
**Target:** `src/internal/spec/tasks.go`
**Context:** [REQ-1, REQ-2]

**Action Steps:**
- Implement `ValidateTasks(ctx context.Context, projectRoot, slug string) ([]string, error)`.
- Use line-by-line parsing to track Phases (H3) and Tasks (H4).
- Validate that Tasks appear after a Phase.
- Validate Task IDs match the parent Phase ID and are sequential.
- For each Task block, validate the presence and non-emptiness of:
    - `**Target:**`
    - `**Context:**`
    - `**Action Steps:**` (must have at least one `-` item)
    - `**Verification (TDD):**`
- Validate no Phase is empty.

**Verification (TDD):**
- Create a new test file `src/internal/spec/tasks_validation_test.go` with table-driven tests covering happy paths, malformed hierarchies, naming violations, and missing mandatory task block sections.

### Phase 2: Service Integration

- [x] T2.1: [CODE] Integrate validator into `GetStatus`
**Target:** `src/internal/spec/status.go`
**Context:** [REQ-3]

**Action Steps:**
- In `GetStatus`, if the `tasks.md` file exists, call `ValidateTasks`.
- Populate `ArtifactStatus.ValidationErrors` for the "tasks" artifact.
- Set `SpecStatus.IsValid` to `false` if any validation errors are found.

**Verification (TDD):**
- Update `src/internal/spec/status_test.go` to include a test case where a malformed `tasks.md` is present and verify `IsValid` is false.

### Phase 3: CLI & TUI Enhancement

- [x] T3.1: [CODE] Update TUI rendering for errors
**Target:** `src/internal/tui/spec.go`
**Context:** [REQ-3]

**Action Steps:**
- Update `RenderSpecStatus` to check for `ValidationErrors` in each artifact.
- If errors exist, render them below the artifact description in a red-indented list.

**Verification (TDD):**
- Run `go run ./src/cmd/specforce/main.go spec status <slug>` on a spec with a malformed `tasks.md` and visually confirm the error reporting.

## 3. Pre-emptive Mitigations
- **Risk:** Parsing logic might be slow for very large `tasks.md` files -> **Mitigation:** Use a simple scanner and early context cancellation checks.
- **Risk:** Existing `tasks.md` files might already be slightly non-compliant (e.g., missing specific ID format) -> **Mitigation:** Ensure the validator is strict but provides very clear error messages on how to fix the format.
