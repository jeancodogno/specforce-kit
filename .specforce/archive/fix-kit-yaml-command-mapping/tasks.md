---
slug: fix-kit-yaml-command-mapping
lens: Backend-heavy
---

# Implementation Roadmap: Unique Command Names in Skill Headers

## 1. Execution Strategy
- **Gravity Order:** TDD Setup -> Core Logic Update -> Metadata Enhancement -> Verification.
- We first reproduce the issue by adding a failing test case to the translator service. Then we apply the logic fix and finally enhance the source metadata to ensure a robust solution.

## 2. Tasks

### Phase 1: Reproduction & Regression Testing

#### T1.1: [TEST] Add Failing Case for Skill Name Leakage
**State:** `[FINISHED]`
**Target:** `src/internal/agent/translator_test.go`
**Context:** [REQ-1]

**Action Steps:**
- Add a new test case to `TestInjectYAMLHeader` (or create a new test function) that simulates a mapping with `Name: "SKILL"` and an empty metadata name.
- Assert that the resulting content contains `name: SKILL` (proving the bug).

**Verification (TDD):**
`go test -v src/internal/agent/translator_test.go` - The test MUST fail with the current code.

### Phase 2: Core Logic Implementation

#### T2.1: [CODE] Fix Generic Name Fallback in Translator
**State:** `[FINISHED]`
**Target:** `src/internal/agent/translator.go`
**Context:** [REQ-1]

**Action Steps:**
- Modify `injectYAMLHeader` to check if `displayName` is `"SKILL"`.
- If it is, fallback to `bp.ID`.
- For `category == "commands"`, prefix the ID with `"spf."`.

**Verification (TDD):**
`go test -v src/internal/agent/translator_test.go` - The test added in T1.1 MUST now pass.

#### T2.2: [CODE] Implement Header Name Uniqueness Validator
**State:** `[FINISHED]`
**Target:** `src/internal/agent/translator.go`
**Context:** [REQ-3]

**Action Steps:**
- Add a mechanism (e.g., a shared map in the translator or registry context) to track all `displayName` values generated in a session.
- If a name is reused for a different `bp.ID`, return an error.

**Verification (TDD):**
Add a test case in `translator_test.go` that attempts to translate two different blueprints that result in the same name, and verify it returns an error.

### Phase 3: Metadata Enhancement

#### T3.1: [CONFIG] Update Command Metadata
**State:** `[FINISHED]`
**Target:** `Global Scope`
**Context:** [REQ-1, REQ-2]

**Action Steps:**
- Update `src/internal/agent/kit/commands/archive.yaml`: Add `name: spf.archive` at root.
- Update `src/internal/agent/kit/commands/constitution.yaml`: Add `name: spf.constitution` at root.
- Update `src/internal/agent/kit/commands/implement.yaml`: Add `name: spf.implement` at root.
- Update `src/internal/agent/kit/commands/spec.yaml`: Add `name: spf.spec` at root.

**Verification (TDD):**
Inspect the files to ensure valid YAML structure and the presence of the `name` field.

### Phase 4: Final Verification

#### T4.1: [VERIFY] Full Integration Check
**State:** `[FINISHED]`
**Target:** `Global Scope`
**Context:** [REQ-1]

**Action Steps:**
- Run all project tests to ensure no regressions.
- (Optional) Run `go run ./src/cmd/specforce/main.go spec refresh` if applicable to see generated outputs in a test environment.

**Verification (TDD):**
`go test ./src/internal/agent/...`
