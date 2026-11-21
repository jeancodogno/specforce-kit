---
slug: 20260520-1434-fix-config-init-logic
lens: Bugfix
---

# Bugfix: fix config initialization logic

## 1. Issue Description
When executing `specforce init`, the configuration file `.specforce/config.yaml` is only created if the `.specforce/` directory does not already exist. If the directory exists (e.g., from a partial initialization or manual creation), the CLI enters the "update" flow which skips the `config.yaml` creation, leaving the project in an inconsistent state.

## 2. Evidence & Observations
- **Symptom:** After running `specforce init` in a directory where `.specforce/` exists but `config.yaml` is missing, the file is still missing.
- **Trace:** `src/internal/cli/cli.go:114` calls `core.EnsureConfigExists(".")` inside `handleNewInitFlow`, but `HandleInit` (line 55) redirects to `handleUpdateFlow` if the directory exists, which lacks this call.

## 3. Reproduction Steps
1. Create a dummy project directory.
2. Manually create a `.specforce/` directory: `mkdir .specforce`.
3. Run `specforce init`.
4. Observe that `.specforce/config.yaml` was NOT created.
5. **Expected Outcome:** `.specforce/config.yaml` should be created even if the directory already exists.

## 4. Root Cause Analysis (RCA)
The CLI implementation in `src/internal/cli/cli.go` uses `project.IsInitialized(".")` to decide between `handleNewInitFlow` and `handleUpdateFlow`. `handleNewInitFlow` contains the logic to ensure `config.yaml` exists, while `handleUpdateFlow` (and the `ProjectService.InitializeProject` or `UpdateTools` methods) does not guarantee its presence. Since `core.EnsureConfigExists` is idempotent, it should be called regardless of the flow.

## 5. Acceptance Criteria (Regression Tests)

### [FIX-1] Idempotent Config Creation
**Scenario: [Regression]**
GIVEN a directory with an existing `.specforce/` folder but no `config.yaml`
WHEN `specforce init` is executed
THEN the `.specforce/config.yaml` file is successfully created with default content.

### [FIX-2] Update Flow Integrity
**Scenario: [Regression]**
GIVEN a directory that is already initialized (has `.specforce/`)
WHEN `specforce init` is run to update agents
THEN it should still ensure that `config.yaml` exists.

## 6. Technical Constraints (NFR)
- **[Safety]:** Use `core.EnsureConfigExists` which uses `SecurePath` to prevent directory traversal.
- **[Observability]:** Ensure any error during config creation is reported as a warning (consistent with current behavior).
