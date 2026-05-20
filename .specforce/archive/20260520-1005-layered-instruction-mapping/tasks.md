---
slug: 20260520-1005-layered-instruction-mapping
lens: Integration
---

# Implementation Roadmap: Layered Instruction Mapping (TDD)

## 1. Execution Strategy
The implementation follows a strict TDD (Red-Green-Refactor) cycle. We first update the configuration template, then establish a failing test harness that proves the layered mapping requirements, and finally implement the logic until all tests pass.

- **Gravity Order:** Config Update -> Failing Tests (Red) -> Implementation (Green) -> Refactor.

## 2. Tasks

### Phase 1: Foundation & Red Stage (Failing Tests)

- [x] T1.1: Update DefaultConfigContent in config.go
**Target:** `src/internal/core/config.go`
**Context:** [US-1]
**Action Steps:**
- Locate the `DefaultConfigContent` constant.
- Update the `instructions:` block comments to document that keys can be generic base types (requirements, design, tasks) or specific prefixed names (e.g., feature-requirements).
**Verification (TDD):**
- Run `go build ./...` to ensure no syntax errors.

- [x] T1.2: Create Failing Unit Tests for Layered Mapping
**Target:** `src/internal/spec/service_test.go`
**Context:** [US-1, US-2]
**Action Steps:**
- Add new test cases to `TestGetArtifact` (or create a specialized test) that assert:
    - `requirements` + `feature-requirements` are merged correctly.
    - `bugfix-requirements` resolves to `requirements` instructions if specific ones are missing.
    - `tasks-for-design` resolves to `design` instructions (Right-to-Left inference).
    - Duplicate instructions are deduplicated.
**Verification (TDD):**
- Run `go test -v ./src/internal/spec/...` and verify the new tests FAIL (Red Stage).

### Phase 2: Green Stage (Implementation)

- [x] T2.1: Implement Base Type Inference
**Target:** `src/internal/spec/service.go`
**Context:** [US-2]
**Action Steps:**
- Implement `inferBaseType(name string) string` using right-to-left keyword matching.
**Verification (TDD):**
- Run tests; ensure the inference-related test cases pass.

- [x] T2.2: Implement Layered Instruction Resolution
**Target:** `src/internal/spec/service.go`
**Context:** [US-1]
**Action Steps:**
- Implement `resolveInstructions(conf *core.ProjectConfig, name string) []string` to merge base and specific rules with deduplication.
**Verification (TDD):**
- Run tests; ensure the resolution logic passes.

- [x] T2.3: Refactor GetArtifact to use Layered Logic
**Target:** `src/internal/spec/service.go`
**Context:** [US-1]
**Action Steps:**
- Update `GetArtifact` to use `resolveInstructions`.
**Verification (TDD):**
- Run `go test -v ./src/internal/spec/...` and verify all tests pass (Green Stage).

### Phase 3: Refactor & Final Verification

- [x] T3.1: Code Cleanup and Refactoring
**Target:** `src/internal/spec/service.go`
**Context:** [US-1]
**Action Steps:**
- Improve naming and structure of the new private methods.
- Ensure no regression in existing spec status logic.
**Verification (TDD):**
- Run all tests and ensure they remain GREEN.

- [x] T3.2: Final E2E Verification
**Target:** `Global Scope`
**Context:** [US-1]
**Action Steps:**
- Run `specforce spec artifact <slug> --json` with a test `config.yaml` containing layered rules.
- Verify the output JSON contains both layers correctly.
**Verification (TDD):**
- Manually inspect JSON output.
