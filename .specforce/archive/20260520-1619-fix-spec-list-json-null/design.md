---
slug: 20260520-1619-fix-spec-list-json-null
lens: Bugfix
---

# Technical Design: Fix Spec List JSON Null Output (Fix Blueprint)

## 1. Code Path Inventory
- `src/internal/spec/discovery.go` -> Change `var specs []SpecInfo` to `specs := []SpecInfo{}` in `ListActiveSpecs` to ensure an empty slice is returned instead of `nil`. This prevents `encoding/json` from outputting `null`.

## 2. Regression Strategy (Verification Plan)
- **Unit Tests:** Add a test case in `src/internal/spec/discovery_test.go` (create if missing) to verify that `ListActiveSpecs` returns a non-nil empty slice (`[]SpecInfo{}`) when the `.specforce/specs` directory is empty or contains no valid specs.
- **Integration Tests:** Update `src/internal/cli/spec_test.go` to capture stdout during `HandleSpecList(ctx, ui, true)` and assert that the output is exactly `[]` (standard JSON for empty array) when no specs exist.
- **Manual Verification:** Execute `go run src/cmd/specforce/main.go spec list --json` in a project with no active specs and verify the output is `[]`.

## 3. Side Effects & Risks
- **Performance:** Negligible; explicit slice initialization is a standard Go practice with no measurable overhead.
- **Compatibility:** Positive impact; standardizes the JSON response for consumers (automated scripts, CI/CD tools) expecting a consistent list structure even when empty.

## 4. Proposed Fix (Abstract Logic)
```go
// src/internal/spec/discovery.go

func ListActiveSpecs(ctx context.Context, projectRoot string) ([]SpecInfo, error) {
    // ... (directory reading logic)

    // Initialize as empty slice, not nil, to ensure correct JSON marshaling to []
    specs := []SpecInfo{} 
    
    // ... (loop and append logic)
    
    return specs, nil
}
```
