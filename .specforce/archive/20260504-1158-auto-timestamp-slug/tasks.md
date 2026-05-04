---
slug: 20260504-1158-auto-timestamp-slug
lens: Backend-heavy
---

# Implementation Roadmap: Auto-Timestamped Spec Slugs

## 1. Execution Strategy
- **Gravity Order:** Domain Logic (Slug Service) -> Unit Tests -> CLI Integration -> Integration Testing.

## 2. Tasks

### Phase 1: Core Logic Implementation

#### T1.1: [SCAFFOLD] Create Slug Preparation Service
**State:** [FINISHED]
**Target:** src/internal/spec/slug.go
**Context:** [REQ-1], [REQ-2]
**Action Steps:**
- Create `src/internal/spec/slug.go` file.
- Define `PrepareSlug(rawSlug string) string` function.
- Implement logic to split path into segments using `filepath.ToSlash` and `strings.Split`.
- Implement regex check `^\d{8}-\d{4}-` on the final segment to detect existing timestamps.
- Generate local timestamp using `time.Now().Format("20060102-1504")`.
- Prepend timestamp to the final segment if missing.
- Sanitize the resulting segment to prevent double hyphens (e.g., `-feature` -> `{TS}-feature`).
- Join segments back using `filepath.FromSlash`.
**Verification (TDD):** `go test -v ./src/internal/spec/slug_test.go` (Note: Test file created in T1.2)

#### T1.2: [TEST] Implement Unit Tests for PrepareSlug
**State:** [FINISHED]
**Target:** src/internal/spec/slug_test.go
**Context:** [REQ-1], [REQ-2]
**Action Steps:**
- Create `src/internal/spec/slug_test.go`.
- Implement test cases covering:
    - Standard slug: `my-feature` -> `YYYYMMDD-HHMM-my-feature`.
    - Nested slug: `team-a/feature` -> `team-a/YYYYMMDD-HHMM-feature`.
    - Already timestamped: `20260101-1200-feature` -> `20260101-1200-feature`.
    - Hyphen prefix: `-feature` -> `YYYYMMDD-HHMM-feature`.
    - Multiple slashes: `team-a//feature` -> `team-a/YYYYMMDD-HHMM-feature`.
**Verification (TDD):** `go test -v ./src/internal/spec/slug_test.go`

### Phase 2: CLI Integration

#### T2.1: [CODE] Integrate Slug Preparation in CLI
**State:** [FINISHED]
**Target:** src/internal/cli/spec.go
**Context:** [REQ-1], [REQ-2]
**Action Steps:**
- In `HandleSpecInit`, call `spec.PrepareSlug(slug)` to transform the input.
- Assign the result back to the `slug` variable.
- Ensure the updated `slug` is used for `spec.SpecExists`, `filepath.Join`, and logging.
**Verification (TDD):** `go test -v ./src/internal/cli/cli_test.go`

#### T2.2: [UI] Verify CLI Output and Directory Creation
**State:** [FINISHED]
**Target:** CLI Command Execution
**Context:** [REQ-1]
**Action Steps:**
- Run `go run src/cmd/specforce/main.go spec init manual-verification-test`.
- Confirm the output displays the timestamped path.
- Confirm the directory `.specforce/specs/<TS>-manual-verification-test` exists.
**Verification (TDD):** `ls -d .specforce/specs/*manual-verification-test`

### Phase 3: Integration Testing

#### T3.1: [INT] Integration Test for Full Lifecycle
**State:** [FINISHED]
**Target:** src/internal/agent/integration_test.go
**Context:** [REQ-1], [REQ-2]
**Action Steps:**
- Add a new integration test case in `src/internal/agent/integration_test.go` that specifically checks the timestamping behavior during `spec init`.
- Verify both root-level and nested-level spec initialization.
**Verification (TDD):** `go test -v ./src/internal/agent/integration_test.go`
