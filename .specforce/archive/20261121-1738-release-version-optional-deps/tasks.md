---
slug: 20261121-1738-release-version-optional-deps
lens: Integration
---

# Implementation Roadmap: Native Package Versions in the Release Gate

## 1. Execution Strategy
- **Scope (small spec, no requirements file):** the goals below are the requirement references used by every task.
  - **[US-1] Gate:** `node scripts/release-tools.js check-version <tag>` also fails, before anything is built or published, when `package.json` does not list exactly the five native packages `@jeancodogno/specforce-kit-{linux-x64,linux-arm64,darwin-x64,darwin-arm64,win32-x64}` in `optionalDependencies`, each at the release version. The single-line error names every package that is missing, unexpected (another `@jeancodogno/specforce-kit-*` entry) or at the wrong version. Unrelated dependencies are ignored.
  - **[US-2] Sync:** `node scripts/sync-version.js` sets the five native package versions in `package.json` to the `package.json` version (adding missing ones, removing other `@jeancodogno/specforce-kit-*` entries), together with the existing `core.Version` update, keeping the file's 2-space JSON format and trailing newline. Running it twice changes nothing the second time.
- **Decisions (approved):** the five platform names are the same list as `PLATFORMS` in the `Makefile`; the pure functions live in `scripts/release-tools.js` (exported) and `scripts/sync-version.js` reuses them, so both are covered by `node:test` per the engineering rule "Release & CI Tooling".
- **Module impact:** none (release tooling only).
- **Gravity Order:** tests first (RED), then the shared functions and the gate (GREEN), then the sync script.

## 2. Tasks

### Phase 1: Gate and sync for native package versions

- [x] T1.1: [RED] Tests for native package version checks and sync
**Target:** `scripts/release-tools.test.js`
**Context:** [US-1], [US-2]

**Constraints:**
- Cases for the check: all five at the version passes; one at an old version fails naming it and both versions; one missing fails naming it; an extra `@jeancodogno/specforce-kit-foo` fails naming it; an unrelated optional dependency is ignored; no `optionalDependencies` at all fails naming all five; the check runs inside the existing version agreement so `check-version` reports it as one line.
- Cases for the sync: sets all five to the version; adds a missing one; removes an extra `@jeancodogno/specforce-kit-*`; keeps unrelated entries; idempotent (second run returns identical text); output is 2-space JSON with a trailing newline.
- Use only in-memory objects/strings; do not read or write the real `package.json`.

**Acceptance Check:**
`node --test scripts/*.test.js` exits non-zero because the new functions do not exist yet (RED state).

- [x] T1.2: [GREEN] Native package checks in the release tooling
**Target:** `scripts/release-tools.js`
**Context:** [US-1], [US-2]

**Constraints:**
- `check-version` reads `optionalDependencies` from the same `package.json` it already reads; keep the existing single-line `Error: ...` format and exit code 1.
- Existing tests and subcommands must keep their behaviour.

**Acceptance Check:**
`node --test scripts/*.test.js` passes all tests, and `node scripts/release-tools.js check-version v1.0.0-beta.3` prints `version ok: v1.0.0-beta.3`.

- [x] T1.3: [CODE] Version sync updates the native package versions
**Target:** `scripts/sync-version.js`
**Context:** [US-2]

**Constraints:**
- Reuse the exported sync function from `scripts/release-tools.js`; keep the existing `constants.go` update and its messages; write `package.json` only when its content changes.

**Acceptance Check:**
`cp package.json /tmp/pkg.bak && node scripts/sync-version.js && cmp -s package.json /tmp/pkg.bak && node scripts/release-tools.js check-version v1.0.0-beta.3 && echo OK` prints `OK` (the current tree is already in sync, so nothing changes).
