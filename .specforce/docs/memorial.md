# Memorial (Active Memory)

> **FOR AI AGENTS: RULES OF ENGAGEMENT**
> This is your cross-session memory file. You MUST obey these rules:
> 1. **READ THIS FIRST** before starting any task, debugging, or writing code in this repository.
> 2. **TRIAGE:** Read `Critical Now` first, then `Last Actions`. Only read lessons if they match your current scope.
> 3. **STRICT LIMITS (DO NOT BLOAT):**
>    - Max 3 `Critical Now` bullets.
>    - Max 5 `Last Actions` entries.
>    - Max 7 `Active Lessons`.
>    - Max 5 `Pending Decisions`.
> 4. **GARBAGE COLLECTION:** If a list is full, you MUST delete the oldest/least relevant item before adding a new one.
> 5. **DISTILLATION:** If a lesson or decision becomes a permanent rule, you MUST move it to the official Constitution files (`engineering.md`, `architecture.md`, etc.) and completely delete it from here.

## Critical Now
- [2026-05-04] [CLI] Auto-timestamped spec slugs implemented. `specforce spec init` now prepends `YYYYMMDD-HHMM-` to slugs automatically.
- [2026-05-04] [Docs] Multi-language documentation structure (EN, PT, ES) implemented. All docs migrated to `docs/{en,pt,es}/`.
- [2026-05-02] [Security] "Zero Scripts" policy implemented. `postinstall` and `prepare` scripts removed from `package.json`.

## Last Actions
- **Date:** 2026-05-04
- **Scope:** Auto-Timestamped Spec Slugs
- **Completed:** Implemented `PrepareSlug` logic in `spec` package, integrated into `HandleSpecInit`, and added unit/integration tests. Created new standard for spec directory naming.
- **Next:** Monitor for any issues with nested directory timestamping in CI.
- **Relevant Files:** src/internal/spec/slug.go, src/internal/cli/spec.go, src/internal/cli/integration_test.go

- **Date:** 2026-05-04
- **Scope:** Multi-language Documentation Expansion
...
## Active Lessons & Anti-Patterns
- **First Seen:** 2026-05-04
- **Last Seen:** 2026-05-04
- **Scope:** CLI / UX
- **Symptom:** Users might not notice if a slug is transformed (timestamped) unless explicitly told.
- **Avoid:** Transforming user input silently without feedback.
- **Do Instead:** Always print the final resolved path/slug in the success message of `spec init` to maintain transparency.
- **Recurrence Count:** 1
- **Status:** Active
- **Distill To:** engineering.md

- **Completed:** Restructured `docs/` folder into `en/`, `pt/`, and `es/`. Migrated English docs, added language selector to root README, and created PT/ES localized versions of root documents.
- **Next:** Continue translating specific documentation files within the language folders.
- **Relevant Files:** README.md, README.pt.md, README.es.md, docs/{en,pt,es}/*

- **Date:** 2026-05-02
- **Scope:** Security Hardening (Socket.dev Alerts)
- **Completed:** Removed `postinstall`/`prepare` scripts, hardened `index.js` proxy with absolute path validation, and updated `security.md`.
- **Next:** Monitor for any user confusion regarding the lack of automatic builds.
- **Relevant Files:** index.js, package.json, .specforce/docs/security.md

- **Date:** 2026-05-01
- **Scope:** macOS CI & Coverage
- **Completed:** Fixed `ScanProject` path discrepancies using `filepath.EvalSymlinks`, added resilience to `os.UserHomeDir` tests, and increased `upgrade` coverage to >80%.
- **Next:** Monitor macOS-specific path issues in future features.
- **Relevant Files:** src/internal/spec/scanner.go, src/internal/agent/translator_test.go, src/internal/upgrade/*.go

## Active Lessons & Anti-Patterns
- **First Seen:** 2026-05-04
- **Last Seen:** 2026-05-04
- **Scope:** Documentation / Migration
- **Symptom:** Links breaking during documentation folder restructuring.
- **Avoid:** Moving documentation files without a recursive link verification step.
- **Do Instead:** Use a systematic approach to update all relative paths in `README.md` and internal `.md` files immediately after moving files to nested subdirectories.
- **Recurrence Count:** 1
- **Status:** Active
- **Distill To:** engineering.md

- **First Seen:** 2026-05-01
- **Last Seen:** 2026-05-01
- **Scope:** Testing / Platform
- **Symptom:** Tests relying on `os.Unsetenv("HOME")` to force `os.UserHomeDir` failure fail on macOS because of system fallback.
- **Avoid:** Assuming `os.UserHomeDir` will always fail if environment variables are removed.
- **Do Instead:** Explicitly check if the function still succeeds after `Unsetenv` and use `t.Skip` to handle platforms with persistent home directory resolution.
- **Recurrence Count:** 1
- **Status:** Active
- **Distill To:** engineering.md
- **First Seen:** 2026-04-30
- **Last Seen:** 2026-04-30
- **Scope:** Proxy / Environment
- **Symptom:** Users encounter "command not found" even after successful global install if npm bin is missing from PATH.
- **Avoid:** Silently failing or giving generic errors when native binaries are missing.
- **Do Instead:** Run an automated diagnostic check that identifies the exact missing PATH entry and provides the specific OS-level fix command.
- **Recurrence Count:** 1
- **Status:** Active
- **Distill To:** engineering.md

## Pending Decisions (Need Distillation)
- **Date:** 2026-04-28
- **Scope:** Architecture / Kit
- **Decision:** Support `MappingConfigs` as a slice in `kit.yaml` to allow multiple destinations for the same category.
- **Why:** Enables "Commands as Skills" without logic duplication.
- **Validate By:** `src/internal/agent/dual_install_test.go`
- **Distill To:** architecture.md
