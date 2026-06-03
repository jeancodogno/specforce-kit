# Proposal: Centralize Version Management

## Context & Problem
Currently, the application version (`1.0.0-alpha.2`) is hardcoded in at least 6 different locations:
- `package.json`
- `AGENTS.md` (will be ignored per user request)
- `src/cmd/specforce/main.go`
- `src/internal/agent/registry.go`
- `src/internal/core/constants.go`
- `src/internal/tui/logo.go`

This duplication makes the release process error-prone and tedious.

## Proposed Solution
Establish `package.json` as the **Single Source of Truth** for the project version and automate its propagation to the Go codebase.

### Technical Strategy
1.  **Consolidate Go Source:**
    - Standardize all components to use `core.Version` from `src/internal/core/constants.go`.
    - Remove redundant `version` or `AppVersion` variables in other packages.
2.  **Automate Synchronization:**
    - Use NPM lifecycle hooks (`version`) to automatically update `src/internal/core/constants.go` whenever `npm version` is called.
    - Create a small utility script (`scripts/sync-version.js`) to perform the surgical update of the Go constant.
3.  **Build-Time Reinforcement:**
    - Update `Makefile` to inject the version from `package.json` into the Go binary using `-ldflags`.
    - This ensures that even without the sync script (e.g., during local development), the binary reflects the version in `package.json`.

## Expected Benefits
- **Zero Duplication:** Version is changed in one place only.
- **Reliability:** Automated sync prevents "forgotten" files.
- **Consistency:** Binaries and NPM packages will always report the same version.
