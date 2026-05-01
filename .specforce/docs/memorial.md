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
- [2026-04-28] [Autonomous Workflows] Proactive Mandate in `AGENTS.md` and dual-installation (Commands as Skills) are now live. Agents should automatically trigger `/spec` or `/implement`.
- [2026-04-30] [Worktrees] Multi-root Git Worktree Support added to Scanner and Console. Specs are now aggregated across branches.
- [2026-04-30] [Environment] Global bin PATH issues resolved with automated builds and stylized diagnostics in `index.js`.

## Last Actions
- **Date:** 2026-04-30
- **Scope:** Fix Command Not Found (Issue #3)
- **Completed:** Added `prepare`/`postinstall` scripts to `package.json`, refactored `index.js` with a stylized diagnostic tool for PATH issues, and updated troubleshooting docs.
- **Next:** Monitor for similar reports in other OS environments.
- **Relevant Files:** index.js, package.json, README.md, docs/getting-started.md

- **Date:** 2026-04-30
- **Scope:** Git Worktree Support
- **Completed:** Implemented porcelain-based worktree discovery, updated scanner to aggregate Active/Archived specs across roots, and added UI labels in Console.
- **Next:** Execute `/archive` to formalize worktree support.
- **Relevant Files:** src/internal/spec/scanner.go, src/internal/tui/console.go, src/internal/spec/scanner_test.go

## Active Lessons & Anti-Patterns
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
