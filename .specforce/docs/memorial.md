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
- [2026-04-26] [Linter] `.golangci.yml` requires `version: "2"` for compatibility with current `golangci-lint` versions.

## Last Actions
- **Date:** 2026-04-26
- **Scope:** Agent / Project Setup
- **Completed:** Corrected `AGENTS.md` hook names and implemented automated platform-specific configurations (Gemini settings, Antigravity/Claude symlinks).
- **Next:** Monitor agent discovery performance across different platforms with the new automated links.
- **Relevant Files:** src/internal/project/agents_md.go, docs/configuration.md, docs/supported-tools.md

## Active Lessons & Anti-Patterns
- **First Seen:** 2026-04-26
- **Last Seen:** 2026-04-26
- **Scope:** Project Initialization / AI Onboarding
- **Symptom:** AI agents (Gemini, Claude Code) failing to discover project rules after `specforce init`.
- **Avoid:** Manual setup instructions for each platform.
- **Do Instead:** Automated environment-specific configuration (symlinks/settings files) triggered during `EnsureAgentsMD`.
- **Recurrence Count:** 1
- **Status:** Resolved
- **Distill To:** engineering.md

## Pending Decisions (Need Distillation)
- **Date:**
- **Scope:**
- **Decision:**
- **Why:**
- **Validate By:**
- **Distill To:**
