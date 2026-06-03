# Tasks: Centralized Version Management

## 1. Execution Strategy
- **Gravity Order:** Infrastructure (Sync Script & Hooks) -> Source Unification (Go Refactoring) -> Build Injection (Makefile). We start by creating the automation script and NPM hooks to establish the SSoT mechanism, then refactor the Go source code to use a single variable, and finally update the Makefile to ensure build-time consistency.

## 2. Tasks

### Phase 1: Infrastructure & Automation

- [x] T1.1: [SCAFFOLD] Create Version Sync Script
**Target:** `scripts/sync-version.js`
**Context:** [US-2]
**Parallel With:** T1.2

**Action Steps:**
- Create the `scripts/` directory in the project root.
- Implement `sync-version.js` to read the `version` from `package.json`.
- Use `fs.readFileSync` and `fs.writeFileSync` with a regex to surgically update `var Version = "..."` in `src/internal/core/constants.go`.

**Acceptance Check:**
Run `node scripts/sync-version.js` and verify that the content of `src/internal/core/constants.go` is updated to match the version in `package.json`.

- [x] T1.2: [CONFIG] Configure package.json Version Hook
**Target:** `package.json`
**Context:** [US-2]
**Parallel With:** T1.1

**Action Steps:**
- Locate the `"scripts"` section in `package.json`.
- Add a `"version"` lifecycle hook: `"node scripts/sync-version.js && git add src/internal/core/constants.go"`.
- Verify the script is correctly registered.

**Acceptance Check:**
Run `npm version 9.9.9 --no-git-tag-version` (temporary test) and verify that `constants.go` is updated and staged. Revert changes after test.

### Phase 2: Go Source Unification

- [x] T2.1: [CODE] Refactor Version Constant to Variable
**Target:** `src/internal/core/constants.go`
**Context:** [US-1]

**Action Steps:**
- Change `const Version = "1.0.0-alpha.2"` to `var Version = "1.0.0-alpha.2"`.
- Add a comment: `// Version is the current tool version, synced from package.json and optionally overwritten at build time.`
- Ensure the package name remains `core`.

**Acceptance Check:**
`go build ./src/internal/core/...` passes without errors.

- [x] T2.2: [CODE] Unify CLI Entrypoint Version
**Target:** `src/cmd/specforce/main.go`
**Context:** [US-1]

**Action Steps:**
- Import `github.com/jeancodogno/specforce-kit/src/internal/core`.
- Remove the local `var version = "..."` declaration.
- Update any reference to the version flag/display to use `core.Version`.

**Acceptance Check:**
`go build ./src/cmd/specforce/main.go` works and `grep "version =" src/cmd/specforce/main.go` returns no results.

- [x] T2.3: [CODE] Unify Agent Registry Version
**Target:** `src/internal/agent/registry.go`
**Context:** [US-1]

**Action Steps:**
- Ensure `github.com/jeancodogno/specforce-kit/src/internal/core` is imported.
- Locate the `AgentMetadata` initialization in `loadAgents`.
- Replace the hardcoded `"1.0.0-alpha.2"` with `core.Version`.

**Acceptance Check:**
`go test ./src/internal/agent/...` passes.

- [x] T2.4: [CODE] Unify TUI Logo Version
**Target:** `src/internal/tui/logo.go`
**Context:** [US-1]

**Action Steps:**
- Import `github.com/jeancodogno/specforce-kit/src/internal/core`.
- Remove the local `var AppVersion = "..."` declaration.
- Update `GenerateLogo` to use `core.Version` instead of `AppVersion`.

**Acceptance Check:**
Run a local build and verify the logo displays the correct version.

### Phase 3: Build & Verification

- [x] T3.1: [CONFIG] Update Makefile for Version Injection
**Target:** `Makefile`
**Context:** [US-3]

**Action Steps:**
- Ensure `VERSION=$(shell node -p "require('./package.json').version")` is defined at the top.
- Update the `build` target to include `-ldflags="-s -w -X 'github.com/jeancodogno/specforce-kit/src/internal/core.Version=$(VERSION)'"`.
- Verify the linker path matches the project structure.

**Acceptance Check:**
Run `make build` and then `./specforce --version` to confirm it matches `package.json`.

- [x] T3.2: [VERIFY] End-to-End Release Simulation
**Target:** `Global Scope`
**Context:** [US-3]

**Action Steps:**
- Update `package.json` version manually or via `npm version`.
- Verify `constants.go` updated automatically (if using npm version).
- Run `make build`.
- Execute `./specforce --version`.

**Acceptance Check:**
The final binary reports the exact version defined in `package.json`.
