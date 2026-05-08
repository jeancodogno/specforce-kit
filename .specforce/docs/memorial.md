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
- [2026-05-08] [Security] Security Hardening v2026-05-08 implemented. Upgraded to Go 1.26.3 to fix critical net/http vulnerabilities. Enforced 0600 file permissions and removed global gosec exclusions in favor of granular #nosec justifications.
- [2026-05-05] [Agent/Protocol] SDD Trigger Optimization implemented. Agents now have explicit mandates for triggering `spf.discovery` (vague intent/bug investigation) and `spf.spec` (planning). Direct edits to codebase without an approved spec are now strictly prohibited in `AGENTS.md`.
- [2026-05-05] [Project/Init] Conditional tool folder creation implemented. `.gemini`, `.claude`, and `.agent` directories are now only created if selected or already existing.

## Last Actions
- **Date:** 2026-05-08
- **Scope:** Security Hardening (v2026-05-08)
- **Completed:** Upgraded Go to 1.26.3 across toolchain and CI. Hardened file permissions for sensitive metadata (upgrade state, spec metadata, AGENTS.md) to 0600. Refined `gosec` scanning by removing global exclusions and implementing granular `#nosec G304/G306/G703` justifications linked to `core.SecurePath` and `os.OpenRoot`.
- **Relevant Files:** go.mod, Makefile, src/internal/upgrade/state.go, src/internal/spec/metadata.go, src/internal/project/agents_md.go, src/internal/upgrade/service.go, src/internal/spec/tasks.go

- **Date:** 2026-05-07
- **Scope:** Automated Background Updates (Self-Upgrade Engine)
...
## Active Lessons & Anti-Patterns
- **First Seen:** 2026-05-08
- **Last Seen:** 2026-05-08
- **Scope:** Security / Static Analysis
- **Symptom:** Global security scan exclusions in `Makefile` hide new vulnerabilities and reduce visibility into the project's security posture.
- **Avoid:** Using `-exclude` flags in `gosec` for broad categories like path traversal (G304).
- **Do Instead:** Use granular `// #nosec` comments at the specific code site with a clear justification (e.g., "Path validated by core.SecurePath").
- **Recurrence Count:** 1
- **Status:** Active
- **Distill To:** engineering.md

- **First Seen:** 2026-05-08
- **Last Seen:** 2026-05-08
- **Scope:** Security / Permissions
- **Symptom:** Standard `0644` permissions on configuration or state files allow other local users to read sensitive framework data.
- **Avoid:** Using default `0644` permissions for internal state files.
- **Do Instead:** Explicitly use `0600` for files and `0750` for directories created by the framework.
- **Recurrence Count:** 1
- **Status:** Active
- **Distill To:** security.md

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
