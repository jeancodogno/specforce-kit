---
slug: fix-command-not-found
lens: Backend-heavy
---

# Implementation Roadmap: Fix Command Not Found

## 1. Execution Strategy
- **Gravity Order:** Documentation -> Configuration -> Implementation -> Verification.

## 2. Tasks

### Phase 1: Documentation

#### T1.1: [DOCS] Add Troubleshooting to README.md
**State:** `[FINISHED]`
**Target:** `README.md`
**Context:** [REQ-3]
**Verification:** Assert `README.md` contains a "Troubleshooting" section specifically mentioning the "Command Not Found" error and linking to detailed documentation.

#### T1.2: [DOCS] Add Detailed Troubleshooting to docs/getting-started.md
**State:** `[FINISHED]`
**Target:** `docs/getting-started.md`
**Context:** [REQ-3]
**Verification:** Assert `docs/getting-started.md` includes OS-specific steps for manual PATH verification (export for Unix, setx for Windows) and `make build` instructions.

### Phase 2: Configuration

#### T2.1: [CONFIG] Add build scripts to package.json
**State:** `[FINISHED]`
**Target:** `package.json`
**Context:** [REQ-2]
**Verification:** Assert `package.json` includes `"prepare": "make build || true"` and `"postinstall": "make build || true"`. Run `npm run postinstall` and verify it exits with 0 even if `make` is unavailable.

### Phase 3: Implementation

#### T3.1: [CODE] Implement Diagnostic Logic in index.js
**State:** `[FINISHED]`
**Target:** `index.js`
**Context:** [REQ-1]
**Verification:** Execute `node index.js` with a mocked failed binary resolution and verify it correctly identifies whether the NPM bin directory is missing from `process.env.PATH`.

#### T3.2: [CODE] Implement OS-specific instructions in index.js
**State:** `[FINISHED]`
**Target:** `index.js`
**Context:** [REQ-1]
**Verification:** Execute the diagnostic logic on a Unix-like system and verify it suggests `export PATH=...`. Execute on Windows and verify it suggests `setx PATH ...`.

### Phase 4: Verification

#### T4.1: [TEST] Verify Diagnostic Output
**State:** `[FINISHED]`
**Target:** `CLI`
**Context:** [REQ-1]
**Verification:** Remove the native binary, run `node index.js`, and verify the final diagnostic output matches the "Ghost in the Machine" aesthetic defined in the Design Spec.
