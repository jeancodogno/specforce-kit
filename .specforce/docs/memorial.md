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

---

## Critical Now
- [2026-05-02] [Security] "Zero Scripts" policy implemented. `postinstall` and `prepare` scripts removed from `package.json` to mitigate supply chain risk alerts.
- [2026-04-30] [Worktrees] Multi-root Git Worktree Support added to Scanner and Console. Specs are now aggregated across branches.
- [2026-04-30] [Environment] Global bin PATH issues resolved with stylized diagnostics in `index.js`.

## Last Actions
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
- **First Seen:** 2026-05-01
- **Last Seen:** 2026-05-01
- **Scope:** CI / Platform
- **Symptom:** `TestScanProject_WithWorktrees` failing on macOS but passing on Linux due to `/var` vs `/private/var` symlink.
- **Avoid:** Direct comparison of `filepath.Abs` results across different environments without evaluating symlinks.
- **Do Instead:** Use a helper (e.g., `evalPath`) that applies `filepath.EvalSymlinks` to canonicalize paths before equality checks.
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
