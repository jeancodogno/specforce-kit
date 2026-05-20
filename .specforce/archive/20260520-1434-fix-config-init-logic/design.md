---
slug: 20260520-1434-fix-config-init-logic
lens: Bugfix
---

# Technical Design: fix config initialization logic (Fix Blueprint)

## 1. Code Path Inventory
- `src/internal/cli/cli.go` -> Move `core.EnsureConfigExists(".")` to a shared location or ensure it is called in both `handleNewInitFlow` and `handleUpdateFlow`.
- `src/internal/project/service.go` -> (Optional but recommended) Add `core.EnsureConfigExists` to `InitializeProject` to ensure structural integrity at the service level.

## 2. Regression Strategy (Verification Plan)
- **Unit Tests:** Add a test case to `src/internal/cli/cli_test.go` named `TestHandleInit_EnsuresConfigExistsInUpdateFlow` that mocks an existing `.specforce` dir and checks for `config.yaml` after execution.
- **Manual Verification:** 
  ```bash
  mkdir .specforce
  go run src/cmd/specforce/main.go init
  ls .specforce/config.yaml
  ```

## 3. Side Effects & Risks
- **Performance:** `core.EnsureConfigExists` performs a few file system stats and a write if missing; overhead is negligible for an `init` command.
- **Compatibility:** No impact on existing projects as the function does not overwrite existing configurations.

## 4. Proposed Fix (Abstract Logic)
```go
// In src/internal/cli/cli.go
func (e *Executor) HandleInit(...) error {
    // ... agent resolution ...
    
    // Ensure config exists early or in both branches
    defer func() {
        _ = core.EnsureConfigExists(".")
    }()
    
    if project.IsInitialized(".") {
        return e.handleUpdateFlow(...)
    }
    return e.handleNewInitFlow(...)
}
```
*Note: Better to call it explicitly in a common point or inside both handlers to handle errors consistently.*
