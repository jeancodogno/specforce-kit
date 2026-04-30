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
- [2026-04-19] [Tasks] Hierarchical task organization (Phases H3 / Tasks H4) implemented. Agents should now use this for all new feature roadmaps.
- [2026-04-28] [Autonomous Workflows] Proactive Mandate in `AGENTS.md` and dual-installation (Commands as Skills) are now live. Agents should automatically trigger `/spec` or `/implement`.
- [2026-04-30] [Worktrees] Multi-root Git Worktree Support added to Scanner and Console. Specs are now aggregated across branches.

## Last Actions
- **Date:** 2026-04-30
- **Scope:** Git Worktree Support
- **Completed:** Implemented porcelain-based worktree discovery, updated scanner to aggregate Active/Archived specs across roots, and added UI labels in Console.
- **Next:** Execute `/archive` to formalize worktree support.
- **Relevant Files:** src/internal/spec/scanner.go, src/internal/tui/console.go, src/internal/spec/scanner_test.go

- **Date:** 2026-04-29
- **Scope:** Discovery Command
- **Completed:** Implemented `spf.discovery` (Specforce Scout), updated `AGENTS.md` template, and enriched prompt with unique branding and diagnostic workflows.
- **Next:** Monitor user engagement with Discovery mode to refine 'Thinking Partner' prompts.
- **Relevant Files:** src/internal/agent/kit/commands/discovery.yaml, src/internal/project/agents_md.go, README.md

## Active Lessons & Anti-Patterns
- **First Seen:** 2026-04-30
- **Last Seen:** 2026-04-30
- **Scope:** Scanner / Performance
- **Symptom:** Aggregating specs from dozens of worktrees could cause disk I/O bottlenecks or UI lag if scanned synchronously.
- **Avoid:** Synchronous scanning of potentially many external project roots during UI refresh.
- **Do Instead:** Use `git worktree list --porcelain` for fast path discovery and implement a robust exclusion logic (e.g., skip Constitution for external roots) to minimize I/O.
- **Recurrence Count:** 1
- **Status:** Active
- **Distill To:** architecture.md

## Pending Decisions (Need Distillation)
- **Date:** 2026-04-28
- **Scope:** Architecture / Kit
- **Decision:** Support `MappingConfigs` as a slice in `kit.yaml` to allow multiple destinations for the same category.
- **Why:** Enables "Commands as Skills" without logic duplication.
- **Validate By:** `src/internal/agent/dual_install_test.go`
- **Distill To:** architecture.md
