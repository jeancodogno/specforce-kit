---
slug: 20260520-1619-fix-spec-list-json-null
lens: Bugfix
---

# Bugfix: Fix Spec List JSON Null Output

## 1. Issue Description
When executing `specforce spec list --json` in a workspace with no active specifications, the command returns the literal string `null` instead of an empty JSON array `[]`. This breaks machine-readable pipelines that expect a list structure.

## 2. Evidence & Observations
- **Symptom:** `specforce spec list --json` outputs `null`.
- **Trace:** The issue originates in `src/internal/spec/discovery.go` within the `ListActiveSpecs` function and is surfaced in `src/internal/cli/spec.go` during JSON marshaling.

## 3. Reproduction Steps
1. Ensure the directory `.specforce/specs` exists but contains no subdirectories (only files like `.gitkeep`).
2. Run the command: `go run src/cmd/specforce/main.go spec list --json`.
3. Observe the output is `null` instead of `[]`.

## 4. Root Cause Analysis (RCA)
- In Go, `var s []T` declares a `nil` slice.
- The `encoding/json` package marshals `nil` slices as `null`.
- `ListActiveSpecs` uses `var specs []SpecInfo`, and when no directories are found to append, it returns this `nil` slice.

## 5. Acceptance Criteria (Regression Tests)

### [FIX-1] Empty Spec List JSON
**Scenario: [Regression]**
GIVEN a project with no active specifications in `.specforce/specs`
WHEN running `specforce spec list --json`
THEN the output MUST be a valid empty JSON array `[]`.

### [FIX-2] Single Spec List JSON
**Scenario: [Sanity]**
GIVEN a project with one active specification (e.g., `test-spec`)
WHEN running `specforce spec list --json`
THEN the output MUST be a JSON array containing the spec info `[{"slug": "test-spec"}]`.

## 6. Technical Constraints (NFR)
- **[Safety]:** Ensure the slice is initialized as empty `[]SpecInfo{}` to guarantee non-nil JSON marshaling.
- **[Compatibility]:** Maintain existing `SpecInfo` struct field names.
