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
- [2026-05-05] [Agent/Protocol] SDD Trigger Optimization implemented. Agents now have explicit mandates for triggering `spf.discovery` (vague intent/bug investigation) and `spf.spec` (planning). Direct edits to codebase without an approved spec are now strictly prohibited in `AGENTS.md`.
- [2026-05-05] [Project/Init] Conditional tool folder creation implemented. `.gemini`, `.claude`, and `.agent` directories are now only created if selected or already existing.
- [2026-05-04] [Agent/Kit] Agent/Skill reorganization (v1.x) implemented. Dynamic Registry with multi-source discovery (Embedded + Local) and variable injection (`{{context}}`) are now active.

## Last Actions
- **Date:** 2026-05-07
- **Scope:** Automated Background Updates (Self-Upgrade Engine)
- **Completed:** Implemented a two-phase background update system (Check/Stage -> Swap). Verified detached background processes, atomic binary swap via 'Rename-to-Old' pattern, and `syscall.Exec` process replacement. Integrated GitHub release provider with checksum verification and 6h throttling.
- **Relevant Files:** src/internal/upgrade/*, src/cmd/specforce/main.go, src/internal/cli/cobra/root.go

- **Date:** 2026-05-06
- **Scope:** Shields.io & Navigation Integration
...
## Active Lessons & Anti-Patterns
- **First Seen:** 2026-05-07
- **Last Seen:** 2026-05-07
- **Scope:** Upgrade / Filesystem
- **Symptom:** `os.Rename` fails when staging a binary from `/tmp` to `~/.specforce` if they reside on different partitions (Invalid cross-device link).
- **Avoid:** Assuming `os.Rename` works globally across the filesystem.
- **Do Instead:** Use a robust `moveFile` helper that attempts `os.Rename` first and falls back to `io.Copy` + `os.Remove` on failure.
- **Recurrence Count:** 1
- **Status:** Active
- **Distill To:** engineering.md

- **First Seen:** 2026-05-07
- **Last Seen:** 2026-05-07
- **Scope:** Upgrade / Process Lifecycle
- **Symptom:** `syscall.Exec` fails or executes the wrong file if `os.Executable()` is called AFTER the binary has been renamed/swapped on disk.
- **Avoid:** Resolving the executable path after mutation.
- **Do Instead:** Capture the absolute path via `os.Executable()` at the very start of the swap operation to ensure the memory-resident process points to the correct entry point for replacement.
- **Recurrence Count:** 1
- **Status:** Active
- **Distill To:** architecture.md

- **First Seen:** 2026-05-06
- **Last Seen:** 2026-05-06
...
## Pending Decisions (Need Distillation)
- **Date:** 2026-05-07
- **Scope:** Architecture / Lifecycle
- **Decision:** Establish the "Atomic Binary Swap" (Rename-to-Old + syscall.Exec) as the standard for CLI self-updates.
- **Why:** Ensures zero-downtime, prevents binary corruption during download, and provides a safe rollback path (.old) if the new version fails to start.
- **Validate By:** `src/internal/upgrade/integration_test.go`
- **Distill To:** architecture.md
