# Technical Design: Standardize Acceptance Check

## 1. Architecture: Parsing & Validation State Machine

The core change involves updating the `tasks.md` parser to transition from the methodology-specific `**Acceptance Check:**` label to the methodology-agnostic `**Acceptance Check:**` label.

### 1.1 Parsing Logic (`src/internal/spec/tasks.go`)
The `updateTaskBlockState` function will be refactored to:
1.  **Enforce New Label:** Recognize `**Acceptance Check:**` as the primary verification header.
2.  **Case-Insensitivity:** Perform prefix checks using `strings.ToLower()` to satisfy the flexibility requirement while maintaining the standard casing in templates.
3.  **Deprecation Tracking:** Explicitly detect the legacy `**Acceptance Check:**` label to provide contextual error messages instead of generic "missing field" errors.

### 1.2 Data Extraction (`src/internal/spec/implementation.go`)
The `parseTaskBlock` function used for `implementation status --json` will be updated to extract the verification command from the new label section using an updated regex.

### 1.3 Validation Feedback
The `ValidateTasks` service will return a specific error when the legacy label is found, suggesting the correct migration path.

## 2. File Inventory

### 2.1 Core Logic & Services
- `src/internal/spec/tasks.go`: Update `updateTaskBlockState` and `validateLastTask`.
- `src/internal/spec/implementation.go`: Update `parseTaskBlock` regex for verification extraction.
- `src/internal/spec/status.go`: Update `validationGuide` example string.

### 2.2 Templates & Artifacts
- `src/internal/agent/artifacts/spec/tasks.yaml`: Update instruction #3 and the Markdown template.
- `src/internal/agent/artifacts/spec/bug-design.yaml`: Update section header.

### 2.3 Agent & Skill Instructions
- `src/internal/agent/kit/agents/specforce-planner.yaml`: Update instruction #4.
- `src/internal/agent/kit/commands/implement.yaml`: Update Step 3.B verification description.
- `src/internal/agent/kit/skills/task-atomic-decomposition/SKILL.yaml`: Update mandatory verification heading and checklist.
- `src/internal/project/agents_md.go`: Update hardcoded instruction string.

### 2.4 Test Suite
- `src/internal/spec/tasks_validation_test.go`
- `src/internal/spec/implementation_test.go`
- `src/internal/cli/implementation_test.go`

## 3. Migration Strategy (Hard Migration)

To maintain system integrity, all existing specifications (active and archived) must be updated. This will be performed via a global batch replacement.

### 3.1 Execution Script
```bash
# Target: All Markdown files in .specforce directory
find .specforce -type f -name "*.md" -exec sed -i 's/\*\*Verification (TDD):\*\*/\*\*Acceptance Check:\*\*/g' {} +
```

### 3.2 Verification
After migration, `specforce spec status --all` must return `IsValid: true` for all migrated specs.

## 4. API & Data Contracts

### 4.1 Internal Task Structure
No changes to the Go `ImplementationTask` struct. The `Verification` field will now simply map to the content under the `**Acceptance Check:**` header.

### 4.2 Regex Contracts
The new regex for extraction in `implementation.go`:
`(?s)\*\*Acceptance Check:\*\*\n(.*?)(?:\n\n|\n###|\n##|$)`

## 5. Security & Integrity

- **Idempotency:** The migration script is safe to run multiple times as it uses exact string matching.
- **Input Validation:** The parser will strictly reject any task that does not use the new label, preventing "methodology drift" in new specs.
- **Data Safety:** `sed -i` performs in-place replacement. Standard Git practices (commit before migration) are assumed for rollback safety.

## 6. UI/UX Refinement (Ghost Protocol)

- **Error State:** When legacy labels are detected, the CLI will output:
  > `[ERROR] Task T1.1 uses deprecated **Acceptance Check:** label. Use **Acceptance Check:** instead.`
- **Success State:** The `status` command will display the new label in the verification summary.
