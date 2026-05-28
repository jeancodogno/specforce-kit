---
slug: 20260519-1803-spec-status-validation-guide
lens: Backend-heavy
---

# Implementation Roadmap: Validation Guide Golden Model

## 1. Execution Strategy
- **Gravity Order:** Data Model (Struct Update) -> Logic (Injection) -> Validation

## 2. Tasks

### Phase 1: Struct and Logic Implementation

- [x] T1.1: [CODE] Update ArtifactStatus Struct
**Target:** `src/internal/spec/status.go`
**Context:** [US-1]

**Action Steps:**
- Add `ValidationGuide string \`json:"validation_guide,omitempty"\`` to the `ArtifactStatus` struct.

**Acceptance Check:**
- Verify the struct compiles successfully.

- [x] T1.2: [CODE] Inject Golden Model
**Target:** `src/internal/spec/status.go`
**Context:** [US-1]

**Action Steps:**
- In `processArtifactStatus`, after calling `ValidateTasks` for `tasks.md`, check if `len(validationErrors) > 0`.
- If true, assign a static, fully valid Markdown string to `ValidationGuide` containing a valid Phase and Task structure.

**Acceptance Check:**
- Run `go test ./src/internal/spec/...` and verify no existing tests break. Create a dummy spec with errors and verify `validation_guide` is present in the `spec status --json` output.