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

## Last Actions
- **Date:** 2026-04-28
- **Scope:** SDD Workflow / AI Proactivity
- **Completed:** Implemented "Autonomous SDD Workflow via Skills" featuring dual command-skill installation and `AGENTS.md` proactive triggers.
- **Next:** Validate agent proactivity in a fresh session.
- **Relevant Files:** src/internal/agent/translator.go, kit.yaml, AGENTS.md, src/internal/project/agents_md.go

## Active Lessons & Anti-Patterns
- **First Seen:** 2026-04-28
- **Last Seen:** 2026-04-28
- **Scope:** Agent / Tool Discovery
- **Symptom:** LLMs failing to recognize slash commands as procedural workflows when they aren't explicitly invoked.
- **Avoid:** Relying on user-initiated `/commands` for core SDD steps.
- **Do Instead:** Map commands to `skills/spf-*/SKILL.md` and instruct proactivity in `AGENTS.md`.
- **Recurrence Count:** 1
- **Status:** Resolved
- **Distill To:** governance.md

## Pending Decisions (Need Distillation)
- **Date:** 2026-04-28
- **Scope:** Architecture / Kit
- **Decision:** Support `MappingConfigs` as a slice in `kit.yaml` to allow multiple destinations for the same category.
- **Why:** Enables "Commands as Skills" without logic duplication.
- **Validate By:** `src/internal/agent/dual_install_test.go`
- **Distill To:** architecture.md
