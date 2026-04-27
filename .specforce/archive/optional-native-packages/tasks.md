---
slug: optional-native-packages
lens: Migration
---

# Implementation Roadmap: Optional Native Packages

## 1. Execution Strategy
- **Gravity Order:** Build System (GoReleaser/Makefile) -> Node.js Proxy -> Main Package JSON -> CI/CD Workflow.

## 2. Tasks

### Phase 1: Build & Packaging

#### T1.1: [SCAFFOLD] Update GoReleaser Configuration
**State:** `[FINISHED]`
**Target:** `.goreleaser.yaml`
**Context:** [REQ-1]

**Action Steps:**
- Update `.goreleaser.yaml` archives section to output binaries in a directory structure compatible with standard NPM packaging.
- Remove legacy `checksums.txt` generation if no longer needed.

**Verification (TDD):**
Run `goreleaser build --snapshot --clean` and ensure binaries are generated correctly.

#### T1.2: [SCAFFOLD] Create Sub-package Generation Target
**State:** `[FINISHED]`
**Target:** `Makefile`
**Context:** [REQ-1]

**Action Steps:**
- Write a Make target that generates `package.json` for each supported OS/Arch combination (linux-x64, darwin-arm64, etc.) pointing to the compiled binaries in the `dist/` folder.

**Verification (TDD):**
Run the Make target locally and verify that a valid `package.json` with correct `os` and `cpu` properties is generated for each target platform.

### Phase 2: Proxy and Main Package

#### T2.1: [CODE] Implement Node.js Proxy
**State:** `[FINISHED]`
**Target:** `index.js`
**Context:** [REQ-2]

**Action Steps:**
- Create `index.js` at the root.
- Map `process.platform` and `process.arch` to the correct native package name.
- Use `require.resolve` to find the binary path in `node_modules`.
- Use `child_process.spawnSync` to execute it with `stdio: "inherit"`.
- Handle `MODULE_NOT_FOUND` exceptions with a clear diagnostic message and `process.exit(1)`.

**Verification (TDD):**
Create a mock executable in a local `node_modules` sub-folder, run `node index.js`, and verify it forwards execution and standard streams properly.

#### T2.2: [CODE] Update Main Package Definition
**State:** `[FINISHED]`
**Target:** `package.json`
**Context:** [REQ-2]

**Action Steps:**
- Remove the `go-npm` dependency.
- Remove the `postinstall` script.
- Add the generated sub-packages to the `optionalDependencies` block.
- Point the `"bin"` field to `"index.js"`.

**Verification (TDD):**
Run `npm pack` and verify the `package.json` inside the tarball no longer references `go-npm` or `postinstall`.

### Phase 3: CI Pipeline

#### T3.1: [SCAFFOLD] Update GitHub Actions Workflow
**State:** `[FINISHED]`
**Target:** `.github/workflows/release.yml`
**Context:** [REQ-1]

**Action Steps:**
- Modify `.github/workflows/release.yml` to publish the generated native sub-packages to NPM before publishing the main `@jeancodogno/specforce-kit` package.

**Verification (TDD):**
Use an action linter or syntax checker to ensure the workflow syntax is correct.
