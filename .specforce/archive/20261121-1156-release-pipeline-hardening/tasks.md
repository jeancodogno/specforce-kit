---
slug: 20261121-1156-release-pipeline-hardening
lens: Integration
---

# Implementation Roadmap: Release Pipeline Hardening for the Stable 1.0.0

## 1. Execution Strategy
- **Gravity Order:** Release decision logic as a tested Node script (RED then GREEN) -> release build configuration and release workflow that call it -> continuous integration that runs everything on every version line.
- **Tooling decision (approved):** The release decisions (version agreement, channel derivation, publish-or-skip, notes selection) live in `scripts/release-tools.js`, a CommonJS Node script next to `scripts/sync-version.js`, tested with the built-in `node:test` runner of Node 20 via `node --test scripts/`. No new npm or Go dependencies. The workflows only call the script's subcommands.
- **Test pairing:** Phase 1 follows RED/GREEN. Phases 2 and 3 change only CI/CD configuration files, which have no unit-test harness in this repository; each of those tasks is proven by a non-interactive check on the file it changes, and the end-to-end behaviour is exercised by the maintainer's `v1.0.0-rc.1` rehearsal (out of scope here).
- **Module impact:** none (no `.specforce/docs/modules/` behaviour changes; `delta.impact: none`).

## 2. Tasks

### Phase 1: Release decision tooling

- [x] T1.1: [RED] Tests for the release decision tooling
**Target:** `scripts/release-tools.test.js`
**Context:** [US-2], [US-3], [US-4], [US-5], [US-6]

**Constraints:**
- Test the exported pure functions and the publish routine with an injected command runner and an injected file reader; no test may touch the network, the real registry or the real repository files.
- Cases to cover: version agreement (`v1.0.0`/`1.0.0`/`1.0.0` passes; `v1.0.0`/`1.0.0-rc.1`/`1.0.0` fails with a message containing all three values; `vnext` fails as not a valid version); channel (`1.0.0` -> `latest`, `1.0.0-beta.3` -> `beta`, `1.1.0-next.1` -> `next`, `1.1.0-1` fails asking for a named channel); publish (absent -> publishes with `--access public --provenance --tag <channel>` and logs `published <name>@<version> (<channel>)`; present -> skips and logs `skipped <name>@<version> (already published)`; existence query error -> fails without publishing; publish error -> fails); notes (non-empty file -> its path; missing file -> nothing; whitespace-only file -> nothing).
- The channel rule MUST match `ChannelOf` in `src/internal/upgrade/provider_npm.go` (first pre-release identifier, `latest` for stable).

**Acceptance Check:**
`node --test scripts/` exits non-zero because `scripts/release-tools.js` does not exist yet (RED state).

- [x] T1.2: [GREEN] Release decision tooling with CLI subcommands
**Target:** `scripts/release-tools.js`
**Context:** [US-2], [US-3], [US-4], [US-5], [US-6]

**Constraints:**
- Subcommands used by the workflow: `check-version <tag>` (reads the version from `package.json` and the `Version` declaration in `src/internal/core/constants.go`, same regex style as `scripts/sync-version.js`), `channel <version>` (prints the channel), `publish <package-dir> <channel>` (existence query via `npm view <name>@<version> version`, where an "E404"/not-found result means absent and any other failure is an error; then `npm publish --access public --provenance --tag <channel>`), `notes <tag>` (prints `.github/release-notes/<tag>.md` when non-empty, prints nothing otherwise).
- Every failure exits with code 1 and a single-line reason on stderr; never use `shell: true` (same rule as `index.js`).

**Acceptance Check:**
`node --test scripts/` passes all tests, and `node scripts/release-tools.js channel 1.1.0-next.1` prints `next`.

### Phase 2: Release build and release workflow

- [x] T2.1: [CONFIG] Pre-release marking and built-in version in the release build
**Target:** `.goreleaser.yaml`
**Context:** [US-3]
**Parallel With:** T2.2

**Constraints:**
- Add `release.prerelease: auto`; replace the ineffective `-X main.version={{.Version}}` with `-X github.com/jeancodogno/specforce-kit/src/internal/core.Version={{.Version}}` (same package path as the `Makefile` build target).

**Acceptance Check:**
`grep -q 'prerelease: auto' .goreleaser.yaml && grep -q 'src/internal/core.Version={{.Version}}' .goreleaser.yaml && ! grep -q 'main.version' .goreleaser.yaml && echo OK` prints `OK`.

- [x] T2.2: [CONFIG] Gated, resumable, provenance-enabled release workflow
**Target:** `.github/workflows/release.yml`
**Context:** [US-2], [US-3], [US-4], [US-5], [US-6]
**Parallel With:** T2.1

**Constraints:**
- A `verify` job runs `node scripts/release-tools.js check-version ${{ github.ref_name }}`, golangci-lint (same action settings as `ci.yml`), `go test ./...` and `node --test scripts/`; the `release` job declares `needs: verify`.
- The `release` job adds `id-token: write` to the existing `contents: write` permissions.
- GoReleaser receives `--release-notes <path>` only when `release-tools.js notes <tag>` printed a path; the log states which source was used.
- Remove the hard-coded alpha/beta/rc branches and the `|| echo "Warning..."`; derive the channel once with `release-tools.js channel` and publish every `npm/*` directory with `release-tools.js publish`, then the root package last, all under `set -euo pipefail`.

**Acceptance Check:**
`f=.github/workflows/release.yml; grep -q 'needs: verify' $f && grep -q 'id-token: write' $f && grep -q 'check-version' $f && grep -q 'release-tools.js publish' $f && grep -q 'release-notes' $f && ! grep -q 'Warning: Failed' $f && ! grep -q '"alpha"\|\*"beta"\*' $f && echo OK` prints `OK`.

- [x] T2.3: [CONFIG] Re-runnable GitHub release with reproducible binaries
**Target:** `.goreleaser.yaml`
**Context:** [US-4]

**Constraints:**
- Sequential after T2.1 (same file). Add `release.replace_existing_artifacts: true` so a re-run replaces assets already uploaded to the tag's release instead of failing.
- Make builds reproducible: add `-trimpath` to the build `flags` and `mod_timestamp: "{{ .CommitTimestamp }}"` to the build; keep the existing ldflags unchanged.

**Acceptance Check:**
`grep -q 'replace_existing_artifacts: true' .goreleaser.yaml && grep -q 'trimpath' .goreleaser.yaml && grep -q 'CommitTimestamp' .goreleaser.yaml && goreleaser check >/dev/null 2>&1 && echo OK` prints `OK`.

### Phase 3: Continuous integration on every version line

- [x] T3.1: [CONFIG] Version line branches, Windows tests, single-run gates
**Target:** `.github/workflows/ci.yml`
**Context:** [US-1]

**Constraints:**
- `push` and `pull_request` both list `main` and `'*.x'`.
- The `test` matrix is `ubuntu-latest`, `macos-latest`, `windows-latest`; on every OS the job runs the unit tests, while the coverage profile and the "Check Coverage" step run only when `matrix.os == 'ubuntu-latest'` (the step relies on bash, awk and bc).
- Add a Linux-only step or job that sets up Node 20 and runs `node --test scripts/`.

**Acceptance Check:**
`f=.github/workflows/ci.yml; grep -q "'\*\.x'" $f && grep -q 'windows-latest' $f && grep -q "matrix.os == 'ubuntu-latest'" $f && grep -q 'node --test scripts/' $f && echo OK` prints `OK`.

## 3. Black-box Scenarios

- **S1:** Mismatched release version is refused
**Context:** [US-2]
**Command:** `node scripts/release-tools.js check-version v9.9.9; echo "exit=$?"`
**Expected:** A single stderr line naming the tag `9.9.9`, the `package.json` version and the built-in version, followed by `exit=1`.

- **S2:** Channel of stable and pre-release versions
**Context:** [US-3]
**Command:** `for v in 1.0.0 1.0.0-rc.1 1.1.0-next.1; do node scripts/release-tools.js channel $v; done`
**Expected:** Prints `latest`, `rc` and `next` on three lines.

- **S3:** Numeric pre-release has no channel
**Context:** [US-3]
**Command:** `node scripts/release-tools.js channel 1.1.0-1; echo "exit=$?"`
**Expected:** A stderr line asking for a named pre-release channel, followed by `exit=1`.

- **S4:** Missing notes file falls back to the generated list
**Context:** [US-6]
**Command:** `out=$(node scripts/release-tools.js notes v0.0.0-none); echo "[$out] exit=$?"`
**Expected:** `[] exit=0`.

## 4. Pre-emptive Mitigations
- **Risk:** The npm "not found" response of `npm view` changes wording and an existing package is treated as absent, making `npm publish` fail on re-run. -> **Mitigation:** treat only the `E404` error code as absent; anything else fails the run, which is safe (never publishes twice, never skips silently).
- **Risk:** A re-run rebuilds the binaries after some were already published to npm; non-reproducible builds would make the replaced GitHub assets and checksums differ from the npm binaries. -> **Mitigation:** reproducible builds (T2.3), verified by building the same commit twice and comparing hashes.
- **Risk:** Provenance requires a public repository and a GitHub-hosted runner; a misconfiguration only shows up at release time. -> **Mitigation:** the maintainer's `v1.0.0-rc.1` rehearsal validates it before the stable tag.
